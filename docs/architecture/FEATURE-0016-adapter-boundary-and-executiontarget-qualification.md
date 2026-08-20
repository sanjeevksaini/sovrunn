# FEATURE-0016 Architecture Boundary: Adapter Boundary and ExecutionTarget Qualification

| Field | Value |
|-------|-------|
| Status | Approved boundary (closed under ADH-2026-058) |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Controlling Decisions | DEC-0036, DEC-0042, DEC-0057 |
| Controlling Handoffs | ADH-2026-025, ADH-2026-040, consolidated ADH-2026-042, ADH-2026-045 (delegation), ADH-2026-058 (executable contract closure) |
| Phase | 2R |
| Depends On | FEATURE-0011 (reuse), FEATURE-0012 (grammar/errors), FEATURE-0013 (decision/audit), FEATURE-0015 (CloudProviderParticipation, InfrastructureStack — read-only prior-feature authority) |

---

## 1. Purpose

Define the closed, executable architecture boundary for FEATURE-0016 so that
requirements, design, and tasks translate one complete observable contract
without inventing routes, fields, states, or error semantics. ADH-2026-058
atomically replaced the FEATURE-0016 placeholder definitions of
`VS0-SCHEMA-015..017` and `VS0-STATE-004` with the register below. This
document restates that approved register as the FEATURE-0016 feature
authority; it introduces no decision beyond ADH-2026-058.

---

## 2. FEATURE-0016 Owned Resources and Activation Boundaries

FEATURE-0016 owns, in full (both introduces and activates):

| Resource | Registry ID | Scope | Profile | Activation Boundary |
|----------|------------|-------|---------|---------------------|
| ExecutionTarget | VS0-SCHEMA-015 | CloudProvider | ManagedResource | identity, closed create fields, derived scope, full status (lifecycle, qualification, maintenanceEpoch, observedGeneration, factSetRef, qualificationResultRef); no persisted `status.availability` |
| NormalizedTargetFactSet | VS0-SCHEMA-016 | CloudProvider | ImmutableRecord | internal-only observation record; never projected |
| TargetQualificationResult | VS0-SCHEMA-017 | CloudProvider | ImmutableRecord | internal-only qualification conclusion; never projected |
| ExecutionTarget lifecycle/qualification/maintenance state machine | VS0-STATE-004 | — | — | Active/Retired x Unqualified/Qualified/Rejected/Indeterminate; Qualifying is a lifecycle-service-only in-flight reservation, never persisted or projected |

FEATURE-0016 owns the sole synthetic observer `sovrunn.synthetic-iaas-observer/v1`,
the `ExecutionTargetLifecycleService` (sole committer of ExecutionTarget status
and internal FactSet/Result records, VS0-WRITER-006), and the five-route HTTP
surface plus its pre-ServeMux transport-only method/path guard.

---

## 3. Field Ownership Boundaries

### 3.1 ExecutionTarget create-field boundary (ADH-2026-058 clause 1/F16-AD-02.1)

The create body accepts exactly `metadata.name`,
`spec.cloudProviderParticipationRef.uid`, `spec.infrastructureStackRef.uid`,
and immutable `spec.targetClass="synthetic-iaas"`. No separate `scopeRef` is
persisted or projected — scope is derived only from
`spec.cloudProviderParticipationRef`. The server alone assigns
`metadata.uid`, `metadata.resourceVersion`, `metadata.generation`,
`status.lifecycle`, `status.qualification`, `status.maintenanceEpoch`,
`status.observedGeneration`, `status.factSetRef`, and
`status.qualificationResultRef`. Unknown fields and `adapterAuthorityRef` are
rejected with `UNKNOWN_FIELD`/400; no `SecretRef` extension is active.

### 3.2 ServiceRegion dependency (unchanged from FEATURE-0015 delegation)

`ServiceRegion.spec.executionTargetRefs` and `ServiceRegion.status.availability`
remain wholly owned by FEATURE-0022; FEATURE-0016 is only a declared
prerequisite and does not introduce or activate either field.

### 3.3 Removed placeholder concepts

Per ADH-2026-058, the following placeholder concepts are removed and must
never reappear as active FEATURE-0016 behavior:

- Persisted `status.availability` (three/four-value enum) on ExecutionTarget.
- `Draining` as an active lifecycle or availability value.
- Publicly projected `Qualifying` as a resource status value.
- Any alias, migration path, or legacy compatibility behavior for the above.
- `spec.participationRef` (renamed to `spec.cloudProviderParticipationRef`).
- `status.conditions` on ExecutionTarget (F0016 defines no public `conditions` member).

---

## 4. Consolidated Canonical Contract (ADH-2026-058 subclauses 1–10)

