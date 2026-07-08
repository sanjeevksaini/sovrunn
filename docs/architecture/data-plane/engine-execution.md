# Engine Execution

Document:
  ID: engine-execution
  Title: Engine Execution
  Parent: sde-data-plane
  Owner: SDE Data Plane
  Layer: SDE Data Plane
  Type: Flow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Engine Runtime and Engine Plugin execution behavior inside SDE Data Plane
  - Define the downstream execution delegation boundary
  - Define how execution fragments become downstream-native operations
  - Define result and error mapping at the Engine Plugin boundary
  - Point detailed engine-stage behavior to focused subflow files

Definition:
  Engine Execution is the SDE Data Plane flow in which Engine Runtime resolves an approved Engine Plugin, validates engine and plugin bindings, delegates execution fragments, and receives canonical Result Model or Error Model output.

  Engine Execution is the only SDE Data Plane path that may invoke a Downstream Datastore.

  Engine Execution must access a Downstream Datastore only through an approved Engine Plugin and approved downstream interface.

Scope:
  In Scope:
    - Engine Runtime request intake
    - Engine metadata resolution
    - Engine Plugin resolution
    - Engine Plugin availability validation
    - Execution fragment translation
    - Downstream datastore invocation
    - Datastore Data Plane boundary
    - Native result mapping
    - Native error mapping
    - Engine execution completion

  Out of Scope:
    - Protocol parsing
    - SIR creation
    - Planning
    - Data Kernel dependency orchestration
    - Downstream datastore lifecycle management
    - Datastore Operator Plugin execution
    - Infrastructure Provider execution
    - SDE Control Plane authoritative state mutation
    - Datastore Data Plane internals

High-Level Flow:
  - Data Kernel sends operation execution request to Engine Runtime.
  - Engine Runtime resolves approved engine metadata.
  - Engine Runtime resolves approved Engine Plugin metadata.
  - Engine Runtime validates engine-plugin compatibility.
  - Engine Runtime creates Engine Plugin invocation context.
  - Engine Plugin translates execution fragment into downstream-native operation.
  - Engine Plugin invokes Downstream Datastore through approved interface.
  - Downstream Datastore executes through Datastore Data Plane.
  - Engine Plugin maps native result to Result Model or native error to Error Model.
  - Engine Runtime returns canonical output to Data Kernel.

Flow Diagram:
  Data Kernel
    ↓
  Engine Runtime
    ↓
  Engine Metadata
    ↓
  Engine Plugin
    ↓
  Downstream-Native Operation
    ↓
  Downstream Datastore
    ↓
  Datastore Data Plane
    ↓
  Native Result or Native Error
    ↓
  Engine Plugin
    ↓
  Result Model or Error Model
    ↓
  Engine Runtime
    ↓
  Data Kernel

Stage Map:
  Engine Request Intake:
    Document: engine-execution/engine-request-intake.md
    Owner: Engine Runtime

  Engine Metadata Resolution:
    Document: engine-execution/engine-metadata-resolution.md
    Owner: Engine Runtime

  Engine Plugin Resolution:
    Document: engine-execution/engine-plugin-resolution.md
    Owner: Engine Runtime

  Execution Fragment Translation:
    Document: engine-execution/execution-fragment-translation.md
    Owner: Engine Plugin

  Downstream Invocation:
    Document: engine-execution/downstream-invocation.md
    Owner: Engine Plugin

  Native Output Mapping:
    Document: engine-execution/native-output-mapping.md
    Owner: Engine Plugin

  Engine Completion:
    Document: engine-execution/engine-completion.md
    Owner: Engine Runtime

Engine Inputs:
  - Operation execution request
  - Execution fragment
  - Execution Context
  - Engine metadata
  - Engine Plugin metadata
  - Capability requirements
  - Runtime configuration
  - Downstream endpoint reference
  - Downstream credential reference where authorized

Engine Outputs:
  Success:
    - Result Model
    - Engine execution metadata where safe

  Failure:
    - Error Model
    - Retry classification
    - Safe causal metadata
    - Unknown outcome marker where required

Engine Runtime Rules:
  - Engine Runtime MUST receive execution requests from Data Kernel.
  - Engine Runtime MUST resolve approved engine metadata.
  - Engine Runtime MUST resolve approved Engine Plugin metadata.
  - Engine Runtime MUST validate engine-plugin compatibility.
  - Engine Runtime MUST validate plugin availability.
  - Engine Runtime MUST delegate downstream work to Engine Plugin.
  - Engine Runtime MUST NOT access Downstream Datastore directly.
  - Engine Runtime MUST NOT invoke Datastore Operator Plugin.
  - Engine Runtime MUST NOT invoke Infrastructure Provider.
  - Engine Runtime MUST NOT manage datastore lifecycle.
  - Engine Runtime MUST NOT modify SDE Control Plane authoritative state.

Engine Plugin Rules:
  - Engine Plugin MUST translate execution fragment into downstream-native operation.
  - Engine Plugin MUST preserve semantic equivalence.
  - Engine Plugin MUST respect declared capability boundaries.
  - Engine Plugin MUST respect Execution Context.
  - Engine Plugin MUST invoke Downstream Datastore through approved interface.
  - Engine Plugin MUST map native result to Result Model.
  - Engine Plugin MUST map native error to Error Model.
  - Engine Plugin MUST protect downstream credentials.
  - Engine Plugin MUST NOT manage datastore lifecycle.
  - Engine Plugin MUST NOT replace Datastore Operator Plugin.
  - Engine Plugin MUST NOT expose raw native result or error directly to Protocol Plugin.

