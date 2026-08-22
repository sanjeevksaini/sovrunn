# Implementation Plan

_FEATURE-0016 — Adapter Boundary and ExecutionTarget Qualification — Tasks Document_

## Overview

This is the Tasks stage for FEATURE-0016. It decomposes the approved design
(`design.md`) into twelve dependency-ordered, independently implementable and
testable tasks in package `internal/executiontarget/` plus five HTTP handlers,
route registration, and a feature-local conformance suite, followed by a final
non-task verification checkpoint. This revision
restructures the task graph (per founder-approved reviewer required changes) so
that: audit assembly exists and is committed before the sole-committer lifecycle
service publishes any create/qualify/retire/maintenance outcome; every
lifecycle/idempotency/scheduler integration point into `lifecycle.go` has an
exact owning task and correct dependency order; no unapproved 501 placeholder
route behavior is ever introduced; the pre-ServeMux guard is proven against a
bounded allowed-method-set invariant on each of the four exact F0016 path
patterns (five method/path registrations total), not only HEAD/PUT; the eight
closed local violation constants have an exact owning task; all 122 conformance
cases are enumerated without omission; and the final verification checkpoint is
strictly verification-only with no self-authorized edit escape hatch.

## 1. Identity, stage, and execution rules

### 1.1 Identity

| Field | Value |
|-------|-------|
| Feature | FEATURE-0016 — Adapter Boundary and ExecutionTarget Qualification |
| Stage | Tasks |
| Kiro slug | `adapter-boundary-and-executiontarget-qualification` |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase | 2R |
| Order | 6 |
| Tier | A (6–12 vertical slices) |
| Owned resources | `ExecutionTarget` (VS0-SCHEMA-015), `NormalizedTargetFactSet` (VS0-SCHEMA-016), `TargetQualificationResult` (VS0-SCHEMA-017) |
| Owned state machine | `VS0-STATE-004` |
| Owned writer | `VS0-WRITER-006` (`ExecutionTargetLifecycleService`) |
| Controlling decisions | DEC-0036, DEC-0042, DEC-0057 |
| Controlling handoffs | ADH-2026-025, ADH-2026-040, ADH-2026-042, ADH-2026-045, ADH-2026-058, ADH-2026-060, ADH-2026-061, ADH-2026-063, ADH-2026-064, ADH-2026-065, ADH-2026-066 |
| Local conformance | VS0-CF-F16-01..128 |

### 1.2 Stage inputs (consumed by reference, not redefined)

This tasks stage translates approved requirements
(`.kiro/specs/adapter-boundary-and-executiontarget-qualification/requirements.md`) and
approved design
(`.kiro/specs/adapter-boundary-and-executiontarget-qualification/design.md`) into
independently implementable and testable vertical slice tasks.

The approved requirements §4 (11 REQ, 12 AC, 122 conformance cases) and approved
design §2–10 (10 design decisions, component/responsibility/path, HTTP surface,
correctness properties, error precedence, security/observability/compatibility,
test strategy), together with the approved higher-precedence corrections
ADH-2026-060, ADH-2026-061, ADH-2026-063, ADH-2026-064, ADH-2026-065, and
ADH-2026-066 are the semantic authorities for this stage. ADH-2026-066 is a
private compatibility bridge, not a reopening of FEATURE-0015.

FEATURE-0012 (ObjectMeta, TypedRef, Problem Details, Condition, media,
If-Match/ETag, validation pipeline), FEATURE-0013 (AuditEvent atomic
append-before-publication), and FEATURE-0015 (CloudProviderParticipation and
InfrastructureStack read-only prior-feature authority, server-resolved
authorization, safe-denial, idempotency conventions) are reused by reference, not
forked.

### 1.3 Closed boundary for this stage

This tasks stage modifies only
`.kiro/specs/adapter-boundary-and-executiontarget-qualification/tasks.md`. It
decomposes approved design into independently implementable and testable vertical
slices. It introduces no semantic choice beyond the requirements/design approved
stop conditions. It emits no source, schema, architecture, prompt, validator,
automation, or test code.

Fail-closed stop conditions available to this stage (exactly the five
manifest-approved conditions; no other stop token may be invented or used):
`ARCHITECTURE_DECISION_REQUIRED`, `REQUIREMENT_CLARIFICATION_REQUIRED`,
`BOUNDARY_CHANGE_REQUIRED`, `DEPENDENCY_APPROVAL_REQUIRED`,
`SECURITY_REVIEW_REQUIRED`. If a task or its verification step discovers drift,
an unspecified observable outcome, or a needed edit outside its declared writable
paths, the task MUST report the exact failure and stop. It may name one of these
five conditions only when that condition's defined cause actually applies; it
must not relabel an ordinary verification failure as an architecture or security
stop condition.

### 1.4 Task execution rules

Each task below is independently implementable and testable in the order listed.
Each task cites exact requirements, design decisions, architecture decisions,
applicable risks, lists exact writable files, tests, verification commands,
acceptance criteria, security/observability impact, explicit exclusions, and an
exact single-change commit message. A task is complete only when its acceptance
criteria pass and the task-scoped verification commands pass. No task may begin
before its declared prerequisites are complete. The final repository-wide
verification checkpoint (not a numbered task) is always last and is strictly
verification-only.

**Reservation signaling protocol:** Every successful completion and every abort
(audit failure, stale fence, cancellation, panic, Maintenance/Retire win, or
shutdown) follows one exact two-phase sequence. While holding the
`ExecutionTargetLifecycleService` mutex, the service finalizes and detaches the
applicable reservation from in-process state and captures its `done` channel.
It then unlocks and closes that captured channel exactly once, making its waiters
runnable only after the protected state is stable. No channel is closed while
the lifecycle mutex is held.

### 1.5 Task graph summary (see also `## Task Dependency Graph` below)

Audit assembly (TASK-F16-02) exists and is complete before the sole-committer
lifecycle service (TASK-F16-04) is written, so every lifecycle commit method
(create, retire, qualify-commit, maintenance enter/clear) calls audit append
inside the same commit path from the moment that method is written — audit is
never bolted on after the fact. The qualification engine (TASK-F16-05) and the
maintenance/expiry scheduler (TASK-F16-06) each depend on the lifecycle service
and are each given `lifecycle.go` in their writable paths (with correct
prerequisite ordering) so they can integrate their own commit dispatch directly
into the sole-committer type rather than being blocked from doing so. The
idempotency table (TASK-F16-03) is built first with its package-private
completion/abort primitives. The lifecycle service (TASK-F16-04) invokes those
primitives atomically with publication without modifying `idempotency.go`. Route
registration for each of the five method/path registrations (across four path
patterns) is performed by the same task that implements that route's real
handler (TASK-F16-09/10/11), never by the transport guard task (TASK-F16-08),
so no placeholder 501 state is ever created or committed.

---

## Task Dependency Graph

The graph below is the corrected, acyclic execution-wave ordering. It resolves
the prior structural defect where audit assembly, idempotency integration, and
lifecycle-owned commit dispatch could not be reached from the tasks that needed
them. Audit (`TASK-F16-02`) and idempotency (`TASK-F16-03`) are built as
standalone components in wave 1, immediately after domain types/violations/store
(`TASK-F16-01`, wave 0). The sole-committer lifecycle service
(`TASK-F16-04`, wave 2) depends on both and owns `lifecycle.go`, so its
create commit method calls audit append and idempotency completion inside
the same commit path from the moment it is written; TASK-F16-04's retire
commit method, as implemented at this wave, handles only the Active-state
commit path with no Qualifying-reservation integration. The qualification
engine (`TASK-F16-05`, wave 3) depends on the lifecycle service and is also
given `lifecycle.go` in its writable paths to add the qualify-commit dispatch
method and to update/extend TASK-F16-04's retire-commit method with the
production integration that checks for and aborts a real in-flight Qualifying
reservation (retire remains valid from every Active combination, including
while Maintenance is active). The scheduler (`TASK-F16-06`, wave 4) depends on the lifecycle service and the
qualification engine (it aborts an in-flight `Qualifying` reservation for the
specific target receiving a Maintenance trigger, if any) and on idempotency
(it aborts that target's linked idempotency reservation only when one is
in flight, per ADH-2026-061 — no idempotency reservation is cancelled when no
qualification is in flight), and is also given `lifecycle.go` to add
maintenance-enter/clear commit methods.
Projection (`TASK-F16-07`, wave 5) depends on qualification and scheduler
outputs. The transport guard (`TASK-F16-08`) has no dependency on lifecycle
internals and runs in wave 1 alongside audit/idempotency. Route registration for
each of the five method/path registrations (across four path patterns) is
performed by the task that implements that route's real handler
(`TASK-F16-09`/`10`/`11`), sequenced across waves 6–8 so no two tasks touch
`internal/server/routes.go` in the same wave. The conformance
suite (`TASK-F16-12`, wave 9) depends on all implementation tasks. TASK-F16-11
also owns the one production composition point: it constructs and injects the
shared lifecycle service and scheduler, starts expiry and Maintenance-trigger
intake, installs the guard once around the completed F0016 mux, and invokes the
lifecycle shutdown hook before HTTP shutdown. The final
verification checkpoint is the terminal step: it is not itself a numbered task
or graph node, runs only after all 10 waves below (all twelve tasks) complete,
and depends on everything with no writable paths of its own.

```json
{
  "waves": [
    { "id": 0, "tasks": ["TASK-F16-01"] },
    { "id": 1, "tasks": ["TASK-F16-02", "TASK-F16-03", "TASK-F16-08"] },
    { "id": 2, "tasks": ["TASK-F16-04"] },
    { "id": 3, "tasks": ["TASK-F16-05"] },
    { "id": 4, "tasks": ["TASK-F16-06"] },
    { "id": 5, "tasks": ["TASK-F16-07"] },
    { "id": 6, "tasks": ["TASK-F16-09"] },
    { "id": 7, "tasks": ["TASK-F16-10"] },
    { "id": 8, "tasks": ["TASK-F16-11"] },
    { "id": 9, "tasks": ["TASK-F16-12"] }
  ]
}
```

The final verification checkpoint runs after wave 9 completes. It is
intentionally omitted from the JSON graph above because it is not a
dependency-ordered implementation task — it is the terminal gate described in
`## Final Verification Checkpoint` below.

Same-file write ordering enforced by the waves above:
- `internal/executiontarget/lifecycle.go`: written by TASK-F16-04 (wave 2),
  TASK-F16-05 (wave 3), TASK-F16-06 (wave 4) — three distinct waves.
- `internal/executiontarget/idempotency.go`: written only by TASK-F16-03 (wave
  1); TASK-F16-04 invokes its package-private primitives without editing it.
- `internal/server/routes.go`: written by TASK-F16-09 (wave 6), TASK-F16-10
  (wave 7), TASK-F16-11 (wave 8) — three distinct waves.
- `internal/server/server.go` and `cmd/sovrunn-api/main.go`: written only by
  TASK-F16-11 (wave 8) for final F0016 composition and ordered shutdown.

---

## Tasks

## 2. Dependency-ordered implementation tasks

### TASK-F16-01 — Domain types, closed violation constants, and in-memory store

**Prerequisites:** None (first task).

**Requirement authority:** REQ-F16-01, REQ-F16-02, REQ-F16-10 (domain types,
store, closed violation set).

**Design authority:** DD-01 (domain service package placement), design §3.1
(new F0016-owned components), design §4.1–4.4 (ExecutionTarget /
NormalizedTargetFactSet / TargetQualificationResult data models, provenance),
design §5.2 (closed local violation set, exactly eight new codes).

**Architecture decisions:** DEC-0042 (ExecutionTarget canonical owner),
ADH-2026-058 clauses 1, 7, 10 (closed create fields, no persisted
status.availability, response-only effectiveAvailability, eight new violation
codes only).

**Applicable risks:** Risk 2 (no external effect, proven by restart clearing
F0016-owned in-memory state; design §6.5).

**Writable paths:**
- `internal/executiontarget/model/types.go`
- `internal/executiontarget/model/enums.go`
- `internal/executiontarget/violations.go`
- `internal/executiontarget/store.go`
- `internal/executiontarget/store_test.go`
- `internal/executiontarget/violations_test.go`

**Tests:**
Unit tests for ExecutionTarget in-memory struct, NormalizedTargetFactSet,
TargetQualificationResult, viability fingerprint, current-Maintenance marker,
scope-unique name reservation (retained after Retired), live tuple
`(participation UID, stack UID, targetClass)` at-most-one index, and
package-private staging/index primitives. They do not publish state. Unit tests
prove model initialization values for a lifecycle-owned create proposal:
derived CloudProvider scope, Active/Unqualified, observedGeneration, and
maintenanceEpoch=0. Unit tests proving `violations.go` defines
exactly the eight closed local violation constants (`VS0_TARGET_RETIRED`,
`VS0_TARGET_MAINTENANCE`, `VS0_TARGET_QUALIFICATION_IN_PROGRESS`,
`VS0_TARGET_EPOCH_STALE`, `VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE`,
`VS0_EXECUTION_TARGET_STACK_UNAVAILABLE`, `VS0_EXECUTION_TARGET_SCOPE_MISMATCH`,
`VS0_EXECUTION_TARGET_VIABILITY_STALE`) as string constants with no other new
code, and that no inherited violation code (`VS0_STATUS_FIELD_WRITE`,
`VS0_SYSTEM_OWNED_FIELD_WRITE`, `VS0_AUTHORIZATION_SAFE_DENIAL`,
`VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`) is redefined here.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -race ./internal/executiontarget/...
```

**Acceptance criteria:**
- ExecutionTarget, NormalizedTargetFactSet, TargetQualificationResult types exist
  with the exact required fields per design §4.1–4.3.
- `internal/executiontarget/violations.go` defines exactly the eight closed local
  violation constants per design §5.2; no other F0016-local violation code
  exists; no inherited code is redefined.
- Scope-unique name reservation is enforced (retained after retirement).
- At most one live non-Retired `(participation UID, stack UID, targetClass)` tuple
  is indexed.
- Store mutation helpers are package-private staging/index primitives; they do
  not publish a target or persist lifecycle state independently. TASK-F16-04
  invokes them only while its lifecycle-service mutex is held, establishing the
  sole publication authority and observable create behavior.

**Security impact:** Establishes the closed violation surface consumed by every
later handler task; no client status write occurs at this task (no handlers
exist yet).

**Observability impact:** None at this task (audit and projection tasks follow).

**Exclusions:** No lifecycle service, no HTTP handlers, no idempotency, no
observer, no audit, no qualification engine, no HTTP routes, no pre-ServeMux
guard, no projection, no maintenance scheduler, no expiry worker. These are
future tasks.

**Commit message:**
```
feat(f16): TASK-F16-01 — domain types, closed violation constants, in-memory store

Implements REQ-F16-01, REQ-F16-02, REQ-F16-10 (domain types, store, closed
violation set).

Creates ExecutionTarget, NormalizedTargetFactSet, TargetQualificationResult
in-memory types with exactly the approved VS0-SCHEMA-015..017 fields; scope-unique
name reservation (retained after retirement); at-most-one live tuple index;
Active/Unqualified initial state at epoch 0.

Creates internal/executiontarget/violations.go defining exactly the eight closed
local violation constants per design §5.2; no other F0016-local code; no
inherited code redefined.

No lifecycle service, HTTP, observer, idempotency, audit, projection, or
qualification engine.

Relates: FEATURE-0016, VS0-SCHEMA-015..017, design §5.2.
```

---

### TASK-F16-02 — Build-only AuditEvent assembly

**Prerequisites:** TASK-F16-01 complete.

**Requirement authority:** REQ-F16-09 (audit matrix, append-failure behavior),
requirements §5.4 (audit inherited FEATURE-0013 semantics).

**Design authority:** Design §3.1 (audit.go), design §6.3 (observability audit
matrix and append-failure), design §6.4 (compatibility reuses FEATURE-0013).

**Architecture decisions:** ADH-2026-058 clause 8 (audit matrix for
success/denial/safe-denial, expiry sole background exception).

**Applicable risks:** Risk 3 (audit-before-publication: design §5.3 + §6.3 +
VS0-CF-F16-39,40,82,99..104; the assembly function constructs an in-memory
AuditEvent, then its caller attempts the inherited append boundary. A failed
append has no durable AuditEvent and no publication; the lifecycle-level guarantee that a failed append
leaves target/ETag unmutated, removes the linked InFlight idempotency
reservation, and creates no completion is proven by TASK-F16-04, which
integrates this assembly into the sole-committer commit path).

**Writable paths:**
- `internal/executiontarget/audit.go`
- `internal/executiontarget/audit_test.go`

**Tests:**
Unit tests for audit assembly building redacted FEATURE-0013 AuditEvents for:
successful create, completed qualification, retirement, maintenance enter/clear,
authenticated authorization denial, safe denial, and expiry. Assembly for each
of these (including the function that builds the expiry AuditEvent struct) is
a distinct, separately-testable function that only assembles the event; the
caller (the lifecycle service for state transitions, the owning handler for
audited denials, or the TASK-F16-06 scheduler for expiry) is responsible for one append of the assembled
event, in exactly this order: assemble the event first, then append it. This
task owns only build-only assembly functions. Exactly one append owner is the
lifecycle service for state transitions, TASK-F16-06 scheduler for expiry, or
the owning HTTP handler for authorization/safe-denial events; assembly never
calls, owns, or retries append.

Request-triggered required append failure: the assembly function may have
constructed and handed an in-memory event to the append boundary; the append
attempt fails with no event durably appended or persisted and no publication.
The error never discloses suppressed
403/404 details. The lifecycle-level guarantee that this failure translates
to `INTERNAL_ERROR`/500 with the target and ETag left unmutated, the linked
InFlight idempotency reservation removed (not left unchanged), no completion
created, and any waiters woken only after the sole-committer mutex is released
(unlock) is proven by TASK-F16-04, which integrates this assembly into the
sole-committer commit path.

Expiry AuditEvent assembly: this task provides only the assembly function that
builds the expiry AuditEvent struct, called by the TASK-F16-06 scheduler on
each retry attempt; on an assemble-then-append failure for expiry, no expiry
event is appended or persisted. This task does not own, describe, or test the
timing or cadence of expiry retries (e.g. "retrying once per injected-clock
second") — that scheduling behavior belongs solely to TASK-F16-06.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -race ./internal/executiontarget/...
```

