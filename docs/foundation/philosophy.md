---
doc_type: foundation
title: Sovrunn Philosophy
status: draft
phase: 0
ai_load_priority: always
ai_summary: Defines the operating beliefs that guide Sovrunn design and AI-assisted development.
---

# Sovrunn Philosophy

## 1. Core Philosophy

Sovrunn turns proven open-source cloud-native building blocks into a governed sovereign PaaS platform.

Sovrunn does not rebuild the ecosystem.

Sovrunn unifies, governs, productizes, extends, and safely automates it.

```text
Reuse proven infrastructure.
Own the sovereign product layer.
Centralize governance.
Isolate tenant consumption.
Expose services through plugins.
Use AI for acceleration.
Keep execution deterministic.
Make every action auditable.
```

## 2. Organization-First

Sovrunn is organization-first, not tenant-only.

Canonical hierarchy:

```text
Organization
  -> OrganizationUnit
      -> Tenant
          -> Project
              -> ServiceInstance
```

The Organization owns central governance.

OrganizationUnit receives delegated authority.

Tenant consumes isolated services.

Project separates environments and workloads.

## 3. Sovereign by Design

Sovereignty means:

- local control,
- auditability,
- portability,
- clear ownership,
- local operability,
- no hidden dependency on one hyperscaler,
- no opaque control plane,
- no unnecessary lock-in.

Sovrunn should run on government cloud, local cloud, colocation, private cloud, on-premise Kubernetes, and edge/regional data centers.

## 4. Open-Source Substrate, Sovrunn Product Layer

Sovrunn should use open-source systems for common infrastructure capabilities:

- Kubernetes for orchestration,
- Argo CD or Flux for GitOps,
- Keycloak for identity,
- Kyverno or OPA for policy,
- OpenTelemetry for telemetry,
- Prometheus, Grafana, Loki, Tempo for observability,
- External Secrets and Vault for secrets,
- Cilium/Calico for networking,
- MetalLB and gateway projects for local cloud networking,
- service operators for PostgreSQL, Redis/Valkey, MinIO, Kafka/Redpanda, Knative, and others.

Sovrunn should build:

- organization and tenant model,
- service catalog,
- service plans,
- operation framework,
- plugin registry,
- capability registry,
- ServiceOps SDK,
- policy inheritance,
- audit aggregation,
- backup and archival governance,
- AI-assisted operations,
- SDE.

## 5. Plugin Ecosystem Philosophy

Every managed service family should use the ServiceOps model.

Plugin families:

- dStoreOps for databases,
- cacheOps for cache,
- objectOps for object storage,
- streamOps for messaging and streaming,
- gatewayOps for API gateways,
- lbOps for load balancers,
- faasOps for serverless/FaaS,
- bigDataOps for big data processing,
- sdeOps for SDE services.

Common lifecycle:

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

## 6. AI-Assisted, Not AI-Uncontrolled

AI is valuable for:

- generating plans,
- explaining failures,
- creating manifests,
- recommending remediation,
- producing runbooks,
- summarizing observability signals,
- assisting plugin development,
- analyzing migration feasibility,
- supporting SDE transformation planning.

AI must not directly mutate production systems without policy, validation, approval, and audit.

Operating rule:

```text
AI proposes.
Policy validates.
Humans approve when needed.
Controllers execute.
Audit records everything.
```

## 7. Determinism and Fail-Fast Behavior

Sovrunn must be predictable.

Sovrunn should prefer:

- deterministic controllers,
- explicit resources,
- capability checks,
- policy validation,
- clear failure states,
- auditable operations.

Sovrunn should fail fast when:

- tenant lacks entitlement,
- service plan is unsupported,
- plugin lacks capability,
- policy blocks the request,
- quota is exceeded,
- dependency is missing,
- SDE transformation is semantically unsafe.

## 8. Data-Path Minimalism

Latency-sensitive runtime paths must avoid synchronous dependencies on:

- AI agents,
- Control Plane,
- Management Plane,
- remote policy services,
- remote plugin calls,
- telemetry backends,
- governance dashboards.

SDE runtime should use local snapshots, local policy context, in-process plugins, and precomputed capability data where required.

## 9. Semantic Honesty

SDE must be semantically honest.

If a transformation is safe, capability-gated, and mapped, SDE may support it.

If semantics are unsupported, SDE must reject, route, or require migration design.

SDE must not pretend compatibility where correctness cannot be guaranteed.

## 10. Founder Engineering Philosophy

Sovrunn should be built using a repeatable platform factory:

```text
RFC
ADR
resource schema
API contract
controller
plugin interface
tests
docs
demo
AI coding prompt
```

Every major feature should produce:

- working code,
- automated tests,
- documentation,
- demo flow,
- decision traceability,
- operational runbook.