### 4.1 Boundary and canonical resource

FEATURE-0016 owns one CloudProvider-confidential, CloudProvider-scoped
`ExecutionTarget` plus internal observation/result records. It qualifies
infrastructure-backed `synthetic-iaas` only. It does not place, provision,
realize, or expose IaaS objects; no adapter resource, credential, external
call, Crossplane dependency, or customer projection is active.

### 4.2 Safe reference access and viability

Safe access precedes graph validation. Create/qualify require a FEATURE-0015
effective-Active `CloudProviderParticipation` (accepted, with neither
suspension hold set) and an Active `InfrastructureStack` in the same
CloudProvider scope. Name is scope-unique for the process lifetime, including
after retirement (retirement releases only the tuple). At most one
non-Retired `(participation UID, stack UID, targetClass)` tuple is live. The
viability fingerprint is only the participation UID/effective-active state
and the stack UID/phase. Outcomes:

| Condition | Outcome |
|---|---|
| Inaccessible backing reference | Safe `RESOURCE_NOT_FOUND`/404 + `VS0_AUTHORIZATION_SAFE_DENIAL` |
| Authorized ineffective participation | `CONFLICT`/409 + `VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE` |
| Authorized non-Active stack | `CONFLICT`/409 + `VS0_EXECUTION_TARGET_STACK_UNAVAILABLE` |
| Authorized cross-CloudProvider graph | `VALIDATION_FAILED`/422 + `VS0_EXECUTION_TARGET_SCOPE_MISMATCH` |
| Duplicate name or live tuple | `ALREADY_EXISTS`/409 |

### 4.3 Closed HTTP surface (five routes, one transport guard)

Exactly five Go 1.22 `http.ServeMux` registrations exist:

| Method | Path |
|---|---|
| POST | `/apis/execution.sovrunn.io/v1alpha1/execution-targets` |
| GET | `/apis/execution.sovrunn.io/v1alpha1/execution-targets` |
| GET | `/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}` |
| POST | `/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify` |
| POST | `/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire` |

No PATCH, PUT, DELETE, HEAD, watch, filter, pagination, fact, result, or
adapter-selection route exists. A fixed pre-ServeMux F0016 method/path guard
runs before those registrations:

- HEAD on the collection path: transport-only 405, `Allow: GET, POST`.
- HEAD on the item path: transport-only 405, `Allow: GET`.
- HEAD on either action path: transport-only 405, `Allow: POST`.
- Every trailing-slash variant: transport-only 404 without redirect.
- Every other unmatched method on an exact F0016 path: transport-only 405.

The guard adds no public route or handler and causes no Problem body,
authentication, authorization, audit, idempotency, or lifecycle effect.
Every response, including a transport-only 404/405 and a replay, sends
`X-Content-Type-Options: nosniff`. Valid registered method/path pairs bypass
the guard and begin at authentication.

### 4.4 Authorization and safe denial

Create/retire use `executiontarget.write`; reads use `.read`; qualify uses
`.qualify`. Item operations bind the target UID and its derived CloudProvider
scope. LIST contains all accessible targets, including Retired, in ascending
`metadata.uid` order; no read grant returns audited 403 rather than empty 200.
Missing/invalid authentication is 401 with no audit. For item GET and either
action, phase-one safe resolution reads only the path UID needed to derive
its CloudProvider scope: an inaccessible target is safe 404; only then is
the target-bound grant evaluated, so an accessible target without that grant
is audited 403. Neither outcome discloses target details.

### 4.5 Request grammar and representation

Every POST requires `application/json` and exactly one `Idempotency-Key`.
Qualify/retire also require exactly one strong opaque `If-Match` and a
zero-byte body. Header-form validation and current-version comparison are
distinct steps. Create/item GET/qualify/retire return `ETag`; LIST does not.
Success uses `application/json`, errors `application/problem+json`, and
every response sends `X-Content-Type-Options: nosniff`. F0016 ignores
`Accept` and emits its fixed response media. No F0016 route accepts a query
parameter.

The safe ExecutionTarget JSON projection is exactly: `metadata.uid`,
`metadata.name`, `metadata.generation`, `metadata.resourceVersion`, the
three immutable spec fields, `status.lifecycle`, `status.qualification`,
and response-only `effectiveAvailability`. F0016 defines no public
`conditions` member. Facts, results, observer data, handles, and internal
record links are omitted. Create returns 201; item GET/qualify/retire
return 200; LIST returns an ascending JSON array of the projection. No
`Location` header is defined.

### 4.6 Executable request precedence

