package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func sampleServiceBinding(name, instanceRef, consumerKind, consumerName string) resources.ServiceBinding {
	return resources.ServiceBinding{
		APIVersion: resources.ServiceBindingAPIVersion,
		Kind:       resources.ServiceBindingKind,
		Metadata: resources.Metadata{
			Name:        name,
			DisplayName: "Display",
			Labels:      map[string]string{"env": "test"},
			Annotations: map[string]string{"note": "x"},
		},
		Spec: resources.ServiceBindingSpec{
			ServiceInstanceRef: instanceRef,
			ConsumerRef: &resources.ConsumerRef{
				Kind: consumerKind,
				Name: consumerName,
			},
			BindingType: resources.BindingTypeCredentials,
		},
		Status: resources.ServiceBindingStatus{
			Phase:     "Ready",
			Message:   "Registered only; no real provisioning in Phase 1",
			SecretRef: "stub-secret-ref",
		},
	}
}

func TestCreateServiceBinding_Stores(t *testing.T) {
	reg := NewServiceBindingRegistry()
	created, err := reg.CreateServiceBinding(context.Background(),
		sampleServiceBinding("pg-bind", "pg-prod", "Application", "payments-api"))
	if err != nil {
		t.Fatalf("CreateServiceBinding() error = %v", err)
	}
	if created.Metadata.Name != "pg-bind" {
		t.Errorf("Name = %q, want pg-bind", created.Metadata.Name)
	}
	got, err := reg.GetServiceBinding(context.Background(), "pg-bind")
	if err != nil {
		t.Fatalf("GetServiceBinding() error = %v", err)
	}
	if got.Metadata.Name != "pg-bind" ||
		got.Spec.ServiceInstanceRef != "pg-prod" ||
		got.Spec.ConsumerRef == nil ||
		got.Spec.ConsumerRef.Kind != "Application" ||
		got.Spec.ConsumerRef.Name != "payments-api" ||
		got.Spec.BindingType != resources.BindingTypeCredentials {
		t.Errorf("got unexpected resource: %+v", got)
	}
}

func TestCreateServiceBinding_Duplicate(t *testing.T) {
	reg := NewServiceBindingRegistry()
	original := sampleServiceBinding("pg-bind", "pg-prod", "Application", "payments-api")
	original.Metadata.Labels = map[string]string{"env": "original"}
	if _, err := reg.CreateServiceBinding(context.Background(), original); err != nil {
		t.Fatalf("first CreateServiceBinding() error = %v", err)
	}
	dup := sampleServiceBinding("pg-bind", "pg-prod", "Application", "other-api")
	dup.Metadata.Labels = map[string]string{"env": "changed"}
	_, err := reg.CreateServiceBinding(context.Background(), dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
	got, err := reg.GetServiceBinding(context.Background(), "pg-bind")
	if err != nil {
		t.Fatalf("GetServiceBinding() error = %v", err)
	}
	if got.Metadata.Labels["env"] != "original" {
		t.Errorf("Labels[env] = %q, want original (unchanged)", got.Metadata.Labels["env"])
	}
	if got.Spec.ConsumerRef.Name != "payments-api" {
		t.Errorf("ConsumerRef.Name = %q, want payments-api (unchanged)", got.Spec.ConsumerRef.Name)
	}
}

func TestCreateServiceBinding_DuplicateAcrossServiceInstances(t *testing.T) {
	reg := NewServiceBindingRegistry()
	if _, err := reg.CreateServiceBinding(context.Background(),
		sampleServiceBinding("pg-bind", "pg-prod", "Application", "payments-api")); err != nil {
		t.Fatalf("first CreateServiceBinding() error = %v", err)
	}
	_, err := reg.CreateServiceBinding(context.Background(),
		sampleServiceBinding("pg-bind", "redis-prod", "Application", "other-api"))
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists (global name uniqueness)", err)
	}
}

