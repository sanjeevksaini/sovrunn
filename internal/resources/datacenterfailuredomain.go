package resources

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// ProviderDatacenterRef is the typed constrained parent-reference alias for
// DatacenterFailureDomain.spec.providerDatacenterRef (DD-08, F14-REQ-09,
// F14-REQ-21). It reuses the FEATURE-0012 TypedRef grammar unchanged
// (apiVersion, kind, name, optional resolved uid) and is constrained to
// kind ProviderDatacenter. Kind and immutability enforcement belong to
// later validation stages (design §7.1); this type only shapes the single
// typed reference.
type ProviderDatacenterRef struct {
	apiref.TypedRef // anonymous embed promotes apiVersion/kind/name/uid
}

// DatacenterFailureDomain is the FEATURE-0014 failure-domain ManagedResource
// (F14-REQ-06, F14-REQ-08, F14-REQ-09, F14-REQ-21, F14-REQ-22).
// Governance scope is the containing Provider via immutable
// metadata.scopeRef. Topology parenthood is solely
// spec.providerDatacenterRef (DD-08); no second topology authority exists.
// This kind records identity and containment only — no correlated-risk,
// availability, quorum, resilience, or placement field.
type DatacenterFailureDomain struct {
	apimeta.TypeMeta                               // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta            `json:"metadata"`
	Spec             DatacenterFailureDomainSpec   `json:"spec"`
	Status           DatacenterFailureDomainStatus `json:"status,omitempty"`
}

// DatacenterFailureDomainSpec carries exactly one required immutable typed
// immediate-parent reference (DD-02, DD-08). No second parent, connectivity,
// capacity, capability, resilience, or native field is authorized.
type DatacenterFailureDomainSpec struct {
	ProviderDatacenterRef ProviderDatacenterRef `json:"providerDatacenterRef"`
}

// DatacenterFailureDomainStatus is system-owned observed state. It holds
// observedGeneration and current-fact conditions Valid and
// TopologyComplete only (DD-03). It is not history.
type DatacenterFailureDomainStatus struct {
	ObservedGeneration int64               `json:"observedGeneration,omitempty"`
	Conditions         []apicond.Condition `json:"conditions,omitempty"`
}
