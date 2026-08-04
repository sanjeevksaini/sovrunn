---
doc_type: mvp_scope
title: MVP-001 Out of Scope
status: draft
phase: 3
ai_load_priority: high
ai_summary: Explicit non-goals for the first governed PostgreSQL PaaS MVP.
---

# MVP-001 Out of Scope

MVP-001 must not include:

- real AWS/Azure/VMware/OpenStack integration,
- real multi-provider placement execution,
- real cross-cloud failover,
- full PostgreSQL HA implementation by Sovrunn,
- full backup/restore automation,
- full global traffic management,
- full autoscaling,
- spot capacity execution,
- full billing/chargeback,
- compliance certification,
- autonomous AI execution,
- full UI/portal,
- production-grade plugin sandbox,
- custom identity provider,
- custom secret manager,
- custom workflow engine.

These are later-phase capabilities or mature systems to reuse/wrap.

## Phase 2R Additional Non-Goals

**ADH-2026-042**

MVP-001 also does not include:

- mandatory ResourcePool or ProviderCapability in core,
- native AWS/OCI/OpenStack/Kubernetes/OpenShift objects in customer schemas,
- dual authority between alpha and canonical models,
- capacity scheduling or provider-wide capability truth,
- autonomous AI decision authority,
- stable API promotion during Phase 2R.