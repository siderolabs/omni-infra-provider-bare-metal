// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package controllers implements COSI controllers for the bare metal provider.
package controllers

import (
	"context"
	"time"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/generic"
	"github.com/cosi-project/runtime/pkg/controller/generic/qtransform"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/siderolabs/gen/xerrors"
	omnispecs "github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	agentpb "github.com/siderolabs/talos-metal-agent/api/agent"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
)

const (
	// IPMIUsername is the username used for IPMI.
	IPMIUsername = "talos-agent"

	// IPMIPasswordLength is the length of the IPMI password.
	IPMIPasswordLength = 16

	// IPMIDefaultPort is the default port for IPMI.
	IPMIDefaultPort = 623
)

// AgentClient is the interface for interacting with the Talos agent over the reverse GRPC tunnel.
type AgentClient interface {
	GetPowerManagement(ctx context.Context, id string) (*agentpb.GetPowerManagementResponse, error)
	SetPowerManagement(ctx context.Context, id string, req *agentpb.SetPowerManagementRequest) error
	WipeDisks(ctx context.Context, id string) error
	AllConnectedMachines() map[string]struct{}
	IsAccessible(ctx context.Context, machineID string) (bool, error)
}

// BMCClientFactory is the interface for creating BMC clients.
type BMCClientFactory interface {
	GetClient(ctx context.Context, bmcConfiguration *resources.BMCConfiguration) (bmc.Client, error)
}

// handleInput reads the additional input resource and automatically manages finalizers.
func handleInput[T generic.ResourceWithRD, S generic.ResourceWithRD](ctx context.Context, r controller.ReaderWriter, finalizer string, main S) (T, error) {
	var zero T

	res, err := safe.ReaderGetByID[T](ctx, r, main.Metadata().ID())
	if err != nil {
		if state.IsNotFoundError(err) {
			return zero, nil
		}

		return zero, err
	}

	if res.Metadata().Phase() == resource.PhaseTearingDown || main.Metadata().Phase() == resource.PhaseTearingDown {
		if err = r.RemoveFinalizer(ctx, res.Metadata(), finalizer); err != nil && !state.IsNotFoundError(err) {
			return zero, err
		}

		if res.Metadata().Phase() == resource.PhaseTearingDown {
			return zero, nil
		}

		return res, nil
	}

	if !res.Metadata().Finalizers().Has(finalizer) {
		if err = r.AddFinalizer(ctx, res.Metadata(), finalizer); err != nil {
			return zero, err
		}
	}

	return res, nil
}

func validateInfraMachine(infraMachine *infra.Machine, logger *zap.Logger) error {
	if infraMachine.TypedSpec().Value.AcceptanceStatus != omnispecs.InfraMachineConfigSpec_ACCEPTED {
		logger.Debug("machine not accepted, skip")

		return xerrors.NewTaggedf[qtransform.SkipReconcileTag]("machine not accepted")
	}

	if infraMachine.TypedSpec().Value.Cordoned {
		logger.Debug("machine cordoned, skip")

		return xerrors.NewTaggedf[qtransform.SkipReconcileTag]("machine is cordoned")
	}

	return nil
}

func getTimeSinceLastPowerOn(powerOperation *resources.PowerOperation, rebootStatus *resources.RebootStatus) time.Duration {
	var lastPowerOnTime time.Time

	if powerOperation != nil && powerOperation.TypedSpec().Value.LastPowerOnTimestamp != nil {
		lastPowerOnTime = powerOperation.TypedSpec().Value.LastPowerOnTimestamp.AsTime()
	}

	if rebootStatus != nil && rebootStatus.TypedSpec().Value.LastRebootTimestamp != nil {
		lastRebootTime := rebootStatus.TypedSpec().Value.LastRebootTimestamp.AsTime()
		if lastRebootTime.After(lastPowerOnTime) {
			lastPowerOnTime = lastRebootTime
		}
	}

	return time.Since(lastPowerOnTime)
}
