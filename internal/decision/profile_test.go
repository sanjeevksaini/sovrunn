package decision

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func sampleSecurityExceptionRef() SecurityExceptionRef {
	return SecurityExceptionRef{
		ApprovalRef: apimeta.TypedRef{
			APIVersion: "governance.sovrunn.io/v1alpha1",
			Kind:       "SecurityExceptionApproval",
			Name:       "break-glass-placement-1",
			UID:        "exception-approval-uid-1",
		},
		ExceptionScopeRef: &apimeta.ScopeRef{
			TypedRef: apimeta.TypedRef{
				APIVersion: "core.sovrunn.io/v1alpha1",
				Kind:       "Project",
				Name:       "payments-production",
				UID:        "project-123",
			},
		},
		OwnerRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       "Organization",
			Name:       "acme",
			UID:        "org-uid-1",
		},
		ApprovingAuthorityRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       "OrganizationUnit",
			Name:       "security-board",
			UID:        "ou-uid-1",
		},
		CompensatingControls: []string{
			"ctrl:dual-control-review",
			"ctrl:time-boxed-session",
		},
		EffectiveFrom:       "2026-07-01T00:00:00Z",
		EffectiveUntil:      "2026-07-31T23:59:59Z",
		AuditTreatment:      "audit:mandatory-durable",
		ReassessmentTrigger: "trigger:quarterly-security-review",
		CoveredFailureModes: []string{
			"timeout",
			"evaluator-unavailable",
		},
		Purpose: "bounded break-glass placement under dual control",
	}
}

func sampleDecisionProfile() DecisionProfile {
	return DecisionProfile{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: APIVersionDecisionProfile,
			Kind:       KindDecisionProfile,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "placement-decision",
			UID:  "profile-uid-1",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: "core.sovrunn.io/v1alpha1",
					Kind:       "Organization",
					Name:       "acme",
					UID:        "org-uid-1",
				},
			},
		},
		Spec: ProfileSpec{
			Name:            "PlacementDecision",
			Family:          "placement",
			Version:         "1.0.0",
			Owner:           "platform-governance",
			LifecycleStatus: "Active",
			PrimaryForm:     DecisionFormSelection,
			SecondaryFacets: []DecisionForm{
				DecisionFormAssessment,
			},
			AllowedAuthorityLevels: []DecisionAuthority{
				DecisionAuthorityAuthoritative,
				DecisionAuthorityAdvisory,
			},
			InputSchemaRef:       "schemas/placement-input@1",
			TypedResultSchemaRef: "schemas/placement-result@1",
			RationaleSchemaRef:   "schemas/placement-rationale@1",
			ObligationSchemaRef:  "schemas/placement-obligation@1",
			AcceptedEvaluationTypes: []AcceptedEvaluationType{
				{Type: "capacity-evaluator", Version: "1.2.0"},
			},
			AcceptedCompositionStrategies: []string{"all-of-v1"},
			ScopeSemantics: ProfileScopeSemantics{
				AllowedScopes: []apimeta.ScopeKind{
					apimeta.ScopeProject,
					apimeta.ScopeTenant,
				},
				PlatformAllowed: false,
			},
			ActorSemantics:   ProfileActorSemantics{RequireActor: true},
			SubjectSemantics: ProfileSubjectSemantics{RequireSubject: true, MaxSubjects: 8},
			PolicySemantics:  ProfilePolicySemantics{RequireEffectivePolicyContext: true},
			EvidenceSemantics: ProfileEvidenceSemantics{
				RequireEvidence: true,
				AcceptedEvidenceKinds: []string{
					"InputSnapshot",
				},
			},
			AuditSemantics: ProfileAuditSemantics{RequireDurableRecord: true},
			DeterminismBehavior: DeterminismBehavior{
				RequireDeterministic:    true,
				AllowNondeterministicAI: false,
			},
			FailureBehavior: FailureBehavior{
				FailPosture:          FailPostureFailClosed,
				TimeoutBehavior:      "fail-closed",
				InsufficientEvidence: "fail-closed",
				DeclaredFailureModes: []string{"timeout", "evaluator-unavailable"},
			},
			ValidityRules: ValidityRules{
				AllowCorrection:   true,
				AllowSupersession: true,
				AllowRevocation:   true,
				MaxChainDepth:     4,
			},
			SensitivityCeiling: SensitivityInternal,
			ContentCategories:  []string{"capacity", "placement"},
			ProjectionRules: []ProjectionRule{
				{
					Audience: "operator",
					Includes: []string{"/record/result"},
					Excludes: []string{"/record/result/typedResult/internalScore"},
				},
			},
			IdempotencyIdentity: "placement-v1",
			Limits: ProfileLimits{
				MaxNodes:        16,
				MaxEdges:        32,
				MaxDepth:        8,
				MaxWidth:        8,
				MaxFanOut:       4,
				MaxPayloadBytes: 65536,
				MaxLatencyMs:    2000,
			},
			CompatibilityPolicy: CompatibilityPolicy{
				Owner:                "platform-governance",
				MinCompatibleVersion: "1.0.0",
				BackwardCompatible:   true,
			},
			FixtureRefs:          []string{"fixtures/placement-positive-1"},
			ReassessmentTriggers: []string{"trigger:new-evaluator-family"},
		},
	}
}

func TestDecisionProfileJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := sampleDecisionProfile()
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out DecisionProfile
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if out.APIVersion != APIVersionDecisionProfile || out.Kind != KindDecisionProfile {
		t.Fatalf("TypeMeta = {%q,%q}, want {%q,%q}", out.APIVersion, out.Kind, APIVersionDecisionProfile, KindDecisionProfile)
	}
	if out.Metadata.Name != "placement-decision" || out.Metadata.UID != "profile-uid-1" {
		t.Fatalf("metadata identity mismatch: %+v", out.Metadata)
	}
	if out.Metadata.ScopeRef == nil || out.Metadata.ScopeRef.UID != "org-uid-1" {
		t.Fatalf("metadata.scopeRef must remain sole scope authority, got %#v", out.Metadata.ScopeRef)
	}
	if out.Spec.Name != "PlacementDecision" || out.Spec.Version != "1.0.0" {
		t.Fatalf("spec identity = {%q,%q}", out.Spec.Name, out.Spec.Version)
	}
	if out.Spec.PrimaryForm != DecisionFormSelection {
		t.Fatalf("primaryForm = %q", out.Spec.PrimaryForm)
	}
	if out.Spec.FailureBehavior.FailPosture != FailPostureFailClosed {
		t.Fatalf("failPosture = %q", out.Spec.FailureBehavior.FailPosture)
	}
	if out.Spec.FailureBehavior.SecurityExceptionRef != nil {
		t.Fatalf("securityExceptionRef must be omitted for fail-closed, got %#v", out.Spec.FailureBehavior.SecurityExceptionRef)
	}
	if out.Spec.Limits.MaxNodes != 16 || out.Spec.Limits.MaxPayloadBytes != 65536 {
		t.Fatalf("limits = %+v", out.Spec.Limits)
	}
	if len(out.Spec.ProjectionRules) != 1 || out.Spec.ProjectionRules[0].Audience != "operator" {
		t.Fatalf("projectionRules = %#v", out.Spec.ProjectionRules)
	}
	if out.Spec.SensitivityCeiling != SensitivityInternal {
		t.Fatalf("sensitivityCeiling = %q", out.Spec.SensitivityCeiling)
	}
}

func TestDecisionProfileOmitemptyOptionalFields(t *testing.T) {
	t.Parallel()

	in := sampleDecisionProfile()
	in.Spec.SecondaryFacets = nil
	in.Spec.AcceptedEvaluationTypes = nil
	in.Spec.AcceptedCompositionStrategies = nil
	in.Spec.ContentCategories = nil
	in.Spec.ResidencyProfile = nil
	in.Spec.RetentionProfile = nil
	in.Spec.LegalHoldProfile = nil
	in.Spec.ProjectionRules = nil
	in.Spec.FixtureRefs = nil
	in.Spec.ReassessmentTriggers = nil
	in.Spec.FailureBehavior.DeclaredFailureModes = nil
	in.Spec.FailureBehavior.SecurityExceptionRef = nil
	in.Spec.SubjectSemantics.AllowedKinds = nil
	in.Spec.EvidenceSemantics.AcceptedEvidenceKinds = nil
	in.Spec.ValidityRules.DefaultValidityMs = 0
	in.Spec.ValidityRules.MaxChainDepth = 0
	in.Spec.Limits = ProfileLimits{} // all zero/empty → omitempty fields absent
	in.Spec.CompatibilityPolicy.Owner = ""
	in.Spec.CompatibilityPolicy.MinCompatibleVersion = ""

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	raw := string(data)

	for _, key := range []string{
		`"secondaryFacets"`,
		`"acceptedEvaluationTypes"`,
		`"acceptedCompositionStrategies"`,
		`"contentCategories"`,
		`"residencyProfile"`,
		`"retentionProfile"`,
		`"legalHoldProfile"`,
		`"projectionRules"`,
		`"fixtureRefs"`,
		`"reassessmentTriggers"`,
		`"declaredFailureModes"`,
		`"securityExceptionRef"`,
		`"maxNodes"`,
		`"maxPayloadBytes"`,
	} {
		if strings.Contains(raw, key) {
			t.Fatalf("optional field %s must be omitted when unset; json=%s", key, raw)
		}
	}

	for _, key := range []string{
		`"apiVersion"`,
		`"kind"`,
		`"metadata"`,
		`"spec"`,
		`"name"`,
		`"family"`,
		`"version"`,
		`"owner"`,
		`"lifecycleStatus"`,
		`"primaryForm"`,
		`"allowedAuthorityLevels"`,
		`"failureBehavior"`,
		`"failPosture"`,
		`"limits"`,
		`"idempotencyIdentity"`,
		`"compatibilityPolicy"`,
	} {
		if !strings.Contains(raw, key) {
			t.Fatalf("required field %s missing; json=%s", key, raw)
		}
	}
}

