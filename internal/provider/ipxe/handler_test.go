// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/cosi-project/runtime/pkg/state"
	"github.com/cosi-project/runtime/pkg/state/impl/inmem"
	"github.com/cosi-project/runtime/pkg/state/impl/namespaced"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/controllers"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/ipxe"
)

// stubImageFactory returns a fixed PXE URL so the agent-mode boot path produces a body without network access.
type stubImageFactory struct{}

func (stubImageFactory) SchematicIPXEURL(context.Context, bool, string, string, []string, []string) (string, error) {
	return "https://factory.example/pxe/abc", nil
}

func serveBoot(ctx context.Context, h *ipxe.Handler, uuid string) *httptest.ResponseRecorder {
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/ipxe/boot.ipxe?uuid="+uuid+"&arch=amd64", nil)
	req.SetPathValue("script", ipxe.BootScriptName)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	return rec
}

// TestBootScriptServesScriptAndReportsBootEvent checks that a boot request returns the flushed script with a
// Content-Length and hands the PXE boot event to the channel.
func TestBootScriptServesScriptAndReportsBootEvent(t *testing.T) {
	ch := make(chan controllers.PXEBootEvent, 1)
	h := ipxe.NewTestHandler(stubImageFactory{}, state.WrapCore(namespaced.NewState(inmem.Build)), ch, zaptest.NewLogger(t))

	rec := serveBoot(t.Context(), h, "machine-0")

	require.Equal(t, http.StatusOK, rec.Code)
	require.True(t, rec.Flushed, "response was not flushed to the client")

	body := rec.Body.String()
	require.Contains(t, body, "https://factory.example/pxe/abc")
	require.Equal(t, strconv.Itoa(len(body)), rec.Header().Get("Content-Length"))

	select {
	case ev := <-ch:
		require.Equal(t, "machine-0", ev.MachineID)
	default:
		t.Fatal("PXE boot event was not delivered to the channel")
	}
}

// TestBootScriptDoesNotBlockUnderHerd fires many concurrent boot requests at a channel that nothing drains,
// mimicking the controller loop tied up polling BMCs. Every request must still return its script promptly; the
// boot events are dropped by design once the buffer is full.
func TestBootScriptDoesNotBlockUnderHerd(t *testing.T) {
	const herd = 300

	ch := make(chan controllers.PXEBootEvent) // unbuffered and undrained: sends can never succeed
	h := ipxe.NewTestHandler(stubImageFactory{}, state.WrapCore(namespaced.NewState(inmem.Build)), ch, zaptest.NewLogger(t))

	var wg sync.WaitGroup

	recs := make([]*httptest.ResponseRecorder, herd)

	for i := range herd {
		wg.Go(func() {
			recs[i] = serveBoot(t.Context(), h, fmt.Sprintf("machine-%d", i))
		})
	}

	served := make(chan struct{})

	go func() { wg.Wait(); close(served) }()

	select {
	case <-served:
	case <-time.After(5 * time.Second):
		t.Fatal("boot.ipxe handlers blocked under a thundering herd")
	}

	for _, rec := range recs {
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "https://factory.example/pxe/abc")
	}
}
