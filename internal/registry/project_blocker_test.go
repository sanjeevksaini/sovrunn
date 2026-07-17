package registry

import (
	"context"
	"errors"
	"testing"
)

func TestProjectChildBlockerChecker_EmptyWhenZero(t *testing.T) {
	reg := NewProjectRegistry()
	blocker := NewProjectChildBlockerChecker(reg)

	blockers, err := blocker.BlockedByTenantChildren(context.Background(), "nic", "ministry-health", "payments")
	if err != nil {
		t.Fatalf("BlockedByTenantChildren() error = %v", err)
	}
	if len(blockers) != 0 {
		t.Errorf("blockers = %+v, want empty", blockers)
	}
}

func TestProjectChildBlockerChecker_BlocksWhenReferenced(t *testing.T) {
	reg := NewProjectRegistry()
	_, _ = reg.CreateProject(context.Background(), sampleProject("nic", "ministry-health", "payments", "prod"))
	_, _ = reg.CreateProject(context.Background(), sampleProject("nic", "ministry-health", "payments", "staging"))
	_, _ = reg.CreateProject(context.Background(), sampleProject("nic", "ministry-health", "billing", "prod"))

	blocker := NewProjectChildBlockerChecker(reg)
	blockers, err := blocker.BlockedByTenantChildren(context.Background(), "nic", "ministry-health", "payments")
	if err != nil {
		t.Fatalf("BlockedByTenantChildren() error = %v", err)
	}
	if len(blockers) != 1 {
		t.Fatalf("blockers = %+v, want exactly one entry", blockers)
	}
	if blockers[0].Kind != "Project" {
		t.Errorf("Kind = %q, want Project", blockers[0].Kind)
	}
	if blockers[0].Count != 2 {
		t.Errorf("Count = %d, want 2", blockers[0].Count)
	}
}

func TestProjectChildBlockerChecker_PropagatesRegistryError(t *testing.T) {
	wantErr := errors.New("count failed")
	blocker := NewProjectChildBlockerChecker(failingProjectCountRegistry{err: wantErr})

	blockers, err := blocker.BlockedByTenantChildren(context.Background(), "nic", "ministry-health", "payments")
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if blockers != nil {
		t.Fatalf("blockers = %+v, want nil", blockers)
	}
}

type failingProjectCountRegistry struct {
	ProjectRegistryIface
	err error
}

func (r failingProjectCountRegistry) CountByTenant(
	ctx context.Context, orgName, ouName, tenantName string,
) (int, error) {
	return 0, r.err
}
