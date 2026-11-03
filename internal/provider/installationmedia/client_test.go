// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package installationmedia_test

import (
	"testing"

	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/cosi-project/runtime/pkg/state/impl/inmem"
	"github.com/cosi-project/runtime/pkg/state/impl/namespaced"
	"github.com/siderolabs/image-factory/pkg/schematic"
	omniconstants "github.com/siderolabs/omni/client/pkg/constants"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/installationmedia"
)

func TestAgentModeTalosVersion(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	st := state.WrapCore(namespaced.NewState(inmem.Build))

	for _, version := range []string{"1.13.10", "1.14.0", "1.14.1"} {
		talosVersion := omni.NewTalosVersion(version)
		talosVersion.TypedSpec().Value.Version = version

		require.NoError(t, st.Create(ctx, talosVersion))
	}

	// an Omni that labels no version as its default
	version, err := installationmedia.AgentModeTalosVersion(ctx, st, "")
	require.NoError(t, err)
	assert.Equal(t, omniconstants.DefaultTalosVersion, version)

	defaultVersion := omni.NewTalosVersion("1.14.0")

	_, err = safe.StateUpdateWithConflicts(ctx, st, defaultVersion.Metadata(), func(res *omni.TalosVersion) error {
		res.Metadata().Labels().Set("omni.sidero.dev/default-version", "")

		return nil
	})
	require.NoError(t, err)

	version, err = installationmedia.AgentModeTalosVersion(ctx, st, "")
	require.NoError(t, err)
	assert.Equal(t, "1.14.0", version)

	version, err = installationmedia.AgentModeTalosVersion(ctx, st, "v1.13.10")
	require.NoError(t, err)
	assert.Equal(t, "v1.13.10", version)
}

func TestWithMachineConfig(t *testing.T) {
	t.Parallel()

	const (
		machineConfig   = "apiVersion: v1alpha1\nkind: SideroLinkConfig\napiUrl: https://omni.example.org:8090?jointoken=token\n"
		configKernelArg = "talos.config=http://172.20.0.1:50042/config?u=${uuid}"
	)

	base := schematic.Schematic{Customization: schematic.Customization{ExtraKernelArgs: []string{"console=ttyS0"}}}

	embedded := installationmedia.WithMachineConfig(base, "1.12.0", machineConfig, configKernelArg)
	assert.Equal(t, machineConfig, embedded.Customization.EmbeddedMachineConfiguration)
	assert.Equal(t, []string{"console=ttyS0"}, embedded.Customization.ExtraKernelArgs, "the machine config is not fetched from the provider from Talos 1.12 on")

	fetched := installationmedia.WithMachineConfig(base, "1.11.6", machineConfig, configKernelArg)
	assert.Empty(t, fetched.Customization.EmbeddedMachineConfiguration)
	assert.Equal(t, []string{configKernelArg, "console=ttyS0"}, fetched.Customization.ExtraKernelArgs, "Talos below 1.12 fetches the machine config from the provider")
}
