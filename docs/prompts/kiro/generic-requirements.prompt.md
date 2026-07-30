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
- Consume previous-feature contracts by exact reference. Never copy, rename,
  specialize, or transfer their ownership unless the architecture explicitly
  authorizes it.
- Keep every normative statement under exactly one canonical owner.
- Use exact stable requirement, decision, and risk IDs where the architecture
  defines them. Reference shared obligations instead of duplicating prose.
- Preserve all adjacent-feature exclusions in enforceable non-goals and
  negative acceptance cases.
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
architecture decision was introduced.

Finish with exactly one receipt:

`STAGE_STATUS: COMPLETE`

or

`STAGE_STATUS: BLOCKED <APPROVED_STOP_CONDITION>`
