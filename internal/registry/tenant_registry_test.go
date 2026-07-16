package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func sampleTenant(orgName, ouName, name string) resources.Tenant {
	return resources.Tenant{
		APIVersion: resources.TenantAPIVersion,
		Kind:       resources.TenantKind,
		Metadata: resources.Metadata{
			Name:        name,
			DisplayName: "Display",
			Labels:      map[string]string{"env": "test"},
			Annotations: map[string]string{"note": "x"},
		},
		Spec: resources.TenantSpec{
			OrganizationName:     orgName,
			OrganizationUnitName: ouName,
			Description:          "desc",
		},
		Status: resources.TenantStatus{Phase: resources.PhaseActive},
	}
}

func TestCreateTenant_Stores(t *testing.T) {
	reg := NewTenantRegistry()
	created, err := reg.CreateTenant(context.Background(), sampleTenant("nic", "ministry-health", "prod"))
	if err != nil {
		t.Fatalf("CreateTenant() error = %v", err)
	}
	if created.Metadata.Name != "prod" {
		t.Errorf("Name = %q, want %q", created.Metadata.Name, "prod")
	}
	got, err := reg.GetTenant(context.Background(), "nic", "ministry-health", "prod")
	if err != nil {
		t.Fatalf("GetTenant() error = %v", err)
	}
	if got.Metadata.Name != "prod" || got.Spec.OrganizationName != "nic" || got.Spec.OrganizationUnitName != "ministry-health" {
		t.Errorf("got %q/%q/%q, want nic/ministry-health/prod",
			got.Spec.OrganizationName, got.Spec.OrganizationUnitName, got.Metadata.Name)
	}
}

