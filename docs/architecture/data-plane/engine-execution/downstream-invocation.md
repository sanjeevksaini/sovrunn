# Downstream Invocation

Document:
  ID: downstream-invocation
  Title: Downstream Invocation
  Parent: engine-execution
  Owner: Engine Plugin
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Downstream Invocation stage of Engine Execution
  - Keep engine-execution.md concise
  - Define focused engine-stage rules

Definition:
  Downstream Invocation defines how Engine Plugin invokes a Downstream Datastore through approved interface.

Flow:
  - Engine Plugin resolves downstream endpoint reference.
  - Engine Plugin obtains authorized credential reference where required.
  - Engine Plugin invokes approved downstream interface.
  - Downstream Datastore executes through Datastore Data Plane.
  - Engine Plugin receives native result or native error.

Rules:
  - Engine Plugin MUST protect downstream credentials.
  - Engine Plugin MUST preserve tenant boundary.
  - Engine Plugin MUST NOT bypass approved downstream interface.
  - Engine Plugin MUST NOT manage datastore lifecycle.

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
