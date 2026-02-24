# Release And Installer Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add tag-triggered GitHub Releases for `tg` and a `curl ... | bash` installer for macOS/Linux that installs the latest (or specified) version.

**Architecture:** Use GoReleaser to build/publish multi-arch binaries on Git tags, orchestrated by GitHub Actions. Add a POSIX shell installer (`scripts/install.sh`) that detects OS/arch, resolves version (latest by default), downloads the right asset, and installs into `${INSTALL_DIR:-$HOME/.local/bin}`.

**Tech Stack:** Go, GoReleaser, GitHub Actions, POSIX shell, shellcheck.

---

### Task 1: Add GoReleaser configuration for release artifacts

**Files:**
- Create: `.goreleaser.yaml`

**Step 1: Write the failing packaging check**

Run:
```bash
goreleaser check
```

Expected: FAIL (config missing).

**Step 2: Add minimal `.goreleaser.yaml`**

Include:
- project name `tg`
- build from `./cmd/tg`
- targets:
  - `darwin_amd64`
  - `darwin_arm64`
  - `linux_amd64`
  - `linux_arm64`
- archive format `tar.gz`
- binary name `tg`
- release notes mode set to minimal/auto

**Step 3: Re-run validation**

Run:
```bash
goreleaser check
```

Expected: PASS.

**Step 4: Commit**

```bash
git add .goreleaser.yaml
git commit -m "build(release): add goreleaser config for tg binaries"
```

### Task 2: Add tag-triggered GitHub Actions release workflow

**Files:**
- Create: `.github/workflows/release.yml`

**Step 1: Write failing workflow sanity check**

Run:
```bash
test -f .github/workflows/release.yml
```

Expected: FAIL (file missing).

**Step 2: Add release workflow**

Workflow requirements:
- Trigger on push tags `v*`
- Permissions for release publishing (`contents: write`)
- Setup Go
- Checkout code
- Run `go test ./...`
- Run GoReleaser release action with `--clean`

**Step 3: Validate workflow syntax**

Run:
```bash
python - <<'PY'
import yaml, sys
yaml.safe_load(open(".github/workflows/release.yml"))
print("ok")
PY
```

Expected: PASS (`ok`).

**Step 4: Commit**

```bash
git add .github/workflows/release.yml
git commit -m "ci(release): publish tg binaries on version tags"
```

### Task 3: Implement installer script with latest-version default

**Files:**
- Create: `scripts/install.sh`
- Create: `scripts/install_test.sh`

**Step 1: Write failing installer tests**

Add shell tests covering:
- OS/arch mapping (`darwin/linux`, `amd64/arm64`)
- latest release URL resolution when no version provided
- explicit version override path
- unsupported platform exits non-zero
- install path default and override behavior

Run:
```bash
bash scripts/install_test.sh
```

Expected: FAIL (installer missing).

**Step 2: Implement `scripts/install.sh`**

Behavior:
- `set -eu`
- parse optional version arg
- detect platform via `uname -s` / `uname -m`
- resolve version:
  - explicit arg if provided
  - otherwise query GitHub Releases latest API
- construct asset name from version + os + arch
- download via `curl -fsSL`
- extract via `tar -xzf`
- install binary to `${INSTALL_DIR:-$HOME/.local/bin}`
- chmod +x
- print success + PATH hint if install dir absent from `PATH`

**Step 3: Re-run installer tests**

Run:
```bash
bash scripts/install_test.sh
```

Expected: PASS.

**Step 4: Run shellcheck**

Run:
```bash
shellcheck scripts/install.sh
```

Expected: PASS (no errors).

**Step 5: Commit**

```bash
git add scripts/install.sh scripts/install_test.sh
git commit -m "feat(install): add curl-install script for macOS and Linux"
```

### Task 4: Document release and installer usage

**Files:**
- Modify: `README.md`

**Step 1: Add release/install documentation**

Document:
- tag-based release model
- installer command examples:
  - latest
  - explicit version
  - custom `INSTALL_DIR`
- supported platforms matrix for v0

**Step 2: Verify docs include required snippets**

Run:
```bash
rg -n "install.sh|GitHub Releases|v\\*|INSTALL_DIR|darwin|linux" README.md
```

Expected: matching lines found.

**Step 3: Commit**

```bash
git add README.md
git commit -m "docs: add release and installer instructions"
```

### Task 5: End-to-end verification before completion

**Files:**
- No new files; verification only

**Step 1: Run full project verification**

Run:
```bash
go test ./...
go build ./cmd/tg
goreleaser check
bash scripts/install_test.sh
shellcheck scripts/install.sh
```

Expected: all commands PASS.

**Step 2: Local installer smoke test**

Run:
```bash
tmp="$(mktemp -d)"
INSTALL_DIR="$tmp/bin" bash scripts/install.sh <explicit-version-or-fixture>
```

Expected: `tg` binary installed executable in `$tmp/bin` (for local CI simulation, use mocked download path if needed).

**Step 3: Final commit (only if verification-induced edits)**

```bash
git add -A
git commit -m "chore: finalize release/install pipeline verification fixes"
```

## Execution Notes

- Keep implementation minimal for v0; do not add Homebrew/Scoop/Windows in this plan.
- Preserve backward compatibility of current Go CLI commands.
- If local `goreleaser` or `shellcheck` is missing, document that in execution output and validate via available checks plus CI configuration.
