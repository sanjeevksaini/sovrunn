package resources

// Project is the workload or environment grouping boundary under a Tenant. It
// follows the canonical metadata/spec/status shape.
type Project struct {
	APIVersion string        `json:"apiVersion"`
	Kind       string        `json:"kind"`
	Metadata   Metadata      `json:"metadata"`
	Spec       ProjectSpec   `json:"spec"`
	Status     ProjectStatus `json:"status"`
}

// ProjectSpec is the desired-state payload for a Project. OrganizationName,
// OrganizationUnitName, and TenantName identify the immutable parent.
type ProjectSpec struct {
	OrganizationName     string `json:"organizationName"`
	OrganizationUnitName string `json:"organizationUnitName"`
	TenantName           string `json:"tenantName"`
	Description          string `json:"description,omitempty"`
}

// ProjectStatus is system-owned observed state.
// Clients must NOT submit status in create/update requests.
type ProjectStatus struct {
	Phase   string `json:"phase"`
	Message string `json:"message,omitempty"`
}

// APIVersion and Kind constants are set by the server, never from client input.
const (
	ProjectAPIVersion = "platform.sovrunn.io/v1alpha1"
	ProjectKind       = "Project"
)
