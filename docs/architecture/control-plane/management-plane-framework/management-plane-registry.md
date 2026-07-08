# Management Plane Registry

Document:
  ID: management-plane-registry
  Title: Management Plane Registry
  Parent: management-plane-framework
  Owner: SDE Control Plane
  Layer: SDE Control Plane
  Type: CONTRACT
  Version: 1.0
  Status: Draft

Purpose:
  - Define registry responsibilities for pluggable management planes
  - Track admitted management planes, versions, capabilities, lifecycle state, and compatibility

Definition:
  Management Plane Registry is the authoritative Control Plane registry for pluggable management plane metadata.

Responsibilities:
  - Store management plane identity
  - Store version and compatibility metadata
  - Store manifest reference
  - Store lifecycle state
  - Store required Foundation Service dependencies
  - Store controller runtime requirements
  - Store admission status
  - Store conformance status

Lifecycle States:
  - Registered
  - Validating
  - Admitted
  - Enabled
  - Disabled
  - Deprecated
  - Rejected
  - Removed

Rules:
  - Only admitted management planes may be enabled.
  - Production enablement requires conformance status.
  - Registry state changes must be audited.
  - Registry admission must be policy-governed.

Invariants:
  - Registry is authoritative for management plane lifecycle metadata.
  - Registry does not execute management workflows.
  - Registry does not execute tenant data requests.
