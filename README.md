# Traefik Forward Auth Service

Simple [Forward Auth](https://doc.traefik.io/traefik/middlewares/http/forwardauth/) service for
[Traefik](https://github.com/traefik/traefik/) (and possibly other compatible proxies), written in
[Go](https://go.dev/).

This Forward Auth service implements the following features:

* [OAuth2 Token Introspection](https://datatracker.ietf.org/doc/html/rfc7662) validation.
* Introspection endpoint discovery via [OpenID Connect Discovery](https://openid.net/specs/openid-connect-discovery-1_0.html).

## Usage

Usage examples can be found in the [`examples/`](examples/) directory.

## Building

To build a release Docker image for the project:
```
git checkout vX.Y.Z
docker buildx build -t ghcr.io/hhromic/traefik-fwdauth:vX.Y.Z .
```

> **Note:** Ready-to-use images are available in the
> [GitHub Container Registry](https://github.com/users/hhromic/packages/container/package/traefik-fwdauth).

To build a snapshot locally Using [GoReleaser](https://goreleaser.com/):
```
goreleaser build --clean --single-target --snapshot
```

## License

This project is licensed under the [Apache License Version 2.0](LICENSE).
