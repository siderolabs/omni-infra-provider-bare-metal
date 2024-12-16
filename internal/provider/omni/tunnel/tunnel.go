// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package tunnel provides the reverse GRPC tunnel to Omni.
package tunnel

import (
	"context"
	"errors"
	"time"

	"github.com/cosi-project/runtime/pkg/state"
	"github.com/jhump/grpctunnel"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	providerpb "github.com/siderolabs/omni-infra-provider-bare-metal/api/provider"
	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/service"
)

// Tunneler is the interface for the reverse GRPC tunnel.
type Tunneler interface {
	Tunnel() *grpctunnel.ReverseTunnelServer
}

// PowerClientFactory is the interface for creating power clients.
type PowerClientFactory interface {
	GetClient(powerManagement *specs.PowerManagement) (power.Client, error)
}

// Tunnel represents the reverse GRPC tunnel to Omni.
type Tunnel struct {
	tunneler           Tunneler
	powerClientFactory PowerClientFactory
	logger             *zap.Logger
	state              state.State
}

// New creates a new Tunnel.
func New(powerClientFactory PowerClientFactory, state state.State, tunneler Tunneler, logger *zap.Logger) *Tunnel {
	return &Tunnel{
		powerClientFactory: powerClientFactory,
		state:              state,
		tunneler:           tunneler,
		logger:             logger,
	}
}

// Run starts the reverse GRPC tunnel to Omni.
func (tunnel *Tunnel) Run(ctx context.Context) error {
	reverseTunnelServer := tunnel.tunneler.Tunnel()
	providerServiceServer := service.NewProviderServiceServer(tunnel.powerClientFactory, tunnel.state, tunnel.logger)

	providerpb.RegisterProviderServiceServer(reverseTunnelServer, providerServiceServer)

	// Open the reverse tunnel and serve requests.
	for {
		if _, err := reverseTunnelServer.Serve(ctx); err != nil {
			if status.Code(err) == codes.Canceled || errors.Is(err, context.Canceled) {
				return nil
			}

			tunnel.logger.Error("failed to serve reverse tunnel", zap.Error(err))

			select {
			case <-ctx.Done():
				return nil
			case <-time.After(10 * time.Second):
			}
		}
	}
}
