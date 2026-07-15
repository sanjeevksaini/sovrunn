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
Phase 1: Sovrunn Platform Core
```

## Current Implementation Rule

Implement one feature at a time.

Feature order:

```text
FEATURE-0001 Organization Resource and Registry
FEATURE-0002 OrganizationUnit Resource
FEATURE-0003 Tenant Resource
FEATURE-0004 Project Resource
FEATURE-0005 Operation Resource
FEATURE-0006 ServiceClass and ServicePlan
FEATURE-0007 Plugin and Capability Registry
FEATURE-0008 ServiceInstance and ServiceBinding
FEATURE-0009 API server health/readiness
FEATURE-0010 Basic CLI/API demo flow
```

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
