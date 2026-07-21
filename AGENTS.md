# Sovrunn AI Agent Instructions

## Source of Truth

The Git repository and terminal verification commands are the source of truth.

AI tools may propose, generate, edit, or review code, but correctness is proven only by:

```bash
git status
git diff
make fmt
make test
make vet
```

## Tool Responsibilities

```text
Kiro      = architecture, specs, task planning, feature breakdown
Cursor    = code editing, refactoring, debugging, tests
ChatGPT   = deep architecture review, exact patches, debugging help
Terminal  = source of truth: build, test, run, Git, runtime behavior
```

No tool owns a separate architecture.

## Current Phase

```text
Phase 2: Reuse-First PaaS Fabric Foundation
```

Phase 0 and Phase 1 documents remain valid baseline records. Phase 2 extends them with reuse-first, adapter-first, provider-neutral, decision-first, audit-first, and plugin-taxonomy-first architecture.

## Current Implementation Rule

Implement one feature at a time.

Current Phase 2 feature order is authoritative in:

```text
docs/phase2/PHASE2_FEATURE_SEQUENCE.md
```

All-phase roadmap placeholders are maintained in:

```text
docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md
docs/features/FEATURE_INDEX.md
```

Future roadmap features may be referenced for scope awareness only. They must not be implemented during Phase 2 or Phase 3 unless a formal decision changes the phase boundary.

## Authoritative Context Priority

Context priority and classification are defined centrally in:

```text
docs/engineering/ai-context-loading-standard.md
```

Individual file front matter such as `ai_load_priority` is only a local hint.

If local metadata conflicts with `docs/engineering/ai-context-loading-standard.md`, the central standard wins.

## Canonical Ownership

If files conflict, canonical ownership wins:

```text
context loading
  docs/engineering/ai-context-loading-standard.md

Go coding
  docs/engineering/go-coding-guardrails.md

terminology
  docs/glossary.md

accepted decisions
  docs/decisions/DECISION_INDEX.md

API behavior
  docs/api/API_CONTRACT_PHASE1.md

resource model
  docs/resource-specs/RESOURCE_MODEL_PHASE1.md

feature scope
  current docs/features/FEATURE-xxxx file

phase roadmap
  docs/architecture/development-phases.md
  docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md
  docs/features/FEATURE_INDEX.md

master AI operating contract
  AGENTS.md
```

## Required Context Before Coding

Before implementing any Go feature, load:

```text
AGENTS.md
README.md
docs/foundation/constitution.md
docs/decisions/DECISION_INDEX.md
docs/glossary.md
docs/features/FEATURE_SEQUENCE.md
docs/phase2/PHASE2_FEATURE_SEQUENCE.md when working on Phase 2
docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md for scope awareness only
docs/resource-specs/RESOURCE_MODEL_PHASE1.md
docs/api/API_CONTRACT_PHASE1.md
docs/engineering/ai-context-loading-standard.md
docs/engineering/go-coding-guardrails.md
docs/architecture/controller-reconciliation-model.md
docs/architecture/observability-and-audit-baseline.md
current FEATURE-xxxx file
```

Do not load all repository files by default.

Do not load future feature files unless explicitly requested.

## AI Context Loading

AI agents must follow:

```text
docs/engineering/ai-context-loading-standard.md
```

Default task-specific loading model:

```text
ALWAYS
+ current FEATURE file
+ GO_IMPLEMENTATION when coding Go
+ role-specific prompt or steering files when needed
```

## Go Guardrails

All Go implementation must follow:

```text
docs/engineering/go-coding-guardrails.md
```

Priority order:

```text
correctness
security
latency
performance
observability
horizontal scalability
serverless readiness
maintainability
testability
```

## Architecture Rules

- Sovrunn is the sovereign cloud-native PaaS platform.
- SDE is a major differentiated capability inside Sovrunn, not the entire platform.
- Use canonical terminology from `docs/glossary.md`.
- Do not invent new resource kinds without approval.
- Do not rename accepted concepts.
- Do not change the feature sequence without approval.
- Do not implement future-phase capabilities early.

## Phase 1 Non-Goals

Do not implement in Phase 1 unless explicitly approved:

```text
persistent database storage
Kubernetes CRDs
GitOps controller
ServiceOps plugin execution
real infrastructure provisioning
UI portal
billing engine
AI agent execution
OPA/Gatekeeper integration
OpenTelemetry collector deployment
multi-cluster federation
SDE runtime transformation
```

## Coding Rules

- Use Go for Phase 1 platform core.
- Use in-memory registry first.
- Keep dependencies minimal.
- Prefer standard library where practical.
- Use metadata/spec/status resource shape.
- Treat `spec` as desired state.
- Treat `status` as system-owned observed state.
- Reject user-authored status where applicable.
- Add deterministic validation.
- Add tests for happy paths and failure paths.
- Use stable error codes from the API contract.
- Do not log secrets.
- Keep storage replaceable for future phases.
- Keep handlers short and context-aware.
- Use request IDs and structured observability fields.
- Do not add unapproved dependencies.

