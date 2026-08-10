# VS-000 Slice 0 Exact Contract Specification

| Field | Value |
|-------|-------|
| Status | Experimental |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Controlling Adoption | ADH-2026-042 |
| Machine-Readable Registry | `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` |

## Scope Statement

This document is the normative specification for the exact minimal synthetic Slice 0 integration profile only.

It is **not**: a full schema for FEATURE-0015 through FEATURE-0026; not a feature requirements, design, or tasks artifact; not a future slice definition.

The machine-readable normative registry (`VS-000-contract-registry.yaml`) is the authoritative source for field-level detail. This specification defines the rules that govern that registry.

FEATURE-0012 (common grammar, Problem Details) and FEATURE-0013 (DecisionRecord, AuditEvent) schemas are reused by reference. Parallel envelopes for the same concerns are forbidden.

---

## 1. Authority and Precedence

1. Canonical data model (`docs/architecture/canonical/sovrunn-finalized-data-model.md`) is semantic owner of all resource kinds.
2. Canonical contract catalog (`docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`) is semantic owner of all API contracts.
3. This specification defines the Slice 0 narrowing profile; it cannot broaden canonical semantics.
4. The registry YAML is the machine-readable expression of this specification.
5. Precedence: active baseline > canonical model/catalog > accepted DEC/ADH records > VS-000 charter > this specification > registry YAML > approved feature stage > implementation code.
6. Seven canonical scope kinds apply: Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider.

---

## 2. Schema-Profile Rules

Every kind or DTO exchanged or persisted within Slice 0 must have a registry entry containing:

| Attribute | Requirement |
|-----------|-------------|
| Registry ID | Unique `VS0-SCHEMA-<NNN>` identifier |
| Identity | Primary key fields and uniqueness constraint |
| Profile | Slice 0 narrowing profile name |
| Boundary | Which layer owns creation/mutation |
| Scope | One of seven canonical scope kinds |
| Fields | Name, type, requiredness, enums, defaults, limits |
| References | Foreign-key relationships to other registry entries |
| Classification/Redaction | Sensitivity label; redaction rule for projection |
| Mutability/Retention | Immutable-after-create, mutable fields, retention class |
| Validation | Deterministic validation rule references |
| Projection | Safe customer-facing projection subset |

FEATURE-0012 common grammar (request/response envelope, pagination, error shape) and FEATURE-0013 DecisionRecord/AuditEvent schemas are incorporated by reference to their canonical definitions. No Slice 0 entry may redefine them.

Slice 0 narrowing may restrict cardinality, constrain enums, or tighten limits. It must not add fields, relax constraints, or introduce semantics absent from the canonical model.

---

## 3. Exact Writer Rules

| Rule ID | Rule |
|---------|------|
| VS0-W01 | Every mutable field has exactly one authoritative writer domain. |
| VS0-W02 | Clients (Portal, CLI, API consumers) never write `status`. |
| VS0-W03 | The deployment planner records an already-authorized target set; it never selects targets, executes steps, or changes decisions. |
| VS0-W04 | Decision services never mutate their inputs; they emit DecisionRecords. |
| VS0-W05 | Adapters write only normalized observations or execution results. |
| VS0-W06 | Projection endpoints emit read-only safe subsets; no write path. |
| VS0-W07 | Binding writers validate consumer entitlement before credential issuance. |
| VS0-W08 | Conflicts between concurrent writers fail closed (reject, do not merge). |

Writer domain table (normative rows in registry under `VS0-WRITER-<NNN>`):

| Writer Domain | Writes To | Never Writes |
|---------------|-----------|--------------|
| Client | `spec` of owned resources | `status`, DecisionRecord |
| Controller | `status`, Operation state | `spec` |
| Deployment planner | Immutable plan from final decisions and pinned inputs | Target selection, execution results, credentials |
| Decision Service | DecisionRecord | Input resources, `status` |
| Adapter | Observations, execution results | `spec`, DecisionRecord |
| Projection controller | Append-only ServicePlacement and transient explanation projection | Canonical decisions, provider topology, secret values |

---

## 4. Exact State-Machine Rules

1. Registry transitions (under `VS0-STATE-<NNN>`) are authoritative for allowed state changes and guards.
2. Milestones are observable checkpoints; they are not lifecycle phases.
3. Immutable Slice 0 records are persisted directly as `FINAL` and are append-only thereafter; `Created` is not a persisted record state.
4. Operation owns `Pending → Running → Succeeded | Failed | Cancelled`; pending, retrying and waiting are never DecisionRecord states.
5. Deletion of a ServiceInstance revokes all associated ServiceBindings before completing.
6. Deletion retains immutable audit/decision records; it does not purge history.
7. Invalid transitions are rejected; the system does not silently coerce state.

State-machine summary (exact definitions in registry):

