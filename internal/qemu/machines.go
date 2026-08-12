// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package qemu

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/netip"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	talosnet "github.com/siderolabs/net"
	clientconfig "github.com/siderolabs/talos/pkg/machinery/client/config"
	"github.com/siderolabs/talos/pkg/provision"
	"github.com/siderolabs/talos/pkg/provision/providers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapio"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/util"
)

const (
	providerName = "qemu"
	diskDriver   = "virtio"

	// TalosctlBinary is the default name of the talosctl binary.
	TalosctlBinary = "talosctl"
)

// Machines represents a set of Talos QEMU machines.
type Machines struct {
	logger      *zap.Logger
	subnetCIDR  netip.Prefix
	stateDir    string
	cniDir      string
	nameservers []netip.Addr
	options     Options
	nanoCPUs    int64
}

// New creates a new set of Talos QEMU machines.
func New(options Options, logger *zap.Logger) (*Machines, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	talosDir, err := clientconfig.GetTalosDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to get Talos directory: %w", err)
	}

	stateDir := filepath.Join(talosDir, "clusters")
	cniDir := filepath.Join(talosDir, "cni")

	logger = logger.With(zap.String("cluster_name", options.Name), zap.String("state_dir", stateDir))

	if options.TalosctlPath == "" {
		if options.TalosctlPath, err = findTalosctl(); err != nil {
			return nil, fmt.Errorf("failed to find talosctl binary: %w", err)
		}
	}

	cidr, err := netip.ParsePrefix(options.CIDR)
	if err != nil {
		return nil, fmt.Errorf("failed to parse subnet CIDR: %w", err)
	}

	nss, err := parseNameservers(options.Nameservers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse nameservers: %w", err)
	}

	nanoCPUs, err := parseCPUShare(options.CPUs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CPU share: %w", err)
	}

	return &Machines{
		options:     options,
		stateDir:    stateDir,
		cniDir:      cniDir,
		subnetCIDR:  cidr,
		nameservers: nss,
		nanoCPUs:    nanoCPUs,
		logger:      logger,
	}, nil
}

// Run starts the provisioner by provisioning a QEMU cluster with the given number of machines. If the cluster already exists, it is loaded instead.
func (machines *Machines) Run(ctx context.Context) error {
	machines.logger.Info(
		"provision a new set of machines",
		zap.String("name", machines.options.Name),
		zap.String("subnet_cidr", machines.subnetCIDR.String()),
		zap.Int("num_machines", machines.options.NumMachines),
		zap.Int64("mem_size", machines.options.MemSize),
		zap.Uint64("disk_size", machines.options.DiskSize),
		zap.Int64("nano_cpus", machines.nanoCPUs),
		zap.Strings("nameservers", machines.options.Nameservers),
		zap.Int("mtu", machines.options.MTU),
		zap.String("talosctl_path", machines.options.TalosctlPath),
		zap.String("cni_bundle_url", machines.options.CNIBundleURL),
	)

	qemuProvisioner, err := providers.Factory(ctx, providerName)
	if err != nil {
		return fmt.Errorf("failed to create provisioner: %w", err)
	}

	existingCluster, err := machines.loadExisting(ctx, qemuProvisioner)
	if err != nil {
		return err
	}

	if existingCluster != nil {
		if machines.options.PXEBootViaDHCPD {
			machines.logger.Warn("an existing machine set was loaded, so the PXE-boot-via-DHCPD option has no effect - destroy and recreate the machines for it to apply")
		}

		// The emulated BMCs live inside this qemu-up process, so a fresh qemu-up
		// cannot re-attach to an existing set: it would bind new ports and start
		// with an empty user store, orphaning any machine already accepted against
		// the previous run's BMCs.
		if machines.options.VirtualBMC {
			return errors.New("virtual BMC machines cannot be re-attached - destroy and recreate the machines")
		}

		return nil
	}

	logWriter := &zapio.Writer{
		Log:   machines.logger,
		Level: zapcore.InfoLevel,
	}

	defer util.LogClose(logWriter, machines.logger)

	cluster, err := machines.createNew(ctx, qemuProvisioner, logWriter)
	if err != nil {
		return err
	}

	// With virtual BMC, qemu-up becomes the BMC supervisor: it hosts one emulated
	// BMC per machine in-process and blocks here until signaled.
	if machines.options.VirtualBMC {
		return machines.serveBMCs(ctx, cluster)
	}

	return nil
}

