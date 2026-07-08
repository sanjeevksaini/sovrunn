# AI Control Plane

Document:
  ID: ai-control-plane
  Title: AI Control Plane
  Parent: sde-control-plane
  Owner: SDE Control Plane
  Layer: SDE Control Plane
  Type: Architecture Placeholder
  Version: 1.0
  Status: Reserved

Purpose:
  - Reserve AI Control Plane as a future pluggable SDE Control Plane extension
  - Define AI governance boundaries without over-designing AI internals
  - Preserve future tenant-facing AI Agent capability
  - Ensure future AI capabilities integrate through approved Control Plane services

Definition:
  AI Control Plane is an optional, pluggable, policy-governed SDE Control Plane extension.

  It may provide AI-assisted observation, recommendation, validation, tenant assistance, workflow initiation, and platform stabilization.

  Detailed AI Control Plane internals are intentionally deferred.

Core Principle:
  AI Control Plane is a governed Control Plane extension.

  AI may observe, explain, recommend, validate, and initiate approved workflows.

  AI must not bypass SDE Control Plane services, Datastore Management Plane, Datastore Operator Plugins, Infrastructure Providers, authorization, policy, workflow, or audit.

Scope:
  In Scope:
    - Architectural reservation
    - Governance boundaries
    - Tenant AI Agent placeholder
    - Integration principles
    - Non-bypass rules
    - Future RFC range reservation

  Out of Scope:
    - AI model implementation
    - Prompt architecture
    - AI agent runtime implementation
    - Autonomous operations internals
    - Plugin code generation
    - AI-generated plugin deployment
    - Detailed tenant AI workflows

Tenant AI Agent:
  Tenant AI Agent is a future tenant-facing AI interface that may help customers configure, integrate, monitor, and operate tenant-scoped SDE resources.

  Tenant AI Agent must operate through approved tenant-scoped Control Plane APIs and workflows.

  Tenant AI Agent must not directly manage Downstream Datastores, infrastructure, plugins, secrets, or tenant data.

Allowed Future Capabilities:
  - Explain SDE concepts and tenant configuration
  - Analyze tenant-scoped telemetry
  - Recommend datastore configuration
  - Generate tenant configuration drafts
  - Generate connection examples
  - Generate workflow requests
  - Recommend non-destructive tuning
  - Suggest remediations
  - Initiate approved workflows
  - Track workflow status
  - Summarize operational health

Explicitly Excluded From Initial Scope:
  - On-the-fly Protocol Plugin generation
  - On-the-fly Engine Plugin generation
  - On-the-fly Datastore Operator Plugin generation
  - Automatic plugin registration
  - Automatic production plugin deployment
  - Destructive datastore operations
  - Direct infrastructure changes
  - Direct credential access
  - Bypassing policy or approval workflows

Integration Boundaries:
  AI Control Plane MAY integrate with:
    - Identity Service
    - Authorization Service
    - Policy Service
    - Workflow Service
    - Audit Service
    - Observability Service
    - Eventing Service
    - Secrets Service through approved references
    - Tenant Management Service
    - Datastore Management Plane
    - Core Control Plane registries

  AI Control Plane MUST NOT bypass:
    - Authorization Service
    - Policy Service
    - Workflow Service
    - Audit Service
    - Datastore Management Plane
    - Datastore Operator Plugins
    - Infrastructure Providers
    - Secrets Service controls
    - Plugin Registry admission
    - Engine Registry governance
    - Capability Governance

Tenant Datastore Management Principle:
  A future Tenant AI Agent may help configure, provision, and manage tenant-specific datastores in tenant-specific namespaces.

  Such actions must flow through:
    Tenant AI Agent
      → Tenant-scoped Control Plane API
      → Authorization and Policy
      → Workflow Service
      → Datastore Management Plane
      → Datastore Operator Plugin
      → Infrastructure Provider where required
      → Tenant-specific Downstream Datastore

  AI Agent must not call Datastore Operator Plugins or infrastructure APIs directly.

Risk Classes:
  Class 0:
    Meaning: Read-only observation and explanation.

  Class 1:
    Meaning: Generate draft recommendations or configuration artifacts.

  Class 2:
    Meaning: Validate or simulate proposed action.

  Class 3:
    Meaning: Apply reversible, non-destructive, policy-approved tuning.

  Class 4:
    Meaning: Controlled operational change requiring workflow approval.

  Class 5:
    Meaning: Destructive or irreversible action requiring explicit human approval and safety workflow.

Rules:
  - Class 0 and Class 1 may be automated when tenant is authorized.
  - Class 2 may be automated when validation environment is approved.
  - Class 3 may be automated only when policy allows rollback-safe action.
  - Class 4 requires approved workflow.
  - Class 5 requires explicit human approval, safety workflow, and audit.
  - AI Control Plane must preserve tenant isolation.
  - AI Control Plane must preserve auditability.
  - AI Control Plane must not mutate authoritative state outside approved workflows.

Invariants:
  - AI Control Plane is optional and pluggable.
  - AI Control Plane is part of SDE Control Plane, not SDE Data Plane.
  - AI Control Plane does not bypass Control Plane services.
  - Tenant AI Agent is tenant-scoped.
  - AI-generated recommendations are not authoritative until approved.
  - AI-generated artifacts are untrusted until validated.
  - AI does not directly manage Downstream Datastores.
  - AI does not directly invoke Datastore Operator Plugins.
  - AI does not directly invoke Infrastructure Providers.

Relationships:
  Parent:
    - control-plane.md
  Depends On:
    - control-plane-foundation.md
    - management-plane.md
    - foundation-services/foundation-services.md
    - datastore-management-plane/datastore-management-plane.md
    - core-control-plane/plugin-registry.md
    - core-control-plane/capability-governance.md
  Used By:
    - Future AI RFCs
    - Tenant AI Agent design
    - Autonomous operations design
    - Tenant datastore workflow design

References:
  - control-plane.md
  - control-plane-map.md
  - control-plane-foundation.md
  - management-plane.md
  - foundation-services/foundation-services.md
  - datastore-management-plane/datastore-management-plane.md
  - core-control-plane/plugin-registry.md
  - core-control-plane/capability-governance.md
