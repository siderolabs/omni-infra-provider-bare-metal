// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package main implements the main entrypoint for the Omni bare metal infra provider.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/siderolabs/talos-metal-agent/pkg/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/constants"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/pxe"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/ipxe"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/version"
)

var (
	providerOptions = provider.DefaultOptions()
	debug           bool
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:     version.Name,
	Short:   "Run the Omni bare metal infra provider",
	Version: version.Tag,
	Args:    cobra.NoArgs,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		cmd.SilenceUsage = true // if the args are parsed fine, no need to show usage
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		logger, err := initLogger()
		if err != nil {
			return fmt.Errorf("failed to create logger: %w", err)
		}

		defer logger.Sync() //nolint:errcheck

		return run(cmd.Context(), logger)
	},
}

func initLogger() (*zap.Logger, error) {
	var loggerConfig zap.Config

	if debug {
		loggerConfig = zap.NewDevelopmentConfig()
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		loggerConfig.Level.SetLevel(zap.DebugLevel)
	} else {
		loggerConfig = zap.NewProductionConfig()
		loggerConfig.Level.SetLevel(zap.InfoLevel)
	}

	return loggerConfig.Build(zap.AddStacktrace(zapcore.FatalLevel)) // only print stack traces for fatal errors)
}

func run(ctx context.Context, logger *zap.Logger) error {
	prov := provider.New(providerOptions, logger)

	if err := prov.Run(ctx); err != nil {
		return fmt.Errorf("failed to run provider: %w", err)
	}

	return nil
}

func main() {
	if err := runCmd(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}

func runCmd() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	return rootCmd.ExecuteContext(ctx)
}

