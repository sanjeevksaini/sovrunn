package validation

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateProviderDatacenter_PositiveOneValidParentRef(t *testing.T) {
	t.Parallel()

	raw := mustProviderDatacenterJSON(t, validProviderDatacenterDoc("dc-bangalore-1", "loc-in"))
	if prob := ValidateProviderDatacenter(context.Background(), raw, testProviderDatacenterStructural(t)); prob != nil {
		t.Fatalf("one valid providerLocationRef must pass: %#v", prob)
	}
}

func TestValidateProviderDatacenter_NegativeMissingReference(t *testing.T) {
	t.Parallel()

	raw := mustProviderDatacenterJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderDatacenter,
		"metadata": map[string]any{
			"name": "dc-missing-ref",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{},
	})
	prob := ValidateProviderDatacenter(context.Background(), raw, testProviderDatacenterStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/providerLocationRef", apiproblem.ViolationCode(apischema.CodeRequiredField))
}

func TestValidateProviderDatacenter_NegativeTwoReferences(t *testing.T) {
	t.Parallel()

	// Cardinality is exactly one singular object. An array is rejected at
	// decode (type mismatch) before semantic checks.
	raw := mustProviderDatacenterJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderDatacenter,
		"metadata": map[string]any{
			"name": "dc-two-refs",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{
			"providerLocationRef": []any{
				map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       resources.KindProviderLocation,
					"name":       "loc-a",
				},
				map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       resources.KindProviderLocation,
					"name":       "loc-b",
				},
			},
		},
	})
	prob := ValidateProviderDatacenter(context.Background(), raw, testProviderDatacenterStructural(t))
	if prob == nil {
		t.Fatal("two providerLocationRef values must be rejected")
	}
	assertProblemViolation(t, prob, apiproblem.CodeMalformedRequest, "/spec/providerLocationRef", apiproblem.ViolationCode(apiproblem.CodeMalformedRequest))
}

func TestValidateProviderDatacenter_NegativeWrongKindReference(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		parentKind string
	}{
		{name: "provider-skipped-level", parentKind: resources.KindProvider},
		{name: "failure-domain-skipped-level", parentKind: resources.KindDatacenterFailureDomain},
		{name: "datacenter-self-kind", parentKind: resources.KindProviderDatacenter},
	}

	structural := testProviderDatacenterStructural(t)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			raw := mustProviderDatacenterJSON(t, map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       resources.KindProviderDatacenter,
				"metadata": map[string]any{
					"name": "dc-wrong-kind",
					"scopeRef": map[string]any{
						"apiVersion": resources.FabricAPIVersion,
						"kind":       string(apimeta.ScopeProvider),
						"name":       "sovereign-provider-a",
					},
				},
				"spec": map[string]any{
					"providerLocationRef": map[string]any{
						"apiVersion": resources.FabricAPIVersion,
						"kind":       tc.parentKind,
						"name":       "not-a-location",
					},
				},
			})
			prob := ValidateProviderDatacenter(context.Background(), raw, structural)
			assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/providerLocationRef/kind", apiproblem.ViolationCode(apiref.CodeKindNotAllowed))
		})
	}
}

func TestValidateProviderDatacenter_NegativeWrongScopeKind(t *testing.T) {
	t.Parallel()

	raw := mustProviderDatacenterJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderDatacenter,
		"metadata": map[string]any{
			"name": "dc-wrong-scope",
			"scopeRef": map[string]any{
				"apiVersion": "core.sovrunn.io/v1alpha1",
				"kind":       string(apimeta.ScopeOrganization),
				"name":       "acme",
			},
		},
		"spec": map[string]any{
			"providerLocationRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       resources.KindProviderLocation,
				"name":       "loc-in",
			},
		},
	})
	prob := ValidateProviderDatacenter(context.Background(), raw, testProviderDatacenterStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/scopeRef/kind", apiproblem.ViolationCode(apiref.CodeScopeNotAllowed))
}

