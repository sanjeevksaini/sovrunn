# FEATURE-0010: Basic CLI/API Demo Flow — Design

## Overview

This design specifies the hardening and completion of `scripts/demo_phase1.sh` as the
repeatable, end-to-end Phase 1 demo for the Sovrunn platform. The existing script already
contains curl-based calls for all Phase 1 resources. This feature adds:

- Portable shebang (`#!/usr/bin/env bash`) for macOS and Linux compatibility
- Configurable `BASE_URL` with sensible default (`http://127.0.0.1:8080`)
- Step-level human-readable labels printed before each API call
- Fail-fast behavior with clear error output on any unexpected HTTP status
- HTTP status code validation per step (expected 200/201)
- GET verification of created ServiceInstance and ServiceBinding by name
- Operation list verification (non-empty assertion)
- A second Capability registration (`postgres-basic-bind`) per the demo flow spec
- Non-zero exit code on any failure
- `make demo` target integration (already exists, no Makefile changes needed)
- README update referencing the demo

No Go code changes are required. No new packages or dependencies are introduced.
The demo script uses bash built-ins and curl only — no grep, sed, tail, jq, or other
external text-processing tools.

---

## Resolved Design Decisions

### 1. Verification depth (requirements Q1)

**Answer**: The script validates HTTP status codes for all mutating calls AND validates
response body content for the two GET-by-name verification steps (ServiceInstance and
ServiceBinding). This balances robustness against script complexity.

Implementation: Each curl call captures the HTTP status code using `-w '\n%{http_code}'`
and splits the response using bash parameter expansion (see Helper Function Design below).
The expected value is 201 for creates and 200 for GETs/lists. The two GET verification
steps additionally use bash substring matching (`[[ "$body" == *"pattern"* ]]`) to confirm
the exact resource name string appears in the response body:
- ServiceInstance verification asserts the literal string `nhm-prod-postgres` is present.
- ServiceBinding verification asserts the literal string `nhm-app-postgres-binding` is present.

No `jq` dependency is introduced. No external tools (grep, sed, tail) are used.

### 2. Second Capability registration (requirements Q2)

**Answer**: Both capabilities are demonstrated. The feature file (section 2, step 10 and 11)
explicitly lists "Register Capability Provision" and "Register Capability Bind" as separate
steps. The demo registers:
- `postgres-basic-provision` (operation: Provision)
- `postgres-basic-bind` (operation: Bind)

This demonstrates that a plugin can declare multiple lifecycle capabilities.

### 3. List filtering (requirements Q3)

**Answer**: The demo uses unfiltered list endpoints only. The Phase 1 API contract does not
define query-parameter-based filtering for list endpoints. The demo calls:
- `GET /v1/operations` (unfiltered)
- `GET /v1/service-instances/{name}` (get by name, not a filtered list)
- `GET /v1/service-bindings/{name}` (get by name, not a filtered list)

Filtered list calls are deferred until filtering is explicitly implemented and documented.

### 4. Output format (requirements Q4)

**Answer**: Raw curl output. No `jq` dependency. The script prints raw JSON response bodies
to stdout. This preserves zero-dependency portability (bash + curl only). Human readability
is achieved through step labels printed before each call. For the operations check, the
script prints the raw response body only — no operation count is extracted or printed.

### 5. Server lifecycle (requirements Q5)

**Answer**: Server management remains manual. The user starts the server separately before
running the demo. The script validates server availability via `/healthz` and `/readyz`
checks at the top. If the server is not running, the script fails immediately with a clear
error message. This keeps the script simple and avoids background process management
complexity.

### 6. Shebang portability (requirements Q6)

**Answer**: Change to `#!/usr/bin/env bash`. This is the standard portable shebang for bash
scripts. The existing `#!/opt/homebrew/bin/bash` is macOS Homebrew-specific and fails on
Linux or systems without Homebrew.

### 7. Cleanup step (requirements Q7)

**Answer**: No cleanup step. The requirements non-goals explicitly state "Idempotency via
update-or-create logic (409 on re-run is acceptable)" and "Testing of DELETE endpoints in
the demo flow". Users restart the server to reset in-memory state. This avoids adding
DELETE calls that are not part of the demo's educational purpose.

### 8. README integration (requirements Q8)

**Answer**: Add a new "## Running the Demo" section in README.md, placed after the existing
"Local Validation" section. This section documents:
- Prerequisites (server must be running)
- The `make demo` command
- Expected behavior
- How to reset state (restart the server)

