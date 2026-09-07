---
doc_type: requirements
feature: FEATURE-0018
title: Governance, IAM, Approval and Exception Foundation
stage: requirements
status: draft
baseline: ARCH-2026.08-PHASE2R-CANONICAL
controlling_handoff: ADH-2026-070
conformance_handoff: ADH-2026-071
audit_taxonomy_correction_handoff: ADH-2026-073
accepted_decision: DEC-0060
change_request: ACR-2026-002
sole_architecture_authority: docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md
ai_load_priority: feature
ai_summary: Requirements translation of the approved FEATURE-0018 architecture (F18-RD-01..24). Observable behavior only; no design, schema, route, or implementation choices. Every requirement maps one-to-one to an approved REQ ledger row and every acceptance criterion maps to an approved AC ledger row.
---

# Requirements Document

<!-- Canonical spec-format heading. The authoritative feature title, structure, and
content are unchanged and continue below under the domain-specific numbered sections. -->

# FEATURE-0018 Requirements: Governance, IAM, Approval and Exception Foundation

> This document translates approved, observable behavior only. It owns intent,
> actors, scenarios, invariants, validation outcomes, security/privacy,
> compatibility, non-goals, and acceptance criteria. It does not own packages,
> files, routes, storage, algorithms, libraries, internal interfaces, or task
> decomposition. The sole behavioral authority is F18-RD-01 through F18-RD-24 as
> carried by the ADH-2026-070 three-file package and the authoritative
> FEATURE-0018 architecture. Any value not transcribable from that authority is a
> blocker, not a design or requirement choice.

## Introduction

This is a thin pointer to the authoritative content that follows; it introduces no
new requirements. The feature intent and scope are stated in the intro blockquote
above and in section 2 (Purpose and use cases).

## Requirements

This is a thin pointer to the authoritative content that follows; it introduces no
new requirements or identifiers. The complete, canonical requirement set is the
Canonical requirement ledger (REQ-F18-01..24) together with the normative details
and acceptance scenarios in section 4 (Normative requirements and acceptance
scenarios).

## Glossary

This is a thin pointer to section 3 (Terms introduced by this feature), which is the
authoritative glossary for FEATURE-0018 terms; it introduces no new terms.

## 1. Identity and stage

| Field | Value |
|---|---|
| Feature | FEATURE-0018 — Governance, IAM, Approval and Exception Foundation |
| Stage | Requirements |
| Phase / order | Phase 2R / order 8 |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Sole architecture authority | `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md` and its content-bound ADH-2026-070 package (core, Appendix A, Appendix B) |
| Accepted change records | ACR-2026-002; DEC-0060 (Accepted 2026-08-30) |
| Direct dependencies | FEATURE-0012 and FEATURE-0017; adopting registrations in FEATURE-0013 and consumed action registrations from FEATURE-0016 |
| Public routes / persistence / real adapters | None authorized |
| Stage boundary | Modifies only this requirements file; generates no design, tasks, schema, route, automation, or implementation |

This stage may translate approved behavior but may not repair, extend, reopen,
or reinterpret architecture. Any conflict with a higher-precedence authority
halts with one approved stop condition (see section 10).

## 2. Purpose and use cases

### 2.1 Purpose

FEATURE-0018 supplies Sovrunn's implementation-neutral control-plane
authorization, approval, privileged-access, access-review, and immutable
exception-evidence foundation for already-authenticated Human, Workload, and
System principals. It defines deterministic, side-effect-free, in-memory
contracts and feature-local conformance under F18-RD-01 through F18-RD-24. It is
a control-plane authorization foundation only. It is not an identity provider,
external group synchronizer, provider-native IAM bridge, policy engine, workflow
engine, effective-governance resolver, execution engine, production persistence
design, or deployment topology (F18-RD-01; architecture §2).

### 2.2 Actors

| Actor | Role in this feature |
|---|---|
| Human principal | Already-authenticated end user, administrator, requester, approver, or reviewer referenced by an exact `PrincipalRef` |
| Workload principal | Already-authenticated non-human workload identity with a responsible Human |
| System principal | Already-authenticated internal controller/provisioner identity with a responsible Human |
| Authorized access/identity/security administrator | Submits exact immutable intents for AccessGroup, Membership, RoleDefinition, RoleAssignment, ApprovalPolicy, GovernanceProfile, and AccessReview within delegated authority |
| Governance administrator | Drafts and publishes GovernanceProfile and approval/rule evidence within delegated authority |
| Owning controllers | RoleAssignment, identity/membership, publication, approval, privileged-access, review, exception, and authorization controllers as sole writers of their registered fields |
| Dependency features | FEATURE-0012 (grammar), FEATURE-0013 (evidence envelopes), FEATURE-0016 (ExecutionTarget actions), FEATURE-0017 (policy seam) — consumed by exact reference only |

### 2.3 Primary use cases

The six approved progressive-disclosure user journeys (F18-RD-24; architecture §12)
are the primary use cases and normally hide security machinery:

1. routine end-customer operation;
2. ordinary access assignment;
3. privileged-access request;
4. approver decision;
5. access review; and
6. exception request.

Every referenced use case resolves through the exact contracts, actions, writers,
and validation precedence defined by the normative requirements in section 4 and
is bounded by the exclusions in section 7.

## 3. Terms introduced by this feature

These terms are introduced or completed by FEATURE-0018 and are defined solely by
their controlling F18-RD authority. This section names them; it does not redefine
inherited dependency contracts.

### 3.1 Persistent resource contracts (independently addressable)

| Term | Meaning (authority) |
|---|---|
| `AccessGroup` | Non-authenticating, access-only authorization subject scoped to exactly one Organization or CloudProvider, with one responsible Human owner and direct principal membership only (F18-RD-04) |
| `Membership` | Non-authorizing record of Organization/CloudProvider belonging or one direct same-scope AccessGroup relationship (F18-RD-05) |
| `RoleDefinition` | Immutable published version composed only of registered canonical Sovrunn actions (F18-RD-07) |
| `RoleAssignment` | Scoped, explicitly valid grant of one exact published role version to one role holder (F18-RD-08) |
| `PrivilegedAccessRequest` | Justified, individual, time-bounded JIT or break-glass privileged-access request (F18-RD-14, F18-RD-15) |
| `AccessReview` | Immutable-snapshot review with retained usage/remediation evidence (F18-RD-16) |
| `ApprovalPolicy` | Immutable published bounded approval stages, eligibility, quorum, expiry, and separation of duties (F18-RD-12) |
| `ApprovalRequest` | Retained approval evidence for one bounded FEATURE-0018 intent or proposal (F18-RD-13) |
| `ExceptionGrant` | Immutable, scoped, time-bounded exception or linked revocation evidence (F18-RD-17) |
| `GovernanceProfile` v1 | Versioned composition of only FEATURE-0018-owned governance references and constraints (F18-RD-18) |

### 3.2 Supporting values (no endpoint, lifecycle, storage authority, or user journey) (F18-RD-02)

`PrincipalRef`, `RoleHolderRef`, `AccessGroupRef`, `AssuranceEvidence`,
`AuthorizationInput`, `AuthorizationResult`, `ApprovalRequirement`,
`EligibilityRef`, `RoleAssignmentProposal`, `MembershipEnabledGrantEnvelope`,
`ExceptionProposal`, `UsageEvidenceSummary`, `GovernanceApplicability`, and
`ActionTargetBinding`.

### 3.3 Closed initial supporting enums (F18-RD-02)

```text
MembershipContextKind = Organization | CloudProvider | AccessGroup
MembershipType        = Standard | Guest
PrivilegedAccessMode  = JIT | BreakGlass
ApprovalDecision      = Approve | Deny
ExceptionGrantEffect  = Grant | Revoke
AssignmentValidityMode = Standing | TimeBound
ActionTargetMode      = ExactResource | CreateParent | ScopeOnly
```

These are canonical contract values, not personas, workflow states,
provider-native values, or extensible strings.

### 3.4 Introduced composition concepts

| Term | Meaning (authority) |
|---|---|
| Grant union + guardrail intersection | Authorization is the union of independently applicable grants constrained by the intersection of mandatory guardrails; restrictions never synthesize a grant (F18-RD-09) |
| `MembershipEnabledGrantEnvelope` | Operation-local set of exact group-held assignment action/reach/validity tuples newly enabled or extended by an AccessGroup Membership assignment-effect expansion (F18-RD-09) |
| Grantor ceiling | Per-action delegable-witness containment of scope, target-binding mode, exact resource reach, validity, and delegation capability with no cross-assignment synthesis (F18-RD-09) |
| `MembershipRelationshipKey` | Deterministic uniqueness key over principal, context kind, scope, and (for AccessGroup) group ref (F18-RD-05) |
| `ReviewBeneficiaryHumanSet` | Bounded set of Humans who benefit from reviewed authority and are excluded from access-preserving/replacing review decisions (F18-RD-16) |
| Two-plane native-IAM boundary | Effective native effect = valid Sovrunn control-plane authorization intersected with independent native IAM authorization (F18-RD-06) |

## Canonical requirement ledger

The rows below are copied exactly, once, from the approved FEATURE-0018 canonical
requirement generation ledger (feature file §6.1). Each identifier maps one-to-one
to its controlling architecture decision group and must not be repurposed,
renumbered, merged, split, or paraphrased.

| Requirement ID | Architecture decision | Exact title |
|---|---|---|
| REQ-F18-01 | F18-RD-01 | Current-feature-only authority |
| REQ-F18-02 | F18-RD-02 | Closed contract inventory and profiles |
| REQ-F18-03 | F18-RD-03 | Stable principal identity |
| REQ-F18-04 | F18-RD-04 | Direct principal and scoped AccessGroup assignments |
| REQ-F18-05 | F18-RD-05 | Membership is non-authorizing |
| REQ-F18-06 | F18-RD-06 | Distributed action ownership and central role composition |
| REQ-F18-07 | F18-RD-07 | Versioned RoleDefinition |
| REQ-F18-08 | F18-RD-08 | Scoped RoleAssignment |
| REQ-F18-09 | F18-RD-09 | Deterministic scoped authorization composition |
| REQ-F18-10 | F18-RD-10 | FEATURE-0017 adoption without reinterpretation |
| REQ-F18-11 | F18-RD-11 | Bounded FEATURE-0017 subject/target use |
| REQ-F18-12 | F18-RD-12 | Bounded ApprovalPolicy |
| REQ-F18-13 | F18-RD-13 | Immutable terminal approval evidence |
| REQ-F18-14 | F18-RD-14 | JIT privileged access authorizes one temporary grant |
| REQ-F18-15 | F18-RD-15 | Constrained break-glass access |
| REQ-F18-16 | F18-RD-16 | Snapshot-based AccessReview |
| REQ-F18-17 | F18-RD-17 | Bounded immutable exception evidence |
| REQ-F18-18 | F18-RD-18 | FEATURE-0018-limited GovernanceProfile v1 |
| REQ-F18-19 | F18-RD-19 | FEATURE-0013 adoption |
| REQ-F18-20 | F18-RD-20 | Audit before authorization-changing publication |
| REQ-F18-21 | F18-RD-21 | Deterministic validation and safe denial |
| REQ-F18-22 | F18-RD-22 | Deterministic in-memory foundation and local conformance |
| REQ-F18-23 | F18-RD-23 | Standards-validation gate |
| REQ-F18-24 | F18-RD-24 | Progressive and normally hidden user friction |

## Canonical acceptance ledger

The rows below are copied exactly, once, from the approved FEATURE-0018 canonical
acceptance generation ledger (feature file §6.2). `AC-F18-38` preserves the retired
duplicate as excluded and must never be treated as an active acceptance obligation.
Detailed scenarios in section 4.2 may expand these rows but may not change their
meaning.

| Acceptance ID | Architecture scenario | Exact scenario title | Disposition |
|---|---|---|---|
| AC-F18-01 | F18-SCN-01 | Active employee receives project-scoped developer access | Active |
| AC-F18-02 | F18-SCN-02 | Employee email changes | Active |
| AC-F18-03 | F18-SCN-03 | Employee is suspended during an active assignment | Active |
| AC-F18-04 | F18-SCN-04 | Employee leaves and later rejoins | Active |
| AC-F18-05 | F18-SCN-05 | Workload calls a control-plane operation | Active |
| AC-F18-06 | F18-SCN-06 | System controller performs an owned transition | Active |
| AC-F18-07 | F18-SCN-07 | Consultant belongs to supplier and customer contexts | Active |
| AC-F18-08 | F18-SCN-08 | External IdP group claim and provisioned AccessGroup are distinguished | Active |
| AC-F18-09 | F18-SCN-09 | Organization administrator delegates project authority | Active |
| AC-F18-10 | F18-SCN-10 | ExecutionTarget operator reads and qualifies one target | Active |
| AC-F18-11 | F18-SCN-11 | Target operator tries to retire with qualifier role | Active |
| AC-F18-12 | F18-SCN-12 | Organization wants create-only but not retire authority | Active |
| AC-F18-13 | F18-SCN-13 | Engineer requests two-hour production DB administration | Active |
| AC-F18-14 | F18-SCN-14 | Engineer approves own privileged request | Active |
| AC-F18-15 | F18-SCN-15 | Required approver is revoked before deciding | Active |
| AC-F18-16 | F18-SCN-16 | Approval and expiry occur concurrently | Active |
| AC-F18-17 | F18-SCN-17 | Same approved JIT request is activated twice | Active |
| AC-F18-18 | F18-SCN-18 | Emergency access is needed without normal approver availability | Active |
| AC-F18-19 | F18-SCN-19 | Access review finds stale administrator access | Active |
| AC-F18-20 | F18-SCN-20 | Assignment changes during review | Active |
| AC-F18-21 | F18-SCN-21 | Incident requires change-window exception plus DB role | Active |
| AC-F18-22 | F18-SCN-22 | Request attempts to except tenant isolation | Active |
| AC-F18-23 | F18-SCN-23 | Exception is revoked before its planned expiry | Active |
| AC-F18-24 | F18-SCN-24 | Two exception grants overlap | Active |
| AC-F18-25 | F18-SCN-25 | FEATURE-0017 returns Allow but role is expired | Active |
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
  or projections. Idempotency binds actor, operation, exact target, scope, and request digest; replay rechecks current
  authentication, authorization, and safe target visibility; same key/different content conflicts; required audit failure
  publishes no state change or completed replay result.

