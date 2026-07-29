# Current Architecture Baseline

Status: Approved Phase 2 start baseline.

Architecture baseline: `ARCH-2026.07-PHASE2-START`

FEATURE-0013 completion status: `ADH-2026-017` is the Approved replacement that captures the complete FEATURE-0013 architecture in one content-bound approval envelope. FEATURE-0013 is implemented and merged through PR #15 on 2026-07-29.

Historical FEATURE-0013 updates: `ADH-2026-014`, `ADH-2026-015`, and `ADH-2026-016` record the path to the consolidation. They are retained for provenance, excluded from downstream instructions, and normatively superseded by approved `ADH-2026-017`.

## Product Position

Sovrunn is a cloud-native sovereign PaaS platform for local cloud providers, MSPs, and on-premise cloud operators.

Sovrunn provides governed service catalog, organization/tenant/project governance, provider-neutral placement, plugin-based service lifecycle, decision/audit/evidence records, and AI-assisted operations.

Sovrunn Data Engine is a future managed service inside the broader Sovrunn platform, not the whole product.

## Approved Architecture Principles

- Reuse before build is mandatory and applies across Sovrunn phases.
- Provider-neutral core is mandatory.
- Adapter boundaries must exist before deep integration.
- Policy logic must go through a `PolicyEngineAdapter` boundary.
- OPA is the preferred first real policy adapter candidate.
- Cedar may be evaluated later for authorization-style decisions.
- Customer-facing APIs must not expose low-level IaaS complexity.
- Provider-facing, internal, plugin-facing, and customer-facing APIs must remain separate.
- AI may recommend and explain, but must not bypass policy, approval, or audit.

## Approved FEATURE-0012 Architecture Baseline

ADH-2026-012 approves `docs/architecture/api-resource-standard.md` as the controlling baseline for FEATURE-0012 Kiro specifications. It establishes provider-neutral resource profiles, scope/reference semantics, boundary and ownership rules, status/condition grammar, strict validation, stable errors, compatibility, conformance, migration, and reassessment requirements.

The standard is cross-phase in effect but remains draft until FEATURE-0012 implementation and final review complete.

## Approved Phase 2 Scope

Phase 2 builds only:

- model foundation,
- resource/API standards,
- decision and audit standards,
- reuse assessment standard,
- adapter boundary foundation,
- provider-neutral resource model,
- governance/security/data/cost policy context foundation,
- service runtime profile foundation,
- placement request and placement decision v0,
- plugin taxonomy foundation,
- AI-readable decision context,
- Phase 2 simulation/demo.

## Phase 2 Explicit Non-Goals

Phase 2 does not build:

- real provider provisioning,
- real Kubernetes workload provisioning,
- real PostgreSQL runtime provisioning,
- full OPA/Cedar integration,
- full Keycloak/Vault/Temporal/Argo integration,
- production workflow engine,
- global traffic execution,
- autoscaling execution,
- DR/failover execution,
- billing/chargeback engine,
- full compliance evidence engine,
- autonomous AI operations,
- full UI/portal.

## Approved Phase 3 Direction

Phase 3 builds the first executable PaaS plugin chain:

- plugin execution contract v0,
- operation controller v0,
- PostgreSQL management plane plugin v0,
- Kubernetes/local substrate plugin v0,
- PostgreSQL runtime plugin v0,
- ServiceInstance provisioning v0,
- ServiceBinding and SecretRef integration,
- end-to-end MVP demo.

Phase 3 should reuse Kubernetes APIs, Helm or an existing PostgreSQL operator where practical.

## MVP Definition

MVP-001: Governed PostgreSQL PaaS Placement and Provisioning on one substrate.

The MVP must demonstrate:

- customer service request,
- entitlement/policy evaluation,
- placement decision,
- operation creation,
- plugin-chain execution,
- service instance status,
- service binding,
- audit event,
- AI-readable explanation.

## Not Approved

The following are not approved architecture directions:

- building a custom policy engine,
- building a PostgreSQL HA/failover controller from scratch,
- putting PostgreSQL lifecycle logic inside Sovrunn core,
- exposing raw IaaS implementation details as the primary customer API,
- treating logs as audit records,
- storing raw credentials in Sovrunn resource records,
- allowing plugins to bypass policy, placement, or audit,
- allowing AI recommendations to execute without policy/approval validation.

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

## Current Execution Focus

Current active phase: Phase 2.

Completed and merged: `FEATURE-0011: Reuse Assessment Standard`.

Completed and merged: `FEATURE-0012: API, Resource Naming, Status, and Validation Standard` — checkpoint 18 passed, final human approval granted (2026-07-24), and merged through PR #14 as commit `a1b74fb` into `phase2-reuse-first-paas-fabric-foundation`.

Completed and merged: `FEATURE-0013: Decision Record and AuditEvent Standard` — implemented through PR #15 on 2026-07-29. `ADH-2026-017` remains the approved single replacement handoff for the complete FEATURE-0013 architecture.

Active next architecture target: `FEATURE-0014: Provider-Neutral Resource Model`. FEATURE-0014 must consume FEATURE-0011 reuse-assessment governance and FEATURE-0012 API/resource grammar. It must not redefine FEATURE-0013 `DecisionRecord`, `AuditEvent`, `DecisionProfile`, `EvaluationResult`, `metadata.scopeRef`, `ScopeKind`, or ServiceInstance-as-subject semantics. It must not own FEATURE-0015 `ResourcePool`/`ProviderCapability` or FEATURE-0016 adapter interfaces.

## Approved Consolidated FEATURE-0013 Architecture Baseline

`ADH-2026-017` approves `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md` as one complete, authoritative architecture package. It consolidates the common immutable provider-neutral `DecisionRecord` envelope, versioned `DecisionProfile` extension, atomic `EvaluationResult` normalization, deterministic bounded composition, `AuditEvent` linkage, structural security and projection controls, conformance cases, Matrix E risks, guardrails, and the downstream adoption contract.

The consolidated scope model preserves FEATURE-0012 semantics:

1. `ScopeKind` remains the six-value governance vocabulary: Platform, Organization, OrganizationUnit, Tenant, Project, and Provider.
2. `metadata.scopeRef` is the sole scope authority. Platform uses the canonical absent scope form; every non-Platform scope uses a canonical non-empty reference.
3. `ServiceInstance` is lifecycle-contained by its owner and is represented as a typed `subjectRef`, normally under its governing Project scope. It is not a `ScopeKind`.
4. Scope identifies the governance boundary; subject identifies the entity about which the decision or audit statement was made. Neither grants authorization.

`ADH-2026-014`, `ADH-2026-015`, and `ADH-2026-016` remain historical decision records. They are normatively superseded and must not be loaded as separate Kiro, reviewer, design, task, or implementation instructions.

FEATURE-0013 requirements, design, tasks, implementation, conformance fixtures, Matrix E evidence, and final approval were completed from approved `ADH-2026-017` and merged through PR #15. Later features consume the merged contract; they do not reload predecessor FEATURE-0013 handoffs as active instructions.

## Change Control

Any proposed change to this baseline must be classified as:

- clarification,
- extension,
- correction,
- replacement,
- new decision.

Replacement or new decision requires explicit human approval and updates to impacted docs, DEC/RFC records, feature sequence, and traceability matrix.