Closed transport method/path guard → authentication → method/media/header-form
validation → phase-one safe references → authorization → safe access → strict
classification → graph/reference validation → replay/reservation → current
If-Match version comparison → lifecycle state → atomic commit. Create alone
has client graph/reference fields and no If-Match; zero-byte actions alone
have If-Match and no client graph. A malformed If-Match header-form failure
therefore precedes its non-zero-body failure; current-version comparison
occurs after replay/reservation and before Maintenance state evaluation.

### 4.7 Idempotency and replay

The namespace is principal UID, route pattern, derived CloudProvider UID,
concrete action-target UID when applicable, and `Idempotency-Key`. The
digest is SHA-256 of recursively sorted allowed create JSON, or SHA-256 of
zero bytes for actions; `If-Match`, `Accept`, authentication, and
request/correlation IDs are excluded. Invalid/unknown/duplicate bodies are
never digested or reserved. Idempotency applies only to create, qualify, and
retire; GET and LIST neither require, inspect, reserve, wait on, read, nor
replay an idempotency record. Successful create/qualification/retirement
retain for 24 hours/10,000 entries only in-process; aborted outcomes do not
replay; restart clears F0016-owned in-memory state but never inherited
FEATURE-0013 AuditEvents. Completed entries evict at expiry, then
earliest-expiry/lexical namespace when capacity is exceeded; InFlight never
evicts. A waiter rechecks current authorization and safe access immediately
after waking; a changed digest is 409 `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`.
The same principal/key/digest on a different action target, or on a create
in another derived CloudProvider scope, forms an independent reservation.

### 4.8 Observer and fact contract

The only observer is `sovrunn.synthetic-iaas-observer/v1`: target-bound,
clock-driven, external-effect-free. Exactly once each of `compute.vm`,
`storage.block`, `storage.object`, and `network.private` is Supported,
Unsupported, or Unknown. Any Unsupported rejects; otherwise any Unknown is
Indeterminate; otherwise Qualified. A missing fixture and a fixture-declared
logical timeout are separate deterministic outcomes, each yielding four
Unknown facts at `now` through `now+60s`; neither waits on a network call.
A FactSet carries target ref, referenced InfrastructureStack generation,
maintenance epoch, viability fingerprint, observer ID/revision, fact
version, `observedAt`/`expiresAt`, and the four named truth values. A Result
carries target/fact-set refs, the same fences, profile version, outcome,
reason codes, and evaluation time. Duplicate/missing facts or malformed
provenance are `INTERNAL_ERROR`/500 observer faults and publish no result.
Fixed provenance: `observerID=sovrunn.synthetic-iaas-observer/v1`,
`profileVersion=synthetic-iaas/v1`, `factSchemaVersion=v1`.

### 4.9 Lifecycle, qualification, and effective availability

`ExecutionTargetLifecycleService` alone commits target state and internal
fact/result records. Lifecycle is Active/Retired; persisted public
qualification is Unqualified/Qualified/Rejected/Indeterminate; `Qualifying`
is a lifecycle-service-only in-flight reservation, never a projected
resource status. The persisted target has no `status.availability`: the
server returns response-only `effectiveAvailability` (Available,
Unavailable, Maintenance) computed from committed state, freshness,
maintenance, and viability. It has no independently persisted writer and
does not alter ETag.

Qualify is synchronous and may start only from Active
Unqualified/Qualified/Rejected/Indeterminate. It creates a Qualifying
reservation that leaves the committed target summary, current
FactSet/Result links, and resourceVersion/ETag unchanged until a completed
conclusion can be audited and published. Create persists Active/Unqualified
at maintenance epoch 0, sets `observedGeneration` to the referenced
InfrastructureStack generation, has no record refs, and projects
`effectiveAvailability=Unavailable`.

Every completed conclusion atomically replaces both record refs, sets
`observedGeneration` to its captured InfrastructureStack generation,
advances resourceVersion/ETag, and appends the required AuditEvent before
publication. Malformed observer output, owner cancellation, required
AuditEvent append failure, and every stale captured fence abort the
reservation without target mutation, AuditEvent, or completion. A stale
qualification proposal never suppresses an independently committed winning
Maintenance or Retirement transition.

Maintenance entered during an in-flight Qualifying reservation aborts that
reservation, clears current links, and persists Unqualified before
projecting Maintenance; otherwise it preserves the last completed
conclusion while projecting Maintenance. Clear persists Active/Unqualified
and projects Unavailable. Retire is allowed from every Active combination,
including while Qualifying is in flight, and persists Retired/Unqualified
while projecting Unavailable.

