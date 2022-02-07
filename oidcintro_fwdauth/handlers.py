"""Handlers module."""

from aiohttp import hdrs, web


async def auth_handler(request):
    """Handler for authentication requests."""
    authorization_headers = request.headers.getall(hdrs.AUTHORIZATION, [])
    if not authorization_headers:
        raise web.HTTPUnauthorized(text="This endpoint requires client authentication")

    client_ids = (set(request.query.get("client_ids").split(","))
                  if request.query.get("client_ids") else None)

    token_type_hint = request.query.get("token_type_hint", None)

    oidc_service = request.config_dict["oidc_service"]
    for authorization in authorization_headers:
        if " " not in authorization:
            continue
        scheme, parameters = authorization.split(" ", 1)
        if scheme.lower() != "bearer":
            continue

        token_info = await oidc_service.introspection(parameters, token_type_hint=token_type_hint)
        if "active" in token_info and not token_info["active"]:
            continue

        if client_ids and ("client_id" not in token_info or
                           token_info["client_id"] not in client_ids):
            continue

        response = web.Response()
        if "sub" in token_info:
            response.headers["X-Forwarded-User"] = token_info["sub"]
        if "client_id" in token_info:
            response.headers["X-Forwarded-Oidc-ClientId"] = token_info["client_id"]
        return response

    raise web.HTTPUnauthorized(text="Invalid client or client credentials")
