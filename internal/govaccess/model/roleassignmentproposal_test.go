package model_test

import (
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

func makeProposalSpec(privileged bool) model.RoleAssignmentSpec {
	holder := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	mode := model.AssignmentValidityStanding
	var nb, exp *time.Time
	if privileged {
		mode = model.AssignmentValidityTimeBound
		n := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
		e := n.Add(2 * time.Hour)
		nb = &n
		exp = &e
	}
	return model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{Kind: model.RoleHolderPrincipal, Principal: &holder},
		RoleDefinitionRef: apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "ops-admin", UID: "rd-1",
		},
		RoleDefinitionVersion: "v1",
		Validity:              mode,
		NotBefore:             nb,
		ExpiresAt:             exp,
	}
}

func TestRoleAssignmentProposalDeepCopyAndValidation(t *testing.T) {
	t.Parallel()
	spec := makeProposalSpec(true)
	ref := apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "ApprovalPolicy", Name: "jit", UID: "ap-1"}
	due := time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC)
	p, ok := model.NewRoleAssignmentProposal(
		"proposal-1",
		"ra-1",
		"ra-one",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: string(apimeta.ScopeProject), Name: "p1", UID: "scope-1"}},
		spec,
		"1",
		"2",
		true,
		&ref,
		&due,
	)
	if !ok || !p.Sealed() {
		t.Fatal("expected valid proposal")
	}
	got := p.Spec()
	got.RoleDefinitionVersion = "mutated"
	if p.Spec().RoleDefinitionVersion == "mutated" {
		t.Fatal("spec accessor must return deep copy")
	}
	nextDue := p.NextReviewDueAt()
	*nextDue = nextDue.Add(3 * time.Hour)
	if p.NextReviewDueAt().Equal(*nextDue) {
		t.Fatal("due accessor must return copy")
	}
	pr, has := p.PolicyRef()
	if !has || pr.UID != "ap-1" {
		t.Fatal("expected policy ref")
	}
	pr.UID = "mutated"
	pr2, _ := p.PolicyRef()
	if pr2.UID == "mutated" {
		t.Fatal("policy accessor must return copy")
	}
}

func TestRoleAssignmentProposalDeterministicAndNoSetter(t *testing.T) {
	t.Parallel()
	fn := func(resource string) bool {
		if resource == "" {
			resource = "ra-fixed"
		}
		spec := makeProposalSpec(false)
		a, oka := model.NewRoleAssignmentProposal(
			"p", resource, "name", apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: string(apimeta.ScopeProject), Name: "p1", UID: "scope-1"}},
			spec, "1", "2", false, nil, nil,
		)
		b, okb := model.NewRoleAssignmentProposal(
			"p", resource, "name", apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: string(apimeta.ScopeProject), Name: "p1", UID: "scope-1"}},
			spec, "1", "2", false, nil, nil,
		)
		return oka == okb && a.Sealed() == b.Sealed() && a.ResourceUID() == b.ResourceUID()
	}
	if err := quick.Check(fn, &quick.Config{MaxCount: 32}); err != nil {
		t.Fatal(err)
	}

	typ := reflect.TypeOf(model.RoleAssignmentProposal{})
	for i := 0; i < typ.NumMethod(); i++ {
		if len(typ.Method(i).Name) >= 3 && typ.Method(i).Name[:3] == "Set" {
			t.Fatalf("setter method forbidden: %s", typ.Method(i).Name)
		}
	}
}
