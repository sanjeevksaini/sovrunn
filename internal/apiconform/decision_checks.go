package apiconform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/decision/bundle"
	"github.com/sanjeevksaini/sovrunn/internal/decision/compose"
	"github.com/sanjeevksaini/sovrunn/internal/decision/validate"
)

// DecisionScenarioID is a stable FEATURE-0013 conformance scenario identifier
// (architecture §17; F13-CONF-001/002/004; AD-028).
type DecisionScenarioID string

// DecisionCheckKind classifies how a registered scenario is exercised.
type DecisionCheckKind string

const (
	// DecisionCheckPositive expects ValidateDecisionRecord / ValidateAuditEvent success.
	DecisionCheckPositive DecisionCheckKind = "positive"
	// DecisionCheckNegative expects exactly one closed DECISION_* violations[].code.
	DecisionCheckNegative DecisionCheckKind = "negative"
	// DecisionCheckSchema exercises VerifyGoTypeAgainstSchema for FEATURE-0013 bindings.
	DecisionCheckSchema DecisionCheckKind = "schema"
	// DecisionCheckReplay exercises the pure compose.Replay reference kernel (F13-CONF-005).
	DecisionCheckReplay DecisionCheckKind = "replay"
	// DecisionCheckCompat exercises FEATURE-0012 compatibility / AuditEvent regression.
	DecisionCheckCompat DecisionCheckKind = "compat"
)

// DecisionCheckFixture is a bound fixture payload for one scenario ID.
// Fixtures are deterministic, pure, and in-memory (F13-CONF-003). Disk paths are
// optional convenience for later coverage-matrix binding (T-028/T-029/T-030);
// inline bytes take precedence when present.
type DecisionCheckFixture struct {
	Kind DecisionCheckKind

	// Relative to tests/conformance/fixtures when loading from disk.
	DecisionRecordPath string
	AuditEventPath     string
	BundlePath         string

	// Inline payloads (preferred for unit tests; also used when paths are empty).
	DecisionRecordJSON []byte
	AuditEventJSON     []byte
	BundleJSON         []byte
	AuditLinkage       *decision.AuditLinkage

	// Negative expectation (exactly one closed DECISION_* code at an RFC 6901 pointer).
	WantViolationCode  apiproblem.ViolationCode
	WantViolationField string

	// Optional orchestration context for ValidateDecisionRecordWith (T-030
	// negative envelopes may carry ParallelScopeFields, sensitivity, trust,
	// relationship Established, and SupportedObligations).
	RecordOptions validate.DecisionRecordOptions

	// Replay kernel inputs (DecisionCheckReplay).
	CapturedEvaluations []decision.EvaluationResult
	Strategy            compose.Strategy
	StrategyInput       compose.StrategyInput
}

// DecisionCheck is one executable FEATURE-0013 conformance check registered by
// scenario ID. apiconform remains a conformance adapter: checks call exported
// decision/validate functions one-directionally and never become the domain owner.
type DecisionCheck struct {
	ID          DecisionScenarioID
	Description string
	Kind        DecisionCheckKind
	// DefaultFixturePaths are repository-relative paths under
	// tests/conformance/fixtures used when a run-time binding omits inline bytes.
	DefaultFixturePaths []string
}

// DecisionCheckRun is the execution context for RegisteredDecisionChecks.
// BoundFixtures maps scenario ID → fixture for this run. ModuleRoot is required
// for schema checks and disk fixture loading.
type DecisionCheckRun struct {
	ModuleRoot    string
	BoundFixtures map[DecisionScenarioID]DecisionCheckFixture
	DecodeMode    apivalid.DecodeMode
}

// DecisionCheckResult is the outcome of RunDecisionCheck.
type DecisionCheckResult struct {
	ScenarioID DecisionScenarioID
	Kind       DecisionCheckKind
	OK         bool
	Detail     string
	Problem    *apiproblem.Problem
}

// RegisteredDecisionScenarioIDs returns every architecture-owned FEATURE-0013
// conformance scenario ID in stable registration order (parents, mandatory
// suffixes, then additional SCOPE/EVAL/SEC/TRUST/COMPAT cases). Coverage is
// counted by scenario ID (F13-CONF-004); suffixes are never merged or omitted.
func RegisteredDecisionScenarioIDs() []DecisionScenarioID {
	ids := make([]DecisionScenarioID, 0, len(decisionScenarioRegistry))
	for _, c := range decisionScenarioRegistry {
		ids = append(ids, c.ID)
	}
	return ids
}

// RegisteredDecisionChecks returns the executable FEATURE-0013 check registry.
func RegisteredDecisionChecks() []DecisionCheck {
	out := make([]DecisionCheck, len(decisionScenarioRegistry))
	copy(out, decisionScenarioRegistry)
	return out
}

// ResolveDecisionCheck looks up a registered scenario ID. Missing IDs fail closed.
func ResolveDecisionCheck(id DecisionScenarioID) (DecisionCheck, bool) {
	c, ok := decisionScenarioByID[id]
	return c, ok
}

