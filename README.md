# TaskGraph

TaskGraph is a local-first CLI for capturing tasks and surfacing task lists. The current reboot baseline focuses on fast daily capture with minimal setup.

**Current status (Go reboot v0):**

- **`init` command:** Initializes `.taskgraph/` in the current directory with `config.yml` and `tasks.md`.
- **`add` / `create` command:** Adds one task as a markdown checklist line to `.taskgraph/tasks.md`. Auto-initializes in the current directory if no `.taskgraph` is found while walking up parent directories.
- **`list` command:** Prints checklist task lines from `.taskgraph/tasks.md`.

## v0 usage

```bash
go build -o tg ./cmd/tg
./tg init
./tg add "buy milk"
./tg create "book dentist"
./tg list
```

## Development

```bash
go test ./...
go build ./cmd/tg
```

## Developer notes

### 2026-02-24: implementation language direction

We are rebooting implementation in Go.

Reasoning:

- Primary goal is very easy cross-platform installation for daily dogfooding.
- Go gives us straightforward single-binary distribution for macOS, Linux, and Windows.
- It preserves fast iteration speed for early CLI/product shaping (`tg add` / `tg create`) better than Rust at this stage.
- It avoids requiring end users to install and manage a Node runtime.

Decision: move forward with Go as the primary implementation language for the reboot.
