// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package machineconfig builds the machine configuration for the bare-metal provider.
package machineconfig

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/siderolabs/omni/client/pkg/jointoken"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	siderolinkres "github.com/siderolabs/omni/client/pkg/omni/resources/siderolink"
	"github.com/siderolabs/omni/client/pkg/siderolink"
	"github.com/siderolabs/talos/pkg/machinery/config/config"
	"github.com/siderolabs/talos/pkg/machinery/config/types/security"

	providermeta "github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/tls"
)

// Build builds the machine configuration for the bare-metal provider.
func Build(ctx context.Context, r controller.Reader, certs *tls.Certs) ([]byte, error) {
	var extraDocs []config.Document

	if certs != nil && certs.CACertPEM != "" {
		trustedRootsConfig := security.NewTrustedRootsConfigV1Alpha1()
		trustedRootsConfig.MetaName = "infra-provider-ca"
		trustedRootsConfig.Certificates = certs.CACertPEM

		extraDocs = append(extraDocs, trustedRootsConfig)
	}

	providerJoinConfigRD, err := safe.ReaderGetByID[*meta.ResourceDefinition](ctx, r, strings.ToLower(siderolinkres.ProviderJoinConfigType))
	if err != nil && !state.IsNotFoundError(err) {
		return nil, err
	}

	// V2 flow
	if providerJoinConfigRD != nil {
		var providerJoinConfig *siderolinkres.ProviderJoinConfig

		providerJoinConfig, err = safe.ReaderGetByID[*siderolinkres.ProviderJoinConfig](ctx, r, providermeta.ProviderID.String())
		if err != nil {
			return nil, err
		}

		return []byte(providerJoinConfig.TypedSpec().Value.Config.Config), nil
	}

	// keeping this code for compatibility with the older Omni versions
	connectionParams, err := safe.ReaderGetByID[*siderolinkres.ConnectionParams](ctx, r, siderolinkres.ConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection params: %w", err)
	}

	opts, err := siderolink.NewJoinOptions(
		siderolink.WithJoinToken(connectionParams.TypedSpec().Value.JoinToken),
		siderolink.WithMachineAPIURL(connectionParams.TypedSpec().Value.ApiEndpoint),
		siderolink.WithEventSinkPort(int(connectionParams.TypedSpec().Value.EventsPort)),
		siderolink.WithLogServerPort(int(connectionParams.TypedSpec().Value.LogsPort)),
		siderolink.WithProvider(infra.NewProvider(providermeta.ProviderID.String())),
		siderolink.WithJoinTokenVersion(jointoken.Version1),
	)
	if err != nil {
		return nil, err
	}

	return opts.RenderJoinConfig(extraDocs...)
}
