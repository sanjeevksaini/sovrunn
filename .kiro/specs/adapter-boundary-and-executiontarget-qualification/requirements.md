# Requirements Document

## Introduction

This document is the FEATURE-0016 Requirements stage for Adapter Boundary and
ExecutionTarget Qualification. It translates the approved FEATURE-0016 feature
authority, the architecture boundary, and ADH-2026-058 into observable-behavior
requirements, introducing no decision beyond those authorities. Full identity,
stage, purpose, actors, use cases, and terms appear in Sections 1-3 below; the
canonical requirement and acceptance ledgers appear in Section 4.

## Glossary

Canonical terminology for this feature — including `ExecutionTarget`,
`NormalizedTargetFactSet`, `TargetQualificationResult`, `synthetic-iaas`, the
four named facts, qualification outcomes, `effectiveAvailability`, maintenance
epoch, the current-Maintenance marker, viability fingerprint, backing tuple,
the pre-ServeMux transport guard, and the eight new local violation codes — is
defined verbatim in Section 3, "Terms introduced by this feature," below.

## Requirements

The normative requirements and acceptance criteria for this feature are
defined verbatim in Section 4, "Normative requirements and acceptance
scenarios," below. Section 4 contains the canonical REQ-F16-01..11 requirement
ledger, the AC-F16-01..12 acceptance ledger, the normative detail headings for
each requirement (§4.1), the acceptance-scenario coverage mapping (§4.2), and
the exact conformance semantics ledger covering the FEATURE-0016 local conformance suite (§4.3).

# FEATURE-0016 Requirements: Adapter Boundary and ExecutionTarget Qualification

## 1. Identity and stage

| Field | Value |
|-------|-------|
| Feature | FEATURE-0016 — Adapter Boundary and ExecutionTarget Qualification |
| Stage | Requirements |
| Kiro slug | `adapter-boundary-and-executiontarget-qualification` |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase | 2R |
| Order | 6 |
| Status | Approved scope (closed under ADH-2026-058) |
| Owned resources | `ExecutionTarget` (VS0-SCHEMA-015), `NormalizedTargetFactSet` (VS0-SCHEMA-016), `TargetQualificationResult` (VS0-SCHEMA-017) |
| Owned state machine | `VS0-STATE-004` |
| Owned writer | `VS0-WRITER-006` (`ExecutionTargetLifecycleService`) |
| Controlling decisions | DEC-0036, DEC-0042, DEC-0057 |
| Controlling handoffs | ADH-2026-025, ADH-2026-040, ADH-2026-042, ADH-2026-045, ADH-2026-058 |
| Local conformance | VS0-CF-F16-01..99, VS0-CF-F16-100..122 |

This document translates the approved FEATURE-0016 feature authority
(`docs/features/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`),
the architecture boundary
(`docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`),
and ADH-2026-058 into observable-behavior requirements. It introduces no
decision beyond those authorities. It owns observable behavior only: intent,
actors, scenarios, invariants, validation outcomes, security/privacy,
compatibility, non-goals, and acceptance criteria. It does not own packages,
files, routes as code, storage structures, algorithms, libraries, internal
interfaces, or task decomposition.

Stage boundary: this requirements stage modifies only this file. Design,
tasks, architecture, schemas, source, prompts, validators, and automation are
out of scope for this stage.

## 2. Purpose and use cases

### 2.1 Purpose

FEATURE-0016 activates one deterministic, infrastructure-backed
`synthetic-iaas` ExecutionTarget qualification profile. It establishes a
Sovrunn-owned observation-and-qualification boundary with no real external
calls, credentials, placement, provisioning, adapter selection, or customer
projection. It qualifies infrastructure-backed `synthetic-iaas` only; it does
not place, provision, realize, or expose IaaS objects.

### 2.2 Actors

| Actor | Role in FEATURE-0016 |
|-------|----------------------|
| CloudProvider operator principal | Authenticated caller that creates, lists, reads, qualifies, and retires ExecutionTargets within an accessible derived CloudProvider scope, subject to `executiontarget.write`, `executiontarget.read`, and `executiontarget.qualify` grants |
| `ExecutionTargetLifecycleService` | Sole committer (VS0-WRITER-006) of persisted ExecutionTarget status and internal `NormalizedTargetFactSet`/`TargetQualificationResult` records; owner of the in-process current-Maintenance marker and the effectiveAvailability projection |
| `sovrunn.synthetic-iaas-observer/v1` | The sole observer: target-bound, clock-driven, external-effect-free; proposes four named facts and never writes directly |
| Injected clock | Deterministic time source that drives fact expiry and retry timing |
| Synthetic-observer fixture maintenance trigger | Deterministic, target-bound trigger that is the only path for maintenance enter/clear (never an HTTP route or generic event bus) |

### 2.3 Primary use cases

1. Register an ExecutionTarget over an effective-Active FEATURE-0015
   `CloudProviderParticipation` and an Active `InfrastructureStack` in the same
   CloudProvider scope, persisting `Active/Unqualified` at maintenance epoch 0.
2. Synchronously qualify an Active target through the synthetic observer,
   producing a `Qualified`, `Rejected`, or `Indeterminate` conclusion with
   atomically replaced internal records.
3. List and read accessible targets through a closed safe projection that
   never exposes facts, results, observer identity, handles, or internal record
   links.
4. Retire a target as the sole lifecycle exit, releasing its backing tuple
   while retaining the scope-unique name for the process lifetime.
5. Enter and clear maintenance only through the deterministic target-bound
   fixture trigger, and expire facts only through the injected clock.
6. Project response-only `effectiveAvailability` (Available, Unavailable,
   Maintenance) computed from committed state, freshness, maintenance, and
   viability, with no persisted `status.availability`.

## 3. Terms introduced by this feature

| Term | Meaning within FEATURE-0016 |
|------|-----------------------------|
| ExecutionTarget | CloudProvider-confidential, CloudProvider-scoped resource (VS0-SCHEMA-015) whose create body accepts exactly `metadata.name`, `spec.cloudProviderParticipationRef.uid`, `spec.infrastructureStackRef.uid`, and immutable `spec.targetClass=synthetic-iaas` |
| NormalizedTargetFactSet | Internal-only immutable observation record (VS0-SCHEMA-016) carrying target ref, InfrastructureStack generation, maintenance epoch, viability fingerprint, observer ID/revision, the persisted `factVersion=v1` field, `observedAt`/`expiresAt`, and four named truth values; never projected |
| TargetQualificationResult | Internal-only immutable qualification conclusion (VS0-SCHEMA-017) carrying target/fact-set refs, only the registered `infrastructureStackGeneration` and `maintenanceEpoch` fences (referencing the FactSet by ref rather than duplicating its fields, and carrying no `viabilityFingerprint`), profile version, outcome, reason codes, and evaluation time; never projected |
| `synthetic-iaas` | The only `spec.targetClass` value; an immutable const |
| `sovrunn.synthetic-iaas-observer/v1` | The sole observer identity; fixed provenance with `profileVersion=synthetic-iaas/v1`; `factSchemaVersion=v1`, when named, is fixed observer provenance describing the observer's emitted fact schema, not an additional persisted FactSet field (the persisted FactSet version field is `factVersion=v1`) |
| Four named facts | `compute.vm`, `storage.block`, `storage.object`, `network.private`, each exactly one of Supported/Unsupported/Unknown |
| Qualification outcome | Persisted public qualification: Unqualified, Qualified, Rejected, Indeterminate |
| Qualifying reservation | A lifecycle-service-only in-flight reservation over an Active state; never persisted as a public status and never projected |
| effectiveAvailability | Response-only projection (Available, Unavailable, Maintenance) with no independent persisted writer; does not change resourceVersion/ETag |
| Maintenance epoch | Monotonic fence incremented on every maintenance enter/clear; a fence, not proof that Maintenance is current |
| Current-Maintenance marker | Target-bound, in-process, lifecycle-service-owned marker (target UID, maintenance epoch, active flag); never client-writable or projected; cleared on restart |
| Viability fingerprint | Value composed only of participation UID/effective-active state and stack UID/phase |
| Backing tuple | The `(participation UID, stack UID, targetClass)` tuple; at most one non-Retired tuple is live; retirement releases only the tuple, not the scope-unique name |
| Pre-ServeMux transport guard | A fixed method/path guard that supplies transport-only 405/404 outcomes with no Problem body, authentication, authorization, audit, idempotency, or lifecycle effect |

New local violation codes introduced by this feature (see §5 and §8):
`VS0_TARGET_RETIRED`, `VS0_TARGET_MAINTENANCE`,
`VS0_TARGET_QUALIFICATION_IN_PROGRESS`, `VS0_TARGET_EPOCH_STALE`,
`VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE`,
`VS0_EXECUTION_TARGET_STACK_UNAVAILABLE`,
`VS0_EXECUTION_TARGET_SCOPE_MISMATCH`,
`VS0_EXECUTION_TARGET_VIABILITY_STALE`.

## 4. Normative requirements and acceptance scenarios

### Canonical requirement ledger

The following table copies every approved REQ row from the FEATURE-0016 feature
authority exactly once, with identical columns and cell text. It is an
immutable semantic ledger; the normative detail headings that follow expand
these rows without changing their meaning.

