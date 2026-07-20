# Sovrunn Context Pack

This is the compact source-of-truth context for new ChatGPT, Kiro, Cursor, and reviewer sessions.

Use with `docs/context/CURRENT_ARCHITECTURE_BASELINE.md` and `docs/context/CHATGPT_ARCHITECTURE_SESSION_PROMPT.md`.

## Product Definition

Sovrunn is a cloud-native sovereign PaaS platform for local cloud providers, MSPs, and on-premise cloud operators.

It provides governed service catalog, organization/tenant/project governance, provider-neutral placement, plugin-based service lifecycle, decision/audit/evidence records, and AI-assisted operations.

Sovrunn Data Engine is a future managed service inside Sovrunn.

## Core Build Boundary

Sovrunn builds:

- governance,
- policy context,
- decision models,
- placement decisions,
- plugin contracts,
- operation tracking,
- audit and evidence,
- AI-readable explanation context,
- customer/provider PaaS experience.

Sovrunn reuses or wraps mature infrastructure such as Kubernetes, OPA, Cedar, Keycloak, Vault, OpenTelemetry, Prometheus, Grafana, Temporal/Argo, PostgreSQL operators, ingress controllers, and backup tools where appropriate.

## Current Architecture Spine

```text
Service request
  -> entitlement check
  -> effective policy context
  -> policy evaluation
  -> placement request
  -> placement decision
  -> operation
  -> plugin execution
  -> service instance status
  -> service binding
  -> audit event
  -> AI-readable explanation
```

## Current Phase

Current active phase: Phase 2.

Phase 2 builds model, decision, audit, adapter, plugin taxonomy, and placement simulation foundations.

Phase 2 does not build real provider provisioning or real PostgreSQL runtime provisioning.

## MVP

MVP-001: Governed PostgreSQL PaaS Placement and Provisioning on one substrate.

## Plugin Planes

1. Provider/Substrate Plugin: infrastructure execution.
2. PaaS Service Management Plane Plugin: service lifecycle planning.
3. PaaS Service Runtime Plugin: runtime configuration, binding, readiness, and status.

## Reuse Targets

- Policy: OPA/Cedar through `PolicyEngineAdapter`.
- Identity: Keycloak/Dex through `IdentityProviderAdapter`.
- Secrets: Vault/External Secrets through `SecretProviderAdapter`.
- Workflow: Temporal/Argo through `OperationEngineAdapter`.
- Observability: OpenTelemetry, Prometheus, Grafana, Loki, Tempo.
- PostgreSQL runtime: CloudNativePG, Crunchy Postgres Operator, or Helm wrapper.

## Current Phase 2 Features

`FEATURE-0011` through `FEATURE-0026` define Phase 2.

## Roadmap Rule

Roadmap placeholders for Phase 4+ are directional scope references. They must be revalidated after Phase 2 and Phase 3 before detailed design or implementation.

## Architecture Rule

Every new idea must either:

- fit the current approved architecture,
- become an explicit open question,
- or become a formal DEC/RFC proposal.

Chat history is not source of truth. The repo is source of truth.