| Kind Group | States | Terminal | Registry Prefix |
|------------|--------|----------|-----------------|
| Participation and enrollment | Pending, Active, Suspended, Terminating, Terminated | Terminated | VS0-STATE-001..002 |
| Published definitions | Draft, Published, Retired | Retired | VS0-STATE-003 |
| Target qualification/availability | Orthogonal qualified/available combinations in the registry | None | VS0-STATE-004 |
| Quota reservation | None, Reserved, Committed, Released | Released | VS0-STATE-005 |
| ServiceInstance | Pending, Provisioning, Ready, Failed, Deleting, Deleted | Deleted | VS0-STATE-006 |
| ServiceBinding | Pending, Bound, Revoking, Revoked, Failed | Revoked | VS0-STATE-007 |
| Operation | Pending, Running, Succeeded, Failed, Cancelled | Succeeded, Failed, Cancelled | VS0-STATE-008 |
| PluginExecution | Pending, Running, Succeeded, Failed | Succeeded, Failed | VS0-STATE-009 |
| Immutable records | FINAL | FINAL | VS0-STATE-010 |

---

## 5. Errors

### RFC 9457 Problem Details Members

All error responses use RFC 9457 Problem Details with the FEATURE-0012 Sovrunn extensions. Required contract members are `type`, `title`, `status`, `detail`, `instance`, `code` and `requestId`; `violations` is present when field-level details apply.

### Top-Level Error Code Mapping

Exact closed set (normative rows in the registry under `problemCodes`):

| Codes | HTTP |
|-------|------|
| `MALFORMED_REQUEST`, `UNKNOWN_FIELD`, `DUPLICATE_FIELD`, `REQUEST_TOO_LARGE` | 400 |
| `AUTH_REQUIRED` | 401 |
| `AUTHORIZATION_DENIED` | 403 |
| `RESOURCE_NOT_FOUND` | 404 |
| `CONFLICT`, `ALREADY_EXISTS`, `DELETE_BLOCKED` | 409 |
| `STALE_RESOURCE_VERSION` | 412 |
| `UNSUPPORTED_MEDIA_TYPE` | 415 |
| `VALIDATION_FAILED` | 422 |
| `INTERNAL_ERROR` | 500 |
| `DEPENDENCY_UNAVAILABLE` | 503 |

For every code, `type` is exactly `urn:sovrunn:problem:<lowercase-kebab-code>`. Slice 0 does not add a top-level code or HTTP mapping. Quota exhaustion is a `CONFLICT`/409 with violation `VS0_QUOTA_EXHAUSTED`; inaccessible cross-scope resources use safe `RESOURCE_NOT_FOUND`/404.

### Validation Detail Violation Codes

Validation details use only FEATURE-0012 shared violations, the inherited FEATURE-0013 decision violations, or the closed Slice 0 entries under registry `violationCodes.slice0`. A violation carries its exact RFC 6901 JSON Pointer or `null` when a safe field path cannot be disclosed.

### Inherited Schemas

FEATURE-0013 registry defines AuditEvent and DecisionRecord error semantics. They are inherited, not duplicated.

### Safety Rules

- Error responses must not confirm existence of resources the caller cannot access (safe denial).
- No error `detail` or `title` string is parsed programmatically by clients; codes and types are the contract.

---

## 6. Conformance

### Conformance Classes

Each conformance requirement has a stable registry identifier. `VS0-CF-HP01` owns the end-to-end happy path; `VS0-CF-F01` through `VS0-CF-F20` map one-to-one to the charter failures; `X`, `L`, `Z`, `T`, `I` and `D` cases cover scope, leakage, external effects, traceability, concurrency and deletion ordering. Feature-local pre-integration evidence uses `VS0-CF-F<feature>-<case>` and must not reuse a downstream-owned scenario. FEATURE-0015 migration proof uses `VS0-CF-MIG01..MIG02` plus the exact negative cases `VS0-CF-MIGF01..MIGF03`.

| Registry range | Class | Verification |
|----------------|-------|--------------|
| `VS0-CF-HP01` | Positive integration | All 20 charter steps complete deterministically. |
| `VS0-CF-F01..F20` | Required failures | Exact state, Problem/status code and side-effect assertions for every charter failure. |
| `VS0-CF-X01..X03` | Scope isolation | Cross-Organization, cross-Project and cross-provider references deny safely. |
| `VS0-CF-L01` | Leakage prevention | Responses, logs, audit and explanation contain no prohibited data. |
| `VS0-CF-Z01` | Zero external effect | Fake execution records `externalCallCount=0`. |
| `VS0-CF-T01` | Correlation | Request, decisions, plan, operation, executions and audit form one complete trace graph. |
| `VS0-CF-I01..I02` | Idempotency and races | Same payload converges; different payload conflicts; one quota/operation wins. |
| `VS0-CF-D01` | Deletion ordering | Binding revocation precedes cleanup, quota release and instance finalization. |
| `VS0-CF-F15-01..F15-10` | FEATURE-0015 local evidence | Canonical resource/scope validation, participation acceptance and uniqueness, writer denial and separation, plan immutability, and signed-backup/verified-restore gating. |
| `VS0-CF-MIG01..MIG02` | Migration positive evidence | Signed plan, signed backup evidence, verified restore evidence, and append-only records deterministically reach the `Completed` milestone without dual authority. |
| `VS0-CF-MIGF01..MIGF03` | Migration failure evidence | Invalid milestone order, dual authority and unresolved references fail with their registered Problem/violation mapping and zero cutover side effects. |

