---
doc_type: feature
id: FEATURE-0009
title: API Server Health and Readiness
status: draft
phase: 1
depends_on: [FEATURE-0001]
ai_load_priority: feature
ai_summary: Implements basic API server health and readiness endpoints.
---

# FEATURE-0009 API Server Health/Readiness
> **Phase 1 baseline note:** This document remains valid as the Phase 1 baseline. Phase 2 extends Sovrunn with reuse-first architecture, adapter boundaries, provider-neutral resource modeling, policy evaluation abstraction, decision/audit standards, plugin taxonomy, and governed placement. Do not treat this Phase 1 document as the complete Phase 2 scope.

## 1. Objective

Implement health, readiness, and basic server metadata endpoints.

Although this feature is listed late, implement a minimal health server at the beginning of coding and finalize it here.

## 2. Endpoints

| Method | Path | Response |
|---|---|---|
| GET | `/healthz` | `ok` |
| GET | `/readyz` | `ready` |
| GET | `/version` | JSON version metadata |

## 3. Version Response

```json
{
  "name": "sovrunn-api",
  "version": "0.1.0",
  "phase": "phase-1",
  "status": "development"
}
```

## 4. Readiness Rules

In Phase 1, readiness means:

- HTTP server started.
- Registry initialized.
- Required config loaded.
- No fatal initialization error.

## 5. Acceptance Criteria

- `/healthz` returns 200.
- `/readyz` returns 200 after registry initialization.
- `/version` returns version JSON.
- Tests cover endpoints.
