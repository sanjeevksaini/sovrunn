package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// newTestProjectHandler builds a ProjectHandler backed by fresh registries and
// returns the handler along with the Tenant and Project registries so tests
// can seed parent Tenants directly. The TenantRegistry satisfies
// registry.TenantLookup.
func newTestProjectHandler() (*ProjectHandler, *registry.TenantRegistry, *registry.ProjectRegistry) {
	tenantRegistry := registry.NewTenantRegistry()
	projectRegistry := registry.NewProjectRegistry()
	handler := NewProjectHandler(projectRegistry, tenantRegistry, nil)
	return handler, tenantRegistry, projectRegistry
}

// seedTenant creates a parent Tenant directly in the registry.
func seedTenant(t *testing.T, reg *registry.TenantRegistry, orgName, ouName, name string) {
	t.Helper()
	tenant := resources.Tenant{
		APIVersion: resources.TenantAPIVersion,
		Kind:       resources.TenantKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.TenantSpec{
			OrganizationName:     orgName,
			OrganizationUnitName: ouName,
		},
		Status: resources.TenantStatus{Phase: resources.PhaseActive},
	}
	if _, err := reg.CreateTenant(context.Background(), tenant); err != nil {
		t.Fatalf("seedTenant(%s/%s/%s): %v", orgName, ouName, name, err)
	}
}

// createProject issues a POST create through the handler and returns the recorder.
func createProject(h *ProjectHandler, orgName, ouName, tenantName, name, desc string) *httptest.ResponseRecorder {
	body := map[string]any{
		"metadata": map[string]any{"name": name},
		"spec": map[string]any{
			"organizationName":     orgName,
			"organizationUnitName": ouName,
			"tenantName":           tenantName,
			"description":          desc,
		},
	}
	req := jsonRequest(http.MethodPost, "/v1/projects", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	return rec
}

func TestProjectHandler_Create_Valid(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	rec := createProject(h, "nic", "ministry-health", "payments", "prod", "Production")
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var project resources.Project
	if err := json.NewDecoder(rec.Body).Decode(&project); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if project.APIVersion != resources.ProjectAPIVersion || project.Kind != resources.ProjectKind {
		t.Errorf("apiVersion/kind not set by server: %+v", project)
	}
	if project.Status.Phase != resources.PhaseActive {
		t.Errorf("phase = %q, want Active", project.Status.Phase)
	}
	if project.Metadata.Name != "prod" || project.Spec.OrganizationName != "nic" ||
		project.Spec.OrganizationUnitName != "ministry-health" || project.Spec.TenantName != "payments" {
		t.Errorf("unexpected resource: %+v", project)
	}
}

func TestProjectHandler_Create_Duplicate(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	if rec := createProject(h, "nic", "ministry-health", "payments", "prod", ""); rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d", rec.Code)
	}
	rec := createProject(h, "nic", "ministry-health", "payments", "prod", "")
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceAlreadyExists {
		t.Errorf("code = %q, want RESOURCE_ALREADY_EXISTS", errBody.Code)
	}
}

func TestProjectHandler_Create_InvalidName(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	rec := createProject(h, "nic", "ministry-health", "payments", "INVALID", "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "metadata.name" {
		t.Errorf("error = %+v, want VALIDATION_FAILED metadata.name", errBody)
	}
}

func TestProjectHandler_Create_MissingOrganizationName(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	body := map[string]any{
		"metadata": map[string]any{"name": "prod"},
		"spec": map[string]any{
			"organizationUnitName": "ministry-health",
			"tenantName":           "payments",
		},
	}
	req := jsonRequest(http.MethodPost, "/v1/projects", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.organizationName" {
		t.Errorf("error = %+v, want VALIDATION_FAILED spec.organizationName", errBody)
	}
}

func TestProjectHandler_Create_MissingOrganizationUnitName(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	body := map[string]any{
		"metadata": map[string]any{"name": "prod"},
		"spec": map[string]any{
			"organizationName": "nic",
			"tenantName":       "payments",
		},
	}
	req := jsonRequest(http.MethodPost, "/v1/projects", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.organizationUnitName" {
		t.Errorf("error = %+v, want VALIDATION_FAILED spec.organizationUnitName", errBody)
	}
}

