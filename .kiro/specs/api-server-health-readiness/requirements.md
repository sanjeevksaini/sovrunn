# Requirements Document

## Introduction

FEATURE-0009 finalizes the API server health, readiness, and version endpoints for Sovrunn
Phase 1. These endpoints enable infrastructure orchestrators (Kubernetes liveness/readiness
probes, load balancers, monitoring systems, and CI/CD pipelines) to determine whether the
sovrunn-api process is alive, ready to serve traffic, and which build is running.

A minimal bootstrap implementation already exists in `internal/api/bootstrap_handler.go` and
`internal/health/readiness.go`. This feature formalizes the contract, hardens the behavior,
adds missing capabilities (HTTP method enforcement, structured readiness checks, startup
metadata, graceful degradation), ensures comprehensive test coverage, and documents the
acceptance criteria that external integrators can rely on.

This feature depends on FEATURE-0001 (Organization Resource and Registry) because the
readiness check validates that the in-memory registry subsystem is initialized. It is
compatible with all previously implemented Phase 1 features (FEATURE-0001 through
FEATURE-0008) and does not alter any resource CRUD behavior.

## Glossary

| Term | Definition |
|---|---|
| Liveness Probe | Infrastructure check that determines whether the process is alive and should not be restarted. Mapped to `/healthz`. |
| Readiness Probe | Infrastructure check that determines whether the process can serve traffic. Mapped to `/readyz`. |
| Version Endpoint | Metadata endpoint that returns build and runtime identification. Mapped to `/version`. |
| ReadinessState | Atomic boolean flag tracking whether server initialization completed successfully. |
| Readiness Check | A named sub-check that contributes to the aggregate readiness decision. |
| Startup Metadata | Static build-time and runtime information returned by `/version`. |
| bootstrapChain | Middleware chain applied to health/readiness/version endpoints consisting of request-ID assignment followed by structured logging. Defined in `internal/server/server.go`. |

## Decisions (Resolved)

The following decisions were made during requirements review and are authoritative:

1. **Version injection target**: The authoritative ldflags symbol is
   `github.com/sanjeevksaini/sovrunn/internal/api.buildVersion`. This corresponds to the
   package-level variable `var buildVersion` in `internal/api/bootstrap_handler.go`.
2. **Error envelope**: All error responses from health/readiness/version endpoints SHALL use
   `resources.APIErrorEnvelope` (defined in `internal/resources/errors.go`). This is the same
   envelope used by all resource API errors. **Exception**: the `/readyz` 503 not-ready
   response is NOT an error in the API-contract sense — it is an expected operational state.
   The 503 from `/readyz` MUST use the domain shape
   `{"status":"not_ready","message":"initializing|shutting_down"}` and MUST NOT use
   `resources.APIErrorEnvelope`. Only true error conditions (e.g. 405 Method Not Allowed)
   use the error envelope on these endpoints.
3. **Middleware ordering**: Requests pass through bootstrapChain (request-ID → logging) first,
   then method enforcement runs, then endpoint-specific handler logic.
4. **Not-ready message content**: The `message` field in 503 responses SHALL use a fixed set
   of generic reason strings: `"initializing"` or `"shutting_down"`. Named sub-check
   identifiers are NOT exposed. This avoids information leakage while remaining useful for
   operators.
5. **Body logging**: bootstrapChain does NOT log request or response bodies. The logging
   middleware logs only `request_id`, `method`, `path`, `status_code`, and `latency_ms`.
   No special-case code is needed to suppress body logging for these endpoints.

## Requirements

---

### Requirement 1: Liveness Endpoint (`/healthz`)

**User Story:** As an infrastructure operator, I want a lightweight liveness endpoint so that
orchestrators can detect whether the sovrunn-api process is alive and responsive.

#### Acceptance Criteria

1. `GET /healthz` SHALL return HTTP 200 with JSON body `{"status":"ok"}`.
2. The endpoint SHALL NOT perform any registry, database, or external-system checks.
3. The endpoint SHALL respond within 10ms under normal conditions (no blocking I/O).
4. The response `Content-Type` header SHALL be `application/json`.
5. Non-GET methods (POST, PUT, DELETE, PATCH, HEAD) SHALL return HTTP 405 Method Not Allowed.
6. The endpoint SHALL be available immediately after the HTTP listener starts, even before
   readiness is set to true.
