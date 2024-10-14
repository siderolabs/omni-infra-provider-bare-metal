// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package qemu

// Options for the set of machines.
type Options struct {
	Name         string
	CIDR         string
	CNIBundleURL string
	TalosctlPath string
	CPUs         string

	Nameservers []string

	NumMachines int
	MTU         int

	DiskSize uint64
	MemSize  int64
}

// DefaultOptions are the default options for the set of machines.
var DefaultOptions = Options{
	Name:         "bare-metal",
	CIDR:         "172.42.0.0/24",
	CNIBundleURL: "https://github.com/siderolabs/talos/releases/latest/download/talosctl-cni-bundle-amd64.tar.gz",
	NumMachines:  4,
	Nameservers:  []string{"1.1.1.1", "1.0.0.1"},
	MTU:          1440,

	CPUs:     "3",
	DiskSize: 6 * 1024 * 1024 * 1024,
	MemSize:  3072 * 1024 * 1024,
}
