Feature:
{{FEATURE_ID}} {{FEATURE_TITLE}}

Generate design.md only.

Do not generate tasks.md yet.
Do not implement code.
Do not modify source files.
Do not change requirements.md unless you find a blocking contradiction.

Input requirements:
{{REQUIREMENTS_PATH}}

Tool-output safety constraints:
- Use fs_write only for chunks of 50 lines or fewer.
- For files longer than 50 lines, create the file with fs_write using the first chunk, then use fs_append in chunks of 50 lines or fewer.
- Do not write the entire design.md in one fs_write call.
- Do not use one very large str_replace edit.
- Split content into logical sections.
- Write one section at a time.
- After writing, read the file back and verify it is complete.

Use these repo context files:
- AGENTS.md
- README.md
- docs/engineering/go-coding-guardrails.md
- docs/engineering/ai-context-loading-standard.md
- existing implementations for similar resources

For FEATURE-0013, load the approved consolidated architecture and only its
approved ADH-2026-017 consolidation handoff. ADH-2026-014/015/016 are
historical provenance, not instructions.
Architecture section 27.8 and section 27.9 are closed. Resolve only
representation and implementation mechanics explicitly delegated to design. If
requirements ask design to choose a closed semantic, ownership, vocabulary,
mapping, public error family, resource profile, scope authority, conformance
identity, or cryptographic selection, stop with
`ARCHITECTURE_DECISION_REQUIRED` instead of resolving it.

FEATURE-0013 section 27.9 anti-wandering controls are mandatory:
- `SecurityExceptionRef` is the only approved fail-open exception evidence
  representation; strategy metadata is not an independent fail-open authority.
- `GraphEdge` requires a stable edge ID distinct from `from`/`to`; graph version
  validation belongs to bundle graph resolution before in-memory graph build.
- Relationship conflict behavior must be deterministic and use only approved
  `DECISION_RELATIONSHIP_*` codes.
- Reuse the FEATURE-0012 baseline workflow exactly: `BASELINE_MANIFEST.json`
  and `BASELINE_APPROVALS.json`; do not introduce parallel baseline paths.
- FEATURE-0013 AuditEvent domain extension types live in a FEATURE-0013
  domain/shared package; `apiconform` is binding/conformance support only.
- Versioned registries use `(id, version)` keys unless explicitly unversioned.
- Present-empty `TrustCarrier` never satisfies provenance, required trust,
  identity generation, or compatibility evidence.
- Projection pointers resolve against the profile-declared canonical view;
  overlap includes equality and ancestor/descendant containment.
- Matrix E automation must not calculate, infer, downgrade, or upgrade residual
  risk levels.
- JSON pre-scan uses JSON tokens; YAML pre-scan uses `yaml.Node` before
  normalization.
Do not read, copy, summarize, or reuse
`.kiro/specs/decision-object-and-auditevent-standard/design.superseded-patched-draft.md`
as a design source.

Preserve the inherited limit hierarchy in every relevant component:
FEATURE-0012 platform/schema ceiling >= DecisionProfile ceiling >=
evaluator-specific ceiling; the most restrictive applicable value controls.

Preserve architecture section 15's complete public error binding and closed
violation-code registry. Design may map every approved code to validators and
safe RFC 6901 fields, but may not add, rename, alias, or reinterpret a public
code, type URI, top-level code, or HTTP meaning.

Preserve the closed request/outcome boundary. `DecisionRequest` remains a
calling-domain-owned conceptual input and has no FEATURE-0013 shared schema,
Go type, registry entry, resource, or persistence model. Design may show how a
domain-owned input conforms to `TransientRequestResult`, but may not create the
shared contract. Completed handling uses `DecisionRecord`; asynchronous
acceptance uses the existing FEATURE-0012 `Operation`; rejection uses inherited
Problem Details. Do not design a pending-decision response type or status.
The generic `Operation` contract and `LongRunningOperation` payload remain
owned by FEATURE-0012/ADH-2026-013; design may reference but not redefine them.

Keep risk identifiers disjoint. FEATURE-0012 owns `F12-R01` through `F12-R16`;
FEATURE-0013 owns `F13-R01` through `F13-R31`. Design may reuse the risk-control
method but must not merge, renumber, alias, inherit, or ordinally map IDs.

Preserve architecture sections 6.5-6.7: canonical `DecisionRecord` and
`AuditEvent` objects are append-only and are never rewritten for correction,
revocation, retention, or lawful erasure. Design may model only a linked
tombstone/redaction record and separately controlled payload-custody boundaries;
key destruction and erasure execution remain deferred.

Design must include overview, resolved implementation decisions, architecture,
files, data models, interfaces, validation, API/handler design where applicable,
registry/storage design where applicable, operation/audit behavior, error
mapping, security/privacy, testing, verification, non-goals, resolved design
questions, and a level-two `Architecture traceability` section. For
FEATURE-0013, that traceability must cite architecture sections 5.4, 6.1, 6.8,
7.1, 9, 12.3, 12.4, 15, 17, and 27.8.

FEATURE-0013 traceability must also enumerate every individual top-level
architecture section `section 1` through `section 28`, `AD-001` through
`AD-045`, `F13-R01` through `F13-R31`, and every stable conformance ID in the
architecture (`F13-CF-01` through `F13-CF-28`, `F13-SCOPE-01` through
`F13-SCOPE-11`, `F13-EVAL-01` through `F13-EVAL-08`, `F13-SEC-01` through
`F13-SEC-08`, `F13-TRUST-01` through `F13-TRUST-04`, and `F13-COMPAT-01`
through `F13-COMPAT-09`). Each exact
ID must appear; range shorthand is insufficient. Record the design disposition
for each without reopening its semantics.

Implementation class is architecture-owned. The `Architecture traceability`
section must contain exactly one explicit row for each `AD-001` through
`AD-045` and reproduce its exact section 19 class (`CONTRACT_NOW`,
`INVARIANT_FOR_LATER`, or `DEFERRED`). Design must not reclassify an AD or
create components for `INVARIANT_FOR_LATER` or `DEFERRED` rows.

Hard constraints:
- Follow docs/engineering/go-version-standard.md.
- Do not introduce external dependencies unless requirements explicitly demand them.
- internal/api must not import internal/server.
- Do not add unrelated future scope.

For FEATURE-0014, load only the approved architecture and ADH-2026-018 plus
ADH-2026-019 replacement-clarification package
plus approved requirements and inherited context allowed by architecture
section 2.2. Design may choose only representation mechanics delegated by
architecture sections 14 and 15. It must enumerate every F14 decision and risk,
map every requirement to representation/validation/evidence, preserve the
single-owner and normalization ledgers, and add no uncited normative behavior.
Do not design ResourcePool, ProviderCapability, connectivity, adapters,
discovery, provider calls, decision/audit/operation behavior, or speculative
extension points. A semantic gap stops the entire stage with
`ARCHITECTURE_DECISION_REQUIRED`.

Generate design.md only.

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
