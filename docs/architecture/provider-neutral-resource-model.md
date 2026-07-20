---
doc_type: architecture
title: Provider-Neutral Resource Model
status: draft
phase: 2
ai_load_priority: always
ai_summary: Provider-neutral infrastructure model for capability-driven placement across local cloud, private cloud, hyperscalers, and hybrid substrates.
---

# Provider-Neutral Resource Model

## Purpose

Sovrunn must model heterogeneous provider infrastructure without hardcoding providers, clouds, or substrates into the core.

## Core Model

```text
Provider
  -> ProviderLocation / ProviderRegion
      -> ProviderDatacenter
          -> DatacenterFailureDomain
              -> IaaSStack
                  -> ResourcePool
                      -> ProviderCapability
```

## Boundaries

| Resource | Boundary |
|---|---|
| Provider | Operator/business boundary. |
| ProviderLocation / ProviderRegion | Sovereignty/geography boundary. |
| ProviderDatacenter | Physical site boundary. |
| DatacenterFailureDomain | Failure-isolation boundary inside a datacenter. |
| IaaSStack | Infrastructure implementation boundary. |
| ResourcePool | Placement boundary. |
| ProviderCapability | Compatibility boundary. |
| ProviderPlugin | Execution boundary. |

## Capability Status

ProviderCapability status may be:

```text
declared
validated
certified
degraded
disabled
```

## Rule

Placement must depend on capabilities and policy context, not on hardcoded provider names.
