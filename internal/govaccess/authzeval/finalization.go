package authzeval

import (
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// ReevaluateCandidate directs callers to recompute authorization before
// publication because candidate predicates are no longer current.
type ReevaluateCandidate struct {
	sealed bool
	reason string
}

func (r ReevaluateCandidate) Sealed() bool   { return r.sealed }
func (r ReevaluateCandidate) Reason() string { return r.reason }

// FinalizedCandidate is the publication-time validated candidate used for
// post-commit release.
type FinalizedCandidate struct {
	sealed bool

	candidate     AuthorizationCandidate
	publicationAt time.Time
}

func (c FinalizedCandidate) Sealed() bool                      { return c.sealed }
func (c FinalizedCandidate) Candidate() AuthorizationCandidate { return c.candidate }
func (c FinalizedCandidate) PublicationAt() time.Time          { return c.publicationAt }

// FinalizeCandidateOutcome is a closed union of finalized, reevaluate, or
// mechanical failure.
type FinalizeCandidateOutcome struct {
	kind       uint8
	finalized  FinalizedCandidate
	reevaluate ReevaluateCandidate
	failure    operation.MechanicalFailure
}

const (
	finalizeNone uint8 = iota
	finalizeSuccess
	finalizeReevaluate
	finalizeFailure
)

func (o FinalizeCandidateOutcome) AsFinalized() (FinalizedCandidate, bool) {
	if o.kind != finalizeSuccess {
		return FinalizedCandidate{}, false
	}
	return o.finalized, true
}

func (o FinalizeCandidateOutcome) AsReevaluate() (ReevaluateCandidate, bool) {
	if o.kind != finalizeReevaluate {
		return ReevaluateCandidate{}, false
	}
	return o.reevaluate, true
}

func (o FinalizeCandidateOutcome) AsMechanicalFailure() (operation.MechanicalFailure, bool) {
	if o.kind != finalizeFailure {
		return operation.MechanicalFailure{}, false
	}
	return o.failure, true
}

// FinalizeCandidateAt revalidates time/membership/lifecycle predicates at the
// transaction publication instant.
func FinalizeCandidateAt(
	candidate AuthorizationCandidate,
	admission evidence.AuthorizationAdmissionProof,
	publicationAt time.Time,
) FinalizeCandidateOutcome {
	if !candidate.Sealed() {
		return FinalizeCandidateOutcome{kind: finalizeFailure, failure: operation.NewMechanicalFailure("candidate_unsealed")}
	}
	if !admission.Sealed() {
		return FinalizeCandidateOutcome{kind: finalizeFailure, failure: operation.NewMechanicalFailure("admission_unsealed")}
	}
	if publicationAt.IsZero() {
		return FinalizeCandidateOutcome{kind: finalizeFailure, failure: operation.NewMechanicalFailure("publication_time_zero")}
	}
	in := candidate.Input()
	revalidated, ok := model.NewAuthorizationInput(
		in.Actor(),
		in.Action(),
		in.ScopeRef(),
		in.TargetRef(),
		publicationAt.UTC(),
		in.RequestID(),
		in.DirectAssignments(),
		in.GroupAssignments(),
		in.ActiveGroupUIDs(),
		in.RoleActions(),
		in.PrivilegedRoleUIDs(),
		in.SuspendedRoleUIDs(),
		in.ScopeMembershipCurrent(),
		in.GuardrailAllows(),
		in.PolicyOutcome(),
	)
	if !ok {
		return FinalizeCandidateOutcome{kind: finalizeFailure, failure: operation.NewMechanicalFailure("candidate_input_invalid")}
	}
	decision, _, _ := evaluateDecision(revalidated)
	if decision != candidate.Decision() {
		return FinalizeCandidateOutcome{
			kind: finalizeReevaluate,
			reevaluate: ReevaluateCandidate{
				sealed: true,
				reason: "authorization_predicates_changed",
			},
		}
	}
	return FinalizeCandidateOutcome{
		kind: finalizeSuccess,
		finalized: FinalizedCandidate{
			sealed:        true,
			candidate:     candidate,
			publicationAt: publicationAt.UTC(),
		},
	}
}
