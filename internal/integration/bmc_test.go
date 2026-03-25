// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//go:build integration_ipmi || integration_redfish

package integration_test

import (
	"context"
	"errors"
	"flag"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	goipmi "github.com/bougou/go-ipmi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/pxe"
)

var (
	bmcAddress  = flag.String("bmc-address", "", "BMC address")
	bmcUsername = flag.String("bmc-username", "", "BMC username")
	bmcPassword = flag.String("bmc-password", "", "BMC password")
)

const (
	retryInterval   = 3 * time.Second
	retryTimeout    = 10 * time.Minute
	powerStateDelay = 10 * time.Second
)

// retryableRedfishMessages contains iLO/Redfish error message substrings that are transient.
var retryableRedfishMessages = []string{
	"UnableToModifyDuringSystemPOST",
	"NodeBusy",
	"ResourceInUse",
	"ResetInProgress",
	// HTTP 412 Precondition Failed — typically a transient ETag/concurrency issue.
	"412:",
	// HTTP 409 Conflict — resource state conflict, often transient.
	"409:",
	// HTTP 503 Service Unavailable — BMC temporarily overloaded.
	"503:",
}

// isRetryable returns true if the error is transient and worth retrying.
func isRetryable(err error) bool {
	if errors.Is(err, errUnexpectedPowerState) {
		return true
	}

	// Context deadline exceeded from inner timeouts (e.g. Redfish per-request timeout) is retryable.
	// This does NOT conflict with the retry loop's own context check — that context is checked
	// separately via the select on ctx.Done().
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	// IPMI completion codes.
	var respErr *goipmi.ResponseError

	if errors.As(err, &respErr) {
		switch respErr.CompletionCode() {
		case goipmi.CompletionCodeNodeBusy,
			goipmi.CompletionCodeProcessTimeout,
			goipmi.CompletionCodeOutOfSpace,
			goipmi.CompletionCodeCannotProvideResponseSDRRInUpdate,
			goipmi.CompletionCodeCannotProvideResponseFirmwareUpdate,
			goipmi.CompletionCodeCannotProvideResponseBMCInitialize:
			return true
		}
	}

	// Redfish/iLO transient error messages.
	errMsg := err.Error()

	for _, msg := range retryableRedfishMessages {
		if strings.Contains(errMsg, msg) {
			return true
		}
	}

	return false
}

// sleep waits for the given duration or until the context is canceled.
func sleep(ctx context.Context, t *testing.T, d time.Duration) {
	t.Helper()

	select {
	case <-ctx.Done():
		require.NoError(t, ctx.Err(), "context canceled during sleep")
	case <-time.After(d):
	}
}

// retry retries a function until it succeeds, the context is canceled, or a non-retryable error occurs.
// Retryable errors (like "Node busy") are retried; all others fail immediately.
func retry(ctx context.Context, t *testing.T, description string, fn func() error) {
	t.Helper()

	deadline := time.Now().Add(retryTimeout)

	var err error

	for attempt := 1; ; attempt++ {
		if err = fn(); err == nil {
			return
		}

		if !isRetryable(err) {
			require.NoError(t, err, "%s failed with non-retryable error", description)
		}

		t.Logf("attempt %d for %s failed (retryable): %v", attempt, description, err)

		if time.Now().Add(retryInterval).After(deadline) {
			require.NoError(t, err, "%s timed out after %s", description, retryTimeout)
		}

		select {
		case <-ctx.Done():
			require.NoError(t, ctx.Err(), "%s canceled", description)
		case <-time.After(retryInterval):
		}
	}
}

// clientWrapper wraps a bmc.Client and tracks which methods have been called.
// It implements bmc.Client — if a new method is added to that interface, this wrapper
// will fail to compile until the method is added here.
type clientWrapper struct {
	inner  bmc.Client
	called map[string]bool
	mu     sync.Mutex
}

// compile-time assertion: clientWrapper must implement bmc.Client.
var _ bmc.Client = (*clientWrapper)(nil)

func newClientWrapper(inner bmc.Client) *clientWrapper {
	return &clientWrapper{
		inner:  inner,
		called: make(map[string]bool),
	}
}

func (w *clientWrapper) markCalled(name string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.called[name] = true
}

func (w *clientWrapper) Close(ctx context.Context) error {
	w.markCalled("Close")

	return w.inner.Close(ctx)
}

func (w *clientWrapper) Reboot(ctx context.Context) error {
	w.markCalled("Reboot")

	return w.inner.Reboot(ctx)
}

func (w *clientWrapper) IsPoweredOn(ctx context.Context) (bool, error) {
	w.markCalled("IsPoweredOn")

	return w.inner.IsPoweredOn(ctx)
}

func (w *clientWrapper) PowerOn(ctx context.Context) error {
	w.markCalled("PowerOn")

	return w.inner.PowerOn(ctx)
}

func (w *clientWrapper) PowerOff(ctx context.Context) error {
	w.markCalled("PowerOff")

	return w.inner.PowerOff(ctx)
}

func (w *clientWrapper) SetPXEBootOnce(ctx context.Context, mode pxe.BootMode) error {
	w.markCalled("SetPXEBootOnce")

	return w.inner.SetPXEBootOnce(ctx, mode)
}

func (w *clientWrapper) ResetBootDevice(ctx context.Context) error {
	w.markCalled("ResetBootDevice")

	return w.inner.ResetBootDevice(ctx)
}

