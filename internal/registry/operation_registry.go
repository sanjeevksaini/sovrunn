package registry

import (
	"context"
	"sort"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// OperationRegistryIface is the storage contract for Operation records. The
// registry is storage-only: it does not generate Operation IDs and has no
// dependency on other registries. ID generation belongs to the emitter layer.
type OperationRegistryIface interface {
	CreateOperation(ctx context.Context, op resources.Operation) (resources.Operation, error)
	GetOperation(ctx context.Context, id string) (resources.Operation, error)
	ListOperations(ctx context.Context) ([]resources.Operation, error)
}

// OperationRegistry is the Phase 1 in-memory implementation of
// OperationRegistryIface. All public methods are safe for concurrent use. The
// registry holds no package-level global state and is keyed by Operation ID
// (metadata.name).
type OperationRegistry struct {
	mu    sync.RWMutex
	store map[string]resources.Operation
}

// NewOperationRegistry returns a ready-to-use registry.
func NewOperationRegistry() *OperationRegistry {
	return &OperationRegistry{
		store: make(map[string]resources.Operation),
	}
}

// deepCopyOperation returns a fully independent copy of op, duplicating the
// Labels and Annotations maps so that callers cannot mutate the registry's
// internal state. OperationSpec and OperationStatus are string-only structs, so
// a plain struct copy suffices for them.
func deepCopyOperation(op resources.Operation) resources.Operation {
	cp := op
	if op.Metadata.Labels != nil {
		cp.Metadata.Labels = make(map[string]string, len(op.Metadata.Labels))
		for k, v := range op.Metadata.Labels {
			cp.Metadata.Labels[k] = v
		}
	}
	if op.Metadata.Annotations != nil {
		cp.Metadata.Annotations = make(map[string]string, len(op.Metadata.Annotations))
		for k, v := range op.Metadata.Annotations {
			cp.Metadata.Annotations[k] = v
		}
	}
	return cp
}

// CreateOperation stores a deep copy of op keyed by its Operation ID
// (metadata.name). It returns ErrMissingOperationID if the ID is empty and
// ErrAlreadyExists if the ID is already present; in both error cases nothing is
// stored or overwritten.
func (r *OperationRegistry) CreateOperation(
	ctx context.Context, op resources.Operation,
) (resources.Operation, error) {
	if op.Metadata.Name == "" {
		return resources.Operation{}, ErrMissingOperationID
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[op.Metadata.Name]; ok {
		return resources.Operation{}, ErrAlreadyExists
	}
	stored := deepCopyOperation(op)
	r.store[op.Metadata.Name] = stored
	return deepCopyOperation(stored), nil
}

// GetOperation returns a deep copy of the stored Operation identified by id, or
// ErrNotFound if absent.
func (r *OperationRegistry) GetOperation(
	ctx context.Context, id string,
) (resources.Operation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	op, ok := r.store[id]
	if !ok {
		return resources.Operation{}, ErrNotFound
	}
	return deepCopyOperation(op), nil
}

// ListOperations returns a new slice of deep copies sorted by
// status.createdAt ascending, then metadata.name ascending as a tie-breaker. It
// returns a non-nil empty slice when no Operations are stored.
func (r *OperationRegistry) ListOperations(
	ctx context.Context,
) ([]resources.Operation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]resources.Operation, 0, len(r.store))
	for _, op := range r.store {
		items = append(items, deepCopyOperation(op))
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Status.CreatedAt != items[j].Status.CreatedAt {
			return items[i].Status.CreatedAt < items[j].Status.CreatedAt
		}
		return items[i].Metadata.Name < items[j].Metadata.Name
	})
	return items, nil
}
