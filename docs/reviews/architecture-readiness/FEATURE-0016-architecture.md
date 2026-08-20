---
doc_type: proposed_feature_architecture
feature: FEATURE-0016
title: Adapter Boundary and ExecutionTarget Qualification
status: proposed-not-approved-architecture
baseline: ARCH-2026.08-PHASE2R-CANONICAL
decision_input: FEATURE-0016-architecture-starting-dossier.md
approval_required: founder-approved-ADH
---

# FEATURE-0016 Architecture: Adapter Boundary and ExecutionTarget Qualification

## 1. Status and use

This is a proposed architecture draft for discussion and Architecture Decision
Handoff preparation. It does not alter the approved baseline, registry,
feature authority, requirements, design, tasks, or implementation scope.

Its recommendations become authoritative only after a founder-approved ADH is
applied through the established architecture-update process.

## 2. Purpose

FEATURE-0016 establishes the smallest canonical boundary through which
Sovrunn can determine whether an internal realization environment is usable
for later PaaS service execution.

It qualifies an `ExecutionTarget`; it does not provision resources, schedule
capacity, place a customer service, execute a plugin, or expose native
infrastructure concepts to a customer.

## 3. Approved inputs consumed by reference

- FEATURE-0015 supplies CloudPlatform, CloudProvider,
  CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, and
  InfrastructureStack. FEATURE-0016 neither changes nor re-owns them.
- FEATURE-0012 supplies API grammar, typed references, Problem Details,
  validation conventions, and optimistic-concurrency conventions.
- FEATURE-0013 supplies AuditEvent and DecisionRecord contracts.
- The current canonical model delegates ExecutionTarget, normalized target
  facts, target qualification results, and target lifecycle to FEATURE-0016.
- Phase 2R remains deterministic, in-memory, and side-effect free.

## 4. Proposed boundary

### 4.1 What FEATURE-0016 owns

- `ExecutionTarget` identity, internal/operator-facing schema and lifecycle.
- `NormalizedTargetFactSet` as a normalized, redacted observation record.
- `TargetQualificationResult` as an immutable, attributable qualification
  conclusion.
- A narrow adapter boundary used to observe target facts.
- Qualification, availability, maintenance-epoch fencing, and a deterministic
  fake adapter.
- F0016-local route, state, authorization, audit, concurrency, and
  no-external-effect proof.

### 4.2 What FEATURE-0016 must not own

- VM, volume, object bucket, VPC, subnet, security group, Kubernetes
  namespace, CloudProvider account, or IaaS management-plane resource.
- A ResourcePool, CloudProvider-wide capability inventory, capacity scheduler,
  placement decision, customer service, customer target projection, or billing
  concept.
- Policy evaluation, entitlement, quota, PluginExecution, service
  provisioning, real adapter calls, external secret storage, or federation.

### 4.3 Minimum active Phase 2R realization

The recommended active F0016 realization is one **synthetic IaaS target**.
It is backed by an existing `InfrastructureStack`, inherits the resolved
CloudProvider boundary, and deterministically reports this minimum normalized
capability profile:

- VM compute;
- block storage;
- object storage; and
- private networking.

The fake creates none of these resources and performs no network, CloudProvider,
Kubernetes, secret-manager, or external API call. Its purpose is to prove the
same qualification, freshness, fencing, redaction, and denial contracts that a
future real IaaS adapter must satisfy.

## 5. Consolidated proposed architecture decision register

This register and its subclauses are the sole proposed semantic authority for
FEATURE-0016. ADH-2026-058 is a bounded application and proof package: it may
identify authority changes, annex proof, and Kiro actions, but must not create
a competing semantic rule.

Every row below is a proposed resolution, not approved architecture. An ADH
must approve or explicitly defer every row before requirements generation.
Neither requirements nor design may choose an omitted observable behavior.

