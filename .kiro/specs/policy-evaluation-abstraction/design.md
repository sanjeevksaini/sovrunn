# FEATURE-0017 Policy Evaluation Abstraction — Design

## 1. Identity, stage, inputs, and closed boundary

| Field | Value |
|---|---|
| Feature | FEATURE-0017 — Policy Evaluation Abstraction |
| Stage | Design |
| Kiro slug | `policy-evaluation-abstraction` |
| Phase / order | Phase 2R / 7 |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Sole architecture authority | `docs/architecture/policy-evaluation-abstraction.md` |
| Controlling handoffs | ADH-2026-067 (closure), ADH-2026-068 (structural mapper), ADH-2026-069 (timing/field/cardinality) |
| Approved requirements | `.kiro/specs/policy-evaluation-abstraction/requirements.md` (REQ-F17-01..06, AC-F17-01..24) |
| Owned schemas | `VS0-SCHEMA-018` (`PolicyEvaluationRequest`), `VS0-SCHEMA-019` (`PolicyEvaluationResult`) |
| Public routes / persistence / controller / real engine | None |

This document translates the approved FEATURE-0017 requirements and the six
closed decision groups (CDG-F17-01 through CDG-F17-06) into an implementable Go
representation. It chooses only the mechanics explicitly delegated to design by
the sole architecture authority Section 9 ("Architecture-leak prevention") and
requirements Section 9: concrete Go package layout, private carrier and type
shapes, interface shapes, constructors, immutable fixture mechanics, the private
representation of the transient success return pairing, fixture serialization,
and goroutine/concurrency mechanics that preserve the exact observable
precedence and at-most-once invocation.

It adds no product semantics, transfers no ownership, reopens no closed
decision, and defines no mechanism owned by an excluded feature. Any conflict
between this document and a higher-precedence authority stops work with
`ARCHITECTURE_DECISION_REQUIRED`.

Closed boundary restated for design control:

- The digest identity, request fields, precedence, result/failure vocabulary,
  reason-code grammar, mapper populated/omitted fields, and side-effect
  boundaries are closed by CDG-F17-01..06 and are not design choices.
- FEATURE-0017 owns no public route, store, controller, idempotency repository,
  adapter selector, routing rule, mutable fixture reload, automatic retry,
  production engine, DecisionRecord, AuditEvent, or external call.
- Prior-feature contracts (FEATURE-0011/0012/0013/0015/0016) are consumed by
  exact reference only and are never forked, redefined, or re-owned.

### Verified repository conventions consumed

Design paths and reuse below were verified against the live repository:

- Go module path `github.com/sanjeevksaini/sovrunn`, Go 1.22, minimal
  dependencies (`go.mod` requires only `gopkg.in/yaml.v3`).
- Feature domains are single `internal/<domain>` packages with an optional
  `<domain>/model` subpackage (`internal/executiontarget`,
  `internal/cloudmodel`, `internal/decision`).
- FEATURE-0013 `EvaluationResult`, `EvaluatorIdentity`, `TimingEnvelope`,
  `OpaqueIntegrityCarrier`, `EvaluationResultStatus` (`SUCCESS`, `NON_MATCH`,
  `INDETERMINATE`, `TIMEOUT`, `ERROR`) live in `internal/decision`
  (`internal/decision/evaluationresult.go`).
- `TypedRef` is defined in `internal/apimeta` (`reference.go`) and re-exported
  by `internal/apiref`; structural reference validation is
  `apiref.Constraint.ValidateRef` returning package-local `apiref.RefIssue`.
- FEATURE-0016's injected clock precedent is `internal/executiontarget.Clock`
  (`Now() time.Time`) with a deterministic `FakeClock`; it is a precedent
  pattern only and is not imported by FEATURE-0017 (see §2, D-06).
- No RFC 8785 JCS canonicalizer or shared canonical-JSON digest helper exists in
  the repository today; FEATURE-0017 supplies its own stdlib-only v1 serializer
  (see §2, D-03).

## 2. Resolved design decisions

Each decision below is a design mechanic permitted by the delegation surface. No
decision alters a closed CDG. Decisions are numbered `D-01..D-14`.

- **D-01 — Single feature package `internal/policyeval`.** All FEATURE-0017 code
  lives in one new package `internal/policyeval`, matching the single-package
  convention of `internal/executiontarget` and `internal/decision` and honoring
  the architecture's "minimal private packages" directive. No new top-level
  module, no `cmd/` entry, no server wiring, and no route registration are
  added. The deterministic fake lives in the same package (file `fake.go`)
  because it is in-process conformance machinery, not a separately deployed
  component; a separate package is unnecessary and would not reduce coupling.

- **D-02 — Closed Go values plus one feature-owned strict JSON decoder.**
  `PolicyEvaluationRequest` and `PolicyEvaluationResult` are plain Go structs
  with exact lower-camel JSON tags and no `TypeMeta`/`ObjectMeta`. Both are
  `TransientRequestResult`-profile values (VS0-SCHEMA-018/019); following the
  FEATURE-0013 `EvaluationResult` EmbeddedValue/transient precedent, they carry
  no artificial identity, scope field, or status. Byte-originated requests
  enter only through `DecodePolicyEvaluationRequestJSON`. Its private wire DTO
  uses `json.RawMessage` presence carriers for every field, rejects unknown
  members and duplicate object-member names at the root and inside every
  reference, rejects trailing JSON, and distinguishes omitted from explicit
  `null` before constructing the public Go value. A token preflight maintains a
  separate member-name set for every object; any repeated name fails closed
  before DTO decoding. A lexical preflight rejects invalid UTF-8 and
  lone/malformed UTF-16 surrogate escapes before `encoding/json` can replace
  them. Direct Go construction represents semantic values only: nil
  `ContextRef` is absence, and nil/empty `CandidateRefs` are the approved empty
  set. `PolicyEvaluationRequest.UnmarshalJSON` delegates to that same private
  strict decoder, so ordinary `json.Unmarshal` cannot bypass it. YAML request
  decoding is unsupported.

- **D-03 — Stdlib-only, domain-bounded RFC 8785 JCS v1 serializer inside
  `internal/policyeval`.** The normalized v1 logical object contains only
  strings, fixed-name objects, and arrays of reference objects—no numbers,
  booleans, or nulls—so FEATURE-0017 implements exactly that JCS subset with a
  byte builder plus `crypto/sha256` and `encoding/hex`; it does not delegate
  canonical string emission to `encoding/json`. Root keys are emitted in JCS
  UTF-16 order (`action`, `candidateRefs`, optional `contextRef`, `profileRefs`,
  `schema`, `subjectRef`) and reference keys in order (`apiVersion`, `kind`,
  `name`, optional `uid`). These keys are fixed ASCII, so their precomputed
  order is also their RFC 8785 order. Every key and value passes the same string
  emitter: invalid UTF-8 and non-scalar Unicode are rejected; no Unicode
  normalization occurs; quote is
  escaped as `\"` and reverse-solidus as `\\`; U+0008/U+0009/U+000A/U+000C/U+000D use the
  JSON short escapes; every other U+0000..U+001F code point uses lowercase
  `\u00xx`; and all other scalar values are emitted as their original UTF-8,
  including `<`, `>`, `&`, `/`, U+2028, and U+2029. The output has no
  insignificant whitespace, BOM, trailing newline, HTML escaping, or external
  prefix/suffix. Tests cover the frozen ASCII vector plus controls,
  quote/reverse-solidus, HTML-sensitive characters, non-ASCII BMP and
  supplementary characters, U+2028/U+2029, invalid UTF-8, and composed versus
  decomposed strings. No third-party JCS dependency is added.