---

## Architecture

### Component Interaction

```text
User
  │
  └─ make demo
       │
       └─ scripts/demo_phase1.sh
            │
            ├─ curl GET /healthz          → 200 ok
            ├─ curl GET /readyz           → 200 ready
            ├─ curl POST /v1/organizations       → 201 Created
            ├─ curl POST /v1/organization-units  → 201 Created
            ├─ curl POST /v1/tenants             → 201 Created
            ├─ curl POST /v1/projects            → 201 Created
            ├─ curl POST /v1/service-classes     → 201 Created
            ├─ curl POST /v1/service-plans       → 201 Created
            ├─ curl POST /v1/plugins             → 201 Created
            ├─ curl POST /v1/capabilities        → 201 Created (Provision)
            ├─ curl POST /v1/capabilities        → 201 Created (Bind)
            ├─ curl POST /v1/service-instances   → 201 Created
            ├─ curl POST /v1/service-bindings    → 201 Created
            ├─ curl GET /v1/operations           → 200 (verify non-empty)
            ├─ curl GET /v1/service-instances/nhm-prod-postgres   → 200 (verify name)
            └─ curl GET /v1/service-bindings/nhm-app-postgres-binding → 200 (verify name)
```

### System Boundaries

```text
scripts/demo_phase1.sh (bash)
  └─ communicates with: sovrunn-api HTTP server (started separately)
       └─ uses: in-memory registry (no persistence)
```

The demo script is external to the Go codebase. It exercises the public HTTP API only.
No Go code changes are required for this feature.

### Package Responsibilities

| Component | Role in FEATURE-0010 |
|-----------|---------------------|
| `scripts/demo_phase1.sh` | Demo script — exercised and hardened |
| `Makefile` | `demo` target invokes the script (already exists) |
| `README.md` | Documentation update referencing the demo |

No Go packages are modified. No new Go code is written.

---

## Files Changed

| File | Action | Purpose |
|------|--------|---------|
| `scripts/demo_phase1.sh` | Modify | Harden script: portable shebang, status validation, step labels, GET verification via bash substring matching, second capability, operations non-empty assertion |
| `README.md` | Modify | Add "Running the Demo" section after "Local Validation" |

No new files are created. No Go source files are modified.

---

## Data Models

Not applicable. This feature does not introduce or modify any Go structs, resource shapes,
or data structures. The demo script sends JSON payloads that conform to existing resource
models defined in FEATURE-0001 through FEATURE-0008.

### Demo Payloads Reference

The demo uses these exact resource payloads (matching existing API validation):

| Step | Kind | metadata.name | Key spec fields |
|------|------|--------------|-----------------|
| Organization | Organization | `nic` | `sovereignLocations: ["in-delhi-1","in-mumbai-1"]` |
| OrganizationUnit | OrganizationUnit | `ministry-health` | `organizationRef: "nic"` |
| Tenant | Tenant | `national-health-mission` | `organizationRef: "nic"`, `organizationUnitRef: "ministry-health"`, `isolationProfile: "namespace"` |
| Project | Project | `production` | `tenantRef: "national-health-mission"`, `organizationRef: "nic"`, `organizationUnitRef: "ministry-health"`, `environmentType: "production"` |
| ServiceClass | ServiceClass | `datastore.postgresql` | `category: "datastore"`, `requiredCapabilities: [...]` |
| ServicePlan | ServicePlan | `postgres-small-ha` | `serviceClassRef: "datastore.postgresql"` |
| Plugin | Plugin | `postgres.dstoreops.basic` | `serviceClassRefs: ["datastore.postgresql"]`, `pluginType: "dStoreOps"` |
| Capability (1) | Capability | `postgres-basic-provision` | `pluginRef: "postgres.dstoreops.basic"`, `operation: "Provision"` |
| Capability (2) | Capability | `postgres-basic-bind` | `pluginRef: "postgres.dstoreops.basic"`, `operation: "Bind"` |
| ServiceInstance | ServiceInstance | `nhm-prod-postgres` | `serviceClassRef`, `servicePlanRef`, `projectRef`, full hierarchy refs |
| ServiceBinding | ServiceBinding | `nhm-app-postgres-binding` | `serviceInstanceRef: "nhm-prod-postgres"` |

---

## Interfaces

Not applicable. No Go interfaces are introduced or modified. The demo script interacts
exclusively with the HTTP API — the interface is the REST contract defined in
`docs/api/API_CONTRACT_PHASE1.md`.

