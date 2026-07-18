package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func samplePlugin(name string) resources.Plugin {
	return resources.Plugin{
		APIVersion: resources.PluginAPIVersion,
		Kind:       resources.PluginKind,
		Metadata: resources.Metadata{
			Name:        name,
			Labels:      map[string]string{"env": "test"},
			Annotations: map[string]string{"note": "x"},
		},
		Spec: resources.PluginSpec{
			PluginType:       resources.PluginTypeDStoreOps,
			Version:          "0.1.0",
			ServiceClassRefs: []string{"datastore.postgresql", "datastore.mysql"},
			DeploymentMode:   resources.DeploymentModeCompiledIn,
			Description:      "desc",
			Tags:             []string{"db", "sql"},
		},
		Status: resources.PluginStatus{Phase: resources.PhaseActive, Message: "ok"},
	}
}

func TestCreatePlugin_Stores(t *testing.T) {
	reg := NewPluginRegistry()
	created, err := reg.CreatePlugin(context.Background(), samplePlugin("postgres-basic"))
	if err != nil {
		t.Fatalf("CreatePlugin() error = %v", err)
	}
	if created.Metadata.Name != "postgres-basic" {
		t.Errorf("Name = %q, want postgres-basic", created.Metadata.Name)
	}
	got, err := reg.GetPlugin(context.Background(), "postgres-basic")
	if err != nil {
		t.Fatalf("GetPlugin() error = %v", err)
	}
	if got.Metadata.Name != "postgres-basic" || got.Spec.PluginType != resources.PluginTypeDStoreOps {
		t.Errorf("got name=%q pluginType=%q, want postgres-basic/dStoreOps", got.Metadata.Name, got.Spec.PluginType)
	}
}

