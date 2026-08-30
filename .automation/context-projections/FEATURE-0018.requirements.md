---
doc_type: ai_context_projection
feature: FEATURE-0018
stage: requirements
authority: non-authoritative-exact-excerpts
generated: true
---

# FEATURE-0018 Requirements Context Projection

This is a non-authoritative, mechanically generated projection of exact
approved-source excerpts. The FEATURE architecture and controlling handoff
package remain authoritative. A source-hash mismatch blocks regeneration.

## Source: `docs/architecture/canonical/sovrunn-finalized-data-model.md`

Approved SHA-256: `e08ad2713bba45233b553eaa8caaca1a51e911ca0af68e7e36e8784f3b240ef9`

### Canonical identity, authorization and governance model

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
  │     └── Project: payments-production
  │           └── ServiceInstance: payments-postgresql
  └── Tenant: Non-Production
        └── Project: developer-sandbox
```

## 8. Identity and access model

| Concept | Canonical purpose | Required behavior |
|---|---|---|
| `PrincipalRef` | Embedded stable reference to an already-authenticated Human, Workload or System identity | Identity lifecycle remains external; issuer/subject, not email, is durable identity |
| `AccessGroup` | Non-authenticating, access-only authorization subject | Explicit owner; direct auditable principal membership only; external, nested or dynamic groups grant nothing |
| `Membership` | Records Organization/CloudProvider belonging or one direct same-scope AccessGroup relationship | Grants no permission or approval/review/privileged eligibility; provenance, freshness, guest expiry and responsibility are explicit; every assignment-effect expansion that activates group-held authority passes the exact non-synthesized membership-enabled grant envelope |
| `RoleDefinition` | Immutable published version of canonical Sovrunn actions | Backend IAM actions are prohibited; privilege classification is monotonic and emergency suspension is explicit |
| `RoleAssignment` | Grants one role version to one `RoleHolderRef` at one registered scope and optional exact resource | `Standing` or `TimeBound`; per-action delegated reach containment; reviewable, revocable and written only by its controller |
| `PrivilegedAccessRequest` | Requests justified JIT or break-glass elevation | Individual; requester eligibility is a non-empty exact-Human PrincipalRef list and cannot derive from roles, assignments, groups or Membership; activation requires separate submit authority and fresh phishing-resistant AAL2-or-higher assurance; bounded, no-grace expiry and fully audited |
| `AccessReview` | Reviews an immutable access snapshot and retained usage/remediation evidence | Exact review modes and timing; reviewer eligibility is a non-empty exact-Human PrincipalRef list with separate decide authority; roles, assignments, groups and Membership never qualify; beneficiary conflict set includes accountable Humans behind direct Workload/System AccessGroup members and prevents direct/indirect self-certification; only exact-version standing certification advances the due date |
| `ApprovalPolicy` / `ApprovalRequest` | Versioned admission rules and retained approval evidence | Approver eligibility is a non-empty exact-Human PrincipalRef list with separate decide authority; roles, assignments, groups and Membership never qualify; approval is current at effect publication, never continuing bearer authority |
| `ExceptionGrant` | Immutable approved exception or linked revocation evidence | Exact subject/scope/control/interval; effective resolution remains FEATURE-0020-owned |

### 8.1 Persona realization and authorization

A persona describes job intent, responsibilities and a user journey. It is non-normative documentation, not an API resource, identity type or authorization input. One persona may map to several role templates, and one principal may hold different roles at different scopes.

```text
Persona
  describes job intent, responsibilities and user journey
  non-normative; not an API resource
        ↓ maps to
RoleDefinition
  defines permitted canonical Sovrunn actions
  may be supplied as a versioned default template
        ↓ applied through
RoleAssignment
  binds RoleHolderRef + pinned RoleDefinition version + governed scope
  optionally narrows to one exact resource; Standing or TimeBound
        ↓ evaluated by
AuthorizationResult
  evaluates actor, canonical action, exact target binding, scope,
  assignments, dependencies, policy guardrails, assurance and time
        ↓ recorded through
