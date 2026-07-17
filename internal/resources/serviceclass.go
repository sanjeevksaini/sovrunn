package resources

// ServiceClass is a global platform catalog resource describing a service type
// (for example PostgreSQL or Redis). It is a catalog definition only and does
// not provision, bind, or execute anything. ServiceClass is NOT scoped to the
// Organization/OrganizationUnit/Tenant/Project hierarchy; its identity is
// metadata.name. It follows the canonical metadata/spec/status shape.
type ServiceClass struct {
	APIVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Metadata   Metadata           `json:"metadata"`
	Spec       ServiceClassSpec   `json:"spec"`
	Status     ServiceClassStatus `json:"status"`
}

// ServiceClassSpec is the desired-state payload for a ServiceClass. Category and
// Lifecycle are required; the remaining fields are optional catalog metadata.
type ServiceClassSpec struct {
	DisplayName     string   `json:"displayName,omitempty"`
	Description     string   `json:"description,omitempty"`
	Category        string   `json:"category"`
	Provider        string   `json:"provider,omitempty"`
	Lifecycle       string   `json:"lifecycle"`
	DefaultPlanName string   `json:"defaultPlanName,omitempty"`
	Tags            []string `json:"tags,omitempty"`
}

// ServiceClassStatus is system-owned observed state.
// Clients must NOT submit status in create/update requests.
type ServiceClassStatus struct {
	Phase   string `json:"phase"`
	Message string `json:"message,omitempty"`
}

// APIVersion and Kind constants are set by the server, never from client input.
const (
	ServiceClassAPIVersion = "platform.sovrunn.io/v1alpha1"
	ServiceClassKind       = "ServiceClass"
)

// ServiceClass category constants.
const (
	CategoryDatabase      = "Database"
	CategoryCache         = "Cache"
	CategoryObjectStorage = "ObjectStorage"
	CategoryStream        = "Stream"
	CategoryGateway       = "Gateway"
	CategoryFunction      = "Function"
	CategoryAnalytics     = "Analytics"
	CategoryOther         = "Other"
)

// Lifecycle constants shared by ServiceClass and ServicePlan.
const (
	LifecyclePreview    = "Preview"
	LifecycleActive     = "Active"
	LifecycleDeprecated = "Deprecated"
	LifecycleRetired    = "Retired"
)
