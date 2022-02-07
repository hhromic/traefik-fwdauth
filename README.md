# traefik-oidcintro-fwdauth

Simple OIDC introspection forward auth service for [Traefik](https://github.com/traefik/traefik/)
(and possibly other compatible proxies), written in [Python 3](https://www.python.org/downloads/)
and [asyncio](https://docs.python.org/3/library/asyncio.html).

## Usage

Usage examples can be found in the [`examples`](examples/) directory.

## Building

To build a Docker image for the project:

    docker build -t traefik-oidcintro-fwdauth .

> **Note:** Ready-to-use images are available in the
> [GitHub Container Registry](https://github.com/users/hhromic/packages/container/package/traefik-oidcintro-fwdauth).

## License

This project is licensed under the [Apache License Version 2.0](LICENSE).
