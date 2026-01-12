// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"fmt"
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
//
//nolint:dupl,cyclop
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
		qtransform.WithExtraMappedInput[*resources.MachineStatus](qtransform.MapperSameID[*infra.Machine]()),
		qtransform.WithExtraMappedInput[*resources.BMCConfiguration](qtransform.MapperSameID[*infra.Machine]()),
		qtransform.WithExtraMappedInput[*resources.WipeStatus](qtransform.MapperSameID[*infra.Machine]()),
		qtransform.WithExtraMappedInput[*resources.PowerOperation](qtransform.MapperSameID[*infra.Machine]()),
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

//nolint:gocyclo,cyclop
func (helper *rebootStatusControllerHelper) transform(ctx context.Context, r controller.ReaderWriter, logger *zap.Logger, infraMachine *infra.Machine, rebootStatus *resources.RebootStatus) error {
	if helper.options.PostTransformFunc != nil {
		defer helper.options.PostTransformFunc()
	}

	machineStatus, err := safe.ReaderGetByID[*resources.MachineStatus](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	if machineStatus != nil && !machineStatus.Metadata().Finalizers().Has(helper.controllerName) {
		if err = r.AddFinalizer(ctx, machineStatus.Metadata(), helper.controllerName); err != nil {
			return fmt.Errorf("failed to add finalizer to machine status: %w", err)
		}
	}

	bmcConfiguration, err := safe.ReaderGetByID[*resources.BMCConfiguration](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	if bmcConfiguration != nil && !bmcConfiguration.Metadata().Finalizers().Has(helper.controllerName) {
		if err = r.AddFinalizer(ctx, bmcConfiguration.Metadata(), helper.controllerName); err != nil {
			return fmt.Errorf("failed to add finalizer to bmc configuration: %w", err)
		}
	}

	wipeStatus, err := safe.ReaderGetByID[*resources.WipeStatus](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	powerOperation, err := safe.ReaderGetByID[*resources.PowerOperation](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
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
		return helper.reboot(ctx, infraMachine, bmcConfiguration, powerOperation, requiresPXEBoot, rebootStatus, logger)
	}

	if rebootStatus.TypedSpec().Value.LastRebootId != infraMachine.TypedSpec().Value.RequestedRebootId {
		logger.Debug("reboot machine by user request")

		return helper.reboot(ctx, infraMachine, bmcConfiguration, powerOperation, requiresPXEBoot, rebootStatus, logger)
	}

	return nil
}

func (helper *rebootStatusControllerHelper) finalizerRemoval(ctx context.Context, r controller.ReaderWriter, logger *zap.Logger, infraMachine *infra.Machine) error {
	// machine is removed, we try our best to get into the agent mode
	machineStatus, err := safe.ReaderGetByID[*resources.MachineStatus](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	bmcConfiguration, err := safe.ReaderGetByID[*resources.BMCConfiguration](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	helper.attemptRebootOnRemoval(ctx, machineStatus, bmcConfiguration, infraMachine, logger)

	if machineStatus != nil && machineStatus.Metadata().Finalizers().Has(helper.controllerName) {
		if err = r.RemoveFinalizer(ctx, machineStatus.Metadata(), helper.controllerName); err != nil {
			return fmt.Errorf("failed to remove finalizer from machine status: %w", err)
		}
	}

	if bmcConfiguration != nil && bmcConfiguration.Metadata().Finalizers().Has(helper.controllerName) {
		if err = r.RemoveFinalizer(ctx, bmcConfiguration.Metadata(), helper.controllerName); err != nil {
			return fmt.Errorf("failed to remove finalizer from bmc configuration: %w", err)
		}
	}

	return nil
}

func (helper *rebootStatusControllerHelper) attemptRebootOnRemoval(ctx context.Context, machineStatus *resources.MachineStatus, bmcConfiguration *resources.BMCConfiguration,
	infraMachine *infra.Machine, logger *zap.Logger,
) {
	logger.Info("attempt to reboot the removed infra machine")

	if bmcConfiguration == nil {
		logger.Warn("bmc configuration does not exist, skip reboot")

		return
	}

	if machineStatus == nil {
		logger.Warn("machine status does not exist, skip reboot")

		return
	}

	if machineStatus.TypedSpec().Value.AgentAccessible {
		logger.Info("agent is accessible, no need to reboot")

		return
	}

	if err := helper.reboot(ctx, infraMachine, bmcConfiguration, nil, true, nil, logger); err != nil {
		logger.Error("failed to reboot the removed infra machine", zap.Error(err))

		return
	}

	logger.Info("rebooted the removed infra machine")
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

	logger.Info("reboot machine to switch boot mode")

	bmcClient, err := helper.bmcClientFactory.GetClient(ctx, bmcConfiguration, logger)
	if err != nil {
		return err
	}

	defer util.LogCloseContext(ctx, bmcClient, logger)

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
