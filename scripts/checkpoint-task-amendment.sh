#!/opt/homebrew/bin/bash
# Task-only amendment path for an active implementation checkpoint (ADH-2026-077 §3).
#
# This is the ONLY permitted repair path for uncommitted tasks while a feature
# is in an implementation checkpoint. It runs the deterministic preflight and,
# on success, reports the amended tasks.md hash to bind a fresh task/executable
# -plan approval. It never records approval, never mutates feature state, and
# never triggers a design review.
set -euo pipefail
source "$(dirname "$0")/common.sh"

FEATURE=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --feature) FEATURE="$2"; shift 2;;
    *) fail "unknown arg: $1";;
  esac
done
[[ -n "$FEATURE" ]] || fail "--feature required"
cd "$(repo_root)"
ensure_feature_state "$FEATURE"

info "Validating control manifest for $FEATURE"
./scripts/feature-control.py validate --feature "$FEATURE" >/dev/null

info "Running task-only amendment preflight (ADH-2026-077 §3)"
PYTHONDONTWRITEBYTECODE=1 python3 ./scripts/checkpoint-task-amendment-preflight.py --feature "$FEATURE"

info "Preflight passed. This is a task-only correction; it does not trigger a design review."
info "Record a fresh hash-bound task/executable-plan approval, then resume Cursor at the amended task."