| ID | Behavior | DEC/ADH | VS0 IDs |
|----|----------|---------|---------|
| REQ-F16-01 | Register ExecutionTarget as a CloudProvider-scoped resource; create accepts exactly `metadata.name`, `spec.cloudProviderParticipationRef.uid`, `spec.infrastructureStackRef.uid`, and immutable `spec.targetClass=synthetic-iaas`; scope derives solely from the participation reference; no separate scopeRef persisted or projected | DEC-0042/0057; ADH-058 clause 1 | VS0-SCHEMA-015 |
| REQ-F16-02 | Safe reference access precedes graph validation; create/qualify require FEATURE-0015 effective-Active participation and Active InfrastructureStack in the same CloudProvider scope; at most one non-Retired (participation, stack, targetClass) tuple is live; name is scope-unique for the process lifetime including after retirement | DEC-0042; ADH-058 clause 2 | VS0-SCHEMA-015, VS0-CF-F16-09..13,65,98 |
| REQ-F16-03 | Exactly five Go 1.22 http.ServeMux registrations (collection POST/GET, item GET, qualify POST, retire POST); a fixed pre-ServeMux method/path guard supplies transport-only 405/404 for HEAD and trailing-slash forms with no Problem body, authentication, audit, idempotency, or lifecycle effect | DEC-0042; ADH-058 clause 3 | VS0-CF-F16-01,14,16,20..22,36,44..45,61..64,66,90,93..96,105..110,113..118 |
| REQ-F16-04 | Every POST requires application/json and exactly one Idempotency-Key; qualify/retire require exactly one strong opaque If-Match and a zero-byte body; ETag on create/item GET/qualify/retire only; F0016 ignores Accept and emits fixed response media; closed safe projection with no public conditions member | DEC-0042; ADH-058 clause 4 | VS0-CF-F16-01..06,14,16,20..22,24..25,36,44..45,48,51..60,70..74,83..86,89..91,105..110,113..118 |
| REQ-F16-05 | Executable request precedence: transport guard, authentication, method/media/header-form validation, phase-one safe references, authorization, safe access, strict classification, graph/reference validation, replay/reservation, current If-Match comparison, lifecycle state, atomic commit | DEC-0042; ADH-058 clause 5 | VS0-CF-F16-02..13,18..30,51..60,67..75,84..86,90..91,105..111,113..118 |
| REQ-F16-06 | Idempotency namespace is principal, route, derived CloudProvider scope, concrete target UID for actions, and Idempotency-Key; digest is SHA-256 of sorted create JSON or zero bytes for actions; applies only to create/qualify/retire; 24h/10,000-entry in-process retention; independent reservations across targets and scopes | DEC-0042; ADH-058 clause 5 | VS0-CF-F16-01,20..23,31..35,46..47,49..50,76..78,81..82,86..89,99 |
| REQ-F16-07 | The only observer is sovrunn.synthetic-iaas-observer/v1: target-bound, clock-driven, external-effect-free; exactly four named facts each Supported/Unsupported/Unknown; any Unsupported rejects, otherwise any Unknown is Indeterminate, otherwise Qualified; missing fixture and logical timeout each yield four Unknown facts at now/+60s | DEC-0036/0042; ADH-058 clause 6 | VS0-CF-F16-20..23,46,67..69,90 |
| REQ-F16-08 | ExecutionTargetLifecycleService is the sole committer of persisted target state and internal fact/result records; lifecycle Active/Retired; persisted qualification Unqualified/Qualified/Rejected/Indeterminate; Qualifying is a lifecycle-service-only in-flight reservation never persisted or projected; no status.availability field; response-only effectiveAvailability computed from committed state/freshness/maintenance/viability | DEC-0042/0057; ADH-058 clause 7 | VS0-SCHEMA-015, VS0-STATE-004, VS0-CF-F16-01,20..22,26..30,36,39,41..43,75,77,79..83,87..88,99..101,119..122 |
| REQ-F16-09 | Expiry is injected-clock driven; Retired/Maintenance outrank it; maintenance enters/clears only through a deterministic target-bound synthetic-observer fixture trigger, never an HTTP route or generic event bus; Retire is the sole lifecycle exit; audit matrix for success/denial/safe-denial; expiry is the sole background retry exception; shutdown ordering stops triggers, then aborts reservations, then stops HTTP serving | DEC-0042/0057; ADH-058 clause 8 | VS0-CF-F16-01,29..30,36,38..43,66,75,77,79..83,87..88,99..104,111..112 |
| REQ-F16-10 | The only new local violations are VS0_TARGET_RETIRED, VS0_TARGET_MAINTENANCE, VS0_TARGET_QUALIFICATION_IN_PROGRESS, VS0_TARGET_EPOCH_STALE, VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE, VS0_EXECUTION_TARGET_STACK_UNAVAILABLE, VS0_EXECUTION_TARGET_SCOPE_MISMATCH, and VS0_EXECUTION_TARGET_VIABILITY_STALE; FEATURE-0012 owns top-level Problems | DEC-0042; ADH-058 clause 9 | VS0-CF-F16-10..12,26..30,38,75 |
| REQ-F16-11 | F0024/Phase 3 owns realization/plugin taxonomy; Crossplane is future-only; no F0016 dependency on Crossplane, real adapter, credential, external call, placement, provisioning, plugin execution, or customer IaaS exposure | DEC-0042; ADH-058 clause 10 | VS0-CF-F16-46 |

### Canonical acceptance ledger

The following table copies every approved AC row from the FEATURE-0016 feature
authority exactly once, with identical columns and cell text. It is an
immutable semantic ledger; the acceptance scenarios in §4.1 map to these rows
without changing their meaning.

| ID | Criterion | Conformance ID |
|----|-----------|----------------|
| AC-F16-01 | ExecutionTarget create accepts exactly the four closed fields, derives scope solely from participation, and persists Active/Unqualified at epoch 0 | VS0-CF-F16-01,02,57,84 |
| AC-F16-02 | Safe reference access precedes graph validation; inaccessible/ineffective/non-Active/cross-scope backing return their exact outcomes | VS0-CF-F16-09..13,65,98 |
| AC-F16-03 | Exactly five routes are registered; the pre-ServeMux guard supplies all HEAD/trailing-slash/unmatched-method transport outcomes with no Problem body or side effect | VS0-CF-F16-44,61..64,66,93..96,105..118 |
| AC-F16-04 | Every POST requires application/json and Idempotency-Key; qualify/retire require If-Match and zero-byte body; header-form failures precede body classification | VS0-CF-F16-05,06,24,25,58..60,70..74,85..86 |
| AC-F16-05 | Idempotency replay/reuse/isolation behaves exactly as specified across create, qualify, and retire; GET/LIST never touch idempotency | VS0-CF-F16-14,16,31..35,46..47,49..50,76..78,81..82,86..89,99 |
| AC-F16-06 | The synthetic observer yields exactly four named facts with Supported/Unsupported/Unknown truth and deterministic missing-fixture/timeout behavior | VS0-CF-F16-20..23,67..69,90 |
| AC-F16-07 | ExecutionTargetLifecycleService is the sole state/record committer; Qualifying is never projected; effectiveAvailability follows the exact truth table with no persisted status.availability | VS0-CF-F16-01,20..22,26..30,36,39,41..43,75,77,79..83,87..88,99..101,119..122 |
| AC-F16-08 | Retire wins over in-flight qualification with VS0_TARGET_RETIRED; Maintenance entry wins with VS0_TARGET_MAINTENANCE; both audited, no qualification AuditEvent | VS0-CF-F16-26..28,75,80,83,88..89 |
| AC-F16-09 | Expiry retries once per injected-clock second on audit failure and never returns 500 to a caller; shutdown ordering stops triggers before aborting reservations before stopping HTTP serving | VS0-CF-F16-41,77,79 |
| AC-F16-10 | Audit matrix covers success/denial/safe-denial exactly as specified; append failures return INTERNAL_ERROR/500 with no disclosure | VS0-CF-F16-39,40,82,99..104,112 |
| AC-F16-11 | Only the eight closed local violation codes are introduced; no other F0016-local violation code exists | VS0-CF-F16-10..12,26..30,38,75 |
| AC-F16-12 | No real adapter, credential, external call, Crossplane dependency, placement, provisioning, plugin execution, or customer IaaS exposure exists in active F0016 behavior | VS0-CF-F16-46 |

### 4.1 Normative detail headings

Each approved REQ has exactly one normative detail heading, in approved order.
Each restates observable behavior only and defers all mechanics to design.

#### REQ-F16-01 — ExecutionTarget registration and closed create-field boundary

The ExecutionTarget is a CloudProvider-scoped resource. Its create body accepts
exactly `metadata.name`, `spec.cloudProviderParticipationRef.uid`,
`spec.infrastructureStackRef.uid`, and immutable
`spec.targetClass=synthetic-iaas`. The server alone assigns `metadata.uid`,
`metadata.resourceVersion`, `metadata.generation`, `status.lifecycle`,
`status.qualification`, `status.maintenanceEpoch`, `status.observedGeneration`,
`status.factSetRef`, and `status.qualificationResultRef`. Scope derives solely
from `spec.cloudProviderParticipationRef`; no separate `scopeRef` is persisted
or projected. Unknown fields, `adapterAuthorityRef`, a `SecretRef` extension,
and any client-supplied status/identity field are rejected per §5. A successful
create returns 201 and persists `Active/Unqualified` at maintenance epoch 0 with
`observedGeneration` equal to the referenced InfrastructureStack generation.
Proof: VS0-CF-F16-01,02,04,57,84,91,92.

#### REQ-F16-02 — Safe reference access, backing viability, and uniqueness

Safe reference access precedes graph validation. Create and qualify require a
FEATURE-0015 effective-Active `CloudProviderParticipation` (accepted, with
neither suspension hold set) and an Active `InfrastructureStack` in the same
CloudProvider scope. The name is scope-unique for the process lifetime,
including after retirement; retirement releases only the backing tuple. At most
one non-Retired `(participation UID, stack UID, targetClass)` tuple is live. An
inaccessible backing reference returns a safe outcome; authorized-but-ineffective
participation, authorized non-Active stack, authorized cross-CloudProvider
graph, duplicate name, and live duplicate tuple each return their exact outcome
per §5. Proof: VS0-CF-F16-09..13,65,98.

#### REQ-F16-03 — Closed five-route surface and pre-ServeMux transport guard

Exactly five method/path pairs are the public surface: collection POST, collection
GET, item GET, qualify POST, and retire POST. A fixed pre-ServeMux method/path
guard supplies transport-only 405 for HEAD and other unmatched methods on an
exact F0016 path, and transport-only 404 without redirect for every
trailing-slash variant, each with `X-Content-Type-Options: nosniff` and no
Problem body, authentication, authorization, audit, idempotency, or lifecycle
effect. No PATCH, PUT, DELETE, HEAD, watch, filter, pagination, fact, result, or
adapter-selection route exists. Valid registered pairs bypass the guard and begin
at authentication. Proof: VS0-CF-F16-44,108,113..118 and the transport rows
enumerated in the canonical ledger.

#### REQ-F16-04 — Request grammar, representation, and safe projection

Every POST requires `application/json` and exactly one `Idempotency-Key`.
Qualify and retire additionally require exactly one strong opaque `If-Match` and
a zero-byte body. Header-form validation and current-version comparison are
distinct steps. `ETag` is returned on create, item GET, qualify, and retire, and
not on LIST. Success uses `application/json`; errors use
`application/problem+json`; every response sends `X-Content-Type-Options:
nosniff`. F0016 ignores `Accept` and emits its fixed response media; no route
accepts a query parameter. The safe ExecutionTarget projection is exactly
`metadata.uid`, `metadata.name`, `metadata.generation`,
`metadata.resourceVersion`, the three immutable spec fields, `status.lifecycle`,
`status.qualification`, and response-only `effectiveAvailability`; there is no
public `conditions` member and facts/results/observer data/handles/internal
record links are omitted. Create returns 201; item GET/qualify/retire return
200; LIST returns an ascending JSON array of the projection with no `Location`
header. Proof: VS0-CF-F16-05,06,16,24,25,45,48,51..60,70..74,85,86,105..107,109,110.

