# Implementation Plan

Feature: FEATURE-0011 — Reuse Assessment Standard
Phase: Phase 2 — Reuse-First PaaS Fabric Foundation
Stage: Tasks

## Overview

This plan decomposes the approved design without making new design or
architecture decisions. Tasks are ordered, dependency-aware, and sized for
sequential implementation and review. All Phase 2 non-goals are preserved:
no runtime resource, no vendor selection, no provider integration, no
plugin execution, no persistence, no billing, no failover, and no
autonomous AI operations.

Legend: each task lists its target paths, the Requirements it satisfies
(Req N), the design sections or rule identifiers it implements, its
dependencies, objective completion criteria, and verification.

## FEATURE-0011 reuse summary

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| Reuse Assessment Standard governance contract | Extend | Extends the existing reuse-before-build baseline, architecture-decision and RFC review practices, risk-control registers, and the existing draft Phase 2 reuse assessment format, rather than building a new governance mechanism. | Approved | ADH-2026-011 (DEC-0026, RFC-0021) |

Sovrunn owns the four-disposition vocabulary, capability-level assessment
rules, sovereign-deployment criteria, provider-neutrality checks,
adapter-boundary requirements, Phase 2 scope controls, future-feature
mitigation requirements, architecture traceability, and feature-gate
structure. General architecture-decision, software-selection, and
risk-management practices are reused or extended.

This task-stage summary does not replace the complete FEATURE-0011
reuse-assessment instance created by Task 2.

## Tasks

- [ ] 1. Consolidate and version the complete canonical standard
  - Target: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
  - Implement the complete approved canonical contract from Requirements
    1–18: feature-level and capability-level assessment formats; controlled
    vocabularies; responsibility and boundary fields; sovereignty, security,
    operations, licensing, and portability considerations; phase impact and
    non-goals; Build rejection rationale; mitigation fields; traceability;
    decision status and change control; conceptual examples; structural
    validation; consistency validation; human semantic review; repository
    consistency; reassessment lifecycle; and strict feature-gate
    enforcement expectations.
  - Add front-matter field `reuse_assessment_format_version: 1.0.0`.
  - Requirements: 1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18
  - Design: Canonical artifact architecture; Data Models; Controlled
    vocabularies; Validation architecture; Versioning and alignment;
    Reassessment lifecycle; Repository-level risk-control matrix; Feature
    gate integration
  - Dependencies: none
  - Completion criteria: canonical file declares version 1.0.0; every
    Requirement 1–18 element is represented without omission or conflicting
    duplication; normative field definitions and controlled vocabularies
    occur in one canonical section; explanatory references and traceability
    mappings may reference those definitions elsewhere; no supporting
    section redefines the canonical schema; the document remains a draft
    (status not set to Approved) until recorded human approval.
  - Verification: `grep -n "reuse_assessment_format_version: 1.0.0"
    docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`; manual field,
    vocabulary, and section review against requirements and design.

