# nsl

CLI for registering direct HTTP services across Not-So-Localhost machines.

## Install

The latest tagged release predates distributed node enrollment. Until a release
containing `enrollment-token` is tagged, install from a current checkout:

```sh
go install ./cmd/nsl
export PATH="$HOME/go/bin:$PATH"
nsl enrollment-token --help
```

The CLI uses the local registry at `http://localhost:7272` by default. Override
it with `NSL_API_URL` or `--api-url`.

Set `NSL_API_TOKEN` to the node's configured `REGISTRY_API_TOKEN` before
performing mutations. The CLI never prints this token.

## Register a service

```sh
nsl add --name landing --target-url http://host.docker.internal:7310
```

The registry assigns the local machine UUID and publishes:

```text
https://landing--<node-name>.joedodge.dev
```

Exposure policies:

```sh
# Keycloak browser login for the whole service
nsl add --name dashboard --target-url http://dashboard:3000 --policy browser

# The upstream validates its own API credentials
nsl add --name api --target-url http://api:8080 --policy upstream

# Keycloak for /ui and LiteLLM keys for /v1
nsl add --name litellm --target-url http://host.docker.internal:4000 --policy litellm
```

NSL does not start applications, databases, Swagger UI, or pgweb containers.
The application and its data remain owned by the machine running them.

## Inventory

```sh
nsl
nsl list
nsl nodes
nsl remove <id-or-exact-name>
```

Output uses TOON for compact, deterministic agent consumption.

## Enroll another machine

Prepare the new `not-so-localhost` checkout's ignored root, registry, and backup
environment files, then prebuild its images before issuing a short-lived token:

```sh
docker compose build
```

On an operator machine, load and export the broker values without printing
them, then issue the token:

```sh
source /path/to/not-so-localhost/.broker-secrets.env
export NSL_BROKER_URL NSL_BROKER_ADMIN_TOKEN
nsl enrollment-token --node-name laptop3
```

Set the returned value as `NSL_ENROLLMENT_TOKEN` in the new checkout's root
`.env`, confirm `NODE_NAME=laptop3` and `ENROLLMENT_BROKER_URL` use the same
broker, then start the already-built stack promptly:

```sh
docker compose up -d
```

The token is single-use and expires after at most 15 minutes.
