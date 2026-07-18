package registry

import (
	"context"
	"sort"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// CapabilityRegistryIface is the storage contract for Capability resources.
// The registry is storage-only: it does not depend on other registries and
// performs no cross-resource checks. Parent/child checks belong to the
// API/service layer. Capability is immutable after create (no Update).
type CapabilityRegistryIface interface {
	CreateCapability(ctx context.Context, c resources.Capability) (resources.Capability, error)
	GetCapability(ctx context.Context, name string) (resources.Capability, error)
	ListCapabilities(ctx context.Context, pluginRef, serviceClassRef string) ([]resources.Capability, error)
	DeleteCapability(ctx context.Context, name string) error
	CountByPlugin(ctx context.Context, pluginName string) (int, error)
}

// CapabilityRegistry is the Phase 1 in-memory implementation of
// CapabilityRegistryIface. All public methods are safe for concurrent use.
// The registry holds no package-level global state and is keyed by
// metadata.name.
type CapabilityRegistry struct {
	mu    sync.RWMutex
	store map[string]resources.Capability
}

// NewCapabilityRegistry returns a ready-to-use registry.
func NewCapabilityRegistry() *CapabilityRegistry {
	return &CapabilityRegistry{
		store: make(map[string]resources.Capability),
	}
}

// deepCopyCapability returns a fully independent copy of c, duplicating the
// Metadata Labels/Annotations maps so that callers cannot mutate the
// registry's internal state. CapabilitySpec contains only scalar fields, so
// no slice/map copy is needed for spec.
func deepCopyCapability(c resources.Capability) resources.Capability {
	cp := c
	if c.Metadata.Labels != nil {
		cp.Metadata.Labels = make(map[string]string, len(c.Metadata.Labels))
		for k, v := range c.Metadata.Labels {
			cp.Metadata.Labels[k] = v
		}
	}
	if c.Metadata.Annotations != nil {
		cp.Metadata.Annotations = make(map[string]string, len(c.Metadata.Annotations))
		for k, v := range c.Metadata.Annotations {
			cp.Metadata.Annotations[k] = v
		}
	}
	return cp
}

// CreateCapability stores a deep copy of c keyed by metadata.name. It returns
// ErrAlreadyExists if the name is already present.
func (r *CapabilityRegistry) CreateCapability(
	ctx context.Context, c resources.Capability,
) (resources.Capability, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[c.Metadata.Name]; ok {
		return resources.Capability{}, ErrAlreadyExists
	}
	stored := deepCopyCapability(c)
	r.store[c.Metadata.Name] = stored
	return deepCopyCapability(stored), nil
}

// GetCapability returns a deep copy of the stored Capability identified by
// name, or ErrNotFound if absent.
func (r *CapabilityRegistry) GetCapability(
	ctx context.Context, name string,
) (resources.Capability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.store[name]
	if !ok {
		return resources.Capability{}, ErrNotFound
	}
	return deepCopyCapability(c), nil
}

// ListCapabilities returns a new slice of deep copies sorted by metadata.name
// ascending. When pluginRef is non-empty, only entries with matching
// Spec.PluginRef are included. When serviceClassRef is non-empty, only
// entries with matching Spec.ServiceClassRef are included. When both are
// non-empty, both filters apply (AND). Empty filter strings are not applied.
// Returns a non-nil empty slice when no Capabilities match.
func (r *CapabilityRegistry) ListCapabilities(
	ctx context.Context, pluginRef, serviceClassRef string,
) ([]resources.Capability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]resources.Capability, 0, len(r.store))
	for _, c := range r.store {
		if pluginRef != "" && c.Spec.PluginRef != pluginRef {
			continue
		}
		if serviceClassRef != "" && c.Spec.ServiceClassRef != serviceClassRef {
			continue
		}
		items = append(items, deepCopyCapability(c))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Metadata.Name < items[j].Metadata.Name
	})
	return items, nil
}

// DeleteCapability removes the entry identified by name, or returns
// ErrNotFound if absent.
func (r *CapabilityRegistry) DeleteCapability(
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

// CountByPlugin returns the number of stored Capabilities whose
// Spec.PluginRef equals pluginName. Used by the Plugin delete blocker.
func (r *CapabilityRegistry) CountByPlugin(
	ctx context.Context, pluginName string,
) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, c := range r.store {
		if c.Spec.PluginRef == pluginName {
			count++
		}
	}
	return count, nil
}
