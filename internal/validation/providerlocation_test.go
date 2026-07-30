package validation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

func TestValidateProviderLocation_PositiveWithAndWithoutGeo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	structural := testProviderLocationStructural(t)

	withoutGeo := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-no-geo",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{},
	})
	if prob := ValidateProviderLocation(ctx, withoutGeo, structural); prob != nil {
		t.Fatalf("without geo: unexpected problem: %#v", prob)
	}

	withGeo := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-with-geo",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{
			"geo": map[string]any{
				"countryCode":     "IN",
				"subdivisionCode": "IN-KA",
			},
		},
	})
	if prob := ValidateProviderLocation(ctx, withGeo, structural); prob != nil {
		t.Fatalf("with geo: unexpected problem: %#v", prob)
	}
}

func TestValidateProviderLocation_NegativeWrongScopeKind(t *testing.T) {
	t.Parallel()

	raw := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-wrong-scope",
			"scopeRef": map[string]any{
				"apiVersion": "core.sovrunn.io/v1alpha1",
				"kind":       string(apimeta.ScopeOrganization),
				"name":       "acme",
			},
		},
		"spec": map[string]any{},
	})
	prob := ValidateProviderLocation(context.Background(), raw, testProviderLocationStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/scopeRef/kind", apiproblem.ViolationCode(apiref.CodeScopeNotAllowed))
}

func TestValidateProviderLocation_NegativeUserAuthoredStatus(t *testing.T) {
	t.Parallel()

	raw := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-with-status",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec":   map[string]any{},
		"status": map[string]any{"observedGeneration": 1},
	})
	prob := ValidateProviderLocation(context.Background(), raw, testProviderLocationStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/status", violationValidationFail)
}

func TestValidateProviderLocation_NegativeUnknownField(t *testing.T) {
	t.Parallel()

	raw := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-unknown",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{
			"region": "apac",
		},
	})
	prob := ValidateProviderLocation(context.Background(), raw, testProviderLocationStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeUnknownField, "/region", apiproblem.ViolationUnknownField)
}

func TestValidateProviderLocation_NegativeMissingSpec(t *testing.T) {
	t.Parallel()

	raw := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-missing-spec",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
	})
	prob := ValidateProviderLocation(context.Background(), raw, testProviderLocationStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec", apiproblem.ViolationCode(apischema.CodeRequiredField))
}

func TestValidateProviderLocation_NegativeMissingScopeRefAPIVersion(t *testing.T) {
	t.Parallel()

	raw := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-missing-scope-apiversion",
			"scopeRef": map[string]any{
				"kind": string(apimeta.ScopeProvider),
				"name": "sovereign-provider-a",
			},
		},
		"spec": map[string]any{},
	})
	// Pass-through structural isolates the apiref.Constraint.ValidateRef path
	// that requires scopeRef.apiVersion (canonical schema also requires it).
	prob := ValidateProviderLocation(context.Background(), raw, passThroughStructural{})
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/scopeRef/apiVersion", apiproblem.ViolationCode(apiref.CodeMissingAPIVersion))
}

func TestValidateProviderLocation_NegativeMissingScopeRefName(t *testing.T) {
	t.Parallel()

	raw := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-missing-scope-name",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
			},
		},
		"spec": map[string]any{},
	})
	prob := ValidateProviderLocation(context.Background(), raw, passThroughStructural{})
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/scopeRef/name", apiproblem.ViolationCode(apiref.CodeMissingName))
}

func TestValidateProviderLocation_NegativeNilStructuralFailsClosed(t *testing.T) {
	t.Parallel()

	raw := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-nil-structural",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{},
	})
	prob := ValidateProviderLocation(context.Background(), raw, nil)
	if prob == nil {
		t.Fatal("nil structural validator must fail closed")
	}
	if prob.Code != apiproblem.CodeInternalError {
		t.Fatalf("problem code = %q, want %q", prob.Code, apiproblem.CodeInternalError)
	}
}

func TestValidateProviderLocation_NegativeUnavailableStructuralFailsClosed(t *testing.T) {
	t.Parallel()

	raw := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-unavailable-structural",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{},
	})
	prob := ValidateProviderLocation(context.Background(), raw, unavailableStructural{})
	if prob == nil {
		t.Fatal("unavailable structural validator must fail closed")
	}
	if prob.Code != apiproblem.CodeInternalError {
		t.Fatalf("problem code = %q, want %q", prob.Code, apiproblem.CodeInternalError)
	}
}

