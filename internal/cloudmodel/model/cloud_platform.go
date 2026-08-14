package model

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Canonical CloudPlatform TypeMeta values (VS0-SCHEMA-008).
const (
	APIVersionCloudPlatform = "core.sovrunn.io/v1alpha1"
	KindCloudPlatform       = "CloudPlatform"
)

// CloudPlatformPhase is the closed status.phase vocabulary for CloudPlatform.
type CloudPlatformPhase string

const (
	CloudPlatformPhaseActive CloudPlatformPhase = "Active"
)

// CloudPlatform is the FEATURE-0015 legal-entity cloud platform resource
// (VS0-SCHEMA-008; REQ-F15-01; ADH-2026-045).
type CloudPlatform struct {
	apimeta.TypeMeta                     // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta  `json:"metadata"`
	Spec             CloudPlatformSpec   `json:"spec"`
	Status           CloudPlatformStatus `json:"status,omitempty"`
}

// CloudPlatformSpec carries immutable owner registration and optional description.
type CloudPlatformSpec struct {
	OwnerRegistration OwnerRegistration `json:"ownerRegistration"`
	Description       string            `json:"description,omitempty"`
}

// OwnerRegistration is the immutable legal-entity registration identity
// (ADH-2026-045). JurisdictionCode is an assigned ISO-3166-1 alpha-2 code.
type OwnerRegistration struct {
	LegalName              string `json:"legalName"`
	RegistrationIdentifier string `json:"registrationIdentifier"`
	JurisdictionCode       string `json:"jurisdictionCode"`
}

// CloudPlatformStatus is api-server-owned observed state.
type CloudPlatformStatus struct {
	Phase      CloudPlatformPhase  `json:"phase,omitempty"`
	Conditions []apicond.Condition `json:"conditions,omitempty"`
}
