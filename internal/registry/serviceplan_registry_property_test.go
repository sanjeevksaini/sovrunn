package registry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func validServicePlan(className, name, display, desc string) resources.ServicePlan {
	return resources.ServicePlan{
		APIVersion: "platform.sovrunn.io/v1alpha1",
		Kind:       resources.ServicePlanKind,
		Metadata: resources.Metadata{
			Name:        name,
			Labels:      map[string]string{"k": "v"},
			Annotations: map[string]string{"a": "b"},
		},
		Spec: resources.ServicePlanSpec{
			ServiceClassName: className,
			DisplayName:      display,
			Description:      desc,
			Tier:             resources.TierSmall,
			Lifecycle:        resources.LifecycleActive,
			Parameters:       map[string]string{"region": "us-east"},
			Tags:             []string{"tag"},
		},
		Status: resources.ServicePlanStatus{Phase: resources.PhaseActive},
	}
}

// Feature: serviceclass-serviceplan, Property 5: ServicePlan composite identity is stable
func TestProperty_ServicePlanRegistry_CompositeIdentityStable(t *testing.T) {
	f := func(className, name, display, desc string) bool {
		if !isValidName(className) || !isValidName(name) {
			return true
		}
		otherClass := className + "-other"
		if !isValidName(otherClass) {
			otherClass = "other-" + className
			if !isValidName(otherClass) {
				return true
			}
		}
		reg := NewServicePlanRegistry()
		first := validServicePlan(className, name, display, desc)
		if _, err := reg.CreateServicePlan(context.Background(), first); err != nil {
			return false
		}
		got, err := reg.GetServicePlan(context.Background(), className, name)
		if err != nil {
			return false
		}
		if got.Metadata.Name != name ||
			got.Spec.ServiceClassName != className ||
			got.Spec.DisplayName != display ||
			got.Spec.Description != desc ||
			got.Spec.Parameters["region"] != "us-east" {
			return false
		}
		second := validServicePlan(otherClass, name, "other", "other-desc")
		if _, err := reg.CreateServicePlan(context.Background(), second); err != nil {
			return false
		}
		gotOther, err := reg.GetServicePlan(context.Background(), otherClass, name)
		if err != nil {
			return false
		}
		original, err := reg.GetServicePlan(context.Background(), className, name)
		if err != nil {
			return false
		}
		return gotOther.Spec.ServiceClassName == otherClass &&
			original.Spec.ServiceClassName == className &&
			original.Spec.Description == desc
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceclass-serviceplan, Property 6: List ordering is deterministic
func TestProperty_ServicePlanRegistry_ListSortedOrder(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%20) + 2
		reg := NewServicePlanRegistry()
		for i := 0; i < n; i++ {
			className := fmt.Sprintf("class-%02d", i%3)
			name := fmt.Sprintf("plan-%02d", i)
			if _, err := reg.CreateServicePlan(context.Background(), validServicePlan(className, name, "d", "desc")); err != nil {
				return false
			}
		}
		items, err := reg.ListServicePlans(context.Background())
		if err != nil {
			return false
		}
		for i := 1; i < len(items); i++ {
			prev, cur := items[i-1], items[i]
			if prev.Spec.ServiceClassName > cur.Spec.ServiceClassName {
				return false
			}
			if prev.Spec.ServiceClassName == cur.Spec.ServiceClassName &&
				prev.Metadata.Name >= cur.Metadata.Name {
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
func TestProperty_ServicePlanRegistry_ReturnsDeepCopies(t *testing.T) {
	f := func(className, name string) bool {
		if !isValidName(className) || !isValidName(name) {
			return true
		}
		reg := NewServicePlanRegistry()
		if _, err := reg.CreateServicePlan(context.Background(), validServicePlan(className, name, "d", "desc")); err != nil {
			return false
		}
		got, err := reg.GetServicePlan(context.Background(), className, name)
		if err != nil {
			return false
		}
		got.Metadata.Labels["k"] = "mutated"
		got.Metadata.Annotations["a"] = "mutated"
		got.Spec.Parameters["region"] = "mutated"
		got.Spec.Tags[0] = "mutated"

		list, err := reg.ListServicePlans(context.Background())
		if err != nil {
			return false
		}
		for i := range list {
			list[i].Metadata.Labels["k"] = "list-mutated"
			list[i].Spec.Parameters["region"] = "list-mutated"
			list[i].Spec.Tags[0] = "list-mutated"
		}

		after, err := reg.GetServicePlan(context.Background(), className, name)
		if err != nil {
			return false
		}
		return after.Metadata.Labels["k"] == "v" &&
			after.Metadata.Annotations["a"] == "b" &&
			after.Spec.Parameters["region"] == "us-east" &&
			after.Spec.Tags[0] == "tag"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceclass-serviceplan, Property 8: Duplicate create never overwrites
func TestProperty_ServicePlanRegistry_DuplicateCreateError(t *testing.T) {
	f := func(className, name string) bool {
		if !isValidName(className) || !isValidName(name) {
			return true
		}
		reg := NewServicePlanRegistry()
		first := validServicePlan(className, name, "first", "first-desc")
		if _, err := reg.CreateServicePlan(context.Background(), first); err != nil {
			return false
		}
		second := validServicePlan(className, name, "second", "second-desc")
		_, err := reg.CreateServicePlan(context.Background(), second)
		if !errors.Is(err, ErrAlreadyExists) {
			return false
		}
		got, err := reg.GetServicePlan(context.Background(), className, name)
		if err != nil {
			return false
		}
		return got.Spec.DisplayName == "first" && got.Spec.Description == "first-desc"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