### Script Interface (shell functions)

The hardened script uses internal helper functions for clarity:

```bash
# step LABEL
#   Prints a human-readable step label to stdout.
step() { echo "==> $1"; }

# api_call METHOD URL EXPECTED_STATUS [DATA]
#   Executes curl with method, URL, and optional JSON body.
#   Captures response body and HTTP status code.
#   Validates HTTP status matches EXPECTED_STATUS.
#   On success: prints the response body.
#   On failure: prints step label, expected vs actual status, response body, exits 1.
#   Returns: sets global variable RESPONSE_BODY for further assertions.
api_call() {
  local method="$1" url="$2" expected_status="$3" data="${4:-}"
  # ... implementation ...
}

# assert_contains BODY EXPECTED_SUBSTRING LABEL
#   Verifies that BODY contains EXPECTED_SUBSTRING using bash substring matching.
#   Uses: [[ "$body" == *"$expected"* ]]
#   On failure: prints LABEL and exits 1.
assert_contains() {
  local body="$1" expected="$2" label="$3"
  if [[ "$body" != *"$expected"* ]]; then
    echo "FAIL: $label — expected substring '$expected' not found"
    exit 1
  fi
}

# assert_not_contains BODY UNEXPECTED_SUBSTRING LABEL
#   Verifies that BODY does NOT contain UNEXPECTED_SUBSTRING.
#   Uses: [[ "$body" != *"$unexpected"* ]]
#   On failure: prints LABEL and exits 1.
#   Used for: operations non-empty check (fail if '"items":[]' is present).
assert_not_contains() {
  local body="$1" unexpected="$2" label="$3"
  if [[ "$body" == *"$unexpected"* ]]; then
    echo "FAIL: $label — unexpected substring '$unexpected' found"
    exit 1
  fi
}
```

All calls use the single `api_call` function. Convenience wrappers are NOT used — every
call site passes `METHOD`, `URL`, `EXPECTED_STATUS`, and optionally `DATA` explicitly.
This removes ambiguity and keeps a single code path for response handling.

These functions are internal to the script and not exported.

---

## Validation

No resource validation logic is added or modified in Go code. The demo script relies on
existing validation implemented in FEATURE-0001 through FEATURE-0008.

### Script-level validation

The script performs these validations:

| Check | Method | Failure behavior |
|-------|--------|-----------------|
| Server reachable | `api_call GET /healthz 200` | Exit 1 with "Server not reachable" message |
| Server ready | `api_call GET /readyz 200` | Exit 1 with "Server not ready" message |
| Create returns 201 | `api_call POST url 201 data` | Exit 1 with step label + actual status |
| GET returns 200 | `api_call GET url 200` | Exit 1 with step label + actual status |
| List returns 200 | `api_call GET url 200` | Exit 1 with step label + actual status |
| Operations non-empty | `assert_not_contains` checks body does NOT contain `"items":[]` — fails if the exact substring `"items":[]` is present | Exit 1 with "No operations found" |
| GET ServiceInstance name | `assert_contains` checks body contains `nhm-prod-postgres` | Exit 1 with assertion failure |
| GET ServiceBinding name | `assert_contains` checks body contains `nhm-app-postgres-binding` | Exit 1 with assertion failure |

---

## API / Handler Design

No API handlers are added or modified. The demo exercises existing endpoints:

| Endpoint | Method | Expected Status | Purpose in Demo |
|----------|--------|-----------------|-----------------|
| `/healthz` | GET | 200 | Verify server is alive |
| `/readyz` | GET | 200 | Verify server is ready |
| `/v1/organizations` | POST | 201 | Create Organization |
| `/v1/organization-units` | POST | 201 | Create OrganizationUnit |
| `/v1/tenants` | POST | 201 | Create Tenant |
| `/v1/projects` | POST | 201 | Create Project |
| `/v1/service-classes` | POST | 201 | Register ServiceClass |
| `/v1/service-plans` | POST | 201 | Register ServicePlan |
| `/v1/plugins` | POST | 201 | Register Plugin |
| `/v1/capabilities` | POST | 201 | Register Capability (Provision) |
| `/v1/capabilities` | POST | 201 | Register Capability (Bind) |
| `/v1/service-instances` | POST | 201 | Create ServiceInstance |
| `/v1/service-bindings` | POST | 201 | Create ServiceBinding |
| `/v1/operations` | GET | 200 | List operations, verify non-empty |
| `/v1/service-instances/nhm-prod-postgres` | GET | 200 | Verify created instance |
| `/v1/service-bindings/nhm-app-postgres-binding` | GET | 200 | Verify created binding |

