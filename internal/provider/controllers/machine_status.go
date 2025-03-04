// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/generic"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	omniresources "github.com/siderolabs/omni/client/pkg/omni/resources"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/util"
)

// AgentConnectionEvent represents an agent connection/disconnection event.
type AgentConnectionEvent struct {
	MachineID resource.ID
	Connected bool
}

// PXEBootEvent represents a PXE boot event for a machine served by the provider.
type PXEBootEvent struct {
	MachineID resource.ID
}

// MachineStatusController manages machine status.
type MachineStatusController struct {
	bmcClientFactory       BMCClientFactory
	agentClient            AgentClient
	agentConnectionEventCh <-chan AgentConnectionEvent
	pxeBootEventCh         <-chan PXEBootEvent
	generic.NamedController
	pollInterval time.Duration
}

// NewMachineStatusController creates a new MachineStatusController.
func NewMachineStatusController(bmcClientFactory BMCClientFactory, agentClient AgentClient,
	agentConnectionEventCh <-chan AgentConnectionEvent, pxeBootEventCh <-chan PXEBootEvent, pollInterval time.Duration,
) *MachineStatusController {
	return &MachineStatusController{
		bmcClientFactory:       bmcClientFactory,
		agentClient:            agentClient,
		agentConnectionEventCh: agentConnectionEventCh,
		pxeBootEventCh:         pxeBootEventCh,
		pollInterval:           pollInterval,
		NamedController: generic.NamedController{
			ControllerName: meta.ProviderID.String() + ".MachineStatusController",
		},
	}
}

// Settings implements the controller.QController interface.
func (ctrl *MachineStatusController) Settings() controller.QSettings {
	return controller.QSettings{
		Inputs: []controller.Input{
			{
				Namespace: omniresources.InfraProviderNamespace,
				Type:      infra.InfraMachineType,
				Kind:      controller.InputQPrimary,
			},
			{
				Namespace: resources.Namespace(),
				Type:      resources.BMCConfigurationType(),
				Kind:      controller.InputQMapped,
			},
		},
		Outputs: []controller.Output{
			{
				Type: resources.MachineStatusType(),
				Kind: controller.OutputExclusive,
			},
		},
		RunHook: func(ctx context.Context, logger *zap.Logger, r controller.QRuntime) error {
			ticker := time.NewTicker(ctrl.pollInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return nil
				case <-ticker.C:
					ctrl.pollAll(ctx, r, logger)
				case event := <-ctrl.agentConnectionEventCh:
					if err := ctrl.handleAgentConnectionEvent(ctx, r, event, logger); err != nil {
						return err
					}
				case pxeBootEvent := <-ctrl.pxeBootEventCh:
					if err := ctrl.handlePXEBootEvent(ctx, r, pxeBootEvent, logger); err != nil {
						return err
					}
				}
			}
		},
	}
}

// Reconcile implements the controller.QController interface.
func (ctrl *MachineStatusController) Reconcile(ctx context.Context, _ *zap.Logger, r controller.QRuntime, ptr resource.Pointer) error {
	infraMachine, err := safe.ReaderGet[*infra.Machine](ctx, r, ptr)
	if err != nil {
		if !state.IsNotFoundError(err) {
			return err
		}

		// resource is not found, so we prepare a fake resource to trigger teardown logic
		infraMachine = infra.NewMachine(ptr.ID())
		infraMachine.Metadata().SetPhase(resource.PhaseTearingDown)
	}

	if _, err = handleInput[*resources.BMCConfiguration](ctx, r, ctrl.Name(), infraMachine); err != nil {
		return err
	}

	if infraMachine.Metadata().Phase() == resource.PhaseTearingDown {
		return ctrl.reconcileTearingDown(ctx, r, infraMachine)
	}

	return ctrl.reconcileRunning(ctx, r, infraMachine)
}

