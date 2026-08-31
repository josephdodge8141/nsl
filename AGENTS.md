# Project agent memory

This file is committed project-intrinsic agent knowledge: build, test, release, architecture, and sharp-edge notes that should travel with the code.

- Build: `go build -ldflags="-X main.Version=$VERSION" ./cmd/nsl/`
- Install: `go install github.com/josephdodge8141/nsl/cmd/nsl@v0.1.0`
- Lint: `go vet ./...` and `gofmt -d .`
- Uses stdlib `flag`; every workflow is non-interactive and agent-safe
- API server default: `http://localhost:7272`, configurable via `--api-url` or `NSL_API_URL`
- CLI calls `${api}/api/v2/...` routes.
- Version embedded at build time via `-ldflags="-X main.Version=..."`. Default: `"dev"`.
- Registry server version endpoint: `GET /api/v2/version` returns `{"version":"..."}`.

## Skills

- `.agents/skills/nsl/SKILL.md` — CLI reference, stack architecture, common workflows

## Registry server (not-so-localhost)

The `nsl` CLI talks to the registry server's HTTP API. The server lives in the separate `not-so-localhost` repo and handles:
- S3 CAS updates to the shared node/app registry
- assignment of apps to persistent node UUIDs
- per-node Traefik route reconciliation
- Cloudflare DNS provisioning through the enrollment broker

The registry binds port 7272 to host loopback for the local CLI.
