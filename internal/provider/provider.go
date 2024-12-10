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
	"strings"
	"time"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/runtime"
	runtimeoptions "github.com/cosi-project/runtime/pkg/controller/runtime/options"
	"github.com/cosi-project/runtime/pkg/resource/protobuf"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/hashicorp/go-multierror"
	"github.com/siderolabs/omni/client/pkg/client"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/constants"
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
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/pxe"
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
	pxeBootMode, err := pxe.ParseBootMode(p.options.IPMIPXEBootMode)
	if err != nil {
		return fmt.Errorf("failed to parse IPMI PXE boot mode: %w", err)
	}

	apiAdvertiseAddress, err := p.determineAPIAdvertiseAddress()
	if err != nil {
		return fmt.Errorf("failed to determine API advertise address: %w", err)
	}

	dhcpProxyIfaceOrIP := p.options.DHCPProxyIfaceOrIP
	if dhcpProxyIfaceOrIP == "" {
		dhcpProxyIfaceOrIP = apiAdvertiseAddress
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
		if constants.IsDebugBuild {
			if err = p.clearState(ctx, omniState); err != nil {
				return fmt.Errorf("failed to clear state: %w", err)
			}

			p.logger.Info("state cleared")
		} else {
			p.logger.Warn("clear state is requested, but this is not a debug build, skipping")
		}
	}

	if err = omniClient.EnsureProviderStatus(ctx, p.options.Name, p.options.Description, icon); err != nil {
		return fmt.Errorf("failed to create/update provider status: %w", err)
	}

	imageFactoryClient, err := imagefactory.NewClient(p.options.ImageFactoryBaseURL, p.options.ImageFactoryPXEBaseURL, p.options.AgentModeTalosVersion)
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

	parsedMachineLabels, err := p.parseLabels()
	if err != nil {
		return fmt.Errorf("failed to parse machine labels: %w", err)
	}

	srvr := server.New(ctx, p.options.APIListenAddress, p.options.APIPort, p.options.UseLocalBootAssets, configHandler, ipxeHandler, p.logger.With(zap.String("component", "server")))
	agentService := agent.NewService(srvr, omniState, p.options.WipeWithZeroes, p.logger.With(zap.String("component", "agent_service"))) //nolint:contextcheck // false positive
	machineStatusPoller := machinestatus.NewPoller(agentService, omniState, p.logger.With(zap.String("component", "machine_status_poller")))
	dhcpProxy := dhcp.NewProxy(apiAdvertiseAddress, p.options.APIPort, dhcpProxyIfaceOrIP, p.logger.With(zap.String("component", "dhcp_proxy")))
	tftpServer := tftp.NewServer(p.logger.With(zap.String("component", "tftp_server")))
	apiPowerManager := powerapi.NewPowerManager(p.options.APIPowerMgmtStateDir)

	// todo: enable if we re-enable reverse tunnel on Omni: https://github.com/siderolabs/omni/pull/746
	// reverseTunnel := tunnel.New(omniState, omniAPIClient, p.logger.With(zap.String("component", "reverse_tunnel")))

	for _, qController := range []controller.QController{
		controllers.NewInfraMachineStatusController(agentService, apiPowerManager, omniState, pxeBootMode, 1*time.Minute, p.options.MinRebootInterval, parsedMachineLabels),
		controllers.NewPowerStatusController(omniState),
	} {
		if err = cosiRuntime.RegisterQController(qController); err != nil {
			return fmt.Errorf("failed to register QController: %w", err)
		}
	}

	return p.runComponents(ctx, []component{
		{cosiRuntime.Run, "COSI runtime"},
		{machineStatusPoller.Run, "machine status poller"},
		{srvr.Run, "server"},
		{dhcpProxy.Run, "DHCP proxy"},
		{tftpServer.Run, "TFTP server"},
		// {reverseTunnel.Run, "reverse tunnel"},
	})
}

func (p *Provider) buildCOSIRuntime(omniAPIClient *client.Client) (*runtime.Runtime, error) {
	omniState := omniAPIClient.Omni().State()

	var options []runtimeoptions.Option

	if err := protobuf.RegisterResource(baremetal.MachineStatusType, &baremetal.MachineStatus{}); err != nil {
		return nil, fmt.Errorf("failed to register protobuf resource: %w", err)
	}

	if err := protobuf.RegisterResource(baremetal.PowerStatusType, &baremetal.PowerStatus{}); err != nil {
		return nil, fmt.Errorf("failed to register protobuf resource: %w", err)
	}

	if p.options.EnableResourceCache {
		options = append(options,
			safe.WithResourceCache[*baremetal.MachineStatus](),
			safe.WithResourceCache[*baremetal.PowerStatus](),
		)
	}

	cosiRuntime, err := runtime.NewRuntime(omniState, p.logger.With(zap.String("component", "cosi_runtime")), options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create runtime: %w", err)
	}

	return cosiRuntime, nil
}

type component struct {
	run  func(context.Context) error
	name string
}

func (p *Provider) runComponents(ctx context.Context, components []component) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	for _, comp := range components {
		logger := p.logger.With(zap.String("component", comp.name))

		eg.Go(func() error {
			defer cancel() // cancel the parent context, so all other components are also stopped

			logger.Info("start component")

			err := comp.run(ctx)
			if err != nil {
				logger.Error("failed to run component", zap.Error(err))

				return err
			}

			logger.Info("component stopped")

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

func (p *Provider) parseLabels() (map[string]string, error) {
	labels := make(map[string]string, len(p.options.MachineLabels))

	for _, l := range p.options.MachineLabels {
		parts := strings.Split(l, "=")
		if len(parts) > 2 {
			return nil, fmt.Errorf("malformed label %s", l)
		}

		value := ""

		if len(parts) > 1 {
			value = parts[1]
		}

		labels[parts[0]] = value
	}

	return labels, nil
}
