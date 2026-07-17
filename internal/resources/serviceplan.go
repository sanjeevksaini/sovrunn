package resources

// ServicePlan is a global platform catalog resource describing a tier or shape
// under a ServiceClass. It is a catalog definition only and does not provision,
// bind, or execute anything. ServicePlan is NOT scoped to the
// Organization/OrganizationUnit/Tenant/Project hierarchy; its identity is the
// composite serviceClassName/name. It follows the canonical
// metadata/spec/status shape.
type ServicePlan struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   Metadata          `json:"metadata"`
	Spec       ServicePlanSpec   `json:"spec"`
	Status     ServicePlanStatus `json:"status"`
}

// ServicePlanSpec is the desired-state payload for a ServicePlan.
// ServiceClassName identifies the immutable parent ServiceClass; Tier and
// Lifecycle are required. Parameters are simple string key/value catalog
// settings and must never carry secrets.
type ServicePlanSpec struct {
	ServiceClassName string            `json:"serviceClassName"`
	DisplayName      string            `json:"displayName,omitempty"`
	Description      string            `json:"description,omitempty"`
	Tier             string            `json:"tier"`
	Lifecycle        string            `json:"lifecycle"`
	Parameters       map[string]string `json:"parameters,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
}

// ServicePlanStatus is system-owned observed state.
// Clients must NOT submit status in create/update requests.
type ServicePlanStatus struct {
	Phase   string `json:"phase"`
	Message string `json:"message,omitempty"`
}

// APIVersion and Kind constants are set by the server, never from client input.
const (
	ServicePlanAPIVersion = "platform.sovrunn.io/v1alpha1"
	ServicePlanKind       = "ServicePlan"
)

// ServicePlan tier constants.
const (
	TierDev        = "Dev"
	TierSmall      = "Small"
	TierMedium     = "Medium"
	TierLarge      = "Large"
	TierProduction = "Production"
	TierCustom     = "Custom"
)
