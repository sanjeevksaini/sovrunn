package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func sampleServicePlan(serviceClassName, name string) resources.ServicePlan {
	return resources.ServicePlan{
		APIVersion: "platform.sovrunn.io/v1alpha1",
		Kind:       resources.ServicePlanKind,
		Metadata: resources.Metadata{
			Name:        name,
			Labels:      map[string]string{"env": "test"},
			Annotations: map[string]string{"note": "x"},
		},
		Spec: resources.ServicePlanSpec{
			ServiceClassName: serviceClassName,
			DisplayName:      "Display",
			Description:      "desc",
			Tier:             resources.TierSmall,
			Lifecycle:        resources.LifecycleActive,
			Parameters:       map[string]string{"region": "us-east"},
			Tags:             []string{"tier", "small"},
		},
		Status: resources.ServicePlanStatus{Phase: resources.PhaseActive, Message: "ok"},
	}
}

func TestCreateServicePlan_Stores(t *testing.T) {
	reg := NewServicePlanRegistry()
	created, err := reg.CreateServicePlan(context.Background(), sampleServicePlan("postgres", "small"))
	if err != nil {
		t.Fatalf("CreateServicePlan() error = %v", err)
	}
	if created.Metadata.Name != "small" || created.Spec.ServiceClassName != "postgres" {
		t.Errorf("got %s/%s, want postgres/small", created.Spec.ServiceClassName, created.Metadata.Name)
	}
	got, err := reg.GetServicePlan(context.Background(), "postgres", "small")
	if err != nil {
		t.Fatalf("GetServicePlan() error = %v", err)
	}
	if got.Spec.Tier != resources.TierSmall {
		t.Errorf("Tier = %q, want Small", got.Spec.Tier)
	}
}

