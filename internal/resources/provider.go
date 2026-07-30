package resources

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Shared FEATURE-0014 fabric API identity (DD-01). Declared once in this
// root file for all five kinds; do not duplicate elsewhere.
const (
	FabricAPIVersion = "fabric.sovrunn.io/v1alpha1"

	KindProvider                = "Provider"
	KindProviderLocation        = "ProviderLocation"
	KindProviderDatacenter      = "ProviderDatacenter"
	KindDatacenterFailureDomain = "DatacenterFailureDomain"
	KindInfrastructureStack     = "InfrastructureStack"
)

// Provider is the FEATURE-0014 supply-boundary ManagedResource (F14-REQ-01,
// F14-REQ-04, F14-REQ-05, F14-REQ-21, F14-REQ-22). Governance scope is
// expressed only through inherited metadata.scopeRef (Platform or one
// Organization). Spec carries no domain fields (DD-02).
type Provider struct {
	apimeta.TypeMeta                    // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta `json:"metadata"`
	Spec             ProviderSpec       `json:"spec"`
	Status           ProviderStatus     `json:"status,omitempty"`
}

// ProviderSpec is intentionally empty: identity and governance scope live
// in metadata; no domain, owner, native, connectivity, or capacity field
// is authorized (DD-02, F14-AD-003, F14-AD-017–F14-AD-019).
type ProviderSpec struct{}

// ProviderStatus is system-owned observed state. It holds
// observedGeneration and current-fact conditions Valid and
// TopologyComplete only (DD-03). It is not history.
type ProviderStatus struct {
	ObservedGeneration int64               `json:"observedGeneration,omitempty"`
	Conditions         []apicond.Condition `json:"conditions,omitempty"`
}
