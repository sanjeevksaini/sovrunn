# FEATURE-0015 — Canonical Cloud Model Foundation: Design

## 1. Identity, stage, inputs, and closed boundary

| Field | Value |
|-------|-------|
| Feature | FEATURE-0015 — Canonical Cloud Model Foundation |
| Spec | `.kiro/specs/canonical-cloud-model-and-alpha-migration/` |
| Stage | Design |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase | 2R |
| Architecture authority | `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`; `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md` |
| Approved requirements | `.kiro/specs/canonical-cloud-model-and-alpha-migration/requirements.md` |
| Controlling Decisions | DEC-0037, DEC-0041, DEC-0042, DEC-0054, DEC-0059 |
| Controlling Handoffs | ADH-2026-020, ADH-2026-024, ADH-2026-025, ADH-2026-037, ADH-2026-042, ADH-2026-043 (non-migration portions), ADH-2026-045, ADH-2026-046, ADH-2026-047, ADH-2026-048, ADH-2026-049, ADH-2026-050, ADH-2026-051, ADH-2026-052, ADH-2026-053, ADH-2026-054, ADH-2026-055, ADH-2026-056, ADH-2026-057 |
| Owned resources | CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, InfrastructureStack |
| Module path | `github.com/sanjeevksaini/sovrunn` (Go 1.22) |

### 1.1 Purpose of this stage

This design translates the approved FEATURE-0015 requirements into an implementable representation. It selects only the deterministic mechanics that the architecture and requirements explicitly delegate (architecture §2.2, §7.12, §9; requirements §9): registry data structures, handler wiring and `http.ServeMux` registration, deterministic validation plumbing, bootstrap-grant resolution internals, idempotency reservation-state mechanics, and audit/mutation coordination. It introduces no product semantics, transfers no ownership, and designs no mechanism owned by an excluded feature.

Authority precedence is honored exactly: active baseline > canonical model/catalog > accepted DEC/ADH > VS-000 charter > contract specification > registry YAML > approved requirements > this design. Where any tension appeared, this design resolves it by reference to the higher authority and never by interpretation.

### 1.2 Closed boundary (inputs consumed by reference)

- **FEATURE-0011** — the Reuse Assessment Standard governance contract is consumed by reference; §7.4 records the reuse disposition for the retained FEATURE-0014 alpha assets. No FEATURE-0011 contract is redefined.
- **FEATURE-0012** — Problem Details (`internal/apiproblem`), ObjectMeta/TypedRef/ScopeRef/ScopeKind (`internal/apimeta`), Condition (`internal/apicond`), reference constraints (`internal/apiref`), and the layer 5–7 validation pipeline (`internal/apivalid`) are consumed by exact reference. Registry grammar `VS0-SCHEMA-001` (ObjectMeta/scope vocabulary), `VS0-SCHEMA-002` (TypedRef), `VS0-SCHEMA-003` (Condition), and `VS0-SCHEMA-004` (Problem) are not forked, renamed, or specialized.
- **FEATURE-0013** — DecisionRecord and AuditEvent (`internal/decision`; `VS0-SCHEMA-007`) are consumed by exact reference for audit evidence. FEATURE-0015 appends redacted AuditEvents through the FEATURE-0013 writer (`VS0-WRITER-011`); it does not redefine the AuditEvent envelope or linkage.
- **FEATURE-0014** — the alpha Provider/ProviderLocation/ProviderDatacenter/DatacenterFailureDomain/InfrastructureStack model in `internal/resources` is a retained repository asset (DEC-0059; ADH-2026-045). It is not runtime-migrated, not imported at runtime, and its types are not reused as canonical types. The explicit stale-concept exclusions are in §10.1.

### 1.3 Excluded adjacent-feature mechanisms (not designed here)

ExecutionTarget and its lifecycle (FEATURE-0016; `VS0-SCHEMA-015`, `VS0-SCHEMA-016`, `VS0-SCHEMA-017`), policy evaluation (FEATURE-0017), identity governance (FEATURE-0018), sovereignty (FEATURE-0019), governance context (FEATURE-0020), enrollment/entitlement/quota and provider-selection intent (FEATURE-0021), the service catalog and ServiceRegion (FEATURE-0022), placement/decision profiles (FEATURE-0023), plugin execution (FEATURE-0024), operation-context projection (FEATURE-0025), and the integration demo (FEATURE-0026) are visible in this design only as negative boundaries (§10). No field, route, status, writer, or state owned by those features is introduced, initialized, defaulted, validated, persisted, or exposed.

---

## 2. Resolved design decisions

Each decision below is a mechanic delegated by the approved authorities. None reopens a semantic or ownership question.

