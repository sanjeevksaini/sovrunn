// Package operation holds FEATURE-0018 neutral leaf protocol values (DD-01).
//
// Types carry no authorization, approval, exception, or lifecycle meaning.
// Construction is sealed; mutable inputs are deep-copied; nil/empty inputs
// fail closed as MechanicalFailure.
package operation

import "bytes"

// MechanicalFailure identifies a local invariant, evidence-construction, or
// publication-mechanism failure. It is never a domain denial outcome.
type MechanicalFailure struct {
	reason string
}

// NewMechanicalFailure returns a MechanicalFailure with a non-empty reason.
// An empty reason yields a closed default reason so callers never observe an
// empty failure value.
func NewMechanicalFailure(reason string) MechanicalFailure {
	if reason == "" {
		reason = "mechanical_failure"
	}
	return MechanicalFailure{reason: reason}
}

// Reason returns the stable mechanical failure reason.
func (f MechanicalFailure) Reason() string {
	return f.reason
}

// ResultMaterial is an immutable deep-copied canonical successful result payload.
type ResultMaterial struct {
	sealed bool
	bytes  []byte
}

// NewResultMaterial seals a deep copy of canonicalSuccessfulResultBytes.
// Nil or empty input fails closed.
func NewResultMaterial(canonicalSuccessfulResultBytes []byte) (ResultMaterial, MechanicalFailure) {
	if len(canonicalSuccessfulResultBytes) == 0 {
		return ResultMaterial{}, NewMechanicalFailure("result_material_empty")
	}
	cp := make([]byte, len(canonicalSuccessfulResultBytes))
	copy(cp, canonicalSuccessfulResultBytes)
	return ResultMaterial{sealed: true, bytes: cp}, MechanicalFailure{}
}

// Bytes returns a defensive copy of the sealed canonical result bytes.
func (m ResultMaterial) Bytes() []byte {
	if !m.sealed || len(m.bytes) == 0 {
		return nil
	}
	cp := make([]byte, len(m.bytes))
	copy(cp, m.bytes)
	return cp
}

// Sealed reports whether m was constructed through NewResultMaterial.
func (m ResultMaterial) Sealed() bool {
	return m.sealed
}

// Equal reports whether two sealed ResultMaterial values carry identical bytes.
func (m ResultMaterial) Equal(other ResultMaterial) bool {
	if !m.sealed || !other.sealed {
		return false
	}
	return bytes.Equal(m.bytes, other.bytes)
}

// ParticipantBinding is an opaque participant identity binding for applied links.
type ParticipantBinding struct {
	sealed bool
	value  []byte
}

// NewParticipantBinding seals a non-empty binding value.
func NewParticipantBinding(value []byte) (ParticipantBinding, MechanicalFailure) {
	if len(value) == 0 {
		return ParticipantBinding{}, NewMechanicalFailure("participant_binding_empty")
	}
	cp := make([]byte, len(value))
	copy(cp, value)
	return ParticipantBinding{sealed: true, value: cp}, MechanicalFailure{}
}

// PublicationBinding is an opaque publication-context binding.
type PublicationBinding struct {
	sealed bool
	value  []byte
}

// NewPublicationBinding seals a non-empty binding value.
func NewPublicationBinding(value []byte) (PublicationBinding, MechanicalFailure) {
	if len(value) == 0 {
		return PublicationBinding{}, NewMechanicalFailure("publication_binding_empty")
	}
	cp := make([]byte, len(value))
	copy(cp, value)
	return PublicationBinding{sealed: true, value: cp}, MechanicalFailure{}
}

// ShadowDeltaDigest is an opaque shadow-delta digest.
type ShadowDeltaDigest struct {
	sealed bool
	value  []byte
}

