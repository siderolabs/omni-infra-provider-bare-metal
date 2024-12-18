# omni-infra-provider-bare-metal

This repo contains the code for the Omni Bare Metal Infra Provider.

## Requirements

To run the provider, you need:

- A running Omni instance
- An Omni infra provider service account matching the ID you'll use with this provider (`bare-metal` by default).
  To create it, run:

  ```bash
  omnictl serviceaccount create --use-user-role=false --role=InfraProvider infra-provider:bare-metal
  ```

  Replace `bare-metal` with your desired provider ID.
- A DHCP server: This provider runs a DHCP proxy to provide DHCP responses for iPXE boot, so a DHCP server must be running in the same network as the provider.
- Access to an [Image Factory](https://www.talos.dev/v1.8/learn-more/image-factory/).

## Development

For local development using Talos running on QEMU, follow these steps:

1. Set up a `buildx` builder instance with host network access, if you don't have one already:

   ```bash
   docker buildx create --driver docker-container --driver-opt network=host --name local1 --buildkitd-flags '--allow-insecure-entitlement security.insecure' --use
   ```

2. Start a local image registry if you don't have one running:

   ```bash
   docker run -d -p 5005:5000 --restart always --name local registry:2
   ```

3. Build `qemu-up` command line tool, and use it to start some QEMU machines:

   ```bash
   make qemu-up
   sudo -E _out/qemu-up-linux-amd64
   ```

4. (Optional) If you have made local changes to the [Talos Metal agent](https://github.com/siderolabs/talos-metal-agent), follow these steps to use your local version:
    1. Build and push Talos Metal Agent boot assets image following [these instructions](https://github.com/siderolabs/talos-metal-agent/blob/main/README.md).
    2. Replace the `ghcr.io/siderolabs/talos-metal-agent-boot-assets` image reference in [.kres.yaml](.kres.yaml) with your built image,
       e.g., `127.0.0.1:5005/siderolabs/talos-metal-agent-boot-assets:v1.9.0-agent-v0.1.0-beta.1-1-gbf1282b-dirty`.
    3. Re-kres the project to propagate this change into `Dockerfile`:

       ```bash
       make rekres
       ```

5. Build a local provider image:

   ```bash
   make image-provider PLATFORM=linux/amd64 REGISTRY=127.0.0.1:5005 PUSH=true TAG=local-dev
   docker pull 127.0.0.1:5005/siderolabs/omni-infra-provider-bare-metal:local-dev
   ```

6. Start the provider with your Omni API address and service account credentials:

   ```bash
   export OMNI_ENDPOINT=<your-omni-api-address>
   export OMNI_SERVICE_ACCOUNT_KEY=<your-omni-service-account-key>

   docker run --name=omni-bare-metal-provider --network host --rm -it \
     -v "$HOME/.talos/clusters/talos-default:/api-power-mgmt-state:ro" \
     -e OMNI_ENDPOINT -e OMNI_SERVICE_ACCOUNT_KEY \
     127.0.0.1:5005/siderolabs/omni-infra-provider-bare-metal:local-dev \
     --insecure-skip-tls-verify \
     --api-advertise-address=<provider-ip-to-advertise> \
     --use-local-boot-assets \
     --agent-test-mode \
     --api-power-mgmt-state-dir=/api-power-mgmt-state \
     --dhcp-proxy-iface-or-ip=172.42.0.1 \
     --debug
   ```

   Important flags:
    - `--use-local-boot-assets`: Makes the provider serve the boot assets image embedded in the provider image.
      This is useful for testing local Talos Metal Agent boot assets.
      Omit this flag to use the upstream agent version, which will forward agent mode PXE boot requests to the image factory.
    - `--agent-test-mode`: Boots the agent in test mode when booting a Talos node in agent mode, enabling API-based power management instead of IPMI/RedFish.
      This is necessary for QEMU development,
      as it uses the power management API run by the `talosctl cluster create` command.
    - The volume mount `-v "$HOME/.talos/clusters/talos-default:/api-power-mgmt-state:ro"`
      mounts the directory containing API-based power management state information generated by `talosctl cluster create`.
    - `--api-power-mgmt-state-dir`: Specifies where to read the API power management address of the nodes.
    - `--dhcp-proxy-iface-or-ip`: Specifies the IP address or interface name for running the DHCP proxy
      (e.g., the IP address of the QEMU bridge interface).
      The tool `qemu-up` uses the subnet `172.42.0.0/24` by default, and the bridge IP address on the host is `172.42.0.1`.

7. When you are done with the development/testing, destroy all QEMU machines and their network bridge:

   ```bash
   sudo -E _out/qemu-up-linux-amd64 --destroy
   ```
