package indexer

import (
	"database/sql"
	"path/filepath"
	"slices"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRebuildSQLiteWritesNodes(t *testing.T) {
	root := t.TempDir()
	dbPath := filepath.Join(root, "taskgraph.db")
	nodes := []Node{
		{
			ID:         "a",
			Kind:       "file",
			Title:      "notes",
			State:      "unknown",
			Path:       "notes.md",
			Line:       0,
			ParentID:   "",
			Context:    "notes",
			SearchText: "notes",
			Source:     "scan",
		},
		{
			ID:         "b",
			Kind:       "checklist",
			Title:      "ship",
			State:      "open",
			Path:       "notes.md",
			Line:       4,
			ParentID:   "a",
			Context:    "notes > ship",
			SearchText: "notes ship",
			Source:     "scan",
		},
	}

	if err := RebuildSQLite(dbPath, nodes); err != nil {
		t.Fatalf("RebuildSQLite returned error: %v", err)
	}

	db := openSQLiteDB(t, dbPath)
	defer db.Close()

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM index_nodes").Scan(&count); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 rows, got %d", count)
	}
}

func TestRebuildSQLiteClearsPreviousRows(t *testing.T) {
	root := t.TempDir()
	dbPath := filepath.Join(root, "taskgraph.db")

	if err := RebuildSQLite(dbPath, []Node{{
		ID: "old", Kind: "file", Title: "old", State: "unknown", Path: "old.md", Line: 0, Context: "old", SearchText: "old", Source: "scan",
	}}); err != nil {
		t.Fatalf("first rebuild failed: %v", err)
	}

	if err := RebuildSQLite(dbPath, []Node{{
		ID: "new", Kind: "file", Title: "new", State: "unknown", Path: "new.md", Line: 0, Context: "new", SearchText: "new", Source: "scan",
	}}); err != nil {
		t.Fatalf("second rebuild failed: %v", err)
	}

	db := openSQLiteDB(t, dbPath)
	defer db.Close()

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM index_nodes").Scan(&count); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row after rebuild, got %d", count)
	}
}

