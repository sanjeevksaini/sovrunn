# Redis Protocol Specification

Document:
  ID: redis-protocol-specification
  Title: Redis Protocol Specification
  Parent: protocol-specification
  Owner: SDE Specification Layer
  Layer: Specification
  Type: Protocol Contract
  Version: 1.0
  Status: Draft

Purpose:
  - Define the Redis protocol integration contract for Sovrunn Data Engine
  - Define Protocol Plugin responsibilities for Redis
  - Define request normalization, session behavior, transaction behavior, result mapping, and error mapping rules
  - Preserve protocol-visible semantics while keeping SDE runtime protocol-neutral

Definition:
  Redis Protocol Specification defines how SDE integrates Redis-compatible clients through a Redis Protocol Plugin while preserving command, key, value, transaction-like, pub/sub, and stream-visible behavior where supported.

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
  Protocol: Redis
  Canonical Identifier: sde.protocol.redis
  Plugin Type: Protocol Plugin

Protocol Plugin Responsibilities:
  MUST:
    - Parse Redis input
    - Validate protocol version compatibility
    - Preserve protocol-visible semantics
    - Produce protocol-normalized intent
    - Preserve session semantics where applicable
    - Preserve transaction semantics where applicable
    - Map Result Model to Redis-compatible response
    - Map Error Model to Redis-compatible error response
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
  - RESP command
  - Key-value command
  - Hash command
  - List command
  - Set command
  - Sorted set command
  - Pub/Sub command where supported
  - Stream command where supported
  - MULTI/EXEC transaction-like command where supported

Request Parsing Rules:
  - Protocol Plugin MUST parse protocol input according to Redis rules.
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
  - Connection-scoped state MUST be declared.
  - Selected logical database behavior MUST be explicit.
  - Authentication context MUST be preserved.
  - Pub/Sub subscription state MUST be isolated by client and tenant.
  - Invalid session context MUST fail closed.

Transaction Rules:
  - MULTI/EXEC behavior MUST be declared when supported.
  - WATCH behavior MUST be declared when supported.
  - Redis transaction-like semantics MUST not be misrepresented as full ACID transaction semantics.
  - Unsupported transaction-like behavior MUST fail deterministically.

Capability Declarations:
  - cache
  - streaming
  - security
  - transactions where MULTI/EXEC is supported

Result Mapping Rules:
  - Protocol Plugin MUST map Result Model to protocol-compatible success response.
  - Protocol response MUST preserve protocol-visible type behavior.
  - Protocol response MUST preserve affected-count or status metadata where applicable.
  - Protocol response MUST preserve cursor, stream, or continuation semantics where applicable.
  - Protocol response MUST not expose raw downstream-native result.
  - Partial result behavior MUST be explicit when supported by protocol.
  - RESP-compatible response shape MUST be preserved.
  - Nil, bulk string, array, integer, and status response behavior MUST be mapped safely.

Error Mapping Rules:
  - Protocol Plugin MUST map Error Model to protocol-compatible error response.
  - Protocol error response MUST preserve safe error semantics.
  - Protocol error response MUST redact unsafe internal details.
  - Trace Identifier MAY be exposed only when protocol and policy allow it.
  - Raw downstream-native error MUST not be exposed.
  - Failure MUST not be represented as success.
  - Redis-compatible error response mapping SHOULD be used where safe and supported.

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
  - Redis Protocol Plugin owns Redis-specific parsing and mapping.
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
