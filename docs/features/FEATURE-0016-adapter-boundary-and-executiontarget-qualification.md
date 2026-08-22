---
doc_type: feature
id: FEATURE-0016
title: Adapter Boundary and ExecutionTarget Qualification
status: approved
phase: 2R
reuse_assessment_format_version: 1.0.0
canonical_architecture: docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md
kiro_slug: adapter-boundary-and-executiontarget-qualification
---

# FEATURE-0016: Adapter Boundary and ExecutionTarget Qualification

| Field | Value |
|-------|-------|
| Status | Approved scope (closed under ADH-2026-058) |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase | 2R |
| Order | 6 |
| Depends On | FEATURE-0011, FEATURE-0012, FEATURE-0013, FEATURE-0015 (read-only prior-feature authority) |
| Depended On By | FEATURE-0017, FEATURE-0019, FEATURE-0022, FEATURE-0023, FEATURE-0024 (future consumers; no activation or modification) |
| Architecture Boundary | docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md |
| Controlling Decisions | DEC-0036, DEC-0042, DEC-0057 |
| Controlling Handoffs | ADH-2026-025, ADH-2026-040, ADH-2026-042, ADH-2026-045, ADH-2026-058, ADH-2026-060, ADH-2026-061, ADH-2026-063, ADH-2026-064, ADH-2026-065, ADH-2026-066 |

---

## 1. Feature Summary

FEATURE-0016 activates one deterministic, infrastructure-backed
`synthetic-iaas` ExecutionTarget qualification profile. It establishes a
Sovrunn-owned observation-and-qualification boundary, with no real external
calls, credentials, placement, provisioning, adapter selection, or customer
projection. ADH-2026-058 replaces the FEATURE-0016 placeholder registry
definitions `VS0-SCHEMA-015..017` and `VS0-STATE-004` with one complete
executable contract.

---

## Feature-level reuse summary

This summary follows the canonical
`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` format. It records approved
dispositions only; it neither selects an external product nor changes the
approved FEATURE-0016 boundary.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 API grammar, Problem Details, media, and optimistic-concurrency contract | Reuse | These platform-wide contracts already provide the closed HTTP, error, ETag, and representation behavior required by F0016. | Approved | DEC-0026; ADH-2026-058 |
| FEATURE-0013 AuditEvent atomic-publication contract | Reuse | F0016 needs durable redacted audit evidence before publication, not a second audit model. | Approved | DEC-0026; ADH-2026-058 |
| FEATURE-0015 CloudProviderParticipation and InfrastructureStack backing authority | Reuse | F0016 consumes backing state without re-owning the resources, routes, or lifecycle. | Approved | DEC-0042; ADH-2026-058 |
| F0016-owned coherent backing access over FEATURE-0015 authority | Extend | A private, read-only store-backed lease supplies principal-aware safe access and immutable paired backing inputs for F0016 qualification/publication; it changes no F0015 public behavior or ownership. | Approved | ADH-2026-066 |
| Deterministic target observation, qualification, lifecycle, and safe projection | Build | No existing candidate supplies Sovrunn's closed target model, fenced qualification, sole-committer, audit, and Phase 2R no-external-effect contract. | Approved | DEC-0036; DEC-0042; ADH-2026-058; ADH-2026-060; ADH-2026-061 |
| In-process synthetic observation boundary | Build | The deterministic `synthetic-iaas` observer is the minimum replaceable anti-corruption boundary; it performs no real integration in Phase 2R. | Approved | DEC-0036; ADH-2026-058 |

## Capability assessment: deterministic target observation and qualification

### Identity

| Field | Value |
|---|---|
| Feature identity | FEATURE-0016 |
| Capability or decision-unit identity | Deterministic `synthetic-iaas` target observation, qualification, lifecycle publication, and safe projection |
| Assessment owner | Sanjeev Kumar, Sovrunn Architecture Owner |

### Classification

| Field | Value |
|---|---|
| Disposition | Build |
| Decision status | Approved |

