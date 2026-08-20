# Design Document

## Overview

This section orients readers to the document outline; the substantive overview
is given by the Introduction immediately below and by §1 (Identity, stage,
inputs, and closed boundary), which together state the stage, the authorities
this design implements without modification, and the closed boundary within
which every mechanic below is chosen. This section adds no content beyond the
Introduction and §1.

---

## Introduction

This document is the FEATURE-0016 Design stage for Adapter Boundary and
ExecutionTarget Qualification. It translates the approved FEATURE-0016
requirements (`.kiro/specs/adapter-boundary-and-executiontarget-qualification/requirements.md`),
the architecture boundary
(`docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`),
and ADH-2026-058 into an implementable representation. It introduces no product
semantics, transfers no ownership, and designs no excluded adjacent-feature
mechanism. It chooses only the mechanics explicitly delegated to design by
requirements §9 and the architecture boundary: in-memory registry data
structures, handler wiring behind the fixed pre-ServeMux guard, idempotency
reservation/waiter internal data structures, injected-clock scheduling
internals, package/file layout, locks, and internal helper composition.

Every observable behavior below is fixed by the authorities above. Where this
document names a struct, package, or field name, that name is a design-local
mechanic; it never adds, renames, or removes an observable field, route, state,
error code, projection member, audit effect, or violation code.

---

## 1. Identity, stage, inputs, and closed boundary

### 1.1 Identity

| Field | Value |
|-------|-------|
| Feature | FEATURE-0016 — Adapter Boundary and ExecutionTarget Qualification |
| Stage | Design |
| Kiro slug | `adapter-boundary-and-executiontarget-qualification` |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase | 2R |
| Order | 6 |
| Owned resources | `ExecutionTarget` (VS0-SCHEMA-015), `NormalizedTargetFactSet` (VS0-SCHEMA-016), `TargetQualificationResult` (VS0-SCHEMA-017) |
| Owned state machine | `VS0-STATE-004` |
| Owned writer | `VS0-WRITER-006` (`ExecutionTargetLifecycleService`) |
| Controlling decisions | DEC-0036, DEC-0042, DEC-0057 |
| Controlling handoffs | ADH-2026-025, ADH-2026-040, ADH-2026-042, ADH-2026-045, ADH-2026-058, ADH-2026-060 |
| Local conformance | VS0-CF-F16-01..122 |

### 1.2 Stage inputs (consumed by reference, not redefined)

- Approved requirements (this spec's `requirements.md`): the sole source of REQ
  and AC meanings and the 122-case conformance ledger.
- `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`
  together with ADH-2026-058 and ADH-2026-060: the current semantic authority
  for VS0-SCHEMA-015..017 and VS0-STATE-004.
- FEATURE-0012 (`internal/apimeta`, `internal/apiproblem`, `internal/apivalid`,
  `internal/apiref`, `internal/apicond`, `internal/apischema`): ObjectMeta,
  TypedRef, Problem Details, Condition, media semantics, optimistic-concurrency
  conventions — reused, not forked.
- FEATURE-0013 (`internal/decision`): AuditEvent atomic append-before-publication
  semantics — reused, not redefined.
- FEATURE-0015 (`internal/cloudmodel`): `CloudProviderParticipation` and
  `InfrastructureStack` read-only prior-feature authority, plus the
  server-resolved authorization, safe-denial, and idempotency conventions to be
  reused — never modified.

### 1.3 Closed boundary for this stage

This design modifies only
`.kiro/specs/adapter-boundary-and-executiontarget-qualification/design.md`. It
generates no source code, schema, task, automation, architecture, prompt, or
validator. It decides only mechanics delegated by requirements §9 and the
architecture boundary §7. It emits exactly one approved stop condition and halts
if a semantic choice were found missing; §10.3 records that none was
encountered.

Fail-closed stop conditions available to this stage:
`ARCHITECTURE_DECISION_REQUIRED`, `REQUIREMENT_CLARIFICATION_REQUIRED`,
`BOUNDARY_CHANGE_REQUIRED`, `DEPENDENCY_APPROVAL_REQUIRED`,
`SECURITY_REVIEW_REQUIRED`.

---

## Architecture

This section points to the architectural content already established in §2
(Resolved design decisions) and §3 (Components and repository paths) below.
DD-01 through DD-10 in §2 fix the package placement, sole-committer
discipline, synthetic-observer boundary, synchronous fenced-commit
qualification flow, idempotency ownership, injected-clock scheduling,
maintenance trigger intake, transport-guard placement, request-pipeline
ordering, and envelope-reuse strategy that together constitute this feature's
architecture. §3 enumerates the concrete components and their repository
paths. No architectural content is introduced here beyond what §2/§3 already
state.

---

## 2. Resolved design decisions

Each decision below is a mechanic that realizes an already-fixed observable
behavior. No decision resolves a product semantic; all product semantics are
supplied by the authorities in §1.2.

### DD-01 — Domain service package placement

FEATURE-0015 realized its canonical cloud-model domain service in a single
cohesive package `internal/cloudmodel/` (`store.go`, `lifecycle.go`,
`idempotency.go`, `scheduler.go`, `grants.go`, `audit.go`, `model/`,
`validate/`). Following that verified live convention, FEATURE-0016 places its
domain service in a new sibling package `internal/executiontarget/`. This keeps
the ExecutionTarget lifecycle, observer, qualification, expiry/maintenance
scheduling, idempotency reservation table, projection, and audit assembly in one
small, testable package, and keeps FEATURE-0015's package unmodified (satisfying
the read-only prior-feature-authority constraint). Rationale: minimal
dependencies, explicit structs, and clear ownership boundaries per the
engineering steering; it also physically separates the `execution.sovrunn.io`
group from the FEATURE-0015 `governance.sovrunn.io` group.

### DD-02 — Sole committer is a single in-process service type

`ExecutionTargetLifecycleService` (VS0-WRITER-006) is the only type that mutates
persisted ExecutionTarget state and the internal `NormalizedTargetFactSet` /
`TargetQualificationResult` records, owns the target-bound in-process
current-Maintenance marker, and computes the response-only `effectiveAvailability`
projection. All create/qualify/retire/maintenance/expiry mutations funnel through
this type under one service-owned `sync.Mutex` guarding all F0016 targets,
indexes, internal records, Maintenance markers, and reservations. There are no
nested F0016 locks or per-target locks. Observation and profile evaluation run
with this mutex unlocked. Commit reacquires it, rechecks current state/fences,
appends the required AuditEvent, publishes atomically, and then releases waiters
after unlocking. Clients and HTTP handlers never write status.

### DD-03 — Synthetic observer is a target-bound fixture reader

`sovrunn.synthetic-iaas-observer/v1` is implemented as an external-effect-free,
clock-driven component that reads a target-bound in-memory fixture and returns
either four named facts (each Supported/Unsupported/Unknown) or a classified
deterministic outcome (missing fixture; fixture-declared logical timeout). It
performs no network I/O and never blocks on a real call; the "logical timeout"
is a fixture-declared outcome resolved against the injected clock, not an elapsed
wait. It proposes facts only; it never commits.

