# Downstream Delegation

Document:
  ID: downstream-delegation
  Title: Downstream Delegation
  Parent: request-flow
  Owner: Engine Runtime
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Downstream Delegation stage of SDE Data Plane Request Flow
  - Keep request-flow.md concise
  - Define focused stage rules

Definition:
  Downstream Delegation defines how execution reaches a Downstream Datastore.

Flow:
  - Data Kernel sends execution fragment to Engine Runtime.
  - Engine Runtime resolves approved Engine Plugin.
  - Engine Runtime validates engine and plugin binding.
  - Engine Plugin translates fragment to downstream-native operation.
  - Engine Plugin invokes Downstream Datastore.
  - Datastore Data Plane executes native operation.
  - Engine Plugin maps native output to Result Model or Error Model.

Rules:
  - Engine Runtime MUST NOT invoke Datastore Operator Plugin.
  - Engine Plugin MUST preserve semantic equivalence.
  - Engine Plugin MUST NOT manage datastore lifecycle.
  - Raw native output MUST NOT bypass Result Model or Error Model.

Failure Rules:
  - Failure MUST produce or propagate Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.
  - Failure MUST NOT hide partial or uncertain execution state.

Relationships:
  Parent:
    - ../request-flow.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/runtime.md
  Used By:
    - ../request-flow.md

References:
  - ../request-flow.md
  - ../data-plane.md
  - ../../runtime/runtime.md