- **D-04 — Reuse the exact FEATURE-0012 reference validator; do not fork.**
  Every reference field uses `apimeta.TypedRef`. Singular references pass to
  `apiref.Constraint.ValidateRef`; collections pass once to
  `apiref.Refs.Validate` with their approved limits. Both shared paths enforce
  required `apiVersion`/`kind`/`name` and provider-native-ID rejection. The
  `contextRef` constraint has
  `AllowedKinds: []string{"EffectiveGovernanceContext"}` and then enforces a
  non-empty UID; every `profileRefs` element likewise enforces UID after the
  shared structural pass. `subjectRef` and `candidateRefs` use the unconstrained
  kind set registered by VS0-SCHEMA-018. FEATURE-0017 adds only its approved UID,
  semantic-set duplicate, and ordering rules; it neither copies
  `ValidateRef` rules nor resolves a reference.

- **D-05 — Normalized-reference identity for sorting and duplicate detection.**
  A normalized reference is the tuple `(apiVersion, kind, name, uid-or-empty)`.
  `profileRefs` and `candidateRefs` sort by that tuple using bytewise UTF-8
  lexical ordering (Go native `<` on the raw UTF-8 strings, compared
  field-by-field in that order). Duplicate identity is equality of the whole
  normalized reference value; duplicates are rejected, never deduplicated.

- **D-06 — Local injected UTC time source; FEATURE-0016 clock is precedent
  only.** FEATURE-0017 defines a minimal in-package time-source seam
  (`type TimeSource interface { NowUTC() time.Time }` with a production UTC
  implementation and an immutable fixed test implementation). It does **not**
  import `internal/executiontarget.Clock`, avoiding cross-feature coupling to
  FEATURE-0016 lifecycle machinery. `NewBoundary` rejects nil and typed-nil time
  sources, and an implementation supplied for concurrent use must itself be
  concurrency-safe. The boundary normalizes each returned `time.Time` with
  `.UTC()` and serializes it using `time.RFC3339Nano`. The seam is a
  dependency-injection pattern, not a resource, store, or platform service.

- **D-07 — Unnamed transient success pairing.** Success is represented only by
  two adjacent return values: the complete `PolicyEvaluationResult` and the
  exact boundary-captured FEATURE-0013 `TimingEnvelope`. No exported or private
  carrier struct is declared. The separate mapper accepts those two values as
  separate arguments. The pairing therefore has no named Go contract, API
  identity, schema, resource kind, persistence, projection, retention, or
  lifecycle.

- **D-08 — Normalized non-result category is an internal typed value, never a
  Problem code.** The five categories (`RequestInvalid`, `Canceled`,
  `DeadlineExceeded`, `AdapterFailure`, `InvalidAdapterResult`) are represented
  by an in-package typed enum returned alongside zero result and zero timing
  values. They are not
  FEATURE-0012 top-level Problem codes, Problem violation IDs, HTTP status
  mappings, or a fifth outcome, matching the requirements zero-inventory
  traceability. No `apiproblem` type is produced by this feature.

- **D-09 — Constructed boundary around a conclusion-or-failure adapter port.**
  The port is
  `type PolicyEngineAdapter interface { Evaluate(context.Context, NormalizedInput, string) (Conclusion, error) }`.
  A non-nil adapter error is only the internal signal for normalized
  `AdapterFailure`; the boundary never wraps, returns, or logs that error.
  `Conclusion` carries only `{Outcome, ReasonCodes}`—no digest, time, timing,
  obligations, DecisionRecord identity, or AuditEvent identity. `NewBoundary`
  captures exactly one adapter and one time source, rejects nil and typed-nil
  dependencies, and returns no partially usable boundary. The adapter contract
  requires context awareness, eventual return, and concurrency safety when one
  boundary is called concurrently.

- **D-10 — Boundary owns conclusion validation, defensive ownership, digest,
  and time.** After a nil-error return, the boundary first copies the adapter's
  entire conclusion, then validates outcome membership and reason-code
  grammar/count/uniqueness, sorts only its owned copy, and supplies the
  precomputed digest and injected UTC completion time. Any invalid or duplicate
  reason code, or any literal outcome outside the closed set, yields
  `InvalidAdapterResult` with no result. Exit-criteria prose `unknown` maps to
  `Indeterminate` only as vocabulary reconciliation, never as an adapter-input
  alias.

- **D-11 — Direct-call invocation linearization with cancellation-first return
  precedence.** After the immediate pre-invocation `ctx.Err()` check succeeds,
  the boundary captures `startedAt` and directly calls `adapter.Evaluate` on the
  current goroutine; entry to that single call is the invocation linearization
  point. When the call returns, the boundary checks `ctx.Err()` *before*
  inspecting either the conclusion or adapter error. An observed cancellation
  or deadline therefore deterministically wins every return-time tie and the
  returned adapter output is discarded. If the context is not done, adapter
  failure precedes conclusion validation. No boundary goroutine, result channel,
  retry, or second call exists. Go cannot forcibly terminate an adapter that
  never returns; eventual return after context cancellation is an explicit port
  liveness requirement. A violating adapter can block its caller and is outside
  the conforming adapter contract, but cannot cause a FEATURE-0017 goroutine
  leak because the boundary creates none.

- **D-12 — Immutable, digest-keyed fake with explicit copy ownership.** The fake
  constructor accepts `[]FakeFixture`, never a map, so repeated digests survive
  input long enough to be rejected. Each fixture contains one `InputDigest`, an
  optional `*Conclusion`, and `AdapterFailure bool`; exactly one of conclusion
  or failure must be selected. `NewDeterministicFake([]FakeFixture)` deep-copies
  every entry and nested reason-code slice into a private map, then validates the
  copied values. It fails without a partial fake on duplicate/malformed digest,
  invalid one-of selection, or malformed conclusion. Each invocation
  returns a fresh copy of the configured conclusion and reason-code slice, so
  boundary sorting or caller mutation cannot alter shared fixture state. A
  configured failure and an unconfigured digest both signal `AdapterFailure`,
  never `Indeterminate`. The fake holds no time, mutable reload path, or
  conditional interpreting domain meaning. Fixture strings are treated as
  trusted, repository-controlled, pre-redacted configuration; grammar checks
  shape but do not claim to detect semantic sensitivity. Construction errors
  are fixed generic diagnostics and never echo a digest, reason code, or fixture
  value.

