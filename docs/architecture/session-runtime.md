# Session Runtime

Document
- ID: session-runtime
- Version: 1.0
- Status: Stable

Purpose
- Define Session Runtime
- Define session lifecycle
- Define session responsibilities
- Define session boundaries

Definition

Session Runtime manages the lifecycle and execution context of Sovrunn sessions.

A Sovrunn Session is independent of downstream engine sessions.

Principles

MUST

- Preserve session identity
- Preserve execution context
- Preserve protocol independence
- Preserve engine independence
- Support concurrent sessions
- Support stateless runtime deployment

MUST NOT

- Own downstream engine sessions
- Own business data
- Modify SIR semantics
- Manage transaction execution

Responsibilities

- Create session
- Resolve session
- Maintain session context
- Propagate session context
- Expire session
- Destroy session

Session Model

Session

MUST contain

- Session Identifier
- Identity Context
- Security Context
- Runtime Context
- Capability Context
- Configuration Context

MAY contain

- User Context
- Locale Context
- Preference Context
- Extension Context

Session Lifecycle

Create

↓

Authenticate

↓

Initialize Context

↓

Active

↓

Idle

↓

Resume

↓

Terminate

Session Context

Identity Context

- Principal
- Tenant
- Roles

Security Context

- Authentication State
- Authorization Context

Runtime Context

- Runtime Configuration
- Session Variables

Capability Context

- Available Runtime Capabilities

Configuration Context

- Runtime Configuration
- Execution Preferences

Propagation

Session Runtime

MUST

- Propagate session context across runtime components
- Preserve session identity
- Preserve session isolation

MUST NOT

- Expose internal runtime state
- Leak session context across sessions

Engine Interaction

Session Runtime

MUST

- Remain independent of downstream engine sessions
- Provide session context to Engine Runtime when required

Engine Plugin

MAY

- Create engine-native session
- Reuse engine-native session
- Destroy engine-native session

Runtime Characteristics

MUST

- Be stateless
- Support concurrent sessions
- Support horizontal scaling
- Support session recovery

Failure Handling

Session Runtime

MUST

- Detect invalid sessions
- Expire inactive sessions
- Preserve session isolation
- Recover session context when possible

MUST NOT

- Corrupt session state
- Corrupt runtime state

Ownership

Sovrunn owns

- Session Runtime
- Session lifecycle
- Session context
- Session propagation

Engine Plugin owns

- Engine session management
- Engine connection management

Downstream Engine owns

- Native session implementation

References

- architecture.md
- runtime.md
- protocol-runtime.md
- sir-runtime.md
- transaction-runtime.md
