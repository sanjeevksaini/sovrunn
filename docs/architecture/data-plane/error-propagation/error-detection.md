# Error Detection

Document:
  ID: error-detection
  Title: Error Detection
  Parent: error-propagation
  Owner: Detecting Component
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Error Detection stage of Error Propagation
  - Keep error-propagation.md concise
  - Define focused error-stage rules

Definition:
  Error Detection defines how runtime components detect and classify failures before Error Model creation or propagation.

Flow:
  - Runtime component detects failure.
  - Component classifies source and failure category.
  - Component captures safe context.
  - Component preserves Trace Identifier where available.
  - Component creates or forwards Error Model input.

Rules:
  - Failure MUST be detected at source where possible.
  - Failure source MUST be preserved.
  - Unsafe internal details MUST not be captured for client exposure.
  - Failure MUST NOT be converted into success.

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