#### REQ-F18-22 — Deterministic in-memory foundation and local conformance (F18-RD-22)

- FEATURE-0018 uses synthetic fixtures immutable to consumers, rejects invalid fixture sets before evaluation, obtains
  time from a deterministic UTC source, and performs zero network or external I/O. Fakes cover the enumerated
  approval/JIT/break-glass/expiry/revocation/review/exception paths without simulating an IdP, policy engine, workflow
  engine, or CloudProvider.
- Feature-local conformance covers exactly the positive, negative, race, replay, audit-failure, and redaction vectors
  enumerated by F18-RD-22, including the membership-enabled grant envelope, grantor ceiling, Human-only EligibilityRef,
  the three FEATURE-0016 ExecutionTarget bindings and intrinsic classes, mode-specific review timing, mandatory
  DecisionRecord emission, the indirect-beneficiary conflict set, and the Sovrunn/native-IAM intersection boundary. Only
  FEATURE-0018-owned conformance identifiers and local proof count as feature acceptance; no integration proof substitutes
  for local proof. Executable conformance is a downstream feature-gate obligation; the registered Slice-0 conformance
  cases with exact fields are transcribed in the Exact conformance semantics ledger below.

#### REQ-F18-23 — Standards-validation gate (F18-RD-23)

- Before final architecture approval, every applicable invariant maps to FEATURE-0018 proof, a reused dependency, or a
  named exclusion under the standards enumerated by F18-RD-23 (NIST SP 800-207, 800-53 Rev. 5, 800-63-4, 800-162; ISO/IEC
  27001/27002; CSA CCM; CIS Controls 5/6/8; OWASP ASVS). OIDC, OAuth, SCIM, WebAuthn, and SPIFFE are paper-only
  future-adapter compatibility checks that authorize no adapter design or implementation.
- The applicable ASVS profile and CIS/CSA shared-responsibility mappings are closed before final architecture approval;
  formal ISO conformity claims remain outside this handoff without a licensed clause-level assessment. This is an
  architecture-gate obligation; requirements neither restate standard clauses as behavior nor claim unassessed conformity.

#### REQ-F18-24 — Progressive and normally hidden user friction (F18-RD-24)

- Routine authorization resolves silently; users do not manage schema versions, policy digests, issuer/subject,
  provenance, external object refs, exact role/policy versions, authorization values, evidence refs, or correlation IDs
  (all system-owned). The six user journeys are exactly those in F18-RD-24, each exposing the smallest
  risk-proportional projection. Security administrators select exact eligible Humans in the applicable policy or rule;
  end customers do not manage EligibilityRef values or group/role-policy linkage.
- PrivilegedAccessRequest/ApprovalRequest/temporary RoleAssignment/evidence are one requester-facing journey, and
  ExceptionProposal/ApprovalRequest/ExceptionGrant are one requester-facing journey; separate protected resources remain
  internally auditable, and no facade resource is added or exact underlying contract weakened.

### 4.2 Acceptance-scenario coverage mapping

Every active approved AC is mapped explicitly below to its primary controlling
REQ(s). This mapping is outside the canonical acceptance ledger; the ledger row and
this mapping together constitute the completeness claim. `AC-F18-38` is Excluded and
is enforced as a negative acceptance case (section 7.3), never as an active obligation.

| Acceptance ID | Primary REQ mapping | Expected observable outcome (abbreviated; full meaning is the ledger row) |
|---|---|---|
| AC-F18-01 | REQ-F18-03, 04, 05, 08, 09 | Exact-scope grant succeeds; unrelated scope denies |
| AC-F18-02 | REQ-F18-03 | Identity remains stable through issuer/subject when email changes |
| AC-F18-03 | REQ-F18-05, 08, 09 | Suspension makes access ineffective immediately; history remains |
| AC-F18-04 | REQ-F18-05, 08 | Revoked/expired assignments do not silently reactivate on rejoin |
| AC-F18-05 | REQ-F18-03, 08, 09 | Workload identity is independently scoped and audited |
| AC-F18-06 | REQ-F18-02, 03, 20 | System identity cannot override sole-writer rules |
| AC-F18-07 | REQ-F18-05, 08 | Guest sponsorship, target scope, expiry, and attribution are unambiguous |
| AC-F18-08 | REQ-F18-04, 05, 09 | Raw claim alone denies; provisioned canonical group + active scoped assignment may allow |
| AC-F18-09 | REQ-F18-08, 09 | Grantor cannot exceed one independently applicable per-action witness |
| AC-F18-10 | REQ-F18-06 | Consumed FEATURE-0016 actions + optional resourceRef authorize only target-bound operations |
| AC-F18-11 | REQ-F18-06 | Denied because `executiontarget.write` is absent |
| AC-F18-12 | REQ-F18-06 | Recorded create/retire granularity limitation; no invented action |
| AC-F18-13 | REQ-F18-12, 13, 14 | Approved request creates one temporary assignment that expires exactly |
| AC-F18-14 | REQ-F18-12 | Self-approval denied and audited without activation |
| AC-F18-15 | REQ-F18-12, 13 | Current eligibility recheck denies the revoked approver's decision |
| AC-F18-16 | REQ-F18-13, 21 | First valid terminal CAS wins; no unauthorized assignment |
| AC-F18-17 | REQ-F18-14, 21 | One temporary assignment; safe replay or deterministic conflict |
| AC-F18-18 | REQ-F18-15, 16 | Pre-authorized break-glass lasts at most one hour with immediate audit and 24h review |
| AC-F18-19 | REQ-F18-16 | Explicit exactly-once revocation with retained evidence |
| AC-F18-20 | REQ-F18-16 | Explicit Stale outcome performs no remediation |
| AC-F18-21 | REQ-F18-08, 17 | Both bounded exception and applicable role are required |
| AC-F18-22 | REQ-F18-17 | Non-exceptionable control denies |
| AC-F18-23 | REQ-F18-17 | Immutable effect=Revoke record makes the exact grant ineffective |
| AC-F18-24 | REQ-F18-17 | Active overlap rejected for same control/subject/scope/interval |
| AC-F18-25 | REQ-F18-09, 10 | FEATURE-0017 Allow does not override an expired role; final authorization denies |
| AC-F18-26 | REQ-F18-10, 11 | RequiresApproval evaluated only for a compatible pinned subject; request only via explicit flow |
| AC-F18-27 | REQ-F18-10, 21 | Indeterminate fails closed with safe decision/audit behavior |
| AC-F18-28 | REQ-F18-20 | Audit append failure publishes no temporary grant or completed replay |
| AC-F18-29 | REQ-F18-24, 09 | Identity/grants/guardrails/provenance/audit resolve silently; only business input/result |
| AC-F18-30 | REQ-F18-21 | Safe denial reveals no target existence or sensitive policy detail |
| AC-F18-31 | REQ-F18-08, 16 | Policy-permitted Standing assignment remains revocable and enters periodic review |
| AC-F18-32 | REQ-F18-05, 08 | Guest assignment without both TimeBound timestamps is denied |
| AC-F18-33 | REQ-F18-09 | Grant union cannot bypass a restrictive guardrail; denies with exact provenance |
| AC-F18-34 | REQ-F18-07 | Suspended version makes its assignments ineffective; other versions unaffected; restoration audited |
| AC-F18-35 | REQ-F18-09 | AuthorizationResult is non-bearer and cannot authorize another target/principal |
| AC-F18-36 | REQ-F18-05 | Stale provenance/absent owner fails safely; raw claim cannot repair it |
| AC-F18-37 | REQ-F18-16 | Absent usage is not proof of non-use; no silent workload revocation |
| AC-F18-39 | REQ-F18-24, 08 | Administrator supplies four inputs; refs/version/provenance/audit are system-managed; forged system fields rejected |
| AC-F18-40 | REQ-F18-24, 14 | Requester sees one request status and expiry; internal request/evidence separation hidden |
| AC-F18-41 | REQ-F18-24, 12, 13 | Approver sees actionable context; exact versions and concurrency controls hidden |
| AC-F18-42 | REQ-F18-24, 16 | Reviewer sees holder/role/scope/validity/usage summary and Retain/Revoke/Replace; snapshot mechanics hidden |
| AC-F18-43 | REQ-F18-24, 17 | User sees one request status and effective period; grant/evidence records internal |
| AC-F18-44 | REQ-F18-08, 09 | `ExactResource(A)` cannot grant resource B or scope-wide reach; no partial-assignment synthesis |
| AC-F18-45 | REQ-F18-16 | Direct/indirect beneficiary receives `REVIEWER_CONFLICT`; unresolved expansion fails closed; Revoke available |
| AC-F18-46 | REQ-F18-05, 09 | TimeBound admin cannot enable Standing group-held effect; mismatched witnesses deny; empty-group grants nothing |
| AC-F18-47 | REQ-F18-05, 12, 14, 16 | AccessGroupRef/Membership/group-derived authority structurally rejected as eligibility; exact Human still needs the action |
| AC-F18-48 | REQ-F18-16 | Non-empty OR list accepts exact Humans only; every indirect/stale/incompatible case denies; rechecked before remediation |
| AC-F18-49 | REQ-F18-12, 14, 16 | Role grant/reference cannot change approver/reviewer/JIT/break-glass eligibility; only exact Humans qualify |

## 5. Security, privacy, compatibility, and operational requirements

### 5.1 Security

| ID | Requirement | Authority |
|---|---|---|
| SEC-01 | Default deny: no grant, invalid/expired evidence, raw external group claim, nested/dynamic group membership, unresolved or cross-tree scope, resource mismatch, conflicting evidence, or `Indeterminate` denies. Policy Deny overrides a role. | REQ-F18-09 |
| SEC-02 | Least privilege and no synthesis: authority is grant union constrained by guardrail intersection; restrictions from different assignments never synthesize a grant; delegation obeys the per-action grantor ceiling. | REQ-F18-08, REQ-F18-09 |
| SEC-03 | Uniform activation floor: every interactive privileged activation (normal JIT and break-glass) requires phishing-resistant AAL2-or-higher AssuranceEvidence no more than fifteen minutes old; a rule may strengthen but never weaken it. | REQ-F18-14, REQ-F18-15 |
| SEC-04 | Separation of duties and no self-dealing: approver conflict sets and the `ReviewBeneficiaryHumanSet` exclude direct and indirect beneficiaries (including accountable Humans behind direct Workload/System AccessGroup members) and are rechecked at decision acceptance and effect publication. | REQ-F18-12, REQ-F18-16 |
| SEC-05 | Single writer: each resource has one accepted-intent/spec writer and one status/result writer; the RoleAssignment controller is the sole RoleAssignment publisher; clients never write status. | REQ-F18-02, REQ-F18-08 |
| SEC-06 | Eligibility is exact-Human only: approver, requester, and reviewer `EligibilityRef` lists are non-empty, OR-composed, exact Human PrincipalRefs; roles, assignments, groups, Membership, and claims never satisfy or create eligibility, which never grants the separately required canonical action. | REQ-F18-12, REQ-F18-14, REQ-F18-16 |
| SEC-07 | Time-bounded privilege and no grace: privileged, guest, JIT, and break-glass access is TimeBound and becomes ineffective at exact expiry; audit/projection failure never extends effectiveness. | REQ-F18-08, REQ-F18-14, REQ-F18-20 |
| SEC-08 | Two-plane native boundary: a native effect requires both Sovrunn authorization and independent native IAM authorization; neither Allow substitutes for nor overrides the other's Deny; no provider-native principal/role/policy/permission/credential is a FEATURE-0018 resource or output. | REQ-F18-06 |
| SEC-09 | Audit before authorization-changing publication: required protected AuditEvent-obligation acceptance precedes every enumerated publication; failure publishes nothing. | REQ-F18-20 |

### 5.2 Privacy and redaction

| ID | Requirement | Authority |
|---|---|---|
| PRIV-01 | Safe denial precedes detailed inaccessible-reference disclosure; error responses never confirm existence of an inaccessible resource. | REQ-F18-21 |
| PRIV-02 | Tokens, assertions, provider claims, confidential policy input, sensitive justification, credentials, secrets, and evaluator diagnostics never enter unsafe errors, logs, or projections; raw authentication material is never retained. | REQ-F18-03, REQ-F18-21 |
| PRIV-03 | Authorization has a protected full-evidence projection and a redacted safe-denial projection; AuditEvent/DecisionRecord evidence uses FEATURE-0013 redaction and projection semantics. | REQ-F18-19 |

### 5.3 Compatibility (inherited-contract reuse by exact reference only)

| ID | Requirement | Authority |
|---|---|---|
| COMPAT-01 | FEATURE-0012 metadata, seven ScopeKinds, `metadata.scopeRef` as sole scope authority, typed references, validation, safe denial, redaction, optimistic concurrency, idempotency, and Problem behavior are reused unchanged; FEATURE-0018 adds no second scope authority. | REQ-F18-02, REQ-F18-21 |
| COMPAT-02 | FEATURE-0013 DecisionRecord, DecisionProfile, EvaluationResult, and AuditEvent envelopes are reused unchanged; FEATURE-0018 only registers the three profiles and the bounded AuditEvent taxonomy. | REQ-F18-19, REQ-F18-20 |
| COMPAT-03 | FEATURE-0016 ExecutionTarget action meanings and lifecycle are consumed unchanged; the combined create/retire granularity of `executiontarget.write` is a recorded limitation; any privileged-handling classification is additive metadata, not a reinterpretation. | REQ-F18-06 |
| COMPAT-04 | FEATURE-0017 PolicyEvaluationRequest/Result, PolicyEngineAdapter, deterministic digest-keyed fake, and pure mapper are reused unchanged within the v1 subject/target limitation; policy results constrain but never create authority. | REQ-F18-10, REQ-F18-11 |
| COMPAT-05 | Backward compatibility impact is documentation/canonical-contract only; FEATURE-0018 has no implemented runtime objects requiring data migration and introduces no dual authority. | ADH-2026-070 phase impact |

### 5.4 Operational

