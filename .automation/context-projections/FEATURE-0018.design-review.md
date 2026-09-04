---
doc_type: ai_context_projection
feature: FEATURE-0018
stage: design-review
authority: non-authoritative-exact-excerpts
generated: true
---

# FEATURE-0018 Design-Review Context Projection

This is a non-authoritative, mechanically generated projection of exact
approved-source excerpts. The FEATURE architecture and controlling handoff
package remain authoritative. A source-hash mismatch blocks regeneration.

## Source: `docs/reviews/architecture-decision-handoffs/ADH-2026-072-feature-0018-design-authority-and-audited-evaluation-clarification.md`

Approved SHA-256: `f6b50ea9faba3d8dc09d86124b574ccebefad27733db1ac3c4eaf98361acd5ca`

### Approved design authority, audited-evaluation boundary and reviewer classification

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
## Decision or proposed decision

### 1. FEATURE-0013 and FEATURE-0018 ownership boundary

FEATURE-0013 owns the shared DecisionProfile, EvaluationResult, DecisionRecord,
and AuditEvent contracts; their carrier types; their common validation, projection,
validity, authority, reason, obligation, and audit-linkage rules; and the approved
extension mechanism used to register the three FEATURE-0018 profiles and the bounded
AuditEvent taxonomy.

FEATURE-0013 does not provide FEATURE-0018 with a transaction, acceptance service,
publisher, durable store, outbox, queue, allocator, lock manager, idempotency store,
or prepared-evidence commit API. The phrase “FEATURE-0013 local atomic boundary” in
F18-RD-20 means that FEATURE-0018 must preserve the inherited evidence contract and
its acceptance-before-publication invariant. It does not name an existing callable
FEATURE-0013 service and does not transfer mutation or transaction ownership to
FEATURE-0013.

FEATURE-0018 may implement a deterministic, feature-local, process-local,
non-durable in-memory publication mechanism that stages and publishes, as one
all-or-nothing aggregate change:

```text
FEATURE-0018 domain mutation(s), when applicable
+ required FEATURE-0013-compatible DecisionRecord/AuditEvent carrier values
+ completed caller idempotency result, when applicable
```

That mechanism is FEATURE-0018 implementation and conformance machinery. It must
consume FEATURE-0013 carrier construction and validation rules without creating a
second evidence envelope or writer, must perform zero external I/O, and must never
claim production durability or a FEATURE-0013 acceptance service.

### 2. Authorization candidate and audited evaluation boundary

Stages before concurrency/idempotency may compute one sealed internal authorization
candidate. A candidate is provisional implementation data only. It is not an
`AuthorizationResult`, DecisionRecord, AuditEvent, bearer value, completed
authorization evaluation, or downstream authority.

A candidate becomes a completed authorization evaluation only after every
operation-applicable check through F18-RD-21's concurrency/idempotency stage has
succeeded and the operation is ready to produce one exact current `Allow | Deny`
AuthorizationResult. At that point, the required `authorization.evaluated`
obligation and any mandatory DecisionRecord are prepared, validated, accepted, and
published before the result is released. Evidence failure returns the existing safe
internal error, publishes no result or completed replay, and grants no downstream
authority.

The following stage-10 terminal or coordination outcomes do not produce an
AuthorizationResult and therefore do not complete an authorization evaluation:

- request-binding `Conflict`;
- terminal optimistic-concurrency/CAS loss;
- `Wait` or retry/re-evaluation direction; and
- controller effective-once `NoOp`.

Their already-registered conflict, lifecycle, concurrency, or mutation-attempt
behavior remains unchanged; this handoff adds no new AuditEvent type or error.
Discarding a provisional candidate on one of these paths does not violate F18-RD-19,
because no completed `Allow | Deny` evaluation exists.

A completed-result replay is different: after current authentication,
authorization, and safe target visibility are re-evaluated and stage 10 selects the
matching stored result, that current evaluation must accept and publish exactly its
required authorization evidence before the stored successful result is returned.
Replay never extends the original effect or validity.

### 3. Design-owned implementation mechanics

Within the exact F18-RD, REQ, AC, route, writer, error, dependency, and Phase 2R
boundaries, FEATURE-0018 design owns the concrete implementation mechanics necessary
to make the approved behavior executable. This includes:

- package decomposition and acyclic import direction;
- sealed interfaces and owner-private constructors;
- immutable fixture construction, copying, validation, and deterministic clocks;
- process-local state roots, snapshots, locks, reservations, waiters, optimistic
  version checks, and idempotency tables;
- the FEATURE-0018-local non-durable unit-of-work protocol, transaction modes,
  context-bound change receipts, seal invariants, commit/abort, panic cleanup, and
  shutdown coordination;
- FEATURE-0013-compatible carrier preparation and validation wiring;
- literal route shapes for the already-closed caller operation surface; and
- static architecture-checker symbols, package/call rules, diagnostics, fixtures,
  and build integration.

These are design choices, not new product semantics. Design may not use them to add
or change a resource, action, field authority, scope, target binding, lifecycle,
eligibility rule, approval rule, audit obligation, error, route operation, writer,
dependency, persistence guarantee, external integration, or adjacent-feature
behavior. If a mechanism requires such a choice, design must stop with
`ARCHITECTURE_DECISION_REQUIRED`.

The requirements statement that no other concern is delegated means no other
observable product semantic or ownership decision is delegated. It does not prohibit
the design mechanics listed above.

### 4. Authority consistency and stale transcription

F18-RD-02, the approved requirements, the canonical contract catalog, the canonical
model, and the VS-000 registry must agree exactly for the fields they register. No
design may select one as the winner. Any mismatch stops with
`ARCHITECTURE_DECISION_REQUIRED`; it is never resolved by source precedence inside
FEATURE-0018 design.

The existing design sentence saying that the registry wins is a stale transcription
and must be removed. Correcting it does not reopen F18-RD-02.

### 5. FEATURE-0018 design-review classification rule

Every blocking or revision finding in a FEATURE-0018 design review must use exactly
one of these classifications:

| Classification | Exact use |
|---|---|
| `STALE_TRANSCRIPTION` | The design contradicts or duplicates a closed authority. Correct the design by exact citation/transcription; do not reopen the decision. |
| `DESIGN_EXECUTABILITY` | Approved behavior is closed, but the proposed package, API, state, concurrency, idempotency, evidence, checker, or test mechanism is incomplete or internally inconsistent. Correct design only. |
| `REQUIREMENT_GAP` | An approved architecture statement is absent or weakened in requirements. Stop design and reconcile requirements without inventing semantics. |
| `ARCHITECTURE_CLARIFICATION_REQUIRED` | Loaded approved authorities do not determine an observable semantic or ownership outcome. Stop and obtain bounded human-approved clarification. |
| `OUT_OF_SCOPE` | The proposed or requested concern belongs to an excluded/future feature, dependency owner, production runtime, or external system. Remove or defer it without designing the future feature. |

A reviewer may not classify a closed decision as a requirement or architecture gap
without citing the exact contradictory or missing authority and the exact target
text. A finding about an implementable mechanism must use
`DESIGN_EXECUTABILITY`, even when the mechanism protects a security invariant.
Repeated wording, preference, or an alternative implementation is not evidence that
a closed architecture decision is open.
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

## Source: `docs/reviews/architecture-decision-handoffs/ADH-2026-073-feature-0018-exceptionproposal-terminal-decision-audit-event.md`

Approved SHA-256: `3cae0d0f138aa80f11d8db44db2b3e6b7a49e4153a5eb5b0aacbb1577ed4b08b`

### Approved terminal-ExceptionProposal AuditEvent registration

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
## Decision or proposed decision

Register exactly one additional FEATURE-0018 AuditEvent taxonomy type:

```text
exceptionproposal.decided
```

`exceptionproposal.decided` is emitted exactly once for every terminal
`ExceptionProposal` `Grant | Deny` outcome. It is linked to the exact immutable
proposal, its terminal reason, the mandatory `exception-decision/v1` DecisionRecord,
and the containing operation/correlation evidence using the existing FEATURE-0013
AuditEvent envelope and linkage semantics.

