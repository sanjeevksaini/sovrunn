---
doc_type: architecture_validation_matrix
feature: FEATURE-0018
title: Governance, IAM, Approval and Exception Foundation Industry Validation Matrix
status: architecture-validated-and-approved
baseline: ARCH-2026.08-PHASE2R-CANONICAL
prepared: 2026-08-27
updated: 2026-08-30
companion_digest: docs/reviews/architecture-readiness/FEATURE-0018-architecture-digest.md
---

# FEATURE-0018 Industry Validation Matrix

> **Approved closure record — not architecture authority or a conformity
> claim.** This matrix records the architecture owner's 2026-08-28 decisions
> and validates the approved-for-handoff FEATURE-0018 model against
> architecture-quality, identity, access-management, zero-trust, cloud-control,
> and application-security references. The approval does not authorize
> requirements, design, tasks, or implementation. Formal ISO clause-level
> assessment requires authorized access to the complete applicable standards.
> Kiro may apply only a subsequently approved Architecture Decision Handoff.

## 1. Purpose

This matrix evaluates whether the canonical FEATURE-0018 output inventory is:

- **simple** — every resource has one necessary responsibility and no duplicate
  authority;
- **complete** — every FEATURE-0018-owned enterprise IAM, approval, review, and
  exception concern is covered or recorded as a gap;
- **accurate** — terminology, scope, lifecycle, security behavior, and
  dependency use agree with canonical authority and applicable standards; and
- **evidence-ready** — every material invariant can produce positive,
  negative, expiry/revocation, concurrency, redaction, and failure proof.

Industry standards do not prescribe a Sovrunn resource count. They define
concerns, outcomes, controls, terminology, and assessment practices. This
matrix maps those concerns to the approved-for-handoff resource model and tests
whether a resource is missing, unnecessary, duplicated, or leaking another
feature's authority.

## 2. Assessment boundary

### 2.1 In scope

- `PrincipalRef` use in FEATURE-0018;
- `AccessGroup`;
- `Membership`;
- `RoleDefinition`;
- `RoleAssignment`;
- `PrivilegedAccessRequest`;
- `AccessReview`;
- `GovernanceProfile`;
- `ApprovalPolicy`;
- `ApprovalRequest`;
- `ExceptionGrant`;
- scoped authorization composition;
- FEATURE-0013 decision/audit adoption;
- FEATURE-0017 policy-evaluation adoption;
- deterministic in-memory fixtures and FEATURE-0018-local conformance; and
- provider-neutral future adapter compatibility.

### 2.2 Out of scope

- identity-provider selection or integration;
- OIDC token verification and SCIM synchronization;
- real policy or workflow engines;
- provider-native IAM and credentials;
- profile assignment and `EffectiveGovernanceContext` resolution;
- sovereignty, placement, entitlement, quota, billing, execution, and AI
  explanation; and
- certification of any implementation that does not yet exist.

An out-of-scope concern is valid only when its owner or explicit deferral is
recorded. A FEATURE-0018-owned security gap cannot pass by being labelled
future work.

## 3. Standards and assurance register

