# Result Aggregation

Document:
  ID: result-aggregation
  Title: Result Aggregation
  Parent: kernel-execution
  Owner: Data Kernel
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Result Aggregation stage of Kernel Execution
  - Keep kernel-execution.md concise
  - Define focused kernel-stage rules

Definition:
  Result Aggregation defines how Data Kernel combines operation outputs.

Flow:
  - Data Kernel receives operation output.
  - Data Kernel classifies Result Model, Error Model, or partial output.
  - Data Kernel aggregates successful results according to Execution Plan.
  - Data Kernel preserves partial result state.
  - Data Kernel records failure state where present.

Rules:
  - Result aggregation MUST preserve Execution Plan semantics.
  - Partial result state MUST be explicit.
  - Failure MUST NOT be converted into success.
  - Raw downstream-native result MUST NOT be exposed.

Failure Rules:
  - Failure MUST produce or propagate Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST preserve partial execution state where applicable.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.

Relationships:
  Parent:
    - ../kernel-execution.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/data-kernel.md
    - ../../runtime/error-model.md
  Used By:
    - ../kernel-execution.md

References:
  - ../kernel-execution.md
  - ../data-plane.md
  - ../../runtime/data-kernel.md
  - ../../runtime/error-model.md
