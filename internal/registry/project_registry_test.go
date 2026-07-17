package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func sampleProject(orgName, ouName, tenantName, name string) resources.Project {
	return resources.Project{
		APIVersion: resources.ProjectAPIVersion,
		Kind:       resources.ProjectKind,
		Metadata: resources.Metadata{
			Name:        name,
			DisplayName: "Display",
			Labels:      map[string]string{"env": "test"},
			Annotations: map[string]string{"note": "x"},
		},
		Spec: resources.ProjectSpec{
			OrganizationName:     orgName,
			OrganizationUnitName: ouName,
			TenantName:           tenantName,
			Description:          "desc",
		},
		Status: resources.ProjectStatus{Phase: resources.PhaseActive},
	}
}

func TestCreateProject_Stores(t *testing.T) {
	reg := NewProjectRegistry()
	created, err := reg.CreateProject(context.Background(), sampleProject("nic", "ministry-health", "payments", "prod"))
	if err != nil {
		t.Fatalf("CreateProject() error = %v", err)
	}
	if created.Metadata.Name != "prod" {
		t.Errorf("Name = %q, want %q", created.Metadata.Name, "prod")
	}
	got, err := reg.GetProject(context.Background(), "nic", "ministry-health", "payments", "prod")
	if err != nil {
		t.Fatalf("GetProject() error = %v", err)
	}
	if got.Metadata.Name != "prod" || got.Spec.OrganizationName != "nic" ||
		got.Spec.OrganizationUnitName != "ministry-health" || got.Spec.TenantName != "payments" {
		t.Errorf("got %q/%q/%q/%q, want nic/ministry-health/payments/prod",
			got.Spec.OrganizationName, got.Spec.OrganizationUnitName, got.Spec.TenantName, got.Metadata.Name)
	}
}

