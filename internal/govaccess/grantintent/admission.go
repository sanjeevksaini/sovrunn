package grantintent

import (
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/approvalreq"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// DirectGrant carries one non-delegable direct grant submission.
type DirectGrant struct {
	ResourceUID     string
	ResourceName    string
	ScopeRef        apimeta.ScopeRef
	Spec            model.RoleAssignmentSpec
	ExpectedVersion string
	NextVersion     string
	Privileged      bool
	PolicyRef       *apimeta.TypedRef
}

// AdmittedGrant is the immutable admitted grant/proposal union.
type AdmittedGrant struct {
	sealed bool

	participantID state.ParticipantID
	requirement   model.ApprovalRequirement
	proposal      model.RoleAssignmentProposal
}

func (a AdmittedGrant) Sealed() bool { return a.sealed }
func (a AdmittedGrant) ParticipantID() state.ParticipantID {
	return a.participantID
}
func (a AdmittedGrant) Requirement() model.ApprovalRequirement { return a.requirement }
func (a AdmittedGrant) Proposal() model.RoleAssignmentProposal { return a.proposal }

// AdmitDirectGrant admits one direct non-delegable role grant.
func AdmitDirectGrant(
	participantID state.ParticipantID,
	in DirectGrant,
	deriver approvalreq.RoleGrantRequirementDeriverPort,
) (AdmittedGrant, operation.MechanicalFailure) {
	proposal, ok := model.NewRoleAssignmentProposal(
		"direct-grant",
		in.ResourceUID,
		in.ResourceName,
		in.ScopeRef,
		in.Spec,
		in.ExpectedVersion,
		in.NextVersion,
		in.Privileged,
		in.PolicyRef,
		nil,
	)
	if !ok {
		return AdmittedGrant{}, operation.NewMechanicalFailure("direct_grant_invalid")
	}
	return AdmitRoleAssignmentProposal(participantID, proposal, deriver)
}

// AdmitRoleAssignmentProposal admits exactly one model.RoleAssignmentProposal.
func AdmitRoleAssignmentProposal(
	participantID state.ParticipantID,
	proposal model.RoleAssignmentProposal,
	deriver approvalreq.RoleGrantRequirementDeriverPort,
) (AdmittedGrant, operation.MechanicalFailure) {
	if participantID == "" {
		return AdmittedGrant{}, operation.NewMechanicalFailure("participant_id_empty")
	}
	if !proposal.Sealed() {
		return AdmittedGrant{}, operation.NewMechanicalFailure("proposal_unsealed")
	}
	if deriver == nil {
		return AdmittedGrant{}, operation.NewMechanicalFailure("deriver_nil")
	}
	policy, hasPolicy := proposal.PolicyRef()
	var policyRef *apimeta.TypedRef
	if hasPolicy {
		p := policy
		policyRef = &p
	}
	requirement, ok := deriver.DeriveForRoleGrant(approvalreq.RoleGrantInput{
		Proposal:   proposal,
		Privileged: proposal.Privileged(),
		PolicyRef:  policyRef,
	})
	if !ok || !requirement.Valid() {
		return AdmittedGrant{}, operation.NewMechanicalFailure("approval_requirement_invalid")
	}
	return AdmittedGrant{
		sealed:        true,
		participantID: participantID,
		requirement:   requirement,
		proposal:      proposal,
	}, operation.MechanicalFailure{}
}