| ID | Canonical decision | Proposed resolution | Example |
|---|---|---|---|
| F16-AD-01 | Boundary and canonical resource | F0016 owns one CloudProvider-scoped, CloudProvider-confidential `ExecutionTarget` plus internal observation/result records. It qualifies infrastructure-backed `synthetic-iaas`; it does not place, provision, realize, or expose IaaS objects. | A qualified target cannot create a VM. |
| F16-AD-02 | Create contract and identity | Create accepts only `metadata.name`, `spec.cloudProviderParticipationRef.uid`, `spec.infrastructureStackRef.uid`, and immutable `spec.targetClass=synthetic-iaas`. Scope is derived only from participation and is not persisted as a separate resource field; name is scope-unique for the process lifetime; at most one non-retired target may use a participation/stack/class tuple. | Retirement releases a tuple but not its name. |
| F16-AD-03 | References and viability | Safe access precedes graph validation. Create/qualify require active participation and Active InfrastructureStack in the same CloudProvider scope. Qualification records a narrow viability fingerprint; only its listed viability changes invalidate qualification. | A visible cross-CloudProvider stack is 422; suspension makes a target effectively unavailable. |
| F16-AD-04 | Closed HTTP surface | Exactly five explicit Go 1.22 `http.ServeMux` registrations exist: collection POST/GET, item GET, qualify POST, retire POST. A fixed F0016 pre-ServeMux method/path guard rejects unsupported HEAD and trailing-slash forms before mux dispatch; it does not add a public route or handler. No PATCH, PUT, DELETE, HEAD, watch, filter, pagination, fact/result, or adapter-selection route exists. | `HEAD /execution-targets/{uid}` is transport-only 405. |
| F16-AD-05 | Authorization and safe denial | Create/retire use `executiontarget.write`; GET/LIST use `.read`; qualify uses `.qualify`. Item access additionally binds target UID and derived CloudProvider scope. | A read principal cannot retire; inaccessible objects safe-404. |
| F16-AD-06 | Request and representation grammar | POST requires `application/json` and one Idempotency-Key; qualify/retire also require one strong opaque If-Match ETag and a zero-byte body; create rejects If-Match. Header-form validation and current-version comparison are distinct. The five exact non-trailing-slash paths accept no query parameters. Create, item GET, qualify, and retire return ETag; LIST does not. | `{}` on qualify, `If-Match: *`, and `?limit=10` are rejected. |
| F16-AD-07 | Executable request precedence | Order: closed transport method/path guard; authentication; method/media/header-form validation; phase-one safe references; authorization; safe access; strict classification; graph/reference validation; replay/reservation; current If-Match version comparison; lifecycle state; atomic commit. The guard is limited to unmatched method/path forms; valid registrations begin at authentication. Create has graph/reference validation and no If-Match; zero-byte actions have If-Match and no client graph. | A malformed If-Match wins over a non-zero body; a valid stale If-Match wins over Maintenance. |
| F16-AD-08 | Idempotency and replay | Namespace is principal, route, derived CloudProvider scope UID, concrete target UID for actions, and Idempotency-Key. The separately stored digest is SHA-256 of canonical allowed create JSON, or SHA-256 of zero bytes for actions; a changed digest in the same namespace conflicts. If-Match and all headers are excluded. Retain completed responses 24 hours/10,000 entries only during one process lifetime. Only successful create, completed qualification, and retirement replay; abort outcomes do not. | A Target A key cannot replay for Target B; a changed body under the same create key conflicts. |
| F16-AD-09 | Adapter and authority boundary | No Adapter resource, vendor type, credential field, SecretRef extension, external call, or native error is active. Server configuration binds the one synthetic observer to the active class and supplies only a non-secret handle. | `adapterAuthorityRef` is absent from schema and rejected. |
| F16-AD-10 | Observation and deterministic fake | Observation receives target/fence context and returns normalized facts or classified outcomes. The sole observer is `sovrunn.synthetic-iaas-observer/v1`; fixtures are target-bound, clock-driven, and make no external call. | A missing fixture returns four Unknown facts, never a qualified default. |
| F16-AD-11 | Fact contract and qualification | Exactly four facts occur once each: compute.vm, storage.block, storage.object, network.private. Truth is Supported/Unsupported/Unknown. Any Unsupported rejects; otherwise any Unknown is Indeterminate; otherwise Qualified. | Unsupported block plus Unknown object storage is Rejected. |
| F16-AD-12 | Provenance and time validity | Fact sets bind target, referenced InfrastructureStack generation, maintenance epoch, fixed observer/profile/fact versions, observedAt, expiresAt, and viability fingerprint. They require `observedAt <= now < expiresAt`; missing-fixture and fixture-declared logical-timeout Unknown facts use now and now+60 seconds. | A reversed timestamp is malformed observer output (500). |
| F16-AD-13 | Internal records and restart | Facts/results are immutable, internal, redacted from projections, and retained only as the current pair for each non-retired target. Restart clears F0016-owned in-memory targets, records, idempotency reservations, and pending work; it does not delete inherited FEATURE-0013 AuditEvents. | Retirement removes the pair but retains a redacted AuditEvent summary. |
| F16-AD-14 | Single lifecycle service | `ExecutionTargetLifecycleService` is the only committer of persisted target state and fact/result records. Handlers submit commands; observer/engine propose; freshness, maintenance, and viability submit triggers. | The engine cannot write target status directly. |
| F16-AD-15 | Lifecycle and availability model | Lifecycle is Active/Retired. Persisted public qualification is Unqualified/Qualified/Rejected/Indeterminate; Qualifying is a lifecycle-service-only in-flight reservation, never a projected resource status. Public availability is an effective lifecycle-service response projection of committed state, freshness, maintenance, and viability; it is not separately persisted. Retired is Unavailable; Maintenance is Maintenance; only Active/Qualified/fresh/viable is Available; all other Active combinations are Unavailable. | Expired facts can make GET return Unavailable with unchanged ETag. |
| F16-AD-16 | Qualification completion and rollback | Qualify is synchronous. Starting it creates an internal Qualifying reservation while retaining the last committed target summary/links and ETag. Valid outcomes atomically commit Qualified, Rejected, or Indeterminate after required AuditEvent append succeeds. Malformed, cancelled, audit-failed, and stale work aborts with no target mutation; Retire or Maintenance remains the winner when it transitions first. | Retire during qualification returns `VS0_TARGET_RETIRED`; the qualifier cannot restore Available. |
| F16-AD-17 | Freshness and expiry | Every current fact set expires. Retired and Maintenance outrank expiry; otherwise expiry preserves the last completed qualification and projects effectiveAvailability Unavailable. An expiry audit failure keeps reads unavailable and retries on the next one-second injected-clock tick; restart removes pending retry. | A 10:00 expiry retry occurs at 10:00:01 after audit failure. |
| F16-AD-18 | Maintenance and trigger priority | The only F0016 maintenance ingress is a deterministic, target-bound synthetic-observer fixture trigger containing target, expected InfrastructureStack generation/epoch, operation, system identity, and correlation; it is not an HTTP route or generic event bus. Priority is Retired, Maintenance, expiry, qualification completion. Retired maintenance triggers are ignored. | A fixture enters maintenance at 10:00, advances epoch 7 to 8, and discards epoch-7 work. |
| F16-AD-19 | Retirement | Retire is the sole lifecycle exit and is allowed for Active targets, including while a Qualifying reservation is in flight, with current If-Match. It releases the tuple. New-key retire of Retired is 409 `VS0_TARGET_RETIRED`; same-key replay returns 200. | A replacement can use the retired target's tuple. |
| F16-AD-20 | Stale classification and concurrency | Invalid returned facts/provenance are 500. A stale qualification proposal whose referenced InfrastructureStack generation, maintenance epoch, or viability fingerprint changed at commit is 412 with the matching stale violation and no qualification mutation or conclusion publication. An independently committed winning Maintenance or Retirement transition remains valid and audited. Same-key in-flight calls wait; different-key Qualifying calls conflict. | Participation suspension during observation returns 412 `VS0_EXECUTION_TARGET_VIABILITY_STALE`. |
| F16-AD-21 | Audit boundary | Audit successful create, completed qualification, retirement, maintenance enter/clear, authenticated authorization denial, and safe denial, including replay-time denials. Request-triggered append failure yields 500 and no committed mutation/completion; denial audit failure overrides the original 403/404. Expiry is the sole background retry exception. | A failed replay safe-denial audit returns 500 without target disclosure. |
| F16-AD-22 | F0016 Problem/violation set | Reuse FEATURE-0012 Problem Details. F0016 adds retired, maintenance, qualification-in-progress, epoch-stale, participation-unavailable, stack-unavailable, scope-mismatch, and viability-stale violations only. | A visible inactive stack returns 409 with its local violation. |
| F16-AD-23 | Privacy and observability | Operator reads expose only the closed safe projection: metadata.uid/name/generation/resourceVersion, immutable spec refs/class, lifecycle, qualification, and effective availability. They omit public conditions, facts, results, observer identity/revision, handles, and internal refs. | GET shows Qualified but not a raw fact payload. |
| F16-AD-24 | Future boundary | Realization/plugin taxonomy is F0024/Phase 3. Crossplane is a future realization/observation candidate only, never an F0016 dependency or model owner. | No Crossplane CRD or credential enters F0016. |
| F16-AD-25 | Local proof and formality | F0016 owns local conformance for every route, denial, transition, audit/replay, fake, and stale outcome. TLC proves lifecycle/freshness/retirement and stale non-publication. | A model proves old work cannot restore availability. |
| F16-AD-26 | Approved authority/proof annex | Before founder approval, ADH-058 must contain one matrix mapping every canonical decision to authority, owner/trigger, local proof ID, and readiness assertion, plus a complete five-route outcome annex. | No Kiro step may choose an omitted outcome or proof meaning. |
| F16-AD-27 | Requirements and design admission | Requirements begin only after architecture checks, formal models, context headroom, and the annex pass. Design requires a fail-closed five-route executability audit. | A missing media or precedence outcome blocks design. |
| F16-AD-28 | Tasks and implementation admission | Tasks require a source-derived writable-path/dependency inventory and acyclic graph. Cursor cannot repair authority gaps or silently broaden scope. | A missing fixture or lifecycle path blocks task approval. |

