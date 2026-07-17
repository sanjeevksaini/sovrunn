package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// emissionTestEnv wires all four mutating handlers to a shared, real
// OperationRegistry-backed emitter so emission can be asserted end-to-end.
type emissionTestEnv struct {
	orgRegistry       *registry.OrganizationRegistry
	ouRegistry        *registry.OrganizationUnitRegistry
	tenantRegistry    *registry.TenantRegistry
	projectRegistry   *registry.ProjectRegistry
	operationRegistry *registry.OperationRegistry
	emitter           OperationEmitter
	orgHandler        *OrgHandler
	ouHandler         *OUHandler
	tenantHandler     *TenantHandler
	projectHandler    *ProjectHandler
}

func newEmissionTestEnv(t *testing.T) *emissionTestEnv {
	t.Helper()
	orgRegistry := registry.NewOrganizationRegistry()
	ouRegistry := registry.NewOrganizationUnitRegistry()
	tenantRegistry := registry.NewTenantRegistry()
	projectRegistry := registry.NewProjectRegistry()
	operationRegistry := registry.NewOperationRegistry()
	emitter := NewRegistryEmitter(operationRegistry, nil)

	ouBlocker := registry.NewOUChildBlockerChecker(ouRegistry)
	tenantBlocker := registry.NewTenantChildBlockerChecker(tenantRegistry)
	projectBlocker := registry.NewProjectChildBlockerChecker(projectRegistry)

	return &emissionTestEnv{
		orgRegistry:       orgRegistry,
		ouRegistry:        ouRegistry,
		tenantRegistry:    tenantRegistry,
		projectRegistry:   projectRegistry,
		operationRegistry: operationRegistry,
		emitter:           emitter,
		orgHandler:        NewOrgHandler(orgRegistry, ouBlocker, emitter),
		ouHandler:         NewOUHandler(ouRegistry, orgRegistry, tenantBlocker, emitter),
		tenantHandler:     NewTenantHandler(tenantRegistry, ouRegistry, projectBlocker, emitter),
		projectHandler:    NewProjectHandler(projectRegistry, tenantRegistry, emitter),
	}
}

func recordedOperations(t *testing.T, env *emissionTestEnv) []resources.Operation {
	t.Helper()
	items, err := env.operationRegistry.ListOperations(context.Background())
	if err != nil {
		t.Fatalf("ListOperations() error = %v", err)
	}
	return items
}

func requireSingleOperation(t *testing.T, env *emissionTestEnv) resources.Operation {
	t.Helper()
	items := recordedOperations(t, env)
	if len(items) != 1 {
		t.Fatalf("recorded operations = %d, want exactly 1", len(items))
	}
	return items[0]
}

// seedProjectDirect stores a Project directly in the registry (no emission),
// used to prepare Update/Delete cases without polluting the OperationRegistry.
func seedProjectDirect(t *testing.T, reg *registry.ProjectRegistry, orgName, ouName, tenantName, name string) {
	t.Helper()
	project := resources.Project{
		APIVersion: resources.ProjectAPIVersion,
		Kind:       resources.ProjectKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.ProjectSpec{
			OrganizationName:     orgName,
			OrganizationUnitName: ouName,
			TenantName:           tenantName,
		},
		Status: resources.ProjectStatus{Phase: resources.PhaseActive},
	}
	if _, err := reg.CreateProject(context.Background(), project); err != nil {
		t.Fatalf("seedProjectDirect(%s/%s/%s/%s): %v", orgName, ouName, tenantName, name, err)
	}
}

