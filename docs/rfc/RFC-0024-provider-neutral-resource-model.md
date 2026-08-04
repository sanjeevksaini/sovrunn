---
doc_type: rfc
title: RFC-0024 Provider-Neutral Resource Model
status: approved
phase: 2
ai_load_priority: high
ai_summary: RFC for the five-resource provider-neutral topology model controlled by ADH-2026-018; placement inventory, capability, adapters, connectivity, and execution remain outside FEATURE-0014.
---

# RFC-0024: Provider-Neutral Resource Model

Controlling architecture: `docs/architecture/provider-neutral-resource-model.md`.

Controlling handoff: `docs/reviews/architecture-decision-handoffs/ADH-2026-018-feature-0014-provider-neutral-resource-model.md`.

## Decision

FEATURE-0014 defines exactly Provider, ProviderLocation, ProviderDatacenter, DatacenterFailureDomain, and InfrastructureStack. It reuses the existing Organization resource as optional supply owner and FEATURE-0012 resource grammar unchanged.

FEATURE-0014 does not define ResourcePool, ProviderCapability, connectivity, adapters, provider integration, decision/audit behavior, or execution. Physical containment implies neither network connectivity nor isolation.

## Phase 2R Migration Note

**ADH-2026-042 / DEC-0042:** The ResourcePool/ProviderCapability placement spine originally defined in this RFC is superseded. Qualified ExecutionTarget is the canonical realization boundary per DEC-0042. FEATURE-0014's five-resource topology (Provider, ProviderLocation, ProviderDatacenter, DatacenterFailureDomain, InfrastructureStack) remains completed history. Provider migrates to CloudPlatform/CloudProvider classification per DEC-0037. ProviderLocation becomes registered HostingLocation fact per DEC-0041. The resource model defined by ADH-2026-018 remains intact as implementation history; the active placement and capability contracts are now governed by the canonical model in `docs/architecture/canonical/sovrunn-finalized-data-model.md`.
