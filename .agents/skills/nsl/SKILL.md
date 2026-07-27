# nsl — Not So Localhost app registry CLI

Manage apps in the Not So Localhost stack through the registry server's HTTP API.

## Quick reference

```bash
# List all registered apps
nsl list

# Add an app (interactive — prompts for everything)
nsl add

# Add a frontend app (non-interactive)
nsl add --name my-app --type fe --target-url http://container-name:3000

# Add a backend app with Swagger docs
nsl add --name my-api --type be --docs-url http://api:8080/swagger

# Add a database app
nsl add --name order-db --type db --connection-string postgres://user:pass@host:5432/db

# Remove by ID or name prefix
nsl remove my-api

# Remove interactively (fuzzy select)
nsl remove
```

## CLI reference

### nsl list

Prints a table of all apps: abbreviated ID, name, type, route rule, enabled.

### nsl add

Dual-mode: no flags starts interactive prompts; any flag starts hybrid mode
(prompts only for missing required fields).

Flags:

| Flag                  | Env            | Description                     |
|-----------------------|----------------|---------------------------------|
| `--api-url`           | `NSL_API_URL`  | Registry API URL (default localhost:7272) |
| `--name, -n`          |                | App name — alphanumeric + hyphens |
| `--type, -t`          |                | App type: fe, be, db            |
| `--target-url, -u`    |                | HTTP URL to the app's frontend    |
| `--docs-url`          |                | OpenAPI doc URL (be only)       |
| `--connection-string` |                | postgres:// URL (db only)       |
| `--description, -d`   |                | Free-text description           |
| `--no-auth`           |                | Skip oauth2-proxy auth          |
| `--disabled`          |                | Register as disabled            |

**Type-specific requirements:**
- `fe` — `--target-url` is required
- `be` — `--docs-url` is required, `--target-url` is optional
- `db` — `--connection-string` is required

All type-specific fields are validated in real time during interactive mode.

### nsl remove

- `nsl remove <id-or-name>` — match by ID (full or prefix), then by name prefix
- `nsl remove` — interactive fuzzy select from live apps with confirmation

## Stack architecture

The Not So Localhost stack runs in Docker behind a Cloudflare Tunnel:

```
Cloudflare Tunnel -> Traefik:80
  ├── auth.YOUR_DOMAIN  -> keycloak:8080
  ├── t.YOUR_DOMAIN     -> oauth2-proxy -> host.docker.internal:7681 (ttyd)
  ├── apps.YOUR_DOMAIN  -> oauth2-proxy -> registry:7272
  └── *.YOUR_DOMAIN     -> oauth2-proxy -> registry:7272 (catchall)
```

Key services (`docker compose ps`):

| Container | Purpose |
|-----------|---------|
| `not-so-localhost-postgres-1` | Central PostgreSQL — app metadata + Keycloak + app DBs |
| `not-so-localhost-registry-1` | Registry HTTP API (port 7272) — nsl talks to this |
| `not-so-localhost-keycloak-1` | OIDC auth provider |
| `not-so-localhost-traefik-1` | Reverse proxy (port 80) |
| `not-so-localhost-oauth2-proxy-1` | Forward-auth middleware |
| `not-so-localhost-cloudflared-1` | Cloudflare Tunnel client |
| `not-so-localhost-terminal-1` | SSH terminal container |
| `not-so-localhost-backup-1` | DB backup runner |

### Sidecar containers (auto-deployed by registry)

When you add a `be` or `db` app, the registry deploys a sidecar container:

- **`<name>-swagger`** — Swagger UI for backend docs (image: `swaggerapi/swagger-ui`)
- **`<name>-pgweb`** — PGWeb for database browsing (image: `sosedoff/pgweb`)

Find them with `docker ps --filter network=not-so-localhost_edge`.

## Common workflows

### Add a frontend app

```bash
# 1. Spin up your app container on the edge network:
docker run -d --name my-app --network not-so-localhost_edge my-image

# 2. Register it in the registry:
nsl add --name my-app --type fe --target-url http://my-app:3000

# 3. Access at: https://apps.YOUR_DOMAIN/my-app
```

### Add a backend API

```bash
# The registry auto-deploys a Swagger UI sidecar.
nsl add --name my-api --type be --docs-url http://api-container:8080/swagger
# Swagger at: https://my-api.YOUR_DOMAIN
```

### Add a database

```bash
# Deploy your database container and register it:
nsl add --name order-db --type db --connection-string postgres://user:pass@host:5432/db
# PGWeb (DB browser) at: https://order-db.YOUR_DOMAIN
```

### Remove an app and its sidecars

```bash
nsl remove order-db
# Removes the PGWeb sidecar container, deletes the DB route, and removes
# the app from the registry.
```

### Check running containers

```bash
# All stack services:
docker compose ps

# Sidecars for registered apps:
docker ps --filter network=not-so-localhost_edge

# App containers by type:
docker ps --filter network=not-so-localhost_edge --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'
```

The registry writes Traefik routes to `traefik/dynamic/managed.yml` — this file is auto-generated from enabled apps and should not be edited by hand.
