package membership

import (
	"crypto/sha256"
	"encoding/json"
	"errors"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

type FinalizedMembershipChange struct {
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

func (c FinalizedMembershipChange) Sealed() bool                             { return c.sealed }
func (c FinalizedMembershipChange) ParticipantClaim() state.ParticipantClaim { return c.claim }
func (c FinalizedMembershipChange) ContextBinding() evidence.PublicationContextBinding {
	return c.ctxBinding
}
func (c FinalizedMembershipChange) PermitBinding() state.PermitBinding { return c.permitBind }
func (c FinalizedMembershipChange) PublicationKind() state.PublicationKind {
	return state.ResourceMutation
}
func (c FinalizedMembershipChange) EvidenceDescriptors() []evidence.MutationDescriptor {
	return append([]evidence.MutationDescriptor(nil), c.descriptors...)
}
func (c FinalizedMembershipChange) ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool) {
	return evidence.DomainConclusionDescriptor{}, false
}
func (c FinalizedMembershipChange) DomainDecisionMaterials() []evidence.DomainDecisionMaterial {
	return nil
}
func (c FinalizedMembershipChange) CanonicalResult() (operation.ResultMaterial, bool) {
	return c.result, c.sealed
}
func (c FinalizedMembershipChange) ApplyTo(ed state.StateEditor) error {
	if !c.sealed {
		return errors.New("membership: change_unsealed")
	}
	if ed == nil {
		return errors.New("membership: editor_nil")
	}
	return ed.PutResource(c.resourceUID, c.version, append([]byte(nil), c.payload...))
}

func (i PreparedMembershipIntent) FinalizeAt(
	permit state.MutationFinalizationPermit,
) operation.FinalizationOutcome[FinalizedMembershipChange, DomainOutcome] {
	if !i.sealed {
		return operation.Failure[FinalizedMembershipChange, DomainOutcome](operation.NewMechanicalFailure("intent_unsealed"))
	}
	if !permit.Sealed() || !permit.ParticipantClaim().Equal(i.claim) {
		return operation.Failure[FinalizedMembershipChange, DomainOutcome](operation.NewMechanicalFailure("permit_mismatch"))
	}
	ctx := permit.PublicationContext()
	ctxBind, fail := evidence.BindPublicationContext(ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedMembershipChange, DomainOutcome](fail)
	}
	status := model.MembershipStatus{Phase: model.MembershipPhaseActive}
	if i.assignmentEffectSize == 0 && i.spec.ContextKind == model.MembershipContextAccessGroup {
		status.Phase = model.MembershipPhaseSuspended
	}
	resource := model.Membership{
		TypeMeta:  apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "Membership"},
		Metadata:  apimeta.ObjectMeta{Name: i.resourceName, UID: i.resourceUID, ScopeRef: &i.scopeRef},
		Spec:      i.spec,
		Protected: i.protected,
		Status:    status,
	}
	payload, err := json.Marshal(resource)
	if err != nil {
		return operation.Failure[FinalizedMembershipChange, DomainOutcome](operation.NewMechanicalFailure("resource_marshal_failed"))
	}
	after := sha256.Sum256(payload)
	facts, ok := model.NewMutationEventFacts(
		"membership.published",
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "Membership", Name: i.resourceName, UID: i.resourceUID},
		i.nextVersion,
		"membership.write",
		after[:],
	)
	if !ok {
		return operation.Failure[FinalizedMembershipChange, DomainOutcome](operation.NewMechanicalFailure("facts_invalid"))
	}
	desc, fail := evidence.NewMutationDescriptor(facts, ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedMembershipChange, DomainOutcome](fail)
	}
	result, fail := operation.NewResultMaterial(payload)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedMembershipChange, DomainOutcome](fail)
	}
	return operation.Finalized[FinalizedMembershipChange, DomainOutcome](FinalizedMembershipChange{
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