| ID | Requirement | Authority |
|---|---|---|
| OPS-01 | Time is obtained only from a deterministic authoritative UTC source; all expiry/deadline/due calculations use it. | REQ-F18-03, REQ-F18-20, REQ-F18-22 |
| OPS-02 | The foundation is fully in-process and side-effect-free: zero network or external I/O; synthetic fixtures are immutable to consumers and invalid fixture sets are rejected before evaluation. | REQ-F18-22 |
| OPS-03 | Idempotency and concurrency behavior is deterministic: a same key/same digest replay returns the stored successful result only after the F18-RD-21 replay rechecks pass — current authentication, current authorization, and safe target visibility are re-evaluated on every replay and a failed recheck yields the corresponding safe error rather than the stored result; same key/different content conflicts; a required audit-obligation failure publishes no state change and no completed replay result; concurrent grant-producing operations conflict or retry against a coherent snapshot. | REQ-F18-13, REQ-F18-21 |
| OPS-04 | Owning controllers retry only permitted retained status, event, or non-authoritative projection materialization after audit/projection failure; effectiveness never depends on successful status/event materialization. | REQ-F18-20 |

## 6. Edge cases

These edge cases are observable-behavior expectations derived from the cited
authorities; each expands an approved AC row without changing its meaning.

| Edge case | Expected outcome | Authority / AC |
|---|---|---|
| Approval decision and subject expiry occur concurrently | First valid terminal CAS wins; no unauthorized assignment is published | REQ-F18-13, REQ-F18-21 / AC-F18-16 |
| Same approved JIT request activated twice | Exactly one temporary assignment; safe replay or deterministic conflict; validity is never extended | REQ-F18-14, REQ-F18-21 / AC-F18-17 |
| Required AuditEvent append fails during JIT activation | No temporary grant, no `Active` request, and no completed replay is published | REQ-F18-20 / AC-F18-28 |
| Two ExceptionGrants overlap on the same control/subject/scope/interval | The new overlapping grant is rejected; no merge or precedence is inferred | REQ-F18-17 / AC-F18-24 |
| Exception revoked before planned expiry | Linked immutable `effect=Revoke` shortens the effective interval without mutation; revocation at/after expiry is rejected | REQ-F18-17 / AC-F18-23 |
| FEATURE-0017 returns Indeterminate | Successfully mapped to Deny with stable protected provenance; safe decision/audit behavior | REQ-F18-10, REQ-F18-09 / AC-F18-27 |
| Published role version suspended during use | Every assignment pinned to that exact version becomes ineffective immediately; other versions unaffected; restoration is auditable | REQ-F18-07 / AC-F18-34 |
| AuthorizationResult replayed for a different target/principal | Rejected; a result is non-bearer and cannot authorize another actor/action/target/scope/instant | REQ-F18-09 / AC-F18-35 |
| Federated membership provenance stale or responsible owner absent | Fails safely where current provenance/ownership is required; a raw external claim cannot repair it | REQ-F18-05 / AC-F18-36 |
| Access review has no or incomplete usage telemetry | Absence is not proof of non-use unless coverage is complete; no silent production-workload revocation | REQ-F18-16 / AC-F18-37 |
| Standing beneficiary attempts to certify retained access | `REVIEWER_CONFLICT`; unresolved beneficiary expansion fails closed with `REVIEW_BENEFICIARY_UNRESOLVED`; independently authorized Revoke remains available | REQ-F18-16 / AC-F18-45 |
| Temporary (TimeBound) membership administrator attempts durable group-derived access | The derived envelope requires Standing witnesses; a TimeBound administrator cannot enable a Standing group-held effect; mismatched/missing witnesses deny; empty-group addition grants nothing | REQ-F18-05, REQ-F18-09 / AC-F18-46 |
| Membership/role administrator attempts to manufacture approval/reviewer/JIT/break-glass eligibility | AccessGroupRef, Membership, and role/assignment-derived authority are structurally rejected as eligibility; an exact named Human still requires the separate canonical action | REQ-F18-12, REQ-F18-14, REQ-F18-16 / AC-F18-47, AC-F18-49 |
| Guest assignment omits expiry or a TimeBound timestamp | Validation denies; guest assignments must be TimeBound with both timestamps and cannot exceed the pinned guest Membership `expiresAt` | REQ-F18-05, REQ-F18-08 / AC-F18-32 |
| Reviewer eligibility missing, indirect, stale, or scope/target-incompatible | Denies; only exact Humans on a non-empty OR list qualify; eligibility and beneficiary conflict are rechecked before remediation | REQ-F18-16 / AC-F18-48 |
| Snapshot empty, changed, or unresolvable on entry to InProgress | Fails closed; the campaign does not silently complete | REQ-F18-16 |
| ApprovalPolicy runtime eligible set smaller than quorum | No approval; the stage remains pending until eligibility recovers or the request reaches its immutable `expiresAt` and becomes `Expired` | REQ-F18-12 |

## 7. Non-goals and adjacent-feature exclusions

This section is the sole location that names prohibited or stale concepts. Active
requirements and acceptance cases in sections 4–6 refer to this section rather
than repeating prohibited names.

### 7.1 Feature-level non-goals (F18-RD-01; architecture §2, §9.2)

FEATURE-0018 is not, and does not build: an identity provider; token processing
or validation; federation or SCIM synchronization; external group synchronization;
nested or dynamic AccessGroup membership; external claims as direct authorization
grants; provider-, Kubernetes-, or engine-native principals, groups, roles,
policies, permissions, credentials, or schemas; a real policy engine or workflow
engine, adapter selection, routing, retry, or credential flow; a second decision,
audit, authorization, scope, GovernanceProfile, effective-governance, or
RoleAssignment writer authority; user-authored status, evidence, provenance,
arbitrary condition expressions, wildcard permissions, or deny assignments;
production persistence, a public facade resource, a raw secret, a protected native
handle, an external call, or an infrastructure effect; and any requirements,
design, tasks, implementation, or future-feature artifact produced from this
stage. No `ExceptionRequest`, generic workflow, facade, condition language,
external-group resource, provider-native IAM resource, or credential resource is
introduced.

### 7.2 Adjacent-feature exclusions (enforceable)

| Excluded owner | Excluded concepts (not implemented, referenced, or pulled forward) |
|---|---|
| FEATURE-0019 | Sovereignty profiles, regulatory policy bundles, sovereignty facts and evidence, placement-facing sovereignty interpretation |
| FEATURE-0020 | ProfileAssignment, governance inheritance and conflict resolution, exception application, `EffectiveGovernanceContext` construction |
| FEATURE-0021 | CloudEnrollment, entitlement, quota, provider-selection intent |
| FEATURE-0022 | Service requirements, service-region eligibility, placement-profile semantics |
| FEATURE-0023 | Candidate-set policy adoption, ranking and placement, DecisionRecord publication for downstream placement |
| FEATURE-0024 | Production adapter selection, plugin execution, provider-native IAM execution, provisioning and realization |
| FEATURE-0025 | Customer-safe, CloudProvider-safe, and AI-readable explanation projection |
| FEATURE-0026 | Cross-feature orchestration and Slice 0 integration-proof ownership |

FEATURE-0018 owns the `GovernanceProfile` v1 envelope and its component semantics
only; FEATURE-0020 retains selection and resolution. No downstream status field,
metric, side effect, controller, adapter/plugin execution, or runtime proof is
imported merely because it is visible in a shared Slice 0 document.

### 7.3 Negative acceptance cases

| Negative case | Required outcome | Authority / AC |
|---|---|---|
| Raw external IdP group claim used alone as authorization | Denied; only trusted-provisioned canonical AccessGroup + Membership + active scoped assignment may allow | REQ-F18-04, REQ-F18-05, REQ-F18-09 / AC-F18-08 |
| Attempt to except a non-exceptionable control (e.g. tenant/scope isolation) | Denied as non-exceptionable | REQ-F18-17 / AC-F18-22 |
| Self-approval of one's own privileged request | Denied and audited without activation | REQ-F18-12 / AC-F18-14 |
| Delegating `ExactResource(A)` authority to resource B or scope-wide | Denied; no partial-assignment synthesis broadens reach | REQ-F18-08, REQ-F18-09 / AC-F18-44 |
| Attempt to retire an ExecutionTarget with only the qualifier role | Denied because `executiontarget.write` is absent | REQ-F18-06 / AC-F18-11 |
| `AC-F18-38` treated as an active obligation | Rejected; it is the retired duplicate of F18-SCN-29, carries no independent scenario/evidence authority, and all routine low-risk journey evidence is attributed to AC-F18-29 | Canonical acceptance ledger / AC-F18-38 (Excluded) |
| Introducing any prohibited or stale concept from section 7.1/7.2 into an active contract | Rejected as an architecture-drift/leakage failure | REQ-F18-01, REQ-F18-02 |

## 8. Architecture/decision/risk traceability by exact ID

### 8.1 Controlling authorities

- Decision/change records: DEC-0060 (Accepted 2026-08-30); ACR-2026-002.
- Controlling handoff package: ADH-2026-070 core, Appendix A (semantic contract),
  Appendix B (registries and evidence); ADH-2026-071 controls conformance IDs
  and traceability only; ADH-2026-073 (Accepted 2026-09-03) registers the single
  `exceptionproposal.decided` AuditEvent taxonomy type for the terminal
  `ExceptionProposal` decision boundary through FEATURE-0013's existing extension
  mechanism (REQ-F18-19, REQ-F18-20).
- Sole architecture authority: `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md`.
- Baseline: ARCH-2026.08-PHASE2R-CANONICAL.

### 8.2 Requirement-to-decision traceability

Every REQ-F18-NN maps one-to-one to F18-RD-NN of the same number (see the Canonical
requirement ledger). Reused-authority references, without redefinition:

| Reused authority | Consumed by | Exact reference IDs |
|---|---|---|
| FEATURE-0012 common contracts | REQ-F18-02, 21 | DEC-0026; ADH-2026-012; DEC-0037 (seven scopes) |
| FEATURE-0013 evidence envelopes | REQ-F18-19, 20 | DEC-0043; ADH-2026-017; DEC-0044 (immutable published definitions) |
| FEATURE-0016 ExecutionTarget actions | REQ-F18-06 | ADH-2026-058 through ADH-2026-066 |
| FEATURE-0017 evaluation seam | REQ-F18-10, 11 | DEC-0028; DEC-0036; ADH-2026-067 through ADH-2026-069 |
| GovernanceProfile envelope | REQ-F18-18 | DEC-0050; ADH-2026-033 |
| Enterprise governance / non-weakenable governance / customer simplicity / extensibility | REQ-F18-04, 08, 18, 24 | DEC-0018; DEC-0019; DEC-0020; DEC-0035; DEC-0048 |
| Future-adapter deferral (non-authoritative) | REQ-F18-01, 23 | DEC-0022; DEC-0036; RFC-0012; RFC-0021; RFC-0022; RFC-0023 |

### 8.3 Schema/writer/state/error traceability (Slice-0 registry)

| Concern | Exact IDs |
|---|---|
| Owned schema IDs | VS0-SCHEMA-020 (PrincipalRef), 021 (Membership), 022 (RoleDefinition), 023 (RoleAssignment), 024 (GovernanceProfile), 062 (AccessGroup), 063 (PrivilegedAccessRequest), 064 (AccessReview), 065 (ApprovalPolicy), 066 (ApprovalRequest), 067 (ExceptionGrant), 068 (FEATURE-0018 supporting values) |
| Owned writer IDs | VS0-WRITER-008 (RoleAssignment-controller, sole publisher of `RoleAssignment.spec`/`RoleAssignment.status`), VS0-WRITER-022 (identity-membership-controller, sole writer of `Membership.spec`/`Membership.protected`/`Membership.status`) |
| Consumed/shared writer IDs | VS0-WRITER-007 (authorized-governance-publisher) is consumed only on its manifest-pinned registered FEATURE-0018 paths `GovernanceProfile.spec` and `RoleDefinition.spec`; VS0-WRITER-011 (authorized-decision-or-audit-producer) is consumed unchanged for FEATURE-0013 `DecisionRecord.record`/`AuditEvent.record` publication (no competing writer is introduced) |
| Out-of-scope writer IDs | VS0-WRITER-023 (delegated-governance-admin, `ProfileAssignment.spec`) is FEATURE-0020-owned and is not FEATURE-0018-owned or consumed |
| Decision profiles (FEATURE-0013) | `authorization-decision/v1`, `approval-decision/v1`, `exception-decision/v1` (REQ-F18-19) |
| Error/failure mappings | Exactly two inherited security failure mappings apply one-to-one: VS0-F01 → VS0-CF-F01 → `AUTH_REQUIRED`/401/`VS0_AUTH_REQUIRED`; VS0-F02 → VS0-CF-F02 → `RESOURCE_NOT_FOUND`/404/`VS0_AUTHORIZATION_SAFE_DENIAL`. VS0-F01 and VS0-F02 are never associated with any local `VS0-CF-F18-*` case; local cases carry their own exact registered outcomes. Only FEATURE-0012 top-level codes are used; Slice-0 violation codes (e.g. `REVIEWER_CONFLICT`, `REVIEW_BENEFICIARY_UNRESOLVED`) appear only in `violations[].code`; no `429`/quota HTTP code. |
| State authorities | Resource lifecycles and effective projections are owned by their REQ/F18-RD (REQ-F18-04, 05, 07, 08, 12, 13, 14, 16, 18); no state-machine authority is substituted by a conformance case. |

### 8.4 Risk and finding traceability (exact IDs)

Independent-review findings closed by the approved package (architecture §1;
ADH-2026-070 summary): F18-SEC-001 through F18-SEC-005; renewal findings F18-REN-001,
F18-REN-002; F18-REN2-001, F18-REN2-002; F18-REN3-001, F18-REN3-002; F18-REN4-001
(renewal-05 verified with no blocking findings). Applicable architecture risks and
their preventive/detective/corrective controls are the risk-mitigation rows of the
feature file §5 and are enforced through REQ-F18-01 through REQ-F18-24. Residual risk
is Medium until independent security review and executable proof; expected Low-to-Medium
after feature-gate evidence.

### 8.5 Conformance and scenario evidence (exact IDs)

Linked acceptance criteria: F18-IVM-A01 through F18-IVM-I06 and F18-SCN-01 through
F18-SCN-49 with F18-SCN-38 retired into F18-SCN-29 (mirrored one-to-one by the Canonical
acceptance ledger AC-F18-01..49). F18-SCN-29 and F18-SCN-39..43 are architecture journey
evidence; F18-SCN-44..49 are security-closure evidence; executable conformance is a
downstream feature-gate obligation.

### 8.6 Per-requirement architecture traceability matrix (REQ-F18-01..24)

