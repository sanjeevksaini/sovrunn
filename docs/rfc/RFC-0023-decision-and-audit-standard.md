---
doc_type: rfc
title: RFC-0023 Decision and Audit Standard
status: draft
phase: 2
controlling_handoff: ADH-2026-014
ai_load_priority: high
ai_summary: RFC for common DecisionRecord and AuditEvent structures. DecisionRecord replaces the former term DecisionObject per ADH-2026-014. AuditEvent supports six governance scopes.
---

# RFC-0023: Decision and Audit Standard

See `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`.

## Controlling handoff

ADH-2026-014 (Approved, 2026-07-27) establishes the normative architecture for this RFC.

## Terminology

`DecisionRecord` is the canonical name for the common governed-conclusion envelope. The former term `DecisionObject` is superseded by this controlled correction.

## Decision

All major Sovrunn decisions must use structured `DecisionRecord` instances with authority, forms, typed results, reason codes, human-readable reasons, alternatives, policy references, suggested actions, obligations, provenance, and audit linkage.

`AuditEvent` supports six formal governance scopes: Platform, Organization, OrganizationUnit, Tenant, Project, and ServiceInstance. Scope does not grant authorization.

`DecisionRecord`, `EvaluationResult`, `AuditEvent`, `Operation`, and `DecisionContext` remain distinct concepts that reference each other but are never collapsed.

## Phase 2 constraint

FEATURE-0013 is contract-only within Phase 2. Schemas, vocabularies, validation, conformance fixtures, and documentation only. No production services, persistence, or runtime implementation.
