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

func newTestServiceClassHandler() *ServiceClassHandler {
	return NewServiceClassHandler(registry.NewServiceClassRegistry(), nil, nil)
}

func newServiceClassHandlerWithBlocker(blocker registry.ServiceClassChildBlocker) *ServiceClassHandler {
	return NewServiceClassHandler(registry.NewServiceClassRegistry(), blocker, nil)
}

func validServiceClassBody(name string) map[string]any {
	return map[string]any{
		"metadata": map[string]any{"name": name},
		"spec": map[string]any{
			"category":  resources.CategoryDatabase,
			"lifecycle": resources.LifecycleActive,
		},
	}
}

func createServiceClass(h *ServiceClassHandler, name string) *httptest.ResponseRecorder {
	req := jsonRequest(http.MethodPost, "/v1/service-classes", validServiceClassBody(name), "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	return rec
}

func TestServiceClassHandler_Create_Valid(t *testing.T) {
	h := newTestServiceClassHandler()
	rec := createServiceClass(h, "postgres")
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var sc resources.ServiceClass
	if err := json.NewDecoder(rec.Body).Decode(&sc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sc.APIVersion != resources.ServiceClassAPIVersion || sc.Kind != resources.ServiceClassKind {
		t.Errorf("apiVersion/kind = %q/%q, want server-owned", sc.APIVersion, sc.Kind)
	}
	if sc.Metadata.Name != "postgres" {
		t.Errorf("name = %q, want postgres", sc.Metadata.Name)
	}
	if sc.Status.Phase != resources.PhaseActive {
		t.Errorf("phase = %q, want Active", sc.Status.Phase)
	}
}

func TestServiceClassHandler_Create_Duplicate(t *testing.T) {
	h := newTestServiceClassHandler()
	if rec := createServiceClass(h, "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d", rec.Code)
	}
	rec := createServiceClass(h, "postgres")
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceAlreadyExists {
		t.Errorf("code = %q, want RESOURCE_ALREADY_EXISTS", errBody.Code)
	}
}

func TestServiceClassHandler_Create_InvalidFields(t *testing.T) {
	h := newTestServiceClassHandler()
	body := map[string]any{
		"metadata": map[string]any{"name": "INVALID"},
		"spec": map[string]any{
			"category":  resources.CategoryDatabase,
			"lifecycle": resources.LifecycleActive,
		},
	}
	req := jsonRequest(http.MethodPost, "/v1/service-classes", body, "application/json")
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

func TestServiceClassHandler_Create_StatusFieldRejected(t *testing.T) {
	h := newTestServiceClassHandler()
	payload := `{"metadata":{"name":"postgres"},"spec":{"category":"Database","lifecycle":"Active"},"status":{}}`
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/service-classes", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "status" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestServiceClassHandler_Create_BadJSON(t *testing.T) {
	h := newTestServiceClassHandler()
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/service-classes", strings.NewReader("{")), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestServiceClassHandler_Create_UnknownField(t *testing.T) {
	h := newTestServiceClassHandler()
	payload := `{"metadata":{"name":"postgres"},"spec":{"category":"Database","lifecycle":"Active"},"bogus":true}`
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/service-classes", strings.NewReader(payload)), "id")
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

func TestServiceClassHandler_Create_OversizedBody(t *testing.T) {
	h := newTestServiceClassHandler()
	large := strings.Repeat("a", 1<<20+1)
	payload := fmt.Sprintf(
		`{"metadata":{"name":"postgres"},"spec":{"category":"Database","lifecycle":"Active","description":"%s"}}`,
		large,
	)
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/service-classes", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

func TestServiceClassHandler_Get_Exists(t *testing.T) {
	h := newTestServiceClassHandler()
	if rec := createServiceClass(h, "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	req := jsonRequest(http.MethodGet, "/v1/service-classes/postgres", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var sc resources.ServiceClass
	if err := json.NewDecoder(rec.Body).Decode(&sc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sc.Metadata.Name != "postgres" || sc.APIVersion == "" || sc.Kind == "" || sc.Status.Phase == "" {
		t.Errorf("incomplete resource: %+v", sc)
	}
}

func TestServiceClassHandler_Get_NotFound(t *testing.T) {
	h := newTestServiceClassHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-classes/missing", nil, "")
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

func TestServiceClassHandler_Get_InvalidPathSegment(t *testing.T) {
	h := newTestServiceClassHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-classes/INVALID", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "metadata.name" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestServiceClassHandler_Get_WrongPathShape(t *testing.T) {
	h := newTestServiceClassHandler()
	for _, path := range []string{"/v1/service-classes/", "/v1/service-classes/postgres/extra"} {
		req := jsonRequest(http.MethodGet, path, nil, "")
		rec := httptest.NewRecorder()
		h.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("path %s: status = %d, want 404", path, rec.Code)
		}
	}
}

func TestServiceClassHandler_List_Empty(t *testing.T) {
	h := newTestServiceClassHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-classes", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"items":[]}` {
		t.Errorf("body = %q, want {\"items\":[]}", got)
	}
}

func TestServiceClassHandler_List_Sorted(t *testing.T) {
	h := newTestServiceClassHandler()
	for _, name := range []string{"zebra", "alpha", "mongo"} {
		if rec := createServiceClass(h, name); rec.Code != http.StatusCreated {
			t.Fatalf("create %s status = %d", name, rec.Code)
		}
	}
	req := jsonRequest(http.MethodGet, "/v1/service-classes", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var top map[string]json.RawMessage
	if err := json.NewDecoder(rec.Body).Decode(&top); err != nil {
		t.Fatalf("decode top: %v", err)
	}
	if len(top) != 1 {
		t.Fatalf("top-level keys = %d, want only items", len(top))
	}
	var items []resources.ServiceClass
	if err := json.Unmarshal(top["items"], &items); err != nil {
		t.Fatalf("decode items: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("items = %d, want 3", len(items))
	}
	for i := 1; i < len(items); i++ {
		if items[i-1].Metadata.Name >= items[i].Metadata.Name {
			t.Fatalf("not sorted: %+v", items)
		}
	}
}

func TestServiceClassHandler_Update_Valid(t *testing.T) {
	h := newTestServiceClassHandler()
	if rec := createServiceClass(h, "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := map[string]any{
		"metadata": map[string]any{"name": "postgres"},
		"spec": map[string]any{
			"displayName": "PostgreSQL",
			"description": "updated",
			"category":    resources.CategoryDatabase,
			"lifecycle":   resources.LifecycleDeprecated,
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/service-classes/postgres", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var sc resources.ServiceClass
	if err := json.NewDecoder(rec.Body).Decode(&sc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sc.Spec.Description != "updated" || sc.Spec.Lifecycle != resources.LifecycleDeprecated {
		t.Errorf("mutable fields not updated: %+v", sc.Spec)
	}
}

func TestServiceClassHandler_Update_NotFound(t *testing.T) {
	h := newTestServiceClassHandler()
	req := jsonRequest(http.MethodPut, "/v1/service-classes/missing", validServiceClassBody("missing"), "application/json")
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

func TestServiceClassHandler_Update_NameAbsent(t *testing.T) {
	h := newTestServiceClassHandler()
	if rec := createServiceClass(h, "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := map[string]any{
		"metadata": map[string]any{},
		"spec": map[string]any{
			"category":  resources.CategoryDatabase,
			"lifecycle": resources.LifecycleActive,
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/service-classes/postgres", body, "application/json")
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

func TestServiceClassHandler_Update_NameMismatch(t *testing.T) {
	h := newTestServiceClassHandler()
	if rec := createServiceClass(h, "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := validServiceClassBody("other")
	req := jsonRequest(http.MethodPut, "/v1/service-classes/postgres", body, "application/json")
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

func TestServiceClassHandler_Update_InvalidFields(t *testing.T) {
	h := newTestServiceClassHandler()
	if rec := createServiceClass(h, "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := map[string]any{
		"metadata": map[string]any{"name": "postgres"},
		"spec": map[string]any{
			"category":  "NotACategory",
			"lifecycle": resources.LifecycleActive,
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/service-classes/postgres", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.category" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestServiceClassHandler_Update_PreservesServerOwnedFields(t *testing.T) {
	h := newTestServiceClassHandler()
	if rec := createServiceClass(h, "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	// status key is rejected before update; use a body without status but with
	// tampered apiVersion/kind to prove registry preservation.
	payload := `{
		"apiVersion":"tampered/v0",
		"kind":"Tampered",
		"metadata":{"name":"postgres"},
		"spec":{"category":"Database","lifecycle":"Active","description":"changed"}
	}`
	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/service-classes/postgres", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var sc resources.ServiceClass
	if err := json.NewDecoder(rec.Body).Decode(&sc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sc.APIVersion != resources.ServiceClassAPIVersion {
		t.Errorf("APIVersion = %q, want preserved", sc.APIVersion)
	}
	if sc.Kind != resources.ServiceClassKind {
		t.Errorf("Kind = %q, want ServiceClass", sc.Kind)
	}
	if sc.Status.Phase != resources.PhaseActive {
		t.Errorf("Status.Phase = %q, want Active", sc.Status.Phase)
	}
	if sc.Metadata.Name != "postgres" {
		t.Errorf("Metadata.Name = %q, want postgres", sc.Metadata.Name)
	}
	if sc.Spec.Description != "changed" {
		t.Errorf("Description = %q, want changed", sc.Spec.Description)
	}
}

func TestServiceClassHandler_Delete_Success(t *testing.T) {
	h := newTestServiceClassHandler()
	if rec := createServiceClass(h, "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	req := jsonRequest(http.MethodDelete, "/v1/service-classes/postgres", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

func TestServiceClassHandler_Delete_NotFound(t *testing.T) {
	h := newTestServiceClassHandler()
	req := jsonRequest(http.MethodDelete, "/v1/service-classes/missing", nil, "")
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

func TestServiceClassHandler_Delete_InvalidPathSegment(t *testing.T) {
	h := newTestServiceClassHandler()
	req := jsonRequest(http.MethodDelete, "/v1/service-classes/BAD", nil, "")
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

func TestServiceClassHandler_Delete_WrongPathShape(t *testing.T) {
	h := newTestServiceClassHandler()
	req := jsonRequest(http.MethodDelete, "/v1/service-classes/postgres/extra", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestServiceClassHandler_Delete_BlockedByServicePlan(t *testing.T) {
	scReg := registry.NewServiceClassRegistry()
	spReg := registry.NewServicePlanRegistry()
	blocker := registry.NewServicePlanChildBlockerChecker(spReg)
	h := NewServiceClassHandler(scReg, blocker, nil)

	if rec := createServiceClass(h, "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create ServiceClass status = %d", rec.Code)
	}
	plan := resources.ServicePlan{
		APIVersion: "platform.sovrunn.io/v1alpha1",
		Kind:       resources.ServicePlanKind,
		Metadata:   resources.Metadata{Name: "small"},
		Spec: resources.ServicePlanSpec{
			ServiceClassName: "postgres",
			Tier:             resources.TierSmall,
			Lifecycle:        resources.LifecycleActive,
		},
		Status: resources.ServicePlanStatus{Phase: resources.PhaseActive},
	}
	if _, err := spReg.CreateServicePlan(context.Background(), plan); err != nil {
		t.Fatalf("seed ServicePlan: %v", err)
	}

	req := jsonRequest(http.MethodDelete, "/v1/service-classes/postgres", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeDeleteBlocked {
		t.Errorf("code = %q, want DELETE_BLOCKED", errBody.Code)
	}
	if !strings.Contains(errBody.Message, "ServicePlan") {
		t.Errorf("message = %q, want it to mention ServicePlan", errBody.Message)
	}
}

func TestServiceClassHandler_Delete_NilBlockerAllows(t *testing.T) {
	h := newServiceClassHandlerWithBlocker(nil)
	if rec := createServiceClass(h, "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	req := jsonRequest(http.MethodDelete, "/v1/service-classes/postgres", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
}

func TestServiceClassHandler_NilEmitterNoPanic(t *testing.T) {
	h := NewServiceClassHandler(registry.NewServiceClassRegistry(), nil, nil)
	if rec := createServiceClass(h, "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := validServiceClassBody("postgres")
	req := jsonRequest(http.MethodPut, "/v1/service-classes/postgres", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200", rec.Code)
	}
	del := jsonRequest(http.MethodDelete, "/v1/service-classes/postgres", nil, "")
	delRec := httptest.NewRecorder()
	h.HandleItem(delRec, del)
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204", delRec.Code)
	}
}

func TestServiceClassHandler_UnsupportedMethods(t *testing.T) {
	h := newTestServiceClassHandler()
	req := jsonRequest(http.MethodPatch, "/v1/service-classes", nil, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("collection status = %d, want 405", rec.Code)
	}

	itemReq := jsonRequest(http.MethodPatch, "/v1/service-classes/postgres", nil, "application/json")
	itemRec := httptest.NewRecorder()
	h.HandleItem(itemRec, itemReq)
	if itemRec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("item status = %d, want 405", itemRec.Code)
	}
}
