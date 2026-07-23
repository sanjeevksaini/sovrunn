package apiconform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

// Task 10.2 canonical contract schemas (D-01, D-08, D-17; F12-NAMING-005/006,
// F12-FIXTURE-002, F12-PROFILE-001, F12-SCOPE-002, F12-SEC-001).
var canonicalSchemaFiles = []string{
	"project.json",
	"resource-pool.json",
	"discovered-database.json",
	"plugin-definition.json",
	"adapter-configuration.json",
	"placement-evaluation-request.json",
	"operation.json",
	"audit-event.json",
}

// expectedCanonicalAnnotations maps each canonical schema file to its Matrix D
// profile / boundary / allowed-scopes / stability contract.
var expectedCanonicalAnnotations = map[string]struct {
	profile   apimeta.Profile
	boundary  apimeta.Boundary
	scopes    []apimeta.ScopeKind
	stability apimeta.Stability
}{
	"project.json": {
		profile:   apimeta.ProfileManagedResource,
		boundary:  apimeta.BoundaryCustomerFacing,
		scopes:    []apimeta.ScopeKind{apimeta.ScopeTenant},
		stability: apimeta.StabilityAlpha,
	},
	"resource-pool.json": {
		profile:   apimeta.ProfileManagedResource,
		boundary:  apimeta.BoundaryOperatorFacing,
		scopes:    []apimeta.ScopeKind{apimeta.ScopeProvider},
		stability: apimeta.StabilityAlpha,
	},
	"discovered-database.json": {
		profile:   apimeta.ProfileObservedExternalResource,
		boundary:  apimeta.BoundaryAdapterFacing,
		scopes:    []apimeta.ScopeKind{apimeta.ScopeProvider},
		stability: apimeta.StabilityAlpha,
	},
	"plugin-definition.json": {
		profile:   apimeta.ProfileVersionedDefinition,
		boundary:  apimeta.BoundaryPluginFacing,
		scopes:    []apimeta.ScopeKind{apimeta.ScopePlatform},
		stability: apimeta.StabilityAlpha,
	},
	"adapter-configuration.json": {
		profile:   apimeta.ProfileManagedResource,
		boundary:  apimeta.BoundaryAdapterFacing,
		scopes:    []apimeta.ScopeKind{apimeta.ScopeProvider},
		stability: apimeta.StabilityAlpha,
	},
	"placement-evaluation-request.json": {
		profile:   apimeta.ProfileTransientRequestResult,
		boundary:  apimeta.BoundaryInternalEngineFacing,
		scopes:    []apimeta.ScopeKind{apimeta.ScopeProject},
		stability: apimeta.StabilityAlpha,
	},
	"operation.json": {
		profile:  apimeta.ProfileLongRunningOperation,
		boundary: apimeta.BoundaryPluginFacing,
		scopes: []apimeta.ScopeKind{
			apimeta.ScopePlatform,
			apimeta.ScopeOrganization,
			apimeta.ScopeOrganizationUnit,
			apimeta.ScopeTenant,
			apimeta.ScopeProject,
			apimeta.ScopeProvider,
		},
		stability: apimeta.StabilityAlpha,
	},
	"audit-event.json": {
		profile:   apimeta.ProfileImmutableRecord,
		boundary:  apimeta.BoundaryGovernanceOnly,
		scopes:    []apimeta.ScopeKind{apimeta.ScopeOrganization},
		stability: apimeta.StabilityAlpha,
	},
}

func canonicalSchemasDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(moduleRoot(t), CanonicalSchemasDir)
}

func loadCanonicalSchema(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join(canonicalSchemasDir(t), name)
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return body
}

func TestCanonicalSchemasPresent(t *testing.T) {
	t.Parallel()

	dir := canonicalSchemasDir(t)
	for _, name := range canonicalSchemaFiles {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("required canonical schema missing: %s (%v)", path, err)
		}
	}
}

func TestCanonicalSchemasValidateSchemaSupport(t *testing.T) {
	t.Parallel()

	for _, name := range canonicalSchemaFiles {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			issues := apischema.ValidateSchemaSupport(loadCanonicalSchema(t, name))
			if len(issues) != 0 {
				t.Fatalf("ValidateSchemaSupport(%s) issues=%#v", name, issues)
			}
		})
	}
}

