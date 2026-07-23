package apiconform

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

func TestSchemaIssuesToViolationsMapsPathCodeMessage(t *testing.T) {
	issues := []apischema.SchemaIssue{
		{
			Path:    "/spec/name",
			Code:    apischema.CodeRequiredField,
			Message: "required property is missing",
		},
		{
			Path:    "/spec/sizeGiB",
			Code:    apischema.CodeTypeMismatch,
			Message: "value type does not match schema type",
		},
	}

	got := SchemaIssuesToViolations(issues)
	if len(got) != len(issues) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(issues))
	}

	for i, issue := range issues {
		v := got[i]
		if v.Field != issue.Path {
			t.Fatalf("got[%d].Field = %q, want %q", i, v.Field, issue.Path)
		}
		if string(v.Code) != issue.Code {
			t.Fatalf("got[%d].Code = %q, want %q", i, v.Code, issue.Code)
		}
		if v.Message != issue.Message {
			t.Fatalf("got[%d].Message = %q, want %q", i, v.Message, issue.Message)
		}
	}

	if got[0].Code != apiproblem.ViolationCode(apischema.CodeRequiredField) {
		t.Fatalf("first code type = %T value %q", got[0].Code, got[0].Code)
	}
}

func TestSchemaIssuesToViolationsNilAndEmpty(t *testing.T) {
	if got := SchemaIssuesToViolations(nil); got != nil {
		t.Fatalf("nil input: got %#v, want nil", got)
	}
	if got := SchemaIssuesToViolations([]apischema.SchemaIssue{}); got != nil {
		t.Fatalf("empty input: got %#v, want nil", got)
	}
}

func TestSchemaIssuesToViolationsCopiesSlice(t *testing.T) {
	issues := []apischema.SchemaIssue{{
		Path:    "/status/phase",
		Code:    apischema.CodeUnknownField,
		Message: "additional property is not allowed",
	}}
	got := SchemaIssuesToViolations(issues)
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}

	issues[0].Path = "/mutated"
	issues[0].Code = "MUTATED"
	issues[0].Message = "mutated"

	if got[0].Field != "/status/phase" {
		t.Fatalf("mutation leaked into Field: %q", got[0].Field)
	}
	if string(got[0].Code) != apischema.CodeUnknownField {
		t.Fatalf("mutation leaked into Code: %q", got[0].Code)
	}
	if got[0].Message != "additional property is not allowed" {
		t.Fatalf("mutation leaked into Message: %q", got[0].Message)
	}
}

func TestSchemaIssuesToViolationsPreservesAllStableCodes(t *testing.T) {
	codes := []string{
		apischema.CodeUnsupportedKeyword,
		apischema.CodeMalformedSchema,
		apischema.CodeUnresolvedRef,
		apischema.CodeSchemaFalse,
		apischema.CodeTypeMismatch,
		apischema.CodeRequiredField,
		apischema.CodeEnumMismatch,
		apischema.CodeOutOfRange,
		apischema.CodePatternMismatch,
		apischema.CodeUnknownField,
	}
	issues := make([]apischema.SchemaIssue, len(codes))
	for i, code := range codes {
		issues[i] = apischema.SchemaIssue{
			Path:    "/spec/field",
			Code:    code,
			Message: "stable code " + code,
		}
	}

	got := SchemaIssuesToViolations(issues)
	if len(got) != len(codes) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(codes))
	}
	for i, code := range codes {
		if string(got[i].Code) != code {
			t.Fatalf("got[%d].Code = %q, want %q", i, got[i].Code, code)
		}
		if got[i].Field != "/spec/field" {
			t.Fatalf("got[%d].Field = %q, want /spec/field", i, got[i].Field)
		}
	}
}
