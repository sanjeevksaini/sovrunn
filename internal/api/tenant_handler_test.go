package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// newTestTenantHandler builds a TenantHandler backed by fresh registries and
// returns the handler along with the OU and Tenant registries so tests can
// seed parent OrganizationUnits directly. The OrganizationUnitRegistry
// satisfies registry.OrganizationUnitLookup.
func newTestTenantHandler() (*TenantHandler, *registry.OrganizationUnitRegistry, *registry.TenantRegistry) {
	ouRegistry := registry.NewOrganizationUnitRegistry()
	tenantRegistry := registry.NewTenantRegistry()
	handler := NewTenantHandler(tenantRegistry, ouRegistry, nil, nil)
	return handler, ouRegistry, tenantRegistry
}

// seedOU creates a parent OrganizationUnit directly in the registry.
func seedOU(t *testing.T, reg *registry.OrganizationUnitRegistry, orgName, name string) {
	t.Helper()
	ou := resources.OrganizationUnit{
		APIVersion: resources.OUAPIVersion,
		Kind:       resources.OUKind,
		Metadata:   resources.Metadata{Name: name},
		Spec:       resources.OrganizationUnitSpec{OrganizationName: orgName},
		Status:     resources.OrganizationUnitStatus{Phase: resources.PhaseActive},
	}
	if _, err := reg.CreateOrganizationUnit(context.Background(), ou); err != nil {
		t.Fatalf("seedOU(%s/%s): %v", orgName, name, err)
	}
}

// createTenant issues a POST create through the handler and returns the recorder.
func createTenant(h *TenantHandler, orgName, ouName, name, desc string) *httptest.ResponseRecorder {
	body := map[string]any{
		"metadata": map[string]any{"name": name},
		"spec": map[string]any{
			"organizationName":     orgName,
			"organizationUnitName": ouName,
			"description":          desc,
		},
	}
	req := jsonRequest(http.MethodPost, "/v1/tenants", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	return rec
}

func TestTenantHandler_Create_Valid(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")

	rec := createTenant(h, "nic", "ministry-health", "prod", "Production")
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var tenant resources.Tenant
	if err := json.NewDecoder(rec.Body).Decode(&tenant); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if tenant.APIVersion != resources.TenantAPIVersion || tenant.Kind != resources.TenantKind {
		t.Errorf("apiVersion/kind not set by server: %+v", tenant)
	}
	if tenant.Status.Phase != resources.PhaseActive {
		t.Errorf("phase = %q, want Active", tenant.Status.Phase)
	}
	if tenant.Metadata.Name != "prod" || tenant.Spec.OrganizationName != "nic" ||
		tenant.Spec.OrganizationUnitName != "ministry-health" {
		t.Errorf("unexpected resource: %+v", tenant)
	}
}

func TestTenantHandler_Create_Duplicate(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")

	if rec := createTenant(h, "nic", "ministry-health", "prod", ""); rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d", rec.Code)
	}
	rec := createTenant(h, "nic", "ministry-health", "prod", "")
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceAlreadyExists {
		t.Errorf("code = %q, want RESOURCE_ALREADY_EXISTS", errBody.Code)
	}
}

func TestTenantHandler_Create_InvalidName(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")

	rec := createTenant(h, "nic", "ministry-health", "INVALID", "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "metadata.name" {
		t.Errorf("error = %+v, want VALIDATION_FAILED metadata.name", errBody)
	}
}

