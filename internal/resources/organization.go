package resources

// Organization is the top-level governance and ownership boundary in
// Sovrunn. It follows the canonical metadata/spec/status shape so that
// it can evolve toward Kubernetes-compatible desired-state reconciliation.
type Organization struct {
	APIVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Metadata   Metadata           `json:"metadata"`
	Spec       OrganizationSpec   `json:"spec"`
	Status     OrganizationStatus `json:"status"`
}

// Metadata carries identity and classification fields for any Sovrunn
// resource. Users may author all four fields; system-owned fields
// (createdAt, resourceVersion, etc.) are not included in Phase 1.
type Metadata struct {
	Name        string            `json:"name"`
	DisplayName string            `json:"displayName,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// OrganizationSpec is the desired-state payload for an Organization.
// All fields are optional in Phase 1 except they are accepted
// and stored without server-side coercion.
type OrganizationSpec struct {
	Description          string   `json:"description,omitempty"`
	SovereignLocations   []string `json:"sovereignLocations,omitempty"`
	DefaultPolicyProfile string   `json:"defaultPolicyProfile,omitempty"`
}

// OrganizationStatus is system-owned observed state.
// Clients must NOT submit status in create/update requests.
type OrganizationStatus struct {
	Phase   string `json:"phase"`
	Message string `json:"message,omitempty"`
}

// OrganizationPhase constants for the Status.Phase field.
const (
	PhaseActive   = "Active"
	PhaseInactive = "Inactive"
	PhaseDeleting = "Deleting"
	PhaseFailed   = "Failed"
)

// APIVersion and Kind constants — set by server, never from client.
const (
	OrgAPIVersion    = "platform.sovrunn.io/v1alpha1"
	OrgKind          = "Organization"
	OrganizationKind = "Organization"
)
