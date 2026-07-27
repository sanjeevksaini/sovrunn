---
doc_type: rfc
title: RFC-0023 Decision and Audit Standard
status: draft
phase: 2
controlling_handoff: ADH-2026-014
controlling_handoffs:
  - ADH-2026-014
  - ADH-2026-015
ai_load_priority: high
ai_summary: RFC for common DecisionRecord and AuditEvent structures. DecisionRecord replaces the former term DecisionObject per ADH-2026-014. Per ADH-2026-015, ScopeKind is the single canonical seven-value scope vocabulary (adds ServiceInstance, retains Provider) and AuditEvent permits all seven values via metadata.scopeRef as sole scope identity.
---

# RFC-0023: Decision and Audit Standard

See `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`.

## Controlling handoff

ADH-2026-014 and ADH-2026-015 (both Approved, 2026-07-27) are the joint controlling handoffs for this RFC. ADH-2026-015 supersedes only the earlier six-scope `AuditEvent` statement.

## Terminology

`DecisionRecord` is the canonical name for the common governed-conclusion envelope. The former term `DecisionObject` is superseded by this controlled correction.

## Decision

All major Sovrunn decisions must use structured `DecisionRecord` instances with authority, forms, typed results, reason codes, human-readable reasons, alternatives, policy references, suggested actions, obligations, provenance, and audit linkage.

`AuditEvent` supports the seven canonical `ScopeKind` values: Platform, Organization, OrganizationUnit, Tenant, Project, Provider, and ServiceInstance (per ADH-2026-015; `ServiceInstance` is additive and `Provider` is retained). `AuditEvent` uses `metadata.scopeRef` as its sole canonical scope identity, with no separate scope enum. Scope does not grant authorization.

`DecisionRecord`, `EvaluationResult`, `AuditEvent`, `Operation`, and `DecisionContext` remain distinct concepts that reference each other but are never collapsed.

## Phase 2 constraint

FEATURE-0013 is contract-only within Phase 2. Schemas, vocabularies, validation, conformance fixtures, and documentation only. No production services, persistence, or runtime implementation.
