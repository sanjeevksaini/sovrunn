---
doc_type: requirements
feature_id: FEATURE-0010
title: Basic CLI/API Demo Flow
status: draft
phase: 1
depends_on: [FEATURE-0001, FEATURE-0002, FEATURE-0003, FEATURE-0004, FEATURE-0005, FEATURE-0006, FEATURE-0007, FEATURE-0008, FEATURE-0009]
---

# FEATURE-0010 Basic CLI/API Demo Flow — Requirements

## 1. Introduction

FEATURE-0010 delivers a repeatable, end-to-end demo that proves the Phase 1
resource model works correctly across all Sovrunn platform resources.

The demo is implemented as a bash script (`scripts/demo_phase1.sh`) that
exercises the REST API using curl. A real CLI binary is out of scope for
Phase 1.

This feature depends on all prior Phase 1 features (FEATURE-0001 through
FEATURE-0009) being complete and stable.

### 1.1 Context

An existing `scripts/demo_phase1.sh` already exists with curl-based calls.
This requirements document defines the full acceptance surface for the demo
feature, including script behavior, output validation, error handling,
idempotency expectations, and documentation integration.

## 2. Glossary

No new concepts are introduced. All terms used in this feature are defined
in `docs/glossary.md`. Key terms exercised:

| Term | Relevance |
|---|---|
| Organization | Created as governance root |
| OrganizationUnit | Created under Organization |
| Tenant | Created under OrganizationUnit |
| Project | Created under Tenant |
| ServiceClass | Registered in catalog |
| ServicePlan | Registered under ServiceClass |
| Plugin | Registered in plugin registry |
| Capability | Declared for plugin |
| ServiceInstance | Provisioned in project |
| ServiceBinding | Bound to ServiceInstance |
| Operation | Audit trail for mutations |

## 3. User Stories

### US-1: Platform operator runs the full demo

As a platform operator, I want to run `make demo` and see the complete
Phase 1 resource lifecycle exercised end-to-end so that I can verify the
platform core works correctly.

### US-2: Developer validates API after changes

As a developer, I want the demo script to fail fast with clear output on
any API error so that I can identify regressions immediately after code
changes.

### US-3: Evaluator understands platform capabilities

As a prospective evaluator, I want the demo script to print human-readable
step descriptions and results so that I can understand what Sovrunn offers
without reading source code.

### US-4: CI pipeline uses demo as smoke test

As a CI maintainer, I want the demo script to exit with non-zero status on
failure so that it can be used as a basic smoke test after build.

### US-5: README references the demo

As a new contributor, I want the README to reference `make demo` so that I
have a clear quick-start path to see the platform in action.

## 4. Acceptance Criteria

### AC-1: Script executes full resource lifecycle

The demo script must perform these steps in order:

1. Check `/healthz` returns 200.
2. Check `/readyz` returns 200.
3. Create Organization `nic`.
4. Create OrganizationUnit `ministry-health` under `nic`.
5. Create Tenant `national-health-mission` under `ministry-health`.
6. Create Project `production` under `national-health-mission`.
7. Register ServiceClass `datastore.postgresql`.
8. Register ServicePlan `postgres-small-ha` for that ServiceClass.
9. Register Plugin `postgres.dstoreops.basic`.
10. Register Capability `postgres-basic-provision` for that Plugin.
11. Create ServiceInstance `nhm-prod-postgres`.
12. Create ServiceBinding `nhm-app-postgres-binding`.
13. List Operations and verify at least one operation exists.
14. GET the created ServiceInstance by name and verify response.
15. GET the created ServiceBinding by name and verify response.

### AC-2: Script exits non-zero on failure

If any API call returns an unexpected HTTP status, the script must exit
immediately with a non-zero exit code.

### AC-3: Script prints step-level output

Each step must print a human-readable label (e.g., "Creating Organization...")
before the curl call, and the response body or a success indicator after.

### AC-4: Script is runnable via `make demo`

The Makefile `demo` target must invoke the script. The target already exists.

### AC-5: Script uses configurable base URL

The base URL must default to `http://127.0.0.1:8080` and be overridable via
`BASE_URL` environment variable.

### AC-6: Operations are recorded

After all mutations, `GET /v1/operations` must return a non-empty list. The
demo must print the operation count or list.

### AC-7: Script is portable bash

The script must work on macOS and Linux with bash 4+ and curl installed.
No other runtime dependencies.

## 5. Non-Goals

