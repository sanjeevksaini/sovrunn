---
doc_type: architecture
title: Adapter Boundary Model
status: draft
phase: 2
ai_load_priority: always
ai_summary: Adapter boundaries that allow Sovrunn to reuse mature OSS while preventing future recoding.
---

# Adapter Boundary Model

## Purpose

Sovrunn must reuse mature infrastructure without coupling core business logic directly to any single implementation.

## Adapter Boundaries

| Area | Adapter | Reuse Direction |
|---|---|---|
| Policy | PolicyEngineAdapter | OPA, Cedar |
| Identity | IdentityProviderAdapter | Keycloak, Dex, OIDC/SAML |
| Secrets | SecretProviderAdapter | Vault, External Secrets, Kubernetes Secrets for local MVP |
| Operations | OperationEngineAdapter | Temporal, Argo Workflows, Kubernetes Jobs, controller-runtime |
| Observability | ObservabilityAdapter | OpenTelemetry, Prometheus, Grafana, Loki, Tempo |
| Events | EventBusAdapter | Redpanda/Kafka, NATS, CloudEvents |
| Persistence | Repository interfaces | PostgreSQL, YugabyteDB, SQLite for tests, CRDs for some states |

## Rule

MVP implementations may be simple, but they must sit behind adapter interfaces where replacement is expected.
