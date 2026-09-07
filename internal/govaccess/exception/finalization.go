package exception

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

type FinalizedExceptionChange struct {
	sealed      bool
	claim       state.ParticipantClaim
	ctxBinding  evidence.PublicationContextBinding
	permitBind  state.PermitBinding
	kind        state.PublicationKind
	descriptors []evidence.MutationDescriptor
	conclusion  evidence.DomainConclusionDescriptor
	hasConc     bool
	materials   []evidence.DomainDecisionMaterial
	result      operation.ResultMaterial
	hasResult   bool
	resourceUID string
	version     string
	payload     []byte
}

func (c FinalizedExceptionChange) Sealed() bool                             { return c.sealed }
func (c FinalizedExceptionChange) ParticipantClaim() state.ParticipantClaim { return c.claim }
func (c FinalizedExceptionChange) ContextBinding() evidence.PublicationContextBinding {
	return c.ctxBinding
}
func (c FinalizedExceptionChange) PermitBinding() state.PermitBinding     { return c.permitBind }
func (c FinalizedExceptionChange) PublicationKind() state.PublicationKind { return c.kind }
func (c FinalizedExceptionChange) EvidenceDescriptors() []evidence.MutationDescriptor {
	return append([]evidence.MutationDescriptor(nil), c.descriptors...)
}
func (c FinalizedExceptionChange) ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool) {
	return c.conclusion, c.hasConc
}
func (c FinalizedExceptionChange) DomainDecisionMaterials() []evidence.DomainDecisionMaterial {
	return append([]evidence.DomainDecisionMaterial(nil), c.materials...)
}
func (c FinalizedExceptionChange) CanonicalResult() (operation.ResultMaterial, bool) {
	return c.result, c.hasResult
}
func (c FinalizedExceptionChange) ApplyTo(ed state.StateEditor) error {
	if !c.sealed {
		return errors.New("exception: change_unsealed")
	}
	if ed == nil {
		return errors.New("exception: editor_nil")
	}
	if c.kind == state.DomainConclusion {
		return nil
	}
	return ed.PutResource(c.resourceUID, c.version, append([]byte(nil), c.payload...))
}

func (i PreparedExceptionIntent) FinalizeAt(
	permit state.MutationFinalizationPermit,
) operation.FinalizationOutcome[FinalizedExceptionChange, DomainOutcome] {
	if !i.sealed {
		return operation.Failure[FinalizedExceptionChange, DomainOutcome](operation.NewMechanicalFailure("intent_unsealed"))
	}
	if !permit.Sealed() || !permit.ParticipantClaim().Equal(i.claim) {
		return operation.Failure[FinalizedExceptionChange, DomainOutcome](operation.NewMechanicalFailure("permit_mismatch"))
	}
	ctx := permit.PublicationContext()
	ctxBind, fail := evidence.BindPublicationContext(ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedExceptionChange, DomainOutcome](fail)
	}
	material, fail := evidence.NewExceptionDecisionMaterial(i.decisionFacts, ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedExceptionChange, DomainOutcome](fail)
	}

	change := FinalizedExceptionChange{
		sealed:     true,
		claim:      i.claim,
		ctxBinding: ctxBind,
		permitBind: permit.Binding(),
		materials:  []evidence.DomainDecisionMaterial{material},
	}

	if i.decisionFacts.Result() == model.ExceptionTerminalDeny {
		conclusion, fail := evidence.NewDomainConclusionDescriptor(i.decisionFacts, ctx)
		if fail.Reason() != "" {
			return operation.Failure[FinalizedExceptionChange, DomainOutcome](fail)
		}
		if fail := evidence.ValidateExceptionTerminalDescriptorSet(i.decisionFacts, nil, &conclusion, change.materials, ctx); fail.Reason() != "" {
			return operation.Failure[FinalizedExceptionChange, DomainOutcome](fail)
		}
		change.kind = state.DomainConclusion
		change.conclusion = conclusion
		change.hasConc = true
		return operation.Finalized[FinalizedExceptionChange, DomainOutcome](change)
	}

	resource := model.ExceptionGrant{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionGov, Kind: model.KindExceptionGrant},
		Metadata: apimeta.ObjectMeta{Name: i.resourceName, UID: i.resourceUID, ScopeRef: &i.scopeRef},
		Record:   i.record,
	}
	payload, err := json.Marshal(resource)
	if err != nil {
		return operation.Failure[FinalizedExceptionChange, DomainOutcome](operation.NewMechanicalFailure("resource_marshal_failed"))
	}
	after := sha256.Sum256(payload)
	factsDecided, ok := model.NewMutationEventFacts(
		evidence.EventTypeExceptionProposalDecided,
		apimeta.TypedRef{APIVersion: model.APIVersionGov, Kind: model.KindExceptionGrant, Name: i.resourceName, UID: i.resourceUID},
		i.nextVersion,
		"exceptionproposal.decide",
		after[:],
	)
	if !ok {
		return operation.Failure[FinalizedExceptionChange, DomainOutcome](operation.NewMechanicalFailure("decided_facts_invalid"))
	}
	decided, fail := evidence.NewMutationDescriptor(factsDecided, ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedExceptionChange, DomainOutcome](fail)
	}
	factsIssued, ok := model.NewMutationEventFacts(
		evidence.EventTypeExceptionGrantIssued,
		apimeta.TypedRef{APIVersion: model.APIVersionGov, Kind: model.KindExceptionGrant, Name: i.resourceName, UID: i.resourceUID},
		i.nextVersion,
		"exceptiongrant.write",
		after[:],
	)
	if !ok {
		return operation.Failure[FinalizedExceptionChange, DomainOutcome](operation.NewMechanicalFailure("issued_facts_invalid"))
	}
	issued, fail := evidence.NewMutationDescriptor(factsIssued, ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedExceptionChange, DomainOutcome](fail)
	}
	descs := evidence.SortMutationDescriptors([]evidence.MutationDescriptor{decided, issued})
	if fail := evidence.ValidateExceptionTerminalDescriptorSet(i.decisionFacts, descs, nil, change.materials, ctx); fail.Reason() != "" {
		return operation.Failure[FinalizedExceptionChange, DomainOutcome](fail)
	}
	result, fail := operation.NewResultMaterial(payload)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedExceptionChange, DomainOutcome](fail)
	}
	change.kind = state.ResourceMutation
	change.descriptors = descs
	change.result = result
	change.hasResult = true
	change.resourceUID = i.resourceUID
	change.version = i.nextVersion
	change.payload = payload
	return operation.Finalized[FinalizedExceptionChange, DomainOutcome](change)
}
