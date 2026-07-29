Feature:
{{FEATURE_ID}} {{FEATURE_TITLE}}

Generate tasks.md only.

Do not implement code.
Do not modify source files.
Do not modify requirements.md.
Do not modify design.md unless you find a blocking contradiction.

Inputs:
- {{REQUIREMENTS_PATH}}
- {{DESIGN_PATH}}
- AGENTS.md
- docs/engineering/go-coding-guardrails.md
- docs/engineering/go-version-standard.md
- docs/engineering/go-observability-standard.md

For FEATURE-0013, load the approved consolidated architecture and only its
approved ADH-2026-017 consolidation handoff. ADH-2026-014/015/016 are
historical provenance, not instructions.
Architecture section 27.8 and section 27.9 are closed. Tasks may implement only
approved `CONTRACT_NOW` design components and must not introduce a semantic
choice, runtime for an `INVARIANT_FOR_LATER` item, or work for a `DEFERRED`
item. Stop with `ARCHITECTURE_DECISION_REQUIRED` if the approved design is
incomplete or contradicts the closed architecture.

FEATURE-0013 section 27.9 anti-wandering controls are mandatory for tasks:
- Plan implementation of `SecurityExceptionRef` only as structural evidence; do
  not create a security-exception approval workflow or service.
- Preserve strategy fail-open gating through the governing profile
  `SecurityExceptionRef`; no independent strategy approval path.
- Plan `GraphEdge` with stable edge ID distinct from `from`/`to`; graph version
  validation belongs to bundle resolution before in-memory graph build.
- Use only approved relationship codes and the deterministic relationship
  conflict matrix from design.
- Reuse `api/schemas/baseline/BASELINE_MANIFEST.json` and
  `api/schemas/baseline/BASELINE_APPROVALS.json`; do not create
  `api/schemas/SCHEMA_BASELINE_MANIFEST.json`, `api/schemas/diffs/`, or
  `api/schemas/approvals/`.
- Keep FEATURE-0013 AuditEvent domain extension types in the FEATURE-0013
  domain/shared package; `apiconform` is binding/conformance support only.
- Use `(id, version)` registry keys where design says versioned.
- Enforce the design-selected present-empty `TrustCarrier` rule.
- Keep projection conformance structural only; no redaction/projection runtime.
- Matrix E tasks may stage architecture-owned values and blank human fields only;
  do not calculate, infer, downgrade, or upgrade residual risk.
- Reuse FEATURE-0012 decode/pre-scan primitives where package DAG permits; YAML
  checks occur on `yaml.Node` before normalization.
- Do not use superseded or patched design drafts as task input.

Tasks must preserve and test the inherited limit hierarchy: FEATURE-0012
platform/schema ceiling >= DecisionProfile ceiling >= evaluator-specific
ceiling, with the most restrictive applicable value controlling.

Tasks must implement and test architecture section 15's complete closed public
violation-code registry and inherited Problem Details binding. No task may
create or reinterpret a public code, type URI, top-level code, or HTTP meaning.

Tasks must test append-only correction/revocation and linked
tombstone/redaction record semantics, and must not implement in-place
canonical-record mutation, payload deletion, key destruction, or an erasure
runtime. Separately controlled payload custody and cryptographic erasure remain
later-feature boundaries.

Tasks must not create a shared FEATURE-0013 `DecisionRequest` schema/type,
registry entry, resource, or persistence model. This calling-domain-owned
conceptual input conforms only to inherited `TransientRequestResult`
boundaries. Tasks must use only
`DecisionRecord` for completion, existing FEATURE-0012 `Operation` for accepted
asynchronous handling, and inherited Problem Details for rejection; no
pending-decision response envelope or status may be implemented.
The generic `Operation` contract and `LongRunningOperation` payload remain
owned by FEATURE-0012/ADH-2026-013; no FEATURE-0013 task may redefine them.

Tasks must preserve disjoint risk namespaces: FEATURE-0012 owns `F12-R01`
through `F12-R16`, while FEATURE-0013 owns `F13-R01` through `F13-R31`. No task
may merge, renumber, alias, inherit, or ordinally map identifiers between them.

Tool-output safety constraints:
- Use fs_write only for chunks of 50 lines or fewer.
- For files longer than 50 lines, create the file with fs_write using the first chunk, then use fs_append in chunks of 50 lines or fewer.
- Do not write the entire tasks.md in one fs_write call.
- Do not use one very large str_replace edit.
- Split content into logical sections.
- Write one section at a time.
- After writing, read the file back and verify it is complete.

Task generation rules:
- Create implementation tasks only.
- Each task must be small enough for Cursor to implement safely.
- Prefer one focused code area per task.
- Add tests close to the code they verify.
- Do not combine model, registry, handler, server wiring, and integration tests into one task.
- Each task must include objective, files, notes, tests, acceptance criteria, and commit message.
- Each task must cite its controlling requirement, design section, architecture
  section, and implementation class. Include a level-two `Architecture traceability`
  section covering architecture sections 5.4, 6.1, 6.8, 7.1, 9,
  12.3, 12.4, 15, 17, and 27.8.
- For FEATURE-0013, that traceability must enumerate every individual
  top-level architecture section `section 1` through `section 28`, `AD-001`
  through `AD-045`, `F13-R01` through `F13-R31`, and every stable conformance
  ID in the architecture (`F13-CF-01` through `F13-CF-28`, `F13-SCOPE-01`
  through `F13-SCOPE-11`, `F13-EVAL-01` through `F13-EVAL-08`, `F13-SEC-01`
  through `F13-SEC-08`, `F13-TRUST-01` through `F13-TRUST-04`, and
  `F13-COMPAT-01` through `F13-COMPAT-09`). Every exact ID must appear; range shorthand is
  insufficient. Mark non-task-producing entries as invariant, deferred, or
  governance-only rather than inventing work.
- Implementation class is architecture-owned. Include exactly one explicit
  traceability row for every `AD-001` through `AD-045`, reproducing the exact
  class from architecture section 19. Generate implementation tasks only for
  `CONTRACT_NOW`; `INVARIANT_FOR_LATER` and `DEFERRED` rows produce no source
  task. Only `CONTRACT_NOW` items produce implementation tasks.
- Final task must include full Docker verification, guardrails, artifact cleanup, and clean git status.

Standard Docker verification command:

docker run --rm -v "$PWD":/src -w /src ${GO_DOCKER_IMAGE:-golang:1.22} sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'

Final Docker verification command:

docker run --rm -v "$PWD":/src -w /src ${GO_DOCKER_IMAGE:-golang:1.22} sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./... && go build ./cmd/sovrunn-api'

Final guardrails:
- rm -f sovrunn-api
- rm -rf bin
- no TODO({{FEATURE_ID}}) under internal or cmd
- no internal/api import of internal/server
- git status clean

For FEATURE-0014, tasks may implement only approved design elements traced to
approved requirements and the ADH-2026-018 closed decisions. Tasks must preserve
all F14 decision/risk/evidence traceability and leave no semantic choice to the
implementer. Do not add FEATURE-0015 ResourcePool/ProviderCapability,
FEATURE-0016 adapter/integration, FEATURE-0053 connectivity, or FEATURE-0013
decision/audit/operation work. A missing or contradictory design mapping stops
the entire stage with `ARCHITECTURE_DECISION_REQUIRED`.

Generate tasks.md only.

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
- defined observability behavior,
- request/operation correlation where applicable,
- no secret or credential logging,
- preserved adapter boundaries.
