// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package provider implements the bare metal infra provider.
package provider

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"time"

	"github.com/cosi-project/runtime/pkg/controller/runtime"
	runtimeoptions "github.com/cosi-project/runtime/pkg/controller/runtime/options"
	"github.com/cosi-project/runtime/pkg/resource/protobuf"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/hashicorp/go-multierror"
	"github.com/siderolabs/omni/client/pkg/client"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/agent"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/baremetal"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/config"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/controllers"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/dhcp"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/imagefactory"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/ip"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/ipxe"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/machinestatus"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/omni"
	powerapi "github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/api"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/server"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/tftp"
)

//go:embed data/icon.svg
var icon []byte

// Provider implements the bare metal infra provider.
type Provider struct {
	logger *zap.Logger

	options Options
}

// New creates a new Provider.
func New(options Options, logger *zap.Logger) *Provider {
	return &Provider{
		options: options,
		logger:  logger,
	}
}

// Run runs the provider.
func (p *Provider) Run(ctx context.Context) error {
	apiAdvertiseAddress, err := p.determineAPIAdvertiseAddress()
	if err != nil {
		return fmt.Errorf("failed to determine API advertise address: %w", err)
	}

	p.logger.Info("starting provider",
		zap.String("api_listen_address", p.options.APIListenAddress),
		zap.String("api_advertise_address", apiAdvertiseAddress),
		zap.Int("api_port", p.options.APIPort))

	omniAPIClient, err := p.buildOmniAPIClient(p.options.OmniAPIEndpoint, p.options.InsecureSkipTLSVerify)
	if err != nil {
		return fmt.Errorf("failed to build omni client: %w", err)
	}

	defer omniAPIClient.Close() //nolint:errcheck

	cosiRuntime, err := p.buildCOSIRuntime(omniAPIClient)
	if err != nil {
		return fmt.Errorf("failed to build COSI runtime: %w", err)
	}

	omniState := state.WrapCore(cosiRuntime.CachedState())
	omniClient := omni.BuildClient(omniState)

	if p.options.ClearState {
		if err = p.clearState(ctx, omniState); err != nil {
			return fmt.Errorf("failed to clear state: %w", err)
		}

		p.logger.Info("state cleared")
	}

	if err = omniClient.EnsureProviderStatus(ctx, p.options.Name, p.options.Description, icon); err != nil {
		return fmt.Errorf("failed to create/update provider status: %w", err)
	}

	imageFactoryClient, err := imagefactory.NewClient(p.options.ImageFactoryBaseURL, p.options.ImageFactoryPXEBaseURL, p.options.AgentModeTalosVersion, p.options.MachineLabels)
	if err != nil {
		return fmt.Errorf("failed to create image factory client: %w", err)
	}

	ipxeHandler, err := ipxe.NewHandler(imageFactoryClient, omniState, ipxe.HandlerOptions{
		APIAdvertiseAddress: apiAdvertiseAddress,
		APIPort:             p.options.APIPort,
		UseLocalBootAssets:  p.options.UseLocalBootAssets,
		AgentTestMode:       p.options.AgentTestMode,
		BootFromDiskMethod:  p.options.BootFromDiskMethod,
	}, p.logger.With(zap.String("component", "ipxe_handler")))
	if err != nil {
		return fmt.Errorf("failed to create iPXE handler: %w", err)
	}

	configHandler, err := config.NewHandler(ctx, omniClient, p.logger.With(zap.String("component", "config_handler")))
	if err != nil {
		return fmt.Errorf("failed to create config handler: %w", err)
	}

	srvr := server.New(ctx, p.options.APIListenAddress, p.options.APIPort, p.options.UseLocalBootAssets, configHandler, ipxeHandler, p.logger.With(zap.String("component", "server")))
	agentController := agent.NewController(srvr, omniState, p.options.WipeWithZeroes, p.logger.With(zap.String("component", "controller"))) //nolint:contextcheck // false positive
	machineStatusPoller := machinestatus.NewPoller(agentController, omniState, p.logger.With(zap.String("component", "machine_status_poller")))
	dhcpProxy := dhcp.NewProxy(apiAdvertiseAddress, p.options.APIPort, p.options.DHCPProxyIfaceOrIP, p.logger.With(zap.String("component", "dhcp_proxy")))
	tftpServer := tftp.NewServer(p.logger.With(zap.String("component", "tftp_server")))
	apiPowerManager := powerapi.NewPowerManager(p.options.APIPowerMgmtStateDir)

	// todo: enable if we re-enable reverse tunnel on Omni: https://github.com/siderolabs/omni/pull/746
	// reverseTunnel := tunnel.New(omniState, omniAPIClient, p.logger.With(zap.String("component", "reverse_tunnel")))

	if err = cosiRuntime.RegisterQController(controllers.NewInfraMachineController(agentController, apiPowerManager, omniState, 1*time.Minute)); err != nil {
		return fmt.Errorf("failed to register controller: %w", err)
	}

	return p.runComponents(ctx, map[string]func(context.Context) error{
		"COSI runtime":          cosiRuntime.Run,
		"machine status poller": machineStatusPoller.Run,
		"server":                srvr.Run,
		// "reverse tunnel":        reverseTunnel.Run,
		"DHCP proxy":  dhcpProxy.Run,
		"TFTP server": tftpServer.Run,
	})
}

