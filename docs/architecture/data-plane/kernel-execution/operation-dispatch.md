# Operation Dispatch

Document:
  ID: operation-dispatch
  Title: Operation Dispatch
  Parent: kernel-execution
  Owner: Data Kernel
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Operation Dispatch stage of Kernel Execution
  - Keep kernel-execution.md concise
  - Define focused kernel-stage rules

Definition:
  Operation Dispatch defines how Data Kernel delegates executable operations to Engine Runtime.

Flow:
  - Data Kernel selects executable operation.
  - Data Kernel builds operation execution request.
  - Data Kernel attaches Execution Context.
  - Data Kernel invokes Engine Runtime.
  - Engine Runtime returns Result Model or Error Model.

Rules:
  - Data Kernel MUST dispatch through Engine Runtime.
  - Data Kernel MUST NOT invoke Engine Plugin directly.
  - Data Kernel MUST NOT access Downstream Datastore directly.
  - Data Kernel MUST pass only plan-authorized execution fragments.

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
