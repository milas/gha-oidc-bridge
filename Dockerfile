# syntax=docker/dockerfile:1
FROM --platform=${BUILDPLATFORM} golang:1.20 AS base

ARG TARGETOS
ARG TARGETARCH
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

WORKDIR /src

FROM base AS deps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM deps AS build
COPY cmd/ cmd/
COPY pkg/ pkg/

FROM build AS builder

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o /out/gha-oidc-bridge ./cmd

FROM scratch as gha-oidc-bridge-bin

COPY --from=builder --link --chmod=0755 /out/gha-oidc-bridge /

FROM gcr.io/distroless/base-debian11 as gha-oidc-bridge

COPY --from=gha-oidc-bridge-bin --link / /

ENTRYPOINT ["/gha-oidc-bridge"]
