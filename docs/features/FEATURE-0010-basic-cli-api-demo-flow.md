---
doc_type: feature
id: FEATURE-0010
title: Basic CLI/API Demo Flow
status: draft
phase: 1
depends_on: [FEATURE-0001, FEATURE-0002, FEATURE-0003, FEATURE-0004, FEATURE-0005, FEATURE-0006, FEATURE-0007, FEATURE-0008, FEATURE-0009]
ai_load_priority: feature
ai_summary: Provides the end-to-end demo flow for Phase 1.
---

# FEATURE-0010 Basic CLI/API Demo Flow

## 1. Objective

Create a repeatable demo that proves the Phase 1 resource model works end-to-end.

The demo may be implemented as curl scripts first. A real CLI can come later.

## 2. Demo Flow

```text
1. Start Sovrunn API.
2. Check health/readiness.
3. Create Organization nic.
4. Create OrganizationUnit ministry-health.
5. Create Tenant national-health-mission.
6. Create Project production.
7. Register ServiceClass datastore.postgresql.
8. Register ServicePlan postgres-small-ha.
9. Register Plugin postgres.dstoreops.basic.
10. Register Capability Provision.
11. Register Capability Bind.
12. Create ServiceInstance nhm-prod-postgres.
13. Create ServiceBinding nhm-app-postgres-binding.
14. List Operations.
15. List tenant/project ServiceInstances.
```

## 3. Demo Script

Provide:

```text
scripts/demo_phase1.sh
```

The script must be idempotent where practical.

## 4. CLI Shape Later

Future CLI examples:

```bash
sovrunn org create nic --display-name "National Informatics Centre"
sovrunn ou create ministry-health --org nic
sovrunn tenant create national-health-mission --org nic --ou ministry-health
sovrunn project create production --tenant national-health-mission
sovrunn service instance create nhm-prod-postgres --class datastore.postgresql --plan postgres-small-ha
```

## 5. Acceptance Criteria

- Demo script completes successfully against local API.
- Created resources can be listed.
- Operation records are visible.
- Demo can be included in README.
