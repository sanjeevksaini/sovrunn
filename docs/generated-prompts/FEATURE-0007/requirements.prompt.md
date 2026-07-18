Feature name:
FEATURE-0007 Plugin and Capability Registry

## Model recommendation

Tool: kiro
Profile: requirements

Use the first available model from this prioritized list:

1. Architecture-heavy (`claude-opus-4.8`), effort: high — Use for requirements because this stage defines scope, constraints, goals, non-goals, and platform boundaries.
2. Fallback (`claude-sonnet-4.5`), effort: medium — Use if Opus is unavailable; keep the prompt strict and bounded.

If the first model is unavailable, use the next available model in priority order.
Do not use a fourth model unless all three are unavailable or blocked by governance.

At the end of execution, include exactly this report:

```text
Model Execution Report:
- Tool: kiro
- Stage or task: <stage/task id>
- Recommended priority list: <copy the three model labels from this prompt>
- Selected model: <actual model selected>
- Effort/reasoning setting: <actual setting if visible>
- Fallback used: yes/no
- Fallback reason: <unavailable/cost/latency/user override/none>
```


Start with requirements.md only.

Do not generate design.md yet.
Do not generate tasks.md yet.
Do not implement code.
Do not modify source files except the requirements file under:
.kiro/specs/plugin-capability-registry/requirements.md

Tool-output safety constraints:
- Use fs_write only for chunks of 50 lines or fewer.
- For files longer than 50 lines, create the file with fs_write using the first chunk, then use fs_append in chunks of 50 lines or fewer.
- Do not write the entire requirements.md in one fs_write call.
- Do not use one very large str_replace edit.
- Split content into logical sections.
- Write one section at a time.
- After writing, read the file back and verify it is complete.

Context:
Sovrunn is an AI-first sovereign cloud-native PaaS platform.
This feature belongs on branch feature-0007-plugin-capability-registry and should be scoped to Phase 1 unless explicitly stated otherwise.

Use these repo context files:
- AGENTS.md
- README.md
- docs/foundation/constitution.md
- docs/decisions/DECISION_INDEX.md
- docs/glossary.md
- docs/features/FEATURE_SEQUENCE.md
- docs/resource-specs/RESOURCE_MODEL_PHASE1.md
- docs/api/API_CONTRACT_PHASE1.md
- docs/engineering/ai-context-loading-standard.md
- docs/engineering/go-coding-guardrails.md

Requirements must include:
1. Introduction
2. Glossary if new concepts are introduced
3. User stories
4. Acceptance criteria
5. Non-goals
6. Edge cases
7. Security/privacy requirements
8. Compatibility with already completed Phase 1 features
9. Design questions to resolve later in design.md

Keep requirements concise, precise, implementation-aware, phase-scoped, and free of scope creep.

Do not implement code.
Do not generate design.md.
Do not generate tasks.md.
