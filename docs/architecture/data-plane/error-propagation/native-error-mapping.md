# Native Error Mapping

Document:
  ID: native-error-mapping
  Title: Native Error Mapping
  Parent: error-propagation
  Owner: Engine Plugin
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Native Error Mapping stage of Error Propagation
  - Keep error-propagation.md concise
  - Define focused error-stage rules

Definition:
  Native Error Mapping defines how Engine Plugin maps downstream-native errors into SDE Error Model.

Flow:
  - Downstream Datastore returns native error.
  - Engine Plugin receives native error.
  - Engine Plugin classifies native error safely.
  - Engine Plugin determines retry classification where possible.
  - Engine Plugin maps native error to Error Model.
  - Engine Plugin redacts unsafe native details.

Rules:
  - Native error MUST map to Error Model.
  - Raw native error MUST NOT bypass Error Model.
  - Unsafe native details MUST be redacted.
  - Unknown downstream outcome MUST be explicit.

Failure Rules:
  - Failure during error handling MUST still produce or preserve Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST NOT expose unsafe details.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.
  - Failure MUST NOT be converted into success.

Relationships:
  Parent:
    - ../error-propagation.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/error-model.md
  Used By:
    - ../error-propagation.md

References:
  - ../error-propagation.md
  - ../data-plane.md
  - ../../runtime/error-model.md
