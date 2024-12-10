// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package power provides power management functionality for machines.
package power

import (
	"context"
	"fmt"
	"io"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/api"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/ipmi"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/pxe"
)

// ErrNoPowerManagementInfo is returned when there is no power management info present yet for a machine.
var ErrNoPowerManagementInfo = fmt.Errorf("no power management info found")

// Client is the interface to interact with a single machine to send power commands to it.
type Client interface {
	io.Closer
	Reboot(ctx context.Context) error
	IsPoweredOn(ctx context.Context) (bool, error)
	PowerOn(ctx context.Context) error
	PowerOff(ctx context.Context) error
	SetPXEBootOnce(ctx context.Context, mode pxe.BootMode) error
}

// GetClient returns a power management client for the given bare metal machine.
func GetClient(mgmt *specs.PowerManagement) (Client, error) {
	if mgmt == nil {
		return nil, ErrNoPowerManagementInfo
	}

	apiInfo := mgmt.Api
	if apiInfo != nil {
		return api.NewClient(apiInfo)
	}

	ipmiInfo := mgmt.Ipmi
	if ipmiInfo != nil {
		return ipmi.NewClient(ipmiInfo)
	}

	return nil, ErrNoPowerManagementInfo
}
