---
doc_type: demo
title: Phase 1 Demo Flow
status: draft
phase: 1
ai_load_priority: reference
ai_summary: Human-readable demo flow for Phase 1 resources.
---

# Phase 1 Demo Flow

Run after the API server is implemented and listening on `127.0.0.1:8080`.

```bash
chmod +x scripts/demo_phase1.sh
./scripts/demo_phase1.sh
```

Expected outcome:

```text
Organization created
OrganizationUnit created
Tenant created
Project created
ServiceClass created
ServicePlan created
Plugin created
Capability created
ServiceInstance created
ServiceBinding created
Operations listed
```