#### REQ-F16-05 — Executable request precedence

The request pipeline is ordered: closed transport method/path guard →
authentication → method/media/header-form validation → phase-one safe references
→ authorization → safe access → strict classification → graph/reference
validation → replay/reservation → current `If-Match` version comparison →
lifecycle state → atomic commit. Create alone carries client graph/reference
fields and no `If-Match`; zero-byte actions alone carry `If-Match` and no client
graph. A malformed `If-Match` header-form failure therefore precedes its
non-zero-body failure; current-version comparison occurs after replay/reservation
and before Maintenance state evaluation. Proof:
VS0-CF-F16-07..12,18,19,24,25,70..74,85,86,105..111.

#### REQ-F16-06 — Idempotency namespace, digest, retention, and isolation

The idempotency namespace is principal UID, route pattern, derived CloudProvider
UID, concrete action-target UID when applicable, and `Idempotency-Key`. The
digest is SHA-256 of recursively sorted allowed create JSON, or SHA-256 of zero
bytes for actions; `If-Match`, `Accept`, authentication, and request/correlation
IDs are excluded. Idempotency applies only to create, qualify, and retire;
GET and LIST never require, inspect, reserve, wait on, read, or replay a record.
Successful create/qualification/retirement retain for 24 hours and 10,000
entries only in-process; aborted outcomes do not replay; completed entries evict
at expiry then earliest-expiry/lexical namespace at capacity; InFlight never
evicts. A waiter rechecks current authorization and safe access after waking; a
changed digest is a reuse-mismatch conflict. The same principal/key/digest on a
different action target, or on a create in another derived CloudProvider scope,
forms an independent reservation. Proof:
VS0-CF-F16-14,16,31..35,46,47,49,50,76..78,81,82,87..89,99.

#### REQ-F16-07 — Synthetic observer and fact contract

The only observer is `sovrunn.synthetic-iaas-observer/v1`: target-bound,
clock-driven, external-effect-free. Exactly once each of `compute.vm`,
`storage.block`, `storage.object`, and `network.private` is Supported,
Unsupported, or Unknown. Any Unsupported yields Rejected; otherwise any Unknown
yields Indeterminate; otherwise Qualified. A missing fixture and a
fixture-declared logical timeout are separate deterministic outcomes, each
yielding four Unknown facts at `now` through `now+60s` with no network wait.
Duplicate facts, missing facts, and malformed provenance are observer faults
that fault the qualification and publish no result. Fixed provenance:
`observerID=sovrunn.synthetic-iaas-observer/v1`,
`profileVersion=synthetic-iaas/v1`, with fixed observer provenance
`factSchemaVersion=v1` describing the observer's emitted fact schema; the
persisted FactSet version field is `factVersion=v1` and `factSchemaVersion` is
not an additional persisted FactSet field. Proof:
VS0-CF-F16-20..23,46,67..69,90.

#### REQ-F16-08 — Sole lifecycle committer, qualification, and effective availability

`ExecutionTargetLifecycleService` alone commits target state and internal
fact/result records (VS0-WRITER-006). Lifecycle is Active/Retired; persisted
public qualification is Unqualified/Qualified/Rejected/Indeterminate; Qualifying
is a lifecycle-service-only in-flight reservation, never persisted or projected.
The persisted target has no `status.availability`; the server returns
response-only `effectiveAvailability` (Available, Unavailable, Maintenance) with
no independent persisted writer and no ETag effect. Qualify is synchronous and
may start only from Active Unqualified/Qualified/Rejected/Indeterminate; it
creates a Qualifying reservation that leaves the committed summary, current
record links, and resourceVersion/ETag unchanged until a completed conclusion is
audited and published. Every completed conclusion atomically replaces both
record refs, refreshes `observedGeneration`, advances resourceVersion/ETag, and
appends the required AuditEvent before publication. Malformed observer output,
owner cancellation, required AuditEvent append failure, and every stale captured
fence abort the reservation without target mutation, AuditEvent, or completion.
`effectiveAvailability` is ordered: Retired → Unavailable; else Active with
current Maintenance → Maintenance; else only Active Qualified with fresh current
facts and viable backing → Available; every other Active combination →
Unavailable. Proof:
VS0-CF-F16-01,20..23,26..30,36,39,41..43,75,79..83,87,88,99..101,119..122.

#### REQ-F16-09 — Expiry, maintenance trigger, retirement, shutdown, and audit matrix

Expiry is injected-clock driven; Retired and Maintenance outrank it. Maintenance
enters/clears only through a deterministic, target-bound synthetic-observer
fixture trigger, never an HTTP route or generic event bus. Retire is the sole
lifecycle exit and releases the tuple. The audit matrix audits successful
create, completed qualification, retirement, maintenance enter/clear,
authenticated authorization denial, and safe denial. A request-triggered
required append failure returns INTERNAL_ERROR/500 with no AuditEvent, mutation,
or completion. Expiry is the sole background exception: no request caller, no 500
return, no expiry AuditEvent until append succeeds, and a failed append retries
once per injected-clock second until one succeeds. Service shutdown first stops
acceptance of expiry and Maintenance triggers, then aborts in-flight
qualification and idempotency reservations and wakes waiters, and only then
stops HTTP serving; a trigger accepted after that boundary neither mutates state
nor appends an AuditEvent. Proof:
VS0-CF-F16-36,38,41..43,66,75,77,79,80,83,88,89,99..104,111,112.

#### REQ-F16-10 — Closed local violation set

The only new local violation codes are `VS0_TARGET_RETIRED`,
`VS0_TARGET_MAINTENANCE`, `VS0_TARGET_QUALIFICATION_IN_PROGRESS`,
`VS0_TARGET_EPOCH_STALE`, `VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE`,
`VS0_EXECUTION_TARGET_STACK_UNAVAILABLE`, `VS0_EXECUTION_TARGET_SCOPE_MISMATCH`,
and `VS0_EXECUTION_TARGET_VIABILITY_STALE`. These appear only in `violations[]`,
never as top-level Problem codes; FEATURE-0012 owns the top-level Problem codes
and their HTTP/URN mappings. No other F0016-local violation code exists. Proof:
VS0-CF-F16-10..12,26..30,38,75.

#### REQ-F16-11 — Future boundary and no real external effect

F0024/Phase 3 owns realization and plugin taxonomy; Crossplane is future-only.
FEATURE-0016 has no dependency on Crossplane, a real adapter, a credential, an
external call, placement, provisioning, plugin execution, or customer IaaS
exposure. Process restart clears F0016-owned in-memory targets, records,
idempotency reservations, and pending work, but never inherited FEATURE-0013
AuditEvents, and makes no external observer call. Proof: VS0-CF-F16-46.

### 4.2 Acceptance-scenario coverage mapping

Each approved AC is mapped here to its acceptance scenario and its exact
conformance IDs (see §4.3 for the full conformance ledger). This mapping is the
explicit acceptance evidence outside the canonical acceptance ledger.

| AC | Acceptance scenario (observable) | Conformance IDs |
|----|----------------------------------|-----------------|
| AC-F16-01 | A valid create with exactly the four closed fields persists Active/Unqualified at epoch 0 with participation-derived scope; unknown/adapterAuthorityRef/SecretRef fields are rejected | VS0-CF-F16-01,02,57,84 |
| AC-F16-02 | Inaccessible, ineffective, non-Active, cross-scope, duplicate-name, and live-duplicate-tuple backings each return their exact outcome, with safe access before graph validation | VS0-CF-F16-09,10,11,12,13,65,98 |
| AC-F16-03 | Exactly five routes register; HEAD, PUT, trailing-slash, unmatched-method, and missing/invalid-auth transport cases produce transport-only outcomes or 401 with no side effect | VS0-CF-F16-44,61,62,63,64,66,93,94,95,96,105,106,107,108,109,110,111,112,113,114,115,116,117,118 |
| AC-F16-04 | POST media/key header-form failures, action If-Match/zero-body failures, and their ordering (header-form before body classification) return their exact outcomes | VS0-CF-F16-05,06,24,25,58,59,60,70,71,72,73,74,85,86 |
| AC-F16-05 | Replay, reuse-mismatch, isolation, expiry, eviction, panic/cancellation, and waiter recheck behave exactly as specified; GET/LIST ignore idempotency | VS0-CF-F16-14,16,31,32,33,34,35,46,47,49,50,76,77,78,81,82,86,87,88,89,99 |
| AC-F16-06 | Four Supported/Unsupported/Unknown facts, malformed/duplicate/missing facts, missing fixture, and logical timeout each produce their exact conclusion | VS0-CF-F16-20,21,22,23,67,68,69,90 |
| AC-F16-07 | Sole committer, non-projected Qualifying, effectiveAvailability truth table, viability-only projection changes, and stale-fence aborts hold across all state cases | VS0-CF-F16-01,20,21,22,26,27,28,29,30,36,39,41,42,43,75,77,79,80,81,82,83,87,88,99,100,101,119,120,121,122 |
| AC-F16-08 | Retire wins with VS0_TARGET_RETIRED and Maintenance entry wins with VS0_TARGET_MAINTENANCE over in-flight qualification; both audited, no qualification AuditEvent | VS0-CF-F16-26,27,28,75,80,83,88,89 |
| AC-F16-09 | Fact expiry with successful or failing append, and shutdown ordering, behave exactly as specified with no 500 to a caller | VS0-CF-F16-41,77,79 |
| AC-F16-10 | Success/denial/safe-denial audit and every append-failure case return INTERNAL_ERROR/500 with no disclosure | VS0-CF-F16-39,40,82,99,100,101,102,103,104,112 |
| AC-F16-11 | Only the eight closed local violation codes appear; ineffective/non-Active/cross-scope/retired/maintenance/in-progress/epoch-stale/viability-stale cases carry them | VS0-CF-F16-10,11,12,26,27,28,29,30,38,75 |
| AC-F16-12 | Process restart demonstrates no external observer call and no persisted external effect | VS0-CF-F16-46 |

### 4.3 Exact conformance semantics ledger

