package uow

import (
	"context"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/authzeval"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// FinalizedIntentOutcome is the closed finalization outcome consumed by uow.
type FinalizedIntentOutcome struct {
	change    state.ContextBoundChange
	hasChange bool

	domainCode string
	hasDomain  bool

	failure operation.MechanicalFailure
}

// NewFinalizedChangeOutcome constructs a finalized-change outcome.
func NewFinalizedChangeOutcome(change state.ContextBoundChange) FinalizedIntentOutcome {
	return FinalizedIntentOutcome{change: change, hasChange: true}
}

// NewDomainOutcome constructs a non-publication domain outcome.
func NewDomainOutcome(code string) FinalizedIntentOutcome {
	return FinalizedIntentOutcome{domainCode: code, hasDomain: true}
}

// NewFailureOutcome constructs a mechanical-failure outcome.
func NewFailureOutcome(f operation.MechanicalFailure) FinalizedIntentOutcome {
	return FinalizedIntentOutcome{failure: f}
}

// AsChange returns the finalized change when present.
func (o FinalizedIntentOutcome) AsChange() (state.ContextBoundChange, bool) {
	return o.change, o.hasChange
}

// AsDomain returns the domain code when present.
func (o FinalizedIntentOutcome) AsDomain() (string, bool) {
	return o.domainCode, o.hasDomain
}

// AsFailure returns mechanical failure when present.
func (o FinalizedIntentOutcome) AsFailure() (operation.MechanicalFailure, bool) {
	if o.failure.Reason() == "" {
		return operation.MechanicalFailure{}, false
	}
	return o.failure, true
}

// FinalizableIntent is the stage-10/11 participant intent consumed by uow.
type FinalizableIntent interface {
	ParticipantClaim() state.ParticipantClaim
	FinalizeAt(permit state.MutationFinalizationPermit) FinalizedIntentOutcome
}

// EvidenceFinalizedBuilder builds evidence.FinalizedAuthorizationCandidate from
// finalized authzeval output and the transaction admission/context.
type EvidenceFinalizedBuilder func(
	finalized authzeval.FinalizedCandidate,
	admission evidence.AuthorizationAdmissionProof,
	ctx evidence.PublicationContext,
) (evidence.FinalizedAuthorizationCandidate, operation.MechanicalFailure)

// CallerMutationRequest coordinates an authorization-bearing caller mutation.
type CallerMutationRequest struct {
	LookupKey     idempotency.IdempotencyLookupKey
	Binding       idempotency.RequestBinding
	Admission     state.MutationAdmission
	Candidate     authzeval.OperationCandidate
	BuildEvidence EvidenceFinalizedBuilder
	Participants  []FinalizableIntent
}

// ControllerMutationRequest coordinates a controller mutation.
type ControllerMutationRequest struct {
	Admission     state.MutationAdmission
	Candidate     *authzeval.OperationCandidate // required for authorization-bearing controller mutations
	BuildEvidence EvidenceFinalizedBuilder
	Participants  []FinalizableIntent
}

// EvidenceOnlyRequest coordinates an evidence-only transaction.
type EvidenceOnlyRequest struct {
	Candidate     authzeval.OperationCandidate
	BuildEvidence EvidenceFinalizedBuilder
}

// MutationOutcomeKind is the closed uow outcome.
type MutationOutcomeKind string

const (
	OutcomeReleased      MutationOutcomeKind = "Released"
	OutcomeReplay        MutationOutcomeKind = "Replay"
	OutcomeConflict      MutationOutcomeKind = "Conflict"
	OutcomeWaitTimeout   MutationOutcomeKind = "WaitTimeout"
	OutcomeDomain        MutationOutcomeKind = "Domain"
	OutcomeReevaluate    MutationOutcomeKind = "Reevaluate"
	OutcomeMechanicalErr MutationOutcomeKind = "MechanicalFailure"
)

// MutationOutcome is a closed result for caller/controller/evidence paths.
type MutationOutcome struct {
	kind MutationOutcomeKind

	replay    idempotency.CompletedRecord
	hasReplay bool

	authzResult model.AuthorizationResult
	hasAuthz    bool

	commitReceipt state.CommitReceipt
	hasCommit     bool

	domainCode string
	failure    operation.MechanicalFailure
}

// Kind returns the closed outcome kind.
func (o MutationOutcome) Kind() MutationOutcomeKind { return o.kind }

// Replay returns replay record when outcome is Replay.
func (o MutationOutcome) Replay() (idempotency.CompletedRecord, bool) {
	return o.replay, o.hasReplay
}

// AuthorizationResult returns released AuthorizationResult when present.
func (o MutationOutcome) AuthorizationResult() (model.AuthorizationResult, bool) {
	return o.authzResult, o.hasAuthz
}

// CommitReceipt returns commit receipt when present.
func (o MutationOutcome) CommitReceipt() (state.CommitReceipt, bool) {
	return o.commitReceipt, o.hasCommit
}

// DomainCode returns domain code when outcome is Domain.
func (o MutationOutcome) DomainCode() (string, bool) {
	if o.kind != OutcomeDomain || o.domainCode == "" {
		return "", false
	}
	return o.domainCode, true
}

// MechanicalFailure returns failure when outcome is MechanicalFailure.
func (o MutationOutcome) MechanicalFailure() (operation.MechanicalFailure, bool) {
	if o.kind != OutcomeMechanicalErr {
		return operation.MechanicalFailure{}, false
	}
	return o.failure, true
}

// MutationCoordinatorPort is the narrow stage-10/11/12 coordinator surface.
type MutationCoordinatorPort interface {
	CoordinateCallerMutation(ctx context.Context, req CallerMutationRequest) MutationOutcome
	CoordinateControllerMutation(ctx context.Context, req ControllerMutationRequest) MutationOutcome
	CoordinateEvidenceOnly(ctx context.Context, req EvidenceOnlyRequest) MutationOutcome
}
