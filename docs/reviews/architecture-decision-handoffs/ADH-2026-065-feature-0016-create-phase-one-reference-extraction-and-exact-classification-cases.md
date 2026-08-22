# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-065
- Date: 2026-08-21
- Source discussion: Codex architecture session
- Related feature: FEATURE-0016
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Correct FEATURE-0016 create phase-one reference extraction and exact
classification conformance cases

## Summary

ADH-2026-063 correctly establishes that strict collection-create body
classification happens only after authentication, phase-one safe-reference
extraction, authorization, and safe backing access. Two small contract gaps
remain:

1. `VS0-CF-F16-125` combines malformed JSON and duplicate top-level-member
   JSON despite their different registered `Problem.code` values.
2. The contract does not state the fail-closed outcome when phase one cannot
   extract exactly one usable UID for each required backing reference.

This correction makes each observable conformance case exact. It neither adds
a route, resource field, state, dependency, top-level Problem code, nor an
external effect.

## Classification

Correction.

## Existing approved baseline

- ADH-2026-063 is the controlling correction for collection-create body
  precedence.
- Strict classification remains after authentication, bounded phase-one
  extraction, authorization, and safe backing access.
- `MALFORMED_REQUEST`, `DUPLICATE_FIELD`, and `AUTHORIZATION_DENIED` are
  existing FEATURE-0012 Problem codes; no code catalog change is authorized.
- Existing `VS0-CF-F16-123`, `124`, and `126` remain exact and unchanged.

## Decision

### 1. Split the divergent strict-classification family

`VS0-CF-F16-125` means only:

> An authenticated, authorized caller with safe-accessible backing submits a
> retained collection-create body containing malformed JSON. Strict
> classification returns `400 MALFORMED_REQUEST`; no mutation, AuditEvent,
> idempotency digest, reservation, or completion occurs.

Add `VS0-CF-F16-127` for only:

> An authenticated, authorized caller with safe-accessible backing submits a
> retained collection-create body containing a duplicate top-level member.
> Strict classification returns `400 DUPLICATE_FIELD`; no mutation,
> AuditEvent, idempotency digest, reservation, or completion occurs.

The cases retain the existing F16-51 and F16-52 behavior respectively. They
are distinct exact cases, not a mixed equivalence family.

### 2. Fail closed when phase-one references are not both extractable

After successful authentication and the existing bounded body read, phase one
attempts to extract only:

- `spec.cloudProviderParticipationRef.uid`; and
- `spec.infrastructureStackRef.uid`.

If phase one cannot extract exactly one syntactically usable UID at both paths,
it must fail closed before derived-scope authorization and safe backing access:

- return the existing audited `403 AUTHORIZATION_DENIED` outcome;
- disclose no body-classification detail and no backing-resource existence;
- perform no strict classification, canonicalization, digest, idempotency
  reservation, observer invocation, mutation, publication, or completion.

Add `VS0-CF-F16-128` as the exact proof case for that outcome. It reuses the
existing audited authorization-denial behavior; it adds no violation or
top-level Problem code.

### 3. Retained exact denial families

`VS0-CF-F16-123` and `VS0-CF-F16-124` may continue to name malformed JSON or a
duplicate top-level member because every member of each family has the same
observable denial outcome: respectively the existing audited authorization
denial and the existing audited safe denial. `VS0-CF-F16-126` remains the
bounded request-size outcome before phase one.

## Rationale

Conformance rows may use an equivalence family only when every member has the
same status, Problem code, violation where applicable, audit effect, and state
effect. Splitting F16-125 restores that invariant. An unextractable reference
cannot provide the derived scope or backing identity needed for authorization
and safe access, so an audited authorization denial is the fail-closed outcome
consistent with ADH-063's non-disclosure ordering.

## Scope and non-goals

- No new route, resource, field, state, writer, top-level Problem code,
  conformance behavior outside F16-125/F16-127/F16-128, implementation, or
  external effect is authorized.
- This does not alter F16-123, F16-124, or F16-126.
- This handoff authorizes architecture/specification updates only. Go handler
  and test work follows only from a newly approved implementation plan.

## Required action

- Correct F16-125 and add F16-127/F16-128 in the contract registry,
  specification, traceability, and FEATURE-0016 architecture authority.
- Update the FEATURE-0016 control handoff evidence and readiness/contract
  checks to fail closed on the three exact cases.
- Mechanically regenerate the requirements conformance ledger and map
  F16-125, F16-127, and F16-128 to `REQ-F16-05` and affected acceptance proof.
- In requirements traceability, include ADH-2026-061, ADH-2026-064, and
  ADH-2026-065; do not represent ADH-2026-058 as the sole current semantic
  authority.

## Strict Kiro read/write boundary

| Path | Access | Purpose |
|---|---|---|
| `docs/reviews/architecture-decision-handoffs/ADH-2026-065-feature-0016-create-phase-one-reference-extraction-and-exact-classification-cases.md` | read | Controlling approved correction |
| `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md` | write | Record exact phase-one and classification behavior |
| `docs/features/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md` | write | Keep feature authority aligned |
| `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` | write | Correct F16-125 and add F16-127/F16-128 |
| `docs/architecture/vertical-slices/VS-000-contract-specification.md` | write | Record exact error/precedence contract |
| `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md` | write | Map exact proof and ADH provenance |
| `.kiro/steering/slice0-contract.md` | write | Keep stage steering aligned |
| `docs/reviews/architecture-readiness/FEATURE-0016-architecture-closure-matrix.md` | write | Record resolved closure |
| `.automation/features/FEATURE-0016.control.json` | write | Add ADH-065 only to `feature.handoffs` |
| `scripts/feature-0016-architecture-readiness-check.py` | write | Add fail-closed exact-case checks |
| `scripts/feature-contract-check.py` | write | Add feature-level exact-case checks |
| `scripts/vs000-contract-check.py` | write | Add registry exact-case checks |
| `.kiro/specs/adapter-boundary-and-executiontarget-qualification/requirements.md` | write | Mechanically regenerate affected ledger and traceability only |

## Acceptance criteria for Kiro update

- [ ] F16-125 is malformed JSON only and returns `400 MALFORMED_REQUEST`.
- [ ] F16-127 is duplicate top-level member only and returns `400 DUPLICATE_FIELD`.
- [ ] F16-128 is missing/unextractable required phase-one reference and returns
  the existing audited `403 AUTHORIZATION_DENIED` outcome without body detail
  or backing disclosure.
- [ ] F16-123, F16-124, and F16-126 remain unchanged.
- [ ] Registry, contract specification, traceability, architecture authority,
  requirements ledger, control evidence, and fail-closed checks agree.
- [ ] No route, field, state, top-level Problem code, implementation, formal
  model, or Phase 2R boundary changes.

## Explicit instructions to Kiro

- Apply this approved correction before any design regeneration.
- Preserve the exact canonical REQ/AC ledgers in `requirements.md`; regenerate
  only affected supplemental conformance and traceability content.
- Do not modify `design.md`, `tasks.md`, Go source, tests, formal models,
  Structurizr, or unrelated FEATURE-0016 authorities.
- Stop with `ARCHITECTURE_DECISION_REQUIRED` if an additional observable
  create-precedence outcome is required outside this allowlist.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-21
- Notes: Founder approved exact splitting of strict-classification outcomes and
  the fail-closed authorization-denial behavior for unextractable required
  phase-one references.
