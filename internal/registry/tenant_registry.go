package registry

import (
	"context"
	"sort"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// TenantRegistryIface is the storage contract for Tenant resources. The
// registry is storage-only: it does not depend on OrganizationUnitRegistry or
// OrganizationRegistry and performs no parent-existence checks. Those checks
// belong to the API/service layer.
type TenantRegistryIface interface {
	CreateTenant(ctx context.Context, t resources.Tenant) (resources.Tenant, error)
	GetTenant(ctx context.Context, orgName, ouName, name string) (resources.Tenant, error)
	ListTenants(ctx context.Context) ([]resources.Tenant, error)
	UpdateTenant(ctx context.Context, t resources.Tenant) (resources.Tenant, error)
	DeleteTenant(ctx context.Context, orgName, ouName, name string) error
	CountByOrganizationUnit(ctx context.Context, orgName, ouName string) (int, error)
}

// OrganizationUnitLookup is a narrow interface for verifying parent
// OrganizationUnit existence. The existing *OrganizationUnitRegistry already
// satisfies it via GetOrganizationUnit. It is intended for injection into the
// TenantHandler in a later task and is NOT used inside TenantRegistry.
type OrganizationUnitLookup interface {
	GetOrganizationUnit(ctx context.Context, orgName, name string) (resources.OrganizationUnit, error)
}

// TenantRegistry is the Phase 1 in-memory implementation of
// TenantRegistryIface. All public methods are safe for concurrent use. The
// registry holds no package-level global state and uses the composite key
// "organizationName/organizationUnitName/name" as the map key.
type TenantRegistry struct {
	mu    sync.RWMutex
	store map[string]resources.Tenant
}

// NewTenantRegistry returns a ready-to-use registry.
func NewTenantRegistry() *TenantRegistry {
	return &TenantRegistry{
		store: make(map[string]resources.Tenant),
	}
}

// tenantCompositeKey builds the map key from orgName, ouName, and name.
func tenantCompositeKey(orgName, ouName, name string) string {
	return orgName + "/" + ouName + "/" + name
}

// deepCopyTenant returns a fully independent copy of t, duplicating the Labels
// and Annotations maps so that callers cannot mutate the registry's internal
// state.
func deepCopyTenant(t resources.Tenant) resources.Tenant {
	cp := t
	if t.Metadata.Labels != nil {
		cp.Metadata.Labels = make(map[string]string, len(t.Metadata.Labels))
		for k, v := range t.Metadata.Labels {
			cp.Metadata.Labels[k] = v
		}
	}
	if t.Metadata.Annotations != nil {
		cp.Metadata.Annotations = make(map[string]string, len(t.Metadata.Annotations))
		for k, v := range t.Metadata.Annotations {
			cp.Metadata.Annotations[k] = v
		}
	}
	return cp
}

// CreateTenant stores a deep copy of t keyed by the composite key
// "organizationName/organizationUnitName/name". It returns ErrAlreadyExists if
// the composite key is already present.
func (r *TenantRegistry) CreateTenant(
	ctx context.Context, t resources.Tenant,
) (resources.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := tenantCompositeKey(t.Spec.OrganizationName, t.Spec.OrganizationUnitName, t.Metadata.Name)
	if _, ok := r.store[key]; ok {
		return resources.Tenant{}, ErrAlreadyExists
	}
	stored := deepCopyTenant(t)
	r.store[key] = stored
	return deepCopyTenant(stored), nil
}

// GetTenant returns a deep copy of the stored Tenant identified by the
// composite key, or ErrNotFound if absent.
func (r *TenantRegistry) GetTenant(
	ctx context.Context, orgName, ouName, name string,
) (resources.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.store[tenantCompositeKey(orgName, ouName, name)]
	if !ok {
		return resources.Tenant{}, ErrNotFound
	}
	return deepCopyTenant(t), nil
}

// ListTenants returns a new slice of deep copies sorted by
// spec.organizationName ascending, then spec.organizationUnitName ascending,
// then metadata.name ascending. It returns a non-nil empty slice when no
// Tenants are stored.
func (r *TenantRegistry) ListTenants(
	ctx context.Context,
) ([]resources.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]resources.Tenant, 0, len(r.store))
	for _, t := range r.store {
		items = append(items, deepCopyTenant(t))
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Spec.OrganizationName != items[j].Spec.OrganizationName {
			return items[i].Spec.OrganizationName < items[j].Spec.OrganizationName
		}
		if items[i].Spec.OrganizationUnitName != items[j].Spec.OrganizationUnitName {
			return items[i].Spec.OrganizationUnitName < items[j].Spec.OrganizationUnitName
		}
		return items[i].Metadata.Name < items[j].Metadata.Name
	})
	return items, nil
}

// UpdateTenant derives the composite key from the submitted Tenant, looks up
// the existing stored entry, replaces only the mutable fields, and returns a
// deep copy of the updated Tenant. Immutable and system-owned fields
// (apiVersion, kind, status, metadata.name, spec.organizationName,
// spec.organizationUnitName) are preserved from the stored entry. Returns
// ErrNotFound if the composite key is absent.
func (r *TenantRegistry) UpdateTenant(
	ctx context.Context, t resources.Tenant,
) (resources.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := tenantCompositeKey(t.Spec.OrganizationName, t.Spec.OrganizationUnitName, t.Metadata.Name)
	existing, ok := r.store[key]
	if !ok {
		return resources.Tenant{}, ErrNotFound
	}
	merged := existing
	merged.Metadata.DisplayName = t.Metadata.DisplayName
	merged.Metadata.Labels = t.Metadata.Labels
	merged.Metadata.Annotations = t.Metadata.Annotations
	merged.Spec.Description = t.Spec.Description
	stored := deepCopyTenant(merged)
	r.store[key] = stored
	return deepCopyTenant(stored), nil
}

// DeleteTenant removes the entry identified by the composite key, or returns
// ErrNotFound if absent.
func (r *TenantRegistry) DeleteTenant(
	ctx context.Context, orgName, ouName, name string,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := tenantCompositeKey(orgName, ouName, name)
	if _, ok := r.store[key]; !ok {
		return ErrNotFound
	}
	delete(r.store, key)
	return nil
}

// CountByOrganizationUnit returns the number of stored Tenants whose
// spec.organizationName and spec.organizationUnitName match the given parent.
func (r *TenantRegistry) CountByOrganizationUnit(
	ctx context.Context, orgName, ouName string,
) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, t := range r.store {
		if t.Spec.OrganizationName == orgName && t.Spec.OrganizationUnitName == ouName {
			count++
		}
	}
	return count, nil
}
