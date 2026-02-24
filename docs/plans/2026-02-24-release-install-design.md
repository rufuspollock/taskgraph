# Release And Install Design (v0)

## Context

TaskGraph has a Go reboot baseline (`init`, `add/create`, `list`) and now needs a cross-platform install path that is easy for daily use. The immediate priority is macOS and Linux, with a simple `curl ... | bash` flow.

## Goals

- Publish versioned binaries automatically on Git tags.
- Support one-line installer usage with sensible defaults.
- Keep the initial solution minimal and reliable.

## Non-Goals

- Windows support in this pass.
- Homebrew/Scoop packaging in this pass.
- Complex package-manager integration before core release flow is stable.

## Selected Approach

Use tag-triggered GitHub Releases built by GoReleaser, plus a shell installer script.

Why this approach:

- Tag-driven releases keep versioning explicit and low-noise.
- GoReleaser gives reproducible multi-platform artifacts quickly.
- A single install script provides the fastest onboarding for dogfooding.

## Scope

- GitHub Actions release workflow triggered by tags matching `v*`.
- GoReleaser config that publishes tarballs for:
  - `darwin/amd64`
  - `darwin/arm64`
  - `linux/amd64`
  - `linux/arm64`
- `scripts/install.sh` that:
  - resolves version from optional arg or defaults to latest release
  - detects OS/arch
  - downloads matching asset
  - installs `tg` to `${INSTALL_DIR:-$HOME/.local/bin}`
  - prints next-step hints (including PATH warning when needed)

## Architecture

### Release pipeline

1. Developer pushes tag (for example `v0.2.0`).
2. GitHub Actions runs GoReleaser in release mode.
3. Release artifacts are attached to GitHub Release.

### Installer pipeline

1. User runs installer script via curl and shell.
2. Script detects platform and target version.
3. Script downloads and extracts matching tarball.
4. Script places `tg` in install directory and marks executable.
5. Script prints completion and verification guidance.

## Error Handling

- Unsupported OS/arch: clear message with supported matrix.
- Missing release or asset: include repo/tag/asset details in error.
- Network/download errors: include URL attempted.
- Permission errors writing install dir: suggest `INSTALL_DIR` override.

## Testing Strategy

- CI verification on push/PR:
  - `go test ./...`
  - `go build ./cmd/tg`
- Installer quality checks:
  - `shellcheck scripts/install.sh`
  - Script smoke tests (platform mapping and version resolution behavior)
- Manual smoke checks:
  - install latest release
  - install explicit version
  - verify `tg --help` (or command invocation) post-install

## Future Extensions (Deferred)

- Homebrew tap automation.
- Scoop manifest generation.
- Windows installer path.
