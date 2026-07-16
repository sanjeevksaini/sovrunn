package resources

// Tenant is the primary isolation and security boundary under an
// OrganizationUnit. It follows the canonical metadata/spec/status shape.
type Tenant struct {
	APIVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   Metadata     `json:"metadata"`
	Spec       TenantSpec   `json:"spec"`
	Status     TenantStatus `json:"status"`
}

// TenantSpec is the desired-state payload for a Tenant.
// OrganizationName and OrganizationUnitName identify the immutable parent.
type TenantSpec struct {
	OrganizationName     string `json:"organizationName"`
	OrganizationUnitName string `json:"organizationUnitName"`
	Description          string `json:"description,omitempty"`
}

// TenantStatus is system-owned observed state.
// Clients must NOT submit status in create/update requests.
type TenantStatus struct {
	Phase   string `json:"phase"`
	Message string `json:"message,omitempty"`
}

// APIVersion and Kind constants are set by the server, never from client input.
const (
	TenantAPIVersion = "platform.sovrunn.io/v1alpha1"
	TenantKind       = "Tenant"
)