func TestReadChecklistNodesFiltersByLabels(t *testing.T) {
	root := t.TempDir()
	dbPath := filepath.Join(root, "taskgraph.db")

	nodes := []Node{
		{
			ID:         "a",
			Kind:       "checklist",
			Title:      "ship #flowershow #abc",
			State:      "open",
			Path:       "notes.md",
			Line:       1,
			Context:    "notes > ship",
			SearchText: "notes ship",
			Source:     "scan",
			Labels:     []string{"flowershow", "abc"},
		},
		{
			ID:         "b",
			Kind:       "checklist",
			Title:      "other #flowershow",
			State:      "open",
			Path:       "other.md",
			Line:       1,
			Context:    "other > other",
			SearchText: "other other",
			Source:     "scan",
			Labels:     []string{"flowershow"},
		},
	}

	if err := RebuildSQLite(dbPath, nodes); err != nil {
		t.Fatalf("RebuildSQLite returned error: %v", err)
	}

	got, err := ReadChecklistNodes(dbPath, false, []string{"flowershow", "abc"})
	if err != nil {
		t.Fatalf("ReadChecklistNodes returned error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("unexpected filtered nodes: %#v", got)
	}
}

func TestReadGraphNodesReturnsTreeMetadataAndLabels(t *testing.T) {
	root := t.TempDir()
	dbPath := filepath.Join(root, "taskgraph.db")

	nodes := []Node{
		{
			ID:         "file-1",
			Kind:       "file",
			Title:      "notes",
			State:      "unknown",
			Path:       "notes.md",
			Line:       0,
			Context:    "notes",
			SearchText: "notes",
			Source:     "scan",
		},
		{
			ID:         "heading-root",
			Kind:       "heading",
			Title:      "Landscape",
			State:      "unknown",
			Path:       "notes.md",
			Line:       1,
			ParentID:   "file-1",
			Context:    "notes > Landscape",
			SearchText: "notes landscape",
			Source:     "scan",
		},
		{
			ID:         "checklist-typed",
			Kind:       "checklist",
			Title:      "Project Alpha #t-project",
			State:      "open",
			Path:       "notes.md",
			Line:       2,
			ParentID:   "heading-root",
			Context:    "notes > Landscape > Project Alpha",
			SearchText: "notes landscape project alpha",
			Source:     "scan",
			Labels:     []string{"t-project"},
		},
		{
			ID:         "checklist-child-a",
			Kind:       "checklist",
			Title:      "Task A",
			State:      "open",
			Path:       "notes.md",
			Line:       3,
			ParentID:   "checklist-typed",
			Context:    "notes > Landscape > Project Alpha > Task A",
			SearchText: "task a",
			Source:     "scan",
		},
		{
			ID:         "checklist-child-b",
			Kind:       "checklist",
			Title:      "Task B",
			State:      "open",
			Path:       "notes.md",
			Line:       4,
			ParentID:   "checklist-typed",
			Context:    "notes > Landscape > Project Alpha > Task B",
			SearchText: "task b",
			Source:     "scan",
		},
		{
			ID:         "checklist-epic",
			Kind:       "checklist",
			Title:      "Solo Epic #t-epic",
			State:      "open",
			Path:       "notes.md",
			Line:       5,
			ParentID:   "heading-root",
			Context:    "notes > Landscape > Solo Epic",
			SearchText: "solo epic",
			Source:     "scan",
			Labels:     []string{"t-epic"},
		},
	}

	if err := RebuildSQLite(dbPath, nodes); err != nil {
		t.Fatalf("RebuildSQLite returned error: %v", err)
	}

	got, err := ReadGraphNodes(dbPath)
	if err != nil {
		t.Fatalf("ReadGraphNodes returned error: %v", err)
	}

	if len(got) != len(nodes) {
		t.Fatalf("got %d nodes want %d", len(got), len(nodes))
	}

	typed := findNodeByID(t, got, "checklist-typed")
	if typed.ParentID != "heading-root" {
		t.Fatalf("typed parent = %q want %q", typed.ParentID, "heading-root")
	}
	if !slices.Equal(typed.Labels, []string{"t-project"}) {
		t.Fatalf("typed labels = %#v want %#v", typed.Labels, []string{"t-project"})
	}

	children := childrenOf(got, "checklist-typed")
	if len(children) != 2 {
		t.Fatalf("got %d children want 2", len(children))
	}
	if children[0].ID != "checklist-child-a" || children[1].ID != "checklist-child-b" {
		t.Fatalf("children order = %#v", []string{children[0].ID, children[1].ID})
	}
}

func TestReadProjectNodes(t *testing.T) {
	root := t.TempDir()
	dbPath := filepath.Join(root, "taskgraph.db")

	nodes := []Node{
		// Project A file node with t-project label, mtime=200 (newer)
		{
			ID:            "proj-a",
			Kind:          "file",
			Title:         "Project A",
			State:         "unknown",
			Path:          "project-a.md",
			Line:          0,
			Context:       "project-a",
			SearchText:    "project a",
			Source:        "scan",
			Labels:        []string{"t-project"},
			SourceMTimeUnix: 200,
		},
		// Project A: open checklist child 1
		{
			ID:       "pa-task-1",
			Kind:     "checklist",
			Title:    "Task 1",
			State:    "open",
			Path:     "project-a.md",
			Line:     1,
			ParentID: "proj-a",
			Context:  "project-a > Task 1",
			SearchText: "task 1",
			Source:   "scan",
		},
		// Project A: open checklist child 2
		{
			ID:       "pa-task-2",
			Kind:     "checklist",
			Title:    "Task 2",
			State:    "open",
			Path:     "project-a.md",
			Line:     2,
			ParentID: "proj-a",
			Context:  "project-a > Task 2",
			SearchText: "task 2",
			Source:   "scan",
		},
		// Project A: closed checklist child (should not count)
		{
			ID:       "pa-task-3",
			Kind:     "checklist",
			Title:    "Task 3",
			State:    "closed",
			Path:     "project-a.md",
			Line:     3,
			ParentID: "proj-a",
			Context:  "project-a > Task 3",
			SearchText: "task 3",
			Source:   "scan",
		},
		// Project B file node with t-project label, mtime=100 (older)
		{
			ID:            "proj-b",
			Kind:          "file",
			Title:         "Project B",
			State:         "unknown",
			Path:          "project-b.md",
			Line:          0,
			Context:       "project-b",
			SearchText:    "project b",
			Source:        "scan",
			Labels:        []string{"t-project"},
			SourceMTimeUnix: 100,
		},
		// Project B: open checklist child
		{
			ID:       "pb-task-1",
			Kind:     "checklist",
			Title:    "Task 1",
			State:    "open",
			Path:     "project-b.md",
			Line:     1,
			ParentID: "proj-b",
			Context:  "project-b > Task 1",
			SearchText: "task 1",
			Source:   "scan",
		},
		// Inbox file node without t-project label (should not appear)
		{
			ID:         "inbox",
			Kind:       "file",
			Title:      "Inbox",
			State:      "unknown",
			Path:       "inbox.md",
			Line:       0,
			Context:    "inbox",
			SearchText: "inbox",
			Source:     "scan",
		},
	}

	if err := RebuildSQLite(dbPath, nodes); err != nil {
		t.Fatalf("RebuildSQLite returned error: %v", err)
	}

	got, err := ReadProjectNodes(dbPath)
	if err != nil {
		t.Fatalf("ReadProjectNodes returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 projects, got %d: %#v", len(got), got)
	}

	// First should be Project A (mtime 200, newer)
	if got[0].ID != "proj-a" {
		t.Fatalf("expected first project ID %q, got %q", "proj-a", got[0].ID)
	}
	if got[0].OpenTaskCount != 2 {
		t.Fatalf("expected first project OpenTaskCount 2, got %d", got[0].OpenTaskCount)
	}

	// Second should be Project B (mtime 100, older)
	if got[1].ID != "proj-b" {
		t.Fatalf("expected second project ID %q, got %q", "proj-b", got[1].ID)
	}
	if got[1].OpenTaskCount != 1 {
		t.Fatalf("expected second project OpenTaskCount 1, got %d", got[1].OpenTaskCount)
	}
}

func findNodeByID(t *testing.T, nodes []Node, id string) Node {
	t.Helper()
	for _, node := range nodes {
		if node.ID == id {
			return node
		}
	}
	t.Fatalf("missing node %q", id)
	return Node{}
}

func childrenOf(nodes []Node, parentID string) []Node {
	out := make([]Node, 0)
	for _, node := range nodes {
		if node.ParentID == parentID {
			out = append(out, node)
		}
	}
	return out
}

func openSQLiteDB(t *testing.T, dbPath string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	return db
}
