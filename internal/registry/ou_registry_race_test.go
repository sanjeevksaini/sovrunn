package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// TestOURegistry_ConcurrentAccess launches many goroutines performing a mix
// of Create, Get, List, Update, Delete, and CountByOrganization operations
// against a single shared registry. It verifies no panic occurs and that any
// returned errors are only the expected concurrency-valid sentinels. Run
// under `go test -race` to detect data races.
func TestOURegistry_ConcurrentAccess(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	ctx := context.Background()
	var wg sync.WaitGroup

	const goroutines = 16
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			orgName := fmt.Sprintf("org-%d", id%4)
			name := fmt.Sprintf("ou-%d", id%8)
			ou := resources.OrganizationUnit{
				APIVersion: resources.OUAPIVersion,
				Kind:       resources.OUKind,
				Metadata: resources.Metadata{
					Name:        name,
					Labels:      map[string]string{"k": "v"},
					Annotations: map[string]string{"a": "b"},
				},
				Spec: resources.OrganizationUnitSpec{
					OrganizationName: orgName,
					Description:      "desc",
				},
				Status: resources.OrganizationUnitStatus{Phase: resources.PhaseActive},
			}

			if _, err := reg.CreateOrganizationUnit(ctx, ou); err != nil && !errors.Is(err, ErrAlreadyExists) {
				t.Errorf("Create unexpected error: %v", err)
			}
			if _, err := reg.GetOrganizationUnit(ctx, orgName, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Get unexpected error: %v", err)
			}
			if _, err := reg.ListOrganizationUnits(ctx); err != nil {
				t.Errorf("List unexpected error: %v", err)
			}
			update := ou
			update.Metadata.DisplayName = fmt.Sprintf("display-%d", id)
			if _, err := reg.UpdateOrganizationUnit(ctx, orgName, name, update); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Update unexpected error: %v", err)
			}
			if _, err := reg.CountByOrganization(ctx, orgName); err != nil {
				t.Errorf("CountByOrganization unexpected error: %v", err)
			}
			if err := reg.DeleteOrganizationUnit(ctx, orgName, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Delete unexpected error: %v", err)
			}
		}(i)
	}
	wg.Wait()
}
