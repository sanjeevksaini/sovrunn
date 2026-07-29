package validate

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

func validProfile() decision.DecisionProfile {
	return decision.DecisionProfile{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: decision.APIVersionDecisionProfile,
			Kind:       decision.KindDecisionProfile,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "placement-decision",
			UID:  "profile-uid-1",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: "core.sovrunn.io/v1alpha1",
					Kind:       string(apimeta.ScopeProject),
					Name:       "payments-production",
					UID:        "project-123",
				},
			},
		},
		Spec: decision.ProfileSpec{
			Name:            "PlacementDecision",
			Family:          "placement",
			Version:         "1.0.0",
			Owner:           "platform-governance",
			LifecycleStatus: LifecycleStatusActive,
			PrimaryForm:     decision.DecisionFormSelection,
			AllowedAuthorityLevels: []decision.DecisionAuthority{
				decision.DecisionAuthorityAuthoritative,
			},
			InputSchemaRef:       "schemas/placement-input@1",
			TypedResultSchemaRef: "schemas/placement-result@1",
			RationaleSchemaRef:   "schemas/placement-rationale@1",
			ObligationSchemaRef:  "schemas/placement-obligation@1",
			AcceptedEvaluationTypes: []decision.AcceptedEvaluationType{
				{Type: "capacity-evaluator", Version: "1.2.0"},
			},
			AcceptedCompositionStrategies: []string{"all-of-v1"},
			ScopeSemantics: decision.ProfileScopeSemantics{
				AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeProject},
			},
			FailureBehavior: decision.FailureBehavior{
				FailPosture:          decision.FailPostureFailClosed,
				TimeoutBehavior:      "fail-closed",
				InsufficientEvidence: "fail-closed",
				DeclaredFailureModes: []string{"timeout", "evaluator-unavailable"},
			},
			SensitivityCeiling:  decision.SensitivityInternal,
			IdempotencyIdentity: "placement-v1",
			Limits: decision.ProfileLimits{
				MaxNodes:        16,
				MaxEdges:        32,
				MaxDepth:        8,
				MaxWidth:        8,
				MaxFanOut:       4,
				MaxPayloadBytes: 65536,
				MaxLatencyMs:    2000,
			},
		},
	}
}

func validException(ownerName string) decision.SecurityExceptionRef {
	return decision.SecurityExceptionRef{
		ApprovalRef: apimeta.TypedRef{
			APIVersion: "governance.sovrunn.io/v1alpha1",
			Kind:       "SecurityExceptionApproval",
			Name:       "break-glass-1",
			UID:        "exception-approval-uid-1",
		},
		ExceptionScopeRef: &apimeta.ScopeRef{
			TypedRef: apimeta.TypedRef{
				APIVersion: "core.sovrunn.io/v1alpha1",
				Kind:       string(apimeta.ScopeProject),
				Name:       "payments-production",
				UID:        "project-123",
			},
		},
		OwnerRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       "Organization",
			Name:       ownerName,
			UID:        "owner-uid-1",
		},
		ApprovingAuthorityRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       "OrganizationUnit",
			Name:       "security-board",
			UID:        "ou-uid-1",
		},
		CompensatingControls: []string{"ctrl:dual-control"},
		EffectiveFrom:        "2026-07-01T00:00:00Z",
		EffectiveUntil:       "2026-07-31T23:59:59Z",
		AuditTreatment:       "audit:mandatory-durable",
		ReassessmentTrigger:  "trigger:quarterly-review",
		CoveredFailureModes:  []string{"timeout"},
		Purpose:              "bounded break-glass under dual control",
	}
}

func assertProfileViolation(t *testing.T, p *apiproblem.Problem, wantCode apiproblem.ViolationCode, wantField string) {
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

func TestValidateProfile_HappyPath(t *testing.T) {
	t.Parallel()

	p := ValidateProfileDocument(validProfile())
	if p != nil {
		t.Fatalf("valid profile must pass, got %#v", p)
	}
}

func TestValidateProfile_Unknown(t *testing.T) {
	t.Parallel()

	p := ValidateProfile(ProfileInput{
		CheckLookup: true,
		Ref:         decision.ProfileRef{Name: "MissingProfile", Version: "1.0.0"},
		IDExists:    false,
		Found:       false,
	})
	assertProfileViolation(t, p, CodeProfileUnknown, ptrProfileRefName)
}

func TestValidateProfile_VersionUnsupported(t *testing.T) {
	t.Parallel()

	p := ValidateProfile(ProfileInput{
		CheckLookup: true,
		Ref:         decision.ProfileRef{Name: "PlacementDecision", Version: "9.9.9"},
		IDExists:    true,
		Found:       false,
	})
	assertProfileViolation(t, p, CodeProfileVersionUnsupported, ptrProfileRefVersion)
}

func TestValidateProfile_Inactive(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	prof.Spec.LifecycleStatus = "Inactive"
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileInactive, ptrSpecLifecycleStatus)
}

