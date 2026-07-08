# Engine Request Intake

Document:
  ID: engine-request-intake
  Title: Engine Request Intake
  Parent: engine-execution
  Owner: Engine Runtime
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Engine Request Intake stage of Engine Execution
  - Keep engine-execution.md concise
  - Define focused engine-stage rules

Definition:
  Engine Request Intake defines how Engine Runtime receives an operation execution request from Data Kernel.

Flow:
  - Data Kernel creates operation execution request.
  - Data Kernel attaches Execution Context.
  - Engine Runtime receives execution request.
  - Engine Runtime validates request identity.
  - Engine Runtime validates execution fragment presence.

Rules:
  - Engine Runtime MUST accept requests from Data Kernel only.
  - Execution Context MUST be present.
  - Execution fragment MUST be authorized by Execution Plan.
  - Invalid request MUST fail before plugin resolution.

Failure Rules:
  - Failure MUST produce or propagate Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST preserve retry classification where applicable.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.
  - Failure MUST NOT mutate datastore lifecycle state.

Relationships:
  Parent:
    - ../engine-execution.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/engine-runtime.md
    - ../../runtime/error-model.md
  Used By:
    - ../engine-execution.md

References:
  - ../engine-execution.md
  - ../data-plane.md
  - ../../runtime/engine-runtime.md
  - ../../runtime/error-model.md