### 5.1 Binding subclauses of the canonical decisions

These are subclauses of the rows above, not additional decisions.

#### F16-AD-02.1 — identity and tuple boundary

The create body is exactly `metadata.name`,
`spec.cloudProviderParticipationRef.uid`,
`spec.infrastructureStackRef.uid`, and
`spec.targetClass="synthetic-iaas"`; both references are UID-pinned. No other
metadata member or spec member is client-writable. The server assigns
`metadata.uid`, `metadata.resourceVersion`, `metadata.generation`,
`status.lifecycle`, `status.qualification`, `status.maintenanceEpoch`,
`status.observedGeneration`, `status.factSetRef`, and
`status.qualificationResultRef`. No separate `scopeRef` is persisted or
projected: scope is always derived from the participation reference. There is
at most one non-Retired target for (participation UID, InfrastructureStack UID,
targetClass). Retirement releases the tuple but not the process-lifetime
scope/name reservation. A duplicate name or live tuple returns
`ALREADY_EXISTS`/409.

#### F16-AD-03.1 — effective backing and viability fingerprint

“Active participation” means FEATURE-0015 **effective Active**: accepted,
with neither independent suspension hold set. The viability fingerprint is
only the participation UID/effective-active value and the InfrastructureStack
UID/`status.phase`; unrelated descriptions or resource-version changes do
not invalidate a qualification. An inaccessible backing reference is safe
`RESOURCE_NOT_FOUND`/404 with `VS0_AUTHORIZATION_SAFE_DENIAL`; an
authorized inactive participation is `CONFLICT`/409 with
`VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE`; an authorized inactive
stack is `CONFLICT`/409 with
`VS0_EXECUTION_TARGET_STACK_UNAVAILABLE`; a same-visible but
cross-CloudProvider graph is `VALIDATION_FAILED`/422 with
`VS0_EXECUTION_TARGET_SCOPE_MISMATCH`.

