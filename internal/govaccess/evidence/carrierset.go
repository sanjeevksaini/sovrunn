package evidence

import (
	"fmt"

	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// CarrierKind identifies a staged evidence carrier class.
type CarrierKind string

const (
	CarrierAuditEvent     CarrierKind = "AuditEvent"
	CarrierDecisionRecord CarrierKind = "DecisionRecord"
)

// CarrierRef is a deterministic preallocated carrier identity.
type CarrierRef struct {
	Kind CarrierKind
	ID   string
}

// CarrierSet is the exact sealed set of carriers for a prepared evidence change.
type CarrierSet struct {
	sealed    bool
	carriers  []CarrierRef
	audit     []decision.AuditLinkage
	decisions []decision.ProfileRef
}

// Sealed reports whether s was package-constructed.
func (s CarrierSet) Sealed() bool { return s.sealed }

// Carriers returns a defensive copy of carrier refs.
func (s CarrierSet) Carriers() []CarrierRef {
	return append([]CarrierRef(nil), s.carriers...)
}

// DecisionProfiles returns pinned DecisionProfile refs.
func (s CarrierSet) DecisionProfiles() []decision.ProfileRef {
	return append([]decision.ProfileRef(nil), s.decisions...)
}

func buildAuthorizedCarrierSet(plan AuthorizedMutationPlan) (CarrierSet, operation.MechanicalFailure) {
	carriers := []CarrierRef{
		{Kind: CarrierAuditEvent, ID: carrierID(plan.context, "authorization.evaluated", 0)},
	}
	decisions := []decision.ProfileRef{}
	ordinal := 1
	for _, proof := range plan.proofs {
		for _, d := range proof.descriptors {
			carriers = append(carriers, CarrierRef{
				Kind: CarrierAuditEvent,
				ID:   carrierID(plan.context, d.EventType(), ordinal),
			})
			ordinal++
		}
		for _, m := range proof.materials {
			carriers = append(carriers, CarrierRef{
				Kind: CarrierDecisionRecord,
				ID:   carrierID(plan.context, string(m.profileID), ordinal),
			})
			decisions = append(decisions, decision.ProfileRef{
				Name:    string(m.profileID),
				Version: m.profileVer,
			})
			ordinal++
		}
	}
	// Mandatory authorization-decision/v1 when evaluation exists (F18-RD-19).
	carriers = append(carriers, CarrierRef{
		Kind: CarrierDecisionRecord,
		ID:   carrierID(plan.context, ProfileAuthorizationDecision, ordinal),
	})
	decisions = append(decisions, decision.ProfileRef{
		Name:    ProfileAuthorizationDecision,
		Version: ProfileVersionV1,
	})
	return CarrierSet{sealed: true, carriers: carriers, decisions: decisions}, operation.MechanicalFailure{}
}

func buildAutomaticCarrierSet(plan AutomaticMutationPlan) (CarrierSet, operation.MechanicalFailure) {
	carriers := []CarrierRef{}
	decisions := []decision.ProfileRef{}
	ordinal := 0
	for _, proof := range plan.proofs {
		for _, d := range proof.descriptors {
			carriers = append(carriers, CarrierRef{
				Kind: CarrierAuditEvent,
				ID:   carrierID(plan.context, d.EventType(), ordinal),
			})
			ordinal++
		}
		for _, m := range proof.materials {
			carriers = append(carriers, CarrierRef{
				Kind: CarrierDecisionRecord,
				ID:   carrierID(plan.context, string(m.profileID), ordinal),
			})
			decisions = append(decisions, decision.ProfileRef{
				Name:    string(m.profileID),
				Version: m.profileVer,
			})
			ordinal++
		}
	}
	if len(carriers) == 0 {
		return CarrierSet{}, operation.NewMechanicalFailure("automatic_carrier_set_empty")
	}
	return CarrierSet{sealed: true, carriers: carriers, decisions: decisions}, operation.MechanicalFailure{}
}

func buildDomainConclusionCarrierSet(plan DomainConclusionPlan) (CarrierSet, operation.MechanicalFailure) {
	carriers := []CarrierRef{
		{Kind: CarrierAuditEvent, ID: carrierID(plan.context, "authorization.evaluated", 0)},
	}
	decisions := []decision.ProfileRef{
		{Name: ProfileAuthorizationDecision, Version: ProfileVersionV1},
	}
	ordinal := 1
	for _, proof := range plan.proofs {
		carriers = append(carriers, CarrierRef{
			Kind: CarrierAuditEvent,
			ID:   carrierID(plan.context, proof.descriptor.EventType(), ordinal),
		})
		ordinal++
		for _, m := range proof.materials {
			if m.profileID != DomainDecisionProfileID(ProfileExceptionDecision) {
				return CarrierSet{}, operation.NewMechanicalFailure("conclusion_carrier_requires_exception_material")
			}
			carriers = append(carriers, CarrierRef{
				Kind: CarrierDecisionRecord,
				ID:   carrierID(plan.context, string(m.profileID), ordinal),
			})
			decisions = append(decisions, decision.ProfileRef{
				Name:    string(m.profileID),
				Version: m.profileVer,
			})
			ordinal++
		}
	}
	return CarrierSet{sealed: true, carriers: carriers, decisions: decisions}, operation.MechanicalFailure{}
}

func carrierID(ctx PublicationContext, label string, ordinal int) string {
	return fmt.Sprintf("f18-%d-%s-%d", ctx.publicationSequence, label, ordinal)
}
