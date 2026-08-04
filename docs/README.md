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

## Phase 2R Canonical Model Source-of-Truth

For Phase 2R and Phase 3 development, AI agents must load these before generating features:

- `docs/architecture/canonical/sovrunn-finalized-data-model.md` — canonical semantic model
- `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md` — contract catalog
- `docs/architecture/reference-flows/postgresql-end-to-end.md` — PostgreSQL reference flow
- `docs/phase2/PHASE2R_REBASELINE.md` — Phase 2R feature authority
- `docs/architecture/vertical-slices/VS-000-core-skeleton.md` — cross-feature acceptance
- `docs/architecture/vertical-slices/VS-000-contract-specification.md` — exact Slice 0 integration rules
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` — machine-readable Slice 0 schemas, writers, states, errors and conformance
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md` — Slice 0 repository traceability
- `docs/architecture/FEATURE-0015-canonical-cloud-model-and-alpha-migration-foundation.md` — closed FEATURE-0015 architecture boundary
- `docs/features/FEATURE-0015-canonical-cloud-model-and-alpha-migration-foundation.md` — executable FEATURE-0015 scope before Kiro stage generation
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/phase2/PHASE2_ACCEPTANCE_GATES.md`
- `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
- `docs/architecture/api-resource-standard.md`
- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- `docs/architecture/development-phases.md`

## Superseded Documents (Retained as History)

The following remain for audit trail but are not active authorities:

- `docs/architecture/provider-neutral-resource-model.md` — alpha model, superseded by canonical model
- Any reference to ResourcePool, ProviderCapability, or generic Provider as active targets

## Architecture Operating System

Sovrunn architecture is evolved through a durable Architecture Operating System.

Key files:

- `docs/context/ARCHITECTURE_VERSION.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/SOVRUNN_CONTEXT_PACK.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/context/CURRENT_DECISION_SUMMARY.md`
- `docs/context/OPEN_QUESTIONS.md`
- `docs/context/CHATGPT_ARCHITECTURE_SESSION_PROMPT.md`
- `docs/governance/ARCHITECTURE_CHANGE_CONTROL.md`
- `docs/governance/ARCHITECTURE_OWNERSHIP.md`
- `docs/governance/REVIEW_GATES.md`
- `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md`
- `docs/traceability/DECISION_TRACEABILITY_MATRIX.md`

## Handoff and Diagrams

- `docs/templates/ARCHITECTURE_DECISION_HANDOFF.md`
- `docs/prompts/chatgpt/architecture-decision-handoff.prompt.md`
- `docs/prompts/kiro/architecture-update.prompt.md`
- `docs/reviews/architecture-decision-handoffs/README.md`
- `docs/diagrams/structurizr/workspace.dsl`
- `docs/diagrams/structurizr/README.md`

## Roadmap and Phase Context

- `docs/architecture/development-phases.md`
- `docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md`
- `docs/features/FEATURE_INDEX.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/phase2/PHASE2R_REBASELINE.md`

## Generated Artifacts Policy

Generated artifacts such as `site/`, `docs/generated-prompts/`, generated context packs, logs, and zip archives are intentionally ignored and must not be treated as architecture source of truth.
