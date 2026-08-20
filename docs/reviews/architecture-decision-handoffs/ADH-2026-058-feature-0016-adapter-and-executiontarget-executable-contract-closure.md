# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-058
- Date: 2026-08-19
- Source discussion: Codex architecture session
- Related feature: FEATURE-0016
- Related phase: Phase 2R
- Author: Codex (proposed for founder review)
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Close FEATURE-0016 adapter and ExecutionTarget executable contract

## Summary

FEATURE-0016 activates one deterministic, infrastructure-backed
`synthetic-iaas` ExecutionTarget qualification profile. It establishes a
Sovrunn-owned observation-and-qualification boundary, without real external
calls, credentials, placement, provisioning, adapter selection, or a customer
projection. This handoff corrects the placeholder F0016 registry schemas and
state machine so requirements, design, and tasks can translate one complete
observable contract without creating architecture decisions downstream.

## Classification

New decision.

## Existing approved baseline

- FEATURE-0015 owns the canonical CloudPlatform, CloudProvider,
  CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, and
  InfrastructureStack resources, and explicitly delegates the complete
  ExecutionTarget contract to FEATURE-0016.
- DEC-0036 / ADH-2026-042 require implementation-neutral adapter boundaries
  that do not expose native types, raw credentials, or vendor errors in the
  Sovrunn core/customer contract.
- VS0-SCHEMA-015..017 and VS0-STATE-004 reserve FEATURE-0016 ownership, but
  they are placeholder contracts: they cannot express the required
  three-valued facts, immutable internal result records, lifecycle retirement,
  expiry semantics, or fully assigned maintenance fencing.
- FEATURE-0012 owns common API grammar, typed references, status rules,
  Problem Details, media semantics, and optimistic concurrency. FEATURE-0013
  owns AuditEvent and its atomic append-before-publication boundary.
- Phase 2R permits deterministic in-process simulation only. It forbids real
  external infrastructure effects and durable or external persistence
  dependencies.

## Decision or proposed decision

This proposed handoff replaces the accumulated F0016 draft decisions with the
single canonical 28-row register in FEATURE-0016 architecture.md §5. These
are the only proposed F0016 architecture decisions; subclauses carry detail
instead of creating later decision numbers.

`FEATURE-0016 architecture.md` §5 and its subclauses are the sole proposed
semantic authority. This handoff is an application package: its contract
extract and annex operationalize that register, but must not redefine or
extend it.

When founder-approved, this handoff atomically replaces—not extends—the active
placeholder definitions of `VS0-SCHEMA-015..017` and `VS0-STATE-004`. The
applied authorities must contain no active persisted `status.availability`,
`Draining`, or publicly projected `Qualifying` value, and must provide no
alias, migration path, or legacy compatibility behavior for them. Before that
approval and application, the current baseline remains unchanged.

| Domain | Canonical rows | Proposed contract |
|---|---|---|
| Boundary/resource | 01 | One CloudProvider-confidential, CloudProvider-scoped ExecutionTarget; no placement, provisioning, realization, or IaaS exposure. |
| Create/backing | 02–03 | Closed create fields; participation-derived scope; scope-name/tuple uniqueness; safe active same-scope backing references. |
| Route/security | 04–07 | Five routes; exact grants/media/ETag/zero-byte actions; fail-closed precedence. |
| Idempotency | 08 | Principal, route, scope, digest, and action-target binding; bounded process-local replay. |
| Observer/facts | 09–13 | No adapter resource/credential/external call; one synthetic observer; four normalized facts; redacted internal records. |
| State/freshness | 14–20 | Sole lifecycle committer; effective availability; qualification, expiry, maintenance, retirement, stale-work rules. |
| Audit/projection | 21–23 | Append-before-publication; closed local violations; safe redacted projection. |
| Future/proof/admission | 24–28 | No F0016 Crossplane/realization; local proof; fail-closed downstream gates. |

### 1. Canonical contract subclauses

1. ExecutionTarget is the sole F0016 target. Create accepts exactly
   metadata.name, spec.cloudProviderParticipationRef.uid,
   spec.infrastructureStackRef.uid, and immutable
   spec.targetClass="synthetic-iaas". Scope derives from participation.
   The server alone assigns metadata.uid, metadata.resourceVersion,
   metadata.generation, status.lifecycle, status.qualification,
   status.maintenanceEpoch, status.observedGeneration, status.factSetRef, and
   status.qualificationResultRef. Unknown fields and adapterAuthorityRef are
   rejected. No separate scopeRef is persisted or projected: scope is derived
   only from spec.cloudProviderParticipationRef.uid. No SecretRef extension is
   active.

2. Safe reference access precedes graph validation. Create/qualify require a
   FEATURE-0015 effective-Active participation (accepted with neither
   suspension hold) and Active InfrastructureStack in the same CloudProvider
   scope. Name remains scope-unique for the process lifetime, including after
   retirement; retirement releases only the tuple. At most one non-Retired
   (participation UID, stack UID, targetClass) tuple is live. The viability
   fingerprint contains only participation UID/effective-active state and
   stack UID/phase. Inaccessible backing is safe 404; authorized ineffective
   participation, non-Active stack, and cross-scope graph use the exact annex
   outcomes.

3. The five exact Go 1.22 http.ServeMux registrations are POST/GET
   /apis/execution.sovrunn.io/v1alpha1/execution-targets, GET
   /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}, POST
   /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify,
   and POST
   /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire.
   No PATCH, PUT, DELETE, HEAD, watch, filter, pagination, fact, result, or
   adapter-selection route exists. A fixed pre-ServeMux F0016 method/path
   guard runs before those registrations. It rejects HEAD on the collection
   path with transport-only 405 and Allow: GET, POST; HEAD on the item path
   with 405 and Allow: GET; and HEAD on either action path with 405 and
   Allow: POST. It rejects every trailing-slash variant with transport-only
   404 without redirect. The guard adds no public route or handler and causes
   no Problem body, authentication, authorization, audit, idempotency, or
   lifecycle effect. Other unmatched methods on an exact F0016 path are also
   transport-only 405. Valid registered method/path pairs bypass the guard
   and begin at authentication.
   Create/retire use executiontarget.write; reads use .read; qualify uses
   .qualify. Item operations bind target UID and derived CloudProvider scope.
   LIST contains all accessible targets, including Retired, in ascending
   metadata.uid order; no read grant returns audited 403 rather than empty 200.
   Missing/invalid authentication is 401 with no audit; direct inaccessible
   item operations are safe 404 with the approved denial audit. For item GET
   and either action, phase-one safe resolution reads only the path UID needed to derive its
   CloudProvider scope: an inaccessible target is safe 404; only then is the
   target-bound qualify/write grant evaluated, so an accessible target without
   that grant is audited 403. Neither outcome discloses target details.

