package approvalreq_test

import (
	"reflect"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/approvalreq"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/grantintent"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

func TestTypedDeriverPortSurfaceAndRoleGrantImplementation(t *testing.T) {
	t.Parallel()
	var _ approvalreq.ApprovalRequirementDeriver = grantintent.RoleGrantRequirementDeriver{}
	var _ approvalreq.RoleGrantRequirementDeriverPort = grantintent.RoleGrantRequirementDeriver{}

	deriver := grantintent.RoleGrantRequirementDeriver{}
	spec := model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{
			Kind: model.RoleHolderPrincipal,
			Principal: &model.PrincipalRef{
				Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman,
			},
		},
		RoleDefinitionRef: apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "ops", UID: "rd-1",
		},
		RoleDefinitionVersion: "v1",
		Validity:              model.AssignmentValidityStanding,
	}
	proposal, ok := model.NewRoleAssignmentProposal(
		"p-1", "ra-1", "ra-one",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "proj", UID: "scope-1",
		}},
		spec, "1", "2", false, nil, nil,
	)
	if !ok {
		t.Fatal("proposal")
	}
	req, ok := deriver.DeriveForRoleGrant(approvalreq.RoleGrantInput{Proposal: proposal})
	if !ok || req.Mode() != model.ApprovalRequirementNotRequired {
		t.Fatal("expected role grant deriver")
	}

	if reflect.TypeOf((*approvalreq.MembershipRequirementDeriverPort)(nil)).Elem().NumMethod() == 0 {
		t.Fatal("membership typed port must exist")
	}
	if reflect.TypeOf((*approvalreq.PrivilegedRequirementDeriverPort)(nil)).Elem().NumMethod() == 0 {
		t.Fatal("privileged typed port must exist")
	}
	if reflect.TypeOf((*approvalreq.ExceptionRequirementDeriverPort)(nil)).Elem().NumMethod() == 0 {
		t.Fatal("exception typed port must exist")
	}
}
