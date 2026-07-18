#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""; STAGE=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --feature) FEATURE="$2"; shift 2;;
    --stage) STAGE="$2"; shift 2;;
    *) fail "unknown arg: $1";;
  esac
done
[[ -n "$FEATURE" ]] || fail "--feature required"
[[ -n "$STAGE" ]] || fail "--stage required"
cd "$(repo_root)"; ensure_feature_state "$FEATURE"
case "$STAGE" in
  requirements) REQUIRED_TOKEN="APPROVED_FOR_DESIGN"; NEXT_STAGE="design"; NEXT_STATUS="requirements_approved";;
  design) REQUIRED_TOKEN="APPROVED_FOR_TASKS"; NEXT_STAGE="tasks"; NEXT_STATUS="design_approved";;
  tasks) REQUIRED_TOKEN="APPROVED_FOR_CURSOR"; NEXT_STAGE="cursor"; NEXT_STATUS="tasks_approved";;
  *) fail "stage must be requirements, design, or tasks";;
esac
REVIEW_FILE=".automation/reviews/$FEATURE/${STAGE}.review.json"
[[ -f "$REVIEW_FILE" ]] || fail "missing review file: $REVIEW_FILE"
python3 - "$REVIEW_FILE" "$REQUIRED_TOKEN" <<'PYAPPROVE'
import json, sys
path, required = sys.argv[1:]
review = json.load(open(path))
status = review.get('status')
token = review.get('approval_token')
issues = review.get('blocking_issues') or []
changes = review.get('required_changes') or []
if status != 'APPROVED' or token != required:
    print(f'ERROR: stage is not approved: status={status} token={token} required={required}', file=sys.stderr)
    sys.exit(1)
if issues or changes:
    print('ERROR: approved review must not include blocking_issues or required_changes', file=sys.stderr)
    sys.exit(1)
print(f'approved: {token}')
PYAPPROVE
./scripts/feature-state.py set --feature "$FEATURE" --key status --value "$NEXT_STATUS"
./scripts/feature-state.py set --feature "$FEATURE" --key current_stage --value "$NEXT_STAGE"
./scripts/feature-state.py set --feature "$FEATURE" --key human_gate_required --value false
info "$STAGE approved with token $REQUIRED_TOKEN; next stage: $NEXT_STAGE"
