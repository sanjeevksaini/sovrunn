---
doc_type: feature
id: FEATURE-0014
title: Provider-Neutral Resource Model
status: implemented_and_merged
phase: 2
reuse_assessment_format_version: 1.0.0
depends_on:
  - FEATURE-0011
  - FEATURE-0012
ai_load_priority: feature
ai_summary: Implemented and merged through PR #16 as commit 1ed47ac on 2026-07-30; preserves the approved five-resource provider-neutral topology under ADH-2026-018 as clarified by ADH-2026-019.
controlling_handoff: ADH-2026-018
canonical_architecture: docs/architecture/provider-neutral-resource-model.md
kiro_slug: provider-neutral-resource-model
---

# FEATURE-0014 — Provider-Neutral Resource Model

## Purpose

FEATURE-0014 defines exactly Provider, ProviderLocation, ProviderDatacenter,
DatacenterFailureDomain, and InfrastructureStack. It reuses the existing
Organization resource and FEATURE-0012 resource grammar. It does not define
ResourcePool, ProviderCapability, connectivity, adapters, provider integration,
decision/audit behavior, or runtime execution.

Canonical architecture: `docs/architecture/provider-neutral-resource-model.md`

Controlling handoff: ADH-2026-018 (Approved)

Canonical reuse standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

## Acceptance Criteria

1. The repository-readiness gate passes before requirements generation.
2. Exactly five FEATURE-0014 resource kinds are present with no active aliases.
3. Organization and FEATURE-0012 contracts are reused without redefinition.
4. All F14-AD-001 through F14-AD-021 decisions remain closed downstream.
5. All F14-R01 through F14-R30 risks retain decision, mitigation, evidence,
   residual-target, owner, and reassessment traceability.
6. Requirements, design, and tasks follow the stage-ownership and single-owner
   overlap controls in the approved architecture.
7. Adjacent-feature and provider-native semantics remain outside FEATURE-0014.
8. Human approval is required between every Kiro stage.

## Feature-level reuse summary

Feature identity: FEATURE-0014 — Provider-Neutral Resource Model.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| Organization owner identity and FEATURE-0012 resource grammar | Extend | Reuse existing Organization, profiles, metadata, scope, references, boundaries, validation, errors, concurrency, and conformance; add only constrained FEATURE-0014 topology semantics. | Approved | ADH-2026-018; ADH-2026-012 |
| Five-resource provider-neutral topology semantics | Build | No external model supplies Sovrunn's exact owner/provider separation, immutable hierarchy, completeness, one-failure-domain stack identity, and non-connectivity inference boundary. | Approved | ADH-2026-018; RFC-0024 |
| Provider integration and discovery | Wrap | Future provider-native translation belongs behind FEATURE-0016 adapter contracts; FEATURE-0014 creates no adapter or runtime. | Deferred | ADH-2026-018; DEC-0036 |

## Capability assessment: FEATURE-0012 grammar and Organization reuse

### Identity

| Field | Value |
|---|---|
| Feature identity | FEATURE-0014 |
| Capability or decision-unit identity | FEATURE-0012 grammar and Organization reuse |
| Assessment owner | Sovrunn Architecture Owner |

### Classification

| Field | Value |
|---|---|
| Disposition | Extend |
| Decision status | Approved |

### Analysis

| Field | Value |
|---|---|
| Assessment scope | Reuse existing Organization ownership and FEATURE-0012 resource grammar for five FEATURE-0014 resources. |
| Candidate category | Existing Sovrunn governance and API/resource standards. |
| Mature candidates / applicable standards | Phase 1 Organization; FEATURE-0012; JSON Schema 2020-12; RFC 9457; RFC 6901; ETag/If-Match; selected Kubernetes-style declarative conventions. |
| Relevant candidate strengths | Existing identity, scope, reference, status, validation, error, compatibility, and conformance semantics are approved and implemented. |
| Material candidate constraints | Generic grammar does not define provider-topology hierarchy or domain invariants. |
| Rationale | Extend approved internal standards rather than create parallel metadata, scope, ownership, reference, or validation contracts. |
| Selected foundation or approach | FEATURE-0012 ManagedResource/operator-facing grammar plus existing Organization and Provider ScopeKind. |
| Why Reuse is insufficient | Domain-constrained topology relationships and completeness semantics are new. |
| Why Wrap is insufficient | No external runtime or engine is involved. |
| Why Build is insufficient | Rebuilding shared grammar would conflict with approved earlier features. |
| Protected Sovrunn differentiation and long-term ownership | FEATURE-0012 retains common-contract ownership; FEATURE-0014 owns only topology meaning. |