- Real CLI binary (`sovrunn` command) — deferred to future phase.
- Interactive terminal UI or colored output requiring ncurses.
- Persistent state across demo runs (in-memory registry resets on restart).
- Load testing or performance benchmarking.
- Authentication or authorization headers (Phase 1 has no auth).
- Demo against remote/production clusters.
- Docker-based demo orchestration.
- Automated server start/stop within the script (user starts server separately).
- Idempotency via update-or-create logic (409 on re-run is acceptable).
- Testing of DELETE endpoints in the demo flow.
- Multi-organization or multi-tenant demo scenarios.
- ServiceOps plugin execution or real provisioning.

## 6. Edge Cases

| # | Edge Case | Expected Behavior |
|---|---|---|
| 1 | API server not running | Script fails at health check with clear error and non-zero exit |
| 2 | Server returns non-JSON error (e.g., 502 from proxy) | Script fails with curl error; `set -e` halts execution |
| 3 | Re-running demo without server restart | 409 Conflict on creates; script fails fast, user must restart server |
| 4 | Network timeout | curl fails with timeout error; script exits non-zero |
| 5 | Partial demo completion (ctrl+c) | No cleanup needed; in-memory state is ephemeral |
| 6 | Missing curl binary | Script fails immediately with command-not-found |
| 7 | BASE_URL set to wrong port | Health check fails; script exits before mutations |
| 8 | Unexpected JSON shape in response | Not validated in Phase 1 script; future enhancement |

## 7. Security and Privacy Requirements

- The demo script must not contain real credentials, tokens, or secrets.
- All resource payloads must use example/fictional data only.
- The script must not send requests to external endpoints.
- The script must not store output to files that could leak into version
  control (all output goes to stdout/stderr only).
- No Authorization headers are included (Phase 1 has no auth layer).
- SecretRef values in ServiceBinding responses are not validated in the
  demo flow but must not be printed if they were real.

## 8. Compatibility with Completed Phase 1 Features

| Feature | Compatibility Requirement |
|---|---|
| FEATURE-0001 Organization | Demo creates Organization with valid spec including `sovereignLocations` |
| FEATURE-0002 OrganizationUnit | Demo creates OU with valid `organizationRef` |
| FEATURE-0003 Tenant | Demo creates Tenant with valid `organizationRef`, `organizationUnitRef`, and `isolationProfile` |
| FEATURE-0004 Project | Demo creates Project with valid `tenantRef`, `organizationRef`, `organizationUnitRef`, and `environmentType` |
| FEATURE-0005 Operation | Demo verifies operations are emitted for mutating calls |
| FEATURE-0006 ServiceClass/Plan | Demo registers both with valid cross-references |
| FEATURE-0007 Plugin/Capability | Demo registers Plugin with `serviceClassRefs` and Capability with `pluginRef` |
| FEATURE-0008 ServiceInstance/Binding | Demo creates instance with all required refs and binding with `serviceInstanceRef` |
| FEATURE-0009 Health/Readiness | Demo checks `/healthz` and `/readyz` before any mutations |

The demo payloads must match the exact field names, validation rules, and
reference constraints enforced by the existing API handlers. If existing
validation rejects a demo payload, the payload in the script must be fixed
to conform — not the validation.

## 9. Design Questions to Resolve in design.md

1. **Verification depth**: Should the script validate response JSON fields
   (e.g., check that returned `metadata.name` matches), or just check HTTP
   status codes? Tradeoff: robustness vs. script complexity.

2. **Second Capability registration**: The FEATURE-0010 spec mentions
   registering two capabilities (Provision and Bind). Should both be
   demonstrated, or is one sufficient for the demo?

3. **List filtering**: Should the demo call list endpoints with query
   filters (e.g., list service instances by project) or just unfiltered
   lists? Depends on whether filter support exists in FEATURE-0008.

4. **Output format**: Should the script pretty-print JSON responses (requires
   `jq` dependency) or print raw curl output? Tradeoff: readability vs.
   zero-dependency portability.

5. **Server lifecycle**: Should the demo script start the server, wait for
   readiness, run the flow, and stop the server? Or should server management
   remain manual? Tradeoff: convenience vs. simplicity and portability.

6. **Shebang portability**: The existing script uses
   `#!/opt/homebrew/bin/bash`. Should this be changed to `#!/usr/bin/env bash`
   for Linux compatibility?

7. **Cleanup step**: Should the script include an optional cleanup phase
   (delete all created resources) to support re-runs without server restart?

8. **README integration**: Where exactly should the demo instructions be
   placed in README.md — under a new section, or inline with existing
   "Local Validation" section?