### DD-04 — Qualification is synchronous with a captured-fence reservation

Qualify is executed synchronously on the request goroutine. The lifecycle
service opens a `Qualifying` in-flight reservation (never persisted, never
projected) that captures the fences (referenced InfrastructureStack generation,
maintenance epoch, viability fingerprint) at reservation time, invokes the
observer, and evaluates the closed profile without the service mutex. It then
reacquires that mutex.

At commit, the Maintenance-vs-epoch-stale race is resolved by this exact
ordered predicate (ADH-2026-060):

1. First, if an active current-Maintenance marker exists, abort qualification
   with `CONFLICT`/409 + `VS0_TARGET_MAINTENANCE`.
2. Otherwise, if the captured maintenance epoch differs from the current
   maintenance epoch, abort with `STALE_RESOURCE_VERSION`/412 +
   `VS0_TARGET_EPOCH_STALE`.
3. A changed referenced InfrastructureStack generation remains the independent
   `STALE_RESOURCE_VERSION`/412 + `VS0_TARGET_EPOCH_STALE` case.

These three outcomes are mutually exclusive: at most one applies to a given
commit attempt. In every one of these aborts, no qualification conclusion, no
qualification AuditEvent, and no idempotency completion is published; if a
Maintenance transition independently won the race, that transition's own
AuditEvent remains audited regardless of the qualification abort.

Only after this ordered decision, every other captured-fence recheck, and the
required AuditEvent append succeeds does the service atomically replace both
current record refs, refresh `observedGeneration`, advance
resourceVersion/ETag, mark idempotency complete, and publish. It then unlocks
and wakes same-key waiters. Mutex reacquisition, observer fault handling, and
panic/shutdown/cancellation aborts remain as already specified: any of these,
or any stale fence, aborts the reservation before publication, removes it, and
wakes waiters without a completion.

### DD-05 — Idempotency reservation table is F0016-owned and in-process

FEATURE-0016 reuses the FEATURE-0012/0015 idempotency conventions (SHA-256 of
recursively sorted allowed create JSON, or SHA-256 of zero bytes for actions;
`If-Match`, `Accept`, authentication, and request/correlation IDs excluded) but
owns its own in-process reservation table keyed by the exact namespace
`(principal UID, route pattern, derived CloudProvider UID, concrete
action-target UID when applicable, Idempotency-Key)`. The table holds `InFlight`
and `Completed` entries; every InFlight reservation owns one unbuffered `done`
channel. The lifecycle service mutex protects reservation ownership and is the
only authority that closes a `done` channel, exactly once, on completion or
abort. A waiter selects its reservation `done` channel or its request context;
cancellation only detaches that waiter. After a wake-up, a waiter repeats the
current authorization, safe-access, validation, and reservation sequence rather
than receiving an unpublished owner outcome. Owner cancellation, panic,
stale-fence failure, audit failure, and shutdown remove the reservation and
wake waiters without a completion. Retention is 24 hours and 10,000 completed
entries; eviction is at-expiry then
earliest-expiry/lexical-namespace at capacity; `InFlight` never evicts. This is
F0016-owned state cleared on restart; it does not touch FEATURE-0015's table.

### DD-06 — Injected clock drives expiry retry and fact freshness

A single injected clock (a `Clock` interface with a real monotonic
implementation for production and a deterministic fake for tests) drives fact
`expiresAt` freshness checks and the once-per-clock-second expiry-audit retry
loop. Expiry runs as the sole background worker with no request caller; it never
returns 500 and writes exactly one expiry AuditEvent per expiry event once its
append succeeds.

### DD-07 — Maintenance enters/clears only through a target-bound fixture trigger

Maintenance enter/clear is delivered exclusively through a deterministic,
target-bound synthetic-observer fixture trigger carrying target UID, expected
InfrastructureStack generation, expected maintenance epoch, operation, system
identity, and correlation. It is not an HTTP route and not a generic event bus.
The lifecycle service validates the fence (matching current epoch/generation)
before applying; stale-fenced triggers produce no state change, AuditEvent, or
record publication.

### DD-08 — Pre-ServeMux transport guard wraps the F0016 mux registrations

Following the FEATURE-0015 live convention of centralized route registration in
`internal/server/routes.go`, FEATURE-0016 registers exactly five Go 1.22
`http.ServeMux` method/path patterns and wraps them with a fixed pre-ServeMux
method/path guard (an `http.Handler` decorator installed ahead of the mux). The
guard supplies transport-only 405/404 outcomes for HEAD, PUT, other unmatched
methods, and every trailing-slash variant, each with
`X-Content-Type-Options: nosniff` and no Problem body, authentication,
authorization, audit, idempotency, or lifecycle effect. Valid registered pairs
bypass the guard and enter the handler at authentication.

### DD-09 — Request pipeline is an explicit ordered sequence

The handler realizes the fixed precedence from REQ-F16-05 as an explicit,
short-circuiting sequence of stages, reusing FEATURE-0012 pipeline stages
(`internal/apivalid`) for decode/classification/structural/reference validation
and FEATURE-0015 conventions (`internal/cloudmodel` grants + safe resolver) for
authentication, authorization, and safe access. The order is fixed and testable;
see §5.

### DD-10 — Reuse inherited envelopes, errors, and helpers by reference

Problem Details, ObjectMeta, TypedRef, Condition, media handling, ETag/If-Match
handling, and the AuditEvent append boundary are reused from their canonical
owners (`internal/apiproblem`, `internal/apimeta`, `internal/apiref`,
`internal/apicond`, `internal/apivalid`, `internal/decision`). FEATURE-0016
defines no parallel envelope and adds only the eight new violation codes to the
violation surface (top-level Problem codes remain FEATURE-0012-owned).

---

## Components and Interfaces

This section points to §3 (Components and repository paths) and §4 (Data/API
representation) below, which are the authoritative statement of this
feature's components, their responsibilities, their repository paths, and
their HTTP interface. §3.1 lists the new F0016-owned components; §3.2 lists
reused prior-feature components consumed by reference only; §3.3 lists
excluded components. §4.5 fixes the five-route HTTP surface. No component or
interface is introduced here beyond what §3/§4 already state.

---

## 3. Components and repository paths

All paths below follow verified live conventions (FEATURE-0015 domain package
`internal/cloudmodel/`; handlers `internal/api/*_collection.go`,
`*_item.go`, `*_actions.go`; mux `internal/server/routes.go`; conformance
`tests/conformance/`). Concrete file decomposition within a package is a
tasks-stage mechanic; the component responsibilities below are authoritative.

### 3.1 New F0016-owned components (IMPLEMENT)

