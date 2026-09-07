package accessgroup

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

type FinalizedAccessGroupChange struct {
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

func (c FinalizedAccessGroupChange) ParticipantClaim() state.ParticipantClaim {
	return c.claim
}
func (c FinalizedAccessGroupChange) ContextBinding() evidence.PublicationContextBinding {
	return c.ctxBinding
}
func (c FinalizedAccessGroupChange) PermitBinding() state.PermitBinding { return c.permitBind }
func (c FinalizedAccessGroupChange) PublicationKind() state.PublicationKind {
	return state.ResourceMutation
}
func (c FinalizedAccessGroupChange) EvidenceDescriptors() []evidence.MutationDescriptor {
	return append([]evidence.MutationDescriptor(nil), c.descriptors...)
}
func (c FinalizedAccessGroupChange) ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool) {
	return evidence.DomainConclusionDescriptor{}, false
}
func (c FinalizedAccessGroupChange) DomainDecisionMaterials() []evidence.DomainDecisionMaterial {
	return nil
}
func (c FinalizedAccessGroupChange) CanonicalResult() (operation.ResultMaterial, bool) {
	return c.result, c.sealed
}
func (c FinalizedAccessGroupChange) Sealed() bool { return c.sealed }

func (c FinalizedAccessGroupChange) ApplyTo(ed state.StateEditor) error {
	if !c.sealed {
		return errors.New("accessgroup: change_unsealed")
	}
	if ed == nil {
		return errors.New("accessgroup: editor_nil")
	}
	return ed.PutResource(c.resourceUID, c.version, append([]byte(nil), c.payload...))
}

func (i PreparedAccessGroupIntent) FinalizeAt(
	permit state.MutationFinalizationPermit,
) operation.FinalizationOutcome[FinalizedAccessGroupChange, DomainOutcome] {
	if !i.sealed {
		return operation.Failure[FinalizedAccessGroupChange, DomainOutcome](operation.NewMechanicalFailure("intent_unsealed"))
	}
	if !permit.Sealed() || !permit.ParticipantClaim().Equal(i.claim) {
		return operation.Failure[FinalizedAccessGroupChange, DomainOutcome](operation.NewMechanicalFailure("permit_mismatch"))
	}
	ctx := permit.PublicationContext()
	binding, fail := evidence.BindPublicationContext(ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedAccessGroupChange, DomainOutcome](fail)
	}
	resource := model.AccessGroup{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: model.KindAccessGroup},
		Metadata: apimeta.ObjectMeta{Name: i.resourceName, UID: i.resourceUID, ScopeRef: &i.scopeRef},
		Spec:     i.spec,
		Status:   model.AccessGroupStatus{Phase: model.AccessGroupPhaseActive},
	}
	payload, err := json.Marshal(resource)
	if err != nil {
		return operation.Failure[FinalizedAccessGroupChange, DomainOutcome](operation.NewMechanicalFailure("resource_marshal_failed"))
	}
	after := sha256.Sum256(payload)
	facts, ok := model.NewMutationEventFacts(
		"accessgroup.published",
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: model.KindAccessGroup, Name: i.resourceName, UID: i.resourceUID},
		i.nextVersion,
		"accessgroup.write",
		after[:],
	)
	if !ok {
		return operation.Failure[FinalizedAccessGroupChange, DomainOutcome](operation.NewMechanicalFailure("facts_invalid"))
	}
	desc, fail := evidence.NewMutationDescriptor(facts, ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedAccessGroupChange, DomainOutcome](fail)
	}
	result, fail := operation.NewResultMaterial(payload)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedAccessGroupChange, DomainOutcome](fail)
	}
	return operation.Finalized[FinalizedAccessGroupChange, DomainOutcome](FinalizedAccessGroupChange{
		sealed:      true,
		claim:       i.claim,
		ctxBinding:  binding,
		permitBind:  permit.Binding(),
		descriptors: []evidence.MutationDescriptor{desc},
		result:      result,
		resourceUID: i.resourceUID,
		version:     i.nextVersion,
		payload:     payload,
	})
}
