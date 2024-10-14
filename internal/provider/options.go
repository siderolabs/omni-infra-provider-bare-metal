// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package provider

import "github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/ipxe"

// Options contains the provider options.
type Options struct {
	Name                   string
	Description            string
	OmniAPIEndpoint        string
	ImageFactoryBaseURL    string
	ImageFactoryPXEBaseURL string
	AgentModeTalosVersion  string
	APIListenAddress       string
	APIAdvertiseAddress    string
	APIPowerMgmtStateDir   string
	DHCPProxyIfaceOrIP     string
	BootFromDiskMethod     string
	MachineLabels          []string
	APIPort                int

	EnableResourceCache   bool
	AgentTestMode         bool
	InsecureSkipTLSVerify bool
	UseLocalBootAssets    bool
	ClearState            bool
	WipeWithZeroes        bool
}

// DefaultOptions returns the default provider options.
var DefaultOptions = Options{
	Name:                   "Bare Metal",
	Description:            "Bare metal infrastructure provider",
	OmniAPIEndpoint:        "",
	ImageFactoryBaseURL:    "https://factory.talos.dev",
	ImageFactoryPXEBaseURL: "https://pxe.factory.talos.dev",
	AgentModeTalosVersion:  "v1.9.0-alpha.2",
	BootFromDiskMethod:     string(ipxe.BootIPXEExit),
	APIPort:                50042,
}
