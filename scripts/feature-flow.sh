#!/opt/homebrew/bin/bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

FEATURE=""
SLUG=""
TITLE=""
PHASE_BRANCH=""
START_TASK="${START_TASK:-}"

usage() {
  cat <<'EOF'
Usage: feature-flow.sh --feature FEATURE-NNNN [options]

Runs the manifest-controlled, resumable feature flow:
  start -> architecture gate -> requirements -> design -> executable-plan gate
  -> tasks -> spec commit -> Cursor tasks -> final verification -> final gate

Options:
  --feature ID          Required feature identifier.
  --slug SLUG           Optional when present in the control manifest.
  --title TITLE         Optional when present in the control manifest.
  --phase-branch NAME   Optional when present in the control manifest.
  --start-task N        Optional first Cursor task; sequence checks still apply.

The command never grants a human approval. At a human gate it exits with status
3, prints the approval command, and can be rerun unchanged after approval.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --feature) FEATURE="$2"; shift 2;;
    --slug) SLUG="$2"; shift 2;;
    --title) TITLE="$2"; shift 2;;
    --phase-branch) PHASE_BRANCH="$2"; shift 2;;
    --start-task) START_TASK="$2"; shift 2;;
    --help|-h) usage; exit 0;;
    *) fail "unknown arg: $1";;
  esac
done

[[ -n "$FEATURE" ]] || fail "--feature is required"
cd "$(repo_root)"

CONTROL_FILE=".automation/features/${FEATURE}.control.json"
[[ -f "$CONTROL_FILE" ]] || fail "full feature flow requires a control manifest: $CONTROL_FILE"

# Static validation is possible before feature-start creates runtime state.
./scripts/feature-control.py validate \
  --feature "$FEATURE" \
  --manifest "$CONTROL_FILE" >/dev/null

manifest_value() {
  python3 - "$CONTROL_FILE" "$1" <<'PY'
import json, sys
data = json.load(open(sys.argv[1]))
print(data["feature"][sys.argv[2]])
PY
}

resolve_identity() {
  local label="$1" supplied="$2" expected="$3"
  if [[ -n "$supplied" && "$supplied" != "$expected" ]]; then
    fail "$label does not match control manifest: supplied '$supplied', expected '$expected'"
  fi
  printf '%s\n' "$expected"
}

SLUG="$(resolve_identity slug "$SLUG" "$(manifest_value slug)")"
TITLE="$(resolve_identity title "$TITLE" "$(manifest_value title)")"
PHASE_BRANCH="$(resolve_identity phase-branch "$PHASE_BRANCH" "$(manifest_value phase_branch)")"
FEATURE_BRANCH="$(manifest_value feature_branch)"
SPEC_PATH="$(manifest_value spec_path)"
STATE_FILE=".automation/state/${FEATURE}.json"

state_value() {
  python3 - "$STATE_FILE" "$1" <<'PY'
import json, sys
print(json.load(open(sys.argv[1])).get(sys.argv[2], ""))
PY
}

stop_for_gate() {
  local gate="$1"
  echo
  echo "==> HUMAN GATE REQUIRED: $gate"
  echo "==> Review the generated evidence, then approve explicitly:"
  echo "    make ff-approve-human-gate FEATURE=$FEATURE GATE=$gate APPROVED_BY='<name>'"
  echo "==> Rerun the same make ff-feature-flow command after approval."
  exit 3
}

echo "============================================================"
echo "==> Manifest-controlled Feature Factory flow"
echo "==> Feature:      $FEATURE"
echo "==> Slug:         $SLUG"
echo "==> Title:        $TITLE"
echo "==> Phase branch: $PHASE_BRANCH"
echo "============================================================"

