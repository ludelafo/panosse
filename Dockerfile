## Documentation
# - https://docs.docker.com/build/guide/multi-platform/
# - https://docs.docker.com/build/dockerfile/frontend/
# syntax=docker/dockerfile:1

# Default versions
ARG GO_VERSION=1.21
ARG ALPINE_VERSION=3.20

## Base stage
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base

WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download -x

## Build stage
FROM base AS build

ARG TARGETOS
ARG TARGETARCH

# Copy and build the source code
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,target=. \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/panosse ./main.go

## Runtime stage
FROM alpine:${ALPINE_VERSION} AS runtime

# Install required packages
RUN apk add --no-cache flac

# Copy the binary from the build stage
COPY --from=build /bin/panosse /bin/

# Default entrypoint
ENTRYPOINT [ "panosse" ]

# Default command passed to the entrypoint
CMD [ "--help" ]
