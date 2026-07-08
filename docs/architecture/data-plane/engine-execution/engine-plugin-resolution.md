# Engine Plugin Resolution

Document:
  ID: engine-plugin-resolution
  Title: Engine Plugin Resolution
  Parent: engine-execution
  Owner: Engine Runtime
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Engine Plugin Resolution stage of Engine Execution
  - Keep engine-execution.md concise
  - Define focused engine-stage rules

Definition:
  Engine Plugin Resolution defines how Engine Runtime resolves and validates an approved Engine Plugin.

Flow:
  - Engine Runtime reads Engine Plugin binding.
  - Engine Runtime resolves approved Plugin Registry metadata.
  - Engine Runtime validates plugin lifecycle state.
  - Engine Runtime validates plugin compatibility.
  - Engine Runtime prepares plugin invocation context.

Rules:
  - Only approved Engine Plugins may be used.
  - Plugin compatibility MUST be validated.
  - Unavailable plugin MUST fail with retry classification when retry is safe.
  - Engine Runtime MUST NOT invoke Datastore Operator Plugin.

Failure Rules:
  - Failure MUST produce or propagate Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST preserve retry classification where applicable.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.
  - Failure MUST NOT mutate datastore lifecycle state.

Relationships:
  Parent:
    - ../engine-execution.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/engine-runtime.md
    - ../../runtime/error-model.md
  Used By:
    - ../engine-execution.md

References:
  - ../engine-execution.md
  - ../data-plane.md
  - ../../runtime/engine-runtime.md
  - ../../runtime/error-model.md
