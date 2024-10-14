// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package machinestatus

import (
	"context"
	"errors"
	"time"

	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/baremetal"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power"
)

// AgentController is the interface for controlling Talos agent.
type AgentController interface {
	AllConnectedMachines() map[string]struct{}
	IsAccessible(ctx context.Context, machineID string) (bool, error)
}

// Poller polls the machines periodically and updates their statuses.
type Poller struct {
	agentController AgentController
	state           state.State
	logger          *zap.Logger
}

// NewPoller creates a new Poller.
func NewPoller(agentController AgentController, state state.State, logger *zap.Logger) *Poller {
	return &Poller{
		agentController: agentController,
		state:           state,
		logger:          logger,
	}
}

// Run starts the Poller.
func (m *Poller) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil
			}

			return ctx.Err()
		case <-ticker.C:
		}

		m.poll(ctx)
	}
}

func (m *Poller) poll(ctx context.Context) {
	statusList, err := safe.StateListAll[*baremetal.MachineStatus](ctx, m.state)
	if err != nil {
		m.logger.Error("failed to get all machine statuses", zap.Error(err))

		return
	}

	connectedMachines := m.agentController.AllConnectedMachines()
	machineIDSet := map[string]struct{}{}

	var numAgentConnected, numAgentDisconnected, numPoweredOn, numPoweredOff, numPowerUnknown int

	for status := range statusList.All() {
		machineIDSet[status.Metadata().ID()] = struct{}{}

		powerState := m.getPowerState(ctx, status, m.logger)

		switch powerState {
		case specs.PowerState_POWER_STATE_ON:
			numPoweredOn++
		case specs.PowerState_POWER_STATE_OFF:
			numPoweredOff++
		case specs.PowerState_POWER_STATE_UNKNOWN:
			numPowerUnknown++
		}

		agentConnected := false

		agentConnected, err = m.agentConnected(ctx, connectedMachines, status.Metadata().ID())
		if err != nil {
			connectionError := errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || grpcstatus.Code(err) == codes.Canceled

			if !connectionError {
				m.logger.Error("failed to check connection", zap.String("machine_id", status.Metadata().ID()), zap.Error(err))

				continue
			}
		}

		if agentConnected {
			numAgentConnected++
		} else {
			numAgentDisconnected++
		}

		if _, saveErr := Modify(ctx, m.state, status.Metadata().ID(), func(status *baremetal.MachineStatus) error {
			status.TypedSpec().Value.PowerState = powerState

			if agentConnected { // if the agent is connected, we know for sure that we are in the agent mode, so we update the status accordingly
				status.TypedSpec().Value.BootMode = specs.BootMode_BOOT_MODE_AGENT_PXE
			}

			return nil
		}); saveErr != nil {
			m.logger.Error("failed to save status", zap.String("machine_id", status.Metadata().ID()), zap.Error(saveErr))
		}
	}

	m.logger.Info("polled machine statuses", zap.Int("agent_connected", numAgentConnected), zap.Int("agent_disconnected", numAgentDisconnected),
		zap.Int("powered_on", numPoweredOn), zap.Int("powered_off", numPoweredOff), zap.Int("power_unknown", numPowerUnknown))
}

func (m *Poller) getPowerState(ctx context.Context, status *baremetal.MachineStatus, logger *zap.Logger) specs.PowerState {
	powerClient, err := power.GetClient(status.TypedSpec().Value.PowerManagement)
	if err != nil {
		if errors.Is(err, power.ErrNoPowerManagementInfo) {
			logger.Debug("no power management info yet, skip update", zap.String("machine_id", status.Metadata().ID()))
		} else {
			logger.Error("failed to get power client", zap.String("machine_id", status.Metadata().ID()), zap.Error(err))
		}

		return specs.PowerState_POWER_STATE_UNKNOWN
	}

	poweredOn, err := powerClient.IsPoweredOn(ctx)
	if err != nil {
		logger.Error("failed to get power state", zap.String("machine_id", status.Metadata().ID()), zap.Error(err))

		return specs.PowerState_POWER_STATE_UNKNOWN
	}

	if poweredOn {
		return specs.PowerState_POWER_STATE_ON
	}

	return specs.PowerState_POWER_STATE_OFF
}

func (m *Poller) agentConnected(ctx context.Context, connectedMachines map[string]struct{}, machineID string) (bool, error) {
	if _, connected := connectedMachines[machineID]; !connected { // no tunnel connection, is definitely not connected
		return false, nil
	}

	// attempt to ping
	accessible, err := m.agentController.IsAccessible(ctx, machineID)
	if err != nil {
		return false, err
	}

	return accessible, nil
}