package model

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Canonical HostingLocation TypeMeta values (VS0-SCHEMA-011).
const (
	APIVersionHostingLocation = "infrastructure.sovrunn.io/v1alpha1"
	KindHostingLocation       = "HostingLocation"
)

// HostingLocationPhase is the closed status.phase vocabulary for HostingLocation.
type HostingLocationPhase string

const (
	HostingLocationPhaseActive HostingLocationPhase = "Active"
)

// HostingLocation is the FEATURE-0015 topology root location resource
// (VS0-SCHEMA-011; REQ-F15-04; DEC-0041). It has no parent reference.
type HostingLocation struct {
	apimeta.TypeMeta                       // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta    `json:"metadata"`
	Spec             HostingLocationSpec   `json:"spec"`
	Status           HostingLocationStatus `json:"status,omitempty"`
}

// HostingLocationSpec carries country/locality identity. CountryCode is an
// assigned uppercase ISO-3166-1 alpha-2 code. There is no parent reference.
type HostingLocationSpec struct {
	CountryCode            string `json:"countryCode"`
	Locality               string `json:"locality"`
	AdministrativeAreaCode string `json:"administrativeAreaCode,omitempty"`
	Description            string `json:"description,omitempty"`
}

// HostingLocationStatus is api-server-owned observed state.
type HostingLocationStatus struct {
	Phase      HostingLocationPhase `json:"phase,omitempty"`
	Conditions []apicond.Condition  `json:"conditions,omitempty"`
}
