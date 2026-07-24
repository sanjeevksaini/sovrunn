package apivalid

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"gopkg.in/yaml.v3"
)

// permissivePolicy accepts status, system-owned metadata, and spec so
// duplicate/unknown-field tests are not confounded by FieldPolicy.
var permissivePolicy = PolicyFor(ModeReadRepresentation)

var testLimits = DefaultLimits()

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

func TestDecodeYAMLHappyPathEquivalentToJSON(t *testing.T) {
	const jsonRaw = `{"apiVersion":"platform.sovrunn.io/v1","kind":"Project","metadata":{"name":"demo"},"spec":{"displayName":"Demo"}}`
	const yamlRaw = `
apiVersion: platform.sovrunn.io/v1
kind: Project
metadata:
  name: demo
spec:
  displayName: Demo
`
	var fromJSON, fromYAML decodeSample
	if prob := DecodeJSON([]byte(jsonRaw), testLimits, permissivePolicy, &fromJSON); prob != nil {
		t.Fatalf("DecodeJSON: %#v", prob)
	}
	if prob := DecodeYAML([]byte(yamlRaw), testLimits, permissivePolicy, &fromYAML); prob != nil {
		t.Fatalf("DecodeYAML: %#v", prob)
	}
	assertDecodeSamplesEqual(t, fromYAML, fromJSON)
}

func TestDecodeYAMLFlowStyleEquivalentToJSON(t *testing.T) {
	const raw = `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{}}`
	var fromJSON, fromYAML decodeSample
	if prob := DecodeJSON([]byte(raw), testLimits, permissivePolicy, &fromJSON); prob != nil {
		t.Fatalf("DecodeJSON: %#v", prob)
	}
	if prob := DecodeYAML([]byte(raw), testLimits, permissivePolicy, &fromYAML); prob != nil {
		t.Fatalf("DecodeYAML flow: %#v", prob)
	}
	assertDecodeSamplesEqual(t, fromYAML, fromJSON)
}

func assertDecodeSamplesEqual(t *testing.T, got, want decodeSample) {
	t.Helper()
	if got.APIVersion != want.APIVersion || got.Kind != want.Kind || got.Metadata != want.Metadata {
		t.Fatalf("identity mismatch:\ngot=%#v\nwant=%#v", got, want)
	}
	if len(got.Spec) != len(want.Spec) {
		t.Fatalf("spec len got=%d want=%d\ngot=%#v\nwant=%#v", len(got.Spec), len(want.Spec), got.Spec, want.Spec)
	}
	for k, wv := range want.Spec {
		gv, ok := got.Spec[k]
		if !ok || fmt.Sprint(gv) != fmt.Sprint(wv) {
			t.Fatalf("spec[%q] got=%#v want=%#v", k, gv, wv)
		}
	}
	if len(got.Status) != len(want.Status) {
		t.Fatalf("status len got=%d want=%d", len(got.Status), len(want.Status))
	}
	for k, wv := range want.Status {
		gv, ok := got.Status[k]
		if !ok || fmt.Sprint(gv) != fmt.Sprint(wv) {
			t.Fatalf("status[%q] got=%#v want=%#v", k, gv, wv)
		}
	}
}

func TestDecodeYAMLRejectsEachYAMLOnlyFeature(t *testing.T) {
	cases := []struct {
		name     string
		raw      string
		wantCode apiproblem.ErrorCode
		wantPath string // empty means any root-ish pointer is acceptable
	}{
		{
			name:     "alias",
			raw:      "a: &x 1\nb: *x\n",
			wantCode: apiproblem.CodeMalformedRequest,
		},
		{
			name:     "anchor",
			raw:      "a: &x 1\nb: 2\n",
			wantCode: apiproblem.CodeMalformedRequest,
		},
		{
			name:     "merge-key",
			raw:      "obj:\n  <<:\n    k: 1\n  m: 2\n",
			wantCode: apiproblem.CodeMalformedRequest,
			wantPath: "/obj/<<",
		},
		{
			name:     "explicit-tag",
			raw:      "x: !!str 123\n",
			wantCode: apiproblem.CodeMalformedRequest,
			wantPath: "/x",
		},
		{
			name:     "custom-tag",
			raw:      "x: !foo bar\n",
			wantCode: apiproblem.CodeMalformedRequest,
			wantPath: "/x",
		},
		{
			name:     "nan",
			raw:      "x: .nan\n",
			wantCode: apiproblem.CodeMalformedRequest,
			wantPath: "/x",
		},
		{
			name:     "inf",
			raw:      "x: .inf\n",
			wantCode: apiproblem.CodeMalformedRequest,
			wantPath: "/x",
		},
		{
			name:     "neg-inf",
			raw:      "x: -.inf\n",
			wantCode: apiproblem.CodeMalformedRequest,
			wantPath: "/x",
		},
		{
			name:     "timestamp",
			raw:      "x: 2020-01-01\n",
			wantCode: apiproblem.CodeMalformedRequest,
			wantPath: "/x",
		},
		{
			name:     "binary",
			raw:      "x: !!binary aGVsbG8=\n",
			wantCode: apiproblem.CodeMalformedRequest,
			wantPath: "/x",
		},
		{
			name:     "multiple-documents",
			raw:      "a: 1\n---\nb: 2\n",
			wantCode: apiproblem.CodeMalformedRequest,
			wantPath: "/",
		},
		{
			name:     "non-string-int-key",
			raw:      "1: hello\n",
			wantCode: apiproblem.CodeMalformedRequest,
			wantPath: "/",
		},
		{
			name:     "non-string-bool-key",
			raw:      "true: hello\n",
			wantCode: apiproblem.CodeMalformedRequest,
			wantPath: "/",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var dst map[string]any
			prob := DecodeYAML([]byte(tc.raw), testLimits, permissivePolicy, &dst)
			if prob == nil {
				t.Fatal("expected Problem, got nil")
			}
			if prob.Code != tc.wantCode {
				t.Fatalf("Code = %q, want %q (detail=%q)", prob.Code, tc.wantCode, prob.Detail)
			}
			if len(prob.Violations) != 1 {
				t.Fatalf("Violations = %#v", prob.Violations)
			}
			if tc.wantPath != "" && prob.Violations[0].Field != tc.wantPath {
				t.Fatalf("Field = %q, want %q", prob.Violations[0].Field, tc.wantPath)
			}
			if !strings.HasPrefix(prob.Violations[0].Field, "/") {
				t.Fatalf("Field %q is not an RFC 6901 JSON Pointer", prob.Violations[0].Field)
			}
		})
	}
}

