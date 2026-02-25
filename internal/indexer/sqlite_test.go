package indexer

import (
	"database/sql"
	"path/filepath"
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

func openSQLiteDB(t *testing.T, dbPath string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	return db
}