- **D-13 — Fixture serialization: in-memory Go config, optional strict YAML for
  tests only.** Fixtures are supplied to the constructor as Go values. Test-only
  helper `decodeFakeFixturesYAMLForTest([]byte)` lives in a `_test.go` file and
  uses the existing `gopkg.in/yaml.v3` dependency with `KnownFields(true)`,
  duplicate-key rejection, exactly one document, and rejection of trailing
  content. It returns the duplicate-preserving fixture slice, which still passes
  through `NewDeterministicFake`. YAML is not accepted for
  `PolicyEvaluationRequest`, and no runtime reload is added.

- **D-14 — Pure mapper with exact FEATURE-0012/0013 structural-contract reuse.**
  `MapToEvaluationResult(result, timing, evaluator, ident)` accepts the two
  unnamed success values separately. It validates the FEATURE-0017 result,
  delegates every TypedRef check to `apiref.Constraint.ValidateRef`, constructs
  the closed candidate `decision.EvaluationResult`, and calls the approved
  FEATURE-0013 dependency extension
  `decisionvalidate.IsEvaluationResultStructurallyValid(candidate)`. That
  exported boolean wrapper delegates directly to FEATURE-0013's existing private
  `checkEvaluationResultStructure`; it adds no rule and exposes no Problem or
  violation. The mapper deliberately does not call FEATURE-0013
  `ValidateEvaluation`, because that adoption validator also requires a
  `DecisionProfile` and DecisionRecord scope that CDG-F17-05 excludes.
  The mapper copies snapshot/integrity pointers, timing, result reason codes,
  obligations (normally absent), and the marshaled `json.RawMessage`; its output
  aliases no caller-owned mutable memory. It performs no time read, lookup,
  profile/scope compatibility validation, write, persistence, or external
  effect.

  Dependency-extension approval: Sanjeev Kumar explicitly approved this
  structural-only boolean wrapper on 2026-08-25. The extension adds one wrapper
  file and tests in the FEATURE-0013-owned validation package; it changes no
  existing validation behavior or ownership.

## 3. Components and repository paths

FEATURE-0017 code remains in one new package. The sole approved dependency
extension is one new boolean-wrapper file (plus tests) in the existing
FEATURE-0013 validation package; no existing FEATURE-0013 file or behavior is
modified.

| Component | Path | Responsibility |
|---|---|---|
| Package doc | `internal/policyeval/doc.go` | Package purpose, non-goals, and closed-boundary note. |
| Request type + strict JSON decode + validation | `internal/policyeval/request.go` | Exact JSON tags, presence-aware strict decoder, duplicate-member/unknown/trailing rejection, field bounds, shared reference validation, UID/kind/set rules (REQ-F17-01). |
| Canonicalization + digest | `internal/policyeval/canonical.go` | v1 normalized logical object, RFC 8785 JCS serialization, SHA-256 lowercase-hex digest, absent/empty handling (REQ-F17-02). |
| Result + vocabulary | `internal/policyeval/result.go` | `PolicyEvaluationResult`, `Outcome` enum, reason-code grammar, non-result category enum (REQ-F17-03). |
| Adapter port + normalized input | `internal/policyeval/adapter.go` | `PolicyEngineAdapter`, `NormalizedInput`, `Conclusion`, and non-nil-error adapter-failure signal (REQ-F17-03). |
| Time source | `internal/policyeval/timesource.go` | `TimeSource` interface, UTC production impl, fixed deterministic impl (REQ-F17-03/06). |
| Evaluation boundary operation | `internal/policyeval/evaluate.go` | Dependency-checked constructor, direct-call linearization, exact precedence, timing, sanitized diagnostics, result construction, unnamed result/timing return (REQ-F17-03/05/06). |
| Deterministic fake | `internal/policyeval/fake.go` | `FakeFixture`, duplicate-preserving constructor input, immutable digest-keyed fake, one-of and construction validation (REQ-F17-04). |
| Test-only fixture decoder | `internal/policyeval/fake_fixture_test.go` | Strict one-document YAML-to-`[]FakeFixture` helper; no production input surface. |
| Pure FEATURE-0013 mapper | `internal/policyeval/mapper.go` | Exact `InputIdentity`, independent-copy construction, shared structural-validator call, closed status/field/timing mapping (REQ-F17-05). |
| FEATURE-0013 structural reuse wrapper | `internal/decision/validate/evaluation_structure_export.go` | Approved bool-only delegate to existing private FEATURE-0013 structure validation; no Problem exposure or new rule. |
| FEATURE-0013 wrapper parity tests | `internal/decision/validate/evaluation_structure_export_test.go` | Prove the wrapper returns the exact existing structural-pass result. |
| Unit/property/conformance tests | `internal/policyeval/*_test.go` | AC-F17-01..24 proof (see §7). |

Dependency direction (all inward, acyclic): `internal/policyeval` imports
`internal/apimeta`, `internal/apiref`, `internal/decision`, and
`internal/decision/validate` (as `decisionvalidate`), plus stdlib
(`bytes`, `context`, `crypto/sha256`, `encoding/hex`, `encoding/json`, `errors`,
`io`, `reflect`, `sort`, `strings`, `time`, `unicode/utf8`).
It does not import `internal/executiontarget`, `internal/cloudmodel`,
`internal/server`, `internal/api`, or `internal/apiproblem`.

## 4. Data / API representation

No wire API is defined (no route, no handler). "API" here means the in-process
Go contract. Field names and JSON tags mirror the closed schema vocabulary.

### 4.1 `PolicyEvaluationRequest` (VS0-SCHEMA-018)

```go
type PolicyEvaluationRequest struct {
    SubjectRef    apimeta.TypedRef   `json:"subjectRef"`
    Action        string             `json:"action"`
    ContextRef    *apimeta.TypedRef  `json:"contextRef,omitempty"`
    ProfileRefs   []apimeta.TypedRef `json:"profileRefs"`
    CandidateRefs []apimeta.TypedRef `json:"candidateRefs,omitempty"`
    RequestID     string             `json:"requestId"`
}
```

- `DecodePolicyEvaluationRequestJSON([]byte)` is the only byte decoder. A
  raw-byte `utf8.Valid` check and a string-token preflight for correctly paired
  surrogate escapes run before `encoding/json` can replace invalid input. A
  recursive `json.Decoder.Token` walk keeps a fresh set of member names per
  object and rejects a repeated name anywhere, including at the root and within
  every reference. A private root DTO stores every member as `json.RawMessage`; a decoder
  with `DisallowUnknownFields` reads exactly one object and a second decode must
  return `io.EOF`. Each reference is decoded through the same one-value,
  unknown-member-rejecting helper; reference arrays first decode to
  `[]json.RawMessage` so nested unknown fields cannot bypass strictness.
