# Execution Fragment Translation

Document:
  ID: execution-fragment-translation
  Title: Execution Fragment Translation
  Parent: engine-execution
  Owner: Engine Plugin
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Execution Fragment Translation stage of Engine Execution
  - Keep engine-execution.md concise
  - Define focused engine-stage rules

Definition:
  Execution Fragment Translation defines how Engine Plugin converts SDE execution fragment into downstream-native operation.

Flow:
  - Engine Plugin receives execution fragment.
  - Engine Plugin validates capability boundary.
  - Engine Plugin reads downstream dialect or API mapping.
  - Engine Plugin translates fragment into downstream-native operation.
  - Engine Plugin validates semantic equivalence.

Rules:
  - Translation MUST preserve semantic equivalence.
  - Unsupported native operation MUST fail deterministically.
  - Engine Plugin MUST NOT silently emulate unsupported capability.
  - Downstream-native operation MUST NOT become SDE platform contract.

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
