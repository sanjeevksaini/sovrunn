package resources

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// DatacenterFailureDomainRef is the typed constrained parent-reference alias
// for InfrastructureStack.spec.datacenterFailureDomainRef (DD-08, F14-REQ-09,
// F14-REQ-16, F14-REQ-21). It reuses the FEATURE-0012 TypedRef grammar
// unchanged (apiVersion, kind, name, optional resolved uid) and is
// constrained to kind DatacenterFailureDomain. Kind and immutability
// enforcement belong to later validation stages (design §7.1); this type
// only shapes the single typed reference. A stack must not carry more than
// one failure-domain reference.
type DatacenterFailureDomainRef struct {
	apiref.TypedRef // anonymous embed promotes apiVersion/kind/name/uid
}

// InfrastructureStack is the FEATURE-0014 leaf ManagedResource
// (F14-REQ-03, F14-REQ-06, F14-REQ-08, F14-REQ-09, F14-REQ-15, F14-REQ-16,
// F14-REQ-17, F14-REQ-18, F14-REQ-21, F14-REQ-22). Governance scope is the
// containing Provider via immutable metadata.scopeRef. Topology parenthood
// is solely spec.datacenterFailureDomainRef (DD-08); no second topology
// authority exists. Distinct immutable identity is metadata.uid plus scoped
// identity; shared technology never implies shared identity.
type InfrastructureStack struct {
	apimeta.TypeMeta                           // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta        `json:"metadata"`
	Spec             InfrastructureStackSpec   `json:"spec"`
	Status           InfrastructureStackStatus `json:"status,omitempty"`
}

// InfrastructureStackSpec carries exactly one required immutable typed
// immediate-parent reference plus an optional bounded descriptive
// technology string (DD-02, DD-05, DD-08). Technology is not identity, a
// closed enum, or a lookup/reference key; length and character-set
// validation belong to a later task. No second failure-domain reference,
// connectivity, capacity, capability, or native field is authorized.
type InfrastructureStackSpec struct {
	DatacenterFailureDomainRef DatacenterFailureDomainRef `json:"datacenterFailureDomainRef"`
	Technology                 string                     `json:"technology,omitempty"`
}

// InfrastructureStackStatus is system-owned observed state. It holds
// observedGeneration and current-fact conditions Valid and
// TopologyComplete only (DD-03). It is not history.
type InfrastructureStackStatus struct {
	ObservedGeneration int64               `json:"observedGeneration,omitempty"`
	Conditions         []apicond.Condition `json:"conditions,omitempty"`
}