| Component | Proposed path | Responsibility |
|-----------|---------------|----------------|
| ExecutionTarget domain types | `internal/executiontarget/model/` | In-memory structs for ExecutionTarget, NormalizedTargetFactSet, TargetQualificationResult, the current-Maintenance marker, and the viability fingerprint; enum consts for lifecycle/qualification/effectiveAvailability/fact truth values |
| In-memory store | `internal/executiontarget/store.go` | Process-local registry; scope-unique name reservation (retained after retirement); at-most-one live `(participation UID, stack UID, targetClass)` tuple index; ascending-`metadata.uid` LIST ordering |
| Lifecycle service (VS0-WRITER-006) | `internal/executiontarget/lifecycle.go` | Sole committer of target state + internal records; Qualifying reservation; atomic conclusion publication; maintenance enter/clear; retirement; current-Maintenance marker; effectiveAvailability computation |
| Synthetic observer | `internal/executiontarget/observer.go` | `sovrunn.synthetic-iaas-observer/v1`; target-bound, clock-driven, external-effect-free fact proposal and classified missing-fixture/logical-timeout outcomes; fixed provenance |
| Qualification engine | `internal/executiontarget/qualify.go` | Closed profile evaluation: any Unsupported → Rejected; else any Unknown → Indeterminate; else Qualified; observer-fault detection (duplicate/missing/malformed) |
| Idempotency reservation table | `internal/executiontarget/idempotency.go` | F0016-owned namespace/digest/reservation/waiter table; retention/eviction; waiter recheck of current authorization + safe access |
| Expiry/maintenance scheduler | `internal/executiontarget/scheduler.go` | Injected-clock expiry worker with once-per-clock-second audit retry; target-bound fixture maintenance-trigger intake; shutdown ordering |
| Safe projection | `internal/executiontarget/projection.go` | Closed safe ExecutionTarget JSON projection incl. response-only effectiveAvailability; omits facts/results/observer/handles/links |
| Audit assembly | `internal/executiontarget/audit.go` | Builds redacted FEATURE-0013 AuditEvents for the F0016 audit matrix; delegates append to the inherited FEATURE-0013 append boundary |
| HTTP handlers | `internal/api/executiontarget_collection.go`, `internal/api/executiontarget_item.go`, `internal/api/executiontarget_actions.go` | Collection POST/GET, item GET, qualify POST, retire POST; pipeline wiring per §5 |
| Route registration + transport guard | `internal/server/routes.go` (five registrations) + guard helper (e.g. `internal/server/executiontarget_guard.go`) | Five exact ServeMux patterns and the fixed pre-ServeMux method/path guard |

### 3.2 Reused prior-feature components (CONTRACT_ONLY/NO_TASK — consumed by reference)

| Component | Path | Reuse |
|-----------|------|-------|
| Problem Details + HTTP/URN map | `internal/apiproblem` | Top-level Problem codes and mappings; F0016 adds no top-level code |
| ObjectMeta / TypedRef | `internal/apimeta`, `internal/apiref` | Identity/version fields and typed-reference shape validation |
| Validation pipeline | `internal/apivalid` | Decode, strict classification, structural/semantic/reference stages, concurrency (`If-Match`) helpers |
| Condition | `internal/apicond` | Inherited Condition type (not projected by F0016) |
| Authorization / safe resolver / idempotency conventions | `internal/cloudmodel` (grants, safe denial, idempotency digest) | Server-resolved grants, safe-denial semantics, digest algorithm conventions |
| AuditEvent append boundary | `internal/decision` | Atomic append-before-publication; F0016 does not redefine |

### 3.3 Excluded components (EXCLUDED — never designed here)

Real adapters, provider-native types, credentials/SecretRef, external calls,
external persistence, controllers beyond the lifecycle service, placement,
provisioning, plugin execution, customer target projection, maintenance-notice
resources, IAM resources, the concepts named in the non-goals list in §10.1, and
all FEATURE-0017/0019/0022/0023/0024 concepts (§10.1).

---

## Data Models

This section points to §4 (Data/API representation) below, which is the
authoritative statement of this feature's data models: §4.1 fixes the
ExecutionTarget (VS0-SCHEMA-015) create/safe-projection shape, §4.2 fixes the
internal-only NormalizedTargetFactSet (VS0-SCHEMA-016), §4.3 fixes the
internal-only TargetQualificationResult (VS0-SCHEMA-017), §4.4 fixes provenance
constants, and §4.6 fixes the response-only effectiveAvailability projection.
No data model is introduced here beyond what §4 already states.

---

## 4. Data/API representation

### 4.1 ExecutionTarget (VS0-SCHEMA-015)

Create body accepts exactly (closed set): `metadata.name`,
`spec.cloudProviderParticipationRef.uid`, `spec.infrastructureStackRef.uid`, and
immutable `spec.targetClass="synthetic-iaas"`. Server-assigned only:
`metadata.uid`, `metadata.resourceVersion`, `metadata.generation`,
`status.lifecycle`, `status.qualification`, `status.maintenanceEpoch`,
`status.observedGeneration`, `status.factSetRef`, `status.qualificationResultRef`.
Scope derives solely from `spec.cloudProviderParticipationRef`; no `scopeRef` is
persisted or projected.

Safe projection (create 201 / item GET 200 / qualify 200 / retire 200; LIST
ascending-UID array, no ETag) is exactly: `metadata.uid`, `metadata.name`,
`metadata.generation`, `metadata.resourceVersion`, the three immutable spec
fields, `status.lifecycle`, `status.qualification`, and response-only
`effectiveAvailability`. No public `conditions` member; facts, results, observer
data, handles, and internal record links are omitted. No `Location` header.

### 4.2 NormalizedTargetFactSet (VS0-SCHEMA-016) — internal-only, never projected

Fields: target ref, referenced InfrastructureStack generation, maintenance
epoch, viability fingerprint, observer ID/revision, `factVersion=v1` (the sole
persisted FactSet version field), `observedAt`/`expiresAt`, and the four named
truth values for `compute.vm`, `storage.block`, `storage.object`,
`network.private`, each exactly one of Supported/Unsupported/Unknown.
`factSchemaVersion=v1` is fixed observer provenance describing the emitted fact
schema and is not a second persisted FactSet field.

### 4.3 TargetQualificationResult (VS0-SCHEMA-017) — internal-only, never projected

Fields: target ref and fact-set ref, only the registered
`infrastructureStackGeneration` and `maintenanceEpoch` fences (the FactSet is
referenced by ref; `viabilityFingerprint` is FactSet-only and is not duplicated
on Result), profile version, outcome (Qualified/Rejected/Indeterminate), reason
codes, and evaluation time. Each record is immutable once created, but a target
retains links only to its current FactSet/Result pair: successful qualification
atomically replaces both links, while retirement clears both links.

### 4.4 Fixed provenance and observer identity

`observerID=sovrunn.synthetic-iaas-observer/v1`,
`profileVersion=synthetic-iaas/v1`, `factSchemaVersion=v1`, and the immutable
selected fixture revision.

### 4.5 HTTP surface (five routes; ADH-2026-058 clause 3)

| Method | Path | Handler | ETag | Idempotency |
|--------|------|---------|------|-------------|
| POST | `/apis/execution.sovrunn.io/v1alpha1/execution-targets` | collection create | yes (201) | yes |
| GET | `/apis/execution.sovrunn.io/v1alpha1/execution-targets` | collection LIST | no | never |
| GET | `/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}` | item GET | yes (200) | never |
| POST | `/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify` | qualify | yes (200) | yes |
| POST | `/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire` | retire | yes (200) | yes |