Datastore Boundary Rules:
  - Downstream Datastore owns Datastore Data Plane.
  - SDE Data Plane reaches Datastore Data Plane only through Engine Plugin.
  - Downstream native execution semantics remain owned by the Downstream Datastore.
  - SDE platform semantics must not be silently redefined by downstream-native behavior.
  - Engine Plugin may map native behavior into SDE canonical models only when semantic equivalence is preserved.

Capability Rules:
  - Engine Execution MUST remain within capabilities approved during Planning.
  - Engine Plugin MUST execute within declared Capability Manifest boundaries.
  - Engine Runtime MUST reject execution when required capability binding is unavailable.
  - Engine Plugin MUST NOT silently emulate unsupported capability.
  - Capability mismatch MUST produce Error Model.

Credential Rules:
  - Downstream credentials MUST be accessed only through authorized secret references.
  - Raw credentials MUST NOT be stored in Execution Plan.
  - Raw credentials MUST NOT be exposed in Result Model or Error Model.
  - Credential failure MUST fail closed.
  - Credential-related errors MUST redact unsafe details.

Result Rules:
  - Native result MUST be mapped to Result Model at Engine Plugin boundary.
  - Result Model MUST preserve type information.
  - Result Model MUST preserve schema metadata where applicable.
  - Result Model MUST preserve affected count where applicable.
  - Result Model MUST preserve stream, cursor, or continuation references where applicable.
  - Raw native result MUST NOT bypass Result Model.

Error Rules:
  - Native error MUST be mapped to Error Model at Engine Plugin boundary.
  - Error Model MUST preserve safe Code, Category, Severity, Source, State, Trace Identifier, Timestamp, and retry classification.
  - Unsafe native error details MUST be redacted.
  - Unknown downstream outcome MUST be reported explicitly.
  - Failure MUST NOT be converted into success.
  - Native error MUST NOT bypass Error Model.

Failure Rules:
  - Missing engine metadata MUST produce Error Model.
  - Missing Engine Plugin MUST produce Error Model.
  - Incompatible Engine Plugin MUST fail closed.
  - Plugin unavailable failure MUST preserve retry classification.
  - Translation failure MUST produce Error Model.
  - Downstream timeout MUST produce Error Model with retry classification.
  - Unknown downstream outcome MUST be explicit.
  - Engine Execution failure MUST NOT mutate SDE Control Plane authoritative state.
  - Engine Execution failure MUST NOT mutate datastore lifecycle state.

Concurrency Rules:
  - Engine Runtime MUST preserve execution isolation.
  - Engine Plugin invocation MUST preserve tenant isolation.
  - Concurrent invocations MUST not share unsafe mutable plugin state.
  - Downstream connection pooling MUST preserve tenant and credential boundaries.
  - Streaming output MUST preserve request and cursor boundaries.
  - Transaction-associated operations MUST preserve transaction boundaries.

Invariants:
  - Engine Execution starts with Data Kernel delegation.
  - Engine Execution ends with Result Model or Error Model.
  - Engine Runtime is the runtime delegation boundary.
  - Engine Plugin is the downstream execution boundary.
  - Datastore Data Plane remains owned by the Downstream Datastore.
  - Engine Runtime does not access Downstream Datastore directly.
  - Engine Plugin does not manage datastore lifecycle.
  - Datastore Operator Plugin is not part of Engine Execution.
  - Infrastructure Provider is not part of Engine Execution.

Relationships:
  Parent:
    - data-plane.md
  Children:
    - engine-execution/engine-request-intake.md
    - engine-execution/engine-metadata-resolution.md
    - engine-execution/engine-plugin-resolution.md
    - engine-execution/execution-fragment-translation.md
    - engine-execution/downstream-invocation.md
    - engine-execution/native-output-mapping.md
    - engine-execution/engine-completion.md
  Depends On:
    - data-plane-map.md
    - request-flow.md
    - kernel-execution.md
    - ../runtime/engine-runtime.md
    - ../runtime/plugin-runtime.md
    - ../runtime/execution-context.md
    - ../runtime/result-model.md
    - ../runtime/error-model.md
    - ../control-plane/core-control-plane/engine-registry.md
    - ../control-plane/core-control-plane/plugin-registry.md
    - ../control-plane/core-control-plane/capability-governance.md
  Used By:
    - kernel-execution.md
    - result-propagation.md
    - error-propagation.md
    - Engine Plugin specifications

References:
  - data-plane.md
  - data-plane-map.md
  - request-flow.md
  - kernel-execution.md
  - engine-execution/engine-request-intake.md
  - engine-execution/engine-metadata-resolution.md
  - engine-execution/engine-plugin-resolution.md
  - engine-execution/execution-fragment-translation.md
  - engine-execution/downstream-invocation.md
  - engine-execution/native-output-mapping.md
  - engine-execution/engine-completion.md
  - ../runtime/engine-runtime.md
  - ../runtime/plugin-runtime.md
  - ../runtime/execution-context.md
  - ../runtime/result-model.md
  - ../runtime/error-model.md
  - ../control-plane/core-control-plane/engine-registry.md
  - ../control-plane/core-control-plane/plugin-registry.md
  - ../control-plane/core-control-plane/capability-governance.md
