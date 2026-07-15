package registry

import (
	"context"
	"sort"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// OrganizationRegistry is the Phase 1 in-memory implementation of
// OrganizationRegistryIface. All public methods are safe for concurrent
// use. The registry holds no package-level global state.
type OrganizationRegistry struct {
	mu    sync.RWMutex
	store map[string]resources.Organization
}

// NewOrganizationRegistry returns a ready-to-use registry.
func NewOrganizationRegistry() *OrganizationRegistry {
	return &OrganizationRegistry{
		store: make(map[string]resources.Organization),
	}
}

// deepCopyOrganization returns a fully independent copy of org,
// duplicating the Labels, Annotations maps and SovereignLocations slice
// so that callers cannot mutate the registry's internal state.
func deepCopyOrganization(org resources.Organization) resources.Organization {
	cp := org
	if org.Metadata.Labels != nil {
		cp.Metadata.Labels = make(map[string]string, len(org.Metadata.Labels))
		for k, v := range org.Metadata.Labels {
			cp.Metadata.Labels[k] = v
		}
	}
	if org.Metadata.Annotations != nil {
		cp.Metadata.Annotations = make(map[string]string, len(org.Metadata.Annotations))
		for k, v := range org.Metadata.Annotations {
			cp.Metadata.Annotations[k] = v
		}
	}
	if org.Spec.SovereignLocations != nil {
		cp.Spec.SovereignLocations = make([]string, len(org.Spec.SovereignLocations))
		copy(cp.Spec.SovereignLocations, org.Spec.SovereignLocations)
	}
	return cp
}

// CreateOrganization stores a deep copy of org keyed by org.Metadata.Name.
func (r *OrganizationRegistry) CreateOrganization(
	ctx context.Context, org resources.Organization,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[org.Metadata.Name]; ok {
		return ErrAlreadyExists
	}
	r.store[org.Metadata.Name] = deepCopyOrganization(org)
	return nil
}

// GetOrganization returns a deep copy of the stored Organization.
func (r *OrganizationRegistry) GetOrganization(
	ctx context.Context, name string,
) (resources.Organization, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	org, ok := r.store[name]
	if !ok {
		return resources.Organization{}, ErrNotFound
	}
	return deepCopyOrganization(org), nil
}

// ListOrganizations returns a new slice of deep copies sorted
// ascending by metadata.name.
func (r *OrganizationRegistry) ListOrganizations(
	ctx context.Context,
) ([]resources.Organization, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]resources.Organization, 0, len(r.store))
	for _, org := range r.store {
		items = append(items, deepCopyOrganization(org))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Metadata.Name < items[j].Metadata.Name
	})
	return items, nil
}

// UpdateOrganization replaces the mutable fields of the stored entry and
// returns a deep copy of the updated Organization.
func (r *OrganizationRegistry) UpdateOrganization(
	ctx context.Context, name string, updated resources.Organization,
) (resources.Organization, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.store[name]
	if !ok {
		return resources.Organization{}, ErrNotFound
	}
	updated.Metadata.Name = existing.Metadata.Name
	updated.Status = existing.Status
	updated.APIVersion = resources.OrgAPIVersion
	updated.Kind = resources.OrgKind
	stored := deepCopyOrganization(updated)
	r.store[name] = stored
	return deepCopyOrganization(stored), nil
}

// DeleteOrganization removes the entry.
func (r *OrganizationRegistry) DeleteOrganization(
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
