package apimeta

// Profile is the Matrix A resource-profile vocabulary (F12-PROFILE-001).
// Every externally exchanged object selects exactly one approved profile.
type Profile string

const (
	ProfileManagedResource          Profile = "ManagedResource"
	ProfileObservedExternalResource Profile = "ObservedExternalResource"
	ProfileVersionedDefinition      Profile = "VersionedDefinition"
	ProfileImmutableRecord          Profile = "ImmutableRecord"
	ProfileLongRunningOperation     Profile = "LongRunningOperation"
	ProfileTransientRequestResult   Profile = "TransientRequestResult"
	ProfileEmbeddedValue            Profile = "EmbeddedValue"
	ProfileListEnvelope             Profile = "ListEnvelope"
)

// AllProfiles returns the closed Matrix A profile set in stable order.
func AllProfiles() []Profile {
	return []Profile{
		ProfileManagedResource,
		ProfileObservedExternalResource,
		ProfileVersionedDefinition,
		ProfileImmutableRecord,
		ProfileLongRunningOperation,
		ProfileTransientRequestResult,
		ProfileEmbeddedValue,
		ProfileListEnvelope,
	}
}

// Valid reports whether p is one of the eight Matrix A profiles.
func (p Profile) Valid() bool {
	switch p {
	case ProfileManagedResource,
		ProfileObservedExternalResource,
		ProfileVersionedDefinition,
		ProfileImmutableRecord,
		ProfileLongRunningOperation,
		ProfileTransientRequestResult,
		ProfileEmbeddedValue,
		ProfileListEnvelope:
		return true
	default:
		return false
	}
}

// Boundary is the Matrix C1 API-boundary vocabulary (F12-BOUNDARY-001).
// Every schema declares exactly one boundary.
type Boundary string

const (
	BoundaryCustomerFacing       Boundary = "customer-facing"
	BoundaryOperatorFacing       Boundary = "operator-facing"
	BoundaryInternalEngineFacing Boundary = "internal-engine-facing"
	BoundaryAdapterFacing        Boundary = "adapter-facing"
	BoundaryPluginFacing         Boundary = "plugin-facing"
	BoundaryGovernanceOnly       Boundary = "governance-only"
)

// AllBoundaries returns the closed Matrix C1 boundary set in stable order.
func AllBoundaries() []Boundary {
	return []Boundary{
		BoundaryCustomerFacing,
		BoundaryOperatorFacing,
		BoundaryInternalEngineFacing,
		BoundaryAdapterFacing,
		BoundaryPluginFacing,
		BoundaryGovernanceOnly,
	}
}

// Valid reports whether b is one of the six Matrix C1 boundaries.
func (b Boundary) Valid() bool {
	switch b {
	case BoundaryCustomerFacing,
		BoundaryOperatorFacing,
		BoundaryInternalEngineFacing,
		BoundaryAdapterFacing,
		BoundaryPluginFacing,
		BoundaryGovernanceOnly:
		return true
	default:
		return false
	}
}

// Stability is the maturity / compatibility expectation vocabulary
// (x-sovrunn-stability; F12-NAMING-006).
type Stability string

const (
	StabilityAlpha  Stability = "alpha"
	StabilityBeta   Stability = "beta"
	StabilityStable Stability = "stable"
)

// AllStabilities returns the closed stability set in stable order.
func AllStabilities() []Stability {
	return []Stability{
		StabilityAlpha,
		StabilityBeta,
		StabilityStable,
	}
}

// Valid reports whether s is an approved stability value.
func (s Stability) Valid() bool {
	switch s {
	case StabilityAlpha, StabilityBeta, StabilityStable:
		return true
	default:
		return false
	}
}

// DataClassification is the closed F12-SEC-002 classification vocabulary.
// Exactly seven values are permitted.
type DataClassification string

const (
	ClassPublic               DataClassification = "Public"
	ClassCustomerVisible      DataClassification = "Customer-visible"
	ClassTenantConfidential   DataClassification = "Tenant-confidential"
	ClassOperatorConfidential DataClassification = "Operator-confidential"
	ClassInternal             DataClassification = "Internal"
	ClassSensitive            DataClassification = "Sensitive"
	ClassSecretReferenceOnly  DataClassification = "Secret-reference-only"
)

// AllDataClassifications returns the closed F12-SEC-002 set in stable order.
func AllDataClassifications() []DataClassification {
	return []DataClassification{
		ClassPublic,
		ClassCustomerVisible,
		ClassTenantConfidential,
		ClassOperatorConfidential,
		ClassInternal,
		ClassSensitive,
		ClassSecretReferenceOnly,
	}
}

// Valid reports whether c is one of the seven approved classifications.
func (c DataClassification) Valid() bool {
	switch c {
	case ClassPublic,
		ClassCustomerVisible,
		ClassTenantConfidential,
		ClassOperatorConfidential,
		ClassInternal,
		ClassSensitive,
		ClassSecretReferenceOnly:
		return true
	default:
		return false
	}
}
