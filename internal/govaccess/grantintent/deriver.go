package grantintent

import (
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/approvalreq"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// RoleGrantRequirementDeriver is the sole Task-4 ApprovalRequirementDeriver implementation.
type RoleGrantRequirementDeriver struct{}

func (RoleGrantRequirementDeriver) DeriverKind() string { return "RoleGrant" }

func (RoleGrantRequirementDeriver) DeriveForRoleGrant(in approvalreq.RoleGrantInput) (model.ApprovalRequirement, bool) {
	if !in.Proposal.Sealed() {
		return model.ApprovalRequirement{}, false
	}
	if in.Privileged {
		ref := in.PolicyRef
		if ref == nil {
			if proposalRef, ok := in.Proposal.PolicyRef(); ok {
				ref = &proposalRef
			}
		}
		if ref == nil {
			return model.ApprovalRequirement{}, false
		}
		return model.NewApprovalRequirementRequired(*ref)
	}
	req := model.NewApprovalRequirementNotRequired()
	return req, req.Valid()
}
