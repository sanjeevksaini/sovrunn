package registry

import (
	"context"
	"sort"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// ServiceClassRegistryIface is the storage contract for ServiceClass resources.
// The registry is storage-only: it does not depend on other registries and
// performs no cross-resource checks. Parent/child checks belong to the
// API/service layer.
type ServiceClassRegistryIface interface {
	CreateServiceClass(ctx context.Context, sc resources.ServiceClass) (resources.ServiceClass, error)
	GetServiceClass(ctx context.Context, name string) (resources.ServiceClass, error)
	ListServiceClasses(ctx context.Context) ([]resources.ServiceClass, error)
	UpdateServiceClass(ctx context.Context, sc resources.ServiceClass) (resources.ServiceClass, error)
	DeleteServiceClass(ctx context.Context, name string) error
}

// ServiceClassLookup is a narrow interface for verifying parent ServiceClass
// existence. The existing *ServiceClassRegistry already satisfies it via
// GetServiceClass. It is intended for injection into the ServicePlanHandler and
// is NOT used inside ServiceClassRegistry or ServicePlanRegistry.
type ServiceClassLookup interface {
	GetServiceClass(ctx context.Context, name string) (resources.ServiceClass, error)
}

// ServiceClassRegistry is the Phase 1 in-memory implementation of
// ServiceClassRegistryIface. All public methods are safe for concurrent use.
// The registry holds no package-level global state and is keyed by
// metadata.name.
type ServiceClassRegistry struct {
	mu    sync.RWMutex
	store map[string]resources.ServiceClass
}

// NewServiceClassRegistry returns a ready-to-use registry.
func NewServiceClassRegistry() *ServiceClassRegistry {
	return &ServiceClassRegistry{
		store: make(map[string]resources.ServiceClass),
	}
}

// deepCopyServiceClass returns a fully independent copy of sc, duplicating the
// Tags slice and Metadata Labels/Annotations maps so that callers cannot mutate
// the registry's internal state.
func deepCopyServiceClass(sc resources.ServiceClass) resources.ServiceClass {
	cp := sc
	if sc.Metadata.Labels != nil {
		cp.Metadata.Labels = make(map[string]string, len(sc.Metadata.Labels))
		for k, v := range sc.Metadata.Labels {
			cp.Metadata.Labels[k] = v
		}
	}
	if sc.Metadata.Annotations != nil {
		cp.Metadata.Annotations = make(map[string]string, len(sc.Metadata.Annotations))
		for k, v := range sc.Metadata.Annotations {
			cp.Metadata.Annotations[k] = v
		}
	}
	if sc.Spec.Tags != nil {
		cp.Spec.Tags = make([]string, len(sc.Spec.Tags))
		copy(cp.Spec.Tags, sc.Spec.Tags)
	}
	return cp
}

// CreateServiceClass stores a deep copy of sc keyed by metadata.name. It
// returns ErrAlreadyExists if the name is already present.
func (r *ServiceClassRegistry) CreateServiceClass(
	ctx context.Context, sc resources.ServiceClass,
) (resources.ServiceClass, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[sc.Metadata.Name]; ok {
		return resources.ServiceClass{}, ErrAlreadyExists
	}
	stored := deepCopyServiceClass(sc)
	r.store[sc.Metadata.Name] = stored
	return deepCopyServiceClass(stored), nil
}

// GetServiceClass returns a deep copy of the stored ServiceClass identified by
// name, or ErrNotFound if absent.
func (r *ServiceClassRegistry) GetServiceClass(
	ctx context.Context, name string,
) (resources.ServiceClass, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sc, ok := r.store[name]
	if !ok {
		return resources.ServiceClass{}, ErrNotFound
	}
	return deepCopyServiceClass(sc), nil
}

// ListServiceClasses returns a new slice of deep copies sorted by
// metadata.name ascending. It returns a non-nil empty slice when no
// ServiceClasses are stored.
func (r *ServiceClassRegistry) ListServiceClasses(
	ctx context.Context,
) ([]resources.ServiceClass, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]resources.ServiceClass, 0, len(r.store))
	for _, sc := range r.store {
		items = append(items, deepCopyServiceClass(sc))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Metadata.Name < items[j].Metadata.Name
	})
	return items, nil
}

// UpdateServiceClass derives the key from sc.Metadata.Name, looks up the
// existing stored entry, replaces only the mutable fields, and returns a deep
// copy of the updated ServiceClass. Immutable and system-owned fields
// (apiVersion, kind, status, metadata.name) are preserved from the stored
// entry. Returns ErrNotFound if the name is absent.
func (r *ServiceClassRegistry) UpdateServiceClass(
	ctx context.Context, sc resources.ServiceClass,
) (resources.ServiceClass, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.store[sc.Metadata.Name]
	if !ok {
		return resources.ServiceClass{}, ErrNotFound
	}
	merged := existing
	merged.Metadata.Labels = sc.Metadata.Labels
	merged.Metadata.Annotations = sc.Metadata.Annotations
	merged.Spec.DisplayName = sc.Spec.DisplayName
	merged.Spec.Description = sc.Spec.Description
	merged.Spec.Category = sc.Spec.Category
	merged.Spec.Provider = sc.Spec.Provider
	merged.Spec.Lifecycle = sc.Spec.Lifecycle
	merged.Spec.DefaultPlanName = sc.Spec.DefaultPlanName
	merged.Spec.Tags = sc.Spec.Tags
	stored := deepCopyServiceClass(merged)
	r.store[sc.Metadata.Name] = stored
	return deepCopyServiceClass(stored), nil
}

// DeleteServiceClass removes the entry identified by name, or returns
// ErrNotFound if absent.
func (r *ServiceClassRegistry) DeleteServiceClass(
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
