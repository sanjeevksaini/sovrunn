#!/opt/homebrew/bin/bash
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
GO_DOCKER_IMAGE="${GO_DOCKER_IMAGE:-golang:1.22}"
echo "==> Running Docker Go verification with ${GO_DOCKER_IMAGE}"
docker run --rm -v "$PWD":/src -w /src "$GO_DOCKER_IMAGE" sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'