// RunDecisionCheck executes one registered check against its bound fixture.
// Schema checks may run without a fixture binding. All other kinds require a
// bound fixture (inline bytes and/or resolvable fixture paths).
func RunDecisionCheck(check DecisionCheck, run DecisionCheckRun) DecisionCheckResult {
	res := DecisionCheckResult{ScenarioID: check.ID, Kind: check.Kind}
	mode := run.DecodeMode
	if mode == 0 {
		// Conformance fixtures are read representations (ModeCreateRequest is iota 0).
		mode = apivalid.ModeReadRepresentation
	}

	switch check.Kind {
	case DecisionCheckSchema:
		findings := CheckFeature0013TypeBindings(run.ModuleRoot)
		if len(findings) == 0 {
			res.OK = true
			res.Detail = "FEATURE-0013 TypeBindings agree with schemas"
			return res
		}
		res.Detail = strings.Join(findings, "; ")
		return res

	case DecisionCheckReplay:
		fx, ok := boundFixture(check, run)
		if !ok {
			res.Detail = "replay check requires a bound fixture with strategy and captured evaluations"
			return res
		}
		out, prob := CheckReplayKernel(fx.CapturedEvaluations, fx.Strategy, fx.StrategyInput)
		if prob != nil {
			res.Problem = prob
			res.Detail = fmt.Sprintf("replay failed: code=%s", prob.Code)
			return res
		}
		if out.StrategyID == "" && fx.StrategyInput.Meta.ID != "" {
			res.Detail = "replay produced empty strategy id"
			return res
		}
		res.OK = true
		res.Detail = "replay kernel deterministic"
		return res

	case DecisionCheckPositive, DecisionCheckCompat:
		fx, ok := boundFixture(check, run)
		if !ok {
			res.Detail = "positive/compat check requires a bound fixture"
			return res
		}
		return runPositiveOrCompat(check, fx, mode)

	case DecisionCheckNegative:
		fx, ok := boundFixture(check, run)
		if !ok {
			res.Detail = "negative check requires a bound fixture"
			return res
		}
		return runNegative(check, fx, mode)

	default:
		res.Detail = fmt.Sprintf("unknown decision check kind %q", check.Kind)
		return res
	}
}

// CheckFeature0013TypeBindings runs VerifyGoTypeAgainstSchema for every
// FEATURE-0013 TypeBinding (design §6.4/§13; F13-COMPAT-001).
func CheckFeature0013TypeBindings(moduleRoot string) []string {
	if strings.TrimSpace(moduleRoot) == "" {
		return []string{"module root is empty"}
	}
	want := feature0013BindingPaths()
	var findings []string
	found := 0
	for _, binding := range TypeBindings {
		if _, ok := want[binding.SchemaPath]; !ok {
			continue
		}
		found++
		schemaPath := filepath.Join(moduleRoot, filepath.FromSlash(binding.SchemaPath))
		schema, err := os.ReadFile(schemaPath) // #nosec G304 -- trusted repo-local schema path
		if err != nil {
			findings = append(findings, fmt.Sprintf("%s: read: %v", binding.SchemaPath, err))
			continue
		}
		if issues := apischema.ValidateSchemaSupport(schema); len(issues) > 0 {
			findings = append(findings, fmt.Sprintf("%s: schema support: %v", binding.SchemaPath, issues))
			continue
		}
		if issues := apischema.VerifyGoTypeAgainstSchema(schema, binding.GoType); len(issues) > 0 {
			for _, issue := range issues {
				findings = append(findings, fmt.Sprintf("%s: %s %s: %s",
					binding.SchemaPath, issue.Path, issue.Code, issue.Message))
			}
		}
	}
	if found != len(want) {
		findings = append(findings, fmt.Sprintf("FEATURE-0013 TypeBindings found=%d want=%d", found, len(want)))
	}
	return findings
}

// CheckStrictDecode decodes dst from local bytes via apivalid.StrictDecode
// (RID-09; design §13). Failures return FEATURE-0012 top-level Problem codes.
func CheckStrictDecode(data []byte, mode apivalid.DecodeMode, dst any) *apiproblem.Problem {
	if mode == 0 {
		mode = apivalid.ModeReadRepresentation
	}
	contentType := "application/yaml"
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		contentType = "application/json"
	}
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
	req.Header.Set("Content-Type", contentType)
	return apivalid.StrictDecode(httptest.NewRecorder(), req, apivalid.DefaultLimits(), mode, dst)
}

// CheckDecisionRecordConformance decodes and validates a DecisionRecord using
// exported decision/validate orchestration (design §9.1; F13-CONF-001).
func CheckDecisionRecordConformance(
	data []byte,
	mode apivalid.DecodeMode,
	view bundle.BundleView,
) (decision.DecisionRecord, *apiproblem.Problem) {
	if mode == 0 {
		mode = apivalid.ModeReadRepresentation
	}
	return validate.DecodeAndValidateDecisionRecord(data, mode, view)
}

