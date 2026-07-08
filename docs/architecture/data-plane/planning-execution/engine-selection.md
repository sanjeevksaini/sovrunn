# Engine Selection

Document:
  ID: engine-selection
  Title: Engine Selection
  Parent: planning-execution
  Owner: Planning
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Engine Selection stage of Planning Execution
  - Keep planning-execution.md concise
  - Define focused planning-stage rules

Definition:
  Engine Selection defines how Planning selects compatible execution targets from approved engine metadata.

Flow:
  - Planning identifies candidate engines from approved Engine Registry metadata.
  - Planning validates Engine Plugin binding.
  - Planning filters candidates by required capabilities.
  - Planning filters candidates by policy constraints.
  - Planning validates client engine preference when present.
  - Planning selects compatible execution target or target set.

Rules:
  - Engine candidates MUST come from approved engine metadata.
  - Engine Plugin binding MUST be approved.
  - Client preference MUST be validated.
  - No compatible candidate MUST fail deterministically.
  - Planning MUST NOT access Downstream Datastore during selection.

Failure Rules:
  - Failure MUST produce or propagate Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST NOT emit Execution Plan when planning validation fails.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.

Relationships:
  Parent:
    - ../planning-execution.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/planning.md
    - ../../runtime/error-model.md
  Used By:
    - ../planning-execution.md

References:
  - ../planning-execution.md
  - ../data-plane.md
  - ../../runtime/planning.md
  - ../../runtime/error-model.md
