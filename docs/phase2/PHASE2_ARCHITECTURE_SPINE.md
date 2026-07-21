# Sovrunn Phase 2 Architecture Spine

**Status:** Approved
**Phase:** Phase 2 — Reuse-First PaaS Fabric Foundation
**Approved date:** 2026-07-21
**Immediate feature:** FEATURE-0011 — Reuse Assessment Standard
**Scope:** Architecture and decision preparation only
**Controlling strategy:** `docs/phase2/PHASE2_EXECUTION_STRATEGY.md`

---

## Executive Architecture Position

Phase 2 establishes a provider-neutral, policy-aware, reusable, auditable, and explainable decision fabric without performing infrastructure execution.

The controlling architecture flow is:

```text
Customer service intent
    -> effective governance and policy context
    -> entitlement check
    -> service runtime requirements
    -> provider-neutral capability matching
    -> policy evaluation
    -> placement decision
    -> audit record
    -> AI-readable explanation
```

During Phase 2, this flow ends at an explainable simulation result.

Phase 2 must not:

- invoke a provider;
- provision PostgreSQL;
- execute plugins;
- create production infrastructure;
- introduce production persistence;
- implement billing;
- perform failover or disaster recovery;
- perform autonomous AI operations.

The architecture spine protects these separations:

```text
Customer intent != infrastructure implementation
Policy context != policy engine implementation
Placement decision != provisioning execution
Adapter boundary != plugin execution boundary
Decision record != audit record != operation record
```

---

# 1. Phase 2 Architecture Spine

## 1.1 Phase 2 End State

At the completion of FEATURE-0026, Sovrunn should be able to simulate this scenario:

1. A tenant or project requests a customer-facing PostgreSQL `ServicePlan`.
2. Sovrunn resolves the applicable governance, security, data-placement, cost, and entitlement context.
3. The `ServicePlan` is translated into provider-neutral runtime requirements.
4. Available `ResourcePools` are evaluated through declared `ProviderCapabilities`.
5. A policy-engine-neutral evaluation result is produced.
6. Sovrunn issues an explainable `PlacementDecision`.
7. The decision is linked to an `AuditEvent`.
8. An AI-readable context explains why the request was allowed, denied, or marked as requiring approval.
9. No provider or runtime operation is executed.

This is a **governed decision simulation**, not a provisioning workflow.

## 1.2 Binding Architecture Invariants

### Invariant A — Reuse before build

Every feature must classify each significant capability as:

```text
Reuse
Wrap
Extend
Build
```

The classification must state:

- what Sovrunn owns;
- what remains delegated;
- whether an adapter is required;
- what is explicitly deferred;
- what would cause the decision to be reassessed.

### Invariant B — Provider-neutral core

Core policy, placement, and governance logic must not depend on:

- Kubernetes resource types;
- OpenStack, VMware, AWS, or Azure concepts;
- provider names;
- provider-specific status values;
- PostgreSQL operator APIs;
- vendor SDK objects.

Provider-specific concepts may appear only behind provider-facing resources, adapters, or future plugin implementations.

### Invariant C — Customer API boundary protection

Customer-facing resources express desired PaaS outcomes.

Provider details, IaaS topology, infrastructure credentials, internal failure-domain models, and low-level capabilities remain provider-facing, internal, or plugin-facing unless explicitly exposed through a reviewed advanced API.

### Invariant D — Decision before execution

No future provisioning or lifecycle operation may begin without an applicable approved decision context.

Phase 2 produces decisions but performs no execution.

### Invariant E — Engine-neutral policy contracts

Sovrunn owns:

- policy inputs;
- effective policy context;
- evaluation requests;
- normalized evaluation results;
- decision composition;
- audit linkage.

OPA, Cedar, or another policy engine may evaluate rules later, but its native types must not become Sovrunn core domain types.

### Invariant F — Adapter before integration

External engines expected to evolve or be replaced must be accessed through Sovrunn-owned adapter contracts.

Adapters translate between Sovrunn concepts and external systems. They must not become a second domain model or duplicate the external engine.

### Invariant G — Plugin metadata before plugin execution

Phase 2 may define:

- plugin taxonomy;
- plugin manifests;
- declared capabilities;
- trust metadata;
- compatibility metadata;
- execution boundaries.

Phase 2 must not define or invoke the complete plugin execution chain.

### Invariant H — Decision, audit, and operation remain distinct

| Record | Responsibility |
|---|---|
| `DecisionObject` | Captures what Sovrunn decided, why it decided it, and which inputs were considered. |
| `AuditEvent` | Captures accountability: who or what acted, on which subject, and with what outcome. |
| `Operation` | Tracks an asynchronous lifecycle action or attempted change. |

A record may reference the others, but they must not be collapsed into one object.

### Invariant I — Explainability is structured

Reason codes, policy references, rejected alternatives, and suggested corrective actions must be machine-readable.

Human-readable text is supplementary and must not be the only explanation.

### Invariant J — Earlier features own shared concepts

A later Phase 2 feature may consume or specialize an earlier contract but may not silently redefine it.

Changing an earlier contract requires an Architecture Decision Handoff and, where appropriate, an Architecture Change Request, DEC, or RFC update.

