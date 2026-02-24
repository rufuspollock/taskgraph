#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck disable=SC1091
source "$ROOT_DIR/scripts/install.sh"

pass() {
  printf 'PASS: %s\n' "$1"
}

fail() {
  printf 'FAIL: %s\n' "$1"
  exit 1
}

assert_eq() {
  local got="$1"
  local want="$2"
  local label="$3"
  if [[ "$got" != "$want" ]]; then
    fail "$label (want=$want got=$got)"
  fi
  pass "$label"
}

assert_nonzero() {
  local label="$1"
  shift
  set +e
  "$@" >/dev/null 2>&1
  local rc=$?
  set -e
  if [[ $rc -eq 0 ]]; then
    fail "$label (expected non-zero exit)"
  fi
  pass "$label"
}

assert_eq "$(map_os Darwin)" "darwin" "map_os Darwin"
assert_eq "$(map_os Linux)" "linux" "map_os Linux"
assert_eq "$(map_arch x86_64)" "amd64" "map_arch x86_64"
assert_eq "$(map_arch arm64)" "arm64" "map_arch arm64"

assert_nonzero "map_os unsupported" map_os FreeBSD
assert_nonzero "map_arch unsupported" map_arch mips

TG_INSTALL_LATEST_JSON='{"tag_name":"v9.9.9"}'
assert_eq "$(resolve_version "")" "v9.9.9" "resolve_version latest"
assert_eq "$(resolve_version "v1.2.3")" "v1.2.3" "resolve_version explicit"

assert_eq "$(asset_name v1.2.3 darwin amd64)" "tg_v1.2.3_darwin_amd64.tar.gz" "asset_name"

unset INSTALL_DIR
assert_eq "$(target_dir)" "$HOME/.local/bin" "target_dir default"
INSTALL_DIR="/tmp/custom-bin"
assert_eq "$(target_dir)" "/tmp/custom-bin" "target_dir env"

echo "All installer tests passed"