func TestValidateProviderDatacenter_NegativeUserAuthoredStatus(t *testing.T) {
	t.Parallel()

	doc := validProviderDatacenterDoc("dc-with-status", "loc-in")
	doc["status"] = map[string]any{"observedGeneration": 1}
	raw := mustProviderDatacenterJSON(t, doc)
	prob := ValidateProviderDatacenter(context.Background(), raw, testProviderDatacenterStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/status", violationValidationFail)
}

func TestValidateProviderDatacenter_BoundaryReferenceCardinalityExactlyOne(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	structural := testProviderDatacenterStructural(t)

	t.Run("exactly-one-accepted", func(t *testing.T) {
		t.Parallel()
		raw := mustProviderDatacenterJSON(t, validProviderDatacenterDoc("dc-card-one", "loc-in"))
		if prob := ValidateProviderDatacenter(ctx, raw, structural); prob != nil {
			t.Fatalf("exactly one providerLocationRef must pass: %#v", prob)
		}
	})

	t.Run("zero-rejected", func(t *testing.T) {
		t.Parallel()
		raw := mustProviderDatacenterJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindProviderDatacenter,
			"metadata": map[string]any{
				"name": "dc-card-zero",
				"scopeRef": map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       string(apimeta.ScopeProvider),
					"name":       "sovereign-provider-a",
				},
			},
			"spec": map[string]any{},
		})
		prob := ValidateProviderDatacenter(ctx, raw, structural)
		assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/providerLocationRef", apiproblem.ViolationCode(apischema.CodeRequiredField))
	})

	t.Run("multiple-rejected", func(t *testing.T) {
		t.Parallel()
		raw := mustProviderDatacenterJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindProviderDatacenter,
			"metadata": map[string]any{
				"name": "dc-card-multi",
				"scopeRef": map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       string(apimeta.ScopeProvider),
					"name":       "sovereign-provider-a",
				},
			},
			"spec": map[string]any{
				"providerLocationRef": []any{
					map[string]any{
						"apiVersion": resources.FabricAPIVersion,
						"kind":       resources.KindProviderLocation,
						"name":       "loc-a",
					},
					map[string]any{
						"apiVersion": resources.FabricAPIVersion,
						"kind":       resources.KindProviderLocation,
						"name":       "loc-b",
					},
				},
			},
		})
		prob := ValidateProviderDatacenter(ctx, raw, structural)
		if prob == nil {
			t.Fatal("multiple providerLocationRef values must be rejected")
		}
		assertProblemViolation(t, prob, apiproblem.CodeMalformedRequest, "/spec/providerLocationRef", apiproblem.ViolationCode(apiproblem.CodeMalformedRequest))
	})
}

func validProviderDatacenterDoc(name, locationName string) map[string]any {
	return map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderDatacenter,
		"metadata": map[string]any{
			"name": name,
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{
			"providerLocationRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       resources.KindProviderLocation,
				"name":       locationName,
			},
		},
	}
}

func mustProviderDatacenterJSON(t *testing.T, doc map[string]any) []byte {
	t.Helper()
	return mustProviderLocationJSON(t, doc)
}

// testProviderDatacenterStructural builds the existing apiconform structural
// validator for tests only. Production ValidateProviderDatacenter depends
// solely on the injected apivalid.StructuralValidator interface.
func testProviderDatacenterStructural(t *testing.T) apivalid.StructuralValidator {
	t.Helper()

	root := providerLocationModuleRoot(t)
	reg, err := apiconform.NewRepositorySchemaRegistry(filepath.Join(root, apiconform.CanonicalSchemasDir))
	if err != nil {
		t.Fatalf("NewRepositorySchemaRegistry: %v", err)
	}
	resolver, err := apiconform.NewLocalRefResolver(reg, apiconform.DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}
	cfg, err := apiconform.NewStructuralValidatorConfig(reg, resolver)
	if err != nil {
		t.Fatalf("NewStructuralValidatorConfig: %v", err)
	}
	v, err := apiconform.NewStructuralValidator(cfg)
	if err != nil {
		t.Fatalf("NewStructuralValidator: %v", err)
	}
	return v
}