## 1.3 Architecture Planes

| Plane | Responsibility | Primary features |
|---|---|---|
| Architecture governance plane | Reuse assessment, architecture controls, and feature-level decision discipline | FEATURE-0011 |
| Resource contract plane | Resource grammar, API boundaries, references, validation, status, and conditions | FEATURE-0012 |
| Decision and accountability plane | Common decision structure, `AuditEvent` structure, and linkage rules | FEATURE-0013 |
| Provider-neutral substrate plane | Provider topology, `ResourcePools`, and capability declarations | FEATURE-0014–0015 |
| Integration boundary plane | Replaceable boundaries for external engines and stores | FEATURE-0016 |
| Policy-context plane | Policy evaluation abstraction, profiles, assignments, effective context, and entitlement | FEATURE-0017–0021 |
| Service-intent plane | Translation from customer `ServicePlan` to internal runtime requirements | FEATURE-0022 |
| Placement plane | Candidate evaluation and explainable placement outcome | FEATURE-0023 |
| Plugin classification plane | Metadata-only plugin roles and trust/execution boundaries | FEATURE-0024 |
| Explanation plane | AI-readable, policy-governed decision context | FEATURE-0025 |
| Integration simulation plane | Demonstrates the complete side-effect-free decision flow | FEATURE-0026 |

## 1.4 Core Phase 2 Control Flow

```text
Existing Phase 1 governance scope
Organization -> OrganizationUnit -> Tenant -> Project
                           |
                           v
Customer-facing ServicePlan and service request
                           |
                           v
ServiceRuntimeProfile
                           |
          +----------------+----------------+
          |                                 |
          v                                 v
EffectivePolicyContext              ServiceEntitlement
          |                                 |
          +----------------+----------------+
                           |
                           v
PlacementRequest
                           |
            +--------------+--------------+
            |                             |
            v                             v
ResourcePool catalogue          ProviderCapabilities
            |                             |
            +--------------+--------------+
                           |
                           v
PolicyEvaluationRequest
                           |
                           v
PolicyEvaluationResult
                           |
                           v
PlacementDecision
              +------------+------------+
              |                         |
              v                         v
          AuditEvent          AI-readable DecisionContext
```

The flow must not continue to a provider adapter or plugin execution during Phase 2.

## 1.5 Architecture Ownership Rule

Sovrunn core owns the meaning of:

```text
governance scope
customer service intent
effective policy context
runtime requirements
provider-neutral capabilities
placement request
decision outcome
audit linkage
explanation context
```

External systems may own enforcement or execution later, including:

```text
identity authentication
policy rule evaluation
secret storage
workflow execution
provider provisioning
database lifecycle management
telemetry storage
event transport
production persistence
```

Sovrunn must not duplicate these mature engines merely to maintain architectural control. Architectural control is preserved through Sovrunn-owned contracts, adapters, decisions, and audit records.

---

# 2. Phase 2 Dependency Graph

FEATURE-0011 is a governance dependency for every later Phase 2 feature.

It does not create a runtime model dependency, but every later feature must conform to its assessment contract.

```text
FEATURE-0011 Reuse Assessment Standard
    |
    +--> FEATURE-0012 API, Resource Naming, Status and Validation Standard
    |        |
    |        +--> FEATURE-0013 Decision Object and AuditEvent Standard
    |        |
    |        +--> FEATURE-0014 Provider-Neutral Resource Model
    |        |        |
    |        |        +--> FEATURE-0015 ResourcePool and ProviderCapability
    |        |
    |        +--> FEATURE-0016 Adapter Boundary Foundation
    |
    +--> FEATURE-0017 Policy Evaluation Abstraction
             |
             +--> FEATURE-0018 GovernanceProfile and SecurityProfile
             |
             +--> FEATURE-0019 DataPlacementPolicy and CostGuardrail
                      |
                      +--> FEATURE-0020 ProfileAssignment and EffectivePolicyContext
                               |
                               +--> FEATURE-0021 ServiceEntitlement and Quota Placeholder

FEATURE-0015 + FEATURE-0021
    |
    +--> FEATURE-0022 ServiceRuntimeProfile

FEATURE-0013
FEATURE-0015
FEATURE-0020
FEATURE-0021
FEATURE-0022
    |
    +--> FEATURE-0023 PlacementRequest and PlacementDecision

FEATURE-0016 + FEATURE-0022
    |
    +--> FEATURE-0024 Plugin Taxonomy Foundation

FEATURE-0013 + FEATURE-0023
    |
    +--> FEATURE-0025 AI-Readable Decision Context

FEATURE-0023 + FEATURE-0024 + FEATURE-0025
    |
    +--> FEATURE-0026 Phase 2 Integration Demo
```

## 2.1 Dependency Interpretation

### Standards dependencies

FEATURE-0012 and FEATURE-0013 establish grammar used by later resources.

Later features should not invent independent:

- metadata;
- reference;
- status;
- condition;
- validation;
- decision;
- audit conventions.

### Provider dependencies

FEATURE-0015 depends on FEATURE-0014 because `ResourcePool` and `ProviderCapability` require a provider-neutral location and infrastructure context.

