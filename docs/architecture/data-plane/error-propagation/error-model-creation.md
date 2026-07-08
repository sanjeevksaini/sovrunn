# Error Model Creation

Document:
  ID: error-model-creation
  Title: Error Model Creation
  Parent: error-propagation
  Owner: Detecting Component
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Error Model Creation stage of Error Propagation
  - Keep error-propagation.md concise
  - Define focused error-stage rules

Definition:
  Error Model Creation defines mandatory fields and creation rules for canonical SDE failures.

Flow:
  - Detecting component receives failure context.
  - Component assigns Error Identifier.
  - Component sets Code, Category, Message, Severity, Source, and State.
  - Component sets Retry Classification.
  - Component preserves Trace Identifier.
  - Component sets Timestamp at detection time.
  - Component attaches safe details and safe cause chain.

Rules:
  - Timestamp is mandatory.
  - Trace Identifier MUST be preserved when available.
  - Raw secrets MUST NOT be included.
  - Unsafe details MUST be redacted.
  - Unknown failure MUST still produce Error Model.

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
