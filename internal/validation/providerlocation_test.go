package validation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateProviderLocation_PositiveWithAndWithoutGeo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

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
	if prob := ValidateProviderLocation(ctx, withoutGeo); prob != nil {
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
	if prob := ValidateProviderLocation(ctx, withGeo); prob != nil {
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
	prob := ValidateProviderLocation(context.Background(), raw)
	assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/scopeRef/kind", apivalid.ViolationInvalidEnum)
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
	prob := ValidateProviderLocation(context.Background(), raw)
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
	prob := ValidateProviderLocation(context.Background(), raw)
	assertProblemViolation(t, prob, apiproblem.CodeUnknownField, "/region", apiproblem.ViolationUnknownField)
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
			wantCode:  apiproblem.ViolationOutOfRange,
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
			wantCode:  apiproblem.ViolationOutOfRange,
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
			prob := ValidateProviderLocation(context.Background(), raw)
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
	prob := ValidateProviderLocation(context.Background(), raw)
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
	prob := ValidateProviderLocation(context.Background(), raw)
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

	t.Run("countryCode-exactly-2", func(t *testing.T) {
		t.Parallel()
		ok := mustProviderLocationJSON(t, baseLocDoc(map[string]any{
			"geo": map[string]any{"countryCode": "US"},
		}))
		if prob := ValidateProviderLocation(ctx, ok); prob != nil {
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
		if prob := ValidateProviderLocation(ctx, ok); prob != nil {
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
		prob := ValidateProviderLocation(ctx, bad)
		assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/spec/geo/subdivisionCode", apiproblem.ViolationOutOfRange)
	})
}

func TestValidateProviderLocation_BoundaryNameLabelAnnotationLimits(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	lim := apivalid.DefaultLimits()

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
		if prob := ValidateProviderLocation(ctx, raw); prob != nil {
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
		prob := ValidateProviderLocation(ctx, raw)
		assertProblemViolation(t, prob, apiproblem.CodeValidationFailed, "/metadata/name", apivalid.ViolationInvalidResourceName)
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
		if prob := ValidateProviderLocation(ctx, raw); prob != nil {
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
		prob := ValidateProviderLocation(ctx, raw)
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
		prob := ValidateProviderLocation(ctx, raw)
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
	prob := ValidateProviderLocation(context.Background(), raw)
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
	prob := ValidateProviderLocation(context.Background(), raw)
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
