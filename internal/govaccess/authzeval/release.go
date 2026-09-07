package authzeval

import (
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// ReleaseAfterPublication materializes AuthorizationResult only after a sealed
// state commit receipt is present.
func ReleaseAfterPublication(
	finalized FinalizedCandidate,
	receipt state.CommitReceipt,
) (model.AuthorizationResult, operation.MechanicalFailure) {
	if !finalized.Sealed() {
		return model.AuthorizationResult{}, operation.NewMechanicalFailure("finalized_candidate_unsealed")
	}
	if !receipt.Sealed() {
		return model.AuthorizationResult{}, operation.NewMechanicalFailure("commit_receipt_unsealed")
	}
	if !receipt.PublicationContext().PublicationInstant().Equal(finalized.PublicationAt()) {
		return model.AuthorizationResult{}, operation.NewMechanicalFailure("publication_instant_mismatch")
	}
	candidate := finalized.Candidate()
	result, ok := model.NewAuthorizationResult(
		candidate.Decision(),
		candidate.ReasonCodes(),
		candidate.Input().RequestID(),
		candidate.Input().Action(),
		candidate.Input().TargetRef(),
		candidate.Input().EvaluatedAt(),
		receipt.PublicationContext().PublicationInstant(),
		candidate.ContributingGrants(),
	)
	if !ok {
		return model.AuthorizationResult{}, operation.NewMechanicalFailure("authorization_result_invalid")
	}
	return result, operation.MechanicalFailure{}
}
