package validation

import (
	"context"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateDatacenterFailureDomain_PositiveOneValidParentRef(t *testing.T) {
	t.Parallel()

	raw := mustDatacenterFailureDomainJSON(t, validDatacenterFailureDomainDoc("fd-bangalore-az1", "dc-bangalore-1"))
	if prob := ValidateDatacenterFailureDomain(context.Background(), raw, testDatacenterFailureDomainStructural(t)); prob != nil {
		t.Fatalf("one valid providerDatacenterRef must pass: %#v", prob)
	}
}

func TestValidateDatacenterFailureDomain_NegativeMissingReference(t *testing.T) {
	t.Parallel()

	raw := mustDatacenterFailureDomainJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindDatacenterFailureDomain,
		"metadata": map[string]any{
			"name": "fd-missing-ref",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{},
	})
	prob := ValidateDatacenterFailureDomain(context.Background(), raw, testDatacenterFailureDomainStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/providerDatacenterRef", apiproblem.ViolationCode(apischema.CodeRequiredField))
}

func TestValidateDatacenterFailureDomain_NegativeTwoReferences(t *testing.T) {
	t.Parallel()

	// Cardinality is exactly one singular object. An array is rejected at
	// decode (type mismatch) before semantic checks.
	raw := mustDatacenterFailureDomainJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindDatacenterFailureDomain,
		"metadata": map[string]any{
			"name": "fd-two-refs",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{
			"providerDatacenterRef": []any{
				map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       resources.KindProviderDatacenter,
					"name":       "dc-a",
				},
				map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       resources.KindProviderDatacenter,
					"name":       "dc-b",
				},
			},
		},
	})
	prob := ValidateDatacenterFailureDomain(context.Background(), raw, testDatacenterFailureDomainStructural(t))
	if prob == nil {
		t.Fatal("two providerDatacenterRef values must be rejected")
	}
	assertProblemViolation(t, prob, apiproblem.CodeMalformedRequest, "/spec/providerDatacenterRef", apiproblem.ViolationCode(apiproblem.CodeMalformedRequest))
}

func TestValidateDatacenterFailureDomain_NegativeWrongKindReference(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		parentKind string
	}{
		{name: "provider-skipped-level", parentKind: resources.KindProvider},
		{name: "location-skipped-level", parentKind: resources.KindProviderLocation},
		{name: "failure-domain-self-kind", parentKind: resources.KindDatacenterFailureDomain},
		{name: "stack-skipped-level", parentKind: resources.KindInfrastructureStack},
	}

	structural := testDatacenterFailureDomainStructural(t)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			raw := mustDatacenterFailureDomainJSON(t, map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       resources.KindDatacenterFailureDomain,
				"metadata": map[string]any{
					"name": "fd-wrong-kind",
					"scopeRef": map[string]any{
						"apiVersion": resources.FabricAPIVersion,
						"kind":       string(apimeta.ScopeProvider),
						"name":       "sovereign-provider-a",
					},
				},
				"spec": map[string]any{
					"providerDatacenterRef": map[string]any{
						"apiVersion": resources.FabricAPIVersion,
						"kind":       tc.parentKind,
						"name":       "not-a-datacenter",
					},
				},
			})
			prob := ValidateDatacenterFailureDomain(context.Background(), raw, structural)
			assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/providerDatacenterRef/kind", apiproblem.ViolationCode(apiref.CodeKindNotAllowed))
		})
	}
}

func TestValidateDatacenterFailureDomain_NegativeWrongScopeKind(t *testing.T) {
	t.Parallel()

	raw := mustDatacenterFailureDomainJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindDatacenterFailureDomain,
		"metadata": map[string]any{
			"name": "fd-wrong-scope",
			"scopeRef": map[string]any{
				"apiVersion": "core.sovrunn.io/v1alpha1",
				"kind":       string(apimeta.ScopeOrganization),
				"name":       "acme",
			},
		},
		"spec": map[string]any{
			"providerDatacenterRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       resources.KindProviderDatacenter,
				"name":       "dc-bangalore-1",
			},
		},
	})
	prob := ValidateDatacenterFailureDomain(context.Background(), raw, testDatacenterFailureDomainStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/scopeRef/kind", apiproblem.ViolationCode(apiref.CodeScopeNotAllowed))
}

