---
doc_type: feature
id: FEATURE-0018
title: Governance, IAM, Approval and Exception Foundation
status: implemented_and_merged
phase: 2R
reuse_assessment_format_version: 1.0.0
canonical_architecture: docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md
kiro_slug: governance-iam-approval-exception-foundation
controlling_handoff: ADH-2026-070
audit_taxonomy_correction_handoff: ADH-2026-073
accepted_decision: DEC-0060
ai_load_priority: feature
ai_summary: Implemented and merged through PR #20 as commit 2295e9e on 2026-09-07; approved architecture and handoffs remain controlling.
merged_pr: #20
merged_at: 2026-09-07
merged_commit: 2295e9e8ace09843ddb17bb8c43521bb3d828ea6
---

# FEATURE-0018: Governance, IAM, Approval and Exception Foundation

| Field | Value |
|---|---|
| Status | Architecture approved; fresh requirements authorization required after ADH-2026-073 |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase / order | Phase 2R / 8 |
| Direct dependencies | FEATURE-0012 and FEATURE-0017; adopting registrations in FEATURE-0013 and FEATURE-0016 |
| Sole architecture authority | `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md`, its content-bound ADH-2026-070 package, and approved ADH-2026-073 taxonomy correction |
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
| Related DEC / RFC / ADH references | DEC-0026, DEC-0036, DEC-0043, DEC-0050, DEC-0060; RFC-0012, RFC-0021, RFC-0022, RFC-0023; ADH-2026-070 package; ADH-2026-071 conformance reconciliation; ADH-2026-072 design-authority and audited-evaluation clarification; ADH-2026-073 terminal-exception audit-event correction |
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

### 6.1 Canonical requirement generation ledger

The rows below are subordinate generation-control mirrors of the 24 approved
architecture decision groups. They create no independent architecture or
requirement semantics. Each identifier maps one-to-one to its controlling
decision group and must not be repurposed, renumbered, merged, split or
paraphrased during requirements generation.

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

### 6.2 Canonical acceptance generation ledger

The rows below are subordinate generation-control mirrors of the approved
F18-SCN-01..49 enterprise scenario ledger. They create no independent scenario,
evidence or architecture authority. `AC-F18-38` preserves the retired duplicate
as excluded and must never be treated as an active acceptance obligation.

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

### 6.3 Canonical conformance adoption ledger

The four Slice-0 cases below are registered as FEATURE-0018-owned security
proof. Their exact inputs, expected state, error, side effects and gate remain
owned by `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` and
must be copied without reinterpretation into the requirements conformance
ledger. This adoption adds no public route, authentication implementation,
identity-provider integration or new error contract. In particular,
`VS0-CF-F01` proves rejection when an already-authenticated operation context is
absent; FEATURE-0018 still does not validate tokens or contact an identity
provider.

| Conformance ID | FEATURE-0018 adoption boundary |
|---|---|
| VS0-CF-F01 | Presence of the already-authenticated operation context; exact registry semantics only |
| VS0-CF-F02 | Authorization safe-denial behavior; exact registry semantics only |
| VS0-CF-X01 | Cross-Organization isolation; exact registry semantics only |
| VS0-CF-X02 | Cross-Project isolation; exact registry semantics only |

### 6.4 FEATURE-0018-local executable conformance ledger

ADH-2026-071 registers the exact local range below without changing any
F18-RD, REQ, AC, resource, action, state, writer, error, route or other runtime
semantic:

Active local registrations are `VS0-CF-F18-01..37` and
`VS0-CF-F18-39..54`.

```text
VS0-CF-F18-01..37  active; one-to-one with AC-F18-01..37
VS0-CF-F18-38      permanent excluded tombstone for AC-F18-38
VS0-CF-F18-39..49  active; one-to-one with AC-F18-39..49
VS0-CF-F18-50..54  active requirement-only architecture/contract proof
```

The exact `owner`, `inputs`, `expectedState`, `expectedError`,
`expectedSideEffects`, and `gate` fields are owned solely by
`docs/architecture/vertical-slices/VS-000-contract-registry.yaml`. Generated
requirements must copy every active FEATURE-0018-owned row exactly. Shared
`VS0-CF-F01`, `VS0-CF-F02`, `VS0-CF-X01`, and `VS0-CF-X02` remain supplementary;
FEATURE-0015-owned `VS0-CF-X03` never counts as FEATURE-0018-local proof.

### 6.5 Exact AC-to-conformance mapping

Each active AC maps to the active local conformance ID with the same numeric
suffix. `AC-F18-38` remains excluded and maps to no active case. Generated
stage artifacts must render all 48 active mappings as individual rows.

### 6.6 Exact REQ-to-conformance mapping

| Requirement | Exact FEATURE-0018-local conformance IDs |
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
| REQ-F18-22 | VS0-CF-F18-01..37, VS0-CF-F18-39..53 |
| REQ-F18-23 | VS0-CF-F18-54 |
| REQ-F18-24 | VS0-CF-F18-29, VS0-CF-F18-39, VS0-CF-F18-40, VS0-CF-F18-41, VS0-CF-F18-42, VS0-CF-F18-43 |
