# Foundation Services

Document
- ID: foundation-services-index
- Version: 1.0
- Status: Stable

Purpose
- Provide Foundation Services index
- Define common Foundation Service rules
- Point to one-service-per-file contracts

Definition
Foundation Services are stable SDE Control Plane contracts exposed by Control Plane Foundation.

Responsibilities

MUST
- Provide stable consumer-facing contracts
- Hide Foundation Provider implementation details
- Preserve tenant isolation, policy, audit, and configuration boundaries
- Report failures through canonical error handling

MUST NOT
- Expose provider-specific APIs as canonical contracts
- Store secrets outside Secrets Service
- Create domain lifecycle semantics owned by Management Planes

Relationships
- Foundation Providers implement Foundation Services.
- Management Planes consume Foundation Services.
- SDE runtime components may consume Foundation Services only when explicitly authorized.

References
- foundation-services/foundation-services.md
- foundation-services/identity-service.md
- foundation-services/authorization-service.md
- foundation-services/tenant-management-service.md
- foundation-services/configuration-service.md
- foundation-services/policy-service.md
- foundation-services/secrets-service.md
- foundation-services/audit-service.md
- foundation-services/workflow-service.md
- foundation-services/eventing-service.md
- foundation-services/observability-service.md
- foundation-services/registry-framework.md
- foundation-services/plugin-framework.md
