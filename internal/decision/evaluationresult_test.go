package decision

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func sampleEvaluationResult() EvaluationResult {
	return EvaluationResult{
		Evaluator: EvaluatorIdentity{
			Type:    "capacity-evaluator",
			Version: "1.2.0",
		},
		InputSnapshotRef: &apimeta.TypedRef{
			APIVersion: "governance.sovrunn.io/v1alpha1",
			Kind:       "InputSnapshot",
			Name:       "placement-input-1",
			UID:        "input-uid-1",
		},
		ResultStatus: EvaluationResultStatusSuccess,
		Result:       json.RawMessage(`{"fit":true,"pool":"pool-a"}`),
		Timing: TimingEnvelope{
			StartedAt:   "2026-07-29T10:00:00Z",
			CompletedAt: "2026-07-29T10:00:00.050Z",
			DurationMs:  50,
		},
		EvaluatedAt:     "2026-07-29T10:00:00.050Z",
		ExecutionConfig: json.RawMessage(`{"mode":"deterministic"}`),
		TrustBoundary:   "local-verified",
		SafetyPolicyFilters: []string{
			"deny-cross-jurisdiction",
		},
		StructuralTrustState: "trusted",
	}
}

func TestEvaluationResultJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := sampleEvaluationResult()
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out EvaluationResult
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if out.Evaluator != (EvaluatorIdentity{Type: "capacity-evaluator", Version: "1.2.0"}) {
		t.Fatalf("evaluator = %+v", out.Evaluator)
	}
	if out.InputSnapshotRef == nil || out.InputSnapshotRef.UID != "input-uid-1" {
		t.Fatalf("inputSnapshotRef = %#v", out.InputSnapshotRef)
	}
	if out.InputIntegrity != nil {
		t.Fatalf("inputIntegrity must be nil when omitted, got %#v", out.InputIntegrity)
	}
	if out.ResultStatus != EvaluationResultStatusSuccess {
		t.Fatalf("resultStatus = %q", out.ResultStatus)
	}
	if string(out.Result) != `{"fit":true,"pool":"pool-a"}` {
		t.Fatalf("result = %s", out.Result)
	}
	if out.Timing.DurationMs != 50 || out.Timing.StartedAt == "" || out.Timing.CompletedAt == "" {
		t.Fatalf("timing = %+v", out.Timing)
	}
	if out.EvaluatedAt != "2026-07-29T10:00:00.050Z" {
		t.Fatalf("evaluatedAt = %q", out.EvaluatedAt)
	}
	if out.StructuralTrustState != "trusted" {
		t.Fatalf("structuralTrustState = %q", out.StructuralTrustState)
	}
	if out.ScopeEvidence != nil {
		t.Fatalf("scopeEvidence must be nil by default (no independent scope), got %#v", out.ScopeEvidence)
	}
}

func TestEvaluationResultOmitemptyIntegrityCarrier(t *testing.T) {
	t.Parallel()

	in := sampleEvaluationResult()
	in.InputIntegrity = nil
	in.ExecutionConfig = nil
	in.TrustBoundary = ""
	in.SafetyPolicyFilters = nil
	in.StructuralTrustState = ""
	in.Result = nil
	in.Timing.DurationMs = 0
	in.ScopeEvidence = nil

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	raw := string(data)

	for _, key := range []string{
		`"inputIntegrity"`,
		`"executionConfig"`,
		`"trustBoundary"`,
		`"safetyPolicyFilters"`,
		`"structuralTrustState"`,
		`"result"`,
		`"durationMs"`,
		`"scopeEvidence"`,
	} {
		if strings.Contains(raw, key) {
			t.Fatalf("optional field %s must be omitted when unset; json=%s", key, raw)
		}
	}

	// Required carriers must remain present.
	for _, key := range []string{
		`"evaluator"`,
		`"resultStatus"`,
		`"timing"`,
		`"evaluatedAt"`,
		`"inputSnapshotRef"`,
	} {
		if !strings.Contains(raw, key) {
			t.Fatalf("required field %s missing; json=%s", key, raw)
		}
	}
}