**Acceptance criteria:**
- Audit assembly functions are build-only: they build redacted FEATURE-0013 AuditEvents for the exact
  audit matrix per REQ-F16-09, including a distinct expiry-event assembly
  function that builds the expiry AuditEvent struct; each assembly function
  returns it to the lifecycle service, scheduler, or owning HTTP handler. That
  caller makes exactly one append attempt — assembly and append are explicit,
  ordered, separate steps with one append owner.
- Request-triggered required append failure: an in-memory AuditEvent may be
  constructed and passed to the append boundary; append failure leaves no
  durable AuditEvent and no publication. The error never discloses
  suppressed 403/404. The lifecycle-level outcome (`INTERNAL_ERROR`/500,
  target/ETag left unmutated, the linked InFlight idempotency reservation
  removed rather than left unchanged, no completion, waiters woken only after
  unlock) is proven by TASK-F16-04, which integrates this assembly into the
  sole-committer commit path — this task does not itself own or touch a
  target, ETag, or idempotency reservation.
- The expiry-event assembly function is called by the TASK-F16-06 scheduler on
  each retry attempt; on assemble-then-append failure, no expiry event is
  appended or persisted. This task owns only the assembly function, not the
  retry timing/cadence — TASK-F16-06 is the sole owner of when and how often
  an expiry retry occurs (e.g. "retrying once per injected-clock second").
- Unit tests prove each assembly function constructs the exact redacted event;
  TASK-F16-04/05/06 prove lifecycle/scheduler caller append behavior and
  TASK-F16-09/10/11 prove exactly one handler-side append attempt for their
  authorization and safe-denial events.

**Security impact:** Establishes build-only audit assembly, preventing a helper
from double-appending or omitting ownership of the security-critical append;
the unit-level half of the
audit-before-publication barrier (Risk 3); the lifecycle-level barrier
preventing partial target/ETag/idempotency mutation on append failure is proven
by TASK-F16-04, which integrates this assembly into the commit path.

**Observability impact:** Provides the audit matrix (including expiry-event
assembly) for all F0016 success/denial outcomes; the retry cadence that
eventually ensures audit without blocking a caller is owned solely by
TASK-F16-06, which calls this task's expiry-event assembly function on each
attempt.

**Response boundary:** An AuditEvent is an appended side effect only. It is never
a member of the closed ExecutionTarget HTTP response body; handlers return only
the approved target projection or Problem envelope.

**Exclusions:** No lifecycle service, no HTTP handlers, no routes, no
idempotency handlers, no retry timing/scheduling of any kind (owned solely by
TASK-F16-06). These are future tasks or other tasks' ownership. This task
provides the audit assembly consumed by the lifecycle service task
(TASK-F16-04), the scheduler task (TASK-F16-06), and later handler tasks — it
exists and is complete *before* the lifecycle service is written, so
audit-before-publication is built into the commit path from its first line, not
bolted on afterward.

**Commit message:**
```
feat(f16): TASK-F16-02 — build-only AuditEvent assembly

Implements REQ-F16-09 (audit matrix, append-failure behavior).

Creates audit assembly functions building redacted FEATURE-0013 AuditEvents
for: successful create, completed qualification, retirement, maintenance
enter/clear, authenticated authorization denial, safe denial, and expiry
(including a distinct expiry-event assembly function that builds the expiry
AuditEvent struct). Each function is build-only and returns the event to its
lifecycle-service, scheduler, or owning HTTP-handler caller, which performs
exactly one inherited append attempt. TASK-F16-09/10/11 prove the handler-owned
authorization/safe-denial append paths. Request-triggered required append failure: the caller attempts
append after assembly; append failure leaves no durable AuditEvent or
publication (unit-level property); never discloses suppressed
403/404. The expiry-event assembly function is called by TASK-F16-06's
scheduler on each retry attempt; this task owns only the assembly, not the
retry timing/cadence, which belongs solely to TASK-F16-06.

Proves the unit-level half of Risk 3 (assemble-then-append discipline of each
assembly function alone). The lifecycle-level guarantee that a failed append
leaves target/ETag unmutated, removes the linked InFlight idempotency
reservation (not left unchanged), creates no completion, and wakes waiters
only after unlock, returning INTERNAL_ERROR/500, is proven by TASK-F16-04,
which integrates this assembly into the sole-committer commit path. Exists
before the sole-committer lifecycle service (TASK-F16-04) so every commit
method calls audit append inside the same commit path.

No lifecycle service, HTTP, routes, idempotency handlers, or retry
timing/scheduling.

Relates: FEATURE-0016, REQ-F16-09, ADH-2026-058, FEATURE-0013 reuse.
```

---

### TASK-F16-03 — Idempotency reservation table and waiter management

**Prerequisites:** TASK-F16-01 complete.

**Requirement authority:** REQ-F16-06 (idempotency namespace, digest, retention,
isolation).

**Design authority:** DD-05 (F0016-owned in-process idempotency table), design §3.1
(idempotency.go), design §5.1 (request precedence step 9), design §6.5
(operational 24h/10,000 retention).

**Architecture decisions:** ADH-2026-058 clause 5 (idempotency namespace, digest,
applies to create/qualify/retire only), ADH-2026-054 (waiter recheck of current
authorization + safe access).

**Applicable risks:** None direct; supports the lifecycle service integration in
TASK-F16-04.

**Writable paths:**
- `internal/executiontarget/idempotency.go`
- `internal/executiontarget/idempotency_test.go`

**Tests:**
Unit tests for F0016-owned namespace `(principal UID, route pattern, derived
CloudProvider UID, concrete action-target UID when applicable, Idempotency-Key)`;
SHA-256 digest of recursively sorted allowed create JSON or SHA-256 of zero bytes
for actions; `If-Match`, `Accept`, authentication, request/correlation IDs
excluded; InFlight and Completed entry states; InFlight owns one unbuffered `done`
channel per reservation, with package-private completion/abort helpers callable
only by the lifecycle service. The service finalizes/detaches under its mutex,
captures the channel, then unlocks before exactly-once close; waiter selects its reservation `done` channel or
request context; cancellation only detaches waiter; retention 24h and 10,000
completed entries in-process; eviction at-expiry then earliest-expiry/lexical
namespace at capacity; InFlight never evicts; independent reservations across
different action targets and derived CloudProvider scopes; same principal/key/digest
on another target/scope forms independent reservation.

The table has no independent mutex and handlers never mutate it directly.
`ExecutionTargetLifecycleService` owns the exact protocol under its single
DD-02 mutex: reserve-or-inspect-replay, attach/detach waiter, finalize
completion, and abort. After wake, the **handler** re-enters at authentication
and repeats authorization, safe access, validation, and reservation with no
lifecycle mutex held. TASK-F16-03
provides only package-private primitives used by those service methods.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -race ./internal/executiontarget/...
```

**Acceptance criteria:**
- F0016-owned idempotency table exists with exact namespace, digest algorithm, and
  retention/eviction rules per REQ-F16-06.
- InFlight reservation owns one unbuffered `done` channel. Package-private
  completion/abort helpers let the lifecycle service finalize and detach it
  under its mutex, capture the channel, unlock, then close it exactly once;
  no handler or caller receives an exported reservation-close capability.
- Waiter selects `done` or request context; cancellation detaches only.
- 24h/10,000-entry completed retention; at-expiry then
  earliest-expiry/lexical-namespace eviction; InFlight never evicts.
- Independent reservations across action targets and derived scopes proven by unit
  tests.
- The table has no independent lock; lifecycle-service methods exclusively own
  reservation acquisition, replay lookup, waiter attachment/detachment,
  completion, and abort under the single DD-02 mutex. A woken handler re-enters
  at authentication with no lifecycle mutex held.

**Security impact:** Waiter recheck of current authorization + safe access after
wake-up (ADH-2026-054) prevents stored-result disclosure when grants change; this
task provides the table, the actual owner-abort/wake-waiters wiring into the
lifecycle service commit path is completed in TASK-F16-04.

**Observability impact:** None at this task.

**Exclusions:** No lifecycle service, no HTTP handlers, no routes, no audit, no
projection. These are future tasks. This task provides the idempotency table
consumed by the lifecycle service task (TASK-F16-04) and the HTTP handler tasks.

**Commit message:**
```
feat(f16): TASK-F16-03 — idempotency reservation table and waiter management

Implements REQ-F16-06 (idempotency namespace, digest, retention, isolation).

Creates F0016-owned in-process idempotency table with exact namespace (principal
UID, route pattern, derived CloudProvider UID, concrete action-target UID when
applicable, Idempotency-Key); SHA-256 digest of recursively sorted allowed create
JSON or zero bytes for actions; InFlight/Completed states; InFlight owns one
unbuffered done channel with package-private exactly-once completion/abort helpers
callable only by the lifecycle service under its mutex; waiter selects
done or request context; cancellation detaches only; 24h/10,000-entry completed
retention; at-expiry then earliest-expiry/lexical-namespace eviction; InFlight
never evicts; independent reservations across action targets and derived scopes.

Supports ADH-2026-054 waiter-recheck security rule. Owner-abort/wake-waiters
wiring into the sole-committer lifecycle service commit path is completed in
TASK-F16-04.

No lifecycle service, HTTP, audit, or projection.

Relates: FEATURE-0016, REQ-F16-06, ADH-2026-058, ADH-2026-054.
```

---

### TASK-F16-04 — Sole-committer lifecycle service (create, retire, audit and idempotency integration)

**Prerequisites:** TASK-F16-01, TASK-F16-02, TASK-F16-03 complete.

**Requirement authority:** REQ-F16-01, REQ-F16-02, REQ-F16-08 (sole committer +
lifecycle), REQ-F16-09 (audit matrix for create/retire).

**Design authority:** DD-02 (sole committer service type), design §3.1
(lifecycle.go), design §4.1–4.4 (data models, provenance), design §5.3
(deterministic single-writer commit, audit-before-publication).

**Architecture decisions:** DEC-0042 (ExecutionTarget canonical owner), DEC-0057
(qualification/availability), ADH-2026-058 clauses 1, 5, 7, 8 (closed create
fields, sole committer, no persisted status.availability, response-only
effectiveAvailability, idempotency completion tied to publication, audit matrix).

**Applicable risks:** Risk 1 (sole-committer status writes must exclude clients;
design enforcement: DD-02, VS0-WRITER-006); Risk 3 (audit-before-publication is
now built into the commit path from this task's first line, because TASK-F16-02
already exists); Risk 2 (no external effect, proven by restart clearing
F0016-owned in-memory state; design §6.5).

**Writable paths:**
- `internal/executiontarget/lifecycle.go`
- `internal/executiontarget/lifecycle_test.go`

**Tests:**
Unit tests proving `ExecutionTargetLifecycleService` exists with a sole-committer
mutex guarding all F0016 targets/indexes/internal records/Maintenance
markers/reservations (DD-02). Create commit method: validates the closed
create-field boundary, derives CloudProvider scope, persists Active/Unqualified
at epoch 0, obtains the event from TASK-F16-02's build-only assembly, performs
exactly one inherited append attempt *before* publication, and only then publishes the target and marks
the idempotency reservation (TASK-F16-03) complete via its package-private
completion helper — a failed required audit append leaves no target, ETag,
scope-name reservation, live tuple index, current links, or marker mutation,
removes the linked InFlight idempotency reservation (per design.md: audit
failure is one of the removal-triggering conditions alongside owner
cancellation, panic, stale-fence failure, and shutdown), creates no completion
record, wakes any waiters only after the sole-committer mutex is released
(unlock), and returns `INTERNAL_ERROR`/500. Retire commit method, AS
IMPLEMENTED BY THIS TASK: allowed from every Active combination including
while Maintenance is active, handling only the Active-state commit path with
no Qualifying-reservation integration — this task's retire-commit method does
not call into or trigger any mechanism owned by TASK-F16-05. Its own
audit-before-publication and idempotency-completion discipline (including the
failed-append reservation-removal behavior) is fully provable now with a
synthetic in-flight-reservation test double that never calls a real abort
mechanism. The production integration of retire against an actual in-flight
Qualifying reservation (retire's commit checking for and aborting a real
in-flight Qualifying reservation, with the qualifying caller receiving 409
CONFLICT + `VS0_TARGET_RETIRED`) is added later by TASK-F16-05, which updates
this method in `lifecycle.go`. Tests prove sole-committer write path and
clients-never-write-status invariant.

This task owns the lifecycle-service callable reservation API used later by
handlers: `ReserveOrInspectReplay`, `AttachWaiter`, `DetachWaiter`,
`AwaitReservation`, and bounded `AbortReservation`. `CompleteReservation` is
commit-internal: only create/qualify/retire audited publication paths finalize a
Completed record after successful append and publication. Each API operation
holds the one DD-02 mutex only for table operations; handlers perform post-wake
authentication and repeat the request pipeline outside that mutex.

This task also owns create and retire owner-cancellation/panic cleanup. Tests
prove each path finalizes and detaches its linked InFlight reservation without a
completion or publication, unlocks, then closes its captured `done` channel once.
This task uses only synthetic/no-qualification retirement tests. TASK-F16-05
owns the real in-flight-qualification retirement append-failure proof.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -race ./internal/executiontarget/...
```

**Acceptance criteria:**
- `ExecutionTargetLifecycleService` exists with sole-committer mutex guarding all
  F0016 targets/indexes/internal records/Maintenance markers/reservations (DD-02).
- Lifecycle-service reservation API owns reserve/replay inspection, waiter
  attach/detach/wait, and bounded abort under the one DD-02 mutex; handlers
  cannot mutate or complete the table. Completion is commit-internal after
  audit append and publication; handlers re-enter authentication after wake.
- Create commit method obtains the required AuditEvent from TASK-F16-02 and
  performs exactly one inherited append attempt before publication; it then marks the idempotency reservation (via TASK-F16-03
  table) complete only after publication; a failed required append leaves
  target, ETag, scope-name reservation, live tuple index, current links, and
  marker state unmutated, removes the linked InFlight idempotency
  reservation (it is not left in an unchanged InFlight state), creates no
  completion, and wakes any waiters only after the sole-committer mutex is
  released (unlock) — this is the lifecycle-level integration proof for the
  unit-level property TASK-F16-02 establishes for the assembly function alone.
- Retire commit method, AS IMPLEMENTED BY THIS TASK, has the same
  audit-before-publication, idempotency-completion, and failed-append
  reservation-removal discipline: an append failure preserves the existing
  FactSet/Result links and the live tuple while publishing neither retirement
  nor tuple release. It handles only the Active-state commit path
  (including while Maintenance is active) with no Qualifying-reservation
  integration; it does not call into or trigger any mechanism owned by
  TASK-F16-05. The production integration where retire's commit checks for and
  aborts a real in-flight Qualifying reservation, and the qualifying-caller-side
  409 CONFLICT + `VS0_TARGET_RETIRED` outcome, is added later by TASK-F16-05,
  which updates this method in `lifecycle.go`.
- Scope-unique name reservation is enforced (retained after retirement).
- At most one live non-Retired `(participation UID, stack UID, targetClass)` tuple
  is indexed.
- Create and retire owner cancellation/panic clean up only their linked InFlight
  reservation: finalize/detach under lock, no completion/publication, unlock,
  then close the captured `done` channel exactly once. Real in-flight
  qualification interaction is deferred to TASK-F16-05.
- Tests prove sole-committer write path and clients-never-write-status invariant.

**Security impact:** Establishes sole-committer status-write barrier
(VS0-WRITER-006) preventing client status mutations (Risk 1); establishes
audit-before-publication (Risk 3) from the first commit method written, with a
failed append removing the linked idempotency reservation and waking waiters
only after unlock rather than leaving state unchanged.

**Observability impact:** Create and retire are audited via the matrix in
TASK-F16-02, wired into the commit path here.

**Exclusions:** No qualify-commit dispatch (added by TASK-F16-05 into this same
file), no maintenance enter/clear commit methods (added by TASK-F16-06 into this
same file), no HTTP handlers, no HTTP routes, no pre-ServeMux guard, no
projection, no observer. These are future tasks.

**Commit message:**
```
feat(f16): TASK-F16-04 — sole-committer lifecycle service (create, retire, audit/idempotency integration)

Implements REQ-F16-01, REQ-F16-02, REQ-F16-08, REQ-F16-09 (sole committer,
lifecycle, audit matrix for create/retire).

Creates ExecutionTargetLifecycleService as sole committer (VS0-WRITER-006) with
one service-owned mutex guarding all F0016 targets/indexes/internal
records/Maintenance markers/reservations; scope-unique name reservation
(retained after retirement); at-most-one live tuple index; Active/Unqualified
initial state at epoch 0.

Create and retire commit methods obtain events from TASK-F16-02 build-only
assembly, each perform exactly one inherited append attempt before publication, and mark the TASK-F16-03 idempotency
reservation complete only after publication via its package-private exactly-once helper;
a failed required append leaves target, ETag, scope-name reservation, live tuple
index, current links, and marker state unmutated, removes the linked
InFlight idempotency reservation (it is not left unchanged), creates no
completion, wakes any waiters only after the sole-committer mutex is released
(unlock), and returns INTERNAL_ERROR/500. Retire commit method, as implemented
by this task, is allowed from every Active combination including while
Maintenance is active, and handles only the Active-state commit path with no
Qualifying-reservation integration — it does not call into or trigger any
mechanism owned by TASK-F16-05. The production integration where retire's
commit checks for and aborts a real in-flight Qualifying reservation (with the
qualifying caller receiving 409 CONFLICT + VS0_TARGET_RETIRED) is added later
by TASK-F16-05, which updates this method in lifecycle.go.

Proves Risk 1 (clients-never-write-status), Risk 3 (audit-before-publication
built in from this task's first commit method, with failed-append reservation
removal rather than unchanged-state retention), Risk 2 (no external effect).

No qualify-commit dispatch, maintenance enter/clear, HTTP, observer, or
Qualifying-reservation integration in retire (added by TASK-F16-05/06 into
this same file; TASK-F16-09/10/11 for HTTP).

Relates: FEATURE-0016, VS0-WRITER-006, VS0-STATE-004, REQ-F16-09.
```

