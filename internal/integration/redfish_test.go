// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// This is a manual integration test that exercises Redfish BMC operations against a real BMC endpoint.
// It is behind the "integration_redfish" build tag and is excluded from regular test runs.
//
// Run with:
//
//	go test -tags integration_redfish -v -timeout 30m ./internal/integration/... \
//	  -bmc-address <bmc-ip> -bmc-username <user> -bmc-password <pass>
//
// The test will power cycle the target machine. Make sure it is safe to do so.

//go:build integration_redfish

package integration_test

import (
	"flag"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/redfish"
)

var (
	redfishPort                      = flag.Int("redfish-port", 443, "Redfish port")
	redfishUseHTTPS                  = flag.Bool("redfish-use-https", true, "Use HTTPS for Redfish")
	redfishInsecureSkipTLS           = flag.Bool("redfish-insecure-skip-tls", true, "Skip TLS verification for Redfish")
	redfishSetBootSourceOverrideMode = flag.Bool("redfish-set-boot-source-override-mode", true, "Set boot source override mode for Redfish")
)

func TestRedfish(t *testing.T) {
	requireBMCFlags(t)

	inner := redfish.NewClient(
		redfish.Options{
			UseHTTPS:                  *redfishUseHTTPS,
			InsecureSkipTLSVerify:     *redfishInsecureSkipTLS,
			Port:                      *redfishPort,
			SetBootSourceOverrideMode: *redfishSetBootSourceOverrideMode,
		},
		*bmcAddress,
		*bmcUsername,
		*bmcPassword,
		zaptest.NewLogger(t),
	)

	runBMCTestWithCleanup(t, inner)
}