func TestCreatePlugin_Duplicate(t *testing.T) {
	reg := NewPluginRegistry()
	original := samplePlugin("postgres-basic")
	original.Spec.Description = "original"
	if _, err := reg.CreatePlugin(context.Background(), original); err != nil {
		t.Fatalf("first CreatePlugin() error = %v", err)
	}
	dup := samplePlugin("postgres-basic")
	dup.Spec.Description = "changed"
	_, err := reg.CreatePlugin(context.Background(), dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
	got, err := reg.GetPlugin(context.Background(), "postgres-basic")
	if err != nil {
		t.Fatalf("GetPlugin() error = %v", err)
	}
	if got.Spec.Description != "original" {
		t.Errorf("Description = %q, want original (unchanged)", got.Spec.Description)
	}
}

func TestGetPlugin_ByName(t *testing.T) {
	reg := NewPluginRegistry()
	_, _ = reg.CreatePlugin(context.Background(), samplePlugin("postgres-basic"))
	got, err := reg.GetPlugin(context.Background(), "postgres-basic")
	if err != nil {
		t.Fatalf("GetPlugin() error = %v", err)
	}
	if got.Spec.Description != "desc" {
		t.Errorf("Description = %q, want desc", got.Spec.Description)
	}
}

func TestGetPlugin_NotFound(t *testing.T) {
	reg := NewPluginRegistry()
	_, err := reg.GetPlugin(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestListPlugins_Empty(t *testing.T) {
	reg := NewPluginRegistry()
	items, err := reg.ListPlugins(context.Background())
	if err != nil {
		t.Fatalf("ListPlugins() error = %v", err)
	}
	if items == nil {
		t.Fatal("got nil slice, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Fatalf("got %d items, want 0", len(items))
	}
}

func TestListPlugins_Sorted(t *testing.T) {
	reg := NewPluginRegistry()
	for _, name := range []string{"zebra", "alpha", "mongo", "postgres"} {
		if _, err := reg.CreatePlugin(context.Background(), samplePlugin(name)); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
	}
	items, err := reg.ListPlugins(context.Background())
	if err != nil {
		t.Fatalf("ListPlugins() error = %v", err)
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

func TestUpdatePlugin_MutableFields(t *testing.T) {
	reg := NewPluginRegistry()
	_, _ = reg.CreatePlugin(context.Background(), samplePlugin("postgres-basic"))
	update := samplePlugin("postgres-basic")
	update.Metadata.Labels = map[string]string{"tier": "gold"}
	update.Metadata.Annotations = map[string]string{"reviewed": "yes"}
	update.Spec.PluginType = resources.PluginTypeCacheOps
	update.Spec.Version = "0.2.0"
	update.Spec.ServiceClassRefs = []string{"cache.redis"}
	update.Spec.DeploymentMode = resources.DeploymentModeCompiledIn
	update.Spec.Description = "new desc"
	update.Spec.Tags = []string{"cache"}
	got, err := reg.UpdatePlugin(context.Background(), update)
	if err != nil {
		t.Fatalf("UpdatePlugin() error = %v", err)
	}
	if got.Spec.PluginType != resources.PluginTypeCacheOps || got.Spec.Version != "0.2.0" {
		t.Errorf("PluginType/Version = %q/%q", got.Spec.PluginType, got.Spec.Version)
	}
	if len(got.Spec.ServiceClassRefs) != 1 || got.Spec.ServiceClassRefs[0] != "cache.redis" {
		t.Errorf("ServiceClassRefs = %v, want [cache.redis]", got.Spec.ServiceClassRefs)
	}
	if got.Spec.Description != "new desc" {
		t.Errorf("Description = %q, want new desc", got.Spec.Description)
	}
	if len(got.Spec.Tags) != 1 || got.Spec.Tags[0] != "cache" {
		t.Errorf("Tags = %v, want [cache]", got.Spec.Tags)
	}
	if got.Metadata.Labels["tier"] != "gold" || got.Metadata.Annotations["reviewed"] != "yes" {
		t.Errorf("Labels/Annotations = %v/%v", got.Metadata.Labels, got.Metadata.Annotations)
	}
}

func TestUpdatePlugin_PreservesImmutableFields(t *testing.T) {
	reg := NewPluginRegistry()
	original := samplePlugin("postgres-basic")
	original.Status = resources.PluginStatus{Phase: resources.PhaseActive, Message: "ok"}
	_, _ = reg.CreatePlugin(context.Background(), original)

	update := samplePlugin("postgres-basic")
	update.APIVersion = "tampered/v0"
	update.Kind = "Tampered"
	update.Status = resources.PluginStatus{Phase: resources.PhaseFailed, Message: "hacked"}
	update.Spec.Description = "changed"

	got, err := reg.UpdatePlugin(context.Background(), update)
	if err != nil {
		t.Fatalf("UpdatePlugin() error = %v", err)
	}
	if got.Metadata.Name != "postgres-basic" {
		t.Errorf("Metadata.Name = %q, want postgres-basic", got.Metadata.Name)
	}
	if got.APIVersion != resources.PluginAPIVersion {
		t.Errorf("APIVersion = %q, want preserved", got.APIVersion)
	}
	if got.Kind != resources.PluginKind {
		t.Errorf("Kind = %q, want Plugin", got.Kind)
	}
	if got.Status.Phase != resources.PhaseActive || got.Status.Message != "ok" {
		t.Errorf("Status = %+v, want {Active ok}", got.Status)
	}
	if got.Spec.Description != "changed" {
		t.Errorf("Description = %q, want changed", got.Spec.Description)
	}
}

func TestUpdatePlugin_NotFound(t *testing.T) {
	reg := NewPluginRegistry()
	_, err := reg.UpdatePlugin(context.Background(), samplePlugin("missing"))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestDeletePlugin_Exists(t *testing.T) {
	reg := NewPluginRegistry()
	_, _ = reg.CreatePlugin(context.Background(), samplePlugin("postgres-basic"))
	if err := reg.DeletePlugin(context.Background(), "postgres-basic"); err != nil {
		t.Fatalf("DeletePlugin() error = %v", err)
	}
	_, err := reg.GetPlugin(context.Background(), "postgres-basic")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound after delete", err)
	}
}

func TestDeletePlugin_NotFound(t *testing.T) {
	reg := NewPluginRegistry()
	err := reg.DeletePlugin(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestPluginRegistry_DeepCopyImmutability(t *testing.T) {
	reg := NewPluginRegistry()
	ctx := context.Background()

	created, err := reg.CreatePlugin(ctx, samplePlugin("postgres-basic"))
	if err != nil {
		t.Fatalf("CreatePlugin() error = %v", err)
	}
	created.Metadata.Labels["env"] = "mutated"
	created.Metadata.Annotations["note"] = "mutated"
	created.Spec.ServiceClassRefs[0] = "mutated"
	created.Spec.Tags[0] = "mutated"

	got, err := reg.GetPlugin(ctx, "postgres-basic")
	if err != nil {
		t.Fatalf("GetPlugin() error = %v", err)
	}
	if got.Metadata.Labels["env"] != "test" || got.Metadata.Annotations["note"] != "x" ||
		got.Spec.ServiceClassRefs[0] != "datastore.postgresql" || got.Spec.Tags[0] != "db" {
		t.Errorf("Create return shares mutable state with store: labels=%v annotations=%v refs=%v tags=%v",
			got.Metadata.Labels, got.Metadata.Annotations, got.Spec.ServiceClassRefs, got.Spec.Tags)
	}

	got.Metadata.Labels["env"] = "mutated-again"
	got.Metadata.Annotations["note"] = "mutated-again"
	got.Spec.ServiceClassRefs[0] = "mutated-again"
	got.Spec.Tags[0] = "mutated-again"

	items, err := reg.ListPlugins(ctx)
	if err != nil {
		t.Fatalf("ListPlugins() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("ListPlugins() len = %d, want 1", len(items))
	}
	if items[0].Metadata.Labels["env"] != "test" || items[0].Metadata.Annotations["note"] != "x" ||
		items[0].Spec.ServiceClassRefs[0] != "datastore.postgresql" || items[0].Spec.Tags[0] != "db" {
		t.Errorf("Get return shares mutable state with store: labels=%v annotations=%v refs=%v tags=%v",
			items[0].Metadata.Labels, items[0].Metadata.Annotations, items[0].Spec.ServiceClassRefs, items[0].Spec.Tags)
	}

	items[0].Metadata.Labels["env"] = "list-mutated"
	items[0].Metadata.Annotations["note"] = "list-mutated"
	items[0].Spec.ServiceClassRefs[0] = "list-mutated"
	items[0].Spec.Tags[0] = "list-mutated"

	update := samplePlugin("postgres-basic")
	update.Spec.Description = "updated"
	updated, err := reg.UpdatePlugin(ctx, update)
	if err != nil {
		t.Fatalf("UpdatePlugin() error = %v", err)
	}
	updated.Metadata.Labels["env"] = "update-mutated"
	updated.Spec.ServiceClassRefs[0] = "update-mutated"
	updated.Spec.Tags[0] = "update-mutated"

	after, err := reg.GetPlugin(ctx, "postgres-basic")
	if err != nil {
		t.Fatalf("GetPlugin() error = %v", err)
	}
	if after.Metadata.Labels["env"] != "test" || after.Metadata.Annotations["note"] != "x" ||
		after.Spec.ServiceClassRefs[0] != "datastore.postgresql" || after.Spec.Tags[0] != "db" {
		t.Errorf("Update/List return share mutable state with store: labels=%v annotations=%v refs=%v tags=%v",
			after.Metadata.Labels, after.Metadata.Annotations, after.Spec.ServiceClassRefs, after.Spec.Tags)
	}
}
