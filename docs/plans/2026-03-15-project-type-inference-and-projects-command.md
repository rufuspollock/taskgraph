# Project Type Inference and `tg projects` Command

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make `tg index` automatically infer standalone markdown files as projects (type `t-project`), and add a `tg projects` command that shows a summary view of all projects with open task counts and last-touched dates.

**Architecture:** Extend `BuildNodes()` in the indexer to attach a `t-project` label to file nodes (excluding `.taskgraph/issues.md`). Add a new `ReadProjectNodes()` query in the SQLite layer that returns file nodes labelled `t-project` with aggregated stats. Add a `tg projects` CLI command that renders the summary. Keep markdown as source of truth; the `t-project` label is inferred at index time, not written back to files.

**Tech Stack:** Go, stdlib, existing `internal/cli`, `internal/indexer`, `internal/tasks`, SQLite via `modernc.org/sqlite`

**Context:** See `docs/projects-vs-tasks-2026-03-15.md` for the design rationale. The key insight: we need a binary project/task distinction to support the core movement (start high, pick a project, descend to leaf tasks). File-based inference is the simplest reliable heuristic.

---

### Task 1: Add failing test for file-node project inference in the indexer

**Files:**
- Modify: `internal/indexer/indexer_test.go`
- Inspect: `internal/indexer/indexer.go`

**Step 1: Write the failing test**

Add a test `TestBuildNodesInfersProjectForFiles` that:
- Creates a temp dir with two markdown files: `projects/my-project.md` (with some checklist items inside) and `.taskgraph/issues.md` (with some inbox items)
- Calls `BuildNodes(root)`
- Asserts that the file node for `projects/my-project.md` has label `t-project` in its `Labels` slice
- Asserts that the file node for `.taskgraph/issues.md` does NOT have label `t-project`

```go
func TestBuildNodesInfersProjectForFiles(t *testing.T) {
	root := t.TempDir()
	mustMkdirAll(t, filepath.Join(root, "projects"))
	mustMkdirAll(t, filepath.Join(root, ".taskgraph"))
	mustWrite(t, filepath.Join(root, "projects", "my-project.md"),
		"# My Project\n\n- [ ] task one\n- [ ] task two\n")
	mustWrite(t, filepath.Join(root, ".taskgraph", "issues.md"),
		"- [ ] inbox item\n")

	nodes, err := indexer.BuildNodes(root)
	if err != nil {
		t.Fatal(err)
	}

	var projectFileNode, issuesFileNode *indexer.Node
	for i, n := range nodes {
		if n.Kind == "file" && n.Path == "projects/my-project.md" {
			projectFileNode = &nodes[i]
		}
		if n.Kind == "file" && n.Path == ".taskgraph/issues.md" {
			issuesFileNode = &nodes[i]
		}
	}

	if projectFileNode == nil {
		t.Fatal("expected file node for projects/my-project.md")
	}
	if !hasLabel(projectFileNode.Labels, "t-project") {
		t.Errorf("expected projects/my-project.md to have t-project label, got %v", projectFileNode.Labels)
	}

	if issuesFileNode == nil {
		t.Fatal("expected file node for .taskgraph/issues.md")
	}
	if hasLabel(issuesFileNode.Labels, "t-project") {
		t.Errorf("expected .taskgraph/issues.md NOT to have t-project label, got %v", issuesFileNode.Labels)
	}
}

func hasLabel(labels []string, target string) bool {
	for _, l := range labels {
		if l == target {
			return true
		}
	}
	return false
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/indexer -run TestBuildNodesInfersProjectForFiles -v`

Expected: FAIL because file nodes currently have nil Labels.

**Step 3: Commit**

```bash
git add internal/indexer/indexer_test.go
git commit -m "test: add failing test for file-node project inference"
```

---

### Task 2: Implement file-node project inference in BuildNodes

**Files:**
- Modify: `internal/indexer/indexer.go`

**Step 1: Implement the inference**

