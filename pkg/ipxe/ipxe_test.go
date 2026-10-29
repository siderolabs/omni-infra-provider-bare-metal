// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/omni-infra-provider-bare-metal/pkg/ipxe"
)

func TestScriptsHandler(t *testing.T) {
	t.Parallel()

	initScript, err := ipxe.BuildInitScript("192.168.1.100", 50042)
	require.NoError(t, err)

	var bootQuery url.Values

	bootScriptHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		bootQuery = req.URL.Query()

		w.Write([]byte("#!ipxe\nboot")) //nolint:errcheck
	})

	mux := http.NewServeMux()
	mux.Handle("/"+ipxe.ScriptsURLPath+"/", ipxe.ScriptsHandler(initScript, bootScriptHandler))

	get := func(t *testing.T, target string) *httptest.ResponseRecorder {
		t.Helper()

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, target, nil))

		return rec
	}

	// the init script chains to the boot script, so the URL it chains to must reach the boot script handler
	_, chainLine, ok := strings.Cut(string(initScript), "chain --replace ")
	require.True(t, ok, "the init script does not chain")

	chainURL := strings.ReplaceAll(strings.Fields(chainLine)[0], "${", "") // iPXE variables, filled in by iPXE at boot time

	parsed, err := url.Parse(chainURL)
	require.NoError(t, err)

	rec := get(t, parsed.RequestURI())
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "#!ipxe\nboot", rec.Body.String())
	assert.Contains(t, bootQuery, "uuid", "the boot script handler must get the machine's details in the query")

	rec = get(t, "/"+ipxe.ScriptsURLPath+"/init.ipxe")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, string(initScript), rec.Body.String())

	rec = get(t, "/"+ipxe.ScriptsURLPath+"/unknown.ipxe")
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