func TestValidateDatacenterFailureDomain_NegativeUserAuthoredStatus(t *testing.T) {
	t.Parallel()

	doc := validDatacenterFailureDomainDoc("fd-with-status", "dc-bangalore-1")
	doc["status"] = map[string]any{"observedGeneration": 1}
	raw := mustDatacenterFailureDomainJSON(t, doc)
	prob := ValidateDatacenterFailureDomain(context.Background(), raw, testDatacenterFailureDomainStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/status", violationValidationFail)
}

func TestValidateDatacenterFailureDomain_BoundaryReferenceCardinalityExactlyOne(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	structural := testDatacenterFailureDomainStructural(t)

	t.Run("exactly-one-accepted", func(t *testing.T) {
		t.Parallel()
		raw := mustDatacenterFailureDomainJSON(t, validDatacenterFailureDomainDoc("fd-card-one", "dc-bangalore-1"))
		if prob := ValidateDatacenterFailureDomain(ctx, raw, structural); prob != nil {
			t.Fatalf("exactly one providerDatacenterRef must pass: %#v", prob)
		}
	})

	t.Run("zero-rejected", func(t *testing.T) {
		t.Parallel()
		raw := mustDatacenterFailureDomainJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindDatacenterFailureDomain,
			"metadata": map[string]any{
				"name": "fd-card-zero",
				"scopeRef": map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       string(apimeta.ScopeProvider),
					"name":       "sovereign-provider-a",
				},
			},
			"spec": map[string]any{},
		})
		prob := ValidateDatacenterFailureDomain(ctx, raw, structural)
		assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/providerDatacenterRef", apiproblem.ViolationCode(apischema.CodeRequiredField))
	})

	t.Run("multiple-rejected", func(t *testing.T) {
		t.Parallel()
		raw := mustDatacenterFailureDomainJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindDatacenterFailureDomain,
			"metadata": map[string]any{
				"name": "fd-card-multi",
				"scopeRef": map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       string(apimeta.ScopeProvider),
					"name":       "sovereign-provider-a",
				},
			},
			"spec": map[string]any{
				"providerDatacenterRef": []any{
					map[string]any{
						"apiVersion": resources.FabricAPIVersion,
						"kind":       resources.KindProviderDatacenter,
						"name":       "dc-a",
					},
					map[string]any{
						"apiVersion": resources.FabricAPIVersion,
						"kind":       resources.KindProviderDatacenter,
						"name":       "dc-b",
					},
				},
			},
		})
		prob := ValidateDatacenterFailureDomain(ctx, raw, structural)
		if prob == nil {
			t.Fatal("multiple providerDatacenterRef values must be rejected")
		}
		assertProblemViolation(t, prob, apiproblem.CodeMalformedRequest, "/spec/providerDatacenterRef", apiproblem.ViolationCode(apiproblem.CodeMalformedRequest))
	})
}

func TestValidateDatacenterFailureDomain_BoundaryNoConnectivityFieldAccepted(t *testing.T) {
	t.Parallel()

	// Defense in depth for F14-REQ-19/F14-REQ-20: the Go type/schema carry no
	// connectivity field (primary evidence is Task 9 schema + Task 20 deny-list).
	// An authored connectivity field must be rejected as unknown; acceptance
	// must not invent connectivity semantics.
	specType := reflect.TypeOf(resources.DatacenterFailureDomainSpec{})
	for i := 0; i < specType.NumField(); i++ {
		name := specType.Field(i).Name
		lower := strings.ToLower(name)
		if strings.Contains(lower, "connect") ||
			strings.Contains(lower, "adjacency") ||
			strings.Contains(lower, "reachab") ||
			strings.Contains(lower, "peer") ||
			strings.Contains(lower, "latency") ||
			strings.Contains(lower, "bandwidth") {
			t.Fatalf("DatacenterFailureDomainSpec must not carry connectivity field %q", name)
		}
	}

	raw := mustDatacenterFailureDomainJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindDatacenterFailureDomain,
		"metadata": map[string]any{
			"name": "fd-connectivity",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{
			"providerDatacenterRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       resources.KindProviderDatacenter,
				"name":       "dc-bangalore-1",
			},
			"connectivity": map[string]any{
				"reachable": true,
			},
		},
	})
	prob := ValidateDatacenterFailureDomain(context.Background(), raw, testDatacenterFailureDomainStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeUnknownField, "/connectivity", apiproblem.ViolationUnknownField)

	// Valid document still carries no connectivity after round-trip.
	ok := mustDatacenterFailureDomainJSON(t, validDatacenterFailureDomainDoc("fd-no-connect", "dc-bangalore-1"))
	if prob := ValidateDatacenterFailureDomain(context.Background(), ok, testDatacenterFailureDomainStructural(t)); prob != nil {
		t.Fatalf("valid document must pass: %#v", prob)
	}
	var fd resources.DatacenterFailureDomain
	if err := json.Unmarshal(ok, &fd); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	blob, err := json.Marshal(fd.Spec)
	if err != nil {
		t.Fatalf("marshal spec: %v", err)
	}
	lower := strings.ToLower(string(blob))
	for _, banned := range []string{"connectivity", "reachable", "adjacency", "latency", "bandwidth"} {
		if strings.Contains(lower, banned) {
			t.Fatalf("accepted spec must not encode connectivity token %q: %s", banned, blob)
		}
	}
}

func validDatacenterFailureDomainDoc(name, datacenterName string) map[string]any {
	return map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindDatacenterFailureDomain,
		"metadata": map[string]any{
			"name": name,
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{
			"providerDatacenterRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       resources.KindProviderDatacenter,
				"name":       datacenterName,
			},
		},
	}
}

func mustDatacenterFailureDomainJSON(t *testing.T, doc map[string]any) []byte {
	t.Helper()
	return mustProviderLocationJSON(t, doc)
}

// testDatacenterFailureDomainStructural builds the existing apiconform
// structural validator for tests only. Production
// ValidateDatacenterFailureDomain depends solely on the injected
// apivalid.StructuralValidator interface.
func testDatacenterFailureDomainStructural(t *testing.T) apivalid.StructuralValidator {
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
