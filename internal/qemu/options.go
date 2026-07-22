// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package qemu

// Options for the set of machines.
type Options struct {
	Name             string
	CIDR             string
	CNIBundleURL     string
	TalosctlPath     string
	CPUs             string
	DefaultBootOrder string

	Nameservers []string

	NumMachines int
	MTU         int

	DiskSize uint64
	MemSize  int64

	UEFIEnabled bool

	// PXEBootViaDHCPD makes the provisioner's own DHCP server hand out the PXE boot server and boot
	// file directly, instead of relying on the provider's ProxyDHCP responses. Some UEFI firmware
	// PXE stacks (e.g., EDK2 in QEMU) do not accept ProxyDHCP offers, so this makes them boot via
	// plain single-server PXE. The boot file is chosen for the machine firmware and architecture
	// and served by the provider's TFTP server on the gateway address. Only applies when the
	// machines are created, not when an existing set is loaded.
	PXEBootViaDHCPD bool
}

// DefaultOptions returns the default options for the set of machines.
func DefaultOptions() Options {
	return Options{
		Name:         "bare-metal",
		CIDR:         "172.29.0.0/24",
		CNIBundleURL: "https://github.com/siderolabs/talos/releases/latest/download/talosctl-cni-bundle-amd64.tar.gz",
		NumMachines:  4,
		Nameservers:  []string{"1.1.1.1", "1.0.0.1"},
		MTU:          1440,

		CPUs:     "3",
		DiskSize: 6 * 1024 * 1024 * 1024,
		MemSize:  3072 * 1024 * 1024,

		DefaultBootOrder: "cn",
	}
}
