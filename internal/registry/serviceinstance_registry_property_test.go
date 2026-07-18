package registry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func validServiceInstance(name, org, ou, tenant, project, class, plan, display, paramVal string) resources.ServiceInstance {
	return resources.ServiceInstance{
		APIVersion: resources.ServiceInstanceAPIVersion,
		Kind:       resources.ServiceInstanceKind,
		Metadata: resources.Metadata{
			Name:        name,
			DisplayName: display,
			Labels:      map[string]string{"k": "v"},
			Annotations: map[string]string{"a": "b"},
		},
		Spec: resources.ServiceInstanceSpec{
			OrganizationRef:     org,
			OrganizationUnitRef: ou,
			TenantRef:           tenant,
			ProjectRef:          project,
			ServiceClassRef:     class,
			ServicePlanRef:      plan,
			Parameters:          map[string]string{"region": paramVal},
		},
		Status: resources.ServiceInstanceStatus{
			Phase:   "Ready",
			Message: "Registered only; no real provisioning in Phase 1",
		},
	}
}

// Feature: serviceinstance-servicebinding, Property 6: Create/Get round-trip preserves data
func TestProperty_ServiceInstanceRegistry_CreateGetRoundTrip(t *testing.T) {
	f := func(name, org, ou, tenant, project, class, plan, display, paramVal string) bool {
		if !isValidName(name) || !isValidName(org) || !isValidName(tenant) ||
			!isValidName(project) || !isValidName(class) || !isValidName(plan) {
			return true
		}
		if ou != "" && !isValidName(ou) {
			return true
		}
		reg := NewServiceInstanceRegistry()
		si := validServiceInstance(name, org, ou, tenant, project, class, plan, display, paramVal)
		if _, err := reg.CreateServiceInstance(context.Background(), si); err != nil {
			return false
		}
		got, err := reg.GetServiceInstance(context.Background(), name)
		if err != nil {
			return false
		}
		return got.Metadata.Name == name &&
			got.Metadata.DisplayName == display &&
			got.Spec.OrganizationRef == org &&
			got.Spec.OrganizationUnitRef == ou &&
			got.Spec.TenantRef == tenant &&
			got.Spec.ProjectRef == project &&
			got.Spec.ServiceClassRef == class &&
			got.Spec.ServicePlanRef == plan &&
			got.Spec.Parameters["region"] == paramVal &&
			got.Metadata.Labels["k"] == "v" &&
			got.Metadata.Annotations["a"] == "b" &&
			got.APIVersion == resources.ServiceInstanceAPIVersion &&
			got.Kind == resources.ServiceInstanceKind &&
			got.Status.Phase == "Ready"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceinstance-servicebinding, Property 7: List sort invariant
func TestProperty_ServiceInstanceRegistry_ListSortedOrder(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%20) + 2
		reg := NewServiceInstanceRegistry()
		for i := 0; i < n; i++ {
			name := fmt.Sprintf("si-%02d", (n-1-i)%n)
			si := validServiceInstance(name, "org-a", "", "tenant-a", "project-a", "class-a", "plan-a", "d", "us-east")
			if _, err := reg.CreateServiceInstance(context.Background(), si); err != nil {
				return false
			}
		}
		items, err := reg.ListServiceInstances(context.Background(), "", "")
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

// Feature: serviceinstance-servicebinding, Property 8: Deep-copy immutability
func TestProperty_ServiceInstanceRegistry_ReturnsDeepCopies(t *testing.T) {
	f := func(name string) bool {
		if !isValidName(name) {
			return true
		}
		reg := NewServiceInstanceRegistry()
		if _, err := reg.CreateServiceInstance(context.Background(),
			validServiceInstance(name, "org-a", "", "tenant-a", "project-a", "class-a", "plan-a", "d", "us-east")); err != nil {
			return false
		}
		got, err := reg.GetServiceInstance(context.Background(), name)
		if err != nil {
			return false
		}
		got.Metadata.Labels["k"] = "mutated"
		got.Metadata.Annotations["a"] = "mutated"
		got.Spec.Parameters["region"] = "mutated"

		list, err := reg.ListServiceInstances(context.Background(), "", "")
		if err != nil {
			return false
		}
		for i := range list {
			list[i].Metadata.Labels["k"] = "list-mutated"
			list[i].Spec.Parameters["region"] = "list-mutated"
		}

		after, err := reg.GetServiceInstance(context.Background(), name)
		if err != nil {
			return false
		}
		return after.Metadata.Labels["k"] == "v" &&
			after.Metadata.Annotations["a"] == "b" &&
			after.Spec.Parameters["region"] == "us-east"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceinstance-servicebinding, Property 9: Duplicate create never overwrites
func TestProperty_ServiceInstanceRegistry_DuplicateCreateError(t *testing.T) {
	f := func(name string) bool {
		if !isValidName(name) {
			return true
		}
		reg := NewServiceInstanceRegistry()
		first := validServiceInstance(name, "org-a", "", "tenant-a", "project-a", "class-a", "plan-a", "first", "original")
		if _, err := reg.CreateServiceInstance(context.Background(), first); err != nil {
			return false
		}
		second := validServiceInstance(name, "org-b", "ou-b", "tenant-b", "project-b", "class-b", "plan-b", "second", "changed")
		_, err := reg.CreateServiceInstance(context.Background(), second)
		if !errors.Is(err, ErrAlreadyExists) {
			return false
		}
		got, err := reg.GetServiceInstance(context.Background(), name)
		if err != nil {
			return false
		}
		return got.Metadata.DisplayName == "first" &&
			got.Spec.OrganizationRef == "org-a" &&
			got.Spec.Parameters["region"] == "original"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceinstance-servicebinding, Property 10: List filter correctness
func TestProperty_ServiceInstanceRegistry_FilterCorrectness(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%10) + 3
		reg := NewServiceInstanceRegistry()
		for i := 0; i < n; i++ {
			name := fmt.Sprintf("si-%02d", i)
			tenant := fmt.Sprintf("tenant-%d", i%2)
			project := fmt.Sprintf("project-%d", i%3)
			si := validServiceInstance(name, "org-a", "", tenant, project, "class-a", "plan-a", "d", "us-east")
			if _, err := reg.CreateServiceInstance(context.Background(), si); err != nil {
				return false
			}
		}

		all, err := reg.ListServiceInstances(context.Background(), "", "")
		if err != nil || len(all) != n {
			return false
		}

		tenantFiltered, err := reg.ListServiceInstances(context.Background(), "tenant-0", "")
		if err != nil {
			return false
		}
		for _, item := range tenantFiltered {
			if item.Spec.TenantRef != "tenant-0" {
				return false
			}
		}

		projectFiltered, err := reg.ListServiceInstances(context.Background(), "", "project-1")
		if err != nil {
			return false
		}
		for _, item := range projectFiltered {
			if item.Spec.ProjectRef != "project-1" {
				return false
			}
		}

		bothFiltered, err := reg.ListServiceInstances(context.Background(), "tenant-0", "project-1")
		if err != nil {
			return false
		}
		for _, item := range bothFiltered {
			if item.Spec.TenantRef != "tenant-0" || item.Spec.ProjectRef != "project-1" {
				return false
			}
		}

		expectedBoth := 0
		for _, item := range all {
			if item.Spec.TenantRef == "tenant-0" && item.Spec.ProjectRef == "project-1" {
				expectedBoth++
			}
		}
		if len(bothFiltered) != expectedBoth {
			return false
		}

		for i := 1; i < len(tenantFiltered); i++ {
			if tenantFiltered[i-1].Metadata.Name >= tenantFiltered[i].Metadata.Name {
				return false
			}
		}
		return true
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceinstance-servicebinding, Property 11: CountByServicePlan correctness
func TestProperty_ServiceInstanceRegistry_CountByServicePlan(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%10) + 3
		reg := NewServiceInstanceRegistry()
		expected := 0
		for i := 0; i < n; i++ {
			name := fmt.Sprintf("si-%02d", i)
			class := fmt.Sprintf("class-%d", i%3)
			plan := fmt.Sprintf("plan-%d", i%2)
			si := validServiceInstance(name, "org-a", "", "tenant-a", "project-a", class, plan, "d", "us-east")
			if _, err := reg.CreateServiceInstance(context.Background(), si); err != nil {
				return false
			}
			if class == "class-0" && plan == "plan-0" {
				expected++
			}
		}

		got, err := reg.CountByServicePlan(context.Background(), "class-0", "plan-0")
		if err != nil {
			return false
		}
		if got != expected {
			return false
		}

		// Same plan name under a different class must not be counted.
		falsePositive, err := reg.CountByServicePlan(context.Background(), "class-1", "plan-0")
		if err != nil {
			return false
		}
		manual := 0
		all, err := reg.ListServiceInstances(context.Background(), "", "")
		if err != nil {
			return false
		}
		for _, item := range all {
			if item.Spec.ServiceClassRef == "class-1" && item.Spec.ServicePlanRef == "plan-0" {
				manual++
			}
		}
		return falsePositive == manual
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceinstance-servicebinding, Property 12: CountByProject correctness
func TestProperty_ServiceInstanceRegistry_CountByProject(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%10) + 4
		reg := NewServiceInstanceRegistry()
		expectedEmptyOU := 0
		expectedWithOU := 0
		for i := 0; i < n; i++ {
			name := fmt.Sprintf("si-%02d", i)
			org := fmt.Sprintf("org-%d", i%2)
			tenant := fmt.Sprintf("tenant-%d", i%3)
			project := fmt.Sprintf("project-%d", i%2)
			ou := ""
			if i%4 == 0 {
				ou = "ou-a"
			}
			si := validServiceInstance(name, org, ou, tenant, project, "class-a", "plan-a", "d", "us-east")
			if _, err := reg.CreateServiceInstance(context.Background(), si); err != nil {
				return false
			}
			if org == "org-0" && ou == "" && tenant == "tenant-0" && project == "project-0" {
				expectedEmptyOU++
			}
			if org == "org-0" && ou == "ou-a" && tenant == "tenant-0" && project == "project-0" {
				expectedWithOU++
			}
		}

		gotEmpty, err := reg.CountByProject(context.Background(), "org-0", "", "tenant-0", "project-0")
		if err != nil || gotEmpty != expectedEmptyOU {
			return false
		}

		gotWithOU, err := reg.CountByProject(context.Background(), "org-0", "ou-a", "tenant-0", "project-0")
		if err != nil || gotWithOU != expectedWithOU {
			return false
		}

		// Empty-OU and non-empty-OU scopes are distinct; counts must not be conflated.
		return true
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
