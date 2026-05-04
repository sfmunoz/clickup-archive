# Agent Instructions

## Project Overview

`clickup-archive` is a Go CLI for exporting ClickUp workspace data into a local archive and browsing that archive with a terminal UI. The project is still early development, but the main fetch and TUI flows are implemented.

The CLI is built with Cobra. Fetching uses ClickUp API v2 and requires `CLICKUP_TOKEN`. Archive data is stored under:

```text
$HOME/.archive/clickup/
```

The archive root is expected to exist before fetch or TUI commands run.

The tool should not modify or delete ClickUp data. Treat the local archive as a snapshot, not a live sync target.

## Repository Layout

- `main.go`: CLI entrypoint.
- `cmd/`: Cobra commands.
  - `root.go`: root command and archive path helper.
  - `fetch*.go`: fetch command tree.
  - `tui.go`: TUI launch command.
- `internal/api/`: ClickUp API response/data types.
- `internal/fetch/`: API client and fetch orchestration.
- `internal/archive/`: local archive paths, save/load, JSON persistence, done markers, and attachment/comment archive handling.
- `internal/tui/`: Bubble Tea TUI models, searchable/expandable item list, content views, delegates, and stats.

## Common Commands

Run these from the repository root:

```sh
go test ./...
go build ./...
go run . --help
go run . fetch --help
go run . tui
```

Create the archive root when needed:

```sh
mkdir -p "$HOME/.archive/clickup"
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
- `fetch comments` and `fetch attachments` use per-task `.done` marker files and rebuild that task's comments or attachments directory when the marker is absent.
- `fetch attachments` re-fetches each task because the list endpoint omits attachment metadata, skips deleted attachments, saves metadata, and writes the downloaded file next to `index.json` when a download URL is available.
- Add focused tests for new archive path/persistence behavior and fetch orchestration that can be validated without live ClickUp credentials.
- Avoid tests that require `CLICKUP_TOKEN` by default. If live API checks are added, gate them behind an explicit environment variable and skip by default.
- For TUI changes, follow the existing Bubble Tea and Lip Gloss patterns in `internal/tui`. Keep the sidebar tree's filtering and `+`/`-` expand-collapse behavior working.

## Current Assumptions

- Go version is declared as `1.25.0` in `go.mod`.
- Primary dependencies include Cobra, Bubble Tea v2, Bubbles v2, Lip Gloss v2, and `github.com/sfmunoz/logit`.
- The implemented tree traversal is workspaces, spaces, folders, lists, tasks, then recursive subtasks.
- ClickUp lists directly under a space are not currently fetched; lists are loaded through folders.
- The TUI displays details for workspaces, spaces, folders, lists, tasks, and comments. Attachment metadata is loaded by the archive layer but is not currently shown in the TUI detail pane or stats overlay.

## Before Finishing

- Run `gofmt` on changed Go files.
- Run `go test ./...` when Go code changes.
- Run `go build ./...` for CLI-facing changes.
- Check `git status --short` and report any commands that could not be run.
