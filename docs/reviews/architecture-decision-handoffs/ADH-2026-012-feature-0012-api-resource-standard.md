---
doc_type: architecture_decision_handoff
handoff_id: ADH-2026-012
feature: FEATURE-0012
title: API, Resource Naming, Status, and Validation Standard
status: Approved
classification: Extension
phase: 2
approval_date: 2026-07-22
approving_role: Sovrunn Architecture Owner
---

# ADH-2026-012 — FEATURE-0012 API and Resource Standard

## Metadata

- Handoff ID: ADH-2026-012
- Date: 2026-07-22
- Related feature: FEATURE-0012
- Related phase: Phase 2, with cross-phase platform applicability
- Author: ChatGPT architecture synthesis with Sovrunn project owner
- Human approver: Sovrunn Architecture Owner
- Approval status: Approved

## Decision title

Approve the Sovrunn-owned, provider-neutral API and resource meta-model as the baseline for FEATURE-0012 Kiro requirements, design, and tasks.

## Summary

Sovrunn will extend mature API and schema standards into a coherent platform contract for resource profiles, naming, metadata, scope, references, boundaries, ownership, status, conditions, validation, errors, compatibility, and conformance. The architecture is broad enough to accommodate later resources, plugins, adapters, data sources, and consumers, while FEATURE-0012 implementation remains limited to shared primitives, schemas, validation, conformance, and Phase 1 compatibility analysis.

Canonical architecture:

`docs/architecture/api-resource-standard.md`

Controlling reuse standard:

`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

FEATURE-0011 is cross-phase governance. Its current canonical path is retained for compatibility; relocation requires a separately reviewed change because validators and historical evidence depend on it.

## Classification

- Extension

## Existing approved baseline

Relevant references:

- `docs/foundation/constitution.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md`
- `docs/phase2/PHASE2_EXECUTION_STRATEGY.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
- `docs/decisions/DEC-0026-reuse-before-build.md`
- `docs/decisions/DEC-0027-phase2-scope.md`
- `docs/decisions/DEC-0036-adapter-boundaries.md`
- `docs/rfc/RFC-0022-api-resource-standard.md`

## Approved decision

1. Use a Kubernetes-inspired declarative model without adopting Kubernetes namespace, CRD, controller, or storage assumptions as Sovrunn core contracts.
2. Use domain API groups and maturity versions; preserve Phase 1 compatibility until separately migrated.
3. Define eight object profiles: ManagedResource, ObservedExternalResource, VersionedDefinition, ImmutableRecord, LongRunningOperation, TransientRequestResult, EmbeddedValue, and ListEnvelope.
4. Give each persistent resource one immutable primary scope, except the platform root; model provider supply as related to but separate from the customer governance hierarchy.
5. Keep `scopeRef`, `ownerRef`, `location`, `sourceRef`, and `subjectRef` distinct.
6. Use typed constrained references with optional immutable UID; deny cross-tenant and cross-organization references by default without existence disclosure.
7. Classify contracts as customer-facing, operator-facing, internal-engine-facing, adapter-facing, plugin-facing, or governance-only.
8. Assign exactly one authoritative writer to every mutable field and condition type.
9. Use optional phase plus tri-state current-fact conditions; event history belongs to FEATURE-0013 records.
10. Reject duplicate and unknown fields by default; apply ordered structural, semantic, reference, authorization, and later-domain validation.
11. Use RFC 9457 Problem Details with stable Sovrunn codes and RFC 6901 JSON Pointer violation paths.
12. Use opaque resource versions with ETag/If-Match semantics for protected updates; defer PATCH semantics.
13. Permit only registered, namespaced, bounded extensions; extensions cannot silently become core decision inputs.
14. Require boundary ledgers, architecture fitness functions, compatibility checks, migration paths, and reassessment triggers.
15. Implement only shared architecture foundation and conformance work in FEATURE-0012; do not implement later-feature domain behavior.

## Reuse-before-build assessment

