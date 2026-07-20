# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-EXAMPLE-001
- Date: YYYY-MM-DD
- Source discussion: ChatGPT Project
- Related feature: FEATURE-0017
- Related phase: Phase 2
- Author: ChatGPT Architecture Governor
- Human approver: Example only
- Approval status: Proposed

## Decision title

FEATURE-0017 PolicyEngineAdapter scope

## Summary

FEATURE-0017 should define the policy evaluation abstraction and placeholders for OPA and Cedar adapters. It should not implement full OPA or Cedar integration during Phase 2.

## Classification

Clarification

## Existing approved baseline

Phase 2 is model, decision, audit, adapter-boundary, plugin-taxonomy, and placement-simulation foundation only. Full OPA/Cedar integration is deferred beyond Phase 2.

## Decision or proposed decision

FEATURE-0017 defines `PolicyEngineAdapter`, `PolicyEvaluationRequest`, `PolicyEvaluationResult`, `PolicyInput`, `PolicyContext`, `PolicyBundleRef`, and adapter placeholders for OPA and Cedar. It does not implement full OPA/Cedar policy execution in Phase 2.

## Rationale

This preserves reuse-before-build while preventing premature coupling to a policy language before Sovrunn policy inputs and decision outputs stabilize.

## Reuse-before-build assessment

- Decision: Wrap
- Mature reusable options considered:
  - OPA/Rego
  - Cedar
- Sovrunn-owned responsibility:
  - policy context,
  - decision result contract,
  - audit linkage,
  - adapter boundary.
- Non-goals:
  - custom policy engine,
  - full OPA/Cedar integration in Phase 2.

## Phase impact

- Current phase allowed? Yes
- If not current phase, target phase: N/A
- Current phase boundary impact:
  - Keeps Phase 2 as abstraction/foundation only.

## Conflict check

- Conflicts with accepted DEC/RFC? No
- Conflicting decisions, if any:
  - None
- Resolution required:
  - Update architecture docs and Kiro feature spec.

## Required action

- Update architecture doc
- Update Kiro requirements.md
- Update Kiro design.md
- Update Kiro tasks.md
- Update traceability matrix

## Impacted files

- `docs/architecture/policy-evaluation-abstraction.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md`
- `.kiro/specs/policy-evaluation-abstraction/requirements.md`
- `.kiro/specs/policy-evaluation-abstraction/design.md`
- `.kiro/specs/policy-evaluation-abstraction/tasks.md`

## Impacted features

- FEATURE-0017: clarifies Phase 2 scope and implementation boundary.

## Acceptance criteria for Kiro update

- [ ] Handoff validated against Architecture Operating System files
- [ ] No unapproved baseline change
- [ ] Phase scope respected
- [ ] Reuse-before-build section preserved
- [ ] Feature requirements/design/tasks updated if applicable
- [ ] Traceability matrix updated if applicable

## Explicit instructions to Kiro

- Do not introduce new decisions beyond this handoff.
- Do not update `CURRENT_ARCHITECTURE_BASELINE.md` unless required by validation.
- Do not implement full OPA/Cedar integration.
- Do not modify Go code.

## Human approval

- Approval status: Proposed
- Approved by:
- Date:
- Notes: Example file only. Do not treat as approved architecture.
