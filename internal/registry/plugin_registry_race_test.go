package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// TestPluginRegistry_ConcurrentAccess launches many goroutines performing
// a mix of Create, Get, List, Update, and Delete against a single shared
// registry. It verifies no panic occurs and that any returned errors are only
// the expected concurrency-valid sentinels. Run under `go test -race`.
func TestPluginRegistry_ConcurrentAccess(t *testing.T) {
	reg := NewPluginRegistry()
	ctx := context.Background()
	var wg sync.WaitGroup

	const goroutines = 16
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			name := fmt.Sprintf("plugin-%d", id%8)
			p := resources.Plugin{
				APIVersion: resources.PluginAPIVersion,
				Kind:       resources.PluginKind,
				Metadata: resources.Metadata{
					Name:        name,
					Labels:      map[string]string{"k": "v"},
					Annotations: map[string]string{"a": "b"},
				},
				Spec: resources.PluginSpec{
					PluginType:       resources.PluginTypeDStoreOps,
					Version:          "0.1.0",
					ServiceClassRefs: []string{"datastore.postgresql"},
					DeploymentMode:   resources.DeploymentModeCompiledIn,
					Tags:             []string{"tag"},
				},
				Status: resources.PluginStatus{Phase: resources.PhaseActive},
			}

			if _, err := reg.CreatePlugin(ctx, p); err != nil && !errors.Is(err, ErrAlreadyExists) {
				t.Errorf("Create unexpected error: %v", err)
			}
			if _, err := reg.GetPlugin(ctx, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Get unexpected error: %v", err)
			}
			if _, err := reg.ListPlugins(ctx); err != nil {
				t.Errorf("List unexpected error: %v", err)
			}
			update := p
			update.Spec.Description = fmt.Sprintf("desc-%d", id)
			if _, err := reg.UpdatePlugin(ctx, update); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Update unexpected error: %v", err)
			}
			if err := reg.DeletePlugin(ctx, name); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Delete unexpected error: %v", err)
			}
		}(i)
	}
	wg.Wait()
}