func TestValidateProviderLocation_NegativeMalformedGeo(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		geo       map[string]any
		wantField string
		wantCode  apiproblem.ViolationCode
	}{
		{
			name:      "lowercase-country",
			geo:       map[string]any{"countryCode": "in"},
			wantField: "/spec/geo/countryCode",
			wantCode:  apiproblem.ViolationCode(apischema.CodePatternMismatch),
		},
		{
			name:      "country-too-long",
			geo:       map[string]any{"countryCode": "IND"},
			wantField: "/spec/geo/countryCode",
			wantCode:  apiproblem.ViolationOutOfRange,
		},
		{
			name:      "country-too-short",
			geo:       map[string]any{"countryCode": "I"},
			wantField: "/spec/geo/countryCode",
			wantCode:  apiproblem.ViolationOutOfRange,
		},
		{
			name: "subdivision-malformed",
			geo: map[string]any{
				"countryCode":     "IN",
				"subdivisionCode": "IN_KA",
			},
			wantField: "/spec/geo/subdivisionCode",
			wantCode:  apiproblem.ViolationCode(apischema.CodePatternMismatch),
		},
		{
			name: "subdivision-too-long",
			geo: map[string]any{
				"countryCode":     "IN",
				"subdivisionCode": "IN-KARN",
			},
			wantField: "/spec/geo/subdivisionCode",
			wantCode:  apiproblem.ViolationOutOfRange,
		},
	}

	structural := testProviderLocationStructural(t)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			raw := mustProviderLocationJSON(t, map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       resources.KindProviderLocation,
				"metadata": map[string]any{
					"name": "loc-bad-geo",
					"scopeRef": map[string]any{
						"apiVersion": resources.FabricAPIVersion,
						"kind":       string(apimeta.ScopeProvider),
						"name":       "sovereign-provider-a",
					},
				},
				"spec": map[string]any{"geo": tc.geo},
			})
			prob := ValidateProviderLocation(context.Background(), raw, structural)
			assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, tc.wantField, tc.wantCode)
		})
	}
}

func TestValidateProviderLocation_NegativePrefixInconsistentGeo(t *testing.T) {
	t.Parallel()

	raw := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-prefix-mismatch",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{
			"geo": map[string]any{
				"countryCode":     "IN",
				"subdivisionCode": "US-CA",
			},
		},
	})
	prob := ValidateProviderLocation(context.Background(), raw, testProviderLocationStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/geo/subdivisionCode", violationValidationFail)
}

func TestValidateProviderLocation_BoundaryUnassignedGeoAcceptedWithoutInference(t *testing.T) {
	t.Parallel()

	// Syntactically valid but unassigned user-assigned codes (ISO 3166
	// AA/QM–QZ/XA–XZ/ZZ style). Acceptance is format-only: the validator
	// performs no membership lookup and must not invent residency,
	// compliance, authoritative-assignment, or placement meaning.
	raw := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-unassigned-geo",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{
			"geo": map[string]any{
				"countryCode":     "XZ",
				"subdivisionCode": "XZ-A1",
			},
		},
	})
	prob := ValidateProviderLocation(context.Background(), raw, testProviderLocationStructural(t))
	if prob != nil {
		t.Fatalf("unassigned syntactically valid geo must be accepted, got %#v", prob)
	}

	// Re-decode the same document: validation must not mutate or inject
	// status / residency / compliance / placement fields.
	var loc resources.ProviderLocation
	if err := json.Unmarshal(raw, &loc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if loc.Status.ObservedGeneration != 0 || len(loc.Status.Conditions) != 0 {
		t.Fatalf("acceptance must not invent status facts: %+v", loc.Status)
	}
	if loc.Spec.Geo == nil || loc.Spec.Geo.CountryCode != "XZ" || loc.Spec.Geo.SubdivisionCode != "XZ-A1" {
		t.Fatalf("geo must remain declared metadata only: %+v", loc.Spec.Geo)
	}
	if _, ok := loc.Metadata.Annotations["residency"]; ok {
		t.Fatal("must not infer residency annotation")
	}
	if _, ok := loc.Metadata.Annotations["compliance"]; ok {
		t.Fatal("must not infer compliance annotation")
	}
	if _, ok := loc.Metadata.Labels["placement"]; ok {
		t.Fatal("must not infer placement label")
	}
}

