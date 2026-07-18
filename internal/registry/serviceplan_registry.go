package registry

import (
	"context"
	"sort"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// ServicePlanRegistryIface is the storage contract for ServicePlan resources.
// The registry is storage-only: it does not depend on ServiceClassRegistry or
// any other registry and performs no parent-existence checks. Those checks
// belong to the API/service layer via ServiceClassLookup.
type ServicePlanRegistryIface interface {
	CreateServicePlan(ctx context.Context, sp resources.ServicePlan) (resources.ServicePlan, error)
	GetServicePlan(ctx context.Context, serviceClassName, name string) (resources.ServicePlan, error)
	ListServicePlans(ctx context.Context) ([]resources.ServicePlan, error)
	UpdateServicePlan(ctx context.Context, sp resources.ServicePlan) (resources.ServicePlan, error)
	DeleteServicePlan(ctx context.Context, serviceClassName, name string) error
	CountByServiceClass(ctx context.Context, serviceClassName string) (int, error)
}

// ServicePlanLookup is a narrow interface for verifying ServicePlan existence
// and its association with a ServiceClass. The existing *ServicePlanRegistry
// already satisfies it via GetServicePlan. It is intended for injection into
// the ServiceInstanceHandler and is NOT used inside ServicePlanRegistry.
type ServicePlanLookup interface {
	GetServicePlan(ctx context.Context, serviceClassName, name string) (resources.ServicePlan, error)
}

// ServicePlanRegistry is the Phase 1 in-memory implementation of
// ServicePlanRegistryIface. All public methods are safe for concurrent use.
// The registry holds no package-level global state and uses the composite key
// "serviceClassName/name" as the map key.
type ServicePlanRegistry struct {
	mu    sync.RWMutex
	store map[string]resources.ServicePlan
}

// NewServicePlanRegistry returns a ready-to-use registry.
func NewServicePlanRegistry() *ServicePlanRegistry {
	return &ServicePlanRegistry{
		store: make(map[string]resources.ServicePlan),
	}
}

// servicePlanCompositeKey builds the map key from serviceClassName and name.
func servicePlanCompositeKey(serviceClassName, name string) string {
	return serviceClassName + "/" + name
}

// deepCopyServicePlan returns a fully independent copy of sp, duplicating the
// Parameters map, Tags slice, and Metadata Labels/Annotations maps so that
// callers cannot mutate the registry's internal state.
func deepCopyServicePlan(sp resources.ServicePlan) resources.ServicePlan {
	cp := sp
	if sp.Metadata.Labels != nil {
		cp.Metadata.Labels = make(map[string]string, len(sp.Metadata.Labels))
		for k, v := range sp.Metadata.Labels {
			cp.Metadata.Labels[k] = v
		}
	}
	if sp.Metadata.Annotations != nil {
		cp.Metadata.Annotations = make(map[string]string, len(sp.Metadata.Annotations))
		for k, v := range sp.Metadata.Annotations {
			cp.Metadata.Annotations[k] = v
		}
	}
	if sp.Spec.Parameters != nil {
		cp.Spec.Parameters = make(map[string]string, len(sp.Spec.Parameters))
		for k, v := range sp.Spec.Parameters {
			cp.Spec.Parameters[k] = v
		}
	}
	if sp.Spec.Tags != nil {
		cp.Spec.Tags = make([]string, len(sp.Spec.Tags))
		copy(cp.Spec.Tags, sp.Spec.Tags)
	}
	return cp
}

// CreateServicePlan stores a deep copy of sp keyed by the composite key
// "serviceClassName/name". It returns ErrAlreadyExists if the composite key is
// already present. The same metadata.name under different serviceClassName
// values yields distinct keys and does not conflict.
func (r *ServicePlanRegistry) CreateServicePlan(
	ctx context.Context, sp resources.ServicePlan,
) (resources.ServicePlan, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := servicePlanCompositeKey(sp.Spec.ServiceClassName, sp.Metadata.Name)
	if _, ok := r.store[key]; ok {
		return resources.ServicePlan{}, ErrAlreadyExists
	}
	stored := deepCopyServicePlan(sp)
	r.store[key] = stored
	return deepCopyServicePlan(stored), nil
}

// GetServicePlan returns a deep copy of the stored ServicePlan identified by
// the composite key, or ErrNotFound if absent.
func (r *ServicePlanRegistry) GetServicePlan(
	ctx context.Context, serviceClassName, name string,
) (resources.ServicePlan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sp, ok := r.store[servicePlanCompositeKey(serviceClassName, name)]
	if !ok {
		return resources.ServicePlan{}, ErrNotFound
	}
	return deepCopyServicePlan(sp), nil
}

// ListServicePlans returns a new slice of deep copies sorted by
// spec.serviceClassName ascending, then metadata.name ascending. It returns a
// non-nil empty slice when no ServicePlans are stored.
func (r *ServicePlanRegistry) ListServicePlans(
	ctx context.Context,
) ([]resources.ServicePlan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]resources.ServicePlan, 0, len(r.store))
	for _, sp := range r.store {
		items = append(items, deepCopyServicePlan(sp))
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Spec.ServiceClassName != items[j].Spec.ServiceClassName {
			return items[i].Spec.ServiceClassName < items[j].Spec.ServiceClassName
		}
		return items[i].Metadata.Name < items[j].Metadata.Name
	})
	return items, nil
}

// UpdateServicePlan derives the composite key from the submitted ServicePlan,
// looks up the existing stored entry, replaces only the mutable fields, and
// returns a deep copy of the updated ServicePlan. Immutable and system-owned
// fields (apiVersion, kind, status, metadata.name, spec.serviceClassName) are
// preserved from the stored entry so a plan is never moved between classes.
// Returns ErrNotFound if the composite key is absent.
func (r *ServicePlanRegistry) UpdateServicePlan(
	ctx context.Context, sp resources.ServicePlan,
) (resources.ServicePlan, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := servicePlanCompositeKey(sp.Spec.ServiceClassName, sp.Metadata.Name)
	existing, ok := r.store[key]
	if !ok {
		return resources.ServicePlan{}, ErrNotFound
	}
	merged := existing
	merged.Metadata.Labels = sp.Metadata.Labels
	merged.Metadata.Annotations = sp.Metadata.Annotations
	merged.Spec.DisplayName = sp.Spec.DisplayName
	merged.Spec.Description = sp.Spec.Description
	merged.Spec.Tier = sp.Spec.Tier
	merged.Spec.Lifecycle = sp.Spec.Lifecycle
	merged.Spec.Parameters = sp.Spec.Parameters
	merged.Spec.Tags = sp.Spec.Tags
	stored := deepCopyServicePlan(merged)
	r.store[key] = stored
	return deepCopyServicePlan(stored), nil
}

// DeleteServicePlan removes the entry identified by the composite key, or
// returns ErrNotFound if absent.
func (r *ServicePlanRegistry) DeleteServicePlan(
	ctx context.Context, serviceClassName, name string,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := servicePlanCompositeKey(serviceClassName, name)
	if _, ok := r.store[key]; !ok {
		return ErrNotFound
	}
	delete(r.store, key)
	return nil
}

// CountByServiceClass returns the number of stored ServicePlans whose
// spec.serviceClassName matches the given parent ServiceClass.
func (r *ServicePlanRegistry) CountByServiceClass(
	ctx context.Context, serviceClassName string,
) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, sp := range r.store {
		if sp.Spec.ServiceClassName == serviceClassName {
			count++
		}
	}
	return count, nil
}
