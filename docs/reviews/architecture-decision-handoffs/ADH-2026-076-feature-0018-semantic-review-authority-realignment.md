# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-076
- Date: 2026-09-05
- Source discussion: ChatGPT Project / Codex
- Related feature: FEATURE-0018 — Governance, IAM, Approval and Exception Foundation
- Related phase: Phase 2R
- Author: ChatGPT / Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

FEATURE-0018 semantic-review authority realignment

## Summary

FEATURE-0018's product architecture, requirements, routes, ownership, security controls and Phase 2R boundary are closed. ADH-2026-075 establishes a hash-bound, machine-checkable Design Mechanics Contract for Go-level package, API, locking, factory, transaction and checker mechanics. This handoff realigns design review: deterministic checks and implementation tests are the authority for those mechanics; the LLM review is the authority for semantic architecture coherence, security, ownership, traceability and scope. This prevents repeated prose-level rediscovery of compiler-level details without relaxing any security or verification gate.

## Classification

- Correction

## Existing approved baseline

DEC-0060 and ADH-2026-070 through ADH-2026-075 remain unchanged. In particular, ADH-2026-072 assigns FEATURE-0018-local deterministic, non-durable publication mechanics to FEATURE-0018 and FEATURE-0013 only shared decision/audit contracts and validation. ADH-2026-075 makes a hash-bound Design Mechanics Contract and its deterministic checker the authoritative representation of those mechanics.

The latest approved product boundary remains: no provider-native IAM manipulation, no real IdP, no external I/O, no durable persistence and no execution-target effect. Task 1–3 checkpoints are complete; Cursor Task 4+ remains paused pending the renewed design, task and executable-plan approvals.

Relevant baseline references:

