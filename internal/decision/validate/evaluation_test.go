package validate

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

func validEvaluationProfile() decision.DecisionProfile {
	prof := validProfile()
	prof.Spec.AcceptedEvaluationTypes = []decision.AcceptedEvaluationType{
		{Type: "capacity-evaluator", Version: "1.2.0"},
	}
	prof.Spec.Limits.MaxLatencyMs = 2000
	prof.Spec.Limits.MaxPayloadBytes = 65536
	return prof
}

func projectRecScope(uid string) apimeta.ScopeIdentity {
	return apimeta.ScopeIdentity{Kind: apimeta.ScopeProject, UID: uid}
}

func validEvaluation() decision.EvaluationResult {
	return decision.EvaluationResult{
		Evaluator: decision.EvaluatorIdentity{
			Type:    "capacity-evaluator",
			Version: "1.2.0",
		},
		InputSnapshotRef: &apimeta.TypedRef{
			APIVersion: "governance.sovrunn.io/v1alpha1",
			Kind:       "InputSnapshot",
			Name:       "placement-input-1",
			UID:        "input-uid-1",
		},
		ResultStatus: decision.EvaluationResultStatusSuccess,
		Result:       json.RawMessage(`{"fit":true}`),
		Timing: decision.TimingEnvelope{
			StartedAt:   "2026-07-29T10:00:00Z",
			CompletedAt: "2026-07-29T10:00:00.050Z",
			DurationMs:  50,
		},
		EvaluatedAt: "2026-07-29T10:00:00.050Z",
	}
}

func assertEvaluationViolation(t *testing.T, p *apiproblem.Problem, wantCode apiproblem.ViolationCode, wantField string) {
	t.Helper()
	if p == nil {
		t.Fatalf("expected Problem with %s at %s, got nil", wantCode, wantField)
	}
	if p.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("top-level code = %q, want %q", p.Code, apiproblem.CodeValidationFailed)
	}
	if p.Status != 422 {
		t.Fatalf("status = %d, want 422", p.Status)
	}
	if p.Type != "urn:sovrunn:problem:validation-failed" {
		t.Fatalf("type = %q, want urn:sovrunn:problem:validation-failed", p.Type)
	}
	if len(p.Violations) != 1 {
		t.Fatalf("violations len=%d, want 1: %#v", len(p.Violations), p.Violations)
	}
	v := p.Violations[0]
	if v.Code != wantCode {
		t.Fatalf("violations[0].code = %q, want %q", v.Code, wantCode)
	}
	if v.Field != wantField {
		t.Fatalf("violations[0].field = %q, want %q", v.Field, wantField)
	}
	if v.Message == "" {
		t.Fatal("violations[0].message must be non-empty (redactable)")
	}
}

func TestValidateEvaluation_HappyPath(t *testing.T) {
	t.Parallel()

	p := ValidateEvaluation(validEvaluation(), projectRecScope("project-123"), validEvaluationProfile())
	if p != nil {
		t.Fatalf("valid evaluation must pass, got %#v", p)
	}
}

func TestValidateEvaluation_HappyPathEqualScopeEvidence(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.ScopeEvidence = &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeProject),
			Name:       "payments-production",
			UID:        "project-123",
		},
	}
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	if p != nil {
		t.Fatalf("equal scopeEvidence must pass, got %#v", p)
	}
}

func TestValidateEvaluation_HappyPathNarrowingScopeEvidence(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.ScopeEvidence = &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeProject),
			Name:       "child-project",
			UID:        "project-child-1",
		},
	}
	// Record is Tenant; Project evidence is structural narrowing.
	rec := apimeta.ScopeIdentity{Kind: apimeta.ScopeTenant, UID: "tenant-1"}
	p := ValidateEvaluation(ev, rec, validEvaluationProfile())
	if p != nil {
		t.Fatalf("structural narrowing scopeEvidence must pass, got %#v", p)
	}
}

func TestValidateEvaluation_TypeUnsupported(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.Evaluator.Type = "unknown-evaluator"
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	assertEvaluationViolation(t, p, CodeEvaluationTypeUnsupported, ptrEvalType)
}

func TestValidateEvaluation_TypeVersionUnsupported(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.Evaluator.Version = "9.9.9"
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	assertEvaluationViolation(t, p, CodeEvaluationTypeUnsupported, ptrEvalType)
}

func TestValidateEvaluation_ResultInvalidEmptyType(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.Evaluator.Type = ""
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	assertEvaluationViolation(t, p, CodeEvaluationResultInvalid, ptrEvalType)
}

func TestValidateEvaluation_ResultInvalidStatus(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.ResultStatus = "PENDING"
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	assertEvaluationViolation(t, p, CodeEvaluationResultInvalid, ptrEvalResultStatus)
}

func TestValidateEvaluation_ResultInvalidMissingInput(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.InputSnapshotRef = nil
	ev.InputIntegrity = nil
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	assertEvaluationViolation(t, p, CodeEvaluationResultInvalid, ptrEvalInputSnapshot)
}

func TestValidateEvaluation_ResultInvalidEmptyIntegrity(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.InputSnapshotRef = nil
	ev.InputIntegrity = &decision.OpaqueIntegrityCarrier{}
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	assertEvaluationViolation(t, p, CodeEvaluationResultInvalid, ptrEvalInputIntegrity)
}

func TestValidateEvaluation_ResultInvalidSuccessWithoutResult(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.Result = nil
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	assertEvaluationViolation(t, p, CodeEvaluationResultInvalid, ptrEvalResult)
}

