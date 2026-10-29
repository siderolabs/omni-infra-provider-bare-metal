// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/ipxe"
)

func TestValidateBootAssets(t *testing.T) {
	logger := zaptest.NewLogger(t)

	writeAssets := func(t *testing.T, dir, arch string, size int) {
		t.Helper()

		for _, name := range []string{"kernel-" + arch, "initramfs-metal-" + arch + ".xz", "cmdline-metal-" + arch} {
			require.NoError(t, os.WriteFile(filepath.Join(dir, name), bytes.Repeat([]byte{'x'}, size), 0o644))
		}
	}

	t.Run("single complete arch is enough", func(t *testing.T) {
		dir := t.TempDir()
		writeAssets(t, dir, "amd64", 16)

		assert.NoError(t, ipxe.ValidateBootAssets(dir, logger))
	})

	t.Run("both arches complete", func(t *testing.T) {
		dir := t.TempDir()
		writeAssets(t, dir, "amd64", 16)
		writeAssets(t, dir, "arm64", 16)

		assert.NoError(t, ipxe.ValidateBootAssets(dir, logger))
	})

	t.Run("empty directory", func(t *testing.T) {
		assert.ErrorContains(t, ipxe.ValidateBootAssets(t.TempDir(), logger), "no complete set")
	})

	t.Run("missing directory", func(t *testing.T) {
		assert.ErrorContains(t, ipxe.ValidateBootAssets(filepath.Join(t.TempDir(), "nope"), logger), "no complete set")
	})

	t.Run("empty files are unusable", func(t *testing.T) {
		dir := t.TempDir()
		writeAssets(t, dir, "amd64", 0)

		assert.ErrorContains(t, ipxe.ValidateBootAssets(dir, logger), "is empty")
	})

	t.Run("directory in place of a file", func(t *testing.T) {
		dir := t.TempDir()
		writeAssets(t, dir, "amd64", 16)
		require.NoError(t, os.Remove(filepath.Join(dir, "kernel-amd64")))
		require.NoError(t, os.Mkdir(filepath.Join(dir, "kernel-amd64"), 0o755))

		assert.ErrorContains(t, ipxe.ValidateBootAssets(dir, logger), "not a regular file")
	})
}
