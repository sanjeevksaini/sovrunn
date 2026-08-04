---
doc_type: architecture
title: FEATURE-0014 Provider-Neutral Resource Model Architecture (Historical)
status: historical-implementation
phase: 2
ai_load_priority: reference
ai_summary: FEATURE-0014 alpha-model implementation history. Active canonical terminology is implementation-neutral. See docs/architecture/canonical/sovrunn-finalized-data-model.md for the current authority.
---

> **Migration Notice (ADH-2026-042):** This document records the FEATURE-0014 implementation
> history using alpha-model terminology (Provider, ProviderLocation, ProviderDatacenter,
> DatacenterFailureDomain, InfrastructureStack). It is retained for audit. The active canonical
> model uses CloudPlatform, CloudProvider, HostingLocation, Datacenter, FaultDomain,
> InfrastructureStack, and qualified ExecutionTarget. The canonical adjective is
> "implementation-neutral", not "provider-neutral". See:
> - `docs/architecture/canonical/sovrunn-finalized-data-model.md`
> - `docs/phase2/PHASE2R_REBASELINE.md`
> - DEC-0037 through DEC-0042

# FEATURE-0014 Provider-Neutral Resource Model Architecture

**Status:** Approved for Kiro requirements
**Feature:** FEATURE-0014 — Provider-Neutral Resource Model
**Phase:** Phase 2 — Reuse-First PaaS Fabric Foundation
**Kiro stage:** Requirements authorized; requirements are intentionally not authored by architecture
**Controlling strategy:** `docs/phase2/PHASE2_EXECUTION_STRATEGY.md`

---

## 1. Purpose

FEATURE-0014 defines the minimum provider-neutral substrate topology needed by later Phase 2 features without introducing placement inventory, capability declarations, integration contracts, or execution behavior.

The feature owns exactly these resource kinds:

- `Provider`
- `ProviderLocation`
- `ProviderDatacenter`
- `DatacenterFailureDomain`
- `InfrastructureStack`

These resources describe where provider supply exists and how it is structurally partitioned. They do not state what capacity or service capability is available, select a placement target, or communicate with provider systems.

## 2. Controlling inputs

This architecture consumes, without redefining:

- the Phase 2 architecture spine and execution strategy;
- FEATURE-0012 resource profiles, type identity, metadata, Provider scope, typed references, boundary classification, status/condition grammar, validation ordering, error contract, concurrency, extension, compatibility, and conformance rules;
- the approved consolidated FEATURE-0013 architecture and sole controlling handoff `ADH-2026-017`, without semantic adoption or redefinition;
- the first attached hierarchy sketch showing Provider → ProviderLocation → Datacenter → Failure Domain → Infrastructure Stack after the approved terminology normalization; and
- the second attached hierarchy sketch showing an existing owner `Organization` above one or more `Provider` resources.

The authoritative local repository records FEATURE-0013 as implemented and merged through PR #15 on 2026-07-29. Its sole normative architecture is `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`, controlled by `ADH-2026-017`; predecessor handoffs `ADH-2026-014`, `ADH-2026-015`, and `ADH-2026-016` are provenance only and must not be loaded as downstream instructions.

### 2.1 Earlier- and adjacent-feature compatibility matrix

| Contract area | Canonical owner | Contract inherited by FEATURE-0014 | FEATURE-0014 may do | FEATURE-0014 must not do | Required compatibility evidence |
|---|---|---|---|---|---|
| Reuse governance | FEATURE-0011 | Canonical reuse-assessment structure, dispositions, approval, and reassessment rules | Classify FEATURE-0014 capabilities and record triggers | Invent a second reuse vocabulary or approval mechanism | Reuse-section structural and semantic review |
| Organization identity and governance hierarchy | Phase 1 Organization model | Existing `Organization` identity, scope, authorization boundary, and lifecycle | Reference Organization as Provider supply owner through Provider `scopeRef` | Redefine Organization, copy its fields, or create a parallel owner resource | Existing-Organization reference fixtures; same-party and distinct-party scenarios |
| Resource grammar | FEATURE-0012 | Profiles, type identity, ObjectMeta, six `ScopeKind` values, typed references, boundaries, ownership, status/conditions, validation, errors, concurrency, extensions, and conformance | Select approved profiles/boundaries and define constrained domain references and semantic validation | Add or reinterpret common metadata, scope, reference, status, condition, error, concurrency, or extension grammar | FEATURE-0012 regression suite; schema diff; common-grammar conformance |
| Provider scope | FEATURE-0012 | Provider is a formal supply scope under Platform or Organization; `metadata.scopeRef` is governance authority | Define the Provider resource and Provider-scoped topology semantics | Add a seventh ScopeKind, use `ownerRef` as scope, or create a second scope source | Platform/Organization Provider fixtures; sole-scope-authority checks |
| Decision and audit contracts | FEATURE-0013 / `ADH-2026-017` | Closed definitions of `DecisionRecord`, `DecisionProfile`, `EvaluationResult`, `AuditEvent`, scope/subject distinctions, and downstream-adoption test | Declare `NOT_APPLICABLE` and preserve the closed boundary | Produce, consume, specialize, project, audit, or redefine a governed FEATURE-0013 concept | FEATURE-0013 adoption check; forbidden-concept scan |
| Provider topology | FEATURE-0014 | Five resources, hierarchy, identity, containment, completeness, non-connectivity inference, and deletion semantics | Own only the decisions in `F14-AD-001`–`F14-AD-021` | Export topology ownership to a later feature or absorb adjacent-feature semantics | Closed-register traceability; topology conformance fixtures |
| Placement inventory and compatibility | FEATURE-0015 | Future `ResourcePool` and `ProviderCapability` contracts | Provide stable topology references for later consumption | Define pools, capabilities, capacity, eligibility, compatibility, or matching | Schema/field deny-list; FEATURE-0015 ownership review |
| Adapter and integration boundaries | FEATURE-0016 | Future adapter interfaces and provider translation boundaries | Preserve provider-neutral contracts suitable for later translation | Define adapter interfaces, provider clients, discovery, credentials, repositories, calls, or runtime integration | Interface/dependency/SDK scan; changed-file review |
| Explicit network connectivity | FEATURE-0053 roadmap placeholder | Future validated `NetworkConnectivityProfile` semantics | Preserve stable topology identity and declare connectivity Unknown when absent | Define or infer links, isolation, routes, reachability, latency, bandwidth, trust, or network health | Connectivity-field absence and non-inference fixtures |

An inherited contract remains owned by its canonical feature. Referencing or specializing it does not transfer ownership to FEATURE-0014.

### 2.2 Requirements-generation context allowlist

The repository-wide context rules in `AGENTS.md` and `docs/engineering/ai-context-loading-standard.md` remain controlling. Within that framework, FEATURE-0014 requirements generation may use the following task-specific normative inputs:

1. `AGENTS.md`
2. `docs/engineering/ai-context-loading-standard.md`
3. `docs/foundation/constitution.md`
4. `docs/decisions/DECISION_INDEX.md`
5. `docs/phase2/PHASE2_EXECUTION_STRATEGY.md`
6. `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md`
7. `docs/phase2/PHASE2_SCOPE.md`
8. `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
9. `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
10. `docs/architecture/api-resource-standard.md`
11. `docs/reviews/architecture-decision-handoffs/ADH-2026-012-feature-0012-api-resource-standard.md`
12. `docs/reviews/architecture-decision-handoffs/ADH-2026-013-operation-allowed-scopes.md`
13. `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
14. `docs/reviews/architecture-decision-handoffs/ADH-2026-017-feature-0013-consolidated-architecture.md`
15. `docs/architecture/provider-neutral-resource-model.md` after it is replaced by and content-matches this approved architecture
16. `docs/reviews/architecture-decision-handoffs/ADH-2026-018-feature-0014-provider-neutral-resource-model.md`
17. `docs/features/FEATURE-0014-provider-neutral-resource-model.md`
18. `docs/features/FEATURE_INDEX.md`
19. `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
20. `docs/context/CURRENT_DECISION_SUMMARY.md`
21. `docs/glossary.md`

Requirements generation must exclude as semantic inputs:

