package resources

// Plugin is a global platform registry resource declaring an implementation
// unit that performs lifecycle operations for a service family or provider.
// Identity: metadata.name (simple key). It is NOT scoped to the
// Organization/OrganizationUnit/Tenant/Project hierarchy. Follows the
// canonical metadata/spec/status shape.
type Plugin struct {
	APIVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   Metadata     `json:"metadata"`
	Spec       PluginSpec   `json:"spec"`
	Status     PluginStatus `json:"status"`
}

// PluginSpec is the desired-state payload for a Plugin. PluginType, Version,
// ServiceClassRefs, and DeploymentMode are required; Description and Tags are
// optional catalog metadata. No secrets may be stored in any field.
type PluginSpec struct {
	PluginType       string   `json:"pluginType"`
	Version          string   `json:"version"`
	ServiceClassRefs []string `json:"serviceClassRefs"`
	DeploymentMode   string   `json:"deploymentMode"`
	Description      string   `json:"description,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

// PluginStatus is system-owned observed state.
// Clients must NOT submit status in create/update requests.
type PluginStatus struct {
	Phase   string `json:"phase"`
	Message string `json:"message,omitempty"`
}

// APIVersion and Kind constants are set by the server, never from client input.
const (
	PluginAPIVersion = "platform.sovrunn.io/v1alpha1"
	PluginKind       = "Plugin"
)

// PluginType constants.
const (
	PluginTypeDStoreOps  = "dStoreOps"
	PluginTypeCacheOps   = "cacheOps"
	PluginTypeStreamOps  = "streamOps"
	PluginTypeObjectOps  = "objectOps"
	PluginTypeGatewayOps = "gatewayOps"
	PluginTypeFaasOps    = "faasOps"
	PluginTypeLBOps      = "lbOps"
	PluginTypeK8sOps     = "k8sOps"
	PluginTypeBigDataOps = "bigDataOps"
	PluginTypeSdeOps     = "sdeOps"
)

// DeploymentMode constants. Phase 1 accepts only CompiledIn.
const (
	DeploymentModeCompiledIn = "compiled-in"
)
