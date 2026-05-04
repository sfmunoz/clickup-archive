# Agent Instructions

## Project Overview

`clickup-archive` is a Go CLI for exporting ClickUp workspace data into a local archive and browsing that archive with a terminal UI. The project is early development and the README currently notes that it is not yet functional.

The CLI is built with Cobra. Fetching uses ClickUp API v2 and requires `CLICKUP_TOKEN`. Archive data is stored under:

```text
$HOME/.archive/clickup/
```

The tool should not modify or delete ClickUp data. Treat the local archive as a snapshot, not a live sync target.

## Repository Layout

- `main.go`: CLI entrypoint.
- `cmd/`: Cobra commands.
  - `root.go`: root command and archive path helper.
  - `fetch*.go`: fetch command tree.
  - `tui.go`: TUI launch command.
- `internal/api/`: ClickUp API response/data types.
- `internal/fetch/`: API client and fetch orchestration.
- `internal/archive/`: local archive paths, save/load, and JSON persistence.
- `internal/tui/`: Bubble Tea TUI models, views, delegates, and stats.

## Common Commands

Run these from the repository root:

```sh
go test ./...
go build ./...
go run . --help
go run . fetch --help
go run . tui
```

Use `gofmt` on edited Go files:

```sh
gofmt -w path/to/file.go
```

## Development Guidelines

- Prefer existing package boundaries: API structs in `internal/api`, fetch behavior in `internal/fetch`, archive persistence in `internal/archive`, and terminal UI code in `internal/tui`.
- Keep fetch logic read-only with respect to ClickUp. Do not add API calls that mutate remote ClickUp state unless the project requirements explicitly change.
- Preserve the archive location contract unless there is a deliberate migration plan.
- Use structured JSON encoding/decoding rather than string-building JSON.
- Keep filesystem writes inside the archive layer where practical.
- Add focused tests for new archive path/persistence behavior and fetch orchestration that can be validated without live ClickUp credentials.
- Avoid tests that require `CLICKUP_TOKEN` by default. If live API checks are added, gate them behind an explicit environment variable and skip by default.
- For TUI changes, follow the existing Bubble Tea and Lip Gloss patterns in `internal/tui`.

## Current Assumptions

- Go version is declared as `1.25.0` in `go.mod`.
- Primary dependencies include Cobra, Bubble Tea v2, Bubbles v2, Lip Gloss v2, and `github.com/sfmunoz/logit`.
- The traversal described by the CLI is workspaces, spaces, folders, lists, tasks, then recursive subtasks.

## Before Finishing

- Run `gofmt` on changed Go files.
- Run `go test ./...` when Go code changes.
- Run `go build ./...` for CLI-facing changes.
- Check `git status --short` and report any commands that could not be run.
