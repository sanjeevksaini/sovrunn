---
doc_type: feature
id: FEATURE-0018
title: Governance, IAM, Approval and Exception Foundation
status: architecture-approved
phase: 2R
reuse_assessment_format_version: 1.0.0
canonical_architecture: docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md
kiro_slug: governance-iam-approval-exception-foundation
controlling_handoff: ADH-2026-070
accepted_decision: DEC-0060
ai_load_priority: feature
ai_summary: Approved scope, reuse and generation-control contract for FEATURE-0018; requirements, design, tasks and implementation remain separately gated.
---

# FEATURE-0018: Governance, IAM, Approval and Exception Foundation

| Field | Value |
|---|---|
| Status | Architecture approved; requirements authorization pending |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase / order | Phase 2R / 8 |
| Direct dependencies | FEATURE-0012 and FEATURE-0017; adopting registrations in FEATURE-0013 and FEATURE-0016 |
| Sole architecture authority | `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md` plus its content-bound ADH-2026-070 package |
| Approved change records | ACR-2026-002 and DEC-0060 |
| Public routes / persistence / real adapters | None authorized at architecture stage |

## 1. Purpose and authority boundary

FEATURE-0018 supplies Sovrunn's provider-neutral control-plane authorization,
approval, privileged-access, access-review and immutable exception-evidence
foundation. It supports Human, Workload and System principals, direct
PrincipalRef assignments and directly provisioned AccessGroup assignments.

This file is the feature-scope, reuse-assessment and generation-control
contract. It does not restate the 24 semantic decision groups. Any conflict
with the canonical architecture or ADH package stops work with
`ARCHITECTURE_DECISION_REQUIRED`; requirements, design and tasks may translate
approved behavior but may not repair or extend architecture.

## 2. Owned scope

The closed persistent inventory is `AccessGroup`, `Membership`,
`RoleDefinition`, `RoleAssignment`, `PrivilegedAccessRequest`, `AccessReview`,
`ApprovalPolicy`, `ApprovalRequest`, `ExceptionGrant`, and the FEATURE-0018 v1
`GovernanceProfile`. `PrincipalRef` and the supporting contracts listed by
F18-RD-02 are embedded or transient values, not independently addressable
resources.

FEATURE-0018 owns deterministic validation, authorization composition, exact
action-target consumption, single-writer materialization, evidence linkage,
fixed-time in-memory fixtures and feature-local conformance for those contracts.

## 3. Dependency and exclusion boundaries

- FEATURE-0012 retains resource grammar, seven ScopeKinds, reference, safe-
  denial, validation, redaction, concurrency, idempotency and Problem authority.
- FEATURE-0013 retains DecisionRecord, DecisionProfile and AuditEvent envelopes;
  FEATURE-0018 only registers three profiles and its event taxonomy.
- FEATURE-0016 retains the meanings and lifecycle of ExecutionTarget actions;
  FEATURE-0018 only consumes their exact registrations.
- FEATURE-0017 retains the engine-neutral policy seam and fake; policy results
  constrain but never create permission.
- FEATURE-0020 retains ProfileAssignment, inheritance, conflict/exception
  application and EffectiveGovernanceContext resolution.
- Real IdP/SCIM, policy/workflow engines, provider-native IAM, credentials,
  persistence, external execution, UI, service/placement/sovereignty behavior
  and every later-feature output are excluded.

