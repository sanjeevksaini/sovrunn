# Context Binding

Document:
  ID: context-binding
  Title: Context Binding
  Parent: kernel-execution
  Owner: Data Kernel
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Context Binding stage of Kernel Execution
  - Keep kernel-execution.md concise
  - Define focused kernel-stage rules

Definition:
  Context Binding defines how Execution Context is attached to Kernel Execution.

Flow:
  - SDE Runtime provides Execution Context.
  - Data Kernel validates Execution Identifier.
  - Data Kernel validates Request Identifier and Trace Identifier.
  - Data Kernel validates tenant and security context.
  - Data Kernel binds session and transaction references where present.

Rules:
  - Execution Context MUST be immutable during execution.
  - Execution Context MUST not contain raw secrets.
  - Tenant mismatch MUST fail closed.
  - Invalid session or transaction reference MUST fail closed.

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