func TestCreateProject_Duplicate(t *testing.T) {
	reg := NewProjectRegistry()
	project := sampleProject("nic", "ministry-health", "payments", "prod")
	project.Spec.Description = "original"
	if _, err := reg.CreateProject(context.Background(), project); err != nil {
		t.Fatalf("first CreateProject() error = %v", err)
	}
	dup := sampleProject("nic", "ministry-health", "payments", "prod")
	dup.Spec.Description = "changed"
	_, err := reg.CreateProject(context.Background(), dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
	got, err := reg.GetProject(context.Background(), "nic", "ministry-health", "payments", "prod")
	if err != nil {
		t.Fatalf("GetProject() error = %v", err)
	}
	if got.Spec.Description != "original" {
		t.Errorf("Description = %q, want original (unchanged)", got.Spec.Description)
	}
}

func TestCreateProject_SameNameDifferentTenants(t *testing.T) {
	reg := NewProjectRegistry()
	if _, err := reg.CreateProject(context.Background(), sampleProject("nic", "ministry-health", "payments", "prod")); err != nil {
		t.Fatalf("create nic/ministry-health/payments/prod: %v", err)
	}
	if _, err := reg.CreateProject(context.Background(), sampleProject("nic", "ministry-health", "billing", "prod")); err != nil {
		t.Fatalf("create nic/ministry-health/billing/prod: %v", err)
	}
	if _, err := reg.CreateProject(context.Background(), sampleProject("nic", "ministry-finance", "payments", "prod")); err != nil {
		t.Fatalf("create nic/ministry-finance/payments/prod: %v", err)
	}
	if _, err := reg.CreateProject(context.Background(), sampleProject("state-gov", "ministry-health", "payments", "prod")); err != nil {
		t.Fatalf("create state-gov/ministry-health/payments/prod: %v", err)
	}
	if _, err := reg.GetProject(context.Background(), "nic", "ministry-health", "payments", "prod"); err != nil {
		t.Errorf("get nic/ministry-health/payments/prod: %v", err)
	}
	if _, err := reg.GetProject(context.Background(), "nic", "ministry-health", "billing", "prod"); err != nil {
		t.Errorf("get nic/ministry-health/billing/prod: %v", err)
	}
	if _, err := reg.GetProject(context.Background(), "nic", "ministry-finance", "payments", "prod"); err != nil {
		t.Errorf("get nic/ministry-finance/payments/prod: %v", err)
	}
	if _, err := reg.GetProject(context.Background(), "state-gov", "ministry-health", "payments", "prod"); err != nil {
		t.Errorf("get state-gov/ministry-health/payments/prod: %v", err)
	}
}

func TestGetProject_ByCompositeKey(t *testing.T) {
	reg := NewProjectRegistry()
	_, _ = reg.CreateProject(context.Background(), sampleProject("nic", "ministry-health", "payments", "prod"))
	got, err := reg.GetProject(context.Background(), "nic", "ministry-health", "payments", "prod")
	if err != nil {
		t.Fatalf("GetProject() error = %v", err)
	}
	if got.Spec.Description != "desc" {
		t.Errorf("Description = %q, want %q", got.Spec.Description, "desc")
	}
}

func TestGetProject_NotFound(t *testing.T) {
	reg := NewProjectRegistry()
	_, err := reg.GetProject(context.Background(), "nic", "ministry-health", "payments", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestListProjects_Empty(t *testing.T) {
	reg := NewProjectRegistry()
	items, err := reg.ListProjects(context.Background())
	if err != nil {
		t.Fatalf("ListProjects() error = %v", err)
	}
	if items == nil {
		t.Fatal("got nil slice, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Fatalf("got %d items, want 0", len(items))
	}
}

func TestListProjects_Sorted(t *testing.T) {
	reg := NewProjectRegistry()
	inputs := []struct{ org, ou, tenant, name string }{
		{"zebra", "unit-b", "tenant-a", "beta"},
		{"alpha", "unit-b", "tenant-a", "delta"},
		{"alpha", "unit-a", "tenant-b", "charlie"},
		{"alpha", "unit-a", "tenant-a", "bravo"},
		{"alpha", "unit-a", "tenant-a", "alpha"},
		{"nic", "ministry-health", "payments", "prod"},
	}
	for _, in := range inputs {
		if _, err := reg.CreateProject(context.Background(), sampleProject(in.org, in.ou, in.tenant, in.name)); err != nil {
			t.Fatalf("create %s/%s/%s/%s: %v", in.org, in.ou, in.tenant, in.name, err)
		}
	}
	items, err := reg.ListProjects(context.Background())
	if err != nil {
		t.Fatalf("ListProjects() error = %v", err)
	}
	if len(items) != 6 {
		t.Fatalf("got %d items, want 6", len(items))
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
			prev.Spec.TenantName > cur.Spec.TenantName {
			t.Fatalf("not sorted by tenantName within org/ou: %v", items)
		}
		if prev.Spec.OrganizationName == cur.Spec.OrganizationName &&
			prev.Spec.OrganizationUnitName == cur.Spec.OrganizationUnitName &&
			prev.Spec.TenantName == cur.Spec.TenantName &&
			prev.Metadata.Name >= cur.Metadata.Name {
			t.Fatalf("not sorted by name within org/ou/tenant: %v", items)
		}
	}
	first := items[0]
	if first.Spec.OrganizationName != "alpha" || first.Spec.OrganizationUnitName != "unit-a" ||
		first.Spec.TenantName != "tenant-a" || first.Metadata.Name != "alpha" {
		t.Errorf("first item = %s/%s/%s/%s, want alpha/unit-a/tenant-a/alpha",
			first.Spec.OrganizationName, first.Spec.OrganizationUnitName, first.Spec.TenantName, first.Metadata.Name)
	}
}

func TestUpdateProject_MutableFields(t *testing.T) {
	reg := NewProjectRegistry()
	_, _ = reg.CreateProject(context.Background(), sampleProject("nic", "ministry-health", "payments", "prod"))
	update := sampleProject("nic", "ministry-health", "payments", "prod")
	update.Metadata.DisplayName = "New Display"
	update.Metadata.Labels = map[string]string{"tier": "gold"}
	update.Metadata.Annotations = map[string]string{"reviewed": "yes"}
	update.Spec.Description = "new desc"
	got, err := reg.UpdateProject(context.Background(), update)
	if err != nil {
		t.Fatalf("UpdateProject() error = %v", err)
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

func TestUpdateProject_PreservesImmutableFields(t *testing.T) {
	reg := NewProjectRegistry()
	original := sampleProject("nic", "ministry-health", "payments", "prod")
	original.Status = resources.ProjectStatus{Phase: resources.PhaseActive, Message: "ok"}
	_, _ = reg.CreateProject(context.Background(), original)

	// Attempt to change immutable/system fields via the update payload while
	// keeping the composite-key fields intact so the lookup still matches.
	update := sampleProject("nic", "ministry-health", "payments", "prod")
	update.Status = resources.ProjectStatus{Phase: resources.PhaseFailed, Message: "hacked"}
	update.APIVersion = "tampered/v0"
	update.Kind = "Tampered"
	update.Spec.Description = "changed"

	got, err := reg.UpdateProject(context.Background(), update)
	if err != nil {
		t.Fatalf("UpdateProject() error = %v", err)
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
	if got.Spec.TenantName != "payments" {
		t.Errorf("Spec.TenantName = %q, want payments", got.Spec.TenantName)
	}
	if got.Status.Phase != resources.PhaseActive || got.Status.Message != "ok" {
		t.Errorf("Status = %+v, want {Active ok}", got.Status)
	}
	if got.APIVersion != resources.ProjectAPIVersion {
		t.Errorf("APIVersion = %q, want %q", got.APIVersion, resources.ProjectAPIVersion)
	}
	if got.Kind != resources.ProjectKind {
		t.Errorf("Kind = %q, want %q", got.Kind, resources.ProjectKind)
	}
	if got.Spec.Description != "changed" {
		t.Errorf("Description = %q, want changed", got.Spec.Description)
	}
}

func TestUpdateProject_NotFound(t *testing.T) {
	reg := NewProjectRegistry()
	_, err := reg.UpdateProject(context.Background(), sampleProject("nic", "ministry-health", "payments", "missing"))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestDeleteProject_Exists(t *testing.T) {
	reg := NewProjectRegistry()
	_, _ = reg.CreateProject(context.Background(), sampleProject("nic", "ministry-health", "payments", "prod"))
	if err := reg.DeleteProject(context.Background(), "nic", "ministry-health", "payments", "prod"); err != nil {
		t.Fatalf("DeleteProject() error = %v", err)
	}
	_, err := reg.GetProject(context.Background(), "nic", "ministry-health", "payments", "prod")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound after delete", err)
	}
}

func TestDeleteProject_NotFound(t *testing.T) {
	reg := NewProjectRegistry()
	err := reg.DeleteProject(context.Background(), "nic", "ministry-health", "payments", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestCountByTenant(t *testing.T) {
	reg := NewProjectRegistry()
	_, _ = reg.CreateProject(context.Background(), sampleProject("nic", "ministry-health", "payments", "prod"))
	_, _ = reg.CreateProject(context.Background(), sampleProject("nic", "ministry-health", "payments", "staging"))
	_, _ = reg.CreateProject(context.Background(), sampleProject("nic", "ministry-health", "billing", "prod"))
	_, _ = reg.CreateProject(context.Background(), sampleProject("nic", "ministry-finance", "payments", "prod"))
	_, _ = reg.CreateProject(context.Background(), sampleProject("state-gov", "ministry-health", "payments", "prod"))

	count, err := reg.CountByTenant(context.Background(), "nic", "ministry-health", "payments")
	if err != nil {
		t.Fatalf("CountByTenant() error = %v", err)
	}
	if count != 2 {
		t.Errorf("count for nic/ministry-health/payments = %d, want 2", count)
	}

	count, err = reg.CountByTenant(context.Background(), "nic", "ministry-health", "billing")
	if err != nil {
		t.Fatalf("CountByTenant() error = %v", err)
	}
	if count != 1 {
		t.Errorf("count for nic/ministry-health/billing = %d, want 1", count)
	}

	count, err = reg.CountByTenant(context.Background(), "nic", "ministry-health", "unknown")
	if err != nil {
		t.Fatalf("CountByTenant() error = %v", err)
	}
	if count != 0 {
		t.Errorf("count for nic/ministry-health/unknown = %d, want 0", count)
	}
}
