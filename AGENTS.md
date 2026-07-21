# omni-infra-provider-bare-metal — agent guide

This is the ever-growing knowledge base for the project.
Maintain it as you go: whenever you learn something durable about how this repo works, add it here.
It is not append-only, so fix or delete anything that becomes wrong or outdated.
The goal is for this file to keep growing and keep getting more correct over time.
Only capture timeless, project-general knowledge here, never in-flight work or an individual's local setup.

Markdown note: this file is linted (`markdownlint` with the `sentences-per-line` rule), so keep one sentence per line and the first line as the single top-level heading.

`CLAUDE.md` and `GEMINI.md` just import this file, so this is the one place to edit.

## What this repo is

This repo is the Omni bare-metal infra provider ("BM provider"), a static infrastructure provider that bridges Omni to physical machines.
It PXE-boots bare-metal machines into Talos "agent mode", registers them with Omni, drives their BMCs for power control, wipes disks, and hands machines off to be installed and joined to clusters.
It is the automated alternative to manually registering a bare-metal machine (Matchbox or a downloaded ISO).

Infra providers come in two flavors.
Dynamic (also called provisioning) providers create and destroy machines on some infrastructure, while a static provider works with a fixed pool of machines.
This is currently the only static provider, and bare metal makes that unlikely to change, since the machines exist regardless of the provider.

It is the successor of Sidero Metal, which is in maintenance mode, and all new bare-metal work happens here.
Unlike Sidero Metal, which was built on Kubernetes Cluster API, this provider depends on Omni.

## The cast and who publishes/consumes what

- Omni is the control plane.
  It owns cluster state, machine acceptance, SideroLink, and config generation.
  The provider authenticates to it with a service account key and reads/writes `infra.*` resources.
- The BM provider (this repo) runs the PXE stack (DHCP proxy, TFTP, iPXE HTTP), the BMC control (IPMI/Redfish/API), and the COSI controllers, and it drives the agent over a reverse tunnel.
- `talos-metal-agent` ("agent mode") is a daemon that runs inside Talos booted in agent mode on the machine.
  It dials the provider and exposes a small gRPC API.
- Talos Linux is the OS, and it has a dedicated "agent mode" boot path.
- The Image Factory builds and serves boot assets for a given Talos version plus schematic (extension set plus kernel args).
- `siderolabs/extensions` is the source and catalog of Talos system extensions, including `guest-agents/metal-agent`, which packages the agent binary.

## Agent mode and why it exists

Agent mode is a special Talos boot mode where Talos does not boot any on-disk install, connects to Omni over SideroLink, and runs the metal-agent extension, all before the machine is installed or joined to a cluster.

It exists so Omni and the provider get a foothold on a physical machine before any cluster config or install exists, and without the provider needing inbound network access to the machine.
The agent solves reachability by dialing out from the machine to the provider and opening a reverse gRPC tunnel (`github.com/jhump/grpctunnel`), so the provider can invoke RPCs on the machine over that single outbound connection.
Routing to a specific machine uses an affinity key equal to the machine UUID carried in gRPC metadata.

The agent exposes five RPCs (`talos-metal-agent/api/agent/agent.proto`): `Hello`, `GetPowerManagement`, `SetPowerManagement`, `Reboot`, and `WipeDisks`.
The power-management RPCs let the agent bootstrap out-of-band BMC access from inside the OS (read the local BMC IP, create an IPMI user), so the provider gains standalone IPMI/Redfish control that survives even when the agent is not running.

## The two boot paths

The iPXE handler (`internal/provider/ipxe/handler.go`, `bootIntoAgentMode`) chooses between two paths based on `--use-local-boot-assets`.

The image factory path is the production default.
The provider asks the Image Factory for a schematic (`internal/provider/imagefactory/client.go`, `SchematicIPXEURL` with `agentMode=true`), forcing `talosVersion = AgentModeTalosVersion` and a fixed extension set of firmware extensions plus `siderolabs/metal-agent` with no version.
The factory resolves `siderolabs/metal-agent` against its per-Talos-version official-extensions catalog and errors if the extension is not published for that Talos version.
So the agent version served this way is whatever the extensions catalog pins for `AgentModeTalosVersion`, and advancing it means getting `AgentModeTalosVersion` onto a Talos version whose catalog pins the desired agent version, by waiting for such a catalog or bumping the setting to one.

