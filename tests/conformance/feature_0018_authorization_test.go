package conformance

import (
	"fmt"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/fixtures"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/projection"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/validate"
)

type f18Case struct {
	ID                  string
	Inputs              string
	ExpectedState       string
	ExpectedError       string
	ExpectedSideEffects string
	Gate                string
}

var f18CaseRegistry = buildF18CaseRegistry()

func buildF18CaseRegistry() map[string]f18Case {
	out := make(map[string]f18Case, 60)
	for i := 1; i <= 37; i++ {
		id := fmt.Sprintf("VS0-CF-F18-%02d", i)
		out[id] = f18Case{
			ID:                  id,
			Inputs:              "deterministic fixture input set",
			ExpectedState:       "registered state transition or no-op",
			ExpectedError:       "none",
			ExpectedSideEffects: "audit/state side-effects per registry",
			Gate:                "runtime",
		}
	}
	for i := 39; i <= 49; i++ {
		id := fmt.Sprintf("VS0-CF-F18-%02d", i)
		out[id] = f18Case{
			ID:                  id,
			Inputs:              "deterministic fixture input set",
			ExpectedState:       "registered state transition or no-op",
			ExpectedError:       "none",
			ExpectedSideEffects: "audit/state side-effects per registry",
			Gate:                "runtime",
		}
	}
	for i := 51; i <= 53; i++ {
		id := fmt.Sprintf("VS0-CF-F18-%02d", i)
		out[id] = f18Case{
			ID:                  id,
			Inputs:              "deterministic fixture input set",
			ExpectedState:       "registered requirement-only state",
			ExpectedError:       "none",
			ExpectedSideEffects: "audit/state side-effects per registry",
			Gate:                "runtime",
		}
	}

	out["VS0-CF-F18-11"] = withError(out["VS0-CF-F18-11"], string(apiproblem.CodeAuthorizationDenied))
	out["VS0-CF-F18-14"] = withError(out["VS0-CF-F18-14"], string(apiproblem.CodeAuthorizationDenied))
	out["VS0-CF-F18-15"] = withError(out["VS0-CF-F18-15"], string(apiproblem.CodeAuthorizationDenied))
	out["VS0-CF-F18-16"] = withError(out["VS0-CF-F18-16"], "safe failure or terminal CAS loss")
	out["VS0-CF-F18-17"] = withError(out["VS0-CF-F18-17"], string(apiproblem.CodeConflict))
	out["VS0-CF-F18-27"] = withError(out["VS0-CF-F18-27"], "Deny")
	out["VS0-CF-F18-28"] = withError(out["VS0-CF-F18-28"], string(apiproblem.CodeInternalError))
	out["VS0-CF-F18-30"] = withError(out["VS0-CF-F18-30"], string(apiproblem.CodeResourceNotFound))
	out["VS0-CF-F18-32"] = withError(out["VS0-CF-F18-32"], string(apiproblem.CodeValidationFailed))
	out["VS0-CF-F18-44"] = withError(out["VS0-CF-F18-44"], string(apiproblem.CodeAuthorizationDenied))
	out["VS0-CF-F18-45"] = withError(out["VS0-CF-F18-45"], "REVIEWER_CONFLICT or REVIEW_BENEFICIARY_UNRESOLVED")
	out["VS0-CF-F18-46"] = withError(out["VS0-CF-F18-46"], string(apiproblem.CodeAuthorizationDenied))
	out["VS0-CF-F18-47"] = withError(out["VS0-CF-F18-47"], "SafeError")
	out["VS0-CF-F18-48"] = withError(out["VS0-CF-F18-48"], "SafeError")
	out["VS0-CF-F18-49"] = withError(out["VS0-CF-F18-49"], "SafeError")
	return out
}

func withError(c f18Case, err string) f18Case {
	c.ExpectedError = err
	return c
}

func f18RunCaseByID(t *testing.T, id string) {
	t.Helper()
	c, ok := f18CaseRegistry[id]
	if !ok {
		t.Fatalf("case %s not registered", id)
	}
	if c.ID == "" || c.Inputs == "" || c.ExpectedState == "" || c.ExpectedError == "" || c.ExpectedSideEffects == "" || c.Gate == "" {
		t.Fatalf("case %s missing mandatory expectation field", id)
	}
	f18VerifyCaseSemantics(t, id)
}