For terminal `Grant`, `exceptionproposal.decided` and `exceptiongrant.issued` are
separate required events in the same final exception atomic boundary:

```text
ExceptionProposal terminal Grant
  -> exceptionproposal.decided
  -> exceptiongrant.issued
  -> exactly one exception-decision/v1 DecisionRecord
  -> exact immutable ExceptionGrant
```

For terminal `Deny`, the boundary publishes no `ExceptionGrant` and no
`exceptiongrant.issued` event:

```text
ExceptionProposal terminal Deny
  -> exceptionproposal.decided
  -> exactly one exception-decision/v1 DecisionRecord
  -> no ExceptionGrant
```

The existing event meanings remain fixed:

```text
exceptiongrant.proposed  -> accepted proposal submission only
exceptiongrant.issued    -> immutable ExceptionGrant issuance only
exceptiongrant.revoked   -> linked revocation evidence only
exceptiongrant.expired   -> expiry evidence only
```

This is a registration through FEATURE-0013's already-approved extension mechanism.
It does not add a FEATURE-0013 carrier type, validation rule, service, transaction,
writer, durable store, or dependency change. FEATURE-0018 remains responsible for
its feature-local, deterministic, non-durable publication mechanics and for the
conformance evidence of this registered type.
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

## Source: `docs/reviews/architecture-decision-handoffs/ADH-2026-071-feature-0018-conformance-executability-reconciliation.md`

Approved SHA-256: `231b8872d83c5ad46331ef8a0e2e24127cb2ef177f0685f829617e648387b0b6`

### Conformance authority, non-semantic boundary and identifier invariants

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
### 1. Authority and non-semantic boundary

ADH-2026-071 authorizes only conformance identifiers, traceability mappings, and
faithful transcription of already-approved semantics. A registration below is a
stable name for existing F18-RD/F18-SCN behavior; it is not a new behavioral
source. If application would require Kiro to choose or change an observable
runtime result, Kiro must stop with `ARCHITECTURE_DECISION_REQUIRED` rather than
infer that result from this handoff.

For the registration ledgers below:

- `expectedError` names an already-approved decision/validation outcome. `Deny`
  is an `AuthorizationResult` or domain decision, not a new Problem code.
- `none` means that the registered successful variant has no error. It does not
  suppress an error required by an approved negative variant.
- `existing safe error` means the exact already-approved FEATURE-0012/F18-RD
  error selected by the cited semantic authority; it expressly authorizes no
  new top-level code, violation code, HTTP status, or precedence rule.
- all state, audit, decision, idempotency, and publication effects are exactly
  those closed by the mapped F18-RD and F18-SCN authority.

### 2. Identifier allocation and invariants

1. Register `VS0-CF-F18-01..VS0-CF-F18-37` and
   `VS0-CF-F18-39..VS0-CF-F18-54` as FEATURE-0018-owned local conformance
   identifiers.
2. `VS0-CF-F18-01..37` and `VS0-CF-F18-39..49` map one-to-one to the active
   `AC-F18-*` having the same suffix.
3. `VS0-CF-F18-38` is a permanent tombstone corresponding to excluded
   `AC-F18-38` / retired duplicate F18-SCN-38. It is never active proof and its
   identifier must never be reused.
4. `VS0-CF-F18-50..54` are requirement-only architecture/contract conformance
   cases. They close exact proof for REQ-F18-01, REQ-F18-18, REQ-F18-19,
   REQ-F18-22, and REQ-F18-23, whose obligations are not completely represented
   by one end-user acceptance scenario.
5. Existing shared cases `VS0-CF-F01`, `VS0-CF-F02`, `VS0-CF-X01`, and
   `VS0-CF-X02` remain unchanged and may be cited as supplementary shared proof.
   They do not replace a local `VS0-CF-F18-*` mapping.
6. `VS0-CF-X03` remains FEATURE-0015-owned inherited safe-denial evidence. It is
   never FEATURE-0018-local acceptance proof and must not be deferred to
   FEATURE-0018 design.
7. The already-existing global `VS0-CF-F18` identifier owned by FEATURE-0026 is
   distinct from the new hyphenated FEATURE-0018-local range
   `VS0-CF-F18-NN`; neither identifier changes ownership or meaning.
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

## Source: `.kiro/steering/slice0-contract.md`

Approved SHA-256: `13c9901ef6e242d12080e27539e602ad67859105677a0f0dceb6f6a690764b38`

### Writer, state, error and general conformance anti-drift rules

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
## 6. Writer Anti-Drift

- Every mutable path maps to exactly one `VS0-WRITER-<NNN>` registry entry.
- Clients never write `status`.
- No competing condition producers for the same condition type.
- No plan/decision/execution authority collapse (planner ≠ decision-service ≠ executor).
- No secret values in any writer path; `SecretRef` identifier only.
- Conflicts fail closed per registered conflict code.
- CanonicalMigrationPlan/Record, the migration-controller, and the approved-migration-plan-publisher are retired (DEC-0059 supersedes DEC-0058; ADH-2026-045 canonical bootstrap). VS0-WRITER-020 and VS0-WRITER-021 are permanently retired tombstones and must never be reused or reintroduced as active writers.

## 7. State Anti-Drift

- Exact `VS0-STATE-<NNN>` transitions and guards are authoritative.
- Milestones are observable checkpoints, not lifecycle phases.
- DecisionRecord is always persisted as `FINAL`; no intermediate record states.
- Operation owns async lifecycle: `Pending → Running → Succeeded | Failed | Cancelled`.
- Immutable records are append-only; corrections create linked records, never mutate.
- Invalid transitions are rejected with the registered error/violation.
- VS0-STATE-011 (formerly the CanonicalMigrationRecord milestone sequence) is retired (DEC-0059 supersedes DEC-0058; ADH-2026-045 canonical bootstrap) and must never be reused or reintroduced.

## 8. Error Anti-Drift

- Only exact FEATURE-0012 top-level codes and their HTTP/URN mappings are permitted.
- Slice 0 violation codes appear only in `violations[].code`, never as top-level Problem codes.
- HTTP status codes are not Problem `code` values; they are transport metadata.
- Safe denial: error responses never confirm existence of inaccessible resources.
- No `429` or quota-specific HTTP code; quota exhaustion uses `CONFLICT`/409 + `VS0_QUOTA_EXHAUSTED`.

## 9. Conformance and Traceability

- Every requirement maps to: ADH/DEC authority + owner + schema/writer/state/error IDs + `VS0-CF-<ID>`.
- Failure mappings `VS0-F01..F20` have one-to-one conformance cases `VS0-CF-F01..F20`.
- Migration conformance (`VS0-MIG-F01..F03`, `VS0-CF-MIG01..MIG02`, `VS0-CF-MIGF01..MIGF03`) is retired (DEC-0059 supersedes DEC-0058; ADH-2026-045 canonical bootstrap) and must never be reused or reintroduced.
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

### Prohibited concepts, stop conditions and generated-output contract

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
## 10. Prohibited Active Concepts

Never use in active Slice 0 artifacts:

```text
ResourcePool
ProviderCapability
generic Provider as combined owner/operator
ServiceClass as canonical catalog
EffectivePolicyContext
provider-neutral (as adjective)
six-scope authority
SovrunnInstallation as active Slice 0 resource
```

## 11. Stop Conditions

Emit `ARCHITECTURE_DECISION_REQUIRED` and halt when:

- A semantic choice is missing from loaded authorities.
- Conflicting owners exist for the same path or condition.
- A schema, state transition, or error code is required but unregistered.
- A requirement cannot be mapped to existing registry entries.

Kiro may resolve ordinary mechanics (CRUD wiring, validation plumbing, projection assembly) already explicit in the registry. Kiro must never invent semantics.

## 12. Generated Output Contract

Every stage output (requirements, design, or tasks) must include:

- **Architecture Traceability section**: lists consumed ADH/DEC/RFC, schema IDs, writer IDs, state IDs, error codes.
- **Conformance Mapping section**: maps each acceptance criterion to `VS0-CF-<ID>`.
- No orphan IDs (every referenced ID must exist in registry).
- IDs are never renumbered or reused across stages.

## 13. Implementation Exclusion

