package validate

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

func assertSensitivityViolation(t *testing.T, p *apiproblem.Problem, wantCode apiproblem.ViolationCode, wantField string) {
	t.Helper()
	if p == nil {
		t.Fatalf("expected Problem with %s at %s, got nil", wantCode, wantField)
	}
	if p.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("top-level code = %q, want %q", p.Code, apiproblem.CodeValidationFailed)
	}
	if p.Status != 422 {
		t.Fatalf("status = %d, want 422", p.Status)
	}
	if p.Type != "urn:sovrunn:problem:validation-failed" {
		t.Fatalf("type = %q, want urn:sovrunn:problem:validation-failed", p.Type)
	}
	if len(p.Violations) != 1 {
		t.Fatalf("violations len=%d, want 1: %#v", len(p.Violations), p.Violations)
	}
	v := p.Violations[0]
	if v.Code != wantCode {
		t.Fatalf("violations[0].code = %q, want %q", v.Code, wantCode)
	}
	if v.Field != wantField {
		t.Fatalf("violations[0].field = %q, want %q", v.Field, wantField)
	}
	if v.Message == "" {
		t.Fatal("violations[0].message must be non-empty (redactable)")
	}
}

func sensitivityProfile(ceiling decision.Sensitivity) decision.DecisionProfile {
	p := validProfile()
	p.Spec.SensitivityCeiling = ceiling
	p.Spec.ContentCategories = []string{"capacity", "placement"}
	p.Spec.ProjectionRules = []decision.ProjectionRule{
		{
			Audience: "operator",
			Includes: []string{"/record/result/typedResult", "/record/result/rationale"},
			Excludes: []string{"/record/trust", "/record/sovereignty"},
		},
	}
	return p
}

func sensitivityCanonicalView() []string {
	return []string{
		"/record/result/typedResult",
		"/record/result/rationale",
		"/record/trust",
		"/record/sovereignty",
		"/record/evaluationResults",
	}
}

func TestValidateSensitivity_HappyPath(t *testing.T) {
	t.Parallel()

	p := ValidateSensitivity(SensitivityInput{
		Profile:              sensitivityProfile(decision.SensitivityConfidential),
		Classifications:      []apimeta.DataClassification{apimeta.ClassInternal},
		DeclaredSensitivity:  decision.SensitivityConfidential, // raised above INTERNAL floor
		ContentCategories:    []string{"capacity"},
		ProhibitedCategories: []string{"raw-secret", "unredacted-pii"},
		CanonicalView:        sensitivityCanonicalView(),
		AllowedAudiences:     []string{"operator", "auditor", "ai"},
	})
	if p != nil {
		t.Fatalf("valid sensitivity/projection input must pass, got %#v", p)
	}
}

func TestValidateSensitivity_HappyPathNoCapturedContent(t *testing.T) {
	t.Parallel()

	// No classifications, declared sensitivity, or content categories: projection-only.
	profile := sensitivityProfile(decision.SensitivityInternal)
	p := ValidateSensitivity(SensitivityInput{
		Profile:       profile,
		CanonicalView: sensitivityCanonicalView(),
	})
	if p != nil {
		t.Fatalf("projection-only validation must pass, got %#v", p)
	}
}

func TestValidateSensitivity_CeilingBoundaryAccepted(t *testing.T) {
	t.Parallel()

	// Declared sensitivity exactly at the ceiling is accepted (F13-SEC-01).
	p := ValidateSensitivity(SensitivityInput{
		Profile:             sensitivityProfile(decision.SensitivityInternal),
		Classifications:     []apimeta.DataClassification{apimeta.ClassPublic},
		DeclaredSensitivity: decision.SensitivityInternal,
		CanonicalView:       sensitivityCanonicalView(),
	})
	if p != nil {
		t.Fatalf("sensitivity at ceiling must pass, got %#v", p)
	}
}

