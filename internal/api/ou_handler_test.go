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

// newTestOUHandler builds an OUHandler backed by fresh registries and
// returns the handler along with both registries so tests can seed parent
// Organizations directly. orgRegistry satisfies registry.OrganizationLookup.
func newTestOUHandler() (*OUHandler, *registry.OrganizationRegistry, *registry.OrganizationUnitRegistry) {
	orgRegistry := registry.NewOrganizationRegistry()
	ouRegistry := registry.NewOrganizationUnitRegistry()
	handler := NewOUHandler(ouRegistry, orgRegistry, nil)
	return handler, orgRegistry, ouRegistry
}

// seedOrg creates a parent Organization directly in the registry.
func seedOrg(t *testing.T, reg *registry.OrganizationRegistry, name string) {
	t.Helper()
	org := resources.Organization{
		APIVersion: resources.OrgAPIVersion,
		Kind:       resources.OrgKind,
		Metadata:   resources.Metadata{Name: name},
		Status:     resources.OrganizationStatus{Phase: resources.PhaseActive},
	}
	if err := reg.CreateOrganization(context.Background(), org); err != nil {
		t.Fatalf("seedOrg(%s): %v", name, err)
	}
}

// createOU issues a POST create through the handler and returns the recorder.
func createOU(h *OUHandler, orgName, name, desc string) *httptest.ResponseRecorder {
	body := map[string]any{
		"metadata": map[string]any{"name": name},
		"spec":     map[string]any{"organizationName": orgName, "description": desc},
	}
	req := jsonRequest(http.MethodPost, "/v1/organization-units", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	return rec
}

// newBlockerWiring builds an OrgHandler and OUHandler that share a single
// OrganizationUnitRegistry through an OUChildBlockerChecker, mirroring the
// production wiring in cmd/sovrunn-api/main.go. This lets HTTP-level tests
// exercise Organization delete blocking end-to-end.
func newBlockerWiring() (*OrgHandler, *OUHandler, *registry.OrganizationRegistry) {
	orgRegistry := registry.NewOrganizationRegistry()
	ouRegistry := registry.NewOrganizationUnitRegistry()
	ouBlocker := registry.NewOUChildBlockerChecker(ouRegistry)
	orgHandler := NewOrgHandler(orgRegistry, ouBlocker)
	ouHandler := NewOUHandler(ouRegistry, orgRegistry, nil)
	return orgHandler, ouHandler, orgRegistry
}

func TestOrgDeleteBlockedByOU(t *testing.T) {
	orgHandler, ouHandler, orgReg := newBlockerWiring()
	seedOrg(t, orgReg, "nic")
	if rec := createOU(ouHandler, "nic", "ministry-health", ""); rec.Code != http.StatusCreated {
		t.Fatalf("create OU status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}

	req := jsonRequest(http.MethodDelete, "/v1/organizations/nic", nil, "")
	rec := httptest.NewRecorder()
	orgHandler.HandleItem(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeDeleteBlocked {
		t.Errorf("code = %q, want DELETE_BLOCKED", errBody.Code)
	}
	if !strings.Contains(errBody.Message, "OrganizationUnit") {
		t.Errorf("message = %q, want it to identify OrganizationUnit", errBody.Message)
	}
}

func TestOrgDeleteAllowedWhenNoOUs(t *testing.T) {
	orgHandler, _, orgReg := newBlockerWiring()
	seedOrg(t, orgReg, "empty-org")

	req := jsonRequest(http.MethodDelete, "/v1/organizations/empty-org", nil, "")
	rec := httptest.NewRecorder()
	orgHandler.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
}

func TestOUHandler_Create_Valid(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")

	rec := createOU(h, "nic", "ministry-health", "Health")
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var ou resources.OrganizationUnit
	if err := json.NewDecoder(rec.Body).Decode(&ou); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if ou.APIVersion != resources.OUAPIVersion || ou.Kind != resources.OUKind {
		t.Errorf("apiVersion/kind not set by server: %+v", ou)
	}
	if ou.Status.Phase != resources.PhaseActive {
		t.Errorf("phase = %q, want Active", ou.Status.Phase)
	}
	if ou.Metadata.Name != "ministry-health" || ou.Spec.OrganizationName != "nic" {
		t.Errorf("unexpected resource: %+v", ou)
	}
}

func TestOUHandler_Create_Duplicate(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")

	if rec := createOU(h, "nic", "ministry-health", ""); rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d", rec.Code)
	}
	rec := createOU(h, "nic", "ministry-health", "")
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceAlreadyExists {
		t.Errorf("code = %q, want RESOURCE_ALREADY_EXISTS", errBody.Code)
	}
}

func TestOUHandler_Create_InvalidName(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")

	rec := createOU(h, "nic", "INVALID", "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "metadata.name" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestOUHandler_Create_MissingOrganizationName(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")

	body := map[string]any{
		"metadata": map[string]any{"name": "ministry-health"},
		"spec":     map[string]any{"description": "no org"},
	}
	req := jsonRequest(http.MethodPost, "/v1/organization-units", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.organizationName" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestOUHandler_Create_NonExistentParent(t *testing.T) {
	h, _, _ := newTestOUHandler()

	rec := createOU(h, "ghost", "ministry-health", "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.organizationName" {
		t.Errorf("error = %+v, want VALIDATION_FAILED spec.organizationName", errBody)
	}
}

func TestOUHandler_Create_StatusFieldRejected(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")

	cases := []string{
		`{"metadata":{"name":"ministry-health"},"spec":{"organizationName":"nic"},"status":{}}`,
		`{"metadata":{"name":"ministry-health"},"spec":{"organizationName":"nic"},"status":null}`,
		`{"metadata":{"name":"ministry-health"},"spec":{"organizationName":"nic"},"status":{"phase":""}}`,
	}
	for _, payload := range cases {
		req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/organization-units", strings.NewReader(payload)), "id")
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.HandleCollection(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("payload %s: status = %d, want 400", payload, rec.Code)
		}
		errBody := decodeAPIError(t, rec)
		if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "status" {
			t.Errorf("payload %s: error = %+v, want VALIDATION_FAILED status", payload, errBody)
		}
	}
}

func TestOUHandler_Create_BadJSON(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")

	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/organization-units", strings.NewReader("{")), "id")
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

func TestOUHandler_Create_OversizedBody(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")

	large := strings.Repeat("a", 1<<20+1)
	payload := fmt.Sprintf(`{"metadata":{"name":"ministry-health"},"spec":{"organizationName":"nic","description":"%s"}}`, large)
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/organization-units", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

// NOTE: POST/PUT 415 (Unsupported Media Type) is produced by
// contentTypeMiddleware in package server, not by OUHandler. Testing it here
// would require importing package server (which imports api), creating an
// import cycle. 415 behavior is therefore covered at the server-wiring level
// in Task 9/10, consistent with org_handler_test.go which also defers 415.

func TestOUHandler_Get_Exists(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")
	if rec := createOU(h, "nic", "ministry-health", "Health"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	req := jsonRequest(http.MethodGet, "/v1/organization-units/nic/ministry-health", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var ou resources.OrganizationUnit
	if err := json.NewDecoder(rec.Body).Decode(&ou); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if ou.Metadata.Name != "ministry-health" || ou.Spec.OrganizationName != "nic" ||
		ou.APIVersion == "" || ou.Kind == "" || ou.Status.Phase == "" {
		t.Errorf("incomplete resource: %+v", ou)
	}
}

func TestOUHandler_Get_NotFound(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")

	req := jsonRequest(http.MethodGet, "/v1/organization-units/nic/missing", nil, "")
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

func TestOUHandler_Get_InvalidPathSegments(t *testing.T) {
	h, _, _ := newTestOUHandler()

	req := jsonRequest(http.MethodGet, "/v1/organization-units/nic/INVALID", nil, "")
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

func TestOUHandler_Get_BareItemPath_NotFound(t *testing.T) {
	h, _, _ := newTestOUHandler()

	req := jsonRequest(http.MethodGet, "/v1/organization-units/", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestOUHandler_Get_SingleSegment_NotFound(t *testing.T) {
	h, _, _ := newTestOUHandler()

	req := jsonRequest(http.MethodGet, "/v1/organization-units/nic", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestOUHandler_Get_ExtraSegment_NotFound(t *testing.T) {
	h, _, _ := newTestOUHandler()

	req := jsonRequest(http.MethodGet, "/v1/organization-units/nic/ministry-health/extra", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestOUHandler_List_Sorted(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "alpha")
	seedOrg(t, orgReg, "zebra")

	inputs := []struct{ org, name string }{
		{"zebra", "beta"},
		{"alpha", "delta"},
		{"alpha", "charlie"},
	}
	for _, in := range inputs {
		if rec := createOU(h, in.org, in.name, ""); rec.Code != http.StatusCreated {
			t.Fatalf("create %s/%s status = %d", in.org, in.name, rec.Code)
		}
	}

	req := jsonRequest(http.MethodGet, "/v1/organization-units", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp organizationUnitListResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Items) != 3 {
		t.Fatalf("items = %d, want 3", len(resp.Items))
	}
	want := []struct{ org, name string }{
		{"alpha", "charlie"},
		{"alpha", "delta"},
		{"zebra", "beta"},
	}
	for i, w := range want {
		if resp.Items[i].Spec.OrganizationName != w.org || resp.Items[i].Metadata.Name != w.name {
			t.Errorf("item[%d] = %s/%s, want %s/%s", i,
				resp.Items[i].Spec.OrganizationName, resp.Items[i].Metadata.Name, w.org, w.name)
		}
	}
}

func TestOUHandler_List_Empty(t *testing.T) {
	h, _, _ := newTestOUHandler()

	req := jsonRequest(http.MethodGet, "/v1/organization-units", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"items":[]}` {
		t.Errorf("body = %q, want {\"items\":[]}", got)
	}
}

func TestOUHandler_Update_Valid(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")
	if rec := createOU(h, "nic", "ministry-health", "old"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	body := map[string]any{
		"metadata": map[string]any{"name": "ministry-health", "displayName": "New Display"},
		"spec":     map[string]any{"organizationName": "nic", "description": "updated"},
	}
	req := jsonRequest(http.MethodPut, "/v1/organization-units/nic/ministry-health", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var ou resources.OrganizationUnit
	if err := json.NewDecoder(rec.Body).Decode(&ou); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if ou.Metadata.DisplayName != "New Display" || ou.Spec.Description != "updated" {
		t.Errorf("mutable fields not updated: %+v", ou)
	}
	if ou.Metadata.Name != "ministry-health" || ou.Spec.OrganizationName != "nic" {
		t.Errorf("immutable fields changed: %+v", ou)
	}
}

func TestOUHandler_Update_NotFound(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")

	body := map[string]any{
		"metadata": map[string]any{"name": "missing"},
		"spec":     map[string]any{"organizationName": "nic"},
	}
	req := jsonRequest(http.MethodPut, "/v1/organization-units/nic/missing", body, "application/json")
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

func TestOUHandler_Update_NameAbsent(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")
	_ = createOU(h, "nic", "ministry-health", "")

	body := map[string]any{
		"metadata": map[string]any{},
		"spec":     map[string]any{"organizationName": "nic"},
	}
	req := jsonRequest(http.MethodPut, "/v1/organization-units/nic/ministry-health", body, "application/json")
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

func TestOUHandler_Update_NameMismatch(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")
	_ = createOU(h, "nic", "ministry-health", "")

	body := map[string]any{
		"metadata": map[string]any{"name": "other-name"},
		"spec":     map[string]any{"organizationName": "nic"},
	}
	req := jsonRequest(http.MethodPut, "/v1/organization-units/nic/ministry-health", body, "application/json")
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

func TestOUHandler_Update_OrganizationNameAbsent(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")
	_ = createOU(h, "nic", "ministry-health", "")

	body := map[string]any{
		"metadata": map[string]any{"name": "ministry-health"},
		"spec":     map[string]any{"description": "no org"},
	}
	req := jsonRequest(http.MethodPut, "/v1/organization-units/nic/ministry-health", body, "application/json")
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

func TestOUHandler_Update_OrganizationNameMismatch(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")
	_ = createOU(h, "nic", "ministry-health", "")

	body := map[string]any{
		"metadata": map[string]any{"name": "ministry-health"},
		"spec":     map[string]any{"organizationName": "other-org"},
	}
	req := jsonRequest(http.MethodPut, "/v1/organization-units/nic/ministry-health", body, "application/json")
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

func TestOUHandler_Update_StatusField(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")
	_ = createOU(h, "nic", "ministry-health", "")

	payload := `{"metadata":{"name":"ministry-health"},"spec":{"organizationName":"nic"},"status":{}}`
	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/organization-units/nic/ministry-health", strings.NewReader(payload)), "id")
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

func TestOUHandler_Update_BadJSON_TargetExists(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")
	_ = createOU(h, "nic", "ministry-health", "")

	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/organization-units/nic/ministry-health", strings.NewReader("{")), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestOUHandler_Update_OversizedBody(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")
	_ = createOU(h, "nic", "ministry-health", "")

	large := strings.Repeat("a", 1<<20+1)
	payload := fmt.Sprintf(`{"metadata":{"name":"ministry-health"},"spec":{"organizationName":"nic","description":"%s"}}`, large)
	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/organization-units/nic/ministry-health", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

func TestOUHandler_Delete_Success(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")
	_ = createOU(h, "nic", "ministry-health", "")

	req := jsonRequest(http.MethodDelete, "/v1/organization-units/nic/ministry-health", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

func TestOUHandler_Delete_BlockerReturningTenantBlocksDelete(t *testing.T) {
	orgRegistry := registry.NewOrganizationRegistry()
	ouRegistry := registry.NewOrganizationUnitRegistry()
	h := NewOUHandler(ouRegistry, orgRegistry, staticOUBlocker{
		blockers: []registry.BlockedBy{{Kind: "Tenant", Count: 2}},
	})
	seedOrg(t, orgRegistry, "nic")
	_ = createOU(h, "nic", "ministry-health", "")

	req := jsonRequest(http.MethodDelete, "/v1/organization-units/nic/ministry-health", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
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
	if _, err := ouRegistry.GetOrganizationUnit(context.Background(), "nic", "ministry-health"); err != nil {
		t.Fatalf("blocked delete removed resource, GetOrganizationUnit() error = %v", err)
	}
}

func TestOUHandler_Delete_EmptyBlockerAllowsDelete(t *testing.T) {
	orgRegistry := registry.NewOrganizationRegistry()
	ouRegistry := registry.NewOrganizationUnitRegistry()
	h := NewOUHandler(ouRegistry, orgRegistry, staticOUBlocker{})
	seedOrg(t, orgRegistry, "nic")
	_ = createOU(h, "nic", "ministry-health", "")

	req := jsonRequest(http.MethodDelete, "/v1/organization-units/nic/ministry-health", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
}

func TestOUHandler_Delete_NotFound(t *testing.T) {
	h, orgReg, _ := newTestOUHandler()
	seedOrg(t, orgReg, "nic")

	req := jsonRequest(http.MethodDelete, "/v1/organization-units/nic/missing", nil, "")
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

func TestOUHandler_Delete_InvalidPathSegments(t *testing.T) {
	h, _, _ := newTestOUHandler()

	req := jsonRequest(http.MethodDelete, "/v1/organization-units/nic/BAD", nil, "")
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

type staticOUBlocker struct {
	blockers []registry.BlockedBy
	err      error
}

func (b staticOUBlocker) BlockedByOUChildren(
	ctx context.Context, orgName, ouName string,
) ([]registry.BlockedBy, error) {
	return b.blockers, b.err
}
