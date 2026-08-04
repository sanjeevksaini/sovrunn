# Sovrunn Integrated Vertical-Slice Delivery Plan

**Status:** Architecture-owner approved target delivery rebaseline — repository integration required
**Version:** 1.2
**Date:** 4 August 2026
**Repository reviewed:** `/Users/sanjeevkumar/SoftwareDevelopment/sovrunn`
**Reviewed branch:** `phase2-reuse-first-paas-fabric-foundation`
**Reviewed commit:** `7bced1068f2837ead9709e78d3fc33f8e6e032c5`
**Change classification:** Coordinated correction, replacement and extension of the roadmap after FEATURE-0014

**Companion architecture:** [Finalized Data Model](sovrunn-finalized-data-model.md) · [Canonical Contract Catalog](sovrunn-final-canonical-contract-catalog.md) · [Canonical ADR Index](sovrunn-final-canonical-adrs/README.md) · [PostgreSQL Reference Flow](sovrunn-postgresql-reference-flow.md)

## 1. Decision

Sovrunn will use the finalized canonical architecture as the design guardrail and implement it incrementally through production-shaped vertical slices.

```text
Canonical architecture
        ↓ constrains
Phases
        ↓ fund outcomes
Features
        ↓ deliver capabilities
Vertical slices
        ↓ prove integration and customer value
Conformance evidence
        ↓ permits API promotion
Stable platform contracts
```

This avoids both failure modes:

- implementing the entire data model before proving customer value; and
- allowing each feature or slice to invent its own incompatible model.

The canonical model is defined before implementation. Exact APIs, state machines, controllers and persistence behavior become normative only when the relevant slice implements and validates them.

## 2. Authority and approval status

This plan does not change the live repository baseline by itself.

The repository currently identifies Phase 2 as active and FEATURE-0015 `ResourcePool and ProviderCapability Model` as next. The finalized target model instead separates `CloudPlatform` product ownership from `CloudProvider` supply participation, adopts externally grounded hosting facts and `ExecutionTarget`, and excludes a mandatory `ResourcePool` or provider-wide infrastructure capability resource.

Therefore:

1. FEATURE-0001–0014 remain immutable implementation history.
2. Existing resources through FEATURE-0014 are migrated only through approved handoffs and an explicit alpha compatibility plan.
3. FEATURE-0015 requirements must not start from the current placeholder.
4. ADH-2026-020 through ADH-2026-041 are approved as a coordinated architecture amendment through `ARCH-APPROVAL-2026-004`.
5. The repository architecture baseline, decisions, roadmap, feature index, phase sequence and traceability matrix must now be updated together through `ACR-2026-001` and `ADH-2026-042`.
6. Phase 2R feature names and allocations are approved target architecture but become repository-authoritative only after that controlled integration is merged.

## 3. Planning-unit definitions

| Unit | Purpose | Normative output |
|---|---|---|
| Canonical architecture | Defines stable meanings, ownership, boundaries, relationships and invariants | Data model, contract catalog and ADRs |
| Phase | Funds and governs a measurable platform or business outcome | Phase charter and exit gate |
| Feature | Delivers one independently reviewable capability | Architecture contract, requirements, design, tasks, implementation and feature-gate evidence |
| Vertical slice | Proves multiple features through one end-to-end provider/customer journey | Slice charter, fixtures, executable scenario and conformance result |
| Conformance gate | Demonstrates that contracts work together without violating boundaries | Positive, negative, security, lifecycle and operability evidence |
| API maturity | Controls compatibility promises based on evidence | Experimental, Alpha, Beta or Stable designation |

A feature may support more than one slice. A slice normally depends on several features. A phase may contain multiple slices and parallel hardening streams.

## 4. Verified repository baseline

