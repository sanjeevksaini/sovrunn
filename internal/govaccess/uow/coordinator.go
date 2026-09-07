package uow

import (
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/authzeval"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

const waitRetryInterval = 2 * time.Millisecond

// Coordinator is the sole stage-10/11/12 mutation coordinator.
type Coordinator struct {
	store *state.Store
}

// NewCoordinator constructs the stage-10/11/12 coordinator.
func NewCoordinator(store *state.Store) *Coordinator {
	return &Coordinator{store: store}
}

func mechanicalOutcome(reason string) MutationOutcome {
	return MutationOutcome{kind: OutcomeMechanicalErr, failure: operation.NewMechanicalFailure(reason)}
}

func finalizeAuthz(
	candidate authzeval.OperationCandidate,
	admission evidence.AuthorizationAdmissionProof,
	ctx evidence.PublicationContext,
) (authzeval.FinalizedCandidate, MutationOutcomeKind, operation.MechanicalFailure) {
	out := authzeval.FinalizeCandidateAt(candidate.Candidate(), admission, ctx.PublicationInstant())
	if fail, ok := out.AsMechanicalFailure(); ok {
		return authzeval.FinalizedCandidate{}, OutcomeMechanicalErr, fail
	}
	if _, ok := out.AsReevaluate(); ok {
		return authzeval.FinalizedCandidate{}, OutcomeReevaluate, operation.MechanicalFailure{}
	}
	finalized, ok := out.AsFinalized()
	if !ok {
		return authzeval.FinalizedCandidate{}, OutcomeMechanicalErr, operation.NewMechanicalFailure("finalize_outcome_invalid")
	}
	return finalized, OutcomeReleased, operation.MechanicalFailure{}
}
