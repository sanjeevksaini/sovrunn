#!/opt/homebrew/bin/bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""; TOOL=""; STAGE=""; TASK=""; SELECTED_MODEL=""; EFFORT=""; FALLBACK_USED="no"; FALLBACK_REASON="none"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --feature) FEATURE="$2"; shift 2;;
    --tool) TOOL="$2"; shift 2;;
    --stage) STAGE="$2"; shift 2;;
    --task) TASK="$2"; shift 2;;
    --selected-model) SELECTED_MODEL="$2"; shift 2;;
    --effort) EFFORT="$2"; shift 2;;
    --fallback-used) FALLBACK_USED="$2"; shift 2;;
    --fallback-reason) FALLBACK_REASON="$2"; shift 2;;
    *) fail "unknown arg: $1";;
  esac
done
[[ -n "$FEATURE" ]] || fail "--feature required"
[[ -n "$TOOL" ]] || fail "--tool required"
[[ -n "$SELECTED_MODEL" ]] || fail "--selected-model required"
cd "$(repo_root)"
mkdir -p ".automation/model-usage/$FEATURE"
OUT=".automation/model-usage/$FEATURE/usage.jsonl"
python3 - "$OUT" "$FEATURE" "$TOOL" "$STAGE" "$TASK" "$SELECTED_MODEL" "$EFFORT" "$FALLBACK_USED" "$FALLBACK_REASON" <<'PY'
import json, sys, datetime
out, feature, tool, stage, task, selected_model, effort, fallback_used, fallback_reason = sys.argv[1:]
row = {
  'ts': datetime.datetime.utcnow().isoformat(timespec='seconds') + 'Z',
  'feature': feature,
  'tool': tool,
  'stage': stage,
  'task': task,
  'selected_model': selected_model,
  'effort': effort,
  'fallback_used': fallback_used,
  'fallback_reason': fallback_reason,
}
with open(out, 'a') as f:
    f.write(json.dumps(row, sort_keys=True) + '\n')
print(out)
PY
