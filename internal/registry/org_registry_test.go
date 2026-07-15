package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func sampleOrg(name string) resources.Organization {
	return resources.Organization{
		APIVersion: resources.OrgAPIVersion,
		Kind:       resources.OrgKind,
		Metadata: resources.Metadata{
			Name:        name,
			DisplayName: "Display",
			Labels:      map[string]string{"env": "test"},
			Annotations: map[string]string{"note": "x"},
		},
		Spec: resources.OrganizationSpec{
			Description:          "desc",
			SovereignLocations:   []string{"in"},
			DefaultPolicyProfile: "default",
		},
		Status: resources.OrganizationStatus{Phase: resources.PhaseActive},
	}
}

func TestCreateOrganization_Stores(t *testing.T) {
	reg := NewOrganizationRegistry()
	org := sampleOrg("nic")
	if err := reg.CreateOrganization(context.Background(), org); err != nil {
		t.Fatalf("CreateOrganization() error = %v", err)
	}
	got, err := reg.GetOrganization(context.Background(), "nic")
	if err != nil {
		t.Fatalf("GetOrganization() error = %v", err)
	}
	if got.Metadata.Name != "nic" {
		t.Errorf("Name = %q, want %q", got.Metadata.Name, "nic")
	}
}

func TestCreateOrganization_Duplicate(t *testing.T) {
	reg := NewOrganizationRegistry()
	org := sampleOrg("nic")
	_ = reg.CreateOrganization(context.Background(), org)
	err := reg.CreateOrganization(context.Background(), org)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
}

func TestGetOrganization_NotFound(t *testing.T) {
	reg := NewOrganizationRegistry()
	_, err := reg.GetOrganization(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestListOrganizations_Empty(t *testing.T) {
	reg := NewOrganizationRegistry()
	items, err := reg.ListOrganizations(context.Background())
	if err != nil {
		t.Fatalf("ListOrganizations() error = %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("got %d items, want 0", len(items))
	}
}

func TestListOrganizations_Sorted(t *testing.T) {
	reg := NewOrganizationRegistry()
	names := []string{"zebra", "alpha", "mike"}
	for _, n := range names {
		if err := reg.CreateOrganization(context.Background(), sampleOrg(n)); err != nil {
			t.Fatalf("CreateOrganization(%s): %v", n, err)
		}
	}
	items, err := reg.ListOrganizations(context.Background())
	if err != nil {
		t.Fatalf("ListOrganizations() error = %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("got %d items, want 3", len(items))
	}
	for i := 1; i < len(items); i++ {
		if items[i-1].Metadata.Name >= items[i].Metadata.Name {
			t.Fatalf("not sorted: %v", items)
		}
	}
}

func TestUpdateOrganization_MutableFields(t *testing.T) {
	reg := NewOrganizationRegistry()
	_ = reg.CreateOrganization(context.Background(), sampleOrg("nic"))
	updated := sampleOrg("nic")
	updated.Metadata.DisplayName = "New Display"
	updated.Spec.Description = "new desc"
	got, err := reg.UpdateOrganization(context.Background(), "nic", updated)
	if err != nil {
		t.Fatalf("UpdateOrganization() error = %v", err)
	}
	if got.Metadata.DisplayName != "New Display" {
		t.Errorf("DisplayName = %q, want %q", got.Metadata.DisplayName, "New Display")
	}
	if got.Spec.Description != "new desc" {
		t.Errorf("Description = %q, want %q", got.Spec.Description, "new desc")
	}
	if got.Metadata.Name != "nic" {
		t.Errorf("Name changed to %q", got.Metadata.Name)
	}
	if got.Status.Phase != resources.PhaseActive {
		t.Errorf("Status.Phase = %q, want Active", got.Status.Phase)
	}
}

func TestUpdateOrganization_NotFound(t *testing.T) {
	reg := NewOrganizationRegistry()
	_, err := reg.UpdateOrganization(context.Background(), "missing", sampleOrg("missing"))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestDeleteOrganization_Exists(t *testing.T) {
	reg := NewOrganizationRegistry()
	_ = reg.CreateOrganization(context.Background(), sampleOrg("nic"))
	if err := reg.DeleteOrganization(context.Background(), "nic"); err != nil {
		t.Fatalf("DeleteOrganization() error = %v", err)
	}
	_, err := reg.GetOrganization(context.Background(), "nic")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound after delete", err)
	}
}

func TestDeleteOrganization_NotFound(t *testing.T) {
	reg := NewOrganizationRegistry()
	err := reg.DeleteOrganization(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestGetOrganization_DeepCopy(t *testing.T) {
	reg := NewOrganizationRegistry()
	_ = reg.CreateOrganization(context.Background(), sampleOrg("nic"))
	got, _ := reg.GetOrganization(context.Background(), "nic")
	got.Metadata.Labels["env"] = "mutated"
	got.Spec.SovereignLocations[0] = "us"
	got2, _ := reg.GetOrganization(context.Background(), "nic")
	if got2.Metadata.Labels["env"] != "test" {
		t.Error("labels mutation affected stored value")
	}
	if got2.Spec.SovereignLocations[0] != "in" {
		t.Error("slice mutation affected stored value")
	}
}