### Policy dependencies

Profiles cannot produce `EffectivePolicyContext` until the policy evaluation abstraction and profile types exist.

Entitlement consumes the resolved context rather than independently resolving governance inheritance.

### Placement dependencies

FEATURE-0023 must not begin until the following inputs exist:

- `DecisionObject` and `AuditEvent` standard;
- `ResourcePool`;
- `ProviderCapability`;
- `EffectivePolicyContext`;
- `ServiceEntitlement`;
- quota placeholder;
- `ServiceRuntimeProfile`.

### Plugin dependency

Plugin taxonomy depends on the adapter boundary and runtime-profile concepts, but it must not alter their semantics.

### Integration dependency

FEATURE-0026 consumes all prior outputs.

It must not compensate for missing contracts by adding private demo-only models.

---

# 3. Phase 2 Shared Object and Model Boundaries

## 3.1 Boundary Map

| Model family | Owner | Primary consumers | Boundary rule |
|---|---|---|---|
| `Organization`, `OrganizationUnit`, `Tenant`, `Project` | Phase 1 baseline | Policy context, entitlement, requests, and audit | Phase 2 references these objects; it does not redesign the governance hierarchy. |
| `ServiceClass` and `ServicePlan` | Phase 1 baseline | `ServiceRuntimeProfile` and service request | Remain customer-facing catalog concepts. Infrastructure requirements must not be added directly to the customer API by default. |
| `ReuseAssessment` | Architecture Operating System | Every Phase 2 feature | Documentation and review contract, not a runtime API resource. |
| Common resource grammar | FEATURE-0012 | Every later resource model | Owns metadata, references, status, validation, and API boundary classification. |
| `DecisionObject` | FEATURE-0013 | Policy, placement, and future decision types | Common decision envelope. Specialized decisions extend or compose it without redefining common semantics. |
| `AuditEvent` | FEATURE-0013 | Decisions, operations, and future execution | Accountability record. It is not a log line or operation status. |
| Provider hierarchy | FEATURE-0014 | `ResourcePool` and placement | Describes provider-neutral topology and sovereignty context. |
| `ResourcePool` | FEATURE-0015 | Runtime compatibility and placement | Unit considered for placement. It is not necessarily a Kubernetes cluster or cloud region. |
| `ProviderCapability` | FEATURE-0015 | Runtime matching and placement | Normalized compatibility declaration. Provider-native capability data remains outside core. |
| Adapter contracts | FEATURE-0016 | Policy, identity, secrets, operations, observability, events, and repositories | Core depends on Sovrunn contracts; adapters own translation to external APIs. |
| `PolicyEvaluationRequest` and `PolicyEvaluationResult` | FEATURE-0017 | Profile evaluation and placement | Engine-neutral. No OPA or Cedar native objects in the core contract. |
| `GovernanceProfile` and `SecurityProfile` | FEATURE-0018 | `ProfileAssignment` and `EffectivePolicyContext` | Declarative policy inputs, not executable policy engines. |
| `DataPlacementPolicy` and `CostGuardrail` | FEATURE-0019 | `EffectivePolicyContext` and placement | Minimal decision inputs only; no production compliance or billing engine. |
| `ProfileAssignment` | FEATURE-0020 | Policy-context resolver | Associates profiles with governance scopes. It does not duplicate profile contents. |
| `EffectivePolicyContext` | FEATURE-0020 | Entitlement, policy evaluation, and placement | Resolved, read-only decision input for a particular request context. |
| `ServiceEntitlement` and quota placeholder | FEATURE-0021 | `PlacementRequest` validation | Answers whether the service request may proceed to placement consideration. It does not implement billing or full quota accounting. |
| `ServiceRuntimeProfile` | FEATURE-0022 | Placement and plugin classification | Internal bridge from `ServicePlan` to provider-neutral runtime and capability requirements. |
| `PlacementRequest` | FEATURE-0023 | Placement evaluator | Immutable evaluation input assembled from previously owned contracts. |
| `PlacementDecision` | FEATURE-0023 | Audit, API response, and AI explanation | Specialized `DecisionObject` representing placement outcome. It does not execute provisioning. |
| Plugin taxonomy models | FEATURE-0024 | Future plugin execution | Classification, manifest, and trust metadata only during Phase 2. |
| `DecisionContext` | FEATURE-0025 | AI explanation and human review | Sanitized structured projection of decisions and evidence. It must not contain secrets or unrestricted tenant data. |

## 3.2 Adapter Versus Plugin Boundary

Adapters and plugins solve different architectural problems.

| Adapter | Plugin |
|---|---|
| Hides or normalizes an external technical system. | Represents a Sovrunn lifecycle or execution role. |
| Called by core through a replaceable interface. | Participates in a future governed operation chain. |
| Examples: policy engine, identity provider, repository. | Examples: provider plugin, service management plugin, service runtime plugin. |
| May be infrastructure used by multiple domains. | Has declared service or provider responsibilities. |
| FEATURE-0016 establishes contracts. | FEATURE-0024 establishes taxonomy and metadata only. |

A Provider/Substrate Plugin may internally use provider SDK adapters in a future phase, but the concepts must not be merged.

