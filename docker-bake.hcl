variable "TAG" {
  default = "latest"
}

target "default" {
  tags = ["ghcr.io/hhromic/traefik-fwdauth:${TAG}"]
}

target "snapshot" {
  inherits = ["default"]
  args = {
    GORELEASER_EXTRA_ARGS = "--snapshot"
  }
  tags = ["ghcr.io/hhromic/traefik-fwdauth:snapshot"]
}
