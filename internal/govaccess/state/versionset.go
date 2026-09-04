package state

import (
	"bytes"
	"crypto/sha256"
	"sort"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// ExpectedVersionPredicate is one resource/version CAS predicate.
type ExpectedVersionPredicate struct {
	ResourceUID     string
	ResourceVersion string
}

// ExpectedVersionSet is a sealed set of version predicates.
type ExpectedVersionSet struct {
	sealed     bool
	predicates []ExpectedVersionPredicate
	digest     []byte
}

// NewExpectedVersionSet seals a non-empty predicate set.
func NewExpectedVersionSet(predicates []ExpectedVersionPredicate) (ExpectedVersionSet, error) {
	if len(predicates) == 0 {
		return ExpectedVersionSet{}, errInvalid("expected_version_set_empty")
	}
	cp := append([]ExpectedVersionPredicate(nil), predicates...)
	sort.SliceStable(cp, func(i, j int) bool {
		if cp[i].ResourceUID != cp[j].ResourceUID {
			return cp[i].ResourceUID < cp[j].ResourceUID
		}
		return cp[i].ResourceVersion < cp[j].ResourceVersion
	})
	seen := map[string]bool{}
	for _, p := range cp {
		if p.ResourceUID == "" || p.ResourceVersion == "" {
			return ExpectedVersionSet{}, errInvalid("expected_version_predicate_invalid")
		}
		if seen[p.ResourceUID] {
			return ExpectedVersionSet{}, errInvalid("expected_version_predicate_duplicate")
		}
		seen[p.ResourceUID] = true
	}
	return ExpectedVersionSet{
		sealed:     true,
		predicates: cp,
		digest:     hashPredicates(cp),
	}, nil
}

// Sealed reports whether s was constructed through NewExpectedVersionSet.
func (s ExpectedVersionSet) Sealed() bool { return s.sealed }

// Predicates returns a defensive copy.
func (s ExpectedVersionSet) Predicates() []ExpectedVersionPredicate {
	return append([]ExpectedVersionPredicate(nil), s.predicates...)
}

// Digest returns a defensive copy of the set digest.
func (s ExpectedVersionSet) Digest() []byte { return copyBytes(s.digest) }

// AuthorizationDependencyVersionSet is the complete candidate dependency set.
type AuthorizationDependencyVersionSet struct {
	sealed                  bool
	expectedSnapshotVersion string
	predicates              []ExpectedVersionPredicate
	digest                  []byte
}

// NewAuthorizationDependencyVersionSet seals a dependency-version set.
func NewAuthorizationDependencyVersionSet(
	expectedSnapshotVersion string,
	predicates []ExpectedVersionPredicate,
) (AuthorizationDependencyVersionSet, error) {
	if expectedSnapshotVersion == "" {
		return AuthorizationDependencyVersionSet{}, errInvalid("authz_dependency_snapshot_empty")
	}
	base, err := NewExpectedVersionSet(predicates)
	if err != nil {
		return AuthorizationDependencyVersionSet{}, err
	}
	d := sha256.Sum256(append([]byte(expectedSnapshotVersion+"\x00"), base.digest...))
	return AuthorizationDependencyVersionSet{
		sealed:                  true,
		expectedSnapshotVersion: expectedSnapshotVersion,
		predicates:              base.predicates,
		digest:                  d[:],
	}, nil
}

// Sealed reports whether s was package-constructed.
func (s AuthorizationDependencyVersionSet) Sealed() bool { return s.sealed }

// ExpectedSnapshotVersion returns the sealed snapshot version.
func (s AuthorizationDependencyVersionSet) ExpectedSnapshotVersion() string {
	return s.expectedSnapshotVersion
}

// Predicates returns a defensive copy.
func (s AuthorizationDependencyVersionSet) Predicates() []ExpectedVersionPredicate {
	return append([]ExpectedVersionPredicate(nil), s.predicates...)
}

// Digest returns a defensive copy.
func (s AuthorizationDependencyVersionSet) Digest() []byte { return copyBytes(s.digest) }

// AsDependencyVersionDigest converts the set digest into an evidence digest.
func (s AuthorizationDependencyVersionSet) AsDependencyVersionDigest() (evidence.DependencyVersionDigest, operation.MechanicalFailure) {
	if !s.sealed {
		return evidence.DependencyVersionDigest{}, operation.NewMechanicalFailure("dependency_set_unsealed")
	}
	return evidence.NewDependencyVersionDigest(s.digest)
}

func hashPredicates(preds []ExpectedVersionPredicate) []byte {
	h := sha256.New()
	for _, p := range preds {
		_, _ = h.Write([]byte(p.ResourceUID))
		_, _ = h.Write([]byte{0})
		_, _ = h.Write([]byte(p.ResourceVersion))
		_, _ = h.Write([]byte{0})
	}
	return h.Sum(nil)
}

func copyBytes(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}
	out := make([]byte, len(in))
	copy(out, in)
	return out
}

func predicatesMatch(current map[string]string, preds []ExpectedVersionPredicate) bool {
	for _, p := range preds {
		v, ok := current[p.ResourceUID]
		if !ok || v != p.ResourceVersion {
			return false
		}
	}
	return true
}

func digestEqual(a, b []byte) bool { return bytes.Equal(a, b) }
