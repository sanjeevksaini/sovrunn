package registry

import (
	"context"
	"testing"
)

func TestServiceInstancePlanBlockerChecker_NotBlockedWithoutInstances(t *testing.T) {
	blocker := NewServiceInstancePlanBlockerChecker(NewServiceInstanceRegistry())

	blockedBy, err := blocker.BlockedByServicePlanInstances(
		context.Background(), "datastore-postgresql", "small",
	)
	if err != nil {
		t.Fatalf("BlockedByServicePlanInstances() error = %v", err)
	}
	if blockedBy != nil {
		t.Fatalf("BlockedByServicePlanInstances() = %+v, want nil", blockedBy)
	}
}

func TestServiceInstancePlanBlockerChecker_BlockedByOneInstance(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	createServicePlanBlockerInstance(t, reg, "postgres-one", "datastore-postgresql", "small")
	blocker := NewServiceInstancePlanBlockerChecker(reg)

	blockedBy, err := blocker.BlockedByServicePlanInstances(
		context.Background(), "datastore-postgresql", "small",
	)
	if err != nil {
		t.Fatalf("BlockedByServicePlanInstances() error = %v", err)
	}
	if len(blockedBy) != 1 {
		t.Fatalf("BlockedByServicePlanInstances() = %+v, want one blocker", blockedBy)
	}
	if blockedBy[0].Kind != "ServiceInstance" || blockedBy[0].Count != 1 {
		t.Errorf("blocker = %+v, want {Kind:ServiceInstance Count:1}", blockedBy[0])
	}
}

func TestServiceInstancePlanBlockerChecker_BlockedByMultipleInstances(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	createServicePlanBlockerInstance(t, reg, "postgres-one", "datastore-postgresql", "small")
	createServicePlanBlockerInstance(t, reg, "postgres-two", "datastore-postgresql", "small")
	createServicePlanBlockerInstance(t, reg, "postgres-large", "datastore-postgresql", "large")
	blocker := NewServiceInstancePlanBlockerChecker(reg)

	blockedBy, err := blocker.BlockedByServicePlanInstances(
		context.Background(), "datastore-postgresql", "small",
	)
	if err != nil {
		t.Fatalf("BlockedByServicePlanInstances() error = %v", err)
	}
	if len(blockedBy) != 1 {
		t.Fatalf("BlockedByServicePlanInstances() = %+v, want one blocker", blockedBy)
	}
	if blockedBy[0].Kind != "ServiceInstance" || blockedBy[0].Count != 2 {
		t.Errorf("blocker = %+v, want {Kind:ServiceInstance Count:2}", blockedBy[0])
	}
}

func TestServiceInstancePlanBlockerChecker_DifferentServiceClassNotCounted(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	createServicePlanBlockerInstance(t, reg, "redis-small", "cache-redis", "small")
	blocker := NewServiceInstancePlanBlockerChecker(reg)

	blockedBy, err := blocker.BlockedByServicePlanInstances(
		context.Background(), "datastore-postgresql", "small",
	)
	if err != nil {
		t.Fatalf("BlockedByServicePlanInstances() error = %v", err)
	}
	if blockedBy != nil {
		t.Fatalf("BlockedByServicePlanInstances() = %+v, want nil", blockedBy)
	}
}

func createServicePlanBlockerInstance(
	t *testing.T,
	reg *ServiceInstanceRegistry,
	name, serviceClass, plan string,
) {
	t.Helper()
	instance := sampleServiceInstance(
		name, "nic", "", "payments", "prod", serviceClass, plan,
	)
	if _, err := reg.CreateServiceInstance(context.Background(), instance); err != nil {
		t.Fatalf("CreateServiceInstance(%q) error = %v", name, err)
	}
}
