# Result Model Mapping

Document:
  ID: result-model-mapping
  Title: Result Model Mapping
  Parent: result-propagation
  Owner: Engine Plugin
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Result Model Mapping stage of Result Propagation
  - Keep result-propagation.md concise
  - Define focused result-stage rules

Definition:
  Result Model Mapping defines how Engine Plugin maps native result into SDE Result Model.

Flow:
  - Engine Plugin reads native result.
  - Engine Plugin applies result mapping rules.
  - Engine Plugin maps type and schema metadata.
  - Engine Plugin maps affected count where applicable.
  - Engine Plugin maps cursor, stream, or continuation references.
  - Engine Plugin emits Result Model.

Rules:
  - Result Model MUST preserve semantic equivalence.
  - Result Model MUST not include raw secrets.
  - Unsafe native metadata MUST be redacted.
  - Failure MUST NOT be represented as Result Model.

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
