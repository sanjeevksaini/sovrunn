# Native Protocol Specification

Document:
  ID: native-protocol-specification
  Title: Native Protocol Specification
  Parent: protocol-specification
  Owner: SDE Specification Layer
  Layer: Specification
  Type: Protocol Contract
  Version: 1.0
  Status: Draft

Purpose:
  - Define the Native protocol integration contract for Sovrunn Data Engine
  - Define Protocol Plugin responsibilities for Native
  - Define request normalization, session behavior, transaction behavior, result mapping, and error mapping rules
  - Preserve protocol-visible semantics while keeping SDE runtime protocol-neutral

Definition:
  Native Protocol Specification defines the SDE-native protocol contract used by trusted SDE clients, internal tools, SDKs, or agents that communicate with SDE using canonical request and response models.

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
  Protocol: Native
  Canonical Identifier: sde.protocol.native
  Plugin Type: Protocol Plugin

Protocol Plugin Responsibilities:
  MUST:
    - Parse Native input
    - Validate protocol version compatibility
    - Preserve protocol-visible semantics
    - Produce protocol-normalized intent
    - Preserve session semantics where applicable
    - Preserve transaction semantics where applicable
    - Map Result Model to Native-compatible response
    - Map Error Model to Native-compatible error response
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
  - Canonical request
  - Canonical query request
  - Canonical mutation request
  - Capability-aware request
  - Execution preference request
  - Streaming request
  - Administrative read-only metadata request where allowed

Request Parsing Rules:
  - Protocol Plugin MUST parse protocol input according to Native rules.
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
  - Native protocol MAY be stateless or session-aware depending on client contract.
  - Session references MUST be explicit.
  - Tenant context MUST be explicit or resolved through approved security context.
  - Request context MUST preserve traceability.
  - Invalid session reference MUST fail closed.

Transaction Rules:
  - Transaction intent MUST be explicit.
  - Transaction references MUST be explicit.
  - Native protocol MUST not silently create transaction scope.
  - Unsupported transaction semantics MUST fail deterministically.
  - Unknown transaction outcome MUST be explicit.

Capability Declarations:
  - transactions
  - security
  - object
  - cache
  - search
  - indexing
  - streaming
  - federation
  - vector
  - graph

Result Mapping Rules:
  - Result Model MAY be returned directly only when client is authorized for canonical models.
  - Protocol response MUST preserve canonical type and metadata behavior.
  - Streaming and continuation behavior MUST be explicit.
  - Unsafe internal metadata MUST be redacted unless internal policy allows it.

Error Mapping Rules:
  - Error Model MAY be returned directly only when client is authorized for canonical models.
  - Unsafe details MUST be redacted.
  - Trace and execution identifiers MAY be exposed according to policy.
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
  - Native Protocol Plugin owns Native-specific parsing and mapping.
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
