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

step "Registering ServiceClass datastore-postgresql..."
api_call POST "$BASE_URL/v1/service-classes" 201 '{
  "apiVersion":"platform.sovrunn.io/v1alpha1",
  "kind":"ServiceClass",
  "metadata":{"name":"datastore-postgresql","displayName":"PostgreSQL"},
  "spec":{"category":"Database","description":"Managed PostgreSQL datastore","lifecycle":"Active"}
}'

step "Registering ServicePlan postgres-small-ha..."
api_call POST "$BASE_URL/v1/service-plans" 201 '{
  "apiVersion":"platform.sovrunn.io/v1alpha1",
  "kind":"ServicePlan",
  "metadata":{"name":"postgres-small-ha"},
  "spec":{"serviceClassName":"datastore-postgresql","description":"Small HA PostgreSQL plan","tier":"Small","lifecycle":"Active"}
}'

step "Registering Plugin postgres-dstoreops-basic..."
api_call POST "$BASE_URL/v1/plugins" 201 '{
  "apiVersion":"platform.sovrunn.io/v1alpha1",
  "kind":"Plugin",
  "metadata":{"name":"postgres-dstoreops-basic","displayName":"PostgreSQL Basic dStoreOps"},
  "spec":{"pluginType":"dStoreOps","version":"1.0.0","serviceClassRefs":["datastore-postgresql"],"deploymentMode":"compiled-in"}
}'

step "Registering Capability postgres-basic-provision..."
api_call POST "$BASE_URL/v1/capabilities" 201 '{
  "apiVersion":"platform.sovrunn.io/v1alpha1",
  "kind":"Capability",
  "metadata":{"name":"postgres-basic-provision"},
  "spec":{"pluginRef":"postgres-dstoreops-basic","serviceClassRef":"datastore-postgresql","operation":"Provision","supported":true}
}'

step "Registering Capability postgres-basic-bind..."
api_call POST "$BASE_URL/v1/capabilities" 201 '{
  "apiVersion":"platform.sovrunn.io/v1alpha1",
  "kind":"Capability",
  "metadata":{"name":"postgres-basic-bind"},
  "spec":{"pluginRef":"postgres-dstoreops-basic","serviceClassRef":"datastore-postgresql","operation":"Bind","supported":true}
}'

step "Creating ServiceInstance nhm-prod-postgres..."
api_call POST "$BASE_URL/v1/service-instances" 201 '{
  "apiVersion":"platform.sovrunn.io/v1alpha1",
  "kind":"ServiceInstance",
  "metadata":{"name":"nhm-prod-postgres","displayName":"NHM Production PostgreSQL"},
  "spec":{
    "organizationRef":"nic",
    "organizationUnitRef":"ministry-health",
    "tenantRef":"national-health-mission",
    "projectRef":"production",
    "serviceClassRef":"datastore-postgresql",
    "servicePlanRef":"postgres-small-ha",
    "parameters":{"storage":"100Gi"}
  }
}'

step "Creating ServiceBinding nhm-app-postgres-binding..."
api_call POST "$BASE_URL/v1/service-bindings" 201 '{
  "apiVersion":"platform.sovrunn.io/v1alpha1",
  "kind":"ServiceBinding",
  "metadata":{"name":"nhm-app-postgres-binding","displayName":"NHM Application PostgreSQL Binding"},
  "spec":{
    "serviceInstanceRef":"nhm-prod-postgres",
    "consumerRef":{"kind":"Application","name":"nhm-app"},
    "bindingType":"credentials"
  }
}'
