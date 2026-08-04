---
doc_type: mvp_scenarios
title: MVP-001 Customer Demo Scenarios
status: draft
phase: 3
ai_load_priority: high
ai_summary: Customer demo scenarios for the first governed PostgreSQL PaaS MVP.
---

# MVP-001 Customer Demo Scenarios

## Scenario 1: Allowed PostgreSQL Placement

Request:

```text
PostgreSQL basic/small
private endpoint
approved India location
regulated profile
```

Expected result:

```text
ALLOWED
PlacementDecision selects a qualified ExecutionTarget
AuditEvent is recorded
AI-readable explanation states why allowed
```

## Scenario 2: Denied Public Endpoint

Expected result:

```text
DENIED
SecurityProfile requires private endpoint
Suggested action: select private endpoint plan or request exception
```

## Scenario 3: Denied Location

Expected result:

```text
DENIED
DataPlacementPolicy allows only approved locations
Suggested action: choose approved provider/location
```

## Scenario 4: Denied Capability Mismatch

Expected result:

```text
DENIED
ServiceRequirementSet requires capability not satisfied by any qualified ExecutionTarget
Suggested action: choose compatible plan or request target qualification review
```

## Scenario 5: Provisioned Service

Expected result:

```text
ALLOWED
Operation created
PostgreSQL runtime delegated to reused operator/Helm wrapper
ServiceBinding created with SecretRef/CredentialRef
AuditEvent recorded
```

## Phase 2R Canonical Context

**ADH-2026-042**

Demo scenarios use canonical model terminology:

- ServiceOffering/ServicePlan replaces ServiceClass/ServicePlan.
- Qualified ExecutionTarget replaces ResourcePool.
- EffectiveGovernanceContext replaces EffectivePolicyContext.
- ServicePlacement is the safe customer projection.
- DecisionRecord profiles handle sovereignty and placement conclusions.
- ServiceBinding is SecretRef-only, per-consumer, separately revocable.