func TestValidateProviderLocation_BoundaryGeoLengthsAndPrefix(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	structural := testProviderLocationStructural(t)

	t.Run("countryCode-exactly-2", func(t *testing.T) {
		t.Parallel()
		ok := mustProviderLocationJSON(t, baseLocDoc(map[string]any{
			"geo": map[string]any{"countryCode": "US"},
		}))
		if prob := ValidateProviderLocation(ctx, ok, structural); prob != nil {
			t.Fatalf("exactly-2 countryCode must pass: %#v", prob)
		}
	})

	t.Run("subdivisionCode-at-6-with-matching-prefix", func(t *testing.T) {
		t.Parallel()
		ok := mustProviderLocationJSON(t, baseLocDoc(map[string]any{
			"geo": map[string]any{
				"countryCode":     "IN",
				"subdivisionCode": "IN-KA1", // 6 chars, matching prefix
			},
		}))
		if prob := ValidateProviderLocation(ctx, ok, structural); prob != nil {
			t.Fatalf("6-char matching subdivision must pass: %#v", prob)
		}
	})

	t.Run("subdivisionCode-over-6-rejected", func(t *testing.T) {
		t.Parallel()
		bad := mustProviderLocationJSON(t, baseLocDoc(map[string]any{
			"geo": map[string]any{
				"countryCode":     "IN",
				"subdivisionCode": "IN-KA12", // 7 chars
			},
		}))
		prob := ValidateProviderLocation(ctx, bad, structural)
		assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/geo/subdivisionCode", apiproblem.ViolationOutOfRange)
	})
}

func TestValidateProviderLocation_BoundaryNameLabelAnnotationLimits(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	lim := apivalid.DefaultLimits()
	structural := testProviderLocationStructural(t)

	t.Run("name-at-63", func(t *testing.T) {
		t.Parallel()
		name63 := strings.Repeat("a", 63)
		raw := mustProviderLocationJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindProviderLocation,
			"metadata": map[string]any{
				"name": name63,
				"scopeRef": map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       string(apimeta.ScopeProvider),
					"name":       "sovereign-provider-a",
				},
			},
			"spec": map[string]any{},
		})
		if prob := ValidateProviderLocation(ctx, raw, structural); prob != nil {
			t.Fatalf("63-char name must pass: %#v", prob)
		}
	})

	t.Run("name-over-63", func(t *testing.T) {
		t.Parallel()
		raw := mustProviderLocationJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindProviderLocation,
			"metadata": map[string]any{
				"name": strings.Repeat("a", 64),
				"scopeRef": map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       string(apimeta.ScopeProvider),
					"name":       "sovereign-provider-a",
				},
			},
			"spec": map[string]any{},
		})
		prob := ValidateProviderLocation(ctx, raw, structural)
		assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/name", apiproblem.ViolationOutOfRange)
	})

	t.Run("labels-at-MaxLabels", func(t *testing.T) {
		t.Parallel()
		labels := make(map[string]any, lim.MaxLabels)
		for i := 0; i < lim.MaxLabels; i++ {
			labels[fmt.Sprintf("k%02d", i)] = "v"
		}
		raw := mustProviderLocationJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindProviderLocation,
			"metadata": map[string]any{
				"name":   "loc-labels-ok",
				"labels": labels,
				"scopeRef": map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       string(apimeta.ScopeProvider),
					"name":       "sovereign-provider-a",
				},
			},
			"spec": map[string]any{},
		})
		if prob := ValidateProviderLocation(ctx, raw, structural); prob != nil {
			t.Fatalf("MaxLabels edge must pass: %#v", prob)
		}
	})

	t.Run("labels-over-MaxLabels", func(t *testing.T) {
		t.Parallel()
		labels := make(map[string]any, lim.MaxLabels+1)
		for i := 0; i < lim.MaxLabels+1; i++ {
			labels[fmt.Sprintf("k%02d", i)] = "v"
		}
		raw := mustProviderLocationJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindProviderLocation,
			"metadata": map[string]any{
				"name":   "loc-labels-over",
				"labels": labels,
				"scopeRef": map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       string(apimeta.ScopeProvider),
					"name":       "sovereign-provider-a",
				},
			},
			"spec": map[string]any{},
		})
		prob := ValidateProviderLocation(ctx, raw, structural)
		assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/labels", apiproblem.ViolationOutOfRange)
	})

	t.Run("annotations-over-MaxAnnotationsBytes", func(t *testing.T) {
		t.Parallel()
		// One annotation whose key+value exceed MaxAnnotationsBytes.
		oversized := strings.Repeat("x", lim.MaxAnnotationsBytes+1)
		raw := mustProviderLocationJSON(t, map[string]any{
			"apiVersion": resources.FabricAPIVersion,
			"kind":       resources.KindProviderLocation,
			"metadata": map[string]any{
				"name": "loc-ann-over",
				"annotations": map[string]any{
					"note": oversized,
				},
				"scopeRef": map[string]any{
					"apiVersion": resources.FabricAPIVersion,
					"kind":       string(apimeta.ScopeProvider),
					"name":       "sovereign-provider-a",
				},
			},
			"spec": map[string]any{},
		})
		prob := ValidateProviderLocation(ctx, raw, structural)
		assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/annotations", apiproblem.ViolationOutOfRange)
	})
}

