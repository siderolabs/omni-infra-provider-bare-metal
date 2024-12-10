// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package api provides power management functionality using an HTTP API, e.g., the HTTP API run by 'talosctl cluster create'.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/pxe"
)

// Client is an API power management client: it communicates with an HTTP API to send power management commands.
type Client struct {
	address string
}

// Close implements the power.Client interface.
func (c *Client) Close() error {
	return nil
}

// Reboot implements the power.Client interface.
func (c *Client) Reboot(ctx context.Context) error {
	return c.doPost(ctx, "/reboot")
}

// PowerOn implements the power.Client interface.
func (c *Client) PowerOn(ctx context.Context) error {
	return c.doPost(ctx, "/poweron")
}

// PowerOff implements the power.Client interface.
func (c *Client) PowerOff(ctx context.Context) error {
	return c.doPost(ctx, "/poweroff")
}

// SetPXEBootOnce implements the power.Client interface.
func (c *Client) SetPXEBootOnce(ctx context.Context, _ pxe.BootMode) error {
	return c.doPost(ctx, "/pxeboot")
}

func (c *Client) doPost(ctx context.Context, path string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	endpoint := "http://" + c.address + path

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request %q: %w", path, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request %q: %w", path, err)
	}

	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code while resetting machine: %d", resp.StatusCode)
	}

	return nil
}

// IsPoweredOn implements the power.Client interface.
func (c *Client) IsPoweredOn(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	endpoint := "http://" + c.address + "/status"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	defer closeBody(resp)

	var status struct {
		PoweredOn bool
	}

	if err = json.NewDecoder(resp.Body).Decode(&status); err != nil { //nolint:musttag
		return false, err
	}

	return status.PoweredOn, nil
}

//nolint:errcheck
func closeBody(resp *http.Response) {
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

// NewClient creates a new API power management client.
func NewClient(info *specs.PowerManagement_API) (*Client, error) {
	return &Client{address: info.Address}, nil
}
