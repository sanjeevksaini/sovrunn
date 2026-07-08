# Plugin Runtime

Document
- ID: plugin-runtime
- Version: 1.0
- Status: Stable

Purpose
- Define plugin runtime architecture
- Define plugin categories
- Define plugin lifecycle
- Define plugin boundaries

Definition

Plugin Runtime manages runtime extensions that integrate external systems with Sovrunn without modifying Core runtime components

Plugin Categories

Protocol Plugin
- Integrates external client protocol
- Consumed by Protocol Runtime
- Produces protocol specific input for SIR Runtime

Engine Plugin
- Integrates downstream Engine
- Consumed by Engine Runtime
- Executes engine specific operations
- Publishes Capability Manifest

Rules

MUST
- Support Protocol Plugins
- Support Engine Plugins
- Discover plugins
- Load plugins
- Validate plugins
- Isolate plugin failures
- Preserve runtime boundaries

MUST NOT
- Own protocol specifications
- Own downstream engines
- Modify SIR semantics
- Bypass runtime contracts
- Expose plugin internals as platform API

Lifecycle

Discover

↓

Load

↓

Validate

↓

Register

↓

Activate

↓

Observe

↓

Deactivate

↓

Unload

Protocol Plugin Responsibilities

MUST
- Declare supported Protocol
- Preserve client semantics
- Produce valid protocol output for Protocol Runtime

MUST NOT
- Access Engine
- Produce Execution Plan
- Modify SIR semantics

Engine Plugin Responsibilities

MUST
- Declare downstream Engine
- Publish Capability Manifest
- Translate execution fragments to native engine operations
- Preserve semantic equivalence

MUST NOT
- Own downstream Engine
- Modify Execution Plan semantics
- Invent capabilities

Capability Manifest

MUST
- Use canonical capability identifiers
- Declare supported capabilities
- Declare capability versions
- Declare engine version
- Declare plugin version

MUST NOT
- Use display names as capability identifiers
- Declare unsupported capabilities

Failure Handling

Plugin Runtime

MUST
- Detect plugin failure
- Isolate failed plugin
- Mark plugin unavailable
- Preserve runtime availability

MUST NOT
- Corrupt runtime state
- Corrupt SIR
- Corrupt Execution Plan

Ownership

Sovrunn owns
- Plugin Runtime
- Plugin contracts
- Plugin lifecycle
- Plugin validation

Protocol Plugin owns
- Protocol integration

Engine Plugin owns
- Engine integration
- Capability Manifest

Downstream Engine owns
- Storage
- Native execution
- Native optimization
- Native transaction management

References
- architecture.md
- runtime.md
- protocol-runtime.md
- engine-runtime.md
- capability-registry.md
- specifications/sir/capabilities.md
- specifications/sir/capability-model.md
