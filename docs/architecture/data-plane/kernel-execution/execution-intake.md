# Execution Intake

Document:
  ID: execution-intake
  Title: Execution Intake
  Parent: kernel-execution
  Owner: Data Kernel
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Execution Intake stage of Kernel Execution
  - Keep kernel-execution.md concise
  - Define focused kernel-stage rules

Definition:
  Execution Intake defines how Data Kernel receives and validates an Execution Plan before execution begins.

Flow:
  - Planning emits immutable Execution Plan.
  - Data Kernel receives Execution Plan.
  - Data Kernel validates plan identity and version.
  - Data Kernel validates operation graph presence.
  - Data Kernel rejects invalid or incomplete plan.

Rules:
  - Data Kernel MUST consume Execution Plan only after Planning completes.
  - Execution Plan MUST be immutable.
  - Invalid Execution Plan MUST fail before downstream delegation.
  - Data Kernel MUST NOT modify plan semantics.

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