### Analysis

| Field | Value |
|---|---|
| Assessment scope | The F0016-owned ExecutionTarget, NormalizedTargetFactSet, TargetQualificationResult, synthetic observation, qualification, maintenance, expiry, and safe-projection boundary only. |
| Candidate category | In-process target observation and qualification control-plane boundary. |
| Mature candidates / applicable standards | FEATURE-0012 API contract, FEATURE-0013 AuditEvent contract, FEATURE-0015 backing-resource authority, and the approved internal adapter-boundary principle in DEC-0036. |
| Relevant candidate strengths | The inherited contracts supply stable request grammar, authorization, Problem handling, audit publication, and backing-resource semantics without duplicating them. |
| Material candidate constraints | None supplies the F0016-specific closed resource model, deterministic synthetic facts, current-marker fencing, audit-before-publication, or Phase 2R external-effect prohibition. |
| Rationale | Build only the Sovrunn-owned differentiation while reusing all applicable shared contracts and retaining a thin observation boundary for later replacement. |
| Selected foundation or approach | `ExecutionTargetLifecycleService` is the sole committer; `sovrunn.synthetic-iaas-observer/v1` proposes normalized facts; F0016 reuses prior-feature contracts by reference. |

### Boundary

| Field | Value |
|---|---|
| Sovrunn-owned responsibility | Normalize four deterministic facts; evaluate the closed qualification profile; fence stale work; commit current records and target state; and project safe response-only availability. |
| Reused or extended responsibility | FEATURE-0012 owns shared API, Problem, media, authorization, and ETag contracts; FEATURE-0013 owns AuditEvent semantics; FEATURE-0015 owns participation and InfrastructureStack authority; F0016 owns its private `BackingAccessProvider` over the approved read-only lease. |
| Responsibility/control boundary | FEATURE-0016 composes inherited FEATURE-0012/0013 contracts and consumes FEATURE-0015's `CloudProviderParticipation`/`InfrastructureStack` by reference only. It does not redefine those contracts. Later realization, placement, plugin execution, IAM, and customer-projection responsibilities remain outside this feature. |
| Data crossing the boundary | Target-bound synthetic fixture input and normalized fact proposals enter the lifecycle service; internal FactSet and Result records remain unprojected. |
| Control crossing the boundary | Qualification and fixture-maintenance triggers are submitted to the lifecycle service; only it publishes committed state after the inherited audit append succeeds. |
| Adapter required | Yes |
| Adapter rationale | DEC-0036 requires a Sovrunn-owned anti-corruption boundary before any external integration; the in-process observer keeps F0016 deterministic and replaceable without coupling core state to an external system. |
| Adapter or contract identifier | `sovrunn.synthetic-iaas-observer/v1` |
| Vendor-native types allowed | No |

### Suitability

| Field | Value |
|---|---|
| Sovereignty and deployment fit | Fully in-process and deterministic; suitable for disconnected and air-gapped deployments because it has no network, credential, or vendor dependency. |
| Security and trust | No raw credential or external call crosses the boundary; F0012 authentication/safe denial and F0013 redacted audit evidence remain authoritative. |
| Operational and supportability | Injected-clock behavior, bounded in-process idempotency, deterministic fixtures, closed outcomes, and conformance cases permit repeatable local diagnosis. |
| Licensing and supply-chain | Standard-library implementation and inherited repository contracts only; no new third-party runtime dependency is selected. |
| Portability and provider-neutrality impact | The target model is CloudProvider-scoped and provider-neutral; no provider-native type or external API becomes canonical. |

### Phase and scope

