# Label Support Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add basic markdown-native label support to `tg` so tasks can store labels as `#tags`, `tg add` can append them via `--labels`, and `tg inbox` / `tg list` can filter by labels using `br`-style `--label` semantics.

**Architecture:** Markdown checklist lines in `.taskgraph/issues.md` remain the source of truth. Label parsing is shared between task writing, inbox filtering, and indexing, while SQLite stores only derived label data for fast `tg list` queries and can be rebuilt from markdown without loss.

**Tech Stack:** Go, stdlib, modernc SQLite, existing `taskgraph` CLI/indexer packages

---

### Task 1: Add label parsing helpers in tasks package

**Files:**
- Modify: `internal/tasks/tasks.go`
- Test: `internal/tasks/tasks_test.go`

**Step 1: Write the failing tests**

Add tests covering:

```go
func TestNormalizeLabelsFromCSV(t *testing.T)
func TestExtractLabelsFromText(t *testing.T)
func TestMergeTextAndFlagLabelsDeduplicates(t *testing.T)
```
```go
func TestExtractLabelsFromText(t *testing.T) {
    got := ExtractLabels("prep venue notes #flowershow #ABC")
    want := []string{"flowershow", "abc"}
    if !reflect.DeepEqual(got, want) {
        t.Fatalf("got %v want %v", got, want)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/tasks -run 'TestNormalizeLabelsFromCSV|TestExtractLabelsFromText|TestMergeTextAndFlagLabelsDeduplicates' -v`
Expected: FAIL with undefined helper errors.

**Step 3: Write minimal implementation**

Add small helpers in `internal/tasks/tasks.go`:

