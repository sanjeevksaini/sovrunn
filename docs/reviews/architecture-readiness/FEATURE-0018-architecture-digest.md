---
doc_type: architecture_readiness_digest
feature: FEATURE-0018
title: Governance, IAM, Approval and Exception Foundation
status: architecture-approved
baseline: ARCH-2026.08-PHASE2R-CANONICAL
prepared: 2026-08-26
updated: 2026-08-30
---

# FEATURE-0018 Architecture Digest

> **Non-normative readiness index.** This digest does not define FEATURE-0018
> semantics. The jointly approved ADH core and its two normative appendices are
> the sole FEATURE-0018 decision payload. If this index disagrees with that
> payload, the payload wins and this index must be corrected. This file does not
> authorize requirements, design, tasks, or implementation.

## 1. Purpose and current gate

FEATURE-0018 defines the provider-neutral enterprise foundation for deciding
which authenticated Human, Workload, or System principal may perform a
canonical Sovrunn control-plane action, against which scoped target, under
which approval, exception, validity, review, assurance, and audit constraints.

The semantic authority is:

1. [ADH-2026-070 core](../architecture-decision-handoffs/ADH-2026-070-feature-0018-governance-iam-approval-exception-foundation.md);
2. [Appendix A — semantic contract](../architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md); and
3. [Appendix B — registries and evidence](../architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md).

Renewal-02 through renewal-04 remain immutable rejection records. Renewal-04
found direct RoleDefinition qualification as a hidden eligibility-delegation
channel. The exact-Human PrincipalRef-only `EligibilityRef` correction closed
that path, and fresh non-author renewal-05 passed the exact manifest with no
blocking findings. Final architecture and reuse approval were recorded on
2026-08-30.

## 2. Starting point and approved extension

The approved Phase 2R baseline inherited nine FEATURE-0018 output resources:
`GovernanceProfile`, `Membership`, `RoleDefinition`, `RoleAssignment`,
`PrivilegedAccessRequest`, `AccessReview`, `ApprovalPolicy`, `ApprovalRequest`,
and `ExceptionGrant`.

ADH-2026-070 adopts one canonical extension: `AccessGroup`, plus its bounded
supporting values such as `AccessGroupRef` and the `roleHolderRef` union. The
approved post-application inventory therefore has ten persistent resources.
`PrincipalRef`, `RoleHolderRef`, authorization inputs/results, `EligibilityRef`,
`MembershipEnabledGrantEnvelope`, proposals, provenance, and evidence summaries
are supporting values, not additional persistent resources.

This distinction is deliberate: the inherited baseline did not contain
`AccessGroup`; the approved architecture does. External group claims remain
non-authorizing unless trusted provisioning creates current canonical Sovrunn
resources. Nested and dynamic groups are excluded.

## 3. Scope and ownership

FEATURE-0018 owns canonical, provider-neutral semantics for:

- stable principal references and contextual, non-authorizing Membership;
- direct-principal and direct-member AccessGroup role holding;
- immutable versioned RoleDefinition and scoped RoleAssignment;
- deterministic authorization composition and provenance;
- approval, privileged-access, access-review, and exception evidence;
- the FEATURE-0018 GovernanceProfile component registrations;
- FEATURE-0013 decision/audit adoption; and
- deterministic in-memory fixtures and local conformance contracts.

FEATURE-0018 does not own authentication, federation, OIDC/SCIM processing,
provider-native IAM, credential acquisition, a policy/workflow engine,
governance assignment or effective resolution, infrastructure execution,
production persistence, deployment topology, or a new customer facade.

An ExecutionTarget effect requires the intersection of a valid Sovrunn
authorization and the target provider's independent authorization of a
separately owned least-privileged adapter identity. Neither layer's Allow
substitutes for the other, and neither overrides the other's Deny.

## 4. Dependency and reuse boundaries