---

### TASK-F16-05 — Synthetic observer and qualification engine (commit dispatch integrated into lifecycle service)

**Prerequisites:** TASK-F16-01, TASK-F16-04 complete.

**Requirement authority:** REQ-F16-02 (lifecycle current-viability recheck
before Qualifying work), REQ-F16-07 (synthetic observer, four named facts,
missing-fixture/timeout), REQ-F16-08 (qualification outcomes, fenced commit,
atomic record replacement, no Qualifying projection), REQ-F16-09 (required
qualification AuditEvent append before publication).

**Design authority:** DD-03 (synthetic observer as target-bound fixture reader),
DD-04 (synchronous qualification with captured-fence reservation), design §3.1
(observer.go, qualify.go), design §4.2 (NormalizedTargetFactSet facts), design §4.4
(fixed provenance), design §5.1 (request precedence step 12 fence recheck).

**Architecture decisions:** DEC-0036 (implementation-neutral adapter boundary),
DEC-0042 (no real external call), ADH-2026-058 clauses 6, 7 (sole observer,
profile evaluation, Qualifying reservation), ADH-2026-060 (Maintenance-race
outcome ordering).

**Applicable risks:** Risk 2 (no external effect: the synthetic observer is
external-effect-free, reads in-memory fixture, never waits on network at the
unit level; design §6.5. The full process-restart proof (VS0-CF-F16-46) is
proven by TASK-F16-12's feature-local conformance suite); Risk 6 (stale-fence
non-publication).

**Writable paths:**
- `internal/executiontarget/observer.go`
- `internal/executiontarget/qualify.go`
- `internal/executiontarget/lifecycle.go` (adding the qualify-commit dispatch
  method that the lifecycle service's sole-committer mutex guards, including
  the Qualifying reservation's own abort logic; AND updating/extending the
  retire-commit method that TASK-F16-04 added, so that its commit path checks
  for and aborts a real in-flight Qualifying reservation for the same target
  via the abort logic added here — this task is explicitly allowed to
  update/extend TASK-F16-04's retire-commit method in this same file, not only
  add a new qualify-commit dispatch method; no change to the create-commit
  method owned by TASK-F16-04)
- `internal/executiontarget/observer_test.go`
- `internal/executiontarget/qualify_test.go`
- `internal/executiontarget/fixtures_test.go` (test-only target-bound fixtures)

**Tests:**
Unit tests for `sovrunn.synthetic-iaas-observer/v1` reading target-bound in-memory
fixtures; four named facts (`compute.vm`, `storage.block`, `storage.object`,
`network.private`) each Supported/Unsupported/Unknown; closed profile evaluation:
any Unsupported → Rejected, else any Unknown → Indeterminate, else Qualified;
missing fixture and logical timeout each yield four Unknown facts at `now` /
`now+60s` with no network wait; duplicate/missing/malformed fact is an observer
fault (faults qualification, publishes no result); fixed provenance
(`observerID=sovrunn.synthetic-iaas-observer/v1`,
`profileVersion=synthetic-iaas/v1`, `factSchemaVersion=v1`, and registered
`observerRevision`); the selected fixture revision is immutable for the
qualification attempt and is retained as FactSet provenance. The persisted FactSet
version field is `factVersion=v1`.

The action handler in TASK-F16-11 owns caller-facing REQ-F16-02 safe backing
resolution and admission, including its exact 404/409/422 responses, before
idempotency replay/reservation. This task does not repeat that public safe-access
pipeline. Before it opens internal Qualifying work, however, the lifecycle
service rechecks under its mutex that the already-resolved participation remains
effective Active, the InfrastructureStack remains Active, and their derived
CloudProvider scopes still match. A changed backing state rejects without opening
the Qualifying reservation, observation, AuditEvent, publication, or completion.
Because this happens after the action idempotency reservation exists, this
TASK-F16-05 lifecycle recheck is the sole cleanup owner: it detaches the linked
InFlight action reservation under the lifecycle mutex, creates no completion,
unlocks, then closes its captured `done` channel exactly once. The handler only
maps the result; it does not perform a second abort. Tests cover ineffective
participation and inactive InfrastructureStack (`409`) and authorized scope
mismatch (`422`), proving no observer/AuditEvent/publication, one post-unlock
signal, and waiter re-entry for each result.

The action idempotency reservation is distinct from the internal Qualifying
reservation. After the handler has admitted the request and obtained/inspected
the idempotency reservation, the lifecycle service captures the internal
Qualifying reservation and its fences (InfrastructureStack generation,
maintenance epoch, viability fingerprint) under the mutex, unlocks for observer
execution, then reacquires the mutex in a separate qualify-commit dispatch and
applies the exact ordered commit predicate (ADH-2026-060): if active
current-Maintenance marker exists → abort with `CONFLICT`/409 +
`VS0_TARGET_MAINTENANCE`; else if captured maintenance epoch differs from current
→ abort with `STALE_RESOURCE_VERSION`/412 + `VS0_TARGET_EPOCH_STALE`; else if
InfrastructureStack generation changed → abort with `STALE_RESOURCE_VERSION`/412 +
`VS0_TARGET_EPOCH_STALE`; else if viability fingerprint changed → abort with
`STALE_RESOURCE_VERSION`/412 + `VS0_EXECUTION_TARGET_VIABILITY_STALE`; else
atomically replace both current record refs, refresh observedGeneration, advance
resourceVersion, obtain the required AuditEvent from TASK-F16-02 build-only
assembly, perform exactly one inherited append attempt before publication, and mark the TASK-F16-03 idempotency reservation
complete only after publication. A failed required audit append leaves the
target and ETag unmutated, removes the linked InFlight idempotency reservation
(TASK-F16-03 table; it is not left in an unchanged InFlight state), creates no
completion, and wakes any waiters only after the sole-committer mutex is
released (unlock).

Observer fault, panic, shutdown, cancellation, and every stale fence abort the
Qualifying reservation and remove the linked InFlight idempotency reservation
without target mutation, AuditEvent, or completion, waking any waiters only
after unlock.

This task owns and unit-tests the Retire-wins rule, and this task is the one
that wires it into production: it updates/extends the retire-commit method
that TASK-F16-04 added in `lifecycle.go` so that, when retire commits while a
Qualifying reservation is in flight for the same target (retire is valid from
every Active combination, including while Maintenance is active), that commit
now checks for and calls into the reservation-abort mechanism added here. That
Qualifying reservation aborts (no qualification result, completion, or
qualification AuditEvent) and the qualifying caller's outcome is 409 CONFLICT +
`VS0_TARGET_RETIRED`; the retirement itself is audited via TASK-F16-04's
retire-commit path, not by this task. This is the unit-level proof of the
rule and the production integration point; TASK-F16-11 later re-exercises it
end-to-end over the live HTTP retire and qualify actions as an integration
re-test, citing this task as the unit-level owner.

Tests also prove that a failed retirement AuditEvent append while a real
Qualifying reservation exists does not commit retirement and leaves that
qualification/reservation intact; only a successful retirement applies the
Retire-wins abort.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -race ./internal/executiontarget/...
```

**Acceptance criteria:**
- `sovrunn.synthetic-iaas-observer/v1` exists as a target-bound, clock-driven,
  external-effect-free component with no network I/O.
- Four named facts are produced with Supported/Unsupported/Unknown truth values.
- Closed profile evaluation (any Unsupported → Rejected; else any Unknown →
  Indeterminate; else Qualified) is proven by unit tests.
- Missing fixture and logical timeout each yield four Unknown facts at `now` /
  `now+60s` with no elapsed network wait.
- Duplicate/missing/malformed facts fault qualification and publish no result.
- FactSet provenance includes the registered `observerRevision`; the selected
  fixture revision is immutable for the attempt and cannot be replaced during
  observation or commit.
- The lifecycle service rechecks current backing viability under its mutex
  before opening internal Qualifying work; caller-facing safe backing admission
  and its exact 404/409/422 outcomes remain TASK-F16-11 ownership. If this
  post-idempotency recheck rejects participation/stack/scope viability, this
  task alone detaches the linked InFlight action reservation under the mutex,
  creates no completion/observation/AuditEvent/publication, unlocks, and closes
  its captured `done` channel exactly once; tests prove each returned 409/422
  outcome wakes waiters for full re-entry without stranding a reservation.
- The lifecycle service's qualify path (added to `lifecycle.go` here) keeps the
  action idempotency reservation distinct from the internal Qualifying
  reservation: it captures the latter and fences under lock, observes unlocked,
  then commits in a separate relocked dispatch. It
  applies the exact ADH-2026-060 ordered commit predicate, appends the required
  AuditEvent before publication, marks idempotency complete only after
  publication, and atomically replaces both record refs on success. A failed
  required append leaves target/ETag unmutated, removes the linked InFlight
  idempotency reservation, creates no completion, and wakes waiters only
  after unlock.
- Observer fault, stale fence, panic, shutdown, and cancellation each abort
  the Qualifying reservation and remove the linked InFlight idempotency
  reservation without target mutation, AuditEvent, or completion, waking
  waiters only after unlock.
- Retire-wins is owned and unit-proved here, and this task performs the
  production wiring: it updates/extends TASK-F16-04's retire-commit method in
  `lifecycle.go` so that method now calls into this task's reservation-abort
  mechanism while a Qualifying reservation is in flight for the same target
  (retire remains valid from every Active combination, including while
  Maintenance is active); that reservation aborts with no qualification
  result, completion, or qualification AuditEvent, and the qualifying caller's
  outcome is 409 CONFLICT + `VS0_TARGET_RETIRED`; the retirement itself is
  audited via TASK-F16-04.
- `externalCallCount=0` is proven at the unit level for this task's own
  observer/qualify code (no external observer call on any invocation); the
  full process-restart proof (VS0-CF-F16-46) is proven by TASK-F16-12's
  feature-local conformance suite, which this task's design makes possible by
  keeping all state in-process and external-effect-free.

**Security impact:** None (internal only; no network exposure or external effect).

**Observability impact:** Completed qualification is audited via the matrix in
TASK-F16-02, wired into the qualify-commit dispatch here.

**Exclusions:** No HTTP handlers, no HTTP routes, no maintenance enter/clear
(added by TASK-F16-06 into this same lifecycle.go file), no projection. These
are future tasks.

**Commit message:**
```
feat(f16): TASK-F16-05 — synthetic observer and qualification engine, commit dispatch integrated into lifecycle service

Implements REQ-F16-02, REQ-F16-07, REQ-F16-08, REQ-F16-09 (lifecycle viability
recheck, observer, qualification, fenced audited commit).

Creates sovrunn.synthetic-iaas-observer/v1 as target-bound, clock-driven,
external-effect-free component; four named facts (compute.vm, storage.block,
storage.object, network.private) with Supported/Unsupported/Unknown truth; closed
profile evaluation (any Unsupported → Rejected, else any Unknown → Indeterminate,
else Qualified); missing fixture and logical timeout yield four Unknown facts at
now/+60s with no network wait; duplicate/missing/malformed facts fault
qualification.

Adds the qualify-commit dispatch method to lifecycle.go (sole-committer mutex):
distinct action-idempotency and internal Qualifying reservations; internal
Qualifying/fence capture under lock, unlocked observation, then a separate
relocked commit dispatch. A failing under-mutex backing viability recheck before
the Qualifying reservation opens is this task's sole cleanup path: detach the
linked action reservation, no completion/observation/AuditEvent/publication,
unlock, then exactly-once signal. Exact ADH-2026-060 ordered
commit predicate (active Maintenance marker → CONFLICT/409, else changed epoch →
STALE_RESOURCE_VERSION/412, else changed generation/viability →
STALE_RESOURCE_VERSION/412); obtains build-only TASK-F16-02 event then performs exactly one append before publication;
idempotency completion via TASK-F16-03 table only after publication; atomic
record replacement on success; observer fault/stale-fence/panic/shutdown/
cancellation abort the Qualifying reservation, remove the linked InFlight
idempotency reservation (not left unchanged), create no completion, and wake
waiters only after unlock, without target mutation or AuditEvent.

Owns and unit-proves Retire-wins, and performs its production wiring: updates
TASK-F16-04's retire-commit method in lifecycle.go so it now calls into the
reservation-abort mechanism added here while a Qualifying reservation is in
flight for the same target (retire remains valid from every Active
combination, including while Maintenance is active); that reservation aborts
with no qualification result, completion, or qualification AuditEvent, and the
qualifying caller's outcome is 409 CONFLICT + VS0_TARGET_RETIRED (retirement
itself audited via TASK-F16-04). TASK-F16-11 later re-exercises this
end-to-end as an integration re-test.

Proves Risk 2 (externalCallCount=0; no network wait), Risk 6 (stale-fence
non-publication).

No HTTP, maintenance enter/clear, or projection.

Relates: FEATURE-0016, DEC-0036, DEC-0042, ADH-2026-058, ADH-2026-060.
```

---

### TASK-F16-06 — Injected-clock expiry scheduler and maintenance trigger intake (maintenance commit methods integrated into lifecycle service)

**Prerequisites:** TASK-F16-01, TASK-F16-02, TASK-F16-03, TASK-F16-04,
TASK-F16-05 complete.

**Requirement authority:** REQ-F16-09 (expiry clock-driven, maintenance
target-bound trigger, shutdown ordering).

**Design authority:** DD-06 (injected clock drives expiry retry and freshness),
DD-07 (maintenance enters/clears only through target-bound fixture trigger),
design §3.1 (scheduler.go), design §6.5 (operational shutdown ordering).

**Architecture decisions:** DEC-0057 (maintenance epoch fence), ADH-2026-058
clause 8 (Retired/Maintenance outrank expiry, maintenance trigger never HTTP route
or event bus, shutdown ordering), ADH-2026-061 (unconditional link clearing and
target-specific abort scope).

**Applicable risks:** Risk 2 (no external effect: maintenance trigger is
deterministic target-bound fixture intake, never an external call; shutdown clears
F0016-owned in-memory state).

**Writable paths:**
- `internal/executiontarget/scheduler.go`
- `internal/executiontarget/clock.go` (injected Clock interface + real/fake
  implementations)
- `internal/executiontarget/lifecycle.go` (adding the maintenance-enter and
  maintenance-clear commit methods, and the shutdown-ordering hook that aborts
  every F0016 InFlight idempotency reservation for create, qualify, and retire,
  plus any in-flight qualification reservation; no change to
  the create/retire/qualify-commit methods owned by TASK-F16-04/05)
- `internal/executiontarget/scheduler_test.go`
- `internal/executiontarget/clock_test.go`

**Tests:**
Unit tests for injected Clock interface with real monotonic implementation and
deterministic fake; expiry worker runs once per injected-clock second and updates
only fact freshness/marker inputs when facts expire. TASK-F16-07 is the sole
owner of effectiveAvailability computation and projection assertions. No expiry
AuditEvent until append succeeds, retrying once per
injected-clock second; Retired/Maintenance outrank expiry.

Maintenance enter/clear is delivered exclusively through deterministic
target-bound synthetic-observer fixture trigger carrying target UID, expected
InfrastructureStack generation, expected maintenance epoch, operation, system
identity, correlation; the lifecycle service's maintenance-enter/clear commit
methods (added here) validate the fence (matching current epoch/generation)
before applying; stale-fenced triggers produce no state change, AuditEvent, or
record publication.

Maintenance enter attempts the required AuditEvent append before it publishes
any Maintenance transition. A failed Maintenance-entry append returns
`INTERNAL_ERROR`/500 with target, ETag, marker, FactSet/Result links, indexes,
and tuple unchanged; it does not abort an in-flight qualification or its linked
idempotency reservation. Tests explicitly prove F16-100 and F16-101. On a
successful append, Maintenance enter sets the current-Maintenance marker (target UID,
maintenance epoch, active flag) active at the new epoch, increments maintenance
epoch and resourceVersion; TASK-F16-07 computes the resulting effectiveAvailability.
Per ADH-2026-061, every successful Maintenance entry unconditionally clears
both the target's current NormalizedTargetFactSet and TargetQualificationResult
links and persists Active/Unqualified — regardless of whether a qualification
is in flight for that target; it never preserves the last completed
conclusion. In addition, if a qualification is in flight for the *specific
target* receiving the trigger at that moment, maintenance enter aborts only
that one Qualifying reservation for that target (via the qualify-commit
dispatch added in TASK-F16-05; no other target's reservation is touched) and
its linked idempotency reservation (via the TASK-F16-03 table's package-private
exactly-once helper), waking only that reservation's waiters, without giving it a
completion. If no qualification is in flight for that target, no idempotency
reservation is cancelled, and the Maintenance entry is still audited. Entry is
audited via the TASK-F16-02 assembly before publication in every case. When a
Qualifying reservation aborts because Maintenance entry won the race, the
qualifying caller's outcome is 409 CONFLICT + `VS0_TARGET_MAINTENANCE` with no
qualification result, completion, or qualification AuditEvent — this task owns
and unit-tests that qualifying-caller-side outcome as a unit-level property;
TASK-F16-11 later re-exercises it end-to-end over the live HTTP actions as an
integration re-test citing this task as the unit-level owner.

Maintenance clear removes matching current-Maintenance marker, clears current
FactSet/Result links, increments maintenance epoch and resourceVersion, and
audits successful clear before
publication. Stale-fenced maintenance enter/clear produce no state change,
AuditEvent, or publication.

This task provides the lifecycle shutdown hook: under the lifecycle mutex it
first atomically closes new lifecycle/idempotency reservation admission, then
detaches every existing F0016 InFlight reservation for create, qualify, and
retire (and any in-flight qualification reservation); it unlocks and wakes all
captured waiters exactly once. No new InFlight reservation can appear after the
abort sweep begins. ADH-2026-061's target-specific limitation applies only to
successful Maintenance entry, not shutdown. TASK-F16-11 owns production
composition and invokes this hook before HTTP shutdown; a trigger or handler
reaching admission after the stop boundary creates no reservation, mutation, or
AuditEvent.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -race ./internal/executiontarget/...
```

**Acceptance criteria:**
- Injected Clock interface exists with real monotonic and deterministic fake
  implementations.
