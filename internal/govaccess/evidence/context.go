package evidence

import (
	"bytes"
	"crypto/sha256"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// PublicationContext is the authoritative publication instant and sequence
// captured by a state transaction under stateMu (design §5.4).
type PublicationContext struct {
	sealed              bool
	publicationInstant  time.Time
	publicationSequence uint64
}

// NewPublicationContext seals a publication context. Instant must be UTC and non-zero.
func NewPublicationContext(instant time.Time, sequence uint64) (PublicationContext, operation.MechanicalFailure) {
	if instant.IsZero() {
		return PublicationContext{}, operation.NewMechanicalFailure("publication_context_instant_zero")
	}
	if sequence == 0 {
		return PublicationContext{}, operation.NewMechanicalFailure("publication_context_sequence_zero")
	}
	return PublicationContext{
		sealed:              true,
		publicationInstant:  instant.UTC(),
		publicationSequence: sequence,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether c was constructed through NewPublicationContext.
func (c PublicationContext) Sealed() bool { return c.sealed }

// PublicationInstant returns the sealed UTC publication instant.
func (c PublicationContext) PublicationInstant() time.Time { return c.publicationInstant }

// PublicationSequence returns the sealed publication sequence.
func (c PublicationContext) PublicationSequence() uint64 { return c.publicationSequence }

// Equal reports sealed equality.
func (c PublicationContext) Equal(other PublicationContext) bool {
	return c.sealed && other.sealed &&
		c.publicationInstant.Equal(other.publicationInstant) &&
		c.publicationSequence == other.publicationSequence
}

// PublicationContextBinding is an opaque binding derived from PublicationContext.
type PublicationContextBinding struct {
	sealed bool
	digest []byte
}

// BindPublicationContext derives a sealed binding digest from a context.
func BindPublicationContext(ctx PublicationContext) (PublicationContextBinding, operation.MechanicalFailure) {
	if !ctx.sealed {
		return PublicationContextBinding{}, operation.NewMechanicalFailure("publication_context_unsealed")
	}
	sum := sha256.Sum256([]byte(ctx.publicationInstant.UTC().Format(time.RFC3339Nano) + ":" + formatUint(ctx.publicationSequence)))
	return PublicationContextBinding{sealed: true, digest: copyBytes(sum[:])}, operation.MechanicalFailure{}
}

// Sealed reports whether b was constructed through BindPublicationContext.
func (b PublicationContextBinding) Sealed() bool { return b.sealed }

// Digest returns a defensive copy of the binding digest.
func (b PublicationContextBinding) Digest() []byte { return copyBytes(b.digest) }

// Equal reports sealed digest equality.
func (b PublicationContextBinding) Equal(other PublicationContextBinding) bool {
	return b.sealed && other.sealed && bytes.Equal(b.digest, other.digest)
}

// CandidateDigest is an opaque sealed candidate digest.
type CandidateDigest struct {
	sealed bool
	value  []byte
}

// NewCandidateDigest seals a non-empty digest.
func NewCandidateDigest(value []byte) (CandidateDigest, operation.MechanicalFailure) {
	if len(value) == 0 {
		return CandidateDigest{}, operation.NewMechanicalFailure("candidate_digest_empty")
	}
	return CandidateDigest{sealed: true, value: copyBytes(value)}, operation.MechanicalFailure{}
}

// Sealed reports whether d was constructed through NewCandidateDigest.
func (d CandidateDigest) Sealed() bool { return d.sealed }

// Bytes returns a defensive copy.
func (d CandidateDigest) Bytes() []byte { return copyBytes(d.value) }

// Equal reports sealed equality.
func (d CandidateDigest) Equal(other CandidateDigest) bool {
	return d.sealed && other.sealed && bytes.Equal(d.value, other.value)
}

// DependencyVersionDigest is an opaque sealed dependency-version digest.
type DependencyVersionDigest struct {
	sealed bool
	value  []byte
}

// NewDependencyVersionDigest seals a non-empty digest.
func NewDependencyVersionDigest(value []byte) (DependencyVersionDigest, operation.MechanicalFailure) {
	if len(value) == 0 {
		return DependencyVersionDigest{}, operation.NewMechanicalFailure("dependency_version_digest_empty")
	}
	return DependencyVersionDigest{sealed: true, value: copyBytes(value)}, operation.MechanicalFailure{}
}

// Sealed reports whether d was constructed through NewDependencyVersionDigest.
func (d DependencyVersionDigest) Sealed() bool { return d.sealed }

// Bytes returns a defensive copy.
func (d DependencyVersionDigest) Bytes() []byte { return copyBytes(d.value) }

// Equal reports sealed equality.
func (d DependencyVersionDigest) Equal(other DependencyVersionDigest) bool {
	return d.sealed && other.sealed && bytes.Equal(d.value, other.value)
}

// AdmissionDigest is an opaque sealed admission digest.
type AdmissionDigest struct {
	sealed bool
	value  []byte
}

// Bytes returns a defensive copy.
func (d AdmissionDigest) Bytes() []byte { return copyBytes(d.value) }

// Equal reports sealed equality.
func (d AdmissionDigest) Equal(other AdmissionDigest) bool {
	return d.sealed && other.sealed && bytes.Equal(d.value, other.value)
}

func copyBytes(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}
	out := make([]byte, len(in))
	copy(out, in)
	return out
}

func formatUint(v uint64) string {
	const digits = "0123456789"
	if v == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = digits[v%10]
		v /= 10
	}
	return string(buf[i:])
}

func sha256Concat(parts ...[]byte) []byte {
	h := sha256.New()
	for _, p := range parts {
		_, _ = h.Write(p)
		_, _ = h.Write([]byte{0})
	}
	return h.Sum(nil)
}
