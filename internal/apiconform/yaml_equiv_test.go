package apiconform

import (
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// Task 8.4: full-pipeline JSON/YAML equivalence with test-local schemas (D-03a).
// DecodeJSON/DecodeYAML + StructuralValidator adapter (task 8.2). Not canonical
// fixtures from task 14.

const (
	yamlEquivSchemaID   = "api/schemas/test-local-project.json"
	yamlEquivCommonMeta = "api/schemas/_common/object-meta.json"
)

// yamlEquivTyped is a typed decode destination used for DisallowUnknownFields
// and FieldPolicy coverage (maps do not reject unknown keys at decode time).
type yamlEquivTyped struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   yamlEquivMeta    `json:"metadata"`
	Spec       yamlEquivSpec    `json:"spec"`
	Status     *yamlEquivStatus `json:"status,omitempty"`
}

type yamlEquivMeta struct {
	Name            string `json:"name"`
	UID             string `json:"uid,omitempty"`
	Generation      int64  `json:"generation,omitempty"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
}

type yamlEquivSpec struct {
	DisplayName string `json:"displayName"`
	Tier        string `json:"tier,omitempty"`
	Count       int64  `json:"count,omitempty"`
}

type yamlEquivStatus struct {
	Phase string `json:"phase,omitempty"`
}

func yamlEquivSchemas() map[string][]byte {
	return map[string][]byte{
		yamlEquivSchemaID: []byte(`{
			"type": "object",
			"properties": {
				"apiVersion": { "type": "string", "minLength": 1 },
				"kind": { "type": "string", "enum": ["Project"] },
				"metadata": { "$ref": "_common/object-meta.json" },
				"spec": {
					"type": "object",
					"properties": {
						"displayName": { "type": "string", "minLength": 1 },
						"tier": { "type": "string", "enum": ["small", "large"] },
						"count": { "type": "integer", "minimum": 0 }
					},
					"required": ["displayName"],
					"additionalProperties": false
				},
				"status": {
					"type": "object",
					"properties": {
						"phase": { "type": "string", "enum": ["Pending", "Ready"] }
					},
					"additionalProperties": false
				}
			},
			"required": ["apiVersion", "kind", "metadata", "spec"],
			"additionalProperties": false
		}`),
		yamlEquivCommonMeta: []byte(`{
			"type": "object",
			"properties": {
				"name": { "type": "string", "minLength": 1 },
				"uid": { "type": "string", "minLength": 1 },
				"generation": { "type": "integer", "minimum": 0 },
				"resourceVersion": { "type": "string", "minLength": 1 }
			},
			"required": ["name"],
			"additionalProperties": false
		}`),
	}
}

func yamlEquivValidator(t *testing.T) *StructuralValidator {
	t.Helper()
	return testStructuralValidator(t, yamlEquivSchemas())
}

func yamlEquivLimits() apivalid.Limits {
	return apivalid.DefaultLimits()
}

func yamlEquivPermissivePolicy() apivalid.FieldPolicy {
	return apivalid.PolicyFor(apivalid.ModeReadRepresentation)
}

// decodeAndStructurallyValidate runs DecodeJSON or DecodeYAML, then the
// StructuralValidator adapter when decode succeeds (task 8.4 full path).
func decodeAndStructurallyValidate(
	raw []byte,
	asYAML bool,
	lim apivalid.Limits,
	pol apivalid.FieldPolicy,
	dst any,
	v *StructuralValidator,
	schemaID string,
) (decodeProb *apiproblem.Problem, violations []apiproblem.Violation, validateErr error) {
	if asYAML {
		decodeProb = apivalid.DecodeYAML(raw, lim, pol, dst)
	} else {
		decodeProb = apivalid.DecodeJSON(raw, lim, pol, dst)
	}
	if decodeProb != nil {
		return decodeProb, nil, nil
	}
	violations, validateErr = v.Validate(dst, schemaID)
	return nil, violations, validateErr
}

func TestYAMLEquivTypedValuesDecodeAndValidate(t *testing.T) {
	t.Parallel()

	v := yamlEquivValidator(t)
	lim := yamlEquivLimits()
	pol := yamlEquivPermissivePolicy()

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
			name:    "nested-scalars-enum-and-integer",
			jsonRaw: `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{"displayName":"Demo","tier":"small","count":3}}`,
			yamlRaw: `
apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  displayName: Demo
  tier: small
  count: 3
`,
		},
		{
			name:    "with-status-and-system-owned",
			jsonRaw: `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","generation":2,"resourceVersion":"7"},"spec":{"displayName":"Demo"},"status":{"phase":"Ready"}}`,
			yamlRaw: `
apiVersion: v1
kind: Project
metadata:
  name: demo
  uid: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  generation: 2
  resourceVersion: "7"
spec:
  displayName: Demo
status:
  phase: Ready
`,
		},
		{
			name:    "flow-style-identical-bytes",
			jsonRaw: `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{"displayName":"Demo"}}`,
			yamlRaw: `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{"displayName":"Demo"}}`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var fromJSON, fromYAML yamlEquivTyped
			jsonProb, jsonViolations, jsonErr := decodeAndStructurallyValidate(
				[]byte(tc.jsonRaw), false, lim, pol, &fromJSON, v, yamlEquivSchemaID)
			yamlProb, yamlViolations, yamlErr := decodeAndStructurallyValidate(
				[]byte(tc.yamlRaw), true, lim, pol, &fromYAML, v, yamlEquivSchemaID)

			if jsonProb != nil || yamlProb != nil {
				t.Fatalf("decode failed: json=%#v yaml=%#v", jsonProb, yamlProb)
			}
			if jsonErr != nil || yamlErr != nil {
				t.Fatalf("structural validate error: json=%v yaml=%v", jsonErr, yamlErr)
			}
			if len(jsonViolations) != 0 || len(yamlViolations) != 0 {
				t.Fatalf("want no violations; json=%#v yaml=%#v", jsonViolations, yamlViolations)
			}
			if !reflect.DeepEqual(fromYAML, fromJSON) {
				t.Fatalf("typed values diverge:\nYAML=%#v\nJSON=%#v", fromYAML, fromJSON)
			}
		})
	}
}

func TestYAMLEquivMapValuesDecodeAndValidate(t *testing.T) {
	t.Parallel()

	v := yamlEquivValidator(t)
	lim := yamlEquivLimits()
	pol := yamlEquivPermissivePolicy()

	jsonRaw := `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{"displayName":"Demo","tier":"large","count":1}}`
	yamlRaw := `
apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  displayName: Demo
  tier: large
  count: 1
`

	var fromJSON, fromYAML map[string]any
	jsonProb, jsonViolations, jsonErr := decodeAndStructurallyValidate(
		[]byte(jsonRaw), false, lim, pol, &fromJSON, v, yamlEquivSchemaID)
	yamlProb, yamlViolations, yamlErr := decodeAndStructurallyValidate(
		[]byte(yamlRaw), true, lim, pol, &fromYAML, v, yamlEquivSchemaID)

	if jsonProb != nil || yamlProb != nil {
		t.Fatalf("decode failed: json=%#v yaml=%#v", jsonProb, yamlProb)
	}
	if jsonErr != nil || yamlErr != nil {
		t.Fatalf("structural validate error: json=%v yaml=%v", jsonErr, yamlErr)
	}
	if len(jsonViolations) != 0 || len(yamlViolations) != 0 {
		t.Fatalf("want no violations; json=%#v yaml=%#v", jsonViolations, yamlViolations)
	}
	if !reflect.DeepEqual(fromYAML, fromJSON) {
		t.Fatalf("map values diverge:\nYAML=%#v\nJSON=%#v", fromYAML, fromJSON)
	}
}

