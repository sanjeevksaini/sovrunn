package evidence

import (
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// PreparedEvidenceChange is the sealed evidence change accepted by a state
// transaction before Seal/Commit (VS0-WRITER-011).
type PreparedEvidenceChange struct {
	sealed  bool
	context PublicationContext
	set     CarrierSet
	kind    preparedKind
}

type preparedKind uint8

const (
	preparedNone preparedKind = iota
	preparedAuthorizedMutation
	preparedAutomaticMutation
	preparedAuthorizationOnly
	preparedDomainConclusion
)

// Sealed reports whether c was package-constructed.
func (c PreparedEvidenceChange) Sealed() bool { return c.sealed }

// PublicationContext returns the sealed publication context.
func (c PreparedEvidenceChange) PublicationContext() PublicationContext { return c.context }

// CarrierSet returns the sealed carrier set.
func (c PreparedEvidenceChange) CarrierSet() CarrierSet { return c.set }

// PrepareMutationCarrierSetAuthorized constructs PreparedEvidenceChange from an
// AuthorizedMutationPlan.
func PrepareMutationCarrierSetAuthorized(plan AuthorizedMutationPlan) (PreparedEvidenceChange, operation.MechanicalFailure) {
	if !plan.sealed {
		return PreparedEvidenceChange{}, operation.NewMechanicalFailure("prepare_mutation_authorized_unsealed")
	}
	set, fail := buildAuthorizedCarrierSet(plan)
	if fail.Reason() != "" {
		return PreparedEvidenceChange{}, fail
	}
	return PreparedEvidenceChange{
		sealed:  true,
		context: plan.context,
		set:     set,
		kind:    preparedAuthorizedMutation,
	}, operation.MechanicalFailure{}
}

// PrepareMutationCarrierSetAutomatic constructs PreparedEvidenceChange from an AutomaticMutationPlan.
func PrepareMutationCarrierSetAutomatic(plan AutomaticMutationPlan) (PreparedEvidenceChange, operation.MechanicalFailure) {
	if !plan.sealed {
		return PreparedEvidenceChange{}, operation.NewMechanicalFailure("prepare_mutation_automatic_unsealed")
	}
	set, fail := buildAutomaticCarrierSet(plan)
	if fail.Reason() != "" {
		return PreparedEvidenceChange{}, fail
	}
	return PreparedEvidenceChange{
		sealed:  true,
		context: plan.context,
		set:     set,
		kind:    preparedAutomaticMutation,
	}, operation.MechanicalFailure{}
}

// PrepareMutationCarrierSet dispatches authorized or automatic mutation plans.
// Exactly one sealed plan variant may be supplied.
func PrepareMutationCarrierSet(
	authorized *AuthorizedMutationPlan,
	automatic *AutomaticMutationPlan,
) (PreparedEvidenceChange, operation.MechanicalFailure) {
	switch {
	case authorized != nil && automatic == nil:
		return PrepareMutationCarrierSetAuthorized(*authorized)
	case automatic != nil && authorized == nil:
		return PrepareMutationCarrierSetAutomatic(*automatic)
	default:
		return PreparedEvidenceChange{}, operation.NewMechanicalFailure("prepare_mutation_plan_ambiguous")
	}
}

// PrepareAuthorizationCarrierSet constructs evidence-only denied/read/replay carriers.
func PrepareAuthorizationCarrierSet(
	evaluation CompletedAuthorizationEvaluation,
) (PreparedEvidenceChange, operation.MechanicalFailure) {
	if !evaluation.sealed {
		return PreparedEvidenceChange{}, operation.NewMechanicalFailure("prepare_authz_evaluation_unsealed")
	}
	set := CarrierSet{
		sealed: true,
		carriers: []CarrierRef{
			{Kind: CarrierAuditEvent, ID: carrierID(evaluation.context, "authorization.evaluated", 0)},
			{Kind: CarrierDecisionRecord, ID: carrierID(evaluation.context, ProfileAuthorizationDecision, 1)},
		},
		decisions: []decision.ProfileRef{
			{Name: ProfileAuthorizationDecision, Version: ProfileVersionV1},
		},
	}
	return PreparedEvidenceChange{
		sealed:  true,
		context: evaluation.context,
		set:     set,
		kind:    preparedAuthorizationOnly,
	}, operation.MechanicalFailure{}
}

// PrepareDomainConclusionCarrierSet constructs carriers for DomainConclusionPlan.
func PrepareDomainConclusionCarrierSet(plan DomainConclusionPlan) (PreparedEvidenceChange, operation.MechanicalFailure) {
	if !plan.sealed {
		return PreparedEvidenceChange{}, operation.NewMechanicalFailure("prepare_conclusion_unsealed")
	}
	set, fail := buildDomainConclusionCarrierSet(plan)
	if fail.Reason() != "" {
		return PreparedEvidenceChange{}, fail
	}
	return PreparedEvidenceChange{
		sealed:  true,
		context: plan.context,
		set:     set,
		kind:    preparedDomainConclusion,
	}, operation.MechanicalFailure{}
}

// MutationDigest is a canonical aggregate mutation digest derived from a plan.
type MutationDigest struct {
	sealed bool
	value  []byte
}

// NewMutationDigestFromProofs derives a sealed mutation digest.
func NewMutationDigestFromProofs(proofs []AppliedMutationProof) (MutationDigest, operation.MechanicalFailure) {
	if len(proofs) == 0 {
		return MutationDigest{}, operation.NewMechanicalFailure("mutation_digest_empty")
	}
	var parts [][]byte
	for _, p := range proofs {
		if !p.sealed {
			return MutationDigest{}, operation.NewMechanicalFailure("mutation_digest_proof_unsealed")
		}
		for _, d := range p.descriptors {
			parts = append(parts, []byte(descriptorOrderKey(d)))
			parts = append(parts, d.facts.AfterStateDigest())
		}
	}
	return MutationDigest{sealed: true, value: sha256Concat(parts...)}, operation.MechanicalFailure{}
}

// Sealed reports whether d was package-constructed.
func (d MutationDigest) Sealed() bool { return d.sealed }

// Bytes returns a defensive copy.
func (d MutationDigest) Bytes() []byte { return copyBytes(d.value) }

// CompletionBindingForPlan builds a neutral operation.CompletionBinding from a
// validated mutation digest and publication context (evidence -> operation edge).
func CompletionBindingForPlan(ctx PublicationContext, digest MutationDigest) (operation.CompletionBinding, operation.MechanicalFailure) {
	if !ctx.sealed {
		return operation.CompletionBinding{}, operation.NewMechanicalFailure("completion_binding_context_unsealed")
	}
	if !digest.sealed {
		return operation.CompletionBinding{}, operation.NewMechanicalFailure("completion_binding_digest_unsealed")
	}
	pub, fail := operation.NewPublicationBinding(sha256Concat([]byte("pub"), []byte(formatUint(ctx.publicationSequence))))
	if fail.Reason() != "" {
		return operation.CompletionBinding{}, fail
	}
	md, fail := operation.NewMutationDigest(digest.value)
	if fail.Reason() != "" {
		return operation.CompletionBinding{}, fail
	}
	return operation.NewCompletionBinding(pub, md)
}