func TestCreateServicePlan_Duplicate(t *testing.T) {
	reg := NewServicePlanRegistry()
	original := sampleServicePlan("postgres", "small")
	original.Spec.Description = "original"
	if _, err := reg.CreateServicePlan(context.Background(), original); err != nil {
		t.Fatalf("first CreateServicePlan() error = %v", err)
	}
	dup := sampleServicePlan("postgres", "small")
	dup.Spec.Description = "changed"
	_, err := reg.CreateServicePlan(context.Background(), dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
	got, err := reg.GetServicePlan(context.Background(), "postgres", "small")
	if err != nil {
		t.Fatalf("GetServicePlan() error = %v", err)
	}
	if got.Spec.Description != "original" {
		t.Errorf("Description = %q, want original (unchanged)", got.Spec.Description)
	}
}

func TestCreateServicePlan_SameNameDifferentClasses(t *testing.T) {
	reg := NewServicePlanRegistry()
	if _, err := reg.CreateServicePlan(context.Background(), sampleServicePlan("postgres", "small")); err != nil {
		t.Fatalf("create postgres/small: %v", err)
	}
	if _, err := reg.CreateServicePlan(context.Background(), sampleServicePlan("redis", "small")); err != nil {
		t.Fatalf("create redis/small: %v", err)
	}
	if _, err := reg.GetServicePlan(context.Background(), "postgres", "small"); err != nil {
		t.Errorf("get postgres/small: %v", err)
	}
	if _, err := reg.GetServicePlan(context.Background(), "redis", "small"); err != nil {
		t.Errorf("get redis/small: %v", err)
	}
}

func TestGetServicePlan_ByCompositeKey(t *testing.T) {
	reg := NewServicePlanRegistry()
	_, _ = reg.CreateServicePlan(context.Background(), sampleServicePlan("postgres", "small"))
	got, err := reg.GetServicePlan(context.Background(), "postgres", "small")
	if err != nil {
		t.Fatalf("GetServicePlan() error = %v", err)
	}
	if got.Spec.Description != "desc" {
		t.Errorf("Description = %q, want desc", got.Spec.Description)
	}
}

func TestGetServicePlan_NotFound(t *testing.T) {
	reg := NewServicePlanRegistry()
	_, err := reg.GetServicePlan(context.Background(), "postgres", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestListServicePlans_Empty(t *testing.T) {
	reg := NewServicePlanRegistry()
	items, err := reg.ListServicePlans(context.Background())
	if err != nil {
		t.Fatalf("ListServicePlans() error = %v", err)
	}
	if items == nil {
		t.Fatal("got nil slice, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Fatalf("got %d items, want 0", len(items))
	}
}

func TestListServicePlans_Sorted(t *testing.T) {
	reg := NewServicePlanRegistry()
	inputs := []struct{ class, name string }{
		{"zebra", "beta"},
		{"alpha", "delta"},
		{"alpha", "charlie"},
		{"alpha", "bravo"},
		{"mongo", "small"},
	}
	for _, in := range inputs {
		if _, err := reg.CreateServicePlan(context.Background(), sampleServicePlan(in.class, in.name)); err != nil {
			t.Fatalf("create %s/%s: %v", in.class, in.name, err)
		}
	}
	items, err := reg.ListServicePlans(context.Background())
	if err != nil {
		t.Fatalf("ListServicePlans() error = %v", err)
	}
	if len(items) != 5 {
		t.Fatalf("got %d items, want 5", len(items))
	}
	for i := 1; i < len(items); i++ {
		prev, cur := items[i-1], items[i]
		if prev.Spec.ServiceClassName > cur.Spec.ServiceClassName {
			t.Fatalf("not sorted by serviceClassName: %v", items)
		}
		if prev.Spec.ServiceClassName == cur.Spec.ServiceClassName &&
			prev.Metadata.Name >= cur.Metadata.Name {
			t.Fatalf("not sorted by name within class: %v", items)
		}
	}
	first := items[0]
	if first.Spec.ServiceClassName != "alpha" || first.Metadata.Name != "bravo" {
		t.Errorf("first = %s/%s, want alpha/bravo", first.Spec.ServiceClassName, first.Metadata.Name)
	}
}

func TestUpdateServicePlan_MutableFields(t *testing.T) {
	reg := NewServicePlanRegistry()
	_, _ = reg.CreateServicePlan(context.Background(), sampleServicePlan("postgres", "small"))
	update := sampleServicePlan("postgres", "small")
	update.Metadata.Labels = map[string]string{"tier": "gold"}
	update.Metadata.Annotations = map[string]string{"reviewed": "yes"}
	update.Spec.DisplayName = "New Display"
	update.Spec.Description = "new desc"
	update.Spec.Tier = resources.TierLarge
	update.Spec.Lifecycle = resources.LifecycleDeprecated
	update.Spec.Parameters = map[string]string{"region": "eu-west"}
	update.Spec.Tags = []string{"large"}
	got, err := reg.UpdateServicePlan(context.Background(), update)
	if err != nil {
		t.Fatalf("UpdateServicePlan() error = %v", err)
	}
	if got.Spec.DisplayName != "New Display" || got.Spec.Description != "new desc" {
		t.Errorf("DisplayName/Description = %q/%q", got.Spec.DisplayName, got.Spec.Description)
	}
	if got.Spec.Tier != resources.TierLarge || got.Spec.Lifecycle != resources.LifecycleDeprecated {
		t.Errorf("Tier/Lifecycle = %q/%q", got.Spec.Tier, got.Spec.Lifecycle)
	}
	if got.Spec.Parameters["region"] != "eu-west" {
		t.Errorf("Parameters = %v, want region=eu-west", got.Spec.Parameters)
	}
	if len(got.Spec.Tags) != 1 || got.Spec.Tags[0] != "large" {
		t.Errorf("Tags = %v, want [large]", got.Spec.Tags)
	}
}

func TestUpdateServicePlan_PreservesImmutableFields(t *testing.T) {
	reg := NewServicePlanRegistry()
	original := sampleServicePlan("postgres", "small")
	original.Status = resources.ServicePlanStatus{Phase: resources.PhaseActive, Message: "ok"}
	_, _ = reg.CreateServicePlan(context.Background(), original)

	update := sampleServicePlan("postgres", "small")
	update.APIVersion = "tampered/v0"
	update.Kind = "Tampered"
	update.Status = resources.ServicePlanStatus{Phase: resources.PhaseFailed, Message: "hacked"}
	update.Spec.Description = "changed"

	got, err := reg.UpdateServicePlan(context.Background(), update)
	if err != nil {
		t.Fatalf("UpdateServicePlan() error = %v", err)
	}
	if got.Metadata.Name != "small" {
		t.Errorf("Metadata.Name = %q, want small", got.Metadata.Name)
	}
	if got.Spec.ServiceClassName != "postgres" {
		t.Errorf("Spec.ServiceClassName = %q, want postgres", got.Spec.ServiceClassName)
	}
	if got.APIVersion != "platform.sovrunn.io/v1alpha1" {
		t.Errorf("APIVersion = %q, want preserved", got.APIVersion)
	}
	if got.Kind != resources.ServicePlanKind {
		t.Errorf("Kind = %q, want ServicePlan", got.Kind)
	}
	if got.Status.Phase != resources.PhaseActive || got.Status.Message != "ok" {
		t.Errorf("Status = %+v, want {Active ok}", got.Status)
	}
	if got.Spec.Description != "changed" {
		t.Errorf("Description = %q, want changed", got.Spec.Description)
	}
}

func TestUpdateServicePlan_NotFound(t *testing.T) {
	reg := NewServicePlanRegistry()
	_, err := reg.UpdateServicePlan(context.Background(), sampleServicePlan("postgres", "missing"))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestDeleteServicePlan_Exists(t *testing.T) {
	reg := NewServicePlanRegistry()
	_, _ = reg.CreateServicePlan(context.Background(), sampleServicePlan("postgres", "small"))
	if err := reg.DeleteServicePlan(context.Background(), "postgres", "small"); err != nil {
		t.Fatalf("DeleteServicePlan() error = %v", err)
	}
	_, err := reg.GetServicePlan(context.Background(), "postgres", "small")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound after delete", err)
	}
}

func TestDeleteServicePlan_NotFound(t *testing.T) {
	reg := NewServicePlanRegistry()
	err := reg.DeleteServicePlan(context.Background(), "postgres", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestCountByServiceClass(t *testing.T) {
	reg := NewServicePlanRegistry()
	_, _ = reg.CreateServicePlan(context.Background(), sampleServicePlan("postgres", "small"))
	_, _ = reg.CreateServicePlan(context.Background(), sampleServicePlan("postgres", "large"))
	_, _ = reg.CreateServicePlan(context.Background(), sampleServicePlan("redis", "small"))

	count, err := reg.CountByServiceClass(context.Background(), "postgres")
	if err != nil {
		t.Fatalf("CountByServiceClass() error = %v", err)
	}
	if count != 2 {
		t.Errorf("count for postgres = %d, want 2", count)
	}

	count, err = reg.CountByServiceClass(context.Background(), "redis")
	if err != nil {
		t.Fatalf("CountByServiceClass() error = %v", err)
	}
	if count != 1 {
		t.Errorf("count for redis = %d, want 1", count)
	}

	count, err = reg.CountByServiceClass(context.Background(), "missing")
	if err != nil {
		t.Fatalf("CountByServiceClass() error = %v", err)
	}
	if count != 0 {
		t.Errorf("count for missing = %d, want 0", count)
	}
}

func TestServicePlanRegistry_DeepCopyImmutability(t *testing.T) {
	reg := NewServicePlanRegistry()
	ctx := context.Background()

	created, err := reg.CreateServicePlan(ctx, sampleServicePlan("postgres", "small"))
	if err != nil {
		t.Fatalf("CreateServicePlan() error = %v", err)
	}
	created.Metadata.Labels["env"] = "mutated"
	created.Metadata.Annotations["note"] = "mutated"
	created.Spec.Parameters["region"] = "mutated"
	created.Spec.Tags[0] = "mutated"

	got, err := reg.GetServicePlan(ctx, "postgres", "small")
	if err != nil {
		t.Fatalf("GetServicePlan() error = %v", err)
	}
	if got.Metadata.Labels["env"] != "test" || got.Metadata.Annotations["note"] != "x" ||
		got.Spec.Parameters["region"] != "us-east" || got.Spec.Tags[0] != "tier" {
		t.Errorf("Create return shares mutable state with store: labels=%v annotations=%v params=%v tags=%v",
			got.Metadata.Labels, got.Metadata.Annotations, got.Spec.Parameters, got.Spec.Tags)
	}

	got.Metadata.Labels["env"] = "mutated-again"
	got.Spec.Parameters["region"] = "mutated-again"
	got.Spec.Tags[0] = "mutated-again"

	items, err := reg.ListServicePlans(ctx)
	if err != nil {
		t.Fatalf("ListServicePlans() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("ListServicePlans() len = %d, want 1", len(items))
	}
	if items[0].Metadata.Labels["env"] != "test" || items[0].Spec.Parameters["region"] != "us-east" || items[0].Spec.Tags[0] != "tier" {
		t.Errorf("Get return shares mutable state with store: labels=%v params=%v tags=%v",
			items[0].Metadata.Labels, items[0].Spec.Parameters, items[0].Spec.Tags)
	}

	items[0].Metadata.Labels["env"] = "list-mutated"
	items[0].Spec.Parameters["region"] = "list-mutated"
	items[0].Spec.Tags[0] = "list-mutated"

	update := sampleServicePlan("postgres", "small")
	update.Spec.Description = "updated"
	updated, err := reg.UpdateServicePlan(ctx, update)
	if err != nil {
		t.Fatalf("UpdateServicePlan() error = %v", err)
	}
	updated.Metadata.Labels["env"] = "update-mutated"
	updated.Spec.Parameters["region"] = "update-mutated"
	updated.Spec.Tags[0] = "update-mutated"

	after, err := reg.GetServicePlan(ctx, "postgres", "small")
	if err != nil {
		t.Fatalf("GetServicePlan() error = %v", err)
	}
	if after.Metadata.Labels["env"] != "test" || after.Metadata.Annotations["note"] != "x" ||
		after.Spec.Parameters["region"] != "us-east" || after.Spec.Tags[0] != "tier" {
		t.Errorf("Update/List return share mutable state with store: labels=%v annotations=%v params=%v tags=%v",
			after.Metadata.Labels, after.Metadata.Annotations, after.Spec.Parameters, after.Spec.Tags)
	}
}
