package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func sampleCapability(name, pluginRef, serviceClassRef string) resources.Capability {
	return resources.Capability{
		APIVersion: resources.CapabilityAPIVersion,
		Kind:       resources.CapabilityKind,
		Metadata: resources.Metadata{
			Name:        name,
			Labels:      map[string]string{"env": "test"},
			Annotations: map[string]string{"note": "x"},
		},
		Spec: resources.CapabilitySpec{
			PluginRef:       pluginRef,
			ServiceClassRef: serviceClassRef,
			Operation:       resources.CapOpProvision,
			Supported:       true,
			Description:     "desc",
		},
		Status: resources.CapabilityStatus{Phase: resources.PhaseActive, Message: "ok"},
	}
}

func TestCreateCapability_Stores(t *testing.T) {
	reg := NewCapabilityRegistry()
	created, err := reg.CreateCapability(context.Background(), sampleCapability("postgres-provision", "postgres-basic", "datastore.postgresql"))
	if err != nil {
		t.Fatalf("CreateCapability() error = %v", err)
	}
	if created.Metadata.Name != "postgres-provision" {
		t.Errorf("Name = %q, want postgres-provision", created.Metadata.Name)
	}
	got, err := reg.GetCapability(context.Background(), "postgres-provision")
	if err != nil {
		t.Fatalf("GetCapability() error = %v", err)
	}
	if got.Metadata.Name != "postgres-provision" || got.Spec.PluginRef != "postgres-basic" {
		t.Errorf("got name=%q pluginRef=%q, want postgres-provision/postgres-basic", got.Metadata.Name, got.Spec.PluginRef)
	}
}

func TestCreateCapability_Duplicate(t *testing.T) {
	reg := NewCapabilityRegistry()
	original := sampleCapability("postgres-provision", "postgres-basic", "datastore.postgresql")
	original.Spec.Description = "original"
	if _, err := reg.CreateCapability(context.Background(), original); err != nil {
		t.Fatalf("first CreateCapability() error = %v", err)
	}
	dup := sampleCapability("postgres-provision", "postgres-basic", "datastore.postgresql")
	dup.Spec.Description = "changed"
	_, err := reg.CreateCapability(context.Background(), dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
	got, err := reg.GetCapability(context.Background(), "postgres-provision")
	if err != nil {
		t.Fatalf("GetCapability() error = %v", err)
	}
	if got.Spec.Description != "original" {
		t.Errorf("Description = %q, want original (unchanged)", got.Spec.Description)
	}
}

func TestGetCapability_ByName(t *testing.T) {
	reg := NewCapabilityRegistry()
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("postgres-provision", "postgres-basic", "datastore.postgresql"))
	got, err := reg.GetCapability(context.Background(), "postgres-provision")
	if err != nil {
		t.Fatalf("GetCapability() error = %v", err)
	}
	if got.Spec.Description != "desc" {
		t.Errorf("Description = %q, want desc", got.Spec.Description)
	}
}

