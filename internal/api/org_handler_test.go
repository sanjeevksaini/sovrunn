package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/requestctx"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func newTestHandler() *OrgHandler {
	return NewOrgHandler(registry.NewOrganizationRegistry(), registry.NoopChildBlockerChecker{}, nil)
}

func withRequestID(r *http.Request, id string) *http.Request {
	return r.WithContext(requestctx.WithRequestID(r.Context(), id))
}

func jsonRequest(method, path string, body any, contentType string) *http.Request {
	var reader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, path, reader)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return withRequestID(req, "test-req-id")
}

func decodeAPIError(t *testing.T, rec *httptest.ResponseRecorder) resources.APIError {
	t.Helper()
	var envelope resources.APIErrorEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	return envelope.Error
}

func TestOrgHandler_Create_Valid(t *testing.T) {
	h := newTestHandler()
	body := map[string]any{
		"metadata": map[string]any{"name": "nic"},
		"spec":     map[string]any{"description": "NIC"},
	}
	req := jsonRequest(http.MethodPost, "/v1/organizations", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var org resources.Organization
	if err := json.NewDecoder(rec.Body).Decode(&org); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if org.APIVersion != resources.OrgAPIVersion || org.Kind != resources.OrgKind {
		t.Errorf("apiVersion/kind not set by server")
	}
	if org.Status.Phase != resources.PhaseActive {
		t.Errorf("phase = %q, want Active", org.Status.Phase)
	}
}

func TestOrgHandler_Create_Duplicate(t *testing.T) {
	h := newTestHandler()
	body := map[string]any{"metadata": map[string]any{"name": "nic"}}
	for i := 0; i < 2; i++ {
		req := jsonRequest(http.MethodPost, "/v1/organizations", body, "application/json")
		rec := httptest.NewRecorder()
		h.HandleCollection(rec, req)
		if i == 0 && rec.Code != http.StatusCreated {
			t.Fatalf("first create status = %d", rec.Code)
		}
		if i == 1 {
			if rec.Code != http.StatusConflict {
				t.Fatalf("status = %d, want 409", rec.Code)
			}
			errBody := decodeAPIError(t, rec)
			if errBody.Code != resources.ErrCodeResourceAlreadyExists {
				t.Errorf("code = %q", errBody.Code)
			}
		}
	}
}

func TestOrgHandler_Create_InvalidName(t *testing.T) {
	h := newTestHandler()
	body := map[string]any{"metadata": map[string]any{"name": "INVALID"}}
	req := jsonRequest(http.MethodPost, "/v1/organizations", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "metadata.name" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestOrgHandler_Create_StatusFieldRejected(t *testing.T) {
	h := newTestHandler()
	cases := []string{
		`{"metadata":{"name":"nic"},"status":{}}`,
		`{"metadata":{"name":"nic"},"status":null}`,
		`{"metadata":{"name":"nic"},"status":{"phase":""}}`,
	}
	for _, payload := range cases {
		req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/organizations", strings.NewReader(payload)), "id")
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.HandleCollection(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("payload %s: status = %d, want 400", payload, rec.Code)
		}
		errBody := decodeAPIError(t, rec)
		if errBody.Field != "status" {
			t.Errorf("field = %q, want status", errBody.Field)
		}
	}
}

func TestOrgHandler_Create_BadJSON(t *testing.T) {
	h := newTestHandler()
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/organizations", strings.NewReader("{")), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestOrgHandler_Create_OversizedBody(t *testing.T) {
	h := newTestHandler()
	large := strings.Repeat("a", 1<<20+1)
	payload := fmt.Sprintf(`{"metadata":{"name":"nic"},"spec":{"description":"%s"}}`, large)
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/organizations", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

func TestOrgHandler_Get_Exists(t *testing.T) {
	h := newTestHandler()
	create := jsonRequest(http.MethodPost, "/v1/organizations", map[string]any{
		"metadata": map[string]any{"name": "nic"},
	}, "application/json")
	h.HandleCollection(httptest.NewRecorder(), create)

	req := jsonRequest(http.MethodGet, "/v1/organizations/nic", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var org resources.Organization
	_ = json.NewDecoder(rec.Body).Decode(&org)
	if org.Metadata.Name != "nic" || org.APIVersion == "" || org.Kind == "" || org.Status.Phase == "" {
		t.Errorf("incomplete resource: %+v", org)
	}
}

func TestOrgHandler_Get_NotFound(t *testing.T) {
	h := newTestHandler()
	req := jsonRequest(http.MethodGet, "/v1/organizations/missing", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestOrgHandler_Get_InvalidPathName(t *testing.T) {
	h := newTestHandler()
	req := jsonRequest(http.MethodGet, "/v1/organizations/INVALID", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestOrgHandler_List_Sorted(t *testing.T) {
	h := newTestHandler()
	for _, name := range []string{"zebra", "alpha", "mike"} {
		req := jsonRequest(http.MethodPost, "/v1/organizations", map[string]any{
			"metadata": map[string]any{"name": name},
		}, "application/json")
		h.HandleCollection(httptest.NewRecorder(), req)
	}
	req := jsonRequest(http.MethodGet, "/v1/organizations", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp organizationListResponse
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if len(resp.Items) != 3 {
		t.Fatalf("items = %d, want 3", len(resp.Items))
	}
	for i := 1; i < len(resp.Items); i++ {
		if resp.Items[i-1].Metadata.Name >= resp.Items[i].Metadata.Name {
			t.Fatalf("not sorted: %+v", resp.Items)
		}
	}
}

func TestOrgHandler_List_Empty(t *testing.T) {
	h := newTestHandler()
	req := jsonRequest(http.MethodGet, "/v1/organizations", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	var resp organizationListResponse
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Items == nil || len(resp.Items) != 0 {
		t.Fatalf("items = %+v, want empty slice", resp.Items)
	}
}

func TestOrgHandler_Update_Valid(t *testing.T) {
	h := newTestHandler()
	h.HandleCollection(httptest.NewRecorder(), jsonRequest(http.MethodPost, "/v1/organizations", map[string]any{
		"metadata": map[string]any{"name": "nic"},
	}, "application/json"))

	body := map[string]any{
		"metadata": map[string]any{"name": "nic", "displayName": "New"},
		"spec":     map[string]any{"description": "updated"},
	}
	req := jsonRequest(http.MethodPut, "/v1/organizations/nic", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestOrgHandler_Update_NotFound(t *testing.T) {
	h := newTestHandler()
	body := map[string]any{"metadata": map[string]any{"name": "missing"}}
	req := jsonRequest(http.MethodPut, "/v1/organizations/missing", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestOrgHandler_Update_NameMismatch(t *testing.T) {
	h := newTestHandler()
	h.HandleCollection(httptest.NewRecorder(), jsonRequest(http.MethodPost, "/v1/organizations", map[string]any{
		"metadata": map[string]any{"name": "nic"},
	}, "application/json"))
	body := map[string]any{"metadata": map[string]any{"name": "other"}}
	req := jsonRequest(http.MethodPut, "/v1/organizations/nic", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "metadata.name" {
		t.Errorf("field = %q", errBody.Field)
	}
}

func TestOrgHandler_Update_NameAbsent(t *testing.T) {
	h := newTestHandler()
	h.HandleCollection(httptest.NewRecorder(), jsonRequest(http.MethodPost, "/v1/organizations", map[string]any{
		"metadata": map[string]any{"name": "nic"},
	}, "application/json"))
	body := map[string]any{"metadata": map[string]any{}}
	req := jsonRequest(http.MethodPut, "/v1/organizations/nic", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestOrgHandler_Update_StatusField(t *testing.T) {
	h := newTestHandler()
	h.HandleCollection(httptest.NewRecorder(), jsonRequest(http.MethodPost, "/v1/organizations", map[string]any{
		"metadata": map[string]any{"name": "nic"},
	}, "application/json"))
	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/organizations/nic",
		strings.NewReader(`{"metadata":{"name":"nic"},"status":{}}`)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestOrgHandler_Delete_Success(t *testing.T) {
	h := newTestHandler()
	h.HandleCollection(httptest.NewRecorder(), jsonRequest(http.MethodPost, "/v1/organizations", map[string]any{
		"metadata": map[string]any{"name": "nic"},
	}, "application/json"))
	req := jsonRequest(http.MethodDelete, "/v1/organizations/nic", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
}

func TestOrgHandler_Delete_NotFound(t *testing.T) {
	h := newTestHandler()
	req := jsonRequest(http.MethodDelete, "/v1/organizations/missing", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestOrgHandler_Delete_InvalidPathName(t *testing.T) {
	h := newTestHandler()
	req := jsonRequest(http.MethodDelete, "/v1/organizations/BAD", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestOrgHandler_BareItemPath_NotFound(t *testing.T) {
	h := newTestHandler()
	req := jsonRequest(http.MethodGet, "/v1/organizations/", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}
