# Sovrunn

Sovrunn is a sovereign cloud-native PaaS platform for organizations that need governed, multi-tenant, open-standard cloud platform capabilities across their own infrastructure, local cloud providers, colocation providers, and regulated datacenter environments.

Sovrunn builds on proven open-source and open-standard technologies instead of replacing them. It provides the missing product layer for organization governance, tenant isolation, service catalog, service plans, lifecycle operations, policy inheritance, auditability, observability, backup governance, plugin-based service management, and AI-assisted operations.


## Current Roadmap Context

Sovrunn now uses a reuse-first phased roadmap. Phase 2 and Phase 3 are the current execution focus. Later phase features are maintained as scope placeholders only and must be rebaselined after Phase 2 and Phase 3 complete.

Authoritative roadmap files:

```text
docs/architecture/development-phases.md
docs/phase2/PHASE2_SCOPE.md
docs/phase2/PHASE2_FEATURE_SEQUENCE.md
docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md
docs/features/FEATURE_INDEX.md
```

Note: `docs/features/FEATURE_ROADMAP_ALL_PHASES.md` is only a pointer to the canonical roadmap to avoid duplicated roadmap sources.

Current MVP anchor:

```text
MVP-001: Governed PostgreSQL PaaS Placement and Provisioning on one substrate.
```

Roadmap rule:

```text
Use future features for scope awareness.
Do not implement future-phase features during Phase 2 or Phase 3 unless a formal decision changes the phase boundary.
```

## Positioning

Sovrunn is the parent platform.

Sovrunn Data Engine, abbreviated as SDE, is a major differentiated capability inside Sovrunn. SDE is not the entire platform.

```text
Sovrunn
  -> sovereign cloud-native PaaS platform

SDE
  -> interoperable data engine capability inside Sovrunn
```

## Vision

Sovrunn enables large organizations, government platforms, regulated enterprises, local cloud providers, and colocation providers to offer governed cloud-native platform services from their own sovereign infrastructure.

A single enterprise-grade Sovrunn deployment should provide centralized cloud management across multiple sovereign datacenter locations, clusters, zones, and infrastructure accounts.

## Core Platform Model

```text
Users and Access Channels
  -> Portal, CLI, API, GitOps, SDKs, AI Assistant

Organization Management Layer
  -> Organization, OrganizationUnit, Tenant, Project
  -> sovereign location context, multi-account management,
     billing and cost management, centralized cloud operations,
     cross-account resource sharing

Cloud Management Plane
  -> API server, resource registry, service catalog,
     Service Management Plane registry, plugin registry,
     capability registry, operation framework

Service Management Planes
  -> datastore, cache, object storage, stream, gateway,
     load balancer, FaaS, big data, SDE

Execution Substrate
  -> Kubernetes, operators, GitOps, policy, observability,
     identity, secrets, storage, networking
```

## Phase 1 Scope

Phase 1 builds the platform grammar and implementation foundation.

```text
Organization
OrganizationUnit
Tenant
Project
Operation
ServiceClass
ServicePlan
Plugin
Capability
ServiceInstance
ServiceBinding
```

Phase 1 starts with a simple Go-based platform core, REST APIs, in-memory registry, deterministic validation, testable resource behavior, and repeatable demo flows.

## Phase 1 Non-Goals

Phase 1 does not implement:

```text
production UI
persistent database storage
Kubernetes CRDs
GitOps controller
ServiceOps plugin execution
real infrastructure provisioning
billing engine
marketplace
multi-cluster federation
AI agent execution
SDE runtime transformation
```

These are future capabilities after the Phase 1 resource model and platform core are stable.

## Key Concepts

### Organization

Top-level governance owner.

### OrganizationUnit

Delegated governance unit inside an Organization.

### Tenant

Primary isolated consumption boundary.

### Project

Workload or environment grouping inside a Tenant.

### ServiceClass

Type of service that Sovrunn can offer.

### ServicePlan

Approved configuration and policy shape for a service.

### ServiceInstance

Tenant/project-scoped requested service.

### ServiceBinding

Consumption relationship between an application or workload and a ServiceInstance.

### Operation

Lifecycle record of what happened, to which resource, under which context, and with what result.

### Plugin

Registered implementation provider for a Service Management Plane.

### Capability

