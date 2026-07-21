#!/usr/bin/env bash
set -euo pipefail

DATA_DIR="${STRUCTURIZR_DATA_DIR:-docs/diagrams/structurizr}"
PORT="${STRUCTURIZR_PORT:-8080}"

if [[ ! -f "$DATA_DIR/workspace.dsl" ]]; then
  echo "ERROR: missing $DATA_DIR/workspace.dsl"
  exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
  echo "ERROR: docker is required to run Structurizr Lite"
  exit 1
fi

echo "Starting Structurizr Lite on http://localhost:$PORT"
echo "Data directory: $DATA_DIR"

docker run -it --rm \
  -p "$PORT:8080" \
  -v "$(pwd)/$DATA_DIR:/usr/local/structurizr" \
  structurizr/lite
