---
doc_type: engineering_standard
title: Sovrunn Go Coding Guardrails
status: draft
phase: 1
ai_load_priority: always
ai_summary: Mandatory Go coding guardrails for AI agents implementing Sovrunn with focus on correctness, security, latency, performance, observability, horizontal scaling, and serverless readiness.
---

# Sovrunn Go Coding Guardrails

## 1. Purpose

This document defines Go-specific coding rules and guardrails for Sovrunn.

AI coding agents such as Kiro, Cursor, ChatGPT, Claude Code, and other assistants must follow these rules when generating, editing, refactoring, or reviewing Go code.

The priorities are:

```text
correctness
security
latency
performance
observability
horizontal scalability
serverless readiness
maintainability
testability
```

Correctness and security come before optimization. Performance work must be measured, not guessed.

## 2. Scope

These guardrails apply to Go code in Sovrunn, including:

```text
API server
resource model
in-memory registry
validation
operation framework
audit hooks
health/readiness handlers
future control plane services
future management plane services
future ServiceOps framework
future SDE control-plane integration
```

These guardrails do not authorize implementation of future capabilities. They only define how Go code should be written when a feature is approved.

## 3. Core Go Principles

Use idiomatic Go.

Prefer:

```text
simple code
explicit structs
small interfaces
context-aware functions
deterministic validation
structured errors
clear package boundaries
minimal dependencies
good tests
boring implementation
```

Avoid:

```text
clever abstractions
large frameworks
reflection-heavy code
global mutable state
hidden goroutines
unbounded queues
unbounded memory growth
panic-based control flow
unmeasured optimization
```

Sovrunn Go code should be easy for humans and AI agents to reason about.

## 4. Package Design

Use internal packages for platform implementation.

Recommended Phase 1 structure:

```text
cmd/sovrunn-api/
internal/api/
internal/audit/
internal/config/
internal/health/
internal/operation/
internal/registry/
internal/resources/
internal/server/
internal/validation/
tests/integration/
```

Rules:

```text
cmd/ contains application entrypoints only.
internal/ contains implementation packages.
pkg/ is avoided unless a stable external SDK is intentionally exposed.
resources/ defines resource structs.
registry/ stores resource state.
validation/ validates resource specs and references.
api/ exposes HTTP handlers.
operation/ records lifecycle operations.
audit/ records accountability events.
config/ loads and validates configuration.
health/ exposes health and readiness checks.
server/ owns HTTP server lifecycle.
```

Do not create new packages unless they have a clear responsibility.

Package names should be short, lowercase, and meaningful.

Good:

```text
resources
registry
validation
operation
server
```

Avoid:

```text
utils
common
helpers
misc
manager
processor
```

Generic package names hide responsibilities and confuse AI agents.

## 5. Feature Boundary Rules

Implement one approved feature at a time.

AI coding agents must not implement future features early.

For example, during `FEATURE-0001 Organization Resource and Registry`, do not implement:

```text
OrganizationUnit
Tenant
Project
ServiceClass
ServicePlan
Plugin
Capability
ServiceInstance
ServiceBinding
persistent storage
Kubernetes CRDs
ServiceOps execution
UI
AI agent execution
```

Small shared scaffolding is allowed only when required by the current feature, such as:

```text
common metadata structs
common error type
basic server setup
basic health/readiness endpoint
basic validation helpers
```

Shared scaffolding must not become hidden implementation of future features.

## 6. Context Rules

All request-scoped and potentially blocking functions must accept `context.Context`.

Good:

```go
func (r *Registry) CreateOrganization(ctx context.Context, org resources.Organization) error
```

Bad:

```go
func (r *Registry) CreateOrganization(org resources.Organization) error
```

Rules:

```text
context.Context must be the first parameter.
Do not store context.Context in structs.
Do not pass nil context.
Use request context in HTTP handlers.
Respect ctx.Done() in long-running operations.
Do not create background contexts inside request-handling code except at process boundaries.
```

Use:

```go
ctx := r.Context()
```

inside HTTP handlers.

Do not use context values as a general dependency injection mechanism. Use them only for request-scoped values such as request ID, actor identity, trace context, or deadlines.

## 7. Error Handling

Do not use panic for normal errors.

Use explicit, stable error codes.

Recommended error shape:

