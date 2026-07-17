package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func sampleServiceClass(name string) resources.ServiceClass {
	return resources.ServiceClass{
		APIVersion: "platform.sovrunn.io/v1alpha1",
		Kind:       resources.ServiceClassKind,
		Metadata: resources.Metadata{
			Name:        name,
			Labels:      map[string]string{"env": "test"},
			Annotations: map[string]string{"note": "x"},
		},
		Spec: resources.ServiceClassSpec{
			DisplayName:     "Display",
			Description:     "desc",
			Category:        resources.CategoryDatabase,
			Provider:        "sovrunn",
			Lifecycle:       resources.LifecycleActive,
			DefaultPlanName: "small",
			Tags:            []string{"db", "sql"},
		},
		Status: resources.ServiceClassStatus{Phase: resources.PhaseActive, Message: "ok"},
	}
}

func TestCreateServiceClass_Stores(t *testing.T) {
	reg := NewServiceClassRegistry()
	created, err := reg.CreateServiceClass(context.Background(), sampleServiceClass("postgres"))
	if err != nil {
		t.Fatalf("CreateServiceClass() error = %v", err)
	}
	if created.Metadata.Name != "postgres" {
		t.Errorf("Name = %q, want postgres", created.Metadata.Name)
	}
	got, err := reg.GetServiceClass(context.Background(), "postgres")
	if err != nil {
		t.Fatalf("GetServiceClass() error = %v", err)
	}
	if got.Metadata.Name != "postgres" || got.Spec.Category != resources.CategoryDatabase {
		t.Errorf("got name=%q category=%q, want postgres/Database", got.Metadata.Name, got.Spec.Category)
	}
}

