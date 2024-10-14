// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package boot provides boot mode determination.
package boot

import (
	"github.com/cosi-project/runtime/pkg/resource"
	omnispecs "github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/baremetal"
)

// Mode describes the required boot mode.
type Mode struct {
	PendingWipeID           string
	BootMode                specs.BootMode
	Installed               bool
	RequiresPowerMgmtConfig bool
}

// DetermineRequiredMode determines the required boot mode.
func DetermineRequiredMode(infraMachine *infra.Machine, status *baremetal.MachineStatus, installStatus *infra.MachineState, logger *zap.Logger) (Mode, error) {
	acceptanceStatus := omnispecs.InfraMachineConfigSpec_PENDING
	tearingDown := false
	allocated := false
	requiresPowerMgmtConfig := true
	installed := false
	wipeID := "initial"
	lastWipeID := ""

	if infraMachine != nil {
		acceptanceStatus = infraMachine.TypedSpec().Value.AcceptanceStatus
		tearingDown = infraMachine.Metadata().Phase() == resource.PhaseTearingDown
		allocated = infraMachine.TypedSpec().Value.ClusterTalosVersion != ""

		if infraMachine.TypedSpec().Value.WipeId != "" {
			wipeID = infraMachine.TypedSpec().Value.WipeId
		}
	}

	if status != nil {
		requiresPowerMgmtConfig = status.TypedSpec().Value.PowerManagement == nil
		lastWipeID = status.TypedSpec().Value.LastWipeId
	}

	if installStatus != nil {
		installed = installStatus.TypedSpec().Value.Installed
	}

	acceptancePending := acceptanceStatus == omnispecs.InfraMachineConfigSpec_PENDING
	rejected := acceptanceStatus == omnispecs.InfraMachineConfigSpec_REJECTED

	pendingWipeID := ""

	if wipeID != lastWipeID {
		pendingWipeID = wipeID
	}

	requiresWipe := pendingWipeID != ""
	bootIntoAgentMode := tearingDown || acceptancePending || !allocated || requiresPowerMgmtConfig || requiresWipe

	var requiredBootMode specs.BootMode

	switch {
	case rejected:
		requiredBootMode = specs.BootMode_BOOT_MODE_TALOS_DISK
	case bootIntoAgentMode:
		requiredBootMode = specs.BootMode_BOOT_MODE_AGENT_PXE
	case installed:
		requiredBootMode = specs.BootMode_BOOT_MODE_TALOS_DISK
	default:
		requiredBootMode = specs.BootMode_BOOT_MODE_TALOS_PXE
	}

	logger.With(
		zap.Bool("tearing_down", tearingDown),
		zap.Bool("requires_power_mgmt_config", requiresPowerMgmtConfig),
		zap.Bool("installed", installed),
		zap.String("wipe_id", wipeID),
		zap.String("last_wipe_id", lastWipeID),
		zap.Stringer("acceptance_status", acceptanceStatus),
		zap.Stringer("required_boot_mode", requiredBootMode),
	).Debug("determined boot mode")

	return Mode{
		PendingWipeID:           pendingWipeID,
		BootMode:                requiredBootMode,
		Installed:               installed,
		RequiresPowerMgmtConfig: requiresPowerMgmtConfig,
	}, nil
}
