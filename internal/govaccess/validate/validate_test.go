package validate_test

import (
	"testing"
	"testing/quick"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/validate"
)

func humanPrincipal() model.PrincipalRef {
	return model.PrincipalRef{
		Issuer: "https://idp.example", Subject: "user-1", PrincipalType: model.PrincipalTypeHuman,
	}
}

func orgScope(uid string) *apimeta.ScopeRef {
	return &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: "tenancy.sovrunn.io/v1alpha1", Kind: string(apimeta.ScopeOrganization), Name: "org", UID: uid,
	}}
}

func TestValidatePrincipalRefHappyAndNegative(t *testing.T) {
	t.Parallel()
	if prob := validate.ValidatePrincipalRef(humanPrincipal()); prob != nil {
		t.Fatalf("happy path: %v", prob)
	}
	if prob := validate.ValidatePrincipalRef(model.PrincipalRef{}); prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("expected VALIDATION_FAILED, got %#v", prob)
	}
	if prob := validate.ValidatePrincipalRef(model.PrincipalRef{
		Issuer: "iss", Subject: "sub", PrincipalType: "Email",
	}); prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("expected invalid enum, got %#v", prob)
	}
}

func TestValidateEligibilityRefHumanOnly(t *testing.T) {
	t.Parallel()
	if prob := validate.ValidateEligibilityRef(model.EligibilityRef{Principal: humanPrincipal()}); prob != nil {
		t.Fatalf("human eligibility: %v", prob)
	}
	if prob := validate.ValidateEligibilityRef(model.EligibilityRef{Principal: model.PrincipalRef{
		Issuer: "iss", Subject: "wl", PrincipalType: model.PrincipalTypeWorkload,
	}}); prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("workload eligibility must fail, got %#v", prob)
	}
}

func TestValidateRoleHolderRefUnion(t *testing.T) {
	t.Parallel()
	p := humanPrincipal()
	if prob := validate.ValidateRoleHolderRef(model.RoleHolderRef{
		Kind: model.RoleHolderPrincipal, Principal: &p,
	}); prob != nil {
		t.Fatalf("principal holder: %v", prob)
	}
	g := model.AccessGroupRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionIAM, Kind: model.KindAccessGroup, Name: "g", UID: "uid-g",
	}}
	if prob := validate.ValidateRoleHolderRef(model.RoleHolderRef{
		Kind: model.RoleHolderAccessGroup, AccessGroup: &g,
	}); prob != nil {
		t.Fatalf("group holder: %v", prob)
	}
	if prob := validate.ValidateRoleHolderRef(model.RoleHolderRef{
		Kind: model.RoleHolderPrincipal, Principal: &p, AccessGroup: &g,
	}); prob == nil {
		t.Fatal("mixed holder must fail")
	}
}

func TestValidateAccessGroupHappyAndScope(t *testing.T) {
	t.Parallel()
	ok := model.AccessGroup{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: model.KindAccessGroup},
		Metadata: apimeta.ObjectMeta{Name: "ops", ScopeRef: orgScope("org-1")},
		Spec:     model.AccessGroupSpec{OwnerRef: humanPrincipal()},
		Status:   model.AccessGroupStatus{Phase: model.AccessGroupPhaseActive},
	}
	if prob := validate.ValidateAccessGroup(ok); prob != nil {
		t.Fatalf("happy: %v", prob)
	}
	badScope := ok
	badScope.Metadata.ScopeRef = &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		Kind: string(apimeta.ScopeProject), Name: "p", UID: "p1",
	}}
	if prob := validate.ValidateAccessGroup(badScope); prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("project scope must fail for AccessGroup, got %#v", prob)
	}
}

func TestValidateMembershipHappyAndGuestRules(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	ok := model.Membership{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: validate.MembershipKind},
		Metadata: apimeta.ObjectMeta{Name: "m1", ScopeRef: orgScope("org-1")},
		Spec: model.MembershipSpec{
			PrincipalRef: humanPrincipal(), ContextKind: model.MembershipContextOrganization,
			MembershipType: model.MembershipTypeStandard,
		},
		Protected: model.MembershipProtected{SourceAuthority: "trusted-provisioner", ProvisionedAt: now},
		Status:    model.MembershipStatus{Phase: model.MembershipPhaseActive},
	}
	if prob := validate.ValidateMembership(ok); prob != nil {
		t.Fatalf("happy: %v", prob)
	}
	guest := ok
	exp := now.Add(24 * time.Hour)
	sponsor := humanPrincipal()
	guest.Spec.MembershipType = model.MembershipTypeGuest
	guest.Spec.ExpiresAt = &exp
	guest.Spec.SponsorRef = &sponsor
	if prob := validate.ValidateMembership(guest); prob != nil {
		t.Fatalf("guest happy: %v", prob)
	}
	guestBad := guest
	guestBad.Spec.SponsorRef = nil
	if prob := validate.ValidateMembership(guestBad); prob == nil {
		t.Fatal("guest without sponsor must fail")
	}
}

