package evidence

import (
	"sort"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// MutationDescriptor is a sealed mutation AuditEvent descriptor derived from
// model.MutationEventFacts (design §3.2 evidence-carrier input).
type MutationDescriptor struct {
	sealed  bool
	facts   model.MutationEventFacts
	context PublicationContext
	ordinal int
}

// NewMutationDescriptor seals a mutation descriptor from MutationEventFacts.
func NewMutationDescriptor(
	facts model.MutationEventFacts,
	ctx PublicationContext,
) (MutationDescriptor, operation.MechanicalFailure) {
	if !facts.Sealed() {
		return MutationDescriptor{}, operation.NewMechanicalFailure("mutation_descriptor_facts_unsealed")
	}
	if !ctx.sealed {
		return MutationDescriptor{}, operation.NewMechanicalFailure("mutation_descriptor_context_unsealed")
	}
	if facts.EventType() == "" {
		return MutationDescriptor{}, operation.NewMechanicalFailure("mutation_descriptor_event_empty")
	}
	return MutationDescriptor{
		sealed:  true,
		facts:   facts,
		context: ctx,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether d was package-constructed.
func (d MutationDescriptor) Sealed() bool { return d.sealed }

// EventType returns the sealed taxonomy event type.
func (d MutationDescriptor) EventType() string { return d.facts.EventType() }

// Facts returns the sealed MutationEventFacts carrier input.
func (d MutationDescriptor) Facts() model.MutationEventFacts { return d.facts }

// PublicationContext returns the sealed publication context.
func (d MutationDescriptor) PublicationContext() PublicationContext { return d.context }

// DomainConclusionDescriptor is the sealed exceptionproposal.decided descriptor
// for a publishable domain conclusion (ADH-2026-073).
type DomainConclusionDescriptor struct {
	sealed    bool
	facts     model.ExceptionDecisionFacts
	context   PublicationContext
	eventType string
}

// EventTypeExceptionProposalDecided is the ADH-2026-073 terminal taxonomy event.
const EventTypeExceptionProposalDecided = "exceptionproposal.decided"

// EventTypeExceptionGrantIssued is the Grant-path issuance taxonomy event.
const EventTypeExceptionGrantIssued = "exceptiongrant.issued"

// NewDomainConclusionDescriptor seals an exceptionproposal.decided descriptor.
func NewDomainConclusionDescriptor(
	facts model.ExceptionDecisionFacts,
	ctx PublicationContext,
) (DomainConclusionDescriptor, operation.MechanicalFailure) {
	if !facts.Sealed() {
		return DomainConclusionDescriptor{}, operation.NewMechanicalFailure("conclusion_descriptor_facts_unsealed")
	}
	if !ctx.sealed {
		return DomainConclusionDescriptor{}, operation.NewMechanicalFailure("conclusion_descriptor_context_unsealed")
	}
	if facts.Result() != model.ExceptionTerminalDeny {
		return DomainConclusionDescriptor{}, operation.NewMechanicalFailure("conclusion_descriptor_requires_deny")
	}
	return DomainConclusionDescriptor{
		sealed:    true,
		facts:     facts,
		context:   ctx,
		eventType: EventTypeExceptionProposalDecided,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether d was package-constructed.
func (d DomainConclusionDescriptor) Sealed() bool { return d.sealed }

// EventType returns the sealed taxonomy event type.
func (d DomainConclusionDescriptor) EventType() string { return d.eventType }

// Facts returns the sealed ExceptionDecisionFacts carrier input.
func (d DomainConclusionDescriptor) Facts() model.ExceptionDecisionFacts { return d.facts }

// PublicationContext returns the sealed publication context.
func (d DomainConclusionDescriptor) PublicationContext() PublicationContext { return d.context }

// ValidateExceptionTerminalDescriptorSet validates the ADH-2026-073 Grant/Deny
// matrix and FEATURE-0013 linkage against ExceptionDecisionFacts.
func ValidateExceptionTerminalDescriptorSet(
	facts model.ExceptionDecisionFacts,
	mutationDescriptors []MutationDescriptor,
	conclusion *DomainConclusionDescriptor,
	materials []DomainDecisionMaterial,
	ctx PublicationContext,
) operation.MechanicalFailure {
	if !facts.Sealed() {
		return operation.NewMechanicalFailure("exception_terminal_facts_unsealed")
	}
	if !ctx.sealed {
		return operation.NewMechanicalFailure("exception_terminal_context_unsealed")
	}
	if len(materials) != 1 || !materials[0].sealed || !materials[0].hasException {
		return operation.NewMechanicalFailure("exception_terminal_material_required")
	}
	matFacts, ok := materials[0].ExceptionFacts()
	if !ok || !matFacts.Equal(facts) {
		return operation.NewMechanicalFailure("exception_terminal_material_mismatch")
	}
	if !materials[0].context.Equal(ctx) {
		return operation.NewMechanicalFailure("exception_terminal_material_context_mismatch")
	}

	switch facts.Result() {
	case model.ExceptionTerminalGrant:
		if conclusion != nil {
			return operation.NewMechanicalFailure("exception_terminal_grant_forbids_conclusion")
		}
		if len(mutationDescriptors) != 2 {
			return operation.NewMechanicalFailure("exception_terminal_grant_descriptor_count")
		}
		var decided, issued *MutationDescriptor
		for i := range mutationDescriptors {
			d := &mutationDescriptors[i]
			if !d.sealed || !d.context.Equal(ctx) {
				return operation.NewMechanicalFailure("exception_terminal_descriptor_context_mismatch")
			}
			switch d.EventType() {
			case EventTypeExceptionProposalDecided:
				if decided != nil {
					return operation.NewMechanicalFailure("exception_terminal_grant_descriptor_order")
				}
				decided = d
			case EventTypeExceptionGrantIssued:
				if issued != nil {
					return operation.NewMechanicalFailure("exception_terminal_grant_descriptor_order")
				}
				issued = d
			default:
				return operation.NewMechanicalFailure("exception_terminal_grant_descriptor_order")
			}
		}
		if decided == nil || issued == nil {
			return operation.NewMechanicalFailure("exception_terminal_grant_descriptor_order")
		}
		if facts.GrantUID() == "" {
			return operation.NewMechanicalFailure("exception_terminal_grant_uid_required")
		}
		return operation.MechanicalFailure{}
	case model.ExceptionTerminalDeny:
		if conclusion == nil || !conclusion.sealed {
			return operation.NewMechanicalFailure("exception_terminal_deny_conclusion_required")
		}
		if len(mutationDescriptors) != 0 {
			return operation.NewMechanicalFailure("exception_terminal_deny_forbids_mutation")
		}
		if conclusion.EventType() != EventTypeExceptionProposalDecided {
			return operation.NewMechanicalFailure("exception_terminal_deny_event")
		}
		if !conclusion.facts.Equal(facts) || !conclusion.context.Equal(ctx) {
			return operation.NewMechanicalFailure("exception_terminal_deny_linkage")
		}
		if facts.GrantUID() != "" {
			return operation.NewMechanicalFailure("exception_terminal_deny_forbids_grant")
		}
		return operation.MechanicalFailure{}
	default:
		return operation.NewMechanicalFailure("exception_terminal_result_invalid")
	}
}

func descriptorOrderKey(d MutationDescriptor) string {
	return d.EventType() + "|" + d.facts.ResourceRef().UID + "|" + d.facts.ResourceVersion()
}

// SortMutationDescriptors returns descriptors in canonical event/resource/version order.
func SortMutationDescriptors(in []MutationDescriptor) []MutationDescriptor {
	out := append([]MutationDescriptor(nil), in...)
	sort.SliceStable(out, func(i, j int) bool {
		return descriptorOrderKey(out[i]) < descriptorOrderKey(out[j])
	})
	return out
}
