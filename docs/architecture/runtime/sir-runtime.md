# SIR Runtime

Document:
  ID: sir-runtime
  Title: SIR Runtime
  Parent: runtime
  Owner: SIR Runtime
  Layer: SDE Data Plane
  Type: Component Architecture
  Version: 1.0
  Status: Stable

Purpose:
  - Define SIR Runtime
  - Define runtime position
  - Define execution responsibilities
  - Define boundaries and failure rules

Definition:
  SIR Runtime creates, validates, and manages live SIR instances before Planning.

Runtime Position:
  - Receives semantic intent from Protocol Runtime.
  - Produces valid SIR for Planning.
  - Preserves SIR semantics.

Responsibilities:
  MUST:
    - Create SIR instance
    - Validate SIR structure
    - Validate SIR version
    - Validate semantic references
    - Forward valid SIR to Planning

  MUST NOT:
    - Plan execution
    - Select Downstream Datastore
    - Invoke Engine Plugin
    - Modify semantic intent for optimization

Inputs:
  - Protocol-normalized intent
  - SIR schema version
  - Tenant context

Outputs:
  - Validated SIR instance
  - SIR validation error

State:
  - Transient SIR instance state
  - Validation metadata

Execution Rules:
  - Reject invalid SIR deterministically.
  - Preserve semantic intent.
  - Use canonical SIR versioning rules.

Failure Rules:
  - Return Error Model entry for validation failure.
  - Preserve source location where available.

Concurrency Rules:
  - Preserve execution isolation.
  - Avoid shared mutable execution state unless explicitly synchronized.
  - Preserve tenant isolation across concurrent executions.
  - Preserve session and transaction boundaries.

Security Rules:
  - Enforce authorized execution context.
  - Avoid exposing secrets.
  - Preserve safe error behavior.
  - Preserve trace and audit correlation where applicable.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - planning.md
    - error-model.md
  Used By:
    - execution-flow.md
    - planning.md

References:
  - runtime.md
  - runtime-map.md
  - planning.md
  - error-model.md
