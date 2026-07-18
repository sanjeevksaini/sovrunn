#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""; TASK=""; MESSAGE=""
while [[ $# -gt 0 ]]; do case "$1" in --feature) FEATURE="$2"; shift 2;; --task) TASK="$2"; shift 2;; --message) MESSAGE="$2"; shift 2;; *) fail "unknown arg: $1";; esac; done
[[ -n "$FEATURE" ]] || fail "--feature required"; [[ -n "$TASK" ]] || fail "--task required"
cd "$(repo_root)"
rm -f sovrunn-api; rm -rf bin
[[ -n "$MESSAGE" ]] || MESSAGE="feat(${FEATURE}): complete task ${TASK}"
[[ -n "$(git status --short)" ]] || fail "nothing to commit"
git add -A
git reset -q docs/generated-prompts .automation/reviews .automation/state .automation/pr || true
git commit -m "$MESSAGE"
./scripts/feature-state.py set --feature "$FEATURE" --key last_committed_task --value "$TASK" >/dev/null
