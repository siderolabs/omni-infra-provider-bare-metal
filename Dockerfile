# syntax = docker/dockerfile-upstream:1.16.0-labs

# THIS FILE WAS AUTOMATICALLY GENERATED, PLEASE DO NOT EDIT.
#
# Generated on 2025-06-25T07:26:01Z by kres 5128bc1.

ARG TOOLCHAIN

FROM ghcr.io/siderolabs/talos-metal-agent-boot-assets:v1.9.6-agent-v0.1.3 AS assets

FROM ghcr.io/siderolabs/ca-certificates:v1.10.0 AS image-ca-certificates

FROM ghcr.io/siderolabs/fhs:v1.10.0 AS image-fhs

FROM ghcr.io/siderolabs/ipxe:v1.11.0-alpha.0-25-g8c4603e AS ipxe

FROM --platform=linux/amd64 ghcr.io/siderolabs/ipxe:v1.11.0-alpha.0-25-g8c4603e AS ipxe-linux-amd64

FROM --platform=linux/arm64 ghcr.io/siderolabs/ipxe:v1.11.0-alpha.0-25-g8c4603e AS ipxe-linux-arm64

FROM ghcr.io/siderolabs/liblzma:v1.11.0-alpha.0-25-g8c4603e AS liblzma

# runs markdownlint
FROM docker.io/oven/bun:1.2.15-alpine AS lint-markdown
WORKDIR /src
RUN bun i markdownlint-cli@0.45.0 sentences-per-line@0.3.0
COPY .markdownlint.json .
COPY ./README.md ./README.md
RUN bunx markdownlint --ignore "CHANGELOG.md" --ignore "**/node_modules/**" --ignore '**/hack/chglog/**' --rules sentences-per-line .

FROM ghcr.io/siderolabs/musl:v1.11.0-alpha.0-25-g8c4603e AS musl

# collects proto specs
FROM scratch AS proto-specs
ADD api/specs/specs.proto /api/specs/

# base toolchain image
FROM --platform=${BUILDPLATFORM} ${TOOLCHAIN} AS toolchain
RUN apk --update --no-cache add bash curl build-base protoc protobuf-dev

# build tools
FROM --platform=${BUILDPLATFORM} toolchain AS tools
ENV GO111MODULE=on
ARG CGO_ENABLED
ENV CGO_ENABLED=${CGO_ENABLED}
ARG GOTOOLCHAIN
ENV GOTOOLCHAIN=${GOTOOLCHAIN}
ARG GOEXPERIMENT
ENV GOEXPERIMENT=${GOEXPERIMENT}
ENV GOPATH=/go
ARG GOIMPORTS_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg go install golang.org/x/tools/cmd/goimports@v${GOIMPORTS_VERSION}
RUN mv /go/bin/goimports /bin
ARG GOMOCK_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg go install go.uber.org/mock/mockgen@v${GOMOCK_VERSION}
RUN mv /go/bin/mockgen /bin
ARG PROTOBUF_GO_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOBUF_GO_VERSION}
RUN mv /go/bin/protoc-gen-go /bin
ARG GRPC_GO_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${GRPC_GO_VERSION}
RUN mv /go/bin/protoc-gen-go-grpc /bin
ARG GRPC_GATEWAY_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v${GRPC_GATEWAY_VERSION}
RUN mv /go/bin/protoc-gen-grpc-gateway /bin
ARG VTPROTOBUF_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@v${VTPROTOBUF_VERSION}
RUN mv /go/bin/protoc-gen-go-vtproto /bin
ARG DEEPCOPY_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg go install github.com/siderolabs/deep-copy@${DEEPCOPY_VERSION} \
	&& mv /go/bin/deep-copy /bin/deep-copy
ARG GOLANGCILINT_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCILINT_VERSION} \
	&& mv /go/bin/golangci-lint /bin/golangci-lint
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg go install golang.org/x/vuln/cmd/govulncheck@latest \
	&& mv /go/bin/govulncheck /bin/govulncheck
ARG GOFUMPT_VERSION
RUN go install mvdan.cc/gofumpt@${GOFUMPT_VERSION} \
	&& mv /go/bin/gofumpt /bin/gofumpt

# tools and sources
FROM tools AS base
WORKDIR /src
COPY go.mod go.mod
COPY go.sum go.sum
RUN cd .
RUN --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg go mod download
RUN --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg go mod verify
COPY ./api ./api
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg go list -mod=readonly all >/dev/null

FROM tools AS embed-generate
ARG SHA
ARG TAG
WORKDIR /src
RUN mkdir -p internal/version/data && \
    echo -n ${SHA} > internal/version/data/sha && \
    echo -n ${TAG} > internal/version/data/tag

# runs protobuf compiler
FROM tools AS proto-compile
COPY --from=proto-specs / /
RUN protoc -I/api --go_out=paths=source_relative:/api --go-grpc_out=paths=source_relative:/api --go-vtproto_out=paths=source_relative:/api --go-vtproto_opt=features=marshal+unmarshal+size+equal+clone /api/specs/specs.proto
RUN rm /api/specs/specs.proto
RUN goimports -w -local github.com/siderolabs/omni-infra-provider-bare-metal /api
RUN gofumpt -w /api

