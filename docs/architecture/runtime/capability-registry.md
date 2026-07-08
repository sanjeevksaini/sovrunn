# Capability Registry

Document:
  ID: capability-registry
  Title: Capability Registry
  Parent: runtime
  Owner: Planning
  Layer: SDE Data Plane
  Type: Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Capability Registry contract
  - Define fields, invariants, and usage rules
  - Support selective AI retrieval

Definition:
  Capability Registry is the runtime-facing lookup of approved capability metadata.

Responsibilities:
  MUST:
    - Expose approved capability metadata
    - Support Planning lookup
    - Preserve capability identifiers
    - Reject unapproved capabilities

  MUST NOT:
    - Approve capabilities
    - Invent capabilities
    - Replace Capability Governance
    - Store raw plugin secrets

Fields:
  - Capability Identifier
  - Capability Version
  - Engine Binding
  - Support Level
  - Constraints
  - Compatibility Metadata

Rules:
  - Populated only from approved SDE Control Plane state.
  - Consumed by Planning.
  - Must preserve version and compatibility.

Failure Rules:
  - Missing capability must produce deterministic planning error.
  - Inconsistent metadata must fail closed.

Invariants:
  - Contract meaning MUST remain stable across runtime components.
  - Consumers MUST treat the contract as canonical for its scope.
  - Extensions MUST NOT weaken mandatory fields or safety rules.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - ../control-plane/core-control-plane/capability-governance.md
  Used By:
    - planning.md

References:
  - runtime.md
  - runtime-map.md
  - ../control-plane/core-control-plane/capability-governance.md
