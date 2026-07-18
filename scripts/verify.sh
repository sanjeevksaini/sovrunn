#!/opt/homebrew/bin/bash
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
echo "==> Running Docker Go verification"
docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'