func TestProjectHandler_Create_MissingTenantName(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	body := map[string]any{
		"metadata": map[string]any{"name": "prod"},
		"spec": map[string]any{
			"organizationName":     "nic",
			"organizationUnitName": "ministry-health",
		},
	}
	req := jsonRequest(http.MethodPost, "/v1/projects", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.tenantName" {
		t.Errorf("error = %+v, want VALIDATION_FAILED spec.tenantName", errBody)
	}
}

func TestProjectHandler_Create_NonExistentParent(t *testing.T) {
	h, _, _ := newTestProjectHandler()

	rec := createProject(h, "nic", "ministry-health", "ghost-tenant", "prod", "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed {
		t.Errorf("code = %q, want VALIDATION_FAILED", errBody.Code)
	}
	if errBody.Field != "spec.tenantName" {
		t.Errorf("field = %q, want spec.tenantName", errBody.Field)
	}
	if !strings.Contains(errBody.Message, "nic/ministry-health/ghost-tenant") {
		t.Errorf("message = %q, want it to include full parent reference nic/ministry-health/ghost-tenant", errBody.Message)
	}
}

func TestProjectHandler_Create_StatusFieldRejected(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	payload := `{"metadata":{"name":"prod"},"spec":{"organizationName":"nic","organizationUnitName":"ministry-health","tenantName":"payments"},"status":{}}`
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "status" {
		t.Errorf("error = %+v, want VALIDATION_FAILED status", errBody)
	}
}

func TestProjectHandler_Create_BadJSON(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader("{")), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed {
		t.Errorf("code = %q, want VALIDATION_FAILED", errBody.Code)
	}
}

func TestProjectHandler_Create_UnknownField(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	payload := `{"metadata":{"name":"prod"},"spec":{"organizationName":"nic","organizationUnitName":"ministry-health","tenantName":"payments"},"bogus":true}`
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed {
		t.Errorf("code = %q, want VALIDATION_FAILED", errBody.Code)
	}
}