Declared lifecycle action supported by a plugin.

## Sovrunn Data Engine

SDE is the interoperable data-engine capability inside Sovrunn.

SDE focuses on:

```text
protocol transparency
semantic data access
datastore portability
runtime transformation
SIR-based execution
data-plane optimization
```

SDE is important, but it is one capability within the broader Sovrunn platform.

## Repository Structure

```text
.
├── AGENTS.md
├── Makefile
├── mkdocs.yml
├── docs/
│   ├── context/          # current architecture baseline, context pack, phase context
│   ├── governance/       # change control, ownership, review gates
│   ├── architecture/     # approved architecture source-of-truth docs
│   ├── decisions/        # DEC index and individual decision records
│   ├── rfc/              # architecture RFCs
│   ├── phase2/           # current execution phase scope and gates
│   ├── roadmap/          # canonical all-phase feature roadmap
│   ├── traceability/     # feature/decision traceability matrices
│   ├── templates/        # ACR, handoff, DEC, RFC, review templates
│   ├── diagrams/         # Structurizr architecture-as-code workspace
│   ├── engineering/      # Go, observability, context, and AI development standards
│   ├── features/         # Phase 1 feature docs and feature index
│   ├── prompts/          # Kiro, Cursor, ChatGPT, and reviewer prompt templates
│   ├── reviews/          # architecture, feature, phase, and monthly reviews
│   └── demo/             # demo flows
├── configs/
├── scripts/
└── tests/
```

Generated artifacts such as `site/`, `docs/generated-prompts/`, generated context packs, logs, and zip archives are intentionally ignored and must not be treated as architecture source of truth.

## Documentation Entry Points

Start here:

```text
docs/foundation/vision.md
docs/foundation/philosophy.md
docs/foundation/constitution.md
docs/decisions/DECISION_INDEX.md
docs/glossary.md
docs/features/FEATURE_SEQUENCE.md
docs/resource-specs/RESOURCE_MODEL_PHASE1.md
docs/api/API_CONTRACT_PHASE1.md
```

For AI-assisted development:

```text
AGENTS.md
docs/engineering/ai-controlled-development.md
docs/engineering/context-engineering-standard.md
docs/prompts/AI_IMPLEMENTATION_PROMPT_PHASE1.md
```

For tool-specific guidance:

```text
.kiro/steering/
.cursor/rules/
docs/prompts/CHATGPT_REVIEW_PROMPT.md
docs/prompts/CHATGPT_PATCH_PROMPT.md
docs/prompts/CHATGPT_DEBUG_PROMPT.md
```

## Development Workflow

This repository is designed to be used across multiple tools while keeping Git and terminal verification as the source of truth.

```text
Same Git repo

├── Kiro       = architecture/spec/task planning
├── Cursor     = code editing/refactoring/debugging
├── ChatGPT    = deep architecture review + exact patches
└── Terminal   = source of truth: build/test/runtime
```

## AI Development Rules

AI tools may generate code, tests, docs, patches, and reviews.

AI tools must not:

```text
invent architecture
rename canonical terms
change feature sequence
implement future features early
skip tests
introduce unapproved dependencies
bypass validation
remove audit or operation hooks
```

Before implementing any feature, load:

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

## Local Validation

Clean generated and local-only files:

```bash
find . -name ".DS_Store" -delete
```

Run formatting, tests, and vet:

```bash
make fmt
make test
make vet
```

Run local API server, after implementation exists:

```bash
make run
```

Verify health, after implementation exists:

```bash
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/readyz
```

Run demo, after implementation exists:

```bash
make demo
```

## Running the Demo

Prerequisites:
- **bash 4+** required. macOS ships bash 3.2 by default. Install a modern bash via
  Homebrew (`brew install bash`) and either invoke with `/usr/local/bin/bash` (Intel)
  or `/opt/homebrew/bin/bash` (Apple Silicon), or add the Homebrew bash to your PATH.
  Verify with `bash --version`.
- **curl** must be available (pre-installed on macOS and most Linux distributions).
- The `sovrunn-api` server must be running (`make run`).

Run the full Phase 1 demo flow:

```bash
make run &
sleep 2
make demo
```

The demo exercises all Phase 1 resources end-to-end. To re-run, restart the server first:

```bash
kill %1
make run &
sleep 2
make demo
```