| Baseline | Verified status | Delivery implication |
|---|---|---|
| Phase 0 | Foundation baseline exists | Retain source-of-truth, AI workflow, architecture governance and feature-factory controls |
| Phase 1 | FEATURE-0001–0010 implemented | Reuse Organization, OrganizationUnit, Tenant, Project, Operation, catalog, plugin, ServiceInstance, ServiceBinding, health and demo foundations; migrate rather than duplicate |
| Phase 2 partial | FEATURE-0011–0014 implemented and merged | Reuse assessment, common resource grammar, DecisionRecord/AuditEvent and current topology are implementation assets |
| Current repository placeholder | FEATURE-0015 ResourcePool and ProviderCapability | Superseded by the approved Phase 2R target; must not be used to generate requirements |
| Current Phase 3 anchor | FEATURE-0027–0034 PostgreSQL execution chain | Retain the narrow PostgreSQL objective, but align contracts and names with the finalized model |
| Current MVP | Governed PostgreSQL placement and provisioning on one substrate | Use as the first real service slice |
| Phase 4+ | Roadmap placeholders | Rebaseline from vertical-slice evidence before implementation |

### 4.1 Assets retained

- Reuse-before-build assessment.
- Common metadata, typed references, status, condition and error grammar.
- DecisionRecord and AuditEvent envelope.
- Organization, Tenant, Project, Operation, ServiceInstance and ServiceBinding implementations.
- Feature gates, architecture drift checks and conformance-fixture approach.
- PostgreSQL as the first real managed-service proof.

### 4.2 Required migrations

- Classify generic `Provider` data into separate `CloudPlatform` product ownership and `CloudProvider` operation, joined through `CloudProviderParticipation`; do not preserve a combined authority.
- Provider-prefixed fixed geography to `HostingLocation`, `Datacenter`, `FaultDomain` and optional `InfrastructureStack` references.
- Mandatory `ResourcePool`/`ProviderCapability` plan to qualified `ExecutionTarget`, normalized facts and adapter qualification.
- Global `ServiceClass` semantics to reusable `ServiceTypeDefinition` plus CloudPlatform-scoped `ServiceOffering`.
- `EffectivePolicyContext` to `EffectiveGovernanceContext`.
- Placement and sovereignty outcomes to registered DecisionRecord profiles.
- Existing ServiceBinding to per-consumer, SecretRef-only issuance, rotation and revocation.

## 5. Target phase and slice map

| Phase | Status | Feature count | Slice coverage | Exit outcome |
|---|---|---:|---|---|
| Phase 0 — Foundation and AI development system | Completed baseline | Not numbered as product features | Input to all slices | Architecture and AI-assisted delivery remain source-controlled, human-approved and test-gated |
| Phase 1 — Platform core skeleton | Completed | 10 | Foundation input | Common governance and service-resource skeleton exists |
| Phase 2 partial — Current fabric foundation | Completed through FEATURE-0014 | 4 | Foundation input | Reuse, API grammar, decision/audit and topology implementation assets exist |
| Phase 2R — Canonical alignment and executable skeleton | Approved target; repository merge pending | 12 approved feature slots | Slice 0 | Corrected canonical model is executable through a synthetic service without real backend provisioning |
| Phase 3 — First real service execution | Proposed rebaseline of existing eight slots | 8 candidate feature slots | Slice 1 | PostgreSQL 17 is provisioned on one qualified Yotta OpenShift target and safely bound to one consumer |
| Phase 4 — First Production MVP | Proposed replacement of current hardening-only scope | 16 candidate feature groups | Slices 2, 3 and 4 plus the parallel platform-lifecycle slice | PostgreSQL is sovereign, HA-capable, Day-2 manageable and supported by a recoverable production control plane |
| Phase 5 — Relationships and portability | Proposed | 6 candidate feature groups | Slice 5 | The same core contract supports service relationships, a second backend and a second service family |
| Phase 6 onward — Product expansion and scale | Directional | Revalidated per phase | Additional service and platform slices | Provider scale, advanced resilience, cost, assurance, AI operations, industry products and SDE expand from proven contracts |

Feature count means independently gated capability groups, not implementation tasks or person-months.

## 6. Phase 2R — canonical alignment and Slice 0

### 6.1 Goal

Correct the model before FEATURE-0015 and prove the common framework with an in-process synthetic service, fake plugin and fake adapter. No real provider or PostgreSQL provisioning occurs.

### 6.2 Approved target feature mapping

