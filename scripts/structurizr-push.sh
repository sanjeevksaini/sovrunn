#!/usr/bin/env bash
set -euo pipefail

WORKSPACE_ID="${STRUCTURIZR_WORKSPACE_ID:?STRUCTURIZR_WORKSPACE_ID is required}"
API_KEY="${STRUCTURIZR_API_KEY:?STRUCTURIZR_API_KEY is required}"
API_SECRET="${STRUCTURIZR_API_SECRET:?STRUCTURIZR_API_SECRET is required}"
WORKSPACE_FILE="${1:-docs/diagrams/structurizr/workspace.dsl}"

if [[ ! -f "$WORKSPACE_FILE" ]]; then
  echo "ERROR: missing workspace file: $WORKSPACE_FILE"
  exit 1
fi

if ! command -v structurizr >/dev/null 2>&1; then
  echo "ERROR: structurizr CLI is required for push"
  exit 1
fi

structurizr push \
  -id "$WORKSPACE_ID" \
  -key "$API_KEY" \
  -secret "$API_SECRET" \
  -workspace "$WORKSPACE_FILE"