In `BuildNodes()`, after building all nodes for a file, add the `t-project` label to the file node if:
- The file's source is `"scan"` (i.e., not `.taskgraph/issues.md`)

This is the simplest possible heuristic: every scanned markdown file is a project candidate. The `.taskgraph/issues.md` inbox is excluded.

In the `indexMarkdown` function, add a `source` parameter check is not needed since `indexMarkdown` does not know about source. Instead, do it in `BuildNodes()` after calling `indexMarkdown`, by finding the file node (always the first node returned) and appending the label.

```go
// In BuildNodes(), after: fileNodes := indexMarkdown(...)
// Add project label to file nodes from scanned sources (not inbox)
if source == "scan" && len(fileNodes) > 0 && fileNodes[0].Kind == "file" {
	fileNodes[0].Labels = append(fileNodes[0].Labels, tasks.TypeLabel("project"))
}
```

**Step 2: Run the test to verify it passes**

Run: `go test ./internal/indexer -run TestBuildNodesInfersProjectForFiles -v`

Expected: PASS

**Step 3: Run all tests to check for regressions**

Run: `go test ./... -v`

Expected: All tests pass. Some graph tests may need updating if they now see `t-project` labels on file nodes that they did not expect.

**Step 4: Commit**

```bash
git add internal/indexer/indexer.go
git commit -m "feat: infer t-project label for scanned file nodes at index time"
```

---

### Task 3: Add failing test for ReadProjectNodes query

**Files:**
- Modify: `internal/indexer/sqlite_test.go`
- Inspect: `internal/indexer/sqlite.go`

**Step 1: Write the failing test**

Add a test `TestReadProjectNodes` that:
- Creates nodes including two file nodes with `t-project` label and one without
- Adds checklist children under each project file
- Stores them in SQLite via `RebuildSQLite`
- Calls a new `ReadProjectNodes(dbPath)` function
- Asserts it returns exactly 2 project nodes
- Asserts each has correct `OpenTaskCount` and `SourceMTimeUnix`

```go
func TestReadProjectNodes(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	now := time.Now().Unix()

	nodes := []indexer.Node{
		{ID: "f1", Kind: "file", Title: "project-alpha", State: "unknown", Path: "project-alpha.md", Line: 0, ParentID: "", Labels: []string{"t-project"}, SourceMTimeUnix: now},
		{ID: "t1", Kind: "checklist", Title: "task one", State: "open", Path: "project-alpha.md", Line: 2, ParentID: "f1", SourceMTimeUnix: now},
		{ID: "t2", Kind: "checklist", Title: "task two", State: "closed", Path: "project-alpha.md", Line: 3, ParentID: "f1", SourceMTimeUnix: now},
		{ID: "t3", Kind: "checklist", Title: "task three", State: "open", Path: "project-alpha.md", Line: 4, ParentID: "f1", SourceMTimeUnix: now},
		{ID: "f2", Kind: "file", Title: "project-beta", State: "unknown", Path: "project-beta.md", Line: 0, ParentID: "", Labels: []string{"t-project"}, SourceMTimeUnix: now - 100},
		{ID: "t4", Kind: "checklist", Title: "task four", State: "open", Path: "project-beta.md", Line: 2, ParentID: "f2", SourceMTimeUnix: now - 100},
		{ID: "f3", Kind: "file", Title: "not-a-project", State: "unknown", Path: ".taskgraph/issues.md", Line: 0, ParentID: "", Labels: nil, SourceMTimeUnix: now},
		{ID: "t5", Kind: "checklist", Title: "inbox item", State: "open", Path: ".taskgraph/issues.md", Line: 1, ParentID: "f3", SourceMTimeUnix: now},
	}

	if err := indexer.RebuildSQLite(dbPath, nodes); err != nil {
		t.Fatal(err)
	}

	projects, err := indexer.ReadProjectNodes(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}

	// Ordered by source_mtime_unix DESC (most recent first)
	if projects[0].Title != "project-alpha" {
		t.Errorf("expected first project to be project-alpha, got %s", projects[0].Title)
	}
	if projects[0].OpenTaskCount != 2 {
		t.Errorf("expected project-alpha to have 2 open tasks, got %d", projects[0].OpenTaskCount)
	}
	if projects[1].Title != "project-beta" {
		t.Errorf("expected second project to be project-beta, got %s", projects[1].Title)
	}
	if projects[1].OpenTaskCount != 1 {
		t.Errorf("expected project-beta to have 1 open task, got %d", projects[1].OpenTaskCount)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/indexer -run TestReadProjectNodes -v`