material DecisionRecord profile and/or AuditEvent
```

Default roles may be packaged for common platform, CloudProvider and customer journeys, but they are not permanent core enums. Providers and customers may create narrower roles from the canonical Sovrunn action vocabulary. Published role-template versions are immutable; changing permissions creates a new version and assignments do not silently acquire it.

`AuthorizationResult` is a transient, non-bearer `Allow | Deny` result with exact
decision provenance, not a persistent competing decision envelope. Applicable
grants combine by union; all guardrails and independent dependencies combine by
intersection. Missing, ambiguous, suspended, revoked, expired or denied evidence
fails closed. Material, denied or privileged decisions use the registered
`authorization-decision/v1` profile when required; every evaluation accepts its
protected audit obligation before returning a result.

The RoleAssignment controller is the sole assignment publisher and sole writer
of its specification, lifecycle, status, revocation, replacement, expiry and
review-due effects. Administrators and grant, approval, privileged-access and
review controllers may submit exact immutable intents but cannot directly write
any part of a RoleAssignment. Delegated publication requires independently
applicable per-action witnesses covering the proposal's scope, target-binding
mode, optional exact resource, validity and delegation capability; authority
over one resource cannot grant another resource or scope-wide reach.

### 8.2 Zero-trust invariants

- Deny by default and authorize every operation at the target resource and repository boundary.
- Give every user, controller, plugin and adapter a distinct identity.
- Use the smallest practical role, scope and validity period.
- Require stronger authentication and separation of duties for sensitive operations.
- Use target-scoped, short-lived execution credentials where the backend permits.
- Never provide target credentials to managed-service customers.
- Record authorization, approval, execution and material policy decisions immutably.
- Authorization depends only on authenticated identity, exact registered actions and targets, scoped assignments, explicit dependencies and current guardrails—not persona names, UI labels, job titles or unprovisioned external group claims.
- A native ExecutionTarget effect requires both Sovrunn authorization and independent native-IAM authorization of a least-privileged adapter/controller identity. Neither layer's Allow substitutes for the other and neither overrides a Deny.

## 9. Governance policy model

| Contract or record | Purpose |
|---|---|
| `GovernanceProfile` | Versioned v1 composition of only FEATURE-0018-owned privileged-access, access-review, exception and audit rule references and constraints |
| `ProfileAssignment` | Applies a versioned profile to CloudProvider, portfolio, organization, tenant, project or service scope |
| `EffectiveGovernanceContext` | Immutable resolution of applicable profiles, entitlements, restrictions and approved exceptions for one decision or operation |
| `ApprovalPolicy` | Defines actions requiring approval, qualified approvers and separation-of-duties rules |
| `ApprovalRequest` | Records the request, rationale, reviewers, decision and validity |
| `ExceptionGrant` | Narrow, time-bound deviation from an explicitly exceptionable control, with compensating controls and evidence |
| `DecisionRecord` | Explains why a request was allowed, denied or constrained and records the policy versions used |
| `AuditEvent` | Records actor, action, scope, time, result and correlation |

### 9.1 Policy resolution

FEATURE-0018 owns the `GovernanceProfile` envelope and its v1 rule-component
semantics only. FEATURE-0020 owns profile assignment, inheritance, conflict and
exception application, and construction of `EffectiveGovernanceContext`.
Security, sovereignty, placement, backup, cost and other future-domain semantics
remain with their owning features and are not activated by FEATURE-0018.

- Allowed sets intersect.
- Prohibitions accumulate.
- Stronger minimum assurance requirements win.
- Shorter evidence-freshness limits win.
- Customer scopes may narrow provider permission but may not enlarge it.
- An exception may override only a control explicitly marked exceptionable.
- Unresolved conflicts or missing mandatory evidence deny the operation.
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

## Source: `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`

Approved SHA-256: `a74e2e54f974ceb974bb203acb28d239c07a03559e52b56b770a0292e347f772`

### Canonical FEATURE-0018 contract ledger and composition path

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->

`Persona` is a non-normative description of job intent, responsibilities and user journey. It is not an API resource, identity kind, role, scope or authorization input. Persona documentation may map a journey to one or more default role templates; actual authority exists only through the contracts below.

| Concept / API kind | Purpose and ownership | Profile, scope and writers | Mutability and lifecycle | Boundary, relationships and invariants |
|---|---|---|---|---|
| `PrincipalRef` | Embedded stable reference to an already-authenticated Human, Workload or System identity | EV only; identity lifecycle remains external; resolver validates issuer/subject and current identity state | Immutable value; identity change produces a new reference | Email/display name are not authority; unresolved identity fails closed |
| `AccessGroup` | Non-authenticating, access-only authorization subject | MR; Organization or CloudProvider scope; authorized administrator writes intent; AccessGroup controller writes status | Owner/display intent may update; suspend/resume/retire explicit | Direct auditable principal members only; not a ScopeKind; external, nested and dynamic groups grant nothing |
| `Membership` | Records Organization/CloudProvider belonging or direct same-scope AccessGroup membership without granting permission or workflow eligibility | MR; Organization or CloudProvider scope; administrator/trusted provisioner submits intent; membership controller alone validates grant envelope and writes canonical state | Relationship immutable; synchronize only system-owned provenance/freshness; suspend/reactivate/revoke; Guest expiry terminal | Exact relationship-key uniqueness; Standard or Guest; an AccessGroup assignment-effect expansion derives `MembershipEnabledGrantEnvelope` and passes complete non-synthesized action/scope/target/resource/validity dominance; privilege reduction remains available |
| `RoleDefinition` | Versioned permission template containing only registered canonical Sovrunn actions | VD; Platform, Organization or CloudProvider scope; security administrator drafts; publication controller owns lifecycle status | Draft mutable; published version immutable; supersede, emergency suspend/restore, retire; retained while referenced | Privilege classification is monotonic; no backend IAM action, persona, arbitrary condition or wildcard authority |
| `RoleAssignment` | Grants one exact published RoleDefinition version to one `RoleHolderRef` at one registered scope, optionally narrowed to one exact resource | MR; seven registered resource scopes; accepted workflows submit exact intents; RoleAssignment controller alone writes spec/status/lifecycle and every revoke/replace/expiry/review-due effect | Immutable published spec; explicit revoke/expiry/review effects; replacement creates a new assignment | `Standing` or `TimeBound`; exact Membership/responsibility/review pins; grantor ceiling includes per-action scope/target/resource-reach/validity/delegation dominance; no cross-assignment synthesis or workflow write |
| `PrivilegedAccessRequest` | Requests one justified individual JIT or break-glass temporary grant | LRO; seven registered scopes; requester writes accepted intent; privileged controller owns state | Immutable request, activation deadline and duration; cancel/block/activate/expire/revoke retained | JIT and break-glass activation require fresh phishing-resistant AAL2-or-higher assurance, non-empty exact-Human PrincipalRef-only `EligibilityRef` requester list, separate submit authority, separation of duties, no-grace expiry and linked review where required |
| `ApprovalPolicy` | Versioned bounded stages, eligibility, quorum, expiry and separation-of-duties policy | VD; Platform, Organization or CloudProvider scope; security/governance administrator drafts; publication controller owns lifecycle | Draft mutable; published immutable; supersede/retire; retained while referenced | Non-empty exact-Human PrincipalRef-only `EligibilityRef` approver list with separate decide authority and decision/effect recheck; never authorizes break-glass bypass; operation-local `ApprovalRequirement` selects one exact version |
| `ApprovalRequest` | Retains approval evidence for one bounded FEATURE-0018 intent or immutable proposal | LRO; seven registered scopes; originating controller materializes; approval controller owns result | Immutable subject/proposal; exactly one terminal decision; cancel/expire retained | Approval is current at downstream publication, then immutable provenance rather than continuing bearer authority |
| `AccessReview` | Reviews an immutable access snapshot with usage and remediation evidence | LRO; seven registered scopes; campaign owner or mandated controller originates; review controller owns state | Snapshot immutable after start; item decisions/completion/remediation retained | StandingCertification, Manual or BreakGlassRetrospective variants; non-empty exact-Human PrincipalRef-only `EligibilityRef` reviewer list with independent action and decision/effect recheck; only exact-version StandingCertification Retain advances `nextReviewDueAt`; beneficiary conflict set includes accountable Humans behind direct Workload/System AccessGroup members and prevents direct or indirect self-certification |
| `ExceptionGrant` | Immutable approved exception grant or linked revocation evidence | IR; seven registered scopes; exception controller is sole writer after approval | Append-only Grant/Revoke evidence; expiry is time-effective without mutation | Exact control/subject/scope/interval; deterministic overlap and narrowing; FEATURE-0020 alone applies it to effective governance |
| `AuthorizationInput` / `AuthorizationResult` | Transient exact authorization request and non-bearer `Allow | Deny` result | TRR; authorization controller writes result; exact provenance retained in required evidence | Evaluation-instant result only; never reusable authority | Grant union plus guardrail/dependency intersection; exact target binding; safe denial and audit-obligation acceptance |

The remaining closed FEATURE-0018 supporting-value inventory is
`RoleHolderRef`, `AssuranceEvidence`, `ApprovalRequirement`, `EligibilityRef`,
`RoleAssignmentProposal`, `MembershipEnabledGrantEnvelope`, `ExceptionProposal`, `UsageEvidenceSummary`,
`GovernanceApplicability`, and `ActionTargetBinding`. Together with
`AuthorizationInput` and `AuthorizationResult`, these are EmbeddedValue or
TransientRequestResult contracts only: no endpoint, independent lifecycle,
storage authority or user journey exists.

Scope applicability is closed per version and consumed generically. AccessGroup
and Membership allow Organization/CloudProvider; RoleDefinition,
ApprovalPolicy and GovernanceProfile allow Platform/Organization/CloudProvider;
RoleAssignment, PrivilegedAccessRequest, AccessReview, ApprovalRequest and
ExceptionGrant allow all seven canonical ScopeKinds. Customer policy may narrow
but not expand this registry. A new ScopeKind or containment rule is a core
architecture change. Exact publication-scope/reference compatibility and the
complete operation/action/target-binding registries remain F18-RD-02/F18-RD-06
authority and must be transcribed without inference into executable schemas.

Canonical processing path:

```text
Persona documentation
    → maps journey to default RoleDefinition templates

Authenticated PrincipalRef
  + canonical Sovrunn action
  + target resource and resolved scope
  + pinned RoleAssignment and RoleDefinition version
  + Membership/responsibility, policy guardrails, assurance and current time
      → AuthorizationResult
          → allow or deny
          → material DecisionRecord profile and/or AuditEvent
