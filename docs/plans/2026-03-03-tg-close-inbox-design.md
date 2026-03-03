# tg close Inbox-Only Design

**Goal:** Add a first-pass `tg close` command that closes an inbox task in `.taskgraph/issues.md` by ID, checks its markdown checkbox, and appends a dated completion reason note.

## Scope

This first pass only supports closing tasks from `.taskgraph/issues.md`.

The command shape is:

```bash
tg close <id> "<reason>"
```

Example:

```text
before: - [ ] ➕2026-03-03 [tg-abc] call Alice #home
after:  - [x] ➕2026-03-03 [tg-abc] call Alice #home **✅2026-03-03 done on phone**
```

## Why Inbox-Only

`tg` currently creates inbox IDs in `.taskgraph/issues.md`, and `tg inbox` already treats that file as the source of truth.

Searching arbitrary markdown files by ID would introduce ambiguity and a fragile write path before the index/database stores stable source locations for write-back.

Using SQLite for close operations is the likely future direction, but it should wait until the index is explicitly trusted for write location resolution.

## Command Behavior

`tg close` should:

- require an item ID and a non-empty reason
- locate the current `.taskgraph` root from the current working directory
- open `.taskgraph/issues.md`
- find the first checklist line containing the exact bracketed ID, e.g. `[tg-abc]`
- require the matching task to be open (`- [ ] `)
- rewrite that line to closed (`- [x] `)
- append a completion annotation in this format: `**✅YYYY-MM-DD reason**`
- rebuild the SQLite index after the file update
- print a concise confirmation

## Error Handling

The command should return clear errors for:

- missing ID
- missing reason
- missing `.taskgraph`
- missing `.taskgraph/issues.md`
- ID not found in inbox
- task already closed

Usage text should be:

```text
usage: tg close <id> <reason>
```

## Markdown Rewrite Rules

The rewrite should preserve the original task line content as much as possible:

- keep the original task text
- keep the original ID
- keep labels and other trailing text already on the line
- only change the checkbox marker from `[ ]` to `[x]`
- append the completion note at the end of the line

For the first pass, the command can assume the matched task line does not already contain a completion annotation. If it does, the command should still reject the task as already closed based on the checkbox state alone.

## Architecture

The close logic should live in `internal/tasks` because it is markdown file mutation logic, analogous to `AppendTask` and `ReadChecklistLines`.

The CLI layer in `internal/cli` should:

- parse `close` arguments
- locate the taskgraph root
- call the task mutation helper
- rebuild the index
- print success or usage errors

This keeps markdown parsing and rewriting out of the CLI package and avoids duplicating file handling logic later if more task mutations are added.

## Testing

Tests should cover:

- parsing and successful close via CLI
- failure on missing reason
- failure on unknown ID
- failure on already closed item
- markdown rewrite format including checkbox flip and completion note
- index rebuild after close so the task no longer appears as open in indexed queries

Tests should use today’s date dynamically rather than hard-coding a calendar date.

## Deferred Follow-Up

After this lands, add a new inbox item to track a richer future close flow:

- close tasks outside inbox using indexed location data and an explicit write-back model

That follow-up should depend on a design that defines how SQLite location data is trusted and how write conflicts are handled.