The local boot-assets path is for dev and airgap and is enabled by `--use-local-boot-assets`.
The boot-assets image is baked into the provider image at `/assets` at build time, pinned in `.kres.yaml` as a `copyFrom` stage and copied in the generated `Dockerfile`.
The flag only decides whether those already-baked files get served.
Important: when `--use-local-boot-assets` is set, `--agent-mode-talos-version` has no effect, because the agent binary comes from the baked-in initramfs, not the factory.
When debugging "my agent fix did not take", first check the provider version and whether local assets are in use.

## Machine lifecycle

The heart of the lifecycle is `machine.RequiredBootMode` (`internal/provider/machine/machine.go`), which resolves to one of three boot modes: `agent-pxe` (PXE into agent mode), `talos-pxe` (PXE into Talos to install), and `talos-disk` (boot the installed system from disk).

1. The provider starts, authenticates to Omni, patches iPXE binaries, and starts the COSI runtime, TFTP (69), DHCP proxy (67 and 4011), and the API server (default 50042).
2. An unknown machine network-boots.
   The DHCP proxy (it does not replace the site DHCP, it only answers PXE clients) points it at the provider's TFTP, which chainloads iPXE, which fetches the boot script, and the provider serves the agent-mode boot.
3. Talos boots agent mode, connects SideroLink to Omni (appearing as a pending `InfraMachine`), and the agent opens its reverse tunnel.
   The machine shows up in Omni's pending list.
4. On acceptance, unless BMC credentials were supplied manually, the provider discovers and provisions them via the agent, then wipes all disks and powers the machine off, leaving it a ready pool member.
   Concretely, the provider asks the agent for power management information (the agent reads the BMC address from inside the OS), generates IPMI credentials, and has the agent create the IPMI user.
   Acceptance is destructive by design.
5. On allocation to a cluster, the provider powers the machine on with a one-time PXE boot into Talos (non-agent), serving a Talos schematic built from the cluster's Talos version, extensions, and kernel args.
   Talos boots maintenance mode, and Omni drives the install to disk and the reboot into the installed system without the provider's involvement, after which the machine joins the cluster.