4. Every POST requires application/json and exactly one Idempotency-Key.
   Qualify/retire also require exactly one strong opaque If-Match and a
   zero-byte body. Header-form validation and current-version comparison are
   distinct. Create/item GET/qualify/retire return ETag; LIST does not.
   Success uses application/json, errors application/problem+json, and every
   response, including a transport-only 404/405 and replay, sends
   X-Content-Type-Options: nosniff. Missing/bad key is 400,
   non-JSON media is 415, malformed action If-Match form is 412 before body
   classification, a valid but non-current If-Match is 412 after replay, and
   a non-zero action body is 400. After writer authorization, status write is audited 403
   VS0_STATUS_FIELD_WRITE; server-owned identity/metadata write is audited 403
   VS0_SYSTEM_OWNED_FIELD_WRITE; unknown field (including adapterAuthorityRef)
   is 400 UNKNOWN_FIELD; invalid declared spec is 422 VALIDATION_FAILED.
   F0016 ignores Accept and emits its fixed response media. ETag represents
   persisted target resource version only; F0016 has no conditional GET, so
   expiry projection may change with unchanged ETag. Create, item GET,
   qualify, and retire return a safe ExecutionTarget JSON projection; create
   is 201 and the others are 200, with no Location header. LIST returns an
   ascending JSON array of that projection. It contains exactly metadata.uid,
   metadata.name, metadata.generation, metadata.resourceVersion, the three
   immutable spec fields, status.lifecycle, status.qualification, and
   response-only effectiveAvailability; F0016 defines no public conditions
   member. Facts, results, observer data, handles, and internal record links
   are omitted. Malformed JSON,
   duplicate top-level members, and oversized bodies use their exact inherited
   400 Problem outcomes; missing required create fields, invalid typed
   reference shapes, and a targetClass other than synthetic-iaas are
   VALIDATION_FAILED/422 with no publication. If-Match on collection create is
   MALFORMED_REQUEST/400. No F0016 route accepts a query parameter: page,
   limit, filter, and unknown names are MALFORMED_REQUEST/400. Only the exact
   non-trailing-slash paths are registered; the pre-ServeMux guard supplies
   the transport-only trailing-slash and HEAD outcomes in clause 3.

5. Request precedence is closed transport method/path guard; authentication;
   method/media/header-form validation; phase-one safe references; authorization; safe access;
   strict classification;
   graph/reference validation; replay/reservation; current If-Match version comparison; lifecycle state;
   atomic commit. Create alone has client graph/reference fields and no
   If-Match; zero-byte actions alone have If-Match and no client graph. An
   action header-form failure therefore precedes its non-zero-body failure;
   current-version comparison occurs after replay/reservation and before
   Maintenance state evaluation. The namespace is principal UID, route pattern,
   derived CloudProvider UID, concrete action-target UID when applicable, and
   Idempotency-Key. The digest is SHA-256 of recursively sorted allowed create
   JSON or SHA-256 of zero bytes for actions; If-Match, Accept,
   authentication, request/correlation IDs, and all other headers are excluded.
   Invalid/unknown/duplicate bodies are never digested or reserved. Idempotency
   applies only to the three POST operations: create, qualify, and retire.
   GET and LIST neither require nor inspect Idempotency-Key and never reserve,
   wait on, read, or replay a record. Replay binds principal, route, scope UID, canonical digest,
   and action target UID. Successful create/qualification/retirement retain
   for 24 hours/10,000 entries only in this process; aborted outcomes do not
   replay; restart clears F0016-owned in-memory targets, records, idempotency
   reservations, and pending work but never inherited FEATURE-0013 AuditEvents.
   Completed entries evict at expiry, then
   earliest-expiry/lexical namespace when capacity is exceeded; InFlight never
   evicts. Expired/evicted requests process as new. Panic, shutdown, owner
   cancellation, and every abort wake waiters with no completion; a cancelled
   waiter only detaches. A waiter rechecks current authorization
   and safe target/scope access immediately after waking; changed digest is
   409 with VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH. The same principal, key, and
   digest on a different action target, or on a create in another derived
   CloudProvider scope, forms an independent reservation and never replays the
   other outcome.

6. The only observer is sovrunn.synthetic-iaas-observer/v1. It is
   target-bound, clock-driven, external-effect-free, and yields facts or
   classified outcomes. Exactly once each of compute.vm, storage.block,
   storage.object, and network.private is Supported, Unsupported, or Unknown.
   Any Unsupported rejects; otherwise any Unknown is Indeterminate; otherwise
   Qualified. Missing fixture and fixture-declared logical timeout are separate
   deterministic outcomes, each yielding four Unknown facts at now through
   now+60 seconds; neither waits on a network call. A FactSet carries target ref, referenced InfrastructureStack generation,
   maintenance epoch, viability fingerprint, observer ID/revision, fact
   version, observedAt/expiresAt, and the four named truth values. A Result
   carries target/fact-set refs, the same fences, profile version, outcome,
   reason codes, and evaluation time. Duplicate/missing facts or malformed
   provenance are INTERNAL_ERROR/500 observer faults and publish no result.
   Fixed provenance is observerID sovrunn.synthetic-iaas-observer/v1,
   profileVersion synthetic-iaas/v1, factSchemaVersion v1, and the immutable
   selected fixture revision.

7. ExecutionTargetLifecycleService alone commits target state and internal
   fact/result records. Lifecycle is Active/Retired; persisted public
   qualification is Unqualified/Qualified/Rejected/Indeterminate; Qualifying
   is a lifecycle-service-only in-flight reservation, never a projected
   resource status. The persisted
   target has no status.availability: the server returns response-only
   effectiveAvailability (Available, Unavailable, Maintenance), computed from
   committed state, freshness, maintenance, and viability. It has no
   independently persisted writer and does not alter ETag. Qualify is synchronous; stale generation,
   maintenance epoch, or viability completion returns 412 and publishes no
   conclusion. Create persists Active/Unqualified at epoch 0, sets
   observedGeneration to the referenced InfrastructureStack metadata.generation,
   and has no record refs; its response projects effectiveAvailability
   Unavailable. The target retains optional current factSetRef and
   qualificationResultRef; every completed conclusion atomically replaces
   both, sets observedGeneration to its captured InfrastructureStack
   generation, advances resourceVersion/ETag, and appends AuditEvent before
   publication; Retire clears both. Qualify starts only from Active
   Unqualified/Qualified/Rejected/Indeterminate: it creates a Qualifying
   reservation but leaves the committed target summary, current FactSet/Result
   links, and resourceVersion/ETag unchanged until a completed conclusion can
   be audited and published. Malformed observer output, owner cancellation,
   required qualification AuditEvent append failure, and every stale captured
   fence abort that reservation without target mutation, AuditEvent, or
   completion. A stale qualification proposal never suppresses an
   independently committed winning Maintenance or Retirement transition:
   that transition remains valid and audited. Maintenance entered during an
   in-flight Qualifying reservation aborts that reservation,
   clears current links, and persists
   Unqualified before projecting Maintenance; otherwise it preserves the last
   completed conclusion while projecting Maintenance. Clear persists
   Active/Unqualified and projects Unavailable; Retire is allowed from every
   Active combination, including while a Qualifying reservation is in flight,
   and persists Retired/Unqualified while projecting
   Unavailable. Every successful
   maintenance enter or clear increments maintenance epoch and resource
   version/ETag, clears both record links, and leaves no pre-fence record
   current; clear also sets qualification Unqualified. maintenanceEpoch is a
   fence, not proof that Maintenance is current. The lifecycle service owns a
   target-bound, in-process current-Maintenance marker containing target UID,
   maintenance epoch, and an active flag. It is never client-writable or
   projected and is cleared on restart. Entry establishes the marker; clear
   requires its matching current epoch and removes it. effectiveAvailability
   is Maintenance exactly while that marker is active.

   EffectiveAvailability is ordered: Retired projects Unavailable; otherwise
   Active with current Maintenance projects Maintenance; otherwise only Active
   Qualified with fresh current facts and viable backing projects Available;
   every other Active combination projects Unavailable. A viability-only
   change does not mutate the target or change its ETag.

   If Retire wins while a qualification reservation is in flight, the caller
   receives CONFLICT/409 plus VS0_TARGET_RETIRED. If Maintenance entry wins,
   the qualifying caller receives CONFLICT/409 plus VS0_TARGET_MAINTENANCE. Both abort qualification
   without result, completion, or qualification AuditEvent; the winning
   transition itself is audited.

