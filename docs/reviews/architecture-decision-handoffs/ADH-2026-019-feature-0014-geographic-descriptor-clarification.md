# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-019
- Date: 2026-07-30
- Source discussion: Sovrunn Architecture Governor ChatGPT Project — Feature 0014 Architecture
- Related feature: FEATURE-0014
- Related phase: Phase 2
- Author: Codex draft for architecture-owner review
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

FEATURE-0014 geographic descriptors are syntax-normalized declarations, not authoritative reference data

## Summary

Remove authoritative geographic membership from FEATURE-0014. Preserve the optional descriptor shape, validate only syntax, finite bounds, and country-prefix agreement, and defer any authoritative geographic reference capability to separately approved reusable platform work.

## Classification

Replacement clarification for the geographic portion of `F14-AD-011`, `F14-R16`, and downstream `F14-REQ-27`. All other ADH-2026-018 decisions remain unchanged.

## Existing approved baseline

ADH-2026-018 and the approved FEATURE-0014 architecture define five provider-neutral topology resources and delegate normalized optional geographic representation to design. The previously approved requirements and design interpreted “normalized” as authoritative ISO dataset membership. This handoff replaces only that interpretation.

## Decision or proposed decision

1. `ProviderLocation.spec.geo` remains optional descriptive declared topology.
2. `countryCode`, when present, is syntax-normalized as exactly two uppercase ASCII letters.
3. `subdivisionCode`, when present, is syntax-normalized as the same two-letter country prefix, a hyphen, and one to three uppercase ASCII letters or digits.
4. FEATURE-0014 validates syntax, bounds, and country-prefix agreement only.
5. FEATURE-0014 does not assert that either value is assigned by ISO, currently active, sovereign, resident, compliant, or suitable for placement.
6. FEATURE-0014 owns no ISO dataset, embedded snapshot, generator, manifest, refresh process, runtime lookup, or third-party geographic-code dependency.
7. A syntactically valid but unassigned code remains valid declared metadata; it is not described as an authoritative ISO code.
8. If authoritative geographic membership or classification becomes necessary, it must be introduced as a separately approved reusable platform capability and consumed without changing FEATURE-0014 resource identity or hierarchy.

## Rationale

Authoritative geographic reference data changes independently of provider topology and is reusable across features. Owning or selecting that dataset in FEATURE-0014 would add lifecycle, licensing, dependency, and update responsibilities that do not contribute to its five-resource hierarchy. Syntax-only normalization preserves interoperable shape while preventing residency, compliance, placement, or sovereignty inference.

## Reuse-before-build assessment

- Disposition: Reuse/Defer.
- FEATURE-0014 reuses regular-expression, finite-bound, and deterministic validation primitives already supplied by FEATURE-0012.
- FEATURE-0014 neither builds nor selects a geographic reference dataset or library.
- A future authoritative reference capability must complete its own reuse assessment and serve multiple consuming features.

## Phase impact

Allowed in Phase 2. The clarification removes work and dependencies; it does not change feature sequence, resource hierarchy, public identity, or adjacent-feature ownership.

## Conflict check

No DEC or RFC conflict is introduced. This handoff intentionally replaces the downstream membership interpretation of F14-REQ-27 and design DD-04 while retaining ADH-2026-018 otherwise unchanged.

## Required downstream changes

- Clarify `F14-AD-011` and `F14-R16` in the canonical architecture.
- Revise `F14-REQ-27` so “invalid” means syntactically invalid or prefix-inconsistent; remove unknown-membership rejection.
- Remove the embedded-dataset decision from design DD-04 and section 6.2.
- Remove Task 11 and all dataset, manifest, membership, library, and refresh work.
- Keep format and prefix validation in the ProviderLocation validator task.
- Revalidate and reapprove requirements, design, and tasks before Cursor execution.

## Required action

Update architecture doc, requirements, design, tasks, reviewers, validators, approval digests, and Cursor task context as one atomic repository change.

## Impacted files

- `docs/architecture/provider-neutral-resource-model.md`
- `.kiro/specs/provider-neutral-resource-model/requirements.md`
- `.kiro/specs/provider-neutral-resource-model/design.md`
- `.kiro/specs/provider-neutral-resource-model/tasks.md`
- FEATURE-0014 prompts, reviewers, validators, state, and Cursor automation

## Impacted features

- FEATURE-0014 only.
- A future reusable geographic reference feature may consume authoritative datasets; no feature ID is assigned here.

## Acceptance criteria for Kiro update

- F14-REQ-27 rejects malformed and prefix-inconsistent descriptors only.
- Design contains no dataset, library, membership, generator, manifest, or refresh implementation.
- Task 11 is absent and Task 12 owns syntax/prefix validation.
- All existing decision/risk IDs and non-geographic boundaries remain traceable.

## Explicit instructions to Kiro

Do not generate a dataset, dependency, lookup, generator, manifest, update job, or runtime fetch. Do not reinterpret a syntactically valid value as an authoritative ISO assignment. Make no unrelated requirement, design, or task change.

## Human approval

Approved by Sanjeev Kumar in the FEATURE-0014 architecture discussion on 2026-07-30.

## Unchanged boundaries

Exactly five owned kinds, FEATURE-0012 grammar reuse, FEATURE-0013 `NOT_APPLICABLE`, connectivity non-inference, and FEATURE-0015/0016 ownership remain unchanged.
