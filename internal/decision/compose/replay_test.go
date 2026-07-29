package compose

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

// recordingStrategy returns a pure Strategy that records the evaluations it
// received and produces a deterministic output from them. It never invokes an
// evaluator.
func recordingStrategy(seen *[]decision.EvaluationResult) Strategy {
	return func(in StrategyInput) (StrategyOutput, *apiproblem.Problem) {
		if seen != nil {
			cp := make([]decision.EvaluationResult, len(in.Evaluations))
			copy(cp, in.Evaluations)
			*seen = cp
		}
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
			out.Rationale.ReasonCodes = append(out.Rationale.ReasonCodes, "CAPTURED:"+ev.Evaluator.Type+"/"+ev.Evaluator.Version)
		}
		allowed := decision.AdjudicationOutcomeAllowed
		out.Adjudication = &allowed
		out.Rationale.Reasons = append(out.Rationale.Reasons, "replayed from captured evaluations")
		return out, nil
	}
}

func sampleCaptured() []decision.EvaluationResult {
	return []decision.EvaluationResult{
		{
			Evaluator:    decision.EvaluatorIdentity{Type: "capacity", Version: "1.0.0"},
			ResultStatus: decision.EvaluationResultStatusSuccess,
			Result:       json.RawMessage(`{"capacity":"ok","score":42}`),
			Timing: decision.TimingEnvelope{
				StartedAt:   "2026-07-29T10:00:00Z",
				CompletedAt: "2026-07-29T10:00:01Z",
				DurationMs:  1000,
			},
			EvaluatedAt:          "2026-07-29T10:00:01Z",
			StructuralTrustState: "TRUSTED",
			InputSnapshotRef: &apimeta.TypedRef{
				APIVersion: "governance.sovrunn.io/v1alpha1",
				Kind:       "InputSnapshot",
				Name:       "snap-1",
				UID:        "snap-uid-1",
			},
		},
		{
			Evaluator:    decision.EvaluatorIdentity{Type: "ai-advisory", Version: "2.1.0"},
			ResultStatus: decision.EvaluationResultStatusSuccess,
			Result:       json.RawMessage(`{"advice":"prefer-zone-a"}`),
			Timing: decision.TimingEnvelope{
				StartedAt:   "2026-07-29T10:00:00Z",
				CompletedAt: "2026-07-29T10:00:02Z",
				DurationMs:  2000,
			},
			EvaluatedAt:          "2026-07-29T10:00:02Z",
			StructuralTrustState: "CAPTURED",
			ExecutionConfig:      json.RawMessage(`{"temperature":0}`),
			SafetyPolicyFilters:  []string{"no-pii", "ai-non-authoritative"},
		},
	}
}

func TestReplayDeterministicIdenticalOutput(t *testing.T) {
	t.Parallel()

	captured := sampleCaptured()
	in := StrategyInput{
		Meta: sampleMeta(false),
		FailureBehavior: decision.FailureBehavior{
			FailPosture: decision.FailPostureFailClosed,
		},
		InputRefs: []string{"eval-capacity", "eval-ai"},
	}
	s := recordingStrategy(nil)

	out1, prob1 := Replay(captured, s, in)
	if prob1 != nil {
		t.Fatalf("first Replay problem = %#v", prob1)
	}
	out2, prob2 := Replay(captured, s, in)
	if prob2 != nil {
		t.Fatalf("second Replay problem = %#v", prob2)
	}

	if !reflect.DeepEqual(out1, out2) {
		t.Fatalf("identical captured input must yield identical output;\n out1=%#v\n out2=%#v", out1, out2)
	}
	if out1.StrategyID != "all-of" || out1.StrategyVersion != "1.0.0" {
		t.Fatalf("strategy identity = %s/%s", out1.StrategyID, out1.StrategyVersion)
	}
	if out1.Adjudication == nil || *out1.Adjudication != decision.AdjudicationOutcomeAllowed {
		t.Fatalf("adjudication = %#v", out1.Adjudication)
	}
	if string(out1.TypedResult) != `{"advice":"prefer-zone-a"}` {
		t.Fatalf("typedResult = %s", out1.TypedResult)
	}
	if len(out1.Rationale.ReasonCodes) != 2 {
		t.Fatalf("reasonCodes = %#v", out1.Rationale.ReasonCodes)
	}
}

