# Kernel Result Aggregation

Document:
  ID: kernel-result-aggregation
  Title: Kernel Result Aggregation
  Parent: result-propagation
  Owner: Data Kernel
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Kernel Result Aggregation stage of Result Propagation
  - Keep result-propagation.md concise
  - Define focused result-stage rules

Definition:
  Kernel Result Aggregation defines how Data Kernel aggregates operation results.

Flow:
  - Data Kernel receives Result Model from Engine Runtime.
  - Data Kernel updates operation graph state.
  - Data Kernel aggregates result according to Execution Plan.
  - Data Kernel preserves ordering where required.
  - Data Kernel preserves partial, cursor, stream, or continuation state.

Rules:
  - Aggregation MUST preserve Execution Plan semantics.
  - Partial result state MUST be explicit.
  - Result ordering MUST be preserved where required.
  - Data Kernel MUST NOT expose raw downstream-native result.

Failure Rules:
  - Failure MUST produce or propagate Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST NOT expose raw downstream-native result.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.

Relationships:
  Parent:
    - ../result-propagation.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/result-model.md
    - ../../runtime/error-model.md
  Used By:
    - ../result-propagation.md

References:
  - ../result-propagation.md
  - ../data-plane.md
  - ../../runtime/result-model.md
  - ../../runtime/error-model.md