func createOrg(h *OrgHandler, name, desc string) *httptest.ResponseRecorder {
	body := map[string]any{
		"metadata": map[string]any{"name": name},
		"spec":     map[string]any{"description": desc},
	}
	req := jsonRequest(http.MethodPost, "/v1/organizations", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	return rec
}

func assertOperation(t *testing.T, op resources.Operation, want resources.OperationSpec) {
	t.Helper()
	if op.APIVersion != resources.OperationAPIVersion {
		t.Errorf("APIVersion = %q, want %q", op.APIVersion, resources.OperationAPIVersion)
	}
	if op.Kind != resources.OperationKind {
		t.Errorf("Kind = %q, want %q", op.Kind, resources.OperationKind)
	}
	if op.Metadata.Name == "" {
		t.Error("Metadata.Name is empty, want generated ID")
	}
	if op.Spec.Type != want.Type {
		t.Errorf("Spec.Type = %q, want %q", op.Spec.Type, want.Type)
	}
	if op.Spec.ResourceKind != want.ResourceKind {
		t.Errorf("Spec.ResourceKind = %q, want %q", op.Spec.ResourceKind, want.ResourceKind)
	}
	if op.Spec.ResourceName != want.ResourceName {
		t.Errorf("Spec.ResourceName = %q, want %q", op.Spec.ResourceName, want.ResourceName)
	}
	if op.Spec.OrganizationName != want.OrganizationName {
		t.Errorf("Spec.OrganizationName = %q, want %q", op.Spec.OrganizationName, want.OrganizationName)
	}
	if op.Spec.OrganizationUnitName != want.OrganizationUnitName {
		t.Errorf("Spec.OrganizationUnitName = %q, want %q", op.Spec.OrganizationUnitName, want.OrganizationUnitName)
	}
	if op.Spec.TenantName != want.TenantName {
		t.Errorf("Spec.TenantName = %q, want %q", op.Spec.TenantName, want.TenantName)
	}
	if op.Spec.ProjectName != want.ProjectName {
		t.Errorf("Spec.ProjectName = %q, want %q", op.Spec.ProjectName, want.ProjectName)
	}
	if op.Spec.Actor != "system" {
		t.Errorf("Spec.Actor = %q, want system", op.Spec.Actor)
	}
	// jsonRequest injects request ID "test-req-id" into the request context.
	if op.Spec.RequestID != "test-req-id" {
		t.Errorf("Spec.RequestID = %q, want test-req-id", op.Spec.RequestID)
	}
	if op.Status.Phase != resources.OperationPhaseSucceeded {
		t.Errorf("Status.Phase = %q, want Succeeded", op.Status.Phase)
	}
	if op.Status.CreatedAt == "" {
		t.Error("Status.CreatedAt is empty")
	}
	if _, err := time.Parse(time.RFC3339, op.Status.CreatedAt); err != nil {
		t.Errorf("Status.CreatedAt %q not RFC3339: %v", op.Status.CreatedAt, err)
	}
	if op.Status.UpdatedAt != op.Status.CreatedAt {
		t.Errorf("Status.UpdatedAt = %q, want == CreatedAt %q", op.Status.UpdatedAt, op.Status.CreatedAt)
	}
	if op.Status.CompletedAt != op.Status.CreatedAt {
		t.Errorf("Status.CompletedAt = %q, want == CreatedAt %q", op.Status.CompletedAt, op.Status.CreatedAt)
	}
}

func TestEmission_AllMutatingPaths(t *testing.T) {
	orgUpdateBody := map[string]any{
		"metadata": map[string]any{"name": "nic"},
		"spec":     map[string]any{"description": "updated"},
	}
	ouUpdateBody := map[string]any{
		"metadata": map[string]any{"name": "ou1"},
		"spec":     map[string]any{"organizationName": "nic", "description": "updated"},
	}
	tenantUpdateBody := map[string]any{
		"metadata": map[string]any{"name": "t1"},
		"spec": map[string]any{
			"organizationName":     "nic",
			"organizationUnitName": "ou1",
			"description":          "updated",
		},
	}
	projectUpdateBody := map[string]any{
		"metadata": map[string]any{"name": "p1"},
		"spec": map[string]any{
			"organizationName":     "nic",
			"organizationUnitName": "ou1",
			"tenantName":           "t1",
			"description":          "updated",
		},
	}

	cases := []struct {
		name       string
		setup      func(env *emissionTestEnv)
		action     func(env *emissionTestEnv) *httptest.ResponseRecorder
		wantStatus int
		wantSpec   resources.OperationSpec
	}{
		{
			name:  "Organization Create",
			setup: func(env *emissionTestEnv) {},
			action: func(env *emissionTestEnv) *httptest.ResponseRecorder {
				return createOrg(env.orgHandler, "nic", "")
			},
			wantStatus: http.StatusCreated,
			wantSpec: resources.OperationSpec{
				Type: resources.OpCreateOrganization, ResourceKind: resources.OrganizationKind,
				ResourceName: "nic", OrganizationName: "nic",
			},
		},
		{
			name: "Organization Update",
			setup: func(env *emissionTestEnv) {
				seedOrg(t, env.orgRegistry, "nic")
			},
			action: func(env *emissionTestEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodPut, "/v1/organizations/nic", orgUpdateBody, "application/json")
				rec := httptest.NewRecorder()
				env.orgHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusOK,
			wantSpec: resources.OperationSpec{
				Type: resources.OpUpdateOrganization, ResourceKind: resources.OrganizationKind,
				ResourceName: "nic", OrganizationName: "nic",
			},
		},
		{
			name: "Organization Delete",
			setup: func(env *emissionTestEnv) {
				seedOrg(t, env.orgRegistry, "nic")
			},
			action: func(env *emissionTestEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodDelete, "/v1/organizations/nic", nil, "")
				rec := httptest.NewRecorder()
				env.orgHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusNoContent,
			wantSpec: resources.OperationSpec{
				Type: resources.OpDeleteOrganization, ResourceKind: resources.OrganizationKind,
				ResourceName: "nic", OrganizationName: "nic",
			},
		},
		{
			name: "OrganizationUnit Create",
			setup: func(env *emissionTestEnv) {
				seedOrg(t, env.orgRegistry, "nic")
			},
			action: func(env *emissionTestEnv) *httptest.ResponseRecorder {
				return createOU(env.ouHandler, "nic", "ou1", "")
			},
			wantStatus: http.StatusCreated,
			wantSpec: resources.OperationSpec{
				Type: resources.OpCreateOrganizationUnit, ResourceKind: resources.OrganizationUnitKind,
				ResourceName: "ou1", OrganizationName: "nic", OrganizationUnitName: "ou1",
			},
		},
		{
			name: "OrganizationUnit Update",
			setup: func(env *emissionTestEnv) {
				seedOrg(t, env.orgRegistry, "nic")
				seedOU(t, env.ouRegistry, "nic", "ou1")
			},
			action: func(env *emissionTestEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodPut, "/v1/organization-units/nic/ou1", ouUpdateBody, "application/json")
				rec := httptest.NewRecorder()
				env.ouHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusOK,
			wantSpec: resources.OperationSpec{
				Type: resources.OpUpdateOrganizationUnit, ResourceKind: resources.OrganizationUnitKind,
				ResourceName: "ou1", OrganizationName: "nic", OrganizationUnitName: "ou1",
			},
		},
		{
			name: "OrganizationUnit Delete",
			setup: func(env *emissionTestEnv) {
				seedOrg(t, env.orgRegistry, "nic")
				seedOU(t, env.ouRegistry, "nic", "ou1")
			},
			action: func(env *emissionTestEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodDelete, "/v1/organization-units/nic/ou1", nil, "")
				rec := httptest.NewRecorder()
				env.ouHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusNoContent,
			wantSpec: resources.OperationSpec{
				Type: resources.OpDeleteOrganizationUnit, ResourceKind: resources.OrganizationUnitKind,
				ResourceName: "ou1", OrganizationName: "nic", OrganizationUnitName: "ou1",
			},
		},
		{
			name: "Tenant Create",
			setup: func(env *emissionTestEnv) {
				seedOrg(t, env.orgRegistry, "nic")
				seedOU(t, env.ouRegistry, "nic", "ou1")
			},
			action: func(env *emissionTestEnv) *httptest.ResponseRecorder {
				return createTenant(env.tenantHandler, "nic", "ou1", "t1", "")
			},
			wantStatus: http.StatusCreated,
			wantSpec: resources.OperationSpec{
				Type: resources.OpCreateTenant, ResourceKind: resources.TenantKind,
				ResourceName: "t1", OrganizationName: "nic", OrganizationUnitName: "ou1", TenantName: "t1",
			},
		},
		{
			name: "Tenant Update",
			setup: func(env *emissionTestEnv) {
				seedOrg(t, env.orgRegistry, "nic")
				seedOU(t, env.ouRegistry, "nic", "ou1")
				seedTenant(t, env.tenantRegistry, "nic", "ou1", "t1")
			},
			action: func(env *emissionTestEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodPut, "/v1/tenants/nic/ou1/t1", tenantUpdateBody, "application/json")
				rec := httptest.NewRecorder()
				env.tenantHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusOK,
			wantSpec: resources.OperationSpec{
				Type: resources.OpUpdateTenant, ResourceKind: resources.TenantKind,
				ResourceName: "t1", OrganizationName: "nic", OrganizationUnitName: "ou1", TenantName: "t1",
			},
		},
		{
			name: "Tenant Delete",
			setup: func(env *emissionTestEnv) {
				seedOrg(t, env.orgRegistry, "nic")
				seedOU(t, env.ouRegistry, "nic", "ou1")
				seedTenant(t, env.tenantRegistry, "nic", "ou1", "t1")
			},
			action: func(env *emissionTestEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodDelete, "/v1/tenants/nic/ou1/t1", nil, "")
				rec := httptest.NewRecorder()
				env.tenantHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusNoContent,
			wantSpec: resources.OperationSpec{
				Type: resources.OpDeleteTenant, ResourceKind: resources.TenantKind,
				ResourceName: "t1", OrganizationName: "nic", OrganizationUnitName: "ou1", TenantName: "t1",
			},
		},
		{
			name: "Project Create",
			setup: func(env *emissionTestEnv) {
				seedOrg(t, env.orgRegistry, "nic")
				seedOU(t, env.ouRegistry, "nic", "ou1")
				seedTenant(t, env.tenantRegistry, "nic", "ou1", "t1")
			},
			action: func(env *emissionTestEnv) *httptest.ResponseRecorder {
				return createProject(env.projectHandler, "nic", "ou1", "t1", "p1", "")
			},
			wantStatus: http.StatusCreated,
			wantSpec: resources.OperationSpec{
				Type: resources.OpCreateProject, ResourceKind: resources.ProjectKind,
				ResourceName: "p1", OrganizationName: "nic", OrganizationUnitName: "ou1", TenantName: "t1", ProjectName: "p1",
			},
		},
		{
			name: "Project Update",
			setup: func(env *emissionTestEnv) {
				seedOrg(t, env.orgRegistry, "nic")
				seedOU(t, env.ouRegistry, "nic", "ou1")
				seedTenant(t, env.tenantRegistry, "nic", "ou1", "t1")
				seedProjectDirect(t, env.projectRegistry, "nic", "ou1", "t1", "p1")
			},
			action: func(env *emissionTestEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodPut, "/v1/projects/nic/ou1/t1/p1", projectUpdateBody, "application/json")
				rec := httptest.NewRecorder()
				env.projectHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusOK,
			wantSpec: resources.OperationSpec{
				Type: resources.OpUpdateProject, ResourceKind: resources.ProjectKind,
				ResourceName: "p1", OrganizationName: "nic", OrganizationUnitName: "ou1", TenantName: "t1", ProjectName: "p1",
			},
		},
		{
			name: "Project Delete",
			setup: func(env *emissionTestEnv) {
				seedOrg(t, env.orgRegistry, "nic")
				seedOU(t, env.ouRegistry, "nic", "ou1")
				seedTenant(t, env.tenantRegistry, "nic", "ou1", "t1")
				seedProjectDirect(t, env.projectRegistry, "nic", "ou1", "t1", "p1")
			},
			action: func(env *emissionTestEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodDelete, "/v1/projects/nic/ou1/t1/p1", nil, "")
				rec := httptest.NewRecorder()
				env.projectHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusNoContent,
			wantSpec: resources.OperationSpec{
				Type: resources.OpDeleteProject, ResourceKind: resources.ProjectKind,
				ResourceName: "p1", OrganizationName: "nic", OrganizationUnitName: "ou1", TenantName: "t1", ProjectName: "p1",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := newEmissionTestEnv(t)
			tc.setup(env)

			rec := tc.action(env)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}

			op := requireSingleOperation(t, env)
			assertOperation(t, op, tc.wantSpec)
		})
	}
}

