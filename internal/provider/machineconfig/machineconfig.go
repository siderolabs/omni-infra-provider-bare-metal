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
	"github.com/hashicorp/go-multierror"
	"github.com/siderolabs/omni/client/pkg/jointoken"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	siderolinkres "github.com/siderolabs/omni/client/pkg/omni/resources/siderolink"
	"github.com/siderolabs/omni/client/pkg/siderolink"
	"github.com/siderolabs/talos/pkg/machinery/config/config"
	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
	"github.com/siderolabs/talos/pkg/machinery/config/container"
	"github.com/siderolabs/talos/pkg/machinery/config/encoder"
	"github.com/siderolabs/talos/pkg/machinery/config/types/runtime"
	"github.com/siderolabs/talos/pkg/machinery/config/types/security"
	siderolinktalos "github.com/siderolabs/talos/pkg/machinery/config/types/siderolink"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"

	providermeta "github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/tls"
)

const (
	kmsgLogName         = "siderolink-agent-kmsg-log"
	infraProviderCAName = "infra-provider-ca"
)

// Build builds the machine configuration for the bare-metal provider.
func Build(ctx context.Context, r controller.Reader, certs *tls.Certs, extraMachineConfigPath string) ([]byte, error) {
	var extraDocs []config.Document

	if certs != nil && certs.CACertPEM != "" {
		trustedRootsConfig := security.NewTrustedRootsConfigV1Alpha1()
		trustedRootsConfig.MetaName = infraProviderCAName
		trustedRootsConfig.Certificates = certs.CACertPEM

		extraDocs = append(extraDocs, trustedRootsConfig)
	}

	if extraMachineConfigPath != "" {
		additionalDocs, err := parseAdditionalDocuments(extraMachineConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load extra machine config: %w", err)
		}

		extraDocs = append(extraDocs, additionalDocs...)
	}

	providerJoinConfigRD, err := safe.ReaderGetByID[*meta.ResourceDefinition](ctx, r, strings.ToLower(siderolinkres.ProviderJoinConfigType))
	if err != nil && !state.IsNotFoundError(err) {
		return nil, err
	}

	// V2 flow
	if providerJoinConfigRD != nil {
		return buildV2MachineConfig(ctx, r, extraDocs)
	}

	// keeping this code for compatibility with the older Omni versions
	connectionParams, err := safe.ReaderGetByID[*siderolinkres.ConnectionParams](ctx, r, siderolinkres.ConfigID) //nolint:staticcheck
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

func buildV2MachineConfig(ctx context.Context, r controller.Reader, extraDocs []config.Document) ([]byte, error) {
	providerJoinConfig, err := safe.ReaderGetByID[*siderolinkres.ProviderJoinConfig](ctx, r, providermeta.ProviderID.String())
	if err != nil {
		return nil, err
	}

	parsedConfig, err := configloader.NewFromReader(strings.NewReader(providerJoinConfig.TypedSpec().Value.Config.Config))
	if err != nil {
		return nil, fmt.Errorf("failed to load provider join config: %w", err)
	}

	docs := parsedConfig.Documents()
	docs = append(docs, extraDocs...)

	configContainer, err := container.New(docs...)
	if err != nil {
		return nil, fmt.Errorf("failed to create config container: %w", err)
	}

	return configContainer.EncodeBytes(encoder.WithComments(encoder.CommentsDisabled))
}

func parseAdditionalDocuments(extraMachineConfigPath string) ([]config.Document, error) {
	loadedConfig, err := configloader.NewFromFile(extraMachineConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load extra machine config: %w", err)
	}

	var errs error

	for _, document := range loadedConfig.Documents() {
		switch document.Kind() {
		case v1alpha1.Version, siderolinktalos.Kind, runtime.EventSinkKind:
			errs = multierror.Append(errs, fmt.Errorf("extra machine config must not contain %s documents", document.Kind()))
		case runtime.KmsgLogKind:
			kmsgLogConfig, ok := document.(*runtime.KmsgLogV1Alpha1)
			if !ok {
				errs = multierror.Append(errs, fmt.Errorf("expected %s document, got %T", runtime.KmsgLogKind, document))

				continue
			}

			if kmsgLogConfig.MetaName == kmsgLogName {
				errs = multierror.Append(errs, fmt.Errorf("extra machine config must not contain %s document with name %q", runtime.KmsgLogKind, kmsgLogName))
			}
		case security.TrustedRootsConfig:
			trustedRootsConfig, ok := document.(*security.TrustedRootsConfigV1Alpha1)
			if !ok {
				errs = multierror.Append(errs, fmt.Errorf("expected %s document, got %T", security.TrustedRootsConfig, document))

				continue
			}

			if trustedRootsConfig.MetaName == infraProviderCAName {
				errs = multierror.Append(errs, fmt.Errorf("extra machine config must not contain %s document with name %q", security.TrustedRootsConfig, infraProviderCAName))
			}
		}
	}

	if errs != nil {
		return nil, fmt.Errorf("invalid extra machine config: %w", errs)
	}

	return loadedConfig.Documents(), nil
}
