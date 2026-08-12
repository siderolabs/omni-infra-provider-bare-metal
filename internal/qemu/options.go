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

	// VirtualBMCUsername and VirtualBMCPassword are the seeded admin credentials
	// installed on every emulated BMC.
	VirtualBMCUsername string
	VirtualBMCPassword string

	Nameservers []string

	MTU         int
	MemSize     int64
	DiskSize    uint64
	NumMachines int

	UEFIEnabled bool

	// PXEBootViaDHCPD makes the provisioner's own DHCP server hand out the PXE boot server and boot
	// file directly, instead of relying on the provider's ProxyDHCP responses. Some UEFI firmware
	// PXE stacks (e.g., EDK2 in QEMU) do not accept ProxyDHCP offers, so this makes them boot via
	// plain single-server PXE. The boot file is chosen for the machine firmware and architecture
	// and served by the provider's TFTP server on the gateway address. Only applies when the
	// machines are created, not when an existing set is loaded.
	PXEBootViaDHCPD bool

	// VirtualBMC gives each machine an emulated IPMI BMC so the provider exercises
	// its real IPMI code path instead of the fake HTTP power backend. qemu-up then
	// stays running to host the BMCs. Each BMC listens on an OS-assigned loopback
	// port and advertises it in-band. The out-of-band IPMI-over-LAN BMC runs on any
	// host. The in-band KCS device (so the agent gets a working IPMI device) is
	// added only on linux/amd64, so elsewhere the BMC is LAN-only. Only applies
	// when the machines are created.
	VirtualBMC bool
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

		VirtualBMCUsername: "ADMIN",
		VirtualBMCPassword: "ADMIN",
	}
}
