package policyeval

import (
	"encoding/json"
)

// Outcome is the closed PolicyEvaluationResult outcome vocabulary (CDG-F17-03).
// Exactly four values are valid; a fifth outcome is never introduced via NonResult.
type Outcome string

const (
	OutcomeAllow            Outcome = "Allow"
	OutcomeDeny             Outcome = "Deny"
	OutcomeIndeterminate    Outcome = "Indeterminate"
	OutcomeRequiresApproval Outcome = "RequiresApproval"
)

// AllOutcomes returns the closed Outcome set in stable order.
func AllOutcomes() []Outcome {
	return []Outcome{
		OutcomeAllow,
		OutcomeDeny,
		OutcomeIndeterminate,
		OutcomeRequiresApproval,
	}
}

// Valid reports whether o is one of the four closed Outcome values.
func (o Outcome) Valid() bool {
	switch o {
	case OutcomeAllow, OutcomeDeny, OutcomeIndeterminate, OutcomeRequiresApproval:
		return true
	default:
		return false
	}
}

// Reason-code grammar and count bounds (CDG-F17-03 / VS0-SCHEMA-019).
// Each reason code is 1..63 characters matching ReasonCodePattern.
// A result carries 1..32 unique lexically sorted codes. Validation of count,
// uniqueness, and grammar is owned by the evaluation boundary (TASK-F17-06).
const (
	ReasonCodePattern = `^[A-Z][A-Z0-9_]{0,62}$`
	MinReasonCodes    = 1
	MaxReasonCodes    = 32
	MaxReasonCodeLen  = 63
)

// reasonCodeMatcher implements the approved ASCII-only reason-code grammar
// without adding regexp to the closed FEATURE-0017 production import surface.
type reasonCodeMatcher struct{}

// MatchString reports whether code matches ^[A-Z][A-Z0-9_]{0,62}$.
func (reasonCodeMatcher) MatchString(code string) bool {
	if len(code) == 0 || len(code) > MaxReasonCodeLen {
		return false
	}
	if code[0] < 'A' || code[0] > 'Z' {
		return false
	}
	for i := 1; i < len(code); i++ {
		c := code[i]
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' {
			return false
		}
	}
	return true
}

// String returns the canonical grammar for diagnostics and parity tests.
func (reasonCodeMatcher) String() string {
	return ReasonCodePattern
}

// ReasonCodeRegexp retains the existing matcher surface while using the
// feature-owned deterministic ASCII implementation above.
var ReasonCodeRegexp reasonCodeMatcher

// PolicyEvaluationResult is the FEATURE-0017 engine-facing evaluation output
// envelope (VS0-SCHEMA-019). It is a TransientRequestResult-profile value with
// no TypeMeta, ObjectMeta, identity, scope, or status.
//
// Phase 2R fake output never produces obligations: Obligations remains nil so
// the field is omitted from JSON serialization (AC-F17-18). Typed obligation
// semantics remain deferred.
type PolicyEvaluationResult struct {
	Outcome     Outcome           `json:"outcome"`
	ReasonCodes []string          `json:"reasonCodes"`
	InputDigest string            `json:"inputDigest"`
	EvaluatedAt string            `json:"evaluatedAt"`
	Obligations []json.RawMessage `json:"obligations,omitempty"`
}

// NonResult is the closed internal evaluation-operation failure category
// (D-08 / CDG-F17-03). Values are never FEATURE-0012 Problem codes, Problem
// violation IDs, HTTP statuses, or a fifth Outcome.
type NonResult string

const (
	NonResultRequestInvalid       NonResult = "RequestInvalid"
	NonResultCanceled             NonResult = "Canceled"
	NonResultDeadlineExceeded     NonResult = "DeadlineExceeded"
	NonResultAdapterFailure       NonResult = "AdapterFailure"
	NonResultInvalidAdapterResult NonResult = "InvalidAdapterResult"
)

// AllNonResults returns the closed NonResult set in stable order.
func AllNonResults() []NonResult {
	return []NonResult{
		NonResultRequestInvalid,
		NonResultCanceled,
		NonResultDeadlineExceeded,
		NonResultAdapterFailure,
		NonResultInvalidAdapterResult,
	}
}

// Valid reports whether n is one of the five closed NonResult categories.
func (n NonResult) Valid() bool {
	switch n {
	case NonResultRequestInvalid,
		NonResultCanceled,
		NonResultDeadlineExceeded,
		NonResultAdapterFailure,
		NonResultInvalidAdapterResult:
		return true
	default:
		return false
	}
}
