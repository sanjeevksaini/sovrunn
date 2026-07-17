package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestServicePlanChildBlockerChecker_EmptyWhenZero(t *testing.T) {
	reg := NewServicePlanRegistry()
	blocker := NewServicePlanChildBlockerChecker(reg)

	blockers, err := blocker.BlockedByServiceClassChildren(context.Background(), "postgres")
	if err != nil {
		t.Fatalf("BlockedByServiceClassChildren() error = %v", err)
	}
	if blockers != nil {
		t.Errorf("blockers = %+v, want nil", blockers)
	}
}

func TestServicePlanChildBlockerChecker_BlocksWhenReferenced(t *testing.T) {
	reg := NewServicePlanRegistry()
	_, _ = reg.CreateServicePlan(context.Background(), sampleServicePlan("postgres", "small"))
	_, _ = reg.CreateServicePlan(context.Background(), sampleServicePlan("postgres", "large"))
	_, _ = reg.CreateServicePlan(context.Background(), sampleServicePlan("redis", "small"))

	blocker := NewServicePlanChildBlockerChecker(reg)
	blockers, err := blocker.BlockedByServiceClassChildren(context.Background(), "postgres")
	if err != nil {
		t.Fatalf("BlockedByServiceClassChildren() error = %v", err)
	}
	if len(blockers) != 1 {
		t.Fatalf("blockers = %+v, want exactly one entry", blockers)
	}
	if blockers[0].Kind != "ServicePlan" {
		t.Errorf("Kind = %q, want ServicePlan", blockers[0].Kind)
	}
	if blockers[0].Count != 2 {
		t.Errorf("Count = %d, want 2", blockers[0].Count)
	}

	// Confirm count is scoped to the requested ServiceClass only.
	redisBlockers, err := blocker.BlockedByServiceClassChildren(context.Background(), "redis")
	if err != nil {
		t.Fatalf("BlockedByServiceClassChildren(redis) error = %v", err)
	}
	if len(redisBlockers) != 1 || redisBlockers[0].Kind != "ServicePlan" || redisBlockers[0].Count != 1 {
		t.Errorf("redis blockers = %+v, want [{ServicePlan 1}]", redisBlockers)
	}
}

func TestServicePlanChildBlockerChecker_PropagatesRegistryError(t *testing.T) {
	wantErr := errors.New("count failed")
	blocker := NewServicePlanChildBlockerChecker(failingServicePlanCountRegistry{err: wantErr})

	blockers, err := blocker.BlockedByServiceClassChildren(context.Background(), "postgres")
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if blockers != nil {
		t.Fatalf("blockers = %+v, want nil", blockers)
	}
}

func TestServicePlanChildBlockerChecker_RetiredLifecycleStillBlocks(t *testing.T) {
	reg := NewServicePlanRegistry()
	retired := sampleServicePlan("postgres", "legacy")
	retired.Spec.Lifecycle = resources.LifecycleRetired
	if _, err := reg.CreateServicePlan(context.Background(), retired); err != nil {
		t.Fatalf("CreateServicePlan() error = %v", err)
	}

	blocker := NewServicePlanChildBlockerChecker(reg)
	blockers, err := blocker.BlockedByServiceClassChildren(context.Background(), "postgres")
	if err != nil {
		t.Fatalf("BlockedByServiceClassChildren() error = %v", err)
	}
	if len(blockers) != 1 {
		t.Fatalf("blockers = %+v, want exactly one entry", blockers)
	}
	if blockers[0].Kind != "ServicePlan" {
		t.Errorf("Kind = %q, want ServicePlan", blockers[0].Kind)
	}
	if blockers[0].Count != 1 {
		t.Errorf("Count = %d, want 1", blockers[0].Count)
	}
}

type failingServicePlanCountRegistry struct {
	ServicePlanRegistryIface
	err error
}

func (r failingServicePlanCountRegistry) CountByServiceClass(
	ctx context.Context, serviceClassName string,
) (int, error) {
	return 0, r.err
}
