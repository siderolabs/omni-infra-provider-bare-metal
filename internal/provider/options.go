// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package provider

import (
	"time"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/agent"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/pxe"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/redfish"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/ipxe"
)

// Options contains the provider options.
type Options struct {
	Name                   string
	Description            string
	OmniAPIEndpoint        string
	ImageFactoryBaseURL    string
	ImageFactoryPXEBaseURL string
	AgentModeTalosVersion  string // todo: get this from Omni. Warning: needs to be Talos 1.9 with agent code inside
	APIListenAddress       string
	APIAdvertiseAddress    string
	APIPowerMgmtStateDir   string
	DHCPProxyIfaceOrIP     string
	BootFromDiskMethod     string
	IPMIPXEBootMode        string
	MachineLabels          []string
	APIPort                int

	EnableResourceCache   bool
	AgentTestMode         bool
	InsecureSkipTLSVerify bool
	UseLocalBootAssets    bool
	ClearState            bool
	DisableDHCPProxy      bool
	SecureBootEnabled     bool

	TLS         TLSOptions
	Redfish     redfish.Options
	AgentClient agent.ClientOptions

	MinRebootInterval time.Duration
}

// TLSOptions contains the TLS options.
type TLSOptions struct {
	APIPort         int
	CATTL           time.Duration
	CertTTL         time.Duration
	Enabled         bool
	AgentSkipVerify bool
}

// DefaultOptions returns the default provider options.
func DefaultOptions() Options {
	return Options{
		Name:                   "Bare Metal",
		Description:            "Bare metal infrastructure provider",
		ImageFactoryBaseURL:    "https://factory.talos.dev",
		ImageFactoryPXEBaseURL: "https://pxe.factory.talos.dev",
		AgentModeTalosVersion:  "v1.9.3",
		BootFromDiskMethod:     string(ipxe.BootIPXEExit),
		IPMIPXEBootMode:        string(pxe.BootModeUEFI),
		APIPort:                50042,
		MinRebootInterval:      5 * time.Minute,
		Redfish:                redfish.DefaultOptions(),
		TLS: TLSOptions{
			Enabled:         false,
			APIPort:         50043,
			AgentSkipVerify: false,
			CATTL:           30 * 365 * 24 * time.Hour, // 30 years
			CertTTL:         24 * time.Hour,
		},
		AgentClient: agent.DefaultClientOptions(),
	}
}
