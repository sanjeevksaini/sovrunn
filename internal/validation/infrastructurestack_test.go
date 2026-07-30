package validation

import (
	"context"
	"path/filepath"
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

func TestValidateInfrastructureStack_PositiveWithAndWithoutTechnology(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	structural := testInfrastructureStackStructural(t)

	withoutTech := mustInfrastructureStackJSON(t, validInfrastructureStackDoc("stack-no-tech", "fd-bangalore-az1", ""))
	if prob := ValidateInfrastructureStack(ctx, withoutTech, structural); prob != nil {
		t.Fatalf("without technology: unexpected problem: %#v", prob)
	}

	withTech := mustInfrastructureStackJSON(t, validInfrastructureStackDoc("stack-with-tech", "fd-bangalore-az1", "Apache CloudStack"))
	if prob := ValidateInfrastructureStack(ctx, withTech, structural); prob != nil {
		t.Fatalf("with valid technology: unexpected problem: %#v", prob)
	}
}

func TestValidateInfrastructureStack_PositiveIdenticalTechnologyIndependentIdentity(t *testing.T) {
	t.Parallel()

	// Two stacks may share the same descriptive technology; identity is uid /
	// scoped name, not technology (F14-REQ-15, F14-REQ-17).
	ctx := context.Background()
	structural := testInfrastructureStackStructural(t)
	tech := "Red Hat OpenShift"

	a := mustInfrastructureStackJSON(t, validInfrastructureStackDoc("stack-a", "fd-bangalore-az1", tech))
	b := mustInfrastructureStackJSON(t, validInfrastructureStackDoc("stack-b", "fd-bangalore-az2", tech))

	if prob := ValidateInfrastructureStack(ctx, a, structural); prob != nil {
		t.Fatalf("stack-a with shared technology must pass: %#v", prob)
	}
	if prob := ValidateInfrastructureStack(ctx, b, structural); prob != nil {
		t.Fatalf("stack-b with shared technology must pass independently: %#v", prob)
	}
}

func TestValidateInfrastructureStack_NegativeMissingReference(t *testing.T) {
	t.Parallel()

	raw := mustInfrastructureStackJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindInfrastructureStack,
		"metadata": map[string]any{
			"name": "stack-missing-ref",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{},
	})
	prob := ValidateInfrastructureStack(context.Background(), raw, testInfrastructureStackStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/datacenterFailureDomainRef", apiproblem.ViolationCode(apischema.CodeRequiredField))
}

func TestValidateInfrastructureStack_NegativeTwoReferences(t *testing.T) {
	t.Parallel()

	// Cardinality is exactly one singular object. An array is rejected at
	// decode (type mismatch) before semantic checks — no spanning (F14-REQ-16).
	raw := mustInfrastructureStackJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindInfrastructureStack,
		"metadata": map[string]any{
			"name": "stack-two-refs",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{
			"datacenterFailureDomainRef": []any{
				map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       resources.KindDatacenterFailureDomain,
					"name":       "fd-a",
				},
				map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       resources.KindDatacenterFailureDomain,
					"name":       "fd-b",
				},
			},
		},
	})
	prob := ValidateInfrastructureStack(context.Background(), raw, testInfrastructureStackStructural(t))
	if prob == nil {
		t.Fatal("two datacenterFailureDomainRef values must be rejected")
	}
	assertProblemViolation(t, prob, apiproblem.CodeMalformedRequest, "/spec/datacenterFailureDomainRef", apiproblem.ViolationCode(apiproblem.CodeMalformedRequest))
}

func TestValidateInfrastructureStack_NegativeWrongKindReference(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		parentKind string
	}{
		{name: "provider-skipped-level", parentKind: resources.KindProvider},
		{name: "location-skipped-level", parentKind: resources.KindProviderLocation},
		{name: "datacenter-skipped-level", parentKind: resources.KindProviderDatacenter},
		{name: "stack-self-kind", parentKind: resources.KindInfrastructureStack},
	}

	structural := testInfrastructureStackStructural(t)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			raw := mustInfrastructureStackJSON(t, map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       resources.KindInfrastructureStack,
				"metadata": map[string]any{
					"name": "stack-wrong-kind",
					"scopeRef": map[string]any{
						"apiVersion": resources.FabricAPIVersion,
						"kind":       string(apimeta.ScopeProvider),
						"name":       "sovereign-provider-a",
					},
				},
				"spec": map[string]any{
					"datacenterFailureDomainRef": map[string]any{
						"apiVersion": resources.FabricAPIVersion,
						"kind":       tc.parentKind,
						"name":       "not-a-failure-domain",
					},
				},
			})
			prob := ValidateInfrastructureStack(context.Background(), raw, structural)
			assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/datacenterFailureDomainRef/kind", apiproblem.ViolationCode(apiref.CodeKindNotAllowed))
		})
	}
}

