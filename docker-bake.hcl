variable "TAG" {
  default = "latest"
}

variable "IMAGE_CREATED" {
  default = timestamp()
  type = string
}

variable "_metadata" {
  default = {
    "org.opencontainers.image.authors" = "https://github.com/hhromic"
    "org.opencontainers.image.created" = IMAGE_CREATED
    "org.opencontainers.image.description" = "Simple Forward Auth service for Traefik (and possibly other compatible proxies), written in Go."
    "org.opencontainers.image.documentation" = "https://github.com/hhromic/traefik-fwdauth"
    "org.opencontainers.image.licenses" = "Apache-2.0"
    "org.opencontainers.image.source" = "https://github.com/hhromic/traefik-fwdauth"
    "org.opencontainers.image.title" = "Traefik Forward Auth Service"
    "org.opencontainers.image.url" = "https://github.com/hhromic/traefik-fwdauth"
    "org.opencontainers.image.vendor" = "https://github.com/hhromic"
  }
  type = map(string)
}

target "default" {
  tags = ["ghcr.io/hhromic/traefik-fwdauth:${TAG}"]
  annotations = [for k, v in _metadata: "${k}=${v}"]
  labels = _metadata
  attest = [
    {"type" = "provenance", "disabled" = true},
    {"type" = "sbom", "disabled" = true},
  ]
}

target "snapshot" {
  inherits = ["default"]
  args = {
    GORELEASER_EXTRA_ARGS = "--snapshot"
  }
  tags = ["ghcr.io/hhromic/traefik-fwdauth:snapshot"]
}