func init() {
	const (
		apiPowerMgmtStateDirFlag = "api-power-mgmt-state-dir"
		ipmiPXEBootModeFlag      = "ipmi-pxe-boot-mode"
		useLocalBootAssetsFlag   = "use-local-boot-assets"
		apiAdvertiseAddressFlag  = "api-advertise-address"
		caCertFileFlag           = "ca-cert-file"
		certFileFlag             = "cert-file"
		keyFileFlag              = "key-file"
	)

	rootCmd.Flags().Var(&meta.ProviderID, "id", "The id of the infra provider, it is used to match the resources with the infra provider label.")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug mode & logs.")

	rootCmd.Flags().StringVar(&providerOptions.APIListenAddress, "api-listen-address", providerOptions.APIListenAddress,
		"The IP address to listen on. If not specified, the server will listen on all interfaces.")
	rootCmd.Flags().StringVar(&providerOptions.APIAdvertiseAddress, apiAdvertiseAddressFlag, providerOptions.APIAdvertiseAddress,
		"The IP address to advertise. Required if the server has more than a single routable IP address. If not specified, the single routable IP address will be used.")
	rootCmd.Flags().IntVar(&providerOptions.APIPort, "api-port", providerOptions.APIPort, "The port to run the api server on.")
	rootCmd.Flags().StringVar(&providerOptions.OmniAPIEndpoint, "omni-api-endpoint", os.Getenv("OMNI_ENDPOINT"),
		"The endpoint of the Omni API, if not set, defaults to OMNI_ENDPOINT env var.")
	rootCmd.Flags().StringVar(&providerOptions.Name, "provider-name", providerOptions.Name, "Provider name as it appears in Omni")
	rootCmd.Flags().StringVar(&providerOptions.Description, "provider-description", providerOptions.Description, "Provider description as it appears in Omni")
	rootCmd.Flags().BoolVar(&providerOptions.UseLocalBootAssets, useLocalBootAssetsFlag, providerOptions.UseLocalBootAssets,
		"Use local boot assets for iPXE booting. If set, the iPXE server will use the kernel and initramfs from the local assets "+
			"instead of forwarding the request to the image factory to boot into agent mode.")
	rootCmd.Flags().StringVar(&providerOptions.DHCPProxyIfaceOrIP, "dhcp-proxy-iface-or-ip", providerOptions.DHCPProxyIfaceOrIP,
		"The interface name or the IP address on the interface to run the DHCP proxy server on. "+
			"If it is an IP address, the DHCP proxy server will run on the interface that has the IP address. "+
			"If not specified, defaults to the API advertise address.")
	rootCmd.Flags().StringVar(&providerOptions.ImageFactoryBaseURL, "image-factory-base-url", providerOptions.ImageFactoryBaseURL,
		"The base URL of the image factory.")
	rootCmd.Flags().StringVar(&providerOptions.ImageFactoryPXEBaseURL, "image-factory-pxe-base-url", providerOptions.ImageFactoryPXEBaseURL,
		"The base URL of the image factory PXE server.")
	rootCmd.Flags().StringVar(&providerOptions.AgentModeTalosVersion, "agent-mode-talos-version", providerOptions.AgentModeTalosVersion,
		"The default Talos version to when forwarding iPXE requests to the image factory to boot into Talos agent.")
	rootCmd.Flags().BoolVar(&providerOptions.AgentTestMode, "agent-test-mode", providerOptions.AgentTestMode,
		fmt.Sprintf("Enable agent test mode. In this mode, the Talos agent will be booted into the test mode via the kernel arg %q. "+
			`In this mode, you probably want to set the "--%s" flag, as the test mode agents are probably QEMU machines whose power is managed over the HTTP API.`,
			config.TestModeKernelArg, apiPowerMgmtStateDirFlag))
	rootCmd.Flags().StringVar(&providerOptions.APIPowerMgmtStateDir, apiPowerMgmtStateDirFlag, providerOptions.APIPowerMgmtStateDir,
		"The directory to read the power management API endpoints and ports, to be used to manage the power state of the machines which are managed via API "+
			"(e.g., QEMU VMs created by 'qemu-up' or 'talosctl cluster create') Mainly used for testing purposes.")
	rootCmd.Flags().StringVar(&providerOptions.BootFromDiskMethod, "boot-from-disk-method", providerOptions.BootFromDiskMethod,
		fmt.Sprintf("Default method to use to boot server from disk if it hits iPXE endpoint after install. Valid values are: %v",
			[]ipxe.BootFromDiskMethod{ipxe.BootIPXEExit, ipxe.Boot404, ipxe.BootSANDisk}))
	rootCmd.Flags().StringVar(&providerOptions.IPMIPXEBootMode, ipmiPXEBootModeFlag, providerOptions.IPMIPXEBootMode,
		fmt.Sprintf("Default boot mode to use when PXE booting a machine via IPMI. Valid values are: %v",
			[]pxe.BootMode{pxe.BootModeBIOS, pxe.BootModeUEFI}))
	rootCmd.Flags().StringSliceVar(&providerOptions.MachineLabels, "machine-labels", providerOptions.MachineLabels,
		"Comma separated list of key=value pairs to be set to the machine. Example: key1=value1,key2,key3=value3")
	rootCmd.Flags().BoolVar(&providerOptions.InsecureSkipTLSVerify, "insecure-skip-tls-verify", providerOptions.InsecureSkipTLSVerify,
		"Skip TLS verification when connecting to the Omni API.")
	rootCmd.Flags().DurationVar(&providerOptions.MinRebootInterval, "min-reboot-interval", providerOptions.MinRebootInterval,
		"the minimum interval between reboots of the machine issued by the provider. This is to prevent the provider from issuing reboots too frequently.")
	rootCmd.Flags().BoolVar(&providerOptions.SecureBootEnabled, "secure-boot-enabled", providerOptions.SecureBootEnabled,
		fmt.Sprintf("Serve secure boot UKI from the iPXE endpoint. The UKI can be used to boot a machine without secure boot, "+
			`but it is required to boot a machine with secure boot. When enabled, "--%s" must be set to %q and "--%s" must be set to false.`,
			ipmiPXEBootModeFlag, pxe.BootModeUEFI, useLocalBootAssetsFlag),
	)

	if constants.IsDebugBuild {
		rootCmd.Flags().BoolVar(&providerOptions.ClearState, "clear-state", providerOptions.ClearState, "Clear the state of the provider on startup.")
	}

	rootCmd.Flags().BoolVar(&providerOptions.EnableResourceCache, "enable-resource-cache", providerOptions.EnableResourceCache,
		"Enable controller runtime resource cache.")
	rootCmd.Flags().BoolVar(&providerOptions.DisableDHCPProxy, "disable-dhcp-proxy", providerOptions.DisableDHCPProxy,
		"Disable the DHCP proxy server.")

	// TLS options
	rootCmd.Flags().BoolVar(&providerOptions.TLS.Enabled, "tls-enabled", providerOptions.TLS.Enabled,
		"Enable TLS for the API server.")
	rootCmd.Flags().IntVar(&providerOptions.TLS.APIPort, "tls-api-port", providerOptions.TLS.APIPort,
		"The port to run the API server on when using TLS.")
	rootCmd.Flags().BoolVar(&providerOptions.TLS.AgentSkipVerify, "tls-agent-skip-verify", providerOptions.TLS.AgentSkipVerify,
		"Make the Talos agent GRPC client skip TLS verification when connecting to the provider.")
	rootCmd.Flags().DurationVar(&providerOptions.TLS.CATTL, "tls-ca-ttl", providerOptions.TLS.CATTL,
		fmt.Sprintf("CA certificate TTL. Ignored if --%s, --%s or --%s is set.", caCertFileFlag, certFileFlag, keyFileFlag))
	rootCmd.Flags().DurationVar(&providerOptions.TLS.CertTTL, "tls-cert-ttl", providerOptions.TLS.CertTTL,
		fmt.Sprintf("TTL for the generated ephemeral certificates using the CA certificate. Ignored if --%s, --%s or --%s is set.", caCertFileFlag, certFileFlag, keyFileFlag))

	rootCmd.Flags().StringVar(&providerOptions.TLS.CACertFile, "tls-ca-cert-file", providerOptions.TLS.CACertFile,
		fmt.Sprintf("The CA certificate file. "+
			"If specified, the provider will not generate a CA certificate to issue ephemeral, short-lived TLS certificates, and this will be used instead."+
			"When set, --%s and --%s must also be set.", certFileFlag, keyFileFlag))
	rootCmd.Flags().StringVar(&providerOptions.TLS.CertFile, "tls-cert-file", providerOptions.TLS.CertFile,
		fmt.Sprintf("The TLS certificate file. It MUST to be valid for the host specified in --%s via SAN DNS names or IP addresses. "+
			"If specified, the provider will not generate a CA certificate to issue ephemeral, short-lived TLS certificates, and this will be used instead. "+
			"When set, --%s is also required. Required if --%s is set.", apiAdvertiseAddressFlag, keyFileFlag, caCertFileFlag))
	rootCmd.Flags().StringVar(&providerOptions.TLS.KeyFile, "tls-key-file", providerOptions.TLS.KeyFile,
		fmt.Sprintf("The TLS key file. "+
			"If specified, the provider will not generate a CA certificate to issue ephemeral, short-lived TLS certificates, and this will be used instead. "+
			"When set, --%s is also required. Required if --%s is set.", certFileFlag, caCertFileFlag))

	// RedFish options
	rootCmd.Flags().BoolVar(&providerOptions.Redfish.UseAlways, "redfish-use-always", providerOptions.Redfish.UseAlways,
		"Always use Redfish for power management.")
	rootCmd.Flags().BoolVar(&providerOptions.Redfish.UseWhenAvailable, "redfish-use-when-available", providerOptions.Redfish.UseWhenAvailable,
		"Use Redfish for power management when available.")
	rootCmd.Flags().BoolVar(&providerOptions.Redfish.UseHTTPS, "redfish-use-https", providerOptions.Redfish.UseHTTPS,
		"Use HTTPS for Redfish connections.")
	rootCmd.Flags().BoolVar(&providerOptions.Redfish.InsecureSkipTLSVerify, "redfish-insecure-skip-tls-verify", providerOptions.Redfish.InsecureSkipTLSVerify,
		"Skip TLS verification when connecting to Redfish.")
	rootCmd.Flags().IntVar(&providerOptions.Redfish.Port, "redfish-port", providerOptions.Redfish.Port,
		"The port to connect to Redfish.")
	rootCmd.Flags().BoolVar(&providerOptions.Redfish.SetBootSourceOverrideMode, "redfish-set-boot-source-override-mode", providerOptions.Redfish.SetBootSourceOverrideMode,
		"Set the boot source override mode field when using Redfish for power management. Some Redfish implementations require this field to be unset.")

	// Agent Client options
	rootCmd.Flags().BoolVar(&providerOptions.AgentClient.WipeWithZeroes, "wipe-with-zeroes", providerOptions.AgentClient.WipeWithZeroes,
		"When wiping a machine, write zeroes to the whole disk instead doing a fast wipe.")
	rootCmd.Flags().DurationVar(&providerOptions.AgentClient.CallTimeout, "agent-call-timeout", providerOptions.AgentClient.CallTimeout,
		"Timeout for agent calls.")
	rootCmd.Flags().DurationVar(&providerOptions.AgentClient.FastWipeTimeout, "agent-fast-wipe-timeout", providerOptions.AgentClient.FastWipeTimeout,
		"Timeout for fast wipe operation (without zeroes) call to the agent.")
	rootCmd.Flags().DurationVar(&providerOptions.AgentClient.ZeroesWipeTimeout, "agent-zeroes-wipe-timeout", providerOptions.AgentClient.ZeroesWipeTimeout,
		"Timeout for slow wipe operation (with zeroes) call to the agent.")
}