8. Expiry is injected-clock driven; Retired and Maintenance outrank it.
   Maintenance enters/clears only through a deterministic, target-bound
   synthetic-observer fixture trigger carrying target, expected
   InfrastructureStack generation/maintenance epoch, operation, system
   identity, and correlation. It is not an HTTP route or generic event bus.
   Retire is the sole lifecycle exit and releases the tuple. Audit successful
   create, completed qualification, retirement, maintenance enter/clear,
   authenticated authorization denial, and safe denial. Request-triggered
   required append failure is INTERNAL_ERROR/500 with no AuditEvent, mutation,
   or completion publication. Expiry is the sole background exception: it has
   no request caller, returns no 500, and publishes no expiry AuditEvent until
   its append succeeds. Expiry leaves the last qualification and ETag unchanged, projects
   Unavailable, and writes one expiry AuditEvent; a failed append retries once
   per injected-clock second until one append succeeds.

   Service shutdown first stops acceptance of expiry and Maintenance triggers,
   then aborts in-flight qualification and idempotency reservations and wakes
   waiters, and only then stops HTTP serving. A trigger accepted after that
   stop boundary must neither mutate target state nor append an AuditEvent.
   Shutdown preserves the last committed target projection, ETag,
   FactSet/Result links, and inherited AuditEvents.

   The lifecycle TLC model proves Retired terminality, the
   effectiveAvailability truth table, Maintenance and Retirement priority over
   expiry, and that Qualifying is never publicly persisted. The
   fence/publication TLC model proves that a qualification conclusion and
   replay completion require a successful AuditEvent append, while stale,
   cancelled, and failed proposals cannot mutate the target or become
   replayable completions.

9. The only new local violations are VS0_TARGET_RETIRED,
   VS0_TARGET_MAINTENANCE, VS0_TARGET_QUALIFICATION_IN_PROGRESS,
   VS0_TARGET_EPOCH_STALE, VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE,
   VS0_EXECUTION_TARGET_STACK_UNAVAILABLE,
   VS0_EXECUTION_TARGET_SCOPE_MISMATCH, and
   VS0_EXECUTION_TARGET_VIABILITY_STALE. FEATURE-0012 owns top-level Problems.

10. F0024/Phase 3 owns realization/plugin taxonomy. Crossplane is future-only.
    Requirements, design, tasks, and Cursor may not choose an omitted
    observable behavior, authority, proof meaning, or implementation scope.

### 2. Actual F0016 proof/outcome annex (proposed registry assignments)

Every row is one proposed F0016-local registry case. Each must record its exact
state, code/HTTP/type, violation, audit effect, replay effect, and readiness
assertion. A row may instead be a finite equivalence family only when it names
every member and every member has the same declared outcome and side effects;
the readiness checker must execute every member. It may not conceal a
route-, grant-, safe-access-, or replay-specific difference. Unless a row
itself tests authentication, request grammar, or safe access, its earlier
pipeline preconditions are valid and satisfied.

In this annex, `no publication` means no target or internal-record mutation,
ETag change, or idempotency completion. Every row separately declares whether
an AuditEvent is written; an audited denial is not made unaudited by the
shorthand.

