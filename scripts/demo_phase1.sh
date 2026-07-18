#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"
RESPONSE_BODY=""

step() {
  echo "==> $1"
}

api_call() {
  local method="$1"
  local url="$2"
  local expected_status="$3"
  local data="${4:-}"
  local raw
  local http_code
  local body

  if [[ -n "$data" ]]; then
    if ! raw=$(curl -s -w '\n%{http_code}' -X "$method" "$url" \
      -H "Content-Type: application/json" \
      -d "$data"); then
      echo "Connection error: unable to reach $url" >&2
      exit 1
    fi
  else
    if ! raw=$(curl -s -w '\n%{http_code}' -X "$method" "$url" \
      -H "Content-Type: application/json"); then
      echo "Connection error: unable to reach $url" >&2
      exit 1
    fi
  fi

  http_code="${raw##*$'\n'}"
  body="${raw%$'\n'*}"

  if [[ "$http_code" != "$expected_status" ]]; then
    echo "API call failed: $method $url" >&2
    echo "Expected HTTP status $expected_status, got $http_code" >&2
    echo "$body" >&2
    exit 1
  fi

  RESPONSE_BODY="$body"
  echo "$body"
}

assert_contains() {
  local body="$1"
  local substring="$2"
  local label="$3"

  if [[ "$body" == *"$substring"* ]]; then
    return 0
  fi

  echo "FAIL: $label — expected substring '$substring' not found" >&2
  exit 1
}

assert_not_contains() {
  local body="$1"
  local substring="$2"
  local label="$3"

  if [[ "$body" != *"$substring"* ]]; then
    return 0
  fi

  echo "FAIL: $label — unexpected substring '$substring' found" >&2
  exit 1
}

# === Demo flow begins below ===

step "Checking server health..."
api_call GET "$BASE_URL/healthz" 200

step "Checking server readiness..."
api_call GET "$BASE_URL/readyz" 200

step "Creating Organization nic..."
api_call POST "$BASE_URL/v1/organizations" 201 '{
  "apiVersion":"platform.sovrunn.io/v1alpha1",
  "kind":"Organization",
  "metadata":{"name":"nic","displayName":"National Informatics Centre"},
  "spec":{"description":"Central government cloud organization","sovereignLocations":["in-delhi-1","in-mumbai-1"]}
}'

step "Creating OrganizationUnit ministry-health..."
api_call POST "$BASE_URL/v1/organization-units" 201 '{
  "apiVersion":"platform.sovrunn.io/v1alpha1",
  "kind":"OrganizationUnit",
  "metadata":{"name":"ministry-health","displayName":"Ministry of Health"},
  "spec":{"organizationName":"nic","description":"Health ministry OU"}
}'

step "Creating Tenant national-health-mission..."
api_call POST "$BASE_URL/v1/tenants" 201 '{
  "apiVersion":"platform.sovrunn.io/v1alpha1",
  "kind":"Tenant",
  "metadata":{"name":"national-health-mission","displayName":"National Health Mission"},
  "spec":{"organizationName":"nic","organizationUnitName":"ministry-health"}
}'

step "Creating Project production..."
api_call POST "$BASE_URL/v1/projects" 201 '{
  "apiVersion":"platform.sovrunn.io/v1alpha1",
  "kind":"Project",
  "metadata":{"name":"production","displayName":"Production"},
  "spec":{"organizationName":"nic","organizationUnitName":"ministry-health","tenantName":"national-health-mission"}
}'
