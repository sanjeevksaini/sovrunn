# Neo4j Engine Specification

Document:
  ID: neo4j-engine-specification
  Title: Neo4j Engine Specification
  Parent: engine-specification
  Owner: SDE Specification Layer
  Layer: Specification
  Type: Engine Contract
  Version: 1.0
  Status: Draft

Purpose:
  - Define the Neo4j Downstream Engine integration contract for Sovrunn Data Engine
  - Define Engine Plugin responsibilities for Neo4j
  - Define execution fragment translation, downstream invocation, native result mapping, and native error mapping rules
  - Preserve separation between SDE execution semantics and Neo4j native behavior

Definition:
  Neo4j Engine Specification defines how SDE integrates a Neo4j-compatible graph Downstream Datastore through a Neo4j Engine Plugin.

Scope:
  In Scope:
    - Engine identity
    - Engine Plugin manifest requirements
    - Supported execution forms
    - Capability declarations
    - Execution fragment translation
    - Downstream invocation
    - Native result mapping
    - Native error mapping
    - Credential reference handling
    - Datastore Data Plane boundary

  Out of Scope:
    - Protocol parsing
    - SIR creation
    - Planning internals
    - Data Kernel dependency orchestration
    - Datastore lifecycle management
    - Datastore Operator Plugin behavior
    - Infrastructure Provider behavior
    - Downstream datastore internal implementation

Engine Identity:
  Downstream Engine: Neo4j
  Canonical Identifier: sde.engine.neo4j
  Plugin Type: Engine Plugin

Engine Plugin Responsibilities:
  MUST:
    - Receive execution fragment from Engine Runtime
    - Validate execution fragment compatibility
    - Validate capability boundary
    - Translate fragment into Neo4j-native operation
    - Preserve semantic equivalence
    - Invoke Downstream Datastore through approved Neo4j interface
    - Map native result to Result Model
    - Map native error to Error Model
    - Protect credential references
    - Preserve Request Identifier, Trace Identifier, and Execution Identifier where available

  MUST NOT:
    - Parse client protocol directly
    - Produce Execution Plan
    - Manage datastore lifecycle
    - Replace Datastore Operator Plugin
    - Invoke Infrastructure Provider
    - Modify SDE Control Plane authoritative state
    - Expose raw native result or error directly to Protocol Plugin

Engine Plugin Manifest Requirements:
  Manifest MUST include:
    - Engine identifier
    - Engine version compatibility
    - Plugin identifier
    - Plugin version
    - Supported execution forms
    - Unsupported execution forms
    - Supported capability identifiers
    - Supported capability versions
    - Unsupported capabilities
    - Translation behavior
    - Result mapping behavior
    - Error mapping behavior
    - Credential reference requirements
    - Downstream interface requirements
    - Known semantic gaps
    - Compatibility metadata

Supported Execution Forms:
  - Graph node read
  - Graph edge read
  - Graph traversal
  - Pattern matching query
  - Graph mutation where supported
  - Path query
  - Transaction-bound graph operation where supported

Capability Declarations:
  - graph
  - security
  - indexing
  - transactions where supported
  - streaming

Execution Fragment Rules:
  - Execution fragment MUST be authorized by Execution Plan.
  - Execution fragment MUST be bound to Execution Context.
  - Execution fragment MUST remain within approved capability boundaries.
  - Execution fragment MUST NOT contain raw downstream credentials.
  - Execution fragment MUST NOT contain datastore lifecycle instructions.
  - Execution fragment MUST NOT mutate SDE Control Plane authoritative state.

Translation Rules:
  - Translation MUST preserve semantic equivalence.
  - Unsupported native operation MUST fail deterministically.
  - Translation MUST remain within declared capability support.
  - Engine Plugin MUST NOT silently emulate unsupported capability.
  - Downstream-native operation MUST NOT become SDE platform contract.
  - Graph pattern translation MUST preserve graph semantics.
  - Path and traversal semantics MUST be explicit.
  - Unsupported graph construct MUST fail deterministically.

Downstream Invocation Rules:
  - Engine Plugin MUST invoke Downstream Datastore through approved interface.
  - Engine Plugin MUST preserve tenant boundary.
  - Engine Plugin MUST preserve request and trace correlation where possible.
  - Engine Plugin MUST protect downstream credential references.
  - Engine Plugin MUST NOT bypass approved downstream interface.
  - Engine Plugin MUST NOT manage datastore lifecycle.
  - Graph transaction context MUST be preserved where supported.
  - Graph access MUST preserve tenant and policy constraints.