func TestYAMLEquivDecodeErrorCodeAndPointerEquivalence(t *testing.T) {
	t.Parallel()

	v := yamlEquivValidator(t)
	lim := yamlEquivLimits()
	createPol := apivalid.PolicyFor(apivalid.ModeCreateRequest)
	replacePol := apivalid.PolicyFor(apivalid.ModeReplaceRequest)
	permissive := yamlEquivPermissivePolicy()

	cases := []struct {
		name      string
		jsonRaw   string
		yamlRaw   string
		policy    apivalid.FieldPolicy
		wantCode  apiproblem.ErrorCode
		wantField string
		wantVCode apiproblem.ViolationCode
	}{
		{
			name:      "unknown-field",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{"displayName":"Demo"},"mystery":true}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: demo\nspec:\n  displayName: Demo\nmystery: true\n",
			policy:    permissive,
			wantCode:  apiproblem.CodeUnknownField,
			wantField: "/mystery",
			wantVCode: apiproblem.ViolationUnknownField,
		},
		{
			name:      "unknown-nested-field",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo","mystery":1},"spec":{"displayName":"Demo"}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: demo\n  mystery: 1\nspec:\n  displayName: Demo\n",
			policy:    permissive,
			wantCode:  apiproblem.CodeUnknownField,
			wantField: "/mystery",
			wantVCode: apiproblem.ViolationUnknownField,
		},
		{
			name:      "duplicate-field",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"a","name":"b"},"spec":{"displayName":"Demo"}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: a\n  name: b\nspec:\n  displayName: Demo\n",
			policy:    permissive,
			wantCode:  apiproblem.CodeDuplicateField,
			wantField: "/metadata/name",
			wantVCode: apiproblem.ViolationDuplicateField,
		},
		{
			name:      "field-policy-status-create",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{"displayName":"X"},"status":{"phase":"Ready"}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: x\nspec:\n  displayName: X\nstatus:\n  phase: Ready\n",
			policy:    createPol,
			wantCode:  apiproblem.CodeValidationFailed,
			wantField: "/status",
			wantVCode: apiproblem.ViolationCode(apiproblem.CodeValidationFailed),
		},
		{
			name:      "field-policy-status-replace",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{"displayName":"X"},"status":{"phase":"Ready"}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: x\nspec:\n  displayName: X\nstatus:\n  phase: Ready\n",
			policy:    replacePol,
			wantCode:  apiproblem.CodeValidationFailed,
			wantField: "/status",
			wantVCode: apiproblem.ViolationCode(apiproblem.CodeValidationFailed),
		},
		{
			name:      "field-policy-system-owned-uid",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"spec":{"displayName":"X"}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: x\n  uid: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\nspec:\n  displayName: X\n",
			policy:    createPol,
			wantCode:  apiproblem.CodeValidationFailed,
			wantField: "/metadata/uid",
			wantVCode: apiproblem.ViolationCode(apiproblem.CodeValidationFailed),
		},
		{
			name:      "field-policy-system-owned-resourceVersion",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x","resourceVersion":"9"},"spec":{"displayName":"X"}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: x\n  resourceVersion: \"9\"\nspec:\n  displayName: X\n",
			policy:    createPol,
			wantCode:  apiproblem.CodeValidationFailed,
			wantField: "/metadata/resourceVersion",
			wantVCode: apiproblem.ViolationCode(apiproblem.CodeValidationFailed),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var jsonDst, yamlDst yamlEquivTyped
			jsonProb, jsonViolations, jsonErr := decodeAndStructurallyValidate(
				[]byte(tc.jsonRaw), false, lim, tc.policy, &jsonDst, v, yamlEquivSchemaID)
			yamlProb, yamlViolations, yamlErr := decodeAndStructurallyValidate(
				[]byte(tc.yamlRaw), true, lim, tc.policy, &yamlDst, v, yamlEquivSchemaID)

			assertYAMLEquivProblems(t, jsonProb, yamlProb)
			assertYAMLEquivProblemCodeField(t, jsonProb, tc.wantCode, tc.wantField, tc.wantVCode)
			assertYAMLEquivProblemCodeField(t, yamlProb, tc.wantCode, tc.wantField, tc.wantVCode)

			// Decode failures must stop before structural validation.
			if jsonErr != nil || yamlErr != nil {
				t.Fatalf("structural validate must not run after decode failure: json=%v yaml=%v", jsonErr, yamlErr)
			}
			if len(jsonViolations) != 0 || len(yamlViolations) != 0 {
				t.Fatalf("structural violations must be empty after decode failure: json=%#v yaml=%#v",
					jsonViolations, yamlViolations)
			}
		})
	}
}

