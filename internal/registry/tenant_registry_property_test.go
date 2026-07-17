package registry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// validTenant builds a Tenant with system fields set, used by the property
// tests. It relies on isValidName (defined in the package's organization
// property test) for name validity checks.
func validTenant(orgName, ouName, name, display, desc string) resources.Tenant {
	return resources.Tenant{
		APIVersion: resources.TenantAPIVersion,
		Kind:       resources.TenantKind,
		Metadata: resources.Metadata{
			Name:        name,
			DisplayName: display,
			Labels:      map[string]string{"k": "v"},
			Annotations: map[string]string{"a": "b"},
		},
		Spec: resources.TenantSpec{
			OrganizationName:     orgName,
			OrganizationUnitName: ouName,
			Description:          desc,
		},
		Status: resources.TenantStatus{Phase: resources.PhaseActive},
	}
}

// Feature: tenant-resource, Property 3: Create then Get is a data-preserving round trip
func TestProperty_TenantRegistry_CreateGetRoundTrip(t *testing.T) {
	f := func(orgName, ouName, name, desc string) bool {
		if !isValidName(orgName) || !isValidName(ouName) || !isValidName(name) {
			return true
		}
		reg := NewTenantRegistry()
		tnt := validTenant(orgName, ouName, name, "display", desc)
		if _, err := reg.CreateTenant(context.Background(), tnt); err != nil {
			return false
		}
		got, err := reg.GetTenant(context.Background(), orgName, ouName, name)
		if err != nil {
			return false
		}
		return got.Metadata.Name == name &&
			got.Spec.OrganizationName == orgName &&
			got.Spec.OrganizationUnitName == ouName &&
			got.Spec.Description == desc &&
			got.Metadata.DisplayName == "display" &&
			got.Metadata.Labels["k"] == "v" &&
			got.Metadata.Annotations["a"] == "b" &&
			got.APIVersion == resources.TenantAPIVersion &&
			got.Kind == resources.TenantKind &&
			got.Status.Phase == resources.PhaseActive
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: tenant-resource, Property 4: List returns tenants in correct composite sort order
func TestProperty_TenantRegistry_ListSortedOrder(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%20) + 2
		reg := NewTenantRegistry()
		for i := 0; i < n; i++ {
			orgName := fmt.Sprintf("org-%02d", i%3)
			ouName := fmt.Sprintf("ou-%02d", i%4)
			name := fmt.Sprintf("tenant-%02d", i)
			tnt := validTenant(orgName, ouName, name, "display", "desc")
			if _, err := reg.CreateTenant(context.Background(), tnt); err != nil {
				return false
			}
		}
		items, err := reg.ListTenants(context.Background())
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

// Feature: tenant-resource, Property 5: Registry returns deep copies — mutations don't affect stored state
func TestProperty_TenantRegistry_GetReturnsDeepCopy(t *testing.T) {
	f := func(orgName, ouName, name string) bool {
		if !isValidName(orgName) || !isValidName(ouName) || !isValidName(name) {
			return true
		}
		reg := NewTenantRegistry()
		tnt := validTenant(orgName, ouName, name, "display", "desc")
		if _, err := reg.CreateTenant(context.Background(), tnt); err != nil {
			return false
		}
		got, err := reg.GetTenant(context.Background(), orgName, ouName, name)
		if err != nil {
			return false
		}
		if got.Metadata.Labels != nil {
			got.Metadata.Labels["k"] = "mutated"
		}
		if got.Metadata.Annotations != nil {
			got.Metadata.Annotations["a"] = "mutated"
		}

		list, err := reg.ListTenants(context.Background())
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

		after, err := reg.GetTenant(context.Background(), orgName, ouName, name)
		if err != nil {
			return false
		}
		return after.Metadata.Labels["k"] == "v" && after.Metadata.Annotations["a"] == "b"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: tenant-resource, Property 6: Duplicate create returns ErrAlreadyExists and original unchanged
func TestProperty_TenantRegistry_DuplicateCreateError(t *testing.T) {
	f := func(orgName, ouName, name string) bool {
		if !isValidName(orgName) || !isValidName(ouName) || !isValidName(name) {
			return true
		}
		reg := NewTenantRegistry()
		first := validTenant(orgName, ouName, name, "first", "first-desc")
		if _, err := reg.CreateTenant(context.Background(), first); err != nil {
			return false
		}
		second := validTenant(orgName, ouName, name, "second", "second-desc")
		_, err := reg.CreateTenant(context.Background(), second)
		if !errors.Is(err, ErrAlreadyExists) {
			return false
		}
		got, err := reg.GetTenant(context.Background(), orgName, ouName, name)
		if err != nil {
			return false
		}
		return got.Metadata.DisplayName == "first" && got.Spec.Description == "first-desc"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