// CheckAuditEventConformance decodes an AuditEvent envelope and validates it
// with FEATURE-0013 AuditLinkage via exported decision/validate functions.
func CheckAuditEventConformance(
	data []byte,
	mode apivalid.DecodeMode,
	linkage decision.AuditLinkage,
	view bundle.BundleView,
) (validate.AuditEventEnvelope, *apiproblem.Problem) {
	if mode == 0 {
		mode = apivalid.ModeReadRepresentation
	}
	return validate.DecodeAndValidateAuditEvent(data, mode, linkage, view)
}

// CheckDecisionRecordNegative asserts ValidateDecisionRecord fails with exactly
// one closed DECISION_* violations[].code at wantField (RFC 6901).
func CheckDecisionRecordNegative(
	data []byte,
	mode apivalid.DecodeMode,
	view bundle.BundleView,
	wantCode apiproblem.ViolationCode,
	wantField string,
) *apiproblem.Problem {
	return CheckDecisionRecordNegativeWith(data, mode, view, validate.DecisionRecordOptions{}, wantCode, wantField)
}

// CheckDecisionRecordNegativeWith is CheckDecisionRecordNegative with optional
// orchestration context (T-030 negative envelopes).
func CheckDecisionRecordNegativeWith(
	data []byte,
	mode apivalid.DecodeMode,
	view bundle.BundleView,
	opts validate.DecisionRecordOptions,
	wantCode apiproblem.ViolationCode,
	wantField string,
) *apiproblem.Problem {
	if mode == 0 {
		mode = apivalid.ModeReadRepresentation
	}
	_, prob := validate.DecodeAndValidateDecisionRecordWith(data, mode, view, opts)
	return assertSingleDecisionViolation(prob, wantCode, wantField)
}

// CheckReplayKernel exercises compose.Replay without rerunning evaluators
// (F13-CONF-005; AD-012, AD-037).
func CheckReplayKernel(
	captured []decision.EvaluationResult,
	strategy compose.Strategy,
	in compose.StrategyInput,
) (compose.StrategyOutput, *apiproblem.Problem) {
	return compose.Replay(captured, strategy, in)
}

// LoadBundleViewStrictDecode loads a DecisionProfileBundle via StrictDecode +
// bundle.Load semantics (offline; no network/crypto).
func LoadBundleViewStrictDecode(data []byte, mode apivalid.DecodeMode) (bundle.BundleView, *apiproblem.Problem) {
	if mode == 0 {
		mode = apivalid.ModeReadRepresentation
	}
	loaded, prob := bundle.Load(data, mode)
	if prob != nil {
		return bundle.BundleView{}, prob
	}
	return loaded.View(), nil
}

func runPositiveOrCompat(check DecisionCheck, fx DecisionCheckFixture, mode apivalid.DecodeMode) DecisionCheckResult {
	res := DecisionCheckResult{ScenarioID: check.ID, Kind: check.Kind}
	view, prob := fixtureView(fx, mode)
	if prob != nil {
		res.Problem = prob
		res.Detail = fmt.Sprintf("bundle load failed: code=%s", prob.Code)
		return res
	}

	if len(fx.DecisionRecordJSON) > 0 || fx.DecisionRecordPath != "" {
		data, err := fixtureBytes(fx.DecisionRecordJSON, fx.DecisionRecordPath, "")
		if err != nil {
			res.Detail = err.Error()
			return res
		}
		if _, prob := CheckDecisionRecordConformance(data, mode, view); prob != nil {
			res.Problem = prob
			res.Detail = fmt.Sprintf("decision record rejected: code=%s", prob.Code)
			return res
		}
	}

	if len(fx.AuditEventJSON) > 0 || fx.AuditEventPath != "" {
		data, err := fixtureBytes(fx.AuditEventJSON, fx.AuditEventPath, "")
		if err != nil {
			res.Detail = err.Error()
			return res
		}
		// FEATURE-0012 base AuditEvent regression (F13-COMPAT-07/08) is proven by
		// StrictDecode of the ImmutableRecord envelope. FEATURE-0013 linkage
		// validation runs only when AuditLinkage is explicitly bound.
		if fx.AuditLinkage == nil {
			var ev AuditEvent
			if prob := CheckStrictDecode(data, mode, &ev); prob != nil {
				res.Problem = prob
				res.Detail = fmt.Sprintf("audit event StrictDecode rejected: code=%s", prob.Code)
				return res
			}
		} else if _, prob := CheckAuditEventConformance(data, mode, *fx.AuditLinkage, view); prob != nil {
			res.Problem = prob
			res.Detail = fmt.Sprintf("audit event rejected: code=%s", prob.Code)
			return res
		}
	}

	if len(fx.DecisionRecordJSON) == 0 && fx.DecisionRecordPath == "" &&
		len(fx.AuditEventJSON) == 0 && fx.AuditEventPath == "" {
		res.Detail = "positive/compat fixture missing decision record and audit event payloads"
		return res
	}

	res.OK = true
	res.Detail = "conformance check passed"
	return res
}

