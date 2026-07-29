---
doc_type: rfc
title: RFC-0023 Decision and Audit Standard
status: draft
phase: 2
controlling_handoff: ADH-2026-017
historical_handoffs:
  - ADH-2026-014
  - ADH-2026-015
  - ADH-2026-016
ai_load_priority: high
ai_summary: RFC for the approved consolidated FEATURE-0013 DecisionRecord and AuditEvent architecture under ADH-2026-017. It preserves FEATURE-0012's six governance scopes, uses metadata.scopeRef as sole scope authority, models ServiceInstance as a typed subject, bounds validation to structural checks, and defers cryptographic selections.
---

# RFC-0023: Decision and Audit Standard

See `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`.

## Controlling handoff

`ADH-2026-017` is the approved single replacement handoff for this RFC and the complete FEATURE-0013 architecture, bound to the exact canonical architecture digest. ADH-2026-014, ADH-2026-015, and ADH-2026-016 remain historical provenance and are not separate downstream instructions.

## Terminology

`DecisionRecord` is the canonical name for the common governed-conclusion envelope. The former term `DecisionObject` is superseded by this controlled correction.

## Decision

All major Sovrunn decisions must use structured `DecisionRecord` instances with authority, forms, typed results, reason codes, human-readable reasons, alternatives, policy references, suggested actions, obligations, provenance, and audit linkage.

`DecisionRecord` uses `metadata.scopeRef` as its sole logical and serialized scope authority. No top-level `DecisionRecord` `scopeRef` or parallel scope source exists; `Platform` is the absent/nil `metadata.scopeRef` form and every non-Platform scope carries a non-nil `metadata.scopeRef` with a canonical `ScopeKind` and UID.

`AuditEvent` supports FEATURE-0012's six canonical governance scopes: Platform, Organization, OrganizationUnit, Tenant, Project, and Provider. It uses `metadata.scopeRef` as its sole scope identity, with no separate scope enum. A `ServiceInstance` is carried as a typed subject reference, normally within Project scope; it is not a scope value. Scope does not grant authorization.

`DecisionRecord`, `EvaluationResult`, `AuditEvent`, `Operation`, and `DecisionContext` remain distinct concepts that reference each other but are never collapsed.

## Consolidated contract boundaries

- Sensitivity classification uses the closed ordered provider-neutral vocabulary `PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED`; jurisdiction-specific labels are versioned profile mappings; profiles permitting captured output declare a sensitivity ceiling, content-category rules, and maximum fields/bytes/depth/value-types.
- FEATURE-0013 security validation is bounded to structural conformance; comprehensive semantic content scanning belongs to later approved security features.
- Algorithm-agile carrier fields and structural trust metadata are contract-now; the canonicalization algorithm/profile, digest-covered fields, signature algorithm, and cryptographic services are DEFERRED under ADR-F13-002; RFC 8785 is illustrative only.
- The conformance scenarios map one-to-one to stable IDs `F13-CF-01` through `F13-CF-28`.
- Historical ADH-specific compatibility files are provenance only. The consolidated architecture and its ADH-2026-017 digest binding are the sole downstream architecture input after approval.

## Phase 2 constraint

FEATURE-0013 is contract-only within Phase 2. Schemas, vocabularies, validation, conformance fixtures, and documentation only. No production services, persistence, or runtime implementation.