FEATURE-0015–0026 are unimplemented placeholders and are the least disruptive rebaselining window. Their target names are approved through `ARCH-APPROVAL-2026-004`; repository authority begins after `ADH-2026-042` is integrated and merged.

| Feature ID | Approved target feature | Principal contracts | Depends on | Primary evidence |
|---|---|---|---|---|
| FEATURE-0015 | Canonical Cloud Model and Alpha Migration Foundation | CloudPlatform, CloudProvider, CloudProviderParticipation, canonical scope vocabulary, CanonicalMigrationPlan and migration reports | Approved ADH-020/037/041 and FEATURE-0012 | Deterministic classification and cutover fixtures; no dual Provider authority; completed history preserved |
| FEATURE-0016 | Adapter Boundary and ExecutionTarget Qualification | Adapter contracts, normalized facts, qualification status, target-scoped credentials | FEATURE-0015 | Fake adapter and target qualification tests |
| FEATURE-0017 | Policy Evaluation Abstraction | Policy evaluation request/result, PolicyEngineAdapter, DecisionRecord linkage | FEATURE-0013/0016 | Deterministic allow/deny/unknown fixtures |
| FEATURE-0018 | Governance, IAM, Approval and Exception Foundation | GovernanceProfile, Membership, RoleDefinition, RoleAssignment, PrivilegedAccessRequest, AccessReview, ApprovalPolicy, ApprovalRequest, ExceptionGrant | FEATURE-0012/0017 | Scope, least-privilege, just-in-time access, review and separation-of-duties tests |
| FEATURE-0019 | Sovereignty Facts, Evidence and Policy Foundation | SovereigntyProfile, RegulatoryPolicyBundle, SovereigntyFactSet, EvidenceRecord | FEATURE-0013/0016 | Versioning, provenance, freshness and no-geography-equals-sovereignty tests |
| FEATURE-0020 | Assignment and Effective Governance Resolution | ProfileAssignment, EffectiveGovernanceContext | FEATURE-0018 | Deterministic inheritance, conflict and exception fixtures |
| FEATURE-0021 | CloudEnrollment, Personal Onboarding, Entitlement and Quota | Personal Organization/default Tenant/Project, CloudEnrollment, EntitlementPackage, ServiceEntitlement, QuotaPolicy/reservation, provider-selection intent | Catalog migration and FEATURE-0020 | Simple personal and enterprise onboarding; CloudPlatform enrollment remains independent from provider participation, entitlement and quota decisions |
| FEATURE-0022 | Service Product, Runtime and Requirement Foundation | ServicePortfolio, ServiceTypeDefinition, ServiceOffering, ServicePlan migration, ServicePlacementProfile, ServiceRuntimeProfile, ServiceRequirementSet | ADH-023/027/032 and FEATURE-0016 | Synthetic service type and portfolio added without core schema change |
| FEATURE-0023 | Sovereignty and Placement Decision v0 | Sovereignty and placement DecisionRecord profiles, ServicePlacement safe projection | FEATURE-0019–0022 | Explainable selected/denied/indeterminate simulation |
| FEATURE-0024 | Plugin Taxonomy and Synthetic Execution Boundary | Service plugin, execution adapter, manifest, compatibility and protected-handle boundaries | FEATURE-0016/0022 | Fake plugin/adapter conformance |
| FEATURE-0025 | AI-Readable Decision and Operation Context | Bounded explanation projection from decisions, evidence, operations and audit | FEATURE-0013/0023 | AI can explain but cannot decide, approve or execute |
| FEATURE-0026 | Slice 0 Integration and Conformance Demo | Synthetic ServiceInstance → decisions → plan → Operation → fake execution → projection/binding | All Phase 2R features | End-to-end positive, negative and leakage tests |

The coordinated migration of completed FEATURE-0006–0014 is part of the approved architecture migration package; it must not be disguised as unrelated new resources inside FEATURE-0015.

### 6.3 Slice 0 charter

**Outcome:** An authenticated user creates one synthetic managed service and traces it through canonical contracts using a fake plugin and adapter.

**Must prove:**

