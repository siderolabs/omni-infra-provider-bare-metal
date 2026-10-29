// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	provideripxe "github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/ipxe"
	"github.com/siderolabs/omni-infra-provider-bare-metal/pkg/ipxe"
)

// TestPatchBinariesKeySet locks the contract between the served file names and the names the DHCP
// proxy and DHCPD advertise: the flat, arch-suffixed TFTP names and the per-arch HTTP names.
func TestPatchBinariesKeySet(t *testing.T) {
	script := []byte("#!ipxe\nchain http://example.org/boot")

	files, err := ipxe.PatchBinaries(provideripxe.Binaries(), script, nil, zaptest.NewLogger(t))
	require.NoError(t, err)

	expected := []string{
		"ipxe.efi", "snp.efi",
		"ipxe-arm64.efi", "snp-arm64.efi",
		"amd64/ipxe.efi", "amd64/snp.efi",
		"arm64/ipxe.efi", "arm64/snp.efi",
		"undionly.kpxe", "undionly.kpxe.0",
	}

	assert.Len(t, files, len(expected))

	for _, name := range expected {
		contents, ok := files[name]

		assert.True(t, ok, "missing file %q", name)
		assert.NotEmpty(t, contents, "file %q is empty", name)
	}

	// the EFI binaries are patched uncompressed, so the script must appear in them verbatim
	assert.True(t, bytes.Contains(files["snp.efi"], script))
}

func TestPatchBinariesCustomCA(t *testing.T) {
	// The CA hash is the SHA-256 of the DER bytes in the PEM block, which are not parsed further.
	der := []byte("custom CA certificate")
	customCAPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	customCAHash := sha256.Sum256(der)

	files, err := ipxe.PatchBinaries(provideripxe.Binaries(), []byte("#!ipxe"), customCAPEM, zaptest.NewLogger(t))
	require.NoError(t, err)

	assert.True(t, bytes.Contains(files["snp.efi"], customCAHash[:]), "the custom CA hash is not in the patched binary")

	_, err = ipxe.PatchBinaries(provideripxe.Binaries(), []byte("#!ipxe"), []byte("not PEM"), zaptest.NewLogger(t))
	assert.ErrorContains(t, err, "failed to decode PEM block")
}

func TestPatchBytes(t *testing.T) {
	logger := zaptest.NewLogger(t)

	placeholder := []byte("head # *PLACEHOLDER START*\n0123456789\n# *PLACEHOLDER END* tail")

	t.Run("success", func(t *testing.T) {
		patched, err := ipxe.PatchBytes(bytes.Clone(placeholder), []byte("#!ipxe"), nil, logger)
		require.NoError(t, err)

		assert.True(t, bytes.Contains(patched, []byte("#!ipxe")))
		assert.Len(t, patched, len(placeholder), "patching must preserve the binary size")
	})

	t.Run("placeholder missing", func(t *testing.T) {
		_, err := ipxe.PatchBytes([]byte("no placeholder here"), []byte("#!ipxe"), nil, logger)
		assert.ErrorContains(t, err, "placeholder start not found")
	})

	t.Run("end before start", func(t *testing.T) {
		_, err := ipxe.PatchBytes([]byte("# *PLACEHOLDER END* # *PLACEHOLDER START*"), []byte("#!ipxe"), nil, logger)
		assert.ErrorContains(t, err, "placeholder end before start")
	})

	t.Run("oversized script", func(t *testing.T) {
		script := bytes.Repeat([]byte{'x'}, len(placeholder)+1)

		_, err := ipxe.PatchBytes(bytes.Clone(placeholder), script, nil, logger)
		assert.ErrorContains(t, err, "larger than placeholder space")
	})
}
