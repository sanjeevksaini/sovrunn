package apischema

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func validAnnotatedSchema() map[string]any {
	return map[string]any{
		"$schema":                  "https://json-schema.org/draft/2020-12/schema",
		"$id":                      "https://sovrunn.example/schemas/demo.json",
		"title":                    "Demo",
		"type":                     "object",
		"x-sovrunn-profile":        "ManagedResource",
		"x-sovrunn-boundary":       "customer-facing",
		"x-sovrunn-allowed-scopes": []any{"Tenant", "Project"},
		"x-sovrunn-stability":      "stable",
		"properties": map[string]any{
			"name": map[string]any{
				"type": "string",
				"x-sovrunn-field-policy": map[string]any{
					"classification":    "Public",
					"authorizedWriter":  "creator",
					"authorizedReaders": []any{"customer"},
					"mutability":        "immutable",
					"retention":         "standard",
					"redaction":         "none",
					"residency":         "any",
					"auditRequired":     true,
				},
			},
			"oneOf": map[string]any{
				"type": "string",
			},
		},
		"additionalProperties": false,
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestReadAnnotationsValid(t *testing.T) {
	t.Parallel()

	meta, issues := ReadAnnotations(mustJSON(t, validAnnotatedSchema()))
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
	if meta.Profile != apimeta.ProfileManagedResource {
		t.Fatalf("profile=%q", meta.Profile)
	}
	if meta.Boundary != apimeta.BoundaryCustomerFacing {
		t.Fatalf("boundary=%q", meta.Boundary)
	}
	if meta.Stability != apimeta.StabilityStable {
		t.Fatalf("stability=%q", meta.Stability)
	}
	if len(meta.AllowedScopes) != 2 ||
		meta.AllowedScopes[0] != apimeta.ScopeTenant ||
		meta.AllowedScopes[1] != apimeta.ScopeProject {
		t.Fatalf("allowed scopes=%#v", meta.AllowedScopes)
	}
	fp, ok := meta.FieldPolicies["/properties/name"]
	if !ok {
		t.Fatalf("missing field policy for /properties/name; got keys %#v", fieldPolicyKeys(meta.FieldPolicies))
	}
	if fp.Classification != apimeta.ClassPublic ||
		fp.AuthorizedWriter != WriterCreator ||
		len(fp.AuthorizedReaders) != 1 || fp.AuthorizedReaders[0] != ReaderCustomer ||
		fp.Mutability != MutabilityImmutable ||
		fp.Retention != RetentionStandard ||
		fp.Redaction != RedactionNone ||
		fp.Residency != ResidencyAny ||
		!fp.AuditRequired {
		t.Fatalf("unexpected field policy: %#v", fp)
	}
}

func TestReadAnnotationsMissingRequiredAnnotation(t *testing.T) {
	t.Parallel()

	for _, missing := range []string{
		ExtProfile,
		ExtBoundary,
		ExtAllowedScopes,
		ExtStability,
	} {
		missing := missing
		t.Run(missing, func(t *testing.T) {
			t.Parallel()
			schema := validAnnotatedSchema()
			delete(schema, missing)
			_, issues := ReadAnnotations(mustJSON(t, schema))
			if !hasAnnotationIssue(issues, CodeAnnotationMissing, "/"+missing) {
				t.Fatalf("expected missing %s issue, got %#v", missing, issues)
			}
		})
	}
}

func TestReadAnnotationsInvalidVocabulary(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		mut  func(map[string]any)
		code string
		path string
	}{
		{
			name: "invalid-profile",
			mut:  func(m map[string]any) { m[ExtProfile] = "NotAProfile" },
			code: CodeAnnotationInvalid,
			path: "/" + ExtProfile,
		},
		{
			name: "invalid-boundary",
			mut:  func(m map[string]any) { m[ExtBoundary] = "public-api" },
			code: CodeAnnotationInvalid,
			path: "/" + ExtBoundary,
		},
		{
			name: "invalid-stability",
			mut:  func(m map[string]any) { m[ExtStability] = "ga" },
			code: CodeAnnotationInvalid,
			path: "/" + ExtStability,
		},
		{
			name: "invalid-scope-kind",
			mut:  func(m map[string]any) { m[ExtAllowedScopes] = []any{"Cluster"} },
			code: CodeAnnotationInvalid,
			path: "/" + ExtAllowedScopes + "/0",
		},
		{
			name: "invalid-classification",
			mut: func(m map[string]any) {
				props := m["properties"].(map[string]any)
				name := props["name"].(map[string]any)
				fp := name[ExtFieldPolicy].(map[string]any)
				fp["classification"] = "TopSecret"
			},
			code: CodeFieldPolicyInvalid,
			path: "/properties/name/" + ExtFieldPolicy + "/classification",
		},
		{
			name: "invalid-writer",
			mut: func(m map[string]any) {
				props := m["properties"].(map[string]any)
				name := props["name"].(map[string]any)
				fp := name[ExtFieldPolicy].(map[string]any)
				fp["authorizedWriter"] = "anyone"
			},
			code: CodeFieldPolicyInvalid,
			path: "/properties/name/" + ExtFieldPolicy + "/authorizedWriter",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			schema := validAnnotatedSchema()
			tc.mut(schema)
			_, issues := ReadAnnotations(mustJSON(t, schema))
			if !hasAnnotationIssue(issues, tc.code, tc.path) {
				t.Fatalf("expected code=%s path=%s, got %#v", tc.code, tc.path, issues)
			}
		})
	}
}

