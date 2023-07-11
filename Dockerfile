# Start a new stage for building the application
FROM golang:1.20.6-bullseye AS builder

# Install GoReleaser
RUN go install github.com/goreleaser/goreleaser@v1.19.2

# Set a well-known building directory
WORKDIR /build

# Download and verify application dependencies
COPY go.mod go.sum ./
RUN go mod download \
    && go mod verify

# Copy application sources and build the application
COPY . .
ENV CGO_ENABLED=0
RUN goreleaser build --clean --single-target --output traefik-fwdauth

# Start a new stage for the final application image
FROM gcr.io/distroless/static-debian11:nonroot AS final

# Configure image labels
LABEL org.opencontainers.image.source=https://github.com/hhromic/traefik-fwdauth \
      org.opencontainers.image.description="Simple Forward Auth service for Traefik (and possibly other compatible proxies), written in Go." \
      org.opencontainers.image.licenses=Apache-2.0

# Configure default entrypoint and exposed port of the application
ENTRYPOINT ["/traefik-fwdauth"]
EXPOSE 9878

# Copy application binary
COPY --from=builder /build/traefik-fwdauth /traefik-fwdauth
