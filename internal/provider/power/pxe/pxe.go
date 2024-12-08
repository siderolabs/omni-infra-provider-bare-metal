// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package pxe contains types related to PXE booting.
package pxe

import "fmt"

// BootMode is the PXE boot mode to be used.
type BootMode string

const (
	// BootModeBIOS is the mode to boot from disk using BIOS.
	BootModeBIOS BootMode = "bios"

	// BootModeUEFI is the mode to boot from disk using UEFI.
	BootModeUEFI BootMode = "uefi"
)

// ParseBootMode parses a boot mode.
func ParseBootMode(mode string) (BootMode, error) {
	switch mode {
	case string(BootModeBIOS):
		return BootModeBIOS, nil
	case string(BootModeUEFI):
		return BootModeUEFI, nil
	default:
		return "", fmt.Errorf("unknown boot mode: %s", mode)
	}
}
