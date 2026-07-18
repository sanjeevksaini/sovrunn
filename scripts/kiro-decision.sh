#!/opt/homebrew/bin/bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""; STAGE=""; QUESTION=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --feature) FEATURE="$2"; shift 2;;
    --stage) STAGE="$2"; shift 2;;
    --question) QUESTION="$2"; shift 2;;
    *) fail "unknown arg: $1";;
  esac
done
[[ -n "$FEATURE" ]] || fail "--feature required"
[[ -n "$STAGE" ]] || STAGE="unknown"
if [[ -z "$QUESTION" ]]; then
  QUESTION="$(cat)"
fi
python3 - "$FEATURE" "$STAGE" "$QUESTION" <<'PYDECISION'
import json, sys
feature, stage, question = sys.argv[1:]
q = question.lower()

def contains_any(words):
    return any(w.lower() in q for w in words)

action = "pause"
answer = "Pause"
reason = "This decision is not covered by the safe auto-accept/auto-reject rules."
instruction = "Pause for architecture review before answering Kiro."

if contains_any(["implement", "write code", "modify source", "create source files"]):
    action = "reject"
    answer = "Reject"
    reason = "Kiro is in the spec stage. Implementation is handled later by Cursor task-by-task."
    instruction = "Do not implement code. Generate or revise the current spec document only."
elif contains_any(["ServiceInstance", "ServiceBinding", "provisioning", "persistence", "auth", "RBAC", "UI", "billing", "marketplace", "ServiceOps execution"]):
    action = "reject"
    answer = "Reject"
    reason = "This looks like scope expansion beyond the current feature unless explicitly listed in the approved requirements."
    instruction = "Keep this as a non-goal or future feature unless the approved current feature scope explicitly includes it."
elif contains_any(["build a feature", "build feature"]):
    action = "accept"
    answer = "Accept"
    reason = "Sovrunn features use Kiro Build a Feature workflow."
    instruction = "Use Build a Feature."
elif contains_any(["requirements.md", "generate requirements"]):
    action = "accept" if stage in {"requirements", "requirements_pending", "unknown"} else "pause"
    answer = "Accept" if action == "accept" else "Pause"
    reason = "Requirements generation is the first spec stage." if action == "accept" else "Requirements generation was requested outside the expected stage."
    instruction = "Generate requirements.md only. Do not generate design.md, tasks.md, or code." if action == "accept" else "Pause for stage-state review."
elif contains_any(["design.md", "generate design"]):
    action = "pause"
    answer = "Pause"
    reason = "Design generation requires APPROVED_FOR_DESIGN token from reviewer."
    instruction = "Continue only if the stage review JSON has status=APPROVED and approval_token=APPROVED_FOR_DESIGN."
elif contains_any(["tasks.md", "generate tasks"]):
    action = "pause"
    answer = "Pause"
    reason = "Tasks generation requires APPROVED_FOR_TASKS token from reviewer."
    instruction = "Continue only if the stage review JSON has status=APPROVED and approval_token=APPROVED_FOR_TASKS."

print(json.dumps({
    "feature": feature,
    "stage": stage,
    "action": action,
    "decision_response": f"Decision: {question}\n\nAnswer: {answer}\n\nReason:\n{reason}\n\nInstruction to Kiro:\n{instruction}",
}, indent=2))
PYDECISION