// NewShadowDeltaDigest seals a non-empty digest value.
func NewShadowDeltaDigest(value []byte) (ShadowDeltaDigest, MechanicalFailure) {
	if len(value) == 0 {
		return ShadowDeltaDigest{}, NewMechanicalFailure("shadow_delta_digest_empty")
	}
	cp := make([]byte, len(value))
	copy(cp, value)
	return ShadowDeltaDigest{sealed: true, value: cp}, MechanicalFailure{}
}

// IndexDeltaDigest is an opaque index-delta digest.
type IndexDeltaDigest struct {
	sealed bool
	value  []byte
}

// NewIndexDeltaDigest seals a non-empty digest value.
func NewIndexDeltaDigest(value []byte) (IndexDeltaDigest, MechanicalFailure) {
	if len(value) == 0 {
		return IndexDeltaDigest{}, NewMechanicalFailure("index_delta_digest_empty")
	}
	cp := make([]byte, len(value))
	copy(cp, value)
	return IndexDeltaDigest{sealed: true, value: cp}, MechanicalFailure{}
}

// MutationDigest is an opaque mutation digest used by completion bindings.
type MutationDigest struct {
	sealed bool
	value  []byte
}

// NewMutationDigest seals a non-empty digest value.
func NewMutationDigest(value []byte) (MutationDigest, MechanicalFailure) {
	if len(value) == 0 {
		return MutationDigest{}, NewMechanicalFailure("mutation_digest_empty")
	}
	cp := make([]byte, len(value))
	copy(cp, value)
	return MutationDigest{sealed: true, value: cp}, MechanicalFailure{}
}

// AppliedChangeLink is a sealed mechanical link for a ResourceMutation application.
type AppliedChangeLink struct {
	sealed      bool
	participant ParticipantBinding
	publication PublicationBinding
	shadow      ShadowDeltaDigest
	index       IndexDeltaDigest
	hasResult   bool
	result      ResultMaterial
}

// NewAppliedChangeLink seals a mechanical applied-change link. Optional result
// may be omitted (zero ResultMaterial with hasResult=false).
func NewAppliedChangeLink(
	participant ParticipantBinding,
	publication PublicationBinding,
	shadow ShadowDeltaDigest,
	index IndexDeltaDigest,
	result ResultMaterial,
	hasResult bool,
) (AppliedChangeLink, MechanicalFailure) {
	if !participant.sealed {
		return AppliedChangeLink{}, NewMechanicalFailure("participant_binding_unsealed")
	}
	if !publication.sealed {
		return AppliedChangeLink{}, NewMechanicalFailure("publication_binding_unsealed")
	}
	if !shadow.sealed {
		return AppliedChangeLink{}, NewMechanicalFailure("shadow_delta_digest_unsealed")
	}
	if !index.sealed {
		return AppliedChangeLink{}, NewMechanicalFailure("index_delta_digest_unsealed")
	}
	if hasResult && !result.sealed {
		return AppliedChangeLink{}, NewMechanicalFailure("result_material_unsealed")
	}
	link := AppliedChangeLink{
		sealed:      true,
		participant: participant,
		publication: publication,
		shadow:      shadow,
		index:       index,
		hasResult:   hasResult,
	}
	if hasResult {
		link.result = result
	}
	return link, MechanicalFailure{}
}

// Sealed reports whether the link was constructed through NewAppliedChangeLink.
func (l AppliedChangeLink) Sealed() bool { return l.sealed }

// HasResult reports whether a ResultMaterial was bound.
func (l AppliedChangeLink) HasResult() bool { return l.hasResult }

// Result returns the optional sealed ResultMaterial.
func (l AppliedChangeLink) Result() (ResultMaterial, bool) {
	if !l.sealed || !l.hasResult {
		return ResultMaterial{}, false
	}
	return l.result, true
}

// AppliedConclusionLink is a sealed mechanical link for a DomainConclusion
// publication (no shadow/index delta).
type AppliedConclusionLink struct {
	sealed      bool
	participant ParticipantBinding
	publication PublicationBinding
}

