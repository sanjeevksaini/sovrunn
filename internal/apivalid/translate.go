package apivalid

import (
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// RefIssuesToViolations translates package-local apiref.RefIssue values into
// apiproblem.Violation values (F12-VALIDATION-006, D-02, D-04).
//
// Mapping is field-for-field:
//
//	RefIssue.Path    → Violation.Field  (RFC 6901 JSON Pointer)
//	RefIssue.Code    → Violation.Code   (stable machine contract)
//	RefIssue.Message → Violation.Message (informational; must not carry secrets)
//
// The returned slice is newly allocated so callers cannot mutate the input
// through the result. A nil or empty input yields nil.
func RefIssuesToViolations(issues []apiref.RefIssue) []apiproblem.Violation {
	if len(issues) == 0 {
		return nil
	}
	out := make([]apiproblem.Violation, len(issues))
	for i, issue := range issues {
		out[i] = apiproblem.Violation{
			Field:   issue.Path,
			Code:    apiproblem.ViolationCode(issue.Code),
			Message: issue.Message,
		}
	}
	return out
}