Each active requirement maps below to its controlling F18-RD/DEC authority, its
owner, and the applicable **registered** schema/writer/state/error/conformance
identifiers by their exact, fully qualified IDs. Where no registered
schema/writer/state/error ID applies, the cell states `n/a — <reason>` and cites
the controlling F18-RD separately; controller names, lifecycle prose, and inferred
runtime errors are never substituted for a registered ID.

Registered writers relevant to FEATURE-0018 are exactly four. Two are owned:
VS0-WRITER-008 (RoleAssignment-controller, paths `RoleAssignment.spec`/
`RoleAssignment.status`) and VS0-WRITER-022 (identity-membership-controller, paths
`Membership.spec`/`Membership.protected`/`Membership.status`). Two are consumed
unchanged: VS0-WRITER-007 (authorized-governance-publisher), whose registered
FEATURE-0018 paths are exactly `GovernanceProfile.spec` and `RoleDefinition.spec`;
and VS0-WRITER-011 (authorized-decision-or-audit-producer), whose registered paths
`DecisionRecord.record` and `AuditEvent.record` carry FEATURE-0018's FEATURE-0013
decision/audit publication. VS0-WRITER-023 (delegated-governance-admin,
`ProfileAssignment.spec`) is FEATURE-0020-owned and out of scope. No registered
writer exists for AccessGroup, ApprovalPolicy, ApprovalRequest,
PrivilegedAccessRequest, AccessReview, or ExceptionGrant, so those write paths are
`n/a — no registered writer ID` with the controlling F18-RD cited (each resource's
sole-writer authority is fixed by its F18-RD, not by a Slice-0 writer registration).

Error identifiers are distinguished from conformance IDs and domain outcomes. Each
Error IDs cell contains either an exact registered Problem/failure identifier when
one applies, or an explicit `n/a — <reason>` for domain decisions, success, or
intentionally unspecified inherited safe errors; a conformance ID, HTTP status, or
domain outcome is never presented as an error identifier. The two inherited security
failure mappings are one-to-one and are never associated with any local
`VS0-CF-F18-*` case: VS0-F01 maps only to VS0-CF-F01 (`AUTH_REQUIRED`/401, absent
authentication only, never invalid/expired assurance); VS0-F02 maps only to
VS0-CF-F02 (`RESOURCE_NOT_FOUND`/404 safe denial, never combined with
`AUTHORIZATION_DENIED`). A local case that registers a domain `Deny` is classified
as a terminal domain `Deny` outcome (not a Problem code); a local case that registers
an exact FEATURE-0012 top-level code (`AUTHORIZATION_DENIED`, `RESOURCE_NOT_FOUND`,
`CONFLICT`, `VALIDATION_FAILED`, `INTERNAL_ERROR`) is cited by that exact code; a
local case that registers a violation code (`REVIEWER_CONFLICT`,
`REVIEW_BENEFICIARY_UNRESOLVED`) is cited as a `violations[].code` value, not a
top-level Problem code; and no error is inferred beyond an exact registered
failure/conformance mapping. Canonical
owners are cited by reference; contract text from the Slice-0 registry (section 8.3),
the finalized data model, and the canonical contract catalog is not duplicated here.
Owner for every FEATURE-0018 requirement is FEATURE-0018, consistent with
`docs/phase2/PHASE2_FEATURE_SEQUENCE.md` order 8.

| REQ | Authority | Owner | Schema IDs | Writer IDs | State IDs | Error IDs | Conformance IDs |
|---|---|---|---|---|---|---|---|
| REQ-F18-01 | F18-RD-01; DEC-0060; ACR-2026-002 | FEATURE-0018 | n/a — no owned schema; scope/exclusion authority only (F18-RD-01) | n/a — no registered writer; defines no mutable path (F18-RD-01) | n/a — no lifecycle (F18-RD-01) | n/a — architecture-drift/structural-validation rejection, no registered top-level code (VS0-CF-F18-50) | VS0-CF-F18-50 |
| REQ-F18-02 | F18-RD-02; DEC-0026; ADH-2026-012; DEC-0037 | FEATURE-0018 | VS0-SCHEMA-020, VS0-SCHEMA-021, VS0-SCHEMA-022, VS0-SCHEMA-023, VS0-SCHEMA-024, VS0-SCHEMA-062, VS0-SCHEMA-063, VS0-SCHEMA-064, VS0-SCHEMA-065, VS0-SCHEMA-066, VS0-SCHEMA-067, VS0-SCHEMA-068 | VS0-WRITER-008, VS0-WRITER-022 (VS0-WRITER-007 and VS0-WRITER-011 consumed; VS0-WRITER-023 out of scope) | n/a — inventory/profile authority, no single lifecycle (per-contract lifecycles owned by F18-RD-04,05,07,08,12,13,14,16,18) | n/a — no registered failure/conformance error for the inventory rule itself; only FEATURE-0012 top-level codes apply where an operation errors (VS0-CF-F18-06) | VS0-CF-F18-06 |
| REQ-F18-03 | F18-RD-03 | FEATURE-0018 | VS0-SCHEMA-020, VS0-SCHEMA-068 | n/a — no registered writer; consumes an already-authenticated PrincipalRef (F18-RD-03) | n/a — operation-local value, no lifecycle (F18-RD-03) | `AUTHORIZATION_DENIED` (VS0-CF-F18-06 as `none or AUTHORIZATION_DENIED`, the sole-writer/system-identity enforcement outcome, which is distinct from and does not arise from REQ-F18-03 identity-stability semantics); the identity-stability local cases register no top-level Problem code: VS0-CF-F18-02 registers success (`null`); VS0-CF-F18-01 and VS0-CF-F18-05 register domain `Deny`/`none` by fixture. Absent authentication is the inherited VS0-F01 → VS0-CF-F01 (`AUTH_REQUIRED`/401) case only; invalid/expired/mismatched assurance fails closed under F18-RD-03/F18-RD-21 with no registered assurance-specific top-level code | VS0-CF-F18-01, VS0-CF-F18-02, VS0-CF-F18-05, VS0-CF-F18-06 |
| REQ-F18-04 | F18-RD-04; DEC-0018; DEC-0019; DEC-0020 | FEATURE-0018 | VS0-SCHEMA-062, VS0-SCHEMA-023 | VS0-WRITER-008 (RoleAssignment); n/a — no registered AccessGroup writer ID (AccessGroup sole-writer authority fixed by F18-RD-04) | n/a — no registered state-machine ID; AccessGroup lifecycle fixed by F18-RD-04 | n/a — no local case registers a top-level Problem code: VS0-CF-F18-01 registers domain `Deny` for the unrelated scope (with `Allow` for the requested scope); VS0-CF-F18-08 registers domain `Deny` for the raw-claim fixture (`none` when the canonical relationship applies). VS0-F02 → VS0-CF-F02 (`RESOURCE_NOT_FOUND`/404) is the separate inherited safe-denial case and is not either local case | VS0-CF-F18-01, VS0-CF-F18-08 |
| REQ-F18-05 | F18-RD-05 | FEATURE-0018 | VS0-SCHEMA-021 | VS0-WRITER-022 | n/a — no registered state-machine ID; Membership lifecycle and GuestExpired/Stale projections fixed by F18-RD-05 | `VALIDATION_FAILED`/422 (VS0-CF-F18-32, guest TimeBound violation); domain `Deny` (VS0-CF-F18-01 unrelated scope, VS0-CF-F18-03, VS0-CF-F18-04, VS0-CF-F18-08 raw claim, VS0-CF-F18-36, VS0-CF-F18-46); n/a — no exact top-level or violation identifier is registered: VS0-CF-F18-47 is an intentionally unspecified existing structural-validation or safe-denial error (ADH-2026-071 authorizes no new top-level code, violation code, HTTP status, or precedence rule); n/a — success or intentionally unspecified inherited safe error (VS0-CF-F18-07 `none or existing safe error`). VS0-F02 → VS0-CF-F02 (`RESOURCE_NOT_FOUND`/404) is the separate inherited safe-denial case and is not any local case here | VS0-CF-F18-01, VS0-CF-F18-03, VS0-CF-F18-04, VS0-CF-F18-07, VS0-CF-F18-08, VS0-CF-F18-32, VS0-CF-F18-36, VS0-CF-F18-46, VS0-CF-F18-47 |
| REQ-F18-06 | F18-RD-06; ADH-2026-058..066 | FEATURE-0018 | VS0-SCHEMA-068 | n/a — no registered writer; consumes FEATURE-0016 action registrations unchanged (F18-RD-06) | n/a — action composition, no lifecycle (F18-RD-06) | `AUTHORIZATION_DENIED` (VS0-CF-F18-11, qualifier-only retire); n/a — no exact top-level or violation identifier is registered: VS0-CF-F18-12 is an intentionally unspecified existing safe error, recorded granularity limitation (ADH-2026-071 authorizes no new top-level code, violation code, HTTP status, or precedence rule); n/a — success (VS0-CF-F18-10 `null`) | VS0-CF-F18-10, VS0-CF-F18-11, VS0-CF-F18-12 |
| REQ-F18-07 | F18-RD-07; DEC-0044 | FEATURE-0018 | VS0-SCHEMA-022 | VS0-WRITER-007 (consumed; registered path `RoleDefinition.spec`) | n/a — no registered state-machine ID; RoleDefinition lifecycle fixed by F18-RD-07 | domain `Deny` (VS0-CF-F18-34, during suspension) | VS0-CF-F18-34 |
| REQ-F18-08 | F18-RD-08; DEC-0018; DEC-0019 | FEATURE-0018 | VS0-SCHEMA-023 | VS0-WRITER-008 | n/a — no registered state-machine ID; RoleAssignment lifecycle and effective projection fixed by F18-RD-08 | `VALIDATION_FAILED`/422 (VS0-CF-F18-32, guest TimeBound violation); `AUTHORIZATION_DENIED` (VS0-CF-F18-44; VS0-CF-F18-09 and VS0-CF-F18-39 as `none or AUTHORIZATION_DENIED`); domain `Deny` (VS0-CF-F18-01 unrelated scope, VS0-CF-F18-03, VS0-CF-F18-04); n/a — success or intentionally unspecified inherited safe error (VS0-CF-F18-05, VS0-CF-F18-07, VS0-CF-F18-21, VS0-CF-F18-31 as `none or existing safe error`/`Deny`). VS0-F02 → VS0-CF-F02 (`RESOURCE_NOT_FOUND`/404) is the separate inherited safe-denial case and is not any local case here | VS0-CF-F18-01, VS0-CF-F18-03, VS0-CF-F18-04, VS0-CF-F18-05, VS0-CF-F18-07, VS0-CF-F18-09, VS0-CF-F18-21, VS0-CF-F18-31, VS0-CF-F18-32, VS0-CF-F18-39, VS0-CF-F18-44 |
| REQ-F18-09 | F18-RD-09 | FEATURE-0018 | VS0-SCHEMA-023, VS0-SCHEMA-068 | VS0-WRITER-008, VS0-WRITER-022 | n/a — authorization algebra is operation-local, no persisted lifecycle (F18-RD-09) | domain `Deny` (VS0-CF-F18-01 unrelated scope, VS0-CF-F18-03, VS0-CF-F18-25, VS0-CF-F18-33, VS0-CF-F18-35); `AUTHORIZATION_DENIED` (VS0-CF-F18-09, VS0-CF-F18-44, VS0-CF-F18-46); n/a — success or intentionally unspecified inherited safe error (VS0-CF-F18-05, VS0-CF-F18-08, VS0-CF-F18-29 as `none or Deny/existing safe error`). VS0-F02 → VS0-CF-F02 (`RESOURCE_NOT_FOUND`/404) is the separate inherited safe-denial case and is not any local case here | VS0-CF-F18-01, VS0-CF-F18-03, VS0-CF-F18-05, VS0-CF-F18-08, VS0-CF-F18-09, VS0-CF-F18-25, VS0-CF-F18-29, VS0-CF-F18-33, VS0-CF-F18-35, VS0-CF-F18-44, VS0-CF-F18-46 |
| REQ-F18-10 | F18-RD-10; DEC-0028; DEC-0036; ADH-2026-067..069 | FEATURE-0018 | VS0-SCHEMA-065, VS0-SCHEMA-068 | n/a — no registered writer; consumes FEATURE-0017 seam unchanged (F18-RD-10) | n/a — operation-local requirement, no lifecycle (F18-RD-10) | domain `Deny` (VS0-CF-F18-25, VS0-CF-F18-27, Indeterminate mapped to Deny); domain `RequiresApproval` (VS0-CF-F18-26, no implicit request) | VS0-CF-F18-25, VS0-CF-F18-26, VS0-CF-F18-27 |
| REQ-F18-11 | F18-RD-11; DEC-0028; DEC-0036 | FEATURE-0018 | VS0-SCHEMA-063 | n/a — no registered writer; bounded FEATURE-0017 consumption (F18-RD-11) | n/a — evaluation-seam bound, no lifecycle (F18-RD-11) | domain `RequiresApproval` (VS0-CF-F18-26, compatible pinned subject only) | VS0-CF-F18-26 |
| REQ-F18-12 | F18-RD-12; DEC-0018; DEC-0020 | FEATURE-0018 | VS0-SCHEMA-065, VS0-SCHEMA-068 | n/a — no registered ApprovalPolicy writer ID (ApprovalPolicy sole-writer authority fixed by F18-RD-12; not a VS0-WRITER-007 registered path) | n/a — no registered state-machine ID; ApprovalPolicy lifecycle fixed by F18-RD-12 | `AUTHORIZATION_DENIED` (VS0-CF-F18-14, VS0-CF-F18-15; VS0-CF-F18-41 as `none or AUTHORIZATION_DENIED`); n/a — no exact top-level or violation identifier is registered: VS0-CF-F18-47 and VS0-CF-F18-49 are intentionally unspecified existing structural-validation or safe-denial errors (ADH-2026-071 authorizes no new top-level code, violation code, HTTP status, or precedence rule); n/a — success (VS0-CF-F18-13 `null`) | VS0-CF-F18-13, VS0-CF-F18-14, VS0-CF-F18-15, VS0-CF-F18-41, VS0-CF-F18-47, VS0-CF-F18-49 |
| REQ-F18-13 | F18-RD-13 | FEATURE-0018 | VS0-SCHEMA-066 | VS0-WRITER-008 (sole publisher of any downstream RoleAssignment effect); n/a — no registered ApprovalRequest writer ID (ApprovalRequest sole-writer authority fixed by F18-RD-13) | n/a — no registered state-machine ID; ApprovalRequest lifecycle fixed by F18-RD-13 | `CONFLICT`/409 (VS0-CF-F18-16, losing terminal CAS; none for the winner); `AUTHORIZATION_DENIED` (VS0-CF-F18-15; VS0-CF-F18-41 as `none or AUTHORIZATION_DENIED`); n/a — success (VS0-CF-F18-13 `null`) | VS0-CF-F18-13, VS0-CF-F18-15, VS0-CF-F18-16, VS0-CF-F18-41 |
| REQ-F18-14 | F18-RD-14 | FEATURE-0018 | VS0-SCHEMA-063, VS0-SCHEMA-023 | VS0-WRITER-008 (sole publisher of the linked TimeBound RoleAssignment); n/a — no registered PrivilegedAccessRequest writer ID (sole-writer authority fixed by F18-RD-14) | n/a — no registered state-machine ID; PrivilegedAccessRequest lifecycle fixed by F18-RD-14 | `CONFLICT`/409 (VS0-CF-F18-17, safe replay or conflict); n/a — no exact top-level or violation identifier is registered: VS0-CF-F18-47 and VS0-CF-F18-49 are intentionally unspecified existing structural-validation or safe-denial errors (ADH-2026-071 authorizes no new top-level code, violation code, HTTP status, or precedence rule); n/a — success (VS0-CF-F18-13, VS0-CF-F18-40 `null`) | VS0-CF-F18-13, VS0-CF-F18-17, VS0-CF-F18-40, VS0-CF-F18-47, VS0-CF-F18-49 |
| REQ-F18-15 | F18-RD-15 | FEATURE-0018 | VS0-SCHEMA-063, VS0-SCHEMA-064 | VS0-WRITER-008 (sole publisher of the linked TimeBound RoleAssignment); n/a — no registered PrivilegedAccessRequest/AccessReview writer ID (sole-writer authority fixed by F18-RD-15/16) | n/a — no registered state-machine ID; break-glass path and readiness fixed by F18-RD-15 | n/a — success (VS0-CF-F18-18 `null`, successful break-glass path); an activation-blocking required-audit failure is the separate `INTERNAL_ERROR`/500 case registered by VS0-CF-F18-28 under F18-RD-20 | VS0-CF-F18-18 |
| REQ-F18-16 | F18-RD-16 | FEATURE-0018 | VS0-SCHEMA-064, VS0-SCHEMA-068 | VS0-WRITER-008 (sole publisher of RoleAssignment revocation/replacement), VS0-WRITER-022 (Membership disposition); n/a — no registered AccessReview writer ID (sole-writer authority fixed by F18-RD-16) | n/a — no registered state-machine ID; AccessReview lifecycle fixed by F18-RD-16 | `violations[].code` `REVIEWER_CONFLICT` or `REVIEW_BENEFICIARY_UNRESOLVED` (VS0-CF-F18-45, not top-level Problem codes); n/a — no exact top-level or violation identifier is registered: VS0-CF-F18-48 is an intentionally unspecified existing F18-RD-16 safe error, and VS0-CF-F18-47 and VS0-CF-F18-49 are intentionally unspecified existing structural-validation or safe-denial errors (ADH-2026-071 authorizes no new top-level code, violation code, HTTP status, or precedence rule); n/a — success or intentionally unspecified inherited safe error (VS0-CF-F18-18, VS0-CF-F18-19, VS0-CF-F18-20, VS0-CF-F18-31, VS0-CF-F18-37, VS0-CF-F18-42) | VS0-CF-F18-18, VS0-CF-F18-19, VS0-CF-F18-20, VS0-CF-F18-31, VS0-CF-F18-37, VS0-CF-F18-42, VS0-CF-F18-45, VS0-CF-F18-47, VS0-CF-F18-48, VS0-CF-F18-49 |
| REQ-F18-17 | F18-RD-17 | FEATURE-0018 | VS0-SCHEMA-067, VS0-SCHEMA-068 | n/a — no registered ExceptionGrant writer ID (ExceptionGrant is an ImmutableRecord whose sole-writer authority is fixed by F18-RD-17; not a VS0-WRITER-007 registered path) | n/a — no registered state-machine ID; ExceptionGrant is append-only with terminal domain `Grant\|Deny` per F18-RD-17 | terminal domain `Deny` (VS0-CF-F18-22, non-exceptionable control, using the registered exception-decision outcome; VS0-CF-F18-21 and VS0-CF-F18-43 as `none or (approved) terminal Deny`); `CONFLICT`/409 (VS0-CF-F18-24, overlap); n/a — success (VS0-CF-F18-23 `null`) | VS0-CF-F18-21, VS0-CF-F18-22, VS0-CF-F18-23, VS0-CF-F18-24, VS0-CF-F18-43 |
| REQ-F18-18 | F18-RD-18; DEC-0050; ADH-2026-033 | FEATURE-0018 | VS0-SCHEMA-024 | VS0-WRITER-007 (consumed; registered path `GovernanceProfile.spec`); VS0-WRITER-023 out of scope (FEATURE-0020 selection/resolution) | n/a — no registered state-machine ID; GovernanceProfile lifecycle fixed by F18-RD-18 | n/a — success or intentionally unspecified existing F18-RD-18 validation failure (VS0-CF-F18-51 `none or existing F18-RD-18 validation failure`); no new top-level Problem code is introduced | VS0-CF-F18-51 |
| REQ-F18-19 | F18-RD-19; DEC-0043; ADH-2026-017; DEC-0044; ADH-2026-073 | FEATURE-0018 | n/a — no owned schema; FEATURE-0013 envelopes reused by reference and profiles `authorization-decision/v1`, `approval-decision/v1`, `exception-decision/v1` registered through the FEATURE-0013 extension process (F18-RD-19) | VS0-WRITER-011 (consumed; FEATURE-0013 DecisionRecord/AuditEvent producer reused unchanged; no competing writer) | n/a — no registered state-machine ID; DecisionRecord persisted `FINAL` is inherited FEATURE-0013 behavior; AccessReview LRO fixed by F18-RD-16 | `INTERNAL_ERROR`/500 (VS0-CF-F18-52, required-append failure; n/a — success otherwise) | VS0-CF-F18-52 |
| REQ-F18-20 | F18-RD-20; DEC-0043; ADH-2026-017; ADH-2026-073 | FEATURE-0018 | n/a — no owned schema; FEATURE-0013 AuditEvent taxonomy reused by reference (F18-RD-20) | VS0-WRITER-011 (consumed; FEATURE-0013 AuditEvent producer reused unchanged); audit-obligation acceptance precedes VS0-WRITER-008/VS0-WRITER-022/VS0-WRITER-007 publication | n/a — audit-obligation boundary, no owned lifecycle (F18-RD-20) | `INTERNAL_ERROR`/500 (VS0-CF-F18-28, VS0-CF-F18-52, required audit failure publishes nothing); `AUTHORIZATION_DENIED` (VS0-CF-F18-06 as `none or AUTHORIZATION_DENIED`) | VS0-CF-F18-06, VS0-CF-F18-28, VS0-CF-F18-52 |
| REQ-F18-21 | F18-RD-21; DEC-0026; ADH-2026-012 | FEATURE-0018 | n/a — validation precedence over all owned schemas; no distinct schema (F18-RD-21) | n/a — precedence/publication ordering, not a registered writer (F18-RD-21) | n/a — no owned lifecycle (F18-RD-21) | `RESOURCE_NOT_FOUND`/404 safe denial (VS0-CF-F18-30); `CONFLICT`/409 (VS0-CF-F18-16 losing terminal CAS; VS0-CF-F18-17 safe replay or conflict); domain `Deny` (VS0-CF-F18-27, Indeterminate mapped to Deny) | VS0-CF-F18-16, VS0-CF-F18-17, VS0-CF-F18-27, VS0-CF-F18-30 |
| REQ-F18-22 | F18-RD-22 | FEATURE-0018 | n/a — deterministic fixtures over owned schemas; no distinct schema (F18-RD-22) | n/a — conformance harness, not a runtime writer (F18-RD-22) | n/a — no owned lifecycle (F18-RD-22) | n/a — fixture-validation failure is a conformance-gate outcome, not a registered runtime Problem code (VS0-CF-F18-50, 53) | VS0-CF-F18-01, VS0-CF-F18-02, VS0-CF-F18-03, VS0-CF-F18-04, VS0-CF-F18-05, VS0-CF-F18-06, VS0-CF-F18-07, VS0-CF-F18-08, VS0-CF-F18-09, VS0-CF-F18-10, VS0-CF-F18-11, VS0-CF-F18-12, VS0-CF-F18-13, VS0-CF-F18-14, VS0-CF-F18-15, VS0-CF-F18-16, VS0-CF-F18-17, VS0-CF-F18-18, VS0-CF-F18-19, VS0-CF-F18-20, VS0-CF-F18-21, VS0-CF-F18-22, VS0-CF-F18-23, VS0-CF-F18-24, VS0-CF-F18-25, VS0-CF-F18-26, VS0-CF-F18-27, VS0-CF-F18-28, VS0-CF-F18-29, VS0-CF-F18-30, VS0-CF-F18-31, VS0-CF-F18-32, VS0-CF-F18-33, VS0-CF-F18-34, VS0-CF-F18-35, VS0-CF-F18-36, VS0-CF-F18-37, VS0-CF-F18-39, VS0-CF-F18-40, VS0-CF-F18-41, VS0-CF-F18-42, VS0-CF-F18-43, VS0-CF-F18-44, VS0-CF-F18-45, VS0-CF-F18-46, VS0-CF-F18-47, VS0-CF-F18-48, VS0-CF-F18-49, VS0-CF-F18-50, VS0-CF-F18-51, VS0-CF-F18-52, VS0-CF-F18-53 |
| REQ-F18-23 | F18-RD-23 | FEATURE-0018 | n/a — standards-gate obligation; no owned schema (F18-RD-23) | n/a — architecture-gate obligation, not a registered writer (F18-RD-23) | n/a — no owned lifecycle (F18-RD-23) | n/a — conformance-gate failure, not a registered runtime Problem code (VS0-CF-F18-54) | VS0-CF-F18-54 |
| REQ-F18-24 | F18-RD-24; DEC-0035; DEC-0048 | FEATURE-0018 | n/a — projection over owned schemas; no distinct schema (F18-RD-24) | n/a — progressive-disclosure projection, not a registered writer (F18-RD-24) | n/a — no owned lifecycle; system-owned fields hidden (F18-RD-24) | n/a — no new error; reuses existing safe-denial/redaction behavior (VS0-CF-F18-29, VS0-CF-F18-39, VS0-CF-F18-40, VS0-CF-F18-41, VS0-CF-F18-42, VS0-CF-F18-43 register none or existing safe error) | VS0-CF-F18-29, VS0-CF-F18-39, VS0-CF-F18-40, VS0-CF-F18-41, VS0-CF-F18-42, VS0-CF-F18-43 |