func TestValidateRoleDefinitionAndAssignment(t *testing.T) {
	t.Parallel()
	rd := model.RoleDefinition{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: validate.RoleDefinitionKind},
		Metadata: apimeta.ObjectMeta{Name: "reader", ScopeRef: nil}, // Platform
		Spec: model.RoleDefinitionSpec{
			Version: "1", Actions: []string{"executiontarget.read"},
			PublicationState:   model.PublicationStatePublished,
			BaseClassification: model.RoleClassificationOrdinary,
		},
	}
	if prob := validate.ValidateRoleDefinition(rd); prob != nil {
		t.Fatalf("role definition: %v", prob)
	}
	p := humanPrincipal()
	nb := time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC)
	exp := nb.Add(time.Hour)
	ra := model.RoleAssignment{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "RoleAssignment"},
		Metadata: apimeta.ObjectMeta{Name: "a1", ScopeRef: orgScope("org-1")},
		Spec: model.RoleAssignmentSpec{
			RoleHolderRef:         model.RoleHolderRef{Kind: model.RoleHolderPrincipal, Principal: &p},
			RoleDefinitionRef:     apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: validate.RoleDefinitionKind, Name: "reader", UID: "rd-1"},
			RoleDefinitionVersion: "1",
			Validity:              model.AssignmentValidityTimeBound,
			NotBefore:             &nb, ExpiresAt: &exp,
		},
		Status: model.RoleAssignmentStatus{Phase: "Current"},
	}
	if prob := validate.ValidateRoleAssignment(ra, nil); prob != nil {
		t.Fatalf("role assignment: %v", prob)
	}
}

func TestAccessGroupHolderScopeCompatibility(t *testing.T) {
	t.Parallel()
	groupOrg := orgScope("org-1")
	assignOrg := orgScope("org-1")
	if prob := validate.AccessGroupHolderScopeCompatibility(groupOrg, assignOrg); prob != nil {
		t.Fatalf("org compatible: %v", prob)
	}
	assignProject := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		Kind: string(apimeta.ScopeProject), Name: "p", UID: "p1",
	}}
	if prob := validate.AccessGroupHolderScopeCompatibility(groupOrg, assignProject); prob != nil {
		t.Fatalf("org descendant project should be compatible: %v", prob)
	}
	assignPlatform := (*apimeta.ScopeRef)(nil)
	if prob := validate.AccessGroupHolderScopeCompatibility(groupOrg, assignPlatform); prob == nil {
		t.Fatal("Platform AccessGroup-held assignment must fail")
	}
	cpGroup := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		Kind: string(apimeta.ScopeCloudProvider), Name: "cp", UID: "cp-1",
	}}
	cpAssign := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		Kind: string(apimeta.ScopeCloudProvider), Name: "cp", UID: "cp-1",
	}}
	if prob := validate.AccessGroupHolderScopeCompatibility(cpGroup, cpAssign); prob != nil {
		t.Fatalf("cloud provider exact match: %v", prob)
	}
	cpOther := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		Kind: string(apimeta.ScopeCloudProvider), Name: "cp", UID: "cp-2",
	}}
	if prob := validate.AccessGroupHolderScopeCompatibility(cpGroup, cpOther); prob == nil {
		t.Fatal("mismatched CloudProvider UID must fail")
	}
}

func TestSafeDenialInheritedCodes(t *testing.T) {
	t.Parallel()
	prob := validate.SafeDenyInaccessible()
	if prob == nil || prob.Code != apiproblem.CodeResourceNotFound || prob.Status != 404 {
		t.Fatalf("safe denial must be RESOURCE_NOT_FOUND/404, got %#v", prob)
	}
	if validate.ResolveReferenceVisibility(true, true) != nil {
		t.Fatal("authorized visible ref must succeed")
	}
	if p := validate.ResolveReferenceVisibility(false, true); p == nil || p.Code != apiproblem.CodeResourceNotFound {
		t.Fatalf("missing must safe-deny, got %#v", p)
	}
	if p := validate.ResolveReferenceVisibility(true, false); p == nil || p.Code != apiproblem.CodeResourceNotFound {
		t.Fatalf("unauthorized must safe-deny, got %#v", p)
	}
}

