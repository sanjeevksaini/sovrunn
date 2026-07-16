package resources

// OrganizationUnit is a delegated governance boundary under an Organization.
// It follows the canonical metadata/spec/status shape so that it can evolve
// toward Kubernetes-compatible desired-state reconciliation.
type OrganizationUnit struct {
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   Metadata               `json:"metadata"`
	Spec       OrganizationUnitSpec   `json:"spec"`
	Status     OrganizationUnitStatus `json:"status"`
}

// OrganizationUnitSpec is the desired-state payload for an OrganizationUnit.
// OrganizationName is required and immutable after creation.
type OrganizationUnitSpec struct {
	OrganizationName string `json:"organizationName"`
	Description      string `json:"description,omitempty"`
}

// OrganizationUnitStatus is system-owned observed state.
// Clients must NOT submit status in create/update requests.
type OrganizationUnitStatus struct {
	Phase   string `json:"phase"`
	Message string `json:"message,omitempty"`
}

// APIVersion and Kind constants — set by server, never from client.
const (
	OUAPIVersion = "platform.sovrunn.io/v1alpha1"
	OUKind       = "OrganizationUnit"
)
