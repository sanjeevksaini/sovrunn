# REST Protocol Specification

Document:
  ID: rest-protocol-specification
  Title: REST Protocol Specification
  Parent: protocol-specification
  Owner: SDE Specification Layer
  Layer: Specification
  Type: Protocol Contract
  Version: 1.0
  Status: Draft

Purpose:
  - Define the REST protocol integration contract for Sovrunn Data Engine
  - Define Protocol Plugin responsibilities for REST
  - Define request normalization, session behavior, transaction behavior, result mapping, and error mapping rules
  - Preserve protocol-visible semantics while keeping SDE runtime protocol-neutral

Definition:
  REST Protocol Specification defines how SDE integrates HTTP/REST clients through a REST Protocol Plugin while preserving resource-oriented request, response, status, header, and error semantics.

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
  Protocol: REST
  Canonical Identifier: sde.protocol.rest
  Plugin Type: Protocol Plugin

Protocol Plugin Responsibilities:
  MUST:
    - Parse REST input
    - Validate protocol version compatibility
    - Preserve protocol-visible semantics
    - Produce protocol-normalized intent
    - Preserve session semantics where applicable
    - Preserve transaction semantics where applicable
    - Map Result Model to REST-compatible response
    - Map Error Model to REST-compatible error response
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
  - HTTP GET
  - HTTP POST
  - HTTP PUT
  - HTTP PATCH
  - HTTP DELETE
  - Query parameters
  - Path parameters
  - Request headers
  - Request body
  - Content negotiation

Request Parsing Rules:
  - Protocol Plugin MUST parse protocol input according to REST rules.
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
  - REST requests are generally stateless unless explicit session token or cookie semantics are configured.
  - Authentication context MUST be extracted safely.
  - Tenant context MUST be resolved deterministically.
  - Idempotency keys MUST be represented where supported.
  - Invalid session token MUST fail closed.

Transaction Rules:
  - REST transaction behavior MUST be explicit when exposed.
  - HTTP request boundaries MUST not imply datastore transaction boundaries unless contract says so.
  - Idempotency behavior MUST be declared.
  - Unsupported transactional request semantics MUST fail deterministically.

Capability Declarations:
  - security
  - object
  - streaming
  - search
  - federation

Result Mapping Rules:
  - Result Model MUST map to HTTP status, headers, and response body.
  - Content type MUST be explicit.
  - Pagination and continuation MUST be represented safely.
  - Partial response semantics MUST follow HTTP contract.
  - Raw downstream-native result MUST not be exposed.

Error Mapping Rules:
  - Error Model MUST map to HTTP status, safe error body, and safe headers.
  - Unsafe details MUST be redacted.
  - Trace Identifier MAY be exposed in header or body only when policy allows it.
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
  - REST Protocol Plugin owns REST-specific parsing and mapping.
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