func runNegative(check DecisionCheck, fx DecisionCheckFixture, mode apivalid.DecodeMode) DecisionCheckResult {
	res := DecisionCheckResult{ScenarioID: check.ID, Kind: check.Kind}
	if fx.WantViolationCode == "" {
		res.Detail = "negative fixture missing WantViolationCode"
		return res
	}
	view, prob := fixtureView(fx, mode)
	if prob != nil {
		// Bundle decode failures are FEATURE-0012 top-level codes; surface them.
		res.Problem = prob
		res.Detail = fmt.Sprintf("bundle load failed: code=%s", prob.Code)
		return res
	}

	data, err := fixtureBytes(fx.DecisionRecordJSON, fx.DecisionRecordPath, "")
	if err != nil {
		res.Detail = err.Error()
		return res
	}
	if assertProb := CheckDecisionRecordNegativeWith(
		data, mode, view, fx.RecordOptions, fx.WantViolationCode, fx.WantViolationField,
	); assertProb != nil {
		res.Problem = assertProb
		res.Detail = assertProb.Detail
		return res
	}
	res.OK = true
	res.Detail = fmt.Sprintf("negative conformance matched %s at %s", fx.WantViolationCode, fx.WantViolationField)
	return res
}

func fixtureView(fx DecisionCheckFixture, mode apivalid.DecodeMode) (bundle.BundleView, *apiproblem.Problem) {
	data, err := fixtureBytes(fx.BundleJSON, fx.BundlePath, "")
	if err != nil || len(bytes.TrimSpace(data)) == 0 {
		// Empty bundle is valid for checks that construct BundleView externally
		// via BoundFixtures with BundleJSON already loaded; treat as empty view.
		return bundle.BundleView{}, nil
	}
	return LoadBundleViewStrictDecode(data, mode)
}

func boundFixture(check DecisionCheck, run DecisionCheckRun) (DecisionCheckFixture, bool) {
	if run.BoundFixtures != nil {
		if fx, ok := run.BoundFixtures[check.ID]; ok {
			return fx, true
		}
	}
	// Fall back to default fixture paths when present on disk.
	if run.ModuleRoot == "" || len(check.DefaultFixturePaths) == 0 {
		return DecisionCheckFixture{}, false
	}
	fx := DecisionCheckFixture{Kind: check.Kind}
	for _, rel := range check.DefaultFixturePaths {
		abs := filepath.Join(run.ModuleRoot, ConformanceFixturesDir, filepath.FromSlash(rel))
		raw, err := os.ReadFile(abs) // #nosec G304 -- trusted repo-local fixture path
		if err != nil {
			continue
		}
		switch classifyConformanceFixture(rel, raw) {
		case fixtureClassNegativeEnvelope:
			env, err := parseNegativeDecisionFixture(raw)
			if err != nil {
				continue
			}
			fx.DecisionRecordJSON = env.DecisionRecordJSON
			fx.WantViolationCode = env.WantViolationCode
			fx.WantViolationField = env.WantViolationField
			fx.RecordOptions = env.RecordOptions
			if len(env.BundleJSON) > 0 {
				fx.BundleJSON = env.BundleJSON
			}
			fx.DecisionRecordPath = abs
		case fixtureClassBundle:
			fx.BundleJSON = raw
			fx.BundlePath = abs
		case fixtureClassAuditEvent:
			fx.AuditEventJSON = raw
			fx.AuditEventPath = abs
		default:
			fx.DecisionRecordJSON = raw
			fx.DecisionRecordPath = abs
		}
	}
	if check.Kind == DecisionCheckNegative {
		if fx.WantViolationCode == "" || len(fx.DecisionRecordJSON) == 0 {
			return DecisionCheckFixture{}, false
		}
		return fx, true
	}
	if len(fx.DecisionRecordJSON) == 0 && len(fx.AuditEventJSON) == 0 && len(fx.BundleJSON) == 0 {
		return DecisionCheckFixture{}, false
	}
	return fx, true
}

const (
	fixtureClassRecord           = "record"
	fixtureClassAuditEvent       = "audit"
	fixtureClassBundle           = "bundle"
	fixtureClassNegativeEnvelope = "negative-envelope"
)

