# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-070
- Date: 2026-08-28
- Source discussion: Codex FEATURE-0018 architecture session; architecture-owner
  approved FEATURE-0018 architecture digest and industry validation matrix
- Related feature: FEATURE-0018
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

FEATURE-0018 governance, IAM, approval, access-review, privileged-access, and
exception foundation

## Summary

FEATURE-0018 establishes Sovrunn's implementation-neutral enterprise authorization
foundation for Human, Workload, and System principals. It defines stable
principal use, non-authorizing Membership, direct and canonical AccessGroup
role holders, membership-instance-bound direct grants, immutable versioned
roles, scoped and explicitly valid role
assignments, contextual membership eligibility, governed versioned scope
applicability, monotonic privilege classification, deterministic zero-trust
authorization composition, explicit operation-local approval requirements,
bounded approval and privileged-access flows,
snapshot-based access review, immutable exception evidence, a FEATURE-0018-
limited GovernanceProfile envelope, FEATURE-0013 decision/audit adoption,
deterministic in-memory fixtures, and FEATURE-0018-local conformance.

The renewed package also closes the first independent security review's five
blocking findings: RoleAssignment has one canonical writer; delegated grants
cannot exceed per-action exact resource reach; access beneficiaries cannot
self-certify retained or replacement authority; normal JIT and break-glass
share a phishing-resistant AAL2-or-higher activation floor; and all readiness
mirrors must match the normative F18-RD vocabulary and boundary.

The renewal-02 corrections additionally treat every AccessGroup Membership
assignment-effect expansion as an atomic grant-producing boundary subject to the
same complete action/scope/target/resource/validity dominance as direct grants,
without adding a resource or workflow. They also restore the already-accepted
DEC-0059 supersession of DEC-0058 in active traceability and require a
fail-closed status-consistency drift check.

The renewal-03 corrections close group-derived eligibility and reviewer
evaluation semantics. Renewal-04 then identified RoleDefinition qualification
as a remaining hidden delegation channel. The v1 closure therefore makes
approver, privileged-requester, and reviewer eligibility one non-resource
`EligibilityRef` containing exactly one Human PrincipalRef. RoleDefinition,
RoleAssignment, AccessGroupRef, Membership, claim, and other relationships
cannot satisfy or create eligibility. Lists remain non-empty, OR-composed,
independently action-authorized, re-evaluated at decision/effect publication,
and fail-closed; policy/rule applicability plus the separate canonical action
alone determine scope and target reach.

The feature does not implement an identity provider, external group sync,
policy engine, workflow engine, provider-native IAM, effective-governance
resolution, infrastructure execution, production persistence, or a new
customer-facing facade resource. Security internals remain normally hidden
behind six progressive-disclosure user journeys.

## Classification

New decision.

This is the primary classification required by the handoff template because
the decision adds canonical `AccessGroup` and `AccessGroupRef` contracts and
changes the canonical `RoleAssignment` holder model. The same package includes
compatible extensions and corrections to existing IAM, governance,
decision-profile, lifecycle, and catalog text. Those subordinate changes do
not reduce the required new-decision change-control path.

An Architecture Change Request and a new accepted Decision Record are
required. Kiro must validate whether an RFC amendment is also required. The
canonical model, contract catalog, glossary, traceability, and current
architecture baseline must be reconciled atomically before the decision is
treated as applied.

## Existing approved baseline

- `ARCH-2026.08-PHASE2R-CANONICAL` is the active approved baseline.
- Phase 2R permits a deterministic, implementation-neutral Governance, IAM,
  Approval and Exception Foundation with no real external integration or
  infrastructure effect.
- The canonical FEATURE-0018 inventory currently contains
  `GovernanceProfile`, `Membership`, `RoleDefinition`, `RoleAssignment`,
  `PrivilegedAccessRequest`, `AccessReview`, `ApprovalPolicy`,
  `ApprovalRequest`, and `ExceptionGrant`. It does not contain `AccessGroup`.
- `PrincipalRef` is a stable reference to an authenticated Human, Workload, or
  System identity using issuer and subject; email is not durable identity.
- Membership records Organization or CloudProvider belonging and grants no
  permission by itself.
- `RoleDefinition` is a versioned definition containing canonical Sovrunn
  actions. Published versions are immutable under DEC-0044.
- `RoleAssignment` is the canonical scoped role grant. The current canonical
  text binds one principal and permits bounded conditions; it does not define
  the approved role-holder union, exact resource narrowing, or explicit
  Standing/TimeBound modes.
