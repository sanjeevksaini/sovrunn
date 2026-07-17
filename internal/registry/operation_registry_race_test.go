package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
)

// TestOperationRegistry_ConcurrentAccess launches many goroutines performing a
// mix of Create, Get, and List operations against a single shared registry. It
// verifies no panic occurs and that any returned errors are only the expected
// concurrency-valid sentinels. Run under `go test -race` to detect data races.
func TestOperationRegistry_ConcurrentAccess(t *testing.T) {
	reg := NewOperationRegistry()
	ctx := context.Background()
	var wg sync.WaitGroup

	const goroutines = 16
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Unique ID per goroutine so Create never collides.
			opID := fmt.Sprintf("op-%d", id)
			createdAt := fmt.Sprintf("2026-01-01T00:00:%02dZ", id)

			if _, err := reg.CreateOperation(ctx, testOperation(opID, createdAt)); err != nil {
				t.Errorf("Create unexpected error: %v", err)
			}
			// Read own ID plus a neighbor's ID, which may or may not exist yet.
			if _, err := reg.GetOperation(ctx, opID); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Get unexpected error: %v", err)
			}
			neighbor := fmt.Sprintf("op-%d", (id+1)%goroutines)
			if _, err := reg.GetOperation(ctx, neighbor); err != nil && !errors.Is(err, ErrNotFound) {
				t.Errorf("Get neighbor unexpected error: %v", err)
			}
			if _, err := reg.ListOperations(ctx); err != nil {
				t.Errorf("List unexpected error: %v", err)
			}
		}(i)
	}
	wg.Wait()
}