#### F16-AD-04/05.1 — read membership and denial

LIST contains every target, including Retired targets, readable through the
caller’s derived CloudProvider grants and sorts the resulting set by ascending
`metadata.uid`. A caller with no applicable read grant receives audited
`AUTHORIZATION_DENIED`/403, not an empty list. An inaccessible direct GET
is safe `RESOURCE_NOT_FOUND`/404.

#### F16-AD-04/05.2 — exact paths and grant-denial coverage

The five and only five registrations are
`POST /apis/execution.sovrunn.io/v1alpha1/execution-targets`,
`GET /apis/execution.sovrunn.io/v1alpha1/execution-targets`,
`GET /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}`,
`POST /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify`,
and
`POST /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire`.
Missing or invalid authentication on any one is `AUTH_REQUIRED`/401 with no
AuditEvent. An authenticated caller lacking the applicable create, qualify, or
retire grant receives audited `AUTHORIZATION_DENIED`/403. An inaccessible
direct item GET, qualify, or retire is safe `RESOURCE_NOT_FOUND`/404 with the
approved safe-denial audit behavior.

#### F16-AD-06.1 — exact request grammar and ETag meaning

Missing, empty, malformed, or overlong `Idempotency-Key` is
`MALFORMED_REQUEST`/400. Non-`application/json` POST is
`UNSUPPORTED_MEDIA_TYPE`/415. Missing, malformed, weak, wildcard, or multiple
action `If-Match` is a header-form failure and returns
`STALE_RESOURCE_VERSION`/412 before body classification. Any non-zero-byte
action body is `MALFORMED_REQUEST`/400. After replay/reservation, a
syntactically valid but non-current `If-Match` also returns
`STALE_RESOURCE_VERSION`/412, before lifecycle evaluation. ETag represents the
persisted target resource version; it is used only by qualify and retire.
F0016 defines no conditional GET. An effective availability change caused only
by fact expiry may therefore be projected with an unchanged ETag.

#### F16-AD-06.2 — closed field and response negotiation outcomes

After a caller has passed the applicable writer authorization, an attempted
client `status` write returns `AUTHORIZATION_DENIED`/403 with
`VS0_STATUS_FIELD_WRITE` and exactly one redacted AuditEvent; an attempted
server-owned identity or metadata write returns `AUTHORIZATION_DENIED`/403
with `VS0_SYSTEM_OWNED_FIELD_WRITE` and the same audit rule. An unknown field,
including `adapterAuthorityRef`, returns `UNKNOWN_FIELD`/400 with no audit.
Registry-declared invalid required or semantic spec values return
`VALIDATION_FAILED`/422. F0016 does not negotiate representations: any
`Accept` value is ignored and the fixed success/problem media type is emitted.

#### F16-AD-04/06.3 — exact successful representations

Create, item GET, qualify, and retire return the safe ExecutionTarget
projection as their JSON body; create returns 201 and the other three return
200. F0016 defines no `Location` header. LIST returns an ascending JSON array
of the same projection. That projection is exactly `metadata.uid`,
`metadata.name`, `metadata.generation`, `metadata.resourceVersion`, the three
immutable `spec` fields in F16-AD-02.1, `status.lifecycle`,
`status.qualification`, and response-only `effectiveAvailability`; F0016
defines no public `conditions` member. It omits facts, results, observer data,
handles, and internal record links.

#### F16-AD-08.1 — replay after waiting

Every request performs current authorization and safe target/scope access
before reservation lookup. A same-key waiter repeats both checks immediately
after waking and before receiving a completed response. A changed digest in
the same namespace is `CONFLICT`/409 with
`VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`. Replay returns the stored successful
status/body and `Content-Type: application/json`, but fresh request
correlation only.

#### F16-AD-08.2 — completed-record expiry and abort

Completed records expire at 24 hours. When capacity exceeds 10,000, eviction
is earliest expiry then lexical namespace; InFlight records are never evicted.
After expiry or eviction, the next valid request is processed as new. Panic,
server shutdown, cancellation of the owner, and an audit/stale/validation
abort remove the reservation and wake waiters without a completion. A
cancelled waiter detaches without aborting the still-running owner.

#### F16-AD-08.3 — namespace-separation proof

The same principal, key, and canonical digest on a different concrete action
target creates an independent reservation and cannot replay another target’s
result. The same principal, key, and canonical create digest in a different
derived CloudProvider scope is likewise independent. Both are mandatory local
proof cases, not merely namespace implementation detail.

#### F16-AD-08.4 — canonical digest inputs

The namespace is principal UID, route pattern, derived CloudProvider UID,
concrete target UID for an action, and `Idempotency-Key`. The digest is
SHA-256 of canonical JSON containing only the allowed create fields for
create; object members are sorted recursively and no unknown, duplicate, or
invalid body is digested or reserved. Qualify and retire digest zero bytes.
`If-Match`, `Accept`, authentication, request/correlation IDs, and every other
header are excluded. A valid same-key completed action may therefore replay
after the target's ETag has changed, subject to its current authorization and
safe-access checks.

#### F16-AD-08.5 — idempotency applicability

