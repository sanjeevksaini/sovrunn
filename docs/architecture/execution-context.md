# Execution Context

Document

- ID: execution-context
- Version: 1.0
- Status: Stable

Purpose

- Define Execution Context
- Define execution-scoped runtime state
- Define execution context boundaries
- Define context propagation

Definition

Execution Context is the immutable execution-scoped context that accompanies an Execution Plan throughout the Sovrunn Runtime.

Execution Context provides runtime information required to execute semantic intent without modifying the Execution Plan.

Principles

MUST

- Be immutable
- Be execution scoped
- Preserve semantic intent
- Preserve execution isolation
- Be protocol independent
- Be engine independent

MUST NOT

- Own session lifecycle
- Own transaction lifecycle
- Own business data
- Modify Execution Plan
- Modify SIR

Responsibilities

Execution Context

- Associate runtime state with an execution
- Propagate runtime metadata
- Reference runtime contexts
- Preserve execution identity

Execution Context Model

Execution Context

MUST contain

- Execution Identifier
- Request Identifier
- Trace Identifier
- Session Reference
- Transaction Reference
- Security Context
- Tenant Context
- Runtime Context

MAY contain

- Deadline
- Timeout
- Execution Options
- Locale
- Correlation Identifier
- Extension Context

Referenced Contexts

Session Context

Owns

- Identity
- Session Variables
- Preferences

Transaction Context

Owns

- Transaction State
- Isolation
- Consistency

Security Context

Contains

- Principal
- Roles
- Permissions

Tenant Context

Contains

- Tenant Identifier
- Sovereignty Policy
- Data Residency Policy

Runtime Context

Contains

- Runtime Configuration
- Execution Limits
- Feature Flags

Lifecycle

Create

↓

Initialize

↓

Propagate

↓

Consume

↓

Dispose

Propagation

Execution Context

MUST

- Propagate unchanged across runtime components
- Preserve execution identity
- Preserve traceability

MUST NOT

- Be modified by runtime components
- Leak across independent executions

Runtime Interaction

Execution Context

Flows through

- Protocol Runtime
- SIR Runtime
- Planning
- Data Kernel
- Engine Runtime

Engine Plugin

MAY

- Read execution context
- Use execution options

MUST NOT

- Modify execution context

Characteristics

MUST

- Be immutable
- Be serializable
- Be thread safe
- Be execution scoped

Ownership

Execution Context owns

- Execution-scoped runtime state

Session Runtime owns

- Session Context

Transaction Runtime owns

- Transaction Context

References

- runtime.md
- session-runtime.md
- transaction-runtime.md
- execution-plan.md
- protocol-runtime.md
- data-kernel.md
