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

// NewMachineStatus creates a new MachineStatus.
func NewMachineStatus(id string) *MachineStatus {
	return typed.NewResource[MachineStatusSpec, MachineStatusExtension](
		resource.NewMetadata(Namespace(), MachineStatusType(), id, resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.MachineStatusSpec{}),
	)
}

// MachineStatusType returns the type of MachineStatus resource.
func MachineStatusType() string {
	return infra.ResourceType("BareMetalMachineStatus", providermeta.ProviderID.String())
}

// MachineStatus describes machine status configuration.
type MachineStatus = typed.Resource[MachineStatusSpec, MachineStatusExtension]

// MachineStatusSpec wraps specs.MachineStatusSpec.
type MachineStatusSpec = protobuf.ResourceSpec[specs.MachineStatusSpec, *specs.MachineStatusSpec]

// MachineStatusExtension providers auxiliary methods for MachineStatus resource.
type MachineStatusExtension struct{}

// ResourceDefinition implements [typed.Extension] interface.
func (MachineStatusExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{
		Type:             MachineStatusType(),
		Aliases:          []resource.Type{},
		DefaultNamespace: Namespace(),
		PrintColumns:     []meta.PrintColumn{},
	}
}
