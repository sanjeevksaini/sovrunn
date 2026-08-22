# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-066
- Date: 2026-08-22
- Source discussion: Codex architecture session
- Related feature: FEATURE-0016
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Add the F0016-owned coherent backing-access compatibility bridge

## Summary

FEATURE-0016 must safely resolve `CloudProviderParticipation` and
`InfrastructureStack` together, evaluate caller-relative access, and retain a
backing fence through its audited qualification publication. The existing
FEATURE-0015 individual `GetParticipation` and `GetInfrastructureStack`
operations cannot supply that guarantee. This handoff authorizes one narrow,
internal compatibility bridge co-located with the existing in-memory store but
owned as FEATURE-0016 integration mechanics. It does not reopen FEATURE-0015
as a product feature and changes no FEATURE-0015 public behavior.

## Classification

Extension: a private compatibility mechanism required to realize existing
FEATURE-0016 safety guarantees, with no new observable product behavior.

## Existing approved baseline

- FEATURE-0015 owns `CloudProviderParticipation` and `InfrastructureStack`.
  FEATURE-0016 consumes them read-only.
- FEATURE-0016 owns `ExecutionTarget`, qualification, safe projection,
  idempotency, and sole target-state publication.
- ADH-2026-058 requires safe backing access, effective-Active participation,
  Active InfrastructureStack, fence recheck, and stale-work non-publication.
- FEATURE-0012 remains the authority for authentication, authorization,
  safe-denial envelope, and top-level Problem codes.
- FEATURE-0013 remains the authority for AuditEvent append-before-publication.

## Decision

### 1. Ownership and non-reopening rule

The bridge is an F0016 integration implementation detail. Its source may be
co-located in the existing `internal/cloudmodel` store package solely because
that package owns the in-memory synchronization protecting the two backing
resources. This placement does **not** transfer resource ownership and does
**not** create new FEATURE-0015 product scope.

No FEATURE-0015 requirements, design, tasks, control manifest, feature
authority, conformance, route, schema, writer, lifecycle, or approval record
is reopened or changed by this handoff.

### 2. Closed internal bridge contract

F0016 owns an internal `BackingAccessProvider`. It accepts the authenticated
principal, F0016 route intent, the participation UID, the InfrastructureStack
UID, and the existing FEATURE-0012 grant evaluator. It is the sole
principal-aware safe-resolution component for F0016 backing reads.

The provider uses one private store-supported backing-read lease that, under
one existing store lock acquisition, supplies immutable copies of both backing
resources to the F0016 provider. The provider, not FEATURE-0015, evaluates the
existing F0012 grant against the derived CloudProvider UID while the lease is
held. It returns only one of these internal results:

- safe denial, mapped by F0016 to the existing audited safe 404;
- authorization denial, mapped by F0016 to the existing audited 403;
- an immutable allowed snapshot containing only:
  - participation UID and derived CloudProvider UID;
  - participation effective-Active value;
  - InfrastructureStack UID, Active value, and generation;
  - viability fingerprint inputs composed from both UIDs and viability values.

No raw backing-existence result crosses the F0016 handler boundary. The bridge
does not persist a scope field or index on `ExecutionTarget`; any process-local
derived-scope cache is non-authoritative and restart-cleared.

### 3. Qualification fence through publication

Initial create/qualify validation obtains an immutable allowed snapshot through
the provider. Qualification observation happens with no backing lease held.

For final qualification commit, F0016 obtains a fresh allowed backing lease,
compares only the registered fences:

- referenced InfrastructureStack generation; and
- viability fingerprint composed from the participation UID/effective-Active
  value and InfrastructureStack UID/phase.

If either differs, F0016 aborts without AuditEvent, publication, or replayable
completion using the existing registered stale outcome. Participation
generation is not a fence and creates no outcome.

If fences match, the same backing lease remains held through F0016's existing
AuditEvent append-before-publication critical section. The fixed lock order is:

```text
FEATURE-0015 store backing-read lease
    → FEATURE-0016 lifecycle-service mutex
        → inherited audit append
            → F0016 target publication
```

No code entered while holding the F0016 lifecycle mutex may acquire the
FEATURE-0015 store lease. FEATURE-0015 resource writers never acquire the
F0016 lifecycle mutex. This is the only permitted cross-package lock order.

### 4. Read projections and prospective admission clarification

Item GET and LIST acquire a fresh F0016 `BackingAccessProvider` result outside
the lifecycle mutex after their approved authorization/safe-resolution stages.
They use it only to compute response-only `effectiveAvailability`; no ETag or
persisted target state changes.

`CloudProviderParticipation` suspension is a prospective admission-control
condition. It prevents new ExecutionTarget qualification and future resource
placement/provisioning through that participation. It does not command, imply,
or cause termination, suspension, draining, deletion, or lifecycle mutation
of an existing InfrastructureStack, ExecutionTarget, or already-realized
workload. In FEATURE-0016, `effectiveAvailability: Unavailable` means
ineligible for new admission, never that existing workloads have stopped.

### 5. Explicit non-goals

This handoff does not authorize:

- a public FEATURE-0015 or FEATURE-0016 route;
- a resource schema, field, state, writer, status, or conformance-case change;
- a FEATURE-0015 lifecycle, access-control, audit, or authorization redesign;
- external persistence, external calls, real adapter behavior, provisioning,
  placement, plugin execution, or customer IaaS exposure;
