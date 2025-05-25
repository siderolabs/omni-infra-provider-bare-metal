// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers_test

import (
	"context"
	"testing"
	"time"

	"github.com/cosi-project/runtime/pkg/controller/runtime"
	"github.com/cosi-project/runtime/pkg/resource/rtestutils"
	"github.com/cosi-project/runtime/pkg/state"
	omnispecs "github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/pxe"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/controllers"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
)

func TestPowerOn(t *testing.T) {
	t.Parallel()

	powerOnCh := make(chan struct{})
	setPXEBootOnceCh := make(chan pxe.BootMode)

	bmcClientFactory := &bmcClientFactoryMock{
		bmcClient: &bmcClientMock{
			poweredOn:        false,
			powerOnCh:        powerOnCh,
			setPXEBootOnceCh: setPXEBootOnceCh,
		},
	}

	pxeBootMode := pxe.BootModeUEFI

	now := time.Now()
	nowFunc := func() time.Time { return now }

	withRuntime(t,
		func(_ context.Context, _ state.State, rt *runtime.Runtime, _ *zap.Logger) {
			controller := controllers.NewPowerOperationController(nowFunc, bmcClientFactory, 0, pxeBootMode)

			require.NoError(t, rt.RegisterQController(controller))
		},
		func(ctx context.Context, st state.State, _ *runtime.Runtime, _ *zap.Logger) {
			bmcConfiguration := resources.NewBMCConfiguration("test-machine")

			require.NoError(t, st.Create(ctx, bmcConfiguration))

			infraMachine := infra.NewMachine("test-machine")

			infraMachine.TypedSpec().Value.AcceptanceStatus = omnispecs.InfraMachineConfigSpec_ACCEPTED

			require.NoError(t, st.Create(ctx, infraMachine))

			// expect a SetPXEBootOnce call
			mode := requireChReceive(ctx, t, setPXEBootOnceCh)
			require.Equal(t, pxeBootMode, mode)

			// expect a PowerOn call
			requireChReceive(ctx, t, powerOnCh)

			rtestutils.AssertResource(ctx, t, st, infraMachine.Metadata().ID(), func(res *resources.PowerOperation, assertion *assert.Assertions) {
				assertion.Equal(specs.PowerState_POWER_STATE_ON, res.TypedSpec().Value.LastPowerOperation)
				assertion.Equal(now.Unix(), res.TypedSpec().Value.LastPowerOnTimestamp.AsTime().Unix())
			})
		},
	)
}