## Exact conformance semantics ledger

Every referenced `VS0-CF-*` identifier is transcribed below with its exact
registry fields (ID, owner, inputs, expectedState, expectedError,
expectedSideEffects, gate) from the approved Slice-0 contract registry. A
conformance case proves only its registered scenario and is never used as a
substitute for a state-machine ID or reinterpreted as another feature's behavior.
These are the FEATURE-0018-owned registered Slice-0 conformance cases with exact
fields; the broader FEATURE-0018-local conformance inventory is governed by
REQ-F18-22 (F18-RD-22) and remains executable-conformance-pending.

| ID | owner | inputs | expectedState | expectedError | expectedSideEffects | gate |
|---|---|---|---|---|---|---|
| VS0-CF-F01 | FEATURE-0018 | unauthenticated-create | unchanged | AUTH_REQUIRED | none | security |
| VS0-CF-F02 | FEATURE-0018 | unauthorized-project-create | unchanged | RESOURCE_NOT_FOUND | none | security |
| VS0-CF-X01 | FEATURE-0018 | cross-organization-ref | unchanged | RESOURCE_NOT_FOUND | none | security |
| VS0-CF-X02 | FEATURE-0018 | cross-project-ref | unchanged | RESOURCE_NOT_FOUND | none | security |
| VS0-CF-F18-01 | FEATURE-0018 | Active employee, published Developer role, active project assignment; same action at unrelated project | requested project Allow; unrelated project unchanged | none; unrelated project Deny | exact protected authorization evidence and audit; no unrelated-scope effect | security |
| VS0-CF-F18-02 | FEATURE-0018 | Same issuer/subject principal with changed email display attribute | principal and existing authority remain bound to issuer/subject | null | attribution retains stable principal identity | identity |
| VS0-CF-F18-03 | FEATURE-0018 | Active assignment whose required Membership/principal becomes suspended | assignment retained but immediately ineffective | Deny | retained history and protected denial audit; no access effect | security |
| VS0-CF-F18-04 | FEATURE-0018 | Revoked/expired relationship followed by a later rejoin | new relationship does not reactivate old assignment | Deny until newly authorized grant | retained old evidence; no silent reactivation | lifecycle |
| VS0-CF-F18-05 | FEATURE-0018 | Workload principal with responsible Human and exact scoped assignment invokes one action | applicable action may Allow; unrelated action/scope denies | none or Deny by fixture | independent workload attribution and protected audit | security |
| VS0-CF-F18-06 | FEATURE-0018 | System controller submits an owned intent and separately attempts a non-owned write | owned controller transition may publish; non-owned write unchanged | none or AUTHORIZATION_DENIED | sole writer publishes; direct competing write has no effect | writer |
| VS0-CF-F18-07 | FEATURE-0018 | Consultant has distinct supplier and customer Guest Membership fixtures | only exact active target context is eligible until exact expiry | none or existing safe error | sponsorship, scope, expiry, and responsible Human remain attributable | security |
| VS0-CF-F18-08 | FEATURE-0018 | Raw external group claim; then trusted-provisioned AccessGroup, current direct Membership, and active scoped assignment | raw claim grants nothing; canonical relationship may allow | Deny; then none when all canonical evidence applies | no native/external group authority; canonical provenance audited | security |
| VS0-CF-F18-09 | FEATURE-0018 | Administrator delegates project action/reach within and beyond one independently applicable per-action witness | contained grant intent may publish; broader intent unchanged | none or AUTHORIZATION_DENIED | RoleAssignment controller is sole publisher; no witness synthesis | security |
| VS0-CF-F18-10 | FEATURE-0018 | Exact FEATURE-0016 ExecutionTarget read and qualify bindings for one accessible target | read/qualify apply only to registered target binding | null | protected action/target provenance; no provider-native effect | integration-contract |
| VS0-CF-F18-11 | FEATURE-0018 | Qualifier-only holder attempts ExecutionTarget retirement | target remains unchanged | AUTHORIZATION_DENIED | denial audit; no retirement | security |
| VS0-CF-F18-12 | FEATURE-0018 | Holder requests create-only ExecutionTarget authority under combined existing write action | no invented create-only grant/action is published | existing safe error | recorded granularity limitation; no action-registry mutation | compatibility |
| VS0-CF-F18-13 | FEATURE-0018 | Exact Human requests approved two-hour JIT production DB role before activation deadline | one TimeBound RoleAssignment becomes active then ineffective at exact expiry | null | immutable approval/activation provenance and required audit | privileged-access |
| VS0-CF-F18-14 | FEATURE-0018 | Requester attempts to approve own privileged request | request not approved or activated | AUTHORIZATION_DENIED | denied decision audited; no RoleAssignment intent/publication | security |
| VS0-CF-F18-15 | FEATURE-0018 | Listed approver loses current eligibility before decision acceptance | request remains without that accepted approval | AUTHORIZATION_DENIED | denial audited; no downstream effect | security |
| VS0-CF-F18-16 | FEATURE-0018 | Approval decision races ApprovalRequest expiry | exactly one valid terminal state wins; no unauthorized assignment | none for winner; CONFLICT for loser | one terminal evidence publication; no duplicate effect | concurrency |
| VS0-CF-F18-17 | FEATURE-0018 | Same approved JIT request activated twice with same and conflicting replay fixtures | exactly one temporary RoleAssignment; validity never extends | safe replay or CONFLICT | one activation publication and one completed result | idempotency |
| VS0-CF-F18-18 | FEATURE-0018 | Eligible Human invokes approved break-glass path with current assurance evidence | individual TimeBound access lasts no more than one hour | null | immediate protected audit and retrospective review due within 24 hours | privileged-access |
| VS0-CF-F18-19 | FEATURE-0018 | Standing certification finds stale administrator assignment and selects Revoke | exact assignment becomes revoked exactly once | null | retained immutable review, usage, decision, and remediation evidence | review |
| VS0-CF-F18-20 | FEATURE-0018 | Reviewed assignment version changes after immutable snapshot | review becomes Stale; assignment unchanged | null | no remediation or review-due advancement | concurrency |
| VS0-CF-F18-21 | FEATURE-0018 | Incident action requires both applicable DB role and bounded change-window exception | Allow only when both independent authorities are current | none or Deny | exact role and exception provenance; neither synthesizes the other | security |
| VS0-CF-F18-22 | FEATURE-0018 | Proposal attempts to except registered non-exceptionable tenant isolation control | no ExceptionGrant published | Deny using existing registered outcome | terminal protected denial evidence and audit | security |
| VS0-CF-F18-23 | FEATURE-0018 | Active ExceptionGrant receives valid linked early-revocation proposal | original record remains immutable and exact grant becomes ineffective at revoke instant | null | one linked effect=Revoke record and audit | exception |
| VS0-CF-F18-24 | FEATURE-0018 | Proposed ExceptionGrant overlaps same control/subject/scope/interval | existing grant unchanged; new grant absent | CONFLICT | no merge, precedence inference, or partial publication | exception |
| VS0-CF-F18-25 | FEATURE-0018 | FEATURE-0017 returns Allow while applicable RoleAssignment is expired | final authorization is Deny | Deny | protected decision provenance records both inputs; no effect | security |
| VS0-CF-F18-26 | FEATURE-0018 | Compatible pinned FEATURE-0017 evaluation returns RequiresApproval | no effect; request exists only through the explicit approved flow | RequiresApproval | protected evaluation/requirement provenance; no implicit ApprovalRequest | policy-seam |
| VS0-CF-F18-27 | FEATURE-0018 | FEATURE-0017 result is Indeterminate | final authorization is Deny | Deny | redacted protected decision/audit; no effect | security |
| VS0-CF-F18-28 | FEATURE-0018 | Required AuditEvent append fails at JIT activation publication boundary | request is not Active; no RoleAssignment or completed replay published | INTERNAL_ERROR | atomic rollback/no publication | audit |
| VS0-CF-F18-29 | FEATURE-0018 | Authorized end customer performs a routine low-risk action | business result returned without governance ceremony | null | identity/grants/guardrails/provenance/audit resolve internally | usability |
| VS0-CF-F18-30 | FEATURE-0018 | Denied operation references an inaccessible target | protected state unchanged | RESOURCE_NOT_FOUND safe denial | redacted audit; no existence or policy-detail disclosure | security |
| VS0-CF-F18-31 | FEATURE-0018 | Production Workload/System receives policy-permitted narrow Standing assignment | assignment effective, immediately revocable, and review-due | null | responsible Human, provenance, and review schedule retained | review |
| VS0-CF-F18-32 | FEATURE-0018 | Guest RoleAssignment omits notBefore or expiresAt | no assignment published | VALIDATION_FAILED | no audit-changing publication or completed replay | contract |
| VS0-CF-F18-33 | FEATURE-0018 | Two grants allow action while a mandatory applicable guardrail excludes it | final authorization is Deny | Deny | exact grant and guardrail provenance; no effect | security |
| VS0-CF-F18-34 | FEATURE-0018 | Published RoleDefinition version is suspended, then validly restored | pinned assignments ineffective during suspension; other versions unaffected | Deny during suspension | suspension/restoration DecisionRecord and audit; no assignment rewrite | lifecycle |
| VS0-CF-F18-35 | FEATURE-0018 | AuthorizationResult is replayed for a different principal/action/target/scope | result grants nothing and target state remains unchanged | Deny | fresh authorization/audit path; no bearer reuse | security |
| VS0-CF-F18-36 | FEATURE-0018 | Federated provenance is stale or responsible Human is absent at a required boundary | privilege-increasing operation unchanged | Deny | safe protected denial; raw claim cannot repair canonical evidence | security |
| VS0-CF-F18-37 | FEATURE-0018 | AccessReview usage telemetry is absent, incomplete, and complete in separate fixtures | absence is not non-use; no silent workload revocation | null | immutable coverage/usage evidence and explicit reviewer outcome | review |
| VS0-CF-F18-39 | FEATURE-0018 | Administrator supplies holder, role version, scope, and validity; separately forges system-owned fields | valid exact intent may publish; forged request unchanged | none or AUTHORIZATION_DENIED | RoleAssignment controller supplies refs/provenance/audit as sole publisher | usability |
| VS0-CF-F18-40 | FEATURE-0018 | Engineer submits exact time-bound privileged request | one requester-facing request status and expiry projection | null | internal request/evidence/resource separation retained | usability |
| VS0-CF-F18-41 | FEATURE-0018 | Exact eligible approver decides privileged/exception request; includes cross-stage reuse attempt | valid decision affects only current stage; one Human cannot count in two stages | none or AUTHORIZATION_DENIED | immutable decision evidence; downstream effect remains separate | approval |
| VS0-CF-F18-42 | FEATURE-0018 | Eligible independent reviewer selects Retain, Revoke, or Replace on exact snapshot | exact approved disposition or explicit stale outcome | none or existing safe error | usage summary, beneficiary set, decision, and remediation evidence retained | review |
| VS0-CF-F18-43 | FEATURE-0018 | User proposes bounded exceptionable control exception | one request status/effective-period projection; grant only after approved terminal processing | none or approved terminal Deny | proposal/approval/grant evidence remains internally distinct | usability |
| VS0-CF-F18-44 | FEATURE-0018 | ExactResource(A) grantor attempts resource B and ScopeOnly delegation | no broader RoleAssignment published | AUTHORIZATION_DENIED | no cross-assignment synthesis or partial grant | security |
| VS0-CF-F18-45 | FEATURE-0018 | Direct or indirect Standing beneficiary attempts Retain/Replace self-certification; unresolved expansion fixture included | no access-preserving/replacing effect or due advancement | REVIEWER_CONFLICT or REVIEW_BENEFICIARY_UNRESOLVED | independently authorized Revoke remains available and audited | security |
| VS0-CF-F18-46 | FEATURE-0018 | TimeBound membership administrator attempts to enable durable group-held assignment | no durable derived access; empty group addition grants nothing | AUTHORIZATION_DENIED | exact envelope/witness denial evidence; no Membership publication when expansion exceeds ceiling | security |
| VS0-CF-F18-47 | FEATURE-0018 | Membership administrator uses group/relationship/role evidence to manufacture approval, JIT, break-glass, or review eligibility | eligibility remains unsatisfied | existing structural-validation or safe-denial error | no ApprovalRequest decision, activation, review effect, or membership-derived eligibility | security |
| VS0-CF-F18-48 | FEATURE-0018 | Reviewer list is empty, indirect, stale, scope-incompatible, target-incompatible, or changes before remediation | review effect unchanged | existing F18-RD-16 safe error | eligibility and beneficiary conflict rechecked; no remediation/due advancement | security |
| VS0-CF-F18-49 | FEATURE-0018 | Role grantor or later rule publisher attempts to make a role/group/relationship qualify for eligibility | eligibility remains exact-Human-only and unsatisfied | existing structural-validation or safe-denial error | no approval, activation, Retain/Replace, or eligibility side effect | security |
| VS0-CF-F18-50 | FEATURE-0018 | Active FEATURE-0018 contract fixture containing one future-owned/prohibited resource, field, adapter, credential flow, or provider-native IAM object | fixture rejected before evaluation/publication | existing architecture-drift or structural-validation failure | zero external I/O and no canonical/runtime publication | architecture |
| VS0-CF-F18-51 | FEATURE-0018 | Valid exact GovernanceProfile v1 component/applicability fixture; then unknown, ambiguous, scope-incompatible, broadened, and retired-version fixtures | valid fixture accepted as pinned rule evidence only; invalid fixtures rejected | none or existing F18-RD-18 validation failure | no profile assignment/effective-resolution effect and no silent successor migration | contract |
| VS0-CF-F18-52 | FEATURE-0018 | Exact three DecisionProfile registrations and mandatory authorization/approval/exception decision and AuditEvent fixtures, including required-append failure | only approved profiles/taxonomy accepted; required publication is atomic with evidence | none or INTERNAL_ERROR for required-append failure | exact immutable FEATURE-0013 evidence; no competing envelope or fourth profile | audit |
| VS0-CF-F18-53 | FEATURE-0018 | Valid and invalid immutable synthetic fixture graphs with deterministic UTC clock and external-call counter | invalid graph rejected before evaluation; valid repeated evaluation is identical | none or existing fixture-validation failure | externalCallCount=0; no IdP, policy-engine, workflow-engine, CloudProvider, or network I/O | conformance |
| VS0-CF-F18-54 | FEATURE-0018 | Standards matrix with every applicable invariant mapped to local proof, reused dependency, or named exclusion; then one unmapped invariant | complete matrix passes; incomplete matrix fails final architecture approval gate | conformance-gate failure, not a runtime Problem code | no runtime state or conformity claim created | standards |

