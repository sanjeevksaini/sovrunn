package resources

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// ProviderLocationRef is the typed constrained parent-reference alias for
// ProviderDatacenter.spec.providerLocationRef (DD-08, F14-REQ-09, F14-REQ-21).
// It reuses the FEATURE-0012 TypedRef grammar unchanged (apiVersion, kind,
// name, optional resolved uid) and is constrained to kind ProviderLocation.
// Kind and immutability enforcement belong to later validation stages
// (design §7.1); this type only shapes the single typed reference.
type ProviderLocationRef struct {
	apiref.TypedRef // anonymous embed promotes apiVersion/kind/name/uid
}

// ProviderDatacenter is the FEATURE-0014 datacenter-level ManagedResource
// (F14-REQ-06, F14-REQ-08, F14-REQ-09, F14-REQ-21, F14-REQ-22).
// Governance scope is the containing Provider via immutable
// metadata.scopeRef. Topology parenthood is solely
// spec.providerLocationRef (DD-08); no second topology authority exists.
type ProviderDatacenter struct {
	apimeta.TypeMeta                          // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta       `json:"metadata"`
	Spec             ProviderDatacenterSpec   `json:"spec"`
	Status           ProviderDatacenterStatus `json:"status,omitempty"`
}

// ProviderDatacenterSpec carries exactly one required immutable typed
// immediate-parent reference (DD-02, DD-08). No second parent, connectivity,
// capacity, capability, or native field is authorized.
type ProviderDatacenterSpec struct {
	ProviderLocationRef ProviderLocationRef `json:"providerLocationRef"`
}

// ProviderDatacenterStatus is system-owned observed state. It holds
// observedGeneration and current-fact conditions Valid and
// TopologyComplete only (DD-03). It is not history.
type ProviderDatacenterStatus struct {
	ObservedGeneration int64               `json:"observedGeneration,omitempty"`
	Conditions         []apicond.Condition `json:"conditions,omitempty"`
}