func TestSecurityExceptionRefRoundTripAndRequiredness(t *testing.T) {
	t.Parallel()

	in := sampleDecisionProfile()
	exc := sampleSecurityExceptionRef()
	in.Spec.FailureBehavior.FailPosture = FailPostureFailOpen
	in.Spec.FailureBehavior.SecurityExceptionRef = &exc

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out DecisionProfile
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	got := out.Spec.FailureBehavior.SecurityExceptionRef
	if got == nil {
		t.Fatal("securityExceptionRef must be present for fail-open sample")
	}

	// UID pinning requiredness (structural presence in contract).
	if got.ApprovalRef.UID == "" {
		t.Fatal("approvalRef.uid must be present (uid pinning)")
	}
	if got.OwnerRef.UID == "" {
		t.Fatal("ownerRef.uid must be present (uid pinning)")
	}
	if got.ApprovingAuthorityRef.UID == "" {
		t.Fatal("approvingAuthorityRef.uid must be present (uid pinning)")
	}

	// CompensatingControls and CoveredFailureModes are required non-empty.
	if len(got.CompensatingControls) == 0 {
		t.Fatal("compensatingControls must be non-empty")
	}
	for i, c := range got.CompensatingControls {
		if c == "" {
			t.Fatalf("compensatingControls[%d] must be non-empty", i)
		}
	}
	if len(got.CoveredFailureModes) == 0 {
		t.Fatal("coveredFailureModes must be non-empty")
	}
	for i, m := range got.CoveredFailureModes {
		if m == "" {
			t.Fatalf("coveredFailureModes[%d] must be non-empty", i)
		}
	}

	if got.Purpose == "" || got.AuditTreatment == "" || got.ReassessmentTrigger == "" {
		t.Fatalf("required string fields missing: %+v", got)
	}
	if got.EffectiveFrom == "" || got.EffectiveUntil == "" {
		t.Fatalf("effective interval missing: %+v", got)
	}

	// No ControlRef / approval-workflow fields invented.
	raw := string(data)
	for _, banned := range []string{
		`"controlRef"`,
		`"approvalWorkflow"`,
		`"issuance"`,
		`"revocation"`,
		`"exceptionService"`,
	} {
		if strings.Contains(raw, banned) {
			t.Fatalf("approval-workflow field %s must not be invented; json=%s", banned, raw)
		}
	}
}

func TestSecurityExceptionRefExceptionScopeDoesNotOverrideMetadataScope(t *testing.T) {
	t.Parallel()

	in := sampleDecisionProfile()
	exc := sampleSecurityExceptionRef()
	in.Spec.FailureBehavior.FailPosture = FailPostureFailOpen
	in.Spec.FailureBehavior.SecurityExceptionRef = &exc

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out DecisionProfile
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if out.Metadata.ScopeRef == nil || out.Metadata.ScopeRef.UID != "org-uid-1" {
		t.Fatalf("metadata.scopeRef sole authority lost: %#v", out.Metadata.ScopeRef)
	}
	got := out.Spec.FailureBehavior.SecurityExceptionRef
	if got == nil || got.ExceptionScopeRef == nil {
		t.Fatal("exceptionScopeRef should be present as evidence only")
	}
	// Evidence may coincide with profile scope for compatibility checks, but
	// it remains a nested failureBehavior field — never a second authority.
	if got.ExceptionScopeRef.UID != "project-123" {
		t.Fatalf("exceptionScopeRef evidence uid = %q", got.ExceptionScopeRef.UID)
	}
	raw := string(data)
	if !strings.Contains(raw, `"scopeRef"`) {
		t.Fatal("metadata.scopeRef must serialize")
	}
	if !strings.Contains(raw, `"exceptionScopeRef"`) {
		t.Fatal("exceptionScopeRef must serialize as evidence-only field")
	}
	// Ensure no invented second top-level scope alias.
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		t.Fatalf("unmarshal top: %v", err)
	}
	if _, ok := top["scopeRef"]; ok {
		t.Fatal("top-level scopeRef alias is prohibited; sole authority is metadata.scopeRef")
	}
	if _, ok := top["auditScope"]; ok {
		t.Fatal("AuditScope / auditScope must not be introduced")
	}
}

