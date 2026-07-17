package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// TestServicePlanRegistry_ConcurrentAccess launches many goroutines performing
// a mix of Create, Get, List, Update, Delete, and CountByServiceClass against a
// single shared registry. It verifies no panic occurs and that any returned
// errors are only the expected concurrency-valid sentinels. Run under
// `go test -race`.
func TestServicePlanRegistry_ConcurrentAccess(t *testing.T) {
	reg := NewServicePlanRegistry()
	ctx := context.Background()
	var wg sync.WaitGroup

	const goroutines = 16
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			className := fmt.Sprintf("class-%d", id%4)
			name := fmt.Sprintf("plan-%d", id%8)
			sp := resources.ServicePlan{
				APIVersion: "platform.sovrunn.io/v1alpha1",
				Kind:       resources.ServicePlanKind,
				Metadata: resources.Metadata{
					Name:        name,
					Labels:      map[string]string{"k": "v"},
					Annotations: map[string]string{"a": "b"},
				},
				Spec: resources.ServicePlanSpec{
					ServiceClassName: className,
					Tier:             resources.TierSmall,
					Lifecycle:        resources.LifecycleActive,
					Parameters:       map[string]string{"region": "us-east"},
					Tags:             []string{"tag"},
				},
				Status: resources.ServicePlanStatus{Phase: resources.PhaseActive},
			}

			if _, err := reg.CreateServicePlan(ctx, sp); err != nil && !errors.Is(err, ErrAlreadyExists) {
				t.Errorf("Create unexpected error: %v", err)
			}
			if _, err := reg.GetServicePlan(ctx, className, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Get unexpected error: %v", err)
			}
			if _, err := reg.ListServicePlans(ctx); err != nil {
				t.Errorf("List unexpected error: %v", err)
			}
			update := sp
			update.Spec.Description = fmt.Sprintf("desc-%d", id)
			if _, err := reg.UpdateServicePlan(ctx, update); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Update unexpected error: %v", err)
			}
			if _, err := reg.CountByServiceClass(ctx, className); err != nil {
				t.Errorf("CountByServiceClass unexpected error: %v", err)
			}
			if err := reg.DeleteServicePlan(ctx, className, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Delete unexpected error: %v", err)
			}
		}(i)
	}
	wg.Wait()
}
