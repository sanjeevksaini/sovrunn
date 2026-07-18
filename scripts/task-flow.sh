#!/opt/homebrew/bin/bash
set -euo pipefail

FEATURE=""
START_TASK="${START_TASK:-1}"

fail() {
  echo "ERROR: $*" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --feature)
      FEATURE="$2"
      shift 2
      ;;
    --start-task)
      START_TASK="$2"
      shift 2
      ;;
    *)
      fail "unknown arg: $1"
      ;;
  esac
done

[[ -n "$FEATURE" ]] || fail "--feature is required"

FEATURE_FILE=".automation/features/${FEATURE}.yaml"
STATE_FILE=".automation/state/${FEATURE}.json"

SPEC_SLUG=""

if [[ -f "$FEATURE_FILE" ]]; then
  SPEC_SLUG="$(grep -E '^slug:' "$FEATURE_FILE" | head -n 1 | sed 's/^slug:[[:space:]]*//' | tr -d '"'\''')"
fi

if [[ -z "$SPEC_SLUG" && -f "$STATE_FILE" ]]; then
  SPEC_SLUG="$(python3 - <<PY
import json
from pathlib import Path

p = Path("$STATE_FILE")
data = json.loads(p.read_text())
print(data.get("slug") or data.get("spec_slug") or data.get("feature_slug") or "")
PY
)"
fi

if [[ -z "$SPEC_SLUG" ]]; then
  fail "could not resolve spec slug for $FEATURE from $FEATURE_FILE or $STATE_FILE"
fi

TASKS_FILE=".kiro/specs/${SPEC_SLUG}/tasks.md"

[[ -f "$TASKS_FILE" ]] || fail "tasks.md not found: $TASKS_FILE"

TASK_IDS="$(awk '
  /^[[:space:]]*- \\[[ xX]\\][[:space:]]+[0-9]+(\\.|[[:space:]])/ {
    line=$0
    sub(/^[[:space:]]*- \\[[ xX]\\][[:space:]]+/, "", line)
    sub(/[^0-9].*$/, "", line)
    if (line != "" && !seen[line]++) print line
  }
  /^[[:space:]]*#{1,6}[[:space:]]*Task[[:space:]]+[0-9]+[:.]?/ {
    line=$0
    sub(/^[[:space:]]*#{1,6}[[:space:]]*Task[[:space:]]+/, "", line)
    sub(/[^0-9].*$/, "", line)
    if (line != "" && !seen[line]++) print line
  }
' "$TASKS_FILE")"

[[ -n "$TASK_IDS" ]] || fail "could not detect numbered tasks in $TASKS_FILE"

echo "==> Feature: $FEATURE"
echo "==> Spec slug: $SPEC_SLUG"
echo "==> Tasks file: $TASKS_FILE"
echo "==> Starting from task: $START_TASK"

for TASK in $TASK_IDS; do
  if (( TASK < START_TASK )); then
    continue
  fi

  echo
  echo "============================================================"
  echo "==> Running Cursor task $TASK"
  echo "============================================================"

  make ff-cursor-task FEATURE="$FEATURE" TASK="$TASK"
  make ff-verify FEATURE="$FEATURE"
  make ff-commit-task FEATURE="$FEATURE" TASK="$TASK"
done

echo
echo "==> All detected Cursor tasks completed."
