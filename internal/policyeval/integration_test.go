package policyeval

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

// TestIntegration_AllFourOutcomesEndToEnd proves request → validation → digest →
// fake → result → timing → mapper → EvaluationResult for every closed Outcome
// (AC-F17-11, AC-F17-19) and obligations remain absent (AC-F17-18).
func TestIntegration_AllFourOutcomesEndToEnd(t *testing.T) {
	t.Parallel()

	fixed := time.Date(2026, 8, 26, 7, 30, 0, 0, time.UTC)
	wantAt := fixed.Format(time.RFC3339Nano)

	cases := []struct {
		outcome    Outcome
		reasonCode string
		wantStatus decision.EvaluationResultStatus
	}{
		{OutcomeAllow, "POLICY_ALLOW", decision.EvaluationResultStatusSuccess},
		{OutcomeDeny, "POLICY_DENY", decision.EvaluationResultStatusSuccess},
		{OutcomeRequiresApproval, "POLICY_REQUIRES_APPROVAL", decision.EvaluationResultStatusSuccess},
		{OutcomeIndeterminate, "POLICY_INDETERMINATE", decision.EvaluationResultStatusIndeterminate},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(string(tc.outcome), func(t *testing.T) {
			t.Parallel()

			fake := mustFake(t, []FakeFixture{{
				InputDigest: frozenInputDigest,
				Conclusion: &Conclusion{
					Outcome:     tc.outcome,
					ReasonCodes: []string{tc.reasonCode},
				},
			}})
			b := mustBoundary(t, fake, NewFixedTimeSource(fixed))

			result, timing, nr, err := b.Evaluate(context.Background(), frozenVectorRequest())
			if err != nil || nr != "" {
				t.Fatalf("Evaluate: nr=%q err=%v", nr, err)
			}
			if result.Outcome != tc.outcome {
				t.Fatalf("outcome=%q, want %q", result.Outcome, tc.outcome)
			}
			if result.InputDigest != frozenInputDigest {
				t.Fatalf("inputDigest=%q, want frozen vector", result.InputDigest)
			}
			if len(result.ReasonCodes) != 1 || result.ReasonCodes[0] != tc.reasonCode {
				t.Fatalf("reasonCodes=%v, want [%s]", result.ReasonCodes, tc.reasonCode)
			}
			if result.EvaluatedAt != wantAt || timing.CompletedAt != wantAt || timing.StartedAt != wantAt {
				t.Fatalf("timing mismatch: result=%q timing=%+v want=%q", result.EvaluatedAt, timing, wantAt)
			}
			if result.Obligations != nil {
				t.Fatalf("Phase 2R fake must leave Obligations nil, got %#v", result.Obligations)
			}

			mapped, err := MapToEvaluationResult(
				result,
				timing,
				validMapperEvaluator(),
				validMapperSnapshotIdent(),
			)
			if err != nil {
				t.Fatalf("MapToEvaluationResult: %v", err)
			}
			if mapped.ResultStatus != tc.wantStatus {
				t.Fatalf("resultStatus=%q, want %q", mapped.ResultStatus, tc.wantStatus)
			}
			if mapped.EvaluatedAt != timing.CompletedAt {
				t.Fatalf("evaluatedAt=%q != completedAt=%q", mapped.EvaluatedAt, timing.CompletedAt)
			}
			if mapped.Timing.DurationMs != 0 {
				t.Fatalf("DurationMs must remain omitted/zero, got %d", mapped.Timing.DurationMs)
			}
			if mapped.ExecutionConfig != nil {
				t.Fatalf("ExecutionConfig must be omitted, got %s", mapped.ExecutionConfig)
			}
			if mapped.TrustBoundary != "" {
				t.Fatalf("TrustBoundary must be omitted, got %q", mapped.TrustBoundary)
			}
			if mapped.SafetyPolicyFilters != nil {
				t.Fatalf("SafetyPolicyFilters must be omitted, got %v", mapped.SafetyPolicyFilters)
			}
			if mapped.StructuralTrustState != "" {
				t.Fatalf("StructuralTrustState must be omitted, got %q", mapped.StructuralTrustState)
			}
			if mapped.ScopeEvidence != nil {
				t.Fatalf("ScopeEvidence must be omitted, got %+v", mapped.ScopeEvidence)
			}

			var wire PolicyEvaluationResult
			if err := json.Unmarshal(mapped.Result, &wire); err != nil {
				t.Fatalf("unmarshal mapped Result: %v", err)
			}
			if wire.Outcome != tc.outcome || wire.Obligations != nil {
				t.Fatalf("mapped wire result=%+v, want outcome=%q obligations=nil", wire, tc.outcome)
			}
			if strings.Contains(string(mapped.Result), `"obligations"`) {
				t.Fatalf("mapped Result JSON must omit obligations: %s", mapped.Result)
			}
		})
	}
}