7. The endpoint SHALL pass through bootstrapChain (request-ID → logging) and then method
   enforcement before handler logic executes.
8. The endpoint SHALL NOT require authentication in Phase 1.

---

### Requirement 2: Readiness Endpoint (`/readyz`)

**User Story:** As an infrastructure operator, I want a readiness endpoint so that load
balancers and orchestrators only route traffic to instances that have completed initialization.

#### Acceptance Criteria

1. `GET /readyz` SHALL return HTTP 200 with JSON body `{"status":"ready"}` when all readiness
   conditions are satisfied.
2. `GET /readyz` SHALL return HTTP 503 Service Unavailable with JSON body
   `{"status":"not_ready","message":"<reason>"}` when any readiness condition is not met.
   The `message` field SHALL be one of the following fixed generic strings:
   - `"initializing"` — server has not yet completed startup.
   - `"shutting_down"` — graceful shutdown has begun.
   Named sub-check identifiers (e.g. registry name, config subsystem) SHALL NOT appear in
   the message. This prevents information leakage while remaining operationally useful.
3. Readiness conditions for Phase 1 SHALL be:
   - In-memory registry subsystem is initialized.
   - Configuration has been loaded without fatal error.
   Note: "HTTP server listener is bound" is intentionally excluded as a readiness condition
   because it is inherently satisfied by the time any request reaches `/readyz`. The server
   cannot serve the readiness probe unless the listener is already bound, making this check
   circular and redundant.
4. The ReadinessState SHALL transition to `true` only after all Phase 1 readiness conditions
   pass during server startup.
5. The ReadinessState SHALL transition to `false` when graceful shutdown begins.
6. Non-GET methods (POST, PUT, DELETE, PATCH, HEAD) SHALL return HTTP 405 Method Not Allowed.
7. The endpoint SHALL pass through bootstrapChain (request-ID → logging) and then method
   enforcement before handler logic executes.
8. The endpoint SHALL NOT require authentication in Phase 1.
9. The response `Content-Type` header SHALL be `application/json`.
10. The endpoint SHALL NOT call slow external systems or perform expensive computations.

---

### Requirement 3: Version Endpoint (`/version`)

**User Story:** As a platform operator or developer, I want a version endpoint so that I can
identify which build is running and confirm the deployment phase.

#### Acceptance Criteria

1. `GET /version` SHALL return HTTP 200 with a JSON object containing at minimum:
   - `name` (string): `"sovrunn-api"`
   - `version` (string): semantic version or `"dev"` for local builds
   - `phase` (string): `"1"`
   - `status` (string): `"alpha"`
2. The `version` field SHALL be injectable at build time via `-ldflags` targeting the
   package-level variable `internal/api.buildVersion`. The authoritative ldflags argument is:
   `-X 'github.com/sanjeevksaini/sovrunn/internal/api.buildVersion=$(VERSION)'`.
3. Non-GET methods (POST, PUT, DELETE, PATCH, HEAD) SHALL return HTTP 405 Method Not Allowed.
4. The endpoint SHALL pass through bootstrapChain (request-ID → logging) and then method
   enforcement before handler logic executes.
5. The endpoint SHALL NOT require authentication in Phase 1.
6. The response `Content-Type` header SHALL be `application/json`.
7. The response SHALL NOT include secrets, tokens, internal file paths, or host-specific
   information that could aid reconnaissance.

---

### Requirement 4: HTTP Method Enforcement

**User Story:** As a security-conscious operator, I want health/readiness/version endpoints to
reject non-GET requests so that the surface area is minimal and predictable.

#### Acceptance Criteria

1. `/healthz`, `/readyz`, and `/version` SHALL accept only the `GET` HTTP method.
2. All other methods — including POST, PUT, DELETE, PATCH, and HEAD — SHALL receive HTTP 405
   Method Not Allowed.
3. For non-HEAD methods receiving 405, the response SHALL include a JSON error body using
   `resources.APIErrorEnvelope`:
   `{"error":{"code":"METHOD_NOT_ALLOWED","message":"only GET is supported"}}`.
