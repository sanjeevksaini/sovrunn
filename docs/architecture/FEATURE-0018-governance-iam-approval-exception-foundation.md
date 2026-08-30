---
doc_type: architecture
feature: FEATURE-0018
title: Governance, IAM, Approval and Exception Foundation
status: approved
phase: 2R
baseline: ARCH-2026.08-PHASE2R-CANONICAL
controlling_handoff: ADH-2026-070
conformance_handoff: ADH-2026-071
controlling_core_sha256: 489ccb5869ca4cc43835ef79de2d2b8f006f580cc22fdad8a6d674b112fdd769
normative_appendix_a_sha256: 66479426f30759a21202c8beeb08989241d912b7ba69de3619619f45823fe485
normative_appendix_b_sha256: 2424cff0130bf82d7b413069fc60580af0ab6e1a4bd00982fc669bebb80d8b49
updated: 2026-08-30
ai_load_priority: always
ai_summary: Approved FEATURE-0018 architecture authority generated from the approved ADH-2026-070 package; it indexes the immutable F18-RD decision payload, preserves ownership boundaries, and does not independently authorize requirements or implementation.
---

# FEATURE-0018 Governance, IAM, Approval and Exception Foundation

## 1. Authority and current status

This document is the compact FEATURE-0018 architecture index and generation
boundary produced from the human-authorized `ADH-2026-070` package. It adds no
architecture decision beyond that package.

The exact normative decision text remains in the two immutable handoff
appendices. This document organizes those decisions for canonical application,
architecture-leakage proof, final review, and later requirements generation.
The core handoff and both appendices must be loaded with this document:

- [ADH-2026-070 core](../reviews/architecture-decision-handoffs/ADH-2026-070-feature-0018-governance-iam-approval-exception-foundation.md);
- [normative Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md); and
- [normative Appendix B](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md).

[ADH-2026-071](../reviews/architecture-decision-handoffs/ADH-2026-071-feature-0018-conformance-executability-reconciliation.md)
is the approved non-semantic conformance-executability correction. It registers
stable local proof identifiers and REQ/AC mappings only; ADH-2026-070 and
DEC-0060 remain the sole runtime-semantic authority.

The SHA-256 values in front matter pin the package used to generate this
document. A digest mismatch, missing file, duplicated decision group, or
inconsistency is an architecture blocker. No file has precedence; Kiro must
stop and return the complete package for renewed human review.

Current architecture status is **Approved**, including the ADH-2026-071
conformance reconciliation. Mutable Feature Factory stage authorization remains
outside this architecture document. This document does not independently
authorize design, tasks, implementation, or source changes.
The first independent review rejected the prior payload with F18-SEC-001
through F18-SEC-005. Renewal-01 closed F18-SEC-002 and F18-SEC-004 but rejected
the payload with F18-REN-001 and F18-REN-002. Renewal-02 closed all of those
reported defects and identified F18-REN2-001 and F18-REN2-002. Renewal-03
closed the decision-status defect and the exact group-held RoleAssignment path,
then identified the remaining group-derived eligibility path as F18-REN3-001
and the undefined reviewer-eligibility contract as F18-REN3-002. Renewal-04
closed both defects and F18-REN2-001 end-to-end, then identified relationship-
derived role eligibility as hidden delegation channel F18-REN4-001. Its exact-
Human PrincipalRef-only correction is incorporated here. Renewal-05 verified
the regenerated exact payload with no blocking findings. ACR-2026-002,
DEC-0060, the reuse boundary and exact package were approved on 2026-08-30.

Mutable Feature Factory stage status must remain outside this architecture.

## 2. Architecture outcome

FEATURE-0018 establishes Sovrunn's implementation-neutral control-plane
authorization foundation for Human, Workload, and System principals. It owns:

- stable use of already-authenticated principals;
- non-authorizing Organization, CloudProvider, and direct AccessGroup
  membership relationships;
