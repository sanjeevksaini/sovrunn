# gRPC Protocol Specification

Document:
  ID: grpc-protocol-specification
  Title: gRPC Protocol Specification
  Parent: protocol-specification
  Owner: SDE Specification Layer
  Layer: Specification
  Type: Protocol Contract
  Version: 1.0
  Status: Draft

Purpose:
  - Define the gRPC protocol integration contract for Sovrunn Data Engine
  - Define Protocol Plugin responsibilities for gRPC
  - Define request normalization, session behavior, transaction behavior, result mapping, and error mapping rules
  - Preserve protocol-visible semantics while keeping SDE runtime protocol-neutral

Definition:
  gRPC Protocol Specification defines how SDE integrates gRPC clients through a gRPC Protocol Plugin while preserving service, method, metadata, message, streaming, status, and error semantics.

Scope:
  In Scope:
    - Protocol identity
    - Protocol version compatibility
    - Protocol Plugin manifest requirements
    - Request parsing and validation
    - Protocol-normalized intent creation
    - Session semantics
    - Transaction semantics
    - Result Model mapping
    - Error Model mapping
    - Capability declaration

  Out of Scope:
    - SIR internals
    - Planning internals
    - Data Kernel orchestration
    - Engine Runtime delegation
    - Engine Plugin implementation
    - Downstream datastore lifecycle management
    - Datastore Data Plane internals

Protocol Identity:
  Protocol: gRPC
  Canonical Identifier: sde.protocol.grpc
  Plugin Type: Protocol Plugin

Protocol Plugin Responsibilities:
  MUST:
    - Parse gRPC input
    - Validate protocol version compatibility
    - Preserve protocol-visible semantics
    - Produce protocol-normalized intent
    - Preserve session semantics where applicable
    - Preserve transaction semantics where applicable
    - Map Result Model to gRPC-compatible response
    - Map Error Model to gRPC-compatible error response
    - Redact unsafe internal details

  MUST NOT:
    - Select Downstream Datastore directly
    - Produce Execution Plan
    - Invoke Engine Runtime
    - Invoke Engine Plugin
    - Invoke Datastore Operator Plugin
    - Manage datastore lifecycle
    - Access Datastore Data Plane directly
    - Modify SDE Control Plane authoritative state

Plugin Manifest Requirements:
  Manifest MUST include:
    - Protocol identifier
    - Protocol version support
    - Plugin identifier
    - Plugin version
    - Supported request forms
    - Unsupported request forms
    - Session behavior
    - Transaction behavior
    - Result mapping behavior
    - Error mapping behavior
    - Capability declarations
    - Compatibility metadata
    - Known semantic gaps

Supported Request Forms:
  - Unary RPC
  - Server streaming RPC
  - Client streaming RPC
  - Bidirectional streaming RPC
  - Metadata handling
  - Deadline handling
  - Cancellation handling

Request Parsing Rules:
  - Protocol Plugin MUST parse protocol input according to gRPC rules.
  - Protocol Plugin MUST reject malformed input deterministically.
  - Protocol Plugin MUST preserve client-visible semantics.
  - Protocol Plugin MUST preserve protocol feature flags where applicable.
  - Unsupported feature MUST fail deterministically unless explicit compatibility behavior exists.
  - Protocol Plugin MUST NOT perform semantic optimization.
  - Protocol Plugin MUST NOT produce Execution Plan.

Protocol-Normalized Intent:
  Protocol-normalized intent MUST include:
    - Protocol identity
    - Protocol operation kind
    - Parsed request structure
    - Protocol semantic modifiers
    - Session reference where applicable
    - Transaction intent where applicable
    - Request metadata
    - Tenant context when available
    - Security context when available

  Protocol-normalized intent MUST NOT include:
    - Raw secrets
    - Engine Plugin invocation data
    - Downstream datastore credentials
    - Downstream-native execution commands
    - SDE Control Plane mutation commands

Session Rules:
  - gRPC metadata MUST be interpreted safely.
  - Authentication context MUST be extracted from approved metadata.
  - Request context MUST preserve deadline and cancellation where supported.
  - Stream context MUST be isolated by tenant and request.
  - Invalid context MUST fail closed.

