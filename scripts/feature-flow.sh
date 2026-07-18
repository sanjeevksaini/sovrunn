#!/opt/homebrew/bin/bash
set -euo pipefail

FEATURE=""
SLUG=""
TITLE=""
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
    --slug)
      SLUG="$2"
      shift 2
      ;;
    --title)
      TITLE="$2"
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
[[ -n "$SLUG" ]] || fail "--slug is required"
[[ -n "$TITLE" ]] || fail "--title is required"

echo "============================================================"
echo "==> Starting full Feature Factory flow"
echo "==> Feature: $FEATURE"
echo "==> Slug:    $SLUG"
echo "==> Title:   $TITLE"
echo "============================================================"

echo
echo "==> Step 1/5: ff-start"
make ff-start FEATURE="$FEATURE" SLUG="$SLUG" TITLE="$TITLE"

echo
echo "==> Step 2/5: ff-spec-flow"
make ff-spec-flow FEATURE="$FEATURE"

echo
echo "==> Step 3/5: ff-commit-spec"
make ff-commit-spec FEATURE="$FEATURE"

echo
echo "==> Step 4/5: ff-task-flow"
START_TASK="$START_TASK" make ff-task-flow FEATURE="$FEATURE"

echo
echo "==> Step 5/5: ff-final"
make ff-final FEATURE="$FEATURE"

echo
echo "============================================================"
echo "==> Feature flow completed successfully"
echo "==> Feature: $FEATURE"
echo "==> Branch:  $(git branch --show-current)"
echo "============================================================"
