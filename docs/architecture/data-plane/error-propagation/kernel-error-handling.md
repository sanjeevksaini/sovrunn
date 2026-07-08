# Kernel Error Handling

Document:
  ID: kernel-error-handling
  Title: Kernel Error Handling
  Parent: error-propagation
  Owner: Data Kernel
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Kernel Error Handling stage of Error Propagation
  - Keep error-propagation.md concise
  - Define focused error-stage rules

Definition:
  Kernel Error Handling defines how Data Kernel handles operation failures during execution orchestration.

Flow:
  - Data Kernel receives Error Model.
  - Data Kernel records failed operation state.
  - Data Kernel evaluates dependent operation impact.
  - Data Kernel preserves partial result state where applicable.
  - Data Kernel stops or continues execution according to Execution Plan.
  - Data Kernel returns final Error Model or partial output.

Rules:
  - Dependent operations MUST stop unless continuation is explicit.
  - Partial failure MUST remain explicit.
  - Unknown operation outcome MUST be explicit.
  - Failure MUST NOT be hidden as success.

Failure Rules:
  - Failure during error handling MUST still produce or preserve Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST NOT expose unsafe details.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.
  - Failure MUST NOT be converted into success.

Relationships:
  Parent:
    - ../error-propagation.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/error-model.md
  Used By:
    - ../error-propagation.md

References:
  - ../error-propagation.md
  - ../data-plane.md
  - ../../runtime/error-model.md
