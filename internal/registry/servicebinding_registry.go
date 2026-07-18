package registry

import (
	"context"
	"sort"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// ServiceBindingRegistryIface is the storage contract for ServiceBinding
// resources. The registry is storage-only: it does not depend on other
// registries and performs no cross-resource reference checks. Those checks
// belong to the API/service layer.
type ServiceBindingRegistryIface interface {
	CreateServiceBinding(ctx context.Context, sb resources.ServiceBinding) (resources.ServiceBinding, error)
	GetServiceBinding(ctx context.Context, name string) (resources.ServiceBinding, error)
	ListServiceBindings(ctx context.Context, serviceInstanceRef string) ([]resources.ServiceBinding, error)
	DeleteServiceBinding(ctx context.Context, name string) error
	CountByServiceInstance(ctx context.Context, instanceName string) (int, error)
}

// ServiceBindingInstanceBlocker is a narrow interface for counting ServiceBindings
// that reference a given ServiceInstance. The existing *ServiceBindingRegistry
// already satisfies it via CountByServiceInstance. It is intended for injection
// into the ServiceInstance handler and is NOT used inside ServiceBindingRegistry.
type ServiceBindingInstanceBlocker interface {
	CountByServiceInstance(ctx context.Context, instanceName string) (int, error)
}

// ServiceBindingRegistry is the Phase 1 in-memory implementation of
// ServiceBindingRegistryIface. All public methods are safe for concurrent
// use. The registry holds no package-level global state and is keyed by
// metadata.name (global uniqueness).
type ServiceBindingRegistry struct {
	mu    sync.RWMutex
	store map[string]resources.ServiceBinding
}

// Compile-time checks that *ServiceBindingRegistry satisfies the interfaces.
var (
	_ ServiceBindingRegistryIface   = (*ServiceBindingRegistry)(nil)
	_ ServiceBindingInstanceBlocker = (*ServiceBindingRegistry)(nil)
)

// NewServiceBindingRegistry returns a ready-to-use registry.
func NewServiceBindingRegistry() *ServiceBindingRegistry {
	return &ServiceBindingRegistry{
		store: make(map[string]resources.ServiceBinding),
	}
}

// deepCopyServiceBinding returns a fully independent copy of sb, duplicating
// Metadata Labels/Annotations maps and the ConsumerRef pointer so callers
// cannot mutate the registry's internal state.
func deepCopyServiceBinding(sb resources.ServiceBinding) resources.ServiceBinding {
	cp := sb
	if sb.Metadata.Labels != nil {
		cp.Metadata.Labels = make(map[string]string, len(sb.Metadata.Labels))
		for k, v := range sb.Metadata.Labels {
			cp.Metadata.Labels[k] = v
		}
	}
	if sb.Metadata.Annotations != nil {
		cp.Metadata.Annotations = make(map[string]string, len(sb.Metadata.Annotations))
		for k, v := range sb.Metadata.Annotations {
			cp.Metadata.Annotations[k] = v
		}
	}
	if sb.Spec.ConsumerRef != nil {
		ref := *sb.Spec.ConsumerRef
		cp.Spec.ConsumerRef = &ref
	}
	return cp
}

// CreateServiceBinding stores a deep copy of sb keyed by metadata.name. It
// returns ErrAlreadyExists if the name is already present.
func (r *ServiceBindingRegistry) CreateServiceBinding(
	ctx context.Context, sb resources.ServiceBinding,
) (resources.ServiceBinding, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[sb.Metadata.Name]; ok {
		return resources.ServiceBinding{}, ErrAlreadyExists
	}
	stored := deepCopyServiceBinding(sb)
	r.store[sb.Metadata.Name] = stored
	return deepCopyServiceBinding(stored), nil
}

// GetServiceBinding returns a deep copy of the stored ServiceBinding
// identified by name, or ErrNotFound if absent.
func (r *ServiceBindingRegistry) GetServiceBinding(
	ctx context.Context, name string,
) (resources.ServiceBinding, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sb, ok := r.store[name]
	if !ok {
		return resources.ServiceBinding{}, ErrNotFound
	}
	return deepCopyServiceBinding(sb), nil
}

// ListServiceBindings returns a new slice of deep copies sorted by
// metadata.name ascending. When serviceInstanceRef is non-empty, only entries
// with matching Spec.ServiceInstanceRef are included. Empty filter strings
// are not applied. Returns a non-nil empty slice when no ServiceBindings match.
func (r *ServiceBindingRegistry) ListServiceBindings(
	ctx context.Context, serviceInstanceRef string,
) ([]resources.ServiceBinding, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]resources.ServiceBinding, 0, len(r.store))
	for _, sb := range r.store {
		if serviceInstanceRef != "" && sb.Spec.ServiceInstanceRef != serviceInstanceRef {
			continue
		}
		items = append(items, deepCopyServiceBinding(sb))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Metadata.Name < items[j].Metadata.Name
	})
	return items, nil
}

// DeleteServiceBinding removes the entry identified by name, or returns
// ErrNotFound if absent.
func (r *ServiceBindingRegistry) DeleteServiceBinding(
	ctx context.Context, name string,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[name]; !ok {
		return ErrNotFound
	}
	delete(r.store, name)
	return nil
}

// CountByServiceInstance returns the number of stored ServiceBindings whose
// Spec.ServiceInstanceRef matches instanceName.
func (r *ServiceBindingRegistry) CountByServiceInstance(
	ctx context.Context, instanceName string,
) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, sb := range r.store {
		if sb.Spec.ServiceInstanceRef == instanceName {
			count++
		}
	}
	return count, nil
}
