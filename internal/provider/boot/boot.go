// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package boot provides boot mode determination.
package boot

import (
	"errors"

	"github.com/cosi-project/runtime/pkg/resource"
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
	if infraMachine == nil {
		return Mode{}, errors.New("infra machine is nil")
	}

	if status == nil {
		return Mode{}, errors.New("machine status is nil")
	}

	tearingDown := infraMachine.Metadata().Phase() == resource.PhaseTearingDown
	accepted := infraMachine.TypedSpec().Value.Accepted
	requiredPowerMgmtConfig := status.TypedSpec().Value.PowerManagement == nil
	installed := installStatus != nil && installStatus.TypedSpec().Value.Installed
	allocated := infraMachine.TypedSpec().Value.ClusterTalosVersion != ""
	pendingWipeID := ""

	wipeID := infraMachine.TypedSpec().Value.WipeId
	if wipeID == "" {
		wipeID = "initial"
	}

	lastWipeID := status.TypedSpec().Value.LastWipeId

	if wipeID != lastWipeID {
		pendingWipeID = wipeID
	}

	bootIntoAgentMode := tearingDown || !accepted || !allocated || requiredPowerMgmtConfig || pendingWipeID != ""

	var requiredBootMode specs.BootMode

	switch {
	case bootIntoAgentMode:
		requiredBootMode = specs.BootMode_BOOT_MODE_AGENT_PXE
	case installed:
		requiredBootMode = specs.BootMode_BOOT_MODE_TALOS_DISK
	default:
		requiredBootMode = specs.BootMode_BOOT_MODE_TALOS_PXE
	}

	logger.With(
		zap.Bool("tearing_down", tearingDown),
		zap.Bool("accepted", accepted),
		zap.Bool("required_power_mgmt_config", requiredPowerMgmtConfig),
		zap.Bool("installed", installed),
		zap.String("wipe_id", wipeID),
		zap.String("last_wipe_id", lastWipeID),
		zap.Stringer("required_boot_mode", requiredBootMode),
	).Debug("determined boot mode")

	return Mode{
		PendingWipeID:           pendingWipeID,
		BootMode:                requiredBootMode,
		Installed:               installed,
		RequiresPowerMgmtConfig: requiredPowerMgmtConfig,
	}, nil
}
