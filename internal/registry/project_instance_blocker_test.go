package registry

import (
	"context"
	"testing"
)

func TestServiceInstanceProjectBlockerChecker_NotBlockedWithoutInstances(t *testing.T) {
	blocker := NewServiceInstanceProjectBlockerChecker(NewServiceInstanceRegistry())

	blockedBy, err := blocker.BlockedByProjectInstances(
		context.Background(), "nic", "ministry-health", "payments", "prod",
	)
	if err != nil {
		t.Fatalf("BlockedByProjectInstances() error = %v", err)
	}
	if blockedBy != nil {
		t.Fatalf("BlockedByProjectInstances() = %+v, want nil", blockedBy)
	}
}

func TestServiceInstanceProjectBlockerChecker_BlockedByOneInstance(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	createProjectBlockerInstance(t, reg, "postgres-one", "nic", "ministry-health", "payments", "prod")
	blocker := NewServiceInstanceProjectBlockerChecker(reg)

	blockedBy, err := blocker.BlockedByProjectInstances(
		context.Background(), "nic", "ministry-health", "payments", "prod",
	)
	if err != nil {
		t.Fatalf("BlockedByProjectInstances() error = %v", err)
	}
	if len(blockedBy) != 1 {
		t.Fatalf("BlockedByProjectInstances() = %+v, want one blocker", blockedBy)
	}
	if blockedBy[0].Kind != "ServiceInstance" || blockedBy[0].Count != 1 {
		t.Errorf("blocker = %+v, want {Kind:ServiceInstance Count:1}", blockedBy[0])
	}
}

func TestServiceInstanceProjectBlockerChecker_BlockedByMultipleInstances(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	createProjectBlockerInstance(t, reg, "postgres-one", "nic", "ministry-health", "payments", "prod")
	createProjectBlockerInstance(t, reg, "postgres-two", "nic", "ministry-health", "payments", "prod")
	createProjectBlockerInstance(t, reg, "postgres-dev", "nic", "ministry-health", "payments", "dev")
	blocker := NewServiceInstanceProjectBlockerChecker(reg)

	blockedBy, err := blocker.BlockedByProjectInstances(
		context.Background(), "nic", "ministry-health", "payments", "prod",
	)
	if err != nil {
		t.Fatalf("BlockedByProjectInstances() error = %v", err)
	}
	if len(blockedBy) != 1 {
		t.Fatalf("BlockedByProjectInstances() = %+v, want one blocker", blockedBy)
	}
	if blockedBy[0].Kind != "ServiceInstance" || blockedBy[0].Count != 2 {
		t.Errorf("blocker = %+v, want {Kind:ServiceInstance Count:2}", blockedBy[0])
	}
}

func TestServiceInstanceProjectBlockerChecker_DifferentTenantNotCounted(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	createProjectBlockerInstance(t, reg, "postgres-other", "nic", "ministry-health", "billing", "prod")
	blocker := NewServiceInstanceProjectBlockerChecker(reg)

	blockedBy, err := blocker.BlockedByProjectInstances(
		context.Background(), "nic", "ministry-health", "payments", "prod",
	)
	if err != nil {
		t.Fatalf("BlockedByProjectInstances() error = %v", err)
	}
	if blockedBy != nil {
		t.Fatalf("BlockedByProjectInstances() = %+v, want nil", blockedBy)
	}
}

func TestServiceInstanceProjectBlockerChecker_EmptyOUIsolation(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	createProjectBlockerInstance(t, reg, "postgres-with-ou", "nic", "ministry-health", "payments", "prod")
	createProjectBlockerInstance(t, reg, "postgres-empty-ou", "nic", "", "payments", "prod")
	blocker := NewServiceInstanceProjectBlockerChecker(reg)

	blockedWithOU, err := blocker.BlockedByProjectInstances(
		context.Background(), "nic", "ministry-health", "payments", "prod",
	)
	if err != nil {
		t.Fatalf("BlockedByProjectInstances(with OU) error = %v", err)
	}
	if len(blockedWithOU) != 1 || blockedWithOU[0].Count != 1 {
		t.Fatalf("BlockedByProjectInstances(with OU) = %+v, want count 1", blockedWithOU)
	}

	blockedEmptyOU, err := blocker.BlockedByProjectInstances(
		context.Background(), "nic", "", "payments", "prod",
	)
	if err != nil {
		t.Fatalf("BlockedByProjectInstances(empty OU) error = %v", err)
	}
	if len(blockedEmptyOU) != 1 || blockedEmptyOU[0].Count != 1 {
		t.Fatalf("BlockedByProjectInstances(empty OU) = %+v, want count 1", blockedEmptyOU)
	}
}

func createProjectBlockerInstance(
	t *testing.T,
	reg *ServiceInstanceRegistry,
	name, org, ou, tenant, project string,
) {
	t.Helper()
	instance := sampleServiceInstance(
		name, org, ou, tenant, project, "datastore-postgresql", "small",
	)
	if _, err := reg.CreateServiceInstance(context.Background(), instance); err != nil {
		t.Fatalf("CreateServiceInstance(%q) error = %v", name, err)
	}
}
