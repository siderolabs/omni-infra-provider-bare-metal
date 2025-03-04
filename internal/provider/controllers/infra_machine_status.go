// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"strings"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/generic/qtransform"
	omnispecs "github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/machine"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
)

// InfraMachineStatusController manages InfraMachine resource lifecycle.
type InfraMachineStatusController = qtransform.QController[*infra.Machine, *infra.MachineStatus]

// NewInfraMachineStatusController initializes InfraMachineStatusController.
func NewInfraMachineStatusController(machineLabels map[string]string) *InfraMachineStatusController {
	controllerName := meta.ProviderID.String() + ".InfraMachineStatusController"
	helper := &infraMachineStatusControllerHelper{
		controllerName: controllerName,
		machineLabels:  machineLabels,
	}

	return qtransform.NewQController(
		qtransform.Settings[*infra.Machine, *infra.MachineStatus]{
			Name: controllerName,
			MapMetadataFunc: func(infraMachine *infra.Machine) *infra.MachineStatus {
				return infra.NewMachineStatus(infraMachine.Metadata().ID())
			},
			UnmapMetadataFunc: func(infraMachineStatus *infra.MachineStatus) *infra.Machine {
				return infra.NewMachine(infraMachineStatus.Metadata().ID())
			},
			TransformExtraOutputFunc:        helper.transform,
			FinalizerRemovalExtraOutputFunc: helper.finalizerRemoval,
		},
		qtransform.WithConcurrency(4),
		qtransform.WithExtraMappedInput(qtransform.MapperSameID[*resources.MachineStatus, *infra.Machine]()),
		qtransform.WithExtraMappedInput(qtransform.MapperSameID[*resources.RebootStatus, *infra.Machine]()),
		qtransform.WithExtraMappedInput(qtransform.MapperSameID[*resources.WipeStatus, *infra.Machine]()),
		qtransform.WithExtraMappedInput(qtransform.MapperSameID[*resources.BMCConfiguration, *infra.Machine]()),
	)
}

type infraMachineStatusControllerHelper struct {
	machineLabels  map[string]string
	controllerName string
}

func (helper *infraMachineStatusControllerHelper) transform(ctx context.Context, r controller.ReaderWriter, logger *zap.Logger,
	infraMachine *infra.Machine, infraMachineStatus *infra.MachineStatus,
) error {
	machineStatus, err := handleInput[*resources.MachineStatus](ctx, r, helper.controllerName, infraMachine)
	if err != nil {
		return err
	}

	rebootStatus, err := handleInput[*resources.RebootStatus](ctx, r, helper.controllerName, infraMachine)
	if err != nil {
		return err
	}

	wipeStatus, err := handleInput[*resources.WipeStatus](ctx, r, helper.controllerName, infraMachine)
	if err != nil {
		return err
	}

	bmcConfiguration, err := handleInput[*resources.BMCConfiguration](ctx, r, helper.controllerName, infraMachine)
	if err != nil {
		return err
	}

	if err = validateInfraMachine(infraMachine, logger); err != nil {
		return err
	}

	// we do not need to call validateInfraMachine here, as this controller does not
	// do any operations/modifications on the machine itself

	// clear existing labels
	for k := range infraMachineStatus.Metadata().Labels().Raw() {
		if strings.HasPrefix(k, omni.SystemLabelPrefix) {
			continue
		}

		infraMachineStatus.Metadata().Labels().Delete(k)
	}

	// set the new labels
	for k, v := range helper.machineLabels {
		infraMachineStatus.Metadata().Labels().Set(k, v)
	}

	if machineStatus != nil {
		switch machineStatus.TypedSpec().Value.PowerState {
		case specs.PowerState_POWER_STATE_UNKNOWN:
			infraMachineStatus.TypedSpec().Value.PowerState = omnispecs.InfraMachineStatusSpec_POWER_STATE_UNKNOWN
		case specs.PowerState_POWER_STATE_OFF:
			infraMachineStatus.TypedSpec().Value.PowerState = omnispecs.InfraMachineStatusSpec_POWER_STATE_OFF
		case specs.PowerState_POWER_STATE_ON:
			infraMachineStatus.TypedSpec().Value.PowerState = omnispecs.InfraMachineStatusSpec_POWER_STATE_ON
		}
	}

	if rebootStatus != nil {
		infraMachineStatus.TypedSpec().Value.LastRebootId = rebootStatus.TypedSpec().Value.LastRebootId
		infraMachineStatus.TypedSpec().Value.LastRebootTimestamp = rebootStatus.TypedSpec().Value.LastRebootTimestamp
	}

	installed := machine.IsInstalled(infraMachine, wipeStatus)
	requiresWipe := machine.RequiresWipe(infraMachine, wipeStatus)
	bmcConfigurationConfigured := bmcConfiguration != nil

	infraMachineStatus.TypedSpec().Value.Installed = installed
	infraMachineStatus.TypedSpec().Value.ReadyToUse = bmcConfigurationConfigured && !requiresWipe

	if wipeStatus != nil {
		infraMachineStatus.TypedSpec().Value.WipedNodeUniqueToken = wipeStatus.TypedSpec().Value.WipedNodeUniqueToken
	}

	logger.Debug("machine status",
		zap.Bool("installed", infraMachineStatus.TypedSpec().Value.Installed),
		zap.Bool("ready_to_use", infraMachineStatus.TypedSpec().Value.ReadyToUse))

	return nil
}

func (helper *infraMachineStatusControllerHelper) finalizerRemoval(ctx context.Context, r controller.ReaderWriter, _ *zap.Logger, infraMachine *infra.Machine) error {
	if _, err := handleInput[*resources.MachineStatus](ctx, r, helper.controllerName, infraMachine); err != nil {
		return err
	}

	if _, err := handleInput[*resources.RebootStatus](ctx, r, helper.controllerName, infraMachine); err != nil {
		return err
	}

	if _, err := handleInput[*resources.WipeStatus](ctx, r, helper.controllerName, infraMachine); err != nil {
		return err
	}

	if _, err := handleInput[*resources.BMCConfiguration](ctx, r, helper.controllerName, infraMachine); err != nil {
		return err
	}

	return nil
}