- **DD-01 — Dedicated canonical package, no clash with retained alpha assets.** The retained FEATURE-0014 alpha model already occupies generic type names (including an alpha `InfrastructureStack`) in `internal/resources`. To keep the canonical model self-contained and avoid symbol collision with the retained asset, FEATURE-0015 canonical resource types, in-memory store, deterministic validation, lifecycle, idempotency, and audit coordination live under a new `internal/cloudmodel` subtree, following the self-contained-subtree precedent set by `internal/decision` (FEATURE-0013). HTTP handlers are new files in `internal/api` (the established handler package); route registration is wired through `internal/server`.
- **DD-02 — Reuse the shared grammar and validation pipeline.** All request bodies flow through the FEATURE-0012 `internal/apivalid` layer 5–7 pipeline (`DefaultingStage`, `ValidationStage`, `StageSet`). Errors are emitted only as FEATURE-0012 `apiproblem` codes. `metadata.scopeRef` is represented with `apimeta.ScopeRef`; references use `apimeta.TypedRef` with `internal/apiref` constraints; conditions use `internal/apicond`. No parallel grammar is created.
- **DD-03 — Canonicalize the shared ScopeKind grammar without preserving a legacy canonical surface.** `internal/apimeta.ScopeKind` is a shared metadata type and must be aligned to the approved seven canonical values: `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `CloudPlatform`, and `CloudProvider`. FEATURE-0015 tasks explicitly update the shared grammar, its validators, fixtures, and tests. The retired `Provider` value is removed from active ScopeKind validation and is rejected by canonical paths. Retained alpha resources may retain their historical resource-kind names in repository assets, but their scope metadata and tests must compile against canonical `CloudProvider`; no legacy `Provider` ScopeKind compatibility API is exposed. This is a compatibility implementation required by the already-approved `VS0-SCHEMA-001` and DEC-0059 canonical bootstrap, not a new product semantic or a grammar fork.
- **DD-04 — api-server is the sole status writer; no controllers.** All owned-resource `status` is written only by the api-server request path and the deterministic scheduler expiry transition (`VS0-WRITER-005`; ADH-2026-046 decision 1). No topology, stack, provider, or participation controller exists. The deterministic scheduler is a system actor that invokes the api-server expiry transition; it is not an additional writer or controller authority.
- **DD-05 — In-memory store only, no durable/external dependency.** The store is an in-memory registry consistent with the requirements' no-durable-or-external-persistence-dependency boundary. This is why `DEPENDENCY_UNAVAILABLE`/503 is never used for the audit-append-failure outcome; that outcome maps to inherited `INTERNAL_ERROR`/500 (REQ-F15-15).
- **DD-06 — Exact route model: 22 logical paths, 35 registrations.** Handlers register exactly 35 explicit Go 1.22 `http.ServeMux` method/path patterns for 22 logical endpoint paths (architecture §7.13; ADH-2026-051). No path-only or catch-all registration, reflection, or internal HTTP-method dispatch is used; the approved complete `{uid}` wildcard segment is required where declared. Unsupported methods return HTTP 405 without introducing a top-level Problem code and create no new endpoint path.
- **DD-07 — Deterministic scope derivation, no client-supplied scope.** Each owned kind derives `metadata.scopeRef` from exactly one server-side source (architecture §7.8). A client may never supply or select `metadata.scopeRef`; any attempt is rejected under the closed create contract.
- **DD-08 — Canonical idempotency namespace, digest, reservation, and retention.** Idempotency is scoped by `(authenticated principal ID, normalized registered method/path pattern, concrete item target UID when the route has {uid}, Idempotency-Key)`. The method/path component is the explicit Go 1.22 registration pattern; collection creates have no item target UID and remain differentiated by their classified digest.

  The canonical digest is SHA-256 of deterministic canonical JSON produced only after phase-two strict decode and closed-contract classification. It includes the normalized accepted request body and the normalized registered method/path pattern. For an existing-participation item action, it also includes the normalized syntactically valid `If-Match` token; collection creates have no `If-Match` component. It excludes `Idempotency-Key`, authorization credentials and grants, request correlation identifiers, transport-only headers, server-assigned metadata, derived `metadata.scopeRef`, status, timestamps, generated UID, and current `resourceVersion`. Object keys are serialized in stable sorted order; arrays retain their classified request order.

  A namespace record moves only through `InFlight → Completed` or `InFlight → Aborted`. The first valid caller atomically reserves `InFlight`. A concurrent same-key/same-digest caller waits until the reservation reaches a terminal state: on `Completed`, it receives the stored successful status/body with `Content-Type: application/json`, without a second mutation, lifecycle transition, or AuditEvent; on `Aborted`, the coordinator atomically removes the record under its lock and the caller retries reservation from current request state. A same-key/different-digest caller returns `CONFLICT` / 409 with `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`.

  Before any completed replay is returned, current authentication, current server-resolved authorization, and safe target/reference access are repeated for the incoming request. A revoked, narrowed, or otherwise insufficient grant receives its established denial and no stored result is disclosed. Only successful collection creates and successful participation item-action transitions are replayable. Validation, uniqueness, lifecycle-source-state, stale-version, audit-append failure, recovered panic, and graceful shutdown outcomes abort the reservation, wake waiters, and never create a completed replay record.

  Current-version comparison occurs after same-digest replay lookup. Therefore a same-key/same-normalized-`If-Match` replay returns the original result before the current `If-Match` version comparison, but only after current authentication, current authorization, safe reference access, and strict classification. The same key with a different normalized action `If-Match` is a different digest and returns `CONFLICT`; all eight action routes have this route-level test.

  Validation failure, current-version failure, audit-append failure, recovered panic, and graceful server shutdown abort an `InFlight` reservation, wake all waiters, and leave no completed replay record unless the outcome was explicitly completed as replayable. Completed records are retained for the lesser of 24 hours from completion and process lifetime and capped at 10,000 per process. On completed-record insertion, evict the record with the earliest expiry; if expiry ties, evict lexical namespace first. Never evict `InFlight`. FEATURE-0015 introduces no durable/external idempotency store and makes no cross-process replay guarantee. Every received request gets fresh request/correlation identity; `requestId`, correlation, `ETag`, `Location`, and every other entity/transport header are never replayed. Concurrency tests cover same-key/same-digest winners and waiters, same-key/different-digest conflicts (including changed action `If-Match`), abort/re-reservation, authorization revocation before replay, and action-versus-expiry races.
- **DD-09 — Audit/mutation and audited-denial coordination.** One publication coordinator stages every required audited success and every required audited denial. It appends the required FEATURE-0013 AuditEvent before publishing a mutation, lifecycle/status transition, idempotency completion, or audited-denial response. On append failure it suppresses the staged outcome and substitutes `INTERNAL_ERROR`/500 with safe non-disclosure (REQ-F15-15; ADH-2026-049/051). No durable cross-store transaction is claimed; an AuditEvent already appended before a later in-process publication failure may remain.
- **DD-10 — Repository-owned, version-pinned ISO-3166-1 alpha-2 dataset.** A compiled-in static dataset of assigned ISO-3166-1 alpha-2 codes as of 2026-08-12 (`internal/cloudmodel/isocodes`) backs assigned-code validation for `jurisdictionCode`, `operatingMarkets[]`, and `countryCode`. It requires no network request and adds no third-party runtime dependency. `administrativeAreaCode` uses syntax-plus-country-prefix validation only; no ISO-3166-2 membership dataset is introduced (ADH-2026-048 decision 1).

---

## 3. Components and repository paths

Paths below reflect the verified live repository convention (flat `internal/*` packages; handlers in `internal/api`; router wiring in `internal/server`; self-contained subtree precedent in `internal/decision`). Concrete file split is a delegated mechanic; the package boundaries and reuse points are fixed by this design.

### 3.1 New FEATURE-0015 components (IMPLEMENT)

| Component | Path (package) | Responsibility |
|-----------|----------------|----------------|
| Canonical resource types | `internal/cloudmodel/model` | metadata/spec/status structs for the seven owned kinds; JSON tags; no alpha-type reuse |
| Assigned ISO dataset | `internal/cloudmodel/isocodes` | compiled-in, version-pinned assigned ISO-3166-1 alpha-2 set (2026-08-12); membership lookup |
| Deterministic validation | `internal/cloudmodel/validate` | pure functions over `cloudmodel/model` types: closed create-contract classification, scope-kind subset, immutable-field, reference/containment, assigned-code, header/body-input checks |
| In-memory store | `internal/cloudmodel` (store file) | imports `cloudmodel/model` and `validate`; per-kind registries; name-uniqueness within server-derived scope; non-terminal participation pair uniqueness; resourceVersion; publication coordinator |
| Participation lifecycle | `internal/cloudmodel` (lifecycle file) | `VS0-STATE-001` transitions; independent holds; derived effective phase; scheduler expiry transition |
| Idempotency coordinator | `internal/cloudmodel` (idempotency file) | reservation states `InFlight`/`Completed`/`Aborted`; digest comparison; waiter/abort/release |
| Bootstrap-grant resolver | `internal/cloudmodel` (grant file) | server-resolved deterministic grants; header/body grant claims ignored; no persisted roles/memberships/assignments |
| Audit/mutation coordinator | `internal/cloudmodel` (audit file) | orders required AuditEvent append before publication; append-failure → unpublished + `INTERNAL_ERROR` |
| HTTP handlers | `internal/api` (new files per kind) | collection/item/action handlers; media-type gate; header preconditions; Problem responses |
| Route registration | `internal/server` | 35 explicit method/path patterns for 22 logical paths (DD-06) |

### 3.2 Reused prior-feature components (CONTRACT_ONLY/NO_TASK — consumed, not modified)

| Reused contract | Path (package) | Owner |
|-----------------|----------------|-------|
| Problem Details, ErrorCode, Violation | `internal/apiproblem` | FEATURE-0012 (`VS0-SCHEMA-004`) |
| ObjectMeta, TypedRef, ScopeRef, ScopeKind, scope normalization | `internal/apimeta` | FEATURE-0012 (`VS0-SCHEMA-001`, `VS0-SCHEMA-002`) |
| Condition | `internal/apicond` | FEATURE-0012 (`VS0-SCHEMA-003`) |
| Reference constraints | `internal/apiref` | FEATURE-0012 |
| Validation pipeline (layers 5–7) | `internal/apivalid` | FEATURE-0012 |
| AuditEvent envelope + linkage + writer | `internal/decision` | FEATURE-0013 (`VS0-SCHEMA-007`, `VS0-WRITER-011`) |
| Request correlation / request ID | `internal/requestctx`, `internal/server/middleware` | prior features |
| Observability baseline (structured logs/metrics) | existing baseline | prior features |

FEATURE-0015 extends the shared ScopeKind implementation only as the explicitly required compatibility prerequisite in DD-03. All other FEATURE-0012 packages and contracts are consumed by reference, not renamed, forked, or transferred.

---

## 4. Data/API representation

### 4.1 Owned resource shapes

All owned kinds use the canonical `metadata`/`spec`/`status` shape. `metadata` (identity, `scopeRef`, uid, resourceVersion, generation, timestamps) and all `status` fields are api-server-owned. Clients supply only the fields listed in §4.3.

| Kind | Registry ID | scopeRef.kind | Key spec fields (F0015-owned) | Key status fields (api-server-owned) |
|------|-------------|---------------|-------------------------------|--------------------------------------|
| CloudPlatform | `VS0-SCHEMA-008` | Platform | immutable `spec.ownerRegistration{legalName,registrationIdentifier,jurisdictionCode}`; optional `spec.description` | `status.phase=Active` |
| CloudProvider | `VS0-SCHEMA-009` | Platform | non-empty unique uppercase `spec.operatingMarkets[]`; optional `spec.displayName` | `status.phase=Active` |
| CloudProviderParticipation | `VS0-SCHEMA-010` | CloudPlatform | `spec.cloudPlatformRef` (UID-pinned), `spec.cloudProviderRef` (UID-pinned), `spec.environment` (`development`); no mutable spec fields | `status.phase`, `status.platformSuspended`, `status.providerSuspended`, `status.requestExpiresAt` |
| HostingLocation | `VS0-SCHEMA-011` | CloudProvider | uppercase `spec.countryCode`, non-empty `spec.locality`, optional `spec.administrativeAreaCode`, optional `spec.description` | `status.phase=Active` |
| Datacenter | `VS0-SCHEMA-012` | CloudProvider | immutable `spec.hostingLocationRef`; optional `spec.description` | `status.phase=Active` |
| FaultDomain | `VS0-SCHEMA-013` | CloudProvider | immutable `spec.datacenterRef`; optional `spec.description` | `status.phase=Active` |
| InfrastructureStack | `VS0-SCHEMA-014` | CloudProvider | immutable `spec.faultDomainRef`; optional `spec.description` | `status.phase=Active` |

Field shapes, types, limits, enums, defaults, references, and mutability are taken exactly from the registry entries `VS0-SCHEMA-008` through `VS0-SCHEMA-014`; this design adds none. `spec.providerSelectionModes` and `spec.permittedHostingLocationRefs` are FEATURE-0021-owned and are never accepted, stored, defaulted, validated, or exposed (§10).

### 4.2 Route model (22 logical paths; 35 explicit Go registrations)

All routes require authentication. Each item or action handler obtains its resource UID through `r.PathValue("uid")`. No path-only or catch-all registration, reflection, or internal HTTP-method dispatch is used; the approved complete `{uid}` wildcard segment is not a catch-all.

| Logical path | Explicit Go 1.22 registration(s) | PathValue | Handler responsibility | Unsupported-method behavior |
|---|---|---|---|---|
| CloudPlatform collection | `GET /apis/core.sovrunn.io/v1alpha1/cloud-platforms`; `POST /apis/core.sovrunn.io/v1alpha1/cloud-platforms` | — | list; create | Only GET/POST registered |
| CloudPlatform item | `GET /apis/core.sovrunn.io/v1alpha1/cloud-platforms/{uid}`; `PATCH /apis/core.sovrunn.io/v1alpha1/cloud-platforms/{uid}` | `uid` | get; permitted PATCH | PUT/DELETE return 405; PATCH only with `application/merge-patch+json` |
| CloudProvider collection | `GET /apis/core.sovrunn.io/v1alpha1/cloud-providers`; `POST /apis/core.sovrunn.io/v1alpha1/cloud-providers` | — | list; create | Only GET/POST registered |
| CloudProvider item | `GET /apis/core.sovrunn.io/v1alpha1/cloud-providers/{uid}`; `PATCH /apis/core.sovrunn.io/v1alpha1/cloud-providers/{uid}` | `uid` | get; permitted PATCH | PUT/DELETE return 405; PATCH only with `application/merge-patch+json` |
| CloudProviderParticipation collection | `GET /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations`; `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations` | — | list; participation create | Only GET/POST registered |
| CloudProviderParticipation item | `GET /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}` | `uid` | get | No PATCH registration; only the explicit action routes below mutate lifecycle |
| HostingLocation collection | `GET /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations`; `POST /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations` | — | list; create | Only GET/POST registered |
| HostingLocation item | `GET /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/{uid}`; `PATCH /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/{uid}` | `uid` | get; permitted PATCH | PUT/DELETE return 405; PATCH only with `application/merge-patch+json` |
| Datacenter collection | `GET /apis/infrastructure.sovrunn.io/v1alpha1/datacenters`; `POST /apis/infrastructure.sovrunn.io/v1alpha1/datacenters` | — | list; create | Only GET/POST registered |
| Datacenter item | `GET /apis/infrastructure.sovrunn.io/v1alpha1/datacenters/{uid}`; `PATCH /apis/infrastructure.sovrunn.io/v1alpha1/datacenters/{uid}` | `uid` | get; permitted PATCH | PUT/DELETE return 405; PATCH only with `application/merge-patch+json` |
| FaultDomain collection | `GET /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains`; `POST /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains` | — | list; create | Only GET/POST registered |
| FaultDomain item | `GET /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/{uid}`; `PATCH /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/{uid}` | `uid` | get; permitted PATCH | PUT/DELETE return 405; PATCH only with `application/merge-patch+json` |
| InfrastructureStack collection | `GET /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks`; `POST /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks` | — | list; create | Only GET/POST registered |
| InfrastructureStack item | `GET /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/{uid}`; `PATCH /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/{uid}` | `uid` | get; permitted PATCH | PUT/DELETE return 405; PATCH only with `application/merge-patch+json` |
| Participation accept | `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/accept` | `uid` | Pending → Active | Empty JSON body; If-Match + Idempotency-Key |
| Participation reject | `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/reject` | `uid` | Pending → Rejected | Empty JSON body; If-Match + Idempotency-Key |
| Participation withdraw | `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/withdraw` | `uid` | Pending → Withdrawn | Empty JSON body; If-Match + Idempotency-Key |
| Participation suspend | `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/suspend` | `uid` | Set acting party hold | Empty JSON body; If-Match + Idempotency-Key |
| Participation resume | `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/resume` | `uid` | Clear acting party hold | Empty JSON body; If-Match + Idempotency-Key |
| Participation request release | `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/request-release` | `uid` | Active/Suspended → Terminating | Empty JSON body; If-Match + Idempotency-Key |
| Participation accept release | `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/accept-release` | `uid` | Terminating → Terminated | Empty JSON body; If-Match + Idempotency-Key |
| Participation decline release | `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/decline-release` | `uid` | Terminating → prior effective Active/Suspended | Empty JSON body; If-Match + Idempotency-Key |

For each participation item-action route, `{uid}` is a complete Go 1.22 `http.ServeMux` wildcard segment, extracted with `r.PathValue("uid")`; the retired `/{uid}:<action>` form is not registered.

Registration arithmetic is exact:

- Seven collection paths × GET and POST = 14 registrations.
- Six PATCHable item paths × GET and PATCH = 12 registrations.
- One CloudProviderParticipation item GET path = 1 registration.
- Eight participation item-action POST paths = 8 registrations.
- Total = 35 explicit Go 1.22 method/path registrations across 22 logical paths.

The PATCH allowlist is exact:

| Kind | Permitted PATCH fields |
|---|---|
| CloudPlatform | `spec.description` |
| CloudProvider | `spec.displayName`, `spec.operatingMarkets` |
| HostingLocation | `spec.description` |
| Datacenter | `spec.description` |
| FaultDomain | `spec.description` |
| InfrastructureStack | `spec.description` |
| CloudProviderParticipation | None; no PATCH surface |

For CloudPlatform, CloudProvider, and topology resources, PUT and DELETE return 405 and only the listed PATCH fields are mutable (`VS0-CF-F15-07`, `VS0-CF-F15-13`, `VS0-CF-F15-14`). CloudProviderParticipation lifecycle mutation is available only through the eight explicit item-action routes.

### 4.3 Closed collection-create request contract (ADH-2026-047 decision 1)

Per architecture §7.7 and §7.10. The client supplies only the columns below. The client must not supply `metadata.uid`, generation, resourceVersion, timestamps, `metadata.scopeRef`, any `status` field, a field owned by another feature, or an unknown field.

| Kind | Client-required create fields | Client-optional create fields | Server-assigned outcome |
|---|---|---|---|
| CloudPlatform | `metadata.name`; `spec.ownerRegistration.legalName`; `spec.ownerRegistration.registrationIdentifier`; `spec.ownerRegistration.jurisdictionCode` | `metadata.displayName`; `spec.description` | Platform-root `scopeRef`; `status.phase=Active`; `201` |
| CloudProvider | `metadata.name`; non-empty `spec.operatingMarkets[]` | `metadata.displayName`; `spec.displayName` | Platform-root `scopeRef`; `status.phase=Active`; `201` |
| HostingLocation | `metadata.name`; `spec.countryCode`; `spec.locality` | `spec.administrativeAreaCode`; `spec.description` | CloudProvider `scopeRef` from `topology.write` grant; `status.phase=Active`; `201` |
| Datacenter | `metadata.name`; `spec.hostingLocationRef` | `spec.description` | CloudProvider `scopeRef` from resolved parent; `status.phase=Active`; `201` |
| FaultDomain | `metadata.name`; `spec.datacenterRef` | `spec.description` | CloudProvider `scopeRef` from resolved parent; `status.phase=Active`; `201` |
| InfrastructureStack | `metadata.name`; `spec.faultDomainRef` | `spec.description` | CloudProvider `scopeRef` from resolved parent; `status.phase=Active`; `201` |
| CloudProviderParticipation | `metadata.name`; `spec.cloudPlatformRef`; `spec.cloudProviderRef`; `spec.environment` (`development`) | (none) | CloudPlatform `scopeRef` from `spec.cloudPlatformRef`; `status.phase=Pending`; both holds `false`; `status.requestExpiresAt=createdAt+7d`; `201` |

A successful PATCH returns `200` with the updated resource. Topology `metadata.name` is immutable identity; only `spec.description` is PATCHable for topology resources (ADH-2026-047 decision 3).

### 4.4 Scope derivation (single source per kind; ADH-2026-047 decision 2)

- CloudPlatform, CloudProvider → immutable deployment Platform-root scope, server-derived.
- CloudProviderParticipation → immutable CloudPlatform scope from resolved `spec.cloudPlatformRef`; referenced CloudPlatform must exist and its UID must equal the resulting `metadata.scopeRef.uid`.
- HostingLocation → immutable CloudProvider scope from the CloudProvider UID bound to the authenticated, server-resolved `topology.write` grant.
- Datacenter/FaultDomain/InfrastructureStack → immutable CloudProvider scope from the resolved immutable parent reference.

### 4.5 Participation lifecycle representation (`VS0-STATE-001`)

`status.phase` is api-server-owned and stores only current lifecycle facts (phase, `platformSuspended`, `providerSuspended`, retained `requestExpiresAt`, standard system status fields); transition history is durable FEATURE-0013 AuditEvent evidence and is not duplicated. Transitions exactly as architecture §7.4:

```text
absent → Pending (create; participation.request)
Pending → Active (accept) | Rejected (reject) | Withdrawn (withdraw) | Expired (scheduler at requestExpiresAt)
Active/Suspended: effective Suspended when either hold true; Active only when both holds false (resume clears own hold)
Active/Suspended → Terminating (request-release)
Terminating → Terminated (accept-release) | prior effective Active/Suspended (decline-release)
```

For accepted participation, the server stores `phase=Active` when both holds are false and `phase=Suspended` when either hold is true. Pending, Terminating, and terminal phases override the hold-derived phase; holds remain retained but immutable outside Active/Suspended. Exactly one non-terminal participation may exist per `(cloudPlatformUID, cloudProviderUID)` pair; a new request after a terminal record receives a new UID.

### 4.6 Deterministic participation-expiry scheduler

`ParticipationExpiryScheduler` is an internal FEATURE-0015 component, not a public route, controller, or additional status writer. It is constructed with an injected clock, injected ticker, scheduler-execution-ID generator, canonical store, and the api-server participation-expiry operation.

The api-server starts the scheduler during server startup. Graceful server shutdown stops future ticks and waits for any active tick to finish before the store and AuditEvent writer are released.

The scheduler's next wake-up is the earliest Pending `status.requestExpiresAt`; a create or successful lifecycle change signals recomputation. If no Pending participation exists it waits for a signal. The injected ticker/clock test harness may advance directly to that due time, but production does not use an arbitrary polling interval. A wake-up processes every due candidate in the stable order below before computing the next due time.

Each tick:

1. obtains Pending participation candidates in stable ascending UID order;
2. considers only records whose `status.requestExpiresAt` is at or before the injected current time;
3. acquires the same concurrency guard used by participation item-action handlers;
4. rechecks that the record is still Pending, still expired, and still at the version observed for this expiry attempt;
5. invokes the api-server expiry transition rather than writing status directly.

If a participation action wins the concurrency race before the recheck, expiry is a no-op and produces no AuditEvent. If expiry wins, the api-server performs exactly one `Pending → Expired` transition and appends exactly one correlated, redacted FEATURE-0013 AuditEvent naming the system actor and scheduler execution ID (`VS0-CF-F15-16`). Repeated ticks and retries are idempotent: an already non-Pending participation is never transitioned or audited again.

The expiry transition follows the same audit-before-publication rule as an authenticated lifecycle action. If its required AuditEvent append fails, it returns/records `INTERNAL_ERROR` / 500 for the attempted operation and publishes neither the Expired status transition nor an idempotency completion (`VS0-CF-F15-24`). The scheduler does not create a public response, public route, external persistence dependency, or separate controller authority.

---

## 5. Validation and deterministic error behavior

### 5.1 Per-route pipeline order (fail-closed)

Every owned route follows the shared prefix, then exactly one route-specific branch below. The first failing step stops processing. No later step, resource mutation, status/lifecycle write, idempotency completion, or AuditEvent occurs unless the approved audited-denial boundary expressly requires one.

**Shared prefix**

1. **Authentication** — Missing or invalid authentication returns `AUTH_REQUIRED` / 401 with no mutation, lifecycle/status write, idempotency record, or AuditEvent (`REQ-F15-24`; `VS0-CF-F15-31`).
2. **Route and actual-method gate** — The handler validates the registered method/path pattern, then checks `r.Method`. HEAD matched by ServeMux to GET returns 405 with no side effect; CloudPlatform, CloudProvider, and topology PUT/DELETE return 405; CloudProviderParticipation has no PATCH registration (`VS0-CF-F15-07`, `VS0-CF-F15-34`).
3. **Reserved and ordinary headers** — After authentication, a present `X-Sovrunn-Bootstrap-Grant` is denied/audited before media or body validation. Otherwise, validate `Idempotency-Key` only on a collection-create or participation item-action POST; reject `If-Match` on participation collection create; require syntactically valid `If-Match` on PATCH and existing-participation actions (`VS0-CF-F15-17`, `30`, `35`, `38`).
4. **Media and safe phase-one decoding** — Without the reserved header, PATCH requires `application/merge-patch+json` before body decode. Participation item actions accept only EOF/zero bytes; whitespace, `{}`, or another non-empty body is `MALFORMED_REQUEST`, except a valid duplicate-free top-level `bootstrapGrant`, which is denied/audited under F15-22. Apply bounded-body, duplicate-member, and syntax-safe decoding; the valid reserved body member is denied/audited before authorization. Otherwise extract only route-safe identifiers and references.
5. **Coarse action-grant lookup** — Resolve the principal’s current server-side grant by the route’s required action only. This check does not yet evaluate a resource-derived scope or optional target-UID restriction; no header/body claim can expand or substitute for it (`VS0-CF-F15-22`).
6. **CloudPlatform-root gate when creating a dependent kind** — Only a CloudProvider, topology, or participation collection create checks that a CloudPlatform exists before reference resolution. Absence returns the audited `CONFLICT` / `VS0_CLOUDPLATFORM_ROOT_REQUIRED` (`VS0-CF-F15-12`).
7. **Safe resolution, scope derivation, and exact authorization** — Resolve a direct GET target or body/parent reference only through the caller’s current access boundary; derive its server-owned resource scope; then evaluate the grant’s exact scope and optional target-UID restriction. Inaccessible direct GET/reference returns safe `RESOURCE_NOT_FOUND` / 404. A visible-but-invalid topology graph proceeds to structural validation. LIST evaluates its registered read action and scope filter here: no grant returns audited 403, a scoped grant returns only visible resources, and an authorized empty scope returns 200 empty (`VS0-CF-F15-21`, `23`, `39`, `X03`).

**Read branches**

- **GET and LIST** end after step 7 (and projection assembly). They never require, look up, reserve, complete, or abort idempotency state.

**Mutation branches**

8. **Strict phase-two decode and closed-contract classification** — Collection creates and PATCHes strictly decode/classify the complete body against their registered request contract; client `metadata.scopeRef`, server-owned metadata, `status`, FEATURE-0021-deferred participation fields, and unknown fields are rejected. An invalid or unclassified body is never digested, stored, or replayed. Participation item actions have already completed their empty-body contract at step 4 (`VS0-CF-F15-26`, `28`).
9. **Conditional idempotency branch** — Only the seven collection-create POSTs and eight participation item-action POSTs compute a digest and look up/reserve `(principal, registered method/path, target UID when present, Idempotency-Key)`. An item-action digest includes normalized `If-Match`; collection-create digests do not. Current authentication, exact authorization, and safe target/reference access have already succeeded. Same-key/same-digest completed replay returns the stored success before current-version comparison; a changed body or changed normalized action `If-Match` returns `CONFLICT` / `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`. PATCH never requires, looks up, reserves, completes, or aborts idempotency state (`VS0-CF-F15-18`, `19`).

Completed records retain for the lesser of 24 hours and process lifetime and are capped at 10,000 per process. On insertion, evict the completed record with the earliest expiry, then lexical namespace on ties; never evict `InFlight`. Expiry/eviction makes the next request new; cancellation detaches only its waiter; every post-reservation failure aborts and wakes waiters. Remove an Aborted namespace atomically before re-reservation. A replay returns only stored successful status/body plus `Content-Type: application/json`; request correlation is fresh and no entity/transport headers are copied (`VS0-CF-F15-36`, `40`).

10. **Server-side scope derivation for creates** — Each create derives `metadata.scopeRef` from its sole registered source; the client cannot supply or select scope (`VS0-CF-F15-27`).
11. **Semantic, reference-integrity, and uniqueness validation** — Validate schema constraints, scope subset, ISO semantics, relationship references, topology graph, name/pair uniqueness, and action authorization. A collection-create/action failure aborts its reservation; PATCH has no reservation.
12. **Current resource-version comparison and publication-lock recheck** — Every PATCH and existing-participation action compares `If-Match` after the applicable validation. Actions compare before lifecycle-source-state validation. PATCH rechecks under the publication lock immediately before required AuditEvent append/publication. Missing/malformed/stale values return 412 with no mutation, AuditEvent, or idempotency record (`VS0-CF-F15-17`, `38`).
13. **Lifecycle source-state validation** — An item action with current `If-Match` validates its source state. An invalid state returns `CONFLICT` / `VS0_PARTICIPATION_STATE_INVALID`; it aborts the action reservation and creates no status write or AuditEvent (`VS0-CF-F15-20`).
14. **Publication coordinator: under-lock recheck, AuditEvent append, and publication** — One coordinator serializes every audited success and audited denial. For collection create it takes the publication lock, rechecks name uniqueness or non-terminal participation-pair uniqueness against current state, then stages the mutation and eligible idempotency completion. A conflicting different-key creator receives its registered conflict outcome; it never publishes a success AuditEvent. For every required audited success or denial, the coordinator appends the redacted FEATURE-0013 AuditEvent before publishing the staged response. An append failure suppresses the staged success or denial and substitutes `INTERNAL_ERROR` / 500; neither the original response nor any idempotency completion is published. An already-appended event may remain after a later in-process publication failure; no durable cross-store transaction is claimed (`VS0-CF-F15-24`).

### 5.1.1 Two-phase decoding and exact input outcomes

Phase one is bounded, syntax-safe decoding used only to establish request safety before authorization and safe reference access. It enforces the inherited body-size limit, rejects malformed JSON and duplicate JSON members, and extracts only route-safe identifiers and references. It never persists, defaults, validates, or accepts the complete request body.

Phase two occurs only after current authorization, the root gate, and safe reference resolution. It strictly decodes the complete body, rejects unknown or duplicate fields, and classifies every field against the registered route request contract before canonical digesting or idempotency lookup.

| Input condition | Exact outcome | Durable AuditEvent |
|---|---|---|
| Body exceeds the inherited request-size limit | `REQUEST_TOO_LARGE` / 400 | No |
| Malformed JSON | `MALFORMED_REQUEST` / 400 | No |
| Duplicate JSON member | `DUPLICATE_FIELD` / 400 | No |
| Unknown field | `UNKNOWN_FIELD` / 400 | No |
| FEATURE-0021-deferred `spec.providerSelectionModes` or `spec.permittedHostingLocationRefs` on participation create | `UNKNOWN_FIELD` / 400 | No |
| Client-supplied `metadata.scopeRef` | `VALIDATION_FAILED` / 422; client scope never selects authority | No |
| Client-supplied `status` field after authentication | `AUTHORIZATION_DENIED` / 403 with `VS0_STATUS_FIELD_WRITE` | Exactly one redacted AuditEvent |
| Client-supplied api-server-owned identity or metadata field after authentication | `AUTHORIZATION_DENIED` / 403 with `VS0_SYSTEM_OWNED_FIELD_WRITE` | Exactly one redacted AuditEvent |
| Present `X-Sovrunn-Bootstrap-Grant` after authentication | `AUTHORIZATION_DENIED` / 403 before body decode | Exactly one redacted AuditEvent |
| Valid duplicate-free top-level `bootstrapGrant` member | `AUTHORIZATION_DENIED` / 403 before authorization/reference access and phase-two unknown-field classification | Exactly one redacted AuditEvent |
| Syntactically valid but unassigned ISO-3166-1 alpha-2 value | `VALIDATION_FAILED` / 422 | No |
| Missing, empty, malformed, or over-length required `Idempotency-Key` | `MALFORMED_REQUEST` / 400 | No |
| `If-Match` supplied on participation collection create | `MALFORMED_REQUEST` / 400 | No |
| Missing or malformed `If-Match` on an existing participation action | `STALE_RESOURCE_VERSION` / 412 | No |
| Non-current `If-Match` on an existing participation action | `STALE_RESOURCE_VERSION` / 412 | No |
| Non-empty JSON body on a participation item-action route without valid duplicate-free top-level `bootstrapGrant` | `MALFORMED_REQUEST` / 400 | No |

The approved audited-denial categories remain closed: status write, api-server-owned identity/metadata write, forged-grant attempt, CloudPlatform-root denial, LIST without read grant, and safe inaccessible-reference denial. This table does not broaden that boundary.

### 5.2 Error-code reuse (FEATURE-0012 closed set only)

All top-level Problem codes are existing FEATURE-0012 codes (`internal/apiproblem`): `MALFORMED_REQUEST` / 400, `UNKNOWN_FIELD` / 400, `DUPLICATE_FIELD` / 400, `REQUEST_TOO_LARGE` / 400, `AUTH_REQUIRED` / 401, `AUTHORIZATION_DENIED` / 403, `RESOURCE_NOT_FOUND` / 404, `ALREADY_EXISTS` / 409, `CONFLICT` / 409, `STALE_RESOURCE_VERSION` / 412, `UNSUPPORTED_MEDIA_TYPE` / 415, `VALIDATION_FAILED` / 422, and `INTERNAL_ERROR` / 500. PUT and DELETE return 405 for CloudPlatform, CloudProvider, and topology resources; successful collection create returns 201; successful permitted PATCH returns 200.

Slice 0 violation codes (`VS0_PARTICIPATION_STATE_INVALID`, `VS0_PARTICIPATION_DUPLICATE`, `VS0_STATUS_FIELD_WRITE`, `VS0_SYSTEM_OWNED_FIELD_WRITE`, `VS0_SCOPE_KIND_INVALID`, `VS0_PATCH_IMMUTABLE_FIELD`, `VS0_CLOUDPLATFORM_ROOT_REQUIRED`, `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`, `VS0_SCOPE_REFERENCE_MISMATCH`, `VS0_TOPOLOGY_PROVIDER_MISMATCH`, `VS0_AUTHORIZATION_SAFE_DENIAL`, and `VS0_AUTH_REQUIRED`) appear only in `violations[].code`, never as top-level Problem codes. No new top-level Problem code or violation code is introduced (`REQ-F15-12`; ADH-2026-048 decision 2).

### 5.3 Determinism

All validation functions in `internal/cloudmodel/validate` are pure and deterministic (no I/O, no wall-clock reads except the injected `createdAt` for `requestExpiresAt`, no map-iteration-order dependence in outputs). The scheduler expiry transition is idempotent.

---

## 6. Security, privacy, observability, compatibility, and versioning

### 6.1 Security and authorization

Authentication is required on every owned method/path pattern (REQ-F15-24). Authorization uses server-resolved, deterministic bootstrap grants only; grant claims in a header or body are ignored and never expand or substitute for the server-resolved grant (REQ-F15-17; VS0-CF-F15-22). FEATURE-0015 persists no roles, memberships, or assignments; identity-governance ownership remains with FEATURE-0018 by reference. Writer boundaries are consumed by reference: `VS0-WRITER-002` (cloud-platform-admin CloudPlatform.spec), `VS0-WRITER-003` (cloud-provider-admin CloudProvider/topology spec), `VS0-WRITER-004` (participation actions, api-server), `VS0-WRITER-005` (status, api-server sole writer). Clients never write status or api-server-owned identity/metadata. Cross-provider and inaccessible references use safe denial with no existence disclosure.

**Security note (network exposure):** all owned routes are authenticated request paths behind the existing api-server; no unauthenticated endpoint is introduced. The bootstrap principal receives narrowly scoped, server-resolved grants only.

| Registrations | Exact server-resolved action | Scope / optional UID narrowing | Failure behavior |
|---|---|---|---|
| CloudPlatform POST/PATCH | `cloudplatform.write` | Platform root; PATCH may narrow to item UID | existing writer denial; PATCH token failure is 412 |
| CloudProvider POST/PATCH | `cloudprovider.write` | Platform root; PATCH may narrow to item UID | existing writer denial; PATCH token failure is 412 |
| HostingLocation, Datacenter, FaultDomain, InfrastructureStack POST/PATCH | `topology.write` | CloudProvider UID; descendants resolve that provider through their parent | inaccessible reference is safe 404; visible mismatch is 422 |
| Participation POST create | `participation.request.provider` | referenced CloudProvider UID plus referenced CloudPlatform relation | root/pair/field rules before publication |
| accept, reject, request-release | `participation.accept.platform`, `participation.reject.platform`, `participation.request-release.platform` | existing participation's CloudPlatform UID | current access then lifecycle rule |
| withdraw, accept-release, decline-release | `participation.withdraw.provider`, `participation.accept-release.provider`, `participation.decline-release.provider` | existing participation's CloudProvider UID | current access then lifecycle rule |
| suspend, resume | either `participation.suspend.platform` / `participation.resume.platform`, or `participation.suspend.provider` / `participation.resume.provider` | existing participation's CloudPlatform UID for platform grant or CloudProvider UID for provider grant | resolve exactly one matching current scoped grant; zero or multiple matches fail closed; the resolved actor changes only its own hold |
| CloudPlatform GET/LIST | `cloudplatform.read` | Platform root; GET may narrow to target UID | LIST without grant 403/audited; inaccessible GET safe 404 |
| CloudProvider GET/LIST | `cloudprovider.read` | Platform root; GET may narrow to target UID | LIST without grant 403/audited; inaccessible GET safe 404 |
| Topology GET/LIST | `topology.read` | CloudProvider UID; GET may narrow to target UID | LIST without grant 403/audited; inaccessible GET safe 404 |
| Participation GET/LIST | `participation.read` | derived CloudPlatform UID; GET may narrow to target UID | LIST without grant 403/audited; inaccessible GET safe 404 |

### 6.1.1 PATCH and collection mechanics

PATCH parses `application/merge-patch+json` as RFC 7396 into a staged clone. Before merge it classifies every patch member against the per-kind allowlist. `null` removes only `spec.description` or `spec.displayName` where those optional mutable fields exist; `null` for required `spec.operatingMarkets` is `VALIDATION_FAILED`; immutable identity/reference fields, server-owned metadata/status, deferred fields, and unknown fields retain their established rejection. PATCH requires a parsed `If-Match`; after merged-clone validation the publication coordinator takes its lock, compares `If-Match` to the then-current stored `resourceVersion`, appends the required audit event, and publishes only if the lock recheck succeeds. Successful LIST filters by its exact current read action/scope and sorts the surviving resources by ascending `metadata.uid` (`VS0-CF-F15-37`, `VS0-CF-F15-38`, `VS0-CF-F15-39`).

The import graph is one-way: `internal/cloudmodel/model` imports only shared grammar packages; `internal/cloudmodel/validate` imports `cloudmodel/model` plus shared grammar packages; `internal/cloudmodel` imports `cloudmodel/model` and `validate`; `internal/api` imports both; `internal/server` imports `api`. Canonical types have one definition in `cloudmodel/model`; validation does not duplicate DTOs or import the store. Tests may import outward through these layers but production packages may not reverse an arrow.

### 6.2 Privacy, classification, and redaction

AuditEvents are redacted with request correlation and no secret or inaccessible-resource disclosure (REQ-F15-15). No credential values or protected handles appear in any audit record; references are UID-pinned identifiers only. Classification follows the registry (CloudPlatform Internal with `ownerRegistration` never customer-visible; CloudProvider/topology/participation Provider-confidential). Non-Active/non-visible participation phases are never end-user visible; customer/Organization visibility eligibility is out of scope (§10).

### 6.3 Observability

Redacted structured logs and metrics use the existing observability baseline and are explicitly not AuditEvents. Request correlation (requestId; resource/participation UID) is carried on both logs and the required AuditEvent. The durable-audit categories and non-audited categories are exactly those in architecture §7.11 / requirements §4.15.

### 6.4 Compatibility and versioning

CloudPlatform and CloudProvider are served at `core.sovrunn.io/v1alpha1`; CloudProviderParticipation is served at `governance.sovrunn.io/v1alpha1`; HostingLocation, Datacenter, FaultDomain, and InfrastructureStack are served at `infrastructure.sovrunn.io/v1alpha1`. Inherited contracts (Problem Details, ObjectMeta, TypedRef, Condition, DecisionRecord, AuditEvent, reuse governance) are consumed by exact reference and are not versioned, redefined, or transferred by this feature. The update surface is PATCH-only (`application/merge-patch+json`); for CloudPlatform, CloudProvider, and topology resources, `PUT`/`DELETE` return 405. No migration plan/record/controller/cutover mechanism is introduced or reintroduced (REQ-F15-13; retired tombstones `VS0-SCHEMA-060`, `VS0-SCHEMA-061`, `VS0-WRITER-020`, `VS0-WRITER-021`, `VS0-STATE-011` are never reused).

---

## 7. Test and conformance strategy

### 7.1 Test packages and levels

Following the verified convention (`_test.go` beside implementation; property and race tests where concurrency-sensitive):

- `internal/cloudmodel/validate` — unit and property tests for closed-contract classification; exact per-kind create and PATCH allowlists; canonical seven-value `ScopeKind` validation; immutable and system-owned fields; scope/reference containment; assigned ISO-code validation; and malformed header/body inputs. Tests cover unknown fields, duplicate JSON members, malformed JSON, oversized body, forbidden `status`/`metadata.scopeRef`, deferred FEATURE-0021 participation fields, malformed Idempotency-Key, forbidden If-Match on participation create, and non-empty item-action bodies.

- `internal/cloudmodel` — store tests for server-derived-scope name uniqueness, non-terminal participation-pair uniqueness, resourceVersion, and all lifecycle transitions from `VS0-STATE-001`. Test independent holds and derived phase, zero/multiple matching suspend/resume grant rejection, deterministic scheduler expiry, scheduler start/stop, stable expiry ordering, and repeated/concurrent scheduler ticks. Test 24-hour/10,000-entry completed-record retention, earliest-expiry then lexical-namespace eviction with `InFlight` never evicted, cancellation detachment, replay correlation-header regeneration, changed-`If-Match` same-key action conflict on all eight actions, abort/release after audit failure, and the invariant that an audit-append failure publishes neither an audited success/denial response nor an idempotency replay record. Run different-key concurrent create-name and participation-pair tests with `go test -race`; exactly one winner may append/publish success.

- `internal/api` — table-driven handler tests for every one of the 22 logical API paths and all 35 registered method/path patterns. Assert actual HEAD-to-GET handling returns 405, empty versus whitespace/non-empty action body behavior, forged-header-plus-invalid-media precedence, successful 201/200 response shape, and exact Problem Details outcome. Assert the shared prefix (authentication; reserved-header/media/body classification; coarse action-grant lookup; root gate where applicable; safe resolution, derived scope, and exact authorization) precedes each branch. Assert only the seven collection-create and eight participation-action POSTs perform digest/reservation/replay; GET, LIST, and PATCH neither require nor create idempotency state. Assert PATCH/action version checking, lifecycle validation, audit append, and publication order separately.

- Compatibility tests — update retained alpha fixtures and tests to compile against the canonical shared `ScopeKind` vocabulary. Assert canonical FEATURE-0015 paths accept `CloudPlatform` and `CloudProvider` scope kinds and reject the retired `Provider` scope kind. No legacy scope-kind compatibility surface is exposed.

- Security and safe-denial tests — cover missing/invalid authentication, narrowed or revoked grants on a replay attempt, CloudPlatform-root gating before reference resolution, inaccessible cross-provider references returning `RESOURCE_NOT_FOUND`, and visible-but-invalid cross-provider references returning `VALIDATION_FAILED`. Verify audited versus non-audited denial boundaries and ensure secrets never appear in AuditEvent or Problem Details.

- Conformance tests (existing `internal/apiconform` / `tests/conformance` convention) — implement one local scenario for every FEATURE-0015-owned case in §9, including `VS0-CF-X03` and `VS0-CF-F15-01` through `VS0-CF-F15-41`. Anti-drift/non-runtime acceptance gates remain verified by their named repository checks rather than invented runtime conformance.

### 7.2 Acceptance-scenario coverage mapping

Every approved AC is mapped here to its detailed acceptance scenario and its FEATURE-0015-local conformance case(s), outside the canonical acceptance ledger. Anti-drift/non-runtime items are labelled per ADH-2026-046 decision 3 and carry no runtime proof case.

| AC ID | Acceptance scenario (observable outcome) | FEATURE-0015-local proof case(s) |
|-------|-------------------------------------------|----------------------------------|
| AC-F15-01 | CloudPlatform, CloudProvider, and topology-chain collection create, GET/LIST, and permitted PATCH behavior succeed with validation at their registry-declared scope subsets; the feature ends at InfrastructureStack. | VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-38, VS0-CF-F15-39 |
| AC-F15-02 | Participation lifecycle with provider request plus CloudPlatform acceptance guard is enforced; pair uniqueness is enforced. | VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-09 |
| AC-F15-03 | No ExecutionTarget identity/schema/route/status/writer/conformance/target-lifecycle behavior is introduced. | VS0-CF-F15-01, VS0-CF-F15-25 |
| AC-F15-04 | A cross-provider target reference is denied safely (no existence disclosure). | VS0-CF-X03, VS0-CF-F15-23 |
| AC-F15-05 | No migration plan/record/controller/cutover state machine exists. **Anti-drift/non-runtime verification gate** (excluded per ADH-2026-046 decision 3). | — |
| AC-F15-06 | PUT/DELETE return 405 for CloudPlatform, CloudProvider, and topology; PATCH accepts only `application/merge-patch+json`. | VS0-CF-F15-07 |
| AC-F15-07 | Every output passes the canonical stale-concept anti-drift gate. **Anti-drift/non-runtime verification gate** (excluded per ADH-2026-046 decision 3). | — |
| AC-F15-08 | Unauthorized writes to CloudPlatform, CloudProvider, or topology specs, status, and system-owned fields are denied; cross-provider references deny safely. | VS0-CF-F15-05, VS0-CF-F15-06, VS0-CF-F15-41, VS0-CF-X03 |
| AC-F15-09 | All Problem Details errors use FEATURE-0012 codes from the closed set; pure HTTP transport outcomes including 405 introduce no top-level Problem code. **Contract-reuse claim**, not a standalone runtime scenario. | Closed-set proof: requirements §4.28 exact conformance semantics ledger |
| AC-F15-10 | Independent holds behave correctly; suspend/resume denied in terminal/Pending/Terminating phases; clearing one hold does not reactivate while the other remains true. | VS0-CF-F15-08, VS0-CF-F15-10 |
| AC-F15-11 | Participation exposes no mutable spec fields; create requires Idempotency-Key only; existing-participation lifecycle changes require If-Match and Idempotency-Key. | VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-15, VS0-CF-F15-17, VS0-CF-F15-18, VS0-CF-F15-19, VS0-CF-F15-40 |
| AC-F15-12 | Scope-reference UID invariant enforced; mismatch rejected with VS0_SCOPE_REFERENCE_MISMATCH. | VS0-CF-F15-11 |
| AC-F15-13 | Lifecycle/hold-changing actions, resource create/PATCH, and security denials produce correlated, redacted AuditEvent with no secrets. | VS0-CF-F15-24 |
| AC-F15-14 | CloudPlatform-root requirement enforced with CONFLICT/409 VS0_CLOUDPLATFORM_ROOT_REQUIRED until at least one CloudPlatform exists. | VS0-CF-F15-12 |
| AC-F15-15 | Every collection-create kind enforces its closed field boundary; server-owned/status/scope/unknown/deferred fields rejected; each successful create returns 201 with exact initial status. | VS0-CF-F15-26 |
| AC-F15-16 | Every resource kind derives scope from exactly one server-side source; a client-supplied or mismatched scope is never honored. | VS0-CF-F15-27 |
| AC-F15-17 | Participation create accepts exactly the four decision-4 fields, initializes Pending/false holds/seven-day expiry, and rejects FEATURE-0021-owned fields; item actions accept only an empty JSON body plus required headers. | VS0-CF-F15-28 |
| AC-F15-18 | A syntactically valid but unassigned ISO-3166-1 alpha-2 value is rejected with VALIDATION_FAILED; administrativeAreaCode remains syntax-plus-country-prefix only. | VS0-CF-F15-29 |
| AC-F15-19 | A malformed/missing Idempotency-Key, an If-Match on participation create, and an ordinary non-empty participation item-action body return MALFORMED_REQUEST with no mutation/idempotency record/AuditEvent. The narrow valid duplicate-free top-level `bootstrapGrant` exception returns audited AUTHORIZATION_DENIED instead. | VS0-CF-F15-30, VS0-CF-F15-22 |
| AC-F15-20 | Missing or invalid authentication on any owned method/path pattern returns AUTH_REQUIRED with no mutation, lifecycle/status write, idempotency record, or AuditEvent. | VS0-CF-F15-31 |

### 7.3 Verification commands

make feature-0015-architecture-readiness
make feature-contract-check FEATURE=FEATURE-0015
make feature-0015-formal-check
make vs000-contract-check
make phase2r-drift-check
mkdocs build --strict
make structurizr-check
make fmt
make test
make vet
go test -race ./...
make ff-feature-gate FEATURE=FEATURE-0015

### 7.4 Reuse disposition

Disposition: Extend. FEATURE-0015 reuses FEATURE-0012 Problem Details, ObjectMeta, TypedRef, Condition, the validation pipeline, and the shared `ScopeKind` type; it extends the shared `ScopeKind` implementation to the approved canonical seven-value set. It reuses the FEATURE-0013 AuditEvent envelope and writer without changing their ownership. Retained alpha assets remain repository history and must compile, but no alpha route, generic `Provider` scope kind, or alpha runtime behavior is exposed by FEATURE-0015.

---

## 8. Implementation classification ledger

This ledger separates the three required dispositions.

### 8.1 IMPLEMENT (design-owned mechanics realized from approved requirements)

- Canonical resource types and JSON shapes for the seven owned kinds (§4.1).
- Shared ScopeKind canonicalization compatibility: update `internal/apimeta` and affected retained-alpha fixtures/tests to the canonical seven-value vocabulary; canonical F0015 paths reject retired `Provider` scope kind (DD-03).
- In-memory store: name uniqueness in server-derived scope; non-terminal participation pair uniqueness; resourceVersion; concurrency guard (§3.1, §7.1).
- Deterministic validation functions and per-route pipeline order (§5.1, §5.3).
- Participation lifecycle state machine, independent holds, derived phase, and deterministic scheduler expiry (§4.5).
- Closed create-contract classification and single-source scope derivation (§4.3, §4.4).
- Idempotency reservation-state coordinator (`InFlight`/`Completed`/`Aborted`) (DD-08).
- Bootstrap-grant resolver (server-resolved; header/body claims ignored) (DD-07, §6.1).
- Audit/mutation atomic coordinator and append-failure mapping (DD-09, §5.1 step 14).
- Assigned ISO-3166-1 alpha-2 dataset and lookup (DD-10).
- HTTP handlers and 35 explicit `http.ServeMux` registrations for 22 logical paths (DD-06, §4.2).

### 8.2 CONTRACT_ONLY/NO_TASK (consumed by reference; no FEATURE-0015 task changes them)

- FEATURE-0011 reuse governance contract (§7.4 reference).
- FEATURE-0012 Problem Details, ObjectMeta/TypedRef/ScopeRef, Condition, reference constraints, and validation pipeline (§3.2); ScopeKind is the explicit DD-03 compatibility implementation exception.
- FEATURE-0013 DecisionRecord/AuditEvent envelope, linkage, and writer `VS0-WRITER-011` (§3.2).
- Writer boundaries `VS0-WRITER-002`, `VS0-WRITER-003`, `VS0-WRITER-004`, `VS0-WRITER-005` (§6.1).
- REQ-F15-12 (existing-codes-only) and REQ-F15-13 (no-migration) are anti-drift/non-runtime obligations proven by construction and by absence; no runtime task realizes new behavior for them.
- The existing observability baseline (redacted logs/metrics), consumed as-is (§6.3).

### 8.3 EXCLUDED (must not implement; negative boundary only)

ExecutionTarget and its lifecycle/schema/routes/status/writer/conformance (FEATURE-0016); PolicyEvaluationRequest/Result and PolicyEngineAdapter (FEATURE-0017); GovernanceProfile/Membership/RoleDefinition/RoleAssignment (FEATURE-0018); SovereigntyProfile/SovereigntyFactSet/EvidenceRecord (FEATURE-0019); ProfileAssignment/EffectiveGovernanceContext (FEATURE-0020); CloudEnrollment/EntitlementPackage/ServiceEntitlement/QuotaPolicy, provider-selection intent, `spec.providerSelectionModes`, `spec.permittedHostingLocationRefs`, Organization visibility eligibility (FEATURE-0021); ServiceTypeDefinition/ServiceOffering/ServicePlan/ServiceRuntimeProfile/ServiceRequirementSet/ServiceRegion (FEATURE-0022); DecisionProfiles/ServicePlacement/PlacementDecision (FEATURE-0023); PluginExecution/ServiceDeploymentPlan/fake harness (FEATURE-0024); DecisionOperationContext/AI-readable projection (FEATURE-0025); the Slice 0 integration demo (FEATURE-0026); and all superseded/stale concepts enumerated in §10.

---

## 9. Requirement/decision/risk traceability

### 9.1 Architecture Traceability

- **Consumed decisions:** DEC-0037, DEC-0041, DEC-0042, DEC-0054, DEC-0059.
- **Consumed handoffs:** ADH-2026-020, ADH-2026-024, ADH-2026-025, ADH-2026-037, ADH-2026-042, ADH-2026-043 (non-migration portions), ADH-2026-045, ADH-2026-046, ADH-2026-047, ADH-2026-048, ADH-2026-049, ADH-2026-050, ADH-2026-051, ADH-2026-052, ADH-2026-053, ADH-2026-054, ADH-2026-055, ADH-2026-056, ADH-2026-057.
- **Schema IDs (owned):** `VS0-SCHEMA-008`, `VS0-SCHEMA-009`, `VS0-SCHEMA-010`, `VS0-SCHEMA-011`, `VS0-SCHEMA-012`, `VS0-SCHEMA-013`, `VS0-SCHEMA-014`.
- **Schema IDs (consumed by reference):** `VS0-SCHEMA-001`, `VS0-SCHEMA-002`, `VS0-SCHEMA-003`, `VS0-SCHEMA-004`, `VS0-SCHEMA-007`.
- **Writer IDs:** `VS0-WRITER-002`, `VS0-WRITER-003`, `VS0-WRITER-004`, `VS0-WRITER-005`, `VS0-WRITER-011` (FEATURE-0013).
- **State IDs:** `VS0-STATE-001`.
- **Error codes:** the FEATURE-0012 closed set enumerated in §5.2.

### 9.2 Conformance Mapping

Local proof cases used by this design (registry owner FEATURE-0015 unless noted): `VS0-CF-X03`, `VS0-CF-F15-01` through `VS0-CF-F15-41`. Each retains its exact registry semantics (requirements §4.28); this design cites but does not restate registry rows.

### 9.3 Risk / anti-drift controls

Architecture §11 anti-drift rules 1–13 are the design's risk controls and are satisfied by construction: stale-concept prohibition (§10); no ExecutionTarget activation (§1.3, §10); no provider-selection evaluation (§10); no platform lifecycle (§10); no migration reintroduction (§6.4); PATCH-only surface (§4.2, §6.4); no participation visibility eligibility (§6.2); no client-supplied scopeRef (§4.4); topology name immutability (§4.3); no FEATURE-0021 participation fields (§4.1, §10); assigned-ISO semantics with no new codes (DD-10, §5.2); exact audit/route boundary with local AUTH proof (§4.2, §5.1, §7.2); and exact name-uniqueness / schema-constraint proofs (§7.2).

---

## Canonical coverage ledger

Every approved REQ and AC appears exactly once below with its design disposition. No ID is created, renumbered, merged, split, reinterpreted, or omitted.

| ID | Disposition | Design coverage |
|----|-------------|-----------------|
| REQ-F15-01 | IMPLEMENT | CloudPlatform type, identity/uniqueness/scope/immutable-ownerRegistration validation (§4.1, §4.3, §5.1). |
| REQ-F15-02 | IMPLEMENT | CloudProvider type, identity/scope/operatingMarkets validation with assigned-code semantics (§4.1, §5.1, DD-10). |
| REQ-F15-03 | IMPLEMENT | CloudProviderParticipation type, pair uniqueness, and lifecycle store (§4.1, §4.5). |
| REQ-F15-04 | IMPLEMENT | Topology chain types, immutable parent references, containment and same-provider UID validation (§4.1, §5.1). |
| REQ-F15-05 | EXCLUDED (boundary) | Feature ends at InfrastructureStack; no target contract introduced (§1.3, §10). |
| REQ-F15-06 | IMPLEMENT | Lifecycle actions and independent suspension holds with derived phase (§4.5). |
| REQ-F15-07 | IMPLEMENT | No mutable participation spec; create needs Idempotency-Key only; item actions need both preconditions (§4.3, §5.1). |
| REQ-F15-08 | IMPLEMENT | PATCH-only merge-patch surface; PUT/DELETE return 405 (§4.2, §5.1). |
| REQ-F15-09 | IMPLEMENT | Safe cross-provider denial with no disclosure (§5.1 step 7). |
| REQ-F15-10 | IMPLEMENT/CONTRACT_ONLY | Seven-scope vocabulary consumed; per-kind scope subsets validated (DD-03, §4.4). |
| REQ-F15-11 | IMPLEMENT | Writer enforcement for CloudPlatform, CloudProvider, and topology mutable spec by reference to registered writers; wrong-administrator PATCH is denied before semantic processing/publication (§6.1; VS0-CF-F15-41). |
| REQ-F15-12 | CONTRACT_ONLY | Existing FEATURE-0012 codes only; anti-drift obligation proven by construction (§5.2). |
| REQ-F15-13 | CONTRACT_ONLY | No migration mechanism; retired tombstones never reused; proven by absence (§6.4). |
| REQ-F15-14 | IMPLEMENT | cloudPlatformRef.uid = scopeRef.uid invariant enforced (§4.4, §5.1 step 10). |
| REQ-F15-15 | IMPLEMENT | Audit categories, atomic publication, and append-failure mapping to INTERNAL_ERROR (§5.1 step 14, §6.2, DD-09). |
| REQ-F15-16 | IMPLEMENT | CloudPlatform-root requirement with audited denial (§5.1 step 6). |
| REQ-F15-17 | IMPLEMENT | Server-resolved deterministic grants; no persisted roles/memberships/assignments (§6.1, DD-07). |
| REQ-F15-18 | IMPLEMENT | Idempotency replay and changed-digest conflict across all creates/actions (§5.1 step 9, DD-08). |
| REQ-F15-19 | IMPLEMENT | Closed collection-create contract with 201/200 outcomes (§4.3, §5.1 step 8). |
| REQ-F15-20 | IMPLEMENT | Single-source server-side scope derivation; no client scopeRef (§4.4). |
| REQ-F15-21 | IMPLEMENT | Participation create body and empty item-action body; FEATURE-0021 fields rejected (§4.1, §4.3, §10). |
| REQ-F15-22 | IMPLEMENT | Assigned ISO-3166-1 alpha-2 dataset; administrativeAreaCode syntax-plus-prefix only (DD-10, §5.1 step 11). |
| REQ-F15-23 | IMPLEMENT | Four malformed/prohibited-input outcomes using existing codes (§5.1 steps 3–4). |
| REQ-F15-24 | IMPLEMENT | AUTH_REQUIRED on every owned pattern with no side effect (§5.1 step 1). |
| AC-F15-01 | IMPLEMENT | Owned-kind create/GET/LIST/PATCH with scope subsets; ends at InfrastructureStack (§4.1, §4.2, §5.1). The separate method-boundary acceptance criterion covers PUT/DELETE 405. |
| AC-F15-02 | IMPLEMENT | Participation lifecycle, acceptance guard, and pair uniqueness (§4.5, §5.1 step 13). |
| AC-F15-03 | EXCLUDED (boundary) | No ExecutionTarget behavior introduced (§1.3, §10). |
| AC-F15-04 | IMPLEMENT | Cross-provider reference denied safely (§5.1 step 7). |
| AC-F15-05 | CONTRACT_ONLY | Anti-drift/non-runtime gate: no migration mechanism exists (§6.4). |
| AC-F15-06 | IMPLEMENT | PUT/DELETE 405; PATCH accepts only merge-patch media type (§4.2, §5.1 step 2). |
| AC-F15-07 | CONTRACT_ONLY | Anti-drift/non-runtime stale-concept gate satisfied by construction (§10). |
| AC-F15-08 | IMPLEMENT | Denials for wrong-administrator mutable-spec, status, and system-owned writes plus safe cross-provider denial (§5.1 steps 5 and 7, §6.1; VS0-CF-F15-41). |
| AC-F15-09 | CONTRACT_ONLY | All errors use existing FEATURE-0012 Problem Details codes (§5.2). |
| AC-F15-10 | IMPLEMENT | Independent holds; suspend/resume denied in terminal/Pending/Terminating (§4.5, §5.1 step 13). |
| AC-F15-11 | IMPLEMENT | No mutable participation spec; header-precondition rules plus replay-response behavior for create versus actions (§4.3, §5.1 steps 3 and 9). |
| AC-F15-12 | IMPLEMENT | Scope-reference UID invariant enforced (§4.4, §5.1 step 10). |
| AC-F15-13 | IMPLEMENT | Correlated redacted AuditEvents for actions, create/PATCH, and security denials (§5.1 step 14, §6.2). |
| AC-F15-14 | IMPLEMENT | Root requirement enforced until a CloudPlatform exists (§5.1 step 6). |
| AC-F15-15 | IMPLEMENT | Closed field boundary per create kind; server/status/scope/unknown/deferred rejected; 201 with exact status (§4.3). |
| AC-F15-16 | IMPLEMENT | Exactly one server-side scope source per kind; client scope never honored (§4.4). |
| AC-F15-17 | IMPLEMENT | Participation create four-field body with Pending/false-holds/seven-day expiry; empty item-action bodies (§4.3, §4.5). |
| AC-F15-18 | IMPLEMENT | Unassigned ISO value rejected; administrativeAreaCode syntax-plus-prefix only (DD-10, §5.1 step 11). |
| AC-F15-19 | IMPLEMENT | Malformed Idempotency-Key, If-Match on create, and ordinary non-empty action body return MALFORMED_REQUEST; the valid duplicate-free bootstrapGrant exception is audited AUTHORIZATION_DENIED (§5.1 steps 3–4). |
| AC-F15-20 | IMPLEMENT | Missing/invalid auth on any owned pattern returns AUTH_REQUIRED with no side effect (§5.1 step 1). |

---

## 10. Non-goals, absence ledger, and unresolved report

### 10.1 Prohibited stale concepts (named only to exclude; never active behavior)

FEATURE-0015 does not implement, activate, or reference as active behavior any superseded concept. `ResourcePool` (superseded by DEC-0042) must not appear as active behavior; `ProviderCapability` (superseded by DEC-0042) must not appear as active behavior; a generic combined Provider (superseded by DEC-0037) must not be used; `ServiceClass` as a canonical catalog (superseded by DEC-0049) must not be used; `EffectivePolicyContext` (superseded by DEC-0050) must not be used; the six-scope vocabulary (superseded by DEC-0037) must not be used; `SovrunnInstallation` as an active Slice 0 resource (DEC-0053 deferred; `VS0-SCHEMA-057` tombstone) must not be used; and `CanonicalMigrationPlan`/`CanonicalMigrationRecord`/migration controller/cutover state machine (DEC-0059; ADH-2026-045; retired tombstones `VS0-SCHEMA-060`, `VS0-SCHEMA-061`, `VS0-WRITER-020`, `VS0-WRITER-021`, `VS0-STATE-011`) must never be reused or reintroduced.

### 10.2 Adjacent-feature exclusions (absence ledger)

| Excluded element | Owner | Absence guarantee |
|------------------|-------|-------------------|
| ExecutionTarget identity/schema/routes/status/writer/conformance/target lifecycle; NormalizedTargetFactSet; TargetQualificationResult (`VS0-SCHEMA-015`, `VS0-SCHEMA-016`, `VS0-SCHEMA-017`) | FEATURE-0016 | Not introduced; proven by VS0-CF-F15-01 and VS0-CF-F15-25 |
| PolicyEvaluationRequest/Result, PolicyEngineAdapter | FEATURE-0017 | Not referenced; route/scope surface bounded by VS0-CF-F15-02, VS0-CF-F15-25 |
| GovernanceProfile, Membership, RoleDefinition, RoleAssignment | FEATURE-0018 | No roles/memberships/assignments persisted; local AUTH proof VS0-CF-F15-31 (not inherited VS0-CF-F01) |
| SovereigntyProfile, SovereigntyFactSet, EvidenceRecord | FEATURE-0019 | Not introduced |
| ProfileAssignment, EffectiveGovernanceContext | FEATURE-0020 | Not introduced |
| CloudEnrollment, EntitlementPackage, ServiceEntitlement, QuotaPolicy, provider-selection intent, `spec.providerSelectionModes`, `spec.permittedHostingLocationRefs`, Organization visibility eligibility | FEATURE-0021 | Provider-selection fields rejected on participation create (VS0-CF-F15-28); eligibility not stored |
| ServiceTypeDefinition, ServiceOffering, ServicePlan, ServiceRuntimeProfile, ServiceRequirementSet, ServiceRegion | FEATURE-0022 | Not introduced |
| DecisionProfiles, ServicePlacement, PlacementDecision | FEATURE-0023 | Not introduced; VS0-CF-F09 is an integration reference only |
| PluginExecution, ServiceDeploymentPlan, fake harness | FEATURE-0024 | Not introduced |
| DecisionOperationContext, AI-readable projection | FEATURE-0025 | Not introduced |
| Slice 0 integration demo, cross-feature orchestration | FEATURE-0026 | Not introduced; VS0-CF-HP01 and VS0-CF-T01 are integration references only |
| Federation, remote deployment, signed remote offer, cross-control-plane replication; platform lifecycle | No Phase 2R owner | Explicitly excluded (ADH-2026-045 decision 12; DEC-0053 deferred) |
| Real provisioning | Phase 3 | Excluded per PHASE2R non-goals |

FEATURE-0014's alpha model is a retained repository asset assessed for reuse under FEATURE-0011; it is not runtime-migrated. The explicit stale-concept exclusions are in §10.1.

### 10.3 Unresolved report

All design questions delegated by architecture (§7.12) and requirements (§9) are resolved in §2 and §5. Every approved REQ and AC maps to a design disposition (Canonical coverage ledger). No semantic or ownership choice remains open; every scope-kind, writer, state, schema, and error reference resolves to an existing registry entry or accepted decision. No `ARCHITECTURE_DECISION_REQUIRED`, `REQUIREMENT_CLARIFICATION_REQUIRED`, `BOUNDARY_CHANGE_REQUIRED`, `DEPENDENCY_APPROVAL_REQUIRED`, or `SECURITY_REVIEW_REQUIRED` condition was triggered.

---

## Model Execution Report:
- Tool: kiro
- Stage or task: FEATURE-0015 Design (`.kiro/specs/canonical-cloud-model-and-alpha-migration/design.md`)
- Recommended priority list: 1) System design (`claude-opus-4.8`), effort: high; 2) Fallback (`claude-sonnet-4.5`), effort: medium
- Selected model: claude-opus-4.8
- Effort/reasoning setting: high
- Fallback used: no
- Fallback reason: none

STAGE_STATUS: COMPLETE
