package policyeval

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

type evaluateSpyAdapter struct {
	mu         sync.Mutex
	calls      int
	lastIn     NormalizedInput
	lastDigest string
	conclusion Conclusion
	err        error
	onCall     func(ctx context.Context)
}

func (s *evaluateSpyAdapter) Evaluate(ctx context.Context, in NormalizedInput, digest string) (Conclusion, error) {
	s.mu.Lock()
	s.calls++
	s.lastIn = in
	s.lastDigest = digest
	fn := s.onCall
	conclusion := s.conclusion
	err := s.err
	s.mu.Unlock()
	if fn != nil {
		fn(ctx)
	}
	return conclusion, err
}

func (s *evaluateSpyAdapter) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
}

type errCountContext struct {
	context.Context
	mu    sync.Mutex
	n     int32
	after int32
	err   error
}

func (c *errCountContext) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
	if c.n >= c.after {
		return c.err
	}
	return nil
}

func (c *errCountContext) Done() <-chan struct{} {
	return nil
}

func (c *errCountContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (c *errCountContext) Value(key any) any {
	return nil
}

type ptrTimeSource struct {
	at time.Time
}

func (p *ptrTimeSource) NowUTC() time.Time {
	return p.at.UTC()
}

func validEvaluateRequest() PolicyEvaluationRequest {
	return frozenVectorRequest()
}

func validAllowConclusion() Conclusion {
	return Conclusion{
		Outcome:     OutcomeAllow,
		ReasonCodes: []string{"POLICY_ALLOW"},
	}
}

func mustBoundary(t *testing.T, adapter PolicyEngineAdapter, ts TimeSource) *Boundary {
	t.Helper()
	b, err := NewBoundary(adapter, ts)
	if err != nil {
		t.Fatalf("NewBoundary: %v", err)
	}
	return b
}

func assertZeroResultTiming(t *testing.T, result PolicyEvaluationResult, timing decision.TimingEnvelope) {
	t.Helper()
	if result.Outcome != "" || result.InputDigest != "" || result.EvaluatedAt != "" ||
		result.ReasonCodes != nil || result.Obligations != nil {
		t.Fatalf("result not zero: %+v", result)
	}
	if timing != (decision.TimingEnvelope{}) {
		t.Fatalf("timing not zero: %+v", timing)
	}
}

func assertZeroNonResult(
	t *testing.T,
	result PolicyEvaluationResult,
	timing decision.TimingEnvelope,
	got NonResult,
	want NonResult,
	err error,
	wantErr error,
	calls int,
) {
	t.Helper()
	if got != want {
		t.Fatalf("NonResult=%q, want %q", got, want)
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("err=%v, want %v", err, wantErr)
	}
	assertZeroResultTiming(t, result, timing)
	if calls != 0 {
		t.Fatalf("adapter calls=%d, want 0", calls)
	}
}

func TestEvaluate_NilContextRequestInvalid(t *testing.T) {
	t.Parallel()

	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	result, timing, nr, err := b.Evaluate(nil, validEvaluateRequest())
	assertZeroNonResult(t, result, timing, nr, NonResultRequestInvalid, err, errRequestInvalid, spy.callCount())
}

func TestEvaluate_PreCallCanceled(t *testing.T) {
	t.Parallel()

	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, timing, nr, err := b.Evaluate(ctx, validEvaluateRequest())
	assertZeroNonResult(t, result, timing, nr, NonResultCanceled, err, errCanceled, spy.callCount())
}

func TestEvaluate_PreCallDeadlineExceeded(t *testing.T) {
	t.Parallel()

	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	result, timing, nr, err := b.Evaluate(ctx, validEvaluateRequest())
	assertZeroNonResult(t, result, timing, nr, NonResultDeadlineExceeded, err, errDeadlineExceeded, spy.callCount())
}

func TestEvaluate_RequestValidationFailure(t *testing.T) {
	t.Parallel()

	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	req := validEvaluateRequest()
	req.Action = "" // bound violation

	result, timing, nr, err := b.Evaluate(context.Background(), req)
	assertZeroNonResult(t, result, timing, nr, NonResultRequestInvalid, err, errRequestInvalid, spy.callCount())

	if err != nil && strings.Contains(err.Error(), req.RequestID) {
		t.Fatalf("sanitized error leaked requestId: %v", err)
	}
}

func TestEvaluate_PreInvocationCanceled(t *testing.T) {
	t.Parallel()

	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	// First ctx.Err() (step 1) succeeds; second (step 4) returns Canceled.
	ctx := &errCountContext{after: 2, err: context.Canceled}

	result, timing, nr, err := b.Evaluate(ctx, validEvaluateRequest())
	assertZeroNonResult(t, result, timing, nr, NonResultCanceled, err, errCanceled, spy.callCount())
}

func TestEvaluate_PreInvocationDeadlineExceeded(t *testing.T) {
	t.Parallel()

	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	ctx := &errCountContext{after: 2, err: context.DeadlineExceeded}

	result, timing, nr, err := b.Evaluate(ctx, validEvaluateRequest())
	assertZeroNonResult(t, result, timing, nr, NonResultDeadlineExceeded, err, errDeadlineExceeded, spy.callCount())
}

func TestEvaluate_ExactlyOneAdapterCall(t *testing.T) {
	t.Parallel()

	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
	fixed := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	b := mustBoundary(t, spy, NewFixedTimeSource(fixed))

	result, timing, nr, err := b.Evaluate(context.Background(), validEvaluateRequest())
	if err != nil || nr != "" {
		t.Fatalf("Evaluate: nr=%q err=%v", nr, err)
	}
	if spy.callCount() != 1 {
		t.Fatalf("adapter calls=%d, want 1", spy.callCount())
	}
	if result.Outcome != OutcomeAllow {
		t.Fatalf("outcome=%q, want Allow", result.Outcome)
	}
	if timing.StartedAt == "" || timing.CompletedAt == "" {
		t.Fatalf("timing missing: %+v", timing)
	}
}

func TestEvaluate_PostReturnCanceledDiscardsLateOutput(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	spy := &evaluateSpyAdapter{
		conclusion: validAllowConclusion(),
		onCall: func(context.Context) {
			cancel()
		},
	}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	result, timing, nr, err := b.Evaluate(ctx, validEvaluateRequest())
	if spy.callCount() != 1 {
		t.Fatalf("adapter calls=%d, want 1", spy.callCount())
	}
	if nr != NonResultCanceled {
		t.Fatalf("NonResult=%q, want Canceled", nr)
	}
	if !errors.Is(err, errCanceled) {
		t.Fatalf("err=%v, want errCanceled", err)
	}
	assertZeroResultTiming(t, result, timing)
}

func TestEvaluate_PostReturnDeadlineDiscardsLateOutput(t *testing.T) {
	t.Parallel()

	// Steps 1 and 4 see nil; step 6 (third Err check) observes deadline.
	deadlineCtx := &errCountContext{after: 3, err: context.DeadlineExceeded}
	spy := &evaluateSpyAdapter{
		conclusion: Conclusion{Outcome: Outcome("Bogus"), ReasonCodes: []string{"bad"}},
		err:        errors.New("SECRET_ADAPTER_DETAIL"),
	}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	result, timing, nr, err := b.Evaluate(deadlineCtx, validEvaluateRequest())
	if spy.callCount() != 1 {
		t.Fatalf("adapter calls=%d, want 1", spy.callCount())
	}
	if nr != NonResultDeadlineExceeded {
		t.Fatalf("NonResult=%q, want DeadlineExceeded", nr)
	}
	if !errors.Is(err, errDeadlineExceeded) {
		t.Fatalf("err=%v, want errDeadlineExceeded", err)
	}
	assertZeroResultTiming(t, result, timing)
	if err != nil && strings.Contains(err.Error(), "SECRET_ADAPTER_DETAIL") {
		t.Fatalf("sanitized error leaked adapter detail: %v", err)
	}
	if err != nil && strings.Contains(err.Error(), "Bogus") {
		t.Fatalf("sanitized error leaked outcome: %v", err)
	}
}

func TestEvaluate_AdapterFailurePrecedencesMalformedConclusion(t *testing.T) {
	t.Parallel()

	spy := &evaluateSpyAdapter{
		conclusion: Conclusion{
			Outcome:     Outcome("NotAnOutcome"),
			ReasonCodes: []string{"bad-code", "bad-code"},
		},
		err: errors.New("SECRET_ADAPTER_DETAIL"),
	}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	result, timing, nr, err := b.Evaluate(context.Background(), validEvaluateRequest())
	if spy.callCount() != 1 {
		t.Fatalf("adapter calls=%d, want 1", spy.callCount())
	}
	if nr != NonResultAdapterFailure {
		t.Fatalf("NonResult=%q, want AdapterFailure (precedence over InvalidAdapterResult)", nr)
	}
	if !errors.Is(err, errAdapterFailed) {
		t.Fatalf("err=%v, want errAdapterFailed", err)
	}
	assertZeroResultTiming(t, result, timing)
	if strings.Contains(err.Error(), "SECRET_ADAPTER_DETAIL") {
		t.Fatalf("sanitized error leaked adapter detail: %v", err)
	}
	if strings.Contains(err.Error(), "NotAnOutcome") || strings.Contains(err.Error(), "bad-code") {
		t.Fatalf("sanitized error leaked conclusion values: %v", err)
	}
}

func TestEvaluate_InvalidOutcome(t *testing.T) {
	t.Parallel()

	spy := &evaluateSpyAdapter{
		conclusion: Conclusion{
			Outcome:     Outcome("Unknown"),
			ReasonCodes: []string{"POLICY_ALLOW"},
		},
	}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	result, timing, nr, err := b.Evaluate(context.Background(), validEvaluateRequest())
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
	if strings.Contains(err.Error(), "Unknown") {
		t.Fatalf("sanitized error leaked outcome: %v", err)
	}
}

func TestEvaluate_InvalidReasonCodes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		codes []string
	}{
		{name: "grammar_lowercase", codes: []string{"allow"}},
		{name: "grammar_leading_digit", codes: []string{"9START"}},
		{name: "duplicate", codes: []string{"POLICY_ALLOW", "POLICY_ALLOW"}},
		{name: "count_zero", codes: []string{}},
		{name: "count_over_32", codes: reasonCodesN(33)},
		{name: "nil_codes", codes: nil},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			spy := &evaluateSpyAdapter{
				conclusion: Conclusion{
					Outcome:     OutcomeAllow,
					ReasonCodes: tc.codes,
				},
			}
			b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))
			result, timing, nr, err := b.Evaluate(context.Background(), validEvaluateRequest())
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
			for _, code := range tc.codes {
				if code != "" && strings.Contains(err.Error(), code) {
					t.Fatalf("sanitized error leaked reason code %q: %v", code, err)
				}
			}
		})
	}
}

