package validation

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateProvider_PositivePlatformNilScope(t *testing.T) {
	t.Parallel()

	raw := mustProviderJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProvider,
		"metadata": map[string]any{
			"name": "sovereign-provider-a",
		},
		"spec": map[string]any{},
	})
	if prob := ValidateProvider(context.Background(), raw, testProviderStructural(t)); prob != nil {
		t.Fatalf("nil/absent Platform scope must pass: %#v", prob)
	}
}

func TestValidateProvider_PositiveExplicitPlatformScope(t *testing.T) {
	t.Parallel()

	raw := mustProviderJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProvider,
		"metadata": map[string]any{
			"name": "sovereign-provider-b",
			"scopeRef": map[string]any{
				"apiVersion": "core.sovrunn.io/v1alpha1",
				"kind":       string(apimeta.ScopePlatform),
				"name":       "platform",
			},
		},
		"spec": map[string]any{},
	})
	if prob := ValidateProvider(context.Background(), raw, testProviderStructural(t)); prob != nil {
		t.Fatalf("explicit complete Platform scope must pass: %#v", prob)
	}
}

func TestValidateProvider_PositiveOrganizationScope(t *testing.T) {
	t.Parallel()

	raw := mustProviderJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProvider,
		"metadata": map[string]any{
			"name": "sovereign-provider-c",
			"scopeRef": map[string]any{
				"apiVersion": "core.sovrunn.io/v1alpha1",
				"kind":       string(apimeta.ScopeOrganization),
				"name":       "acme",
			},
		},
		"spec": map[string]any{},
	})
	if prob := ValidateProvider(context.Background(), raw, testProviderStructural(t)); prob != nil {
		t.Fatalf("Organization scope must pass: %#v", prob)
	}
}

func TestValidateProvider_NegativeWrongScopeKind(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		kind string
	}{
		{name: "tenant", kind: string(apimeta.ScopeTenant)},
		{name: "project", kind: string(apimeta.ScopeProject)},
		{name: "provider", kind: string(apimeta.ScopeProvider)},
		{name: "organization-unit", kind: string(apimeta.ScopeOrganizationUnit)},
	}

	structural := testProviderStructural(t)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			raw := mustProviderJSON(t, map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       resources.KindProvider,
				"metadata": map[string]any{
					"name": "prov-wrong-scope",
					"scopeRef": map[string]any{
						"apiVersion": "core.sovrunn.io/v1alpha1",
						"kind":       tc.kind,
						"name":       "acme",
					},
				},
				"spec": map[string]any{},
			})
			prob := ValidateProvider(context.Background(), raw, structural)
			assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/scopeRef/kind", apiproblem.ViolationCode(apiref.CodeScopeNotAllowed))
		})
	}
}

func TestValidateProvider_NegativeOwnerRefAsScope(t *testing.T) {
	t.Parallel()

	// ownerRef must not replace governance scope. ObjectMeta rejects it as an
	// unknown field; status/scope remain solely via scopeRef (F14-REQ-05).
	raw := mustProviderJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProvider,
		"metadata": map[string]any{
			"name": "prov-ownerref",
			"ownerRef": map[string]any{
				"apiVersion": "core.sovrunn.io/v1alpha1",
				"kind":       string(apimeta.ScopeOrganization),
				"name":       "acme",
			},
		},
		"spec": map[string]any{},
	})
	prob := ValidateProvider(context.Background(), raw, testProviderStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeUnknownField, "/ownerRef", apiproblem.ViolationUnknownField)
}

func TestValidateProvider_NegativeOwnerField(t *testing.T) {
	t.Parallel()

	raw := mustProviderJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProvider,
		"metadata": map[string]any{
			"name":  "prov-owner-field",
			"owner": "acme",
		},
		"spec": map[string]any{},
	})
	prob := ValidateProvider(context.Background(), raw, testProviderStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeUnknownField, "/owner", apiproblem.ViolationUnknownField)
}

func TestValidateProvider_NegativeDomainSpecField(t *testing.T) {
	t.Parallel()

	raw := mustProviderJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProvider,
		"metadata": map[string]any{
			"name": "prov-domain-spec",
		},
		"spec": map[string]any{
			"region": "apac",
		},
	})
	prob := ValidateProvider(context.Background(), raw, testProviderStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeUnknownField, "/region", apiproblem.ViolationUnknownField)
}

func TestValidateProvider_NegativeUserAuthoredStatus(t *testing.T) {
	t.Parallel()

	raw := mustProviderJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProvider,
		"metadata": map[string]any{
			"name": "prov-with-status",
		},
		"spec":   map[string]any{},
		"status": map[string]any{"observedGeneration": 1},
	})
	prob := ValidateProvider(context.Background(), raw, testProviderStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/status", violationValidationFail)
}

func TestValidateProvider_NegativeNilStructuralFailsClosed(t *testing.T) {
	t.Parallel()

	raw := mustProviderJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProvider,
		"metadata": map[string]any{
			"name": "prov-nil-structural",
		},
		"spec": map[string]any{},
	})
	prob := ValidateProvider(context.Background(), raw, nil)
	if prob == nil {
		t.Fatal("nil structural validator must fail closed")
	}
	if prob.Code != apiproblem.CodeInternalError {
		t.Fatalf("problem code = %q, want %q", prob.Code, apiproblem.CodeInternalError)
	}
}

