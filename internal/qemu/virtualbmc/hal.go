// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package virtualbmc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/bougou/go-ipmi/pkg/hal"
	"github.com/bougou/go-ipmi/pkg/types"
)

// launcherHAL implements hal.HAL for a single QEMU machine. Chassis power
// operations are backed by the talos launcher's HTTP power API, and the network
// configuration reports the machine's emulated BMC LAN address (read in-band by
// the agent). Every other subsystem is absent.
type launcherHAL struct {
	chassis *chassisHAL
	network *networkHAL
}

func newLauncherHAL(powerAPIURL string, ipCfg hal.IPConfig) *launcherHAL {
	return &launcherHAL{
		chassis: &chassisHAL{
			baseURL: powerAPIURL,
			client:  &http.Client{Timeout: 10 * time.Second},
		},
		network: &networkHAL{cfg: ipCfg},
	}
}

func (h *launcherHAL) Chassis() hal.ChassisHAL { return h.chassis }
func (h *launcherHAL) Network() hal.NetworkHAL { return h.network }
func (h *launcherHAL) Sensors() hal.SensorHAL  { return nil }
func (h *launcherHAL) Storage() hal.StorageHAL { return nil }
func (h *launcherHAL) GPIO() hal.GPIOHAL       { return nil }
func (h *launcherHAL) I2C() hal.I2CHAL         { return nil }
func (h *launcherHAL) Console() hal.ConsoleHAL { return nil }
func (h *launcherHAL) Close() error            { return nil }

// chassisHAL maps IPMI chassis power and boot-flag operations onto the talos
// launcher's HTTP power API (/poweron, /poweroff, /reboot, /pxeboot, /status).
// That API is the only way to power a QEMU machine back on, because the launcher
// runs QEMU with -no-reboot and relaunches it from its own control loop.
type chassisHAL struct {
	client    *http.Client
	bootFlags *types.BootOptionParam_BootFlags
	baseURL   string
	mu        sync.Mutex
}

func (c *chassisHAL) post(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body) //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("power API %s returned %s: %s", path, resp.Status, bytes.TrimSpace(body))
	}

	return nil
}

// PowerState reports whether the launcher currently has the VM powered on.
func (c *chassisHAL) PowerState(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/status", nil)
	if err != nil {
		return false, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("power API /status returned %s", resp.Status)
	}

	var status struct {
		PoweredOn bool
	}

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return false, fmt.Errorf("decode power status: %w", err)
	}

	return status.PoweredOn, nil
}

// SetPower powers the VM on or off.
func (c *chassisHAL) SetPower(ctx context.Context, on bool) error {
	if on {
		return c.post(ctx, "/poweron")
	}

	return c.post(ctx, "/poweroff")
}

// PowerCycle, ColdReset and WarmReset all map to the launcher's reboot, which
// stops and relaunches the VM. The provider only ever issues a power cycle. The
// reset variants are mapped for completeness.
func (c *chassisHAL) PowerCycle(ctx context.Context) error { return c.post(ctx, "/reboot") }
func (c *chassisHAL) ColdReset(ctx context.Context) error  { return c.post(ctx, "/reboot") }
func (c *chassisHAL) WarmReset(ctx context.Context) error  { return c.post(ctx, "/reboot") }

// Identify is a no-op: QEMU has no chassis identify LED to pulse.
func (c *chassisHAL) Identify(context.Context, uint8) error { return nil }

// IntrusionState is not modeled.
func (c *chassisHAL) IntrusionState(context.Context) (bool, error) {
	return false, hal.ErrNotSupported
}

// SetBootFlags arms a one-time PXE boot when the selector requests ForcePXE.
// The launcher only distinguishes "next boot is network" from "next boot uses
// the configured order", so any other selector (in particular NoOverride, which
// the provider sends to clear the override) is a no-op: the launcher already
// defaults to its configured boot order. The full flags are retained so
// Get System Boot Options can read them back.
func (c *chassisHAL) SetBootFlags(ctx context.Context, flags *types.BootOptionParam_BootFlags) error {
	c.mu.Lock()
	c.bootFlags = flags
	c.mu.Unlock()

	if flags != nil && flags.BootDeviceSelector == types.BootDeviceSelectorForcePXE {
		return c.post(ctx, "/pxeboot")
	}

	return nil
}

// GetBootFlags returns the last boot flags written via SetBootFlags.
func (c *chassisHAL) GetBootFlags(context.Context) (*types.BootOptionParam_BootFlags, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.bootFlags == nil {
		return nil, hal.ErrNotSupported
	}

	flags := *c.bootFlags

	return &flags, nil
}

// SetBootInfoAcknowledge is accepted and dropped: the launcher does not track
// boot initiator identity.
func (c *chassisHAL) SetBootInfoAcknowledge(context.Context, *types.BootOptionParam_BootInfoAcknowledge) error {
	return nil
}

// GetBootInfoAcknowledge is not modeled.
func (c *chassisHAL) GetBootInfoAcknowledge(context.Context) (*types.BootOptionParam_BootInfoAcknowledge, error) {
	return nil, hal.ErrNotSupported
}

// networkHAL reports the emulated BMC's fixed LAN configuration. The agent reads
// it in-band (Get LAN Configuration Parameters) to learn the address on which
// the provider will later reach the BMC over IPMI-over-LAN.
type networkHAL struct {
	cfg hal.IPConfig
}

func (n *networkHAL) GetConfig(context.Context) (*hal.IPConfig, error) {
	cfg := n.cfg

	return &cfg, nil
}

// SetConfig is rejected: the emulated BMC's network configuration is fixed at
// launch and the agent only reads it.
func (n *networkHAL) SetConfig(context.Context, *hal.IPConfig) error {
	return hal.ErrNotSupported
}
