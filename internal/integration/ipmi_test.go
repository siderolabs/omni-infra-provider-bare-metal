// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// This is a manual integration test that exercises IPMI BMC operations against a real BMC endpoint.
// It is behind the "integration_ipmi" build tag and is excluded from regular test runs.
//
// Run with:
//
//	go test -tags integration_ipmi -v -timeout 30m ./internal/integration/... \
//	  -bmc-address <bmc-ip> -bmc-username <user> -bmc-password <pass>
//
// The test will power cycle the target machine. Make sure it is safe to do so.

//go:build integration_ipmi

package integration_test

import (
	"context"
	"flag"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/ipmi"
)

var ipmiPort = flag.Uint("ipmi-port", 623, "IPMI port")

func TestIPMI(t *testing.T) {
	requireBMCFlags(t)

	inner, err := ipmi.NewClient(context.Background(), &specs.BMCConfigurationSpec_IPMI{
		Address:  *bmcAddress,
		Port:     uint32(*ipmiPort),
		Username: *bmcUsername,
		Password: *bmcPassword,
	})
	require.NoError(t, err, "failed to create IPMI client")

	runBMCTestWithCleanup(t, inner)
}
