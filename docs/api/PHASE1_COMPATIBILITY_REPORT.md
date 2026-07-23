---
doc_type: compatibility_report
title: Phase 1 Compatibility Report
status: draft
phase: 1
feature: FEATURE-0012
ai_load_priority: phase
ai_summary: Records Phase 1 API/resource conforming behavior, explicit exceptions, and migration candidates under the FEATURE-0012 grammar without triggering a rewrite.
---

# Phase 1 Compatibility Report

> **Governing decisions:** D-13, D-09; requirements F12-COMPAT-001, F12-COMPAT-002, F12-COMPAT-003, F12-COMPAT-005.
> **Source baseline:** `docs/api/API_CONTRACT_PHASE1.md`, Phase 1 resource types in `internal/resources`, handlers in `internal/api`, demo flow in `docs/demo/PHASE1_DEMO_FLOW.md` / `scripts/demo_phase1.sh`.

## Purpose

FEATURE-0012 introduces the Sovrunn-owned API/resource grammar (profiles, metadata, scope, references, conditions, Problem Details, validation pipeline, and conformance fixtures). This report documents how completed Phase 1 contracts relate to that grammar.

This report:

- covers every Phase 1 resource and endpoint required by F12-COMPAT-001;
- records conforming behavior, explicit exceptions, and migration candidates per contract (F12-COMPAT-002);
- states that Phase 1 routes and wire behavior are retained unchanged until a separately approved migration (F12-COMPAT-003);
- does **not** trigger a wholesale rewrite of Phase 1 APIs, handlers, or demo flow.

## Rewrite posture (F12-COMPAT-003)

| Decision | Status |
| --- | --- |
| Phase 1 `/v1/...` resource routes | Retained unchanged |
| Phase 1 `/healthz`, `/readyz`, `/version` | Retained unchanged |
| Phase 1 `{"error": ...}` envelope | Retained unchanged (exception recorded) |
| Phase 1 demo flow (`scripts/demo_phase1.sh`) | Retained unchanged |
| FEATURE-0012 runtime route registration | None (grammar/conformance only) |
| Wholesale Phase 1 rewrite | Not triggered |

New contracts that adopt FEATURE-0012 MUST use `/apis/<group>/<version>/<plural-kebab>` and the shared grammar. Existing Phase 1 contracts coexist during migration; any breaking change requires the versioning/migration/approval rules in F12-EVOLVE-* and F12-COMPAT-004.

## Coverage inventory (F12-COMPAT-001)

The following Phase 1 contracts are covered by this report:

| Contract ID | Kind / endpoint family | Phase 1 routes (unchanged) |
| --- | --- | --- |
| Organization | Organization | `/v1/organizations` |
| OrganizationUnit | OrganizationUnit | `/v1/organization-units` |
| Tenant | Tenant | `/v1/tenants` |
| Project | Project | `/v1/projects` |
| Operation | Operation | `/v1/operations` (GET list/get only) |
| ServiceClass | ServiceClass | `/v1/service-classes` |
| ServicePlan | ServicePlan | `/v1/service-plans` |
| Plugin | Plugin | `/v1/plugins` |
| Capability | Capability | `/v1/capabilities` |
| ServiceInstance | ServiceInstance | `/v1/service-instances` |
| ServiceBinding | ServiceBinding | `/v1/service-bindings` |
| health/readiness | health/readiness (+ `/version`) | `/healthz`, `/readyz`, `/version` |
| demo-flow | demo-flow | `scripts/demo_phase1.sh` (exercises Phase 1 routes) |

## Shared Phase 1 baseline (applies to all resource contracts)

### Conforming behavior (generalized by FEATURE-0012)

