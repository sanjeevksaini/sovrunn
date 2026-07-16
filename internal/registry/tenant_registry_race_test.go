package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// TestTenantRegistry_ConcurrentAccess launches many goroutines performing a
// mix of Create, Get, List, Update, Delete, and CountByOrganizationUnit
// operations against a single shared registry. It verifies no panic occurs and
// that any returned errors are only the expected concurrency-valid sentinels.
// Run under `go test -race` to detect data races.
func TestTenantRegistry_ConcurrentAccess(t *testing.T) {
	reg := NewTenantRegistry()
	ctx := context.Background()
	var wg sync.WaitGroup

	const goroutines = 16
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			orgName := fmt.Sprintf("org-%d", id%4)
			ouName := fmt.Sprintf("ou-%d", id%3)
			name := fmt.Sprintf("tenant-%d", id%8)
			tnt := resources.Tenant{
				APIVersion: resources.TenantAPIVersion,
				Kind:       resources.TenantKind,
				Metadata: resources.Metadata{
					Name:        name,
					Labels:      map[string]string{"k": "v"},
					Annotations: map[string]string{"a": "b"},
				},
				Spec: resources.TenantSpec{
					OrganizationName:     orgName,
					OrganizationUnitName: ouName,
					Description:          "desc",
				},
				Status: resources.TenantStatus{Phase: resources.PhaseActive},
			}

			if _, err := reg.CreateTenant(ctx, tnt); err != nil && !errors.Is(err, ErrAlreadyExists) {
				t.Errorf("Create unexpected error: %v", err)
			}
			if _, err := reg.GetTenant(ctx, orgName, ouName, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Get unexpected error: %v", err)
			}
			if _, err := reg.ListTenants(ctx); err != nil {
				t.Errorf("List unexpected error: %v", err)
			}
			update := tnt
			update.Metadata.DisplayName = fmt.Sprintf("display-%d", id)
			if _, err := reg.UpdateTenant(ctx, update); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Update unexpected error: %v", err)
			}
			if _, err := reg.CountByOrganizationUnit(ctx, orgName, ouName); err != nil {
				t.Errorf("CountByOrganizationUnit unexpected error: %v", err)
			}
			if err := reg.DeleteTenant(ctx, orgName, ouName, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Delete unexpected error: %v", err)
			}
		}(i)
	}
	wg.Wait()
}
