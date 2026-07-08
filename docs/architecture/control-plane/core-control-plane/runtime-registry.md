# Runtime Registry

Document:
  ID: runtime-registry
  Title: Runtime Registry
  Parent: core-control-plane
  Owner: Core Control Plane
  Layer: SDE Control Plane
  Type: Control Plane Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Runtime Registry
  - Define responsibilities and boundaries
  - Support selective AI retrieval

Definition:
  Tracks SDE runtime instances, runtime topology, health, version, and availability metadata.

Responsibilities:
  MUST:
    - Register runtime instance
    - Track runtime health
    - Expose approved runtime topology
    - Support runtime compatibility checks
    - Use Foundation Services for shared control-plane concerns

  MUST NOT:
    - Execute SDE Data Plane requests
    - Own downstream datastore lifecycle
    - Bypass SDE Control Plane policy

Inputs:
  - Authorized registry or governance request
  - Policy context
  - Tenant context

Outputs:
  - Approved metadata
  - Registry state
  - Governance decision

State:
  - Authoritative metadata
  - Version state
  - Lifecycle state

Failure Rules:
  - Preserve authoritative state consistency
  - Reject invalid metadata deterministically

Relationships:
  Parent:
    - core-control-plane
  Depends On:
    - ../foundation-services/registry-framework.md
    - ../foundation-services/plugin-framework.md
  Used By:
    - SDE Data Plane
    - Planning
    - Engine Runtime
    - Management Planes

References:
  - core-control-plane.md
