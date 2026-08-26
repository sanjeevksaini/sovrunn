# FEATURE-0017 Policy Evaluation Abstraction — Tasks

## 1. Identity, stage, and execution rules

| Field | Value |
|---|---|
| Feature | FEATURE-0017 — Policy Evaluation Abstraction |
| Stage | Tasks |
| Kiro slug | `policy-evaluation-abstraction` |
| Phase / order | Phase 2R / 7 |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Sole architecture authority | `docs/architecture/policy-evaluation-abstraction.md` |
| Controlling handoffs | ADH-2026-067, ADH-2026-068, ADH-2026-069 |
| Approved requirements | `.kiro/specs/policy-evaluation-abstraction/requirements.md` (REQ-F17-01..06, AC-F17-01..24) |
| Approved design | `.kiro/specs/policy-evaluation-abstraction/design.md` (D-01..D-14) |
| Owned schemas | `VS0-SCHEMA-018` (`PolicyEvaluationRequest`), `VS0-SCHEMA-019` (`PolicyEvaluationResult`) |
| Public routes / persistence / controller / real engine | None |

This document decomposes the approved design into independently testable
implementation units. Each implementation task owns concrete paths, tests,
verification commands, acceptance criteria, and a commit message; the final
read-only verification task intentionally owns neither paths nor a commit. Tasks
do not make design or architecture decisions. Every task cites exact
requirements, design sections, and architecture decisions. Tasks sequence real
dependencies and mark independent work explicitly.

Cursor may execute from this `tasks.md` only after the independent tasks-stage
approval and the separately recorded founder executable-plan approval both
exist. This document neither records nor self-asserts either approval. Tasks may
not add semantics, public routes, stores, controllers, production engines,
external effects, or adjacent-feature behavior absent from the approved design.

### Cursor execution prerequisites

Before coding any task, Cursor must follow the complete implementation-context
and verification contract in `AGENTS.md`. The following task-local list is
additive and is not a replacement for any context required by `AGENTS.md`:

1. `AGENTS.md`
2. `docs/engineering/go-coding-guardrails.md`
3. `docs/engineering/go-observability-standard.md`
4. This `tasks.md` file
5. The approved `design.md` for FEATURE-0017

Cursor must not perform cleanup, refactoring, or fixes outside the current
task's writable-path allowlist. Each task completion requires baseline-relative
`git status` and `git diff` review to confirm only intended changes were made and
pre-existing worktree state is preserved. Cursor must run the focused race
commands assigned to implementation tasks; TASK-F17-12 runs the required final
full-repository `go test -race ./...` verification.

### Shared conclusion-validation ownership

Conclusion validation logic (outcome membership, reason-code grammar/count/
uniqueness) is implemented once in TASK-F17-06 (boundary). TASK-F17-07 (fake
constructor) reuses the same package-private helper for configured-conclusion
checking. TASK-F17-06 owns exactly one unexported
`validateConclusion(Conclusion) error` implementation in
`internal/policyeval/evaluate.go`; no exported validation API or optional helper
file is authorized. TASK-F17-06 must complete before TASK-F17-07 calls it.

### Feature-gate and completion

The founder executable-plan approval is a separate external prerequisite to
Cursor execution and must be recorded by the Feature Factory; Kiro and this
document must not fabricate or infer that receipt. After authorized Cursor
execution, this feature is not complete until:

1. All tasks pass their verification commands.
2. `make ff-feature-gate FEATURE=FEATURE-0017` passes.
3. The repository records every separately required post-implementation review
   or human-gate receipt.

No task or section of this document may self-declare feature completion or
Cursor-handoff readiness.

## 2. Dependency-ordered implementation tasks

### TASK-F17-01 — Package foundation and time-source abstraction

**Owner:** FEATURE-0017  
**Dependencies:** None (foundational; all subsequent tasks depend on this)  
**Independent from:** None (must complete first)

**Architecture traceability:**

- CDG-F17-03 (injected deterministic UTC time source pattern)
- CDG-F17-06 (in-process operation)
- Design D-01 (single package `internal/policyeval`)
- Design D-06 (local injected UTC time source)

**Requirements:** REQ-F17-03, REQ-F17-06

**Description:**

Create the new `internal/policyeval` package with package documentation, a
`TimeSource` interface abstracting deterministic time for testing, a production
UTC implementation, and a fixed test implementation. Establish the minimal
foundation without policy logic.

**Writable paths:**

- `internal/policyeval/doc.go`
- `internal/policyeval/timesource.go`
- `internal/policyeval/timesource_test.go`

**Tests:**

- Production time source returns valid UTC `time.Time`.
- Fixed test time source returns the exact configured time, unchanged by
  repeated calls.
- Fixed time normalizes supplied `time.Time` with `.UTC()`.
- Serialization to RFC3339Nano is correct for both implementations.

**Verification commands:**

```bash
go test ./internal/policyeval/... -v
go vet ./internal/policyeval/...
```

**Acceptance criteria:**

- AC-F17-16 foundation (deterministic time source exists).
- No FEATURE-0016 `Clock` import or cross-feature coupling.
- Package doc clarifies this is an engine-neutral evaluation seam, not a policy
  engine or policy system.

**Security / observability impact:**

- Time source is in-process dependency injection; no new service or external
  dependency.
- Foundation for deterministic replay in later tasks.

**Exclusions:**

- No adapter, request, result, validation, digest, or fake yet.
- No FEATURE-0016 clock import.

**Commit message:**

```
feat(policyeval): add package foundation and time-source abstraction

Introduce internal/policyeval package with TimeSource interface,
production UTC impl, and fixed test impl per REQ-F17-03/06, D-06.

No adapter, request, result, or validation yet.

Supports: TASK-F17-01, AC-F17-16 foundation
```

---

### TASK-F17-02 — Request contract and strict JSON decoder

**Owner:** FEATURE-0017  
**Dependencies:** TASK-F17-01 (package exists)  
**Independent from:** TASK-F17-04 (parallel with result types)

**Architecture traceability:**

- CDG-F17-01 (minimal request boundary)
- CDG-F17-02 (absent/empty handling)
- Design D-02 (closed Go values plus strict JSON decoder, including lexical preflight)
- Design D-04 (reuse FEATURE-0012 reference validator)

**Requirements:** REQ-F17-01

**Description:**

Define `PolicyEvaluationRequest` with exact JSON tags, create the strict byte
decoder using presence-aware `json.RawMessage` carriers. Implement a lexical
preflight that rejects invalid UTF-8 (via `utf8.Valid`) and lone, malformed, or
unpaired UTF-16 surrogate escapes (\\uD800–\\uDFFF that do not form a valid
surrogate pair) before `encoding/json` can silently replace them. After the
preflight passes, reject unknown/duplicate root and nested members using a
per-object member-name set during a `json.Decoder.Token` walk, distinguish
omitted from `null`, enforce required fields, validate all bounds, and delegate
reference structural validation to `apiref.Constraint.ValidateRef` /
`apiref.Refs.Validate`. Add `PolicyEvaluationRequest.UnmarshalJSON` delegating
to the same strict decoder.

**Writable paths:**

- `internal/policyeval/request.go`
- `internal/policyeval/request_test.go`

**Tests:**