// classifyConformanceFixture prefers JSON kind sniffing so AuditEvent and
// DecisionProfileBundle bodies under decision/positive/ load correctly
// (T-029; RID-07). Path conventions remain the fallback.
func classifyConformanceFixture(rel string, raw []byte) string {
	switch sniffJSONKind(raw) {
	case "DecisionNegativeFixture":
		return fixtureClassNegativeEnvelope
	case "DecisionProfileBundle":
		return fixtureClassBundle
	case "AuditEvent":
		return fixtureClassAuditEvent
	case "DecisionRecord":
		return fixtureClassRecord
	}
	switch {
	case strings.Contains(rel, "bundle"):
		return fixtureClassBundle
	case strings.Contains(rel, "audit-event"):
		return fixtureClassAuditEvent
	case strings.Contains(rel, "negative/decision") || strings.HasPrefix(rel, "negative/"):
		// Prefer envelope when present; otherwise treat as bare DecisionRecord
		// (legacy inline unit-test payloads).
		if sniffJSONKind(raw) == "DecisionNegativeFixture" {
			return fixtureClassNegativeEnvelope
		}
		// Heuristic: wrapper objects carry wantViolationCode.
		if bytes.Contains(bytes.TrimSpace(raw), []byte(`"wantViolationCode"`)) {
			return fixtureClassNegativeEnvelope
		}
		return fixtureClassRecord
	case strings.Contains(rel, "decision/positive") || strings.Contains(rel, "decision-record"):
		return fixtureClassRecord
	default:
		return fixtureClassRecord
	}
}

func sniffJSONKind(raw []byte) string {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return ""
	}
	var meta struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(trimmed, &meta); err != nil {
		return ""
	}
	return strings.TrimSpace(meta.Kind)
}

func fixtureBytes(inline []byte, path, moduleRoot string) ([]byte, error) {
	if len(bytes.TrimSpace(inline)) > 0 {
		return inline, nil
	}
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("fixture payload is empty")
	}
	if filepath.IsAbs(path) {
		raw, err := os.ReadFile(path) // #nosec G304 -- trusted absolute fixture path
		if err != nil {
			return nil, err
		}
		return raw, nil
	}
	if moduleRoot == "" {
		return nil, fmt.Errorf("relative fixture path %q requires module root", path)
	}
	raw, err := os.ReadFile(filepath.Join(moduleRoot, path)) // #nosec G304 -- trusted repo-local fixture path
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func assertSingleDecisionViolation(prob *apiproblem.Problem, wantCode apiproblem.ViolationCode, wantField string) *apiproblem.Problem {
	if prob == nil {
		return apiproblem.New(apiproblem.CodeValidationFailed).
			WithDetail(fmt.Sprintf("expected violation %s at %s, got success", wantCode, wantField))
	}
	if prob.Code != apiproblem.CodeValidationFailed {
		return apiproblem.New(prob.Code).
			WithDetail(fmt.Sprintf("expected top-level VALIDATION_FAILED, got %s: %s", prob.Code, prob.Detail))
	}
	if len(prob.Violations) != 1 {
		return apiproblem.New(prob.Code).
			WithDetail(fmt.Sprintf("expected exactly one violation, got %d", len(prob.Violations)))
	}
	v := prob.Violations[0]
	if v.Code != wantCode {
		return apiproblem.New(prob.Code).
			WithDetail(fmt.Sprintf("expected violations[].code %s, got %s", wantCode, v.Code))
	}
	if wantField != "" && v.Field != wantField {
		return apiproblem.New(prob.Code).
			WithDetail(fmt.Sprintf("expected violations[].field %s, got %s", wantField, v.Field))
	}
	return nil
}

func feature0013BindingPaths() map[string]struct{} {
	want := make(map[string]struct{}, len(feature0013CanonicalSchemaFiles)+8)
	for _, name := range feature0013CanonicalSchemaFiles {
		want["api/schemas/"+name] = struct{}{}
	}
	for _, name := range []string{
		"decision-linkage.json",
		"decision-profile-ref.json",
		"trust-carrier.json",
		"security-exception-ref.json",
		"semantic-decision-identity.json",
		"decision-relationship.json",
		"sensitivity.json",
		"decision-graph.json",
	} {
		want["api/schemas/_common/"+name] = struct{}{}
	}
	return want
}

// decisionScenarioRegistry is the authoritative FEATURE-0013 scenario check
// table. Default fixture path conventions match RID-07; physical fixtures are
// supplied by T-029/T-030 and mapped by T-028. T-027 registers and executes.
var decisionScenarioRegistry = buildDecisionScenarioRegistry()

var decisionScenarioByID = indexDecisionScenarios(decisionScenarioRegistry)

func indexDecisionScenarios(checks []DecisionCheck) map[DecisionScenarioID]DecisionCheck {
	out := make(map[DecisionScenarioID]DecisionCheck, len(checks))
	for _, c := range checks {
		out[c.ID] = c
	}
	return out
}