| Field | Value |
|---|---|
| Allowed in current phase | Yes |
| Current-phase work | Implement deterministic synthetic observation, qualification, lifecycle, audit, idempotency, five-route HTTP behavior, and conformance proof only. |
| Deferred work | Real observation integration, realization, placement, provisioning, plugin execution, credential handling, and customer IaaS exposure require a separately approved future feature. |
| Explicit non-goals | No real adapter, provider-native type, credential, external call, external persistence, placement, provisioning, plugin execution, generic event bus, or customer target API. |
| Exit or migration boundary | A later approved feature may substitute a real observation implementation behind a compatible boundary; it must not rewrite F0016 lifecycle, qualification, audit, or safe-projection semantics without a new ADH. |
| Phase 2 non-goal acknowledgement | Phase 2R permits deterministic observation and qualification only; real infrastructure execution and external integrations remain non-goals. |

### Build justification

| Field | Value |
|---|---|
| Why Reuse is insufficient | Shared F0012/F0013/F0015 contracts are reused, but none owns an ExecutionTarget qualification model or its fenced lifecycle semantics. |
| Why Wrap is insufficient | There is no selected mature engine to wrap in Phase 2R, and wrapping one would introduce prohibited external coupling. |
| Why Extend is insufficient | Extending a prior resource or lifecycle would violate the approved F0015/F0016 ownership boundary and create competing writers. |
| Protected Sovrunn differentiation and long-term ownership | Sovrunn owns the provider-neutral target model, normalized fact vocabulary, qualification conclusions, fencing, audit boundary, and later-replaceable observation contract. |

### Risk mitigation

#### Applicable architecture risks

#### Risk-control matrix

| Risk | Preventive control | Detection control | Corrective path |
|---|---|---|---|
| External-coupling or credential leakage | Closed `synthetic-iaas` contract; no credential/external-call fields; Phase 2R exclusions. | Feature-contract, Phase 2R drift, and conformance checks. | Remove prohibited coupling and require an approved ADH before any real integration. |
| Incorrect or stale qualification publication | Sole lifecycle committer; fence recheck; ADH-2026-060 precedence; ADH-2026-061 link-clearing rule; audit-before-publication. | F16 conformance, formal models, race tests, and audit-failure cases. | Abort the proposal, preserve committed state, and correct the lifecycle implementation under the existing contract. |
| Ownership or safe-projection leakage | FEATURE-0015 references remain read-only; the ADH-2026-066 lease is private and F0016-owned; internal FactSet/Result records are never projected; closed five-route surface. | Registry, traceability, projection, safe-denial, and backing-access readiness checks. | Restore the approved ownership boundary and reject unauthorized field/route expansion. |

- Residual risk: Low. Synthetic fixtures cannot validate a real environment,
  and F0016 in-memory state intentionally resets on restart.
- Replacement risk: Medium
- Reassessment triggers: Approval of a real observation or realization feature,
  a new credential or external dependency requirement, or a change to lifecycle
  states, fence inputs, audit semantics, restart/retention policy, schema,
  route, or cross-feature ownership.

### Traceability

| Field | Value |
|---|---|
| Related DEC / RFC / ADH references | DEC-0026, DEC-0036, DEC-0042, DEC-0057, RFC-0021, ADH-2026-058, ADH-2026-060, ADH-2026-061, ADH-2026-063, ADH-2026-064, ADH-2026-065, ADH-2026-066 |
| Linked acceptance criteria | AC-F16-01 through AC-F16-12; VS0-CF-F16-01 through VS0-CF-F16-122 |
| Validation and review evidence | Approved ADH-2026-058/060/061; `docs/reviews/reuse-assessments/FEATURE-0016-approval-evidence.md`; FEATURE-0016 readiness, feature-contract, VS-000, Phase 2R drift, formal, conformance, and final feature-gate evidence. |

### Human-approval evidence

| Field | Value |
|---|---|
| Structured approval-evidence record | `docs/reviews/reuse-assessments/FEATURE-0016-approval-evidence.md` |
| Approval status | Approved |
| Approving person or role | Sanjeev Kumar, Sovrunn Architecture Owner |
| Approval date | 2026-08-20 |
| Approval basis | ADH-2026-058 and the recorded FEATURE-0016 reuse-assessment approval evidence confirm the disposition and responsibility boundary. |

---

## 2. Scope Classification

### 2.1 REQUIREMENTS-Owned Observable Behavior

