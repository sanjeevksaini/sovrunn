#!/usr/bin/env bash
set -euo pipefail

CONFIG="${CONFIG:-configs/sovrunn-api.local.yaml}"

echo "Starting Sovrunn API with config: ${CONFIG}"
go run ./cmd/sovrunn-api --config "${CONFIG}"