This steering governs spec/design/task generation only. No Go code or runtime test implementation derives from this steering alone. Implementation proceeds only from approved feature tasks.
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

## Source: `.kiro/specs/governance-iam-approval-exception-foundation/requirements.md`

Approved SHA-256: `aed8f347ea154c30292cdd0d8f43841fad6d419b8cf789becacc8b57d9c2a9e4`

### Approved normative requirement details and acceptance-scenario coverage

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
| AC-F18-26 | F18-SCN-26 | FEATURE-0017 returns RequiresApproval | Active |
| AC-F18-27 | F18-SCN-27 | FEATURE-0017 returns Indeterminate | Active |
| AC-F18-28 | F18-SCN-28 | Required AuditEvent append fails during JIT activation | Active |
| AC-F18-29 | F18-SCN-29 | End customer performs a routine low-risk operation | Active |
| AC-F18-30 | F18-SCN-30 | User receives denial involving inaccessible target | Active |
| AC-F18-31 | F18-SCN-31 | Production workload needs durable narrow access | Active |
| AC-F18-32 | F18-SCN-32 | Guest assignment omits expiry | Active |
| AC-F18-33 | F18-SCN-33 | Two assignments grant an action but a mandatory guardrail excludes it | Active |
| AC-F18-34 | F18-SCN-34 | A published role version is suspended during use | Active |
| AC-F18-35 | F18-SCN-35 | AuthorizationResult is replayed for a different target | Active |
| AC-F18-36 | F18-SCN-36 | Federated membership provenance is stale or responsible owner is absent | Active |
| AC-F18-37 | F18-SCN-37 | Access review has no or incomplete usage telemetry | Active |
| AC-F18-38 | F18-SCN-38 | Retired duplicate of F18-SCN-29 | Excluded |
| AC-F18-39 | F18-SCN-39 | Access administrator grants ordinary project access | Active |
| AC-F18-40 | F18-SCN-40 | Engineer requests time-bound privileged access | Active |
| AC-F18-41 | F18-SCN-41 | Approver decides a privileged or exception request | Active |
| AC-F18-42 | F18-SCN-42 | Reviewer certifies or remediates access | Active |
| AC-F18-43 | F18-SCN-43 | User requests a bounded control exception | Active |
| AC-F18-44 | F18-SCN-44 | Resource-narrowed administrator tries to delegate another resource or the whole scope | Active |
| AC-F18-45 | F18-SCN-45 | Standing beneficiary attempts to certify retained access | Active |
| AC-F18-46 | F18-SCN-46 | Temporary membership administrator attempts to create durable group-derived access | Active |
| AC-F18-47 | F18-SCN-47 | Membership administrator attempts to manufacture approval or privileged eligibility | Active |
| AC-F18-48 | F18-SCN-48 | Reviewer eligibility is missing, indirect, stale or scope-incompatible | Active |
| AC-F18-49 | F18-SCN-49 | Role grantor or later rule publisher attempts to manufacture eligibility | Active |

## 4. Normative requirements and acceptance scenarios

### 4.1 Normative requirement details

Each approved REQ has exactly one normative detail heading below, in approved
order. Each heading states observable behavior only and cites its sole
controlling F18-RD authority. Where a statement depends on an inherited
dependency contract, it references that contract without redefining it.

#### REQ-F18-01 — Current-feature-only authority (F18-RD-01)

- FEATURE-0018 defines only the contracts and observable behavior required for its
  own Phase 2R development; it may name implementation-neutral future adapter
  boundaries but selects, designs, or implements no future adapter or consumer.
- The exclusions in F18-RD-01 are closed and authorize no future resource,
  adapter, vendor, protocol, credential flow, deployment boundary, requirement,
  task, or implementation. These prohibited concepts are enumerated only in
  section 7; active requirements reference that exclusion set without repeating it.

#### REQ-F18-02 — Closed contract inventory and profiles (F18-RD-02)

- The persistent resource inventory is exactly the ten contracts in section 3.1;
  only rows with persistent resource profiles are independently addressable. The
  supporting values in section 3.2 have no endpoint, independent lifecycle,
  storage authority, or user journey.
- Each contract has exactly one accepted-intent/immutable-spec writer and exactly
  one status/result/record writer as registered by F18-RD-02. The RoleAssignment
  controller is the sole canonical RoleAssignment publisher and sole writer of its
  specification, lifecycle, status, revocation, replacement, expiry, and
  review-due effects; all other authorities submit only exact immutable intents.
- Scope applicability is closed and deterministic per versioned contract and is
  consumed generically. Allowed `metadata.scopeRef` kinds per contract are exactly
  the registered set: AccessGroup and Membership at Organization/CloudProvider;
  RoleDefinition, ApprovalPolicy, and GovernanceProfile at
  Platform/Organization/CloudProvider; RoleAssignment, PrivilegedAccessRequest,
  AccessReview, ApprovalRequest, and ExceptionGrant at all seven canonical
  ScopeKinds. Customer policy may narrow but never expand a registered contract.
- Definition publication scope and reference compatibility are separate from
  authorization containment and never grant an action: Platform definitions are
  usable at any registered consuming scope; Organization definitions only within
  their Organization tree; CloudProvider definitions only at that exact provider;
  a Platform GovernanceProfile may reference only a Platform policy. CloudPlatform
  RoleDefinition/ApprovalPolicy/GovernanceProfile publication is not part of the
  initial contract. Undefined or cross-tree reference compatibility fails closed.
- No FEATURE-0018 contract supports generic `PUT`, unrestricted `PATCH`, hard
  deletion, or caller-authored status. Schedulers may invoke an owning controller
  but never become field writers. Adding a value, ScopeKind, containment, or
  applicability combination requires versioned contract change and conformance
  evidence per F18-RD-02; it is not a requirements or design choice.

#### REQ-F18-03 — Stable principal identity (F18-RD-03)

- Durable external identity is `issuer + subject`. Email, display name, username,
  persona, job title, and external role label are non-authorizing attributes.
- The closed principal categories are Human, Workload, and System. FEATURE-0018
  consumes an already-authenticated `PrincipalRef` and canonical
  `AssuranceEvidence`; it validates no token and contacts no identity provider.
- `AssuranceEvidence` is an operation-local value binding the exact PrincipalRef,
  authoritative `authenticatedAt`, closed level `AAL1 | AAL2 | AAL3`, boolean
  `phishingResistant`, trusted `sourceAuthority`, opaque integrity reference, and
  optional evidence expiry. Its PrincipalRef must equal the authenticated actor;
  age is computed only from authoritative UTC time. Missing, expired, unrecognized,
  mismatched, or untrusted evidence fails closed. Raw tokens, assertions,
  authentication secrets, and provider-native method payloads are never retained.

#### REQ-F18-04 — Direct principal and scoped AccessGroup assignments (F18-RD-04)

- `RoleAssignment.roleHolderRef` is exactly one `PrincipalRef` or `AccessGroupRef`.
- An `AccessGroup` is scoped to exactly one Organization or CloudProvider, has one
  responsible Human PrincipalRef owner, cannot authenticate, contains only direct
  `PrincipalRef` membership in v1, and rejects nested and dynamic groups. It never
  derives access from an external group claim. It becomes authorization-ineligible
  when suspended and terminally ineligible when retired (retired prevents new
  membership/assignment and retains immutable history). Lifecycle is
  `Active -> Suspended -> Active` and `Active | Suspended -> Retired`; creation
  publishes Active only after authorization and required audit-obligation acceptance.
- The responsible owner must resolve to a current Human at group creation and every
  privilege-increasing operation. An inactive/unresolved owner blocks adding or
  reactivating members, expanding the group, and creating or expanding group-held
  assignments, but is not a runtime dependency for already-effective access.
  Privilege-reducing transitions and authorized owner recovery remain permitted;
  recovery cannot derive authority from the inactive owner.
- Creating/reactivating a direct AccessGroup Membership, or restoring/extending its
  assignment-effect interval through provenance or freshness, is a
  privilege-increasing operation whenever it makes a group-held RoleAssignment newly
  or longer effective, and must pass the membership-enabled grant ceiling (REQ-F18-09)
  at the Membership controller's atomic publication boundary. No ownership,
  administrator, trusted-provisioner, or controller identity bypasses that ceiling.