---

## Registry / Storage Design

Not applicable. The demo script does not modify registry or storage logic. It exercises
the existing in-memory registry through the HTTP API. All state is ephemeral — restarting
the server resets the registry.

---

## Operation / Audit Behavior

The demo **verifies** that operations are emitted but does not modify operation recording
logic. After all mutations:

1. Script calls `api_call GET /v1/operations 200`
2. Script asserts the response body does NOT contain the exact substring `"items":[]`.
   If `"items":[]` is present, the check fails with "No operations found" and exits 1.
   If `"items":[]` is absent (meaning the array has at least one element), the check passes.
3. Script prints the raw operations response body to stdout. No operation count is
   extracted or printed — the raw response is the only output.

**Synchronous recording assumption**: In Phase 1, operation records are written
synchronously by the API handler before it returns the HTTP response. Therefore, by the
time the demo receives a 201 response, the corresponding Operation record already exists
in the in-memory registry. No retry, sleep, or polling is needed for the operations
verification step.

Expected behavior (based on existing FEATURE-0005 implementation):
- Each POST that creates a resource emits an Operation record synchronously
- The operations list should contain entries for all 11 create/register calls
- Operation records include `resourceKind`, `resourceName`, `action`, and `status`

The demo does not validate individual operation field values — only that operations exist.

---

## Error Mapping

The demo script maps HTTP status codes to pass/fail decisions:

| HTTP Status | Script Interpretation | Action |
|-------------|----------------------|--------|
| 200 | Success (GET/List) | Continue |
| 201 | Success (POST create) | Continue |
| 400 | Validation failure | Print error, exit 1 |
| 404 | Resource not found | Print error, exit 1 |
| 409 | Conflict (duplicate on re-run) | Print error, exit 1 |
| 500 | Internal server error | Print error, exit 1 |
| Connection refused | Server not running | Print "Server not reachable", exit 1 |
| Timeout | Network timeout | curl fails, `set -e` halts, exit non-zero |

The script does not attempt retry or recovery. Any unexpected status causes immediate
termination per AC-2.

---

## Security and Privacy

1. **No real credentials**: All resource payloads use fictional example data (NIC,
   Ministry of Health, National Health Mission). No real government identifiers, API keys,
   tokens, or secrets are included.

2. **No external requests**: The script only communicates with `$BASE_URL` (default
   `http://127.0.0.1:8080`). No requests to external services or the internet.

3. **No file output**: All output goes to stdout/stderr. No response bodies are written
   to files that could leak into version control.

4. **No auth headers**: Phase 1 has no authentication layer. The script sends no
   `Authorization`, cookie, or token headers.

5. **SecretRef safety**: ServiceBinding responses may contain `secretRef` fields. The demo
   prints raw responses which in Phase 1 contain only placeholder data. In future phases
   with real secrets, the script would need modification to redact sensitive fields.

6. **No privilege escalation**: The script runs as the invoking user. It does not require
   root, sudo, or elevated permissions.

7. **Minimal attack surface**: The script invokes only two external programs: `bash`
   (the interpreter) and `curl`. No other binaries are executed. All string processing
   uses bash built-ins, eliminating shell injection vectors from piping to external tools.

---

## Testing Strategy

### Script Validation (manual and CI)

The demo script itself serves as a smoke test for the entire Phase 1 platform. Testing the
script involves running it against a live API server.

| Test Scenario | Method | Expected Result |
|---------------|--------|-----------------|
| Full demo against fresh server | `make run` then `make demo` | Script exits 0, all steps print success |
| Demo against stopped server | `make demo` without server | Script exits 1 at health check |
| Demo re-run without restart | `make demo` twice without restart | Script exits 1 with 409 on first create |
| Demo with wrong BASE_URL | `BASE_URL=http://127.0.0.1:9999 make demo` | Script exits 1 at health check |
| Demo against server mid-shutdown | Send SIGTERM to server, then `make demo` | Script exits 1 at readiness check |

### Portability Validation

| Platform | Validation |
|----------|-----------|
| macOS (bash 4+ via Homebrew) | Run `make demo` using Homebrew bash. Verify `bash --version` shows 4+ |
| Linux (bash 4+) | Run `make demo` on Linux (or CI) |
| bash 3.x (macOS default) | Not supported; script uses `$'\n'` parameter expansion requiring bash 4+. Document in README |

