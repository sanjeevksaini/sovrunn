package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// TestServiceInstanceRegistry_ConcurrentAccess launches many goroutines
// performing a mix of Create, Get, List, Update, Delete, CountByServicePlan,
// and CountByProject against a single shared registry. It verifies no panic
// occurs and that any returned errors are only the expected concurrency-valid
// sentinels. Run under `go test -race`.
func TestServiceInstanceRegistry_ConcurrentAccess(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	ctx := context.Background()
	var wg sync.WaitGroup

	const goroutines = 16
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			name := fmt.Sprintf("si-%d", id%8)
			org := fmt.Sprintf("org-%d", id%4)
			ou := fmt.Sprintf("ou-%d", id%3)
			tenant := fmt.Sprintf("tenant-%d", id%5)
			project := fmt.Sprintf("project-%d", id%4)
			class := fmt.Sprintf("class-%d", id%3)
			plan := fmt.Sprintf("plan-%d", id%2)
			si := resources.ServiceInstance{
				APIVersion: resources.ServiceInstanceAPIVersion,
				Kind:       resources.ServiceInstanceKind,
				Metadata: resources.Metadata{
					Name:        name,
					Labels:      map[string]string{"k": "v"},
					Annotations: map[string]string{"a": "b"},
				},
				Spec: resources.ServiceInstanceSpec{
					OrganizationRef:     org,
					OrganizationUnitRef: ou,
					TenantRef:           tenant,
					ProjectRef:          project,
					ServiceClassRef:     class,
					ServicePlanRef:      plan,
					Parameters:          map[string]string{"region": "us-east"},
				},
				Status: resources.ServiceInstanceStatus{Phase: "Ready"},
			}

			if _, err := reg.CreateServiceInstance(ctx, si); err != nil && !errors.Is(err, ErrAlreadyExists) {
				t.Errorf("Create unexpected error: %v", err)
			}
			if _, err := reg.GetServiceInstance(ctx, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Get unexpected error: %v", err)
			}
			if _, err := reg.ListServiceInstances(ctx, tenant, project); err != nil {
				t.Errorf("List unexpected error: %v", err)
			}
			update := si
			update.Metadata.DisplayName = fmt.Sprintf("display-%d", id)
			update.Spec.Parameters = map[string]string{"region": fmt.Sprintf("r-%d", id)}
			if _, err := reg.UpdateServiceInstance(ctx, name, update); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Update unexpected error: %v", err)
			}
			if _, err := reg.CountByServicePlan(ctx, class, plan); err != nil {
				t.Errorf("CountByServicePlan unexpected error: %v", err)
			}
			if _, err := reg.CountByProject(ctx, org, ou, tenant, project); err != nil {
				t.Errorf("CountByProject unexpected error: %v", err)
			}
			if err := reg.DeleteServiceInstance(ctx, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Delete unexpected error: %v", err)
			}
		}(i)
	}
	wg.Wait()
}