func TestYAMLEquivFieldPolicyAcceptanceThenValidate(t *testing.T) {
	t.Parallel()

	v := yamlEquivValidator(t)
	lim := yamlEquivLimits()
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
	modes := []apivalid.DecodeMode{
		apivalid.ModeInternalObject,
		apivalid.ModeReadRepresentation,
	}
	for _, mode := range modes {
		mode := mode
		t.Run(mode.String(), func(t *testing.T) {
			t.Parallel()

			pol := apivalid.PolicyFor(mode)
			var fromJSON, fromYAML yamlEquivTyped
			jsonProb, jsonViolations, jsonErr := decodeAndStructurallyValidate(
				[]byte(rawJSON), false, lim, pol, &fromJSON, v, yamlEquivSchemaID)
			yamlProb, yamlViolations, yamlErr := decodeAndStructurallyValidate(
				[]byte(rawYAML), true, lim, pol, &fromYAML, v, yamlEquivSchemaID)

			if jsonProb != nil || yamlProb != nil {
				t.Fatalf("decode under %s: json=%#v yaml=%#v", mode, jsonProb, yamlProb)
			}
			if jsonErr != nil || yamlErr != nil {
				t.Fatalf("validate under %s: json=%v yaml=%v", mode, jsonErr, yamlErr)
			}
			if len(jsonViolations) != 0 || len(yamlViolations) != 0 {
				t.Fatalf("want no violations under %s; json=%#v yaml=%#v", mode, jsonViolations, yamlViolations)
			}
			if !reflect.DeepEqual(fromYAML, fromJSON) {
				t.Fatalf("typed values diverge under %s:\nYAML=%#v\nJSON=%#v", mode, fromYAML, fromJSON)
			}
		})
	}
}

func TestYAMLEquivStructuralViolationEquivalence(t *testing.T) {
	t.Parallel()

	v := yamlEquivValidator(t)
	lim := yamlEquivLimits()
	pol := yamlEquivPermissivePolicy()

	cases := []struct {
		name      string
		jsonRaw   string
		yamlRaw   string
		wantCode  string
		wantField string
	}{
		{
			name:      "missing-required-displayName",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: demo\nspec: {}\n",
			wantCode:  apischema.CodeRequiredField,
			wantField: "/spec/displayName",
		},
		{
			name:      "type-mismatch-count",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{"displayName":"Demo","count":"nope"}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: demo\nspec:\n  displayName: Demo\n  count: \"nope\"\n",
			wantCode:  apischema.CodeTypeMismatch,
			wantField: "/spec/count",
		},
		{
			name:      "enum-mismatch-tier",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{"displayName":"Demo","tier":"xlarge"}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: demo\nspec:\n  displayName: Demo\n  tier: xlarge\n",
			wantCode:  apischema.CodeEnumMismatch,
			wantField: "/spec/tier",
		},
		{
			name:      "schema-unknown-field-additionalProperties",
			jsonRaw:   `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{"displayName":"Demo","unexpected":true}}`,
			yamlRaw:   "apiVersion: v1\nkind: Project\nmetadata:\n  name: demo\nspec:\n  displayName: Demo\n  unexpected: true\n",
			wantCode:  apischema.CodeUnknownField,
			wantField: "/spec/unexpected",
		},
		{
			name:      "enum-mismatch-kind",
			jsonRaw:   `{"apiVersion":"v1","kind":"Tenant","metadata":{"name":"demo"},"spec":{"displayName":"Demo"}}`,
			yamlRaw:   "apiVersion: v1\nkind: Tenant\nmetadata:\n  name: demo\nspec:\n  displayName: Demo\n",
			wantCode:  apischema.CodeEnumMismatch,
			wantField: "/kind",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// map destination so decode accepts keys that only the schema rejects.
			var fromJSON, fromYAML map[string]any
			jsonProb, jsonViolations, jsonErr := decodeAndStructurallyValidate(
				[]byte(tc.jsonRaw), false, lim, pol, &fromJSON, v, yamlEquivSchemaID)
			yamlProb, yamlViolations, yamlErr := decodeAndStructurallyValidate(
				[]byte(tc.yamlRaw), true, lim, pol, &fromYAML, v, yamlEquivSchemaID)

			if jsonProb != nil || yamlProb != nil {
				t.Fatalf("decode must succeed before structural findings: json=%#v yaml=%#v", jsonProb, yamlProb)
			}
			if jsonErr != nil || yamlErr != nil {
				t.Fatalf("structural validate error: json=%v yaml=%v", jsonErr, yamlErr)
			}
			assertYAMLEquivViolations(t, jsonViolations, yamlViolations)
			if !hasViolation(jsonViolations, tc.wantCode, tc.wantField) {
				t.Fatalf("json want %s at %s, got %#v", tc.wantCode, tc.wantField, jsonViolations)
			}
			if !hasViolation(yamlViolations, tc.wantCode, tc.wantField) {
				t.Fatalf("yaml want %s at %s, got %#v", tc.wantCode, tc.wantField, yamlViolations)
			}
			if !reflect.DeepEqual(fromYAML, fromJSON) {
				t.Fatalf("decoded maps diverge:\nYAML=%#v\nJSON=%#v", fromYAML, fromJSON)
			}
		})
	}
}

