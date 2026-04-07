// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package dhcp_test

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/dhcp"
)

func TestIsBootDHCP(t *testing.T) {
	t.Parallel()

	pxeDiscover := newPXEPacket(t, dhcpv4.MessageTypeDiscover)
	pxeRequest := newPXEPacket(t, dhcpv4.MessageTypeRequest)
	nonPXEDiscover := newNonPXEPacket(t, dhcpv4.MessageTypeDiscover)

	tests := []struct {
		pkt     *dhcpv4.DHCPv4
		name    string
		wantErr string
		port    int
	}{
		{
			name: "port 67 accepts DHCPDISCOVER",
			pkt:  pxeDiscover,
			port: dhcp.Port67,
		},
		{
			name:    "port 67 rejects DHCPREQUEST",
			pkt:     pxeRequest,
			port:    dhcp.Port67,
			wantErr: "not DISCOVER",
		},
		{
			name: "port 4011 accepts DHCPREQUEST",
			pkt:  pxeRequest,
			port: dhcp.Port4011,
		},
		{
			name:    "port 4011 rejects DHCPDISCOVER",
			pkt:     pxeDiscover,
			port:    dhcp.Port4011,
			wantErr: "not REQUEST",
		},
		{
			name:    "port 67 rejects non-PXE DHCPDISCOVER",
			pkt:     nonPXEDiscover,
			port:    dhcp.Port67,
			wantErr: "missing option 93",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := dhcp.IsBootDHCP(tt.pkt, tt.port)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOfferDHCP(t *testing.T) {
	t.Parallel()

	const (
		apiAddr = "192.168.1.100"
		apiPort = 50042
	)

	t.Run("port 67 responds with DHCPOFFER", func(t *testing.T) {
		t.Parallel()

		req := newPXEPacket(t, dhcpv4.MessageTypeDiscover)
		resp, err := dhcp.OfferDHCP(req, apiAddr, apiPort, dhcp.FirmwareX86EFI, dhcp.Port67)
		require.NoError(t, err)

		assert.Equal(t, dhcpv4.MessageTypeOffer, resp.MessageType())
		assert.Equal(t, "snp.efi", resp.BootFileNameOption())
		assert.NotNil(t, resp.GetOneOption(dhcpv4.OptionClassIdentifier))
	})

	t.Run("port 4011 responds with DHCPACK", func(t *testing.T) {
		t.Parallel()

		req := newPXEPacket(t, dhcpv4.MessageTypeRequest)
		resp, err := dhcp.OfferDHCP(req, apiAddr, apiPort, dhcp.FirmwareX86EFI, dhcp.Port4011)
		require.NoError(t, err)

		assert.Equal(t, dhcpv4.MessageTypeAck, resp.MessageType())
		assert.Equal(t, "snp.efi", resp.BootFileNameOption())
		assert.NotNil(t, resp.GetOneOption(dhcpv4.OptionClassIdentifier))
	})

	t.Run("firmware types produce correct boot filenames", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			wantFile string
			fwtype   dhcp.Firmware
		}{
			{fwtype: dhcp.FirmwareX86PC, wantFile: "undionly.kpxe"},
			{fwtype: dhcp.FirmwareX86EFI, wantFile: "snp.efi"},
			{fwtype: dhcp.FirmwareARMEFI, wantFile: "snp-arm64.efi"},
			{fwtype: dhcp.FirmwareX86Ipxe, wantFile: "tftp://192.168.1.100/undionly.kpxe"},
			{fwtype: dhcp.FirmwareX86HTTP, wantFile: "http://192.168.1.100:50042/tftp/amd64/snp.efi"},
			{fwtype: dhcp.FirmwareARMHTTP, wantFile: "http://192.168.1.100:50042/tftp/arm64/snp.efi"},
		}

		for _, tt := range tests {
			req := newPXEPacket(t, dhcpv4.MessageTypeDiscover)
			resp, err := dhcp.OfferDHCP(req, apiAddr, apiPort, tt.fwtype, dhcp.Port67)
			require.NoError(t, err)

			assert.Equal(t, tt.wantFile, resp.BootFileNameOption(), "firmware type %d", tt.fwtype)
		}
	})
}

func newPXEPacket(t *testing.T, msgType dhcpv4.MessageType) *dhcpv4.DHCPv4 {
	t.Helper()

	pkt, err := dhcpv4.New(
		dhcpv4.WithMessageType(msgType),
		dhcpv4.WithHwAddr([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
	)
	require.NoError(t, err)

	// Option 93: Client System Architecture (EFI x86_64)
	pkt.UpdateOption(dhcpv4.OptClientArch(iana.EFI_X86_64))

	return pkt
}

func newNonPXEPacket(t *testing.T, msgType dhcpv4.MessageType) *dhcpv4.DHCPv4 {
	t.Helper()

	pkt, err := dhcpv4.New(
		dhcpv4.WithMessageType(msgType),
		dhcpv4.WithHwAddr([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
	)
	require.NoError(t, err)

	return pkt
}
