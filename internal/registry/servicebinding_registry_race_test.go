package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// TestServiceBindingRegistry_ConcurrentAccess launches more than ten
// goroutines performing Create, Get, List, Delete, and CountByServiceInstance
// against one shared registry. Run under `go test -race`.
func TestServiceBindingRegistry_ConcurrentAccess(t *testing.T) {
	reg := NewServiceBindingRegistry()
	ctx := context.Background()
	var wg sync.WaitGroup

	const goroutines = 16
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			name := fmt.Sprintf("binding-%d", id%8)
			instanceRef := fmt.Sprintf("instance-%d", id%4)
			sb := resources.ServiceBinding{
				APIVersion: resources.ServiceBindingAPIVersion,
				Kind:       resources.ServiceBindingKind,
				Metadata: resources.Metadata{
					Name:        name,
					Labels:      map[string]string{"worker": fmt.Sprintf("%d", id)},
					Annotations: map[string]string{"source": "race-test"},
				},
				Spec: resources.ServiceBindingSpec{
					ServiceInstanceRef: instanceRef,
					ConsumerRef: &resources.ConsumerRef{
						Kind: "Application",
						Name: fmt.Sprintf("consumer-%d", id),
					},
					BindingType: resources.BindingTypeCredentials,
				},
				Status: resources.ServiceBindingStatus{
					Phase:     "Ready",
					SecretRef: "stub-secret-ref",
				},
			}

			if _, err := reg.CreateServiceBinding(ctx, sb); err != nil && !errors.Is(err, ErrAlreadyExists) {
				t.Errorf("CreateServiceBinding() unexpected error: %v", err)
			}
			if _, err := reg.GetServiceBinding(ctx, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("GetServiceBinding() unexpected error: %v", err)
			}
			if _, err := reg.ListServiceBindings(ctx, instanceRef); err != nil {
				t.Errorf("ListServiceBindings() unexpected error: %v", err)
			}
			if _, err := reg.CountByServiceInstance(ctx, instanceRef); err != nil {
				t.Errorf("CountByServiceInstance() unexpected error: %v", err)
			}
			if err := reg.DeleteServiceBinding(ctx, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("DeleteServiceBinding() unexpected error: %v", err)
			}
		}(i)
	}
	wg.Wait()
}
