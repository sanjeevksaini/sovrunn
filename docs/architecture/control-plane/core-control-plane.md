# Core Control Plane

Document
- ID: core-control-plane
- Version: 1.0
- Status: Stable

Purpose
- Define built-in runtime governance Management Plane
- Separate runtime governance from request execution
- Point to registry and governance contracts

Definition
Core Control Plane is the built-in Management Plane responsible for SDE runtime governance.

Responsibilities

MUST
- Govern runtime, plugin, engine, capability, and deployment metadata
- Publish approved state for SDE Data Plane consumption
- Preserve compatibility, versioning, and rollout governance

MUST NOT
- Execute client data requests
- Invoke Engine Plugins
- Invoke Datastore Operator Plugins
- Manage downstream datastore lifecycle

Boundaries

Owns
- Runtime Registry
- Plugin Registry
- Engine Registry
- Capability Governance
- Deployment Governance

Does Not Own
- SDE Data Plane execution
- Datastore lifecycle management
- Foundation Provider implementation

References
- core-control-plane/core-control-plane.md
- core-control-plane/runtime-registry.md
- core-control-plane/plugin-registry.md
- core-control-plane/engine-registry.md
- core-control-plane/capability-governance.md
- core-control-plane/deployment-governance.md