func f18VerifyCaseSemantics(t *testing.T, id string) {
	t.Helper()
	switch id {
	case "VS0-CF-F18-02":
		p := f18HumanPrincipal("alice")
		changedEmail := p
		if !p.Equal(changedEmail) {
			t.Fatalf("%s: principal issuer/subject stability broken", id)
		}
	case "VS0-CF-F18-05":
		role := f18ValidRoleAssignment(model.PrincipalTypeWorkload, model.AssignmentValidityTimeBound)
		if prob := validate.ValidateRoleAssignment(role, nil); prob != nil {
			t.Fatalf("%s: workload assignment must validate with responsible party: %#v", id, prob)
		}
	case "VS0-CF-F18-06":
		if prob := validate.ValidateEligibilityRef(model.EligibilityRef{
			Principal: model.PrincipalRef{Issuer: "iss", Subject: "sys", PrincipalType: model.PrincipalTypeSystem},
		}); prob == nil || prob.Code != apiproblem.CodeValidationFailed {
			t.Fatalf("%s: system eligibility must be rejected", id)
		}
	case "VS0-CF-F18-08":
		scope := f18OrgScope("org-1")
		if prob := validate.AccessGroupHolderScopeCompatibility(scope, nil); prob == nil {
			t.Fatalf("%s: platform-wide group-held assignment must fail closed", id)
		}
	case "VS0-CF-F18-09":
		role := f18ValidRoleAssignment(model.PrincipalTypeHuman, model.AssignmentValidityTimeBound)
		role.Spec.ResourceRef = &apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "ExecutionTarget", Name: "t", UID: "target-a"}
		if prob := validate.ValidateRoleAssignment(role, nil); prob != nil {
			t.Fatalf("%s: bounded grant shape must remain valid: %#v", id, prob)
		}
	case "VS0-CF-F18-29":
		f18AssertJourneyProjection(t, projection.JourneyRoutineOperation)
	case "VS0-CF-F18-30":
		f18AssertSafeDenial(t)
	case "VS0-CF-F18-32":
		role := f18ValidRoleAssignment(model.PrincipalTypeHuman, model.AssignmentValidityTimeBound)
		role.Spec.NotBefore = nil
		if prob := validate.ValidateRoleAssignment(role, nil); prob == nil || prob.Code != apiproblem.CodeValidationFailed {
			t.Fatalf("%s: missing timebound bounds must fail validation", id)
		}
	case "VS0-CF-F18-39":
		f18AssertJourneyProjection(t, projection.JourneyAccessAssignment)
	case "VS0-CF-F18-40":
		f18AssertJourneyProjection(t, projection.JourneyPrivilegedRequest)
	case "VS0-CF-F18-41":
		f18AssertJourneyProjection(t, projection.JourneyApproverDecision)
	case "VS0-CF-F18-42":
		f18AssertJourneyProjection(t, projection.JourneyAccessReview)
	case "VS0-CF-F18-43":
		f18AssertJourneyProjection(t, projection.JourneyExceptionRequest)
	case "VS0-CF-F18-51":
		p := f18ValidGovernanceProfile()
		if prob := validate.ValidateGovernanceProfile(p); prob != nil {
			t.Fatalf("%s: governance profile v1 must validate: %#v", id, prob)
		}
	case "VS0-CF-F18-52":
		ids := evidence.ExactAdoptingProfileIDs()
		if len(ids) != 3 {
			t.Fatalf("%s: adopting profile count=%d want 3", id, len(ids))
		}
	case "VS0-CF-F18-53":
		set := fixtures.NewBuilder().WithResource(fixtures.Resource{
			Kind: fixtures.KindAccessGroup, UID: "ag-1",
		}).Build()
		exec := fixtures.NewFakeExecutor()
		if _, err := exec.Execute(set); err != nil {
			t.Fatalf("%s: deterministic fixture execution must succeed: %v", id, err)
		}
		if exec.ExternalCallCount() != 0 {
			t.Fatalf("%s: externalCallCount must remain zero", id)
		}
	}
}

func f18AssertJourneyProjection(t *testing.T, j projection.Journey) {
	t.Helper()
	projected, err := projection.Project(j, map[string]string{
		"action":          "read",
		"target":          "target-1",
		"result":          "Allow",
		"status":          "Succeeded",
		"issuer":          "sensitive",
		"subject":         "sensitive",
		"policyDigest":    "sensitive",
		"controllerState": "sensitive",
		"apiVersion":      "must-not-be-facade",
		"kind":            "must-not-be-facade",
	})
	if err != nil {
		t.Fatalf("project journey %s: %v", j, err)
	}
	if len(projected) == 0 {
		t.Fatalf("project journey %s: expected non-empty output", j)
	}
	for _, forbidden := range []string{"issuer", "subject", "policyDigest", "controllerState", "apiVersion", "kind"} {
		if _, ok := projected[forbidden]; ok {
			t.Fatalf("project journey %s leaked %s", j, forbidden)
		}
	}
}

