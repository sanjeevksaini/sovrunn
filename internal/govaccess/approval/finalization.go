package approval

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

type FinalizedApprovalChange struct {
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

func (c FinalizedApprovalChange) Sealed() bool                             { return c.sealed }
func (c FinalizedApprovalChange) ParticipantClaim() state.ParticipantClaim { return c.claim }
func (c FinalizedApprovalChange) ContextBinding() evidence.PublicationContextBinding {
	return c.ctxBinding
}
func (c FinalizedApprovalChange) PermitBinding() state.PermitBinding { return c.permitBind }
func (c FinalizedApprovalChange) PublicationKind() state.PublicationKind {
	return state.ResourceMutation
}
func (c FinalizedApprovalChange) EvidenceDescriptors() []evidence.MutationDescriptor {
	return append([]evidence.MutationDescriptor(nil), c.descriptors...)
}
func (c FinalizedApprovalChange) ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool) {
	return evidence.DomainConclusionDescriptor{}, false
}
func (c FinalizedApprovalChange) DomainDecisionMaterials() []evidence.DomainDecisionMaterial {
	return nil
}
func (c FinalizedApprovalChange) CanonicalResult() (operation.ResultMaterial, bool) {
	return c.result, c.sealed
}
func (c FinalizedApprovalChange) ApplyTo(ed state.StateEditor) error {
	if !c.sealed {
		return errors.New("approval: change_unsealed")
	}
	if ed == nil {
		return errors.New("approval: editor_nil")
	}
	return ed.PutResource(c.resourceUID, c.version, append([]byte(nil), c.payload...))
}

func (i PreparedApprovalIntent) FinalizeAt(
	permit state.MutationFinalizationPermit,
) operation.FinalizationOutcome[FinalizedApprovalChange, DomainOutcome] {
	if !i.sealed {
		return operation.Failure[FinalizedApprovalChange, DomainOutcome](operation.NewMechanicalFailure("intent_unsealed"))
	}
	if !permit.Sealed() || !permit.ParticipantClaim().Equal(i.claim) {
		return operation.Failure[FinalizedApprovalChange, DomainOutcome](operation.NewMechanicalFailure("permit_mismatch"))
	}
	// REQ-F18-13 immutable terminal approval: once terminal decision exists, deny mutation.
	if i.request.Status.Decision.Valid() {
		return operation.Domain[FinalizedApprovalChange](DomainOutcome{sealed: true, code: "terminal_immutable"})
	}
	ctx := permit.PublicationContext()
	_ = ctx.PublicationInstant()
	ctxBind, fail := evidence.BindPublicationContext(ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedApprovalChange, DomainOutcome](fail)
	}
	resource := i.request
	resource.TypeMeta = apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "ApprovalRequest"}
	resource.Metadata = apimeta.ObjectMeta{Name: i.resourceName, UID: i.resourceUID, ScopeRef: &i.scopeRef}
	payload, err := json.Marshal(resource)
	if err != nil {
		return operation.Failure[FinalizedApprovalChange, DomainOutcome](operation.NewMechanicalFailure("resource_marshal_failed"))
	}
	after := sha256.Sum256(payload)
	facts, ok := model.NewMutationEventFacts(
		"approvalrequest.published",
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "ApprovalRequest", Name: i.resourceName, UID: i.resourceUID},
		i.nextVersion,
		"approvalrequest.write",
		after[:],
	)
	if !ok {
		return operation.Failure[FinalizedApprovalChange, DomainOutcome](operation.NewMechanicalFailure("facts_invalid"))
	}
	desc, fail := evidence.NewMutationDescriptor(facts, ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedApprovalChange, DomainOutcome](fail)
	}
	result, fail := operation.NewResultMaterial(payload)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedApprovalChange, DomainOutcome](fail)
	}
	return operation.Finalized[FinalizedApprovalChange, DomainOutcome](FinalizedApprovalChange{
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