func TestGetServiceBinding_NotFound(t *testing.T) {
	reg := NewServiceBindingRegistry()
	_, err := reg.GetServiceBinding(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestListServiceBindings_Empty(t *testing.T) {
	reg := NewServiceBindingRegistry()
	items, err := reg.ListServiceBindings(context.Background(), "")
	if err != nil {
		t.Fatalf("ListServiceBindings() error = %v", err)
	}
	if items == nil {
		t.Fatal("got nil slice, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Fatalf("got %d items, want 0", len(items))
	}
}

func TestListServiceBindings_Sorted(t *testing.T) {
	reg := NewServiceBindingRegistry()
	for _, name := range []string{"zebra", "alpha", "mongo", "postgres"} {
		if _, err := reg.CreateServiceBinding(context.Background(),
			sampleServiceBinding(name, "pg-prod", "Application", "app-"+name)); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
	}
	items, err := reg.ListServiceBindings(context.Background(), "")
	if err != nil {
		t.Fatalf("ListServiceBindings() error = %v", err)
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

func TestListServiceBindings_FilterByServiceInstanceRef(t *testing.T) {
	reg := NewServiceBindingRegistry()
	_, _ = reg.CreateServiceBinding(context.Background(),
		sampleServiceBinding("sb-a", "pg-prod", "Application", "app-a"))
	_, _ = reg.CreateServiceBinding(context.Background(),
		sampleServiceBinding("sb-b", "pg-prod", "Application", "app-b"))
	_, _ = reg.CreateServiceBinding(context.Background(),
		sampleServiceBinding("sb-c", "redis-prod", "Application", "app-c"))

	items, err := reg.ListServiceBindings(context.Background(), "pg-prod")
	if err != nil {
		t.Fatalf("ListServiceBindings() error = %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
	if items[0].Metadata.Name != "sb-a" || items[1].Metadata.Name != "sb-b" {
		t.Errorf("got %q/%q, want sb-a/sb-b", items[0].Metadata.Name, items[1].Metadata.Name)
	}
}

func TestListServiceBindings_NoFilter(t *testing.T) {
	reg := NewServiceBindingRegistry()
	_, _ = reg.CreateServiceBinding(context.Background(),
		sampleServiceBinding("sb-a", "pg-prod", "Application", "app-a"))
	_, _ = reg.CreateServiceBinding(context.Background(),
		sampleServiceBinding("sb-b", "redis-prod", "Application", "app-b"))

	items, err := reg.ListServiceBindings(context.Background(), "")
	if err != nil {
		t.Fatalf("ListServiceBindings() error = %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
}

func TestDeleteServiceBinding_Exists(t *testing.T) {
	reg := NewServiceBindingRegistry()
	_, _ = reg.CreateServiceBinding(context.Background(),
		sampleServiceBinding("pg-bind", "pg-prod", "Application", "payments-api"))
	if err := reg.DeleteServiceBinding(context.Background(), "pg-bind"); err != nil {
		t.Fatalf("DeleteServiceBinding() error = %v", err)
	}
	_, err := reg.GetServiceBinding(context.Background(), "pg-bind")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound after delete", err)
	}
}

func TestDeleteServiceBinding_NotFound(t *testing.T) {
	reg := NewServiceBindingRegistry()
	err := reg.DeleteServiceBinding(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestCountByServiceInstance(t *testing.T) {
	reg := NewServiceBindingRegistry()
	_, _ = reg.CreateServiceBinding(context.Background(),
		sampleServiceBinding("sb-a", "pg-prod", "Application", "app-a"))
	_, _ = reg.CreateServiceBinding(context.Background(),
		sampleServiceBinding("sb-b", "pg-prod", "Application", "app-b"))
	_, _ = reg.CreateServiceBinding(context.Background(),
		sampleServiceBinding("sb-c", "redis-prod", "Application", "app-c"))

	count, err := reg.CountByServiceInstance(context.Background(), "pg-prod")
	if err != nil {
		t.Fatalf("CountByServiceInstance() error = %v", err)
	}
	if count != 2 {
		t.Errorf("count for pg-prod = %d, want 2", count)
	}

	count, err = reg.CountByServiceInstance(context.Background(), "redis-prod")
	if err != nil {
		t.Fatalf("CountByServiceInstance() error = %v", err)
	}
	if count != 1 {
		t.Errorf("count for redis-prod = %d, want 1", count)
	}

	count, err = reg.CountByServiceInstance(context.Background(), "missing")
	if err != nil {
		t.Fatalf("CountByServiceInstance() error = %v", err)
	}
	if count != 0 {
		t.Errorf("count for missing = %d, want 0", count)
	}
}

func TestServiceBindingRegistry_DeepCopyImmutability(t *testing.T) {
	reg := NewServiceBindingRegistry()
	ctx := context.Background()

	created, err := reg.CreateServiceBinding(ctx,
		sampleServiceBinding("pg-bind", "pg-prod", "Application", "payments-api"))
	if err != nil {
		t.Fatalf("CreateServiceBinding() error = %v", err)
	}
	created.Metadata.Labels["env"] = "mutated"
	created.Metadata.Annotations["note"] = "mutated"
	created.Spec.ConsumerRef.Kind = "Mutated"
	created.Spec.ConsumerRef.Name = "mutated-name"

	got, err := reg.GetServiceBinding(ctx, "pg-bind")
	if err != nil {
		t.Fatalf("GetServiceBinding() error = %v", err)
	}
	if got.Metadata.Labels["env"] != "test" ||
		got.Metadata.Annotations["note"] != "x" ||
		got.Spec.ConsumerRef.Kind != "Application" ||
		got.Spec.ConsumerRef.Name != "payments-api" {
		t.Errorf("Create return shares mutable state with store: labels=%v annotations=%v consumer=%+v",
			got.Metadata.Labels, got.Metadata.Annotations, got.Spec.ConsumerRef)
	}

	got.Metadata.Labels["env"] = "mutated-again"
	got.Metadata.Annotations["note"] = "mutated-again"
	got.Spec.ConsumerRef.Kind = "MutatedAgain"
	got.Spec.ConsumerRef.Name = "mutated-again"

	items, err := reg.ListServiceBindings(ctx, "")
	if err != nil {
		t.Fatalf("ListServiceBindings() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("ListServiceBindings() len = %d, want 1", len(items))
	}
	if items[0].Metadata.Labels["env"] != "test" ||
		items[0].Metadata.Annotations["note"] != "x" ||
		items[0].Spec.ConsumerRef.Kind != "Application" ||
		items[0].Spec.ConsumerRef.Name != "payments-api" {
		t.Errorf("Get return shares mutable state with store: labels=%v annotations=%v consumer=%+v",
			items[0].Metadata.Labels, items[0].Metadata.Annotations, items[0].Spec.ConsumerRef)
	}

	items[0].Metadata.Labels["env"] = "list-mutated"
	items[0].Spec.ConsumerRef.Name = "list-mutated"

	after, err := reg.GetServiceBinding(ctx, "pg-bind")
	if err != nil {
		t.Fatalf("GetServiceBinding() error = %v", err)
	}
	if after.Metadata.Labels["env"] != "test" {
		t.Errorf("Labels[env] = %q, want test", after.Metadata.Labels["env"])
	}
	if after.Spec.ConsumerRef.Name != "payments-api" {
		t.Errorf("ConsumerRef.Name = %q, want payments-api", after.Spec.ConsumerRef.Name)
	}
}