- Valid request at every registered bound accepted (AC-F17-01).
- Invalid UTF-8 byte sequence rejected before JSON decode (AC-F17-02, D-02 preflight).
- Lone high surrogate escape (e.g., `\uD800`) rejected (AC-F17-02, D-02 preflight).
- Lone low surrogate escape (e.g., `\uDC00`) rejected (AC-F17-02, D-02 preflight).
- Malformed surrogate pair (high followed by non-surrogate) rejected (AC-F17-02, D-02 preflight).
- Unpaired surrogate in a string value rejected (AC-F17-02, D-02 preflight).
- Valid surrogate pair (e.g., `\uD83D\uDE00`) accepted (D-02 preflight).
- Unknown root field rejected (AC-F17-02).
- Unknown nested reference field rejected (AC-F17-02).
- Duplicate root member rejected (AC-F17-02).
- Duplicate nested reference member rejected (AC-F17-02).
- Trailing JSON content rejected (AC-F17-02).
- Malformed reference rejected (AC-F17-02).
- Absent `contextRef` accepted; explicit `null`, `{}`, partial rejected (AC-F17-04,
  AC-F17-05).
- Omitted `candidateRefs` constructs empty slice; `null` rejected; `[]` accepted
  (AC-F17-05).
- Empty required strings, empty `profileRefs`, `null` required fields rejected
  (AC-F17-05).
- FEATURE-0012 structural validation applied correctly to every reference.
- Direct Go construction represents semantic values (nil = absence, nil/empty
  candidates = empty set).

**Verification commands:**

```bash
go test ./internal/policyeval/... -v -run TestRequest
go vet ./internal/policyeval/...
```

**Acceptance criteria:**

- AC-F17-01 (valid request at bounds).
- AC-F17-02 (unknown/malformed rejection).
- AC-F17-04 (optional absent and valid-present `contextRef`).
- AC-F17-05 (`null`, empty, partial rules).

**Security / observability impact:**

- Strict decode prevents unknown-field injection and duplicate-member attacks.
- No request values logged or echoed into errors (AC-F17-17 foundation).
- Tenant-confidential contract; no public route added (AC-F17-22).

**Exclusions:**

- No digest calculation yet.
- No adapter invocation yet.
- No UID enforcement or duplicate-reference detection yet (TASK-F17-05).

**Commit message:**

```
feat(policyeval): add request contract and strict JSON decoder

Implement PolicyEvaluationRequest per CDG-F17-01, D-02/D-04 with
presence-aware strict decode, unknown/duplicate rejection, and
FEATURE-0012 structural reference validation.

No digest or validation beyond structural yet.

Supports: TASK-F17-02, AC-F17-01/02/04/05
```

---

### TASK-F17-03 — RFC 8785 JCS v1 canonicalization and SHA-256 digest

**Owner:** FEATURE-0017  
**Dependencies:** TASK-F17-01 (package exists), TASK-F17-02 (request type exists)  
**Independent from:** TASK-F17-04 (parallel with result types)

**Architecture traceability:**

- CDG-F17-02 (exact v1 canonicalization and digest)
- Design D-03 (stdlib-only RFC 8785 JCS v1 serializer)
- Design D-05 (normalized-reference identity for sorting)

**Requirements:** REQ-F17-02

**Description:**

Implement domain-bounded RFC 8785 JCS v1 serializer with byte builder,
`crypto/sha256`, and `encoding/hex`. Emit root keys in JCS order
(`action`, `candidateRefs`, optional `contextRef`, `profileRefs`, `schema`,
`subjectRef`). Emit reference keys in order (`apiVersion`, `kind`, `name`,
optional `uid`). Implement exact string emitter: reject invalid UTF-8 and
non-scalar Unicode; no normalization; escape quote, reverse-solidus, and
controls; emit supplementary code points and HTML-sensitive characters unchanged.
Deep-copy request into a snapshot, sort `profileRefs` and `candidateRefs` by
`(apiVersion, kind, name, uid-or-empty)` tuple, serialize, hash with SHA-256, and
encode as lowercase 64-char hex. No BOM, trailing newline, or external
prefix/suffix.

**Writable paths:**

- `internal/policyeval/canonical.go`
- `internal/policyeval/canonical_test.go`

**Tests:**

- Frozen ASCII vector reproduces exact required digest (AC-F17-10).
- Controls, quote/reverse-solidus, HTML-sensitive, BMP, supplementary,
  U+2028/U+2029, invalid UTF-8, lone surrogates, composed vs decomposed strings.
- Reordered `profileRefs`/`candidateRefs` produce same digest (AC-F17-07).
- Semantic change alters digest (AC-F17-08).
- `requestId`-only change preserves digest (AC-F17-09).
- Omitted vs `[]` `candidateRefs` both serialize to `"candidateRefs":[]`
  (AC-F17-05).
- Omitted `contextRef` is absent from normalized object; present included
  (AC-F17-04, AC-F17-05).
- No whitespace, BOM, trailing newline, or external prefix/suffix.

**Verification commands:**

```bash
go test ./internal/policyeval/... -v -run TestCanonical
go vet ./internal/policyeval/...
```

**Acceptance criteria:**

- AC-F17-07 (list-order-independent digest).
- AC-F17-08 (semantic change alters digest).
- AC-F17-09 (`requestId`-only change preserves digest).
- AC-F17-10 (exact RFC 8785/SHA-256 test vector).

**Security / observability impact:**

- Deterministic digest enables replay testing and fixture lookup.
- No request values logged; digest is internal identity only (AC-F17-17
  foundation).

**Exclusions:**

- No duplicate-reference rejection yet (TASK-F17-05).
- No adapter invocation yet.

**Commit message:**

```
feat(policyeval): add RFC 8785 JCS v1 canonicalization and SHA-256 digest

Implement stdlib-only domain-bounded JCS serializer with exact
string emitter, sorting, and lowercase hex SHA-256 digest per
CDG-F17-02, D-03/D-05.

Frozen vector and comprehensive string test suite included.

Supports: TASK-F17-03, AC-F17-07/08/09/10
```

---

### TASK-F17-04 — Result types, adapter port, and non-result categories

**Owner:** FEATURE-0017  
**Dependencies:** TASK-F17-01 (package exists)  
**Independent from:** TASK-F17-02, TASK-F17-03 (parallel with request and digest)

**Architecture traceability:**

- CDG-F17-03 (adapter, result, timing, failure normalization)
- Design D-08 (normalized non-result category is internal typed value)
- Design D-09 (constructed boundary with conclusion-or-failure adapter port)
- Design D-10 (boundary owns conclusion validation)

**Requirements:** REQ-F17-03

**Description:**

Define `Outcome` enum (`Allow`, `Deny`, `Indeterminate`, `RequiresApproval`),
`PolicyEvaluationResult` with exact JSON tags and absence of `obligations`,
`NonResult` enum (five categories), `NormalizedInput`, `Conclusion`,
`PolicyEngineAdapter` interface, and reason-code grammar constants. Document that
Phase 2R fake produces no obligations and that adapter error is normalized to
`AdapterFailure`.

**Writable paths:**

- `internal/policyeval/result.go`
- `internal/policyeval/adapter.go`
- `internal/policyeval/result_test.go`

**Tests:**

- `Outcome` enum contains exactly four valid values (AC-F17-11 foundation).
- `NonResult` enum contains exactly five categories (AC-F17-21 foundation).
- Reason-code grammar regex `^[A-Z][A-Z0-9_]{0,62}$` is correct.
- `obligations` field serializes as omitted when nil/empty (AC-F17-18).

**Verification commands:**

```bash
go test ./internal/policyeval/... -v -run TestResult
go vet ./internal/policyeval/...
```

