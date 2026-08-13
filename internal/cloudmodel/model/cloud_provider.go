package model

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Canonical CloudProvider TypeMeta values (VS0-SCHEMA-009).
const (
	APIVersionCloudProvider = "core.sovrunn.io/v1alpha1"
	KindCloudProvider       = "CloudProvider"
)

// CloudProviderPhase is the closed status.phase vocabulary for CloudProvider.
type CloudProviderPhase string

const (
	CloudProviderPhaseActive CloudProviderPhase = "Active"
)

// CloudProvider is the FEATURE-0015 supply-boundary cloud provider resource
// (VS0-SCHEMA-009; REQ-F15-02; DEC-0037).
type CloudProvider struct {
	apimeta.TypeMeta                     // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta  `json:"metadata"`
	Spec             CloudProviderSpec   `json:"spec"`
	Status           CloudProviderStatus `json:"status,omitempty"`
}

// CloudProviderSpec carries operating markets and optional display name.
// OperatingMarkets values are assigned uppercase ISO-3166-1 alpha-2 codes.
type CloudProviderSpec struct {
	DisplayName      string   `json:"displayName,omitempty"`
	OperatingMarkets []string `json:"operatingMarkets"`
}

// CloudProviderStatus is api-server-owned observed state.
type CloudProviderStatus struct {
	Phase      CloudProviderPhase  `json:"phase,omitempty"`
	Conditions []apicond.Condition `json:"conditions,omitempty"`
}
