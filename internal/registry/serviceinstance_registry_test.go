package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func sampleServiceInstance(name, org, ou, tenant, project, class, plan string) resources.ServiceInstance {
	return resources.ServiceInstance{
		APIVersion: resources.ServiceInstanceAPIVersion,
		Kind:       resources.ServiceInstanceKind,
		Metadata: resources.Metadata{
			Name:        name,
			DisplayName: "Display",
			Labels:      map[string]string{"env": "test"},
			Annotations: map[string]string{"note": "x"},
		},
		Spec: resources.ServiceInstanceSpec{
			OrganizationRef:     org,
			OrganizationUnitRef: ou,
			TenantRef:           tenant,
			ProjectRef:          project,
			ServiceClassRef:     class,
			ServicePlanRef:      plan,
			Parameters:          map[string]string{"size": "small"},
		},
		Status: resources.ServiceInstanceStatus{
			Phase:   "Ready",
			Message: "Registered only; no real provisioning in Phase 1",
		},
	}
}

func TestCreateServiceInstance_Stores(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	created, err := reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("pg-prod", "nic", "ministry-health", "payments", "prod", "datastore.postgresql", "small"))
	if err != nil {
		t.Fatalf("CreateServiceInstance() error = %v", err)
	}
	if created.Metadata.Name != "pg-prod" {
		t.Errorf("Name = %q, want pg-prod", created.Metadata.Name)
	}
	got, err := reg.GetServiceInstance(context.Background(), "pg-prod")
	if err != nil {
		t.Fatalf("GetServiceInstance() error = %v", err)
	}
	if got.Metadata.Name != "pg-prod" ||
		got.Spec.OrganizationRef != "nic" ||
		got.Spec.TenantRef != "payments" ||
		got.Spec.ProjectRef != "prod" ||
		got.Spec.ServiceClassRef != "datastore.postgresql" ||
		got.Spec.ServicePlanRef != "small" {
		t.Errorf("got unexpected resource: %+v", got)
	}
}

func TestCreateServiceInstance_Duplicate(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	original := sampleServiceInstance("pg-prod", "nic", "ministry-health", "payments", "prod", "datastore.postgresql", "small")
	original.Spec.Parameters = map[string]string{"size": "original"}
	if _, err := reg.CreateServiceInstance(context.Background(), original); err != nil {
		t.Fatalf("first CreateServiceInstance() error = %v", err)
	}
	dup := sampleServiceInstance("pg-prod", "nic", "ministry-health", "payments", "prod", "datastore.postgresql", "small")
	dup.Spec.Parameters = map[string]string{"size": "changed"}
	_, err := reg.CreateServiceInstance(context.Background(), dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
	got, err := reg.GetServiceInstance(context.Background(), "pg-prod")
	if err != nil {
		t.Fatalf("GetServiceInstance() error = %v", err)
	}
	if got.Spec.Parameters["size"] != "original" {
		t.Errorf("Parameters[size] = %q, want original (unchanged)", got.Spec.Parameters["size"])
	}
}

func TestCreateServiceInstance_DuplicateAcrossGovernanceRefs(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	if _, err := reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("pg-prod", "nic", "ministry-health", "payments", "prod", "datastore.postgresql", "small")); err != nil {
		t.Fatalf("first CreateServiceInstance() error = %v", err)
	}
	_, err := reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("pg-prod", "other-org", "other-ou", "other-tenant", "other-proj", "cache.redis", "basic"))
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists (global name uniqueness)", err)
	}
}

