# ADH-2026-014 Terminology Reconciliation Record

> **⚠️ SUPERSESSION NOTICE — six-scope AuditEvent material superseded by ADH-2026-015.**
> This record is **preserved as historical evidence** of the ADH-2026-014
> terminology reconciliation and remains valid for the `DecisionObject` →
> `DecisionRecord` correction. However, **all six-scope `AuditEvent` material in
> this document is superseded by ADH-2026-015** (Approved, 2026-07-27), which is
> the authoritative source for `AuditEvent` scope. Under ADH-2026-015, the single
> canonical `ScopeKind` vocabulary has exactly seven values — `Platform`,
> `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, and
> `ServiceInstance` — where `Provider` is retained (the earlier six-scope list
> omitted it) and `ServiceInstance` is additive; `AuditEvent` permits all seven
> through `metadata.scopeRef` as its sole scope authority, and the canonical
> `Platform` form is an absent/nil `metadata.scopeRef`.
>
> This notice does **not** rewrite the historical evidence below: the
> "Six-scope AuditEvent correction" section and any six-scope wording are
> retained verbatim as the record of what ADH-2026-014 originally stated. They
> are not restated as though they described seven scopes.
>
> Authoritative successor material:
> - `docs/reviews/architecture-decision-handoffs/ADH-2026-015-feature-0013-scope-vocabulary-clarification.md`
> - `docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md`
> - `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md` (section 29)

## Purpose

Documents the controlled replacement of `DecisionObject` with `DecisionRecord` across the Sovrunn repository baseline, as mandated by ADH-2026-014 (Approved, 2026-07-27).

## Controlling reference

- ADH-2026-014 — FEATURE-0013 Decision Record and AuditEvent Standard
- Approval: Sanjeev Kumar, 2026-07-27
- Classification: Extension with controlled correction

## Terminology correction

| Prior term | Corrected term | Rationale |
|---|---|---|
| `DecisionObject` | `DecisionRecord` | "Object" is ambiguous; "Record" communicates immutability, governed conclusion, and audit commitment. ADH-2026-014 section "Classification" point 1. |

## Six-scope AuditEvent correction

| Prior scope | Corrected scope | Rationale |
|---|---|---|
| Organization-only (FEATURE-0012 alpha) | Platform, Organization, OrganizationUnit, Tenant, Project, ServiceInstance | Required for governed decisions at all governance boundaries. ADH-2026-014 section "Classification" point 2. |

## Corrected normative files

| File | Prior content | Correction applied | Rationale |
|---|---|---|---|
| `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md` | `DecisionObject` in Invariant H, dependency graph, shared boundary map, policy boundary | Replaced with `DecisionRecord`; added six-scope AuditEvent note | Normative Phase 2 spine uses canonical terminology |
| `docs/context/CURRENT_ARCHITECTURE_BASELINE.md` | FEATURE-0012 described as "Kiro requirements generation"; no FEATURE-0013 baseline | Added FEATURE-0013 baseline section; corrected FEATURE-0012 status; added DecisionRecord and six-scope decisions | Baseline must reflect current approved state |
| `docs/context/CURRENT_DECISION_SUMMARY.md` | No ADH-2026-014 entries; FEATURE-0012 described without ADH-2026-013 | Added ADH-2026-014 decisions; added FEATURE-0012 approval status | Decision summary must be current |
| `docs/context/CURRENT_PHASE_CONTEXT.md` | No FEATURE-0012 completion or FEATURE-0013 active status | Added completed/active feature tables; updated build scope wording | Phase context must reflect execution state |
| `docs/rfc/RFC-0023-decision-and-audit-standard.md` | Generic "decision objects" language; no controlling handoff | Updated with DecisionRecord terminology, six-scope AuditEvent, ADH-2026-014 reference | RFC must align with approved architecture |
| `docs/features/FEATURE_INDEX.md` | FEATURE-0013 named "Decision Object and AuditEvent Standard" | Updated name to "Decision Record and AuditEvent Standard"; added ADH-2026-014 reference | Feature index uses canonical names |
| `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md` | FEATURE-0012 status "Architecture Approved / Kiro requirements pending"; FEATURE-0013 status "Planned" | Updated FEATURE-0012 to Implemented and Merged through PR #14 as commit `a1b74fb` into `phase2-reuse-first-paas-fabric-foundation`; updated FEATURE-0013 with ADH-2026-014, DecisionRecord terminology correction, six-scope AuditEvent, controlling architecture approved 2026-07-27 | Traceability matrix must reflect current feature status and approved corrections |
| `docs/phase2/PHASE2_ACCEPTANCE_GATES.md` | No FEATURE-0013 specific gates | Added section 6 with FEATURE-0013 gate criteria | Gates must cover active features |
| `docs/glossary.md` | `DecisionObject` entry | Replaced with `DecisionRecord` entry including ADH-2026-014 reference | Glossary is normative terminology source |
| `docs/architecture/observability-and-audit-baseline.md` | `DecisionObject` in Phase 2 audit section | Replaced with `DecisionRecord`; added six-scope AuditEvent note | Architecture docs use canonical terminology |
| `.kiro/agents/sovrunn-spec.md` | No FEATURE-0013 or ADH-2026-014 resources | Added FEATURE-0013 and ADH-2026-014 resource references | Agent must load controlling architecture |

## Files intentionally NOT corrected

| File | Reason |
|---|---|
| `docs/architecture/api-resource-standard.md` | Uses `DecisionObject` as a deferred-domain reference in non-goal context ("DecisionObject or AuditEvent domain payloads owned by FEATURE-0013"). These references describe what FEATURE-0012 does NOT own and will be updated when FEATURE-0013 implementation introduces the actual DecisionRecord schemas. Correcting them now would change the semantic meaning of a completed feature's non-goal boundary. |
| `.kiro/specs/api-resource-naming-status-and-validation-standard/requirements.md` | Same as above — FEATURE-0012 non-goal and deferred references. Owned by completed FEATURE-0012 scope. |
| `.kiro/specs/api-resource-naming-status-and-validation-standard/design.md` | Same as above. |
| `docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md` | Non-goal reference in completed feature spec. |
| `docs/reviews/architecture-decision-handoffs/ADH-2026-012-*.md` | Historical handoff document records the state at time of approval. |
| `docs/reviews/architecture-decision-handoffs/ADH-2026-014-*.md` | Self-referencing (describes the correction being made). |
| `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md` | Contains its own ADH-2026-014 terminology history section. No correction needed. |

## Backward compatibility treatment

The Kiro slug `decision-object-and-auditevent-standard` is preserved for directory compatibility. The slug is a filesystem artifact, not a normative terminology assertion.

The FEATURE-0013 canonical feature title is "Decision Record and AuditEvent Standard". The existing Kiro directory slug `decision-object-and-auditevent-standard` is retained as a stable legacy identifier to avoid unnecessary path churn.

## Verification

- PRE-01: Complete — all normative `DecisionObject` occurrences in Phase 2 spine, current baseline, decision summary, RFC, feature index, traceability matrix, and acceptance gates have been audited.
- PRE-02: Complete — each occurrence corrected to `DecisionRecord`.
- PRE-03: Complete — this document constitutes the required traceability record.

## Human approval reference

ADH-2026-014 approval by Sanjeev Kumar on 2026-07-27 authorizes this terminology reconciliation.
