# Boundary Ledger

<!-- Generated from docs/api/boundary-ledger.yaml. Do not edit by hand. -->

Machine-readable source of truth: `docs/api/boundary-ledger.yaml`.
This Markdown file is a regenerable human view (D-12, F12-LEDGER-001).

## Document

| Field | Value |
| --- | --- |
| apiVersion | sovrunn.io/v1alpha1 |
| kind | BoundaryLedger |
| name | feature-0012-matrix-c1 |
| feature | FEATURE-0012 |

## Description

Trust and API boundary ledger for FEATURE-0012 canonical schemas. Records purpose, ownership, producers/consumers, data controls, authorization, audit, observability, failure behavior, versioning, and evolution paths for each Matrix C1 boundary.

## Boundaries

### customer-facing

#### Purpose

Expose product intent, safe observed status, and actionable errors to tenant-facing clients without revealing provider internals or inaccessible tenant data.

#### Owner

Sovrunn Platform API Architecture

#### Producers

- api/schemas/project.json
- Sovrunn core control-plane (desired-state validation and status ownership)

#### Consumers

- Portal
- CLI
- SDK
- GitOps
- Tenant automation

#### Allowed data

- Product intent and customer-authored desired state (spec)
- Safe customer-visible status and conditions
- Actionable RFC 9457 problems with stable codes and JSON Pointers
- Typed references using Sovrunn kinds and scopes only

#### Prohibited data

- Provider-native identifiers, SDK types, and substrate handles
- Raw secrets, credentials, tokens, private keys, and connection strings
- Existence or detail of inaccessible cross-tenant objects
- Operator-confidential diagnostics not authorized for the tenant

#### Authorization

Scope-aware authorization against the resource scopeRef; cross-scope denial uses SafeDenial (404 RESOURCE_NOT_FOUND) without existence disclosure. Customer mutation modes reject system-owned and status fields.

#### Audit

Mutating customer-facing operations must be auditable with actor, action, resource identity, scope, request_id, and outcome. Audit payloads must not contain raw secrets.

#### Observability

Structured logs and traces carry request_id and, when applicable, operation_id. Latency and error_code may be recorded. Secrets, credentials, tokens, private keys, and connection strings are never logged.

#### Failure behavior

Validation and decode failures return stable problem codes with JSON Pointer field paths. Authorization denials that must not disclose existence return byte-identical 404 SafeDenial responses. Oversized or unsupported media inputs fail closed at the boundary.

#### Versioning

Contracts version via apiVersion group/maturity (for example core.sovrunn.io/v1alpha1). Breaking schema changes require a new version or approved baseline evidence under the schema-diff gate.

#### Replacement path

Replace a customer-facing contract by introducing a new versioned schema and TypeBinding, migrating clients deliberately, and retiring the prior version through the maturity and compatibility policy.

#### Migration path

Phase 1 customer resources coexist unchanged; migration candidates are recorded in the Phase 1 compatibility report and executed only through separately approved features.

#### Reassessment trigger

Reassess before stable API promotion, cross-organization sharing, a new privileged customer view, regulated/classified customer workloads, or any object that cannot select an approved profile at this boundary.


### operator-facing

#### Purpose

Provide administrative and normalized infrastructure contracts for platform, MSP, and provider operators without exposing raw secrets or unrestricted customer data.

#### Owner

Sovrunn Platform API Architecture

#### Producers

- api/schemas/resource-pool.json
- Sovrunn operator control surfaces and pool lifecycle controllers

#### Consumers

- Platform operators
- MSP operators
- Provider operators

#### Allowed data

- Administrative desired state and normalized infrastructure descriptors
- Operator-visible diagnostics and pool readiness conditions
- Provider-neutral capability class and jurisdiction metadata

#### Prohibited data

- Raw secrets and unrestricted customer confidential payloads
- Provider SDK types as shared core models
- Customer-facing product contracts served from this boundary

#### Authorization

Operator principals are authorized within Provider (or declared operator) scope. Access to customer-confidential data requires an explicit approved boundary view; SafeDenial applies for unauthorized cross-scope access.

#### Audit

Operator mutations and privileged reads that affect pools or operator state must emit audit records with actor, scope, request_id, and outcome without logging secret material.

#### Observability

Operator-facing operations propagate request_id and operation_id where applicable. Diagnostics may include normalized error codes; raw credentials and connection strings remain redacted.

#### Failure behavior

Fail closed on unknown fields, duplicate keys, unsupported schema keywords, and missing required stages. Stale concurrent writes return 412 STALE_RESOURCE_VERSION when If-Match protection is required.