// loadExisting returns the already-provisioned cluster, or nil when none exists.
func (machines *Machines) loadExisting(ctx context.Context, qemuProvisioner provision.Provisioner) (provision.Cluster, error) {
	existingCluster, err := qemuProvisioner.Reflect(ctx, machines.options.Name, machines.stateDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to load existing cluster: %w", err)
		}

		machines.logger.Info("no existing set of machines found")

		return nil, nil //nolint:nilnil // nil cluster means "none exists", not an error
	}

	machines.logger.Info("loaded existing set of machines")

	// The machines are PXE nodes, which the provisioner reports in ExtraNodes, so
	// count both to compare against the requested number.
	info := existingCluster.Info()

	if existing := len(info.Nodes) + len(info.ExtraNodes); existing != machines.options.NumMachines {
		machines.logger.Warn("number of existing machines does not match the requested number of machines",
			zap.Int("requested", machines.options.NumMachines),
			zap.Int("existing", existing))
	}

	return existingCluster, nil
}

func (machines *Machines) createNew(ctx context.Context, qemuProvisioner provision.Provisioner, logWriter io.Writer) (provision.Cluster, error) {
	machines.logger.Info("create a new set of machines")

	gatewayAddr, err := talosnet.NthIPInNetwork(machines.subnetCIDR, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway address: %w", err)
	}

	var tftpServer, ipxeBootFilename string

	if machines.options.PXEBootViaDHCPD {
		tftpServer = gatewayAddr.String()

		// This mirrors the firmware detection the provider's ProxyDHCP does, which this mode
		// bypasses. The guest arch matches runtime.GOARCH because qemu-up always creates
		// host-native VMs.
		switch {
		case !machines.options.UEFIEnabled:
			ipxeBootFilename = "undionly.kpxe"
		case runtime.GOARCH == "arm64":
			ipxeBootFilename = "snp-arm64.efi"
		default:
			ipxeBootFilename = "snp.efi"
		}
	}

	nodes := make([]provision.NodeRequest, 0, machines.options.NumMachines)

	// Attach the in-band KCS device only where QEMU has the IPMI devices. Probe
	// the reconnect chardev option once for the whole set.
	attachBMCKCS := machines.options.VirtualBMC && attachKCS()

	bmcReconnectOpt := ""
	if attachBMCKCS {
		bmcReconnectOpt = qemuReconnectOpt()
	}

	for i := range machines.options.NumMachines {
		nodeUUID := uuid.New()

		machines.logger.Info("generated node UUID", zap.String("uuid", nodeUUID.String()))

		ip, ipErr := talosnet.NthIPInNetwork(machines.subnetCIDR, i+2)
		if ipErr != nil {
			return nil, fmt.Errorf("failed to calculate offset %d from CIDR %s: %w", i+2, machines.subnetCIDR, ipErr)
		}

		nodeReq := provision.NodeRequest{
			Name: nodeUUID.String(),

			IPs:      []netip.Addr{ip},
			Memory:   machines.options.MemSize,
			NanoCPUs: machines.nanoCPUs,
			Disks: []*provision.Disk{
				{
					Size:   machines.options.DiskSize,
					Driver: diskDriver,
				},
			},
			SkipInjectingConfig: true,
			UUID:                &nodeUUID,
			PXEBooted:           true,
			DefaultBootOrder:    machines.options.DefaultBootOrder, // first disk, then network
			TFTPServer:          tftpServer,
			IPXEBootFilename:    ipxeBootFilename,
		}

		if attachBMCKCS {
			nodeReq.ExtraQEMUArgs = bmcExtraQEMUArgs(bmcSocketPath(machines.statePath(), nodeReq.Name), bmcReconnectOpt)
		}

		nodes = append(nodes, nodeReq)
	}

	request := provision.ClusterRequest{
		Name: machines.options.Name,

		Network: provision.NetworkRequest{
			Name:         machines.options.Name,
			CIDRs:        []netip.Prefix{machines.subnetCIDR},
			GatewayAddrs: []netip.Addr{gatewayAddr},
			MTU:          machines.options.MTU,
			Nameservers:  machines.nameservers,
			CNI: provision.CNIConfig{
				BinPath:   []string{filepath.Join(machines.cniDir, "bin")},
				ConfDir:   filepath.Join(machines.cniDir, "conf.d"),
				CacheDir:  filepath.Join(machines.cniDir, "cache"),
				BundleURL: machines.options.CNIBundleURL,
			},
		},

		SelfExecutable: machines.options.TalosctlPath,
		StateDirectory: machines.stateDir,

		Nodes: nodes,
	}

	cluster, err := qemuProvisioner.Create(
		ctx, request,
		provision.WithBootloader(true),
		provision.WithLogWriter(logWriter),
		provision.WithUEFI(machines.options.UEFIEnabled), // Note: UEFI doesn't work correctly on PXE timeout in QEMU, as it drops to UEFI shell
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create machines: %w", err)
	}

	return cluster, nil
}