- `AccessGroup` is a relationship target, never a ScopeKind. A group Membership
  retains its Organization/CloudProvider `metadata.scopeRef` and one exact
  UID-pinned `accessGroupRef`; the Membership and group must share the same
  canonical owning scope; cross-scope group membership is rejected.
- Role-holder reference compatibility is exact: an Organization-scoped AccessGroup
  may hold a RoleAssignment at that Organization or its OrganizationUnit/Tenant/Project
  descendants; a CloudProvider-scoped AccessGroup only at that exact CloudProvider.
  AccessGroup-held Platform and CloudPlatform assignments are prohibited in v1.
  AccessGroup and group-derived assignments never supply approver, requester, or
  reviewer eligibility in v1. Compatibility never grants an action; policy may
  narrow but never broaden it.

#### REQ-F18-05 — Membership is non-authorizing (F18-RD-05)

- A `Membership` associates one principal with an Organization or CloudProvider
  belonging context, or records one direct AccessGroup relationship within that
  same scope. It never grants an action; never establishes approver, requester,
  reviewer, JIT, or break-glass eligibility; never acts as a role; and never stores
  `roleAssignmentRefs`. A closed `MembershipContextKind` distinguishes the three
  context kinds; only AccessGroup context carries `accessGroupRef`.
- Membership retains immutable `sourceAuthority`, optional opaque `sourceObjectRef`,
  `provisionedAt`, and, when synchronized, `lastSynchronizedAt` and system-owned
  `freshUntil` (a synchronized membership is authorization-ineligible at
  `freshUntil`). `membershipType` is `Standard | Guest`; Guest is Human-only and
  requires a Human sponsor and `expiresAt`; Workload and System must be Standard and
  require `responsiblePartyRef`. Raw tokens, assertions, credentials, claims, and
  provider-native group schemas are prohibited.
- `MembershipRelationshipKey` = principalRef + contextKind + metadata.scopeRef +
  accessGroupRef (when AccessGroup). At most one non-terminal Membership exists per
  key; creating another record cannot bypass suspension or staleness. Retained
  lifecycle is `Active -> Suspended -> Active` and `Active | Suspended -> Revoked`.
  Guest expiry and synchronized freshness expiry are authoritative-time eligibility
  projections, not mutable transitions. `GuestExpired` is terminal for that UID;
  `Stale` is recoverable only through an authorized update to system-owned
  freshness/provenance fields.
- For a Workload/System principal, the responsible party must resolve to a current
  Human when a Membership is created and when an AccessReview retains it; this is
  retained accountability evidence, not a per-request authorization condition. Later
  responsible-party inactivity blocks new creation/retention and triggers protected
  review/escalation but does not by itself make an otherwise-current production grant
  ineffective; explicit policy suspension/revocation may.
- Membership replacement and renewal are not v1 operations; synchronization updates
  only system-owned freshness/provenance fields. Contexts are independently
  evaluated; multiple memberships never union identity, scope, or authority.
  Suspension/revocation disables only authorization depending on that exact
  membership context. Where an assignment requires Membership, it pins the exact
  authorizing Membership UID; a new/later Membership never revives an assignment
  bound to a revoked Membership; guest-assignment `expiresAt` cannot exceed the
  pinned guest Membership `expiresAt`. Membership remains non-authorizing even though
  an AccessGroup assignment-effect expansion is grant-producing under REQ-F18-09.

#### REQ-F18-06 — Distributed action ownership and central role composition (F18-RD-06)

- The feature owning an operation owns its canonical action meaning. FEATURE-0018
  validates registered actions and composes them into roles without renaming,
  splitting, or reinterpreting them. FEATURE-0016's `executiontarget.read`,
  `executiontarget.qualify`, and `executiontarget.write` are consumed unchanged; the
  combined create/retire granularity of `executiontarget.write` remains an explicit
  recorded dependency limitation. A missing or ambiguous required intrinsic
  classification is fail-safe Privileged (REQ-F18-07).
- The exact initial assignable action registry and its intrinsic Ordinary/Privileged
  classifications are as registered by F18-RD-06 and are consumed generically. Each
  canonical action registers one or more exact discriminated `ActionTargetBinding`
  variants — `ExactResource { action, targetKind }`, `CreateParent { action,
  parentScopeKind }`, or `ScopeOnly { action, scopeKind }` — with exactly one variant
  applying per registration. A mixed variant, missing discriminator or required kind,
  empty, wildcard, unknown, or ambiguous registration fails closed. An absent
  FEATURE-0016 registration makes that action unassignable and unusable in
  FEATURE-0018 rather than inferred or reinterpreted.
- FEATURE-0018 authorizes canonical Sovrunn control-plane actions only. Provider-native
  principals, groups, roles, policies, permissions, and credentials are neither
  FEATURE-0018 resources nor outputs, and are never created, published, synchronized,
  translated, or exposed. A native ExecutionTarget effect requires both valid Sovrunn
  authorization and independent native IAM authorization of a least-privileged
  adapter/controller identity: the two layers compose by intersection; neither Allow
  substitutes for the other and neither overrides the other's Deny. Direct out-of-band
  IaaS administration remains CloudProvider-owned and is never represented as a Sovrunn
  RoleAssignment effect.

#### REQ-F18-07 — Versioned RoleDefinition (F18-RD-07)

- A RoleDefinition contains only registered canonical Sovrunn actions; provider-native
  IAM permissions and unrestricted wildcards are prohibited. Lifecycle is
  `Draft -> Published`, `Published -> Suspended`, `Suspended -> Published`, and
  `Published | Suspended -> Retired`. Draft is mutable; published versions are
  immutable; changes create a new version; version lineage is linear in v1.
- Supersession is an informational version-lineage relationship written only by the
  publication controller that publishes the successor; it atomically records exact
  `supersededByRef` on the previously current version and emits the registered
  supersession AuditEvent. Superseded versions receive no new references while valid
  existing exact pins remain effective.
- Emergency suspension makes every assignment pinned to that exact version ineffective
  immediately without mutating the definition or assignment and is reversible by
  restoration. Retirement is terminal: a retired version receives no new references and
  every assignment pinned to it is immediately and permanently
  authorization-ineligible; retirement is not a migration mechanism. Suspension,
  restoration, and retirement require current authority, justification, mandatory
  `authorization-decision/v1` DecisionRecord evidence, and protected audit-obligation
  acceptance before lifecycle publication.
- Each published version declares an immutable base classification `Ordinary` or
  `Privileged`; its effective classification is `Privileged` when the definition
  declares it, any registered action requires privileged handling, an applicable
  system-selected `privilegedAccessRule` elevates it, or any required classification
  input is missing or ambiguous. Classification sources may elevate but never
  downgrade. A privileged role cannot receive a Standing or AccessGroup assignment and
  uses an approved direct-PrincipalRef TimeBound flow; the production Workload/System
  Standing exception (REQ-F18-08) applies only to effectively Ordinary roles.

#### REQ-F18-08 — Scoped RoleAssignment (F18-RD-08)

- One RoleAssignment binds one `roleHolderRef` (PrincipalRef or AccessGroupRef), one
  published RoleDefinition identity and exact version, `metadata.scopeRef` as sole
  canonical target scope, optional exact `resourceRef` governed by that scope, the exact
  authorizing `membershipRef` when Membership is required for a direct PrincipalRef
  holder, exact `responsiblePartyRef` for a direct Workload/System holder, and explicit
  validity mode with mode-specific timestamps. `resourceRef` narrows to one exact
  resource and never becomes a second scope authority; it contributes only to an
  `ExactResource` authorization with matching UID and a target-binding-permitted kind.
  Assignment creation rejects when no action of the pinned version can apply to the
  resource kind.
- Holder, exact role version, scope, resource restriction, required Membership UID,
  responsible party, pinned Standing review rule, and validity are immutable; change
  occurs only through replacement or revocation. `membershipRef` is prohibited for an
  AccessGroup holder and for Platform/CloudPlatform authorization. Responsible-party
  evidence supplies accountability, not permission or a per-request condition.
