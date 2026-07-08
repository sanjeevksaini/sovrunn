# Error Propagation

Document:
  ID: error-propagation
  Title: Error Propagation
  Parent: sde-data-plane
  Owner: SDE Data Plane
  Layer: SDE Data Plane
  Type: Flow
  Version: 1.0
  Status: Stable

Purpose:
  - Define failure propagation inside SDE Data Plane
  - Define how runtime and downstream failures become SDE Error Model
  - Define how Error Model moves through Engine Runtime, Data Kernel, Protocol Runtime, and Protocol Plugin
  - Define protocol error response mapping boundary
  - Point detailed error-stage behavior to focused subflow files

Definition:
  Error Propagation is the SDE Data Plane flow that carries execution failure from any runtime component or Downstream Datastore through canonical Error Model and finally to a protocol-compatible client error response.

  Error Propagation begins when a failure is detected by Protocol Runtime, Protocol Plugin, SIR Runtime, Planning, Data Kernel, Engine Runtime, Engine Plugin, or Downstream Datastore.

  Error Propagation ends when Protocol Runtime returns a protocol-compatible error response to the client or records a response-delivery failure.

  Error Propagation does not define successful result behavior. Successful and partial result behavior is defined in result-propagation.md.

Scope:
  In Scope:
    - Runtime failure detection
    - Downstream native error intake
    - Native error mapping
    - Error Model creation
    - Error propagation through Engine Runtime
    - Error propagation through Data Kernel
    - Partial failure state preservation
    - Protocol error mapping
    - Client error response return
    - Unknown outcome handling
    - Retry classification preservation
    - Safe redaction

  Out of Scope:
    - Successful result propagation
    - Protocol request parsing details
    - Planning success path
    - Execution Plan production success path
    - Downstream datastore lifecycle management
    - Datastore Data Plane internals
    - SDE Control Plane authoritative state mutation

High-Level Flow:
  - Failure is detected by runtime component or Downstream Datastore.
  - Detecting component creates or propagates Error Model.
  - Engine Plugin maps downstream-native error to Error Model where applicable.
  - Engine Runtime validates and propagates Error Model.
  - Data Kernel records failure and preserves partial execution state.
  - Protocol Runtime receives Error Model.
  - Protocol Plugin maps Error Model to protocol-compatible error response.
  - Protocol Runtime returns error response to client.

Flow Diagram:
  Runtime Component or Downstream Datastore
    ↓
  Failure Detection
    ↓
  Error Model
    ↓
  Engine Runtime or Data Kernel
    ↓
  Protocol Runtime
    ↓
  Protocol Plugin
    ↓
  Protocol-Compatible Error Response
    ↓
  Client

Stage Map:
  Error Detection:
    Document: error-propagation/error-detection.md
    Owner: Detecting Component

  Native Error Mapping:
    Document: error-propagation/native-error-mapping.md
    Owner: Engine Plugin

  Error Model Creation:
    Document: error-propagation/error-model-creation.md
    Owner: Detecting Component

  Kernel Error Handling:
    Document: error-propagation/kernel-error-handling.md
    Owner: Data Kernel

  Protocol Error Mapping:
    Document: error-propagation/protocol-error-mapping.md
    Owner: Protocol Plugin

  Client Error Return:
    Document: error-propagation/client-error-return.md
    Owner: Protocol Runtime

Error Sources:
  Runtime Sources:
    - Protocol Runtime
    - Protocol Plugin
    - SIR Runtime
    - Planning
    - Data Kernel
    - Engine Runtime
    - Plugin Runtime
    - Session Runtime
    - Transaction Runtime

  Downstream Sources:
    - Engine Plugin
    - Downstream Datastore
    - Datastore Data Plane

  Control Boundary Sources:
    - Missing approved runtime metadata
    - Missing approved plugin metadata
    - Missing approved engine metadata
    - Missing approved capability metadata
    - Missing approved configuration
    - Policy denial

Error Model Requirements:
  Error Model MUST preserve:
    - Error Identifier
    - Code
    - Category
    - Message
    - Severity
    - Source
    - State
    - Retry Classification
    - Trace Identifier
    - Timestamp

  Error Model MAY preserve:
    - Execution Identifier
    - Request Identifier
    - Tenant-safe details
    - Safe cause chain
    - Safe downstream error class
    - Unknown outcome marker
    - Partial execution metadata

Error Classification Rules:
  - Protocol error MUST identify protocol source.
  - SIR validation error MUST identify semantic validation source.
  - Planning error MUST identify planning source.
  - Capability error MUST identify capability source.
  - Policy denial MUST identify policy source.
  - Kernel orchestration error MUST identify Data Kernel source.
  - Engine Runtime error MUST identify Engine Runtime source.
  - Engine Plugin error MUST identify Engine Plugin source.
  - Downstream native error MUST identify downstream source safely.
  - Unknown outcome MUST be explicit.

Retry Classification Rules:
  - Retry classification MUST be explicit.
  - Transient failures SHOULD be classified as retryable only when safe.
  - Non-idempotent operations MUST NOT be retried blindly.
  - Unknown downstream outcome MUST NOT be marked blindly retryable.
  - Policy denial MUST NOT be retryable unless policy state can change.
  - Validation errors SHOULD NOT be retryable without input change.
  - Plugin unavailable errors MAY be retryable when plugin recovery is safe.
  - Downstream timeout MAY be retryable only when operation semantics allow it.

