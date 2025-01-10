// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package controllers implements COSI controllers for the bare metal provider.
package controllers

import (
	"context"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power"
)

// PowerClientFactory is the interface for creating power clients.
type PowerClientFactory interface {
	GetClient(ctx context.Context, powerManagement *specs.PowerManagement) (power.Client, error)
}
