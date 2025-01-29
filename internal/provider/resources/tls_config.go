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

// TLSConfigID is the ID of the TLSConfig resource.
const TLSConfigID = "tls-config"

// NewTLSConfig creates a new TLSConfig.
func NewTLSConfig() *TLSConfig {
	return typed.NewResource[TLSConfigSpec, TLSConfigExtension](
		resource.NewMetadata(Namespace(), TLSConfigType(), TLSConfigID, resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.TLSConfigSpec{}),
	)
}

// TLSConfigType is the type of TLSConfig resource.
func TLSConfigType() string {
	return infra.ResourceType("TLSConfig", providermeta.ProviderID.String())
}

// TLSConfig describes the resource configuration.
type TLSConfig = typed.Resource[TLSConfigSpec, TLSConfigExtension]

// TLSConfigSpec wraps specs.TLSConfigSpec.
type TLSConfigSpec = protobuf.ResourceSpec[specs.TLSConfigSpec, *specs.TLSConfigSpec]

// TLSConfigExtension providers auxiliary methods for TLSConfig resource.
type TLSConfigExtension struct{}

// ResourceDefinition implements [typed.Extension] interface.
func (TLSConfigExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{
		Type:             TLSConfigType(),
		Aliases:          []resource.Type{},
		DefaultNamespace: Namespace(),
		PrintColumns:     []meta.PrintColumn{},
	}
}