For every referenced case in the FEATURE-0016 local conformance suite this
ledger copies the registered fields: ID, owner, inputs, expectedState,
expectedError, expectedViolation (where registered), expectedSideEffects, and
gate. Each conformance case proves
only its registered scenario. A conformance case is never a substitute for a
state-machine ID (`VS0-STATE-004`) and no downstream or shared case
(`VS0-CF-F10`, `VS0-CF-X03`, `VS0-CF-HP01`) is used as FEATURE-0016-local proof.

| ID | owner | inputs | expectedState | expectedError | expectedViolation | expectedSideEffects | gate |
|----|-------|--------|---------------|---------------|--------------------|----------------------|------|
| VS0-CF-F16-01 | FEATURE-0016 | Valid create | 201 safe projection | null | — | 201 safe projection; persists Active/Unqualified, observedGeneration equals referenced InfrastructureStack generation, maintenance epoch 0, and projects effectiveAvailability Unavailable; ETag, AuditEvent, completed replay. | feature |
| VS0-CF-F16-02 | FEATURE-0016 | Authorized create with one unknown field | unchanged | UNKNOWN_FIELD | — | 400 UNKNOWN_FIELD; no mutation, audit, or completion. | feature |
| VS0-CF-F16-03 | FEATURE-0016 | Authorized create with client status | unchanged | AUTHORIZATION_DENIED | VS0_STATUS_FIELD_WRITE | 403 AUTHORIZATION_DENIED plus VS0_STATUS_FIELD_WRITE; audited denial. | feature |
| VS0-CF-F16-04 | FEATURE-0016 | Authorized create with client metadata.uid | unchanged | AUTHORIZATION_DENIED | VS0_SYSTEM_OWNED_FIELD_WRITE | 403 AUTHORIZATION_DENIED plus VS0_SYSTEM_OWNED_FIELD_WRITE; audited denial. | feature |
| VS0-CF-F16-05 | FEATURE-0016 | Create/action with a missing Idempotency-Key | unchanged | MALFORMED_REQUEST | — | 400 MALFORMED_REQUEST; no publication. | feature |
| VS0-CF-F16-06 | FEATURE-0016 | POST with non-application/json media | unchanged | UNSUPPORTED_MEDIA_TYPE | — | 415 UNSUPPORTED_MEDIA_TYPE; no publication. | feature |
| VS0-CF-F16-07 | FEATURE-0016 | Create with missing authentication | unchanged | AUTH_REQUIRED | — | 401 AUTH_REQUIRED; no mutation, audit, or replay record. | feature |
| VS0-CF-F16-08 | FEATURE-0016 | Authenticated create without executiontarget.write | unchanged | AUTHORIZATION_DENIED | — | 403 AUTHORIZATION_DENIED; audited denial; no publication. | feature |
| VS0-CF-F16-09 | FEATURE-0016 | Inaccessible create backing reference | Safe 404 RESOURCE_NOT_FOUND | RESOURCE_NOT_FOUND | VS0_AUTHORIZATION_SAFE_DENIAL | Safe 404 RESOURCE_NOT_FOUND plus VS0_AUTHORIZATION_SAFE_DENIAL; audited denial. | feature |
| VS0-CF-F16-10 | FEATURE-0016 | Authorized ineffective participation | unchanged | CONFLICT | VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE | 409 CONFLICT plus VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE; no publication. | feature |
| VS0-CF-F16-11 | FEATURE-0016 | Authorized non-Active InfrastructureStack | unchanged | CONFLICT | VS0_EXECUTION_TARGET_STACK_UNAVAILABLE | 409 CONFLICT plus VS0_EXECUTION_TARGET_STACK_UNAVAILABLE; no publication. | feature |
| VS0-CF-F16-12 | FEATURE-0016 | Authorized cross-CloudProvider participation/stack graph | unchanged | VALIDATION_FAILED | VS0_EXECUTION_TARGET_SCOPE_MISMATCH | 422 VALIDATION_FAILED plus VS0_EXECUTION_TARGET_SCOPE_MISMATCH; no publication. | feature |
| VS0-CF-F16-13 | FEATURE-0016 | Create with a duplicate derived name | unchanged | ALREADY_EXISTS | — | 409 ALREADY_EXISTS; no publication. | feature |
| VS0-CF-F16-14 | FEATURE-0016 | Authorized LIST with an arbitrary Idempotency-Key | 200 accessible Active and Retired targets, ascending UID, no ETag | null | — | 200 accessible Active and Retired targets, ascending UID, no ETag; ignores the key and creates, waits on, reads, and replays no idempotency record. | feature |
| VS0-CF-F16-15 | FEATURE-0016 | LIST without applicable executiontarget.read | unchanged | AUTHORIZATION_DENIED | — | 403 AUTHORIZATION_DENIED; audited denial. | feature |
| VS0-CF-F16-16 | FEATURE-0016 | Authorized item GET with an arbitrary Idempotency-Key | 200 ETag safe projection without facts, result, observer, or handle | null | — | 200 ETag safe projection without facts, result, observer, or handle; ignores the key and creates, waits on, reads, and replays no idempotency record. | feature |
| VS0-CF-F16-17 | FEATURE-0016 | Inaccessible direct item GET | Safe 404 RESOURCE_NOT_FOUND | RESOURCE_NOT_FOUND | VS0_AUTHORIZATION_SAFE_DENIAL | Safe 404 RESOURCE_NOT_FOUND plus VS0_AUTHORIZATION_SAFE_DENIAL; audited denial. | feature |
| VS0-CF-F16-18 | FEATURE-0016 | Qualify without executiontarget.qualify | unchanged | AUTHORIZATION_DENIED | — | 403 AUTHORIZATION_DENIED; audited denial; no publication. | feature |
| VS0-CF-F16-19 | FEATURE-0016 | Inaccessible direct qualify action | Safe 404 RESOURCE_NOT_FOUND | RESOURCE_NOT_FOUND | VS0_AUTHORIZATION_SAFE_DENIAL | Safe 404 RESOURCE_NOT_FOUND plus VS0_AUTHORIZATION_SAFE_DENIAL; audited denial. | feature |
| VS0-CF-F16-20 | FEATURE-0016 | Qualify, four Supported facts | 200 Qualified/effectiveAvailability Available | null | — | 200 Qualified/effectiveAvailability Available; FactSet/Result replacement, observedGeneration refresh, new ETag, AuditEvent, completed replay. | feature |
| VS0-CF-F16-21 | FEATURE-0016 | Qualify, any Unsupported fact | 200 Rejected/effectiveAvailability Unavailable | null | — | 200 Rejected/effectiveAvailability Unavailable; records, observedGeneration refresh, new ETag, AuditEvent, completed replay. | feature |
| VS0-CF-F16-22 | FEATURE-0016 | Qualify with a missing target-bound fixture | 200 Indeterminate/effectiveAvailability Unavailable | null | — | 200 Indeterminate/effectiveAvailability Unavailable; four Unknown facts now/+60s, records, observedGeneration refresh, new ETag, AuditEvent, completed replay. | feature |
| VS0-CF-F16-23 | FEATURE-0016 | Qualify from Qualified receives one malformed fact while all captured fences remain current | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; abort reservation with unchanged committed target, ETag, and FactSet/Result links; publish no AuditEvent or idempotency completion. | feature |
| VS0-CF-F16-24 | FEATURE-0016 | Action with a missing If-Match | unchanged | STALE_RESOURCE_VERSION | — | 412 STALE_RESOURCE_VERSION; no publication; header-form failure precedes body classification. | feature |
| VS0-CF-F16-25 | FEATURE-0016 | Action valid If-Match with non-zero body | unchanged | MALFORMED_REQUEST | — | 400 MALFORMED_REQUEST; no publication. | feature |
| VS0-CF-F16-26 | FEATURE-0016 | Qualify Retired target | unchanged | CONFLICT | VS0_TARGET_RETIRED | 409 CONFLICT plus VS0_TARGET_RETIRED; no publication. | feature |
| VS0-CF-F16-27 | FEATURE-0016 | Qualify Maintenance target | unchanged | CONFLICT | VS0_TARGET_MAINTENANCE | 409 CONFLICT plus VS0_TARGET_MAINTENANCE; no publication. | feature |
| VS0-CF-F16-28 | FEATURE-0016 | Different-key qualify while Qualifying | unchanged | CONFLICT | VS0_TARGET_QUALIFICATION_IN_PROGRESS | 409 CONFLICT plus VS0_TARGET_QUALIFICATION_IN_PROGRESS; no publication. | feature |
| VS0-CF-F16-29 | FEATURE-0016 | Referenced InfrastructureStack generation changes at qualification commit | unchanged | STALE_RESOURCE_VERSION | VS0_TARGET_EPOCH_STALE | 412 STALE_RESOURCE_VERSION plus VS0_TARGET_EPOCH_STALE; abort reservation with unchanged committed target, ETag, and FactSet/Result links; publish no qualification AuditEvent or completion. | feature |
| VS0-CF-F16-30 | FEATURE-0016 | Viability fingerprint changes at qualification commit | unchanged | STALE_RESOURCE_VERSION | VS0_EXECUTION_TARGET_VIABILITY_STALE | 412 STALE_RESOURCE_VERSION plus VS0_EXECUTION_TARGET_VIABILITY_STALE; abort reservation with unchanged committed target, ETag, and FactSet/Result links; projection follows current viability and publishes no qualification AuditEvent or completion. | feature |
| VS0-CF-F16-31 | FEATURE-0016 | Same-key completed replay after current authorization and safe access | Original successful status/body, application/json, fresh correlation only | null | — | Original successful status/body, application/json, fresh correlation only; no second audit. | feature |
| VS0-CF-F16-32 | FEATURE-0016 | Same-key replay requester loses its grant while target remains safe-accessible | Current 403 AUTHORIZATION_DENIED plus one redacted AuditEvent | AUTHORIZATION_DENIED | — | Current 403 AUTHORIZATION_DENIED plus one redacted AuditEvent; no stored-result disclosure; audit append failure is INTERNAL_ERROR/500 with no 403 disclosure. | feature |
| VS0-CF-F16-33 | FEATURE-0016 | Same namespace with changed digest | unchanged | CONFLICT | VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH | 409 CONFLICT plus VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH; original completion unchanged. | feature |
| VS0-CF-F16-34 | FEATURE-0016 | Completed record expires after 24 hours | Record unavailable | null | — | Record unavailable; next valid request is processed as new; no InFlight eviction. | feature |
| VS0-CF-F16-35 | FEATURE-0016 | Owner panic after reservation | Reservation aborted and waiters woken without completion | null | — | Reservation aborted and waiters woken without completion; cancelled waiter only detaches. | feature |
| VS0-CF-F16-36 | FEATURE-0016 | Valid retire from any Active combination | 200 Retired/Unqualified/effectiveAvailability Unavailable, ETag, refs cleared, AuditEvent, tuple release, completed replay | null | — | 200 Retired/Unqualified/effectiveAvailability Unavailable, ETag, refs cleared, AuditEvent, tuple release, completed replay. | feature |
| VS0-CF-F16-37 | FEATURE-0016 | Authenticated retire without executiontarget.write | unchanged | AUTHORIZATION_DENIED | — | 403 AUTHORIZATION_DENIED; audited denial; no publication. | feature |
| VS0-CF-F16-38 | FEATURE-0016 | New-key retire after retirement | unchanged | CONFLICT | VS0_TARGET_RETIRED | 409 CONFLICT plus VS0_TARGET_RETIRED; no resurrection. | feature |
| VS0-CF-F16-39 | FEATURE-0016 | Valid create required AuditEvent append fails | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; no AuditEvent, target creation, or idempotency completion. | feature |
| VS0-CF-F16-40 | FEATURE-0016 | Create authorization-denial AuditEvent append fails | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; no AuditEvent; the original 403 is not disclosed. | feature |
| VS0-CF-F16-41 | FEATURE-0016 | Fact expiry with successful required AuditEvent append | effectiveAvailability Unavailable with unchanged ETag | null | — | effectiveAvailability Unavailable with unchanged ETag; exactly one expiry AuditEvent. | feature |
| VS0-CF-F16-42 | FEATURE-0016 | Current fenced maintenance-entry trigger | effectiveAvailability Maintenance | null | — | effectiveAvailability Maintenance; maintenance epoch and resourceVersion/ETag increment; current FactSet/Result links clear; lifecycle-service current-Maintenance marker is active at the new epoch; successful entry audited. | feature |
| VS0-CF-F16-43 | FEATURE-0016 | Fenced maintenance clear | Active/Unqualified/effectiveAvailability Unavailable | null | — | Active/Unqualified/effectiveAvailability Unavailable; maintenance epoch and resourceVersion/ETag increment; current FactSet/Result links clear; matching current-Maintenance marker removed; successful clear audited. | feature |
| VS0-CF-F16-44 | FEATURE-0016 | HEAD on the item GET path | Transport-only 405 with Allow: GET and X-Content-Type-Options: nosniff | null | — | Transport-only 405 with Allow: GET and X-Content-Type-Options: nosniff; no Problem body, mutation, audit, or idempotency. | feature |
| VS0-CF-F16-45 | FEATURE-0016 | Successful GET with any Accept value | Server ignores Accept and emits application/json with nosniff | null | — | Server ignores Accept and emits application/json with nosniff. | feature |
| VS0-CF-F16-46 | FEATURE-0016 | Process restart | F0016-owned in-memory targets, records, idempotency reservations, and pending work are absent | null | — | F0016-owned in-memory targets, records, idempotency reservations, and pending work are absent; inherited FEATURE-0013 AuditEvents are not deleted; no external observer call. | feature |
| VS0-CF-F16-47 | FEATURE-0016 | Same-key replay target is no longer safe-accessible | Current safe 404 RESOURCE_NOT_FOUND plus one redacted AuditEvent | RESOURCE_NOT_FOUND | — | Current safe 404 RESOURCE_NOT_FOUND plus one redacted AuditEvent; no stored-result disclosure; audit append failure is INTERNAL_ERROR/500 with no 404 disclosure. | feature |
| VS0-CF-F16-48 | FEATURE-0016 | Problem response with any Accept value | Server ignores Accept and emits application/problem+json with nosniff | null | — | Server ignores Accept and emits application/problem+json with nosniff. | feature |
| VS0-CF-F16-49 | FEATURE-0016 | Same principal/key/digest action on another target UID | Independent reservation and outcome | null | — | Independent reservation and outcome; never replays the first target. | feature |
| VS0-CF-F16-50 | FEATURE-0016 | Same principal/key/digest create in another derived CloudProvider scope | Independent reservation and outcome | null | — | Independent reservation and outcome; never replays the first scope. | feature |
| VS0-CF-F16-51 | FEATURE-0016 | Malformed JSON request body | unchanged | MALFORMED_REQUEST | — | 400 MALFORMED_REQUEST; no mutation, audit, or completion. | feature |
| VS0-CF-F16-52 | FEATURE-0016 | Duplicate top-level JSON member | unchanged | DUPLICATE_FIELD | — | 400 DUPLICATE_FIELD; no mutation, audit, or completion. | feature |
| VS0-CF-F16-53 | FEATURE-0016 | Oversized request body | unchanged | REQUEST_TOO_LARGE | — | 400 REQUEST_TOO_LARGE; no mutation, audit, or completion. | feature |
| VS0-CF-F16-54 | FEATURE-0016 | Missing required create field | unchanged | VALIDATION_FAILED | — | 422 VALIDATION_FAILED; no mutation, audit, or completion. | feature |
| VS0-CF-F16-55 | FEATURE-0016 | Invalid typed reference shape | unchanged | VALIDATION_FAILED | — | 422 VALIDATION_FAILED; no mutation, audit, or completion. | feature |
| VS0-CF-F16-56 | FEATURE-0016 | targetClass other than synthetic-iaas | unchanged | VALIDATION_FAILED | — | 422 VALIDATION_FAILED; no mutation, audit, or completion. | feature |
| VS0-CF-F16-57 | FEATURE-0016 | Authorized create with adapterAuthorityRef | unchanged | UNKNOWN_FIELD | — | 400 UNKNOWN_FIELD; no mutation, audit, or completion. | feature |
| VS0-CF-F16-58 | FEATURE-0016 | Create/action with an empty Idempotency-Key | unchanged | MALFORMED_REQUEST | — | 400 MALFORMED_REQUEST; no publication. | feature |
| VS0-CF-F16-59 | FEATURE-0016 | Create/action with a malformed Idempotency-Key | unchanged | MALFORMED_REQUEST | — | 400 MALFORMED_REQUEST; no publication. | feature |
| VS0-CF-F16-60 | FEATURE-0016 | Create/action with an overlong Idempotency-Key | unchanged | MALFORMED_REQUEST | — | 400 MALFORMED_REQUEST; no publication. | feature |
| VS0-CF-F16-61 | FEATURE-0016 | LIST with missing authentication | unchanged | AUTH_REQUIRED | — | 401 AUTH_REQUIRED; no mutation, audit, or replay record. | feature |
| VS0-CF-F16-62 | FEATURE-0016 | Item GET with missing authentication | unchanged | AUTH_REQUIRED | — | 401 AUTH_REQUIRED; no mutation, audit, or replay record. | feature |
| VS0-CF-F16-63 | FEATURE-0016 | Qualify with missing authentication | unchanged | AUTH_REQUIRED | — | 401 AUTH_REQUIRED; no mutation, audit, or replay record. | feature |
| VS0-CF-F16-64 | FEATURE-0016 | Retire with missing authentication | unchanged | AUTH_REQUIRED | — | 401 AUTH_REQUIRED; no mutation, audit, or replay record. | feature |
| VS0-CF-F16-65 | FEATURE-0016 | Create with a live duplicate backing tuple | unchanged | ALREADY_EXISTS | — | 409 ALREADY_EXISTS; no publication. | feature |
| VS0-CF-F16-66 | FEATURE-0016 | Inaccessible direct retire action | Safe 404 RESOURCE_NOT_FOUND | RESOURCE_NOT_FOUND | VS0_AUTHORIZATION_SAFE_DENIAL | Safe 404 RESOURCE_NOT_FOUND plus VS0_AUTHORIZATION_SAFE_DENIAL; audited denial. | feature |
| VS0-CF-F16-67 | FEATURE-0016 | Qualify receives a duplicate named fact while all captured fences remain current | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; abort reservation with unchanged committed target, ETag, and FactSet/Result links; publish no AuditEvent or completion. | feature |
| VS0-CF-F16-68 | FEATURE-0016 | Qualify receives a missing named fact while all captured fences remain current | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; abort reservation with unchanged committed target, ETag, and FactSet/Result links; publish no AuditEvent or completion. | feature |
| VS0-CF-F16-69 | FEATURE-0016 | Qualify receives malformed fact provenance while all captured fences remain current | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; abort reservation with unchanged committed target, ETag, and FactSet/Result links; publish no AuditEvent or completion. | feature |
| VS0-CF-F16-70 | FEATURE-0016 | Action with a malformed If-Match | unchanged | STALE_RESOURCE_VERSION | — | 412 STALE_RESOURCE_VERSION; no publication; header-form failure precedes body classification. | feature |
| VS0-CF-F16-71 | FEATURE-0016 | Action with a weak If-Match | unchanged | STALE_RESOURCE_VERSION | — | 412 STALE_RESOURCE_VERSION; no publication; header-form failure precedes body classification. | feature |
| VS0-CF-F16-72 | FEATURE-0016 | Action with wildcard If-Match | unchanged | STALE_RESOURCE_VERSION | — | 412 STALE_RESOURCE_VERSION; no publication; header-form failure precedes body classification. | feature |
| VS0-CF-F16-73 | FEATURE-0016 | Action with multiple If-Match headers | unchanged | STALE_RESOURCE_VERSION | — | 412 STALE_RESOURCE_VERSION; no publication; header-form failure precedes body classification. | feature |
| VS0-CF-F16-74 | FEATURE-0016 | Action with a syntactically valid stale If-Match | unchanged | STALE_RESOURCE_VERSION | — | 412 STALE_RESOURCE_VERSION; no publication; current-version comparison follows replay/reservation and precedes lifecycle evaluation. | feature |
| VS0-CF-F16-75 | FEATURE-0016 | Maintenance epoch changes because Maintenance entry wins during qualification | unchanged | STALE_RESOURCE_VERSION | VS0_TARGET_EPOCH_STALE | 412 STALE_RESOURCE_VERSION plus VS0_TARGET_EPOCH_STALE; no qualification conclusion, AuditEvent, or completion; the independently winning Maintenance transition is audited and supplies its own target mutation and ETag. | feature |
| VS0-CF-F16-76 | FEATURE-0016 | Completed record is selected for deterministic capacity eviction | Record unavailable | null | — | Record unavailable; next valid request is processed as new; no InFlight eviction. | feature |
| VS0-CF-F16-77 | FEATURE-0016 | Server shutdown after reservation | Stop expiry/Maintenance trigger acceptance, abort qualification and idempotency reservations, wake waiters, then stop HTTP serving | null | — | Stop expiry/Maintenance trigger acceptance, abort qualification and idempotency reservations, wake waiters, then stop HTTP serving; no completion or post-stop trigger AuditEvent; last committed target projection, ETag, and links remain unchanged. | feature |
| VS0-CF-F16-78 | FEATURE-0016 | Owner cancellation after reservation | Reservation aborted and waiters woken without completion | null | — | Reservation aborted and waiters woken without completion; cancelled waiter only detaches. | feature |
| VS0-CF-F16-79 | FEATURE-0016 | Fact expiry with required AuditEvent append failure | effectiveAvailability remains Unavailable with unchanged ETag | null | — | effectiveAvailability remains Unavailable with unchanged ETag; one expiry audit retries each injected-clock second until success. | feature |
| VS0-CF-F16-80 | FEATURE-0016 | Stale fenced maintenance-entry trigger | No state change, AuditEvent, or record publication | null | — | No state change, AuditEvent, or record publication. | feature |
| VS0-CF-F16-81 | FEATURE-0016 | Same-key waiter wakes after its owner aborts on a stale InfrastructureStack-generation qualification | No cached completion | null | — | No cached completion; waiter repeats current authorization and safe access, then may establish a new reservation if its current request remains valid; the owner's target remains unchanged as in F16-29. | feature |
| VS0-CF-F16-82 | FEATURE-0016 | Required qualification AuditEvent append failure after reservation while all captured fences remain current | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; abort reservation with unchanged committed target, ETag, and FactSet/Result links; wake waiters without AuditEvent or completion. | feature |
| VS0-CF-F16-83 | FEATURE-0016 | Stale fenced maintenance-clear trigger | No state change, AuditEvent, or record publication | null | — | No state change, AuditEvent, or record publication. | feature |
| VS0-CF-F16-84 | FEATURE-0016 | Authorized create with a SecretRef extension field | unchanged | UNKNOWN_FIELD | — | 400 UNKNOWN_FIELD; no mutation, audit, or completion. | feature |
| VS0-CF-F16-85 | FEATURE-0016 | Qualify Maintenance target with syntactically valid stale If-Match | unchanged | STALE_RESOURCE_VERSION | — | 412 STALE_RESOURCE_VERSION; current-version comparison precedes Maintenance state evaluation; no publication. | feature |
| VS0-CF-F16-86 | FEATURE-0016 | Qualify with malformed If-Match and non-zero body | unchanged | STALE_RESOURCE_VERSION | — | 412 STALE_RESOURCE_VERSION; header-form failure precedes body classification; no publication. | feature |
| VS0-CF-F16-87 | FEATURE-0016 | Same-key completed qualify after target ETag changes | Original successful 200 JSON response | null | — | Original successful 200 JSON response; If-Match exclusion permits replay after current authorization and safe access; no second audit. | feature |
| VS0-CF-F16-88 | FEATURE-0016 | Retire wins while qualify is in flight | Qualifying caller gets 409 CONFLICT | CONFLICT | VS0_TARGET_RETIRED | Qualifying caller gets 409 CONFLICT plus VS0_TARGET_RETIRED; no qualification result, completion, or qualification AuditEvent; winning retirement is audited. | feature |
| VS0-CF-F16-89 | FEATURE-0016 | Maintenance entry wins while qualify is in flight | Qualifying caller gets 409 CONFLICT | CONFLICT | VS0_TARGET_MAINTENANCE | Qualifying caller gets 409 CONFLICT plus VS0_TARGET_MAINTENANCE; no qualification result, completion, or qualification AuditEvent; winning maintenance entry is audited. | feature |
| VS0-CF-F16-90 | FEATURE-0016 | Qualify with fixture-declared logical timeout | 200 Indeterminate/effectiveAvailability Unavailable | null | — | 200 Indeterminate/effectiveAvailability Unavailable; four Unknown facts now/+60s, records, AuditEvent, completed replay; no elapsed network wait. | feature |
| VS0-CF-F16-91 | FEATURE-0016 | Authorized create with client metadata.resourceVersion | unchanged | AUTHORIZATION_DENIED | VS0_SYSTEM_OWNED_FIELD_WRITE | 403 AUTHORIZATION_DENIED plus VS0_SYSTEM_OWNED_FIELD_WRITE; audited denial. | feature |
| VS0-CF-F16-92 | FEATURE-0016 | Authorized create with client metadata.generation | unchanged | AUTHORIZATION_DENIED | VS0_SYSTEM_OWNED_FIELD_WRITE | 403 AUTHORIZATION_DENIED plus VS0_SYSTEM_OWNED_FIELD_WRITE; audited denial. | feature |
| VS0-CF-F16-93 | FEATURE-0016 | Create with invalid authentication | unchanged | AUTH_REQUIRED | — | 401 AUTH_REQUIRED; no mutation, audit, or replay record. | feature |
| VS0-CF-F16-94 | FEATURE-0016 | LIST with invalid authentication | unchanged | AUTH_REQUIRED | — | 401 AUTH_REQUIRED; no mutation, audit, or replay record. | feature |
| VS0-CF-F16-95 | FEATURE-0016 | Item GET with invalid authentication | unchanged | AUTH_REQUIRED | — | 401 AUTH_REQUIRED; no mutation, audit, or replay record. | feature |
| VS0-CF-F16-96 | FEATURE-0016 | Qualify with invalid authentication | unchanged | AUTH_REQUIRED | — | 401 AUTH_REQUIRED; no mutation, audit, or replay record. | feature |
| VS0-CF-F16-97 | FEATURE-0016 | Retire with invalid authentication | unchanged | AUTH_REQUIRED | — | 401 AUTH_REQUIRED; no mutation, audit, or replay record. | feature |
| VS0-CF-F16-98 | FEATURE-0016 | Create with a name reserved by a Retired target | unchanged | ALREADY_EXISTS | — | 409 ALREADY_EXISTS; retirement released the backing tuple but not the scope/name reservation. | feature |
| VS0-CF-F16-99 | FEATURE-0016 | Valid retirement required AuditEvent append fails | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; no AuditEvent, retirement, tuple release, ETag advance, or idempotency completion. | feature |
| VS0-CF-F16-100 | FEATURE-0016 | Current maintenance-entry required AuditEvent append fails | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; no AuditEvent, maintenance epoch/resourceVersion advance, link clear, or transition publication. | feature |
| VS0-CF-F16-101 | FEATURE-0016 | Current maintenance-clear required AuditEvent append fails | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; no AuditEvent, maintenance epoch/resourceVersion advance, link clear, or transition publication. | feature |
| VS0-CF-F16-102 | FEATURE-0016 | Inaccessible direct item GET required AuditEvent append fails | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; no AuditEvent; the original safe 404 is not disclosed. | feature |
| VS0-CF-F16-103 | FEATURE-0016 | Same-key replay authorization denial AuditEvent append fails | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; no AuditEvent; the original replay-time 403 is not disclosed. | feature |
| VS0-CF-F16-104 | FEATURE-0016 | Same-key replay safe-denial AuditEvent append fails | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; no AuditEvent; the original replay-time 404 is not disclosed. | feature |
| VS0-CF-F16-105 | FEATURE-0016 | Collection create with an If-Match header | unchanged | MALFORMED_REQUEST | — | 400 MALFORMED_REQUEST; no mutation, AuditEvent, or idempotency completion. | feature |
| VS0-CF-F16-106 | FEATURE-0016 | Authenticated URI-query equivalence family: POST collection, GET collection, GET item, POST qualify, and POST retire, each with `?x=1` | For every named valid method/path pair, 400 MALFORMED_REQUEST before safe resolution | MALFORMED_REQUEST | — | For every named valid method/path pair, 400 MALFORMED_REQUEST before safe resolution; no reference lookup, mutation, AuditEvent, idempotency reservation, or completion. | feature |
| VS0-CF-F16-107 | FEATURE-0016 | Authenticated item GET with a malformed canonical UID segment | before safe resolution | MALFORMED_REQUEST | — | 400 MALFORMED_REQUEST before safe resolution; no reference lookup, mutation, AuditEvent, idempotency reservation, or completion. | feature |
| VS0-CF-F16-108 | FEATURE-0016 | Item GET with trailing slash | Transport-only 404 without redirect or Problem body, with X-Content-Type-Options: nosniff | null | — | Transport-only 404 without redirect or Problem body, with X-Content-Type-Options: nosniff; no authentication, mutation, audit, or idempotency. | feature |
| VS0-CF-F16-109 | FEATURE-0016 | Authenticated qualify action with a malformed canonical UID segment | before safe resolution | MALFORMED_REQUEST | — | 400 MALFORMED_REQUEST before safe resolution; no reference lookup, mutation, AuditEvent, idempotency reservation, or completion. | feature |
| VS0-CF-F16-110 | FEATURE-0016 | Authenticated retire action with a malformed canonical UID segment | before safe resolution | MALFORMED_REQUEST | — | 400 MALFORMED_REQUEST before safe resolution; no reference lookup, mutation, AuditEvent, idempotency reservation, or completion. | feature |
| VS0-CF-F16-111 | FEATURE-0016 | Direct item GET of a safe-accessible target without executiontarget.read | After UID-only safe resolution derives the target scope, 403 AUTHORIZATION_DENIED | AUTHORIZATION_DENIED | — | After UID-only safe resolution derives the target scope, 403 AUTHORIZATION_DENIED; one redacted AuditEvent; no target disclosure, mutation, or idempotency record. | feature |
| VS0-CF-F16-112 | FEATURE-0016 | Denial AuditEvent append-failure equivalence family: LIST without read; direct GET without read; qualify without qualify grant; retire without write grant; inaccessible direct qualify; inaccessible direct retire | For every named member, 500 INTERNAL_ERROR | INTERNAL_ERROR | — | For every named member, 500 INTERNAL_ERROR; its original 403/404 and violation are not disclosed; no AuditEvent, mutation, lifecycle/record publication, or idempotency completion. | feature |
| VS0-CF-F16-113 | FEATURE-0016 | HEAD on collection path | Transport-only 405 with Allow: GET, POST and X-Content-Type-Options: nosniff | null | — | Transport-only 405 with Allow: GET, POST and X-Content-Type-Options: nosniff; no Problem body, authentication, audit, idempotency, mutation, or lifecycle effect. | feature |
| VS0-CF-F16-114 | FEATURE-0016 | HEAD equivalence family: qualify action and retire action paths | For each named action path, transport-only 405 with Allow: POST and X-Content-Type-Options: nosniff | null | — | For each named action path, transport-only 405 with Allow: POST and X-Content-Type-Options: nosniff; no Problem body, authentication, audit, idempotency, mutation, or lifecycle effect. | feature |
| VS0-CF-F16-115 | FEATURE-0016 | Trailing-slash equivalence family: POST collection, GET collection, POST qualify action, and POST retire action | For each named method/path pair, transport-only 404 without redirect or Problem body, with X-Content-Type-Options: nosniff | null | — | For each named method/path pair, transport-only 404 without redirect or Problem body, with X-Content-Type-Options: nosniff; no authentication, audit, idempotency, mutation, or lifecycle effect. | feature |
| VS0-CF-F16-116 | FEATURE-0016 | PUT on collection path | Transport-only 405 with Allow: GET, POST and X-Content-Type-Options: nosniff | null | — | Transport-only 405 with Allow: GET, POST and X-Content-Type-Options: nosniff; no Problem body, authentication, audit, idempotency, mutation, or lifecycle effect. | feature |
| VS0-CF-F16-117 | FEATURE-0016 | PUT on item path | Transport-only 405 with Allow: GET and X-Content-Type-Options: nosniff | null | — | Transport-only 405 with Allow: GET and X-Content-Type-Options: nosniff; no Problem body, authentication, audit, idempotency, mutation, or lifecycle effect. | feature |
| VS0-CF-F16-118 | FEATURE-0016 | PUT equivalence family: qualify action and retire action paths | For each named action path, transport-only 405 with Allow: POST and X-Content-Type-Options: nosniff | null | — | For each named action path, transport-only 405 with Allow: POST and X-Content-Type-Options: nosniff; no Problem body, authentication, audit, idempotency, mutation, or lifecycle effect. | feature |
| VS0-CF-F16-119 | FEATURE-0016 | Fresh Qualified target backing becomes non-viable | GET projects effectiveAvailability Unavailable with unchanged target ETag and no F0016 mutation or AuditEvent | null | — | GET projects effectiveAvailability Unavailable with unchanged target ETag and no F0016 mutation or AuditEvent. | feature |
| VS0-CF-F16-120 | FEATURE-0016 | Fresh Qualified target backing returns viable before fact expiry | GET projects effectiveAvailability Available with unchanged target ETag and no F0016 mutation or AuditEvent | null | — | GET projects effectiveAvailability Available with unchanged target ETag and no F0016 mutation or AuditEvent. | feature |
| VS0-CF-F16-121 | FEATURE-0016 | Qualify starts from Active/Qualified with current fences | Create an internal Qualifying reservation | null | — | Create an internal Qualifying reservation; retain the public committed target projection, ETag, and FactSet/Result links, with no AuditEvent or idempotency completion. | feature |
| VS0-CF-F16-122 | FEATURE-0016 | Current-fence observer fault after qualification starts from Active/Qualified | unchanged | INTERNAL_ERROR | — | 500 INTERNAL_ERROR; abort reservation with unchanged committed target, ETag, and FactSet/Result links; publish no AuditEvent or idempotency completion. | feature |