`{uid}` is a complete Go 1.22 `http.ServeMux` wildcard read with
`r.PathValue("uid")`. Success media `application/json`; error media
`application/problem+json`; every response (including transport-only 404/405 and
replay) sends `X-Content-Type-Options: nosniff`. `Accept` is ignored; no route
accepts a query parameter.

### 4.6 effectiveAvailability projection (response-only; no persisted writer; no ETag effect)

Ordered truth table (from committed state, freshness, maintenance, viability):

1. Retired → `Unavailable`.
2. else Active with current Maintenance marker active → `Maintenance`.
3. else only Active + Qualified + fresh current facts + viable backing → `Available`.
4. every other Active combination → `Unavailable`.

A viability-only change does not mutate the target or change its ETag
(VS0-CF-F16-119,120).

---

## Correctness Properties

This section points to the correctness content already established in §9.2
(Risk-relevant invariants and enforcement) and §7.4 (Formal proof references)
below, and to the fence/commit rules in DD-04 and §5.1 step 12. §9.2 states
the invariants that must hold across all valid executions (sole-committer
status writes, no external effect, audit-before-publication, safe-denial
non-disclosure, the closed five-route surface, and stale-fence non-publication)
together with their design enforcement point and conformance IDs. §7.4 points
to the architecture-owned formal models (`ExecutionTargetLifecycle.tla`,
`ExecutionTargetFence.tla`) that these invariants are proved against. No new
correctness property is introduced here beyond what §9.2/§7.4 already state.

Each property heading below restates, verbatim in substance, one row of the
§9.2 invariants table and adds no new invariant.

### Property 1: Sole-committer status writes

For all requests and all ExecutionTarget/NormalizedTargetFactSet/
TargetQualificationResult records, only `ExecutionTargetLifecycleService`
mutates status/internal records; no client or handler write occurs (DD-02;
VS0-CF-F16-03,04,91,92).

**Validates: Requirements 16.08**

### Property 2: No external effect

For all observer invocations, no network I/O, credential use, or external
persistence occurs (DD-03; §8.3; VS0-CF-F16-46).

**Validates: Requirements 16.11**

### Property 3: Audit-before-publication

For all committing mutations, the required AuditEvent append succeeds before
the mutation is published; a failed append leaves target, ETag, and record
links unchanged (DD-04, §5.3; VS0-CF-F16-39,40,82,99..104).

**Validates: Requirements 16.09**

### Property 4: Safe-denial non-disclosure

For all requests against an inaccessible target or backing reference, the
response never confirms existence (§5.1 step 6, §6.2; VS0-CF-F16-09,17,19,66).

**Validates: Requirements 16.02**

### Property 5: Closed route surface with side-effect-free guard

For all requests, exactly five registered routes exist and the pre-ServeMux
transport guard produces no authentication, authorization, audit, idempotency,
or lifecycle effect (DD-08; VS0-CF-F16-44,113..118).

**Validates: Requirements 16.03**

### Property 6: Stale-fence non-publication

For all qualification/maintenance commits, a stale captured fence at commit
time aborts with no target mutation, AuditEvent, or completion (DD-04, §5.1
step 12; VS0-CF-F16-29,30,75,81,82).

**Validates: Requirements 16.08**

---

## Error Handling

This section points to §5 (Validation and deterministic error behavior)
immediately below, together with §5.2 (closed local violation set), §5.3
(deterministic single-writer commit), and §6.3 (Observability — audit matrix
and append-failure behavior). §5 fixes the exact ordered precedence by which
every request either proceeds or fails closed with a specific Problem status
and, where registered, a specific violation code; §5.2 closes the local
violation set at exactly eight new codes; §5.3 fixes the all-or-nothing commit
discipline; §6.3 fixes which mutations are audited and the no-disclosure
behavior on a required audit-append failure. No error-handling behavior is
introduced here beyond what §5/§6.3 already state.

---

## 5. Validation and deterministic error behavior

### 5.1 Request precedence (REQ-F16-05; ADH clause 5) — fixed order

1. Closed transport method/path guard (pre-ServeMux; DD-08).
2. Authentication → `AUTH_REQUIRED`/401 on missing/invalid, no audit, no record.
3. Method/media/header-form validation (application/json for POST; exactly one
   Idempotency-Key; qualify/retire: exactly one strong opaque If-Match + zero-byte
   body; query parameter or If-Match on collection create → `MALFORMED_REQUEST`/400;
   non-JSON media → `UNSUPPORTED_MEDIA_TYPE`/415; malformed UID segment →
   `MALFORMED_REQUEST`/400 before safe resolution). Header-form failure precedes
   body classification.
4. Phase-one safe references: read only the path UID needed to derive
   CloudProvider scope.
5. Authorization: create/retire `executiontarget.write`; reads
   `executiontarget.read`; qualify `executiontarget.qualify`.
6. Safe access: inaccessible target/backing → safe `RESOURCE_NOT_FOUND`/404 +
   `VS0_AUTHORIZATION_SAFE_DENIAL` (audited); never discloses existence.
7. Strict classification: unknown field / `adapterAuthorityRef` / SecretRef →
   `UNKNOWN_FIELD`/400; client status write → `AUTHORIZATION_DENIED`/403 +
   `VS0_STATUS_FIELD_WRITE`; client identity/metadata write →
   `AUTHORIZATION_DENIED`/403 + `VS0_SYSTEM_OWNED_FIELD_WRITE`; malformed/duplicate/
   oversized body → inherited `MALFORMED_REQUEST`/`DUPLICATE_FIELD`/`REQUEST_TOO_LARGE`/400.
8. Graph/reference validation: missing required field / invalid typed-ref shape /
   non-`synthetic-iaas` targetClass → `VALIDATION_FAILED`/422; ineffective
   participation → `CONFLICT`/409 + `VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE`;
   non-Active stack → `CONFLICT`/409 + `VS0_EXECUTION_TARGET_STACK_UNAVAILABLE`;
   cross-CloudProvider graph → `VALIDATION_FAILED`/422 +
   `VS0_EXECUTION_TARGET_SCOPE_MISMATCH`; duplicate name or live duplicate tuple →
   `ALREADY_EXISTS`/409.
9. Replay/reservation (create/qualify/retire only): completed replay →
   original successful status/body + fresh correlation, no second audit;
   changed digest → `CONFLICT`/409 + `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`;
   independent reservations across action targets and derived scopes.
10. Current If-Match version comparison (actions): syntactically valid stale →
    `STALE_RESOURCE_VERSION`/412 (after replay/reservation, before Maintenance
    state evaluation).
11. Lifecycle state: Retired → `CONFLICT`/409 + `VS0_TARGET_RETIRED`;
    Maintenance → `CONFLICT`/409 + `VS0_TARGET_MAINTENANCE`; different-key qualify
    while Qualifying → `CONFLICT`/409 + `VS0_TARGET_QUALIFICATION_IN_PROGRESS`.
