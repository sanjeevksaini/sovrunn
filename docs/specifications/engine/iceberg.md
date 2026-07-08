# Iceberg Engine Specification

Document:
  ID: iceberg-engine-specification
  Title: Iceberg Engine Specification
  Parent: engine-specification
  Owner: SDE Specification Layer
  Layer: Specification
  Type: Engine Contract
  Version: 1.0
  Status: Draft

Purpose:
  - Define the Iceberg Downstream Engine integration contract for Sovrunn Data Engine
  - Define Engine Plugin responsibilities for Iceberg
  - Define execution fragment translation, downstream invocation, native result mapping, and native error mapping rules
  - Preserve separation between SDE execution semantics and Iceberg native behavior

Definition:
  Iceberg Engine Specification defines how SDE integrates Apache Iceberg table metadata and data operations through an Iceberg Engine Plugin.

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
  Downstream Engine: Iceberg
  Canonical Identifier: sde.engine.iceberg
  Plugin Type: Engine Plugin

Engine Plugin Responsibilities:
  MUST:
    - Receive execution fragment from Engine Runtime
    - Validate execution fragment compatibility
    - Validate capability boundary
    - Translate fragment into Iceberg-native operation
    - Preserve semantic equivalence
    - Invoke Downstream Datastore through approved Iceberg interface
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
  - Table metadata read
  - Snapshot read
  - Manifest read
  - Partition planning
  - Scan planning
  - Append operation where supported
  - Delete operation where supported
  - Time travel read where supported

Capability Declarations:
  - object
  - indexing
  - security
  - streaming
  - federation

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
  - Table, snapshot, manifest, and partition semantics MUST be explicit.
  - Time travel semantics MUST be preserved where supported.
  - Unsupported table operation MUST fail deterministically.

Downstream Invocation Rules:
  - Engine Plugin MUST invoke Downstream Datastore through approved interface.
  - Engine Plugin MUST preserve tenant boundary.
  - Engine Plugin MUST preserve request and trace correlation where possible.
  - Engine Plugin MUST protect downstream credential references.
  - Engine Plugin MUST NOT bypass approved downstream interface.
  - Engine Plugin MUST NOT manage datastore lifecycle.
  - Catalog and object storage references MUST preserve tenant and policy constraints.
  - Snapshot selection MUST be deterministic.

Datastore Data Plane Boundary:
  - Iceberg Downstream Datastore owns its Datastore Data Plane.
  - SDE Data Plane reaches Iceberg Datastore Data Plane only through approved Engine Plugin.
  - Native execution semantics remain owned by Iceberg.
  - SDE platform semantics MUST not be silently redefined by Iceberg native behavior.
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
  - Table metadata, snapshot metadata, manifests, scan tasks, and rows where applicable MUST map to Result Model.
  - Continuation and split metadata MUST be safe.

Error Mapping Rules:
  - Native error MUST map to Error Model.
  - Error Model MUST preserve safe native error category.
  - Error Model MUST preserve retry classification where known.
  - Unsafe native details MUST be redacted.
  - Unknown downstream outcome MUST be explicit.
  - Raw native error MUST NOT bypass Error Model.
  - Catalog, snapshot, manifest, object storage, and consistency errors MUST map safely to Error Model.

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
  - Iceberg Engine Plugin is the downstream execution boundary.
  - Engine Runtime invokes Iceberg Engine Plugin.
  - Iceberg Engine Plugin invokes Iceberg Downstream Datastore.
  - Iceberg Engine Plugin does not manage datastore lifecycle.
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