- a participation-generation fence or any ninth F0016 public violation.

## Rationale

An F0016-only mutex cannot prevent a concurrent FEATURE-0015 backing mutation.
The bridge reuses the existing store's synchronization without altering its
resource model or external API. A single lease gives F0016 one verified view
and a valid final fence, while the principal-aware F0016 provider preserves the
existing non-disclosure behavior.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend the existing in-memory backing-store synchronization;
  reuse FEATURE-0012 grant/safe-denial and FEATURE-0013 audit contracts.
- Summary of mature candidates / applicable standards: Existing
  `internal/cloudmodel.Store` lock and copy-returning model; no external
  library, database transaction, or provider adapter is selected.
- Sovrunn-owned responsibility summary: F0016 owns principal-aware backing
  access composition, fence comparison, lifecycle publication, and safe
  projection. FEATURE-0015 continues to own backing resources only.
- Non-goals summary: No public API, resource model, lifecycle, execution, or
  external-effect expansion.

## Phase impact

- Current phase allowed? Yes.
- Current phase boundary impact: In-process, deterministic, read-only reuse
  only. No real provider operation or external side effect is introduced.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting decisions: none.
- Resolution required: approved ADH application; no DEC, RFC, ACR, or baseline
  update.

## Required action

- Update FEATURE-0016 architecture and feature authority with this private
  bridge and prospective-admission clarification.
- Update only direct F0016 traceability, closure, control, and fail-closed
  readiness/contract assertions.
- After the architecture update, regenerate/review F0016 design and tasks.
- Implement later as one F0016 remediation task with exact code/test paths.

## Strict Kiro read/write boundary

| Path | Access | Purpose |
|---|---|---|
| `docs/reviews/architecture-decision-handoffs/ADH-2026-066-feature-0016-coherent-backing-access-compatibility-bridge.md` | read | Controlling approved decision |
| `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md` | write | Record bridge, lock order, and admission clarification |
| `docs/features/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md` | write | Record F0016 integration ownership and reuse disposition |
| `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md` | write | Add ADH-066 provenance only to directly affected F0016 evidence |
| `docs/reviews/architecture-readiness/FEATURE-0016-architecture-closure-matrix.md` | write | Record resolved backing-fence closure |
| `.automation/features/FEATURE-0016.control.json` | write | Add ADH-066 only to `feature.handoffs` |
| `scripts/feature-0016-architecture-readiness-check.py` | write | Fail closed on sequential backing reads, unregistered fence, or cross-lock reversal |
| `scripts/feature-contract-check.py` | write | Confirm bridge preserves F0016 ownership and registered fences |
| `scripts/vs000-contract-check.py` | write | Confirm no registry semantics or F0015 ownership drift |

The architecture-update pass must not modify FEATURE-0015 documentation,
control evidence, requirements, design, tasks, source, tests, registry,
contract specification, formal models, Structurizr, or any F0016
requirements/design/tasks artifact.

## Implementation scope reservation

After an approved F0016 design/task remediation, the only source paths that
may be added to the F0016 remediation task are:

- `internal/cloudmodel/store.go` and `internal/cloudmodel/store_test.go` for
  the private backing-read lease only;
- `internal/executiontarget/` for the F0016-owned provider, lifecycle, and
  scheduler integration;
- `internal/api/executiontarget_*.go` and corresponding tests for F0016
  handler wiring only;
- `tests/conformance/feature_0016_test.go` for exact F0016 proof.

No other FEATURE-0015 implementation path may be edited. The remediation task
must run both FEATURE-0015 regression and FEATURE-0016 feature gates.

## Acceptance criteria for Kiro update

- [ ] FEATURE-0015 remains unchanged in all feature authority, spec, control,
  conformance, and review artifacts.
- [ ] F0016 owns the principal-aware provider; the store-supported lease is a
  private, resource-model-neutral integration primitive.
- [ ] The only final backing fences are InfrastructureStack generation and
  viability fingerprint.
- [ ] The backing lease is held through audit-before-publication and obeys the
  stated acyclic lock order.
- [ ] Suspension means unavailable for new admission, not terminated existing
  workloads.
- [ ] No public route, schema, field, state, violation, conformance behavior,
  formal model, or Phase 2R boundary changes.
- [ ] Readiness/contract checks fail closed if F0016 returns to sequential
  backing reads, adds a participation-generation fence, or reverses lock order.

## Explicit instructions to Kiro

- Validate this handoff against the Architecture Operating System before any
  edit.
- Apply only the listed write allowlist. Do not modify FEATURE-0015 artifacts.
- Do not edit `requirements.md`, `design.md`, `tasks.md`, Go source, tests,
  registry semantics, formal models, or Structurizr in this pass.
- Stop with `ARCHITECTURE_DECISION_REQUIRED` if any additional observable
  behavior is necessary; stop with `DEPENDENCY_APPROVAL_REQUIRED` if this
  closed private bridge cannot be represented without a FEATURE-0015 product
  artifact change.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-22
- Notes: Founder approved the strict private compatibility bridge. FEATURE-0015
  remains closed as a product feature; no follow-up FEATURE-0015 feature loop
  or public behavior change is permitted.