Plugin execution is explicitly outside Phase 2.

## 3.3 Decision Versus Policy Result Boundary

A `PolicyEvaluationResult` is an input to a Sovrunn decision.

It is not automatically the final `DecisionObject` because Sovrunn may need to combine:

- multiple policy evaluations;
- entitlement;
- runtime compatibility;
- capability state;
- approval requirements;
- candidate comparison;
- rejected alternatives.

The placement component owns the resulting `PlacementDecision`.

The policy adapter owns only the normalized policy evaluation result.

---

# 4. Phase 2 Cross-Feature Contracts

## 4.1 Purpose

These contracts define how independently implemented Phase 2 features remain architecturally compatible.

They are not runtime workflows by themselves.

Each contract includes a conceptual example to improve AI reasoning and human understanding.

> **Conceptual-example rule:** Every example in this section is illustrative only. The examples do not authorize implementation, runtime provisioning, vendor selection, provider calls, plugin execution, or expansion of the approved Phase 2 scope.

## P2-C01 — Mandatory Reuse Assessment

Every FEATURE-0011-and-later architecture contract, Kiro requirements, design, tasks, and final review must include a reuse assessment.

The assessment must be approved before implementation starts.

### Conceptual example — not execution scope

FEATURE-0017 introduces a policy evaluation abstraction.

Its assessment might state:

```text
Build:
Sovrunn-owned PolicyEvaluationRequest and PolicyEvaluationResult contracts.

Wrap:
A future external policy engine through PolicyEngineAdapter.

Reuse:
A mature policy language or engine implementation.

Deferred:
Selecting and integrating the first production policy engine.
```

This example explains architectural responsibility only. It does not select OPA, Cedar, or another product.

## P2-C02 — Capability-Level Classification

Reuse classification applies to an architectural capability or decision unit, not blindly to an entire feature.

A feature may legitimately contain several dispositions.

### Conceptual example — not execution scope

FEATURE-0020 may contain:

```text
Build:
The EffectivePolicyContext domain contract.

Extend:
The existing Phase 1 governance hierarchy through profile references.

Reuse:
Common reference and metadata standards from FEATURE-0012.

Wrap:
A future external policy source if profiles are later loaded externally.
```

The feature should not receive one broad `Build` label because that would hide reused and wrapped components.

## P2-C03 — No Silent Redefinition

A feature may consume an earlier model but must not change its semantics through private fields, duplicate types, or feature-local terminology.

Required changes must return through the Architecture Decision Handoff process.

### Conceptual example — not execution scope

FEATURE-0023 may consume a capability such as:

```text
storage.encryption.atRest = true
```

It may not introduce a second placement-only provider-feature model with different meanings for encryption, sovereignty region, or availability.

If the earlier capability contract is insufficient, FEATURE-0023 must request a reviewed change.

## P2-C04 — Reference Instead of Duplication

Shared objects should normally be linked by stable references rather than embedded copies.

Snapshotting may be used for decision reproducibility, but references and evaluation snapshots must be explicitly distinguished.

### Conceptual example — not execution scope

A `PlacementRequest` may reference:

```text
ServiceRuntimeProfile
EffectivePolicyContext
ServiceEntitlement
```

A `PlacementDecision` may also contain an immutable summary of the values evaluated at decision time.

The reference identifies the source object. The snapshot records the evaluated state.

## P2-C05 — Provider-Neutral Domain Contract

Core objects use normalized concepts and capabilities.

Provider-native identifiers belong in provider-facing extensions, adapter state, plugin-private configuration, or explicitly opaque metadata.

### Conceptual example — not execution scope

Sovrunn core may express:

```text
Residency jurisdiction: IN
Storage encryption at rest: required
Minimum availability domains: 2
Managed backup capability: required
```

Core should not require an AWS availability-zone ID, OpenStack aggregate UUID, Kubernetes `StorageClass`, or VMware datastore-cluster name.

## P2-C06 — Engine-Neutral Policy Contract

All policy evaluation goes through the `PolicyEngineAdapter` contract.

Handlers, placement logic, and profile resolvers must not contain an independent hidden policy engine.

### Conceptual example — not execution scope

A normalized policy request may state:

```text
Subject: Tenant A
Action: Request PostgreSQL plan
Context:
- Data must remain in India
- Maximum approved service tier is Standard
```

A normalized result may state:

```text
Outcome: Allowed with approval
Reason code: COST_APPROVAL_REQUIRED
```

Core must not depend on engine-native query syntax or package paths.

## P2-C07 — Effective Context Is the Resolution Boundary

Policy inheritance and profile assignment are resolved before placement.

Placement consumes `EffectivePolicyContext`; it does not independently traverse the governance hierarchy.

### Conceptual example — not execution scope

Assume:

```text
Organization: requires encryption
OrganizationUnit: requires India residency
Tenant: restricts service tier to Standard
Project: requires manual approval above a cost threshold
```

FEATURE-0020 resolves these into one `EffectivePolicyContext`.

FEATURE-0023 consumes that resolved context rather than repeating hierarchy traversal.

## P2-C08 — ServicePlan Remains Customer-Facing

`ServicePlan` expresses a customer-consumable offering.