func TestEmission_FailureIsolation(t *testing.T) {
	// A stub emitter that always errors must not affect the primary response.
	newHandlers := func() *emissionTestEnv {
		orgRegistry := registry.NewOrganizationRegistry()
		ouRegistry := registry.NewOrganizationUnitRegistry()
		ouBlocker := registry.NewOUChildBlockerChecker(ouRegistry)
		return &emissionTestEnv{
			orgRegistry: orgRegistry,
			ouRegistry:  ouRegistry,
			orgHandler:  NewOrgHandler(orgRegistry, ouBlocker, &failingEmitter{}),
		}
	}

	t.Run("create still returns 201", func(t *testing.T) {
		env := newHandlers()
		rec := createOrg(env.orgHandler, "nic", "")
		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
		}
	})

	t.Run("update still returns 200", func(t *testing.T) {
		env := newHandlers()
		seedOrg(t, env.orgRegistry, "nic")
		body := map[string]any{
			"metadata": map[string]any{"name": "nic"},
			"spec":     map[string]any{"description": "updated"},
		}
		req := jsonRequest(http.MethodPut, "/v1/organizations/nic", body, "application/json")
		rec := httptest.NewRecorder()
		env.orgHandler.HandleItem(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
	})

	t.Run("delete still returns 204", func(t *testing.T) {
		env := newHandlers()
		seedOrg(t, env.orgRegistry, "nic")
		req := jsonRequest(http.MethodDelete, "/v1/organizations/nic", nil, "")
		rec := httptest.NewRecorder()
		env.orgHandler.HandleItem(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
		}
	})
}

