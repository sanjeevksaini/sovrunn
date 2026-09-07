package privileged

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

type FinalizedPrivilegedChange struct {
	sealed      bool
	claim       state.ParticipantClaim
	ctxBinding  evidence.PublicationContextBinding
	permitBind  state.PermitBinding
	descriptors []evidence.MutationDescriptor
	result      operation.ResultMaterial
	resourceUID string
	version     string
	payload     []byte
}

func (c FinalizedPrivilegedChange) Sealed() bool                             { return c.sealed }
func (c FinalizedPrivilegedChange) ParticipantClaim() state.ParticipantClaim { return c.claim }
func (c FinalizedPrivilegedChange) ContextBinding() evidence.PublicationContextBinding {
	return c.ctxBinding
}
func (c FinalizedPrivilegedChange) PermitBinding() state.PermitBinding { return c.permitBind }
func (c FinalizedPrivilegedChange) PublicationKind() state.PublicationKind {
	return state.ResourceMutation
}
func (c FinalizedPrivilegedChange) EvidenceDescriptors() []evidence.MutationDescriptor {
	return append([]evidence.MutationDescriptor(nil), c.descriptors...)
}
func (c FinalizedPrivilegedChange) ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool) {
	return evidence.DomainConclusionDescriptor{}, false
}
func (c FinalizedPrivilegedChange) DomainDecisionMaterials() []evidence.DomainDecisionMaterial {
	return nil
}
func (c FinalizedPrivilegedChange) CanonicalResult() (operation.ResultMaterial, bool) {
	return c.result, c.sealed
}
func (c FinalizedPrivilegedChange) ApplyTo(ed state.StateEditor) error {
	if !c.sealed {
		return errors.New("privileged: change_unsealed")
	}
	if ed == nil {
		return errors.New("privileged: editor_nil")
	}
	return ed.PutResource(c.resourceUID, c.version, append([]byte(nil), c.payload...))
}

func (i PreparedPrivilegedIntent) FinalizeAt(
	permit state.MutationFinalizationPermit,
) operation.FinalizationOutcome[FinalizedPrivilegedChange, DomainOutcome] {
	if !i.sealed {
		return operation.Failure[FinalizedPrivilegedChange, DomainOutcome](operation.NewMechanicalFailure("intent_unsealed"))
	}
	if !permit.Sealed() || !permit.ParticipantClaim().Equal(i.claim) {
		return operation.Failure[FinalizedPrivilegedChange, DomainOutcome](operation.NewMechanicalFailure("permit_mismatch"))
	}
	pubInstant := permit.PublicationContext().PublicationInstant().UTC()
	if i.policyEvaluation.Outcome == "Deny" {
		return operation.Domain[FinalizedPrivilegedChange](DomainOutcome{sealed: true, code: "policy_denied"})
	}
	// REQ-F18-14/15 activation floor AAL2+ and phishing-resistant.
	if i.assurance.Level != model.AssuranceLevelAAL2 && i.assurance.Level != model.AssuranceLevelAAL3 {
		return operation.Domain[FinalizedPrivilegedChange](DomainOutcome{sealed: true, code: "aal2_required"})
	}
	if !i.assurance.PhishingResistant {
		return operation.Domain[FinalizedPrivilegedChange](DomainOutcome{sealed: true, code: "phishing_resistance_required"})
	}
	if pubInstant.Sub(i.assurance.AuthenticatedAt.UTC()) > 15*time.Minute {
		return operation.Domain[FinalizedPrivilegedChange](DomainOutcome{sealed: true, code: "assurance_too_old"})
	}
	ctx := permit.PublicationContext()
	ctxBind, fail := evidence.BindPublicationContext(ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedPrivilegedChange, DomainOutcome](fail)
	}
	resource := model.PrivilegedAccessRequest{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "PrivilegedAccessRequest"},
		Metadata: apimeta.ObjectMeta{Name: i.resourceName, UID: i.resourceUID, ScopeRef: &i.scopeRef},
		Spec:     i.spec,
		Status:   model.PrivilegedAccessRequestStatus{Phase: "PendingApproval"},
	}
	payload, err := json.Marshal(resource)
	if err != nil {
		return operation.Failure[FinalizedPrivilegedChange, DomainOutcome](operation.NewMechanicalFailure("resource_marshal_failed"))
	}
	after := sha256.Sum256(payload)
	facts, ok := model.NewMutationEventFacts(
		"privilegedaccessrequest.published",
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "PrivilegedAccessRequest", Name: i.resourceName, UID: i.resourceUID},
		i.nextVersion,
		"privilegedaccessrequest.write",
		after[:],
	)
	if !ok {
		return operation.Failure[FinalizedPrivilegedChange, DomainOutcome](operation.NewMechanicalFailure("facts_invalid"))
	}
	desc, fail := evidence.NewMutationDescriptor(facts, ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedPrivilegedChange, DomainOutcome](fail)
	}
	result, fail := operation.NewResultMaterial(payload)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedPrivilegedChange, DomainOutcome](fail)
	}
	return operation.Finalized[FinalizedPrivilegedChange, DomainOutcome](FinalizedPrivilegedChange{
		sealed:      true,
		claim:       i.claim,
		ctxBinding:  ctxBind,
		permitBind:  permit.Binding(),
		descriptors: []evidence.MutationDescriptor{desc},
		result:      result,
		resourceUID: i.resourceUID,
		version:     i.nextVersion,
		payload:     payload,
	})
}
