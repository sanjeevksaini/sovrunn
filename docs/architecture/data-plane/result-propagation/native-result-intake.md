# Native Result Intake

Document:
  ID: native-result-intake
  Title: Native Result Intake
  Parent: result-propagation
  Owner: Engine Plugin
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Native Result Intake stage of Result Propagation
  - Keep result-propagation.md concise
  - Define focused result-stage rules

Definition:
  Native Result Intake defines how Engine Plugin receives downstream-native successful or partial output.

Flow:
  - Downstream Datastore completes native operation.
  - Datastore Data Plane returns native result.
  - Engine Plugin receives native result.
  - Engine Plugin classifies result kind.
  - Engine Plugin preserves execution correlation.

Rules:
  - Native result MUST enter SDE through Engine Plugin only.
  - Engine Plugin MUST classify success, partial success, stream, cursor, or continuation.
  - Engine Plugin MUST preserve Trace Identifier where available.
  - Engine Plugin MUST NOT expose native result directly.

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