// NewAppliedConclusionLink seals a mechanical applied-conclusion link.
func NewAppliedConclusionLink(
	participant ParticipantBinding,
	publication PublicationBinding,
) (AppliedConclusionLink, MechanicalFailure) {
	if !participant.sealed {
		return AppliedConclusionLink{}, NewMechanicalFailure("participant_binding_unsealed")
	}
	if !publication.sealed {
		return AppliedConclusionLink{}, NewMechanicalFailure("publication_binding_unsealed")
	}
	return AppliedConclusionLink{
		sealed:      true,
		participant: participant,
		publication: publication,
	}, MechanicalFailure{}
}

// Sealed reports whether the link was constructed through NewAppliedConclusionLink.
func (l AppliedConclusionLink) Sealed() bool { return l.sealed }

// CompletionBinding is a sealed mechanical completion binding for idempotency.
type CompletionBinding struct {
	sealed      bool
	publication PublicationBinding
	mutation    MutationDigest
}

// NewCompletionBinding seals a completion binding from publication and mutation digests.
func NewCompletionBinding(
	publication PublicationBinding,
	mutation MutationDigest,
) (CompletionBinding, MechanicalFailure) {
	if !publication.sealed {
		return CompletionBinding{}, NewMechanicalFailure("publication_binding_unsealed")
	}
	if !mutation.sealed {
		return CompletionBinding{}, NewMechanicalFailure("mutation_digest_unsealed")
	}
	return CompletionBinding{
		sealed:      true,
		publication: publication,
		mutation:    mutation,
	}, MechanicalFailure{}
}

// Sealed reports whether the binding was constructed through NewCompletionBinding.
func (b CompletionBinding) Sealed() bool { return b.sealed }

type outcomeKind uint8

const (
	outcomeNone outcomeKind = iota
	outcomeFinalized
	outcomeDomain
	outcomeMechanical
)

// FinalizationOutcome is a closed discriminated union of Finalized(T), Domain(D),
// or MechanicalFailure. Exactly one variant is active.
type FinalizationOutcome[T, D any] struct {
	kind      outcomeKind
	finalized T
	domain    D
	failure   MechanicalFailure
}

// Finalized constructs a Finalized(T) outcome.
func Finalized[T, D any](value T) FinalizationOutcome[T, D] {
	return FinalizationOutcome[T, D]{kind: outcomeFinalized, finalized: value}
}

// Domain constructs a Domain(D) non-publication outcome.
func Domain[T, D any](value D) FinalizationOutcome[T, D] {
	return FinalizationOutcome[T, D]{kind: outcomeDomain, domain: value}
}

// Failure constructs a MechanicalFailure outcome.
func Failure[T, D any](failure MechanicalFailure) FinalizationOutcome[T, D] {
	if failure.reason == "" {
		failure = NewMechanicalFailure("mechanical_failure")
	}
	return FinalizationOutcome[T, D]{kind: outcomeMechanical, failure: failure}
}

// AsFinalized returns the Finalized variant when active.
func (o FinalizationOutcome[T, D]) AsFinalized() (T, bool) {
	var zero T
	if o.kind != outcomeFinalized {
		return zero, false
	}
	return o.finalized, true
}

// AsDomain returns the Domain variant when active.
func (o FinalizationOutcome[T, D]) AsDomain() (D, bool) {
	var zero D
	if o.kind != outcomeDomain {
		return zero, false
	}
	return o.domain, true
}

// AsMechanicalFailure returns the MechanicalFailure variant when active.
func (o FinalizationOutcome[T, D]) AsMechanicalFailure() (MechanicalFailure, bool) {
	if o.kind != outcomeMechanical {
		return MechanicalFailure{}, false
	}
	return o.failure, true
}

// KindName returns a stable discriminant label for tests and diagnostics.
func (o FinalizationOutcome[T, D]) KindName() string {
	switch o.kind {
	case outcomeFinalized:
		return "Finalized"
	case outcomeDomain:
		return "Domain"
	case outcomeMechanical:
		return "MechanicalFailure"
	default:
		return "Invalid"
	}
}
