# Go v0 Init/Add/List Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Reboot TaskGraph as a Go CLI with `tg init`, `tg add`/`tg create`, and `tg list`, using `.taskgraph/tasks.md` as the v0 storage backend.

**Architecture:** Implement a small Go CLI with clear seams: command dispatch in `internal/cli`, project root discovery/init in `internal/project`, and task file IO in `internal/tasks`. `tg add` and `tg list` walk upward for `.taskgraph`; `tg add` auto-initializes in CWD if missing.

**Tech Stack:** Go (std lib), `testing` package, temp-dir based integration tests.

---

## Scope (approved)

- Commands: `tg init`, `tg add`, `tg create` (alias), `tg list`
- Storage: `.taskgraph/config.yml` and `.taskgraph/tasks.md`
- `tg add` appends one checklist line: `- [ ] <task text>`
- `tg add` auto-runs init in CWD if no `.taskgraph` found while walking up
- `tg list` prints raw checklist lines from `.taskgraph/tasks.md`

## Deferred (v0.2+)

- Richer `br create` parity flags (`--priority`, tags, due date, etc.)
- `br`-style richer list formatting/output controls
- Configurable task storage path outside `.taskgraph/tasks.md`

## Task 1: Initialize Go module and executable entrypoint

**Files:**
- Create: `go.mod`
- Create: `cmd/tg/main.go`

**Step 1: Create Go module file**

Set module path and Go version.

**Step 2: Create minimal CLI entrypoint**

`main.go` calls `internal/cli.Run(os.Args[1:], os.Stdout, os.Stderr)` and exits non-zero on error.

**Step 3: Commit bootstrap**

```bash
git add go.mod cmd/tg/main.go
git commit -m "chore(go): bootstrap tg CLI entrypoint"
```

## Task 2: Add failing tests for project discovery and init behavior

**Files:**
- Create: `internal/project/project_test.go`
- Create: `internal/project/project.go`

**Step 1: Write failing tests**

Cover:
- walk-up finds nearest ancestor `.taskgraph`
- missing `.taskgraph` returns not-found signal
- `InitAt(dir)` creates `.taskgraph/config.yml` and `.taskgraph/tasks.md`
- `InitAt(dir)` is idempotent

**Step 2: Run tests and verify failure**

Run: `go test ./internal/project -v`  
Expected: FAIL (functions not implemented)

**Step 3: Implement minimal `internal/project`**

Add:
- `FindTaskgraphRoot(startDir string) (root string, found bool, err error)`
- `InitAt(dir string) (root string, created bool, err error)`

**Step 4: Re-run tests**

Run: `go test ./internal/project -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/project/project.go internal/project/project_test.go
git commit -m "feat(go): add .taskgraph discovery and init primitives"
```

## Task 3: Add failing tests for task file append/read behavior

**Files:**
- Create: `internal/tasks/tasks_test.go`
- Create: `internal/tasks/tasks.go`

**Step 1: Write failing tests**

Cover:
- append task writes `- [ ] text` line
- append preserves existing lines and appends newline safely
- list returns raw checklist lines in file order

**Step 2: Run tests and verify failure**

Run: `go test ./internal/tasks -v`  
Expected: FAIL

**Step 3: Implement minimal tasks file functions**

Add:
- `AppendTask(tasksFile, text string) error`
- `ReadChecklistLines(tasksFile string) ([]string, error)`

**Step 4: Re-run tests**

Run: `go test ./internal/tasks -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/tasks/tasks.go internal/tasks/tasks_test.go
git commit -m "feat(go): implement tasks.md append and list primitives"
```

## Task 4: Add failing CLI integration tests for init/add/create/list

**Files:**
- Create: `internal/cli/cli_test.go`
- Create: `internal/cli/cli.go`

**Step 1: Write failing CLI integration tests**

Cover:
- `tg init` creates `.taskgraph/config.yml` and `.taskgraph/tasks.md`
- `tg add "x"` auto-inits in CWD when missing and reports initialization
- `tg add` uses nearest ancestor `.taskgraph` when present
- `tg create "x"` behaves exactly like `tg add "x"`
- `tg list` prints raw checklist lines
- `tg add` with empty text exits non-zero with usage/error message

**Step 2: Run tests and verify failure**

Run: `go test ./internal/cli -v`  
Expected: FAIL

**Step 3: Implement minimal `internal/cli.Run` command dispatch**

Implement:
- `init`
- `add`
- `create` (alias)
- `list`

Use `internal/project` and `internal/tasks`.

**Step 4: Re-run CLI tests**

Run: `go test ./internal/cli -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/cli.go internal/cli/cli_test.go
git commit -m "feat(go): add init/add/create/list CLI behavior"
```

## Task 5: Wire full test run and clean developer UX

**Files:**
- Modify: `README.md`
- Create: `Makefile` (optional; include only if useful for local test commands)

**Step 1: Add minimal developer run instructions**

Document:
- build: `go build ./cmd/tg`
- test: `go test ./...`
- run: `go run ./cmd/tg --help` (or equivalent command behavior)

**Step 2: Execute full test suite**

Run: `go test ./...`  
Expected: PASS

**Step 3: Commit**

```bash
git add README.md Makefile
git commit -m "docs(go): add local build and test instructions"
```

## Task 6: Remove Node/TypeScript implementation after Go parity is present

**Files:**
- Delete/modify as needed: `src/`, `tests/`, `package.json`, `pnpm-lock.yaml`, `tsconfig.json`, `vitest.config.ts`, `bin/taskgraph.mjs`

**Step 1: Remove old Node runtime and test scaffolding**

Keep docs that still provide historical context unless explicitly obsolete.

**Step 2: Verify repository still builds/tests with Go only**

Run:
- `go test ./...`
- `go build ./cmd/tg`

Expected: PASS

**Step 3: Commit**

```bash
git add -A
git commit -m "refactor: remove legacy Node CLI after Go reboot baseline"
```

## Verification checklist before merge

- `tg init` is idempotent
- `tg add` auto-inits in CWD when no `.taskgraph` exists
- `tg add` and `tg create` are equivalent
- `tg list` returns raw checklist lines from `.taskgraph/tasks.md`
- walk-up behavior selects nearest ancestor `.taskgraph`
- `go test ./...` passes cleanly
