package state

import (
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
)

// AppliedChangeReceipt proves an applied state change for one participant
// (DD-10; design §5.4). Constructible only by a state transaction.
type AppliedChangeReceipt struct {
	sealed           bool
	claim            ParticipantClaim
	context          evidence.PublicationContext
	binding          PermitBinding
	kind             PublicationKind
	mutationProof    evidence.AppliedMutationProof
	hasMutationProof bool
	conclusionProof  evidence.AppliedConclusionProof
	hasConclusion    bool
}

// ParticipantClaim returns the sealed participant claim.
func (r AppliedChangeReceipt) ParticipantClaim() ParticipantClaim { return r.claim }

// PublicationContext returns the sealed publication context.
func (r AppliedChangeReceipt) PublicationContext() evidence.PublicationContext { return r.context }

// PermitBinding returns the sealed permit binding.
func (r AppliedChangeReceipt) PermitBinding() PermitBinding { return r.binding }

// PublicationKind returns the sealed publication kind.
func (r AppliedChangeReceipt) PublicationKind() PublicationKind { return r.kind }

// EvidenceProof returns the mutation proof when PublicationKind is ResourceMutation.
func (r AppliedChangeReceipt) EvidenceProof() (evidence.AppliedMutationProof, bool) {
	if !r.sealed || !r.hasMutationProof {
		return evidence.AppliedMutationProof{}, false
	}
	return r.mutationProof, true
}

// ConclusionProof returns the conclusion proof when PublicationKind is DomainConclusion.
func (r AppliedChangeReceipt) ConclusionProof() (evidence.AppliedConclusionProof, bool) {
	if !r.sealed || !r.hasConclusion {
		return evidence.AppliedConclusionProof{}, false
	}
	return r.conclusionProof, true
}

// CommitReceipt proves post-install acceptance of prepared evidence.
type CommitReceipt struct {
	sealed   bool
	context  evidence.PublicationContext
	evidence evidence.PreparedEvidenceChange
	pins     []decision.ProfileRef
}

// PublicationContext returns the sealed publication context.
func (r CommitReceipt) PublicationContext() evidence.PublicationContext { return r.context }

// AcceptedEvidence returns the sealed accepted PreparedEvidenceChange.
func (r CommitReceipt) AcceptedEvidence() evidence.PreparedEvidenceChange { return r.evidence }

// ProfilePins returns the FEATURE-0013 DecisionProfile pins accepted with the evidence.
func (r CommitReceipt) ProfilePins() []decision.ProfileRef {
	return append([]decision.ProfileRef(nil), r.pins...)
}

// Sealed reports whether r was constructed by a committing transaction.
func (r CommitReceipt) Sealed() bool { return r.sealed }
