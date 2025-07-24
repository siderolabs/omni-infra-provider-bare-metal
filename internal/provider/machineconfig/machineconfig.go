// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package machineconfig builds the machine configuration for the bare-metal provider.
package machineconfig

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/hashicorp/go-multierror"
	"github.com/siderolabs/omni/client/pkg/jointoken"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
	siderolinkres "github.com/siderolabs/omni/client/pkg/omni/resources/siderolink"
	"github.com/siderolabs/talos/pkg/machinery/config/config"
	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
	"github.com/siderolabs/talos/pkg/machinery/config/container"
	"github.com/siderolabs/talos/pkg/machinery/config/encoder"
	"github.com/siderolabs/talos/pkg/machinery/config/types/meta"
	"github.com/siderolabs/talos/pkg/machinery/config/types/runtime"
	"github.com/siderolabs/talos/pkg/machinery/config/types/security"
	"github.com/siderolabs/talos/pkg/machinery/config/types/siderolink"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"

	providermeta "github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/tls"
)

const (
	siderolinkAddress   = "fdae:41e4:649b:9303::1"
	kmsgLogName         = "omni-kmsg"
	infraProviderCAName = "infra-provider-ca"
)

// Build builds the machine configuration for the bare-metal provider.
func Build(ctx context.Context, r controller.Reader, certs *tls.Certs, extraMachineConfigPath string) ([]byte, error) {
	connectionParams, err := safe.ReaderGetByID[*siderolinkres.ConnectionParams](ctx, r, siderolinkres.ConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection params: %w", err)
	}

	siderolinkAPIURL, err := getSiderolinkAPIURL(connectionParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get siderolink API URL: %w", err)
	}

	apiURL, err := url.Parse(siderolinkAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %w", err)
	}

	siderolinkConfig := siderolink.NewConfigV1Alpha1()
	siderolinkConfig.APIUrlConfig = meta.URL{
		URL: apiURL,
	}

	eventSinkConfig := runtime.NewEventSinkV1Alpha1()
	eventSinkConfig.Endpoint = net.JoinHostPort(siderolinkAddress, strconv.Itoa(int(connectionParams.TypedSpec().Value.EventsPort)))

	kmsgLogURL, err := url.Parse("tcp://" + net.JoinHostPort(siderolinkAddress, strconv.Itoa(int(connectionParams.TypedSpec().Value.LogsPort))))
	if err != nil {
		return nil, fmt.Errorf("failed to parse kmsg log URL: %w", err)
	}

	kmsgLogConfig := runtime.NewKmsgLogV1Alpha1()
	kmsgLogConfig.MetaName = kmsgLogName
	kmsgLogConfig.KmsgLogURL = meta.URL{
		URL: kmsgLogURL,
	}

	documents := []config.Document{
		siderolinkConfig,
		eventSinkConfig,
		kmsgLogConfig,
	}

	if certs != nil {
		trustedRootsConfig := security.NewTrustedRootsConfigV1Alpha1()
		trustedRootsConfig.MetaName = "infra-provider-ca"
		trustedRootsConfig.Certificates = certs.CACertPEM

		documents = append(documents, trustedRootsConfig)
	}

	if extraMachineConfigPath != "" {
		var additionalDocuments []config.Document

		if additionalDocuments, err = parseAdditionalDocuments(extraMachineConfigPath); err != nil {
			return nil, fmt.Errorf("failed to load extra machine config: %w", err)
		}

		documents = append(documents, additionalDocuments...)
	}

	configContainer, err := container.New(documents...)
	if err != nil {
		return nil, fmt.Errorf("failed to create config container: %w", err)
	}

	return configContainer.EncodeBytes(encoder.WithComments(encoder.CommentsDisabled))
}

func getSiderolinkAPIURL(connectionParams *siderolinkres.ConnectionParams) (string, error) {
	token, err := jointoken.NewWithExtraData(connectionParams.TypedSpec().Value.JoinToken, map[string]string{
		omni.LabelInfraProviderID: providermeta.ProviderID.String(), // go to omni, don't do the check of MachineReqStatus
	})
	if err != nil {
		return "", fmt.Errorf("failed to create siderolink token: %w", err)
	}

	tokenString, err := token.Encode()
	if err != nil {
		return "", fmt.Errorf("failed to encode the siderolink token: %w", err)
	}

	apiURL, err := siderolinkres.APIURL(connectionParams, siderolinkres.WithJoinToken(tokenString))
	if err != nil {
		return "", fmt.Errorf("failed to build API URL: %w", err)
	}

	return apiURL, nil
}

func parseAdditionalDocuments(extraMachineConfigPath string) ([]config.Document, error) {
	loadedConfig, err := configloader.NewFromFile(extraMachineConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load extra machine config: %w", err)
	}

	var errs error

	for _, document := range loadedConfig.Documents() {
		switch document.Kind() {
		case v1alpha1.Version, siderolink.Kind, runtime.EventSinkKind:
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