// assertAllMethodsTested checks that every method on the bmc.Client interface was called
// at least once during the test. If a new method is added to bmc.Client and not tested,
// this will fail.
func (w *clientWrapper) assertAllMethodsTested(t *testing.T) {
	t.Helper()

	w.mu.Lock()
	defer w.mu.Unlock()

	clientType := reflect.TypeOf((*bmc.Client)(nil)).Elem()

	for i := range clientType.NumMethod() {
		method := clientType.Method(i)
		assert.True(t, w.called[method.Name], "bmc.Client method %q was not exercised by the integration test", method.Name)
	}
}

// requireBMCFlags fails the test if the shared BMC flags are not set.
func requireBMCFlags(t *testing.T) {
	t.Helper()

	if *bmcAddress == "" {
		t.Fatal("-bmc-address flag is required")
	}

	if *bmcUsername == "" {
		t.Fatal("-bmc-username flag is required")
	}

	if *bmcPassword == "" {
		t.Fatal("-bmc-password flag is required")
	}
}

// runBMCTestWithCleanup sets up the client wrapper, registers cleanup, and runs the shared BMC test steps.
func runBMCTestWithCleanup(t *testing.T, inner bmc.Client) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	t.Cleanup(cancel)

	client := newClientWrapper(inner)

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cleanupCancel()

		// Reset the boot device override back to default.
		t.Log("cleanup: resetting boot device to no-override")

		retry(cleanupCtx, t, "cleanup: reset boot device", func() error {
			return client.ResetBootDevice(cleanupCtx)
		})

		// Power off the server.
		t.Log("cleanup: powering off")

		retry(cleanupCtx, t, "cleanup: power off", func() error {
			return inner.PowerOff(cleanupCtx)
		})

		// Close is best-effort, the BMC may be unreachable after power off.
		if closeErr := client.Close(cleanupCtx); closeErr != nil {
			t.Logf("cleanup: failed to close BMC client (best-effort): %v", closeErr)
		}

		if !t.Failed() {
			client.assertAllMethodsTested(t)
		}
	})

	runBMCTests(ctx, t, client)
}

// ensurePoweredOn powers on the machine if it is not already on.
func ensurePoweredOn(ctx context.Context, t *testing.T, client *clientWrapper) {
	t.Helper()

	var poweredOn bool

	retry(ctx, t, "check power state", func() error {
		var err error
		poweredOn, err = client.IsPoweredOn(ctx)

		return err
	})

	if poweredOn {
		return
	}

	t.Log("machine is off, powering on first")

	retry(ctx, t, "PowerOn", func() error {
		return client.PowerOn(ctx)
	})

	t.Logf("waiting %s for power state to settle", powerStateDelay)
	sleep(ctx, t, powerStateDelay)
}

func runBMCTests(ctx context.Context, t *testing.T, client *clientWrapper) {
	t.Helper()

	// Step 1: Check the current power state.
	t.Run("IsPoweredOn", func(t *testing.T) {
		var (
			poweredOn bool
			err       error
		)

		retry(ctx, t, "IsPoweredOn", func() error {
			poweredOn, err = client.IsPoweredOn(ctx)

			return err
		})

		t.Logf("initial power state: powered_on=%v", poweredOn)
	})

	// Step 2: Power off the machine, then verify it is off.
	t.Run("PowerOff", func(t *testing.T) {
		// Ensure the machine is on first so we actually test PowerOff.
		ensurePoweredOn(ctx, t, client)

		retry(ctx, t, "PowerOff", func() error {
			return client.PowerOff(ctx)
		})

		t.Logf("waiting %s for power state to settle", powerStateDelay)
		sleep(ctx, t, powerStateDelay)

		retry(ctx, t, "verify powered off", func() error {
			poweredOn, powerErr := client.IsPoweredOn(ctx)
			if powerErr != nil {
				return powerErr
			}

			if poweredOn {
				return errUnexpectedPowerState
			}

			return nil
		})
	})

	// Step 3: Power on the machine, then verify it is on.
	t.Run("PowerOn", func(t *testing.T) {
		retry(ctx, t, "PowerOn", func() error {
			return client.PowerOn(ctx)
		})

		t.Logf("waiting %s for power state to settle", powerStateDelay)
		sleep(ctx, t, powerStateDelay)

		retry(ctx, t, "verify powered on", func() error {
			poweredOn, powerErr := client.IsPoweredOn(ctx)
			if powerErr != nil {
				return powerErr
			}

			if !poweredOn {
				return errUnexpectedPowerState
			}

			return nil
		})
	})

	// Step 4: Set PXE boot once with BIOS mode.
	t.Run("SetPXEBootOnce/BIOS", func(t *testing.T) {
		retry(ctx, t, "SetPXEBootOnce BIOS", func() error {
			return client.SetPXEBootOnce(ctx, pxe.BootModeBIOS)
		})
	})

	// Step 5: Set PXE boot once with UEFI mode.
	t.Run("SetPXEBootOnce/UEFI", func(t *testing.T) {
		retry(ctx, t, "SetPXEBootOnce UEFI", func() error {
			return client.SetPXEBootOnce(ctx, pxe.BootModeUEFI)
		})
	})

	// Step 6: Reboot the machine.
	t.Run("Reboot", func(t *testing.T) {
		// Ensure the machine is on — it may have powered off after PXE boot with no PXE server.
		ensurePoweredOn(ctx, t, client)

		retry(ctx, t, "Reboot", func() error {
			return client.Reboot(ctx)
		})
	})
}

type unexpectedPowerStateError struct{}

func (unexpectedPowerStateError) Error() string {
	return "unexpected power state"
}

var errUnexpectedPowerState = unexpectedPowerStateError{}
