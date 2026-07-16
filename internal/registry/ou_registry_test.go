package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func sampleOU(orgName, name string) resources.OrganizationUnit {
	return resources.OrganizationUnit{
		APIVersion: resources.OUAPIVersion,
		Kind:       resources.OUKind,
		Metadata: resources.Metadata{
			Name:        name,
			DisplayName: "Display",
			Labels:      map[string]string{"env": "test"},
			Annotations: map[string]string{"note": "x"},
		},
		Spec: resources.OrganizationUnitSpec{
			OrganizationName: orgName,
			Description:      "desc",
		},
		Status: resources.OrganizationUnitStatus{Phase: resources.PhaseActive},
	}
}

func TestCreateOrganizationUnit_Stores(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	created, err := reg.CreateOrganizationUnit(context.Background(), sampleOU("nic", "ministry-health"))
	if err != nil {
		t.Fatalf("CreateOrganizationUnit() error = %v", err)
	}
	if created.Metadata.Name != "ministry-health" {
		t.Errorf("Name = %q, want %q", created.Metadata.Name, "ministry-health")
	}
	got, err := reg.GetOrganizationUnit(context.Background(), "nic", "ministry-health")
	if err != nil {
		t.Fatalf("GetOrganizationUnit() error = %v", err)
	}
	if got.Metadata.Name != "ministry-health" || got.Spec.OrganizationName != "nic" {
		t.Errorf("got %q/%q, want nic/ministry-health", got.Spec.OrganizationName, got.Metadata.Name)
	}
}

func TestCreateOrganizationUnit_Duplicate(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	ou := sampleOU("nic", "ministry-health")
	if _, err := reg.CreateOrganizationUnit(context.Background(), ou); err != nil {
		t.Fatalf("first CreateOrganizationUnit() error = %v", err)
	}
	_, err := reg.CreateOrganizationUnit(context.Background(), ou)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
}

func TestCreateOrganizationUnit_SameNameDifferentOrgs(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	if _, err := reg.CreateOrganizationUnit(context.Background(), sampleOU("nic", "ministry-health")); err != nil {
		t.Fatalf("create nic/ministry-health: %v", err)
	}
	if _, err := reg.CreateOrganizationUnit(context.Background(), sampleOU("state-gov", "ministry-health")); err != nil {
		t.Fatalf("create state-gov/ministry-health: %v", err)
	}
	if _, err := reg.GetOrganizationUnit(context.Background(), "nic", "ministry-health"); err != nil {
		t.Errorf("get nic/ministry-health: %v", err)
	}
	if _, err := reg.GetOrganizationUnit(context.Background(), "state-gov", "ministry-health"); err != nil {
		t.Errorf("get state-gov/ministry-health: %v", err)
	}
}

func TestGetOrganizationUnit_ByCompositeKey(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	_, _ = reg.CreateOrganizationUnit(context.Background(), sampleOU("nic", "ministry-health"))
	got, err := reg.GetOrganizationUnit(context.Background(), "nic", "ministry-health")
	if err != nil {
		t.Fatalf("GetOrganizationUnit() error = %v", err)
	}
	if got.Spec.Description != "desc" {
		t.Errorf("Description = %q, want %q", got.Spec.Description, "desc")
	}
}

