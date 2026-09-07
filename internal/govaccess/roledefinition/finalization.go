package roledefinition

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

type FinalizedRoleDefinitionChange struct {
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

func (c FinalizedRoleDefinitionChange) ParticipantClaim() state.ParticipantClaim { return c.claim }
func (c FinalizedRoleDefinitionChange) ContextBinding() evidence.PublicationContextBinding {
	return c.ctxBinding
}
func (c FinalizedRoleDefinitionChange) PermitBinding() state.PermitBinding { return c.permitBind }
func (c FinalizedRoleDefinitionChange) PublicationKind() state.PublicationKind {
	return state.ResourceMutation
}
func (c FinalizedRoleDefinitionChange) EvidenceDescriptors() []evidence.MutationDescriptor {
	return append([]evidence.MutationDescriptor(nil), c.descriptors...)
}
func (c FinalizedRoleDefinitionChange) ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool) {
	return evidence.DomainConclusionDescriptor{}, false
}
func (c FinalizedRoleDefinitionChange) DomainDecisionMaterials() []evidence.DomainDecisionMaterial {
	return nil
}
func (c FinalizedRoleDefinitionChange) CanonicalResult() (operation.ResultMaterial, bool) {
	return c.result, c.sealed
}
func (c FinalizedRoleDefinitionChange) Sealed() bool { return c.sealed }
func (c FinalizedRoleDefinitionChange) ApplyTo(ed state.StateEditor) error {
	if !c.sealed {
		return errors.New("roledefinition: change_unsealed")
	}
	if ed == nil {
		return errors.New("roledefinition: editor_nil")
	}
	return ed.PutResource(c.resourceUID, c.version, append([]byte(nil), c.payload...))
}

func (i PreparedRoleDefinitionIntent) FinalizeAt(
	permit state.MutationFinalizationPermit,
) operation.FinalizationOutcome[FinalizedRoleDefinitionChange, DomainOutcome] {
	if !i.sealed {
		return operation.Failure[FinalizedRoleDefinitionChange, DomainOutcome](operation.NewMechanicalFailure("intent_unsealed"))
	}
	if !permit.Sealed() || !permit.ParticipantClaim().Equal(i.claim) {
		return operation.Failure[FinalizedRoleDefinitionChange, DomainOutcome](operation.NewMechanicalFailure("permit_mismatch"))
	}
	ctx := permit.PublicationContext()
	ctxBind, fail := evidence.BindPublicationContext(ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedRoleDefinitionChange, DomainOutcome](fail)
	}
	spec := i.spec
	spec.BaseClassification = effectiveClassification(spec)
	resource := model.RoleDefinition{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "RoleDefinition"},
		Metadata: apimeta.ObjectMeta{Name: i.resourceName, UID: i.resourceUID, ScopeRef: &i.scopeRef},
		Spec:     spec,
	}
	payload, err := json.Marshal(resource)
	if err != nil {
		return operation.Failure[FinalizedRoleDefinitionChange, DomainOutcome](operation.NewMechanicalFailure("resource_marshal_failed"))
	}
	after := sha256.Sum256(payload)
	facts, ok := model.NewMutationEventFacts(
		"roledefinition.published",
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: i.resourceName, UID: i.resourceUID},
		i.nextVersion,
		"roledefinition.write",
		after[:],
	)
	if !ok {
		return operation.Failure[FinalizedRoleDefinitionChange, DomainOutcome](operation.NewMechanicalFailure("facts_invalid"))
	}
	desc, fail := evidence.NewMutationDescriptor(facts, ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedRoleDefinitionChange, DomainOutcome](fail)
	}
	result, fail := operation.NewResultMaterial(payload)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedRoleDefinitionChange, DomainOutcome](fail)
	}
	return operation.Finalized[FinalizedRoleDefinitionChange, DomainOutcome](FinalizedRoleDefinitionChange{
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
