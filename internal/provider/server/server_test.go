// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package server_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/server"
)

func TestMultiHandlerRouting(t *testing.T) {
	tftpDir := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(tftpDir, "amd64"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(tftpDir, "amd64", "snp.efi"), []byte("snp-efi-content"), 0o644))

	configCalled := false
	configHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		configCalled = true

		w.WriteHeader(http.StatusOK)
	})

	ipxeCalled := false
	ipxeHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		ipxeCalled = true

		w.WriteHeader(http.StatusOK)
	})

	grpcHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := zaptest.NewLogger(t)
	handler := server.NewMultiHandler(configHandler, ipxeHandler, grpcHandler, false, tftpDir, logger)

	t.Run("config path routes to configHandler", func(t *testing.T) {
		configCalled = false

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/config", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		assert.True(t, configCalled)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ipxe path routes to ipxeHandler", func(t *testing.T) {
		ipxeCalled = false

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/ipxe/boot.ipxe", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		assert.True(t, ipxeCalled)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("tftp path serves files from tftpDir", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/tftp/amd64/snp.efi", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "snp-efi-content", rec.Body.String())
	})

	t.Run("tftp path does not serve files from a different dir", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/tftp/nonexistent.efi", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