- FEATURE-0013 predecessor handoffs `ADH-2026-014`, `ADH-2026-015`, and `ADH-2026-016`;
- superseded or pre-migration FEATURE-0014 architecture, RFC, glossary, roadmap, feature-index, or context text;
- generated prompts, prior generated requirements/design/tasks, automation logs, review scratch files, site output, examples, and attached sketches as independent normative authority;
- FEATURE-0015, FEATURE-0016, and FEATURE-0053 design/specification artifacts beyond their approved ownership summaries;
- provider/vendor documentation as authority for Sovrunn core semantics.

If an allowed input still contains the superseded stack-kind name, an alternate location-level term, or a semantic conflict with the approved FEATURE-0014 package, requirements generation must stop with `ARCHITECTURE_DECISION_REQUIRED`; it must not reconcile the conflict locally.

### 2.3 Source-precedence rule

Global canonical ownership defined by `AGENTS.md` remains authoritative. For FEATURE-0014 task-specific semantics, apply this order:

```text
1. Approved FEATURE-0014 architecture + approved ADH-2026-018 as one package
2. Inherited FEATURE-0012 architecture + approved ADH-2026-012/013
3. FEATURE-0013 consolidated architecture + ADH-2026-017, boundary check only
4. Approved Phase 2 strategy, spine, scope, and feature sequence
5. Current FEATURE-0014 feature file, feature index, baseline, decision summary, and glossary
6. RFC, roadmap, explanatory examples, and attached sketches
```

Rules:

- A lower-precedence source may summarize but cannot broaden, narrow, rename, or override a higher-precedence contract.
- The architecture and its approved ADH are content-bound; a conflict between them is not resolved by choosing one. Generation stops pending correction.
- Earlier-feature canonical ownership still wins for inherited contracts even when FEATURE-0014 restates them.
- Roadmap placeholders, RFC summaries, examples, and diagrams are non-authoritative when they conflict with the approved package.
- Kiro must report the exact files and conflicting statements when stopping; silent precedence resolution is prohibited for semantic conflicts.

### 2.4 Single-owner overlap check

Every normative concept, field family, validation rule, status fact, reference, and behavior in requirements must have exactly one canonical owning feature.

| Concept | Single owner | FEATURE-0014 treatment | Overlap failure condition |
|---|---|---|---|
| Organization identity/hierarchy | Phase 1 | Typed reference only | FEATURE-0014 defines owner fields or lifecycle |
| Common metadata and `metadata.scopeRef` | FEATURE-0012 | Reuse unchanged | A second scope/metadata authority appears |
| `ownerRef` semantics | FEATURE-0012 | Lifecycle containment only | It is used for governance or authorization |
| Resource profiles/boundaries | FEATURE-0012 | Select existing values | A new or modified profile/boundary is defined |
| Status/conditions/errors/concurrency | FEATURE-0012 | Reuse grammar; add topology facts only | Common grammar is forked or history is stored in status |
| Decision/audit vocabulary | FEATURE-0013 | `NOT_APPLICABLE` | Any governed FEATURE-0013 concept is produced or consumed |
| Five topology resources and invariants | FEATURE-0014 | Define and own | Another feature's semantics enter these contracts |
| ResourcePool/capability/capacity/eligibility | FEATURE-0015 | Explicitly excluded | Any corresponding field, enum, status, or validation appears |
| Adapter/discovery/provider integration | FEATURE-0016 | Explicitly excluded | An interface, SDK dependency, credential, endpoint, or call appears |
| Network connectivity | FEATURE-0053 | Unknown and non-inferable | Any connection/isolation assertion or inference appears |

Requirements generation must produce an overlap ledger containing, for every normative requirement:

- requirement ID;
- owning feature;
- cited `F14-AD-*` or inherited contract;
- affected resource/field/behavior;
- confirmation that no second owner exists.

The overlap check fails when a requirement has no owner, more than one owner, cites only a lower-precedence summary, or introduces semantics owned by an adjacent feature. Failure produces `ARCHITECTURE_DECISION_REQUIRED` and stops the affected generation.

## 3. Architecture boundary

### 3.1 In scope

- Provider-supply topology and stable identity.
- Reuse of the existing `Organization` resource as the optional supply-owner and governance parent of `Provider`.
- Provider-owned location grouping.
- Provider datacenter identity and containment.
- Failure-domain identity within a datacenter.
- Provider-neutral identification of an IaaS control/compute substrate within a failure domain.
- Immutable topology relationships and provider-scope validation.
- Operator-facing declaration and safe read models using FEATURE-0012 grammar.
- Minimal current-state facts needed to distinguish registered, valid, and available-for-catalogue topology from malformed or unavailable topology.
- Explicit non-inference rules ensuring topology containment never implies network connectivity or reachability.

### 3.2 Locked non-goals

- `ResourcePool` or any placement-candidate model; owned by FEATURE-0015.
- `ProviderCapability` or capability vocabulary, discovery, validation, certification, or matching; owned by FEATURE-0015.
- Adapter interfaces, provider clients, discovery protocols, repositories, credentials, or provider-native translation; owned by FEATURE-0016 or later work.
- New or altered `DecisionRecord`, `DecisionProfile`, `EvaluationResult`, `AuditEvent`, reason-code, linkage, actor, subject, projection, or scope contracts; owned by FEATURE-0013.
- Placement requests or decisions, policy evaluation, entitlement, runtime profiles, plugins, operations, provisioning, reconciliation, failover, disaster recovery, capacity, cost, billing, or production persistence.
- Kubernetes, OpenStack, VMware, AWS, Azure, or other vendor objects as core domain types.
- A new `Owner`, `ProviderOwner`, or `SupplyOwner` resource kind; ownership reuses the existing `Organization` resource.
- Network links, routes, peering, reachability, latency, bandwidth, trust zones, connectivity health, or a `NetworkConnectivityProfile`; these require a separately owned later feature.

## 4. Canonical topology

```text
Organization (existing resource; optional supply owner)
  └── Provider
        └── ProviderLocation
              └── ProviderDatacenter
                    └── DatacenterFailureDomain
                          └── InfrastructureStack
```

`Organization` is context consumed from the existing governance model. It is not a sixth FEATURE-0014 resource.

Cardinality distinguishes registration from topology completeness:

| Parent | Child | Registered-state cardinality | Complete-topology cardinality |
|---|---|---:|---:|
| `Provider` | `ProviderLocation` | 0..* | 1..* |
| `ProviderLocation` | `ProviderDatacenter` | 0..* | 1..* |
| `ProviderDatacenter` | `DatacenterFailureDomain` | 0..* | 1..* |
| `DatacenterFailureDomain` | `InfrastructureStack` | 0..* | 1..* |

The zero-child form is intentional and supports ordered onboarding, maintenance, decommissioning, and safe deletion because parent and child resources have independent lifecycles. It means a registered parent may temporarily be empty; it does not mean the hierarchy can be skipped.

Strict containment rules:

- A Provider with no ProviderLocation is valid only as an incomplete registered shell. It cannot contain a ProviderDatacenter directly.
- A ProviderLocation with no ProviderDatacenter is valid only as an incomplete registered shell. A DatacenterFailureDomain cannot attach directly to it.
- A ProviderDatacenter with no DatacenterFailureDomain is valid only as an incomplete registered shell. An InfrastructureStack cannot attach directly to it.
- A DatacenterFailureDomain with no InfrastructureStack is valid only as an incomplete registered shell.
- A complete topology path contains every level: `Provider` → `ProviderLocation` → `ProviderDatacenter` → `DatacenterFailureDomain` → `InfrastructureStack`.
- Only a complete path may be advertised as topology-complete. This is an administrative topology fact, not FEATURE-0015 capability, capacity, placement eligibility, or service readiness.
- The exact condition name and transition mechanics are deferred to design, but requirements must preserve the distinction between registered existence and topology completeness.

The hierarchy is descriptive containment, not a customer governance hierarchy and not a placement graph.

### 4.1 Topology does not imply connectivity

Provider topology and network connectivity are orthogonal models.

Two `DatacenterFailureDomain` resources may be connected, isolated, partially connected, conditionally reachable, or not yet assessed regardless of whether they are:

