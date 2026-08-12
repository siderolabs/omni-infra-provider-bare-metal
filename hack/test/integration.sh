#!/bin/bash

set -eou pipefail

TALOS_VERSION=1.13.0
SUBNET_CIDR=172.29.0.0/24
GATEWAY_IP=172.29.0.1
ARTIFACTS=_out
NUM_MACHINES=8
USE_LOCAL_BOOT_ASSETS=false
IMAGE_FACTORY_BASE_DOMAIN=factory.talos.dev
IMAGE_FACTORY_PXE_DOMAIN=pxe.factory.talos.dev

# The integration test drives the provider's real IPMI path against per-machine
# emulated BMCs (hosted in-process by qemu-up) instead of the fake HTTP power
# backend. Keep TALOSCTL_IMAGE in sync with the talos version in go.mod: the
# launcher comes from this image while qemu-up links the provision library, and a
# version mismatch fails silently because the launch config ignores unknown fields.
TALOSCTL_IMAGE=${TALOSCTL_IMAGE:-ghcr.io/siderolabs/talosctl:v1.14.0-rc.2-1-g322de8bf2}

echo "OMNI_IMAGE: $OMNI_IMAGE"
echo "OMNI_INTEGRATION_TEST_IMAGE: $OMNI_INTEGRATION_TEST_IMAGE"
echo "SKIP_CLEANUP: $SKIP_CLEANUP"

TEST_OUTPUTS_DIR=${GITHUB_WORKSPACE:-/tmp}/integration-test
mkdir -p "$TEST_OUTPUTS_DIR"

echo "Using test outputs dir: $TEST_OUTPUTS_DIR"

docker pull "$OMNI_IMAGE"
docker pull "$OMNI_INTEGRATION_TEST_IMAGE"

echo "Build and push provider image to the temp registry $TEMP_REGISTRY..."

make image-provider REGISTRY="$TEMP_REGISTRY" TAG=test PUSH=true

PROVIDER_IMAGE="$TEMP_REGISTRY/siderolabs/omni-infra-provider-bare-metal:test"

docker pull "$PROVIDER_IMAGE"

echo "Export talosctl (${TALOSCTL_IMAGE})..."

mkdir -p ${ARTIFACTS}

# Always re-export talosctl fresh, so a copy cached from a prior run cannot shadow
# a talos version bump made here.
TALOSCTL=$(realpath "${ARTIFACTS}/talosctl")
crane export "${TALOSCTL_IMAGE}" | tar x -C ${ARTIFACTS}

QEMU_UP="${ARTIFACTS}/qemu-up-linux-amd64 --talosctl-path=${TALOSCTL} --cidr $SUBNET_CIDR --num-machines=$NUM_MACHINES --virtual-bmc"

echo "Register cleanup script..."

# stop_supervisors signals any lingering qemu-up supervisor and waits for it to
# exit, so its launcher children are reparented and reaped before a destroy runs.
# A destroy polls the launcher PIDs, and a killed launcher whose qemu-up parent is
# still alive lingers as a zombie that the poll would wait on until it times out.
function stop_supervisors() {
  pkill -f qemu-up-linux-amd64 || true

  for _ in $(seq 1 20); do
    pgrep -f qemu-up-linux-amd64 >/dev/null || return 0
    sleep 0.5
  done

  # Still alive after the grace period: force it, so a stuck supervisor cannot keep
  # its launcher children parented and stall the destroy.
  pkill -KILL -f qemu-up-linux-amd64 || true
}

function cleanup() {
  local exit_code=$? # preserve the original exit code

  chown -R "${SUDO_USER:-$(whoami)}" ${ARTIFACTS} || true

  if [[ "$SKIP_CLEANUP" == "true" ]]; then
    echo "Skipping cleanup..."
    exit $exit_code
  fi

  rm -rf ./omnictl || true

  echo "Stop containers"
  docker stop omni provider vault-dev || true

  # Artifact collection is best-effort: it runs under `set -e`, and a failure here
  # (e.g., a container that never started) must not abort cleanup before qemu-up and
  # the machines are torn down below, which would leak the resident supervisor.
  echo "Gather container logs"
  docker logs omni &>"$TEST_OUTPUTS_DIR/omni.log" || true
  docker logs provider &>"$TEST_OUTPUTS_DIR/provider.log" || true

  echo "Gather machine logs and configs"
  machines_dir="$TEST_OUTPUTS_DIR/machines/"
  mkdir -p "$machines_dir"
  find "$HOME/.talos/clusters/bare-metal" -type f -name "*.log" ! -name "dhcpd.log" ! -name "lb.log" -exec cp {} "$machines_dir" \; || true
  find "$HOME/.talos/clusters/bare-metal" -type f -name "*.config" -exec cp {} "$machines_dir" \; || true

  # qemu-up hosts the emulated BMCs in-process, so stop it and WAIT for it to exit
  # before destroying: its launcher children must be reparented and reaped first,
  # otherwise a concurrent destroy polls zombie PIDs forever.
  if [[ -n "${QEMU_UP_PID:-}" ]]; then
    kill "$QEMU_UP_PID" 2>/dev/null || true
    wait "$QEMU_UP_PID" 2>/dev/null || true
  fi
  stop_supervisors
  ${QEMU_UP} --destroy || true
  pkill -f talosctl || true

  echo "Remove containers and Omni artifacts"
  docker rm -f omni provider vault-dev || true
  rm -rf $ARTIFACTS/omni/ || true

  exit $exit_code
}

