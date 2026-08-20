#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
formal_dir="$root/docs/formal/feature-0016"
state_dir="$(mktemp -d "${TMPDIR:-/tmp}/feature-0016-tlc.XXXXXX")"
trap 'rm -rf "$state_dir"' EXIT

run_tlc() {
  local model="$1"
  local cfg="$2"
  if command -v tlc >/dev/null 2>&1; then
    tlc -metadir "$state_dir" -config "$cfg" "$model"
    return
  fi
  if [[ -n "${TLA2TOOLS_JAR:-}" && -f "${TLA2TOOLS_JAR}" ]]; then
    java -cp "$TLA2TOOLS_JAR" tlc2.TLC -metadir "$state_dir" -config "$cfg" "$model"
    return
  fi
  echo "FAIL: TLC is required for FEATURE-0016 formal checks." >&2
  echo "Install the TLA+ tools and expose 'tlc', or set TLA2TOOLS_JAR to tla2tools.jar." >&2
  exit 1
}

run_tlc "$formal_dir/ExecutionTargetLifecycle.tla" "$formal_dir/ExecutionTargetLifecycle.cfg"
run_tlc "$formal_dir/ExecutionTargetFence.tla" "$formal_dir/ExecutionTargetFence.cfg"
echo "PASS: FEATURE-0016 formal lifecycle and fence/publication models"