## Verification Commands

Run these before marking Go work complete:

```bash
make fmt
make test
make vet
```

When the API server exists:

```bash
make run
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/readyz
```

For concurrency-sensitive changes:

```bash
go test -race ./...
```

For documentation changes:

```bash
mkdocs build --strict
```

## Output Expected From AI Tools

When making changes, report:

```text
feature or task implemented
files changed
why each file changed
tests added
validation added
security considerations
observability considerations
performance considerations
commands run
command results
acceptance criteria satisfied
non-goals intentionally not implemented
known limitations
next feature boundary
```

## Final Rule

AI accelerates implementation.

Architecture remains spec-first, founder-controlled, test-gated, terminal-verified, and Git-reviewed.

## Architecture Operating System

Sovrunn uses an Architecture Operating System to keep long-term architecture stable across ChatGPT, Kiro, Cursor, reviewers, and multi-developer teams.

Before architecture work, agents must load:

- `docs/context/ARCHITECTURE_VERSION.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/SOVRUNN_CONTEXT_PACK.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/context/CHATGPT_ARCHITECTURE_SESSION_PROMPT.md`

Before implementation work, agents must also load:

- `docs/governance/REVIEW_GATES.md`
- `docs/engineering/go-observability-standard.md`
- `docs/architecture/observability-and-audit-baseline.md`

### Source-of-Truth Priority

1. Current architecture baseline
2. Accepted DEC files and Decision Index
3. Approved RFC files
4. Architecture docs
5. Phase scope docs
6. Feature specs
7. Roadmap placeholders
8. Chat discussion

Roadmap placeholders are directional only. They do not override current baseline, accepted decisions, or phase scope.

### Architecture Change Control

Agents must not change approved Sovrunn architecture casually.

Any proposed architecture change must be classified as:

- clarification
- extension
- correction
- replacement
- new decision

A replacement or new decision requires:

- architecture change request,
- impacted docs listed,
- impacted features listed,
- backward compatibility impact,
- phase impact,
- explicit human approval,
- updated DEC/RFC records,
- updated current architecture baseline.

No implementation may proceed from an unapproved architecture change.

### Feature Gate Rule

Before starting the next feature, the current feature must pass:

```bash
make ff-feature-gate FEATURE=<FEATURE-ID>
```

A feature is not complete unless the feature gate passes.

No agent may proceed to the next feature when:

- tests fail,
- lint fails,
- security checks fail,
- reuse assessment is missing,
- architecture drift checks are missing,
- acceptance criteria are missing,
- generated artifacts are staged,
- Phase 2 scope boundaries are violated,
- approval review is missing or not approved in strict team mode.

## ChatGPT-to-Kiro Architecture Handoff Rule

Architecture discussion may happen in a ChatGPT Project, but repository updates must be applied through Kiro using an approved Architecture Decision Handoff.

Required handoff contract:

- `docs/templates/ARCHITECTURE_DECISION_HANDOFF.md`
- `docs/prompts/chatgpt/architecture-decision-handoff.prompt.md`
- `docs/prompts/kiro/architecture-update.prompt.md`

Operating flow:

1. ChatGPT discusses tradeoffs and produces an Architecture Decision Handoff.
2. Human approves, rejects, or defers the handoff.
3. Kiro validates the approved handoff against the Architecture Operating System.
4. Kiro updates architecture docs, DEC/RFC files, traceability, and Kiro specs.
5. Cursor implements only from approved Kiro `tasks.md`.
6. Feature gate validates before moving to the next feature.

Rules:

- ChatGPT handoff is not approval by itself.
- Kiro must not apply an unapproved handoff.
- Kiro must not introduce new architecture beyond the handoff.
- Cursor must not change architecture while implementing Go code.
- `CURRENT_ARCHITECTURE_BASELINE.md` changes require accepted DEC/RFC or explicit baseline review approval.

## Structurizr Architecture Diagrams

Structurizr DSL is the approved architecture-as-code representation for Sovrunn C4 views.

When changing approved architecture structure, agents must update:

```text
docs/diagrams/structurizr/workspace.dsl
```

Update the Structurizr workspace when changes affect:

```text
system boundaries
major containers
plugin boundaries
external OSS/reuse relationships
deployment/runtime relationships
major dynamic flows
ChatGPT -> Kiro -> Cursor handoff workflow
```

Do not use ad hoc diagrams as the durable source of truth.
Markdown architecture docs, DEC/RFC records, and `workspace.dsl` must remain aligned.

Useful commands:

```bash
make structurizr-check
make structurizr-lite
```

Kiro owns architecture/spec updates, including Structurizr DSL updates when an approved Architecture Decision Handoff changes the architecture model.
Cursor must not change Structurizr DSL unless the approved task explicitly requires diagram/model maintenance.


Note: `docs/features/FEATURE_ROADMAP_ALL_PHASES.md` is a pointer to the canonical roadmap only. Do not treat it as a separate source of truth.