func TestRemainingValidatorsHappyPaths(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	elig := model.EligibilityRef{Principal: humanPrincipal()}

	par := model.PrivilegedAccessRequest{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "PrivilegedAccessRequest"},
		Metadata: apimeta.ObjectMeta{Name: "jit-1", ScopeRef: orgScope("org-1")},
		Spec: model.PrivilegedAccessRequestSpec{
			Mode: model.PrivilegedAccessModeJIT, RequesterRef: humanPrincipal(),
			RoleDefinitionRef:  apimeta.TypedRef{Kind: validate.RoleDefinitionKind, Name: "admin", UID: "rd"},
			ActivationDeadline: now.Add(time.Hour), RequestedDuration: "30m",
		},
	}
	if prob := validate.ValidatePrivilegedAccessRequest(par); prob != nil {
		t.Fatalf("privileged: %v", prob)
	}

	ar := model.AccessReview{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "AccessReview"},
		Metadata: apimeta.ObjectMeta{Name: "rev-1", ScopeRef: orgScope("org-1")},
		Spec: model.AccessReviewSpec{
			AccessReviewRuleRef: apimeta.TypedRef{Name: "rule", UID: "rule-1"},
			ReviewerEligibility: []model.EligibilityRef{elig},
		},
	}
	if prob := validate.ValidateAccessReview(ar); prob != nil {
		t.Fatalf("access review: %v", prob)
	}

	ap := model.ApprovalPolicy{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "ApprovalPolicy"},
		Metadata: apimeta.ObjectMeta{Name: "pol", ScopeRef: nil},
		Spec: model.ApprovalPolicySpec{
			Version: "1", PublicationState: model.PublicationStatePublished,
			ApproverEligibility: []model.EligibilityRef{elig},
		},
	}
	if prob := validate.ValidateApprovalPolicy(ap); prob != nil {
		t.Fatalf("approval policy: %v", prob)
	}

	aq := model.ApprovalRequest{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "ApprovalRequest"},
		Metadata: apimeta.ObjectMeta{Name: "req", ScopeRef: orgScope("org-1")},
		Spec: model.ApprovalRequestSpec{
			ApprovalPolicyRef: apimeta.TypedRef{Kind: "ApprovalPolicy", Name: "pol", UID: "pol-1"},
			SubjectKind:       "PrivilegedAccessRequest", SubjectUID: "jit-1",
		},
	}
	if prob := validate.ValidateApprovalRequest(aq); prob != nil {
		t.Fatalf("approval request: %v", prob)
	}

	eg := model.ExceptionGrant{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionGov, Kind: model.KindExceptionGrant},
		Metadata: apimeta.ObjectMeta{Name: "ex", ScopeRef: orgScope("org-1")},
		Record: model.ExceptionGrantRecord{
			Effect:     model.ExceptionGrantEffectGrant,
			ControlRef: apimeta.TypedRef{Name: "ctl", UID: "ctl-1"},
			SubjectUID: "sub-1", NotBefore: now, ExpiresAt: now.Add(24 * time.Hour),
		},
	}
	if prob := validate.ValidateExceptionGrant(eg); prob != nil {
		t.Fatalf("exception grant: %v", prob)
	}

	gp := model.GovernanceProfile{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionGov, Kind: "GovernanceProfile"},
		Metadata: apimeta.ObjectMeta{Name: "gp", ScopeRef: nil},
		Spec: model.GovernanceProfileSpec{
			Version: "1", PublicationState: model.PublicationStatePublished,
			ApprovalPolicyRefs: []apimeta.TypedRef{{Kind: "ApprovalPolicy", Name: "pol", UID: "pol-1"}},
		},
	}
	if prob := validate.ValidateGovernanceProfile(gp); prob != nil {
		t.Fatalf("governance profile: %v", prob)
	}
}

func TestPropertyValidationDeterminism(t *testing.T) {
	t.Parallel()
	fn := func(subject string) bool {
		if subject == "" {
			subject = "s"
		}
		ref := model.PrincipalRef{Issuer: "iss", Subject: subject, PrincipalType: model.PrincipalTypeHuman}
		a := validate.ValidatePrincipalRef(ref)
		b := validate.ValidatePrincipalRef(ref)
		if a == nil && b == nil {
			return true
		}
		if a == nil || b == nil {
			return false
		}
		return a.Code == b.Code && a.Status == b.Status
	}
	if err := quick.Check(fn, &quick.Config{MaxCount: 64}); err != nil {
		t.Fatal(err)
	}
}

func TestPropertyCrossTypeSubstitutionFailsValidation(t *testing.T) {
	t.Parallel()
	fn := func(uid string) bool {
		if uid == "" {
			uid = "uid"
		}
		// AccessGroupRef-shaped values must not validate as PrincipalRef.
		fake := model.PrincipalRef{Issuer: "", Subject: uid, PrincipalType: model.PrincipalTypeHuman}
		return validate.ValidatePrincipalRef(fake) != nil
	}
	if err := quick.Check(fn, &quick.Config{MaxCount: 32}); err != nil {
		t.Fatal(err)
	}
}

func TestActionTargetBindingValidation(t *testing.T) {
	t.Parallel()
	ok := model.ActionTargetBinding{
		Mode: model.ActionTargetExactResource, Action: "executiontarget.read", TargetKind: "ExecutionTarget",
	}
	if prob := validate.ValidateActionTargetBinding(ok); prob != nil {
		t.Fatalf("exact: %v", prob)
	}
	mixed := ok
	mixed.ScopeKind = apimeta.ScopeProject
	if prob := validate.ValidateActionTargetBinding(mixed); prob == nil {
		t.Fatal("mixed binding must fail")
	}
}

func TestInheritedProblemCodesOnly(t *testing.T) {
	t.Parallel()
	probs := []*apiproblem.Problem{
		validate.ValidatePrincipalRef(model.PrincipalRef{}),
		validate.SafeDenyInaccessible(),
		validate.ValidateAccessGroup(model.AccessGroup{}),
	}
	for _, p := range probs {
		if p == nil {
			continue
		}
		if !p.Code.Valid() {
			t.Fatalf("non-inherited problem code %q", p.Code)
		}
	}
}
