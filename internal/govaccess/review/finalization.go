package review

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

type FinalizedReviewChange struct {
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

func (c FinalizedReviewChange) Sealed() bool                             { return c.sealed }
func (c FinalizedReviewChange) ParticipantClaim() state.ParticipantClaim { return c.claim }
func (c FinalizedReviewChange) ContextBinding() evidence.PublicationContextBinding {
	return c.ctxBinding
}
func (c FinalizedReviewChange) PermitBinding() state.PermitBinding { return c.permitBind }
func (c FinalizedReviewChange) PublicationKind() state.PublicationKind {
	return state.ResourceMutation
}
func (c FinalizedReviewChange) EvidenceDescriptors() []evidence.MutationDescriptor {
	return append([]evidence.MutationDescriptor(nil), c.descriptors...)
}
func (c FinalizedReviewChange) ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool) {
	return evidence.DomainConclusionDescriptor{}, false
}
func (c FinalizedReviewChange) DomainDecisionMaterials() []evidence.DomainDecisionMaterial {
	return nil
}
func (c FinalizedReviewChange) CanonicalResult() (operation.ResultMaterial, bool) {
	return c.result, c.sealed
}
func (c FinalizedReviewChange) ApplyTo(ed state.StateEditor) error {
	if !c.sealed {
		return errors.New("review: change_unsealed")
	}
	if ed == nil {
		return errors.New("review: editor_nil")
	}
	return ed.PutResource(c.resourceUID, c.version, append([]byte(nil), c.payload...))
}

func (i PreparedReviewIntent) FinalizeAt(
	permit state.MutationFinalizationPermit,
) operation.FinalizationOutcome[FinalizedReviewChange, DomainOutcome] {
	if !i.sealed {
		return operation.Failure[FinalizedReviewChange, DomainOutcome](operation.NewMechanicalFailure("intent_unsealed"))
	}
	if !permit.Sealed() || !permit.ParticipantClaim().Equal(i.claim) {
		return operation.Failure[FinalizedReviewChange, DomainOutcome](operation.NewMechanicalFailure("permit_mismatch"))
	}
	for _, elig := range i.spec.ReviewerEligibility {
		if hasBeneficiaryConflict(elig.Principal, i.beneficiarySet) {
			return operation.Domain[FinalizedReviewChange](DomainOutcome{sealed: true, code: "reviewer_conflict"})
		}
	}
	ctx := permit.PublicationContext()
	ctxBind, fail := evidence.BindPublicationContext(ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedReviewChange, DomainOutcome](fail)
	}
	resource := model.AccessReview{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "AccessReview"},
		Metadata: apimeta.ObjectMeta{Name: i.resourceName, UID: i.resourceUID, ScopeRef: &i.scopeRef},
		Spec:     i.spec,
		Status:   model.AccessReviewStatus{Phase: "Open"},
	}
	payload, err := json.Marshal(resource)
	if err != nil {
		return operation.Failure[FinalizedReviewChange, DomainOutcome](operation.NewMechanicalFailure("resource_marshal_failed"))
	}
	after := sha256.Sum256(payload)
	facts, ok := model.NewMutationEventFacts(
		"accessreview.published",
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "AccessReview", Name: i.resourceName, UID: i.resourceUID},
		i.nextVersion,
		"accessreview.write",
		after[:],
	)
	if !ok {
		return operation.Failure[FinalizedReviewChange, DomainOutcome](operation.NewMechanicalFailure("facts_invalid"))
	}
	desc, fail := evidence.NewMutationDescriptor(facts, ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedReviewChange, DomainOutcome](fail)
	}
	result, fail := operation.NewResultMaterial(payload)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedReviewChange, DomainOutcome](fail)
	}
	return operation.Finalized[FinalizedReviewChange, DomainOutcome](FinalizedReviewChange{
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
