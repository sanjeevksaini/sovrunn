# Dependency Evaluation

Document:
  ID: dependency-evaluation
  Title: Dependency Evaluation
  Parent: kernel-execution
  Owner: Data Kernel
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Dependency Evaluation stage of Kernel Execution
  - Keep kernel-execution.md concise
  - Define focused kernel-stage rules

Definition:
  Dependency Evaluation defines how Data Kernel determines executable operations from the Execution Plan.

Flow:
  - Data Kernel reads operation graph.
  - Data Kernel identifies dependency-free operations.
  - Data Kernel evaluates completed predecessor operations.
  - Data Kernel marks executable operations.
  - Data Kernel blocks dependent operations until prerequisites complete.

Rules:
  - Dependency ordering MUST be preserved.
  - Independent operations MAY run concurrently only when plan allows it.
  - Failed dependency MUST block dependent operations unless continuation is explicit.
  - Dependency evaluation MUST be deterministic.

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