- inside the same `ProviderDatacenter`;
- in different datacenters within the same `ProviderLocation`;
- in different locations belonging to the same Provider; or
- in different Providers or owner Organizations.

Therefore:

- common ancestry must not imply network reachability;
- different ancestry must not imply network isolation;
- sibling order or physical proximity must not imply latency, bandwidth, or failure correlation;
- an `InfrastructureStack` being contained by a failure domain must not imply connectivity to another stack or domain;
- absence of an explicit future connectivity assertion means `Unknown`, not connected and not disconnected;
- FEATURE-0014 must not add a `connected` boolean, adjacency list, route, endpoint, peer reference, or connectivity status to topology resources;
- later placement, resilience, or data-movement features must consume an explicitly owned and validated connectivity contract rather than infer connectivity from this hierarchy.

The roadmap currently assigns a `NetworkConnectivityProfile` foundation to FEATURE-0053. FEATURE-0014 preserves room for that later model but neither defines nor pre-designs it.

## 5. Scope and containment model

FEATURE-0012 distinguishes governance/security scope from lifecycle or topology containment.

- `Provider` is a formal provider-supply scope resource. A provider is created under either `Platform` or `Organization`, as allowed by FEATURE-0012.
- When an organization owns or governs a supply catalogue operated by another party, the `Provider` has that existing `Organization` as its immutable primary `scopeRef`.
- When the same real-world party is both supply owner and infrastructure provider, it is represented by two role-specific resources: an `Organization` and a `Provider` scoped to that organization. Equality of names or legal identity does not collapse their domain roles.
- `ProviderLocation`, `ProviderDatacenter`, `DatacenterFailureDomain`, and `InfrastructureStack` each have the containing `Provider` as their immutable primary `scopeRef`.
- Each descendant also has one immutable, typed parent reference for its immediate topology parent.
- A parent reference never replaces `scopeRef`, grants authorization, or changes tenant/organization isolation.
- Every immediate parent and child must resolve to the same Provider scope UID.
- Cross-provider topology references are invalid and must use FEATURE-0012 no-existence-disclosure behavior.
- Provider-native identifiers must not be used as Sovrunn identity or core references.
- FEATURE-0012 `ownerRef` must not represent the supply-owner organization. `ownerRef` is lifecycle containment, while the owner organization is the Provider's governance/security `scopeRef`.

Conceptual relationship fields:

| Resource | Primary scope | Immediate parent relationship |
|---|---|---|
| `Provider` | `Platform` or `Organization` | None |
| `ProviderLocation` | `Provider` | Provider is already identified by `scopeRef`; no duplicate unconstrained provider ID |
| `ProviderDatacenter` | `Provider` | constrained `providerLocationRef` |
| `DatacenterFailureDomain` | `Provider` | constrained `providerDatacenterRef` |
| `InfrastructureStack` | `Provider` | constrained `datacenterFailureDomainRef` |

Whether immediate-parent containment is represented by a domain-specific parent reference alone or additionally mirrored by FEATURE-0012 `ownerRef` is a design-level representation choice. It must not create two competing authorities for topology.

### 5.1 Supply owner and provider operator roles

The architecture distinguishes the party governing provider supply from the party operating infrastructure:

| Role | Resource | Meaning |
|---|---|---|
| Supply owner / administrative authority | Existing `Organization` | Governs which provider supply is registered and visible within its boundary. |
| Infrastructure supplier / operator | `Provider` | Identifies the party under which provider topology is declared. |

Examples:

```text
Organization: OwnerOrganization-A
  ├── Provider: ProviderOperator-A
  │     └── ProviderOperator-A topology registered for OwnerOrganization-A
  └── Provider: ProviderOperator-B
        └── ProviderOperator-B topology registered for OwnerOrganization-A

Organization: ProviderOperator-A
  └── Provider: ProviderOperator-A
        └── ProviderOperator-A's own provider topology
```

The same operator may be represented by separately scoped Provider resources when different owner organizations independently govern their catalogues. Those resources must not be inferred to be interchangeable merely because they describe the same external legal entity.

## 6. Resource semantics

### 6.1 Provider

Represents an infrastructure supplier/operator or provider-supply boundary known to Sovrunn. Its primary `scopeRef` identifies the Platform or existing `Organization` that governs the registered supply. It is not itself the owner organization, adapter configuration, credential holder, capability aggregate, customer organization, or provider plugin.

### 6.2 ProviderLocation

Represents a named geographic or operational provider grouping. `ProviderLocation` is the only permitted kind and term for this level of the FEATURE-0014 model. No alias resource, alias field, or alternative location-level vocabulary is permitted.

Location may carry normalized geographic/sovereignty descriptors appropriate to topology, such as country or jurisdiction codes, only when their semantics are provider-neutral and bounded. Policy evaluation over those descriptors belongs to later features.

### 6.3 ProviderDatacenter

Represents a provider-operated physical datacenter within one `ProviderLocation`. It is not a cluster, resource pool, or placement target.

### 6.4 DatacenterFailureDomain

Represents a provider-declared isolation boundary inside one datacenter whose members may share a common failure mode. FEATURE-0014 records identity and containment only; it does not calculate correlated risk, availability guarantees, quorum, resilience, or placement suitability.

### 6.5 InfrastructureStack

Represents one uniquely identified deployed infrastructure stack contained by exactly one failure domain. It identifies a substrate instance without exposing vendor SDK objects, endpoints, credentials, native availability-zone identifiers, capacity, or capabilities.

Illustrative technologies or environments associated with an InfrastructureStack include Apache CloudStack, Cloud Foundry, Red Hat OpenShift, AWS IaaS, and Azure IaaS. These names are examples, not resource kinds, provider-native identities, or a closed core enum. The exact bounded technology descriptor is deferred to design and must remain descriptive rather than a compatibility or capability contract.

Each InfrastructureStack has its own immutable FEATURE-0012 identity: server-generated `metadata.uid` plus its scoped resource identity. Multiple InfrastructureStack resources may use the same technology while remaining distinct stack instances. For example, two OpenShift-based stacks in two failure domains have different UIDs and are never treated as the same stack.

An InfrastructureStack must not span, belong to, or reference multiple `DatacenterFailureDomain` resources. Moving a stack to another failure domain requires recreate-and-migrate rather than changing its immutable scope or parent relationship.

An `InfrastructureStack` is not an adapter, plugin, cluster, `ResourcePool`, or `ProviderCapability`. Any provider/vendor classification must be a bounded normalized descriptor or provider-private metadata and must not drive core compatibility decisions in this feature.

## 7. FEATURE-0012 resource grammar adoption

All five kinds use FEATURE-0012 unchanged:

- canonical type identity with a domain API group and maturity version;
- singular PascalCase kinds, lowerCamelCase fields, plural kebab-case collections, and immutable lowercase kebab-case names;
- one canonical machine-readable schema per contract;
- `ManagedResource` as the proposed profile for operator-declared topology;
- `operator-facing` as the primary boundary, with separately reviewed projections if another audience is required;
- common metadata, opaque UID/resource version, ETag/If-Match, and authoritative field writers;
- client/operator-owned desired `spec` and Sovrunn-owned `status`;
- tri-state conditions for current facts only, never event history;
- typed constrained references with optional resolved UID;
- ordered structural, semantic, reference, authorization, and later-domain validation;
- RFC 9457 Problem Details, stable Sovrunn codes, RFC 6901 violation paths, redaction, and finite bounds;
- rejection of unknown and duplicate fields by default;
- registered namespaced extensions only.

No new FEATURE-0012 resource profile, scope kind, reference grammar, status grammar, error shape, or compatibility rule is introduced.

## 8. State, ownership, and validation

### 8.1 Authoritative ownership

- An authorized provider/platform operator owns topology declarations in `spec`.
- Sovrunn owns identity, resolved scope, generation, resource version, timestamps, validation status, and conditions.
- One registered producer owns each condition type.
- No adapter or external provider system is an authoritative writer in FEATURE-0014.

### 8.2 Minimal state posture

State describes current catalogue validity and administrative availability only. It must not claim discovered provider health, capacity, capabilities, placement eligibility, or infrastructure execution state.

