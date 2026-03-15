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

CREATE TABLE IF NOT EXISTS index_node_labels (
    node_id TEXT NOT NULL,
    label TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_node_labels_node_id ON index_node_labels(node_id);
CREATE INDEX IF NOT EXISTS idx_node_labels_label ON index_node_labels(label);
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
	if _, err := tx.Exec("DELETE FROM index_node_labels"); err != nil {
		return fmt.Errorf("clear index_node_labels: %w", err)
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

	labelStmt, err := tx.Prepare(`
INSERT INTO index_node_labels
	(node_id, label)
VALUES
	(?, ?)
`)
	if err != nil {
		return fmt.Errorf("prepare label insert: %w", err)
	}
	defer labelStmt.Close()

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
		for _, label := range n.Labels {
			if _, err := labelStmt.Exec(n.ID, label); err != nil {
				return fmt.Errorf("insert node label %s/%s: %w", n.ID, label, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func ReadChecklistNodes(dbPath string, includeClosed bool, requiredLabels []string) ([]Node, error) {
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
	if len(requiredLabels) > 0 {
		query += `
 AND id IN (
	SELECT node_id
	FROM index_node_labels
	WHERE label IN (` + placeholders(len(requiredLabels)) + `)
	GROUP BY node_id
	HAVING COUNT(DISTINCT label) = ?
)`
		for _, label := range requiredLabels {
			args = append(args, label)
		}
		args = append(args, len(requiredLabels))
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

func ReadGraphNodes(dbPath string) ([]Node, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(`
SELECT id, kind, title, state, path, line, COALESCE(parent_id, ''), context, search_text, source, source_mtime_unix
FROM index_nodes
ORDER BY path ASC, line ASC, id ASC
`)
	if err != nil {
		return nil, fmt.Errorf("query graph nodes: %w", err)
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
			return nil, fmt.Errorf("scan graph node: %w", err)
		}
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate graph nodes: %w", err)
	}

	labelsByNodeID, err := readLabelsByNodeID(db)
	if err != nil {
		return nil, err
	}
	for i := range out {
		out[i].Labels = labelsByNodeID[out[i].ID]
	}
	return out, nil
}

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

func readLabelsByNodeID(db *sql.DB) (map[string][]string, error) {
	rows, err := db.Query(`
SELECT node_id, label
FROM index_node_labels
ORDER BY node_id ASC, label ASC
`)
	if err != nil {
		return nil, fmt.Errorf("query node labels: %w", err)
	}
	defer rows.Close()

	out := make(map[string][]string)
	for rows.Next() {
		var nodeID string
		var label string
		if err := rows.Scan(&nodeID, &label); err != nil {
			return nil, fmt.Errorf("scan node label: %w", err)
		}
		out[nodeID] = append(out[nodeID], label)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate node labels: %w", err)
	}
	return out, nil
}

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	out := "?"
	for i := 1; i < n; i++ {
		out += ", ?"
	}
	return out
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
