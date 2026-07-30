Feature name:
{{FEATURE_ID}} {{FEATURE_TITLE}}

{{MODEL_RECOMMENDATIONS}}

Start with requirements.md only.

Do not generate design.md yet.
Do not generate tasks.md yet.
Do not implement code.
Do not modify source files except the requirements file under:
{{SPEC_PATH}}/requirements.md

Tool-output safety constraints:
- Use fs_write only for chunks of 50 lines or fewer.
- For files longer than 50 lines, create the file with fs_write using the first chunk, then use fs_append in chunks of 50 lines or fewer.
- Do not write the entire requirements.md in one fs_write call.
- Do not use one very large str_replace edit.
- Split content into logical sections.
- Write one section at a time.
- After writing, read the file back and verify it is complete.

Context:
Sovrunn is an AI-first sovereign cloud-native PaaS platform.
This feature belongs on branch {{FEATURE_BRANCH}} and must use the phase and scope resolved from docs/features/FEATURE_INDEX.md and the active feature assessment; do not default Phase 2 work to Phase 1.

Use these repo context files:
- AGENTS.md
- README.md
- docs/foundation/constitution.md
- docs/decisions/DECISION_INDEX.md
- docs/glossary.md
- docs/features/FEATURE_SEQUENCE.md
- docs/resource-specs/RESOURCE_MODEL_PHASE1.md
- docs/api/API_CONTRACT_PHASE1.md
- docs/engineering/ai-context-loading-standard.md
- docs/engineering/go-coding-guardrails.md
- docs/features/FEATURE_INDEX.md
- the active feature assessment resolved from FEATURE_INDEX.md
- the controlling ADH referenced by that assessment
- any canonical architecture document referenced by the assessment or ADH
- docs/phase2/PHASE2_FEATURE_SEQUENCE.md and docs/phase2/PHASE2_EXECUTION_STRATEGY.md for Phase 2 features

## FEATURE-0013 consolidated architecture boundary

When `{{FEATURE_ID}}` is `FEATURE-0013`, load all of these exact controlling
inputs before writing:

- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- the single approved
  `docs/reviews/architecture-decision-handoffs/ADH-2026-017-feature-0013-*.md`
  consolidation handoff
