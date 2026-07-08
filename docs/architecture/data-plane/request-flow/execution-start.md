# Execution Start

Document:
  ID: execution-start
  Title: Execution Start
  Parent: request-flow
  Owner: Data Kernel
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Execution Start stage of SDE Data Plane Request Flow
  - Keep request-flow.md concise
  - Define focused stage rules

Definition:
  Execution Start defines how Data Kernel begins Execution Plan execution.

Flow:
  - Data Kernel receives Execution Plan.
  - SDE Runtime creates Execution Context.
  - Data Kernel validates plan and context readiness.
  - Data Kernel initializes in-flight execution state.
  - Data Kernel determines executable operation set.
  - Data Kernel prepares Engine Runtime requests.

Rules:
  - Execution Plan MUST be immutable.
  - Execution Context MUST be immutable during execution.
  - Data Kernel MUST preserve operation dependency order.
  - Data Kernel MUST NOT access Downstream Datastore directly.

Failure Rules:
  - Failure MUST produce or propagate Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.
  - Failure MUST NOT hide partial or uncertain execution state.

Relationships:
  Parent:
    - ../request-flow.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/runtime.md
  Used By:
    - ../request-flow.md

References:
  - ../request-flow.md
  - ../data-plane.md
  - ../../runtime/runtime.md
