# Help Command Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement a robust help UX for TaskGraph via no-args help and `-h`/`--help`.

**Architecture:** Add a dedicated help text renderer in `internal/cli`, dispatch help paths early in `Run`, and verify behavior with CLI unit tests.

**Tech Stack:** Go stdlib, existing `testing` suite.

---

### Task 1: Add failing tests for help behavior

**Files:**
- Modify: `internal/cli/cli_test.go`

**Step 1: Write failing tests**

Add tests for:
- no-args run returns nil and prints help text
- `-h` and `--help` return nil and print help text
- unknown command error contains help hint

**Step 2: Run tests to verify failure**

Run:
```bash
go test ./internal/cli -v
```

Expected: FAIL for new tests.

### Task 2: Implement help behavior

**Files:**
- Modify: `internal/cli/cli.go`

**Step 1: Add help text function**

Create `helpText()` string output with short ASCII banner and command docs.

**Step 2: Update command dispatch**

Implement:
- no args => print help and return nil
- `-h`, `--help`, `help` => print help and return nil
- unknown command => include `Run 'tg --help'` hint

**Step 3: Run tests**

Run:
```bash
go test ./internal/cli -v
go test ./...
```

Expected: PASS.

### Task 3: Commit

**Step 1: Commit code and tests**

```bash
git add internal/cli/cli.go internal/cli/cli_test.go
git commit -m "feat(cli): add full help output and help flags"
```