func reasonCodesN(n int) []string {
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = "C" + padIndex(i)
	}
	return out
}

func padIndex(i int) string {
	const digits = "0123456789"
	if i < 10 {
		return "0" + string(digits[i])
	}
	return string(digits[i/10]) + string(digits[i%10])
}

func TestEvaluate_ValidConclusion(t *testing.T) {
	t.Parallel()

	fixed := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	spy := &evaluateSpyAdapter{
		conclusion: Conclusion{
			Outcome:     OutcomeDeny,
			ReasonCodes: []string{"ZULU", "ALPHA", "MIKE"},
		},
	}
	b := mustBoundary(t, spy, NewFixedTimeSource(fixed))
	req := validEvaluateRequest()

	result, timing, nr, err := b.Evaluate(context.Background(), req)
	if err != nil || nr != "" {
		t.Fatalf("Evaluate: nr=%q err=%v", nr, err)
	}
	if spy.callCount() != 1 {
		t.Fatalf("adapter calls=%d, want 1", spy.callCount())
	}

	wantDigest, digErr := digestPolicyEvaluationRequest(req)
	if digErr != nil {
		t.Fatalf("digestPolicyEvaluationRequest: %v", digErr)
	}
	wantAt := fixed.Format(time.RFC3339Nano)

	if result.Outcome != OutcomeDeny {
		t.Fatalf("outcome=%q, want Deny", result.Outcome)
	}
	wantCodes := []string{"ALPHA", "MIKE", "ZULU"}
	if len(result.ReasonCodes) != len(wantCodes) {
		t.Fatalf("reasonCodes=%v, want %v", result.ReasonCodes, wantCodes)
	}
	for i := range wantCodes {
		if result.ReasonCodes[i] != wantCodes[i] {
			t.Fatalf("reasonCodes=%v, want %v", result.ReasonCodes, wantCodes)
		}
	}
	if result.InputDigest != wantDigest {
		t.Fatalf("inputDigest=%q, want %q", result.InputDigest, wantDigest)
	}
	if result.EvaluatedAt != wantAt {
		t.Fatalf("evaluatedAt=%q, want %q", result.EvaluatedAt, wantAt)
	}
	if result.EvaluatedAt != timing.CompletedAt {
		t.Fatalf("evaluatedAt=%q != completedAt=%q", result.EvaluatedAt, timing.CompletedAt)
	}
	if timing.StartedAt != wantAt || timing.CompletedAt != wantAt {
		t.Fatalf("timing=%+v, want started/completed %q", timing, wantAt)
	}
	if timing.DurationMs != 0 {
		t.Fatalf("DurationMs=%d, want 0 (absent)", timing.DurationMs)
	}
	if result.Obligations != nil {
		t.Fatalf("Obligations must remain nil, got %v", result.Obligations)
	}

	// Adapter-owned slice must not be sorted in place.
	if spy.lastDigest != wantDigest {
		t.Fatalf("adapter digest=%q, want %q", spy.lastDigest, wantDigest)
	}
	if got := spy.conclusion.ReasonCodes; len(got) != 3 || got[0] != "ZULU" {
		t.Fatalf("adapter conclusion mutated: %v", got)
	}
}