func TestEmission_NoEmissionOnFailedAction(t *testing.T) {
	t.Run("invalid create payload", func(t *testing.T) {
		env := newEmissionTestEnv(t)
		rec := createOrg(env.orgHandler, "INVALID_NAME", "")
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		if ops := recordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0 on validation failure", len(ops))
		}
	})

	t.Run("duplicate create", func(t *testing.T) {
		env := newEmissionTestEnv(t)
		seedOrg(t, env.orgRegistry, "nic")
		rec := createOrg(env.orgHandler, "nic", "")
		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", rec.Code)
		}
		if ops := recordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0 on duplicate create", len(ops))
		}
	})

	t.Run("update missing target", func(t *testing.T) {
		env := newEmissionTestEnv(t)
		body := map[string]any{
			"metadata": map[string]any{"name": "nic"},
			"spec":     map[string]any{"description": "updated"},
		}
		req := jsonRequest(http.MethodPut, "/v1/organizations/nic", body, "application/json")
		rec := httptest.NewRecorder()
		env.orgHandler.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ops := recordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0 on update miss", len(ops))
		}
	})

	t.Run("delete missing target", func(t *testing.T) {
		env := newEmissionTestEnv(t)
		req := jsonRequest(http.MethodDelete, "/v1/organizations/nic", nil, "")
		rec := httptest.NewRecorder()
		env.orgHandler.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ops := recordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0 on delete miss", len(ops))
		}
	})

	t.Run("blocked delete", func(t *testing.T) {
		env := newEmissionTestEnv(t)
		seedOrg(t, env.orgRegistry, "nic")
		seedOU(t, env.ouRegistry, "nic", "ou1")
		req := jsonRequest(http.MethodDelete, "/v1/organizations/nic", nil, "")
		rec := httptest.NewRecorder()
		env.orgHandler.HandleItem(rec, req)
		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
		}
		if ops := recordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0 on blocked delete", len(ops))
		}
	})
}