func TestValidateInfrastructureStack_NegativeWrongScopeKind(t *testing.T) {
	t.Parallel()

	raw := mustInfrastructureStackJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindInfrastructureStack,
		"metadata": map[string]any{
			"name": "stack-wrong-scope",
			"scopeRef": map[string]any{
				"apiVersion": "core.sovrunn.io/v1alpha1",
				"kind":       string(apimeta.ScopeOrganization),
				"name":       "acme",
			},
		},
		"spec": map[string]any{
			"datacenterFailureDomainRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       resources.KindDatacenterFailureDomain,
				"name":       "fd-bangalore-az1",
			},
		},
	})
	prob := ValidateInfrastructureStack(context.Background(), raw, testInfrastructureStackStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/scopeRef/kind", apiproblem.ViolationCode(apiref.CodeScopeNotAllowed))
}

func TestValidateInfrastructureStack_NegativeUserAuthoredStatus(t *testing.T) {
	t.Parallel()

	doc := validInfrastructureStackDoc("stack-with-status", "fd-bangalore-az1", "Apache CloudStack")
	doc["status"] = map[string]any{"observedGeneration": 1}
	raw := mustInfrastructureStackJSON(t, doc)
	prob := ValidateInfrastructureStack(context.Background(), raw, testInfrastructureStackStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/status", violationValidationFail)
}

func TestValidateInfrastructureStack_NegativeTechnology(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		tech      string
		wantField string
		wantCode  apiproblem.ViolationCode
	}{
		{
			name:      "control-character",
			tech:      "Apache\tCloudStack",
			wantField: "/spec/technology",
			wantCode:  apiproblem.ViolationCode(apischema.CodePatternMismatch),
		},
		{
			name:      "non-ascii",
			tech:      "Apache Café",
			wantField: "/spec/technology",
			wantCode:  apiproblem.ViolationCode(apischema.CodePatternMismatch),
		},
		{
			name:      "over-100-chars",
			tech:      strings.Repeat("a", 101),
			wantField: "/spec/technology",
			wantCode:  apiproblem.ViolationOutOfRange,
		},
		{
			name:      "leading-trailing-space-after-normalization",
			tech:      " Apache CloudStack ",
			wantField: "/spec/technology",
			wantCode:  apiproblem.ViolationCode(apischema.CodePatternMismatch),
		},
	}

	structural := testInfrastructureStackStructural(t)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			raw := mustInfrastructureStackJSON(t, validInfrastructureStackDoc("stack-bad-tech", "fd-bangalore-az1", tc.tech))
			prob := ValidateInfrastructureStack(context.Background(), raw, structural)
			assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, tc.wantField, tc.wantCode)
		})
	}
}

func TestValidateInfrastructureStack_NegativeTechnologyAllSpacesSemantic(t *testing.T) {
	t.Parallel()

	// After §6.3 trim, all-spaces becomes empty and must be rejected. Use
	// pass-through structural so the semantic normalizer owns the finding.
	raw := mustInfrastructureStackJSON(t, validInfrastructureStackDoc("stack-spaces-tech", "fd-bangalore-az1", "   "))
	prob := ValidateInfrastructureStack(context.Background(), raw, passThroughStructural{})
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/technology", apiproblem.ViolationOutOfRange)
}