- Expiry worker runs once per injected-clock second, updates freshness inputs,
  and retries expiry audit once per second on append
  failure.
- Retired/Maintenance outrank expiry.
- Maintenance enter/clear is delivered only through deterministic target-bound
  fixture trigger, dispatched to the lifecycle service's maintenance commit
  methods added in this task; stale-fenced triggers produce no state change,
  AuditEvent, or publication.
- Maintenance enter always sets active current-Maintenance marker and
  increments maintenance epoch/resourceVersion, and implements the exact
  ADH-2026-061 rule: every successful Maintenance entry unconditionally clears
  both current NormalizedTargetFactSet and TargetQualificationResult links and
  persists Active/Unqualified, whether or not a qualification is in flight for
  that target (never preserving the last completed conclusion). If a
  qualification is in flight for that specific target, entry additionally
  aborts only that one Qualifying reservation for that target (not any other
  target's reservation) and its linked idempotency reservation without a
  completion, waking only that reservation's waiters. If no qualification is
  in flight for that target, no idempotency reservation is cancelled. Entry is
  audited before publication in every case.
- Required audit append failure for Maintenance entry or clear proves
  F16-100/F16-101 exactly: no target state, ETag, marker, FactSet/Result link,
  index, tuple, qualification reservation, or idempotency reservation changes;
  in particular, failed entry does not abort in-flight qualification work.
- This task owns and unit-proves the Maintenance-wins qualifying-caller-side
  outcome: when Maintenance entry aborts an in-flight Qualifying reservation,
  the qualifying caller receives 409 CONFLICT + `VS0_TARGET_MAINTENANCE` with
  no qualification result, completion, or qualification AuditEvent;
  TASK-F16-11 re-exercises this end-to-end as an integration re-test citing
  this task as the unit-level owner.
- Maintenance clear removes matching marker, clears current FactSet/Result links,
  increments maintenance epoch/resourceVersion, and audits clear before
  publication. Its committed state/freshness inputs are projected as
  `effectiveAvailability: Unavailable` only by TASK-F16-07.
- Deterministic concurrent tests prove shutdown closes new reservation admission
  under the lifecycle mutex before detaching every create/qualify/retire
  InFlight idempotency reservation, then wakes all corresponding waiters after
  unlock; no InFlight reservation can appear after the abort sweep begins.
  TASK-F16-11 proves this hook precedes HTTP shutdown in the composed server;
  post-stop triggers/handlers produce no reservation, mutation, or AuditEvent.
- No external call or network I/O; `externalCallCount=0` remains proven (Risk 2).

**Security impact:** None (internal clock/trigger only).

**Observability impact:** Expiry audit retry ensures eventual audit without 500 to
caller; maintenance enter/clear audited before publication.

**Exclusions:** No HTTP handlers, no HTTP routes, no projection. These are future
tasks. This task completes all lifecycle-service commit methods
(create/retire/qualify/maintenance-enter/maintenance-clear) in `lifecycle.go`.

**Commit message:**
```
feat(f16): TASK-F16-06 — injected-clock expiry scheduler and maintenance trigger intake, maintenance commit methods integrated into lifecycle service

Implements REQ-F16-09 (expiry clock-driven, maintenance target-bound trigger,
shutdown ordering).

Creates injected Clock interface with real monotonic and deterministic fake
implementations; expiry worker runs once per injected-clock second and updates
fact freshness inputs when facts expire; TASK-F16-07 alone computes effectiveAvailability,
retries expiry audit once per second on append failure; Retired/Maintenance
outrank expiry.

Adds maintenance-enter and maintenance-clear commit methods to lifecycle.go
(sole-committer mutex): delivered exclusively through deterministic target-bound
synthetic-observer fixture trigger (never HTTP route or event bus); stale-fenced
triggers produce no state change, AuditEvent, or publication. Maintenance enter
always sets active current-Maintenance marker and increments
epoch/resourceVersion, and implements the exact ADH-2026-061 rule: every
successful Maintenance entry unconditionally clears both current
NormalizedTargetFactSet and TargetQualificationResult links and persists
Active/Unqualified, whether or not a qualification is in flight (never
preserving the last completed conclusion); if a qualification is in flight for
that specific target, entry additionally aborts only that one Qualifying
reservation for that target and its linked idempotency reservation without
completion, waking only that reservation's waiters; if no qualification is in
flight, no idempotency reservation is cancelled. Owns and unit-proves the
Maintenance-wins qualifying-caller-side outcome: 409 CONFLICT +
VS0_TARGET_MAINTENANCE with no qualification result, completion, or
qualification AuditEvent (TASK-F16-11 re-exercises this end-to-end). Entry is
audited before publication in every case by obtaining the TASK-F16-02 event then
performing exactly one append attempt. Maintenance clear removes marker, clears
links, increments epoch/resourceVersion, and audits clear before publication.
TASK-F16-07 alone projects the resulting state/freshness inputs as
`effectiveAvailability: Unavailable`.

Adds the lifecycle shutdown hook: atomically close new lifecycle/idempotency
reservation admission under the lifecycle mutex → detach every existing
create/qualify/retire InFlight reservation → unlock and wake captured waiters.
Deterministic concurrency tests prove no InFlight reservation appears after the
abort sweep begins. TASK-F16-11 owns invoking this hook before HTTP shutdown;
post-stop triggers/handlers create no reservation, mutation, or AuditEvent.

Proves Risk 2 (externalCallCount=0; no network I/O).

Completes all lifecycle.go commit methods. No HTTP or projection.

Relates: FEATURE-0016, REQ-F16-09, DEC-0057, ADH-2026-058, ADH-2026-061.
```

---

### TASK-F16-07 — Safe projection and effectiveAvailability computation

**Prerequisites:** TASK-F16-01, TASK-F16-04, TASK-F16-05, TASK-F16-06 complete.

**Requirement authority:** REQ-F16-04 (closed safe projection), REQ-F16-08
(response-only effectiveAvailability truth table, no persisted status.availability).

**Design authority:** Design §3.1 (projection.go), design §4.1 (closed safe
ExecutionTarget projection), design §4.6 (effectiveAvailability ordered truth
table).

**Architecture decisions:** ADH-2026-058 clauses 4, 7 (closed safe projection with
no public conditions, response-only effectiveAvailability).

**Applicable risks:** Risk 4 (safe-denial non-disclosure: design §5.1 step 6 +
§6.2; safe projection omits facts/results/observer identity/handles/internal
record links).

**Writable paths:**
- `internal/executiontarget/projection.go`
- `internal/executiontarget/projection_test.go`

**Tests:**
Unit tests for closed safe ExecutionTarget JSON projection containing exactly:
`metadata.uid`, `metadata.name`, `metadata.generation`, `metadata.resourceVersion`,
the three immutable spec fields (`cloudProviderParticipationRef`,
`infrastructureStackRef`, `targetClass`), `status.lifecycle`, `status.qualification`,
and response-only `effectiveAvailability`. No public `conditions` member; facts,
results, observer identity, handles, internal record links omitted.

The lifecycle service provides the projection a coherent immutable snapshot of
target lifecycle/qualification, current marker, FactSet/Result links, fact
freshness, and backing-viability inputs captured under the DD-02 mutex. This
task never reads mutable store maps, markers, facts, or viability separately.

`effectiveAvailability` ordered truth table:

1. Retired → `Unavailable`.
2. else Active with current Maintenance marker active → `Maintenance`.
3. else only Active + Qualified + fresh current facts + viable backing →
   `Available`.
4. every other Active combination → `Unavailable`.

A viability-only change (participation effective-active state or stack phase)
changes effectiveAvailability projection with no target mutation or ETag change.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -race ./internal/executiontarget/...
```

**Acceptance criteria:**
- Closed safe projection proven by unit tests: exactly the approved fields per
  design §4.1, no public conditions, no facts/results/observer/handles/links.
- `effectiveAvailability` ordered truth table implemented and proven by unit
  tests against a lifecycle-service coherent snapshot captured under the DD-02
  mutex; no unsynchronized component reads are permitted.
- A viability-only change projects changed effectiveAvailability with no target
  mutation or ETag change proven by unit tests.
- Safe-denial non-disclosure at this task's own scope: the projection function
  never includes facts/results/observer identity/handles/internal record links
  in its output, for any target passed to it. The full Risk 4 guarantee that
  an inaccessible target/backing is never disclosed at all (safe 404 before
  this projection function is ever called) is proven by TASK-F16-09/10, which
  integrate safe reference access ahead of this projection.

**Security impact:** Establishes this task's own half of the safe-denial
non-disclosure barrier (Risk 4): never exposes facts/results/observer/handles/
links in the projected output. TASK-F16-09/10 integrate safe reference access
before this projection is reached, completing the full Risk 4 guarantee.

**Observability impact:** None (projection only; no new audit or log).

**Exclusions:** No HTTP handlers, no HTTP routes. These are future tasks. This
task provides the projection function consumed by HTTP handler tasks.

**Commit message:**
```
feat(f16): TASK-F16-07 — safe projection and effectiveAvailability computation

Implements REQ-F16-04, REQ-F16-08 (closed safe projection, response-only
effectiveAvailability).

Creates closed safe ExecutionTarget JSON projection containing exactly:
metadata.uid/name/generation/resourceVersion, the three immutable spec fields,
status.lifecycle, status.qualification, and response-only effectiveAvailability;
no public conditions; facts/results/observer identity/handles/internal record
links omitted.

effectiveAvailability ordered truth table: Retired → Unavailable; else Active with
current Maintenance marker active → Maintenance; else only Active + Qualified +
fresh current facts + viable backing → Available; every other Active combination →
Unavailable. A viability-only change projects changed effectiveAvailability with
no target mutation or ETag change.

Proves Risk 4 (safe-denial non-disclosure; never exposes
facts/results/observer/handles/links).

No HTTP or routes.

Relates: FEATURE-0016, REQ-F16-04, REQ-F16-08, ADH-2026-058.
```

---

### TASK-F16-08 — Pre-ServeMux transport guard (guard only; no route registration, no placeholder handlers)

**Prerequisites:** TASK-F16-01 complete.

**Requirement authority:** REQ-F16-03 (four F0016 path patterns registering five
method/path registrations total, pre-ServeMux transport guard).

**Design authority:** DD-08 (pre-ServeMux guard wraps the five registrations),
design §3.1 (routes.go + guard helper), design §4.5 (HTTP surface: four path
patterns, five method/path registrations — the collection path registers both
`POST` and `GET`).

**Architecture decisions:** ADH-2026-058 clause 3 (exactly four path patterns
registering five method/path registrations total, behind pre-ServeMux guard
supplying transport-only 405/404 with no Problem body, authentication, audit,
idempotency, or lifecycle effect).

**Applicable risks:** Risk 5 (closed four-path/five-registration surface with
side-effect-free guard: design DD-08 + VS0-CF-F16-44,113..118; unmatched-method/
trailing-slash transport outcomes with no authentication, authorization, audit,
idempotency, or lifecycle effect).

**Writable paths:**
- `internal/server/executiontarget_guard.go`
- `internal/server/executiontarget_guard_test.go`

**Tests:**
Unit tests for the pre-ServeMux method/path guard supplying, for each of the four
exact F0016 path patterns, a bounded allowed-method-set invariant rather than an
open-ended enumeration, with an exact literal `Allow` header serialization (not
an unordered set) per path: collection path
(`/apis/execution.sovrunn.io/v1alpha1/execution-targets`): allowed methods
`GET, POST` with literal header value `Allow: GET, POST` (this one path pattern
registers both method/path registrations); item-GET path
(`.../execution-targets/{uid}`): allowed method `GET` with literal header value
`Allow: GET`; qualify-action path (`.../execution-targets/{uid}/actions/qualify`):
allowed method `POST` with literal header value `Allow: POST`; retire-action
path (`.../execution-targets/{uid}/actions/retire`): allowed method `POST` with
literal header value `Allow: POST` — four path patterns, five method/path
registrations total. The guard returns transport-only 405 for any method value
not in that path's allowed-method set (a general "not-in-allowed-set → 405"
invariant proven for representative non-allowed methods including but not
limited to HEAD, PUT, DELETE, PATCH, OPTIONS, TRACE, and CONNECT — not an
exhaustive enumerated list of every conceivable method name), each with the
exact literal `Allow` header value shown above for that path and
`X-Content-Type-Options: nosniff`; and a transport-only 404 without redirect
for every trailing-slash variant of each of the four path patterns, each with
no Problem body, authentication, authorization, audit, idempotency, or
lifecycle effect. The guard itself is a standalone `http.Handler` decorator in
this task; it is not yet wired to the five real ServeMux registrations, which
do not exist until TASK-F16-09/10/11 register them together with their real
handlers. Tests also prove compatibility: each valid F0016 method/path pair
delegates to the next handler exactly once without guard-authored response
mutation, and every non-F0016 method/path delegates unchanged exactly once.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -race ./internal/server/...
```

**Acceptance criteria:**
- Pre-ServeMux guard decorator exists, ready to wrap exactly the five Go 1.22
  `http.ServeMux` registrations across four path patterns per design §4.5, once
  those registrations exist.
- For each of the four exact F0016 path patterns, the guard implements the
  bounded allowed-method-set invariant with the exact literal `Allow` header
  serialization: collection path → `Allow: GET, POST`; item-GET path →
  `Allow: GET`; qualify-action path → `Allow: POST`; retire-action path →
  `Allow: POST`. Every method not in that path's allowed-method set returns
  transport-only 405 with that path's exact literal `Allow` header value and
  `X-Content-Type-Options: nosniff`, proven for representative non-allowed
  methods (HEAD, PUT, DELETE, PATCH, OPTIONS, TRACE, CONNECT, and others) as an
  instance of the general invariant, not as an exhaustive enumerated list.
- Trailing-slash variants return transport-only 404 without redirect or Problem
  body.
- All transport outcomes have no authentication, authorization, audit, idempotency,
  or lifecycle effect proven by unit tests.
- For each valid F0016 method/path pair, the decorator delegates to its next
  handler exactly once and adds no guard response; every non-F0016 route also
  passes through unchanged exactly once. This protects inherited routes while
  retaining the closed F0016 transport surface.
- Guard-only method/path and transport behavior proven here (Risk 5); the
  completed five-registration surface is proven only by TASK-F16-11/12.
- No 501 or any other placeholder status is produced by this task for any path;
  this task registers no route and defines no handler for the five real
  registrations.

**Security impact:** Establishes closed route surface (Risk 5) with side-effect-free
transport guard.

**Observability impact:** None (transport guard only; no audit or log).

**Exclusions:** No route registration (`internal/server/routes.go` is not
written by this task), no handler implementation, and no placeholder response of
any kind for the five real F0016 registrations. Each real registration is added
by the same task that implements that route's real handler (TASK-F16-09/10/11),
so a registered path always begins at authentication and never observably
returns a placeholder status.

**Commit message:**
```
feat(f16): TASK-F16-08 — pre-ServeMux transport guard (guard only)

Implements REQ-F16-03 (pre-ServeMux transport guard).

Creates pre-ServeMux method/path guard decorator supplying, for each of the
four exact F0016 path patterns (five method/path registrations total, since
the collection path registers both POST and GET), a bounded
allowed-method-set invariant with exact literal Allow header serialization:
collection path → Allow: GET, POST; item-GET path → Allow: GET; qualify-action
path → Allow: POST; retire-action path → Allow: POST; any method not in a
path's allowed set returns transport-only 405 with that path's exact literal
Allow header value and X-Content-Type-Options: nosniff, and transport-only 404
without redirect for every trailing-slash variant; no Problem body,
authentication, authorization, audit, idempotency, or lifecycle effect.

Proves Risk 5 (closed four-path/five-registration surface with side-effect-free
guard) via the general not-in-allowed-set invariant, not a narrowed HEAD/PUT-only
or open-ended enumerated-method proof. For valid F0016 method/path pairs it
delegates to the next handler exactly once without guard-authored response
mutation, and it passes every non-F0016 route through unchanged exactly once.

Registers no route and defines no placeholder handler; the five real
registrations are added by TASK-F16-09/10/11 together with their real handlers,
so no interim 501 or other placeholder state is ever created or committed.

Relates: FEATURE-0016, REQ-F16-03, ADH-2026-058, Risk 5.
```

---

### TASK-F16-09 — Collection create and LIST handlers with explicit pipeline and route registration

**Prerequisites:** TASK-F16-01, TASK-F16-04, TASK-F16-07, TASK-F16-08 complete.

**Requirement authority:** REQ-F16-01 (closed create-field boundary, derived
scope), REQ-F16-02 (safe reference access before graph validation), REQ-F16-04
(request grammar, ETag, ignored Accept), REQ-F16-05 (executable request
precedence), REQ-F16-06 (idempotency applies to create only), REQ-F16-09
(create and collection-denial AuditEvent behavior).

**Design authority:** DD-09 (explicit ordered request pipeline), design §3.1
(executiontarget_collection.go, routes.go), design §4.1 (create body closed set,
LIST ascending-UID no ETag), design §5.1 (request precedence ordered sequence,
stages matrix).

**Architecture decisions:** ADH-2026-058 clauses 1, 3, 4, 5 (closed create fields,
POST requires application/json + Idempotency-Key, LIST ignores idempotency, the
two collection registrations are wired behind the pre-ServeMux guard from
TASK-F16-08).

**Applicable risks:** All risks proven by earlier tasks (sole-committer, no
external effect, audit-before-publication, safe-denial, closed route surface)
are integrated through this task's explicit pipeline and route wiring.

**Writable paths:**
- `internal/api/executiontarget_collection.go`
- `internal/api/executiontarget_collection_test.go`
- `internal/server/routes.go` (registering exactly the two collection patterns —
  `POST` and `GET` on
  `/apis/execution.sovrunn.io/v1alpha1/execution-targets` — behind the
  unexposed F0016 mux builder; wired directly to this task's real handlers,
  never to a placeholder. TASK-F16-11 alone exposes the completed builder
  behind the TASK-F16-08 guard.)

**Tests:**
Unit and integration tests for collection POST and GET.

POST (create): closed transport guard → authentication → POST media/key grammar →
phase-one extraction of only `spec.cloudProviderParticipationRef.uid` and
`spec.infrastructureStackRef.uid` from closed envelope → derived-scope
authorization → safe backing access → strict classification → graph/reference and
uniqueness validation → reservation/replay → audit-before-publication commit
(via the TASK-F16-04 lifecycle service). Returns 201 with safe projection and ETag;
it appends the AuditEvent as a side effect, never in the response body. Completed replay persists Active/Unqualified at epoch 0 with
observedGeneration equal to referenced InfrastructureStack generation, projects
effectiveAvailability Unavailable.

Accepts exactly: `metadata.name`, `spec.cloudProviderParticipationRef.uid`,
`spec.infrastructureStackRef.uid`, immutable `spec.targetClass=synthetic-iaas`.
Rejects: `metadata.uid`, generation, resourceVersion, timestamps, `metadata.scopeRef`,
any status field, `adapterAuthorityRef`, SecretRef extension, unknown field, field
owned by another feature.

GET (LIST): closed transport guard → authentication → read authorization →
accessible-target enumeration and ascending-UID projection. Returns 200 with
accessible Active and Retired targets ascending `metadata.uid`, no ETag; ignores
any Idempotency-Key and creates, waits on, reads, replays no idempotency record.

Route-builder test: the two real collection registrations are added to the
unexposed F0016 mux builder with no placeholder. TASK-F16-11 alone proves
process-start exposure behind the one installed guard.

All failure paths per design §5.1: missing/invalid authentication → 401
AUTH_REQUIRED; missing Idempotency-Key or non-application/json media → 400
MALFORMED_REQUEST / 415 UNSUPPORTED_MEDIA_TYPE; without executiontarget.write →
403 AUTHORIZATION_DENIED; inaccessible backing → safe 404 RESOURCE_NOT_FOUND +
VS0_AUTHORIZATION_SAFE_DENIAL; ineffective participation → 409 CONFLICT +
VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE; non-Active stack → 409 CONFLICT +
VS0_EXECUTION_TARGET_STACK_UNAVAILABLE; cross-CloudProvider graph → 422
VALIDATION_FAILED + VS0_EXECUTION_TARGET_SCOPE_MISMATCH; duplicate name or live
tuple → 409 ALREADY_EXISTS; changed digest → 409 CONFLICT +
VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH.

This task owns collection audited-denial append behavior: tests cover F16-40
(create authorization denial) and the LIST-without-read member of F16-112. The
collection handler builds the denial event, performs exactly one inherited
append attempt, and returns `500 INTERNAL_ERROR` with no durable event or
publication if it fails; it never discloses the suppressed 403.

This task also owns collection HTTP contract cases F16-45, F16-48, F16-105,
and the POST/GET collection members of F16-106: success ignores `Accept` and
emits `application/json` plus `X-Content-Type-Options: nosniff`; every Problem
ignores `Accept` and emits `application/problem+json` plus nosniff; create
rejects `If-Match`; and any query is rejected before safe resolution. Completed
create replay returns only the stored successful body/media plus a fresh
correlation value—never raw Idempotency-Key, stored correlation, credentials,
or request headers.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -race ./internal/api/...
go test -race ./internal/server/...
go test -race ./internal/executiontarget/...
```

**Acceptance criteria:**
- Collection POST implements exact ordered pipeline per design §5.1; returns 201
  with safe projection and ETag, while AuditEvent is appended only as a side
  effect; completed replay persists
  Active/Unqualified at epoch 0.
- Closed create-field boundary enforced: only approved client fields accepted;
  server-owned/status/scope/unknown/deferred rejected.
- Derived CloudProvider scope solely from resolved participation reference.
- Safe reference access before graph validation proven.
- Collection GET implements ascending-UID LIST with no ETag, no idempotency.
- The two collection patterns are added to the unexposed F0016 mux builder and
  wired to real handlers with no placeholder. TASK-F16-11 exposes the completed
  five-registration mux behind the TASK-F16-08 guard.
- All failure paths return exact Problem status/violation per design §5.1.
- Task-local tests prove F16-40 and the collection LIST member of F16-112:
  successful denial append retains the suppressed response externally, while
  append failure returns only `500 INTERNAL_ERROR` with no durable AuditEvent
  or publication.
- Handler integration tests exercise TASK-F16-04 create owner cancellation and
  panic cleanup: no completion/publication, detach under lock, then one
  post-unlock waiter signal.
- Task-local tests prove F16-45, F16-48, F16-105, and collection F16-106 with
  exact fixed media, nosniff, query precedence, and replay header/correlation
  filtering.
- Integration tests prove end-to-end create/LIST with idempotency replay and
  authorization/safe-denial.

**Security impact:** Enforces authentication, authorization, safe-denial, closed
field boundary, and idempotency waiter recheck.

**Observability impact:** Successful create and authorization/safe denials audited.

**Exclusions:** No item GET, no qualify, no retire, no registration of those
three patterns. These are TASK-F16-10/11.

**Commit message:**
```
feat(f16): TASK-F16-09 — collection create and LIST handlers with route registration

Implements REQ-F16-01, REQ-F16-02, REQ-F16-04, REQ-F16-05, REQ-F16-06,
REQ-F16-09 (create/LIST pipeline, idempotency, audit behavior, route
registration).

Creates collection POST (create) with exact ordered pipeline: closed transport
guard → authentication → POST media/key grammar → phase-one extraction → derived-scope
authorization → safe backing access → strict classification → graph/reference and
uniqueness validation → reservation/replay → audit-before-publication commit;
returns 201 with safe projection and ETag; AuditEvent is appended only as a side
effect, never in the response body; completed replay persists
Active/Unqualified at epoch 0 with observedGeneration equal to referenced
InfrastructureStack generation, projects effectiveAvailability Unavailable.

Accepts exactly: metadata.name, spec.cloudProviderParticipationRef.uid,
spec.infrastructureStackRef.uid, immutable spec.targetClass=synthetic-iaas.
Rejects: metadata.uid, generation, resourceVersion, timestamps, metadata.scopeRef,
any status field, adapterAuthorityRef, SecretRef extension, unknown field, field
owned by another feature.

Creates collection GET (LIST) with ascending-UID projection, no ETag, no
idempotency; ignores any Idempotency-Key.

Adds the two collection registrations to the unexposed F0016 mux builder,
wired directly to real handlers with no placeholder. TASK-F16-11 exposes the
completed five-registration builder once behind the TASK-F16-08 guard.

All failure paths return exact Problem status/violation per design §5.1.

Integration tests prove end-to-end create/LIST with idempotency replay,
authorization, safe-denial.

No item GET, qualify, or retire (TASK-F16-10/11).

Relates: FEATURE-0016, REQ-F16-01..06, ADH-2026-058.
```

---

### TASK-F16-10 — Item GET handler with safe projection integration and route registration

**Prerequisites:** TASK-F16-01, TASK-F16-07, TASK-F16-08, TASK-F16-09 complete.

**Requirement authority:** REQ-F16-04 (item GET returns ETag safe projection, no
idempotency), REQ-F16-05 (executable request precedence), REQ-F16-09 (item safe-
denial AuditEvent behavior).

**Design authority:** DD-09 (explicit ordered request pipeline), design §3.1
(executiontarget_item.go, routes.go), design §5.1 (request precedence ordered
sequence, stages matrix for item GET).

**Architecture decisions:** ADH-2026-058 clauses 3, 4 (ETag on item GET, ignored
Accept, the item registration is wired behind the pre-ServeMux guard from
TASK-F16-08).

**Applicable risks:** Risk 4 (safe-denial non-disclosure; inaccessible direct GET
returns safe 404 with no disclosure).

**Writable paths:**
- `internal/api/executiontarget_item.go`
- `internal/api/executiontarget_item_test.go`
- `internal/server/routes.go` (registering exactly the item-GET pattern —
  `GET` on `/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}` —
  on the unexposed F0016 mux builder, wired directly to this task's real
  handler, never to a placeholder; this task's edit is additive to the two
  collection registrations from TASK-F16-09. TASK-F16-11 alone exposes the
  completed builder behind the TASK-F16-08 guard.)

**Tests:**
Unit and integration tests for item GET.

GET: closed transport guard → authentication → phase-one path-UID safe resolution
→ derived-scope read authorization → safe access → safe projection. Returns 200
with ETag safe projection without facts, result, observer, or handle; ignores any
Idempotency-Key and creates, waits on, reads, replays no idempotency record.

Route-builder test: the real item registration is added to the unexposed F0016
mux builder with no placeholder. TASK-F16-11 alone proves process-start
exposure behind the one installed guard.

All failure paths: missing/invalid authentication → 401 AUTH_REQUIRED;
inaccessible direct item GET → safe 404 RESOURCE_NOT_FOUND +
VS0_AUTHORIZATION_SAFE_DENIAL; without executiontarget.read after UID-only safe
resolution derives target scope → 403 AUTHORIZATION_DENIED; malformed UID segment
→ 400 MALFORMED_REQUEST before safe resolution.

This task owns direct-item audited-denial behavior: it covers F16-102 and the
direct-GET member of F16-112. The item handler builds the denial event and
performs exactly one inherited append attempt; failed append returns only
`500 INTERNAL_ERROR`, with no durable AuditEvent or disclosure of the safe 404.

This task owns item-GET HTTP contract cases F16-45, F16-48, F16-107, F16-108,
and the item member of F16-106: it ignores `Accept`, emits fixed
`application/json` on success and `application/problem+json` on Problem with
nosniff, rejects query before safe resolution, rejects malformed canonical UID
before lookup, and preserves transport-only trailing-slash 404 behavior.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -race ./internal/api/...
go test -race ./internal/server/...
```

