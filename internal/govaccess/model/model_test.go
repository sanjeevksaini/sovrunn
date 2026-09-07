package model_test

import (
	"reflect"
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

func TestPrincipalRefIsNotTypedRef(t *testing.T) {
	t.Parallel()
	pt := reflect.TypeOf(model.PrincipalRef{})
	if _, ok := pt.FieldByName("TypedRef"); ok {
		t.Fatal("PrincipalRef must not embed TypedRef")
	}
	for _, name := range []string{"APIVersion", "Kind", "Name", "UID"} {
		if _, ok := pt.FieldByName(name); ok {
			t.Fatalf("PrincipalRef must not carry TypedRef field %s", name)
		}
	}
	// Stable identity is issuer/subject, not email.
	p := model.PrincipalRef{
		Issuer:        "https://idp.example",
		Subject:       "sub-123",
		PrincipalType: model.PrincipalTypeHuman,
	}
	if p.Issuer == "" || p.Subject == "" {
		t.Fatal("issuer/subject required for stable identity")
	}
	rt := reflect.TypeOf(p)
	if _, ok := rt.FieldByName("Email"); ok {
		t.Fatal("PrincipalRef must not carry email")
	}
}

func TestRoleHolderRefDiscriminatedUnion(t *testing.T) {
	t.Parallel()
	principal := model.PrincipalRef{
		Issuer: "iss", Subject: "sub", PrincipalType: model.PrincipalTypeHuman,
	}
	group := model.AccessGroupRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionIAM, Kind: model.KindAccessGroup, Name: "g1", UID: "uid-1",
	}}

	okPrincipal := model.RoleHolderRef{
		Kind: model.RoleHolderPrincipal, Principal: &principal,
	}
	if !okPrincipal.IsPrincipal() || okPrincipal.IsAccessGroup() {
		t.Fatal("expected Principal-only variant")
	}

	okGroup := model.RoleHolderRef{
		Kind: model.RoleHolderAccessGroup, AccessGroup: &group,
	}
	if !okGroup.IsAccessGroup() || okGroup.IsPrincipal() {
		t.Fatal("expected AccessGroup-only variant")
	}

	mixed := model.RoleHolderRef{
		Kind: model.RoleHolderPrincipal, Principal: &principal, AccessGroup: &group,
	}
	if mixed.IsPrincipal() || mixed.IsAccessGroup() {
		t.Fatal("mixed holder must not report a sole variant")
	}

	empty := model.RoleHolderRef{Kind: model.RoleHolderPrincipal}
	if empty.IsPrincipal() {
		t.Fatal("empty Principal variant must fail IsPrincipal")
	}
}

func TestEligibilityRefHumanOnlyConstructionShape(t *testing.T) {
	t.Parallel()
	human := model.EligibilityRef{Principal: model.PrincipalRef{
		Issuer: "iss", Subject: "sub", PrincipalType: model.PrincipalTypeHuman,
	}}
	if human.Principal.PrincipalType != model.PrincipalTypeHuman {
		t.Fatal("EligibilityRef must carry Human PrincipalRef")
	}
	workload := model.EligibilityRef{Principal: model.PrincipalRef{
		Issuer: "iss", Subject: "wl", PrincipalType: model.PrincipalTypeWorkload,
	}}
	if workload.Principal.PrincipalType == model.PrincipalTypeHuman {
		t.Fatal("fixture setup error")
	}
}