- Parse comma-separated labels from CLI input.
- Extract inline `#labels` from text.
- Normalize to lowercase without `#`.
- Deduplicate while preserving stable order.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/tasks -run 'TestNormalizeLabelsFromCSV|TestExtractLabelsFromText|TestMergeTextAndFlagLabelsDeduplicates' -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/tasks/tasks.go internal/tasks/tasks_test.go
git commit -m "feat: add markdown label parsing helpers"
```

### Task 2: Teach task creation to append markdown labels

**Files:**
- Modify: `internal/tasks/tasks.go`
- Modify: `internal/cli/cli.go`
- Test: `internal/tasks/tasks_test.go`
- Test: `internal/cli/cli_test.go`

**Step 1: Write the failing tests**

Add tests covering:

```go
func TestAppendTaskAppendsNormalizedLabels(t *testing.T)
func TestAddSupportsLabelsFlag(t *testing.T)
func TestAddDeduplicatesExistingInlineLabels(t *testing.T)
```
```go
func TestAddSupportsLabelsFlag(t *testing.T) {
    stdout, stderr, err := run([]string{"add", "prep venue notes", "--labels", "flowershow,ABC"})
    if err != nil {
        t.Fatalf("err=%v stderr=%q", err, stderr)
    }
    content := readFile(t, filepath.Join(dir, ".taskgraph", "issues.md"))
    if !strings.Contains(content, "prep venue notes #flowershow #abc") {
        t.Fatalf("got %q", content)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/tasks ./internal/cli -run 'TestAppendTaskAppendsNormalizedLabels|TestAddSupportsLabelsFlag|TestAddDeduplicatesExistingInlineLabels' -v`
Expected: FAIL because `AppendTask` and `runAdd` do not accept label input yet.

**Step 3: Write minimal implementation**

Update:

- `tasks.AppendTask` to accept label input and append normalized `#label` tokens.
- `runAdd` argument parsing so `--labels` works before or after task text.
- Help text and usage error text for `tg add`.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/tasks ./internal/cli -run 'TestAppendTaskAppendsNormalizedLabels|TestAddSupportsLabelsFlag|TestAddDeduplicatesExistingInlineLabels' -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/tasks/tasks.go internal/tasks/tasks_test.go internal/cli/cli.go internal/cli/cli_test.go
git commit -m "feat: support labels on tg add"
```

### Task 3: Add inbox label filtering from markdown source

**Files:**
- Modify: `internal/cli/cli.go`
- Test: `internal/cli/cli_test.go`

**Step 1: Write the failing tests**

Add tests covering:

```go
func TestInboxFiltersByLabels(t *testing.T)
func TestInboxFiltersByLabelsAndAll(t *testing.T)
func TestInboxRejectsInvalidLabelsUsage(t *testing.T)
```
```go
func TestInboxFiltersByLabels(t *testing.T) {
    mustWrite(t, issuesPath, "- [ ] a #flowershow\n- [ ] b #other\n")
    stdout, stderr, err := run([]string{"inbox", "--labels", "flowershow"})
    if err != nil {
        t.Fatalf("err=%v stderr=%q", err, stderr)
    }
    if stdout != "- [ ] a #flowershow\n" {
        t.Fatalf("got %q", stdout)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run 'TestInboxFiltersByLabels|TestInboxFiltersByLabelsAndAll|TestInboxRejectsInvalidLabelsUsage' -v`
Expected: FAIL because inbox args only support `--all`.

**Step 3: Write minimal implementation**

Update inbox argument parsing to support:

- `--label flowershow`
- repeated `--label` flags for AND semantics

Use shared label parsing and require all requested labels to be present on each matching line.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli -run 'TestInboxFiltersByLabels|TestInboxFiltersByLabelsAndAll|TestInboxRejectsInvalidLabelsUsage' -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/cli.go internal/cli/cli_test.go
git commit -m "feat: add label filters to tg inbox"
```

### Task 4: Extend index schema and list filtering for labels

**Files:**
- Modify: `internal/indexer/indexer.go`
- Modify: `internal/indexer/sqlite.go`
- Modify: `internal/indexer/indexer_test.go`
- Modify: `internal/indexer/sqlite_test.go`
- Modify: `internal/cli/cli.go`
- Test: `internal/cli/cli_test.go`

**Step 1: Write the failing tests**

Add tests covering:

```go
func TestBuildNodesExtractsLabelsFromChecklistLines(t *testing.T)
func TestReadChecklistNodesFiltersByLabels(t *testing.T)
func TestListFiltersByLabels(t *testing.T)
```
```go
func TestListFiltersByLabels(t *testing.T) {
    mustWrite(t, filepath.Join(dir, "alpha.md"), "- [ ] Alpha #flowershow\n")
    mustWrite(t, filepath.Join(dir, "beta.md"), "- [ ] Beta #other\n")
    _, _, _ = run([]string{"index"})
    stdout, stderr, err := run([]string{"list", "--labels", "flowershow"})
    if err != nil {
        t.Fatalf("err=%v stderr=%q", err, stderr)
    }
    if !strings.Contains(stdout, "Alpha") || strings.Contains(stdout, "Beta") {
        t.Fatalf("got %q", stdout)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/indexer ./internal/cli -run 'TestBuildNodesExtractsLabelsFromChecklistLines|TestReadChecklistNodesFiltersByLabels|TestListFiltersByLabels' -v`
Expected: FAIL because nodes and DB schema do not carry labels yet.

**Step 3: Write minimal implementation**

Update:

- Node building to parse labels from checklist titles/content.
- SQLite schema/storage to persist labels as derived data.
- Query path to filter by required labels.
- `tg list` argument parsing to accept repeated `--label` flags with AND semantics.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/indexer ./internal/cli -run 'TestBuildNodesExtractsLabelsFromChecklistLines|TestReadChecklistNodesFiltersByLabels|TestListFiltersByLabels' -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/indexer/indexer.go internal/indexer/sqlite.go internal/indexer/indexer_test.go internal/indexer/sqlite_test.go internal/cli/cli.go internal/cli/cli_test.go
git commit -m "feat: index and filter markdown labels"
```

### Task 5: Document source-of-truth rule and label usage

**Files:**
- Create: `docs/DESIGN.md`
- Modify: `README.md`
- Modify: `internal/cli/cli.go`

**Step 1: Write the failing doc/test check**

Define the required documentation updates:

- `docs/DESIGN.md` states markdown is authoritative and SQLite is derived/rebuildable.
- `README.md` includes a label example for `tg add`, `tg inbox`, and `tg list`.
- CLI help text uses `--labels` for add examples and `--label` for list/inbox filtering examples.

**Step 2: Run verification to confirm docs are missing**

Run: `rg -n "source of truth|derived state|--labels" docs/DESIGN.md README.md internal/cli/cli.go`
Expected: Missing or incomplete matches before edits.

**Step 3: Write minimal implementation**

Add the design note and CLI examples without introducing new behavior.

**Step 4: Run verification to confirm docs are present**

Run: `rg -n "source of truth|derived state|--labels" docs/DESIGN.md README.md internal/cli/cli.go`
Expected: Matches in all intended files.

**Step 5: Commit**

```bash
git add docs/DESIGN.md README.md internal/cli/cli.go
git commit -m "docs: describe markdown labels and derived index"
```

### Task 6: Full verification

**Files:**
- Modify: none
- Test: `internal/tasks/tasks_test.go`
- Test: `internal/cli/cli_test.go`
- Test: `internal/indexer/indexer_test.go`
- Test: `internal/indexer/sqlite_test.go`

**Step 1: Run targeted package tests**

Run: `go test ./internal/tasks ./internal/cli ./internal/indexer -v`
Expected: PASS

**Step 2: Run full test suite**

Run: `go test ./...`
Expected: PASS

**Step 3: Inspect workspace state**

Run: `git status --short`
Expected: Only intended files changed.

**Step 4: Commit final verification checkpoint**

```bash
git add docs/DESIGN.md README.md internal/tasks/tasks.go internal/tasks/tasks_test.go internal/cli/cli.go internal/cli/cli_test.go internal/indexer/indexer.go internal/indexer/sqlite.go internal/indexer/indexer_test.go internal/indexer/sqlite_test.go docs/plans/2026-03-03-label-support-design.md docs/plans/2026-03-03-label-support.md
git commit -m "feat: add markdown-native label support"
```
