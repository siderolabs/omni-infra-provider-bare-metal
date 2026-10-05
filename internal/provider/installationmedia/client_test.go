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
