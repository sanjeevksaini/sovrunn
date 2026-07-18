package resources

// ServiceBinding records how a consumer uses a ServiceInstance. It contains
// only a secret reference; no credential values are stored in this resource.
type ServiceBinding struct {
	APIVersion string               `json:"apiVersion"`
	Kind       string               `json:"kind"`
	Metadata   Metadata             `json:"metadata"`
	Spec       ServiceBindingSpec   `json:"spec"`
	Status     ServiceBindingStatus `json:"status"`
}

// ServiceBindingSpec is the desired-state payload for a ServiceBinding.
type ServiceBindingSpec struct {
	ServiceInstanceRef string       `json:"serviceInstanceRef"`
	ConsumerRef        *ConsumerRef `json:"consumerRef"`
	BindingType        string       `json:"bindingType"`
}

// ConsumerRef identifies the consumer of a ServiceBinding.
type ConsumerRef struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

// ServiceBindingStatus is system-owned observed state.
// Clients must NOT submit status in create requests.
type ServiceBindingStatus struct {
	Phase     string `json:"phase"`
	Message   string `json:"message,omitempty"`
	SecretRef string `json:"secretRef,omitempty"`
}

// APIVersion and Kind constants are set by the server, never from client input.
const (
	ServiceBindingAPIVersion = "platform.sovrunn.io/v1alpha1"
	ServiceBindingKind       = "ServiceBinding"
)

// Allowed ServiceBinding types for Phase 1.
const (
	BindingTypeCredentials = "credentials"
)
