# Engine Specification

Document:
  ID: engine-specification
  Title: Engine Specification
  Parent: specifications
  Owner: SDE Specification Layer
  Layer: Specification
  Type: Contract
  Version: 1.0
  Status: Draft

Purpose:
  - Define how Downstream Engines integrate with SDE Data Plane
  - Define Engine Plugin contract boundaries
  - Define execution fragment translation and native output mapping
  - Preserve separation between SDE execution semantics and downstream-native datastore behavior

Definition:
  Engine Specification defines the contract by which a Downstream Engine is integrated into SDE through an Engine Plugin.

  An Engine Plugin translates SDE execution fragments into downstream-native operations, invokes the Downstream Datastore through approved interfaces, and maps native result or native error into canonical SDE models.

Scope:
  In Scope:
    - Engine identity
    - Engine Plugin manifest
    - Engine capability declaration
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
    - Downstream datastore internal architecture

Engine Plugin Responsibilities:
  MUST:
    - Receive execution fragment from Engine Runtime
    - Validate execution fragment compatibility
    - Validate capability boundary
    - Translate fragment into downstream-native operation
    - Preserve semantic equivalence
    - Invoke Downstream Datastore through approved interface
    - Map native result to Result Model
    - Map native error to Error Model
    - Protect credential references
    - Preserve Trace Identifier and Execution Identifier where available

  MUST NOT:
    - Parse client protocol directly
    - Produce Execution Plan
    - Manage datastore lifecycle
    - Replace Datastore Operator Plugin
    - Invoke Infrastructure Provider
    - Modify SDE Control Plane authoritative state
    - Expose raw native result or error directly to Protocol Plugin

Engine Plugin Manifest:
  Manifest MUST include:
    - Engine identifier
    - Engine version compatibility
    - Plugin identifier
    - Plugin version
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

Execution Fragment:
  Execution fragment is an Engine Runtime request payload derived from Execution Plan.

  It is:
    - Plan-authorized
    - Engine-targeted
    - Request-scoped
    - Semantics-preserving
    - Bound to Execution Context

  It is not:
    - SIR
    - Full Execution Plan
    - Protocol request
    - Datastore lifecycle instruction
    - SDE Control Plane mutation request

Translation Rules:
  - Engine Plugin MUST translate only plan-authorized fragments.
  - Translation MUST preserve semantic equivalence.
  - Unsupported native operation MUST fail deterministically.
  - Translation MUST remain within declared capabilities.
  - Engine Plugin MUST NOT silently emulate unsupported capability.
  - Downstream-native operation MUST NOT become SDE platform contract.

Downstream Invocation Rules:
  - Engine Plugin MUST invoke Downstream Datastore through approved interface.
  - Engine Plugin MUST protect downstream credentials.
  - Engine Plugin MUST preserve tenant boundary.
  - Engine Plugin MUST preserve transaction boundary where applicable.
  - Engine Plugin MUST preserve request and trace correlation where possible.
  - Engine Plugin MUST NOT bypass approved downstream interface.

Datastore Data Plane Boundary:
  - Downstream Datastore owns Datastore Data Plane.
  - SDE reaches Datastore Data Plane only through Engine Plugin.
  - Downstream native execution semantics remain owned by Downstream Datastore.
  - SDE platform semantics must not be silently redefined by downstream-native behavior.
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
  - Result Model MUST preserve schema metadata where applicable.
  - Result Model MUST preserve affected count where applicable.
  - Result Model MUST preserve cursor, stream, or continuation references where applicable.
  - Raw native result MUST NOT bypass Result Model.

Error Mapping Rules:
  - Native error MUST map to Error Model.
  - Error Model MUST preserve safe native error category.
  - Error Model MUST preserve retry classification where known.
  - Unsafe native details MUST be redacted.
  - Unknown downstream outcome MUST be explicit.
  - Failure MUST NOT be converted into success.
  - Raw native error MUST NOT bypass Error Model.

Capability Rules:
  - Engine Plugin MUST declare supported capabilities.
  - Engine capability support MUST be explicit.
  - Planning MUST validate required capabilities before Engine Execution.
  - Engine Runtime MUST reject execution when approved capability binding is unavailable.
  - Capability mismatch MUST produce Error Model.

Invariants:
  - Engine Plugin is the downstream execution boundary.
  - Engine Runtime invokes Engine Plugin.
  - Engine Plugin invokes Downstream Datastore.
  - Engine Plugin does not manage datastore lifecycle.
  - Datastore Operator Plugin is not part of Engine Execution.
  - Infrastructure Provider is not part of Engine Execution.
  - Native result and native error never bypass canonical SDE models.

Relationships:
  Parent:
    - ../specifications
  Children:
    - postgresql.md
    - mysql.md
    - mongodb.md
    - redis.md
    - cassandra.md
    - opensearch.md
    - neo4j.md
    - milvus.md
    - s3.md
    - iceberg.md
    - delta-lake.md
    - parquet.md
  Depends On:
    - ../versioning/versioning.md
    - ../serialization/serialization.md
    - ../capability/capability.md
    - ../../architecture/data-plane/engine-execution.md
    - ../../architecture/runtime/engine-runtime.md
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
  - ../versioning/versioning.md
  - ../serialization/serialization.md
  - ../capability/capability.md
  - ../../architecture/data-plane/engine-execution.md
  - ../../architecture/runtime/engine-runtime.md
  - ../../architecture/runtime/result-model.md
  - ../../architecture/runtime/error-model.md
  - ../../architecture/control-plane/core-control-plane/engine-registry.md
  - ../../architecture/control-plane/core-control-plane/plugin-registry.md
