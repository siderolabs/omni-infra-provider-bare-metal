// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package bmc provides BMC functionality for machines.
package bmc

import (
	"context"
	"fmt"
	"io"
	"sync"

	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/api"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/ipmi"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/pxe"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/redfish"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
)

// Client is the interface to interact with a single machine to send BMC commands to it.
type Client interface {
	io.Closer
	Reboot(ctx context.Context) error
	IsPoweredOn(ctx context.Context) (bool, error)
	PowerOn(ctx context.Context) error
	PowerOff(ctx context.Context) error
	SetPXEBootOnce(ctx context.Context, mode pxe.BootMode) error
}

// ClientFactory is a factory to create BMC clients.
type ClientFactory struct {
	addressToRedfishAvailability   map[string]bool
	options                        ClientFactoryOptions
	addressToRedfishAvailabilityMu sync.Mutex
}

// ClientFactoryOptions contains options for the client factory.
type ClientFactoryOptions struct {
	RedfishOptions redfish.Options
}

// NewClientFactory creates a new BMC client factory.
func NewClientFactory(options ClientFactoryOptions) *ClientFactory {
	return &ClientFactory{
		options:                      options,
		addressToRedfishAvailability: map[string]bool{},
	}
}

// GetClient returns a BMC client for the given bare metal machine.
func (factory *ClientFactory) GetClient(ctx context.Context, config *resources.BMCConfiguration, logger *zap.Logger) (Client, error) {
	if config == nil {
		return nil, fmt.Errorf("cannot get BMC client: config is nil")
	}

	spec := config.TypedSpec().Value

	if spec.Ipmi == nil && spec.Api == nil {
		return nil, fmt.Errorf("invalid BMC config: both IPMI and API fields are nil")
	}

	if spec.Api != nil {
		apiClient, err := api.NewClient(spec.Api)
		if err != nil {
			return nil, err
		}

		return &loggingClient{client: apiClient, logger: logger.With(zap.String("bmc_client", "api"))}, nil
	}

	useRedfish := factory.options.RedfishOptions.UseAlways || (factory.options.RedfishOptions.UseWhenAvailable && factory.redfishAvailable(ctx, spec.Ipmi, logger))

	if useRedfish {
		logger = logger.With(zap.String("bmc_client", "redfish"))
		redfishClient := redfish.NewClient(factory.options.RedfishOptions, spec.Ipmi.Address, spec.Ipmi.Username, spec.Ipmi.Password, logger)

		return &loggingClient{client: redfishClient, logger: logger}, nil
	}

	ipmiClient, err := ipmi.NewClient(ctx, spec.Ipmi)
	if err != nil {
		return nil, err
	}

	return &loggingClient{client: ipmiClient, logger: logger.With(zap.String("bmc_client", "ipmi"))}, nil
}

func (factory *ClientFactory) redfishAvailable(ctx context.Context, ipmiInfo *specs.BMCConfigurationSpec_IPMI, logger *zap.Logger) bool {
	factory.addressToRedfishAvailabilityMu.Lock()
	defer factory.addressToRedfishAvailabilityMu.Unlock()

	address := ipmiInfo.Address

	available, ok := factory.addressToRedfishAvailability[address]
	if ok {
		return available
	}

	logger.Debug("probe redfish availability", zap.String("address", address))

	redfishClient := redfish.NewClient(factory.options.RedfishOptions, address, ipmiInfo.Username, ipmiInfo.Password, logger)

	if _, err := redfishClient.IsPoweredOn(ctx); err != nil {
		logger.Debug("redfish is not available on address", zap.String("address", address), zap.Error(err))

		factory.addressToRedfishAvailability[address] = false

		return false
	}

	logger.Debug("redfish is available on address", zap.String("address", address))

	factory.addressToRedfishAvailability[address] = true

	return true
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