## 4. Feature-level reuse summary

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 resource/scope/reference/validation foundation | Reuse | Existing common contracts provide the canonical shared behavior. | Approved | DEC-0026; ADH-2026-012; ADH-2026-070 |
| FEATURE-0013 DecisionRecord/DecisionProfile/AuditEvent envelopes | Reuse | One evidence-envelope owner prevents competing audit and decision authority. | Approved | DEC-0043; ADH-2026-017; F18-RD-19/20 |
| FEATURE-0016 ExecutionTarget action meanings and lifecycle | Reuse | Existing actions are consumed unchanged with owner-supplied binding/class metadata. | Approved | ADH-2026-058..066; F18-RD-06 |
| FEATURE-0017 evaluation seam, result vocabulary, mapper and fake | Reuse | Existing engine-neutral evaluation is sufficient within its bounded v1 subject/target fit. | Approved | DEC-0028; DEC-0036; ADH-2026-067..069 |
| Canonical FEATURE-0018 resource/value contracts and authorization composition | Build | No external model owns Sovrunn's complete provider-neutral canonical semantics, writers and evidence boundary. | Approved | ACR-2026-002; DEC-0060; F18-RD-01..21 |
| Deterministic FEATURE-0018 controllers/fixtures/conformance fakes | Build | Side-effect-free Phase 2R proof must exercise the same Sovrunn-owned boundaries without external systems. | Approved | F18-RD-22 |
| GovernanceProfile v1 and progressive-disclosure projections | Build | Sovrunn owns its bounded governance envelope and simple user journeys; FEATURE-0020 retains resolution. | Approved | DEC-0050; F18-RD-18/24 |
| OIDC/SCIM/WebAuthn/SPIFFE compatibility | Reuse | Standards are used for paper compatibility only; no runtime integration is selected. | Deferred | F18-RD-23 |
| Real IdP, policy-engine or workflow integration | Wrap | A future owner may wrap a selected mature system behind an approved adapter after fresh assessment. | Deferred | DEC-0036; RFC-0012; F18-RD-01/23 |

## 5. Capability assessment: canonical authorization and evidence foundation

### Identity

| Field | Value |
|---|---|
| Feature identity | FEATURE-0018 |
| Capability or decision-unit identity | Provider-neutral governance, IAM, approval, review, exception contracts and deterministic conformance foundation |
| Assessment owner | Sovrunn Architecture Owner |

### Classification

| Field | Value |
|---|---|
| Disposition | Build |
| Decision status | Approved |

### Analysis

| Field | Value |
|---|---|
| Assessment scope | Closed F18-RD-02 inventory, authorization composition, lifecycle, writers, bounded workflows, evidence adoption, fixtures and local proof only. |
| Candidate category | Enterprise IAM/RBAC/IGA/PAM, approval/workflow, policy and cloud-provider authorization foundations. |
| Mature candidates / applicable standards | NIST SP 800-207, 800-53 Rev. 5, 800-63-4 and 800-162; CIS Controls 5/6/8; CSA CCM 4.1; OWASP ASVS 5.0; ISO/IEC 27001/27002 concepts; AWS IAM/Organizations/Identity Center patterns; OIDC, SCIM, WebAuthn and SPIFFE compatibility. |
| Relevant candidate strengths | Mature frameworks establish zero trust, least privilege, separation of duties, lifecycle, strong authentication, audit, review and shared-responsibility practices; protocols provide portable future integration boundaries. |
| Material candidate constraints | No candidate owns Sovrunn's canonical seven-scope model, action vocabulary, resource profiles, user projections, dependency ownership, provider-neutral/native-IAM intersection or Phase 2R deterministic proof. Cloud-native IAM models are provider-specific; standards are controls, not a complete product contract. |
| Rationale | Reuse the frameworks and prior-feature primitives, but build the smallest Sovrunn-owned domain layer needed to make authorization deterministic, portable, auditable and simple for users. |
| Selected foundation or approach | Scoped RBAC with direct-member access groups; grant union plus guardrail intersection; explicit validity; bounded approval/PAM/review/exception resources; immutable provenance; audit-before-publication; progressive disclosure. |

### Boundary

