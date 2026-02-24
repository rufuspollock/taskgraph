#!/usr/bin/env bash
set -euo pipefail

REPO="${TG_INSTALL_REPO:-rufuspollock/taskgraph}"

map_os() {
  local raw="$1"
  case "$raw" in
    Darwin) echo "darwin" ;;
    Linux) echo "linux" ;;
    *)
      echo "Unsupported operating system: $raw (supported: Darwin, Linux)" >&2
      return 1
      ;;
  esac
}

map_arch() {
  local raw="$1"
  case "$raw" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *)
      echo "Unsupported architecture: $raw (supported: x86_64/amd64, arm64/aarch64)" >&2
      return 1
      ;;
  esac
}

api_url_latest() {
  echo "https://api.github.com/repos/${REPO}/releases/latest"
}

parse_tag_name() {
  local json="$1"
  # Extract first tag_name value from GitHub API JSON
  echo "$json" | grep -Eo '"tag_name"[[:space:]]*:[[:space:]]*"[^"]+"' | head -n1 | sed -E 's/.*"([^"]+)"$/\1/'
}

resolve_latest_version() {
  if [[ -n "${TG_INSTALL_LATEST_JSON:-}" ]]; then
    parse_tag_name "$TG_INSTALL_LATEST_JSON"
    return
  fi

  local json
  json="$(curl -fsSL "$(api_url_latest)")"
  parse_tag_name "$json"
}

resolve_version() {
  local requested="${1:-}"
  if [[ -n "$requested" ]]; then
    echo "$requested"
    return
  fi

  local latest
  latest="$(resolve_latest_version)"
  if [[ -z "$latest" ]]; then
    echo "Failed to resolve latest release tag for ${REPO}" >&2
    return 1
  fi
  echo "$latest"
}

asset_name() {
  local version="$1"
  local os="$2"
  local arch="$3"
  echo "tg_${version}_${os}_${arch}.tar.gz"
}

target_dir() {
  echo "${INSTALL_DIR:-$HOME/.local/bin}"
}

release_asset_url() {
  local version="$1"
  local os="$2"
  local arch="$3"
  local asset
  asset="$(asset_name "$version" "$os" "$arch")"
  echo "https://github.com/${REPO}/releases/download/${version}/${asset}"
}

install_binary() {
  local version="$1"
  local os="$2"
  local arch="$3"

  local out_dir
  out_dir="$(target_dir)"
  mkdir -p "$out_dir"

  local tmpdir archive url
  tmpdir="$(mktemp -d)"
  archive="$tmpdir/tg.tar.gz"
  trap 'rm -rf "$tmpdir"' EXIT

  url="$(release_asset_url "$version" "$os" "$arch")"
  echo "Downloading ${url}"
  curl -fsSL "$url" -o "$archive"

  tar -xzf "$archive" -C "$tmpdir"
  if [[ ! -f "$tmpdir/tg" ]]; then
    echo "Downloaded archive did not contain tg binary" >&2
    return 1
  fi

  cp "$tmpdir/tg" "$out_dir/tg"
  chmod +x "$out_dir/tg"

  echo "Installed tg to $out_dir/tg"
  if [[ ":$PATH:" != *":$out_dir:"* ]]; then
    echo "Warning: $out_dir is not in PATH"
  fi
  echo "Run: tg --help"
}

main() {
  local requested_version="${1:-}"
  local raw_os raw_arch os arch version

  raw_os="$(uname -s)"
  raw_arch="$(uname -m)"
  os="$(map_os "$raw_os")"
  arch="$(map_arch "$raw_arch")"
  version="$(resolve_version "$requested_version")"

  install_binary "$version" "$os" "$arch"
}

if [[ "${TG_INSTALL_LIB_ONLY:-0}" != "1" ]]; then
  main "$@"
fi
