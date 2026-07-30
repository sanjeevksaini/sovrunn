package resources

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// ProviderLocation is the FEATURE-0014 location-level ManagedResource
// (F14-REQ-02, F14-REQ-06, F14-REQ-21, F14-REQ-22, F14-REQ-27).
// The immediate parent Provider is identified only by the inherited
// immutable metadata.scopeRef; there is no separate parent field
// (design §4.2, DD-08).
type ProviderLocation struct {
	apimeta.TypeMeta                        // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta     `json:"metadata"`
	Spec             ProviderLocationSpec   `json:"spec"`
	Status           ProviderLocationStatus `json:"status,omitempty"`
}

// ProviderLocationSpec carries only the optional descriptive geo
// descriptor (DD-02, DD-04). No parent, alias, connectivity, capacity,
// or native field is authorized.
type ProviderLocationSpec struct {
	Geo *ProviderLocationGeo `json:"geo,omitempty"`
}

// ProviderLocationGeo is optional descriptive declared topology only
// (DD-04, F14-REQ-27). It carries no residency, compliance,
// authoritative-assignment, or placement meaning. Syntax and
// country-prefix validation belong to a later task; this type is
// representation only and performs no dataset or membership lookup.
type ProviderLocationGeo struct {
	CountryCode     string `json:"countryCode"`
	SubdivisionCode string `json:"subdivisionCode,omitempty"`
}

// ProviderLocationStatus is system-owned observed state. It holds
// observedGeneration and current-fact conditions Valid and
// TopologyComplete only (DD-03). It is not history.
type ProviderLocationStatus struct {
	ObservedGeneration int64               `json:"observedGeneration,omitempty"`
	Conditions         []apicond.Condition `json:"conditions,omitempty"`
}
