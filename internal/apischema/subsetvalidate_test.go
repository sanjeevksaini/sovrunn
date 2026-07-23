package apischema

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSupportedKeywordsExactVocabulary(t *testing.T) {
	t.Parallel()

	want := make(map[string]struct{}, len(CoreSupportedKeywords)+len(RegisteredExtensionKeywords))
	for _, k := range CoreSupportedKeywords {
		want[k] = struct{}{}
	}
	for _, k := range RegisteredExtensionKeywords {
		want[k] = struct{}{}
	}

	if len(SupportedKeywords) != len(want) {
		t.Fatalf("SupportedKeywords len = %d, want %d", len(SupportedKeywords), len(want))
	}
	for k := range want {
		if !IsSupportedKeyword(k) {
			t.Fatalf("missing supported keyword %q", k)
		}
	}
	for k := range SupportedKeywords {
		if _, ok := want[k]; !ok {
			t.Fatalf("unexpected supported keyword %q", k)
		}
	}

	// $defs is explicitly prohibited and must never be in the supported set.
	if IsSupportedKeyword("$defs") {
		t.Fatal("$defs must not be a supported keyword")
	}
	if IsSupportedKeyword("oneOf") || IsSupportedKeyword("if") || IsSupportedKeyword("then") || IsSupportedKeyword("else") {
		t.Fatal("combinators must not be supported keywords")
	}

	for _, k := range RegisteredExtensionKeywords {
		if !IsRegisteredExtension(k) {
			t.Fatalf("registered extension %q must be recognized", k)
		}
	}
	if IsRegisteredExtension("x-sovrunn-foo") {
		t.Fatal("unknown x-sovrunn-* must not be a registered extension")
	}
}

func TestValidateSchemaSupportAcceptsDocumentMetadata(t *testing.T) {
	t.Parallel()

	schema := []byte(`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$id": "https://sovrunn.io/schemas/example.json",
		"title": "Example",
		"description": "Document metadata keywords are accepted",
		"type": "object"
	}`)
	if issues := ValidateSchemaSupport(schema); len(issues) != 0 {
		t.Fatalf("metadata schema must pass, got %#v", issues)
	}
}

func TestValidateSchemaSupportAcceptsSupportedSubset(t *testing.T) {
	t.Parallel()

	schema := []byte(`{
		"type": "object",
		"properties": {
			"name": {
				"type": "string",
				"minLength": 1,
				"maxLength": 63,
				"pattern": "^[a-z0-9-]+$",
				"default": "demo",
				"examples": ["demo", "payments"],
				"x-sovrunn-field-policy": {
					"classification": "Public",
					"authorizedWriter": "creator",
					"authorizedReaders": ["customer"],
					"mutability": "immutable",
					"retention": "standard",
					"redaction": "none",
					"residency": "any",
					"auditRequired": true
				}
			},
			"count": {
				"type": "integer",
				"minimum": 0,
				"maximum": 100
			},
			"items": {
				"type": "array",
				"items": { "type": "string" }
			},
			"tags": {
				"type": "object",
				"additionalProperties": { "type": "string" }
			},
			"phase": {
				"type": "string",
				"enum": ["Pending", "Ready"]
			},
			"ref": { "$ref": "../_common/typed-ref.json" }
		},
		"required": ["name"],
		"additionalProperties": false,
		"x-sovrunn-profile": "ManagedResource",
		"x-sovrunn-boundary": "customer-facing",
		"x-sovrunn-allowed-scopes": ["Tenant"],
		"x-sovrunn-stability": "alpha"
	}`)
	if issues := ValidateSchemaSupport(schema); len(issues) != 0 {
		t.Fatalf("supported-subset schema must pass, got %#v", issues)
	}
}

func TestValidateSchemaSupportRejectsUnsupportedKeywords(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		raw  string
		path string
	}{
		{
			name: "oneOf",
			raw:  `{"oneOf":[{"type":"string"},{"type":"integer"}]}`,
			path: "/oneOf",
		},
		{
			name: "if",
			raw:  `{"if":{"type":"object"},"then":{"required":["a"]},"else":{"required":["b"]}}`,
			path: "/if",
		},
		{
			name: "defs",
			raw:  `{"$defs":{"Name":{"type":"string"}}}`,
			path: "/$defs",
		},
		{
			name: "allOf",
			raw:  `{"allOf":[{"type":"object"}]}`,
			path: "/allOf",
		},
		{
			name: "format",
			raw:  `{"type":"string","format":"email"}`,
			path: "/format",
		},
		{
			name: "unknown_x_sovrunn",
			raw:  `{"x-sovrunn-foo":true}`,
			path: "/x-sovrunn-foo",
		},
		{
			name: "nested_anyOf",
			raw:  `{"type":"object","properties":{"v":{"anyOf":[{"type":"string"}]}}}`,
			path: "/properties/v/anyOf",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			issues := ValidateSchemaSupport([]byte(tc.raw))
			if !hasIssue(issues, CodeUnsupportedKeyword, tc.path) {
				t.Fatalf("want %s at %s, got %#v", CodeUnsupportedKeyword, tc.path, issues)
			}
		})
	}
}