Permanent excluded tombstone `VS0-CF-F18-38` is not an active ledger
row. `VS0-CF-X03` remains FEATURE-0015-owned inherited evidence and is never
FEATURE-0018-local acceptance proof.

## Acceptance-to-conformance mapping

Every active AC maps exactly once to the local conformance ID with the same
numeric suffix. `AC-F18-38` is excluded and therefore has no active row here.

| Acceptance ID | Exact local conformance ID |
|---|---|
| AC-F18-01 | VS0-CF-F18-01 |
| AC-F18-02 | VS0-CF-F18-02 |
| AC-F18-03 | VS0-CF-F18-03 |
| AC-F18-04 | VS0-CF-F18-04 |
| AC-F18-05 | VS0-CF-F18-05 |
| AC-F18-06 | VS0-CF-F18-06 |
| AC-F18-07 | VS0-CF-F18-07 |
| AC-F18-08 | VS0-CF-F18-08 |
| AC-F18-09 | VS0-CF-F18-09 |
| AC-F18-10 | VS0-CF-F18-10 |
| AC-F18-11 | VS0-CF-F18-11 |
| AC-F18-12 | VS0-CF-F18-12 |
| AC-F18-13 | VS0-CF-F18-13 |
| AC-F18-14 | VS0-CF-F18-14 |
| AC-F18-15 | VS0-CF-F18-15 |
| AC-F18-16 | VS0-CF-F18-16 |
| AC-F18-17 | VS0-CF-F18-17 |
| AC-F18-18 | VS0-CF-F18-18 |
| AC-F18-19 | VS0-CF-F18-19 |
| AC-F18-20 | VS0-CF-F18-20 |
| AC-F18-21 | VS0-CF-F18-21 |
| AC-F18-22 | VS0-CF-F18-22 |
| AC-F18-23 | VS0-CF-F18-23 |
| AC-F18-24 | VS0-CF-F18-24 |
| AC-F18-25 | VS0-CF-F18-25 |
| AC-F18-26 | VS0-CF-F18-26 |
| AC-F18-27 | VS0-CF-F18-27 |
| AC-F18-28 | VS0-CF-F18-28 |
| AC-F18-29 | VS0-CF-F18-29 |
| AC-F18-30 | VS0-CF-F18-30 |
| AC-F18-31 | VS0-CF-F18-31 |
| AC-F18-32 | VS0-CF-F18-32 |
| AC-F18-33 | VS0-CF-F18-33 |
| AC-F18-34 | VS0-CF-F18-34 |
| AC-F18-35 | VS0-CF-F18-35 |
| AC-F18-36 | VS0-CF-F18-36 |
| AC-F18-37 | VS0-CF-F18-37 |
| AC-F18-39 | VS0-CF-F18-39 |
| AC-F18-40 | VS0-CF-F18-40 |
| AC-F18-41 | VS0-CF-F18-41 |
| AC-F18-42 | VS0-CF-F18-42 |
| AC-F18-43 | VS0-CF-F18-43 |
| AC-F18-44 | VS0-CF-F18-44 |
| AC-F18-45 | VS0-CF-F18-45 |
| AC-F18-46 | VS0-CF-F18-46 |
| AC-F18-47 | VS0-CF-F18-47 |
| AC-F18-48 | VS0-CF-F18-48 |
| AC-F18-49 | VS0-CF-F18-49 |

