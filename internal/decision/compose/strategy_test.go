package compose

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

func sampleException() decision.SecurityExceptionRef {
	return decision.SecurityExceptionRef{
		ApprovalRef: apimeta.TypedRef{
			APIVersion: "governance.sovrunn.io/v1alpha1",
			Kind:       "SecurityExceptionApproval",
			Name:       "break-glass-1",
			UID:        "exception-approval-uid-1",
		},
		OwnerRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       "Organization",
			Name:       "acme",
			UID:        "org-uid-1",
		},
		ApprovingAuthorityRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       "OrganizationUnit",
			Name:       "security-board",
			UID:        "ou-uid-1",
		},
		CompensatingControls: []string{"ctrl:dual-control"},
		EffectiveFrom:        "2026-07-01T00:00:00Z",
		EffectiveUntil:       "2026-07-31T23:59:59Z",
		AuditTreatment:       "audit:mandatory-durable",
		ReassessmentTrigger:  "trigger:quarterly-review",
		CoveredFailureModes:  []string{"timeout", "error"},
		Purpose:              "bounded break-glass under dual control",
	}
}

func sampleMeta(requestsFailOpen bool) StrategyMetadata {
	return StrategyMetadata{
		ID:                  "all-of",
		Version:             "1.0.0",
		AcceptedInputTypes:  []string{"evaluation-result"},
		Ordering:            "declaration",
		ShortCircuit:        true,
		MissingBehavior:     "fail-closed",
		TimeoutBehavior:     "fail-closed",
		ConflictBehavior:    "fail-closed",
		ErrorBehavior:       "fail-closed",
		RequestsFailOpen:    requestsFailOpen,
		DeterministicOutput: true,
		MaxFanOut:           8,
		ExplanationRules:    "record-all-inputs",
		EvidenceRules:       "retain-captured-evaluations",
	}
}

func TestStrategyPurityNoIOImports(t *testing.T) {
	t.Parallel()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse dir: %v", err)
	}

	bannedPrefixes := []string{
		"net",
		"os",
		"database",
		"crypto",
		"plugin",
		"syscall",
		"github.com/sanjeevksaini/sovrunn/internal/server",
		"github.com/sanjeevksaini/sovrunn/internal/api",
	}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			for _, imp := range f.Imports {
				path := strings.Trim(imp.Path.Value, `"`)
				for _, banned := range bannedPrefixes {
					if path == banned || strings.HasPrefix(path, banned+"/") {
						t.Fatalf("compose package must remain pure; banned import %q in %s", path, fset.File(f.Pos()).Name())
					}
				}
			}
		}
	}
}

func TestStrategyFuncIsPureInMemory(t *testing.T) {
	t.Parallel()

	// A Strategy that only transforms captured evaluations — no I/O.
	var pure Strategy = func(in StrategyInput) (StrategyOutput, *apiproblem.Problem) {
		out := StrategyOutput{
			StrategyID:      in.Meta.ID,
			StrategyVersion: in.Meta.Version,
			InputRefs:       append([]string(nil), in.InputRefs...),
		}
		for _, ev := range in.Evaluations {
			if ev.ResultStatus != decision.EvaluationResultStatusSuccess {
				return ApplyFailure(in, string(ev.ResultStatus))
			}
			if len(ev.Result) > 0 {
				out.TypedResult = append(json.RawMessage(nil), ev.Result...)
			}
		}
		allowed := decision.AdjudicationOutcomeAllowed
		out.Adjudication = &allowed
		out.Rationale = decision.DecisionRationale{
			ReasonCodes: []string{"ALL_SUCCESS"},
			Reasons:     []string{"all captured evaluations succeeded"},
		}
		return out, nil
	}

	in := StrategyInput{
		Meta: sampleMeta(false),
		FailureBehavior: decision.FailureBehavior{
			FailPosture: decision.FailPostureFailClosed,
		},
		Evaluations: []decision.EvaluationResult{
			{
				Evaluator:    decision.EvaluatorIdentity{Type: "capacity", Version: "1.0.0"},
				ResultStatus: decision.EvaluationResultStatusSuccess,
				Result:       json.RawMessage(`{"ok":true}`),
				Timing: decision.TimingEnvelope{
					StartedAt:   "2026-07-29T10:00:00Z",
					CompletedAt: "2026-07-29T10:00:01Z",
				},
				EvaluatedAt: "2026-07-29T10:00:01Z",
			},
		},
		InputRefs: []string{"eval-1"},
	}

	out, prob := pure(in)
	if prob != nil {
		t.Fatalf("pure strategy problem = %#v", prob)
	}
	if out.StrategyID != "all-of" || out.StrategyVersion != "1.0.0" {
		t.Fatalf("strategy identity = %s/%s", out.StrategyID, out.StrategyVersion)
	}
	if out.FailOpenApplied {
		t.Fatal("fail-open must not apply on success path")
	}
	if out.Adjudication == nil || *out.Adjudication != decision.AdjudicationOutcomeAllowed {
		t.Fatalf("adjudication = %#v", out.Adjudication)
	}
	if string(out.TypedResult) != `{"ok":true}` {
		t.Fatalf("typedResult = %s", out.TypedResult)
	}
}