**Acceptance criteria:**

- AC-F17-11 foundation (four outcomes defined).
- AC-F17-18 (obligations absent in Phase 2R fake).
- AC-F17-21 foundation (non-result categories defined).

**Security / observability impact:**

- Internal types; no route or public projection (AC-F17-22).
- Reason-code grammar constrains shape; content redaction deferred to fake
  configuration and fixture review (AC-F17-17 preparation).

**Exclusions:**

- No boundary or fake yet.
- No actual validation or invocation yet.

**Commit message:**

```
feat(policyeval): add result types, adapter port, and non-result categories

Define Outcome, PolicyEvaluationResult, NonResult, NormalizedInput,
Conclusion, PolicyEngineAdapter per CDG-F17-03, D-08/09/10.

No obligations in Phase 2R; adapter error is AdapterFailure signal.

Supports: TASK-F17-04, AC-F17-11/18/21 foundation
```

---

### TASK-F17-05 — Request validation including UID and duplicate enforcement

**Owner:** FEATURE-0017  
**Dependencies:** TASK-F17-01 (package exists), TASK-F17-02 (request exists), TASK-F17-03 (normalization exists)  
**Independent from:** TASK-F17-04 (parallel with result types)

**Architecture traceability:**

- CDG-F17-01 (minimal request boundary, duplicate rejection)
- CDG-F17-02 (normalized-reference identity)
- Design D-05 (normalized-reference identity and duplicate detection)

**Requirements:** REQ-F17-01

**Description:**

Add request validation enforcing `contextRef.kind == EffectiveGovernanceContext`
plus required UID for `contextRef` and every `profileRefs` entry. Detect
duplicate `profileRefs` and duplicate `candidateRefs` using normalized-reference
equality `(apiVersion, kind, name, uid-or-empty)`. Reject duplicates, never
deduplicate silently. Validation occurs after structural validation and before
digest calculation.

**Writable paths:**

- `internal/policyeval/request.go` (add validation functions)
- `internal/policyeval/request_test.go` (add validation tests)

**Tests:**

- Duplicate `profileRefs` rejected (AC-F17-06).
- Duplicate `candidateRefs` rejected (AC-F17-06).
- `contextRef` with wrong kind rejected.
- `contextRef` missing UID rejected.
- `profileRefs` entry missing UID rejected.
- No false-positive duplicate when UIDs differ.

**Verification commands:**

```bash
go test ./internal/policyeval/... -v -run TestRequestValidation
go vet ./internal/policyeval/...
```

**Acceptance criteria:**

- AC-F17-06 (duplicate reference rejection).
- AC-F17-04 extension (UID-pinned `contextRef<EffectiveGovernanceContext>`).

**Security / observability impact:**

- Duplicate rejection prevents semantic ambiguity and potential confusion
  attacks.
- No leakage of request values into errors (AC-F17-17 foundation).

**Exclusions:**

- No existence lookup or domain validation.
- No adapter invocation yet.

**Commit message:**

```
feat(policyeval): add UID and duplicate-reference validation

Enforce contextRef.kind, required UIDs, and duplicate profileRefs/
candidateRefs rejection using normalized-reference identity per
CDG-F17-01/02, D-05.

Supports: TASK-F17-05, AC-F17-06
```

---

### TASK-F17-06 — Evaluation boundary with exact precedence and at-most-once invocation

**Owner:** FEATURE-0017  
**Dependencies:** TASK-F17-01 (time source), TASK-F17-02 (request), TASK-F17-03 (digest), TASK-F17-04 (result/adapter), TASK-F17-05 (validation)  
**Independent from:** None (on the critical path; TASK-F17-07 depends on its conclusion-validation helper)

**Architecture traceability:**

- CDG-F17-03 (boundary owns validation, timing, conclusion validation, result
  construction)
- CDG-F17-06 (exact precedence, at-most-once invocation, no retry)
- Design D-09 (constructed boundary with adapter and time source)
- Design D-10 (boundary owns conclusion validation, defensive ownership, digest,
  time)
- Design D-11 (direct-call linearization with cancellation-first return
  precedence)

**Requirements:** REQ-F17-03, REQ-F17-06

**Description:**

Implement `NewBoundary(adapter, ts)` with nil/typed-nil dependency rejection.
Implement `Evaluate(ctx, req)` performing the exact nine-step precedence: (1)
pre-check `ctx.Err()`; (2) request validation; (3) digest; (4) pre-invocation
`ctx.Err()`; (5) capture `startedAt`, direct-call adapter on current goroutine;
(6) post-return `ctx.Err()` before inspecting output, discard late output; (7)
non-nil adapter error → `AdapterFailure`; (8) conclusion validation (outcome
membership, reason-code grammar/count/uniqueness); (9) capture `completedAt`,
sort reason codes, set `evaluatedAt = completedAt`, construct result, return
result and timing. On any normalized non-result, return zero result, zero timing,
category, sanitized error. Return unnamed transient pairing of result and
`decision.TimingEnvelope`.

Additionally, implement exactly one package-private
`validateConclusion(Conclusion) error` helper in `evaluate.go`. It checks
outcome membership, reason-code grammar (`^[A-Z][A-Z0-9_]{0,62}$`), count
(1..32), and uniqueness. It is the single source of truth reused by both the
boundary (step 8) and TASK-F17-07's fake constructor. No exported validation
API or additional production file is permitted.

**Writable paths:**

- `internal/policyeval/evaluate.go`
- `internal/policyeval/evaluate_test.go`

**Tests:**

- Nil context → `RequestInvalid`, zero result, zero timing, zero invocations
  (AC-F17-03).
- Pre-call cancel/deadline → zero invocations, zero result, zero timing,
  `Canceled`/`DeadlineExceeded` (AC-F17-03, AC-F17-14).
- Request validation failure → `RequestInvalid`, zero result, zero timing, zero
  invocations, sanitized error (AC-F17-03, AC-F17-14).
- Pre-invocation cancel/deadline → zero result, zero timing, zero invocations
  (AC-F17-14).
- Exactly one adapter call after all pre-invocation checks pass (AC-F17-14).
- Post-return cancel/deadline discards late output, returns zero result, zero
  timing, `Canceled`/`DeadlineExceeded` (AC-F17-14).
- Adapter returning both a non-nil error and a malformed conclusion →
  `AdapterFailure`, zero result, zero timing, and sanitized error; this
  observably proves adapter-failure precedence without inspecting private helper
  call counts (AC-F17-12 partial).
- Invalid outcome → `InvalidAdapterResult`, zero result, zero timing, sanitized
  error (AC-F17-15).
- Invalid reason codes (grammar, duplicate, count) → `InvalidAdapterResult`,
  zero result, zero timing, sanitized error (AC-F17-15, AC-F17-17 partial).
- Valid conclusion → correct result, timing, `evaluatedAt == completedAt`
  (AC-F17-11 partial, AC-F17-16 partial).
- Spy adapter proves at-most-once invocation and no retry (AC-F17-14).
- Nil boundary dependency rejected at construction.
- Typed-nil boundary dependency rejected at construction.
- Every NonResult category returns exactly zero `PolicyEvaluationResult` (all
  fields at zero value) and zero `TimingEnvelope`.

**Verification commands:**

```bash
go test ./internal/policyeval/... -v -run TestEvaluate
go test ./internal/policyeval/... -race -run TestEvaluate
go vet ./internal/policyeval/...
```

**Acceptance criteria:**

