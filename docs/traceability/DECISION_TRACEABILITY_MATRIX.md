# Decision Traceability Matrix

This matrix links decisions to architecture docs, RFCs, features, and validation state.

| Decision | Status | Architecture Docs | RFCs | Features | Validation |
|---|---|---|---|---|---|
| DEC-0026 Reuse Before Build | Accepted | constitution, reuse-first architecture | RFC-0021 | FEATURE-0011 | Pending |
| DEC-0027 Phase 2 Scope | Accepted | development-phases, PHASE2_SCOPE | RFC-0021 | FEATURE-0011..0026 | Pending |
| DEC-0028 Policy Engine Abstraction | Accepted | policy-evaluation-abstraction | RFC-0025 | FEATURE-0017 | Pending |
| DEC-0029 Plugin Taxonomy | Accepted | plugin-taxonomy-and-boundaries | RFC-0027 | FEATURE-0024 | Pending |
| DEC-0030 PostgreSQL MVP | Accepted | MVP_001_GOVERNED_POSTGRESQL_PAAS | RFC-0029 | FEATURE-0027..0034 | Pending |
| DEC-0032 ResourcePool as Placement Boundary | Accepted | provider-neutral-resource-model, placement-decision-engine | RFC-0024, RFC-0026 | FEATURE-0015, FEATURE-0023 | Pending |
| DEC-0033 ProviderCapability as Compatibility Boundary | Accepted | provider-neutral-resource-model | RFC-0024 | FEATURE-0015 | Pending |
| DEC-0034 PlacementDecision Required Before Provisioning | Accepted | placement-decision-engine | RFC-0026 | FEATURE-0023 | Pending |

## Rule

A decision becomes `Validated` only when the related feature gate or phase gate confirms it through implementation, tests, or accepted documentation evidence.