**Acceptance criteria:**
- Item GET implements exact ordered pipeline per design §5.1; returns 200 with
  ETag safe projection, no idempotency.
- Inaccessible direct GET returns safe 404 RESOURCE_NOT_FOUND +
  VS0_AUTHORIZATION_SAFE_DENIAL with no existence disclosure.
- UID-only safe resolution derives target scope before executiontarget.read check.
- Malformed UID segment returns 400 MALFORMED_REQUEST before safe resolution.
- Task-local tests prove F16-102 and the direct-GET F16-112 member: denial
  append is attempted exactly once; failure returns only `500 INTERNAL_ERROR`
  with no durable AuditEvent, publication, or safe-404 disclosure.
- Task-local tests prove F16-45, F16-48, F16-107, F16-108, and the item F16-106
  member with exact fixed media, nosniff, query precedence, malformed-UID
  ordering, and transport-only trailing-slash behavior.
- The item-GET registration is added to the unexposed F0016 mux builder and
  wired to its real handler with no placeholder. TASK-F16-11 exposes the
  completed five-registration builder once behind the TASK-F16-08 guard.
- Integration tests prove end-to-end item GET with authorization and safe-denial.

**Security impact:** Enforces safe-denial non-disclosure (Risk 4) for inaccessible
direct GET.

**Observability impact:** Inaccessible direct GET safe-denial audited (with
append-failure returning INTERNAL_ERROR/500 with no 404 disclosure).

**Exclusions:** No qualify, no retire, no registration of those two patterns.
These are TASK-F16-11.

**Commit message:**
```
feat(f16): TASK-F16-10 — item GET handler with route registration

Implements REQ-F16-04, REQ-F16-05, REQ-F16-09 (item GET pipeline, audited safe
denial, route registration).

Creates item GET with exact ordered pipeline: closed transport guard →
authentication → phase-one path-UID safe resolution → derived-scope read
authorization → safe access → safe projection; returns 200 with ETag safe
projection without facts, result, observer, or handle; ignores any Idempotency-Key
and creates, waits on, reads, replays no idempotency record.

Inaccessible direct GET returns safe 404 RESOURCE_NOT_FOUND +
VS0_AUTHORIZATION_SAFE_DENIAL with no existence disclosure. UID-only safe
resolution derives target scope before executiontarget.read check. Malformed UID
segment returns 400 MALFORMED_REQUEST before safe resolution.

Adds the item-GET ServeMux pattern to the unexposed F0016 mux builder, wired
directly to this real handler; no placeholder 501 or other interim state is
ever created or committed. TASK-F16-11 alone exposes the completed builder
behind the TASK-F16-08 guard.

Integration tests prove end-to-end item GET with authorization, safe-denial.

No qualify or retire (TASK-F16-11).

Relates: FEATURE-0016, REQ-F16-04, REQ-F16-05, ADH-2026-058, Risk 4.
```

---

### TASK-F16-11 — Qualify and retire action handlers with fenced commit and route registration

**Prerequisites:** TASK-F16-01, TASK-F16-04, TASK-F16-05, TASK-F16-06,
TASK-F16-07, TASK-F16-08, TASK-F16-09, TASK-F16-10 complete.

**Requirement authority:** REQ-F16-02 (qualify safe backing admission before
reservation), REQ-F16-04 (action grammar: If-Match + zero-byte body), REQ-F16-05
(executable request precedence), REQ-F16-06 (idempotency applies to
qualify/retire), REQ-F16-08 (synchronous fenced qualify, atomic commit), REQ-F16-09
(Retire wins, Maintenance entry wins, audit matrix, shutdown ordering).

**Design authority:** DD-04 (synchronous qualification with captured-fence
reservation and exact ADH-2026-060 ordered commit predicate), DD-09 (explicit
ordered request pipeline), design §3.1 (executiontarget_actions.go, routes.go),
design §5.1 (request precedence ordered sequence, stages matrix for actions).

**Architecture decisions:** ADH-2026-058 clauses 3, 4, 5, 7, 8 (action If-Match +
zero-byte body, the two action registrations wired behind the pre-ServeMux
guard from TASK-F16-08, fenced commit, Retire/Maintenance win over in-flight
qualification, audit matrix), ADH-2026-060 (Maintenance-race outcome ordering),
ADH-2026-061 (unconditional Maintenance link clearing and target-specific abort).

**Applicable risks:** All risks proven by earlier tasks; this task integrates them
through action pipelines, fenced commit, and route wiring.

**Writable paths:**
- `internal/api/executiontarget_actions.go`
- `internal/api/executiontarget_actions_test.go`
- `internal/server/routes.go` (registering exactly the two action patterns —
  `POST .../actions/qualify` and `POST .../actions/retire` — behind the
  unexposed F0016 mux builder, wired directly to this task's real handlers,
  never to a placeholder; this task's edit is additive to the three
  registrations from TASK-F16-09/10, completing all five registrations before
  this task exposes the builder once behind the TASK-F16-08 guard)
- `internal/server/server.go` (the sole production composition point: construct
  the shared ExecutionTargetLifecycleService and scheduler, inject the shared
  service into every F0016 handler, install the TASK-F16-08 guard once around
  the completed F0016 mux, start expiry and Maintenance-trigger intake, and
  invoke the lifecycle shutdown hook before HTTP shutdown)
- `internal/server/server_test.go` (composition and ordered-shutdown integration)
- `cmd/sovrunn-api/main.go` (pass the composed server configuration at process start)

**Tests:**
Unit and integration tests for qualify and retire actions.

Production-composition integration test: one shared lifecycle service and one
scheduler are constructed once, injected into collection/item/action handlers,
and wrapped once by the pre-ServeMux guard after all five registrations exist.
On shutdown it proves: stop expiry and Maintenance-trigger acceptance → abort
new lifecycle/idempotency reservation admission atomically under the lifecycle
mutex → detach every existing F0016 InFlight idempotency reservation for create,
qualify, and retire, plus any in-flight qualification reservation → wake all
corresponding waiters after unlock → invoke HTTP shutdown. A deterministic
concurrent test proves no InFlight reservation can appear after the abort sweep
begins; no trigger or handler reaching admission after the stop boundary can
reserve, mutate state, or append an AuditEvent.

Qualify: closed transport guard → authentication → action media/key/If-Match/zero-byte
grammar → phase-one path-UID safe resolution → derived-scope qualify authorization
→ safe backing access/current viability admission → replay/reservation → current If-Match → lifecycle state → fenced
observer/qualification commit and audit-before-publication (via the
TASK-F16-04/05 lifecycle service qualify-commit dispatch). Returns 200 with safe
projection and ETag; AuditEvent is appended only as a side effect, never in the
response body. Completed replay atomically replaces both record
refs, refreshes observedGeneration, advances resourceVersion.

Starts only from Active Unqualified/Qualified/Rejected/Indeterminate; creates
internal Qualifying reservation with captured fences (InfrastructureStack
generation, maintenance epoch, viability fingerprint); invokes observer; applies
exact ADH-2026-060 ordered commit predicate: if active current-Maintenance marker
exists → abort with `CONFLICT`/409 + `VS0_TARGET_MAINTENANCE`; else if captured
maintenance epoch differs from current → abort with `STALE_RESOURCE_VERSION`/412 +
`VS0_TARGET_EPOCH_STALE`; else if InfrastructureStack generation changed → abort
with `STALE_RESOURCE_VERSION`/412 + `VS0_TARGET_EPOCH_STALE`; else if viability
fingerprint changed → abort with `STALE_RESOURCE_VERSION`/412 +
`VS0_EXECUTION_TARGET_VIABILITY_STALE`.

Before replay/reservation and observation, tests prove REQ-F16-02 admission:
inaccessible backing → safe `404 RESOURCE_NOT_FOUND`; ineffective participation
→ `409 CONFLICT` + `VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE`; non-Active
InfrastructureStack → `409 CONFLICT` +
`VS0_EXECUTION_TARGET_STACK_UNAVAILABLE`; scope mismatch → `422
VALIDATION_FAILED` + `VS0_EXECUTION_TARGET_SCOPE_MISMATCH`. Inaccessible backing
is the exception: the safe `404 RESOURCE_NOT_FOUND` denial assembles and appends
exactly one redacted AuditEvent; append failure returns only `500 INTERNAL_ERROR`
without disclosing the 404. The ineffective-participation, inactive-stack, and
authorized scope-mismatch 409/422 admission denials are unaudited and make no
reservation, observer call, AuditEvent, or completion.

If the TASK-F16-05 under-mutex backing recheck rejects after replay/reservation
but before an internal Qualifying reservation opens, TASK-F16-05 alone detaches
and signals the action reservation; this handler maps the returned 409/422
outcome and must not issue a duplicate abort.