| ID | Request or trigger | Exact outcome and side effects |
|---|---|---|
| VS0-CF-F16-01 | Valid create | 201 safe projection; persists Active/Unqualified, observedGeneration equals referenced InfrastructureStack generation, maintenance epoch 0, and projects effectiveAvailability Unavailable; ETag, AuditEvent, completed replay. |
| VS0-CF-F16-02 | Authorized create with one unknown field | 400 UNKNOWN_FIELD; no mutation, audit, or completion. |
| VS0-CF-F16-03 | Authorized create with client status | 403 AUTHORIZATION_DENIED plus VS0_STATUS_FIELD_WRITE; audited denial. |
| VS0-CF-F16-04 | Authorized create with client metadata.uid | 403 AUTHORIZATION_DENIED plus VS0_SYSTEM_OWNED_FIELD_WRITE; audited denial. |
| VS0-CF-F16-05 | Create/action with a missing Idempotency-Key | 400 MALFORMED_REQUEST; no publication. |
| VS0-CF-F16-06 | POST with non-application/json media | 415 UNSUPPORTED_MEDIA_TYPE; no publication. |
| VS0-CF-F16-07 | Create with missing authentication | 401 AUTH_REQUIRED; no mutation, audit, or replay record. |
| VS0-CF-F16-08 | Authenticated create without executiontarget.write | 403 AUTHORIZATION_DENIED; audited denial; no publication. |
| VS0-CF-F16-09 | Inaccessible create backing reference | Safe 404 RESOURCE_NOT_FOUND plus VS0_AUTHORIZATION_SAFE_DENIAL; audited denial. |
| VS0-CF-F16-10 | Authorized ineffective participation | 409 CONFLICT plus VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE; no publication. |
| VS0-CF-F16-11 | Authorized non-Active InfrastructureStack | 409 CONFLICT plus VS0_EXECUTION_TARGET_STACK_UNAVAILABLE; no publication. |
| VS0-CF-F16-12 | Authorized cross-CloudProvider participation/stack graph | 422 VALIDATION_FAILED plus VS0_EXECUTION_TARGET_SCOPE_MISMATCH; no publication. |
| VS0-CF-F16-13 | Create with a duplicate derived name | 409 ALREADY_EXISTS; no publication. |
| VS0-CF-F16-14 | Authorized LIST with an arbitrary Idempotency-Key | 200 accessible Active and Retired targets, ascending UID, no ETag; ignores the key and creates, waits on, reads, and replays no idempotency record. |
| VS0-CF-F16-15 | LIST without applicable executiontarget.read | 403 AUTHORIZATION_DENIED; audited denial. |
| VS0-CF-F16-16 | Authorized item GET with an arbitrary Idempotency-Key | 200 ETag safe projection without facts, result, observer, or handle; ignores the key and creates, waits on, reads, and replays no idempotency record. |
| VS0-CF-F16-17 | Inaccessible direct item GET | Safe 404 RESOURCE_NOT_FOUND plus VS0_AUTHORIZATION_SAFE_DENIAL; audited denial. |
| VS0-CF-F16-18 | Qualify without executiontarget.qualify | 403 AUTHORIZATION_DENIED; audited denial; no publication. |
| VS0-CF-F16-19 | Inaccessible direct qualify action | Safe 404 RESOURCE_NOT_FOUND plus VS0_AUTHORIZATION_SAFE_DENIAL; audited denial. |
| VS0-CF-F16-20 | Qualify, four Supported facts | 200 Qualified/effectiveAvailability Available; FactSet/Result replacement, observedGeneration refresh, new ETag, AuditEvent, completed replay. |
| VS0-CF-F16-21 | Qualify, any Unsupported fact | 200 Rejected/effectiveAvailability Unavailable; records, observedGeneration refresh, new ETag, AuditEvent, completed replay. |
| VS0-CF-F16-22 | Qualify with a missing target-bound fixture | 200 Indeterminate/effectiveAvailability Unavailable; four Unknown facts now/+60s, records, observedGeneration refresh, new ETag, AuditEvent, completed replay. |
| VS0-CF-F16-23 | Qualify from Qualified receives one malformed fact while all captured fences remain current | 500 INTERNAL_ERROR; abort reservation with unchanged committed target, ETag, and FactSet/Result links; publish no AuditEvent or idempotency completion. |
| VS0-CF-F16-24 | Action with a missing If-Match | 412 STALE_RESOURCE_VERSION; no publication; header-form failure precedes body classification. |
| VS0-CF-F16-25 | Action valid If-Match with non-zero body | 400 MALFORMED_REQUEST; no publication. |
| VS0-CF-F16-26 | Qualify Retired target | 409 CONFLICT plus VS0_TARGET_RETIRED; no publication. |
| VS0-CF-F16-27 | Qualify Maintenance target | 409 CONFLICT plus VS0_TARGET_MAINTENANCE; no publication. |
| VS0-CF-F16-28 | Different-key qualify while Qualifying | 409 CONFLICT plus VS0_TARGET_QUALIFICATION_IN_PROGRESS; no publication. |
| VS0-CF-F16-29 | Referenced InfrastructureStack generation changes at qualification commit | 412 STALE_RESOURCE_VERSION plus VS0_TARGET_EPOCH_STALE; abort reservation with unchanged committed target, ETag, and FactSet/Result links; publish no qualification AuditEvent or completion. |
| VS0-CF-F16-30 | Viability fingerprint changes at qualification commit | 412 STALE_RESOURCE_VERSION plus VS0_EXECUTION_TARGET_VIABILITY_STALE; abort reservation with unchanged committed target, ETag, and FactSet/Result links; projection follows current viability and publishes no qualification AuditEvent or completion. |
| VS0-CF-F16-31 | Same-key completed replay after current authorization and safe access | Original successful status/body, application/json, fresh correlation only; no second audit. |
| VS0-CF-F16-32 | Same-key replay requester loses its grant while target remains safe-accessible | Current 403 AUTHORIZATION_DENIED plus one redacted AuditEvent; no stored-result disclosure; audit append failure is INTERNAL_ERROR/500 with no 403 disclosure. |
| VS0-CF-F16-33 | Same namespace with changed digest | 409 CONFLICT plus VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH; original completion unchanged. |
| VS0-CF-F16-34 | Completed record expires after 24 hours | Record unavailable; next valid request is processed as new; no InFlight eviction. |
| VS0-CF-F16-35 | Owner panic after reservation | Reservation aborted and waiters woken without completion; cancelled waiter only detaches. |
| VS0-CF-F16-36 | Valid retire from any Active combination | 200 Retired/Unqualified/effectiveAvailability Unavailable, ETag, refs cleared, AuditEvent, tuple release, completed replay. |
| VS0-CF-F16-37 | Authenticated retire without executiontarget.write | 403 AUTHORIZATION_DENIED; audited denial; no publication. |
| VS0-CF-F16-38 | New-key retire after retirement | 409 CONFLICT plus VS0_TARGET_RETIRED; no resurrection. |
| VS0-CF-F16-39 | Valid create required AuditEvent append fails | 500 INTERNAL_ERROR; no AuditEvent, target creation, or idempotency completion. |
| VS0-CF-F16-40 | Create authorization-denial AuditEvent append fails | 500 INTERNAL_ERROR; no AuditEvent; the original 403 is not disclosed. |
| VS0-CF-F16-41 | Fact expiry with successful required AuditEvent append | effectiveAvailability Unavailable with unchanged ETag; exactly one expiry AuditEvent. |
| VS0-CF-F16-42 | Current fenced maintenance-entry trigger | effectiveAvailability Maintenance; maintenance epoch and resourceVersion/ETag increment; current FactSet/Result links clear; lifecycle-service current-Maintenance marker is active at the new epoch; successful entry audited. |
| VS0-CF-F16-43 | Fenced maintenance clear | Active/Unqualified/effectiveAvailability Unavailable; maintenance epoch and resourceVersion/ETag increment; current FactSet/Result links clear; matching current-Maintenance marker removed; successful clear audited. |
| VS0-CF-F16-44 | HEAD on the item GET path | Transport-only 405 with Allow: GET and X-Content-Type-Options: nosniff; no Problem body, mutation, audit, or idempotency. |
| VS0-CF-F16-45 | Successful GET with any Accept value | Server ignores Accept and emits application/json with nosniff. |
| VS0-CF-F16-46 | Process restart | F0016-owned in-memory targets, records, idempotency reservations, and pending work are absent; inherited FEATURE-0013 AuditEvents are not deleted; no external observer call. |
| VS0-CF-F16-47 | Same-key replay target is no longer safe-accessible | Current safe 404 RESOURCE_NOT_FOUND plus one redacted AuditEvent; no stored-result disclosure; audit append failure is INTERNAL_ERROR/500 with no 404 disclosure. |
| VS0-CF-F16-48 | Problem response with any Accept value | Server ignores Accept and emits application/problem+json with nosniff. |
| VS0-CF-F16-49 | Same principal/key/digest action on another target UID | Independent reservation and outcome; never replays the first target. |
| VS0-CF-F16-50 | Same principal/key/digest create in another derived CloudProvider scope | Independent reservation and outcome; never replays the first scope. |
| VS0-CF-F16-51 | Malformed JSON request body | 400 MALFORMED_REQUEST; no mutation, audit, or completion. |
| VS0-CF-F16-52 | Duplicate top-level JSON member | 400 DUPLICATE_FIELD; no mutation, audit, or completion. |
| VS0-CF-F16-53 | Oversized request body | 400 REQUEST_TOO_LARGE; no mutation, audit, or completion. |
| VS0-CF-F16-54 | Missing required create field | 422 VALIDATION_FAILED; no mutation, audit, or completion. |
| VS0-CF-F16-55 | Invalid typed reference shape | 422 VALIDATION_FAILED; no mutation, audit, or completion. |
| VS0-CF-F16-56 | targetClass other than synthetic-iaas | 422 VALIDATION_FAILED; no mutation, audit, or completion. |
| VS0-CF-F16-57 | Authorized create with adapterAuthorityRef | 400 UNKNOWN_FIELD; no mutation, audit, or completion. |
| VS0-CF-F16-58 | Create/action with an empty Idempotency-Key | 400 MALFORMED_REQUEST; no publication. |
| VS0-CF-F16-59 | Create/action with a malformed Idempotency-Key | 400 MALFORMED_REQUEST; no publication. |
| VS0-CF-F16-60 | Create/action with an overlong Idempotency-Key | 400 MALFORMED_REQUEST; no publication. |
| VS0-CF-F16-61 | LIST with missing authentication | 401 AUTH_REQUIRED; no mutation, audit, or replay record. |
| VS0-CF-F16-62 | Item GET with missing authentication | 401 AUTH_REQUIRED; no mutation, audit, or replay record. |
| VS0-CF-F16-63 | Qualify with missing authentication | 401 AUTH_REQUIRED; no mutation, audit, or replay record. |
| VS0-CF-F16-64 | Retire with missing authentication | 401 AUTH_REQUIRED; no mutation, audit, or replay record. |
| VS0-CF-F16-65 | Create with a live duplicate backing tuple | 409 ALREADY_EXISTS; no publication. |
| VS0-CF-F16-66 | Inaccessible direct retire action | Safe 404 RESOURCE_NOT_FOUND plus VS0_AUTHORIZATION_SAFE_DENIAL; audited denial. |
| VS0-CF-F16-67 | Qualify receives a duplicate named fact while all captured fences remain current | 500 INTERNAL_ERROR; abort reservation with unchanged committed target, ETag, and FactSet/Result links; publish no AuditEvent or completion. |
| VS0-CF-F16-68 | Qualify receives a missing named fact while all captured fences remain current | 500 INTERNAL_ERROR; abort reservation with unchanged committed target, ETag, and FactSet/Result links; publish no AuditEvent or completion. |
| VS0-CF-F16-69 | Qualify receives malformed fact provenance while all captured fences remain current | 500 INTERNAL_ERROR; abort reservation with unchanged committed target, ETag, and FactSet/Result links; publish no AuditEvent or completion. |
| VS0-CF-F16-70 | Action with a malformed If-Match | 412 STALE_RESOURCE_VERSION; no publication; header-form failure precedes body classification. |
| VS0-CF-F16-71 | Action with a weak If-Match | 412 STALE_RESOURCE_VERSION; no publication; header-form failure precedes body classification. |
| VS0-CF-F16-72 | Action with wildcard If-Match | 412 STALE_RESOURCE_VERSION; no publication; header-form failure precedes body classification. |
| VS0-CF-F16-73 | Action with multiple If-Match headers | 412 STALE_RESOURCE_VERSION; no publication; header-form failure precedes body classification. |
| VS0-CF-F16-74 | Action with a syntactically valid stale If-Match | 412 STALE_RESOURCE_VERSION; no publication; current-version comparison follows replay/reservation and precedes lifecycle evaluation. |
| VS0-CF-F16-75 | Maintenance epoch changes because Maintenance entry wins during qualification | 412 STALE_RESOURCE_VERSION plus VS0_TARGET_EPOCH_STALE; no qualification conclusion, AuditEvent, or completion; the independently winning Maintenance transition is audited and supplies its own target mutation and ETag. |
| VS0-CF-F16-76 | Completed record is selected for deterministic capacity eviction | Record unavailable; next valid request is processed as new; no InFlight eviction. |
| VS0-CF-F16-77 | Server shutdown after reservation | Stop expiry/Maintenance trigger acceptance, abort qualification and idempotency reservations, wake waiters, then stop HTTP serving; no completion or post-stop trigger AuditEvent; last committed target projection, ETag, and links remain unchanged. |
| VS0-CF-F16-78 | Owner cancellation after reservation | Reservation aborted and waiters woken without completion; cancelled waiter only detaches. |
| VS0-CF-F16-79 | Fact expiry with required AuditEvent append failure | effectiveAvailability remains Unavailable with unchanged ETag; one expiry audit retries each injected-clock second until success. |
| VS0-CF-F16-80 | Stale fenced maintenance-entry trigger | No state change, AuditEvent, or record publication. |
| VS0-CF-F16-81 | Same-key waiter wakes after its owner aborts on a stale InfrastructureStack-generation qualification | No cached completion; waiter repeats current authorization and safe access, then may establish a new reservation if its current request remains valid; the owner's target remains unchanged as in F16-29. |
| VS0-CF-F16-82 | Required qualification AuditEvent append failure after reservation while all captured fences remain current | 500 INTERNAL_ERROR; abort reservation with unchanged committed target, ETag, and FactSet/Result links; wake waiters without AuditEvent or completion. |
| VS0-CF-F16-83 | Stale fenced maintenance-clear trigger | No state change, AuditEvent, or record publication. |
| VS0-CF-F16-84 | Authorized create with a SecretRef extension field | 400 UNKNOWN_FIELD; no mutation, audit, or completion. |
| VS0-CF-F16-85 | Qualify Maintenance target with syntactically valid stale If-Match | 412 STALE_RESOURCE_VERSION; current-version comparison precedes Maintenance state evaluation; no publication. |
| VS0-CF-F16-86 | Qualify with malformed If-Match and non-zero body | 412 STALE_RESOURCE_VERSION; header-form failure precedes body classification; no publication. |
| VS0-CF-F16-87 | Same-key completed qualify after target ETag changes | Original successful 200 JSON response; If-Match exclusion permits replay after current authorization and safe access; no second audit. |
| VS0-CF-F16-88 | Retire wins while qualify is in flight | Qualifying caller gets 409 CONFLICT plus VS0_TARGET_RETIRED; no qualification result, completion, or qualification AuditEvent; winning retirement is audited. |
| VS0-CF-F16-89 | Maintenance entry wins while qualify is in flight | Qualifying caller gets 409 CONFLICT plus VS0_TARGET_MAINTENANCE; no qualification result, completion, or qualification AuditEvent; winning maintenance entry is audited. |
| VS0-CF-F16-90 | Qualify with fixture-declared logical timeout | 200 Indeterminate/effectiveAvailability Unavailable; four Unknown facts now/+60s, records, AuditEvent, completed replay; no elapsed network wait. |
| VS0-CF-F16-91 | Authorized create with client metadata.resourceVersion | 403 AUTHORIZATION_DENIED plus VS0_SYSTEM_OWNED_FIELD_WRITE; audited denial. |
| VS0-CF-F16-92 | Authorized create with client metadata.generation | 403 AUTHORIZATION_DENIED plus VS0_SYSTEM_OWNED_FIELD_WRITE; audited denial. |
| VS0-CF-F16-93 | Create with invalid authentication | 401 AUTH_REQUIRED; no mutation, audit, or replay record. |
| VS0-CF-F16-94 | LIST with invalid authentication | 401 AUTH_REQUIRED; no mutation, audit, or replay record. |
| VS0-CF-F16-95 | Item GET with invalid authentication | 401 AUTH_REQUIRED; no mutation, audit, or replay record. |
| VS0-CF-F16-96 | Qualify with invalid authentication | 401 AUTH_REQUIRED; no mutation, audit, or replay record. |
| VS0-CF-F16-97 | Retire with invalid authentication | 401 AUTH_REQUIRED; no mutation, audit, or replay record. |
| VS0-CF-F16-98 | Create with a name reserved by a Retired target | 409 ALREADY_EXISTS; retirement released the backing tuple but not the scope/name reservation. |
| VS0-CF-F16-99 | Valid retirement required AuditEvent append fails | 500 INTERNAL_ERROR; no AuditEvent, retirement, tuple release, ETag advance, or idempotency completion. |
| VS0-CF-F16-100 | Current maintenance-entry required AuditEvent append fails | 500 INTERNAL_ERROR; no AuditEvent, maintenance epoch/resourceVersion advance, link clear, or transition publication. |
| VS0-CF-F16-101 | Current maintenance-clear required AuditEvent append fails | 500 INTERNAL_ERROR; no AuditEvent, maintenance epoch/resourceVersion advance, link clear, or transition publication. |
| VS0-CF-F16-102 | Inaccessible direct item GET required AuditEvent append fails | 500 INTERNAL_ERROR; no AuditEvent; the original safe 404 is not disclosed. |
| VS0-CF-F16-103 | Same-key replay authorization denial AuditEvent append fails | 500 INTERNAL_ERROR; no AuditEvent; the original replay-time 403 is not disclosed. |
| VS0-CF-F16-104 | Same-key replay safe-denial AuditEvent append fails | 500 INTERNAL_ERROR; no AuditEvent; the original replay-time 404 is not disclosed. |
| VS0-CF-F16-105 | Collection create with an If-Match header | 400 MALFORMED_REQUEST; no mutation, AuditEvent, or idempotency completion. |
| VS0-CF-F16-106 | Authenticated URI-query equivalence family: POST collection, GET collection, GET item, POST qualify, and POST retire, each with `?x=1` | For every named valid method/path pair, 400 MALFORMED_REQUEST before safe resolution; no reference lookup, mutation, AuditEvent, idempotency reservation, or completion. |
| VS0-CF-F16-107 | Authenticated item GET with a malformed canonical UID segment | 400 MALFORMED_REQUEST before safe resolution; no reference lookup, mutation, AuditEvent, idempotency reservation, or completion. |
| VS0-CF-F16-108 | Item GET with trailing slash | Transport-only 404 without redirect or Problem body, with X-Content-Type-Options: nosniff; no authentication, mutation, audit, or idempotency. |
| VS0-CF-F16-109 | Authenticated qualify action with a malformed canonical UID segment | 400 MALFORMED_REQUEST before safe resolution; no reference lookup, mutation, AuditEvent, idempotency reservation, or completion. |
| VS0-CF-F16-110 | Authenticated retire action with a malformed canonical UID segment | 400 MALFORMED_REQUEST before safe resolution; no reference lookup, mutation, AuditEvent, idempotency reservation, or completion. |
| VS0-CF-F16-111 | Direct item GET of a safe-accessible target without executiontarget.read | After UID-only safe resolution derives the target scope, 403 AUTHORIZATION_DENIED; one redacted AuditEvent; no target disclosure, mutation, or idempotency record. |
| VS0-CF-F16-112 | Denial AuditEvent append-failure equivalence family: LIST without read; direct GET without read; qualify without qualify grant; retire without write grant; inaccessible direct qualify; inaccessible direct retire | For every named member, 500 INTERNAL_ERROR; its original 403/404 and violation are not disclosed; no AuditEvent, mutation, lifecycle/record publication, or idempotency completion. |
| VS0-CF-F16-113 | HEAD on collection path | Transport-only 405 with Allow: GET, POST and X-Content-Type-Options: nosniff; no Problem body, authentication, audit, idempotency, mutation, or lifecycle effect. |
| VS0-CF-F16-114 | HEAD equivalence family: qualify action and retire action paths | For each named action path, transport-only 405 with Allow: POST and X-Content-Type-Options: nosniff; no Problem body, authentication, audit, idempotency, mutation, or lifecycle effect. |
| VS0-CF-F16-115 | Trailing-slash equivalence family: POST collection, GET collection, POST qualify action, and POST retire action | For each named method/path pair, transport-only 404 without redirect or Problem body, with X-Content-Type-Options: nosniff; no authentication, audit, idempotency, mutation, or lifecycle effect. |
| VS0-CF-F16-116 | PUT on collection path | Transport-only 405 with Allow: GET, POST and X-Content-Type-Options: nosniff; no Problem body, authentication, audit, idempotency, mutation, or lifecycle effect. |
| VS0-CF-F16-117 | PUT on item path | Transport-only 405 with Allow: GET and X-Content-Type-Options: nosniff; no Problem body, authentication, audit, idempotency, mutation, or lifecycle effect. |
| VS0-CF-F16-118 | PUT equivalence family: qualify action and retire action paths | For each named action path, transport-only 405 with Allow: POST and X-Content-Type-Options: nosniff; no Problem body, authentication, audit, idempotency, mutation, or lifecycle effect. |
| VS0-CF-F16-119 | Fresh Qualified target backing becomes non-viable | GET projects effectiveAvailability Unavailable with unchanged target ETag and no F0016 mutation or AuditEvent. |
| VS0-CF-F16-120 | Fresh Qualified target backing returns viable before fact expiry | GET projects effectiveAvailability Available with unchanged target ETag and no F0016 mutation or AuditEvent. |
| VS0-CF-F16-121 | Qualify starts from Active/Qualified with current fences | Create an internal Qualifying reservation; retain the public committed target projection, ETag, and FactSet/Result links, with no AuditEvent or idempotency completion. |
| VS0-CF-F16-122 | Current-fence observer fault after qualification starts from Active/Qualified | 500 INTERNAL_ERROR; abort reservation with unchanged committed target, ETag, and FactSet/Result links; publish no AuditEvent or idempotency completion. |

