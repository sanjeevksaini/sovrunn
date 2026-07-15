package registry

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

var dnsLabelRe = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

// Feature: organization-resource-registry, Property 3: Create then Get is a data-preserving round trip
func TestProperty_Registry_CreateGetRoundTrip(t *testing.T) {
	f := func(name string, desc string) bool {
		if !isValidName(name) {
			return true
		}
		reg := NewOrganizationRegistry()
		org := resources.Organization{
			Metadata: resources.Metadata{Name: name},
			Spec:     resources.OrganizationSpec{Description: desc},
		}
		org.APIVersion = resources.OrgAPIVersion
		org.Kind = resources.OrgKind
		org.Status.Phase = resources.PhaseActive
		if err := reg.CreateOrganization(context.Background(), org); err != nil {
			return false
		}
		got, err := reg.GetOrganization(context.Background(), name)
		if err != nil {
			return false
		}
		return got.Metadata.Name == name &&
			got.Spec.Description == desc &&
			got.APIVersion == resources.OrgAPIVersion &&
			got.Kind == resources.OrgKind &&
			got.Status.Phase == resources.PhaseActive
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: organization-resource-registry, Property 4: Registry returns value copies
func TestProperty_Registry_GetReturnsValueCopy(t *testing.T) {
	f := func(name string) bool {
		if !isValidName(name) {
			return true
		}
		reg := NewOrganizationRegistry()
		org := resources.Organization{
			APIVersion: resources.OrgAPIVersion,
			Kind:       resources.OrgKind,
			Metadata: resources.Metadata{
				Name:        name,
				Labels:      map[string]string{"k": "v"},
				Annotations: map[string]string{"a": "b"},
			},
			Spec: resources.OrganizationSpec{
				SovereignLocations: []string{"in"},
			},
			Status: resources.OrganizationStatus{Phase: resources.PhaseActive},
		}
		if err := reg.CreateOrganization(context.Background(), org); err != nil {
			return false
		}
		got, err := reg.GetOrganization(context.Background(), name)
		if err != nil {
			return false
		}
		if got.Metadata.Labels != nil {
			got.Metadata.Labels["k"] = "mutated"
		}
		if got.Metadata.Annotations != nil {
			got.Metadata.Annotations["a"] = "mutated"
		}
		if got.Spec.SovereignLocations != nil {
			got.Spec.SovereignLocations[0] = "us"
		}
		got2, err := reg.GetOrganization(context.Background(), name)
		if err != nil {
			return false
		}
		if got2.Metadata.Labels != nil && got2.Metadata.Labels["k"] != "v" {
			return false
		}
		if got2.Metadata.Annotations != nil && got2.Metadata.Annotations["a"] != "b" {
			return false
		}
		if got2.Spec.SovereignLocations != nil && got2.Spec.SovereignLocations[0] != "in" {
			return false
		}
		return true
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: organization-resource-registry, Property 5: List returns organizations in ascending lexicographic order
func TestProperty_Registry_ListSortedOrder(t *testing.T) {
	f := func(count uint8) bool {
		n := int(count%20) + 2
		reg := NewOrganizationRegistry()
		seen := make(map[string]bool)
		for i := 0; i < n; i++ {
			name := fmt.Sprintf("org-%02d", i)
			if seen[name] {
				continue
			}
			seen[name] = true
			org := resources.Organization{
				APIVersion: resources.OrgAPIVersion,
				Kind:       resources.OrgKind,
				Metadata:   resources.Metadata{Name: name},
				Status:     resources.OrganizationStatus{Phase: resources.PhaseActive},
			}
			if err := reg.CreateOrganization(context.Background(), org); err != nil {
				return false
			}
		}
		items, err := reg.ListOrganizations(context.Background())
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

// Feature: organization-resource-registry, Property 6: Update preserves immutable system fields
func TestProperty_Registry_UpdatePreservesSystemFields(t *testing.T) {
	f := func(name string, newDisplay string) bool {
		if !isValidName(name) {
			return true
		}
		reg := NewOrganizationRegistry()
		original := resources.Organization{
			APIVersion: resources.OrgAPIVersion,
			Kind:       resources.OrgKind,
			Metadata:   resources.Metadata{Name: name, DisplayName: "orig"},
			Status:     resources.OrganizationStatus{Phase: resources.PhaseActive, Message: "ok"},
		}
		if err := reg.CreateOrganization(context.Background(), original); err != nil {
			return false
		}
		update := resources.Organization{
			Metadata: resources.Metadata{Name: name, DisplayName: newDisplay},
			Spec:     resources.OrganizationSpec{Description: "updated"},
		}
		if _, err := reg.UpdateOrganization(context.Background(), name, update); err != nil {
			return false
		}
		got, err := reg.GetOrganization(context.Background(), name)
		if err != nil {
			return false
		}
		return got.Metadata.Name == name &&
			got.Status.Phase == original.Status.Phase &&
			got.Status.Message == original.Status.Message
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

func isValidName(name string) bool {
	return name != "" && len(name) <= 63 && dnsLabelRe.MatchString(name)
}