## Requirement-to-conformance mapping

| Requirement ID | Exact FEATURE-0018-local conformance IDs |
|---|---|
| REQ-F18-01 | VS0-CF-F18-50 |
| REQ-F18-02 | VS0-CF-F18-06 |
| REQ-F18-03 | VS0-CF-F18-01, VS0-CF-F18-02, VS0-CF-F18-05, VS0-CF-F18-06 |
| REQ-F18-04 | VS0-CF-F18-01, VS0-CF-F18-08 |
| REQ-F18-05 | VS0-CF-F18-01, VS0-CF-F18-03, VS0-CF-F18-04, VS0-CF-F18-07, VS0-CF-F18-08, VS0-CF-F18-32, VS0-CF-F18-36, VS0-CF-F18-46, VS0-CF-F18-47 |
| REQ-F18-06 | VS0-CF-F18-10, VS0-CF-F18-11, VS0-CF-F18-12 |
| REQ-F18-07 | VS0-CF-F18-34 |
| REQ-F18-08 | VS0-CF-F18-01, VS0-CF-F18-03, VS0-CF-F18-04, VS0-CF-F18-05, VS0-CF-F18-07, VS0-CF-F18-09, VS0-CF-F18-21, VS0-CF-F18-31, VS0-CF-F18-32, VS0-CF-F18-39, VS0-CF-F18-44 |
| REQ-F18-09 | VS0-CF-F18-01, VS0-CF-F18-03, VS0-CF-F18-05, VS0-CF-F18-08, VS0-CF-F18-09, VS0-CF-F18-25, VS0-CF-F18-29, VS0-CF-F18-33, VS0-CF-F18-35, VS0-CF-F18-44, VS0-CF-F18-46 |
| REQ-F18-10 | VS0-CF-F18-25, VS0-CF-F18-26, VS0-CF-F18-27 |
| REQ-F18-11 | VS0-CF-F18-26 |
| REQ-F18-12 | VS0-CF-F18-13, VS0-CF-F18-14, VS0-CF-F18-15, VS0-CF-F18-41, VS0-CF-F18-47, VS0-CF-F18-49 |
| REQ-F18-13 | VS0-CF-F18-13, VS0-CF-F18-15, VS0-CF-F18-16, VS0-CF-F18-41 |
| REQ-F18-14 | VS0-CF-F18-13, VS0-CF-F18-17, VS0-CF-F18-40, VS0-CF-F18-47, VS0-CF-F18-49 |
| REQ-F18-15 | VS0-CF-F18-18 |
| REQ-F18-16 | VS0-CF-F18-18, VS0-CF-F18-19, VS0-CF-F18-20, VS0-CF-F18-31, VS0-CF-F18-37, VS0-CF-F18-42, VS0-CF-F18-45, VS0-CF-F18-47, VS0-CF-F18-48, VS0-CF-F18-49 |
| REQ-F18-17 | VS0-CF-F18-21, VS0-CF-F18-22, VS0-CF-F18-23, VS0-CF-F18-24, VS0-CF-F18-43 |
| REQ-F18-18 | VS0-CF-F18-51 |
| REQ-F18-19 | VS0-CF-F18-52 |
| REQ-F18-20 | VS0-CF-F18-06, VS0-CF-F18-28, VS0-CF-F18-52 |
| REQ-F18-21 | VS0-CF-F18-16, VS0-CF-F18-17, VS0-CF-F18-27, VS0-CF-F18-30 |
| REQ-F18-22 | VS0-CF-F18-01, VS0-CF-F18-02, VS0-CF-F18-03, VS0-CF-F18-04, VS0-CF-F18-05, VS0-CF-F18-06, VS0-CF-F18-07, VS0-CF-F18-08, VS0-CF-F18-09, VS0-CF-F18-10, VS0-CF-F18-11, VS0-CF-F18-12, VS0-CF-F18-13, VS0-CF-F18-14, VS0-CF-F18-15, VS0-CF-F18-16, VS0-CF-F18-17, VS0-CF-F18-18, VS0-CF-F18-19, VS0-CF-F18-20, VS0-CF-F18-21, VS0-CF-F18-22, VS0-CF-F18-23, VS0-CF-F18-24, VS0-CF-F18-25, VS0-CF-F18-26, VS0-CF-F18-27, VS0-CF-F18-28, VS0-CF-F18-29, VS0-CF-F18-30, VS0-CF-F18-31, VS0-CF-F18-32, VS0-CF-F18-33, VS0-CF-F18-34, VS0-CF-F18-35, VS0-CF-F18-36, VS0-CF-F18-37, VS0-CF-F18-39, VS0-CF-F18-40, VS0-CF-F18-41, VS0-CF-F18-42, VS0-CF-F18-43, VS0-CF-F18-44, VS0-CF-F18-45, VS0-CF-F18-46, VS0-CF-F18-47, VS0-CF-F18-48, VS0-CF-F18-49, VS0-CF-F18-50, VS0-CF-F18-51, VS0-CF-F18-52, VS0-CF-F18-53 |
| REQ-F18-23 | VS0-CF-F18-54 |
| REQ-F18-24 | VS0-CF-F18-29, VS0-CF-F18-39, VS0-CF-F18-40, VS0-CF-F18-41, VS0-CF-F18-42, VS0-CF-F18-43 |

## 9. Design questions explicitly delegated by architecture

The architecture delegates the following implementation-mechanism questions to the
design stage. These are the only concerns the approved authority marks as design
questions; they carry no semantic or ownership latitude and may not alter any REQ,
AC, contract, writer, state, or error meaning.

| Delegated question | Delegating authority |
|---|---|
| Concrete fixture copying, construction, validation, and deterministic UTC time-supply mechanisms | REQ-F18-22 (F18-RD-22): "Design owns the concrete copying, construction, validation, and time-supply mechanisms." |
| Literal HTTP path design for the closed operation surface under FEATURE-0012 domain-grouped versioned route rules, without adding an operation, changing its authority, or converting a controller effect into caller permission | REQ-F18-02 (F18-RD-02) |

No other concern is delegated. Semantic or ownership questions are resolved only by
the controlling authority and, if unresolvable, halt with a stop condition (section 10).

## 10. Completeness and unresolved-decision report

### 10.1 Completeness

- All 24 approved requirement groups (REQ-F18-01..24) have exactly one normative
  detail heading in approved order (section 4.1) and are copied exactly once in the
  Canonical requirement ledger.
- All 49 approved acceptance rows (AC-F18-01..49) are copied exactly once in the
  Canonical acceptance ledger; all 48 active rows are mapped to controlling REQ(s)
  in section 4.2; AC-F18-38 is preserved as Excluded and enforced as a negative
  acceptance case (section 7.3).
- Every active FEATURE-0018-owned conformance row—four shared security cases and
  53 `VS0-CF-F18-*` local cases—is transcribed exactly from the registry. All 48
  active ACs and all 24 REQs have explicit local mappings; F18-38 remains excluded.
- Adjacent-feature exclusions (FEATURE-0019..0026) and feature-level non-goals are
  enforceable in section 7 with negative acceptance cases; prohibited/stale concepts
  are named only there.
- Inherited FEATURE-0012/0013/0016/0017 contracts are consumed by exact reference
  and are not redefined, copied, renamed, specialized, or transferred.

### 10.2 Boundaries observed

- No design choice (package, file, route, storage, algorithm, library, internal
  interface, or task decomposition) is introduced; the two design-delegated
  mechanism questions are recorded in section 9 exactly as delegated.
- No downstream-feature semantic, status field, metric, side effect, controller,
  adapter/plugin execution, or runtime proof is imported.
- No architecture decision is reopened, extended, merged, split, renumbered, or
  paraphrased; no 25th decision group is invented; no missing value is filled by
  inference.

### 10.3 Unresolved decisions

None unresolved. All 24 F18-RD decision groups are fully specified in the approved
package and were transcribable without inference. A prior
`REQUIREMENT_CLARIFICATION_REQUIRED` review condition — raised because the approved
F18-RD-08, F18-RD-16, and F18-RD-20 security/audit safeguards were omitted from the
requirements transcription — was resolved by faithful transcription (see section
10.4). A subsequent `REQUIREMENT_CLARIFICATION_REQUIRED` review condition — raised
because OPS-03 weakened the mandatory F18-RD-21 replay preconditions and because the
per-REQ schema/writer/state/error mapping was absent (section 8.3 listed only
aggregate IDs) — was resolved by qualifying OPS-03 to preserve all F18-RD-21 replay
preconditions and by adding the section 8.6 per-requirement traceability matrix. A
following `REQUIREMENT_CLARIFICATION_REQUIRED` review condition — raised because the
first section 8.6 matrix abbreviated or ranged conformance IDs, mixed controller
names and lifecycle prose into the writer/state columns, and introduced unsupported
or contradictory error mappings (a VS0-F01 assurance generalization, an ambiguous
combined VS0-F02 mapping, a REQ-F18-17 `VALIDATION_FAILED` mapping, and a
REQ-F18-17 VS0-WRITER-007 association) — was resolved by an earlier rebuild of
section 8.6 with fully qualified registered identifiers and justified
`n/a — <reason>` cells. A subsequent `REQUIREMENT_CLARIFICATION_REQUIRED` review
condition — raised because exact traceability remained inconsistent: (1) section 8.6
still associated VS0-F02 with local cases (VS0-CF-F18-01, VS0-CF-F18-08) in
REQ-F18-04/05/08/09 although the controlling mapping is one-to-one VS0-F02 →
VS0-CF-F02 and those local cases produce domain `Deny`, not `RESOURCE_NOT_FOUND`/404;
(2) the Error IDs column still mixed conformance IDs, domain outcomes, HTTP statuses,
and unspecified safe outcomes as though they were error identifiers; and (3) the
writer inventory was internally inconsistent because the section 8.6 preamble said
only three relevant registered writers exist and section 8.3 omitted VS0-WRITER-011
while REQ-F18-19/20 and sections 10.3/10.4 already treated VS0-WRITER-011 as a
consumed registered writer — was resolved by (a) removing every
VS0-F01/VS0-F02 association from local `VS0-CF-F18-*` rows and preserving VS0-F01 →
VS0-CF-F01 and VS0-F02 → VS0-CF-F02 as the only one-to-one inherited mappings;
(b) rebuilding each Error IDs cell to contain an exact registered Problem/failure
identifier, a terminal domain `Deny`/`RequiresApproval` outcome, a `violations[].code`
value, or an explicit `n/a — <reason>` for domain decisions, success, or intentionally
unspecified inherited safe errors; and (c) reconciling one registry-backed writer
inventory across sections 8.3, 8.6, 10.3, and 10.4 that classifies VS0-WRITER-008 and
VS0-WRITER-022 as owned, VS0-WRITER-007 and VS0-WRITER-011 as consumed (confirmed by
the manifest-pinned VS-000 contract registry), and VS0-WRITER-023 as out of scope. A prior
`REQUIREMENT_CLARIFICATION_REQUIRED` review condition — raised because several
registry placeholders had been converted into unapproved violation-channel choices and
because REQ-F18-03's Error IDs cell was internally contradictory: (1) the Error IDs
cells for REQ-F18-05, REQ-F18-06, REQ-F18-12, REQ-F18-14, and REQ-F18-16 classified the
intentionally unspecified existing safe errors VS0-CF-F18-12, VS0-CF-F18-47,
VS0-CF-F18-48, and VS0-CF-F18-49 as `violations[].code` values although ADH-2026-071
states those placeholders authorize no choice of top-level code, violation code, HTTP
status, or precedence; and (2) REQ-F18-03's Error IDs cell asserted that no mapped local
case registers a top-level Problem code while the same cell cited VS0-CF-F18-06 with
`AUTHORIZATION_DENIED` — was resolved by (a) replacing each inferred
`violations[].code` classification of VS0-CF-F18-12, VS0-CF-F18-47, VS0-CF-F18-48, and
VS0-CF-F18-49 with `n/a — no exact top-level or violation identifier is registered`,
retaining the exact ADH-2026-071 wording as an intentionally unspecified existing safe
error; (b) retaining `REVIEWER_CONFLICT` and `REVIEW_BENEFICIARY_UNRESOLVED`
(VS0-CF-F18-45) as `violations[].code` values because they are the only explicitly
registered violation codes; and (c) rewriting REQ-F18-03's Error IDs cell to state the
mapped VS0-CF-F18-06 `AUTHORIZATION_DENIED` outcome as the sole-writer/system-identity
enforcement outcome distinct from REQ-F18-03 identity-stability semantics, so the cell
no longer claims that no mapped case registers a top-level Problem code. The current
`REQUIREMENT_CLARIFICATION_REQUIRED` review condition — raised because the ADH-2026-073
transcription was incomplete and its traceability inconsistent: REQ-F18-19 omitted the
required `exceptionproposal.decided` linkage to the exact immutable proposal, terminal
reason, and containing operation/correlation evidence, and sections 8.1 and 8.6 omitted
ADH-2026-073 from the controlling authority for REQ-F18-19 and REQ-F18-20 — is resolved
in this revision by faithfully transcribing the complete ADH-2026-073 event linkage into
REQ-F18-19 through FEATURE-0013's existing envelope and linkage semantics (preserving
exactly-once emission, the Grant-only separate `exceptiongrant.issued` event, no grant or
issuance event for a terminal `Deny`, and submission-only meaning for
`exceptiongrant.proposed`) and by adding ADH-2026-073 to section 8.1 and to the
REQ-F18-19 and REQ-F18-20 authority cells in section 8.6. All
corrections are transcription/traceability corrections and required no new decision,
error choice, or inferred value; the only registered addition is the single approved
`exceptionproposal.decided` AuditEvent taxonomy type registered through FEATURE-0013's
existing extension mechanism. No `ARCHITECTURE_DECISION_REQUIRED`, unresolved
`REQUIREMENT_CLARIFICATION_REQUIRED`, `BOUNDARY_CHANGE_REQUIRED`,
`DEPENDENCY_APPROVAL_REQUIRED`, or `SECURITY_REVIEW_REQUIRED` condition was
triggered during this requirements stage.

### 10.4 Semantic-delta report (regenerated)