// Destroy destroys the existing set of machines.
func (machines *Machines) Destroy(ctx context.Context) error {
	defer util.LogErr(func() error {
		return os.RemoveAll(filepath.Join(machines.stateDir, machines.options.Name))
	}, machines.logger)

	// The emulated BMCs are goroutines inside the supervising qemu-up process, so
	// there is nothing for destroy to tear down here: whoever runs destroy must
	// stop that process first (which stops its BMCs). Destroy only removes the
	// machines and network.

	qemuProvisioner, err := providers.Factory(ctx, providerName)
	if err != nil {
		return fmt.Errorf("failed to create provisioner: %w", err)
	}

	// attempt to load existing cluster
	cluster, err := qemuProvisioner.Reflect(ctx, machines.options.Name, machines.stateDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("failed to load existing cluster while clearing state: %w", err)
	}

	logWriter := &zapio.Writer{
		Log:   machines.logger,
		Level: zapcore.InfoLevel,
	}
	defer util.LogClose(logWriter, machines.logger)

	if err = qemuProvisioner.Destroy(ctx, cluster, provision.WithDeleteOnErr(true), provision.WithLogWriter(logWriter)); err != nil {
		if strings.Contains(err.Error(), "no such network interface") {
			return nil
		}

		return fmt.Errorf("failed to destroy cluster: %w", err)
	}

	return nil
}

func parseCPUShare(cpus string) (int64, error) {
	cpu, ok := new(big.Rat).SetString(cpus)
	if !ok {
		return 0, fmt.Errorf("failed to parse as a rational number: %s", cpus)
	}

	nano := cpu.Mul(cpu, big.NewRat(1e9, 1))
	if !nano.IsInt() {
		return 0, errors.New("value is too precise")
	}

	return nano.Num().Int64(), nil
}

func findTalosctl() (string, error) {
	// check the current working directory
	if stat, err := os.Stat(TalosctlBinary); err == nil && !stat.IsDir() {
		return TalosctlBinary, nil
	}

	// check the PATH
	path, err := exec.LookPath(TalosctlBinary)
	if err != nil {
		return "", fmt.Errorf("failed to find talosctl binary: %w", err)
	}

	return path, nil
}

func parseNameservers(nameservers []string) ([]netip.Addr, error) {
	nss := make([]netip.Addr, 0, len(nameservers))

	for _, ns := range nameservers {
		addr, err := netip.ParseAddr(ns)
		if err != nil {
			return nil, fmt.Errorf("failed to parse nameserver %q: %w", ns, err)
		}

		nss = append(nss, addr)
	}

	return nss, nil
}
