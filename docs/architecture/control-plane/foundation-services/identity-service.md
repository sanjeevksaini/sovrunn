# Identity Service

Document:
  ID: identity-service
  Title: Identity Service
  Parent: control-plane-foundation
  Owner: Control Plane Foundation
  Layer: SDE Control Plane
  Type: Foundation Service Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Identity Service
  - Define responsibilities and boundaries
  - Support selective AI retrieval

Definition:
  Identity Service authenticates and resolves management actors, service identities, runtime component identities, Management Plane identities, plugin identities, and Foundation Provider identities.

Responsibilities:
  MUST:
    - Authenticate identities
    - Resolve identity context
    - Support identity federation where configured
    - Expose stable identity claims to authorized consumers

  MUST NOT:
    - Authorize actions
    - Store authorization policy
    - Store secrets
    - Execute runtime operations

Inputs:
  - Authentication request
  - Credential reference
  - Tenant context

Outputs:
  - Identity context
  - Authenticated principal
  - Identity error

State:
  - Identity provider binding
  - Identity session metadata where applicable

Failure Rules:
  - Fail closed when identity cannot be established
  - Preserve safe authentication error details

Relationships:
  Parent:
    - control-plane-foundation
  Depends On:
    - foundation-providers/identity-provider.md
  Used By:
    - Authorization Service
    - Management Planes

References:
  - foundation-services.md
  - foundation-providers/identity-provider.md
