# SQLite Index Design (v1)

## Context

TaskGraph currently stores captured tasks in `.taskgraph/tasks.md` and supports `init`, `add/create`, and `list`.
This design adds an explicit index command backed by a local SQLite database, modeled after Beads' project-local database convention (`.beads/beads.db`).

## Decisions

### Storage Location

- Use project-local DB file: `.taskgraph/taskgraph.db`
- Database contents are derived/regenerable index data
- No JSON output artifact in v1

### Command

- Add `tg index`
- Command finds TaskGraph root via existing root discovery
- If no `.taskgraph` exists, return initialization error (same style as `tg list`)

### Source Files

- Recursively scan from TaskGraph project root
- Include only files ending in `.md`
- Exclude paths in dot-folders (any segment beginning with `.`)
- Exclude `node_modules`
- Always include `.taskgraph/tasks.md` explicitly as an input source, even though `.taskgraph/` is excluded from recursive scanning

### Indexed Node Types

- `file`: one node per source markdown file
- `heading`: one node per markdown heading (`#`-`######`)
- `checklist`: one node per checklist item (`- [ ]`, `- [x]`, case-insensitive x)

### SQLite Schema (v1)

Single table:

- `index_nodes`
  - `id TEXT PRIMARY KEY`
  - `kind TEXT NOT NULL`
  - `title TEXT NOT NULL`
  - `state TEXT NOT NULL` (`open|closed|unknown`)
  - `path TEXT NOT NULL` (project-relative)
  - `line INTEGER NOT NULL`
  - `parent_id TEXT`
  - `context TEXT NOT NULL`
  - `search_text TEXT NOT NULL`
  - `source TEXT NOT NULL` (`scan|tasks_md`)

Indexes:

- `idx_nodes_kind` on `(kind)`
- `idx_nodes_path` on `(path)`
- `idx_nodes_state` on `(state)`

### Rebuild Strategy

- `tg index` performs full rebuild:
  - Ensure schema exists
  - `DELETE FROM index_nodes`
  - Bulk insert all new nodes in one transaction

### IDs and Relationships

- `id` is stable hash from relative path + hierarchy + line
- `parent_id` points to file or enclosing heading
- `context` stores a breadcrumb string for fast local relevance/search usage

### Output

- Print a one-line summary:
  - indexed file count
  - indexed node count
  - db path

## Non-goals (v1)

- No `.gitignore` parsing
- No incremental index updates
- No query command over SQLite yet
- No migration framework beyond idempotent `CREATE TABLE/INDEX IF NOT EXISTS`
