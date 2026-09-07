package uow

import (
	"context"
	"errors"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/authzeval"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// CoordinateEvidenceOnly coordinates denied/read/replay evidence-only commit.
func (c *Coordinator) CoordinateEvidenceOnly(_ context.Context, req EvidenceOnlyRequest) MutationOutcome {
	if c == nil || c.store == nil {
		return mechanicalOutcome("coordinator_store_nil")
	}
	if !req.Candidate.Sealed() {
		return mechanicalOutcome("evidence_only_candidate_unsealed")
	}
	if req.BuildEvidence == nil {
		return mechanicalOutcome("evidence_only_build_evidence_nil")
	}
	return c.coordinateDeniedEvidenceOnly(req.Candidate.Currentness(), req.Candidate, req.BuildEvidence)
}

func (c *Coordinator) coordinateDeniedEvidenceOnly(
	claim state.AuthorizationCurrentnessClaim,
	candidate authzeval.OperationCandidate,
	build EvidenceFinalizedBuilder,
) MutationOutcome {
	tx, err := c.store.BeginEvidenceTransaction(claim)
	if err != nil {
		if errors.Is(err, state.ErrRetryEvaluation) {
			return MutationOutcome{kind: OutcomeReevaluate}
		}
		return mechanicalOutcome("evidence_begin_failed")
	}
	defer tx.Abort()

	finalized, kind, fail := finalizeAuthz(candidate, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	if kind == OutcomeReevaluate {
		return MutationOutcome{kind: OutcomeReevaluate}
	}
	evidenceFinalized, fail := build(finalized, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	completed, fail := evidence.CompleteAuthorizationCandidate(evidenceFinalized, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	prepared, fail := evidence.PrepareAuthorizationCarrierSet(completed)
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
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
	result, fail := authzeval.ReleaseAfterPublication(finalized, receipt)
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	return MutationOutcome{
		kind:          OutcomeReleased,
		authzResult:   result,
		hasAuthz:      true,
		commitReceipt: receipt,
		hasCommit:     true,
	}
}
