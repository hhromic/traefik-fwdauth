variable "TAG" {
  default = "latest"
}

target "default" {
  tags = ["ghcr.io/hhromic/traefik-fwdauth:${TAG}"]
  annotations = [
    "org.opencontainers.image.source=https://github.com/hhromic/traefik-fwdauth",
    "org.opencontainers.image.description=Simple Forward Auth service for Traefik (and possibly other compatible proxies), written in Go.",
    "org.opencontainers.image.licenses=Apache-2.0",
  ]
}

target "snapshot" {
  inherits = ["default"]
  args = {
    GORELEASER_EXTRA_ARGS = "--snapshot"
  }
  tags = ["ghcr.io/hhromic/traefik-fwdauth:snapshot"]
}
