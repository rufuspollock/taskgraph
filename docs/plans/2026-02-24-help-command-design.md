# Help Command Design (v0)

## Goal

Add a clear CLI help experience so users can discover TaskGraph quickly:

- `tg` with no args shows full help and exits successfully.
- `tg -h` and `tg --help` show the same help.
- Help includes purpose, commands, examples, and storage notes.

## Scope

- In scope:
  - `-h`, `--help`, and `help` command behavior
  - full text help output with a small ASCII banner
  - improved unknown-command guidance
- Out of scope:
  - animated terminal effects
  - shell completion
  - man page generation

## Behavior

- `tg` -> print full help to `stdout`, exit `0`
- `tg -h` / `tg --help` / `tg help` -> print full help to `stdout`, exit `0`
- Unknown command -> return error with hint: `Run 'tg --help'`

## Content

Help text should include:

- what TaskGraph is (AI-friendly local task graph substrate)
- core commands: `init`, `add`, `create`, `list`, `help`
- quick examples
- storage location note (`.taskgraph/tasks.md`)