```

Authorization must never branch on persona name, UI label, organizational job title, email address, unprovisioned external group claim or backend IAM role. Default roles are packaged conveniences, not permanent enums; providers and customers may create narrower roles but may never grant actions outside their own delegated authority. Applicable grants combine by union and every guardrail/dependency combines by intersection; missing or ambiguous evidence denies.

An AccessGroup Membership assignment-effect expansion is grant-producing even
though Membership itself contributes no action. At its atomic publication boundary,
the Membership controller derives the exact current group-held assignment
action/reach/validity envelope and applies F18-RD-09's Membership-operation,
`roleassignment.grant`, per-action witness, temporal-dominance, no-synthesis,
trusted-provisioner and coherent-snapshot rules. No user-authored Membership
field stores or supplies that authority. Membership cannot establish approver,
reviewer, JIT, or break-glass eligibility. All three eligibility lists use the
same non-empty `EligibilityRef` OR contract containing exact Human
PrincipalRefs only. RoleDefinition, RoleAssignment, AccessGroupRef, Membership,
claim, and other relationships never qualify in v1. Eligibility remains
separate from the required canonical action; exact policy/rule applicability
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

### Canonical GovernanceProfile contract boundary

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
## 9. Governance contracts

### 9.1 Contract matrix

| Concept / API kind | Purpose and ownership | Profile, scope and writers | Mutability and lifecycle | Boundary, relationships and invariants |
|---|---|---|---|---|
| `GovernanceProfile` | v1 composition of only FEATURE-0018-owned privileged-access, access-review, exception and audit rule references/constraints | VD; Platform, CloudProvider or Organization scope; governance administrator drafts; publication controller owns lifecycle | Draft mutable; published immutable; supersession blocks new selection while preserving valid pins; retirement blocks new/current-rule use; retained for audit | GO canonical with safe summaries; no security, sovereignty, placement, backup, cost or other future-domain semantics; FEATURE-0020 owns assignment and effective resolution; DEC-0060/F18-RD-18 |
| `ProfileAssignment` | Applies a specific profile version to a target scope/resource | MR; same scope as target; authorized scope policy administrator owns spec; policy resolver owns status | ProfileRef and targetRef immutable; validity mutable by governed operation; RESTRICT active decision dependencies | GO/OF; allowed sets intersect, prohibitions accumulate, strongest minimum and shortest evidence age win |
| `EffectiveGovernanceContext` | Immutable resolved governance, entitlement, restriction and approved-exception snapshot for one decision or operation | IR; same canonical scope as subject; policy resolver is sole writer | Append-only; supersession creates new record; RETAIN | IE/GO; supersedes the planned `EffectivePolicyContext` term; contains references/versions, not copied secrets; exact contextRef is mandatory for authoritative decisions; `ADH-2026-033` |

### 9.2 Skeletal YAML

```yaml
apiVersion: policy.sovrunn.io/v1alpha1
kind: GovernanceProfile
metadata:
  name: government-production-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
spec:
  version: v1
  approvalPolicyRefs: [government-production-approval-v1]
  privilegedAccessRules: [government-production-privileged-v1]
  accessReviewRules: [government-production-standing-review-v1]
  exceptionRules: [government-production-exception-v1]
  auditRequirements: [government-production-iam-audit-v1]
  publicationState: Published
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

### Canonical writers and retention boundaries

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
## 15. Canonical creation and mutation authority

| Actor | May author | Must not author |
|---|---|---|
| Platform architecture/operator | Platform-scoped definitions, CloudProvider registrations, global policy registries | Customer service intent or provider-native facts without authority |
| Approved platform release publisher | SovrunnRelease drafts and signed published versions within delegated scope | Installation desired state, lifecycle approval, plans, execution or health status |
| Release compatibility authority | ReleaseCompatibilityContract, ComponentManifest and approved ValidationGateDefinition drafts/publication within delegated scope | Installation desired state, lifecycle approval or execution status |
| Platform lifecycle administrator | SovrunnInstallation desired release/policy intent and authorized platform Operation requests | Published release contents, immutable lifecycle plans, approval decisions, agent execution status or underlying infrastructure lifecycle |
| Platform lifecycle planner | Immutable PlatformLifecyclePlan derived from exact release, policy, approval, backup, compatibility and health inputs | Desired installation state, approval authority or execution status |
| Platform lifecycle agent | Execution status and protected handles for an approved PlatformLifecyclePlan and correlated Operation | Release publishing, desired state, lifecycle policy, approval, authorization decisions or unplanned underlying infrastructure changes |
| Infrastructure operator or trusted adapter | InfrastructureMaintenanceNotice desired facts within delegated CloudProviderParticipation/target authority | ExecutionTarget availability/qualification status, owner restrictions or service placement decisions |
| ExecutionTarget lifecycle controller | Target availability state, maintenance epoch/fence and requalification status | Native maintenance facts, customer intent or provider agreement state |
| (retired) Canonical migration controller | None — role retired by DEC-0059; no runtime migration authority exists | Transformation outputs, CanonicalMigrationRecord, or any migration-plan/record write — all removed from active Phase 2R authorities |
| Approved service-contract publisher | ServiceTypeDefinition, ServiceRelationshipDefinition and, when applicable, ServiceCompositionDefinition drafts and published versions within delegated scope | Provider commercial terms, customer instances, decisions or execution state |
| CloudPlatform administrator | Platform catalog, enrollments, entitlements, quota envelope, approved product/governance/placement profiles and participation requests within delegated scope | Provider credentials or native operations, customer Organization ownership or customer Project contents |
| CloudProvider administrator | Participation response, topology declarations, realization mappings, maintenance notices, approved operating profiles and targets within delegated scope | CloudPlatform catalog/enrollment grants, customer Organization ownership or customer Project contents |
| Customer organization administrator | Organization hierarchy, membership/policy assignments and customer restrictions | Provider grants, target facts or wider entitlement |
| Authorized identity/security administrator | Membership intent, RoleDefinition drafts/publication requests, exact scoped RoleAssignment grant intents and authorized AccessReview campaign intents within delegated authority | Direct RoleAssignment publication/status, persona-based grants, permissions outside delegated scope, authorization outcomes or silent mutation of published role versions |
| RoleAssignment controller | Publish/revoke/expire and write status for the exact accepted immutable RoleAssignment intent | Selecting or enlarging holder, role, scope, resource, validity or delegation; originating workflow semantics |
| Project user | ServiceInstance intent, authorized ServiceRelationship intent and ServiceBinding requests within entitlement, quota and role | Status, decisions, evidence, deployment plans, raw credentials, native infrastructure objects or execution credentials |
| Policy resolver | EffectiveGovernanceContext and policy status | Source profiles or customer intent |
| Evidence collector/adapter | SovereigntyFactSet observations and EvidenceRecords within authorized subjects | Desired customer state or decision authority |
| Platform dependency resolver | Immutable PlatformDependencySnapshot from exact installation and dependency facts | Dependency desired state, evidence contents, policy or assessment outcome |
| Sovereignty/placement decision service | FEATURE-0013 DecisionRecords under registered profiles | Policy inputs, facts, evidence or execution state |
| Deployment planner | ServiceDeploymentPlan | Placement authority or backend execution |
| Operation controller | Operation status and orchestration | Customer request fields after acceptance |
| Service plugin/adapter | PluginExecution status and protected backend handles | Entitlement, governance or placement decisions |
| Projection controller | ServicePlacement customer-safe record | Canonical decisions or sensitive provider topology |
| Binding controller | ServiceBinding status, safe endpoint metadata and typed SecretRef | Secret values, provider target credentials or unrelated consumer bindings |
| Relationship controller | ServiceRelationship status, safe rationale and decision/operation references | Source definitions, customer intent, canonical decisions, bindings, secret values or implementation-native handles |

## 16. Cross-resource deletion rules