`ServiceRuntimeProfile` owns internal capability, topology, and runtime requirements.

Placement must not infer infrastructure requirements from plan names.

### Conceptual example — not execution scope

A customer may select:

```text
ServicePlan: postgresql-standard
```

The internal runtime profile may express:

```text
Minimum storage class: durable
Backup capability: required
High availability: not required
Supported PostgreSQL family: 16-compatible
```

Placement evaluates explicit requirements rather than inferring a particular VM size, operator, or storage backend.

## P2-C09 — ResourcePool Is the Placement Boundary

Placement selects, rejects, or requests approval for `ResourcePools`.

It must not select Kubernetes nodes, virtual machines, operator custom resources, or provider-native availability zones directly.

### Conceptual example — not execution scope

Sovrunn may evaluate:

```text
ResourcePool: india-west-managed-services
```

That pool may later map internally to clusters, OpenStack projects, VMware resource groups, or another substrate. Those mappings are outside the Phase 2 placement contract.

## P2-C10 — Capabilities Drive Compatibility

Compatibility is based on normalized `ProviderCapabilities` and validation state, not provider brands or hardcoded substrate types.

### Conceptual example — not execution scope

A runtime profile may require:

```text
encrypted durable storage
automated backup support
two availability domains
PostgreSQL 16 compatibility
```

Multiple providers may satisfy these requirements. Placement evaluates capabilities rather than hardcoded provider allowlists for technical compatibility.

## P2-C11 — Decision and Audit Linkage

Every meaningful Phase 2 decision must be capable of linking to an `AuditEvent`.

The final linkage mechanism belongs to FEATURE-0013.

### Conceptual example — not execution scope

A decision may state:

```text
Decision: Approval required
Reason: The candidate satisfies technical requirements but exceeds the project cost guardrail.
```

The linked audit event may record:

```text
Actor: Placement evaluator
Subject: Placement request 123
Outcome: Approval required
Decision reference: Placement decision 456
```

The decision explains the result. The audit event records accountability.

## P2-C12 — Explainability Contract

Denials and approval requirements must include:

- stable reason codes;
- human-readable reasons;
- relevant policy or requirement references;
- rejected alternatives where applicable;
- safe corrective or next-step suggestions.

### Conceptual example — not execution scope

```text
Reason code: DATA_RESIDENCY_UNSATISFIED

Human explanation:
No eligible ResourcePool satisfies the required India residency policy.

Rejected candidates:
- Pool A: jurisdiction mismatch
- Pool B: capability data not validated

Suggested corrective action:
Select a permitted jurisdiction or onboard a compliant ResourcePool.
```

The reason code supports machines. The explanation supports human review. The suggestion does not automatically change policy.

## P2-C13 — No Side Effects

Models and evaluators built during Phase 2 must be deterministic and side-effect-free except for in-memory or test-oriented recording required by the integration simulation.

No provider calls, plugin execution, or real provisioning are allowed.

### Conceptual example — not execution scope

FEATURE-0026 may evaluate a static set of pools and produce:

```text
PlacementDecision: Pool B selected in simulation
```

It must not create a database, call Kubernetes, allocate a VM, execute a provider plugin, or mutate production infrastructure.

## P2-C14 — AI Is a Consumer, Not an Authority

`DecisionContext` may explain or summarize an existing decision.

It must not:

- change the decision;
- bypass policy;
- insert unvalidated provider facts;
- initiate execution;
- expose secrets or unauthorized tenant information.

### Conceptual example — not execution scope

An AI assistant may say:

```text
The request was denied because none of the validated pools meet the required India data-residency rule.
```

It must not change the result, invent a compliant pool, or trigger provisioning.

## P2-C15 — Traceability

Every feature must trace its architecture to:

- relevant accepted DEC records;
- related RFC or architecture documents;
- dependencies;
- reuse assessment;
- phase scope;
- acceptance criteria;
- feature-gate evidence.

### Conceptual example — not execution scope

FEATURE-0017 may trace:

```text
Architecture basis:
- DEC-0026 Reuse Before Build
- DEC-0036 Adapter Boundaries
- FEATURE-0012 common resource conventions
- FEATURE-0013 decision conventions

Reuse assessment:
- Build Sovrunn policy contract
- Wrap external policy engine

Acceptance evidence:
- Adapter contract tests
- No vendor-native type leakage
- Feature-gate approval
```

This is documentation traceability only and does not select an engine.

---

# 5. Risks if FEATURE-0011 Is Designed Incorrectly

## 5.1 Risk Register

