package registry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func validServiceClass(name, display, desc string) resources.ServiceClass {
	return resources.ServiceClass{
		APIVersion: "platform.sovrunn.io/v1alpha1",
		Kind:       resources.ServiceClassKind,
		Metadata: resources.Metadata{
			Name:        name,
			Labels:      map[string]string{"k": "v"},
			Annotations: map[string]string{"a": "b"},
		},
		Spec: resources.ServiceClassSpec{
			DisplayName: display,
			Description: desc,
			Category:    resources.CategoryDatabase,
			Lifecycle:   resources.LifecycleActive,
			Tags:        []string{"tag"},
		},
		Status: resources.ServiceClassStatus{Phase: resources.PhaseActive},
	}
}

// Feature: serviceclass-serviceplan, Property 4: ServiceClass Create/Get round trip preserves data
func TestProperty_ServiceClassRegistry_CreateGetRoundTrip(t *testing.T) {
	f := func(name, display, desc string) bool {
		if !isValidName(name) {
			return true
		}
		reg := NewServiceClassRegistry()
		sc := validServiceClass(name, display, desc)
		if _, err := reg.CreateServiceClass(context.Background(), sc); err != nil {
			return false
		}
		got, err := reg.GetServiceClass(context.Background(), name)
		if err != nil {
			return false
		}
		return got.Metadata.Name == name &&
			got.Spec.DisplayName == display &&
			got.Spec.Description == desc &&
			got.Spec.Category == resources.CategoryDatabase &&
			got.Spec.Lifecycle == resources.LifecycleActive &&
			got.Metadata.Labels["k"] == "v" &&
			got.Metadata.Annotations["a"] == "b" &&
			len(got.Spec.Tags) == 1 && got.Spec.Tags[0] == "tag" &&
			got.Kind == resources.ServiceClassKind &&
			got.Status.Phase == resources.PhaseActive
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceclass-serviceplan, Property 6: List ordering is deterministic
func TestProperty_ServiceClassRegistry_ListSortedOrder(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%20) + 2
		reg := NewServiceClassRegistry()
		for i := 0; i < n; i++ {
			name := fmt.Sprintf("class-%02d", (n-1-i)%n)
			if _, err := reg.CreateServiceClass(context.Background(), validServiceClass(name, "d", "desc")); err != nil {
				return false
			}
		}
		items, err := reg.ListServiceClasses(context.Background())
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

// Feature: serviceclass-serviceplan, Property 7: Registries return deep copies
func TestProperty_ServiceClassRegistry_ReturnsDeepCopies(t *testing.T) {
	f := func(name string) bool {
		if !isValidName(name) {
			return true
		}
		reg := NewServiceClassRegistry()
		if _, err := reg.CreateServiceClass(context.Background(), validServiceClass(name, "d", "desc")); err != nil {
			return false
		}
		got, err := reg.GetServiceClass(context.Background(), name)
		if err != nil {
			return false
		}
		got.Metadata.Labels["k"] = "mutated"
		got.Metadata.Annotations["a"] = "mutated"
		got.Spec.Tags[0] = "mutated"

		list, err := reg.ListServiceClasses(context.Background())
		if err != nil {
			return false
		}
		for i := range list {
			list[i].Metadata.Labels["k"] = "list-mutated"
			list[i].Spec.Tags[0] = "list-mutated"
		}

		after, err := reg.GetServiceClass(context.Background(), name)
		if err != nil {
			return false
		}
		return after.Metadata.Labels["k"] == "v" &&
			after.Metadata.Annotations["a"] == "b" &&
			after.Spec.Tags[0] == "tag"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceclass-serviceplan, Property 8: Duplicate create never overwrites
func TestProperty_ServiceClassRegistry_DuplicateCreateError(t *testing.T) {
	f := func(name string) bool {
		if !isValidName(name) {
			return true
		}
		reg := NewServiceClassRegistry()
		first := validServiceClass(name, "first", "first-desc")
		if _, err := reg.CreateServiceClass(context.Background(), first); err != nil {
			return false
		}
		second := validServiceClass(name, "second", "second-desc")
		_, err := reg.CreateServiceClass(context.Background(), second)
		if !errors.Is(err, ErrAlreadyExists) {
			return false
		}
		got, err := reg.GetServiceClass(context.Background(), name)
		if err != nil {
			return false
		}
		return got.Spec.DisplayName == "first" && got.Spec.Description == "first-desc"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