4. For HEAD requests receiving 405, the response MUST NOT include a body per HTTP semantics.
   The response SHALL include status code 405, header `Allow: GET`, and header
   `Content-Type: application/json`. No response body SHALL be written.
5. The 405 response SHALL include an `Allow: GET` header (applies to all methods including
   HEAD).
6. Method enforcement SHALL execute after bootstrapChain (request-ID + logging) and before
   any endpoint-specific handler logic. The middleware ordering is:
   request-ID → logging → method enforcement → endpoint handler.
7. The 405 response `Content-Type` header SHALL be `application/json`.

---

### Requirement 5: Graceful Shutdown Interaction

**User Story:** As an infrastructure operator, I want the readiness endpoint to immediately
report not-ready when shutdown begins so that orchestrators stop sending new requests.

#### Acceptance Criteria

1. When the server receives SIGINT or SIGTERM, ReadinessState SHALL be set to `false` before
   the HTTP server stops accepting new connections.
2. In-flight requests to `/readyz` after ReadinessState becomes `false` SHALL return 503
   with message `"shutting_down"`.
3. In-flight requests to `/healthz` SHALL continue returning 200 until the listener closes.
4. The shutdown sequence SHALL allow a configurable drain timeout (from `config.Server.ShutdownTimeout`).
5. After the drain timeout, remaining in-flight connections SHALL be forcibly closed.

---

### Requirement 6: Structured Observability for Health Endpoints

**User Story:** As a platform operator, I want health and readiness endpoint access to be
logged with standard structured fields so that I can monitor probe frequency and failures.

#### Acceptance Criteria

1. Requests to `/healthz`, `/readyz`, and `/version` SHALL be logged with the same structured
   fields as other requests: `request_id`, `method`, `path`, `status_code`, `latency_ms`.
2. Health endpoint logs SHALL NOT include request or response bodies. The bootstrapChain
   logging middleware does not log bodies by design; no special-case suppression is needed.
3. Logging SHALL NOT add measurable latency to health endpoint responses.
4. The request-ID header (`X-Request-ID`) SHALL be propagated or generated for these endpoints.

---

### Requirement 7: No Authentication Required

**User Story:** As an infrastructure operator, I want health endpoints accessible without
credentials so that Kubernetes probes and load balancer checks work without token management.

#### Acceptance Criteria

1. `/healthz`, `/readyz`, and `/version` SHALL NOT require any authentication token, API key,
   or client certificate in Phase 1.
2. When authentication is added in future phases, these endpoints SHALL remain unauthenticated
   unless an explicit security decision changes this.
3. The endpoint handlers SHALL NOT read or validate `Authorization` headers.

---

### Requirement 8: Response Consistency

**User Story:** As a developer integrating with Sovrunn, I want health endpoints to return
consistent JSON structures so that automated tooling can parse responses reliably.

#### Acceptance Criteria

1. All successful health endpoint responses SHALL be valid JSON objects.
2. All error responses (405) from these endpoints SHALL use `resources.APIErrorEnvelope`
   (defined in `internal/resources/errors.go`) with shape:
   `{"error":{"code":"...","message":"..."}}`.
   This is the same envelope used by all resource API errors across the platform.
   **Exception — HEAD method**: HEAD 405 responses MUST NOT include a body per HTTP
   semantics. They SHALL include status code 405, header `Allow: GET`, and header
   `Content-Type: application/json`, but no response body. The APIErrorEnvelope requirement
   applies only to non-HEAD 405 responses.
3. Exception: the 503 not-ready response from `/readyz` uses the domain shape
   `{"status":"not_ready","message":"..."}` because it is not an error in the API-contract
   sense — it is an expected operational state. Only true error conditions (405) use the
   error envelope.
4. Responses SHALL NOT return plain text, HTML, or empty bodies — except for HEAD responses,
   which MUST NOT include a body per HTTP semantics.
5. `Content-Type: application/json` SHALL be set on all responses from these endpoints
   (including HEAD responses where no body is written).
6. Responses SHALL NOT include trailing newlines outside the JSON body (consistent with
   existing `writeJSON` behavior).

---

### Requirement 9: Build-Time Version Injection

**User Story:** As a release engineer, I want the version endpoint to report the exact build
version without requiring runtime configuration so that deployed binaries are self-identifying.

