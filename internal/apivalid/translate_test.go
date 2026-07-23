package apivalid

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

func TestRefIssuesToViolationsMapsPathCodeMessage(t *testing.T) {
	issues := []apiref.RefIssue{
		{
			Path:    "/spec/resourcePoolRef/kind",
			Code:    apiref.CodeKindNotAllowed,
			Message: "reference kind is not in the allowed set for this field",
		},
		{
			Path:    "/metadata/scopeRef",
			Code:    apiref.CodeNameUIDMismatch,
			Message: "reference name and uid must identify the same object",
		},
	}

	got := RefIssuesToViolations(issues)
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

	if got[0].Code != apiproblem.ViolationCode(apiref.CodeKindNotAllowed) {
		t.Fatalf("first code type = %T value %q", got[0].Code, got[0].Code)
	}
}

func TestRefIssuesToViolationsNilAndEmpty(t *testing.T) {
	if got := RefIssuesToViolations(nil); got != nil {
		t.Fatalf("nil input: got %#v, want nil", got)
	}
	if got := RefIssuesToViolations([]apiref.RefIssue{}); got != nil {
		t.Fatalf("empty input: got %#v, want nil", got)
	}
}

func TestRefIssuesToViolationsCopiesSlice(t *testing.T) {
	issues := []apiref.RefIssue{{
		Path:    "/spec/refs/0",
		Code:    apiref.CodeProviderNativeID,
		Message: "provider-native identifiers must not act as core references",
	}}
	got := RefIssuesToViolations(issues)
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}

	issues[0].Path = "/mutated"
	issues[0].Code = "MUTATED"
	issues[0].Message = "mutated"

	if got[0].Field != "/spec/refs/0" {
		t.Fatalf("mutation leaked into Field: %q", got[0].Field)
	}
	if string(got[0].Code) != apiref.CodeProviderNativeID {
		t.Fatalf("mutation leaked into Code: %q", got[0].Code)
	}
	if got[0].Message != "provider-native identifiers must not act as core references" {
		t.Fatalf("mutation leaked into Message: %q", got[0].Message)
	}
}

func TestRefIssuesToViolationsPreservesAllStableCodes(t *testing.T) {
	codes := []string{
		apiref.CodeKindNotAllowed,
		apiref.CodeScopeNotAllowed,
		apiref.CodeDirectionInvalid,
		apiref.CodeNameUIDMismatch,
		apiref.CodeProviderNativeID,
		apiref.CodeMissingAPIVersion,
		apiref.CodeMissingKind,
		apiref.CodeMissingName,
		apiref.CodeRefsExceedLimit,
	}
	issues := make([]apiref.RefIssue, len(codes))
	for i, code := range codes {
		issues[i] = apiref.RefIssue{
			Path:    "/spec/ref",
			Code:    code,
			Message: "stable code " + code,
		}
	}

	got := RefIssuesToViolations(issues)
	if len(got) != len(codes) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(codes))
	}
	for i, code := range codes {
		if string(got[i].Code) != code {
			t.Fatalf("got[%d].Code = %q, want %q", i, got[i].Code, code)
		}
		if got[i].Field != "/spec/ref" {
			t.Fatalf("got[%d].Field = %q, want /spec/ref", i, got[i].Field)
		}
	}
}
