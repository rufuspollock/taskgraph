package indexer

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS index_nodes (
    id TEXT PRIMARY KEY,
    kind TEXT NOT NULL,
    title TEXT NOT NULL,
    state TEXT NOT NULL,
    path TEXT NOT NULL,
    line INTEGER NOT NULL,
    parent_id TEXT,
    context TEXT NOT NULL,
    search_text TEXT NOT NULL,
    source TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_nodes_kind ON index_nodes(kind);
CREATE INDEX IF NOT EXISTS idx_nodes_path ON index_nodes(path);
CREATE INDEX IF NOT EXISTS idx_nodes_state ON index_nodes(state);
`

func RebuildSQLite(dbPath string, nodes []Node) error {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("initialize schema: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec("DELETE FROM index_nodes"); err != nil {
		return fmt.Errorf("clear index_nodes: %w", err)
	}

	stmt, err := tx.Prepare(`
INSERT INTO index_nodes
	(id, kind, title, state, path, line, parent_id, context, search_text, source)
VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()

	for _, n := range nodes {
		var parent any
		if n.ParentID != "" {
			parent = n.ParentID
		}
		if _, err := stmt.Exec(
			n.ID,
			n.Kind,
			n.Title,
			n.State,
			n.Path,
			n.Line,
			parent,
			n.Context,
			n.SearchText,
			n.Source,
		); err != nil {
			return fmt.Errorf("insert node %s: %w", n.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
