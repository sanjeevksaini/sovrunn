# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-068
- Date: 2026-08-25
- Source discussion: Codex FEATURE-0017 Kiro pre-generation review
- Related feature: FEATURE-0017
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved
- Authority payload SHA-256:
  `6c85554764edacf500dc2376ad10ed379d7c7b9873dee86b92d892da066ec3d6`
- Authority payload rule: SHA-256 of the exact UTF-8 bytes between the
  `authority-payload:start` and `authority-payload:end` marker lines, excluding
  both marker lines and with the final newline before the end marker included.

<!-- authority-payload:start -->
## Decision title

FEATURE-0017 structural-only FEATURE-0013 mapper clarification

## Summary

The FEATURE-0017 mapper remains pure and non-persisting. It validates its
complete `PolicyEvaluationResult`, evaluator identity/version, timing, and
structural FEATURE-0013 input identity, then constructs an `EvaluationResult`.
It does not validate adopting-`DecisionProfile` compatibility or DecisionRecord
scope. The later adopting domain performs those validations before adoption.

## Classification

Correction. This closes an internal ambiguity in CDG-F17-05 without changing
the six-group architecture, phase boundary, feature sequence, or baseline.

## Existing approved baseline

- ADH-2026-067 and the sole FEATURE-0017 architecture define a separate,
  caller-invoked pure mapper.
- FEATURE-0013 owns `EvaluationResult`, `DecisionProfile`, DecisionRecord scope,
  profile compatibility, authoritative adoption, and publication semantics.
- FEATURE-0017 owns no profile lookup, DecisionRecord scope authority,
  persistence, DecisionRecord writer, or AuditEvent writer.

## Decision or proposed decision

The mapper inputs are exactly:

- one complete valid `PolicyEvaluationResult`;
- immutable evaluator identity and version;
- the evaluation-boundary-captured FEATURE-0013 timing envelope; and
- structural FEATURE-0013 input identity consisting of a snapshot reference,
  an integrity carrier, or both.

The mapper validates only the structural validity of those inputs and the
closed FEATURE-0017-to-FEATURE-0013 field/status/timing mapping. Malformed
structural input makes mapping fail without changing the valid policy result.

The mapper does not receive or resolve a `DecisionProfile`, does not receive
DecisionRecord scope, and does not decide whether an evaluator, input-identity
arrangement, result schema, limits, or scope are compatible with an adopting
profile. The later adopting domain validates the constructed
`EvaluationResult` against its resolved `DecisionProfile` and DecisionRecord
scope before authoritative adoption.

## Rationale

Profile compatibility and record scope are downstream adoption semantics.
Keeping them out of the mapper preserves FEATURE-0013 ownership, avoids a
profile dependency in the generic evaluation seam, and makes the mapper
deterministic, pure, and usable before an authoritative DecisionRecord exists.

## Reuse-before-build assessment

The ADH-2026-067 feature-level reuse summary remains unchanged. FEATURE-0017
reuses FEATURE-0013 `EvaluationResult` structures and semantics, while the
later adopting domain reuses FEATURE-0013 profile/scope validation. No new
capability, dependency, adapter, or Build disposition is introduced.

## Phase impact

- Current phase allowed: Yes.
- Phase 2R remains deterministic, in-process, transient, and side-effect-free.
- No real policy engine, profile resolver, public route, store, or external
  execution is introduced.
- No sequence or baseline change is required.

## Conflict check

- No accepted DEC/RFC conflict is introduced.
- The correction removes the contradiction between requiring profile-
  compatibility rejection and prohibiting mapper profile lookup while omitting
  a profile from mapper inputs.
- FEATURE-0013 validation and adoption authority remains unchanged.

## Required action

- Update CDG-F17-05 and its conformance wording in the sole FEATURE-0017
  architecture.
- Preserve the six decision groups and every other ADH-2026-067 decision.
- Do not generate requirements, design, tasks, or implementation as part of
  this clarification.

## Impacted files

- `docs/architecture/policy-evaluation-abstraction.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-068-feature-0017-structural-mapper-clarification.md`

## Impacted features

- FEATURE-0013: ownership is preserved; no contract or implementation change.
- FEATURE-0017: mapper validation responsibility is narrowed to structural
  mapping inputs and closed mapping semantics.
- Later adopting features: retain resolved-profile and DecisionRecord-scope
  validation before authoritative adoption.

## Acceptance criteria for Kiro update

- [ ] The mapper input list contains no `DecisionProfile` or DecisionRecord
      scope.
- [ ] The mapper rejects malformed structural inputs only and does not claim
      adopting-profile compatibility validation.
- [ ] The later adopting domain explicitly owns validation against its resolved
      `DecisionProfile` and DecisionRecord scope.
- [ ] Mapping remains separate, caller-invoked, deterministic, pure,
      non-persisting, and side-effect-free.
- [ ] CDG-F17-01 through CDG-F17-06 remain the only decision groups.
- [ ] No requirements, design, tasks, Go code, public route, store, real engine,
      or external integration is introduced.

## Explicit instructions to Kiro

- Apply only this structural-mapper correction to the sole FEATURE-0017
  architecture.
- Do not pass or resolve a `DecisionProfile` in the mapper.
- Do not move FEATURE-0013 adoption, profile compatibility, record-scope
  validation, DecisionRecord publication, or AuditEvent publication into
  FEATURE-0017.
- Preserve every other ADH-2026-067 decision and the Phase 2R exclusions.
<!-- authority-payload:end -->

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-25
- Notes: The user explicitly approved the structural-only mapper correction in
  the source architecture session.

## Application Record

- Application status: Applied to the canonical FEATURE-0017 architecture.
- Application date: 2026-08-25
- Canonical artifact: `docs/architecture/policy-evaluation-abstraction.md`
- Verification: Authority-payload digest verified; the structural-only mapper
  correction is incorporated in CDG-F17-05 and its local conformance inventory.