func TestValidateProviderLocation_NegativeSystemOwnedMetadata(t *testing.T) {
	t.Parallel()

	raw := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name":            "loc-sys",
			"resourceVersion": "3",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": map[string]any{},
	})
	prob := ValidateProviderLocation(context.Background(), raw, testProviderLocationStructural(t))
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/resourceVersion", violationValidationFail)
}

func TestValidateProviderLocation_ErrorsCarryNoSecrets(t *testing.T) {
	t.Parallel()

	raw := mustProviderLocationJSON(t, map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-secret-check",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeOrganization),
				"name":       "acme",
			},
			"annotations": map[string]any{
				"token": "super-secret-token-value",
			},
		},
		"spec": map[string]any{},
	})
	prob := ValidateProviderLocation(context.Background(), raw, testProviderLocationStructural(t))
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

func baseLocDoc(spec map[string]any) map[string]any {
	return map[string]any{
		"apiVersion": resources.FabricAPIVersion,
		"kind":       resources.KindProviderLocation,
		"metadata": map[string]any{
			"name": "loc-boundary",
			"scopeRef": map[string]any{
				"apiVersion": resources.FabricAPIVersion,
				"kind":       string(apimeta.ScopeProvider),
				"name":       "sovereign-provider-a",
			},
		},
		"spec": spec,
	}
}

func mustProviderLocationJSON(t *testing.T, doc map[string]any) []byte {
	t.Helper()
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return raw
}

func assertProblemViolation(
	t *testing.T,
	prob *apiproblem.Problem,
	wantCode apiproblem.ErrorCode,
	wantField string,
	wantVCode apiproblem.ViolationCode,
) {
	t.Helper()
	if prob == nil {
		t.Fatalf("expected problem with field %s, got nil", wantField)
	}
	if prob.Code != wantCode {
		t.Fatalf("problem code = %q, want %q (problem=%#v)", prob.Code, wantCode, prob)
	}
	if len(prob.Violations) == 0 {
		t.Fatalf("expected violations, got none (problem=%#v)", prob)
	}
	found := false
	for _, v := range prob.Violations {
		if v.Field == wantField && v.Code == wantVCode {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("missing violation field=%q code=%q in %#v", wantField, wantVCode, prob.Violations)
	}
}

// testProviderLocationStructural builds the existing apiconform structural
// validator for tests only. Production ValidateProviderLocation depends solely
// on the injected apivalid.StructuralValidator interface.
func testProviderLocationStructural(t *testing.T) apivalid.StructuralValidator {
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

// passThroughStructural is a test-only stub that reports structural success so
// semantic ValidateRef paths can be exercised in isolation.
type passThroughStructural struct{}

func (passThroughStructural) Validate(any, string) ([]apiproblem.Violation, error) {
	return nil, nil
}

type unavailableStructural struct{}

func (unavailableStructural) Validate(any, string) ([]apiproblem.Violation, error) {
	return nil, errors.New("structural validator unavailable")
}

func providerLocationModuleRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("go.mod not found walking up from %s", wd)
		}
		dir = parent
	}
}
