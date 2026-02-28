# TaskGraph

TaskGraph is a local-first CLI for quick task capture and simple task listing. 

Designed as an AI-friendly planning substrate: a lightweight, non-invasive way to build a local, queryable task graph that both humans and AI tools can use to decide what to do next.

Markdown-oriented, it adds a lightweight, natural performant interface to what you already have -- be that checklists scattered across markdown files, a need to capture quickly into an inbox, or AI oriented coding workflows.

## Quick Example

```bash
tg add "buy milk"
tg create "book dentist"
tg list
```

## Install (macOS/Linux)

Latest release:

```bash
curl -fsSL https://raw.githubusercontent.com/rufuspollock/taskgraph/main/scripts/install.sh | bash
```

Specific version:

```bash
curl -fsSL https://raw.githubusercontent.com/rufuspollock/taskgraph/main/scripts/install.sh | bash -s -- v0.1.0
```

Custom install directory:

```bash
curl -fsSL https://raw.githubusercontent.com/rufuspollock/taskgraph/main/scripts/install.sh | INSTALL_DIR="$HOME/bin" bash
```

Installer scope (v0):

- Supported OS: macOS, Linux
- Supported architectures: amd64, arm64
- Binary install location: `${INSTALL_DIR:-$HOME/.local/bin}`

## Usage

Initialize explicitly:

```bash
tg init
```

Add tasks:

```bash
tg add "buy milk"
tg create "book dentist"
```

View inbox captures:

```bash
tg inbox
```

List indexed tasks across markdown (DB-backed):

```bash
tg list
tg list --all
```

Migrate from Beads JSONL:

```bash
tg migrate-beads
```

Notes:

- `tg add` auto-initializes `.taskgraph/` in the current directory if none exists in parent directories.
- Inbox tasks are stored as checklist lines in `.taskgraph/issues.md`.
- Indexed task graph is stored in `.taskgraph/taskgraph.db`.
- `tg migrate-beads` expects both `./.beads/` and `./.taskgraph/` in the current directory.
- `tg migrate-beads` imports from `./.beads/issues.jsonl` into `./.taskgraph/issues.md`.

## Jobs To Be Done (Current Status)

1. Capture quickly: `tg add` / `tg create` - Implemented
2. Process inbox: `tg inbox` - Implemented
3. View indexed task graph: `tg list` - Implemented (v1 checklist view)
4. Graph-native planning: `tg graph` - Partial (data indexed, UX not yet shipped)
5. Suggest best next action: `tg next` - Not yet implemented

---

## Developer Guide

### Recommended Daily Workflow (Active Development)

If you are actively changing code and want `tg` to always run the latest source, use this:

```bash
mkdir -p "$HOME/.local/bin"
ln -sf "$(pwd)/scripts/tg-dev" "$HOME/.local/bin/tg"
ln -sf "$(pwd)/scripts/tg-dev" "$HOME/.local/bin/taskgraph" # optional alias
```

Then just run `tg ...` normally while developing. No `go build` needed each edit.

### Run Tests

```bash
go test ./...
```

### One-Off Local Build

```bash
go build -o tg ./cmd/tg
```

Run the built binary directly:

```bash
./tg add "test task"
```

### `go install` (Optional)

Use this when you want a compiled binary in your Go bin path:

```bash
go install ./cmd/tg
```

Note: this is not the recommended active-dev loop, because it does not auto-refresh after source edits.

### Build-Based Link (Alternative)

Use this if you want faster startup and do not mind rebuilding after edits:

```bash
go build -o tg ./cmd/tg
ln -sf "$(pwd)/tg" "$HOME/.local/bin/tg"
```

If `$HOME/.local/bin` is not in your `PATH`, add it in your shell profile.

### Release Process

- Releases are built automatically by GitHub Actions when a tag matching `v*` is pushed.

```bash
git tag v0.2.0
git push origin v0.2.0
```
