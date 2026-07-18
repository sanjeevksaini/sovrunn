#!/opt/homebrew/bin/bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""; STAGE=""; MODE="${FEATURE_FACTORY_REVIEW_MODE:-auto}"; MAX_REVISIONS="${FEATURE_FACTORY_MAX_REVISIONS:-3}"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --feature) FEATURE="$2"; shift 2;;
    --stage) STAGE="$2"; shift 2;;
    --mode) MODE="$2"; shift 2;;
    --max-revisions) MAX_REVISIONS="$2"; shift 2;;
    *) fail "unknown arg: $1";;
  esac
done
[[ -n "$FEATURE" ]] || fail "--feature required"
[[ -n "$STAGE" ]] || fail "--stage required"
cd "$(repo_root)"; ensure_feature_state "$FEATURE"
./scripts/reviewer-stage.sh --feature "$FEATURE" --stage "$STAGE" --mode "$MODE"
REVIEW_FILE=".automation/reviews/$FEATURE/${STAGE}.review.json"
[[ -f "$REVIEW_FILE" ]] || fail "missing review file: $REVIEW_FILE"
python3 - "$FEATURE" "$STAGE" "$REVIEW_FILE" "$MAX_REVISIONS" <<'PY'
import json, sys
from pathlib import Path
feature, stage, review_file, max_revisions = sys.argv[1], sys.argv[2], Path(sys.argv[3]), int(sys.argv[4])
review = json.loads(review_file.read_text())
status = review.get('status')
token = review.get('approval_token')
if status == 'APPROVED':
    required = {'requirements':'APPROVED_FOR_DESIGN','design':'APPROVED_FOR_TASKS','tasks':'APPROVED_FOR_CURSOR'}[stage]
    if token != required:
        print(f'ERROR: APPROVED review has wrong token: {token}, required {required}', file=sys.stderr)
        sys.exit(1)
    sys.exit(0)
if status == 'NEEDS_REVISION':
    revision_prompt = (review.get('revision_prompt') or '').strip()
    if not revision_prompt:
        print('ERROR: NEEDS_REVISION review must include revision_prompt', file=sys.stderr)
        sys.exit(1)
    gen_dir = Path(f'.automation/generated-prompts/{feature}')
    gen_dir.mkdir(parents=True, exist_ok=True)
    out = gen_dir / f'{stage}.revision.prompt.md'
    stage_lock = {
      'requirements': 'Revise requirements.md only. Do not generate design.md. Do not generate tasks.md. Do not implement code. Do not modify source files.',
      'design': 'Revise design.md only. Do not generate tasks.md. Do not implement code. Do not modify source files. Do not modify requirements.md unless there is a blocking contradiction.',
      'tasks': 'Revise tasks.md only. Do not implement code. Do not modify source files. Do not modify requirements.md or design.md unless there is a blocking contradiction.'
    }[stage]
    wrapped = f'''{stage_lock}\n\nTool-output safety constraints:\n- Use fs_write only for chunks of 50 lines or fewer.\n- For files longer than 50 lines, use fs_write for the first chunk and fs_append in chunks of 50 lines or fewer.\n- Prefer multiple small str_replace edits instead of one large replacement.\n- Read the file back after editing and verify it is complete.\n\nReviewer summary:\n{review.get('summary','')}\n\nRequired changes:\n'''
    for i, ch in enumerate(review.get('required_changes') or [], 1):
        wrapped += f'{i}. {ch}\n'
    if review.get('blocking_issues'):
        wrapped += '\nBlocking issues to fix:\n'
        for i, issue in enumerate(review.get('blocking_issues') or [], 1):
            wrapped += f'{i}. {issue}\n'
    wrapped += f'''\nRevision instruction from reviewer:\n{revision_prompt}\n\nDo not change the main scope unless the reviewer explicitly requires it.\nDo not move to the next stage until a later review returns the required approval token.\n'''
    out.write_text(wrapped)
    count_file = Path(f'.automation/reviews/{feature}/{stage}.revision-count')
    old = int(count_file.read_text().strip()) if count_file.exists() else 0
    new = old + 1
    count_file.write_text(str(new))
    print(out)
    if new > max_revisions:
        print(f'ERROR: max revision attempts exceeded for {stage}: {new}>{max_revisions}', file=sys.stderr)
        sys.exit(1)
    sys.exit(2)
if status == 'BLOCKED':
    print('ERROR: reviewer returned BLOCKED', file=sys.stderr)
    sys.exit(1)
print(f'ERROR: unknown review status: {status}', file=sys.stderr)
sys.exit(1)
PY
rc=$?
if [[ $rc -eq 0 ]]; then
  ./scripts/approve-stage.sh --feature "$FEATURE" --stage "$STAGE"
  exit 0
elif [[ $rc -eq 2 ]]; then
  ./scripts/feature-state.py set --feature "$FEATURE" --key status --value "${STAGE}_revision_required" >/dev/null
  ./scripts/feature-state.py set --feature "$FEATURE" --key human_gate_required --value true >/dev/null
  info "Revision required. Run the generated revision prompt in Kiro, then rerun this stage review."
  exit 2
else
  ./scripts/feature-state.py set --feature "$FEATURE" --key status --value "${STAGE}_blocked" >/dev/null || true
  exit "$rc"
fi