func TestValidateSensitivity_CeilingBoundaryExceeded(t *testing.T) {
	t.Parallel()

	// Captured/declared sensitivity above profile ceiling fails (F13-SEC-005/01).
	p := ValidateSensitivity(SensitivityInput{
		Profile:             sensitivityProfile(decision.SensitivityInternal),
		Classifications:     []apimeta.DataClassification{apimeta.ClassPublic},
		DeclaredSensitivity: decision.SensitivityConfidential,
		CanonicalView:       sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeEvaluationLimitExceeded, ptrDeclaredSensitivity)
}

func TestValidateSensitivity_CeilingRestrictedAcceptsAll(t *testing.T) {
	t.Parallel()

	p := ValidateSensitivity(SensitivityInput{
		Profile:             sensitivityProfile(decision.SensitivityRestricted),
		Classifications:     []apimeta.DataClassification{apimeta.ClassSensitive},
		DeclaredSensitivity: decision.SensitivityRestricted,
		CanonicalView:       sensitivityCanonicalView(),
	})
	if p != nil {
		t.Fatalf("RESTRICTED ceiling must accept RESTRICTED declared, got %#v", p)
	}
}

func TestValidateSensitivity_FloorCannotBeLowered(t *testing.T) {
	t.Parallel()

	// Customer-visible floor is CONFIDENTIAL; declaring INTERNAL lowers it.
	p := ValidateSensitivity(SensitivityInput{
		Profile:             sensitivityProfile(decision.SensitivityRestricted),
		Classifications:     []apimeta.DataClassification{apimeta.ClassCustomerVisible},
		DeclaredSensitivity: decision.SensitivityInternal,
		CanonicalView:       sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, ptrDeclaredSensitivity)
}

func TestValidateSensitivity_FloorRaiseOnlyAccepted(t *testing.T) {
	t.Parallel()

	// Floor CONFIDENTIAL; raising to RESTRICTED is allowed (F13-SEC-003).
	p := ValidateSensitivity(SensitivityInput{
		Profile:             sensitivityProfile(decision.SensitivityRestricted),
		Classifications:     []apimeta.DataClassification{apimeta.ClassTenantConfidential},
		DeclaredSensitivity: decision.SensitivityRestricted,
		CanonicalView:       sensitivityCanonicalView(),
	})
	if p != nil {
		t.Fatalf("raising floor must pass, got %#v", p)
	}
}

func TestValidateSensitivity_UnmappedClassificationFailsClosed(t *testing.T) {
	t.Parallel()

	p := ValidateSensitivity(SensitivityInput{
		Profile:             sensitivityProfile(decision.SensitivityRestricted),
		Classifications:     []apimeta.DataClassification{apimeta.DataClassification("Top-secret")},
		DeclaredSensitivity: decision.SensitivityRestricted,
		CanonicalView:       sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, classificationPtr(0))
}

func TestValidateSensitivity_BelowFloorPublicVsSensitive(t *testing.T) {
	t.Parallel()

	p := ValidateSensitivity(SensitivityInput{
		Profile:             sensitivityProfile(decision.SensitivityRestricted),
		Classifications:     []apimeta.DataClassification{apimeta.ClassSensitive},
		DeclaredSensitivity: decision.SensitivityPublic,
		CanonicalView:       sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, ptrDeclaredSensitivity)
}

func TestValidateSensitivity_MostRestrictiveFloor(t *testing.T) {
	t.Parallel()

	// Public→PUBLIC and Sensitive→RESTRICTED → most restrictive floor RESTRICTED;
	// declaring CONFIDENTIAL is below that floor.
	p := ValidateSensitivity(SensitivityInput{
		Profile: sensitivityProfile(decision.SensitivityRestricted),
		Classifications: []apimeta.DataClassification{
			apimeta.ClassPublic,
			apimeta.ClassSensitive,
		},
		DeclaredSensitivity: decision.SensitivityConfidential,
		CanonicalView:       sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, ptrDeclaredSensitivity)
}

func TestValidateSensitivity_CeilingBelowFloorContradictory(t *testing.T) {
	t.Parallel()

	// Sensitive floor is RESTRICTED; ceiling INTERNAL contradicts.
	p := ValidateSensitivity(SensitivityInput{
		Profile:             sensitivityProfile(decision.SensitivityInternal),
		Classifications:     []apimeta.DataClassification{apimeta.ClassSensitive},
		DeclaredSensitivity: decision.SensitivityRestricted,
		CanonicalView:       sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, ptrSensitivityCeiling)
}

func TestValidateSensitivity_ProhibitedCategoryRejected(t *testing.T) {
	t.Parallel()

	p := ValidateSensitivity(SensitivityInput{
		Profile:              sensitivityProfile(decision.SensitivityConfidential),
		ContentCategories:    []string{"raw-secret"},
		ProhibitedCategories: []string{"raw-secret"},
		CanonicalView:        sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeEvaluationResultInvalid, contentCategoryPtr(0))
}

func TestValidateSensitivity_CategoryNotInAllowedSet(t *testing.T) {
	t.Parallel()

	p := ValidateSensitivity(SensitivityInput{
		Profile:           sensitivityProfile(decision.SensitivityConfidential),
		ContentCategories: []string{"unregistered-category"},
		CanonicalView:     sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeEvaluationResultInvalid, contentCategoryPtr(0))
}

func TestValidateSensitivity_AllowedAndProhibitedContradiction(t *testing.T) {
	t.Parallel()

	profile := sensitivityProfile(decision.SensitivityConfidential)
	profile.Spec.ContentCategories = []string{"capacity", "raw-secret"}
	p := ValidateSensitivity(SensitivityInput{
		Profile:              profile,
		ProhibitedCategories: []string{"raw-secret"},
		CanonicalView:        sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, "/spec/contentCategories")
}

func TestValidateSensitivity_OverlappingIncludesExcludesEquality(t *testing.T) {
	t.Parallel()

	profile := sensitivityProfile(decision.SensitivityConfidential)
	profile.Spec.ProjectionRules = []decision.ProjectionRule{{
		Audience: "operator",
		Includes: []string{"/record/trust"},
		Excludes: []string{"/record/trust"},
	}}
	p := ValidateSensitivity(SensitivityInput{
		Profile:       profile,
		CanonicalView: sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, "/spec/projectionRules/0/includes/0")
}

func TestValidateSensitivity_OverlappingIncludesExcludesContainment(t *testing.T) {
	t.Parallel()

	profile := sensitivityProfile(decision.SensitivityConfidential)
	profile.Spec.ProjectionRules = []decision.ProjectionRule{{
		Audience: "operator",
		Includes: []string{"/record/result"},
		Excludes: []string{"/record/result/typedResult"},
	}}
	view := append(sensitivityCanonicalView(), "/record/result")
	p := ValidateSensitivity(SensitivityInput{
		Profile:       profile,
		CanonicalView: view,
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, "/spec/projectionRules/0/includes/0")
}

func TestValidateSensitivity_OverlappingExcludesContainsInclude(t *testing.T) {
	t.Parallel()

	profile := sensitivityProfile(decision.SensitivityConfidential)
	profile.Spec.ProjectionRules = []decision.ProjectionRule{{
		Audience: "operator",
		Includes: []string{"/record/result/rationale/reasonCodes"},
		Excludes: []string{"/record/result/rationale"},
	}}
	view := append(sensitivityCanonicalView(),
		"/record/result/rationale",
		"/record/result/rationale/reasonCodes",
	)
	p := ValidateSensitivity(SensitivityInput{
		Profile:       profile,
		CanonicalView: view,
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, "/spec/projectionRules/0/includes/0")
}

func TestValidateSensitivity_NonOverlappingSiblingPointersAccepted(t *testing.T) {
	t.Parallel()

	profile := sensitivityProfile(decision.SensitivityConfidential)
	profile.Spec.ProjectionRules = []decision.ProjectionRule{{
		Audience: "operator",
		Includes: []string{"/record/result/typedResult"},
		Excludes: []string{"/record/result/rationale"},
	}}
	p := ValidateSensitivity(SensitivityInput{
		Profile:       profile,
		CanonicalView: sensitivityCanonicalView(),
	})
	if p != nil {
		t.Fatalf("sibling pointers must not overlap, got %#v", p)
	}
}

func TestValidateSensitivity_UnknownProjectionPointer(t *testing.T) {
	t.Parallel()

	profile := sensitivityProfile(decision.SensitivityConfidential)
	profile.Spec.ProjectionRules = []decision.ProjectionRule{{
		Audience: "operator",
		Includes: []string{"/record/notInCanonicalView"},
	}}
	p := ValidateSensitivity(SensitivityInput{
		Profile:       profile,
		CanonicalView: sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, "/spec/projectionRules/0/includes/0")
}

func TestValidateSensitivity_InvalidProjectionPointerSyntax(t *testing.T) {
	t.Parallel()

	profile := sensitivityProfile(decision.SensitivityConfidential)
	profile.Spec.ProjectionRules = []decision.ProjectionRule{{
		Audience: "operator",
		Includes: []string{"record/missing-leading-slash"},
	}}
	p := ValidateSensitivity(SensitivityInput{
		Profile:       profile,
		CanonicalView: []string{"record/missing-leading-slash"},
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, "/spec/projectionRules/0/includes/0")
}

func TestValidateSensitivity_AudienceMismatch(t *testing.T) {
	t.Parallel()

	p := ValidateSensitivity(SensitivityInput{
		Profile:          sensitivityProfile(decision.SensitivityConfidential),
		CanonicalView:    sensitivityCanonicalView(),
		AllowedAudiences: []string{"customer", "auditor"},
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, "/spec/projectionRules/0/audience")
}

func TestValidateSensitivity_EmptyAudienceRejected(t *testing.T) {
	t.Parallel()

	profile := sensitivityProfile(decision.SensitivityConfidential)
	profile.Spec.ProjectionRules = []decision.ProjectionRule{{
		Audience: "",
		Includes: []string{"/record/trust"},
	}}
	p := ValidateSensitivity(SensitivityInput{
		Profile:       profile,
		CanonicalView: sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, "/spec/projectionRules/0/audience")
}

func TestValidateSensitivity_InvalidCeilingVocabulary(t *testing.T) {
	t.Parallel()

	profile := sensitivityProfile(decision.Sensitivity("TOP_SECRET"))
	p := ValidateSensitivity(SensitivityInput{Profile: profile})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, ptrSensitivityCeiling)
}

func TestValidateSensitivity_MissingDeclaredWhenClassificationsPresent(t *testing.T) {
	t.Parallel()

	p := ValidateSensitivity(SensitivityInput{
		Profile:         sensitivityProfile(decision.SensitivityRestricted),
		Classifications: []apimeta.DataClassification{apimeta.ClassPublic},
		CanonicalView:   sensitivityCanonicalView(),
	})
	assertSensitivityViolation(t, p, CodeProfileSchemaInvalid, ptrDeclaredSensitivity)
}

func TestPointersOverlap(t *testing.T) {
	t.Parallel()

	cases := []struct {
		a, b string
		want bool
	}{
		{"/a", "/a", true},
		{"/a", "/a/b", true},
		{"/a/b", "/a", true},
		{"/a", "/ab", false},
		{"/a/b", "/a/c", false},
		{"", "/anything", true},
		{"/x", "/y", false},
	}
	for _, tc := range cases {
		if got := pointersOverlap(tc.a, tc.b); got != tc.want {
			t.Fatalf("pointersOverlap(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestValidJSONPointer(t *testing.T) {
	t.Parallel()

	if !validJSONPointer("") {
		t.Fatal("empty pointer must be valid")
	}
	if !validJSONPointer("/a~0b/c~1d") {
		t.Fatal("escaped pointer must be valid")
	}
	if validJSONPointer("no-slash") {
		t.Fatal("missing leading slash must be invalid")
	}
	if validJSONPointer("/a~") {
		t.Fatal("truncated escape must be invalid")
	}
	if validJSONPointer("/a~2") {
		t.Fatal("unknown escape must be invalid")
	}
}
