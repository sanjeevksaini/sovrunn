# Sovrunn Phase 2R Architecture Spine

**Status:** Approved
**Phase:** Phase 2R — Canonical Model PaaS Fabric Foundation
**Baseline:** `ARCH-2026.08-PHASE2R-CANONICAL`
**Approved date:** 2026-08-04
**Immediate feature:** FEATURE-0015 — Canonical Cloud Model Foundation
**Scope:** Architecture, contract, simulation, and conformance only
**Controlling rebaseline:** `docs/phase2/PHASE2R_REBASELINE.md`
**Canonical model:** `docs/architecture/canonical/sovrunn-finalized-data-model.md`
**Contract catalog:** `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`
**Cross-feature acceptance:** `docs/architecture/vertical-slices/VS-000-core-skeleton.md`
**Decisions:** DEC-0037 through DEC-0059

---

## 1. Executive Architecture Position

Phase 2R establishes an implementation-neutral, sovereignty-aware, auditable, and explainable decision fabric without performing infrastructure execution.

### 1.1 Controlling Architecture Flow

```text
Customer service intent
  → authenticate and authorize through adapter
  → CloudEnrollment, ServiceEntitlement, and QuotaPolicy reservation
  → immutable EffectiveGovernanceContext
  → immutable ServiceRequirementSet
  → qualified ExecutionTarget candidates from adapter-provided normalized facts
  → sovereignty DecisionRecord (registered profile)
  → placement DecisionRecord (registered profile)
  → immutable ServiceDeploymentPlan
  → accepted Operation
  → synthetic PluginExecution (fake in Phase 2R)
  → verified synthetic readiness
  → customer-safe ServicePlacement
  → per-consumer ServiceBinding with SecretRef only
  → correlated AuditEvents
  → bounded AI-readable explanation
```

During Phase 2R this flow ends at a conformance simulation result with no external side effects.

### 1.2 Phase 2R Must Not

- invoke a real provider or infrastructure API;
- provision PostgreSQL or any real database;
- execute real plugins against external systems;
- create production infrastructure;
- introduce durable persistence;
- implement billing or marketplace;
- perform failover, migration, or disaster recovery;
- perform autonomous AI operations;
- produce native AWS, OCI, OpenStack, Kubernetes, or OpenShift objects in customer or core schemas.

### 1.3 Architectural Separations

```text
Customer intent ≠ infrastructure implementation
Governance resolution ≠ policy engine implementation
Sovereignty assessment ≠ geographic location alone
Placement decision ≠ provisioning execution
Adapter boundary ≠ plugin execution boundary
DecisionRecord ≠ AuditEvent ≠ Operation
ServicePlacement (safe projection) ≠ canonical topology access
CloudPlatform (product ownership) ≠ CloudProvider (infrastructure operation)
Entitlement (what may be consumed) ≠ QuotaPolicy (how much)
```

---

## 2. Binding Architecture Invariants

### Invariant A — Reuse before build

Every feature classifies each significant capability as Reuse, Wrap, Extend, or Build with explicit Sovrunn-owned responsibility, external responsibility, adapter decision, phase impact, non-goals, replacement risk, and reassessment triggers.

### Invariant B — Implementation-neutral core

Core governance, placement, and decision logic must not depend on Kubernetes resource types, OpenStack, VMware, AWS, Azure, OCI concepts, provider names, provider-specific status values, PostgreSQL operator APIs, or vendor SDK objects. Provider-specific concepts appear only behind adapter-facing boundaries.

### Invariant C — Customer API boundary protection

Customer-facing resources express desired PaaS outcomes. Infrastructure topology, provider credentials, protected handles, and internal failure-domain models remain internal. Customer APIs expose only safe projections (ServicePlacement per DEC-0046).

### Invariant D — Decision before execution

No provisioning or lifecycle operation may begin without an applicable approved DecisionRecord. Phase 2R produces decisions but performs no execution beyond deterministic fake simulation.

### Invariant E — Engine-neutral policy contracts

Sovrunn owns PolicyEvaluationRequest, PolicyEvaluationResult, EffectiveGovernanceContext resolution, and DecisionRecord composition. OPA, Cedar, or another policy engine evaluates rules later through PolicyEngineAdapter. Engine-native types must not become Sovrunn core domain types.

### Invariant F — Adapter before integration

