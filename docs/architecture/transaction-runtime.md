# Transaction Runtime

Document
- ID: transaction-runtime
- Version: 1.0
- Status: Stable

Purpose
- Define Transaction Runtime
- Define transaction context
- Define transaction lifecycle
- Define transaction boundaries

Definition

Transaction Runtime manages Sovrunn transaction context across runtime execution.

A Sovrunn transaction is independent of downstream engine transaction implementation.

Principles

MUST

- Preserve transaction intent
- Preserve semantic consistency
- Preserve engine independence
- Support capability driven transaction behavior
- Fail deterministically when transaction semantics cannot be preserved

MUST NOT

- Assume downstream transaction capability
- Hide semantic degradation
- Implement downstream engine transaction managers
- Modify SIR semantics

Responsibilities

- Create transaction context
- Resolve transaction context
- Propagate transaction context
- Coordinate transaction lifecycle
- Validate transaction capability
- Complete transaction context

Transaction Model

Transaction

MUST contain

- Transaction Identifier
- Session Identifier
- Transaction State
- Isolation Requirement
- Consistency Requirement
- Participating Resources
- Required Capabilities

Transaction State

- Created
- Active
- Preparing
- Committing
- Committed
- Rolling Back
- Rolled Back
- Failed

Lifecycle

Create

↓

Bind Session

↓

Validate Capabilities

↓

Active

↓

Prepare

↓

Commit

↓

Complete

Rollback Path

Active

↓

Rolling Back

↓

Rolled Back

Capability Interaction

Transaction Runtime

MUST

- Use Capability Registry
- Validate transaction capabilities
- Detect unsupported transaction semantics
- Report deterministic transaction failures

MUST NOT

- Assume ACID support
- Assume distributed transaction support
- Emulate unsupported semantics silently

Engine Interaction

Transaction Runtime

MUST

- Coordinate through Data Kernel and Engine Runtime
- Preserve transaction boundaries
- Preserve semantic consistency

MUST NOT

- Access downstream engines directly
- Manage native engine transactions directly

Engine Plugin

MAY

- Create native transaction
- Commit native transaction
- Rollback native transaction
- Expose transaction capability

Consistency

Transaction Runtime

MUST

- Preserve declared consistency requirement
- Fail when consistency cannot be preserved

MUST NOT

- Downgrade consistency without explicit policy

Failure Handling

Transaction Runtime

MUST

- Detect transaction failure
- Preserve deterministic state transition
- Preserve session isolation
- Report transaction outcome

MUST NOT

- Leave unknown transaction state without reporting
- Corrupt runtime state
- Corrupt semantic intent

Ownership

Sovrunn owns

- Transaction Runtime
- Transaction context
- Transaction lifecycle
- Transaction coordination

Engine Plugin owns

- Native transaction integration

Downstream Engine owns

- Native transaction implementation
- Native isolation
- Native durability

References

- architecture.md
- runtime.md
- session-runtime.md
- data-kernel.md
- engine-runtime.md
- capability-registry.md
- specifications/sir/capabilities.md
- specifications/sir/capability-model.md
