// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package machine provides utilities for determining the required state of a machine.
package machine

import (
	"github.com/cosi-project/runtime/pkg/resource"
	omnispecs "github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
)

// IsInstalled returns true if the machine is installed.
func IsInstalled(infraMachine *infra.Machine, wipeStatus *resources.WipeStatus) bool {
	if infraMachine == nil {
		return false
	}

	installEventID := infraMachine.TypedSpec().Value.InstallEventId
	lastWipeInstallEventID := uint64(0)

	if wipeStatus != nil {
		lastWipeInstallEventID = wipeStatus.TypedSpec().Value.LastWipeInstallEventId
	}

	return installEventID > lastWipeInstallEventID
}

// RequiresWipe returns true if the machine needs to be wiped.
func RequiresWipe(infraMachine *infra.Machine, wipeStatus *resources.WipeStatus) bool {
	// maybe check acceptance here (or here as well)
	if infraMachine == nil || wipeStatus == nil || !wipeStatus.TypedSpec().Value.InitialWipeDone {
		return true
	}

	return infraMachine.TypedSpec().Value.WipeId != wipeStatus.TypedSpec().Value.LastWipeId
}

// RequiredBootMode returns the required boot mode for the machine.
func RequiredBootMode(infraMachine *infra.Machine, bmcConfiguration *resources.BMCConfiguration, wipeStatus *resources.WipeStatus, logger *zap.Logger) specs.BootMode {
	installed := IsInstalled(infraMachine, wipeStatus)
	requiresWipe := RequiresWipe(infraMachine, wipeStatus)
	acceptanceStatus := omnispecs.InfraMachineConfigSpec_PENDING
	infraMachineTearingDown := false
	allocated := false

	if infraMachine != nil {
		acceptanceStatus = infraMachine.TypedSpec().Value.AcceptanceStatus
		infraMachineTearingDown = infraMachine.Metadata().Phase() == resource.PhaseTearingDown
		allocated = infraMachine.TypedSpec().Value.ClusterTalosVersion != ""
	}

	acceptancePending := acceptanceStatus == omnispecs.InfraMachineConfigSpec_PENDING
	rejected := acceptanceStatus == omnispecs.InfraMachineConfigSpec_REJECTED
	requiresPowerMgmtConfig := bmcConfiguration == nil

	bootIntoAgentMode := infraMachineTearingDown || acceptancePending || !allocated || requiresPowerMgmtConfig || requiresWipe

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
		zap.Bool("infra_machine_tearing_down", infraMachineTearingDown),
		zap.Bool("requires_power_mgmt_config", requiresPowerMgmtConfig),
		zap.Bool("installed", installed),
		zap.Stringer("acceptance_status", acceptanceStatus),
		zap.Stringer("required_boot_mode", requiredBootMode),
	).Debug("determined boot mode")

	return requiredBootMode
}