External engines expected to evolve or be replaced must be accessed through Sovrunn-owned adapter contracts (DEC-0036). Adapters translate between Sovrunn concepts and external systems without becoming a second domain model.

### Invariant G — Immutable published definitions

Published product, runtime, governance, sovereignty, and policy definitions are immutable by version (DEC-0044). Once published, a version cannot be altered. New versions are new resources.

### Invariant H — Seven canonical scope kinds

The governance scope vocabulary is Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider (DEC-0037). All scoped resources use metadata.scopeRef as sole scope authority.

### Invariant I — DecisionRecord, AuditEvent, and Operation remain distinct

| Record | Responsibility |
|---|---|
| DecisionRecord | Immutable governed conclusion with registered DecisionProfile extension, authority, evidence versions, typed result, and audit linkage (DEC-0043). |
| AuditEvent | Accountability: who or what acted, on which subject, with what outcome. Uses seven scope kinds through metadata.scopeRef. |
| Operation | Tracks an asynchronous lifecycle action, links to decisions and PluginExecution. |

These records reference each other but must not be collapsed into one object.

### Invariant J — Sovereignty is not geography

Physical geography, governance scope, and sovereignty are independent dimensions (DEC-0041). Geography alone never proves sovereignty. Sovereignty requires evidence-backed assessment through a registered DecisionRecord profile (DEC-0055).

### Invariant K — Explainability is structured

Reason codes, policy references, evidence versions, rejected alternatives, and suggested corrective actions must be machine-readable. Human-readable text is supplementary. AI-readable context is a bounded projection with no decision or execution authority.

### Invariant L — Earlier features own shared concepts

A later feature may consume or specialize an earlier contract but may not silently redefine it. Changing an earlier contract requires an Architecture Decision Handoff.

### Invariant M — Canonical bootstrap, not runtime migration

The first control-plane release exposes canonical resource contracts only (DEC-0059, supersedes DEC-0058). FEATURE-0001–0014 are retained repository assets and reuse input, not live state requiring conversion; there is no dual-writer cutover to manage because there is no alpha authority to coexist with.

---

## 3. Architecture Planes

| Plane | Responsibility | Primary features |
|---|---|---|
| Architecture governance | Reuse assessment, controls, feature-level decision discipline | FEATURE-0011 |
| Resource contract | Resource grammar, API boundaries, references, validation, status, conditions | FEATURE-0012 |
| Decision and accountability | DecisionRecord, DecisionProfile, AuditEvent, scope authority | FEATURE-0013 |
| Canonical cloud model | Seven-scope canonical bootstrap: CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, InfrastructureStack (ends here; no alpha runtime migration) | FEATURE-0015 |
| Integration boundary | Adapter contracts, ExecutionTarget qualification, normalized target facts, fake adapter | FEATURE-0016 |
| Policy context | Policy Evaluation Abstraction: PolicyEvaluationRequest/Result, PolicyEngineAdapter, DecisionRecord linkage | FEATURE-0017 |
| Governance and access | Governance, IAM, Approval and Exception Foundation: GovernanceProfile, Membership, Roles, Approval, Exception | FEATURE-0018 |
| Sovereignty evidence | Sovereignty Facts, Evidence and Policy Foundation: SovereigntyProfile, SovereigntyFactSet, EvidenceRecord, RegulatoryPolicyBundle | FEATURE-0019 |
| Effective governance resolution | ProfileAssignment, immutable EffectiveGovernanceContext | FEATURE-0020 |
| Customer enrollment and catalog access | CloudEnrollment, Personal Organization, EntitlementPackage, QuotaPolicy | FEATURE-0021 |
| Service product and requirements | ServiceTypeDefinition, ServiceOffering, ServicePlan, ServiceRuntimeProfile, ServiceRequirementSet | FEATURE-0022 |
| Sovereignty and placement decisions | Registered DecisionProfiles, DecisionRecords, candidate evaluation, ServicePlacement | FEATURE-0023 |
| Plugin classification and synthetic execution | Plugin/adapter manifests, compatibility, PluginExecution contract, fake harness | FEATURE-0024 |
| Explanation | Bounded audience-safe AI-readable projection | FEATURE-0025 |
| Integration conformance | VS-000 cross-feature orchestration, Slice 0 demo | FEATURE-0026 |

---

## 4. Phase 2R End State

At FEATURE-0026 completion, Sovrunn demonstrates through VS-000:

1. A customer Organization enrolled in a CloudPlatform requests a service using a published ServiceOffering and ServicePlan.
2. Sovrunn resolves CloudEnrollment, ServiceEntitlement, and QuotaPolicy.
3. Sovrunn composes an immutable EffectiveGovernanceContext from assigned profiles with non-weakenable controls.
4. The ServicePlan is translated into an immutable ServiceRequirementSet.
5. Qualified ExecutionTarget candidates are loaded with fresh normalized facts from the fake adapter.
6. A sovereignty DecisionRecord is created using evidence-backed SovereigntyProfile assessment.
7. A placement DecisionRecord selects or denies targets with exact evidence/policy versions.
8. An immutable ServiceDeploymentPlan is created referencing the selected target.
9. An Operation is accepted and executed through the fake plugin/adapter harness.
10. A customer-safe ServicePlacement is published with no topology or credential leakage.
11. A per-consumer ServiceBinding with fake SecretRef is issued.
12. Correlated AuditEvents and bounded AI-readable explanation are produced.
13. No external provider, Kubernetes, or database API is called.
14. Denial, stale-input, unauthorized, idempotency, fencing, redaction, and no-side-effect conformance paths all pass.

---

## 5. Dependency Graph

```text
FEATURE-0011–0014 completed baseline
              │
              ▼
FEATURE-0015 Canonical Cloud Model Foundation
              │
              ▼
FEATURE-0016 Adapter Boundary and ExecutionTarget Qualification
       ┌──────┴───────────┐
       ▼                  ▼
FEATURE-0017           FEATURE-0019
Policy Evaluation      Sovereignty Facts,
Abstraction            Evidence and Policy
       │               Foundation
       ▼                  │
FEATURE-0018              │
Governance, IAM,          │
Approval and Exception    │
Foundation                │
       └──────┬───────────┘
              ▼
FEATURE-0020 Assignment and Effective Governance Resolution
              │
              ▼
FEATURE-0021 CloudEnrollment, Personal Onboarding, Entitlement and Quota
              │
              ▼
FEATURE-0022 Service Product, Runtime and Requirement Foundation
              │
              ▼
FEATURE-0023 Sovereignty and Placement Decision v0
              │
              ▼
FEATURE-0024 Plugin Taxonomy and Synthetic Execution Boundary
              │
              ▼
FEATURE-0025 AI-Readable Decision and Operation Context
              │
              ▼
FEATURE-0026 Slice 0 Integration and Conformance Demo
```

FEATURE-0019 may proceed in parallel with FEATURE-0017/0018 after FEATURE-0016.

---

## 6. Shared Contract Boundaries

| Contract family | Owner | Primary consumers | Boundary rule |
|---|---|---|---|
| Organization, OrganizationUnit, Tenant, Project | Phase 1 retained | Governance resolution, enrollment, audit | Phase 2R references; does not redesign the hierarchy. |
| Common resource grammar, seven-scope vocabulary | FEATURE-0012 retained + DEC-0037 | Every later resource | Owns metadata, scopeRef, references, status, validation, conditions, Problem Details. |
| DecisionRecord, DecisionProfile, AuditEvent | FEATURE-0013 retained + DEC-0043 | Sovereignty, placement, governance, plugin, AI | Common immutable envelope with registered profiles. Specialized decisions extend it. |
| CloudPlatform, CloudProvider, CloudProviderParticipation, topology through InfrastructureStack | FEATURE-0015 | Adapter qualification, enrollment, placement, plugin | Implementation-neutral cloud model. No combined owner/operator concept. ExecutionTarget owned entirely by FEATURE-0016. |
| Adapter contracts, target qualification | FEATURE-0016 | Policy, sovereignty, placement, plugin | Normalized facts. No provider-native leakage into customer/core. |
| PolicyEvaluationRequest/Result | FEATURE-0017 | Governance resolution, placement | Engine-neutral. No OPA/Cedar native objects in core. |
| GovernanceProfile, Membership, Roles, Approval, Exception | FEATURE-0018 | Resolution, entitlement, placement | Declarative governance inputs and access control. |
| SovereigntyProfile, SovereigntyFactSet, EvidenceRecord | FEATURE-0019 | Sovereignty assessment, placement | Evidence-backed. Geography alone is insufficient. |
| ProfileAssignment, EffectiveGovernanceContext | FEATURE-0020 | Entitlement, placement, plugin | Immutable resolved context. Non-weakenable controls. |
| CloudEnrollment, EntitlementPackage, ServiceEntitlement, QuotaPolicy | FEATURE-0021 | Service request validation, placement | Independent evaluation: entitlement (what) and quota (how much) both deny independently. |
| ServiceTypeDefinition, ServiceOffering, ServicePlan, ServiceRequirementSet | FEATURE-0022 | Placement, plugin, binding | Immutable published versions. ServiceOffering references reusable ServiceTypeDefinition. |
| Sovereignty/Placement DecisionRecords, ServicePlacement | FEATURE-0023 | Operation, plugin, binding, explanation | Authoritative decisions. Safe customer projection. No topology leakage. |
| Plugin/adapter manifests, PluginExecution | FEATURE-0024 | Operation, binding | Fake only in Phase 2R. Cannot bypass decisions, operations, or SecretRef boundaries. |
| AI-readable decision/operation context | FEATURE-0025 | Customer UX, operations | Bounded projection. No decision, approval, or execution authority. |

