# tg add Task Types Design

## Goal

Add task type support to `tg add`/`tg create` with `br`/`bd`-style `--type` syntax, while keeping markdown files as the source of truth and SQLite as derived state.

## Decisions

- Add `-t, --type <type>` to `tg add` and `tg create`.
- Exactly one type per task.
- Store type in markdown as a namespaced label tag: `#t-<type>`.
- Keep regular labels in `#label` format.
- Validate types against a strict built-in allowlist plus optional project custom types from config.

## Type Model

### Built-in types

Built-ins are copied into code and not sourced dynamically from docs:

- `idea`
- `initiative`
- `product`
- `epic`
- `feature`
- `task`
- `subtask`
- `bug`
- `chore`
- `decision`

This includes `bd` core types plus additional taskgraph hierarchy types.

### Custom types

Projects can opt in to custom types via `.taskgraph/config.yml`:

```yaml
issue-types: research,spike
```

Effective allowed types = built-ins + normalized configured custom types.

## Storage Model

Type is persisted in `.taskgraph/issues.md` as exactly one type label token:

```text
- [ ] ➕2026-03-07 [tg-abc] map migration plan #t-epic #planning
```

Rules:

- Type label namespace is `#t-`.
- The stored type value is normalized to lowercase kebab-case.
- A task must contain at most one `#t-*` token.
- Type can be provided via `--type` or inline `#t-*` in text.
- If both are provided, they must match after normalization.

## CLI Behavior

### `tg add` / `tg create`

Support:

```bash
tg add "ship rollout checklist" --type task
tg add --type epic "cross-team launch prep"
tg create "capture context" --type idea
```

Validation:

- Missing type value after `--type` is usage error.
- Unknown type is error, with allowed type list shown.
- Multiple `--type` flags are error.
- Multiple inline `#t-*` tags are error.
- Conflicting inline and flag types are error.

## Representation Tradeoff

### Chosen: namespaced label (`#t-epic`)

Pros:

- Keeps markdown-only source-of-truth model.
- Works with existing label extraction/index patterns.
- Avoids collisions with plain domain labels.

Cons:

- Slightly noisier text than plain `#epic`.

### Rejected: plain `#epic`

- Collides with topic labels and makes type uniqueness harder to enforce.

### Rejected: dedicated inline field (`[type:epic]`)

- Introduces new syntax path and parser complexity for little gain.

## Config Handling

Extend project config parsing to read optional `issue-types` from `.taskgraph/config.yml` as comma-separated values.

Compatibility:

- Existing config remains valid.
- Missing `issue-types` means only built-ins are allowed.

## Indexing

No schema changes required for first version.

- Type is represented as label (`t-...`) and already projects through label extraction/index storage.
- `tg list --label t-epic` works immediately via existing label filter behavior.

## Testing

Cover:

- `tg add --type epic` stores `#t-epic`.
- `tg create` alias supports `--type`.
- Unknown type is rejected unless configured custom type.
- Configured custom type is accepted.
- Multiple types rejected (repeated `--type`, multiple inline `#t-*`, conflicts).
- Existing `--labels` behavior remains stable.
