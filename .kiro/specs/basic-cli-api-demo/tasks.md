# FEATURE-0010 Basic CLI/API Demo Flow — Tasks

## Task 1: Rewrite demo script with portable shebang, configuration, and helper functions

### Objective

Replace the existing `scripts/demo_phase1.sh` with the hardened script foundation:
portable shebang, `set -euo pipefail`, configurable `BASE_URL`, and the four helper
functions (`step`, `api_call`, `assert_contains`, `assert_not_contains`). No demo
flow steps yet — only the scaffolding.

### Files

- `scripts/demo_phase1.sh` — full rewrite

### Notes

- Shebang: `#!/usr/bin/env bash`
- `set -euo pipefail`
- `BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"`
- `step()` prints `==> $1`
- `api_call METHOD URL EXPECTED_STATUS [DATA]` uses `curl -s -w '\n%{http_code}'`,
  splits response with bash parameter expansion (`${raw##*$'\n'}` and `${raw%$'\n'*}`),
  validates status, prints body, sets global `RESPONSE_BODY`
- `assert_contains BODY SUBSTRING LABEL` uses `[[ "$body" == *"$substring"* ]]`
- `assert_not_contains BODY SUBSTRING LABEL` uses `[[ "$body" != *"$substring"* ]]`
- No external tools (no grep, sed, tail, jq, awk)
- Add a placeholder comment `# === Demo flow begins below ===` at the end

### Tests

- `bash -n scripts/demo_phase1.sh` passes (syntax check)
- Script is executable (`chmod +x`)
- Running script without server prints a connection error and exits non-zero

### Acceptance Criteria

- [x] Shebang is `#!/usr/bin/env bash`
- [x] `set -euo pipefail` is present
- [x] `BASE_URL` is configurable with default `http://127.0.0.1:8080`
- [x] Four helper functions defined: `step`, `api_call`, `assert_contains`, `assert_not_contains`
- [x] `api_call` uses `-w '\n%{http_code}'` and bash parameter expansion only
- [x] No external text-processing tools used
- [x] `bash -n scripts/demo_phase1.sh` passes

### Commit Message

feat(demo): rewrite script foundation with portable shebang and helper functions

## Task 2: Add health and readiness check steps to demo script

### Objective

Add the first two demo flow steps: `/healthz` and `/readyz` checks using the
`api_call` helper. These gate all subsequent steps — if the server is not reachable,
the script fails immediately with a clear message.

### Files

- `scripts/demo_phase1.sh` — add health check steps after helper functions

### Notes

- `step "Checking server health..."` then `api_call GET "$BASE_URL/healthz" 200`
- `step "Checking server readiness..."` then `api_call GET "$BASE_URL/readyz" 200`
- On failure, `api_call` already exits 1 with status mismatch message
- These are the first lines of the demo flow section

### Tests

- Run script without server → exits non-zero at health check step
- Run script with server running → passes health and readiness, then continues

### Acceptance Criteria

- [x] Script checks `/healthz` with expected status 200
- [x] Script checks `/readyz` with expected status 200
- [x] Failure at health check produces clear error and exit 1
- [x] `bash -n scripts/demo_phase1.sh` passes

### Commit Message

feat(demo): add health and readiness check steps

## Task 3: Add Organization hierarchy creation steps

### Objective

Add the four governance resource creation steps: Organization, OrganizationUnit,
Tenant, and Project. Each uses `api_call POST` with expected status 201.

### Files

- `scripts/demo_phase1.sh` — add steps after health checks

### Notes

- Step: Create Organization `nic` → POST `/v1/organizations` → 201
- Step: Create OrganizationUnit `ministry-health` → POST `/v1/organization-units` → 201
- Step: Create Tenant `national-health-mission` → POST `/v1/tenants` → 201
- Step: Create Project `production` → POST `/v1/projects` → 201
- Payloads must match existing API validation (same fields as current script)
- Use `step` label before each `api_call`
- All JSON payloads inline in the script

### Tests

- Start server fresh, run script → all four creates return 201
- Re-run without restart → 409 at Organization create, script exits 1
- `bash -n scripts/demo_phase1.sh` passes

### Acceptance Criteria

- [x] Organization `nic` created with `sovereignLocations`
- [x] OrganizationUnit `ministry-health` created with `organizationRef: "nic"`
- [x] Tenant `national-health-mission` created with all required refs and `isolationProfile`
- [x] Project `production` created with all required refs and `environmentType`
- [x] Each step uses `step` label and `api_call` helper
- [x] Expected status is 201 for all creates