func TestGetOrganizationUnit_NotFound(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	_, err := reg.GetOrganizationUnit(context.Background(), "nic", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestListOrganizationUnits_Empty(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	items, err := reg.ListOrganizationUnits(context.Background())
	if err != nil {
		t.Fatalf("ListOrganizationUnits() error = %v", err)
	}
	if items == nil {
		t.Fatal("got nil slice, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Fatalf("got %d items, want 0", len(items))
	}
}

func TestListOrganizationUnits_Sorted(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	inputs := []struct{ org, name string }{
		{"zebra", "beta"},
		{"alpha", "delta"},
		{"alpha", "charlie"},
		{"nic", "ministry-health"},
	}
	for _, in := range inputs {
		if _, err := reg.CreateOrganizationUnit(context.Background(), sampleOU(in.org, in.name)); err != nil {
			t.Fatalf("create %s/%s: %v", in.org, in.name, err)
		}
	}
	items, err := reg.ListOrganizationUnits(context.Background())
	if err != nil {
		t.Fatalf("ListOrganizationUnits() error = %v", err)
	}
	if len(items) != 4 {
		t.Fatalf("got %d items, want 4", len(items))
	}
	for i := 1; i < len(items); i++ {
		prev, cur := items[i-1], items[i]
		if prev.Spec.OrganizationName > cur.Spec.OrganizationName {
			t.Fatalf("not sorted by org: %v", items)
		}
		if prev.Spec.OrganizationName == cur.Spec.OrganizationName &&
			prev.Metadata.Name >= cur.Metadata.Name {
			t.Fatalf("not sorted by name within org: %v", items)
		}
	}
	if items[0].Spec.OrganizationName != "alpha" || items[0].Metadata.Name != "charlie" {
		t.Errorf("first item = %s/%s, want alpha/charlie", items[0].Spec.OrganizationName, items[0].Metadata.Name)
	}
}

func TestUpdateOrganizationUnit_MutableFields(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	_, _ = reg.CreateOrganizationUnit(context.Background(), sampleOU("nic", "ministry-health"))
	update := sampleOU("nic", "ministry-health")
	update.Metadata.DisplayName = "New Display"
	update.Metadata.Labels = map[string]string{"tier": "gold"}
	update.Metadata.Annotations = map[string]string{"reviewed": "yes"}
	update.Spec.Description = "new desc"
	got, err := reg.UpdateOrganizationUnit(context.Background(), "nic", "ministry-health", update)
	if err != nil {
		t.Fatalf("UpdateOrganizationUnit() error = %v", err)
	}
	if got.Metadata.DisplayName != "New Display" {
		t.Errorf("DisplayName = %q, want %q", got.Metadata.DisplayName, "New Display")
	}
	if got.Metadata.Labels["tier"] != "gold" {
		t.Errorf("Labels = %v, want tier=gold", got.Metadata.Labels)
	}
	if got.Metadata.Annotations["reviewed"] != "yes" {
		t.Errorf("Annotations = %v, want reviewed=yes", got.Metadata.Annotations)
	}
	if got.Spec.Description != "new desc" {
		t.Errorf("Description = %q, want %q", got.Spec.Description, "new desc")
	}
}

func TestUpdateOrganizationUnit_PreservesImmutableFields(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	original := sampleOU("nic", "ministry-health")
	original.Status = resources.OrganizationUnitStatus{Phase: resources.PhaseActive, Message: "ok"}
	_, _ = reg.CreateOrganizationUnit(context.Background(), original)

	// Attempt to change immutable/system fields via the update payload.
	update := sampleOU("nic", "ministry-health")
	update.Metadata.Name = "tampered"
	update.Spec.OrganizationName = "tampered-org"
	update.Status = resources.OrganizationUnitStatus{Phase: resources.PhaseFailed, Message: "hacked"}
	update.APIVersion = "tampered/v0"
	update.Kind = "Tampered"
	update.Spec.Description = "changed"

	got, err := reg.UpdateOrganizationUnit(context.Background(), "nic", "ministry-health", update)
	if err != nil {
		t.Fatalf("UpdateOrganizationUnit() error = %v", err)
	}
	if got.Metadata.Name != "ministry-health" {
		t.Errorf("Metadata.Name = %q, want ministry-health", got.Metadata.Name)
	}
	if got.Spec.OrganizationName != "nic" {
		t.Errorf("Spec.OrganizationName = %q, want nic", got.Spec.OrganizationName)
	}
	if got.Status.Phase != resources.PhaseActive || got.Status.Message != "ok" {
		t.Errorf("Status = %+v, want {Active ok}", got.Status)
	}
	if got.APIVersion != resources.OUAPIVersion {
		t.Errorf("APIVersion = %q, want %q", got.APIVersion, resources.OUAPIVersion)
	}
	if got.Kind != resources.OUKind {
		t.Errorf("Kind = %q, want %q", got.Kind, resources.OUKind)
	}
	if got.Spec.Description != "changed" {
		t.Errorf("Description = %q, want changed", got.Spec.Description)
	}
}

func TestUpdateOrganizationUnit_NotFound(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	_, err := reg.UpdateOrganizationUnit(context.Background(), "nic", "missing", sampleOU("nic", "missing"))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestDeleteOrganizationUnit_Exists(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	_, _ = reg.CreateOrganizationUnit(context.Background(), sampleOU("nic", "ministry-health"))
	if err := reg.DeleteOrganizationUnit(context.Background(), "nic", "ministry-health"); err != nil {
		t.Fatalf("DeleteOrganizationUnit() error = %v", err)
	}
	_, err := reg.GetOrganizationUnit(context.Background(), "nic", "ministry-health")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound after delete", err)
	}
}

func TestDeleteOrganizationUnit_NotFound(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	err := reg.DeleteOrganizationUnit(context.Background(), "nic", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestCountByOrganization(t *testing.T) {
	reg := NewOrganizationUnitRegistry()
	_, _ = reg.CreateOrganizationUnit(context.Background(), sampleOU("nic", "ministry-health"))
	_, _ = reg.CreateOrganizationUnit(context.Background(), sampleOU("nic", "ministry-finance"))
	_, _ = reg.CreateOrganizationUnit(context.Background(), sampleOU("state-gov", "ministry-health"))

	count, err := reg.CountByOrganization(context.Background(), "nic")
	if err != nil {
		t.Fatalf("CountByOrganization() error = %v", err)
	}
	if count != 2 {
		t.Errorf("count for nic = %d, want 2", count)
	}

	count, err = reg.CountByOrganization(context.Background(), "state-gov")
	if err != nil {
		t.Fatalf("CountByOrganization() error = %v", err)
	}
	if count != 1 {
		t.Errorf("count for state-gov = %d, want 1", count)
	}

	count, err = reg.CountByOrganization(context.Background(), "unknown")
	if err != nil {
		t.Fatalf("CountByOrganization() error = %v", err)
	}
	if count != 0 {
		t.Errorf("count for unknown = %d, want 0", count)
	}
}
