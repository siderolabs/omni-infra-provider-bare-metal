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
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/tls"
)

// Options contains the provider options.
type Options struct {
	IPMIPXEBootMode        string
	DHCPProxyIfaceOrIP     string
	OmniAPIEndpoint        string
	ImageFactoryBaseURL    string
	ImageFactoryPXEBaseURL string
	AgentModeTalosVersion  string // todo: get this from Omni. Warning: needs to be Talos 1.9 with agent code inside
	APIListenAddress       string
	APIAdvertiseAddress    string
	APIPowerMgmtStateDir   string
	Name                   string
	Description            string
	BootFromDiskMethod     string

	MachineLabels []string

	TLS               tls.Options
	AgentClient       agent.ClientOptions
	Redfish           redfish.Options
	MinRebootInterval time.Duration
	DHCPProxyPort     int
	APIPort           int

	UseLocalBootAssets    bool
	AgentTestMode         bool
	InsecureSkipTLSVerify bool
	EnableResourceCache   bool
	ClearState            bool
	DisableDHCPProxy      bool
	SecureBootEnabled     bool
}

// DefaultOptions returns the default provider options.
func DefaultOptions() Options {
	return Options{
		Name:                   "Bare Metal",
		Description:            "Bare metal infrastructure provider",
		ImageFactoryBaseURL:    "https://factory.talos.dev",
		ImageFactoryPXEBaseURL: "https://pxe.factory.talos.dev",
		AgentModeTalosVersion:  "v1.12.0-beta.0",
		DHCPProxyPort:          67,
		BootFromDiskMethod:     string(ipxe.BootIPXEExit),
		IPMIPXEBootMode:        string(pxe.BootModeUEFI),
		APIPort:                50042,
		MinRebootInterval:      15 * time.Minute,
		Redfish:                redfish.DefaultOptions(),
		TLS: tls.Options{
			Enabled:         false,
			APIPort:         50043,
			AgentSkipVerify: false,
			CATTL:           30 * 365 * 24 * time.Hour, // 30 years
			CertTTL:         24 * time.Hour,
		},
		AgentClient: agent.DefaultClientOptions(),
	}
}