func TestCreateTenant_Duplicate(t *testing.T) {
	reg := NewTenantRegistry()
	tnt := sampleTenant("nic", "ministry-health", "prod")
	tnt.Spec.Description = "original"
	if _, err := reg.CreateTenant(context.Background(), tnt); err != nil {
		t.Fatalf("first CreateTenant() error = %v", err)
	}
	dup := sampleTenant("nic", "ministry-health", "prod")
	dup.Spec.Description = "changed"
	_, err := reg.CreateTenant(context.Background(), dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
	got, err := reg.GetTenant(context.Background(), "nic", "ministry-health", "prod")
	if err != nil {
		t.Fatalf("GetTenant() error = %v", err)
	}
	if got.Spec.Description != "original" {
		t.Errorf("Description = %q, want original (unchanged)", got.Spec.Description)
	}
}

func TestCreateTenant_SameNameDifferentOUs(t *testing.T) {
	reg := NewTenantRegistry()
	if _, err := reg.CreateTenant(context.Background(), sampleTenant("nic", "ministry-health", "prod")); err != nil {
		t.Fatalf("create nic/ministry-health/prod: %v", err)
	}
	if _, err := reg.CreateTenant(context.Background(), sampleTenant("nic", "ministry-finance", "prod")); err != nil {
		t.Fatalf("create nic/ministry-finance/prod: %v", err)
	}
	if _, err := reg.CreateTenant(context.Background(), sampleTenant("state-gov", "ministry-health", "prod")); err != nil {
		t.Fatalf("create state-gov/ministry-health/prod: %v", err)
	}
	if _, err := reg.GetTenant(context.Background(), "nic", "ministry-health", "prod"); err != nil {
		t.Errorf("get nic/ministry-health/prod: %v", err)
	}
	if _, err := reg.GetTenant(context.Background(), "nic", "ministry-finance", "prod"); err != nil {
		t.Errorf("get nic/ministry-finance/prod: %v", err)
	}
	if _, err := reg.GetTenant(context.Background(), "state-gov", "ministry-health", "prod"); err != nil {
		t.Errorf("get state-gov/ministry-health/prod: %v", err)
	}
}

func TestGetTenant_ByCompositeKey(t *testing.T) {
	reg := NewTenantRegistry()
	_, _ = reg.CreateTenant(context.Background(), sampleTenant("nic", "ministry-health", "prod"))
	got, err := reg.GetTenant(context.Background(), "nic", "ministry-health", "prod")
	if err != nil {
		t.Fatalf("GetTenant() error = %v", err)
	}
	if got.Spec.Description != "desc" {
		t.Errorf("Description = %q, want %q", got.Spec.Description, "desc")
	}
}

func TestGetTenant_NotFound(t *testing.T) {
	reg := NewTenantRegistry()
	_, err := reg.GetTenant(context.Background(), "nic", "ministry-health", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestListTenants_Empty(t *testing.T) {
	reg := NewTenantRegistry()
	items, err := reg.ListTenants(context.Background())
	if err != nil {
		t.Fatalf("ListTenants() error = %v", err)
	}
	if items == nil {
		t.Fatal("got nil slice, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Fatalf("got %d items, want 0", len(items))
	}
}

func TestListTenants_Sorted(t *testing.T) {
	reg := NewTenantRegistry()
	inputs := []struct{ org, ou, name string }{
		{"zebra", "unit-b", "beta"},
		{"alpha", "unit-b", "delta"},
		{"alpha", "unit-a", "charlie"},
		{"alpha", "unit-a", "bravo"},
		{"nic", "ministry-health", "prod"},
	}
	for _, in := range inputs {
		if _, err := reg.CreateTenant(context.Background(), sampleTenant(in.org, in.ou, in.name)); err != nil {
			t.Fatalf("create %s/%s/%s: %v", in.org, in.ou, in.name, err)
		}
	}
	items, err := reg.ListTenants(context.Background())
	if err != nil {
		t.Fatalf("ListTenants() error = %v", err)
	}
	if len(items) != 5 {
		t.Fatalf("got %d items, want 5", len(items))
	}
	for i := 1; i < len(items); i++ {
		prev, cur := items[i-1], items[i]
		if prev.Spec.OrganizationName > cur.Spec.OrganizationName {
			t.Fatalf("not sorted by organizationName: %v", items)
		}
		if prev.Spec.OrganizationName == cur.Spec.OrganizationName &&
			prev.Spec.OrganizationUnitName > cur.Spec.OrganizationUnitName {
			t.Fatalf("not sorted by organizationUnitName within org: %v", items)
		}
		if prev.Spec.OrganizationName == cur.Spec.OrganizationName &&
			prev.Spec.OrganizationUnitName == cur.Spec.OrganizationUnitName &&
			prev.Metadata.Name >= cur.Metadata.Name {
			t.Fatalf("not sorted by name within org/ou: %v", items)
		}
	}
	first := items[0]
	if first.Spec.OrganizationName != "alpha" || first.Spec.OrganizationUnitName != "unit-a" || first.Metadata.Name != "bravo" {
		t.Errorf("first item = %s/%s/%s, want alpha/unit-a/bravo",
			first.Spec.OrganizationName, first.Spec.OrganizationUnitName, first.Metadata.Name)
	}
}

func TestUpdateTenant_MutableFields(t *testing.T) {
	reg := NewTenantRegistry()
	_, _ = reg.CreateTenant(context.Background(), sampleTenant("nic", "ministry-health", "prod"))
	update := sampleTenant("nic", "ministry-health", "prod")
	update.Metadata.DisplayName = "New Display"
	update.Metadata.Labels = map[string]string{"tier": "gold"}
	update.Metadata.Annotations = map[string]string{"reviewed": "yes"}
	update.Spec.Description = "new desc"
	got, err := reg.UpdateTenant(context.Background(), update)
	if err != nil {
		t.Fatalf("UpdateTenant() error = %v", err)
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

func TestUpdateTenant_PreservesImmutableFields(t *testing.T) {
	reg := NewTenantRegistry()
	original := sampleTenant("nic", "ministry-health", "prod")
	original.Status = resources.TenantStatus{Phase: resources.PhaseActive, Message: "ok"}
	_, _ = reg.CreateTenant(context.Background(), original)

	// Attempt to change immutable/system fields via the update payload while
	// keeping the composite-key fields intact so the lookup still matches.
	update := sampleTenant("nic", "ministry-health", "prod")
	update.Status = resources.TenantStatus{Phase: resources.PhaseFailed, Message: "hacked"}
	update.APIVersion = "tampered/v0"
	update.Kind = "Tampered"
	update.Spec.Description = "changed"

	got, err := reg.UpdateTenant(context.Background(), update)
	if err != nil {
		t.Fatalf("UpdateTenant() error = %v", err)
	}
	if got.Metadata.Name != "prod" {
		t.Errorf("Metadata.Name = %q, want prod", got.Metadata.Name)
	}
	if got.Spec.OrganizationName != "nic" {
		t.Errorf("Spec.OrganizationName = %q, want nic", got.Spec.OrganizationName)
	}
	if got.Spec.OrganizationUnitName != "ministry-health" {
		t.Errorf("Spec.OrganizationUnitName = %q, want ministry-health", got.Spec.OrganizationUnitName)
	}
	if got.Status.Phase != resources.PhaseActive || got.Status.Message != "ok" {
		t.Errorf("Status = %+v, want {Active ok}", got.Status)
	}
	if got.APIVersion != resources.TenantAPIVersion {
		t.Errorf("APIVersion = %q, want %q", got.APIVersion, resources.TenantAPIVersion)
	}
	if got.Kind != resources.TenantKind {
		t.Errorf("Kind = %q, want %q", got.Kind, resources.TenantKind)
	}
	if got.Spec.Description != "changed" {
		t.Errorf("Description = %q, want changed", got.Spec.Description)
	}
}

func TestUpdateTenant_NotFound(t *testing.T) {
	reg := NewTenantRegistry()
	_, err := reg.UpdateTenant(context.Background(), sampleTenant("nic", "ministry-health", "missing"))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestDeleteTenant_Exists(t *testing.T) {
	reg := NewTenantRegistry()
	_, _ = reg.CreateTenant(context.Background(), sampleTenant("nic", "ministry-health", "prod"))
	if err := reg.DeleteTenant(context.Background(), "nic", "ministry-health", "prod"); err != nil {
		t.Fatalf("DeleteTenant() error = %v", err)
	}
	_, err := reg.GetTenant(context.Background(), "nic", "ministry-health", "prod")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound after delete", err)
	}
}

func TestDeleteTenant_NotFound(t *testing.T) {
	reg := NewTenantRegistry()
	err := reg.DeleteTenant(context.Background(), "nic", "ministry-health", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestCountByOrganizationUnit(t *testing.T) {
	reg := NewTenantRegistry()
	_, _ = reg.CreateTenant(context.Background(), sampleTenant("nic", "ministry-health", "prod"))
	_, _ = reg.CreateTenant(context.Background(), sampleTenant("nic", "ministry-health", "staging"))
	_, _ = reg.CreateTenant(context.Background(), sampleTenant("nic", "ministry-finance", "prod"))
	_, _ = reg.CreateTenant(context.Background(), sampleTenant("state-gov", "ministry-health", "prod"))

	count, err := reg.CountByOrganizationUnit(context.Background(), "nic", "ministry-health")
	if err != nil {
		t.Fatalf("CountByOrganizationUnit() error = %v", err)
	}
	if count != 2 {
		t.Errorf("count for nic/ministry-health = %d, want 2", count)
	}

	count, err = reg.CountByOrganizationUnit(context.Background(), "nic", "ministry-finance")
	if err != nil {
		t.Fatalf("CountByOrganizationUnit() error = %v", err)
	}
	if count != 1 {
		t.Errorf("count for nic/ministry-finance = %d, want 1", count)
	}

	count, err = reg.CountByOrganizationUnit(context.Background(), "nic", "unknown")
	if err != nil {
		t.Fatalf("CountByOrganizationUnit() error = %v", err)
	}
	if count != 0 {
		t.Errorf("count for nic/unknown = %d, want 0", count)
	}
}