```text
authenticated request
  → scoped authorization
  → enrollment, entitlement and quota decision
  → effective governance
  → requirements
  → synthetic target qualification
  → sovereignty and placement DecisionRecords
  → immutable deployment plan
  → Operation
  → fake PluginExecution
  → safe placement projection
  → per-consumer fake ServiceBinding
  → correlated AuditEvents and AI-readable explanation
```

**Non-goals:** Real Kubernetes/OpenShift calls, real PostgreSQL, production policy/identity/secrets engines, multi-target HA and autonomous AI.

**Exit gate:** Every canonical writer, reference and boundary used in Slice 0 has deterministic positive, negative, idempotency, authorization and leakage tests.

## 7. Phase 3 — Slice 1: thin PostgreSQL happy path

### 7.1 Goal

Replace the synthetic executor with one real PostgreSQL service path while keeping the customer contract implementation-neutral.

### 7.2 Candidate mapping of existing Phase 3 slots

| Candidate ID | Rebaselined feature | Purpose |
|---|---|---|
| FEATURE-0027 | PluginExecution Contract v0 | Immutable execution request, attempt, normalized result, Operation and target linkage |
| FEATURE-0028 | Operation Orchestration Controller v0 | Idempotent plan execution, retries, cancellation, compensation placeholder, audit and status |
| FEATURE-0029 | PostgreSQL Management Plugin v0 | Interpret PostgreSQL service contract and plan lifecycle without implementing database internals |
| FEATURE-0030 | OpenShift Execution Adapter v0 | Translate approved abstract actions into protected OpenShift/Kubernetes objects |
| FEATURE-0031 | PostgreSQL Runtime Integration v0 | Wrap the selected mature PostgreSQL operator/distribution for create, readiness, endpoint and delete |
| FEATURE-0032 | Service Realization Pipeline v0 | Connect final decisions to immutable plan, Operation, PluginExecutions and ServiceInstance status |
| FEATURE-0033 | ServiceBinding, Workload Identity and SecretRef Integration | Issue one consumer-specific safe endpoint and secret reference; support revoke-on-delete |
| FEATURE-0034 | Slice 1 PostgreSQL Conformance Demo | Run the thin Yotta–Department of Posts PostgreSQL journey end to end |

### 7.3 Slice 1 limits

- One CloudProvider: Yotta.
- One customer Organization: Department of Posts.
- One Tenant and one Project.
- One PostgreSQL offering and plan.
- PostgreSQL major version 17.
- One ServiceRegion.
- One qualified OpenShift ExecutionTarget.
- One selected mature PostgreSQL operator/distribution.
- One governance profile and one bounded sovereignty profile.
- Provision, verify, bind and delete.
- A small critical set of denial, retry, idempotency and secret-leakage tests.

### 7.4 Slice 1 definition of done

```text
Customer request
  → authenticated and authorized
  → enrollment, entitlement and quota verified
  → plan and PostgreSQL version validated
  → governance and bounded sovereignty evaluated
  → one target selected
  → immutable plan created
  → Operation accepted
  → PostgreSQL actually provisioned through the adapter/operator
  → readiness verified
  → customer-safe status and placement published
  → consumer-specific ServiceBinding issued
  → delete revokes access and cleans up safely
```

## 8. Phase 4 — First Production MVP

Phase 4 must become more than demo hardening. A sovereign production PaaS requires evidence-backed placement, resilience, Day-2 service operations and independently recoverable Sovrunn lifecycle management.

Final feature IDs are intentionally not assigned until the Phase 2R amendment is approved.

### 8.1 Slice 2 — evidence-backed sovereign placement

| Candidate feature group | Capability |
|---|---|
| P4-SOV-01 | Production normalized target and deployed-service facts with freshness and provenance |
| P4-SOV-02 | Effective-dated RegulatoryPolicyBundle evaluation and authoritative sovereignty DecisionRecord |
| P4-SOV-03 | Customer-safe sovereignty assurance projection, reassessment and stale/missing-evidence denial |

**Outcome:** Sovrunn can prove why a PostgreSQL placement satisfies the selected sovereignty outcome at a particular time, without claiming that physical geography alone proves sovereignty or legal compliance.

