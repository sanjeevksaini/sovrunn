Generate `{{DESIGN_PATH}}` for {{FEATURE_ID}} — {{FEATURE_TITLE}}.

Stage: Design. Modify that file only. Do not modify requirements, tasks,
architecture, source code, schemas, prompts, validators, or automation.

{{MODEL_RECOMMENDATIONS}}

{{CONTROL_FRAGMENT}}

Design translates approved requirements into implementable representation. It
may choose only mechanics explicitly delegated by architecture or requirements.
It must not add product semantics, transfer ownership, or design excluded
adjacent-feature mechanisms.

Rules:

- Map every requirement to a design disposition or an explicit contract-only
  disposition. Do not silently drop requirements.
- Preserve the approved REQ and AC meanings exactly. Add a section titled
  exactly `Canonical coverage ledger` in which every approved REQ and every
  approved AC appears exactly once with its design disposition. Do not create,
  renumber, merge, split, reinterpret, or omit an ID.
- Reuse previous-feature types, grammar, errors, limits, and test helpers from
  their canonical owners; do not fork them.
- Follow the repository's Go package direction, minimal-dependency policy,
  deterministic-validation rules, stable-error conventions, security rules,
  observability rules, and compatibility requirements.
- State files and packages only after verifying the live repository convention.
- Define public contracts, validation order, failure behavior, versioning,
  compatibility, security, observability, testing, and migration impact where
  applicable.
- Separate `IMPLEMENT`, `CONTRACT_ONLY/NO_TASK`, and `EXCLUDED` elements.
- Do not invent repositories, persistence, controllers, adapters, handlers, or
  runtime behavior unless approved requirements require them.
- Prefer references to inherited standards over copying their contents.
- Resolve all design questions. If a semantic choice remains, stop rather than
  delegating it to tasks or Cursor.
- Do not import downstream fields, metrics, side effects, controllers,
  adapter/plugin execution, or runtime proofs. A referenced `VS0-CF-*` case
  retains its exact registry owner and semantics.

Required sections:

1. Identity, stage, inputs, and closed boundary
2. Resolved design decisions
3. Components and repository paths
4. Data/API representation
5. Validation and deterministic error behavior
6. Security, privacy, observability, compatibility, and versioning
7. Test and conformance strategy
8. Implementation classification ledger
9. Requirement/decision/risk traceability
10. Non-goals, absence ledger, and unresolved report

Before completion, read the entire file and verify that it contains no
unapproved semantic choice and no mechanism owned by an excluded feature. The
automation will run a deterministic semantic check before review; preemptively
satisfy every rule above before returning COMPLETE.

Finish with exactly one receipt:

`STAGE_STATUS: COMPLETE`

or

`STAGE_STATUS: BLOCKED <APPROVED_STOP_CONDITION>`
