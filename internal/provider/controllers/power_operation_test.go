// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers_test

import (
	"context"
	"testing"
	"time"

	"github.com/cosi-project/runtime/pkg/resource/rtestutils"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/cosi-project/runtime/pkg/state/impl/inmem"
	"github.com/cosi-project/runtime/pkg/state/impl/namespaced"
	omnispecs "github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"golang.org/x/sync/errgroup"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/pxe"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/controllers"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
)

func TestPowerOn(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	powerOnCh := make(chan struct{})
	setPXEBootOnceCh := make(chan pxe.BootMode)

	bmcClientFactory := &bmcClientFactoryMock{
		bmcClient: &bmcClientMock{
			poweredOn:        false,
			powerOnCh:        powerOnCh,
			setPXEBootOnceCh: setPXEBootOnceCh,
		},
	}

	logger := zaptest.NewLogger(t)
	st := state.WrapCore(namespaced.NewState(inmem.Build))

	cosiRuntime, err := provider.BuildCOSIRuntime(st, false, logger)
	require.NoError(t, err)

	pxeBootMode := pxe.BootModeUEFI

	now := time.Now()
	nowFunc := func() time.Time { return now }

	controller := controllers.NewPowerOperationController(nowFunc, bmcClientFactory, 0, pxeBootMode)

	require.NoError(t, cosiRuntime.RegisterQController(controller))

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return cosiRuntime.Run(ctx)
	})

	bmcConfiguration := resources.NewBMCConfiguration("test-machine")

	require.NoError(t, st.Create(ctx, bmcConfiguration))

	infraMachine := infra.NewMachine("test-machine")

	infraMachine.TypedSpec().Value.AcceptanceStatus = omnispecs.InfraMachineConfigSpec_ACCEPTED

	require.NoError(t, st.Create(ctx, infraMachine))

	select { // expect a SetPXEBootOnce call
	case mode := <-setPXEBootOnceCh:
		require.Equal(t, pxeBootMode, mode)
	case <-ctx.Done():
		require.Fail(t, "timeout waiting for SetPXEBootOnce")
	}

	select { // expect a PowerOn call
	case <-powerOnCh:
	case <-ctx.Done():
		require.Fail(t, "timeout waiting for PowerOn")
	}

	rtestutils.AssertResource(ctx, t, st, infraMachine.Metadata().ID(), func(res *resources.PowerOperation, assertion *assert.Assertions) {
		assertion.Equal(specs.PowerState_POWER_STATE_ON, res.TypedSpec().Value.LastPowerOperation)
		assertion.Equal(now.Unix(), res.TypedSpec().Value.LastPowerOnTimestamp.AsTime().Unix())
	})

	cancel()

	require.NoError(t, eg.Wait())
}

type bmcClientFactoryMock struct {
	bmcClient bmc.Client
}

func (b *bmcClientFactoryMock) GetClient(context.Context, *resources.BMCConfiguration) (bmc.Client, error) {
	return b.bmcClient, nil
}

type bmcClientMock struct {
	powerOnCh        chan<- struct{}
	setPXEBootOnceCh chan<- pxe.BootMode
	poweredOn        bool
}

func (b *bmcClientMock) Close() error {
	return nil
}

func (b *bmcClientMock) Reboot(context.Context) error {
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