func TestTenantHandler_Create_MissingOrganizationName(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")

	body := map[string]any{
		"metadata": map[string]any{"name": "prod"},
		"spec":     map[string]any{"organizationUnitName": "ministry-health"},
	}
	req := jsonRequest(http.MethodPost, "/v1/tenants", body, "application/json")
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

func TestTenantHandler_Create_MissingOrganizationUnitName(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")

	body := map[string]any{
		"metadata": map[string]any{"name": "prod"},
		"spec":     map[string]any{"organizationName": "nic"},
	}
	req := jsonRequest(http.MethodPost, "/v1/tenants", body, "application/json")
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

func TestTenantHandler_Create_NonExistentParent(t *testing.T) {
	h, _, _ := newTestTenantHandler()

	rec := createTenant(h, "nic", "ghost-unit", "prod", "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed {
		t.Errorf("code = %q, want VALIDATION_FAILED", errBody.Code)
	}
	if errBody.Field != "spec.organizationUnitName" {
		t.Errorf("field = %q, want spec.organizationUnitName", errBody.Field)
	}
	if !strings.Contains(errBody.Message, "nic/ghost-unit") {
		t.Errorf("message = %q, want it to include full parent reference nic/ghost-unit", errBody.Message)
	}
}

func TestTenantHandler_Create_StatusFieldRejected(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")

	payload := `{"metadata":{"name":"prod"},"spec":{"organizationName":"nic","organizationUnitName":"ministry-health"},"status":{}}`
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/tenants", strings.NewReader(payload)), "id")
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

func TestTenantHandler_Create_BadJSON(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")

	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/tenants", strings.NewReader("{")), "id")
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

func TestTenantHandler_Create_UnknownField(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")

	payload := `{"metadata":{"name":"prod"},"spec":{"organizationName":"nic","organizationUnitName":"ministry-health"},"bogus":true}`
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/tenants", strings.NewReader(payload)), "id")
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

func TestTenantHandler_Create_OversizedBody(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")

	large := strings.Repeat("a", 1<<20+1)
	payload := fmt.Sprintf(`{"metadata":{"name":"prod"},"spec":{"organizationName":"nic","organizationUnitName":"ministry-health","description":"%s"}}`, large)
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/tenants", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

func TestTenantHandler_Get_Exists(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")
	if rec := createTenant(h, "nic", "ministry-health", "prod", "Production"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	req := jsonRequest(http.MethodGet, "/v1/tenants/nic/ministry-health/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var tenant resources.Tenant
	if err := json.NewDecoder(rec.Body).Decode(&tenant); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if tenant.Metadata.Name != "prod" || tenant.Spec.OrganizationName != "nic" ||
		tenant.Spec.OrganizationUnitName != "ministry-health" ||
		tenant.APIVersion == "" || tenant.Kind == "" || tenant.Status.Phase == "" {
		t.Errorf("incomplete resource: %+v", tenant)
	}
}

func TestTenantHandler_Get_NotFound(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")

	req := jsonRequest(http.MethodGet, "/v1/tenants/nic/ministry-health/missing", nil, "")
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

func TestTenantHandler_Get_InvalidNameSegment(t *testing.T) {
	h, _, _ := newTestTenantHandler()

	req := jsonRequest(http.MethodGet, "/v1/tenants/nic/ministry-health/INVALID", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "metadata.name" {
		t.Errorf("error = %+v, want VALIDATION_FAILED metadata.name", errBody)
	}
}

func TestTenantHandler_Get_InvalidOrgNameSegment(t *testing.T) {
	h, _, _ := newTestTenantHandler()

	req := jsonRequest(http.MethodGet, "/v1/tenants/INVALID/ministry-health/prod", nil, "")
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

func TestTenantHandler_Get_InvalidOUNameSegment(t *testing.T) {
	h, _, _ := newTestTenantHandler()

	req := jsonRequest(http.MethodGet, "/v1/tenants/nic/INVALID/prod", nil, "")
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

func TestTenantHandler_Get_BareItemPath_NotFound(t *testing.T) {
	h, _, _ := newTestTenantHandler()

	req := jsonRequest(http.MethodGet, "/v1/tenants/", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestTenantHandler_Get_TwoSegments_NotFound(t *testing.T) {
	h, _, _ := newTestTenantHandler()

	req := jsonRequest(http.MethodGet, "/v1/tenants/nic/ministry-health", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestTenantHandler_Get_ExtraSegment_NotFound(t *testing.T) {
	h, _, _ := newTestTenantHandler()

	req := jsonRequest(http.MethodGet, "/v1/tenants/nic/ministry-health/prod/extra", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestTenantHandler_List_Sorted(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "alpha", "unit-a")
	seedOU(t, ouReg, "alpha", "unit-b")
	seedOU(t, ouReg, "zebra", "unit-b")

	inputs := []struct{ org, ou, name string }{
		{"zebra", "unit-b", "beta"},
		{"alpha", "unit-b", "delta"},
		{"alpha", "unit-a", "charlie"},
	}
	for _, in := range inputs {
		if rec := createTenant(h, in.org, in.ou, in.name, ""); rec.Code != http.StatusCreated {
			t.Fatalf("create %s/%s/%s status = %d", in.org, in.ou, in.name, rec.Code)
		}
	}

	req := jsonRequest(http.MethodGet, "/v1/tenants", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp tenantListResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Items) != 3 {
		t.Fatalf("items = %d, want 3", len(resp.Items))
	}
	want := []struct{ org, ou, name string }{
		{"alpha", "unit-a", "charlie"},
		{"alpha", "unit-b", "delta"},
		{"zebra", "unit-b", "beta"},
	}
	for i, wnt := range want {
		got := resp.Items[i]
		if got.Spec.OrganizationName != wnt.org || got.Spec.OrganizationUnitName != wnt.ou || got.Metadata.Name != wnt.name {
			t.Errorf("item[%d] = %s/%s/%s, want %s/%s/%s", i,
				got.Spec.OrganizationName, got.Spec.OrganizationUnitName, got.Metadata.Name,
				wnt.org, wnt.ou, wnt.name)
		}
	}
}

func TestTenantHandler_List_Empty(t *testing.T) {
	h, _, _ := newTestTenantHandler()

	req := jsonRequest(http.MethodGet, "/v1/tenants", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"items":[]}` {
		t.Errorf("body = %q, want {\"items\":[]}", got)
	}
}

func TestTenantHandler_Update_Valid(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")
	if rec := createTenant(h, "nic", "ministry-health", "prod", "old"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	body := map[string]any{
		"metadata": map[string]any{"name": "prod", "displayName": "New Display"},
		"spec": map[string]any{
			"organizationName":     "nic",
			"organizationUnitName": "ministry-health",
			"description":          "updated",
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/tenants/nic/ministry-health/prod", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var tenant resources.Tenant
	if err := json.NewDecoder(rec.Body).Decode(&tenant); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if tenant.Metadata.DisplayName != "New Display" || tenant.Spec.Description != "updated" {
		t.Errorf("mutable fields not updated: %+v", tenant)
	}
	if tenant.Metadata.Name != "prod" || tenant.Spec.OrganizationName != "nic" ||
		tenant.Spec.OrganizationUnitName != "ministry-health" {
		t.Errorf("immutable fields changed: %+v", tenant)
	}
}

func TestTenantHandler_Update_PreservesServerOwnedFields(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")
	if rec := createTenant(h, "nic", "ministry-health", "prod", "old"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	// Attempt to tamper with server-owned apiVersion/kind through the body.
	payload := `{"apiVersion":"tampered/v0","kind":"Tampered","metadata":{"name":"prod"},"spec":{"organizationName":"nic","organizationUnitName":"ministry-health","description":"changed"}}`
	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/tenants/nic/ministry-health/prod", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var tenant resources.Tenant
	if err := json.NewDecoder(rec.Body).Decode(&tenant); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if tenant.APIVersion != resources.TenantAPIVersion || tenant.Kind != resources.TenantKind {
		t.Errorf("server-owned fields not preserved by registry: %+v", tenant)
	}
	if tenant.Status.Phase != resources.PhaseActive {
		t.Errorf("status.phase = %q, want Active", tenant.Status.Phase)
	}
	if tenant.Spec.Description != "changed" {
		t.Errorf("description = %q, want changed", tenant.Spec.Description)
	}
}

func TestTenantHandler_Update_NotFound(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")

	body := map[string]any{
		"metadata": map[string]any{"name": "missing"},
		"spec":     map[string]any{"organizationName": "nic", "organizationUnitName": "ministry-health"},
	}
	req := jsonRequest(http.MethodPut, "/v1/tenants/nic/ministry-health/missing", body, "application/json")
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

func TestTenantHandler_Update_NameMismatch(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")
	_ = createTenant(h, "nic", "ministry-health", "prod", "")

	body := map[string]any{
		"metadata": map[string]any{"name": "other-name"},
		"spec":     map[string]any{"organizationName": "nic", "organizationUnitName": "ministry-health"},
	}
	req := jsonRequest(http.MethodPut, "/v1/tenants/nic/ministry-health/prod", body, "application/json")
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

func TestTenantHandler_Update_OrganizationNameMismatch(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")
	_ = createTenant(h, "nic", "ministry-health", "prod", "")

	body := map[string]any{
		"metadata": map[string]any{"name": "prod"},
		"spec":     map[string]any{"organizationName": "other-org", "organizationUnitName": "ministry-health"},
	}
	req := jsonRequest(http.MethodPut, "/v1/tenants/nic/ministry-health/prod", body, "application/json")
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

func TestTenantHandler_Update_OrganizationUnitNameMismatch(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")
	_ = createTenant(h, "nic", "ministry-health", "prod", "")

	body := map[string]any{
		"metadata": map[string]any{"name": "prod"},
		"spec":     map[string]any{"organizationName": "nic", "organizationUnitName": "other-unit"},
	}
	req := jsonRequest(http.MethodPut, "/v1/tenants/nic/ministry-health/prod", body, "application/json")
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

func TestTenantHandler_Update_StatusField(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")
	_ = createTenant(h, "nic", "ministry-health", "prod", "")

	payload := `{"metadata":{"name":"prod"},"spec":{"organizationName":"nic","organizationUnitName":"ministry-health"},"status":{}}`
	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/tenants/nic/ministry-health/prod", strings.NewReader(payload)), "id")
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

func TestTenantHandler_Update_BadJSON(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")
	_ = createTenant(h, "nic", "ministry-health", "prod", "")

	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/tenants/nic/ministry-health/prod", strings.NewReader("{")), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestTenantHandler_Delete_Success(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")
	_ = createTenant(h, "nic", "ministry-health", "prod", "")

	req := jsonRequest(http.MethodDelete, "/v1/tenants/nic/ministry-health/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

func TestTenantHandler_Delete_NotFound(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")

	req := jsonRequest(http.MethodDelete, "/v1/tenants/nic/ministry-health/missing", nil, "")
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

func TestTenantHandler_Delete_InvalidPathSegment(t *testing.T) {
	h, _, _ := newTestTenantHandler()

	req := jsonRequest(http.MethodDelete, "/v1/tenants/nic/ministry-health/BAD", nil, "")
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

func TestTenantHandler_Collection_MethodNotAllowed(t *testing.T) {
	h, _, _ := newTestTenantHandler()

	req := jsonRequest(http.MethodDelete, "/v1/tenants", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}

func TestTenantHandler_Item_MethodNotAllowed(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")
	_ = createTenant(h, "nic", "ministry-health", "prod", "")

	req := jsonRequest(http.MethodPost, "/v1/tenants/nic/ministry-health/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}

// newTenantBlockerWiring builds an OUHandler and TenantHandler that share a
// single TenantRegistry through a TenantChildBlockerChecker, mirroring the
// production wiring in cmd/sovrunn-api/main.go. This lets HTTP-level tests
// exercise OrganizationUnit delete blocking by Tenants end-to-end.
func newTenantBlockerWiring() (*OUHandler, *TenantHandler, *registry.OrganizationRegistry, *registry.OrganizationUnitRegistry) {
	orgRegistry := registry.NewOrganizationRegistry()
	ouRegistry := registry.NewOrganizationUnitRegistry()
	tenantRegistry := registry.NewTenantRegistry()
	tenantBlocker := registry.NewTenantChildBlockerChecker(tenantRegistry)
	ouHandler := NewOUHandler(ouRegistry, orgRegistry, tenantBlocker, nil)
	tenantHandler := NewTenantHandler(tenantRegistry, ouRegistry, nil, nil)
	return ouHandler, tenantHandler, orgRegistry, ouRegistry
}

func TestOUDeleteBlockedByTenant(t *testing.T) {
	ouHandler, tenantHandler, orgReg, _ := newTenantBlockerWiring()
	seedOrg(t, orgReg, "nic")
	if rec := createOU(ouHandler, "nic", "ministry-health", ""); rec.Code != http.StatusCreated {
		t.Fatalf("create OU status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	if rec := createTenant(tenantHandler, "nic", "ministry-health", "prod", ""); rec.Code != http.StatusCreated {
		t.Fatalf("create Tenant status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}

	req := jsonRequest(http.MethodDelete, "/v1/organization-units/nic/ministry-health", nil, "")
	rec := httptest.NewRecorder()
	ouHandler.HandleItem(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeDeleteBlocked {
		t.Errorf("code = %q, want DELETE_BLOCKED", errBody.Code)
	}
	if !strings.Contains(errBody.Message, "Tenant") {
		t.Errorf("message = %q, want it to identify Tenant", errBody.Message)
	}
}

func TestOUDeleteAllowedWhenNoTenants(t *testing.T) {
	ouHandler, _, orgReg, _ := newTenantBlockerWiring()
	seedOrg(t, orgReg, "nic")
	if rec := createOU(ouHandler, "nic", "ministry-health", ""); rec.Code != http.StatusCreated {
		t.Fatalf("create OU status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}

	req := jsonRequest(http.MethodDelete, "/v1/organization-units/nic/ministry-health", nil, "")
	rec := httptest.NewRecorder()
	ouHandler.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

// stubTenantChildBlocker is a controllable registry.TenantChildBlocker used to
// exercise the FEATURE-0004 Tenant delete blocking path.
type stubTenantChildBlocker struct {
	blockers []registry.BlockedBy
	err      error
}

func (s stubTenantChildBlocker) BlockedByTenantChildren(
	_ context.Context, _, _, _ string,
) ([]registry.BlockedBy, error) {
	return s.blockers, s.err
}

// newTenantHandlerWithBlocker builds a TenantHandler with the given blocker and
// returns it alongside the OU and Tenant registries for seeding.
func newTenantHandlerWithBlocker(blocker registry.TenantChildBlocker) (*TenantHandler, *registry.OrganizationUnitRegistry, *registry.TenantRegistry) {
	ouRegistry := registry.NewOrganizationUnitRegistry()
	tenantRegistry := registry.NewTenantRegistry()
	handler := NewTenantHandler(tenantRegistry, ouRegistry, blocker, nil)
	return handler, ouRegistry, tenantRegistry
}

func TestTenantHandler_Delete_NilBlockerAllows(t *testing.T) {
	h, ouReg, _ := newTestTenantHandler()
	seedOU(t, ouReg, "nic", "ministry-health")
	if rec := createTenant(h, "nic", "ministry-health", "prod", ""); rec.Code != http.StatusCreated {
		t.Fatalf("create Tenant status = %d, want 201", rec.Code)
	}

	req := jsonRequest(http.MethodDelete, "/v1/tenants/nic/ministry-health/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

func TestTenantHandler_Delete_BlockedByProject(t *testing.T) {
	blocker := stubTenantChildBlocker{blockers: []registry.BlockedBy{{Kind: "Project", Count: 1}}}
	h, ouReg, _ := newTenantHandlerWithBlocker(blocker)
	seedOU(t, ouReg, "nic", "ministry-health")
	if rec := createTenant(h, "nic", "ministry-health", "prod", ""); rec.Code != http.StatusCreated {
		t.Fatalf("create Tenant status = %d, want 201", rec.Code)
	}

	req := jsonRequest(http.MethodDelete, "/v1/tenants/nic/ministry-health/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeDeleteBlocked {
		t.Errorf("code = %q, want DELETE_BLOCKED", errBody.Code)
	}
	if !strings.Contains(errBody.Message, "Project") {
		t.Errorf("message = %q, want it to identify Project", errBody.Message)
	}
}

func TestTenantHandler_Delete_EmptyBlockerAllows(t *testing.T) {
	blocker := stubTenantChildBlocker{blockers: nil}
	h, ouReg, _ := newTenantHandlerWithBlocker(blocker)
	seedOU(t, ouReg, "nic", "ministry-health")
	if rec := createTenant(h, "nic", "ministry-health", "prod", ""); rec.Code != http.StatusCreated {
		t.Fatalf("create Tenant status = %d, want 201", rec.Code)
	}

	req := jsonRequest(http.MethodDelete, "/v1/tenants/nic/ministry-health/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

func TestTenantHandler_Delete_BlockerErrorMapsTo500(t *testing.T) {
	blocker := stubTenantChildBlocker{err: errors.New("count failed")}
	h, ouReg, _ := newTenantHandlerWithBlocker(blocker)
	seedOU(t, ouReg, "nic", "ministry-health")
	if rec := createTenant(h, "nic", "ministry-health", "prod", ""); rec.Code != http.StatusCreated {
		t.Fatalf("create Tenant status = %d, want 201", rec.Code)
	}

	req := jsonRequest(http.MethodDelete, "/v1/tenants/nic/ministry-health/prod", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeInternalError {
		t.Errorf("code = %q, want INTERNAL_ERROR", errBody.Code)
	}
}