func TestYAMLEquivCreateModeRejectsStatusBeforeStructural(t *testing.T) {
	t.Parallel()

	v := yamlEquivValidator(t)
	lim := yamlEquivLimits()
	pol := apivalid.PolicyFor(apivalid.ModeCreateRequest)

	jsonRaw := `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{"displayName":"X"},"status":{"phase":"Ready"}}`
	yamlRaw := `
apiVersion: v1
kind: Project
metadata:
  name: x
spec:
  displayName: X
status:
  phase: Ready
`

	var fromJSON, fromYAML map[string]any
	jsonProb, jsonViolations, jsonErr := decodeAndStructurallyValidate(
		[]byte(jsonRaw), false, lim, pol, &fromJSON, v, yamlEquivSchemaID)
	yamlProb, yamlViolations, yamlErr := decodeAndStructurallyValidate(
		[]byte(yamlRaw), true, lim, pol, &fromYAML, v, yamlEquivSchemaID)

	assertYAMLEquivProblems(t, jsonProb, yamlProb)
	assertYAMLEquivProblemCodeField(t, jsonProb, apiproblem.CodeValidationFailed, "/status",
		apiproblem.ViolationCode(apiproblem.CodeValidationFailed))
	assertYAMLEquivProblemCodeField(t, yamlProb, apiproblem.CodeValidationFailed, "/status",
		apiproblem.ViolationCode(apiproblem.CodeValidationFailed))
	if jsonErr != nil || yamlErr != nil || len(jsonViolations) != 0 || len(yamlViolations) != 0 {
		t.Fatalf("FieldPolicy must stop before structural validation")
	}
}

func assertYAMLEquivProblems(t *testing.T, jsonProb, yamlProb *apiproblem.Problem) {
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

func assertYAMLEquivProblemCodeField(t *testing.T, prob *apiproblem.Problem, code apiproblem.ErrorCode, field string, vCode apiproblem.ViolationCode) {
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

func assertYAMLEquivViolations(t *testing.T, jsonVs, yamlVs []apiproblem.Violation) {
	t.Helper()
	if len(jsonVs) != len(yamlVs) {
		t.Fatalf("violation len diverge: json=%d yaml=%d\njson=%#v\nyaml=%#v",
			len(jsonVs), len(yamlVs), jsonVs, yamlVs)
	}
	for i := range jsonVs {
		if jsonVs[i].Field != yamlVs[i].Field {
			t.Fatalf("Violations[%d].Field diverge: json=%q yaml=%q", i, jsonVs[i].Field, yamlVs[i].Field)
		}
		if jsonVs[i].Code != yamlVs[i].Code {
			t.Fatalf("Violations[%d].Code diverge: json=%q yaml=%q", i, jsonVs[i].Code, yamlVs[i].Code)
		}
		if !strings.HasPrefix(jsonVs[i].Field, "/") {
			t.Fatalf("Field %q is not an RFC 6901 JSON Pointer", jsonVs[i].Field)
		}
	}
}