func TestReadAnnotationsUnknownExtensionFailsClosed(t *testing.T) {
	t.Parallel()

	schema := validAnnotatedSchema()
	schema["x-sovrunn-foo"] = "bar"
	_, issues := ReadAnnotations(mustJSON(t, schema))
	if !hasAnnotationIssue(issues, CodeUnknownExtension, "/x-sovrunn-foo") {
		t.Fatalf("expected unknown extension rejection, got %#v", issues)
	}
}

func TestReadAnnotationsUnknownFieldPolicyFieldFailsClosed(t *testing.T) {
	t.Parallel()

	schema := validAnnotatedSchema()
	props := schema["properties"].(map[string]any)
	name := props["name"].(map[string]any)
	fp := name[ExtFieldPolicy].(map[string]any)
	fp["extra"] = "nope"
	_, issues := ReadAnnotations(mustJSON(t, schema))
	if !hasAnnotationIssue(issues, CodeFieldPolicyUnknownField, "/properties/name/"+ExtFieldPolicy+"/extra") {
		t.Fatalf("expected unknown field-policy field rejection, got %#v", issues)
	}
}

func TestReadAnnotationsMissingFieldPolicyFieldFails(t *testing.T) {
	t.Parallel()

	schema := validAnnotatedSchema()
	props := schema["properties"].(map[string]any)
	name := props["name"].(map[string]any)
	fp := name[ExtFieldPolicy].(map[string]any)
	delete(fp, "auditRequired")
	_, issues := ReadAnnotations(mustJSON(t, schema))
	if !hasAnnotationIssue(issues, CodeFieldPolicyInvalid, "/properties/name/"+ExtFieldPolicy+"/auditRequired") {
		t.Fatalf("expected missing field-policy field rejection, got %#v", issues)
	}
}

func TestReadAnnotationsPropertyNameCollisionNotExtension(t *testing.T) {
	t.Parallel()

	schema := validAnnotatedSchema()
	props := schema["properties"].(map[string]any)
	props["x-sovrunn-foo"] = map[string]any{"type": "string"}
	_, issues := ReadAnnotations(mustJSON(t, schema))
	for _, iss := range issues {
		if iss.Code == CodeUnknownExtension && strings.Contains(iss.Path, "/properties/x-sovrunn-foo") {
			t.Fatalf("property identifier must not be treated as extension: %#v", issues)
		}
	}
	if len(issues) != 0 {
		t.Fatalf("expected clean pass, got %#v", issues)
	}
}

func TestReadAnnotationsRejectsNestedDocumentAnnotations(t *testing.T) {
	t.Parallel()

	schema := validAnnotatedSchema()
	props := schema["properties"].(map[string]any)
	name := props["name"].(map[string]any)
	name[ExtProfile] = "ManagedResource"
	_, issues := ReadAnnotations(mustJSON(t, schema))
	if !hasAnnotationIssue(issues, CodeAnnotationInvalid, "/properties/name/"+ExtProfile) {
		t.Fatalf("expected nested document annotation rejection, got %#v", issues)
	}
}

func hasAnnotationIssue(issues []SchemaIssue, code, path string) bool {
	for _, iss := range issues {
		if iss.Code == code && iss.Path == path {
			return true
		}
	}
	return false
}

func fieldPolicyKeys(m map[string]FieldPolicyMeta) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
