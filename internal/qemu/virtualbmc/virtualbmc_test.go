// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package virtualbmc_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/bougou/go-ipmi/pkg/client"
	"github.com/bougou/go-ipmi/pkg/command/app"
	"github.com/bougou/go-ipmi/pkg/command/chassis"
	"github.com/bougou/go-ipmi/pkg/types"
	"github.com/bougou/go-ipmi/pkg/vmproto"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/qemu/virtualbmc"
)

// fakePowerAPI emulates the talos launcher HTTP power API the BMC drives.
type fakePowerAPI struct {
	hits      []string
	mu        sync.Mutex
	poweredOn bool
}

func (f *fakePowerAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.hits = append(f.hits, r.URL.Path)

	switch r.URL.Path {
	case "/poweron", "/reboot":
		f.poweredOn = true
	case "/poweroff":
		f.poweredOn = false
	case "/status":
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"PoweredOn":%v}`, f.poweredOn) //nolint:errcheck

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (f *fakePowerAPI) paths() []string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return slices.Clone(f.hits)
}

// newServedBMC starts an emulated BMC backed by a fake power API and a connected
// lanplus client, and returns them wired up for a test.
func newServedBMC(t *testing.T) (*fakePowerAPI, *client.Client, *virtualbmc.Server) {
	t.Helper()

	const username, password = "ADMIN", "ADMIN"

	fake := &fakePowerAPI{}
	ts := httptest.NewServer(fake)

	t.Cleanup(ts.Close)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	srv, err := virtualbmc.New(ctx, virtualbmc.Options{
		BMCIP:       "127.0.0.1",
		PowerAPIURL: ts.URL,
		Username:    username,
		Password:    password,
		MaxUsers:    8,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	go func() { _ = srv.Serve(ctx) }() //nolint:errcheck

	addr := srv.Addr()

	c, err := client.NewClient(addr.Addr().String(), int(addr.Port()), username, password)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	c.WithInterface(client.InterfaceLanplus).WithTimeout(2 * time.Second).WithRetry(1)

	if err = c.Connect(ctx); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	t.Cleanup(func() { _ = c.Close(ctx) }) //nolint:errcheck

	return fake, c, srv
}

// TestServerLANPath drives the emulated BMC over IPMI-over-LAN exactly as the
// provider does and asserts each chassis command reaches the launcher power API.
func TestServerLANPath(t *testing.T) {
	ctx := context.Background()
	fake, c, _ := newServedBMC(t)

	if _, err := c.ChassisControl(ctx, chassis.ChassisControlPowerUp); err != nil {
		t.Fatalf("power up: %v", err)
	}

	status, err := c.GetChassisStatus(ctx)
	if err != nil {
		t.Fatalf("chassis status: %v", err)
	}

	if !status.PowerIsOn {
		t.Fatalf("expected the machine to report powered on")
	}

	if err = c.SetBootDevice(ctx, types.BootDeviceSelectorForcePXE, types.BIOSBootTypeEFI, false); err != nil {
		t.Fatalf("set boot device: %v", err)
	}

	if _, err = c.ChassisControl(ctx, chassis.ChassisControlPowerDown); err != nil {
		t.Fatalf("power down: %v", err)
	}

	for _, want := range []string{"/poweron", "/status", "/pxeboot", "/poweroff"} {
		if !slices.Contains(fake.paths(), want) {
			t.Errorf("expected power API to be called with %s, got %v", want, fake.paths())
		}
	}
}

// TestServerAdvertisesBoundPort proves the emulated BMC advertises the OS-assigned
// port it bound via Get LAN Configuration Parameters param #8, which is how the
// agent discovers a non-623 loopback BMC in-band.
func TestServerAdvertisesBoundPort(t *testing.T) {
	_, c, srv := newServedBMC(t)

	var portParam types.LanConfigParam_PrimaryRMCPPort

	if err := c.GetLanConfigParamFor(context.Background(), 1, &portParam); err != nil {
		t.Fatalf("GetLanConfigParamFor(primary RMCP port): %v", err)
	}

	if portParam.Port != srv.Addr().Port() {
		t.Errorf("advertised RMCP port = %d, want the bound port %d", portParam.Port, srv.Addr().Port())
	}
}

// vmRequest is the subset of a go-ipmi command request the in-band helper needs:
// its wire identity and its packed body.
type vmRequest interface {
	Command() types.Command
	Pack() []byte
}

// sendInBand issues one request over the VM-protocol (in-band) frontend and fails
// the test unless the BMC answers with the success completion code.
func sendInBand(t *testing.T, vmc *vmproto.Client, req vmRequest) {
	t.Helper()

	cmd := req.Command()

	cc, _, err := vmc.Command(byte(cmd.NetFn), cmd.ID, req.Pack()...)
	if err != nil {
		t.Fatalf("%s in-band: %v", cmd.Name, err)
	}

	if cc != 0 {
		t.Fatalf("%s in-band: completion code 0x%02x", cmd.Name, cc)
	}
}

// TestServerInBandUserReachableOverLAN exercises the full production sequence that
// otherwise only the CI integration run covers: the agent creates an IPMI user
// in-band over the VM-protocol frontend, and the provider then authenticates as
// that user out-of-band over IPMI-over-LAN. Both frontends share one BMC, so a
// user created in-band must be usable over LAN. It also asserts the server shuts
// down cleanly and releases its LAN port.
func TestServerInBandUserReachableOverLAN(t *testing.T) {
	const (
		adminUser, adminPass = "ADMIN", "ADMIN"
		agentUser, agentPass = "agent", "agentpass"
		agentSlot            = 3
	)

	fake := &fakePowerAPI{}
	ts := httptest.NewServer(fake)

	t.Cleanup(ts.Close)

	// The VM-protocol frontend listens on a unix socket. Keep the path short and
	// off the long macOS temp dir, which overflows the sun_path length limit, so
	// t.TempDir() is not usable here.
	socketDir, err := os.MkdirTemp("/tmp", "vbmc") //nolint:usetesting
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}

	t.Cleanup(func() { _ = os.RemoveAll(socketDir) }) //nolint:errcheck

	vmSocket := filepath.Join(socketDir, "b.sock")

	serverCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := virtualbmc.New(serverCtx, virtualbmc.Options{
		BMCIP:       "127.0.0.1",
		VMSocket:    vmSocket,
		PowerAPIURL: ts.URL,
		Username:    adminUser,
		Password:    adminPass,
		MaxUsers:    8,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	served := make(chan error, 1)

	go func() { served <- srv.Serve(serverCtx) }()

	addr := srv.Addr()

	// In-band: create and enable the LAN user the way the agent does.
	dialer := net.Dialer{Timeout: 2 * time.Second}

	conn, err := dialer.DialContext(context.Background(), "unix", vmSocket)
	if err != nil {
		t.Fatalf("dial VM socket: %v", err)
	}

	vmc := vmproto.NewClient(conn, 2*time.Second)

	sendInBand(t, vmc, &app.SetUsernameRequest{UserID: agentSlot, Username: agentUser})
	sendInBand(t, vmc, &app.SetUserAccessRequest{
		EnableChanging:      true,
		EnableIPMIMessaging: true,
		ChannelNumber:       1,
		UserID:              agentSlot,
		MaxPrivLevel:        uint8(types.PrivilegeLevelAdministrator),
	})
	sendInBand(t, vmc, &app.SetUserPasswordRequest{
		UserID:    agentSlot,
		Operation: app.PasswordOperationSetPassword,
		Password:  agentPass,
	})
	sendInBand(t, vmc, &app.SetUserPasswordRequest{UserID: agentSlot, Operation: app.PasswordOperationEnableUser})

	conn.Close() //nolint:errcheck

	// Out-of-band: authenticate as that user over LAN and drive a power command.
	c, err := client.NewClient(addr.Addr().String(), int(addr.Port()), agentUser, agentPass)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	c.WithInterface(client.InterfaceLanplus).WithTimeout(2 * time.Second).WithRetry(1)

	if err = c.Connect(context.Background()); err != nil {
		t.Fatalf("connect over LAN as the in-band user: %v", err)
	}

	if _, err = c.ChassisControl(context.Background(), chassis.ChassisControlPowerUp); err != nil {
		t.Errorf("power up as the in-band user: %v", err)
	}

	c.Close(context.Background()) //nolint:errcheck

	if !slices.Contains(fake.paths(), "/poweron") {
		t.Errorf("expected the power API to be called with /poweron, got %v", fake.paths())
	}

	// The server shuts down cleanly and releases its LAN port.
	cancel()

	select {
	case err = <-served:
		if err != nil {
			t.Fatalf("Serve returned an error on shutdown: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Serve did not return within 5s of cancellation")
	}

	reb, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: int(addr.Port())})
	if err != nil {
		t.Fatalf("LAN port %d was not released after shutdown: %v", addr.Port(), err)
	}

	reb.Close() //nolint:errcheck
}
