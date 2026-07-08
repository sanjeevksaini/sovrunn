# Result Propagation

Document:
  ID: result-propagation
  Title: Result Propagation
  Parent: sde-data-plane
  Owner: SDE Data Plane
  Layer: SDE Data Plane
  Type: Flow
  Version: 1.0
  Status: Stable

Purpose:
  - Define successful and partial result propagation inside SDE Data Plane
  - Define how downstream-native results become SDE Result Model
  - Define how Result Model moves through Engine Runtime, Data Kernel, and Protocol Runtime
  - Define protocol response mapping boundary
  - Point detailed result-stage behavior to focused subflow files

Definition:
  Result Propagation is the SDE Data Plane flow that carries successful or partial execution output from Downstream Datastore through Engine Plugin, Engine Runtime, Data Kernel, Protocol Runtime, and Protocol Plugin to the client.

  Result Propagation begins when an Engine Plugin receives downstream-native successful or partial output.

  Result Propagation ends when Protocol Runtime returns a protocol-compatible response to the client.

  Result Propagation does not handle failures. Failure behavior is defined in error-propagation.md.

Scope:
  In Scope:
    - Downstream native result intake
    - Native result classification
    - Native result mapping to Result Model
    - Engine Runtime result return
    - Data Kernel result aggregation
    - Partial result handling
    - Cursor, stream, and continuation handling
    - Protocol response mapping
    - Client response return

  Out of Scope:
    - Error propagation
    - Protocol request parsing
    - Planning
    - Execution Plan production
    - Downstream datastore lifecycle management
    - Datastore Data Plane internals
    - SDE Control Plane authoritative state mutation

High-Level Flow:
  - Downstream Datastore returns native result.
  - Engine Plugin classifies native result.
  - Engine Plugin maps native result to Result Model.
  - Engine Runtime validates canonical result shape.
  - Engine Runtime returns Result Model to Data Kernel.
  - Data Kernel aggregates operation results according to Execution Plan.
  - Data Kernel preserves partial, cursor, stream, and continuation state.
  - Protocol Runtime receives canonical result.
  - Protocol Plugin maps Result Model to protocol-compatible response.
  - Protocol Runtime returns response to client.

Flow Diagram:
  Downstream Datastore
    ↓
  Datastore Data Plane
    ↓
  Native Result
    ↓
  Engine Plugin
    ↓
  Result Model
    ↓
  Engine Runtime
    ↓
  Data Kernel
    ↓
  Protocol Runtime
    ↓
  Protocol Plugin
    ↓
  Protocol-Compatible Response
    ↓
  Client

Stage Map:
  Native Result Intake:
    Document: result-propagation/native-result-intake.md
    Owner: Engine Plugin

  Result Model Mapping:
    Document: result-propagation/result-model-mapping.md
    Owner: Engine Plugin

  Engine Result Return:
    Document: result-propagation/engine-result-return.md
    Owner: Engine Runtime

  Kernel Result Aggregation:
    Document: result-propagation/kernel-result-aggregation.md
    Owner: Data Kernel

  Protocol Result Mapping:
    Document: result-propagation/protocol-result-mapping.md
    Owner: Protocol Plugin

  Client Response Return:
    Document: result-propagation/client-response-return.md
    Owner: Protocol Runtime

Result Inputs:
  - Downstream native result
  - Execution Context
  - Result shape metadata
  - Execution Plan result expectations
  - Engine Plugin mapping rules
  - Protocol response mapping rules

Result Outputs:
  - Result Model
  - Aggregated Result Model
  - Protocol-compatible response
  - Partial response where supported
  - Cursor, stream, or continuation reference where applicable

Result Model Rules:
  - Native result MUST be mapped to Result Model at Engine Plugin boundary.
  - Result Model MUST preserve result kind.
  - Result Model MUST preserve type information.
  - Result Model MUST preserve schema metadata where applicable.
  - Result Model MUST preserve affected count where applicable.
  - Result Model MUST preserve cursor, stream, or continuation references where applicable.
  - Result Model MUST preserve partial result state where applicable.
  - Result Model MUST NOT contain raw secrets.
  - Result Model MUST NOT expose unsafe downstream-native internals.
  - Result Model MUST NOT represent failure as success.

Engine Plugin Rules:
  - Engine Plugin MUST classify downstream-native result.
  - Engine Plugin MUST map native result to Result Model.
  - Engine Plugin MUST preserve semantic equivalence.
  - Engine Plugin MUST preserve Execution Context correlation.
  - Engine Plugin MUST redact unsafe native details.
  - Engine Plugin MUST NOT return raw native result directly to Protocol Plugin.