### Commit Message

feat(demo): add organization hierarchy creation steps

## Task 4: Add service catalog registration steps

### Objective

Add ServiceClass and ServicePlan registration steps to the demo script.

### Files

- `scripts/demo_phase1.sh` — add steps after Project creation

### Notes

- Step: Register ServiceClass `datastore.postgresql` → POST `/v1/service-classes` → 201
- Step: Register ServicePlan `postgres-small-ha` → POST `/v1/service-plans` → 201
- ServiceClass payload includes `category`, `description`, `requiredCapabilities`
- ServicePlan payload includes `serviceClassRef: "datastore.postgresql"`
- Payloads match existing validation

### Tests

- Start server fresh, run through to this point → 201 for both
- `bash -n scripts/demo_phase1.sh` passes

### Acceptance Criteria

- [x] ServiceClass `datastore.postgresql` registered with `category: "datastore"`
- [x] ServicePlan `postgres-small-ha` registered with `serviceClassRef`
- [x] Expected status 201 for both
- [x] Step labels printed before each call

### Commit Message

feat(demo): add service catalog registration steps

## Task 5: Add plugin and capability registration steps

### Objective

Add Plugin registration and both Capability registrations (Provision and Bind)
to the demo script.

### Files

- `scripts/demo_phase1.sh` — add steps after ServicePlan creation

### Notes

- Step: Register Plugin `postgres.dstoreops.basic` → POST `/v1/plugins` → 201
- Step: Register Capability `postgres-basic-provision` → POST `/v1/capabilities` → 201
  - `operation: "Provision"`, `pluginRef: "postgres.dstoreops.basic"`
- Step: Register Capability `postgres-basic-bind` → POST `/v1/capabilities` → 201
  - `operation: "Bind"`, `pluginRef: "postgres.dstoreops.basic"`
- Plugin payload includes `pluginType`, `version`, `serviceClassRefs`, `deploymentMode`
- Both capabilities include `serviceClassRef: "datastore.postgresql"` and `supported: true`

### Tests

- Start server fresh, run through to this point → 201 for all three
- `bash -n scripts/demo_phase1.sh` passes

### Acceptance Criteria

- [x] Plugin `postgres.dstoreops.basic` registered
- [x] Capability `postgres-basic-provision` registered with `operation: "Provision"`
- [x] Capability `postgres-basic-bind` registered with `operation: "Bind"`
- [x] Expected status 201 for all three creates
- [x] Step labels printed before each call

### Commit Message

feat(demo): add plugin and capability registration steps

## Task 6: Add service consumption steps (ServiceInstance and ServiceBinding)

### Objective

Add ServiceInstance and ServiceBinding creation steps to the demo script.

### Files

- `scripts/demo_phase1.sh` — add steps after Capability registrations

### Notes

- Step: Create ServiceInstance `nhm-prod-postgres` → POST `/v1/service-instances` → 201
  - Includes all hierarchy refs: `organizationRef`, `organizationUnitRef`, `tenantRef`, `projectRef`
  - Includes `serviceClassRef`, `servicePlanRef`, `parameters`
- Step: Create ServiceBinding `nhm-app-postgres-binding` → POST `/v1/service-bindings` → 201
  - Includes `serviceInstanceRef`, `consumerRef`, `bindingType`

### Tests

- Start server fresh, run full flow through to this point → 201 for both
- `bash -n scripts/demo_phase1.sh` passes

### Acceptance Criteria

- [x] ServiceInstance `nhm-prod-postgres` created with full hierarchy refs
- [x] ServiceBinding `nhm-app-postgres-binding` created with `serviceInstanceRef`
- [x] Expected status 201 for both
- [x] Step labels printed before each call

### Commit Message

feat(demo): add service instance and binding creation steps

## Task 7: Add verification steps (operations, GET by name, success summary)

### Objective

Add the three verification steps and the success summary message at the end
of the demo script.

### Files

- `scripts/demo_phase1.sh` — add verification steps after ServiceBinding creation

### Notes

- Step: List Operations → `api_call GET "$BASE_URL/v1/operations" 200`
  - Then `assert_not_contains "$RESPONSE_BODY" '"items":[]' "No operations found"`