### Boundary

| Field | Value |
|---|---|
| Sovrunn-owned responsibility | FEATURE-0014 topology resources and invariants. |
| Reused or extended responsibility | Organization and FEATURE-0012 own governance identity and common resource grammar. |
| Responsibility/control boundary | FEATURE-0014 constrains inherited contracts but does not redefine them. |
| Data crossing the boundary | Typed references, ObjectMeta, spec/status, conditions, and Problem Details. |
| Control crossing the boundary | Local validation only; no provider execution or external control. |
| Adapter required | No |
| Adapter rationale | This assessment unit reuses internal contracts without external integration. |
| Adapter or contract identifier | none |
| Vendor-native types allowed | No |

### Suitability

| Field | Value |
|---|---|
| Sovereignty and deployment fit | Local provider-neutral contracts support sovereign and disconnected environments. |
| Security and trust | Sole scope authority, typed constrained references, redaction, and no-existence disclosure are inherited. |
| Operational and supportability | Stable errors and deterministic validation/conformance are inherited. |
| Licensing and supply-chain | No new dependency or vendor SDK is selected. |
| Portability and provider-neutrality impact | Improves portability by excluding provider-native core identity and types. |

### Phase and scope

| Field | Value |
|---|---|
| Allowed in current phase | Yes |
| Current-phase work | Contract semantics, schemas, validation, conformance, and documentation after approved Kiro stages. |
| Deferred work | ResourcePool, capabilities, adapters, connectivity, provider integration, and execution. |
| Explicit non-goals | No common-grammar redefinition; no owner resource; no adjacent-feature runtime or domain semantics. |
| Exit or migration boundary | Common-contract change returns to FEATURE-0012 ownership through approved architecture change. |
| Phase 2 non-goal acknowledgement | No real provider provisioning or integration is authorized. |

### Risk mitigation

| Field | Value |
|---|---|
| Applicable architecture risks | F14-R03, F14-R04, F14-R05, F14-R08, F14-R15, F14-R17, F14-R24, F14-R26, F14-R28, and F14-R30. |
| Residual risk | Target residual ratings remain pending implementation evidence and human semantic review. |
| Replacement risk | Medium |
| Reassessment triggers | FEATURE-0012 contract change, scope incident, provider-native leakage, duplicate owner model, or context conflict. |

### Traceability

| Field | Value |
|---|---|
| Related DEC / RFC / ADH references | ADH-2026-012; ADH-2026-018; DEC-0026; DEC-0036; RFC-0022; RFC-0024 |
| Linked acceptance criteria | Acceptance criteria 1–8 in this document; F14-AD-002 through F14-AD-004 and F14-AD-014 through F14-AD-019. |
| Validation and review evidence | FEATURE-0014 boundary validator, repository-readiness preflight, FEATURE-0012 regression/conformance evidence, and later stage review evidence. |

### Human-approval evidence

| Field | Value |
|---|---|
| Approving person or role | Sovrunn Architecture Owner |
| Approval date | 2026-07-29 |
| Approved ADH or assessment-review reference | ADH-2026-018 |
| Structured approval-evidence record | docs/reviews/reuse-assessments/FEATURE-0014-approval-evidence.md |
| Approval applies to | Architecture and requirements-generation readiness only; later Kiro and implementation gates remain mandatory. |

## Capability assessment: five-resource topology contract

### Identity

| Field | Value |
|---|---|
| Feature identity | FEATURE-0014 |
| Capability or decision-unit identity | Five-resource provider-neutral topology contract |
| Assessment owner | Sovrunn Architecture Owner |

### Classification

| Field | Value |
|---|---|
| Disposition | Build |
| Decision status | Approved |

### Analysis

| Field | Value |
|---|---|
| Assessment scope | Define the five resource meanings, immutable hierarchy, identity, completeness, deletion, and connectivity non-inference rules. |
| Candidate category | Provider inventory and infrastructure-topology domain models. |
| Mature candidates / applicable standards | Industry provider/location/datacenter/failure-domain concepts and mature geographic code standards where applicable. |
| Relevant candidate strengths | Common concepts improve interoperability and operator comprehension. |
| Material candidate constraints | External provider models are vendor-specific and conflate topology, capability, connectivity, and placement. |
| Rationale | Build a minimal normalized Sovrunn contract while excluding native types and later-feature semantics. |
| Selected foundation or approach | Five typed ManagedResources with immutable Provider scope and immediate-parent references. |
| Why Reuse is insufficient | No candidate preserves all approved Sovrunn scope, hierarchy, completeness, identity, and separation rules. |
| Why Wrap is insufficient | Wrapping provider APIs is FEATURE-0016 and would introduce integration early. |
| Why Extend is insufficient | External domain models do not own Sovrunn's provider-neutral topology contract. |
| Protected Sovrunn differentiation and long-term ownership | Sovrunn owns normalized topology semantics and stable cross-feature references. |