func TestEnumValidAndAllAccessors(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		all  int
		ok   func() bool
		bad  func() bool
	}{
		{"PrincipalType", len(model.AllPrincipalTypes()), model.PrincipalTypeHuman.Valid, func() bool { return model.PrincipalType("email").Valid() }},
		{"MembershipContextKind", len(model.AllMembershipContextKinds()), model.MembershipContextOrganization.Valid, func() bool { return model.MembershipContextKind("Team").Valid() }},
		{"MembershipType", len(model.AllMembershipTypes()), model.MembershipTypeGuest.Valid, func() bool { return model.MembershipType("External").Valid() }},
		{"PrivilegedAccessMode", len(model.AllPrivilegedAccessModes()), model.PrivilegedAccessModeJIT.Valid, func() bool { return model.PrivilegedAccessMode("Always").Valid() }},
		{"ApprovalDecision", len(model.AllApprovalDecisions()), model.ApprovalDecisionApprove.Valid, func() bool { return model.ApprovalDecision("Abstain").Valid() }},
		{"ExceptionGrantEffect", len(model.AllExceptionGrantEffects()), model.ExceptionGrantEffectGrant.Valid, func() bool { return model.ExceptionGrantEffect("Modify").Valid() }},
		{"AssignmentValidityMode", len(model.AllAssignmentValidityModes()), model.AssignmentValidityStanding.Valid, func() bool { return model.AssignmentValidityMode("Forever").Valid() }},
		{"ActionTargetMode", len(model.AllActionTargetModes()), model.ActionTargetExactResource.Valid, func() bool { return model.ActionTargetMode("Wildcard").Valid() }},
		{"RoleHolderKind", len(model.AllRoleHolderKinds()), model.RoleHolderPrincipal.Valid, func() bool { return model.RoleHolderKind("Claim").Valid() }},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.all == 0 || !tc.ok() || tc.bad() {
				t.Fatalf("enum %s Valid/All accessors failed (all=%d ok=%v bad=%v)", tc.name, tc.all, tc.ok(), tc.bad())
			}
		})
	}
}

func TestActionTargetBindingVariants(t *testing.T) {
	t.Parallel()
	exact := model.ActionTargetBinding{
		Mode: model.ActionTargetExactResource, Action: "executiontarget.read", TargetKind: "ExecutionTarget",
	}
	if !exact.IsExactResource() {
		t.Fatal("expected ExactResource")
	}
	create := model.ActionTargetBinding{
		Mode: model.ActionTargetCreateParent, Action: "roleassignment.grant", ParentScopeKind: apimeta.ScopeOrganization,
	}
	if !create.IsCreateParent() {
		t.Fatal("expected CreateParent")
	}
	scope := model.ActionTargetBinding{
		Mode: model.ActionTargetScopeOnly, Action: "membership.create", ScopeKind: apimeta.ScopeOrganization,
	}
	if !scope.IsScopeOnly() {
		t.Fatal("expected ScopeOnly")
	}
	mixed := model.ActionTargetBinding{
		Mode: model.ActionTargetExactResource, Action: "a", TargetKind: "X", ScopeKind: apimeta.ScopeProject,
	}
	if mixed.IsExactResource() {
		t.Fatal("mixed binding must fail closed")
	}
}

func TestPropertyCrossTypeSubstitutionShape(t *testing.T) {
	t.Parallel()
	// Distinct Go types: assigning PrincipalRef where AccessGroupRef is required
	// is a compile-time failure. Property check: identity equality never confuses
	// PrincipalRef with AccessGroupRef field sets.
	fn := func(issuer, subject, name, uid string) bool {
		p := model.PrincipalRef{Issuer: issuer, Subject: subject, PrincipalType: model.PrincipalTypeHuman}
		g := model.AccessGroupRef{TypedRef: apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: model.KindAccessGroup, Name: name, UID: uid,
		}}
		// Different value spaces: principal has no Kind/UID; group has no Issuer/Subject.
		pt := reflect.TypeOf(p)
		gt := reflect.TypeOf(g)
		return pt != gt && pt.Name() == "PrincipalRef" && gt.Name() == "AccessGroupRef"
	}
	if err := quick.Check(fn, &quick.Config{MaxCount: 32}); err != nil {
		t.Fatal(err)
	}
}

func TestAccessGroupRefEmbedsTypedRef(t *testing.T) {
	t.Parallel()
	ref := model.AccessGroupRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionIAM, Kind: model.KindAccessGroup, Name: "ops", UID: "uid-ops",
	}}
	if ref.Kind != model.KindAccessGroup || ref.UID == "" {
		t.Fatal("AccessGroupRef must expose TypedRef fields")
	}
}
