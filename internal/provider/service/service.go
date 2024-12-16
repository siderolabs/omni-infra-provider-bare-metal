// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package service implements the bare metal infra provider GRPC service server.
package service

import (
	"context"

	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/provider"
	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/baremetal"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power"
)

// PowerClientFactory is the interface for creating power clients.
type PowerClientFactory interface {
	GetClient(powerManagement *specs.PowerManagement) (power.Client, error)
}

// ProviderServiceServer is the bare metal infra provider service server.
type ProviderServiceServer struct {
	providerpb.UnimplementedProviderServiceServer

	logger             *zap.Logger
	state              state.State
	powerClientFactory PowerClientFactory
}

// NewProviderServiceServer creates a new ProviderServiceServer.
func NewProviderServiceServer(powerClientFactory PowerClientFactory, state state.State, logger *zap.Logger) *ProviderServiceServer {
	return &ProviderServiceServer{
		powerClientFactory: powerClientFactory,
		state:              state,
		logger:             logger,
	}
}

// RebootMachine reboots a machine.
func (service *ProviderServiceServer) RebootMachine(ctx context.Context, request *providerpb.RebootMachineRequest) (*providerpb.RebootMachineResponse, error) {
	status, err := safe.StateGetByID[*baremetal.MachineStatus](ctx, service.state, request.Id)
	if err != nil {
		return nil, err
	}

	powerClient, err := service.powerClientFactory.GetClient(status.TypedSpec().Value.PowerManagement)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := powerClient.Close(); closeErr != nil {
			service.logger.Error("failed to close power client", zap.Error(closeErr))
		}
	}()

	if err = powerClient.Reboot(ctx); err != nil {
		return nil, err
	}

	return &providerpb.RebootMachineResponse{}, nil
}
