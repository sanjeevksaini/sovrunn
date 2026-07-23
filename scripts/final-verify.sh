#!/opt/homebrew/bin/bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""
while [[ $# -gt 0 ]]; do case "$1" in --feature) FEATURE="$2"; shift 2;; *) fail "unknown arg: $1";; esac; done
[[ -n "$FEATURE" ]] || fail "--feature required"
cd "$(repo_root)"
GO_DOCKER_IMAGE="${GO_DOCKER_IMAGE:-golang:1.22}"
echo "==> Running final Docker verification with ${GO_DOCKER_IMAGE}"
docker run --rm -v "$PWD":/src -w /src "$GO_DOCKER_IMAGE" sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./... && go build ./cmd/sovrunn-api'
rm -f sovrunn-api; rm -rf bin
./scripts/guardrails.sh --feature "$FEATURE"
git status --short
[[ -z "$(git status --short)" ]] || fail "working tree is not clean after final verification"
echo "==> Final verification passed"
