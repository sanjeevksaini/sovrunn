package apiconform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

// Task 10.1 canonical _common schema files (D-01, D-08; F12-NAMING-005/006, F12-SEC-001).
var commonSchemaFiles = commonSubSchemaFiles

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
			issues := CheckFieldPolicyCompleteness(loadCommonSchema(t, name))
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
	issues := CheckFieldPolicyCompleteness(raw)
	if !hasSchemaIssue(issues, CodeFitnessFieldPolicyMissing, "/properties/kind") {
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
	issues := CheckFieldPolicyCompleteness(raw)
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
	issues := CheckFieldPolicyCompleteness(raw)
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

func hasSchemaIssue(issues []apischema.SchemaIssue, code, path string) bool {
	for _, issue := range issues {
		if issue.Code == code && issue.Path == path {
			return true
		}
	}
	return false
}