### Boundary

| Field | Value |
|---|---|
| Sovrunn-owned responsibility | Five resource meanings, hierarchy, identity, completeness, deletion, and non-connectivity inference. |
| Reused or extended responsibility | Generic resource grammar and mature geographic code standards. |
| Responsibility/control boundary | Provider-native inventory remains outside core and later crosses approved adapters only. |
| Data crossing the boundary | Provider-neutral topology declarations and current validation facts. |
| Control crossing the boundary | No provider calls, discovery, placement, connectivity, or execution. |
| Adapter required | No |
| Adapter rationale | FEATURE-0014 is contract-only; FEATURE-0016 owns adapters. |
| Adapter or contract identifier | none |
| Vendor-native types allowed | No |

### Suitability

| Field | Value |
|---|---|
| Sovereignty and deployment fit | Supports local, government, private, hybrid, and disconnected provider catalogues without vendor coupling. |
| Security and trust | Operator-facing topology, scope isolation, redaction, immutable identity, and no-existence disclosure. |
| Operational and supportability | Deterministic completeness, reference, deletion, and concurrency behavior with stable evidence. |
| Licensing and supply-chain | Standard-library/local contract work; no vendor dependency. |
| Portability and provider-neutrality impact | Stable topology can be consumed by later features without provider-name matching. |

### Phase and scope

| Field | Value |
|---|---|
| Allowed in current phase | Yes |
| Current-phase work | Five-resource contract, schema, validation, conformance, and documentation only. |
| Deferred work | Pools, capabilities, adapters, connectivity, discovery, provider calls, and runtime behavior. |
| Explicit non-goals | All FEATURE-0015, FEATURE-0016, FEATURE-0053, and provider-specific execution semantics. |
| Exit or migration boundary | Material semantic change requires replacement ADH and compatibility/migration evidence. |
| Phase 2 non-goal acknowledgement | No infrastructure provisioning, plugin execution, policy, placement, or production persistence. |

### Risk mitigation

| Field | Value |
|---|---|
| Applicable architecture risks | F14-R01 through F14-R30 as controlled by the canonical architecture risk register. |
| Residual risk | No final residual risk is accepted; target residuals require implementation evidence and human semantic review. |
| Replacement risk | Medium |
| Reassessment triggers | Any F14 risk trigger, new topology form, failed mitigation, adjacent-feature dependency, or stable API promotion. |

### Traceability

| Field | Value |
|---|---|
| Related DEC / RFC / ADH references | ADH-2026-018; DEC-0026; DEC-0036; RFC-0024 |
| Linked acceptance criteria | Acceptance criteria 1–8; F14-AD-001 through F14-AD-021; F14-R01 through F14-R30. |
| Validation and review evidence | FEATURE-0014 boundary validator, stage manifests, overlap/normalization ledgers, conformance fixtures, and human semantic review. |

### Human-approval evidence

| Field | Value |
|---|---|
| Approving person or role | Sovrunn Architecture Owner |
| Approval date | 2026-07-29 |
| Approved ADH or assessment-review reference | ADH-2026-018 |
| Structured approval-evidence record | docs/reviews/reuse-assessments/FEATURE-0014-approval-evidence.md |
| Approval applies to | Architecture and requirements-generation readiness only; design, tasks, implementation, and final residual risk remain separately gated. |

## Boundary decision

Decision: Approved

The closed architecture and ADH-2026-018 control requirements generation. Any
semantic conflict or unowned concept stops with `ARCHITECTURE_DECISION_REQUIRED`.
Any repository/tooling preflight failure stops with
`REPOSITORY_CONTEXT_NOT_READY`.

## Human-approval evidence

| Field | Value |
|---|---|
| Approving person or role | Sanjeev Kumar, Sovrunn Architecture Owner |
| Approval date | 2026-07-29 |
| Approved ADH or assessment-review reference | ADH-2026-018 |
| Approval applies to | FEATURE-0014 architecture application and requirements preflight only. |
