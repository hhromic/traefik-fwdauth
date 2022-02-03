"""Application signals module."""

from aiohttp import hdrs
from .services import OIDCService


async def oidc_service_cleanup_ctx(app):
    """Cleanup context for the OIDC service."""
    try:
        app["oidc_service"] = OIDCService(app["oidc_issuer_url"],
                                          app["oidc_client_id"],
                                          app["oidc_client_secret"],
                                          expire_after=app["oidc_expire_after"])
        await app["oidc_service"].discover()
        yield
    finally:
        await app["oidc_service"].close()


async def remove_server_header(_, response):
    """Signal for removing the server header from a response."""
    if hdrs.SERVER in response.headers:
        del response.headers[hdrs.SERVER]
