// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package resources

import (
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/protobuf"
	"github.com/cosi-project/runtime/pkg/resource/typed"
	"github.com/siderolabs/omni/client/pkg/infra"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	providermeta "github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
)

// NewRebootStatus creates a new RebootStatus.
func NewRebootStatus(id string) *RebootStatus {
	return typed.NewResource[RebootStatusSpec, RebootStatusExtension](
		resource.NewMetadata(Namespace(), RebootStatusType(), id, resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.RebootStatusSpec{}),
	)
}

// RebootStatusType is the type of RebootStatus resource.
func RebootStatusType() string {
	return infra.ResourceType("RebootStatus", providermeta.ProviderID.String())
}

// RebootStatus describes the resource configuration.
type RebootStatus = typed.Resource[RebootStatusSpec, RebootStatusExtension]

// RebootStatusSpec wraps specs.RebootStatusSpec.
type RebootStatusSpec = protobuf.ResourceSpec[specs.RebootStatusSpec, *specs.RebootStatusSpec]

// RebootStatusExtension providers auxiliary methods for RebootStatus resource.
type RebootStatusExtension struct{}

// ResourceDefinition implements [typed.Extension] interface.
func (RebootStatusExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{
		Type:             RebootStatusType(),
		Aliases:          []resource.Type{},
		DefaultNamespace: Namespace(),
		PrintColumns:     []meta.PrintColumn{},
	}
}
