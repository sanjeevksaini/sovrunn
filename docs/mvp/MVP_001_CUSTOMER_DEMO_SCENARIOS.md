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
PlacementDecision selects a ResourcePool
AuditEvent is recorded
AI-readable explanation states why allowed
```

## Scenario 2: Denied Public Endpoint

Expected result:

```text
DENIED
SecurityProfile requires private endpoint
Suggested action: select private endpoint plan or approved pool
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
ServiceRuntimeProfile requires capability not present in ResourcePool
Suggested action: choose compatible ResourcePool or plan
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
