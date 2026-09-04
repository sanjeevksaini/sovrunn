package roleassign

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

// FinalizedRoleAssignmentChange is the sealed, permit-bound RoleAssignment
// publication change (DD-01, DD-13 Rule 009). Only FinalizeAt constructs it.
type FinalizedRoleAssignmentChange struct {
	sealed      bool
	claim       state.ParticipantClaim
	ctxBinding  evidence.PublicationContextBinding
	permitBind  state.PermitBinding
	kind        state.PublicationKind
	descriptors []evidence.MutationDescriptor
	result      operation.ResultMaterial
	hasResult   bool
	resourceUID string
	version     string
	payload     []byte
}

// ParticipantClaim returns the sealed participant claim.
func (c FinalizedRoleAssignmentChange) ParticipantClaim() state.ParticipantClaim {
	return c.claim
}

// ContextBinding returns the sealed publication-context binding.
func (c FinalizedRoleAssignmentChange) ContextBinding() evidence.PublicationContextBinding {
	return c.ctxBinding
}

// PermitBinding returns the unforgeable permit binding embedded at FinalizeAt.
func (c FinalizedRoleAssignmentChange) PermitBinding() state.PermitBinding {
	return c.permitBind
}

// PublicationKind returns ResourceMutation for RoleAssignment effects.
func (c FinalizedRoleAssignmentChange) PublicationKind() state.PublicationKind {
	return c.kind
}

// EvidenceDescriptors returns the canonical ordered mutation descriptors.
func (c FinalizedRoleAssignmentChange) EvidenceDescriptors() []evidence.MutationDescriptor {
	return append([]evidence.MutationDescriptor(nil), c.descriptors...)
}

// ConclusionDescriptor is absent for RoleAssignment ResourceMutation changes.
func (c FinalizedRoleAssignmentChange) ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool) {
	return evidence.DomainConclusionDescriptor{}, false
}

// DomainDecisionMaterials is empty for RoleAssignment (approval/exception only).
func (c FinalizedRoleAssignmentChange) DomainDecisionMaterials() []evidence.DomainDecisionMaterial {
	return nil
}

// CanonicalResult returns the originating ResultMaterial when present.
func (c FinalizedRoleAssignmentChange) CanonicalResult() (operation.ResultMaterial, bool) {
	if !c.sealed || !c.hasResult {
		return operation.ResultMaterial{}, false
	}
	return c.result, true
}

// ApplyTo applies the RoleAssignment resource mutation through the private editor.
func (c FinalizedRoleAssignmentChange) ApplyTo(ed state.StateEditor) error {
	if !c.sealed {
		return errors.New("roleassign: finalized_change_unsealed")
	}
	if ed == nil {
		return errors.New("roleassign: state_editor_nil")
	}
	return ed.PutResource(c.resourceUID, c.version, append([]byte(nil), c.payload...))
}

// Sealed reports whether c was constructed through FinalizeAt.
func (c FinalizedRoleAssignmentChange) Sealed() bool { return c.sealed }

