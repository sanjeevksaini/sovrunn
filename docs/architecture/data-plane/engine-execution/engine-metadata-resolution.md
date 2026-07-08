# Engine Metadata Resolution

Document:
  ID: engine-metadata-resolution
  Title: Engine Metadata Resolution
  Parent: engine-execution
  Owner: Engine Runtime
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Engine Metadata Resolution stage of Engine Execution
  - Keep engine-execution.md concise
  - Define focused engine-stage rules

Definition:
  Engine Metadata Resolution defines how Engine Runtime resolves approved engine metadata.

Flow:
  - Engine Runtime reads engine reference from execution request.
  - Engine Runtime resolves approved Engine Registry metadata.
  - Engine Runtime validates engine lifecycle state.
  - Engine Runtime validates tenant and policy compatibility.
  - Engine Runtime resolves downstream endpoint reference where required.

Rules:
  - Engine metadata MUST come from approved Engine Registry state.
  - Engine Runtime MUST NOT invent engine metadata.
  - Unavailable engine metadata MUST fail deterministically.
  - Raw secrets MUST NOT be stored in engine metadata.

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
