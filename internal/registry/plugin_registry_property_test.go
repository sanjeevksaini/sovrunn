package registry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func validPlugin(name, version, desc string) resources.Plugin {
	return resources.Plugin{
		APIVersion: resources.PluginAPIVersion,
		Kind:       resources.PluginKind,
		Metadata: resources.Metadata{
			Name:        name,
			Labels:      map[string]string{"k": "v"},
			Annotations: map[string]string{"a": "b"},
		},
		Spec: resources.PluginSpec{
			PluginType:       resources.PluginTypeDStoreOps,
			Version:          version,
			ServiceClassRefs: []string{"datastore.postgresql"},
			DeploymentMode:   resources.DeploymentModeCompiledIn,
			Description:      desc,
			Tags:             []string{"tag"},
		},
		Status: resources.PluginStatus{Phase: resources.PhaseActive},
	}
}

// Feature: plugin-capability-registry, Property 8: Create/Get round trip preserves data
func TestProperty_PluginRegistry_CreateGetRoundTrip(t *testing.T) {
	f := func(name, version, desc string) bool {
		if !isValidName(name) {
			return true
		}
		if version == "" {
			version = "0.1.0"
		}
		reg := NewPluginRegistry()
		p := validPlugin(name, version, desc)
		if _, err := reg.CreatePlugin(context.Background(), p); err != nil {
			return false
		}
		got, err := reg.GetPlugin(context.Background(), name)
		if err != nil {
			return false
		}
		return got.Metadata.Name == name &&
			got.Spec.Version == version &&
			got.Spec.Description == desc &&
			got.Spec.PluginType == resources.PluginTypeDStoreOps &&
			got.Spec.DeploymentMode == resources.DeploymentModeCompiledIn &&
			got.Metadata.Labels["k"] == "v" &&
			got.Metadata.Annotations["a"] == "b" &&
			len(got.Spec.ServiceClassRefs) == 1 && got.Spec.ServiceClassRefs[0] == "datastore.postgresql" &&
			len(got.Spec.Tags) == 1 && got.Spec.Tags[0] == "tag" &&
			got.Kind == resources.PluginKind &&
			got.Status.Phase == resources.PhaseActive
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: plugin-capability-registry, Property 9: List sort invariant
func TestProperty_PluginRegistry_ListSortedOrder(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%20) + 2
		reg := NewPluginRegistry()
		for i := 0; i < n; i++ {
			name := fmt.Sprintf("plugin-%02d", (n-1-i)%n)
			if _, err := reg.CreatePlugin(context.Background(), validPlugin(name, "0.1.0", "desc")); err != nil {
				return false
			}
		}
		items, err := reg.ListPlugins(context.Background())
		if err != nil {
			return false
		}
		for i := 1; i < len(items); i++ {
			if items[i-1].Metadata.Name >= items[i].Metadata.Name {
				return false
			}
		}
		return true
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: plugin-capability-registry, Property 10: Deep-copy immutability
func TestProperty_PluginRegistry_ReturnsDeepCopies(t *testing.T) {
	f := func(name string) bool {
		if !isValidName(name) {
			return true
		}
		reg := NewPluginRegistry()
		if _, err := reg.CreatePlugin(context.Background(), validPlugin(name, "0.1.0", "desc")); err != nil {
			return false
		}
		got, err := reg.GetPlugin(context.Background(), name)
		if err != nil {
			return false
		}
		got.Metadata.Labels["k"] = "mutated"
		got.Metadata.Annotations["a"] = "mutated"
		got.Spec.ServiceClassRefs[0] = "mutated"
		got.Spec.Tags[0] = "mutated"

		list, err := reg.ListPlugins(context.Background())
		if err != nil {
			return false
		}
		for i := range list {
			list[i].Metadata.Labels["k"] = "list-mutated"
			list[i].Spec.ServiceClassRefs[0] = "list-mutated"
			list[i].Spec.Tags[0] = "list-mutated"
		}

		after, err := reg.GetPlugin(context.Background(), name)
		if err != nil {
			return false
		}
		return after.Metadata.Labels["k"] == "v" &&
			after.Metadata.Annotations["a"] == "b" &&
			after.Spec.ServiceClassRefs[0] == "datastore.postgresql" &&
			after.Spec.Tags[0] == "tag"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: plugin-capability-registry, Property 11: Duplicate create never overwrites
func TestProperty_PluginRegistry_DuplicateCreateError(t *testing.T) {
	f := func(name string) bool {
		if !isValidName(name) {
			return true
		}
		reg := NewPluginRegistry()
		first := validPlugin(name, "0.1.0", "first-desc")
		if _, err := reg.CreatePlugin(context.Background(), first); err != nil {
			return false
		}
		second := validPlugin(name, "0.2.0", "second-desc")
		_, err := reg.CreatePlugin(context.Background(), second)
		if !errors.Is(err, ErrAlreadyExists) {
			return false
		}
		got, err := reg.GetPlugin(context.Background(), name)
		if err != nil {
			return false
		}
		return got.Spec.Version == "0.1.0" && got.Spec.Description == "first-desc"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
