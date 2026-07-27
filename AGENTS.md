# Project agent memory

This file is committed project-intrinsic agent knowledge: build, test, release, architecture, and sharp-edge notes that should travel with the code.

- Build: `go build -ldflags="-X main.Version=$VERSION" ./cmd/nsl/`
- Install: `go install github.com/josephdodge8141/nsl/cmd/nsl@v0.1.0`
- Lint: `go vet ./...` and `gofmt -d .`
- Uses stdlib `flag` for CLI, `survey` for interactive prompts
- API server default: `http://localhost:7272`, configurable via `--api-url` or `NSL_API_URL`
- CLI calls `${api}/api/v1/...` routes. Breaking server changes bump to `/api/v2/`.
- Version embedded at build time via `-ldflags="-X main.Version=..."`. Default: `"dev"`.
- Registry server version endpoint: `GET /api/v1/version` returns `{"version":"..."}`.

## Skills

- `.agents/skills/nsl/SKILL.md` — CLI reference, stack architecture, common workflows

## Registry server (not-so-localhost)

The `nsl` CLI talks to the registry server's HTTP API. The server lives in the separate `not-so-localhost` repo and handles:
- CRUD on the `apps` PostgreSQL table
- Traefik route generation (`traefik/dynamic/managed.yml`)
- Docker sidecar deployment (swagger-ui for `be` apps, pgweb for `db` apps)

The registry exposes port 7272 on the host. No auth on the API — it's localhost-only.