### Integration with CI

The demo script can be used as a CI smoke test:

```yaml
# Example CI step
- name: Build and smoke test
  run: |
    make build
    ./bin/sovrunn-api --config configs/sovrunn-api.local.yaml &
    sleep 2
    make demo
    kill %1
```

No dedicated Go test files are added for this feature. The demo script IS the integration
test.

---

## Verification Commands

```bash
# Verify script is syntactically valid
bash -n scripts/demo_phase1.sh

# Verify script is executable
ls -la scripts/demo_phase1.sh

# Run the full demo (server must be running)
make run &
sleep 2
make demo
kill %1

# Verify exit code on success
echo $?  # should be 0

# Verify exit code on failure (server not running)
make demo
echo $?  # should be non-zero
```

---

## Non-Goals (not implemented in this feature)

1. Real CLI binary (`sovrunn` command) — deferred to future phase
2. Interactive terminal UI or colored output requiring ncurses
3. Persistent state across demo runs
4. Load testing or performance benchmarking
5. Authentication or authorization headers
6. Demo against remote/production clusters
7. Docker-based demo orchestration
8. Automated server start/stop within the script
9. Idempotency via update-or-create (409 on re-run is acceptable)
10. DELETE endpoint testing in the demo flow
11. Multi-organization or multi-tenant demo scenarios
12. ServiceOps plugin execution or real provisioning
13. Pretty-printed JSON output (no `jq` dependency)
14. Response body schema validation (only status code + substring check)
15. Retry logic, polling, or exponential backoff — Phase 1 operations are recorded
    synchronously, so no timing-based verification is needed
16. Parallel curl calls for performance
17. Demo recording or replay capability
18. Counting operations — raw response is printed, no count extracted
19. Use of external text-processing tools (grep, sed, tail, awk) — all handled by bash built-ins

---

## Resolved Design Questions

### Q1: Verification depth

**Answer**: HTTP status code validation for all calls. Additionally, bash substring matching
(`[[ "$body" == *"pattern"* ]]`) for GET-by-name calls — asserting the exact resource name
string (`nhm-prod-postgres`, `nhm-app-postgres-binding`) appears in the response body.
For operations, assert the exact substring `"items":[]` is NOT present (meaning the list is
non-empty). No `jq` dependency. No grep/sed/tail. No full JSON schema validation. This
provides meaningful regression detection using only bash built-ins.

### Q2: Second Capability registration

**Answer**: Yes, both capabilities are registered. The feature file section 2 explicitly
lists steps 10 (Provision) and 11 (Bind). The demo registers:
- `postgres-basic-provision` with `operation: "Provision"`
- `postgres-basic-bind` with `operation: "Bind"`

### Q3: List filtering

**Answer**: No filtered list calls. The Phase 1 API contract does not document query
parameter filtering for list endpoints. The demo uses GET-by-name for verification and
unfiltered GET for the operations list.

### Q4: Output format

**Answer**: Raw JSON from curl. No `jq` dependency. Step labels provide human context.
Raw output is sufficient for Phase 1 evaluation. For the operations list, only the raw
response body is printed — no operation count is extracted or displayed.

### Q5: Server lifecycle

**Answer**: Manual. User starts server, runs demo, stops server. The script validates
availability at startup via `/healthz` and `/readyz`.

### Q6: Shebang portability

**Answer**: `#!/usr/bin/env bash`. Standard portable shebang. Works on macOS and Linux.

### Q7: Cleanup step

**Answer**: None. In-memory registry resets on server restart. No DELETE calls in the demo.

### Q8: README integration

**Answer**: New "## Running the Demo" section after "Local Validation" in README.md.

---

## Implementation Notes

### Script structure

The hardened `scripts/demo_phase1.sh` follows this structure:

```bash
#!/usr/bin/env bash
set -euo pipefail

# Configuration
BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"

# Helper functions
step()            # Print step label
api_call()        # Single HTTP helper: api_call METHOD URL EXPECTED_STATUS [DATA]
assert_contains() # Verify substring in response using bash [[ ]]
assert_not_contains() # Verify substring is absent (used for operations non-empty check)

# Demo flow
# 1. Health checks
# 2. Organization hierarchy (Org → OU → Tenant → Project)
# 3. Service catalog (ServiceClass → ServicePlan)
# 4. Plugin registry (Plugin → Capability × 2)
# 5. Service consumption (ServiceInstance → ServiceBinding)
# 6. Verification (Operations list, GET by name)
# 7. Success summary
```

