// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"strings"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/generic/qtransform"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
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
	helper := &infraMachineStatusControllerHelper{
		machineLabels: machineLabels,
	}

	return qtransform.NewQController(
		qtransform.Settings[*infra.Machine, *infra.MachineStatus]{
			Name: meta.ProviderID.String() + ".InfraMachineStatusController",
			MapMetadataFunc: func(infraMachine *infra.Machine) *infra.MachineStatus {
				return infra.NewMachineStatus(infraMachine.Metadata().ID())
			},
			UnmapMetadataFunc: func(infraMachineStatus *infra.MachineStatus) *infra.Machine {
				return infra.NewMachine(infraMachineStatus.Metadata().ID())
			},
			TransformFunc: helper.transform,
		},
		qtransform.WithConcurrency(4),
		qtransform.WithExtraMappedInput[*resources.MachineStatus](qtransform.MapperSameID[*infra.Machine]()),
		qtransform.WithExtraMappedInput[*resources.RebootStatus](qtransform.MapperSameID[*infra.Machine]()),
		qtransform.WithExtraMappedInput[*resources.WipeStatus](qtransform.MapperSameID[*infra.Machine]()),
		qtransform.WithExtraMappedInput[*resources.BMCConfiguration](qtransform.MapperSameID[*infra.Machine]()),
	)
}

type infraMachineStatusControllerHelper struct {
	machineLabels map[string]string
}

//nolint:gocyclo,cyclop
func (helper *infraMachineStatusControllerHelper) transform(ctx context.Context, r controller.Reader, logger *zap.Logger,
	infraMachine *infra.Machine, infraMachineStatus *infra.MachineStatus,
) error {
	machineStatus, err := safe.ReaderGetByID[*resources.MachineStatus](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	rebootStatus, err := safe.ReaderGetByID[*resources.RebootStatus](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	wipeStatus, err := safe.ReaderGetByID[*resources.WipeStatus](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	bmcConfiguration, err := safe.ReaderGetByID[*resources.BMCConfiguration](ctx, r, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
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

	for _, k := range []string{omni.LabelCluster, omni.LabelMachineSet, omni.LabelControlPlaneRole, omni.LabelWorkerRole} {
		if val, ok := infraMachine.Metadata().Labels().Get(k); ok {
			infraMachineStatus.Metadata().Labels().Set(k, val)
		} else {
			infraMachineStatus.Metadata().Labels().Delete(k)
		}
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
