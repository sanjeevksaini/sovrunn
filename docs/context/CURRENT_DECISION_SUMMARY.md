# Current Decision Summary

This file summarizes the currently binding architecture decisions. The detailed source of truth remains `docs/decisions/DECISION_INDEX.md` and individual DEC/RFC files.

## Accepted Decisions

- Sovrunn is a sovereign cloud-native PaaS platform.
- SDE is a future managed service inside Sovrunn.
- Reuse before build is mandatory (`DEC-0026`).
- Phase 2R is model/decision/audit/adapter/simulation only; no real provider execution (`DEC-0027` extended by Phase 2R rebaseline).
- Policy evaluation must use a policy-engine abstraction (`DEC-0028`).
- Plugin taxonomy has three planes: provider/substrate, service management, service runtime (`DEC-0029`).
- MVP is governed PostgreSQL PaaS placement and provisioning on one substrate (`DEC-0030`).
- Adapter boundaries required before external integration (`DEC-0036`).

### Canonical Model Decisions (DEC-0037 through DEC-0058)

- Seven canonical scope kinds: Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider (`DEC-0037`). Generic Provider is migration input only.
- CloudEnrollment joins customer Organization to CloudPlatform; provider participation is separate (`DEC-0038`).
- Entitlement and quota are independent; both may deny (`DEC-0039`).
- Industry and domain clouds are versioned ServicePortfolios, not new core kinds (`DEC-0040`).
- Physical geography, governance scope, and sovereignty are independent dimensions (`DEC-0041`).
- Qualified ExecutionTarget is the actionable realization boundary; no mandatory ResourcePool or ProviderCapability in core (`DEC-0042`). Supersedes DEC-0032 and DEC-0033.
- Sovereignty and placement use registered FEATURE-0013 DecisionRecord profiles (`DEC-0043`).
- Published definitions are immutable by version (`DEC-0044`).
- One service may span multiple execution targets with explicit cross-target constraints (`DEC-0045`).
- ServicePlacement is a safe customer projection, not canonical topology access (`DEC-0046`).
- Personal onboarding creates Personal Organization, default Tenant, default Project (`DEC-0047`).
- New capabilities enter through data, contracts, plugins, adapters — not core conditionals (`DEC-0048`).
- ServiceOffering is CloudPlatform product referencing reusable ServiceTypeDefinition; global ServiceClass is migration input (`DEC-0049`).
- Governance controls compose into one EffectiveGovernanceContext; sovereignty remains separate (`DEC-0050`).
- ServiceBinding is per-consumer, SecretRef-only, separately revocable (`DEC-0051`).
- ServiceRelationshipDefinition models implementation-neutral cross-service semantics (`DEC-0052`).
- Platform lifecycle uses one external accepted-intent authority, Operation-first activation, independently recoverable agent (`DEC-0053`).
- CloudPlatform ownership, CloudProvider participation, one-provider-per-installation isolation are distinct (`DEC-0054`).
- Sovereignty evaluates SovrunnInstallation and every dependency (`DEC-0055`).
- Release compatibility is directed; recovery actions are explicit and pinned (`DEC-0056`).
- External maintenance has explicit notice authority, epochs, fences, mandatory requalification (`DEC-0057`).
- Alpha migration is one signed cutover, no dual authority, immutable history preserved (`DEC-0058`) — superseded by `DEC-0059`.
- Canonical bootstrap replaces alpha runtime migration: the first control-plane release exposes canonical contracts only; FEATURE-0001–0014 are retained repository assets and reuse input, not live state requiring conversion; no CanonicalMigrationPlan/Record, migration controller, or cutover state machine (`DEC-0059`, supersedes `DEC-0058`).
- FEATURE-0018 provides canonical scoped Sovrunn authorization, direct-member
  AccessGroup assignments, exact-Human approval/review/privileged eligibility,
  immutable exception evidence and a Sovrunn/native-IAM intersection boundary
  (`DEC-0060`, `ACR-2026-002`, approved 2026-08-30).

### Superseded Decisions

- DEC-0032 (ResourcePool as placement boundary) — superseded by DEC-0042.
- DEC-0033 (ProviderCapability as compatibility boundary) — superseded by DEC-0042.
- DEC-0058 (alpha migration as one signed cutover) — superseded by DEC-0059.

### Completed Feature Architecture (Retained History)

- FEATURE-0012 uses approved API/resource baseline (`ADH-2026-012`, `ADH-2026-013`). Merged PR #14.
- FEATURE-0013 uses approved consolidated DecisionRecord/AuditEvent architecture (`ADH-2026-017`). Merged PR #15.
- FEATURE-0014 uses approved provider-neutral resource model (`ADH-2026-018`, `ADH-2026-019`). Merged PR #16.

## Decisions Requiring Future Revalidation

- First real policy adapter choice: OPA expected, Cedar evaluated later for authorization.
- First PostgreSQL runtime reuse choice: CloudNativePG, Crunchy, or Helm.
- Operation engine backend: simple v0 first; Temporal/Argo later.
- Secret backend: Kubernetes Secret for local MVP; Vault/External Secrets later.

## Approved FEATURE-0018 stage boundary

DEC-0060 and its exact architecture/reuse package are accepted. Independent
security-review renewal-05 passed with no blocking findings. Requirements,
design, tasks and implementation remain separate human-gated stages.