func TestReplayNeverInvokesEvaluator(t *testing.T) {
	t.Parallel()

	evaluatorCalls := 0
	liveEvaluator := func() decision.EvaluationResult {
		evaluatorCalls++
		return decision.EvaluationResult{
			Evaluator:    decision.EvaluatorIdentity{Type: "live", Version: "9.9.9"},
			ResultStatus: decision.EvaluationResultStatusSuccess,
			Result:       json.RawMessage(`{"live":true}`),
			Timing: decision.TimingEnvelope{
				StartedAt:   "2026-07-29T11:00:00Z",
				CompletedAt: "2026-07-29T11:00:01Z",
			},
			EvaluatedAt: "2026-07-29T11:00:01Z",
		}
	}

	captured := sampleCaptured()
	// Deliberately place a "live" evaluation on the input and keep a live
	// evaluator available. Replay must use only captured and never call live.
	in := StrategyInput{
		Meta: sampleMeta(false),
		FailureBehavior: decision.FailureBehavior{
			FailPosture: decision.FailPostureFailClosed,
		},
		Evaluations: []decision.EvaluationResult{liveEvaluator()},
		InputRefs:   []string{"eval-capacity", "eval-ai"},
	}
	if evaluatorCalls != 1 {
		t.Fatalf("setup expected one live call, got %d", evaluatorCalls)
	}

	var seen []decision.EvaluationResult
	out, prob := Replay(captured, recordingStrategy(&seen), in)
	if prob != nil {
		t.Fatalf("Replay problem = %#v", prob)
	}
	if evaluatorCalls != 1 {
		t.Fatalf("Replay must never invoke an evaluator; live calls = %d", evaluatorCalls)
	}
	if len(seen) != 2 {
		t.Fatalf("strategy must receive captured evaluations only; seen=%d", len(seen))
	}
	if seen[0].Evaluator.Type != "capacity" || seen[1].Evaluator.Type != "ai-advisory" {
		t.Fatalf("strategy saw %#v, want captured evaluator identities", seen)
	}
	for _, ev := range seen {
		if ev.Evaluator.Type == "live" {
			t.Fatal("Replay must ignore StrategyInput.Evaluations (live path)")
		}
	}
	if string(out.TypedResult) == `{"live":true}` {
		t.Fatal("Replay must not surface live evaluator output")
	}
}

func TestReplayIgnoresInputEvaluationsInFavorOfCaptured(t *testing.T) {
	t.Parallel()

	captured := []decision.EvaluationResult{{
		Evaluator:    decision.EvaluatorIdentity{Type: "captured", Version: "1.0.0"},
		ResultStatus: decision.EvaluationResultStatusSuccess,
		Result:       json.RawMessage(`{"source":"captured"}`),
		Timing: decision.TimingEnvelope{
			StartedAt:   "2026-07-29T10:00:00Z",
			CompletedAt: "2026-07-29T10:00:01Z",
		},
		EvaluatedAt: "2026-07-29T10:00:01Z",
	}}
	stale := []decision.EvaluationResult{{
		Evaluator:    decision.EvaluatorIdentity{Type: "stale-input", Version: "0.0.1"},
		ResultStatus: decision.EvaluationResultStatusSuccess,
		Result:       json.RawMessage(`{"source":"stale"}`),
		Timing: decision.TimingEnvelope{
			StartedAt:   "2026-07-29T09:00:00Z",
			CompletedAt: "2026-07-29T09:00:01Z",
		},
		EvaluatedAt: "2026-07-29T09:00:01Z",
	}}

	var seen []decision.EvaluationResult
	out, prob := Replay(captured, recordingStrategy(&seen), StrategyInput{
		Meta:        sampleMeta(false),
		Evaluations: stale,
		InputRefs:   []string{"captured-only"},
	})
	if prob != nil {
		t.Fatalf("Replay problem = %#v", prob)
	}
	if len(seen) != 1 || seen[0].Evaluator.Type != "captured" {
		t.Fatalf("captured must replace in.Evaluations; seen=%#v", seen)
	}
	if string(out.TypedResult) != `{"source":"captured"}` {
		t.Fatalf("typedResult = %s", out.TypedResult)
	}
}