- AC-F17-03 (invalid request invokes no adapter).
- AC-F17-11 partial (valid outcomes producible).
- AC-F17-14 (cancellation, deadline, at-most-once, late-output discard).
- AC-F17-15 (invalid adapter conclusion rejected).
- AC-F17-16 partial (fixed time source deterministic replay foundation).

**Security / observability impact:**

- Exact precedence prevents bypassing cancellation, validation, or conclusion
  checks.
- Sanitized diagnostics never wrap adapter errors or include request/result
  values (AC-F17-17 foundation).
- No goroutine or channel creation; direct-call cardinality follows D-11 and
  does not reinterpret AC-F17-23.

**Exclusions:**

- No fake yet; use spy adapter for invocation tests.
- No FEATURE-0013 mapper yet.

**Commit message:**

```
feat(policyeval): add evaluation boundary with exact precedence

Implement NewBoundary and Evaluate with nine-step precedence,
at-most-once direct-call invocation, cancellation-first tie,
conclusion validation, and result construction per CDG-F17-03/06,
D-09/10/11.

Returns unnamed transient result+timing pairing.

Supports: TASK-F17-06, AC-F17-03/11/14/15/16 partial
```

---

### TASK-F17-07 — Deterministic fake adapter

**Owner:** FEATURE-0017  
**Dependencies:** TASK-F17-04 (adapter port exists), TASK-F17-06 (shared conclusion-validation helper exists)  
**Independent from:** TASK-F17-09 only (mapper may proceed in parallel;
TASK-F17-08 remains downstream of this task)

**Architecture traceability:**

- CDG-F17-04 (deterministic fake)
- Design D-12 (immutable digest-keyed fake with explicit copy ownership)

**Requirements:** REQ-F17-04

**Description:**

Define `FakeFixture` with digest, optional `*Conclusion`, and
`AdapterFailure bool`. Implement `NewDeterministicFake([]FakeFixture)` that
deep-copies input and nested slices, validates digest grammar, detects duplicate
digests, enforces exactly one of conclusion/failure, validates configured
conclusions using the shared package-private `validateConclusion` helper from
TASK-F17-06, and
fails construction on invalid configuration. Implement `Evaluate` returning
fresh-copied configured conclusion or `AdapterFailure` for missing/configured-
failure digests. No time, no conditionals interpreting domain meaning.

**Writable paths:**

- `internal/policyeval/fake.go`
- `internal/policyeval/fake_test.go`

**Tests:**

- Valid digest mapped to valid conclusion returns that outcome/reason codes
  (AC-F17-11).
- Missing/unconfigured digest → `AdapterFailure` (AC-F17-12).
- Configured `AdapterFailure` → adapter error signal (AC-F17-12).
- Duplicate digest fails construction (AC-F17-13).
- Malformed digest fails construction (AC-F17-13).
- Invalid conclusion (bad one-of, malformed outcome/reason codes) fails
  construction (AC-F17-13).
- Constructor and returned conclusion deep-copy ownership (D-12).
- Fixture mutation after construction or returned conclusion mutation cannot
  affect fake or future returns (D-12).

**Verification commands:**

```bash
go test ./internal/policyeval/... -v -run TestFake
go vet ./internal/policyeval/...
```

**Acceptance criteria:**

- AC-F17-11 (all four valid outcomes).
- AC-F17-12 (missing fixture and configured adapter failure).
- AC-F17-13 (invalid fixture construction failure).
- D-12 conformance proof (copy ownership; no reinterpretation of AC-F17-23).

**Security / observability impact:**

- Fake is deterministic and policy-intelligence-free conformance machinery.
- Fixture strings are trusted, pre-redacted, repository-controlled configuration
  (AC-F17-17 preparation).

**Exclusions:**

- No YAML decoder yet (TASK-F17-08).
- No domain conditionals.

**Commit message:**

```
feat(policyeval): add deterministic fake adapter

Implement FakeFixture and NewDeterministicFake with duplicate-
preserving constructor, immutable digest-keyed lookup, copy
ownership, and construction validation per CDG-F17-04, D-12.

No domain conditionals; deterministic behavior only.

Supports: TASK-F17-07, AC-F17-11/12/13, D-12 copy ownership
```

---

### TASK-F17-08 — Test-only strict YAML fixture decoder

**Owner:** FEATURE-0017  
**Dependencies:** TASK-F17-07 (fake exists)  
**Independent from:** TASK-F17-09 (parallel with mapper)

**Architecture traceability:**

- Design D-13 (fixture serialization: in-memory Go config, optional strict YAML
  for tests)

**Requirements:** REQ-F17-04

**Description:**

Add test-only helper `decodeFakeFixturesYAMLForTest([]byte)` in
`fake_fixture_test.go` using `gopkg.in/yaml.v3` with `KnownFields(true)`,
duplicate-key rejection, exactly one document, and trailing-content rejection.
Return duplicate-preserving slice; all semantic validation remains in
`NewDeterministicFake`.

**Writable paths:**

- `internal/policyeval/fake_fixture_test.go` (add YAML helper and its tests)

**Tests:**

- Valid YAML fixture decoded correctly.
- Unknown YAML field rejected.
- Duplicate YAML key rejected.
- Multiple YAML documents rejected.
- Trailing YAML content rejected.
- Duplicate fixture slice preserved; constructor detects and rejects.

**Verification commands:**

```bash
go test ./internal/policyeval/... -v -run TestFakeFixtureYAML
go vet ./internal/policyeval/...
```

**Acceptance criteria:**

- D-13 conformance proof (strict YAML for test fixtures); this does not extend
  or reinterpret AC-F17-13.
- Test-only helper; no production YAML input surface.

**Security / observability impact:**

- Test-only helper; no runtime YAML parsing or reload surface.
- Fixture content is trusted and pre-redacted (AC-F17-17 preparation).

**Exclusions:**

- No runtime YAML or production fixture reload.

**Commit message:**

```
feat(policyeval): add test-only strict YAML fixture decoder

Add decodeFakeFixturesYAMLForTest with KnownFields, duplicate-key,
one-document, trailing-content rejection per D-13.

Test-only; no production YAML surface.

Supports: TASK-F17-08, D-13
```

---

### TASK-F17-09 — FEATURE-0013 structural-validation export and pure mapper

**Owner:** FEATURE-0017  
**Dependencies:** TASK-F17-04 (result exists), TASK-F17-06 (boundary returns timing)  
**Independent from:** TASK-F17-07, TASK-F17-08 (parallel with fake)

**Architecture traceability:**

- CDG-F17-05 (transient FEATURE-0013 evidence linkage)
- Design D-14 (pure mapper with exact FEATURE-0012/0013 structural-contract
  reuse)
- Approved dependency extension (Sanjeev Kumar, 2026-08-25)

**Requirements:** REQ-F17-05

**Description:**

Add approved FEATURE-0013 structural-validation bool wrapper:
`internal/decision/validate/evaluation_structure_export.go` containing
`IsEvaluationResultStructurallyValid(result)` delegating to existing private
`checkEvaluationResultStructure`, plus parity tests in
`evaluation_structure_export_test.go`.

