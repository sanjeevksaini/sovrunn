# Workflow Service

Document:
  ID: workflow-service
  Title: Workflow Service
  Parent: control-plane-foundation
  Owner: Control Plane Foundation
  Layer: SDE Control Plane
  Type: Foundation Service Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Workflow Service
  - Define responsibilities and boundaries
  - Support selective AI retrieval

Definition:
  Workflow Service orchestrates long-running management operations.

Responsibilities:
  MUST:
    - Run durable workflows
    - Track workflow state
    - Coordinate retries
    - Coordinate compensation where defined

  MUST NOT:
    - Define domain semantics
    - Execute client data requests
    - Replace domain controllers

Inputs:
  - Workflow definition
  - Workflow input

Outputs:
  - Workflow identifier
  - Workflow state
  - Workflow result

State:
  - Workflow execution state

Failure Rules:
  - Return deterministic error
  - Preserve traceability
  - Avoid leaking unsafe internal details

Relationships:
  Parent:
    - control-plane-foundation
  Depends On:
    - foundation-providers/workflow-provider.md
  Used By:
    - Datastore Management Plane

References:
  - foundation-services.md
  - foundation-providers/workflow-provider.md
