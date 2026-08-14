# Sovrunn Context Pack

This is the compact source-of-truth context for new ChatGPT, Kiro, Cursor, and reviewer sessions.

Use with `docs/context/CURRENT_ARCHITECTURE_BASELINE.md` and `docs/context/CHATGPT_ARCHITECTURE_SESSION_PROMPT.md`.

## Product Definition

Sovrunn is a cloud-native sovereign PaaS platform for local cloud providers, MSPs, and on-premise cloud operators.

It provides governed service catalog, organization/tenant/project governance, implementation-neutral placement, plugin-based service lifecycle, decision/audit/evidence records, and AI-assisted operations.

Sovrunn Data Engine is a future managed service inside Sovrunn.

## Core Build Boundary

Sovrunn builds:

- governance,
- policy context,
- decision models,
- placement decisions,
- plugin contracts,
- operation tracking,
- audit and evidence,
- AI-readable explanation context,
- customer/provider PaaS experience,
- sovereign cloud product and participation semantics,
- implementation-neutral requirements, decisions, lifecycle, and safe projections,
- sovereignty evidence composition and plugin/adapter contracts.

Sovrunn reuses or wraps mature infrastructure such as Kubernetes, OPA, Cedar, Keycloak, Vault, OpenTelemetry, Prometheus, Grafana, Temporal/Argo, PostgreSQL operators, ingress controllers, and backup tools where appropriate.

## Current Architecture Spine

```text
Service request
  -> authenticate and authorize through adapter
  -> CloudEnrollment, ServiceEntitlement, and Quota reservation
  -> EffectiveGovernanceContext
  -> immutable ServiceRequirementSet
  -> qualified ExecutionTarget candidates
  -> sovereignty DecisionRecord
  -> placement DecisionRecord
  -> immutable ServiceDeploymentPlan
  -> accepted Operation
  -> PluginExecution (fake in Phase 2R)
  -> verified synthetic readiness
  -> customer-safe ServicePlacement
  -> per-consumer ServiceBinding with SecretRef only
  -> correlated AuditEvents and safe AI-readable explanation
```

## Current Phase

Current active phase: Phase 2R.

Architecture baseline: `ARCH-2026.08-PHASE2R-CANONICAL`.

Next feature: FEATURE-0015 Canonical Cloud Model Foundation.

## Canonical Model Authority

- Semantic model: `docs/architecture/canonical/sovrunn-finalized-data-model.md`
- Contract catalog: `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`
- Reference flow: `docs/architecture/reference-flows/postgresql-end-to-end.md`
- Phase 2R rebaseline: `docs/phase2/PHASE2R_REBASELINE.md`
- VS-000 charter: `docs/architecture/vertical-slices/VS-000-core-skeleton.md`
- Delivery plan: `docs/architecture/vertical-slices/INTEGRATED_DELIVERY_PLAN.md`

## Seven Canonical Scope Kinds

Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider.

## Key Canonical Boundaries

- CloudPlatform (product ownership) ≠ CloudProvider (infrastructure operation)
- CloudEnrollment (customer joins platform) ≠ CloudProviderParticipation (provider participates)
- ExecutionTarget (qualified realization boundary) — no mandatory ResourcePool or ProviderCapability
- EffectiveGovernanceContext (resolved governance) ≠ sovereignty assessment (separate evidence)
- ServiceTypeDefinition (reusable type) → ServiceOffering (platform product) → ServicePlan (versioned, immutable)
- ServicePlacement (safe customer projection) — no topology/credential leakage
- ServiceBinding (per-consumer, SecretRef-only, revocable)