// TestIntegration_NonResultProducesNoMapping proves each NonResult category
// yields zero result, zero timing, the specific category, a sanitized error,
// zero pre-invocation adapter calls where required, and no successful mapping
// (AC-F17-21).
func TestIntegration_NonResultProducesNoMapping(t *testing.T) {
	t.Parallel()

	fixed := NewFixedTimeSource(time.Date(2026, 8, 26, 7, 30, 0, 0, time.UTC))

	t.Run("RequestInvalid", func(t *testing.T) {
		t.Parallel()
		spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
		b := mustBoundary(t, spy, fixed)

		req := frozenVectorRequest()
		req.Action = ""

		result, timing, nr, err := b.Evaluate(context.Background(), req)
		assertZeroNonResult(t, result, timing, nr, NonResultRequestInvalid, err, errRequestInvalid, spy.callCount())
		assertMappingRejectedForZero(t, result, timing)
	})

	t.Run("Canceled", func(t *testing.T) {
		t.Parallel()
		spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
		b := mustBoundary(t, spy, fixed)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		result, timing, nr, err := b.Evaluate(ctx, frozenVectorRequest())
		assertZeroNonResult(t, result, timing, nr, NonResultCanceled, err, errCanceled, spy.callCount())
		assertMappingRejectedForZero(t, result, timing)
	})

	t.Run("DeadlineExceeded", func(t *testing.T) {
		t.Parallel()
		spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
		b := mustBoundary(t, spy, fixed)

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
		defer cancel()

		result, timing, nr, err := b.Evaluate(ctx, frozenVectorRequest())
		assertZeroNonResult(t, result, timing, nr, NonResultDeadlineExceeded, err, errDeadlineExceeded, spy.callCount())
		assertMappingRejectedForZero(t, result, timing)
	})

	t.Run("AdapterFailure", func(t *testing.T) {
		t.Parallel()
		spy := &evaluateSpyAdapter{err: errors.New("raw-adapter-failure-detail")}
		b := mustBoundary(t, spy, fixed)

		result, timing, nr, err := b.Evaluate(context.Background(), frozenVectorRequest())
		if spy.callCount() != 1 {
			t.Fatalf("adapter calls=%d, want 1", spy.callCount())
		}
		if nr != NonResultAdapterFailure {
			t.Fatalf("NonResult=%q, want AdapterFailure", nr)
		}
		if !errors.Is(err, errAdapterFailed) {
			t.Fatalf("err=%v, want errAdapterFailed", err)
		}
		assertZeroResultTiming(t, result, timing)
		assertMappingRejectedForZero(t, result, timing)
	})

	t.Run("InvalidAdapterResult", func(t *testing.T) {
		t.Parallel()
		spy := &evaluateSpyAdapter{
			conclusion: Conclusion{
				Outcome:     Outcome("NotValidOutcome"),
				ReasonCodes: []string{"POLICY_ALLOW"},
			},
		}
		b := mustBoundary(t, spy, fixed)

		result, timing, nr, err := b.Evaluate(context.Background(), frozenVectorRequest())
		if spy.callCount() != 1 {
			t.Fatalf("adapter calls=%d, want 1", spy.callCount())
		}
		if nr != NonResultInvalidAdapterResult {
			t.Fatalf("NonResult=%q, want InvalidAdapterResult", nr)
		}
		if !errors.Is(err, errInvalidAdapterResult) {
			t.Fatalf("err=%v, want errInvalidAdapterResult", err)
		}
		assertZeroResultTiming(t, result, timing)
		assertMappingRejectedForZero(t, result, timing)
	})
}

// TestIntegration_FakeObligationsAbsent is an explicit Phase 2R obligations
// absence check across the fake → boundary path (AC-F17-18).
func TestIntegration_FakeObligationsAbsent(t *testing.T) {
	t.Parallel()

	for _, outcome := range AllOutcomes() {
		outcome := outcome
		t.Run(string(outcome), func(t *testing.T) {
			t.Parallel()
			fake := mustFake(t, []FakeFixture{{
				InputDigest: frozenInputDigest,
				Conclusion: &Conclusion{
					Outcome:     outcome,
					ReasonCodes: []string{"POLICY_CHECK"},
				},
			}})
			b := mustBoundary(t, fake, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))
			result, _, nr, err := b.Evaluate(context.Background(), frozenVectorRequest())
			if err != nil || nr != "" {
				t.Fatalf("Evaluate: nr=%q err=%v", nr, err)
			}
			if result.Obligations != nil {
				t.Fatalf("Obligations=%v, want nil", result.Obligations)
			}
			raw, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("json.Marshal: %v", err)
			}
			if strings.Contains(string(raw), `"obligations"`) {
				t.Fatalf("obligations present in JSON: %s", raw)
			}
		})
	}
}

func assertMappingRejectedForZero(t *testing.T, result PolicyEvaluationResult, timing decision.TimingEnvelope) {
	t.Helper()
	assertZeroResultTiming(t, result, timing)
	_, err := MapToEvaluationResult(result, timing, validMapperEvaluator(), validMapperSnapshotIdent())
	if !errors.Is(err, errMappingInputInvalid) {
		t.Fatalf("mapping zero non-result input: err=%v, want errMappingInputInvalid", err)
	}
}