### 8.3 Domain invariants

- Every descendant has exactly one Provider scope.
- Every non-platform-scoped Provider has exactly one existing `Organization` scope; no new owner kind is accepted.
- A Provider's supply-owner relationship is expressed only by its primary `scopeRef`, never by `ownerRef` or an unconstrained owner identifier.
- Every descendant has exactly one immediate topology parent of the required kind.
- Parent and child resolve to the same Provider scope UID.
- Parent relationships and primary scope are immutable; movement requires recreate-and-migrate.
- Topology cycles are impossible by kind constraints and must also be rejected semantically.
- Deleting or deactivating a parent must not silently reparent descendants.
- Provider-native identifiers, credentials, endpoints, and secrets are prohibited from common metadata, errors, conditions, and audit text.
- A topology resource must not declare capacity or capability fields that belong to FEATURE-0015.
- No topology relationship, shared ancestor, or resource status may be interpreted as evidence of network connectivity.

## 9. FEATURE-0013 Adoption

Applicability: `NOT_APPLICABLE`

Rationale: FEATURE-0014 defines provider-neutral topology resources only. It does not produce, consume, specialize, project, audit, or enforce `EvaluationResult`, `DecisionProfile`, `DecisionRecord`, decision rationale or obligations, decision-linked `AuditEvent`, decision-linked `Operation`, `DecisionContext`, or a decision composition strategy.

FEATURE-0014 nevertheless preserves the closed upstream boundary:

- Current resource facts use FEATURE-0012 status and conditions only.
- FEATURE-0014 introduces no decision or audit payload, profile, result, outcome, reason, obligation, actor, subject, linkage, projection, or composition semantics.
- `metadata.scopeRef` remains the sole governance-scope authority.
- The six-value FEATURE-0012/0013 `ScopeKind` vocabulary remains unchanged: `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, and `Provider`.
- `ownerRef` remains lifecycle containment only.
- If a later feature audits or decides upon topology changes, that later feature must perform its own FEATURE-0013 applicability classification and reuse the approved contract.

## 10. Reuse-before-build assessment

### Assessment scope

Provider-neutral substrate topology and its five resource contracts.

### Existing mature solutions or standards

- selected Kubernetes-style declarative resource conventions already adopted by FEATURE-0012;
- ISO country/jurisdiction identifiers where a normalized geographic field requires them;
- industry concepts of provider locations, datacenters, failure domains, and infrastructure-stack boundaries.
- the existing Phase 1 `Organization` resource for owner/governance identity.

### Decision

- **Reuse:** the existing `Organization` resource and FEATURE-0012 grammar; preserve the approved FEATURE-0013 boundary without semantic adoption.
- **Extend:** the Phase 2 provider-supply scope with five Sovrunn domain resources.
- **Build:** only the Sovrunn-owned normalized topology semantics and constrained schemas.
- **Wrap:** none; adapter work is prohibited in FEATURE-0014.

### Responsibility boundary

Sovrunn owns stable identity, normalized hierarchy, containment invariants, scope safety, boundary classification, and conformance. Providers retain ownership of their native inventory, identifiers, APIs, credentials, topology detail, and operational truth. Translation between the two is deferred to FEATURE-0016 or later approved work.

### Reassessment triggers

Reassess through an Architecture Decision Handoff if a second canonical location kind is proposed, parent mobility is unavoidable, provider scope cannot express required ownership, connectivity must become an earlier placement input, or FEATURE-0014 begins producing or consuming a governed FEATURE-0013 concept. The one-failure-domain-per-InfrastructureStack rule is locked and may change only through an approved replacement decision.

## 11. Architecture fitness checks

- Schema inventory contains exactly the five FEATURE-0014 kinds.
- No schema or document defines `ResourcePool` or `ProviderCapability` beyond explicit non-goal/dependency references.
- No adapter interface, SDK type, provider call, credential, plugin contract, provisioning behavior, or discovery workflow is introduced.
- Every descendant schema allows only Provider primary scope and constrains its immediate parent kind.
- Provider allows only the FEATURE-0012 Platform-or-Organization primary scope; no new owner resource or scope kind exists.
- Fixtures cover both a distinct-owner/operator case (`Organization: OwnerOrganization-A` → `Provider: ProviderOperator-A/ProviderOperator-B`) and a same-party case (`Organization: ProviderOperator-A` → `Provider: ProviderOperator-A`).
- Fixtures cover empty registered parents at every level and reject every skipped-level relationship.
- Fixtures prove that two stacks using the same technology retain distinct immutable UIDs and parents.
- Fixtures reject an InfrastructureStack with zero, ambiguous, or multiple failure-domain parents.
- Conformance fixtures prove that same-datacenter, cross-datacenter, and cross-location relationships carry no implied connectivity semantics.
- Cross-provider references fail without target-existence disclosure.
- API fields and examples contain no provider-native type as a core contract.
- The architecture contains the required FEATURE-0013 `NOT_APPLICABLE` declaration and introduces no governed FEATURE-0013 concept.
- FEATURE-0012 conformance fixtures and compatibility rules remain applicable without modification.

## 12. Risks and controls

Risk scoring uses `Likelihood × Impact` on a 1–5 scale. Residual ratings are architecture targets pending requirements, design, implementation evidence, and human semantic review; this draft does not accept residual risk on behalf of the architecture owner.

| Risk ID | Category | Risk scenario | Inherent | Mapped decisions | Mitigation | Required verification evidence | Target residual | Owner | Reassessment trigger |
|---|---|---|---:|---|---|---|---:|---|---|
| `F14-R01` | Scope | Topology becomes a hidden placement, capacity, or capability model. | 4×5 Critical | `F14-AD-001`, `017` | Keep only topology semantics; FEATURE-0015 owns pools, capabilities, eligibility, and matching. | Schema/field deny-list; changed-file semantic review; negative fixtures. | 1×5 Medium | FEATURE-0014 owner | Any capacity, eligibility, compatibility, or placement field is proposed. |
| `F14-R02` | Terminology | Competing location terms create duplicate kinds, fields, or client interpretations. | 4×4 High | `F14-AD-001`, `006` | Permit only `ProviderLocation`; prohibit aliases and alternate vocabulary. | Repository-wide terminology scan; schema and route inventory. | 1×3 Low | API architecture owner | Alias or alternate location-level term is requested. |
| `F14-R03` | Governance | Containment or `ownerRef` is mistaken for governance/security scope. | 4×5 Critical | `F14-AD-002`–`005`, `014` | Use `metadata.scopeRef` as sole scope authority and typed immediate-parent references for topology. | Positive/negative scope fixtures; cross-scope authorization tests. | 1×5 Medium | Security/API owner | A second scope source or owner-based authorization path appears. |
| `F14-R04` | Governance | A new owner resource duplicates the existing Organization model. | 3×4 High | `F14-AD-002`, `003` | Reuse existing Organization; allow Platform-scoped Provider where appropriate. | Kind inventory; owner-scenario conformance fixtures. | 1×3 Low | Architecture owner | Existing Organization cannot express a validated ownership case. |
| `F14-R05` | Identity | Same real-world party causes Organization and Provider roles to be collapsed. | 3×4 High | `F14-AD-002`, `003` | Preserve role-specific resources even when names or legal identities match. | Same-party and distinct-party fixtures; identity-resolution tests. | 1×3 Low | Domain owner | A workflow assumes Organization UID equals Provider UID. |
| `F14-R06` | Connectivity | Physical hierarchy is interpreted as network reachability or isolation. | 5×5 Critical | `F14-AD-012`, `013` | Make connectivity non-inferable; missing information is Unknown; FEATURE-0053 owns explicit connectivity. | Same- and cross-parent non-inference fixtures; field/schema absence checks. | 1×5 Medium | Network/placement architecture owner | Placement or resilience consumes ancestry as connectivity evidence. |
| `F14-R07` | Data quality | A boolean or static adjacency field creates false connectivity certainty. | 4×5 Critical | `F14-AD-012`, `013`, `017` | Prohibit connectivity fields, graphs, routes, peers, and quality/health assertions. | Schema deny-list and semantic review. | 1×4 Low | Network architecture owner | Early connectivity becomes a validated Phase 2 dependency. |
| `F14-R08` | Provider neutrality | InfrastructureStack leaks vendor SDK types, native IDs, endpoints, or credentials. | 4×5 Critical | `F14-AD-011`, `018`, `019` | Allow bounded descriptive technology only; keep native values behind later boundaries; prohibit secrets. | Provider-neutral schema scan; secret/redaction tests; dependency review. | 1×5 Medium | Security/integration owner | A native value is required for core identity or customer-facing output. |
| `F14-R09` | Lifecycle | Empty registered parents are treated as usable or topology-complete. | 4×4 High | `F14-AD-007`, `015` | Separate registered existence from topology completeness; require an unbroken five-level path. | Empty-parent fixtures at each level; completeness transition tests. | 1×4 Low | Domain owner | Incomplete topology becomes an input to later placement. |
| `F14-R10` | Identity | Shared technology is mistaken for shared InfrastructureStack identity. | 3×4 High | `F14-AD-009`, `011`, `014` | Use immutable UID and scoped identity for every deployed stack instance. | Duplicate-technology/distinct-UID fixtures; name/UID mismatch tests. | 1×3 Low | API owner | Technology is used as a lookup key or reference identity. |
| `F14-R11` | Integrity | One InfrastructureStack is attached to zero, ambiguous, or multiple failure domains. | 4×5 Critical | `F14-AD-008`, `010` | Require exactly one immutable typed failure-domain parent; cross-domain deployments use separate stacks. | Zero/multiple/wrong-kind/cross-scope parent fixtures. | 1×5 Medium | Domain owner | A real deployment claims it cannot be represented as distinct stacks. |
| `F14-R12` | Integrity | A child skips, reorders, or bypasses a hierarchy level. | 4×5 Critical | `F14-AD-005`, `007`, `008` | Constrain each child to exactly one required immediate parent and same Provider scope. | All skipped-level and wrong-parent negative fixtures. | 1×5 Medium | Domain/API owner | A new topology form requires alternate containment. |
| `F14-R13` | Lifecycle | Parent deletion creates orphans, races, cascades, or silent reparenting. | 4×5 Critical | `F14-AD-008`, `020` | Reject deletion while children exist; require leaf-first decommissioning; use FEATURE-0012 concurrency controls. | Child-existence conflict, concurrent delete/create, stale-version, and retry tests. | 2×4 Medium | Lifecycle/API owner | Cascade or bulk decommissioning becomes an approved requirement. |
| `F14-R14` | Freshness | Manually registered topology diverges from provider operational truth. | 4×4 High | `F14-AD-015`, `018` | Status states only Sovrunn-known validation/current facts; do not claim discovery or provider health. | Staleness/unknown-state fixtures; wording and writer-ownership review. | 2×3 Medium | Operator/domain owner | Discovery or external observation is introduced. |
| `F14-R15` | Security | Detailed datacenter and failure-domain topology is disclosed to unauthorized audiences. | 4×5 Critical | `F14-AD-003`, `004`, `014`, `015` | Use operator-facing boundary, least-privilege projections, redaction, and no-existence disclosure. | Cross-scope list/get/reference tests; response/error redaction review. | 1×5 Medium | Security owner | Customer-facing or AI projection of topology is proposed. |
| `F14-R16` | Sovereignty | Geographic descriptors are inaccurate, ambiguous, or treated as proof of authoritative assignment, residency, or compliance. | 4×5 Critical | `F14-AD-011`, `014`, `017` | Validate syntax, bounds, and country-prefix agreement only; treat geography as declared topology and make no assignment, residency, compliance, or policy claim. FEATURE-0014 owns no geographic reference dataset or dependency. | Malformed/prefix-mismatch fixtures; syntactically valid unassigned-code fixture; semantic review proving no authoritative or residency inference. | 2×4 Medium | Sovereignty/domain owner | Authoritative geographic membership becomes a policy, placement, or attestation input. |
| `F14-R17` | Scale | Unbounded hierarchy size, metadata, or traversal causes latency or memory exhaustion. | 4×4 High | `F14-AD-014`, `015` | Inherit FEATURE-0012 finite bounds and pagination; avoid recursive unbounded payloads. | Boundary/property tests; pagination, maximum-size, and traversal benchmarks. | 2×3 Medium | API/performance owner | Expected provider topology exceeds reviewed limits. |
| `F14-R18` | Concurrency | Concurrent onboarding, reparenting attempts, or decommissioning produces stale completeness or lost updates. | 4×4 High | `F14-AD-007`, `008`, `014`, `020` | Use immutable parents, resourceVersion/ETag/If-Match, deterministic recomputation, and conflict errors. | Concurrency, stale-write, generation, and condition-consistency tests. | 2×3 Medium | API/lifecycle owner | Multi-writer or external synchronization is introduced. |
| `F14-R19` | Compatibility | The InfrastructureStack rename is only partially applied, creating dual active contracts. | 4×4 High | `F14-AD-001`, `021` | Perform one atomic terminology migration; authorize no alias or dual-name period. | Repository-wide old-name scan; schema/API diff; docs/spec consistency check. | 1×4 Low | Architecture/API owner | A deployed external consumer requires compatibility migration. |
| `F14-R20` | Semantics | Descriptive stack technology silently becomes a capability or compatibility enum. | 4×4 High | `F14-AD-011`, `017` | Keep technology descriptive and bounded; FEATURE-0015 owns capability semantics. | Enum/decision-use scan; negative compatibility tests; semantic review. | 1×4 Low | Domain/placement owner | Matching or filtering begins using technology as capability evidence. |
| `F14-R21` | Audit/governance | Resource status becomes history, or FEATURE-0013 semantics are introduced without adoption. | 3×4 High | `F14-AD-015`, `016` | Conditions represent current facts only; retain FEATURE-0013 `NOT_APPLICABLE`; escalate later audit behavior. | Status-history field scan; FEATURE-0013 adoption gate. | 1×3 Low | Architecture/audit owner | FEATURE-0014 produces or consumes a governed FEATURE-0013 concept. |
| `F14-R22` | Dependency order | Adapters, discovery, provider calls, plugins, or runtime execution are designed early. | 4×5 Critical | `F14-AD-017`, `018` | Prohibit integration/runtime artifacts; preserve FEATURE-0015/0016 ownership. | Changed-file/dependency/interface scan; architecture drift review. | 1×5 Medium | Architecture owner | An implementation task adds an external-system boundary. |
| `F14-R23` | Referential integrity | Future connectivity or later-feature records reference deleted or recreated topology UIDs. | 3×4 High | `F14-AD-008`, `009`, `020` | Preserve immutable/non-reused UIDs, reject parent deletion with children, and require later consumers to handle stale typed references. | UID non-reuse and stale-reference fixtures; later-feature contract review. | 1×4 Low | API/later-feature owner | FEATURE-0015 or FEATURE-0053 introduces durable references. |
| `F14-R24` | Multi-owner isolation | The same external operator represented under multiple owner organizations is accidentally treated as one authorization domain. | 4×5 Critical | `F14-AD-003`, `004`, `009`, `019` | Treat each scoped Provider and descendant identity independently; never authorize by external operator name/native ID. | Same-operator/different-owner isolation fixtures; list/get/reference denial tests. | 1×5 Medium | Security/domain owner | Cross-owner federation or shared-provider administration is proposed. |
| `F14-R25` | Architecture governance | An unresolved semantic choice is silently resolved in requirements, design, or tasks. | 4×5 Critical | `F14-AD-001`–`021` | Closed register, stage-ownership matrix, and whole-stage `ARCHITECTURE_DECISION_REQUIRED` stop. | Ambiguity fixtures against prompts; stage review proving no unresolved semantic choice or placeholder. | 1×5 Medium | Architecture owner | Any downstream artifact selects a value not closed or explicitly delegated here. |
| `F14-R26` | Context integrity | Stale architecture, RFC, glossary, roadmap, index, prompt, or prior spec conflicts with the approved package. | 5×5 Critical | `F14-AD-001`, `006`, `016`, `017`, `018`, `021` | Atomic repository migration, task-specific allowlist, content-bound package, source precedence, and conflict-stop rule. | Repository terminology/ownership scan; context-manifest inspection; architecture/ADH content verification. | 1×5 Medium | Architecture/tooling owner | Any allowlisted input conflicts or a disallowed input enters generation/review context. |
| `F14-R27` | Specification quality | Requirements duplicate the same obligation across resources, risks, stories, or acceptance sections and drift semantically. | 4×4 High | `F14-AD-001`, `005`, `007`, `014` | Canonical requirement keys, one normative home per obligation, cross-references instead of copies, and duplicate semantic review. | Requirement-normalization ledger; duplicate-key/text/similarity report; human semantic review. | 1×4 Low | Requirements owner | Two requirements have the same owner/resource/behavior/outcome or conflicting acceptance criteria. |
| `F14-R28` | Traceability | A decision, risk, inherited contract, requirement, design element, task, or evidence obligation is omitted or orphaned. | 4×5 Critical | `F14-AD-001`–`021` | Bidirectional exact-ID traceability, single-owner overlap ledger, deterministic completeness gate, and no range-only assertions. | Automated enumeration of all decisions/risks; orphan and fan-in/fan-out report; reviewer manifest. | 1×5 Medium | Traceability/tooling owner | Missing ID, uncited normative statement, orphan design/task, or unowned evidence appears. |
| `F14-R29` | Stage separation | Kiro introduces implementation/design choices in requirements or semantic choices in design/tasks. | 4×5 Critical | `F14-AD-008`, `009`, `011`, `014`, `015`, `020` | Hard stage state machine, per-stage file/output allowlists, delegated-choice matrix, and forbidden-choice scans. | Stage diff; design-language scan in requirements; semantic-delta review across stages. | 1×5 Medium | Kiro workflow owner | A stage modifies an earlier artifact or contains a choice owned by another stage. |
| `F14-R30` | Dependency order | Adjacent-feature semantics leak through generic prompts, reviewers, examples, placeholders, or “future-proofing.” | 5×5 Critical | `F14-AD-012`, `013`, `016`, `017`, `018` | Applicability-aware prompts/reviewers, adjacent-owner deny-list, exact context selection, and feature-specific boundary gate. | Prompt/reviewer tests; forbidden-concept semantic scan; changed-file/dependency/interface review. | 1×5 Medium | Architecture/tooling owner | ResourcePool, capability, connectivity, adapter, decision/audit, or runtime semantics enter active FEATURE-0014 artifacts. |

### 12.1 Risk governance rules

- Every risk must map to at least one closed `F14-AD-*` decision and at least one verifiable evidence item.
- Requirements must preserve the risk scenario, mitigation outcome, and evidence obligation without changing the target residual rating.
- Design must identify which component/schema/check realizes each mitigation and how evidence will be produced.
- Tasks may implement only approved mitigation and evidence work; they may not accept, downgrade, close, or reclassify risk.
- Residual-risk acceptance is human-owned and occurs only after implementation evidence and semantic review.
- A new risk, increased inherent exposure, failed mitigation, or unmet target residual triggers architecture review and, when semantics change, `ARCHITECTURE_DECISION_REQUIRED`.
- Later features inherit only risks applicable to their behavior and must not silently assume FEATURE-0014 accepted them.

## 13. Closed architecture decision register

This register is the closed semantic boundary for downstream requirements, design, and tasks. Downstream stages must trace to these decisions; they may not reinterpret, weaken, extend, or replace them.

| Decision ID | Locked architecture decision | Downstream obligation |
|---|---|---|
| `F14-AD-001` | FEATURE-0014 owns exactly `Provider`, `ProviderLocation`, `ProviderDatacenter`, `DatacenterFailureDomain`, and `InfrastructureStack`. | No additional kind or compatibility alias. |
| `F14-AD-002` | Existing `Organization` is reused as the optional supply owner; it is not redefined by FEATURE-0014. | No new owner resource or scope kind. |
| `F14-AD-003` | Provider primary scope is exactly Platform or Organization; Provider `scopeRef` expresses supply governance. | `ownerRef` and unconstrained owner IDs cannot replace scope. |
| `F14-AD-004` | Every topology descendant is immutably Provider-scoped. | Cross-provider parent references are rejected without existence disclosure. |
| `F14-AD-005` | The only hierarchy is Provider → ProviderLocation → ProviderDatacenter → DatacenterFailureDomain → InfrastructureStack. | No skipped, optional, reordered, or alternate parent level. |
| `F14-AD-006` | `ProviderLocation` is the only location-level kind and term. | No alias kind, alias field, or alternate vocabulary. |
| `F14-AD-007` | Registered parents may have zero children only as topology-incomplete shells; a complete path has all five levels. | Requirements must distinguish registered existence from topology completeness. |
| `F14-AD-008` | Each child has exactly one immutable immediate parent of the required kind. | Reparenting uses recreate-and-migrate; ambiguous or multiple parents fail validation. |
| `F14-AD-009` | Each InfrastructureStack has a distinct immutable FEATURE-0012 UID and scoped identity. | Shared technology never implies shared stack identity. |
| `F14-AD-010` | Each InfrastructureStack belongs to exactly one DatacenterFailureDomain and cannot span failure domains. | Cross-domain deployments use distinct stack resources and UIDs. |
| `F14-AD-011` | Stack technology and optional ProviderLocation geography are descriptive only. Geography is syntax-normalized with country-prefix agreement but has no authoritative membership check; technology examples are not kinds, native identities, capability declarations, or a closed compatibility enum. | Neither descriptor can drive compatibility, placement, residency, compliance, or authoritative assignment in FEATURE-0014; no geographic dataset, library, lookup, or refresh lifecycle is owned here. |
| `F14-AD-012` | Physical containment implies neither network connectivity nor isolation. | No inference from ancestry, proximity, or sibling relationships. |
| `F14-AD-013` | Missing connectivity information means Unknown. | No connectivity boolean, graph, route, peer, reachability, quality, trust, or health contract. |
| `F14-AD-014` | All five resources reuse FEATURE-0012 grammar unchanged. | No new profile, scope vocabulary, metadata, reference, condition, validation, error, concurrency, or extension grammar. |
| `F14-AD-015` | Operator-declared topology uses ManagedResource and operator-facing boundaries; current facts use status/conditions. | Status is system-owned and cannot become history, capability, capacity, or audit. |
| `F14-AD-016` | FEATURE-0013 adoption is `NOT_APPLICABLE`. | No governed FEATURE-0013 concept may appear without a replacement handoff and new applicability decision. |
| `F14-AD-017` | ResourcePool and ProviderCapability are owned by FEATURE-0015. | No placement candidate, capacity, capability, eligibility, or compatibility semantics. |
| `F14-AD-018` | Adapter interfaces and provider integration are owned by FEATURE-0016 or later work. | No discovery, SDK, credential, endpoint, provider call, repository, plugin, or execution contract. |
| `F14-AD-019` | Provider-native identifiers and objects never become Sovrunn core identity or references. | Native values remain private/opaque behind later approved boundaries. |
| `F14-AD-020` | Parent deletion is rejected while children exist; FEATURE-0014 performs no cascade deletion or silent reparenting. | Decommissioning proceeds leaf-first unless a future approved lifecycle decision replaces this rule. |
| `F14-AD-021` | `InfrastructureStack` atomically replaces the superseded stack-kind name in active contracts. | No dual-name compatibility period or active alias is authorized. |

Any change to an `F14-AD-*` decision requires an approved replacement Architecture Decision Handoff before a downstream artifact changes.

## 14. Stage-ownership matrix

| Concern | Architecture ownership | Requirements may | Design may | Tasks may | Prohibited downstream action |
|---|---|---|---|---|---|
| Resource kinds and terminology | `F14-AD-001`, `006`, `021` | State testable conformance and migration outcomes | Map approved kinds to schemas/types/routes | Implement approved mappings | Add, rename, alias, or retain a superseded kind |
| Owner, scope, and containment | `F14-AD-002`–`005`, `008` | Define positive/negative authorization and reference outcomes | Choose conforming reference representations using FEATURE-0012 | Implement validation and fixtures | Change scope authority, parent kinds, or hierarchy |
| Registration and topology completeness | `F14-AD-007` | Define observable completeness semantics and acceptance cases | Select condition representation consistent with FEATURE-0012 | Implement approved conditions/checks | Treat incomplete topology as complete or placement-ready |
| InfrastructureStack identity and technology | `F14-AD-009`–`011` | Define uniqueness, immutability, and bounded descriptive-input requirements | Choose bounded field representation and reusable code-list mechanics | Implement approved schema/validation | Turn technology into identity, capability, or compatibility logic |
| Connectivity | `F14-AD-012`, `013` | Define negative non-inference requirements | Provide no connectivity model | Add only tests proving absence/non-inference | Introduce or infer any connectivity semantic |
| Shared resource grammar | `F14-AD-014`, `015` | Reference FEATURE-0012 behavior and make domain semantics testable | Bind resources to existing grammar and choose implementation structure | Implement approved bindings | Fork or extend FEATURE-0012 grammar locally |
| Decision/audit boundary | `F14-AD-016` | Preserve `NOT_APPLICABLE` | Add no adoption machinery | Add no decision/audit behavior | Produce or consume FEATURE-0013 governed concepts |
| Later-feature boundaries | `F14-AD-017`, `018` | State explicit exclusions and drift tests | Add no placeholder domain/interface | Add no implementation scaffold | Pre-design or implement FEATURE-0015/0016 |
| Native-provider leakage | `F14-AD-019` | Define redaction and rejection outcomes | Choose opaque/private later-boundary representation only if required by an approved requirement | Implement only approved boundary-safe validation | Expose native IDs, SDK objects, endpoints, or credentials in core |
| Deletion behavior | `F14-AD-020` | Require child-existence conflict and leaf-first removal | Choose stable error binding consistent with FEATURE-0012 | Implement approved validation only | Cascade delete or silently reparent descendants |
| Exact API group and routes | FEATURE-0012 constraints plus this feature boundary | Require one domain-grouped alpha API and no unversioned endpoint | Select the exact conforming group and routes | Implement approved design | Create public compatibility semantics not required by architecture |
| Initial domain field inventory | Closed decisions in section 13 | Include only fields necessary to express locked semantics | Select schema composition, indexes, and internal layout | Implement approved design | Add new domain meaning through a convenient field |
| Finite limits and configuration | FEATURE-0012 boundedness | Require all collections/strings/payloads to be finite | Select reviewed initial values and configuration mechanism | Implement approved limits | Leave values unbounded or silently change semantics |

### 14.1 Mandatory escalation rule

If requirements, design, or tasks encounter a semantic choice not resolved by section 13 or explicitly delegated by the stage-ownership matrix, that stage must:

1. emit `ARCHITECTURE_DECISION_REQUIRED` with the exact unresolved question and affected `F14-AD-*` decisions;
2. stop generation for the affected portion;
3. avoid selecting a default, inventing a placeholder, or hiding the choice in examples, schemas, tasks, or implementation notes; and
4. resume only after an approved handoff updates this architecture.

### 14.2 Traceability rule

- Every normative requirement must cite at least one `F14-AD-*` decision or an explicitly inherited FEATURE-0012 contract.
- Design must map every requirement to a conforming representation and may not add uncited normative behavior.
- Every task must cite approved design elements and their originating requirements; tasks may not cite architecture as permission to bypass design approval.
- Any normative requirement, design decision, or task without upstream traceability fails the feature gate.
- Downstream artifacts may quote concise identifiers and outcomes but must not copy this register into a competing source of truth.

## 15. Explicitly delegated representation choices

The following are not open domain semantics. They are bounded representation choices delegated by section 14:

- exact API group and versioned routes, within FEATURE-0012 domain-group and maturity rules;
- the minimal field inventory needed to express the closed decisions;
- exact condition names representing registered validity and topology completeness, using FEATURE-0012 condition grammar;
- normalized geographic code representation using a mature standard where applicable;
- bounded descriptive infrastructure-technology representation that cannot become identity, capability, or compatibility logic;
- initial finite limits and configuration values;
- schema composition, indexes, package layout, and internal implementation structure.

No delegated choice may change a resource's meaning, ownership, cardinality, scope, hierarchy, identity, connectivity posture, deletion semantics, or feature boundary.

## 16. Kiro generation guardrail contract

This section controls Kiro generation for FEATURE-0014. It is normative. Kiro must not optimize, broaden, simplify, reinterpret, or complete the architecture beyond the explicit delegations in sections 14 and 15.

### 16.1 Stage state machine

```text
APPROVED_ARCHITECTURE
  -> GENERATE_REQUIREMENTS_ONLY
  -> HUMAN_REQUIREMENTS_APPROVAL
  -> GENERATE_DESIGN_ONLY
  -> HUMAN_DESIGN_APPROVAL
  -> GENERATE_TASKS_ONLY
  -> HUMAN_TASKS_APPROVAL
  -> IMPLEMENTATION_SEPARATELY_AUTHORIZED