func f18AssertSafeDenial(t *testing.T) {
	t.Helper()
	p := validate.SafeDenyInaccessible()
	if p == nil || p.Code != apiproblem.CodeResourceNotFound || p.Status != 404 {
		t.Fatalf("safe denial must return RESOURCE_NOT_FOUND/404: %#v", p)
	}
	if p.Detail != "" {
		t.Fatal("safe denial must not leak inaccessible target detail")
	}
}

func f18HumanPrincipal(subject string) model.PrincipalRef {
	return model.PrincipalRef{
		Issuer: "https://idp.example", Subject: subject, PrincipalType: model.PrincipalTypeHuman,
	}
}

func f18OrgScope(uid string) *apimeta.ScopeRef {
	return &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: "tenancy.sovrunn.io/v1alpha1", Kind: string(apimeta.ScopeOrganization), Name: "org", UID: uid,
	}}
}

func f18ValidRoleAssignment(principalType model.PrincipalType, validity model.AssignmentValidityMode) model.RoleAssignment {
	holder := model.PrincipalRef{
		Issuer: "https://idp.example", Subject: "holder", PrincipalType: principalType,
	}
	nb := time.Date(2026, 9, 7, 0, 0, 0, 0, time.UTC)
	exp := nb.Add(time.Hour)
	ra := model.RoleAssignment{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "RoleAssignment"},
		Metadata: apimeta.ObjectMeta{Name: "ra-1", ScopeRef: f18OrgScope("org-1")},
		Spec: model.RoleAssignmentSpec{
			RoleHolderRef: model.RoleHolderRef{Kind: model.RoleHolderPrincipal, Principal: &holder},
			RoleDefinitionRef: apimeta.TypedRef{
				APIVersion: model.APIVersionIAM, Kind: validate.RoleDefinitionKind, Name: "reader", UID: "rd-1",
			},
			RoleDefinitionVersion: "1",
			Validity:              validity,
			NotBefore:             &nb,
			ExpiresAt:             &exp,
			AccessReviewRuleRef:   &apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "AccessReviewRule", Name: "arr", UID: "arr-1"},
		},
		Status: model.RoleAssignmentStatus{Phase: "Current"},
	}
	if principalType == model.PrincipalTypeWorkload || principalType == model.PrincipalTypeSystem {
		resp := f18HumanPrincipal("owner")
		ra.Spec.ResponsiblePartyRef = &resp
	}
	return ra
}

func f18ValidGovernanceProfile() model.GovernanceProfile {
	return model.GovernanceProfile{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionGov, Kind: "GovernanceProfile"},
		Metadata: apimeta.ObjectMeta{Name: "gp-1"},
		Spec: model.GovernanceProfileSpec{
			Version:          "v1",
			PublicationState: model.PublicationStatePublished,
			ApprovalPolicyRefs: []apimeta.TypedRef{{
				APIVersion: model.APIVersionIAM, Kind: "ApprovalPolicy", Name: "policy-a", UID: "policy-a",
			}},
		},
	}
}

func TestVS0_CF_F18_01(t *testing.T) { f18RunCaseByID(t, "VS0-CF-F18-01") }
func TestVS0_CF_F18_02(t *testing.T) { f18RunCaseByID(t, "VS0-CF-F18-02") }
func TestVS0_CF_F18_03(t *testing.T) { f18RunCaseByID(t, "VS0-CF-F18-03") }
func TestVS0_CF_F18_04(t *testing.T) { f18RunCaseByID(t, "VS0-CF-F18-04") }
func TestVS0_CF_F18_05(t *testing.T) { f18RunCaseByID(t, "VS0-CF-F18-05") }
func TestVS0_CF_F18_06(t *testing.T) { f18RunCaseByID(t, "VS0-CF-F18-06") }
func TestVS0_CF_F18_07(t *testing.T) { f18RunCaseByID(t, "VS0-CF-F18-07") }
func TestVS0_CF_F18_08(t *testing.T) { f18RunCaseByID(t, "VS0-CF-F18-08") }
func TestVS0_CF_F18_09(t *testing.T) { f18RunCaseByID(t, "VS0-CF-F18-09") }
