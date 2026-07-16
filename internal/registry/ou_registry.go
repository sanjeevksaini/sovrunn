package registry

import (
	"context"
	"sort"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// OrganizationUnitRegistryIface is the storage contract for
// OrganizationUnit resources. The registry is storage-only: it does not
// depend on OrganizationRegistry and does not perform parent Organization
// existence checks. Those checks belong to the API/service layer.
type OrganizationUnitRegistryIface interface {
	CreateOrganizationUnit(ctx context.Context, ou resources.OrganizationUnit) (resources.OrganizationUnit, error)
	GetOrganizationUnit(ctx context.Context, orgName, name string) (resources.OrganizationUnit, error)
	ListOrganizationUnits(ctx context.Context) ([]resources.OrganizationUnit, error)
	UpdateOrganizationUnit(ctx context.Context, orgName, name string, ou resources.OrganizationUnit) (resources.OrganizationUnit, error)
	DeleteOrganizationUnit(ctx context.Context, orgName, name string) error
	CountByOrganization(ctx context.Context, orgName string) (int, error)
}

// OrganizationLookup is a narrow interface for verifying parent
// Organization existence. The existing *OrganizationRegistry already
// satisfies it via GetOrganization. It is intended for injection into the
// OUHandler in a later task and is NOT used inside OrganizationUnitRegistry.
type OrganizationLookup interface {
	GetOrganization(ctx context.Context, name string) (resources.Organization, error)
}

// OrganizationUnitRegistry is the Phase 1 in-memory implementation of
// OrganizationUnitRegistryIface. All public methods are safe for concurrent
// use. The registry holds no package-level global state and uses the
// composite key "organizationName/name" as the map key.
type OrganizationUnitRegistry struct {
	mu    sync.RWMutex
	store map[string]resources.OrganizationUnit
}

// NewOrganizationUnitRegistry returns a ready-to-use registry.
func NewOrganizationUnitRegistry() *OrganizationUnitRegistry {
	return &OrganizationUnitRegistry{
		store: make(map[string]resources.OrganizationUnit),
	}
}

// compositeKey builds the map key from orgName and name.
func compositeKey(orgName, name string) string {
	return orgName + "/" + name
}

// deepCopyOrganizationUnit returns a fully independent copy of ou,
// duplicating the Labels and Annotations maps so that callers cannot
// mutate the registry's internal state.
func deepCopyOrganizationUnit(ou resources.OrganizationUnit) resources.OrganizationUnit {
	cp := ou
	if ou.Metadata.Labels != nil {
		cp.Metadata.Labels = make(map[string]string, len(ou.Metadata.Labels))
		for k, v := range ou.Metadata.Labels {
			cp.Metadata.Labels[k] = v
		}
	}
	if ou.Metadata.Annotations != nil {
		cp.Metadata.Annotations = make(map[string]string, len(ou.Metadata.Annotations))
		for k, v := range ou.Metadata.Annotations {
			cp.Metadata.Annotations[k] = v
		}
	}
	return cp
}

// CreateOrganizationUnit stores a deep copy of ou keyed by the composite
// key "organizationName/name". It returns ErrAlreadyExists if the composite
// key is already present.
func (r *OrganizationUnitRegistry) CreateOrganizationUnit(
	ctx context.Context, ou resources.OrganizationUnit,
) (resources.OrganizationUnit, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := compositeKey(ou.Spec.OrganizationName, ou.Metadata.Name)
	if _, ok := r.store[key]; ok {
		return resources.OrganizationUnit{}, ErrAlreadyExists
	}
	stored := deepCopyOrganizationUnit(ou)
	r.store[key] = stored
	return deepCopyOrganizationUnit(stored), nil
}

// GetOrganizationUnit returns a deep copy of the stored OrganizationUnit
// identified by the composite key, or ErrNotFound if absent.
func (r *OrganizationUnitRegistry) GetOrganizationUnit(
	ctx context.Context, orgName, name string,
) (resources.OrganizationUnit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ou, ok := r.store[compositeKey(orgName, name)]
	if !ok {
		return resources.OrganizationUnit{}, ErrNotFound
	}
	return deepCopyOrganizationUnit(ou), nil
}

// ListOrganizationUnits returns a new slice of deep copies sorted by
// spec.organizationName ascending, then metadata.name ascending. It
// returns an empty slice when no OrganizationUnits are stored.
func (r *OrganizationUnitRegistry) ListOrganizationUnits(
	ctx context.Context,
) ([]resources.OrganizationUnit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]resources.OrganizationUnit, 0, len(r.store))
	for _, ou := range r.store {
		items = append(items, deepCopyOrganizationUnit(ou))
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Spec.OrganizationName != items[j].Spec.OrganizationName {
			return items[i].Spec.OrganizationName < items[j].Spec.OrganizationName
		}
		return items[i].Metadata.Name < items[j].Metadata.Name
	})
	return items, nil
}

// UpdateOrganizationUnit replaces the mutable fields of the stored entry
// and returns a deep copy of the updated OrganizationUnit. Immutable and
// system-owned fields (metadata.name, spec.organizationName, status,
// apiVersion, kind) are preserved from the stored entry.
func (r *OrganizationUnitRegistry) UpdateOrganizationUnit(
	ctx context.Context, orgName, name string, ou resources.OrganizationUnit,
) (resources.OrganizationUnit, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := compositeKey(orgName, name)
	existing, ok := r.store[key]
	if !ok {
		return resources.OrganizationUnit{}, ErrNotFound
	}
	merged := existing
	merged.Metadata.DisplayName = ou.Metadata.DisplayName
	merged.Metadata.Labels = ou.Metadata.Labels
	merged.Metadata.Annotations = ou.Metadata.Annotations
	merged.Spec.Description = ou.Spec.Description
	stored := deepCopyOrganizationUnit(merged)
	r.store[key] = stored
	return deepCopyOrganizationUnit(stored), nil
}

// DeleteOrganizationUnit removes the entry identified by the composite key,
// or returns ErrNotFound if absent.
func (r *OrganizationUnitRegistry) DeleteOrganizationUnit(
	ctx context.Context, orgName, name string,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := compositeKey(orgName, name)
	if _, ok := r.store[key]; !ok {
		return ErrNotFound
	}
	delete(r.store, key)
	return nil
}

// CountByOrganization returns the number of stored OrganizationUnits whose
// spec.organizationName equals orgName.
func (r *OrganizationUnitRegistry) CountByOrganization(
	ctx context.Context, orgName string,
) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, ou := range r.store {
		if ou.Spec.OrganizationName == orgName {
			count++
		}
	}
	return count, nil
}
