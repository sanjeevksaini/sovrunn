package compose

import (
	"encoding/json"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

// Closed composition violation code used when Replay is invoked without a
// registered strategy (architecture §15.2). Canonical constants live in
// internal/decision/validate (T-017); this reproduces the closed string so
// compose can emit Problem Details now.
const violationStrategyUnsupported apiproblem.ViolationCode = "DECISION_COMPOSITION_STRATEGY_UNSUPPORTED"

// Replay is the deterministic composition reference kernel (AD-012, AD-037;
// RID-06; design §8; F13-EVAL-002, F13-SCALE-005).
//
// It applies a pure Strategy to already-captured EvaluationResult values and
// never reruns an evaluator. Captured results are the sole evaluation source;
// any Evaluations present on in are ignored.
//
// Replay is a pure function of versioned StrategyInput metadata (including
// fixture-local bounded configuration on Meta), the captured evaluations, and
// the registered strategy. It holds no session state and performs no
// persistence, delivery, locking, caching, networking, or scaling.
func Replay(captured []decision.EvaluationResult, s Strategy, in StrategyInput) (StrategyOutput, *apiproblem.Problem) {
	if s == nil {
		return StrategyOutput{}, composeProblem(
			violationStrategyUnsupported,
			"/composition/strategy",
			"composition replay requires a registered strategy",
		)
	}

	replayIn := StrategyInput{
		Meta:            cloneStrategyMetadata(in.Meta),
		FailureBehavior: cloneFailureBehavior(in.FailureBehavior),
		Evaluations:     cloneEvaluations(captured),
		InputRefs:       cloneStrings(in.InputRefs),
	}
	return s(replayIn)
}

func cloneStrategyMetadata(m StrategyMetadata) StrategyMetadata {
	out := m
	out.AcceptedInputTypes = cloneStrings(m.AcceptedInputTypes)
	return out
}

func cloneFailureBehavior(fb decision.FailureBehavior) decision.FailureBehavior {
	out := fb
	out.DeclaredFailureModes = cloneStrings(fb.DeclaredFailureModes)
	if fb.SecurityExceptionRef != nil {
		exc := *fb.SecurityExceptionRef
		exc.CompensatingControls = cloneStrings(fb.SecurityExceptionRef.CompensatingControls)
		exc.CoveredFailureModes = cloneStrings(fb.SecurityExceptionRef.CoveredFailureModes)
		out.SecurityExceptionRef = &exc
	}
	return out
}

func cloneEvaluations(in []decision.EvaluationResult) []decision.EvaluationResult {
	if len(in) == 0 {
		return nil
	}
	out := make([]decision.EvaluationResult, len(in))
	for i := range in {
		out[i] = cloneEvaluation(in[i])
	}
	return out
}

func cloneEvaluation(ev decision.EvaluationResult) decision.EvaluationResult {
	c := ev
	if ev.InputSnapshotRef != nil {
		ref := *ev.InputSnapshotRef
		c.InputSnapshotRef = &ref
	}
	if ev.InputIntegrity != nil {
		ic := *ev.InputIntegrity
		c.InputIntegrity = &ic
	}
	if ev.Result != nil {
		c.Result = append(json.RawMessage(nil), ev.Result...)
	}
	if ev.ExecutionConfig != nil {
		c.ExecutionConfig = append(json.RawMessage(nil), ev.ExecutionConfig...)
	}
	if ev.SafetyPolicyFilters != nil {
		c.SafetyPolicyFilters = append([]string(nil), ev.SafetyPolicyFilters...)
	}
	if ev.ScopeEvidence != nil {
		se := *ev.ScopeEvidence
		c.ScopeEvidence = &se
	}
	return c
}
