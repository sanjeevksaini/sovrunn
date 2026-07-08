# Engine Completion

Document:
  ID: engine-completion
  Title: Engine Completion
  Parent: engine-execution
  Owner: Engine Runtime
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Engine Completion stage of Engine Execution
  - Keep engine-execution.md concise
  - Define focused engine-stage rules

Definition:
  Engine Completion defines how Engine Runtime returns canonical engine output to Data Kernel.

Flow:
  - Engine Plugin returns Result Model or Error Model.
  - Engine Runtime validates canonical output shape.
  - Engine Runtime preserves safe plugin execution metadata.
  - Engine Runtime preserves trace and execution identifiers.
  - Engine Runtime returns output to Data Kernel.

Rules:
  - Engine Runtime MUST return canonical output only.
  - Engine Runtime MUST NOT map protocol response.
  - Engine Runtime MUST NOT hide plugin failure.
  - Invalid canonical output MUST produce Error Model.

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
