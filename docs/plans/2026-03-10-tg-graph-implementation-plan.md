# tg graph Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a `tg graph` CLI command that prints a compact tree overview from genuine root nodes, using the existing indexed containment graph and bounded expansion.

**Architecture:** Keep markdown as the source of truth and reuse `.taskgraph/taskgraph.db` as a projection. Extend the index read path to expose enough graph metadata for root selection and rendering, implement `tg graph` in `internal/cli`, and keep the output intentionally simple: typed and structural roots, depth-bounded traversal, and omission markers for large child sets.

**Tech Stack:** Go, stdlib, existing `internal/cli`, `internal/indexer`, SQLite via `modernc.org/sqlite`

---

### Task 1: Add failing graph-query tests in the indexer layer

**Files:**
- Modify: `internal/indexer/sqlite_test.go`
- Inspect: `internal/indexer/sqlite.go`
- Inspect: `internal/indexer/indexer.go`

**Step 1: Write the failing test**

Add tests that build a small markdown tree, store it in SQLite, and then assert that a new graph-query helper returns enough data to:

- identify parent-child relationships
- expose labels for typed nodes
- distinguish top-level nodes from deeper descendants
- preserve source order for children

Use fixtures that include:

- a top-level heading with children
- a top-level heading with no children
- a childless `#t-epic`
- a nested `#t-project`

**Step 2: Run test to verify it fails**

Run: `go test ./internal/indexer -run 'Graph|Root' -v`

Expected: FAIL because the graph-query helper does not exist yet.

**Step 3: Write minimal implementation**

In `internal/indexer/sqlite.go`, add a read helper that loads indexed nodes and labels in a graph-friendly form, including:

- node id
- parent id
- kind
- title
- state
- path
- line
- labels

Keep it read-only and derived from existing tables.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/indexer -run 'Graph|Root' -v`

Expected: PASS

**Step 5: Commit**

```bash
git add internal/indexer/sqlite.go internal/indexer/sqlite_test.go
git commit -m "feat: add graph node query for indexed tree rendering"
```

### Task 2: Add failing CLI tests for `tg graph` root selection

**Files:**
- Modify: `internal/cli/cli_test.go`
- Inspect: `internal/cli/cli.go`
- Inspect: `internal/tasks/tasks.go`

**Step 1: Write the failing test**

Add CLI tests that set up a temporary `.taskgraph`, create markdown content, build the index, and assert:

- `tg graph` includes top-level branching nodes as roots
- `tg graph` excludes top-level childless non-typed nodes
- `tg graph` includes childless typed roots for `idea`, `initiative`, `project`, `product`, and `epic`
- `tg graph` includes nested typed roots when they qualify
- structural ancestors are not duplicated when a typed descendant is the meaningful root

Use exact expected output strings, not partial contains checks, for the final rendered tree.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run 'Graph' -v`

Expected: FAIL because `tg graph` does not exist.

**Step 3: Write minimal implementation**

In `internal/cli/cli.go`:

- register `graph` in help text and command dispatch
- add `runGraph`
- load the indexed graph nodes from SQLite
- implement root selection using the approved heuristics

Keep flag handling minimal at first if needed, but the command must exist and select roots correctly.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli -run 'Graph' -v`

Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/cli.go internal/cli/cli_test.go
git commit -m "feat: add tg graph root selection"
```

### Task 3: Add failing CLI tests for depth and child limiting

**Files:**
- Modify: `internal/cli/cli_test.go`
- Inspect: `internal/cli/cli.go`

**Step 1: Write the failing test**

Add tests covering:

- default depth limit
- `--depth` override
- default `--max-children` limit
- omission marker output as `... N more`
- stable child ordering

Use one fixture with a deep branch and another with more than five children under one parent.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run 'Graph.*Depth|Graph.*Children|Graph.*Omission' -v`

Expected: FAIL because traversal limiting is not implemented yet.

**Step 3: Write minimal implementation**

Extend `runGraph` and helper functions to:

- parse `--depth`
- parse `--max-children`
- stop traversal when depth is exhausted
- cap children per parent
- append `... N more` when hidden children remain

Keep rendering ASCII-only.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli -run 'Graph.*Depth|Graph.*Children|Graph.*Omission' -v`

Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/cli.go internal/cli/cli_test.go
git commit -m "feat: add bounded graph traversal output"
```

### Task 4: Add failing CLI tests for open/closed filtering

**Files:**
- Modify: `internal/cli/cli_test.go`
- Inspect: `internal/cli/cli.go`

**Step 1: Write the failing test**

Add tests asserting:

- closed checklist items are hidden by default
- `--all` includes closed checklist items
- headings remain visible when needed to connect visible descendants
- typed childless roots can still appear even if no open children remain

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run 'Graph.*All|Graph.*Closed|Graph.*Visible' -v`

Expected: FAIL because filtering behavior is incomplete.

**Step 3: Write minimal implementation**

Implement filtering rules in graph rendering:

- hide closed checklist nodes unless `--all`
- keep structural ancestors needed to reach visible descendants
- allow explicitly typed root-worthy roots to remain visible

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli -run 'Graph.*All|Graph.*Closed|Graph.*Visible' -v`

Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/cli.go internal/cli/cli_test.go
git commit -m "feat: filter tg graph output by open state"
```

### Task 5: Document the command and verify behavior end-to-end

**Files:**
- Modify: `README.md`
- Modify: `internal/cli/cli.go`
- Modify: `docs/plans/2026-03-10-tg-graph-design.md` (only if wording needs correction after implementation)

**Step 1: Write the failing doc/help test**

If the help text is covered by existing tests, add or extend one to assert:

- `tg graph` appears in help output
- the help synopsis mentions `--depth` and `--max-children` if those are documented inline

If no explicit help test is needed, write the doc update directly after confirming behavior in tests.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run 'Help|Graph' -v`

Expected: FAIL if help assertions were added, otherwise skip to implementation.

**Step 3: Write minimal implementation**

Update:

- CLI help text
- README examples for `tg graph`
- command documentation wording

Keep examples small and aligned with the actual output style.

**Step 4: Run verification**

Run: `go test ./internal/cli ./internal/indexer -v`

Expected: PASS

Run: `go test ./...`

Expected: PASS

**Step 5: Commit**

```bash
git add README.md internal/cli/cli.go internal/cli/cli_test.go internal/indexer/sqlite.go internal/indexer/sqlite_test.go docs/plans/2026-03-10-tg-graph-design.md docs/plans/2026-03-10-tg-graph-implementation-plan.md
git commit -m "feat: add tg graph command"
```
