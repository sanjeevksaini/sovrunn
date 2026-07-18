package registry

import (
	"context"
	"sort"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// ProjectRegistryIface is the storage contract for Project resources. The
// registry is storage-only: it does not depend on TenantRegistry,
// OrganizationUnitRegistry, or OrganizationRegistry and performs no
// parent-existence checks. Those checks belong to the API/service layer.
type ProjectRegistryIface interface {
	CreateProject(ctx context.Context, p resources.Project) (resources.Project, error)
	GetProject(ctx context.Context, orgName, ouName, tenantName, name string) (resources.Project, error)
	ListProjects(ctx context.Context) ([]resources.Project, error)
	UpdateProject(ctx context.Context, p resources.Project) (resources.Project, error)
	DeleteProject(ctx context.Context, orgName, ouName, tenantName, name string) error
	CountByTenant(ctx context.Context, orgName, ouName, tenantName string) (int, error)
}

// TenantLookup is a narrow interface for verifying parent Tenant existence.
// The existing *TenantRegistry already satisfies it via GetTenant. It is
// intended for injection into the ProjectHandler in a later task and is NOT
// used inside ProjectRegistry.
type TenantLookup interface {
	GetTenant(ctx context.Context, orgName, ouName, name string) (resources.Tenant, error)
}

// ProjectLookup is a narrow interface for verifying parent Project existence.
// The existing *ProjectRegistry already satisfies it via GetProject. It is
// intended for injection into the ServiceInstanceHandler and is NOT used
// inside ProjectRegistry.
type ProjectLookup interface {
	GetProject(ctx context.Context, orgName, ouName, tenantName, name string) (resources.Project, error)
}

// ProjectRegistry is the Phase 1 in-memory implementation of
// ProjectRegistryIface. All public methods are safe for concurrent use. The
// registry holds no package-level global state and uses the composite key
// "organizationName/organizationUnitName/tenantName/name" as the map key.
type ProjectRegistry struct {
	mu    sync.RWMutex
	store map[string]resources.Project
}

// NewProjectRegistry returns a ready-to-use registry.
func NewProjectRegistry() *ProjectRegistry {
	return &ProjectRegistry{
		store: make(map[string]resources.Project),
	}
}

// projectCompositeKey builds the map key from orgName, ouName, tenantName, and
// name.
func projectCompositeKey(orgName, ouName, tenantName, name string) string {
	return orgName + "/" + ouName + "/" + tenantName + "/" + name
}

// deepCopyProject returns a fully independent copy of p, duplicating the Labels
// and Annotations maps so that callers cannot mutate the registry's internal
// state.
func deepCopyProject(p resources.Project) resources.Project {
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
	return cp
}

// CreateProject stores a deep copy of p keyed by the composite key
// "organizationName/organizationUnitName/tenantName/name". It returns
// ErrAlreadyExists if the composite key is already present.
func (r *ProjectRegistry) CreateProject(
	ctx context.Context, p resources.Project,
) (resources.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := projectCompositeKey(p.Spec.OrganizationName, p.Spec.OrganizationUnitName, p.Spec.TenantName, p.Metadata.Name)
	if _, ok := r.store[key]; ok {
		return resources.Project{}, ErrAlreadyExists
	}
	stored := deepCopyProject(p)
	r.store[key] = stored
	return deepCopyProject(stored), nil
}

// GetProject returns a deep copy of the stored Project identified by the
// composite key, or ErrNotFound if absent.
func (r *ProjectRegistry) GetProject(
	ctx context.Context, orgName, ouName, tenantName, name string,
) (resources.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[projectCompositeKey(orgName, ouName, tenantName, name)]
	if !ok {
		return resources.Project{}, ErrNotFound
	}
	return deepCopyProject(p), nil
}

// ListProjects returns a new slice of deep copies sorted by
// spec.organizationName ascending, then spec.organizationUnitName ascending,
// then spec.tenantName ascending, then metadata.name ascending. It returns a
// non-nil empty slice when no Projects are stored.
func (r *ProjectRegistry) ListProjects(
	ctx context.Context,
) ([]resources.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]resources.Project, 0, len(r.store))
	for _, p := range r.store {
		items = append(items, deepCopyProject(p))
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Spec.OrganizationName != items[j].Spec.OrganizationName {
			return items[i].Spec.OrganizationName < items[j].Spec.OrganizationName
		}
		if items[i].Spec.OrganizationUnitName != items[j].Spec.OrganizationUnitName {
			return items[i].Spec.OrganizationUnitName < items[j].Spec.OrganizationUnitName
		}
		if items[i].Spec.TenantName != items[j].Spec.TenantName {
			return items[i].Spec.TenantName < items[j].Spec.TenantName
		}
		return items[i].Metadata.Name < items[j].Metadata.Name
	})
	return items, nil
}

// UpdateProject derives the composite key from the submitted Project, looks up
// the existing stored entry, replaces only the mutable fields, and returns a
// deep copy of the updated Project. Immutable and system-owned fields
// (apiVersion, kind, status, metadata.name, spec.organizationName,
// spec.organizationUnitName, spec.tenantName) are preserved from the stored
// entry. Returns ErrNotFound if the composite key is absent.
func (r *ProjectRegistry) UpdateProject(
	ctx context.Context, p resources.Project,
) (resources.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := projectCompositeKey(p.Spec.OrganizationName, p.Spec.OrganizationUnitName, p.Spec.TenantName, p.Metadata.Name)
	existing, ok := r.store[key]
	if !ok {
		return resources.Project{}, ErrNotFound
	}
	merged := existing
	merged.Metadata.DisplayName = p.Metadata.DisplayName
	merged.Metadata.Labels = p.Metadata.Labels
	merged.Metadata.Annotations = p.Metadata.Annotations
	merged.Spec.Description = p.Spec.Description
	stored := deepCopyProject(merged)
	r.store[key] = stored
	return deepCopyProject(stored), nil
}

// DeleteProject removes the entry identified by the composite key, or returns
// ErrNotFound if absent.
func (r *ProjectRegistry) DeleteProject(
	ctx context.Context, orgName, ouName, tenantName, name string,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := projectCompositeKey(orgName, ouName, tenantName, name)
	if _, ok := r.store[key]; !ok {
		return ErrNotFound
	}
	delete(r.store, key)
	return nil
}

// CountByTenant returns the number of stored Projects whose
// spec.organizationName, spec.organizationUnitName, and spec.tenantName match
// the given parent Tenant.
func (r *ProjectRegistry) CountByTenant(
	ctx context.Context, orgName, ouName, tenantName string,
) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, p := range r.store {
		if p.Spec.OrganizationName == orgName &&
			p.Spec.OrganizationUnitName == ouName &&
			p.Spec.TenantName == tenantName {
			count++
		}
	}
	return count, nil
}