Add `internal/policyeval/mapper.go` defining `InputIdentity` with snapshot/
integrity validation (at least one required, both permitted), implementing
`MapToEvaluationResult(result, timing, evaluator, ident)` with closed status
mapping (`Allow`/`Deny`/`RequiresApproval` → `SUCCESS`, `Indeterminate` →
`INDETERMINATE`), exact populated fields (`Evaluator`, `InputSnapshotRef`,
`InputIntegrity`, `ResultStatus`, marshaled `Result`, `Timing.StartedAt`,
`Timing.CompletedAt`, `EvaluatedAt`), exact omitted fields (`Timing.DurationMs`,
`ExecutionConfig`, `TrustBoundary`, `SafetyPolicyFilters`, `StructuralTrustState`,
`ScopeEvidence`), and equality assertion `result.evaluatedAt ==
timing.completedAt`. Delegate structural validation to
`decisionvalidate.IsEvaluationResultStructurallyValid`. Perform no time read,
lookup, profile/scope validation, write, persistence, or external effect.

**Writable paths:**

- `internal/decision/validate/evaluation_structure_export.go`
- `internal/decision/validate/evaluation_structure_export_test.go`
- `internal/policyeval/mapper.go`
- `internal/policyeval/mapper_test.go`

**Tests:**

- FEATURE-0013 wrapper returns exact existing structural-pass result; no
  behavior change.
- All four outcomes map correctly to `SUCCESS`/`INDETERMINATE` (AC-F17-19).
- Exact populated fields present, exact omitted fields absent (AC-F17-19).
- `evaluatedAt == completedAt` assertion (AC-F17-19).
- Malformed evaluator identity (empty/blank fields) fails mapping (AC-F17-20).
- Invalid timing (zero `StartedAt`, zero `CompletedAt`, `evaluatedAt !=
  completedAt` mismatch) fails mapping (AC-F17-20).
- Invalid snapshot reference (malformed TypedRef) fails mapping (AC-F17-20).
- Invalid integrity carrier (nil `State`, empty `State`, partial carrier) fails
  mapping (AC-F17-20).
- Both `InputSnapshotRef` and `InputIntegrity` nil fails mapping (AC-F17-20).
- Invalid policy-result outcome (not in closed set) fails mapping (AC-F17-20).
- Invalid policy-result reason codes (grammar violation, duplicate, >32 count,
  empty) fail mapping (AC-F17-20).
- Invalid policy-result digest (wrong length, non-hex) fails mapping (AC-F17-20).
- Timestamp validation is pinned to the boundary's canonical UTC
  `time.RFC3339Nano` representation: accept a canonical `Z` value; reject empty,
  malformed, non-UTC-offset, offset-zero `+00:00`, and otherwise non-canonical
  timestamp strings. Reject zero or non-UTC `StartedAt`/`CompletedAt` values
  (AC-F17-20).
- `evaluatedAt`/`completedAt` mismatch fails mapping (AC-F17-20).
- Every non-empty obligations slice is rejected in Phase 2R, including values
  containing valid JSON scalars, objects, or arrays as well as malformed raw
  JSON; nil and empty represent absence (AC-F17-20, D-14 §5.5).
- Malformed structural FEATURE-0013 candidate fails via bool wrapper (AC-F17-20).
- Valid policy result unchanged by mapping failure — no mutation (AC-F17-20).
- Mapper copies all permitted mutable inputs and output raw JSON; mutation of
  caller-owned values cannot change mapped output (D-14). Obligations cannot be
  copied into Phase 2R output because every non-empty value is rejected.

**Verification commands:**

```bash
go test ./internal/decision/validate/... -v -run TestEvaluationStructureExport
go test ./internal/policyeval/... -v -run TestMapper
go vet ./internal/decision/validate/...
go vet ./internal/policyeval/...
```

**Acceptance criteria:**

- AC-F17-19 (successful result/timing transport and pure FEATURE-0013 mapping).
- AC-F17-20 (malformed structural mapper-input failure with no mutation).
- AC-F17-21 (non-result interactions produce no mapping; tested in TASK-F17-11).
- D-14 conformance proof (copy ownership; no reinterpretation of AC-F17-23).

**Security / observability impact:**

- Approved structural-only extension; no Problem exposure or new FEATURE-0013
  rule.
- Pure mapping; no mutation, write, persistence, or external effect (AC-F17-22).

**Exclusions:**

- No DecisionProfile, DecisionRecord scope, adopting validation, publication,
  or AuditEvent.

**Commit message:**

```
feat(policyeval): add approved FEATURE-0013 bool wrapper and pure mapper

Add decision/validate bool wrapper delegating to existing private
structural check, plus MapToEvaluationResult with closed status
mapping, exact populated/omitted fields, timing equality per
CDG-F17-05, D-14.

No DecisionProfile/scope/adoption/publication.

Supports: TASK-F17-09, AC-F17-19/20/21, D-14 copy ownership
```

---

### TASK-F17-10 — Complete unit and determinism tests

**Owner:** FEATURE-0017  
**Dependencies:** All implementation tasks (TASK-F17-01..09)  
**Independent from:** None (TASK-F17-11 depends on this task)

**Architecture traceability:**

- CDG-F17-06 (determinism and concurrency)
- Design §5.3 (determinism and concurrency), §5.5 (mapper validation)

**Requirements:** REQ-F17-06

**Description:**

Add comprehensive unit tests covering bounds, precedence, digest properties,
fake behavior, mapper validation, and concurrent determinism. Prove fixed time
source produces deterministic `evaluatedAt` and results. Prove concurrent
equivalent-call determinism under `-race`. Prove copy ownership by mutating
caller/fixture/adapter-return/output pointers and slices after handoff.

**Writable paths:**

- `internal/policyeval/determinism_test.go`

**Tests:**

- Fixed time source → deterministic `evaluatedAt` and result replay (AC-F17-16).
- Concurrent equivalent calls produce identical results (AC-F17-23).
- Copy ownership: caller/fixture/adapter/output mutations after handoff do not
  affect boundary/fake/result (D-10, D-12, D-14 support for the concurrent
  determinism proof; copy ownership does not redefine AC-F17-23).
- Reason-code sorting is stable and total (AC-F17-17 partial).
- All required bounds checked (AC-F17-01).

**Verification commands:**

```bash
go test ./internal/policyeval/... -v
go test ./internal/policyeval/... -race
go vet ./internal/policyeval/...
```

**Acceptance criteria:**

- AC-F17-01 (complete bounds coverage).
- AC-F17-16 (fixed-time-source deterministic replay).
- AC-F17-23 (concurrent equivalent-call determinism).
- D-10/D-12/D-14 conformance proof (copy ownership).

**Security / observability impact:**

- Determinism proof enables fixture-based conformance and replay.
- Race detector coverage confirms concurrency safety.

**Exclusions:**

- No integration or cross-feature mutation tests yet (TASK-F17-11).

**Commit message:**

```
test(policyeval): add complete unit and determinism tests

Prove fixed-time deterministic replay, concurrent equivalent-call
determinism under -race, copy ownership, and complete bounds
coverage per CDG-F17-06, §5.3.

Supports: TASK-F17-10, AC-F17-01/16/23
```

---

### TASK-F17-11 — Integration, conformance, and zero-external-effect tests

**Owner:** FEATURE-0017  
**Dependencies:** All implementation and unit-test tasks (TASK-F17-01..10)  
**Independent from:** None (final test task before closeout)

**Architecture traceability:**

- CDG-F17-06 (in-process operation, no route/store/controller, no external call)
- Design §6 (security, privacy, observability, compatibility), §7 (test and
  conformance strategy)

**Requirements:** REQ-F17-06

**Description:**

Add integration tests combining boundary, fake, and mapper. Prove non-result
interactions produce no mapping. Add package-level conformance assertions that
are non-vacuous:

- **Closed import-surface test:** parse every production `.go` file in
  `internal/policyeval` (never `_test.go`) and fail on every import not present
  in the exact D-01/D-03/D-14 allowlist: Sovrunn
  `github.com/sanjeevksaini/sovrunn/internal/apimeta`,
  `github.com/sanjeevksaini/sovrunn/internal/apiref`,
  `github.com/sanjeevksaini/sovrunn/internal/decision`, and
  `github.com/sanjeevksaini/sovrunn/internal/decision/validate`; and stdlib
  `bytes`, `context`, `crypto/sha256`, `encoding/hex`, `encoding/json`, `errors`,
  `io`, `reflect`, `sort`, `strings`, `time`, and `unicode/utf8`. The test
  compares parsed import paths, not source keywords, so its own strings cannot
  satisfy or defeat the assertion. Test-only imports, including
  `gopkg.in/yaml.v3`, are outside the production allowlist check.
- **Symbol-absence test:** inspect production `.go` file declarations (not
  `_test.go`) and confirm no exported or unexported symbol references
  DecisionRecord, AuditEvent, route, controller, or store types. Use `go/ast`
  or file-enumeration; do not use grep-in-test-string patterns that would
  vacuously pass.
- **Fixture-corpus review:** no committed fixture artifact is claimed. Parse the
  complete actual test corpus: every `FakeFixture`/`Conclusion` composite
  literal and every YAML literal passed to `decodeFakeFixturesYAMLForTest` in
  `_test.go` files. Fail if a fixture source cannot be classified, and inspect
  every discovered reason code for CloudProvider-native or customer-sensitive
  tokens (`aws`, `azure`, `gcp`, `customer`, `secret`, `credential`, `password`,
  `token`). A future committed fixture artifact would require a separately
  enumerated corpus source; this task does not claim one exists now.

Add leakage tests scanning sentinel inputs to prove request values, reason codes,
raw adapter errors, and strict-decoder member names do not appear in returned
errors or diagnostics. Decoder coverage must include unique sentinel unknown
member names and duplicate member names at both root and nested-reference scope.

**Writable paths:**

- `internal/policyeval/integration_test.go`
- `internal/policyeval/conformance_test.go`

**Tests:**

- End-to-end: request → validation → digest → fake → result → timing → mapper
  → `EvaluationResult` for all four outcomes (AC-F17-11, AC-F17-19).
- Non-result interactions (`RequestInvalid`, `Canceled`, `DeadlineExceeded`,
  `AdapterFailure`, `InvalidAdapterResult`) each produce: zero result, zero
  timing, the specific NonResult category, a sanitized error, and zero adapter
  invocation count for categories before invocation (AC-F17-21).
- No `internal/policyeval` production source file declares or references
  DecisionRecord, AuditEvent, route handler, controller, or store types
  (AC-F17-22) — verified via `go/ast` or production-file enumeration, not
  self-matching test strings.
- Every parsed direct production import belongs to the exact closed allowlist
  above; any undeclared stdlib, Sovrunn, or third-party import fails the test
  (AC-F17-24). Together with the closed in-process API and production-symbol
  inspection, this is the non-vacuous zero-other-external-effects proof.
- Sentinel leakage tests: inject unique sentinel strings as request field values,
  adapter error messages, reason-code text, fixture digest, unknown JSON member
  names, and duplicate JSON member names at root and nested-reference scope;
  assert none appear in returned `error.Error()` strings or NonResult
  diagnostics (AC-F17-17).
- Fixture-corpus review: AST-enumerate the complete actual Go/YAML fixture corpus
  described above, fail closed on unclassified fixture sources, and assert no
  CloudProvider-native or customer-sensitive token matches the curated deny-list
  (AC-F17-17).
- Obligations field absent in all Phase 2R fake returns (AC-F17-18).

**Verification commands:**

```bash
go test ./internal/policyeval/... -v -run TestIntegration
go test ./internal/policyeval/... -v -run TestConformance
go test ./internal/policyeval/... -v -run TestLeakage
go vet ./internal/policyeval/...
```

**Acceptance criteria:**

- AC-F17-11 (all four valid outcomes end-to-end).
- AC-F17-17 (reason-code grammar, sorting, and redaction).
- AC-F17-18 (absence of obligations in Phase 2R fake).
- AC-F17-21 (no mapping for non-result interactions).
- AC-F17-22 (no FEATURE-0017 DecisionRecord/AuditEvent/route/controller/store).
- AC-F17-24 (zero external effects).

**Security / observability impact:**

- Import and symbol surface tests prevent accidental external dependencies or
  effects.
- Leakage tests prevent sensitive value disclosure in errors/diagnostics.
- Trusted fixture review proves Phase 2R conformance codes are internal-only.

**Exclusions:**

- No downstream metric, side effect, controller, adapter/plugin execution, or
  `VS0-CF-*` runtime proof (FEATURE-0017 acceptance is `AC-F17-01..24` only).

**Commit message:**

```
test(policyeval): add integration, conformance, and zero-external-effect tests

Prove end-to-end evaluation + mapping for all four outcomes,
no mapping for non-result, no FEATURE-0017 DecisionRecord/route/
store symbols, zero external imports/effects, and sanitized
diagnostics per CDG-F17-06, §6, §7.

Supports: TASK-F17-11, AC-F17-11/17/18/21/22/24
```

---

### TASK-F17-12 — Read-only repository verification and closeout evidence

**Owner:** FEATURE-0017  
**Dependencies:** All implementation and test tasks (TASK-F17-01..11)  
**Independent from:** None (final task)

**Architecture traceability:**

- Feature-gate requirements
- Baseline-relative worktree and artifact-review requirements

**Requirements:** All REQ-F17-01..06, all AC-F17-01..24

**Description:**

Perform read-only closeout verification after TASK-F17-01..11 have completed.
Run a read-only repository-wide `gofmt -l` cleanliness assertion, followed by
the mandatory `make fmt` verification, `make test`, `make vet`, focused package
checks, and the required full-repository `go test -race ./...` so the
FEATURE-0013 validation extension is also race-covered. Because the preceding
assertion proves formatting cleanliness, `make fmt` must be a no-op. Record the
baseline-relative diff before and after it; if it changes any file, TASK-F17-12
stops, reopens the owning implementation task, and restarts only after that task
corrects formatting. Validate that the owned schemas (`VS0-SCHEMA-018` and
`VS0-SCHEMA-019`) remain documented and that no public route, store, controller,
production engine, generated artifact, temporary file, or stale fixture was
added.

Compare `git status` and `git diff` with the recorded pre-task worktree baseline.
The gate does not require an initially dirty repository to become clean; it
requires every FEATURE-0017 delta to be intended and every pre-existing change
to remain preserved. Confirm all REQ and AC mappings in the coverage ledgers.
TASK-F17-12 makes no edits: any failed check reopens the earliest owning task and
its exact writable-path allowlist, after which TASK-F17-12 restarts.

Run `make ff-feature-gate FEATURE=FEATURE-0017` only after implementation tests
and the required review receipts are present. The separately recorded founder
executable-plan approval is a prerequisite to starting Cursor, not an approval
that Kiro or this task may assert or create.

**Writable paths:**

- None. This task is strictly read-only; failures reopen the owning task.

**Tests:**

