package model

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Canonical CloudProviderParticipation TypeMeta values (VS0-SCHEMA-010).
const (
	APIVersionCloudProviderParticipation = "governance.sovrunn.io/v1alpha1"
	KindCloudProviderParticipation       = "CloudProviderParticipation"
)

// ParticipationEnvironment is the closed create-time environment vocabulary.
type ParticipationEnvironment string

const (
	ParticipationEnvironmentDevelopment ParticipationEnvironment = "development"
)

// ParticipationPhase is the closed status.phase vocabulary for
// CloudProviderParticipation (VS0-STATE-001).
type ParticipationPhase string

const (
	ParticipationPhasePending     ParticipationPhase = "Pending"
	ParticipationPhaseActive      ParticipationPhase = "Active"
	ParticipationPhaseRejected    ParticipationPhase = "Rejected"
	ParticipationPhaseWithdrawn   ParticipationPhase = "Withdrawn"
	ParticipationPhaseExpired     ParticipationPhase = "Expired"
	ParticipationPhaseSuspended   ParticipationPhase = "Suspended"
	ParticipationPhaseTerminating ParticipationPhase = "Terminating"
	ParticipationPhaseTerminated  ParticipationPhase = "Terminated"
)

// CloudProviderParticipation is the FEATURE-0015 participation agreement
// between a CloudPlatform and a CloudProvider (VS0-SCHEMA-010; REQ-F15-03; DEC-0054).
type CloudProviderParticipation struct {
	apimeta.TypeMeta                                  // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta               `json:"metadata"`
	Spec             CloudProviderParticipationSpec   `json:"spec"`
	Status           CloudProviderParticipationStatus `json:"status,omitempty"`
}

// CloudProviderParticipationSpec carries immutable create-time references and
// environment. FEATURE-0021 fields are intentionally absent.
type CloudProviderParticipationSpec struct {
	CloudPlatformRef apimeta.TypedRef         `json:"cloudPlatformRef"`
	CloudProviderRef apimeta.TypedRef         `json:"cloudProviderRef"`
	Environment      ParticipationEnvironment `json:"environment"`
}

// CloudProviderParticipationStatus is api-server-owned lifecycle state with
// independent platform and provider suspension holds.
type CloudProviderParticipationStatus struct {
	Phase             ParticipationPhase  `json:"phase,omitempty"`
	PlatformSuspended bool                `json:"platformSuspended"`
	ProviderSuspended bool                `json:"providerSuspended"`
	RequestExpiresAt  string              `json:"requestExpiresAt,omitempty"`
	Conditions        []apicond.Condition `json:"conditions,omitempty"`
}