- For a direct Workload/System holder, the system-selected `responsiblePartyRef` must
  resolve to a current Human at assignment creation, at any replacement that materializes
  a successor assignment, and at every AccessReview retention of that assignment. Where a
  Membership applies to the assignment, that `responsiblePartyRef` must equal the
  responsible party pinned on the exact authorizing Workload/System Membership. Missing,
  mismatched, or non-Human responsibility evidence fails closed: the creation,
  replacement, or retention is rejected and no RoleAssignment intent is published.
- The RoleAssignment controller is the sole publisher and sole writer of specification,
  lifecycle, status, revocation, replacement, expiry, and review-due effects; it
  revalidates all pins and required decision/audit evidence before publishing and cannot
  enlarge the submitted intent.
- `Standing` has no automatic expiry, requires periodic review, and is immediately
  revocable; `TimeBound` requires both `notBefore` and `expiresAt` and becomes
  ineffective at exact expiry without grace. Privileged, guest, JIT, and break-glass
  assignments are TimeBound. A Standing assignment is valid only through the exact
  registered access-review rule path that permits the holder kind, scope kind, and (when
  resource-narrowed) target kind, with an effectively Ordinary role and periodic review;
  production Workload/System assignments may be Standing only through this path.
- Every Standing assignment pins one exact `accessReviewRuleRef` and has system-owned
  `status.nextReviewDueAt`; only a completed exact-version `StandingCertification` with
  `Retain` advances the due time. For Human or AccessGroup-held Standing access, reaching
  the due time without current certification makes it authorization-ineligible until
  certified or replaced. Retained lifecycle is `Current -> Revoked | Expired` (Expired is
  automatic only for TimeBound at exact `expiresAt`); effective projection is exactly
  `NotYetValid | Effective | InactiveDependency | Revoked | Expired`. There is no
  user-authored condition language, deny assignment, wildcard condition, arbitrary
  expression, or generic ABAC.

#### REQ-F18-09 — Deterministic scoped authorization composition (F18-RD-09)

- Authorization is explicit grant union constrained by guardrail intersection.
  `Allow` requires all of: current authenticated PrincipalRef; every membership required
  by the requested scope and holder is current; requested action is in the union of
  actions from every independently applicable RoleAssignment; exact target and canonical
  scope match; every contributing `resourceRef` matches the exact target; every
  contributing assignment is valid; every mandatory guardrail permits; no policy outcome
  is Deny or Indeterminate; and any operation-local approval/exception evidence required
  now is current and exact. Restrictions from different assignments never synthesize a
  grant. Policy Allow, approval, exception, Membership, controller identity, native cloud
  IAM, and retained results never create Sovrunn authority.
- Membership eligibility is contextual: Organization/OrganizationUnit/Tenant/Project
  authorization requires current belonging in the governing Organization; CloudProvider
  authorization requires current belonging in that exact CloudProvider; AccessGroup-derived
  authorization also requires current direct Membership in the exact same-scope,
  authorization-eligible group; Platform and CloudPlatform authorization use current
  PrincipalRef and RoleAssignment authority without inventing a Membership. The only
  inherited containment chain is Organization -> OrganizationUnit -> Tenant -> Project;
  undefined containment is exact-scope only.
- Delegated administration has a mandatory grantor ceiling: the proposed grant has one
  exact reach from its target-binding mode, scope, and optional exact resource; a grantor
  must hold current `roleassignment.grant` over the complete reach and, per proposed
  action, one independently applicable delegable witness covering the same or broader
  permitted reach, plus `roleassignment.delegate` when the role contains
  `roleassignment.grant`/`roleassignment.delegate`. `ExactResource(A)` covers only
  resource A; scope-wide authority may cover an exact resource governed by that scope.
  Witness fragments from different assignments cannot be synthesized. Grantor authority is
  evaluated at publication and does not cascade to already-published assignments.
- An AccessGroup Membership assignment-effect expansion is a separate grant-producing
  publication boundary. The Membership controller derives one operation-local
  `MembershipEnabledGrantEnvelope` from a coherent snapshot containing every independently
  applicable group-held assignment newly or longer effective, with only the newly enabled
  or extended effect interval (intersection of assignment validity, Guest/`freshUntil`
  bounds, group eligibility, and every dependency). Publishing a non-empty envelope requires
  the applicable Membership-operation authority, current `roleassignment.grant` over the
  complete envelope, one delegable witness per envelope action/reach/validity tuple, and
  accepted concurrency/decision/audit obligations. Standing effects require Standing
  witnesses; a TimeBound administrator/provisioner cannot exceed its temporal ceiling; an
  empty envelope grants no action and establishes no approver/reviewer/JIT/break-glass
  eligibility. Trusted provisioners receive no bypass; concurrent changes conflict or retry
  against a coherent snapshot.
- `AuthorizationInput` and `AuthorizationResult` are operation-local, non-bearer values, not
  resources, tokens, sessions, or retained evidence; a result cannot authorize a different
  request. FEATURE-0017 `Indeterminate` is successfully mapped to `Deny` with stable
  protected provenance; a validation/evaluation/mandatory audit-obligation failure returns a
  safe error and no AuthorizationResult.

#### REQ-F18-10 — FEATURE-0017 adoption without reinterpretation (F18-RD-10)

- FEATURE-0018 consumes FEATURE-0017 outcomes exactly: `Allow` continues remaining checks;
  `Deny` denies; `RequiresApproval` consumes one exact system-selected ApprovalPolicy
  reference and evaluates it and never auto-creates a request; `Indeterminate` fails closed.
  Transport success is not authorization; the FEATURE-0017 fake stays digest-keyed and
  contains no IAM/approval/exception/target/governance logic.
- Every grant-producing operation derives exactly one system-owned operation-local
  `ApprovalRequirement`: `NotRequired`, or `Required` with one exact system-selected
  published ApprovalPolicy version ref. The value is not requester-authored, independently
  stored, or bearer authority. Missing/ambiguous requirement, or `Required` with zero,
  multiple, unpublished, mismatched, or unavailable policy, fails closed and creates no
  ApprovalRequest; `NotRequired` carries no policy ref and an unexpected ref rejects.
- The originating operation controller is the sole derivation authority before intent
  acceptance (grant controller for RoleAssignment grant/proposal; privileged-access
  controller for PrivilegedAccessRequest; exception controller for ExceptionProposal;
  Membership controller for an AccessGroup assignment-effect expansion). JIT
  PrivilegedAccessRequest, effectively Privileged RoleAssignmentProposal, and
  ExceptionProposal derive `Required`; an effectively Ordinary role grant may derive
  `NotRequired` while still requiring `roleassignment.grant`; an AccessGroup Membership
  assignment-effect expansion always derives `NotRequired` in v1 and never weakens the
  membership-enabled grant ceiling or creates eligibility; BreakGlass derives `Required`
  unless the exact rule permits bypass (REQ-F18-15). A requester cannot select or replace the
  requirement or policy; FEATURE-0018 performs no GovernanceProfile resolution and consumes
  trusted preselected evidence.

#### REQ-F18-11 — Bounded FEATURE-0017 subject/target use (F18-RD-11)

- FEATURE-0018 may invoke FEATURE-0017 only for an exact UID- and generation-pinned
  PrivilegedAccessRequest at Project, CloudPlatform, or CloudProvider scope using canonical
  action `privilegedaccessrequest.submit`. It does not invoke FEATURE-0017 for ordinary
  authorization, RoleAssignmentProposal, ExceptionProposal, AccessReview, Membership,
  AccessGroup, or group evaluation. An unsupported subject or scope never causes contract
  emulation; when a mandatory FEATURE-0017 evaluation cannot be performed, the operation fails
  closed.
- The AccessGroup assignment-effect expansion and its grantor ceiling are deterministic
  FEATURE-0018-local authorization algebra, not policy emulation or a new FEATURE-0017
  subject. Ordinary RBAC algebra remains FEATURE-0018-owned and is not a custom policy engine;
  target-aware dynamic policy for ordinary authorization requires a separately approved
  FEATURE-0017 version and must not overload `subjectRef` or emulate the missing contract.

