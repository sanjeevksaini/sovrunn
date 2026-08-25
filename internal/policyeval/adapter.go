package policyeval

import (
	"context"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// NormalizedInput is the semantic adapter input snapshot derived from a
// validated PolicyEvaluationRequest. RequestID is intentionally absent: it is
// not a behavior input to the adapter (CDG-F17-03 / D-09).
//
// ProfileRefs and CandidateRefs are sorted by normalized-reference identity.
// CandidateRefs is an empty (non-nil) slice when the request has no candidates.
// ContextRef is non-nil only when the request supplied a context reference.
type NormalizedInput struct {
	SubjectRef    apimeta.TypedRef
	Action        string
	ContextRef    *apimeta.TypedRef
	ProfileRefs   []apimeta.TypedRef
	CandidateRefs []apimeta.TypedRef
}

// Conclusion is the adapter-owned evaluation conclusion. It carries only
// outcome and reason codes—no digest, time, timing, obligations,
// DecisionRecord identity, or AuditEvent identity (D-09 / CDG-F17-03).
//
// YAML tags support the test-only fixture decoder (D-13); they are not a
// production input surface.
type Conclusion struct {
	Outcome     Outcome  `yaml:"outcome"`
	ReasonCodes []string `yaml:"reasonCodes"`
}

// PolicyEngineAdapter is the engine-neutral in-process port behind the
// evaluation boundary (D-09 / CDG-F17-03).
//
// A non-nil error returned by Evaluate is only the internal signal for the
// normalized NonResult category AdapterFailure. The evaluation boundary
// never wraps, returns, or logs that error text. Phase 2R adapter output
// contains no obligations; the boundary constructs the complete
// PolicyEvaluationResult and always leaves Obligations absent.
//
// A conforming adapter honors cancellation and deadlines, eventually returns,
// and is safe for concurrent use. It transfers ownership of the returned
// Conclusion and must not mutate its reason-code slice concurrently or after
// return.
type PolicyEngineAdapter interface {
	Evaluate(ctx context.Context, in NormalizedInput, inputDigest string) (Conclusion, error)
}
