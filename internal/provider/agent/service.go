// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package agent

import (
	"context"
	"time"

	"github.com/cosi-project/runtime/pkg/state"
	"github.com/jhump/grpctunnel"
	"github.com/jhump/grpctunnel/tunnelpb"
	agentpb "github.com/siderolabs/talos-metal-agent/api/agent"
	agentconstants "github.com/siderolabs/talos-metal-agent/pkg/constants"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/baremetal"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/machinestatus"
)

const (
	timeout           = 30 * time.Second
	zeroesWipeTimeout = 24 * time.Hour
)

// Service controls servers by establishing a reverse GRPC tunnel with them and by sending them commands.
type Service struct {
	logger        *zap.Logger
	grpcServer    grpc.ServiceRegistrar
	tunnelHandler *grpctunnel.TunnelServiceHandler

	wipeWithZeroes bool
}

// NewService creates a new agent service.
func NewService(grpcServer grpc.ServiceRegistrar, state state.State, wipeWithZeroes bool, logger *zap.Logger) *Service {
	tunnelHandler := grpctunnel.NewTunnelServiceHandler(
		grpctunnel.TunnelServiceHandlerOptions{
			OnReverseTunnelOpen: func(channel grpctunnel.TunnelChannel) {
				handleTunnelEvent(channel, state, true, logger)
			},
			OnReverseTunnelClose: func(channel grpctunnel.TunnelChannel) {
				handleTunnelEvent(channel, state, false, logger)
			},
			AffinityKey: func(channel grpctunnel.TunnelChannel) any {
				id, ok := machineIDAffinityKey(channel.Context(), logger)
				if !ok {
					return "invalid"
				}

				return id
			},
		},
	)

	tunnelpb.RegisterTunnelServiceServer(grpcServer, tunnelHandler.Service())

	return &Service{
		logger:         logger,
		grpcServer:     grpcServer,
		tunnelHandler:  tunnelHandler,
		wipeWithZeroes: wipeWithZeroes,
	}
}

func handleTunnelEvent(channel grpctunnel.TunnelChannel, state state.State, connected bool, logger *zap.Logger) {
	affinityKey, ok := machineIDAffinityKey(channel.Context(), logger)
	if !ok {
		logger.Warn("invalid affinity key", zap.String("reason", "no machine ID in metadata"))

		return
	}

	logger = logger.With(zap.String("machine_id", affinityKey), zap.Bool("connected", connected))

	logger.Debug("machine tunnel event")

	if channel.Context().Err() != nil { // context is closed, probably the app is shutting down, nothing to do
		return
	}

	if connected { // if an agent is connected, update the boot mode to PXE
		if _, err := machinestatus.Modify(channel.Context(), state, affinityKey, func(status *baremetal.MachineStatus) error {
			status.TypedSpec().Value.BootMode = specs.BootMode_BOOT_MODE_AGENT_PXE

			return nil
		}); err != nil {
			logger.Error("failed to update machine status", zap.String("machine_id", affinityKey), zap.Error(err))
		}
	}
}

// IsAccessible checks if the agent with the given ID is accessible.
func (c *Service) IsAccessible(ctx context.Context, id string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	channel := c.tunnelHandler.KeyAsChannel(id)
	cli := agentpb.NewAgentServiceClient(channel)

	_, err := cli.Hello(ctx, &agentpb.HelloRequest{})
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetPowerManagement retrieves the IPMI information from the server with the given ID.
func (c *Service) GetPowerManagement(ctx context.Context, id string) (*agentpb.GetPowerManagementResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	channel := c.tunnelHandler.KeyAsChannel(id)
	cli := agentpb.NewAgentServiceClient(channel)

	return cli.GetPowerManagement(ctx, &agentpb.GetPowerManagementRequest{})
}

// SetPowerManagement sets the IPMI information on the server with the given ID.
func (c *Service) SetPowerManagement(ctx context.Context, id string, req *agentpb.SetPowerManagementRequest) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	channel := c.tunnelHandler.KeyAsChannel(id)
	cli := agentpb.NewAgentServiceClient(channel)

	_, err := cli.SetPowerManagement(ctx, req)

	return err
}

// WipeDisks wipes the disks on the server with the given ID.
func (c *Service) WipeDisks(ctx context.Context, id string) error {
	channel := c.tunnelHandler.KeyAsChannel(id)
	cli := agentpb.NewAgentServiceClient(channel)

	wipeTimeout := timeout
	if c.wipeWithZeroes {
		wipeTimeout = zeroesWipeTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, wipeTimeout)
	defer cancel()

	_, err := cli.WipeDisks(ctx, &agentpb.WipeDisksRequest{
		Zeroes: c.wipeWithZeroes,
	})

	return err
}

// AllConnectedMachines returns a set of all connected machines.
func (c *Service) AllConnectedMachines() map[string]struct{} {
	allTunnels := c.tunnelHandler.AllReverseTunnels()

	machines := make(map[string]struct{}, len(allTunnels))

	for _, tunnel := range allTunnels {
		affinityKey, ok := machineIDAffinityKey(tunnel.Context(), c.logger)
		if !ok {
			c.logger.Warn("invalid affinity key", zap.String("reason", "no machine ID in metadata"))

			continue
		}

		machines[affinityKey] = struct{}{}
	}

	return machines
}

func machineIDAffinityKey(ctx context.Context, logger *zap.Logger) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Warn("invalid affinity key", zap.String("reason", "no metadata"))

		return "", false
	}

	machineID := md.Get(agentconstants.MachineIDMetadataKey)
	if len(machineID) == 0 {
		logger.Warn("invalid affinity key", zap.String("reason", "no machine ID in metadata"))

		return "", false
	}

	if len(machineID) > 1 {
		logger.Warn("multiple machine IDs in metadata", zap.Strings("machine_ids", machineID))
	}

	return machineID[0], true
}