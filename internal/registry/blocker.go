package registry

import "context"

// BlockedBy describes a resource kind that prevents deletion.
type BlockedBy struct {
	Kind  string
	Count int
}

// ChildBlockerChecker is injected into the delete handler. Phase 1
// uses NoopChildBlockerChecker which always returns an empty slice.
type ChildBlockerChecker interface {
	BlockedByChildren(ctx context.Context, orgName string) ([]BlockedBy, error)
}

// NoopChildBlockerChecker is the Phase 1 stub — always returns empty.
type NoopChildBlockerChecker struct{}

// BlockedByChildren always returns an empty slice in Phase 1.
func (NoopChildBlockerChecker) BlockedByChildren(
	_ context.Context, _ string,
) ([]BlockedBy, error) {
	return nil, nil
}