12. Atomic commit: fence recheck at commit — active current-Maintenance marker
    (checked first) → `CONFLICT`/409 + `VS0_TARGET_MAINTENANCE`; otherwise a
    changed captured maintenance epoch → `STALE_RESOURCE_VERSION`/412 +
    `VS0_TARGET_EPOCH_STALE`; a changed InfrastructureStack generation remains
    the independent `STALE_RESOURCE_VERSION`/412 + `VS0_TARGET_EPOCH_STALE`
    case; these three outcomes are mutually exclusive per commit attempt
    (ADH-2026-060); stale viability →
    `STALE_RESOURCE_VERSION`/412 + `VS0_EXECUTION_TARGET_VIABILITY_STALE`;
    observer fault (duplicate/missing/malformed) →
    `INTERNAL_ERROR`/500 with reservation abort and no publication; required
    AuditEvent append failure → `INTERNAL_ERROR`/500 with no AuditEvent, mutation,
    or completion and no suppressed-outcome disclosure.

### 5.1.1 Five-route stage matrix

The fixed sequence above applies only where a stage is applicable. This matrix
does not add a route or outcome; it makes the five existing route contracts
executable. A waiter that wakes after an owner abort re-enters at authentication
and repeats every applicable stage before it can establish a new reservation.

| Route | Applicable stages in order | Skipped stages | Exact local conformance |
|-------|----------------------------|----------------|-------------------------|
| Collection create | transport guard → authentication → POST media/key grammar → phase-one extraction of only `spec.cloudProviderParticipationRef.uid` and `spec.infrastructureStackRef.uid` from the closed envelope → derived-scope authorization → safe backing access → strict classification → graph/reference and uniqueness validation → reservation/replay → audit-before-publication commit | path-UID resolution; action If-Match; lifecycle action state | VS0-CF-F16-01..13,46..50,54..60,84,91..93,98,105,112 |
| Collection LIST | transport guard → authentication → read authorization → accessible-target enumeration and ascending-UID projection | body decode/classification; backing graph validation; safe item resolution; idempotency; If-Match; lifecycle/fence commit | VS0-CF-F16-44,63,94,105..110,113..118 |
| Item GET | transport guard → authentication → phase-one path-UID safe resolution → derived-scope read authorization → safe access → safe projection | body decode/classification; idempotency; If-Match; lifecycle/fence commit | VS0-CF-F16-09,17,19,44,64,66,95,102,105..118 |
| Qualify action | transport guard → authentication → action media/key/If-Match/zero-byte grammar → phase-one path-UID safe resolution → derived-scope qualify authorization → safe access → replay/reservation → current If-Match → lifecycle state → fenced observer/qualification commit and audit-before-publication | client graph body fields; create uniqueness | VS0-CF-F16-05,06,14,16,20..35,58..60,67..75,78,81..83,85..90,96,99,103,104 |
| Retire action | transport guard → authentication → action media/key/If-Match/zero-byte grammar → phase-one path-UID safe resolution → derived-scope write authorization → safe access → replay/reservation → current If-Match → lifecycle state → audit-before-publication retirement commit | client graph body fields; observer/profile evaluation; qualification fences | VS0-CF-F16-05,06,14,16,24..26,31..35,37,39,40,58..60,70..74,77,78,82,85..89,97,99 |

### 5.2 Closed local violation set (REQ-F16-10)

Exactly eight new codes, appearing only in `violations[]`:
`VS0_TARGET_RETIRED`, `VS0_TARGET_MAINTENANCE`,
`VS0_TARGET_QUALIFICATION_IN_PROGRESS`, `VS0_TARGET_EPOCH_STALE`,
`VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE`,
`VS0_EXECUTION_TARGET_STACK_UNAVAILABLE`,
`VS0_EXECUTION_TARGET_SCOPE_MISMATCH`,
`VS0_EXECUTION_TARGET_VIABILITY_STALE`. Top-level Problem codes remain
FEATURE-0012-owned; HTTP status is transport metadata, never a Problem `code`;
no `429`/quota code exists. `VS0_STATUS_FIELD_WRITE`,
`VS0_SYSTEM_OWNED_FIELD_WRITE`, `VS0_AUTHORIZATION_SAFE_DENIAL`, and
`VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH` are inherited violation codes reused, not
introduced.

### 5.3 Deterministic single-writer commit

All commits serialize through the lifecycle service (DD-02). The registered
writer conflict for VS0-WRITER-006 is `STALE_RESOURCE_VERSION`. The commit step
is all-or-nothing: audit append precedes publication; a failed append leaves the
committed target, ETag, and record links unchanged.

---

## 6. Security, privacy, observability, compatibility, and versioning

### 6.1 Security

- Authentication precedes all F0016 effects; missing/invalid → 401 with no audit
  or record (VS0-CF-F16-07,61..64,93..97).
- Grants: `executiontarget.write` (create/retire), `executiontarget.read`
  (LIST/GET), `executiontarget.qualify` (qualify), each bound to the derived
  CloudProvider scope of the target UID for item operations.
- Client status/identity writes are audited 403 with the inherited
  `VS0_STATUS_FIELD_WRITE` / `VS0_SYSTEM_OWNED_FIELD_WRITE` violations; clients
  never write status (VS0-WRITER-005/006).
- Network-exposed routes are authenticated and authorized by design; there is no
  unauthenticated mutating path. The pre-ServeMux guard is transport-only and
  performs no authentication decision, so it never exposes resource data.

### 6.2 Privacy / redaction

- ExecutionTarget is Provider-confidential; NormalizedTargetFactSet and
  TargetQualificationResult are internal-only and never projected.
- Safe denial (404 + `VS0_AUTHORIZATION_SAFE_DENIAL`) never confirms existence of
  an inaccessible target or backing reference.
- No secret value on any path; no SecretRef extension is active; no credentials
  exist. Replay discloses only the original successful status/body with fresh
  correlation; denial replays disclose no stored result.

### 6.3 Observability

- Reuse the inherited FEATURE-0013 AuditEvent as the durable audit signal. The
  audit matrix audits: successful create, completed qualification, retirement,
  maintenance enter/clear, authenticated authorization denial, and safe denial.
  A required request-triggered append failure → `INTERNAL_ERROR`/500 with no
  AuditEvent, mutation, or completion (VS0-CF-F16-39,40,82,99..104,112).
- Expiry is the sole background exception: no request caller, no 500, one expiry
  AuditEvent once its append succeeds, retrying once per injected-clock second
  (VS0-CF-F16-41,79).
- Structured request/correlation fields follow the FEATURE-0015 handler
  conventions and the observability-and-audit baseline; no secret is logged.

### 6.4 Compatibility and versioning

- Reuses FEATURE-0012 Problem Details, ObjectMeta, TypedRef, Condition, media,
  and optimistic-concurrency conventions by reference; defines no parallel
  envelope. API group/version is `execution.sovrunn.io/v1alpha1`.
- Consumes FEATURE-0015 `CloudProviderParticipation` and `InfrastructureStack`
  read-only; modifies no FEATURE-0015 document, schema, writer, state, or
  conformance case.
- ADH-2026-058 atomically replaced the placeholder VS0-SCHEMA-015..017 and
  VS0-STATE-004; there is no alias, migration path, or legacy compatibility for
  `status.availability`, `Draining`, or publicly projected `Qualifying`. No
  migration impact: all state is new and in-process.

