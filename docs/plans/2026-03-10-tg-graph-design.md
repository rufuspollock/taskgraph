# tg graph Design

## Purpose

`tg graph` should give a compact command-line overview of the task graph so the user can start from genuine root nodes, orient at a high level, and descend toward concrete work without being overwhelmed by every leaf.

This command directly supports the north-star behavior described in [docs/north-star-task-graph-vision-2026-03-02.md](/Users/rgrp/src/taskgraph/docs/north-star-task-graph-vision-2026-03-02.md): start high, choose a branch, go down fast, find a leaf, and if no useful leaf exists, learn that from the structure.

## Non-goals

- Do not build a full graphical renderer or HTML/D3 visualization in this first version.
- Do not introduce database-only graph state. The SQLite index remains a projection of markdown.
- Do not solve cross-branch dependency visualization yet. First version is based on containment hierarchy already indexed from markdown.

## Command Shape

The command is:

```text
tg graph
```

Suggested first-version flags:

```text
tg graph --depth 4
tg graph --max-children 5
tg graph --all
```

Defaults:

- `--depth`: `4`
- `--max-children`: `5`
- show open checklist items only by default
- still show heading/file structure needed to anchor open descendants

## Data Source

`tg graph` should read from `.taskgraph/taskgraph.db`, reusing the existing index built from markdown files and `.taskgraph/issues.md`.

The command should operate on the indexed containment tree:

- file nodes
- heading nodes
- checklist nodes
- existing `parent_id` relationships
- existing labels, including task-type labels stored as `#t-*`

No new source-of-truth format is required.

## Root Selection

The key product decision is selecting "genuine roots" rather than blindly starting from file nodes.

### A node is a root candidate when either:

1. it is a top-level node with at least one child
2. it has type `epic`, `project`, `product`, `initiative`, or `idea`, even if childless

### Clarifications

- "Top-level" means the node's parent is the file node, or the node otherwise sits at the highest meaningful content level in the indexed containment tree.
- Top-level nodes with no children are excluded unless they match the explicit root type allowlist above.
- Typed roots may appear below other structural nodes and should still be considered roots.

### Root Type Allowlist

The first version should treat these task types as root-worthy:

- `idea`
- `initiative`
- `project`
- `product`
- `epic`

### De-duplication and preference rules

The command should avoid printing the same branch twice when both a structural ancestor and a typed descendant qualify as roots.

Recommended rules:

1. Prefer explicitly typed root-worthy nodes over generic structural nodes.
2. Suppress a structural root if its meaningful branch is already represented by a typed descendant that also qualifies as a root.
3. Avoid showing synthetic file nodes as roots unless no meaningful content root exists for that file.

This keeps output aligned with how the user thinks about work: projects, products, initiatives, epics, and ideas rather than filenames.

## Traversal and Rendering

`tg graph` should render a text tree that is easy to scan in a terminal.

Example shape:

```text
[project] TaskGraph
  next-action UX
    tg graph command
      add root selection heuristic
      add child limiting
      add omission marker
      ... 3 more
```

This is still vertically rendered, but it should feel like left-to-right progression because each indentation level represents moving rightward through the graph.

### Traversal rules

- Start from the selected roots.
- Descend depth-first through containment children.
- Preserve source order where practical so the graph reflects document order.
- Stop when the configured depth limit is reached.
- Limit displayed children per parent to `max-children`.
- If more children exist, show an omission marker such as `... N more`.

### Filtering rules

- Closed checklist items are hidden by default.
- Headings may still be shown if they are required to connect visible descendants.
- `--all` includes closed checklist items.

### Display details

- Prefix typed root-worthy nodes with `[type]`.
- Do not prefix every node with its type; keep the output compact.
- Use stable indentation with simple ASCII connectors or plain indentation.
- Favor readability over decorative formatting.

## Ordering

Roots should be ordered so the most meaningful starting points appear first.

Recommended root ordering:

1. typed root-worthy nodes
2. other top-level branching nodes
3. childless typed roots

Within each group, preserve source/file order.

Children should preserve their indexed source order.

## Why This Design

This design is intentionally narrow:

- It supports the north-star workflow of selecting a branch before drowning in leaves.
- It builds on the current index rather than inventing a second graph model.
- It stays useful in the terminal.
- It creates a data and UX shape that can later feed richer HTML/D3 output.

The key product value is not "draw everything." It is "show me the real starting points and enough of each branch to orient quickly."

## Edge Cases

- A childless typed node should appear as a one-line root.
- A file containing only a single top-level leaf heading should not show that heading as a root unless it has a root-worthy type.
- A typed root nested under a top-level heading should still be eligible as a root.
- If all children under a root are filtered out because they are closed, the root may still be shown if it is explicitly typed and root-worthy.

## Testing Strategy

The first implementation should cover:

- root selection from top-level branching headings
- inclusion of childless typed roots
- exclusion of top-level non-branching leaves
- preference for typed roots over structural ancestors where duplication would occur
- depth limiting
- max-children truncation with omission marker
- default closed-item filtering
- `--all` behavior
- stable output ordering

## Future Extensions

- Add branch-state annotations such as `blocked` or `needs-breakdown`
- Add dependency edges beyond containment
- Add interactive narrowing
- Export the same graph data to HTML/D3 for richer visualization
