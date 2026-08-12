// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package virtualbmc runs an emulated IPMI BMC for a single QEMU machine created
// by the talos provision library. It serves IPMI-over-LAN (RMCP+) on a loopback
// address and, on hosts that attach the QEMU KCS device, the OpenIPMI VM protocol
// on a unix socket, both backed by one shared BMC whose power operations map to
// the talos launcher's HTTP power API. It lets the bare-metal provider exercise
// its real IPMI path against emulated machines instead of the fake HTTP power
// backend.
//
// The BMC runs in-process inside qemu-up (one per machine), so qemu-up learns the
// bound address directly from [New] and no cross-process handshake is needed. The
// RMCP port is not fixed: the BMC binds an OS-assigned loopback port (so any
// number coexist on one host without collision or privilege) and advertises it
// via Get LAN Configuration Parameters param #8, which the agent reads in-band.
package virtualbmc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"os"

	"github.com/bougou/go-ipmi/pkg/bmc"
	"github.com/bougou/go-ipmi/pkg/clock"
	"github.com/bougou/go-ipmi/pkg/hal"
	"github.com/bougou/go-ipmi/pkg/server"
	"github.com/bougou/go-ipmi/pkg/transport/udp"
	"github.com/bougou/go-ipmi/pkg/vmproto"
	"golang.org/x/sync/errgroup"
)

// adminUserSlot is the seeded admin's user slot. Slot 1 is reserved for the
// anonymous/null user, so the well-known admin lives in slot 2. On the full
// (linux/amd64) flow the agent creates its own user in a higher free slot. On the
// LAN-only flow the admin is the credential the provider authenticates with.
const adminUserSlot uint8 = 2

// lanChannel is the LAN channel number the default channel store models, and the
// channel the agent configures its user on.
const lanChannel uint8 = 1

// Options configures the emulated BMC.
type Options struct {
	// BMCIP is the loopback IPv4 address the BMC binds and advertises in-band.
	BMCIP string
	// VMSocket is the unix socket for the QEMU OpenIPMI VM protocol frontend.
	// Empty disables it (LAN-only mode).
	VMSocket string
	// PowerAPIURL is the base URL of the talos launcher HTTP power API.
	PowerAPIURL string
	// Username and Password are the seeded admin credentials.
	Username string
	Password string
	// MaxUsers is the highest user slot the BMC advertises to in-band enumerators.
	MaxUsers int
}

// Server is a bound but not-yet-serving emulated BMC. Create one with [New] to
// learn its bound address, then call [Server.Serve].
type Server struct {
	bmc        *bmc.BMC
	lan        *server.Server
	vmListener net.Listener
	addr       netip.AddrPort
}

// New binds the BMC's loopback UDP port (and the VM protocol unix socket when
// configured) and returns a Server that reports the bound address. Binding up
// front lets the caller fail fast and learn the OS-assigned port before serving.
// Close or Serve releases the sockets. ctx bounds only the bind operations.
func New(ctx context.Context, opts Options) (*Server, error) {
	if opts.PowerAPIURL == "" {
		return nil, errors.New("power API URL is required")
	}

	if opts.MaxUsers < int(adminUserSlot) || opts.MaxUsers > int(bmc.MaxUsers) {
		return nil, fmt.Errorf("max users must be between %d and %d", adminUserSlot, bmc.MaxUsers)
	}

	ip, err := parseIPv4(opts.BMCIP)
	if err != nil {
		return nil, fmt.Errorf("BMC IP: %w", err)
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IP(ip[:])})
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", opts.BMCIP, err)
	}

	boundAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		conn.Close() //nolint:errcheck

		return nil, fmt.Errorf("unexpected LAN listener address type %T", conn.LocalAddr())
	}

	var vmListener net.Listener

	if opts.VMSocket != "" {
		if vmListener, err = listenVMSocket(ctx, opts.VMSocket); err != nil {
			conn.Close() //nolint:errcheck

			return nil, err
		}
	}

	b := bmc.New(deviceInfo(), deviceGUID(), newLauncherHAL(opts.PowerAPIURL, hal.IPConfig{
		IP:   ip,
		Port: uint16(boundAddr.Port),
	}), bmc.WithClock(clock.Real))
	b.Users = bmc.NewUserStore(bmc.WithMaxUsers(uint8(opts.MaxUsers)))

	if err = seedAdmin(b.Users, opts.Username, opts.Password); err != nil {
		conn.Close() //nolint:errcheck

		if vmListener != nil {
			vmListener.Close() //nolint:errcheck
		}

		return nil, fmt.Errorf("seed admin user: %w", err)
	}

	return &Server{
		bmc:        b,
		lan:        server.NewServer(b, udp.Wrap(conn)),
		vmListener: vmListener,
		addr:       boundAddr.AddrPort(),
	}, nil
}