Expected: FAIL because `ReadProjectNodes` does not exist and `ProjectNode` type does not exist.

**Step 3: Commit**

```bash
git add internal/indexer/sqlite_test.go
git commit -m "test: add failing test for ReadProjectNodes query"
```

---

### Task 4: Implement ReadProjectNodes in SQLite layer

**Files:**
- Modify: `internal/indexer/sqlite.go`

**Step 1: Add the ProjectNode type and ReadProjectNodes function**

```go
// ProjectNode is a project-level summary returned by ReadProjectNodes.
type ProjectNode struct {
	ID              string
	Title           string
	Path            string
	OpenTaskCount   int
	SourceMTimeUnix int64
}

func ReadProjectNodes(dbPath string) ([]ProjectNode, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(`
SELECT
    f.id,
    f.title,
    f.path,
    f.source_mtime_unix,
    COUNT(CASE WHEN c.state = 'open' THEN 1 END) AS open_task_count
FROM index_nodes f
JOIN index_node_labels l ON l.node_id = f.id AND l.label = 't-project'
LEFT JOIN index_nodes c ON c.path = f.path AND c.kind = 'checklist'
WHERE f.kind = 'file'
GROUP BY f.id
ORDER BY f.source_mtime_unix DESC, f.path ASC
`)
	if err != nil {
		return nil, fmt.Errorf("query project nodes: %w", err)
	}
	defer rows.Close()

	var out []ProjectNode
	for rows.Next() {
		var p ProjectNode
		if err := rows.Scan(&p.ID, &p.Title, &p.Path, &p.SourceMTimeUnix, &p.OpenTaskCount); err != nil {
			return nil, fmt.Errorf("scan project node: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate project nodes: %w", err)
	}
	return out, nil
}
```

**Step 2: Run test to verify it passes**

Run: `go test ./internal/indexer -run TestReadProjectNodes -v`

Expected: PASS

**Step 3: Run all tests**

Run: `go test ./... -v`

Expected: All pass.

**Step 4: Commit**

```bash
git add internal/indexer/sqlite.go
git commit -m "feat: add ReadProjectNodes query with open task counts"
```

---

### Task 5: Add failing test for `tg projects` CLI command

**Files:**
- Modify: `internal/cli/cli_test.go`

**Step 1: Write the failing test**

Add a test `TestProjectsCommand` that:
- Sets up a temp dir with `.taskgraph/` initialized, a project file with checklist items, and a pre-built index
- Runs `tg projects`
- Asserts output contains the project name and open task count

```go
func TestProjectsCommand(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	mustMkdirAll(t, filepath.Join(dir, ".taskgraph"))
	mustWrite(t, filepath.Join(dir, ".taskgraph", "config.yml"), "prefix: tg\n")
	mustWrite(t, filepath.Join(dir, ".taskgraph", "issues.md"), "")
	mustMkdirAll(t, filepath.Join(dir, "projects"))
	mustWrite(t, filepath.Join(dir, "projects", "alpha.md"),
		"# Alpha\n\n- [ ] do thing one\n- [ ] do thing two\n- [x] done thing\n")

	// Build the index first
	var idxOut, idxErr strings.Builder
	if err := cli.Run([]string{"index"}, &idxOut, &idxErr); err != nil {
		t.Fatalf("index failed: %v: %s", err, idxErr.String())
	}

	var stdout, stderr strings.Builder
	err := cli.Run([]string{"projects"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("projects failed: %v: %s", err, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected output to contain 'alpha', got: %s", out)
	}
	if !strings.Contains(out, "2") {
		t.Errorf("expected output to contain open task count '2', got: %s", out)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run TestProjectsCommand -v`

