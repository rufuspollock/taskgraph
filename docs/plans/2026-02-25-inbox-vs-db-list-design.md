# Inbox vs DB List Design

## Goal

Split queue and graph listing behavior for faster iteration:

- `tg inbox`: show only inbox capture file (`.taskgraph/issues.md`)
- `tg list`: show indexed checklist items from SQLite (`.taskgraph/taskgraph.db`)

## Decisions

- Breaking change is acceptable now (pre-adoption).
- `tg list` default includes only open checklist nodes.
- `tg list --all` includes closed checklist nodes.
- `tg list` sorts newest-first by source markdown file mtime, then by `path`, `line`.
- Keep implementation minimal and iterate quickly.

## Data model changes

Add `source_mtime_unix` to `index_nodes` so ordering can be done in SQL.

## Command semantics

- `tg add`: append to `.taskgraph/issues.md`, then rebuild index (already in place)
- `tg inbox`: read raw checklist lines from `.taskgraph/issues.md`
- `tg list`: query SQLite `index_nodes` for `kind='checklist'`

## Non-goals (this change)

- Advanced filters
- Node-kind mixing in `tg list`
- Rich formatting
