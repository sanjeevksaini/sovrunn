# Architecture Decision Handoff

> Purpose: Standard contract between ChatGPT architecture discussion and Kiro repo/spec update.
>
> This file is not an approval by itself. It captures a candidate or approved architecture outcome so Kiro can validate and apply it against the Sovrunn Architecture Operating System.

## Metadata

- Handoff ID: ADH-YYYY-NNN
- Date: YYYY-MM-DD
- Source discussion: ChatGPT Project / other
- Related feature: FEATURE-XXXX
- Related phase: Phase X
- Author: <name/tool>
- Human approver: <name>
- Approval status: Proposed / Approved / Rejected / Deferred

## Decision title

<Short title>

## Summary

<One paragraph summary of the decision or proposed change.>

## Classification

Select one:

- Explanation
- Clarification
- Extension
- Correction
- Replacement
- New decision

## Existing approved baseline

State what the current approved Sovrunn baseline already says.

Relevant baseline references:

- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/decisions/DECISION_INDEX.md`
- Related DEC/RFC: <ids>

## Decision or proposed decision

State the exact decision in implementation-neutral terms.

## Rationale

Explain why this is recommended.

## Reuse-before-build assessment

- Decision: Reuse / Wrap / Extend / Build
- Mature reusable options considered:
  - <option>
- Sovrunn-owned responsibility:
  - <responsibility>
- Non-goals:
  - <non-goal>

## Phase impact

- Current phase allowed? Yes / No
- If not current phase, target phase: Phase X
- Current phase boundary impact:
  - <impact>

## Conflict check

- Conflicts with accepted DEC/RFC? Yes / No
- Conflicting decisions, if any:
  - <id>
- Resolution required:
  - None / ACR / DEC / RFC / Baseline update

## Required action

Select all that apply:

- No repo change
- Update architecture doc
- Update phase scope doc
- Update roadmap placeholder
- Create/update Open Question
- Create Architecture Change Request
- Create/update DEC
- Create/update RFC
- Update Kiro requirements.md
- Update Kiro design.md
- Update Kiro tasks.md
- Update traceability matrix
- Update feature gate/checks

## Impacted files

List exact files Kiro should update or inspect.

- `...`

## Impacted features

- FEATURE-XXXX: <impact>

## Acceptance criteria for Kiro update

- [ ] Handoff validated against Architecture Operating System files
- [ ] No unapproved baseline change
- [ ] Phase scope respected
- [ ] Reuse-before-build section preserved
- [ ] Required DEC/RFC/ACR created or updated if needed
- [ ] Feature requirements/design/tasks updated if applicable
- [ ] Traceability matrix updated if applicable
- [ ] Open questions updated if deferred

## Explicit instructions to Kiro

- Do not introduce new decisions beyond this handoff.
- Do not update `CURRENT_ARCHITECTURE_BASELINE.md` unless this handoff is approved and requires a baseline update.
- Do not move future-phase implementation into the current phase.
- Do not modify Go code unless the handoff explicitly requires implementation work.

## Human approval

- Approval status: Proposed / Approved / Rejected / Deferred
- Approved by:
- Date:
- Notes:
