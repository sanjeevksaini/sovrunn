---
doc_type: feature
id: FEATURE-0017
title: Policy Evaluation Abstraction
status: architecture-approved
phase: 2R
reuse_assessment_format_version: 1.0.0
canonical_architecture: docs/architecture/policy-evaluation-abstraction.md
kiro_slug: policy-evaluation-abstraction
controlling_handoff: ADH-2026-067
clarifying_handoff: ADH-2026-068
corrective_handoff: ADH-2026-069
ai_load_priority: feature
ai_summary: Architecture-approved scope and generation-control contract for the deterministic, in-process, engine-neutral FEATURE-0017 evaluation seam and fake.
---

# FEATURE-0017: Policy Evaluation Abstraction

| Field | Value |
|---|---|
| Status | Architecture approved; feature-factory architecture gate recorded |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase / order | Phase 2R / 7 |
| Direct dependencies | FEATURE-0013 and FEATURE-0016 |
| Reused prior authority | FEATURE-0011, FEATURE-0012, and FEATURE-0015 |
| Depended on by | FEATURE-0018 directly; FEATURE-0020, FEATURE-0023, FEATURE-0025, and FEATURE-0026 through later owning domains |
| Sole architecture authority | `docs/architecture/policy-evaluation-abstraction.md` |
| Controlling handoffs | ADH-2026-067, ADH-2026-068, and ADH-2026-069 |
| Public routes | None |
| Persistence / controller | None |
| Real policy engine | None in Phase 2R |

## 1. Purpose and authority boundary

FEATURE-0017 supplies a deterministic, in-process, engine-neutral evaluation
seam so Sovrunn can validate and identify policy input, invoke one replaceable
adapter, normalize its conclusion or failure, and return transient audit-ready
evidence without embedding a real policy engine.

This file is the feature scope and generation-control contract. It does not
restate or reinterpret FEATURE-0017 architecture. The six closed decision
groups, all field semantics, ordering, mapping behavior, conformance cases,
and deferrals are owned only by
`docs/architecture/policy-evaluation-abstraction.md`. Any conflict stops work
with `ARCHITECTURE_DECISION_REQUIRED`; Kiro must not repair architecture in a
spec artifact.

## 2. Owned scope

FEATURE-0017 owns realization of `PolicyEvaluationRequest` and
`PolicyEvaluationResult`, the engine-neutral `PolicyEngineAdapter` port, the
evaluation boundary, a deterministic in-process fake, structural validation,
RFC 8785 JCS canonicalization, SHA-256 input identity, normalized
result-versus-failure handling, transient evidence return, and the separate
pure structural mapper to FEATURE-0013 `EvaluationResult`.

The exact owned behavior is the canonical architecture's CDG-F17-01 through
CDG-F17-06. Requirements may translate that behavior into verifiable
statements and acceptance criteria; they may not add another decision group,
policy semantics, or downstream adoption rule.

## 3. Dependency and reuse boundaries

- FEATURE-0011 supplies the mandatory reuse-before-build governance format.
- FEATURE-0012 supplies common immutable-reference, transient-contract,
  validation, and redaction foundations. FEATURE-0017 does not alter them.
- FEATURE-0013 owns `DecisionRecord`, `DecisionProfile`, `EvaluationResult`,
  and `AuditEvent`. FEATURE-0017 supplies transient evidence and a pure
  structural mapping only; it does not inspect a `DecisionProfile`, validate
  adopting-profile compatibility, publish, or persist.
- FEATURE-0015 supplies canonical CloudPlatform and CloudProvider scope
  terminology. FEATURE-0017 does not read or mutate FEATURE-0015 resources.
- FEATURE-0016 supplies adapter-boundary precedent and owns ExecutionTarget
  qualification, normalized target facts, lifecycle, backing access, and
  stores. Candidate references remain opaque to FEATURE-0017.
- FEATURE-0018 and later domains may call or adopt FEATURE-0017 results but
  remain the owners of IAM, governance, sovereignty, placement, approval,
  execution, explanation, and orchestration semantics.

## 4. Feature-level reuse summary

This is the required FEATURE-0017 feature-level reuse summary. It mirrors the
approved disposition in the sole architecture authority and creates no
independent architecture authority.

