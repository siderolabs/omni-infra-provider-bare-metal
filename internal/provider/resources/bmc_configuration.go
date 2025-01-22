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

// NewBMCConfiguration creates a new BMCConfiguration.
func NewBMCConfiguration(id string) *BMCConfiguration {
	return typed.NewResource[BMCConfigurationSpec, BMCConfigurationExtension](
		resource.NewMetadata(Namespace(), BMCConfigurationType(), id, resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.BMCConfigurationSpec{}),
	)
}

// BMCConfigurationType is the type of BMCConfiguration resource.
func BMCConfigurationType() string {
	return infra.ResourceType("BMCConfiguration", providermeta.ProviderID.String())
}

// BMCConfiguration describes the resource configuration.
type BMCConfiguration = typed.Resource[BMCConfigurationSpec, BMCConfigurationExtension]

// BMCConfigurationSpec wraps specs.BMCConfigurationSpec.
type BMCConfigurationSpec = protobuf.ResourceSpec[specs.BMCConfigurationSpec, *specs.BMCConfigurationSpec]

// BMCConfigurationExtension providers auxiliary methods for BMCConfiguration resource.
type BMCConfigurationExtension struct{}

// ResourceDefinition implements [typed.Extension] interface.
func (BMCConfigurationExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{
		Type:             BMCConfigurationType(),
		Aliases:          []resource.Type{},
		DefaultNamespace: Namespace(),
		PrintColumns:     []meta.PrintColumn{},
	}
}
