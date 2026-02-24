# Task Types (Rough Schema)

This is a rough, useful schema rather than a strict hierarchy. `idea` is orthogonal to the hierarchy and can connect to items at multiple levels.

- `idea`: A raw thought, opportunity, or problem statement. Orthogonal to the main tree; usually relates to `initiative`, `product`, or `epic`/`feature`.
- `initiative`: Top-level effort, often company/org scale (for example, creating a new company, org, or major strategic program).
- `product`: A product-level area or surface within an initiative.
- `epic` / `feature`: A substantial chunk of work under a product. These may overlap in meaning and can be treated as near-equivalents for now.
- `task`: A concrete, actionable work item.
- `subtask`: A smaller step within a task; can be nested further as needed.

## Rough hierarchy (not strict)

```text
idea (orthogonal)
  |\
  | +--> can relate to initiative
  | +--> can relate to product
  | +--> can relate to epic / feature
  |
  +--> (not required to fit cleanly in the tree)

initiative (top level; company/org-scale)
  |
  +--> product
         |
         +--> epic / feature
                |
                +--> task
                       |
                       +--> subtask
                              |
                              +--> subtask (deeper as needed)

      [You can insert additional levels anywhere as needed]
```