- `PrivilegedAccessRequest` and `AccessReview` are long-running operations.
- `GovernanceProfile` is the compositional governance envelope established by
  ADH-2026-033. FEATURE-0020 owns `ProfileAssignment`, hierarchy traversal,
  conflict resolution, exception application, and immutable
  `EffectiveGovernanceContext` construction.
- FEATURE-0012 owns common metadata, the seven-scope vocabulary,
  `metadata.scopeRef` as sole scope authority, typed references, validation,
  safe non-disclosure, redaction, optimistic concurrency, idempotency, and
  error behavior.
- FEATURE-0013 owns `DecisionRecord`, `DecisionProfile`, `EvaluationResult`,
  and `AuditEvent`; adopting domains register profiles rather than creating a
  competing decision or audit envelope.
- FEATURE-0016 owns ExecutionTarget actions and lifecycle. Its accepted action
  registry includes `executiontarget.read`, `executiontarget.qualify`, and
  `executiontarget.write`; create and retire remain combined under
  `executiontarget.write`.
- FEATURE-0017 owns the generic `PolicyEvaluationRequest` and
  `PolicyEvaluationResult`, `PolicyEngineAdapter`, deterministic digest-keyed
  fake, normalized outcomes, and pure FEATURE-0013 mapper. `Allow` is not
  execution authority and `RequiresApproval` does not create a request.
- DEC-0026 requires reuse before build. DEC-0028 prohibits embedding custom
  policy-engine behavior in handlers or domain code. DEC-0036 requires adapter
  boundaries before replaceable external integration.
- DEC-0037 defines the seven canonical scope kinds: Platform, Organization,
  OrganizationUnit, Tenant, Project, CloudPlatform, and CloudProvider.
- DEC-0050 and ADH-2026-033 require one resolved
  `EffectiveGovernanceContext`; FEATURE-0018 must not become a second resolver.

Relevant authority:

- `docs/context/ARCHITECTURE_VERSION.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/context/SOVRUNN_CONTEXT_PACK.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
- `docs/decisions/DECISION_INDEX.md`
- `docs/architecture/canonical/sovrunn-finalized-data-model.md`
- `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`
- `docs/glossary.md`
- `docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md`
- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`
- `docs/architecture/policy-evaluation-abstraction.md`
- `docs/reviews/architecture-decision-handoffs/canonical-model/ADH-2026-033-governance-context-profile-consolidation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-067-feature-0017-policy-evaluation-architecture-closure.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-068-feature-0017-structural-mapper-clarification.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-069-feature-0017-timing-transport-and-mapper-field-closure.md`
- `docs/reviews/architecture-readiness/FEATURE-0018-architecture-digest.md`
- `docs/reviews/architecture-readiness/FEATURE-0018-industry-validation-matrix.md`

## Approval-package authority

This core ADH and its two normative appendices form one indivisible approval
package:

- core: `docs/reviews/architecture-decision-handoffs/ADH-2026-070-feature-0018-governance-iam-approval-exception-foundation.md`;
- Appendix A: `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md`; and
- Appendix B: `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md`.

The core carries change control, rationale, reconciliation, impacted artifacts,
acceptance, Kiro instructions, and human-approval metadata. [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) carries the
domain semantic decision groups. [Appendix B](ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) carries the contract/action/governance
registries and the decision/audit, conformance, and standards-validation groups.
Neither appendix is independently approvable or applicable, and neither creates
architecture authority outside ADH-2026-070.

Human approval must cover the exact bytes of all three files as one unit.
Approval becomes effective only when external approval evidence pins one Git
commit containing all three files or records an independent SHA-256 digest for
each file. After approval, any byte change to any package file invalidates the
package approval and returns the complete package to `Proposed` until it is
reviewed and approved again. Kiro must validate and apply all three files
together and must stop if any file is absent, changed, independently approved,
or inconsistent. No package file has precedence over another; inconsistency is
a blocker requiring correction and renewed approval.

## Decision or proposed decision

F18-RD-01 through F18-RD-24 are the complete FEATURE-0018 architecture closure
carried by this package. Their identifiers, titles, and approved semantics are
preserved exactly in the two appendices. Kiro must not add, merge, split,
renumber, paraphrase, or infer another decision group.

The F18-RD groups are the sole normative FEATURE-0018 architecture decision
statements added by this package. The core rationale, conflict, required-action,
impacted-file, acceptance, and instruction sections may impose mandatory
application or process checks, but architecture summaries in those sections
are non-authoritative restatements. They must reference, not redefine, F18-RD
behavior.

Reader map:

