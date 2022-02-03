"""Simple OIDC introspection forward auth service entry-point module."""

import asyncio
import logging
from argparse import ArgumentParser, ArgumentDefaultsHelpFormatter
from pathlib import Path
import uvloop
from yarl import URL
from aiohttp import web
from .handlers import auth_handler
from .signals import oidc_service_cleanup_ctx, remove_server_header
from .version import __version__


DEF_HOST = "0.0.0.0"
DEF_PORT = 4181
DEF_REMOVE_SERVER_HEADER = False
DEF_OIDC_EXPIRE_AFTER = 300
DEF_LOG_LEVEL = "INFO"
DEF_ACCESS_LOG_LEVEL = "WARN"

LOGGER = logging.getLogger(__package__)
LOGGER_FORMAT = "%(asctime)s [%(name)s] %(levelname)s %(message)s"


def main(args):
    """Main entry-point."""
    LOGGER.info("Starting application: version=%s", __version__)

    oidc_client_secret = (Path(args.oidc_client_secret_file).read_text(encoding="utf-8").strip()
                          if args.oidc_client_secret_file else args.oidc_client_secret)
    if not oidc_client_secret:
        raise RuntimeError("no OIDC client secret was provided")

    app = web.Application()
    app["oidc_issuer_url"] = args.oidc_issuer_url
    app["oidc_client_id"] = args.oidc_client_id
    app["oidc_client_secret"] = oidc_client_secret
    app["oidc_expire_after"] = args.oidc_expire_after
    app.add_routes((
        web.route("*", "/auth", auth_handler),
    ))
    app.cleanup_ctx.append(oidc_service_cleanup_ctx)
    if args.remove_server_header:
        app.on_response_prepare.append(remove_server_header)

    LOGGER.info("Running HTTP server on %s:%d", args.host, args.port)
    web.run_app(app, host=args.host, port=args.port, print=None)


if __name__ == "__main__":
    PARSER = ArgumentParser(prog=__package__, description=__doc__,
                            formatter_class=ArgumentDefaultsHelpFormatter)

    HTTP_ARGS = PARSER.add_argument_group("main HTTP server arguments")
    HTTP_ARGS.add_argument("--host", metavar="HOSTNAME",
                           default=DEF_HOST,
                           help="HTTP server listening host")
    HTTP_ARGS.add_argument("--port", metavar="PORT",
                           default=DEF_PORT, type=int,
                           help="HTTP server listening port")
    HTTP_ARGS.add_argument("--remove-server-header",
                           default=DEF_REMOVE_SERVER_HEADER, action="store_true",
                           help="whether to remove or leave the 'Server' header in responses")

    OIDC_ARGS = PARSER.add_argument_group("OIDC service arguments")
    OIDC_ARGS.add_argument("--oidc-issuer-url", metavar="URL",
                           required=True, type=URL,
                           help="OIDC Issuer URL")
    OIDC_ARGS.add_argument("--oidc-client-id", metavar="ID",
                           required=True,
                           help="OIDC Client Id")
    OIDC_ARGS.add_argument("--oidc-client-secret", metavar="SECRET",
                           help="OIDC Client Secret")
    OIDC_ARGS.add_argument("--oidc-client-secret-file",
                           metavar="FILE",
                           help="File with the OIDC Client Secret")
    OIDC_ARGS.add_argument("--oidc-expire-after", metavar="SECONDS",
                           default=DEF_OIDC_EXPIRE_AFTER, type=int,
                           help="time for expiring cached OIDC requests")

    MONITORING_ARGS = PARSER.add_argument_group("application monitoring arguments")
    MONITORING_ARGS.add_argument("--log-level", metavar="LEVEL",
                                 default=DEF_LOG_LEVEL,
                                 help="application logging level")
    MONITORING_ARGS.add_argument("--access-log-level", metavar="LEVEL",
                                 default=DEF_ACCESS_LOG_LEVEL,
                                 help="HTTP server access logging level")

    ARGS = PARSER.parse_args()

    logging.basicConfig(format=LOGGER_FORMAT, level=ARGS.log_level)
    logging.getLogger("aiohttp.access").setLevel(ARGS.access_log_level)

    asyncio.set_event_loop_policy(uvloop.EventLoopPolicy())
    main(ARGS)
