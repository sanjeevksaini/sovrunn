#!/opt/homebrew/bin/bash
set -euo pipefail

FEATURE=""
TASK=""
MESSAGE=""

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
    --task)
      TASK="$2"
      shift 2
      ;;
    --message)
      MESSAGE="$2"
      shift 2
      ;;
    *)
      fail "unknown arg: $1"
      ;;
  esac
done

[[ -n "$FEATURE" ]] || fail "--feature is required"
[[ -n "$TASK" ]] || fail "--task is required"

if [[ -z "$MESSAGE" ]]; then
  MESSAGE="feat(${FEATURE}): complete task ${TASK}"
fi

# Update task state before staging so the state change is committed with the task.
./scripts/feature-state.py set --feature "$FEATURE" --key last_committed_task --value "$TASK" >/dev/null
./scripts/feature-state.py set --feature "$FEATURE" --key current_task --value "$TASK" >/dev/null
./scripts/feature-state.py set --feature "$FEATURE" --key status --value "cursor_task_${TASK}_committed" >/dev/null

[[ -n "$(git status --short)" ]] || fail "nothing to commit"

git add -A

# Keep generated runtime artifacts out of task commits.
git reset -q docs/generated-prompts .automation/reviews .automation/generated-prompts .automation/logs .automation/model-usage .automation/pr || true

git commit -m "$MESSAGE"
