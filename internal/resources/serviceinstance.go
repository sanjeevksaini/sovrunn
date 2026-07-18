package resources

// ServiceInstance is a tenant/project-scoped requested service. It records
// desired configuration only and does not provision infrastructure in Phase 1.
type ServiceInstance struct {
	APIVersion string                `json:"apiVersion"`
	Kind       string                `json:"kind"`
	Metadata   Metadata              `json:"metadata"`
	Spec       ServiceInstanceSpec   `json:"spec"`
	Status     ServiceInstanceStatus `json:"status"`
}

// ServiceInstanceSpec is the desired-state payload for a ServiceInstance.
// Parameters must not contain secrets or credentials.
type ServiceInstanceSpec struct {
	OrganizationRef     string            `json:"organizationRef"`
	OrganizationUnitRef string            `json:"organizationUnitRef,omitempty"`
	TenantRef           string            `json:"tenantRef"`
	ProjectRef          string            `json:"projectRef"`
	ServiceClassRef     string            `json:"serviceClassRef"`
	ServicePlanRef      string            `json:"servicePlanRef"`
	Parameters          map[string]string `json:"parameters,omitempty"`
}

// ServiceInstanceStatus is system-owned observed state.
// Clients must NOT submit status in create/update requests.
type ServiceInstanceStatus struct {
	Phase   string `json:"phase"`
	Message string `json:"message,omitempty"`
}

// APIVersion and Kind constants are set by the server, never from client input.
const (
	ServiceInstanceAPIVersion = "platform.sovrunn.io/v1alpha1"
	ServiceInstanceKind       = "ServiceInstance"
)