Every successful maintenance enter or clear increments maintenance epoch and
resourceVersion/ETag, clears both record links, and leaves no pre-fence
record current; clear also sets qualification Unqualified. The lifecycle
service owns a target-bound, in-process current-Maintenance marker (target
UID, maintenance epoch, active flag). It is never client-writable or
projected and is cleared on restart. `effectiveAvailability` is Maintenance
exactly while that marker is active.

`effectiveAvailability` is ordered: Retired projects Unavailable; otherwise
Active with current Maintenance projects Maintenance; otherwise only Active
Qualified with fresh current facts and viable backing projects Available;
every other Active combination projects Unavailable. A viability-only
change does not mutate the target or change its ETag.

If Retire wins while a qualification reservation is in flight, the caller
receives `CONFLICT`/409 + `VS0_TARGET_RETIRED`. If Maintenance entry wins,
the qualifying caller receives `CONFLICT`/409 + `VS0_TARGET_MAINTENANCE`.
Both abort qualification without result, completion, or qualification
AuditEvent; the winning transition itself is audited.

### 4.10 Expiry, maintenance trigger, retirement, shutdown, and formal proof

Expiry is injected-clock driven; Retired and Maintenance outrank it.
Maintenance enters/clears only through a deterministic, target-bound
synthetic-observer fixture trigger — never an HTTP route or generic event
bus. Retire is the sole lifecycle exit and releases the tuple.

Audit successful create, completed qualification, retirement, maintenance
enter/clear, authenticated authorization denial, and safe denial.
Request-triggered required append failure is `INTERNAL_ERROR`/500 with no
AuditEvent, mutation, or completion publication. Expiry is the sole
background exception: no request caller, no 500 return, no expiry
AuditEvent until append succeeds; a failed append retries once per
injected-clock second until one append succeeds. Expiry leaves the last
qualification and ETag unchanged, projects Unavailable, and writes one
expiry AuditEvent.

Service shutdown first stops acceptance of expiry and Maintenance triggers,
then aborts in-flight qualification and idempotency reservations and wakes
waiters, and only then stops HTTP serving. A trigger accepted after that
stop boundary must neither mutate target state nor append an AuditEvent.

The lifecycle TLC model (`docs/formal/feature-0016/ExecutionTargetLifecycle.tla`)
proves Retired terminality, the effectiveAvailability truth table,
Maintenance/Retirement priority over expiry, and that Qualifying is never
publicly persisted. The fence/publication TLC model
(`docs/formal/feature-0016/ExecutionTargetFence.tla`) proves that a
qualification conclusion and replay completion require a successful
AuditEvent append, while stale, cancelled, and failed proposals cannot
mutate the target or become replayable completions.

### 4.11 Closed local violation set

The only new local violations are `VS0_TARGET_RETIRED`,
`VS0_TARGET_MAINTENANCE`, `VS0_TARGET_QUALIFICATION_IN_PROGRESS`,
`VS0_TARGET_EPOCH_STALE`, `VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE`,
`VS0_EXECUTION_TARGET_STACK_UNAVAILABLE`,
`VS0_EXECUTION_TARGET_SCOPE_MISMATCH`, and
`VS0_EXECUTION_TARGET_VIABILITY_STALE`. FEATURE-0012 owns top-level Problems.

### 4.12 Future boundary

F0024/Phase 3 owns realization/plugin taxonomy. Crossplane is future-only.
Requirements, design, tasks, and Cursor may not choose an omitted observable
behavior, authority, proof meaning, or implementation scope beyond this
document and ADH-2026-058.

---

## 5. Preserve FEATURE-0015 as Read-Only Prior-Feature Authority

FEATURE-0016 consumes FEATURE-0015's `CloudProviderParticipation` and
`InfrastructureStack` contracts by reference only. FEATURE-0016 does not
modify, extend, or reinterpret any FEATURE-0015 document, schema, writer,
state, or conformance case.

---

## 6. Local Conformance (VS0-CF-F16-01..122)

FEATURE-0016 owns an exact, closed, feature-local conformance suite of 122
cases registered in `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`
and cross-referenced in `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`.
No shared or downstream case (`VS0-CF-F10`, `VS0-CF-X03`, `VS0-CF-HP01`) is
counted as F0016-local proof; they remain cross-feature/shared references
only. See ADH-2026-058 §2 for the complete annex and §3 for the
decision-to-authority/proof matrix.

---

## 7. Non-Goals

FEATURE-0016 does not introduce: a real adapter, provider-native type, raw
credential, external call, external persistence, placement, execution,
plugin, customer target API, maintenance-notice resource, IAM resource,
`ResourcePool`, `ProviderCapability`, a generic Provider, or a target
subtype. It does not choose Go data structures, locks, package layout, or
scheduler implementation — those remain design/task mechanics that must
preserve every observable behavior in this document.
