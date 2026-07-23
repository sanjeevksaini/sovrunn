package apiconform

import (
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

// SchemaIssuesToViolations translates package-local apischema.SchemaIssue
// values into apiproblem.Violation values (F12-VALIDATION-006, D-01a, D-02).
//
// Mapping is field-for-field:
//
//	SchemaIssue.Path    → Violation.Field  (RFC 6901 JSON Pointer)
//	SchemaIssue.Code    → Violation.Code   (stable machine contract)
//	SchemaIssue.Message → Violation.Message (informational; must not carry secrets)
//
// The returned slice is newly allocated so callers cannot mutate the input
// through the result. A nil or empty input yields nil.
//
// Translation lives in apiconform (not apischema) so apischema never imports
// apiproblem, preserving the D-02 import-direction boundary.
func SchemaIssuesToViolations(issues []apischema.SchemaIssue) []apiproblem.Violation {
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
