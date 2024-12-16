// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package redfish provides power management functionality using Redfish.
package redfish

import (
	"context"
	"fmt"
	"time"

	"github.com/siderolabs/gen/xslices"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/pxe"
)

// Client is a wrapper around the gofish client.
type Client struct {
	logger                    *zap.Logger
	config                    gofish.ClientConfig
	setBootSourceOverrideMode bool
}

// Close implements the power.Client interface.
func (c *Client) Close() error {
	return nil
}

// Reboot implements the power.Client interface.
func (c *Client) Reboot(ctx context.Context) error {
	return c.withClient(ctx, func(client *gofish.APIClient) error {
		return c.doComputerSystemReset(client, redfish.ForceRestartResetType) // todo: consider making reset type configurable
	})
}

// IsPoweredOn implements the power.Client interface.
func (c *Client) IsPoweredOn(ctx context.Context) (bool, error) {
	poweredOn := false

	if err := c.withClient(ctx, func(client *gofish.APIClient) error {
		system, err := c.getSystem(client)
		if err != nil {
			return err
		}

		poweredOn = system.PowerState == redfish.OnPowerState

		return nil
	}); err != nil {
		return false, err
	}

	return poweredOn, nil
}

// PowerOn implements the power.Client interface.
func (c *Client) PowerOn(ctx context.Context) error {
	return c.withClient(ctx, func(client *gofish.APIClient) error {
		return c.doComputerSystemReset(client, redfish.OnResetType)
	})
}

// PowerOff implements the power.Client interface.
func (c *Client) PowerOff(ctx context.Context) error {
	return c.withClient(ctx, func(client *gofish.APIClient) error {
		return c.doComputerSystemReset(client, redfish.ForceOffResetType)
	})
}

func (c *Client) doComputerSystemReset(client *gofish.APIClient, resetType redfish.ResetType) error {
	system, err := c.getSystem(client)
	if err != nil {
		return err
	}

	return system.Reset(resetType)
}

// SetPXEBootOnce implements the power.Client interface.
func (c *Client) SetPXEBootOnce(ctx context.Context, mode pxe.BootMode) error {
	return c.withClient(ctx, func(client *gofish.APIClient) error {
		system, err := c.getSystem(client)
		if err != nil {
			return err
		}

		boot := redfish.Boot{
			BootSourceOverrideEnabled: redfish.OnceBootSourceOverrideEnabled,
			BootSourceOverrideTarget:  redfish.PxeBootSourceOverrideTarget,
		}

		if c.setBootSourceOverrideMode {
			switch mode {
			case pxe.BootModeBIOS:
				boot.BootSourceOverrideMode = redfish.LegacyBootSourceOverrideMode
			case pxe.BootModeUEFI:
				boot.BootSourceOverrideMode = redfish.UEFIBootSourceOverrideMode
			default:
				return fmt.Errorf("unknown boot mode: %s", mode)
			}
		}

		return system.SetBoot(boot)
	})
}

func (c *Client) withClient(ctx context.Context, f func(client *gofish.APIClient) error) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := gofish.ConnectContext(ctx, c.config)
	if err != nil {
		return err
	}

	defer client.Logout()

	return f(client)
}

func (c *Client) getSystem(client *gofish.APIClient) (*redfish.ComputerSystem, error) {
	systems, err := client.Service.Systems()
	if err != nil {
		return nil, err
	}

	if len(systems) == 0 {
		return nil, fmt.Errorf("no systems found")
	}

	if len(systems) > 1 {
		ids := xslices.Map(systems, func(system *redfish.ComputerSystem) string {
			return system.ID
		})

		c.logger.Warn("multiple systems found, using first one", zap.Strings("system_ids", ids))
	}

	return systems[0], nil
}

// NewClient returns a new Redfish power management client.
func NewClient(address, username, password string, setBootSourceOverrideMode bool, logger *zap.Logger) *Client {
	return &Client{
		config: gofish.ClientConfig{
			Endpoint:            "https://" + address,
			Username:            username,
			Password:            password,
			Insecure:            true,
			TLSHandshakeTimeout: 5, // seconds
			BasicAuth:           true,
		},
		setBootSourceOverrideMode: setBootSourceOverrideMode,
		logger:                    logger,
	}
}
