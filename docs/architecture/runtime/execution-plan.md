# Execution Plan

Document:
  ID: execution-plan
  Title: Execution Plan
  Parent: runtime
  Owner: Planning
  Layer: SDE Data Plane
  Type: Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Execution Plan contract
  - Define fields, invariants, and usage rules
  - Support selective AI retrieval

Definition:
  Execution Plan is the immutable runtime execution contract produced by Planning and consumed by Data Kernel.

Responsibilities:
  MUST:
    - Represent executable runtime work
    - Preserve SIR intent
    - Define operation dependencies
    - Declare required capabilities
    - Remain immutable during execution

  MUST NOT:
    - Contain protocol-specific commands as platform contract
    - Own execution lifecycle
    - Own mutable runtime state
    - Silently change semantic intent

Fields:
  - Plan Identifier
  - SIR Reference
  - Operation Graph
  - Capability Requirements
  - Engine Bindings
  - Execution Constraints
  - Result Shape

Rules:
  - Created only by Planning.
  - Consumed by Data Kernel.
  - Must be immutable after creation.
  - Must fail validation if required capability is unavailable.

Failure Rules:
  - Invalid plan must not execute.
  - Unsupported capability must produce deterministic planning error.

Invariants:
  - Contract meaning MUST remain stable across runtime components.
  - Consumers MUST treat the contract as canonical for its scope.
  - Extensions MUST NOT weaken mandatory fields or safety rules.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - planning.md
    - capability-registry.md
  Used By:
    - data-kernel.md
    - execution-flow.md

References:
  - runtime.md
  - runtime-map.md
  - planning.md
  - capability-registry.md