- **Resource shape:** `apiVersion` / `kind` / `metadata` / `spec` / `status` (F12-COMPAT-005).
- **Status ownership:** `status` is system-owned; customer create/replace must not author status.
- **Stable error codes:** Phase 1 `resources.ErrorCode` values (`VALIDATION_FAILED`, `RESOURCE_NOT_FOUND`, `RESOURCE_ALREADY_EXISTS`, `DELETE_BLOCKED`, `METHOD_NOT_ALLOWED`, `INTERNAL_ERROR`) remain a compatible subset of the FEATURE-0012 stable-code approach.
- **In-memory registry:** Phase 1 storage remains in-memory; FEATURE-0012 does not introduce persistent storage.
- **Request correlation:** Phase 1 propagates request IDs (for example `X-Sovrunn-Request-ID`); FEATURE-0012 Problem Details generalize this via `requestId`.
- **No secrets in wire errors:** error bodies carry codes/messages/field hints, not credentials or connection strings.
- **Provider neutrality of core customer routes:** Phase 1 organization/tenant/project/instance routes do not expose provider SDK types on the customer contract surface.

### Explicit exceptions (shared)

| Exception ID | Description | Why retained |
| --- | --- | --- |
| EX-P1-ROUTE | Routes use `/v1/<plural-kebab>` rather than `/apis/<group>/<version>/<plural-kebab>`. | F12-COMPAT-003; D-09 coexistence. |
| EX-P1-ERROR-ENVELOPE | Errors use `{"error":{"code","message","field","details"}}` rather than RFC 9457 Problem Details (`type`/`title`/`status`/`detail`/`instance`/`code`/`requestId`/`violations[]`). | Documented migration candidate; no silent reinterpretation. |
| EX-P1-FIELD-PATH | Validation field paths are dotted JSON paths (for example `metadata.name`), not RFC 6901 JSON Pointers (`/metadata/name`). | Compatible intent; pointer form is the FEATURE-0012 norm for new contracts. |
| EX-P1-METADATA | Phase 1 `Metadata` is name/displayName/labels/annotations only (no uid, resourceVersion, generation, scopeRef, ownerRef, timestamps). | Phase 1 identity was name-keyed; FEATURE-0012 ObjectMeta is the target for adopters. |
| EX-P1-STATUS | Status is typically `{phase,message}` without `apicond.Condition` arrays / lastTransitionTime semantics. | Conditions are the FEATURE-0012 status grammar for new/adopted contracts. |
| EX-P1-REFS | Parent and catalog links are often plain name strings in `spec` (for example `tenantName`, `organizationRef`) rather than `TypedRef` / `ScopeRef` / `OwnerRef`. | Typed references are mandatory for new FEATURE-0012 contracts; Phase 1 string refs remain as-is. |
| EX-P1-LIST | List responses use a simple `{"items":[...]}` shape without `ListEnvelope` TypeMeta promotion or opaque `page.nextPageToken`. | Pagination grammar is for adopting APIs; Phase 1 lists unchanged. |
| EX-P1-CONCURRENCY | No `If-Match` / `resourceVersion` optimistic concurrency on replace. | F12-UPDATE-002 applies to adopting contracts; Phase 1 replace behavior retained. |
| EX-P1-DECODE | Phase 1 decode/validation is handler-local and does not yet compose the nine-layer FEATURE-0012 pipeline, strict YAML subset, or schema registry. | FEATURE-0012 supplies reusable grammar; Phase 1 handlers are not rewritten here. |

### Shared migration candidates

1. Introduce a versioned `/apis/<group>/<version>/...` surface (or dual-publish) with an approved migration plan; keep `/v1/...` until clients migrate.
2. Map `APIErrorEnvelope` → RFC 9457 Problem Details while preserving stable `code` values and adding JSON Pointer `violations[].field`.
3. Evolve `Metadata` toward FEATURE-0012 `ObjectMeta` (uid, resourceVersion, scopeRef/ownerRef where applicable) behind a new apiVersion.
4. Replace string parent/catalog refs with constrained `TypedRef` / `ScopeRef` where the contract crosses trust or scope boundaries.
5. Adopt `Condition` status grammar where multi-aspect observed state is required; keep phase as a summary field if useful.
6. Adopt `ListEnvelope` + opaque page tokens for list endpoints that need pagination.
7. Add `If-Match` / `resourceVersion` for protected replace once ObjectMeta carries resourceVersion.
8. Compose FEATURE-0012 decode/validate helpers at adoption time without silently changing Phase 1 request/response semantics.

---

## Per-contract records

