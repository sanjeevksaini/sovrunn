# Engine Result Return

Document:
  ID: engine-result-return
  Title: Engine Result Return
  Parent: result-propagation
  Owner: Engine Runtime
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Engine Result Return stage of Result Propagation
  - Keep result-propagation.md concise
  - Define focused result-stage rules

Definition:
  Engine Result Return defines how Engine Runtime returns canonical Result Model to Data Kernel.

Flow:
  - Engine Plugin returns Result Model.
  - Engine Runtime validates canonical result shape.
  - Engine Runtime preserves safe plugin metadata.
  - Engine Runtime preserves execution correlation.
  - Engine Runtime returns Result Model to Data Kernel.

Rules:
  - Engine Runtime MUST accept canonical Result Model only.
  - Invalid result shape MUST produce Error Model.
  - Engine Runtime MUST NOT map protocol response.
  - Engine Runtime MUST NOT expose raw native result.

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
