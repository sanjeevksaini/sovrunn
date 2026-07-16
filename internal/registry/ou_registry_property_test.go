package registry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// validOU builds an OrganizationUnit with system fields set, used by the
// property tests. It relies on isValidName (defined in the package's
// organization property test) for name/orgName validity checks.
func validOU(orgName, name, display, desc string) resources.OrganizationUnit {
	return resources.OrganizationUnit{
		APIVersion: resources.OUAPIVersion,
		Kind:       resources.OUKind,
		Metadata: resources.Metadata{
			Name:        name,
			DisplayName: display,
			Labels:      map[string]string{"k": "v"},
			Annotations: map[string]string{"a": "b"},
		},
		Spec: resources.OrganizationUnitSpec{
			OrganizationName: orgName,
			Description:      desc,
		},
		Status: resources.OrganizationUnitStatus{Phase: resources.PhaseActive},
	}
}

// Feature: organizationunit-resource, Property 3: Create then Get is a data-preserving round trip
func TestProperty_OURegistry_CreateGetRoundTrip(t *testing.T) {
	f := func(orgName, name, desc string) bool {
		if !isValidName(orgName) || !isValidName(name) {
			return true
		}
		reg := NewOrganizationUnitRegistry()
		ou := validOU(orgName, name, "display", desc)
		if _, err := reg.CreateOrganizationUnit(context.Background(), ou); err != nil {
			return false
		}
		got, err := reg.GetOrganizationUnit(context.Background(), orgName, name)
		if err != nil {
			return false
		}
		return got.Metadata.Name == name &&
			got.Spec.OrganizationName == orgName &&
			got.Spec.Description == desc &&
			got.Metadata.DisplayName == "display" &&
			got.Metadata.Labels["k"] == "v" &&
			got.Metadata.Annotations["a"] == "b" &&
			got.APIVersion == resources.OUAPIVersion &&
			got.Kind == resources.OUKind &&
			got.Status.Phase == resources.PhaseActive
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: organizationunit-resource, Property 4: Registry returns value copies — stored state is immutable to callers
func TestProperty_OURegistry_GetReturnsValueCopy(t *testing.T) {
	f := func(orgName, name string) bool {
		if !isValidName(orgName) || !isValidName(name) {
			return true
		}
		reg := NewOrganizationUnitRegistry()
		ou := validOU(orgName, name, "display", "desc")
		if _, err := reg.CreateOrganizationUnit(context.Background(), ou); err != nil {
			return false
		}
		got, err := reg.GetOrganizationUnit(context.Background(), orgName, name)
		if err != nil {
			return false
		}
		if got.Metadata.Labels != nil {
			got.Metadata.Labels["k"] = "mutated"
		}
		if got.Metadata.Annotations != nil {
			got.Metadata.Annotations["a"] = "mutated"
		}

		list, err := reg.ListOrganizationUnits(context.Background())
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

		after, err := reg.GetOrganizationUnit(context.Background(), orgName, name)
		if err != nil {
			return false
		}
		return after.Metadata.Labels["k"] == "v" && after.Metadata.Annotations["a"] == "b"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: organizationunit-resource, Property 5: List returns OrganizationUnits sorted by organizationName then name
func TestProperty_OURegistry_ListSortedOrder(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%20) + 2
		reg := NewOrganizationUnitRegistry()
		for i := 0; i < n; i++ {
			orgName := fmt.Sprintf("org-%02d", i%4)
			name := fmt.Sprintf("ou-%02d", i)
			ou := validOU(orgName, name, "display", "desc")
			if _, err := reg.CreateOrganizationUnit(context.Background(), ou); err != nil {
				return false
			}
		}
		items, err := reg.ListOrganizationUnits(context.Background())
		if err != nil {
			return false
		}
		for i := 1; i < len(items); i++ {
			prev, cur := items[i-1], items[i]
			if prev.Spec.OrganizationName > cur.Spec.OrganizationName {
				return false
			}
			if prev.Spec.OrganizationName == cur.Spec.OrganizationName &&
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

// Feature: organizationunit-resource, Property 6: Update preserves immutable system fields
func TestProperty_OURegistry_UpdatePreservesSystemFields(t *testing.T) {
	f := func(orgName, name, newDisplay, newDesc string) bool {
		if !isValidName(orgName) || !isValidName(name) {
			return true
		}
		reg := NewOrganizationUnitRegistry()
		original := validOU(orgName, name, "orig", "orig-desc")
		original.Status = resources.OrganizationUnitStatus{Phase: resources.PhaseActive, Message: "ok"}
		if _, err := reg.CreateOrganizationUnit(context.Background(), original); err != nil {
			return false
		}
		update := validOU("tampered-org", "tampered-name", newDisplay, newDesc)
		update.Status = resources.OrganizationUnitStatus{Phase: resources.PhaseFailed, Message: "hacked"}
		update.APIVersion = "tampered/v0"
		update.Kind = "Tampered"
		if _, err := reg.UpdateOrganizationUnit(context.Background(), orgName, name, update); err != nil {
			return false
		}
		got, err := reg.GetOrganizationUnit(context.Background(), orgName, name)
		if err != nil {
			return false
		}
		return got.Metadata.Name == name &&
			got.Spec.OrganizationName == orgName &&
			got.Status.Phase == resources.PhaseActive &&
			got.Status.Message == "ok" &&
			got.APIVersion == resources.OUAPIVersion &&
			got.Kind == resources.OUKind
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: organizationunit-resource, Property 7: Duplicate composite key returns ErrAlreadyExists and original entry is unchanged
func TestProperty_OURegistry_DuplicateCreateError(t *testing.T) {
	f := func(orgName, name string) bool {
		if !isValidName(orgName) || !isValidName(name) {
			return true
		}
		reg := NewOrganizationUnitRegistry()
		first := validOU(orgName, name, "first", "first-desc")
		if _, err := reg.CreateOrganizationUnit(context.Background(), first); err != nil {
			return false
		}
		second := validOU(orgName, name, "second", "second-desc")
		_, err := reg.CreateOrganizationUnit(context.Background(), second)
		if !errors.Is(err, ErrAlreadyExists) {
			return false
		}
		got, err := reg.GetOrganizationUnit(context.Background(), orgName, name)
		if err != nil {
			return false
		}
		return got.Metadata.DisplayName == "first" && got.Spec.Description == "first-desc"
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
