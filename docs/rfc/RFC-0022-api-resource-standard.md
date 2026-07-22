---
doc_type: rfc
id: RFC-0022
feature: FEATURE-0012
title: API Resource Standard
status: draft
phase: 2
architecture_status: approved_for_kiro
standard_maturity: draft
kiro_authorization: approved_for_requirements
controlling_handoff: ADH-2026-012
canonical_standard: docs/architecture/api-resource-standard.md
reuse_assessment: docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md
reuse_standard: docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md
ai_load_priority: high
ai_summary: Concise RFC index for the approved FEATURE-0012 provider-neutral API, resource, naming, status, reference, validation, and compatibility architecture.
---

# RFC-0022: API Resource Standard

## 1. Status

The FEATURE-0012 architecture baseline is approved for Kiro requirements generation under `ADH-2026-012`.

The canonical standard remains **Draft** until FEATURE-0012 implementation, conformance validation, final human review, and merge are complete. Approval of this RFC does not authorize design, tasks, or Cursor implementation; those retain separate gates.

## 2. Identifier mapping

`RFC-0022` and `FEATURE-0012` use separate repository numbering sequences:

- `RFC-0022` identifies this architecture decision record.
- `FEATURE-0012` identifies the implementation and governance work item.

## 3. Canonical source

The normative architecture is:

`docs/architecture/api-resource-standard.md`

This RFC is a concise decision index. It must not duplicate or silently redefine the canonical standard.

## 4. Decision

Sovrunn will extend mature HTTP, OpenAPI, JSON Schema, Problem Details, JSON Pointer, conditional-request, and declarative-resource conventions into a **Sovrunn-owned, provider-neutral platform grammar**.

The standard defines:

- eight object and resource profiles;
- naming, identity, metadata, immutable scope, and typed references;
- customer, operator, internal-engine, adapter, plugin, and governance boundaries;
- field ownership, mutability, status, phase, and condition semantics;
- strict decoding and ordered structural, semantic, reference, and policy validation;
- stable machine-readable error contracts;
- compatibility, boundedness, migration, conformance, fitness-function, and reassessment rules.

## 5. Architectural constraints

The core grammar must remain provider-neutral and cross-phase. Provider-, plugin-, adapter-, and vendor-native contracts stay behind explicit, owned, versioned boundaries.

FEATURE-0012 must not implement:

- provider integration or discovery;
- plugin or adapter execution;
- provisioning or reconciliation behavior;
- policy, entitlement, placement, or audit domain behavior owned by later features;
- vendor SDK types in common contracts;
- a wholesale Phase 1 rewrite;
- arbitrary unregistered extension maps.

## 6. Assurance position

The approved baseline has no identified unmitigated conflict with Sovrunn foundation principles or evaluated growth scenarios. Intentional boundaries are explicit, owned, versioned, observable, auditable, replaceable, and testable. Remaining uncertainties require documented migration paths and reassessment triggers.

Conformance evidence must include representative fixtures for managed resources, external observations, versioned definitions, immutable records, long-running operations, transient request/results, embedded values, and list envelopes.

## 7. Ownership boundary

Sovrunn owns:

- the resource-profile taxonomy and common platform grammar;
- scope, identity, reference, ownership, mutability, and boundary semantics;
- provider-neutral status, condition, validation, error, compatibility, and extension rules;
- conformance requirements and architecture fitness functions.

External standards and implementations remain externally owned. Sovrunn owns only their constrained use and composition within the platform.

## 8. Controlling references

- `ADH-2026-012`
- `FEATURE-0012`
- `docs/architecture/api-resource-standard.md`
- `docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md`
- `docs/reviews/reuse-assessments/FEATURE-0012-approval-evidence.md`
- `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md`
- `docs/phase2/PHASE2_EXECUTION_STRATEGY.md`
- `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

The reuse-assessment standard is cross-phase governance. Its current Phase 2 repository path is retained for compatibility until a separately approved migration establishes a new canonical path.