// Addr returns the loopback address the BMC bound (with its OS-assigned port).
func (s *Server) Addr() netip.AddrPort {
	return s.addr
}

// Close releases the bound sockets without serving. It is only needed to abandon
// a Server that will not be served. Serve releases them on return.
func (s *Server) Close() {
	s.lan.Close() //nolint:errcheck

	if s.vmListener != nil {
		s.vmListener.Close() //nolint:errcheck
	}
}

// Serve runs the LAN (and, when configured, in-band VM protocol) frontends until
// ctx is canceled. A VM protocol failure returns an error so the caller can tear
// the whole set down rather than leave a BMC that answers LAN with a dead
// /dev/ipmi0.
func (s *Server) Serve(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()
		s.lan.Close() //nolint:errcheck // unblocks lan.Serve on shutdown

		return nil
	})

	eg.Go(func() error {
		if err := s.lan.Serve(ctx); err != nil && !errors.Is(err, context.Canceled) {
			return fmt.Errorf("serve LAN: %w", err)
		}

		return nil
	})

	if s.vmListener != nil {
		eg.Go(func() error {
			if err := vmproto.NewVMServer(s.bmc).Serve(ctx, s.vmListener); err != nil && ctx.Err() == nil {
				return fmt.Errorf("serve VM protocol: %w", err)
			}

			return nil
		})
	}

	return eg.Wait()
}

// seedAdmin installs the well-known admin user in slot [adminUserSlot] with
// administrator access on the LAN channel.
func seedAdmin(store *bmc.UserStore, username, password string) error {
	user, err := store.Add(adminUserSlot, username)
	if err != nil {
		return err
	}

	user.SetPassword([]byte(password))
	user.Enabled = true
	user.ChannelAccess[lanChannel] = bmc.UserChannelAccess{
		MaxPrivilege: bmc.PrivilegeLevelAdministrator,
		Enabled:      true,
	}

	return nil
}

// listenVMSocket creates the unix listener QEMU's ipmi-bmc-extern chardev
// connects to, removing any stale socket left by a previous run.
func listenVMSocket(ctx context.Context, path string) (net.Listener, error) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("remove stale VM socket %s: %w", path, err)
	}

	var lc net.ListenConfig

	ln, err := lc.Listen(ctx, "unix", path)
	if err != nil {
		return nil, fmt.Errorf("listen on VM socket %s: %w", path, err)
	}

	return ln, nil
}

// deviceInfo returns the BMC identity reported by Get Device ID.
func deviceInfo() bmc.DeviceInfo {
	return bmc.DeviceInfo{
		DeviceID:                32,
		DeviceRevision:          1,
		FirmwareMajor:           1,
		FirmwareMinor:           0,
		IPMIVersion:             0x20,
		ManufacturerID:          0x000157,
		ProductID:               0x0001,
		AdditionalDeviceSupport: 0x39,
	}
}

func deviceGUID() [16]byte {
	var guid [16]byte

	copy(guid[:], "omni-bm-virt-bmc")

	return guid
}

func parseIPv4(s string) ([4]byte, error) {
	var out [4]byte

	ip := net.ParseIP(s)
	if ip == nil {
		return out, fmt.Errorf("invalid IPv4 address %q", s)
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return out, fmt.Errorf("not an IPv4 address: %q", s)
	}

	copy(out[:], ip4)

	return out, nil
}
