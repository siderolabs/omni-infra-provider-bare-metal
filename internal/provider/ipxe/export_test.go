// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe

import (
	"github.com/cosi-project/runtime/pkg/controller"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/controllers"
)

// NewTestHandler builds a Handler with only the fields needed to exercise the boot path in ServeHTTP.
func NewTestHandler(installationMediaClient InstallationMediaClient, reader controller.Reader,
	pxeBootEventCh chan<- controllers.PXEBootEvent, logger *zap.Logger,
) *Handler {
	return &Handler{
		installationMediaClient: installationMediaClient,
		reader:                  reader,
		pxeBootEventCh:          pxeBootEventCh,
		logger:                  logger,
	}
}
