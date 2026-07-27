---
doc_type: rfc
title: RFC-0023 Decision and Audit Standard
status: draft
phase: 2
controlling_handoff: ADH-2026-014
controlling_handoffs:
  - ADH-2026-014
  - ADH-2026-015
  - ADH-2026-016
ai_load_priority: high
ai_summary: RFC for common DecisionRecord and AuditEvent structures. DecisionRecord replaces the former term DecisionObject per ADH-2026-014. Per ADH-2026-015, ScopeKind is the single canonical seven-value scope vocabulary (adds ServiceInstance, retains Provider) and AuditEvent permits all seven values via metadata.scopeRef as sole scope identity. Per ADH-2026-016, DecisionRecord uses metadata.scopeRef as its sole scope authority, sensitivity uses the ordered vocabulary PUBLIC<INTERNAL<CONFIDENTIAL<RESTRICTED, security validation is structural-only, cryptographic selection is deferred (ADR-F13-002), and conformance scenarios are F13-CF-01..F13-CF-28.
---

# RFC-0023: Decision and Audit Standard

See `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`.

## Controlling handoff

ADH-2026-014, ADH-2026-015, and ADH-2026-016 (all Approved, 2026-07-27) are the joint controlling handoffs for this RFC. ADH-2026-015 supersedes only the earlier six-scope `AuditEvent` statement. ADH-2026-016 is a clarification that clarifies previously unresolved contract boundaries without introducing new architecture or authorizing any runtime capability or product/provider/algorithm selection.

## Terminology

`DecisionRecord` is the canonical name for the common governed-conclusion envelope. The former term `DecisionObject` is superseded by this controlled correction.

## Decision

All major Sovrunn decisions must use structured `DecisionRecord` instances with authority, forms, typed results, reason codes, human-readable reasons, alternatives, policy references, suggested actions, obligations, provenance, and audit linkage.

`DecisionRecord` uses `metadata.scopeRef` as its sole logical and serialized scope authority (ADH-2026-016). No top-level `DecisionRecord` `scopeRef` or parallel scope source exists; `Platform` is the absent/nil `metadata.scopeRef` form and every non-Platform scope carries a non-nil `metadata.scopeRef` with a canonical `ScopeKind` and UID.

`AuditEvent` supports the seven canonical `ScopeKind` values: Platform, Organization, OrganizationUnit, Tenant, Project, Provider, and ServiceInstance (per ADH-2026-015; `ServiceInstance` is additive and `Provider` is retained). `AuditEvent` uses `metadata.scopeRef` as its sole canonical scope identity, with no separate scope enum. Scope does not grant authorization.

`DecisionRecord`, `EvaluationResult`, `AuditEvent`, `Operation`, and `DecisionContext` remain distinct concepts that reference each other but are never collapsed.

## ADH-2026-016 contract boundaries

- Sensitivity classification uses the closed ordered provider-neutral vocabulary `PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED`; jurisdiction-specific labels are versioned profile mappings; profiles permitting captured output declare a sensitivity ceiling, content-category rules, and maximum fields/bytes/depth/value-types.
- FEATURE-0013 security validation is bounded to structural conformance; comprehensive semantic content scanning belongs to later approved security features.
- Algorithm-agile carrier fields and structural trust metadata are contract-now; the canonicalization algorithm/profile, digest-covered fields, signature algorithm, and cryptographic services are DEFERRED under ADR-F13-002; RFC 8785 is illustrative only.
- The conformance scenarios map one-to-one to stable IDs `F13-CF-01` through `F13-CF-28`.
- `docs/traceability/ADH-2026-014-six-scope-auditevent-compatibility.md` is SUPERSEDED historical evidence; `docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md` is authoritative for the seven-value contract.

## Phase 2 constraint

FEATURE-0013 is contract-only within Phase 2. Schemas, vocabularies, validation, conformance fixtures, and documentation only. No production services, persistence, or runtime implementation.