func TestFailOpenClosedDefault(t *testing.T) {
	t.Parallel()

	meta := sampleMeta(false)
	fb := decision.FailureBehavior{FailPosture: decision.FailPostureFailClosed}
	if FailOpenEffective(meta, fb) {
		t.Fatal("default must be fail-closed when strategy does not request fail-open")
	}

	// Strategy requests fail-open but profile remains fail-closed → still closed.
	meta.RequestsFailOpen = true
	if FailOpenEffective(meta, fb) {
		t.Fatal("RequestsFailOpen alone must not authorize fail-open")
	}

	// Profile declares FAIL_OPEN without exception evidence → still closed.
	fb.FailPosture = decision.FailPostureFailOpen
	fb.SecurityExceptionRef = nil
	if FailOpenEffective(meta, fb) {
		t.Fatal("FAIL_OPEN without SecurityExceptionRef must fail closed")
	}

	in := StrategyInput{Meta: meta, FailureBehavior: fb, InputRefs: []string{"e1"}}
	_, prob := ApplyFailure(in, "timeout")
	if prob == nil {
		t.Fatal("ApplyFailure must fail closed without valid exception evidence")
	}
	if prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("top-level code = %q", prob.Code)
	}
	if len(prob.Violations) != 1 || prob.Violations[0].Code != violationInputConflict {
		t.Fatalf("violations = %#v", prob.Violations)
	}
}

func TestFailOpenEffectiveOnlyWithValidExceptionEvidence(t *testing.T) {
	t.Parallel()

	meta := sampleMeta(true)
	exc := sampleException()
	fb := decision.FailureBehavior{
		FailPosture:          decision.FailPostureFailOpen,
		SecurityExceptionRef: &exc,
	}
	if !FailOpenEffective(meta, fb) {
		t.Fatal("fail-open must be effective with RequestsFailOpen + FAIL_OPEN + valid SecurityExceptionRef")
	}

	in := StrategyInput{
		Meta:            meta,
		FailureBehavior: fb,
		InputRefs:       []string{"eval-timeout"},
		Evaluations: []decision.EvaluationResult{{
			Evaluator:    decision.EvaluatorIdentity{Type: "capacity", Version: "1.0.0"},
			ResultStatus: decision.EvaluationResultStatusTimeout,
			Timing: decision.TimingEnvelope{
				StartedAt:   "2026-07-29T10:00:00Z",
				CompletedAt: "2026-07-29T10:00:05Z",
			},
			EvaluatedAt: "2026-07-29T10:00:05Z",
		}},
	}
	out, prob := ApplyFailure(in, "timeout")
	if prob != nil {
		t.Fatalf("ApplyFailure with valid exception = %#v", prob)
	}
	if !out.FailOpenApplied {
		t.Fatal("FailOpenApplied must be true when gate is effective")
	}
	if out.StrategyID != "all-of" || out.StrategyVersion != "1.0.0" {
		t.Fatalf("identity = %s/%s", out.StrategyID, out.StrategyVersion)
	}
}

func TestFailOpenRejectedForMalformedExceptionEvidence(t *testing.T) {
	t.Parallel()

	meta := sampleMeta(true)

	cases := []struct {
		name string
		mut  func(*decision.SecurityExceptionRef)
	}{
		{"missing approval uid", func(r *decision.SecurityExceptionRef) { r.ApprovalRef.UID = "" }},
		{"missing owner uid", func(r *decision.SecurityExceptionRef) { r.OwnerRef.UID = "" }},
		{"empty compensating controls", func(r *decision.SecurityExceptionRef) { r.CompensatingControls = nil }},
		{"empty covered modes", func(r *decision.SecurityExceptionRef) { r.CoveredFailureModes = []string{} }},
		{"empty purpose", func(r *decision.SecurityExceptionRef) { r.Purpose = "" }},
		{"non-positive interval", func(r *decision.SecurityExceptionRef) {
			r.EffectiveFrom = "2026-07-31T00:00:00Z"
			r.EffectiveUntil = "2026-07-01T00:00:00Z"
		}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			exc := sampleException()
			tc.mut(&exc)
			fb := decision.FailureBehavior{
				FailPosture:          decision.FailPostureFailOpen,
				SecurityExceptionRef: &exc,
			}
			if FailOpenEffective(meta, fb) {
				t.Fatalf("malformed evidence %q must fail closed", tc.name)
			}
		})
	}
}

func TestStrategyMetadataCarriesNoApprovalAuthority(t *testing.T) {
	t.Parallel()

	meta := sampleMeta(true)
	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	raw := string(data)
	for _, banned := range []string{
		`"securityExceptionRef"`,
		`"approvalRef"`,
		`"approvalWorkflow"`,
		`"issuance"`,
		`"revocation"`,
		`"failPosture"`,
		`"exceptionService"`,
	} {
		if strings.Contains(raw, banned) {
			t.Fatalf("strategy metadata must not carry approval authority field %s; json=%s", banned, raw)
		}
	}
	if !strings.Contains(raw, `"requestsFailOpen":true`) {
		t.Fatalf("requestsFailOpen descriptive flag missing; json=%s", raw)
	}
	// Round-trip preserves registered declaration fields.
	var out StrategyMetadata
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.ID != meta.ID || out.Version != meta.Version || !out.RequestsFailOpen {
		t.Fatalf("round-trip = %+v", out)
	}
}

func TestFailOpenRequiresProfileFailOpenPosture(t *testing.T) {
	t.Parallel()

	meta := sampleMeta(true)
	exc := sampleException()

	for _, posture := range []decision.FailPosture{
		decision.FailPostureFailClosed,
		decision.FailPostureRequiresApproval,
	} {
		fb := decision.FailureBehavior{
			FailPosture:          posture,
			SecurityExceptionRef: &exc,
		}
		if FailOpenEffective(meta, fb) {
			t.Fatalf("posture %q must not make strategy fail-open effective", posture)
		}
	}
}
