package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// TestProjectRegistry_ConcurrentAccess launches many goroutines performing a
// mix of Create, Get, List, Update, Delete, and CountByTenant operations
// against a single shared registry. It verifies no panic occurs and that any
// returned errors are only the expected concurrency-valid sentinels. Run under
// `go test -race` to detect data races.
func TestProjectRegistry_ConcurrentAccess(t *testing.T) {
	reg := NewProjectRegistry()
	ctx := context.Background()
	var wg sync.WaitGroup

	const goroutines = 16
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			orgName := fmt.Sprintf("org-%d", id%4)
			ouName := fmt.Sprintf("ou-%d", id%3)
			tenantName := fmt.Sprintf("tenant-%d", id%5)
			name := fmt.Sprintf("project-%d", id%8)
			project := resources.Project{
				APIVersion: resources.ProjectAPIVersion,
				Kind:       resources.ProjectKind,
				Metadata: resources.Metadata{
					Name:        name,
					Labels:      map[string]string{"k": "v"},
					Annotations: map[string]string{"a": "b"},
				},
				Spec: resources.ProjectSpec{
					OrganizationName:     orgName,
					OrganizationUnitName: ouName,
					TenantName:           tenantName,
					Description:          "desc",
				},
				Status: resources.ProjectStatus{Phase: resources.PhaseActive},
			}

			if _, err := reg.CreateProject(ctx, project); err != nil && !errors.Is(err, ErrAlreadyExists) {
				t.Errorf("Create unexpected error: %v", err)
			}
			if _, err := reg.GetProject(ctx, orgName, ouName, tenantName, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Get unexpected error: %v", err)
			}
			if _, err := reg.ListProjects(ctx); err != nil {
				t.Errorf("List unexpected error: %v", err)
			}
			update := project
			update.Metadata.DisplayName = fmt.Sprintf("display-%d", id)
			if _, err := reg.UpdateProject(ctx, update); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Update unexpected error: %v", err)
			}
			if _, err := reg.CountByTenant(ctx, orgName, ouName, tenantName); err != nil {
				t.Errorf("CountByTenant unexpected error: %v", err)
			}
			if err := reg.DeleteProject(ctx, orgName, ouName, tenantName, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Delete unexpected error: %v", err)
			}
		}(i)
	}
	wg.Wait()
}
