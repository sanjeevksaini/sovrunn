package governanceprofile

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

type FinalizedGovernanceProfileChange struct {
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

func (c FinalizedGovernanceProfileChange) Sealed() bool                             { return c.sealed }
func (c FinalizedGovernanceProfileChange) ParticipantClaim() state.ParticipantClaim { return c.claim }
func (c FinalizedGovernanceProfileChange) ContextBinding() evidence.PublicationContextBinding {
	return c.ctxBinding
}
func (c FinalizedGovernanceProfileChange) PermitBinding() state.PermitBinding { return c.permitBind }
func (c FinalizedGovernanceProfileChange) PublicationKind() state.PublicationKind {
	return state.ResourceMutation
}
func (c FinalizedGovernanceProfileChange) EvidenceDescriptors() []evidence.MutationDescriptor {
	return append([]evidence.MutationDescriptor(nil), c.descriptors...)
}
func (c FinalizedGovernanceProfileChange) ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool) {
	return evidence.DomainConclusionDescriptor{}, false
}
func (c FinalizedGovernanceProfileChange) DomainDecisionMaterials() []evidence.DomainDecisionMaterial {
	return nil
}
func (c FinalizedGovernanceProfileChange) CanonicalResult() (operation.ResultMaterial, bool) {
	return c.result, c.sealed
}
func (c FinalizedGovernanceProfileChange) ApplyTo(ed state.StateEditor) error {
	if !c.sealed {
		return errors.New("governanceprofile: change_unsealed")
	}
	if ed == nil {
		return errors.New("governanceprofile: editor_nil")
	}
	return ed.PutResource(c.resourceUID, c.version, append([]byte(nil), c.payload...))
}

func (i PreparedGovernanceProfileIntent) FinalizeAt(
	permit state.MutationFinalizationPermit,
) operation.FinalizationOutcome[FinalizedGovernanceProfileChange, DomainOutcome] {
	if !i.sealed {
		return operation.Failure[FinalizedGovernanceProfileChange, DomainOutcome](operation.NewMechanicalFailure("intent_unsealed"))
	}
	if !permit.Sealed() || !permit.ParticipantClaim().Equal(i.claim) {
		return operation.Failure[FinalizedGovernanceProfileChange, DomainOutcome](operation.NewMechanicalFailure("permit_mismatch"))
	}
	ctx := permit.PublicationContext()
	ctxBind, fail := evidence.BindPublicationContext(ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedGovernanceProfileChange, DomainOutcome](fail)
	}
	resource := model.GovernanceProfile{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionGov, Kind: "GovernanceProfile"},
		Metadata: apimeta.ObjectMeta{Name: i.resourceName, UID: i.resourceUID, ScopeRef: &i.scopeRef},
		Spec:     i.spec,
	}
	payload, err := json.Marshal(resource)
	if err != nil {
		return operation.Failure[FinalizedGovernanceProfileChange, DomainOutcome](operation.NewMechanicalFailure("resource_marshal_failed"))
	}
	after := sha256.Sum256(payload)
	facts, ok := model.NewMutationEventFacts(
		"governanceprofile.published",
		apimeta.TypedRef{APIVersion: model.APIVersionGov, Kind: "GovernanceProfile", Name: i.resourceName, UID: i.resourceUID},
		i.nextVersion,
		"governanceprofile.write",
		after[:],
	)
	if !ok {
		return operation.Failure[FinalizedGovernanceProfileChange, DomainOutcome](operation.NewMechanicalFailure("facts_invalid"))
	}
	desc, fail := evidence.NewMutationDescriptor(facts, ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedGovernanceProfileChange, DomainOutcome](fail)
	}
	result, fail := operation.NewResultMaterial(payload)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedGovernanceProfileChange, DomainOutcome](fail)
	}
	return operation.Finalized[FinalizedGovernanceProfileChange, DomainOutcome](FinalizedGovernanceProfileChange{
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
