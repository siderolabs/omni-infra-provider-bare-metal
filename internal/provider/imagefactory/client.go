// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package imagefactory

import (
	"context"
	"fmt"
	"time"

	"github.com/siderolabs/image-factory/pkg/schematic"
	"github.com/siderolabs/omni/client/pkg/client"
	omnifactory "github.com/siderolabs/omni/client/pkg/imagefactory"
	"github.com/siderolabs/omni/client/pkg/infra"
	"github.com/siderolabs/omni/client/pkg/infra/provision"
	"github.com/siderolabs/talos/pkg/machinery/constants"
	"go.uber.org/zap"
)

var agentModeExtensions = []string{
	// include all firmware extensions
	"siderolabs/amd-ucode",
	"siderolabs/amdgpu-firmware",
	"siderolabs/bnx2-bnx2x",
	"siderolabs/chelsio-firmware",
	"siderolabs/i915-ucode",
	"siderolabs/intel-ice-firmware",
	"siderolabs/intel-ucode",
	"siderolabs/qlogic-firmware",
	"siderolabs/realtek-firmware",
	// include the agent extension itself
	"siderolabs/metal-agent",
}

// Client resolves a boot request into an iPXE URL through Omni.
//
// Omni picks the factory and authenticates the fetch, so no factory address or credentials live here.
type Client struct {
	omniClient            *client.Client
	logger                *zap.Logger
	agentModeTalosVersion string
	secureBootEnabled     bool
}

// NewClient creates a new image factory client.
func NewClient(omniClient *client.Client, agentModeTalosVersion string, secureBootEnabled bool, logger *zap.Logger) (*Client, error) {
	return &Client{
		omniClient:            omniClient,
		agentModeTalosVersion: agentModeTalosVersion,
		secureBootEnabled:     secureBootEnabled,
		logger:                logger,
	}, nil
}

// SchematicIPXEURL ensures a schematic exists on the image factory and returns the iPXE URL to it.
//
// If agentMode is true, the schematic will be created with the firmware extensions and the metal-agent extension.
func (c *Client) SchematicIPXEURL(ctx context.Context, agentMode bool, talosVersion, arch string, extensions, extraKernelArgs []string) (string, error) {
	logger := c.logger.With(zap.String("talos_version", talosVersion), zap.String("arch", arch),
		zap.Strings("extensions", extensions), zap.Strings("extra_kernel_args", extraKernelArgs))

	logger.Debug("generate schematic iPXE URL")

	var metaValues []schematic.MetaValue

	if !agentMode && talosVersion == "" {
		return "", fmt.Errorf("talosVersion is required when not booting into agent mode")
	}

	if agentMode {
		talosVersion = c.agentModeTalosVersion

		extensions = agentModeExtensions
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	sch := schematic.Schematic{
		Customization: schematic.Customization{
			ExtraKernelArgs: extraKernelArgs,
			Meta:            metaValues,
			SystemExtensions: schematic.SystemExtensions{
				OfficialExtensions: extensions,
			},
		},
	}

	marshaled, err := sch.Marshal()
	if err != nil {
		return "", fmt.Errorf("failed to marshal schematic: %w", err)
	}

	logger.Debug("generated schematic", zap.String("schematic", string(marshaled)))

	media, err := infra.EnsureInstallationMedia(ctx, c.omniClient, talosVersion, sch, provision.MediaSpec{
		MediaSpec: omnifactory.MediaSpec{
			Kind:         omnifactory.InstallationMediaKindPXE,
			Platform:     constants.PlatformMetal,
			Architecture: arch,
			SecureBoot:   c.secureBootEnabled,
		},
		StandaloneURL: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to resolve the iPXE installation media: %w", err)
	}

	logger.Debug("generated schematic iPXE URL",
		zap.String("schematic_id", media.SchematicID),
		zap.String("image_factory_host", media.ImageFactoryHost))

	return media.URL, nil
}