| ID | Behavior | DEC/ADH | VS0 IDs |
|----|----------|---------|---------|
| REQ-F16-01 | Register ExecutionTarget as a CloudProvider-scoped resource; create accepts exactly `metadata.name`, `spec.cloudProviderParticipationRef.uid`, `spec.infrastructureStackRef.uid`, and immutable `spec.targetClass=synthetic-iaas`; scope derives solely from the participation reference; no separate scopeRef persisted or projected | DEC-0042/0057; ADH-058 clause 1 | VS0-SCHEMA-015 |
| REQ-F16-02 | Safe reference access precedes graph validation; create/qualify require FEATURE-0015 effective-Active participation and Active InfrastructureStack in the same CloudProvider scope; at most one non-Retired (participation, stack, targetClass) tuple is live; name is scope-unique for the process lifetime including after retirement | DEC-0042; ADH-058 clause 2 | VS0-SCHEMA-015, VS0-CF-F16-09..13,65,98 |
| REQ-F16-03 | Exactly five Go 1.22 http.ServeMux registrations (collection POST/GET, item GET, qualify POST, retire POST); a fixed pre-ServeMux method/path guard supplies transport-only 405/404 for HEAD and trailing-slash forms with no Problem body, authentication, audit, idempotency, or lifecycle effect | DEC-0042; ADH-058 clause 3 | VS0-CF-F16-01,14,16,20..22,36,44..45,61..64,66,90,93..96,105..110,113..118 |
| REQ-F16-04 | Every POST requires application/json and exactly one Idempotency-Key; qualify/retire require exactly one strong opaque If-Match and a zero-byte body; ETag on create/item GET/qualify/retire only; F0016 ignores Accept and emits fixed response media; closed safe projection with no public conditions member | DEC-0042; ADH-058 clause 4 | VS0-CF-F16-01..06,14,16,20..22,24..25,36,44..45,48,51..60,70..74,83..86,89..91,105..110,113..118 |
| REQ-F16-05 | Executable request precedence: transport guard, authentication, method/media/header-form validation, phase-one safe references, authorization, safe access, strict classification, graph/reference validation, replay/reservation, current If-Match comparison, lifecycle state, atomic commit | DEC-0042; ADH-058 clause 5 | VS0-CF-F16-02..13,18..30,51..60,67..75,84..86,90..91,105..111,113..118 |
| REQ-F16-06 | Idempotency namespace is principal, route, derived CloudProvider scope, concrete target UID for actions, and Idempotency-Key; digest is SHA-256 of sorted create JSON or zero bytes for actions; applies only to create/qualify/retire; 24h/10,000-entry in-process retention; independent reservations across targets and scopes | DEC-0042; ADH-058 clause 5 | VS0-CF-F16-01,20..23,31..35,46..47,49..50,76..78,81..82,86..89,99 |
| REQ-F16-07 | The only observer is sovrunn.synthetic-iaas-observer/v1: target-bound, clock-driven, external-effect-free; exactly four named facts each Supported/Unsupported/Unknown; any Unsupported rejects, otherwise any Unknown is Indeterminate, otherwise Qualified; missing fixture and logical timeout each yield four Unknown facts at now/+60s | DEC-0036/0042; ADH-058 clause 6 | VS0-CF-F16-20..23,46,67..69,90 |
| REQ-F16-08 | ExecutionTargetLifecycleService is the sole committer of persisted target state and internal fact/result records; lifecycle Active/Retired; persisted qualification Unqualified/Qualified/Rejected/Indeterminate; Qualifying is a lifecycle-service-only in-flight reservation never persisted or projected; no status.availability field; response-only effectiveAvailability computed from committed state/freshness/maintenance/viability | DEC-0042/0057; ADH-058 clause 7 | VS0-SCHEMA-015, VS0-STATE-004, VS0-CF-F16-01,20..22,26..30,36,39,41..43,75,77,79..83,87..88,99..101,119..122 |
| REQ-F16-09 | Expiry is injected-clock driven; Retired/Maintenance outrank it; maintenance enters/clears only through a deterministic target-bound synthetic-observer fixture trigger, never an HTTP route or generic event bus; Retire is the sole lifecycle exit; audit matrix for success/denial/safe-denial; expiry is the sole background retry exception; shutdown ordering stops triggers, then aborts reservations, then stops HTTP serving | DEC-0042/0057; ADH-058 clause 8 | VS0-CF-F16-01,29..30,36,38..43,66,75,77,79..83,87..88,99..104,111..112 |
| REQ-F16-10 | The only new local violations are VS0_TARGET_RETIRED, VS0_TARGET_MAINTENANCE, VS0_TARGET_QUALIFICATION_IN_PROGRESS, VS0_TARGET_EPOCH_STALE, VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE, VS0_EXECUTION_TARGET_STACK_UNAVAILABLE, VS0_EXECUTION_TARGET_SCOPE_MISMATCH, and VS0_EXECUTION_TARGET_VIABILITY_STALE; FEATURE-0012 owns top-level Problems | DEC-0042; ADH-058 clause 9 | VS0-CF-F16-10..12,26..30,38,75 |
| REQ-F16-11 | F0024/Phase 3 owns realization/plugin taxonomy; Crossplane is future-only; no F0016 dependency on Crossplane, real adapter, credential, external call, placement, provisioning, plugin execution, or customer IaaS exposure | DEC-0042; ADH-058 clause 10 | VS0-CF-F16-46 |