This revision makes the single approved ADH-2026-073 taxonomy correction in
REQ-F18-19 and its controlling-authority traceability in sections 8.1 and 8.6, plus
the corresponding delta record in section 10.4. It faithfully transcribes
ADH-2026-073 by registering exactly one additional FEATURE-0018 AuditEvent taxonomy
type, `exceptionproposal.decided`, through FEATURE-0013's existing extension
mechanism and by completing the required `exceptionproposal.decided` linkage to the
exact immutable proposal, its terminal reason, the mandatory `exception-decision/v1`
DecisionRecord, and the containing operation/correlation evidence using FEATURE-0013's
existing envelope and linkage semantics. It preserves OPS-03, the restored
F18-RD-08/16/20 safeguards, and every prior traceability-presentation correction from
the prior revisions unchanged, and preserves every REQ, AC, F18-RD semantic,
conformance registration, local mapping, exclusion, owner, and stage boundary. The
manifest-pinned target is unchanged and the ADH-2026-070 package hashes still match;
every delta is a faithful transcription/traceability correction only, introducing no
new resource, route, action, writer, state, error, decision profile, dependency,
persistence, external effect, design, or implementation choice beyond the one approved
`exceptionproposal.decided` taxonomy registration, and altering no REQ, AC, F18-RD
semantic, conformance registration, mapping, exclusion, owner, or stage boundary.

Current-revision delta (review-directed, ADH-2026-073):

| Location | Prior state | Corrected state | Controlling authority |
|---|---|---|---|
| REQ-F18-19 (F18-RD-19) `exceptionproposal.decided` linkage | The `exceptionproposal.decided` AuditEvent was described as linked only to its mandatory `exception-decision/v1` DecisionRecord | The event is transcribed with the complete ADH-2026-073 linkage — exact immutable proposal, terminal reason, mandatory `exception-decision/v1` DecisionRecord, and containing operation/correlation evidence — through FEATURE-0013's existing envelope and linkage semantics, preserving exactly-once emission, the Grant-only separate `exceptiongrant.issued` event, no grant or issuance event for a `Deny`, and submission-only meaning for `exceptiongrant.proposed`; no additional event, resource, writer, or error is added | ADH-2026-073 (Accepted 2026-09-03); F18-RD-17, F18-RD-19, F18-RD-20; FEATURE-0013 extension mechanism |
| Sections 8.1 and 8.6 (REQ-F18-19, REQ-F18-20) controlling authority | ADH-2026-073 was absent from the controlling handoff package in section 8.1 and from the REQ-F18-19/REQ-F18-20 authority cells in section 8.6 | ADH-2026-073 is added to the section 8.1 controlling handoff package and to the REQ-F18-19 and REQ-F18-20 authority cells in section 8.6 as the controlling authority for the `exceptionproposal.decided` registration | ADH-2026-073; ADH-2026-071 (traceability); VS-000 contract registry |

Prior-revision deltas (retained record): The immediately preceding revision made the
following review-directed traceability-presentation corrections; their text is
unchanged in this revision.

| Location | Prior state | Corrected state | Controlling authority |
|---|---|---|---|
| Section 8.6 rows REQ-F18-05, 06, 12, 14, 16 (placeholder violation-channel misclassification) | Each row classified the intentionally unspecified existing safe errors VS0-CF-F18-12 (REQ-F18-06), VS0-CF-F18-47 (REQ-F18-05/12/14/16), VS0-CF-F18-48 (REQ-F18-16), and VS0-CF-F18-49 (REQ-F18-12/14/16) as `violations[].code` values | Each of VS0-CF-F18-12, VS0-CF-F18-47, VS0-CF-F18-48, and VS0-CF-F18-49 is stated as `n/a — no exact top-level or violation identifier is registered`, retaining the exact ADH-2026-071 wording (existing safe error / existing structural-validation or safe-denial error / existing F18-RD-16 safe error) as an intentionally unspecified existing safe error; only the explicitly registered `REVIEWER_CONFLICT` and `REVIEW_BENEFICIARY_UNRESOLVED` (VS0-CF-F18-45) remain `violations[].code` values | ADH-2026-071 (placeholders authorize no top-level code, violation code, HTTP status, or precedence); VS-000 contract registry; .kiro/steering/slice0-contract.md §8–9 |
| Section 8.6 row REQ-F18-03 (Error IDs internal contradiction) | The Error IDs cell asserted that no mapped local case registers a top-level Problem code while the same cell cited VS0-CF-F18-06 with `AUTHORIZATION_DENIED` | The mapped VS0-CF-F18-06 `AUTHORIZATION_DENIED` outcome is stated as the sole-writer/system-identity enforcement outcome, distinct from and not arising from REQ-F18-03 identity-stability semantics; the identity-stability local cases (VS0-CF-F18-01/02/05) are separately stated as registering no top-level Problem code, so the cell is no longer internally contradictory | ADH-2026-071 (traceability); VS-000 contract registry; F18-RD-03; .kiro/steering/slice0-contract.md §8–9 |

Prior-revision deltas (retained record): The revisions preceding those made the
following three review-directed traceability-presentation corrections; their text is
unchanged in this revision.

| Location | Prior state | Corrected state | Controlling authority |
|---|---|---|---|
| Section 8.6 rows REQ-F18-04, 05, 08, 09 (VS0-F02 misassociation) | Each row cited `VS0-F02 → RESOURCE_NOT_FOUND/404 safe denial` against local cases VS0-CF-F18-01 (and VS0-CF-F18-08), implying VS0-F02 maps to those local cases | VS0-F02 is stated only as the separate inherited one-to-one mapping VS0-F02 → VS0-CF-F02 and is no longer associated with any local case; the cited local cases are classified by their exact registered outcome (VS0-CF-F18-01 and VS0-CF-F18-08 as domain `Deny`), and VS0-F01 likewise stays limited to VS0-CF-F01 | ADH-2026-071 (traceability); VS-000 contract registry (`VS0-F01→VS0-CF-F01`, `VS0-F02→VS0-CF-F02`); .kiro/steering/slice0-contract.md §9 |
| Section 8.6 Error IDs column (identifier vs outcome) | The column mixed conformance IDs, domain outcomes, HTTP statuses, and unspecified safe outcomes as though they were error identifiers | Each Error IDs cell contains an exact registered Problem/failure identifier (`AUTHORIZATION_DENIED`, `RESOURCE_NOT_FOUND`, `CONFLICT`, `VALIDATION_FAILED`, `INTERNAL_ERROR`), a terminal domain `Deny`/`RequiresApproval` outcome, a registered `violations[].code` value (`REVIEWER_CONFLICT`, `REVIEW_BENEFICIARY_UNRESOLVED`), or an explicit `n/a — <reason>` for domain decisions, success, or intentionally unspecified inherited safe errors; REQ-F18-17 uses terminal `Deny` plus `CONFLICT`/409 overlap, not `VALIDATION_FAILED`; no error is inferred | ADH-2026-071 (traceability); VS-000 contract registry; F18-RD-17; .kiro/steering/slice0-contract.md §8–9 |
| Sections 8.3 and 8.6 preamble (VS0-WRITER-011 inventory) | Section 8.3 omitted VS0-WRITER-011 and the section 8.6 preamble asserted only three relevant registered writers exist, contradicting REQ-F18-19/20 and sections 10.3/10.4 treating VS0-WRITER-011 as consumed | One registry-backed writer inventory is reconciled across sections 8.3, 8.6, 10.3, and 10.4: VS0-WRITER-008 and VS0-WRITER-022 owned; VS0-WRITER-007 (`GovernanceProfile.spec`/`RoleDefinition.spec`) and VS0-WRITER-011 (`DecisionRecord.record`/`AuditEvent.record`) consumed; VS0-WRITER-023 out of scope; the preamble now names exactly four relevant registered writers | VS-000 contract registry (manifest-pinned VS0-WRITER-007/008/011/022/023); ADH-2026-071 |

Prior-revision deltas (retained record): OPS-03 was qualified to preserve all
F18-RD-21 replay preconditions, and the section 8.6 per-requirement matrix was first
introduced; their controlling authority is REQ-F18-21 / F18-RD-21 and ADH-2026-071.
Additionally, the following three approved safeguards were restored from the
ADH-2026-070 package after having been omitted from the initial requirements
transcription; their text is unchanged in this revision.

| Location | Prior state | Restored approved safeguard | Controlling authority |
|---|---|---|---|
| REQ-F18-08 (F18-RD-08) | Responsible-party accountability stated without resolution/equality/fail-closed checks | System-selected `responsiblePartyRef` must resolve to a current Human at assignment creation, at replacement, and at every AccessReview retention; must equal the responsible party pinned on the exact authorizing Workload/System Membership where Membership applies; missing, mismatched, or non-Human responsibility evidence fails closed with no RoleAssignment intent published | REQ-F18-08 / F18-RD-08 |
| REQ-F18-16 (F18-RD-16) | Beneficiary expansion stated for access-preserving/replacing RoleAssignment dispositions only | Equivalent beneficiary-Human expansion restored for Membership-only `Retain`; retaining a Workload/System Membership or Workload/System-held RoleAssignment requires a current Human `responsiblePartyRef`; failed responsibility or beneficiary validation must not certify and must not silently revoke the item | REQ-F18-16 / F18-RD-16 |
| REQ-F18-20 (F18-RD-20) | Audit-obligation acceptance stated for authorization-changing publications; authorization-evaluation audit-failure consequence incomplete | Protected audit-obligation acceptance required before every accepted mutation represented by the taxonomy even when authorization is not yet changed (Draft create/update, proposal submission, request submission, synchronization, campaign creation); authorization-evaluation audit failure returns a safe internal error, produces no `AuthorizationResult`, and confers no downstream authority | REQ-F18-20 / F18-RD-20 |

No other requirement text changed. The prior `REQUIREMENT_CLARIFICATION_REQUIRED`
condition raised against the omission of F18-RD-08, F18-RD-16, and F18-RD-20
security/audit semantics remains resolved by faithful transcription. The prior
`REQUIREMENT_CLARIFICATION_REQUIRED` condition raised against the weakened OPS-03
replay preconditions and the absent per-REQ mapping remains resolved by the retained
OPS-03 qualification and the section 8.6 matrix. The earlier
`REQUIREMENT_CLARIFICATION_REQUIRED` condition raised against the first section 8.6
matrix's abbreviated conformance IDs and unsupported error/writer mappings remains
resolved by the fully qualified rebuild. The prior
`REQUIREMENT_CLARIFICATION_REQUIRED` condition — raised because exact traceability
remained inconsistent (a VS0-F02 misassociation with local cases in
REQ-F18-04/05/08/09, an Error IDs column that mixed conformance IDs, domain outcomes,
HTTP statuses, and unspecified safe outcomes as error identifiers, and a
VS0-WRITER-011 inventory inconsistency between the section 8.6 preamble, section 8.3,
and REQ-F18-19/20 with sections 10.3/10.4) — remains resolved by the sections 8.3, 8.6,
10.3, and 10.4 corrections recorded above: VS0-F01/VS0-F02 preserve their one-to-one
VS0-CF-F01/VS0-CF-F02 mappings and are absent from local rows; every Error IDs cell
now holds an exact registered identifier or a justified `n/a` distinguishing domain,
violation, success, and inherited-safe outcomes; and one registry-backed writer
inventory classifies VS0-WRITER-008/022 owned, VS0-WRITER-007/011 consumed, and
VS0-WRITER-023 out of scope. A prior `REQUIREMENT_CLARIFICATION_REQUIRED`
condition — raised because several registry placeholders had been converted into
unapproved violation-channel choices (VS0-CF-F18-12, VS0-CF-F18-47, VS0-CF-F18-48,
and VS0-CF-F18-49 classified as `violations[].code` in the section 8.6 Error IDs
cells for REQ-F18-05, 06, 12, 14, and 16) and because REQ-F18-03's Error IDs cell was
internally contradictory (claiming no mapped local case registers a top-level Problem
code while citing VS0-CF-F18-06 with `AUTHORIZATION_DENIED`) — is resolved by the
sections 8.6, 10.3, and 10.4 corrections above: each of VS0-CF-F18-12, VS0-CF-F18-47,
VS0-CF-F18-48, and VS0-CF-F18-49 is now `n/a — no exact top-level or violation
identifier is registered`, retaining the exact ADH-2026-071 wording as an
intentionally unspecified existing safe error; only the explicitly registered
`REVIEWER_CONFLICT` and `REVIEW_BENEFICIARY_UNRESOLVED` (VS0-CF-F18-45) remain
`violations[].code` values; and REQ-F18-03's Error IDs cell now states the mapped
VS0-CF-F18-06 `AUTHORIZATION_DENIED` outcome as the sole-writer/system-identity
enforcement outcome distinct from REQ-F18-03 identity-stability semantics. The current
`REQUIREMENT_CLARIFICATION_REQUIRED` condition — raised because the ADH-2026-073
transcription was incomplete (REQ-F18-19 omitted the required `exceptionproposal.decided`
linkage to the exact immutable proposal, terminal reason, and containing
operation/correlation evidence) and its traceability was inconsistent (sections 8.1 and
8.6 omitted ADH-2026-073 from the controlling authority for REQ-F18-19 and REQ-F18-20) —
is resolved by the REQ-F18-19, section 8.1, section 8.6, and section 10.4 corrections
recorded above: REQ-F18-19 now transcribes the complete ADH-2026-073 event linkage
through FEATURE-0013's existing envelope and linkage semantics (preserving exactly-once
emission, the Grant-only separate `exceptiongrant.issued` event, no grant or issuance
event for a terminal `Deny`, and submission-only meaning for `exceptiongrant.proposed`),
and ADH-2026-073 is added as controlling authority in section 8.1 and in the REQ-F18-19
and REQ-F18-20 authority cells of section 8.6. The only registered addition is the
single approved `exceptionproposal.decided` AuditEvent taxonomy type registered through
FEATURE-0013's existing extension mechanism. No new resource, route, action, writer,
state, error, decision profile, dependency, design, or error choice was introduced.
Requirements are resubmitted for review.

### 10.5 Notes

- `VS0-CF-X03` remains FEATURE-0015-owned inherited safe-denial evidence. It is
  intentionally absent from the FEATURE-0018 exact local ledger and never counts
  as FEATURE-0018 acceptance proof.

## Stage receipt

STAGE_STATUS: COMPLETE
