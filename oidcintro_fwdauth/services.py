"""Services module."""

import logging
from aiohttp_client_cache import CachedSession, CacheBackend
from aiohttp import hdrs
from yarl import URL


LOGGER = logging.getLogger(__name__)


class OIDCService:
    """Service class to communicate with an OIDC issuer."""

    DISCOVERY_PATH = URL("/.well-known/openid-configuration")

    def __init__(self, issuer_url, client_id, client_secret, *,
                 expire_after=None, introspection_endpoint=None):
        LOGGER.info("Initializing OIDC service: "
                    "issuer_url=%s, client_id=%s, expire_after=%d, introspection_endpoint=%s",
                    issuer_url, client_id, expire_after, introspection_endpoint)

        self.issuer_url = issuer_url
        self.client_id = client_id
        self.client_secret = client_secret
        self.discovery_endpoint = issuer_url.join(OIDCService.DISCOVERY_PATH)
        self.introspection_endpoint = introspection_endpoint

        expire_after = 0 if expire_after is None else expire_after
        cache = CacheBackend(allowed_methods=(hdrs.METH_GET, hdrs.METH_POST),
                             expire_after=expire_after)
        self.client_session = CachedSession(cache=cache, raise_for_status=True)

    async def discover(self):
        """Perform a discovery request to populate the internal configuration for this service.
           Running a discovery will override any attribute previously passed to the constructor."""
        LOGGER.info("Performing OIDC discovery: discovery_endpoint=%s", self.discovery_endpoint)

        async with self.client_session.get(self.discovery_endpoint) as response:
            data = await response.json()
            for key in ("introspection_endpoint",):
                if key in data:
                    setattr(self, key, data[key])
                    LOGGER.info("Discovered OIDC attribute: key=%s, value=%s", key, data[key])

    async def introspection(self, token, *, token_type_hint=None):
        """Perform a token introspection request."""
        if self.introspection_endpoint is None:
            raise RuntimeError("introspection endpoint not set, call discover() first")

        LOGGER.debug("Performing introspection: introspection_endpoint=%s, token_type_hint=%s",
                     self.introspection_endpoint, token_type_hint)

        data = {
            "client_id": self.client_id,
            "client_secret": self.client_secret,
            "token": token,
        }
        if token_type_hint:
            data["token_type_hint"] = token_type_hint

        async with self.client_session.post(self.introspection_endpoint, data=data) as response:
            return await response.json()

    async def close(self):
        """Close the HTTP client session of this service."""
        LOGGER.info("Closing OIDC service")
        if self.client_session:
            await self.client_session.close()
            LOGGER.info("HTTP client session closed")
