# clickup-archive

A Go CLI for exporting ClickUp workspace data into a local archive and browsing
that archive in a terminal UI.

The project is still in early development, but the main archive flow is
implemented:

- fetch the ClickUp hierarchy with ClickUp API v2
- fetch task comments from the archived task tree
- download task attachments from the archived task tree
- browse the local archive with a Bubble Tea terminal UI, including searchable
  and expandable/collapsible archive items

The tool is read-only with respect to ClickUp. It does not modify or delete
remote ClickUp data.

## Archive Location

Archive data is written under:

```text
$HOME/.archive/clickup/
```

The archive root must already exist before running fetch or TUI commands:

```sh
mkdir -p "$HOME/.archive/clickup"
```

Each ClickUp entity is stored in its own directory with an `index.json` file.
The hierarchy is:

```text
$HOME/.archive/clickup/
  <workspace-id>/index.json
    <space-id>/index.json
      <folder-id>/index.json
        <list-id>/index.json
          <task-id>/index.json
            comments/<comment-id>/index.json
            attachments/<attachment-id>/index.json
            attachments/<attachment-id>/<downloaded-file>
```

Subtasks are fetched recursively and stored alongside the other tasks in the
same list directory. The current tree fetch walks lists inside folders; ClickUp
folderless lists are not fetched yet.

## Requirements

- Go 1.25.0 or newer
- A ClickUp personal API token
- `CLICKUP_TOKEN` set in the environment for fetch commands

```sh
export CLICKUP_TOKEN="your-clickup-token"
```

## Usage

Build the CLI:

```sh
go build ./...
```

Show command help:

```sh
go run . --help
go run . fetch --help
```

Fetch the workspace tree:

```sh
go run . fetch tree
```

Fetch comments for every archived task:

```sh
go run . fetch comments
```

Download attachments for every archived task:

```sh
go run . fetch attachments
```

Launch the terminal UI:

```sh
go run . tui
```

## Fetch Behavior

`fetch tree` walks:

```text
workspaces -> spaces -> folders -> lists -> tasks -> subtasks
```

Tasks are fetched with `subtasks=true` and paginated until exhausted.

`fetch comments` walks the local task tree and saves each task comment under
`<task-id>/comments/<comment-id>/index.json`. After a task's comments are
successfully fetched, the command writes `<task-id>/comments.done`. Existing
done markers cause the task to be skipped. If the marker is missing, the task's
comments directory is rebuilt.

`fetch attachments` walks the local task tree, re-fetches each task to obtain
attachment metadata, writes metadata under
`<task-id>/attachments/<attachment-id>/index.json`, and downloads the file next
to that metadata. After a task's attachments are successfully fetched, the
command writes `<task-id>/attachments.done`. Existing done markers cause the
task to be skipped. If the marker is missing, the task's attachments directory
is rebuilt.

## Terminal UI

The `tui` command loads the local archive and opens an interactive browser.
The sidebar shows the archive tree, the content pane shows details for the
selected workspace, space, folder, list, task, or comment, and `F1` opens an
archive statistics overlay.

Current controls include:

- `q` or `ctrl-c`: quit
- `q`: close the stats overlay when it is visible
- `F1`: show or hide archive stats
- `+`: expand the selected archive item, or expand the next collapsed level
  below it
- `-`: collapse the deepest expanded level below the selected archive item, or
  collapse the selected item

The item list also supports the standard Bubble Tea list navigation and filter
behavior. Filtering searches each item's title and description.

## Current Limitations

- The archive format is still evolving.
- Fetch commands depend on ClickUp API v2 responses and have not been hardened
  as a fully resumable backup system.
- Fetching currently skips ClickUp lists that are directly under a space rather
  than inside a folder.
- The TUI is functional but basic, and it loads attachment metadata without
  displaying attachments in the detail pane or stats overlay.
- Docs and other ClickUp object types outside the implemented workspace, space,
  folder, list, task, comment, and attachment paths are not archived.
- Tests should avoid requiring live ClickUp credentials by default.

## Development

Common commands:

```sh
go test ./...
go build ./...
go run . --help
go run . fetch --help
go run . tui
```

Format edited Go files with:

```sh
gofmt -w path/to/file.go
```
