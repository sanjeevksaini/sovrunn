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

## Capability assessment: deterministic policy-evaluation seam and fake

This assessment applies the canonical
`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` fields without redefining
them. It consolidates only the dispositions and boundaries already approved by
ADH-2026-067, ADH-2026-068, and ADH-2026-069.

### Identity

| Field | Value |
|---|---|
| Feature identity | FEATURE-0017 |
| Capability or decision-unit identity | Deterministic, engine-neutral policy-evaluation seam, pure FEATURE-0013 mapper, and in-process fake |
| Assessment owner | Sanjeev Kumar, Sovrunn Architecture Owner |

### Classification

| Field | Value |
|---|---|
| Disposition | Build |
| Decision status | Approved |

### Analysis

| Field | Value |
|---|---|
| Assessment scope | The FEATURE-0017 request/result contracts, evaluation boundary, adapter port, canonical input digest, pure FEATURE-0013 mapper, and deterministic fake only. |
| Candidate category | In-process policy-evaluation contract and replaceable engine boundary. |
| Mature candidates / applicable standards | RFC 8785 JCS, SHA-256, FEATURE-0012 reference/validation foundations, FEATURE-0013 evaluation evidence, and the DEC-0036 adapter-boundary principle. No production policy engine is selected. |
| Relevant candidate strengths | The reused standards and prior-feature contracts supply deterministic canonicalization, stable input identity, structural references, timing/evidence carriers, and replaceable integration boundaries. |
| Material candidate constraints | No mature external component owns Sovrunn's complete canonical request/result, validation, failure normalization, digest, transient evidence, and downstream-neutral mapping contract; Phase 2R also prohibits real external engine execution. |
| Rationale | Build only the Sovrunn-owned seam and deterministic conformance fake while reusing applicable standards and prior-feature contracts unchanged. |
| Selected foundation or approach | One in-process `PolicyEngineAdapter` port surrounded by the FEATURE-0017 boundary, RFC 8785/SHA-256 input identity, an injected UTC time source, a pure FEATURE-0013 mapper, and one immutable digest-keyed fake. |

### Boundary

| Field | Value |
|---|---|
| Sovrunn-owned responsibility | Validate and canonicalize the request, compute its digest, invoke one injected adapter at most once, normalize conclusion or failure, construct transient result/timing evidence, and perform the pure structural FEATURE-0013 mapping. |
| Reused or extended responsibility | FEATURE-0012 owns TypedRef validation/redaction foundations; FEATURE-0013 owns EvaluationResult and structural evidence carriers; FEATURE-0015 owns CloudPlatform/CloudProvider terminology; FEATURE-0016 supplies adapter-boundary precedent; RFC 8785 and SHA-256 remain external standards. |
| Responsibility/control boundary | FEATURE-0017 owns only the generic in-process evaluation seam and fake. Prior features retain their canonical contracts, while later domains retain IAM, governance, sovereignty, placement, approval, execution, explanation, adoption, publication, and orchestration authority. |
| Data crossing the boundary | A normalized request plus its precomputed digest enters the injected adapter; only a normalized outcome/reason-code conclusion or adapter-failure signal returns. Tenant-confidential request/result values remain in process. |
| Control crossing the boundary | The evaluation boundary invokes exactly one configured adapter only after cancellation, request-validation, canonicalization, and digest checks succeed; no adapter selects another adapter or performs a retry. |
| Adapter required | Yes |
| Adapter rationale | DEC-0036 requires the canonical Sovrunn contract to remain independent of any replaceable policy engine; the Phase 2R fake implements the same seam without external execution. |
| Adapter or contract identifier | `PolicyEngineAdapter` |
| Vendor-native types allowed | No |

### Suitability

| Field | Value |
|---|---|
| Sovereignty and deployment fit | Fully in-process and deterministic in Phase 2R; suitable for disconnected and air-gapped deployment because no network, credential, CloudProvider SDK, or external engine is used. |
| Security and trust | Request and result values are Tenant-confidential; fixed sanitized failures never echo request, digest, fixture, reason-code, or raw adapter-error values; no credential field or secret reference is introduced. |
| Operational and supportability | Fixed-time replay, frozen digest vectors, immutable fixtures, closed outcomes/non-results, exact precedence, race tests, and conformance tests provide repeatable local diagnosis. |
| Licensing and supply-chain | Production uses the Go standard library and existing Sovrunn contracts only; test-only YAML uses the repository's existing `gopkg.in/yaml.v3` dependency; no production policy-engine dependency is added. |
| Portability and provider-neutrality impact | The seam uses canonical CloudProvider terminology and opaque typed references; no provider-native or engine-native type becomes canonical. |

### Phase and scope

