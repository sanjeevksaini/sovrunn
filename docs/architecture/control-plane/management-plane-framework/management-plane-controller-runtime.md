# Management Plane Controller Runtime

Document:
  ID: management-plane-controller-runtime
  Title: Management Plane Controller Runtime
  Parent: management-plane-framework
  Owner: SDE Control Plane
  Layer: SDE Control Plane
  Type: ARCHITECTURE
  Version: 1.0
  Status: Draft

Purpose:
  - Define runtime responsibilities for hosting pluggable management plane controllers
  - Clarify the difference between a management plane and its controller runtime
  - Clarify DMP Controller Runtime semantics

Definition:
  Management Plane Controller Runtime is the executable runtime that hosts management plane controllers and reconciles management-plane resources and workflows.

  It is not the management plane itself.

DMP Example:
  DMP:
    Meaning:
      - Datastore Management Plane as a pluggable management plane.

  sde-dmp-controller:
    Meaning:
      - DMP Controller Runtime executable that hosts and reconciles DMP.

Responsibilities:
  - Start management-plane controllers
  - Watch management-plane resources
  - Reconcile desired and actual state
  - Invoke approved workflows
  - Use policy and authorization checks
  - Emit audit events
  - Emit metrics, traces, and health status
  - Support graceful shutdown

Must Not:
  - Execute tenant data-plane requests
  - Bypass Workflow Service
  - Bypass Policy Service
  - Bypass Audit Service
  - Invoke unauthorized plugins/providers

Invariants:
  - Controller runtime hosts reconciliation.
  - Management plane defines domain semantics.
  - DMP Controller Runtime is not the whole DMP.
