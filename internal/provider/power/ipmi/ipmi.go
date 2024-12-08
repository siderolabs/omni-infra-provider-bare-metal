// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ipmi provides power management functionality using IPMI.
package ipmi

import (
	"context"
	"fmt"

	goipmi "github.com/pensando/goipmi"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/pxe"
)

const ipmiUsername = "talos-agent"

// Client is a wrapper around the goipmi client.
type Client struct {
	ipmiClient *goipmi.Client
}

// Close implements the power.Client interface.
func (c *Client) Close() error {
	return c.ipmiClient.Close()
}

// Reboot implements the power.Client interface.
func (c *Client) Reboot(context.Context) error {
	return c.ipmiClient.Control(goipmi.ControlPowerCycle)
}

// PowerOff implements the power.Client interface.
func (c *Client) PowerOff(context.Context) error {
	return c.ipmiClient.Control(goipmi.ControlPowerDown)
}

// SetPXEBootOnce implements the power.Client interface.
func (c *Client) SetPXEBootOnce(_ context.Context, mode pxe.BootMode) error {
	switch mode {
	case pxe.BootModeBIOS:
		return c.ipmiClient.SetBootDevice(goipmi.BootDevicePxe)
	case pxe.BootModeUEFI:
		return c.ipmiClient.SetBootDeviceEFI(goipmi.BootDevicePxe)
	default:
		return fmt.Errorf("unsupported mode %q", mode)
	}
}

// IsPoweredOn implements the power.Client interface.
func (c *Client) IsPoweredOn(context.Context) (bool, error) {
	req := &goipmi.Request{
		NetworkFunction: goipmi.NetworkFunctionChassis,
		Command:         goipmi.CommandChassisStatus,
		Data:            goipmi.ChassisStatusRequest{},
	}

	res := &goipmi.ChassisStatusResponse{}

	err := c.ipmiClient.Send(req, res)
	if err != nil {
		return false, err
	}

	return res.IsSystemPowerOn(), nil
}

// NewClient creates a new IPMI client.
func NewClient(info *specs.PowerManagement_IPMI) (*Client, error) {
	conn := &goipmi.Connection{
		Hostname:  info.Address,
		Port:      int(info.Port),
		Username:  ipmiUsername,
		Password:  info.Password,
		Interface: "lanplus",
	}

	client, err := goipmi.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return &Client{ipmiClient: client}, nil
}
