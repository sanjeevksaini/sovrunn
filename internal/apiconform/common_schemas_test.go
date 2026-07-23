package apiconform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

// Task 10.1 canonical _common schema files (D-01, D-08; F12-NAMING-005/006, F12-SEC-001).
var commonSchemaFiles = []string{
	"type-meta.json",
	"object-meta.json",
	"typed-ref.json",
	"scope-ref.json",
	"owner-ref.json",
	"condition.json",
	"problem.json",
	"violation.json",
	"page.json",
}

func commonSchemasDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(moduleRoot(t), CanonicalSchemasDir, "_common")
}

func loadCommonSchema(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join(commonSchemasDir(t), name)
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return body
}

func TestCommonSchemasPresent(t *testing.T) {
	t.Parallel()

	dir := commonSchemasDir(t)
	for _, name := range commonSchemaFiles {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("required _common schema missing: %s (%v)", path, err)
		}
	}
}

func TestCommonSchemasValidateSchemaSupport(t *testing.T) {
	t.Parallel()

	for _, name := range commonSchemaFiles {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			issues := apischema.ValidateSchemaSupport(loadCommonSchema(t, name))
			if len(issues) != 0 {
				t.Fatalf("ValidateSchemaSupport(%s) issues=%#v", name, issues)
			}
		})
	}
}

func TestCommonSchemasFieldPolicyCompleteness(t *testing.T) {
	t.Parallel()

	for _, name := range commonSchemaFiles {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			issues := checkFieldPolicyCompleteness(loadCommonSchema(t, name))
			if len(issues) != 0 {
				t.Fatalf("field-policy completeness failed for %s: %#v", name, issues)
			}
		})
	}
}

func TestCommonSchemasMissingFieldPolicyFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCommonSchema(t, "type-meta.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	props := root["properties"].(map[string]any)
	kind := props["kind"].(map[string]any)
	delete(kind, apischema.ExtFieldPolicy)

	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	issues := checkFieldPolicyCompleteness(raw)
	if !hasSchemaIssue(issues, "FIELD_POLICY_MISSING", "/properties/kind") {
		t.Fatalf("expected FIELD_POLICY_MISSING at /properties/kind, got %#v", issues)
	}
}

func TestCommonSchemasUnknownFieldPolicyFieldFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCommonSchema(t, "type-meta.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	props := root["properties"].(map[string]any)
	kind := props["kind"].(map[string]any)
	fp := kind[apischema.ExtFieldPolicy].(map[string]any)
	fp["extra"] = "nope"

	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	issues := checkFieldPolicyCompleteness(raw)
	if !hasSchemaIssue(issues, apischema.CodeFieldPolicyUnknownField, "/properties/kind/"+apischema.ExtFieldPolicy+"/extra") {
		t.Fatalf("expected unknown field-policy field rejection, got %#v", issues)
	}
}

func TestCommonSchemasMissingFieldPolicyKeyFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCommonSchema(t, "page.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	props := root["properties"].(map[string]any)
	tok := props["nextPageToken"].(map[string]any)
	fp := tok[apischema.ExtFieldPolicy].(map[string]any)
	delete(fp, "auditRequired")

	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	issues := checkFieldPolicyCompleteness(raw)
	if !hasSchemaIssue(issues, apischema.CodeFieldPolicyInvalid, "/properties/nextPageToken/"+apischema.ExtFieldPolicy+"/auditRequired") {
		t.Fatalf("expected missing auditRequired rejection, got %#v", issues)
	}
}

func TestRepositorySchemaRegistryLoadsCommonSchemas(t *testing.T) {
	t.Parallel()

	reg, err := NewRepositorySchemaRegistry(filepath.Join(moduleRoot(t), CanonicalSchemasDir))
	if err != nil {
		t.Fatalf("NewRepositorySchemaRegistry: %v", err)
	}
	for _, name := range commonSchemaFiles {
		id := CanonicalSchemasDir + "/_common/" + name
		body, err := reg.Load(id)
		if err != nil {
			t.Fatalf("Load(%s): %v", id, err)
		}
		if len(body) == 0 {
			t.Fatalf("Load(%s): empty body", id)
		}
	}
}

func TestCommonObjectMetaScopeRefResolves(t *testing.T) {
	t.Parallel()

	reg, err := NewRepositorySchemaRegistry(filepath.Join(moduleRoot(t), CanonicalSchemasDir))
	if err != nil {
		t.Fatalf("NewRepositorySchemaRegistry: %v", err)
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

	instance := map[string]any{
		"name": "payments-production",
		"scopeRef": map[string]any{
			"apiVersion": "tenancy.sovrunn.io/v1alpha1",
			"kind":       "Tenant",
			"name":       "acme-tenant",
			"uid":        "opaque-tenant-uid",
		},
	}
	violations, err := v.Validate(instance, CanonicalSchemasDir+"/_common/object-meta.json")
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("unexpected violations: %#v", violations)
	}
}