| Risk | Architectural consequence |
|---|---|
| One classification is forced for an entire feature | Composite features become falsely classified and important reuse or adapter boundaries are hidden. |
| Classification is vendor-first | Phase 2 prematurely selects products before Sovrunn contracts stabilize. |
| `Reuse` is treated as direct coupling | External APIs and native types leak into core and increase replacement cost. |
| `Wrap` has no ownership boundary | Sovrunn wrappers grow into duplicated policy, workflow, or provider engines. |
| `Extend` permits uncontrolled forks | Sovrunn inherits long-term maintenance and supply-chain risk without explicit approval. |
| `Build` requires no evidence | Build becomes a loophole for unnecessary custom infrastructure. |
| Assessment records only a label | The strict gate becomes a keyword-presence check rather than an architecture control. |
| No sovereignty, security, or operational criteria | A mature component may be technically capable but unsuitable for sovereign or disconnected deployment. |
| No adapter decision | Later features couple directly to chosen implementations. |
| No phase-impact field | Assessments may justify runtime integration during Phase 2. |
| No non-goals | Future integrations expand current feature scope. |
| No replacement-risk analysis | Temporary choices become permanent accidentally. |
| No reassessment triggers | Decisions remain stale after requirements, licensing, maintenance, or deployment changes. |
| Assessment is mistaken for approval | Kiro or Cursor may treat a candidate as an accepted decision. |
| Templates maintain independent formats | The standard, prompts, and gates drift apart. |
| Automation replaces semantic review | A structurally valid but architecturally weak assessment passes. |
| FEATURE-0011 becomes a runtime resource | Architecture governance unnecessarily enters the runtime domain model. |
| Later features override assessments silently | Sequential architecture control is lost. |

## 5.2 Phase 2 Risk-Mitigation Framework

FEATURE-0011 establishes controls that every future feature applies through:

```text
Architecture discussion
    -> Architecture Decision Handoff
    -> human approval
    -> Kiro requirements
    -> Kiro design
    -> Kiro tasks
    -> Cursor implementation
    -> tests
    -> feature gate
    -> commit and pull request
```

### Level 1 — Architecture discussion controls

Before a handoff is produced, identify:

- capability boundaries;
- reuse candidates;
- Sovrunn-owned responsibilities;
- external responsibilities;
- adapter needs;
- phase constraints;
- non-goals;
- replacement risks;
- unresolved questions.

### Level 2 — Architecture Decision Handoff controls

The handoff must state:

- proposed disposition;
- assessment unit;
- alternatives;
- adapter boundaries;
- current and deferred work;
- affected files and features;
- acceptance criteria;
- explicit Kiro instructions.

### Level 3 — Kiro specification controls

Kiro requirements, design, and tasks must preserve the approved assessment.

Kiro must not:

- select an unapproved vendor;
- change a disposition;
- move deferred work into the current phase;
- introduce vendor-native domain types;
- remove non-goals;
- combine unrelated features.

Any required change returns to architecture review.

### Level 4 — Implementation and test controls

Cursor implementation must be checked for:

- conformance with the approved boundary;
- absence of provider-specific leakage;
- absence of unapproved integrations;
- absence of runtime side effects;
- absence of duplicate shared models;
- deterministic behavior;
- contract-level test coverage.

### Level 5 — Feature-gate controls

The strict Phase 2 gate should reject:

- missing reuse assessments;
- invalid disposition values;
- missing non-goals;
- missing adapter decisions;
- missing phase-impact statements;
- missing replacement risks;
- missing mitigation fields;
- missing traceability;
- unapproved architecture changes.

The gate provides structural enforcement. It does not replace human architectural judgment.

### Level 6 — Pull-request and human review controls

The final reviewer verifies:

- implementation matches the approved handoff;
- architecture acceptance criteria are met;
- no later feature was implemented early;
- no deferred integration was introduced;
- no shared contract was silently redefined;
- traceability is current;
- gate evidence is attached.

## 5.3 Risk-by-Risk Mitigation

| Risk | Mandatory mitigation |
|---|---|
| Whole-feature classification | Require feature summary plus capability-level assessments. |
| Vendor-first classification | Define contracts and responsibilities before selecting products. |
| Direct coupling under `Reuse` | Require an adapter decision and prohibit native-type leakage without approved exception. |
| Undefined `Wrap` boundary | Require explicit Sovrunn, external, data, and control boundaries. |
| Uncontrolled `Extend` fork | Treat Extend as supported extension or composition; require separate approval for forks. |
| Unjustified `Build` | Require written rejection of Reuse, Wrap, and Extend and a differentiator rationale. |
| Label-only assessment | Require rationale, boundaries, adapter decision, phase impact, non-goals, mitigation, replacement risk, and traceability. |
| Missing sovereign fit | Require proportional sovereignty, deployment, security, operations, and licensing analysis. |
| Missing adapter | Make adapter requirement explicitly Yes or No with rationale. |
| Phase leakage | Separate current-phase and deferred work and validate against Phase 2 non-goals. |
| Missing non-goals | Require non-goals in the assessment and Kiro specifications. |
| Missing replacement planning | Require risk level, exit boundary, and migration considerations. |
| Missing reassessment triggers | Require triggers such as license change, end-of-life, security failure, or altered deployment targets. |
| Assessment mistaken for approval | Require decision status and human approval before Kiro implementation planning. |
| Template drift | Maintain one canonical standard and align dependent prompts and checks. |
| Automation replacing review | Separate structural validation from human semantic approval. |
| Runtime `ReuseAssessment` | Explicitly prohibit it in scope, non-goals, Kiro instructions, and gates. |
| Silent future override | Require a new Architecture Decision Handoff for changed dispositions or boundaries. |

## 5.4 Future-Feature Mitigation Record

