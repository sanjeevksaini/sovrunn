# Architecture

Document
- ID: architecture
- Version: 1.0
- Status: Stable

Purpose
- Define Sovrunn runtime architecture
- Define architectural boundaries
- Define runtime responsibilities
- Define component interactions

Definition

The Sovrunn Architecture defines the runtime structure that transforms semantic intent into engine execution while preserving semantic equivalence.

Principles

MUST

- Preserve semantic intent
- Separate semantics from execution
- Separate planning from execution
- Separate runtime from engines
- Separate control plane from data plane
- Preserve engine independence
- Preserve protocol independence
- Prefer capability driven execution

MUST NOT

- Couple architecture to protocol
- Couple architecture to engine
- Couple architecture to deployment
- Couple architecture to implementation

Architecture

```text
                Application
                      │
                      ▼
              Protocol Runtime
                      │
                      ▼
                    SIR
                      │
                      ▼
                  Planning
                      │
                      ▼
               Execution Plan
                      │
                      ▼
                Data Kernel
                      │
            ┌─────────┴─────────┐
            ▼                   ▼
     Capability Registry   Engine Runtime
                                    │
                                    ▼
                              Engine Plugin
                                    │
                                    ▼
                                 Engine
```

Components

Protocol Runtime
- Produces SIR
- Consumes Protocols

Planning
- Produces Execution Plan
- Consumes SIR
- Uses Capability Registry

Execution Plan
- Defines executable runtime model

Data Kernel
- Coordinates execution
- Preserves semantic equivalence

Capability Registry
- Publishes engine capabilities
- Supports capability discovery

Engine Runtime
- Translates execution plan
- Delegates to engine plugins

Engine Plugin
- Implements engine integration
- Publishes capability manifest

Engine
- Executes native operations

Boundaries

Semantic Boundary

- Protocol Runtime
- SIR

Planning Boundary

- Planning
- Execution Plan

Execution Boundary

- Data Kernel
- Engine Runtime

Storage Boundary

- Engine

Ownership

Sovrunn owns

- Runtime architecture
- Planning
- Execution Plan
- Data Kernel
- Capability Registry
- Protocol Runtime
- Engine Runtime

Engine owns

- Storage
- Physical execution
- Native optimization

Plugin owns

- Engine integration
- Capability Manifest

References

- foundation/constitution.md
- specifications/sir/sir.md
- specifications/sir/capability-model.md
- specifications/reuse/adopted-architecture-patterns.md
