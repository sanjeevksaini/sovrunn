#!/opt/homebrew/bin/bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""
while [[ $# -gt 0 ]]; do case "$1" in --feature) FEATURE="$2"; shift 2;; *) fail "unknown arg: $1";; esac; done
[[ -n "$FEATURE" ]] || fail "--feature required"
cd "$(repo_root)"; ensure_feature_state "$FEATURE"
PHASE_BRANCH=$(get_feature_value "$FEATURE" phase_branch); FEATURE_BRANCH=$(get_feature_value "$FEATURE" feature_branch); TITLE=$(get_feature_value "$FEATURE" title); SPEC_PATH=$(get_feature_value "$FEATURE" spec_path)
case "$PHASE_BRANCH" in
  phase1-*) PHASE_LABEL="Phase 1";;
  phase2-*) PHASE_LABEL="Phase 2";;
  *) PHASE_LABEL="phase branch ${PHASE_BRANCH}";;
esac
GO_DOCKER_IMAGE="${GO_DOCKER_IMAGE:-golang:1.22}"
PR_DIR=".automation/pr"; PR_BODY="$PR_DIR/${FEATURE}.md"; mkdir -p "$PR_DIR"
require_clean_tree
cat > "$PR_BODY" <<PRBODY
Implements ${TITLE} for Sovrunn ${PHASE_LABEL}.

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

All verification passed in Docker using ${GO_DOCKER_IMAGE}.
PRBODY

git push -u origin "$FEATURE_BRANCH"
gh pr create --base "$PHASE_BRANCH" --head "$FEATURE_BRANCH" --title "${FEATURE}: ${TITLE}" --body-file "$PR_BODY"
