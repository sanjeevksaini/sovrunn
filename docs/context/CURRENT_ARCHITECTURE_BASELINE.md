# Current Architecture Baseline

Status: Approved Phase 2R canonical baseline.

Architecture baseline: `ARCH-2026.08-PHASE2R-CANONICAL`

Predecessor baseline: `ARCH-2026.07-PHASE2-START`

Controlling adoption: `ADH-2026-042`, `ACR-2026-001`, `ARCH-APPROVAL-2026-004`

## Canonical Model Authority

The canonical semantic model is `docs/architecture/canonical/sovrunn-finalized-data-model.md` version 1.8.

The canonical contract catalog is `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md` version 1.8.

The PostgreSQL reference flow is `docs/architecture/reference-flows/postgresql-end-to-end.md`.

The Phase 2R rebaseline is `docs/phase2/PHASE2R_REBASELINE.md`.

The cross-feature acceptance charter is `docs/architecture/vertical-slices/VS-000-core-skeleton.md`.

## Product Position

Sovrunn is a cloud-native sovereign PaaS platform for local cloud providers, MSPs, and on-premise cloud operators.

Sovrunn provides governed service catalog, organization/tenant/project governance, implementation-neutral placement, plugin-based service lifecycle, decision/audit/evidence records, and AI-assisted operations.

Sovrunn Data Engine is a future managed service inside the broader Sovrunn platform, not the whole product.

## Approved Architecture Principles

- Reuse before build is mandatory and applies across Sovrunn phases.
- Implementation-neutral core is mandatory; provider-native objects do not appear in customer or core APIs.
- Adapter boundaries must exist before deep integration.
- Policy logic must go through a `PolicyEngineAdapter` boundary.
- OPA is the preferred first real policy adapter candidate.
- Cedar may be evaluated later for authorization-style decisions.
- Customer-facing APIs must not expose low-level IaaS complexity.
- Provider-facing, internal, plugin-facing, and customer-facing APIs must remain separate.
- AI may recommend and explain, but must not bypass policy, approval, or audit.
- Published definitions are immutable by version (DEC-0044).
- Extensibility enters through data, contracts, plugins, and adapters — not core conditionals (DEC-0048).

## Seven Canonical Scope Kinds

Per DEC-0037, the active governance scope vocabulary is:

```text
Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider
```

The former six-scope vocabulary (Platform, Organization, OrganizationUnit, Tenant, Project, Provider) is superseded. `ServiceInstance` remains a typed subject within its governing scope, not a scope kind.

## Canonical Cloud Model

- **CloudPlatform** — product ownership boundary for a sovereign cloud offering.
- **CloudProvider** — infrastructure/operations supply boundary.
- **CloudProviderParticipation** — relationship between a CloudProvider and a CloudPlatform installation.
- **CloudEnrollment** — customer Organization joins a CloudPlatform (DEC-0038).
- **HostingLocation** — registered geographic/sovereignty fact.
- **Datacenter** — physical site boundary.
- **FaultDomain** — failure-isolation boundary inside a datacenter.
- **InfrastructureStack** — externally operated substrate stack inside a FaultDomain.
- **ExecutionTarget** — qualified actionable realization boundary (DEC-0042). No mandatory ResourcePool or ProviderCapability.

## Canonical Service Model

- **ServiceTypeDefinition** — reusable, implementation-neutral service type definition.
- **ServiceOffering** — CloudPlatform-scoped product referencing a ServiceTypeDefinition (DEC-0049).
- **ServicePlan** — versioned, immutable customer-facing plan under a ServiceOffering.
- **ServiceRuntimeProfile** — runtime/capability requirements bridge.
- **ServiceRequirementSet** — immutable requirements for placement and execution.
- **ServicePlacement** — safe, immutable customer projection (DEC-0046).
- **ServiceBinding** — per-consumer, SecretRef-only access boundary (DEC-0051).
- **ServiceRelationshipDefinition** / **ServiceRelationship** — cross-service semantics (DEC-0052).

## Governance Model

- **EffectiveGovernanceContext** — single resolved governance composition (DEC-0050).
- **SovereigntyProfile** / **SovereigntyFactSet** / **EvidenceRecord** — evidence-backed sovereignty evaluation (DEC-0055).
- **DecisionRecord** with registered profiles for sovereignty, placement, governance (DEC-0043).
- **EntitlementPackage** / **ServiceEntitlement** — what may be consumed (DEC-0039).
- **QuotaPolicy** — how much may be consumed, independently (DEC-0039).
- **FEATURE-0018 authorization foundation** — direct-principal and direct-member
  AccessGroup RoleAssignments, exact-Human workflow eligibility, bounded
  approval/JIT/review/exception evidence, and Sovrunn/native-IAM intersection
  under DEC-0060.

## Completed Feature History

FEATURE-0001 through FEATURE-0014 are completed implementation history:

