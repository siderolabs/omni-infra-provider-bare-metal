// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ipxe provides iPXE functionality.
package ipxe

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"github.com/siderolabs/talos-metal-agent/pkg/config"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/baremetal"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/boot"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/machinestatus"
)

const (
	ipxeScriptTemplateFormat = `#!ipxe
chain --replace %s
`

	ipxeScriptTemplateFormatLocalAssets = `#!ipxe
kernel %s %s
initrd %s
boot
`

	// ipxeBootFromDiskExit script is used to skip PXE booting and boot from disk via exit.
	ipxeBootFromDiskExit = `#!ipxe
exit
`

	// ipxeBootFromDiskSanboot script is used to skip PXE booting and boot from disk via sanboot.
	ipxeBootFromDiskSanboot = `#!ipxe
sanboot --no-describe --drive 0x80
`

	archAmd64 = "amd64"
	archArm64 = "arm64"
)

// ImageFactoryClient represents an image factory client which ensures a schematic exists on image factory, and returns the PXE URL to it.
type ImageFactoryClient interface {
	SchematicIPXEURL(ctx context.Context, agentMode bool, talosVersion, arch string, extensions, extraKernelArgs []string) (string, error)
}

// HandlerOptions represents the options for the iPXE handler.
type HandlerOptions struct {
	APIAdvertiseAddress string
	BootFromDiskMethod  string
	APIPort             int
	UseLocalBootAssets  bool
	AgentTestMode       bool
}

// Handler represents an iPXE handler.
type Handler struct {
	imageFactoryClient ImageFactoryClient
	state              state.State

	logger *zap.Logger

	bootFromDiskMethod BootFromDiskMethod

	defaultKernelArgs []string
	agentKernelArgs   []string

	options HandlerOptions
}

// ServeHTTP serves the iPXE request.
//
// URL pattern: http://ip-of-this-provider:50042/ipxe?uuid=${uuid}&mac=${net${idx}/mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}&arch=${buildarch}
//
// Implements http.Handler interface.
func (handler *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	query := req.URL.Query()
	uuid := query.Get("uuid")
	mac := query.Get("mac")
	arch := query.Get("arch")
	logger := handler.logger.With(zap.String("uuid", uuid), zap.String("mac", mac), zap.String("arch", arch))

	logger.Info("handle iPXE request")

	decision, err := handler.makeBootDecision(ctx, arch, uuid, logger)
	if err != nil {
		handler.logger.Error("failed to check if Talos is installed", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)

		if _, err = w.Write([]byte("failed to check if Talos is installed")); err != nil {
			handler.logger.Error("failed to write error response", zap.Error(err))
		}

		return
	}

	w.WriteHeader(decision.statusCode)

	if _, err = w.Write([]byte(decision.body)); err != nil {
		handler.logger.Error("failed to write ok response", zap.Error(err))

		return
	}

	// save the machine in any case
	if _, err = machinestatus.Modify(ctx, handler.state, uuid, func(status *baremetal.MachineStatus) error {
		status.TypedSpec().Value.BootMode = decision.mode

		return nil
	}); err != nil {
		handler.logger.Error("failed to mark machine as dirty", zap.Error(err))
	}
}

type bootDecision struct {
	body       string
	statusCode int
	mode       specs.BootMode
}

