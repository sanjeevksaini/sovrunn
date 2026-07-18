package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// TestCapabilityRegistry_ConcurrentAccess launches many goroutines performing
// a mix of Create, Get, List, Delete, and CountByPlugin against a single shared
// registry. It verifies no panic occurs and that any returned errors are only
// the expected concurrency-valid sentinels. Run under `go test -race`.
func TestCapabilityRegistry_ConcurrentAccess(t *testing.T) {
	reg := NewCapabilityRegistry()
	ctx := context.Background()
	var wg sync.WaitGroup

	const goroutines = 16
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			name := fmt.Sprintf("cap-%d", id%8)
			pluginRef := fmt.Sprintf("plugin-%d", id%4)
			c := resources.Capability{
				APIVersion: resources.CapabilityAPIVersion,
				Kind:       resources.CapabilityKind,
				Metadata: resources.Metadata{
					Name:        name,
					Labels:      map[string]string{"k": "v"},
					Annotations: map[string]string{"a": "b"},
				},
				Spec: resources.CapabilitySpec{
					PluginRef:       pluginRef,
					ServiceClassRef: "datastore.postgresql",
					Operation:       resources.CapOpProvision,
					Supported:       true,
				},
				Status: resources.CapabilityStatus{Phase: resources.PhaseActive},
			}

			if _, err := reg.CreateCapability(ctx, c); err != nil && !errors.Is(err, ErrAlreadyExists) {
				t.Errorf("Create unexpected error: %v", err)
			}
			if _, err := reg.GetCapability(ctx, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Get unexpected error: %v", err)
			}
			if _, err := reg.ListCapabilities(ctx, "", ""); err != nil {
				t.Errorf("List unexpected error: %v", err)
			}
			if _, err := reg.ListCapabilities(ctx, pluginRef, ""); err != nil {
				t.Errorf("List filter unexpected error: %v", err)
			}
			if _, err := reg.CountByPlugin(ctx, pluginRef); err != nil {
				t.Errorf("CountByPlugin unexpected error: %v", err)
			}
			if err := reg.DeleteCapability(ctx, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Delete unexpected error: %v", err)
			}
		}(i)
	}
	wg.Wait()
}
