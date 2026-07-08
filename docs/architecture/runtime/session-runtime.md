# Session Runtime

Document:
  ID: session-runtime
  Title: Session Runtime
  Parent: runtime
  Owner: Session Runtime
  Layer: SDE Data Plane
  Type: Component Architecture
  Version: 1.0
  Status: Stable

Purpose:
  - Define Session Runtime
  - Define runtime position
  - Define execution responsibilities
  - Define boundaries and failure rules

Definition:
  Session Runtime manages SDE session context and session lifecycle references.

Runtime Position:
  - Used by Protocol Runtime.
  - Referenced by Execution Context.
  - Preserves protocol session semantics where required.

Responsibilities:
  MUST:
    - Create session context
    - Resolve session context
    - Update session-scoped runtime parameters
    - Preserve tenant isolation
    - Expire sessions safely

  MUST NOT:
    - Own execution context
    - Own transaction lifecycle
    - Store secrets
    - Leak session state across tenants

Inputs:
  - Session open request
  - Session identifier
  - Runtime parameter update

Outputs:
  - Session reference
  - Session context
  - Session error

State:
  - Session metadata
  - Runtime parameter state
  - Prepared statement references where applicable

Execution Rules:
  - Externalize session context where required for stateless runtime.
  - Preserve session ownership.
  - Validate tenant boundary on lookup.

Failure Rules:
  - Fail closed on missing or unauthorized session.
  - Preserve cleanup behavior.

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
    - execution-context.md
    - error-model.md
  Used By:
    - protocol-runtime.md
    - session-flow.md

References:
  - runtime.md
  - runtime-map.md
  - execution-context.md
  - error-model.md