The handler owns only failures resolved before it invokes lifecycle qualification
commit dispatch: for a linked idempotency reservation, stale If-Match or
lifecycle denial in VS0-CF-F16-26..28, F16-38, F16-74, and F16-85 calls the
lifecycle service's bounded pre-commit abort. The service detaches under the
DD-02 mutex, creates no completion, unlocks, then closes the captured `done`
channel exactly once. Once commit dispatch begins, TASK-F16-05/06 lifecycle
paths exclusively own observer fault, fence, audit-append, panic, cancellation,
Maintenance/Retire-wins, and shutdown aborts; a late handler cleanup request is
a bounded no-op when the reservation is already finalized. Tests prove unchanged
target/ETag/AuditEvent state for the applicable losing path, one post-unlock
signal, and waiter re-entry; no InFlight reservation is left stuck.

Retire: closed transport guard → authentication → action
media/key/If-Match/zero-byte grammar → phase-one path-UID safe resolution →
derived-scope write authorization → safe access → replay/reservation → current
If-Match → lifecycle state → audit-before-publication retirement commit (via the
TASK-F16-04 lifecycle service retire-commit method, as updated/extended by
TASK-F16-05). Returns 200 with safe projection
(Retired/Unqualified/effectiveAvailability Unavailable), ETag, refs cleared,
AuditEvent, tuple release, completed replay; retirement allowed from every
Active combination including while Qualifying is in flight and including while
Maintenance is active.

Production-composition test: after TASK-F16-11 completes the builder with the
two real action registrations, it exposes the complete five-registration mux
once behind the guard. From process start every registered path begins at
authentication; no placeholder response is observable.

All action failure paths: missing Idempotency-Key or If-Match → 400
MALFORMED_REQUEST or 412 STALE_RESOURCE_VERSION (header-form before body
classification); action valid If-Match with non-zero body → 400 MALFORMED_REQUEST;
malformed/weak/wildcard/multiple/stale If-Match → 412 STALE_RESOURCE_VERSION;
without executiontarget.qualify or executiontarget.write → 403
AUTHORIZATION_DENIED; inaccessible direct action → safe 404 RESOURCE_NOT_FOUND +
VS0_AUTHORIZATION_SAFE_DENIAL.

This task owns action and replay-time audited-denial behavior: F16-103
(same-key replay authorization denial), F16-104 (same-key replay safe denial),
and the qualify/retire/inaccessible-action members of F16-112. Each handler
builds the event then performs exactly one inherited append attempt. On append
failure it returns only `500 INTERNAL_ERROR`, with no durable AuditEvent,
publication, idempotency completion, or disclosure of the suppressed 403/404.

This task owns action HTTP contract cases F16-48, F16-109, F16-110, and the
qualify/retire members of F16-106: every action Problem ignores `Accept` and
emits `application/problem+json` plus nosniff; query is rejected before safe
resolution; malformed canonical UID is rejected before lookup; and successful
same-key replay returns the stored successful `application/json` body with a
fresh correlation only, never a raw Idempotency-Key, stored correlation,
credential, or request header.

Route-outcome split for lifecycle-state denial — the following two outcomes
apply to different actions and must not be conflated: (a) a Retired target
denies **qualify** with 409 CONFLICT + VS0_TARGET_RETIRED (Retire itself is
never denied by target state other than an already-Retired target rejecting a
new-key retire); (b) a target with an active current-Maintenance marker denies
**qualify only** with the registered Maintenance outcome, 409 CONFLICT +
VS0_TARGET_MAINTENANCE (TASK-F16-06 is the unit-level owner of this outcome).
Retire is never denied by an active current-Maintenance marker: retire remains
valid for every Active target, including during Maintenance, as well as while
Qualifying is in flight. Different-key qualify while Qualifying → 409 CONFLICT
+ VS0_TARGET_QUALIFICATION_IN_PROGRESS; changed digest → 409 CONFLICT +
VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH.

Retire wins during in-flight qualification: Qualifying caller gets 409 CONFLICT +
VS0_TARGET_RETIRED; no qualification result, completion, or qualification
AuditEvent; winning retirement audited. This end-to-end HTTP re-test exercises
for real, over the live retire and qualify actions, the rule that TASK-F16-05
already owns and unit-tests (the Qualifying reservation's abort logic and the
qualifying-caller-side 409/VS0_TARGET_RETIRED outcome); this task is not the
first or sole place this rule is established. Maintenance entry wins during
in-flight qualification: Qualifying caller gets 409 CONFLICT +
VS0_TARGET_MAINTENANCE (the registered Maintenance outcome only — qualify is
the only action Maintenance denies); no qualification result, completion, or
qualification AuditEvent; winning maintenance entry audited. This end-to-end
HTTP re-test similarly exercises for real the rule that TASK-F16-06 already
owns and unit-tests (the single-target Qualifying-reservation abort and the
qualifying-caller-side 409/VS0_TARGET_MAINTENANCE outcome). Retire remains
valid for every Active target throughout, including during Maintenance.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -race ./internal/api/...
go test -race ./internal/server/...
go test -race ./internal/executiontarget/...
```

**Acceptance criteria:**
- Qualify implements exact ordered pipeline per design §5.1; starts only from
  Active; creates Qualifying reservation with captured fences; invokes observer;
  applies exact ADH-2026-060 ordered commit predicate; atomically replaces both
  record refs on success; once lifecycle commit dispatch begins, its
  observer-fault/stale-fence/panic/shutdown/cancellation paths exclusively abort
  the Qualifying reservation, remove the linked InFlight idempotency reservation,
  and wake waiters only after unlock, without target mutation.
- Qualify performs exact REQ-F16-02 safe backing/current-viability admission
  before replay/reservation and observation. Inaccessible backing produces an
  audited safe 404 exactly once (or 500 on append failure without 404
  disclosure); ineffective participation, inactive stack, and authorized scope
  mismatch are unaudited 409/422 denials with no reservation, observer call, or
  completion.
- For every pre-commit post-reservation stale-version or lifecycle denial,
  including F16-26..28, F16-38, F16-74, and F16-85, handler integration invokes
  bounded lifecycle abort: detach under the DD-02 mutex, no completion, unlock,
  then one exactly-once `done` close. TASK-F16-05/06 exclusively own aborts
  after commit dispatch begins; a late handler cleanup is a tested no-op for an
  already-finalized reservation. Tests prove target/ETag/AuditEvent invariance
  and waiter wake-up/re-entry.
- Retire implements exact ordered pipeline; allowed from every Active combination
  including while Qualifying is in flight and including while Maintenance is
  active; returns 200 Retired/Unqualified/effectiveAvailability Unavailable,
  clears refs, releases tuple.
- Route-outcome split is unambiguous: Maintenance denies **qualify only** with
  the registered Maintenance outcome (409 CONFLICT + VS0_TARGET_MAINTENANCE,
  TASK-F16-06 is the unit-level owner); retire remains valid for every Active
  target, including during Maintenance, as well as while Qualifying is in
  flight — retire is never denied by an active current-Maintenance marker.
- Retire wins over in-flight qualification with 409 CONFLICT + VS0_TARGET_RETIRED;
  Maintenance entry wins over in-flight qualification with 409 CONFLICT +
  VS0_TARGET_MAINTENANCE; both audited, no qualification AuditEvent. This task
  re-proves both rules end-to-end over the live HTTP actions; TASK-F16-05 is
  the unit-level owner of the Retire-wins/VS0_TARGET_RETIRED rule and
  TASK-F16-06 is the unit-level owner of the
  Maintenance-wins/VS0_TARGET_MAINTENANCE rule — this task is not their first
  or sole proof.
- The exact VS0-CF-F16-42, VS0-CF-F16-43, and VS0-CF-F16-89 registered meanings
  are preserved unchanged; this task's narrative is consistent with them.
- The two action patterns complete the unexposed F0016 mux builder. This task
  alone exposes all five registrations once behind the TASK-F16-08 guard;
  production starts only from that guarded completed mux, with no placeholder
  response observable.
- All action failure paths return exact Problem status/violation per design §5.1.
- Task-local tests prove F16-103, F16-104, and each action member of F16-112:
  one append attempt on an audited denial; failure returns only
  `500 INTERNAL_ERROR` without a durable AuditEvent or suppressed 403/404.
- Task-local tests prove F16-48, F16-109, F16-110, and both action F16-106
  members with fixed Problem media/nosniff, query and malformed-UID precedence,
  and replay header/correlation filtering.
- Handler integration tests exercise TASK-F16-04 retire owner cancellation and
  panic cleanup, and prove a failed retirement append does not abort a real
  in-flight qualification unless retirement actually commits.
- Integration tests prove end-to-end qualify/retire with fenced commit,
  idempotency replay, Retire/Maintenance wins, authorization, safe-denial.
- Production composition constructs one shared service/scheduler, injects it
  into all F0016 handlers, installs the guard exactly once, starts workers and
  trigger intake, and proves atomic admission closure plus reservation abort
  precedes HTTP shutdown, with no post-sweep InFlight reservation.
- Integration tests prove the installed guard delegates each valid F0016 request
  exactly once to its real handler and passes a non-F0016 route unchanged to its
  next handler exactly once.

**Security impact:** Enforces action authorization, safe-denial, idempotency
waiter recheck, and fenced commit preventing stale-work publication.

**Observability impact:** Successful qualification, retirement, and
Retire/Maintenance wins audited; failed qualification and stale-fence aborts
produce no AuditEvent.

**Exclusions:** None; this completes HTTP handler implementation and all five
route registrations.

**Commit message:**
```
feat(f16): TASK-F16-11 — qualify and retire action handlers with fenced commit and route registration

Implements REQ-F16-02, REQ-F16-04..09 (safe backing admission, action grammar,
pipeline, idempotency, fenced commit,
Retire/Maintenance wins, audit matrix, shutdown ordering, route registration).

Creates qualify action with exact ordered pipeline: closed transport guard →
authentication → action media/key/If-Match/zero-byte grammar → phase-one path-UID
safe resolution → derived-scope qualify authorization → safe backing access/current
viability admission →
replay/reservation → current If-Match → lifecycle state → fenced
observer/qualification commit and audit-before-publication; starts only from
Active; creates Qualifying reservation with captured fences; invokes observer;
applies exact ADH-2026-060 ordered commit predicate (active Maintenance marker →
CONFLICT/409, else changed epoch → STALE_RESOURCE_VERSION/412, else changed
generation/viability → STALE_RESOURCE_VERSION/412); atomically replaces both
record refs on success; lifecycle commit dispatch exclusively handles observer
fault/stale-fence/panic/shutdown/cancellation aborts, removes the linked InFlight
idempotency reservation (not left unchanged), creates no completion, and wakes
waiters only after unlock, without target mutation.

Creates retire action with exact ordered pipeline; allowed from every Active
combination including while Qualifying is in flight and including while
Maintenance is active; returns 200
Retired/Unqualified/effectiveAvailability Unavailable, clears refs, releases
tuple.

Retire wins over in-flight qualification with 409 CONFLICT + VS0_TARGET_RETIRED;
Maintenance entry wins with 409 CONFLICT + VS0_TARGET_MAINTENANCE; both audited,
no qualification AuditEvent. Both rules are re-proved here end-to-end over the
live HTTP actions; TASK-F16-05 unit-owns the Retire-wins rule and TASK-F16-06
unit-owns the Maintenance-wins rule.

Registers the two action ServeMux patterns in internal/server/routes.go behind
the TASK-F16-08 guard, wired directly to these real handlers; no placeholder
501 or other interim state is ever created or committed. Completes all five
F0016 ServeMux registrations.

All action failure paths return exact Problem status/violation per design §5.1.
This task also owns final production composition: one shared lifecycle service
and scheduler are constructed, injected into every F0016 handler, wrapped once
by the guard, started for expiry/trigger intake, and shut down through the
TASK-F16-06 hook before HTTP shutdown.

Integration tests prove end-to-end qualify/retire with fenced commit, idempotency
replay, Retire/Maintenance wins, authorization, safe-denial.

Completes HTTP handler implementation.

Relates: FEATURE-0016, REQ-F16-04..09, ADH-2026-058, ADH-2026-060, ADH-2026-061.
```

---

### TASK-F16-12 — Feature-local conformance suite (122 cases, no omissions)

**Prerequisites:** TASK-F16-01 through TASK-F16-11 complete.

**Requirement authority:** Requirements §4.3 (exact conformance semantics ledger:
122 cases VS0-CF-F16-01..122), AC-F16-01 through AC-F16-12.

**Design authority:** Design §7 (test and conformance strategy), design §7.2
(conformance coverage obligation: every one of 122 cases realized), design §7.3
(test categories).

**Architecture decisions:** ADH-2026-058 full closure and ADH-2026-061
Maintenance-entry correction (no shared/downstream case used as F0016-local
proof).

**Applicable risks:** All six risks proven by feature-local conformance coverage:
Risk 1 (sole-committer status writes), Risk 2 (no external effect; restart clears
F0016-owned state, `externalCallCount=0`), Risk 3 (audit-before-publication; a
failed required append creates no durable AuditEvent or publication and preserves
all unrelated state; request-owner abort occurs only where the approved action
requires it, while failed Maintenance entry preserves in-flight qualification and
its linked idempotency reservation), Risk 4 (safe-denial non-disclosure),
Risk 5 (closed four-path/five-registration surface with side-effect-free
guard), Risk 6 (stale-fence non-publication).

**Writable paths:**
- `tests/conformance/feature_0016_test.go`
- `tests/conformance/feature_0016_fixtures.go`

**Tests:**
Feature-local conformance suite realizing all 122 approved cases
VS0-CF-F16-01..122 from requirements §4.3. The category list below is the
corrected, exhaustive grouping (every one of 01..122 appears in at least one
category; several IDs are intentionally cross-listed where design §7.3/§9.3
treats them as belonging to more than one theme, matching the design's own
non-partitioned categorization style):

- Create-field boundary and scope derivation (01,02,04,54..57,84,91,92).
- Authorization denial without required grant (03,15,18,37).
- Safe access before graph validation (09..13,17,19,65,66,98).
- Route/guard transport outcomes (44,45,48,105..110,113..118) and precedence
  (05,06,07,08,24,25,58..60,70..74,85,86).
- Authentication failures across all four path patterns / five registrations
  (61..64,93..97).
- Idempotency replay/reuse/isolation/eviction/panic/waiter
  (14,16,31..35,46,47,49,50,76..78,81,82,87..89,99).
- Malformed/duplicate/oversized body (51..53).
- Observer/facts/conclusions (20..23,67..69,90) and lifecycle/availability
  (26..30,36,38,39,41..43,75,79..83,88,89,100,101,119..122).
- Safe-accessible-without-grant item GET (111).
- Audit matrix and append-failure (39,40,82,99..104,112) and future boundary (46).

Equivalence families (VS0-CF-F16-106,112,114,115,118) execute every named member.
Determinism guaranteed by injected clock and target-bound fixtures; no test relies
on wall-clock elapse or network.

Each conformance test asserts exact registered `expectedState`, `expectedError`,
`expectedViolation` (where registered), and `expectedSideEffects`. No
shared/downstream case (`VS0-CF-F10`, `VS0-CF-X03`, `VS0-CF-HP01`) counted as
F0016-local proof.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -race ./tests/conformance/...
make feature-contract-check FEATURE=FEATURE-0016
```

**Acceptance criteria:**
- All 122 approved conformance cases VS0-CF-F16-01..122 exist as executable tests
  asserting exact registered outcomes; every one of 01..122 appears in at least
  one category grouping above with no gaps (03,15,18,37,51..53,61..64,93..97,111
  are explicitly present, correcting the prior omission).
- Equivalence families execute every named member.
- Deterministic with injected clock and target-bound fixtures; no wall-clock or
  network dependency.
- `make feature-contract-check FEATURE=FEATURE-0016` passes.
- All six risks proven by conformance coverage.

**Security impact:** Proves all security invariants (sole-committer,
safe-denial, closed route surface, waiter recheck, fenced commit).

**Observability impact:** Proves audit matrix and append-failure behavior.

**Exclusions:** None; this completes feature-local conformance.

**Commit message:**
```
test(f16): TASK-F16-12 — feature-local conformance suite (122 cases, no omissions)

Implements requirements §4.3 (122 cases VS0-CF-F16-01..122), AC-F16-01..12.

Creates feature-local conformance suite realizing all 122 approved cases, grouped
into ten categories covering every ID 01..122 with no gaps: create-field boundary
and scope derivation; authorization denial without required grant; safe access
before graph validation; route/guard transport outcomes and precedence;
authentication failures across all four path patterns / five registrations;
idempotency replay/reuse/isolation/eviction/panic/waiter; malformed/duplicate/oversized body;
observer/facts/conclusions and lifecycle/availability; safe-accessible-without-
grant item GET; audit matrix and append-failure and future boundary.

Corrects the prior omission of 03,15,18,37,51..53,61..64,93..97,111 — each now
appears in an explicit category.

Equivalence families execute every named member. Deterministic with injected clock
and target-bound fixtures; no wall-clock or network dependency.

Each test asserts exact registered expectedState, expectedError, expectedViolation,
expectedSideEffects. No shared/downstream case counted as F0016-local proof.

Proves all six risks: Risk 1 (sole-committer status writes), Risk 2 (no external
effect; externalCallCount=0), Risk 3 (audit-before-publication), Risk 4
(safe-denial non-disclosure), Risk 5 (closed four-path/five-registration
surface), Risk 6 (stale-fence non-publication).

Relates: FEATURE-0016, requirements §4.3, AC-F16-01..12, ADH-2026-058.
```

---

### ADH-066 remediation task — F16-R01: coherent backing-access bridge

**Prerequisites:** TASK-F16-01 through TASK-F16-12 complete; ADH-2026-066 is
approved and present in the FEATURE-0016 control manifest.

**Purpose:** Repair the already-implemented F0016 backing-access path without
reopening FEATURE-0015. This is a bounded post-plan remediation task, not a
thirteenth numbered vertical slice and not a change to the approved five-route
surface or the 128-case public contract.

**Requirement authority:** REQ-F16-02, REQ-F16-04, REQ-F16-05, REQ-F16-08, and
REQ-F16-09; existing VS0-CF-F16-09..13,17,19,29,30,65,66,98,119,120.

**Design authority:** DD-02, DD-04, DD-15, DD-17.