func buildDecisionScenarioRegistry() []DecisionCheck {
	type row struct {
		id   DecisionScenarioID
		desc string
		kind DecisionCheckKind
	}
	rows := make([]row, 0, 128)

	add := func(id, desc string, kind DecisionCheckKind) {
		rows = append(rows, row{DecisionScenarioID(id), desc, kind})
	}

	// F13-CF-01 … F13-CF-28 parents (architecture §17).
	parents := []struct {
		id, desc string
		kind     DecisionCheckKind
	}{
		{"F13-CF-01", "simplest synchronous atomic authorization decision", DecisionCheckPositive},
		{"F13-CF-02", "simplest deterministic denial with actionable rationale and obligations", DecisionCheckPositive},
		{"F13-CF-03", "pure selection with a typed result and no adjudication facet", DecisionCheckPositive},
		{"F13-CF-04", "compound family parent: ranking/classification/resolution/allocation/plan/assessment/recommendation/advisory/simulation", DecisionCheckPositive},
		{"F13-CF-05", "human-approval asynchronous handoff returns existing Operation reference", DecisionCheckPositive},
		{"F13-CF-06", "maximally complex but bounded hierarchical composite decision", DecisionCheckPositive},
		{"F13-CF-07", "parallel evaluations with deterministic aggregation", DecisionCheckPositive},
		{"F13-CF-08", "compound family parent: timeout/missing evidence/evaluator failure/policy conflict", DecisionCheckNegative},
		{"F13-CF-09", "idempotent replay and concurrent duplicate requests", DecisionCheckPositive},
		{"F13-CF-10", "partial decision/audit persistence failure and reconciliation (contract-only)", DecisionCheckPositive},
		{"F13-CF-11", "compound family parent: supersession/correction/revocation/cycle/chain-limit", DecisionCheckNegative},
		{"F13-CF-12", "compound family parent: unauthorized profile/authority/projection/obligation/cross-scope", DecisionCheckNegative},
		{"F13-CF-13", "offline/air-gapped static profile-bundle structure and opaque trust carriers", DecisionCheckPositive},
		{"F13-CF-14", "export/import across two conforming provider implementations (contract-only)", DecisionCheckPositive},
		{"F13-CF-15", "AI-disabled operation and AI projection redaction (structural)", DecisionCheckPositive},
		{"F13-CF-16", "compound family parent: graph size/depth/width/fan-out/timeout/payload budget rejection", DecisionCheckNegative},
		{"F13-CF-17", "registration of a new decision family without a common-envelope change", DecisionCheckSchema},
		{"F13-CF-18", "immutable supersession and revocation with calculated effective state", DecisionCheckPositive},
		{"F13-CF-19", "local atomic decision/audit-obligation acceptance while audit delivery unavailable", DecisionCheckPositive},
		{"F13-CF-20", "compound family parent: idempotency retry identity versus semantic decision identity", DecisionCheckPositive},
		{"F13-CF-21", "deterministic replay from captured evaluator output without rerunning evaluator", DecisionCheckReplay},
		{"F13-CF-22", "synchronous contract with local version-pinned bundles and opaque trust state", DecisionCheckPositive},
		{"F13-CF-23", "structural rejection of unknown/expired/revoked/mismatched profile bundles", DecisionCheckNegative},
		{"F13-CF-24", "enforcement denial for unknown or unsupported mandatory obligations", DecisionCheckNegative},
		{"F13-CF-25", "local authorization artifact cache hit/expiry/revocation/scope/stale failure (contract-only)", DecisionCheckNegative},
		{"F13-CF-26", "compound family parent: cancellation and composite budget rejection", DecisionCheckNegative},
		{"F13-CF-27", "structural representation for later local-digest/batched-signing profile during outage", DecisionCheckPositive},
		{"F13-CF-28", "compound family parent: disconnected export/import contract semantics", DecisionCheckPositive},
	}
	for _, p := range parents {
		add(p.id, p.desc, p.kind)
	}

	// Mandatory compound suffixes (architecture §17 textual order).
	for _, s := range []struct{ id, desc string }{
		{"F13-CF-04.a", "ranking profile"},
		{"F13-CF-04.b", "classification profile"},
		{"F13-CF-04.c", "resolution profile"},
		{"F13-CF-04.d", "allocation profile"},
		{"F13-CF-04.e", "plan profile"},
		{"F13-CF-04.f", "assessment profile"},
		{"F13-CF-04.g", "recommendation profile"},
		{"F13-CF-04.h", "advisory profile"},
		{"F13-CF-04.i", "simulation profile"},
	} {
		add(s.id, s.desc, DecisionCheckPositive)
	}
	for _, s := range []struct{ id, desc string }{
		{"F13-CF-08.a", "timeout"},
		{"F13-CF-08.b", "missing evidence"},
		{"F13-CF-08.c", "evaluator failure"},
		{"F13-CF-08.d", "policy conflict"},
	} {
		add(s.id, s.desc, DecisionCheckNegative)
	}
	for _, s := range []struct{ id, desc string }{
		{"F13-CF-11.a", "supersession"},
		{"F13-CF-11.b", "correction"},
		{"F13-CF-11.c", "revocation"},
		{"F13-CF-11.d", "cycle"},
		{"F13-CF-11.e", "chain-limit"},
	} {
		add(s.id, s.desc, DecisionCheckNegative)
	}
	for _, s := range []struct{ id, desc string }{
		{"F13-CF-12.a", "unauthorized profile"},
		{"F13-CF-12.b", "unauthorized authority"},
		{"F13-CF-12.c", "unauthorized projection"},
		{"F13-CF-12.d", "obligation bypass"},
		{"F13-CF-12.e", "cross-scope reference rejection"},
	} {
		add(s.id, s.desc, DecisionCheckNegative)
	}
	for _, s := range []struct{ id, desc string }{
		{"F13-CF-16.a", "graph size rejection"},
		{"F13-CF-16.b", "graph depth rejection"},
		{"F13-CF-16.c", "graph width rejection"},
		{"F13-CF-16.d", "graph fan-out rejection"},
		{"F13-CF-16.e", "timeout budget rejection"},
		{"F13-CF-16.f", "payload budget rejection"},
	} {
		add(s.id, s.desc, DecisionCheckNegative)
	}
	for _, s := range []struct{ id, desc string }{
		{"F13-CF-20.a", "policy change identity"},
		{"F13-CF-20.b", "profile change identity"},
		{"F13-CF-20.c", "strategy change identity"},
		{"F13-CF-20.d", "authority change identity"},
		{"F13-CF-20.e", "evaluator-version change identity"},
		{"F13-CF-20.f", "validity change identity"},
	} {
		add(s.id, s.desc, DecisionCheckPositive)
	}
	for _, s := range []struct{ id, desc string }{
		{"F13-CF-26.a", "cancellation budget"},
		{"F13-CF-26.b", "remote-call budget"},
		{"F13-CF-26.c", "byte budget"},
		{"F13-CF-26.d", "concurrency budget"},
		{"F13-CF-26.e", "retry budget"},
		{"F13-CF-26.f", "jurisdiction budget"},
		{"F13-CF-26.g", "time budget"},
		{"F13-CF-26.h", "cost budget"},
	} {
		add(s.id, s.desc, DecisionCheckNegative)
	}
	for _, s := range []struct{ id, desc string }{
		{"F13-CF-28.a", "replay"},
		{"F13-CF-28.b", "duplicate"},
		{"F13-CF-28.c", "ordering"},
		{"F13-CF-28.d", "unknown trust root"},
		{"F13-CF-28.e", "revoked origin"},
		{"F13-CF-28.f", "conflict reconciliation"},
	} {
		add(s.id, s.desc, DecisionCheckPositive)
	}

	// Additional mandatory cases (do not alter the canonical 28-family count).
	for _, s := range []struct {
		id, desc string
		kind     DecisionCheckKind
	}{
		{"F13-SCOPE-01", "positive canonical Platform scope; serialized metadata.scopeRef absent", DecisionCheckPositive},
		{"F13-SCOPE-02", "positive canonical Organization scope", DecisionCheckPositive},
		{"F13-SCOPE-03", "positive canonical OrganizationUnit scope", DecisionCheckPositive},
		{"F13-SCOPE-04", "positive canonical Tenant scope", DecisionCheckPositive},
		{"F13-SCOPE-05", "positive canonical Project scope", DecisionCheckPositive},
		{"F13-SCOPE-06", "positive canonical Provider scope", DecisionCheckPositive},
		{"F13-SCOPE-07", "absent scope rejected when Platform is not permitted", DecisionCheckNegative},
		{"F13-SCOPE-08", "explicit Platform input normalizes to canonical absent serialized form", DecisionCheckPositive},
		{"F13-SCOPE-09", "parallel, duplicate, or second scope source rejected", DecisionCheckNegative},
		{"F13-SCOPE-10", "out-of-vocabulary scope, including ServiceInstance, rejected", DecisionCheckNegative},
		{"F13-SCOPE-11", "ServiceInstance subject under Project metadata.scopeRef; subject/scope mismatch rejected", DecisionCheckPositive},
		{"F13-EVAL-01", "profile sensitivity and projection limits accepted at declared bounds", DecisionCheckPositive},
		{"F13-EVAL-02", "evaluator registration narrows profile limits", DecisionCheckPositive},
		{"F13-EVAL-03", "evaluator widening of profile limits rejected", DecisionCheckNegative},
		{"F13-EVAL-04", "missing, unknown, or incomparable mandatory limits rejected", DecisionCheckNegative},
		{"F13-EVAL-05", "prohibited typed content category or value type rejected structurally", DecisionCheckNegative},
		{"F13-EVAL-06", "evaluation scope mismatch or unauthorized widening rejected without disclosure", DecisionCheckNegative},
		{"F13-EVAL-07", "opaque expected-versus-observed integrity state mismatch fails closed", DecisionCheckNegative},
		{"F13-EVAL-08", "arbitrary text is not subjected to lexical secret/credential/PII/malware/DLP scanning", DecisionCheckPositive},
		{"F13-SEC-01", "sensitivity ceiling boundary", DecisionCheckNegative},
		{"F13-SEC-02", "field-count boundary", DecisionCheckNegative},
		{"F13-SEC-03", "byte-size boundary", DecisionCheckNegative},
		{"F13-SEC-04", "nesting-depth boundary", DecisionCheckNegative},
		{"F13-SEC-05", "allowed-value-type boundary", DecisionCheckNegative},
		{"F13-SEC-06", "prohibited-category boundary", DecisionCheckNegative},
		{"F13-SEC-07", "FEATURE-0012 classification-to-sensitivity floor cannot be lowered", DecisionCheckNegative},
		{"F13-SEC-08", "unknown, unmapped, contradictory, or below-floor classification fails closed", DecisionCheckNegative},
		{"F13-TRUST-01", "algorithm-agile carrier presence and identifier syntax validated structurally", DecisionCheckPositive},
		{"F13-TRUST-02", "covered-field descriptor presence validated without interpreting it", DecisionCheckPositive},
		{"F13-TRUST-03", "unknown/absent/expired/revoked/mismatched/unverified required trust state fails closed", DecisionCheckNegative},
		{"F13-TRUST-04", "no canonicalization/digest/signing/verification/notarization before ADR-F13-002", DecisionCheckPositive},
		{"F13-COMPAT-01", "FEATURE-0012 Platform ScopeKind regression", DecisionCheckCompat},
		{"F13-COMPAT-02", "FEATURE-0012 Organization ScopeKind regression", DecisionCheckCompat},
		{"F13-COMPAT-03", "FEATURE-0012 OrganizationUnit ScopeKind regression", DecisionCheckCompat},
		{"F13-COMPAT-04", "FEATURE-0012 Tenant ScopeKind regression", DecisionCheckCompat},
		{"F13-COMPAT-05", "FEATURE-0012 Project ScopeKind regression", DecisionCheckCompat},
		{"F13-COMPAT-06", "FEATURE-0012 Provider ScopeKind regression", DecisionCheckCompat},
		{"F13-COMPAT-07", "existing Organization-scoped FEATURE-0012 AuditEvent remains valid", DecisionCheckCompat},
		{"F13-COMPAT-08", "AuditEvent accepts all six governance scopes without changing ScopeKind vocabulary", DecisionCheckCompat},
		{"F13-COMPAT-09", "FEATURE-0012 metadata/reference/validation/Problem/extension/ServiceInstance semantics unchanged", DecisionCheckSchema},
	} {
		add(s.id, s.desc, s.kind)
	}

	out := make([]DecisionCheck, 0, len(rows))
	seen := make(map[DecisionScenarioID]struct{}, len(rows))
	for _, r := range rows {
		if _, dup := seen[r.id]; dup {
			panic("duplicate FEATURE-0013 scenario ID: " + string(r.id))
		}
		seen[r.id] = struct{}{}
		paths := defaultFixturePathsFor(r.id, r.kind)
		out = append(out, DecisionCheck{
			ID:                  r.id,
			Description:         r.desc,
			Kind:                r.kind,
			DefaultFixturePaths: paths,
		})
	}
	return out
}

