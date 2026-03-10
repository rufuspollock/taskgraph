# Design

## Source of Truth

TaskGraph stores authoritative task data in markdown files.

- `.taskgraph/issues.md` is the inbox source of truth.
- Other markdown files remain the source of truth for indexed checklist items.
- Labels are stored inline in markdown as tags such as `#flowershow`.

## Derived State

`.taskgraph/taskgraph.db` is derived state only.

- It exists to make queries and listing fast.
- It must be rebuildable from markdown at any time without data loss.
- New features should keep authoritative data in markdown rather than introducing database-only state.

## Database Schema

Canonical schema is defined in code: `internal/indexer/sqlite.go` (`const schema`).

Current tables:
- `index_nodes`: indexed file/heading/checklist nodes and hierarchy (`parent_id`), plus search/source metadata.
- `index_node_labels`: normalized label rows (`node_id`, `label`) for filtering.
