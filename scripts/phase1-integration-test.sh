#!/opt/homebrew/bin/bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"
CONFIG_FILE="${CONFIG_FILE:-configs/sovrunn-api.local.yaml}"
LOG_FILE="${LOG_FILE:-/tmp/sovrunn-phase1-integration.log}"

echo "==> Phase 1 Integration Test"
echo "==> BASE_URL=$BASE_URL"
echo "==> CONFIG_FILE=$CONFIG_FILE"
echo "==> LOG_FILE=$LOG_FILE"

cleanup() {
  if [[ -n "${SERVER_PID:-}" ]]; then
    echo "==> Stopping server PID $SERVER_PID"
    kill "$SERVER_PID" >/dev/null 2>&1 || true
    wait "$SERVER_PID" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

echo "==> Starting sovrunn-api"
go run ./cmd/sovrunn-api --config "$CONFIG_FILE" >"$LOG_FILE" 2>&1 &
SERVER_PID="$!"

echo "==> Waiting for server"
for i in {1..30}; do
  if curl -fsS "$BASE_URL/healthz" >/dev/null 2>&1; then
    echo "==> Server healthz is reachable"
    break
  fi

  if ! kill -0 "$SERVER_PID" >/dev/null 2>&1; then
    echo "ERROR: server exited early"
    cat "$LOG_FILE"
    exit 1
  fi

  sleep 1

  if [[ "$i" -eq 30 ]]; then
    echo "ERROR: server did not become healthy"
    cat "$LOG_FILE"
    exit 1
  fi
done

echo "==> Checking healthz"
curl -fsS "$BASE_URL/healthz" >/dev/null

echo "==> Checking readyz"
curl -fsS "$BASE_URL/readyz" >/dev/null

echo "==> Running Phase 1 demo"
BASE_URL="$BASE_URL" ./scripts/demo_phase1.sh

echo "==> Verifying list endpoints"
for path in \
  /v1/organizations \
  /v1/organization-units \
  /v1/tenants \
  /v1/projects \
  /v1/service-classes \
  /v1/service-plans \
  /v1/plugins \
  /v1/capabilities \
  /v1/service-instances \
  /v1/service-bindings \
  /v1/operations
do
  echo "==> GET $path"
  curl -fsS "$BASE_URL$path" >/dev/null
done

echo "==> Verifying invalid request fails"
status="$(
  curl -sS -o /tmp/sovrunn-invalid-response.json \
    -w "%{http_code}" \
    -H "Content-Type: application/json" \
    -X POST \
    "$BASE_URL/v1/organizations" \
    -d '{}'
)"

case "$status" in
  400|422)
    echo "==> Invalid request correctly failed with HTTP $status"
    ;;
  *)
    echo "ERROR: expected invalid request to fail with 400 or 422, got $status"
    cat /tmp/sovrunn-invalid-response.json || true
    exit 1
    ;;
esac

echo "==> Phase 1 Integration Test passed"
