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

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/siderolabs/omni/client/pkg/jointoken"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
	siderolinkres "github.com/siderolabs/omni/client/pkg/omni/resources/siderolink"
	"github.com/siderolabs/talos/pkg/machinery/config/config"
	"github.com/siderolabs/talos/pkg/machinery/config/container"
	"github.com/siderolabs/talos/pkg/machinery/config/encoder"
	"github.com/siderolabs/talos/pkg/machinery/config/types/meta"
	"github.com/siderolabs/talos/pkg/machinery/config/types/runtime"
	"github.com/siderolabs/talos/pkg/machinery/config/types/security"
	"github.com/siderolabs/talos/pkg/machinery/config/types/siderolink"

	providermeta "github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/tls"
)

const (
	siderolinkAddress       = "fdae:41e4:649b:9303::1"
	siderolinkEventSinkPort = "8090"
	siderolinkLogPort       = "8092"
)

// Build builds the machine configuration for the bare-metal provider.
func Build(ctx context.Context, r controller.Reader, certs *tls.Certs) ([]byte, error) {
	siderolinkAPIURL, err := getSiderolinkAPIURL(ctx, r)
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
	eventSinkConfig.Endpoint = net.JoinHostPort(siderolinkAddress, siderolinkEventSinkPort)

	kmsgLogURL, err := url.Parse("tcp://" + net.JoinHostPort(siderolinkAddress, siderolinkLogPort))
	if err != nil {
		return nil, fmt.Errorf("failed to parse kmsg log URL: %w", err)
	}

	kmsgLogConfig := runtime.NewKmsgLogV1Alpha1()
	kmsgLogConfig.MetaName = "omni-kmsg"
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

	configContainer, err := container.New(documents...)
	if err != nil {
		return nil, fmt.Errorf("failed to create config container: %w", err)
	}

	return configContainer.EncodeBytes(encoder.WithComments(encoder.CommentsDisabled))
}

func getSiderolinkAPIURL(ctx context.Context, r controller.Reader) (string, error) {
	connectionParams, err := safe.ReaderGetByID[*siderolinkres.ConnectionParams](ctx, r, siderolinkres.ConfigID)
	if err != nil {
		return "", fmt.Errorf("failed to get connection params: %w", err)
	}

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
