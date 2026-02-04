// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package redfish provides BMC functionality using Redfish.
package redfish

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/siderolabs/gen/xslices"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/schemas"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/pxe"
)

// Client is a wrapper around the gofish client.
type Client struct {
	logger                    *zap.Logger
	config                    gofish.ClientConfig
	setBootSourceOverrideMode bool
}

// Close implements the power.Client interface.
func (c *Client) Close(context.Context) error {
	return nil
}

// Reboot implements the power.Client interface.
func (c *Client) Reboot(ctx context.Context) error {
	return c.withClient(ctx, func(client *gofish.APIClient) error {
		return c.doComputerSystemReset(client, schemas.ForceRestartResetType) // todo: consider making reset type configurable
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

		poweredOn = system.PowerState == schemas.OnPowerState

		return nil
	}); err != nil {
		return false, err
	}

	return poweredOn, nil
}

// PowerOn implements the power.Client interface.
func (c *Client) PowerOn(ctx context.Context) error {
	return c.withClient(ctx, func(client *gofish.APIClient) error {
		return c.doComputerSystemReset(client, schemas.OnResetType)
	})
}

// PowerOff implements the power.Client interface.
func (c *Client) PowerOff(ctx context.Context) error {
	return c.withClient(ctx, func(client *gofish.APIClient) error {
		return c.doComputerSystemReset(client, schemas.ForceOffResetType)
	})
}

func (c *Client) doComputerSystemReset(client *gofish.APIClient, resetType schemas.ResetType) error {
	system, err := c.getSystem(client)
	if err != nil {
		return err
	}

	_, err = system.Reset(resetType)

	return err
}

// SetPXEBootOnce implements the power.Client interface.
func (c *Client) SetPXEBootOnce(ctx context.Context, mode pxe.BootMode) error {
	return c.withClient(ctx, func(client *gofish.APIClient) error {
		system, err := c.getSystem(client)
		if err != nil {
			return err
		}

		boot := schemas.Boot{
			BootSourceOverrideEnabled: schemas.OnceBootSourceOverrideEnabled,
			BootSourceOverrideTarget:  schemas.PxeBootSource,
		}

		if c.setBootSourceOverrideMode {
			switch mode {
			case pxe.BootModeBIOS:
				boot.BootSourceOverrideMode = schemas.LegacyBootSourceOverrideMode
			case pxe.BootModeUEFI:
				boot.BootSourceOverrideMode = schemas.UEFIBootSourceOverrideMode
			default:
				return fmt.Errorf("unknown boot mode: %s", mode)
			}
		}

		if err = system.SetBoot(&boot); err != nil {
			if c.isAMIFutureStateError(err) {
				c.logger.Debug("attempting AMI FutureState workaround for boot settings")

				return c.setBootAMIFutureState(client, system, boot)
			}

			return err
		}

		return nil
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

func (c *Client) getSystem(client *gofish.APIClient) (*schemas.ComputerSystem, error) {
	systems, err := client.Service.Systems()
	if err != nil {
		return nil, err
	}

	if len(systems) == 0 {
		return nil, fmt.Errorf("no systems found")
	}

	if len(systems) > 1 {
		ids := xslices.Map(systems, func(system *schemas.ComputerSystem) string {
			return system.ID
		})

		c.logger.Warn("multiple systems found, using first one", zap.Strings("system_ids", ids))
	}

	return systems[0], nil
}

// isAMIFutureStateError checks if the error is the specific AMI error requiring FutureState URI.
func (c *Client) isAMIFutureStateError(err error) bool {
	return strings.Contains(err.Error(), "Ami.1.0.OperationSupportedInFutureStateURI")
}

// setBootAMIFutureState handles boot setting for AMI BMCs using the FutureState URI.
func (c *Client) setBootAMIFutureState(client *gofish.APIClient, system *schemas.ComputerSystem, boot schemas.Boot) error {
	// For AMI BMCs, we need to:
	// 1. GET the current FutureState to obtain ETag
	// 2. PATCH boot settings to /redfish/v1/Systems/{id}/SD (FutureState URI) with If-Match header

	// Construct the FutureState URI
	futureStateURI := system.ODataID + "/SD"

	c.logger.Debug("using AMI FutureState URI for boot settings", zap.String("uri", futureStateURI))

	// First, GET the current FutureState to obtain ETag
	resp, err := client.Get(futureStateURI)
	if err != nil {
		return fmt.Errorf("failed to get current FutureState: %w", err)
	}

	etag := resp.Header.Get("ETag")
	if etag == "" {
		return fmt.Errorf("no ETag found in FutureState response")
	}

	c.logger.Debug("obtained ETag from FutureState", zap.String("etag", etag))

	// PATCH to the FutureState URI with If-Match header
	headers := map[string]string{
		"If-Match": etag,
	}

	// Boot should be a field in the SD object, so we need to wrap it in a Boot object
	// See https://pubs.lenovo.com/tsm/patch_systems_instance_sd for more details
	payload := struct {
		Boot schemas.Boot `json:"Boot"`
	}{Boot: boot}

	_, err = client.PatchWithHeaders(futureStateURI, payload, headers)
	if err != nil {
		return fmt.Errorf("failed to set boot via AMI FutureState URI: %w", err)
	}

	c.logger.Debug("successfully set boot settings via AMI FutureState URI")

	return nil
}

// NewClient returns a new Redfish BMC client.
func NewClient(options Options, address, username, password string, logger *zap.Logger) *Client {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}

	protocol := "http"
	if options.UseHTTPS {
		protocol = "https"
	}

	endpoint := fmt.Sprintf("%s://%s", protocol, net.JoinHostPort(host, strconv.Itoa(options.Port)))

	return &Client{
		config: gofish.ClientConfig{
			Endpoint:  endpoint,
			Username:  username,
			Password:  password,
			Insecure:  options.InsecureSkipTLSVerify,
			BasicAuth: true,
		},

		setBootSourceOverrideMode: options.SetBootSourceOverrideMode,
		logger:                    logger,
	}
}