#### Versioning

Operator contracts follow the same apiVersion maturity ladder as other FEATURE-0012 schemas; breaking changes are gated by baseline approval evidence.

#### Replacement path

Introduce a replacement operator schema/version, update operator tooling bindings, and remove the prior contract only after compatibility review.

#### Migration path

Existing Phase 1 operator-adjacent resources remain; operator-facing grammar adoption proceeds via explicit compatibility exceptions and later approved migrations.

#### Reassessment trigger

Reassess at the first non-Kubernetes provider integration, the first real provider adapter exposing new operator views, or when operator diagnostics would otherwise require unrestricted customer data.


### internal-engine-facing

#### Purpose

Carry normalized internal decision and evaluation contracts for policy, placement, entitlement, and orchestration engines without embedding vendor SDK types as shared models.

#### Owner

Sovrunn Platform API Architecture

#### Producers

- api/schemas/placement-evaluation-request.json
- Internal engines that emit TransientRequestResult profiles

#### Consumers

- Policy evaluation
- Placement evaluation
- Entitlement evaluation
- Orchestration engines

#### Allowed data

- Normalized internal request/result contracts
- Explainable decision inputs using Sovrunn kinds and scopes
- Finite evaluation metadata required for engine coordination

#### Prohibited data

- Vendor SDK types as shared platform models
- Customer portal schemas reused as engine contracts
- Raw secrets in evaluation payloads

#### Authorization

Internal-engine contracts are not a customer or operator public surface. Callers must be platform-authorized engine components; cross-tenant evaluation must still respect scope equality and SafeDenial rules for unauthorized targets.

#### Audit

Engine evaluations that affect placement or entitlement decisions should correlate via request_id/operation_id and leave an auditable decision trail in later FEATURE-0013 records without storing secrets.

#### Observability

Engine-facing calls use structured logs with request_id and operation_id correlation. Evaluation latency and stable error codes are permitted; secret-bearing fields are never logged.

#### Failure behavior

Structural or semantic validation failures stop the pipeline at the failing layer with stable codes. Misconfigured authorization layers fail closed with 500 INTERNAL_ERROR and no silent skip.

#### Versioning

Internal-engine schemas are versioned independently of customer-facing contracts so engine evolution does not force portal/SDK breakage.

#### Replacement path

Replace an engine contract by publishing a new versioned schema, updating engine adapters, and retiring the prior contract after dual-read or dual-write migration as approved.

#### Migration path

FEATURE-0012 defines grammar only; live placement/policy execution migrates in later features while preserving these contract boundaries.

#### Reassessment trigger

Reassess when an object cannot select an approved profile, when an extension becomes required by core decisions or multiple engines, or before high-frequency status updates that stress evaluation contracts.


### adapter-facing

#### Purpose

Isolate external-system translation contracts, provider handles, and observation provenance so provider-native data cannot leak into customer or core schemas.

#### Owner

Sovrunn Platform API Architecture

#### Producers

- api/schemas/discovered-database.json
- api/schemas/adapter-configuration.json
- External-system adapters (future FEATURE-0016+)

#### Consumers

- External-system adapters
- Adapter configuration and discovery controllers

#### Allowed data

- Translation contracts between Sovrunn models and external systems
- Provider handles and adapter-local identifiers
- Provenance, observed time, and freshness for external observations

#### Prohibited data

- Leakage of provider-native fields into customer-facing or core schemas
- Treating adapter handles as customer product identity
- Unredacted runtime credentials in metadata, status, or errors

#### Authorization

Adapter-facing resources are authorized in Provider scope for adapter and operator principals. Customer callers must not consume adapter-native contracts directly.

#### Audit

Adapter configuration changes and discovery ingest events must be auditable with actor, provider scope, request_id, and outcome. Secret values remain externalized via secret references only.

#### Observability

Adapter operations correlate with request_id/operation_id. Provenance and freshness fields support accuracy; credentials and connection strings are never written to logs.

#### Failure behavior

Missing provenance/freshness on observed external resources fails validation. Stale observations must surface as degraded/unknown rather than silent success. Unauthorized customer access to adapter objects uses SafeDenial.

#### Versioning

Adapter contracts version separately from customer APIs so provider replacement does not require customer-contract changes.

#### Replacement path

Replace an adapter contract by introducing a new adapter schema version and migrating the adapter implementation behind the same customer/core facade.

#### Migration path

Provider-native data stays behind this boundary; core/customer migration never imports adapter SDK types. Later adapter features adopt these schemas without rewriting customer contracts.

