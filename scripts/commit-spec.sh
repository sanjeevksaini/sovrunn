#!/opt/homebrew/bin/bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""
while [[ $# -gt 0 ]]; do case "$1" in --feature) FEATURE="$2"; shift 2;; *) fail "unknown arg: $1";; esac; done
[[ -n "$FEATURE" ]] || fail "--feature required"
cd "$(repo_root)"; ensure_feature_state "$FEATURE"
SPEC_PATH=$(get_feature_value "$FEATURE" spec_path); TITLE=$(get_feature_value "$FEATURE" title)
for f in requirements.md design.md tasks.md; do [[ -f "$SPEC_PATH/$f" ]] || fail "missing $SPEC_PATH/$f"; done
if [[ -f ".automation/features/${FEATURE}.control.json" ]]; then
  ./scripts/spec-approval-check.py --feature "$FEATURE"
fi
git add "$SPEC_PATH/requirements.md" "$SPEC_PATH/design.md" "$SPEC_PATH/tasks.md"
[[ -f "$SPEC_PATH/.config.kiro" ]] && git add "$SPEC_PATH/.config.kiro"
[[ -f ".automation/state/${FEATURE}.json" ]] && git add ".automation/state/${FEATURE}.json"
[[ -f ".automation/features/${FEATURE}.yaml" ]] && git add ".automation/features/${FEATURE}.yaml"
git commit -m "docs(${FEATURE}): add ${TITLE} spec"
