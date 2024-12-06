#!/bin/bash

set -eou pipefail

TALOS_VERSION=1.9.0-alpha.3
SUBNET_CIDR=172.42.0.0/24
GATEWAY_IP=172.42.0.1
ARTIFACTS=_out
JOIN_TOKEN=testonly
NUM_MACHINES=4

echo "OMNI_IMAGE: $OMNI_IMAGE"
echo "OMNI_INTEGRATION_TEST_IMAGE: $OMNI_INTEGRATION_TEST_IMAGE"
echo "SKIP_CLEANUP: $SKIP_CLEANUP"

docker pull "$OMNI_IMAGE"
docker pull "$OMNI_INTEGRATION_TEST_IMAGE"

echo "Build and push provider image to the temp registry $TEMP_REGISTRY..."

make image-provider REGISTRY="$TEMP_REGISTRY" TAG=test PUSH=true

PROVIDER_IMAGE="$TEMP_REGISTRY/siderolabs/omni-infra-provider-bare-metal:test"

echo "Download talosctl..."

mkdir -p ${ARTIFACTS}

[ -f ${ARTIFACTS}/talosctl ] || (crane export ghcr.io/siderolabs/talosctl:latest | tar x -C ${ARTIFACTS})

TALOSCTL=$(realpath "${ARTIFACTS}/talosctl")
QEMU_UP="${ARTIFACTS}/qemu-up-linux-amd64 --talosctl-path=${TALOSCTL} --cidr $SUBNET_CIDR --num-machines=$NUM_MACHINES"

echo "Register cleanup script..."

function cleanup() {
  local exit_code=$? # preserve the original exit code

  docker logs omni > /tmp/omni.log || true
  docker logs provider > /tmp/provider.log || true

  if [[ "$SKIP_CLEANUP" == "true" ]]; then
    echo "Skipping cleanup..."
    exit $exit_code
  fi

  pkill -f qemu-up-linux-amd64 || true

  # ${QEMU_UP} --destroy || true # disabled for now, as it removes Talos logs

  echo "Stop and remove Omni, Provider and Vault..."

  docker rm -f omni provider vault-dev || true
  rm -rf $ARTIFACTS/omni/ || true

  exit $exit_code
}

trap cleanup EXIT SIGINT

echo "Bring up some QEMU machines..."

${QEMU_UP} 2>&1 &

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

  for registry in docker.io k8s.gcr.io quay.io gcr.io ghcr.io registry.k8s.io factory.talos.dev; do
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
export BASE_URL="https://localhost:$OMNI_PORT/"
export AUTH_USERNAME="${AUTH0_TEST_USERNAME}"
export AUTH0_CLIENT_ID="${AUTH0_CLIENT_ID}"
export AUTH0_DOMAIN="${AUTH0_DOMAIN}"

docker run -d --network host \
  --name omni \
  -v ./hack/certs:/certs \
  -v "$(pwd)/${ARTIFACTS}/omni:/artifacts" \
  --cap-add=NET_ADMIN \
  --device=/dev/net/tun \
  -e SIDEROLINK_DEV_JOIN_TOKEN="${JOIN_TOKEN}" \
  -e VAULT_TOKEN=dev-o-token \
  -e VAULT_ADDR='http://127.0.0.1:8200' \
  "$OMNI_IMAGE" \
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
  "${REGISTRY_MIRROR_FLAGS[@]}"

docker logs -f omni &

echo "Wait for Omni to listen on ${BASE_URL}..."
timeout 60s bash -c "until curl -s -k -o /dev/null $BASE_URL; do echo 'Waiting for Omni...'; sleep 5; done"
echo "Omni is listening on ${BASE_URL}."

SERVICE_ACCOUNT_KEY_PATH="${ARTIFACTS}/omni/key"

echo "Wait for service account key to be created..."
timeout 60s bash -c "until [ -f '${SERVICE_ACCOUNT_KEY_PATH}' ]; do echo 'Waiting for service account key...'; sleep 5; done"
echo "Service account key is found at ${SERVICE_ACCOUNT_KEY_PATH}."

echo "Determine local IP address..."

LOCAL_IP=$(ip -o route get to 8.8.8.8 | sed -n 's/.*src \([0-9.]\+\).*/\1/p')

echo "Local IP: $LOCAL_IP"

# Export Omni endpoint and service account key
OMNI_ENDPOINT=${BASE_URL}
OMNI_SERVICE_ACCOUNT_KEY=$(cat $SERVICE_ACCOUNT_KEY_PATH)

export OMNI_ENDPOINT OMNI_SERVICE_ACCOUNT_KEY

# Launch infra provider in the background

echo "Launch infra provider..."

docker run -d --network host \
  --name provider \
  -v "$HOME/.talos/clusters/bare-metal:/api-power-mgmt-state:ro" \
  -e OMNI_ENDPOINT -e OMNI_SERVICE_ACCOUNT_KEY \
  "$PROVIDER_IMAGE" \
  --insecure-skip-tls-verify \
  --api-advertise-address="$LOCAL_IP" \
  --use-local-boot-assets \
  --agent-test-mode \
  --api-power-mgmt-state-dir=/api-power-mgmt-state \
  --dhcp-proxy-iface-or-ip=$GATEWAY_IP \
  --debug

docker logs -f provider &

echo "Waiting for provider to listen on $LOCAL_IP..."
timeout 60s bash -c "until curl -s -o /dev/null http://${LOCAL_IP}:50042; do echo 'Waiting for provider...'; sleep 5; done"
echo "Provider is listening on $LOCAL_IP."

echo "Run integration tests..."

docker run --rm --network host \
  --name omni-integration-test \
  -v "$(pwd)/hack/certs:/etc/ssl/certs" \
  -e SSL_CERT_DIR=/etc/ssl/certs \
  -e OMNI_SERVICE_ACCOUNT_KEY \
  "$OMNI_INTEGRATION_TEST_IMAGE" \
  --endpoint=${BASE_URL} \
  --expected-machines=$NUM_MACHINES \
  --talos-version="${TALOS_VERSION}" \
  --test.run "StaticInfraProvider"