- All tests pass: `go test ./internal/policyeval/... -v`
- Race detector clean: `go test ./internal/policyeval/... -race`
- Full-repository race detector clean: `go test -race ./...`
- Vet clean: `go vet ./internal/policyeval/...`
- Repository formatting assertion reports no files; mandatory `make fmt` passes
  as a recorded no-op with no baseline-relative diff; `make test` and `make vet`
  pass
- Feature gate: `make ff-feature-gate FEATURE=FEATURE-0017` passes

**Verification commands:**

```bash
test -z "$(gofmt -l .)"
make fmt
make test
make vet
go test -race ./internal/policyeval/...
go test -race ./...
git status
git diff
make ff-feature-gate FEATURE=FEATURE-0017
```

**Acceptance criteria:**

- All AC-F17-01..24 mapped and traceable.
- All REQ-F17-01..06 covered.
- No generated artifacts committed.
- Baseline-relative `git status`/`git diff` show only intended FEATURE-0017
  deltas while preserving every pre-existing change.
- The read-only formatting assertion passes and the required `make fmt` command
  is recorded as a no-op; any formatting delta reopens its owning task.
- All repository verification commands pass.
- Focused and full-repository `go test -race` pass for this
  concurrency-sensitive feature and its FEATURE-0013 validation extension.
- `make ff-feature-gate FEATURE=FEATURE-0017` passes.

**Security / observability impact:**

- Final verification confirms no leakage, external effects, or accidental
  routes/stores.
- Baseline-relative diff review confirms no sensitive fixture values,
  credentials, or test data were introduced.
- Race detector confirms concurrency safety.

**Exclusions:**

- No downstream integration, `VS0-CF-*` runtime proof, or cross-feature mutation
  tests (those are FEATURE-0026 integration responsibility).
- No cleanup or refactoring outside FEATURE-0017 owned paths.

**Commit message:** None. TASK-F17-12 is verification-only and creates no
commit. Any corrective commit belongs to the reopened owning task.

---

## 3. No-task ledger

The following approved design elements are contract-only, governance, invariants,
or exclusions with no implementation artifact in this feature. They are
acknowledged explicitly and are not assigned to a task.

| Design element | Classification | Rationale |
|---|---|---|
| No public routes | CONTRACT_ONLY | FEATURE-0017 defines no HTTP endpoint; internal operation only (AC-F17-22). |
| No controller | CONTRACT_ONLY | In-process operation; no lifecycle controller (AC-F17-22). |
| No persistence / store | CONTRACT_ONLY | Transient evidence; no FEATURE-0017 result repository (AC-F17-22, CDG-F17-05). |
| No production adapter selection | EXCLUDED | Deferred to later approved production-engine feature (CDG-F17-06, §8). |
| No real OPA / Cedar / engine | EXCLUDED | Phase 2R uses deterministic fake only (CDG-F17-04, §6, AC-F17-24). |
| No DecisionRecord / AuditEvent publication | EXCLUDED | FEATURE-0013-owned; FEATURE-0017 supplies transient evidence only (CDG-F17-05, AC-F17-22). |
| No authentication / authorization | EXCLUDED | FEATURE-0018-owned; FEATURE-0017 is trusted in-process operation (§8, §10). |
| No EffectiveGovernanceContext resolution | EXCLUDED | FEATURE-0020-owned; FEATURE-0017 accepts optional structural reference only (CDG-F17-01, §8). |
| No governance / sovereignty / placement semantics | EXCLUDED | Later owning domains consume FEATURE-0017 results (§8, §10). |
| No entitlement / quota / profile interpretation | EXCLUDED | Later owning domains (FEATURE-0021/0022/0023); fake has no conditionals (CDG-F17-04). |
| No customer-safe / AI-readable explanation | EXCLUDED | FEATURE-0025-owned (§8, §10). |
| No obligations in Phase 2R fake | INVARIANT | Typed obligation semantics deferred (CDG-F17-03, AC-F17-18). |
| No time-dependent policy semantics | INVARIANT | `evaluatedAt` is chronology only; production time policy deferred (CDG-F17-03). |
| No v2 digest or profile-kind extension | GOVERNANCE | Requires approved v2 digest contract and fresh reuse assessment (CDG-F17-02, §6). |
| No modification to FEATURE-0012 / 0013 / 0015 / 0016 ownership | GOVERNANCE | Reuse by reference only; no fork or redefine (§8, D-04, D-14). |
| Frozen RFC 8785 / SHA-256 test vector | INVARIANT | v1 digest identity is pinned (CDG-F17-02, AC-F17-10). |
| Trusted fixture content and pre-redaction | GOVERNANCE | Fixture strings are repository-controlled, pre-redacted configuration (CDG-F17-04, AC-F17-17). |

## 4. Canonical coverage ledger

Every approved REQ and AC appears exactly once below, mapped to its implementing
task or to the no-task ledger. IDs preserve their approved meaning and are
neither created, renumbered, merged, split, reinterpreted, nor omitted.

### 4.1 Requirements coverage

| REQ | Implementing task(s) | Design realization |
|---|---|---|
| REQ-F17-01 | TASK-F17-02 (request), TASK-F17-05 (UID/duplicate) | D-02, D-04, D-05; `request.go` validation |
| REQ-F17-02 | TASK-F17-03 (canonicalization and digest) | D-03, D-05; `canonical.go` |
| REQ-F17-03 | TASK-F17-01 (time source), TASK-F17-04 (result/adapter), TASK-F17-06 (boundary) | D-06, D-08, D-09, D-10, D-11; `timesource.go`, `result.go`, `adapter.go`, `evaluate.go` |
| REQ-F17-04 | TASK-F17-07 (fake), TASK-F17-08 (YAML helper) | D-12, D-13; `fake.go` |
| REQ-F17-05 | TASK-F17-09 (FEATURE-0013 wrapper and mapper) | D-14; `mapper.go`, approved bool wrapper |
| REQ-F17-06 | TASK-F17-01 (foundation), TASK-F17-06 (boundary precedence), TASK-F17-10 (determinism), TASK-F17-11 (zero-external-effect) | D-01, D-09, D-11, §5.1, §5.3, §6; `evaluate.go` |

### 4.2 Acceptance coverage

| AC | Implementing task(s) | Notes |
|---|---|---|
| AC-F17-01 | TASK-F17-02, TASK-F17-10 | Bounds table test |
| AC-F17-02 | TASK-F17-02 | Strict decode, unknown/duplicate/trailing rejection |
| AC-F17-03 | TASK-F17-06 | Invalid request → no adapter invocation |
| AC-F17-04 | TASK-F17-02, TASK-F17-03, TASK-F17-05 | Optional `contextRef` handling and UID enforcement |
| AC-F17-05 | TASK-F17-02, TASK-F17-03 | `null`, empty, partial rules |
| AC-F17-06 | TASK-F17-05 | Duplicate reference rejection |
| AC-F17-07 | TASK-F17-03 | List-order-independent digest |
| AC-F17-08 | TASK-F17-03 | Semantic change alters digest |
| AC-F17-09 | TASK-F17-03 | `requestId`-only change preserves digest |
| AC-F17-10 | TASK-F17-03 | Frozen RFC 8785/SHA-256 vector |
| AC-F17-11 | TASK-F17-04, TASK-F17-06, TASK-F17-07, TASK-F17-11 | All four outcomes |
| AC-F17-12 | TASK-F17-06, TASK-F17-07 | Missing/configured failure → `AdapterFailure` |
| AC-F17-13 | TASK-F17-07 | Fake construction failure; TASK-F17-08 separately proves D-13 without extending this AC |
| AC-F17-14 | TASK-F17-06 | Cancellation, deadline, at-most-once, late-output discard |
| AC-F17-15 | TASK-F17-06 | Invalid adapter conclusion → `InvalidAdapterResult` |
| AC-F17-16 | TASK-F17-01, TASK-F17-06, TASK-F17-10 | Fixed time source deterministic replay |
| AC-F17-17 | TASK-F17-04, TASK-F17-06, TASK-F17-07, TASK-F17-11 | Reason-code grammar, sorting, trusted pre-redaction, leakage tests |
| AC-F17-18 | TASK-F17-04, TASK-F17-11 | Obligations absent in Phase 2R fake |
| AC-F17-19 | TASK-F17-09, TASK-F17-11 | Result/timing transport and pure mapping with exact fields |
| AC-F17-20 | TASK-F17-09 | Malformed structural mapper-input failure |
| AC-F17-21 | TASK-F17-09, TASK-F17-11 | No mapping for non-result interactions |
| AC-F17-22 | TASK-F17-11, No-task ledger | No FEATURE-0017 DecisionRecord/AuditEvent/route/controller/store |
| AC-F17-23 | TASK-F17-10 | Concurrent equivalent-call determinism; copy ownership is separately pinned by D-10/D-12/D-14 |
| AC-F17-24 | TASK-F17-11 | Zero external effects |