Idempotency applies only to collection create, qualify, and retire: the three
F0016 POST operations. Only those operations validate `Idempotency-Key`, may
reserve or wait on a namespace, create a completion, or replay. GET and LIST
neither require nor inspect `Idempotency-Key`; they never reserve, wait on,
read, or replay an idempotency record.

#### F16-AD-07.1 — item target resolution before target-bound grant

For item GET, qualify, and retire, phase-one safe resolution reads only the
path UID needed to derive the target’s CloudProvider scope. If that target is
inaccessible, the request is safe `RESOURCE_NOT_FOUND`/404. Only after
successful safe resolution is the target-bound `executiontarget.read`,
`executiontarget.qualify`, or `executiontarget.write` grant evaluated; an
accessible target without that grant returns audited
`AUTHORIZATION_DENIED`/403. This is not a general reference graph validation
and reveals no target details.

#### F16-AD-11/12/13.1 — exact internal record contract

The target holds optional current `status.factSetRef` and
`status.qualificationResultRef`. Create has neither; every completed
Qualified, Rejected, or Indeterminate conclusion atomically replaces both,
sets `status.observedGeneration` to its captured InfrastructureStack
generation, advances resourceVersion/ETag, and appends the required AuditEvent
before publication. Retire clears both. A FactSet contains target ref, referenced InfrastructureStack generation,
maintenance epoch, viability fingerprint, observer ID/revision, fact version,
observedAt, expiresAt, and exactly one named `compute.vm`, `storage.block`,
`storage.object`, and `network.private` truth value. A Result contains
target/fact-set refs, the same fences, profile version, outcome, reason codes,
and evaluation time. Duplicate/missing facts or malformed provenance are
observer faults: `INTERNAL_ERROR`/500 with no conclusion publication. The
fixed values are `observerID=sovrunn.synthetic-iaas-observer/v1`,
`profileVersion=synthetic-iaas/v1`, and `factSchemaVersion=v1`; fixture
revision is the immutable revision string carried by the selected target-bound
fixture. A timeout is a fixture-declared logical outcome, not a network wait.

#### F16-AD-15/16/19.1 — valid state combinations and visibility

Create persists Active/Unqualified with maintenance epoch 0 and no record
refs, and its response projects effectiveAvailability Unavailable. Qualify may
start from Active Unqualified, Qualified, Rejected, or Indeterminate only. It
creates a lifecycle-service-only Qualifying reservation that preserves the
last committed projection, current FactSet/Result links, and ETag; it neither
changes persisted target status nor emits an AuditEvent or idempotency
completion. Supported
facts conclude persisted Active/Qualified and project Available; an Unsupported
fact concludes Active/Rejected and projects Unavailable; Unknown concludes
Active/Indeterminate and projects Unavailable. Maintenance entered during
Qualifying aborts that qualification reservation, clears current record links,
and persists Unqualified before projecting Maintenance; otherwise it preserves
the last completed qualification while projecting Maintenance. Clear persists
Active/Unqualified and projects Unavailable. Retire is allowed from every
Active combination, including Maintenance and Qualifying, and persists
Retired/Unqualified while projecting Unavailable. Retired targets remain
visible to authorized GET/LIST callers.

#### F16-AD-16/18/19.2 — winning-transition response

If Retire wins while a qualification is in flight, the qualifying caller gets
`CONFLICT`/409 with `VS0_TARGET_RETIRED`. If Maintenance entry wins, that
caller gets `CONFLICT`/409 with `VS0_TARGET_MAINTENANCE`. In either case the
qualification reservation aborts without a qualification result, completion,
or qualification AuditEvent; the winning Retirement or Maintenance transition
is audited under F16-AD-21.

At qualification commit these outcomes are mutually exclusive and ordered
(ADH-2026-060): an active current-Maintenance marker always selects
`CONFLICT`/409 + `VS0_TARGET_MAINTENANCE` first; only when no active marker
exists does a changed captured maintenance epoch select
`STALE_RESOURCE_VERSION`/412 + `VS0_TARGET_EPOCH_STALE`. A changed referenced
InfrastructureStack generation remains the independent epoch-stale case.

#### F16-AD-16.1 — qualification abort and stale-fence disposition

If malformed observer output, owner cancellation, required qualification
AuditEvent append failure, or a stale captured fence aborts a Qualifying
reservation, the lifecycle service leaves the last committed qualification
summary, FactSet/Result links, resourceVersion/ETag, and projection unchanged;
it publishes no AuditEvent and aborts the idempotency reservation. A stale
InfrastructureStack generation or viability fingerprint is reported by its
exact 412 violation without fabricating a current observation. If Retire or
Maintenance has already won, its existing audited transition remains the sole
committed result and no restoration occurs.

#### F16-AD-08/20.1 — stale waiter disposition

A same-key waiter awakened after an owner qualification abort has no cached
completion. It repeats the normal request pipeline from current authorization
and safe target/scope access; it does not inherit the owner's stale, cancelled,
or failed result. If its current request remains valid, it may establish a new
reservation. A cancelled waiter still only detaches and cannot abort a
still-running owner.

#### F16-AD-15.2 — response-only effective availability

The persisted target contains lifecycle, qualification, maintenance epoch,
observed generation, and the optional record refs; it contains **no**
`status.availability` field. The server computes and returns
`effectiveAvailability` as a response-only value of Available, Unavailable,
or Maintenance. It has no independently persisted writer and does not change
the resource ETag. Registry/schema entries must use this projection name and
must not claim a persisted three-value availability status field.

