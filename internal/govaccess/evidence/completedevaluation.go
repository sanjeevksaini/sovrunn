package evidence

import (
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// CompletedAuthorizationEvaluation is the sealed completed current Allow|Deny
// evaluation produced after stage-10 finalization (ADH-2026-072 §2; DD-12).
type CompletedAuthorizationEvaluation struct {
	sealed    bool
	finalized FinalizedAuthorizationCandidate
	admission AuthorizationAdmissionProof
	context   PublicationContext
	outcome   AuthorizationOutcome
}

// CompleteAuthorizationCandidate binds a finalized candidate to admission and
// context, producing the sole completed evaluation value.
func CompleteAuthorizationCandidate(
	finalized FinalizedAuthorizationCandidate,
	admission AuthorizationAdmissionProof,
	ctx PublicationContext,
) (CompletedAuthorizationEvaluation, operation.MechanicalFailure) {
	if !finalized.sealed {
		return CompletedAuthorizationEvaluation{}, operation.NewMechanicalFailure("completed_eval_finalized_unsealed")
	}
	if !admission.sealed {
		return CompletedAuthorizationEvaluation{}, operation.NewMechanicalFailure("completed_eval_admission_unsealed")
	}
	if !ctx.sealed {
		return CompletedAuthorizationEvaluation{}, operation.NewMechanicalFailure("completed_eval_context_unsealed")
	}
	if !finalized.admission.admission.Equal(admission.admission) {
		return CompletedAuthorizationEvaluation{}, operation.NewMechanicalFailure("completed_eval_admission_mismatch")
	}
	if !finalized.context.Equal(ctx) || !admission.context.Equal(ctx) {
		return CompletedAuthorizationEvaluation{}, operation.NewMechanicalFailure("completed_eval_context_mismatch")
	}
	if !finalized.Outcome().Valid() {
		return CompletedAuthorizationEvaluation{}, operation.NewMechanicalFailure("completed_eval_outcome_invalid")
	}
	return CompletedAuthorizationEvaluation{
		sealed:    true,
		finalized: finalized,
		admission: admission,
		context:   ctx,
		outcome:   finalized.Outcome(),
	}, operation.MechanicalFailure{}
}

// Sealed reports whether e was package-constructed.
func (e CompletedAuthorizationEvaluation) Sealed() bool { return e.sealed }

// Outcome returns the sealed Allow|Deny outcome.
func (e CompletedAuthorizationEvaluation) Outcome() AuthorizationOutcome { return e.outcome }

// FinalizedCandidate returns the sealed finalized candidate.
func (e CompletedAuthorizationEvaluation) FinalizedCandidate() FinalizedAuthorizationCandidate {
	return e.finalized
}

// AdmissionProof returns the sealed admission proof.
func (e CompletedAuthorizationEvaluation) AdmissionProof() AuthorizationAdmissionProof {
	return e.admission
}

// PublicationContext returns the sealed publication context.
func (e CompletedAuthorizationEvaluation) PublicationContext() PublicationContext { return e.context }
