# SQLite Index Command Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add `tg index` to build a local SQLite index of markdown file/heading/checklist nodes at `.taskgraph/taskgraph.db`.

**Architecture:** Extend CLI with a new `index` command that resolves taskgraph root and invokes a new index package. The index package scans markdown files, parses headings/checklist lines, and rebuilds a single SQLite table transactionally.

**Tech Stack:** Go stdlib + `modernc.org/sqlite` driver

---

### Task 1: Add failing CLI tests for `tg index`

**Files:**
- Modify: `internal/cli/cli_test.go`

**Step 1: Write the failing tests**
- Add a test that runs `tg index` in an initialized repo and expects success summary text.
- Add a test that verifies DB file `.taskgraph/taskgraph.db` exists after indexing.

**Step 2: Run tests to verify failure**
Run: `go test ./internal/cli -run Index -count=1`
Expected: FAIL due to unknown command.

**Step 3: Implement minimal CLI dispatch**
- Add `index` command branch in CLI dispatch and help text.

**Step 4: Re-run tests**
Run: `go test ./internal/cli -run Index -count=1`
Expected: still FAIL until index backend exists.

### Task 2: Add failing tests for index scanner/parser behavior

**Files:**
- Create: `internal/indexer/indexer_test.go`
- Create: `internal/indexer/testdata/sample.md`

**Step 1: Write failing parser tests**
- Test extraction of `file`, `heading`, `checklist` nodes and checklist states.
- Test path exclusion rules (dot folders + node_modules).

**Step 2: Run tests to verify failure**
Run: `go test ./internal/indexer -count=1`
Expected: FAIL (missing package functionality).

### Task 3: Implement markdown scanning/parsing to in-memory nodes

**Files:**
- Create: `internal/indexer/indexer.go`

**Step 1: Implement minimal parser/scanner**
- Recursive markdown file discovery with exclusion rules.
- Explicit inclusion of `.taskgraph/tasks.md`.
- Build node hierarchy and stable IDs.

**Step 2: Re-run tests**
Run: `go test ./internal/indexer -count=1`
Expected: PASS.

### Task 4: Add SQLite rebuild tests and implementation

**Files:**
- Create: `internal/indexer/sqlite_test.go`
- Create: `internal/indexer/sqlite.go`
- Modify: `go.mod`

**Step 1: Write failing SQLite tests**
- Rebuild writes all nodes into `index_nodes`.
- Rebuild clears previous rows before insert.

**Step 2: Run tests to verify failure**
Run: `go test ./internal/indexer -count=1`
Expected: FAIL.

**Step 3: Implement SQLite schema + rebuild transaction**
- `CREATE TABLE/INDEX IF NOT EXISTS`
- `DELETE` + insert in transaction

**Step 4: Re-run tests**
Run: `go test ./internal/indexer -count=1`
Expected: PASS.

### Task 5: Wire CLI index command end-to-end

**Files:**
- Modify: `internal/cli/cli.go`
- Modify: `internal/cli/cli_test.go`

**Step 1: Implement `runIndex`**
- Discover root
- Build index nodes
- Rebuild sqlite DB
- Print summary

**Step 2: Run focused tests**
Run: `go test ./internal/cli -count=1`
Expected: PASS.

### Task 6: Verify full project

**Files:**
- No file changes

**Step 1: Run full test suite**
Run: `go test ./...`
Expected: PASS.