#### F16-AD-15.3 — effective availability truth table

The response projection is deterministic in this order: Retired projects
Unavailable; otherwise Active with current Maintenance projects Maintenance;
otherwise Active with Qualified qualification, fresh current facts, and a
viable backing projects Available; every other Active combination projects
Unavailable. A viability-only change does not mutate the target or change its
ETag. A fresh Qualified target therefore becomes Unavailable when its
participation loses effective Active state and becomes Available again if that
state returns before the facts expire.

#### F16-AD-15.4 — observed-generation meaning

`status.observedGeneration` is the last successfully committed FactSet's
referenced InfrastructureStack generation. Create initializes it from the
current referenced generation and a successful qualification replaces it with
the captured generation. Expiry, reservation abort, stale completion,
Maintenance entry/clear, and Retirement do not alter it. It is diagnostic
history, never evidence that the currently referenced generation has been
observed.

#### F16-AD-15/18.1 — persisted fence transition values

Create assigns `status.observedGeneration` to the referenced
InfrastructureStack `metadata.generation` and `status.maintenanceEpoch=0`.
The captured value is a backing-resource fence, never the ExecutionTarget's
own immutable post-create generation. Each successful maintenance enter and clear
increments maintenanceEpoch, increments the target resourceVersion, and
therefore changes ETag. Both clear current FactSet/Result links; the
qualification value on enter is retained only as a safe summary, while clear
sets Unqualified. No pre-fence record is current after either transition.

#### F16-AD-15.5 — current Maintenance marker

`status.maintenanceEpoch` is a persisted fence, not proof that Maintenance is
currently active. The lifecycle service owns one target-bound, in-process
current-Maintenance marker containing target UID, maintenance epoch, and an
active flag. It is never client-writable or projected, and restart clears it
with all other F0016-owned in-process state. A successful Maintenance entry
advances the epoch, clears current FactSet/Result links, sets qualification to
Unqualified, establishes the active marker, and appends its AuditEvent before
publication. A clear requires the current marker and matching epoch, advances
the epoch, removes the marker, persists Active/Unqualified, and appends its
AuditEvent before publication. `effectiveAvailability` is Maintenance exactly
while the current marker is active; otherwise the F16-AD-15.3 truth table
applies.

#### F16-AD-06.4 — remaining closed body/schema outcomes

Malformed JSON, duplicate top-level members, and an oversized request body use
their exact inherited 400 Problem outcomes. A missing required create field,
an invalid typed reference shape, or a `targetClass` value other than
`synthetic-iaas` returns `VALIDATION_FAILED`/422. None publishes a mutation,
AuditEvent, or idempotency completion.

#### F16-AD-06.5 — closed create headers, query strings, and path spelling

`If-Match` on collection create is `MALFORMED_REQUEST`/400 with no mutation,
AuditEvent, or idempotency completion. Every F0016 route accepts no query
parameter: `page`, `limit`, `filter`, and unknown query names are each
`MALFORMED_REQUEST`/400 with the same no-publication effect. Only the five
exact non-trailing-slash paths are registered; a trailing-slash variant is
unregistered and returns transport-only 404 with no redirect. The fixed
pre-ServeMux guard runs before the five registrations: it rejects `HEAD` on
the collection path with transport-only 405 and `Allow: GET, POST`, on the
item path with 405 and `Allow: GET`, and on either action path with 405 and
`Allow: POST`. It rejects every trailing-slash variant with transport-only
404 and no redirect, and any other unmatched method on an exact F0016 path
with transport-only 405. These guard outcomes have no Problem body,
authentication, authorization, audit, idempotency, or lifecycle effect.

#### F16-AD-06.6 — exact response media and anti-sniffing

Every successful F0016 representation uses `Content-Type: application/json`;
every F0016 Problem Details response uses
`Content-Type: application/problem+json`; and every F0016 response, including
transport-only 404/405 and a replay, sends
`X-Content-Type-Options: nosniff`. `Accept` is ignored. A transport-only
404/405 still has no Problem body and introduces no Problem code.

#### F16-AD-06/07.2 — authenticated URI grammar

After authentication and before phase-one safe resolution, a valid registered
method/path pair must have an empty query string and, for an item/action path,
a syntactically valid canonical UID segment. A malformed query or UID returns
the inherited 400 malformed-request outcome with no reference lookup,
mutation, AuditEvent, idempotency reservation, or completion. Thus malformed
URI grammar never becomes a safe-404 existence probe, while unauthenticated
requests still stop at 401 first.

#### F16-AD-17/21.1 — expiry and audit retry

Expiry leaves the last completed qualification and ETag unchanged but projects
Unavailable. The lifecycle service appends one expiry AuditEvent. If append
fails, the projection remains Unavailable and the service retries once per
injected-clock one-second tick until one append succeeds; it emits no duplicate
successful expiry event. Expiry has no request caller: it returns no 500,
publishes no expiry AuditEvent until the append succeeds, and is the sole
exception to F16-AD-21’s request-triggered append-failure rule.

#### F16-AD-14/17/20.2 — graceful lifecycle shutdown