if [[ ! -f "$STATE_FILE" ]]; then
  echo
  echo "==> Stage 1/8: initialize feature branch and runtime state"
  [[ -z "$(git status --porcelain --untracked-files=all)" ]] || \
    fail "feature initialization requires a clean working tree, including untracked files"
  make ff-start \
    FEATURE="$FEATURE" \
    SLUG="$SLUG" \
    TITLE="$TITLE" \
    PHASE_BRANCH="$PHASE_BRANCH"
  ./scripts/feature-control.py validate --feature "$FEATURE"

  # Commit only deterministic bootstrap artifacts so Kiro begins from a clean,
  # reviewable baseline. The control manifest already exists on the phase branch.
  BOOTSTRAP_PATHS=("$STATE_FILE" "$SPEC_PATH/.config.kiro")
  [[ -f ".automation/features/${FEATURE}.yaml" ]] && \
    BOOTSTRAP_PATHS+=(".automation/features/${FEATURE}.yaml")
  git add "${BOOTSTRAP_PATHS[@]}"
  git commit -m "chore(${FEATURE}): initialize controlled feature flow"
  stop_for_gate architecture
fi

CURRENT_BRANCH="$(git branch --show-current)"
[[ "$CURRENT_BRANCH" == "$FEATURE_BRANCH" ]] || \
  fail "resume must run on $FEATURE_BRANCH; current branch is $CURRENT_BRANCH"
./scripts/feature-control.py validate --feature "$FEATURE"

if [[ "$(state_value final_approval_token)" == "APPROVED_FOR_MERGE" ]]; then
  UNEXPECTED_FINAL_DIRTY="$(git status --porcelain --untracked-files=all | grep -v -F " .automation/state/${FEATURE}.json" || true)"
  [[ -z "$UNEXPECTED_FINAL_DIRTY" ]] || \
    fail "working tree contains changes other than final approval state: $UNEXPECTED_FINAL_DIRTY"
  if [[ -n "$(git status --porcelain -- "$STATE_FILE")" ]]; then
    git add "$STATE_FILE"
    git commit -m "chore(${FEATURE}): record final feature approval"
  fi
  echo "==> $FEATURE is already final-approved and ready for PR preparation."
  exit 0
fi

if [[ "$(state_value architecture_approval_token)" != "APPROVED_ARCHITECTURE" ]]; then
  stop_for_gate architecture
fi

echo
echo "==> Stage 2/8: generate and machine-review requirements, design, and tasks"
if ./scripts/spec-approval-check.py --feature "$FEATURE" >/dev/null 2>&1; then
  echo "==> Approved spec digests are unchanged; skipping Kiro stages."
else
  set +e
  ./scripts/spec-flow.sh \
    --feature "$FEATURE" \
    --mode "${FEATURE_FACTORY_REVIEW_MODE:-auto}" \
    --kiro-mode "${FEATURE_FACTORY_KIRO_MODE:-auto}"
  SPEC_STATUS=$?
  set -e
  if [[ "$SPEC_STATUS" -eq 3 ]]; then
    stop_for_gate executable_plan
  fi
  [[ "$SPEC_STATUS" -eq 0 ]] || exit "$SPEC_STATUS"
fi

echo
echo "==> Stage 3/8: validate approved specification digests and state"
./scripts/spec-approval-check.py --feature "$FEATURE"

spec_is_committed() {
  local path
  for path in requirements.md design.md tasks.md; do
    git ls-files --error-unmatch "$SPEC_PATH/$path" >/dev/null 2>&1 || return 1
    git diff --quiet HEAD -- "$SPEC_PATH/$path" || return 1
  done
}

echo
echo "==> Stage 4/8: commit the approved Kiro specification"
if spec_is_committed; then
  echo "==> Approved specification is already committed; skipping."
else
  make ff-commit-spec FEATURE="$FEATURE"
fi

echo
echo "==> Stage 5/8: execute all remaining Cursor tasks in controlled batches"
RUN_ALL=1 START_TASK="$START_TASK" make ff-controlled-run FEATURE="$FEATURE"

echo
echo "==> Stage 6/8: revalidate approved specification after implementation"
./scripts/spec-approval-check.py --feature "$FEATURE"

echo
echo "==> Stage 7/8: run final repository verification"
make ff-final FEATURE="$FEATURE"

echo
echo "==> Stage 8/8: founder final approval"
stop_for_gate final
