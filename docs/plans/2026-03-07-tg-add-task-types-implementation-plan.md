# tg add Task Types Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add single-type support to `tg add`/`tg create` using `--type`, store types as markdown `#t-<type>` tags, and validate against built-ins plus optional project custom types.

**Architecture:** Keep markdown checklist lines in `.taskgraph/issues.md` authoritative. Extend add argument parsing to handle one normalized type, validate against `project` config-backed allowed types, and encode type as namespaced label (`t-...`) merged into existing label handling. Reuse existing index label pipeline so type metadata remains derived from markdown.

**Tech Stack:** Go, stdlib, existing `internal/cli`, `internal/tasks`, `internal/project`

---

### Task 1: Add failing tests for type parsing and storage in tasks package

**Files:**
- Modify: `internal/tasks/tasks_test.go`
- Modify: `internal/tasks/tasks.go`

**Step 1: Write the failing tests**

Add tests for:

- Extracting one type from `#t-*` tokens.
- Rejecting multiple inline type labels.
- Appending a normalized type label to `AppendTask` output.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/tasks -run 'Type|AppendTask' -v`
Expected: FAIL because type helpers and enforcement do not exist.

**Step 3: Write minimal implementation**

Implement helpers in `internal/tasks/tasks.go`:

- normalize type values
- map type to namespaced label (`t-...`)
- extract inline type label from text
- detect duplicates/conflicts

Update `AppendTask` signature/logic to accept a single optional type and enforce one-type rule.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/tasks -run 'Type|AppendTask' -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/tasks/tasks.go internal/tasks/tasks_test.go
git commit -m "feat: add typed task helpers and storage"
```

### Task 2: Add failing tests for CLI add/create --type behavior

**Files:**
- Modify: `internal/cli/cli_test.go`
- Modify: `internal/cli/cli.go`
- Modify: `internal/project/project.go`
- Modify: `internal/project/project_test.go`

**Step 1: Write the failing tests**

Add CLI tests for:

- `tg add "x" --type epic` stores `#t-epic`.
- `tg create "x" --type idea` also works.
- repeated `--type` fails.
- conflicting inline `#t-*` and flag type fails.
- unknown type fails.
- configured custom type via `issue-types` passes.

Add project config tests for parsing `issue-types` from `.taskgraph/config.yml`.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli ./internal/project -run 'Type|issue-types|Add|Create' -v`
Expected: FAIL due to missing parser and config support.

**Step 3: Write minimal implementation**

- Extend add arg parser with `-t/--type` single-value handling.
- Add allowed-type loading in `project` package (built-ins + config custom).
- Validate and pass selected type to `tasks.AppendTask`.
- Update usage/help text for add/create examples.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli ./internal/project -run 'Type|issue-types|Add|Create' -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/cli/cli.go internal/cli/cli_test.go internal/project/project.go internal/project/project_test.go
git commit -m "feat: support tg add --type with allowlisted validation"
```

### Task 3: Full verification and docs update

**Files:**
- Modify: `README.md`
- Modify: `internal/cli/cli.go`
- Modify: `docs/plans/2026-03-07-tg-add-task-types-design.md` (if final wording tweak needed)

**Step 1: Add/update docs**

- Add `--type` usage examples in README.
- Ensure `tg --help` text documents `--type` and one-type rule.

**Step 2: Run full verification**

Run: `go test ./...`
Expected: PASS.

**Step 3: Spot-check CLI behavior manually**

Run:

- `go test ./internal/cli -run TestAddSupportsTypeFlag -v`
- `go test ./internal/cli -run TestAddRejectsUnknownType -v`

Expected: PASS.

**Step 4: Commit**

```bash
git add README.md internal/cli/cli.go docs/plans/2026-03-07-tg-add-task-types-design.md docs/plans/2026-03-07-tg-add-task-types-implementation-plan.md
git commit -m "docs: document tg add task types"
```