func TestProjectHandler_Create_OversizedBody(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	large := strings.Repeat("a", 1<<20+1)
	payload := fmt.Sprintf(`{"metadata":{"name":"prod"},"spec":{"organizationName":"nic","organizationUnitName":"ministry-health","tenantName":"payments","description":"%s"}}`, large)
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

func TestProjectHandler_Get_Exists(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")
	if rec := createProject(h, "nic", "ministry-health", "payments", "prod", "Production"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	req := jsonRequest(http.MethodGet, "/v1/projects/nic/ministry-health/payments/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var project resources.Project
	if err := json.NewDecoder(rec.Body).Decode(&project); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if project.Metadata.Name != "prod" || project.Spec.OrganizationName != "nic" ||
		project.Spec.OrganizationUnitName != "ministry-health" || project.Spec.TenantName != "payments" ||
		project.APIVersion == "" || project.Kind == "" || project.Status.Phase == "" {
		t.Errorf("incomplete resource: %+v", project)
	}
}

func TestProjectHandler_Get_NotFound(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	req := jsonRequest(http.MethodGet, "/v1/projects/nic/ministry-health/payments/missing", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}

func TestProjectHandler_Get_InvalidOrgNameSegment(t *testing.T) {
	h, _, _ := newTestProjectHandler()

	req := jsonRequest(http.MethodGet, "/v1/projects/INVALID/ministry-health/payments/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "spec.organizationName" {
		t.Errorf("field = %q, want spec.organizationName", errBody.Field)
	}
}

func TestProjectHandler_Get_InvalidOUNameSegment(t *testing.T) {
	h, _, _ := newTestProjectHandler()

	req := jsonRequest(http.MethodGet, "/v1/projects/nic/INVALID/payments/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "spec.organizationUnitName" {
		t.Errorf("field = %q, want spec.organizationUnitName", errBody.Field)
	}
}

func TestProjectHandler_Get_InvalidTenantNameSegment(t *testing.T) {
	h, _, _ := newTestProjectHandler()

	req := jsonRequest(http.MethodGet, "/v1/projects/nic/ministry-health/INVALID/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "spec.tenantName" {
		t.Errorf("field = %q, want spec.tenantName", errBody.Field)
	}
}

func TestProjectHandler_Get_InvalidNameSegment(t *testing.T) {
	h, _, _ := newTestProjectHandler()

	req := jsonRequest(http.MethodGet, "/v1/projects/nic/ministry-health/payments/INVALID", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "metadata.name" {
		t.Errorf("field = %q, want metadata.name", errBody.Field)
	}
}

func TestProjectHandler_Get_BareItemPath_NotFound(t *testing.T) {
	h, _, _ := newTestProjectHandler()

	req := jsonRequest(http.MethodGet, "/v1/projects/", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestProjectHandler_Get_ThreeSegments_NotFound(t *testing.T) {
	h, _, _ := newTestProjectHandler()

	req := jsonRequest(http.MethodGet, "/v1/projects/nic/ministry-health/payments", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestProjectHandler_Get_ExtraSegment_NotFound(t *testing.T) {
	h, _, _ := newTestProjectHandler()

	req := jsonRequest(http.MethodGet, "/v1/projects/nic/ministry-health/payments/prod/extra", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestProjectHandler_List_Sorted(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "alpha", "unit-a", "tenant-a")
	seedTenant(t, tenantReg, "alpha", "unit-b", "tenant-a")
	seedTenant(t, tenantReg, "zebra", "unit-b", "tenant-c")

	inputs := []struct{ org, ou, tenant, name string }{
		{"zebra", "unit-b", "tenant-c", "beta"},
		{"alpha", "unit-b", "tenant-a", "delta"},
		{"alpha", "unit-a", "tenant-a", "charlie"},
		{"alpha", "unit-a", "tenant-a", "bravo"},
	}
	for _, in := range inputs {
		if rec := createProject(h, in.org, in.ou, in.tenant, in.name, ""); rec.Code != http.StatusCreated {
			t.Fatalf("create %s/%s/%s/%s status = %d", in.org, in.ou, in.tenant, in.name, rec.Code)
		}
	}

	req := jsonRequest(http.MethodGet, "/v1/projects", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp projectListResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Items) != 4 {
		t.Fatalf("items = %d, want 4", len(resp.Items))
	}
	want := []struct{ org, ou, tenant, name string }{
		{"alpha", "unit-a", "tenant-a", "bravo"},
		{"alpha", "unit-a", "tenant-a", "charlie"},
		{"alpha", "unit-b", "tenant-a", "delta"},
		{"zebra", "unit-b", "tenant-c", "beta"},
	}
	for i, wnt := range want {
		got := resp.Items[i]
		if got.Spec.OrganizationName != wnt.org || got.Spec.OrganizationUnitName != wnt.ou ||
			got.Spec.TenantName != wnt.tenant || got.Metadata.Name != wnt.name {
			t.Errorf("item[%d] = %s/%s/%s/%s, want %s/%s/%s/%s", i,
				got.Spec.OrganizationName, got.Spec.OrganizationUnitName, got.Spec.TenantName, got.Metadata.Name,
				wnt.org, wnt.ou, wnt.tenant, wnt.name)
		}
	}
}

func TestProjectHandler_List_Empty(t *testing.T) {
	h, _, _ := newTestProjectHandler()

	req := jsonRequest(http.MethodGet, "/v1/projects", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"items":[]}` {
		t.Errorf("body = %q, want {\"items\":[]}", got)
	}
}

func TestProjectHandler_Update_Valid(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")
	if rec := createProject(h, "nic", "ministry-health", "payments", "prod", "old"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	body := map[string]any{
		"metadata": map[string]any{"name": "prod", "displayName": "New Display"},
		"spec": map[string]any{
			"organizationName":     "nic",
			"organizationUnitName": "ministry-health",
			"tenantName":           "payments",
			"description":          "updated",
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/projects/nic/ministry-health/payments/prod", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var project resources.Project
	if err := json.NewDecoder(rec.Body).Decode(&project); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if project.Metadata.DisplayName != "New Display" || project.Spec.Description != "updated" {
		t.Errorf("mutable fields not updated: %+v", project)
	}
	if project.Metadata.Name != "prod" || project.Spec.OrganizationName != "nic" ||
		project.Spec.OrganizationUnitName != "ministry-health" || project.Spec.TenantName != "payments" {
		t.Errorf("immutable fields changed: %+v", project)
	}
}

func TestProjectHandler_Update_PreservesServerOwnedFields(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")
	if rec := createProject(h, "nic", "ministry-health", "payments", "prod", "old"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	payload := `{"apiVersion":"tampered/v0","kind":"Tampered","metadata":{"name":"prod"},"spec":{"organizationName":"nic","organizationUnitName":"ministry-health","tenantName":"payments","description":"changed"}}`
	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/projects/nic/ministry-health/payments/prod", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var project resources.Project
	if err := json.NewDecoder(rec.Body).Decode(&project); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if project.APIVersion != resources.ProjectAPIVersion || project.Kind != resources.ProjectKind {
		t.Errorf("server-owned fields not preserved by registry: %+v", project)
	}
	if project.Status.Phase != resources.PhaseActive {
		t.Errorf("status.phase = %q, want Active", project.Status.Phase)
	}
	if project.Spec.Description != "changed" {
		t.Errorf("description = %q, want changed", project.Spec.Description)
	}
}

func TestProjectHandler_Update_NotFound(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	body := map[string]any{
		"metadata": map[string]any{"name": "missing"},
		"spec": map[string]any{
			"organizationName":     "nic",
			"organizationUnitName": "ministry-health",
			"tenantName":           "payments",
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/projects/nic/ministry-health/payments/missing", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}

func TestProjectHandler_Update_NameMismatch(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")
	_ = createProject(h, "nic", "ministry-health", "payments", "prod", "")

	body := map[string]any{
		"metadata": map[string]any{"name": "other-name"},
		"spec": map[string]any{
			"organizationName":     "nic",
			"organizationUnitName": "ministry-health",
			"tenantName":           "payments",
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/projects/nic/ministry-health/payments/prod", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "metadata.name" {
		t.Errorf("field = %q, want metadata.name", errBody.Field)
	}
}

func TestProjectHandler_Update_OrganizationNameMismatch(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")
	_ = createProject(h, "nic", "ministry-health", "payments", "prod", "")

	body := map[string]any{
		"metadata": map[string]any{"name": "prod"},
		"spec": map[string]any{
			"organizationName":     "other-org",
			"organizationUnitName": "ministry-health",
			"tenantName":           "payments",
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/projects/nic/ministry-health/payments/prod", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "spec.organizationName" {
		t.Errorf("field = %q, want spec.organizationName", errBody.Field)
	}
}

func TestProjectHandler_Update_OrganizationUnitNameMismatch(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")
	_ = createProject(h, "nic", "ministry-health", "payments", "prod", "")

	body := map[string]any{
		"metadata": map[string]any{"name": "prod"},
		"spec": map[string]any{
			"organizationName":     "nic",
			"organizationUnitName": "other-unit",
			"tenantName":           "payments",
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/projects/nic/ministry-health/payments/prod", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "spec.organizationUnitName" {
		t.Errorf("field = %q, want spec.organizationUnitName", errBody.Field)
	}
}

func TestProjectHandler_Update_TenantNameMismatch(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")
	_ = createProject(h, "nic", "ministry-health", "payments", "prod", "")

	body := map[string]any{
		"metadata": map[string]any{"name": "prod"},
		"spec": map[string]any{
			"organizationName":     "nic",
			"organizationUnitName": "ministry-health",
			"tenantName":           "other-tenant",
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/projects/nic/ministry-health/payments/prod", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "spec.tenantName" {
		t.Errorf("field = %q, want spec.tenantName", errBody.Field)
	}
}

func TestProjectHandler_Update_StatusField(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")
	_ = createProject(h, "nic", "ministry-health", "payments", "prod", "")

	payload := `{"metadata":{"name":"prod"},"spec":{"organizationName":"nic","organizationUnitName":"ministry-health","tenantName":"payments"},"status":{}}`
	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/projects/nic/ministry-health/payments/prod", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "status" {
		t.Errorf("field = %q, want status", errBody.Field)
	}
}

func TestProjectHandler_Update_BadJSON(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")
	_ = createProject(h, "nic", "ministry-health", "payments", "prod", "")

	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/projects/nic/ministry-health/payments/prod", strings.NewReader("{")), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestProjectHandler_Delete_Success(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")
	_ = createProject(h, "nic", "ministry-health", "payments", "prod", "")

	req := jsonRequest(http.MethodDelete, "/v1/projects/nic/ministry-health/payments/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

func TestProjectHandler_Delete_NotFound(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")

	req := jsonRequest(http.MethodDelete, "/v1/projects/nic/ministry-health/payments/missing", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}

func TestProjectHandler_Delete_InvalidPathSegment(t *testing.T) {
	h, _, _ := newTestProjectHandler()

	req := jsonRequest(http.MethodDelete, "/v1/projects/nic/ministry-health/payments/BAD", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed {
		t.Errorf("code = %q, want VALIDATION_FAILED", errBody.Code)
	}
}

func TestProjectHandler_Collection_MethodNotAllowed(t *testing.T) {
	h, _, _ := newTestProjectHandler()

	req := jsonRequest(http.MethodDelete, "/v1/projects", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}

func TestProjectHandler_Item_MethodNotAllowed(t *testing.T) {
	h, tenantReg, _ := newTestProjectHandler()
	seedTenant(t, tenantReg, "nic", "ministry-health", "payments")
	_ = createProject(h, "nic", "ministry-health", "payments", "prod", "")

	req := jsonRequest(http.MethodPost, "/v1/projects/nic/ministry-health/payments/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}

// newProjectBlockerWiring builds a TenantHandler and ProjectHandler that share
// in-memory registries, with the real ProjectChildBlockerChecker wired into the
// Tenant delete path (matching production wiring in cmd/sovrunn-api/main.go).
func newProjectBlockerWiring() (*TenantHandler, *ProjectHandler, *registry.OrganizationUnitRegistry, *registry.TenantRegistry) {
	ouRegistry := registry.NewOrganizationUnitRegistry()
	tenantRegistry := registry.NewTenantRegistry()
	projectRegistry := registry.NewProjectRegistry()
	projectBlocker := registry.NewProjectChildBlockerChecker(projectRegistry)
	tenantHandler := NewTenantHandler(tenantRegistry, ouRegistry, projectBlocker, nil)
	projectHandler := NewProjectHandler(projectRegistry, tenantRegistry, nil)
	return tenantHandler, projectHandler, ouRegistry, tenantRegistry
}

func TestIntegration_TenantDeleteBlockedByProject(t *testing.T) {
	tenantHandler, projectHandler, ouReg, _ := newProjectBlockerWiring()
	seedOU(t, ouReg, "nic", "ministry-health")

	if rec := createTenant(tenantHandler, "nic", "ministry-health", "payments", ""); rec.Code != http.StatusCreated {
		t.Fatalf("create tenant status = %d; body=%s", rec.Code, rec.Body.String())
	}
	if rec := createProject(projectHandler, "nic", "ministry-health", "payments", "prod", ""); rec.Code != http.StatusCreated {
		t.Fatalf("create project status = %d; body=%s", rec.Code, rec.Body.String())
	}

	req := jsonRequest(http.MethodDelete, "/v1/tenants/nic/ministry-health/payments", nil, "")
	rec := httptest.NewRecorder()
	tenantHandler.HandleItem(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeDeleteBlocked {
		t.Errorf("code = %q, want DELETE_BLOCKED", errBody.Code)
	}
	if !strings.Contains(errBody.Message, "Project") {
		t.Errorf("message = %q, want it to contain Project", errBody.Message)
	}
}

func TestIntegration_TenantDeleteAllowedWhenNoProject(t *testing.T) {
	tenantHandler, _, ouReg, _ := newProjectBlockerWiring()
	seedOU(t, ouReg, "nic", "ministry-health")

	if rec := createTenant(tenantHandler, "nic", "ministry-health", "payments", ""); rec.Code != http.StatusCreated {
		t.Fatalf("create tenant status = %d; body=%s", rec.Code, rec.Body.String())
	}

	req := jsonRequest(http.MethodDelete, "/v1/tenants/nic/ministry-health/payments", nil, "")
	rec := httptest.NewRecorder()
	tenantHandler.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}
