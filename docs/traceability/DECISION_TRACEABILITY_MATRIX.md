# Decision Traceability Matrix

This matrix links decisions to architecture docs, RFCs, features, and validation state.

## Pre-Phase 2R Decisions

| Decision | Status | Architecture Docs | RFCs | Features | Validation |
|---|---|---|---|---|---|
| DEC-0026 Reuse Before Build | Accepted | constitution, reuse-first architecture, PHASE2_REUSE_ASSESSMENT_STANDARD | RFC-0021 | FEATURE-0011 | Validated |
| DEC-0036 Adapter Boundaries | Accepted | adapter-boundary-model | RFC-0021 | FEATURE-0011, FEATURE-0016 | Pending |
| DEC-0027 Phase 2 Scope | Accepted (extended by Phase 2R) | development-phases, PHASE2_SCOPE, PHASE2R_REBASELINE | RFC-0021 | FEATURE-0011..0026 | Active |
| DEC-0028 Policy Engine Abstraction | Accepted | `docs/architecture/policy-evaluation-abstraction.md`; ADH-2026-067/068/069 | RFC-0025 (Approved) | FEATURE-0017 | Architecture reconciled; implementation pending |
| DEC-0029 Plugin Taxonomy | Accepted | plugin-taxonomy-and-boundaries | RFC-0027 | FEATURE-0024 | Pending |
| DEC-0030 PostgreSQL MVP | Accepted | MVP_001_GOVERNED_POSTGRESQL_PAAS | RFC-0029 | FEATURE-0027..0034 | Pending |
| DEC-0032 ResourcePool as Placement Boundary | Superseded by DEC-0042 | provider-neutral-resource-model (historical) | RFC-0024 | FEATURE-0015 (old) | Superseded |
| DEC-0033 ProviderCapability as Compatibility Boundary | Superseded by DEC-0042 | provider-neutral-resource-model (historical) | RFC-0024 | FEATURE-0015 (old) | Superseded |
| DEC-0034 PlacementDecision Required | Accepted | placement-decision-engine | RFC-0026 | FEATURE-0023 | Pending |

## Phase 2R Canonical Model Decisions (DEC-0037 through DEC-0058)

| Decision | Status | Canonical Authority | Features | Validation |
|---|---|---|---|---|
| DEC-0037 Seven Canonical Scope Kinds | Accepted | canonical data model, glossary | FEATURE-0015 | Pending |
| DEC-0038 CloudEnrollment Joins Customer to Platform | Accepted | canonical data model | FEATURE-0015, FEATURE-0021 | Pending |
| DEC-0039 Entitlement and Quota Independence | Accepted | canonical data model | FEATURE-0021 | Pending |
| DEC-0040 Industry Portfolios Not Provider Kinds | Accepted | canonical data model | FEATURE-0022 | Pending |
| DEC-0041 Geography/Governance/Sovereignty Independent | Accepted | canonical data model | FEATURE-0015, FEATURE-0019 | Pending |
| DEC-0042 ExecutionTarget No ResourcePool | Accepted | canonical data model, PHASE2R_REBASELINE | FEATURE-0015, FEATURE-0016 | Pending |
| DEC-0043 DecisionRecord Profiles | Accepted | canonical data model, FEATURE-0013 standard | FEATURE-0017, FEATURE-0019, FEATURE-0023 | Pending |
| DEC-0044 Immutable Versioned Definitions | Accepted | canonical data model | FEATURE-0022 | Pending |
| DEC-0045 Multi-Target Placement | Accepted | canonical data model | FEATURE-0023 | Pending |
| DEC-0046 ServicePlacement Safe Projection | Accepted | canonical data model | FEATURE-0023 | Pending |
| DEC-0047 Personal Organization | Accepted | canonical data model | FEATURE-0021 | Pending |
| DEC-0048 Extensibility by Data | Accepted | canonical data model | FEATURE-0022 | Pending |
| DEC-0049 ServiceOffering Replaces ServiceClass | Accepted | canonical data model, contract catalog | FEATURE-0022 | Pending |
| DEC-0050 EffectiveGovernanceContext | Accepted | canonical data model | FEATURE-0018, FEATURE-0020 | Pending |
| DEC-0051 ServiceBinding SecretRef Only | Accepted | canonical data model | FEATURE-0024 | Pending |
| DEC-0052 Service Relationship Boundary | Accepted | canonical data model | FEATURE-0022 | Pending |
| DEC-0053 Platform Lifecycle Boundary | Accepted | canonical data model | Phase 3+ | Pending |
| DEC-0054 CloudPlatform/Provider/Installation Isolation | Accepted | canonical data model | FEATURE-0015 | Pending |
| DEC-0055 Platform Sovereignty Composition | Accepted | canonical data model | FEATURE-0019 | Pending |
| DEC-0056 Release Recovery Compatibility | Accepted | canonical data model | Phase 3+ | Pending |
| DEC-0057 Infrastructure Maintenance Authority | Accepted | canonical data model | FEATURE-0016 | Pending |
| DEC-0058 Alpha Migration No Dual Authority | Accepted | canonical data model, PHASE2R_REBASELINE | FEATURE-0015, FEATURE-0026 | Pending |

## Rule

A decision becomes `Validated` only when the related feature gate or phase gate confirms it through implementation, tests, or accepted documentation evidence.