- `PolicyEvaluationRequest.UnmarshalJSON` calls the same private decoder (which
  decodes into the wire DTO, not recursively into the public type), making
  strict behavior mandatory for all standard-library JSON entry paths.
- A nil raw member means omitted. A raw token equal to `null` is explicit null.
  Required members reject both; omitted `contextRef` is valid while explicit
  null is invalid; omitted `candidateRefs` constructs an owned empty slice,
  explicit null is invalid, and `[]` is valid. `{}` and partial references
  decode but fail the reused structural validator. Empty required strings and
  an empty/null/omitted `profileRefs` fail before digest calculation.
- Direct Go callers invoke the same semantic validator. They cannot represent
  unknown members or a JSON explicit-null token; nil `ContextRef` therefore
  means absence, while nil and empty `CandidateRefs` both mean the approved
  empty set. YAML request input is unsupported.
- `RequestID` is correlation only; it is excluded from the normalized object and
  from any adapter behavior and result identity.
- No `evaluationClass`, `metadata.scopeRef`, status, identity, version, or
  timestamp field exists.

### 4.2 v1 normalized logical object (VS0-SCHEMA-018 → digest)

The normalized object serialized by JCS contains exactly, with object members
emitted in JCS-sorted order:

- `schema`: constant `"sovrunn.policy-evaluation-request/v1"`;
- `subjectRef`: normalized reference;
- `action`: exact string;
- `contextRef`: present only when the input `contextRef` is present;
- `profileRefs`: sorted array of normalized references;
- `candidateRefs`: always present; omitted input serializes as `[]`.

`requestId` is excluded. A normalized reference contains `apiVersion`, `kind`,
`name`, and includes `uid` only when present and valid. Absent/empty handling is
exactly as CDG-F17-02: omitted `contextRef` is valid and absent; `contextRef`
`null`/`{}`/partial is invalid; omitted and `[]` `candidateRefs` are the same
empty set; `candidateRefs: null` is invalid; empty required strings/objects and
an empty `profileRefs` set fail validation before digest.

Normalization first validates UTF-8 in every string and deep-copies the request,
including the context pointer and both slices. Sorting and serialization touch
only that owned snapshot. The caller must honor the registered immutable-input
contract and not mutate a request concurrently with handoff; mutations after
handoff cannot affect the snapshot, digest, or adapter input.

### 4.3 `PolicyEvaluationResult` (VS0-SCHEMA-019)

```go
type Outcome string // Allow | Deny | Indeterminate | RequiresApproval

type PolicyEvaluationResult struct {
    Outcome     Outcome           `json:"outcome"`
    ReasonCodes []string          `json:"reasonCodes"`            // 1..32, unique, lexically sorted, ^[A-Z][A-Z0-9_]{0,62}$
    InputDigest string            `json:"inputDigest"`            // lowercase 64-char sha256 hex
    EvaluatedAt string            `json:"evaluatedAt"`            // UTC RFC3339
    Obligations []json.RawMessage `json:"obligations,omitempty"`  // absent in the Phase 2R fake
}
```

`Obligations` is always empty (nil) in Phase 2R, so it never serializes
(AC-F17-18). A boundary result owns its reason-code slice and aliases neither
adapter nor fixture storage.

### 4.4 Adapter port and normalized input

```go
type NormalizedInput struct {
    SubjectRef    apimeta.TypedRef
    Action        string
    ContextRef    *apimeta.TypedRef  // present only when supplied
    ProfileRefs   []apimeta.TypedRef // sorted
    CandidateRefs []apimeta.TypedRef // sorted; empty slice when none
}

type Conclusion struct {
    Outcome     Outcome  `yaml:"outcome"`
    ReasonCodes []string `yaml:"reasonCodes"`
}

type PolicyEngineAdapter interface {
    Evaluate(ctx context.Context, in NormalizedInput, inputDigest string) (Conclusion, error)
}
```

A non-nil error returned by `Evaluate` is normalized to the `AdapterFailure`
category and is discarded without wrapping, logging, or returning its text.
`RequestID` is not passed to the adapter as a behavior input. Immediately before
the call the boundary gives the adapter a deep copy of its normalized snapshot;
adapter mutation cannot alter the digest snapshot or any result. A conforming
adapter honors cancellation/deadlines, eventually returns, and is safe for the
concurrent use promised by its construction. It transfers ownership of the
returned `Conclusion` and must not mutate its slice concurrently or after
return; the boundary nevertheless copies it before validation or sorting.

### 4.5 Unnamed transient success pairing and non-result

```go
type NonResult string // RequestInvalid | Canceled | DeadlineExceeded | AdapterFailure | InvalidAdapterResult

type Boundary struct { /* exactly one adapter and one time source; fields private */ }

func NewBoundary(adapter PolicyEngineAdapter, ts TimeSource) (*Boundary, error)

func (b *Boundary) Evaluate(ctx context.Context, req PolicyEvaluationRequest) (
    PolicyEvaluationResult,
    decision.TimingEnvelope,
    NonResult,
    error,
)
```

`NewBoundary` rejects nil and typed-nil dependencies using a small stdlib
nil-interface check; it never returns a partially usable value. On success,
`Evaluate` returns a populated result and timing envelope, an empty `NonResult`,
and nil error. On a normalized non-result it returns zero result and zero timing,
the specific category, and a fixed generic diagnostic error. Diagnostics never
wrap adapter errors or include request, reference, digest, outcome, reason-code,
or fixture values. A nil context violates the operation precondition and fails
closed before invocation with the same sanitized request-invalid surface; it is
never substituted with a background context.

### 4.6 FEATURE-0013 `EvaluationResult` mapper output

`InputIdentity` is the complete mapper-owned in-process input-identity contract:

```go
type InputIdentity struct {
    InputSnapshotRef *apimeta.TypedRef
    InputIntegrity   *decision.OpaqueIntegrityCarrier
}
```

Valid combinations are snapshot only, integrity only, or both. Both nil is
invalid. A present snapshot must pass `apiref.Constraint.ValidateRef`; a
present integrity carrier must be non-empty with a non-blank `State`. Empty or
partial pointed-to values are invalid. The mapper copies each present value into
new storage before constructing output, so neither output pointer aliases the
caller. `InputIdentity` has no JSON tags or serialization path.

The mapper signature carries no named success receipt:

```go
func MapToEvaluationResult(
    result PolicyEvaluationResult,
    timing decision.TimingEnvelope,
    evaluator decision.EvaluatorIdentity,
    ident InputIdentity,
) (decision.EvaluationResult, error)
```

`MapToEvaluationResult` populates exactly (CDG-F17-05, ADH-2026-069):