func TestEvaluate_SpyAtMostOnceNoRetry(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32
	spy := &evaluateSpyAdapter{
		err: errors.New("transient"),
		onCall: func(context.Context) {
			calls.Add(1)
		},
	}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	_, _, nr, err := b.Evaluate(context.Background(), validEvaluateRequest())
	if nr != NonResultAdapterFailure || !errors.Is(err, errAdapterFailed) {
		t.Fatalf("nr=%q err=%v", nr, err)
	}
	if calls.Load() != 1 || spy.callCount() != 1 {
		t.Fatalf("calls=%d spy=%d, want exactly 1 (no retry)", calls.Load(), spy.callCount())
	}
}

func TestEvaluate_NilBoundaryDependencyRejected(t *testing.T) {
	t.Parallel()

	ts := NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}

	if _, err := NewBoundary(nil, ts); !errors.Is(err, errDependencyInvalid) {
		t.Fatalf("nil adapter: err=%v, want errDependencyInvalid", err)
	}
	if _, err := NewBoundary(spy, nil); !errors.Is(err, errDependencyInvalid) {
		t.Fatalf("nil time source: err=%v, want errDependencyInvalid", err)
	}
	if _, err := NewBoundary(nil, nil); !errors.Is(err, errDependencyInvalid) {
		t.Fatalf("both nil: err=%v, want errDependencyInvalid", err)
	}
}

