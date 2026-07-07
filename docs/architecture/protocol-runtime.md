# Protocol Runtime

Document
- ID: protocol-runtime
- Version: 1.0
- Status: Stable

Purpose
- Define Protocol Runtime
- Define protocol request lifecycle
- Define protocol boundaries
- Define interaction with Protocol Plugins and SIR Runtime

Definition

Protocol Runtime manages client request processing and transforms protocol-specific requests into validated SIR instances through Protocol Plugins.

Principles

MUST

- Preserve client semantics
- Preserve protocol independence
- Delegate protocol implementation to Protocol Plugins
- Produce valid SIR input
- Preserve session boundaries
- Preserve request ordering

MUST NOT

- Implement protocol specifications
- Perform semantic planning
- Execute operations
- Access downstream engines
- Modify SIR semantics

Runtime Responsibilities

- Accept client connections
- Select Protocol Plugin
- Manage request lifecycle
- Coordinate authentication handoff
- Coordinate session lifecycle
- Coordinate request processing
- Forward validated requests to SIR Runtime
- Return protocol responses

Request Lifecycle

Client Connection

↓

Protocol Identification

↓

Protocol Plugin Selection

↓

Authentication Handoff

↓

Session Resolution

↓

Request Reception

↓

Protocol Parsing

↓

Protocol Validation

↓

SIR Request Generation

↓

SIR Runtime

↓

Protocol Response

↓

Client

Protocol Plugin Interaction

Protocol Runtime

MUST

- Discover available Protocol Plugins
- Select compatible Protocol Plugin
- Invoke Protocol Plugin through published contracts
- Preserve protocol isolation

MUST NOT

- Depend on protocol implementation details
- Access plugin internal state

Session Integration

Protocol Runtime

MUST

- Create session when required
- Resolve existing session
- Associate request with session
- Release session resources

Authentication

Protocol Runtime

MUST

- Delegate authentication to configured provider
- Preserve authenticated identity
- Pass authenticated context to Session Runtime

MUST NOT

- Implement authentication provider
- Persist authentication credentials

Error Handling

Protocol Runtime

MUST

- Normalize protocol errors
- Preserve protocol semantics
- Return deterministic responses

MUST NOT

- Leak runtime implementation details
- Leak downstream engine errors

Runtime Characteristics

MUST

- Support concurrent connections
- Support multiple Protocol Plugins
- Support stateless request processing
- Support graceful connection termination

Ownership

Sovrunn owns

- Protocol Runtime
- Request lifecycle
- Session coordination
- Protocol orchestration

Protocol Plugin owns

- Protocol parsing
- Protocol validation
- Protocol serialization
- Protocol-specific behavior

SIR Runtime owns

- SIR instance creation
- SIR validation
- SIR serialization
- SIR transfer

References

- architecture.md
- runtime.md
- plugin-runtime.md
- sir-runtime.md
- session-runtime.md
