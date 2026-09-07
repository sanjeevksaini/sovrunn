package uow

import (
	"context"
	"errors"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/authzeval"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// CoordinateCallerMutation coordinates stage-10/11/12 for caller mutations.
func (c *Coordinator) CoordinateCallerMutation(ctx context.Context, req CallerMutationRequest) MutationOutcome {
	if c == nil || c.store == nil {
		return mechanicalOutcome("coordinator_store_nil")
	}
	if !req.LookupKey.Sealed() || !req.Binding.Sealed() {
		return mechanicalOutcome("caller_idempotency_unsealed")
	}
	if !req.Candidate.Sealed() {
		return mechanicalOutcome("caller_candidate_unsealed")
	}
	if req.BuildEvidence == nil {
		return mechanicalOutcome("caller_build_evidence_nil")
	}

	var lease *state.CallerReservationLease
	for {
		outcome := c.store.InspectOrReserveCallerMutation(req.LookupKey, req.Binding)
		if outcome.IsConflict() {
			return MutationOutcome{kind: OutcomeConflict}
		}
		if replay, ok := outcome.Replay(); ok {
			return MutationOutcome{kind: OutcomeReplay, replay: replay, hasReplay: true}
		}
		if _, wait := outcome.WaitKey(); wait {
			select {
			case <-ctx.Done():
				return MutationOutcome{kind: OutcomeWaitTimeout}
			case <-time.After(waitRetryInterval):
			}
			continue
		}
		l, ok := outcome.Owner()
		if !ok || l == nil {
			return mechanicalOutcome("caller_reservation_owner_missing")
		}
		lease = l
		break
	}

	tx, err := lease.Begin(req.Admission)
	if err != nil {
		if errors.Is(err, state.ErrConflict) {
			return MutationOutcome{kind: OutcomeConflict}
		}
		return mechanicalOutcome("caller_begin_failed")
	}
	defer tx.Abort()

	finalized, kind, fail := finalizeAuthz(req.Candidate, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	if kind == OutcomeReevaluate {
		return MutationOutcome{kind: OutcomeReevaluate}
	}
	evidenceFinalized, fail := req.BuildEvidence(finalized, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	require := evidence.RequireCurrentAllow(evidenceFinalized, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	if deny, ok := require.AsDenyRoute(); ok && deny.Sealed() {
		tx.Abort()
		return c.coordinateDeniedEvidenceOnly(req.Candidate.Currentness(), req.Candidate, req.BuildEvidence)
	}
	allow, ok := require.AsAllow()
	if !ok {
		if fail, ok := require.AsMechanicalFailure(); ok {
			return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
		}
		return mechanicalOutcome("caller_require_current_allow_invalid")
	}
	if fail := tx.AcceptCurrentAllow(allow); fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}

	proofs, resultMaterial, domain, ok, fail := applyParticipantsCaller(tx, req.Participants)
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	if ok {
		tx.Abort()
		return MutationOutcome{kind: OutcomeDomain, domainCode: domain}
	}

	completed, fail := evidence.CompleteAuthorizationCandidate(evidenceFinalized, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	plan, fail := evidence.NewAuthorizedMutationPlan(tx.PublicationContext(), completed, proofs)
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	preparedEvidence, fail := evidence.PrepareMutationCarrierSet(&plan, nil)
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	if fail := tx.AcceptPreparedEvidence(preparedEvidence); fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	digest, fail := evidence.NewMutationDigestFromProofs(proofs)
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	completion, fail := evidence.CompletionBindingForPlan(tx.PublicationContext(), digest)
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	preparedCompleted, fail := idempotency.StageCompleted(req.LookupKey, req.Binding, resultMaterial, completion)
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	if fail := tx.StageCompletedResult(preparedCompleted); fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	if fail := tx.Seal(); fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	receipt, fail := tx.Commit()
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	released, fail := authzRelease(finalized, receipt)
	if fail.Reason() != "" {
		return MutationOutcome{kind: OutcomeMechanicalErr, failure: fail}
	}
	return MutationOutcome{
		kind:          OutcomeReleased,
		authzResult:   released,
		hasAuthz:      true,
		commitReceipt: receipt,
		hasCommit:     true,
	}
}

func authzRelease(finalized authzeval.FinalizedCandidate, receipt state.CommitReceipt) (model.AuthorizationResult, operation.MechanicalFailure) {
	return authzeval.ReleaseAfterPublication(finalized, receipt)
}

func applyParticipantsCaller(
	tx state.CallerMutationTransaction,
	intents []FinalizableIntent,
) ([]evidence.AppliedMutationProof, operation.ResultMaterial, string, bool, operation.MechanicalFailure) {
	proofs := make([]evidence.AppliedMutationProof, 0, len(intents))
	var result operation.ResultMaterial
	hasResult := false

	for _, intent := range intents {
		permit, fail := tx.FinalizationPermit(intent.ParticipantClaim())
		if fail.Reason() != "" {
			return nil, operation.ResultMaterial{}, "", false, fail
		}
		outcome := intent.FinalizeAt(permit)
		if fail, ok := outcome.AsFailure(); ok {
			return nil, operation.ResultMaterial{}, "", false, fail
		}
		if domain, ok := outcome.AsDomain(); ok {
			return nil, operation.ResultMaterial{}, domain, true, operation.MechanicalFailure{}
		}
		change, ok := outcome.AsChange()
		if !ok {
			return nil, operation.ResultMaterial{}, "", false, operation.NewMechanicalFailure("finalized_change_missing")
		}
		receipt, fail := tx.ApplyFinalized(change)
		if fail.Reason() != "" {
			return nil, operation.ResultMaterial{}, "", false, fail
		}
		if receipt.PublicationKind() != state.ResourceMutation {
			continue
		}
		proof, ok := receipt.EvidenceProof()
		if !ok {
			return nil, operation.ResultMaterial{}, "", false, operation.NewMechanicalFailure("mutation_proof_missing")
		}
		proofs = append(proofs, proof)
		if intent.ParticipantClaim().Role() == state.OriginatingResult {
			res, ok := change.CanonicalResult()
			if !ok || !res.Sealed() {
				return nil, operation.ResultMaterial{}, "", false, operation.NewMechanicalFailure("originating_result_missing")
			}
			result = res
			hasResult = true
		}
	}
	if len(proofs) == 0 {
		return nil, operation.ResultMaterial{}, "", false, operation.NewMechanicalFailure("mutation_proofs_empty")
	}
	if !hasResult {
		return nil, operation.ResultMaterial{}, "", false, operation.NewMechanicalFailure("originating_result_unavailable")
	}
	return proofs, result, "", false, operation.MechanicalFailure{}
}