#### Reassessment trigger

Reassess at the first real provider adapter, the first external discovery source, the first disconnected/federated control plane, or any proposal to expose provider-native identifiers on customer-facing contracts.


### plugin-facing

#### Purpose

Define plugin capability catalogs and long-running operation contracts for the plugin manager and implementations without granting unrestricted control-plane access or policy bypass.

#### Owner

Sovrunn Platform API Architecture

#### Producers

- api/schemas/plugin-definition.json
- api/schemas/operation.json
- Plugin manager and validated plugin result publishers

#### Consumers

- Plugin manager
- Plugin implementations
- Operation lifecycle observers authorized for the operation scope

#### Allowed data

- Versioned plugin capability definitions
- Operation targetRef/scopeRef contracts with six allowed scopes
- Validated plugin results expressed through Sovrunn status grammar

#### Prohibited data

- Unrestricted control-plane access or authorization bypass
- Policy-engine substitution through plugin payloads
- Provider SDK types smuggled as plugin-facing core models

#### Authorization

Plugin-facing operations require an authorized caller for the operation scope. Operation.targetRef governance scope must equal Operation.scopeRef (D-17); unauthorized or unavailable targets use SafeDenial without disclosing mismatch details.

#### Audit

Plugin definition publication and operation lifecycle transitions must emit audit records with actor, plugin/operation identity, scope, request_id, operation_id, and outcome. Result messages must not embed secrets.

#### Observability

Operation and plugin flows propagate request_id and operation_id. Structured status/conditions communicate progress; secret material is excluded from logs and problem details.

#### Failure behavior

Scope/target mismatch yields OPERATION_TARGET_SCOPE_MISMATCH at /metadata/scopeRef when the target is authorized and available. Absent or unauthorized targets return identical 404 SafeDenial. Plugin results cannot redefine core status grammar.

#### Versioning

PluginDefinition and Operation schemas use explicit apiVersion maturity. Capability additions that break consumers require review-classified or breaking change handling under the schema-diff gate.

#### Replacement path

Replace plugin-facing contracts by versioning PluginDefinition/Operation, dual-supporting during migration, then removing the prior version after plugin-manager readiness.

#### Migration path

FEATURE-0012 ships grammar and fixtures only; plugin execution arrives in later features while preserving these contracts and six Operation scopes.

#### Reassessment trigger

Reassess at the first remotely executed or data-path plugin, before plugin taxonomy expansion that changes trust assumptions, or when a plugin result would need to redefine core status ownership.


### governance-only

#### Purpose

Record immutable decisions, assessments, approvals, and traceability for architecture and review workflows without carrying runtime credentials or customer secrets.

#### Owner

Sovrunn Platform API Architecture

#### Producers

- api/schemas/audit-event.json
- Architecture and governance review workflows
- Future decision/audit record publishers (FEATURE-0013+)

#### Consumers

- Architecture review
- Compliance and audit review
- Governance workflow tooling

#### Allowed data

- Immutable audit and decision records
- Assessment and approval evidence
- Traceability identifiers (request_id, operation_id, actor, scope)

#### Prohibited data

- Runtime credentials and customer secrets
- Mutable desired-state specs treated as governance truth
- Provider-native payloads as governance evidence without redaction

#### Authorization

Governance-only records are readable by authorized governance and audit roles within the declared Organization (or tighter) scope. Writers are restricted to approved governance producers; normal customers cannot mutate immutable records.

#### Audit

The ledger entry itself is an audit boundary: append-only immutable records capture who decided what, with correlation IDs. Nested secret values are forbidden; references only.

#### Observability

Governance ingest correlates via request_id/operation_id when produced from a platform action. Observability fields remain structured and non-secret; evidence bodies follow data-classification and redaction policies.

#### Failure behavior

Attempts to mutate immutable governance records fail closed. Missing required audit metadata or secret-like values in record payloads are rejected with stable validation codes.

#### Versioning

AuditEvent and related governance schemas version independently so evidence format evolution does not break customer APIs.

#### Replacement path

Introduce a new governance record version, backfill or dual-write as approved, and retain prior evidence formats for retention-policy duration before retirement.

#### Migration path

FEATURE-0012 defines the AuditEvent grammar contract; FEATURE-0013 and later own durable audit/decision services while preserving this boundary's immutability and redaction rules.

#### Reassessment trigger

Reassess for regulated/classified workloads, when retention or residency rules change, when a new privileged governance consumer appears, or before promoting governance evidence formats to a stable API.