### Helper function design

**`api_call METHOD URL EXPECTED_STATUS [DATA]`**: Single generic HTTP helper used for all
calls. Calls curl with `-s -w '\n%{http_code}'`. Splits the response into body and HTTP
status using bash parameter expansion (see Status Code Capture Technique below). Asserts
the HTTP status matches `EXPECTED_STATUS`. On failure: prints the step label, expected vs
actual status, and the response body, then exits 1. On success: prints the response body
and stores it in a global variable `RESPONSE_BODY` for subsequent assertions.

There are no separate `api_post` or `api_get` wrappers. Every call uses `api_call`
directly with all parameters explicit.

**`assert_contains BODY SUBSTRING LABEL`**: Uses bash `[[ "$body" == *"$substring"* ]]`.
Exits 1 with `LABEL` message if the substring is not found.

**`assert_not_contains BODY SUBSTRING LABEL`**: Uses bash `[[ "$body" != *"$substring"* ]]`.
Exits 1 with `LABEL` message if the substring IS found. Used for the operations non-empty
check: `assert_not_contains "$RESPONSE_BODY" '"items":[]' "No operations found"`.

### Status code capture technique

```bash
# Curl appends HTTP status code on a new line after the body
raw=$(curl -s -w '\n%{http_code}' -X "$method" "$url" \
  -H "Content-Type: application/json" \
  ${data:+-d "$data"})

# Split using bash parameter expansion — no sed/tail/grep needed
http_code="${raw##*$'\n'}"
body="${raw%$'\n'*}"
```

This uses only bash parameter expansion (`${raw##*$'\n'}` extracts the last line,
`${raw%$'\n'*}` removes the last line). No external tools (sed, tail, grep) are invoked.
The approach requires bash 4+ for reliable `$'\n'` handling in parameter expansion.

### Changes from existing script

| Current behavior | New behavior |
|-----------------|-------------|
| `#!/opt/homebrew/bin/bash` | `#!/usr/bin/env bash` |
| No status code validation | Every call validates HTTP status via `api_call` |
| No GET verification | GET ServiceInstance and ServiceBinding by name, assert exact resource name in body |
| One Capability registered | Two Capabilities registered (Provision + Bind) |
| No operation count check | Asserts operations list is non-empty (body does not contain `"items":[]`) |
| No summary output | Prints success summary at end |
| Uses `curl -fsS` (fail silently) | Uses explicit status capture via parameter expansion |
| No helper functions | `step`, `api_call`, `assert_contains`, `assert_not_contains` helpers |
| Implicit dependency on grep/sed | Pure bash + curl only; no coreutils dependency |

### README changes

Add after "Local Validation" section:

```markdown
## Running the Demo

Prerequisites:
- **bash 4+** required. macOS ships bash 3.2 by default. Install a modern bash via
  Homebrew (`brew install bash`) and either invoke with `/usr/local/bin/bash` (Intel)
  or `/opt/homebrew/bin/bash` (Apple Silicon), or add the Homebrew bash to your PATH.
  Verify with `bash --version`.
- **curl** must be available (pre-installed on macOS and most Linux distributions).
- The `sovrunn-api` server must be running (`make run`).

Run the full Phase 1 demo flow:

\```bash
make run &
sleep 2
make demo
\```

The demo exercises all Phase 1 resources end-to-end. To re-run, restart the server first:

\```bash
kill %1
make run &
sleep 2
make demo
\```

Override the base URL if the server runs on a different port:

\```bash
BASE_URL=http://127.0.0.1:9090 make demo
\```
```

---

## Go 1.21 Compatibility

Not applicable. No Go code is written or modified in this feature. The demo script uses
bash and curl only.

---

## Dependency Impact

No new dependencies. The script uses **bash + curl only**. No coreutils tools (grep, sed,
tail, awk) are invoked at runtime. All string operations use bash built-ins:

- Response splitting: bash parameter expansion (`${var##pattern}`, `${var%pattern}`)
- Substring checks: bash conditional `[[ "$var" == *"pattern"* ]]`
- Output: `echo` (bash built-in)

Runtime requirements:
- **bash 4+** — required for reliable `$'\n'` handling in parameter expansion. macOS
  ships bash 3.2; users must install bash 4+ via Homebrew (see README changes above).
- **curl** — available by default on macOS and Linux.
- A running `sovrunn-api` server (from `make run`).

No Python, jq, grep, sed, tail, Node.js, Docker, or other runtime dependencies.

