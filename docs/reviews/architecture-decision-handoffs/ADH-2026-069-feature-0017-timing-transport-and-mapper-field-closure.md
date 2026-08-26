# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-069
- Date: 2026-08-25
- Source discussion: Codex FEATURE-0017 Kiro pre-generation review
- Related feature: FEATURE-0017
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved
- Authority payload SHA-256:
  `38e038f7de4329b2e5012fefa0f5b092a1d1283d499b59eb37f02049b87b614b`
- Authority payload rule: SHA-256 of the exact UTF-8 bytes between the
  `authority-payload:start` and `authority-payload:end` marker lines, excluding
  both marker lines and with the final newline before the end marker included.

<!-- authority-payload:start -->
## Decision title

FEATURE-0017 success timing transport and exact FEATURE-0013 mapper-field
closure

## Summary

FEATURE-0017 closes the success-path data flow required by its separate pure
FEATURE-0013 mapper and removes ambiguity about optional mapper output fields.
A successful evaluation returns its complete `PolicyEvaluationResult` together
with the evaluation-boundary-captured FEATURE-0013 timing envelope as transient
in-process metadata. This pair is not a new canonical type or resource. The
mapper populates an exact closed set of FEATURE-0013 fields and leaves the
remaining optional fields absent. Adapter invocation remains at most once and
occurs exactly once only after every pre-invocation cancellation, deadline,
validation, and digest step succeeds.

## Classification

Correction. This closes two internal data-flow ambiguities and one subordinate
summary error without changing the six decision groups, canonical request or
result schema, feature sequence, phase boundary, or active baseline.

## Existing approved baseline

- ADH-2026-067 defines the six-group FEATURE-0017 architecture and a separate,
  caller-invoked pure FEATURE-0013 mapper.
- ADH-2026-068 keeps that mapper structural-only and assigns adopting-profile
  and DecisionRecord-scope validation to the later adopting domain.
- `PolicyEvaluationResult` contains `evaluatedAt` but does not contain the
  FEATURE-0013 `timing.startedAt` value required by the mapper.
- FEATURE-0013 requires evaluator identity, result status, timing, and
  `evaluatedAt`; structural validation additionally requires
  `inputSnapshotRef`, `inputIntegrity`, or both. FEATURE-0017 already requires
  the complete policy result in `EvaluationResult.result`.
- FEATURE-0017 has no store, public route, controller, or automatic mapping.

## Decision or proposed decision

### Successful evaluation return

On success, the in-process evaluation operation returns:

- one complete canonical `PolicyEvaluationResult`; and
- the exact evaluation-boundary-captured FEATURE-0013 `TimingEnvelope`
  containing `startedAt` and `completedAt`, with `durationMs` absent.

The result and timing envelope travel together as transient in-process return
metadata. Their pairing has no API identity, schema, resource kind, persistence,
retention, public projection, or independent lifecycle. Design may represent
the pair with a private receipt, tuple, or equivalent private carrier, but must
not add a canonical contract. On a normalized non-result failure, the operation
returns neither a `PolicyEvaluationResult` nor a success timing envelope.

The separate pure mapper consumes the exact returned timing envelope. It does
not recapture time, use hidden mutable state, read a store, or reconstruct
`startedAt` from `evaluatedAt`.

### Exact FEATURE-0013 mapper output

For every successfully mapped outcome, FEATURE-0017 populates exactly:

- `evaluator` from the caller-supplied immutable evaluator identity/version;
- `inputSnapshotRef`, `inputIntegrity`, or both, exactly as supplied after
  structural validation;
- `resultStatus` using the closed FEATURE-0017 mapping;
- `result` with the complete `PolicyEvaluationResult`;
- `timing.startedAt` and `timing.completedAt` from the returned boundary timing
  envelope; and
- `evaluatedAt`, identical to both `PolicyEvaluationResult.evaluatedAt` and
  `timing.completedAt`.

FEATURE-0017 mapper output always omits:

- `timing.durationMs`;
- `executionConfig`;
- `trustBoundary`;
- `safetyPolicyFilters`;
- `structuralTrustState`; and
- `scopeEvidence`.

The mapper does not inspect a `DecisionProfile` or conditionally populate those
omitted fields. A later adopting domain requiring any omitted field must use
its own approved `EvaluationResult` construction/adoption path or obtain a
separately approved FEATURE-0017 mapper extension.

This exact field rule supersedes only the ADH-2026-067 statement that other
optional FEATURE-0013 fields might later be populated by this mapper when an
adopting profile requires them. All other ADH-2026-067 decisions remain intact.

### Adapter invocation cardinality

The evaluation boundary invokes the adapter at most once. It invokes exactly
once only when the already-observed cancellation/deadline check, request
validation, digest calculation, and immediate pre-invocation
cancellation/deadline recheck all succeed. Cancellation or deadline observed
before invocation produces zero adapter invocations.

