---
doc_type: rfc
id: RFC-0020
title: ServiceOps Plugin Framework
status: draft
phase: 1
depends_on:
  - DEC-0012
  - DEC-0023
  - constitution.md
ai_load_priority: feature
ai_summary: Defines the common ServiceOps plugin model, plugin manifest, capability declaration, lifecycle operations, and conformance expectations.
---

# RFC-0020: ServiceOps Plugin Framework

## 1. Status

Draft for founder review.

## 2. Purpose

Define the ServiceOps plugin framework used by Sovrunn to manage PaaS service families.

ServiceOps provides a common lifecycle model across:

- databases,
- cache,
- object storage,
- streams,
- gateways,
- load balancers,
- FaaS,
- big data,
- SDE.

## 3. Goals

- Define common plugin model.
- Define plugin manifest.
- Define capability declaration.
- Define lifecycle operations.
- Define conformance expectations.
- Enable first implementation: `postgres.dStoreOps`.

## 4. Non-Goals

This RFC does not implement:

- plugin marketplace,
- remote plugin runtime,
- plugin billing,
- plugin signing,
- full conformance suite,
- production sandboxing,
- SDE data-path plugin framework.

## 5. Definitions

| Term | Definition |
|---|---|
| ServiceOps | Generic lifecycle framework for PaaS service management. |
| Plugin | Implementation provider for a service family or provider. |
| PluginManifest | Metadata document describing plugin identity, version, kind, capabilities, dependencies, and deployment mode. |
| Capability | Declared supported operation or feature. |
| OperationHandler | Plugin function that performs a lifecycle action. |
| ConformanceTest | Test verifying plugin behavior against declared contract. |

## 6. Plugin Families

| Family | Purpose |
|---|---|
| dStoreOps | Datastore operations. |
| cacheOps | Cache operations. |
| objectOps | Object storage operations. |
| streamOps | Stream/messaging operations. |
| gatewayOps | API gateway operations. |
| lbOps | Load balancer operations. |
| faasOps | FaaS operations. |
| bigDataOps | Big data processing operations. |
| sdeOps | SDE management operations. |

## 7. Lifecycle Operations

Common operations:

```text
Validate
Plan
Provision
Configure
Bind
Observe
Scale
Upgrade
Backup
Restore
RotateCredentials
Unbind
Delete
```

Plugins do not need to support every operation.

Unsupported operations must be declared and fail fast with `CapabilityUnsupported`.

## 8. Plugin Manifest

Minimum fields:

| Field | Required | Description |
|---|---:|---|
| name | yes | Plugin name. |
| family | yes | Plugin family. |
| provider | yes | Provider/service implementation. |
| version | yes | Plugin version. |
| deploymentMode | yes | compiled-in, sidecar, remote. |
| supportedOperations | yes | Lifecycle operations supported. |
| requiredPermissions | no | Required permissions. |
| requiredSecrets | no | Secret references needed. |
| dependencies | no | Operators, CRDs, Helm charts, APIs. |
| status | yes | draft, active, deprecated. |

Example:

```yaml
name: postgres.dstoreops.basic
family: dStoreOps
provider: postgresql
version: 0.1.0
deploymentMode: compiled-in
supportedOperations:
  - Validate
  - Plan
  - Provision
  - Bind
  - Observe
  - Delete
dependencies:
  - postgresql-operator
status: draft
```

## 9. Capability Declaration

Minimum fields:

| Field | Required | Description |
|---|---:|---|
| pluginRef | yes | Plugin name/version. |
| capability | yes | Capability name. |
| serviceClass | yes | Associated ServiceClass. |
| servicePlans | no | Supported ServicePlans. |
| constraints | no | Capability constraints. |

## 10. Operation Handler Contract

A plugin operation handler should receive:

- Operation context,
- Organization scope,
- OrganizationUnit scope,
- Tenant scope,
- Project scope,
- ServiceInstance spec,
- ServicePlan,
- resolved parameters,
- SecretRef references,
- correlation ID.

A plugin operation handler should return:

- operation phase,
- reason code if failed,
- message,
- output references,
- status details,
- emitted events.

## 11. Deployment Modes

| Mode | Phase | Description |
|---|---|---|
| compiled-in | early | Plugin compiled into management plane binary. |
| sidecar | later | Plugin deployed alongside management plane. |
| remote | later | Plugin exposed as remote service. |

Phase 1 should model deployment mode but may implement only compiled-in behavior.

## 12. Data-Path Boundary

ServiceOps is for management-plane operations.

ServiceOps plugins must not be called synchronously in SDE SQL hot path.

Data-path plugin rules are separate and stricter.

## 13. API Behavior

Minimum Phase 1 behavior:

- register Plugin,
- list Plugins,
- register Capability,
- list Capabilities,
- resolve Plugin for ServiceClass/ServicePlan,
- fail fast if capability missing.

## 14. Validation Rules

- Plugin name/version must be unique.
- Plugin family must be known.
- deploymentMode must be supported.
- supportedOperations must be known lifecycle operations.
- Capability must reference existing Plugin.
- Capability serviceClass must exist when service catalog is enabled.
- Unsupported operation must return `CapabilityUnsupported`.

## 15. Failure Modes

| Failure | Reason Code |
|---|---|
| Unknown plugin family | PluginFamilyUnsupported |
| Duplicate plugin | AlreadyExists |
| Unsupported deployment mode | DeploymentModeUnsupported |
| Unknown lifecycle operation | OperationUnsupported |
| Capability references missing plugin | PluginNotFound |
| Required capability missing | CapabilityUnsupported |

## 16. Security and Governance

Plugin manifest must declare:

- permissions,
- secret needs,
- dependencies,
- external endpoints where applicable,
- deployment mode.

Do not grant broad cluster permissions by default.

## 17. Observability and Audit

Every plugin operation must be linked to an Operation.

Every plugin failure must include:

- reason code,
- message,
- plugin name,
- operation ID,
- correlation ID.

## 18. Tests

Required tests:

- register plugin success,
- duplicate plugin failure,
- unknown family failure,
- unsupported deployment mode failure,
- register capability success,
- capability missing plugin failure,
- unsupported operation failure,
- plugin resolution success.

## 19. Acceptance Criteria

- PluginManifest model exists.
- Capability model exists.
- Plugin registry exists.
- Capability registry exists.
- ServiceInstance can resolve a plugin later.
- Unsupported capabilities fail fast.
- Phase 1 does not implement remote plugin runtime.

## 20. Related Decisions

- DEC-0012
- DEC-0023
- DEC-0025

## 21. AI Implementation Guidance

- Implement minimal compiled-in plugin model first.
- Do not implement remote plugin runtime.
- Do not build marketplace.
- Do not add billing.
- Keep ServiceOps separate from SDE hot path.
- Add tests for plugin validation and capability resolution.
