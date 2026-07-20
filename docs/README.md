# Sovrunn Documentation

Start here:

- foundation/vision.md
- foundation/constitution.md
- decisions/DECISION_INDEX.md
- glossary.md
- features/FEATURE_SEQUENCE.md
- resource-specs/RESOURCE_MODEL_PHASE1.md
- api/API_CONTRACT_PHASE1.md
- engineering/ai-context-loading-standard.md
- engineering/go-coding-guardrails.md

## Phase 2 Source-of-Truth Additions

For Phase 2 and Phase 3 development, AI agents must load these before generating features:

- `docs/architecture/development-phases.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/phase2/PHASE2_ACCEPTANCE_GATES.md`
- `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
- `docs/architecture/api-resource-standard.md`
- `docs/architecture/decision-and-audit-standard.md`
- `docs/architecture/provider-neutral-resource-model.md`
- `docs/architecture/policy-evaluation-abstraction.md`
- `docs/architecture/placement-decision-engine.md`
- `docs/architecture/plugin-taxonomy-and-boundaries.md`
- `docs/architecture/adapter-boundary-model.md`
- `docs/mvp/MVP_001_GOVERNED_POSTGRESQL_PAAS.md`

## Architecture Operating System

The Architecture Operating System layer controls long-term Sovrunn architecture evolution.

Important folders:

- `context/` — current baseline, context pack, session prompt, checkpoints.
- `governance/` — change control, ownership, review gates.
- `traceability/` — feature and decision traceability matrices.
- `templates/` — architecture change, DEC, RFC, review templates.
- `reviews/monthly/` — recurring architecture baseline reviews.
- `reviews/phase-gates/` — phase readiness and closeout reviews.
- `reviews/feature-gates/` — feature approval reviews.

Chat history is not source of truth. Approved repo docs are source of truth.

## Architecture Decision Handoff

Sovrunn separates architecture discussion from repository updates.

- ChatGPT Project produces an Architecture Decision Handoff.
- Human approves the handoff.
- Kiro validates and applies approved handoff to docs/specs.
- Cursor implements only from approved Kiro tasks.

Important files:

- `docs/templates/ARCHITECTURE_DECISION_HANDOFF.md`
- `docs/prompts/chatgpt/architecture-decision-handoff.prompt.md`
- `docs/prompts/kiro/architecture-update.prompt.md`
- `docs/reviews/architecture-decision-handoffs/README.md`

## Architecture diagrams

Structurizr DSL lives under:

```text
docs/diagrams/structurizr/workspace.dsl
```

Use it for durable C4 architecture views. Markdown docs explain architecture decisions; Structurizr DSL visualizes the approved architecture model.

Commands:

```bash
make structurizr-check
make structurizr-lite
```

## Canonical Roadmap and Feature Index

- Canonical all-phase roadmap: `docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md`
- Feature ID to Kiro slug mapping: `docs/features/FEATURE_INDEX.md`
- `docs/features/FEATURE_ROADMAP_ALL_PHASES.md` is intentionally only a pointer to avoid duplicated roadmap sources.

## Generated Artifact Policy

These paths are generated artifacts and must not be used as source of truth:

```text
site/
docs/generated-prompts/
docs/context/SOVRUNN_CONTEXT_PACK.generated.md
.automation/generated-prompts/
.automation/logs/
.automation/reviews/
*.zip
```