### 8.2 Slice 3 — HA and multi-target placement

| Candidate feature group | Capability |
|---|---|
| P4-HA-01 | Multi-target requirements, failure-domain and explicit connectivity evidence |
| P4-HA-02 | Atomic target-set placement, HA profile and independent backup-target realization |
| P4-HA-03 | InfrastructureMaintenanceNotice, target availability state, epoch/fencing, impact assessment, failover/keep/migrate/stop and requalification |

**Outcome:** PostgreSQL can use at least two failure domains and an independently approved backup target while respecting latency, connectivity and sovereignty constraints.

### 8.3 Slice 4 — Day-2 managed-service lifecycle

| Candidate feature group | Capability |
|---|---|
| P4-LIFE-01 | Backup, point-in-time recovery and verified restore |
| P4-LIFE-02 | Minor update, major upgrade and explicit rollback/restore/forward-recovery behavior |
| P4-LIFE-03 | Credential rotation, evidence-triggered reassessment and governed migration |

**Outcome:** The service is operable after creation and can recover from expected lifecycle and infrastructure events.

### 8.4 Production hardening stream

| Candidate feature group | Capability |
|---|---|
| P4-PROD-01 | Persistent stores, HA control-plane deployment and data migration |
| P4-PROD-02 | Real identity, authorization, secret-provider and workload-identity integrations |
| P4-PROD-03 | Observability, SLOs, audit export, supportability, threat model, load/race/chaos and release gates |

### 8.5 Parallel platform-lifecycle slice

| Candidate feature group | Capability |
|---|---|
| P4-PLAT-01 | SovrunnInstallation, signed SovrunnRelease, ComponentManifest, ReleaseCompatibilityContract and PlatformLifecyclePolicy |
| P4-PLAT-02 | Typed lifecycle references, explicit recovery modes, immutable PlatformLifecyclePlan and Operation-first atomic desired-state activation |
| P4-PLAT-03 | Single external lifecycle authority, independently recoverable agent, credentials and backups |
| P4-PLAT-04 | Install, upgrade, failed-upgrade recovery, restore with control plane unavailable and platform-health projection |

This stream starts in Phase 2R with interface and deployment-topology decisions, proceeds alongside service slices, and is mandatory for the Phase 4 production gate.

### 8.6 First Production MVP release gate

The MVP is production-ready only when all of the following are demonstrated:

1. A real customer-scoped PostgreSQL service can be provisioned, verified, bound and deleted.
2. Placement is based on current facts, explicit evidence, applicable policy and an authoritative sovereignty assessment.
3. HA and backup placement satisfy explicit failure-domain, latency, connectivity and sovereignty constraints.
4. Backup restoration and at least one upgrade/recovery path are proven.
5. Raw secrets and protected backend identifiers never enter customer or canonical core records.
6. Authorization is least-privileged and scope-bound; privileged access is time-bound and audited.
7. The Sovrunn control plane can be upgraded and restored through the external lifecycle agent when the affected API is unavailable.
8. Security, availability, observability, supportability and data-retention expectations have measurable SLOs and evidence.
9. All required feature gates and slice conformance suites pass.
10. A design partner accepts the customer and provider operating journeys and known limitations.

## 9. Phase 5 — Slice 5: relationships and portability

| Candidate feature group | Capability |
|---|---|
| P5-REL-01 | ServiceRelationshipDefinition and Project-scoped ServiceRelationship lifecycle |
| P5-REL-02 | Relationship-derived requirements, authorization and DecisionRecord profiles |
| P5-REL-03 | Implementation-neutral private-connectivity realization through plugins/adapters |
| P5-PORT-01 | Second materially different ExecutionTarget adapter |
| P5-PORT-02 | Second materially different service family through ServiceTypeDefinition |
| P5-PORT-03 | Plugin certification and cross-implementation portability conformance |

**Outcome:** A VM-like application service can privately consume PostgreSQL through the same core relationship contract on two implementations, with no VPC, VCN, subnet, route, security-group or NetworkPolicy types entering the canonical core.

API stabilization should normally wait until at least one canonical extension surface has been exercised by two materially different implementations or service families.