**Architecture authority:** ADH-2026-066 only; it authorizes the private bridge
and explicitly prohibits a FEATURE-0015 product-scope change.

**Writable paths:**

- `internal/cloudmodel/store.go`
- `internal/cloudmodel/store_test.go`
- `internal/executiontarget/`
- `internal/api/executiontarget_*.go`
- `internal/api/executiontarget_*_test.go`
- `tests/conformance/feature_0016_test.go`

**Implementation contract:**

- Add only a private store-backed read-lease primitive in `internal/cloudmodel`
  that acquires the existing store lock once and supplies immutable copies of
  both backing resources to an F0016 callback. It must not add an F0015 route,
  schema field, writer, lifecycle transition, conformance case, or public API.
- Implement the F0016-owned, principal-aware `BackingAccessProvider` on top of
  that lease. It alone evaluates inherited F0012 grants against the derived
  CloudProvider scope and returns safe denial, authorization denial, or the
  minimal immutable allowed snapshot.
- Replace sequential F0015 `GetParticipation`/`GetInfrastructureStack` use in
  F0016 create/qualify/GET/LIST paths. Final qualification holds the fresh
  backing lease through the ordered `lease → lifecycle mutex → AuditEvent
  append → target publication` critical section. No F0016 lifecycle-held path
  may acquire the lease.
- Compare only InfrastructureStack generation and viability fingerprint at
  qualification commit. Participation generation is never stored or evaluated
  as a fence.
- Treat participation suspension as unavailable for new admission only; never
  mutate or stop an existing InfrastructureStack, ExecutionTarget, or workload.

**Tests and acceptance criteria:**

- Deterministically prove a paired backing read cannot be interleaved by a
  FEATURE-0015 backing write while the F0016 provider evaluates safe access.
- Prove safe 404 and authorized 403 remain caller-relative and do not disclose
  raw backing existence.
- Prove a changed InfrastructureStack generation or viability fingerprint aborts
  before AuditEvent, publication, or replayable completion; prove no
  participation-generation fence exists.
- Prove GET/LIST freshness changes only response-only
  `effectiveAvailability`, never target state or ETag.
- Prove the required cross-package lock order with race tests and no deadlock.
- Run F0015 regression (`go test ./internal/cloudmodel/...` and
  `make feature-contract-check FEATURE=FEATURE-0015`) plus the F0016 checks
  below. No F0015 feature artifact is changed.

**Verification commands:**

```bash
make fmt
make vet
make test
go test -race ./internal/cloudmodel/... ./internal/executiontarget/... ./internal/api/... ./tests/conformance/...
make feature-contract-check FEATURE=FEATURE-0015
make feature-contract-check FEATURE=FEATURE-0016
make feature-0016-architecture-readiness
make ff-feature-gate FEATURE=FEATURE-0016
git diff --check
```

**Security and observability:** Safe-denial remains principal-relative; no raw
backing resource leaks. The bridge preserves the inherited audit-before-
publication boundary and adds no new AuditEvent type or public violation.

**Exclusions:** No F0015 product code outside the private `store.go` lease, no
F0015 docs/specs/control changes, and no new route, schema, state, conformance
ID, external call, provisioning, or placement behavior.

**Commit message:**

```text
fix(f16): add coherent principal-aware backing access bridge
```

---

## Final Verification Checkpoint

**Checkpoint:** Repository-wide verification and command execution report
(strictly verification-only; not a numbered task).

**Prerequisites:** TASK-F16-01 through TASK-F16-12 and remediation task F16-R01 complete.

**Requirement authority:** All REQ-F16-01..11, all AC-F16-01..12.

**Design authority:** All design decisions DD-01..10, design §7.5 (verification
commands), design §10.3 (resolved report: no unresolved decision).

**Architecture decisions:** All controlling decisions and handoffs; slice0-contract
anti-drift steering §9.

**Applicable risks:** All six risks proven by repository-wide verification.

**Writable paths:**
- None. This task has zero authority to create, modify, or delete any
  **repository** file, including `.kiro/steering/*`, `scripts/*`, configuration,
  or any tracked/untracked repository path. There is no escape-hatch clause of
  any kind. Commands may use only ephemeral OS temporary files and normal tool
  cache artifacts outside the repository; they grant no authority to remediate
  failures. If verification
  discovers drift, an unspecified observable outcome, or any condition that
  would require a file change to resolve, this task MUST report the exact
  failure and stop rather than self-authorizing any fix. It may name a
  manifest-approved stop condition from §1.3 only when that condition's defined
  cause actually applies.

**Tests:**
Run all repository-wide verification commands and report their exact pass/fail
results. This task performs no test authoring; TASK-F16-12 already authored the
conformance suite.

```bash
gofmt -l .
make test
make vet
go test -race ./...
make feature-contract-check FEATURE=FEATURE-0016
make vs000-contract-check
make phase2r-drift-check
git status
git diff --check

if lsof -n -iTCP:8080 -sTCP:LISTEN >/dev/null 2>&1; then
  echo "ERROR: smoke port 8080 is already occupied" >&2
  exit 1
fi
smoke_dir="$(mktemp -d)"
smoke_bin="$smoke_dir/sovrunn-api"
smoke_log="$smoke_dir/server.log"
smoke_pid=""
cleanup_smoke() {
  if [ -n "$smoke_pid" ] && kill -0 "$smoke_pid" 2>/dev/null; then
    kill -TERM "$smoke_pid" 2>/dev/null || true
  fi
  if [ -n "$smoke_pid" ]; then wait "$smoke_pid" 2>/dev/null || true; fi
  rm -rf "$smoke_dir"
}
trap cleanup_smoke EXIT INT TERM
go build -o "$smoke_bin" ./cmd/sovrunn-api
"$smoke_bin" >"$smoke_log" 2>&1 &
smoke_pid=$!
smoke_ready=0
for attempt in 1 2 3 4 5; do
  if kill -0 "$smoke_pid" 2>/dev/null && \
    curl --fail --connect-timeout 2 --max-time 5 http://127.0.0.1:8080/healthz && \
    curl --fail --connect-timeout 2 --max-time 5 http://127.0.0.1:8080/readyz; then
    smoke_ready=1
    break
  fi
  sleep 1
done
if [ "$smoke_ready" -ne 1 ]; then
  echo "ERROR: launched API-server did not pass both health checks" >&2
  exit 1
fi
kill -TERM "$smoke_pid"
if ! wait "$smoke_pid"; then
  echo "ERROR: launched API-server did not shut down cleanly" >&2
  exit 1
fi
smoke_pid=""
if lsof -n -iTCP:8080 -sTCP:LISTEN >/dev/null 2>&1; then
  echo "ERROR: API-server still listens on smoke port after shutdown" >&2
  exit 1
fi
rm -rf "$smoke_dir"
trap - EXIT INT TERM
```

`gofmt -l .` is a read-only repository formatting check (lists any misformatted
files without modifying them); it replaces `make fmt` here because this
checkpoint has zero repository-write authority and must never run a
repository-mutating formatter.

Confirm by inspection only (no file edits):
- No F0016 code change in any architecture-only document.
- No ResourcePool, ProviderCapability, generic Provider, six-scope vocabulary,
  SovrunnInstallation as active, CanonicalMigrationPlan/Record, or any
  superseded/prohibited concept introduced.
- No shared/downstream conformance case used as F0016-local proof.
- No future-feature concept (FEATURE-0017/0019/0022/0023/0024) introduced or
  activated.
- ExecutionTarget `status.*` activated only by FEATURE-0016.
- ServiceRegion and its fields introduced/activated only by FEATURE-0022.
- No real adapter, credential, external call, Crossplane dependency, placement,
  provisioning, plugin execution, or customer IaaS exposure.
- `git status` shows a clean working tree.
- `git diff --check` reports no whitespace errors or unresolved conflict
  markers.

**Verification commands:**
```bash
gofmt -l .
make test
make vet
go test -race ./...
make feature-contract-check FEATURE=FEATURE-0016
make vs000-contract-check
make phase2r-drift-check
git status
git diff --check

if lsof -n -iTCP:8080 -sTCP:LISTEN >/dev/null 2>&1; then
  echo "ERROR: smoke port 8080 is already occupied" >&2
  exit 1
fi
smoke_dir="$(mktemp -d)"
smoke_bin="$smoke_dir/sovrunn-api"
smoke_log="$smoke_dir/server.log"
smoke_pid=""
cleanup_smoke() {
  if [ -n "$smoke_pid" ] && kill -0 "$smoke_pid" 2>/dev/null; then
    kill -TERM "$smoke_pid" 2>/dev/null || true
  fi
  if [ -n "$smoke_pid" ]; then wait "$smoke_pid" 2>/dev/null || true; fi
  rm -rf "$smoke_dir"
}
trap cleanup_smoke EXIT INT TERM
go build -o "$smoke_bin" ./cmd/sovrunn-api
"$smoke_bin" >"$smoke_log" 2>&1 &
smoke_pid=$!
smoke_ready=0
for attempt in 1 2 3 4 5; do
  if kill -0 "$smoke_pid" 2>/dev/null && \
    curl --fail --connect-timeout 2 --max-time 5 http://127.0.0.1:8080/healthz && \
    curl --fail --connect-timeout 2 --max-time 5 http://127.0.0.1:8080/readyz; then
    smoke_ready=1
    break
  fi
  sleep 1
done
if [ "$smoke_ready" -ne 1 ]; then
  echo "ERROR: launched API-server did not pass both health checks" >&2
  exit 1
fi
kill -TERM "$smoke_pid"
if ! wait "$smoke_pid"; then
  echo "ERROR: launched API-server did not shut down cleanly" >&2
  exit 1
fi
smoke_pid=""
if lsof -n -iTCP:8080 -sTCP:LISTEN >/dev/null 2>&1; then
  echo "ERROR: API-server still listens on smoke port after shutdown" >&2
  exit 1
fi
rm -rf "$smoke_dir"
trap - EXIT INT TERM
```

**Acceptance criteria:**
- `gofmt -l .` reports no misformatted files (empty output); this checkpoint
  never runs a repository-mutating formatter. Its `mktemp -d` server binary and
  log, plus normal tool caches, are ephemeral artifacts outside the repository;
  the known temporary directory is removed after verified server termination and
  they do not grant remediation authority.
- Every command listed above has been run and its exact pass/fail result is
  reported verbatim; this task does not paraphrase or summarize a command result
  as "passed" without showing the command's actual output status.
- `git status` shows a clean working tree.
- `git diff --check` reports no whitespace errors or conflict markers.
- Controlled API-server smoke builds a binary only in an external `mktemp -d`
  directory, starts that binary once (so `smoke_pid` is the actual server PID),
  records success only when both `/healthz` and `/readyz` pass within the five
  bounded attempts (otherwise exits nonzero before cleanup can mask failure),
  sends it `SIGTERM`, checks clean exit, verifies port 8080 is no longer
  listening, and removes the known temporary directory before reporting success.
- No architecture-only doc contains F0016 code changes.
- No prohibited/superseded concept introduced.
- No shared/downstream case used as F0016-local proof.
- No future-feature concept introduced or activated.
- No real external effect or customer IaaS exposure.
- This task's report explicitly states that running these commands and
  reporting their results is this task's entire scope, and that this task does
  not itself constitute or guarantee that
  `make ff-feature-gate FEATURE=FEATURE-0016` will pass — that command depends
  on conditions outside this task's control (human review gates and other
  repository state) and remains a separate, later, human-gated step per the
  manifest's `human_gates: ["architecture", "executable_plan", "final"]`. This
  task neither runs `make ff-feature-gate` as a self-granted completion proof
  nor claims final feature-gate approval.
- If any command fails, this task reports the exact failure and STOPS; it does
  not edit any file to make a failing command pass. For an ordinary
  verification failure (e.g. a failing test, a lint/format finding reported by
  `gofmt -l .`, or a failed build/vet command), stopping and reporting the
  exact failure is sufficient — this does NOT necessarily require emitting one
  of the five manifest-approved stop conditions. One of the five
  manifest-approved stop conditions (`ARCHITECTURE_DECISION_REQUIRED`,
  `REQUIREMENT_CLARIFICATION_REQUIRED`, `BOUNDARY_CHANGE_REQUIRED`,
  `DEPENDENCY_APPROVAL_REQUIRED`, `SECURITY_REVIEW_REQUIRED`) is emitted only
  when that specific condition's defined cause actually applies — for example,
  discovering an unspecified observable outcome or missing architecture
  authority warrants `ARCHITECTURE_DECISION_REQUIRED`, and discovering a
  security invariant violation warrants `SECURITY_REVIEW_REQUIRED` — never as
  an automatic label for every failure. This task does not claim or imply that
  every failure is automatically an architecture or security stop condition.
  The existing "no self-authorized fix" and "zero writable paths" invariants
  remain unchanged regardless of which of these two reporting forms applies.

**Security impact:** Confirms all security invariants hold across the full
codebase by inspection and command execution only; grants itself no authority to
remediate.

**Observability impact:** Confirms audit matrix completeness by running the
listed commands and reporting their results; makes no code or doc change.

**Exclusions:** No repository file creation, modification, or deletion. The
explicitly permitted external temporary files and normal external tool caches
do not grant remediation authority. No self-authorized steering or check-script edit. No running or self-granting of
`make ff-feature-gate FEATURE=FEATURE-0016` as this task's own completion proof.

This checkpoint performs no commit; it is a verification checkpoint, not an
implementation task, and therefore has no commit message block.

---

## 3. No-task ledger

The following design elements require no separate implementation task because they
are contract-only reuse, invariants proven by other tasks, governance constraints,
or explicitly excluded concepts:

### 3.1 Contract-only reuse (NO_TASK — consumed by reference, not forked)

| Element | Owner | Reused By | Disposition |
|---------|-------|-----------|-------------|
| Problem Details + HTTP/URN mapping | FEATURE-0012 / `internal/apiproblem` | All HTTP handlers | Reused; F0016 adds no top-level code |
| ObjectMeta / TypedRef / Condition | FEATURE-0012 / `internal/apimeta`, `apiref`, `apicond` | All domain types | Reused as-is |
| Validation pipeline stages, media, If-Match/ETag | FEATURE-0012 / `internal/apivalid` | All HTTP handlers | Reused; F0016 wires order, adds no stage semantics |
| AuditEvent atomic append-before-publication | FEATURE-0013 / `internal/decision` | TASK-F16-02 (audit assembly) | Reused; not redefined |
| Server-resolved authorization, safe denial, digest conventions | FEATURE-0015 / `internal/cloudmodel` | All HTTP handlers | Reused; FEATURE-0015 unchanged |
| Inherited violation codes (`VS0_STATUS_FIELD_WRITE`, `VS0_SYSTEM_OWNED_FIELD_WRITE`, `VS0_AUTHORIZATION_SAFE_DENIAL`, `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`) | Prior features | All HTTP handlers | Reused; not introduced |

### 3.2 Invariants proven by existing tasks (NO_TASK — proven by implementation, not separately tested)

| Invariant | Proven By | Registry/Conformance |
|-----------|-----------|----------------------|
| Clients never write status; sole committer is the lifecycle service | TASK-F16-04 (DD-02; handlers never mutate status) | VS0-WRITER-006; VS0-CF-F16-03,04,91,92 |
| No external effect, credential, or persistence | TASK-F16-05 (observer external-effect-free), TASK-F16-12 (VS0-CF-F16-46 restart proof) | REQ-F16-11; VS0-CF-F16-46 |
| Audit-before-publication; no partial mutation on append failure | TASK-F16-02 (assembly), TASK-F16-04 (create/retire), TASK-F16-05 (qualification), TASK-F16-06 (expiry/Maintenance), TASK-F16-09/10/11 (handler append owners), TASK-F16-12 (conformance) | REQ-F16-09; design §5.3 |
| Safe denial never discloses inaccessible resources | TASK-F16-09/10/11 (collection, item, and action safe denial), TASK-F16-12 (VS0-CF-F16-09,17,19,66) | REQ-F16-02; design §5.1 step 6 |
| Closed four-path/five-registration surface; transport guard has no side effect | TASK-F16-08 (guard), TASK-F16-09/10/11 (route registration), TASK-F16-12 (VS0-CF-F16-44,113..118) | REQ-F16-03; DD-08 |
| Stale fences cannot publish a conclusion or replayable completion | TASK-F16-05 (fenced commit), TASK-F16-11 (action wiring), TASK-F16-12 (VS0-CF-F16-29,30,75,81,82) | REQ-F16-08; DD-04 |

### 3.3 Governance constraints (NO_TASK — architecture boundary, not implementation)

| Constraint | Authority | Enforcement |
|------------|-----------|-------------|
| Seven canonical scope kinds only | DEC-0037 | Architecture boundary; no F0016 scope extension |
| ExecutionTarget owned entirely by FEATURE-0016 | ADH-2026-045 decision 7 | No FEATURE-0015 ExecutionTarget route/schema/writer/conformance |
| ExecutionTarget `status.*` introduced and activated only by FEATURE-0016 | ADH-2026-058 clause 7 | No prior-feature or future-feature status activation |
| ServiceRegion fields introduced/activated only by FEATURE-0022 | VS0-SCHEMA-056 `fieldOwnership` | FEATURE-0022 dependency; no F0016 activation |
| CloudProviderParticipation `spec.providerSelectionModes` introduced/activated only by FEATURE-0021 | VS0-SCHEMA-010 `fieldOwnership` | FEATURE-0021 dependency; no F0016 activation |
| No active `status.availability`, `Draining`, or publicly projected `Qualifying` | ADH-2026-058 replacement of placeholder | No alias, migration, or legacy compatibility |

### 3.4 Explicit exclusions (EXCLUDED — must never be designed or built by F0016)

