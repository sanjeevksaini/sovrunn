# Current Decision Summary

This file summarizes the currently binding architecture decisions. The detailed source of truth remains `docs/decisions/DECISION_INDEX.md` and individual DEC/RFC files.

## Accepted Decisions

- Sovrunn is a sovereign cloud-native PaaS platform.
- SDE is a future managed service inside Sovrunn.
- Reuse before build is mandatory (`DEC-0026`).
- FEATURE-0012 uses the approved provider-neutral API/resource baseline in `docs/architecture/api-resource-standard.md` (`ADH-2026-012`, `ADH-2026-013`).
- FEATURE-0012 is implemented and merged through PR #14 as commit `a1b74fb` into `phase2-reuse-first-paas-fabric-foundation` (final approval 2026-07-24).
- FEATURE-0013 is implemented and merged through PR #15 (2026-07-29). Its architecture remains consolidated into one complete baseline in `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md` under approved replacement handoff `ADH-2026-017`; later Phase 2 features consume this contract and must not redefine its decision, audit, scope, or ServiceInstance-subject semantics.
- FEATURE-0014 is implemented and merged through PR #16 as commit `1ed47ac` on 2026-07-30. Its architecture remains approved under `ADH-2026-018`, as clarified by `ADH-2026-019`; it owns exactly Provider, ProviderLocation, ProviderDatacenter, DatacenterFailureDomain, and InfrastructureStack, reuses Organization and FEATURE-0012 grammar, declares FEATURE-0013 adoption `NOT_APPLICABLE`, and excludes FEATURE-0015 pools/capabilities, FEATURE-0016 adapters, and FEATURE-0053 connectivity.
- The consolidated architecture uses `DecisionRecord` as the canonical governed-conclusion envelope and preserves FEATURE-0012's six-value `ScopeKind` vocabulary: Platform, Organization, OrganizationUnit, Tenant, Project, and Provider.
- `DecisionRecord` and `AuditEvent` use `metadata.scopeRef` as their sole scope authority. `ServiceInstance` is a typed subject, normally governed by Project scope, and is not a seventh `ScopeKind`.
- `ADH-2026-014`, `ADH-2026-015`, and `ADH-2026-016` are predecessor history. Approved `ADH-2026-017` supersedes them normatively; they are excluded from downstream generation and review inputs.
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