func TestCanonicalSchemasReadAnnotations(t *testing.T) {
	t.Parallel()

	for _, name := range canonicalSchemaFiles {
		name := name
		want := expectedCanonicalAnnotations[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			meta, issues := apischema.ReadAnnotations(loadCanonicalSchema(t, name))
			if len(issues) != 0 {
				t.Fatalf("ReadAnnotations(%s) issues=%#v", name, issues)
			}
			if meta.Profile != want.profile {
				t.Fatalf("profile: got %q want %q", meta.Profile, want.profile)
			}
			if meta.Boundary != want.boundary {
				t.Fatalf("boundary: got %q want %q", meta.Boundary, want.boundary)
			}
			if meta.Stability != want.stability {
				t.Fatalf("stability: got %q want %q", meta.Stability, want.stability)
			}
			if len(meta.AllowedScopes) != len(want.scopes) {
				t.Fatalf("allowed-scopes len: got %d want %d (%v)", len(meta.AllowedScopes), len(want.scopes), meta.AllowedScopes)
			}
			for i := range want.scopes {
				if meta.AllowedScopes[i] != want.scopes[i] {
					t.Fatalf("allowed-scopes[%d]: got %q want %q (full=%v)", i, meta.AllowedScopes[i], want.scopes[i], meta.AllowedScopes)
				}
			}
		})
	}
}

func TestCanonicalSchemasFieldPolicyCompleteness(t *testing.T) {
	t.Parallel()

	for _, name := range canonicalSchemaFiles {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			issues := checkFieldPolicyCompleteness(loadCanonicalSchema(t, name))
			if len(issues) != 0 {
				t.Fatalf("field-policy completeness failed for %s: %#v", name, issues)
			}
		})
	}
}

func TestOperationSchemaExactlySixAllowedScopes(t *testing.T) {
	t.Parallel()

	meta, issues := apischema.ReadAnnotations(loadCanonicalSchema(t, "operation.json"))
	if len(issues) != 0 {
		t.Fatalf("ReadAnnotations(operation.json) issues=%#v", issues)
	}
	want := []apimeta.ScopeKind{
		apimeta.ScopePlatform,
		apimeta.ScopeOrganization,
		apimeta.ScopeOrganizationUnit,
		apimeta.ScopeTenant,
		apimeta.ScopeProject,
		apimeta.ScopeProvider,
	}
	if len(meta.AllowedScopes) != 6 {
		t.Fatalf("Operation allowed-scopes must be exactly six, got %d: %v", len(meta.AllowedScopes), meta.AllowedScopes)
	}
	for i := range want {
		if meta.AllowedScopes[i] != want[i] {
			t.Fatalf("Operation allowed-scopes[%d]=%q want %q", i, meta.AllowedScopes[i], want[i])
		}
	}
}

func TestCanonicalSchemasUnknownExtensionFailsClosed(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "project.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	root["x-sovrunn-foo"] = "bar"
	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	_, issues := apischema.ReadAnnotations(raw)
	if !hasSchemaIssue(issues, apischema.CodeUnknownExtension, "/x-sovrunn-foo") {
		t.Fatalf("expected unknown extension rejection, got %#v", issues)
	}
}

func TestCanonicalSchemasMissingFieldPolicyFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "project.json"), &root); err != nil {
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

func TestRepositorySchemaRegistryLoadsCanonicalSchemas(t *testing.T) {
	t.Parallel()

	reg, err := NewRepositorySchemaRegistry(filepath.Join(moduleRoot(t), CanonicalSchemasDir))
	if err != nil {
		t.Fatalf("NewRepositorySchemaRegistry: %v", err)
	}
	for _, name := range canonicalSchemaFiles {
		id := CanonicalSchemasDir + "/" + name
		body, err := reg.Load(id)
		if err != nil {
			t.Fatalf("Load(%s): %v", id, err)
		}
		if len(body) == 0 {
			t.Fatalf("Load(%s): empty body", id)
		}
	}
}

func TestCanonicalProjectSchemaStructuralValidate(t *testing.T) {
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
		"apiVersion": "core.sovrunn.io/v1alpha1",
		"kind":       "Project",
		"metadata": map[string]any{
			"name": "payments-production",
			"scopeRef": map[string]any{
				"apiVersion": "tenancy.sovrunn.io/v1alpha1",
				"kind":       "Tenant",
				"name":       "acme-tenant",
				"uid":        "opaque-tenant-uid",
			},
		},
		"spec": map[string]any{
			"description": "Payments production project",
		},
	}
	violations, err := v.Validate(instance, CanonicalSchemasDir+"/project.json")
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("unexpected violations: %#v", violations)
	}
}
