---
doc_type: architecture_readiness_dossier
feature: FEATURE-0016
title: Adapter Boundary and ExecutionTarget Qualification
status: planning-input-not-approved-architecture
baseline: ARCH-2026.08-PHASE2R-CANONICAL
prepared: 2026-08-14
---

# FEATURE-0016 Architecture Starting Dossier

## 1. Purpose and authority boundary

This dossier gathers the current repository authority for FEATURE-0016 and
records the architecture decisions that must be closed before requirements are
generated. It is a planning and review artifact only: it neither changes an
approved contract nor authorizes requirements, design, tasks, or source work.

The governing source order is the current Phase 2R baseline, canonical data
model, canonical contract catalog, accepted DEC files, the VS-000 registry,
and then feature-sequence/architecture documents. A later generated artifact
cannot resolve an ambiguity in those sources.

## 2. Approved starting position

| Area | Already approved and usable as architecture input | Authority |
|---|---|---|
| Feature sequence | FEATURE-0016 follows FEATURE-0015 and owns adapter boundary plus ExecutionTarget qualification. | `docs/phase2/PHASE2_FEATURE_SEQUENCE.md` |
| Feature exit evidence | A deterministic fake target qualifies, another is denied, and no implementation-native fields enter canonical or customer contracts. | `docs/phase2/PHASE2R_REBASELINE.md` §3 |
| Resource ownership | FEATURE-0016 fully owns ExecutionTarget identity, schema, routes, status, writers, conformance, qualification, availability, maintenance epoch, NormalizedTargetFactSet, TargetQualificationResult, and target lifecycle. | `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md` §4 |
| Canonical resource boundary | ExecutionTarget is an adapter-addressable realization boundary. It may be infrastructure-backed, an external managed-service API, an edge control plane, or a future federated endpoint. It is not a mandatory child of topology, not customer-visible, and not a capacity scheduler. | `docs/architecture/canonical/sovrunn-finalized-data-model.md` §10 |
| No retired abstractions | No mandatory ResourcePool or CloudProvider-wide capability inventory exists in core. | DEC-0042 |
| Adapter rule | Core depends on Sovrunn contracts; replaceable external systems sit behind adapters. Phase 2R adds an execution-target qualification adapter. | DEC-0036; `docs/architecture/adapter-boundary-model.md` |
| Registered contracts | VS0-SCHEMA-015 ExecutionTarget, VS0-SCHEMA-016 NormalizedTargetFactSet, VS0-SCHEMA-017 TargetQualificationResult; CloudProvider scope; confidential/never-customer-visible or omitted as registered. | `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` |
| Writers | CloudProvider administrator owns target spec. The fake adapter and qualification controller own fact/result output. The target lifecycle controller owns target status. | VS0-WRITER-006 and VS0-SCHEMA-015 |
| Target states | Initial Unqualified+Unavailable; qualification and availability are orthogonal; stale epoch fails with STALE_RESOURCE_VERSION/412 and VS0_TARGET_EPOCH_STALE. | VS0-STATE-004; VS0-CF-F10 |
| Maintenance rule | CloudProvider/trusted adapter announces native maintenance; only the target lifecycle controller changes availability; epoch/fence blocks stale downstream work and requalification is mandatory before availability returns. | DEC-0057; canonical model §10.2 |
| Phase boundary | Phase 2R remains deterministic, in-memory, and side-effect free: no real cloud, Kubernetes, OpenShift, PostgreSQL, external secret system, or CloudProvider API call. | `docs/phase2/PHASE2R_REBASELINE.md` §§5–6 |

## 3. Feature boundary

### Owned by FEATURE-0016

- ExecutionTarget registration and its internal/operator-facing projection.
- Adapter-facing normalized target facts and immutable qualification results.
- Qualification/availability state transition semantics, maintenance epoch, and
  target fencing.
- A deterministic fake qualification adapter, including the target-scoped
  SecretRef boundary without raw credential values.
- FEATURE-0016-local contract, negative, concurrency, redaction, and
  no-external-effect proof.

### Consumed by reference only

- FEATURE-0012 metadata, typed references, Problem Details, request grammar,
  common validation, and optimistic-concurrency conventions.
- FEATURE-0013 AuditEvent and DecisionRecord contracts.
- FEATURE-0015 CloudProviderParticipation, InfrastructureStack, and topology
  facts. They are prerequisites, not editable or re-owned by F0016.

### Explicitly deferred or excluded

- ResourcePool, a CloudProvider-wide capability inventory, a core capacity
  scheduler, or the retired generic supply term as an active concept.
- Customer-facing target/topology projection; FEATURE-0023 owns the safe
  ServicePlacement projection.
- ServiceRegion and `spec.executionTargetRefs`; FEATURE-0022 owns them.
- Policy evaluation, IAM/RoleAssignment, CloudEnrollment, entitlement, quota,
  sovereignty evidence, placement decisions, PluginExecution, and actual
  provisioning.
