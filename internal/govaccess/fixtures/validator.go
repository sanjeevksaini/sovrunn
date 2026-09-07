package fixtures

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidSet is returned when a fixture set violates graph constraints.
	ErrInvalidSet = errors.New("fixtures: invalid fixture set")
	// ErrProhibitedConcept is returned when prohibited concepts are present.
	ErrProhibitedConcept = errors.New("fixtures: prohibited concept")
)

var allowedKinds = map[ResourceKind]bool{
	KindAccessGroup:             true,
	KindMembership:              true,
	KindRoleDefinition:          true,
	KindRoleAssignment:          true,
	KindPrivilegedAccessRequest: true,
	KindAccessReview:            true,
	KindApprovalPolicy:          true,
	KindApprovalRequest:         true,
	KindExceptionGrant:          true,
	KindGovernanceProfileV1:     true,
}

// Validate enforces referential closure, owned-kind boundaries, and prohibited
// concept gates before any evaluation/publication.
func Validate(set FixtureSet) error {
	if len(set.ProhibitedConcepts) > 0 {
		first := set.ProhibitedConcepts[0]
		return fmt.Errorf("%w: %s/%s", ErrProhibitedConcept, first.Category, first.Identifier)
	}

	seen := make(map[string]ResourceKind, len(set.Resources))
	for _, r := range set.Resources {
		if !allowedKinds[r.Kind] {
			return fmt.Errorf("%w: unknown or unowned kind %q", ErrInvalidSet, r.Kind)
		}
		if r.UID == "" {
			return fmt.Errorf("%w: resource uid is required", ErrInvalidSet)
		}
		if prior, ok := seen[r.UID]; ok {
			return fmt.Errorf("%w: duplicate uid %q in %q and %q", ErrInvalidSet, r.UID, prior, r.Kind)
		}
		seen[r.UID] = r.Kind
	}
	for _, r := range set.Resources {
		for _, ref := range r.RefUIDs {
			if ref == "" {
				return fmt.Errorf("%w: empty reference in %q", ErrInvalidSet, r.UID)
			}
			if _, ok := seen[ref]; !ok {
				return fmt.Errorf("%w: missing reference %q from %q", ErrInvalidSet, ref, r.UID)
			}
		}
	}
	return nil
}