func TestValidateSchemaSupportPropertyNamesAreNotKeywords(t *testing.T) {
	t.Parallel()

	// Property identifiers may collide with unsupported keyword names; they
	// must not be rejected as schema keywords.
	schema := []byte(`{
		"type": "object",
		"properties": {
			"oneOf": { "type": "string" },
			"$defs": { "type": "integer" },
			"if": { "type": "boolean" },
			"format": { "type": "string" }
		}
	}`)
	issues := ValidateSchemaSupport(schema)
	if len(issues) != 0 {
		t.Fatalf("property names must not be treated as keywords, got %#v", issues)
	}
}

func TestValidateSchemaSupportExtensionObjectFieldsAreNotKeywords(t *testing.T) {
	t.Parallel()

	// Fields inside x-sovrunn-field-policy are extension-schema fields, not
	// core JSON Schema keywords — even when they reuse keyword-looking names.
	schema := []byte(`{
		"type": "object",
		"properties": {
			"name": {
				"type": "string",
				"x-sovrunn-field-policy": {
					"classification": "Public",
					"authorizedWriter": "creator",
					"authorizedReaders": ["customer"],
					"mutability": "immutable",
					"retention": "standard",
					"redaction": "none",
					"residency": "any",
					"auditRequired": true,
					"oneOf": "not-a-schema-keyword-here",
					"$defs": "also-not-a-keyword-here"
				}
			}
		},
		"x-sovrunn-profile": "ManagedResource",
		"x-sovrunn-boundary": "customer-facing",
		"x-sovrunn-allowed-scopes": ["Tenant"],
		"x-sovrunn-stability": "alpha"
	}`)
	issues := ValidateSchemaSupport(schema)
	if len(issues) != 0 {
		t.Fatalf("extension-object fields must not be treated as keywords, got %#v", issues)
	}
}

func TestValidateSchemaSupportFailClosedNoSilentIgnore(t *testing.T) {
	t.Parallel()

	schema := []byte(`{
		"type": "object",
		"unevaluatedProperties": false,
		"properties": {
			"v": {
				"type": "string",
				"contentEncoding": "base64"
			}
		}
	}`)
	issues := ValidateSchemaSupport(schema)
	if !hasIssue(issues, CodeUnsupportedKeyword, "/unevaluatedProperties") {
		t.Fatalf("unevaluatedProperties must be rejected, got %#v", issues)
	}
	if !hasIssue(issues, CodeUnsupportedKeyword, "/properties/v/contentEncoding") {
		t.Fatalf("contentEncoding must be rejected, got %#v", issues)
	}
	if len(issues) < 2 {
		t.Fatalf("fail-closed must report every unsupported keyword, got %#v", issues)
	}
}

func TestValidateSchemaSupportMalformedJSON(t *testing.T) {
	t.Parallel()

	issues := ValidateSchemaSupport([]byte(`{not-json`))
	if !hasIssue(issues, CodeMalformedSchema, "/") {
		t.Fatalf("malformed JSON must yield %s at /, got %#v", CodeMalformedSchema, issues)
	}

	issues = ValidateSchemaSupport(nil)
	if !hasIssue(issues, CodeMalformedSchema, "/") {
		t.Fatalf("empty schema must yield %s at /, got %#v", CodeMalformedSchema, issues)
	}

	issues = ValidateSchemaSupport([]byte(`["not-a-schema"]`))
	if !hasIssue(issues, CodeMalformedSchema, "/") {
		t.Fatalf("array root must yield %s at /, got %#v", CodeMalformedSchema, issues)
	}
}

func TestValidateSchemaSupportBooleanSchema(t *testing.T) {
	t.Parallel()

	if issues := ValidateSchemaSupport([]byte(`true`)); len(issues) != 0 {
		t.Fatalf("boolean schema true must pass, got %#v", issues)
	}
	if issues := ValidateSchemaSupport([]byte(`false`)); len(issues) != 0 {
		t.Fatalf("boolean schema false must pass, got %#v", issues)
	}
}

func TestValidateSchemaSupportJSONPointerEscaping(t *testing.T) {
	t.Parallel()

	// Property name with / and ~ must escape in the nested unsupported-keyword path.
	schema := []byte(`{
		"type": "object",
		"properties": {
			"a/b~c": { "not": {} }
		}
	}`)
	issues := ValidateSchemaSupport(schema)
	wantPath := "/properties/a~1b~0c/not"
	if !hasIssue(issues, CodeUnsupportedKeyword, wantPath) {
		t.Fatalf("want %s at %s, got %#v", CodeUnsupportedKeyword, wantPath, issues)
	}
}

func TestValidateSchemaSupportDoesNotImportApiproblem(t *testing.T) {
	t.Parallel()

	// Compile-time / package-boundary reminder: this package must stay free of
	// apiproblem. The apiconform imports_test enforces the full matrix; this
	// local check documents the SchemaIssue boundary for task 5.1.
	raw, err := json.Marshal(SchemaIssue{
		Path:    "/oneOf",
		Code:    CodeUnsupportedKeyword,
		Message: "unsupported schema keyword \"oneOf\"",
	})
	if err != nil {
		t.Fatalf("SchemaIssue must marshal without apiproblem: %v", err)
	}
	if !strings.Contains(string(raw), CodeUnsupportedKeyword) {
		t.Fatalf("unexpected marshal output: %s", raw)
	}
}

func hasIssue(issues []SchemaIssue, code, path string) bool {
	for _, iss := range issues {
		if iss.Code == code && iss.Path == path {
			return true
		}
	}
	return false
}
