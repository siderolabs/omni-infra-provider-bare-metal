// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ipmi provides BMC functionality using IPMI.
package ipmi

import (
	"context"
	"fmt"
	"time"

	"github.com/bougou/go-ipmi"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/pxe"
)

const (
	ipmiUsername = "talos-agent"
	timeout      = 30 * time.Second
)

// Client is a wrapper around the goipmi client.
type Client struct {
	ipmiClient *ipmi.Client
}

// Close implements the power.Client interface.
func (c *Client) Close(ctx context.Context) error {
	return c.ipmiClient.Close(ctx)
}

// Reboot implements the power.Client interface.
func (c *Client) Reboot(ctx context.Context) error {
	_, err := c.ipmiClient.ChassisControl(ctx, ipmi.ChassisControlPowerCycle)

	return err
}

// PowerOn implements the power.Client interface.
func (c *Client) PowerOn(ctx context.Context) error {
	_, err := c.ipmiClient.ChassisControl(ctx, ipmi.ChassisControlPowerUp)

	return err
}

// PowerOff implements the power.Client interface.
func (c *Client) PowerOff(ctx context.Context) error {
	_, err := c.ipmiClient.ChassisControl(ctx, ipmi.ChassisControlPowerDown)

	return err
}

// SetPXEBootOnce implements the power.Client interface.
func (c *Client) SetPXEBootOnce(ctx context.Context, mode pxe.BootMode) error {
	var bootType ipmi.BIOSBootType

	switch mode {
	case pxe.BootModeBIOS:
		bootType = ipmi.BIOSBootTypeLegacy
	case pxe.BootModeUEFI:
		bootType = ipmi.BIOSBootTypeEFI
	default:
		return fmt.Errorf("unsupported mode %q", mode)
	}

	return c.ipmiClient.SetBootDevice(ctx, ipmi.BootDeviceSelectorForcePXE, bootType, false)
}

// IsPoweredOn implements the power.Client interface.
func (c *Client) IsPoweredOn(ctx context.Context) (bool, error) {
	resp, err := c.ipmiClient.GetChassisStatus(ctx)
	if err != nil {
		return false, err
	}

	return resp.PowerIsOn, nil
}

// NewClient creates a new IPMI client and connects to the BMC using the provided configuration.
//
// It needs to be closed after use to release resources.
func NewClient(ctx context.Context, info *specs.BMCConfigurationSpec_IPMI) (*Client, error) {
	client, err := ipmi.NewClient(info.Address, int(info.Port), ipmiUsername, info.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create IPMI client: %w", err)
	}

	client = client.WithTimeout(timeout)

	if err = client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect IPMI client: %w", err)
	}

	return &Client{
		ipmiClient: client,
	}, nil
}