- Real adapter/CloudProvider calls, raw credentials, protected native handles in
  customer/canonical contracts, and alpha runtime migration.

## 4. Architecture decisions that are still required

The following are executable-contract gaps, not invitations for requirements,
design, or tasks to choose semantics. They must be resolved in one approved,
bounded ADH before requirements generation.

| ID | Gap that remains after reading current authority | Required approved outcome | Why it blocks requirements |
|---|---|---|---|
| ARC-F16-01 | ExecutionTarget's registered schema has required spec/status fields, but no exact F0016 method/path table, create/PATCH surface, request body, scope derivation, initial status, or local route conformance ledger. | Specify the complete F0016 public/internal API surface, each route's request/headers/result/problem/status/audit effect, and its Go 1.22 feasible registration form. | Requirements cannot translate an unknown API surface or choose client versus server ownership. |
| ARC-F16-02 | `spec.infrastructureStackRef` is required in VS0-SCHEMA-015, while the canonical data model says an ExecutionTarget need not reference InfrastructureStack for external API/edge/federated realization. | Reconcile required/optional backing-environment semantics, reference names, and the allowed realization-pattern contract. | The current authorities permit mutually incompatible valid target shapes. |
| ARC-F16-03 | DEC-0036 requires an adapter boundary, but the adapter contract has no exact operation set, input/output envelope, timeout/cancellation/failure classification, fact freshness rule, fake adapter behavior, or provenance/redaction rule. | Define the minimal F0016 qualification-adapter contract and deterministic fake behavior without choosing a vendor or performing external calls. | Design otherwise invents adapter semantics and leaks implementation-native fields. |
| ARC-F16-04 | Schema 015 says target-lifecycle-controller writes status and writer 006 names fake-adapter-and-qualification-controller for facts/results, but exact status-field ownership and handoff between fact observation, qualification result, and target publication are not enumerated. | Assign every F0016 status/record field to exactly one writer, define initial values and publication/audit behavior, and prohibit direct client status writes. | It prevents split/competing writers and an untestable status lifecycle. |
| ARC-F16-05 | VS0-STATE-004 has a compact six-state form, while canonical maintenance authority describes DrainPending, Draining, InMaintenance, Requalifying, plus Unavailable/Degraded/Restricted branches. InfrastructureMaintenanceNotice has no visible VS-000 schema/owner/conformance entry. | Reconcile the state vocabulary; decide what F0016 actively models now; and, if notice input is F0016 scope, register its schema, writer, routes, audit, retention, and proof. Otherwise declare its owner and F0016's bounded input contract. | Requirements cannot safely choose maintenance states, a notice resource, or an implicit controller. |
| ARC-F16-06 | F10 proves only changed-epoch stale work rejection. No F0016-local cases cover registration, fact freshness/expiry, qualified/rejected/indeterminate result, authorized versus inaccessible CloudProvider references, writer denial, audit failure, or zero external effect. | Add a complete F0016-owned conformance matrix with exact inputs, output state, Problem code/HTTP/violation, side effects, and gate for every observable requirement. | F0015 showed that downstream/shared cases cannot be repurposed as feature-local proof. |
| ARC-F16-07 | The target-scoped SecretRef boundary is stated in the phase plan, but the current F0016 schemas do not identify the allowed reference carrier, writers/readers, rotation behavior, logging/audit redaction, or denial outcome. | Define a reference-only credential/authority boundary using FEATURE-0008/0012 conventions, or explicitly defer a field while preserving fake-adapter operation. | A requirements author would otherwise invent secret storage or omit an essential adapter boundary. |
| ARC-F16-08 | F0016 depends on F0015 but shares cross-CloudProvider target references. Existing X03 is owned by F0015, while F0016 needs exact safe-denial and authorized-mismatch semantics for target/participation/stack resolution. | State whether F0016 consumes X03 unchanged or receives a local case; enumerate safe-resolution order and no-existence-disclosure/audit effects. | It avoids a false claim of F0016-local proof and prevents information disclosure. |
| ARC-F16-09 | Qualification, lifecycle publication, audit append failure, retry/replay behavior, and maintenance-epoch concurrency have no feature-local atomicity/replay matrix. | Reuse the prior audit/idempotency contract by exact reference where applicable and explicitly state which F0016 outcomes are replayable, abort-only, audited, or non-audited. | Design must not choose durable/temporary storage or observable retry/error semantics. |
| ARC-F16-10 | The current baseline/phase-context “next feature” pointers still name F0015, even though F0015 is merged and the approved sequence now selects F0016. | Correct the stale status pointers as an administrative follow-up in the approved F0016 authority package; do not treat the stale text as feature scope authority. | It prevents the wrong feature context from being loaded into Kiro prompts and reviews. |