1. A CloudProvider cannot be deleted while products, enrollments, targets or retained records reference it.
2. Published definitions are retired, not mutated or immediately deleted.
3. A ServiceTypeDefinition or ServiceCompositionDefinition version cannot be removed while an offering, runtime profile, instance history, decision or retained record references it.
4. A ServiceRelationshipDefinition version cannot be removed while a ServiceRelationship, requirement set, decision, plan, operation or retained history references it.
5. Terminating an enrollment does not delete the customer Organization or its hierarchy.
6. Revoking entitlement does not silently delete running services; a governed transition policy decides deny-new, grace, migrate or terminate.
7. Deleting a Project or ServiceInstance requires an impact preview and explicit lifecycle operations.
8. A ServiceInstance cannot finalize deletion while an active required ServiceRelationship or ServiceBinding depends on it; relationships must follow declared lifecycle policy and bindings must be revoked first.
9. Datacenters, stacks and targets must be drained/decommissioned before removal and cannot be removed while active placement records depend on them.
10. Effective contexts, decisions, evidence, deployment plans, operations and placement projections obey retention/legal-hold rules and are never cascade-deleted with a mutable parent.
11. Correction of an immutable record creates a linked correction, supersession or revocation record.
12. FEATURE-0018 resources have no hard-delete surface. AccessGroup and Membership use explicit suspend/reactivate/revoke/retire behavior; RoleAssignment uses explicit revoke/expiry/review effects; published definitions retire and remain resolvable while referenced; ApprovalRequest, AccessReview and ExceptionGrant evidence follows retention/legal hold. No parent deletion or workflow outcome cascade-deletes authorization or audit history.
13. A SovrunnInstallation cannot finalize decommissioning until required backups and restore evidence are current, active platform Operations are terminal and control-plane impact has been approved.
14. A published SovrunnRelease or PlatformLifecyclePolicy version cannot be removed while an installation, plan, operation, evidence or retained audit record references it.
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

## Source: `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`

Approved SHA-256: `461d1c59e544b248dbb9b5fd0bc9471c6e3fbc9d0ce7eae4f07f8787a667b836`

