// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package imagefactory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/siderolabs/image-factory/pkg/client"
	"github.com/siderolabs/image-factory/pkg/schematic"
	"github.com/siderolabs/omni/client/pkg/meta"
	"gopkg.in/yaml.v3"
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

// Client is an image factory client.
type Client struct {
	factoryClient         *client.Client
	pxeBaseURL            string
	agentModeTalosVersion string
	machineLabelsMeta     string
}

// NewClient creates a new image factory client.
func NewClient(baseURL, pxeBaseURL, agentModeTalosVersion string, machineLabels []string) (*Client, error) {
	labelsMeta, err := parseLabels(machineLabels)
	if err != nil {
		return nil, err
	}

	factoryClient, err := client.New(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		pxeBaseURL:            pxeBaseURL,
		agentModeTalosVersion: agentModeTalosVersion,
		machineLabelsMeta:     labelsMeta,

		factoryClient: factoryClient,
	}, nil
}

// SchematicIPXEURL ensures a schematic exists on the image factory and returns the iPXE URL to it.
//
// If agentMode is true, the schematic will be created with the firmware extensions and the metal-agent extension.
func (c *Client) SchematicIPXEURL(ctx context.Context, agentMode bool, talosVersion, arch string, extensions, extraKernelArgs []string) (string, error) {
	var metaValues []schematic.MetaValue

	if c.machineLabelsMeta != "" {
		metaValues = append(metaValues, schematic.MetaValue{
			Key:   meta.LabelsMeta,
			Value: c.machineLabelsMeta,
		})
	}

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

	schematicID, err := c.factoryClient.SchematicCreate(ctx, sch)
	if err != nil {
		return "", fmt.Errorf("failed to create schematic: %w", err)
	}

	ipxeURL := fmt.Sprintf("%s/pxe/%s/%s/metal-%s", c.pxeBaseURL, schematicID, talosVersion, arch)

	return ipxeURL, err
}

func parseLabels(machineLabels []string) (string, error) {
	labels := map[string]string{}

	for _, l := range machineLabels {
		parts := strings.Split(l, "=")
		if len(parts) > 2 {
			return "", fmt.Errorf("malformed label %s", l)
		}

		value := ""

		if len(parts) > 1 {
			value = parts[1]
		}

		labels[parts[0]] = value
	}

	data, err := yaml.Marshal(meta.ImageLabels{
		Labels: labels,
	})
	if err != nil {
		return "", err
	}

	return string(data), nil
}
