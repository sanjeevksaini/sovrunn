# Error Model

Document:
  ID: error-model
  Title: Error Model
  Parent: runtime
  Owner: SDE Runtime
  Layer: SDE Data Plane
  Type: Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Error Model contract
  - Define fields, invariants, and usage rules
  - Support selective AI retrieval

Definition:
  Error Model is the canonical SDE runtime representation of execution and runtime failures.

Responsibilities:
  MUST:
    - Represent failure deterministically
    - Preserve error code, category, severity, source, retry classification, trace identifier, and timestamp
    - Preserve safe details
    - Support protocol error mapping

  MUST NOT:
    - Expose unsafe native errors directly
    - Convert failure into success
    - Lose causal chain
    - Omit timestamp

Fields:
  - Error Identifier
  - Code
  - Category
  - Message
  - Severity
  - Source
  - State
  - Retry Classification
  - Trace Identifier
  - Timestamp
  - Safe Details
  - Cause

Rules:
  - Timestamp generated at detection time and preserved.
  - Engine Plugin maps native errors to Error Model.
  - Protocol Plugin maps Error Model to protocol error response.

Failure Rules:
  - Unknown failure must still produce Error Model entry.
  - Unsafe details must be redacted.

Invariants:
  - Contract meaning MUST remain stable across runtime components.
  - Consumers MUST treat the contract as canonical for its scope.
  - Extensions MUST NOT weaken mandatory fields or safety rules.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - engine-runtime.md
    - protocol-runtime.md
  Used By:
    - error-flow.md
    - data-kernel.md

References:
  - runtime.md
  - runtime-map.md
  - engine-runtime.md
  - protocol-runtime.md
