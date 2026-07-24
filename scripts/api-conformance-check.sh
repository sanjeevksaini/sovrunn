#!/usr/bin/env bash
# FEATURE-0012 API conformance check (F12-VERIFY-002).
# Runs fitness functions, schema-diff, coverage, baseline integrity,
# baseline approval, and boundary-ledger checks via internal/apiconform tests.
# Exit 0 on pass; exit 1 on failure.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${REPO_ROOT}"

echo "==> FEATURE-0012 API conformance check"
echo "==> Running: go test ./internal/apiconform/..."

if ! go test ./internal/apiconform/...; then
  echo "FAIL: API conformance checks failed"
  exit 1
fi

echo "PASS: API conformance checks passed"
exit 0