### 6.5 Operational

- All state is in-process. Restart clears F0016-owned targets, records,
  idempotency reservations, and pending work; it never deletes inherited
  FEATURE-0013 AuditEvents and makes no external observer call (VS0-CF-F16-46).
- Shutdown ordering: stop expiry/Maintenance trigger acceptance → abort in-flight
  qualification and idempotency reservations and wake waiters → stop HTTP serving
  (VS0-CF-F16-77).

---

## Testing Strategy

This section points to §7 (Test and conformance strategy) below, which is the
authoritative statement of this feature's testing approach: §7.1 fixes test
placement and harness reuse; §7.2 fixes the obligation to realize every one of
the 122 FEATURE-0016-local conformance cases (VS0-CF-F16-01..122); §7.3 groups
those cases by category; §7.4 points to the architecture-owned formal proof
references consumed as design constraints; §7.5 lists the tasks-stage
verification commands. No testing approach is introduced here beyond what §7
already states.

---

## 7. Test and conformance strategy

### 7.1 Placement and reuse

Feature-local conformance cases follow the live convention in
`tests/conformance/` (e.g. `feature_0016_test.go`), reusing the FEATURE-0015 test
harness/fixtures patterns. Unit and property tests live beside the implementation
in `internal/executiontarget/` and `internal/api/`, and route/guard tests in
`internal/server/`, matching FEATURE-0015's `*_test.go` co-location.

### 7.2 Conformance coverage obligation

Every one of the 122 FEATURE-0016-local cases (VS0-CF-F16-01..122) is realized as
a test asserting its exact registered `expectedState`, `expectedError`,
`expectedViolation` (where registered), and `expectedSideEffects`. Equivalence
families (VS0-CF-F16-106,112,114,115,118) must execute every named member.
Determinism is guaranteed by the injected clock and target-bound fixtures; no
test relies on wall-clock elapse or network. No shared/downstream case
(`VS0-CF-F10`, `VS0-CF-X03`, `VS0-CF-HP01`) is counted as F0016-local proof.

### 7.3 Test categories

- Create-field boundary and scope derivation (01,02,04,54..57,84,91,92).
- Safe access before graph validation (09..13,17,19,65,66,98).
- Route/guard transport outcomes (44,45,48,105..110,113..118) and precedence
  (05,06,07,08,24,25,58..60,70..74,85,86).
- Idempotency replay/reuse/isolation/eviction/panic/waiter (14,16,31..35,46,47,
  49,50,76..78,81,82,87..89,99).
- Observer/facts/conclusions (20..23,67..69,90) and lifecycle/availability
  (26..30,36,38,39,41..43,75,79..83,88,89,100,101,119..122).
- Audit matrix and append-failure (39,40,82,99..104,112) and future boundary (46).

### 7.4 Formal proof references (owned by architecture; referenced only)

`docs/formal/feature-0016/ExecutionTargetLifecycle.tla` (Retired terminality,
availability truth table, Maintenance/Retirement priority, non-public Qualifying)
and `docs/formal/feature-0016/ExecutionTargetFence.tla` (audit-before-conclusion,
stale/cancelled/failed non-mutation) are consumed as design constraints, not
re-authored.

### 7.5 Verification commands (tasks-stage execution)

`make fmt`, `make test`, `make vet`, `go test -race ./...`,
`make feature-contract-check FEATURE=FEATURE-0016`,
`make vs000-contract-check`, and `make phase2r-drift-check` are component
checks. The authoritative final gate is
`make ff-feature-gate FEATURE=FEATURE-0016`; it must pass after implementation.
Design authors no runtime; these commands run when implementation exists.

---

## 8. Implementation classification ledger

### 8.1 IMPLEMENT (produces tasks / code later)

| Element | Component (§3) | Realizes |
|---------|----------------|----------|
| ExecutionTarget in-memory type + store + tuple/name index | model/, store.go | Create-time persistence and scope-derived, name/tuple-unique storage of ExecutionTarget |
| ExecutionTargetLifecycleService (sole committer) | lifecycle.go | Sole-committer lifecycle state management and expiry/maintenance scheduling |
| Qualifying reservation + captured-fence commit | lifecycle.go, qualify.go | Sole-committer atomic, fence-checked qualification commit |
| Synthetic observer + fixtures | observer.go | Synthetic fact observation feeding qualification |
| Qualification profile evaluation | qualify.go | Synthetic fact observation feeding qualification |
| NormalizedTargetFactSet + TargetQualificationResult records | model/ | Synthetic fact observation and the sole-committer qualification commit it feeds |
| effectiveAvailability projection | projection.go, lifecycle.go | Sole-committer lifecycle state management and its response-only availability projection |
| Idempotency reservation/waiter table | idempotency.go | F0016-owned idempotency handling for create/qualify/retire |
| Injected-clock expiry worker + retry | scheduler.go | Expiry/maintenance scheduling |
| Target-bound maintenance-trigger intake + current-Maintenance marker | scheduler.go, lifecycle.go | Expiry/maintenance scheduling and sole-committer lifecycle state management |
| Shutdown ordering | scheduler.go, lifecycle.go | Expiry/maintenance scheduling shutdown behavior |
| Five handlers + request pipeline | `internal/api/executiontarget_*.go` | The closed HTTP surface and its fixed, ordered request pipeline |
| Five ServeMux registrations + pre-ServeMux guard | `internal/server/*` | The closed HTTP surface |
| Safe projection + closed representation | projection.go | The closed HTTP surface's response representation |
| Audit assembly (redacted, matrix) | audit.go | Sole-committer lifecycle state management and expiry/maintenance scheduling audit effects |
| Eight new violation constants | `internal/executiontarget/violations.go` | F0016-owned `violations[]` constants; `internal/apiproblem` Problem enum/map is reused unchanged |
| F0016-local conformance suite (122 cases) | `tests/conformance/feature_0016_test.go` | Full coverage of the feature's design elements above |

### 8.2 CONTRACT_ONLY/NO_TASK (reused by reference; no new task)

| Element | Owner | Disposition |
|---------|-------|-------------|
| Problem Details + HTTP/URN mapping | FEATURE-0012 / `internal/apiproblem` | Reused; F0016 adds no top-level code |
| ObjectMeta / TypedRef / Condition | FEATURE-0012 / `internal/apimeta`,`apiref`,`apicond` | Reused as-is |
| Validation pipeline stages, media, If-Match/ETag | FEATURE-0012 / `internal/apivalid` | Reused; F0016 wires order, adds no stage semantics |
| AuditEvent atomic append-before-publication | FEATURE-0013 / `internal/decision` | Reused; not redefined |
| Server-resolved authorization, safe denial, digest conventions | FEATURE-0015 / `internal/cloudmodel` | Reused; FEATURE-0015 unchanged |
| Inherited violation codes (`VS0_STATUS_FIELD_WRITE`, `VS0_SYSTEM_OWNED_FIELD_WRITE`, `VS0_AUTHORIZATION_SAFE_DENIAL`, `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`) | prior features | Reused; not introduced |
| Formal models (lifecycle, fence) | architecture / `docs/formal/feature-0016` | Referenced as constraints |