#### REQ-F18-12 — Bounded ApprovalPolicy (F18-RD-12)

- ApprovalPolicy is declarative with one to three linear ordered stages, eligible approvers,
  required approval count per stage, decision expiry, justification rules, self-approval
  restrictions, and separation of duties. It is evaluated only for `ApprovalRequirement.Required`
  and never authorizes break-glass bypass (bypass is defined solely by REQ-F18-15). Its closed
  applicability key is `subjectKind + canonical action + ScopeKind + optional targetKind +
  requestMode`; `subjectKind` is exactly `PrivilegedAccessRequest`, `RoleAssignmentProposal`, or
  `ExceptionProposal`; `requestMode` when applicable is exactly `JIT` or `BreakGlass` (permitted
  only for PrivilegedAccessRequest). Applicability names no exact target instance and performs no
  hierarchy/profile resolution.
- The eligible-approver list is non-empty and every entry is one `EligibilityRef` containing
  exactly one Human PrincipalRef, satisfied only by exact principal equality. RoleDefinition,
  RoleAssignment, AccessGroupRef, Membership, external/nested/dynamic group never satisfy or
  create approver eligibility; the list uses OR semantics; eligibility never broadens scope or
  target and never grants `approvalrequest.decide` (that action is independently required).
  Empty, duplicate, malformed, unsupported, inactive, stale, ambiguous, scope-incompatible, or
  target-incompatible evidence fails closed.
- Approver eligibility and separation of duties are re-evaluated when a decision is accepted and
  immediately before downstream effect publication. One principal contributes at most one
  immutable effective decision to a stage; same-decision replay is idempotent and a different
  decision conflicts; decisions for an inactive/completed stage reject. Separation of duties uses
  an exact current conflict set (requester; every direct beneficiary; a beneficiary AccessGroup's
  owner and current direct principal members; the responsible Human for a beneficiary
  Workload/System principal; and every beneficiary resolved through an exact
  RoleAssignmentRef/PrivilegedAccessRequestRef); a conflicted principal cannot approve, and a
  conflict discovered before publication blocks the effect and requires a fresh request. Stages
  execute strictly in order; any valid `Deny` terminates the request as `Denied`; a stage
  succeeds only at its positive-integer quorum met by distinct eligible principals. A Human
  principal whose approval counts toward one stage cannot count toward another stage of the
  same ApprovalRequest. Decision
  expiry is greater than zero and no more than twenty-four hours (narrower policy may shorten
  only). No abstain, branching, loop, delegation chain, scripted escalation, expression, or
  general workflow engine exists. Lifecycle is `Draft -> Published -> Retired` with immutable
  published versions and informational supersession.

#### REQ-F18-13 — Immutable terminal approval evidence (F18-RD-13)

- ApprovalRequest lifecycle is `Pending -> Approved | Denied | Cancelled | Expired`; terminal
  outcomes are immutable. The subject is exactly one of a `PrivilegedAccessRequestRef`, embedded
  immutable `RoleAssignmentProposal`, or embedded immutable `ExceptionProposal`, each pinning its
  exact required fields; an effectively Privileged role proposal must be direct-PrincipalRef and
  TimeBound (AccessGroup or Standing privileged proposals reject). Arbitrary or future-domain
  payloads reject.
- A request is created only for `ApprovalRequirement.Required` and only from an explicit,
  successfully accepted FEATURE-0018 business intent; a `NotRequired` ordinary grant creates no
  request. The Approval-request controller is the sole creator and pins the exact subject, one
  system-selected published ApprovalPolicy version, scope, requester, correlation, and one
  immutable `expiresAt` equal to the earliest of creation time plus policy decision-expiry and the
  applicable subject bound (a PrivilegedAccessRequest `activationDeadline`, a TimeBound proposal's
  proposed `expiresAt`, or an ExceptionProposal's requested `expiresAt`; a Standing proposal
  contributes none). No replay or stage transition extends it. An Approved request stays
  historically Approved but its evidence is current only while before `expiresAt`.
- Approver eligibility and separation of duties are rechecked at decision acceptance and downstream
  effect publication (REQ-F18-12). The first valid terminal transition wins under inherited
  FEATURE-0012 optimistic concurrency. One originating intent creates at most one ApprovalRequest;
  idempotency binds originating subject identity, actor, scope, and request digest. Approval is
  evidence and never performs the effect: after approval the owning controller re-evaluates current
  authorization and exact evidence, and only the RoleAssignment controller may publish either
  RoleAssignment intent; denial, cancellation, or expiry authorizes no effect.

#### REQ-F18-14 — JIT privileged access authorizes one temporary grant (F18-RD-14)

- The requester and resulting holder are the same exact Human PrincipalRef; Workload/System
  time-bounded administration uses an authorized RoleAssignmentProposal with responsiblePartyRef, not
  interactive JIT or break-glass. Every PrivilegedAccessRequest pins system-owned immutable
  `submittedAt`, `mode` (`JIT | BreakGlass`), `activationDeadline`, and requested duration. The atomic
  publication instant of both the linked TimeBound RoleAssignment and the request's `Active` transition
  must be strictly before `activationDeadline`. For JIT the deadline is no later than
  `submittedAt +` the selected policy decision-expiry (bounded by the twenty-four-hour maximum); for
  BreakGlass no later than fifteen minutes after `submittedAt`. Requested duration begins only at
  successful activation.
- The exact system-selected privileged-access rule carries a non-empty `requesterEligibilityRefs` list
  of `EligibilityRef` values (exact Human PrincipalRefs, exact-requester equality, OR semantics);
  roles, assignments, groups, Membership, and external/nested/dynamic groups never satisfy or create
  it. Eligibility is re-evaluated at admission and atomic activation, is necessary but never grants
  `privilegedaccessrequest.submit` or the requested privileged actions, and fails closed on empty,
  duplicate, malformed, unsupported, inactive, stale, ambiguous, or scope/target-incompatible evidence.
- Lifecycle is `Submitted -> PendingApproval | Active | Cancelled | Expired`, `PendingApproval ->
  Approved | Denied | Cancelled | Expired`, `Approved -> Active | Cancelled | Expired`, and
  `Active -> Expired | Revoked` (Denied/Cancelled/Expired/Revoked terminal). `Approved` is non-terminal
  evidence; failed fresh activation checks keep the request `Approved` with readiness `Blocked`.
  Approval-path activation requires every contextually required Membership, current
  `privilegedaccessrequest.submit` authority and exact eligibility, one exact current ApprovalPolicy,
  an approved request, current phishing-resistant AAL2-or-higher AssuranceEvidence no more than fifteen
  minutes old, bounded duration, and no conflicting/expired evidence; a rule may require AAL3 or shorter
  age but cannot weaken the floor. Activation performs a fresh authorization evaluation and may submit
  exactly one individual TimeBound RoleAssignment materialization intent linked immutably to the request
  and to approval/bypass evidence; only the RoleAssignment controller publishes it, and the request
  becomes Active only when that publication succeeds. Normal JIT defaults to one hour with an eight-hour
  platform maximum (policy may shorten only); the published assignment `expiresAt` is the earliest of
  `notBefore +` approved duration, pinned Membership expiry, and applicable rule/policy ceilings. Delayed
  activation never extends any bound; activation after the deadline requires a new request.

#### REQ-F18-15 — Constrained break-glass access (F18-RD-15)

- Break-glass is a PrivilegedAccessRequest mode, not a resource. It is individual only and requires
  pre-authorized eligibility, fresh phishing-resistant AAL2-or-higher AssuranceEvidence, explicit
  emergency justification, immediate protected audit, exact scope, automatic expiry, and retrospective
  AccessReview. It requires one exact system-selected privileged-access rule; the rule's
  `breakGlassAllowed` boolean is the sole authority permitting bypass of normal pre-approval, and
  break-glass never bypasses authentication, scope isolation, expiry, or audit.