Every later Phase 2 feature must document:

```text
Relevant architecture risks:
Which FEATURE-0011 risks apply?

Preventive controls:
What design constraints prevent them?

Detection controls:
Which tests, reviews, or gates detect violations?

Corrective path:
What happens if a violation or gap is found?

Residual risk:
What remains after controls?

Reassessment trigger:
What change requires renewed architecture review?
```

### Conceptual example — not execution scope

For FEATURE-0017:

```text
Relevant risk:
Policy-engine types leak into Sovrunn core.

Preventive control:
Define engine-neutral request and result contracts.

Detection control:
Review domain models and tests for vendor-native types.

Corrective path:
Return through Architecture Decision Handoff if the abstraction is insufficient.

Residual risk:
Some engine capabilities may not map cleanly to the initial contract.

Reassessment trigger:
A required policy feature cannot be represented without engine-specific behavior.
```

---

# 6. Current Repository Risks and Mitigation

## 6.1 Template Divergence

### Risk

The canonical draft may contain fields omitted by Feature Factory or reviewer prompts, creating competing definitions.

### Mitigation

1. Declare `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` canonical.
2. Add a standard version identifier.
3. Update templates and prompts to reference it.
4. Avoid independent shortened definitions.
5. Add consistency checks where duplication is unavoidable.
6. Require dependent-template updates with canonical-field changes.
7. Include template-alignment evidence in the FEATURE-0011 gate.

### Future-feature control

Later features must use the canonical format and may not introduce feature-local assessment structures.

## 6.2 Superficial Feature-Gate Validation

### Risk

A gate may only check for a `Reuse Assessment` heading.

### Mitigation

The gate should validate at least:

```text
Assessment scope
Valid disposition
Decision status
Sovrunn responsibility
External responsibility
Adapter decision
Phase impact
Non-goals
Mitigation controls
Replacement risk
Traceability
```

Conditional checks should include:

```text
Build:
Requires reasons for rejecting Reuse, Wrap, and Extend.

External replaceable engine:
Requires adapter rationale or approved exception.

Phase 2 feature:
Must acknowledge applicable Phase 2 non-goals.
```

Human review remains responsible for semantic credibility.

## 6.3 Traceability Mismatch

### Risk

FEATURE-0011 may reference unrelated decisions, causing humans and AI tools to reason from the wrong baseline.

### Mitigation

1. Associate FEATURE-0011 with DEC-0026 and RFC-0021.
2. Associate adapter requirements with DEC-0036.
3. Verify constitution references against the canonical constitution.
4. Update feature and decision traceability matrices.
5. Check referenced architecture files exist.
6. Include relevance review in human approval.

## 6.4 Duplicated Reuse Vocabulary

### Risk

Inconsistent terms such as adopt, leverage, customize, or integrate may obscure the controlled disposition.

### Mitigation

Use exactly:

```text
Reuse
Wrap
Extend
Build
```

Other words may appear in rationale but not as primary disposition values.

## 6.5 Unclear Approval State

### Risk

Kiro or Cursor may interpret a strongly worded candidate as approved.

### Mitigation

Every assessment must use one status:

```text
Proposed
Approved
Deferred
Rejected
Superseded
```

Only `Approved` with recorded human approval may authorize the next stage.

## 6.6 Missing Reassessment Lifecycle

### Risk

A decision may become stale without a defined review trigger.

### Mitigation

Require triggers such as:

- license change;
- governance or ownership change;
- end-of-life;
- loss of maintenance;
- security incident;
- sovereignty incompatibility;
- unsupported disconnected deployment;
- material API incompatibility;
- unacceptable operational cost;
- changed phase objective;
- changed provider requirements.

A triggered reassessment requires a new Architecture Decision Handoff before implementation changes.

## 6.7 Repository Risk Mitigation Matrix

| Repository risk | Preventive control | Detection control | Corrective action |
|---|---|---|---|
| Template divergence | Canonical versioned standard | Consistency check and review checklist | Update dependent templates |
| Superficial gate | Mandatory structured fields | Strict gate checks | Reject feature until complete |
| Incorrect traceability | Controlled DEC/RFC references | Existence and relevance review | Correct matrices and metadata |
| Vocabulary drift | Four controlled values | Reject unknown dispositions | Normalize terminology |
| Unclear approval | Mandatory status | Kiro precondition check | Return for human approval |
| Stale decisions | Reassessment triggers | Review during affected feature planning | Create new handoff |
| Silent architecture change | Approved handoff baseline | PR comparison | Stop merge and return to review |
| Runtime scope leakage | Explicit Phase 2 non-goals | Design review, tests, gate | Remove or defer execution behavior |
| Vendor-native leakage | Adapter and neutrality rules | Model review and contract tests | Replace with normalized contract |
| AI overreach | AI-consumer-only rule | Decision authority review | Remove mutation or execution ability |

---

# 7. Recommended Architecture Direction for FEATURE-0011

## 7.1 Architecture Decision

FEATURE-0011 is an **Architecture Operating System governance standard**, not a Sovrunn runtime resource or Go service.

Its purpose is to establish a mandatory, reviewable, and partially machine-verifiable decision contract for Phase 2 and later features.