# Run cleanup once, from the EXIT trap only. The signal traps just exit, which
# then fires EXIT, so a SIGINT/SIGTERM does not run cleanup a second time on top of
# the normal exit path (which truncated logs and raced container removal).
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

echo "Stop any existing QEMU machines..."

stop_supervisors
${QEMU_UP} --destroy || true
pkill -f talosctl || true

echo "Bring up some QEMU machines (qemu-up stays running to host the emulated BMCs)..."

${QEMU_UP} >"$TEST_OUTPUTS_DIR/qemu-up.log" 2>&1 &
QEMU_UP_PID=$!

echo "Wait for qemu-up to report all machines and virtual BMCs ready..."
timeout 300s bash -c "until grep -q 'all machines and virtual BMCs ready' '$TEST_OUTPUTS_DIR/qemu-up.log'; do
  kill -0 $QEMU_UP_PID 2>/dev/null || { echo 'qemu-up exited before becoming ready:'; cat '$TEST_OUTPUTS_DIR/qemu-up.log'; exit 1; }
  sleep 2
done"

echo "Wait for IP address $GATEWAY_IP to appear..."
timeout 60s bash -c "until ip a | grep -q '${GATEWAY_IP}'; do echo 'Waiting for IP address...'; sleep 5; done"
echo "IP address $GATEWAY_IP is up."

echo "Start Vault..."

docker run --rm -d --cap-add=IPC_LOCK -p 8200:8200 -e 'VAULT_DEV_ROOT_TOKEN_ID=dev-o-token' --name vault-dev hashicorp/vault:1.15

sleep 10

echo "Load private key into Vault..."

docker cp hack/certs/key.private vault-dev:/tmp/key.private
docker exec -e VAULT_ADDR='http://0.0.0.0:8200' -e VAULT_TOKEN=dev-o-token vault-dev \
  vault kv put -mount=secret omni-private-key \
  private-key=@/tmp/key.private

sleep 5

echo "Build registry mirror args..."

if [[ "${CI:-false}" == "true" ]]; then
  REGISTRY_MIRROR_FLAGS=()

  for registry in docker.io k8s.gcr.io quay.io gcr.io ghcr.io registry.k8s.io $IMAGE_FACTORY_BASE_DOMAIN; do
    service="registry-${registry//./-}.ci.svc"
    addr=$(python3 -c "import socket; print(socket.gethostbyname('${service}'))")

    REGISTRY_MIRROR_FLAGS+=("--registry-mirror=${registry}=http://${addr}:5000")
  done
else
  # use the value from the environment, if present
  REGISTRY_MIRROR_FLAGS=("${REGISTRY_MIRROR_FLAGS:-}")
fi

echo "Launch Omni..."

export OMNI_PORT=8099
export BASE_URL="https://localhost:$OMNI_PORT"
export AUTH_USERNAME="${AUTH0_TEST_USERNAME}"
export AUTH0_CLIENT_ID="${AUTH0_CLIENT_ID}"
export AUTH0_DOMAIN="${AUTH0_DOMAIN}"

sqlite_path="${TEST_OUTPUTS_DIR}/sqlite.db"

docker run -d --network host \
  --name omni \
  -v ./hack/certs:/certs \
  -v "$(pwd)/${ARTIFACTS}/omni:/artifacts" \
  --cap-add=NET_ADMIN \
  --device=/dev/net/tun \
  -e SIDEROLINK_DEV_JOIN_TOKEN=testonly \
  -e VAULT_TOKEN=dev-o-token \
  -e VAULT_ADDR='http://127.0.0.1:8200' \
  "$OMNI_IMAGE" \
  --eula-accept-name="Test User" \
  --eula-accept-email="test-user@siderolabs.com" \
  --siderolink-wireguard-advertised-addr=${GATEWAY_IP}:50180 \
  --siderolink-api-advertised-url="grpc://${GATEWAY_IP}:8090" \
  --machine-api-bind-addr=0.0.0.0:8090 \
  --siderolink-wireguard-bind-addr=0.0.0.0:50180 \
  --event-sink-port=8091 \
  --auth-auth0-enabled=true \
  --advertised-api-url="${BASE_URL}" \
  --auth-auth0-client-id="${AUTH0_CLIENT_ID}" \
  --auth-auth0-domain="${AUTH0_DOMAIN}" \
  --initial-users="${AUTH_USERNAME}" \
  --private-key-source="vault://secret/omni-private-key" \
  --public-key-files="/certs/key.public" \
  --bind-addr="0.0.0.0:$OMNI_PORT" \
  --enable-talos-pre-release-versions \
  --key=/certs/localhost-key.pem \
  --cert=/certs/localhost.pem \
  --etcd-embedded-unsafe-fsync=true \
  --embedded-discovery-service-snapshots-enabled=false \
  --create-initial-service-account \
  --initial-service-account-key-path=/artifacts/key \
  --join-tokens-mode=strict \
  --image-factory-address="https://${IMAGE_FACTORY_BASE_DOMAIN}" \
  --image-factory-pxe-address="https://${IMAGE_FACTORY_PXE_DOMAIN}" \
  --sqlite-storage-path="${sqlite_path}" \
  "${REGISTRY_MIRROR_FLAGS[@]}"