Each subsection below is a required coverage entry. Headings use the exact F12-COMPAT-001 contract identifiers.

### Organization

- **Phase 1 routes (unchanged):** `POST/GET/PUT/DELETE /v1/organizations`, `GET /v1/organizations/{name}`
- **Conforming behavior:** metadata/spec/status shape; system-owned status; DNS-label name validation; stable error codes; in-memory registry; request ID on responses.
- **Explicit exceptions:** EX-P1-ROUTE, EX-P1-ERROR-ENVELOPE, EX-P1-FIELD-PATH, EX-P1-METADATA, EX-P1-STATUS, EX-P1-LIST, EX-P1-CONCURRENCY, EX-P1-DECODE. Organization is platform-root; Phase 1 has no `scopeRef` (canonical platform nil-scope is a FEATURE-0012 concept not emitted on Phase 1 Organization).
- **Migration candidates:** Adopt FEATURE-0012 Organization (or equivalent) schema under `/apis/...` with ObjectMeta, Problem Details, and optional Condition status; retain `/v1/organizations` until approved cutover.

### OrganizationUnit

- **Phase 1 routes (unchanged):** `POST/GET/PUT/DELETE /v1/organization-units`, `GET /v1/organization-units/{name}`
- **Conforming behavior:** metadata/spec/status; parent Organization referenced; status system-owned; delete-blocked / missing-reference behavior via stable codes.
- **Explicit exceptions:** Shared EX-P1-* set. Parent link is a string name in spec rather than `TypedRef`/`OwnerRef` (EX-P1-REFS). Governance scope is implicit via organizationName, not `metadata.scopeRef`.
- **Migration candidates:** Typed parent/scope references; ObjectMeta; Problem Details; route form `/apis/.../organization-units` under approved versioning.

### Tenant

- **Phase 1 routes (unchanged):** `POST/GET/PUT/DELETE /v1/tenants`, `GET /v1/tenants/{name}`
- **Conforming behavior:** metadata/spec/status; hierarchy refs to Organization/OrganizationUnit; status ownership; stable validation/conflict codes.
- **Explicit exceptions:** Shared EX-P1-* set; string hierarchy refs (EX-P1-REFS); no FEATURE-0012 `ScopeRef` identity tuple on the wire.
- **Migration candidates:** Express tenant governance scope via `scopeRef`/`TypedRef`; migrate error envelope and list envelope with a new apiVersion.

### Project

- **Phase 1 routes (unchanged):** `POST/GET/PUT/DELETE /v1/projects`, `GET /v1/projects/{name}`
- **Conforming behavior:** metadata/spec/status; immutable parent names in spec (`organizationName`, `organizationUnitName`, `tenantName`); status system-owned. Aligns in intent with FEATURE-0012 ManagedResource / customer-facing Project fixture family (schema fit is proven for the *new* grammar Project, not a silent reinterpretation of Phase 1 Project).
- **Explicit exceptions:** Shared EX-P1-* set; Phase 1 Project must not be treated as already identical to `api/schemas/project.json` without an explicit versioned migration (Matrix D “Phase 1 resource migrated” requires explicit group/version maturity on the *new* contract).
- **Migration candidates:** Dual-publish or replace under `/apis/core.sovrunn.io/<version>/projects` using the FEATURE-0012 Project schema (Tenant scopeRef, ObjectMeta, field-policy annotations); keep `/v1/projects` until clients migrate.

### Operation

- **Phase 1 routes (unchanged):** `GET /v1/operations`, `GET /v1/operations/{name}` (create is server-emitted; no public POST in Phase 1 contract)
- **Conforming behavior:** metadata/spec/status lifecycle trace; records action type and non-secret resource references; carries `requestId` in spec when available; status phases include Succeeded/Failed (and reserved Pending/Running); immutable emission pattern.
- **Explicit exceptions:** Shared EX-P1-* set. Phase 1 Operation uses flat string resource fields (`resourceKind`, `resourceName`, optional hierarchy name fields) rather than FEATURE-0012 `targetRef` + `scopeRef` (D-17). Phase 1 Operation is not the six-scope LongRunningOperation schema; do not silently reinterpret Phase 1 Operation under D-17.
- **Migration candidates:** Adopt FEATURE-0012 Operation schema (`targetRef`, canonical `scopeRef` matching target governance scope, optional `ownerRef` without replacing scope) on a new versioned route family; map emission sites to typed refs; retain GET `/v1/operations` until cutover.

