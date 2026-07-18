#!/opt/homebrew/bin/bash
set -euo pipefail

echo "Running go fmt..."
go fmt ./...

echo "Running go test..."
go test ./...

echo "Running go vet..."
go vet ./...

echo "All checks passed."
