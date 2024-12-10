// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package baremetal

import (
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/protobuf"
	"github.com/cosi-project/runtime/pkg/resource/typed"
	"github.com/siderolabs/omni/client/pkg/infra"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	providermeta "github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
)

// NewPowerStatus creates a new PowerStatus.
func NewPowerStatus(id string) *PowerStatus {
	return typed.NewResource[PowerStatusSpec, PowerStatusExtension](
		resource.NewMetadata(Namespace, PowerStatusType, id, resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.PowerStatusSpec{}),
	)
}

// PowerStatusType is the type of PowerStatus resource.
var PowerStatusType = infra.ResourceType("BareMetalPowerStatus", providermeta.ProviderID)

// PowerStatus describes power status configuration.
type PowerStatus = typed.Resource[PowerStatusSpec, PowerStatusExtension]

// PowerStatusSpec wraps specs.PowerStatusSpec.
type PowerStatusSpec = protobuf.ResourceSpec[specs.PowerStatusSpec, *specs.PowerStatusSpec]

// PowerStatusExtension providers auxiliary methods for PowerStatus resource.
type PowerStatusExtension struct{}

// ResourceDefinition implements [typed.Extension] interface.
func (PowerStatusExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{
		Type:             PowerStatusType,
		Aliases:          []resource.Type{},
		DefaultNamespace: Namespace,
		PrintColumns:     []meta.PrintColumn{},
	}
}
