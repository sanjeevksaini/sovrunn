# Capability Resolution

Document:
  ID: capability-resolution
  Title: Capability Resolution
  Parent: planning-execution
  Owner: Planning
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Capability Resolution stage of Planning Execution
  - Keep planning-execution.md concise
  - Define focused planning-stage rules

Definition:
  Capability Resolution defines how Planning validates required capabilities against approved capability metadata.

Flow:
  - Planning extracts required capabilities from validated SIR.
  - Planning reads approved Capability Registry.
  - Planning resolves candidate capability matches.
  - Planning validates capability versions and constraints.
  - Planning records capability decisions for Execution Plan production.

Rules:
  - Planning MUST validate required capabilities.
  - Planning MUST reject unsupported required capabilities.
  - Planning MUST NOT silently downgrade capability requirements.
  - Planning MUST NOT consume unapproved Capability Manifest data.
  - Planning MUST preserve canonical capability identifiers.

Failure Rules:
  - Failure MUST produce or propagate Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST NOT emit Execution Plan when planning validation fails.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.

Relationships:
  Parent:
    - ../planning-execution.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/planning.md
    - ../../runtime/error-model.md
  Used By:
    - ../planning-execution.md

References:
  - ../planning-execution.md
  - ../data-plane.md
  - ../../runtime/planning.md
  - ../../runtime/error-model.md
