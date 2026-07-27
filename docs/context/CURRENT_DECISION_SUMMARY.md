# Current Decision Summary

This file summarizes the currently binding architecture decisions. The detailed source of truth remains `docs/decisions/DECISION_INDEX.md` and individual DEC/RFC files.

## Accepted Decisions

- Sovrunn is a sovereign cloud-native PaaS platform.
- SDE is a future managed service inside Sovrunn.
- Reuse before build is mandatory (`DEC-0026`).
- FEATURE-0012 uses the approved provider-neutral API/resource baseline in `docs/architecture/api-resource-standard.md` (`ADH-2026-012`, `ADH-2026-013`).
- FEATURE-0012 is implemented and merged through PR #14 as commit `a1b74fb` into `phase2-reuse-first-paas-fabric-foundation` (final approval 2026-07-24).
- FEATURE-0013 uses the approved DecisionRecord and AuditEvent standard in `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md` (`ADH-2026-014`).
- `DecisionRecord` replaces the former term `DecisionObject` as a controlled terminology correction (`ADH-2026-014`).
- `AuditEvent` scope is expanded from Organization-only to six governance scopes: Platform, Organization, OrganizationUnit, Tenant, Project, ServiceInstance (`ADH-2026-014`). **(Superseded by `ADH-2026-015`.)**
- `ScopeKind` is the single canonical shared scope vocabulary with exactly seven values: Platform, Organization, OrganizationUnit, Tenant, Project, Provider, ServiceInstance (`ADH-2026-015`). `ServiceInstance` is additive; `Provider` is retained.
- `AuditEvent` permits all seven canonical `ScopeKind` values and uses `metadata.scopeRef` as its sole canonical scope identity; no separate `AuditEvent` scope enum is permitted (`ADH-2026-015`). Existing FEATURE-0012 serialized values and Organization-scoped `AuditEvent`s remain valid. FEATURE-0013 requirements are returned to fresh independent review; design/tasks/implementation remain unauthorized.
- Phase 2 is model/decision/audit/adapter only (`DEC-0027`).
- Policy evaluation must use a policy-engine abstraction (`DEC-0028`).
- Provider-neutral resource model is required before provider integrations.
- ResourcePool is the placement boundary (`DEC-0032`).
- ProviderCapability is the compatibility boundary (`DEC-0033`).
- PlacementDecision is required before provisioning (`DEC-0034`).
- Plugin taxonomy has three planes: provider/substrate, service management, service runtime (`DEC-0029`).
- MVP is governed PostgreSQL PaaS placement and provisioning on one substrate (`DEC-0030`).

## Decisions Requiring Future Revalidation

- First real policy adapter choice: OPA expected, Cedar evaluated later for authorization.
- First PostgreSQL runtime reuse choice: CloudNativePG, Crunchy, or Helm.
- Operation engine backend: simple v0 first; Temporal/Argo later.
- Secret backend: Kubernetes Secret for local MVP; Vault/External Secrets later.
