package model

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Canonical Datacenter TypeMeta values (VS0-SCHEMA-012).
const (
	APIVersionDatacenter = "infrastructure.sovrunn.io/v1alpha1"
	KindDatacenter       = "Datacenter"
)

// DatacenterPhase is the closed status.phase vocabulary for Datacenter.
type DatacenterPhase string

const (
	DatacenterPhaseActive DatacenterPhase = "Active"
)

// Datacenter is the FEATURE-0015 topology datacenter resource
// (VS0-SCHEMA-012; REQ-F15-04; DEC-0041).
type Datacenter struct {
	apimeta.TypeMeta                    // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta `json:"metadata"`
	Spec             DatacenterSpec     `json:"spec"`
	Status           DatacenterStatus   `json:"status,omitempty"`
}

// DatacenterSpec carries an immutable parent HostingLocation reference.
type DatacenterSpec struct {
	HostingLocationRef apimeta.TypedRef `json:"hostingLocationRef"`
	Description        string           `json:"description,omitempty"`
}

// DatacenterStatus is api-server-owned observed state.
type DatacenterStatus struct {
	Phase      DatacenterPhase     `json:"phase,omitempty"`
	Conditions []apicond.Condition `json:"conditions,omitempty"`
}