### 3. Decision-to-authority/proof matrix

| AD | Authority target | Owner / trigger | Annex proof | Fail-closed assertion |
|---|---|---|---|---|
| 01 | VS0-SCHEMA-015 / feature authority | F0016 feature authority; API server projection | 01, 16, 20–22, 36 | sole canonical target/projection |
| 02 | VS0-SCHEMA-015 / registry | API server create command; lifecycle service | 01–04, 10–13, 54–57, 65, 90–91, 97 | closed create fields and uniqueness |
| 03 | VS0-SCHEMA-015 plus state / registry | API server create/qualify safe resolver; response projection | 09–12, 29–30, 119–120 | safe backing and viability |
| 04 | contract specification / registry routes | Pre-ServeMux guard and API server mux | 01, 14, 16, 20–22, 36, 44–45, 61–64, 66, 90, 93–96, 105–110, 113–118 | five exact registrations; guard owns all unmatched-method and trailing-slash transport outcomes |
| 05 | registry authorization contract | API server authorization and safe resolver | 07–09, 15, 17–19, 32, 37, 47, 61–66, 92–96, 111–112 | grants and safe denial |
| 06 | contract specification / registry | API server request/response boundary | 01–06, 14, 16, 20–22, 24–25, 36, 44–45, 48, 51–60, 70–74, 83–86, 89–91, 105–110, 113–118 | grammar, ETag, representation |
| 07 | contract specification / readiness checker | Pre-ServeMux guard and API server request pipeline | 02–13, 18–30, 51–60, 67–75, 84–86, 90–91, 105–111, 113–118 | applicable-stage precedence and transport guard |
| 08 | registry / conformance matrix | Idempotency coordinator | 01, 20–23, 31–35, 46–47, 49–50, 76–78, 81–82, 86–89, 99 | replay, waiters, retention |
| 09 | VS0-SCHEMA-015 / feature authority | Server configuration observer binding | 02, 46, 57, 85 | no adapter authority/credentials |
| 10 | VS0-SCHEMA-016 / fake contract | Synthetic observer fixture and injected clock | 20–23, 46, 67–69, 90 | one deterministic no-call observer |
| 11 | VS0-SCHEMA-016/017 / registry | Qualification engine proposal; lifecycle-service commit | 20–23, 67–69, 87–90 | facts and conclusions |
| 12 | VS0-SCHEMA-016/017 / registry | Synthetic observer; injected-clock expiry trigger | 20–23, 41, 67–69, 79, 90 | provenance and freshness |
| 13 | VS0-SCHEMA-015/016/017 / feature authority | Lifecycle service; process restart | 16, 20–23, 36, 46 | links, redaction, restart |
| 14 | VS0-STATE-004 / feature authority | ExecutionTargetLifecycleService | 01, 20–22, 26–30, 36, 39, 41–43, 75, 77, 79–83, 87–88, 99–101, 121–122 | sole lifecycle committer |
| 15 | VS0-STATE-004 / registry | Lifecycle service; response projection | 01, 20–22, 26, 36, 38, 41–43, 79–83, 87–88, 99–101, 119–122 | valid axes, fences, and projection |
| 16 | VS0-STATE-004 / registry | Lifecycle-service qualify command | 20–30, 75, 81–82, 87–88, 121–122 | qualification transitions |
| 17 | VS0-STATE-004 / formal model | Injected-clock expiry trigger; audit coordinator; shutdown boundary | 41, 77, 79 | expiry/audit retry and shutdown ordering |
| 18 | VS0-STATE-004 / registry | Synthetic-observer fixture maintenance trigger; lifecycle service | 01, 29–30, 42–43, 75, 80, 83, 88, 100–101 | fences and maintenance |
| 19 | VS0-STATE-004 / registry | Lifecycle-service retire command | 26, 36, 38, 87, 97 | retirement |
| 20 | VS0-STATE-004 / formal model | Lifecycle service and idempotency coordinator | 23, 28–35, 75–78, 81–82, 87–88, 121–122 | stale/concurrent work |
| 21 | audit matrix / registry | Audit coordinator | 01, 03–04, 08–09, 15, 17–23, 32, 36–43, 47, 66, 79–83, 87–88, 99–104, 111–112 | audit ordering and finite denial-family coverage |
| 22 | problem/violation registry | API server and lifecycle service | 10–12, 26–30, 38, 75 | closed F0016-local violations |
| 23 | schema/projection contract | API server projection writer | 01, 14, 16, 20–22, 36, 45, 48 | redaction/representation |
| 24 | feature authority / Structurizr | F0016 feature authority | 46 | no real adapter/Crossplane |
| 25 | formal models plus conformance matrix | TLC for state/fence invariants; conformance suite for HTTP cases | 01–122 | all local behavior has one proof mechanism |
| 26 | closure matrix / readiness checker | Architecture-readiness checker | 01–122 | every case exact or explicitly enumerated equivalent family |
| 27 | feature control / route audit | Requirements/design admission gates | 01–122 | requirements/design admission |
| 28 | task inventory / feature control | Task-scope inventory gate | 01–122 | task admission |

