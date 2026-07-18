#!/opt/homebrew/bin/bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""
while [[ $# -gt 0 ]]; do case "$1" in --feature) FEATURE="$2"; shift 2;; *) fail "unknown arg: $1";; esac; done
[[ -n "$FEATURE" ]] || fail "--feature required"
cd "$(repo_root)"
info "Checking no TODO markers for $FEATURE"
if grep -R "TODO(${FEATURE})" internal cmd 2>/dev/null; then fail "found TODO(${FEATURE}) marker"; fi
info "Checking internal/api does not import internal/server"
if grep -R "internal/server" internal/api 2>/dev/null; then fail "internal/api imports internal/server"; fi
info "Checking build artifacts are absent"
if [[ -e sovrunn-api || -d bin ]]; then fail "build artifact present: remove sovrunn-api and/or bin/"; fi
info "Guardrails passed"
