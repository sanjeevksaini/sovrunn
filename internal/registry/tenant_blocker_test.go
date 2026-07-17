package registry

import (
	"context"
	"errors"
	"testing"
)

func TestTenantChildBlockerChecker_EmptyWhenZero(t *testing.T) {
	reg := NewTenantRegistry()
	blocker := NewTenantChildBlockerChecker(reg)

	blockers, err := blocker.BlockedByOUChildren(context.Background(), "nic", "ministry-health")
	if err != nil {
		t.Fatalf("BlockedByOUChildren() error = %v", err)
	}
	if len(blockers) != 0 {
		t.Errorf("blockers = %+v, want empty", blockers)
	}
}

func TestTenantChildBlockerChecker_BlocksWhenReferenced(t *testing.T) {
	reg := NewTenantRegistry()
	_, _ = reg.CreateTenant(context.Background(), sampleTenant("nic", "ministry-health", "prod"))
	_, _ = reg.CreateTenant(context.Background(), sampleTenant("nic", "ministry-health", "staging"))
	_, _ = reg.CreateTenant(context.Background(), sampleTenant("nic", "ministry-finance", "prod"))

	blocker := NewTenantChildBlockerChecker(reg)
	blockers, err := blocker.BlockedByOUChildren(context.Background(), "nic", "ministry-health")
	if err != nil {
		t.Fatalf("BlockedByOUChildren() error = %v", err)
	}
	if len(blockers) != 1 {
		t.Fatalf("blockers = %+v, want exactly one entry", blockers)
	}
	if blockers[0].Kind != "Tenant" {
		t.Errorf("Kind = %q, want Tenant", blockers[0].Kind)
	}
	if blockers[0].Count != 2 {
		t.Errorf("Count = %d, want 2", blockers[0].Count)
	}
}

func TestTenantChildBlockerChecker_PropagatesRegistryError(t *testing.T) {
	wantErr := errors.New("count failed")
	blocker := NewTenantChildBlockerChecker(failingTenantCountRegistry{err: wantErr})

	blockers, err := blocker.BlockedByOUChildren(context.Background(), "nic", "ministry-health")
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if blockers != nil {
		t.Fatalf("blockers = %+v, want nil", blockers)
	}
}

type failingTenantCountRegistry struct {
	TenantRegistryIface
	err error
}

func (r failingTenantCountRegistry) CountByOrganizationUnit(
	ctx context.Context, orgName, ouName string,
) (int, error) {
	return 0, r.err
}
