#!/opt/homebrew/bin/bash
set -euo pipefail

fail() { echo "ERROR: $*" >&2; exit 1; }
info() { echo "==> $*"; }
repo_root() { git rev-parse --show-toplevel 2>/dev/null || fail "not inside a git repository"; }
require_clean_tree() {
  if ! git diff --quiet || ! git diff --cached --quiet; then
    fail "working tree has uncommitted tracked changes"
  fi
}
feature_branch_from_id_slug() {
  local feature="$1" slug="$2" lower
  lower=$(echo "$feature" | tr '[:upper:]' '[:lower:]')
  echo "${lower}-${slug}"
}
state_file() { echo ".automation/state/${1}.json"; }
get_feature_value() { ./scripts/feature-state.py get-value --feature "$1" --key "$2"; }
ensure_feature_state() { test -f "$(state_file "$1")" || fail "missing feature state for $1; run feature-start first"; }