- Decision: Extend
- Mature foundations: HTTP semantics, OpenAPI 3.1, JSON Schema 2020-12, RFC 9457, RFC 6901, ETag/If-Match, and selected Kubernetes API conventions.
- Sovrunn-owned responsibility: Sovrunn owns the resource-profile taxonomy, API and naming conventions, common metadata, identity, scope and reference semantics, API-boundary classification, field ownership and mutability rules, status and condition grammar, validation and error contracts, provider-neutrality constraints, compatibility policy, conformance rules, and reassessment triggers.
- Reused or extended responsibility: HTTP semantics, OpenAPI 3.1, JSON Schema 2020-12, RFC 9457 Problem Details, RFC 6901 JSON Pointer, ETag/If-Match concurrency semantics, and selected Kubernetes API conventions are reused or extended.
- Responsibility/control boundary: External standards own their generic syntax and semantics. Sovrunn owns the constrained sovereign PaaS contract and conformance policy. Provider-, plugin-, adapter-, and vendor-native types remain behind classified boundaries and do not become customer-facing or core resource contracts.
- Adapter required: No. FEATURE-0012 defines adapter-facing grammar but introduces no external runtime integration.

## Foundation assurance

The approved architecture was checked against accuracy, transparency, security, sovereignty, auditability, flexibility, provider neutrality, organization-first governance, tenant isolation, reuse-before-build, plugin/adapter separation, and governed AI consumption.

Approval statement:

> The FEATURE-0012 architecture has no identified unmitigated conflict with Sovrunn foundation principles or evaluated growth scenarios. Its intentional boundaries are explicit, owned, versioned, observable, auditable, replaceable, and testable. Remaining uncertainties have documented migration paths and reassessment triggers. Approved as the baseline for Kiro requirements, design, and tasks.

## Phase impact

- Current phase allowed: Yes
- Phase 2 output: normative standard, shared primitives, schema conventions, validators, conformance fixtures, compatibility report, and gate support.
- Cross-phase effect: later APIs and resources consume the standard unless an approved architecture change supersedes it.
- No Phase 2 sequence change.

## Non-goals

- No functional provider or substrate model.
- No provider integration or infrastructure provisioning.
- No plugin execution or adapter runtime protocol.
- No policy or placement engine.
- No DecisionObject or AuditEvent domain payload.
- No persistence selection, billing, failover, or autonomous AI behavior.
- No wholesale Phase 1 API rewrite.
- No unrestricted arbitrary extension map.

## Required action

- Replace the placeholder `docs/architecture/api-resource-standard.md` with the approved baseline.
- Create the FEATURE-0012 reuse assessment and structured approval evidence.
- Align RFC, current context, traceability, and Phase 2 acceptance-gate references.
- Initialize the Kiro spec directory resolved by `docs/features/FEATURE_INDEX.md`.
- Generate `requirements.md` only, then wait for `APPROVED_FOR_DESIGN`.

## Impacted files

- `docs/architecture/api-resource-standard.md`
- `docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md`
- `docs/reviews/reuse-assessments/FEATURE-0012-approval-evidence.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-012-feature-0012-api-resource-standard.md`
- `docs/rfc/RFC-0022-api-resource-standard.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_DECISION_SUMMARY.md`
- `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md`
- `docs/phase2/PHASE2_ACCEPTANCE_GATES.md`
- `.kiro/specs/api-resource-naming-status-and-validation-standard/.config.kiro`

## Acceptance criteria for Kiro requirements

- [ ] File contains `FEATURE-0012` and `Stage: Requirements`.
- [ ] Every normative architecture decision is represented by testable requirements and acceptance criteria.
- [ ] Resource, scope, boundary/ownership, compatibility-scenario, and risk matrices are preserved by reference or equivalent requirement coverage.
- [ ] Shared grammar is separated from FEATURE-0013+ domain semantics.
- [ ] Phase 1 compatibility is addressed without requiring a rewrite.
- [ ] Provider-neutrality, isolation, ownership, redaction, boundedness, audit linkage, and compatibility are testable.
- [ ] Non-goals and reassessment triggers are explicit.
- [ ] Requirements do not choose implementation libraries or create source code.

## Explicit instructions to Kiro

- Generate `requirements.md` only.
- Use the Kiro slug from `docs/features/FEATURE_INDEX.md`: `api-resource-naming-status-and-validation-standard`.
- Treat `docs/architecture/api-resource-standard.md` and this ADH as controlling architecture.
- Do not introduce new resource profiles, scope kinds, boundaries, or runtime behavior.
- Record unresolved semantic gaps as `ARCHITECTURE_DECISION_REQUIRED`.
- Do not generate design, tasks, or implementation until separately authorized.

## Human approval

- Approval status: Approved
- Approved by: Sovrunn Architecture Owner
- Date: 2026-07-22
- Notes: Architecture baseline and Kiro requirements generation approved. Later stage gates remain mandatory.
