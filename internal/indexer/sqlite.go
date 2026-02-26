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
    source TEXT NOT NULL,
    source_mtime_unix INTEGER NOT NULL DEFAULT 0
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
	if err := ensureColumn(db, "index_nodes", "source_mtime_unix", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return fmt.Errorf("ensure source_mtime_unix column: %w", err)
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
	(id, kind, title, state, path, line, parent_id, context, search_text, source, source_mtime_unix)
VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
			n.SourceMTimeUnix,
		); err != nil {
			return fmt.Errorf("insert node %s: %w", n.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func ReadChecklistNodes(dbPath string, includeClosed bool) ([]Node, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	query := `
SELECT id, kind, title, state, path, line, COALESCE(parent_id, ''), context, search_text, source, source_mtime_unix
FROM index_nodes
WHERE kind = 'checklist'
`
	var args []any
	if !includeClosed {
		query += " AND state = ?"
		args = append(args, "open")
	}
	query += " ORDER BY source_mtime_unix DESC, path ASC, line ASC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query checklist nodes: %w", err)
	}
	defer rows.Close()

	out := []Node{}
	for rows.Next() {
		var n Node
		if err := rows.Scan(
			&n.ID,
			&n.Kind,
			&n.Title,
			&n.State,
			&n.Path,
			&n.Line,
			&n.ParentID,
			&n.Context,
			&n.SearchText,
			&n.Source,
			&n.SourceMTimeUnix,
		); err != nil {
			return nil, fmt.Errorf("scan checklist node: %w", err)
		}
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate checklist nodes: %w", err)
	}
	return out, nil
}

func ensureColumn(db *sql.DB, table, column, columnDef string) error {
	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull int
		var dflt sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return err
		}
		if name == column {
			return nil
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE " + table + " ADD COLUMN " + column + " " + columnDef)
	return err
}
