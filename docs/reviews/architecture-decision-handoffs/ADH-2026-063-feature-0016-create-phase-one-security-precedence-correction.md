# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-063
- Date: 2026-08-21
- Source discussion: Codex architecture session
- Related feature: FEATURE-0016
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Correct FEATURE-0016 collection-create phase-one security precedence

## Summary

This correction closes a discovered precedence gap between the approved F0016
request pipeline and its implementation. For collection create, phase one may
capture a bounded request body and extract only the two references needed for
scope derivation. It must not emit a malformed-JSON, duplicate-member, unknown-
field, status-write, or other body-classification Problem. Authentication,
derived-scope authorization, and safe backing access occur first; strict
classification then returns the existing body outcome for an authorized,
safe-accessible request.

No resource, route, field, state, top-level Problem code, dependency, or
future-phase behavior is added.

## Classification

Correction.

## Existing approved baseline

- ADH-2026-058 and the F0016 architecture define the order: transport guard,
  authentication, method/media/header-form validation, phase-one safe
  references, authorization, safe access, strict classification, then graph
  validation and publication.
- The same authority requires safe 404 for inaccessible backing and audited 403
  for a safe-accessible caller lacking `executiontarget.write`.
- F16-51, F16-52, and F16-53 define existing malformed, duplicate, and
  oversized-body outcomes; F16-08 and F16-09 define authorization and safe-
  denial outcomes. Their combined-input precedence was not recorded.

## Decision

For `POST /apis/execution.sovrunn.io/v1alpha1/execution-targets`, the fixed
order is:

1. transport guard, authentication, and existing method/media/header-form
   validation;
2. bounded single read and phase-one extraction of only
   `spec.cloudProviderParticipationRef.uid` and
   `spec.infrastructureStackRef.uid`, without emitting any body-classification
   Problem;
3. derived-CloudProvider `executiontarget.write` authorization and safe backing
   access;
4. strict classification of the retained same request bytes, followed by the
   existing graph/reference, idempotency, audit, and publication stages.

Consequently:

- an authenticated caller without `executiontarget.write` and a malformed or
  duplicate body receives the existing audited 403 authorization denial;
- an authenticated caller whose backing is inaccessible and whose body is
  malformed or duplicate receives the existing audited safe 404;
- an authorized, safe-accessible caller with malformed JSON receives 400
  `MALFORMED_REQUEST`, and with a duplicate member receives 400
  `DUPLICATE_FIELD`;
- an oversized body retains its existing bounded `REQUEST_TOO_LARGE` outcome
  before phase-one extraction, authorization, or safe access.

Phase one neither classifies nor digests/reserves a body. Strict decoding uses
the retained same bytes exactly once after safe access; invalid, unknown, and
duplicate bodies are never digested or reserved.

## Rationale

This preserves the approved safe-denial boundary: an untrusted body cannot
select a different result before authorization or disclose backing existence.
The request-size bound remains an unavoidable transport/resource limit and is
not a body-classification outcome.

## Reuse-before-build assessment

- Disposition: Clarify existing reuse.
- FEATURE-0012 body limits, strict decode/classification, and Problem envelopes
  are reused unchanged.
- FEATURE-0013 audited authorization/safe-denial semantics are reused unchanged.
- FEATURE-0015 derived-scope authorization and safe backing access remain
  read-only authority.
- No real adapter, credential, external call, durable persistence, controller,
  plugin execution, placement, provisioning, or ExecutionTarget subtype is
  introduced.

## Phase impact

- Current phase allowed? Yes.
- Current phase boundary impact: none; this is a deterministic in-process
  precedence correction within the existing F0016 HTTP contract.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting authority: the current create decoder emits strict body errors
  before the approved authorization/safe-access stages.
- Resolution required: founder-approved ADH application only; no DEC, RFC, or
  architecture-baseline update.

## Required action

- Update the F0016 architecture authority, registry/conformance catalog,
  contract specification, traceability, steering, closure evidence, control
  manifest, and fail-closed checks.
- Regenerate the exact F0016 requirements ledger and revise design/tasks to
  consume the corrected precedence.
- After the revised tasks plan is approved, update the create handler and tests
  in a bounded implementation task. Do not apply Go changes during the
  architecture-update pass.

## Strict Kiro read/write boundary

Kiro may read only this handoff and the following files while applying it. It
may write only files marked **write**.

| Path | Access | Purpose |
|---|---|---|
| `docs/reviews/architecture-decision-handoffs/ADH-2026-063-feature-0016-create-phase-one-security-precedence-correction.md` | read | Controlling approved decision |
| `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md` | write | Add controlling handoff and exact create precedence |
| `docs/reviews/architecture-readiness/FEATURE-0016-architecture.md` | write | Align decision register and proof matrix |
| `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` | write | Add four exact combined-input conformance rows and amend affected evidence |
| `docs/architecture/vertical-slices/VS-000-contract-specification.md` | write | Record phase-one and strict-classification order |
| `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md` | write | Add ADH-063 provenance and proof links |
| `.kiro/steering/slice0-contract.md` | write | Record the safe precedence invariant |
| `docs/reviews/architecture-readiness/FEATURE-0016-architecture-closure-matrix.md` | write | Add resolved closure row |
| `.automation/features/FEATURE-0016.control.json` | write | Add ADH-063 only to `feature.handoffs` |
| `scripts/feature-0016-architecture-readiness-check.py` | write | Fail closed on pre-auth body classification |
| `scripts/feature-contract-check.py` | write | Verify F0016 precedence evidence |
| `scripts/vs000-contract-check.py` | write | Verify exact registry semantics |
| `.kiro/specs/adapter-boundary-and-executiontarget-qualification/requirements.md` | write | Regenerate exact ledger and supplemental traceability only |

## Acceptance criteria for Kiro update

- [ ] Four combined-input cases exist: unauthorized+malformed, inaccessible+
  malformed, authorized/safe malformed-or-duplicate, and oversized body.
- [ ] Phase one returns no strict body-classification Problem and body bytes are
  not digested/reserved before strict classification.
- [ ] Existing F16-08, F16-09, F16-51, F16-52, and F16-53 meanings remain
  unchanged.
- [ ] No field, route, state, top-level Problem code, dependency, or phase scope
  changes.
- [ ] Requirements ledger is mechanically regenerated; canonical REQ/AC ledgers
  remain untouched.
- [ ] Architecture readiness, feature-contract, VS-000, Phase 2R drift, formal,
  documentation, Structurizr, and diff checks pass.

## Explicit instructions to Kiro

- This handoff is founder-approved and must be applied only through the strict
  Kiro architecture-update workflow.
- Do not modify Go source, design.md, tasks.md, or formal models in this pass.
- Do not introduce a new public route, error code, resource, state, dependency,
  or future-phase behavior.
- Preserve the existing oversized-body outcome and the existing action If-Match
  header-form precedence.
- If the registry cannot express the four proofs without changing a published
  semantic outside this handoff, stop with `ARCHITECTURE_DECISION_REQUIRED`.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-21
- Notes: Founder selected authorization and safe access before strict
  collection-create body classification, with the bounded request-size outcome
  retained before phase-one extraction.
