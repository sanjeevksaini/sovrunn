package registry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// validProject builds a Project with system fields set, used by the property
// tests. It relies on isValidName (defined in the package's organization
// property test) for name validity checks.
func validProject(orgName, ouName, tenantName, name, display, desc string) resources.Project {
	return resources.Project{
		APIVersion: resources.ProjectAPIVersion,
		Kind:       resources.ProjectKind,
		Metadata: resources.Metadata{
			Name:        name,
			DisplayName: display,
			Labels:      map[string]string{"k": "v"},
			Annotations: map[string]string{"a": "b"},
		},
		Spec: resources.ProjectSpec{
			OrganizationName:     orgName,
			OrganizationUnitName: ouName,
			TenantName:           tenantName,
			Description:          desc,
		},
		Status: resources.ProjectStatus{Phase: resources.PhaseActive},
	}
}

// Feature: project-resource, Property 3: Create then Get is a data-preserving round trip
func TestProperty_ProjectRegistry_CreateGetRoundTrip(t *testing.T) {
	f := func(orgName, ouName, tenantName, name, desc string) bool {
		if !isValidName(orgName) || !isValidName(ouName) || !isValidName(tenantName) || !isValidName(name) {
			return true
		}
		reg := NewProjectRegistry()
		project := validProject(orgName, ouName, tenantName, name, "display", desc)
		if _, err := reg.CreateProject(context.Background(), project); err != nil {
			return false
		}
		got, err := reg.GetProject(context.Background(), orgName, ouName, tenantName, name)
		if err != nil {
			return false
		}
		return got.Metadata.Name == name &&
			got.Spec.OrganizationName == orgName &&
			got.Spec.OrganizationUnitName == ouName &&
			got.Spec.TenantName == tenantName &&
			got.Spec.Description == desc &&
			got.Metadata.DisplayName == "display" &&
			got.Metadata.Labels["k"] == "v" &&
			got.Metadata.Annotations["a"] == "b" &&
			got.APIVersion == resources.ProjectAPIVersion &&
			got.Kind == resources.ProjectKind &&
			got.Status.Phase == resources.PhaseActive
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: project-resource, Property 4: List returns projects in correct four-level composite sort order
func TestProperty_ProjectRegistry_ListSortedOrder(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%20) + 2
		reg := NewProjectRegistry()
		for i := 0; i < n; i++ {
			orgName := fmt.Sprintf("org-%02d", i%3)
			ouName := fmt.Sprintf("ou-%02d", i%4)
			tenantName := fmt.Sprintf("tenant-%02d", i%5)
			name := fmt.Sprintf("project-%02d", i)
			project := validProject(orgName, ouName, tenantName, name, "display", "desc")
			if _, err := reg.CreateProject(context.Background(), project); err != nil {
				return false
			}
		}
		items, err := reg.ListProjects(context.Background())
		if err != nil {
			return false
		}
		for i := 1; i < len(items); i++ {
			prev, cur := items[i-1], items[i]
			if prev.Spec.OrganizationName > cur.Spec.OrganizationName {
				return false
			}
			if prev.Spec.OrganizationName == cur.Spec.OrganizationName &&
				prev.Spec.OrganizationUnitName > cur.Spec.OrganizationUnitName {
				return false
			}
			if prev.Spec.OrganizationName == cur.Spec.OrganizationName &&
				prev.Spec.OrganizationUnitName == cur.Spec.OrganizationUnitName &&
				prev.Spec.TenantName > cur.Spec.TenantName {
				return false
			}
			if prev.Spec.OrganizationName == cur.Spec.OrganizationName &&
				prev.Spec.OrganizationUnitName == cur.Spec.OrganizationUnitName &&
				prev.Spec.TenantName == cur.Spec.TenantName &&
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

// Feature: project-resource, Property 5: Registry returns deep copies — mutations don't affect stored state
func TestProperty_ProjectRegistry_GetReturnsDeepCopy(t *testing.T) {
	f := func(orgName, ouName, tenantName, name string) bool {
		if !isValidName(orgName) || !isValidName(ouName) || !isValidName(tenantName) || !isValidName(name) {
			return true
		}
		reg := NewProjectRegistry()
		project := validProject(orgName, ouName, tenantName, name, "display", "desc")
		if _, err := reg.CreateProject(context.Background(), project); err != nil {
			return false
		}
		got, err := reg.GetProject(context.Background(), orgName, ouName, tenantName, name)
		if err != nil {
			return false
		}
		if got.Metadata.Labels != nil {
			got.Metadata.Labels["k"] = "mutated"
		}
		if got.Metadata.Annotations != nil {
			got.Metadata.Annotations["a"] = "mutated"
		}

		list, err := reg.ListProjects(context.Background())
		if err != nil {
			return false
		}
		for i := range list {
			if list[i].Metadata.Labels != nil {
				list[i].Metadata.Labels["k"] = "list-mutated"
			}
			if list[i].Metadata.Annotations != nil {
				list[i].Metadata.Annotations["a"] = "list-mutated"
			}
		}

		after, err := reg.GetProject(context.Background(), orgName, ouName, tenantName, name)
		if err != nil {
			return false
		}
		return after.Metadata.Labels["k"] == "v" && after.Metadata.Annotations["a"] == "b"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: project-resource, Property 6: Duplicate create returns ErrAlreadyExists and original unchanged
func TestProperty_ProjectRegistry_DuplicateCreateError(t *testing.T) {
	f := func(orgName, ouName, tenantName, name string) bool {
		if !isValidName(orgName) || !isValidName(ouName) || !isValidName(tenantName) || !isValidName(name) {
			return true
		}
		reg := NewProjectRegistry()
		first := validProject(orgName, ouName, tenantName, name, "first", "first-desc")
		if _, err := reg.CreateProject(context.Background(), first); err != nil {
			return false
		}
		second := validProject(orgName, ouName, tenantName, name, "second", "second-desc")
		_, err := reg.CreateProject(context.Background(), second)
		if !errors.Is(err, ErrAlreadyExists) {
			return false
		}
		got, err := reg.GetProject(context.Background(), orgName, ouName, tenantName, name)
		if err != nil {
			return false
		}
		return got.Metadata.DisplayName == "first" && got.Spec.Description == "first-desc"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