### 4. Mandatory downstream admission

1. Registry, contract specification, traceability, feature authority, closure
   matrix, readiness checker, contract checker, and Structurizr must represent
   only this register and annex.
2. Requirements may copy canonical ledgers, not decide fields, status, writer,
   route, problem/violation, audit effect, replay class, or missing outcome.
3. Design must map the complete five-route input-to-outcome pipeline and
   fail-closed annex proof before tasks.
4. Tasks must derive source-based writable paths and an acyclic graph.
   Implementation cannot repair an authority gap.

## Rationale

The current registry reserves F0016 resource names but leaves observable
behavior incomplete. That forced equivalent decisions into F0015 requirements,
design, and tasks and caused repeated review loops. This closure establishes
one bounded, production-shaped fake IaaS qualification flow while keeping
credentials, external calls, placement, plugin execution, and real lifecycle
coordination explicitly out of F0016.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse FEATURE-0012 API grammar, typed references, Problem Details, media
    semantics, and concurrency contract.
  - Reuse FEATURE-0013 AuditEvent atomic-publication semantics.
  - Reuse FEATURE-0015 server-resolved authorization, safe-denial, and
    idempotency profile; do not reuse its resource routes or lifecycle.
  - Crossplane remains a future realization candidate only, with no F0016
    dependency.
