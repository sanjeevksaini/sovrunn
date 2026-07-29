package decision

// RetryIdempotencyKey is the caller-supplied retry idempotency identity
// (F13-SCALE-001; AD-013 carrier; AD-036). DecisionBody carries it as the
// retryKey string field.
//
// It is intentionally distinct from SemanticDecisionIdentity: equal retry
// keys do not imply equal semantic identity, and semantic-identity fields
// must never be used as a retry key.
type RetryIdempotencyKey string

// SemanticDecisionIdentity is the complete semantic identity descriptor
// (design §7.9; F13-SCALE-002; AD-036).
//
// CanonicalScope is derived from metadata.scopeRef (and profile scope rules)
// and is never an independent scope authority. It never overrides, replaces,
// or competes with metadata.scopeRef, which remains the sole scope authority.
// Digest is an opaque structural identity carrier; no algorithm is selected
// or executed.
type SemanticDecisionIdentity struct {
	Basis          string   `json:"basis"`            // profile-declared identity basis id; required non-empty
	InputRefs      []string `json:"inputRefs"`        // bounded ordered input identity refs; required non-empty
	CanonicalScope string   `json:"canonicalScope"`   // derived scope descriptor; required non-empty
	Digest         string   `json:"digest,omitempty"` // opaque structural carrier; no algorithm executed
}