func TestCreateServiceClass_Duplicate(t *testing.T) {
	reg := NewServiceClassRegistry()
	original := sampleServiceClass("postgres")
	original.Spec.Description = "original"
	if _, err := reg.CreateServiceClass(context.Background(), original); err != nil {
		t.Fatalf("first CreateServiceClass() error = %v", err)
	}
	dup := sampleServiceClass("postgres")
	dup.Spec.Description = "changed"
	_, err := reg.CreateServiceClass(context.Background(), dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
	got, err := reg.GetServiceClass(context.Background(), "postgres")
	if err != nil {
		t.Fatalf("GetServiceClass() error = %v", err)
	}
	if got.Spec.Description != "original" {
		t.Errorf("Description = %q, want original (unchanged)", got.Spec.Description)
	}
}

func TestGetServiceClass_ByName(t *testing.T) {
	reg := NewServiceClassRegistry()
	_, _ = reg.CreateServiceClass(context.Background(), sampleServiceClass("postgres"))
	got, err := reg.GetServiceClass(context.Background(), "postgres")
	if err != nil {
		t.Fatalf("GetServiceClass() error = %v", err)
	}
	if got.Spec.Description != "desc" {
		t.Errorf("Description = %q, want desc", got.Spec.Description)
	}
}

func TestGetServiceClass_NotFound(t *testing.T) {
	reg := NewServiceClassRegistry()
	_, err := reg.GetServiceClass(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestListServiceClasses_Empty(t *testing.T) {
	reg := NewServiceClassRegistry()
	items, err := reg.ListServiceClasses(context.Background())
	if err != nil {
		t.Fatalf("ListServiceClasses() error = %v", err)
	}
	if items == nil {
		t.Fatal("got nil slice, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Fatalf("got %d items, want 0", len(items))
	}
}

func TestListServiceClasses_Sorted(t *testing.T) {
	reg := NewServiceClassRegistry()
	for _, name := range []string{"zebra", "alpha", "mongo", "postgres"} {
		if _, err := reg.CreateServiceClass(context.Background(), sampleServiceClass(name)); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
	}
	items, err := reg.ListServiceClasses(context.Background())
	if err != nil {
		t.Fatalf("ListServiceClasses() error = %v", err)
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

func TestUpdateServiceClass_MutableFields(t *testing.T) {
	reg := NewServiceClassRegistry()
	_, _ = reg.CreateServiceClass(context.Background(), sampleServiceClass("postgres"))
	update := sampleServiceClass("postgres")
	update.Metadata.Labels = map[string]string{"tier": "gold"}
	update.Metadata.Annotations = map[string]string{"reviewed": "yes"}
	update.Spec.DisplayName = "New Display"
	update.Spec.Description = "new desc"
	update.Spec.Category = resources.CategoryCache
	update.Spec.Provider = "other"
	update.Spec.Lifecycle = resources.LifecycleDeprecated
	update.Spec.DefaultPlanName = "large"
	update.Spec.Tags = []string{"cache"}
	got, err := reg.UpdateServiceClass(context.Background(), update)
	if err != nil {
		t.Fatalf("UpdateServiceClass() error = %v", err)
	}
	if got.Spec.DisplayName != "New Display" || got.Spec.Description != "new desc" {
		t.Errorf("DisplayName/Description = %q/%q", got.Spec.DisplayName, got.Spec.Description)
	}
	if got.Spec.Category != resources.CategoryCache || got.Spec.Provider != "other" {
		t.Errorf("Category/Provider = %q/%q", got.Spec.Category, got.Spec.Provider)
	}
	if got.Spec.Lifecycle != resources.LifecycleDeprecated || got.Spec.DefaultPlanName != "large" {
		t.Errorf("Lifecycle/DefaultPlanName = %q/%q", got.Spec.Lifecycle, got.Spec.DefaultPlanName)
	}
	if len(got.Spec.Tags) != 1 || got.Spec.Tags[0] != "cache" {
		t.Errorf("Tags = %v, want [cache]", got.Spec.Tags)
	}
	if got.Metadata.Labels["tier"] != "gold" || got.Metadata.Annotations["reviewed"] != "yes" {
		t.Errorf("Labels/Annotations = %v/%v", got.Metadata.Labels, got.Metadata.Annotations)
	}
}

func TestUpdateServiceClass_PreservesImmutableFields(t *testing.T) {
	reg := NewServiceClassRegistry()
	original := sampleServiceClass("postgres")
	original.Status = resources.ServiceClassStatus{Phase: resources.PhaseActive, Message: "ok"}
	_, _ = reg.CreateServiceClass(context.Background(), original)

	update := sampleServiceClass("postgres")
	update.APIVersion = "tampered/v0"
	update.Kind = "Tampered"
	update.Status = resources.ServiceClassStatus{Phase: resources.PhaseFailed, Message: "hacked"}
	update.Spec.Description = "changed"

	got, err := reg.UpdateServiceClass(context.Background(), update)
	if err != nil {
		t.Fatalf("UpdateServiceClass() error = %v", err)
	}
	if got.Metadata.Name != "postgres" {
		t.Errorf("Metadata.Name = %q, want postgres", got.Metadata.Name)
	}
	if got.APIVersion != "platform.sovrunn.io/v1alpha1" {
		t.Errorf("APIVersion = %q, want preserved", got.APIVersion)
	}
	if got.Kind != resources.ServiceClassKind {
		t.Errorf("Kind = %q, want ServiceClass", got.Kind)
	}
	if got.Status.Phase != resources.PhaseActive || got.Status.Message != "ok" {
		t.Errorf("Status = %+v, want {Active ok}", got.Status)
	}
	if got.Spec.Description != "changed" {
		t.Errorf("Description = %q, want changed", got.Spec.Description)
	}
}

func TestUpdateServiceClass_NotFound(t *testing.T) {
	reg := NewServiceClassRegistry()
	_, err := reg.UpdateServiceClass(context.Background(), sampleServiceClass("missing"))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestDeleteServiceClass_Exists(t *testing.T) {
	reg := NewServiceClassRegistry()
	_, _ = reg.CreateServiceClass(context.Background(), sampleServiceClass("postgres"))
	if err := reg.DeleteServiceClass(context.Background(), "postgres"); err != nil {
		t.Fatalf("DeleteServiceClass() error = %v", err)
	}
	_, err := reg.GetServiceClass(context.Background(), "postgres")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound after delete", err)
	}
}

func TestDeleteServiceClass_NotFound(t *testing.T) {
	reg := NewServiceClassRegistry()
	err := reg.DeleteServiceClass(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestServiceClassRegistry_DeepCopyImmutability(t *testing.T) {
	reg := NewServiceClassRegistry()
	ctx := context.Background()

	created, err := reg.CreateServiceClass(ctx, sampleServiceClass("postgres"))
	if err != nil {
		t.Fatalf("CreateServiceClass() error = %v", err)
	}
	created.Metadata.Labels["env"] = "mutated"
	created.Metadata.Annotations["note"] = "mutated"
	created.Spec.Tags[0] = "mutated"

	got, err := reg.GetServiceClass(ctx, "postgres")
	if err != nil {
		t.Fatalf("GetServiceClass() error = %v", err)
	}
	if got.Metadata.Labels["env"] != "test" || got.Metadata.Annotations["note"] != "x" || got.Spec.Tags[0] != "db" {
		t.Errorf("Create return shares mutable state with store: labels=%v annotations=%v tags=%v",
			got.Metadata.Labels, got.Metadata.Annotations, got.Spec.Tags)
	}

	got.Metadata.Labels["env"] = "mutated-again"
	got.Metadata.Annotations["note"] = "mutated-again"
	got.Spec.Tags[0] = "mutated-again"

	items, err := reg.ListServiceClasses(ctx)
	if err != nil {
		t.Fatalf("ListServiceClasses() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("ListServiceClasses() len = %d, want 1", len(items))
	}
	if items[0].Metadata.Labels["env"] != "test" || items[0].Metadata.Annotations["note"] != "x" || items[0].Spec.Tags[0] != "db" {
		t.Errorf("Get return shares mutable state with store: labels=%v annotations=%v tags=%v",
			items[0].Metadata.Labels, items[0].Metadata.Annotations, items[0].Spec.Tags)
	}

	items[0].Metadata.Labels["env"] = "list-mutated"
	items[0].Metadata.Annotations["note"] = "list-mutated"
	items[0].Spec.Tags[0] = "list-mutated"

	update := sampleServiceClass("postgres")
	update.Spec.Description = "updated"
	updated, err := reg.UpdateServiceClass(ctx, update)
	if err != nil {
		t.Fatalf("UpdateServiceClass() error = %v", err)
	}
	updated.Metadata.Labels["env"] = "update-mutated"
	updated.Spec.Tags[0] = "update-mutated"

	after, err := reg.GetServiceClass(ctx, "postgres")
	if err != nil {
		t.Fatalf("GetServiceClass() error = %v", err)
	}
	if after.Metadata.Labels["env"] != "test" || after.Metadata.Annotations["note"] != "x" || after.Spec.Tags[0] != "db" {
		t.Errorf("Update/List return share mutable state with store: labels=%v annotations=%v tags=%v",
			after.Metadata.Labels, after.Metadata.Annotations, after.Spec.Tags)
	}
}