| Field | Value |
|---|---|
| Allowed in current phase | Yes |
| Current-phase work | Implement only the deterministic in-process request/result contracts, adapter port and boundary, canonical digest, injected time seam, pure mapper, immutable fake, and FEATURE-0017-local conformance proof. |
| Deferred work | Selecting or executing a real OPA, Cedar, or other engine; production adapter selection/registration; profile publication; external policy loading; credentials; persistence; routes; and later IAM, governance, sovereignty, placement, approval, execution, explanation, and orchestration semantics. |
| Explicit non-goals | No real policy engine, external call, credential, CloudProvider SDK, route, handler, controller, store, persistence, retry, plugin execution, provisioning, DecisionRecord/AuditEvent publication, or downstream domain-policy interpretation. |
| Exit or migration boundary | A separately approved later feature may implement a real engine adapter behind the unchanged `PolicyEngineAdapter` port. Any v1 semantic-input change requires an approved v2 digest contract; any ownership or canonical-contract change requires a new ADH. |
| Phase 2 non-goal acknowledgement | Phase 2R authorizes deterministic in-process evaluation and fake conformance only; real external policy-engine execution and infrastructure effects remain non-goals. |

### Build justification

| Field | Value |
|---|---|
| Why Reuse is insufficient | Existing standards and prior-feature contracts are reused, but none supplies the complete Sovrunn policy-evaluation seam, normalized failure model, transient evidence linkage, and exact Phase 2R boundary. |
| Why Wrap is insufficient | No production engine is selected or allowed in FEATURE-0017, so there is no approved external runtime to wrap. |
| Why Extend is insufficient | Extending FEATURE-0012, FEATURE-0013, or FEATURE-0016 would violate their ownership and couple generic policy evaluation to reference, decision-adoption, or ExecutionTarget lifecycle semantics. |
| Protected Sovrunn differentiation and long-term ownership | Sovrunn owns the provider-neutral request/result contract, adapter port, validation/canonicalization boundary, normalized outcome/failure evidence, pure mapper, and deterministic fake conformance behavior. |

### Risk mitigation

#### Applicable architecture risks

#### Risk-control matrix

| Risk | Preventive control | Detection control | Corrective path |
|---|---|---|---|
| Architecture leakage or premature engine selection | Closed six-group architecture; exact exclusions; engine-neutral adapter; deterministic fake only. | Requirements/design/tasks semantic review, Phase 2R scope checks, and production-import conformance. | Remove leaked semantics or dependencies; require an approved ADH before changing the engine or responsibility boundary. |
| Nondeterministic input identity or result timing | Frozen RFC 8785/SHA-256 vector, sorted semantic sets, injected UTC time, immutable fixtures, and exact precedence. | Digest, replay, concurrency, race, cancellation, and malformed-input tests. | Restore the approved v1 canonicalization/timing contract; require an approved v2 digest decision for semantic-input changes. |
| Tenant-confidential data or external-effect leakage | Sanitized fixed diagnostics, defensive copying, closed production import surface, no route/store/controller, and no credential or external-call surface. | Leakage sentinels, fixture-corpus review, import/declaration conformance, Phase 2R scope checks, and feature gate. | Remove the leaking value/effect, restore the closed import and diagnostic boundary, and reassess through an ADH if external behavior is required. |
| Prior-feature ownership drift | Pure structural FEATURE-0013 mapper, shared FEATURE-0012 validators, opaque candidate/profile references, and no adopting-profile validation or publication. | Mapper structural/parity tests, dependency-direction checks, spec semantic review, and traceability review. | Revert duplicated or downstream semantics to the owning feature and require approval for any dependency extension. |

- Residual risk: Low. The deterministic fake proves the seam but cannot prove
  the behavior, performance, or operational suitability of an unselected real
  policy engine.
- Replacement risk: Medium
- Reassessment triggers: selection of a real engine or new runtime dependency;
  a v1 request/digest or result vocabulary change; obligations becoming active;
  external calls, credentials, persistence, or routes; a mapper ownership
  change; or a material licensing, maintenance, security, sovereignty, or
  deployment-context change.

### Traceability

| Field | Value |
|---|---|
| Related DEC / RFC / ADH references | DEC-0026, DEC-0028, DEC-0036, DEC-0043, RFC-0021, RFC-0025, ADH-2026-067, ADH-2026-068, ADH-2026-069 |
| Linked acceptance criteria | AC-F17-01 through AC-F17-24 |
| Validation and review evidence | Approved ADH-2026-067/068/069; recorded FEATURE-0017 architecture and executable-plan gates; independent requirements/design/tasks reviews; FEATURE-0017 semantic, boundary, Phase 2R, conformance, deterministic, race, and repository verification evidence. |

### Human-approval evidence

| Field | Value |
|---|---|
| Structured approval-evidence record | `docs/reviews/reuse-assessments/FEATURE-0017-approval-evidence.md` |
| Approval status | Approved |
| Approving person or role | Sanjeev Kumar, Sovrunn Architecture Owner |
| Approval date | 2026-08-25 |
| Approval basis | ADH-2026-067, ADH-2026-068, ADH-2026-069, and the recorded FEATURE-0017 architecture-gate approval apply to the listed dispositions and responsibility boundary. |

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
