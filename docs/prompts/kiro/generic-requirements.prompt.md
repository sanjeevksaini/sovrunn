Generate `{{REQUIREMENTS_PATH}}` for {{FEATURE_ID}} — {{FEATURE_TITLE}}.

Stage: Requirements. Modify that file only. Do not create or modify design,
tasks, architecture, schemas, source code, prompts, validators, or automation.

{{MODEL_RECOMMENDATIONS}}

{{CONTROL_FRAGMENT}}

Requirements own observable behavior: intent, actors, scenarios, invariants,
validation outcomes, security/privacy, compatibility, non-goals, and acceptance
criteria. They do not own packages, files, routes, storage, algorithms,
libraries, internal interfaces, or task decomposition.

Rules:

- Translate approved architecture decisions; do not reopen or extend them.
- Treat the current feature file's REQ table and acceptance-criteria table as
  immutable semantic ledgers. Preserve every ID-to-meaning mapping one-to-one:
  never repurpose, renumber, merge, split, omit, or paraphrase a ledger row.
- Add a section titled exactly `Canonical requirement ledger` that copies every
  approved REQ table row, with the same columns and cell text, exactly once.
  Add `Canonical acceptance ledger` with the same exact-copy rule for every AC.
  Detailed scenarios may expand those rows but may not change their meaning.
- Consume previous-feature contracts by exact reference. Never copy, rename,
  specialize, or transfer their ownership unless the architecture explicitly
  authorizes it.
- Keep every normative statement under exactly one canonical owner.
- Use exact stable requirement, decision, and risk IDs where the architecture
  defines them. Reference shared obligations instead of duplicating prose.
- Give every approved REQ exactly one normative detail heading, in approved
  order. Map every approved AC explicitly outside the canonical ledger in an
  acceptance-scenario or coverage mapping; a completeness claim without that
  mapping and its exact ledger row is a failure.
- Add `Exact conformance semantics ledger`. For every referenced `VS0-CF-*`
  ID, copy exactly these registry fields: ID, owner, inputs, expectedState,
  expectedError, expectedSideEffects, and gate. A conformance case proves only
  its registered scenario. Never use a conformance case as a substitute for a
  state-machine ID or reinterpret a downstream case as this feature's behavior.
- Preserve all adjacent-feature exclusions in enforceable non-goals and
  negative acceptance cases.
- Name prohibited stale concepts only inside an explicit `Non-goals` or
  `Out of Scope` section. Active requirements and acceptance cases must refer
  to that exclusion section without repeating prohibited names.
- Do not import a downstream feature's status field, metric, side effect,
  controller, adapter/plugin execution, or runtime proof merely because it is
  visible in a shared Slice 0 document.
- Record implementation questions only when architecture explicitly delegates
  them to design. A semantic or ownership question is not a design question.
- Do not create requirements solely to satisfy a generic template when the
  feature architecture marks the concern not applicable.
- Prefer concise tables and references over repeated narrative.

Required sections:

1. Identity and stage
2. Purpose and use cases
3. Terms introduced by this feature
4. Normative requirements and acceptance scenarios
5. Security, privacy, compatibility, and operational requirements
6. Edge cases
7. Non-goals and adjacent-feature exclusions
8. Architecture/decision/risk traceability by exact ID
9. Design questions explicitly delegated by architecture
10. Completeness and unresolved-decision report

Before completion, read the entire file and verify that no design choice,
future-feature semantic, duplicated inherited contract, or unresolved
architecture decision was introduced. The automation will run a deterministic
repository semantic check before review; preemptively satisfy every rule above
before returning COMPLETE.

Finish with exactly one receipt:

`STAGE_STATUS: COMPLETE`

or

`STAGE_STATUS: BLOCKED <APPROVED_STOP_CONDITION>`
