Feature name:
{{FEATURE_ID}} {{FEATURE_TITLE}}

{{MODEL_RECOMMENDATIONS}}

Start with requirements.md only.

Do not generate design.md yet.
Do not generate tasks.md yet.
Do not implement code.
Do not modify source files except the requirements file under:
{{SPEC_PATH}}/requirements.md

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
This feature belongs on branch {{FEATURE_BRANCH}} and should be scoped to Phase 1 unless explicitly stated otherwise.

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

## Phase 2 Reuse and Drift Gates

Every FEATURE-0011-and-later feature must include a reuse assessment that
conforms to the canonical standard:

`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

Do not duplicate or redefine the assessment field schema in this document.
Populate the feature-level reuse summary and capability-level assessments
using the canonical fields and controlled vocabularies.

Architecture drift checks:

- no provider-specific hardcoding in core,
- no Kubernetes-only assumptions in core,
- no PostgreSQL lifecycle logic in core placement engine,
- no custom policy engine embedded in handlers,
- no raw secret storage,
- no customer-facing IaaS leakage,
- explainable decision object,
- defined audit behavior,
- preserved adapter boundaries.
