# Datastore Management Plane

Document
- ID: datastore-management-plane
- Version: 1.0
- Status: Stable

Purpose
- Define Datastore Management Plane
- Define dstoreOps relationship
- Separate datastore lifecycle from SDE Data Plane execution
- Point to registries, plugins, providers, and controllers

Definition
Datastore Management Plane is an optional pluggable Management Plane that powers dstoreOps and manages downstream datastore lifecycle.

Responsibilities

MUST
- Manage downstream datastore lifecycle through authorized workflows
- Use Datastore Operator Plugins for datastore-specific lifecycle operations
- Use Infrastructure Providers for infrastructure environment integration
- Publish approved datastore metadata through SDE Control Plane interfaces when required

MUST NOT
- Execute client data requests
- Process SIR
- Invoke Engine Plugins as execution path
- Replace SDE Data Plane
- Replace Datastore Data Plane

Boundaries

Owns
- dstoreOps
- Datastore Registry
- Datastore Operator Registry
- Infrastructure Provider Registry
- Lifecycle controllers

Does Not Own
- SDE runtime execution
- Engine Plugin execution
- Downstream datastore native query execution

References
- datastore-management-plane/datastore-management-plane.md
- datastore-management-plane/dstoreops.md
- datastore-management-plane/datastore-registry.md
- datastore-management-plane/datastore-operator-registry.md
- datastore-management-plane/infrastructure-provider-registry.md
- datastore-management-plane/datastore-operator-plugin.md
- datastore-management-plane/infrastructure-provider.md
- datastore-management-plane/lifecycle-controller.md
- datastore-management-plane/provisioning-controller.md
- datastore-management-plane/configuration-controller.md
- datastore-management-plane/scaling-controller.md
- datastore-management-plane/backup-controller.md
- datastore-management-plane/restore-controller.md
- datastore-management-plane/patch-controller.md
- datastore-management-plane/upgrade-controller.md
- datastore-management-plane/monitoring-controller.md
- datastore-management-plane/retirement-controller.md
