# clickup-archive

A CLI tool to export and archive your ClickUp workspace data locally.

## Purpose

ClickUp holds a lot of valuable information — tasks, comments, docs, attachments — but accessing it depends on an active subscription and an internet connection. **clickup-archive** fetches your workspace data via the ClickUp API and saves it to a local folder as plain files, giving you a portable, offline snapshot of your workspace that you own and control.

## What it does

- Connects to the ClickUp API using a personal API token
- Exports spaces, folders, lists, tasks, comments, and attachments
- Saves everything to a structured local directory
- Designed to be re-run to refresh the archive incrementally

## What it does NOT do

- It does not modify or delete any data in ClickUp
- It does not import data into any other tool
- The output folder is treated as **read-only** — it is a snapshot, not a live sync

## Use cases

- Keep a local backup of your ClickUp workspace
- Preserve data before cancelling or downgrading a subscription
- Browse and search your tasks and history without an internet connection
- Use the archive as a data source for custom tooling or reporting