// checkFieldPolicyCompleteness asserts every named property under "properties"
// explicitly declares a complete x-sovrunn-field-policy (D-08, F12-SEC-001).
// FEATURE-0012 uses no inheritance algorithm: policies must be present on each
// property schema object in the document under test.
func checkFieldPolicyCompleteness(schema []byte) []apischema.SchemaIssue {
	var root any
	if err := json.Unmarshal(schema, &root); err != nil {
		return []apischema.SchemaIssue{{
			Path:    "/",
			Code:    apischema.CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}
	var issues []apischema.SchemaIssue
	walkFieldPolicyCompleteness(root, "", &issues)
	sort.SliceStable(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		return issues[i].Code < issues[j].Code
	})
	return issues
}

func walkFieldPolicyCompleteness(node any, path string, issues *[]apischema.SchemaIssue) {
	obj, ok := node.(map[string]any)
	if !ok {
		return
	}

	if rawProps, hasProps := obj["properties"]; hasProps {
		props, ok := rawProps.(map[string]any)
		if !ok {
			*issues = append(*issues, apischema.SchemaIssue{
				Path:    joinTestPointer(path, "properties"),
				Code:    apischema.CodeMalformedSchema,
				Message: "properties must be an object",
			})
			return
		}
		names := make([]string, 0, len(props))
		for name := range props {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			propPath := joinTestPointer(joinTestPointer(path, "properties"), name)
			propSchema, ok := props[name].(map[string]any)
			if !ok {
				*issues = append(*issues, apischema.SchemaIssue{
					Path:    propPath,
					Code:    apischema.CodeMalformedSchema,
					Message: "property schema must be an object",
				})
				continue
			}
			rawPolicy, hasPolicy := propSchema[apischema.ExtFieldPolicy]
			if !hasPolicy {
				*issues = append(*issues, apischema.SchemaIssue{
					Path:    propPath,
					Code:    "FIELD_POLICY_MISSING",
					Message: "boundary-crossing property missing explicit x-sovrunn-field-policy",
				})
			} else if rawPolicy == nil {
				*issues = append(*issues, apischema.SchemaIssue{
					Path:    joinTestPointer(propPath, apischema.ExtFieldPolicy),
					Code:    apischema.CodeFieldPolicyInvalid,
					Message: "x-sovrunn-field-policy must be an object",
				})
			} else {
				// Reuse ReadAnnotations vocabulary checks by wrapping a tiny
				// annotated document around this property schema.
				probe := map[string]any{
					"type":                     "object",
					"x-sovrunn-profile":        "EmbeddedValue",
					"x-sovrunn-boundary":       "customer-facing",
					"x-sovrunn-allowed-scopes": []any{"Platform"},
					"x-sovrunn-stability":      "stable",
					"properties": map[string]any{
						"probe": propSchema,
					},
				}
				raw, err := json.Marshal(probe)
				if err != nil {
					*issues = append(*issues, apischema.SchemaIssue{
						Path:    propPath,
						Code:    apischema.CodeMalformedSchema,
						Message: "could not marshal field-policy probe",
					})
				} else {
					_, annIssues := apischema.ReadAnnotations(raw)
					for _, issue := range annIssues {
						// Remap probe paths back to the real property path.
						mapped := issue
						if len(issue.Path) >= len("/properties/probe") &&
							issue.Path[:len("/properties/probe")] == "/properties/probe" {
							mapped.Path = propPath + issue.Path[len("/properties/probe"):]
						}
						*issues = append(*issues, mapped)
					}
				}
			}
			walkFieldPolicyCompleteness(propSchema, propPath, issues)
		}
	}

	if rawItems, hasItems := obj["items"]; hasItems {
		walkFieldPolicyCompleteness(rawItems, joinTestPointer(path, "items"), issues)
	}
	if rawAP, hasAP := obj["additionalProperties"]; hasAP {
		if _, isBool := rawAP.(bool); !isBool {
			walkFieldPolicyCompleteness(rawAP, joinTestPointer(path, "additionalProperties"), issues)
		}
	}
}

func joinTestPointer(base, key string) string {
	if base == "" || base == "/" {
		return "/" + escapeJSONPointerToken(key)
	}
	return base + "/" + escapeJSONPointerToken(key)
}

func escapeJSONPointerToken(s string) string {
	r := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '~':
			r = append(r, '~', '0')
		case '/':
			r = append(r, '~', '1')
		default:
			r = append(r, s[i])
		}
	}
	return string(r)
}

func hasSchemaIssue(issues []apischema.SchemaIssue, code, path string) bool {
	for _, issue := range issues {
		if issue.Code == code && issue.Path == path {
			return true
		}
	}
	return false
}