| Excluded Element | Reason | Negative Acceptance |
|------------------|--------|---------------------|
| Real adapter/provider-native type/credential/SecretRef/external call/external persistence | REQ-F16-11; DEC-0036/0042 | VS0-CF-F16-46 (restart proves no external call) |
| Placement, provisioning, plugin execution, realization, plugin taxonomy, Crossplane dependency | F0024/Phase 3 owns realization | VS0-CF-F16-46 |
| Customer target projection | FEATURE-0023 owns ServicePlacement | No customer-facing target route |
| Maintenance-notice resource, IAM resource | Out of F0016 scope | No maintenance-notice/IAM schema |
| Persisted `status.availability`, `Draining`, public `Qualifying`, public `conditions`, `scopeRef`, `spec.participationRef`, `adapterAuthorityRef` | ADH-2026-058 replacement | No alias or legacy field |
| FEATURE-0017 PolicyEvaluationRequest/Result/PolicyEngineAdapter | FEATURE-0017 owner | No F0016 policy route/field/record |
| FEATURE-0019 SovereigntyProfile/SovereigntyFactSet/EvidenceRecord | FEATURE-0019 owner | No F0016 sovereignty route/field/record |
| FEATURE-0022 ServiceRegion / `spec.executionTargetRefs` / `ServiceRegion.status.availability` | FEATURE-0022 owner | F0016 prerequisite only; no F0016 activation |
| FEATURE-0023 DecisionProfiles/ServicePlacement/PlacementDecision/customer target projection | FEATURE-0023 owner | No customer-facing target projection |
| FEATURE-0024 PluginExecution/ServiceDeploymentPlan/real realization | FEATURE-0024 owner | VS0-CF-F16-46 (externalCallCount=0) |

---

## 4. Canonical coverage ledger

Requirement/design/decision/risk coverage ledger (every REQ and AC mapped
exactly once).

### 4.1 Requirements coverage (every REQ mapped exactly once)

| REQ ID | Requirement | Implemented By | Classification |
|--------|-------------|----------------|----------------|
| REQ-F16-01 | Register ExecutionTarget; closed create fields; derived scope | TASK-F16-01, TASK-F16-04, TASK-F16-09 | IMPLEMENT |
| REQ-F16-02 | Safe reference access before graph validation; backing viability; uniqueness | TASK-F16-01, TASK-F16-04, TASK-F16-05, TASK-F16-09, TASK-F16-10, TASK-F16-11 | IMPLEMENT |
| REQ-F16-03 | Closed four-path/five-registration surface; pre-ServeMux transport guard | TASK-F16-08, TASK-F16-09, TASK-F16-10, TASK-F16-11 | IMPLEMENT |
| REQ-F16-04 | Request grammar; representation; safe projection | TASK-F16-07, TASK-F16-09, TASK-F16-10, TASK-F16-11 | IMPLEMENT |
| REQ-F16-05 | Executable request precedence | TASK-F16-09, TASK-F16-10, TASK-F16-11 | IMPLEMENT |
| REQ-F16-06 | Idempotency namespace; digest; retention; isolation | TASK-F16-03, TASK-F16-04, TASK-F16-09, TASK-F16-11 | IMPLEMENT |
| REQ-F16-07 | Synthetic observer; four named facts; missing-fixture/timeout | TASK-F16-05 | IMPLEMENT |
| REQ-F16-08 | Sole committer; lifecycle; qualification; effectiveAvailability | TASK-F16-04, TASK-F16-05, TASK-F16-07, TASK-F16-11 | IMPLEMENT |
| REQ-F16-09 | Expiry; maintenance trigger; retirement; shutdown; audit matrix | TASK-F16-02, TASK-F16-04, TASK-F16-05, TASK-F16-06, TASK-F16-09, TASK-F16-10, TASK-F16-11 | IMPLEMENT |
| REQ-F16-10 | Closed local violation set | TASK-F16-01, all HTTP handler tasks | IMPLEMENT (violation surface) |
| REQ-F16-11 | No external effect; restart clears F0016-owned state | TASK-F16-05, TASK-F16-12 | IMPLEMENT (restart proof) |

### 4.2 Acceptance criteria coverage (every AC mapped exactly once)

| AC ID | Criterion | Proven By | Conformance |
|-------|-----------|-----------|-------------|
| AC-F16-01 | Closed create-field boundary; derived scope; Active/Unqualified at epoch 0 | TASK-F16-01, TASK-F16-04, TASK-F16-09, TASK-F16-12 | VS0-CF-F16-01,02,57,84 |
| AC-F16-02 | Safe reference access before graph validation | TASK-F16-09, TASK-F16-12 | VS0-CF-F16-09..13,65,98 |
| AC-F16-03 | Four path patterns, five method/path registrations; pre-ServeMux guard transport outcomes | TASK-F16-08, TASK-F16-09, TASK-F16-10, TASK-F16-11, TASK-F16-12 | VS0-CF-F16-44,61..64,66,93..96,105..118 |
| AC-F16-04 | POST media/key; action If-Match/zero-body; header-form precedence | TASK-F16-09, TASK-F16-11, TASK-F16-12 | VS0-CF-F16-05,06,24,25,58..60,70..74,85..86 |
| AC-F16-05 | Replay/reuse/isolation across create/qualify/retire; GET/LIST ignore idempotency | TASK-F16-03, TASK-F16-09, TASK-F16-10, TASK-F16-11, TASK-F16-12 | VS0-CF-F16-14,16,31..35,46..47,49..50,76..78,81..82,86..89,99 |
| AC-F16-06 | Four facts; deterministic missing-fixture/timeout | TASK-F16-05, TASK-F16-12 | VS0-CF-F16-20..23,67..69,90 |
| AC-F16-07 | Sole committer; non-projected Qualifying; effectiveAvailability truth table | TASK-F16-04, TASK-F16-05, TASK-F16-07, TASK-F16-11, TASK-F16-12 | VS0-CF-F16-01,20..22,26..30,36,39,41..43,75,77,79..83,87..88,99..101,119..122 |
| AC-F16-08 | Retire/Maintenance win over in-flight qualification | TASK-F16-05, TASK-F16-06, TASK-F16-11, TASK-F16-12 | VS0-CF-F16-26..28,75,80,83,88..89 |
| AC-F16-09 | Expiry retry; shutdown ordering | TASK-F16-06, TASK-F16-11, TASK-F16-12 | VS0-CF-F16-41,77,79 |
| AC-F16-10 | Audit matrix; append-failure behavior | TASK-F16-02, TASK-F16-04, TASK-F16-05, TASK-F16-06, TASK-F16-09, TASK-F16-10, TASK-F16-11, TASK-F16-12 | VS0-CF-F16-39,40,82,99..104,112 |
| AC-F16-11 | Only eight closed local violation codes | TASK-F16-01, all HTTP handler tasks, TASK-F16-12 | VS0-CF-F16-10..12,26..30,38,75 |
| AC-F16-12 | No real external effect; restart proof | TASK-F16-05, TASK-F16-12 | VS0-CF-F16-46 |

### 4.3 Design decisions coverage (every DD mapped exactly once)

| DD ID | Decision | Implemented By | Mechanics |
|-------|----------|----------------|-----------|
| DD-01 | Domain service package placement | TASK-F16-01 | `internal/executiontarget/` sibling to FEATURE-0015 |
| DD-02 | Sole committer is a single in-process service type | TASK-F16-04, TASK-F16-05, TASK-F16-06 | `ExecutionTargetLifecycleService` with one service-owned mutex; commit methods added across three tasks into the same file |
| DD-03 | Synthetic observer is a target-bound fixture reader | TASK-F16-05 | `sovrunn.synthetic-iaas-observer/v1` external-effect-free |
| DD-04 | Qualification is synchronous with captured-fence reservation | TASK-F16-05, TASK-F16-11 | Qualifying reservation; ADH-2026-060 ordered commit predicate |
| DD-05 | Idempotency reservation table is F0016-owned and in-process | TASK-F16-03, TASK-F16-04 | F0016-owned namespace/digest/retention/eviction; completion wiring integrated into lifecycle commit |
| DD-06 | Injected clock drives expiry retry and fact freshness | TASK-F16-06 | Injected Clock interface with real/fake implementations |
| DD-07 | Maintenance enters/clears only through target-bound fixture trigger | TASK-F16-06 | Deterministic target-bound fixture trigger intake; commit methods added into lifecycle.go |
| DD-08 | Pre-ServeMux transport guard wraps F0016 mux registrations | TASK-F16-08, TASK-F16-09, TASK-F16-10, TASK-F16-11 | Fixed pre-ServeMux method/path guard built standalone, then wrapped around the five registrations as each is added |
| DD-09 | Request pipeline is an explicit ordered sequence | TASK-F16-09, TASK-F16-10, TASK-F16-11 | Explicit short-circuiting pipeline stages per design §5.1 |
| DD-10 | Reuse inherited envelopes, errors, and helpers by reference | All HTTP handler tasks | Problem Details, ObjectMeta, TypedRef, Condition, media, validation, AuditEvent reused |

### 4.4 Architecture decisions coverage (all controlling DEC/ADH)

| Authority | Role in FEATURE-0016 | Implemented By |
|-----------|----------------------|----------------|
| DEC-0036 | Implementation-neutral adapter boundary | TASK-F16-05 (no native types, raw credentials, vendor errors) |
| DEC-0042 | Canonical cloud model; superseded ResourcePool/ProviderCapability/generic Provider | TASK-F16-01 (ExecutionTarget canonical owner); no ResourcePool/ProviderCapability |
| DEC-0057 | ExecutionTarget qualification/availability semantics | TASK-F16-04, TASK-F16-05, TASK-F16-07 (effectiveAvailability, maintenance epoch) |
| ADH-2026-025, ADH-2026-040, ADH-2026-042, ADH-2026-045 | Prior controlling handoffs | All tasks (baseline + canonical model + ExecutionTarget delegation to F0016) |
| ADH-2026-058 | Executable contract closure; replaced placeholder VS0-SCHEMA-015..017/VS0-STATE-004 | All tasks (baseline executable closure; clauses mapped to exact tasks) |
| ADH-2026-060 | Maintenance-race outcome correction | TASK-F16-05, TASK-F16-11 (exact ordered commit predicate; F16-75/F16-89 mutually exclusive) |
| ADH-2026-061 | Maintenance-entry link-clearing correction | TASK-F16-06, TASK-F16-11 (every successful entry clears both links; only the linked in-flight qualification reservation aborts) |
| ADH-2026-062 | Automation task-inventory path-classification correction | Tasks-stage inventory checker only; no product or implementation semantic change |

### 4.5 Risk coverage (all six risks proven)

| Risk | Statement | Enforcement | Proven By |
|------|-----------|-------------|-----------|
| Risk 1 | Clients never write status; sole committer only | DD-02; VS0-WRITER-006; handlers never mutate status | TASK-F16-04, TASK-F16-12 (VS0-CF-F16-03,04,91,92) |
| Risk 2 | No external effect/credential/persistence | DD-03; observer external-effect-free; restart clears F0016-owned state | TASK-F16-05, TASK-F16-12 (VS0-CF-F16-46 externalCallCount=0) |
| Risk 3 | Audit-before-publication; no partial mutation on append failure | DD-04; design §5.3; lifecycle service/scheduler/handlers append before publication or return denial | TASK-F16-02, TASK-F16-04, TASK-F16-05, TASK-F16-06, TASK-F16-09, TASK-F16-10, TASK-F16-11, TASK-F16-12 (VS0-CF-F16-39,40,82,99..104) |
| Risk 4 | Safe denial never discloses inaccessible resources | Design §5.1 step 6; §6.2; safe projection omits facts/results/observer/handles/links | TASK-F16-07, TASK-F16-09, TASK-F16-10, TASK-F16-11, TASK-F16-12 (VS0-CF-F16-09,17,19,66) |
| Risk 5 | Closed four-path/five-registration surface; transport guard has no side effect | DD-08; design §4.5; transport guard before authentication | TASK-F16-08, TASK-F16-09, TASK-F16-10, TASK-F16-11, TASK-F16-12 (VS0-CF-F16-44,113..118) |
| Risk 6 | Stale fences cannot publish conclusion or replayable completion | DD-04; design §5.1 step 12; qualification fence recheck at commit | TASK-F16-05, TASK-F16-11, TASK-F16-12 (VS0-CF-F16-29,30,75,81,82) |

---

## 5. Orphan and unresolved report

### 5.1 Orphan requirements/AC check

All 11 REQ-F16-* and all 12 AC-F16-* identifiers from approved requirements are
mapped exactly once in §4.1–4.2 above. No orphan requirement or acceptance
criterion exists.

### 5.2 Orphan design decision check

All 10 design decisions DD-01..DD-10 from approved design are mapped exactly once
in §4.3 above. No orphan design decision exists.

### 5.3 Orphan conformance case check

All 122 conformance cases VS0-CF-F16-01..122 from approved requirements §4.3 are
realized in TASK-F16-12, with every ID appearing in at least one category
grouping (see TASK-F16-12's corrected category list, which explicitly restores
03,15,18,37,51..53,61..64,93..97,111 that were omitted in the prior revision). No
orphan conformance case exists. No shared/downstream case (`VS0-CF-F10`,
`VS0-CF-X03`, `VS0-CF-HP01`) is counted as F0016-local proof.

### 5.4 Unresolved decisions

None. ADH-2026-060, ADH-2026-061, and ADH-2026-062 were approved and applied to the manifest-controlled
architecture authorities (registry, architecture doc, traceability, steering,
requirements.md, design.md, and the readiness/contract-check/vs000-check scripts)
before this tasks stage began. The Maintenance-race outcome between F16-75
(inactive-marker, changed-maintenance-epoch, `STALE_RESOURCE_VERSION`/412) and
F16-89 (active-marker, Maintenance-wins, `CONFLICT`/409) is mutually exclusive and
ordered per the predicate fixed in DD-04, design §5.1 step 12, and realized in
TASK-F16-05/11. ADH-2026-061 is the higher-precedence correction consumed by
TASK-F16-06/11: every successful Maintenance entry clears both current links,
and only its linked in-flight qualification/idempotency reservation aborts.
No unresolved architecture, requirements, design, or task
decision remains for this feature.

### 5.5 Task graph restructuring change log (this revision)

This revision restructured the task graph per founder-approved reviewer required
changes, without changing any approved requirement, design decision, route,
field, state, error code, or conformance-case semantic meaning:

- Old TASK-F16-01 (domain types + store + lifecycle service in one task) was
  split into new TASK-F16-01 (domain types, violations, store only) and new
  TASK-F16-04 (lifecycle service, depending on audit and idempotency).
- Old TASK-F16-04 (audit assembly) was renumbered to new TASK-F16-02 and moved
  earlier so it exists before the lifecycle service that depends on it.
- Old TASK-F16-03 (idempotency table) was renumbered to new TASK-F16-03 (no
  change in position; now explicitly feeds into new TASK-F16-04).
- Old TASK-F16-02 (observer + qualification engine) was renumbered to new
  TASK-F16-05 and now writes into `lifecycle.go` (adding the qualify-commit
  dispatch method) in addition to `observer.go`/`qualify.go`.
- Old TASK-F16-05 (scheduler) was renumbered to new TASK-F16-06 and now writes
  into `lifecycle.go` (adding the maintenance-enter/clear commit methods and the
  shutdown-ordering hook) in addition to `scheduler.go`/`clock.go`.
- Old TASK-F16-06 (projection) was renumbered to new TASK-F16-07 with updated
  prerequisites (now depends on the lifecycle service, qualification engine, and
  scheduler, since it projects their real outputs).
- Old TASK-F16-07 (guard + route registration with 501 placeholders) was split
  into new TASK-F16-08 (guard only, no route registration, no placeholder) and
  the route-registration responsibility was distributed into new
  TASK-F16-09/10/11 (each registers exactly the ServeMux patterns for the route
  it implements).
- Old TASK-F16-08/09/10 (collection/item/action handlers) were renumbered to new
  TASK-F16-09/10/11 respectively, each gaining `internal/server/routes.go` in
  its writable paths.
- Old TASK-F16-11 (conformance suite) was renumbered to new TASK-F16-12 with a
  corrected, exhaustive category list.
- Old TASK-F16-12 (repository-wide verification) was converted into the final
  verification checkpoint (no longer a numbered `TASK-F16-*` task, and no
  longer carrying a commit-message block), with its writable-paths escape
  hatch removed and `git diff --check` added to its verification commands.

The current TASK-F16-01 through TASK-F16-12 identifiers are the authoritative
execution plan. This history records a prior draft's split (old TASK-F16-01
became separate domain/store and lifecycle tasks), task redistribution, and
renumbering; it is explanatory only and does not claim identifier preservation
or a universal 1:1 mapping to the prior draft.

---

## Notes

- Tasks in this document are not marked with the `*` optional-test postfix
  convention because every test sub-bullet here is a required sub-part of its
  parent task's acceptance criteria (per the approved design's Correctness
  Properties in design.md, all listed unit/integration/conformance tests are
  required proof of a Property or Risk, not optional coverage); no task in this
  plan may be marked complete without its listed tests passing.
- TASK-F16-08 (guard) and TASK-F16-09/10/11 (route registration) are
  deliberately split so that no route is ever registered before its real
  handler exists; this is the mechanism that removes the previously-present 501
  placeholder state.
- TASK-F16-04, TASK-F16-05, and TASK-F16-06 each add methods to the same
  `internal/executiontarget/lifecycle.go` file in three separate, ordered waves;
  this is intentional sole-committer integration, not accidental file
  contention — the dependency graph places them in three different waves so
  they never conflict.
- The final verification checkpoint is strictly verification-only: it has no
  writable paths and no authority to edit steering or check scripts under any
  condition, including discovered drift. It is intentionally not a
  `TASK-F16-*` heading and carries no commit-message block, since it performs
  no commit.
- Checkpoints are folded into the final verification checkpoint's own
  acceptance criteria (run all verification commands, report exact results)
  rather than being a separate top-level task, since this stage's twelve-task
  count and structure were fixed by the reviewer's required restructuring.

## 6. Model execution report

- **Tool:** kiro
- **Stage or task:** FEATURE-0016 Tasks
  (adapter-boundary-and-executiontarget-qualification/tasks.md) — revision per
  founder-approved reviewer required changes
- **Recommended priority list:** Task decomposition (claude-sonnet-4.5, effort
  medium); Fallback (claude-opus-4.8, effort high); (third model not specified in
  prompt — only two labels provided)
- **Selected model:** claude-sonnet-4.5
- **Effort/reasoning setting:** medium
- **Fallback used:** no
- **Fallback reason:** none

---

STAGE_STATUS: COMPLETE
