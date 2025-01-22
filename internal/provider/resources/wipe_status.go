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

// NewWipeStatus creates a new WipeStatus.
func NewWipeStatus(id string) *WipeStatus {
	return typed.NewResource[WipeStatusSpec, WipeStatusExtension](
		resource.NewMetadata(Namespace(), WipeStatusType(), id, resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.WipeStatusSpec{}),
	)
}

// WipeStatusType is the type of WipeStatus resource.
func WipeStatusType() string {
	return infra.ResourceType("WipeStatus", providermeta.ProviderID.String())
}

// WipeStatus describes the resource configuration.
type WipeStatus = typed.Resource[WipeStatusSpec, WipeStatusExtension]

// WipeStatusSpec wraps specs.WipeStatusSpec.
type WipeStatusSpec = protobuf.ResourceSpec[specs.WipeStatusSpec, *specs.WipeStatusSpec]

// WipeStatusExtension providers auxiliary methods for WipeStatus resource.
type WipeStatusExtension struct{}

// ResourceDefinition implements [typed.Extension] interface.
func (WipeStatusExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{
		Type:             WipeStatusType(),
		Aliases:          []resource.Type{},
		DefaultNamespace: Namespace(),
		PrintColumns:     []meta.PrintColumn{},
	}
}
