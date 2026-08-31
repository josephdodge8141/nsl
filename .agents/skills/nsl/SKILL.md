---
name: nsl
description: Register, list, verify, and remove HTTP services across Not-So-Localhost machines. Use whenever the user asks to register or expose an app, publish an x--node.joedodge.dev URL, configure LiteLLM access, inspect the shared apps portal, or enroll another NSL laptop.
---

# NSL expert

Execute NSL operations for the user instead of only suggesting commands.

Before mutations, ensure `NSL_API_TOKEN` is set from the node's local
`REGISTRY_API_TOKEN` configuration without printing either value.

## Register an app

Infer the app name and target from the current project when possible. A native
host process must use `host.docker.internal`, not `localhost`, because Traefik
runs in Docker.

```sh
nsl add --name <name> --target-url http://host.docker.internal:<port>
```

Choose one policy:

- `browser`: Keycloak protects every route. This is the default.
- `upstream`: the application validates its own API credentials.
- `litellm`: Keycloak protects browser/UI routes while `/v1` relies on LiteLLM
  virtual keys.

```sh
nsl add --name litellm \
  --target-url http://host.docker.internal:4000 \
  --policy litellm
```

Do not register databases. NSL has no Swagger UI or pgweb sidecars and does not
provision app workloads or app-owned databases.

After mutation, run `nsl list` and make a harmless request to the public URL.
Never print API keys used for verification.

## Shared inventory

`nsl list` returns every app from the shared S3 registry. Apps include an owning
node UUID; only that node renders their Traefik routes.

Use `nsl nodes` to map node UUIDs to names. CLI registration always uses the
local node. Use the authenticated browser portal for an explicitly selected
remote owner.

Repeated registration with the same node, name, and target is a successful
no-op. A conflicting target must be edited deliberately through the portal.

## Remove

Resolve the exact app first, then run:

```sh
nsl remove <id-or-exact-name>
```

When the same name exists on multiple nodes, use the UUID.

## Enroll a node

The enrollment broker must already be deployed. Read the admin secret from the
environment; never ask the user to paste it into chat.

The Cloudflare API token currently stored in the enrollment broker as
`CLOUDFLARE_API_TOKEN` expires on **2027-08-28**. Check or rotate this Worker
secret first when enrollment, tunnel creation, or DNS provisioning begins
failing near that date.

```sh
NSL_BROKER_ADMIN_TOKEN=<secret> \
  nsl enrollment-token --node-name <canonical-node-name>
```

The returned token is single-use, bound to that name, and valid for at most 15
minutes. Put it in the new node's root `.env` with `NODE_NAME` and
`ENROLLMENT_BROKER_URL`, then run
`docker compose up -d --build`. Node initialization exchanges it for persistent
node and tunnel credentials and never needs it again.

## Output

Commands emit TOON on stdout. Treat non-zero exit status as failure. Do not
parse stderr for application data.