## 5. Required F0016 architecture admission artifacts

Before a requirements prompt is generated, the approved package must contain:

1. a new FEATURE-0016 architecture authority and a feature-level reuse summary;
2. a closed registry/schema/writer/state/route/conformance update for every
   approved F0016 observable behavior;
3. a single table that classifies each field as client-provided,
   server-derived, adapter-observed, controller-owned, or deferred;
4. exact reference-resolution and authorization ordering, including safe
   denial before any existence-revealing structural validation;
5. an adapter contract with a deterministic fake and a no-external-effect
   test oracle;
6. an audit matrix and a target/epoch concurrency matrix;
7. a resource/route/requirement/acceptance/conformance proof matrix with no
   downstream ID counted as F0016-local proof; and
8. deterministic readiness checks plus small formal models for target-state
   safety and stale-epoch fencing, before Kiro reads the requirements prompt.

## 6. FEATURE-0015 lessons converted into hard gates

The objective is not to trust a reviewer to remember previous gaps. Each item
below must be mechanically checked before a stage can advance.

| F0015 learning | F0016 preventive gate |
|---|---|
| A required status field without a single named writer created ambiguity. | Every F0016 status/record field has one initial value, writer, transition source, audit rule, and conformance case. |
| Route prose was not enough; one action URI was impossible for Go ServeMux. | Route table is validated against exact method/path registrations and `PathValue` use before requirements generation. |
| Field ownership, create body, scope derivation, and success status were discovered too late. | Each route has a closed input classification table and result/error/audit declaration. |
| Generic idempotency language let different concrete targets share a replay namespace. | Any replayable F0016 operation defines principal, registered pattern, concrete target/scope, canonical digest, authorization-before-replay, retention, eviction, abort, and response-header rules. |
| Audit append failure and denial audit policy were under-specified. | An outcome matrix states mutation, state, idempotency, AuditEvent, append-failure, and public error effects for success, denial, validation, concurrency, panic, and shutdown. |
| A conformance ID was repeatedly mapped beyond its actual input semantics. | The checker validates each REQ/AC-to-conformance mapping by owner and exact case input; shared/downstream IDs are reference-only. |
| Requirements leaked design choices such as in-memory coordination. | Requirements state observable phase constraints only; data structures, locks, and internal scheduling stay in design unless they change an observable outcome. |
| Architecture terms appeared in old/stale forms outside explicit exclusion context. | The architecture authority contains an explicit active-vs-exclusion vocabulary check, with F0016's prohibited terms rejected as active behavior. |
| Implementation tasks omitted dependent writable paths and caused repeated Cursor blocks. | Before tasks approval, generate a deterministic implementation path inventory from every task's code/test/dependency surface and validate it against the Cursor allowlist. |
| Formal reasoning was added after concurrency semantics were already unstable. | Run target lifecycle and stale-epoch fence TLC models as part of F0016 architecture readiness, before Kiro requirements generation. |
| Context budgets were exceeded after authorities grew. | Manifest context is measured after the approved architecture package is assembled; trim by authority classification, not by omitting required contracts. |

## 7. Stage leak-prevention contract

| Stage | May do | Must not do | Deterministic admission rule |
|---|---|---|---|
| Architecture | Resolve observable product/API/authorization/state/error/audit/conformance choices through approved authority. | Leave an observable route, writer, error, lifecycle, or adapter behavior implicit. | All ARC-F16 rows resolved or explicitly deferred to a named later owner; readiness checker and formal models pass. |
| Requirements | Translate the approved authority verbatim into REQ/AC and proof mappings. | Select API mechanics, storage, adapter internals, or future-feature semantics. | Canonical ledgers and one-to-one local-proof mappings validate; every referenced authority digest matches. |
| Design | Choose internal code mechanics that preserve approved observable behavior. | Change routes, contract fields, writers, error precedence, lifecycle, audit, or conformance scope. | Executability audit proves every approved route and outcome has an implementation path and test plan. |
| Tasks | Decompose the approved design into dependency-ordered, path-complete work. | Invent missing behavior, silently broaden writable paths, or use shared/downstream IDs as local acceptance proof. | Task graph is acyclic; all writable/test/lifecycle paths and proof cases validate before Cursor. |

## 8. Recommended next architecture action

Prepare one bounded ADH titled **“ADH-2026-058 — FEATURE-0016 adapter and
ExecutionTarget executable-contract closure.”** It should resolve only
ARC-F16-01 through ARC-F16-10, list exact impacted authorities and checks, and
contain no vendor selection, real external integration, or later-feature
activation. Founder approval is required before Kiro applies it.

## 9. Candidate ExecutionTarget taxonomy for architecture discussion

This is a planning inventory, not an approved enum, API field, or activation
decision. A candidate becomes active only when an approved ADH assigns its
realization pattern, normalized fact vocabulary, adapter contract, projection,
and conformance evidence.