func (ctrl *MachineStatusController) reconcileTearingDown(ctx context.Context, r controller.QRuntime, infraMachine *infra.Machine) error {
	md := resources.NewMachineStatus(infraMachine.Metadata().ID()).Metadata()

	ready, err := r.Teardown(ctx, md)
	if err != nil {
		if state.IsNotFoundError(err) {
			return r.RemoveFinalizer(ctx, infraMachine.Metadata(), ctrl.Name())
		}

		return err
	}

	if !ready {
		return nil
	}

	if err = r.Destroy(ctx, md); err != nil {
		return err
	}

	return r.RemoveFinalizer(ctx, infraMachine.Metadata(), ctrl.Name())
}

func (ctrl *MachineStatusController) reconcileRunning(ctx context.Context, r controller.QRuntime, infraMachine *infra.Machine) error {
	if err := r.AddFinalizer(ctx, infraMachine.Metadata(), ctrl.Name()); err != nil {
		return err
	}

	return nil
}

// MapInput implements the controller.QController interface.
func (ctrl *MachineStatusController) MapInput(_ context.Context, _ *zap.Logger, _ controller.QRuntime, ptr resource.Pointer) ([]resource.Pointer, error) {
	switch ptr.Type() {
	case resources.BMCConfigurationType(), infra.InfraMachineType:
		return []resource.Pointer{infra.NewMachine(ptr.ID()).Metadata()}, nil
	}

	return nil, fmt.Errorf("unexpected resource type %q", ptr.Type())
}

func (ctrl *MachineStatusController) handleAgentConnectionEvent(ctx context.Context, r controller.QRuntime, event AgentConnectionEvent, logger *zap.Logger) error {
	logger.Info("handle agent connection event", zap.String("machine_id", event.MachineID), zap.Bool("connected", event.Connected))

	return safe.WriterModify(ctx, r, resources.NewMachineStatus(event.MachineID), func(res *resources.MachineStatus) error {
		if res.Metadata().Phase() == resource.PhaseTearingDown {
			return nil
		}

		res.TypedSpec().Value.AgentAccessible = event.Connected

		return nil
	})
}

func (ctrl *MachineStatusController) pollAll(ctx context.Context, r controller.QRuntime, logger *zap.Logger) {
	logger.Info("poll machine statuses")

	bmcConfigurationList, err := safe.ReaderListAll[*resources.BMCConfiguration](ctx, r)
	if err != nil {
		logger.Error("failed to get all power management info", zap.Error(err))

		return
	}

	bmcConfigurationMap := make(map[resource.ID]*resources.BMCConfiguration, bmcConfigurationList.Len())

	for bmcConfiguration := range bmcConfigurationList.All() {
		bmcConfigurationMap[bmcConfiguration.Metadata().ID()] = bmcConfiguration
	}

	machineStatusList, err := safe.ReaderListAll[*resources.MachineStatus](ctx, r)
	if err != nil {
		logger.Error("failed to get all machine statuses", zap.Error(err))

		return
	}

	connectedMachines := ctrl.agentClient.AllConnectedMachines()

	var numAgentConnected, numAgentDisconnected, numPoweredOn, numPoweredOff, numPowerUnknown int

	for machineStatus := range machineStatusList.All() {
		id := machineStatus.Metadata().ID()

		if machineStatus.Metadata().Phase() == resource.PhaseTearingDown {
			continue
		}

		if machineStatus, err = ctrl.pollSingle(ctx, id, bmcConfigurationMap[id], r, connectedMachines, logger); err != nil {
			logger.Error("failed to poll machine status", zap.String("machine_id", id), zap.Error(err))

			continue
		}

		switch machineStatus.TypedSpec().Value.PowerState {
		case specs.PowerState_POWER_STATE_ON:
			numPoweredOn++
		case specs.PowerState_POWER_STATE_OFF:
			numPoweredOff++
		case specs.PowerState_POWER_STATE_UNKNOWN:
			numPowerUnknown++
		}

		if machineStatus.TypedSpec().Value.AgentAccessible {
			numAgentConnected++
		} else {
			numAgentDisconnected++
		}
	}

	logger.Info("polled machine statuses", zap.Int("agent_connected", numAgentConnected), zap.Int("agent_disconnected", numAgentDisconnected),
		zap.Int("powered_on", numPoweredOn), zap.Int("powered_off", numPoweredOff), zap.Int("power_unknown", numPowerUnknown))
}

