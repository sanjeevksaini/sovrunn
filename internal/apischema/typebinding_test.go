package apischema

import (
	"reflect"
	"testing"
)

// correctSample is a deliberately small Go type that matches sampleSchema.
type correctSample struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Name       string            `json:"name"`
	Labels     map[string]string `json:"labels,omitempty"`
	Count      int64             `json:"count,omitempty"`
	Enabled    bool              `json:"enabled,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	Status     sampleStatus      `json:"status,omitempty"`
	Nested     sampleNested      `json:"nested,omitempty"`
}

type sampleStatus string

type sampleNested struct {
	Value string `json:"value"`
}

// embeddedParent promotes EmbeddedMeta fields into the JSON object.
type embeddedParent struct {
	EmbeddedMeta
	Name string `json:"name"`
}

type EmbeddedMeta struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

// mismatchedSample deliberately diverges from sampleSchema.
type mismatchedSample struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	// Wrong JSON tag name for schema property "name".
	Title string `json:"title"`
	// Extra field not in schema.
	Extra string `json:"extra,omitempty"`
	// Required-looking optional without omitempty/pointer for schema-optional count.
	Count int64 `json:"count"`
	// Wrong primitive type for enabled.
	Enabled string `json:"enabled,omitempty"`
}

const sampleSchema = `{
  "type": "object",
  "properties": {
    "apiVersion": { "type": "string" },
    "kind": { "type": "string" },
    "name": { "type": "string" },
    "labels": {
      "type": "object",
      "additionalProperties": { "type": "string" }
    },
    "count": { "type": "integer" },
    "enabled": { "type": "boolean" },
    "tags": {
      "type": "array",
      "items": { "type": "string" }
    },
    "status": {
      "type": "string",
      "enum": ["Ready", "Pending"]
    },
    "nested": {
      "type": "object",
      "properties": {
        "value": { "type": "string" }
      },
      "required": ["value"],
      "additionalProperties": false
    }
  },
  "required": ["apiVersion", "kind", "name"],
  "additionalProperties": false
}`

const embeddedSchema = `{
  "type": "object",
  "properties": {
    "apiVersion": { "type": "string" },
    "kind": { "type": "string" },
    "name": { "type": "string" }
  },
  "required": ["apiVersion", "kind", "name"],
  "additionalProperties": false
}`

const refPropertySchema = `{
  "type": "object",
  "properties": {
    "metadata": {
      "$ref": "_common/object-meta.json"
    }
  },
  "required": ["metadata"],
  "additionalProperties": false
}`

func TestVerifyGoTypeAgainstSchemaAcceptsMatchingType(t *testing.T) {
	t.Parallel()

	issues := VerifyGoTypeAgainstSchema([]byte(sampleSchema), reflect.TypeOf(correctSample{}))
	if len(issues) != 0 {
		t.Fatalf("matching Go type must pass, got %#v", issues)
	}
}

func TestVerifyGoTypeAgainstSchemaAcceptsPointerType(t *testing.T) {
	t.Parallel()

	issues := VerifyGoTypeAgainstSchema([]byte(sampleSchema), reflect.TypeOf(&correctSample{}))
	if len(issues) != 0 {
		t.Fatalf("pointer to matching Go type must pass, got %#v", issues)
	}
}

func TestVerifyGoTypeAgainstSchemaAcceptsEmbeddedPromotion(t *testing.T) {
	t.Parallel()

	issues := VerifyGoTypeAgainstSchema([]byte(embeddedSchema), reflect.TypeOf(embeddedParent{}))
	if len(issues) != 0 {
		t.Fatalf("embedded promoted fields must pass, got %#v", issues)
	}
}

func TestVerifyGoTypeAgainstSchemaRejectsMismatchedType(t *testing.T) {
	t.Parallel()

	issues := VerifyGoTypeAgainstSchema([]byte(sampleSchema), reflect.TypeOf(mismatchedSample{}))
	if len(issues) == 0 {
		t.Fatal("deliberately mismatched Go type must fail")
	}

	if !hasIssue(issues, CodeGoFieldMissing, "/name") {
		t.Fatalf("expected GO_FIELD_MISSING for /name, got %#v", issues)
	}
	if !hasIssue(issues, CodeGoFieldExtra, "/title") && !hasIssue(issues, CodeGoFieldExtra, "/extra") {
		// title is extra relative to schema; extra is also extra
		t.Fatalf("expected GO_FIELD_EXTRA for undeclared fields, got %#v", issues)
	}
	if !hasIssue(issues, CodeGoRequiredMismatch, "/count") {
		t.Fatalf("expected GO_REQUIRED_MISMATCH for optional /count without omitempty, got %#v", issues)
	}
	if !hasIssue(issues, CodeGoTypeMismatch, "/enabled") {
		t.Fatalf("expected GO_TYPE_MISMATCH for /enabled, got %#v", issues)
	}
}

func TestVerifyGoTypeAgainstSchemaRejectsMissingRequiredOmitempty(t *testing.T) {
	t.Parallel()

	type badRequired struct {
		APIVersion string `json:"apiVersion,omitempty"`
		Kind       string `json:"kind"`
		Name       string `json:"name"`
	}
	schema := []byte(`{
		"type": "object",
		"properties": {
			"apiVersion": { "type": "string" },
			"kind": { "type": "string" },
			"name": { "type": "string" }
		},
		"required": ["apiVersion", "kind", "name"],
		"additionalProperties": false
	}`)
	issues := VerifyGoTypeAgainstSchema(schema, reflect.TypeOf(badRequired{}))
	if !hasIssue(issues, CodeGoRequiredMismatch, "/apiVersion") {
		t.Fatalf("required field with omitempty must fail, got %#v", issues)
	}
}

func TestVerifyGoTypeAgainstSchemaRejectsMapForClosedObject(t *testing.T) {
	t.Parallel()

	issues := VerifyGoTypeAgainstSchema([]byte(embeddedSchema), reflect.TypeOf(map[string]string{}))
	if !hasIssue(issues, CodeGoTypeMismatch, "/") {
		t.Fatalf("map must not satisfy closed object schema, got %#v", issues)
	}
}

func TestVerifyGoTypeAgainstSchemaAcceptsAdditionalPropertiesMap(t *testing.T) {
	t.Parallel()

	schema := []byte(`{
		"type": "object",
		"additionalProperties": { "type": "string" }
	}`)
	issues := VerifyGoTypeAgainstSchema(schema, reflect.TypeOf(map[string]string{}))
	if len(issues) != 0 {
		t.Fatalf("map[string]string must match additionalProperties string schema, got %#v", issues)
	}

	bad := VerifyGoTypeAgainstSchema(schema, reflect.TypeOf(map[string]int{}))
	if !hasIssue(bad, CodeGoTypeMismatch, "/") {
		t.Fatalf("map[string]int must fail string additionalProperties, got %#v", bad)
	}
}

func TestVerifyGoTypeAgainstSchemaRefPropertyRequiresStruct(t *testing.T) {
	t.Parallel()

	type withMeta struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
	}
	issues := VerifyGoTypeAgainstSchema([]byte(refPropertySchema), reflect.TypeOf(withMeta{}))
	if len(issues) != 0 {
		t.Fatalf("$ref property with struct field must pass shallowly, got %#v", issues)
	}

	type withStringMeta struct {
		Metadata string `json:"metadata"`
	}
	bad := VerifyGoTypeAgainstSchema([]byte(refPropertySchema), reflect.TypeOf(withStringMeta{}))
	if !hasIssue(bad, CodeGoTypeMismatch, "/metadata") {
		t.Fatalf("$ref property with string field must fail, got %#v", bad)
	}
}

func TestVerifyGoTypeAgainstSchemaNilAndMalformed(t *testing.T) {
	t.Parallel()

	if issues := VerifyGoTypeAgainstSchema(nil, reflect.TypeOf(correctSample{})); !hasIssue(issues, CodeMalformedSchema, "/") {
		t.Fatalf("nil schema must be malformed, got %#v", issues)
	}
	if issues := VerifyGoTypeAgainstSchema([]byte(sampleSchema), nil); !hasIssue(issues, CodeGoTypeMismatch, "/") {
		t.Fatalf("nil Go type must fail, got %#v", issues)
	}
	if issues := VerifyGoTypeAgainstSchema([]byte(`{`), reflect.TypeOf(correctSample{})); !hasIssue(issues, CodeMalformedSchema, "/") {
		t.Fatalf("invalid JSON schema must be malformed, got %#v", issues)
	}
}

func TestTypeBindingHoldsSchemaPathAndGoType(t *testing.T) {
	t.Parallel()

	binding := TypeBinding{
		SchemaPath: "api/schemas/_common/page.json",
		GoType: reflect.TypeOf(struct {
			NextPageToken string `json:"nextPageToken,omitempty"`
		}{}),
	}
	if binding.SchemaPath == "" || binding.GoType == nil {
		t.Fatal("TypeBinding must carry SchemaPath and GoType")
	}
}