- immutable versioned roles expressed only through canonical Sovrunn actions;
- scoped direct-principal and AccessGroup role assignments;
- deterministic grant-union and guardrail-intersection authorization;
- bounded approval, JIT, break-glass, review, and exception evidence;
- a FEATURE-0018-limited GovernanceProfile composition envelope;
- exact FEATURE-0013 decision and audit adoption;
- deterministic in-memory fixtures and local conformance; and
- six progressive-disclosure user journeys that normally hide security
  machinery from end customers.

FEATURE-0018 is a control-plane authorization foundation. It is not an
identity provider, external group synchronizer, provider-native IAM bridge,
policy engine, workflow engine, effective-governance resolver, execution
engine, production persistence design, or deployment topology.

## 3. Phase, dependency, and ownership boundary

| Boundary | FEATURE-0018 treatment | Controlling authority |
|---|---|---|
| FEATURE-0011 | Reuse assessment format and evidence are adopted unchanged | DEC-0026; canonical reuse standard |
| FEATURE-0012 | Metadata, seven scopes, typed references, validation, safe denial, redaction, concurrency, idempotency, and problems are reused | ADH-2026-012 |
| FEATURE-0013 | DecisionRecord, DecisionProfile, EvaluationResult, and AuditEvent remain single-owner contracts; FEATURE-0018 only registers profiles and taxonomy | DEC-0043; ADH-2026-017 |
| FEATURE-0016 | ExecutionTarget actions and lifecycle remain FEATURE-0016-owned; FEATURE-0018 consumes exact registrations without reinterpretation | ADH-2026-058 through ADH-2026-066 |
| FEATURE-0017 | PolicyEvaluationRequest/Result, PolicyEngineAdapter, deterministic fake, and pure mapper are reused unchanged | ADH-2026-067 through ADH-2026-069 |
| FEATURE-0018 | Owns only the resources, supporting values, composition, flows, fixtures, and proof identified by F18-RD-01 through F18-RD-24 | ADH-2026-070 |
| FEATURE-0019 | Sovereignty facts, evidence, policy, and placement semantics remain excluded | Phase 2R sequence |
| FEATURE-0020 | ProfileAssignment, hierarchy, conflict/exception application, and EffectiveGovernanceContext resolution remain excluded | DEC-0050; ADH-2026-033 |
| FEATURE-0021 and later | Enrollment, entitlement, quota, service, placement, plugin, execution, and explanation behavior remain excluded | Phase 2R sequence |

The current phase permits deterministic, in-process, side-effect-free contract
and simulation behavior only. No real IdP, OIDC/SCIM adapter, external group
provider, policy/workflow engine, provider SDK, native IAM mutation, database,
Kubernetes call, provisioning, credential acquisition, or infrastructure
effect is authorized.

## 4. Normative decision register

F18-RD-01 through F18-RD-24 are the complete decision set. The appendices are
allocated by whole decision group so tables and behavioral contracts are not
split across files. Titles below are navigation labels, not restatements.

