// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package omni provides Omni-related functionality.
package omni

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/siderolabs/omni/client/pkg/jointoken"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
	"github.com/siderolabs/omni/client/pkg/omni/resources/siderolink"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
)

// Client is a wrapper around the Omni client.
type Client struct {
	state state.State
}

// BuildClient creates a new Omni client, wrapping the given Omni API client and state.
func BuildClient(state state.State) *Client {
	return &Client{
		state: state,
	}
}

// GetSiderolinkAPIURL returns the SideroLink API URL.
func (c *Client) GetSiderolinkAPIURL(ctx context.Context) (string, error) {
	connectionParams, err := safe.StateGetByID[*siderolink.ConnectionParams](ctx, c.state, siderolink.ConfigID)
	if err != nil {
		return "", fmt.Errorf("failed to get connection params: %w", err)
	}

	token, err := jointoken.NewWithExtraData(connectionParams.TypedSpec().Value.JoinToken, map[string]string{
		omni.LabelInfraProviderID: meta.ProviderID.String(), // go to omni, don't do the check of MachineReqStatus
	})
	if err != nil {
		return "", err
	}

	tokenString, err := token.Encode()
	if err != nil {
		return "", fmt.Errorf("failed to encode the siderolink token: %w", err)
	}

	apiURL, err := siderolink.APIURL(connectionParams, siderolink.WithJoinToken(tokenString))
	if err != nil {
		return "", fmt.Errorf("failed to build API URL: %w", err)
	}

	return apiURL, nil
}

// EnsureProviderStatus makes sure that the infra.ProviderStatus resource exists and is up to date for this provider.
func (c *Client) EnsureProviderStatus(ctx context.Context, name, description string, rawIcon []byte) error {
	populate := func(res *infra.ProviderStatus) {
		res.Metadata().Labels().Set(omni.LabelIsStaticInfraProvider, "")

		res.TypedSpec().Value.Name = name
		res.TypedSpec().Value.Description = description
		res.TypedSpec().Value.Icon = base64.RawStdEncoding.EncodeToString(rawIcon)
	}

	providerStatus := infra.NewProviderStatus(meta.ProviderID.String())

	populate(providerStatus)

	if err := c.state.Create(ctx, providerStatus); err != nil {
		if !state.IsConflictError(err) {
			return err
		}

		if _, err = safe.StateUpdateWithConflicts(ctx, c.state, providerStatus.Metadata(), func(res *infra.ProviderStatus) error {
			populate(res)

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
