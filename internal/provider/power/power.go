// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package power provides power management functionality for machines.
package power

import (
	"context"
	"fmt"
	"io"

	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/api"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/ipmi"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/pxe"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/redfish"
)

// ErrNoPowerManagementInfo is returned when there is no power management info present yet for a machine.
var ErrNoPowerManagementInfo = fmt.Errorf("no power management info found")

// Client is the interface to interact with a single machine to send power commands to it.
type Client interface {
	io.Closer
	Reboot(ctx context.Context) error
	IsPoweredOn(ctx context.Context) (bool, error)
	PowerOn(ctx context.Context) error
	PowerOff(ctx context.Context) error
	SetPXEBootOnce(ctx context.Context, mode pxe.BootMode) error
}

// ClientFactory is a factory to create power management clients.
type ClientFactory struct {
	logger  *zap.Logger
	options ClientFactoryOptions
}

// ClientFactoryOptions contains options for the client factory.
type ClientFactoryOptions struct {
	ExperimentalUseRedfish           bool
	RedfishSetBootSourceOverrideMode bool
}

// NewClientFactory creates a new power management client factory.
func NewClientFactory(options ClientFactoryOptions, logger *zap.Logger) *ClientFactory {
	return &ClientFactory{
		options: options,
		logger:  logger,
	}
}

// GetClient returns a power management client for the given bare metal machine.
func (factory *ClientFactory) GetClient(mgmt *specs.PowerManagement) (Client, error) {
	if mgmt == nil {
		return nil, ErrNoPowerManagementInfo
	}

	apiInfo := mgmt.Api
	if apiInfo != nil {
		client, err := api.NewClient(apiInfo)
		if err != nil {
			return nil, err
		}

		return &loggingClient{client: client, logger: factory.logger.With(zap.String("power_client", "api"))}, nil
	}

	ipmiInfo := mgmt.Ipmi
	if ipmiInfo != nil {
		if factory.options.ExperimentalUseRedfish {
			logger := factory.logger.With(zap.String("power_client", "redfish"))
			client := redfish.NewClient(ipmiInfo.Address, ipmiInfo.Username, ipmiInfo.Password, factory.options.RedfishSetBootSourceOverrideMode, logger)

			return &loggingClient{client: client, logger: logger}, nil
		}

		client, err := ipmi.NewClient(ipmiInfo)
		if err != nil {
			return nil, err
		}

		return &loggingClient{client: client, logger: factory.logger.With(zap.String("power_client", "ipmi"))}, nil
	}

	return nil, ErrNoPowerManagementInfo
}

type loggingClient struct {
	client Client
	logger *zap.Logger
}

func (client *loggingClient) Close() error {
	client.logger.Debug("close client")

	return client.client.Close()
}

func (client *loggingClient) Reboot(ctx context.Context) error {
	client.logger.Debug("reboot")

	return client.client.Reboot(ctx)
}

func (client *loggingClient) IsPoweredOn(ctx context.Context) (bool, error) {
	client.logger.Debug("is powered on")

	return client.client.IsPoweredOn(ctx)
}

func (client *loggingClient) PowerOn(ctx context.Context) error {
	client.logger.Debug("power on")

	return client.client.PowerOn(ctx)
}

func (client *loggingClient) PowerOff(ctx context.Context) error {
	client.logger.Debug("power off")

	return client.client.PowerOff(ctx)
}

func (client *loggingClient) SetPXEBootOnce(ctx context.Context, mode pxe.BootMode) error {
	client.logger.Debug("set PXE boot once", zap.String("mode", string(mode)))

	return client.client.SetPXEBootOnce(ctx, mode)
}
