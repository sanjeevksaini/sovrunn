# Plugin Runtime

Document:
  ID: plugin-runtime
  Title: Plugin Runtime
  Parent: runtime
  Owner: Plugin Runtime
  Layer: SDE Data Plane
  Type: Component Architecture
  Version: 1.0
  Status: Stable

Purpose:
  - Define Plugin Runtime
  - Define runtime position
  - Define execution responsibilities
  - Define boundaries and failure rules

Definition:
  Plugin Runtime loads, isolates, and governs runtime plugin instances used by SDE Runtime.

Runtime Position:
  - Supports Protocol Plugins and Engine Plugins.
  - Consumes approved plugin metadata.
  - Provides runtime plugin lifecycle.

Responsibilities:
  MUST:
    - Load approved runtime plugins
    - Validate plugin compatibility
    - Manage plugin lifecycle
    - Provide plugin isolation
    - Expose plugin availability

  MUST NOT:
    - Approve plugin governance
    - Bypass Plugin Registry
    - Execute datastore lifecycle operations
    - Expose unsafe plugin internals

Inputs:
  - Approved plugin metadata
  - Runtime configuration
  - Plugin lifecycle action

Outputs:
  - Plugin instance
  - Plugin availability state
  - Plugin lifecycle error

State:
  - Plugin instance state
  - Compatibility metadata
  - Lifecycle state

Execution Rules:
  - Load only approved plugins.
  - Enforce plugin isolation.
  - Fail closed when plugin metadata is invalid.

Failure Rules:
  - Return plugin lifecycle error.
  - Quarantine failed plugin where required.

Concurrency Rules:
  - Preserve execution isolation.
  - Avoid shared mutable execution state unless explicitly synchronized.
  - Preserve tenant isolation across concurrent executions.
  - Preserve session and transaction boundaries.

Security Rules:
  - Enforce authorized execution context.
  - Avoid exposing secrets.
  - Preserve safe error behavior.
  - Preserve trace and audit correlation where applicable.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - ../control-plane/core-control-plane/plugin-registry.md
  Used By:
    - protocol-runtime.md
    - engine-runtime.md

References:
  - runtime.md
  - runtime-map.md
  - ../control-plane/core-control-plane/plugin-registry.md
