# Datastore Operator Plugin

Document:
  ID: datastore-operator-plugin
  Title: Datastore Operator Plugin
  Parent: datastore-management-plane
  Owner: Datastore Management Plane
  Layer: SDE Control Plane
  Type: Control Plane Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Datastore Operator Plugin
  - Define responsibilities and boundaries
  - Support selective AI retrieval

Definition:
  Manages datastore-specific lifecycle operations for a downstream datastore technology.

Responsibilities:
  MUST:
    - Perform datastore-specific lifecycle operations
    - Declare lifecycle capabilities
    - Report lifecycle state
    - Use Secrets Service for credentials
    - Operate under SDE Control Plane authority
    - Use Foundation Services for shared concerns
    - Preserve lifecycle state consistency

  MUST NOT:
    - Execute client data requests
    - Process SIR
    - Replace SDE Data Plane
    - Replace Datastore Data Plane
    - Execute SDE Execution Plans
    - Return query results
    - Replace Engine Plugin

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