func (handler *Handler) makeBootDecision(ctx context.Context, arch, uuid string, logger *zap.Logger) (bootDecision, error) {
	switch arch { // https://ipxe.org/cfg/buildarch
	case archArm64:
		arch = archArm64
	default:
		arch = archAmd64
	}

	machineResources, err := handler.getResources(ctx, uuid)
	if err != nil {
		return bootDecision{statusCode: http.StatusInternalServerError}, fmt.Errorf("failed to get resources: %w", err)
	}

	var userExtraKernelArgs []string

	if machineResources.infraMachine != nil {
		userExtraKernelArgs = strings.Fields(machineResources.infraMachine.TypedSpec().Value.ExtraKernelArgs)
	}

	mode, err := boot.DetermineRequiredMode(machineResources.infraMachine, machineResources.status, machineResources.machineState, logger)
	if err != nil {
		return bootDecision{statusCode: http.StatusInternalServerError}, fmt.Errorf("failed to determine required boot mode: %w", err)
	}

	requiredBootMode := mode.BootMode

	logger = logger.With(zap.Stringer("required_boot_mode", requiredBootMode))

	switch requiredBootMode {
	case specs.BootMode_BOOT_MODE_AGENT_PXE:
		body, statusCode, agentErr := handler.bootIntoAgentMode(ctx, arch, userExtraKernelArgs)
		if agentErr != nil {
			return bootDecision{statusCode: http.StatusInternalServerError}, fmt.Errorf("failed to boot into agent mode: %w", agentErr)
		}

		return bootDecision{
			mode:       requiredBootMode,
			body:       body,
			statusCode: statusCode,
		}, nil
	case specs.BootMode_BOOT_MODE_TALOS_PXE:
		logger.Info("boot Talos over iPXE")

		consoleKernelArgs := handler.consoleKernelArgs(arch)
		extraKernelArgs := slices.Concat(handler.defaultKernelArgs, consoleKernelArgs, userExtraKernelArgs)
		talosVersion := machineResources.infraMachine.TypedSpec().Value.ClusterTalosVersion
		extensions := machineResources.infraMachine.TypedSpec().Value.Extensions

		var ipxeURL string

		ipxeURL, err = handler.imageFactoryClient.SchematicIPXEURL(ctx, false, talosVersion, arch, extensions, extraKernelArgs)
		if err != nil {
			return bootDecision{statusCode: http.StatusInternalServerError}, fmt.Errorf("failed to get schematic IPXE URL: %w", err)
		}

		ipxeScript := fmt.Sprintf(ipxeScriptTemplateFormat, ipxeURL)

		return bootDecision{
			mode:       requiredBootMode,
			body:       ipxeScript,
			statusCode: http.StatusOK,
		}, nil
	case specs.BootMode_BOOT_MODE_TALOS_DISK:
		logger.Info("boot from disk")

		switch handler.bootFromDiskMethod {
		case Boot404:
			return bootDecision{statusCode: http.StatusNotFound}, nil
		case BootSANDisk:
			return bootDecision{body: ipxeBootFromDiskSanboot, statusCode: http.StatusOK}, nil
		case BootIPXEExit:
			fallthrough
		default:
			return bootDecision{
				mode:       requiredBootMode,
				body:       ipxeBootFromDiskExit,
				statusCode: http.StatusOK,
			}, nil
		}
	case specs.BootMode_BOOT_MODE_UNKNOWN:
		fallthrough
	default:
		return bootDecision{statusCode: http.StatusInternalServerError}, fmt.Errorf("unknown boot mode: %s", requiredBootMode)
	}
}

type resources struct {
	infraMachine *infra.Machine
	status       *baremetal.MachineStatus
	machineState *infra.MachineState
}

func (handler *Handler) getResources(ctx context.Context, id string) (resources, error) {
	infraMachine, err := safe.StateGetByID[*infra.Machine](ctx, handler.state, id)
	if err != nil && !state.IsNotFoundError(err) {
		return resources{}, fmt.Errorf("failed to get infra machine: %w", err)
	}

	status, err := machinestatus.Modify(ctx, handler.state, id, nil)
	if err != nil {
		return resources{}, fmt.Errorf("failed to get bare metal machine status: %w", err)
	}

	machineState, err := safe.StateGetByID[*infra.MachineState](ctx, handler.state, id)
	if err != nil && !state.IsNotFoundError(err) {
		return resources{}, fmt.Errorf("failed to get infra machine install status: %w", err)
	}

	return resources{
		infraMachine: infraMachine,
		status:       status,
		machineState: machineState,
	}, nil
}