func TestReplayNilStrategyFailsClosed(t *testing.T) {
	t.Parallel()

	_, prob := Replay(sampleCaptured(), nil, StrategyInput{Meta: sampleMeta(false)})
	if prob == nil {
		t.Fatal("nil strategy must fail closed")
	}
	if prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("top-level code = %q", prob.Code)
	}
	if len(prob.Violations) != 1 || prob.Violations[0].Code != violationStrategyUnsupported {
		t.Fatalf("violations = %#v", prob.Violations)
	}
	if prob.Violations[0].Field != "/composition/strategy" {
		t.Fatalf("field = %q", prob.Violations[0].Field)
	}
}

func TestReplayHoldsNoSessionState(t *testing.T) {
	t.Parallel()

	s := recordingStrategy(nil)
	in := StrategyInput{
		Meta: sampleMeta(false),
		FailureBehavior: decision.FailureBehavior{
			FailPosture: decision.FailPostureFailClosed,
		},
		InputRefs: []string{"a"},
	}

	first := []decision.EvaluationResult{{
		Evaluator:    decision.EvaluatorIdentity{Type: "first", Version: "1"},
		ResultStatus: decision.EvaluationResultStatusSuccess,
		Result:       json.RawMessage(`{"n":1}`),
		Timing:       decision.TimingEnvelope{StartedAt: "2026-07-29T10:00:00Z", CompletedAt: "2026-07-29T10:00:01Z"},
		EvaluatedAt:  "2026-07-29T10:00:01Z",
	}}
	second := []decision.EvaluationResult{{
		Evaluator:    decision.EvaluatorIdentity{Type: "second", Version: "1"},
		ResultStatus: decision.EvaluationResultStatusSuccess,
		Result:       json.RawMessage(`{"n":2}`),
		Timing:       decision.TimingEnvelope{StartedAt: "2026-07-29T10:00:00Z", CompletedAt: "2026-07-29T10:00:01Z"},
		EvaluatedAt:  "2026-07-29T10:00:01Z",
	}}

	out1, prob1 := Replay(first, s, in)
	if prob1 != nil {
		t.Fatalf("first = %#v", prob1)
	}
	out2, prob2 := Replay(second, s, in)
	if prob2 != nil {
		t.Fatalf("second = %#v", prob2)
	}
	if string(out1.TypedResult) != `{"n":1}` || string(out2.TypedResult) != `{"n":2}` {
		t.Fatalf("session bleed: out1=%s out2=%s", out1.TypedResult, out2.TypedResult)
	}

	// Replay clones captured evaluations so caller mutation after return
	// cannot rewrite what the strategy observed during the call.
	captured := sampleCaptured()
	var seen []decision.EvaluationResult
	if _, prob := Replay(captured, recordingStrategy(&seen), in); prob != nil {
		t.Fatalf("isolation Replay = %#v", prob)
	}
	if len(seen) == 0 {
		t.Fatal("expected strategy to observe captured evaluations")
	}
	original := string(captured[0].Result)
	seen[0].Result = json.RawMessage(`{"mutated":true}`)
	if string(captured[0].Result) != original {
		t.Fatal("Replay must clone captured evaluations; caller slice must remain unchanged")
	}
}

func TestReplayEmptyCapturedIsDeterministic(t *testing.T) {
	t.Parallel()

	s := recordingStrategy(nil)
	in := StrategyInput{
		Meta:      sampleMeta(false),
		InputRefs: []string{},
	}
	out1, prob1 := Replay(nil, s, in)
	out2, prob2 := Replay([]decision.EvaluationResult{}, s, in)
	if prob1 != nil || prob2 != nil {
		t.Fatalf("empty captured problems: %#v %#v", prob1, prob2)
	}
	if !reflect.DeepEqual(out1, out2) {
		t.Fatalf("nil and empty captured must compose identically;\n out1=%#v\n out2=%#v", out1, out2)
	}
}