// FinalizeAt performs publication-time RoleAssignment finalization using only
// the permit's PublicationContext (DD-01, design §5.4). It never reads clock,
// state, locks, or external I/O.
func (i PreparedRoleAssignmentIntent) FinalizeAt(
	permit state.MutationFinalizationPermit,
) operation.FinalizationOutcome[FinalizedRoleAssignmentChange, DomainOutcome] {
	if !i.sealed {
		return operation.Failure[FinalizedRoleAssignmentChange, DomainOutcome](
			operation.NewMechanicalFailure("intent_unsealed"),
		)
	}
	if !permit.Sealed() {
		return operation.Failure[FinalizedRoleAssignmentChange, DomainOutcome](
			operation.NewMechanicalFailure("permit_unsealed"),
		)
	}
	if !permit.ParticipantClaim().Equal(i.claim) {
		return operation.Failure[FinalizedRoleAssignmentChange, DomainOutcome](
			operation.NewMechanicalFailure("permit_claim_mismatch"),
		)
	}
	ctx := permit.PublicationContext()
	if !ctx.Sealed() {
		return operation.Failure[FinalizedRoleAssignmentChange, DomainOutcome](
			operation.NewMechanicalFailure("publication_context_unsealed"),
		)
	}
	pubInstant := ctx.PublicationInstant().UTC()

	resource, eventType, action, domain, fail := i.deriveFinalResource(pubInstant)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedRoleAssignmentChange, DomainOutcome](fail)
	}
	if domain.Sealed() {
		return operation.Domain[FinalizedRoleAssignmentChange](domain)
	}

	payload, err := json.Marshal(resource)
	if err != nil {
		return operation.Failure[FinalizedRoleAssignmentChange, DomainOutcome](
			operation.NewMechanicalFailure("final_resource_marshal"),
		)
	}
	afterDigest := sha256.Sum256(payload)

	facts, ok := model.NewMutationEventFacts(
		eventType,
		apimeta.TypedRef{
			APIVersion: model.APIVersionIAM,
			Kind:       KindRoleAssignment,
			Name:       i.resourceName,
			UID:        i.resourceUID,
		},
		i.nextVersion,
		action,
		afterDigest[:],
	)
	if !ok {
		return operation.Failure[FinalizedRoleAssignmentChange, DomainOutcome](
			operation.NewMechanicalFailure("mutation_facts_invalid"),
		)
	}
	desc, fail := evidence.NewMutationDescriptor(facts, ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedRoleAssignmentChange, DomainOutcome](fail)
	}
	ctxBind, fail := evidence.BindPublicationContext(ctx)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedRoleAssignmentChange, DomainOutcome](fail)
	}
	result, fail := operation.NewResultMaterial(payload)
	if fail.Reason() != "" {
		return operation.Failure[FinalizedRoleAssignmentChange, DomainOutcome](fail)
	}

	change := FinalizedRoleAssignmentChange{
		sealed:      true,
		claim:       i.claim,
		ctxBinding:  ctxBind,
		permitBind:  permit.Binding(),
		kind:        state.ResourceMutation,
		descriptors: []evidence.MutationDescriptor{desc},
		result:      result,
		hasResult:   true,
		resourceUID: i.resourceUID,
		version:     i.nextVersion,
		payload:     payload,
	}
	return operation.Finalized[FinalizedRoleAssignmentChange, DomainOutcome](change)
}

