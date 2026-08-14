Generate `{{TASKS_PATH}}` for {{FEATURE_ID}} — {{FEATURE_TITLE}}.

Stage: Tasks. Modify that file only. Do not modify requirements, design,
architecture, source code, schemas, prompts, validators, or automation.

{{MODEL_RECOMMENDATIONS}}

{{CONTROL_FRAGMENT}}

Tasks decompose the approved design; they do not make design or architecture
decisions. Stop if requirements and design do not determine an implementation.

Rules:

- Generate tasks only for design elements classified `IMPLEMENT`.
- Add a section titled exactly `Canonical coverage ledger`. Map every approved
  REQ and AC exactly once to a task or to the no-task ledger; preserve its
  approved meaning and never create, renumber, merge, split, reinterpret, or
  omit an ID.
- Target 6–12 independently testable vertical slices for a Tier A feature.
- Prefer a resource/behavior slice containing representation, validation and
  tests over separate all-types/all-schemas/all-tests phases.
- Sequence real dependencies; mark independent tasks explicitly.
- Every task must cite exact requirements, design sections/decisions,
  architecture decisions, and applicable risks.
- Every task must list exact writable files or directory prefixes, tests,
  verification commands, acceptance criteria, security/observability impact,
  exclusions, and an exact commit message.
- Use these exact field labels once in every task: `Writable paths:`, `Tests:`,
  `Verification commands:`, `Acceptance criteria:`, and `Commit message:`.
- No task may create a semantic choice, external dependency, production
  mechanism, or adjacent-feature behavior absent from the approved design.
- Add a no-task ledger for contract-only, invariant, governance, and excluded
  design elements instead of inventing implementation work.
- The final task performs full repository verification, boundary validation,
  artifact cleanup, and clean-tree confirmation.
- Prefer concise references; do not repeat full requirement/design prose.

Required sections:

1. Identity, stage, and execution rules
2. Dependency-ordered implementation tasks
3. No-task ledger
4. Requirement/design/decision/risk coverage ledger
5. Orphan and unresolved report

Before completion, verify that every `IMPLEMENT` element has exactly one task,
every task has an approved owner, and no task implements a non-goal. The
automation will run a deterministic semantic check before review; preemptively
satisfy every rule above before returning COMPLETE.

Finish with exactly one receipt:

`STAGE_STATUS: COMPLETE`

or

`STAGE_STATUS: BLOCKED <APPROVED_STOP_CONDITION>`
