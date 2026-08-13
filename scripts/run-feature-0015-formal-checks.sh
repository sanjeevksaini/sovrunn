#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
formal_dir="$root/docs/formal/feature-0015"

run_tlc() {
  local model="$1"
  local cfg="$2"
  if command -v tlc >/dev/null 2>&1; then
    tlc -config "$cfg" "$model"
    return
  fi
  if [[ -n "${TLA2TOOLS_JAR:-}" && -f "${TLA2TOOLS_JAR}" ]]; then
    java -cp "$TLA2TOOLS_JAR" tlc2.TLC -config "$cfg" "$model"
    return
  fi
  echo "FAIL: TLC is required for FEATURE-0015 formal checks." >&2
  echo "Install the TLA+ tools and expose 'tlc', or set TLA2TOOLS_JAR to tla2tools.jar." >&2
  exit 1
}

run_tlc "$formal_dir/ParticipationLifecycle.tla" "$formal_dir/ParticipationLifecycle.cfg"
run_tlc "$formal_dir/IdempotencyAuditPublication.tla" "$formal_dir/IdempotencyAuditPublication.cfg"
echo "PASS: FEATURE-0015 formal lifecycle and idempotency/audit models"