| Area | Authoritative decisions | Normative location |
|---|---|---|
| Authority and inventory | F18-RD-01 through F18-RD-02 | Appendices A and B as indexed below |
| Identity and Membership | F18-RD-03 through F18-RD-05 | Appendix A |
| Roles and authorization | F18-RD-06 through F18-RD-11 | Appendices A and B as indexed below |
| Approval and privileged access | F18-RD-12 through F18-RD-15 | Appendix A |
| Review and exception | F18-RD-16 through F18-RD-17 | Appendix A |
| Governance and evidence | F18-RD-18 through F18-RD-23 | Appendices A and B as indexed below |
| User simplicity | F18-RD-24 | Appendix A |

Exact decision index:

| Decision group | Title | Normative location |
|---|---|---|
| F18-RD-01 | Current-feature-only authority | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-02 | Closed contract inventory and profiles | [Appendix B](ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-03 | Stable principal identity | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-04 | Direct principal and scoped AccessGroup assignments | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-05 | Membership is non-authorizing | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-06 | Distributed action ownership and central role composition | [Appendix B](ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-07 | Versioned RoleDefinition | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-08 | Scoped RoleAssignment | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-09 | Deterministic scoped authorization composition | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-10 | FEATURE-0017 adoption without reinterpretation | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-11 | Bounded FEATURE-0017 subject/target use | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-12 | Bounded ApprovalPolicy | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-13 | Immutable terminal approval evidence | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-14 | JIT privileged access authorizes one temporary grant | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-15 | Constrained break-glass access | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-16 | Snapshot-based AccessReview | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-17 | Bounded immutable exception evidence | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-18 | FEATURE-0018-limited GovernanceProfile v1 | [Appendix B](ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-19 | FEATURE-0013 adoption | [Appendix B](ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-20 | Audit before authorization-changing publication | [Appendix B](ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-21 | Deterministic validation and safe denial | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-22 | Deterministic in-memory foundation and local conformance | [Appendix B](ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-23 | Standards-validation gate | [Appendix B](ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-24 | Progressive and normally hidden user friction | [Appendix A](ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |

## Rationale

The model separates independent identity, definition, grant, request,
approval, review, and retained-exception lifecycles while composing them into
simple risk-proportional journeys. This avoids both a single overloaded IAM
object and unnecessary user-visible resources.

Direct-only assignment is insufficient for enterprise operation; canonical
AccessGroup provides auditable bulk administration without trusting external
claims. A role-holder union is clearer than overloading `PrincipalRef`, because
an AccessGroup cannot authenticate. Direct membership, no nesting, and no
dynamic evaluation keep v1 deterministic and prevent transitive privilege
surprises.

Grant-union plus mandatory-guardrail intersection matches mature cloud IAM
practice while preserving Sovrunn's deny-by-default, least-privilege,
resource-bound, current-evidence, and non-bearer-result invariants. Explicit
Standing and TimeBound validity avoids accidental workload outages while
requiring stronger controls for human, guest, JIT, and emergency access.

Separate bounded approval, review, and immutable exception evidence preserves
separation of duties, effective-once behavior, audit integrity, and independent
retention. Keeping GovernanceProfile v1 limited to FEATURE-0018 semantics
prevents sovereignty, placement, cost, backup, entitlement, and effective-
resolution leakage from future owning features.

Alternatives rejected by this decision include:

- direct-principal-only authorization with no enterprise grouping;
- authorization directly from OIDC/external group claims;
- nested or dynamic AccessGroups in v1;
- treating Membership as a permission grant;
- using persona, job title, email, or provider-native IAM role as authority;
- mutable published roles or silent assignment upgrade to a new role version;
- a generic user-authored RoleAssignment condition/ABAC language;
- policy Allow, approval, or exception as an independent permission grant;
- unbounded workflow semantics or a new workflow resource family;
- mutable review usage fields on RoleAssignment;
- an additional ExceptionRequest or break-glass resource;
- FEATURE-0018 ownership of effective governance resolution; and
- premature real IdP, policy, workflow, provider, or persistence integration.

## Reuse-before-build assessment

This is the package-level FEATURE-0011 summary. The authoritative FEATURE-0018
architecture must carry the complete feature-level summary and capability
assessments required by
`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` before requirements are
authorized.

| Capability / decision unit | Disposition | Decision status | Rationale | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 metadata, scope, reference, validation, error, redaction, safe-denial, concurrency, and idempotency foundations | Reuse | Proposed | Existing canonical common contracts satisfy the shared API/resource behavior and must not be duplicated. | DEC-0026; ADH-2026-012 |
| FEATURE-0013 DecisionRecord, DecisionProfile, EvaluationResult, and AuditEvent envelopes | Reuse | Proposed | FEATURE-0013 remains the single evidence-envelope owner; FEATURE-0018 registers profiles and taxonomy only. | DEC-0043; ADH-2026-017 |
| FEATURE-0016 ExecutionTarget actions and lifecycle authority | Reuse | Proposed | FEATURE-0018 composes registered actions without splitting or reinterpreting them. | ADH-2026-058 through ADH-2026-066 |
| FEATURE-0017 policy-evaluation seam, result vocabulary, mapper, and deterministic fake | Reuse | Proposed | The existing engine-neutral boundary is consumed unchanged within its v1 subject/target limitation. | DEC-0028; DEC-0036; ADH-2026-067 through ADH-2026-069 |
| NIST, ISO, CSA, CIS, OWASP, OIDC, OAuth, SCIM, WebAuthn, and SPIFFE standards | Reuse | Proposed | Standards constrain invariants, assurance, proof, and future adapter compatibility without becoming vendor-native canonical schemas. | F18-RD-23 and FEATURE-0018 industry validation matrix |
| Sovrunn canonical IAM and authorization algebra | Build | Proposed | No external product owns Sovrunn's seven-scope actions, role-holder model, grant/guardrail composition, evidence linkage, or implementation-neutral contracts. | This package, F18-RD-02 through F18-RD-11 |
| Sovrunn approval, JIT, break-glass, review, and exception semantics plus deterministic fakes | Build | Proposed | These are Sovrunn control-plane contracts and conformance machinery; a workflow product cannot own their canonical meaning. | This package, F18-RD-12 through F18-RD-17 and F18-RD-22 |
| FEATURE-0018 GovernanceProfile v1 composition and progressive user projections | Build | Proposed | Sovrunn owns its governance envelope and user experience while FEATURE-0020 retains resolution. | ADH-2026-033; this package, F18-RD-18 and F18-RD-24 |
| Real identity-provider/federation/group-provisioning integration | Wrap | Deferred | Wrap is a non-authoritative candidate only. The future owning feature must perform a fresh FEATURE-0011 assessment and may select any controlled disposition supported by then-current evidence. | DEC-0022; DEC-0036 |
| Real policy-engine adapter | Wrap | Deferred | Non-authoritative compatibility candidate only. FEATURE-0017 selects no engine; the future owner performs a fresh FEATURE-0011 assessment and may select any supported disposition without a Wrap presumption. | DEC-0028; DEC-0036 |
| Real workflow-engine adapter | Wrap | Deferred | Non-authoritative compatibility candidate only. FEATURE-0018 selects no runtime; the future owner performs a fresh FEATURE-0011 assessment and may select any supported disposition without a Wrap presumption. | DEC-0026; DEC-0036 |

Every `Wrap`/`Deferred` row above is a non-authoritative compatibility
hypothesis. It authorizes no product, protocol, adapter design, interface
expansion, mapping schema, credential flow, deployment boundary, requirement,
task, or implementation. The future owning feature must repeat FEATURE-0011
using then-current requirements and evidence and may select `Reuse`, `Wrap`,
`Extend`, or `Build`.

For Build units, Reuse is insufficient because no mature component owns the
canonical Sovrunn resource identities, scopes, actions, lifecycles, evidence,
and progressive projections together. Wrap is insufficient because Phase 2R
selects no external identity, policy, or workflow runtime. Extend is
insufficient because extending FEATURE-0012, FEATURE-0013, FEATURE-0016, or
FEATURE-0017 would violate their existing ownership. Sovrunn builds only the
implementation-neutral contracts, deterministic composition, fakes, and local proof;
it does not build an IdP, policy engine, or workflow engine.

## Phase impact

- Current phase allowed: Yes.
- FEATURE-0018 remains order 8 in Phase 2R and depends on FEATURE-0012 and
  FEATURE-0017, while explicitly adopting FEATURE-0013 and consuming
  FEATURE-0016 actions.
- Phase 2R work remains deterministic, in-memory, implementation-neutral,
  adapter-first, decision-first, audit-first, and side-effect-free.
- No real IdP, OIDC/SCIM adapter, policy engine, workflow engine, provider IAM,
  database, Kubernetes, provisioning, or execution is introduced.
- No FEATURE-0019 or later implementation is pulled forward.
- FEATURE-0020 continues to depend on FEATURE-0018 and remains the sole owner
  of ProfileAssignment and EffectiveGovernanceContext resolution.
- Backward compatibility impact is documentation/canonical-contract only:
  FEATURE-0018 has no implemented runtime objects requiring data migration.
  Canonical bootstrap remains intact and no dual authority is permitted.
- The current baseline requires an approved canonical-model amendment because
  `AccessGroup`, `roleHolderRef`, RoleAssignment restriction/validity,
  RoleDefinition lifecycle, GovernanceProfile v1 composition, and new decision
  profiles alter or complete canonical contract authority.

## Conflict check

- Conflicts with accepted DEC/RFC: Yes, at canonical contract wording and
  inventory level; no conflict requires reversing the product principles.
- Primary resolution: ACR plus new DEC, canonical model/catalog/glossary
  update, traceability update, and current baseline reconciliation.
- RFC action: Kiro must determine whether the IAM/governance RFC set can be
  amended or a focused FEATURE-0018 RFC is required. No RFC semantics may be
  invented beyond this package.

Canonical and dependent-artifact conflicts to reconcile:

| ID | Affected authority or artifact | Conflict to reconcile | Controlling authority and required action |
|---|---|---|---|
| 1 | Canonical inventory | Omits AccessGroup, AccessGroupRef, supporting values, enums, and group suspension/retirement effects | Apply F18-RD-02 and F18-RD-04. |
| 2 | Canonical Membership | Supports belonging only and lacks same-scope direct AccessGroup relationships | Apply F18-RD-04 and F18-RD-05 without making AccessGroup a ScopeKind. |
| 3 | Canonical RoleAssignment | Binds principalRef and omits the holder union, required Membership/responsibility pins, Standing review-rule pin, and group compatibility | Apply F18-RD-04 and F18-RD-08. |
| 4 | Canonical RoleAssignment restrictions | Retains conditions and competing target/scope examples | Apply F18-RD-02, F18-RD-06, F18-RD-08, and F18-RD-09. |
| 5 | Canonical validity | Omits Standing/TimeBound and exact no-grace expiry | Apply F18-RD-08. |
| 6 | Canonical definition lifecycle | Conflates retirement and supersession and lacks emergency RoleDefinition suspension | Apply F18-RD-07 and the referenced definition-specific lifecycle rules. |
| 7 | Canonical Membership/glossary | Omits provenance/freshness, guest semantics, System ownership, responsibility, and effective expiry projections | Apply F18-RD-05. |
| 8 | Canonical GovernanceProfile | Activates future-domain semantics | Limit v1 to F18-RD-18 and preserve FEATURE-0020 resolution. |
| 9 | Canonical approval/exception catalog | ApprovalPolicy, ApprovalRequest, and ExceptionGrant contracts and writer/lifecycle behavior are incomplete | Apply F18-RD-12, F18-RD-13, and F18-RD-17. |
| 10 | FEATURE-0013 registry | Lacks the three adopting profiles and mandatory-emission semantics | Apply F18-RD-19 and F18-RD-20 without an AccessReview DecisionProfile. |
| 11 | FEATURE-0017 adoption | Compatible v1 subject/target use is not bounded | Preserve FEATURE-0017 and apply only F18-RD-10 and F18-RD-11. |
| 12 | Lower-authority experimental material | Omits System, authorizes through Membership, stores roleAssignmentRefs, uses inconsistent scopes, or assigns future semantics | Correct or retire it against F18-RD-01 through F18-RD-24. |
| 13 | Canonical scope registrations | Lacks the initial applicability registry, generic extension rule, and exact reference compatibility | Apply F18-RD-02. |
| 14 | Canonical action/role contracts | Lacks monotonic privilege classification, fail-safe ambiguity, and provenance | Apply F18-RD-06 and F18-RD-07. |
| 15 | Approval applicability/bypass | Lacks operation-local ApprovalRequirement and could imply ApprovalPolicy break-glass bypass | Apply F18-RD-10, F18-RD-12, F18-RD-15, and F18-RD-18. |
| 16 | Canonical AccessReview | Lacks bounded modes, snapshot/usage evidence, disposition, timing, containment, and due-advancement semantics | Apply F18-RD-16. |
| 17 | Canonical effect writers | Approval and workflow text can imply multiple RoleAssignment writers | Apply F18-RD-02, F18-RD-08, F18-RD-13, F18-RD-14, and F18-RD-16. |
| 18 | Canonical ExceptionGrant | Lacks exact subject/scope, interval, overlap, revocation, terminal/retry, expiry, and FEATURE-0020 boundaries | Apply F18-RD-17. |
| 19 | Canonical PrivilegedAccessRequest | Lacks exact approval mapping, mode/deadline/duration, readiness, break-glass, and linked-review semantics | Apply F18-RD-14 and F18-RD-15. |
| 20 | Canonical operation/action registry | Lacks the closed operation surface, action classifications, target-binding variants, and sole RoleAssignment writer | Apply F18-RD-02, F18-RD-06, and F18-RD-08. |
| 21 | Canonical assurance | Lacks the closed provider-neutral carrier and deterministic thresholds | Apply F18-RD-03. |
| 22 | Canonical lifecycle/state text | Lifecycle, informational supersession, campaign state, and item dispositions are inconsistent or incomplete | Apply the exact resource-specific F18-RD lifecycle rules. |
| 23 | Governance/exception typed registrations | GovernanceProfile rule schemas and ExceptionGrant control-owned registrations are incomplete | Apply F18-RD-17 and F18-RD-18. |
| 24 | FEATURE-0013 AuditEvent extension | Lacks the exact F18 taxonomy and pre-publication obligation boundary | Apply F18-RD-19 and F18-RD-20. |
| 25 | FEATURE-0018 architecture digest | Resolves references before structural/local validation | Replace its sequence with F18-RD-21. |
| 26 | Matrix C11/C17 | Omits ScopeOnly from target binding | Validate all F18-RD-06 variants; retain NOT READY until evidence exists. |
| 27 | Matrix assessment boundary | Omits AccessGroup from the explicit in-scope inventory | Add F18-RD-04 AccessGroup and keep external/nested/dynamic groups excluded. |
| 28 | Matrix H01 | Orders safe resolution before structural/local validation | Replace its sequence with F18-RD-21. |
| 29 | Matrix D14/F18-SCN-18 | Implies ApprovalPolicy bypass and optional phishing resistance | Apply F18-RD-15 privileged-rule ownership and mandatory phishing resistance. |
| 30 | Matrix B06/B13/F18-SCN-36 | Omits Membership uniqueness/expiry distinctions and owner-inactivity boundaries | Apply F18-RD-04 and F18-RD-05. |
| 31 | Matrix C05/F18-SCN-09/F18-SCN-31 | Omits publication-time delegation semantics and uses qualitative Standing scope | Apply F18-RD-08, F18-RD-09, and F18-RD-18. |
| 32 | Matrix D06/D11 | Omits indirect-beneficiary conflicts, exact eligibility, and activation-deadline atomicity | Apply F18-RD-12 and F18-RD-14. |
| 33 | Matrix E03 and exception checks | Conflates review state/disposition and omits review variants and exception terminal/retry/narrowing rules | Apply F18-RD-16 and F18-RD-17. |
| 34 | Matrix F02 | Does not distinguish GovernanceProfile supersession and retirement | Apply F18-RD-18. |
| 35 | Matrix C14/H10 | Omits the two-plane Sovrunn/native-IAM intersection | Apply F18-RD-06 without introducing provider implementation. |
| 36 | Matrix H08/process text | Defers ASVS too late and uses stale handoff/readiness language | Apply F18-RD-23 and the approved ADH → approval → Kiro application → final evidence sequence. |

Accepted authorities preserved:

- DEC-0018/0019 enterprise Organization governance;
- DEC-0020 non-weakenable governance, whose effective resolution stays in
  FEATURE-0020;
- DEC-0026 reuse before build;
- DEC-0028 policy-engine boundary and no custom policy engine;
- DEC-0035 customer simplicity;
- DEC-0036 replaceable integration boundaries;
- DEC-0037 seven canonical scopes;
- DEC-0043 FEATURE-0013 decision profiles;
- DEC-0044 immutable published definitions;
- DEC-0048 extensibility through governed data/contracts/adapters; and
- DEC-0050/ADH-2026-033 single EffectiveGovernanceContext resolution.

## Required action

- Create `ACR-2026-002` or the next available ACR for the FEATURE-0018
  canonical IAM/governance change.
- Create `DEC-0060` or the next available DEC and update
  `docs/decisions/DECISION_INDEX.md`.
- Create or amend the controlling IAM/governance RFC only if Kiro's authority
  review requires it.
- Update the current architecture baseline after, and only after, the ACR/DEC
  and this package are approved.
- Update the canonical semantic model, contract catalog, glossary, feature
  sequence output, feature index, roadmap summary, and traceability.
- Create the authoritative FEATURE-0018 architecture and feature-scope files.
- Create the complete FEATURE-0011 reuse assessment and approval-evidence
  record using the canonical standard.
- Register the exact FEATURE-0013 profiles and AuditEvent taxonomy through
  FEATURE-0013's extension process.
- Produce exact resource, operation, action, governed versioned scope-
  applicability registration, reference, writer, mutability, lifecycle,
  projection, validation, idempotency, decision/audit, conformance, and future-
  exclusion ledgers solely by transcribing and organizing the approved F18-RD
  semantics. Kiro must stop rather than fill a missing semantic value. Scope
  registrations must remain generic to core machinery and preserve the initial
  registry, reference-compatibility relationship, and extension/change-control
  rules in F18-RD-02.
- Reconcile the FEATURE-0018 architecture digest and industry-validation matrix
  exactly as enumerated in conflict-ledger items 25 through 36. Preserve
  every matrix `NOT READY` gate and unapplied `DEPENDENCY` row until its
  explicitly required architecture-stage, dependency, or downstream executable
  evidence exists; approval of this package alone must not mark a row `PASS`.
- Inspect Structurizr alignment. No new container, external integration, or
  deployment boundary is authorized, so `workspace.dsl` should change only if
  Kiro identifies a genuine approved structural view impact.
- Do not generate requirements, design, tasks, prompts, or implementation as
  part of applying this package.

## Impacted files

Kiro must validate this inventory and report additions or removals before
application.

Normative approval package (already created; approve and apply only as one unit):

- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md`

Create:

- `docs/governance/architecture-change-requests/ACR-2026-002-feature-0018-governance-iam-approval-exception-foundation.md` or next available ACR
- `docs/decisions/DEC-0060-feature-0018-governance-iam-approval-exception-foundation.md` or next available DEC
- `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `docs/features/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/reuse-assessments/FEATURE-0018-approval-evidence.md`
- `.automation/features/FEATURE-0018.control.json` when the architecture control
  manifest is prepared

Update:

- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/ARCHITECTURE_VERSION.md` only if the accepted change requires a
  baseline version increment
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/context/SOVRUNN_CONTEXT_PACK.md`
- `docs/decisions/DECISION_INDEX.md`
- `docs/glossary.md`
- `docs/architecture/canonical/sovrunn-finalized-data-model.md`
- `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md`
- `docs/features/FEATURE_INDEX.md`
- `docs/traceability/DECISION_TRACEABILITY_MATRIX.md`
- `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md`
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`
- `docs/reviews/architecture-readiness/FEATURE-0018-architecture-digest.md`
- `docs/reviews/architecture-readiness/FEATURE-0018-industry-validation-matrix.md`
- this core ADH, only to finalize human-approval metadata before exact package approval

Inspect and update only if the approved structural model requires it:

- `docs/diagrams/structurizr/workspace.dsl`
- controlling IAM/governance RFC files identified by Kiro
- FEATURE-0013 profile registry and conformance authorities identified by Kiro
- FEATURE-0016 action registry/traceability, adding only the approved exact
  target-binding registrations and privileged-handling classification metadata
  without changing action meaning
- FEATURE-0017 traceability, without changing the FEATURE-0017 contract
- VS-000 scenario and conformance inventories

## Impacted features

- FEATURE-0011: supplies the reuse assessment format unchanged; FEATURE-0018
  must add its own conforming assessment and approval evidence.
- FEATURE-0012: common resource, reference, scope, validation, safe-denial,
  redaction, concurrency, and idempotency authority is reused unchanged;
  governed resource scope-applicability registrations must be consumed
  generically without adding a second scope authority.
- FEATURE-0013: registers three new adopting-domain DecisionProfiles and the
  bounded AuditEvent taxonomy; envelope semantics remain unchanged.
- FEATURE-0016: its ExecutionTarget actions and lifecycle authority remain
  unchanged; create/retire action granularity is a recorded limitation and any
  privileged-handling classification is additive security metadata rather than
  a reinterpretation of action meaning.
- FEATURE-0017: its seam and fake are reused unchanged and its v1 target-aware
  limitation is preserved.
- FEATURE-0018: primary owner and only current implementation scope.
- FEATURE-0019: no sovereignty fact, evidence, policy, or placement semantics
  are moved into FEATURE-0018.
- FEATURE-0020: continues to own ProfileAssignment, inheritance, conflict and
  exception application, and EffectiveGovernanceContext resolution.
- FEATURE-0021 and later: no entitlement, quota, service, placement, plugin,
  execution, explanation, or future adapter behavior is pulled forward.
- FEATURE-0026: will later consume FEATURE-0018 conformance in VS-000; it gains
  no current implementation scope.

## Acceptance criteria for Kiro update

- [ ] This exact three-file package has explicit human approval before application.
- [ ] The change is classified and applied through the required ACR, DEC,
      Decision Index, traceability, and explicit RFC disposition; baseline files
      change only after that approval chain is valid.
- [ ] Canonical and dependency artifacts reconcile conflict-ledger items 1
      through 24 solely from their controlling F18-RD references; Kiro records
      and stops on any value that cannot be transcribed without inference.
- [ ] The architecture digest and industry-validation matrix reconcile items 25
      through 36 exactly, while every unevidenced gate remains `NOT READY`, every
      unapplied dependency remains `DEPENDENCY`, and no prose-only `PASS` is
      recorded.
- [ ] Resource, supporting-value, scope/reference, action/target-binding,
      operation, writer, mutability, lifecycle, projection, validation,
      concurrency/idempotency, decision, and audit ledgers demonstrate exact
      conformance to F18-RD-02 through F18-RD-21 without creating a second
      authority.
- [ ] FEATURE-0013, FEATURE-0016, and FEATURE-0017 changes are limited to the
      registrations and adoption boundaries explicitly assigned by F18-RD-06,
      F18-RD-10, F18-RD-11, F18-RD-19, and F18-RD-20; their owning contracts are
      not reinterpreted.
- [ ] The complete FEATURE-0011 reuse assessment and separate human approval
      evidence use only the controlled disposition and decision-status
      vocabularies.
- [ ] F18-RD-22 conformance preserves F18-IVM-A01 through F18-IVM-I06 and
      F18-SCN-01 through F18-SCN-49, with F18-SCN-38 retired into F18-SCN-29;
      architecture-stage inventories, deterministic vectors, evidence owners,
      journey projections, and traceability exist before requirements, while
      executable proof remains a downstream feature-gate obligation.
- [ ] F18-RD-23 standards mappings, applicable ASVS profile, CIS/CSA
      responsibility evidence, future-adapter paper mappings, and independent
      security review are completed at their approved architecture gates without
      claiming unassessed conformity.
- [ ] F18-RD-24's six progressive-disclosure journeys are preserved without a
      facade resource or user-authored system/provenance fields.
- [ ] F18-RD-01 exclusions remain closed: no future-owned resource or field,
      real adapter, external call, persistence, provider execution, deployment
      topology, requirements, design, task, prompt, or implementation is
      introduced by this update.
- [ ] `git diff --check`, `mkdocs build --strict`, relevant architecture-drift
      checks, and `make structurizr-check` when structural DSL changes all pass.

## Explicit instructions to Kiro

- Treat this package as `Proposed` until the architecture owner approves this
  exact three-file package; approval of source discussions or earlier artifacts is not
  approval of this package revision.
- Validate the approved package against every higher-authority source before
  application and report any incompatibility.
- Preserve F18-RD-01 through F18-RD-24 as the sole behavioral authority. Do not
  add a 25th group, infer a missing value, or resolve inconsistent summaries by
  precedence; stop and surface the conflict.
- Apply conflict-ledger items 1 through 36 only through their cited controlling
  decisions. Preserve the existing ownership of FEATURE-0012, FEATURE-0013,
  FEATURE-0016, FEATURE-0017, FEATURE-0020, and every future feature.
- Update the current architecture baseline only after the ACR, DEC, this
  package, and the canonical updates are approved and internally consistent.
- Do not mark an industry-validation row `PASS` without its named evidence.
- Do not generate requirements, design, tasks, generated prompts,
  implementation, routes, persistence, adapters, external calls, deployment
  topology, Structurizr changes without an approved structural impact, or Go
  code as part of this architecture update.
- Record exact human approval and authority metadata without turning discussion
  chronology into architecture authority.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar, Sovrunn Architecture Owner
- Date: 2026-08-30
- Approved package evidence: final-approval manifest SHA-256
  `9f69ea2fb89ba04aa253a37a0cf5f63124753c2b1e9cf08cf0a5c028369d8efc`
- Independent security evidence: renewal-05 SHA-256
  `2f95272aebb23fbfdca507de6434461f5e66378ddc206e9ff5b382848a905ccf`
- Notes: Approval covers the exact security-reviewed three-file package and its
  canonical application. Requirements, design, tasks and implementation remain
  separately gated. Detailed discussion history remains non-authoritative
  session provenance.

This Architecture Decision Handoff is approved and has been applied through the
controlled repository reconciliation. It does not itself authorize a later
Feature Factory stage.