func TestValidateEvaluation_ResultInvalidBadTiming(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.Timing.StartedAt = "not-a-timestamp"
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	assertEvaluationViolation(t, p, CodeEvaluationResultInvalid, ptrEvalTimingStarted)
}

func TestValidateEvaluation_ScopeMismatchWidening(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.ScopeEvidence = &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeOrganization),
			Name:       "acme",
			UID:        "org-1",
		},
	}
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	assertEvaluationViolation(t, p, CodeEvaluationScopeMismatch, ptrEvalScopeEvidence)
	if strings.Contains(strings.ToLower(p.Violations[0].Message), "exist") {
		t.Fatalf("scope mismatch must not disclose existence; message=%q", p.Violations[0].Message)
	}
	if strings.Contains(strings.ToLower(p.Violations[0].Message), "not found") {
		t.Fatalf("scope mismatch must not disclose absence; message=%q", p.Violations[0].Message)
	}
}

func TestValidateEvaluation_ScopeMismatchDifferentUID(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.ScopeEvidence = &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeProject),
			Name:       "other",
			UID:        "project-other",
		},
	}
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	assertEvaluationViolation(t, p, CodeEvaluationScopeMismatch, ptrEvalScopeEvidence)
}

func TestValidateEvaluation_ScopeMismatchUnknownKind(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.ScopeEvidence = &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       "ServiceInstance",
			Name:       "si-1",
			UID:        "si-uid-1",
		},
	}
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	assertEvaluationViolation(t, p, CodeEvaluationScopeMismatch, ptrEvalScopeEvidence)
	msg := strings.ToLower(p.Violations[0].Message)
	if strings.Contains(msg, "exist") || strings.Contains(msg, "not found") {
		t.Fatalf("unknown scope must not disclose existence; message=%q", p.Violations[0].Message)
	}
}

func TestValidateEvaluation_ScopeMismatchMissingUID(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.ScopeEvidence = &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeProject),
			Name:       "payments-production",
		},
	}
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	assertEvaluationViolation(t, p, CodeEvaluationScopeMismatch, ptrEvalScopeEvidence)
}

func TestValidateEvaluation_LimitExceededDuration(t *testing.T) {
	t.Parallel()

	prof := validEvaluationProfile()
	prof.Spec.Limits.MaxLatencyMs = 10
	ev := validEvaluation()
	ev.Timing.DurationMs = 50
	p := ValidateEvaluation(ev, projectRecScope("project-123"), prof)
	assertEvaluationViolation(t, p, CodeEvaluationLimitExceeded, ptrEvalTimingDuration)
}

func TestValidateEvaluation_LimitExceededPayload(t *testing.T) {
	t.Parallel()

	prof := validEvaluationProfile()
	prof.Spec.Limits.MaxPayloadBytes = 4
	ev := validEvaluation()
	ev.Result = json.RawMessage(`{"fit":true,"pool":"pool-a"}`)
	p := ValidateEvaluation(ev, projectRecScope("project-123"), prof)
	assertEvaluationViolation(t, p, CodeEvaluationLimitExceeded, ptrEvalResult)
}

func TestValidateEvaluation_EvaluatorWideningRejected(t *testing.T) {
	t.Parallel()

	prof := validEvaluationProfile()
	wide := decision.ProfileLimits{
		MaxPayloadBytes: prof.Spec.Limits.MaxPayloadBytes + 1,
	}
	prof.Spec.AcceptedEvaluationTypes = []decision.AcceptedEvaluationType{
		{Type: "capacity-evaluator", Version: "1.2.0", Limits: &wide},
	}
	p := ValidateEvaluation(validEvaluation(), projectRecScope("project-123"), prof)
	assertEvaluationViolation(t, p, CodeEvaluationLimitExceeded, "/evaluator")
}

func TestValidateEvaluation_EvaluatorNarrowingOK(t *testing.T) {
	t.Parallel()

	prof := validEvaluationProfile()
	narrow := decision.ProfileLimits{
		MaxLatencyMs:    500,
		MaxPayloadBytes: 4096,
	}
	prof.Spec.AcceptedEvaluationTypes = []decision.AcceptedEvaluationType{
		{Type: "capacity-evaluator", Version: "1.2.0", Limits: &narrow},
	}
	ev := validEvaluation()
	ev.Timing.DurationMs = 100
	p := ValidateEvaluation(ev, projectRecScope("project-123"), prof)
	if p != nil {
		t.Fatalf("evaluator narrowing must pass, got %#v", p)
	}
}

func TestValidateEvaluation_IntegrityCarrierAloneOK(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.InputSnapshotRef = nil
	ev.InputIntegrity = &decision.OpaqueIntegrityCarrier{
		State: "trusted",
	}
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	if p != nil {
		t.Fatalf("integrity-only input must pass, got %#v", p)
	}
}

func TestValidateEvaluation_NonSuccessWithoutResultOK(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.ResultStatus = decision.EvaluationResultStatusTimeout
	ev.Result = nil
	p := ValidateEvaluation(ev, projectRecScope("project-123"), validEvaluationProfile())
	if p != nil {
		t.Fatalf("TIMEOUT without result payload must pass, got %#v", p)
	}
}

func TestValidateEvaluation_PlatformRecordEqualEvidenceOK(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.ScopeEvidence = nil
	p := ValidateEvaluation(ev, apimeta.ScopeIdentity{
		Kind: apimeta.ScopePlatform,
		UID:  apimeta.PlatformScopeUID,
	}, validEvaluationProfile())
	if p != nil {
		t.Fatalf("Platform record without scopeEvidence must pass, got %#v", p)
	}
}
