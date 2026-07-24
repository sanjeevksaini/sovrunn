package apiproblem

import (
	"strings"
)

// Problem is an RFC 9457 Problem Details response with Sovrunn extensions
// (F12-ERROR-001, D-05).
//
// Stable machine contracts: Type, Code, and Violation.Code.
// Human-readable Detail and Violation.Message MAY evolve and MUST NOT be
// parsed by clients (F12-ERROR-003).
//
// Responses MUST NOT expose credentials, stack traces, raw provider errors,
// sensitive policy inputs, or inaccessible resource details (F12-ERROR-004).
type Problem struct {
	Type       string      `json:"type"`
	Title      string      `json:"title"`
	Status     int         `json:"status"`
	Detail     string      `json:"detail,omitempty"`
	Instance   string      `json:"instance,omitempty"`
	Code       ErrorCode   `json:"code"`
	RequestID  string      `json:"requestId,omitempty"`
	Violations []Violation `json:"violations,omitempty"`
}

// Violation identifies one invalid field by RFC 6901 JSON Pointer
// (F12-ERROR-001, F12-VALIDATION-006).
type Violation struct {
	Field   string        `json:"field"`   // RFC 6901 JSON Pointer, e.g. /spec/storage/sizeGiB
	Code    ViolationCode `json:"code"`    // stable violation code
	Message string        `json:"message"` // human-readable; redacted of sensitive detail
}

// MediaTypeProblemJSON is the RFC 9457 recommended media type for Problem
// Details responses. HTTP adapters set this; the grammar package itself has
// no HTTP dependency beyond status integers in httpmap.
const MediaTypeProblemJSON = "application/problem+json"

// TypeURNPrefix is the stable URI prefix for Sovrunn problem type URNs.
const TypeURNPrefix = "urn:sovrunn:problem:"

// TypeURN returns the stable RFC 9457 type URI for code
// (e.g. VALIDATION_FAILED → urn:sovrunn:problem:validation-failed).
func TypeURN(code ErrorCode) string {
	return TypeURNPrefix + codeToKebab(string(code))
}

// New builds a Problem for code using the baseline title and HTTP status
// from the F12-ERROR-002 registry. Detail, Instance, RequestID, and
// Violations are left empty for the caller to set.
func New(code ErrorCode) *Problem {
	return &Problem{
		Type:   TypeURN(code),
		Title:  TitleFor(code),
		Status: StatusForCode(code),
		Code:   code,
	}
}

// WithDetail returns a shallow copy of p with Detail set.
func (p *Problem) WithDetail(detail string) *Problem {
	if p == nil {
		return nil
	}
	out := *p
	out.Detail = detail
	return &out
}

// WithRequestID returns a shallow copy of p with RequestID set for
// request/operation correlation (observability baseline).
func (p *Problem) WithRequestID(requestID string) *Problem {
	if p == nil {
		return nil
	}
	out := *p
	out.RequestID = requestID
	return &out
}

// WithInstance returns a shallow copy of p with Instance set.
func (p *Problem) WithInstance(instance string) *Problem {
	if p == nil {
		return nil
	}
	out := *p
	out.Instance = instance
	return &out
}

// WithViolations returns a shallow copy of p with Violations set.
// The violations slice is copied so callers cannot mutate the Problem's
// backing array through the input slice.
func (p *Problem) WithViolations(violations []Violation) *Problem {
	if p == nil {
		return nil
	}
	out := *p
	if len(violations) == 0 {
		out.Violations = nil
		return &out
	}
	out.Violations = make([]Violation, len(violations))
	copy(out.Violations, violations)
	return &out
}

func codeToKebab(code string) string {
	return strings.ToLower(strings.ReplaceAll(code, "_", "-"))
}