func TestSecurityExceptionRefOmitemptyExceptionScope(t *testing.T) {
	t.Parallel()

	exc := sampleSecurityExceptionRef()
	exc.ExceptionScopeRef = nil // absent = Platform where allowed (evidence only)

	data, err := json.Marshal(exc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	raw := string(data)
	if strings.Contains(raw, `"exceptionScopeRef"`) {
		t.Fatalf("exceptionScopeRef must be omitted when unset; json=%s", raw)
	}

	// Required SecurityExceptionRef carriers must remain present.
	for _, key := range []string{
		`"approvalRef"`,
		`"ownerRef"`,
		`"approvingAuthorityRef"`,
		`"compensatingControls"`,
		`"effectiveFrom"`,
		`"effectiveUntil"`,
		`"auditTreatment"`,
		`"reassessmentTrigger"`,
		`"coveredFailureModes"`,
		`"purpose"`,
	} {
		if !strings.Contains(raw, key) {
			t.Fatalf("required field %s missing; json=%s", key, raw)
		}
	}

	// UID pinning present on typed refs.
	if !strings.Contains(raw, `"uid":"exception-approval-uid-1"`) {
		t.Fatalf("approvalRef.uid missing; json=%s", raw)
	}
	if !strings.Contains(raw, `"uid":"org-uid-1"`) {
		t.Fatalf("ownerRef.uid missing; json=%s", raw)
	}
	if !strings.Contains(raw, `"uid":"ou-uid-1"`) {
		t.Fatalf("approvingAuthorityRef.uid missing; json=%s", raw)
	}
}

func TestProjectionRuleRoundTrip(t *testing.T) {
	t.Parallel()

	in := ProjectionRule{
		Audience: "auditor",
		Includes: []string{"/record/result/rationale"},
		Excludes: []string{"/record/evaluationResults"},
	}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out ProjectionRule
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Audience != "auditor" || len(out.Includes) != 1 || len(out.Excludes) != 1 {
		t.Fatalf("projection rule = %+v", out)
	}

	minimal := ProjectionRule{Audience: "operator"}
	data, err = json.Marshal(minimal)
	if err != nil {
		t.Fatalf("marshal minimal: %v", err)
	}
	raw := string(data)
	if strings.Contains(raw, `"includes"`) || strings.Contains(raw, `"excludes"`) {
		t.Fatalf("includes/excludes must omit when empty; json=%s", raw)
	}
}

func TestFailPostureEnumeration(t *testing.T) {
	t.Parallel()

	want := []FailPosture{
		FailPostureFailClosed,
		FailPostureRequiresApproval,
		FailPostureFailOpen,
	}
	got := AllFailPostures()
	if len(got) != len(want) {
		t.Fatalf("AllFailPostures len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllFailPostures[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("%q must be Valid()", got[i])
		}
	}
	for _, bad := range []FailPosture{"", "fail-closed", "OPEN", "ALLOWED"} {
		if bad.Valid() {
			t.Fatalf("%q must be rejected", bad)
		}
	}
}

func TestAcceptedEvaluationTypeLimitNarrowingCarrier(t *testing.T) {
	t.Parallel()

	in := sampleDecisionProfile()
	in.Spec.AcceptedEvaluationTypes = []AcceptedEvaluationType{
		{
			Type:    "capacity-evaluator",
			Version: "1.2.0",
			Limits: &ProfileLimits{
				MaxLatencyMs:    500, // narrows profile MaxLatencyMs=2000
				MaxPayloadBytes: 4096,
			},
		},
	}

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out DecisionProfile
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.Spec.AcceptedEvaluationTypes) != 1 {
		t.Fatalf("acceptedEvaluationTypes len=%d", len(out.Spec.AcceptedEvaluationTypes))
	}
	lim := out.Spec.AcceptedEvaluationTypes[0].Limits
	if lim == nil || lim.MaxLatencyMs != 500 || lim.MaxPayloadBytes != 4096 {
		t.Fatalf("evaluator limits = %#v", lim)
	}
	// Profile ceilings remain the outer profile-owned values.
	if out.Spec.Limits.MaxLatencyMs != 2000 {
		t.Fatalf("profile limits must remain unchanged, got %+v", out.Spec.Limits)
	}
}