- `Evaluator` from caller-supplied `decision.EvaluatorIdentity`;
- `InputSnapshotRef`, `InputIntegrity`, or both, exactly as supplied after
  structural validation (at least one required);
- `ResultStatus`: `Allow`/`Deny`/`RequiresApproval` → `SUCCESS`,
  `Indeterminate` → `INDETERMINATE`;
- `Result`: the complete `PolicyEvaluationResult` marshaled to `json.RawMessage`;
- `Timing.StartedAt` and `Timing.CompletedAt` from the supplied timing envelope;
- `EvaluatedAt` identical to both `PolicyEvaluationResult.EvaluatedAt` and
  `Timing.CompletedAt`.

It always leaves absent: `Timing.DurationMs`, `ExecutionConfig`, `TrustBoundary`,
`SafetyPolicyFilters`, `StructuralTrustState`, and `ScopeEvidence`. These are
never conditionally populated.

### 4.7 Deterministic fake configuration and constructor

```go
type FakeFixture struct {
    InputDigest   string      `yaml:"inputDigest"`
    Conclusion    *Conclusion `yaml:"conclusion,omitempty"`
    AdapterFailure bool       `yaml:"adapterFailure,omitempty"`
}

type DeterministicFake struct { /* private map[string]fakeBehavior */ }

func NewDeterministicFake(fixtures []FakeFixture) (*DeterministicFake, error)

func (f *DeterministicFake) Evaluate(
    ctx context.Context,
    in NormalizedInput,
    inputDigest string,
) (Conclusion, error)
```

The slice is intentionally duplicate-preserving. An empty slice is a valid fake
whose every digest is unconfigured. For each entry, exactly one of non-nil
`Conclusion` or `AdapterFailure == true` is required; neither or both is invalid.
The constructor checks the digest grammar, duplicate digest, and configured
conclusion, copying input and nested slices before retaining them. A configured
failure or absent digest returns a zero `Conclusion` and a private generic
non-nil error used only as the adapter-failure signal. Each configured conclusion
return is freshly copied. No map-taking constructor or mutation/reload method
exists.

The test-only helper has the exact signature:

```go
func decodeFakeFixturesYAMLForTest(data []byte) ([]FakeFixture, error)
```

It performs the strict one-document YAML rules from D-13 and returns a slice;
all semantic validation remains centralized in `NewDeterministicFake`.

### 4.8 Approved FEATURE-0013 structural-validation extension

The exact added surface in package `internal/decision/validate` is:

```go
// IsEvaluationResultStructurallyValid exposes the existing FEATURE-0013
// structural pass without profile, record-scope, limit, Problem, or violation
// output.
func IsEvaluationResultStructurallyValid(result decision.EvaluationResult) bool {
    return checkEvaluationResultStructure(result) == nil
}
```

It is a thin delegate to the existing private function in the same package.
The wrapper contains no condition other than the nil result check, and
FEATURE-0013 parity tests cover representative valid and invalid structures.
FEATURE-0017 treats `false` as a generic sanitized mapping error and cannot
observe the private Problem value.

## 5. Validation and deterministic error behavior

### 5.1 Exact observable precedence

`Evaluate` implements the closed nine-step precedence without reordering
(CDG-F17-06). The category returned at each failing step:

1. Already-observed cancellation or expired deadline (`ctx.Err()` before any
   work) → `Canceled` or `DeadlineExceeded`; zero adapter invocations.
2. Request structural / reference / bound validation → `RequestInvalid` on
   failure; no adapter invocation.
3. Canonicalization and SHA-256 digest of the normalized object.
4. Cancellation/deadline recheck immediately before invocation → `Canceled` or
   `DeadlineExceeded`; zero adapter invocations.
5. Capture `startedAt` from the injected time source and perform exactly one
   direct `adapter.Evaluate` call on the current goroutine; call entry is the
   invocation linearization point.
6. Immediately after return, inspect `ctx.Err()` before adapter output. If
   cancellation/deadline occurred during invocation, discard all returned
   output and return `Canceled`/`DeadlineExceeded`.
7. Only when context remains active, a non-nil adapter error becomes the
   sanitized `AdapterFailure` category; the raw error is discarded.
8. Conclusion validation (outcome membership; reason-code grammar, count 1..32,
   uniqueness) → `InvalidAdapterResult` on failure.
9. Capture `completedAt`, lexically sort reason codes, set `evaluatedAt =
   completedAt`, construct the complete result, and return the transient
   unnamed result/timing pairing.

`Canceled` versus `DeadlineExceeded` is distinguished by `ctx.Err()`
(`context.Canceled` vs `context.DeadlineExceeded`). The post-call context check
is the deterministic tie rule: observed cancellation/deadline wins over both a
simultaneously available adapter conclusion and adapter error. A conforming
adapter must eventually return after its context is done; the boundary cannot
terminate a non-returning call and creates no helper goroutine that could leak.

### 5.2 Request validation (step 2) order

For byte input, strict decoding precedes the operation: one root object, no
duplicate or unknown root/nested-reference members, exactly one JSON value, and
explicit presence/null classification. Within step 2 the deterministic order is:
(a) required presence/null/empty and string bounds; (b) collection bounds;
(c) exact FEATURE-0012 structural validation of every reference through
`apiref.Constraint.ValidateRef`; (d) `contextRef.kind ==
EffectiveGovernanceContext` plus its required UID and each profile UID;
(e) duplicate detection over normalized references. Unknown/extra input,
explicit invalid null, an empty/partial reference, or any step-2 failure yields
only sanitized `RequestInvalid`, no digest, and no adapter invocation.

### 5.3 Determinism and concurrency

- Given the same request, adapter configuration, and injected time source,
  results and `evaluatedAt` are identical on repeated and concurrent calls
  (AC-F17-16, AC-F17-23). The fake's frozen map and the pure canonicalizer hold
  no mutable state; `Evaluate` shares no mutable state across calls.
- The boundary invokes directly and creates no goroutine or channel. The single
  guarded call site proves at-most-once cardinality; the adapter's explicit
  context/liveness contract bounds cancellation completion.
- Ownership is closed at every mutable boundary: strict decode creates owned
  slices; normalization deep-copies request pointers/slices; adapter handoff is
  another deep copy; the fake copies construction input and each returned
  conclusion; boundary validation/sorting uses a fresh reason-code copy; and
  mapper construction copies all pointers, slices, obligations, and raw JSON.
  No sort is performed on caller-, adapter-, or fixture-owned storage.
- Sorting is total and stable over the `(apiVersion, kind, name, uid-or-empty)`
  tuple, so canonical bytes are input-order-independent (AC-F17-07) while a
  semantic change alters the digest (AC-F17-08) and a `requestId`-only change
  does not (AC-F17-09).

### 5.4 Non-result semantics vs Problem codes

