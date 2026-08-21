# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-062
- Date: 2026-08-21
- Source discussion: Codex task-plan review
- Related feature: FEATURE-0016
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Correct FEATURE-0016 task-scope inventory path classification

## Summary

This automation-only correction fixes a false-positive task-admission failure.
The current F0016 task-scope inventory rejects any occurrence of `internal/`,
`cmd/`, `api/schemas/`, or `tests/` in `tasks.md` because it treats documented
implementation paths as Kiro-generated outputs. The approved Tasks artifact is
required to enumerate exact future Cursor writable and test paths, so that rule
makes the mandatory inventory structurally incompatible with its own contract.

The correction distinguishes documentation of declared implementation inventory
from an out-of-bound Kiro write. It introduces no F0016 resource, route, field,
state, error, conformance meaning, implementation behavior, or phase change.

## Classification

Correction — automation boundary interpretation.

## Existing approved baseline

- The FEATURE-0016 control manifest correctly limits Kiro Tasks-stage writes to
  `.kiro/specs/adapter-boundary-and-executiontarget-qualification/tasks.md`.
- Its `forbidden_outputs` list correctly prohibits Kiro from creating source,
  schema, or test outputs in the Tasks stage.
- The approved task-plan format requires exact `Writable paths` and `Tests`
  inventory for later Cursor implementation.
- `scripts/feature-0016-task-scope-inventory.py` currently scans all text for
  forbidden-output prefixes and therefore rejects those required inventory
  references even though no Kiro output was created outside `tasks.md`.

## Decision

The F0016 task-scope inventory must validate, not reject, repository paths
declared inside the Tasks document's `Writable paths` or `Tests` sections.

1. The task-stage Kiro output boundary remains unchanged: only `tasks.md` is
   writable by Kiro.
2. References to `internal/`, `cmd/`, `api/schemas/`, and `tests/` within
   task-inventory fields are declarations of future Cursor scope, not Kiro
   outputs. They are permitted and must be syntactically validated as relative,
   repository-contained paths.
3. The inventory must fail closed for malformed paths, absolute paths,
   traversal (`..`), empty inventory entries, and any task lacking required
   writable/test inventory fields.
4. The inventory must still reject an actual Tasks-stage output outside the
   manifest's sole writable path through existing stage-boundary checks. This
   correction does not weaken that enforcement.

## Example

The following Task 11 declarations are valid plan inventory, not Kiro source
writes:

```text
internal/api/executiontarget_actions.go
internal/server/server.go
cmd/sovrunn-api/main.go
internal/server/server_test.go
```

Kiro may document those paths in `tasks.md`, but cannot create or edit them.
Cursor may later edit them only when the approved task prompt permits them.

## Required action

- Update the F0016 task-scope inventory script to parse and validate task
  inventory fields instead of textually rejecting their path prefixes.
- Add deterministic self-test coverage for accepted repository-relative
  inventory paths and rejected malformed/traversal/absolute paths.
- Update the F0016 control manifest only if needed to record this approved
  handoff; do not alter its Kiro Tasks-stage writable path or
  `forbidden_outputs` behavior.
- Do not modify requirements, design, tasks semantics, architecture authority,
  Go implementation, or formal models while applying this correction.

## Strict Kiro read/write boundary

| Path | Access | Purpose |
|---|---|---|
| `docs/reviews/architecture-decision-handoffs/ADH-2026-062-feature-0016-task-scope-inventory-path-classification-correction.md` | read | Controlling approved correction |
| `scripts/feature-0016-task-scope-inventory.py` | write | Parse/validate declared task inventory without treating it as Kiro output |
| `tests/feature_factory/test_feature_control.py` | write | Add deterministic parser/path-classification coverage if the existing feature-factory test location is applicable |
| `.automation/features/FEATURE-0016.control.json` | write | Add ADH-2026-062 only to `feature.handoffs`, if required by manifest provenance |

## Acceptance criteria

- [ ] Existing valid F0016 task inventory entries under `internal/`, `cmd/`,
  `api/schemas/`, and `tests/` pass the scope inventory check.
- [ ] Absolute, traversal, malformed, and empty declared inventory paths fail.
- [ ] The task-stage Kiro write boundary remains exactly `tasks.md`.
- [ ] No source file is created or edited by this architecture-update run.
- [ ] The task-scope check, semantic check, feature-control validation, and
  `git diff --check` pass.

## Explicit instructions to Kiro

- Apply only this automation-boundary correction.
- Do not weaken `forbidden_outputs`; distinguish declared plan inventory from
  actual Kiro outputs.
- Do not edit `tasks.md` while applying this handoff.
- Do not modify requirements, design, feature architecture, registry, contract
  specification, traceability, Go source, or formal models.
- If the checker cannot make this distinction without weakening actual
  stage-boundary enforcement, stop with `SECURITY_REVIEW_REQUIRED`.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-21
- Notes: Founder approved the bounded task-scope inventory path-classification
  correction before further F0016 task-plan review.
