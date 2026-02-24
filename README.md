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

List tasks:

```bash
tg list
```

Notes:

- `tg add` auto-initializes `.taskgraph/` in the current directory if none exists in parent directories.
- Tasks are stored as checklist lines in `.taskgraph/tasks.md`.

## Developer Guide

Run tests:

```bash
go test ./...
```

Build locally:

```bash
go build -o tg ./cmd/tg
```

Use local build in your shell:

```bash
./tg add "test task"
```

Optional local install while developing:

```bash
go install ./cmd/tg
```

Go equivalent of `npm link` (local command linking):

```bash
# 1) Build from current checkout
go build -o tg ./cmd/tg

# 2) Link into your PATH (example target)
mkdir -p "$HOME/.local/bin"
ln -sf "$(pwd)/tg" "$HOME/.local/bin/tg"

# 3) Optional alias command name
ln -sf "$(pwd)/tg" "$HOME/.local/bin/taskgraph"
```

If `$HOME/.local/bin` is not in your `PATH`, add it in your shell profile.

Release process:

- Releases are built automatically by GitHub Actions when a tag matching `v*` is pushed.

```bash
git tag v0.2.0
git push origin v0.2.0
```
