// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ipmi provides power management functionality using IPMI.
package ipmi

import (
	"context"
	"fmt"
	"time"

	"github.com/bougou/go-ipmi"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/pxe"
)

// Client is a wrapper around the goipmi client.
type Client struct {
	client *ipmi.Client
}

// Close implements the power.Client interface.
func (c *Client) Close() error {
	return c.client.Close()
}

// Reboot implements the power.Client interface.
func (c *Client) Reboot(context.Context) error {
	_, err := c.client.ChassisControl(ipmi.ChassisControlPowerCycle)

	return err
}

// PowerOn implements the power.Client interface.
func (c *Client) PowerOn(context.Context) error {
	_, err := c.client.ChassisControl(ipmi.ChassisControlPowerUp)

	return err
}

// PowerOff implements the power.Client interface.
func (c *Client) PowerOff(context.Context) error {
	_, err := c.client.ChassisControl(ipmi.ChassisControlPowerDown)

	return err
}

// SetPXEBootOnce implements the power.Client interface.
func (c *Client) SetPXEBootOnce(_ context.Context, mode pxe.BootMode) error {
	var biosBootType ipmi.BIOSBootType

	switch mode {
	case pxe.BootModeBIOS:
		biosBootType = ipmi.BIOSBootTypeLegacy
	case pxe.BootModeUEFI:
		biosBootType = ipmi.BIOSBootTypeEFI
	default:
		return fmt.Errorf("unknown boot mode: %v", mode)
	}

	return c.client.SetBootDevice(ipmi.BootDeviceSelectorForcePXE, biosBootType, false)
}

// IsPoweredOn implements the power.Client interface.
func (c *Client) IsPoweredOn(context.Context) (bool, error) {
	status, err := c.client.GetChassisStatus()
	if err != nil {
		return false, err
	}

	return status.PowerIsOn, nil
}

// NewClient creates a new IPMI client.
func NewClient(info *specs.PowerManagement_IPMI) (*Client, error) {
	if info.Port == 0 {
		info.Port = 623
	}

	client, err := ipmi.NewClient(info.Address, int(info.Port), info.Username, info.Password)
	if err != nil {
		return nil, err
	}

	client = client.WithTimeout(30 * time.Second) // todo: rework here, so that context is respected in all calls

	if err = client.Connect(); err != nil {
		client.Close() //nolint:errcheck

		return nil, err
	}

	return &Client{client: client}, nil
}
