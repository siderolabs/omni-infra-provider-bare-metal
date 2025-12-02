// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package machineconfig_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/cosi-project/runtime/pkg/state/impl/inmem"
	"github.com/cosi-project/runtime/pkg/state/impl/namespaced"
	"github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/siderolink"
	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/machineconfig"
	providermeta "github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
)

func TestV2ExtraDocs(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	t.Cleanup(cancel)

	st := state.WrapCore(namespaced.NewState(inmem.Build))

	providerJoinConfig := siderolink.NewProviderJoinConfig(providermeta.ProviderID.String())
	providerJoinConfig.TypedSpec().Value.Config = &specs.JoinConfig{
		Config: `apiVersion: v1alpha1
kind: SideroLinkConfig
apiUrl: grpc://127.0.0.1:8090?jointoken=test
---
apiVersion: v1alpha1
kind: EventSinkConfig
endpoint: '[fdae:41e4:649b:9303::1]:8090'
---
apiVersion: v1alpha1
kind: KmsgLogConfig
name: omni-kmsg
url: tcp://[fdae:41e4:649b:9303::1]:8092`,
	}

	require.NoError(t, st.Create(ctx, providerJoinConfig))

	resDef, err := meta.NewResourceDefinition(providerJoinConfig.ResourceDefinition())
	require.NoError(t, err)

	require.NoError(t, st.Create(ctx, resDef))

	configPath := filepath.Join(t.TempDir(), "extra-config.yaml")

	err = os.WriteFile(configPath, []byte(`apiVersion: v1alpha1
kind: KmsgLogConfig
name: extra-doc
url: tcp://[fdae:41e4:649b:9303::1]:8092`), 0o644)
	require.NoError(t, err)

	config, err := machineconfig.Build(ctx, st, nil, configPath)
	require.NoError(t, err)

	returnedConfig, err := configloader.NewFromBytes(config)
	require.NoError(t, err)

	require.Len(t, returnedConfig.Documents(), 4)
	require.Contains(t, string(config), "name: extra-doc")
}