| Reference | Version used | Validation role | Authority or link |
|---|---|---|---|
| ISO/IEC/IEEE 42010 | 2022 | Architecture-description integrity: stakeholders, concerns, viewpoints, relationships, and decision traceability | [ISO 42010:2022](https://www.iso.org/standard/74393.html) |
| ISO/IEC 25010 | 2023 | Product-quality lens for functional suitability, interaction, maintainability, reliability, security, and acceptance criteria | [ISO 25010:2023](https://www.iso.org/standard/78176.html) |
| ISO/IEC 24760 series | 2025 | Identity concepts, terminology, reference architecture, requirements, and identity lifecycle | [ISO 24760 package](https://www.iso.org/publication/PUB200247.html) |
| ISO/IEC 29146 | 2024 | Access-management concepts, accountability, components, and distributed access boundary | [ISO 29146:2024](https://www.iso.org/standard/86013.html) |
| NIST Cybersecurity Framework | 2.0 | Enterprise risk outcomes, governance, identity/access, protection, detection, response, and recovery context | [NIST CSF 2.0](https://www.nist.gov/publications/nist-cybersecurity-framework-csf-20) |
| NIST SP 800-37 | Rev. 2 | Risk-management process, control selection, authorization, assessment, and monitoring | [NIST SP 800-37 Rev. 2](https://csrc.nist.gov/pubs/sp/800/37/r2/final) |
| NIST SP 800-53 | Rev. 5, current NIST release/errata at assessment time | Access Control, Identification and Authentication, Audit and Accountability, Assessment/Authorization/Monitoring, Personnel Security, and Risk Assessment | [NIST SP 800-53 Rev. 5](https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final) |
| NIST SP 800-53A | Rev. 5, current NIST assessment release at assessment time | Security/privacy control-assessment procedures and evidence method | [NIST SP 800-53A Rev. 5](https://csrc.nist.gov/pubs/sp/800/53/a/r5/final) |
| NIST SP 800-63 suite | Revision 4, 2025 | Identity risk, identity proofing, authentication assurance, federation, privacy, customer experience, and redress | [NIST SP 800-63-4](https://csrc.nist.gov/pubs/sp/800/63/4/final) |
| NIST SP 800-207 | 2020 | Zero-trust resource-centric authorization and removal of implicit trust | [NIST SP 800-207](https://csrc.nist.gov/pubs/sp/800/207/final) |
| NIST SP 800-207A | 2023 | Identity-tier access control for cloud-native and multi-cloud applications | [NIST SP 800-207A](https://csrc.nist.gov/pubs/sp/800/207/a/final) |
| NIST SP 800-162 | Updated 2019 | Attribute-based access-control subject, object, action, environment, and policy semantics | [NIST SP 800-162](https://csrc.nist.gov/pubs/sp/800/162/upd2/final) |
| AWS IAM, Azure RBAC, and Google Cloud IAM | Current official documentation at assessment time | Comparative CSP validation of principal/role/scope assignment, grant union, restrictive boundaries, deny precedence, temporary access, and access analysis without importing provider-native schemas | [AWS evaluation logic](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_evaluation-logic.html), [Azure role assignments](https://learn.microsoft.com/en-us/azure/role-based-access-control/role-assignments), [Google principal access boundaries](https://docs.cloud.google.com/iam/docs/principal-access-boundary-policies), [AWS IAM Access Analyzer](https://docs.aws.amazon.com/IAM/latest/UserGuide/what-is-access-analyzer.html) |
| CSA Cloud Controls Matrix | 4.1, 2026 | Cloud IAM, governance, logging, interoperability, and shared-responsibility assurance | [CSA CCM 4.1](https://cloudsecurityalliance.org/artifacts/cloud-controls-matrix-v4-1) |
| CIS Critical Security Controls | 8.1 | Operational account, access-control, and audit-log cross-check | [CIS Controls v8.1](https://www.cisecurity.org/controls/v8-1) |
| OWASP ASVS | 5.0.0 | Later API/application verification for authentication, access control, validation, and logging | [OWASP ASVS](https://owasp.org/www-project-application-security-verification-standard/) |
| OIDC, OAuth Security BCP, SCIM, WebAuthn, SPIFFE | Current approved standards at adapter assessment time | Future adapter compatibility only; not FEATURE-0018 implementation scope | [OAuth Security BCP](https://www.rfc-editor.org/info/rfc9700/), [SCIM Core Schema](https://www.rfc-editor.org/info/rfc7643/) |

ISO references in this matrix are validation headings rather than assertions of
clause-level conformity. A formal conformity review must use licensed/current
texts and record exact clauses in the approval evidence.

## 4. Status vocabulary

| Status | Meaning |
|---|---|
| `INHERITED` | Closed by a higher-precedence approved Sovrunn authority; FEATURE-0018 must preserve it |
| `APPROVED_FOR_HANDOFF` | Historical pre-application state: explicitly dispositioned by the architecture owner but not yet canonically applied |
| `PROPOSED` | Recommended in the companion digest; human disposition and canonical application are pending |
| `GAP` | Material ambiguity or missing contract that must be resolved before architecture closure |
| `DEPENDENCY` | Satisfied only through an approved dependency contract that FEATURE-0018 consumes unchanged |
| `EXCLUDED` | Intentionally outside FEATURE-0018 with a named owner or explicit deferral |
| `NOT_ASSESSED` | Requires complete standards text, scenario review, or conformance evidence not yet available |
| `EVIDENCE_PREPARED` | Historical pre-approval state: named architecture evidence existed but independent and/or final human review was pending |
| `ARCHITECTURE_PASS` | Canonical application, required architecture evidence, independent security review and final human approval are complete; executable proof remains a later feature-gate obligation |
| `PASS` | Approved architecture and required evidence both satisfy the row |
| `FAIL` | Approved scope owns the concern but the architecture or evidence does not satisfy it |

No row becomes `PASS` from prose alone. Approval and the specified evidence are
both required.

## 5. Gate A — Architecture integrity and simplicity

| ID | Validation concern | Reference | Resource or boundary | Pass condition | Required evidence | Current status |
|---|---|---|---|---|---|---|
| F18-IVM-A01 | Every resource addresses a named stakeholder concern | ISO 42010 | All FEATURE-0018 resources | Each resource has purpose, stakeholders, owner, scope, writers, lifecycle, relationships, and exclusions | Approved contract responsibility ledger | `ARCHITECTURE_PASS` |
| F18-IVM-A02 | Every resource is necessary | ISO 42010; ISO 25010 | All FEATURE-0018 resources | Removing the resource would lose independent identity, lifecycle, authorization, or retained evidence; separate contracts do not imply separate user journeys | Approved resource remove/merge analysis and later contract ledger | `ARCHITECTURE_PASS` |
| F18-IVM-A03 | No two resources own the same responsibility | ISO 42010; ISO 25010 | All FEATURE-0018 resources | Access grouping, belonging, role definition, grant, elevation, review, policy, approval transaction, profile, and exception evidence remain distinct while projections compose them for users | Responsibility-overlap and projection matrix | `ARCHITECTURE_PASS` |
| F18-IVM-A04 | One authoritative scope source | ISO 29146; NIST SP 800-53 AC | Scoped resources | `metadata.scopeRef` is the sole ownership scope; relationship refs cannot become competing authorities | Scope-source ledger and mismatch tests | `INHERITED` |
| F18-IVM-A05 | One writer per system-owned field | ISO 42010; NIST SP 800-53 AC/AU | Managed and LRO resources | Every system-owned field has one named writer; schedulers invoke that authority rather than becoming writers | Writer ledger and wrong-writer proof | `ARCHITECTURE_PASS` |
| F18-IVM-A06 | Status does not duplicate history | ISO 25010; NIST SP 800-53 AU | Membership, RoleAssignment, PAR, AccessReview, ApprovalRequest | Status contains current facts; DecisionRecord/AuditEvent retain history | Schema review and no-duplicate-history proof | `ARCHITECTURE_PASS` |
| F18-IVM-A07 | No generic container resource | ISO 25010 | GovernanceProfile, ApprovalRequest | Typed, bounded fields only; ApprovalRequest accepts only PAR ref, RoleAssignmentProposal, or ExceptionProposal; GovernanceProfile activates only its approved v1 components | Field-by-field necessity ledger and unknown-payload proof | `ARCHITECTURE_PASS` |
| F18-IVM-A08 | No unapproved resource kinds | ISO 42010; Sovrunn architecture rules | AccessGroup, ExceptionRequest, workflow resources | AccessGroup is an approved-for-handoff canonical extension; no ExceptionRequest, facade, or workflow kind is added | Schema/kind inventory and applied-change check | `ARCHITECTURE_PASS` |
| F18-IVM-A09 | No future-owned fields are activated early | ISO 42010; Sovrunn leakage rules | All resources | Future-owned fields are rejected rather than stored, defaulted, validated, or projected | Unknown/deferred-field negative tests | `ARCHITECTURE_PASS` |
| F18-IVM-A10 | User complexity is proportional to risk | ISO 25010; NIST SP 800-63-4 customer experience | User-facing projections | Routine operations are silent; ordinary assignment asks who/role/scope/validity; privileged, review, and exception flows expose only actionable business information | Six approved architecture journeys plus later executable projection/usability evidence | `ARCHITECTURE_PASS` |
| F18-IVM-A11 | Internal provenance and evidence are not manual user inputs | ISO 25010; NIST SP 800-63-4 customer experience; NIST SP 800-53 AU | User-facing projections + internal contracts | Issuer/subject, exact versions, provenance timestamps/object refs, AuthorizationInput/Result, usage evidence, decision/audit refs, and correlation IDs are system-managed; only non-derivable business inputs are requested | Projection field ledger, forged-system-field negatives, and six journey cases | `ARCHITECTURE_PASS` |

## 6. Gate B — Identity and Membership completeness

| ID | Validation concern | Reference | Resource or boundary | Pass condition | Required evidence | Current status |
|---|---|---|---|---|---|---|
| F18-IVM-B01 | Stable identity key | ISO 24760; NIST SP 800-63-4 | PrincipalRef | Issuer and subject are durable; email/display name are not identity authority | Stable-identity and changed-email scenarios | `INHERITED` |
| F18-IVM-B02 | Identity and authentication remain separate from authorization | ISO 24760; NIST SP 800-63-4; NIST SP 800-207 | PrincipalRef boundary | FEATURE-0018 accepts authenticated identity and provider-neutral assurance evidence but does not authenticate tokens or retain raw factors | Boundary review and zero-token-processing proof | `ARCHITECTURE_PASS` |
| F18-IVM-B03 | Human, Workload, and System identities | ISO 24760; NIST SP 800-207A | PrincipalRef | Closed categories support human and non-human callers without provider-native identity types | One scenario and negative case per principal category | `INHERITED` |
| F18-IVM-B04 | Membership is non-authorizing | ISO 29146; NIST SP 800-207 | Membership | Organization, CloudProvider, or AccessGroup membership alone grants no action and cannot substitute for RoleAssignment | Membership-only denial proof for every context | `ARCHITECTURE_PASS` |
| F18-IVM-B05 | Membership context is bounded | ISO 24760; ISO 29146 | Membership | CloudProvider, Organization, and approved AccessGroup contexts are independent; AccessGroup is scoped to one Organization or CloudProvider | Scope matrix, applied-change record, and invalid-scope tests | `ARCHITECTURE_PASS` |
| F18-IVM-B06 | Join/suspend/reactivate/revoke lifecycle | ISO 24760 lifecycle; NIST SP 800-53 AC/PS | Membership | Lifecycle, guest expiry, synchronized freshUntil, and terminal behavior are deterministic and auditable | Positive, invalid-transition, expiry, freshness, and concurrency cases | `ARCHITECTURE_PASS` |
| F18-IVM-B07 | Suspension has immediate authorization effect | NIST SP 800-207; NIST SP 800-53 AC | Membership + RoleAssignment | Dependent assignments become ineffective without rewriting historical grants | Authorization-before/after suspension proof | `ARCHITECTURE_PASS` |
| F18-IVM-B08 | Reactivation does not revive expired/revoked grants | ISO 29146; least privilege | Membership + RoleAssignment | Reactivation re-enables only independently valid active grants | Suspension/reactivation/expiry sequence proof | `ARCHITECTURE_PASS` |
| F18-IVM-B09 | Multiple memberships remain independently scoped | ISO 24760; ISO 29146 | Membership | Multiple memberships never union identity, scope, or authority | Cross-Organization and cross-CloudProvider cases | `ARCHITECTURE_PASS` |
| F18-IVM-B10 | Enterprise group behavior is explicit | ISO 24760; SCIM compatibility | AccessGroup/Membership/RoleAssignment boundary | Direct PrincipalRef and scoped AccessGroup assignments are explicit; only directly provisioned canonical membership contributes; raw IdP claims, nested groups, and dynamic groups grant nothing | Applied AccessGroup change, direct/group allow-deny pairs, and raw-claim denial | `ARCHITECTURE_PASS` |
| F18-IVM-B11 | Guest/contractor membership is bounded | ISO 24760; NIST SP 800-63-4 federation | Membership | Guest requires sponsor and expiresAt; assignments remain TimeBound; approval responsibility and audit attribution are exact | Guest-access scenario suite | `ARCHITECTURE_PASS` |
| F18-IVM-B12 | Personal data is minimized and redacted | ISO 24760 privacy; NIST SP 800-63-4 | PrincipalRef/Membership projections | Raw assertions/factors are absent; system provenance is protected; customer-safe projections contain only required identity display | Data-classification and redaction review | `ARCHITECTURE_PASS` |
| F18-IVM-B13 | Membership provenance and responsible ownership are explicit | NIST SP 800-53 AC-2; CIS Controls 5/6; ISO 24760 | Membership | Source authority/object, provisioning and synchronization time are retained where applicable; guest, Workload, and System memberships identify a responsible party; stale provenance grants nothing silently | Provenance freshness, missing-owner, suspension, and redaction cases | `ARCHITECTURE_PASS` |
| F18-IVM-B14 | AccessGroup assignment-effect expansion cannot bypass delegated-grant reach | NIST SP 800-53 AC-2/AC-3/AC-6; NIST SP 800-207; CIS Control 6 | Membership + AccessGroup + RoleAssignment | Creation, reactivation, or freshness/provenance extension derives the exact current `MembershipEnabledGrantEnvelope`; non-empty effect requires applicable Membership authority, `roleassignment.grant`, and one non-synthesized per-action scope/target/resource/validity witness; Standing effect requires Standing witnesses; trusted provisioner has no bypass and privilege reduction remains available | Empty/non-empty envelope, resource/scope/time mismatch, JIT-to-Standing laundering, trusted-provisioner overreach, stale refresh, self/accomplice add, and membership/assignment race cases | `ARCHITECTURE_PASS` |
| F18-IVM-B15 | Membership cannot manufacture workflow or privileged eligibility | NIST SP 800-53 AC-2/AC-5/AC-6; NIST SP 800-207; CIS Control 6 | EligibilityRef + Membership + approval/review/privileged rules | Approver, requester and reviewer eligibility accept exact Human PrincipalRefs only; RoleDefinition, RoleAssignment, AccessGroupRef, Membership, claim and other relationships reject; eligibility remains separate from action authorization and is rechecked at decision/effect publication | Empty-group policy/rule, self/accomplice, unsupported-kind, scope/target/action, and publication-race cases | `ARCHITECTURE_PASS` |
| F18-IVM-B16 | Role publication or assignment cannot manufacture workflow or privileged eligibility | NIST SP 800-53 AC-3/AC-5/AC-6; NIST SP 800-207; CIS Control 6 | EligibilityRef + RoleDefinition + RoleAssignment + governance rules | RoleDefinition and RoleAssignment are structurally invalid EligibilityRef kinds; assigning a role or later publishing a policy/rule cannot qualify existing holders; policy/rule applicability and separately authorized canonical action alone determine scope and target reach | Self/accomplice ordinary-role assignment, later-rule reference, broad existing-holder, normal-JIT, break-glass, approval, reviewer and concurrent publication cases | `ARCHITECTURE_PASS` |

## 7. Gate C — Roles and scoped authorization completeness

| ID | Validation concern | Reference | Resource or boundary | Pass condition | Required evidence | Current status |
|---|---|---|---|---|---|---|
| F18-IVM-C01 | Canonical action ownership | ISO 29146; NIST SP 800-53 AC | RoleDefinition/action registry | Resource-owning feature defines action meaning; FEATURE-0018 validates and composes without reinterpretation | Registered-action provenance check | `ARCHITECTURE_PASS` |
| F18-IVM-C02 | Provider-native permissions do not become canonical | ISO 29146; CSA CCM interoperability | RoleDefinition | AWS/Azure/OpenStack/Kubernetes/IdP roles and claims are absent from canonical roles | Schema/import/string sentinel checks | `INHERITED` |
| F18-IVM-C03 | Role publication is immutable and version-pinned | ISO 29146; NIST SP 800-53 CM/AC | RoleDefinition | Draft is mutable; published version immutable; changes create new version; retained refs remain resolvable | Publish/mutate/retire/reference tests | `INHERITED` |
| F18-IVM-C04 | RoleAssignment is the grant | ISO 29146; NIST SP 800-53 AC-3/AC-6 | RoleAssignment | Assignment binds exactly one `roleHolderRef`, pinned role version, metadata.scopeRef, optional exact narrowing resourceRef, and explicit validity | Direct/group holder, resource-specific, invalid-holder, and absence-of-grant tests | `ARCHITECTURE_PASS` |
| F18-IVM-C05 | Grantor cannot exceed delegated authority | ISO 29146; NIST SP 800-53 AC-5/AC-6 | RoleAssignment | Grantor has roleassignment.grant over the complete proposed reach and one independently applicable delegable witness per action covering scope, target-binding mode, exact resource reach, validity and delegation; restrictions from different assignments never synthesize broader reach | Wrong-admin, overbroad action/scope/duration/delegation, resource-A-to-B and resource-to-scope denial tests | `ARCHITECTURE_PASS` |
| F18-IVM-C06 | Scope matching is deterministic | ISO 29146; NIST SP 800-162 | RoleAssignment | Organization→OrganizationUnit→Tenant→Project is the inherited containment chain; Platform, CloudPlatform, and CloudProvider do not implicitly contain one another; undefined containment is exact | Seven-scope applicability and cross-tree denial matrix | `ARCHITECTURE_PASS` |
| F18-IVM-C07 | RoleAssignment v1 has no parallel condition language | NIST SP 800-162; ISO 25010 | RoleAssignment + FEATURE-0017 boundary | Scope and explicit validity are the only assignment-local restrictions; contextual restrictions are mandatory FEATURE-0017 guardrails; condition, wildcard, deny-assignment, arbitrary-expression, and generic-ABAC fields are rejected | Schema-absence, unknown-field, and guardrail-non-granting proof | `ARCHITECTURE_PASS` |
| F18-IVM-C08 | Deny-by-default | NIST SP 800-207; NIST SP 800-53 AC | Authorization composition | Missing, invalid, expired, ambiguous, conflicting, cross-tree, or resource-mismatched evidence denies | Negative matrix for every missing/mismatched input | `ARCHITECTURE_PASS` |
| F18-IVM-C09 | Multiple grants compose deterministically | ISO 29146; NIST SP 800-162; AWS/Azure/Google IAM comparison | Authorization composition | Independently applicable grants union; mandatory guardrails intersect; policy Deny or Indeterminate overrides; no guardrail, Allow, approval, or exception creates permission | Multi-grant, boundary-intersection, deny-override, and non-granting-evidence cases | `ARCHITECTURE_PASS` |
| F18-IVM-C10 | Policy Deny cannot be bypassed by a role | NIST SP 800-207; NIST SP 800-162 | RoleAssignment + FEATURE-0017 evidence | Deny overrides role grant; Allow continues remaining checks; F18 does not emulate missing target-aware FEATURE-0017 policy | Allow/Deny/Indeterminate/RequiresApproval and no-custom-engine matrix | `ARCHITECTURE_PASS` |
| F18-IVM-C11 | Authorization is target-bound | NIST SP 800-207; NIST SP 800-53 AC-3 | Authorization composition | Input binds actor and action to exactly one registered `ExactResource(targetKind)`, `CreateParent(parentScopeKind)`, or `ScopeOnly(scopeKind)` variant plus current time, assurance, request and correlation; result cannot authorize another target or request | Cross-target, target-UID, create-parent, scope-only, and result-replay isolation cases | `ARCHITECTURE_PASS` |
| F18-IVM-C12 | ExecutionTarget actions preserve FEATURE-0016 meaning | ISO 42010; prior-feature ownership | RoleDefinition + ExecutionTarget | `read`, `qualify`, and `write` are composed unchanged; status writer remains FEATURE-0016-owned | F16 compatibility review | `DEPENDENCY` |
| F18-IVM-C13 | ExecutionTarget action granularity risk is explicit | NIST SP 800-53 AC-6 | `executiontarget.write` | Create/retire bundling is accepted as a recorded dependency limitation; FEATURE-0018 does not split or reinterpret the FEATURE-0016 action | Architecture-owner disposition and F16 compatibility proof | `ARCHITECTURE_PASS` |
| F18-IVM-C14 | External IaaS authority is not granted | NIST SP 800-207A; CSA CCM IAM | ExecutionTarget authorization boundary | Sovrunn grants control-plane action only; a native effect requires independent least-privileged adapter identity authorization, and the two layers intersect so neither Allow substitutes for the other or overrides a Deny | Zero-native-credential, two-plane allow/deny truth table, and no-external-call proof | `ARCHITECTURE_PASS` |
| F18-IVM-C15 | Assignment validity mode is explicit | NIST SP 800-53 AC-2/AC-6; NIST SP 800-207 | RoleAssignment | `Standing` has no automatic expiry but is revocable and periodically reviewed; `TimeBound` requires `notBefore` and `expiresAt` and denies at expiry without grace; omission cannot change modes implicitly | Mode-schema, missing-timestamp, exact-boundary, revocation, and retained-evidence tests | `ARCHITECTURE_PASS` |
| F18-IVM-C16 | Published role versions have an emergency kill switch | NIST SP 800-53 AC-2/AC-6/IR; CIS Control 6 | RoleDefinition + RoleAssignment | Suspending an exact published role version immediately makes every referencing assignment ineffective without mutating definition or assignment history; restoration is authorized and audited | Suspend/use race, cross-version isolation, restore, and audit-failure cases | `ARCHITECTURE_PASS` |
| F18-IVM-C17 | Authorization input, result, and provenance are exact | NIST SP 800-162/207; ISO 29146 | Authorization value boundary | Input binds actor/action and one ExactResource, CreateParent, or ScopeOnly variant plus time/assurance/request; result binds outcome, contributors, exact versions, scope relationship, evidence, time, and stable reasons; result is not reusable authority | Three-variant mismatch, replay denial, provenance completeness, projection, and no-bearer-authority cases | `ARCHITECTURE_PASS` |
| F18-IVM-C18 | Superseded role version is not an authorization state | ISO 25010; NIST SP 800-53 CM/AC | RoleDefinition | Effective lifecycle is Draft, Published, Suspended, and Retired; supersession only reports a newer version and never changes pinned assignments or authorization | Successor-publication, cross-version, suspend, retire, and projection cases | `ARCHITECTURE_PASS` |

## 8. Gate D — Approval and privileged-access completeness

| ID | Validation concern | Reference | Resource or boundary | Pass condition | Required evidence | Current status |
| F18-IVM-D01 | ApprovalPolicy has a bounded purpose | ISO 42010; NIST CSF Govern | ApprovalPolicy | Policy defines applicability, eligible approvers, one-to-three linear stages/quorum, expiry, justification, and separation of duties; branching, loops, delegation chains, scripted escalation, and arbitrary workflow expressions are absent | Contract matrix, field-necessity review, stage-boundary tests, and curated template fixtures | `ARCHITECTURE_PASS` |
| F18-IVM-D02 | ApprovalPolicy is versioned and immutable after publication | ISO 29146; NIST SP 800-53 CM/AC | ApprovalPolicy | Request pins an exact published policy version; changes do not alter in-flight evidence | Version-change/in-flight-request tests | `ARCHITECTURE_PASS` |
| F18-IVM-D03 | ApprovalRequest has a bounded subject | ISO 42010; NIST SP 800-53 AC | ApprovalRequest | Subject is exactly PAR ref, immutable RoleAssignmentProposal, or immutable ExceptionProposal; arbitrary payloads are rejected | Subject-kind and unknown-payload tests | `ARCHITECTURE_PASS` |
| F18-IVM-D04 | Approval terminal outcomes are immutable | NIST SP 800-53 AU/AC | ApprovalRequest | Approved, Denied, Cancelled, or Expired is terminal and retained | State-machine and invalid-transition proof | `ARCHITECTURE_PASS` |
| F18-IVM-D05 | Approver eligibility is current | NIST SP 800-207; NIST SP 800-53 AC | ApprovalRequest + Membership/RoleAssignment | Eligibility is re-evaluated at decision time | Approver-revoked-during-request case | `ARCHITECTURE_PASS` |
| F18-IVM-D06 | Self-approval and conflicts of interest are denied | NIST SP 800-53 AC-5 | ApprovalPolicy + ApprovalRequest | Conflict set includes requester/direct beneficiary, AccessGroup owner/current direct members, Workload/System responsible party, and beneficiaries resolved through assignment/request refs; it is re-evaluated at decision and effect publication | Direct, group-derived and responsible-party conflict cases at both boundaries | `ARCHITECTURE_PASS` |
| F18-IVM-D07 | Concurrent terminal actions resolve exactly once | NIST SP 800-53 AU; ISO 25010 reliability | ApprovalRequest | First valid terminal CAS wins; same-key/same-digest replays; same-key/different-content conflicts; audit precedes publication | Race, idempotency, and audit-failure proof | `ARCHITECTURE_PASS` |
| F18-IVM-D08 | Approval does not itself grant access | ISO 29146; NIST SP 800-207 | ApprovalRequest | Approval is evidence; only a freshly authorized RoleAssignment grants action | Approved-without-assignment and changed-state denial proof | `ARCHITECTURE_PASS` |
| F18-IVM-D09 | Privileged access is strongly authenticated | NIST SP 800-63B-4; NIST SP 800-207 | PrivilegedAccessRequest | Normal JIT and break-glass activation require provider-neutral phishing-resistant AssuranceEvidence at AAL2 or higher and no more than 15 minutes old; the applicable privileged-access rule may strengthen but not weaken the floor; raw factors are absent | JIT and break-glass missing/stale/non-phishing-resistant/insufficient assurance cases | `ARCHITECTURE_PASS` |
| F18-IVM-D10 | JIT request is immutable after acceptance | NIST SP 800-53 AC-6; ISO 29146 | PrivilegedAccessRequest | Principal, role version, target, duration, and justification cannot change in flight | Mutation and replacement tests | `ARCHITECTURE_PASS` |
| F18-IVM-D11 | JIT activation creates at most one temporary assignment | NIST SP 800-53 AC-6; ISO 25010 reliability | PAR + ApprovalRequest + RoleAssignment | Atomic publication occurs strictly before immutable activationDeadline and produces exactly one individual bounded assignment; once active, its RoleAssignment expiry—not activationDeadline—controls access | Before/at/after deadline, replay, changed-state, and concurrent activation proof | `ARCHITECTURE_PASS` |
| F18-IVM-D12 | JIT assignment is time-bound and automatically expires | NIST SP 800-207; NIST SP 800-53 AC-2/AC-6 | PAR + temporary RoleAssignment | JIT always creates a `TimeBound` assignment with mandatory `notBefore` and `expiresAt`; access is ineffective at exact expiry without grace | Fixed-clock boundary and expiry-race tests | `ARCHITECTURE_PASS` |
| F18-IVM-D13 | Renewal requires fresh authorization and approval | NIST SP 800-207; least privilege | PAR | Active duration is not extended by mutation; renewal is a new request | Renewal/replay negative tests | `ARCHITECTURE_PASS` |
| F18-IVM-D14 | Break-glass remains controlled | NIST SP 800-53 AC/IA/AU; NIST SP 800-63-4 | PAR + privilegedAccessRule + AccessReview | `privilegedAccessRule.breakGlassAllowed` is the sole bypass authority; ApprovalPolicy never bypasses approval. Eligibility is individual and pre-authorized, phishing-resistant assurance is mandatory, duration is at most 1 hour, audit is immediate, mutation renewal is forbidden, and review is due within 24 hours | Bypass-owner, assurance, duration, review and audit-failure cases | `ARCHITECTURE_PASS` |
| F18-IVM-D15 | Privileged-duration ceiling is bounded | NIST SP 800-53 AC-2/AC-6; NIST SP 800-207 | ApprovalPolicy + PAR | Normal JIT defaults to 1 hour and has an 8-hour platform maximum; policy may only narrow it | Boundary-value, default, and over-maximum denial tests | `ARCHITECTURE_PASS` |
| F18-IVM-D16 | Request resources compose into one business journey | ISO 25010; NIST SP 800-63-4 customer experience | PAR + ApprovalRequest + temporary RoleAssignment | Requester sees one correlated request/status/expiry journey; approvers see one actionable task; protected contracts remain separate; no facade resource or second requester-managed request is introduced | Requester/approver projection ledger and correlated journey cases | `ARCHITECTURE_PASS` |

## 9. Gate E — Access review and exception completeness

| ID | Validation concern | Reference | Resource or boundary | Pass condition | Required evidence | Current status |
|---|---|---|---|---|---|---|
| F18-IVM-E01 | Review snapshot is immutable | NIST SP 800-53 AC/CA; ISO 29146 | AccessReview | Review captures exact membership/assignment identities, versions, scopes, validity, and reviewers at start | Snapshot mutation and stale-input tests | `ARCHITECTURE_PASS` |
| F18-IVM-E02 | Reviewer independence is enforced | NIST SP 800-53 AC-5; CSA CCM IAM | AccessReview | No principal may preserve, replace or advance authority from which they benefit; for an AccessGroup holder the Human conflict set includes its owner, direct Human members, and the responsible Human of every direct Workload/System member; current and proposed sets are rechecked at atomic decision/effect publication and unresolved expansion fails closed for Retain/Replace without blocking Revoke | Direct-holder, group owner/member, group-to-workload-to-responsible-Human, replacement-union, unresolved-expansion, changed-conflict and Revoke cases | `ARCHITECTURE_PASS` |
| F18-IVM-E03 | Review state and item disposition are distinct | NIST SP 800-53 AC/CA | AccessReview | Campaign lifecycle and item disposition are separate closed vocabularies; StandingCertification, Manual and BreakGlassRetrospective use their exact mode-specific timing, and only exact-version StandingCertification `Retain` advances the due date | Mode/state/disposition matrix and due-time vectors | `ARCHITECTURE_PASS` |
| F18-IVM-E04 | Revocation remediation is explicit and exactly once | NIST SP 800-53 AC/AU | AccessReview + RoleAssignment | Revoke/Replace is exact-assignment-bound, audit-before-publication, idempotent, and retained | Duplicate remediation and audit-failure tests | `ARCHITECTURE_PASS` |
| F18-IVM-E05 | Stale snapshot cannot modify a changed assignment silently | ISO 25010 reliability; NIST SP 800-53 CA | AccessReview | Changed assignment yields Stale and no remediation | Concurrent review/replacement proof | `ARCHITECTURE_PASS` |
| F18-IVM-E06 | Incomplete review never counts as certification | NIST SP 800-53 CA-7; CIS account/access controls | AccessReview | Human privileged/guest access expires independently; standing production workload/system access becomes overdue/escalated and is not revoked solely for missing telemetry or incomplete review absent prepublished suspension policy | Incomplete/expiry/workload-safety cases | `ARCHITECTURE_PASS` |
| F18-IVM-E07 | Exception proposal has a bounded carrier | ISO 42010; NIST CSF Govern | ApprovalRequest boundary | No ExceptionRequest or arbitrary payload; immutable proposal pins control, subject, scope, interval, justification, and compensating controls | Proposal grammar and unknown-field tests | `ARCHITECTURE_PASS` |
| F18-IVM-E08 | Only explicitly exceptionable controls can be overridden | NIST SP 800-37; NIST SP 800-53 RA/CA | ExceptionGrant | Control-owning versioned definition declares exceptionability and maximum duration; F18 consumes without redefining semantics | Exceptionability registry/owner and denial tests | `ARCHITECTURE_PASS` |
| F18-IVM-E09 | Security invariants remain non-exceptionable | NIST SP 800-207; NIST SP 800-53 AC/AU | ExceptionGrant | Authentication integrity, scope/tenant isolation, audit integrity, immutable evidence, and unresolved conflict cannot be overridden | Attempted-invariant-exception cases | `ARCHITECTURE_PASS` |
| F18-IVM-E10 | Exception is narrow and time-bound | NIST SP 800-37; CSA CCM GRC | ExceptionGrant | Exact subject/control/scope/interval/justification/compensating controls/evidence are required; platform maximum is 90 days and policy may narrow | Missing/broad/expired/over-maximum cases | `ARCHITECTURE_PASS` |
| F18-IVM-E11 | Exception does not grant a role | ISO 29146; NIST SP 800-207 | ExceptionGrant + RoleAssignment | Applicable role grant remains independently required | Exception-without-role denial proof | `ARCHITECTURE_PASS` |
| F18-IVM-E12 | Immutable grant supports early revocation coherently | NIST SP 800-53 AU; ISO 25010 | ExceptionGrant | Linked immutable same-kind effect=Revoke evidence names exact revokedGrantRef without mutating history or adding a kind | Revocation schema and early-revocation proof | `ARCHITECTURE_PASS` |
| F18-IVM-E13 | Overlap and conflict resolve deterministically | NIST SP 800-162; NIST SP 800-37 | ExceptionGrant | F18 rejects active overlap for same control/subject/scope/interval and infers no precedence; F20 owns cross-profile effective resolution | Pairwise overlap and ownership-boundary matrix | `ARCHITECTURE_PASS` |
| F18-IVM-E14 | FEATURE-0020 application boundary is preserved | ISO 42010; Sovrunn ownership | ExceptionGrant | FEATURE-0018 validates/issues evidence only; effective governance application remains FEATURE-0020-owned | Dependency-direction and no-resolver proof | `INHERITED` |
| F18-IVM-E15 | AccessReview usage evidence is immutable and qualified | NIST SP 800-53 AC-2/CA-7; CIS Controls 5/6; AWS IAM Access Analyzer comparison | AccessReview + FEATURE-0013 evidence | Review may pin versioned usage evidence with window, last observed use, action set, generation time, and coverage quality; absence/incompleteness is not silently treated as non-use; RoleAssignment has no mutable usage field | Complete/stale/absent/incomplete evidence and production-workload safety cases | `ARCHITECTURE_PASS` |

## 10. Gate F — GovernanceProfile completeness and leakage prevention

| ID | Validation concern | Reference | Resource or boundary | Pass condition | Required evidence | Current status |
|---|---|---|---|---|---|---|
| F18-IVM-F01 | GovernanceProfile has a bounded FEATURE-0018 purpose | ISO 42010; NIST CSF Govern | GovernanceProfile | v1 typed components are approvalPolicyRefs, privilegedAccessRules, accessReviewRules, exceptionRules, and auditRequirements; F18 owns envelope/lifecycle, not future-domain semantics | Exact v1 field/component ledger | `ARCHITECTURE_PASS` |
| F18-IVM-F02 | Published profiles are immutable | ISO 25010; NIST SP 800-53 CM | GovernanceProfile | Draft mutable and Published immutable; supersession prevents new selection but preserves valid exact pins; retirement is terminal and blocks new selection and operations requiring current rule evidence; neither silently migrates retained artifacts | Publication, supersession-pin, retirement-current-evidence and no-migration tests | `ARCHITECTURE_PASS` |
| F18-IVM-F03 | Lower scope cannot weaken mandatory controls | NIST SP 800-53 AC/PM; CSA CCM GRC | GovernanceProfile boundary | Narrowing is permitted; weakening higher-scope mandatory controls is prohibited | Profile compatibility/weakening cases | `INHERITED` |
| F18-IVM-F04 | Sovereignty remains separate | NIST CSF Govern; Sovrunn canonical model | GovernanceProfile | No sovereignty facts, evidence, or assessment semantics enter FEATURE-0018 | Schema and terminology sentinel checks | `INHERITED` |
| F18-IVM-F05 | Future domains do not leak into v1 | ISO 42010; ISO 25010 | GovernanceProfile | No active placement, backup-engine, cost-engine, entitlement, quota, sovereignty, or execution semantics | Future-field rejection tests | `ARCHITECTURE_PASS` |
| F18-IVM-F06 | Profile assignment and resolution remain external | ISO 42010 | GovernanceProfile | No ProfileAssignment writer, hierarchy traversal, conflict resolver, or EffectiveGovernanceContext construction | Dependency-direction/no-resolver tests | `INHERITED` |
| F18-IVM-F07 | ApprovalPolicy relationship is exact and version-pinned | ISO 29146; NIST SP 800-53 AC | GovernanceProfile + ApprovalPolicy | Zero or more exact published refs are allowed; applicability keys are unique and ambiguous overlap is rejected | Contract matrix, duplicate-key, and invalid-version cases | `ARCHITECTURE_PASS` |
| F18-IVM-F08 | GovernanceProfile is not a routine IAM user surface | ISO 25010; NIST SP 800-63-4 customer experience | GovernanceProfile projections | Only governance administrators manage curated typed profiles; routine authorization, assignment, approval, and review journeys neither create nor select them | Audience/projection matrix and routine-journey absence proof | `ARCHITECTURE_PASS` |

## 11. Gate G — Decision, audit, policy-evaluation, and failure accuracy

| ID | Validation concern | Reference | Resource or boundary | Pass condition | Required evidence | Current status |
|---|---|---|---|---|---|---|
| F18-IVM-G01 | FEATURE-0013 remains sole evidence-envelope owner | ISO 42010; NIST SP 800-53 AU | DecisionRecord/AuditEvent adoption | No competing F18 decision or audit envelope exists | Type/schema/import inventory check | `DEPENDENCY` |
| F18-IVM-G02 | Decision profiles are registered and bounded | NIST SP 800-53 AU/CA | authorization/approval/exception profiles | Exact profiles are authorization-decision/v1, approval-decision/v1, and exception-decision/v1; AccessReview remains its own retained LRO | FEATURE-0013 registration, field, and projection evidence | `ARCHITECTURE_PASS` |
| F18-IVM-G03 | Audit event taxonomy is complete but not excessive | NIST SP 800-53 AU-2/AU-3; CIS Control 8 | All lifecycle resources | One bounded category taxonomy covers authorization, membership/role, approval, privileged, review, and exception lifecycle; structural noise and protected data are excluded | Event matrix and projection review | `ARCHITECTURE_PASS` |
| F18-IVM-G04 | Audit evidence is protected | NIST SP 800-53 AU-9; CSA CCM LOG | AuditEvent adoption | Unauthorized alteration/deletion is prevented by inherited contract; confidential fields are redacted | Dependency proof and redaction cases | `DEPENDENCY` |
| F18-IVM-G05 | Required audit failure prevents sensitive publication | NIST SP 800-53 AU; ISO 25010 reliability | State-changing operations | No role grant, JIT activation, terminal approval, remediation, revocation, or exception publishes or completes replay without required evidence | Audit-failure atomicity cases | `ARCHITECTURE_PASS` |
| F18-IVM-G06 | FEATURE-0017 outcomes are interpreted exactly | NIST SP 800-162; dependency authority | FEATURE-0017 adoption | Allow continues checks; Deny denies; RequiresApproval evaluates policy; Indeterminate fails closed | Four-outcome composition tests | `ARCHITECTURE_PASS` |
| F18-IVM-G07 | FEATURE-0017 SUCCESS is not authorization | NIST SP 800-207 | FEATURE-0017 adoption | Result transport success cannot bypass role, resource, scope, validity, approval, or exception checks | SUCCESS-with-invalid-grant denial tests | `ARCHITECTURE_PASS` |
| F18-IVM-G08 | FEATURE-0017 fake remains policy-free | ISO 42010; reuse boundary | FEATURE-0017 fake | Digest lookup only; no IAM, approval, exception, target, or governance interpretation | Fixture and production-import review | `DEPENDENCY` |
| F18-IVM-G09 | Principal/target gap is not silently overloaded | ISO 42010; NIST SP 800-207 | FEATURE-0017 request boundary | F18 uses only compatible UID/generation-pinned subjects; target-aware dynamic policy requires separate F17 change and is not emulated | Contract-fit, no-overload, and no-custom-engine tests | `ARCHITECTURE_PASS` |
| F18-IVM-G10 | Safe denial prevents existence disclosure | NIST SP 800-207; NIST SP 800-53 AC | Authorization/reference resolution | Inaccessible target/reference receives inherited safe outcome before detailed validation | Cross-scope safe-denial matrix | `DEPENDENCY` |

## 12. Gate H — Determinism, conformance, and operational assurance

| ID | Validation concern | Reference | Resource or boundary | Pass condition | Required evidence | Current status |
|---|---|---|---|---|---|---|
| F18-IVM-H01 | Validation order is closed | ISO 25010; NIST SP 800-53A; OWASP ASVS | All request boundaries | Bounded transport → authentication → structural validation → local semantic/cross-field validation → safe resolution → assurance/contextual Membership → assignment/action/scope/resource/delegation and membership-enabled-grant envelope → required EligibilityRef and separation of duties → approval/exception evidence → concurrency/idempotency → decision/audit-obligation acceptance → publication | Ordered negative-case matrix | `ARCHITECTURE_PASS` |
| F18-IVM-H02 | Time behavior is deterministic | ISO 25010 reliability; NIST SP 800-53 AC | Validity, approval, JIT, review, exception | Injected UTC clock; exact-boundary expiry; approved JIT/break-glass/exception ceilings; no wall-clock fixture dependence | Fixed-clock vectors and expiry races | `ARCHITECTURE_PASS` |
| F18-IVM-H03 | Fixtures are immutable and synthetic | NIST SP 800-53A; OWASP ASVS | In-memory fixtures | Constructor validation, defensive copies, no customer data, no shared mutable state | Fixture mutation/leakage tests | `ARCHITECTURE_PASS` |
| F18-IVM-H04 | Positive and negative proof are paired | NIST SP 800-53A | Every invariant | Each material allow has applicable missing, malformed, unauthorized, expired, revoked, stale, cross-scope, replay, audit-failure, and redaction cases | Conformance traceability matrix | `ARCHITECTURE_PASS` |
| F18-IVM-H05 | Race behavior is proven | ISO 25010 reliability; NIST SP 800-53A | Approval/JIT/review/exception | Terminal approval, JIT activation, revoke/use, review/remediation, and exception revocation are deterministic and race-clean | Race tests and operation-state matrices | `ARCHITECTURE_PASS` |
| F18-IVM-H06 | Idempotency is target and principal bound | NIST SP 800-53 AU/AC; prior-feature precedent | State-changing operations | Key binds actor, operation, exact target, scope, and digest; replay rechecks current auth and safe visibility; changed content conflicts | Replay-isolation, changed-access, and changed-content tests | `ARCHITECTURE_PASS` |
| F18-IVM-H07 | No external effects exist in Phase 2R proof | NIST SP 800-53A; CSA CCM | Fakes and tests | No network, IdP, policy/workflow engine, CloudProvider SDK, database, Kubernetes, provisioning, or execution | Import, declaration, network, and side-effect sentinels | `ARCHITECTURE_PASS` |
| F18-IVM-H08 | API verification profile is pinned without architecture drift | OWASP ASVS 5.0 | Architecture/conformance boundary | ASVS 5.0 Level 2 applies platform-wide; Level 3 applies privileged, approval, exception and security-administration surfaces. Requirements may add proof, not alter canonical semantics | FEATURE-0018 standards mapping and independent review | `ARCHITECTURE_PASS` |
| F18-IVM-H09 | Operational account/access/log controls are cross-checked | CIS Controls 5, 6, 8 | Membership, assignments, reviews, audit | CIS mapping is required before final architecture approval without creating new canonical authority | CIS mapping review | `ARCHITECTURE_PASS` |
| F18-IVM-H10 | Cloud shared responsibility is explicit | CSA CCM 4.1 | CloudPlatform/CloudProvider/Organization boundaries | Identity, membership, role, approval, audit, and external-provider responsibilities are mapped before final architecture approval | CSA CCM responsibility matrix | `ARCHITECTURE_PASS` |
| F18-IVM-H11 | Active decision-status mirrors are consistent | ISO 42010; Sovrunn architecture change control | Decision Index + traceability + Phase 2R rebaseline | A lower-precedence active matrix cannot mark a superseded decision Accepted; superseding DEC identity must match the Decision Index and rebaseline while historical records remain immutable | Accepted/proposed/superseded status and replacement-ID drift cases | `ARCHITECTURE_PASS` |

## 13. Gate I — Future adapter compatibility without scope leakage

| ID | Validation concern | Reference | Resource or boundary | Pass condition | Required evidence | Current status |
|---|---|---|---|---|---|---|
| F18-IVM-I01 | OIDC identity maps without schema change | OIDC Core; NIST SP 800-63C-4 | PrincipalRef boundary | Paper mapping proves issuer/subject and assurance compatibility without raw tokens, IdP-native roles, adapter choice, or implementation | Paper mapping and incompatible-claim cases | `ARCHITECTURE_PASS` |
| F18-IVM-I02 | OAuth integration follows current security BCP | RFC 9700 | Future authorization adapter | No insecure OAuth assumptions are embedded in FEATURE-0018; adapter can enforce current BCP | Future adapter assessment | `EXCLUDED` |
| F18-IVM-I03 | SCIM users/groups do not become canonical authority accidentally | RFC 7643/7644 | AccessGroup/Membership boundary | Paper mapping preserves trusted provisioning into canonical resources; raw external claims, nested/dynamic membership, and unprovisioned groups grant nothing | Paper mapping, provisioning provenance, and raw-claim denial | `ARCHITECTURE_PASS` |
| F18-IVM-I04 | WebAuthn assurance can support privileged access | WebAuthn; NIST SP 800-63B-4 | PAR assurance evidence | Paper mapping proves provider-neutral strength/freshness/phishing-resistance representation without WebAuthn types or adapter design in core | Paper mapping | `ARCHITECTURE_PASS` |
| F18-IVM-I05 | Workload identity can map through SPIFFE-compatible adapter | SPIFFE; NIST SP 800-207A | Workload PrincipalRef | Paper mapping proves issuer/subject and assurance compatibility without SPIFFE types or adapter implementation in canonical schema | Paper mapping | `ARCHITECTURE_PASS` |
| F18-IVM-I06 | Real adapters require fresh FEATURE-0011 assessment | Sovrunn reuse standard; CSA CCM | All future adapters | Selection evaluates current security, licensing, sovereignty, deployment, lifecycle, and migration characteristics | Approved future reuse assessment | `EXCLUDED` |

## 14. Enterprise scenario validation matrix

| ID | Enterprise scenario | Primary resources | Expected architecture outcome | Evidence required | Current status |
|---|---|---|---|---|---|
| F18-SCN-01 | Active employee receives project-scoped developer access | PrincipalRef, Membership, RoleDefinition, RoleAssignment | Exact-scope grant succeeds; unrelated scope denies | Positive and cross-scope cases | `ARCHITECTURE_PASS` |
| F18-SCN-02 | Employee email changes | PrincipalRef | Identity remains stable through issuer/subject | Identity-stability case | `INHERITED` |
| F18-SCN-03 | Employee is suspended during an active assignment | Membership, RoleAssignment | Access becomes ineffective immediately; history remains | Suspension/use race | `ARCHITECTURE_PASS` |
| F18-SCN-04 | Employee leaves and later rejoins | Membership, RoleAssignment | Revoked/expired assignments do not silently reactivate | Rejoin sequence | `ARCHITECTURE_PASS` |
| F18-SCN-05 | Workload calls a control-plane operation | PrincipalRef, RoleAssignment | Workload identity is independently scoped and audited | Workload allow/deny cases | `ARCHITECTURE_PASS` |
| F18-SCN-06 | System controller performs an owned transition | PrincipalRef, writer contract | System identity does not override sole-writer rules | Correct/wrong writer cases | `ARCHITECTURE_PASS` |
| F18-SCN-07 | Consultant belongs to supplier and customer contexts | Membership, RoleAssignment | Sponsorship, target scope, expiry, and attribution are unambiguous | Guest/multiple-membership suite | `ARCHITECTURE_PASS` |
| F18-SCN-08 | External IdP group claim and provisioned AccessGroup are distinguished | PrincipalRef, AccessGroup, Membership, RoleAssignment | Raw claim alone denies; directly provisioned canonical membership plus an active scoped group assignment may allow | Raw-claim denial and canonical-group allow/deny pair | `ARCHITECTURE_PASS` |
| F18-SCN-09 | Organization administrator delegates project authority | RoleDefinition, RoleAssignment | Grantor cannot assign actions, scope, target mode, exact resource reach, duration, or delegation beyond one independently applicable per-action witness | Delegation/overreach and cross-assignment-synthesis cases | `ARCHITECTURE_PASS` |
| F18-SCN-10 | ExecutionTarget operator reads and qualifies one target | RoleDefinition, RoleAssignment | Existing FEATURE-0016 actions plus optional exact resourceRef authorize only target-bound Sovrunn operations | F16 compatibility and resource-binding cases | `ARCHITECTURE_PASS` |
| F18-SCN-11 | Target operator tries to retire with qualifier role | RoleDefinition, RoleAssignment | Denied because `executiontarget.write` is absent | Action-negative case | `ARCHITECTURE_PASS` |
| F18-SCN-12 | Organization wants create-only but not retire authority | RoleDefinition boundary | Existing F16 granularity limitation is accepted and surfaced; no invented action | Recorded residual-risk decision | `ARCHITECTURE_PASS` |
| F18-SCN-13 | Engineer requests two-hour production DB administration | PAR, ApprovalPolicy, ApprovalRequest, RoleAssignment | Approved request creates one temporary assignment and expires exactly | End-to-end deterministic flow | `ARCHITECTURE_PASS` |
| F18-SCN-14 | Engineer approves own privileged request | ApprovalPolicy, ApprovalRequest | Denied and audited without activation | Self-approval case | `ARCHITECTURE_PASS` |
| F18-SCN-15 | Required approver is revoked before deciding | Membership, RoleAssignment, ApprovalRequest | Current eligibility check denies the decision | Eligibility-change case | `ARCHITECTURE_PASS` |
| F18-SCN-16 | Approval and expiry occur concurrently | ApprovalRequest | First valid terminal CAS wins; no unauthorized assignment | Race proof | `ARCHITECTURE_PASS` |
| F18-SCN-17 | Same approved JIT request is activated twice | PAR, RoleAssignment | One temporary assignment; safe replay or deterministic conflict | Exactly-once activation case | `ARCHITECTURE_PASS` |
| F18-SCN-18 | Emergency access is needed without normal approver availability | PAR, ApprovalPolicy, AccessReview | Individually pre-authorized break-glass lasts at most one hour and requires immediate audit and review within 24 hours | Break-glass suite | `ARCHITECTURE_PASS` |
| F18-SCN-19 | Access review finds stale administrator access | AccessReview, RoleAssignment | Explicit exactly-once revocation with retained evidence | Review/remediation flow | `ARCHITECTURE_PASS` |
| F18-SCN-20 | Assignment changes during review | AccessReview, RoleAssignment | Explicit Stale outcome performs no remediation | Review/replace race | `ARCHITECTURE_PASS` |
| F18-SCN-21 | Incident requires change-window exception plus DB role | ApprovalRequest, ExceptionGrant, RoleAssignment | Both bounded exception and applicable role are required | Missing-role/missing-exception pair | `ARCHITECTURE_PASS` |
| F18-SCN-22 | Request attempts to except tenant isolation | ExceptionGrant | Denied as non-exceptionable | Non-exceptionable-control case | `ARCHITECTURE_PASS` |
| F18-SCN-23 | Exception is revoked before its planned expiry | ExceptionGrant | Immutable effect=Revoke record makes the exact grant ineffective | Linked-revocation flow | `ARCHITECTURE_PASS` |
| F18-SCN-24 | Two exception grants overlap | ExceptionGrant | F18 rejects active overlap for the same control/subject/scope/interval; F20 owns cross-profile resolution | Pairwise overlap and boundary cases | `ARCHITECTURE_PASS` |
| F18-SCN-25 | FEATURE-0017 returns Allow but role is expired | RoleAssignment, FEATURE-0017 result | Final authorization denies | Policy-success-not-authority case | `ARCHITECTURE_PASS` |
| F18-SCN-26 | FEATURE-0017 returns RequiresApproval | FEATURE-0017 result, ApprovalPolicy | F18 evaluates only a compatible pinned subject and creates a request only through a valid explicit flow | RequiresApproval composition flow | `ARCHITECTURE_PASS` |
| F18-SCN-27 | FEATURE-0017 returns Indeterminate | FEATURE-0017 result | Fail closed with safe decision/audit behavior | Indeterminate case | `ARCHITECTURE_PASS` |
| F18-SCN-28 | Required AuditEvent append fails during JIT activation | PAR, RoleAssignment, AuditEvent | No temporary grant or completed replay is published | Audit atomicity case | `ARCHITECTURE_PASS` |
| F18-SCN-29 | End customer performs a routine low-risk operation | Membership, RoleAssignment, authorization composition | Identity, grants, guardrails, provenance, and audit resolve silently; user supplies only business input and receives result or safe denial | Routine projection and no-IAM-input case | `ARCHITECTURE_PASS` |
| F18-SCN-30 | User receives denial involving inaccessible target | Authorization boundary | Safe result reveals no target existence or sensitive policy detail | Non-disclosure case | `ARCHITECTURE_PASS` |
| F18-SCN-31 | Production workload needs durable narrow access | RoleAssignment, AccessReview | Explicitly policy-permitted `Standing` assignment avoids accidental outage, remains revocable, and enters periodic review | Standing-mode allow, revocation, and review-due cases | `ARCHITECTURE_PASS` |
| F18-SCN-32 | Guest assignment omits expiry | Membership, RoleAssignment | Validation denies because guest assignments must be `TimeBound` with both timestamps | Missing-mode and missing-timestamp negative cases | `ARCHITECTURE_PASS` |
| F18-SCN-33 | Two assignments grant an action but a mandatory guardrail excludes it | RoleAssignment, authorization composition | Grant union cannot bypass the restrictive boundary; final result denies with exact provenance | Grant-union/boundary-intersection pair | `ARCHITECTURE_PASS` |
| F18-SCN-34 | A published role version is suspended during use | RoleDefinition, RoleAssignment | Every assignment pinned to that version becomes ineffective; other versions are unaffected; authorized restoration is auditable | Suspend/use race, version isolation, and restoration cases | `ARCHITECTURE_PASS` |
| F18-SCN-35 | AuthorizationResult is replayed for a different target | AuthorizationInput, AuthorizationResult | Result is not bearer authority and cannot authorize another actor, action, target, scope, or instant | Cross-target and cross-principal replay denial | `ARCHITECTURE_PASS` |
| F18-SCN-36 | Federated membership provenance is stale or responsible owner is absent | Membership | Authorization fails safely where current provenance/ownership is required; raw external claim cannot repair it | Staleness, missing-owner, and raw-claim cases | `ARCHITECTURE_PASS` |
| F18-SCN-37 | Access review has no or incomplete usage telemetry | AccessReview | Absence is not treated as proof of non-use and does not silently revoke a production workload | Complete/absent/stale/incomplete usage-evidence cases | `ARCHITECTURE_PASS` |
| F18-SCN-38 | Retired duplicate of F18-SCN-29 | — | No independent scenario or evidence authority remains; all routine low-risk journey evidence is attributed to F18-SCN-29 | Duplicate-absence check | `EXCLUDED` |
| F18-SCN-39 | Access administrator grants ordinary project access | AccessGroup/PrincipalRef, RoleDefinition, RoleAssignment | Administrator supplies holder, named role, scope, and explicit validity; exact refs, version, provenance, and audit are system-managed | Four-input assignment journey and forged-system-field negatives | `ARCHITECTURE_PASS` |
| F18-SCN-40 | Engineer requests time-bound privileged access | PAR, ApprovalPolicy, ApprovalRequest, RoleAssignment | Engineer supplies role, scope, duration, and justification and sees one request status and expiry; internal request/evidence separation remains hidden | Single-request projection and correlation case | `ARCHITECTURE_PASS` |
| F18-SCN-41 | Approver decides a privileged or exception request | ApprovalPolicy, ApprovalRequest | Approver sees requester, access/control, scope, duration, risk context, quorum progress, and required justification; exact versions and concurrency controls remain hidden | Approver projection and system-field ownership case | `ARCHITECTURE_PASS` |
| F18-SCN-42 | Reviewer certifies or remediates access | AccessReview, RoleAssignment | Reviewer sees holder, role, scope, validity, qualified usage summary, and Retain/Revoke/Replace actions; immutable snapshot/evidence mechanics remain hidden | Review projection with complete/stale/absent evidence cases | `ARCHITECTURE_PASS` |
| F18-SCN-43 | User requests a bounded control exception | ApprovalRequest, ExceptionGrant | User supplies exact control/scope/duration/justification/compensating controls and sees one request status and effective period; grant and evidence records remain internal | Exception journey, non-exceptionable control, and no-role denial cases | `ARCHITECTURE_PASS` |
| F18-SCN-44 | Resource-narrowed administrator tries to delegate another resource or the whole scope | RoleAssignment | `ExactResource(A)` authority cannot grant resource B or scope-wide reach; no combination of partial assignments broadens it | Resource-A-to-B, resource-to-scope, scope-to-resource and per-action witness cases | `ARCHITECTURE_PASS` |
| F18-SCN-45 | Standing beneficiary attempts to certify retained access | AccessReview, RoleAssignment, AccessGroup | Direct Human holder, group owner/direct Human member, direct Workload/System responsible Human, or responsible Human behind a direct Workload/System group member receives `REVIEWER_CONFLICT`; unresolved expansion fails closed and no Retain/Replace effect or due advancement publishes, while independently authorized Revoke remains available | Direct and compounded-indirect self-certification, replacement-union, unresolved-expansion, changed-conflict and Revoke cases | `ARCHITECTURE_PASS` |
| F18-SCN-46 | Temporary membership administrator attempts to create durable group-derived access | Membership, AccessGroup, RoleAssignment | Membership controller derives the exact newly enabled envelope; a TimeBound administrator cannot enable a Standing group-held effect, missing or mismatched witnesses deny, trusted provisioning has no bypass, and a concurrent assignment/member change conflicts or retries; empty-group addition grants nothing | Empty/non-empty envelope, bounded success, JIT-to-Standing laundering, scope/resource mismatch, stale refresh, trusted-provisioner and race cases | `ARCHITECTURE_PASS` |
| F18-SCN-47 | Membership administrator attempts to manufacture approval or privileged eligibility | EligibilityRef, Membership, ApprovalPolicy, privilegedAccessRule | AccessGroupRef, Membership and group-derived authority are structurally rejected; an empty-group membership cannot create approver, JIT or break-glass qualification; an exact named Human still requires the separate canonical action | Empty-group direct-policy/rule, self/accomplice approver, JIT, break-glass, trusted-provisioner and concurrent rule/group change cases | `ARCHITECTURE_PASS` |
| F18-SCN-48 | Reviewer eligibility is missing, indirect, stale or scope-incompatible | EligibilityRef, accessReviewRule, AccessReview | Non-empty OR list accepts exact Human PrincipalRefs only; every role/assignment/group/Membership/claim kind, inactive Human, scope/target-incompatible rule/action or action-missing evidence denies; eligibility and beneficiary conflict are rechecked before remediation | Positive exact-Human and negative empty/unknown/role/assignment/group/Membership/claim/inactive/scope/target/action/SoD/publication-race cases | `ARCHITECTURE_PASS` |
| F18-SCN-49 | Role grantor or later rule publisher attempts to manufacture eligibility | EligibilityRef, RoleDefinition, RoleAssignment, ApprovalPolicy, privilegedAccessRule, accessReviewRule | RoleDefinition and RoleAssignment reject as EligibilityRef; granting an ordinary role or later referencing any role cannot change approver, reviewer, JIT or break-glass eligibility; only exact named Humans qualify and still require the separate action | Self/accomplice role grant, already broad role, later policy/rule publication, approval/review/JIT/break-glass and concurrent role/rule cases | `ARCHITECTURE_PASS` |

F18-SCN-29 and F18-SCN-39 through F18-SCN-43 are approved architecture journey
evidence, not implementation, usability-test, or conformance proof. F18-SCN-44
through F18-SCN-49 are security-closure evidence. F18-SCN-38 is retired as a
duplicate. These decisions close the intended architecture interaction boundary;
executable usability and conformance proof remains a later feature-gate duty.

## 15. Preliminary findings

### 15.1 Directionally strong areas

- The inherited nine resources and approved AccessGroup have distinguishable
  primary responsibilities under the reconciled canonical model.
- Membership and RoleAssignment are correctly separated.
- Standing access, JIT elevation, access review, approval, and exception
  evidence are modeled separately.
- Provider-native IAM types remain outside canonical contracts.
- FEATURE-0013 and FEATURE-0017 ownership can be reused without duplicate
  evidence or policy-engine contracts.
- The split between FEATURE-0018 definitions and FEATURE-0020 resolution
  prevents governance-resolution leakage.
- Resource separation is preserved without exposing every resource as a user
  journey; six progressive-disclosure journeys are approved for handoff.
- RoleAssignment v1 has no user-authored condition language, approval stages
  are linearly bounded, provenance/evidence are system-managed, and
  GovernanceProfile is advanced-only.

### 15.2 Approved decisions and remaining executable work

The architecture owner approved the behavioral decisions on 2026-08-28 and the
exact applied architecture on 2026-08-30 after renewal-05 passed. Remaining
work is executable conformance rather than permission to invent architecture:

1. translate exact F18-RD contracts into requirements after separate human authorization;
2. implement the scenario, race, idempotency, safe-denial, redaction,
   audit-failure and zero-external-effect proof inventory; and
3. pass the later FEATURE-0018 feature gate without weakening dependency,
   provider or future-feature boundaries.

## 16. Gate summary

| Gate | Subject | Current outcome | Closure condition |
|---|---|---|---|
| A | Architecture integrity and simplicity | `ARCHITECTURE_PASS` | Resource/writer ledgers, six journey boundaries and final review passed |
| B | Identity and Membership completeness | `ARCHITECTURE_PASS` | Canonical reconciliation passed; executable proof remains a later feature gate |
| C | Roles and scoped authorization | `ARCHITECTURE_PASS` | Canonical/action/dependency registrations are accepted and traced |
| D | Approval and privileged access | `ARCHITECTURE_PASS` | Exact controlling RDs and scenario inventory passed independent review |
| E | Access review and exceptions | `ARCHITECTURE_PASS` | Exact controlling RDs and scenario inventory passed independent review |
| F | GovernanceProfile boundary | `ARCHITECTURE_PASS` | v1 envelope reconciliation and future-domain sentinels passed |
| G | Decision/audit/policy adoption | `ARCHITECTURE_PASS` | FEATURE-0013/0016 registrations and bounded FEATURE-0017 adoption are approved |
| H | Determinism and conformance | `ARCHITECTURE_PASS` | Exact ordering and local conformance inventory approved; executable proof remains |
| I | Adapter compatibility | `ARCHITECTURE_PASS` | Paper mappings and shared-responsibility boundary passed review |

The aggregate FEATURE-0018 industry-validation outcome is **PASS FOR FINAL
ARCHITECTURE APPROVAL**. DEC/reuse approval and renewal-05 are complete.
Executable proof remains a later feature-gate obligation; architecture pass
does not claim implemented control effectiveness.

## 17. Hard pass/fail rules

FEATURE-0018 passes this matrix only when:

1. every resource passes necessity and single-responsibility review;
2. every FEATURE-0018-owned concern maps to a contract and evidence;
3. every standards concern is supported, inherited, explicitly excluded, or
   assigned to a named owner;
4. no high-risk FEATURE-0018 concern is deferred without architecture-owner
   approval;
5. every lifecycle covers creation, valid transition, expiry/revocation,
   retained evidence, and invalid/concurrent transition behavior;
6. every privileged flow has positive, unauthorized, stale, forged, expiry,
   revocation, race, audit-failure, and redaction proof as applicable;
7. no resource duplicates another feature's scope, writer, decision, audit, or
   lifecycle authority;
8. future protocols map through adapters without entering canonical contracts;
9. all `GAP` and `NOT_ASSESSED` rows are dispositioned;
10. a qualified security reviewer records independent assessment evidence; and
11. the architecture owner explicitly approves the final matrix and handoff.

An aggregate percentage score is prohibited because it could hide one critical
identity, isolation, authorization, audit, or exception failure.

## 18. Evidence package required before final architecture approval

The completed validation package must include:

- exact standards applicability and available control mapping, with named
  exclusions where licensed/current clause-level assessment is unavailable;
- resource necessity and responsibility-overlap analysis;
- resource profile, scope, writer, mutability, and lifecycle matrix;
- enterprise scenario disposition for F18-SCN-01 through F18-SCN-49, with
  F18-SCN-38 retired into F18-SCN-29 and F18-SCN-29/F18-SCN-39–43 preserved as
  architecture journey evidence until executable proof exists;
- threat and abuse-case review;
- FEATURE-0013 profile and AuditEvent adoption matrix;
- FEATURE-0017 contract-fit decision;
- FEATURE-0011 reuse assessment and human approval evidence;
- FEATURE-0018-local conformance registry proposal;
- future adapter paper mappings;
- independent security-review record; and
- explicit architecture-owner approval.

## 19. Human disposition and next approval boundary

On 2026-08-28 the architecture owner approved the closure package represented
by the `ARCHITECTURE_PASS` rows. ADH-2026-070, its ACR/DEC, canonical
application and evidence package were then prepared. Renewal-05 passed, and the
architecture owner approved the exact package on 2026-08-30. This is not an
implemented-control or conformity claim and does not independently authorize
requirements, design, tasks or implementation.

The controlled sequence is: prepare ADH → human approves the exact ADH → Kiro
validates/applies it → final evidence review and exact architecture approval.
The repository now contains the approved application and security record.
Requirements generation is the next separately authorized stage.