func TestEvaluate_TypedNilBoundaryDependencyRejected(t *testing.T) {
	t.Parallel()

	ts := NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	var typedNilAdapter *evaluateSpyAdapter
	var typedNilTS *ptrTimeSource
	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}

	if _, err := NewBoundary(typedNilAdapter, ts); !errors.Is(err, errDependencyInvalid) {
		t.Fatalf("typed-nil adapter: err=%v, want errDependencyInvalid", err)
	}
	if _, err := NewBoundary(spy, typedNilTS); !errors.Is(err, errDependencyInvalid) {
		t.Fatalf("typed-nil time source: err=%v, want errDependencyInvalid", err)
	}
	if _, err := NewBoundary(typedNilAdapter, typedNilTS); !errors.Is(err, errDependencyInvalid) {
		t.Fatalf("both typed-nil: err=%v, want errDependencyInvalid", err)
	}
}

func TestEvaluate_EveryNonResultReturnsZeroResultAndTiming(t *testing.T) {
	t.Parallel()

	fixed := NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	req := validEvaluateRequest()

	cases := []struct {
		name      string
		setup     func() (context.Context, *evaluateSpyAdapter, PolicyEvaluationRequest)
		wantNR    NonResult
		wantErr   error
		wantCalls int
	}{
		{
			name: "RequestInvalid",
			setup: func() (context.Context, *evaluateSpyAdapter, PolicyEvaluationRequest) {
				bad := req
				bad.Action = ""
				return context.Background(), &evaluateSpyAdapter{conclusion: validAllowConclusion()}, bad
			},
			wantNR:    NonResultRequestInvalid,
			wantErr:   errRequestInvalid,
			wantCalls: 0,
		},
		{
			name: "Canceled",
			setup: func() (context.Context, *evaluateSpyAdapter, PolicyEvaluationRequest) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, &evaluateSpyAdapter{conclusion: validAllowConclusion()}, req
			},
			wantNR:    NonResultCanceled,
			wantErr:   errCanceled,
			wantCalls: 0,
		},
		{
			name: "DeadlineExceeded",
			setup: func() (context.Context, *evaluateSpyAdapter, PolicyEvaluationRequest) {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Minute))
				t.Cleanup(cancel)
				return ctx, &evaluateSpyAdapter{conclusion: validAllowConclusion()}, req
			},
			wantNR:    NonResultDeadlineExceeded,
			wantErr:   errDeadlineExceeded,
			wantCalls: 0,
		},
		{
			name: "AdapterFailure",
			setup: func() (context.Context, *evaluateSpyAdapter, PolicyEvaluationRequest) {
				return context.Background(), &evaluateSpyAdapter{err: errors.New("boom")}, req
			},
			wantNR:    NonResultAdapterFailure,
			wantErr:   errAdapterFailed,
			wantCalls: 1,
		},
		{
			name: "InvalidAdapterResult",
			setup: func() (context.Context, *evaluateSpyAdapter, PolicyEvaluationRequest) {
				return context.Background(), &evaluateSpyAdapter{
					conclusion: Conclusion{Outcome: Outcome("nope"), ReasonCodes: []string{"OK"}},
				}, req
			},
			wantNR:    NonResultInvalidAdapterResult,
			wantErr:   errInvalidAdapterResult,
			wantCalls: 1,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, spy, r := tc.setup()
			b := mustBoundary(t, spy, fixed)
			result, timing, nr, err := b.Evaluate(ctx, r)
			if nr != tc.wantNR {
				t.Fatalf("NonResult=%q, want %q", nr, tc.wantNR)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("err=%v, want %v", err, tc.wantErr)
			}
			assertZeroResultTiming(t, result, timing)
			if spy.callCount() != tc.wantCalls {
				t.Fatalf("calls=%d, want %d", spy.callCount(), tc.wantCalls)
			}
		})
	}
}

