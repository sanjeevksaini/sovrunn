#!/opt/homebrew/bin/bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""
while [[ $# -gt 0 ]]; do case "$1" in --feature) FEATURE="$2"; shift 2;; *) fail "unknown arg: $1";; esac; done
[[ -n "$FEATURE" ]] || fail "--feature required"
cd "$(repo_root)"; ensure_feature_state "$FEATURE"
PHASE_BRANCH=$(get_feature_value "$FEATURE" phase_branch); FEATURE_BRANCH=$(get_feature_value "$FEATURE" feature_branch); TITLE=$(get_feature_value "$FEATURE" title); SPEC_PATH=$(get_feature_value "$FEATURE" spec_path)
CONTROL_FILE=".automation/features/${FEATURE}.control.json"
if [[ -f "$CONTROL_FILE" ]]; then
  python3 - "$FEATURE" <<'PY'
import json, subprocess, sys
from pathlib import Path
feature = sys.argv[1]
state = json.loads(Path(f'.automation/state/{feature}.json').read_text())
control = Path(f'.automation/features/{feature}.control.json')
import hashlib
if state.get('control_approved_sha256') != hashlib.sha256(control.read_bytes()).hexdigest():
    raise SystemExit('ERROR: feature control differs from the founder-approved digest')
if state.get('final_approval_token') != 'APPROVED_FOR_MERGE':
    raise SystemExit('ERROR: explicit final human approval is required before PR creation')
verified = state.get('final_approved_commit', '')
if not verified:
    raise SystemExit('ERROR: final approved commit is missing')
if subprocess.run(['git', 'merge-base', '--is-ancestor', verified, 'HEAD']).returncode:
    raise SystemExit('ERROR: final approved commit is not an ancestor of HEAD')
changed = subprocess.check_output(['git', 'diff', '--name-only', f'{verified}..HEAD'], text=True).splitlines()
allowed = {f'.automation/state/{feature}.json'}
unexpected = [path for path in changed if path not in allowed]
if unexpected:
    raise SystemExit('ERROR: non-state changes occurred after final approval: ' + ', '.join(unexpected))
PY
fi
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
