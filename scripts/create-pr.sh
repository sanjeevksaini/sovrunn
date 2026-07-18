#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""
while [[ $# -gt 0 ]]; do case "$1" in --feature) FEATURE="$2"; shift 2;; *) fail "unknown arg: $1";; esac; done
[[ -n "$FEATURE" ]] || fail "--feature required"
cd "$(repo_root)"; ensure_feature_state "$FEATURE"
PHASE_BRANCH=$(get_feature_value "$FEATURE" phase_branch); FEATURE_BRANCH=$(get_feature_value "$FEATURE" feature_branch); TITLE=$(get_feature_value "$FEATURE" title); SPEC_PATH=$(get_feature_value "$FEATURE" spec_path)
PR_DIR=".automation/pr"; PR_BODY="$PR_DIR/${FEATURE}.md"; mkdir -p "$PR_DIR"
require_clean_tree
cat > "$PR_BODY" <<PRBODY
Implements ${TITLE} for Sovrunn Phase 1.

Spec files:
- ${SPEC_PATH}/requirements.md
- ${SPEC_PATH}/design.md
- ${SPEC_PATH}/tasks.md

Verification:
- gofmt clean
- go vet ./...
- go test ./...
- go test -race ./...
- go build ./cmd/sovrunn-api

All verification passed in Docker using golang:1.21.
PRBODY

git push -u origin "$FEATURE_BRANCH"
gh pr create --base "$PHASE_BRANCH" --head "$FEATURE_BRANCH" --title "${FEATURE}: ${TITLE}" --body-file "$PR_BODY"
