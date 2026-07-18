package registry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func validCapability(name, pluginRef, serviceClassRef, desc string) resources.Capability {
	return resources.Capability{
		APIVersion: resources.CapabilityAPIVersion,
		Kind:       resources.CapabilityKind,
		Metadata: resources.Metadata{
			Name:        name,
			Labels:      map[string]string{"k": "v"},
			Annotations: map[string]string{"a": "b"},
		},
		Spec: resources.CapabilitySpec{
			PluginRef:       pluginRef,
			ServiceClassRef: serviceClassRef,
			Operation:       resources.CapOpProvision,
			Supported:       true,
			Description:     desc,
		},
		Status: resources.CapabilityStatus{Phase: resources.PhaseActive},
	}
}

// Feature: plugin-capability-registry, Property 8: Create/Get round trip preserves data
func TestProperty_CapabilityRegistry_CreateGetRoundTrip(t *testing.T) {
	f := func(name, pluginRef, serviceClassRef, desc string) bool {
		if !isValidName(name) || !isValidName(pluginRef) || !isValidName(serviceClassRef) {
			return true
		}
		reg := NewCapabilityRegistry()
		c := validCapability(name, pluginRef, serviceClassRef, desc)
		if _, err := reg.CreateCapability(context.Background(), c); err != nil {
			return false
		}
		got, err := reg.GetCapability(context.Background(), name)
		if err != nil {
			return false
		}
		return got.Metadata.Name == name &&
			got.Spec.PluginRef == pluginRef &&
			got.Spec.ServiceClassRef == serviceClassRef &&
			got.Spec.Description == desc &&
			got.Spec.Operation == resources.CapOpProvision &&
			got.Spec.Supported &&
			got.Metadata.Labels["k"] == "v" &&
			got.Metadata.Annotations["a"] == "b" &&
			got.Kind == resources.CapabilityKind &&
			got.Status.Phase == resources.PhaseActive
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: plugin-capability-registry, Property 9: List sort invariant
func TestProperty_CapabilityRegistry_ListSortedOrder(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%20) + 2
		reg := NewCapabilityRegistry()
		for i := 0; i < n; i++ {
			name := fmt.Sprintf("cap-%02d", (n-1-i)%n)
			if _, err := reg.CreateCapability(context.Background(), validCapability(name, "plugin-a", "class-a", "desc")); err != nil {
				return false
			}
		}
		items, err := reg.ListCapabilities(context.Background(), "", "")
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
func TestProperty_CapabilityRegistry_ReturnsDeepCopies(t *testing.T) {
	f := func(name string) bool {
		if !isValidName(name) {
			return true
		}
		reg := NewCapabilityRegistry()
		if _, err := reg.CreateCapability(context.Background(), validCapability(name, "plugin-a", "class-a", "desc")); err != nil {
			return false
		}
		got, err := reg.GetCapability(context.Background(), name)
		if err != nil {
			return false
		}
		got.Metadata.Labels["k"] = "mutated"
		got.Metadata.Annotations["a"] = "mutated"

		list, err := reg.ListCapabilities(context.Background(), "", "")
		if err != nil {
			return false
		}
		for i := range list {
			list[i].Metadata.Labels["k"] = "list-mutated"
			list[i].Metadata.Annotations["a"] = "list-mutated"
		}

		after, err := reg.GetCapability(context.Background(), name)
		if err != nil {
			return false
		}
		return after.Metadata.Labels["k"] == "v" &&
			after.Metadata.Annotations["a"] == "b"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: plugin-capability-registry, Property 11: Duplicate create never overwrites
func TestProperty_CapabilityRegistry_DuplicateCreateError(t *testing.T) {
	f := func(name string) bool {
		if !isValidName(name) {
			return true
		}
		reg := NewCapabilityRegistry()
		first := validCapability(name, "plugin-a", "class-a", "first-desc")
		if _, err := reg.CreateCapability(context.Background(), first); err != nil {
			return false
		}
		second := validCapability(name, "plugin-b", "class-b", "second-desc")
		_, err := reg.CreateCapability(context.Background(), second)
		if !errors.Is(err, ErrAlreadyExists) {
			return false
		}
		got, err := reg.GetCapability(context.Background(), name)
		if err != nil {
			return false
		}
		return got.Spec.PluginRef == "plugin-a" && got.Spec.Description == "first-desc"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: plugin-capability-registry, Property 12: Capability filter correctness
func TestProperty_CapabilityRegistry_FilterCorrectness(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%10) + 3
		reg := NewCapabilityRegistry()
		for i := 0; i < n; i++ {
			name := fmt.Sprintf("cap-%02d", i)
			pluginRef := fmt.Sprintf("plugin-%d", i%2)
			serviceClassRef := fmt.Sprintf("class-%d", i%3)
			if _, err := reg.CreateCapability(context.Background(), validCapability(name, pluginRef, serviceClassRef, "desc")); err != nil {
				return false
			}
		}

		all, err := reg.ListCapabilities(context.Background(), "", "")
		if err != nil || len(all) != n {
			return false
		}

		pluginFiltered, err := reg.ListCapabilities(context.Background(), "plugin-0", "")
		if err != nil {
			return false
		}
		for _, item := range pluginFiltered {
			if item.Spec.PluginRef != "plugin-0" {
				return false
			}
		}

		classFiltered, err := reg.ListCapabilities(context.Background(), "", "class-1")
		if err != nil {
			return false
		}
		for _, item := range classFiltered {
			if item.Spec.ServiceClassRef != "class-1" {
				return false
			}
		}

		bothFiltered, err := reg.ListCapabilities(context.Background(), "plugin-0", "class-1")
		if err != nil {
			return false
		}
		for _, item := range bothFiltered {
			if item.Spec.PluginRef != "plugin-0" || item.Spec.ServiceClassRef != "class-1" {
				return false
			}
		}

		expectedBoth := 0
		for _, item := range all {
			if item.Spec.PluginRef == "plugin-0" && item.Spec.ServiceClassRef == "class-1" {
				expectedBoth++
			}
		}
		if len(bothFiltered) != expectedBoth {
			return false
		}

		for i := 1; i < len(pluginFiltered); i++ {
			if pluginFiltered[i-1].Metadata.Name >= pluginFiltered[i].Metadata.Name {
				return false
			}
		}
		return true
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
