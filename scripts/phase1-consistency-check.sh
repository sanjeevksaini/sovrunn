#!/opt/homebrew/bin/bash
set -euo pipefail

echo "==> Phase 1 Consistency Check"

resources=(
  organization
  organizationunit
  tenant
  project
  serviceclass
  serviceplan
  plugin
  capability
  serviceinstance
  servicebinding
)

missing=0

check_file() {
  local file="$1"
  if [[ ! -f "$file" ]]; then
    echo "MISSING: $file"
    missing=1
  fi
}

check_one_of() {
  local label="$1"
  shift

  for file in "$@"; do
    if [[ -f "$file" ]]; then
      echo "OK: $label -> $file"
      return 0
    fi
  done

  echo "MISSING: $label"
  printf '  expected one of:\n'
  for file in "$@"; do
    printf '  - %s\n' "$file"
  done

  missing=1
}

warn_one_of() {
  local label="$1"
  shift

  for file in "$@"; do
    if [[ -f "$file" ]]; then
      echo "OK: $label -> $file"
      return 0
    fi
  done

  echo "WARN: missing optional $label"
  printf '  expected one of:\n'
  for file in "$@"; do
    printf '  - %s\n' "$file"
  done
}

echo "==> Checking resource model files"
for r in "${resources[@]}"; do
  check_file "internal/resources/${r}.go"
done

echo "==> Checking validation files"
for r in "${resources[@]}"; do
  check_file "internal/validation/${r}.go"
done

echo "==> Checking API handler/decode files"
for r in "${resources[@]}"; do
  case "$r" in
    organization)
      check_one_of "organization API handler" \
        "internal/api/organization_handler.go" \
        "internal/api/org_handler.go"
      warn_one_of "organization API decode" \
        "internal/api/organization_decode.go" \
        "internal/api/org_decode.go"
      ;;
    organizationunit)
      check_one_of "organizationunit API handler" \
        "internal/api/organizationunit_handler.go" \
        "internal/api/ou_handler.go"
      warn_one_of "organizationunit API decode" \
        "internal/api/organizationunit_decode.go" \
        "internal/api/ou_decode.go"
      ;;
    *)
      check_file "internal/api/${r}_handler.go"
      warn_one_of "${r} API decode" "internal/api/${r}_decode.go"
      ;;
  esac
done

echo "==> Checking registry files"
for r in "${resources[@]}"; do
  case "$r" in
    organization)
      check_one_of "organization registry" \
        "internal/registry/organization_registry.go" \
        "internal/registry/org_registry.go"
      ;;
    organizationunit)
      check_one_of "organizationunit registry" \
        "internal/registry/organizationunit_registry.go" \
        "internal/registry/ou_registry.go"
      ;;
    *)
      check_file "internal/registry/${r}_registry.go"
      ;;
  esac
done

echo "==> Checking gofmt"
unformatted="$(gofmt -l .)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt required:"
  echo "$unformatted"
  exit 1
fi

echo "==> Checking tests"
go test ./...

echo "==> Checking race safety"
go test -race ./...

echo "==> Checking lint"
golangci-lint run ./...

echo "==> Checking security"
gosec ./...

if [[ "$missing" -ne 0 ]]; then
  echo "ERROR: required consistency files missing"
  exit 1
fi

echo "==> Phase 1 Consistency Check passed"