func TestValidateProvider_NegativeUnavailableStructuralFailsClosed(t *testing.T) {
	t.Parallel()

	raw := mustProviderJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProvider,
		"metadata": map[string]any{
			"name": "prov-unavailable-structural",
		},
		"spec": map[string]any{},
	})
	prob := ValidateProvider(context.Background(), raw, unavailableStructural{})
	if prob == nil {
		t.Fatal("unavailable structural validator must fail closed")
	}
	if prob.Code != apiproblem.CodeInternalError {
		t.Fatalf("problem code = %q, want %q", prob.Code, apiproblem.CodeInternalError)
	}
}

func TestValidateProvider_NegativeIncompleteOrganizationScope(t *testing.T) {
	t.Parallel()

	raw := mustProviderJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProvider,
		"metadata": map[string]any{
			"name": "prov-incomplete-org",
			"scopeRef": map[string]any{
				"kind": string(apimeta.ScopeOrganization),
				"name": "acme",
			},
		},
		"spec": map[string]any{},
	})
	// Pass-through structural isolates the apiref.Constraint.ValidateRef path
	// that requires scopeRef.apiVersion.
	prob := ValidateProvider(context.Background(), raw, passThroughStructural{})
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/scopeRef/apiVersion", apiproblem.ViolationCode(apiref.CodeMissingAPIVersion))
}

func TestValidateProvider_BoundaryNameLabelAnnotationLimits(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	lim := apivalid.DefaultLimits()
	structural := testProviderStructural(t)

	t.Run("name-at-63", func(t *testing.T) {
		t.Parallel()
		raw := mustProviderJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindProvider,
			"metadata": map[string]any{
				"name": strings.Repeat("a", 63),
			},
			"spec": map[string]any{},
		})
		if prob := ValidateProvider(ctx, raw, structural); prob != nil {
			t.Fatalf("63-char name must pass: %#v", prob)
		}
	})

	t.Run("name-over-63", func(t *testing.T) {
		t.Parallel()
		raw := mustProviderJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindProvider,
			"metadata": map[string]any{
				"name": strings.Repeat("a", 64),
			},
			"spec": map[string]any{},
		})
		prob := ValidateProvider(ctx, raw, structural)
		assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/name", apiproblem.ViolationOutOfRange)
	})

	t.Run("labels-at-MaxLabels", func(t *testing.T) {
		t.Parallel()
		labels := make(map[string]any, lim.MaxLabels)
		for i := 0; i < lim.MaxLabels; i++ {
			labels[fmt.Sprintf("k%02d", i)] = "v"
		}
		raw := mustProviderJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindProvider,
			"metadata": map[string]any{
				"name":   "prov-labels-ok",
				"labels": labels,
			},
			"spec": map[string]any{},
		})
		if prob := ValidateProvider(ctx, raw, structural); prob != nil {
			t.Fatalf("MaxLabels edge must pass: %#v", prob)
		}
	})

	t.Run("labels-over-MaxLabels", func(t *testing.T) {
		t.Parallel()
		labels := make(map[string]any, lim.MaxLabels+1)
		for i := 0; i < lim.MaxLabels+1; i++ {
			labels[fmt.Sprintf("k%02d", i)] = "v"
		}
		raw := mustProviderJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindProvider,
			"metadata": map[string]any{
				"name":   "prov-labels-over",
				"labels": labels,
			},
			"spec": map[string]any{},
		})
		prob := ValidateProvider(ctx, raw, structural)
		assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/labels", apiproblem.ViolationOutOfRange)
	})

	t.Run("annotations-over-MaxAnnotationsBytes", func(t *testing.T) {
		t.Parallel()
		oversized := strings.Repeat("x", lim.MaxAnnotationsBytes+1)
		raw := mustProviderJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindProvider,
			"metadata": map[string]any{
				"name": "prov-ann-over",
				"annotations": map[string]any{
					"note": oversized,
				},
			},
			"spec": map[string]any{},
		})
		prob := ValidateProvider(ctx, raw, structural)
		assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/annotations", apiproblem.ViolationOutOfRange)
	})
}

func TestValidateProvider_ErrorsCarryNoSecrets(t *testing.T) {
	t.Parallel()

	raw := mustProviderJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProvider,
		"metadata": map[string]any{
			"name": "prov-secret-check",
			"scopeRef": map[string]any{
				"apiVersion": "core.sovrunn.io/v1alpha1",
				"kind":       string(apimeta.ScopeTenant),
				"name":       "acme",
			},
			"annotations": map[string]any{
				"token": "super-secret-token-value",
			},
		},
		"spec": map[string]any{},
	})
	prob := ValidateProvider(context.Background(), raw, testProviderStructural(t))
	if prob == nil {
		t.Fatal("expected validation failure")
	}
	blob, err := json.Marshal(prob)
	if err != nil {
		t.Fatalf("marshal problem: %v", err)
	}
	if strings.Contains(string(blob), "super-secret-token-value") {
		t.Fatalf("problem must not echo secrets: %s", blob)
	}
}

func mustProviderJSON(t *testing.T, doc map[string]any) []byte {
	t.Helper()
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return raw
}

// testProviderStructural builds the existing apiconform structural validator
// for tests only. Production ValidateProvider depends solely on the injected
// apivalid.StructuralValidator interface.
func testProviderStructural(t *testing.T) apivalid.StructuralValidator {
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
