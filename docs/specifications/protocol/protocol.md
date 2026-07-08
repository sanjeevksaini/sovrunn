# Protocol Specification

Document:
  ID: protocol-specification
  Title: Protocol Specification
  Parent: specifications
  Owner: SDE Specification Layer
  Layer: Specification
  Type: Contract
  Version: 1.0
  Status: Draft

Purpose:
  - Define how client protocols integrate with SDE Data Plane
  - Define Protocol Plugin contract boundaries
  - Define protocol request normalization and response mapping rules
  - Preserve protocol transparency without coupling SDE semantics to one protocol

Definition:
  Protocol Specification defines the contract by which a client protocol is integrated into SDE through a Protocol Plugin.

  A Protocol Plugin parses protocol input, preserves protocol-visible semantics, produces protocol-normalized intent, and maps Result Model or Error Model into protocol-compatible response.

Scope:
  In Scope:
    - Protocol identity
    - Protocol version compatibility
    - Protocol Plugin manifest
    - Protocol request parsing
    - Protocol-normalized intent
    - Session semantics
    - Transaction semantics
    - Result response mapping
    - Error response mapping
    - Protocol capability declaration

  Out of Scope:
    - Engine Plugin execution
    - Downstream datastore native execution
    - SDE Planning internals
    - Data Kernel orchestration
    - Datastore lifecycle management
    - Datastore Operator Plugin behavior

Protocol Plugin Responsibilities:
  MUST:
    - Parse protocol input
    - Validate protocol version compatibility
    - Preserve protocol-visible semantics
    - Produce protocol-normalized intent
    - Interpret protocol session behavior where applicable
    - Interpret protocol transaction commands where applicable
    - Map Result Model to protocol-compatible response
    - Map Error Model to protocol-compatible error response
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

Protocol Plugin Manifest:
  Manifest MUST include:
    - Protocol identifier
    - Protocol version support
    - Plugin identifier
    - Plugin version
    - Supported client features
    - Unsupported client features
    - Session behavior
    - Transaction behavior
    - Result mapping behavior
    - Error mapping behavior
    - Capability declarations
    - Compatibility metadata

Protocol-Normalized Intent:
  Protocol-normalized intent is an intermediate protocol-layer output consumed by SIR Runtime.

  It is:
    - Protocol-aware
    - Request-scoped
    - Semantics-preserving
    - Suitable for SIR Runtime input

  It is not:
    - SIR
    - Execution Plan
    - Downstream-native operation
    - Engine Plugin invocation

Protocol-Normalized Intent MUST include:
  - Protocol identity
  - Protocol operation kind
  - Parsed request structure
  - Protocol semantic modifiers
  - Session reference where applicable
  - Transaction intent where applicable
  - Request metadata
  - Tenant context when available
  - Security context when available

Parsing Rules:
  - Protocol Plugin MUST parse according to protocol rules.
  - Protocol Plugin MUST reject malformed protocol input deterministically.
  - Unsupported protocol feature MUST fail deterministically unless explicit compatibility behavior exists.
  - Protocol Plugin MUST preserve protocol feature flags where applicable.
  - Protocol Plugin MUST not perform semantic optimization.

Session Rules:
  - Protocol Plugin MAY interpret protocol-visible session behavior.
  - Session Runtime owns session context.
  - Protocol Runtime manages protocol session boundary.
  - Session state MUST be tenant-isolated.
  - Invalid session MUST fail closed.

Transaction Rules:
  - Protocol Plugin MAY interpret protocol-visible transaction commands.
  - Transaction Runtime owns transaction context.
  - Protocol Runtime preserves transaction references.
  - Unsupported transaction semantics MUST fail deterministically.
  - Protocol Plugin MUST NOT execute transaction operations directly against Downstream Datastore.

Result Mapping Rules:
  - Protocol Plugin MUST map Result Model to protocol-compatible response.
  - Protocol response MUST preserve protocol-visible type behavior.
  - Protocol response MUST preserve affected count where applicable.
  - Protocol response MUST preserve cursor, stream, or continuation semantics where applicable.
  - Raw downstream-native result MUST NOT be exposed.

Error Mapping Rules:
  - Protocol Plugin MUST map Error Model to protocol-compatible error response.
  - Protocol error response MUST preserve safe error semantics.
  - Protocol error response MUST redact unsafe details.
  - Trace Identifier MAY be exposed only when protocol and policy allow it.
  - Raw downstream-native error MUST NOT be exposed.
  - Failure MUST NOT be represented as success.

Capability Rules:
  - Protocol Plugin MUST declare supported protocol capabilities.
  - Protocol capability support MUST be explicit.
  - Unsupported required protocol feature MUST fail deterministically.
  - Protocol capabilities MUST be validated by Planning where they affect execution.

Invariants:
  - Protocol Plugin is required for protocol-specific behavior.
  - Protocol Specification does not define downstream datastore execution.
  - Protocol Plugin does not produce Execution Plan.
  - Protocol Plugin does not invoke Engine Plugin.
  - Protocol Plugin does not manage datastore lifecycle.
  - Protocol Plugin maps canonical models to protocol response.

Relationships:
  Parent:
    - ../specifications
  Children:
    - postgresql.md
    - mysql.md
    - mongodb.md
    - redis.md
    - rest.md
    - grpc.md
    - native.md
  Depends On:
    - ../versioning/versioning.md
    - ../serialization/serialization.md
    - ../capability/capability.md
    - ../../architecture/data-plane/protocol-execution.md
    - ../../architecture/runtime/protocol-runtime.md
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
  - ../versioning/versioning.md
  - ../serialization/serialization.md
  - ../capability/capability.md
  - ../../architecture/data-plane/protocol-execution.md
  - ../../architecture/runtime/protocol-runtime.md
  - ../../architecture/runtime/result-model.md
  - ../../architecture/runtime/error-model.md
