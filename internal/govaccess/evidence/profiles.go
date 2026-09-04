package evidence

import (
	"encoding/json"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/decision/bundle"
	decisionvalidate "github.com/sanjeevksaini/sovrunn/internal/decision/validate"
)

// Exact F18-RD-19 profile identities (DD-09). No fourth profile exists.
const (
	ProfileAuthorizationDecision = "authorization-decision"
	ProfileApprovalDecision      = "approval-decision"
	ProfileExceptionDecision     = "exception-decision"
	ProfileVersionV1             = "v1"
)

// ProfileSemantics captures the F18-RD-19 authority, input/result, validity,
// obligation, projection, and audit-linkage semantics for one adopting profile.
type ProfileSemantics struct {
	ID                  string
	Version             string
	Authority           decision.DecisionAuthority
	PrimaryForm         decision.DecisionForm
	InputSummary        string
	ResultSummary       string
	ValiditySummary     string
	ObligationSummary   string
	ProjectionSummary   string
	AuditLinkageSummary string
}

// LocalDecisionProfileBundleBytes returns deterministic strict-local
// DecisionProfileBundle bytes containing exactly the three FEATURE-0013
// adopting profiles. Evidence constructs these bytes but performs no
// registration; root composition is the sole decision/bundle.Load registration
// point (design §FEATURE-0013 Adoption).
func LocalDecisionProfileBundleBytes() []byte {
	doc := bundle.Bundle{
		APIVersion: bundle.APIVersionDecisionProfileBundle,
		Kind:       bundle.KindDecisionProfileBundle,
		Spec: bundle.BundleSpec{
			Profiles: []decision.DecisionProfile{
				f18Profile(ProfileAuthorizationDecision, ProfileVersionV1, authzSemantics()),
				f18Profile(ProfileApprovalDecision, ProfileVersionV1, approvalSemantics()),
				f18Profile(ProfileExceptionDecision, ProfileVersionV1, exceptionSemantics()),
			},
			Forms: []bundle.VersionedEntry{
				{ID: string(decision.DecisionFormAdjudication), Version: "1.0.0"},
			},
			Authorities: []bundle.VersionedEntry{
				{ID: string(decision.DecisionAuthorityAuthoritative), Version: "1.0.0"},
			},
			Obligations: []bundle.VersionedEntry{
				{ID: "audit-retain", Version: "1.0.0"},
			},
			Projections: []bundle.VersionedEntry{
				{ID: "protected-full", Version: "1.0.0"},
				{ID: "safe-denial", Version: "1.0.0"},
			},
			EventTypes: []bundle.VersionedEntry{
				{ID: "authorization.evaluated", Version: "1.0.0"},
				{ID: EventTypeExceptionProposalDecided, Version: "1.0.0"},
				{ID: EventTypeExceptionGrantIssued, Version: "1.0.0"},
			},
			Trust: &decision.TrustCarrier{
				State: "f18-local-trusted",
			},
			ExpectedTrustState: "f18-local-trusted",
		},
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		// Deterministic construction must not fail; surface empty so Load fails closed.
		return nil
	}
	return raw
}

// ValidateLocalDecisionProfileBundle validates the local bytes through existing
// FEATURE-0013 bundle.Load and ValidateProfile. It creates no registry and
// performs no network lookup or fallback version selection.
func ValidateLocalDecisionProfileBundle(data []byte) (bundle.BundleView, *apiproblem.Problem) {
	loaded, prob := bundle.Load(data, apivalid.ModeReadRepresentation)
	if prob != nil {
		return bundle.BundleView{}, prob
	}
	view := loaded.View()
	required := []struct{ id, version string }{
		{ProfileAuthorizationDecision, ProfileVersionV1},
		{ProfileApprovalDecision, ProfileVersionV1},
		{ProfileExceptionDecision, ProfileVersionV1},
	}
	for _, r := range required {
		p, ok := view.Profile(r.id, r.version)
		if !ok {
			return bundle.BundleView{}, decisionvalidate.ValidateProfile(decisionvalidate.ProfileInput{
				CheckLookup: true,
				Ref:         decision.ProfileRef{Name: r.id, Version: r.version},
				Found:       false,
				IDExists:    false,
			})
		}
		if prob := decisionvalidate.ValidateProfile(decisionvalidate.ProfileInput{Profile: p}); prob != nil {
			return bundle.BundleView{}, prob
		}
	}
	return view, nil
}

// LookupProfileSemantics returns the exact F18-RD-19 semantics for a profile.
func LookupProfileSemantics(id, version string) (ProfileSemantics, bool) {
	if version != ProfileVersionV1 {
		return ProfileSemantics{}, false
	}
	switch id {
	case ProfileAuthorizationDecision:
		return authzSemantics(), true
	case ProfileApprovalDecision:
		return approvalSemantics(), true
	case ProfileExceptionDecision:
		return exceptionSemantics(), true
	default:
		return ProfileSemantics{}, false
	}
}

// ExactAdoptingProfileIDs returns the closed three-profile set in stable order.
func ExactAdoptingProfileIDs() []string {
	return []string{
		ProfileAuthorizationDecision,
		ProfileApprovalDecision,
		ProfileExceptionDecision,
	}
}

