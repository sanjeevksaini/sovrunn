package policyeval

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

// Sanitized evaluation-operation diagnostics (D-08 / design §5.4).
// errRequestInvalid is defined in request.go with the same RequestInvalid text.
var (
	errDependencyInvalid    = errors.New("policy evaluation dependency invalid")
	errCanceled             = errors.New("policy evaluation canceled")
	errDeadlineExceeded     = errors.New("policy evaluation deadline exceeded")
	errAdapterFailed        = errors.New("policy evaluation adapter failed")
	errInvalidAdapterResult = errors.New("policy evaluation adapter result invalid")
)

// Boundary is the constructed evaluation operation around exactly one adapter
// and one time source (D-09 / CDG-F17-03).
type Boundary struct {
	adapter PolicyEngineAdapter
	ts      TimeSource
}

// NewBoundary captures adapter and ts. It rejects nil and typed-nil
// dependencies and never returns a partially usable boundary.
func NewBoundary(adapter PolicyEngineAdapter, ts TimeSource) (*Boundary, error) {
	if isNilDependency(adapter) || isNilDependency(ts) {
		return nil, errDependencyInvalid
	}
	return &Boundary{adapter: adapter, ts: ts}, nil
}

// Evaluate performs the closed nine-step precedence (CDG-F17-06 / D-11):
// (1) pre-check ctx.Err(); (2) request validation; (3) digest;
// (4) pre-invocation ctx.Err(); (5) capture startedAt and direct-call adapter;
// (6) post-return ctx.Err() before inspecting output; (7) adapter error;
// (8) conclusion validation; (9) capture completedAt, sort reason codes,
// construct result and timing.
//
// Success returns a populated result, timing envelope, empty NonResult, and
// nil error. Any normalized non-result returns zero result, zero timing, the
// category, and a fixed sanitized diagnostic.
func (b *Boundary) Evaluate(
	ctx context.Context,
	req PolicyEvaluationRequest,
) (PolicyEvaluationResult, decision.TimingEnvelope, NonResult, error) {
	var zeroResult PolicyEvaluationResult
	var zeroTiming decision.TimingEnvelope

	// Nil context is an operation-precondition failure (D-08 / §4.5).
	if ctx == nil {
		return zeroResult, zeroTiming, NonResultRequestInvalid, errRequestInvalid
	}

	// (1) Already-observed cancellation or expired deadline.
	if nr, err := classifyContextErr(ctx.Err()); err != nil {
		return zeroResult, zeroTiming, nr, err
	}

	// (2) Request structural / reference / bound validation.
	if err := req.Validate(); err != nil {
		return zeroResult, zeroTiming, NonResultRequestInvalid, errRequestInvalid
	}

	// (3) Canonicalization and SHA-256 digest.
	digest, err := digestPolicyEvaluationRequest(req)
	if err != nil {
		return zeroResult, zeroTiming, NonResultRequestInvalid, errRequestInvalid
	}

	// (4) Cancellation/deadline recheck immediately before invocation.
	if nr, err := classifyContextErr(ctx.Err()); err != nil {
		return zeroResult, zeroTiming, nr, err
	}

	// (5) Capture startedAt and perform exactly one direct adapter call.
	startedAt := b.ts.NowUTC().UTC()
	in := normalizedInputFromRequest(req)
	conclusion, adapterErr := b.adapter.Evaluate(ctx, in, digest)

	// (6) Post-return ctx.Err() before inspecting adapter output; discard late output.
	if nr, err := classifyContextErr(ctx.Err()); err != nil {
		return zeroResult, zeroTiming, nr, err
	}

	// (7) Non-nil adapter error → AdapterFailure; raw error discarded.
	if adapterErr != nil {
		return zeroResult, zeroTiming, NonResultAdapterFailure, errAdapterFailed
	}

	// (8) Copy conclusion, then validate outcome and reason codes.
	owned := copyConclusion(conclusion)
	if err := validateConclusion(owned); err != nil {
		return zeroResult, zeroTiming, NonResultInvalidAdapterResult, errInvalidAdapterResult
	}

	// (9) Capture completedAt, sort owned reason codes, construct result+timing.
	completedAt := b.ts.NowUTC().UTC()
	sort.Strings(owned.ReasonCodes)
	completedAtStr := completedAt.Format(time.RFC3339Nano)
	result := PolicyEvaluationResult{
		Outcome:     owned.Outcome,
		ReasonCodes: owned.ReasonCodes,
		InputDigest: digest,
		EvaluatedAt: completedAtStr,
	}
	timing := decision.TimingEnvelope{
		StartedAt:   startedAt.Format(time.RFC3339Nano),
		CompletedAt: completedAtStr,
	}
	return result, timing, "", nil
}

// validateConclusion is the single package-private source of truth for
// conclusion validation (outcome membership; reason-code grammar, count 1..32,
// uniqueness). Reused by the boundary (step 8) and the deterministic fake
// constructor (TASK-F17-07).
func validateConclusion(c Conclusion) error {
	if !c.Outcome.Valid() {
		return errInvalidAdapterResult
	}
	n := len(c.ReasonCodes)
	if n < MinReasonCodes || n > MaxReasonCodes {
		return errInvalidAdapterResult
	}
	seen := make(map[string]struct{}, n)
	for _, code := range c.ReasonCodes {
		if !ReasonCodeRegexp.MatchString(code) {
			return errInvalidAdapterResult
		}
		if _, dup := seen[code]; dup {
			return errInvalidAdapterResult
		}
		seen[code] = struct{}{}
	}
	return nil
}

func classifyContextErr(err error) (NonResult, error) {
	if err == nil {
		return "", nil
	}
	// DeadlineExceeded first: it must not be collapsed into Canceled.
	if errors.Is(err, context.DeadlineExceeded) {
		return NonResultDeadlineExceeded, errDeadlineExceeded
	}
	if errors.Is(err, context.Canceled) {
		return NonResultCanceled, errCanceled
	}
	// Fail closed for unexpected ctx.Err() values.
	return NonResultCanceled, errCanceled
}

func isNilDependency(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
		return rv.IsNil()
	default:
		return false
	}
}

func copyConclusion(c Conclusion) Conclusion {
	out := Conclusion{Outcome: c.Outcome}
	if c.ReasonCodes != nil {
		out.ReasonCodes = append([]string(nil), c.ReasonCodes...)
	}
	return out
}

// normalizedInputFromRequest builds a deep-copied, sorted NormalizedInput
// snapshot for adapter handoff. RequestID is intentionally excluded.
func normalizedInputFromRequest(req PolicyEvaluationRequest) NormalizedInput {
	snap := snapshotPolicyEvaluationRequest(req)
	sortTypedRefs(snap.ProfileRefs)
	sortTypedRefs(snap.CandidateRefs)

	in := NormalizedInput{
		SubjectRef:    snap.SubjectRef,
		Action:        snap.Action,
		ProfileRefs:   snap.ProfileRefs,
		CandidateRefs: snap.CandidateRefs,
	}
	if snap.ContextRef != nil {
		cp := *snap.ContextRef
		in.ContextRef = &cp
	}
	return in
}
