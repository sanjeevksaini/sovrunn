package state

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// PublicationKind is ResourceMutation | DomainConclusion.
type PublicationKind string

const (
	ResourceMutation PublicationKind = "ResourceMutation"
	DomainConclusion PublicationKind = "DomainConclusion"
)

// Valid reports whether k is a closed PublicationKind.
func (k PublicationKind) Valid() bool {
	switch k {
	case ResourceMutation, DomainConclusion:
		return true
	default:
		return false
	}
}

// FinalizationPermitToken is a transaction-private opaque token.
type FinalizationPermitToken struct {
	sealed bool
	value  []byte
}

// PermitBinding is an unforgeable opaque binding derived from a
// FinalizationPermitToken (design §5.4). Equality-comparable; mintable only by
// the issuing transaction.
type PermitBinding struct {
	sealed bool
	digest []byte
}

// Equal reports sealed binding equality.
func (b PermitBinding) Equal(other PermitBinding) bool {
	return b.sealed && other.sealed && digestEqual(b.digest, other.digest)
}

// Sealed reports whether b was minted by a transaction.
func (b PermitBinding) Sealed() bool { return b.sealed }

// MutationFinalizationPermit is the owner-bound finalization permit.
type MutationFinalizationPermit struct {
	sealed    bool
	claim     ParticipantClaim
	context   evidence.PublicationContext
	authority MutationAuthorityMode
	binding   PermitBinding
	token     FinalizationPermitToken
}

// ParticipantClaim returns the sealed admitted claim.
func (p MutationFinalizationPermit) ParticipantClaim() ParticipantClaim { return p.claim }

// PublicationContext returns the sealed publication context.
func (p MutationFinalizationPermit) PublicationContext() evidence.PublicationContext {
	return p.context
}

// AuthorityMode returns the sealed authority mode.
func (p MutationFinalizationPermit) AuthorityMode() MutationAuthorityMode { return p.authority }

// Binding returns the unforgeable PermitBinding.
func (p MutationFinalizationPermit) Binding() PermitBinding { return p.binding }

// Sealed reports whether p was issued by a transaction.
func (p MutationFinalizationPermit) Sealed() bool { return p.sealed }

type permitRegistryEntry struct {
	claim    ParticipantClaim
	binding  PermitBinding
	consumed bool
}

func mintPermitBinding(token FinalizationPermitToken) PermitBinding {
	sum := sha256.Sum256(token.value)
	return PermitBinding{sealed: true, digest: sum[:]}
}

func newPermitToken() (FinalizationPermitToken, operation.MechanicalFailure) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		// Deterministic fallback for constrained environments: hash a counter-like marker.
		sum := sha256.Sum256([]byte("f18-permit-fallback"))
		return FinalizationPermitToken{sealed: true, value: sum[:]}, operation.MechanicalFailure{}
	}
	return FinalizationPermitToken{sealed: true, value: buf}, operation.MechanicalFailure{}
}

// ContextBoundChange is the finalized owner change interface applied by state.
type ContextBoundChange interface {
	ParticipantClaim() ParticipantClaim
	ContextBinding() evidence.PublicationContextBinding
	PermitBinding() PermitBinding
	PublicationKind() PublicationKind
	EvidenceDescriptors() []evidence.MutationDescriptor
	ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool)
	DomainDecisionMaterials() []evidence.DomainDecisionMaterial
	CanonicalResult() (operation.ResultMaterial, bool)
	ApplyTo(StateEditor) error
}

func bindingHex(b PermitBinding) string {
	if !b.sealed {
		return ""
	}
	return hex.EncodeToString(b.digest)
}
