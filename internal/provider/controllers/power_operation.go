// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"time"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/generic/qtransform"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/siderolabs/gen/xerrors"
	omnispecs "github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/pxe"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/machine"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/util"
)

// NowFunc is a function that returns the current time.
type NowFunc = func() time.Time

// PowerOperationController manages InfraMachine resource lifecycle.
type PowerOperationController = qtransform.QController[*infra.Machine, *resources.PowerOperation]

// NewPowerOperationController creates a new PowerOperationController.
//
//nolint:dupl
func NewPowerOperationController(nowFunc NowFunc, bmcClientFactory BMCClientFactory, minRebootInterval time.Duration, pxeBootMode pxe.BootMode) *PowerOperationController {
	helper := &powerOperationControllerHelper{
		nowFunc:           nowFunc,
		bmcClientFactory:  bmcClientFactory,
		minRebootInterval: minRebootInterval,
		pxeBootMode:       pxeBootMode,
	}

	return qtransform.NewQController(
		qtransform.Settings[*infra.Machine, *resources.PowerOperation]{
			Name: meta.ProviderID.String() + ".PowerOperationController",
			MapMetadataFunc: func(infraMachine *infra.Machine) *resources.PowerOperation {
				return resources.NewPowerOperation(infraMachine.Metadata().ID())
			},
			UnmapMetadataFunc: func(powerOperation *resources.PowerOperation) *infra.Machine {
				return infra.NewMachine(powerOperation.Metadata().ID())
			},
			TransformFunc: helper.transform,
		},
		qtransform.WithExtraMappedInput[*resources.BMCConfiguration](qtransform.MapperSameID[*infra.Machine]()),
		qtransform.WithExtraMappedInput[*resources.WipeStatus](qtransform.MapperSameID[*infra.Machine]()),
		qtransform.WithExtraMappedInput[*resources.RebootStatus](qtransform.MapperSameID[*infra.Machine]()),
		qtransform.WithExtraMappedInput[*resources.MachineStatus](qtransform.MapperSameID[*infra.Machine]()),
		qtransform.WithConcurrency(4),
	)
}

type powerOperationControllerHelper struct {
	bmcClientFactory  BMCClientFactory
	nowFunc           NowFunc
	pxeBootMode       pxe.BootMode
	minRebootInterval time.Duration
}

//nolint:gocyclo,cyclop
func (helper *powerOperationControllerHelper) transform(ctx context.Context, r controller.Reader, logger *zap.Logger,
	infraMachine *infra.Machine, powerOperation *resources.PowerOperation,
) error {
	bmcConfiguration, err := safe.ReaderGetByID[*resources.BMCConfiguration](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	wipeStatus, err := safe.ReaderGetByID[*resources.WipeStatus](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	// this controller needs to wake up after a reboot to bring the machine again to the desired power state
	rebootStatus, err := safe.ReaderGetByID[*resources.RebootStatus](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	if err = validateInfraMachine(infraMachine, logger); err != nil {
		return err
	}

	if bmcConfiguration == nil {
		logger.Debug("machine has no power management configuration")

		return xerrors.NewTaggedf[qtransform.SkipReconcileTag]("machine has no power management configuration")
	}

	requiresPowerOn := machine.RequiresPowerOn(infraMachine, wipeStatus)

	logger.Info("power operation",
		zap.Bool("installed", machine.IsInstalled(infraMachine, wipeStatus)),
		zap.Bool("allocated", infraMachine.TypedSpec().Value.ClusterTalosVersion != ""),
		zap.Bool("requires_wipe", machine.RequiresWipe(infraMachine, wipeStatus)),
		zap.Bool("requires_power_on", requiresPowerOn),
	)

	bmcClient, err := helper.bmcClientFactory.GetClient(ctx, bmcConfiguration, logger)
	if err != nil {
		return err
	}

	defer util.LogCloseContext(ctx, bmcClient, logger)

	isPoweredOn, err := bmcClient.IsPoweredOn(ctx)
	if err != nil {
		return err
	}

	preferredPowerState := infraMachine.TypedSpec().Value.PreferredPowerState

	logger = logger.With(zap.Bool("is_powered_on", isPoweredOn), zap.Stringer("preferred_power_state", preferredPowerState))

	switch {
	case !isPoweredOn && (requiresPowerOn || preferredPowerState == omnispecs.InfraMachineSpec_POWER_STATE_ON):
		logger.Debug("power on machine")

		requiredBootMode := machine.RequiredBootMode(infraMachine, bmcConfiguration, wipeStatus, logger)
		if machine.RequiresPXEBoot(requiredBootMode) {
			if err = bmcClient.SetPXEBootOnce(ctx, helper.pxeBootMode); err != nil {
				return err
			}
		}

		if err = bmcClient.PowerOn(ctx); err != nil {
			return err
		}

		powerOperation.TypedSpec().Value.LastPowerOperation = specs.PowerState_POWER_STATE_ON
		powerOperation.TypedSpec().Value.LastPowerOnTimestamp = timestamppb.New(helper.nowFunc())
	case isPoweredOn && (!requiresPowerOn && preferredPowerState == omnispecs.InfraMachineSpec_POWER_STATE_OFF):
		timeSinceLastPowerOn := getTimeSinceLastPowerOn(powerOperation, rebootStatus)
		if timeSinceLastPowerOn < helper.minRebootInterval {
			logger.Debug("we are in power off cooldown period, requeue", zap.Duration("elapsed", timeSinceLastPowerOn), zap.Duration("min_reboot_interval", helper.minRebootInterval))

			return controller.NewRequeueInterval(helper.minRebootInterval - timeSinceLastPowerOn + time.Second)
		}

		logger.Debug("power off machine")

		if err = bmcClient.PowerOff(ctx); err != nil {
			return err
		}

		powerOperation.TypedSpec().Value.LastPowerOperation = specs.PowerState_POWER_STATE_OFF
	default:
		logger.Debug("machine power state is already as desired")
	}

	return nil
}
