// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package imagefactory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/siderolabs/image-factory/pkg/schematic"
	"github.com/siderolabs/omni/client/pkg/client"
	omniconstants "github.com/siderolabs/omni/client/pkg/constants"
	omnifactory "github.com/siderolabs/omni/client/pkg/imagefactory"
	"github.com/siderolabs/omni/client/pkg/infra"
	"github.com/siderolabs/omni/client/pkg/infra/provision"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
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
	resolvedAt            time.Time
	omniClient            *client.Client
	logger                *zap.Logger
	agentModeTalosVersion string
	resolvedVersion       string
	mu                    sync.Mutex
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

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if agentMode {
		var err error

		if talosVersion, err = c.resolveAgentModeTalosVersion(ctx); err != nil {
			return "", err
		}

		logger.Debug("resolved the agent mode Talos version", zap.String("agent_mode_talos_version", talosVersion))

		extensions = agentModeExtensions
	}

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

// agentModeTalosVersionTTL is how long a resolved agent mode Talos version is reused, matching the interval Omni refreshes its Talos versions in.
const agentModeTalosVersionTTL = 15 * time.Minute

// resolveAgentModeTalosVersion returns AgentModeTalosVersion, reusing the last result for agentModeTalosVersionTTL.
func (c *Client) resolveAgentModeTalosVersion(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.resolvedVersion != "" && time.Since(c.resolvedAt) < agentModeTalosVersionTTL {
		return c.resolvedVersion, nil
	}

	version, err := AgentModeTalosVersion(ctx, c.omniClient.Omni().State(), c.agentModeTalosVersion)
	if err != nil {
		return "", err
	}

	c.resolvedVersion, c.resolvedAt = version, time.Now()

	return version, nil
}

// defaultVersionLabel marks the TalosVersion that Omni uses as its default.
//
// todo: use omni.LabelDefaultVersion once the Omni client exposes it.
const defaultVersionLabel = "omni.sidero.dev/default-version"

// AgentModeTalosVersion returns the Talos version to boot the agent mode with.
//
// It is the override when set, otherwise the version Omni labels as its default.
// An Omni that labels no version falls back to the default of the Omni client library.
func AgentModeTalosVersion(ctx context.Context, st state.State, override string) (string, error) {
	if override != "" {
		return override, nil
	}

	versions, err := safe.StateListAll[*omni.TalosVersion](ctx, st, state.WithLabelQuery(resource.LabelExists(defaultVersionLabel)))
	if err != nil {
		return "", fmt.Errorf("failed to get the default Talos version from Omni: %w", err)
	}

	if versions.Len() > 0 {
		return versions.Get(0).TypedSpec().Value.Version, nil
	}

	return omniconstants.DefaultTalosVersion, nil
}