```go
type ErrorCode string

const (
    ErrValidationFailed       ErrorCode = "VALIDATION_FAILED"
    ErrResourceNotFound       ErrorCode = "RESOURCE_NOT_FOUND"
    ErrResourceAlreadyExists  ErrorCode = "RESOURCE_ALREADY_EXISTS"
    ErrDeleteBlocked          ErrorCode = "DELETE_BLOCKED"
    ErrReferenceInvalid       ErrorCode = "REFERENCE_INVALID"
    ErrPolicyDenied           ErrorCode = "POLICY_DENIED"
    ErrUnauthorized           ErrorCode = "UNAUTHORIZED"
    ErrForbidden              ErrorCode = "FORBIDDEN"
    ErrConflict               ErrorCode = "CONFLICT"
    ErrInternal               ErrorCode = "INTERNAL_ERROR"
)
```

Recommended API error structure:

```go
type APIError struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
    Field   string    `json:"field,omitempty"`
    Details string    `json:"details,omitempty"`
}
```

Rules:

```text
Errors must be actionable.
Errors must not leak secrets.
Errors must map cleanly to HTTP status codes.
Validation errors should identify the invalid field.
Reference errors should identify the missing referenced resource.
Internal errors must not expose stack traces or internal implementation details.
```

Use error wrapping for internal diagnostics:

```go
fmt.Errorf("create organization: %w", err)
```

But do not return wrapped internal details directly to API clients.

## 8. HTTP API Rules

Use the standard library `net/http` unless a small router is explicitly approved.

HTTP handlers must:

```text
read request context
generate or propagate request ID
limit request body size
decode JSON safely
validate input
call registry/service layer
return structured JSON errors
set Content-Type consistently
set request ID response header
set operation ID response header after FEATURE-0005
never expose stack traces
```

Recommended request body limit for Phase 1:

```go
http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB
```

HTTP status guidance:

```text
201 Created     for create
200 OK          for get, list, update
202 Accepted    for future long-running operations
204 No Content  for delete when no body is returned
400 Bad Request for malformed JSON or validation failure
401 Unauthorized for missing/invalid authentication, future
403 Forbidden   for policy/RBAC denial, future
404 Not Found   for missing resource
409 Conflict    for duplicate, delete-blocked, or reference conflict
413 Payload Too Large for oversized request body
415 Unsupported Media Type for invalid content type
500 Internal Server Error for unexpected server failures
```

Do not block HTTP requests for long-running infrastructure operations in future phases. Return an `Operation` instead.

## 9. JSON Rules

Use explicit JSON tags.

Good:

```go
type Metadata struct {
    Name        string            `json:"name"`
    DisplayName string           `json:"displayName,omitempty"`
    Labels      map[string]string `json:"labels,omitempty"`
    Annotations map[string]string `json:"annotations,omitempty"`
}
```

Rules:

```text
Do not expose internal-only fields in JSON.
Do not accept user-authored status in create/update requests.
Use omitempty only where absence is meaningful.
Keep API stable once accepted.
Use explicit request/response structs when API behavior differs from internal structs.
```

Avoid using:

```go
map[string]any
```

except at clear boundaries where schema-free data is intentional.

Do not use JSON marshal/unmarshal cycles to copy structs in hot paths.

## 10. Resource Shape Rules

All Sovrunn resources must follow the canonical shape:

```text
apiVersion
kind
metadata
spec
status
```

Rules:

```text
metadata = identity and classification
spec = desired state
status = observed state
operation = lifecycle trace
```

User input may set:

```text
apiVersion
kind
metadata.name
metadata.displayName
metadata.labels
metadata.annotations
spec
```

User input must not set:

```text
status
resourceVersion
generation
createdAt
updatedAt
operationRef
```

If a client submits system-owned fields, handlers must reject them or ignore them according to the API contract. Prefer rejection for clarity in Phase 1.

## 11. Naming and Identifier Rules

Resource names must be deterministic and safe.

Recommended naming rule:

```text
lowercase DNS-compatible name
letters, numbers, and hyphen
must start with a letter or number
must end with a letter or number
reasonable maximum length
```

Example validation pattern:

```text
^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
```

Rules:

```text
Do not allow spaces in resource names.
Do not allow path separators in names.
Do not allow hidden normalization that changes user intent.
Display names may be human-readable.
Labels and annotations must be bounded.
```

Reject ambiguous or unsafe names early.

## 12. Validation Rules

Validation must be deterministic and explicit.

Validate:

```text
required fields
DNS-compatible names
duplicate names
reference existence
parent-child consistency
allowed enum values
delete-blocked conditions
status not user-authored
label and annotation size limits
unsupported kind/apiVersion
```

Do not rely on database errors, map behavior, or downstream execution failures for validation.

Validation functions should be testable without starting the API server.

Good:

```go
func ValidateOrganization(org resources.Organization) []FieldError
```

Avoid validation that mutates input unexpectedly.

