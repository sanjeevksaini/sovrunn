#!/opt/homebrew/bin/bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""; SLUG=""; TITLE=""; PHASE_BRANCH="phase1-foundation"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --feature) FEATURE="$2"; shift 2;; --slug) SLUG="$2"; shift 2;; --title) TITLE="$2"; shift 2;; --phase-branch) PHASE_BRANCH="$2"; shift 2;; *) fail "unknown arg: $1";;
  esac
done
[[ -n "$FEATURE" ]] || fail "--feature required"; [[ -n "$SLUG" ]] || fail "--slug required"; [[ -n "$TITLE" ]] || fail "--title required"
cd "$(repo_root)"; require_clean_tree
BRANCH=$(feature_branch_from_id_slug "$FEATURE" "$SLUG")
info "Starting $FEATURE: $TITLE"; info "Base branch: $PHASE_BRANCH"; info "Feature branch: $BRANCH"
git checkout "$PHASE_BRANCH"; git pull origin "$PHASE_BRANCH"; git checkout -b "$BRANCH"
mkdir -p ".kiro/specs/$SLUG" ".automation/state" ".automation/features" "docs/generated-prompts/$FEATURE" ".automation/reviews/$FEATURE" ".automation/pr"
./scripts/feature-state.py init --feature "$FEATURE" --slug "$SLUG" --title "$TITLE" --phase-branch "$PHASE_BRANCH" --branch "$BRANCH"
cat > ".kiro/specs/$SLUG/.config.kiro" <<CFG
feature_id=$FEATURE
slug=$SLUG
title=$TITLE
CFG
info "Feature started. Next: make ff-prompt-requirements FEATURE=$FEATURE"