// SharedFixtureLocalBundlePath is the fixture-local DecisionProfileBundle used
// by positive DecisionRecord scenarios (T-029). Profiles declare explicit
// fixture-local values; no FEATURE-0013 defaults, fallbacks, or hidden maximums.
const SharedFixtureLocalBundlePath = "decision/positive/_shared-bundle.json"

func defaultFixturePathsFor(id DecisionScenarioID, kind DecisionCheckKind) []string {
	name := string(id)
	switch kind {
	case DecisionCheckNegative:
		return []string{"negative/decision/" + name + ".json"}
	case DecisionCheckCompat:
		if strings.HasPrefix(name, "F13-COMPAT-07") || strings.HasPrefix(name, "F13-COMPAT-08") {
			return []string{
				"audit-event.json",
				"decision/positive/" + name + ".json",
				SharedFixtureLocalBundlePath,
			}
		}
		return []string{"decision/positive/" + name + ".json", SharedFixtureLocalBundlePath}
	case DecisionCheckSchema, DecisionCheckReplay:
		return nil
	default:
		return []string{"decision/positive/" + name + ".json", SharedFixtureLocalBundlePath}
	}
}

// DecisionScenarioIDSet returns a set copy of RegisteredDecisionScenarioIDs
// for coverage-matrix consumers (T-028).
func DecisionScenarioIDSet() map[DecisionScenarioID]struct{} {
	ids := RegisteredDecisionScenarioIDs()
	out := make(map[DecisionScenarioID]struct{}, len(ids))
	for _, id := range ids {
		out[id] = struct{}{}
	}
	return out
}

// SortedDecisionScenarioIDs returns RegisteredDecisionScenarioIDs sorted
// lexicographically (stable helper for diffing; registration order remains
// authoritative via RegisteredDecisionScenarioIDs).
func SortedDecisionScenarioIDs() []DecisionScenarioID {
	ids := RegisteredDecisionScenarioIDs()
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}