- When the exact rule has `breakGlassAllowed=true` and every precondition passes, the activation trigger
  derives `ApprovalRequirement.NotRequired`, may authorize submission of exactly one linked TimeBound
  RoleAssignment materialization intent, and creates no ApprovalRequest; only the RoleAssignment
  controller publishes it, and the exceptional `Submitted -> Active` transition occurs only when the
  assignment publication and linked-review acceptance succeed. When `breakGlassAllowed=false`, an
  otherwise-valid request derives `Required` and follows the normal `Submitted -> PendingApproval` path.
  Missing, multiple, unpublished, mismatched, or unavailable rule evidence fails closed and never becomes
  an approval fallback; authentication, structural, eligibility, scope, authorization, review-rule, or
  audit-obligation failure never falls back to approval or is repaired by an approver.
- An accepted bypass-authorized request whose current assurance, authorization, eligibility,
  duration, required audit-obligation acceptance, or other activation evidence later becomes
  stale, unavailable, failed, or conflicting remains non-active in `Submitted` with system-owned
  activation readiness `Blocked`. It creates no ApprovalRequest and authorizes no publishable
  RoleAssignment intent. Correction may trigger an idempotent retry strictly before
  `activationDeadline` without extending any validity bound.
- Maximum duration is one hour, assurance age at most fifteen minutes, mutation renewal prohibited, and
  review due within twenty-four hours (a rule may only shorten these); audit failure blocks activation.
  Successful direct activation atomically materializes exactly one linked Pending AccessReview with
  immutable reviewed-subject, activation/grant references, system-selected `accessReviewRuleRef`,
  originating `privilegedAccessRule` evidence, and the exact `dueAt` from REQ-F18-16; failure to accept the
  review or its audit obligation blocks activation; retry is idempotent and creates no second review. An
  overdue review never extends the independently expiring access.

#### REQ-F18-16 — Snapshot-based AccessReview (F18-RD-16)

- Campaign mode is exactly `StandingCertification`, `BreakGlassRetrospective`, or `Manual`. The first two
  each review exactly one linked RoleAssignment at that assignment's exact scope; a Manual campaign carries
  a non-empty, deduplicated list of exact UID- and resourceVersion-pinned Membership/RoleAssignment
  references with no dynamic query and honors exact containment per scope (cross-tree and undefined
  containment reject). Every campaign pins exactly one system-selected published `accessReviewRuleRef`;
  `dueAt` derives per mode (`StandingCertification` = assignment `nextReviewDueAt`; `Manual` = `createdAt +
  reviewWindow`; `BreakGlassRetrospective` = `activationAt + min(24h, retrospectiveReviewDeadline)`).
- The rule's non-empty `reviewerEligibilityRefs` list contains only `EligibilityRef` exact Human
  PrincipalRefs (exact-reviewer equality, OR semantics); roles, assignments, groups, Membership, and
  external/nested/dynamic groups never satisfy or create reviewer eligibility. Eligibility is necessary but
  never grants `accessreview.decide`, is re-evaluated at decision acceptance and immediately before atomic
  remediation publication, and fails closed on empty/duplicate/malformed/unsupported/inactive/stale/
  scope-incompatible/target-incompatible evidence.
- Lifecycle is `Pending -> InProgress | Cancelled | ExpiredIncomplete` and `InProgress -> Completed |
  Cancelled | ExpiredIncomplete` (Completed/Cancelled/ExpiredIncomplete immutable terminal). The immutable
  snapshot is captured exactly once on entry to InProgress; an empty, changed, or unresolvable required
  snapshot fails closed. Per-item dispositions are closed by kind (Membership: `Retain | Revoke | Stale`;
  RoleAssignment: `Retain | Revoke | Replace | Stale`). Revoke/Replace authorize the review controller to
  submit exact intents; only the RoleAssignment controller publishes revocation or atomically publishes the
  new assignment and revokes the exact old one. Replace must preserve the holder and may only reduce
  authority; a changed item yields `Stale` with no remediation.
- No principal may retain, replace, or advance review eligibility for authority from which it currently
  benefits. For every access-preserving/replacing disposition the controller derives the exact
  `ReviewBeneficiaryHumanSet` (Human holder -> that Human; Workload/System holder -> its current responsible
  Human; AccessGroup holder -> current owner + every current direct Human member + the responsible Human of
  every current direct Workload/System member), bounded to current direct memberships; a replacement uses the
  union across reviewed and proposed assignments; the acting reviewer must be a current Human outside that
  union. Unresolved evidence fails closed with `REVIEW_BENEFICIARY_UNRESOLVED` (no effect, no due advancement);
  a conflict returns `REVIEWER_CONFLICT` (no effect, no `nextReviewDueAt`, item left for a different eligible
  reviewer); independently authorized Revoke remains available. For a Membership-only
  `Retain` disposition the controller derives the equivalent beneficiary-Human expansion
  (Human principal -> that Human; Workload/System principal -> its current responsible
  Human; AccessGroup relationship -> current owner + every current direct Human member +
  the responsible Human of every current direct Workload/System member), bounded to
  current direct memberships, and the acting reviewer must be a current Human outside that
  set. Retaining a Workload/System Membership or a Workload/System-held RoleAssignment
  additionally requires that its `responsiblePartyRef` resolve to a current Human; failed
  responsibility or beneficiary validation must not certify the item and must not silently
  revoke it (the item is left undecided for a different eligible reviewer or an
  independently authorized Revoke). A completed exact-version `StandingCertification
  + Retain` advances `nextReviewDueAt` to `completedAt + reviewInterval`; Manual and BreakGlassRetrospective
  never advance it. Missing usage is not proof of non-use unless coverage is complete; incomplete review never
  certifies. Certification and due advancement publish under the same local audit/concurrency boundary.
- Before snapshot acceptance, unresolved required evidence leaves the campaign `Pending` with
  readiness `Blocked` until correction, cancellation, or `dueAt`; at `dueAt` it becomes immutable
  `ExpiredIncomplete`. A campaign completes only when every snapshotted item has one final
  disposition; `Stale` completes an item without certification or remediation, and any undecided
  item at `dueAt` makes the campaign `ExpiredIncomplete`. Human and AccessGroup-held Standing
  assignments become authorization-ineligible at their review due time until a completed
  exact-version `StandingCertification + Retain` or authorized replacement exists. An incomplete
  production Workload/System review escalates and follows its pinned `overdueEffect`; it is not
  silently revoked solely because telemetry is missing or incomplete unless a prepublished rule
  requires suspension. Failure to materialize a due campaign never extends assignment validity or
  effectiveness.

#### REQ-F18-17 — Bounded immutable exception evidence (F18-RD-17)

- No `ExceptionRequest` resource exists. The flow is an immutable `ExceptionProposal` inside an ApprovalRequest
  followed by an approved `ExceptionGrant` pinning exact control, subject, scope, bounded override, interval,
  justification, compensating controls, approval/decision evidence, and linkage. A current Approved request plus
  successful final validation produces `Grant`; ApprovalRequest Denied/Cancelled/Expired produces `Deny` with the
  corresponding reason; final proposal validation may produce terminal `Deny` with exactly one stable reason from
  the registered set; `ConcurrencyConflict | DependencyUnavailable | AuditObligationUnavailable` are retryable
  failures using the same immutable proposal identity, becoming terminal `Deny` with `ApprovalExpired` if
  unresolved at `expiresAt`.
- `subjectRef` is exactly one PrincipalRef, RoleAssignmentRef, PrivilegedAccessRequestRef, or UID-pinned canonical
  ResourceRef of a control-registered kind. The versioned control definition declares its identity/version,
  exceptionability, allowed subject kinds, typed bounded-override schema, required compensating-control kinds, and
  maximum duration; `boundedOverride` must validate against that exact schema. The platform ceiling is ninety days
  (policy may shorten). Authentication integrity, tenant/scope isolation, audit integrity, immutable evidence, and
  unresolved conflict are non-exceptionable. Final constraint composition is exact (minimum durations, union of
  compensating controls, intersections of subjects/scopes/overrides); the controller rejects a nonconforming
  proposal and never silently rewrites it.