The five `NonResult` values are internal evaluation-operation categories. They
are never emitted as FEATURE-0012 top-level Problem codes, Problem violation
IDs, or HTTP statuses, and they are never a fifth `Outcome`. FEATURE-0017
introduces no public or registry-addressable violation identifier (requirements
zero-inventory traceability).

Returned diagnostic text is selected only from this closed safe table; values
are unexported errors, not registry identifiers or Problem codes:

| Situation | Exact diagnostic text |
|---|---|
| Invalid boundary dependency | `policy evaluation dependency invalid` |
| `RequestInvalid` | `policy evaluation request invalid` |
| `Canceled` | `policy evaluation canceled` |
| `DeadlineExceeded` | `policy evaluation deadline exceeded` |
| `AdapterFailure` | `policy evaluation adapter failed` |
| `InvalidAdapterResult` | `policy evaluation adapter result invalid` |
| Invalid fake construction | `policy evaluation fake configuration invalid` |
| Invalid mapper input | `policy evaluation mapping input invalid` |

No underlying error or input value is joined to these messages.

### 5.5 Mapper validation

`MapToEvaluationResult` first validates the complete FEATURE-0017 result
(closed outcome; exact reason-code, digest, obligations-absent, and UTC time
rules) and applies `apiref.Constraint.ValidateRef` to any snapshot reference. It
checks the FEATURE-0017-specific equality `result.evaluatedAt ==
timing.completedAt`, constructs the exact closed candidate
`decision.EvaluationResult`, then calls
`decisionvalidate.IsEvaluationResultStructurallyValid(candidate)`. That shared
FEATURE-0013 wrapper is the sole evaluator/input/timing/result structural pass;
`internal/policyeval` does not reproduce its rules. The mapper does not call the
broader FEATURE-0013 `ValidateEvaluation`, because that function also performs
adopting `DecisionProfile`, DecisionRecord-scope, and limit checks explicitly
excluded by ADH-2026-068.

Malformed input returns a fixed sanitized mapping error, leaves the caller's
valid policy result unchanged, and produces no partial output or side effect.
On success, every mutable field is copied into independently owned output. The
later adopting domain remains responsible for profile compatibility,
DecisionRecord scope, limits, adoption, and publication.

## 6. Security, privacy, observability, compatibility, and versioning

Security and privacy:

- `PolicyEvaluationRequest` and `PolicyEvaluationResult` are Tenant-confidential,
  internal-engine-facing values (VS0-SCHEMA-018/019, `redaction: redact`); the
  design defines no public projection, response body, or route.
- Request strings and references are themselves Tenant-confidential and may not
  be logged, echoed, wrapped into errors, or projected. Credentials and
  CloudProvider-native identifiers remain prohibited contract input; the shared
  FEATURE-0012 validator rejects provider-native reference shapes.
- Reason-code grammar constrains shape only; it cannot prove that a syntactically
  valid token is non-sensitive. Phase 2R fixture configuration is trusted,
  repository-controlled, and pre-redacted before construction. The fake exposes
  no free-text explanation field, and conformance tests scan committed fixtures
  plus sentinel adapter/request inputs to prove that returned errors and any
  permitted diagnostics do not reproduce sensitive values.
- Every dependency/configuration/request/adapter/mapping failure returns its
  closed `NonResult` (where applicable) plus a generic fixed diagnostic. No raw
  adapter error is wrapped with `%w`/`%v`; no request field, reference value,
  digest, reason code, fixture value, or engine diagnostic is interpolated.
- The seam performs no authentication, authorization, or principal resolution
  and exposes no network endpoint, so no route-level access control is
  introduced (excluded; see §10). It is a trusted in-process operation.
- `requestId` is correlation metadata only, never identity or an idempotency
  key.

Observability:

- The seam emits no AuditEvent, no DecisionRecord, and no persistent record
  (AC-F17-22). It follows the inherited observability and audit baseline
  (`docs/architecture/observability-and-audit-baseline.md`) for any in-process
  structured diagnostics by reference only. Diagnostics may contain only the
  operation name and normalized non-result category; they contain no request,
  result, reason-code, digest, adapter-error, fixture, or reference values. No
  metric, span, or side channel owned by a downstream feature is introduced.

Compatibility and versioning:

- FEATURE-0012 `TypedRef`, transient-contract, structural-validation, and
  redaction foundations are reused by reference (`internal/apimeta`,
  `internal/apiref`) and never redefined.
- FEATURE-0013 `EvaluationResult`, `EvaluatorIdentity`, `TimingEnvelope`, and
  `OpaqueIntegrityCarrier` are reused by reference (`internal/decision`) for the
  mapper. `DecisionProfile`, `DecisionRecord`, and `AuditEvent` remain
  FEATURE-0013-owned; FEATURE-0017 neither receives nor constructs them.
- The v1 digest identity is fixed and pinned by the frozen test vector. A future
  semantic field requires an approved v2 digest contract and a new `schema`
  constant; v1 is never silently changed (REQ-F17-02).
- The `PolicyEngineAdapter` port is stable so a later approved real adapter can
  be wrapped behind it without changing this contract; production adapter
  selection is out of scope.

Migration impact: none. This is a new, additive, in-process package. There is no
persistence, schema migration, route change, or data backfill. Existing packages
are unchanged.

## 7. Test and conformance strategy

Tests live in `internal/policyeval/*_test.go` and run under `go test`,
`go test -race`, and `go vet` (repository verification commands). Every approved
acceptance case maps to explicit coverage. No test asserts a downstream metric,
side effect, controller, adapter/plugin execution, or `VS0-CF-*` runtime proof;
FEATURE-0017's acceptance inventory is `AC-F17-01..24`, not a `VS0-CF-*` range.

### Conformance Mapping (AC → design coverage)