func TestEvaluate_NormalizedInputSortedAndRequestIDExcluded(t *testing.T) {
	t.Parallel()

	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	req := PolicyEvaluationRequest{
		SubjectRef: apimeta.TypedRef{
			APIVersion: "services.sovrunn.io/v1alpha1",
			Kind:       "ServiceInstance",
			Name:       "reporting-api",
			UID:        "service-001",
		},
		Action: "service.read",
		ProfileRefs: []apimeta.TypedRef{
			{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "writer", UID: "role-002"},
			{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "reader", UID: "role-001"},
		},
		CandidateRefs: []apimeta.TypedRef{
			{APIVersion: "topology.sovrunn.io/v1alpha1", Kind: "CloudProvider", Name: "b", UID: "cp-b"},
			{APIVersion: "topology.sovrunn.io/v1alpha1", Kind: "CloudProvider", Name: "a", UID: "cp-a"},
		},
		RequestID: "corr-secret-should-not-reach-adapter",
	}

	_, _, nr, err := b.Evaluate(context.Background(), req)
	if err != nil || nr != "" {
		t.Fatalf("Evaluate: nr=%q err=%v", nr, err)
	}
	in := spy.lastIn
	if in.ProfileRefs[0].Name != "reader" || in.ProfileRefs[1].Name != "writer" {
		t.Fatalf("ProfileRefs not sorted: %+v", in.ProfileRefs)
	}
	if in.CandidateRefs[0].Name != "a" || in.CandidateRefs[1].Name != "b" {
		t.Fatalf("CandidateRefs not sorted: %+v", in.CandidateRefs)
	}
	if in.Action != "service.read" {
		t.Fatalf("Action=%q", in.Action)
	}
}

func TestEvaluate_ConcurrentNoSharedMutableState(t *testing.T) {
	t.Parallel()

	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
	b := mustBoundary(t, spy, NewFixedTimeSource(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))
	req := validEvaluateRequest()

	const goroutines = 32
	var wg sync.WaitGroup
	errs := make(chan error, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, timing, nr, err := b.Evaluate(context.Background(), req)
			if err != nil || nr != "" {
				errs <- errors.New("unexpected non-result")
				return
			}
			if result.Outcome != OutcomeAllow || timing.CompletedAt == "" {
				errs <- errors.New("unexpected success payload")
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
	if spy.callCount() != goroutines {
		t.Fatalf("calls=%d, want %d", spy.callCount(), goroutines)
	}
}

func TestEvaluate_ValidateConclusionHelper(t *testing.T) {
	t.Parallel()

	if err := validateConclusion(validAllowConclusion()); err != nil {
		t.Fatalf("valid conclusion: %v", err)
	}
	if err := validateConclusion(Conclusion{Outcome: OutcomeAllow, ReasonCodes: nil}); err == nil {
		t.Fatal("nil reason codes must fail")
	}
	if err := validateConclusion(Conclusion{Outcome: Outcome("x"), ReasonCodes: []string{"A"}}); err == nil {
		t.Fatal("invalid outcome must fail")
	}
}
