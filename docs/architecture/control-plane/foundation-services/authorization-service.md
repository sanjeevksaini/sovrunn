# Authorization Service

Document:
  ID: authorization-service
  Title: Authorization Service
  Parent: control-plane-foundation
  Owner: Control Plane Foundation
  Layer: SDE Control Plane
  Type: Foundation Service Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Authorization Service
  - Define responsibilities and boundaries
  - Support selective AI retrieval

Definition:
  Authorization Service authorizes management and runtime-adjacent actions using identity context, tenant context, policy context, and requested operation.

Responsibilities:
  MUST:
    - Evaluate authorization decisions
    - Consume Policy Service where required
    - Return allow, deny, or indeterminate
    - Preserve decision auditability

  MUST NOT:
    - Authenticate identities
    - Own identity lifecycle
    - Store secrets
    - Execute domain operations

Inputs:
  - Identity context
  - Requested action
  - Resource context
  - Policy context

Outputs:
  - Authorization decision
  - Decision reason
  - Authorization error

State:
  - Policy bindings
  - Decision metadata

Failure Rules:
  - Return deterministic error
  - Preserve traceability
  - Avoid leaking unsafe internal details

Relationships:
  Parent:
    - control-plane-foundation
  Depends On:
    - foundation-providers/authorization-provider.md
    - policy-service.md
  Used By:
    - Management Planes
    - Foundation Services

References:
  - foundation-services.md
  - foundation-providers/authorization-provider.md
  - policy-service.md
