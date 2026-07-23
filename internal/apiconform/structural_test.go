package apiconform

import (
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

const (
	testSchemaID       = "api/schemas/example.json"
	testCommonTypeMeta = "api/schemas/_common/type-meta.json"
)

func testStructuralValidator(t *testing.T, schemas map[string][]byte) *StructuralValidator {
	t.Helper()
	reg, err := NewMemorySchemaRegistry(schemas)
	if err != nil {
		t.Fatalf("NewMemorySchemaRegistry: %v", err)
	}
	resolver, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}
	cfg, err := NewStructuralValidatorConfig(reg, resolver)
	if err != nil {
		t.Fatalf("NewStructuralValidatorConfig: %v", err)
	}
	v, err := NewStructuralValidator(cfg)
	if err != nil {
		t.Fatalf("NewStructuralValidator: %v", err)
	}
	return v
}

func testExampleSchemas() map[string][]byte {
	return map[string][]byte{
		testSchemaID: []byte(`{
			"type": "object",
			"properties": {
				"name": { "type": "string", "minLength": 1 },
				"meta": { "$ref": "_common/type-meta.json" }
			},
			"required": ["name", "meta"],
			"additionalProperties": false
		}`),
		testCommonTypeMeta: []byte(`{
			"type": "object",
			"properties": {
				"apiVersion": { "type": "string", "minLength": 1 },
				"kind": { "type": "string", "minLength": 1 }
			},
			"required": ["apiVersion", "kind"],
			"additionalProperties": false
		}`),
	}
}

func hasViolation(vs []apiproblem.Violation, code, field string) bool {
	for _, v := range vs {
		if string(v.Code) == code && v.Field == field {
			return true
		}
	}
	return false
}

func TestNewStructuralValidator_NilRegistryRejected(t *testing.T) {
	t.Parallel()

	_, err := NewStructuralValidator(StructuralValidatorConfig{})
	if !errors.Is(err, ErrStructuralValidator) {
		t.Fatalf("zero config: err=%v, want ErrStructuralValidator", err)
	}
}

func TestNewStructuralValidator_NilResolverRejected(t *testing.T) {
	t.Parallel()

	reg, err := NewMemorySchemaRegistry(map[string][]byte{
		testCommonTypeMeta: []byte(`{"type":"object"}`),
	})
	if err != nil {
		t.Fatalf("NewMemorySchemaRegistry: %v", err)
	}
	// Bypass NewStructuralValidatorConfig so we can isolate nil-resolver
	// rejection on the adapter constructor itself.
	cfg := StructuralValidatorConfig{registry: reg, resolver: nil}
	_, err = NewStructuralValidator(cfg)
	if !errors.Is(err, ErrStructuralValidator) {
		t.Fatalf("nil resolver: err=%v, want ErrStructuralValidator", err)
	}
}

func TestStructuralValidator_AcceptsValidInstance(t *testing.T) {
	t.Parallel()

	v := testStructuralValidator(t, testExampleSchemas())
	instance := map[string]any{
		"name": "demo",
		"meta": map[string]any{
			"apiVersion": "platform.sovrunn.io/v1",
			"kind":       "Project",
		},
	}

	violations, err := v.Validate(instance, testSchemaID)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("valid instance: violations=%#v, want none", violations)
	}
}

func TestStructuralValidator_RejectsInvalidInstanceWithViolations(t *testing.T) {
	t.Parallel()

	v := testStructuralValidator(t, testExampleSchemas())
	instance := map[string]any{
		"meta": map[string]any{
			"apiVersion": "platform.sovrunn.io/v1",
			"kind":       "Project",
		},
	}

	violations, err := v.Validate(instance, testSchemaID)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if !hasViolation(violations, apischema.CodeRequiredField, "/name") {
		t.Fatalf("want %s at /name, got %#v", apischema.CodeRequiredField, violations)
	}
}

func TestStructuralValidator_MissingSchemaReturnsError(t *testing.T) {
	t.Parallel()

	v := testStructuralValidator(t, testExampleSchemas())
	violations, err := v.Validate(map[string]any{"name": "demo"}, "api/schemas/missing.json")
	if !errors.Is(err, ErrStructuralValidator) {
		t.Fatalf("missing schema: err=%v, want ErrStructuralValidator", err)
	}
	if !errors.Is(err, ErrSchemaNotFound) {
		t.Fatalf("missing schema: err=%v, want ErrSchemaNotFound cause", err)
	}
	if violations != nil {
		t.Fatalf("missing schema: violations=%#v, want nil", violations)
	}
}

func TestStructuralValidator_NilRegistryOnValidateReturnsError(t *testing.T) {
	t.Parallel()

	// Zero-value adapter (not constructed via NewStructuralValidator) must
	// fail closed rather than panic or silently skip structural checks.
	var v StructuralValidator
	violations, err := v.Validate(map[string]any{"name": "demo"}, testSchemaID)
	if !errors.Is(err, ErrStructuralValidator) {
		t.Fatalf("nil registry validate: err=%v, want ErrStructuralValidator", err)
	}
	if violations != nil {
		t.Fatalf("nil registry validate: violations=%#v, want nil", violations)
	}
}

func TestStructuralValidator_RefResolutionFailureReturnsError(t *testing.T) {
	t.Parallel()

	schemas := map[string][]byte{
		testSchemaID: []byte(`{
			"type": "object",
			"properties": {
				"meta": { "$ref": "_common/missing.json" }
			}
		}`),
	}
	v := testStructuralValidator(t, schemas)
	violations, err := v.Validate(map[string]any{"meta": map[string]any{}}, testSchemaID)
	if !errors.Is(err, ErrStructuralValidator) {
		t.Fatalf("ref failure: err=%v, want ErrStructuralValidator", err)
	}
	if !errors.Is(err, ErrSchemaNotFound) && !errors.Is(err, ErrRefRejected) {
		t.Fatalf("ref failure: err=%v, want ErrSchemaNotFound or ErrRefRejected cause", err)
	}
	if violations != nil {
		t.Fatalf("ref failure: violations=%#v, want nil", violations)
	}
}

func TestStructuralValidator_UnsupportedSchemaReturnsViolations(t *testing.T) {
	t.Parallel()

	schemas := map[string][]byte{
		testSchemaID: []byte(`{"oneOf":[{"type":"string"}]}`),
	}
	v := testStructuralValidator(t, schemas)
	violations, err := v.Validate("x", testSchemaID)
	if err != nil {
		t.Fatalf("unsupported schema must be ordinary violations, err=%v", err)
	}
	if !hasViolation(violations, apischema.CodeUnsupportedKeyword, "/oneOf") {
		t.Fatalf("want %s at /oneOf, got %#v", apischema.CodeUnsupportedKeyword, violations)
	}
}

func TestStructuralValidator_NilReceiverReturnsError(t *testing.T) {
	t.Parallel()

	var v *StructuralValidator
	violations, err := v.Validate(map[string]any{}, testSchemaID)
	if !errors.Is(err, ErrStructuralValidator) {
		t.Fatalf("nil receiver: err=%v, want ErrStructuralValidator", err)
	}
	if violations != nil {
		t.Fatalf("nil receiver: violations=%#v, want nil", violations)
	}
}

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
