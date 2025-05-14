# Start a new stage for building the application
FROM golang:1.24.3 AS builder

# Download and install GoReleaser
ADD https://github.com/goreleaser/goreleaser/releases/download/v2.9.0/goreleaser_Linux_x86_64.tar.gz goreleaser.tar.gz
RUN tar zxf goreleaser.tar.gz --one-top-level=/usr/bin goreleaser

# Set a well-known building directory
WORKDIR /build

# Download and verify application dependencies
COPY go.mod go.sum ./
RUN go mod download \
    && go mod verify

# Copy application sources and build the application
COPY . .
ARG GORELEASER_EXTRA_ARGS
RUN GOTOOLCHAIN=local \
    CGO_ENABLED=0 \
    goreleaser build --clean --single-target ${GORELEASER_EXTRA_ARGS}

# Start a new stage for the final application image
FROM cgr.dev/chainguard/static:latest AS final

# Configure image labels
LABEL org.opencontainers.image.source=https://github.com/hhromic/traefik-fwdauth \
      org.opencontainers.image.description="Simple Forward Auth service for Traefik (and possibly other compatible proxies), written in Go." \
      org.opencontainers.image.licenses=Apache-2.0

# Configure default entrypoint and exposed port of the application
ENTRYPOINT ["/traefik-fwdauth"]
EXPOSE 9878

# Copy application binary
COPY --from=builder /build/dist/traefik-fwdauth /traefik-fwdauth