## 5. Security, privacy, compatibility, and operational requirements

### 5.1 Authentication and authorization

- Missing or invalid authentication on any registered method/path pair returns
  `AUTH_REQUIRED`/401 with no mutation, audit, or replay record
  (VS0-CF-F16-07,61..64,93..97).
- Create and retire require `executiontarget.write`; reads require
  `executiontarget.read`; qualify requires `executiontarget.qualify`. Item
  operations bind the target UID and its derived CloudProvider scope.
- An authenticated principal lacking the applicable grant receives an audited
  `AUTHORIZATION_DENIED`/403 (VS0-CF-F16-08,15,18,37,111).
- Client attempts to write server-owned status or identity/metadata are audited
  `AUTHORIZATION_DENIED`/403 with `VS0_STATUS_FIELD_WRITE` or
  `VS0_SYSTEM_OWNED_FIELD_WRITE` (VS0-CF-F16-03,04,91,92). The status-write
  prohibition is inherited (VS0-WRITER-005/VS0-WRITER-006); clients never write
  status.

### 5.2 Safe denial and non-disclosure

- Phase-one safe resolution reads only the path UID needed to derive scope; an
  inaccessible target or backing reference returns safe `RESOURCE_NOT_FOUND`/404
  with `VS0_AUTHORIZATION_SAFE_DENIAL` and never discloses existence or details
  (VS0-CF-F16-09,17,19,66).
