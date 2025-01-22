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

// NewPowerOperation creates a new PowerOperation.
func NewPowerOperation(id string) *PowerOperation {
	return typed.NewResource[PowerOperationSpec, PowerOperationExtension](
		resource.NewMetadata(Namespace(), PowerOperationType(), id, resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.PowerOperationSpec{}),
	)
}

// PowerOperationType is the type of PowerOperation resource.
func PowerOperationType() string {
	return infra.ResourceType("PowerOperation", providermeta.ProviderID.String())
}

// PowerOperation describes power status configuration.
type PowerOperation = typed.Resource[PowerOperationSpec, PowerOperationExtension]

// PowerOperationSpec wraps specs.PowerOperationSpec.
type PowerOperationSpec = protobuf.ResourceSpec[specs.PowerOperationSpec, *specs.PowerOperationSpec]

// PowerOperationExtension providers auxiliary methods for PowerOperation resource.
type PowerOperationExtension struct{}

// ResourceDefinition implements [typed.Extension] interface.
func (PowerOperationExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{
		Type:             PowerOperationType(),
		Aliases:          []resource.Type{},
		DefaultNamespace: Namespace(),
		PrintColumns:     []meta.PrintColumn{},
	}
}