- Subject/scope compatibility is exact; an Organization-scoped exception does not cover descendants; all listed
  scopes are exact-scope only. ExceptionGrant is append-only using the half-open interval `[notBefore, expiresAt)`;
  two grants overlap when the same exact control identity/version, UID-pinned subject, and scope have intersecting
  effective intervals, and a new overlap is rejected with no merge. Early revocation is linked same-kind immutable
  evidence with `effect=Revoke`, exact `revokedGrantRef`, and authoritative `revokedAt`, shortening the effective
  interval without mutation. An exception never grants a role; FEATURE-0020 alone owns effective application and
  cross-profile resolution. ExceptionGrant has no mutable lifecycle/status; effectiveness is computed from time, the
  immutable interval, and valid linked revocation; expiry emits the registered audit and may refresh a
  non-authoritative projection without mutating the grant.

#### REQ-F18-18 — FEATURE-0018-limited GovernanceProfile v1 (F18-RD-18)

- GovernanceProfile is the versioned compositional envelope from ADH-2026-033, not an assignment or effective
  context. v1 activates only `approvalPolicyRefs`, `privilegedAccessRules`, `accessReviewRules`, `exceptionRules`,
  and `auditRequirements`, with the exact closed component schemas and applicability registrations of F18-RD-18.
  ApprovalPolicy refs pin exact published versions; applicability keys are unique and ambiguous overlaps reject; no
  component is syntactically mandatory; generic maps and speculative future-domain fields are prohibited.
- Rule values may only narrow the platform ceilings and mandatory controls in REQ-F18-12 through REQ-F18-20;
  `Ineligible` is mandatory for overdue Human and AccessGroup-held Standing access while production Workload/System
  rules may select `Ineligible` or `Escalate`. A syntactically optional component becomes operationally required when
  an operation depends on it; absence, ambiguity, scope incompatibility, or an unsupported value fails closed. These
  schemas define local validation only; FEATURE-0020 exclusively selects and resolves effective profile evidence.
  `auditRequirements.DecisionRecordRequirement` may add DecisionRecord cases but `Never` cannot suppress any record
  made mandatory by REQ-F18-19.
- Lifecycle is `Draft -> Published -> Retired` with immutable published versions and informational supersession; a
  retired version is unavailable to an operation requiring current rule evidence. GovernanceProfile is an advanced
  administrator contract absent from routine user journeys, and any rule used by FEATURE-0018 is supplied as exact
  system-selected evidence (deterministic fixtures before FEATURE-0020 exists).

#### REQ-F18-19 — FEATURE-0013 adoption (F18-RD-19)

- FEATURE-0018 registers exactly three FEATURE-0013 DecisionProfiles through the approved extension process —
  `authorization-decision/v1`, `approval-decision/v1`, and `exception-decision/v1` — with the closed adopting-domain
  semantics of F18-RD-19. All three reuse FEATURE-0013's existing envelope, authority, reason, obligation, projection,
  validity, and audit-linkage mechanics unchanged; authorization has a protected full-evidence projection and a
  redacted safe-denial projection. No adopting profile changes FEATURE-0013 or performs a downstream effect, and
  FEATURE-0018 creates no competing envelope.
- Every authorization evaluation emits one protected outcome AuditEvent; material, denied, and privileged
  authorization may additionally use `authorization-decision/v1` when policy requires; RoleDefinition
  suspension/restoration/retirement remain mandatory `authorization-decision/v1` cases. Every terminal Approved/Denied
  ApprovalRequest emits exactly one mandatory `approval-decision/v1` DecisionRecord (Cancelled/Expired remain
  lifecycle AuditEvents unless an auditRequirement adds a record). Every accepted ExceptionProposal reaches an exact
  terminal `Grant | Deny` and emits exactly one mandatory `exception-decision/v1` DecisionRecord (`Deny` creates no
  ExceptionGrant). AccessReview remains its own retained LRO and adds no fourth profile; status contains only current
  state with history in AuditEvent/DecisionRecord evidence.
- Every terminal `ExceptionProposal` `Grant | Deny` emits exactly one
  `exceptionproposal.decided` AuditEvent linked, through the existing FEATURE-0013
  AuditEvent envelope and linkage semantics, to the exact immutable proposal, its
  terminal reason, the mandatory `exception-decision/v1` DecisionRecord, and the
  containing operation/correlation evidence. A terminal `Grant` additionally emits
  the separate `exceptiongrant.issued` event for its immutable ExceptionGrant in
  the same final exception atomic boundary; a terminal `Deny` publishes no
  ExceptionGrant and emits neither an issuance event nor a second decision event.
  `exceptiongrant.proposed` remains proposal-submission evidence only. This
  registration adds no additional event, resource, writer, or error beyond the one
  `exceptionproposal.decided` taxonomy type registered through FEATURE-0013's
  existing extension mechanism (ADH-2026-073).

#### REQ-F18-20 — Audit before authorization-changing publication (F18-RD-20)

- Required protected AuditEvent-obligation acceptance follows FEATURE-0013's local atomic boundary and precedes every
  authorization-changing publication enumerated by F18-RD-20 (AccessGroup, Membership, RoleDefinition, RoleAssignment,
  ApprovalPolicy/GovernanceProfile, ApprovalRequest, privileged activation/blocked/expiry/revocation, AccessReview,
  ExceptionProposal/ExceptionGrant, and completed idempotency result). Failure returns a safe internal failure and
  publishes none of the state change or completed replay result; application logs never substitute for the protected
  obligation or AuditEvent.
- Required protected AuditEvent-obligation acceptance also precedes every accepted mutation represented by the
  registered FEATURE-0018 AuditEvent taxonomy even when authorization is not yet changed, including RoleDefinition,
  ApprovalPolicy, and GovernanceProfile Draft create/update; RoleAssignmentProposal and ExceptionProposal submission;
  PrivilegedAccessRequest submission; Membership synchronization; and AccessReview campaign creation. Failure to accept
  the protected obligation for any such accepted mutation returns a safe internal failure and publishes none of the
  state change or completed replay result.
- Every authorization evaluation accepts its protected audit obligation before returning Allow or Deny; an
  authorization-evaluation audit-obligation failure returns a safe internal error, produces no `AuthorizationResult`,
  and confers no downstream authority.
- For an AccessGroup Membership assignment-effect expansion, and for approver/privileged-requester/reviewer eligibility,
  the protected evidence pins the exact enumerated fields and belongs in protected DecisionRecord/AuditEvent provenance;
  it is never copied into user-authored Membership fields and never becomes bearer authority. Approval publication and
  downstream effect publication are separate atomic operations; the approval boundary publishes the terminal state,
  mandatory `approval-decision/v1`, and audit evidence and publishes no downstream effect. Privileged activation and the
  exception controller's final `Grant | Deny` are each later, freshly authorized, all-or-nothing atomic boundaries. Every
  authorization evaluation accepts its protected audit obligation before returning Allow or Deny. The initial registered
  FEATURE-0018 AuditEvent taxonomy is exactly as enumerated by F18-RD-20 as
  corrected by approved ADH-2026-073; `exceptionproposal.decided` is its one
  added terminal-exception decision type, registered through FEATURE-0013's
  existing extension mechanism with FEATURE-0018 conformance evidence.
- Authoritative-time effectiveness is independent of successful status/event materialization: reaching any enumerated
  expiry/deadline/due boundary fails closed without grace, and audit or projection failure never extends access or
  exception effect.

#### REQ-F18-21 — Deterministic validation and safe denial (F18-RD-21)

- Validation precedence is exactly the F18-RD-21 order: bounded transport -> authentication -> structural validation ->
  local semantic/cross-field validation -> safe scope/reference resolution -> operation-required assurance and
  contextually required current membership -> assignment/action/scope/resource/delegation and
  membership-enabled-grant-envelope evaluation -> operation-required EligibilityRef and separation-of-duties evaluation
  -> approval-requirement/policy/approval/exception evidence evaluation -> concurrency and idempotency -> required
  decision/audit-obligation evidence -> publication. Unknown, deferred, and future-owned fields reject.
- Safe denial occurs before detailed inaccessible-reference disclosure; tokens, assertions, provider claims, confidential
  policy input, sensitive justification, credentials, secrets, and evaluator diagnostics never enter unsafe errors, logs,
<!-- END EXACT APPROVED-SOURCE EXCERPT -->
