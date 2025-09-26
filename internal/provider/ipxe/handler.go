// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/klauspost/compress/zstd"
	omnispecs "github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	agentconfig "github.com/siderolabs/talos-metal-agent/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/constants"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/controllers"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/machine"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/tls"
)

const (
	ipxeScriptTemplateFormat = `#!ipxe
chain --replace %s
`

	ipxeScriptTemplateFormatLocalAssets = `#!ipxe
imgfree
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

	// initScriptName is the name of the iPXE init script served by the HTTP server.
	//
	// Some UEFIs with built-in iPXE require the script URL to be in the form of a filename ending with ".ipxe", hence we serve it under this path.
	initScriptName = "init.ipxe"

	// bootScriptName is the name of the iPXE boot script served by the HTTP server.
	//
	// Some UEFIs with built-in iPXE require the script URL to be in the form of a filename ending with ".ipxe", hence we serve it under this path.
	bootScriptName = "boot.ipxe"
)

// ImageFactoryClient represents an image factory client which ensures a schematic exists on image factory, and returns the PXE URL to it.
type ImageFactoryClient interface {
	SchematicIPXEURL(ctx context.Context, agentMode bool, talosVersion, arch string, extensions, extraKernelArgs []string) (string, error)
}

// HandlerOptions represents the options for the iPXE handler.
type HandlerOptions struct {
	APIAdvertiseAddress string
	BootFromDiskMethod  string
	TLS                 tls.Options
	APIPort             int
	UseLocalBootAssets  bool
	AgentTestMode       bool
}

// Handler represents an iPXE handler.
type Handler struct {
	imageFactoryClient ImageFactoryClient
	reader             controller.Reader
	logger             *zap.Logger
	pxeBootEventCh     chan<- controllers.PXEBootEvent
	bootFromDiskMethod BootFromDiskMethod
	defaultKernelArgs  []string
	agentKernelArgs    []string
	initScript         []byte
	options            HandlerOptions
}

// ServeHTTP serves the iPXE request.
//
// URL pattern: http://ip-of-this-provider:50042/ipxe/boot.ipxe?uuid=${uuid}&mac=${net${idx}/mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}&arch=${buildarch}
//
// Implements http.Handler interface.
func (handler *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.PathValue("script") {
	default:
		handler.logger.Error("invalid iPXE script", zap.String("script", req.PathValue("script")))

		w.WriteHeader(http.StatusNotFound)

		return
	case initScriptName:
		handler.handleInitScript(w)

		return
	case bootScriptName:
	}

	ctx := req.Context()
	query := req.URL.Query()
	uuid := query.Get("uuid")
	mac := query.Get("mac")
	arch := query.Get("arch")
	logger := handler.logger.With(zap.String("uuid", uuid), zap.String("mac", mac), zap.String("arch", arch))

	logger.Info("handle iPXE boot request")

	decision, err := handler.makeBootDecision(ctx, arch, uuid, logger)
	if err != nil {
		handler.logger.Error("failed to make boot decision", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)

		if _, err = w.Write([]byte("failed to make boot decision: " + err.Error())); err != nil {
			handler.logger.Error("failed to write error response", zap.Error(err))
		}

		return
	}

	w.WriteHeader(decision.statusCode)

	if _, err = w.Write([]byte(decision.body)); err != nil {
		handler.logger.Error("failed to write response", zap.Error(err))

		return
	}

	select {
	case <-ctx.Done():
		handler.logger.Error("failed to send PXE boot event", zap.String("uuid", uuid))
	case handler.pxeBootEventCh <- controllers.PXEBootEvent{MachineID: uuid}:
	}
}

func (handler *Handler) handleInitScript(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")

	if _, err := w.Write(handler.initScript); err != nil {
		handler.logger.Error("failed to write init script", zap.Error(err))
	}
}

type bootDecision struct {
	body       string
	statusCode int
}

//nolint:gocyclo,cyclop
func (handler *Handler) makeBootDecision(ctx context.Context, arch, uuid string, logger *zap.Logger) (bootDecision, error) {
	if uuid == "" {
		return bootDecision{statusCode: http.StatusBadRequest}, fmt.Errorf("missing uuid")
	}

	switch arch { // https://ipxe.org/cfg/buildarch
	case archArm64:
		arch = archArm64
	default:
		arch = archAmd64
	}

	infraMachine, err := safe.ReaderGetByID[*infra.Machine](ctx, handler.reader, uuid)
	if err != nil && !state.IsNotFoundError(err) {
		return bootDecision{}, err
	}

	bmcConfiguration, err := safe.ReaderGetByID[*resources.BMCConfiguration](ctx, handler.reader, uuid)
	if err != nil && !state.IsNotFoundError(err) {
		return bootDecision{}, err
	}

	wipeStatus, err := safe.ReaderGetByID[*resources.WipeStatus](ctx, handler.reader, uuid)
	if err != nil && !state.IsNotFoundError(err) {
		return bootDecision{}, err
	}

	if infraMachine != nil {
		logger = logger.With(zap.String("machine_id", infraMachine.Metadata().ID()))

		if infraMachine.TypedSpec().Value.Cordoned {
			logger.Info("machine is cordoned, skip making a boot decision")

			return bootDecision{body: "machine is cordoned", statusCode: http.StatusNotFound}, nil
		}

		if infraMachine.TypedSpec().Value.AcceptanceStatus == omnispecs.InfraMachineConfigSpec_REJECTED {
			logger.Info("machine is rejected, return not found")

			return bootDecision{statusCode: http.StatusNotFound}, nil
		}
	}

	requiredBootMode := machine.RequiredBootMode(infraMachine, bmcConfiguration, wipeStatus, logger)

	var userExtraKernelArgs []string

	if infraMachine != nil {
		userExtraKernelArgs = strings.Fields(infraMachine.TypedSpec().Value.ExtraKernelArgs)

		if infraMachine.TypedSpec().Value.Cordoned {
			logger.Info("machine is cordoned, skip making a boot decision")

			return bootDecision{body: "machine is cordoned", statusCode: http.StatusNotFound}, nil
		}
	}

	switch requiredBootMode {
	case machine.BootModeAgentPXE:
		logger.Info("boot machine: Talos agent mode")

		body, statusCode, agentErr := handler.bootIntoAgentMode(ctx, arch, userExtraKernelArgs)
		if agentErr != nil {
			return bootDecision{statusCode: http.StatusInternalServerError}, fmt.Errorf("failed to boot into agent mode: %w", agentErr)
		}

		return bootDecision{
			body:       body,
			statusCode: statusCode,
		}, nil
	case machine.BootModeTalosPXE:
		logger.Info("boot machine: Talos over iPXE")

		consoleKernelArgs := handler.consoleKernelArgs(arch)
		extraKernelArgs := slices.Concat(handler.defaultKernelArgs, consoleKernelArgs, userExtraKernelArgs)
		talosVersion := infraMachine.TypedSpec().Value.ClusterTalosVersion
		extensions := infraMachine.TypedSpec().Value.Extensions

		var ipxeURL string

		ipxeURL, err = handler.imageFactoryClient.SchematicIPXEURL(ctx, false, talosVersion, arch, extensions, extraKernelArgs)
		if err != nil {
			return bootDecision{statusCode: http.StatusInternalServerError}, fmt.Errorf("failed to get schematic IPXE URL: %w", err)
		}

		ipxeScript := fmt.Sprintf(ipxeScriptTemplateFormat, ipxeURL)

		return bootDecision{
			body:       ipxeScript,
			statusCode: http.StatusOK,
		}, nil
	case machine.BootModeTalosDisk:
		logger.Info("boot machine: from the disk")

		switch handler.bootFromDiskMethod {
		case Boot404:
			return bootDecision{statusCode: http.StatusNotFound}, nil
		case BootSANDisk:
			return bootDecision{body: ipxeBootFromDiskSanboot, statusCode: http.StatusOK}, nil
		case BootIPXEExit:
			fallthrough
		default:
			return bootDecision{
				body:       ipxeBootFromDiskExit,
				statusCode: http.StatusOK,
			}, nil
		}
	default:
		return bootDecision{statusCode: http.StatusInternalServerError}, fmt.Errorf("unknown boot mode: %s", requiredBootMode)
	}
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
func NewHandler(ctx context.Context, imageFactoryClient ImageFactoryClient, machineConfig []byte, r controller.Reader,
	pxeBootEventCh chan<- controllers.PXEBootEvent, options HandlerOptions, logger *zap.Logger,
) (*Handler, error) {
	bootFromDiskMethod, err := parseBootFromDiskMethod(options.BootFromDiskMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to parse boot from disk method: %w", err)
	}

	initScript, err := buildInitScript(options.APIAdvertiseAddress, options.APIPort)
	if err != nil {
		return nil, fmt.Errorf("failed to build init script: %w", err)
	}

	logger.Info("patch iPXE binaries")

	if err = patchBinaries(ctx, initScript, options.TLS.CustomIPXECACertFile, logger); err != nil {
		return nil, err
	}

	logger.Info("successfully patched iPXE binaries")

	apiHostPort := net.JoinHostPort(options.APIAdvertiseAddress, strconv.Itoa(options.APIPort))
	talosConfigURL := fmt.Sprintf("http://%s/config?u=${uuid}", apiHostPort)
	defaultKernelArgs := []string{
		fmt.Sprintf("%s=%s", constants.KernelParamConfig, talosConfigURL),
	}

	var providerAddress string

	if options.TLS.Enabled {
		providerAddress = "https://" + net.JoinHostPort(options.APIAdvertiseAddress, strconv.Itoa(options.TLS.APIPort))

		var inlineConfig string

		inlineConfig, err = compressInlineConfig(machineConfig, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to build inline config: %w", err)
		}

		logger.Debug("built inline config", zap.String("config", inlineConfig), zap.Int("length", len(inlineConfig)))

		defaultKernelArgs = append(defaultKernelArgs, fmt.Sprintf("%s=%s", constants.KernelParamConfigInline, inlineConfig))
	} else {
		providerAddress = net.JoinHostPort(options.APIAdvertiseAddress, strconv.Itoa(options.APIPort))
	}

	agentExtraKernelArgs := []string{
		fmt.Sprintf("%s=%s", agentconfig.MetalProviderAddressKernelArg, providerAddress),
	}

	if options.AgentTestMode {
		agentExtraKernelArgs = append(agentExtraKernelArgs,
			fmt.Sprintf("%s=%s", agentconfig.TestModeKernelArg, "1"),
		)
	}

	if options.TLS.AgentSkipVerify {
		agentExtraKernelArgs = append(agentExtraKernelArgs,
			fmt.Sprintf("%s=%s", agentconfig.TLSSkipVerifyKernelArg, "1"),
		)
	}

	agentKernelArgs := slices.Concat(defaultKernelArgs, agentExtraKernelArgs)

	return &Handler{
		pxeBootEventCh:     pxeBootEventCh,
		reader:             r,
		imageFactoryClient: imageFactoryClient,
		options:            options,
		defaultKernelArgs:  defaultKernelArgs,
		agentKernelArgs:    agentKernelArgs,
		bootFromDiskMethod: bootFromDiskMethod,
		initScript:         initScript,
		logger:             logger,
	}, nil
}

func compressInlineConfig(config []byte, logger *zap.Logger) (string, error) {
	var buf bytes.Buffer

	zstdEncoder, err := zstd.NewWriter(&buf, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		return "", fmt.Errorf("failed to create zstd encoder: %w", err)
	}

	closeFunc := sync.OnceValue(zstdEncoder.Close)

	defer func() {
		if closeErr := closeFunc(); closeErr != nil {
			logger.Error("failed to close zstd encoder", zap.Error(closeErr))
		}
	}()

	if _, err = zstdEncoder.Write(config); err != nil {
		return "", fmt.Errorf("failed to write zstd data: %w", err)
	}

	if err = closeFunc(); err != nil {
		return "", fmt.Errorf("failed to close zstd encoder: %w", err)
	}

	configBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return configBase64, nil
}
