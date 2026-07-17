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

func newTestServicePlanHandler() (*ServicePlanHandler, *registry.ServiceClassRegistry, *registry.ServicePlanRegistry) {
	scReg := registry.NewServiceClassRegistry()
	spReg := registry.NewServicePlanRegistry()
	h := NewServicePlanHandler(spReg, scReg, nil)
	return h, scReg, spReg
}

func seedServiceClass(t *testing.T, reg *registry.ServiceClassRegistry, name string) {
	t.Helper()
	sc := resources.ServiceClass{
		APIVersion: resources.ServiceClassAPIVersion,
		Kind:       resources.ServiceClassKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.ServiceClassSpec{
			Category:  resources.CategoryDatabase,
			Lifecycle: resources.LifecycleActive,
		},
		Status: resources.ServiceClassStatus{Phase: resources.PhaseActive},
	}
	if _, err := reg.CreateServiceClass(context.Background(), sc); err != nil {
		t.Fatalf("seedServiceClass(%s): %v", name, err)
	}
}

func validServicePlanBody(serviceClassName, name string) map[string]any {
	return map[string]any{
		"metadata": map[string]any{"name": name},
		"spec": map[string]any{
			"serviceClassName": serviceClassName,
			"tier":             resources.TierSmall,
			"lifecycle":        resources.LifecycleActive,
		},
	}
}

func createServicePlan(h *ServicePlanHandler, serviceClassName, name string) *httptest.ResponseRecorder {
	req := jsonRequest(http.MethodPost, "/v1/service-plans", validServicePlanBody(serviceClassName, name), "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	return rec
}

func TestServicePlanHandler_Create_Valid(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")

	rec := createServicePlan(h, "postgres", "small")
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var sp resources.ServicePlan
	if err := json.NewDecoder(rec.Body).Decode(&sp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sp.APIVersion != resources.ServicePlanAPIVersion || sp.Kind != resources.ServicePlanKind {
		t.Errorf("apiVersion/kind = %q/%q, want server-owned", sp.APIVersion, sp.Kind)
	}
	if sp.Metadata.Name != "small" || sp.Spec.ServiceClassName != "postgres" {
		t.Errorf("identity = %s/%s, want postgres/small", sp.Spec.ServiceClassName, sp.Metadata.Name)
	}
	if sp.Status.Phase != resources.PhaseActive {
		t.Errorf("phase = %q, want Active", sp.Status.Phase)
	}
}

func TestServicePlanHandler_Create_Duplicate(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d", rec.Code)
	}
	rec := createServicePlan(h, "postgres", "small")
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceAlreadyExists {
		t.Errorf("code = %q, want RESOURCE_ALREADY_EXISTS", errBody.Code)
	}
}

func TestServicePlanHandler_Create_SameNameDifferentClasses(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	seedServiceClass(t, scReg, "redis")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create postgres/small status = %d", rec.Code)
	}
	if rec := createServicePlan(h, "redis", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create redis/small status = %d", rec.Code)
	}
}

func TestServicePlanHandler_Create_InvalidFields(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	body := map[string]any{
		"metadata": map[string]any{"name": "INVALID"},
		"spec": map[string]any{
			"serviceClassName": "postgres",
			"tier":             resources.TierSmall,
			"lifecycle":        resources.LifecycleActive,
		},
	}
	req := jsonRequest(http.MethodPost, "/v1/service-plans", body, "application/json")
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

func TestServicePlanHandler_Create_MissingParent(t *testing.T) {
	h, _, _ := newTestServicePlanHandler()
	rec := createServicePlan(h, "postgres", "small")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.serviceClassName" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestServicePlanHandler_Create_StatusFieldRejected(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	payload := `{"metadata":{"name":"small"},"spec":{"serviceClassName":"postgres","tier":"Small","lifecycle":"Active"},"status":{}}`
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/service-plans", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "status" {
		t.Errorf("field = %q, want status", errBody.Field)
	}
}

