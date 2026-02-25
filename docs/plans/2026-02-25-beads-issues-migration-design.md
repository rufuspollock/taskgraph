# Beads Issues Migration Design

**Date:** 2026-02-25

## Goal

Add a first-pass migration command that imports Beads issues from `./.beads/issues.jsonl` into TaskGraph checklist storage at `./.taskgraph/issues.md`.

## Scope

- In scope:
  - Read newline-delimited JSON from `./.beads/issues.jsonl`.
  - Append mapped checklist lines to `./.taskgraph/issues.md`.
  - Expose command as `tg migrate-beads`.
  - Require current working directory to contain both `./.beads/` and `./.taskgraph/`.
- Out of scope:
  - Config migration.
  - Idempotency/deduplication.
  - Dependency or label migration.
  - SQLite DB direct import.

## UX and Constraints

- Command runs relative to current working directory only.
- If either `./.beads/` or `./.taskgraph/` is missing, command fails with a clear expectation message.
- If `./.beads/issues.jsonl` is missing or unreadable, command fails with explicit file path.
- Duplicates are allowed on re-run (append-only behavior).

## Data Mapping

Each JSONL record is mapped as follows:

- Task text: `[beads:<id>] <title>`
- Checklist state:
  - `closed`, `done`, `resolved` -> `- [x]`
  - `tombstone` -> skipped
  - any other status -> `- [ ]`

Only records with non-empty `id` and `title` are imported. Invalid JSON lines fail the command with line number context.

## Architecture

Add a new package `internal/migrate` to keep migration logic separate from CLI dispatch and task creation logic:

- `ImportBeadsIssues(cwd string) (ImportSummary, error)`
- Helpers to validate required directories/files, parse JSONL stream, and append checklist lines.

CLI wiring:

- Add `migrate-beads` in `internal/cli/cli.go` command switch.
- Print summary of imported/skipped counts on success.

## Error Handling

- Missing required dirs: error states command expects both `./.beads` and `./.taskgraph` in current directory.
- JSON parse errors include line number.
- Missing `id` or `title` are counted as skipped invalid records and reported.
- IO errors from read/write bubble up with context.

## Testing Strategy

- Package tests for migration logic:
  - imports open/closed records correctly.
  - skips tombstones.
  - errors on missing required directories/files.
  - errors on malformed JSONL with line number.
- CLI tests:
  - `tg migrate-beads` happy path appends to `.taskgraph/issues.md`.
  - command emits expected missing-directory error.

## Risks

- Non-idempotent import may create duplicates when re-run (accepted for v1).
- Checklist text escaping is minimal; future versions may sanitize whitespace more strictly.
