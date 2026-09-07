package approvalreq

import (
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// RoleGrantInput is the typed role-grant derivation input.
type RoleGrantInput struct {
	Proposal   model.RoleAssignmentProposal
	Privileged bool
	PolicyRef  *apimeta.TypedRef
}

// MembershipInput is reserved for Task-5 owner implementation.
type MembershipInput struct{}

// PrivilegedInput is reserved for Task-5 owner implementation.
type PrivilegedInput struct{}

// ExceptionInput is reserved for Task-5 owner implementation.
type ExceptionInput struct{}

// ApprovalRequirementDeriver is the shared marker interface for typed derivers.
type ApprovalRequirementDeriver interface {
	DeriverKind() string
}

// RoleGrantRequirementDeriverPort is the typed role-grant derivation port.
type RoleGrantRequirementDeriverPort interface {
	ApprovalRequirementDeriver
	DeriveForRoleGrant(RoleGrantInput) (model.ApprovalRequirement, bool)
}

// MembershipRequirementDeriverPort is reserved for Task 5.
type MembershipRequirementDeriverPort interface {
	ApprovalRequirementDeriver
	DeriveForMembership(MembershipInput) (model.ApprovalRequirement, bool)
}

// PrivilegedRequirementDeriverPort is reserved for Task 5.
type PrivilegedRequirementDeriverPort interface {
	ApprovalRequirementDeriver
	DeriveForPrivileged(PrivilegedInput) (model.ApprovalRequirement, bool)
}

// ExceptionRequirementDeriverPort is reserved for Task 5.
type ExceptionRequirementDeriverPort interface {
	ApprovalRequirementDeriver
	DeriveForException(ExceptionInput) (model.ApprovalRequirement, bool)
}
