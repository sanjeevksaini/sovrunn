# ChatGPT Patch Prompt

Use this prompt when asking ChatGPT for exact file-level changes.

## Prompt

You are producing an exact patch for Sovrunn Phase 1.

Use these as source of truth:

```text
AGENTS.md
docs/foundation/constitution.md
docs/decisions/DECISION_INDEX.md
docs/glossary.md
docs/features/FEATURE_SEQUENCE.md
docs/resource-specs/RESOURCE_MODEL_PHASE1.md
docs/api/API_CONTRACT_PHASE1.md
current feature file
```

Task:

```text
<describe exact bug/fix/change>
```

Rules:

```text
keep patch minimal
do not expand scope
do not introduce future features
do not rename canonical terms
do not add dependencies unless necessary
preserve tests or add tests
explain commands to run
```

Return files to change, exact code or diff, why the change is needed, tests to run, and expected result.