- Sovrunn-owned responsibility summary:
  - Normalize target observations, evaluate the closed qualification profile,
    publish current target status safely, and fence stale work.
- Non-goals summary:
  - No real target adapter, provider-native type, raw secret, external call,
    external persistence, placement, execution, plugin, customer target API,
    maintenance-notice resource, or IAM resource.

## Phase impact

- Current phase allowed? Yes.
- Current phase boundary impact:
  - Phase 2R gains only deterministic observation and qualification. Real
    realization remains FEATURE-0024/Phase 3 work, and target credentials are
    deferred until a real target class is separately approved.

## Conflict check

- Conflicts with accepted DEC/RFC? No, after correcting the placeholder
  F0016 registry/schema/state entries.
- Resolution required: Approved ADH and coherent authority update; no baseline
  replacement, DEC, or RFC is required.

## Required action

- Update architecture doc and FEATURE-0016 feature authority.
- Update registry, contract specification, traceability, closure matrix,
  Slice 0 steering, Structurizr dynamic/container view as necessary, and
  FEATURE-0016 control manifest.
- Create/update deterministic readiness and generic feature-contract checks.
- Create/update the deterministic F0016 route-executability and task-scope
  inventory checks.
- Add F0016 formal models and formal-check wiring.
- Do **not** update Kiro requirements.md, design.md, tasks.md, or Go code.

## Impacted files

The following strict read/write boundary is the complete impacted-file list for
this handoff.

### Strict read/write boundary

Kiro must first read the Architecture Operating System sources required by
`AGENTS.md`. Apart from those mandatory governance reads, it may read only:

- this handoff;
- `docs/reviews/architecture-readiness/FEATURE-0016-architecture.md`;
- `docs/reviews/architecture-readiness/FEATURE-0016-architecture-starting-dossier.md`;
- `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`;
- `docs/features/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`;
- `.automation/features/FEATURE-0016.control.json`;
- `docs/reviews/architecture-readiness/FEATURE-0016-architecture-closure-matrix.md`;
- `docs/reviews/reuse-assessments/FEATURE-0016-reuse-assessment.md`;
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`;
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`;
- `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/architecture/adapter-boundary-model.md`;
- `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`;
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`;
- `.kiro/steering/slice0-contract.md`;
- `docs/diagrams/structurizr/workspace.dsl`;
- `scripts/feature-contract-check.py`;
- `scripts/vs000-contract-check.py`;
- `scripts/feature-0015-architecture-readiness-check.py`.
- `scripts/feature-control.py`.
- `scripts/generic-feature-boundary-check.py`.

If approved, Kiro may write only:

- this handoff;
- `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`;
- `docs/features/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`;
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`;
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`;
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`;
- `docs/reviews/architecture-readiness/FEATURE-0016-architecture-closure-matrix.md`;
- `docs/reviews/reuse-assessments/FEATURE-0016-reuse-assessment.md`;
- `.kiro/steering/slice0-contract.md`;
- `.automation/features/FEATURE-0016.control.json`;
- `docs/diagrams/structurizr/workspace.dsl`;
- `scripts/feature-0016-architecture-readiness-check.py`;
- `scripts/feature-0016-route-executability-check.py`;
- `scripts/feature-0016-task-scope-inventory.py`;
- `scripts/feature-contract-check.py`;
- `scripts/vs000-contract-check.py`;
- `scripts/run-feature-0016-formal-checks.sh`;
- `docs/formal/feature-0016/ExecutionTargetLifecycle.tla`;
- `docs/formal/feature-0016/ExecutionTargetLifecycle.cfg`;
- `docs/formal/feature-0016/ExecutionTargetFence.tla`;
- `docs/formal/feature-0016/ExecutionTargetFence.cfg`.

Kiro must not read or modify `.kiro/specs/**`, Go source, existing feature
state files, DEC files, baseline files, roadmap files, or any file outside
these lists. It must stop if a required authority cannot be expressed within
this boundary.

## Impacted features

- FEATURE-0016: receives the complete active target-observation and
  qualification contract.
- FEATURE-0015: consumed only as a fixed reference owner; no change.
- FEATURE-0017, FEATURE-0019, FEATURE-0022, FEATURE-0023, FEATURE-0024:
  future consumers; no activation or modification.

## Acceptance criteria for Kiro update

- [ ] VS0-SCHEMA-015 has exactly the F0016 create-field boundary, derived
  CloudProvider scope, Active/Retired lifecycle, four-value persisted public
  qualification plus the lifecycle-service-only Qualifying reservation,
  maintenance epoch, observed generation, and current FactSet/Result links.
  It has no persisted `status.availability`; `effectiveAvailability` is the
  three-value response projection computed by the lifecycle service.
- [ ] The update replaces the current `VS0-SCHEMA-015..017` and
  `VS0-STATE-004` placeholders rather than retaining parallel or compatibility
  definitions. Active `status.availability`, `Draining`, and public
  `Qualifying` are absent, rejected, and covered by fail-closed checks.
- [ ] VS0-SCHEMA-016 expresses all four named facts with exactly
  Supported/Unsupported/Unknown truth and complete target/generation/epoch/
  observer/freshness provenance.
- [ ] VS0-SCHEMA-017 is an immutable internal result record, not a transient
  public result; it binds exact target, fact set, generation, epoch, profile,
  outcome, reasons, and evaluation time.
- [ ] VS0-STATE-004 expresses every decision above, including the exact
  valid-combination and transition-authority matrix, retirement, expiry,
  maintenance fence entry/clear, and stale non-publication; it does not
  activate Draining.
- [ ] The registry/specification define five routes, exact headers/media/
  statuses, safe graph resolution, F0016 Problem/violation outcomes,
  idempotency retention, audit matrix, projections, and no 405 Problem code.
- [ ] The registry contains an exact F0016-local conformance suite; no shared
  or downstream case is counted as F0016-local proof.
- [ ] Deterministic readiness and feature-contract checks fail closed on all
  decision clauses. The lifecycle TLC model proves Retired terminality,
  availability truth/priority, and non-public Qualifying; the fence/publication
  model proves audit-before-conclusion/completion and stale/cancelled/failed
  non-mutation. Both run as part of the architecture gate.
- [ ] The authority compatibility/proof matrix has exactly one F0016-local
  proof target for every active observable decision clause, and the
  requirements context manifest has measured headroom.
- [ ] The route-executability audit is wired to the design admission gate and
  fails on a missing route outcome or precedence rule.
- [ ] The source-derived task-scope inventory is wired to the tasks admission
  gate and fails on a missing writable/test/dependency path or a cyclic graph.
- [ ] `make feature-0016-architecture-readiness`,
  `make feature-contract-check FEATURE=FEATURE-0016`, formal checks,
  `make vs000-contract-check`, `make phase2r-drift-check`,
  `git diff --check`, `mkdocs build --strict`, and `make structurizr-check`
  pass.
- [ ] No requirements, design, tasks, Go code, baseline, DEC, roadmap, or
  unallowlisted file is changed.

## Explicit instructions to Kiro

- Do not apply this handoff until founder approval is recorded below.
- Do not introduce a SecretRef extension, real adapter, external call,
  persistence dependency, placement behavior, plugin execution, maintenance
  notice, or F0018 IAM resource.
- Do not reintroduce ResourcePool, ProviderCapability, a generic Provider, or
  a target subtype.
- Do not choose requirements/design mechanics such as Go data structures,
  locks, package layout, or scheduler implementation.
- Do not introduce any decision beyond this handoff. Stop and report a bounded
  authority gap instead.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-20
- Notes: Founder approved the bounded FEATURE-0016 adapter and ExecutionTarget contract closure for strict allowlisted application.
## Application Record

- Applied by: Kiro (architecture/spec update agent)
- Application date: 2026-08-20
- Written artifacts: `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`, `docs/features/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`, `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`, `docs/architecture/vertical-slices/VS-000-contract-specification.md`, `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`, `docs/reviews/architecture-readiness/FEATURE-0016-architecture-closure-matrix.md`, `docs/reviews/reuse-assessments/FEATURE-0016-reuse-assessment.md`, `.kiro/steering/slice0-contract.md`, `.automation/features/FEATURE-0016.control.json`, `scripts/feature-0016-architecture-readiness-check.py`, `scripts/feature-0016-route-executability-check.py`, `scripts/feature-0016-task-scope-inventory.py`, `scripts/feature-contract-check.py` (reviewed; no FEATURE-0016 catalog added — see bounded-authority note below), `scripts/vs000-contract-check.py`, `scripts/run-feature-0016-formal-checks.sh`, `docs/formal/feature-0016/ExecutionTargetLifecycle.tla`, `docs/formal/feature-0016/ExecutionTargetLifecycle.cfg`, `docs/formal/feature-0016/ExecutionTargetFence.tla`, `docs/formal/feature-0016/ExecutionTargetFence.cfg`.
- `docs/diagrams/structurizr/workspace.dsl` left unchanged: ADH-2026-058 establishes a Sovrunn-owned, no-adapter/no-external-relationship observation boundary inside the existing `sovrunn.placementEngine`/`sovrunn.registry` containers; it adds no new system boundary, container, plugin plane, external relationship, or dynamic flow requiring a diagram update.
- No Makefile edit was required: `feature-0016-architecture-readiness`, `feature-contract-check`, and `vs000-contract-check` targets already existed and point at the corresponding allowlisted scripts. `feature-0016-formal-check` and a Makefile wrapper for `run-feature-0016-formal-checks.sh` are not yet wired into the root `Makefile`; see the bounded-authority-gap note in the Kiro final report for this application run.
- No Kiro `.kiro/specs/**`, Go source, `.automation/state/**`, DEC file, `CURRENT_ARCHITECTURE_BASELINE.md`, or roadmap file was read or modified.
