# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-075
- Date: 2026-09-05
- Source discussion: ChatGPT Project / Codex
- Related feature: FEATURE-0018 — Governance, IAM, Approval and Exception Foundation
- Related phase: Phase 2R
- Author: ChatGPT / Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

FEATURE-0018 executable-design normalization and deterministic contract check

## Summary

FEATURE-0018's approved product architecture remains unchanged, but its generated
design has accumulated repeated and sometimes contradictory implementation-mechanics
prose. This handoff authorizes a bounded normalization: retain one authoritative,
machine-checkable design-mechanics contract; make the remaining design cite it rather
than restate it; and validate that contract before an LLM semantic review. It does not
authorize a full design regeneration, a runtime behavior change, or modification of
the completed Task 1–3 Cursor checkpoints.

## Classification

- Correction

## Existing approved baseline

`ARCH-2026.08-PHASE2R-CANONICAL`, DEC-0060, ADH-2026-070 and its appendices
remain the FEATURE-0018 semantic authority. ADH-2026-071 fixes local conformance
executability; ADH-2026-072 confines FEATURE-0013 to decision/audit contracts and
validation while FEATURE-0018 owns local, deterministic, non-durable publication
machinery; ADH-2026-073 closes the terminal ExceptionProposal event taxonomy; and
ADH-2026-074 closes narrow sequence, completion-binding, reservation and gate-path
mechanics.

The current implementation has completed Task 1–3 checkpoints. No future Cursor task
may rely on a revised design until the normalized design, its task plan, and the
required hash-bound approvals are renewed.

Relevant baseline references:

- `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `docs/decisions/DEC-0060-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-071-feature-0018-conformance-executability-reconciliation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-072-feature-0018-design-authority-and-audited-evaluation-clarification.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-073-feature-0018-exceptionproposal-terminal-decision-audit-event.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-074-feature-0018-local-publication-checkpoint-reconciliation.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/requirements.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/design.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/tasks.md`

## Decision or proposed decision

### 1. Normalize, do not regenerate, the design

Kiro must preserve every approved FEATURE-0018 semantic statement, resource, action,
route, writer, error, DecisionProfile, AuditEvent, REQ, AC, conformance identifier,
dependency boundary and Phase 2R exclusion. It must not regenerate the design from a
new inferred model.

Kiro may replace duplicate or historical implementation-mechanics prose with exact
citations to one authoritative **Design Mechanics Contract**. The main design retains
feature scope, component intent, routes, ownership, test intent and traceability; the
contract retains only executable mechanics and their deterministic proof obligations.

### 2. Establish one authoritative Design Mechanics Contract

Kiro must create a machine-readable, hash-bound FEATURE-0018 Design Mechanics
Contract beside the Kiro design specification and make it the single source of truth
for the following non-product-semantic mechanics:

- the exhaustive repository-owned package-edge ledger;
- named Go 1.22 API signatures, result variants, factories and permitted callers;
- `StateEditor` method surface and exact owner-to-method allowlist;
- transaction modes and their Begin, admission, finalization, Seal, Commit, Abort,
  reservation cleanup, waiter notification and result-release state machines;
- deterministic identity, publication-instant, sequence, overflow and failure rules;
- static-checker scope, test-only seams, positive fixtures and isolated negative
  fixtures; and
- exact non-runtime proof artifact, command and ownership mappings for active
  conformance gates.

The contract may use explicit structured fields rather than prose placeholders. No
field may use `any`, `none`, ellipses, unnamed algebraic values, generic package
placeholders or an open-ended caller/owner category. Where a Go concept is modeled,
it must be representable by Go 1.22 and distinguish an actual Go `error` from the
closed FEATURE-0018 mechanical-failure result without changing public Problem
semantics.

### 3. Add a deterministic pre-review contract checker

Kiro must add a standard-library-only FEATURE-0018 design-contract checker and a
named Make target. The checker runs before the LLM design review and rejects at least:

- a mismatch between the component dependency ledger and the exhaustive package-edge
  ledger;
- undeclared or pseudo-Go API types and result variants;
- a missing factory, permitted caller, owner-method rule, failure path or test seam
  for a declared mechanics-contract surface;
- an undefined or duplicate transaction-state transition, lock order, notification
  point or context-identity component;
- a repeated, repurposed or unmapped active FEATURE-0018 conformance identifier; and
- an undefined proof artifact or command for an active non-runtime conformance gate.

The checker validates repository-owned package edges only. It must not prohibit
approved Go standard-library imports. Its test-only seams must be bounded to
`*_test.go` or testdata and must not weaken the production caller restrictions.

### 4. One bounded fixpoint review

Before resubmitting the design, Kiro must execute the deterministic checker, the
existing design semantic guardrail, the Phase 2R drift check and a read-only fixpoint
audit against the Design Mechanics Contract. The audit must classify every finding as
one of `STALE_TRANSCRIPTION`, `DESIGN_EXECUTABILITY`,
`REQUIREMENT_GAP`, `ARCHITECTURE_CLARIFICATION_REQUIRED`, or `OUT_OF_SCOPE`.

The subsequent independent design review is capped at two attempts. If an issue falls
outside the declared checker and audit scope, it must be classified and escalated; it
must not trigger another open-ended prose-repair cycle.

## Rationale

The current `design.md` is 3,968 lines. Repeated inherited semantics, API sketches,
checker descriptions and review-history prose have made a local correction likely to
leave stale restatements elsewhere. A full rewrite would introduce a far larger,
unreviewable delta and risk reopening approved feature semantics or diverging from
the Task 1–3 implementation checkpoint. A single structured contract plus a
deterministic consistency check makes the execution details reviewable without
changing the product architecture.

## Reuse-before-build assessment

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse the repository's feature-control, context-projection, semantic-guardrail,
    Make-target and standard-library testing conventions.
  - Reuse the existing FEATURE-0018 static architecture-checker approach and its
    positive/negative fixture pattern.
  - Reuse FEATURE-0013 only for its approved decision/audit carrier contracts and
    validation; no FEATURE-0013 service is created or extended.
- Sovrunn-owned responsibility summary:
  - FEATURE-0018 owns only its design-spec consistency tooling and deterministic,
    local, non-durable publication-mechanics specification.
- Non-goals summary:
  - No runtime IAM capability, external integration, persistence, provider action,
    resource schema, policy behavior, public error, API route, event taxonomy,
    FEATURE-0013 change or future-feature implementation.

## Phase impact

- Current phase allowed? Yes
- If not current phase, target phase: Not applicable
- Current phase boundary impact:
  - Documentation and deterministic local validation only. The change adds no
    external I/O, durable persistence or execution-target effect.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting decisions, if any:
  - None. Existing duplicated or contradictory design restatements are stale
    transcriptions, not competing accepted decisions.
- Resolution required:
  - Human approval of this handoff, then Kiro reconciliation. No DEC, RFC, ACR or
    architecture-baseline update is required.

## Required action

- Update Kiro design.md
- Create/update a FEATURE-0018 Design Mechanics Contract
- Update feature gate/checks
- Update controlled design-review context and approval hashes as required
- Update Kiro tasks.md only if a design-contract file path or verification command
  must be reflected; do not alter task scope, ordering or implementation semantics
- Do not update requirements.md
- Do not modify Go runtime code

## Impacted files

Kiro must inspect and update only the planning, validation and context artifacts
needed to apply this handoff, including:

- `docs/reviews/architecture-decision-handoffs/ADH-2026-075-feature-0018-design-normalization-and-contract-check.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/design.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/design-mechanics-contract.json`
- `.kiro/specs/governance-iam-approval-exception-foundation/tasks.md` only if needed
  for the contract's required verification command or path
- `.automation/features/FEATURE-0018.control.json`
- `.automation/context-projections/FEATURE-0018.design-review.spec.json`
- `.automation/context-projections/FEATURE-0018.design-review.md`
- `.automation/context-projections/FEATURE-0018.design-review.manifest.json`
- `scripts/feature0018-design-contract-check.py`
- `tests/feature_factory/test_feature0018_design_contract_check.py`
- `Makefile` only to add the named checker target

No `internal/govaccess/` implementation path is in scope under this handoff.

## Impacted features

- FEATURE-0018: design normalization and pre-review validation only.
- FEATURE-0012, FEATURE-0013, FEATURE-0016, FEATURE-0017 and all future features:
  no semantic or implementation change.

## Acceptance criteria for Kiro update

- [ ] Confirm no product-semantic change and no FEATURE-0013 ownership or contract
  expansion.
- [ ] Create the hash-bound structured Design Mechanics Contract with every required
  mechanics category and no placeholder types or owners.
- [ ] Make `design.md` cite that contract and remove redundant implementation-mechanics
  restatements without removing approved architecture or traceability assertions.
- [ ] Implement deterministic contract validation and its positive/negative tests.
- [ ] Add a named Make target and run it before the design semantic review.
- [ ] Keep the checker scoped to repository-owned edges and test-only seams.
- [ ] Run the checker, `./scripts/kiro-semantic-check.py --feature FEATURE-0018 --stage design`,
  `make vs000-contract-check`, and `make phase2r-drift-check` successfully.
- [ ] Produce a classified, read-only fixpoint-audit report.
- [ ] Refresh controlled review context and record the required fresh human approvals.
- [ ] Preserve Task 1–3 checkpoints and do not permit Cursor Task 4+ until the
  normalized design and task plan are reapproved.

## Explicit instructions to Kiro

- Do not introduce new FEATURE-0018 product behavior while normalizing documents.
- Do not change requirements, ADH-2026-070 through ADH-2026-074, DEC-0060, the
  architecture baseline, routes, writers, resource schemas, errors, events,
  conformance semantics, dependency ownership, or Phase 2R scope.
- Do not regenerate the whole design from an inferred architecture. Preserve its
  approved content and replace only duplicated executable-mechanics restatements
  with exact contract references.
- Do not modify `internal/govaccess/` code, provider adapters, database behavior or
  external integrations.
- Stop with `ARCHITECTURE_DECISION_REQUIRED` if the contract reveals a genuine
  semantic conflict instead of resolving it by inference.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-09-05
- Notes: Approved as a narrow FEATURE-0018 design-normalization and deterministic
  pre-review validation correction. No product-semantic or runtime implementation
  change is authorized by this handoff.