func TestDecodeYAMLRejectsDuplicateKey(t *testing.T) {
	const raw = `
apiVersion: v1
kind: Project
metadata:
  name: a
  name: b
spec: {}
`
	var dst decodeSample
	prob := DecodeYAML([]byte(raw), testLimits, permissivePolicy, &dst)
	if prob == nil {
		t.Fatal("expected duplicate-key Problem, got nil")
	}
	if prob.Code != apiproblem.CodeDuplicateField {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeDuplicateField)
	}
	if prob.Violations[0].Field != "/metadata/name" {
		t.Fatalf("Field = %q, want /metadata/name", prob.Violations[0].Field)
	}
	if prob.Violations[0].Code != apiproblem.ViolationDuplicateField {
		t.Fatalf("Violation.Code = %q, want %q", prob.Violations[0].Code, apiproblem.ViolationDuplicateField)
	}
}

func TestDecodeYAMLRejectsUnknownFieldViaDecodeJSON(t *testing.T) {
	const raw = `
apiVersion: v1
kind: Project
metadata:
  name: demo
spec: {}
extraField: true
`
	var dst decodeSample
	prob := DecodeYAML([]byte(raw), testLimits, permissivePolicy, &dst)
	if prob == nil {
		t.Fatal("expected unknown-field Problem, got nil")
	}
	if prob.Code != apiproblem.CodeUnknownField {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeUnknownField)
	}
	if prob.Violations[0].Field != "/extraField" {
		t.Fatalf("Field = %q, want /extraField", prob.Violations[0].Field)
	}
}

func TestDecodeYAMLFieldPolicyRejectsStatus(t *testing.T) {
	const raw = `
apiVersion: v1
kind: Project
metadata:
  name: x
spec: {}
status:
  phase: Ready
`
	pol := FieldPolicy{Mode: ModeCreateRequest, AllowSpecMutation: true}
	var dst decodeSample
	prob := DecodeYAML([]byte(raw), testLimits, pol, &dst)
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

func TestDecodeYAMLQuotedNonFiniteIsString(t *testing.T) {
	// Quoted ".nan" is a JSON string, not a YAML float — must be accepted
	// by the YAML safety pass and then fail only if the destination rejects it.
	const raw = "x: \".nan\"\n"
	var dst map[string]any
	if prob := DecodeYAML([]byte(raw), testLimits, permissivePolicy, &dst); prob != nil {
		t.Fatalf("quoted .nan must decode as string: %#v", prob)
	}
	if dst["x"] != ".nan" {
		t.Fatalf("x = %#v, want \".nan\"", dst["x"])
	}
}

func TestDecodeYAMLRejectsAliasNodeDirectly(t *testing.T) {
	n := &yaml.Node{
		Kind:  yaml.AliasNode,
		Value: "x",
		Alias: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: "1"},
	}
	prob := rejectYAMLOnlyConstructs(n, "/b")
	if prob == nil {
		t.Fatal("expected alias rejection, got nil")
	}
	if prob.Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeMalformedRequest)
	}
	if prob.Violations[0].Field != "/b" {
		t.Fatalf("Field = %q, want /b", prob.Violations[0].Field)
	}
}

func TestDecodeYAMLRejectsEmptyBody(t *testing.T) {
	var dst decodeSample
	prob := DecodeYAML(nil, testLimits, permissivePolicy, &dst)
	if prob == nil {
		t.Fatal("expected MALFORMED_REQUEST, got nil")
	}
	if prob.Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeMalformedRequest)
	}
}

func TestDecodeYAMLRejectsOversizedBody(t *testing.T) {
	lim := Limits{MaxObjectBytes: 8, MaxNestingDepth: 32}
	raw := []byte("apiVersion: v1\nkind: Project\n")
	var dst decodeSample
	prob := DecodeYAML(raw, lim, permissivePolicy, &dst)
	if prob == nil {
		t.Fatal("expected REQUEST_TOO_LARGE, got nil")
	}
	if prob.Code != apiproblem.CodeRequestTooLarge {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeRequestTooLarge)
	}
}