### ServiceClass

- **Phase 1 routes (unchanged):** `POST/GET/PUT/DELETE /v1/service-classes`, `GET /v1/service-classes/{name}`
- **Conforming behavior:** catalog definition only (no provisioning); metadata/spec/status; global name identity; status system-owned; no secrets in catalog fields.
- **Explicit exceptions:** Shared EX-P1-* set. Global/platform catalog identity without FEATURE-0012 profile/boundary annotations on the Phase 1 wire type. Optional `provider` string in spec is catalog metadata, not a provider-native SDK type embedded in core grammar packages.
- **Migration candidates:** Versioned catalog schema with FEATURE-0012 annotations (profile/boundary/allowed-scopes/field-policy); TypedRef for default plan linkage; Problem Details.

### ServicePlan

- **Phase 1 routes (unchanged):** `POST/GET/PUT/DELETE /v1/service-plans`, `GET /v1/service-plans/{name}`
- **Conforming behavior:** catalog plan under a ServiceClass; metadata/spec/status; reference integrity to ServiceClass; status ownership.
- **Explicit exceptions:** Shared EX-P1-* set; string `serviceClass` / class ref rather than TypedRef (EX-P1-REFS).
- **Migration candidates:** Typed ServiceClass reference; FEATURE-0012 schema annotations; versioned `/apis/...` route.

### Plugin

- **Phase 1 routes (unchanged):** `POST/GET/PUT/DELETE /v1/plugins`, `GET /v1/plugins/{name}`
- **Conforming behavior:** registry declaration only (no plugin execution in Phase 1); metadata/spec/status; no secrets in plugin fields; ServiceClassRefs as name list.
- **Explicit exceptions:** Shared EX-P1-* set. Phase 1 Plugin is not the FEATURE-0012 `plugin-definition.json` VersionedDefinition contract; coexistence only—no silent schema equivalence claim.
- **Migration candidates:** Map toward FEATURE-0012 PluginDefinition (or successor) under plugin-facing boundary with typed refs and field policies; keep `/v1/plugins` until approved migration.

### Capability

- **Phase 1 routes (unchanged):** `POST/GET/DELETE /v1/capabilities`, `GET /v1/capabilities/{name}` (Phase 1 contract has no PUT for capabilities)
- **Conforming behavior:** capability registration linked to plugins/service classes; metadata/spec/status; registry-only (no execution).
- **Explicit exceptions:** Shared EX-P1-* set; string refs; Phase 1 update surface differs (no replace) from full CRUD normative set on adopting APIs.
- **Migration candidates:** Typed capability/plugin refs; decide whether Capability remains a first-class external contract or folds into PluginDefinition/ProviderCapability in later features under approved versioning—without rewriting Phase 1 now.

### ServiceInstance

- **Phase 1 routes (unchanged):** `POST/GET/PUT/DELETE /v1/service-instances`, `GET /v1/service-instances/{name}`
- **Conforming behavior:** desired-state request only (no real provisioning); metadata/spec/status; hierarchy and catalog string refs; parameters must not hold secrets; status system-owned.
- **Explicit exceptions:** Shared EX-P1-* set; EX-P1-REFS for organization/tenant/project/class/plan refs; no FEATURE-0012 Operation target-scope equality on the Phase 1 instance itself.
- **Migration candidates:** TypedRef/ScopeRef for hierarchy and catalog; adopt provisioning via FEATURE-0012 Operation grammar in later features without changing Phase 1 instance routes here.

### ServiceBinding

- **Phase 1 routes (unchanged):** `POST/GET/DELETE /v1/service-bindings`, `GET /v1/service-bindings/{name}` (no PUT in Phase 1 contract)
- **Conforming behavior:** binding desired state only; metadata/spec/status; references ServiceInstance; no secret material in binding fields; status system-owned.
- **Explicit exceptions:** Shared EX-P1-* set; string refs; no replace route in Phase 1.
- **Migration candidates:** Typed ServiceInstance reference; versioned binding schema; optional replace semantics only under approved API evolution.