- `docs/features/FEATURE-0011-reuse-assessment-standard.md`
- `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
- `docs/reviews/feature-gates/FEATURE-0011-approval-review.md`
- `docs/architecture/api-resource-standard.md`
- `.kiro/specs/api-resource-naming-status-and-validation-standard/requirements.md`
- `.kiro/specs/api-resource-naming-status-and-validation-standard/design.md`
- `docs/reviews/feature-gates/FEATURE-0012-approval-review.md`

The consolidated architecture main body is authoritative and ADH-2026-017 is
its sole approval envelope. Do not load ADH-2026-014, ADH-2026-015, or
ADH-2026-016 as instructions or parallel authority; they are historical
provenance only. Do not load a superseded
FEATURE-0013 requirements or design draft as controlling context.

Treat architecture section 27.8 as a closed decision boundary. Requirements
must translate those decisions into testable obligations and must not reopen
them as design questions. A missing semantic value or conflict must be returned
as `ARCHITECTURE_DECISION_REQUIRED`; do not invent or choose it.

Preserve the closed limit hierarchy explicitly: FEATURE-0012 platform/schema
ceilings are absolute outer bounds; a `DecisionProfile` may only narrow them;
an evaluator registration may only narrow the profile; and the effective limit
is the most restrictive applicable value. Do not invent FEATURE-0013 numeric
defaults.

Preserve architecture section 15's public error contract exactly. Requirements
must reproduce the complete closed `violations[].code` registry and map
testable cases plus RFC 6901 fields to those existing codes. They must not
invent leaf suffixes, aliases, type URIs, top-level codes, or HTTP meanings.

Preserve the request/outcome boundary exactly. `DecisionRequest` is a
calling-domain-owned conceptual input role, not a FEATURE-0013 shared contract;
it only conforms to inherited `TransientRequestResult` boundaries. Do not
require a common request schema, Go type, registry entry, resource, or
persistence model. Completed handling returns `DecisionRecord`, accepted
asynchronous handling returns the existing FEATURE-0012 `Operation`, and
rejected or unaccepted handling returns inherited Problem Details. Do not
define `PendingDecision`, `DeferredDecision`, `DecisionPending`, a pending
decision status, or any other FEATURE-0013 non-final response envelope.
The generic `Operation` contract and `LongRunningOperation` payload remain
owned by FEATURE-0012/ADH-2026-013; FEATURE-0013 must not redefine them.

Preserve risk-identifier ownership. FEATURE-0012 `F12-R01` through `F12-R16`
remain FEATURE-0012 identifiers. FEATURE-0013 uses only its separate `F13-R01`
through `F13-R31` namespace. Reuse the governance approach, but do not merge,
renumber, alias, inherit, or ordinally map the two identifier sets.

Preserve the immutable erasure boundary from architecture sections 6.5-6.7:
`DecisionRecord` and `AuditEvent` are append-only; correction, revocation, and
erasure never rewrite the canonical record. Requirements must cover separately
controlled payload custody, later key destruction, and linked
tombstone/redaction records as structural alternatives, while authorizing no
erasure or cryptographic runtime in FEATURE-0013.

The generated document must contain a level-two `Architecture traceability`
section mapping every requirement family and every retained design question to
the exact controlling architecture section. At minimum, it must cover sections
5.4, 6.1, 6.8, 7.1, 9, 12.3, 12.4, 15, 17, and 27.8.

For FEATURE-0013, “complete” is deterministic rather than a summary claim. The
Architecture traceability section must contain every top-level architecture
section reference `section 1` through `section 28`, every exact decision ID
`AD-001` through `AD-045`, every exact risk ID `F13-R01` through `F13-R31`, and
every stable conformance ID: `F13-CF-01` through `F13-CF-28`,
`F13-SCOPE-01` through `F13-SCOPE-11`, `F13-EVAL-01` through `F13-EVAL-08`,
`F13-SEC-01` through `F13-SEC-08`, `F13-TRUST-01` through `F13-TRUST-04`, and
`F13-COMPAT-01` through `F13-COMPAT-09`. Each entry must state whether it is
translated now, preserved as a later invariant, deferred, delegated to design,
or governance-only. Ranges such as `AD-001–AD-045` do not substitute for the
individual IDs.

Implementation class is architecture-owned, not a requirements decision. For
each `AD-001` through `AD-045`, include exactly one explicit traceability row
whose implementation class exactly matches the section 19 scorecard:
`CONTRACT_NOW`, `INVARIANT_FOR_LATER`, or `DEFERRED`. Requirements must not
reclassify, combine, qualify, or omit that value. A mismatch must stop with
`ARCHITECTURE_DECISION_REQUIRED`.

Requirements must include the exact active feature identity and `Stage: Requirements`, then:
1. Introduction
2. Glossary if new concepts are introduced
3. User stories
4. Acceptance criteria
5. Non-goals
6. Edge cases
7. Security/privacy requirements
8. Compatibility with already completed Phase 1 features
9. Design questions to resolve later in design.md
10. Architecture traceability

Keep requirements concise, precise, implementation-aware, phase-scoped, and free of scope creep.

## FEATURE-0014 closed architecture boundary

When `{{FEATURE_ID}}` is `FEATURE-0014`, load the exact approved
`docs/architecture/provider-neutral-resource-model.md` and
`docs/reviews/architecture-decision-handoffs/ADH-2026-018-feature-0014-provider-neutral-resource-model.md`
and its approved geographic-descriptor replacement clarification
`docs/reviews/architecture-decision-handoffs/ADH-2026-019-feature-0014-geographic-descriptor-clarification.md`
before writing. Run `make feature-0014-architecture-readiness`; a failure is
`REPOSITORY_CONTEXT_NOT_READY`, not permission to repair or generate partial
requirements.

Treat architecture sections 2.1-2.4 and 13-16 as closed generation controls.
Requirements must enumerate every `F14-AD-001` through `F14-AD-021` and every
`F14-R01` through `F14-R30` individually; range shorthand is insufficient.
Every normative requirement must have exactly one owner and one canonical
semantic key `(owning feature, resource or contract, actor, behavior,
observable outcome)`. Reference shared obligations instead of copying them.

FEATURE-0014 owns only Provider, ProviderLocation, ProviderDatacenter,
DatacenterFailureDomain, and InfrastructureStack. It must not define or infer
ResourcePool, ProviderCapability, capacity, compatibility, placement
eligibility, connectivity, adapters, discovery, credentials, provider calls,
plugins, operations, provisioning, runtime execution, or provider-native core
identity. FEATURE-0013 adoption is exactly `NOT_APPLICABLE`; do not invent an
explainable decision, AuditEvent behavior, operation behavior, or decision
projection to satisfy a generic Phase 2 gate.

Requirements may translate closed decisions into testable outcomes only. They
must not choose routes, field layout, packages, storage, algorithms, internal
types, implementation files, or other design mechanics. An unresolved semantic
choice stops the entire stage with `ARCHITECTURE_DECISION_REQUIRED`.

Do not implement code.
Do not generate design.md.
Do not generate tasks.md.

## Phase 2 Reuse and Drift Gates

Every FEATURE-0011-and-later feature must include a reuse assessment that
conforms to the canonical standard:

`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

Do not duplicate or redefine the assessment field schema in this document.
Populate the feature-level reuse summary and capability-level assessments
using the canonical fields and controlled vocabularies.

Architecture drift checks:

- no provider-specific hardcoding in core,
- no Kubernetes-only assumptions in core,
- no PostgreSQL lifecycle logic in core placement engine,
- no custom policy engine embedded in handlers,
- no raw secret storage,
- no customer-facing IaaS leakage,
- explainable `DecisionRecord` where the approved feature applicability requires one,
- defined audit behavior where the approved feature applicability requires it,
- preserved adapter boundaries.