func TestValidateProfile_SchemaInvalidPrimaryForm(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	prof.Spec.PrimaryForm = decision.DecisionForm("NOT_A_FORM")
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileSchemaInvalid, ptrSpecPrimaryForm)
}

func TestValidateProfile_FailOpenMissingException(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	prof.Spec.FailureBehavior.FailPosture = decision.FailPostureFailOpen
	prof.Spec.FailureBehavior.SecurityExceptionRef = nil
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileSchemaInvalid, ptrSpecSecurityException)
}

func TestValidateProfile_SecurityExceptionMissingUID(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	ex := validException(prof.Spec.Owner)
	ex.ApprovalRef.UID = ""
	prof.Spec.FailureBehavior.FailPosture = decision.FailPostureFailOpen
	prof.Spec.FailureBehavior.SecurityExceptionRef = &ex
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileSchemaInvalid, ptrSpecExceptionUID)
}

func TestValidateProfile_SecurityExceptionOwnerMismatch(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	ex := validException("other-owner")
	prof.Spec.FailureBehavior.FailPosture = decision.FailPostureFailOpen
	prof.Spec.FailureBehavior.SecurityExceptionRef = &ex
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileSchemaInvalid, ptrSpecExceptionOwner)
}

func TestValidateProfile_SecurityExceptionScopeMismatch(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	ex := validException(prof.Spec.Owner)
	ex.ExceptionScopeRef = &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeTenant),
			Name:       "other-tenant",
			UID:        "tenant-999",
		},
	}
	prof.Spec.ScopeSemantics.AllowedScopes = []apimeta.ScopeKind{apimeta.ScopeProject, apimeta.ScopeTenant}
	prof.Spec.FailureBehavior.FailPosture = decision.FailPostureFailOpen
	prof.Spec.FailureBehavior.SecurityExceptionRef = &ex
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileSchemaInvalid, ptrSpecExceptionScope)
}

func TestValidateProfile_SecurityExceptionExpiredInterval(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	ex := validException(prof.Spec.Owner)
	ex.EffectiveFrom = "2026-07-31T00:00:00Z"
	ex.EffectiveUntil = "2026-07-01T00:00:00Z"
	prof.Spec.FailureBehavior.FailPosture = decision.FailPostureFailOpen
	prof.Spec.FailureBehavior.SecurityExceptionRef = &ex
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileSchemaInvalid, ptrSpecExceptionUntil)
}

func TestValidateProfile_SecurityExceptionExpiredVsDecisionTime(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	ex := validException(prof.Spec.Owner)
	prof.Spec.FailureBehavior.FailPosture = decision.FailPostureFailOpen
	prof.Spec.FailureBehavior.SecurityExceptionRef = &ex
	p := ValidateProfile(ProfileInput{
		Profile:      prof,
		DecisionTime: "2026-08-01T00:00:00Z",
	})
	assertProfileViolation(t, p, CodeProfileSchemaInvalid, ptrSpecExceptionUntil)
}

func TestValidateProfile_SecurityExceptionEmptyControls(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	ex := validException(prof.Spec.Owner)
	ex.CompensatingControls = nil
	prof.Spec.FailureBehavior.FailPosture = decision.FailPostureFailOpen
	prof.Spec.FailureBehavior.SecurityExceptionRef = &ex
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileSchemaInvalid, ptrSpecExceptionControls)
}

func TestValidateProfile_SecurityExceptionEmptyModes(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	ex := validException(prof.Spec.Owner)
	ex.CoveredFailureModes = []string{}
	prof.Spec.FailureBehavior.FailPosture = decision.FailPostureFailOpen
	prof.Spec.FailureBehavior.SecurityExceptionRef = &ex
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileSchemaInvalid, ptrSpecExceptionModes)
}

func TestValidateProfile_SecurityExceptionUnknownMode(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	ex := validException(prof.Spec.Owner)
	ex.CoveredFailureModes = []string{"not-declared"}
	prof.Spec.FailureBehavior.FailPosture = decision.FailPostureFailOpen
	prof.Spec.FailureBehavior.SecurityExceptionRef = &ex
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileSchemaInvalid, ptrSpecExceptionModes+"/0")
}

func TestValidateProfile_SecurityExceptionHappyPath(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	ex := validException(prof.Spec.Owner)
	prof.Spec.FailureBehavior.FailPosture = decision.FailPostureFailOpen
	prof.Spec.FailureBehavior.SecurityExceptionRef = &ex
	p := ValidateProfile(ProfileInput{
		Profile:      prof,
		DecisionTime: "2026-07-15T12:00:00Z",
	})
	if p != nil {
		t.Fatalf("valid fail-open exception must pass, got %#v", p)
	}
}