```

Rules:

- Only one stage may be generated per authorization.
- Approval of this architecture and ADH authorizes `requirements.md` only.
- Design generation requires explicit human `APPROVED_FOR_DESIGN` after requirements review.
- Task generation requires explicit human `APPROVED_FOR_TASKS` after design review.
- No Kiro stage authorizes source code, schemas outside the approved design stage, generated clients, migrations, tests, or implementation.
- Kiro must not generate a later stage speculatively, in the same response, or as an appendix.
- A later stage must not modify an approved earlier-stage artifact. Required changes return to that artifact's approval gate.
- `ARCHITECTURE_DECISION_REQUIRED` stops the entire current generation stage, not only the affected paragraph.

### 16.2 Per-stage output allowlist

| Authorized stage | Permitted output | Prohibited output |
|---|---|---|
| Requirements | `.kiro/specs/provider-neutral-resource-model/requirements.md` | `design.md`, `tasks.md`, schemas, code, tests, migrations, generated prompts, or implementation plans |
| Design | `.kiro/specs/provider-neutral-resource-model/design.md` | Requirements edits, `tasks.md`, code, tests, migrations, or implementation artifacts |
| Tasks | `.kiro/specs/provider-neutral-resource-model/tasks.md` | Requirements/design edits, code, schemas, tests, migrations, or execution |

Architecture, ADH, baseline, RFC, index, glossary, roadmap, and traceability updates needed to apply this handoff are repository-update work performed before Kiro stage generation; Kiro must not silently modify them while generating a stage.

### 16.3 Requirements guardrails

Requirements must:

- translate every `F14-AD-001`–`F14-AD-021` decision into one or more testable outcomes;
- enumerate `F14-R01`–`F14-R30` individually and preserve mitigation/evidence obligations;
- contain the single-owner overlap ledger from section 2.4;
- distinguish inherited FEATURE-0012 behavior from FEATURE-0014 domain requirements;
- include positive, negative, boundary, isolation, concurrency, deletion, terminology, and architecture-drift acceptance criteria;
- state all locked non-goals as enforceable exclusions;
- use `SHALL`/`MUST` only when traced to a closed decision or inherited contract;
- contain no placeholder such as `TBD`, `TODO`, “implementation-defined,” “as appropriate,” or “for example” where normative meaning is required.

Requirements must not:

- choose packages, storage engines, controllers, provider SDKs, persistence, deployment topology, internal package layout, algorithms, or source files;
- introduce a field, enum, kind, relationship, status meaning, operation, interface, or behavior not required by a cited decision;
- convert examples into normative enum values or compatibility logic;
- resolve an ambiguity through a seemingly harmless acceptance criterion.

### 16.4 Design guardrails

Design must:

- map every approved requirement ID to its representation, validator, ownership, error behavior, and evidence mechanism;
- include a field-ownership and mutability ledger with no field lacking one authoritative writer;
- map every applicable risk control to a component, schema, validation path, fixture, or review check;
- choose only representation details explicitly delegated by sections 14 and 15;
- record absence of ResourcePool, ProviderCapability, connectivity, adapter, provider-integration, and FEATURE-0013 adoption machinery;
- preserve offline structural/semantic validation separately from stateful reference/authorization validation.

Design must not:

- create a new normative requirement or weaken an acceptance criterion;
- change an `F14-AD-*` decision, risk rating, feature owner, hierarchy, cardinality, scope, identity, or deletion rule;
- add “future-proof” fields, generic extension points, placeholder interfaces, provider handles, or speculative abstractions;
- select a default for an unresolved semantic choice.

### 16.5 Task guardrails

Tasks must:

- cite the approved design section and originating requirement IDs;
- remain within FEATURE-0014 changed-file and dependency boundaries established by design;
- separate implementation, positive tests, negative tests, compatibility checks, security checks, and documentation/traceability verification;
- include explicit tasks for all architecture fitness checks and risk evidence;
- leave no architecture question for the implementer to resolve.

Tasks must not:

- restate or reinterpret architecture as task-level decisions;
- introduce optional enhancements, future hooks, cleanup unrelated to FEATURE-0014, or adjacent-feature scaffolding;
- authorize implementation before human task approval;
- combine FEATURE-0015, FEATURE-0016, or FEATURE-0053 work with FEATURE-0014.

### 16.6 Deterministic completeness checks

Before returning any stage, Kiro must prove:

- exactly five owned resource kinds and no active alias kind;
- all 21 closed decisions are individually traced;
- all 30 risks are individually traced;
- every normative statement has exactly one owner;
- every requirement maps upstream, every design element maps to requirements, and every task maps to design plus requirements;
- no prohibited term, excluded feature concept, provider-native type, connectivity inference, or uncited normative statement exists;
- no earlier or later stage file was created or modified;
- no unresolved marker or semantic ambiguity remains.

A failed completeness check stops the stage and reports the exact IDs and files involved. Kiro must not report partial approval readiness.

### 16.7 Example and diagram rule

Examples and diagrams explain locked decisions only. They cannot:

- create resource kinds, fields, enum values, default behavior, or compatibility semantics;
- override prose, the closed decision register, or the stage-ownership matrix;
- be treated as exhaustive valid data;
- resolve an ambiguity not already settled by architecture.

When an example conflicts with normative text, Kiro stops and reports `ARCHITECTURE_DECISION_REQUIRED` rather than choosing either interpretation.

### 16.8 Blocking repository-readiness gate

Requirements generation is prohibited until all checks below pass against the live branch:

1. The approved architecture is installed at `docs/architecture/provider-neutral-resource-model.md` and content-matches the human-approved package.
2. `ADH-2026-018` is installed, human-approved, and passes `make arch-handoff-check HANDOFF=<path>` or the repository-equivalent handoff validator.
3. The former stack-kind name and alternate location-level terminology are absent from active architecture, RFC, glossary, feature index, phase sequence, roadmap, current context, Kiro specs, schemas, and generated-contract sources.
4. The current FEATURE-0014 architecture/RFC no longer includes ResourcePool, ProviderCapability, capability status, connectivity, adapter, or runtime semantics except explicit ownership/non-goal references.
5. FEATURE-0014 automation metadata/state and the canonical Kiro spec path resolve successfully through the feature factory; generation must not guess paths.
6. The generated Kiro prompt context manifest contains the exact approved FEATURE-0014 architecture and ADH and no superseded/disallowed semantic input.
7. Reviewer context contains the exact approved architecture and ADH. Filename globbing alone is insufficient; the canonical architecture path must be explicitly resolved.
8. Requirements, design, tasks, and reviewer prompts are applicability-aware: FEATURE-0014 must not be required to invent an explainable decision, AuditEvent behavior, operation behavior, or FEATURE-0013 adoption content beyond its approved `NOT_APPLICABLE` declaration.
9. A FEATURE-0014-specific architecture-boundary validator exists and passes against the approved architecture package before generation begins.
10. The working tree contains no unrelated or partially applied architecture migration that could enter Kiro context.

Failure of any item is `REPOSITORY_CONTEXT_NOT_READY`, not a requirements revision. Kiro must not generate or repair requirements until the repository-update step resolves the failure and the preflight is rerun.

### 16.9 Requirement normalization and duplication guard

Every normative requirement receives one canonical semantic key:

```text
(owning feature, resource or contract, actor, behavior, observable outcome)
```

Rules:

- Exactly one requirement is the normative home for a semantic key.
- User stories, risk rows, acceptance criteria, edge cases, and traceability tables reference the canonical requirement ID rather than restating its normative text.
- Shared inherited FEATURE-0012 behavior is referenced once as an inherited requirement family and specialized only where FEATURE-0014 adds a domain constraint.
- A requirement applying uniformly to several resources uses one parameterized requirement plus per-resource acceptance examples; it is not copied five times unless behavior materially differs.
- Security, concurrency, deletion, and isolation acceptance criteria attach to the owning requirement instead of creating semantically duplicate requirement families.
- Similar wording is not automatically duplication; the semantic key and observable outcome decide. Conversely, different wording does not hide duplication.
- Conflicting requirements sharing a semantic key stop generation with `ARCHITECTURE_DECISION_REQUIRED`.

The requirements review must produce a normalization ledger containing requirement ID, semantic key, upstream owner/decision, referenced risks, and duplicate disposition. Any duplicate without an explicit `REFERENCE_ONLY` disposition fails review.

### 16.10 FEATURE-0014 automated boundary gate

Before any Kiro stage is accepted, a feature-specific validator must check at minimum:

- exact enumeration of `F14-AD-001`–`F14-AD-021` and `F14-R01`–`F14-R30` with no missing, duplicate, renumbered, or range-only substitution;
- exactly five owned resource kinds and zero active aliases;
- absence of skipped-level, multiple-parent, mutable-parent, or cross-provider containment semantics;
- absence of inferred or explicit connectivity semantics;
- absence of ResourcePool, ProviderCapability, capacity, eligibility, compatibility, adapter, discovery, credential, endpoint, provider-call, plugin, operation, provisioning, or runtime semantics outside marked non-goal/ownership text;
- FEATURE-0013 adoption exactly `NOT_APPLICABLE` and no decision/audit requirement introduced by generic prompt text;
- complete single-owner overlap and requirement-normalization ledgers;
- no design-choice vocabulary in requirements unless the item is explicitly delegated to design and marked non-normative;
- no uncited normative requirement, orphan design element, orphan task, or missing risk evidence mapping;
- stage output and changed-file allowlist compliance;
- old-terminology and conflicting-context scans over every file included in Kiro and reviewer manifests.

The validator must distinguish normative use from explicit non-goal, ownership, migration, and risk text to avoid both false acceptance and superficial keyword failures. A validator configuration error blocks the stage rather than degrading to a warning.

### 16.11 Prompt and reviewer applicability rule

Generic Phase 2 prompts and reviewers may state cross-phase quality principles, but they must apply decision/audit/operation/observability checks only when the feature's approved applicability and behavior require them.

For FEATURE-0014:

- the valid FEATURE-0013 result is `NOT_APPLICABLE`;
- “explainable decision,” “defined audit behavior,” AuditEvent production, and operation lifecycle are not mandatory requirements or design sections;
- observability review is limited to contract validation, stable errors, redaction, and future implementation behavior actually authorized by approved requirements/design;
- a generic prompt or reviewer instruction conflicting with this architecture is a repository tooling defect and triggers `REPOSITORY_CONTEXT_NOT_READY`;
- Kiro must not satisfy a generic gate by adding out-of-scope domain semantics.

### 16.12 Architecture-readiness verdict

Architecture may be marked `READY_FOR_REQUIREMENTS` only when:

- no open semantic decision exists;
- the compatibility, ownership, source-precedence, stage-ownership, decision, and risk registers are internally consistent;
- all section 16.8 repository preflight checks pass;
- the FEATURE-0014 automated boundary gate is executable and passing;
- human review accepts the architecture and target residual-risk posture for requirements generation, without accepting final residual implementation risk.

Until then, status remains `PROPOSED` or `APPROVED_PENDING_REPOSITORY_ALIGNMENT`; Kiro requirements generation is not authorized.

## 17. Approval record

- Architecture status: Approved for Kiro requirements
- Human approver: Sanjeev Kumar
- Approval date: 2026-07-29
- Authorization after approval: Kiro may generate `requirements.md` only, subject to the accompanying ADH and normal stage gates.
