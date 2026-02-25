# Beads Issues Migration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement `tg migrate-beads` to import `./.beads/issues.jsonl` into `./.taskgraph/issues.md` using the approved status mapping and current-directory preconditions.

**Architecture:** Add an isolated migration module (`internal/migrate`) that validates paths, parses JSONL in source order, maps records to checklist lines, and appends to `.taskgraph/issues.md`. Wire the new command in CLI dispatch, keeping command behavior explicit and local to current working directory. Verify behavior through focused unit tests plus CLI integration-style tests.

**Tech Stack:** Go 1.22+, standard library (`bufio`, `encoding/json`, `os`, `path/filepath`), existing CLI/testing framework (`go test`).

---

### Task 1: Add failing tests for migration logic

**Files:**
- Create: `internal/migrate/migrate_test.go`

**Step 1: Write failing tests**

Add tests covering:
- Missing `.beads/` and/or `.taskgraph/` returns clear error.
- Missing `.beads/issues.jsonl` returns clear error.
- Valid JSONL imports open as unchecked and closed as checked.
- `tombstone` status is skipped.
- Malformed JSON line fails with line number.

**Step 2: Run test to verify failure**

Run: `go test ./internal/migrate -run Test -v`
Expected: FAIL (package missing or behavior unimplemented).

**Step 3: Commit**

```bash
git add internal/migrate/migrate_test.go
git commit -m "test(migrate): add failing tests for beads import"
```

### Task 2: Implement minimal migration module

**Files:**
- Create: `internal/migrate/migrate.go`
- Modify: `internal/migrate/migrate_test.go`

**Step 1: Write minimal implementation**

Implement:
- `ImportBeadsIssues(cwd string) (Summary, error)`
- Directory and file precondition checks.
- JSONL parsing line-by-line.
- Mapping rules:
  - `[beads:<id>] <title>`
  - closed/done/resolved -> checked
  - tombstone -> skip
  - others -> unchecked
- Append-only writing to `.taskgraph/issues.md` preserving source order.

**Step 2: Run package tests**

Run: `go test ./internal/migrate -v`
Expected: PASS.

**Step 3: Refactor minimally**

Extract tiny helpers only if needed for readability; keep behavior unchanged.

**Step 4: Commit**

```bash
git add internal/migrate/migrate.go internal/migrate/migrate_test.go
git commit -m "feat(migrate): import beads issues jsonl into taskgraph issues"
```

### Task 3: Add failing CLI tests for `migrate-beads`

**Files:**
- Modify: `internal/cli/cli_test.go`

**Step 1: Write failing tests**

Add tests for:
- `tg migrate-beads` happy path appends expected lines to `.taskgraph/issues.md`.
- `tg migrate-beads` errors when required directories are missing with explicit expectation text.

**Step 2: Run targeted tests to verify failure**

Run: `go test ./internal/cli -run MigrateBeads -v`
Expected: FAIL (command not wired yet).

**Step 3: Commit**

```bash
git add internal/cli/cli_test.go
git commit -m "test(cli): add migrate-beads command tests"
```

### Task 4: Wire CLI command and help output

**Files:**
- Modify: `internal/cli/cli.go`

**Step 1: Minimal command wiring**

- Add `migrate-beads` switch case.
- Implement `runMigrateBeads(stdout, stderr)` that:
  - gets effective cwd.
  - calls `migrate.ImportBeadsIssues(cwd)`.
  - prints imported/skipped summary.

**Step 2: Update help text**

Add `migrate-beads` to command list and examples.

**Step 3: Run targeted CLI tests**

Run: `go test ./internal/cli -run 'MigrateBeads|Help' -v`
Expected: PASS.

**Step 4: Commit**

```bash
git add internal/cli/cli.go
git commit -m "feat(cli): add migrate-beads command"
```

### Task 5: End-to-end verification and cleanup

**Files:**
- Modify (if needed): `README.md`

**Step 1: Decide docs update**

If CLI docs list commands, add `migrate-beads` usage note with current-directory requirement.

**Step 2: Run full verification**

Run: `go test ./...`
Expected: PASS.

**Step 3: Final commit**

```bash
git add README.md internal/cli/cli.go internal/cli/cli_test.go internal/migrate/migrate.go internal/migrate/migrate_test.go
git commit -m "feat: migrate beads issues jsonl into taskgraph issues"
```