| Decision | Normative title | Location |
|---|---|---|
| F18-RD-01 | Current-feature-only authority | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-02 | Closed contract inventory and profiles | [Appendix B](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-03 | Stable principal identity | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-04 | Direct principal and scoped AccessGroup assignments | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-05 | Membership is non-authorizing | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-06 | Distributed action ownership and central role composition | [Appendix B](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-07 | Versioned RoleDefinition | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-08 | Scoped RoleAssignment | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-09 | Deterministic scoped authorization composition | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-10 | FEATURE-0017 adoption without reinterpretation | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-11 | Bounded FEATURE-0017 subject/target use | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-12 | Bounded ApprovalPolicy | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-13 | Immutable terminal approval evidence | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-14 | JIT privileged access authorizes one temporary grant | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-15 | Constrained break-glass access | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-16 | Snapshot-based AccessReview | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-17 | Bounded immutable exception evidence | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-18 | FEATURE-0018-limited GovernanceProfile v1 | [Appendix B](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-19 | FEATURE-0013 adoption | [Appendix B](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-20 | Audit before authorization-changing publication | [Appendix B](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-21 | Deterministic validation and safe denial | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |
| F18-RD-22 | Deterministic in-memory foundation and local conformance | [Appendix B](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-23 | Standards-validation gate | [Appendix B](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md) |
| F18-RD-24 | Progressive and normally hidden user friction | [Appendix A](../reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md) |

No downstream artifact may add a 25th group, merge or split a group, invent a
missing value, or resolve inconsistent summaries by precedence.

## 5. Closed contract model

This section is a navigation summary only. F18-RD-02 owns the exact inventory,
profiles, supporting values, enums, allowed scopes, reference compatibility,
writers, operations, and mutability.

| Contract | Architectural role |
|---|---|
| `PrincipalRef` | Stable reference to an already-authenticated Human, Workload, or System identity |
| `AccessGroup` | Non-authenticating, access-only authorization subject with explicit ownership |
| `Membership` | Non-authorizing belonging or direct AccessGroup relationship |
| `RoleDefinition` | Immutable published version of canonical Sovrunn actions |
| `RoleAssignment` | Scoped, explicitly valid grant of one exact role version to one role holder |
| `PrivilegedAccessRequest` | Justified, approved, time-bounded privileged-access request |
| `AccessReview` | Immutable-snapshot review and retained remediation evidence |
| `ApprovalPolicy` | Immutable published approval-stage, eligibility, quorum, expiry, and separation-of-duties definition |
| `ApprovalRequest` | Retained approval evidence for one bounded FEATURE-0018 intent or proposal |
| `ExceptionGrant` | Immutable, scoped, time-bounded exception or linked revocation evidence |
| `GovernanceProfile` | Versioned composition of FEATURE-0018-owned governance references and constraints only |

Supporting values such as `RoleHolderRef`, `AuthorizationInput`,
`AuthorizationResult`, `ApprovalRequirement`, `EligibilityRef`,
`MembershipEnabledGrantEnvelope`,
proposals, usage evidence, governance applicability, and action-target bindings have no independent
endpoint, lifecycle, storage authority, or user journey.

No `ExceptionRequest`, generic workflow, facade, condition language, external
group resource, provider-native IAM resource, credential resource, or future
governance resolver is introduced.

## 6. Authorization and trust composition

The authorization model is explicit grant union constrained by guardrail
intersection. Each assignment must independently match holder, exact role
version, scope, optional exact resource, validity, required Membership,
responsibility, review, and lifecycle. Restrictions from separate assignments
never synthesize a grant.

```text
authenticated PrincipalRef
  + current direct and/or direct-AccessGroup assignment candidates
  -> independently applicable grant union
  -> exact action/target/scope/resource match
  -> mandatory membership, responsibility, validity and review constraints
  -> FEATURE-0017 policy and contextual guardrail intersection
  -> current operation-local approval/exception evidence when required
  -> Allow or Deny with exact non-bearer provenance
  -> protected decision/audit obligations before sensitive publication
```

Policy `Allow`, approval, exception, Membership, controller identity, native
cloud IAM, and retained authorization results never create Sovrunn authority.
Missing, stale, unresolved, malformed, wildcard, ambiguous, or conflicting
security inputs fail closed according to F18-RD-09 and F18-RD-21.

An AccessGroup Membership creation, reactivation, or freshness/provenance
change that restores or extends group-derived assignment effect is an atomic
grant-producing boundary. The Membership controller derives the exact current
`MembershipEnabledGrantEnvelope`; a non-empty envelope requires the applicable
Membership authority, `roleassignment.grant` over its complete reach, and one
non-synthesized per-action scope/target/resource/validity witness. A Standing
effect requires Standing witnesses. Empty-group membership grants no action;
later group-assignment publication independently passes the ordinary grantor
ceiling. Trusted provisioners receive no bypass, concurrent membership/
assignment changes conflict or retry, and privilege-reducing Membership
transitions remain available. Membership cannot establish approver, reviewer,
JIT, or break-glass eligibility. Those lists use non-empty `EligibilityRef`
OR semantics and contain exact Human PrincipalRefs only. RoleDefinition,
RoleAssignment, AccessGroupRef, Membership, claim, and other relationships
cannot satisfy or create eligibility. Eligibility never supplies the
independently required canonical action and is re-evaluated at decision/effect
publication; policy/rule applicability and the action alone determine scope
and target reach.

### 6.1 ExecutionTarget and native IAM boundary

FEATURE-0018 authorizes canonical Sovrunn control-plane actions only.
Provider-native principals, groups, roles, policies, permissions, and
credentials are neither FEATURE-0018 resources nor outputs.

When a separately owned execution or integration feature carries an authorized
operation to an ExecutionTarget, both independent layers must permit it:

```text
effective native effect
  = valid Sovrunn control-plane authorization
    INTERSECT successful native IAM authorization
```

The responsible later controller or adapter uses its own least-privileged
native identity. Neither layer's Allow substitutes for the other, and neither
layer can override the other's Deny. Direct out-of-band IaaS administration
remains CloudProvider-owned and cannot be represented as the effect of a
Sovrunn RoleAssignment.

FEATURE-0018 defines only this trust boundary. Native identity acquisition,
credential management, permission translation, synchronization, and effect
execution remain outside this feature.

## 7. Single-writer and evidence boundary

Every resource has the exact writer split registered in F18-RD-02. In
particular, the RoleAssignment controller is the sole canonical publisher and
sole writer of specification, lifecycle, status, revocation, replacement,
expiry, and review-due effects. Administrators and grant, approval,
privileged-access, and review controllers may submit only exact immutable
intents; they cannot directly write any part of a RoleAssignment.

Delegated publication requires `roleassignment.grant` over the proposal's
complete target reach and one independently applicable delegable authority
witness for every proposed action. Each witness must cover the proposed scope,
target-binding mode, optional exact resource, validity, and delegation
capability. Resource-A authority cannot grant resource B or scope-wide reach,
and restrictions from different assignments cannot be synthesized.

Every interactive privileged activation, including normal JIT and break-glass,
requires phishing-resistant AssuranceEvidence at AAL2 or higher and no more
than fifteen minutes old; the applicable privileged-access rule may strengthen
but not weaken that floor.
Every access-preserving or replacing review disposition excludes the direct
Human holder, the responsible Human for a Workload/System holder, and, for an
AccessGroup holder, its owner, every current direct Human member, and the
responsible Human of every current direct Workload/System member. The complete
set is re-evaluated at its atomic decision/effect boundary; unresolved evidence
fails closed for Retain/Replace without blocking Revoke.

Approval is admission-time evidence. It must be current at downstream effect
publication, but its later expiry does not silently revoke an already-published
RoleAssignment. Continuing authorization is controlled by assignment validity,
dependencies, review state, guardrails, suspension, and explicit revocation.

Approval publication and downstream effect publication are distinct atomic
operations. Sensitive state cannot publish unless required decision and audit
obligations have been accepted as defined by F18-RD-20.

## 8. Feature-level reuse summary

This is the single FEATURE-0018 architecture-level summary required by the
canonical Reuse Assessment Standard. The capability assessments and approval
evidence were accepted during final architecture review.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 common resource/API foundations | Reuse | Existing canonical metadata, scope, reference, validation, error, redaction, concurrency, and idempotency behavior must not be duplicated | Approved | DEC-0026; ADH-2026-012 |
| FEATURE-0013 decision and audit envelopes | Reuse | FEATURE-0013 remains the single evidence-envelope owner | Approved | DEC-0043; ADH-2026-017 |
| FEATURE-0016 ExecutionTarget actions and lifecycle | Reuse | FEATURE-0018 consumes exact registrations without splitting or reinterpreting action meaning | Approved | ADH-2026-058 through ADH-2026-066 |
| FEATURE-0017 policy-evaluation seam and fake | Reuse | The deterministic engine-neutral boundary is consumed unchanged within its v1 limitation | Approved | DEC-0028; DEC-0036; ADH-2026-067 through ADH-2026-069 |
| Applicable security, identity, cloud, API, and assurance standards | Reuse | Standards constrain invariants and proof without becoming vendor-native canonical schemas | Approved | F18-RD-23; industry-validation matrix |
| Sovrunn IAM and authorization algebra | Build | No mature component owns Sovrunn's seven-scope role-holder, action, guardrail, and evidence model | Approved | F18-RD-02 through F18-RD-11 |
| Sovrunn approval, JIT, break-glass, review, exception, and deterministic fake semantics | Build | These are Sovrunn control-plane contracts; an external workflow runtime cannot own their canonical meaning | Approved | F18-RD-12 through F18-RD-17; F18-RD-22 |
| FEATURE-0018 GovernanceProfile composition and progressive projections | Build | Sovrunn owns this bounded envelope and experience while FEATURE-0020 retains effective resolution | Approved | ADH-2026-033; F18-RD-18; F18-RD-24 |
| Real IdP, federation, or external-group provisioning integration | Wrap | Non-authoritative future compatibility hypothesis requiring fresh FEATURE-0011 assessment | Deferred | DEC-0022; DEC-0036 |
| Real policy-engine adapter | Wrap | Non-authoritative future hypothesis; FEATURE-0017 selects no production engine | Deferred | DEC-0028; DEC-0036 |
| Real workflow-engine adapter | Wrap | Non-authoritative future hypothesis; FEATURE-0018 selects no runtime | Deferred | DEC-0026; DEC-0036 |

Every `Wrap`/`Deferred` row is non-authoritative. It selects no product,
protocol, adapter, credential flow, schema, deployment boundary, requirement,
task, or implementation. The future owning feature must perform a fresh
FEATURE-0011 assessment and may choose any controlled disposition supported by
then-current evidence.

For each Build unit, Reuse is insufficient because no mature component owns
the complete canonical Sovrunn identity, scope, action, lifecycle, evidence,
and progressive-projection model. Wrap is insufficient because Phase 2R
selects no external identity, policy, or workflow runtime. Extend is
insufficient because it would transfer FEATURE-0018 semantics into a prior
feature or external engine. Sovrunn builds only its implementation-neutral
contracts, deterministic composition, fakes, and local proof.

## 9. Architecture-leakage proof contract

Architecture leakage is a release-blocking failure, not a documentation
warning. Before final architecture approval, Kiro must produce and trace the
following architecture-stage evidence without inventing semantics.

| Proof area | Required evidence | Controlling decisions |
|---|---|---|
| Closed inventory | Exact resource/supporting-value profiles, allowed scopes, references, and absence of unapproved kinds | F18-RD-01, F18-RD-02 |
| Identity boundary | Stable issuer/subject identity, no email authority, no authentication lifecycle, and no direct external-group claim grant | F18-RD-03 through F18-RD-05 |
| Writer boundary | One accepted-intent/spec writer and one status/result/record writer per contract; sole RoleAssignment publisher | F18-RD-02, F18-RD-08, F18-RD-13, F18-RD-14, F18-RD-16 |
| Scope and target binding | Exact governed scope-applicability registrations and discriminated ExactResource/CreateParent/ScopeOnly action bindings | F18-RD-02, F18-RD-06 |
| Authorization algebra | Independent grant union, guardrail intersection, exact provenance, ambiguity denial, and no evidence-derived grants | F18-RD-07 through F18-RD-11 |
| Approval and privilege | Exact eligibility, stages, quorum, separation of duties, activation deadline, time bounds, no bypass leakage, and independent publication boundaries | F18-RD-12 through F18-RD-15 |
| Review and exception | Immutable snapshot/usage evidence, mode-specific review rules, deterministic remediation, exact exception narrowing, and FEATURE-0020 non-resolution | F18-RD-16, F18-RD-17 |
| Governance boundary | Only registered FEATURE-0018 v1 components; no sovereignty, placement, backup, cost, entitlement, quota, execution, or effective resolution | F18-RD-18 |
| Decision/audit adoption | Exactly three FEATURE-0013 profiles, bounded audit taxonomy, and audit-before-publication acceptance | F18-RD-19, F18-RD-20 |
| Determinism and failure | Exact validation precedence, bounded inputs, safe reference resolution, idempotency/concurrency behavior, and privacy-preserving denial | F18-RD-21 |
| Local proof | Immutable synthetic fixtures, positive/negative/race/replay/failure vectors, journey evidence, and zero external effects | F18-RD-22, F18-RD-24 |
| Standards boundary | Exact standards applicability, no unassessed conformity claim, future-adapter paper mappings, CIS/CSA responsibility evidence, and independent security review | F18-RD-23 |
| Native IAM boundary | Two-plane authorization intersection and absence of provider-native resources, credentials, translation, or execution | F18-RD-06 |

### 9.1 Required architecture-stage ledgers

The evidence package must contain exact, cross-referenced ledgers for:

- resource and supporting-value inventory;
- resource profile and allowed canonical scope applicability;
- reference kind, scope compatibility, and containment;
- accepted-intent/spec, status/result, and field ownership;
- operation, mutability, lifecycle, and effective projection;
- canonical action, owner, intrinsic classification, and target binding;
- RoleDefinition classification and RoleAssignment applicability;
- approval, privileged-access, review, and exception state/evidence flow;
- GovernanceProfile component applicability;
- validation precedence, concurrency, idempotency, and terminal/retry behavior;
- FEATURE-0013 decision profiles and AuditEvent taxonomy;
- FEATURE-0017 request/result adoption boundary;
- six progressive-disclosure journeys;
- FEATURE-0018-local conformance and deterministic vectors;
- standards/control mappings and independent-review ownership; and
- explicit future-feature and provider-native leakage sentinels.

Each ledger must identify its sole controlling F18-RD source. A ledger may
organize or transcribe approved semantics but cannot become a second authority.
Missing values produce `ARCHITECTURE_DECISION_REQUIRED`; they must not be
filled from examples, implementation preference, vendor behavior, or future
roadmap text.

### 9.2 Mandatory negative leakage sentinels

Architecture review must fail if any active FEATURE-0018 contract, proposed
requirement, or generated artifact introduces:

- identity-provider selection, token validation, federation, SCIM transport,
  or external-group synchronization;
- external claims as direct authorization grants;
- nested or dynamic AccessGroup membership;
- provider-, Kubernetes-, or engine-native principals, groups, roles,
  policies, permissions, credentials, or schemas;
- a real policy or workflow engine, adapter selection, routing, retry,
  credential, or deployment topology;
- a second decision, audit, authorization, scope, GovernanceProfile, effective
  governance, or RoleAssignment writer authority;
- user-authored status, evidence, provenance, arbitrary condition expressions,
  wildcard permissions, or deny assignments;
- ProfileAssignment, inheritance, conflict resolution, exception application,
  or EffectiveGovernanceContext construction;
- sovereignty, placement, backup, cost, enrollment, entitlement, quota,
  billing, provisioning, execution, or AI-explanation semantics;
- production persistence, public facade resource, raw secret, protected native
  handle, external call, or infrastructure effect; or
- requirements, design, tasks, implementation, or future-feature artifacts
  during architecture application.

## 10. Industry-validation and evidence boundary

The [FEATURE-0018 industry-validation matrix](../reviews/architecture-readiness/FEATURE-0018-industry-validation-matrix.md)
is a review and evidence contract, not architecture authority or a conformity
claim. Its `APPROVED_FOR_HANDOFF` rows record disposition; `DEPENDENCY`,
`INHERITED`, `NOT_ASSESSED`, and `NOT READY` remain evidence states.

Approval of ADH-2026-070 or generation of this document must not turn any row
into `PASS`. Final architecture approval requires the matrix evidence package,
including standards applicability, responsibility mappings, threat/abuse
review, exact ledgers, scenario disposition, reuse assessment, adapter paper
mappings, independent security review, and explicit architecture-owner
approval.

Executable conformance remains a later feature-gate obligation. Architecture
journey evidence is not usability-test or runtime proof.

## 11. Requirements-generation boundary

Requirements may be generated only after this architecture and its evidence
package receive explicit final architecture approval. The architecture gate
authorizes `requirements.md` only; design, tasks, Cursor execution, and
implementation retain their separate approvals.

Requirements must translate exact observable behavior from F18-RD-01 through
F18-RD-24 and must trace every requirement to one or more decision groups.
They may organize outcomes by resource or journey, but cannot:

- redefine the resource inventory, enums, scopes, actions, writers,
  lifecycles, state, or validation precedence;
- select Go packages, concrete types, method signatures, persistence,
  controllers, transport routes, deployment topology, or external products;
- convert examples, journey labels, standards mappings, or future
  compatibility hypotheses into new behavior;
- reinterpret FEATURE-0012, FEATURE-0013, FEATURE-0016, FEATURE-0017, or
  FEATURE-0020 authority;
- infer an unregistered action, target binding, scope, lifecycle transition,
  exceptionable control, rule applicability, or evidence value; or
- weaken the future-feature exclusions or negative leakage sentinels.

No requirement may cite the architecture digest, industry matrix, roadmap, or
discussion history as independent semantic authority. Those artifacts provide
context, validation, or provenance only.

### 11.1 Executable conformance mapping

ADH-2026-071 closes the Slice-0 executability correction without changing any
F18-RD semantic:

```text
AC-F18-01..37 -> VS0-CF-F18-01..37 (same suffix)
AC-F18-38     -> excluded; VS0-CF-F18-38 permanent tombstone
AC-F18-39..49 -> VS0-CF-F18-39..49 (same suffix)
REQ-only proof -> VS0-CF-F18-50..54
```

Every active `VS0-CF-F18-*` row is owned by FEATURE-0018 and its exact machine
fields are authoritative only in
`docs/architecture/vertical-slices/VS-000-contract-registry.yaml`. The complete
REQ-to-case mapping is ADH-2026-071 section 6. Existing shared cases
`VS0-CF-F01`, `VS0-CF-F02`, `VS0-CF-X01`, and `VS0-CF-X02` remain supplementary.
FEATURE-0015-owned `VS0-CF-X03` is inherited evidence only and never counts as
FEATURE-0018-local acceptance proof.

## 12. User-simplicity boundary

Security complexity remains normally hidden. The six approved journeys are:

1. routine end-customer operation;
2. ordinary access assignment;
3. privileged-access request;
4. approver decision;
5. access review; and
6. exception request.

Users see the smallest risk-proportional projection necessary for their task.
Security administrators select exact eligible Humans in the applicable policy
or rule; end customers do not manage EligibilityRef values or group/role-policy
linkage. This keeps eligibility individually auditable without adding a
separate customer journey.
They do not author controller state, provenance, evidence, decision/audit
records, native IAM mappings, or internal orchestration resources. This
progressive disclosure does not create a facade resource or weaken the exact
underlying contracts.

## 13. Application reconciliation boundary

ADH-2026-070 enumerates 36 canonical and dependent-artifact conflicts. Kiro
must reconcile all 36 only through their cited F18-RD authorities. In
particular:

- conflicts 1 through 24 control canonical and dependency artifacts;
- conflicts 25 through 36 control the architecture digest and industry matrix;
- no unresolved conflict may be silently deferred into requirements or design;
- no artifact may be marked reconciled without exact evidence; and
- any value not directly transcribable from the approved decision payload is a
  blocker, not a design choice.

Structurizr requires no change merely because this document exists. Kiro must
update `workspace.dsl` only if approved application introduces a genuine system
boundary, major container, plugin plane, external integration, deployment
relationship, or major dynamic-flow change. ADH-2026-070 currently authorizes
none of those structural changes.

## 14. Final architecture gate

This document is ready for final architecture review only when:

- the controlling package hashes still match;
- all 24 decision groups exist exactly once in the immutable appendices;
- the ACR/DEC and canonical reconciliation path is accepted;
- every one of the 36 conflicts has disposition and evidence;
- the complete FEATURE-0011 assessment and human approval evidence exist;
- all architecture-stage ledgers and leakage sentinels pass;
- the industry matrix is reconciled without false `PASS` claims;
- independent security review evidence exists;
- `git diff --check` and `mkdocs build --strict` pass;
- `make structurizr-check` passes if and only if Structurizr changed; and
- the architecture owner explicitly approves this exact architecture package.

Until then, requirements, design, tasks, implementation, and Cursor handoff
remain unauthorized.