### 8.3 EXCLUDED (must never be designed or built by F0016)

Real adapter/provider-native type/credential/SecretRef/external call/external
persistence; placement, provisioning, plugin execution, realization, plugin
taxonomy, Crossplane dependency; customer target projection; maintenance-notice
resource; IAM resource; the explicit non-goals listed in §10.1; persisted
`status.availability`, `Draining`, public `Qualifying`, public `conditions`,
`scopeRef`, `spec.participationRef`, `adapterAuthorityRef`; FEATURE-0017
PolicyEvaluationRequest/Result/PolicyEngineAdapter; FEATURE-0019
SovereigntyProfile/SovereigntyFactSet/EvidenceRecord; FEATURE-0022
ServiceRegion / `spec.executionTargetRefs` / `ServiceRegion.status.availability`;
FEATURE-0023 DecisionProfiles/ServicePlacement/PlacementDecision/customer target
projection; FEATURE-0024 PluginExecution/ServiceDeploymentPlan/real realization.

---

## 9. Requirement/decision/risk traceability

### 9.1 Architecture Traceability

- Consumed decisions/handoffs: DEC-0036, DEC-0042, DEC-0057; ADH-2026-025,
  ADH-2026-040, ADH-2026-042, ADH-2026-045, ADH-2026-058 (sole current semantic
  authority).
- Schema IDs: VS0-SCHEMA-015 (ExecutionTarget), VS0-SCHEMA-016
  (NormalizedTargetFactSet), VS0-SCHEMA-017 (TargetQualificationResult) —
  referenced, not redefined.
- Writer ID: VS0-WRITER-006 (`ExecutionTargetLifecycleService`; conflict
  `STALE_RESOURCE_VERSION`).
- State ID: VS0-STATE-004 (Active/Retired × Unqualified/Qualified/Rejected/
  Indeterminate; Qualifying is an in-flight reservation only).
- Error codes: top-level Problem codes owned by FEATURE-0012; eight new
  `violations[]` codes per §5.2.
- Conformance: VS0-CF-F16-01..122 (all F0016-local).

### 9.2 Risk-relevant invariants and enforcement

| Invariant | Design enforcement |
|-----------|--------------------|
| Clients never write status; single sole committer | DD-02; handlers never mutate status; VS0-CF-F16-03,04,91,92 |
| No external effect/credential/persistence | DD-03; observer external-effect-free; §8.3; VS0-CF-F16-46 |
| Audit-before-publication; no partial mutation on append failure | DD-04, §5.3; VS0-CF-F16-39,40,82,99..104 |
| Safe denial never discloses inaccessible resources | §5.1 step 6, §6.2; VS0-CF-F16-09,17,19,66 |
| Closed five-route surface; transport guard has no side effect | DD-08; VS0-CF-F16-44,113..118 |
| Stale fences cannot publish a conclusion or replayable completion | DD-04, §5.1 step 12; VS0-CF-F16-29,30,75,81,82 |

### 9.3 Conformance Mapping (behavior → design disposition → conformance IDs)

The canonical coverage ledger (below) is the single source of REQ-F16-*/AC-F16-*
identifier mappings; this table restates the same coverage in plain-language
behavior terms, pointing only at conformance IDs.

| Behavior | Design disposition | Conformance IDs |
|----|--------------------|-----------------|
| Create-time persistence and field acceptance | model/,store.go,lifecycle.go create path (§4.1,§5.1) | VS0-CF-F16-01,02,57,84 |
| Safe access ordering before graph validation | safe resolver before graph validation (§5.1 steps 4/6/8) | VS0-CF-F16-09..13,65,98 |
| Closed five-route HTTP surface and transport guard | five registrations + pre-ServeMux guard (DD-08,§4.5) | VS0-CF-F16-44,61..64,66,93..96,105..118 |
| Header-form request grammar validation | header-form validation before body classification (§5.1 step 3) | VS0-CF-F16-05,06,24,25,58..60,70..74,85,86 |
| Idempotency replay and reservation behavior | idempotency table (DD-05,§5.1 step 9) | VS0-CF-F16-14,16,31..35,46,47,49,50,76..78,81,82,86..89,99 |
| Synthetic observer fact behavior | observer + qualification engine (DD-03,DD-04,§4.2) | VS0-CF-F16-20..23,67..69,90 |
| Sole-committer lifecycle and effectiveAvailability projection | lifecycle service + effectiveAvailability (DD-02,§4.6) | VS0-CF-F16-01,20..22,26..30,36,39,41..43,75,77,79..83,87,88,99..101,119..122 |
| Retire/Maintenance precedence over in-flight qualification | retire/maintenance win over Qualifying (§5.1 step 11); the Maintenance-race outcome is mutually exclusive and ordered (DD-04, §5.1 step 12; ADH-2026-060): F16-75 is the inactive-marker changed-maintenance-epoch `STALE_RESOURCE_VERSION`/412 case, F16-89 is the active-Maintenance-entry-wins `CONFLICT`/409 case | VS0-CF-F16-26..28,75,80,83,88,89 |
| Expiry retry and shutdown ordering | expiry retry + shutdown ordering (DD-06,§6.5) | VS0-CF-F16-41,77,79 |
| Audit matrix and append-failure behavior | audit matrix + append-failure (§6.3) | VS0-CF-F16-39,40,82,99..104,112 |
| Closed local violation set | eight closed violations (§5.2) | VS0-CF-F16-10..12,26..30,38,75 |
| Absence of external effect; restart behavior | no external effect; restart clears in-memory (§6.5,§8.3) | VS0-CF-F16-46 |

---

## Canonical coverage ledger

Every approved requirement and every approved acceptance criterion from
`requirements.md` appears below exactly once with its design disposition. No ID
is created, renumbered, merged, split, reinterpreted, or omitted. REQ and AC
meanings are preserved exactly from `requirements.md`.

### Requirements