Transaction Rules:
  - gRPC method semantics MUST not imply transaction semantics unless specified by service contract.
  - Transaction reference MAY be represented in metadata or message only when approved.
  - Unsupported transaction behavior MUST fail deterministically.
  - Cancellation and deadline behavior MUST be mapped safely.

Capability Declarations:
  - security
  - streaming
  - object
  - federation

Result Mapping Rules:
  - Result Model MUST map to gRPC response message or stream.
  - Streaming response ordering MUST be preserved.
  - gRPC trailers and metadata MUST be safe.
  - Partial output behavior MUST be explicit.
  - Raw downstream-native result MUST not be exposed.

Error Mapping Rules:
  - Error Model MUST map to gRPC status and safe error details.
  - Unsafe details MUST be redacted.
  - Deadline exceeded and cancellation MUST map safely.
  - Failure MUST not be represented as success.

Security Rules:
  - Protocol Plugin MUST preserve tenant isolation.
  - Protocol Plugin MUST protect request context.
  - Protocol Plugin MUST protect session and transaction references.
  - Protocol Plugin MUST NOT expose raw secrets.
  - Protocol Plugin MUST NOT expose raw downstream-native result.
  - Protocol Plugin MUST NOT expose raw downstream-native error.
  - Protocol Plugin MUST redact unsafe internal details.
  - Protocol Plugin MUST comply with policy-constrained response shape.

Compatibility Rules:
  - Protocol version compatibility MUST follow Versioning Specification.
  - Protocol Plugin MUST declare supported protocol versions.
  - Unsupported protocol version MUST fail deterministically.
  - Deprecated protocol behavior MUST remain explicit.
  - Compatibility behavior MUST not silently change protocol-visible semantics.

Observability Rules:
  - Protocol Plugin SHOULD emit safe parse and mapping telemetry.
  - Protocol Runtime MUST preserve Request Identifier.
  - Protocol Runtime MUST preserve Trace Identifier.
  - Telemetry MUST NOT contain raw secrets.
  - Telemetry MUST NOT contain unsafe protocol payloads.
  - Telemetry MUST NOT expose unsafe downstream-native details.

Failure Rules:
  - Protocol parsing failure MUST produce Error Model.
  - Unsupported protocol feature MUST produce Error Model.
  - Session resolution failure MUST fail closed.
  - Transaction resolution failure MUST fail closed.
  - Result mapping failure MUST produce protocol-compatible internal error.
  - Error mapping failure MUST produce protocol-compatible internal error.
  - Failure MUST NOT be converted into success.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.

Invariants:
  - gRPC Protocol Plugin owns gRPC-specific parsing and mapping.
  - Protocol Plugin produces protocol-normalized intent, not SIR.
  - Protocol Plugin does not produce Execution Plan.
  - Protocol Plugin does not invoke Engine Plugin.
  - Protocol Plugin does not manage datastore lifecycle.
  - Protocol Plugin maps canonical SDE models to protocol-compatible responses.
  - Raw downstream-native output is never exposed directly to client.

Relationships:
  Parent:
    - protocol.md
  Depends On:
    - protocol.md
    - ../versioning/versioning.md
    - ../serialization/serialization.md
    - ../capability/capability.md
    - ../../architecture/data-plane/protocol-execution.md
    - ../../architecture/runtime/protocol-runtime.md
    - ../../architecture/runtime/session-runtime.md
    - ../../architecture/runtime/transaction-runtime.md
    - ../../architecture/runtime/result-model.md
    - ../../architecture/runtime/error-model.md
  Used By:
    - Protocol Runtime
    - Protocol Plugin
    - SIR Runtime
    - Planning
    - Result Propagation
    - Error Propagation

References:
  - protocol.md
  - ../versioning/versioning.md
  - ../serialization/serialization.md
  - ../capability/capability.md
  - ../../architecture/data-plane/protocol-execution.md
  - ../../architecture/runtime/protocol-runtime.md
  - ../../architecture/runtime/result-model.md
  - ../../architecture/runtime/error-model.md
