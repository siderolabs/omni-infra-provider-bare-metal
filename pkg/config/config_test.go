// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package config_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/siderolabs/omni-infra-provider-bare-metal/pkg/config"
)

// TestKernelArgReachesHandler checks that the URL in the talos.config kernel argument is served by the Handler,
// mounted as its documentation says.
func TestKernelArgReachesHandler(t *testing.T) {
	t.Parallel()

	arg := config.KernelArg("192.168.1.100", 50042)

	configURL, ok := strings.CutPrefix(arg, "talos.config=")
	require.True(t, ok, "unexpected kernel argument %q", arg)

	parsed, err := url.Parse(strings.ReplaceAll(configURL, "${uuid}", "4c4c4544-0000-1000-8000-000000000000"))
	require.NoError(t, err)

	assert.Equal(t, "192.168.1.100:50042", parsed.Host)

	mux := http.NewServeMux()
	mux.Handle("/"+config.URLPath, config.NewHandler([]byte("machine config"), zaptest.NewLogger(t)))

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, parsed.RequestURI(), nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "machine config", rec.Body.String())
}
