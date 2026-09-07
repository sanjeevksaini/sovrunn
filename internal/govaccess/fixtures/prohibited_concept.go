package fixtures

import "fmt"

// ProhibitedCategory captures non-FEATURE-0018 concepts blocked by fixtures.
type ProhibitedCategory string

const (
	ProhibitedFutureResource    ProhibitedCategory = "future-resource"
	ProhibitedFutureField       ProhibitedCategory = "future-field"
	ProhibitedAdapter           ProhibitedCategory = "adapter"
	ProhibitedCredentialFlow    ProhibitedCategory = "credential-flow"
	ProhibitedProviderNativeIAM ProhibitedCategory = "provider-native-iam"
)

// ProhibitedConcept records a blocked concept for fixture gate assertions.
type ProhibitedConcept struct {
	Category   ProhibitedCategory
	Identifier string
}

// NewProhibitedConcept validates category + identifier for deterministic tests.
func NewProhibitedConcept(category ProhibitedCategory, identifier string) (ProhibitedConcept, error) {
	if identifier == "" {
		return ProhibitedConcept{}, fmt.Errorf("%w: identifier required", ErrProhibitedConcept)
	}
	switch category {
	case ProhibitedFutureResource, ProhibitedFutureField, ProhibitedAdapter, ProhibitedCredentialFlow, ProhibitedProviderNativeIAM:
		return ProhibitedConcept{Category: category, Identifier: identifier}, nil
	default:
		return ProhibitedConcept{}, fmt.Errorf("%w: unsupported category %q", ErrProhibitedConcept, category)
	}
}

// FixtureWithProhibitedFutureOwnedResource returns the VS0-CF-F18-50 negative fixture.
func FixtureWithProhibitedFutureOwnedResource() FixtureSet {
	concept, _ := NewProhibitedConcept(ProhibitedFutureResource, "CloudEnrollment")
	return NewBuilder().
		WithResource(Resource{Kind: KindAccessGroup, UID: "ag-1"}).
		WithProhibitedConcept(concept).
		Build()
}