- LIST without a read grant returns audited 403 rather than an empty 200.
- Error responses never confirm the existence of an inaccessible resource
  (inherited FEATURE-0012 safe-denial semantics).

### 5.3 Privacy, redaction, and classification

- ExecutionTarget is `Provider-confidential`; `NormalizedTargetFactSet` and
  `TargetQualificationResult` are internal-only and never appear in any
  projection.
- The safe projection is closed to exactly the fields in REQ-F16-04; facts,
  results, observer identity, handles, and internal record links are omitted.
- No secret value is present on any path; no `SecretRef` extension is active
  (VS0-CF-F16-84). No credentials exist.

### 5.4 Audit (inherited FEATURE-0013 semantics, referenced not redefined)

- FEATURE-0016 consumes FEATURE-0013 AuditEvent atomic append-before-publication
  semantics by reference and does not redefine them.
- Audited events: successful create, completed qualification, retirement,
  maintenance enter/clear, authenticated authorization denial, and safe denial.
- A request-triggered required AuditEvent append failure returns
  `INTERNAL_ERROR`/500 with no AuditEvent, mutation, or completion, and never
  discloses the suppressed 403/404 (VS0-CF-F16-39,40,82,99..104,112).
- Expiry is the sole background exception: no request caller, no 500 return, no
  expiry AuditEvent until append succeeds, retrying once per injected-clock
  second (VS0-CF-F16-79).