### health/readiness

- **Phase 1 routes (unchanged):** `GET /healthz` → `ok`; `GET /readyz` → `ready` (503 when not ready); `GET /version` → JSON server metadata
- **Conforming behavior:** liveness/readiness separation; no resource grammar required; suitable for probes; does not leak secrets.
- **Explicit exceptions:** Unversioned operational endpoints outside `/apis/<group>/<version>/...` (intentional ops surface). Not subject to resource schema/profile annotations. `/version` is informational JSON, not a FEATURE-0012 resource kind.
- **Migration candidates:** Retain unversioned health/ready probes. Optionally add a versioned platform info resource later without removing `/healthz`/`/readyz`. Do not force Probe endpoints through ValidateRoute’s resource-collection pattern.

### demo-flow

- **Phase 1 surface (unchanged):** `scripts/demo_phase1.sh` and `docs/demo/PHASE1_DEMO_FLOW.md` exercise the Phase 1 `/v1/...` and health endpoints end-to-end (Organization → … → ServiceBinding, then Operations list).
- **Conforming behavior:** Demonstrates hierarchy, catalog, plugin/capability registration, instance/binding create, and operation listing against the retained Phase 1 contract; validates live server behavior without FEATURE-0012 runtime routes.
- **Explicit exceptions:** Demo is a scripted client of Phase 1 compatibility APIs, not a FEATURE-0012 conformance fixture runner. It does not exercise RFC 9457, `/apis/...` routes, or canonical schema fixtures.
- **Migration candidates:** Add a separate conformance/demo path for FEATURE-0012 fixtures when adopting APIs exist; keep `demo_phase1.sh` as the Phase 1 compatibility smoke path until Phase 1 routes are formally deprecated.

---

## Consistency with FEATURE-0012 (F12-COMPAT-005)

| Phase 1 concept | FEATURE-0012 generalization | Compatibility stance |
| --- | --- | --- |
| metadata/spec/status | Same shape; richer ObjectMeta/Condition/TypedRef | Compatible foundation |
| In-memory registry | Storage-replaceable; still no DB in this feature | Compatible |
| `resources.ErrorCode` | `apiproblem` stable codes + Problem Details | Codes consistent; envelope excepted (EX-P1-ERROR-ENVELOPE) |
| Request ID header | `requestId` on Problem + structured logs | Compatible correlation model |
| Name-keyed resources | uid + name + scope identity | Migration candidate, not silent change |

## Observability and audit notes (compatibility)

- **Preserved:** Phase 1 request ID propagation and structured, secret-free error surfaces remain valid.
- **Not changed by this report:** No new runtime audit emitters, log fields, or handlers are introduced by FEATURE-0012 task 15.1.
- **Adopter expectation:** When Phase 1 APIs later adopt FEATURE-0012 helpers, they MUST continue to avoid logging secrets, credentials, tokens, private keys, or connection strings; request/operation correlation fields remain preferred over embedding sensitive payloads.

## Non-goals confirmed

This compatibility report does **not**:

- rewrite Phase 1 routes or handlers;
- register FEATURE-0012 HTTP routes;
- claim Phase 1 types are byte-identical to `api/schemas/*.json`;
- implement task 15.2 coverage assertion code (separate task);
- approve residual risk or write human approval tokens.

## Traceability

| Requirement | Evidence in this document |
| --- | --- |
| F12-COMPAT-001 | Coverage inventory + per-contract sections for all listed Phase 1 resources/endpoints |
| F12-COMPAT-002 | Conforming behavior, explicit exceptions, migration candidates per contract; no rewrite |
| F12-COMPAT-003 | Rewrite posture table; routes retained unchanged |
| F12-COMPAT-005 | Consistency table mapping Phase 1 concepts to FEATURE-0012 generalizations |
| D-13 | This file (`docs/api/PHASE1_COMPATIBILITY_REPORT.md`) |
| D-09 | Route-form coexistence notes under Organization…ServiceBinding and health/readiness |