---

## 7. Ownership Rules

### 7.1 Sovrunn Core Owns

```text
seven canonical scope kinds and governance scope vocabulary
customer service intent and product catalog semantics
effective governance resolution
sovereignty assessment coordination
implementation-neutral requirements and placement
decision outcome and audit linkage
operation lifecycle and plugin execution coordination
safe customer projections and SecretRef-only binding
explanation context
```

### 7.2 External Systems May Own (Later Phases)

```text
identity authentication
policy rule evaluation
secret storage and rotation
workflow execution orchestration
provider infrastructure provisioning
database lifecycle management
telemetry storage and alerting
event transport
production persistence
```

Sovrunn must not duplicate these mature engines. Architectural control is preserved through Sovrunn-owned contracts, adapters, decisions, and audit records.

---

## 8. Security Boundaries

- Customer APIs contain no provider-native objects, raw secrets, protected handles, or sensitive topology (DEC-0046).
- ServiceBinding is per-consumer, separately revocable, SecretRef-only (DEC-0051).
- Authorization uses authenticated PrincipalRef, canonical action, scoped RoleAssignment, and applicable policy. Persona labels grant no authority.
- Cross-scope denial uses SafeDenial (404) without existence disclosure.
- Secret values must not appear in responses, logs, audit events, decisions, or AI explanations.
- Provider credential isolation is enforced per SovrunnInstallation (DEC-0054).

---

## 9. Observability and Audit

- Structured logs and traces carry request_id and operation_id.
- Every meaningful decision links to an AuditEvent with actor, action, subject, scope, and outcome.
- DecisionRecords capture exact input versions for reproducibility.
- AuditEvents are append-only and immutable.
- Observability fields never contain secret material.
- PluginExecution records correlated to their parent Operation.

---

## 10. Performance and Scalability Direction

- Governance resolution and placement are designed as deterministic, stateless evaluations over immutable inputs.
- In-memory registry for Phase 2R; persistence is a later-phase concern.
- Published definitions are immutable by version — no cache-invalidation complexity for published resources.
- Horizontal scaling: control-plane components are stateless request handlers over shared registry.
- No synchronous network calls in the hot evaluation path during Phase 2R simulation.

---

## 11. Phase 2R Non-Goals

Phase 2R must not:

- invoke real providers or infrastructure APIs;
- provision real databases or workloads;
- execute real plugins;
- implement production persistence;
- implement billing, marketplace, or UI;
- perform real autoscaling, failover, or DR;
- perform autonomous AI operations;
- produce native provider objects in customer/core schemas;
- introduce mandatory capacity scheduling or provider-wide capability truth;
- create dual authority between alpha and canonical models;
- promote APIs to stable.

---

## 12. Feature Factory Controls

### 12.1 Mandatory Reuse Assessment

Every Phase 2R feature must include a reuse assessment conforming to `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` before implementation.

### 12.2 Architecture Drift Gate

The drift gate must reject:

- active target contracts using superseded concepts (DEC-0042 supersedes mandatory placement pools/capabilities);
- provider-specific hardcoding in core;
- embedded custom policy engines;
- raw secret storage;
- customer-facing topology leakage;
- fewer than seven canonical scope kinds in active authorities;
- generic combined owner/operator concepts;
- global catalog concepts without ServiceTypeDefinition/ServiceOffering structure;
- provider-native objects in customer or core schemas.