Override the base URL if the server runs on a different port:

```bash
BASE_URL=http://127.0.0.1:9090 make demo
```

## Documentation Site

If MkDocs is configured, build the documentation with:

```bash
mkdocs build --strict
```

Serve locally:

```bash
mkdocs serve
```

Then open:

```text
http://127.0.0.1:8000
```

Generated site output should not be committed.

## Branching

Recommended branch model:

```text
main
  stable accepted baseline

phase1-foundation
  Phase 1 architecture, docs, standards, and implementation backlog

feature-0001-organization-registry
  implementation branch for FEATURE-0001

feature-0002-organizationunit-resource
  implementation branch for FEATURE-0002
```

Do not begin feature implementation until `phase1-foundation` is finalized.

## Current Status

The project is in Phase 1 foundation and platform-core implementation preparation.

Next implementation milestone:

```text
FEATURE-0001 Organization Resource and Registry
```

## License

License to be decided.

## Architecture Operating System

Sovrunn architecture is evolved through a durable Architecture Operating System.

The repo is the source of truth. ChatGPT, Kiro, Cursor, and reviewers must work from the approved baseline and decision records, not from chat memory.

Key files:

- `docs/context/ARCHITECTURE_VERSION.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/SOVRUNN_CONTEXT_PACK.md`
- `docs/context/CHATGPT_ARCHITECTURE_SESSION_PROMPT.md`
- `docs/governance/ARCHITECTURE_CHANGE_CONTROL.md`
- `docs/governance/ARCHITECTURE_OWNERSHIP.md`
- `docs/governance/REVIEW_GATES.md`
- `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md`
- `docs/traceability/DECISION_TRACEABILITY_MATRIX.md`

Generate a fresh context pack before major architecture sessions:

```bash
make context-pack
```

Validate a feature before moving to the next feature:

```bash
make ff-feature-gate FEATURE=FEATURE-0011
```

Architecture rule:

```text
ChatGPT can propose.
Git repo decides.
DEC/RFC records approve.
Feature gate enforces.
```

## ChatGPT + Kiro Architecture Operating Model

Sovrunn uses a governed handoff between architecture discussion and repo updates.

```text
ChatGPT Project
  -> architecture tradeoff discussion
  -> Architecture Decision Handoff

Kiro
  -> validates handoff against Architecture Operating System
  -> updates architecture docs / DEC / RFC / requirements / design / tasks

Cursor
  -> implements Go code from approved Kiro tasks only

Feature Gate
  -> validates before moving to the next feature
```

Key files:

- `docs/templates/ARCHITECTURE_DECISION_HANDOFF.md`
- `docs/prompts/chatgpt/architecture-decision-handoff.prompt.md`
- `docs/prompts/kiro/architecture-update.prompt.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/governance/ARCHITECTURE_CHANGE_CONTROL.md`

Before Kiro applies an architecture decision handoff, validate the handoff structure:

```bash
make arch-handoff-check HANDOFF=docs/reviews/architecture-decision-handoffs/ADH-YYYY-NNN.md
```

Then Kiro may use `docs/prompts/kiro/architecture-update.prompt.md` to apply the approved handoff.

## Structurizr Architecture Diagrams

Sovrunn uses Structurizr DSL as the architecture-as-code model for approved C4 views.

Key files:

```text
docs/diagrams/structurizr/workspace.dsl
docs/diagrams/structurizr/README.md
```

Run local validation/checks:

```bash
make structurizr-check
```

Run Structurizr Lite locally:

```bash
make structurizr-lite
```

Then open:

```text
http://localhost:8080
```

Architecture changes that alter system boundaries, containers, plugin planes, major external systems, or dynamic flows should update `workspace.dsl` as part of the Kiro architecture update workflow.


## Minimum Daily Workflow

```text
1. Use ChatGPT Project only for architecture tradeoff discussion when needed.
2. Convert approved architecture discussion into an Architecture Decision Handoff.
3. Use Kiro to validate/apply the handoff to architecture docs and feature specs.
4. Use Cursor to implement only approved Kiro tasks.
5. Run `make ff-feature-gate FEATURE=<FEATURE-ID>` before moving to the next feature.
```

Small clarifications may update docs/specs through a handoff. Major architecture changes require ACR/DEC/RFC updates and human approval.
