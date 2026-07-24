package apivalid

import (
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// Task 8.3: JSON/YAML decode-only equivalence (D-03a).
// These tests exercise DecodeJSON and DecodeYAML only — no pipeline.

func TestDecodeJSONYAMLTypedValueEquivalence(t *testing.T) {
	cases := []struct {
		name    string
		jsonRaw string
		yamlRaw string
	}{
		{
			name:    "minimal-project",
			jsonRaw: `{"apiVersion":"platform.sovrunn.io/v1","kind":"Project","metadata":{"name":"demo"},"spec":{"displayName":"Demo"}}`,
			yamlRaw: `
apiVersion: platform.sovrunn.io/v1
kind: Project
metadata:
  name: demo
spec:
  displayName: Demo
`,
		},
		{
			name:    "nested-scalars-and-array",
			jsonRaw: `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{"enabled":true,"count":3,"tags":["a","b"],"note":null}}`,
			yamlRaw: `
apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  enabled: true
  count: 3
  tags:
    - a
    - b
  note: null
`,
		},
		{
			name:    "with-status-and-system-owned",
			jsonRaw: `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","generation":2,"resourceVersion":"7","createdAt":"2020-01-01T00:00:00Z","updatedAt":"2020-01-02T00:00:00Z"},"spec":{"displayName":"Demo"},"status":{"phase":"Ready"}}`,
			yamlRaw: `
apiVersion: v1
kind: Project
metadata:
  name: demo
  uid: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  generation: 2
  resourceVersion: "7"
  createdAt: "2020-01-01T00:00:00Z"
  updatedAt: "2020-01-02T00:00:00Z"
spec:
  displayName: Demo
status:
  phase: Ready
`,
		},
		{
			name:    "flow-style-identical-bytes",
			jsonRaw: `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{}}`,
			yamlRaw: `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{}}`,
		},
		{
			name:    "quoted-strings-and-bools",
			jsonRaw: `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{"label":"true","flag":false}}`,
			yamlRaw: `
apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  label: "true"
  flag: false
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var fromJSON, fromYAML decodeSample
			if prob := DecodeJSON([]byte(tc.jsonRaw), testLimits, permissivePolicy, &fromJSON); prob != nil {
				t.Fatalf("DecodeJSON: %#v", prob)
			}
			if prob := DecodeYAML([]byte(tc.yamlRaw), testLimits, permissivePolicy, &fromYAML); prob != nil {
				t.Fatalf("DecodeYAML: %#v", prob)
			}
			assertDecodeSamplesEqual(t, fromYAML, fromJSON)
			if !reflect.DeepEqual(fromYAML, fromJSON) {
				t.Fatalf("typed values diverge:\nYAML=%#v\nJSON=%#v", fromYAML, fromJSON)
			}
		})
	}
}

func TestDecodeJSONYAMLErrorCodeAndPointerEquivalence(t *testing.T) {
	createPol := PolicyFor(ModeCreateRequest)
	replacePol := PolicyFor(ModeReplaceRequest)

	cases := []struct {
		name      string
		jsonRaw   string
		yamlRaw   string
		policy    FieldPolicy
		wantCode  apiproblem.ErrorCode
		wantField string
		wantVCode apiproblem.ViolationCode
	}{
		{
			name:      "unknown-field",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{},"extraField":true}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: demo\nspec: {}\nextraField: true\n",
			policy:    permissivePolicy,
			wantCode:  apiproblem.CodeUnknownField,
			wantField: "/extraField",
			wantVCode: apiproblem.ViolationUnknownField,
		},
		{
			// encoding/json DisallowUnknownFields reports the leaf field name only.
			name:      "unknown-nested-field",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo","mystery":1},"spec":{}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: demo\n  mystery: 1\nspec: {}\n",
			policy:    permissivePolicy,
			wantCode:  apiproblem.CodeUnknownField,
			wantField: "/mystery",
			wantVCode: apiproblem.ViolationUnknownField,
		},
		{
			name:      "duplicate-field",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"a","name":"b"},"spec":{}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: a\n  name: b\nspec: {}\n",
			policy:    permissivePolicy,
			wantCode:  apiproblem.CodeDuplicateField,
			wantField: "/metadata/name",
			wantVCode: apiproblem.ViolationDuplicateField,
		},
		{
			name:      "duplicate-root-kind",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","kind":"Tenant","metadata":{"name":"x"},"spec":{}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nkind: Tenant\nmetadata:\n  name: x\nspec: {}\n",
			policy:    permissivePolicy,
			wantCode:  apiproblem.CodeDuplicateField,
			wantField: "/kind",
			wantVCode: apiproblem.ViolationDuplicateField,
		},
		{
			name:      "field-policy-status-create",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{},"status":{"phase":"Ready"}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: x\nspec: {}\nstatus:\n  phase: Ready\n",
			policy:    createPol,
			wantCode:  apiproblem.CodeValidationFailed,
			wantField: "/status",
			wantVCode: apiproblem.ViolationCode(apiproblem.CodeValidationFailed),
		},
		{
			name:      "field-policy-status-replace",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{},"status":{"phase":"Ready"}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: x\nspec: {}\nstatus:\n  phase: Ready\n",
			policy:    replacePol,
			wantCode:  apiproblem.CodeValidationFailed,
			wantField: "/status",
			wantVCode: apiproblem.ViolationCode(apiproblem.CodeValidationFailed),
		},
		{
			name:      "field-policy-system-owned-uid",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"spec":{}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: x\n  uid: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\nspec: {}\n",
			policy:    createPol,
			wantCode:  apiproblem.CodeValidationFailed,
			wantField: "/metadata/uid",
			wantVCode: apiproblem.ViolationCode(apiproblem.CodeValidationFailed),
		},
		{
			name:      "field-policy-system-owned-resourceVersion",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x","resourceVersion":"9"},"spec":{}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: x\n  resourceVersion: \"9\"\nspec: {}\n",
			policy:    createPol,
			wantCode:  apiproblem.CodeValidationFailed,
			wantField: "/metadata/resourceVersion",
			wantVCode: apiproblem.ViolationCode(apiproblem.CodeValidationFailed),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var jsonDst, yamlDst decodeSample
			jsonProb := DecodeJSON([]byte(tc.jsonRaw), testLimits, tc.policy, &jsonDst)
			yamlProb := DecodeYAML([]byte(tc.yamlRaw), testLimits, tc.policy, &yamlDst)
			assertEquivalentProblems(t, jsonProb, yamlProb)
			assertProblemCodeField(t, jsonProb, tc.wantCode, tc.wantField, tc.wantVCode)
			assertProblemCodeField(t, yamlProb, tc.wantCode, tc.wantField, tc.wantVCode)
		})
	}
}

func TestDecodeJSONYAMLFieldPolicyAcceptanceEquivalence(t *testing.T) {
	rawJSON := `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"spec":{"displayName":"X"},"status":{"phase":"Ready"}}`
	rawYAML := `
apiVersion: v1
kind: Project
metadata:
  name: x
  uid: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
spec:
  displayName: X
status:
  phase: Ready
`
	modes := []DecodeMode{
		ModeInternalObject,
		ModeReadRepresentation,
	}
	for _, mode := range modes {
		t.Run(mode.String(), func(t *testing.T) {
			pol := PolicyFor(mode)
			var fromJSON, fromYAML decodeSample
			if prob := DecodeJSON([]byte(rawJSON), testLimits, pol, &fromJSON); prob != nil {
				t.Fatalf("DecodeJSON under %s: %#v", mode, prob)
			}
			if prob := DecodeYAML([]byte(rawYAML), testLimits, pol, &fromYAML); prob != nil {
				t.Fatalf("DecodeYAML under %s: %#v", mode, prob)
			}
			if !reflect.DeepEqual(fromYAML, fromJSON) {
				t.Fatalf("typed values diverge under %s:\nYAML=%#v\nJSON=%#v", mode, fromYAML, fromJSON)
			}
		})
	}
}

func TestDecodeYAMLOnlyConstructsRejectedBeforeJSONNormalization(t *testing.T) {
	// Each case pairs a YAML-only construct with a field that would produce a
	// different Problem if DecodeJSON/FieldPolicy ran first. Rejection must
	// be MALFORMED_REQUEST from the YAML safety pass (D-03a step 2), not the
	// post-normalization DecodeJSON codes.
	createPol := PolicyFor(ModeCreateRequest)

	cases := []struct {
		name           string
		raw            string
		policy         FieldPolicy
		wantField      string
		detailContains string
		// If YAML normalization + DecodeJSON ran first, these would be expected
		// instead — proving they must NOT appear.
		forbiddenCode apiproblem.ErrorCode
	}{
		{
			name: "anchor-before-unknown-field",
			raw: `
apiVersion: v1
kind: Project
metadata:
  name: &anchor demo
spec: {}
extraField: true
`,
			policy:         permissivePolicy,
			wantField:      "/metadata/name",
			detailContains: "YAML anchors",
			forbiddenCode:  apiproblem.CodeUnknownField,
		},
		{
			// Real YAML aliases require an earlier anchor; the safety pass
			// rejects the anchor first — still before JSON normalization.
			name: "alias-document-anchor-before-unknown-field",
			raw: `
apiVersion: v1
kind: Project
metadata:
  name: &n demo
spec:
  displayName: *n
extraField: true
`,
			policy:         permissivePolicy,
			wantField:      "/metadata/name",
			detailContains: "YAML anchors",
			forbiddenCode:  apiproblem.CodeUnknownField,
		},
		{
			name: "merge-key-before-status-policy",
			raw: `
apiVersion: v1
kind: Project
metadata:
  name: x
spec: {}
status:
  <<:
    phase: Ready
`,
			policy:         createPol,
			wantField:      "/status/<<",
			detailContains: "YAML merge keys",
			forbiddenCode:  apiproblem.CodeValidationFailed,
		},
		{
			name: "nan-before-status-policy",
			raw: `
apiVersion: v1
kind: Project
metadata:
  name: x
spec: {}
status:
  phase: Ready
score: .nan
`,
			policy:         createPol,
			wantField:      "/score",
			detailContains: "non-finite YAML numbers",
			forbiddenCode:  apiproblem.CodeValidationFailed,
		},
		{
			name: "explicit-tag-before-unknown-field",
			raw: `
apiVersion: v1
kind: Project
metadata:
  name: !!str demo
spec: {}
extraField: true
`,
			policy:         permissivePolicy,
			wantField:      "/metadata/name",
			detailContains: "YAML explicit tags",
			forbiddenCode:  apiproblem.CodeUnknownField,
		},
		{
			name: "timestamp-before-unknown-field",
			raw: `
apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  created: 2020-01-01
extraField: true
`,
			policy:         permissivePolicy,
			wantField:      "/spec/created",
			detailContains: "YAML timestamp",
			forbiddenCode:  apiproblem.CodeUnknownField,
		},
		{
			name: "non-string-key-before-unknown-field",
			raw: `
apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  1: hello
extraField: true
`,
			policy:         permissivePolicy,
			wantField:      "/spec",
			detailContains: "YAML mapping keys must be strings",
			forbiddenCode:  apiproblem.CodeUnknownField,
		},
		{
			name: "multiple-documents-before-unknown-field",
			raw: `apiVersion: v1
kind: Project
metadata:
  name: demo
spec: {}
extraField: true
---
other: 1
`,
			policy:         permissivePolicy,
			wantField:      "/",
			detailContains: "multiple YAML documents",
			forbiddenCode:  apiproblem.CodeUnknownField,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var dst decodeSample
			prob := DecodeYAML([]byte(tc.raw), testLimits, tc.policy, &dst)
			if prob == nil {
				t.Fatal("expected YAML-only rejection before JSON normalization, got nil")
			}
			if prob.Code != apiproblem.CodeMalformedRequest {
				t.Fatalf("Code = %q, want %q (detail=%q); YAML-only constructs must fail before DecodeJSON",
					prob.Code, apiproblem.CodeMalformedRequest, prob.Detail)
			}
			if prob.Code == tc.forbiddenCode {
				t.Fatalf("got forbidden post-normalization code %q; YAML safety pass must run first", tc.forbiddenCode)
			}
			if len(prob.Violations) != 1 {
				t.Fatalf("Violations = %#v, want one", prob.Violations)
			}
			if prob.Violations[0].Field != tc.wantField {
				t.Fatalf("Field = %q, want %q", prob.Violations[0].Field, tc.wantField)
			}
			if !strings.Contains(prob.Detail, tc.detailContains) {
				t.Fatalf("Detail = %q, want substring %q", prob.Detail, tc.detailContains)
			}
		})
	}
}

func TestDecodeYAMLOnlyFeatureTableRejectedBeforeNormalization(t *testing.T) {
	// Representative YAML-only features from D-03a; each must fail at the
	// safety pass with MALFORMED_REQUEST before any JSON normalization.
	cases := []struct {
		name           string
		raw            string
		detailContains string
	}{
		// Alias documents always declare an anchor first; either message proves
		// the YAML safety pass ran before JSON normalization.
		{name: "alias-with-anchor", raw: "a: &x 1\nb: *x\n", detailContains: "YAML anchors"},
		{name: "anchor", raw: "a: &x 1\nb: 2\n", detailContains: "YAML anchors"},
		{name: "merge-key", raw: "obj:\n  <<:\n    k: 1\n", detailContains: "YAML merge keys"},
		{name: "custom-tag", raw: "x: !foo bar\n", detailContains: "YAML"},
		{name: "explicit-tag", raw: "x: !!str 123\n", detailContains: "YAML explicit tags"},
		{name: "nan", raw: "x: .nan\n", detailContains: "non-finite YAML numbers"},
		{name: "inf", raw: "x: .inf\n", detailContains: "non-finite YAML numbers"},
		{name: "binary", raw: "x: !!binary aGVsbG8=\n", detailContains: "YAML"},
		{name: "timestamp", raw: "x: 2020-01-01\n", detailContains: "YAML timestamp"},
		{name: "non-string-key", raw: "1: hello\n", detailContains: "YAML mapping keys must be strings"},
		{name: "multiple-docs", raw: "a: 1\n---\nb: 2\n", detailContains: "multiple YAML documents"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var dst map[string]any
			prob := DecodeYAML([]byte(tc.raw), testLimits, permissivePolicy, &dst)
			if prob == nil {
				t.Fatal("expected Problem, got nil")
			}
			if prob.Code != apiproblem.CodeMalformedRequest {
				t.Fatalf("Code = %q, want %q (detail=%q)", prob.Code, apiproblem.CodeMalformedRequest, prob.Detail)
			}
			if !strings.Contains(prob.Detail, tc.detailContains) {
				t.Fatalf("Detail = %q, want substring %q", prob.Detail, tc.detailContains)
			}
			// Post-normalization DecodeJSON codes must not appear for YAML-only constructs.
			switch prob.Code {
			case apiproblem.CodeUnknownField, apiproblem.CodeDuplicateField, apiproblem.CodeValidationFailed:
				t.Fatalf("YAML-only construct produced DecodeJSON-stage code %q", prob.Code)
			}
		})
	}
}

func assertEquivalentProblems(t *testing.T, jsonProb, yamlProb *apiproblem.Problem) {
	t.Helper()
	if jsonProb == nil || yamlProb == nil {
		t.Fatalf("expected both Problems non-nil; json=%#v yaml=%#v", jsonProb, yamlProb)
	}
	if jsonProb.Code != yamlProb.Code {
		t.Fatalf("Code diverge: json=%q yaml=%q", jsonProb.Code, yamlProb.Code)
	}
	if len(jsonProb.Violations) != len(yamlProb.Violations) {
		t.Fatalf("Violations len diverge: json=%d yaml=%d\njson=%#v\nyaml=%#v",
			len(jsonProb.Violations), len(yamlProb.Violations), jsonProb.Violations, yamlProb.Violations)
	}
	for i := range jsonProb.Violations {
		jv, yv := jsonProb.Violations[i], yamlProb.Violations[i]
		if jv.Field != yv.Field {
			t.Fatalf("Violations[%d].Field diverge: json=%q yaml=%q", i, jv.Field, yv.Field)
		}
		if jv.Code != yv.Code {
			t.Fatalf("Violations[%d].Code diverge: json=%q yaml=%q", i, jv.Code, yv.Code)
		}
	}
}

func assertProblemCodeField(t *testing.T, prob *apiproblem.Problem, code apiproblem.ErrorCode, field string, vCode apiproblem.ViolationCode) {
	t.Helper()
	if prob == nil {
		t.Fatal("expected Problem, got nil")
	}
	if prob.Code != code {
		t.Fatalf("Code = %q, want %q", prob.Code, code)
	}
	if len(prob.Violations) != 1 {
		t.Fatalf("Violations = %#v, want one", prob.Violations)
	}
	if prob.Violations[0].Field != field {
		t.Fatalf("Field = %q, want %q", prob.Violations[0].Field, field)
	}
	if vCode != "" && prob.Violations[0].Code != vCode {
		t.Fatalf("Violation.Code = %q, want %q", prob.Violations[0].Code, vCode)
	}
	if !strings.HasPrefix(prob.Violations[0].Field, "/") {
		t.Fatalf("Field %q is not an RFC 6901 JSON Pointer", prob.Violations[0].Field)
	}
}