| Capability | Disposition | Decision status | Rationale | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 `TypedRef`, transient contract, validation, and redaction foundations | Reuse | Approved | Existing canonical foundations satisfy the common contract need. | DEC-0026 and ADH-2026-012 |
| FEATURE-0013 `EvaluationResult` | Reuse | Approved | Decision evidence remains owned by FEATURE-0013. | ADH-2026-017 and DEC-0043 |
| RFC 8785 JCS | Reuse | Approved | A stable standard supplies deterministic JSON canonicalization without a Sovrunn-specific algorithm. | ADH-2026-067 |
| SHA-256 | Reuse | Approved | A standard digest supplies deterministic v1 input identity. | ADH-2026-067 |
| Sovrunn `PolicyEngineAdapter` port | Build | Approved | No external component owns the canonical Sovrunn policy-evaluation contract. | DEC-0028, DEC-0036, and ADH-2026-067 |
| Evaluation boundary and pure FEATURE-0013 mapper | Build | Approved | Sovrunn owns validation, normalization, evidence construction, success timing transport, and structural linkage while engine behavior remains outside core. | ADH-2026-067, ADH-2026-068, and ADH-2026-069 |
| Deterministic in-process fake | Build | Approved | Phase 2R needs executable conformance without selecting or embedding a real policy engine. | Phase 2R scope and ADH-2026-067 |
| Injected deterministic UTC time-source pattern | Reuse | Approved | In-process dependency injection supplies deterministic timing without a new platform service. | ADH-2026-067 |
| Real OPA, Cedar, or other engine adapter | Wrap | Deferred | A later feature requires a fresh candidate assessment; no real engine is selected or allowed here. | DEC-0028 and DEC-0036 |

For the Build units, Reuse is insufficient because no external component owns
the Sovrunn canonical seam. Wrap is not applicable because FEATURE-0017
selects no real engine. Extend is not applicable because there is no approved
external implementation to extend. This disposition does not authorize a
custom policy engine.

## 5. Requirements-generation boundary

Requirements must derive only observable behavior and FEATURE-0017-local proof
from the sole architecture authority. In particular, they must preserve:

- deterministic, in-process, side-effect-free Phase 2R behavior;
- at most one fake-adapter invocation, and exactly one only after both
  pre-invocation cancellation/deadline checks, request validation, and digest
  calculation succeed;
- the closed result-versus-non-result distinction;
- optional generic `contextRef` semantics without resolving the context;
- opaque profile and candidate references without domain interpretation;
- the structural, pure, non-persisting FEATURE-0013 mapping boundary;
- successful transient transport of the complete `PolicyEvaluationResult`
  together with its exact boundary-captured FEATURE-0013 timing envelope;
- exact mapper population of `evaluator`, caller-supplied `inputSnapshotRef`
  and/or `inputIntegrity`, `resultStatus`, complete `result`,
  `timing.startedAt`, `timing.completedAt`, and `evaluatedAt`;
- mapper omission of `timing.durationMs`, `executionConfig`, `trustBoundary`,
  `safetyPolicyFilters`, `structuralTrustState`, and `scopeEvidence`;
- all canonical architecture precedence and conformance obligations; and
- zero public route, store, controller, production selector, retry, external
  call, credential, or real engine.

Requirements must not choose Go packages, concrete types, method signatures,
fixture serialization, goroutine mechanics, or production topology. Examples
in the architecture are illustrative and cannot become required fixtures or
business rules by inference.

### 5.1 Canonical requirement generation ledger

The rows below are subordinate generation-control mirrors of the six canonical
architecture decision groups. They give the feature factory stable requirement
identifiers but create no independent architecture authority. The exact
observable behavior remains owned solely by the referenced CDG in
`docs/architecture/policy-evaluation-abstraction.md`.

| Requirement ID | Architecture group | Exact title |
|---|---|---|
| REQ-F17-01 | CDG-F17-01 | Minimal request boundary |
| REQ-F17-02 | CDG-F17-02 | Exact v1 canonicalization and digest |
| REQ-F17-03 | CDG-F17-03 | Adapter, result, timing, and failure normalization |
| REQ-F17-04 | CDG-F17-04 | Deterministic fake |
| REQ-F17-05 | CDG-F17-05 | Transient FEATURE-0013 evidence linkage |
| REQ-F17-06 | CDG-F17-06 | In-process operation, precedence, proof, and governance |

### 5.2 Canonical acceptance generation ledger

