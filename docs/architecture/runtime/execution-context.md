# Execution Context

Document:
  ID: execution-context
  Title: Execution Context
  Parent: runtime
  Owner: SDE Runtime
  Layer: SDE Data Plane
  Type: Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Execution Context contract
  - Define fields, invariants, and usage rules
  - Support selective AI retrieval

Definition:
  Execution Context is the immutable execution-scoped runtime context that accompanies Execution Plan execution.

Responsibilities:
  MUST:
    - Carry execution identity
    - Carry tenant and security context
    - Carry session and transaction references
    - Carry trace context
    - Remain immutable per execution

  MUST NOT:
    - Modify Execution Plan
    - Modify SIR
    - Own session lifecycle
    - Own transaction lifecycle
    - Leak across executions

Fields:
  - Execution Identifier
  - Request Identifier
  - Trace Identifier
  - Session Reference
  - Transaction Reference
  - Security Context
  - Tenant Context
  - Runtime Context
  - Deadline
  - Execution Options

Rules:
  - Created per execution.
  - Passed to runtime components.
  - Must remain immutable.
  - Must not contain raw secrets.

Failure Rules:
  - Missing required context must fail execution.
  - Context mismatch must fail closed.

Invariants:
  - Contract meaning MUST remain stable across runtime components.
  - Consumers MUST treat the contract as canonical for its scope.
  - Extensions MUST NOT weaken mandatory fields or safety rules.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - session-runtime.md
    - transaction-runtime.md
  Used By:
    - data-kernel.md
    - engine-runtime.md
    - execution-flow.md

References:
  - runtime.md
  - runtime-map.md
  - session-runtime.md
  - transaction-runtime.md