### FEATURE-0018 schema, writer and cross-scope registry entries

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
  - {id: VS0-SCHEMA-020, identity: "iam.sovrunn.io/v1alpha1/PrincipalRef", owner: FEATURE-0018, profile: EmbeddedValue, boundary: internal-engine-facing, scopes: [], required: ["issuer:string", "subject:string", "principalType:enum[Human,Workload,System]"], optional: [], refs: {}, classification: Tenant-confidential, redaction: redact, mutability: immutable, retention: containing-contract, slice0Constraint: "identity lifecycle external; unresolved principal fails closed; DEC-0060 Accepted/F18-RD-02/03"}
  - {id: VS0-SCHEMA-021, identity: "iam.sovrunn.io/v1alpha1/Membership", owner: FEATURE-0018, profile: ManagedResource, boundary: operator-facing, scopes: [Organization, CloudProvider], required: ["spec.principalRef:PrincipalRef", "spec.contextKind:enum[Organization,CloudProvider,AccessGroup]", "spec.membershipType:enum[Standard,Guest]", "protected.sourceAuthority", "protected.provisionedAt:RFC3339", "status.phase:enum[Active,Suspended,Revoked]"], optional: ["spec.accessGroupRef:AccessGroupRef", "spec.expiresAt:RFC3339", "spec.sponsorRef:PrincipalRef<Human>", "spec.responsiblePartyRef:PrincipalRef<Human>", "protected.sourceObjectRef", "protected.lastSynchronizedAt:RFC3339", "protected.freshUntil:RFC3339", "status.eligibilityProjection:enum[Eligible,Stale,GuestExpired]"], refs: {spec.accessGroupRef: [AccessGroup]}, classification: Tenant-confidential, redaction: redact, mutability: "relationship spec immutable; trusted provisioner/admin submits intent; membership controller solely validates any membership-enabled grant envelope and writes lifecycle/protected provenance/freshness", retention: audit-policy, slice0Constraint: "no roleAssignmentRefs; grants no action or approval/review/privileged eligibility; AccessGroup assignment-effect expansion requires exact MembershipEnabledGrantEnvelope dominance with no witness synthesis; GuestExpired terminal while Stale is synchronization-recoverable; DEC-0060 Accepted/F18-RD-04/05/09/12/14/16"}
  - {id: VS0-SCHEMA-022, identity: "iam.sovrunn.io/v1alpha1/RoleDefinition", owner: FEATURE-0018, profile: VersionedDefinition, boundary: governance-only, scopes: [Platform, Organization, CloudProvider], required: ["spec.version", "spec.actions:registered-action[](non-empty)", "spec.publicationState:enum[Draft,Published,Suspended,Retired]"], optional: ["spec.supersedesRef"], refs: {}, classification: Internal, redaction: redact, mutability: "draft mutable; published immutable; publication-controller lifecycle", retention: archival, slice0Constraint: "no conditions/wildcards/provider-native actions; DEC-0060 Accepted/F18-RD-06/07"}
  - {id: VS0-SCHEMA-023, identity: "iam.sovrunn.io/v1alpha1/RoleAssignment", owner: FEATURE-0018, profile: ManagedResource, boundary: governance-only, scopes: [Platform, CloudPlatform, CloudProvider, Organization, OrganizationUnit, Tenant, Project], required: ["spec.roleHolderRef:RoleHolderRef", "spec.roleDefinitionRef:TypedRef<RoleDefinition>(uid+version pinned)", "spec.validity:Standing|TimeBound", "status.phase"], optional: ["spec.resourceRef:TypedRef(exact target)", "spec.membershipOrResponsibilityPins", "spec.standingReviewRulePin"], refs: {spec.roleDefinitionRef: [RoleDefinition]}, classification: Tenant-confidential, redaction: redact, mutability: "immutable published spec; RoleAssignment controller sole spec/status/lifecycle/effect writer", retention: audit-policy, slice0Constraint: "no conditions or competing scope field; delegated grant requires per-action scope/target/resource/time/delegation reach containment with no cross-assignment synthesis; DEC-0060 Accepted/F18-RD-08/09"}
  - {id: VS0-SCHEMA-024, identity: "governance.sovrunn.io/v1alpha1/GovernanceProfile", owner: FEATURE-0018, profile: VersionedDefinition, boundary: governance-only, scopes: [Platform, Organization, CloudProvider], required: ["spec.version", "spec.publicationState:enum[Draft,Published,Retired]"], optional: ["spec.approvalPolicyRefs", "spec.privilegedAccessRules", "spec.accessReviewRules", "spec.exceptionRules", "spec.auditRequirements"], refs: {spec.approvalPolicyRefs: [ApprovalPolicy]}, classification: Internal, redaction: redact, mutability: "draft mutable; published immutable; publication-controller lifecycle", retention: archival, slice0Constraint: "FEATURE-0018 components only; FEATURE-0020 owns resolution; DEC-0060 Accepted/F18-RD-18"}
  - {id: VS0-SCHEMA-062, identity: "iam.sovrunn.io/v1alpha1/AccessGroup", owner: FEATURE-0018, profile: ManagedResource, boundary: operator-facing, scopes: [Organization, CloudProvider], slice0Constraint: "access-only; explicit owner; direct principal membership; no nesting/dynamic/external-claim authority; DEC-0060 Accepted/F18-RD-04"}
  - {id: VS0-SCHEMA-063, identity: "iam.sovrunn.io/v1alpha1/PrivilegedAccessRequest", owner: FEATURE-0018, profile: LongRunningOperation, boundary: governance-only, scopes: "$metadata.scopeKinds", slice0Constraint: "JIT|BreakGlass; immutable activationDeadline/duration; activation requires phishing-resistant AAL2-or-higher evidence age<=15m; individual TimeBound effect; DEC-0060 Accepted/F18-RD-14/15"}
  - {id: VS0-SCHEMA-064, identity: "iam.sovrunn.io/v1alpha1/AccessReview", owner: FEATURE-0018, profile: LongRunningOperation, boundary: governance-only, scopes: "$metadata.scopeKinds", slice0Constraint: "immutable snapshot/usage evidence; StandingCertification|Manual|BreakGlassRetrospective; direct/indirect beneficiary conflict set includes accountable Humans behind direct Workload/System AccessGroup members and is rechecked at atomic disposition/effect publication; DEC-0060 Accepted/F18-RD-16"}
  - {id: VS0-SCHEMA-065, identity: "iam.sovrunn.io/v1alpha1/ApprovalPolicy", owner: FEATURE-0018, profile: VersionedDefinition, boundary: governance-only, scopes: [Platform, Organization, CloudProvider], slice0Constraint: "bounded stages/eligibility/quorum/expiry/SoD; no break-glass bypass; DEC-0060 Accepted/F18-RD-12"}
  - {id: VS0-SCHEMA-066, identity: "iam.sovrunn.io/v1alpha1/ApprovalRequest", owner: FEATURE-0018, profile: LongRunningOperation, boundary: governance-only, scopes: "$metadata.scopeKinds", slice0Constraint: "immutable bounded subject/proposal; admission evidence; exactly-once terminal decision; DEC-0060 Accepted/F18-RD-13"}
  - {id: VS0-SCHEMA-067, identity: "governance.sovrunn.io/v1alpha1/ExceptionGrant", owner: FEATURE-0018, profile: ImmutableRecord, boundary: governance-only, scopes: "$metadata.scopeKinds", slice0Constraint: "exact control/subject/scope/interval and linked revocation evidence; FEATURE-0020 resolves effect; DEC-0060 Accepted/F18-RD-17"}
  - {id: VS0-SCHEMA-068, identity: "F18:SupportingValues", owner: FEATURE-0018, profile: "EmbeddedValue|TransientRequestResult", boundary: shared, scopes: [], slice0Constraint: "closed F18-RD-02 inventory including EligibilityRef and operation-local MembershipEnabledGrantEnvelope; EligibilityRef contains exactly one Human PrincipalRef and roles, assignments, groups, Membership or claims never qualify; no endpoint/lifecycle/storage authority"}
  - {id: VS0-WRITER-008, paths: ["RoleAssignment.spec", "RoleAssignment.status"], writer: RoleAssignment-controller, forbidden: [all-clients, delegated-identity-governance-admin, grant-controller, approval-controller, privileged-access-controller, access-review-controller, adapters], conflict: AUTHORIZATION_DENIED, note: "sole canonical publisher/writer for create, revoke, replace, expiry and review-due effects; other authorities submit only exact immutable intents; DEC-0060 Accepted/F18-RD-02/08"}
  - {id: VS0-WRITER-022, paths: ["Membership.spec", "Membership.protected", "Membership.status"], writer: identity-membership-controller, forbidden: [all-clients, adapters, persona-label], conflict: AUTHORIZATION_DENIED, note: "administrators/trusted provisioners submit exact intent only; controller validates MembershipEnabledGrantEnvelope for AccessGroup assignment-effect expansion and writes canonical relationship, protected provenance/freshness, lifecycle and projection; Membership cannot establish approval/review/privileged eligibility; DEC-0060 Accepted/F18-RD-02/05/09/12/14/16"}
  - {id: VS0-WRITER-023, paths: ["ProfileAssignment.spec"], writer: delegated-governance-admin, forbidden: [persona-label, project-user-outside-delegation, adapters], conflict: AUTHORIZATION_DENIED, note: "preserves existing FEATURE-0020 authority; FEATURE-0018 assigns no ProfileAssignment semantics"}
  - {id: VS0-CF-F01, owner: FEATURE-0018, inputs: unauthenticated-create, expectedState: unchanged, expectedError: AUTH_REQUIRED, expectedSideEffects: none, gate: security}
  - {id: VS0-CF-F02, owner: FEATURE-0018, inputs: unauthorized-project-create, expectedState: unchanged, expectedError: RESOURCE_NOT_FOUND, expectedSideEffects: none, gate: security}
  - {id: VS0-CF-X01, owner: FEATURE-0018, inputs: cross-organization-ref, expectedState: unchanged, expectedError: RESOURCE_NOT_FOUND, expectedSideEffects: none, gate: security}
  - {id: VS0-CF-X02, owner: FEATURE-0018, inputs: cross-project-ref, expectedState: unchanged, expectedError: RESOURCE_NOT_FOUND, expectedSideEffects: none, gate: security}
  - {id: "VS0-CF-F18-01", owner: "FEATURE-0018", inputs: "Active employee, published Developer role, active project assignment; same action at unrelated project", expectedState: "requested project Allow; unrelated project unchanged", expectedError: "none; unrelated project Deny", expectedSideEffects: "exact protected authorization evidence and audit; no unrelated-scope effect", gate: "security"}
  - {id: "VS0-CF-F18-02", owner: "FEATURE-0018", inputs: "Same issuer/subject principal with changed email display attribute", expectedState: "principal and existing authority remain bound to issuer/subject", expectedError: null, expectedSideEffects: "attribution retains stable principal identity", gate: "identity"}
  - {id: "VS0-CF-F18-03", owner: "FEATURE-0018", inputs: "Active assignment whose required Membership/principal becomes suspended", expectedState: "assignment retained but immediately ineffective", expectedError: "Deny", expectedSideEffects: "retained history and protected denial audit; no access effect", gate: "security"}
  - {id: "VS0-CF-F18-04", owner: "FEATURE-0018", inputs: "Revoked/expired relationship followed by a later rejoin", expectedState: "new relationship does not reactivate old assignment", expectedError: "Deny until newly authorized grant", expectedSideEffects: "retained old evidence; no silent reactivation", gate: "lifecycle"}
  - {id: "VS0-CF-F18-05", owner: "FEATURE-0018", inputs: "Workload principal with responsible Human and exact scoped assignment invokes one action", expectedState: "applicable action may Allow; unrelated action/scope denies", expectedError: "none or Deny by fixture", expectedSideEffects: "independent workload attribution and protected audit", gate: "security"}
  - {id: "VS0-CF-F18-06", owner: "FEATURE-0018", inputs: "System controller submits an owned intent and separately attempts a non-owned write", expectedState: "owned controller transition may publish; non-owned write unchanged", expectedError: "none or AUTHORIZATION_DENIED", expectedSideEffects: "sole writer publishes; direct competing write has no effect", gate: "writer"}
  - {id: "VS0-CF-F18-07", owner: "FEATURE-0018", inputs: "Consultant has distinct supplier and customer Guest Membership fixtures", expectedState: "only exact active target context is eligible until exact expiry", expectedError: "none or existing safe error", expectedSideEffects: "sponsorship, scope, expiry, and responsible Human remain attributable", gate: "security"}
  - {id: "VS0-CF-F18-08", owner: "FEATURE-0018", inputs: "Raw external group claim; then trusted-provisioned AccessGroup, current direct Membership, and active scoped assignment", expectedState: "raw claim grants nothing; canonical relationship may allow", expectedError: "Deny; then none when all canonical evidence applies", expectedSideEffects: "no native/external group authority; canonical provenance audited", gate: "security"}
  - {id: "VS0-CF-F18-09", owner: "FEATURE-0018", inputs: "Administrator delegates project action/reach within and beyond one independently applicable per-action witness", expectedState: "contained grant intent may publish; broader intent unchanged", expectedError: "none or AUTHORIZATION_DENIED", expectedSideEffects: "RoleAssignment controller is sole publisher; no witness synthesis", gate: "security"}
  - {id: "VS0-CF-F18-10", owner: "FEATURE-0018", inputs: "Exact FEATURE-0016 ExecutionTarget read and qualify bindings for one accessible target", expectedState: "read/qualify apply only to registered target binding", expectedError: null, expectedSideEffects: "protected action/target provenance; no provider-native effect", gate: "integration-contract"}
  - {id: "VS0-CF-F18-11", owner: "FEATURE-0018", inputs: "Qualifier-only holder attempts ExecutionTarget retirement", expectedState: "target remains unchanged", expectedError: "AUTHORIZATION_DENIED", expectedSideEffects: "denial audit; no retirement", gate: "security"}
  - {id: "VS0-CF-F18-12", owner: "FEATURE-0018", inputs: "Holder requests create-only ExecutionTarget authority under combined existing write action", expectedState: "no invented create-only grant/action is published", expectedError: "existing safe error", expectedSideEffects: "recorded granularity limitation; no action-registry mutation", gate: "compatibility"}
  - {id: "VS0-CF-F18-13", owner: "FEATURE-0018", inputs: "Exact Human requests approved two-hour JIT production DB role before activation deadline", expectedState: "one TimeBound RoleAssignment becomes active then ineffective at exact expiry", expectedError: null, expectedSideEffects: "immutable approval/activation provenance and required audit", gate: "privileged-access"}
  - {id: "VS0-CF-F18-14", owner: "FEATURE-0018", inputs: "Requester attempts to approve own privileged request", expectedState: "request not approved or activated", expectedError: "AUTHORIZATION_DENIED", expectedSideEffects: "denied decision audited; no RoleAssignment intent/publication", gate: "security"}
  - {id: "VS0-CF-F18-15", owner: "FEATURE-0018", inputs: "Listed approver loses current eligibility before decision acceptance", expectedState: "request remains without that accepted approval", expectedError: "AUTHORIZATION_DENIED", expectedSideEffects: "denial audited; no downstream effect", gate: "security"}
  - {id: "VS0-CF-F18-16", owner: "FEATURE-0018", inputs: "Approval decision races ApprovalRequest expiry", expectedState: "exactly one valid terminal state wins; no unauthorized assignment", expectedError: "none for winner; CONFLICT for loser", expectedSideEffects: "one terminal evidence publication; no duplicate effect", gate: "concurrency"}
  - {id: "VS0-CF-F18-17", owner: "FEATURE-0018", inputs: "Same approved JIT request activated twice with same and conflicting replay fixtures", expectedState: "exactly one temporary RoleAssignment; validity never extends", expectedError: "safe replay or CONFLICT", expectedSideEffects: "one activation publication and one completed result", gate: "idempotency"}
  - {id: "VS0-CF-F18-18", owner: "FEATURE-0018", inputs: "Eligible Human invokes approved break-glass path with current assurance evidence", expectedState: "individual TimeBound access lasts no more than one hour", expectedError: null, expectedSideEffects: "immediate protected audit and retrospective review due within 24 hours", gate: "privileged-access"}
  - {id: "VS0-CF-F18-19", owner: "FEATURE-0018", inputs: "Standing certification finds stale administrator assignment and selects Revoke", expectedState: "exact assignment becomes revoked exactly once", expectedError: null, expectedSideEffects: "retained immutable review, usage, decision, and remediation evidence", gate: "review"}
  - {id: "VS0-CF-F18-20", owner: "FEATURE-0018", inputs: "Reviewed assignment version changes after immutable snapshot", expectedState: "review becomes Stale; assignment unchanged", expectedError: null, expectedSideEffects: "no remediation or review-due advancement", gate: "concurrency"}
  - {id: "VS0-CF-F18-21", owner: "FEATURE-0018", inputs: "Incident action requires both applicable DB role and bounded change-window exception", expectedState: "Allow only when both independent authorities are current", expectedError: "none or Deny", expectedSideEffects: "exact role and exception provenance; neither synthesizes the other", gate: "security"}
  - {id: "VS0-CF-F18-22", owner: "FEATURE-0018", inputs: "Proposal attempts to except registered non-exceptionable tenant isolation control", expectedState: "no ExceptionGrant published", expectedError: "Deny using existing registered outcome", expectedSideEffects: "terminal protected denial evidence and audit", gate: "security"}
  - {id: "VS0-CF-F18-23", owner: "FEATURE-0018", inputs: "Active ExceptionGrant receives valid linked early-revocation proposal", expectedState: "original record remains immutable and exact grant becomes ineffective at revoke instant", expectedError: null, expectedSideEffects: "one linked effect=Revoke record and audit", gate: "exception"}
  - {id: "VS0-CF-F18-24", owner: "FEATURE-0018", inputs: "Proposed ExceptionGrant overlaps same control/subject/scope/interval", expectedState: "existing grant unchanged; new grant absent", expectedError: "CONFLICT", expectedSideEffects: "no merge, precedence inference, or partial publication", gate: "exception"}
  - {id: "VS0-CF-F18-25", owner: "FEATURE-0018", inputs: "FEATURE-0017 returns Allow while applicable RoleAssignment is expired", expectedState: "final authorization is Deny", expectedError: "Deny", expectedSideEffects: "protected decision provenance records both inputs; no effect", gate: "security"}
  - {id: "VS0-CF-F18-26", owner: "FEATURE-0018", inputs: "Compatible pinned FEATURE-0017 evaluation returns RequiresApproval", expectedState: "no effect; request exists only through the explicit approved flow", expectedError: "RequiresApproval", expectedSideEffects: "protected evaluation/requirement provenance; no implicit ApprovalRequest", gate: "policy-seam"}
  - {id: "VS0-CF-F18-27", owner: "FEATURE-0018", inputs: "FEATURE-0017 result is Indeterminate", expectedState: "final authorization is Deny", expectedError: "Deny", expectedSideEffects: "redacted protected decision/audit; no effect", gate: "security"}
  - {id: "VS0-CF-F18-28", owner: "FEATURE-0018", inputs: "Required AuditEvent append fails at JIT activation publication boundary", expectedState: "request is not Active; no RoleAssignment or completed replay published", expectedError: "INTERNAL_ERROR", expectedSideEffects: "atomic rollback/no publication", gate: "audit"}
  - {id: "VS0-CF-F18-29", owner: "FEATURE-0018", inputs: "Authorized end customer performs a routine low-risk action", expectedState: "business result returned without governance ceremony", expectedError: null, expectedSideEffects: "identity/grants/guardrails/provenance/audit resolve internally", gate: "usability"}
  - {id: "VS0-CF-F18-30", owner: "FEATURE-0018", inputs: "Denied operation references an inaccessible target", expectedState: "protected state unchanged", expectedError: "RESOURCE_NOT_FOUND safe denial", expectedSideEffects: "redacted audit; no existence or policy-detail disclosure", gate: "security"}
  - {id: "VS0-CF-F18-31", owner: "FEATURE-0018", inputs: "Production Workload/System receives policy-permitted narrow Standing assignment", expectedState: "assignment effective, immediately revocable, and review-due", expectedError: null, expectedSideEffects: "responsible Human, provenance, and review schedule retained", gate: "review"}
  - {id: "VS0-CF-F18-32", owner: "FEATURE-0018", inputs: "Guest RoleAssignment omits notBefore or expiresAt", expectedState: "no assignment published", expectedError: "VALIDATION_FAILED", expectedSideEffects: "no audit-changing publication or completed replay", gate: "contract"}
  - {id: "VS0-CF-F18-33", owner: "FEATURE-0018", inputs: "Two grants allow action while a mandatory applicable guardrail excludes it", expectedState: "final authorization is Deny", expectedError: "Deny", expectedSideEffects: "exact grant and guardrail provenance; no effect", gate: "security"}
  - {id: "VS0-CF-F18-34", owner: "FEATURE-0018", inputs: "Published RoleDefinition version is suspended, then validly restored", expectedState: "pinned assignments ineffective during suspension; other versions unaffected", expectedError: "Deny during suspension", expectedSideEffects: "suspension/restoration DecisionRecord and audit; no assignment rewrite", gate: "lifecycle"}
  - {id: "VS0-CF-F18-35", owner: "FEATURE-0018", inputs: "AuthorizationResult is replayed for a different principal/action/target/scope", expectedState: "result grants nothing and target state remains unchanged", expectedError: "Deny", expectedSideEffects: "fresh authorization/audit path; no bearer reuse", gate: "security"}
  - {id: "VS0-CF-F18-36", owner: "FEATURE-0018", inputs: "Federated provenance is stale or responsible Human is absent at a required boundary", expectedState: "privilege-increasing operation unchanged", expectedError: "Deny", expectedSideEffects: "safe protected denial; raw claim cannot repair canonical evidence", gate: "security"}
  - {id: "VS0-CF-F18-37", owner: "FEATURE-0018", inputs: "AccessReview usage telemetry is absent, incomplete, and complete in separate fixtures", expectedState: "absence is not non-use; no silent workload revocation", expectedError: null, expectedSideEffects: "immutable coverage/usage evidence and explicit reviewer outcome", gate: "review"}
  - {id: "VS0-CF-F18-39", owner: "FEATURE-0018", inputs: "Administrator supplies holder, role version, scope, and validity; separately forges system-owned fields", expectedState: "valid exact intent may publish; forged request unchanged", expectedError: "none or AUTHORIZATION_DENIED", expectedSideEffects: "RoleAssignment controller supplies refs/provenance/audit as sole publisher", gate: "usability"}
  - {id: "VS0-CF-F18-40", owner: "FEATURE-0018", inputs: "Engineer submits exact time-bound privileged request", expectedState: "one requester-facing request status and expiry projection", expectedError: null, expectedSideEffects: "internal request/evidence/resource separation retained", gate: "usability"}
  - {id: "VS0-CF-F18-41", owner: "FEATURE-0018", inputs: "Exact eligible approver decides privileged/exception request; includes cross-stage reuse attempt", expectedState: "valid decision affects only current stage; one Human cannot count in two stages", expectedError: "none or AUTHORIZATION_DENIED", expectedSideEffects: "immutable decision evidence; downstream effect remains separate", gate: "approval"}
  - {id: "VS0-CF-F18-42", owner: "FEATURE-0018", inputs: "Eligible independent reviewer selects Retain, Revoke, or Replace on exact snapshot", expectedState: "exact approved disposition or explicit stale outcome", expectedError: "none or existing safe error", expectedSideEffects: "usage summary, beneficiary set, decision, and remediation evidence retained", gate: "review"}
  - {id: "VS0-CF-F18-43", owner: "FEATURE-0018", inputs: "User proposes bounded exceptionable control exception", expectedState: "one request status/effective-period projection; grant only after approved terminal processing", expectedError: "none or approved terminal Deny", expectedSideEffects: "proposal/approval/grant evidence remains internally distinct", gate: "usability"}
  - {id: "VS0-CF-F18-44", owner: "FEATURE-0018", inputs: "ExactResource(A) grantor attempts resource B and ScopeOnly delegation", expectedState: "no broader RoleAssignment published", expectedError: "AUTHORIZATION_DENIED", expectedSideEffects: "no cross-assignment synthesis or partial grant", gate: "security"}
  - {id: "VS0-CF-F18-45", owner: "FEATURE-0018", inputs: "Direct or indirect Standing beneficiary attempts Retain/Replace self-certification; unresolved expansion fixture included", expectedState: "no access-preserving/replacing effect or due advancement", expectedError: "REVIEWER_CONFLICT or REVIEW_BENEFICIARY_UNRESOLVED", expectedSideEffects: "independently authorized Revoke remains available and audited", gate: "security"}
  - {id: "VS0-CF-F18-46", owner: "FEATURE-0018", inputs: "TimeBound membership administrator attempts to enable durable group-held assignment", expectedState: "no durable derived access; empty group addition grants nothing", expectedError: "AUTHORIZATION_DENIED", expectedSideEffects: "exact envelope/witness denial evidence; no Membership publication when expansion exceeds ceiling", gate: "security"}
  - {id: "VS0-CF-F18-47", owner: "FEATURE-0018", inputs: "Membership administrator uses group/relationship/role evidence to manufacture approval, JIT, break-glass, or review eligibility", expectedState: "eligibility remains unsatisfied", expectedError: "existing structural-validation or safe-denial error", expectedSideEffects: "no ApprovalRequest decision, activation, review effect, or membership-derived eligibility", gate: "security"}
  - {id: "VS0-CF-F18-48", owner: "FEATURE-0018", inputs: "Reviewer list is empty, indirect, stale, scope-incompatible, target-incompatible, or changes before remediation", expectedState: "review effect unchanged", expectedError: "existing F18-RD-16 safe error", expectedSideEffects: "eligibility and beneficiary conflict rechecked; no remediation/due advancement", gate: "security"}
  - {id: "VS0-CF-F18-49", owner: "FEATURE-0018", inputs: "Role grantor or later rule publisher attempts to make a role/group/relationship qualify for eligibility", expectedState: "eligibility remains exact-Human-only and unsatisfied", expectedError: "existing structural-validation or safe-denial error", expectedSideEffects: "no approval, activation, Retain/Replace, or eligibility side effect", gate: "security"}
  - {id: "VS0-CF-F18-50", owner: "FEATURE-0018", inputs: "Active FEATURE-0018 contract fixture containing one future-owned/prohibited resource, field, adapter, credential flow, or provider-native IAM object", expectedState: "fixture rejected before evaluation/publication", expectedError: "existing architecture-drift or structural-validation failure", expectedSideEffects: "zero external I/O and no canonical/runtime publication", gate: "architecture"}
  - {id: "VS0-CF-F18-51", owner: "FEATURE-0018", inputs: "Valid exact GovernanceProfile v1 component/applicability fixture; then unknown, ambiguous, scope-incompatible, broadened, and retired-version fixtures", expectedState: "valid fixture accepted as pinned rule evidence only; invalid fixtures rejected", expectedError: "none or existing F18-RD-18 validation failure", expectedSideEffects: "no profile assignment/effective-resolution effect and no silent successor migration", gate: "contract"}
  - {id: "VS0-CF-F18-52", owner: "FEATURE-0018", inputs: "Exact three DecisionProfile registrations and mandatory authorization/approval/exception decision and AuditEvent fixtures, including required-append failure", expectedState: "only approved profiles/taxonomy accepted; required publication is atomic with evidence", expectedError: "none or INTERNAL_ERROR for required-append failure", expectedSideEffects: "exact immutable FEATURE-0013 evidence; no competing envelope or fourth profile", gate: "audit"}
  - {id: "VS0-CF-F18-53", owner: "FEATURE-0018", inputs: "Valid and invalid immutable synthetic fixture graphs with deterministic UTC clock and external-call counter", expectedState: "invalid graph rejected before evaluation; valid repeated evaluation is identical", expectedError: "none or existing fixture-validation failure", expectedSideEffects: "externalCallCount=0; no IdP, policy-engine, workflow-engine, CloudProvider, or network I/O", gate: "conformance"}
  - {id: "VS0-CF-F18-54", owner: "FEATURE-0018", inputs: "Standards matrix with every applicable invariant mapped to local proof, reused dependency, or named exclusion; then one unmapped invariant", expectedState: "complete matrix passes; incomplete matrix fails final architecture approval gate", expectedError: "conformance-gate failure, not a runtime Problem code", expectedSideEffects: "no runtime state or conformity claim created", gate: "standards"}
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

