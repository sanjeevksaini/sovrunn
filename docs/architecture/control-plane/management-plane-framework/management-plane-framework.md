# Management Plane Framework

Document:
  ID: management-plane-framework
  Title: Management Plane Framework
  Parent: sde-control-plane
  Owner: SDE Control Plane
  Layer: SDE Control Plane
  Type: Architecture
  Version: 1.0
  Status: Draft

Purpose:
  - Define the shared framework for pluggable management planes
  - Define how SDE Control Plane hosts management domains
  - Define admission, lifecycle, controller runtime, conformance, and governance requirements
  - Provide the foundation for Datastore Management Plane and future management planes

Definition:
  Management Plane Framework is the SDE Control Plane framework that allows management domains to be added as governed, pluggable planes.

  A pluggable management plane owns domain-specific lifecycle, resources, workflows, and management APIs while using SDE Control Plane Foundation Services for identity, authorization, policy, workflow, audit, eventing, observability, registry, secrets, and configuration.

Core Principle:
  Management planes are not uncontrolled plugins.

  They are governed Control Plane domains admitted through manifest validation, policy, conformance, and lifecycle controls.

Responsibilities:
  - Register management plane types
  - Validate Management Plane Manifest
  - Admit or reject management planes
  - Govern management plane lifecycle
  - Bind required Foundation Services
  - Expose management plane APIs
  - Provide controller runtime integration
  - Enforce policy, workflow, audit, and observability boundaries
  - Validate management plane conformance

Components:
  - Management Plane Registry
  - Management Plane Manifest
  - Management Plane Admission
  - Management Plane Controller Runtime
  - Management Plane Conformance
  - Management Plane Lifecycle

First Management Plane:
  Datastore Management Plane:
    Role:
      - First pluggable management plane.
      - Manages tenant-scoped Downstream Datastore lifecycle and operations.
      - Powers dstoreOps.

Future Management Planes:
  Examples:
    - Cache Management Plane
    - Search Management Plane
    - Vector Management Plane
    - Data Pipeline Management Plane
    - Tenant Integration Management Plane

Boundaries:
  Management Plane Framework Must Not:
    - Contain datastore-specific lifecycle logic
    - Execute tenant data-plane requests
    - Bypass Foundation Services
    - Replace SDE Data Plane runtime
    - Replace plugin registries

Invariants:
  - SDE Control Plane hosts pluggable management planes through this framework.
  - DMP is admitted through this framework.
  - Management planes must be policy-governed and audited.
  - Management planes must not execute tenant data-plane requests.
  - Management planes must define manifests and pass conformance before production admission.

Related Documents:
  - ../control-plane.md
  - ../management-plane.md
  - management-plane-registry.md
  - management-plane-manifest.md
  - management-plane-controller-runtime.md
  - management-plane-admission.md
  - management-plane-conformance.md
  - ../datastore-management-plane/datastore-management-plane.md
