package registry

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func propertyServiceBinding(name, instanceRef, displayName, label, annotation, consumerName string) resources.ServiceBinding {
	return resources.ServiceBinding{
		APIVersion: resources.ServiceBindingAPIVersion,
		Kind:       resources.ServiceBindingKind,
		Metadata: resources.Metadata{
			Name:        name,
			DisplayName: displayName,
			Labels:      map[string]string{"label": label},
			Annotations: map[string]string{"annotation": annotation},
		},
		Spec: resources.ServiceBindingSpec{
			ServiceInstanceRef: instanceRef,
			ConsumerRef: &resources.ConsumerRef{
				Kind: "Application",
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

// Feature: serviceinstance-servicebinding, Property 13: Create/Get round-trip preserves data
func TestProperty_ServiceBindingRegistry_CreateGetRoundTrip(t *testing.T) {
	f := func(seed uint16, displayName, label, annotation, consumerName string) bool {
		name := fmt.Sprintf("binding-%d", seed)
		instanceRef := fmt.Sprintf("instance-%d", seed%7)
		sb := propertyServiceBinding(name, instanceRef, displayName, label, annotation, consumerName)
		reg := NewServiceBindingRegistry()

		created, err := reg.CreateServiceBinding(context.Background(), sb)
		if err != nil || !reflect.DeepEqual(created, sb) {
			return false
		}
		got, err := reg.GetServiceBinding(context.Background(), name)
		return err == nil && reflect.DeepEqual(got, sb)
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceinstance-servicebinding, Property 14: List sort invariant
func TestProperty_ServiceBindingRegistry_ListSortedOrder(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%20) + 2
		reg := NewServiceBindingRegistry()
		for i := n - 1; i >= 0; i-- {
			name := fmt.Sprintf("binding-%02d", i)
			sb := propertyServiceBinding(name, "instance-a", "display", "value", "note", "consumer")
			if _, err := reg.CreateServiceBinding(context.Background(), sb); err != nil {
				return false
			}
		}

		items, err := reg.ListServiceBindings(context.Background(), "")
		if err != nil || len(items) != n {
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

// Feature: serviceinstance-servicebinding, Property 15: Deep-copy immutability
func TestProperty_ServiceBindingRegistry_ReturnsDeepCopies(t *testing.T) {
	f := func(seed uint16) bool {
		name := fmt.Sprintf("binding-%d", seed)
		reg := NewServiceBindingRegistry()
		created, err := reg.CreateServiceBinding(context.Background(),
			propertyServiceBinding(name, "instance-a", "display", "original", "original", "consumer"))
		if err != nil {
			return false
		}
		created.Metadata.Labels["label"] = "mutated"
		created.Metadata.Annotations["annotation"] = "mutated"
		created.Spec.ConsumerRef.Kind = "Mutated"
		created.Spec.ConsumerRef.Name = "mutated"

		got, err := reg.GetServiceBinding(context.Background(), name)
		if err != nil {
			return false
		}
		got.Metadata.Labels["label"] = "get-mutated"
		got.Metadata.Annotations["annotation"] = "get-mutated"
		got.Spec.ConsumerRef.Name = "get-mutated"

		items, err := reg.ListServiceBindings(context.Background(), "")
		if err != nil || len(items) != 1 {
			return false
		}
		items[0].Metadata.Labels["label"] = "list-mutated"
		items[0].Metadata.Annotations["annotation"] = "list-mutated"
		items[0].Spec.ConsumerRef.Kind = "ListMutated"

		after, err := reg.GetServiceBinding(context.Background(), name)
		return err == nil &&
			after.Metadata.Labels["label"] == "original" &&
			after.Metadata.Annotations["annotation"] == "original" &&
			after.Spec.ConsumerRef.Kind == "Application" &&
			after.Spec.ConsumerRef.Name == "consumer"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceinstance-servicebinding, Property 16: Duplicate create never overwrites
func TestProperty_ServiceBindingRegistry_DuplicateCreateError(t *testing.T) {
	f := func(seed uint16, originalValue, replacementValue string) bool {
		name := fmt.Sprintf("binding-%d", seed)
		reg := NewServiceBindingRegistry()
		first := propertyServiceBinding(name, "instance-a", "first", originalValue, originalValue, "consumer-a")
		if _, err := reg.CreateServiceBinding(context.Background(), first); err != nil {
			return false
		}
		second := propertyServiceBinding(name, "instance-b", "second", replacementValue, replacementValue, "consumer-b")
		if _, err := reg.CreateServiceBinding(context.Background(), second); !errors.Is(err, ErrAlreadyExists) {
			return false
		}

		got, err := reg.GetServiceBinding(context.Background(), name)
		return err == nil &&
			got.Spec.ServiceInstanceRef == "instance-a" &&
			got.Metadata.DisplayName == "first" &&
			got.Metadata.Labels["label"] == originalValue &&
			got.Spec.ConsumerRef.Name == "consumer-a"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceinstance-servicebinding, Property 17: List filter correctness
func TestProperty_ServiceBindingRegistry_FilterCorrectness(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%20) + 3
		reg := NewServiceBindingRegistry()
		expected := 0
		for i := 0; i < n; i++ {
			instanceRef := fmt.Sprintf("instance-%d", i%3)
			if instanceRef == "instance-1" {
				expected++
			}
			sb := propertyServiceBinding(
				fmt.Sprintf("binding-%02d", i),
				instanceRef,
				"display",
				"value",
				"note",
				fmt.Sprintf("consumer-%02d", i),
			)
			if _, err := reg.CreateServiceBinding(context.Background(), sb); err != nil {
				return false
			}
		}

		all, err := reg.ListServiceBindings(context.Background(), "")
		if err != nil || len(all) != n {
			return false
		}
		filtered, err := reg.ListServiceBindings(context.Background(), "instance-1")
		if err != nil || len(filtered) != expected {
			return false
		}
		for i, item := range filtered {
			if item.Spec.ServiceInstanceRef != "instance-1" {
				return false
			}
			if i > 0 && filtered[i-1].Metadata.Name >= item.Metadata.Name {
				return false
			}
		}
		return true
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceinstance-servicebinding, Property 18: CountByServiceInstance correctness
func TestProperty_ServiceBindingRegistry_CountByServiceInstance(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%20) + 3
		reg := NewServiceBindingRegistry()
		expected := 0
		for i := 0; i < n; i++ {
			instanceRef := fmt.Sprintf("instance-%d", i%4)
			if instanceRef == "instance-2" {
				expected++
			}
			sb := propertyServiceBinding(
				fmt.Sprintf("binding-%02d", i),
				instanceRef,
				"display",
				"value",
				"note",
				fmt.Sprintf("consumer-%02d", i),
			)
			if _, err := reg.CreateServiceBinding(context.Background(), sb); err != nil {
				return false
			}
		}

		got, err := reg.CountByServiceInstance(context.Background(), "instance-2")
		if err != nil || got != expected {
			return false
		}
		missing, err := reg.CountByServiceInstance(context.Background(), "missing")
		return err == nil && missing == 0
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