| AC | Coverage |
|---|---|
| AC-F17-01 | Request validation table test at every registered bound (`action` 1 and 63; `profileRefs` 1 and 32; `candidateRefs` 0 and 64; `requestId` 1 and 128). |
| AC-F17-02 | Strict JSON tests reject duplicate and unknown root/nested-reference members, trailing values/content, and malformed references. |
| AC-F17-03 | Spy adapter asserting zero invocations on request-validation failure. |
| AC-F17-04 | Absent `contextRef` and valid present UID-pinned `contextRef` both accepted. |
| AC-F17-05 | Presence-aware decode matrix covers omitted vs explicit `null`, `{}`, partial refs, omitted/`[]`/null candidates, and omitted/null/empty profiles before digest. |
| AC-F17-06 | Duplicate `profileRefs`/`candidateRefs` rejected (not deduplicated). |
| AC-F17-07 | Reordered reference lists produce identical digest. |
| AC-F17-08 | Semantic change alters digest. |
| AC-F17-09 | `requestId`-only change preserves digest. |
| AC-F17-10 | Frozen preimage/SHA-256 vector plus JCS string vectors for controls, escaping, UTF-8, U+2028/U+2029, HTML-sensitive characters, supplementary code points, invalid UTF-8/lone surrogates, and no normalization. |
| AC-F17-11 | Fake configured for each of `Allow`/`Deny`/`Indeterminate`/`RequiresApproval` produces correctly normalized outcomes. |
| AC-F17-12 | Missing/unconfigured digest → `AdapterFailure`; configured failure → `AdapterFailure`. |
| AC-F17-13 | Slice-based fake construction rejects duplicate/malformed digest, invalid conclusion/failure one-of state, and malformed conclusion; constructor-copy ownership is verified. |
| AC-F17-14 | Direct-call linearization tests: pre-call cancel/deadline (zero invocations), exactly one guarded call, post-return cancellation-first tie precedence, late-output discard, and context-honoring adapter liveness. |
| AC-F17-15 | Invalid adapter conclusion (bad outcome / bad reason codes) → `InvalidAdapterResult`, no result. |
| AC-F17-16 | Fixed time source → deterministic `evaluatedAt` and result replay for equivalent inputs. |
| AC-F17-17 | Grammar/duplicate/sort tests plus trusted-fixture review and sentinel leakage tests proving request values and raw adapter errors never enter returned errors/diagnostics. |
| AC-F17-18 | Result has no obligations from the fake. |
| AC-F17-19 | Success transport of result+timing and pure mapping for all four outcomes with exact populated and exact omitted `EvaluationResult` fields and `evaluatedAt == completedAt`. |
| AC-F17-20 | Exact FEATURE-0012 reference and FEATURE-0013 evaluator/input/timing/result structural cases fail with sanitized error; caller input remains byte-for-byte unchanged. |
| AC-F17-21 | Each non-result category produces no mapping input and is never mapped. |
| AC-F17-22 | Package-level assertion of no DecisionRecord/AuditEvent/route/controller/store symbols or effects. |
| AC-F17-23 | `-race` concurrent equivalent-call determinism plus copy-ownership tests mutating caller/fixture/adapter-return/output slices and pointers after handoff. |
| AC-F17-24 | Adapter/import-surface test proving zero network/OPA/Cedar/Kubernetes/CloudProvider/database/identity/provisioning imports or effects. |

## 8. Canonical coverage ledger

Every approved REQ and AC appears exactly once with its design disposition. IDs
are neither created, renumbered, merged, split, reinterpreted, nor omitted.
Disposition legend: `IMPLEMENT` = realized by design components above;
`CONTRACT_ONLY/NO_TASK` = contract statement with no code artifact;
`EXCLUDED` = owned by another feature.

### 8.1 Requirements

| REQ | Disposition | Design realization |
|---|---|---|
| REQ-F17-01 | IMPLEMENT | `request.go` structural validation; strict decode; reference/UID/set rules (§4.1, §5.2). |
| REQ-F17-02 | IMPLEMENT | `canonical.go` v1 normalized object, JCS serialization, SHA-256 digest, absent/empty rules, frozen vector (§2 D-03, §4.2). |
| REQ-F17-03 | IMPLEMENT | `adapter.go`, `result.go`, `timesource.go`, `evaluate.go` conclusion validation, reason-code sorting, timing, non-result categories (§4.3–4.5, §5). |
| REQ-F17-04 | IMPLEMENT | `fake.go` duplicate-preserving fixture shape, immutable digest-keyed fake, strict test-only fixture decoder, and construction validation (§2 D-12/13, §4.7). |
| REQ-F17-05 | IMPLEMENT | `evaluate.go` unnamed result/timing success return; exact `InputIdentity`; approved FEATURE-0013 bool wrapper; pure closed mapping (§4.5, §4.6, §4.8, §5.5). |
| REQ-F17-06 | IMPLEMENT | `evaluate.go` single in-process operation, exact precedence, at-most-once invocation, no route/store/controller (§5.1, §5.3). |

### 8.2 Acceptance cases

| AC | Disposition | Notes |
|---|---|---|
| AC-F17-01 | IMPLEMENT | §7; bounds table test. |
| AC-F17-02 | IMPLEMENT | §4.1 strict decode; §5.2. |
| AC-F17-03 | IMPLEMENT | §5.1 step 2; spy adapter. |
| AC-F17-04 | IMPLEMENT | §4.2 `contextRef` optionality. |
| AC-F17-05 | IMPLEMENT | §4.2, §5.2 pre-digest failure. |
| AC-F17-06 | IMPLEMENT | §2 D-05 duplicate rejection. |
| AC-F17-07 | IMPLEMENT | §2 D-05, §5.3 order-independent digest. |
| AC-F17-08 | IMPLEMENT | §4.2 semantic-change digest. |
| AC-F17-09 | IMPLEMENT | §4.1 `requestId` excluded from normalization. |
| AC-F17-10 | IMPLEMENT | §2 D-03 frozen vector and complete domain-bounded JCS string suite. |
| AC-F17-11 | IMPLEMENT | §4.3 four outcomes; §7. |
| AC-F17-12 | IMPLEMENT | §2 D-12 missing/failure → `AdapterFailure`. |
| AC-F17-13 | IMPLEMENT | §2 D-12 and §4.7 duplicate-preserving constructor/one-of failure. |
| AC-F17-14 | IMPLEMENT | §5.1 precedence; §2 D-11 cardinality/discard. |
| AC-F17-15 | IMPLEMENT | §5.1 step 8 `InvalidAdapterResult`. |
| AC-F17-16 | IMPLEMENT | §2 D-06 fixed time source; §5.3. |
| AC-F17-17 | IMPLEMENT | §2 D-10, §6 grammar/sort plus trusted pre-redaction and leakage proof. |
| AC-F17-18 | IMPLEMENT | §4.3 obligations always absent. |
| AC-F17-19 | IMPLEMENT | §4.5, §4.6 exact populated/omitted fields, timing equality. |
| AC-F17-20 | IMPLEMENT | §4.6/4.8 and §5.5 shared structural-validation failure, no side effect. |
| AC-F17-21 | IMPLEMENT | §4.5 non-result returns zero result and zero timing; §7. |
| AC-F17-22 | IMPLEMENT | §3 dependency direction; §6; §7 package assertion. |
| AC-F17-23 | IMPLEMENT | §5.3 determinism, deep-copy ownership, and `-race`. |
| AC-F17-24 | IMPLEMENT | §3 import surface; §7 zero-external-effect test. |

## 9. Requirement / decision / risk traceability

### Architecture Traceability

- **Consumed handoffs:** ADH-2026-067 (six-group closure), ADH-2026-068
  (structural-only mapper), ADH-2026-069 (success-timing transport, exact
  mapper fields, at-most-once cardinality).
- **Approved design dependency extension:** Sanjeev Kumar, 2026-08-25 — expose
  the existing FEATURE-0013 structural pass through the bool-only
  `IsEvaluationResultStructurallyValid` wrapper; no new validation rule,
  profile/scope input, or Problem output.