func TestEvaluationResultIntegrityCarrierRoundTrip(t *testing.T) {
	t.Parallel()

	in := sampleEvaluationResult()
	in.InputIntegrity = &OpaqueIntegrityCarrier{
		State:                  "matched",
		CanonicalizationMethod: "carrier:opaque",
		DigestAlgorithm:        "carrier:opaque",
		CoveredFields:          "inputSnapshotRef",
		SignatureAlgorithm:     "carrier:opaque",
	}

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out EvaluationResult
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.InputIntegrity == nil || out.InputIntegrity.State != "matched" {
		t.Fatalf("inputIntegrity = %#v", out.InputIntegrity)
	}
	if out.InputIntegrity.CoveredFields != "inputSnapshotRef" {
		t.Fatalf("coveredFields = %q", out.InputIntegrity.CoveredFields)
	}
}

func TestEvaluationResultStatusEnumeration(t *testing.T) {
	t.Parallel()

	want := []EvaluationResultStatus{
		EvaluationResultStatusSuccess,
		EvaluationResultStatusNonMatch,
		EvaluationResultStatusIndeterminate,
		EvaluationResultStatusTimeout,
		EvaluationResultStatusError,
	}
	got := AllEvaluationResultStatuses()
	if len(got) != len(want) {
		t.Fatalf("AllEvaluationResultStatuses len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllEvaluationResultStatuses[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("%q must be Valid()", got[i])
		}
	}

	for _, bad := range []EvaluationResultStatus{"", "success", "FAILED", "PENDING", "ALLOWED"} {
		if bad.Valid() {
			t.Fatalf("%q must be rejected", bad)
		}
	}
}

func TestEvaluationResultNoIndependentScopeField(t *testing.T) {
	t.Parallel()

	in := sampleEvaluationResult()
	in.ScopeEvidence = &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       "Project",
			Name:       "payments-production",
			UID:        "project-123",
		},
	}

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	raw := string(data)

	// Evidence field is allowed; a second scope authority field is not.
	if !strings.Contains(raw, `"scopeEvidence"`) {
		t.Fatalf("scopeEvidence (evidence only) missing; json=%s", raw)
	}
	if strings.Contains(raw, `"scopeRef"`) {
		t.Fatalf("EvaluationResult must not carry independent scopeRef; json=%s", raw)
	}
	if strings.Contains(raw, `"metadata"`) {
		t.Fatalf("EvaluationResult must not carry ObjectMeta (EmbeddedValue); json=%s", raw)
	}
	if strings.Contains(raw, `"apiVersion"`) || strings.Contains(raw, `"kind"`) {
		// TypedRef inside inputSnapshotRef / scopeEvidence legitimately carries
		// apiVersion/kind; ensure top-level TypeMeta is absent by checking
		// the unmarshaled shape has empty top-level identity.
		var top map[string]json.RawMessage
		if err := json.Unmarshal(data, &top); err != nil {
			t.Fatalf("unmarshal top: %v", err)
		}
		if _, ok := top["apiVersion"]; ok {
			t.Fatal("top-level apiVersion must be absent (no TypeMeta on EmbeddedValue)")
		}
		if _, ok := top["kind"]; ok {
			t.Fatal("top-level kind must be absent (no TypeMeta on EmbeddedValue)")
		}
	}
}

func TestEvaluationResultStatusRoundTripEach(t *testing.T) {
	t.Parallel()

	for _, status := range AllEvaluationResultStatuses() {
		status := status
		t.Run(string(status), func(t *testing.T) {
			t.Parallel()
			in := sampleEvaluationResult()
			in.ResultStatus = status
			if status != EvaluationResultStatusSuccess {
				in.Result = json.RawMessage(`{"detail":"non-success"}`)
			}
			data, err := json.Marshal(in)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var out EvaluationResult
			if err := json.Unmarshal(data, &out); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if out.ResultStatus != status {
				t.Fatalf("resultStatus = %q, want %q", out.ResultStatus, status)
			}
		})
	}
}

func TestDecisionBodyEvaluationResultsEmbedded(t *testing.T) {
	t.Parallel()

	rec := sampleDecisionRecord()
	rec.Record.EvaluationResults = []EvaluationResult{sampleEvaluationResult()}

	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out DecisionRecord
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.Record.EvaluationResults) != 1 {
		t.Fatalf("evaluationResults len=%d, want 1", len(out.Record.EvaluationResults))
	}
	if out.Record.EvaluationResults[0].ResultStatus != EvaluationResultStatusSuccess {
		t.Fatalf("embedded resultStatus = %q", out.Record.EvaluationResults[0].ResultStatus)
	}
	// Containing record remains sole scope authority.
	if out.Metadata.ScopeRef == nil || out.Metadata.ScopeRef.UID != "project-123" {
		t.Fatalf("record metadata.scopeRef must remain sole scope authority, got %#v", out.Metadata.ScopeRef)
	}
}
