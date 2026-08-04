# Open Questions

This file tracks architecture questions that are not yet approved decisions.

## Phase 2R Open Questions

| ID | Question | Status | Target Decision Point |
|---|---|---|---|
| OQ-0001 | Should OPA be the first real policy adapter after Phase 2R? | Open | Before Phase 3 hardening |
| OQ-0002 | Should Phase 3 PostgreSQL runtime use CloudNativePG, Crunchy, or Helm first? | Open | Before FEATURE-0031 |
| OQ-0003 | Should the first operation engine remain in-process or wrap Kubernetes Jobs? | Open | Before FEATURE-0028 |
| OQ-0004 | Which secrets backend is used after local MVP? Vault or External Secrets first? | Open | Before MVP hardening |

## Resolved by Phase 2R Adoption

The following questions are resolved by DEC-0037 through DEC-0058:

- ResourcePool vs. ExecutionTarget as placement boundary → DEC-0042: ExecutionTarget.
- ProviderCapability as compatibility boundary → DEC-0042: no mandatory ProviderCapability in core.
- Whether governance and sovereignty should be composed or separate → DEC-0050: composed governance, separate sovereignty.
- Whether ServiceClass is the canonical catalog concept → DEC-0049: ServiceTypeDefinition + ServiceOffering replaces it.
- Whether Provider is one boundary or multiple → DEC-0037, DEC-0054: CloudPlatform, CloudProvider, and installation are distinct.
- Exact ServiceClass to ServiceTypeDefinition/ServiceOffering mapping → DEC-0049 and FEATURE-0015 CanonicalMigrationPlan control this deterministically.
- Alpha migration backup/dry-run procedure → DEC-0058 and CanonicalMigrationPlan define verified backup, write freeze, deterministic mapping, no dual authority.

## Rule

Open questions do not change architecture until converted into accepted DEC/RFC records.