# runs gofumpt
FROM base AS lint-gofumpt
RUN FILES="$(gofumpt -l .)" && test -z "${FILES}" || (echo -e "Source code is not formatted with 'gofumpt -w .':\n${FILES}"; exit 1)

# runs golangci-lint
FROM base AS lint-golangci-lint
WORKDIR /src
COPY .golangci.yml .
ENV GOGC=50
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/root/.cache/golangci-lint,id=omni-infra-provider-bare-metal/root/.cache/golangci-lint,sharing=locked --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg golangci-lint run --config .golangci.yml

# runs govulncheck
FROM base AS lint-govulncheck
WORKDIR /src
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg govulncheck ./...

# runs unit-tests with race detector
FROM base AS unit-tests-race
WORKDIR /src
ARG TESTPKGS
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg --mount=type=cache,target=/tmp,id=omni-infra-provider-bare-metal/tmp CGO_ENABLED=1 go test -v -race -count 1 ${TESTPKGS}

# runs unit-tests
FROM base AS unit-tests-run
WORKDIR /src
ARG TESTPKGS
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg --mount=type=cache,target=/tmp,id=omni-infra-provider-bare-metal/tmp go test -v -covermode=atomic -coverprofile=coverage.txt -coverpkg=${TESTPKGS} -count 1 ${TESTPKGS}

FROM embed-generate AS embed-abbrev-generate
WORKDIR /src
ARG ABBREV_TAG
RUN echo -n 'undefined' > internal/version/data/sha && \
    echo -n ${ABBREV_TAG} > internal/version/data/tag

FROM scratch AS unit-tests
COPY --from=unit-tests-run /src/coverage.txt /coverage-unit-tests.txt

# cleaned up specs and compiled versions
FROM scratch AS generate
COPY --from=proto-compile /api/ /api/
COPY --from=embed-abbrev-generate /src/internal/version internal/version

# builds provider-linux-amd64
FROM base AS provider-linux-amd64-build
COPY --from=generate / /
COPY --from=embed-generate / /
WORKDIR /src/cmd/provider
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg GOARCH=amd64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=provider -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /provider-linux-amd64

# builds provider-linux-arm64
FROM base AS provider-linux-arm64-build
COPY --from=generate / /
COPY --from=embed-generate / /
WORKDIR /src/cmd/provider
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg GOARCH=arm64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=provider -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /provider-linux-arm64

# builds qemu-up-linux-amd64
FROM base AS qemu-up-linux-amd64-build
COPY --from=generate / /
COPY --from=embed-generate / /
WORKDIR /src/cmd/qemu-up
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg GOARCH=amd64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=qemu-up -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /qemu-up-linux-amd64

# builds qemu-up-linux-arm64
FROM base AS qemu-up-linux-arm64-build
COPY --from=generate / /
COPY --from=embed-generate / /
WORKDIR /src/cmd/qemu-up
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build,id=omni-infra-provider-bare-metal/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=omni-infra-provider-bare-metal/go/pkg GOARCH=arm64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=qemu-up -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /qemu-up-linux-arm64

FROM scratch AS provider-linux-amd64
COPY --from=provider-linux-amd64-build /provider-linux-amd64 /provider-linux-amd64

FROM scratch AS provider-linux-arm64
COPY --from=provider-linux-arm64-build /provider-linux-arm64 /provider-linux-arm64

FROM scratch AS qemu-up-linux-amd64
COPY --from=qemu-up-linux-amd64-build /qemu-up-linux-amd64 /qemu-up-linux-amd64

FROM scratch AS qemu-up-linux-arm64
COPY --from=qemu-up-linux-arm64-build /qemu-up-linux-arm64 /qemu-up-linux-arm64

FROM provider-linux-${TARGETARCH} AS provider

FROM scratch AS provider-all
COPY --from=provider-linux-amd64 / /
COPY --from=provider-linux-arm64 / /

FROM qemu-up-linux-${TARGETARCH} AS qemu-up

FROM scratch AS qemu-up-all
COPY --from=qemu-up-linux-amd64 / /
COPY --from=qemu-up-linux-arm64 / /

FROM scratch AS image-provider
ARG TARGETARCH
COPY --from=provider provider-linux-${TARGETARCH} /provider
COPY --from=image-fhs / /
COPY --from=image-ca-certificates / /
COPY --from=musl / /
COPY --from=liblzma / /
COPY --from=ipxe /usr/libexec/zbin /usr/bin/zbin
COPY --from=ipxe-linux-amd64 /usr/libexec/ /var/lib/ipxe/amd64
COPY --from=ipxe-linux-arm64 /usr/libexec/ /var/lib/ipxe/arm64
COPY --from=assets / /assets
LABEL org.opencontainers.image.source=https://github.com/siderolabs/omni-infra-provider-bare-metal
ENTRYPOINT ["/provider"]