The rows below are subordinate generation-control mirrors of the canonical
architecture's 24-item local conformance inventory. Their order and text are
copied verbatim. They create no independent acceptance or architecture
authority and may not be repurposed, renumbered, merged, split, or paraphrased.

| Acceptance ID | Exact acceptance item |
|---|---|
| AC-F17-01 | valid request and every registered bound; |
| AC-F17-02 | unknown-field and malformed-reference rejection; |
| AC-F17-03 | invalid request invokes no adapter; |
| AC-F17-04 | optional absent and valid-present `contextRef`; |
| AC-F17-05 | `null`, empty, and partial field rules; |
| AC-F17-06 | duplicate reference rejection; |
| AC-F17-07 | list-order-independent digest; |
| AC-F17-08 | semantic input change alters digest; |
| AC-F17-09 | `requestId`-only change preserves digest; |
| AC-F17-10 | exact RFC 8785/SHA-256 test vector; |
| AC-F17-11 | all four valid outcomes; |
| AC-F17-12 | missing fixture and configured adapter failure; |
| AC-F17-13 | invalid fixture construction failure; |
| AC-F17-14 | cancellation, deadline, zero-invocation-before-call, at-most-once invocation, and late-output discard; |
| AC-F17-15 | invalid adapter conclusion rejection; |
| AC-F17-16 | fixed-time-source deterministic replay; |
| AC-F17-17 | reason-code grammar, duplicate rejection, sorting, and redaction; |
| AC-F17-18 | absence of obligations in the Phase 2R fake; |
| AC-F17-19 | successful result/timing transport and separate pure FEATURE-0013 mapping for all four results, with the exact populated fields and exact omitted fields required by CDG-F17-05; |
| AC-F17-20 | malformed structural mapper-input failure with no mutation or side effect; |
| AC-F17-21 | no mapping for non-result interactions; |
| AC-F17-22 | no FEATURE-0017 DecisionRecord, AuditEvent, route, controller, or store; |
| AC-F17-23 | concurrent equivalent-call determinism; and |
| AC-F17-24 | zero network, OPA, Cedar, Kubernetes, CloudProvider, database, identity, provisioning, or other external effects. |

## 6. Acceptance boundary

The complete feature-level acceptance inventory is the 24-case local
conformance list in Section 10 of the sole architecture authority and its
subordinate generation-control mirror in Section 5.2 above. Later Kiro
requirements must make every case traceable and testable without broadening
it. The FEATURE-0017 contract checker must validate schemas, the six decision
groups, digest vector, adapter/fake boundaries, conformance inventory, and
zero-route/store/controller guarantees; it must not require a route catalog.

The feature is not implementation-complete until the repository feature gate
passes. Architecture-gate approval authorizes requirements generation only;
it does not authorize design, tasks, Cursor execution, or implementation.

## 7. Explicit non-goals

- authentication, PrincipalRef creation, membership, role, or authorization
  resolution;
- governance, sovereignty, assignment, inheritance, exception, approval,
  entitlement, quota, placement, ranking, selection, or execution semantics;
- policy/profile publication, profile-kind taxonomy, bundle distribution,
  loading, or hot reload;
- `EffectiveGovernanceContext` construction or lookup;
- DecisionRecord or AuditEvent publication and any request/result store;
- customer-safe or CloudProvider-safe reason/explanation projection;
- public API routes, handlers, controllers, idempotency repositories, or
  persistent storage;
- production adapter selection, registry, routing, retry, or topology; and
- real OPA, Cedar, Kubernetes, identity, database, network, CloudProvider SDK,
  provisioning, or plugin execution.

## 8. Generation and execution controls

The machine-readable controls are
`.automation/features/FEATURE-0017.control.json` and
`.automation/state/FEATURE-0017.json`. The only permitted Kiro outputs are, in
order, `requirements.md`, `design.md`, and `tasks.md` within
`.kiro/specs/policy-evaluation-abstraction/`, one stage at a time after its
required gate. Cursor may act only from an approved complete `tasks.md` and
must not change architecture.

The generic feature-factory architecture gate pins this feature file, the sole
architecture, ADH-2026-067, ADH-2026-068, ADH-2026-069, and the control
manifest by digest. Any later change to that package invalidates the recorded
approval and requires a fresh explicit architecture-gate approval before
requirements generation resumes.
