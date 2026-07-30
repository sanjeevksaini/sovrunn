Implement {{FEATURE_ID}} Task {{TASK_ID}} only.

{{MODEL_RECOMMENDATIONS}}

{{CONTROL_FRAGMENT}}

## Approved task

{{TASK_BLOCK}}

## Writable paths

{{WRITABLE_PATHS}}

Do not create, modify, rename, or delete any other path. Begin by checking the
working tree and the live implementation conventions in the listed context.

## Go implementation guardrails

- Follow `docs/engineering/go-coding-guardrails.md`, the repository Go version,
  existing package direction, and standard-library-first dependency policy.
- Reuse existing FEATURE-0012 and previous-feature primitives from their
  canonical packages; do not fork shared metadata, references, conditions,
  validation, errors, limits, or conformance helpers.
- Keep APIs small, typed, deterministic, context-aware, race-safe, and easy to
  test. Validate at boundaries and return stable errors without leaking secrets
  or provider-native details.
- Keep handlers thin. Do not add persistence, controllers, adapters, provider
  clients, global registries, background work, or dependencies unless this task
  and approved design explicitly require them.
- Add happy-path, boundary, negative, and concurrency tests applicable to this
  task. Never weaken an existing test to make new code pass.
- Use structured observability fields where the approved design requires them;
  never log credentials, tokens, private keys, connection strings, or raw
  sensitive payloads.
- Do not edit architecture, requirements, design, tasks, prompts, validators,
  automation, or future-task files.
- Do not make a semantic or architecture decision. Stop with an approved stop
  condition if the task cannot be implemented exactly.

Run every verification command required by the task and feature control
manifest. Inspect `git diff` and confirm all paths are allowed.

Finish with a concise implementation report followed by exactly one receipt:

`TASK_STATUS: COMPLETE`

or

`TASK_STATUS: BLOCKED <APPROVED_STOP_CONDITION>`
