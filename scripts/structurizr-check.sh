#!/usr/bin/env bash
set -euo pipefail

WORKSPACE="${1:-docs/diagrams/structurizr/workspace.dsl}"

if [[ ! -f "$WORKSPACE" ]]; then
  echo "FAIL: missing Structurizr workspace: $WORKSPACE"
  exit 1
fi

echo "PASS: Structurizr workspace exists: $WORKSPACE"

if command -v structurizr >/dev/null 2>&1; then
  echo "==> Running Structurizr CLI validation"
  structurizr validate -workspace "$WORKSPACE"
  echo "PASS: Structurizr CLI validation"
else
  echo "WARN: Structurizr CLI not installed; syntax validation skipped"
  echo "      Use 'make structurizr-lite' for local visual validation via Docker."
fi
