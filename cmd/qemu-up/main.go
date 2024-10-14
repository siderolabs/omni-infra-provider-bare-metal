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

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/qemu"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/version"
)

var (
	qemuOptions qemu.Options
	destroy     bool
	debug       bool
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:     version.Name,
	Short:   "Bring QMU Talos machines up with iPXE boot for development and testing",
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
	machines, err := qemu.New(qemuOptions, logger)
	if err != nil {
		return fmt.Errorf("failed to create QEMU machines: %w", err)
	}

	if destroy {
		if err = machines.Destroy(ctx); err != nil {
			return fmt.Errorf("failed to destroy machines: %w", err)
		}

		return nil
	}

	return machines.Run(ctx)
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
	rootCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug mode & logs.")
	rootCmd.Flags().BoolVar(&destroy, "destroy", false, "Destroy existing machines and exit.")

	rootCmd.Flags().StringVar(&qemuOptions.Name, "name", qemu.DefaultOptions.Name, "Name of the cluster (the set of machines).")
	rootCmd.Flags().StringVar(&qemuOptions.CIDR, "cidr", qemu.DefaultOptions.CIDR, "CIDR for the machines' network.")
	rootCmd.Flags().StringVar(&qemuOptions.CNIBundleURL, "cni-bundle-url", qemu.DefaultOptions.CNIBundleURL, "URL to the CNI bundle.")
	rootCmd.Flags().StringVar(&qemuOptions.TalosctlPath, "talosctl-path", qemu.DefaultOptions.TalosctlPath,
		fmt.Sprintf("Path to the talosctl binary. If not specified, the binary %q will be looked up in the current working dir and in the PATH.", qemu.TalosctlBinary))
	rootCmd.Flags().StringVar(&qemuOptions.CPUs, "cpus", qemu.DefaultOptions.CPUs, "Number of CPUs for each machine.")
	rootCmd.Flags().IntVar(&qemuOptions.NumMachines, "num-machines", qemu.DefaultOptions.NumMachines, "Number of machines to bring up.")
	rootCmd.Flags().IntVar(&qemuOptions.MTU, "mtu", qemu.DefaultOptions.MTU, "MTU for the machines' network.")
	rootCmd.Flags().Uint64Var(&qemuOptions.DiskSize, "disk-size", qemu.DefaultOptions.DiskSize, "Disk size for each machine.")
	rootCmd.Flags().Int64Var(&qemuOptions.MemSize, "mem-size", qemu.DefaultOptions.MemSize, "Memory size for each machine.")
	rootCmd.Flags().StringSliceVar(&qemuOptions.Nameservers, "nameservers", qemu.DefaultOptions.Nameservers, "Nameservers for the machines' network.")
}
