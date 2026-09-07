package uow

import (
	"context"
	"errors"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// CoordinateControllerMutation coordinates controller mutation publication.
func (c *Coordinator) CoordinateControllerMutation(_ context.Context, req ControllerMutationRequest) MutationOutcome {
	if c == nil || c.store == nil {
		return mechanicalOutcome("coordinator_store_nil")
	}
	if req.BuildEvidence == nil {
		return mechanicalOutcome("controller_build_evidence_nil")
	}
	tx, err := c.store.BeginControllerMutation(req.Admission)
	if err != nil {
		if errors.Is(err, state.ErrConflict) || errors.Is(err, state.ErrNoOp) {
			return MutationOutcome{kind: OutcomeConflict}
		}
		return mechanicalOutcome("controller_begin_failed")
	}
	defer tx.Abort()

	hasAuthorization := false
	if proof, ok := tx.AuthorizationAdmissionProof(); ok {
		hasAuthorization = true
		if req.Candidate == nil || !req.Candidate.Sealed() {
			return mechanicalOutcome("controller_candidate_required")
		}
		finalized, kind, fail := finalizeAuthz(*req.Candidate, proof, tx.PublicationContext())
		if fail.Reason() != "" {
			return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
		}
		if kind == OutcomeReevaluate {
			return MutationOutcome{kind: OutcomeReevaluate}
		}
		evidenceFinalized, fail := req.BuildEvidence(finalized, proof, tx.PublicationContext())
		if fail.Reason() != "" {
			return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
		}
		require := evidence.RequireCurrentAllow(evidenceFinalized, proof, tx.PublicationContext())
		allow, ok := require.AsAllow()
		if !ok {
			if fail, ok := require.AsMechanicalFailure(); ok {
				return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
			}
			return MutationOutcome{kind: OutcomeDomain, domainCode: "authorization_denied"}
		}
		if fail := tx.AcceptCurrentAllow(allow); fail.Reason() != "" {
			return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
		}
	}

	proofs, domain, hasDomain, fail := applyParticipantsController(tx, req.Participants)
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	if hasDomain {
		tx.Abort()
		return MutationOutcome{kind: OutcomeDomain, domainCode: domain}
	}
	var prepared evidence.PreparedEvidenceChange
	if hasAuthorization {
		proof, ok := tx.AuthorizationAdmissionProof()
		if !ok || req.Candidate == nil {
			return mechanicalOutcome("controller_authz_missing")
		}
		finalized, kind, fail := finalizeAuthz(*req.Candidate, proof, tx.PublicationContext())
		if fail.Reason() != "" {
			return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
		}
		if kind == OutcomeReevaluate {
			return MutationOutcome{kind: OutcomeReevaluate}
		}
		evFin, fail := req.BuildEvidence(finalized, proof, tx.PublicationContext())
		if fail.Reason() != "" {
			return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
		}
		completed, fail := evidence.CompleteAuthorizationCandidate(evFin, proof, tx.PublicationContext())
		if fail.Reason() != "" {
			return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
		}
		plan, fail := evidence.NewAuthorizedMutationPlan(tx.PublicationContext(), completed, proofs)
		if fail.Reason() != "" {
			return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
		}
		prepared, fail = evidence.PrepareMutationCarrierSet(&plan, nil)
		if fail.Reason() != "" {
			return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
		}
	} else {
		plan, fail := evidence.NewAutomaticMutationPlan(tx.PublicationContext(), proofs)
		if fail.Reason() != "" {
			return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
		}
		prepared, fail = evidence.PrepareMutationCarrierSet(nil, &plan)
		if fail.Reason() != "" {
			return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
		}
	}
	if fail := tx.AcceptPreparedEvidence(prepared); fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	if fail := tx.Seal(); fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	receipt, fail := tx.Commit()
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	return MutationOutcome{kind: OutcomeReleased, commitReceipt: receipt, hasCommit: true}
}

func applyParticipantsController(
	tx state.ControllerMutationTransaction,
	intents []FinalizableIntent,
) ([]evidence.AppliedMutationProof, string, bool, operation.MechanicalFailure) {
	proofs := make([]evidence.AppliedMutationProof, 0, len(intents))
	for _, intent := range intents {
		permit, fail := tx.FinalizationPermit(intent.ParticipantClaim())
		if fail.Reason() != "" {
			return nil, "", false, fail
		}
		outcome := intent.FinalizeAt(permit)
		if fail, ok := outcome.AsFailure(); ok {
			return nil, "", false, fail
		}
		if domain, ok := outcome.AsDomain(); ok {
			return nil, domain, true, operation.MechanicalFailure{}
		}
		change, ok := outcome.AsChange()
		if !ok {
			return nil, "", false, operation.NewMechanicalFailure("finalized_change_missing")
		}
		receipt, fail := tx.ApplyFinalized(change)
		if fail.Reason() != "" {
			return nil, "", false, fail
		}
		if receipt.PublicationKind() != state.ResourceMutation {
			continue
		}
		proof, ok := receipt.EvidenceProof()
		if !ok {
			return nil, "", false, operation.NewMechanicalFailure("mutation_proof_missing")
		}
		proofs = append(proofs, proof)
	}
	if len(proofs) == 0 {
		return nil, "", false, operation.NewMechanicalFailure("mutation_proofs_empty")
	}
	return proofs, "", false, operation.MechanicalFailure{}
}