func (handler *Handler) bootIntoAgentMode(ctx context.Context, arch string, extraKernelArgs []string) (string, int, error) {
	agentExtraKernelArgs := slices.Concat(handler.agentKernelArgs, handler.consoleKernelArgs(arch), extraKernelArgs)

	if handler.options.UseLocalBootAssets {
		handler.logger.Info("boot agent mode using local assets")

		return handler.bootAgentModeViaLocalIPXEScript(arch, agentExtraKernelArgs)
	}

	handler.logger.Info("boot agent mode using image factory")

	return handler.bootViaFactoryIPXEScript(ctx, true, arch, agentExtraKernelArgs)
}

func (handler *Handler) bootViaFactoryIPXEScript(ctx context.Context, agentMode bool, arch string, kernelArgs []string) (string, int, error) {
	ipxeURL, err := handler.imageFactoryClient.SchematicIPXEURL(ctx, agentMode, "", arch, nil, kernelArgs)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("failed to get schematic IPXE URL: %w", err)
	}

	ipxeScript := fmt.Sprintf(ipxeScriptTemplateFormat, ipxeURL)

	return ipxeScript, http.StatusOK, nil
}

func (handler *Handler) bootAgentModeViaLocalIPXEScript(arch string, kernelArgs []string) (string, int, error) {
	hostPort := net.JoinHostPort(handler.options.APIAdvertiseAddress, strconv.Itoa(handler.options.APIPort))

	kernel := fmt.Sprintf("http://%s/assets/kernel-%s", hostPort, arch)
	initramfs := fmt.Sprintf("http://%s/assets/initramfs-metal-%s.xz", hostPort, arch)

	// read cmdline
	cmdline, err := os.ReadFile(fmt.Sprintf("/assets/cmdline-metal-%s", arch))
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("failed to read cmdline: %w", err)
	}

	cmdlineStr := strings.TrimSpace(string(cmdline)) + " " + strings.Join(kernelArgs, " ")
	ipxeScript := fmt.Sprintf(ipxeScriptTemplateFormatLocalAssets, kernel, cmdlineStr, initramfs)

	return ipxeScript, http.StatusOK, nil
}

func (handler *Handler) consoleKernelArgs(arch string) []string {
	switch arch {
	case archArm64:
		return []string{"console=tty0", "console=ttyAMA0"}
	default:
		return []string{"console=tty0", "console=ttyS0"}
	}
}

// NewHandler creates a new iPXE server.
func NewHandler(imageFactoryClient ImageFactoryClient, state state.State, options HandlerOptions, logger *zap.Logger) (*Handler, error) {
	bootFromDiskMethod, err := parseBootFromDiskMethod(options.BootFromDiskMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to parse boot from disk method: %w", err)
	}

	logger.Info("patch iPXE binaries")

	if err := patchBinaries(options.APIAdvertiseAddress, options.APIPort); err != nil {
		return nil, err
	}

	logger.Info("successfully patched iPXE binaries")

	apiHostPort := net.JoinHostPort(options.APIAdvertiseAddress, strconv.Itoa(options.APIPort))

	talosConfigURL := fmt.Sprintf("http://%s/config?u=${uuid}", apiHostPort)
	defaultKernelArgs := []string{
		"talos.config=" + talosConfigURL,
	}

	agentExtraKernelArgs := []string{
		fmt.Sprintf("%s=%s", config.MetalProviderAddressKernelArg, apiHostPort),
	}

	if options.AgentTestMode {
		agentExtraKernelArgs = append(agentExtraKernelArgs,
			fmt.Sprintf("%s=%s", config.TestModeKernelArg, "1"),
		)
	}

	agentKernelArgs := slices.Concat(defaultKernelArgs, agentExtraKernelArgs)

	return &Handler{
		state:              state,
		imageFactoryClient: imageFactoryClient,
		options:            options,
		defaultKernelArgs:  defaultKernelArgs,
		agentKernelArgs:    agentKernelArgs,
		bootFromDiskMethod: bootFromDiskMethod,
		logger:             logger,
	}, nil
}
