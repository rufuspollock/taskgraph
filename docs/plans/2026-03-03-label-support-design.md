# Label Support Design

## Goal

Add basic label support to `tg` using markdown-native tags like `#flowershow` as the only source of truth, while keeping SQLite as a rebuildable index only.

## Decisions

- Use inline markdown tags such as `#flowershow` and `#abc` in checklist task text.
- Add `--labels` to `tg add` so labels can be supplied as comma-separated values and appended as tags.
- Use `--label` for filtering in `tg inbox` and `tg list`, matching `br`/beads list behavior.
- Treat markdown files as authoritative data and SQLite as derived cache only.
- Do not add separate `tg label ...` management commands in this first version.

## Storage Model

Tasks remain stored in `.taskgraph/issues.md` as plain markdown checklist lines. Labels are appended to the end of the task text as normalized tags.

Example:

```text
- [ ] ➕2026-03-03 [task-xyz] prep venue notes #flowershow #abc
```

Rules:

- Normalize labels from flags to lowercase.
- Accept comma-separated bare labels from `--labels`.
- Write labels with leading `#`.
- Deduplicate labels across existing inline tags and flag-provided labels.
- Preserve existing task text and task format.

## CLI Behavior

### `tg add`

Support:

```bash
tg add "prep venue notes" --labels flowershow,abc
tg add --labels flowershow,abc "prep venue notes"
```

Behavior:

- Parse `--labels` anywhere in the argument list.
- Merge labels from the flag with any `#tags` already present in the task text.
- Append normalized labels once each at the end of the stored task text.

### `tg inbox`

Support:

```bash
tg inbox --label flowershow
tg inbox --label flowershow --label abc
```

Behavior:

- Filter raw checklist lines from `.taskgraph/issues.md`.
- Require all requested labels to be present.
- Continue excluding closed items by default unless `--all` is passed.

### `tg list`

Support:

```bash
tg list --label flowershow
tg list --label flowershow --label abc
```

Behavior:

- Filter indexed checklist nodes using labels parsed from markdown content.
- Require all requested labels to be present.
- Continue excluding closed items by default unless `--all` is passed.

## Parsing

Implement a small shared parser for markdown labels.

Constraints:

- Recognize tag tokens intended as labels, such as `#flowershow`.
- Normalize extracted labels to lowercase without the `#`.
- Avoid treating arbitrary `#` fragments as labels where practical.
- Keep parsing logic close to task and CLI text handling so index and inbox behavior stay consistent.

## Indexing

SQLite remains a projection of markdown content.

- The indexer may persist parsed labels for fast filtering.
- Rebuilding the database from markdown must fully restore label behavior.
- No label data should exist only in SQLite.

## Documentation

Add a short design note in `docs/DESIGN.md` stating:

- Markdown files are the source of truth.
- SQLite is derived state and may be rebuilt at any time without data loss.

Also update CLI help text and examples to show `--labels`.

## Testing

Cover:

- `tg add --labels` appends tags correctly.
- Existing inline tags and flag labels are deduplicated.
- `tg inbox --labels` filters open tasks correctly.
- `tg inbox --labels ... --all` includes closed matches.
- `tg list --labels` filters indexed tasks correctly after `tg index`.
- Rebuilding the index preserves label queries from markdown alone.
