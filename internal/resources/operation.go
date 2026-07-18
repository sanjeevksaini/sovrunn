package resources

// Operation is an immutable control-plane record of a single lifecycle action.
// It records resource references only, not raw request bodies or secrets.
type Operation struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   Metadata        `json:"metadata"`
	Spec       OperationSpec   `json:"spec"`
	Status     OperationStatus `json:"status"`
}

// OperationSpec records the action and non-sensitive target resource reference.
type OperationSpec struct {
	Type                 string `json:"type"`
	ResourceKind         string `json:"resourceKind"`
	ResourceName         string `json:"resourceName"`
	OrganizationName     string `json:"organizationName,omitempty"`
	OrganizationUnitName string `json:"organizationUnitName,omitempty"`
	TenantName           string `json:"tenantName,omitempty"`
	ProjectName          string `json:"projectName,omitempty"`
	ServiceClassName     string `json:"serviceClassName,omitempty"`
	ServicePlanName      string `json:"servicePlanName,omitempty"`
	PluginName           string `json:"pluginName,omitempty"`
	CapabilityName       string `json:"capabilityName,omitempty"`
	ServiceInstanceName  string `json:"serviceInstanceName,omitempty"`
	ServiceBindingName   string `json:"serviceBindingName,omitempty"`
	Actor                string `json:"actor"`
	RequestID            string `json:"requestId,omitempty"`
}

// OperationStatus is system-owned. Timestamps are RFC3339 UTC strings.
type OperationStatus struct {
	Phase       string `json:"phase"`
	Message     string `json:"message,omitempty"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	CompletedAt string `json:"completedAt,omitempty"`
}

// APIVersion and Kind constants are set by the server.
const (
	OperationAPIVersion = "platform.sovrunn.io/v1alpha1"
	OperationKind       = "Operation"
)

// Operation phase constants. Phase 1 emits Succeeded; the other phases are
// reserved for future workflow states.
const (
	OperationPhasePending   = "Pending"
	OperationPhaseRunning   = "Running"
	OperationPhaseSucceeded = "Succeeded"
	OperationPhaseFailed    = "Failed"
)

// Operation type constants.
const (
	OpCreateOrganization     = "CreateOrganization"
	OpUpdateOrganization     = "UpdateOrganization"
	OpDeleteOrganization     = "DeleteOrganization"
	OpCreateOrganizationUnit = "CreateOrganizationUnit"
	OpUpdateOrganizationUnit = "UpdateOrganizationUnit"
	OpDeleteOrganizationUnit = "DeleteOrganizationUnit"
	OpCreateTenant           = "CreateTenant"
	OpUpdateTenant           = "UpdateTenant"
	OpDeleteTenant           = "DeleteTenant"
	OpCreateProject          = "CreateProject"
	OpUpdateProject          = "UpdateProject"
	OpDeleteProject          = "DeleteProject"
)

// Catalog operation type constants (FEATURE-0006).
const (
	OpCreateServiceClass = "CreateServiceClass"
	OpUpdateServiceClass = "UpdateServiceClass"
	OpDeleteServiceClass = "DeleteServiceClass"
	OpCreateServicePlan  = "CreateServicePlan"
	OpUpdateServicePlan  = "UpdateServicePlan"
	OpDeleteServicePlan  = "DeleteServicePlan"
)

// Plugin and Capability operation type constants (FEATURE-0007).
const (
	OpCreatePlugin     = "CreatePlugin"
	OpUpdatePlugin     = "UpdatePlugin"
	OpDeletePlugin     = "DeletePlugin"
	OpCreateCapability = "CreateCapability"
	OpDeleteCapability = "DeleteCapability"
)

// Service consumption operation type constants (FEATURE-0008).
const (
	OpCreateServiceInstance = "CreateServiceInstance"
	OpUpdateServiceInstance = "UpdateServiceInstance"
	OpDeleteServiceInstance = "DeleteServiceInstance"
	OpCreateServiceBinding  = "CreateServiceBinding"
	OpDeleteServiceBinding  = "DeleteServiceBinding"
)
