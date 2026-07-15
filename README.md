# Sovrunn

Sovrunn is a sovereign cloud-native PaaS platform for organizations that need governed, multi-tenant, open-standard cloud platform capabilities across their own infrastructure, local cloud providers, colocation providers, and regulated datacenter environments.

Sovrunn builds on proven open-source and open-standard technologies instead of replacing them. It provides the missing product layer for organization governance, tenant isolation, service catalog, service plans, lifecycle operations, policy inheritance, auditability, observability, backup governance, plugin-based service management, and AI-assisted operations.

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
├── .kiro/
│   └── steering/
├── .cursor/
│   └── rules/
├── docs/
│   ├── foundation/
│   ├── architecture/
│   ├── decisions/
│   ├── engineering/
│   ├── features/
│   ├── resource-specs/
│   ├── api/
│   ├── prompts/
│   ├── rfc/
│   └── demo/
├── configs/
├── scripts/
└── cmd/
```

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