func (p *Provider) buildCOSIRuntime(omniAPIClient *client.Client) (*runtime.Runtime, error) {
	omniState := omniAPIClient.Omni().State()

	if err := protobuf.RegisterResource(baremetal.MachineStatusType, &baremetal.MachineStatus{}); err != nil {
		return nil, fmt.Errorf("failed to register protobuf resource: %w", err)
	}

	var options []runtimeoptions.Option

	if p.options.EnableResourceCache {
		options = append(options, safe.WithResourceCache[*baremetal.MachineStatus]())
		options = append(options, safe.WithResourceCache[*infra.Machine]())
		options = append(options, safe.WithResourceCache[*infra.MachineStatus]())
	}

	cosiRuntime, err := runtime.NewRuntime(omniState, p.logger.With(zap.String("component", "cosi_runtime")), options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create runtime: %w", err)
	}

	return cosiRuntime, nil
}

func (p *Provider) runComponents(ctx context.Context, components map[string]func(context.Context) error) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	for name, f := range components {
		eg.Go(func() error {
			defer cancel() // cancel the parent context, so all other components are also stopped

			p.logger.Info("start component ", zap.String("name", name))

			err := f(ctx)
			if err != nil {
				p.logger.Error("failed to run component", zap.String("name", name), zap.Error(err))

				return err
			}

			p.logger.Info("component stopped", zap.String("name", name))

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to run components: %w", err)
	}

	return nil
}

func (p *Provider) determineAPIAdvertiseAddress() (string, error) {
	if p.options.APIAdvertiseAddress != "" {
		return p.options.APIAdvertiseAddress, nil
	}

	routableIPs, err := ip.RoutableIPs()
	if err != nil {
		return "", fmt.Errorf("failed to get routable IPs: %w", err)
	}

	if len(routableIPs) != 1 {
		return "", fmt.Errorf(`expected exactly one routable IP, got %d: %v. specify API advertise address explicitly`, len(routableIPs), routableIPs)
	}

	return routableIPs[0], nil
}

// buildOmniAPIClient creates a new Omni API client.
func (p *Provider) buildOmniAPIClient(endpoint string, insecureSkipTLSVerify bool) (*client.Client, error) {
	serviceAccountKey := os.Getenv("OMNI_SERVICE_ACCOUNT_KEY")

	cliOpts := []client.Option{
		client.WithInsecureSkipTLSVerify(insecureSkipTLSVerify),
	}

	if serviceAccountKey != "" {
		cliOpts = append(cliOpts, client.WithServiceAccount(serviceAccountKey))
	}

	return client.New(endpoint, cliOpts...)
}

// clearState clears the persistent state of this provider. Useful for debugging purposes.
func (p *Provider) clearState(ctx context.Context, st state.State) error {
	list, err := st.List(ctx, baremetal.NewMachineStatus("").Metadata())
	if err != nil {
		return fmt.Errorf("failed to list bare metal machinees: %w", err)
	}

	var errs error

	for _, item := range list.Items {
		res, getErr := st.Get(ctx, item.Metadata())
		if getErr != nil {
			errs = multierror.Append(errs, getErr)

			continue
		}

		if destroyErr := st.Destroy(ctx, item.Metadata(), state.WithDestroyOwner(res.Metadata().Owner())); destroyErr != nil {
			errs = multierror.Append(errs, destroyErr)

			continue
		}
	}

	return errs
}
