# Kernel Completion

Document:
  ID: kernel-completion
  Title: Kernel Completion
  Parent: kernel-execution
  Owner: Data Kernel
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Kernel Completion stage of Kernel Execution
  - Keep kernel-execution.md concise
  - Define focused kernel-stage rules

Definition:
  Kernel Completion defines how Data Kernel completes execution and returns canonical output.

Flow:
  - Data Kernel determines final operation graph state.
  - Data Kernel validates completion criteria.
  - Data Kernel finalizes aggregated result or error.
  - Data Kernel preserves trace and execution identifiers.
  - Data Kernel returns canonical output to protocol response path.

Rules:
  - Completion MUST reflect success, partial success, failure, or unknown outcome accurately.
  - Unknown outcome MUST be explicit.
  - Final output MUST be Result Model or Error Model.
  - Completion MUST NOT mutate SDE Control Plane authoritative state.

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