func TestServicePlanHandler_Create_BadJSON(t *testing.T) {
	h, _, _ := newTestServicePlanHandler()
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/service-plans", strings.NewReader("{")), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestServicePlanHandler_Create_UnknownField(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	payload := `{"metadata":{"name":"small"},"spec":{"serviceClassName":"postgres","tier":"Small","lifecycle":"Active"},"bogus":true}`
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/service-plans", strings.NewReader(payload)), "id")
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

func TestServicePlanHandler_Create_OversizedBody(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	large := strings.Repeat("a", 1<<20+1)
	payload := fmt.Sprintf(
		`{"metadata":{"name":"small"},"spec":{"serviceClassName":"postgres","tier":"Small","lifecycle":"Active","description":"%s"}}`,
		large,
	)
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/service-plans", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

func TestServicePlanHandler_Get_Exists(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	req := jsonRequest(http.MethodGet, "/v1/service-plans/postgres/small", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var sp resources.ServicePlan
	if err := json.NewDecoder(rec.Body).Decode(&sp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sp.Metadata.Name != "small" || sp.Spec.ServiceClassName != "postgres" {
		t.Errorf("got %s/%s, want postgres/small", sp.Spec.ServiceClassName, sp.Metadata.Name)
	}
}

func TestServicePlanHandler_Get_NotFound(t *testing.T) {
	h, _, _ := newTestServicePlanHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-plans/postgres/missing", nil, "")
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

func TestServicePlanHandler_Get_InvalidServiceClassNameSegment(t *testing.T) {
	h, _, _ := newTestServicePlanHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-plans/INVALID/small", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "spec.serviceClassName" {
		t.Errorf("field = %q, want spec.serviceClassName", errBody.Field)
	}
}

func TestServicePlanHandler_Get_InvalidNameSegment(t *testing.T) {
	h, _, _ := newTestServicePlanHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-plans/postgres/INVALID", nil, "")
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

func TestServicePlanHandler_Get_WrongPathShape(t *testing.T) {
	h, _, _ := newTestServicePlanHandler()
	for _, path := range []string{
		"/v1/service-plans/",
		"/v1/service-plans/postgres",
		"/v1/service-plans/postgres/small/extra",
	} {
		req := jsonRequest(http.MethodGet, path, nil, "")
		rec := httptest.NewRecorder()
		h.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("path %s: status = %d, want 404", path, rec.Code)
		}
	}
}

func TestServicePlanHandler_List_Empty(t *testing.T) {
	h, _, _ := newTestServicePlanHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-plans", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"items":[]}` {
		t.Errorf("body = %q, want {\"items\":[]}", got)
	}
}

func TestServicePlanHandler_List_Sorted(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	seedServiceClass(t, scReg, "redis")
	seedServiceClass(t, scReg, "alpha")
	inputs := []struct{ class, name string }{
		{"redis", "small"},
		{"postgres", "large"},
		{"postgres", "small"},
		{"alpha", "dev"},
	}
	for _, in := range inputs {
		if rec := createServicePlan(h, in.class, in.name); rec.Code != http.StatusCreated {
			t.Fatalf("create %s/%s status = %d", in.class, in.name, rec.Code)
		}
	}

	req := jsonRequest(http.MethodGet, "/v1/service-plans", nil, "")
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
	var items []resources.ServicePlan
	if err := json.Unmarshal(top["items"], &items); err != nil {
		t.Fatalf("decode items: %v", err)
	}
	if len(items) != 4 {
		t.Fatalf("items = %d, want 4", len(items))
	}
	for i := 1; i < len(items); i++ {
		prev, cur := items[i-1], items[i]
		if prev.Spec.ServiceClassName > cur.Spec.ServiceClassName {
			t.Fatalf("not sorted by serviceClassName: %+v", items)
		}
		if prev.Spec.ServiceClassName == cur.Spec.ServiceClassName &&
			prev.Metadata.Name >= cur.Metadata.Name {
			t.Fatalf("not sorted by name within class: %+v", items)
		}
	}
}

func TestServicePlanHandler_Update_Valid(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := map[string]any{
		"metadata": map[string]any{"name": "small"},
		"spec": map[string]any{
			"serviceClassName": "postgres",
			"displayName":      "Small Plan",
			"description":      "updated",
			"tier":             resources.TierLarge,
			"lifecycle":        resources.LifecycleDeprecated,
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/small", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var sp resources.ServicePlan
	if err := json.NewDecoder(rec.Body).Decode(&sp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sp.Spec.Description != "updated" || sp.Spec.Tier != resources.TierLarge {
		t.Errorf("mutable fields not updated: %+v", sp.Spec)
	}
}

func TestServicePlanHandler_Update_NotFound(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/missing", validServicePlanBody("postgres", "missing"), "application/json")
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

func TestServicePlanHandler_Update_NameAbsent(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := map[string]any{
		"metadata": map[string]any{},
		"spec": map[string]any{
			"serviceClassName": "postgres",
			"tier":             resources.TierSmall,
			"lifecycle":        resources.LifecycleActive,
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/small", body, "application/json")
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

func TestServicePlanHandler_Update_NameMismatch(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := validServicePlanBody("postgres", "other")
	req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/small", body, "application/json")
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

func TestServicePlanHandler_Update_ServiceClassNameAbsent(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := map[string]any{
		"metadata": map[string]any{"name": "small"},
		"spec": map[string]any{
			"tier":      resources.TierSmall,
			"lifecycle": resources.LifecycleActive,
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/small", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "spec.serviceClassName" {
		t.Errorf("field = %q, want spec.serviceClassName", errBody.Field)
	}
}

func TestServicePlanHandler_Update_ServiceClassNameMismatch(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	seedServiceClass(t, scReg, "redis")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := validServicePlanBody("redis", "small")
	req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/small", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "spec.serviceClassName" {
		t.Errorf("field = %q, want spec.serviceClassName", errBody.Field)
	}
}

func TestServicePlanHandler_Update_InvalidFields(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := map[string]any{
		"metadata": map[string]any{"name": "small"},
		"spec": map[string]any{
			"serviceClassName": "postgres",
			"tier":             "NotATier",
			"lifecycle":        resources.LifecycleActive,
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/small", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.tier" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestServicePlanHandler_Update_MissingParent(t *testing.T) {
	scReg := registry.NewServiceClassRegistry()
	spReg := registry.NewServicePlanRegistry()
	h := NewServicePlanHandler(spReg, scReg, nil)
	seedServiceClass(t, scReg, "postgres")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	if err := scReg.DeleteServiceClass(context.Background(), "postgres"); err != nil {
		t.Fatalf("delete parent: %v", err)
	}

	req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/small", validServicePlanBody("postgres", "small"), "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "spec.serviceClassName" {
		t.Errorf("field = %q, want spec.serviceClassName", errBody.Field)
	}
}

func TestServicePlanHandler_Update_PreservesServerOwnedFields(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	payload := `{
		"apiVersion":"tampered/v0",
		"kind":"Tampered",
		"metadata":{"name":"small"},
		"spec":{"serviceClassName":"postgres","tier":"Small","lifecycle":"Active","description":"changed"}
	}`
	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/service-plans/postgres/small", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var sp resources.ServicePlan
	if err := json.NewDecoder(rec.Body).Decode(&sp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sp.APIVersion != resources.ServicePlanAPIVersion {
		t.Errorf("APIVersion = %q, want preserved", sp.APIVersion)
	}
	if sp.Kind != resources.ServicePlanKind {
		t.Errorf("Kind = %q, want ServicePlan", sp.Kind)
	}
	if sp.Status.Phase != resources.PhaseActive {
		t.Errorf("Status.Phase = %q, want Active", sp.Status.Phase)
	}
	if sp.Metadata.Name != "small" || sp.Spec.ServiceClassName != "postgres" {
		t.Errorf("identity = %s/%s, want postgres/small", sp.Spec.ServiceClassName, sp.Metadata.Name)
	}
	if sp.Spec.Description != "changed" {
		t.Errorf("Description = %q, want changed", sp.Spec.Description)
	}
}

func TestServicePlanHandler_Delete_Success(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	req := jsonRequest(http.MethodDelete, "/v1/service-plans/postgres/small", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

func TestServicePlanHandler_Delete_NotFound(t *testing.T) {
	h, _, _ := newTestServicePlanHandler()
	req := jsonRequest(http.MethodDelete, "/v1/service-plans/postgres/missing", nil, "")
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

func TestServicePlanHandler_Delete_InvalidServiceClassNameSegment(t *testing.T) {
	h, _, _ := newTestServicePlanHandler()
	req := jsonRequest(http.MethodDelete, "/v1/service-plans/INVALID/small", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "spec.serviceClassName" {
		t.Errorf("field = %q, want spec.serviceClassName", errBody.Field)
	}
}

func TestServicePlanHandler_Delete_InvalidNameSegment(t *testing.T) {
	h, _, _ := newTestServicePlanHandler()
	req := jsonRequest(http.MethodDelete, "/v1/service-plans/postgres/INVALID", nil, "")
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

func TestServicePlanHandler_Delete_WrongPathShape(t *testing.T) {
	h, _, _ := newTestServicePlanHandler()
	req := jsonRequest(http.MethodDelete, "/v1/service-plans/postgres/small/extra", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestServicePlanHandler_NilEmitterNoPanic(t *testing.T) {
	h, scReg, _ := newTestServicePlanHandler()
	seedServiceClass(t, scReg, "postgres")
	if rec := createServicePlan(h, "postgres", "small"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/small", validServicePlanBody("postgres", "small"), "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200", rec.Code)
	}
	del := jsonRequest(http.MethodDelete, "/v1/service-plans/postgres/small", nil, "")
	delRec := httptest.NewRecorder()
	h.HandleItem(delRec, del)
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204", delRec.Code)
	}
}

func TestServicePlanHandler_NilLookupReturns500(t *testing.T) {
	spReg := registry.NewServicePlanRegistry()
	h := NewServicePlanHandler(spReg, nil, nil)

	rec := createServicePlan(h, "postgres", "small")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("create status = %d, want 500; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeInternalError {
		t.Errorf("create code = %q, want INTERNAL_ERROR", errBody.Code)
	}

	// Seed a plan directly so update path can be exercised without create.
	plan := resources.ServicePlan{
		APIVersion: resources.ServicePlanAPIVersion,
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
		t.Fatalf("seed plan: %v", err)
	}
	req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/small", validServicePlanBody("postgres", "small"), "application/json")
	upd := httptest.NewRecorder()
	h.HandleItem(upd, req)
	if upd.Code != http.StatusInternalServerError {
		t.Fatalf("update status = %d, want 500; body=%s", upd.Code, upd.Body.String())
	}
	updErr := decodeAPIError(t, upd)
	if updErr.Code != resources.ErrCodeInternalError {
		t.Errorf("update code = %q, want INTERNAL_ERROR", updErr.Code)
	}
}

func TestServicePlanHandler_UnsupportedMethods(t *testing.T) {
	h, _, _ := newTestServicePlanHandler()
	req := jsonRequest(http.MethodPatch, "/v1/service-plans", nil, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("collection status = %d, want 405", rec.Code)
	}

	itemReq := jsonRequest(http.MethodPatch, "/v1/service-plans/postgres/small", nil, "application/json")
	itemRec := httptest.NewRecorder()
	h.HandleItem(itemRec, itemReq)
	if itemRec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("item status = %d, want 405", itemRec.Code)
	}
}