| Source | Reused authority | FEATURE-0018 limit |
|---|---|---|
| FEATURE-0011 | Reuse Assessment Standard | Records reuse/adapt/build decisions; it does not supply runtime IAM behavior. |
| FEATURE-0012 | Metadata, references, scopes, validation, errors, safe disclosure, concurrency | FEATURE-0018 registers owned contracts without redefining common envelopes or scope containment. |
| FEATURE-0013 | DecisionRecord, DecisionProfile, EvaluationResult, AuditEvent | FEATURE-0018 adopts exact profiles/taxonomy and audit obligations; it creates no competing decision/audit model. |
| FEATURE-0016 | Existing ExecutionTarget actions and target kinds | FEATURE-0018 registers authorization bindings only; it does not execute or translate provider IAM. |
| FEATURE-0017 | Policy evaluation request/result, adapter boundary, deterministic fake | Policy output is a bounded guardrail input, never execution authority or a replacement IAM system. |
| FEATURE-0020 | Governance assignment and effective resolution | FEATURE-0018 owns the profile envelope and typed components only; it does not publish ProfileAssignment or EffectiveGovernanceContext. |

The full assessment and evidence are in the
[FEATURE-0018 feature contract](../../features/FEATURE-0018-governance-iam-approval-exception-foundation.md)
and [approval evidence](../reuse-assessments/FEATURE-0018-approval-evidence.md).

## 5. Normative decision index

This table is navigation only; the linked F18-RD text is authoritative.

| Decision | Topic | Authority |
|---|---|---|
| F18-RD-01 | Current-feature-only authority | Appendix A |
| F18-RD-02 | Closed contract inventory and profiles | Appendix B |
| F18-RD-03 | Stable principal identity | Appendix A |
| F18-RD-04 | Direct principal and scoped AccessGroup assignments | Appendix A |
| F18-RD-05 | Membership is non-authorizing | Appendix A |
| F18-RD-06 | Distributed action ownership and central role composition | Appendix B |
| F18-RD-07 | Versioned RoleDefinition | Appendix A |
| F18-RD-08 | Scoped RoleAssignment | Appendix A |
| F18-RD-09 | Deterministic scoped authorization composition | Appendix A |
| F18-RD-10 | FEATURE-0017 adoption without reinterpretation | Appendix A |
| F18-RD-11 | Bounded FEATURE-0017 subject/target use | Appendix A |
| F18-RD-12 | Bounded ApprovalPolicy | Appendix A |
| F18-RD-13 | Immutable terminal approval evidence | Appendix A |
| F18-RD-14 | JIT privileged access | Appendix A |
| F18-RD-15 | Constrained break-glass access | Appendix A |
| F18-RD-16 | Snapshot-based AccessReview | Appendix A |
| F18-RD-17 | Bounded immutable exception evidence | Appendix A |
| F18-RD-18 | FEATURE-0018-limited GovernanceProfile v1 | Appendix B |
| F18-RD-19 | FEATURE-0013 adoption | Appendix B |
| F18-RD-20 | Audit before authorization-changing publication | Appendix B |
| F18-RD-21 | Deterministic validation and safe denial | Appendix A |
| F18-RD-22 | Deterministic in-memory foundation and local conformance | Appendix B |
| F18-RD-23 | Standards-validation gate | Appendix B |
| F18-RD-24 | Progressive and normally hidden user friction | Appendix A |

## 6. Security-critical consistency assertions

These assertions are sentinels against stale mirrors; they do not replace the
linked contracts.