### 12.3 Feature Gate

A feature is not complete unless:

- tests pass (happy + failure paths);
- lint/security checks pass;
- reuse assessment is approved;
- architecture drift checks pass;
- acceptance criteria from PHASE2R_REBASELINE are met;
- Phase 2R scope boundaries are respected;
- human approval review is complete.

### 12.4 VS-000 Conformance (FEATURE-0026)

The final integration feature must pass:

- positive flow;
- denial flows (authentication, authorization, enrollment, entitlement, quota, sovereignty, placement);
- stale-input flows (generation mismatch, expired evidence, unqualified target);
- idempotency (same key same request; same key different request);
- fencing (maintenance epoch, draining target);
- redaction (no secret, protected handle, or sensitive topology in any customer response, log, or explanation);
- no-side-effect (no real external calls).

---

## 13. Cross-Feature Contracts

### C01 — No Silent Redefinition

A feature may consume an earlier model but must not change its semantics. Required changes return through Architecture Decision Handoff.

### C02 — Reference Instead of Duplication

Shared objects are linked by stable references. Decision reproducibility uses immutable snapshots (exact input versions), not duplicated objects.

### C03 — Engine-Neutral Policy Contract

All policy evaluation goes through PolicyEngineAdapter. Handlers, placement logic, and resolution must not contain an independent hidden policy engine.

### C04 — EffectiveGovernanceContext Is the Resolution Boundary

Profile inheritance and assignment are resolved by FEATURE-0020 before placement and sovereignty assessment. Consumers receive the immutable resolved context; they do not independently traverse the governance hierarchy.

### C05 — Sovereignty Is Evidence-Based

SovereigntyAssessment uses a registered DecisionRecord profile with explicit SovereigntyFactSet and EvidenceRecord inputs. Missing, stale, or tampered evidence causes indeterminate or deny outcomes, never success.

### C06 — Placement Evaluates Qualified Targets

Placement evaluates qualified ExecutionTarget candidates with adapter-provided normalized facts and availability axes. It does not select Kubernetes nodes, VMs, or provider-native resources directly.

### C07 — ServicePlacement Is a Safe Projection

Customers receive ServicePlacement containing only platform region, selected provider (if disclosure policy permits), and summarized decision outcome. No internal target identifiers, credentials, or protected handles.

### C08 — Entitlement and Quota Are Independent

Entitlement (what may be consumed) and quota (how much) evaluate independently. Both may deny; denial by either is sufficient (DEC-0039).

### C09 — ServicePlan Remains Customer-Facing

ServicePlan expresses a customer-consumable offering. ServiceRuntimeProfile and ServiceRequirementSet own internal capability and placement requirements. Placement evaluates explicit requirements, not plan names.

### C10 — AI Is a Consumer, Not an Authority

AI-readable context may explain or summarize existing authorized decisions. It must not change decisions, bypass policy, insert unvalidated facts, initiate execution, or expose secrets.

### C11 — No Side Effects in Phase 2R

All models and evaluators are deterministic and side-effect-free except for in-memory test recording. No provider calls, plugin execution against real systems, or infrastructure mutation.

---

## 14. Acceptance Criteria Summary

Phase 2R is complete when:

1. VS-000 definition of done and conformance matrix pass.
2. Sovereignty and placement simulation demonstrates selected, denied, requires-approval, and indeterminate outcomes.
3. Every Phase 2R feature passes its feature gate.
4. DecisionRecord and AuditEvent semantics remain owned by FEATURE-0013.
5. Customer contracts contain no provider-native objects, raw secrets, or protected handles.
6. No active desired-state API supports old and canonical kinds/scopes simultaneously.
7. Seven canonical scope kinds are enforced across all active authorities.
8. Phase 3 readiness review approves replacement of fake execution with one real PostgreSQL path.

---

## Approval Record

This architecture spine was approved as part of ADH-2026-042 canonical model adoption on 2026-08-04 under ACR-2026-001 and ARCH-APPROVAL-2026-004. It replaces the original Phase 2 spine.

The authoritative feature sequence is `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`.
The authoritative rebaseline is `docs/phase2/PHASE2R_REBASELINE.md`.
The cross-feature acceptance charter is `docs/architecture/vertical-slices/VS-000-core-skeleton.md`.