Redaction Rules:
  - Raw secrets MUST NOT appear in Error Model.
  - Downstream credentials MUST NOT appear in Error Model.
  - Unsafe downstream-native details MUST be redacted.
  - Unsafe policy internals MUST be redacted.
  - Unsafe protocol payload details MUST be redacted.
  - Tenant identifiers MUST be exposed only according to policy.
  - Stack traces MUST be exposed only in safe internal contexts.

Propagation Rules:
  - Failure MUST NOT be converted into success.
  - Error Model MUST preserve Trace Identifier.
  - Error Model MUST preserve Timestamp.
  - Error Model MUST preserve safe causal chain.
  - Error Model MUST preserve partial execution state where applicable.
  - Error Model MUST preserve unknown outcome marker where required.
  - Error Model MUST NOT mutate SDE Control Plane authoritative state.
  - Error Model MUST NOT mutate downstream datastore lifecycle state.

Engine Plugin Rules:
  - Engine Plugin MUST map native downstream error to Error Model.
  - Engine Plugin MUST preserve safe native error category.
  - Engine Plugin MUST preserve retry classification where known.
  - Engine Plugin MUST redact unsafe native details.
  - Engine Plugin MUST explicitly report unknown downstream outcome.
  - Engine Plugin MUST NOT return raw native error directly to Protocol Plugin.

Data Kernel Rules:
  - Data Kernel MUST receive and propagate Error Model.
  - Data Kernel MUST record failed operation state.
  - Data Kernel MUST preserve partial result state where applicable.
  - Data Kernel MUST stop dependent operations unless Execution Plan explicitly allows continuation.
  - Data Kernel MUST preserve dependency failure causality.
  - Data Kernel MUST NOT hide partial failure.

Protocol Rules:
  - Protocol Plugin MUST map Error Model to protocol-compatible error response.
  - Protocol error response MUST preserve safe error semantics.
  - Protocol error response MUST redact unsafe details.
  - Protocol error response MAY expose Trace Identifier only when protocol and policy allow it.
  - Protocol error response MUST NOT expose raw downstream-native error.
  - Protocol error response MUST NOT represent failure as success.

Observability Rules:
  - Error Propagation MUST preserve Trace Identifier.
  - Error Propagation MUST preserve Request Identifier where available.
  - Error Propagation MUST preserve Execution Identifier where available.
  - Error telemetry MUST be tenant-safe.
  - Error telemetry MUST NOT expose secrets.
  - Error telemetry MUST NOT expose unsafe native details.
  - Runtime telemetry MUST NOT replace Audit Service.

Failure Rules:
  - Error mapping failure MUST still produce Error Model.
  - Protocol error mapping failure MUST produce protocol-compatible internal error.
  - Response delivery failure MUST be recorded in runtime telemetry.
  - Unknown failure MUST produce Error Model.
  - Unknown outcome MUST be explicit.
  - Failure MUST NOT corrupt SDE Control Plane authoritative state.
  - Failure MUST NOT corrupt downstream datastore lifecycle state.

Invariants:
  - Error Model is the canonical SDE failure representation.
  - Every failed request has Error Model.
  - Error Model Timestamp is mandatory.
  - Failure is never represented as Result Model.
  - Raw downstream-native error is never exposed directly to client.
  - Error Propagation does not manage downstream datastore lifecycle.
  - Error Propagation does not mutate SDE Control Plane authoritative state.
  - Error Propagation preserves safe traceability.

Relationships:
  Parent:
    - data-plane.md
  Children:
    - error-propagation/error-detection.md
    - error-propagation/native-error-mapping.md
    - error-propagation/error-model-creation.md
    - error-propagation/kernel-error-handling.md
    - error-propagation/protocol-error-mapping.md
    - error-propagation/client-error-return.md
  Depends On:
    - data-plane-map.md
    - request-flow.md
    - protocol-execution.md
    - planning-execution.md
    - kernel-execution.md
    - engine-execution.md
    - result-propagation.md
    - ../runtime/error-model.md
    - ../runtime/result-model.md
    - ../runtime/engine-runtime.md
    - ../runtime/data-kernel.md
    - ../runtime/protocol-runtime.md
  Used By:
    - request-flow.md
    - protocol-execution.md
    - planning-execution.md
    - kernel-execution.md
    - engine-execution.md
    - Protocol Plugin specifications
    - Engine Plugin specifications

References:
  - data-plane.md
  - data-plane-map.md
  - request-flow.md
  - protocol-execution.md
  - planning-execution.md
  - kernel-execution.md
  - engine-execution.md
  - result-propagation.md
  - error-propagation/error-detection.md
  - error-propagation/native-error-mapping.md
  - error-propagation/error-model-creation.md
  - error-propagation/kernel-error-handling.md
  - error-propagation/protocol-error-mapping.md
  - error-propagation/client-error-return.md
  - ../runtime/error-model.md
  - ../runtime/result-model.md
  - ../runtime/engine-runtime.md
  - ../runtime/data-kernel.md
  - ../runtime/protocol-runtime.md
