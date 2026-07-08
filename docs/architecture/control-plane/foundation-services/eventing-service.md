# Eventing Service

Document:
  ID: eventing-service
  Title: Eventing Service
  Parent: control-plane-foundation
  Owner: Control Plane Foundation
  Layer: SDE Control Plane
  Type: Foundation Service Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Eventing Service
  - Define responsibilities and boundaries
  - Support selective AI retrieval

Definition:
  Eventing Service publishes and consumes management events across SDE Control Plane components.

Responsibilities:
  MUST:
    - Provide event channels
    - Preserve event metadata
    - Support reliable publication where configured
    - Apply event safety rules

  MUST NOT:
    - Become authoritative registry state
    - Replace Workflow Service
    - Expose unsafe payloads

Inputs:
  - Event message

Outputs:
  - Publication acknowledgement
  - Event delivery

State:
  - Topic or stream metadata

Failure Rules:
  - Return deterministic error
  - Preserve traceability
  - Avoid leaking unsafe internal details

Relationships:
  Parent:
    - control-plane-foundation
  Depends On:
    - foundation-providers/eventing-provider.md
  Used By:
    - Management Planes

References:
  - foundation-services.md
  - foundation-providers/eventing-provider.md
