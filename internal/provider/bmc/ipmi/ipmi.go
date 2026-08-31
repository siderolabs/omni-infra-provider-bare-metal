// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ipmi provides BMC functionality using IPMI.
package ipmi

import (
	"context"
	"fmt"
	"time"

	"github.com/bougou/go-ipmi/pkg/client"
	"github.com/bougou/go-ipmi/pkg/command/chassis"
	"github.com/bougou/go-ipmi/pkg/types"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/pxe"
)

const timeout = 30 * time.Second

// Client is a wrapper around the goipmi client.
type Client struct {
	ipmiClient *client.Client
}

// Close implements the power.Client interface.
func (c *Client) Close(ctx context.Context) error {
	return c.ipmiClient.Close(ctx)
}

// Reboot implements the power.Client interface.
func (c *Client) Reboot(ctx context.Context) error {
	_, err := c.ipmiClient.ChassisControl(ctx, chassis.ChassisControlPowerCycle)

	return err
}

// PowerOn implements the power.Client interface.
func (c *Client) PowerOn(ctx context.Context) error {
	_, err := c.ipmiClient.ChassisControl(ctx, chassis.ChassisControlPowerUp)

	return err
}

// PowerOff implements the power.Client interface.
func (c *Client) PowerOff(ctx context.Context) error {
	_, err := c.ipmiClient.ChassisControl(ctx, chassis.ChassisControlPowerDown)

	return err
}

// SetPXEBootOnce implements the power.Client interface.
func (c *Client) SetPXEBootOnce(ctx context.Context, mode pxe.BootMode) error {
	var bootType types.BIOSBootType

	switch mode {
	case pxe.BootModeBIOS:
		bootType = types.BIOSBootTypeLegacy
	case pxe.BootModeUEFI:
		bootType = types.BIOSBootTypeEFI
	default:
		return fmt.Errorf("unsupported mode %q", mode)
	}

	return c.ipmiClient.SetBootDevice(ctx, types.BootDeviceSelectorForcePXE, bootType, false)
}

// ResetBootDevice clears any boot device override, resetting it to the default boot order.
func (c *Client) ResetBootDevice(ctx context.Context) error {
	return c.ipmiClient.SetBootDevice(ctx, types.BootDeviceSelectorNoOverride, types.BIOSBootTypeEFI, false)
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
	ipmiClient, err := client.NewClient(info.Address, int(info.Port), info.Username, info.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create IPMI client: %w", err)
	}

	ipmiClient = ipmiClient.WithTimeout(timeout)

	if err = ipmiClient.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect IPMI client: %w", err)
	}

	return &Client{
		ipmiClient: ipmiClient,
	}, nil
}
