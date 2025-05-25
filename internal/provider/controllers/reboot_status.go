// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"errors"
	"time"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/generic/qtransform"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/siderolabs/gen/xerrors"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/pxe"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/machine"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/util"
)

// RebootStatusControllerOptions defines options for the RebootStatusController.
type RebootStatusControllerOptions struct {
	PostTransformFunc func()
}

// RebootStatusController manages machine power management.
type RebootStatusController = qtransform.QController[*infra.Machine, *resources.RebootStatus]

// NewRebootStatusController creates a new RebootStatusController.
func NewRebootStatusController(bmcClientFactory BMCClientFactory, minRebootInterval time.Duration, pxeBootMode pxe.BootMode, options RebootStatusControllerOptions) *RebootStatusController {
	controllerName := meta.ProviderID.String() + ".RebootStatusController"

	helper := &rebootStatusControllerHelper{
		bmcClientFactory:  bmcClientFactory,
		minRebootInterval: minRebootInterval,
		pxeBootMode:       pxeBootMode,
		controllerName:    controllerName,
		options:           options,
	}

	return qtransform.NewQController(
		qtransform.Settings[*infra.Machine, *resources.RebootStatus]{
			Name: controllerName,
			MapMetadataFunc: func(infraMachine *infra.Machine) *resources.RebootStatus {
				return resources.NewRebootStatus(infraMachine.Metadata().ID())
			},
			UnmapMetadataFunc: func(rebootStatus *resources.RebootStatus) *infra.Machine {
				return infra.NewMachine(rebootStatus.Metadata().ID())
			},
			TransformExtraOutputFunc:        helper.transform,
			FinalizerRemovalExtraOutputFunc: helper.finalizerRemoval,
		},
		qtransform.WithExtraMappedInput(qtransform.MapperSameID[*resources.MachineStatus, *infra.Machine]()),
		qtransform.WithExtraMappedInput(qtransform.MapperSameID[*resources.BMCConfiguration, *infra.Machine]()),
		qtransform.WithExtraMappedInput(qtransform.MapperSameID[*resources.WipeStatus, *infra.Machine]()),
		qtransform.WithExtraMappedInput(qtransform.MapperSameID[*resources.PowerOperation, *infra.Machine]()),
		qtransform.WithConcurrency(4),
	)
}

type rebootStatusControllerHelper struct {
	bmcClientFactory  BMCClientFactory
	options           RebootStatusControllerOptions
	pxeBootMode       pxe.BootMode
	controllerName    string
	minRebootInterval time.Duration
}

func (helper *rebootStatusControllerHelper) transform(ctx context.Context, r controller.ReaderWriter, logger *zap.Logger, infraMachine *infra.Machine, rebootStatus *resources.RebootStatus) error {
	if helper.options.PostTransformFunc != nil {
		defer helper.options.PostTransformFunc()
	}

	machineStatus, err := handleInput[*resources.MachineStatus](ctx, r, helper.controllerName, infraMachine)
	if err != nil {
		return err
	}

	bmcConfiguration, err := handleInput[*resources.BMCConfiguration](ctx, r, helper.controllerName, infraMachine)
	if err != nil {
		return err
	}

	wipeStatus, err := handleInput[*resources.WipeStatus](ctx, r, helper.controllerName, infraMachine)
	if err != nil {
		return err
	}

	powerOperation, err := handleInput[*resources.PowerOperation](ctx, r, helper.controllerName, infraMachine)
	if err != nil {
		return err
	}

	if err = validateInfraMachine(infraMachine, logger); err != nil {
		return err
	}

	if bmcConfiguration == nil {
		logger.Debug("bmc configuration not found, skip")

		return xerrors.NewTaggedf[qtransform.SkipReconcileTag]("bmc configuration not found")
	}

	if machineStatus == nil {
		logger.Debug("machine status not found, skip")

		return xerrors.NewTaggedf[qtransform.SkipReconcileTag]("machine status not found")
	}

	requiredBootMode := machine.RequiredBootMode(infraMachine, bmcConfiguration, wipeStatus, logger)
	requiresPXEBoot := machine.RequiresPXEBoot(requiredBootMode)
	requiresPowerOn := machine.RequiresPowerOn(infraMachine, wipeStatus)
	agentAccessible := machineStatus.TypedSpec().Value.AgentAccessible

	// We decide if a reboot is required or not by checking if any of the following is true:
	// 1. If the machine is required to be in agent mode and the agent is not accessible.
	// 2. If the machine is required to be PXE booted Talos, but the agent is still accessible.
	//
	// If the machine, however, is required to be booted from the disk, we never issue a reboot here due to necessity, as it is Omni's responsibility.
	modeMismatch := (requiredBootMode == machine.BootModeAgentPXE && !agentAccessible) ||
		(requiredBootMode == machine.BootModeTalosPXE && agentAccessible)

	// if the machine does not need to be powered on, we should not reboot it,so we avoid reboot loops.
	requiresReboot := requiresPowerOn && modeMismatch

	logger = logger.With(
		zap.Bool("requires_reboot", requiresReboot),
		zap.Bool("requires_pxe_boot", requiresPXEBoot),
		zap.String("required_boot_mode", string(requiredBootMode)),
		zap.Bool("required_power_on", requiresPowerOn),
		zap.Bool("agent_accessible", agentAccessible),
		zap.Bool("mode_mismatch", modeMismatch),
	)

	if requiresReboot {
		logger.Info("reboot machine to switch boot mode")

		return helper.reboot(ctx, infraMachine, bmcConfiguration, powerOperation, requiresPXEBoot, rebootStatus, logger)
	}

	if rebootStatus.TypedSpec().Value.LastRebootId != infraMachine.TypedSpec().Value.RequestedRebootId {
		logger.Debug("reboot machine by user request")

		return helper.reboot(ctx, infraMachine, bmcConfiguration, powerOperation, requiresPXEBoot, rebootStatus, logger)
	}

	return nil
}