- [ ] 2. Author the actual FEATURE-0011 reuse-assessment instance
  - Target: `docs/features/FEATURE-0011-reuse-assessment-standard.md` [NEW]
  - This file is the actual FEATURE-0011 reuse-assessment instance and the
    assessment artifact validated by the gate. It must contain: Markdown
    front matter `reuse_assessment_format_version: 1.0.0`; purpose and
    acceptance criteria; the exact approved feature-level reuse summary; and
    one complete capability-level assessment for the "Reuse Assessment
    Standard governance contract" with disposition Extend; decision status
    Approved; assessment owner; ADH/DEC/RFC traceability; candidate
    category, strengths, constraints, and selected approach; Sovrunn-owned
    and reused/extended responsibilities; boundary data and control; adapter
    fields and rationale; suitability considerations; allowed-in-current-
    phase value and phase work; deferred work and non-goals; exit or
    migration boundary; architecture risks; preventive, detection, and
    corrective controls; residual risk; replacement risk; reassessment
    triggers; linked acceptance criteria; and validation and review
    evidence.
  - The assessment must include recorded human-approval evidence that
    identifies the approving person or role; the approval date; the approved
    ADH or assessment-review reference; and evidence that the approval
    applies to the recorded Extend disposition and responsibility boundary.
    The evidence may resolve through ADH-2026-011 when that approved
    artifact already contains the approver and approval date.
  - All values must be sourced only from the approved requirements, design,
    ADH-2026-011, and Architecture Spine. Do not invent an approver or date
    or any missing semantic decision; if the approved source does not
    provide them, stop and report `ARCHITECTURE_DECISION_REQUIRED` (or the
    repository's applicable missing-approval-evidence condition).
  - Requirements: 1,2,3,4,5,6,7,8,9,10,11,12,16,17
  - Design: FEATURE-0011 reuse summary; Data Models; Repository artifact map
  - Dependencies: Task 1
  - Completion criteria: file exists, declares the format version, and holds
    a complete capability-level assessment whose values trace to approved
    sources; a resolvable approver/role and approval date (directly or via
    ADH-2026-011) are recorded as approval evidence; decision status
    Approved is recorded here as the assessment's own field, not a change to
    the canonical standard's status.
  - Verification: `grep -n "reuse_assessment_format_version: 1.0.0"
    docs/features/FEATURE-0011-reuse-assessment-standard.md`; confirm every
    mandatory field is present.

- [ ] 3. Align Kiro, Cursor, reviewer prompts and Kiro policy
  - Targets: `docs/prompts/kiro/requirements.prompt.md`,
    `docs/prompts/kiro/design.prompt.md`,
    `docs/prompts/kiro/tasks.prompt.md`,
    `docs/prompts/cursor/task.prompt.md`,
    `docs/prompts/reviewer/spec-review.prompt.md`,
    `docs/prompts/reviewer/approval-review.prompt.md`,
    `docs/automation/KIRO_DECISION_POLICY.md`
  - Each artifact links to `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
    and must not duplicate or redefine the field schema. Inspect and update
    only if inconsistent. Link-only artifacts declare no numeric version.
  - Requirements: 11,15,16
  - Design: Repository consistency and alignment; rules RA-C12, RA-C14
  - Dependencies: Task 1
  - Completion criteria: each artifact references the canonical path and
    contains no duplicated field-definition block.
  - Verification: `grep -n "PHASE2_REUSE_ASSESSMENT_STANDARD.md"` on each
    listed file; manual duplicated-schema review (RA-C12/RA-C14 intent).

- [ ] 4. Align Feature Factory documents and templates
  - Targets: `docs/automation/FEATURE_FACTORY.md`,
    `docs/ai/AI_FEATURE_FACTORY.md`,
    `docs/templates/ARCHITECTURE_DECISION_HANDOFF.md`,
    `docs/templates/ARCHITECTURE_CHANGE_REQUEST.md`,
    `docs/templates/RFC_TEMPLATE.md`,
    `docs/templates/FEATURE_REVIEW_TEMPLATE.md`
  - Link to the canonical standard by reference; do not duplicate the field
    schema. Inspect and update only if inconsistent.
  - Requirements: 15,16
  - Design: Repository consistency and alignment; rules RA-C12, RA-C14
  - Dependencies: Task 1
  - Completion criteria: each artifact references the canonical path with no
    duplicated schema block.
  - Verification: `grep -n "PHASE2_REUSE_ASSESSMENT_STANDARD.md"` on each
    listed file; manual duplicated-schema review.

- [ ] 5. Inspect and align the RFC and Phase 2 acceptance gates
  - Targets: `docs/rfc/RFC-0021-reuse-first-architecture.md`,
    `docs/phase2/PHASE2_ACCEPTANCE_GATES.md`
  - Inspect and update only if inconsistent with the canonical standard; add
    a canonical-source reference where required; do not duplicate the schema.
  - Requirements: 16,18
  - Design: Repository artifact map; Feature gate integration
  - Dependencies: Task 1
  - Completion criteria: both files are consistent with the canonical
    standard and reference it where required.
  - Verification: `grep -n "PHASE2_REUSE_ASSESSMENT_STANDARD.md"` on both
    files; manual consistency review.

- [ ] 6. Update traceability matrices
  - Targets: `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md`,
    `docs/traceability/DECISION_TRACEABILITY_MATRIX.md`
  - Add the FEATURE-0011 entry and its mapping to DEC-0026, DEC-0036, and
    RFC-0021.
  - Requirements: 10,16
  - Design: Repository artifact map; Repository consistency and alignment
  - Dependencies: Task 2
  - Completion criteria: both matrices contain a FEATURE-0011 row with the
    correct decision references.
  - Verification: `grep -n "FEATURE-0011"
    docs/traceability/FEATURE_TRACEABILITY_MATRIX.md
    docs/traceability/DECISION_TRACEABILITY_MATRIX.md`.

- [ ] 7. Create the validator entry point and diagnostic framework
  - Target: `scripts/reuse-assessment-check.sh` [NEW] (Bash entry point that
    may invoke Python 3 for deterministic Markdown parsing; no third-party
    framework or runtime service)
  - Add the shell version marker `# reuse_assessment_format_version=1.0.0`.
    Implement the deterministic, non-mutating interface: inputs (repository
    root, feature identifier, assessment artifact paths, validation mode);
    diagnostic fields (stable rule id, layer, feature id, file path,
    section/field, message, severity, corrective guidance); severity
    vocabulary (error, warning); exit codes (0 pass, 1 failures, 2 config or
    internal error); ordering by file path, then section/field, then rule
    id; strict mode for FEATURE-0011+; legacy exemption for FEATURE-0001–
    0010; fail-safe on unknown/malformed feature identifiers.
  - Requirements: 13,14,15,18
  - Design: Validation contract; Validation architecture; Structural rule
    identifiers and diagnostic behavior
  - Dependencies: Task 1
  - Completion criteria: script runs, never mutates inputs, emits ordered
    diagnostics, returns 0/1/2 as specified, and declares the version
    marker; human approval is never inferred by automation.
  - Verification: `grep -n "reuse_assessment_format_version=1.0.0"
    scripts/reuse-assessment-check.sh`; usage without arguments returns
    exit 2.

- [ ] 8. Implement structural rules RA-S01 through RA-S10
  - Target: `scripts/reuse-assessment-check.sh`
  - RA-S01 required section missing; RA-S02 required field missing; RA-S03
    invalid disposition; RA-S04 invalid decision status; RA-S05 invalid
    adapter value; RA-S06 invalid phase value; RA-S07 invalid
    vendor-native-types value; RA-S08 invalid replacement-risk value; RA-S09
    missing risk-control field; RA-S10 missing traceability field.
  - Requirements: 13 (sources 3,4,5,6,7,9,10)
  - Design: Validation architecture Layer 1; Structural rule identifiers
  - Dependencies: Task 7
  - Completion criteria: each RA-S rule emits its stable id at severity error
    and fails the run with exit 1 on violation.
  - Verification: run against invalid fixtures (Task 14); assert each RA-S id
    appears and exit code is 1.

- [ ] 9. Implement consistency rules (RA-C01–C02, RA-C04–C08, RA-C10–C12,
  RA-C14)
  - Target: `scripts/reuse-assessment-check.sh`
  - Own exactly: RA-C01 adapter rationale mandatory for Yes and No; RA-C02
    adapter-related DEC-0036 reference; RA-C04 exact conceptual-example
    label; RA-C05 DEC/RFC/ADH reference existence; RA-C06 adapter/contract
    identifier mandatory (reserved literal `none` when Adapter required is
    No); RA-C07 Build triple rejection rationale and ownership; RA-C08 risk
    triple controls; RA-C10 feature identifier matches `^FEATURE-[0-9]{4}$`;
    RA-C11 Phase 2 scope acknowledgement present; RA-C12 no duplicated
    canonical schema; RA-C14 required artifact references the canonical
    source. RA-C03 is owned by Task 12; RA-C09 by Task 10; RA-C13 by Task 11.
  - Requirements: 7,14,16 (rule sources 5,8,9,10,12)
  - Design: Validation architecture Layer 2; Consistency rule design
    (detailed); detection boundaries for RA-C12/RA-C14
  - Dependencies: Task 7
  - Completion criteria: each owned RA-C rule emits its stable id at severity
    error and fails with exit 1; RA-C12 and RA-C14 have distinct
    responsibilities; RA-C03, RA-C09, RA-C13 are not implemented here.
  - Verification: run against invalid fixtures (Task 14); assert each owned
    RA-C id and exit 1; verify `^FEATURE-[0-9]{4}$` rejects `XFEATURE-0011`,
    `FEATURE-0011-extra`, `FEATURE-011`, `FEATURE-001A`.

- [ ] 10. Implement RA-C09 version validation and canonical config errors
  - Target: `scripts/reuse-assessment-check.sh`
  - Resolve the canonical version from
    `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`. Treat canonical file
    missing, unreadable, or version missing/malformed as configuration
    errors with exit 2 (never RA-C09). After a valid canonical version is
    resolved, apply RA-C09 (exit 1) to required target version-bearing
    artifacts: each actual reuse assessment (including the FEATURE-0011
    assessment document), the validator version marker, and complete
    fixtures/examples. Link-only artifacts are not version-checked.
  - Requirements: 14,16
  - Design: Validation contract version-resolution flow; Versioning and
    alignment; RA-C09 definition; Error Handling
  - Dependencies: Task 7, Task 9
  - Completion criteria: canonical problems exit 2; target version
    missing/malformed/mismatched fails RA-C09 with exit 1; the canonical
    document itself never fails RA-C09.
  - Verification: fixtures for canonical-config error (exit 2) and target
    version drift (exit 1) via Task 14.

- [ ] 11. Implement RA-C13 evaluation logic (approval enforcement)
  - Target: `scripts/reuse-assessment-check.sh`
  - This task owns only the RA-C13 evaluation logic; it must not implement a
    second Git change-set collector. It accepts the active feature
    identifier; the normalized, deterministic changed-file list supplied by
    the feature gate (Task 13); and the active feature's assessment path as
    resolved by Task 13 via `docs/features/FEATURE_INDEX.md`. It identifies
    implementation-attempt paths in the list (files under `cmd/`,
    `internal/`, `pkg/`, `api/`, `scripts/`, or other recognized executable
    source dirs; documentation and `.kiro/specs/` alone do not count) and,
    when such paths exist for FEATURE-0011+, validates the decision status
    and human-approval evidence in that resolved assessment (Approved status
    plus evidence referencing an approved ADH or assessment-review record
    identifying approver and date). FEATURE-0011's assessment path is valid
    only when the active feature is FEATURE-0011; RA-C13 never validates
    FEATURE-0011's assessment on behalf of a later feature. Returns exit 1
    for approval failures and exit 2 for missing required inputs or
    configuration (for example, no changed-file list, or a resolved
    assessment whose declared feature identity does not match the active
    feature).
  - Requirements: 11,14
  - Design: RA-C13 detection boundary; Change control and approval flow
  - Dependencies: Task 2, Task 7, Task 9 (the resolved assessment path and
    changed-file list are runtime inputs supplied by Task 13; RA-C13 defines
    the input contract and is not built after Task 13)
  - Completion criteria: RA-C13 consumes the gate-supplied list and the
    active feature's resolved assessment without its own Git discovery;
    enforces Approved status and evidence in that assessment for
    implementation-attempt paths; exits 1 on approval failure and exit 2 on
    missing inputs or feature-identity mismatch.
  - Verification (via Task 15): FEATURE-0011 uses its own assessment; a later
    feature (for example FEATURE-0012) uses its own assessment; FEATURE-0012
    cannot pass using FEATURE-0011 approval evidence; an assessment
    feature-identity mismatch exits 2; changed-file lists for committed,
    staged, unstaged, and untracked cases assert exit codes 0, 1, and 2.

- [ ] 12. Integrate RA-C03 with the authoritative scope-text checker
  - Targets: `scripts/reuse-assessment-check.sh`,
    `scripts/phase2-scope-check.sh` (consumed, not modified)
  - RA-C03 is owned solely by this task. It reuses or invokes the existing
    blocked-pattern and allowed-heading logic in
    `scripts/phase2-scope-check.sh` as the authoritative Phase 2 scope-phrase
    source; the validator must not maintain a second blocked-phrase list.
    Future-integration content is flagged only outside deferred-work,
    non-goals, or future-phase headings.
  - Requirements: 7,14
  - Design: detection boundaries (RA-C03); Repository artifact map
  - Dependencies: Task 7, Task 9
  - Completion criteria: RA-C03 relies on the shared scope-check logic; no
    duplicate blocked-phrase list exists in the validator.
  - Verification: scope-placement fixtures (Task 14) produce results
    consistent with `scripts/phase2-scope-check.sh`.

- [ ] 13. Own Git change-set discovery and reconcile/integrate the gate
  - Targets: `scripts/feature-gate.sh`,
    `scripts/reuse-assessment-check.sh` (consumed)
  - Authoritative Git change-set discovery lives here in
    `scripts/feature-gate.sh`: committed changes relative to the phase-branch
    merge base (default `phase2-reuse-first-paas-fabric-foundation`,
    overridable via `PHASE_BRANCH` or gate configuration); staged tracked
    changes; unstaged tracked changes; untracked non-ignored files;
    normalization, de-duplication, and deterministic sorting. Pass that
    exact list to `scripts/reuse-assessment-check.sh`. The gate and the
    validator must not maintain separate implementations of Git change
    discovery.
  - Reconcile the legacy strict checks: the current gate requires Phase 2
    design text containing "Architecture Drift", "Observability",
    "Security", and "Non-goals". Make those heading checks
    applicability-aware so FEATURE-0011 (governance-only) is not forced to
    add irrelevant runtime design sections, while preserving them for
    features where they apply (do not weaken the gate for later runtime
    features). For FEATURE-0011 the gate instead relies on: requirements and
    design artifact presence; the feature-level reuse summary; explicit
    non-goals; the reuse-assessment validator; Phase 2 scope validation; and
    recorded stage and architecture evidence. Automation may validate
    artifact identity, stage labels, controlling references, summaries, and
    non-goals; it must not determine that a Kiro stage was human-approved
    merely because the artifact exists.
  - When the active feature is FEATURE-0011, the gate validates these exact
    Kiro artifacts: `.kiro/specs/reuse-assessment-standard/requirements.md`,
    `.kiro/specs/reuse-assessment-standard/design.md`, and
    `.kiro/specs/reuse-assessment-standard/tasks.md`. For any active feature
    it validates the resolved artifacts for: file presence; active-feature
    identity; correct Requirements, Design, and Tasks stages; reuse-summary
    presence; non-goal preservation; and controlling ADH references.
  - Validator invocation uses the active feature's resolved paths. When the
    active feature is FEATURE-0011, the gate supplies
    `docs/features/FEATURE-0011-reuse-assessment-standard.md`. When the
    active feature is FEATURE-0012 or later, the gate supplies that
    feature's assessment path resolved from `docs/features/FEATURE_INDEX.md`.
    For every FEATURE-0011-and-later execution, the gate supplies to
    `scripts/reuse-assessment-check.sh` in strict mode: the active feature
    identifier; the active feature's resolved assessment path; the active
    feature's resolved Kiro requirements, design, and tasks paths; and the
    authoritative normalized changed-file list. The validator must never
    receive FEATURE-0011's assessment for another active feature. A missing,
    ambiguous, or feature-identity-mismatched resolved path produces a
    configuration diagnostic and exit code 2. Preserve the FEATURE-0001–0010
    legacy exemption unchanged; propagate the validator's non-zero exit as a
    gate failure.
  - Feature-path resolution (single authoritative source): for FEATURE-0011
    and later, `scripts/feature-gate.sh` resolves paths using
    `docs/features/FEATURE_INDEX.md` as the authoritative feature-to-slug
    mapping. Resolution flow: (1) find exactly one index row whose Feature
    field exactly equals the active feature identifier; (2) read and
    validate its Kiro Slug; (3) resolve the Kiro spec directory as
    `.kiro/specs/<kiro_slug>`; (4) resolve the reuse-assessment artifact as
    `docs/features/<FEATURE_ID>-<kiro_slug>.md`; (5) verify the
    requirements, design, tasks, and assessment documents all declare the
    active feature identity. For FEATURE-0011 this resolves exactly to
    `.kiro/specs/reuse-assessment-standard/requirements.md`,
    `.kiro/specs/reuse-assessment-standard/design.md`,
    `.kiro/specs/reuse-assessment-standard/tasks.md`, and
    `docs/features/FEATURE-0011-reuse-assessment-standard.md`. For
    FEATURE-0011 and later, do not use a grep-based Kiro-directory fallback
    and do not fall back to FEATURE-0011 paths; fail with a configuration
    diagnostic and exit code 2 when the index row is missing, duplicated,
    malformed, or resolves to missing or mismatched artifacts. Legacy
    FEATURE-0001 through FEATURE-0010 behavior may remain exempt as approved.
  - The automated gate must not infer human stage approval from file
    presence or content. Requirements-to-Design, Design-to-Tasks, and
    implementation authorization remain separate human-controlled workflow
    decisions under `docs/automation/KIRO_DECISION_POLICY.md`.
  - Exact final-review marker: replace the existing generic case-insensitive
    search for "APPROVED" in the review artifact with exact parsing of the
    field `Final feature-review status: <value>`. The final merge-approval
    gate passes only when the normalized value is exactly `Approved`. It
    must not be satisfied by `Assessment decision status: Approved`, an
    Approved ADH reference, an Approved reuse-summary row, `Final
    feature-review status: Pending`, or any unrelated use of the word
    "Approved". A missing or malformed `Final feature-review status` field
    fails the final merge-approval check.
  - Requirements: 11,13,14,18
  - Design: Feature gate integration; RA-C13 detection boundary
  - Dependencies: Task 2, Task 8, Task 9, Task 10, Task 11, Task 12
  - Completion criteria: the gate owns Git change discovery and supplies the
    normalized list to the validator; no duplicate discovery exists; runtime
    headings are not forced on FEATURE-0011; `docs/features/FEATURE_INDEX.md`
    is the single Phase 2 feature-path resolution source; the gate resolves
    and validates each active feature's own assessment and Kiro artifacts
    with no grep fallback and no FEATURE-0011 fallback; a missing,
    duplicated, or malformed index row, or missing/mismatched resolved
    artifacts, exit 2; legacy features remain exempt; later runtime features
    keep their checks; the gate never infers human stage approval.
  - Verification: `bash scripts/feature-gate.sh FEATURE-0011` resolves via
    the index to `docs/features/FEATURE-0011-reuse-assessment-standard.md`
    and the `reuse-assessment-standard` spec paths without demanding runtime
    headings; a later-feature invocation resolves that feature's own paths
    from the index; an unknown, duplicated, or mismatched feature exits 2;
    `bash scripts/feature-gate.sh FEATURE-0005` skips strict checks; confirm
    the validator receives the gate-supplied changed-file list.

- [ ] 14. Create deterministic fixture data
  - Target: `tests/reuse-assessment/` [NEW]
  - Create deterministic fixture data for: valid and invalid assessment
    Markdown (controlled vocabularies, mandatory fields, adapter rules
    Yes/No and `none`, Build rationale, risk controls, conceptual-example
    labels, feature-identifier anchoring, version drift); DEC/RFC/ADH
    reference records (present and missing); operational artifacts for
    RA-C12 (duplicated schema) and RA-C14 (missing canonical link);
    temporary canonical-document variants (valid, missing version,
    malformed version, unreadable); Phase 2 scope-placement examples; and
    the expected diagnostic records (rule id, severity, exit code) for each
    case. A standalone Markdown fixture alone does not cover every rule;
    rules requiring repository layout or Git state are constructed by the
    Task 15 harness.
  - Requirements: 13,14,16
  - Design: Testing Strategy; Correctness Properties 1–15
  - Dependencies: Task 1
  - Completion criteria: fixture data and expected-diagnostic records exist
    for the file-content rules (RA-S01–RA-S10 and RA-C01–C02, RA-C04–C08,
    RA-C10–C11); complete versioned assessment fixtures declare
    `reuse_assessment_format_version`; fixtures needed for RA-C03, RA-C09,
    RA-C12, RA-C13, and RA-C14 scenarios are provided as inputs the Task 15
    harness assembles.
  - Verification: fixtures are consumed by Task 15 and yield the expected
    rule ids and exit codes.

- [ ] 15. Implement the executable test harness with isolated scenarios
  - Targets: `tests/reuse-assessment/` [NEW] (harness invoking
    `scripts/reuse-assessment-check.sh` and, where needed,
    `scripts/feature-gate.sh`)
  - The harness owns scenario construction and execution for rules that need
    repository layout or Git state: RA-C12 operational duplicated-schema
    scans; RA-C13 changed-file and Git-state behavior; RA-C14 mapped
    canonical-link scans; canonical configuration-error repository layouts;
    and feature-specific assessment and Kiro-spec resolution. It also runs
    the file-content cases from Task 14 fixtures: controlled vocabularies;
    mandatory fields; adapter rules; Build rationale; risk controls;
    DEC/RFC/ADH existence; Phase 2 scope placement; conceptual-example
    labels; version drift; feature-identifier anchoring (reject
    `XFEATURE-0011`, `FEATURE-0011-extra`, `FEATURE-011`, `FEATURE-001A`);
    approval evidence; exit codes 0, 1, and 2; and Phase 1 legacy
    exemptions.
  - Isolation: use `mktemp -d` (or equivalent); copy or construct only the
    fixture files each scenario needs; for RA-C13, initialize an isolated
    phase branch and feature branch, create committed, staged, unstaged, and
    untracked states, and configure the phase branch and merge base inside
    that workspace; clean up with a shell `trap`. The harness must never
    checkout branches, stage, commit, or modify the index in, or otherwise
    modify, the active Sovrunn repository.
  - Feature-path resolution tests (via `docs/features/FEATURE_INDEX.md`):
    FEATURE-0011 resolves its exact assessment and spec paths; a later
    feature resolves its own paths; a feature-identity mismatch fails; a
    missing assessment fails; a missing, duplicated, malformed, or ambiguous
    index row or Kiro spec directory fails; and there is no grep fallback and
    no fallback to FEATURE-0011. Mismatch/missing/duplicated/ambiguous cases
    exit 2.
  - RA-C13 active-feature assessment tests: FEATURE-0011 uses its own
    assessment; FEATURE-0012 uses its own assessment; FEATURE-0012 cannot
    pass using FEATURE-0011 approval evidence; an assessment
    feature-identity mismatch exits 2.
  - Final-review status tests: a pending final review plus an Approved
    assessment does not pass final merge approval; `Final feature-review
    status: Approved` passes; a missing or malformed final-review status
    fails; and unrelated "Approved" text (assessment status, ADH reference,
    or reuse-summary row) does not pass the final merge-approval check.
  - Requirements: 13,14,15,18
  - Design: Testing Strategy; RA-S01–RA-S10; RA-C01–RA-C14; RA-C13 change-set
  - Dependencies: Task 8, Task 9, Task 10, Task 11, Task 12, Task 13, Task 14
  - Completion criteria: together with Task 14, every RA-S01–RA-S10 and
    RA-C01–RA-C14 rule has passing and failing coverage; valid fixtures pass
    (exit 0); invalid fixtures fail with the expected rule id and exit 1;
    canonical-config errors and path-resolution failures exit 2; legacy
    features are exempt; the active repository branch, index, and working
    tree are unchanged after the harness completes.
  - Verification: run the harness (for example `bash
    tests/reuse-assessment/run.sh`); assert expected rule ids and exit
    codes; assert `git rev-parse --abbrev-ref HEAD`, `git status
    --porcelain`, and the index are unchanged before and after the run.

- [ ] 16. Run automated checks and create pending review evidence
  - Targets: `scripts/feature-gate.sh` and
    `scripts/reuse-assessment-check.sh` (executed),
    `docs/reviews/feature-gates/FEATURE-0011-approval-review.md` [NEW]
  - Required sequence: (1) run the automated feature-gate and validator
    checks against
    `docs/features/FEATURE-0011-reuse-assessment-standard.md` before the
    review artifact exists; (2) capture the successful automated evidence
    (commands and results); (3) create the review artifact as pending
    human-review evidence containing that evidence.
  - The review artifact must include: the exact approved FEATURE-0011
    reuse-summary row; the approved Sovrunn-owned and reused/extended
    responsibility statement; a reference to
    `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`; and two distinct
    fields — `Assessment decision status: Approved` and `Final
    feature-review status: Pending`. These are different fields: the reuse
    assessment is Approved through ADH-2026-011, while the final FEATURE-0011
    merge review remains Pending until human review. The artifact must not
    claim merge readiness, and the automated results must not be represented
    as architecture approval. Consistent with the exact final-review marker
    in Task 13, the pending `Final feature-review status: Pending` must not
    satisfy the final merge-approval check. After human review, a separate
    action may change only the final-review status to `Approved`.
  - Requirements: 15,18
  - Design: Feature gate integration; Change control and approval flow
  - Dependencies: Task 1, Task 2, Task 3, Task 4, Task 5, Task 6, Task 13,
    Task 15
  - Completion criteria: automated checks pass and their evidence is captured
    before the artifact exists; the artifact contains the reuse-summary row,
    the responsibility statement, the canonical reference, `Assessment
    decision status: Approved`, and `Final feature-review status: Pending`;
    it makes no merge-readiness claim; the pending artifact cannot satisfy
    the final merge-approval check; the canonical standard remains draft.
  - Verification: capture validator/gate exit codes prior to creating the
    artifact; confirm (1) `Assessment decision status: Approved` is present
    and records the approved assessment, (2) `Final feature-review status:
    Pending` is present and remains Pending, (3) the exact final-review
    status is not `Approved`, and (4) the artifact cannot satisfy the final
    merge-approval check — for example `grep -n "Reuse Assessment Standard
    governance contract"`, `grep -n "Assessment decision status: Approved"`,
    and `grep -n "Final feature-review status: Pending"` on
    `docs/reviews/feature-gates/FEATURE-0011-approval-review.md`; confirm the
    canonical file is not marked Approved.

## Task Dependency Graph

```json
{
  "waves": [
    { "wave": 1, "tasks": ["1"] },
    { "wave": 2, "tasks": ["2", "3", "4", "5", "7", "14"] },
    { "wave": 3, "tasks": ["6", "8", "9"] },
    { "wave": 4, "tasks": ["10", "11", "12"] },
    { "wave": 5, "tasks": ["13"] },
    { "wave": 6, "tasks": ["15"] },
    { "wave": 7, "tasks": ["16"] }
  ],
  "tasks": [
    { "task": "1", "dependsOn": [] },
    { "task": "2", "dependsOn": ["1"] },
    { "task": "3", "dependsOn": ["1"] },
    { "task": "4", "dependsOn": ["1"] },
    { "task": "5", "dependsOn": ["1"] },
    { "task": "6", "dependsOn": ["2"] },
    { "task": "7", "dependsOn": ["1"] },
    { "task": "8", "dependsOn": ["7"] },
    { "task": "9", "dependsOn": ["7"] },
    { "task": "10", "dependsOn": ["7", "9"] },
    { "task": "11", "dependsOn": ["2", "7", "9"] },
    { "task": "12", "dependsOn": ["7", "9"] },
    { "task": "13", "dependsOn": ["2", "8", "9", "10", "11", "12"] },
    { "task": "14", "dependsOn": ["1"] },
    { "task": "15", "dependsOn": ["8", "9", "10", "11", "12", "13", "14"] },
    { "task": "16", "dependsOn": ["1", "2", "3", "4", "5", "6", "13", "15"] }
  ]
}
```

```text
Task 1 (canonical standard, Req 1-18)
  ├── Task 2 (FEATURE-0011 assessment instance) → Task 6 (traceability)
  ├── Task 3 (align prompts + policy)
  ├── Task 4 (align Feature Factory + templates)
  ├── Task 5 (inspect/align RFC + acceptance gates)
  ├── Task 7 (validator entry point + framework)
  │     ├── Task 8 (RA-S01-RA-S10)
  │     └── Task 9 (RA-C01-C02, C04-C08, C10-C12, C14)
  │           ├── Task 10 (RA-C09 + config errors)
  │           ├── Task 11 (RA-C13 evaluation; needs Task 2)
  │           └── Task 12 (RA-C03 scope-check)
  └── Task 14 (fixtures)

Task 13 (Git change-set discovery + gate reconcile/integrate)
  └── needs Tasks 2, 8, 9, 10, 11, 12

Task 15 (test suite, isolated Git state)
  └── needs Tasks 8, 9, 10, 11, 12, 13, 14

Task 16 (automated checks + pending review evidence)
  └── needs Tasks 1, 2, 3, 4, 5, 6, 13, 15
```

Recommended execution order: 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8 → 9 → 10 → 11
→ 12 → 13 → 14 → 15 → 16.

## Notes

- Documentation alignment (Tasks 1–6), validator implementation (Tasks
  7–12), gate reconciliation and integration (Task 13), fixtures (Task 14),
  tests (Task 15), and final automated validation with pending review
  evidence (Task 16) are independently reviewable.
- RA-C03 has exactly one owner (Task 12); RA-C09 is owned by Task 10;
  RA-C13 by Task 11; all other consistency rules by Task 9.
- The FEATURE-0011 assessment document (Task 2) is the actual assessment
  artifact supplied to the validator and gate; the canonical standard
  (Task 1) defines the format only.
- No task selects a vendor, creates a runtime resource, or defers an
  architecture or implementation-engine choice; the approved design is used
  exactly, and all Phase 2 non-goals are preserved.
- Task 16 creates review evidence as pending only; recorded human approval
  and any Approved status transition of the canonical standard remain a
  separate manual step.
- Implementation must not begin without a separate implementation
  authorization after Tasks review.