| Candidate target category | What the target represents | Phase 2R position | Native details that must remain behind the adapter |
|---|---|---|---|
| **Synthetic qualification target** | Deterministic in-process fake used to prove qualification, rejection, stale facts, fencing, and redaction. | **F0016 required test implementation.** It never represents a deployable environment or performs an external call. | All CloudProvider/cluster/account details. |
| **Infrastructure-backed IaaS target** | A CloudProvider-approved IaaS realization boundary with normalized VM compute, block-storage, object-storage, and private-network capability. | **Recommended first production-shaped fake profile.** No VM, volume, bucket, or network is created in Phase 2R. | Account/project IDs, VPC/subnet/tenant IDs, image IDs, security-group rules, quota inventory, and credentials. |
| **Private-cloud infrastructure target** | A target backed by an OpenStack-, VMware-, or similar privately operated infrastructure stack. | Deferred real adapter; the canonical shape must remain identical to the IaaS target. | Native API objects, regions, host aggregates, datastore names, and credentials. |
| **Container-platform target** | A CloudProvider-approved Kubernetes/OpenShift-like realization boundary for later service plugins. | Deferred real adapter and execution. It is not a customer-visible cluster contract. | Cluster/namespace IDs, kubeconfig, manifests, CRDs, nodes, and network policies. |
| **External managed-service API target** | An approved managed-service control-plane endpoint that can realize a service without an InfrastructureStack backing environment. | Canonically allowed; deferred until its optional-backing-reference and adapter rules are approved. | CloudProvider account, endpoint implementation, API tokens, service IDs, and implementation error detail. |
| **Virtualization or bare-metal target** | A qualified boundary exposed by an existing virtualized or bare-metal substrate. | Deferred; may reuse the same IaaS capability vocabulary where sufficient. | Hypervisor/host identifiers, VLANs, inventory, and management credentials. |
| **Edge or disconnected-site target** | A qualified local/edge control-plane boundary with constrained connectivity and sovereignty facts. | Deferred; requires explicit offline/freshness/recovery rules before activation. | Site appliance identity, local network topology, credentials, and operational handles. |
| **Remote or federated target** | A future independently operated remote control plane or federation endpoint. | Future-only. FEATURE-0016 must not introduce federation trust, replication, or remote desired-state authority. | Remote identity, offers, synchronization protocol, credentials, and native resource identities. |

### 9.1 What is never itself an ExecutionTarget

- An individual VM, volume, object bucket, VPC, subnet, security group, or
  Kubernetes namespace.
- A CloudProvider account, an IaaS management plane, or an adapter connection.
- A customer `ServiceInstance`, placement decision, Operation, or PluginExecution.
- A generic ResourcePool, CloudProvider-wide capability inventory, or capacity
  scheduler.

### 9.2 Admission rules for every future target category

Every category must prove all of the following before it is activated:

1. a CloudProvider/participation authorization boundary and safe reference
   resolution;
2. target-scoped reference-only credentials or authority, with no raw secret
   value in a canonical/customer contract, log, or audit event;
3. normalized, versioned, fresh, provenance-bearing facts rather than native
   implementation objects;
4. deterministic Qualified, Rejected, or Indeterminate behavior and a
   fail-closed stale/maintenance epoch fence;
5. no policy, placement, customer-projection, or execution authority inside
   the adapter; and
6. a reusable adapter conformance suite that proves redaction, cancellation,
   safe denial, and zero unintended external effect.

### 9.3 Canonical permanence promotion rule

This proposed F0016 architecture gate prevents an incidental implementation
layer from becoming a permanent Sovrunn platform concept.

A new resource, controller, adapter subtype, registry entry, public route, or
durably retained record may be promoted into the canonical model only when it
has at least one distinct and enduring responsibility that cannot be expressed
by an existing contract:

1. its own lifecycle and terminal/retention rule;
2. its own writer/reader authorization boundary;
3. a stable cross-feature identity or typed-reference contract;
4. a required customer-safe or operator-safe projection; or
5. an independently auditable external effect or decision boundary.

If none applies, the element must remain one of: a field on an existing
resource, a derived view, a transient request/result value, an
adapter-internal opaque handle, an implementation detail, or a test fixture.
It must not receive a canonical schema ID, public route, persistent store, or
separate feature owner merely for implementation convenience.

Every proposed promotion must state its owner, exit/retirement boundary,
replacement risk, and the specific invariant that an existing contract cannot
express. The approved ADH must reject a proposal that only duplicates an
existing lifecycle, capability, authority, or projection.

## 10. Exit condition for this dossier

FEATURE-0016 may enter requirements generation only after the approved ADH is
applied, the F0016 readiness checker and both formal models pass, the context
budget is measured, and a human architecture approval records the resulting
authority digest.