func (ctrl *MachineStatusController) pollSingle(ctx context.Context, id resource.ID, bmcConfiguration *resources.BMCConfiguration,
	r controller.ReaderWriter, connectedMachines map[string]struct{}, logger *zap.Logger,
) (*resources.MachineStatus, error) {
	agentAccessible := false
	powerState := specs.PowerState_POWER_STATE_UNKNOWN

	if bmcConfiguration != nil {
		var err error

		powerState, err = ctrl.getPowerState(ctx, bmcConfiguration, logger)
		if err != nil {
			logger.Error("failed to get power state", zap.Error(err))
		}
	}

	if powerState != specs.PowerState_POWER_STATE_OFF { // if the machine is off, no point in pinging the agent
		var err error

		agentAccessible, err = ctrl.agentConnected(ctx, connectedMachines, id)
		if err != nil {
			logger.Error("failed to check agent connection", zap.Error(err))
		}
	}

	return safe.WriterModifyWithResult(ctx, r, resources.NewMachineStatus(id), func(res *resources.MachineStatus) error {
		if res.Metadata().Phase() == resource.PhaseTearingDown {
			return nil
		}

		res.TypedSpec().Value.PowerState = powerState
		res.TypedSpec().Value.AgentAccessible = agentAccessible

		return nil
	})
}

func (ctrl *MachineStatusController) getPowerState(ctx context.Context, bmcConfiguration *resources.BMCConfiguration, logger *zap.Logger) (specs.PowerState, error) {
	bmcClient, err := ctrl.bmcClientFactory.GetClient(ctx, bmcConfiguration, logger)
	if err != nil {
		return specs.PowerState_POWER_STATE_UNKNOWN, err
	}

	defer util.LogClose(bmcClient, logger)

	poweredOn, err := bmcClient.IsPoweredOn(ctx)
	if err != nil {
		logger.Error("failed to get power state", zap.Error(err))

		return specs.PowerState_POWER_STATE_UNKNOWN, err
	}

	if poweredOn {
		return specs.PowerState_POWER_STATE_ON, nil
	}

	return specs.PowerState_POWER_STATE_OFF, nil
}

func (ctrl *MachineStatusController) agentConnected(ctx context.Context, connectedMachines map[string]struct{}, machineID string) (bool, error) {
	if _, connected := connectedMachines[machineID]; !connected { // no tunnel connection, is definitely not connected
		return false, nil
	}

	// attempt to ping
	accessible, err := ctrl.agentClient.IsAccessible(ctx, machineID)
	if err != nil {
		errCode := grpcstatus.Code(err)

		if errCode == codes.Canceled || errCode == codes.DeadlineExceeded || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return false, nil
		}

		return false, err
	}

	return accessible, nil
}

func (ctrl *MachineStatusController) handlePXEBootEvent(ctx context.Context, r controller.QRuntime, pxeBootEvent PXEBootEvent, logger *zap.Logger) error {
	id := pxeBootEvent.MachineID
	tearingDown := false
	initial := false

	if err := safe.WriterModify(ctx, r, resources.NewMachineStatus(pxeBootEvent.MachineID), func(res *resources.MachineStatus) error {
		if res.Metadata().Phase() == resource.PhaseTearingDown {
			tearingDown = true

			return nil
		}

		initial = !res.TypedSpec().Value.Initialized
		res.TypedSpec().Value.Initialized = true

		return nil
	}); err != nil {
		return err
	}

	if tearingDown || !initial {
		return nil
	}

	logger.Debug("discovered new machine, poll status")

	bmcConfiguration, err := safe.ReaderGetByID[*resources.BMCConfiguration](ctx, r, id)
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	if _, err = ctrl.pollSingle(ctx, id, bmcConfiguration, r, ctrl.agentClient.AllConnectedMachines(), logger); err != nil {
		logger.Error("failed to poll machine status", zap.String("machine_id", id), zap.Error(err))
	}

	return nil
}