## 5. Requirement / design / decision / risk coverage ledger

### 5.1 Architecture traceability

- **Controlling handoffs:** ADH-2026-067 (six-group closure), ADH-2026-068
  (structural mapper), ADH-2026-069 (timing/field/cardinality).
- **Accepted decisions:** DEC-0026 (reuse-before-build), DEC-0028 (policy via
  adapter), DEC-0036 (adapter boundaries), DEC-0043 (DecisionRecord profiles).
- **Canonical model references:** FEATURE-0013 `EvaluationResult` reuse,
  FEATURE-0015 CloudPlatform/CloudProvider scope terminology,
  FEATURE-0016 adapter-boundary precedent.
- **Owned schemas:** `VS0-SCHEMA-018` (`PolicyEvaluationRequest`),
  `VS0-SCHEMA-019` (`PolicyEvaluationResult`).
- **Reuse:** RFC 8785 JCS (Reuse), SHA-256 (Reuse), FEATURE-0012 `TypedRef` /
  transient / validation / redaction (Reuse), FEATURE-0013 `EvaluationResult`
  (Reuse), injected time-source pattern (Reuse).
- **Build:** `PolicyEngineAdapter` port, evaluation boundary, pure FEATURE-0013
  mapper, deterministic fake (all justified by necessity; no custom policy engine
  authorized).
- **Deferred:** Real OPA/Cedar/engine adapter, profile-kind taxonomy, production
  selector, retry, persistence, customer-safe explanation, entitlement/quota/
  sovereignty/placement semantics, obligations, time-dependent policy.

### 5.2 Risk mitigation

| Risk | Mitigation | Task |
|---|---|---|
| Architecture leak into requirements/design/tasks | Approved closed six-group architecture; strict delegation surface; conflict stops work with `ARCHITECTURE_DECISION_REQUIRED`. | All tasks |
| Request-value or adapter-error leakage in diagnostics | Sanitized fixed error table; no wrapping, interpolation, or echoing; leakage tests scan committed fixtures and returned errors. | TASK-F17-06, TASK-F17-11 |
| Unknown field or duplicate member bypass | Strict JSON decode with recursive token walk and per-object member-name set. | TASK-F17-02 |
| Digest instability | Frozen RFC 8785 / SHA-256 test vector; complete string test suite; v2 requires approved contract. | TASK-F17-03 |
| Invocation cardinality violation | Direct-call linearization on current goroutine; spy adapter proves at-most-once. | TASK-F17-06 |
| Cancellation / deadline race | Exact precedence with post-return `ctx.Err()` check before inspecting output; deterministic tie rule. | TASK-F17-06 |
| Fake fixture configuration error | Duplicate-preserving constructor; digest/conclusion validation; construction failure before partial fake. | TASK-F17-07 |
| Mapper aliasing or mutation | Deep-copy permitted mutable inputs and output raw JSON; reject every non-empty Phase 2R obligation value; mutation tests. | TASK-F17-09 |
| Accidental external dependency or effect | Import-surface test; package-level symbol assertion; zero-external-effect test. | TASK-F17-11 |
| Cross-feature ownership violation | Reuse FEATURE-0012/0013 by reference only; approved FEATURE-0013 extension; no earlier-feature mutation. | TASK-F17-02, TASK-F17-09 |

### 5.3 Explicit non-goals restated

No implementation task creates:

- authentication, PrincipalRef creation, membership, role, authorization
  resolution;
- governance, sovereignty, assignment, inheritance, exception, approval,
  entitlement, quota, placement, ranking, selection, execution semantics;
- policy/profile publication, profile-kind taxonomy, bundle distribution,
  loading, hot reload;
- `EffectiveGovernanceContext` construction or lookup;
- DecisionRecord or AuditEvent publication and any request/result store;
- customer-safe or CloudProvider-safe reason/explanation projection;
- public API routes, handlers, controllers, idempotency repositories, or
  persistent storage;
- production adapter selection, registry, routing, retry, or topology; or
- real OPA, Cedar, Kubernetes, identity, database, network, CloudProvider SDK,
  provisioning, or plugin execution.

## 6. Orphan and unresolved report

**Orphan IDs:** None. All REQ-F17-01..06 and AC-F17-01..24 are mapped above.

**Unresolved items:** None. All approved requirements, design decisions, and
acceptance cases have exact task owners or no-task justifications.

**Design elements without tasks:** See Section 3 (No-task ledger) for
contract-only, governance, invariant, and exclusion items.

**Architecture / requirements / design conflicts:** None detected. Any conflict
discovered during implementation must stop work with
`ARCHITECTURE_DECISION_REQUIRED`.

---

## 7. Execution receipt

**Model Execution Report:**

- Tool: kiro
- Stage or task: tasks
- Recommended priority list: claude-sonnet-4.5, claude-opus-4.8, (none listed third)
- Selected model: claude-sonnet-4.5
- Effort/reasoning setting: medium (task decomposition)
- Fallback used: no
- Fallback reason: none

STAGE_STATUS: COMPLETE

---

This tasks document decomposes the approved FEATURE-0017 design into 12
independently testable implementation units, preserving exact traceability to
REQ-F17-01..06, AC-F17-01..24, CDG-F17-01..06, and design decisions D-01..D-14.
Every task cites exact requirements, design sections, architecture decisions, and
applicable risks. Tasks sequence real dependencies (e.g., request before digest,
boundary before integration) and mark independent work explicitly where safe.
Each implementation task owns concrete writable paths, tests, verification
commands, acceptance criteria, and a commit message. TASK-F17-12 is deliberately
read-only and has no commit. No task adds semantics, public routes, stores,
controllers, production engines, external effects, or adjacent-feature behavior
absent from the approved design.

The final task (TASK-F17-12) performs read-only repository verification,
boundary and artifact review, focused and full-repository race passes,
baseline-relative worktree review, and
`make ff-feature-gate FEATURE=FEATURE-0017` after the required review receipts.
All AC-F17-01..24 and REQ-F17-01..06 are mapped; no orphan or unresolved item
remains. Founder executable-plan approval and independent tasks approval are
separately recorded external prerequisites to Cursor execution; this document
does not assert either receipt.
