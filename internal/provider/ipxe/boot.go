// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe

import "fmt"

// BootFromDiskMethod defines a way to boot from disk.
type BootFromDiskMethod string

const (
	// BootIPXEExit is a method to boot from disk using iPXE script with `exit` command.
	BootIPXEExit BootFromDiskMethod = "ipxe-exit"

	// Boot404 is a method to boot from disk using HTTP 404 response to iPXE.
	Boot404 BootFromDiskMethod = "http-404"

	// BootSANDisk is a method to boot from disk using iPXE script with `sanboot` command.
	BootSANDisk BootFromDiskMethod = "ipxe-sanboot"
)

// parseBootFromDiskMethod parses a boot from disk method.
func parseBootFromDiskMethod(method string) (BootFromDiskMethod, error) {
	switch method {
	case string(BootIPXEExit):
		return BootIPXEExit, nil
	case string(Boot404):
		return Boot404, nil
	case string(BootSANDisk):
		return BootSANDisk, nil
	default:
		return "", fmt.Errorf("unknown boot from disk method: %s", method)
	}
}