6. In steady state the machine runs the installed Talos from disk, whether by its own firmware boot order or because the provider hands it off to disk on a PXE request, so it just serves in the cluster until something changes.
7. On deprovision (cluster destroy or scale-down, which look identical from the machine's side) teardown happens in two phases.
   First Omni resets the machine, which removes the machine config but leaves Talos on disk, and the machine reboots into maintenance mode from that on-disk Talos.
   Then the provider does its stronger reset: it reboots the machine one time into agent mode (the only mode where a wipe can run), wipes Talos off the disk entirely, marks the machine ready to use, and powers it off, returning it to the pool as if fresh.
   The provider can afford to erase Talos from disk because it can always PXE-boot the machine again, so PXE is effectively a permanent recovery medium.

## Boot, install, and power decisions

A few high-level facts explain how the provider drives a machine.
The exact predicates live in `internal/provider/machine`, so treat this as the intent and read the code for the precise conditions.

- The provider decides whether a machine is installed from a wipe-aware count of Talos install events, never by reading the disk.
  Omni counts install events, and each wipe remembers the count at wipe time, so a machine counts as installed only once an install event arrives after its most recent wipe.
  This is how the provider can trust that a freshly wiped machine is not installed without ever inspecting its disk.
- A machine becomes ready to use, and eligible for cluster config, once it has power management configured and no pending wipe.
- The boot decision picks one of three outcomes for a machine that asks the provider what to boot: agent mode, maintenance Talos to be installed, or its own disk.
  Boot from disk is served as an explicit hand-off, where the provider tells the machine to move on to its next boot device instead of booting it over the network, so the machine's firmware boot order does not have to put disk first.
  Disk-first is recommended and is what the emulated tests and `talosctl cluster create` use, since then an installed machine boots without even consulting the provider, but network-first works too because the provider hands off to disk explicitly.
- The provider only changes what a machine boots by setting a one-time network-boot flag over the BMC right before a power-on or reboot, which it needs when switching a machine into agent mode to wipe it or into maintenance Talos to install.
- Power follows work: the provider keeps a machine on while it is allocated, installed, or waiting to be wiped, and otherwise honors the configured preferred power state, which defaults to off.
  Idle machines default to off because the provider can power one on the instant a cluster needs it, so there is no reason to keep it running.

## Provider internals

- Provider to Omni: COSI runtime over gRPC, authed via `OMNI_SERVICE_ACCOUNT_KEY`, watching and writing `infra.*` resources.
- Provider to agent: reverse gRPC tunnel on the API port, affinity-routed by machine UUID (`internal/provider/agent`).
- Provider to machines: DHCP proxy (`internal/provider/dhcp`), TFTP (`internal/provider/tftp`), and the HTTP iPXE handler (`internal/provider/ipxe`).
- Provider to Image Factory: HTTP (`internal/provider/imagefactory`).
- Provider to BMC: IPMI, Redfish, or API backends behind one `Client` interface (`internal/provider/bmc`), with `bmc/pxe` holding the BIOS/UEFI PXE-boot abstraction.
  The backend is chosen per machine: the API backend if the machine has an API power management config, otherwise Redfish when enabled and probed as available (cached per address), otherwise IPMI.

The COSI controllers live in `internal/provider/controllers` and their internal resources in `internal/provider/resources`.
The main ones are `MachineStatusController` (polls BMC power state), `BMCConfigurationController` (provisions IPMI creds via the agent), `PowerOperationController` (issues BMC power commands), `RebootStatusController` (reboots and one-time PXE), `WipeStatusController` (drives `WipeDisks`), and `InfraMachineStatusController` (syncs state up to Omni).

## Omni integration and state

The provider connects to Omni through the Omni client API, authenticating with the service account key created by `omnictl infraprovider create <id>` (an infra provider identity is a service account under the hood).
Through that client it obtains a COSI `state.State` backed by Omni over gRPC, with resources marshaled as protobuf (cosi-runtime supports state over gRPC, and the client wraps the gRPC state client into a regular COSI state).
The provider then builds its own COSI runtime on top of that remote state and registers its own resource types and controllers, following the same controller/reconciliation pattern as Omni itself.
Consequently the provider's internal store lives in Omni's state store (etcd), not locally, and the provider is stateless on disk.

An infra provider interacts with two namespaces.
`infra-provider` is the shared namespace through which Omni and providers communicate, and providers have limited access there.
`infra-provider:<id>` is the provider-specific namespace where the provider has full access, and this provider's internal resources live there (`<id>` defaults to `bare-metal` but is configurable).

Dynamic providers are usually built on the provider SDK in `omni/client/pkg/infra`, implementing a `Provisioner` and letting `infra.NewProvider` run the machine request provisioning loop.
This provider does not use that loop, since it provisions nothing.
It builds its own runtime and controllers and reuses only pieces of the SDK, such as the provider health status controller and the `infra` resource types.

## Key config and flags

Flags live in `cmd/provider/main.go`.
The most relevant for this repo's work:

- `--use-local-boot-assets` serves baked-in assets instead of the factory.
- `--agent-mode-talos-version` sets the Talos version used for factory agent-mode schematics, and it has no effect under `--use-local-boot-assets`.
- `--image-factory-base-url` and `--image-factory-pxe-base-url` point at the factory.
- `--secure-boot-enabled` serves a UKI, requires UEFI PXE mode, and rules out local boot assets.
- `--boot-from-disk-method` picks how an installed machine boots from disk (`ipxe-exit`, `http-404`, or `ipxe-sanboot`), for firmware that handles the iPXE exit path differently.
- The `--redfish-*` and `--ipmi-*` flags tune BMC behavior.
- `--agent-test-mode` boots the agent with API-based power management for QEMU test machines.
- The `--tls-*` flags choose between ephemeral auto-generated certs and persistent operator-supplied certs.

## Development and testing

Bare-metal machines can be emulated with QEMU for development and integration testing.
The `qemu-up` command in this repo uses the Talos provision library to create PXE-bootable QEMU machines with blank disks, the same machinery behind `talosctl cluster create`.
Each machine's launcher process serves a small per-machine HTTP power API (power on and off, reboot, one-time PXE boot, status) that emulates a BMC, and its address is recorded in the machine's launch config file under the provisioner state directory.
The provider is then run with `--agent-test-mode` and with `--api-power-mgmt-state-dir` pointing at that state directory.
In agent test mode the provider boots agents with the test-mode kernel argument, the agent reports API-based power management instead of configuring IPMI, and the provider looks up each machine's power API address in the state directory by node UUID.
From there on everything behaves as with real hardware, with the API-backed BMC client standing in for IPMI or Redfish.

`hack/test/integration.sh` wires this together end to end: it builds the provider image, brings up emulated machines with `qemu-up`, starts Vault and Omni in containers, creates the infra provider with `omnictl infraprovider create`, runs the provider in agent test mode, and then runs Omni's integration test suite against the whole stack.
It runs in CI through the `run-integration-test` make target.
Agent test mode is a development-only feature, so it is intentionally not part of the public documentation.

## Generated files and rekres

Many files are generated by kres and must not be hand-edited, among them `Dockerfile`, `Makefile`, `.github/workflows/*`, `.dockerignore`, `.gitignore`, `.golangci.yml`, and `lefthook.yml`.
The list is not exhaustive, and a kres-generated file identifies itself with a generated-file comment near its top.
To change them, edit `.kres.yaml` and run `make rekres`.
`make rekres` pulls `ghcr.io/siderolabs/kres:latest`, so a newer kres than last-generated brings tooling churn, which is expected in a bump-deps or rekres commit.

kres auto-detects tracked root-level `.md` files and wires them into the `Dockerfile` markdown-lint stage and `.dockerignore`.
Personal untracked files, for example a git-ignored `CLAUDE.local.md`, are left out.
That is why every committed root `.md` file, including this one, must pass markdownlint (see the markdown note at the top).

The boot-assets stage image is pinned in `.kres.yaml` as `ghcr.io/siderolabs/talos-metal-agent-boot-assets:<imager>-agent-<agent-version>`.
Repin it and rekres when you want newer default local-dev assets, and always use the upstream `siderolabs` image rather than a fork build.

`AgentModeTalosVersion` (`internal/provider/options.go`) should track a Talos version whose extensions catalog publishes the desired metal-agent version.

## Dependency bumps

Bump direct deps only: `go list -m -f '{{if and (not .Main) (not .Indirect)}}{{.Path}}@upgrade{{end}}' all | xargs go get`, then `go mod tidy`.
Never use `-u` or `all`, since those drag indirect deps past what direct deps require.

Bumping `talos/pkg/machinery` (Talos's public API) is almost always safe, including to an alpha, because it is backwards compatible.
When the latest `image-factory` or `omni/client` require a Talos prerelease, pin the prerelease first (`go get github.com/siderolabs/talos@<prerelease> github.com/siderolabs/talos/pkg/machinery@<prerelease>`), then run the direct-deps bump so it resolves cleanly.
A machinery jump can break a few call sites, so expect to fix them (for example the `image-factory` `SchematicCreate` return signature or renamed Talos `provision` options), and verify with `go build ./...` first.

Gates, in order, using the make targets, which are authoritative: `make lint-fmt`, `make lint` (includes govulncheck and markdownlint), `make generate`, and `make unit-tests`.

## Release process

Releases are cut with the repo's own tooling (`hack/release.sh` and `hack/release.toml`) in two phases.
First a release PR sets `hack/release.toml` `previous` to the latest tag, regenerates version files, runs `hack/release.sh changelog vX` to update the changelog, and creates the `release(vX): prepare release` commit via `hack/release.sh commit vX` (the script adds a DCO sign-off, and commits are GPG-signed when git is configured to sign).
Then, after merge, a signed tag is pushed, which triggers CI to build and push the release image and draft a GitHub release.
The deps bump and the release are separate PRs.

## Commit and PR conventions

Siderolabs PRs often contain a single commit, but multiple commits are fine for separate atomic logical changes (conform sets `maximumOfOneCommit: false`).
A single commit is usually still preferable, since the PR title and body then equal the commit title and body minus the DCO trailer.
Either way, keep PR titles and bodies simple like a commit message, without fancy markdown, and not long like a documentation page.
The commit title follows the conform rules: a `type[(scope)]: imperative summary` line, with types such as `feat`, `fix`, `chore`, `refactor`, `test`, `docs`, and `release`.
Commits carry a DCO `Signed-off-by` trailer (`git commit -s`) and are GPG-signed.
Reference issues inline as sentences, for example `Closes siderolabs/omni-infra-provider-bare-metal#<n>.`

## Related projects

- `siderolabs/booter` is a minimal spin-off carved out of this provider.
  It keeps only the DHCP proxy and the PXE boot serving, with no Omni dependency and no BMC management, to easily PXE-boot Talos machines in any subnet.
- `siderolabs/omni-infra-provider-proxmox` is a good reference for a dynamic provider built on the SDK.
- Sidero Metal (`siderolabs/sidero`) is the predecessor, kept in maintenance mode.

## Reference pointers

- The authoritative end-to-end guide is the Omni bare-metal infra provider tutorial at <https://docs.siderolabs.com/omni/omni-cluster-setup/setting-up-the-bare-metal-infrastructure-provider>.
- The Omni documentation also covers infrastructure providers (static versus dynamic) and manual bare-metal registration by PXE/iPXE or ISO, for contrast with this automated flow.
- The agent, its boot-assets, and the extensions coupling are documented in the `talos-metal-agent` repo's own `AGENTS.md`.
