# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-073
- Date: 2026-09-03
- Source discussion: ChatGPT Project / Codex
- Related feature: FEATURE-0018 — Governance, IAM, Approval and Exception Foundation
- Related phase: Phase 2R
- Author: ChatGPT / Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

FEATURE-0018 ExceptionProposal terminal-decision audit-event registration

## Summary

The approved FEATURE-0018 semantics require every terminal
`ExceptionProposal` `Grant | Deny` to publish required audit evidence and exactly
one `exception-decision/v1` DecisionRecord. The bounded F18-RD-20 AuditEvent
taxonomy contains no event that truthfully represents that terminal decision. This
correction registers one missing taxonomy type, `exceptionproposal.decided`, through
the existing FEATURE-0013 extension-registration mechanism. It changes no resource,
caller action, route, decision profile, error, writer, permission, lifecycle rule,
dependency contract, persistence guarantee, or Phase 2R boundary.

## Classification

- Correction

## Existing approved baseline

`ARCH-2026.08-PHASE2R-CANONICAL`, DEC-0060, and the approved ADH-2026-070
three-file package are authoritative for FEATURE-0018. In particular:

- F18-RD-17 requires every accepted `ExceptionProposal` to reach one terminal
  `Grant | Deny`, with the registered terminal reasons, and requires `Deny` to
  create no `ExceptionGrant`.
- F18-RD-19 requires exactly one mandatory `exception-decision/v1` DecisionRecord
  for each such terminal outcome.
- F18-RD-20 requires protected AuditEvent-obligation acceptance for an
  `ExceptionProposal` terminal decision and defines the taxonomy as closed unless
  it is extended through FEATURE-0013's approved registration mechanism.
- The present taxonomy contains `exceptiongrant.proposed`,
  `exceptiongrant.issued`, `exceptiongrant.revoked`, and
  `exceptiongrant.expired`; none denotes a terminal `ExceptionProposal` decision.

`exceptiongrant.proposed` remains the submission event. It must not be reinterpreted
as a terminal decision, and `exceptiongrant.issued` must not be used to represent a
terminal `Deny`.

Relevant baseline references:

- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `docs/decisions/DEC-0060-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-071-feature-0018-conformance-executability-reconciliation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-072-feature-0018-design-authority-and-audited-evaluation-clarification.md`

## Decision or proposed decision

Register exactly one additional FEATURE-0018 AuditEvent taxonomy type:

```text
exceptionproposal.decided
```

`exceptionproposal.decided` is emitted exactly once for every terminal
`ExceptionProposal` `Grant | Deny` outcome. It is linked to the exact immutable
proposal, its terminal reason, the mandatory `exception-decision/v1` DecisionRecord,
and the containing operation/correlation evidence using the existing FEATURE-0013
AuditEvent envelope and linkage semantics.

For terminal `Grant`, `exceptionproposal.decided` and `exceptiongrant.issued` are
separate required events in the same final exception atomic boundary:

```text
ExceptionProposal terminal Grant
  -> exceptionproposal.decided
  -> exceptiongrant.issued
  -> exactly one exception-decision/v1 DecisionRecord
  -> exact immutable ExceptionGrant
```

For terminal `Deny`, the boundary publishes no `ExceptionGrant` and no
`exceptiongrant.issued` event:

```text
ExceptionProposal terminal Deny
  -> exceptionproposal.decided
  -> exactly one exception-decision/v1 DecisionRecord
  -> no ExceptionGrant
```

The existing event meanings remain fixed:

```text
exceptiongrant.proposed  -> accepted proposal submission only
exceptiongrant.issued    -> immutable ExceptionGrant issuance only
exceptiongrant.revoked   -> linked revocation evidence only
exceptiongrant.expired   -> expiry evidence only
```

This is a registration through FEATURE-0013's already-approved extension mechanism.
It does not add a FEATURE-0013 carrier type, validation rule, service, transaction,
writer, durable store, or dependency change. FEATURE-0018 remains responsible for
its feature-local, deterministic, non-durable publication mechanics and for the
conformance evidence of this registered type.

## Rationale

A terminal decision is different from proposal submission and different from grant
issuance. Without an exact event, a terminal `Deny` has no truthful required
AuditEvent mapping, while relabelling `exceptiongrant.proposed` or
`exceptiongrant.issued` would create false audit meaning. The added type closes the
already-required F18-RD-17/F18-RD-19/F18-RD-20 boundary with the smallest possible
taxonomy correction.

## Reuse-before-build assessment

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse FEATURE-0013's existing DecisionRecord/AuditEvent envelope, validation,
    linkage, projection, and extension-registration mechanism.
  - Reuse the approved F18-RD-17 terminal result and F18-RD-20 atomic-audit
    boundary without reinterpretation.
