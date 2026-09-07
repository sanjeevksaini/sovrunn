package evidence

import (
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// DomainDecisionProfileID identifies the F18-RD-19 profile used for material.
type DomainDecisionProfileID string

// DomainDecisionMaterial is profile-specific sealed domain-decision material
// for mandatory F18-RD-19 DecisionRecords.
type DomainDecisionMaterial struct {
	sealed       bool
	profileID    DomainDecisionProfileID
	profileVer   string
	context      PublicationContext
	approval     model.ApprovalDecisionFacts
	exception    model.ExceptionDecisionFacts
	hasApproval  bool
	hasException bool
}

// NewApprovalDecisionMaterial seals approval-decision/v1 material from the
// exact design §3.2 carrier-input type model.ApprovalDecisionFacts.
func NewApprovalDecisionMaterial(
	facts model.ApprovalDecisionFacts,
	ctx PublicationContext,
) (DomainDecisionMaterial, operation.MechanicalFailure) {
	if !facts.Sealed() {
		return DomainDecisionMaterial{}, operation.NewMechanicalFailure("approval_material_facts_unsealed")
	}
	if !ctx.sealed {
		return DomainDecisionMaterial{}, operation.NewMechanicalFailure("approval_material_context_unsealed")
	}
	return DomainDecisionMaterial{
		sealed:      true,
		profileID:   DomainDecisionProfileID(ProfileApprovalDecision),
		profileVer:  ProfileVersionV1,
		context:     ctx,
		approval:    facts,
		hasApproval: true,
	}, operation.MechanicalFailure{}
}

// NewExceptionDecisionMaterial seals exception-decision/v1 material from the
// exact design §3.2 carrier-input type model.ExceptionDecisionFacts.
func NewExceptionDecisionMaterial(
	facts model.ExceptionDecisionFacts,
	ctx PublicationContext,
) (DomainDecisionMaterial, operation.MechanicalFailure) {
	if !facts.Sealed() {
		return DomainDecisionMaterial{}, operation.NewMechanicalFailure("exception_material_facts_unsealed")
	}
	if !ctx.sealed {
		return DomainDecisionMaterial{}, operation.NewMechanicalFailure("exception_material_context_unsealed")
	}
	return DomainDecisionMaterial{
		sealed:       true,
		profileID:    DomainDecisionProfileID(ProfileExceptionDecision),
		profileVer:   ProfileVersionV1,
		context:      ctx,
		exception:    facts,
		hasException: true,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether m was package-constructed.
func (m DomainDecisionMaterial) Sealed() bool { return m.sealed }

// ProfileID returns the sealed profile id.
func (m DomainDecisionMaterial) ProfileID() DomainDecisionProfileID { return m.profileID }

// ProfileVersion returns the sealed profile version pin.
func (m DomainDecisionMaterial) ProfileVersion() string { return m.profileVer }

// PublicationContext returns the sealed publication context.
func (m DomainDecisionMaterial) PublicationContext() PublicationContext { return m.context }

// ApprovalFacts returns approval facts when this is approval material.
func (m DomainDecisionMaterial) ApprovalFacts() (model.ApprovalDecisionFacts, bool) {
	if !m.sealed || !m.hasApproval {
		return model.ApprovalDecisionFacts{}, false
	}
	return m.approval, true
}

// ExceptionFacts returns exception facts when this is exception material.
func (m DomainDecisionMaterial) ExceptionFacts() (model.ExceptionDecisionFacts, bool) {
	if !m.sealed || !m.hasException {
		return model.ExceptionDecisionFacts{}, false
	}
	return m.exception, true
}