- `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `docs/decisions/DEC-0060-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-072-feature-0018-design-authority-and-audited-evaluation-clarification.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-075-feature-0018-design-normalization-and-contract-check.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/requirements.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/design.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/design-mechanics-contract.json`

## Decision or proposed decision

### 1. Separate semantic review from deterministic mechanics proof

The FEATURE-0018 Design Mechanics Contract and its deterministic checker are the sole authority for implementation-level mechanics: repository-owned package edges, named Go 1.22 API shapes, factory/caller confinement, `StateEditor` permissions, transaction state machines, lock order, identity issuance, test seams, static-checker fixtures and non-runtime gate evidence.

The independent LLM design review must instead assess only:

- conformance of FEATURE-0018 product semantics to approved ADHs, DEC-0060 and requirements;
- canonical ownership, writer boundaries, security controls, audit/decision semantics, authorization composition and Phase 2R exclusions;
- literal route, action and resource traceability;
- contract hash binding, deterministic-check result and whether the contract stays within its declared non-product-semantic authority; and
- evidence that required conformance, test and feature-gate ownership remains present.

### 2. Review finding rules

An LLM reviewer may raise `DESIGN_EXECUTABILITY` only when it identifies one of the following with an exact citation:

- a failed or absent deterministic Design Mechanics Contract check;
- a contract field missing a required authority, owner, test, fixture, result or proof artifact specified by ADH-2026-075;
- a conflict between the contract and approved semantic authority; or
- a product-semantic assertion leaked into the mechanics contract.

An LLM reviewer must not reopen a closed mechanical representation merely because it would choose a different Go type, factory shape, method name, lock implementation or test seam when the checker passes and the contract satisfies ADH-2026-075. It must not require duplicate mechanics prose in `design.md` as corroboration of the contract.

Every finding must be classified as one of `STALE_TRANSCRIPTION`, `DESIGN_EXECUTABILITY`, `REQUIREMENT_GAP`, `ARCHITECTURE_CLARIFICATION_REQUIRED`, or `OUT_OF_SCOPE`, and must cite the contradictory controlling authority before reopening a closed decision.

### 3. Hash-bound semantic review package

Kiro must produce a controlled FEATURE-0018 semantic-review projection within the existing review budget. It must contain:

- the full approved requirements projection and its approval evidence;
- the feature identity, semantic design sections, route/resource/ownership and traceability ledgers required for semantic review;
- the full hash-bound Design Mechanics Contract, its checker result and its contract hash binding to the raw `design.md`; and
- exact ADH/DEC excerpts needed to judge semantic authority.

The projection may omit non-normative historical revision narration, execution reports and contract-governed mechanical restatements from the LLM input only after the deterministic checker verifies that the raw design cites—not overrides—the contract. The raw design and contract hashes must remain visible to the reviewer.

### 4. Bounded exit path

Before semantic review, the following all must pass:

```text
make feature-0018-design-contract-check
./scripts/kiro-semantic-check.py --feature FEATURE-0018 --stage design
make vs000-contract-check
make phase2r-drift-check
make ff-feature-factory-test
```

One fresh semantic-review epoch may have at most two attempts. A finding outside the review authority defined here is recorded as `OUT_OF_SCOPE` or escalated as `ARCHITECTURE_CLARIFICATION_REQUIRED`; it must not begin another open-ended design revision cycle.

## Rationale

The current design reached nearly 4,000 lines because compiler-level mechanics were repeated in components, decisions, protocol descriptions, fixtures, traceability, self-verification and revision history. LLM review of each prose restatement has repeatedly discovered a new representation inconsistency even though FEATURE-0018's IAM architecture is unchanged. Deterministic validation is more reliable, faster and cheaper for closed mechanics; expert review adds the most value when it focuses on semantic security and ownership outcomes.

## Reuse-before-build assessment

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse the ADH-2026-075 Design Mechanics Contract, its deterministic checker, controlled context projections and existing feature-factory review routing.
  - Reuse standard Go compiler, unit-test and race-test evidence during Cursor implementation rather than reproducing compiler judgments in prose review.
- Sovrunn-owned responsibility summary:
  - FEATURE-0018 owns its bounded contract/checker evidence and semantic-review projection only.
- Non-goals summary:
  - No change to IAM behavior, resource schema, action, route, writer, event, public error, FEATURE-0013 contract, external integration, durable persistence or execution-target capability.

## Phase impact

- Current phase allowed? Yes
- If not current phase, target phase: Not applicable
- Current phase boundary impact:
  - Governance and deterministic design-validation change only; no external I/O or runtime platform effect is introduced.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting decisions, if any:
  - None. This clarifies the evidence authority introduced by ADH-2026-075; it does not replace any FEATURE-0018 product-semantic decision.
- Resolution required:
  - Human approval, then Kiro context and review-route reconciliation. No DEC, RFC, ACR or baseline update is required.

## Required action

- Update Kiro design.md only to remove stale review-authority prose if necessary
- Update controlled design-review context/projection and manifest
- Update generic review routing/checks only where needed to enforce this handoff
- Update Kiro tasks.md only if a required verification command must be retained
- Update feature gate/checks
- Do not update requirements.md
- Do not modify FEATURE-0018 runtime Go code

## Impacted files

- `docs/reviews/architecture-decision-handoffs/ADH-2026-076-feature-0018-semantic-review-authority-realignment.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/design.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/design-mechanics-contract.json`
- `.automation/features/FEATURE-0018.control.json`
- `.automation/context-projections/FEATURE-0018.design-review.spec.json`
- `.automation/context-projections/FEATURE-0018.design-review.md`
- `.automation/context-projections/FEATURE-0018.design-review.manifest.json`
- `scripts/feature0018-design-contract-check.py`
- `scripts/feature-control.py` and `scripts/generic-review-prompt.py` only if required for semantic-review context selection and enforcement
- `tests/feature_factory/` only for the corresponding deterministic routing/checker tests
- `Makefile` only if the existing checker target requires review-route wiring

## Impacted features

- FEATURE-0018: review authority, context and deterministic-evidence normalization only.
- FEATURE-0012, FEATURE-0013, FEATURE-0016, FEATURE-0017 and all future features: no semantic or implementation change.

## Acceptance criteria for Kiro update

- [ ] Preserve all approved FEATURE-0018 product semantics and completed Task 1–3 checkpoints.
- [ ] Keep the Design Mechanics Contract hash-bound to raw `design.md` and verify it before every semantic review.
- [ ] Enforce the stated reviewer finding rules in the review route/prompt.
- [ ] Produce a hash-bound semantic-review projection within the existing 650 KB budget without increasing it or suppressing required semantic authority.
- [ ] Ensure raw design mechanics cannot override the contract; fail closed if the raw design makes a conflicting mechanics claim.
- [ ] Run every bounded-exit command successfully before semantic review.
- [ ] Record no more than two attempts in the fresh semantic-review epoch.
- [ ] Do not modify requirements, runtime Go code, feature semantics, routes, writers, errors, events, dependencies or Phase 2R scope.

## Explicit instructions to Kiro

- Do not use this handoff to waive a failed deterministic check, a failed test, a contract/hash mismatch or a genuine semantic security contradiction.
- Do not regenerate FEATURE-0018 architecture, requirements, design or tasks from inference.
- Do not replace deterministic evidence with an LLM assertion.
- Do not add a new product architecture decision.
- Stop with `ARCHITECTURE_DECISION_REQUIRED` if applying this review-authority split requires changing approved IAM behavior or the Phase 2R boundary.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-09-05
- Notes: Approved as a FEATURE-0018 review-authority correction. Deterministic
  contract evidence remains mandatory; no product-semantic or runtime change is
  authorized by this handoff.
