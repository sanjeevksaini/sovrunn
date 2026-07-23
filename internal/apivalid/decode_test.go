package apivalid

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// permissivePolicy accepts status, system-owned metadata, and spec so
// duplicate/unknown-field tests are not confounded by FieldPolicy.
var permissivePolicy = FieldPolicy{
	Mode:              ModeReadRepresentation,
	AllowStatus:       true,
	AllowSystemOwned:  true,
	AllowSpecMutation: true,
}

var testLimits = Limits{
	MaxObjectBytes:  1 << 20,
	MaxNestingDepth: 32,
}

type decodeSample struct {
	APIVersion string         `json:"apiVersion"`
	Kind       string         `json:"kind"`
	Metadata   decodeMeta     `json:"metadata"`
	Spec       map[string]any `json:"spec"`
	Status     map[string]any `json:"status,omitempty"`
}

type decodeMeta struct {
	Name            string `json:"name"`
	UID             string `json:"uid,omitempty"`
	Generation      int64  `json:"generation,omitempty"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
	CreatedAt       string `json:"createdAt,omitempty"`
	UpdatedAt       string `json:"updatedAt,omitempty"`
}

func TestDecodeJSONHappyPath(t *testing.T) {
	const raw = `{"apiVersion":"platform.sovrunn.io/v1","kind":"Project","metadata":{"name":"demo"},"spec":{"displayName":"Demo"}}`
	var dst decodeSample
	if prob := DecodeJSON([]byte(raw), testLimits, permissivePolicy, &dst); prob != nil {
		t.Fatalf("DecodeJSON returned problem: %#v", prob)
	}
	if dst.Kind != "Project" || dst.Metadata.Name != "demo" {
		t.Fatalf("unexpected decode result: %#v", dst)
	}
}

func TestDecodeJSONRejectsDuplicateKey(t *testing.T) {
	// Duplicate "name" under metadata — encoding/json would keep the last
	// value; the token-scan detector must reject with DUPLICATE_FIELD.
	const raw = `{"apiVersion":"v1","kind":"Project","metadata":{"name":"a","name":"b"},"spec":{}}`
	var dst decodeSample
	prob := DecodeJSON([]byte(raw), testLimits, permissivePolicy, &dst)
	if prob == nil {
		t.Fatal("expected duplicate-key Problem, got nil")
	}
	if prob.Code != apiproblem.CodeDuplicateField {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeDuplicateField)
	}
	if len(prob.Violations) != 1 {
		t.Fatalf("Violations = %#v, want one entry", prob.Violations)
	}
	if prob.Violations[0].Field != "/metadata/name" {
		t.Fatalf("Field = %q, want /metadata/name", prob.Violations[0].Field)
	}
	if prob.Violations[0].Code != apiproblem.ViolationDuplicateField {
		t.Fatalf("Violation.Code = %q, want %q", prob.Violations[0].Code, apiproblem.ViolationDuplicateField)
	}
}

func TestDecodeJSONRejectsDuplicateKeyAtRoot(t *testing.T) {
	const raw = `{"apiVersion":"v1","kind":"Project","kind":"Tenant","metadata":{"name":"x"},"spec":{}}`
	var dst decodeSample
	prob := DecodeJSON([]byte(raw), testLimits, permissivePolicy, &dst)
	if prob == nil {
		t.Fatal("expected duplicate-key Problem, got nil")
	}
	if prob.Code != apiproblem.CodeDuplicateField {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeDuplicateField)
	}
	if got := prob.Violations[0].Field; got != "/kind" {
		t.Fatalf("Field = %q, want /kind", got)
	}
}

func TestDecodeJSONRejectsUnknownField(t *testing.T) {
	const raw = `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{},"extraField":true}`
	var dst decodeSample
	prob := DecodeJSON([]byte(raw), testLimits, permissivePolicy, &dst)
	if prob == nil {
		t.Fatal("expected unknown-field Problem, got nil")
	}
	if prob.Code != apiproblem.CodeUnknownField {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeUnknownField)
	}
	if len(prob.Violations) != 1 {
		t.Fatalf("Violations = %#v, want one entry", prob.Violations)
	}
	if prob.Violations[0].Field != "/extraField" {
		t.Fatalf("Field = %q, want /extraField", prob.Violations[0].Field)
	}
	if prob.Violations[0].Code != apiproblem.ViolationUnknownField {
		t.Fatalf("Violation.Code = %q, want %q", prob.Violations[0].Code, apiproblem.ViolationUnknownField)
	}
}

func TestDecodeJSONStableCodeAndJSONPointer(t *testing.T) {
	cases := []struct {
		name      string
		raw       string
		wantCode  apiproblem.ErrorCode
		wantField string
	}{
		{
			name:      "duplicate",
			raw:       `{"apiVersion":"v1","apiVersion":"v2","kind":"Project","metadata":{"name":"x"},"spec":{}}`,
			wantCode:  apiproblem.CodeDuplicateField,
			wantField: "/apiVersion",
		},
		{
			name:      "unknown",
			raw:       `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{},"notInSchema":1}`,
			wantCode:  apiproblem.CodeUnknownField,
			wantField: "/notInSchema",
		},
		{
			name:      "pointer-escape-tilde",
			raw:       `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{},"a~b":1}`,
			wantCode:  apiproblem.CodeUnknownField,
			wantField: "/a~0b",
		},
		{
			name:      "pointer-escape-slash",
			raw:       `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{},"a/b":1}`,
			wantCode:  apiproblem.CodeUnknownField,
			wantField: "/a~1b",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var dst decodeSample
			prob := DecodeJSON([]byte(tc.raw), testLimits, permissivePolicy, &dst)
			if prob == nil {
				t.Fatal("expected Problem, got nil")
			}
			if prob.Code != tc.wantCode {
				t.Fatalf("Code = %q, want %q", prob.Code, tc.wantCode)
			}
			if len(prob.Violations) != 1 {
				t.Fatalf("Violations = %#v", prob.Violations)
			}
			if prob.Violations[0].Field != tc.wantField {
				t.Fatalf("Field = %q, want %q", prob.Violations[0].Field, tc.wantField)
			}
			if !strings.HasPrefix(prob.Violations[0].Field, "/") {
				t.Fatalf("Field %q is not an RFC 6901 JSON Pointer", prob.Violations[0].Field)
			}
		})
	}
}

func TestDecodeJSONFieldPolicyRejectsStatus(t *testing.T) {
	const raw = `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{},"status":{"phase":"Ready"}}`
	pol := FieldPolicy{Mode: ModeCreateRequest, AllowSpecMutation: true}
	var dst decodeSample
	prob := DecodeJSON([]byte(raw), testLimits, pol, &dst)
	if prob == nil {
		t.Fatal("expected status rejection, got nil")
	}
	if prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeValidationFailed)
	}
	if prob.Violations[0].Field != "/status" {
		t.Fatalf("Field = %q, want /status", prob.Violations[0].Field)
	}
}

func TestDecodeJSONFieldPolicyRejectsSystemOwned(t *testing.T) {
	const raw = `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"spec":{}}`
	pol := FieldPolicy{Mode: ModeCreateRequest, AllowSpecMutation: true}
	var dst decodeSample
	prob := DecodeJSON([]byte(raw), testLimits, pol, &dst)
	if prob == nil {
		t.Fatal("expected system-owned rejection, got nil")
	}
	if prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeValidationFailed)
	}
	if prob.Violations[0].Field != "/metadata/uid" {
		t.Fatalf("Field = %q, want /metadata/uid", prob.Violations[0].Field)
	}
}

func TestDecodeJSONFieldPolicyAllowsStatusWhenPermitted(t *testing.T) {
	const raw = `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{},"status":{"phase":"Ready"}}`
	var dst decodeSample
	if prob := DecodeJSON([]byte(raw), testLimits, permissivePolicy, &dst); prob != nil {
		t.Fatalf("permissive policy must accept status: %#v", prob)
	}
	if dst.Status["phase"] != "Ready" {
		t.Fatalf("status not decoded: %#v", dst.Status)
	}
}

func TestDecodeJSONRejectsOversizedBody(t *testing.T) {
	lim := Limits{MaxObjectBytes: 8, MaxNestingDepth: 32}
	raw := []byte(`{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{}}`)
	var dst decodeSample
	prob := DecodeJSON(raw, lim, permissivePolicy, &dst)
	if prob == nil {
		t.Fatal("expected REQUEST_TOO_LARGE, got nil")
	}
	if prob.Code != apiproblem.CodeRequestTooLarge {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeRequestTooLarge)
	}
}

func TestDecodeJSONRejectsEmptyBody(t *testing.T) {
	var dst decodeSample
	prob := DecodeJSON(nil, testLimits, permissivePolicy, &dst)
	if prob == nil {
		t.Fatal("expected MALFORMED_REQUEST, got nil")
	}
	if prob.Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeMalformedRequest)
	}
}

func TestDecodeJSONRejectsNestingDepth(t *testing.T) {
	lim := Limits{MaxObjectBytes: 1 << 20, MaxNestingDepth: 2}
	// depth: root(1) -> a(2) -> b(3) exceeds 2
	const raw = `{"a":{"b":{"c":1}}}`
	var dst map[string]any
	prob := DecodeJSON([]byte(raw), lim, permissivePolicy, &dst)
	if prob == nil {
		t.Fatal("expected nesting rejection, got nil")
	}
	if prob.Code != apiproblem.CodeRequestTooLarge {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeRequestTooLarge)
	}
}

func TestEscapeJSONPointer(t *testing.T) {
	if got := escapeJSONPointer("a/b~c"); got != "a~1b~0c" {
		t.Fatalf("escapeJSONPointer = %q, want a~1b~0c", got)
	}
}