## Source: `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`

Approved SHA-256: `a8428d947de610f37fd47769363c7cd06292ff36004de37de826299a6ef394b2`

### FEATURE-0018 canonical traceability rows

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
| 3 | Accepted decisions and controlling handoffs | `docs/decisions/DECISION_INDEX.md`; ADH-2026-042; ADH-2026-064 (ExecutionTarget ownership correction); ADH-2026-070/071 (FEATURE-0018 semantics and executable conformance) |
| VS0-SCHEMA-020 | PrincipalRef (EmbeddedValue) | FEATURE-0018 | DEC-0050; DEC-0060 Accepted; ADH-2026-070/071 | FEATURE-0018 architecture; canonical data model | VS0-CF-F18-01..08,13..18,29..31,35..49 |
| VS0-SCHEMA-021 | Membership | FEATURE-0018 | DEC-0060 Accepted; ADH-2026-070/071 | FEATURE-0018 architecture; canonical catalog | VS0-CF-F18-01,03,04,07,08,32,36,46,47 |
| VS0-SCHEMA-022 | RoleDefinition | FEATURE-0018 | DEC-0060 Accepted; ADH-2026-070/071 | FEATURE-0018 architecture; canonical catalog | VS0-CF-F18-01,09..13,21,25,31,33,34,39,44,49 |
| VS0-SCHEMA-023 | RoleAssignment | FEATURE-0018 | DEC-0060 Accepted; ADH-2026-070/071 | FEATURE-0018 architecture; canonical catalog | VS0-CF-F18-01,03..05,08,09,13,17..21,25,28,31..35,39,44..46 |
| VS0-SCHEMA-024 | GovernanceProfile | FEATURE-0018 | DEC-0050; DEC-0060 Accepted; ADH-2026-070/071 | FEATURE-0018 architecture; canonical catalog | VS0-CF-F18-51 |
| VS0-SCHEMA-062 | AccessGroup | FEATURE-0018 | DEC-0060 Accepted; ADH-2026-070/071 | FEATURE-0018 architecture; canonical catalog | VS0-CF-F18-08,45..49 |
| VS0-SCHEMA-063 | PrivilegedAccessRequest | FEATURE-0018 | DEC-0060 Accepted; ADH-2026-070/071 | FEATURE-0018 architecture; canonical catalog | VS0-CF-F18-13..18,28,40,41,47,49 |
| VS0-SCHEMA-064 | AccessReview | FEATURE-0018 | DEC-0060 Accepted; ADH-2026-070/071 | FEATURE-0018 architecture; canonical catalog | VS0-CF-F18-18..20,31,37,42,45,47..49 |
| VS0-SCHEMA-065 | ApprovalPolicy | FEATURE-0018 | DEC-0060 Accepted; ADH-2026-070/071 | FEATURE-0018 architecture; canonical catalog | VS0-CF-F18-13..16,18,21,26,40,41,47,49,51 |
| VS0-SCHEMA-066 | ApprovalRequest | FEATURE-0018 | DEC-0060 Accepted; ADH-2026-070/071 | FEATURE-0018 architecture; canonical catalog | VS0-CF-F18-13..18,21,26,28,40,41,43,47,49 |
| VS0-SCHEMA-067 | ExceptionGrant | FEATURE-0018 | DEC-0060 Accepted; ADH-2026-070/071 | FEATURE-0018 architecture; canonical catalog | VS0-CF-F18-21..24,41,43 |
| VS0-SCHEMA-068 | FEATURE-0018 supporting values | FEATURE-0018 | DEC-0060 Accepted; ADH-2026-070/071 | FEATURE-0018 architecture | VS0-CF-F18-01..54 except retired F18-38; exact registry rows control |
| VS0-WRITER-007 | authorized-governance-publisher | FEATURE-0018+FEATURE-0019 | DEC-0041/0043/0050/0055 | VS0-CF-HP01,F07 |
| VS0-WRITER-008 | RoleAssignment-controller (sole canonical publisher/writer; callers and workflow controllers submit exact intents only) | FEATURE-0018 | DEC-0060 Accepted; ADH-2026-070/F18-RD-02/08 | VS0-CF-F01,F02,X01..X03 |
| VS0-WRITER-022 | identity-membership-controller | FEATURE-0018 | DEC-0060 Accepted; ADH-2026-070/F18-RD-02/05 | VS0-CF-F01,F02,X01..X03 |
| VS0-F01 | VS0-CF-F01 | AUTH_REQUIRED | 401 | VS0_AUTH_REQUIRED | FEATURE-0018 | DEC-0050 |
| VS0-F02 | VS0-CF-F02 | RESOURCE_NOT_FOUND | 404 | VS0_AUTHORIZATION_SAFE_DENIAL | FEATURE-0018 | DEC-0050 |
| VS0-CF-X01 | Cross-organization isolation | FEATURE-0018 | security | DEC-0050; DEC-0060 Accepted | VS0-WRITER-008, VS0-SCHEMA-023 |
| VS0-CF-X02 | Cross-project isolation | FEATURE-0018 | security | DEC-0050; DEC-0060 Accepted | VS0-WRITER-008, VS0-SCHEMA-023 |
No runtime conformance test implementation exists yet. All conformance IDs (VS0-CF-HP01, VS0-CF-F01..F20, VS0-CF-X01..X03, VS0-CF-L01, VS0-CF-Z01, VS0-CF-T01, VS0-CF-I01..I02, VS0-CF-D01, VS0-CF-F15-01..41, VS0-CF-F16-01..128, and VS0-CF-F18-01..54 except retired F18-38) are exact test contracts owned by their respective feature tasks. FEATURE-0026 provides the integration proof that exercises the cross-feature Slice 0 contracts end-to-end in the synthetic profile; FEATURE-0015, FEATURE-0016, and FEATURE-0018 each own their local conformance cases only (no migration conformance exists under DEC-0059). FEATURE-0015 owns exactly 22 logical endpoint paths and registers exactly 35 explicit Go 1.22 `http.ServeMux` method/path patterns (ADH-2026-051). FEATURE-0016 owns exactly five explicit Go 1.22 `http.ServeMux` registrations plus a pre-ServeMux transport-only method/path guard (ADH-2026-058 clause 3). ADH-2026-071 registers FEATURE-0018 proof identifiers and mappings only; it adds no route or runtime semantic.
### FEATURE-0018 local conformance mapping
| Inherited-only proof | VS0-CF-X03 remains FEATURE-0015-owned and never counts as FEATURE-0018-local acceptance | ADH-2026-071 §2 |
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

## Source: `docs/reviews/architecture-readiness/FEATURE-0018-industry-validation-matrix.md`

Approved SHA-256: `fc6ee39865534a653727045dda62e19507c115d917bbd6ea4a7ba301656634c6`

### Approved enterprise scenarios F18-SCN-01 through F18-SCN-49

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
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
<!-- END EXACT APPROVED-SOURCE EXCERPT -->