Expected: FAIL because `tg projects` is not a recognized command.

**Step 3: Commit**

```bash
git add internal/cli/cli_test.go
git commit -m "test: add failing test for tg projects command"
```

---

### Task 6: Implement `tg projects` CLI command

**Files:**
- Modify: `internal/cli/cli.go`

**Step 1: Add the command dispatch**

In `Run()`, add a case for `"projects"`:

```go
case "projects":
    return runProjects(stdout, stderr)
```

**Step 2: Implement runProjects**

```go
func runProjects(stdout io.Writer, stderr io.Writer) error {
	cwd, err := effectiveCWD()
	if err != nil {
		return err
	}
	root, found, err := project.FindTaskgraphRoot(cwd)
	if err != nil {
		return err
	}
	if !found {
		fmt.Fprintln(stderr, "No .taskgraph found. Run `tg init` or `tg add \"task text\"`.")
		return errors.New("not initialized")
	}

	dbPath := filepath.Join(root, ".taskgraph", "taskgraph.db")
	projects, err := indexer.ReadProjectNodes(dbPath)
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		fmt.Fprintln(stdout, "No projects found. Run `tg index` to scan markdown files.")
		return nil
	}

	for _, p := range projects {
		fmt.Fprintf(stdout, "%-40s  %d open  (%s)\n", p.Title, p.OpenTaskCount, p.Path)
	}
	return nil
}
```

**Step 3: Update help text**

Add `projects` to the help text and command list in `helpText()`.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli -run TestProjectsCommand -v`

Expected: PASS

**Step 5: Run all tests**

Run: `go test ./... -v`

Expected: All pass.

**Step 6: Commit**

```bash
git add internal/cli/cli.go
git commit -m "feat: add tg projects command showing project summary view"
```

---

### Task 7: Fix any graph test regressions from project inference

**Files:**
- Modify: `internal/cli/cli_test.go` (if needed)
- Inspect: `internal/indexer/indexer_test.go`

**Step 1: Run all tests and identify failures**

Run: `go test ./... -v`

File nodes now carry `t-project` labels, which means `isTypedGraphRoot()` will return true for them. This changes graph root selection behavior — file nodes that were previously structural containers may now appear as typed roots.

**Step 2: Assess whether graph behavior change is correct**

The graph already treats `project` as a root type in `graphRootTypes`. File nodes with `t-project` will now be selected as typed roots by `selectGraphRoots`. This is actually the desired behavior — project files should appear as top-level roots in `tg graph`.

However, `graphVisibility` currently only marks checklist nodes as visible when they are typed roots (line ~702: `if roots[id] && isTypedGraphRoot(node)`). File nodes with `kind == "file"` are only visible if they have visible children. This should still work correctly since project files with open tasks will have visible children.

**Step 3: Fix any failing tests**

Update test expectations that assumed file nodes would not be typed roots. The fix will depend on what specifically fails.

**Step 4: Run all tests**

Run: `go test ./... -v`

Expected: All pass.

**Step 5: Commit**

```bash
git add -A
git commit -m "fix: update graph tests for file-node project inference"
```

---

### Task 8: Manual smoke test

**Step 1: Build and run against a real project**

```bash
go build -o tg ./cmd/tg
./tg index
./tg projects
./tg graph
```

**Step 2: Verify output**

- `tg projects` should list markdown files with their open task counts
- `tg graph` should show project files as roots where they have open tasks
- `.taskgraph/issues.md` should NOT appear as a project

**Step 3: Commit any final adjustments**

```bash
git add -A
git commit -m "chore: smoke test adjustments for project inference"
```
