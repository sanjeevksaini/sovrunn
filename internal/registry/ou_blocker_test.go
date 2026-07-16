package registry

import (
	"context"
	"testing"
)

func TestCountByOrganization_ZeroForNoOUs(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	count, err := reg.CountByOrganization(context.Background(), "nic")
	if err != nil {
		t.Fatalf("CountByOrganization() error = %v", err)
	}
	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}
}

func TestCountByOrganization_MultipleOUs(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	_, _ = reg.CreateOrganizationUnit(context.Background(), sampleOU("nic", "ministry-health"))
	_, _ = reg.CreateOrganizationUnit(context.Background(), sampleOU("nic", "ministry-finance"))
	_, _ = reg.CreateOrganizationUnit(context.Background(), sampleOU("nic", "ministry-defence"))
	_, _ = reg.CreateOrganizationUnit(context.Background(), sampleOU("state-gov", "ministry-health"))

	count, err := reg.CountByOrganization(context.Background(), "nic")
	if err != nil {
		t.Fatalf("CountByOrganization() error = %v", err)
	}
	if count != 3 {
		t.Errorf("count for nic = %d, want 3", count)
	}
}

func TestOUChildBlockerChecker_EmptyWhenZero(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	blocker := NewOUChildBlockerChecker(reg)
	blockers, err := blocker.BlockedByChildren(context.Background(), "nic")
	if err != nil {
		t.Fatalf("BlockedByChildren() error = %v", err)
	}
	if len(blockers) != 0 {
		t.Errorf("blockers = %+v, want empty", blockers)
	}
}

func TestOUChildBlockerChecker_BlocksWhenReferenced(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	_, _ = reg.CreateOrganizationUnit(context.Background(), sampleOU("nic", "ministry-health"))
	_, _ = reg.CreateOrganizationUnit(context.Background(), sampleOU("nic", "ministry-finance"))

	blocker := NewOUChildBlockerChecker(reg)
	blockers, err := blocker.BlockedByChildren(context.Background(), "nic")
	if err != nil {
		t.Fatalf("BlockedByChildren() error = %v", err)
	}
	if len(blockers) != 1 {
		t.Fatalf("blockers = %+v, want exactly one entry", blockers)
	}
	if blockers[0].Kind != "OrganizationUnit" {
		t.Errorf("Kind = %q, want OrganizationUnit", blockers[0].Kind)
	}
	if blockers[0].Count != 2 {
		t.Errorf("Count = %d, want 2", blockers[0].Count)
	}
}