## 10. Phase 6 onward — directional expansion

Later phases remain outcome portfolios until production and customer evidence justify exact features.

| Future phase | Recommended outcome | Candidate capability areas | Rebaseline trigger |
|---|---|---|---|
| Phase 6 — Provider and plugin scale | Repeatable onboarding and certified ecosystem | Plugin trust, manifests, compatibility, health, credentials, provider onboarding and certification | Two providers or adapter partners |
| Phase 7 — Advanced resilience and data movement | Governed regional and cross-cloud service continuity | DR profiles, replication, traffic, cross-location data movement and advanced failover | Production HA evidence and customer RPO/RTO demand |
| Phase 8 — Capacity, autoscaling and economics | Governed scaling and commercial consumption | Capacity facts, quota, autoscaling, cost estimates, guardrails, usage, rating and charge records | Reliable usage and capacity sources |
| Phase 9 — Continuous sovereign assurance | Exportable and continuously refreshed assurance | Control mappings, evidence collectors, exceptions, attestations and audit export | Qualified legal/compliance validation |
| Phase 10 — AI-native operations | Governed recommendations and bounded automation | Recommendations, runbooks, risk, approval, autonomy policy, memory and explanations | Sufficient deterministic telemetry and runbooks |
| Phase 11 — Multi-service and industry clouds | Reusable core powers differentiated products | Application, data, AI, messaging, streaming, DNS, identity and industry portfolios | PostgreSQL and second-service portability proven |
| Phase 12 — SDE, federation and enterprise scale | Advanced data platform and multi-provider operation | SDE, service supply/federation, fleet operations and enterprise scale | Separate product validation and scale evidence |

The existing roadmap's later placeholders should be mined for useful capability intent, but their current phase and resource assumptions are not authoritative after this rebaseline.

## 11. Slice documentation set

### 11.1 Integrated plan

This document owns phase, feature-group and slice mapping. It must not contain detailed API schemas.

### 11.2 One bounded charter per slice

```text
vertical-slices/
  VS-000-core-skeleton.md
  VS-001-postgresql-happy-path.md
  VS-002-sovereign-placement.md
  VS-003-ha-multitarget.md
  VS-004-day2-lifecycle.md
  VS-005-relationships-portability.md
  VS-PLATFORM-production-lifecycle.md
```

Every charter contains:

- provider and customer outcome;
- goals and explicit non-goals;
- participating canonical contracts;
- contributing features and dependencies;
- actors and authorization boundaries;
- test fixture graph;
- happy path and critical failure paths;
- security and information-disclosure tests;
- operational and observability requirements;
- demo script and definition of done;
- evidence required for API maturity promotion.

### 11.3 Feature specification package

Every feature retains the repository feature-factory lifecycle:

```text
approved architecture handoff
  → feature architecture contract
  → reuse assessment
  → requirements
  → design
  → tasks
  → implementation and tests
  → feature gate
```

### 11.4 Contract-to-feature traceability

The repository traceability matrix should add these columns:

| Field | Purpose |
|---|---|
| Canonical contract | Identifies the architecture concept implemented or consumed |
| ADR/ADH | Identifies the governing decision |
| Owning feature | Identifies the sole implementation owner |
| First slice | Identifies the first end-to-end proof |
| Reused by slices | Prevents duplicate feature ownership |
| Conformance fixture | Links to executable evidence |
| API maturity | Shows Experimental, Alpha, Beta or Stable status |
| Supersedes/migrates | Makes compatibility changes explicit |

## 12. API maturity and promotion policy

| Maturity | Minimum evidence | Compatibility promise |
|---|---|---|
| Experimental | Architecture reviewed and synthetic fixture works | May change without migration guarantee |
| Alpha | One real slice passes positive, negative and security conformance | Breaking changes allowed only through explicit migration |
| Beta | Production-shaped use, lifecycle recovery and operational evidence exist | Compatibility preserved unless an approved exceptional migration is required |
| Stable | Multiple implementations/service families prove extensibility and production SLOs | Strong compatibility and deprecation policy |

No contract becomes stable merely because its schema has been written. Promotion depends on conformance evidence.