func authzSemantics() ProfileSemantics {
	return ProfileSemantics{
		ID:                  ProfileAuthorizationDecision,
		Version:             ProfileVersionV1,
		Authority:           decision.DecisionAuthorityAuthoritative,
		PrimaryForm:         decision.DecisionFormAdjudication,
		InputSummary:        "AuthorizationInput with contributing resource/version refs, scope relationship, and guardrail evidence",
		ResultSummary:       "Allow | Deny",
		ValiditySummary:     "evaluation instant only; never bearer authority",
		ObligationSummary:   "protected audit-retain obligation before release",
		ProjectionSummary:   "protected-full and safe-denial projections",
		AuditLinkageSummary: "authorization.evaluated AuditEvent linkage required",
	}
}

func approvalSemantics() ProfileSemantics {
	return ProfileSemantics{
		ID:                  ProfileApprovalDecision,
		Version:             ProfileVersionV1,
		Authority:           decision.DecisionAuthorityAuthoritative,
		PrimaryForm:         decision.DecisionFormAdjudication,
		InputSummary:        "policy version, subject/proposal digest, active stage, eligible PrincipalRefs, counted decisions",
		ResultSummary:       "Approved | Denied",
		ValiditySummary:     "ends no later than ApprovalRequest expiresAt",
		ObligationSummary:   "mandatory approval-decision/v1 DecisionRecord on terminal Approved|Denied",
		ProjectionSummary:   "protected-full projection; no fabricated Cancelled/Expired decisions",
		AuditLinkageSummary: "terminal approval AuditEvent and DecisionRecord linkage",
	}
}

func exceptionSemantics() ProfileSemantics {
	return ProfileSemantics{
		ID:                  ProfileExceptionDecision,
		Version:             ProfileVersionV1,
		Authority:           decision.DecisionAuthorityAuthoritative,
		PrimaryForm:         decision.DecisionFormAdjudication,
		InputSummary:        "control version, subject, scope, interval, typed override, compensating controls, approval evidence",
		ResultSummary:       "Grant | Deny",
		ValiditySummary:     "Deny is decision-instant; Grant bounded by approved exception interval; never bearer authority",
		ObligationSummary:   "mandatory exception-decision/v1 DecisionRecord on terminal Grant|Deny",
		ProjectionSummary:   "protected-full projection; Deny creates no ExceptionGrant",
		AuditLinkageSummary: "exceptionproposal.decided (and exceptiongrant.issued on Grant) linkage",
	}
}

func f18Profile(id, version string, sem ProfileSemantics) decision.DecisionProfile {
	return decision.DecisionProfile{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: decision.APIVersionDecisionProfile,
			Kind:       decision.KindDecisionProfile,
		},
		Metadata: apimeta.ObjectMeta{
			Name: id,
			UID:  "f18-profile-" + id + "-" + version,
		},
		Spec: decision.ProfileSpec{
			Name:                   id,
			Family:                 "feature-0018",
			Version:                version,
			Owner:                  "feature-0018-governance",
			LifecycleStatus:        decisionvalidate.LifecycleStatusActive,
			PrimaryForm:            sem.PrimaryForm,
			AllowedAuthorityLevels: []decision.DecisionAuthority{sem.Authority},
			InputSchemaRef:         "schemas/f18/" + id + "/input@" + version,
			TypedResultSchemaRef:   "schemas/f18/" + id + "/result@" + version,
			RationaleSchemaRef:     "schemas/f18/" + id + "/rationale@" + version,
			ObligationSchemaRef:    "schemas/f18/" + id + "/obligation@" + version,
			ScopeSemantics: decision.ProfileScopeSemantics{
				AllowedScopes: []apimeta.ScopeKind{
					apimeta.ScopeOrganization,
					apimeta.ScopeOrganizationUnit,
					apimeta.ScopeTenant,
					apimeta.ScopeProject,
					apimeta.ScopeCloudPlatform,
					apimeta.ScopeCloudProvider,
				},
				PlatformAllowed: true,
			},
			ActorSemantics:   decision.ProfileActorSemantics{RequireActor: true},
			SubjectSemantics: decision.ProfileSubjectSemantics{RequireSubject: true},
			EvidenceSemantics: decision.ProfileEvidenceSemantics{
				RequireEvidence: true,
			},
			AuditSemantics: decision.ProfileAuditSemantics{
				RequireDurableRecord: true,
			},
			DeterminismBehavior: decision.DeterminismBehavior{
				RequireDeterministic: true,
			},
			FailureBehavior: decision.FailureBehavior{
				FailPosture:          decision.FailPostureFailClosed,
				TimeoutBehavior:      "fail-closed",
				InsufficientEvidence: "fail-closed",
			},
			ValidityRules: decision.ValidityRules{
				AllowCorrection:   false,
				AllowSupersession: true,
				AllowRevocation:   false,
			},
			SensitivityCeiling:  decision.SensitivityInternal,
			IdempotencyIdentity: id + "/" + version,
			Limits: decision.ProfileLimits{
				MaxPayloadBytes: 65536,
				MaxDepth:        8,
				MaxFieldCount:   64,
			},
			ProjectionRules: []decision.ProjectionRule{
				{Audience: "protected-full", Includes: []string{"/"}},
				{Audience: "safe-denial", Includes: []string{"/record/result"}},
			},
			CompatibilityPolicy: decision.CompatibilityPolicy{
				Owner:              "feature-0018-governance",
				BackwardCompatible: false,
			},
		},
	}
}