func (i PreparedRoleAssignmentIntent) deriveFinalResource(
	pubInstant time.Time,
) (model.RoleAssignment, string, string, DomainOutcome, operation.MechanicalFailure) {
	meta := apimeta.ObjectMeta{
		Name: i.resourceName,
		UID:  i.resourceUID,
	}
	if i.scopeRef.Kind != "" || i.scopeRef.Name != "" || i.scopeRef.UID != "" {
		scope := i.scopeRef
		meta.ScopeRef = &scope
	}

	switch i.operation {
	case OpGrant, OpActivate:
		if !i.hasSpec {
			return model.RoleAssignment{}, "", "", DomainOutcome{}, operation.NewMechanicalFailure("grant_spec_missing")
		}
		spec := copyRoleAssignmentSpec(i.spec)
		phase, due, domain, fail := deriveGrantStatus(spec, pubInstant, i.nextReviewDueAt)
		if fail.Reason() != "" {
			return model.RoleAssignment{}, "", "", DomainOutcome{}, fail
		}
		if domain.Sealed() {
			return model.RoleAssignment{}, "", "", domain, operation.MechanicalFailure{}
		}
		return model.RoleAssignment{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: KindRoleAssignment},
			Metadata: meta,
			Spec:     spec,
			Status:   model.RoleAssignmentStatus{Phase: phase, NextReviewDueAt: due},
		}, "roleassignment.created", string(i.operation), DomainOutcome{}, operation.MechanicalFailure{}

	case OpReviewReplace:
		if !i.hasReplacement {
			return model.RoleAssignment{}, "", "", DomainOutcome{}, operation.NewMechanicalFailure("replace_spec_missing")
		}
		spec := copyRoleAssignmentSpec(i.replacementSpec)
		phase, due, domain, fail := deriveGrantStatus(spec, pubInstant, i.nextReviewDueAt)
		if fail.Reason() != "" {
			return model.RoleAssignment{}, "", "", DomainOutcome{}, fail
		}
		if domain.Sealed() {
			return model.RoleAssignment{}, "", "", domain, operation.MechanicalFailure{}
		}
		return model.RoleAssignment{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: KindRoleAssignment},
			Metadata: meta,
			Spec:     spec,
			Status:   model.RoleAssignmentStatus{Phase: phase, NextReviewDueAt: due},
		}, "roleassignment.created", string(i.operation), DomainOutcome{}, operation.MechanicalFailure{}

	case OpPrivilegedRevoke, OpReviewRevoke, OpDirectRevoke:
		return model.RoleAssignment{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: KindRoleAssignment},
			Metadata: meta,
			Spec:     copyRoleAssignmentSpec(i.spec),
			Status:   model.RoleAssignmentStatus{Phase: "Revoked"},
		}, "roleassignment.revoked", string(i.operation), DomainOutcome{}, operation.MechanicalFailure{}

	case OpAdvanceReviewDue:
		if i.nextReviewDueAt == nil {
			return model.RoleAssignment{}, "", "", DomainOutcome{}, operation.NewMechanicalFailure("review_due_missing")
		}
		due := i.nextReviewDueAt.UTC()
		if !due.After(pubInstant) {
			return model.RoleAssignment{}, "", "", newDomainOutcome("review_due_not_after_publication"), operation.MechanicalFailure{}
		}
		return model.RoleAssignment{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: KindRoleAssignment},
			Metadata: meta,
			Spec:     copyRoleAssignmentSpec(i.spec),
			Status: model.RoleAssignmentStatus{
				Phase:           "Effective",
				NextReviewDueAt: &due,
			},
		}, "roleassignment.review-due-advanced", string(i.operation), DomainOutcome{}, operation.MechanicalFailure{}

	default:
		return model.RoleAssignment{}, "", "", DomainOutcome{}, operation.NewMechanicalFailure("operation_unsupported")
	}
}

func deriveGrantStatus(
	spec model.RoleAssignmentSpec,
	pubInstant time.Time,
	reviewDue *time.Time,
) (phase string, due *time.Time, domain DomainOutcome, fail operation.MechanicalFailure) {
	if !spec.Validity.Valid() {
		return "", nil, DomainOutcome{}, operation.NewMechanicalFailure("validity_invalid")
	}
	switch spec.Validity {
	case model.AssignmentValidityTimeBound:
		if spec.NotBefore == nil || spec.ExpiresAt == nil {
			return "", nil, DomainOutcome{}, operation.NewMechanicalFailure("timebound_timestamps_missing")
		}
		nb := spec.NotBefore.UTC()
		exp := spec.ExpiresAt.UTC()
		if !exp.After(nb) {
			return "", nil, DomainOutcome{}, operation.NewMechanicalFailure("timebound_window_invalid")
		}
		if !pubInstant.Before(exp) {
			return "", nil, newDomainOutcome("assignment_expired"), operation.MechanicalFailure{}
		}
		if pubInstant.Before(nb) {
			phase = "NotYetValid"
		} else {
			phase = "Effective"
		}
	case model.AssignmentValidityStanding:
		phase = "Effective"
		if reviewDue != nil {
			t := reviewDue.UTC()
			if !t.After(pubInstant) {
				return "", nil, newDomainOutcome("review_due_not_after_publication"), operation.MechanicalFailure{}
			}
			due = &t
		}
	}
	return phase, due, DomainOutcome{}, operation.MechanicalFailure{}
}