### 2.2 DESIGN-Delegated Mechanics

- In-memory registry data structures for ExecutionTarget, NormalizedTargetFactSet, and TargetQualificationResult
- Handler wiring and router setup behind the pre-ServeMux guard
- Idempotency reservation/waiter internal data structures
- Injected-clock scheduling internals for expiry retry

### 2.3 Downstream Concepts Visible Only as Negative Boundaries

| Element | Introduced and Activated By | FEATURE-0016 Rule |
|---------|-----------------------------|------------------|
| ServiceRegion.spec.executionTargetRefs, ServiceRegion.status.availability | FEATURE-0022 | FEATURE-0016 is only a declared prerequisite; it does not introduce or activate either field |

### 2.4 EXCLUDED (Must NOT Implement)

| Concept | Reason | Reference |
|---------|--------|-----------|
| Persisted status.availability on ExecutionTarget | Replaced by response-only effectiveAvailability | ADH-058 clause 7/F16-AD-15.2 |
| Draining as an active lifecycle/availability value | Removed placeholder | ADH-058 |
| Publicly projected Qualifying | Qualifying is lifecycle-service-only, never projected | ADH-058 clause 7 |
| Real adapter, provider-native type, raw credential | Deferred/future-only | ADH-058 clause 9/10 |
| Crossplane dependency | Future realization candidate only | ADH-058 clause 10 |
| Placement, provisioning, plugin execution | F0023/F0024 scope | ADH-058 clause 10 |
| Customer IaaS exposure, customer target projection | FEATURE-0023 safe projection scope | ADH-058 |
| adapterAuthorityRef, SecretRef extension | Unknown/rejected fields | ADH-058 clause 1/6.1 |
| ResourcePool, ProviderCapability, generic Provider, target subtype | Superseded | DEC-0042 |
| PATCH, PUT, DELETE, HEAD (as public routes), watch, filter, pagination, fact/result, adapter-selection route | Not part of the closed five-route surface | ADH-058 clause 3 |

---

## 3. Acceptance Criteria

