# Native Output Mapping

Document:
  ID: native-output-mapping
  Title: Native Output Mapping
  Parent: engine-execution
  Owner: Engine Plugin
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Native Output Mapping stage of Engine Execution
  - Keep engine-execution.md concise
  - Define focused engine-stage rules

Definition:
  Native Output Mapping defines how Engine Plugin maps downstream native output into canonical SDE models.

Flow:
  - Engine Plugin receives native result or native error.
  - Engine Plugin classifies output as success, partial success, failure, or unknown outcome.
  - Engine Plugin maps native result to Result Model.
  - Engine Plugin maps native error to Error Model.
  - Engine Plugin redacts unsafe native details.

Rules:
  - Native result MUST map to Result Model.
  - Native error MUST map to Error Model.
  - Raw native output MUST NOT bypass canonical models.
  - Unknown outcome MUST be explicit.

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