### 5.5 Concurrency and optimistic concurrency (inherited FEATURE-0012)

- Qualify and retire require exactly one strong opaque `If-Match`; the ETag
  represents persisted target resource version only. F0016 has no conditional
  GET.
- Current-version comparison follows replay/reservation and precedes Maintenance
  state evaluation (VS0-CF-F16-74,85).
- The registered writer conflict for VS0-WRITER-006 is `STALE_RESOURCE_VERSION`.

### 5.6 Compatibility

- FEATURE-0016 reuses FEATURE-0012 Problem Details, ObjectMeta, TypedRef,
  Condition schemas, media semantics, and optimistic-concurrency conventions by
  reference; it defines no parallel envelope.
- FEATURE-0016 consumes FEATURE-0015 `CloudProviderParticipation` and
  `InfrastructureStack` as read-only prior-feature authority; it does not modify,
  extend, or reinterpret any FEATURE-0015 document, schema, writer, state, or
  conformance case.
- ADH-2026-058 atomically replaced the placeholder `VS0-SCHEMA-015..017` and
  `VS0-STATE-004`; there is no alias, migration path, or legacy compatibility for
  `status.availability`, `Draining`, or publicly projected `Qualifying`.

### 5.7 Operational

- All FEATURE-0016 state is in-process only. Process restart clears
  F0016-owned in-memory targets, records, idempotency reservations, and pending
  work, but never inherited FEATURE-0013 AuditEvents, and makes no external
  observer call (VS0-CF-F16-46).
- Idempotency retention is 24 hours and 10,000 completed entries in-process;
  completed entries evict at expiry then earliest-expiry/lexical namespace at
  capacity; InFlight never evicts.
- Service shutdown ordering: stop expiry/Maintenance trigger acceptance → abort
  in-flight qualification and idempotency reservations and wake waiters → stop
  HTTP serving (VS0-CF-F16-77).

## 6. Edge cases

| Edge case | Observable outcome | Conformance |
|-----------|--------------------|-------------|
| Malformed JSON, duplicate top-level member, oversized body | 400 MALFORMED_REQUEST / DUPLICATE_FIELD / REQUEST_TOO_LARGE; no publication | VS0-CF-F16-51,52,53 |
| Missing required create field, invalid typed reference shape, non-`synthetic-iaas` targetClass | 422 VALIDATION_FAILED; no publication | VS0-CF-F16-54,55,56 |
| Empty, malformed, or overlong Idempotency-Key | 400 MALFORMED_REQUEST; no publication | VS0-CF-F16-58,59,60 |
| Collection create with an If-Match header | 400 MALFORMED_REQUEST; no publication | VS0-CF-F16-105 |
| Any route with a query parameter (`?x=1`) | 400 MALFORMED_REQUEST before safe resolution | VS0-CF-F16-106 |
| Malformed canonical UID segment on item GET/qualify/retire | 400 MALFORMED_REQUEST before safe resolution | VS0-CF-F16-107,109,110 |
| Malformed/weak/wildcard/multiple/stale If-Match on actions | 412 STALE_RESOURCE_VERSION; header-form before body classification | VS0-CF-F16-70..74 |
| Qualify from Qualified with malformed/duplicate/missing fact while fences current | 500 INTERNAL_ERROR; reservation aborted; no publication | VS0-CF-F16-23,67,68,69,122 |
| Stale InfrastructureStack generation or viability fingerprint at commit | 412 STALE_RESOURCE_VERSION + VS0_TARGET_EPOCH_STALE / VS0_EXECUTION_TARGET_VIABILITY_STALE | VS0-CF-F16-29,30 |
| Maintenance epoch changes because Maintenance wins during qualification | 412 STALE_RESOURCE_VERSION + VS0_TARGET_EPOCH_STALE; winning transition audited | VS0-CF-F16-75 |
| Retire or Maintenance wins while qualify in flight | 409 CONFLICT + VS0_TARGET_RETIRED / VS0_TARGET_MAINTENANCE; winning transition audited | VS0-CF-F16-88,89 |
| Different-key qualify while Qualifying reservation in flight | 409 CONFLICT + VS0_TARGET_QUALIFICATION_IN_PROGRESS | VS0-CF-F16-28 |
| Owner panic or cancellation after reservation | Reservation aborted; waiters woken; cancelled waiter detaches | VS0-CF-F16-35,78 |
| Same-key replay while requester loses grant / target no longer accessible | Current 403 / safe 404 with one redacted AuditEvent; no stored-result disclosure | VS0-CF-F16-32,47 |
| Same principal/key/digest on another target or another derived scope | Independent reservation; never replays the other outcome | VS0-CF-F16-49,50 |
| Completed record expiry or capacity eviction | Record unavailable; next valid request processed as new; no InFlight eviction | VS0-CF-F16-34,76 |
| Fresh Qualified backing becomes non-viable / returns viable before fact expiry | GET projects Unavailable / Available with unchanged ETag, no mutation or AuditEvent | VS0-CF-F16-119,120 |
| Stale fenced maintenance enter/clear trigger | No state change, AuditEvent, or record publication | VS0-CF-F16-80,83 |
| Create name reserved by a Retired target | 409 ALREADY_EXISTS; tuple released but name/scope reservation retained | VS0-CF-F16-98 |