func TestGetServiceInstance_NotFound(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, err := reg.GetServiceInstance(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestListServiceInstances_Empty(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	items, err := reg.ListServiceInstances(context.Background(), "", "")
	if err != nil {
		t.Fatalf("ListServiceInstances() error = %v", err)
	}
	if items == nil {
		t.Fatal("got nil slice, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Fatalf("got %d items, want 0", len(items))
	}
}

func TestListServiceInstances_Sorted(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	for _, name := range []string{"zebra", "alpha", "mongo", "postgres"} {
		if _, err := reg.CreateServiceInstance(context.Background(),
			sampleServiceInstance(name, "nic", "", "payments", "prod", "datastore.postgresql", "small")); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
	}
	items, err := reg.ListServiceInstances(context.Background(), "", "")
	if err != nil {
		t.Fatalf("ListServiceInstances() error = %v", err)
	}
	if len(items) != 4 {
		t.Fatalf("got %d items, want 4", len(items))
	}
	for i := 1; i < len(items); i++ {
		if items[i-1].Metadata.Name >= items[i].Metadata.Name {
			t.Fatalf("not sorted by name: %v before %v", items[i-1].Metadata.Name, items[i].Metadata.Name)
		}
	}
	if items[0].Metadata.Name != "alpha" {
		t.Errorf("first = %q, want alpha", items[0].Metadata.Name)
	}
}

func TestListServiceInstances_FilterByTenantRef(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-a", "nic", "", "payments", "prod", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-b", "nic", "", "payments", "dev", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-c", "nic", "", "billing", "prod", "datastore.postgresql", "small"))

	items, err := reg.ListServiceInstances(context.Background(), "payments", "")
	if err != nil {
		t.Fatalf("ListServiceInstances() error = %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
	if items[0].Metadata.Name != "si-a" || items[1].Metadata.Name != "si-b" {
		t.Errorf("got %q/%q, want si-a/si-b", items[0].Metadata.Name, items[1].Metadata.Name)
	}
}

func TestListServiceInstances_FilterByProjectRef(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-a", "nic", "", "payments", "prod", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-b", "nic", "", "payments", "dev", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-c", "nic", "", "billing", "prod", "datastore.postgresql", "small"))

	items, err := reg.ListServiceInstances(context.Background(), "", "prod")
	if err != nil {
		t.Fatalf("ListServiceInstances() error = %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
	if items[0].Metadata.Name != "si-a" || items[1].Metadata.Name != "si-c" {
		t.Errorf("got %q/%q, want si-a/si-c", items[0].Metadata.Name, items[1].Metadata.Name)
	}
}

func TestListServiceInstances_FilterByBothAND(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-a", "nic", "", "payments", "prod", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-b", "nic", "", "payments", "dev", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-c", "nic", "", "billing", "prod", "datastore.postgresql", "small"))

	items, err := reg.ListServiceInstances(context.Background(), "payments", "prod")
	if err != nil {
		t.Fatalf("ListServiceInstances() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("got %d items, want 1", len(items))
	}
	if items[0].Metadata.Name != "si-a" {
		t.Errorf("got %q, want si-a", items[0].Metadata.Name)
	}
}

func TestListServiceInstances_NoFilters(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-a", "nic", "", "payments", "prod", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-b", "nic", "", "billing", "dev", "cache.redis", "basic"))

	items, err := reg.ListServiceInstances(context.Background(), "", "")
	if err != nil {
		t.Fatalf("ListServiceInstances() error = %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
}

func TestUpdateServiceInstance_MutableFields(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("pg-prod", "nic", "ministry-health", "payments", "prod", "datastore.postgresql", "small"))

	update := sampleServiceInstance("pg-prod", "changed-org", "changed-ou", "changed-tenant", "changed-proj", "changed-class", "changed-plan")
	update.APIVersion = "changed/v1"
	update.Kind = "Changed"
	update.Status = resources.ServiceInstanceStatus{Phase: "Failed", Message: "should not apply"}
	update.Metadata.DisplayName = "New Display"
	update.Metadata.Labels = map[string]string{"tier": "gold"}
	update.Metadata.Annotations = map[string]string{"reviewed": "yes"}
	update.Spec.Parameters = map[string]string{"size": "large"}

	got, err := reg.UpdateServiceInstance(context.Background(), "pg-prod", update)
	if err != nil {
		t.Fatalf("UpdateServiceInstance() error = %v", err)
	}
	if got.Metadata.DisplayName != "New Display" {
		t.Errorf("DisplayName = %q, want New Display", got.Metadata.DisplayName)
	}
	if got.Metadata.Labels["tier"] != "gold" {
		t.Errorf("Labels[tier] = %q, want gold", got.Metadata.Labels["tier"])
	}
	if got.Metadata.Annotations["reviewed"] != "yes" {
		t.Errorf("Annotations[reviewed] = %q, want yes", got.Metadata.Annotations["reviewed"])
	}
	if got.Spec.Parameters["size"] != "large" {
		t.Errorf("Parameters[size] = %q, want large", got.Spec.Parameters["size"])
	}
	// Immutable governance/catalog refs preserved.
	if got.Spec.OrganizationRef != "nic" ||
		got.Spec.OrganizationUnitRef != "ministry-health" ||
		got.Spec.TenantRef != "payments" ||
		got.Spec.ProjectRef != "prod" ||
		got.Spec.ServiceClassRef != "datastore.postgresql" ||
		got.Spec.ServicePlanRef != "small" {
		t.Errorf("immutable refs changed: %+v", got.Spec)
	}
	if got.APIVersion != resources.ServiceInstanceAPIVersion {
		t.Errorf("APIVersion = %q, want preserved", got.APIVersion)
	}
	if got.Kind != resources.ServiceInstanceKind {
		t.Errorf("Kind = %q, want preserved", got.Kind)
	}
	if got.Metadata.Name != "pg-prod" {
		t.Errorf("Metadata.Name = %q, want pg-prod", got.Metadata.Name)
	}
}

func TestUpdateServiceInstance_PreservesStatus(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	original := sampleServiceInstance("pg-prod", "nic", "", "payments", "prod", "datastore.postgresql", "small")
	original.Status = resources.ServiceInstanceStatus{Phase: "Ready", Message: "original message"}
	_, _ = reg.CreateServiceInstance(context.Background(), original)

	update := sampleServiceInstance("pg-prod", "nic", "", "payments", "prod", "datastore.postgresql", "small")
	update.Status = resources.ServiceInstanceStatus{Phase: "Failed", Message: "attacker"}
	update.Metadata.DisplayName = "Updated"
	got, err := reg.UpdateServiceInstance(context.Background(), "pg-prod", update)
	if err != nil {
		t.Fatalf("UpdateServiceInstance() error = %v", err)
	}
	if got.Status.Phase != "Ready" || got.Status.Message != "original message" {
		t.Errorf("Status = %+v, want {Ready original message}", got.Status)
	}
}

func TestUpdateServiceInstance_NotFound(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, err := reg.UpdateServiceInstance(context.Background(), "missing",
		sampleServiceInstance("missing", "nic", "", "payments", "prod", "datastore.postgresql", "small"))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestDeleteServiceInstance_Exists(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("pg-prod", "nic", "", "payments", "prod", "datastore.postgresql", "small"))
	if err := reg.DeleteServiceInstance(context.Background(), "pg-prod"); err != nil {
		t.Fatalf("DeleteServiceInstance() error = %v", err)
	}
	_, err := reg.GetServiceInstance(context.Background(), "pg-prod")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound after delete", err)
	}
}

func TestDeleteServiceInstance_NotFound(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	err := reg.DeleteServiceInstance(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestCountByServicePlan(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-a", "nic", "", "payments", "prod", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-b", "nic", "", "payments", "dev", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-c", "nic", "", "billing", "prod", "datastore.postgresql", "large"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-d", "nic", "", "billing", "dev", "cache.redis", "small"))

	count, err := reg.CountByServicePlan(context.Background(), "datastore.postgresql", "small")
	if err != nil {
		t.Fatalf("CountByServicePlan() error = %v", err)
	}
	if count != 2 {
		t.Errorf("count for datastore.postgresql/small = %d, want 2", count)
	}

	count, err = reg.CountByServicePlan(context.Background(), "datastore.postgresql", "large")
	if err != nil {
		t.Fatalf("CountByServicePlan() error = %v", err)
	}
	if count != 1 {
		t.Errorf("count for datastore.postgresql/large = %d, want 1", count)
	}

	count, err = reg.CountByServicePlan(context.Background(), "missing", "small")
	if err != nil {
		t.Fatalf("CountByServicePlan() error = %v", err)
	}
	if count != 0 {
		t.Errorf("count for missing = %d, want 0", count)
	}
}

func TestCountByServicePlan_NoFalsePositivesAcrossClasses(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-a", "nic", "", "payments", "prod", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-b", "nic", "", "payments", "dev", "cache.redis", "small"))

	count, err := reg.CountByServicePlan(context.Background(), "datastore.postgresql", "small")
	if err != nil {
		t.Fatalf("CountByServicePlan() error = %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1 (same plan name under different ServiceClass must not match)", count)
	}
}

func TestCountByProject(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-a", "nic", "ministry-health", "payments", "prod", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-b", "nic", "ministry-health", "payments", "prod", "cache.redis", "basic"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-c", "nic", "ministry-health", "payments", "dev", "datastore.postgresql", "small"))

	count, err := reg.CountByProject(context.Background(), "nic", "ministry-health", "payments", "prod")
	if err != nil {
		t.Fatalf("CountByProject() error = %v", err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}

	count, err = reg.CountByProject(context.Background(), "nic", "ministry-health", "payments", "dev")
	if err != nil {
		t.Fatalf("CountByProject() error = %v", err)
	}
	if count != 1 {
		t.Errorf("count for prod→dev = %d, want 1", count)
	}
}

func TestCountByProject_NoFalsePositivesAcrossTenants(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-a", "nic", "ministry-health", "payments", "prod", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-b", "nic", "ministry-health", "billing", "prod", "datastore.postgresql", "small"))

	count, err := reg.CountByProject(context.Background(), "nic", "ministry-health", "payments", "prod")
	if err != nil {
		t.Fatalf("CountByProject() error = %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1 (same project name under different tenant must not match)", count)
	}
}

func TestCountByProject_EmptyOUDoesNotMatchNonEmptyOU(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-empty-ou", "nic", "", "payments", "prod", "datastore.postgresql", "small"))
	_, _ = reg.CreateServiceInstance(context.Background(),
		sampleServiceInstance("si-with-ou", "nic", "ministry-health", "payments", "prod", "datastore.postgresql", "small"))

	count, err := reg.CountByProject(context.Background(), "nic", "", "payments", "prod")
	if err != nil {
		t.Fatalf("CountByProject() error = %v", err)
	}
	if count != 1 {
		t.Errorf("empty OU count = %d, want 1", count)
	}

	count, err = reg.CountByProject(context.Background(), "nic", "ministry-health", "payments", "prod")
	if err != nil {
		t.Fatalf("CountByProject() error = %v", err)
	}
	if count != 1 {
		t.Errorf("non-empty OU count = %d, want 1", count)
	}
}

func TestServiceInstanceRegistry_DeepCopyImmutability(t *testing.T) {
	reg := NewServiceInstanceRegistry()
	ctx := context.Background()

	created, err := reg.CreateServiceInstance(ctx,
		sampleServiceInstance("pg-prod", "nic", "", "payments", "prod", "datastore.postgresql", "small"))
	if err != nil {
		t.Fatalf("CreateServiceInstance() error = %v", err)
	}
	created.Metadata.Labels["env"] = "mutated"
	created.Metadata.Annotations["note"] = "mutated"
	created.Spec.Parameters["size"] = "mutated"

	got, err := reg.GetServiceInstance(ctx, "pg-prod")
	if err != nil {
		t.Fatalf("GetServiceInstance() error = %v", err)
	}
	if got.Metadata.Labels["env"] != "test" ||
		got.Metadata.Annotations["note"] != "x" ||
		got.Spec.Parameters["size"] != "small" {
		t.Errorf("Create return shares mutable state with store: labels=%v annotations=%v params=%v",
			got.Metadata.Labels, got.Metadata.Annotations, got.Spec.Parameters)
	}

	got.Metadata.Labels["env"] = "mutated-again"
	got.Metadata.Annotations["note"] = "mutated-again"
	got.Spec.Parameters["size"] = "mutated-again"

	items, err := reg.ListServiceInstances(ctx, "", "")
	if err != nil {
		t.Fatalf("ListServiceInstances() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("ListServiceInstances() len = %d, want 1", len(items))
	}
	if items[0].Metadata.Labels["env"] != "test" ||
		items[0].Metadata.Annotations["note"] != "x" ||
		items[0].Spec.Parameters["size"] != "small" {
		t.Errorf("Get return shares mutable state with store: labels=%v annotations=%v params=%v",
			items[0].Metadata.Labels, items[0].Metadata.Annotations, items[0].Spec.Parameters)
	}

	items[0].Metadata.Labels["env"] = "list-mutated"
	items[0].Spec.Parameters["size"] = "list-mutated"

	update := sampleServiceInstance("pg-prod", "nic", "", "payments", "prod", "datastore.postgresql", "small")
	update.Metadata.DisplayName = "updated"
	update.Spec.Parameters = map[string]string{"size": "medium"}
	updated, err := reg.UpdateServiceInstance(ctx, "pg-prod", update)
	if err != nil {
		t.Fatalf("UpdateServiceInstance() error = %v", err)
	}
	updated.Metadata.Labels["env"] = "update-mutated"
	updated.Spec.Parameters["size"] = "update-mutated"

	after, err := reg.GetServiceInstance(ctx, "pg-prod")
	if err != nil {
		t.Fatalf("GetServiceInstance() error = %v", err)
	}
	if after.Spec.Parameters["size"] != "medium" {
		t.Errorf("Parameters[size] = %q, want medium", after.Spec.Parameters["size"])
	}
	if after.Metadata.Labels["env"] != "test" {
		t.Errorf("Labels[env] = %q, want test (update returned mutable copy)", after.Metadata.Labels["env"])
	}
}