## 13. Registry Rules

Phase 1 uses in-memory registry.

Registry must be:

```text
thread-safe
deterministic
testable
context-aware
free of hidden external dependencies
replaceable in future phases
```

Use mutexes explicitly.

Example:

```go
type Registry struct {
    mu            sync.RWMutex
    organizations map[string]resources.Organization
}
```

Rules:

```text
Use sync.RWMutex for map-backed registries.
Protect all reads and writes.
Do not return mutable internal maps.
Return copies or immutable views.
Avoid package-level global registries.
Do not start goroutines inside registry methods.
Do not use registry as an event bus.
Do not use local in-memory locks as future distributed locks.
```

Pre-size lists when possible:

```go
items := make([]resources.Organization, 0, len(r.organizations))
```

When returning resources, ensure callers cannot mutate internal state accidentally.

## 14. Operation Rules

Operations represent lifecycle activity.

After `FEATURE-0005`, mutating requests should produce Operation records.

Operation records should support future asynchronous execution, auditability, and troubleshooting.

Rules:

```text
Do not hide lifecycle changes.
Do not mutate important resources without operation trace after operation framework exists.
Operation IDs must be unique.
Operation status must be explicit.
Operation errors must use stable error codes.
```

Future long-running actions must return quickly with an Operation reference instead of blocking HTTP requests.

## 15. Performance Rules

Prioritize simple low-allocation code.

Rules:

```text
Avoid reflection in hot paths.
Avoid unnecessary JSON marshal/unmarshal cycles.
Avoid converting structs to map[string]any unless needed.
Avoid regex in hot validation paths unless precompiled.
Avoid per-request allocation-heavy logging fields.
Avoid unnecessary goroutines per request.
Avoid unbounded buffering.
Avoid defer in very hot loops if measurable overhead matters.
Pre-size maps and slices when size is known.
Use strings.Builder for repeated string construction.
Use bytes.Buffer or strings.Builder instead of repeated string concatenation in loops.
```

Good:

```go
items := make([]Organization, 0, len(r.organizations))
```

Bad:

```go
items := []Organization{}
```

Do not optimize blindly. Add benchmarks for performance-sensitive paths.

## 16. Latency Rules

API handlers should do the minimum required work synchronously.

Rules:

```text
Keep request path short.
Validate before mutation.
Avoid blocking I/O in handlers.
Avoid background goroutine fan-out in request path.
Avoid synchronous remote calls in Phase 1.
Avoid logging massive payloads.
Avoid expensive reflection in hot paths.
Return quickly after registry update and operation recording.
```

Future infrastructure operations should be represented as `Operation`, not held open inside HTTP requests.

## 17. Security Rules

Never log secrets.

Never return internal stack traces.

Never trust user input.

Rules:

```text
limit request body size
validate content type for JSON endpoints
reject unknown or unsupported resource kinds
reject user-authored status
sanitize error messages
do not echo credentials
do not log Authorization headers
do not log cookies
do not log tokens
do not log passwords
do not log private keys
do not log full request bodies by default
```

Use constant-time comparisons for secrets or tokens when implemented later.

Do not implement authentication or authorization ad hoc in Phase 1. Prepare clear boundaries for future identity integration.

Security-sensitive code must be explicit and reviewed by a human.

## 18. Input Handling Rules

HTTP handlers must handle bad input safely.

Rules:

```text
reject malformed JSON
reject unknown resource kinds
reject unsupported apiVersion
reject oversized request bodies
reject invalid content type
reject path/body name mismatch where applicable
reject user-authored system fields
do not silently coerce invalid input
```

For update requests, if path name and body metadata.name are both present, they must match.

## 19. Observability Rules

Every request should have structured observability fields.

Required request log fields:

```text
request_id
method
path
status_code
latency_ms
error_code, on failure
```

When available, include:

```text
actor
resource_kind
resource_name
organization
organizationUnit
tenant
project
operation_id
```

Do not log full resource payloads by default.

For future OpenTelemetry readiness, keep request handling structured as:

```text
decode
validate
authorize, future
execute
record operation
record audit
respond
```

Avoid observability code that strongly couples business logic to one logging or tracing provider.

## 20. Monitoring and Metrics Rules

Metrics should be easy to add without large refactors.

Future metrics should include:

```text
http_requests_total
http_request_duration_seconds
api_errors_total
registry_resources_total
operations_total
operations_failed_total
validation_failures_total
```

Rules:

```text
Keep metric labels bounded.
Do not use resource names as high-cardinality metric labels unless explicitly approved.
Do not expose secrets in metrics.
Do not block request path on metrics export.
```

