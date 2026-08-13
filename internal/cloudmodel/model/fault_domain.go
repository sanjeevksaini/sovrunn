package model

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Canonical FaultDomain TypeMeta values (VS0-SCHEMA-013).
const (
	APIVersionFaultDomain = "infrastructure.sovrunn.io/v1alpha1"
	KindFaultDomain       = "FaultDomain"
)

// FaultDomainPhase is the closed status.phase vocabulary for FaultDomain.
type FaultDomainPhase string

const (
	FaultDomainPhaseActive FaultDomainPhase = "Active"
)

// FaultDomain is the FEATURE-0015 topology fault-domain resource
// (VS0-SCHEMA-013; REQ-F15-04; DEC-0041).
type FaultDomain struct {
	apimeta.TypeMeta                    // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta `json:"metadata"`
	Spec             FaultDomainSpec    `json:"spec"`
	Status           FaultDomainStatus  `json:"status,omitempty"`
}

// FaultDomainSpec carries an immutable parent Datacenter reference.
type FaultDomainSpec struct {
	DatacenterRef apimeta.TypedRef `json:"datacenterRef"`
	Description   string           `json:"description,omitempty"`
}

// FaultDomainStatus is api-server-owned observed state.
type FaultDomainStatus struct {
	Phase      FaultDomainPhase    `json:"phase,omitempty"`
	Conditions []apicond.Condition `json:"conditions,omitempty"`
}