The service shutdown boundary first stops acceptance of expiry and Maintenance
triggers, then aborts in-flight qualification and idempotency reservations and
wakes their waiters, and only then stops HTTP serving. A trigger accepted after
that stop boundary must neither mutate target state nor append an AuditEvent.
Shutdown preserves the last committed target projection, ETag, FactSet/Result
links, and inherited AuditEvents.

#### F16-AD-21.2 — exact AuditEvent append-failure branches

Create, qualification, retirement, maintenance-entry, and maintenance-clear
append failures are separate mutation branches: each returns
`INTERNAL_ERROR`/500 and publishes neither its mutation nor an idempotency
completion. Authorization-denial, safe-denial, replay authorization-denial,
and replay safe-denial append failures are separate denial branches: each
returns `INTERNAL_ERROR`/500, publishes no AuditEvent, and hides the original
403/404 outcome and violation.

#### F16-AD-21/25/26.2 — finite denial-proof families

A named conformance family is permitted only for a finite, explicitly listed
set of requests whose HTTP status, Problem code/type, violation, mutation,
AuditEvent effect, idempotency effect, and replay class are identical. The
annex must name every member and the readiness checker must execute every
member; a family may not hide a route-, grant-, safe-access-, or replay-specific
difference. The denial-audit-append family covers LIST without read, direct
GET without read, qualify without qualify grant, retire without write grant,
and inaccessible direct qualify or retire. For each member, failed required
AuditEvent append returns `INTERNAL_ERROR`/500, publishes no AuditEvent,
mutation, lifecycle/record change, or idempotency completion, and hides its
original 403/404 outcome and violation.

#### F16-AD-20/22/26.1 — exact proof allocation

Each F0016-local conformance row has one concrete request or trigger and one
exact outcome: state, Problem code/HTTP/type, violation, audit effect, replay
effect, and gate. The only exception is the finite, explicitly enumerated
equivalence family permitted by F16-AD-21/25/26.2: every named member must
have the identical declared observable outcome and effects, and the
readiness checker must execute every member. The ADH-058 annex is
authoritative for that inventory. Terms such as “where applicable”, “exact
inherited”, or “conflict” are not an outcome contract.

Within that annex, `no publication` is exact shorthand only for no target or
internal-record mutation, ETag change, or idempotency completion. Each row
must separately state its AuditEvent effect; an audited denial is therefore
not made unaudited by this shorthand.

#### F16-AD-25.1 — required formal invariants

The lifecycle formal model proves Retired terminality, the effectiveAvailability
truth table, Maintenance and Retirement priority over expiry, and that
Qualifying is never publicly persisted. The fence/publication model proves
that a qualification conclusion and replay completion require a successful
AuditEvent append, while stale, cancelled, and failed proposals cannot mutate
the target or become replayable completions.

#### F16-AD-26.3 — atomic placeholder replacement

After founder approval, ADH-058 must atomically replace the active placeholder
definitions of `VS0-SCHEMA-015..017` and `VS0-STATE-004`; it must not append
an alternative F0016 contract. The applied authorities contain no active
`status.availability`, `Draining`, or publicly projected `Qualifying` value,
and accept no alias, migration path, or legacy compatibility behavior for
them. Until that approved replacement is applied, the current baseline remains
unchanged and requirements generation is not authorized.

### Architecture-admission rule

No F0016 requirements prompt may be generated until every F16-AD row is
approved in authority or explicitly deferred to a named later feature with no
active F0016 route, field, state, writer, record, or conformance claim. A
phrase such as “implementation-defined”, “handled by the adapter”, or “reuse
the established behavior” is not closure unless it names the exact inherited
authority and proves semantic compatibility.

## 6. ExecutionTarget realization taxonomy

This is a planning taxonomy, not an approved enum or public API. Every entry
remains the same canonical `ExecutionTarget` resource. A profile may be
activated only after an ADH defines its backing rule, normalized facts,
protected authority boundary, lifecycle effects, and local proof.

| Profile | Meaning | Example | F0016 disposition |
|---|---|---|---|
| `synthetic-iaas` | Deterministic fake shaped like a minimum IaaS environment. | Reports compute, block, object, and private-network readiness. | Active and required. |
| `infrastructure-iaas` | A target backed by a CloudProvider-approved InfrastructureStack. | A local CloudProvider's OpenStack-like substrate. | Future real adapter. |
| `private-cloud` | A privately operated virtualized/IaaS substrate. | VMware-like or OpenStack-like control plane. | Deferred. |
| `container-platform` | A qualified container execution environment used by future service plugins. | Kubernetes/OpenShift-like environment. | Deferred; native CRDs/objects remain internal. |
| `external-managed-control-plane` | An external API that realizes an approved service without direct stack backing. | Managed database control plane. | Deferred pending a backing-reference decision. |
| `edge-site` | A constrained or intermittently connected local environment. | An edge appliance control plane. | Deferred pending offline/freshness/recovery rules. |
| `federated` | An independently operated remote control plane. | A sovereign partner environment. | Future-only; no federation semantics in F0016. |

## 7. Adapter boundary

This section explains the boundary already decided in §5; it adds no route,
state, writer, or observable-contract semantics.

### 7.1 Direction of dependency

```text
Sovrunn canonical contracts
        -> ExecutionTarget adapter contract
        -> adapter implementation
        -> external realization system
```

Core code depends only on Sovrunn-owned contracts. An adapter translates
native observations into normalized facts and classified failures. Native IDs,
CloudProvider-native API types, credentials, CloudProvider-specific capability objects,
and implementation errors remain behind this boundary.