- Step: GET ServiceInstance by name → `api_call GET "$BASE_URL/v1/service-instances/nhm-prod-postgres" 200`
  - Then `assert_contains "$RESPONSE_BODY" "nhm-prod-postgres" "ServiceInstance name not found in response"`
- Step: GET ServiceBinding by name → `api_call GET "$BASE_URL/v1/service-bindings/nhm-app-postgres-binding" 200`
  - Then `assert_contains "$RESPONSE_BODY" "nhm-app-postgres-binding" "ServiceBinding name not found in response"`
- Print success summary: `echo ""; echo "==> Demo completed successfully."`
- Remove the placeholder comment from Task 1 if still present

### Tests

- Start server fresh, run full demo → exits 0 with success summary
- Verify operations list is non-empty
- Verify GET by name returns correct resources
- `bash -n scripts/demo_phase1.sh` passes

### Acceptance Criteria

- [x] Operations list verified non-empty via `assert_not_contains` with `"items":[]`
- [x] ServiceInstance GET-by-name response contains `nhm-prod-postgres`
- [x] ServiceBinding GET-by-name response contains `nhm-app-postgres-binding`
- [x] Success summary printed at the end
- [x] Script exits 0 on full successful run
- [x] No placeholder comments remain

### Commit Message

feat(demo): add verification steps and success summary

## Task 8: Update README.md with "Running the Demo" section

### Objective

Add a "Running the Demo" section to README.md after the "Local Validation" section.
Documents prerequisites, `make demo` usage, re-run instructions, and `BASE_URL` override.

### Files

- `README.md` — add new section

### Notes

- Section title: `## Running the Demo`
- Place after existing "Local Validation" section
- Content per design.md README changes section:
  - Prerequisites: bash 4+, curl, server must be running
  - `make run &` then `sleep 2` then `make demo`
  - Re-run instructions (restart server first)
  - `BASE_URL` override example
- Do not modify any other README sections

### Tests

- `mkdocs build --strict` passes (if MkDocs is configured)
- Section is present after "Local Validation"
- No broken markdown formatting

### Acceptance Criteria

- [x] "Running the Demo" section exists in README.md
- [x] Placed after "Local Validation" section
- [x] Documents bash 4+ requirement
- [x] Documents `make run` prerequisite
- [x] Documents `make demo` command
- [x] Documents re-run with server restart
- [x] Documents `BASE_URL` override
- [x] No other README sections modified

### Commit Message

docs: add "Running the Demo" section to README

## Task 9: End-to-end demo verification and final guardrails

### Objective

Run the full demo end-to-end against a fresh server, verify all acceptance criteria,
run final Docker verification, and ensure clean project state.

### Files

- No file changes expected. This task verifies the work from Tasks 1–8.

### Notes

- Start fresh server: `make run &` then `sleep 2`
- Run demo: `make demo` → must exit 0
- Verify script outputs step labels for all 16 steps
- Verify success summary printed
- Kill server, re-run demo without server → must exit non-zero at health check
- Run syntax check: `bash -n scripts/demo_phase1.sh`
- Run standard Docker verification:
  ```
  docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'
  ```
- Run final Docker verification:
  ```
  docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./... && go build ./cmd/sovrunn-api'
  ```
- Final guardrails:
  - `rm -f sovrunn-api`
  - `rm -rf bin`
  - Verify no `TODO(FEATURE-0010)` under `internal/` or `cmd/`
  - Verify no `internal/api` import of `internal/server`
  - `git status` clean (no untracked artifacts)

### Tests

- Full demo exits 0 against fresh server
- Demo exits non-zero without server
- Docker verification passes (fmt, vet, test, test-race, build)
- No binary artifacts left in repo root or bin/
- No TODO markers for this feature remain

### Acceptance Criteria

- [x] `make demo` exits 0 against fresh server
- [x] `make demo` exits non-zero when server is not running
- [x] `bash -n scripts/demo_phase1.sh` passes
- [x] Standard Docker verification passes
- [x] Final Docker verification passes (includes -race and build)
- [x] `rm -f sovrunn-api` — no binary in repo root
- [x] `rm -rf bin` — no bin directory
- [x] No `TODO(FEATURE-0010)` under `internal/` or `cmd/`
- [x] No `internal/api` import of `internal/server`
- [x] `git status` shows clean working tree
- [x] All 16 demo steps execute with step labels
- [x] Success summary printed at end of demo

### Commit Message

chore(demo): verify FEATURE-0010 end-to-end and clean artifacts

