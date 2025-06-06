// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cosi-project/runtime/pkg/controller/runtime"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/cosi-project/runtime/pkg/state/impl/inmem"
	"github.com/cosi-project/runtime/pkg/state/impl/namespaced"
	"github.com/siderolabs/gen/containers"
	"github.com/siderolabs/gen/pair"
	agentpb "github.com/siderolabs/talos-metal-agent/api/agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"golang.org/x/sync/errgroup"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/pxe"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
)

func init() {
	if err := provider.RegisterResources(); err != nil {
		panic("failed to register resources: " + err.Error())
	}
}

type testFunc func(ctx context.Context, st state.State, rt *runtime.Runtime, logger *zap.Logger)

func withRuntime(t *testing.T, beforeStart, afterStart testFunc) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)

	logger := zaptest.NewLogger(t)
	st := state.WrapCore(namespaced.NewState(inmem.Build))

	cosiRuntime, err := provider.BuildCOSIRuntime(st, false, logger)
	require.NoError(t, err)

	beforeStart(ctx, st, cosiRuntime, logger)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return cosiRuntime.Run(ctx)
	})

	afterStart(ctx, st, cosiRuntime, logger)

	cancel()

	require.NoError(t, eg.Wait())
}

type bmcClientFactoryMock struct {
	bmcClient bmc.Client
}

func (b *bmcClientFactoryMock) GetClient(context.Context, *resources.BMCConfiguration, *zap.Logger) (bmc.Client, error) {
	return b.bmcClient, nil
}

type bmcClientMock struct {
	powerOnCh        chan<- struct{}
	rebootCh         chan<- struct{}
	setPXEBootOnceCh chan<- pxe.BootMode
	poweredOn        bool
}

func (b *bmcClientMock) Close() error {
	return nil
}

func (b *bmcClientMock) Reboot(ctx context.Context) error {
	select {
	case b.rebootCh <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (b *bmcClientMock) IsPoweredOn(context.Context) (bool, error) {
	return b.poweredOn, nil
}

func (b *bmcClientMock) PowerOn(ctx context.Context) error {
	select {
	case b.powerOnCh <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (b *bmcClientMock) PowerOff(context.Context) error {
	return nil
}

func (b *bmcClientMock) SetPXEBootOnce(ctx context.Context, mode pxe.BootMode) error {
	select {
	case b.setPXEBootOnceCh <- mode:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

type agentClientMock struct {
	getPowerMgmtResponseMap *containers.ConcurrentMap[string, *agentpb.GetPowerManagementResponse]
	setPowerMgmtRequestCh   chan<- pair.Pair[string, *agentpb.SetPowerManagementRequest]
}

func (a *agentClientMock) GetPowerManagement(_ context.Context, id string) (*agentpb.GetPowerManagementResponse, error) {
	val, _ := a.getPowerMgmtResponseMap.Get(id)

	return val, nil
}

func (a *agentClientMock) SetPowerManagement(ctx context.Context, id string, req *agentpb.SetPowerManagementRequest) error {
	select {
	case a.setPowerMgmtRequestCh <- pair.MakePair(id, req):
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (a *agentClientMock) WipeDisks(context.Context, string) error {
	return nil
}

func (a *agentClientMock) AllConnectedMachines() map[string]struct{} {
	return nil
}

func (a *agentClientMock) IsAccessible(context.Context, string) (bool, error) {
	return false, nil
}

func requireChReceive[T any](ctx context.Context, t *testing.T, ch chan T) T {
	select {
	case val := <-ch:
		return val
	case <-ctx.Done():
		require.Fail(t, "timeout waiting for channel receive")

		return *new(T) // unreachable
	}
}

type reconcileAsserter struct {
	t             *testing.T
	numReconciles atomic.Uint32
}

func (r *reconcileAsserter) incrementReconcile() {
	r.numReconciles.Add(1)
}

func (r *reconcileAsserter) requireNotZero(ctx context.Context) uint32 {
	var numReconciles uint32

	require.EventuallyWithT(r.t, func(c *assert.CollectT) {
		require.NoError(c, ctx.Err())

		numReconciles = r.numReconciles.Load()

		assert.NotZero(c, numReconciles, "expected at least one reconcile to be called")
	}, 5*time.Second, 100*time.Millisecond, "expected at least one reconcile to be called")

	r.t.Logf("numReconciles: %d", numReconciles)

	return numReconciles
}

func (r *reconcileAsserter) requireReconcile(ctx context.Context, before uint32) {
	require.EventuallyWithT(r.t, func(c *assert.CollectT) {
		require.NoError(c, ctx.Err())

		numReconciles := r.numReconciles.Load()
		if assert.Greater(c, numReconciles, before, "expected at least one reconcile to be called") {
			r.t.Logf("numReconciles = %d", numReconciles)
		}
	}, 5*time.Second, 100*time.Millisecond, "expected at least one reconcile to be called")
}