| Field | Value |
|---|---|
| Sovrunn-owned responsibility | Canonical FEATURE-0018 contracts, actions owned by FEATURE-0018, assignment/authorization semantics, single writers, evidence linkage, projections, deterministic fixtures and conformance. |
| Reused or external responsibility | FEATURE-0012 grammar; FEATURE-0013 evidence envelopes; FEATURE-0016 ExecutionTarget actions; FEATURE-0017 evaluation seam; external authentication lifecycle; provider-native IAM; standards and future protocol adapters. |
| Data crossing the boundary | Already-authenticated PrincipalRef and assurance evidence enter; exact canonical refs, preselected policy/evidence and transient evaluation evidence cross internal seams; protected/redacted decision and audit evidence exits. No raw token, credential or provider-native role crosses the canonical contract. |
| Control crossing the boundary | FEATURE-0018 authorizes a canonical control-plane operation. A later execution owner independently invokes a least-privileged native identity; native IAM separately allows or denies. No FEATURE-0018 controller performs external effects. |
| Adapter required | No |
| Adapter rationale | Phase 2R builds only in-process contracts and deterministic fakes. Every real external integration is deferred and requires its own approved adapter assessment. |
| Adapter or contract identifier | none |
| Vendor-native types allowed | No |

### Suitability

| Field | Value |
|---|---|
| Sovereignty and deployment fit | Provider-neutral and fully in-process in Phase 2R; deterministic proof works disconnected/air-gapped. Future external systems remain replaceable adapters. |
| Security and trust | Default deny, least privilege, exact scopes/targets, strong assurance, separation of duties, immutable provenance, time-bounded privilege, review/revocation and protected audit. |
| Operational and supportability | Explicit owners/writers/lifecycles, fixed reason/evidence contracts, exact clocks, safe projections, idempotency and race proof support diagnosis without exposing sensitive internals. |
| Licensing and supply-chain | No production third-party runtime is selected or added. Standards/framework mappings are documentary; formal ISO conformity is not claimed. |
| Portability and provider-neutrality impact | Canonical resources contain no IdP-, AWS-, Azure-, GCP-, Kubernetes- or other provider-native IAM type. Native authorization remains an independent second plane. |

### Phase and scope

| Field | Value |
|---|---|
| Allowed in current phase | Yes |
| Current-phase work | Contracts, validation, deterministic in-memory controllers/fixtures, positive/negative/race conformance and evidence only. |
| Deferred work | Real IdP/group provisioning, policy/workflow engine, provider-native IAM/credential integration, persistence, routes/UI, production operations and later-feature governance resolution. |
| Explicit non-goals | No identity lifecycle, raw-token validation, external group-claim authority, nested/dynamic groups, generic ABAC, provider-native role/policy/credential, real adapter, network/database call, service execution, future feature or deployment topology. |
| Exit or migration boundary | A future integration requires a new assessment/ADH and must remain behind canonical contracts. New ScopeKind/containment, group form, action meaning or condition language requires core change control. |
| Phase 2 non-goal acknowledgement | Phase 2R remains side-effect-free; no real provider, Kubernetes, IdP, policy/workflow or infrastructure operation is authorized. |

### Build justification

| Field | Value |
|---|---|
| Why Reuse is insufficient | Standards and provider products do not supply Sovrunn's canonical resources, seven-scope semantics, action registry, writer boundaries or user projections. |
| Why Wrap is insufficient | There is no single selected runtime to wrap, and making one external IAM/workflow product canonical would violate portability and phase scope. |
| Why Extend is insufficient | Extending FEATURE-0012/0013/0016/0017 would move IAM-domain semantics into the wrong owner and couple unrelated shared foundations. |
| Protected Sovrunn differentiation and long-term ownership | Sovrunn owns the unified provider-neutral PaaS control-plane authorization, governance evidence and simple cross-provider user experience. |

### Risk mitigation

#### Applicable architecture risks