- **Decision groups:** CDG-F17-01 → REQ-F17-01; CDG-F17-02 → REQ-F17-02;
  CDG-F17-03 → REQ-F17-03; CDG-F17-04 → REQ-F17-04; CDG-F17-05 → REQ-F17-05;
  CDG-F17-06 → REQ-F17-06.
- **Schema IDs:** `VS0-SCHEMA-018` (`PolicyEvaluationRequest`),
  `VS0-SCHEMA-019` (`PolicyEvaluationResult`) — consumed by exact ID; no field,
  bound, enum, classification, or retention change.
- **Accepted decisions / RFC:** DEC-0026 (FEATURE-0012 foundations reused),
  DEC-0028 and DEC-0036 (adapter boundary; no custom engine), DEC-0043
  (FEATURE-0013 owns `EvaluationResult`/`DecisionProfile`/`DecisionRecord`/
  `AuditEvent`), DEC-0050 (FEATURE-0020 sole `EffectiveGovernanceContext`
  resolver), RFC-0025 (subordinate to the sole architecture authority).
- **Writer IDs:** none owned or activated (`VS0-WRITER-*` zero inventory).
- **State IDs:** none owned or activated (`VS0-STATE-*` zero inventory).
- **Error codes:** none owned, introduced, or activated (no FEATURE-0012
  top-level Problem code, no violation ID); the five `NonResult` categories are
  internal only.
- **Conformance IDs:** none owned (`VS0-CF-*` zero inventory); acceptance
  inventory is `AC-F17-01..24`.

### Reuse dispositions realized

| Capability | Disposition | Design realization |
|---|---|---|
| FEATURE-0012 `TypedRef` / validation / redaction | Reuse | `internal/apimeta`, `internal/apiref` imported unchanged (D-04). |
| FEATURE-0013 `EvaluationResult`/`TimingEnvelope`/`EvaluatorIdentity`/`OpaqueIntegrityCarrier` and exact structural pass | Reuse + approved minimal dependency extension | `internal/decision` plus bool-only `internal/decision/validate.IsEvaluationResultStructurallyValid` imported by the mapper (D-14, §4.8). |
| RFC 8785 JCS + SHA-256 | Reuse (standard) | Stdlib-only v1 serializer + `crypto/sha256` (D-03). |
| `PolicyEngineAdapter` port, boundary, mapper, fake | Build | New `internal/policyeval` package (D-01, D-09, D-12, D-14). |
| Injected UTC time source | Reuse (pattern) | Local `TimeSource` seam, FEATURE-0016 clock as precedent only (D-06). |
| Real OPA/Cedar/other engine | Wrap / Deferred | Not selected or implemented; port kept stable. |

### Risk posture

The sole authority defines no numbered risk register for FEATURE-0017. Its risk
posture is the deferrals and leakage controls of the authority, realized here as
the acyclic inward dependency direction (§3), the zero-external-effect and
no-writer/state/route design, and the negative acceptance cases AC-F17-18,
AC-F17-21, AC-F17-22, and AC-F17-24. No new risk ID is invented.

## 10. Non-goals, absence ledger, and unresolved report

### Non-goals (design realizes their absence)

FEATURE-0017 design defines no authentication/authorization/principal
resolution; no governance/sovereignty/assignment/inheritance/exception/approval/
entitlement/quota/placement/ranking/selection/execution semantics; no
policy/profile publication, taxonomy, bundle distribution, loading, or hot
reload; no `EffectiveGovernanceContext` construction, resolution, or lookup; no
DecisionRecord/AuditEvent publication or any request/result store; no
customer-safe or CloudProvider-safe reason/explanation projection; no public
route, handler, controller, idempotency repository, or persistent storage; no
production adapter selection, registry, routing, retry, or topology; and no real
OPA/Cedar/Casbin/OpenFGA/SpiceDB/Cerbos/Kubernetes/identity/database/network/
CloudProvider SDK/provisioning/plugin execution. Prohibited stale vocabulary
(boolean `allowed`, `EffectivePolicyContext`, `PolicyInput`, `PolicyContext`,
`PolicyBundleRef`, `PolicyDecisionReason`, OPA/Cedar Go placeholders) does not
appear in the design.

### Excluded adjacent-feature mechanisms (not designed here)

- FEATURE-0018: PrincipalRef creation, Membership, RoleDefinition, RoleAssignment,
  authorization resolution, ApprovalRequest, ExceptionGrant.
- FEATURE-0019: SovereigntyProfile/RegulatoryPolicyBundle semantics,
  SovereigntyFactSet, EvidenceRecord, governed-dimension interpretation.
- FEATURE-0020: EffectiveGovernanceContext construction/assignment/inheritance/
  conflict/exception resolution.
- FEATURE-0021: entitlement, quota, provider-selection intent.
- FEATURE-0022: ServicePlacementProfile semantics, service requirements,
  ServiceRegion.
- FEATURE-0023: DecisionProfile adoption compatibility, DecisionRecord
  publication, candidate-set evaluation, per-candidate outcomes, ranking,
  placement, selection.
- FEATURE-0024: production adapter selection, plugin execution, real realization,
  provisioning, execution authority.
- FEATURE-0025: customer-safe/CloudProvider-safe explanation, AI-readable
  projection.
- FEATURE-0026: cross-feature orchestration, Slice 0 integration ownership.

### Absence ledger

| Contract / registry family | FEATURE-0017 inventory | Design consequence |
|---|---|---|
| `VS0-WRITER-*` | None | No writer, persistence owner, or publication path is designed. |
| `VS0-STATE-*` | None | No resource lifecycle or state transition is designed. |
| `VS0-CF-*` | None | Acceptance is `AC-F17-01..24`; no `VS0-CF-*` case is claimed. |
| FEATURE-0012 top-level Problem codes | None | No public Problem response, HTTP mapping, or route error contract. |
| Problem violation IDs / local violation IDs | None | No public or registry-addressable violation identifier. |
| Public routes / handlers / controllers / stores | None | No wire API, controller, or storage is designed. |

### Unresolved report

All design questions delegated by the architecture (package layout, private
carrier/type shapes, interface shapes, constructors, immutable fixture
mechanics, transient success-return representation, fixture serialization, and
goroutine/concurrency mechanics) are resolved above (D-01..D-14). No semantic,
contract, ownership, digest-identity, precedence, result/failure vocabulary, or
field-population question remains open. The full file was re-read against the
sole architecture authority, the three controlling handoffs, and the approved
requirements: it contains no unapproved semantic choice and no mechanism owned
by an excluded feature. `DEPENDENCY_APPROVAL_REQUIRED` was triggered for the
missing FEATURE-0013 structural-only callable and was resolved by Sanjeev
Kumar's explicit 2026-08-25 approval of the exact bool-only wrapper in §4.8.
No stop condition remains active; no architecture, boundary, requirement, or
security approval is outstanding.

STAGE_STATUS: COMPLETE
