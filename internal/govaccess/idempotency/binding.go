package idempotency

import (
	"bytes"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// RequestBinding is the immutable request binding compared after lookup
// (exactTarget, scopeRef, canonicalRequestDigest). It is a distinct value
// type from IdempotencyLookupKey (DD-11).
type RequestBinding struct {
	sealed                 bool
	exactTarget            apimeta.TypedRef
	scopeRef               apimeta.ScopeRef
	canonicalRequestDigest []byte
}

// NewRequestBinding seals a structurally valid request binding. Digests are
// deep-copied. Nil/empty digest fails closed.
func NewRequestBinding(
	exactTarget apimeta.TypedRef,
	scopeRef apimeta.ScopeRef,
	canonicalRequestDigest []byte,
) (RequestBinding, operation.MechanicalFailure) {
	if exactTarget.APIVersion == "" || exactTarget.Kind == "" || exactTarget.Name == "" {
		return RequestBinding{}, operation.NewMechanicalFailure("request_binding_exact_target_invalid")
	}
	if scopeRef.Kind == "" || !apimeta.ScopeKind(scopeRef.Kind).Valid() {
		return RequestBinding{}, operation.NewMechanicalFailure("request_binding_scope_kind_invalid")
	}
	if apimeta.ScopeKind(scopeRef.Kind) != apimeta.ScopePlatform {
		if scopeRef.APIVersion == "" || scopeRef.Name == "" {
			return RequestBinding{}, operation.NewMechanicalFailure("request_binding_scope_invalid")
		}
	}
	if len(canonicalRequestDigest) == 0 {
		return RequestBinding{}, operation.NewMechanicalFailure("request_binding_digest_empty")
	}
	maxBytes := apivalid.DefaultLimits().MaxObjectBytes
	if maxBytes > 0 && len(canonicalRequestDigest) > maxBytes {
		return RequestBinding{}, operation.NewMechanicalFailure("request_binding_digest_too_large")
	}
	cp := make([]byte, len(canonicalRequestDigest))
	copy(cp, canonicalRequestDigest)
	return RequestBinding{
		sealed:                 true,
		exactTarget:            exactTarget,
		scopeRef:               scopeRef,
		canonicalRequestDigest: cp,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether the binding was constructed through NewRequestBinding.
func (b RequestBinding) Sealed() bool { return b.sealed }

// ExactTarget returns the sealed exact target reference.
func (b RequestBinding) ExactTarget() apimeta.TypedRef { return b.exactTarget }

// ScopeRef returns the sealed scope reference.
func (b RequestBinding) ScopeRef() apimeta.ScopeRef { return b.scopeRef }

// CanonicalRequestDigest returns a defensive copy of the sealed digest.
func (b RequestBinding) CanonicalRequestDigest() []byte {
	if !b.sealed || len(b.canonicalRequestDigest) == 0 {
		return nil
	}
	cp := make([]byte, len(b.canonicalRequestDigest))
	copy(cp, b.canonicalRequestDigest)
	return cp
}

// Equal reports exact sealed field equality.
func (b RequestBinding) Equal(other RequestBinding) bool {
	if !b.sealed || !other.sealed {
		return false
	}
	return b.exactTarget == other.exactTarget &&
		b.scopeRef == other.scopeRef &&
		bytes.Equal(b.canonicalRequestDigest, other.canonicalRequestDigest)
}