#### Acceptance Criteria

1. The Makefile SHALL define a `VERSION` variable (defaulting to `dev`).
2. The `go build` command in the Makefile SHALL inject the version via `-ldflags`:
   `-X 'github.com/sanjeevksaini/sovrunn/internal/api.buildVersion=$(VERSION)'`.
   This targets the package-level variable `var buildVersion` in
   `internal/api/bootstrap_handler.go`.
3. When built without ldflags (e.g. `go run`), the version SHALL default to `"dev"`.
4. The injected version SHALL appear in the `/version` response `version` field.
5. All examples and references in this document use the authoritative ldflags path
   `github.com/sanjeevksaini/sovrunn/internal/api.buildVersion`. No other injection target
   (e.g. `main.buildVersion`) is valid.

---

### Requirement 10: Test Coverage

**User Story:** As a developer, I want comprehensive tests for all health endpoints so that
regressions are caught immediately.

#### Acceptance Criteria

1. Unit tests SHALL cover:
   - `/healthz` returns 200 with `{"status":"ok"}`.
   - `/readyz` returns 200 when ReadinessState is true.
   - `/readyz` returns 503 with `{"status":"not_ready","message":"initializing"}` when
     ReadinessState is false (pre-startup).
   - `/readyz` returns 503 with `{"status":"not_ready","message":"shutting_down"}` when
     ReadinessState is false (post-shutdown signal).
   - `/version` returns 200 with all required fields.
   - Non-GET methods (POST, PUT, DELETE, PATCH) return 405 for all three endpoints with
     `resources.APIErrorEnvelope` JSON body.
   - HEAD requests return 405 for all three endpoints with correct status code and headers
     (`Allow: GET`, `Content-Type: application/json`) but NO response body. Tests for HEAD
     MUST assert status and headers only and MUST NOT assert presence of a JSON body.
   - 405 responses include `Allow: GET` header (all methods including HEAD).
   - 405 responses for non-HEAD methods use `resources.APIErrorEnvelope` shape.
   - Response Content-Type is `application/json`.
   - `/readyz` 503 responses use domain shape `{"status":"not_ready","message":"..."}`
     and NOT `resources.APIErrorEnvelope`.
2. Integration-level tests SHALL verify:
   - The test starts the server and polls `/readyz` with a bounded timeout (max 5 seconds,
     polling every 50ms) until it returns 200. Only after receiving 200 does it assert
     `/healthz` also returns 200. This avoids startup race conditions.
   - The test triggers graceful shutdown (e.g. sends signal or calls shutdown method), then
     asserts `/readyz` returns 503 with message `"shutting_down"` during the drain window.
   - Tests use deterministic synchronization (poll-until-ready with timeout) rather than
     fixed-duration sleeps.
3. Tests SHALL run deterministically without network or external dependencies.
4. Tests SHALL pass under `go test -race ./...` without data races.

---

## Non-Goals

The following are explicitly out of scope for FEATURE-0009:

1. **Deep health checks**: No database connectivity, external service reachability, or
   dependency-graph health checks. Phase 1 uses in-memory registry only.
2. **Metrics endpoint** (`/metrics`): Prometheus metrics exposure is a future capability.
3. **Authentication on health endpoints**: No token/key/cert requirement.
4. **Custom readiness check registration**: No plugin-based or dynamic readiness sub-checks.
5. **Health endpoint rate limiting**: No request throttling on probe endpoints.
6. **HTML or human-readable health page**: JSON-only responses.
7. **Distributed readiness coordination**: No cross-instance readiness aggregation.
8. **OpenTelemetry health instrumentation**: Standard structured logging only in Phase 1.
9. **Startup probe endpoint** (`/startupz`): Kubernetes startup probes can use `/healthz`.
10. **TLS termination**: Health endpoints are served over plain HTTP in Phase 1.
11. **Verbose readiness mode**: `/readyz?verbose=true` with per-check breakdown is deferred
    to a future feature.

## Edge Cases

1. **Server starts but registry initialization panics**: ReadinessState remains `false`. `/healthz`
   returns 200 (process alive), `/readyz` returns 503 with message `"initializing"`.
   Orchestrator may restart.