- FEATURE-0011: Reuse Assessment Standard (merged)
- FEATURE-0012: API, Resource Naming, Status, and Validation Standard (merged PR #14, 2026-07-24)
- FEATURE-0013: Decision Record and AuditEvent Standard (merged PR #15, 2026-07-29; ADH-2026-017 controlling)
- FEATURE-0014: Provider-Neutral Resource Model (merged PR #16, 2026-07-30; ADH-2026-018/019 controlling)

These features used the alpha model terminology. Their implementation remains intact as retained repository assets and reuse input under FEATURE-0011; they do not constitute live control-plane data requiring conversion. Per DEC-0059 (superseding DEC-0058), the first executable control-plane release exposes canonical contracts only — there is no alpha runtime migration to execute.

## Superseded Decisions

- DEC-0032 (ResourcePool as placement boundary) — superseded by DEC-0042.
- DEC-0033 (ProviderCapability as compatibility boundary) — superseded by DEC-0042.
- DEC-0058 (alpha migration as one signed cutover) — superseded by DEC-0059 (canonical bootstrap; no runtime alpha migration).

## Phase 2R Scope

Phase 2R replaces the unimplemented Phase 2 features (FEATURE-0015 through FEATURE-0026) with the canonical model sequence. It remains a contract, simulation, and side-effect-free phase. No real provider, Kubernetes, OpenShift, or PostgreSQL calls.

## Implemented Phase 2R Features

Completed and merged: `FEATURE-0018: Governance, IAM, Approval and Exception Foundation` — PR #20 merged 2026-09-07 as commit `2295e9e` into `phase2-reuse-first-paas-fabric-foundation`. The approved architecture `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md` and handoff package `ADH-2026-070`, `ADH-2026-070-appendix-a`, `ADH-2026-070-appendix-b`, `ADH-2026-071`, `ADH-2026-072`, `ADH-2026-073`, `ADH-2026-074`, `ADH-2026-075`, `ADH-2026-076` remain controlling. It owns AccessGroup, Membership, RoleDefinition, RoleAssignment, PrivilegedAccessRequest, AccessReview, ApprovalPolicy, ApprovalRequest, ExceptionGrant, GovernanceProfile v1, FEATURE-0018 embedded and transient supporting contracts, Deterministic in-memory fixtures and local conformance. Adjacent ownership remains excluded: FEATURE-0019 (sovereignty profiles, regulatory policy bundles, sovereignty facts and evidence, placement-facing sovereignty interpretation); FEATURE-0020 (ProfileAssignment, governance inheritance and conflict resolution, exception application, EffectiveGovernanceContext construction); FEATURE-0021 (CloudEnrollment, entitlement, quota, provider-selection intent); FEATURE-0022 (service requirements, service-region eligibility, placement profile semantics); FEATURE-0023 (candidate-set policy adoption, ranking and placement, DecisionRecord publication for downstream placement); FEATURE-0024 (production adapter selection, plugin execution, provider-native IAM execution, provisioning and realization); FEATURE-0025 (customer-safe explanation, CloudProvider-safe explanation, AI-readable explanation projection); FEATURE-0026 (cross-feature orchestration, Slice 0 integration proof ownership).

## Phase 2R Next Planned Feature

FEATURE-0015: Canonical Cloud Model Foundation (renamed under ADH-2026-045; no alpha runtime migration scope; ends at InfrastructureStack; ExecutionTarget owned entirely by FEATURE-0016).

## MVP Definition

MVP-001: Governed PostgreSQL PaaS Placement and Provisioning on one substrate.

The MVP demonstrates the canonical end-to-end flow defined in `docs/architecture/reference-flows/postgresql-end-to-end.md`.

## Not Approved

The following remain not approved architecture directions:

- building a custom policy engine,
- building a PostgreSQL HA/failover controller from scratch,
- putting PostgreSQL lifecycle logic inside Sovrunn core,
- exposing raw IaaS implementation details as the primary customer API,
- treating logs as audit records,
- storing raw credentials in Sovrunn resource records,
- allowing plugins to bypass policy, placement, or audit,
- allowing AI recommendations to execute without policy/approval validation,
- mandatory ResourcePool or provider-wide ProviderCapability in core,
- dual write/authority between old and new canonical models,
- native AWS, OCI, OpenStack, Kubernetes, or OpenShift objects in customer/core schemas,
- runtime alpha data migration, a migration controller, or a cutover state machine (DEC-0059).

## Deferred Decisions

Deferred until later phases:

- full OPA integration,
- full Cedar integration,
- full Keycloak/Dex integration,
- full Vault/External Secrets integration,
- Temporal/Argo workflow backend,
- VMware/OpenStack/AWS/Azure provider integrations,
- ResilienceGroup execution,
- GlobalTrafficPolicy execution,
- autoscaling execution,
- spot/preemptible capacity execution,
- full compliance evidence engine,
- AI autonomous remediation,
- production-grade UI/portal,
- billing/chargeback.

## Change Control

Any proposed change to this baseline must be classified as:

- clarification,
- extension,
- correction,
- replacement,
- new decision.

Replacement or new decision requires explicit human approval and updates to impacted docs, DEC/RFC records, feature sequence, and traceability matrix.