## 13. Conformance model

Every slice must produce evidence across applicable gates:

| Gate | What it proves |
|---|---|
| Contract | Schema, references, ownership, mutability and state-machine invariants |
| Authorization | Identity, scope, least privilege, no-existence disclosure and separation of duties |
| Governance | Enrollment, entitlement, quota, policy resolution, approval and exception behavior |
| Sovereignty | Applicable policy, fact/evidence freshness, assessment and safe assurance projection |
| Decision | Deterministic outcomes, exact input versions, alternatives, reasons and audit linkage |
| Execution | Immutable plan, idempotent Operation, PluginExecution and protected native handles |
| Security | Secret redaction, credential isolation, tenant boundary and supply-chain provenance |
| Lifecycle | Create, verify, update, backup, restore, rotate, recover and delete behavior as in scope |
| Operability | Logs, metrics, traces, audit, SLO, diagnostics and support runbooks |
| Portability | Same canonical contract works on a second implementation without core-native branching |

## 14. Dependency and parallel-development rules

### 14.1 Must remain sequential

- Architecture approval precedes migration or new FEATURE-0015 requirements.
- Common resource and scope migration precedes dependent APIs.
- Service contract and governance resolution precede placement.
- Sovereignty assessment precedes final placement when sovereignty applies.
- Final decisions precede immutable deployment planning.
- Operation acceptance precedes PluginExecution.
- Verified readiness precedes ServiceBinding issuance.
- Platform release compatibility, plan and approval precede lifecycle activation.

### 14.2 Safe parallel streams

- Plugin SDK and adapter SDK after their common execution contract is fixed.
- PostgreSQL management plugin, OpenShift adapter and operator integration after FEATURE-0027 boundaries are fixed.
- Sovereignty policy content and evidence collectors after canonical schemas are fixed.
- HA placement and Day-2 lifecycle after Slice 1 produces stable execution contracts.
- Platform lifecycle implementation in parallel with service slices after the external authority boundary is approved.
- Threat modeling, conformance fixtures, observability and documentation alongside each feature rather than at phase end.

Parallel teams may share approved contracts; they may not create parallel versions of those contracts.

## 15. Anti-drift rules

1. Customers request service outcomes, plans, regions and optional placement/sovereignty profiles—not backend-native resources.
2. CloudProviders publish service products and register reusable execution targets; Sovrunn does not claim ownership of underlying infrastructure.
3. Physical hosting facts, desired sovereignty, regulatory interpretation, evidence and assessment remain separate.
4. Core contains no AWS, OCI, OpenStack, OpenShift, Kubernetes or vendor-specific API branches.
5. `ExecutionTarget` is the actionable realization boundary; no mandatory `ResourcePool` is reintroduced without a separately demonstrated requirement and approved ADR.
6. New services extend through ServiceTypeDefinition, plugins and adapters rather than universal core fields.
7. Service relationships and access bindings remain separate.
8. DecisionRecord remains the common decision envelope.
9. AI explains, recommends and assists; deterministic policy, approval and controllers retain authority.
10. Platform lifecycle remains externally recoverable and does not make Sovrunn the owner of underlying infrastructure lifecycle.

## 16. Required decisions before implementation begins

The following must be resolved or explicitly bounded before Slice 0/1 work is accepted:

1. Approval and repository integration of ADH-2026-020–041.
2. Exact alpha migration strategy for Provider → CloudProvider and Phase 1 catalog types.
3. Official feature names, IDs and dependency order after FEATURE-0014.
4. First production ExecutionTarget type and supported OpenShift/Kubernetes versions.
5. PostgreSQL operator/distribution reuse decision and supported PostgreSQL version policy.
6. Policy, identity and secret-provider choices for Slice 1 versus production MVP.
7. Exact typed-reference, condition, reason-code and authoritative-writer registries.
8. Feature-level realization of the now-fixed Operation retry, cancellation, compensation, recovery and concurrency contracts.
9. Physical product selection and deployment of the fixed external lifecycle authority, signing, approval and recovery contracts.
10. First Production MVP SLO values and contractual limitations.