docker logs -f omni &

echo "Wait for Omni to listen on ${BASE_URL}..."
timeout 60s bash -c "until curl -s -k -o /dev/null $BASE_URL; do echo 'Waiting for Omni...'; sleep 5; done"
echo "Omni is listening on ${BASE_URL}."

ADMIN_SERVICE_ACCOUNT_KEY_PATH="${ARTIFACTS}/omni/key"

echo "Wait for service account key to be created..."
timeout 60s bash -c "until [ -f '${ADMIN_SERVICE_ACCOUNT_KEY_PATH}' ]; do echo 'Waiting for admin service account key...'; sleep 5; done"
echo "Admin service account key is found at ${ADMIN_SERVICE_ACCOUNT_KEY_PATH}."

ADMIN_SERVICE_ACCOUNT_KEY=$(cat "$ADMIN_SERVICE_ACCOUNT_KEY_PATH")

export OMNI_SERVICE_ACCOUNT_KEY="${ADMIN_SERVICE_ACCOUNT_KEY}"
export OMNI_ENDPOINT="${BASE_URL}"

echo "Download omnictl..."
curl -fsSL -k -o ./omnictl "${BASE_URL}/api/omnictl/omnictl-linux-amd64"
chmod +x ./omnictl

echo "Create infra provider..."

PROVIDER_SERVICE_ACCOUNT_KEY=$(./omnictl --insecure-skip-tls-verify infraprovider create bare-metal | grep 'OMNI_SERVICE_ACCOUNT_KEY=' | cut -d'=' -f2-)

# Get image factory leaf certificate PEM
openssl s_client -showcerts -connect $IMAGE_FACTORY_PXE_DOMAIN:443 </dev/null | awk '/BEGIN CERTIFICATE/,/END CERTIFICATE/ { print; if (/END CERTIFICATE/) exit }' >factory.crt

echo "Launch infra provider in the background..."

# We run the provider in a container, as its container image contains everything needed by the provider,
# e.g., ipmitool and ipxe binaries, metal agent boot assets etc.
docker run -d --network host \
  --name provider \
  -v ./factory.crt:/factory.crt:ro \
  -e OMNI_ENDPOINT \
  -e OMNI_SERVICE_ACCOUNT_KEY="${PROVIDER_SERVICE_ACCOUNT_KEY}" \
  "$PROVIDER_IMAGE" \
  --insecure-skip-tls-verify \
  --api-advertise-address="$GATEWAY_IP" \
  --use-local-boot-assets=$USE_LOCAL_BOOT_ASSETS \
  --redfish-use-when-available=false \
  --ipmi-pxe-boot-mode=bios \
  --min-reboot-interval=1m \
  --machine-labels=a=b,c \
  --tls-custom-ipxe-ca-cert-file=/factory.crt \
  --debug

docker logs -f provider &

echo "Waiting for provider to listen on $GATEWAY_IP..."
timeout 60s bash -c "until curl -s -o /dev/null http://$GATEWAY_IP:50042; do echo 'Waiting for provider...'; sleep 5; done"
echo "Provider is listening on $GATEWAY_IP."

echo "Run integration tests..."

docker run --rm --network host \
  --name omni-integration-test \
  -v "$(pwd)/hack/certs:/etc/ssl/certs" \
  -v "$(pwd)/hack/test:/var/test" \
  -v "$TEST_OUTPUTS_DIR:/tmp/integration-test/" \
  -e SSL_CERT_DIR=/etc/ssl/certs \
  -e OMNI_SERVICE_ACCOUNT_KEY="$ADMIN_SERVICE_ACCOUNT_KEY" \
  "$OMNI_INTEGRATION_TEST_IMAGE" \
  --omni.endpoint=${BASE_URL} \
  --omni.talos-version="${TALOS_VERSION}" \
  --omni.provision-config-file=/var/test/provisionconfig.yaml \
  --omni.skip-extensions-check-on-create \
  --test.failfast \
  --test.v \
  --test.run "TestIntegration/Suites/(StaticInfraProvider|ConfigPatching)$"
