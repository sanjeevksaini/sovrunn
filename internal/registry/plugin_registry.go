package registry

import (
	"context"
	"sort"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// PluginRegistryIface is the storage contract for Plugin resources.
// The registry is storage-only: it does not depend on other registries and
// performs no cross-resource checks. Parent/child checks belong to the
// API/service layer.
type PluginRegistryIface interface {
	CreatePlugin(ctx context.Context, p resources.Plugin) (resources.Plugin, error)
	GetPlugin(ctx context.Context, name string) (resources.Plugin, error)
	ListPlugins(ctx context.Context) ([]resources.Plugin, error)
	UpdatePlugin(ctx context.Context, p resources.Plugin) (resources.Plugin, error)
	DeletePlugin(ctx context.Context, name string) error
}

// PluginLookup is a narrow interface for verifying Plugin existence.
// The concrete *PluginRegistry satisfies it via GetPlugin. It is intended for
// injection into the CapabilityHandler and is NOT used inside PluginRegistry
// or CapabilityRegistry.
type PluginLookup interface {
	GetPlugin(ctx context.Context, name string) (resources.Plugin, error)
}

// PluginRegistry is the Phase 1 in-memory implementation of
// PluginRegistryIface. All public methods are safe for concurrent use.
// The registry holds no package-level global state and is keyed by
// metadata.name.
type PluginRegistry struct {
	mu    sync.RWMutex
	store map[string]resources.Plugin
}

// NewPluginRegistry returns a ready-to-use registry.
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		store: make(map[string]resources.Plugin),
	}
}

// deepCopyPlugin returns a fully independent copy of p, duplicating the
// ServiceClassRefs and Tags slices and Metadata Labels/Annotations maps so
// that callers cannot mutate the registry's internal state.
func deepCopyPlugin(p resources.Plugin) resources.Plugin {
	cp := p
	if p.Metadata.Labels != nil {
		cp.Metadata.Labels = make(map[string]string, len(p.Metadata.Labels))
		for k, v := range p.Metadata.Labels {
			cp.Metadata.Labels[k] = v
		}
	}
	if p.Metadata.Annotations != nil {
		cp.Metadata.Annotations = make(map[string]string, len(p.Metadata.Annotations))
		for k, v := range p.Metadata.Annotations {
			cp.Metadata.Annotations[k] = v
		}
	}
	if p.Spec.ServiceClassRefs != nil {
		cp.Spec.ServiceClassRefs = make([]string, len(p.Spec.ServiceClassRefs))
		copy(cp.Spec.ServiceClassRefs, p.Spec.ServiceClassRefs)
	}
	if p.Spec.Tags != nil {
		cp.Spec.Tags = make([]string, len(p.Spec.Tags))
		copy(cp.Spec.Tags, p.Spec.Tags)
	}
	return cp
}

// CreatePlugin stores a deep copy of p keyed by metadata.name. It returns
// ErrAlreadyExists if the name is already present.
func (r *PluginRegistry) CreatePlugin(
	ctx context.Context, p resources.Plugin,
) (resources.Plugin, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[p.Metadata.Name]; ok {
		return resources.Plugin{}, ErrAlreadyExists
	}
	stored := deepCopyPlugin(p)
	r.store[p.Metadata.Name] = stored
	return deepCopyPlugin(stored), nil
}

// GetPlugin returns a deep copy of the stored Plugin identified by name, or
// ErrNotFound if absent.
func (r *PluginRegistry) GetPlugin(
	ctx context.Context, name string,
) (resources.Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[name]
	if !ok {
		return resources.Plugin{}, ErrNotFound
	}
	return deepCopyPlugin(p), nil
}

// ListPlugins returns a new slice of deep copies sorted by metadata.name
// ascending. It returns a non-nil empty slice when no Plugins are stored.
func (r *PluginRegistry) ListPlugins(
	ctx context.Context,
) ([]resources.Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]resources.Plugin, 0, len(r.store))
	for _, p := range r.store {
		items = append(items, deepCopyPlugin(p))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Metadata.Name < items[j].Metadata.Name
	})
	return items, nil
}

// UpdatePlugin derives the key from p.Metadata.Name, looks up the existing
// stored entry, replaces only the mutable fields, and returns a deep copy of
// the updated Plugin. Immutable and system-owned fields (apiVersion, kind,
// status, metadata.name) are preserved from the stored entry. Returns
// ErrNotFound if the name is absent.
func (r *PluginRegistry) UpdatePlugin(
	ctx context.Context, p resources.Plugin,
) (resources.Plugin, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.store[p.Metadata.Name]
	if !ok {
		return resources.Plugin{}, ErrNotFound
	}
	merged := existing
	merged.Metadata.Labels = p.Metadata.Labels
	merged.Metadata.Annotations = p.Metadata.Annotations
	merged.Spec.PluginType = p.Spec.PluginType
	merged.Spec.Version = p.Spec.Version
	merged.Spec.ServiceClassRefs = p.Spec.ServiceClassRefs
	merged.Spec.DeploymentMode = p.Spec.DeploymentMode
	merged.Spec.Description = p.Spec.Description
	merged.Spec.Tags = p.Spec.Tags
	stored := deepCopyPlugin(merged)
	r.store[p.Metadata.Name] = stored
	return deepCopyPlugin(stored), nil
}

// DeletePlugin removes the entry identified by name, or returns ErrNotFound
// if absent.
func (r *PluginRegistry) DeletePlugin(
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