Phase 1 may start with logs only, but code should not prevent future metrics.

## 21. Audit Rules

Mutating requests must be auditable.

After `FEATURE-0005`, mutating requests must create Operation records.

Future AuditEvent should capture:

```text
actor
action
resource_kind
resource_name
organization
organizationUnit
tenant
project
request_id
operation_id
status
error_code
message
timestamp
```

Rules:

```text
Audit data must be structured.
Audit events must not contain secrets.
Audit logging must not leak full credentials or tokens.
Audit hooks must not block indefinitely.
```

Phase 1 may log audit-style records even before durable audit storage exists.

## 22. Horizontal Scaling Rules

Design code so API instances can scale horizontally later.

Rules:

```text
Do not depend on local process memory for long-term correctness in future phases.
Keep registry behind an interface.
Keep storage replaceable.
Do not use local filesystem for shared state.
Do not assume single instance except in Phase 1.
Do not use in-process locks as future distributed locks.
Do not generate IDs that can collide across instances.
Do not couple API handlers directly to map storage.
```

Phase 1 in-memory registry is acceptable only as an implementation bootstrap.

Every stateful component must have a future migration path to durable storage.

## 23. Serverless Readiness Rules

Sovrunn API code should be deployable later in serverless or container environments.

Rules:

```text
fast startup
config from environment or config file
no required local disk writes
no background jobs required for startup
graceful shutdown
stateless request handling where possible
bounded memory use
readiness endpoint
health endpoint
clear dependency initialization
```

Do not assume:

```text
long-lived local process state
writable local filesystem
fixed hostname
fixed local IP
manual startup order
single process forever
```

Avoid startup behavior that requires external services unless those dependencies are explicitly part of the current feature.

## 24. Graceful Shutdown

HTTP server must support graceful shutdown.

Use:

```go
srv.Shutdown(ctx)
```

Rules:

```text
handle SIGINT and SIGTERM
stop accepting new requests
allow in-flight requests to complete within timeout
flush logs if applicable
close resources explicitly
avoid goroutine leaks
```

Server shutdown timeout must be bounded.

## 25. Configuration Rules

Configuration must be explicit.

Use config file and environment override where useful.

Rules:

```text
no hardcoded production values
safe local defaults are allowed
validate config at startup
fail fast on invalid config
do not log secrets from config
prefer explicit config structs
document defaults
```

Configuration should support future container/serverless deployment.

Recommended sources:

```text
config file
environment variables
command-line flags for local development
```

Avoid hidden configuration from arbitrary files or current working directory assumptions.

## 26. Concurrency Rules

Use goroutines carefully.

Rules:

```text
do not start unbounded goroutines
do not leak goroutines
use context cancellation
use buffered channels only with clear bounds
do not use channels where mutexes are simpler
protect shared maps with mutexes
run tests with -race for concurrency-sensitive code
do not hold locks during slow I/O
do not call external systems while holding locks
```

Use `sync.RWMutex` for simple map-backed registries.

Use channels for coordination, not as a default abstraction.

## 27. Memory Rules

Avoid avoidable memory growth.

Rules:

```text
bound request body size
avoid storing large request bodies
avoid keeping references to mutable request buffers
copy data before storing if needed
do not expose internal maps or slices
preallocate result slices where possible
avoid unbounded caches
define cache invalidation before adding cache
```

Do not add caching just to improve perceived performance. Caching requires correctness rules.

## 28. Dependency Rules

Dependencies must be minimal and justified.

Before adding a dependency, AI must explain:

```text
what problem it solves
why standard library is insufficient
whether it is actively maintained
whether it impacts security
whether it is required for current feature
```

Avoid dependencies for:

```text
simple validation
simple routing
simple config parsing unless already approved
simple logging unless already approved
small helper functions
```

Do not add dependencies that introduce large transitive dependency trees without approval.

## 29. Logging Rules

Use structured logs.

Required fields for request logs:

```text
request_id
method
path
status_code
latency_ms
```

Error logs should include:

```text
request_id
error_code
safe_message
```

Rules:

```text
do not log secrets
do not log full request bodies by default
do not log authorization headers
do not log cookies
do not log private keys
do not log tokens
keep log volume bounded
avoid log spam in loops
```

Log messages should help operators debug issues without exposing sensitive data.

## 30. Testing Rules

Every feature must include tests.

Minimum tests:

```text
validation tests
registry tests
API handler tests
error mapping tests
delete-blocked tests where applicable
concurrency tests for registry where applicable
configuration validation tests where applicable
```

Run:

```bash
go test ./...
go test -race ./...
go vet ./...
```

