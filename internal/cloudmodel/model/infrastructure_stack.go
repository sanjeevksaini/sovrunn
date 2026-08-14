package model

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Canonical InfrastructureStack TypeMeta values (VS0-SCHEMA-014).
const (
	APIVersionInfrastructureStack = "infrastructure.sovrunn.io/v1alpha1"
	KindInfrastructureStack       = "InfrastructureStack"
)

// InfrastructureStackPhase is the closed status.phase vocabulary for
// InfrastructureStack.
type InfrastructureStackPhase string

const (
	InfrastructureStackPhaseActive InfrastructureStackPhase = "Active"
)

// InfrastructureStack is the FEATURE-0015 topology stack resource
// (VS0-SCHEMA-014; REQ-F15-04; DEC-0041). It is distinct from the retained
// FEATURE-0014 alpha InfrastructureStack in internal/resources.
type InfrastructureStack struct {
	apimeta.TypeMeta                           // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta        `json:"metadata"`
	Spec             InfrastructureStackSpec   `json:"spec"`
	Status           InfrastructureStackStatus `json:"status,omitempty"`
}

// InfrastructureStackSpec carries an immutable parent FaultDomain reference.
type InfrastructureStackSpec struct {
	FaultDomainRef apimeta.TypedRef `json:"faultDomainRef"`
	Description    string           `json:"description,omitempty"`
}

// InfrastructureStackStatus is api-server-owned observed state.
type InfrastructureStackStatus struct {
	Phase      InfrastructureStackPhase `json:"phase,omitempty"`
	Conditions []apicond.Condition      `json:"conditions,omitempty"`
}
