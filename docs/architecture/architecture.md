# Architecture

Document
- ID: architecture
- Version: 1.1
- Status: Stable

Purpose
- Define Sovrunn Data Engine architecture map
- Define top-level architectural domains
- Define domain boundaries
- Define architecture document organization

Definition

Sovrunn Data Engine (SDE) architecture defines how semantic intent is accepted, represented, planned, executed, governed, and integrated with heterogeneous downstream datastores.

SDE is a sovereign semantic execution platform.

SDE is not a data engine platform.

Top-Level Architecture

```text
Sovrunn Data Engine

├── SDE Data Plane
│   ├── Protocol Runtime
│   ├── SIR Runtime
│   ├── Planning
│   ├── Execution Plan
│   ├── Execution Context
│   ├── Data Kernel
│   ├── Engine Runtime
│   └── Engine Plugins
│
├── SDE Control Plane
│   ├── Control Plane Foundation
│   ├── Core Control Plane
│   └── Datastore Management Plane
│
└── Downstream Datastores
    └── Datastore SDE Data Plane
```

Principles

MUST

- Preserve semantic intent
- Separate SDE Data Plane from SDE Control Plane
- Separate SDE Data Plane from Datastore SDE Data Plane
- Separate runtime execution from management governance
- Separate planning from execution
- Separate SDE from Downstream Datastores
- Preserve protocol independence
- Preserve datastore independence
- Prefer capability-driven execution
- Preserve policy and tenant boundaries

MUST NOT

- Treat SDE as a data engine platform
- Couple SDE semantics to a protocol
- Couple SDE semantics to a Downstream Datastore
- Merge SDE Data Plane with Datastore SDE Data Plane
- Let SDE Control Plane execute client data requests
- Let SDE Data Plane manage downstream datastore lifecycle

Architectural Domains

SDE Data Plane

- Executes client data requests
- Transforms protocol intent into SIR
- Plans SIR into Execution Plan
- Coordinates execution through Data Kernel
- Delegates downstream execution through Engine Runtime and Engine Plugins
- Produces Result Model and Error Model

SDE Control Plane

- Provides management authority
- Provides Control Plane Foundation
- Hosts Management Planes
- Governs runtime metadata
- Governs plugins, engines, capabilities, deployment, policy, tenancy, and audit
- Does not execute client data requests

Datastore Management Plane

- Optional Management Plane
- Powers dstoreOps
- Manages downstream datastore lifecycle
- Uses Datastore Operator Plugins and Infrastructure Providers
- Does not execute client data requests

Downstream Datastore

- Owns native storage
- Owns native execution
- Owns native optimization
- Owns Datastore SDE Data Plane
- Is accessed by SDE through Engine Plugins for execution

Boundary Summary

Semantic Boundary

- SIR
- SIR Specification

Execution Boundary

- SDE Data Plane
- Execution Plan
- Execution Context
- Data Kernel
- Engine Runtime

Management Boundary

- SDE Control Plane
- Control Plane Foundation
- Management Planes

Datastore Boundary

- Downstream Datastore
- Datastore SDE Data Plane
- Native datastore APIs

Extension Boundary

- Protocol Plugin
- Engine Plugin
- Foundation Provider
- Management Plane Plugin
- Datastore Operator Plugin
- Infrastructure Provider

Document Map

Runtime Architecture

- runtime/runtime.md
- runtime/protocol-runtime.md
- runtime/sir-runtime.md
- runtime/planning.md
- runtime/execution-plan.md
- runtime/execution-context.md
- runtime/data-kernel.md
- runtime/engine-runtime.md
- runtime/capability-registry.md
- runtime/session-runtime.md
- runtime/transaction-runtime.md
- runtime/result-model.md
- runtime/error-model.md

SDE Control Plane Architecture

- control-plane/control-plane.md
- control-plane/control-plane-foundation.md
- control-plane/foundation-services.md
- control-plane/foundation-providers.md
- control-plane/management-plane.md
- control-plane/core-control-plane.md
- control-plane/datastore-management-plane.md

SDE Data Plane Architecture

- data-plane/data-plane.md
- data-plane/request-flow.md
- data-plane/protocol-execution.md
- data-plane/planning-execution.md
- data-plane/kernel-execution.md
- data-plane/engine-execution.md
- data-plane/result-propagation.md
- data-plane/error-propagation.md

Ownership

Sovrunn owns

- Sovrunn Data Engine architecture
- SDE Data Plane architecture
- SDE Control Plane architecture
- Runtime architecture
- Semantic execution architecture

Downstream Datastore owns

- Native datastore state
- Native datastore storage
- Native datastore execution
- Datastore SDE Data Plane

References

- ../foundation/constitution.md
- ../foundation/glossary.md
- ../foundation/ontology.md
- ../specifications/sir/sir.md
- ../specifications/sir/capability-model.md
- ../specifications/reuse/adopted-architecture-patterns.md
- runtime/runtime.md
- control-plane/control-plane.md
- data-plane/data-plane.md


Control Plane Refactor

- docs/architecture/control-plane/control-plane.md
- docs/architecture/control-plane/foundation-services/foundation-services.md
- docs/architecture/control-plane/foundation-providers/foundation-providers.md
- docs/architecture/control-plane/core-control-plane/core-control-plane.md
- docs/architecture/control-plane/datastore-management-plane/datastore-management-plane.md


Control Plane Architecture v2

Control Plane Architecture v2:

  Map:
    - docs/architecture/control-plane/control-plane-map.md

  Architecture:
    - docs/architecture/control-plane/control-plane.md
    - docs/architecture/control-plane/control-plane-foundation.md
    - docs/architecture/control-plane/management-plane.md
    - docs/architecture/control-plane/foundation-services/foundation-services.md
    - docs/architecture/control-plane/foundation-providers/foundation-providers.md
    - docs/architecture/control-plane/core-control-plane/core-control-plane.md
    - docs/architecture/control-plane/datastore-management-plane/datastore-management-plane.md

  Contracts:
    - docs/architecture/control-plane/foundation-services/*.md
    - docs/architecture/control-plane/foundation-providers/*.md
    - docs/architecture/control-plane/core-control-plane/*.md
    - docs/architecture/control-plane/datastore-management-plane/*.md


Runtime Architecture Refactor

Runtime Architecture Refactor:

  Map:
    - docs/architecture/runtime/runtime-map.md

  Architecture:
    - docs/architecture/runtime/runtime.md

  Components:
    - docs/architecture/runtime/protocol-runtime.md
    - docs/architecture/runtime/sir-runtime.md
    - docs/architecture/runtime/planning.md
    - docs/architecture/runtime/data-kernel.md
    - docs/architecture/runtime/engine-runtime.md
    - docs/architecture/runtime/plugin-runtime.md
    - docs/architecture/runtime/session-runtime.md
    - docs/architecture/runtime/transaction-runtime.md

  Contracts:
    - docs/architecture/runtime/execution-plan.md
    - docs/architecture/runtime/execution-context.md
    - docs/architecture/runtime/capability-registry.md
    - docs/architecture/runtime/result-model.md
    - docs/architecture/runtime/error-model.md

  Flows:
    - docs/architecture/runtime/execution-flow.md
    - docs/architecture/runtime/session-flow.md
    - docs/architecture/runtime/transaction-flow.md
    - docs/architecture/runtime/result-flow.md
    - docs/architecture/runtime/error-flow.md