func TestGetCapability_NotFound(t *testing.T) {
	reg := NewCapabilityRegistry()
	_, err := reg.GetCapability(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestListCapabilities_Empty(t *testing.T) {
	reg := NewCapabilityRegistry()
	items, err := reg.ListCapabilities(context.Background(), "", "")
	if err != nil {
		t.Fatalf("ListCapabilities() error = %v", err)
	}
	if items == nil {
		t.Fatal("got nil slice, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Fatalf("got %d items, want 0", len(items))
	}
}

func TestListCapabilities_Sorted(t *testing.T) {
	reg := NewCapabilityRegistry()
	for _, name := range []string{"zebra", "alpha", "mongo", "postgres"} {
		if _, err := reg.CreateCapability(context.Background(), sampleCapability(name, "plugin-a", "class-a")); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
	}
	items, err := reg.ListCapabilities(context.Background(), "", "")
	if err != nil {
		t.Fatalf("ListCapabilities() error = %v", err)
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

func TestListCapabilities_FilterByPluginRef(t *testing.T) {
	reg := NewCapabilityRegistry()
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-a", "plugin-a", "class-a"))
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-b", "plugin-a", "class-b"))
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-c", "plugin-b", "class-a"))

	items, err := reg.ListCapabilities(context.Background(), "plugin-a", "")
	if err != nil {
		t.Fatalf("ListCapabilities() error = %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
	if items[0].Metadata.Name != "cap-a" || items[1].Metadata.Name != "cap-b" {
		t.Errorf("got %q/%q, want cap-a/cap-b", items[0].Metadata.Name, items[1].Metadata.Name)
	}
}

func TestListCapabilities_FilterByServiceClassRef(t *testing.T) {
	reg := NewCapabilityRegistry()
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-a", "plugin-a", "class-a"))
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-b", "plugin-a", "class-b"))
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-c", "plugin-b", "class-a"))

	items, err := reg.ListCapabilities(context.Background(), "", "class-a")
	if err != nil {
		t.Fatalf("ListCapabilities() error = %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
	if items[0].Metadata.Name != "cap-a" || items[1].Metadata.Name != "cap-c" {
		t.Errorf("got %q/%q, want cap-a/cap-c", items[0].Metadata.Name, items[1].Metadata.Name)
	}
}

func TestListCapabilities_FilterByBothAND(t *testing.T) {
	reg := NewCapabilityRegistry()
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-a", "plugin-a", "class-a"))
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-b", "plugin-a", "class-b"))
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-c", "plugin-b", "class-a"))

	items, err := reg.ListCapabilities(context.Background(), "plugin-a", "class-a")
	if err != nil {
		t.Fatalf("ListCapabilities() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("got %d items, want 1", len(items))
	}
	if items[0].Metadata.Name != "cap-a" {
		t.Errorf("got %q, want cap-a", items[0].Metadata.Name)
	}
}

func TestListCapabilities_NoFilters(t *testing.T) {
	reg := NewCapabilityRegistry()
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-a", "plugin-a", "class-a"))
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-b", "plugin-b", "class-b"))

	items, err := reg.ListCapabilities(context.Background(), "", "")
	if err != nil {
		t.Fatalf("ListCapabilities() error = %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
}

func TestDeleteCapability_Exists(t *testing.T) {
	reg := NewCapabilityRegistry()
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("postgres-provision", "postgres-basic", "datastore.postgresql"))
	if err := reg.DeleteCapability(context.Background(), "postgres-provision"); err != nil {
		t.Fatalf("DeleteCapability() error = %v", err)
	}
	_, err := reg.GetCapability(context.Background(), "postgres-provision")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound after delete", err)
	}
}

func TestDeleteCapability_NotFound(t *testing.T) {
	reg := NewCapabilityRegistry()
	err := reg.DeleteCapability(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestCountByPlugin(t *testing.T) {
	reg := NewCapabilityRegistry()
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-a", "postgres-basic", "datastore.postgresql"))
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-b", "postgres-basic", "datastore.mysql"))
	_, _ = reg.CreateCapability(context.Background(), sampleCapability("cap-c", "redis-basic", "cache.redis"))

	count, err := reg.CountByPlugin(context.Background(), "postgres-basic")
	if err != nil {
		t.Fatalf("CountByPlugin() error = %v", err)
	}
	if count != 2 {
		t.Errorf("count for postgres-basic = %d, want 2", count)
	}

	count, err = reg.CountByPlugin(context.Background(), "redis-basic")
	if err != nil {
		t.Fatalf("CountByPlugin() error = %v", err)
	}
	if count != 1 {
		t.Errorf("count for redis-basic = %d, want 1", count)
	}

	count, err = reg.CountByPlugin(context.Background(), "missing")
	if err != nil {
		t.Fatalf("CountByPlugin() error = %v", err)
	}
	if count != 0 {
		t.Errorf("count for missing = %d, want 0", count)
	}
}

func TestCapabilityRegistry_DeepCopyImmutability(t *testing.T) {
	reg := NewCapabilityRegistry()
	ctx := context.Background()

	created, err := reg.CreateCapability(ctx, sampleCapability("postgres-provision", "postgres-basic", "datastore.postgresql"))
	if err != nil {
		t.Fatalf("CreateCapability() error = %v", err)
	}
	created.Metadata.Labels["env"] = "mutated"
	created.Metadata.Annotations["note"] = "mutated"

	got, err := reg.GetCapability(ctx, "postgres-provision")
	if err != nil {
		t.Fatalf("GetCapability() error = %v", err)
	}
	if got.Metadata.Labels["env"] != "test" || got.Metadata.Annotations["note"] != "x" {
		t.Errorf("Create return shares mutable state with store: labels=%v annotations=%v",
			got.Metadata.Labels, got.Metadata.Annotations)
	}

	got.Metadata.Labels["env"] = "mutated-again"
	got.Metadata.Annotations["note"] = "mutated-again"

	items, err := reg.ListCapabilities(ctx, "", "")
	if err != nil {
		t.Fatalf("ListCapabilities() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("ListCapabilities() len = %d, want 1", len(items))
	}
	if items[0].Metadata.Labels["env"] != "test" || items[0].Metadata.Annotations["note"] != "x" {
		t.Errorf("Get return shares mutable state with store: labels=%v annotations=%v",
			items[0].Metadata.Labels, items[0].Metadata.Annotations)
	}

	items[0].Metadata.Labels["env"] = "list-mutated"
	items[0].Metadata.Annotations["note"] = "list-mutated"

	after, err := reg.GetCapability(ctx, "postgres-provision")
	if err != nil {
		t.Fatalf("GetCapability() error = %v", err)
	}
	if after.Metadata.Labels["env"] != "test" || after.Metadata.Annotations["note"] != "x" {
		t.Errorf("List return shares mutable state with store: labels=%v annotations=%v",
			after.Metadata.Labels, after.Metadata.Annotations)
	}
}
