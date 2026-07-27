# Project agent memory

This file is the project's committed home for project-intrinsic agent knowledge: build, test, release, architecture, and sharp-edge notes that should travel with the code.

- Add durable project-specific notes here as they are discovered through real work.
- Build: `go build ./cmd/nsl/`
- Lint: `go vet ./...` and `gofmt -d .`
- Uses stdlib `flag` for CLI, `survey` for interactive prompts
- API server default: `http://localhost:7272`, configurable via `--api-url` or `NSL_API_URL`
