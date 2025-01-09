// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/generic/qtransform"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	omnispecs "github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/baremetal"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/boot"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/machinestatus"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
)

// PowerStatusController manages InfraMachine resource lifecycle.
type PowerStatusController = qtransform.QController[*infra.Machine, *baremetal.PowerStatus]

// NewPowerStatusController creates a new PowerStatusController.
func NewPowerStatusController(powerClientFactory PowerClientFactory, state state.State) *PowerStatusController {
	helper := &powerStatusControllerHelper{
		powerClientFactory: powerClientFactory,
		state:              state,
	}

	return qtransform.NewQController(
		qtransform.Settings[*infra.Machine, *baremetal.PowerStatus]{
			Name: meta.ProviderID + ".PowerStatusController",
			MapMetadataFunc: func(infraMachine *infra.Machine) *baremetal.PowerStatus {
				return baremetal.NewPowerStatus(infraMachine.Metadata().ID())
			},
			UnmapMetadataFunc: func(powerStatus *baremetal.PowerStatus) *infra.Machine {
				return infra.NewMachine(powerStatus.Metadata().ID())
			},
			TransformFunc: helper.transform,
		},
		qtransform.WithExtraMappedInput(
			func(_ context.Context, _ *zap.Logger, _ controller.QRuntime, status *baremetal.MachineStatus) ([]resource.Pointer, error) {
				ptr := infra.NewMachine(status.Metadata().ID()).Metadata()

				return []resource.Pointer{ptr}, nil
			},
		),
		qtransform.WithConcurrency(4),
	)
}

type powerStatusControllerHelper struct {
	powerClientFactory PowerClientFactory
	state              state.State
}

func (helper *powerStatusControllerHelper) transform(ctx context.Context, _ controller.Reader, logger *zap.Logger, machine *infra.Machine, powerStatus *baremetal.PowerStatus) error {
	machineStatus, err := machinestatus.Get(ctx, helper.state, machine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	machineState, err := safe.StateGetByID[*infra.MachineState](ctx, helper.state, machine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	mode, err := boot.DetermineRequiredMode(machine, machineStatus, machineState, logger)
	if err != nil {
		return err
	}

	if mode.RequiresPowerMgmtConfig {
		logger.Debug("power management is not yet configured")

		powerStatus.TypedSpec().Value.PowerState = specs.PowerState_POWER_STATE_UNKNOWN

		return nil
	}

	powerManagement := machineStatus.TypedSpec().Value.PowerManagement
	preferredPowerState := machine.TypedSpec().Value.PreferredPowerState

	powerClient, err := helper.powerClientFactory.GetClient(powerManagement)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := powerClient.Close(); closeErr != nil {
			logger.Error("failed to close power client", zap.Error(closeErr))
		}
	}()

	isPoweredOn, err := powerClient.IsPoweredOn(ctx)
	if err != nil {
		return err
	}

	logger = logger.With(zap.Bool("is_powered_on", isPoweredOn), zap.Bool("needs_to_be_powered_on", mode.NeedsToBePoweredOn), zap.Stringer("preferred_power_state", preferredPowerState))

	if machine.TypedSpec().Value.Cordoned {
		logger.Debug("machine is cordoned, skip power management")

		return helper.updatePowerState(ctx, powerStatus, isPoweredOn, false)
	}

	switch {
	case !isPoweredOn && (mode.NeedsToBePoweredOn || preferredPowerState == omnispecs.InfraMachineSpec_POWER_STATE_ON):
		logger.Debug("power on machine")

		if err = powerClient.PowerOn(ctx); err != nil {
			return err
		}

		return helper.updatePowerState(ctx, powerStatus, true, true)
	case isPoweredOn && (!mode.NeedsToBePoweredOn && preferredPowerState == omnispecs.InfraMachineSpec_POWER_STATE_OFF):
		logger.Debug("power off machine")

		if err = powerClient.PowerOff(ctx); err != nil {
			return err
		}

		return helper.updatePowerState(ctx, powerStatus, true, false)
	default:
		logger.Debug("machine power state is already as desired")

		return helper.updatePowerState(ctx, powerStatus, isPoweredOn, false)
	}
}

func (helper *powerStatusControllerHelper) updatePowerState(ctx context.Context, status *baremetal.PowerStatus, isPoweredOn, updateLastRebootTimestamp bool) error {
	powerState := specs.PowerState_POWER_STATE_OFF
	if isPoweredOn {
		powerState = specs.PowerState_POWER_STATE_ON
	}

	if _, err := machinestatus.Modify(ctx, helper.state, status.Metadata().ID(), func(machineStatus *baremetal.MachineStatus) error {
		machineStatus.TypedSpec().Value.PowerState = powerState

		if updateLastRebootTimestamp {
			machineStatus.TypedSpec().Value.LastRebootTimestamp = timestamppb.Now()
		}

		return nil
	}); err != nil {
		return err
	}

	status.TypedSpec().Value.PowerState = powerState

	return nil
}