| ID | Design disposition | Classification |
|----|--------------------|----------------|
| REQ-F16-01 | ExecutionTarget in-memory type, store, closed create-field decode, participation-derived scope, server-assigned identity/status; create persists Active/Unqualified at epoch 0 (DD-01,DD-02; §3.1,§4.1,§5.1) | IMPLEMENT |
| REQ-F16-02 | Safe resolver before graph validation; effective-Active participation + Active stack same-scope check; scope-unique name retained after retirement; at-most-one live tuple index (§3.1 store.go,§5.1 steps 4/6/8) | IMPLEMENT |
| REQ-F16-03 | Five ServeMux registrations + fixed pre-ServeMux transport guard (DD-08; §3.1,§4.5) | IMPLEMENT |
| REQ-F16-04 | Request grammar (media/Idempotency-Key/If-Match/zero-byte), ETag placement, ignored Accept, closed safe projection, status codes (§4.1,§4.5,§5.1 step 3) | IMPLEMENT |
| REQ-F16-05 | Explicit ordered request pipeline with short-circuit and fixed precedence (DD-09; §5.1) | IMPLEMENT |
| REQ-F16-06 | F0016-owned idempotency namespace/digest/reservation/waiter/retention/eviction; GET/LIST never touch idempotency (DD-05; §3.1,§5.1 step 9) | IMPLEMENT |
| REQ-F16-07 | Synthetic observer + closed qualification profile + four named facts + missing-fixture/timeout outcomes + fixed provenance + observer-fault faulting (DD-03,DD-04; §3.1,§4.2,§4.4) | IMPLEMENT |
| REQ-F16-08 | Sole-committer lifecycle service; lifecycle/qualification axes; non-projected Qualifying; response-only effectiveAvailability truth table; atomic audited conclusion; stale-fence abort (DD-02,DD-04; §4.6,§5.3) | IMPLEMENT |
| REQ-F16-09 | Injected-clock expiry with retry; target-bound maintenance trigger; retire sole exit; audit matrix; shutdown ordering (DD-06,DD-07; §3.1,§6.3,§6.5) | IMPLEMENT |
| REQ-F16-10 | Eight new `violations[]` codes only; top-level Problems remain FEATURE-0012 (§5.2) | IMPLEMENT (violation surface); top-level codes CONTRACT_ONLY (reused) |
| REQ-F16-11 | No external effect/adapter/credential/Crossplane; restart clears F0016-owned in-memory state, preserves inherited AuditEvents (§6.5,§8.3) | IMPLEMENT (restart/in-memory boundary); external-effect absence CONTRACT_ONLY |

### Acceptance criteria

| ID | Design disposition | Exact conformance cases | Classification |
|----|--------------------|-------------------------|----------------|
| AC-F16-01 | Create path persists Active/Unqualified at epoch 0 with participation-derived scope; unknown/adapterAuthorityRef/SecretRef rejected (§4.1,§5.1) | VS0-CF-F16-01,02,57,84 | IMPLEMENT |
| AC-F16-02 | Safe access before graph validation with exact ineffective/non-Active/cross-scope/duplicate outcomes (§5.1) | VS0-CF-F16-09..13,65,98 | IMPLEMENT |
| AC-F16-03 | Exactly five routes; guard supplies all HEAD/PUT/trailing-slash/unmatched transport outcomes with no side effect (DD-08) | VS0-CF-F16-44,61..64,66,93..96,105..118 | IMPLEMENT |
| AC-F16-04 | POST media/key and action If-Match/zero-body header-form validation precede body classification (§5.1 step 3) | VS0-CF-F16-05,06,24,25,58..60,70..74,85..86 | IMPLEMENT |
| AC-F16-05 | Replay/reuse/isolation across create/qualify/retire; GET/LIST ignore idempotency (DD-05) | VS0-CF-F16-14,16,31..35,46..47,49..50,76..78,81..82,86..89,99 | IMPLEMENT |
| AC-F16-06 | Four Supported/Unsupported/Unknown facts + deterministic missing-fixture/timeout behavior (DD-03) | VS0-CF-F16-20..23,67..69,90 | IMPLEMENT |
| AC-F16-07 | Sole committer; Qualifying never projected; effectiveAvailability truth table; no persisted status.availability (DD-02,§4.6) | VS0-CF-F16-01,20..22,26..30,36,39,41..43,75,77,79..83,87..88,99..101,119..122 | IMPLEMENT |
| AC-F16-08 | Retire/Maintenance win over in-flight qualification with exact violations; both audited, no qualification AuditEvent (§5.1 step 11) | VS0-CF-F16-26..28,75,80,83,88..89 | IMPLEMENT |
| AC-F16-09 | Expiry retries once per injected-clock second, never returns 500 to a caller; shutdown ordering (DD-06,§6.5) | VS0-CF-F16-41,77,79 | IMPLEMENT |
| AC-F16-10 | Audit matrix; append failures → INTERNAL_ERROR/500 with no disclosure (§6.3) | VS0-CF-F16-39,40,82,99..104,112 | IMPLEMENT |
| AC-F16-11 | Only eight closed local violation codes; no other F0016-local code (§5.2) | VS0-CF-F16-10..12,26..30,38,75 | IMPLEMENT |
| AC-F16-12 | No real adapter/credential/external call/Crossplane/placement/provisioning/plugin/customer IaaS exposure; restart proves absence (§6.5,§8.3) | VS0-CF-F16-46 | IMPLEMENT (restart proof); external-effect absence CONTRACT_ONLY |

---

## 10. Non-goals, absence ledger, and unresolved report

### 10.1 Non-goals and adjacent-feature exclusions

FEATURE-0016 does not design or build any element listed in §8.3. It introduces
no real adapter, provider-native type, raw credential, external call, external
persistence, placement, execution, plugin, customer target API,
maintenance-notice resource, or IAM resource. FEATURE-0017/0019/0022/0023/0024
concepts appear only as enforceable exclusions, never as design mechanisms.

The following are the explicit non-goals for this feature:

- Non-goal: FEATURE-0016 must not introduce, activate, or design `ResourcePool`.
- Non-goal: FEATURE-0016 must not introduce, activate, or design `ProviderCapability`.
- Non-goal: FEATURE-0016 must not introduce, activate, or design a generic `Provider`.
- Non-goal: FEATURE-0016 must not introduce or activate any ExecutionTarget subtype as an active platform concept.

### 10.2 Absence ledger (deliberately absent by contract)

- Persisted `status.availability`; `Draining`; publicly projected `Qualifying`.
- Public `conditions` member on ExecutionTarget.
- `scopeRef` (scope derived solely from participation); `spec.participationRef`
  (renamed to `spec.cloudProviderParticipationRef`).
- `adapterAuthorityRef`; any `SecretRef` extension.
- `viabilityFingerprint` on TargetQualificationResult (FactSet-only); a second
  persisted FactSet version field (only `factVersion=v1` persists).
- `Location` header; conditional GET; `Accept` negotiation; any query parameter.
- PATCH/PUT/DELETE/HEAD public routes; watch/filter/pagination/fact/result/
  adapter-selection routes.
- `429`/quota HTTP code; any top-level Problem code introduced by F0016.
- Any alias, migration path, or legacy compatibility for the above.

### 10.3 Resolved report

ADH-2026-060 has been approved and applied to the manifest-controlled
architecture authorities (registry, architecture doc, traceability, steering,
requirements.md, and the readiness/contract-check/vs000-check scripts). The
Maintenance-race outcome between F16-75 (inactive-marker, changed-maintenance-
epoch, `STALE_RESOURCE_VERSION`/412) and F16-89 (active-marker,
Maintenance-wins, `CONFLICT`/409) is now mutually exclusive and ordered per the
predicate fixed in DD-04 and §5.1 step 12. No unresolved architecture,
requirements, or design decision remains for this feature.

---

Model Execution Report:
- Tool: kiro
- Stage or task: FEATURE-0016 Design (adapter-boundary-and-executiontarget-qualification/design.md)
- Recommended priority list: System design (claude-opus-4.8, effort high); Fallback (claude-sonnet-4.5, effort medium); (third model not specified in prompt — only two labels provided)
- Selected model: claude-opus-4.8
- Effort/reasoning setting: high
- Fallback used: no
- Fallback reason: none

STAGE_STATUS: COMPLETE
