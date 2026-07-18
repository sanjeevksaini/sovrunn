package registry

import (
	"context"
	"sort"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// ServiceInstanceRegistryIface is the storage contract for ServiceInstance
// resources. The registry is storage-only: it does not depend on other
// registries and performs no cross-resource reference checks. Those checks
// belong to the API/service layer.
type ServiceInstanceRegistryIface interface {
	CreateServiceInstance(ctx context.Context, si resources.ServiceInstance) (resources.ServiceInstance, error)
	GetServiceInstance(ctx context.Context, name string) (resources.ServiceInstance, error)
	ListServiceInstances(ctx context.Context, tenantRef, projectRef string) ([]resources.ServiceInstance, error)
	UpdateServiceInstance(ctx context.Context, name string, si resources.ServiceInstance) (resources.ServiceInstance, error)
	DeleteServiceInstance(ctx context.Context, name string) error
	CountByServicePlan(ctx context.Context, serviceClassRef, servicePlanRef string) (int, error)
	CountByProject(ctx context.Context, organizationRef, organizationUnitRef, tenantRef, projectRef string) (int, error)
}

// ServiceInstanceLookup is a narrow interface for verifying ServiceInstance
// existence. The existing *ServiceInstanceRegistry already satisfies it via
// GetServiceInstance. It is intended for injection into the ServiceBinding
// handler and is NOT used inside ServiceInstanceRegistry.
type ServiceInstanceLookup interface {
	GetServiceInstance(ctx context.Context, name string) (resources.ServiceInstance, error)
}

// ServiceInstanceRegistry is the Phase 1 in-memory implementation of
// ServiceInstanceRegistryIface. All public methods are safe for concurrent
// use. The registry holds no package-level global state and is keyed by
// metadata.name (global uniqueness).
type ServiceInstanceRegistry struct {
	mu    sync.RWMutex
	store map[string]resources.ServiceInstance
}

// Compile-time checks that *ServiceInstanceRegistry satisfies the interfaces.
var (
	_ ServiceInstanceRegistryIface = (*ServiceInstanceRegistry)(nil)
	_ ServiceInstanceLookup        = (*ServiceInstanceRegistry)(nil)
)

// NewServiceInstanceRegistry returns a ready-to-use registry.
func NewServiceInstanceRegistry() *ServiceInstanceRegistry {
	return &ServiceInstanceRegistry{
		store: make(map[string]resources.ServiceInstance),
	}
}

// deepCopyServiceInstance returns a fully independent copy of si, duplicating
// Metadata Labels/Annotations and Spec.Parameters maps so callers cannot
// mutate the registry's internal state.
func deepCopyServiceInstance(si resources.ServiceInstance) resources.ServiceInstance {
	cp := si
	if si.Metadata.Labels != nil {
		cp.Metadata.Labels = make(map[string]string, len(si.Metadata.Labels))
		for k, v := range si.Metadata.Labels {
			cp.Metadata.Labels[k] = v
		}
	}
	if si.Metadata.Annotations != nil {
		cp.Metadata.Annotations = make(map[string]string, len(si.Metadata.Annotations))
		for k, v := range si.Metadata.Annotations {
			cp.Metadata.Annotations[k] = v
		}
	}
	if si.Spec.Parameters != nil {
		cp.Spec.Parameters = make(map[string]string, len(si.Spec.Parameters))
		for k, v := range si.Spec.Parameters {
			cp.Spec.Parameters[k] = v
		}
	}
	return cp
}

// CreateServiceInstance stores a deep copy of si keyed by metadata.name. It
// returns ErrAlreadyExists if the name is already present.
func (r *ServiceInstanceRegistry) CreateServiceInstance(
	ctx context.Context, si resources.ServiceInstance,
) (resources.ServiceInstance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[si.Metadata.Name]; ok {
		return resources.ServiceInstance{}, ErrAlreadyExists
	}
	stored := deepCopyServiceInstance(si)
	r.store[si.Metadata.Name] = stored
	return deepCopyServiceInstance(stored), nil
}

// GetServiceInstance returns a deep copy of the stored ServiceInstance
// identified by name, or ErrNotFound if absent.
func (r *ServiceInstanceRegistry) GetServiceInstance(
	ctx context.Context, name string,
) (resources.ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	si, ok := r.store[name]
	if !ok {
		return resources.ServiceInstance{}, ErrNotFound
	}
	return deepCopyServiceInstance(si), nil
}

// ListServiceInstances returns a new slice of deep copies sorted by
// metadata.name ascending. When tenantRef is non-empty, only entries with
// matching Spec.TenantRef are included. When projectRef is non-empty, only
// entries with matching Spec.ProjectRef are included. When both are
// non-empty, both filters apply (AND). Empty filter strings are not applied.
// Returns a non-nil empty slice when no ServiceInstances match.
func (r *ServiceInstanceRegistry) ListServiceInstances(
	ctx context.Context, tenantRef, projectRef string,
) ([]resources.ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]resources.ServiceInstance, 0, len(r.store))
	for _, si := range r.store {
		if tenantRef != "" && si.Spec.TenantRef != tenantRef {
			continue
		}
		if projectRef != "" && si.Spec.ProjectRef != projectRef {
			continue
		}
		items = append(items, deepCopyServiceInstance(si))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Metadata.Name < items[j].Metadata.Name
	})
	return items, nil
}

// UpdateServiceInstance looks up the existing stored entry by name, replaces
// only the mutable fields, and returns a deep copy of the updated
// ServiceInstance. Immutable and system-owned fields (apiVersion, kind,
// status, metadata.name, and all governance/catalog spec refs) are preserved
// from the stored entry. Returns ErrNotFound if the name is absent.
func (r *ServiceInstanceRegistry) UpdateServiceInstance(
	ctx context.Context, name string, si resources.ServiceInstance,
) (resources.ServiceInstance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.store[name]
	if !ok {
		return resources.ServiceInstance{}, ErrNotFound
	}
	merged := existing
	merged.Metadata.DisplayName = si.Metadata.DisplayName
	merged.Metadata.Labels = si.Metadata.Labels
	merged.Metadata.Annotations = si.Metadata.Annotations
	merged.Spec.Parameters = si.Spec.Parameters
	stored := deepCopyServiceInstance(merged)
	r.store[name] = stored
	return deepCopyServiceInstance(stored), nil
}

// DeleteServiceInstance removes the entry identified by name, or returns
// ErrNotFound if absent.
func (r *ServiceInstanceRegistry) DeleteServiceInstance(
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

// CountByServicePlan returns the number of stored ServiceInstances whose
// Spec.ServiceClassRef and Spec.ServicePlanRef both match the given refs.
func (r *ServiceInstanceRegistry) CountByServicePlan(
	ctx context.Context, serviceClassRef, servicePlanRef string,
) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, si := range r.store {
		if si.Spec.ServiceClassRef == serviceClassRef &&
			si.Spec.ServicePlanRef == servicePlanRef {
			count++
		}
	}
	return count, nil
}

// CountByProject returns the number of stored ServiceInstances whose four
// governance refs match exactly (including empty organizationUnitRef).
func (r *ServiceInstanceRegistry) CountByProject(
	ctx context.Context, organizationRef, organizationUnitRef, tenantRef, projectRef string,
) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, si := range r.store {
		if si.Spec.OrganizationRef == organizationRef &&
			si.Spec.OrganizationUnitRef == organizationUnitRef &&
			si.Spec.TenantRef == tenantRef &&
			si.Spec.ProjectRef == projectRef {
			count++
		}
	}
	return count, nil
}