Engine Runtime Rules:
  - Engine Runtime MUST accept canonical Result Model from Engine Plugin.
  - Engine Runtime MUST validate result shape where possible.
  - Engine Runtime MUST preserve safe plugin execution metadata.
  - Engine Runtime MUST preserve Trace Identifier and Execution Identifier.
  - Engine Runtime MUST return Result Model to Data Kernel.
  - Engine Runtime MUST NOT map protocol response.

Data Kernel Rules:
  - Data Kernel MUST aggregate Result Model outputs according to Execution Plan.
  - Data Kernel MUST preserve operation dependency semantics.
  - Data Kernel MUST preserve result ordering where required.
  - Data Kernel MUST preserve partial result state.
  - Data Kernel MUST preserve cursor, stream, and continuation references.
  - Data Kernel MUST NOT expose raw downstream-native result.

Protocol Rules:
  - Protocol Plugin MUST map Result Model to protocol-compatible response.
  - Protocol response MUST preserve protocol-visible result semantics.
  - Protocol response MUST preserve type behavior where protocol supports it.
  - Protocol response MUST preserve affected count where applicable.
  - Protocol response MUST preserve cursor, stream, and continuation behavior where applicable.
  - Protocol response MUST NOT expose raw downstream-native result.
  - Protocol response MUST NOT expose unsafe internal metadata.

Partial Result Rules:
  - Partial result state MUST be explicit.
  - Partial result MUST preserve completed operation state where safe.
  - Partial result MUST preserve continuation information where applicable.
  - Partial result MUST NOT hide failed or unknown operation state.
  - Partial result MUST be mapped according to protocol support.
  - Unsupported partial response behavior MUST be mapped safely.

Streaming Rules:
  - Stream identity MUST be request-scoped or cursor-scoped.
  - Stream ordering MUST be preserved where required.
  - Stream lifecycle MUST preserve tenant isolation.
  - Stream continuation MUST not expose raw downstream credentials.
  - Stream closure MUST be observable by runtime.
  - Stream failure after partial output MUST be handled through Error Model.

Security Rules:
  - Result Propagation MUST preserve tenant isolation.
  - Result Propagation MUST protect Execution Context correlation.
  - Result Propagation MUST not expose raw secrets.
  - Result Propagation MUST not expose unsafe downstream-native internals.
  - Result Propagation MUST not expose unauthorized rows, objects, fields, metadata, cursors, or streams.
  - Result Propagation MUST preserve policy-constrained output shape.

Failure Rules:
  - Result mapping failure MUST produce Error Model.
  - Invalid Result Model shape MUST produce Error Model.
  - Result aggregation failure MUST produce Error Model.
  - Protocol result mapping failure MUST produce protocol-compatible internal error.
  - Response delivery failure MUST be recorded in runtime telemetry.
  - Failure MUST NOT be converted into success.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.

Invariants:
  - Native result enters SDE through Engine Plugin only.
  - Result Model is the canonical SDE success representation.
  - Result Propagation never bypasses Result Model.
  - Protocol Plugin is responsible for protocol-specific response mapping.
  - Raw downstream-native result is never exposed directly to client.
  - Failure is not part of Result Model and must use Error Model.
  - Result Propagation does not manage downstream datastore lifecycle.
  - Result Propagation does not mutate SDE Control Plane authoritative state.

Relationships:
  Parent:
    - data-plane.md
  Children:
    - result-propagation/native-result-intake.md
    - result-propagation/result-model-mapping.md
    - result-propagation/engine-result-return.md
    - result-propagation/kernel-result-aggregation.md
    - result-propagation/protocol-result-mapping.md
    - result-propagation/client-response-return.md
  Depends On:
    - data-plane-map.md
    - request-flow.md
    - engine-execution.md
    - kernel-execution.md
    - protocol-execution.md
    - ../runtime/result-model.md
    - ../runtime/error-model.md
    - ../runtime/engine-runtime.md
    - ../runtime/data-kernel.md
    - ../runtime/protocol-runtime.md
  Used By:
    - request-flow.md
    - protocol-execution.md
    - kernel-execution.md
    - engine-execution.md
    - Protocol Plugin specifications
    - Engine Plugin specifications

References:
  - data-plane.md
  - data-plane-map.md
  - request-flow.md
  - engine-execution.md
  - kernel-execution.md
  - protocol-execution.md
  - result-propagation/native-result-intake.md
  - result-propagation/result-model-mapping.md
  - result-propagation/engine-result-return.md
  - result-propagation/kernel-result-aggregation.md
  - result-propagation/protocol-result-mapping.md
  - result-propagation/client-response-return.md
  - ../runtime/result-model.md
  - ../runtime/engine-runtime.md
  - ../runtime/data-kernel.md
  - ../runtime/protocol-runtime.md