The standard itself is classified as:

```text
Extend
```

It extends:

- DEC-0026 reuse before build;
- RFC-0021 reuse-first architecture;
- existing handoff and review templates;
- established architecture-decision practices;

with Sovrunn-specific requirements for:

- sovereignty;
- provider neutrality;
- adapter boundaries;
- capability-level classification;
- phase scope;
- feature-gate traceability;
- future-feature mitigation;
- reassessment lifecycle.

## 7.2 Assessment Granularity

Each assessment identifies its unit of analysis:

```text
feature
sub-capability
external integration
protocol or standard
data store
execution engine
provider integration
runtime controller
```

A feature-level summary is required, but significant components may have separate classifications.

## 7.3 Canonical Assessment Contract

```markdown
## Reuse Assessment

### Assessment scope
- Feature:
- Capability or decision unit:
- Assessment owner:
- Related DEC/RFC/ADH:

### Existing mature solutions or standards
| Candidate | Category | Relevant strengths | Material constraints |
|---|---|---|---|

### Decision
- Disposition: Reuse / Wrap / Extend / Build
- Selected foundation or approach:
- Decision status: Proposed / Approved / Deferred / Rejected / Superseded

### Rationale
- Why this disposition fits:
- Why rejected dispositions do not fit:

### Responsibility boundary
- Sovrunn-owned responsibility:
- Reused or external responsibility:
- Data crossing the boundary:
- Control crossing the boundary:

### Adapter boundary
- Required: Yes / No
- Adapter or contract:
- Reason:
- Vendor-native types allowed in Sovrunn core: No / Approved exception reference

### Fit considerations
- Sovereignty and deployment:
- Security and trust:
- Operations and supportability:
- Licensing and supply chain:
- Portability and replacement:
- Provider-neutrality impact:

### Phase impact
- Allowed in current phase: Yes / No
- Current-phase work:
- Deferred work:

### Non-goals
- ...

### Risks and mitigation
- Applicable architecture risks:
- Preventive controls:
- Detection controls:
- Corrective path:
- Residual risk:

### Replacement and reassessment
- Future replacement risk: Low / Medium / High
- Reassessment triggers:
- Exit or migration boundary:

### Traceability
- Architecture documents:
- Acceptance criteria:
- Validation or review evidence:
```

The assessment must be proportionate to architectural significance; it is not a general procurement study.

## 7.4 Disposition Semantics

### Reuse

Adopt a mature implementation, protocol, or standard substantially as provided.

Sovrunn does not own or fork its core behavior. An adapter may still be required.

### Wrap

Place a Sovrunn-owned contract around a mature capability for governance, neutrality, audit, normalized errors, tenant boundaries, and replaceability.

The wrapper must not recreate the wrapped engine.

### Extend

Add Sovrunn-specific behavior through supported extension, composition, or compatible augmentation.

A maintained fork requires explicit approval and a maintenance assessment.

### Build

Implement a Sovrunn-owned capability because it is a core differentiator or no mature solution provides acceptable fit.

A Build decision must explain why Reuse, Wrap, and Extend are insufficient and define long-term ownership.

## 7.5 Validation Model

### Level 1 — Structural validation

Verify:

- required sections;
- controlled disposition;
- decision status;
- adapter Yes/No;
- replacement risk;
- non-goals;
- mitigation fields;
- architecture references.

### Level 2 — Consistency validation

Where practical, detect:

- invalid feature IDs;
- unknown disposition values;
- missing Build rationale;
- missing adapter rationale;
- missing Phase 2 non-goals;
- inconsistent templates;
- nonexistent DEC/RFC references;
- missing mitigation controls;
- unapproved status before implementation.

Automation must not choose products or approve architecture.

### Level 3 — Human semantic review

The reviewer determines:

- credibility of mature options;
- correctness of ownership;
- sovereign-deployment fit;
- adequacy of adapters;
- justification for Build;
- mitigation quality;
- acceptability of residual risk;
- phase compliance.

## 7.6 Canonical-Source Rule

`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` is the canonical definition.

Templates and prompts must reference it or reproduce an explicitly versioned canonical section.

## 7.7 Change-Control Rule

A reuse assessment records reasoning but does not independently create an accepted decision.

A disposition becomes binding only through approved feature architecture, an approved Architecture Decision Handoff, or an accepted DEC/RFC where required.

Changing an approved disposition, mitigation plan, or responsibility boundary requires a new handoff.

## 7.8 FEATURE-0011 Non-Goals

FEATURE-0011 must not:

- select a production policy engine;
- select an identity or secret backend;
- select a workflow engine;
- select a PostgreSQL operator;
- define provider integrations;
- define plugin execution;
- create a runtime `ReuseAssessment` API resource;
- implement Go production code;
- build a general procurement platform;
- fully assess every later Phase 2 feature in advance.

---

## Approval Record

This architecture spine was approved by the Sovrunn project owner in the architecture review conversation on 2026-07-21.

Approval authorizes preparation of FEATURE-0011 Kiro requirements, design, and tasks through the required sequential approval process.

Approval does not authorize runtime implementation, FEATURE-0012 work, provider integration, or provisioning.
