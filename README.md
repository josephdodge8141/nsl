# nsl

CLI for managing apps in the Not So Localhost registry server.

## Install

    go install ./cmd/nsl

## Usage

The CLI communicates with the registry server at `http://localhost:7272` by default.
Set a custom URL with `--api-url` or the `NSL_API_URL` environment variable.

### List all apps

    nsl list
    nsl list --api-url http://registry.example.com:7272

### Add an app

Interactive mode (no flags):

    nsl add

Non-interactive:

    nsl add --name my-api --type be --docs-url http://my-api:8080/swagger
    nsl add --name fe-app --type fe --target-url http://fe-app:3000 --description "Frontend"
    nsl add --name pg-db --type db --connection-string postgres://user:pass@host:5432/db

Partial flags (prompts for missing fields):

    nsl add --name my-api --type be
    nsl add --name my-api --disabled

Flags:

| Flag                  | Env            | Description                     |
|-----------------------|----------------|---------------------------------|
| `--api-url`           | `NSL_API_URL`  | Registry API URL                |
| `--name, -n`          |                | App name                        |
| `--type, -t`          |                | App type (fe, be, db)           |
| `--target-url, -u`    |                | Target URL (fe/be)              |
| `--docs-url`          |                | Docs URL (be)                   |
| `--connection-string` |                | Postgres connection string (db) |
| `--description, -d`   |                | Description                     |
| `--no-auth`           |                | Disable auth (default false)    |
| `--disabled`          |                | Create disabled (default false) |

### Remove an app

By ID or name:

    nsl remove my-api
    nsl remove abc12345

Interactive fuzzy select (omit argument):

    nsl remove