| Boundary | Consistency assertion | Authority |
|---|---|---|
| Assignment writer | The RoleAssignment controller is the sole canonical spec, status, lifecycle, publication, revocation, and replacement writer; workflow controllers submit immutable intents only. | F18-RD-02, F18-RD-08 |
| Target binding | Every action uses exactly one registered `ExactResource`, `CreateParent`, or `ScopeOnly` variant; missing, wildcard, mixed, or ambiguous registrations fail closed. | F18-RD-06 |
| Membership | Organization-tree operations require current Organization Membership, exact CloudProvider operations require current CloudProvider Membership, and Platform/CloudPlatform operations require no Membership. Membership never grants an action. | F18-RD-05, F18-RD-09 |
| Group membership grant boundary | Creating, reactivating, restoring, or extending AccessGroup assignment effect derives the exact current `MembershipEnabledGrantEnvelope`; a non-empty envelope requires Membership authority, `roleassignment.grant`, and non-synthesized per-action scope/target/resource/validity witnesses. Standing effects require Standing witnesses; trusted provisioners have no bypass and privilege reduction remains available. Membership cannot create workflow eligibility. | F18-RD-04, F18-RD-05, F18-RD-09 |
| Human eligibility boundary | Approver, privileged-requester, and reviewer lists are non-empty OR lists of `EligibilityRef` containing exact Human PrincipalRefs only. RoleDefinition, RoleAssignment, AccessGroupRef, Membership, claim and other relationships never qualify; eligibility is separate from action authorization and rechecked at decision/effect publication. Policy/rule applicability plus the separate action alone determine scope and target reach. | F18-RD-12, F18-RD-14, F18-RD-16, F18-RD-18 |
| Standing access | Eligibility is deterministic from registered holder, scope, optional target, non-privileged role, responsibility, and review rules; no qualitative “narrow” test exists. | F18-RD-08 |
| Review evidence | Usage evidence is an embedded immutable `UsageEvidenceSummary`; FEATURE-0018 creates no usage-evidence resource. | F18-RD-16 |
| Review vocabulary | `ExpiredIncomplete` is a terminal campaign state; per-item dispositions are kind-specific `Retain`, `Revoke`, `Replace`, or `Stale`. | F18-RD-16 |
| Review independence | For AccessGroup-held authority, the conflict set includes the owner, direct Human members, and the responsible Human of every direct Workload/System member. Unresolved expansion fails closed for Retain/Replace; Revoke remains available. | F18-RD-16 |
| Definition lifecycle | `Superseded` is informational, not a lifecycle state. Published definitions remain version-pinned; suspension and retirement follow their exact F18-RD rules. | F18-RD-07, F18-RD-18 |
| Policy boundary | FEATURE-0017 outcomes are bounded inputs. Final authorization still evaluates exact grants, mandatory guardrails, evidence, dependencies, and safe-denial rules. | F18-RD-09 through F18-RD-11 |

## 7. User simplicity and journey evidence

The internal resource separation is normally hidden. Users provide business
intent; Sovrunn manages canonical references, versions, provenance,
concurrency, idempotency, and audit fields. Friction appears progressively only
when risk requires it.

| Journey | User-visible interaction | Scenario |
|---|---|---|
| Ordinary assignment | Choose holder, named role, scope, and validity. | F18-SCN-39 |
| Privileged request | Choose role, scope, duration, and justification; see one request state and expiry. | F18-SCN-40 |
| Approver decision | Review requester, requested access/control, scope, duration, risk context, quorum, and justification. | F18-SCN-41 |
| Access review | Review holder, role, scope, validity, qualified usage summary, and allowed disposition. | F18-SCN-42 |
| Exception request | Supply control, scope, duration, justification, and compensating controls; see request state and effective period. | F18-SCN-43 |

F18-SCN-01 through F18-SCN-49 are the current scenario inventory;
F18-SCN-38 is retired as a duplicate of F18-SCN-29. F18-SCN-44 through
F18-SCN-49 are security-closure scenarios, not additional user journeys.

## 8. Conflict reconciliation and architecture-leakage prevention

All 36 known baseline, dependency, digest, and matrix conflicts are enumerated
in the ADH conflict ledger and tracked in the
[architecture-leakage proof](FEATURE-0018-architecture-leakage-proof.md).
Reconciliation means the lower-precedence artifact points to or exactly agrees
with the controlling F18-RD; it does not promote this digest into a second
authority.

Leakage is prevented by these rules:

- no future-owned resource, field, action, workflow, adapter, or effect is
  introduced for implementation convenience;
- no provider-native principal, group, role, policy, permission, credential,
  or IAM effect becomes a FEATURE-0018 resource or output;
- no raw external claim grants access;
- no arbitrary condition language, nested/dynamic group model, provider SDK,
  external call, database, deployment topology, requirements, design, task, or
  implementation is introduced here; and
- every unresolved semantic conflict stops final approval rather than being
  delegated to downstream requirements or design.

## 9. Remaining delivery gates, not open architecture decisions

The approved recommendations leave no known FEATURE-0018 semantic choice for
requirements or design. Architecture, reuse, renewed security review and final
documentation/drift/hash checks have passed. Requirements generation remains a
separate human-gated next stage; design, tasks, implementation, executable
F18-RD-22 conformance and the FEATURE-0018 feature gate remain later gates.
