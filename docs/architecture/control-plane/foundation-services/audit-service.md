# Audit Service

Document:
  ID: audit-service
  Title: Audit Service
  Parent: control-plane-foundation
  Owner: Control Plane Foundation
  Layer: SDE Control Plane
  Type: Foundation Service Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Audit Service
  - Define responsibilities and boundaries
  - Support selective AI retrieval

Definition:
  Audit Service records governance-relevant management events and security-sensitive actions.

Responsibilities:
  MUST:
    - Record audit events
    - Preserve actor, action, target, timestamp, and outcome
    - Protect audit integrity
    - Support audit query where authorized

  MUST NOT:
    - Replace observability telemetry
    - Execute business logic
    - Store secrets

Inputs:
  - Audit event

Outputs:
  - Audit acknowledgement
  - Audit query result where authorized

State:
  - Audit log metadata

Failure Rules:
  - Return deterministic error
  - Preserve traceability
  - Avoid leaking unsafe internal details

Relationships:
  Parent:
    - control-plane-foundation
  Depends On:
    - foundation-providers/audit-provider.md
  Used By:
    - All Management Planes

References:
  - foundation-services.md
  - foundation-providers/audit-provider.md