func TestValidateInfrastructureStack_BoundaryTechnologyLengthAndSpaces(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	structural := testInfrastructureStackStructural(t)

	t.Run("technology-exactly-100-accepted", func(t *testing.T) {
		t.Parallel()
		tech100 := "A" + strings.Repeat("b", 99)
		if len(tech100) != 100 {
			t.Fatalf("fixture length = %d, want 100", len(tech100))
		}
		raw := mustInfrastructureStackJSON(t, validInfrastructureStackDoc("stack-tech-100", "fd-bangalore-az1", tech100))
		if prob := ValidateInfrastructureStack(ctx, raw, structural); prob != nil {
			t.Fatalf("exactly-100 technology must pass: %#v", prob)
		}
	})

	t.Run("technology-101-rejected", func(t *testing.T) {
		t.Parallel()
		tech101 := "A" + strings.Repeat("b", 100)
		if len(tech101) != 101 {
			t.Fatalf("fixture length = %d, want 101", len(tech101))
		}
		raw := mustInfrastructureStackJSON(t, validInfrastructureStackDoc("stack-tech-101", "fd-bangalore-az1", tech101))
		prob := ValidateInfrastructureStack(ctx, raw, structural)
		assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/technology", apiproblem.ViolationOutOfRange)
	})

	t.Run("single-internal-space-preserved", func(t *testing.T) {
		t.Parallel()
		normalized, ok := normalizeTechnology("Apache CloudStack")
		if !ok {
			t.Fatal("allowed characters must normalize")
		}
		if normalized != "Apache CloudStack" {
			t.Fatalf("single internal space must be preserved, got %q", normalized)
		}
		raw := mustInfrastructureStackJSON(t, validInfrastructureStackDoc("stack-tech-space", "fd-bangalore-az1", "Apache CloudStack"))
		if prob := ValidateInfrastructureStack(ctx, raw, structural); prob != nil {
			t.Fatalf("single-space technology must pass: %#v", prob)
		}
	})

	t.Run("double-internal-space-collapsed", func(t *testing.T) {
		t.Parallel()
		normalized, ok := normalizeTechnology("Apache  CloudStack")
		if !ok {
			t.Fatal("allowed characters must normalize")
		}
		if normalized != "Apache CloudStack" {
			t.Fatalf("double internal space must collapse, got %q", normalized)
		}
		// Pattern allows double spaces structurally; semantic normalization
		// collapses them and the normalized form remains valid.
		raw := mustInfrastructureStackJSON(t, validInfrastructureStackDoc("stack-tech-dbl", "fd-bangalore-az1", "Apache  CloudStack"))
		if prob := ValidateInfrastructureStack(ctx, raw, structural); prob != nil {
			t.Fatalf("double-space technology must pass after normalization: %#v", prob)
		}
	})
}

func TestValidateInfrastructureStack_NoSupersededStackKindName(t *testing.T) {
	t.Parallel()

	// F14-REQ-03 / F14-AD-021: validator must not introduce the superseded name.
	const superseded = "IaaSStack"
	raw := mustInfrastructureStackJSON(t, validInfrastructureStackDoc("stack-active-name", "fd-bangalore-az1", "AWS IaaS"))
	if strings.Contains(string(raw), superseded) {
		t.Fatalf("fixture must not carry superseded kind %q", superseded)
	}
	if prob := ValidateInfrastructureStack(context.Background(), raw, testInfrastructureStackStructural(t)); prob != nil {
		t.Fatalf("active InfrastructureStack name must validate: %#v", prob)
	}
	if resources.KindInfrastructureStack == superseded {
		t.Fatalf("KindInfrastructureStack must not equal superseded name %q", superseded)
	}
}

func validInfrastructureStackDoc(name, failureDomainName, technology string) map[string]any {
	spec := map[string]any{
		"datacenterFailureDomainRef": map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindDatacenterFailureDomain,
			"name":       failureDomainName,
		},
	}
	if technology != "" {
		spec["technology"] = technology
	}
	return map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindInfrastructureStack,
		"metadata": map[string]any{
			"name": name,
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": spec,
	}
}

func mustInfrastructureStackJSON(t *testing.T, doc map[string]any) []byte {
	t.Helper()
	return mustProviderLocationJSON(t, doc)
}

// testInfrastructureStackStructural builds the existing apiconform structural
// validator for tests only. Production ValidateInfrastructureStack depends
// solely on the injected apivalid.StructuralValidator interface.
func testInfrastructureStackStructural(t *testing.T) apivalid.StructuralValidator {
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