- Sovrunn-owned responsibility summary:
  - FEATURE-0018 registers and emits one feature-owned terminal-decision taxonomy
    type and supplies deterministic local conformance evidence.
- Non-goals summary:
  - No FEATURE-0013 contract, schema, implementation, service, or persistence
    change; no new resource, caller action, route, permission, role, approval rule,
    exception reason, error, writer, native-IAM effect, or future-feature behavior.

## Phase impact

- Current phase allowed? Yes
- If not current phase, target phase: Not applicable
- Current phase boundary impact:
  - Canonical contract and deterministic in-memory conformance only. No external
    I/O, production persistence, native-IAM interaction, or execution-target effect
    is introduced.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting decisions, if any:
  - The omission between F18-RD-20's required ExceptionProposal terminal-decision
    audit boundary and its enumerated taxonomy is corrected by this handoff.
- Resolution required:
  - Human approval of this correction, followed by FEATURE-0013 extension
    registration reuse and FEATURE-0018 registry/traceability/design reconciliation.
    No FEATURE-0013 core-contract change, DEC, RFC, ACR, or baseline replacement is
    required.

## Required action

- Update architecture doc
- Update feature authority
- Update the F18-RD-20 registered AuditEvent taxonomy
- Update Kiro requirements.md
- Update Kiro design.md
- Update traceability matrix and local conformance evidence
- Update feature control/review projection and approval hashes as required
- Do not update Kiro tasks.md
- Do not modify Go code

## Impacted files

- `docs/reviews/architecture-decision-handoffs/ADH-2026-073-feature-0018-exceptionproposal-terminal-decision-audit-event.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md`
- `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `docs/features/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`
- `docs/architecture/canonical/sovrunn-finalized-data-model.md` only if it records
  the bounded FEATURE-0018 AuditEvent taxonomy
- `docs/traceability/DECISION_TRACEABILITY_MATRIX.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/requirements.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/design.md`
- `.automation/context-projections/FEATURE-0018.design-review.spec.json`
- `.automation/context-projections/FEATURE-0018.design-review.md`
- `.automation/context-projections/FEATURE-0018.design-review.manifest.json`
- `.automation/features/FEATURE-0018.control.json`
- `.automation/state/FEATURE-0018.json` only for approved context/spec hash reconciliation

## Impacted features

- FEATURE-0018: adds one exact audit taxonomy registration required by its already
  approved terminal-exception decision boundary.
- FEATURE-0013: no core change. Its existing extension-registration mechanism and
  carrier contracts are reused unchanged.
- FEATURE-0012, FEATURE-0016, FEATURE-0017, FEATURE-0019 through FEATURE-0026:
  no change.

## Acceptance criteria for Kiro update

- [ ] Validate this handoff against the Architecture Operating System and confirm
  that `exceptionproposal.decided` is the only new taxonomy entry.
- [ ] Register `exceptionproposal.decided` through the existing FEATURE-0013
  extension mechanism without changing FEATURE-0013 carrier contracts or writers.
- [ ] Preserve every F18-RD-17 terminal reason and the exact `Grant | Deny`
  semantics.
- [ ] Require exactly one `exceptionproposal.decided` event and exactly one
  `exception-decision/v1` DecisionRecord for every terminal proposal decision.
- [ ] Require `exceptiongrant.issued` only for a terminal `Grant`; terminal `Deny`
  creates neither an `ExceptionGrant` nor an issuance event.
- [ ] Preserve `exceptiongrant.proposed` as proposal-submission evidence only.
- [ ] Reconcile the requirements, design, registry, traceability, controlled review
  context, and local conformance evidence atomically.
- [ ] Add no resource, route, action, error, writer, dependency, persistence,
  external effect, or future-feature behavior.
- [ ] Do not generate tasks or modify Go code.

## Explicit instructions to Kiro

- Do not introduce a second event type for the terminal `Grant` or `Deny` path.
- Do not repurpose `exceptiongrant.proposed` or `exceptiongrant.issued` to represent
  the terminal decision.
- Do not alter FEATURE-0013's carrier envelope, common validation, writer authority,
  persistence posture, or extension mechanism.
- Do not change the terminal exception reason registry, ExceptionGrant resource
  shape, caller operation surface, approval semantics, or Phase 2R exclusions.
- Preserve the separate approval and final-exception atomic boundaries.
- Do not generate tasks or Go code from this handoff.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-09-03
- Notes: Approved as the narrow correction required to make the existing
  F18-RD-17/F18-RD-19/F18-RD-20 terminal-exception evidence boundary exact.