func TestValidateProfile_LimitAbovePlatformCeilingPayload(t *testing.T) {
	t.Parallel()

	plat := InheritedPlatformCeilings()
	prof := validProfile()
	prof.Spec.Limits.MaxPayloadBytes = plat.MaxPayloadBytes + 1
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileLimitInvalid, ptrSpecLimits+"/maxPayloadBytes")
}

func TestValidateProfile_LimitAbovePlatformCeilingDepth(t *testing.T) {
	t.Parallel()

	plat := InheritedPlatformCeilings()
	prof := validProfile()
	prof.Spec.Limits.MaxDepth = plat.MaxDepth + 1
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileLimitInvalid, ptrSpecLimits+"/maxDepth")
}

func TestValidateProfile_MissingMandatoryCompositionLimit(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	prof.Spec.Limits.MaxNodes = 0
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileLimitInvalid, ptrSpecLimits+"/maxNodes")
}

func TestValidateProfile_EvaluatorNarrowsOK(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	narrow := decision.ProfileLimits{
		MaxNodes:        8,
		MaxPayloadBytes: 4096,
		MaxLatencyMs:    500,
	}
	prof.Spec.AcceptedEvaluationTypes = []decision.AcceptedEvaluationType{
		{Type: "capacity-evaluator", Version: "1.2.0", Limits: &narrow},
	}
	p := ValidateProfileDocument(prof)
	if p != nil {
		t.Fatalf("evaluator narrowing must pass, got %#v", p)
	}
}

func TestValidateProfile_EvaluatorWideningRejected(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	wide := decision.ProfileLimits{
		MaxPayloadBytes: prof.Spec.Limits.MaxPayloadBytes + 1,
	}
	prof.Spec.AcceptedEvaluationTypes = []decision.AcceptedEvaluationType{
		{Type: "capacity-evaluator", Version: "1.2.0", Limits: &wide},
	}
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileLimitInvalid,
		"/spec/acceptedEvaluationTypes/0/limits/maxPayloadBytes")
}

func TestValidateProfile_EvaluatorIncomparableCostBudget(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	eval := decision.ProfileLimits{MaxCostBudget: "budget:gold"}
	prof.Spec.AcceptedEvaluationTypes = []decision.AcceptedEvaluationType{
		{Type: "capacity-evaluator", Limits: &eval},
	}
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileLimitInvalid,
		"/spec/acceptedEvaluationTypes/0/limits/maxCostBudget")
}

func TestInheritedPlatformCeilingsMatchFEATURE0012(t *testing.T) {
	t.Parallel()

	got := InheritedPlatformCeilings()
	want := apivalid.DefaultLimits()
	if got.MaxPayloadBytes != int64(want.MaxObjectBytes) {
		t.Fatalf("MaxPayloadBytes ceiling = %d, want %d", got.MaxPayloadBytes, want.MaxObjectBytes)
	}
	if got.MaxDepth != want.MaxNestingDepth {
		t.Fatalf("MaxDepth ceiling = %d, want %d", got.MaxDepth, want.MaxNestingDepth)
	}
	if got.MaxFieldCount != want.MaxLabels {
		t.Fatalf("MaxFieldCount ceiling = %d, want %d", got.MaxFieldCount, want.MaxLabels)
	}
}

func TestEffectiveIntLimit_MostRestrictiveWins(t *testing.T) {
	t.Parallel()

	if got := EffectiveIntLimit(32, 16, 8); got != 8 {
		t.Fatalf("EffectiveIntLimit(32,16,8) = %d, want 8", got)
	}
	if got := EffectiveIntLimit(32, 16, 0); got != 16 {
		t.Fatalf("EffectiveIntLimit(32,16,0) = %d, want 16", got)
	}
	if got := EffectiveInt64Limit(1_048_576, 65536, 4096); got != 4096 {
		t.Fatalf("EffectiveInt64Limit = %d, want 4096", got)
	}
}

func TestValidateProfile_LookupThenInactive(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	prof.Spec.LifecycleStatus = "Deprecated"
	p := ValidateProfile(ProfileInput{
		Profile:     prof,
		CheckLookup: true,
		Ref:         decision.ProfileRef{Name: prof.Spec.Name, Version: prof.Spec.Version},
		IDExists:    true,
		Found:       true,
	})
	assertProfileViolation(t, p, CodeProfileInactive, ptrSpecLifecycleStatus)
}

func TestValidateProfile_SecurityClassFailOpenRequiresException(t *testing.T) {
	t.Parallel()

	prof := validProfile()
	prof.Spec.Family = "authorization"
	prof.Spec.FailureBehavior.FailPosture = decision.FailPostureFailOpen
	prof.Spec.FailureBehavior.SecurityExceptionRef = nil
	p := ValidateProfileDocument(prof)
	assertProfileViolation(t, p, CodeProfileSchemaInvalid, ptrSpecSecurityException)
}
