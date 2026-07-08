# SIR Intake

Document:
  ID: sir-intake
  Title: SIR Intake
  Parent: planning-execution
  Owner: Planning
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define SIR Intake stage of Planning Execution
  - Keep planning-execution.md concise
  - Define focused planning-stage rules

Definition:
  SIR Intake defines how Planning receives and accepts validated SIR from SIR Runtime.

Flow:
  - SIR Runtime emits validated SIR.
  - Planning receives validated SIR and request context.
  - Planning verifies SIR validation status.
  - Planning verifies SIR version compatibility.
  - Planning extracts operation, resource, relationship, expression, and constraint references required for planning.

Rules:
  - Planning MUST accept validated SIR only.
  - Planning MUST reject unvalidated SIR.
  - Planning MUST preserve SIR semantic intent.
  - Planning MUST NOT mutate SIR.
  - Planning MUST NOT perform protocol parsing.

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
