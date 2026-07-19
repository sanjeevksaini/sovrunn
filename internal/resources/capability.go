package resources

// Capability is a global platform registry resource declaring a specific
// lifecycle action supported by a plugin for a given ServiceClass.
// Identity: metadata.name (simple key). It is NOT scoped to the
// Organization/OrganizationUnit/Tenant/Project hierarchy. Follows the
// canonical metadata/spec/status shape. Capability is immutable after create
// (no update); delete and recreate instead.
type Capability struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   Metadata         `json:"metadata"`
	Spec       CapabilitySpec   `json:"spec"`
	Status     CapabilityStatus `json:"status"`
}

// CapabilitySpec is the desired-state payload for a Capability. PluginRef,
// ServiceClassRef, and Operation are required. Supported defaults to false
// (Go zero-value). Description is optional. No secrets may be stored.
type CapabilitySpec struct {
	PluginRef       string `json:"pluginRef"`
	ServiceClassRef string `json:"serviceClassRef"`
	Operation       string `json:"operation"`
	Supported       bool   `json:"supported"`
	Description     string `json:"description,omitempty"`
}

// CapabilityStatus is system-owned observed state.
// Clients must NOT submit status in create requests.
type CapabilityStatus struct {
	Phase   string `json:"phase"`
	Message string `json:"message,omitempty"`
}

// APIVersion and Kind constants are set by the server, never from client input.
const (
	CapabilityAPIVersion = "platform.sovrunn.io/v1alpha1"
	CapabilityKind       = "Capability"
)

// CapabilityOperation constants.
const (
	CapOpValidate          = "Validate"
	CapOpPlan              = "Plan"
	CapOpProvision         = "Provision"
	CapOpConfigure         = "Configure"
	CapOpBind              = "Bind"
	CapOpObserve           = "Observe"
	CapOpScale             = "Scale"
	CapOpUpgrade           = "Upgrade"
	CapOpBackup            = "Backup"
	CapOpRestore           = "Restore"
	CapOpRotateCredentials = "RotateCredentials" // #nosec G101 -- operation name, not a hardcoded credential.
	CapOpUnbind            = "Unbind"
	CapOpDelete            = "Delete"
)
