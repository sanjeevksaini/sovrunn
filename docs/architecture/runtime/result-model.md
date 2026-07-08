# Result Model

Document:
  ID: result-model
  Title: Result Model
  Parent: runtime
  Owner: SDE Runtime
  Layer: SDE Data Plane
  Type: Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Result Model contract
  - Define fields, invariants, and usage rules
  - Support selective AI retrieval

Definition:
  Result Model is the canonical SDE runtime representation of successful or partial execution results.

Responsibilities:
  MUST:
    - Represent successful result
    - Represent partial result where applicable
    - Preserve type information
    - Preserve result metadata
    - Support protocol mapping

  MUST NOT:
    - Represent failure as success
    - Expose unsafe downstream-native result directly
    - Hide partial result state

Fields:
  - Result Identifier
  - Result Kind
  - Rows or Values
  - Affected Count
  - Schema Metadata
  - Cursor or Stream Reference
  - Partial State
  - Trace Identifier

Rules:
  - Engine Plugin maps native result to Result Model.
  - Protocol Plugin maps Result Model to protocol response.
  - Result Model must remain separate from Error Model.

Failure Rules:
  - Invalid result mapping must produce Error Model entry.
  - Partial result must be explicit.

Invariants:
  - Contract meaning MUST remain stable across runtime components.
  - Consumers MUST treat the contract as canonical for its scope.
  - Extensions MUST NOT weaken mandatory fields or safety rules.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - engine-runtime.md
    - data-kernel.md
  Used By:
    - result-flow.md
    - protocol-runtime.md

References:
  - runtime.md
  - runtime-map.md
  - engine-runtime.md
  - data-kernel.md
