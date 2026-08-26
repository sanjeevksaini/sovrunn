package policyeval

import (
	"context"
	"errors"
)

// Fixed construction diagnostic (D-08 / design §5.4). Never interpolates
// digest, reason code, or fixture values.
var errFakeConfigurationInvalid = errors.New("policy evaluation fake configuration invalid")

// Private adapter-failure signal returned by Evaluate for missing or
// configured-failure digests. The evaluation boundary discards the raw error
// text and normalizes to NonResultAdapterFailure (D-12 / §4.7).
var errFakeAdapterFailure = errors.New("policy evaluation adapter failed")

// inputDigestLen is the approved v1 lowercase SHA-256 hex length.
const inputDigestLen = 64

// FakeFixture is one digest-keyed fake configuration entry (D-12 / §4.7).
// Exactly one of Conclusion (non-nil) or AdapterFailure (true) must be set.
//
// YAML tags support the test-only fixture decoder (D-13); they are not a
// production input surface.
type FakeFixture struct {
	InputDigest    string      `yaml:"inputDigest"`
	Conclusion     *Conclusion `yaml:"conclusion,omitempty"`
	AdapterFailure bool        `yaml:"adapterFailure,omitempty"`
}

// fakeBehavior is the immutable per-digest lookup entry. fail==true means
// AdapterFailure; otherwise conclusion is non-nil and owned by the fake.
type fakeBehavior struct {
	conclusion *Conclusion
	fail       bool
}

// DeterministicFake is the Phase 2R digest-keyed conformance adapter
// (CDG-F17-04 / D-12). Configuration is immutable after construction. Lookup
// performs no domain-policy interpretation and holds no time source.
type DeterministicFake struct {
	byDigest map[string]fakeBehavior
}

// NewDeterministicFake deep-copies fixtures into an immutable digest-keyed
// fake. The input slice is intentionally duplicate-preserving so repeated
// digests survive long enough to be rejected. An empty slice is valid
// (every digest is unconfigured). Construction fails with a fixed generic
// diagnostic and no partial fake on duplicate/malformed digest, invalid
// conclusion/failure one-of, or malformed configured conclusion.
func NewDeterministicFake(fixtures []FakeFixture) (*DeterministicFake, error) {
	byDigest := make(map[string]fakeBehavior, len(fixtures))

	for i := range fixtures {
		fx := fixtures[i]

		if !validInputDigest(fx.InputDigest) {
			return nil, errFakeConfigurationInvalid
		}
		if _, dup := byDigest[fx.InputDigest]; dup {
			return nil, errFakeConfigurationInvalid
		}

		hasConclusion := fx.Conclusion != nil
		hasFailure := fx.AdapterFailure
		if hasConclusion == hasFailure {
			// Neither or both selected is invalid.
			return nil, errFakeConfigurationInvalid
		}

		if hasFailure {
			byDigest[fx.InputDigest] = fakeBehavior{fail: true}
			continue
		}

		owned := copyConclusion(*fx.Conclusion)
		if err := validateConclusion(owned); err != nil {
			return nil, errFakeConfigurationInvalid
		}
		byDigest[fx.InputDigest] = fakeBehavior{conclusion: &owned}
	}

	return &DeterministicFake{byDigest: byDigest}, nil
}

// Evaluate looks up inputDigest in the immutable fixture map.
// Configured conclusions are returned as a fresh deep copy. Missing and
// configured-failure digests return a zero Conclusion and a non-nil
// adapter-failure signal. NormalizedInput is ignored: the fake never
// interprets domain meaning (CDG-F17-04).
func (f *DeterministicFake) Evaluate(
	_ context.Context,
	_ NormalizedInput,
	inputDigest string,
) (Conclusion, error) {
	var zero Conclusion
	if f == nil {
		return zero, errFakeAdapterFailure
	}
	beh, ok := f.byDigest[inputDigest]
	if !ok || beh.fail {
		return zero, errFakeAdapterFailure
	}
	return copyConclusion(*beh.conclusion), nil
}

// validInputDigest reports whether d is a lowercase 64-character hex string
// matching the approved v1 input-digest grammar.
func validInputDigest(d string) bool {
	if len(d) != inputDigestLen {
		return false
	}
	for i := 0; i < len(d); i++ {
		c := d[i]
		switch {
		case c >= '0' && c <= '9':
		case c >= 'a' && c <= 'f':
		default:
			return false
		}
	}
	return true
}
