package validate

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

func TestEvaluationStructureExport_ParityWithPrivateCheck(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		mut  func(decision.EvaluationResult) decision.EvaluationResult
	}{
		{name: "valid", mut: func(e decision.EvaluationResult) decision.EvaluationResult { return e }},
		{name: "empty_evaluator_type", mut: func(e decision.EvaluationResult) decision.EvaluationResult {
			e.Evaluator.Type = ""
			return e
		}},
		{name: "blank_evaluator_version", mut: func(e decision.EvaluationResult) decision.EvaluationResult {
			e.Evaluator.Version = "   "
			return e
		}},
		{name: "invalid_status", mut: func(e decision.EvaluationResult) decision.EvaluationResult {
			e.ResultStatus = "PENDING"
			return e
		}},
		{name: "missing_input_identity", mut: func(e decision.EvaluationResult) decision.EvaluationResult {
			e.InputSnapshotRef = nil
			e.InputIntegrity = nil
			return e
		}},
		{name: "empty_integrity", mut: func(e decision.EvaluationResult) decision.EvaluationResult {
			e.InputSnapshotRef = nil
			e.InputIntegrity = &decision.OpaqueIntegrityCarrier{}
			return e
		}},
		{name: "blank_integrity_state", mut: func(e decision.EvaluationResult) decision.EvaluationResult {
			e.InputSnapshotRef = nil
			e.InputIntegrity = &decision.OpaqueIntegrityCarrier{
				State:                  "   ",
				CanonicalizationMethod: "jcs",
			}
			return e
		}},
		{name: "success_without_result", mut: func(e decision.EvaluationResult) decision.EvaluationResult {
			e.Result = nil
			return e
		}},
		{name: "bad_timing", mut: func(e decision.EvaluationResult) decision.EvaluationResult {
			e.Timing.StartedAt = "not-a-timestamp"
			return e
		}},
		{name: "empty_evaluated_at", mut: func(e decision.EvaluationResult) decision.EvaluationResult {
			e.EvaluatedAt = ""
			return e
		}},
		{name: "negative_duration", mut: func(e decision.EvaluationResult) decision.EvaluationResult {
			e.Timing.DurationMs = -1
			return e
		}},
		{name: "integrity_only_valid", mut: func(e decision.EvaluationResult) decision.EvaluationResult {
			e.InputSnapshotRef = nil
			e.InputIntegrity = &decision.OpaqueIntegrityCarrier{State: "bound"}
			return e
		}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ev := tc.mut(validEvaluation())
			want := checkEvaluationResultStructure(ev) == nil
			got := IsEvaluationResultStructurallyValid(ev)
			if got != want {
				t.Fatalf("IsEvaluationResultStructurallyValid=%v, want exact private-check parity %v", got, want)
			}
		})
	}
}

func TestEvaluationStructureExport_ValidPasses(t *testing.T) {
	t.Parallel()

	if !IsEvaluationResultStructurallyValid(validEvaluation()) {
		t.Fatal("validEvaluation must be structurally valid")
	}
}

func TestEvaluationStructureExport_InvalidFails(t *testing.T) {
	t.Parallel()

	ev := validEvaluation()
	ev.Evaluator.Type = ""
	if IsEvaluationResultStructurallyValid(ev) {
		t.Fatal("empty evaluator.type must fail structural check")
	}
}