| ID | Criterion | Conformance ID |
|----|-----------|----------------|
| AC-F16-01 | ExecutionTarget create accepts exactly the four closed fields, derives scope solely from participation, and persists Active/Unqualified at epoch 0 | VS0-CF-F16-01,02,57,84 |
| AC-F16-02 | Safe reference access precedes graph validation; inaccessible/ineffective/non-Active/cross-scope backing return their exact outcomes | VS0-CF-F16-09..13,65,98 |
| AC-F16-03 | Exactly five routes are registered; the pre-ServeMux guard supplies all HEAD/trailing-slash/unmatched-method transport outcomes with no Problem body or side effect | VS0-CF-F16-44,61..64,66,93..96,105..118 |
| AC-F16-04 | Every POST requires application/json and Idempotency-Key; qualify/retire require If-Match and zero-byte body; header-form failures precede body classification | VS0-CF-F16-05,06,24,25,58..60,70..74,85..86 |
| AC-F16-05 | Idempotency replay/reuse/isolation behaves exactly as specified across create, qualify, and retire; GET/LIST never touch idempotency | VS0-CF-F16-14,16,31..35,46..47,49..50,76..78,81..82,86..89,99 |
| AC-F16-06 | The synthetic observer yields exactly four named facts with Supported/Unsupported/Unknown truth and deterministic missing-fixture/timeout behavior | VS0-CF-F16-20..23,67..69,90 |
| AC-F16-07 | ExecutionTargetLifecycleService is the sole state/record committer; Qualifying is never projected; effectiveAvailability follows the exact truth table with no persisted status.availability | VS0-CF-F16-01,20..22,26..30,36,39,41..43,75,77,79..83,87..88,99..101,119..122 |
| AC-F16-08 | Retire wins over in-flight qualification with VS0_TARGET_RETIRED; Maintenance entry wins with VS0_TARGET_MAINTENANCE; both audited, no qualification AuditEvent | VS0-CF-F16-26..28,75,80,83,88..89 |
| AC-F16-09 | Expiry retries once per injected-clock second on audit failure and never returns 500 to a caller; shutdown ordering stops triggers before aborting reservations before stopping HTTP serving | VS0-CF-F16-41,77,79 |
| AC-F16-10 | Audit matrix covers success/denial/safe-denial exactly as specified; append failures return INTERNAL_ERROR/500 with no disclosure | VS0-CF-F16-39,40,82,99..104,112 |
| AC-F16-11 | Only the eight closed local violation codes are introduced; no other F0016-local violation code exists | VS0-CF-F16-10..12,26..30,38,75 |
| AC-F16-12 | No real adapter, credential, external call, Crossplane dependency, placement, provisioning, plugin execution, or customer IaaS exposure exists in active F0016 behavior | VS0-CF-F16-46 |

---

## 4. Traceability

Every REQ/AC above maps to `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`
(`VS0-SCHEMA-015..017`, `VS0-STATE-004`, `VS0-CF-F16-01..128`) and
`docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`. The F0016 local
conformance range is `VS0-CF-F16-01..128`: `VS0-CF-F16-123..126` were added by
ADH-2026-063 (create phase-one/strict-classification precedence), and
ADH-2026-065 refined `VS0-CF-F16-125` into the exact non-family malformed-JSON
classification case and added `VS0-CF-F16-127` (exact non-family
duplicate-top-level-member classification) and `VS0-CF-F16-128` (fail-closed
missing/unextractable required phase-one reference denial). `VS0-CF-F16-125`,
`VS0-CF-F16-127`, and `VS0-CF-F16-128` map to `REQ-F16-05`. The current
FEATURE-0016 semantic authorities are ADH-2026-058 together with
ADH-2026-060/061/063/064/065/066; ADH-2026-058 is not the sole current semantic
authority. ADH-2026-066 supplies only the F0016-owned private backing-access
bridge; it changes no public F0015 behavior or F0016 local conformance outcome.
No downstream ID (`VS0-CF-HP01`, `VS0-CF-F09`) is counted as
F0016-local proof; `VS0-CF-F10` and `VS0-CF-X03` remain cross-feature/shared
references only.

---

## 5. Non-Goals (Phase 2R)

Real CloudProvider/Kubernetes/OpenShift/PostgreSQL calls, real adapter
credentials, real placement, real provisioning, real plugin execution, and
any customer-facing target projection remain out of scope for FEATURE-0016
and Phase 2R generally.