Tests must be deterministic.

Avoid tests that depend on external services in Phase 1 unless clearly marked integration tests.

Do not remove tests to make a build pass.

## 31. Benchmarking Rules

For performance-sensitive paths, add benchmarks when appropriate.

Use:

```bash
go test -bench=. ./...
go test -bench=. -benchmem ./...
```

Benchmarks should avoid external services unless clearly marked.

Track:

```text
ns/op
B/op
allocs/op
```

Rules:

```text
Benchmark hot paths, not everything.
Do not optimize based on microbenchmarks alone.
Compare against realistic request flows when possible.
Keep benchmarks deterministic.
```

## 32. Static Analysis Rules

Use basic static checks before completion.

Required:

```bash
go fmt ./...
go test ./...
go vet ./...
```

Recommended for later phases:

```text
staticcheck
govulncheck
gosec
golangci-lint
```

Do not introduce a new mandatory tool into the build without approval.

## 33. API Compatibility Rules

Once an API contract is accepted, do not break it casually.

Rules:

```text
do not rename JSON fields without approval
do not change error codes without approval
do not change resource shape without approval
do not remove fields without approval
prefer additive changes
document breaking changes
```

Phase 1 may still evolve, but AI agents must flag contract changes explicitly.

## 34. Server Lifecycle Rules

The API server must have clear lifecycle ownership.

Rules:

```text
main wires dependencies
server owns HTTP lifecycle
handlers do request work
registry owns state
validation owns validation logic
operation package owns operation records
audit package owns audit records
```

Do not put business logic inside `main.go`.

`main.go` should be thin.

## 35. Health and Readiness Rules

Health and readiness endpoints must be simple and reliable.

Recommended endpoints:

```text
GET /healthz
GET /readyz
GET /version
```

Rules:

```text
/healthz reports process health.
/readyz reports readiness to serve traffic.
/version reports build/runtime version when available.
Do not make health checks expensive.
Do not call slow external systems from health checks unless explicitly required.
```

## 36. SDE Integration Boundary

SDE is a major capability inside Sovrunn, but Phase 1 platform core must not become SDE-only.

Rules:

```text
Do not place SDE-specific assumptions in generic platform packages.
Do not make generic ServiceInstance depend on SDE runtime concepts.
Keep SDE-specific code under clearly named SDE packages in future phases.
Keep platform core generic across service classes.
```

SDE concepts must not leak into Organization, Tenant, Project, ServiceClass, ServicePlan, or generic Operation implementation unless explicitly required.

## 37. ServiceOps Boundary

ServiceOps plugin execution is not part of Phase 1 unless explicitly approved.

Rules:

```text
Plugin and Capability may be modeled as resources.
Do not execute plugins yet.
Do not load dynamic plugins yet.
Do not call external plugin runtimes yet.
Do not build plugin marketplace yet.
```

Design must allow future ServiceOps execution but not implement it early.

## 38. Horizontal and Serverless Design Checklist

Before completing any Go feature, confirm:

```text
Can this work behind a load balancer later?
Can storage be replaced later?
Does the request path avoid long blocking operations?
Does it avoid local filesystem state?
Does it have bounded memory behavior?
Can it shut down gracefully?
Does it expose health/readiness where relevant?
Does it avoid global mutable state?
```

If the answer is no, document the limitation.

## 39. AI Agent Guardrails

AI must not:

```text
optimize before correctness
remove validation to improve benchmark
remove tests to pass build
introduce global state for convenience
hide errors
panic instead of returning errors
add async behavior without explaining lifecycle
add caching without invalidation rules
add goroutines without shutdown logic
add dependencies without justification
change architecture without approval
```

AI should prefer:

```text
clear correctness first
measured optimization second
targeted refactoring third
```

When uncertain, stop and ask for human decision.

## 40. Required Completion Report

For every Go feature, AI must report:

```text
feature implemented
files changed
new structs
new interfaces
new endpoints
validation added
tests added
performance considerations
security considerations
observability fields
commands run
known limitations
non-goals intentionally not implemented
```

## 41. Required Verification Commands

Before marking a Go feature complete, run:

```bash
make fmt
make test
make vet
```

For concurrency-sensitive changes, also run:

```bash
go test -race ./...
```

For performance-sensitive changes, also run:

```bash
go test -bench=. -benchmem ./...
```

Do not mark work complete if tests fail.

## 42. Final Principle

Sovrunn Go code must be boring, fast, secure, observable, horizontally scalable, serverless-ready, and easy for humans to review.

AI-generated Go code is acceptable only when it remains spec-first, test-gated, security-conscious, and terminal-verified.
