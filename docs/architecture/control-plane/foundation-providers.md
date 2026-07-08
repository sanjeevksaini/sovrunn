# Foundation Providers

Document
- ID: foundation-providers-index
- Version: 1.0
- Status: Stable

Purpose
- Provide Foundation Providers index
- Define common provider rules
- Point to one-provider-category-per-file contracts

Definition
Foundation Providers are pluggable implementations of Foundation Services.

Responsibilities

MUST
- Implement exactly one Foundation Service contract category unless explicitly composed
- Remain hidden behind Foundation Service interfaces
- Declare capabilities, configuration, lifecycle, and failure modes
- Preserve tenant, policy, audit, and security boundaries

MUST NOT
- Become consumer-facing contracts
- Replace Foundation Services
- Replace Engine Plugins
- Replace Datastore Operator Plugins
- Be confused with Infrastructure Providers

Relationships
- Foundation Services consume providers through controlled binding.
- Management Planes do not depend directly on providers by default.
- Infrastructure Providers belong to Datastore Management Plane, not Control Plane Foundation.

References
- foundation-providers/foundation-providers.md
- foundation-providers/identity-provider.md
- foundation-providers/authorization-provider.md
- foundation-providers/tenant-management-provider.md
- foundation-providers/configuration-provider.md
- foundation-providers/policy-provider.md
- foundation-providers/secrets-provider.md
- foundation-providers/audit-provider.md
- foundation-providers/workflow-provider.md
- foundation-providers/eventing-provider.md
- foundation-providers/observability-provider.md
- foundation-providers/registry-provider.md
- foundation-providers/plugin-provider.md