Each conformance test specifies exact expected state, expected error (code + HTTP + type), and expected side effects.

---

## 7. Kiro Anti-Drift Contract

When implementing any Slice 0 feature stage, Kiro must:

1. Load: canonical model, canonical catalog, VS-000 charter, this specification, registry YAML, Slice 0 traceability matrix, FEATURE traceability and DECISION_INDEX.
2. Implement one feature at a time, one stage at a time.
3. Never bypass feature-stage approval gates.
4. Never use stale/prohibited concepts (see §9).
5. Never invent schema fields, state transitions, or error codes not in the registry.
6. Any change to schema, state machine, or error contract must update registry, traceability, and conformance tests in the same change.
7. If a genuine semantic gap is encountered (question cannot be answered by loaded authorities), stop and emit `ARCHITECTURE_DECISION_REQUIRED` with the gap description. Do not guess.

Violation of any rule above constitutes architecture drift and blocks merge.

---

## 8. Repository Traceability and Validation

### Traceability Paths

```text
docs/architecture/vertical-slices/VS-000-contract-specification.md  (this file)
docs/architecture/vertical-slices/VS-000-contract-registry.yaml     (machine-readable registry)
docs/architecture/vertical-slices/VS-000-core-skeleton.md            (skeleton design)
docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md             (Slice 0 contract linkage)
docs/traceability/FEATURE_TRACEABILITY_MATRIX.md                     (feature linkage)
docs/traceability/DECISION_TRACEABILITY_MATRIX.md                    (decision linkage)
docs/decisions/DECISION_INDEX.md                                     (accepted decisions)
```

### Validation Target

```bash
make vs000-contract-check
```

This target validates:
- Registry YAML parses without error.
- All `VS0-SCHEMA-*` entries have the required profile attributes.
- All `VS0-STATE-*` transitions reference declared states.
- Every Problem code matches the FEATURE-0012 closed set and HTTP/type mapping.
- Every `VS0-F01..F20` maps to exactly one `VS0-CF-F*` case.
- No prohibited terms appear in active registry entries.
- No field exists in implementation that is absent from registry.

---

## 9. Explicit Non-Goals

This specification does not:

- Define new resource kinds or API routes.
- Serve as requirements, design, or tasks for any FEATURE.
- Specify future slices (VS-001+).
- Replace the canonical model or catalog as semantic owner.
- Define runtime behavior beyond synthetic Slice 0 scope.
- Mandate persistent storage, Kubernetes CRDs, or real provisioning.
- Prescribe UI, billing, marketplace, or multi-cluster behavior.

### Prohibited Terms (stale/superseded — must not appear in active Slice 0 artifacts)

| Term | Superseded By | Reference |
|------|---------------|-----------|
| ResourcePool | (removed) | DEC-0042 |
| ProviderCapability | (removed) | DEC-0042 |
| Generic Provider (combined) | CloudPlatform + CloudProvider | DEC-0037 |
| ServiceClass (as canonical catalog) | ServiceTypeDefinition + ServiceOffering | DEC-0049 |
| EffectivePolicyContext | EffectiveGovernanceContext | DEC-0050 |
| Six-scope vocabulary | Seven canonical scopes | DEC-0037 |
| SovrunnInstallation (active Slice 0) | CloudProviderParticipation boundary proof | No Phase 2R owner; DEC-0053 deferred |

---

## 10. Architecture Decision Required Rule

Ordinary mechanics (CRUD, state transitions, validation, projection, error mapping) must be exact and fully specified in the registry. Implementation must not require interpretation.

When a genuine semantic gap exists — a question that cannot be resolved by the canonical model, catalog, accepted decisions, or this specification — the following applies:

1. The gap is recorded in registry `architectureDecisionRequired` with a stable ID, question and affected features.
2. Downstream code generation, test generation, and implementation for the affected scope stops.
3. An `ARCHITECTURE_DECISION_REQUIRED` marker is emitted.
4. Resolution requires a new DEC, RFC, or ADH through the Architecture Operating System.
5. After resolution, the registry entry is updated, conformance tests are added, and implementation may proceed.

No agent may fill a semantic gap with assumptions, defaults, or analogies from other platforms.

---

*End of specification.*