func (helper *rebootStatusControllerHelper) finalizerRemoval(ctx context.Context, r controller.ReaderWriter, logger *zap.Logger, infraMachine *infra.Machine) (retErr error) {
	defer func() {
		retErr = errors.Join(retErr, helper.removeFinalizers(ctx, r, infraMachine))
	}()

	logger.Info("attempt to reboot the removed infra machine")

	// machine is removed, we try our best to get into the agent mode
	machineStatus, err := safe.ReaderGetByID[*resources.MachineStatus](ctx, r, infraMachine.Metadata().ID())
	if err != nil {
		if state.IsNotFoundError(err) {
			logger.Debug("machine status not found, skip")

			return nil
		}

		return err
	}

	if machineStatus.TypedSpec().Value.AgentAccessible {
		logger.Info("agent is accessible, no need to reboot")

		return nil
	}

	bmcConfiguration, err := safe.ReaderGetByID[*resources.BMCConfiguration](ctx, r, infraMachine.Metadata().ID())
	if err != nil {
		if state.IsNotFoundError(err) {
			logger.Info("bmc configuration does not exist, skip reboot")

			return nil
		}

		return err
	}

	if err = helper.reboot(ctx, infraMachine, bmcConfiguration, nil, true, nil, logger); err != nil {
		logger.Error("failed to reboot the removed machine", zap.Error(err))
	} else {
		logger.Info("rebooted the removed infra machine")
	}

	return nil
}

func (helper *rebootStatusControllerHelper) removeFinalizers(ctx context.Context, r controller.ReaderWriter, infraMachine *infra.Machine) error {
	if _, err := handleInput[*resources.MachineStatus](ctx, r, helper.controllerName, infraMachine); err != nil {
		return err
	}

	if _, err := handleInput[*resources.BMCConfiguration](ctx, r, helper.controllerName, infraMachine); err != nil {
		return err
	}

	if _, err := handleInput[*resources.WipeStatus](ctx, r, helper.controllerName, infraMachine); err != nil {
		return err
	}

	if _, err := handleInput[*resources.PowerOperation](ctx, r, helper.controllerName, infraMachine); err != nil {
		return err
	}

	return nil
}

func (helper *rebootStatusControllerHelper) reboot(ctx context.Context,
	infraMachine *infra.Machine, bmcConfiguration *resources.BMCConfiguration,
	powerOperation *resources.PowerOperation, requiresPXEBoot bool, rebootStatus *resources.RebootStatus, logger *zap.Logger,
) error {
	// check if we are in the cooldown period
	timeSinceLastPowerOn := getTimeSinceLastPowerOn(powerOperation, rebootStatus)
	if timeSinceLastPowerOn < helper.minRebootInterval {
		logger.Debug("we are in reboot cooldown period, requeue", zap.Duration("elapsed", timeSinceLastPowerOn), zap.Duration("min_reboot_interval", helper.minRebootInterval))

		return controller.NewRequeueInterval(helper.minRebootInterval - timeSinceLastPowerOn + time.Second)
	}

	bmcClient, err := helper.bmcClientFactory.GetClient(ctx, bmcConfiguration, logger)
	if err != nil {
		return err
	}

	defer util.LogClose(bmcClient, logger)

	if requiresPXEBoot {
		if err = bmcClient.SetPXEBootOnce(ctx, helper.pxeBootMode); err != nil {
			return err
		}
	}

	if err = bmcClient.Reboot(ctx); err != nil {
		return err
	}

	if rebootStatus != nil {
		rebootStatus.TypedSpec().Value.LastRebootId = infraMachine.TypedSpec().Value.RequestedRebootId
		rebootStatus.TypedSpec().Value.LastRebootTimestamp = timestamppb.Now()
	}

	return nil
}
