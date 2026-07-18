#!/opt/homebrew/bin/bash
set -euo pipefail

FEATURE="${1:-}"

if [[ -z "${FEATURE}" ]]; then
  echo "usage: ./scripts/verify-feature.sh FEATURE-0001"
  exit 1
fi

echo "Verifying ${FEATURE}"

echo "Current branch:"
git branch --show-current

echo
echo "Git status:"
git status --short

echo
echo "Running format..."
go fmt ./...

echo
echo "Running tests..."
go test ./...

echo
echo "Running vet..."
go vet ./...

echo
echo "Verification complete for ${FEATURE}"
