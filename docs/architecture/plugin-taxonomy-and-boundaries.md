---
doc_type: architecture
title: Plugin Taxonomy and Boundaries
status: draft
phase: 2
ai_load_priority: always
ai_summary: Defines provider, service management, runtime, traffic, backup, observability, evidence, and AI-operation plugin boundaries.
---

# Plugin Taxonomy and Boundaries

## Purpose

Sovrunn needs explicit plugin categories so lifecycle logic, infrastructure execution, runtime configuration, and governance decisions do not become mixed.

## Plugin Types

- Provider/Substrate Plugin
- PaaS Service Management Plane Plugin
- PaaS Service Runtime Plugin
- Traffic Management Plugin
- Backup Management Plugin
- Observability Plugin
- Security Plugin
- Compliance Evidence Plugin
- AI Operations Plugin

## Responsibility Split

```text
Sovrunn Core decides.
Policy engine evaluates.
PaaS Service Management Plane Plugin plans.
Provider/Substrate Plugin provisions capacity.
PaaS Service Runtime Plugin configures/checks runtime.
Operation records lifecycle.
Audit records every meaningful action.
AI explains and recommends, but does not bypass policy.
```

## Resources

- PluginType
- PluginManifest
- PluginProfile
- ProviderPluginProfile
- ServiceManagementPluginProfile
- ServiceRuntimePluginProfile
- PluginCapabilityScope
- PluginTrustProfile
- PluginCredentialPolicy
- PluginExecutionBoundary

## Phase 2 Scope

Metadata and taxonomy only. Runtime execution begins in Phase 3.

## Phase 2R Update

**ADH-2026-042**

Plugin taxonomy in Phase 2R includes a synthetic PluginExecution contract (FEATURE-0024). The fake execution harness is conformance machinery only — no provider or workload side effects.

Plugin/adapter roles, manifests, compatibility rules, and protected handles are governed per the canonical contract catalog. Plugins must not bypass final DecisionRecord, Operation, target authority, or SecretRef boundaries.

Canonical plugin boundary model: `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`.
