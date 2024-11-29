// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package baremetal contains bare-metal related resources.
package baremetal

import (
	"github.com/siderolabs/omni/client/pkg/infra"

	providermeta "github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
)

// Namespace is the resource namespace of this provider.
var Namespace = infra.ResourceNamespace(providermeta.ProviderID)