2. **Concurrent shutdown during readiness check**: A `/readyz` request in-flight when shutdown
   signal arrives may see either 200 or 503 depending on timing. This is acceptable because
   the atomic boolean provides memory-safe access and the next probe will see 503.
3. **Rapid repeated probes**: Health endpoints must remain O(1) and allocation-minimal to
   handle high-frequency probing (e.g. every 1–5 seconds per orchestrator).
4. **Request with body on GET**: Health endpoints SHALL ignore any request body on GET requests
   and respond normally. They SHALL NOT return 400 for unexpected body content on GET.
5. **Unknown query parameters**: Health endpoints SHALL ignore unknown query parameters and
   respond normally.
6. **HEAD requests**: Health endpoints SHALL return 405 for HEAD (only GET is allowed), keeping
   the contract simple and predictable. The `Allow: GET` header is included in the response.
   Per HTTP semantics, the 405 response to HEAD MUST NOT include a body; only status code
   and headers (`Allow: GET`, `Content-Type: application/json`) are returned.
7. **Very large request headers**: The HTTP server's existing `MaxHeaderBytes` or Go default
   (1MB) applies. Health endpoints do not add special header-size handling.
8. **Version endpoint when ldflags not set**: SHALL return `"dev"` as the version string.

## Security and Privacy Requirements

1. Health endpoints SHALL NOT expose internal IP addresses, hostnames, file paths, goroutine
   counts, memory statistics, or other system internals.
2. Health endpoints SHALL NOT log or expose any data from request bodies or authorization
   headers.
3. The `/version` response SHALL NOT include Git commit hashes, build machine names, or CI
   pipeline identifiers unless explicitly approved in a future decision.
4. Health endpoints SHALL be safe to expose on a public network without information leakage
   concerns.
5. Health endpoints SHALL NOT be usable as an amplification vector (responses are small,
   fixed-size JSON).
6. The `/readyz` 503 message field SHALL contain only generic state labels (`"initializing"`,
   `"shutting_down"`) and SHALL NOT contain subsystem names, error details, stack traces, or
   internal identifiers.

## Compatibility with Completed Phase 1 Features

| Feature | Compatibility |
|---|---|
| FEATURE-0001 Organization Resource and Registry | No changes to Organization CRUD. Registry initialization is a readiness precondition. |
| FEATURE-0002 OrganizationUnit Resource | No changes. Health endpoints are independent of OU lifecycle. |
| FEATURE-0003 Tenant Resource | No changes. Health endpoints are independent of Tenant lifecycle. |
| FEATURE-0004 Project Resource | No changes. Health endpoints are independent of Project lifecycle. |
| FEATURE-0005 Operation Resource | No changes. Health/readiness/version requests do NOT emit Operation records (they are not mutating governance actions). |
| FEATURE-0006 ServiceClass and ServicePlan | No changes. Health endpoints are independent of catalog resources. |
| FEATURE-0007 Plugin and Capability Registry | No changes. Health endpoints are independent of plugin lifecycle. |
| FEATURE-0008 ServiceInstance and ServiceBinding | No changes. Health endpoints are independent of service consumption resources. |

The existing `BootstrapHandler`, `ReadinessState`, and server route registration remain the
authoritative implementation locations. FEATURE-0009 hardens and completes them; it does not
move or restructure them.

## Design Questions to Resolve in design.md

1. **Method enforcement placement**: Should 405 enforcement be a dedicated middleware wrapping
   each bootstrap endpoint (positioned after bootstrapChain), or inline logic at the top of
   each handler method? The requirement mandates ordering: bootstrapChain → method enforcement
   → handler. Design must honor this.
2. **Verbose readiness mode**: Deferred to a future feature. Not in scope for FEATURE-0009.
3. **Shutdown race window**: Is the current atomic-boolean approach sufficient, or should the
   server use `http.Server.RegisterOnShutdown` callbacks for more deterministic ordering?
4. **Version field extensibility**: Should `/version` use a fixed struct or allow optional
   fields (e.g. `goVersion`, `buildTime`) to be added without breaking clients?
5. **Makefile ldflags target**: Should the Makefile modify the existing `build` target to
   include ldflags, or add a dedicated `release` target? The requirement mandates that
   ldflags are in the `go build` command used by the Makefile.
