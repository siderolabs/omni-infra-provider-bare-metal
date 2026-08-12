// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package qemu

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/siderolabs/talos/pkg/provision"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/qemu/virtualbmc"
)

const (
	// bmcLoopbackIP is the address every emulated BMC binds and advertises. All
	// BMCs share it and are told apart by an OS-assigned port, so no bridge
	// address or privileged setup is needed, and the provider (host-networked on
	// Linux CI, native on macOS) reaches them over loopback. Each BMC advertises
	// its port in-band via Get LAN Config param #8.
	bmcLoopbackIP = "127.0.0.1"

	// bmcMaxUsers is the highest user slot the emulated BMCs advertise: the seeded
	// admin (slot 2) plus headroom for the agent-created user.
	bmcMaxUsers = 8
)

// attachKCS reports whether the guest gets an in-band KCS device (and thus the
// full agent flow). It requires the QEMU IPMI devices, which only exist on the
// x86 machine types. Elsewhere the emulated BMC is LAN-only.
func attachKCS() bool {
	return runtime.GOOS == "linux" && runtime.GOARCH == "amd64"
}

func (machines *Machines) statePath() string {
	return filepath.Join(machines.stateDir, machines.options.Name)
}

func bmcSocketPath(statePath, nodeName string) string {
	return filepath.Join(statePath, nodeName+"-bmc.sock")
}

// bmcExtraQEMUArgs attaches an external-BMC KCS device forwarding to the BMC's
// unix socket, giving the guest kernel a /dev/ipmi0. reconnectOpt lets QEMU
// tolerate the socket not existing yet (qemu-up binds it after Create returns)
// and reconnect on each guest power cycle.
func bmcExtraQEMUArgs(socketPath, reconnectOpt string) []string {
	return []string{
		"-chardev", fmt.Sprintf("socket,id=ipmi0,path=%s,%s", socketPath, reconnectOpt),
		"-device", "ipmi-bmc-extern,id=bmc0,chardev=ipmi0",
		"-device", "isa-ipmi-kcs,bmc=bmc0",
	}
}

// qemuReconnectOpt returns the chardev reconnect option the QEMU binary the
// launcher will run supports: reconnect-ms (QEMU >= 9.2) or the legacy reconnect
// (seconds). QEMU rejects the unsupported spelling outright at parse time, so the
// arg must match the exact binary. It probes the same candidates the talos qemu
// provisioner selects from and, when it cannot tell, prefers reconnect-ms (the
// modern spelling, correct for the overwhelmingly common QEMU >= 9.2).
func qemuReconnectOpt() string {
	binary := ""

	for _, candidate := range []string{"qemu-system-x86_64", "qemu-kvm", "/usr/libexec/qemu-kvm"} {
		if path, err := exec.LookPath(candidate); err == nil {
			binary = path

			break
		}
	}

	if binary == "" {
		return "reconnect-ms=2000"
	}

	out, _ := exec.Command(binary, "-chardev", "socket,help").CombinedOutput() //nolint:errcheck,noctx

	switch {
	case strings.Contains(string(out), "reconnect-ms"):
		return "reconnect-ms=2000"
	case strings.Contains(string(out), "reconnect"):
		return "reconnect=2"
	default:
		return "reconnect-ms=2000"
	}
}

// serveBMCs hosts one emulated BMC per machine in-process and blocks until ctx is
// canceled or a BMC fails. qemu-up stays alive as the BMC supervisor for the
// lifetime of the machine set. When it exits, the BMCs stop with it (the talos
// launcher processes it spawned keep running until a destroy).
func (machines *Machines) serveBMCs(ctx context.Context, cluster provision.Cluster) error {
	info := cluster.Info()

	if len(info.Network.GatewayAddrs) == 0 {
		return errors.New("cluster has no gateway address")
	}

	gateway := info.Network.GatewayAddrs[0]
	statePath := machines.statePath()

	// qemu-up creates its machines as PXE nodes, which the provisioner returns in
	// ExtraNodes rather than Nodes, so serve both.
	nodes := slices.Concat(info.Nodes, info.ExtraNodes)

	servers := make([]*virtualbmc.Server, 0, len(nodes))

	for _, node := range nodes {
		vmSocket := ""
		if attachKCS() {
			vmSocket = bmcSocketPath(statePath, node.Name)
		}

		srv, err := virtualbmc.New(ctx, virtualbmc.Options{
			BMCIP:       bmcLoopbackIP,
			VMSocket:    vmSocket,
			PowerAPIURL: "http://" + net.JoinHostPort(gateway.String(), strconv.Itoa(node.APIPort)),
			Username:    machines.options.VirtualBMCUsername,
			Password:    machines.options.VirtualBMCPassword,
			MaxUsers:    bmcMaxUsers,
		})
		if err != nil {
			for _, started := range servers {
				started.Close()
			}

			return fmt.Errorf("failed to start virtual BMC for node %q: %w", node.Name, err)
		}

		servers = append(servers, srv)

		machines.logger.Info("virtual BMC ready", zap.String("node", node.Name), zap.String("address", srv.Addr().String()))
	}

	machines.logger.Info("all machines and virtual BMCs ready", zap.Int("count", len(servers)))

	eg, egCtx := errgroup.WithContext(ctx)

	for _, srv := range servers {
		eg.Go(func() error {
			return srv.Serve(egCtx)
		})
	}

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
