package evidence

import (
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// AppliedMutationProof is a sealed proof of a ResourceMutation application.
type AppliedMutationProof struct {
	sealed      bool
	link        operation.AppliedChangeLink
	descriptors []MutationDescriptor
	materials   []DomainDecisionMaterial
	context     PublicationContext
}

// BindAppliedChange seals an applied mutation proof from a neutral link.
func BindAppliedChange(
	link operation.AppliedChangeLink,
	descriptors []MutationDescriptor,
	materials []DomainDecisionMaterial,
) (AppliedMutationProof, operation.MechanicalFailure) {
	if !link.Sealed() {
		return AppliedMutationProof{}, operation.NewMechanicalFailure("applied_mutation_link_unsealed")
	}
	if len(descriptors) == 0 {
		return AppliedMutationProof{}, operation.NewMechanicalFailure("applied_mutation_descriptors_empty")
	}
	ordered := SortMutationDescriptors(descriptors)
	ctx := ordered[0].context
	if !ctx.sealed {
		return AppliedMutationProof{}, operation.NewMechanicalFailure("applied_mutation_context_unsealed")
	}
	seen := map[string]bool{}
	for _, d := range ordered {
		if !d.sealed || !d.context.Equal(ctx) {
			return AppliedMutationProof{}, operation.NewMechanicalFailure("applied_mutation_descriptor_invalid")
		}
		key := descriptorOrderKey(d)
		if seen[key] {
			return AppliedMutationProof{}, operation.NewMechanicalFailure("applied_mutation_descriptor_duplicate")
		}
		seen[key] = true
	}
	mats := append([]DomainDecisionMaterial(nil), materials...)
	for _, m := range mats {
		if !m.sealed || !m.context.Equal(ctx) {
			return AppliedMutationProof{}, operation.NewMechanicalFailure("applied_mutation_material_invalid")
		}
	}
	return AppliedMutationProof{
		sealed:      true,
		link:        link,
		descriptors: ordered,
		materials:   mats,
		context:     ctx,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether p was package-constructed.
func (p AppliedMutationProof) Sealed() bool { return p.sealed }

// Descriptors returns a copy of sealed mutation descriptors.
func (p AppliedMutationProof) Descriptors() []MutationDescriptor {
	return append([]MutationDescriptor(nil), p.descriptors...)
}

// Materials returns a copy of sealed domain decision materials.
func (p AppliedMutationProof) Materials() []DomainDecisionMaterial {
	return append([]DomainDecisionMaterial(nil), p.materials...)
}

// PublicationContext returns the sealed publication context.
func (p AppliedMutationProof) PublicationContext() PublicationContext { return p.context }

// AppliedConclusionProof is a sealed proof of a DomainConclusion application.
type AppliedConclusionProof struct {
	sealed     bool
	link       operation.AppliedConclusionLink
	descriptor DomainConclusionDescriptor
	materials  []DomainDecisionMaterial
	context    PublicationContext
}

// BindAppliedConclusion seals an applied conclusion proof.
func BindAppliedConclusion(
	link operation.AppliedConclusionLink,
	descriptor DomainConclusionDescriptor,
	materials []DomainDecisionMaterial,
) (AppliedConclusionProof, operation.MechanicalFailure) {
	if !link.Sealed() {
		return AppliedConclusionProof{}, operation.NewMechanicalFailure("applied_conclusion_link_unsealed")
	}
	if !descriptor.sealed {
		return AppliedConclusionProof{}, operation.NewMechanicalFailure("applied_conclusion_descriptor_unsealed")
	}
	ctx := descriptor.context
	if !ctx.sealed {
		return AppliedConclusionProof{}, operation.NewMechanicalFailure("applied_conclusion_context_unsealed")
	}
	mats := append([]DomainDecisionMaterial(nil), materials...)
	for _, m := range mats {
		if !m.sealed || !m.context.Equal(ctx) {
			return AppliedConclusionProof{}, operation.NewMechanicalFailure("applied_conclusion_material_invalid")
		}
	}
	return AppliedConclusionProof{
		sealed:     true,
		link:       link,
		descriptor: descriptor,
		materials:  mats,
		context:    ctx,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether p was package-constructed.
func (p AppliedConclusionProof) Sealed() bool { return p.sealed }

// Descriptor returns the sealed domain conclusion descriptor.
func (p AppliedConclusionProof) Descriptor() DomainConclusionDescriptor { return p.descriptor }

// Materials returns a copy of sealed domain decision materials.
func (p AppliedConclusionProof) Materials() []DomainDecisionMaterial {
	return append([]DomainDecisionMaterial(nil), p.materials...)
}

// PublicationContext returns the sealed publication context.
func (p AppliedConclusionProof) PublicationContext() PublicationContext { return p.context }

// AuthorizedMutationPlan is a sealed plan for authorization-bearing ResourceMutation.
type AuthorizedMutationPlan struct {
	sealed     bool
	context    PublicationContext
	evaluation CompletedAuthorizationEvaluation
	proofs     []AppliedMutationProof
}

// NewAuthorizedMutationPlan seals an authorized mutation plan.
func NewAuthorizedMutationPlan(
	ctx PublicationContext,
	evaluation CompletedAuthorizationEvaluation,
	proofs []AppliedMutationProof,
) (AuthorizedMutationPlan, operation.MechanicalFailure) {
	if !ctx.sealed {
		return AuthorizedMutationPlan{}, operation.NewMechanicalFailure("authorized_plan_context_unsealed")
	}
	if !evaluation.sealed {
		return AuthorizedMutationPlan{}, operation.NewMechanicalFailure("authorized_plan_evaluation_unsealed")
	}
	if !evaluation.context.Equal(ctx) {
		return AuthorizedMutationPlan{}, operation.NewMechanicalFailure("authorized_plan_context_mismatch")
	}
	if len(proofs) == 0 {
		return AuthorizedMutationPlan{}, operation.NewMechanicalFailure("authorized_plan_proofs_empty")
	}
	cp := append([]AppliedMutationProof(nil), proofs...)
	for _, p := range cp {
		if !p.sealed || !p.context.Equal(ctx) {
			return AuthorizedMutationPlan{}, operation.NewMechanicalFailure("authorized_plan_proof_invalid")
		}
	}
	return AuthorizedMutationPlan{
		sealed:     true,
		context:    ctx,
		evaluation: evaluation,
		proofs:     cp,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether p was package-constructed.
func (p AuthorizedMutationPlan) Sealed() bool { return p.sealed }

// AutomaticMutationPlan is a sealed plan for AutomaticMutationOnly ResourceMutation.
type AutomaticMutationPlan struct {
	sealed  bool
	context PublicationContext
	proofs  []AppliedMutationProof
}

// NewAutomaticMutationPlan seals an automatic mutation plan.
func NewAutomaticMutationPlan(
	ctx PublicationContext,
	proofs []AppliedMutationProof,
) (AutomaticMutationPlan, operation.MechanicalFailure) {
	if !ctx.sealed {
		return AutomaticMutationPlan{}, operation.NewMechanicalFailure("automatic_plan_context_unsealed")
	}
	if len(proofs) == 0 {
		return AutomaticMutationPlan{}, operation.NewMechanicalFailure("automatic_plan_proofs_empty")
	}
	cp := append([]AppliedMutationProof(nil), proofs...)
	for _, p := range cp {
		if !p.sealed || !p.context.Equal(ctx) {
			return AutomaticMutationPlan{}, operation.NewMechanicalFailure("automatic_plan_proof_invalid")
		}
	}
	return AutomaticMutationPlan{sealed: true, context: ctx, proofs: cp}, operation.MechanicalFailure{}
}

// Sealed reports whether p was package-constructed.
func (p AutomaticMutationPlan) Sealed() bool { return p.sealed }

// DomainConclusionPlan is a sealed plan for publishable DomainConclusion.
type DomainConclusionPlan struct {
	sealed     bool
	context    PublicationContext
	evaluation CompletedAuthorizationEvaluation
	proofs     []AppliedConclusionProof
}

// NewDomainConclusionPlan seals a domain-conclusion plan.
func NewDomainConclusionPlan(
	ctx PublicationContext,
	evaluation CompletedAuthorizationEvaluation,
	proofs []AppliedConclusionProof,
) (DomainConclusionPlan, operation.MechanicalFailure) {
	if !ctx.sealed {
		return DomainConclusionPlan{}, operation.NewMechanicalFailure("conclusion_plan_context_unsealed")
	}
	if !evaluation.sealed {
		return DomainConclusionPlan{}, operation.NewMechanicalFailure("conclusion_plan_evaluation_unsealed")
	}
	if !evaluation.context.Equal(ctx) {
		return DomainConclusionPlan{}, operation.NewMechanicalFailure("conclusion_plan_context_mismatch")
	}
	if len(proofs) == 0 {
		return DomainConclusionPlan{}, operation.NewMechanicalFailure("conclusion_plan_proofs_empty")
	}
	cp := append([]AppliedConclusionProof(nil), proofs...)
	for _, p := range cp {
		if !p.sealed || !p.context.Equal(ctx) {
			return DomainConclusionPlan{}, operation.NewMechanicalFailure("conclusion_plan_proof_invalid")
		}
	}
	return DomainConclusionPlan{
		sealed:     true,
		context:    ctx,
		evaluation: evaluation,
		proofs:     cp,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether p was package-constructed.
func (p DomainConclusionPlan) Sealed() bool { return p.sealed }
