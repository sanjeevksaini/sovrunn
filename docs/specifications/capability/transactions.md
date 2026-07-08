# Transactions Capability Specification

Document:
  ID: transactions-specification
  Title: Transactions Capability Specification
  Parent: capability-specification
  Owner: SDE Specification Layer
  Layer: Specification
  Type: Capability Contract
  Version: 1.0
  Status: Draft

Purpose:
  - Define the Transactions capability category for Sovrunn Data Engine
  - Define capability requirements, support levels, validation rules, and failure behavior
  - Provide Planning, Protocol Plugin, Engine Plugin, and Capability Registry with a common contract

Definition:
  Transactions capability defines how SDE represents transaction boundaries, transaction isolation expectations, atomicity requirements, consistency guarantees, and transaction outcome handling across protocols, runtime, Engine Plugins, and Downstream Datastores.

Scope:
  In Scope:
    - Capability identity
    - Capability requirements
    - Capability support levels
    - Capability declaration
    - Planning validation
    - Engine Plugin validation
    - Protocol Plugin interaction
    - Error Model behavior

  Out of Scope:
    - Downstream datastore native implementation
    - Engine Plugin implementation code
    - Protocol Plugin implementation code
    - Datastore lifecycle management
    - Infrastructure provisioning

Capability Identity:
  Category: Transactions
  Canonical Identifier Prefix: sde.capability.transactions

Capability Requirements:
  - Transaction boundary support
  - Explicit transaction begin, commit, and rollback behavior
  - Autocommit behavior where protocol requires it
  - Isolation level support
  - Atomicity scope
  - Read/write consistency expectations
  - Savepoint support where applicable
  - Distributed transaction support where applicable
  - Unknown outcome reporting

Support Levels:
  Required:
    Meaning: Request execution MUST fail if the capability is unavailable.

  Optional:
    Meaning: Request execution MAY continue without the capability when semantics remain valid.

  Preferred:
    Meaning: Planning SHOULD use the capability when available.

  Unsupported:
    Meaning: Selected Protocol Plugin, Engine Plugin, or Downstream Datastore does not support the capability.

  Emulated:
    Meaning: Capability may be provided by SDE or plugin behavior only when semantic equivalence is guaranteed.

Declaration Rules:
  - Capability support MUST be declared explicitly.
  - Capability version MUST be declared.
  - Support level MUST be declared.
  - Unsupported combinations MUST be declared.
  - Known semantic gaps MUST be declared.
  - Safe downgrade behavior MUST be declared when allowed.
  - Capability declaration MUST be approved before Data Plane use.

Planning Rules:
  - Planning MUST validate required Transactions capabilities before Execution Plan emission.
  - Planning MUST use approved Capability Registry metadata.
  - Planning MUST use approved Protocol Plugin and Engine Plugin capability metadata.
  - Planning MUST reject unsupported required capability.
  - Planning MUST NOT silently downgrade required capability.
  - Planning MUST NOT invent capability support.
  - Planning MUST encode capability decisions in Execution Plan where execution behavior depends on them.

Execution Rules:
  - Execution MUST preserve explicit transaction boundaries.
  - Transaction Runtime owns transaction context.
  - Engine Plugin MUST preserve downstream transaction semantics when delegated.
  - Unsupported isolation level MUST fail or downgrade only when explicitly allowed.
  - Unknown transaction outcome MUST be represented explicitly.
  - Autocommit behavior MUST preserve protocol-visible semantics.
  - SDE MUST NOT silently emulate unsupported transaction semantics.

Failure Rules:
  - Unsupported required transaction capability MUST produce Error Model.
  - Invalid transaction reference MUST fail closed.
  - Commit uncertainty MUST produce Error Model with unknown outcome marker.
  - Rollback failure MUST produce Error Model.
  - Isolation mismatch MUST produce Error Model unless explicit safe downgrade exists.

Security Rules:
  - Capability validation MUST preserve tenant isolation.
  - Capability metadata MUST NOT expose raw secrets.
  - Capability failure details MUST redact unsafe internal or downstream-native details.
  - Capability validation MUST fail closed when security impact is unknown.
  - Capability use MUST respect policy constraints.

Compatibility Rules:
  - Capability version compatibility MUST follow Versioning Specification.
  - Breaking semantic change requires MAJOR version change.
  - Optional field addition is allowed only with safe absence behavior.
  - Deprecated capability versions MUST remain documented.
  - Removed capability versions MUST fail deterministically when requested.

Invariants:
  - Transactions capability support is explicit.
  - Required Transactions capability is validated before execution.
  - Unsupported required Transactions capability prevents Execution Plan emission.
  - Capability downgrade is never silent.
  - Capability emulation is never implicit.
  - Capability metadata is governed by SDE Control Plane.
  - SDE Data Plane consumes approved capability metadata only.

Relationships:
  Parent:
    - capability.md
  Depends On:
    - capability.md
    - ../versioning/versioning.md
    - ../serialization/serialization.md
    - ../../architecture/runtime/capability-registry.md
    - ../../architecture/data-plane/planning-execution.md
    - ../../architecture/runtime/error-model.md
  Used By:
    - Planning
    - Execution Plan
    - Capability Registry
    - Capability Governance
    - Protocol Plugin Manifest
    - Engine Plugin Manifest
    - Transaction Runtime
    - Protocol Execution
    - Kernel Execution
    - Engine Execution

References:
  - capability.md
  - ../versioning/versioning.md
  - ../serialization/serialization.md
  - ../../architecture/runtime/capability-registry.md
  - ../../architecture/data-plane/planning-execution.md
  - ../../architecture/runtime/error-model.md