## Rationale

The mapper cannot receive the boundary-captured `startedAt` unless successful
evaluation transports it to the caller. Explicit transient return metadata
closes that data flow without expanding the canonical `PolicyEvaluationResult`
schema or adding persistence. Enumerating populated and omitted FEATURE-0013
fields prevents “all optional fields” from incorrectly excluding the required
input-identity carrier or complete result payload. Exact invocation cardinality
preserves the approved cancellation precedence.

## Reuse-before-build assessment

- Disposition: Reuse and Build, unchanged from ADH-2026-067.
- Applicable standard: reuse the FEATURE-0013 `TimingEnvelope`,
  `EvaluationResult`, evaluator identity, and structural input-identity
  contracts.
- Sovrunn-owned responsibility: build only the private transient success
  carrier and pure closed mapper required by the FEATURE-0017 seam.
- Non-goals: no new canonical resource/type, public API, store, controller,
  clock service, DecisionProfile inspection, DecisionRecord publication, real
  policy engine, or external effect.

## Phase impact

- Current phase allowed: Yes.
- Phase 2R remains deterministic, in-process, transient, audit-first, and
  side-effect-free.
- No schema, route, persistence, production topology, real adapter, or feature
  sequence change is introduced.

## Conflict check

- Accepted DEC/RFC conflict: No.
- ADH-2026-067 conflict: its conditional optional-field sentence is narrowly
  superseded as stated above; its six-group closure otherwise remains intact.
- ADH-2026-068 conflict: No; structural-only validation and later-domain
  adoption ownership remain unchanged.
- Baseline update required: No.

## Required action

- Update architecture doc.
- Update feature scope contract.
- Update DEC/RFC provenance.
- Update traceability matrices.
- Update the FEATURE-0017 control manifest handoff inventory.
- Do not generate requirements, design, tasks, or implementation in this
  correction.

## Impacted files

- `docs/architecture/policy-evaluation-abstraction.md`
- `docs/features/FEATURE-0017-policy-evaluation-abstraction.md`
- `docs/decisions/DEC-0028-policy-engine-adapter.md`
- `docs/rfc/RFC-0025-policy-evaluation-abstraction.md`
- `docs/traceability/DECISION_TRACEABILITY_MATRIX.md`
- `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md`
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`
- `.automation/features/FEATURE-0017.control.json`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-069-feature-0017-timing-transport-and-mapper-field-closure.md`

## Impacted features

- FEATURE-0013: contracts and ownership are reused unchanged; no
  DecisionProfile or DecisionRecord-scope validation moves to FEATURE-0017.
- FEATURE-0017: success timing transport, exact mapper fields, and invocation
  cardinality become unambiguous.
- Later adopting features: continue to own any enriched EvaluationResult path,
  profile/scope compatibility, authoritative adoption, publication, and audit.

## Acceptance criteria for Kiro update

- [ ] Successful evaluation transports the complete policy result and exact
      boundary timing together without creating a canonical type or store.
- [ ] Non-result failure transports neither a result nor success timing.
- [ ] The mapper neither recaptures time nor reconstructs `startedAt`.
- [ ] Populated and omitted FEATURE-0013 fields are enumerated exactly.
- [ ] `inputSnapshotRef` and/or `inputIntegrity`, and the complete `result`, are
      not incorrectly removed as optional fields.
- [ ] Adapter invocation is at most once and respects both pre-invocation
      cancellation/deadline checks.
- [ ] CDG-F17-01 through CDG-F17-06 remain the only decision groups.
- [ ] No requirements, design, tasks, Go code, route, store, controller, real
      engine, or external integration is introduced.

## Explicit instructions to Kiro

- Apply only the timing-transport, exact mapper-field, invocation-cardinality,
  and provenance corrections in this handoff.
- Treat any success carrier as design-private and noncanonical.
- Do not add fields to `PolicyEvaluationResult` or alter FEATURE-0013 schemas.
- Do not inspect a DecisionProfile or validate DecisionRecord scope in the
  FEATURE-0017 mapper.
- Do not introduce requirements, design, tasks, implementation, or future
  policy-engine work during this correction.
<!-- authority-payload:end -->

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-25
- Notes: The user explicitly directed Codex to resolve the three Kiro
  pre-generation content findings and their authority provenance.

## Application Record

- Application status: Applied to the canonical FEATURE-0017 architecture.
- Application date: 2026-08-25
- Canonical artifact: `docs/architecture/policy-evaluation-abstraction.md`
- Verification: Authority-payload digest verified; success timing transport,
  exact mapper fields, and adapter invocation cardinality are incorporated in
  CDG-F17-05 and CDG-F17-06 and their local conformance inventory.
