---
doc_type: foundation
title: Sovrunn Constitution
status: draft
phase: 0
ai_load_priority: always
ai_summary: Non-negotiable rules for Sovrunn governance, ServiceOps, AI, SDE, and engineering.
---

# Sovrunn Constitution

## 1. Purpose

This constitution defines non-negotiable engineering and product principles for Sovrunn.

All RFCs, ADRs, implementation decisions, plugins, and generated code must comply with this file unless a formal superseding decision is accepted.

## 2. Constitutional Principles

### Principle 1: Sovrunn is organization-first

Sovrunn must support Organization and OrganizationUnit layers above Tenant.

Canonical hierarchy:

```text
Organization
  -> OrganizationUnit
      -> Tenant
          -> Project
              -> ServiceInstance
```

### Principle 2: Tenant consumption must be isolated

Tenants must consume services inside explicit isolation profiles.

Supported isolation profiles:

- namespace,
- vCluster,
- dedicated cluster.

### Principle 3: Policy must be inherited and non-weakenable

Policy resolution:

```text
Organization baseline
  -> OrganizationUnit policy
      -> Tenant policy
          -> Project policy
              -> ServiceInstance enforcement
```

Lower layers may strengthen policy.

Lower layers must not weaken Organization baseline policy unless an explicit exception workflow is supported and audited.

### Principle 4: Sovrunn builds on open-source infrastructure

Sovrunn must not rebuild mature open-source infrastructure unless there is a strong architectural reason.

Sovrunn should use proven systems for:

- Kubernetes,
- GitOps,
- identity,
- policy,
- observability,
- secrets,
- networking,
- service operators,
- gateways,
- load balancing,
- FaaS runtimes,
- storage and data services.

### Principle 5: Service delivery must be plugin-based

All managed service families must use a common ServiceOps pattern.

Plugin families include:

- dStoreOps,
- cacheOps,
- objectOps,
- streamOps,
- gatewayOps,
- lbOps,
- faasOps,
- bigDataOps,
- sdeOps.

Plugins must declare capabilities, supported operations, dependencies, versions, security requirements, and conformance status.

### Principle 6: Management-plane plugins may be remote; data-path plugins must be local by default

Remote plugins are acceptable for asynchronous management-plane operations.

Data-path plugins, especially in SDE hot paths, must be in-process by default unless an RFC proves remote execution does not violate latency, reliability, or correctness.

### Principle 7: Operations must be asynchronous, auditable, and traceable

Every meaningful platform change must create or link to an Operation.

Operation must include:

- requester,
- organization,
- organization unit,
- tenant,
- project,
- resource,
- action,
- policy decision,
- plugin used,
- status,
- timestamps,
- logs/events,
- correlation ID,
- audit record.

### Principle 8: AI must operate through governed tools

AI agents must not bypass Sovrunn APIs, policy validation, approval workflow, tenant boundaries, secret handling, or audit logging.

AI may:

- explain,
- recommend,
- draft,
- validate,
- diagnose,
- generate manifests,
- propose operations,
- create runbooks,
- assist plugin development.

### Principle 9: AI must not be a mandatory synchronous dependency in latency-sensitive paths

AI must not be required for SDE SQL execution or other latency-sensitive runtime paths.

### Principle 10: SDE must remain semantically honest

SDE must not claim universal compatibility across all protocols and datastores.

Runtime transformation is allowed only when:

- request can be represented in SIR,
- source semantics are understood,
- target datastore capabilities are declared,
- transformation mapping exists,
- unsupported behavior is rejected or safely routed,
- correctness can be explained.

### Principle 11: SDE hot path must be minimal

SDE data hot path must avoid synchronous calls to:

- Control Plane,
- Management Plane,
- dStoreOps,
- AI agents,
- remote policy engines,
- remote plugin services,
- telemetry backends.

### Principle 12: Fail fast and explain clearly

Sovrunn must fail fast when requests are invalid, unsupported, unsafe, or not entitled.

Failures must be explicit and actionable.

### Principle 13: Open standards are preferred

Sovrunn should prefer:

- Kubernetes API and CRDs,
- OIDC,
- SAML,
- OpenTelemetry,
- Prometheus metrics,
- GitOps,
- OCI images,
- S3-compatible APIs,
- PostgreSQL wire protocol,
- open ServiceOps contracts,
- open SIR and TransformationMapping specifications when mature.

### Principle 14: Commercial features must not break open trust

Sovrunn may use an open-core commercial model.

The open core should be sufficient for learning, local development, plugin creation, and basic service management.

Enterprise features may provide advanced governance, federation, certified plugins, support, compliance, AI/AOE capabilities, and advanced SDE features.

### Principle 15: Documentation and decision traceability are mandatory

Every major platform capability must be backed by:

- DEC entry,
- ADR where tradeoffs exist,
- RFC for detailed design,
- resource/API specification,
- implementation notes,
- tests,
- runbook,
- example manifests.

## 3. Engineering Rules

### Rule 1: Build the platform factory

Every feature should follow:

```text
RFC
ADR
Resource schema
API contract
Controller
Plugin interface
Tests
Docs
Demo
Runbook
```

### Rule 2: Prefer thin vertical slices

Build one service family end-to-end before expanding to many.

Initial wedge:

```text
PostgreSQL PaaS through Sovrunn ServiceOps
```

### Rule 3: Use proven components before custom components

Custom implementation is allowed only when Sovrunn’s differentiator requires it.

### Rule 4: Every service must expose health, status, events, and audit

A service instance is not platform-ready unless Sovrunn can show:

- desired state,
- actual state,
- operation history,
- health,
- logs/metrics/traces,
- policy status,
- backup status where applicable,
- binding status,
- audit history.

### Rule 5: Tenant data and secrets must not leak into AI

AI prompts, logs, traces, and model calls must not expose secrets or unauthorized tenant data.

Redaction and tenant boundary enforcement are mandatory for AI tooling.

## 4. Non-Goals

Sovrunn is not:

- Kubernetes replacement,
- identity provider replacement,
- observability backend replacement,
- database operator replacement,
- hyperscaler clone,
- generic chatbot,
- uncontrolled autonomous operations system,
- universal SQL-to-everything compatibility promise.

## 5. Summary

```text
Organization-first governance.
Tenant-isolated consumption.
Open-source-backed enforcement.
Plugin-based service delivery.
AI-assisted but policy-governed automation.
Auditable operations.
Fail-fast validation.
Minimal hot paths.
Semantic honesty.
Open standards.
Enterprise-ready productization.
```

## 6. Phase 2 Superseding Principles

These principles extend the original Phase 0/1 constitution for the evolved Sovrunn scope.

### Principle 16: Reuse before build

Sovrunn must not reinvent mature open-source infrastructure, protocols, controllers, policy engines, observability stacks, workflow engines, service operators, identity systems, secret managers, networking stacks, or database lifecycle systems unless a formal decision approves it.

Sovrunn should build the sovereign PaaS control plane, decision layer, plugin contracts, governance model, AI-readable operations context, and customer/provider experience that connects mature components into one governed platform.

### Principle 17: Adapter boundary before integration

Whenever Sovrunn expects to use or replace an external engine later, the core must depend on an adapter boundary, not the concrete implementation.

Required adapter boundaries include:

- PolicyEngineAdapter,
- IdentityProviderAdapter,
- SecretProviderAdapter,
- OperationEngineAdapter,
- ObservabilityAdapter,
- EventBusAdapter,
- repository interfaces.

### Principle 18: Policy abstraction before policy engine

Sovrunn must define PolicyEvaluationRequest and PolicyEvaluationResult before binding to OPA, Cedar, or any other engine.

Go-based policy logic is allowed only as a temporary bootstrap adapter behind PolicyEngineAdapter, not as embedded business logic in handlers or placement.

### Principle 19: Decision-first execution

Sovrunn must create explainable decisions before executing platform changes.

No provisioning, scaling, failover, traffic, backup, data movement, or runtime operation should proceed without the relevant decision context.

### Principle 20: Customer API boundary protection

Customer-facing APIs must expose PaaS outcomes, not low-level provider or IaaS details by default.

Provider, resource pool, failure domain, and capability details may be provider/MSP-facing or internal unless explicitly surfaced as an advanced view.