### 7.2 Future Crossplane position

Crossplane is a future implementation candidate for an infrastructure-backed
target. It would be reached only through a later `CrossplaneExecutionAdapter`:

```text
Sovrunn control plane
  -> ExecutionTarget / adapter contract
  -> CrossplaneExecutionAdapter
  -> Crossplane control plane and packages
  -> CloudProvider IaaS management plane
```

Crossplane must not replace the canonical model, define customer-facing APIs,
own CloudProvider participation, or become a required F0016 dependency.

## 8. FEATURE-0016 workflow within Sovrunn

This workflow is explanatory only. The register in §5 remains authoritative
for all observable behavior and ownership.

### 8.1 F0016 qualification workflow

```text
CloudProvider administrator
  -> creates, qualifies, or retires an internal ExecutionTarget
  -> current authorization and safe participation/InfrastructureStack access
  -> TargetObservationAdapter receives only protected target authority handle
  -> emits normalized, provenance-bearing, expiring target facts
  -> TargetQualificationEngine evaluates the closed synthetic-IaaS profile
  -> proposed immutable TargetQualificationResult
  -> ExecutionTargetLifecycleService atomically appends AuditEvent plus publishes current status and records
  -> qualified, fresh, epoch-current target becomes an internal candidate
```

The observation adapter never decides policy, placement, customer visibility,
or realization. The qualification engine never calls an external environment.
The qualification result does not create a VM, volume, bucket, network, or
customer-visible resource.

### 8.2 Relationship to the full Sovrunn flow

```text
FEATURE-0015
CloudPlatform + CloudProvider + participation + topology + InfrastructureStack
  -> FEATURE-0016
     ExecutionTarget observation -> normalized facts -> core qualification
  -> FEATURE-0019 / FEATURE-0017
     sovereignty facts and policy evaluation
  -> FEATURE-0021 / FEATURE-0022
     entitlement, quota, product/runtime requirements, ServiceRegion
  -> FEATURE-0023
     approved candidate selection and safe ServicePlacement
  -> FEATURE-0024
     ServiceDeploymentPlan -> PluginExecution -> synthetic realization adapter
  -> Phase 3
     real realization adapter -> external environment
```

An ExecutionTarget is therefore a necessary qualified candidate boundary, not a
placement decision and not execution authorization. Qualification answers
whether a target meets its technical profile now; later policy and placement
features decide whether it may be selected for a particular service request;
only later PluginExecution and real realization may act on an external system.

### 8.3 Observation and realization separation

| Concern | F0016 observation/qualification | F0024 and Phase 3 realization |
|---|---|---|
| Input | ExecutionTarget, protected authority handle, expected generation/epoch | Authorized ServiceDeploymentPlan/PluginExecution, selected target, protected handle |
| Output | NormalizedTargetFactSet and TargetQualificationResult | Verified execution outcome and protected execution records |
| External effect | None in Phase 2R | Synthetic only in F0024; real effect deferred to Phase 3 |
| Owner of decision | Sovrunn core qualification engine | Approved plan, plugin, operation, and future realization adapter |
| Customer visibility | None | Only later safe ServicePlacement and ServiceBinding projections |

## 9. Required admission evidence before requirements generation

1. Every F16-AD outcome is approved or explicitly deferred
   to a named later owner.
2. Registry, specification, feature authority, traceability, and closure
   matrix express one coherent schema/writer/state/route/conformance contract.
3. Every observable route has closed field classification, authorization,
   reference-resolution order, success/error/precedence, status, audit, and
   replay behavior.
4. Every status and record field has one writer, initial value, transition
   source, freshness/epoch rule, and local proof case.
5. A deterministic fake demonstrates both qualification and denial with zero
   external effect and redacted evidence.
6. F0016-local conformance, architecture-readiness checks, and the lifecycle
   and stale-fence formal models pass.
7. Founder architecture approval records the approved authority digest.
8. A requirements-stage context manifest has measured headroom against its
   configured budget; authority classification, rather than omission of a
   required contract, is the only permitted way to reduce context.
9. Before design approval, a deterministic executability audit maps every
   F0016 route to authentication, authorization, safe resolution, field
   classification, state precondition, idempotency/concurrency, publication,
   response, error, audit, and local conformance evidence.
10. Before tasks approval, a deterministic task-scope inventory derives the
    implementation, fixture, schema, server-lifecycle, and test-import paths
    needed by each task. Every derived path must be in the task's writable
    scope, and the dependency graph must be acyclic before Cursor is invoked.

## 10. Explicit non-goals

- Selecting a real CloudProvider or real adapter implementation.
- Installing or depending on Crossplane, Kubernetes, OpenShift, PostgreSQL,
  or an external secret system in Phase 2R.
- Introducing target subtypes, ResourcePool, a CloudProvider-wide capability
  inventory, uncontrolled polling, placement, or a customer-visible target API.
- Creating a new generic supply term in place of the canonical CloudProvider
  and ExecutionTarget boundaries.
- Generating requirements, design, tasks, or code from this proposed draft.

## 11. Approval path

Use this draft and the architecture starting dossier to prepare the bounded
ADH-058 decision package. Founder approval is required before Kiro updates
approved architecture authorities. Only then may FEATURE-0016 enter
requirements generation.