| Risk | Preventive controls | Detection controls | Corrective path |
|---|---|---|---|
| Privilege escalation or authorization ambiguity | Closed action/target/scope registries, per-action resource-reach-aware grantor ceiling with no cross-assignment synthesis, the same exact dominance for AccessGroup Membership assignment-effect expansion, and exact-Human PrincipalRef-only `EligibilityRef` semantics for approval/review/privileged flows. | Over-action/scope/resource/duration/delegation negatives, empty/non-empty group envelopes, rejection of RoleDefinition/RoleAssignment/AccessGroup/Membership/claim-derived eligibility, exact reviewer scope/target/action/SoD cases, JIT-to-Standing laundering, trusted-provisioner overreach, membership/assignment/rule races, registry drift and security review. | Deny or conflict/retry without action or eligibility publication; privilege-reducing Membership transitions remain available. |
| Stale identity/group/guest/workload authority | Explicit provenance/freshness/expiry/responsibility, direct membership only, review/revocation. | Freshness, expiry, inactive-owner/responsible-party and review-due proof. | Suspend/revoke/recover through owning controller; never trust raw claims. |
| Multiple writers or partial evidence publication | Sole controllers, immutable intents and audit-before-publication atomic boundaries. | Writer ledger, failure injection, idempotency and race proof. | Publish nothing, retry the same authorized operation where permitted, and reconcile only through sole writer. |
| Privileged assurance or review self-dealing | Uniform phishing-resistant AAL2-or-higher activation floor and beneficiary conflict set that expands AccessGroup-held authority through each direct Workload/System member to its responsible Human at atomic review decision/effect publication. | JIT/break-glass assurance negatives plus direct, group, group-workload-responsible-Human, unresolved-expansion and replacement-union self-certification cases. | Deny activation, return `REVIEWER_CONFLICT`, or fail closed with `REVIEW_BENEFICIARY_UNRESOLVED`; never publish Retain/Replace authority effect or review-due advancement, while independently authorized Revoke remains available. |
| Provider-native IAM leakage | Canonical/native intersection boundary; no native types/credentials/effects. | Import/schema/string sentinels and two-plane truth table. | Remove leakage; require separate adapter architecture and security review. |
| User complexity | Progressive disclosure, curated templates and six correlated journeys; no facade resource. | Journey evidence and usability review. | Simplify projections/templates without collapsing protected contracts or writers. |
| Future-feature leakage | Closed exclusions and FEATURE-0020 resolution boundary. | Resource/field/dependency ledgers and architecture-drift checks. | Remove future semantics or process a separate ACR/DEC. |

- Residual risk: Medium until independent security review and executable proof;
  expected Low-to-Medium after feature-gate evidence.
- Replacement risk: Medium
- Reassessment triggers: new ScopeKind/containment; action meaning; group or
  condition model; real adapter/runtime/dependency; provider-native type;
  persistence/public route; new assurance/approval/review/exception semantics;
  standards major-version change; material security finding.

### Traceability

| Field | Value |
|---|---|
| Related DEC / RFC / ADH references | DEC-0026, DEC-0036, DEC-0043, DEC-0050, DEC-0060; RFC-0012, RFC-0021, RFC-0022, RFC-0023; ADH-2026-070 package |
| Linked acceptance criteria | F18-IVM-A01..I06, F18-SCN-01..49 with F18-SCN-38 retired, F18-RD-22 conformance inventory and the 36-item reconciliation ledger |
| Validation and review evidence | FEATURE-0018 digest, industry matrix, leakage proof, standards mapping, passing independent security review renewal-05, strict docs/drift checks and later feature gate |

### Human-approval evidence

| Field | Value |
|---|---|
| Structured approval-evidence record | `docs/reviews/reuse-assessments/FEATURE-0018-approval-evidence.md` |
| Decision status | Approved |
| Approving person or role | Sanjeev Kumar, Sovrunn Architecture Owner |
| Approval date | 2026-08-30 |

## 6. Generation boundary

Architecture approval does not automatically start a Kiro specification or
implementation stage. Once requirements generation is separately authorized,
requirements may derive verifiable statements only from F18-RD-01..24 and the
approved evidence ledgers. Any missing value is a blocker, not permission to
invent architecture.