Datastore Data Plane Boundary:
  - Neo4j Downstream Datastore owns its Datastore Data Plane.
  - SDE Data Plane reaches Neo4j Datastore Data Plane only through approved Engine Plugin.
  - Native execution semantics remain owned by Neo4j.
  - SDE platform semantics MUST not be silently redefined by Neo4j native behavior.
  - Native behavior may be mapped into SDE canonical models only when semantic equivalence is preserved.

Credential Rules:
  - Engine Plugin MUST access credentials only through authorized references.
  - Raw credentials MUST NOT be stored in Execution Plan.
  - Raw credentials MUST NOT be emitted in Result Model.
  - Raw credentials MUST NOT be emitted in Error Model.
  - Credential failure MUST fail closed.
  - Credential-related errors MUST redact unsafe details.

Result Mapping Rules:
  - Native result MUST map to Result Model.
  - Result Model MUST preserve semantic equivalence.
  - Result Model MUST preserve type information where applicable.
  - Result Model MUST preserve schema or metadata where applicable.
  - Result Model MUST preserve affected count where applicable.
  - Raw native result MUST NOT bypass Result Model.
  - Nodes, relationships, paths, and graph result shapes MUST map to Result Model.
  - Path ordering and traversal metadata MUST be preserved where relevant.

Error Mapping Rules:
  - Native error MUST map to Error Model.
  - Error Model MUST preserve safe native error category.
  - Error Model MUST preserve retry classification where known.
  - Unsafe native details MUST be redacted.
  - Unknown downstream outcome MUST be explicit.
  - Raw native error MUST NOT bypass Error Model.
  - Graph query, constraint, and transaction errors MUST map safely to Error Model.

Capability Rules:
  - Engine Plugin MUST declare supported capabilities.
  - Engine capability support MUST be explicit.
  - Planning MUST validate required capabilities before Engine Execution.
  - Engine Runtime MUST reject execution when approved capability binding is unavailable.
  - Capability mismatch MUST produce Error Model.
  - Capability downgrade MUST NOT be silent.
  - Capability emulation MUST be explicit and semantically safe.

Security Rules:
  - Engine Plugin MUST preserve tenant isolation.
  - Engine Plugin MUST preserve security context.
  - Engine Plugin MUST preserve transaction boundary where applicable.
  - Engine Plugin MUST protect credential references.
  - Engine Plugin MUST redact unsafe native details.
  - Engine Plugin MUST not expose raw downstream-native data outside Result Model or Error Model.

Failure Rules:
  - Translation failure MUST produce Error Model.
  - Downstream invocation failure MUST produce Error Model.
  - Native error MUST produce Error Model.
  - Unknown downstream outcome MUST be explicit.
  - Timeout MUST preserve retry classification where known.
  - Failure MUST NOT be converted into success.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.
  - Failure MUST NOT mutate datastore lifecycle state.

Compatibility Rules:
  - Engine version compatibility MUST follow Versioning Specification.
  - Engine Plugin MUST declare supported engine versions.
  - Unsupported engine version MUST fail deterministically.
  - Deprecated engine behavior MUST remain explicit.
  - Compatibility behavior MUST not silently change SDE-visible semantics.

Invariants:
  - Neo4j Engine Plugin is the downstream execution boundary.
  - Engine Runtime invokes Neo4j Engine Plugin.
  - Neo4j Engine Plugin invokes Neo4j Downstream Datastore.
  - Neo4j Engine Plugin does not manage datastore lifecycle.
  - Datastore Operator Plugin is not part of Engine Execution.
  - Infrastructure Provider is not part of Engine Execution.
  - Native result and native error never bypass canonical SDE models.

Relationships:
  Parent:
    - engine.md
  Depends On:
    - engine.md
    - ../versioning/versioning.md
    - ../serialization/serialization.md
    - ../capability/capability.md
    - ../../architecture/data-plane/engine-execution.md
    - ../../architecture/runtime/engine-runtime.md
    - ../../architecture/runtime/execution-context.md
    - ../../architecture/runtime/result-model.md
    - ../../architecture/runtime/error-model.md
    - ../../architecture/control-plane/core-control-plane/engine-registry.md
    - ../../architecture/control-plane/core-control-plane/plugin-registry.md
  Used By:
    - Engine Runtime
    - Engine Plugin
    - Planning
    - Capability Registry
    - Engine Registry
    - Plugin Registry
    - Result Propagation
    - Error Propagation

References:
  - engine.md
  - ../versioning/versioning.md
  - ../serialization/serialization.md
  - ../capability/capability.md
  - ../../architecture/data-plane/engine-execution.md
  - ../../architecture/runtime/engine-runtime.md
  - ../../architecture/runtime/result-model.md
  - ../../architecture/runtime/error-model.md
