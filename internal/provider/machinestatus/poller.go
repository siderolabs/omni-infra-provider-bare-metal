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

// AgentService is the interface for controlling Talos agent.
type AgentService interface {
	AllConnectedMachines() map[string]struct{}
	IsAccessible(ctx context.Context, machineID string) (bool, error)
}

// PowerClientFactory is the interface for creating power clients.
type PowerClientFactory interface {
	GetClient(ctx context.Context, powerManagement *specs.PowerManagement) (power.Client, error)
}

// Poller polls the machines periodically and updates their statuses.
type Poller struct {
	agentService       AgentService
	powerClientFactory PowerClientFactory
	state              state.State
	logger             *zap.Logger
}

// NewPoller creates a new Poller.
func NewPoller(agentService AgentService, powerClientFactory PowerClientFactory, state state.State, logger *zap.Logger) *Poller {
	return &Poller{
		agentService:       agentService,
		powerClientFactory: powerClientFactory,
		state:              state,
		logger:             logger,
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

	connectedMachines := m.agentService.AllConnectedMachines()
	machineIDSet := map[string]struct{}{}

	var numAgentConnected, numAgentDisconnected, numPoweredOn, numPoweredOff, numPowerUnknown int

	for status := range statusList.All() {
		machineIDSet[status.Metadata().ID()] = struct{}{}

		logger := m.logger.With(zap.String("machine_id", status.Metadata().ID()))

		powerState := m.getPowerState(ctx, status, logger)

		switch powerState {
		case specs.PowerState_POWER_STATE_ON:
			numPoweredOn++
		case specs.PowerState_POWER_STATE_OFF:
			numPoweredOff++
		case specs.PowerState_POWER_STATE_UNKNOWN:
			numPowerUnknown++
		}

		logger = logger.With(zap.Stringer("power_state", powerState))

		agentConnected := false

		agentConnected, err = m.agentConnected(ctx, connectedMachines, status.Metadata().ID())
		if err != nil {
			m.logger.Error("failed to check connection", zap.String("machine_id", status.Metadata().ID()), zap.Error(err))
		}

		if agentConnected {
			numAgentConnected++
		} else {
			numAgentDisconnected++
		}

		if _, err = Modify(ctx, m.state, status.Metadata().ID(), func(status *baremetal.MachineStatus) error {
			status.TypedSpec().Value.PowerState = powerState

			if agentConnected { // if the agent is connected, we know for sure that we are in the agent mode, so we update the status accordingly
				status.TypedSpec().Value.BootMode = specs.BootMode_BOOT_MODE_AGENT_PXE
			}

			return nil
		}); err != nil {
			logger.Error("failed to save status", zap.Error(err))
		}
	}

	m.logger.Info("polled machine statuses", zap.Int("agent_connected", numAgentConnected), zap.Int("agent_disconnected", numAgentDisconnected),
		zap.Int("powered_on", numPoweredOn), zap.Int("powered_off", numPoweredOff), zap.Int("power_unknown", numPowerUnknown))
}

func (m *Poller) getPowerState(ctx context.Context, status *baremetal.MachineStatus, logger *zap.Logger) specs.PowerState {
	powerClient, err := m.powerClientFactory.GetClient(ctx, status.TypedSpec().Value.PowerManagement)
	if err != nil {
		if errors.Is(err, power.ErrNoPowerManagementInfo) {
			logger.Debug("no power management info yet, skip update")
		} else {
			logger.Error("failed to get power client", zap.Error(err))
		}

		return specs.PowerState_POWER_STATE_UNKNOWN
	}

	defer func() {
		if closeErr := powerClient.Close(); closeErr != nil {
			logger.Error("failed to close power client", zap.Error(closeErr))
		}
	}()

	poweredOn, err := powerClient.IsPoweredOn(ctx)
	if err != nil {
		logger.Error("failed to get power state", zap.Error(err))

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
	accessible, err := m.agentService.IsAccessible(ctx, machineID)
	if err != nil {
		errCode := grpcstatus.Code(err)

		if errCode == codes.Canceled || errCode == codes.DeadlineExceeded || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return false, nil
		}

		return false, err
	}

	return accessible, nil
}
