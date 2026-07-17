package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// TestServiceClassRegistry_ConcurrentAccess launches many goroutines performing
// a mix of Create, Get, List, Update, and Delete against a single shared
// registry. It verifies no panic occurs and that any returned errors are only
// the expected concurrency-valid sentinels. Run under `go test -race`.
func TestServiceClassRegistry_ConcurrentAccess(t *testing.T) {
	reg := NewServiceClassRegistry()
	ctx := context.Background()
	var wg sync.WaitGroup

	const goroutines = 16
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			name := fmt.Sprintf("class-%d", id%8)
			sc := resources.ServiceClass{
				APIVersion: "platform.sovrunn.io/v1alpha1",
				Kind:       resources.ServiceClassKind,
				Metadata: resources.Metadata{
					Name:        name,
					Labels:      map[string]string{"k": "v"},
					Annotations: map[string]string{"a": "b"},
				},
				Spec: resources.ServiceClassSpec{
					Category:  resources.CategoryDatabase,
					Lifecycle: resources.LifecycleActive,
					Tags:      []string{"tag"},
				},
				Status: resources.ServiceClassStatus{Phase: resources.PhaseActive},
			}

			if _, err := reg.CreateServiceClass(ctx, sc); err != nil && !errors.Is(err, ErrAlreadyExists) {
				t.Errorf("Create unexpected error: %v", err)
			}
			if _, err := reg.GetServiceClass(ctx, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Get unexpected error: %v", err)
			}
			if _, err := reg.ListServiceClasses(ctx); err != nil {
				t.Errorf("List unexpected error: %v", err)
			}
			update := sc
			update.Spec.Description = fmt.Sprintf("desc-%d", id)
			if _, err := reg.UpdateServiceClass(ctx, update); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Update unexpected error: %v", err)
			}
			if err := reg.DeleteServiceClass(ctx, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Delete unexpected error: %v", err)
			}
		}(i)
	}
	wg.Wait()
}