## 17. Completeness and accuracy validation

### 17.1 Source validation performed

This plan was checked against:

- the approved repository vision;
- the current architecture version and baseline;
- the current phase context and architecture session rules;
- the repository development phases, roadmap and feature index;
- Phase 1 and Phase 2 feature sequences;
- the current PostgreSQL MVP definition;
- the finalized data model and canonical contract catalog;
- the canonical ADH package through ADH-2026-036;
- the architecture-closure decisions ADH-2026-037 through ADH-2026-041;
- the repository baseline conflict review; and
- the Yotta–Department of Posts PostgreSQL reference flow.

### 17.2 Coverage validation

| Canonical area | First modeled | First real proof | Production gate | Later portability proof |
|---|---|---|---|---|
| Common resource grammar | Existing/Phase 2R migration | Slice 1 | Phase 4 | Slice 5 |
| CloudPlatform products, provider participation and customer governance | Phase 2R | Slice 1 | Phase 4 | Industry portfolios and additional providers |
| IAM and zero trust | Phase 2R | Slice 1 | Production IAM integration | Provider/federation scale |
| Hosting facts and ExecutionTarget | Phase 2R | Slice 1 | Slice 3 qualification/requalification | Second backend |
| Governance and decisions | Phase 2R | Slice 1 | Slices 2–4 | All later services |
| Sovereignty | Phase 2R foundation | Bounded Slice 1 evaluation | Slice 2 full evidence | Continuous assurance |
| Service realization and binding | Phase 2R simulation | Slice 1 | Slice 4 lifecycle | Second service family |
| Multi-target HA | Contract preserved in model | Not required in Slice 1 | Slice 3 | Advanced resilience |
| Service relationships | Contract preserved in model | Not required in MVP | Not required for first PostgreSQL MVP | Slice 5 |
| Platform lifecycle | Boundary decided in Phase 2R | Parallel agent prototype | Mandatory Phase 4 recovery proof | Fleet management |
| AI | Structured context in Phase 2R | Explanation in Slice 1 | Guarded operational assistance | AI-native operations |
| Composite-service model | Canonical boundary reserved | Not required in Slice 1 | Not required for first PostgreSQL MVP | Validate after relationships using independently governed child services |
| Metering/federation/governed data | Reserved extension boundaries | Deferred | Deferred unless first contract requires it | Promoted by market evidence |

### 17.3 Validation verdict

The plan has governed architecture approval, roadmap rebaselining and Slice 0 charter creation. It is not yet repository- or implementation-authoritative because:

- the approved canonical handoffs are outside the live repository and await controlled integration;
- official post-0014 feature allocation is approved in the target package but not yet merged into repository authorities;
- several reuse selections and production SLOs remain intentionally undecided; and
- exact schemas and state machines belong to feature and slice specifications.

No canonical domain is left without an implementation phase, explicit deferral or promotion gate at roadmap level. Exact schemas, state machines and product SLOs remain intentionally open. The most significant repository incompatibilities are explicitly exposed rather than hidden: Provider/CloudProvider, fixed provider topology, ResourcePool/ProviderCapability, global ServiceClass, governance-context vocabulary and the timing of sovereignty/platform-lifecycle capabilities.

## 18. Immediate next actions

1. **Completed:** approve the canonical ADH-2026-020–041 package.
2. **Completed:** produce the coordinated repository integration manifest and roadmap rebaseline.
3. **Pending repository merge:** replace the current FEATURE-0015 placeholder before its architecture work begins.
4. **Completed as target architecture:** allocate Phase 2R feature names and dependencies.
5. **Completed:** create the Slice 0 charter and contract-to-feature traceability.
6. **Next after repository integration:** generate FEATURE-0015 requirements, then execute Slice 0 incrementally across Phase 2R.
7. Create the thin Slice 1 charter from the larger PostgreSQL reference flow.
8. Start the platform-lifecycle design stream early, but keep its production gate in Phase 4.

---

**Authority note:** This is an architecture-owner approved target delivery and traceability plan. It preserves the repository's authority rules and becomes repository-authoritative only through the reviewed `ADH-2026-042` integration merge.