## 7. Non-goals and adjacent-feature exclusions

Active requirements and acceptance cases refer to this section for prohibited
concepts and never repeat the prohibited names elsewhere.

### 7.1 Feature non-goals (Phase 2R)

FEATURE-0016 does not introduce or activate: a real adapter, provider-native
type, raw credential, external call, external persistence, placement, execution,
plugin, customer target API, maintenance-notice resource, or IAM resource. Real
CloudProvider/Kubernetes/OpenShift/PostgreSQL calls, real adapter credentials,
real placement, real provisioning, real plugin execution, and any customer-facing
target projection are out of scope.

### 7.2 Removed placeholder concepts (must never reappear as active behavior)

- Persisted `status.availability` on ExecutionTarget (replaced by response-only
  `effectiveAvailability`).
- `Draining` as an active lifecycle or availability value.
- Publicly projected `Qualifying` as a resource status value.
- `spec.participationRef` (renamed to `spec.cloudProviderParticipationRef`).
- `status.conditions` / any public `conditions` member on ExecutionTarget.
- `adapterAuthorityRef` and any `SecretRef` extension field.
- Any alias, migration path, or legacy compatibility for the above.

### 7.3 Superseded concepts (DEC-0042)

- `ResourcePool` must never be introduced or activated as FEATURE-0016 behavior; it is prohibited as active FEATURE-0016 behavior, superseded by DEC-0042.
- `ProviderCapability` must never be introduced or activated as FEATURE-0016 behavior; it is prohibited as active FEATURE-0016 behavior, superseded by DEC-0042.
- A generic Provider as a combined owner/operator concept must never be introduced or activated as FEATURE-0016 behavior; it is prohibited as active FEATURE-0016 behavior, superseded by DEC-0042.
- Any ExecutionTarget subtype must never be introduced or activated as FEATURE-0016 behavior; it is prohibited as active FEATURE-0016 behavior, superseded by DEC-0042.

### 7.4 Adjacent-feature exclusions (enforceable non-goals)

| Excluded element | Owner | Negative acceptance |
|------------------|-------|---------------------|
| `PolicyEvaluationRequest`, `PolicyEvaluationResult`, `PolicyEngineAdapter` | FEATURE-0017 | No F0016 route, field, or record introduces or activates these |
| `SovereigntyProfile`, `SovereigntyFactSet`, `EvidenceRecord` | FEATURE-0019 | Not referenced by any F0016 create field, projection, or record |
| `ServiceRegion`, `spec.executionTargetRefs`, `ServiceRegion.status.availability` | FEATURE-0022 | F0016 is only a declared prerequisite; it does not introduce or activate either field |
| `DecisionProfiles`, `ServicePlacement`, `PlacementDecision`, customer target projection | FEATURE-0023 | No customer-facing target projection exists in F0016 |
| `PluginExecution`, `ServiceDeploymentPlan`, real realization/plugin taxonomy, Crossplane dependency | FEATURE-0024 | F0024/Phase 3 owns realization; Crossplane is future-only; proven by VS0-CF-F16-46 |

### 7.5 Negative acceptance cases

- No real adapter, credential, external call, Crossplane dependency, placement,
  provisioning, plugin execution, or customer IaaS exposure exists in active
  F0016 behavior (AC-F16-12; VS0-CF-F16-46).
- Only the eight closed local violation codes are introduced; no other
  F0016-local violation code exists (AC-F16-11).
- No PATCH, PUT, DELETE, HEAD (as a public route), watch, filter, pagination,
  fact, result, or adapter-selection route exists (REQ-F16-03; PUT/HEAD proven
  transport-only by VS0-CF-F16-44,113..118).

## 8. Architecture/decision/risk traceability by exact ID

### 8.1 Decisions and handoffs

| Authority | Role in FEATURE-0016 |
|-----------|----------------------|
| DEC-0036 | Implementation-neutral adapter boundary (no native types, raw credentials, or vendor errors) |
| DEC-0042 | Canonical cloud model; superseded `ResourcePool`/`ProviderCapability`/generic Provider |
| DEC-0057 | ExecutionTarget qualification/availability semantics |
| ADH-2026-025, ADH-2026-040, ADH-2026-042, ADH-2026-045 | Prior controlling handoffs and delegation of the ExecutionTarget contract to FEATURE-0016 |
| ADH-2026-058 | Executable contract closure; replaced placeholder VS0-SCHEMA-015..017 and VS0-STATE-004 (the sole current semantic authority) |

### 8.2 Registry IDs consumed

- Schemas: `VS0-SCHEMA-015` (ExecutionTarget), `VS0-SCHEMA-016`
  (NormalizedTargetFactSet), `VS0-SCHEMA-017` (TargetQualificationResult).
- Writer: `VS0-WRITER-006` (`ExecutionTargetLifecycleService`, sole committer;
  conflict `STALE_RESOURCE_VERSION`).
- State machine: `VS0-STATE-004` (Active/Retired × Unqualified/Qualified/
  Rejected/Indeterminate; Qualifying is an in-flight reservation only).
- Conformance: `VS0-CF-F16-01..99`, `VS0-CF-F16-100..122` (all FEATURE-0016-local).

### 8.3 Inherited error contract (FEATURE-0012 top-level Problem codes)

Top-level Problem codes and their HTTP/URN mappings are owned by FEATURE-0012
(`docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md`)
and are registered in
`docs/architecture/vertical-slices/VS-000-contract-registry.yaml`'s
`problemCodes` section; they are referenced here, not redefined or republished.
FEATURE-0016 introduces no top-level Problem code; its eight new codes appear
only in `violations[]` (REQ-F16-10). HTTP status codes are transport metadata,
never Problem `code` values, and no `429`/quota code is used.

### 8.4 Formal-proof references (owned by architecture, referenced here)

- `docs/formal/feature-0016/ExecutionTargetLifecycle.tla`: Retired terminality,
  the effectiveAvailability truth table, Maintenance/Retirement priority over
  expiry, and non-public Qualifying.
- `docs/formal/feature-0016/ExecutionTargetFence.tla`: qualification conclusion
  and replay completion require a successful AuditEvent append; stale/cancelled/
  failed proposals cannot mutate the target or become replayable completions.

### 8.5 Risk-relevant invariants

| Invariant | Enforced by |
|-----------|-------------|
| Clients never write status; sole committer is the lifecycle service | VS0-WRITER-006; VS0-CF-F16-03,04,91,92 |
| No external effect, credential, or persistence | REQ-F16-11; AC-F16-12; VS0-CF-F16-46 |
| Audit-before-publication; no partial mutation on append failure | §5.4; VS0-CF-F16-39,40,82,99..104 |
| Safe denial never discloses inaccessible resources | §5.2; VS0-CF-F16-09,17,19,66 |
| Closed five-route surface; transport guard has no side effect | REQ-F16-03; VS0-CF-F16-44,113..118 |

## 9. Design questions explicitly delegated by architecture

The architecture boundary (§2.2 DESIGN-Delegated Mechanics and ADH-2026-058
downstream admission) delegates only the following to the design stage. These
are mechanics, not semantics; they must preserve every observable behavior above.

1. In-memory registry data structures for `ExecutionTarget`,
   `NormalizedTargetFactSet`, and `TargetQualificationResult`.
2. Handler wiring and router setup behind the fixed pre-ServeMux guard.
3. Idempotency reservation/waiter internal data structures.
4. Injected-clock scheduling internals for expiry retry.

No semantic or ownership question is delegated. Fields, routes, states, error
codes, audit effects, replay classes, and projections are fully specified by the
authorities above and are not open design questions.

## 10. Completeness and unresolved-decision report

### 10.1 Coverage

- Canonical requirement ledger: all 11 approved REQ rows copied verbatim; each
  has exactly one normative detail heading in approved order (§4.1).
- Canonical acceptance ledger: all 12 approved AC rows copied verbatim; each is
  mapped to an acceptance scenario and exact conformance IDs (§4.2).
- Exact conformance semantics ledger: all 122 cases in the FEATURE-0016 local
  conformance suite copied with ID, owner, inputs, expectedState,
  expectedError, expectedViolation (where registered), expectedSideEffects,
  and gate (§4.3).
- Owned schema (VS0-SCHEMA-015..017), writer (VS0-WRITER-006), and state
  (VS0-STATE-004) IDs are referenced without redefinition.

### 10.2 Boundary checks

- No design choice (data structures, locks, package layout, scheduler
  implementation) is decided here.
- No downstream/future-feature semantic (FEATURE-0017/0019/0022/0023/0024) is
  imported; all appear only as enforceable exclusions in §7.
- No inherited FEATURE-0012/0013/0015 contract is copied, renamed, specialized,
  or re-owned; each is consumed by reference.
- No shared/downstream conformance case (`VS0-CF-F10`, `VS0-CF-X03`,
  `VS0-CF-HP01`) is used as FEATURE-0016-local proof.
- No orphan IDs: every referenced ID exists in the registry.

### 10.3 Unresolved decisions

None. No `ARCHITECTURE_DECISION_REQUIRED`, `REQUIREMENT_CLARIFICATION_REQUIRED`,
`BOUNDARY_CHANGE_REQUIRED`, `DEPENDENCY_APPROVAL_REQUIRED`, or
`SECURITY_REVIEW_REQUIRED` condition was encountered. ADH-2026-058 supplies one
complete, closed, executable contract, and every requirement maps to existing
registry entries.

---

Model Execution Report:
- Tool: kiro
- Stage or task: FEATURE-0016 Requirements (adapter-boundary-and-executiontarget-qualification/requirements.md)
- Recommended priority list: Architecture-heavy (claude-opus-4.8, effort high); Fallback (claude-sonnet-4.5, effort medium); (third model not specified in prompt — only two labels provided)
- Selected model: claude-opus-4.8
- Effort/reasoning setting: high
- Fallback used: no
- Fallback reason: none
STAGE_STATUS: COMPLETE