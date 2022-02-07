# Traefik Example

> **Note:** The example in this directory assumes that you have Docker running
> in [swarm mode](https://docs.docker.com/engine/swarm/).

To run the example, first edit the Traefik stack in `stacks/traefik.yaml` to
customise your OIDC Issuer URL and Client ID. Also edit the secret in
`secrets/oidc-client-secret` to customise your OIDC Client Secret.

Then, deploy the Traefik stack:

    docker secret create oidc-client-secret secrets/oidc-client-secret
    docker stack deploy -c stacks/traefik.yaml traefik

The above will initialise a Traefik instance running on ports `5555` (default
entrypoint for serving requests) and `8080` (the Traefik dashboard).

In addition, this stack also creates two overlay networks: `traefik_default`,
for the `dsproxy`, `fwdauth` and `traefik` services, and `traefik_services`, for
the application services to proxy.

At this point, you can browse the running [Traefik dashboard](http://localhost:8080).

If you navigate to the [HTTP Services](http://localhost:8080/dashboard/#/http/services)
status page, you should be able to find a running `fwdauth@docker` service.

If you navigate to the [HTTP Middlewares](http://localhost:8080/dashboard/#/http/middlewares)
status page, you should be able to find three middlewares: `oidcintro-auth@docker`,
`oidcintro-cl1-auth@docker` and `oidcintro-cl12-auth@docker`. These correspond
to the different authentication configurations defined in the Traefik stack file.

Now, to deploy the example application service stack:

    docker stack deploy -c stacks/app.yaml app

This stack will deploy six [whoami](https://github.com/traefik/whoami)
containers:

* Two using the `oidcintro-auth@docker` ForwardAuth middleware.
* Two using the `oidcintro-cl1-auth@docker` ForwardAuth middleware.
* Two using the `oidcintro-cl12-auth@docker` ForwardAuth middleware.

After deploying, and after a short time, the configured Docker Swarm
autodiscovery in Traefik will find the deployed service containers and
autoconfigure routing/middlewares using the labels.

> **Note:** The Docker Swarm autodiscovery defaults to refreshing data every `15s`.

If you navigate to the [HTTP Services](http://localhost:8080/dashboard/#/http/services)
status page, you should now be able to find three new application services:
`whoami@docker`, `whoami-cl1@docker` and `whoami-cl12@docker`.

If you navigate to the [HTTP Routers](http://localhost:8080/dashboard/#/http/routers)
status page, you should now be able to find three new routers corresponding to
the three deployed application services. You can further navigate into each
router to see their details.

All of these services should be protected by OIDC introspection:
```
$ curl -H "Authorization: Bearer YOUR-TOKEN" http://localhost:5555/whoami
$ curl -H "Authorization: Bearer YOUR-TOKEN" http://localhost:5555/whoami-cl1
$ curl -H "Authorization: Bearer YOUR-TOKEN" http://localhost:5555/whoami-cl12
```

> **Note:** In this example, the Traefik access logs (available on stderr) will
> contain the authenticated user and issuing Client ID via logging of the
> `X-Forwarded-User` and `X-Forwarded-Oidc-ClientId` HTTP headers.

Finally, to remove everything done in this example:

    docker stack rm app
    docker stack rm traefik
    docker secret rm oidc-client-secret
