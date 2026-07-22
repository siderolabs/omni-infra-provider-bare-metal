// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe

import (
	"github.com/cosi-project/runtime/pkg/controller"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/controllers"
)

// BootScriptName exposes the boot script path to tests.
var BootScriptName = bootScriptName

// PatchBinaries exposes patchBinaries to tests.
var PatchBinaries = patchBinaries

// PatchBytes exposes patchBytes to tests.
var PatchBytes = patchBytes

// NewTestHandler builds a Handler with only the fields needed to exercise the boot path in ServeHTTP.
func NewTestHandler(imageFactoryClient ImageFactoryClient, reader controller.Reader,
	pxeBootEventCh chan<- controllers.PXEBootEvent, logger *zap.Logger,
) *Handler {
	return &Handler{
		imageFactoryClient: imageFactoryClient,
		reader:             reader,
		pxeBootEventCh:     pxeBootEventCh,
		logger:             logger,
	}
}
