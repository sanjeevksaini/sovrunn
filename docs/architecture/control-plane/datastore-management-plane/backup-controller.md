# Backup Controller

Document:
  ID: backup-controller
  Title: Backup Controller
  Parent: datastore-management-plane
  Owner: Datastore Management Plane
  Layer: SDE Control Plane
  Type: Control Plane Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Backup Controller
  - Define responsibilities and boundaries
  - Support selective AI retrieval

Definition:
  Coordinates datastore backup workflows.

Responsibilities:
  MUST:
    - Validate backup policy
    - Coordinate backup operation
    - Record backup metadata
    - Operate under SDE Control Plane authority
    - Use Foundation Services for shared concerns
    - Preserve lifecycle state consistency

  MUST NOT:
    - Execute client data requests
    - Process SIR
    - Replace SDE Data Plane
    - Replace Datastore Data Plane

Inputs:
  - Authorized lifecycle request
  - Tenant context
  - Policy context

Outputs:
  - Lifecycle operation result
  - Lifecycle state update

State:
  - Lifecycle metadata
  - Workflow metadata where applicable

Failure Rules:
  - Preserve Datastore Registry consistency
  - Record known lifecycle state
  - Avoid publishing execution-ready metadata until validated

Relationships:
  Parent:
    - datastore-management-plane
  Depends On:
    - ../foundation-services/workflow-service.md
    - ../foundation-services/audit-service.md
    - ../foundation-services/secrets-service.md
  Used By:
    - Datastore Management Plane
    - dstoreOps

References:
  - datastore-management-plane.md
