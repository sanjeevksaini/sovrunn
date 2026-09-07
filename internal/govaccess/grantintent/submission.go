package grantintent

import (
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/roleassign"
)

// SubmitToGrantPort submits an immutable trigger intent to roleassign.GrantPort.
// grantintent never constructs roleassign prepared/finalized change values.
func SubmitToGrantPort(
	port roleassign.GrantPort,
	admitted AdmittedGrant,
) (roleassign.PreparedRoleAssignmentIntent, operation.MechanicalFailure) {
	var zero roleassign.PreparedRoleAssignmentIntent
	if !admitted.Sealed() {
		return zero, operation.NewMechanicalFailure("admitted_grant_unsealed")
	}
	proposal := admitted.Proposal()
	trigger, fail := roleassign.NewGrantTriggerIntent(
		admitted.ParticipantID(),
		proposal.ResourceUID(),
		proposal.ResourceName(),
		proposal.ScopeRef(),
		proposal.Spec(),
		proposal.ExpectedVersion(),
		proposal.NextVersion(),
		nil,
		"",
		proposal.NextReviewDueAt(),
	)
	if fail.Reason() != "" {
		return zero, fail
	}
	return port.Submit(trigger)
}
