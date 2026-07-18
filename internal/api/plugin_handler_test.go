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

func newTestPluginHandler() (*PluginHandler, *registry.ServiceClassRegistry, *registry.PluginRegistry) {
	scReg := registry.NewServiceClassRegistry()
	pluginReg := registry.NewPluginRegistry()
	h := NewPluginHandler(pluginReg, scReg, nil, nil)
	return h, scReg, pluginReg
}

func newPluginHandlerWithBlocker(
	blocker registry.PluginChildBlocker,
) (*PluginHandler, *registry.ServiceClassRegistry, *registry.PluginRegistry) {
	scReg := registry.NewServiceClassRegistry()
	pluginReg := registry.NewPluginRegistry()
	h := NewPluginHandler(pluginReg, scReg, blocker, nil)
	return h, scReg, pluginReg
}

func seedPluginServiceClass(t *testing.T, reg *registry.ServiceClassRegistry, name string) {
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
		t.Fatalf("seedPluginServiceClass(%s): %v", name, err)
	}
}

func validPluginBody(name string, serviceClassRefs ...string) map[string]any {
	if serviceClassRefs == nil {
		serviceClassRefs = []string{"postgres"}
	}
	return map[string]any{
		"metadata": map[string]any{"name": name},
		"spec": map[string]any{
			"pluginType":       resources.PluginTypeDStoreOps,
			"version":          "1.0.0",
			"serviceClassRefs": serviceClassRefs,
			"deploymentMode":   resources.DeploymentModeCompiledIn,
		},
	}
}

func createPlugin(h *PluginHandler, name string, serviceClassRefs ...string) *httptest.ResponseRecorder {
	req := jsonRequest(http.MethodPost, "/v1/plugins", validPluginBody(name, serviceClassRefs...), "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	return rec
}

func TestPluginHandler_Create_Valid(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")

	rec := createPlugin(h, "pg-plugin")
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var p resources.Plugin
	if err := json.NewDecoder(rec.Body).Decode(&p); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if p.APIVersion != resources.PluginAPIVersion || p.Kind != resources.PluginKind {
		t.Errorf("apiVersion/kind = %q/%q, want server-owned", p.APIVersion, p.Kind)
	}
	if p.Metadata.Name != "pg-plugin" {
		t.Errorf("name = %q, want pg-plugin", p.Metadata.Name)
	}
	if p.Status.Phase != resources.PhaseActive {
		t.Errorf("phase = %q, want Active", p.Status.Phase)
	}
	if p.Spec.PluginType != resources.PluginTypeDStoreOps {
		t.Errorf("pluginType = %q, want dStoreOps", p.Spec.PluginType)
	}
}

func TestPluginHandler_Create_Duplicate(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	if rec := createPlugin(h, "pg-plugin"); rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d", rec.Code)
	}
	rec := createPlugin(h, "pg-plugin")
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceAlreadyExists {
		t.Errorf("code = %q, want RESOURCE_ALREADY_EXISTS", errBody.Code)
	}
}

func TestPluginHandler_Create_InvalidFields(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")

	cases := []struct {
		name  string
		body  map[string]any
		field string
	}{
		{
			name: "missing name",
			body: map[string]any{
				"metadata": map[string]any{},
				"spec": map[string]any{
					"pluginType":       resources.PluginTypeDStoreOps,
					"version":          "1.0.0",
					"serviceClassRefs": []string{"postgres"},
					"deploymentMode":   resources.DeploymentModeCompiledIn,
				},
			},
			field: "metadata.name",
		},
		{
			name: "invalid pluginType",
			body: map[string]any{
				"metadata": map[string]any{"name": "pg-plugin"},
				"spec": map[string]any{
					"pluginType":       "notAPluginType",
					"version":          "1.0.0",
					"serviceClassRefs": []string{"postgres"},
					"deploymentMode":   resources.DeploymentModeCompiledIn,
				},
			},
			field: "spec.pluginType",
		},
		{
			name: "empty version",
			body: map[string]any{
				"metadata": map[string]any{"name": "pg-plugin"},
				"spec": map[string]any{
					"pluginType":       resources.PluginTypeDStoreOps,
					"version":          "",
					"serviceClassRefs": []string{"postgres"},
					"deploymentMode":   resources.DeploymentModeCompiledIn,
				},
			},
			field: "spec.version",
		},
		{
			name: "empty serviceClassRefs",
			body: map[string]any{
				"metadata": map[string]any{"name": "pg-plugin"},
				"spec": map[string]any{
					"pluginType":       resources.PluginTypeDStoreOps,
					"version":          "1.0.0",
					"serviceClassRefs": []string{},
					"deploymentMode":   resources.DeploymentModeCompiledIn,
				},
			},
			field: "spec.serviceClassRefs",
		},
		{
			name: "invalid deploymentMode",
			body: map[string]any{
				"metadata": map[string]any{"name": "pg-plugin"},
				"spec": map[string]any{
					"pluginType":       resources.PluginTypeDStoreOps,
					"version":          "1.0.0",
					"serviceClassRefs": []string{"postgres"},
					"deploymentMode":   "sidecar",
				},
			},
			field: "spec.deploymentMode",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := jsonRequest(http.MethodPost, "/v1/plugins", tc.body, "application/json")
			rec := httptest.NewRecorder()
			h.HandleCollection(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
			}
			errBody := decodeAPIError(t, rec)
			if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != tc.field {
				t.Errorf("error = %+v, want field %q", errBody, tc.field)
			}
		})
	}
}

func TestPluginHandler_Create_StatusFieldRejected(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	payload := `{"metadata":{"name":"pg-plugin"},"spec":{"pluginType":"dStoreOps","version":"1.0.0","serviceClassRefs":["postgres"],"deploymentMode":"compiled-in"},"status":{}}`
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/plugins", strings.NewReader(payload)), "id")
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

func TestPluginHandler_Create_BadJSON(t *testing.T) {
	h, _, _ := newTestPluginHandler()
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/plugins", strings.NewReader("{")), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestPluginHandler_Create_UnknownField(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	payload := `{"metadata":{"name":"pg-plugin"},"spec":{"pluginType":"dStoreOps","version":"1.0.0","serviceClassRefs":["postgres"],"deploymentMode":"compiled-in"},"bogus":true}`
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/plugins", strings.NewReader(payload)), "id")
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

func TestPluginHandler_Create_OversizedBody(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	large := strings.Repeat("a", 1<<20+1)
	payload := fmt.Sprintf(
		`{"metadata":{"name":"pg-plugin"},"spec":{"pluginType":"dStoreOps","version":"1.0.0","serviceClassRefs":["postgres"],"deploymentMode":"compiled-in","description":"%s"}}`,
		large,
	)
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/plugins", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

func TestPluginHandler_Create_MissingServiceClassRef(t *testing.T) {
	h, _, _ := newTestPluginHandler()
	rec := createPlugin(h, "pg-plugin", "missing-sc")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.serviceClassRefs" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestPluginHandler_Get_Exists(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	if rec := createPlugin(h, "pg-plugin"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	req := jsonRequest(http.MethodGet, "/v1/plugins/pg-plugin", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var p resources.Plugin
	if err := json.NewDecoder(rec.Body).Decode(&p); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if p.Metadata.Name != "pg-plugin" || p.APIVersion == "" || p.Kind == "" || p.Status.Phase == "" {
		t.Errorf("incomplete resource: %+v", p)
	}
}

func TestPluginHandler_Get_NotFound(t *testing.T) {
	h, _, _ := newTestPluginHandler()
	req := jsonRequest(http.MethodGet, "/v1/plugins/missing", nil, "")
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

func TestPluginHandler_Get_InvalidPathSegment(t *testing.T) {
	h, _, _ := newTestPluginHandler()
	req := jsonRequest(http.MethodGet, "/v1/plugins/INVALID", nil, "")
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

func TestPluginHandler_Get_WrongPathShape(t *testing.T) {
	h, _, _ := newTestPluginHandler()
	for _, path := range []string{"/v1/plugins/", "/v1/plugins/pg-plugin/extra"} {
		req := jsonRequest(http.MethodGet, path, nil, "")
		rec := httptest.NewRecorder()
		h.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("path %s: status = %d, want 404", path, rec.Code)
		}
	}
}

func TestPluginHandler_List_Empty(t *testing.T) {
	h, _, _ := newTestPluginHandler()
	req := jsonRequest(http.MethodGet, "/v1/plugins", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"items":[]}` {
		t.Errorf("body = %q, want {\"items\":[]}", got)
	}
}

func TestPluginHandler_List_Sorted(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	for _, name := range []string{"zebra", "alpha", "mongo"} {
		if rec := createPlugin(h, name); rec.Code != http.StatusCreated {
			t.Fatalf("create %s status = %d", name, rec.Code)
		}
	}
	req := jsonRequest(http.MethodGet, "/v1/plugins", nil, "")
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
	var items []resources.Plugin
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

func TestPluginHandler_Update_Valid(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	if rec := createPlugin(h, "pg-plugin"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := map[string]any{
		"metadata": map[string]any{"name": "pg-plugin"},
		"spec": map[string]any{
			"pluginType":       resources.PluginTypeCacheOps,
			"version":          "2.0.0",
			"serviceClassRefs": []string{"postgres"},
			"deploymentMode":   resources.DeploymentModeCompiledIn,
			"description":      "updated",
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/plugins/pg-plugin", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var p resources.Plugin
	if err := json.NewDecoder(rec.Body).Decode(&p); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if p.Spec.PluginType != resources.PluginTypeCacheOps || p.Spec.Version != "2.0.0" || p.Spec.Description != "updated" {
		t.Errorf("mutable fields not updated: %+v", p.Spec)
	}
}

func TestPluginHandler_Update_NotFound(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	req := jsonRequest(http.MethodPut, "/v1/plugins/missing", validPluginBody("missing"), "application/json")
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

func TestPluginHandler_Update_NameAbsent(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	if rec := createPlugin(h, "pg-plugin"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := map[string]any{
		"metadata": map[string]any{},
		"spec": map[string]any{
			"pluginType":       resources.PluginTypeDStoreOps,
			"version":          "1.0.0",
			"serviceClassRefs": []string{"postgres"},
			"deploymentMode":   resources.DeploymentModeCompiledIn,
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/plugins/pg-plugin", body, "application/json")
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

func TestPluginHandler_Update_NameMismatch(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	if rec := createPlugin(h, "pg-plugin"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := validPluginBody("other")
	req := jsonRequest(http.MethodPut, "/v1/plugins/pg-plugin", body, "application/json")
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

func TestPluginHandler_Update_InvalidFields(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	if rec := createPlugin(h, "pg-plugin"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := map[string]any{
		"metadata": map[string]any{"name": "pg-plugin"},
		"spec": map[string]any{
			"pluginType":       "notAPluginType",
			"version":          "1.0.0",
			"serviceClassRefs": []string{"postgres"},
			"deploymentMode":   resources.DeploymentModeCompiledIn,
		},
	}
	req := jsonRequest(http.MethodPut, "/v1/plugins/pg-plugin", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.pluginType" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestPluginHandler_Update_MissingServiceClassRef(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	if rec := createPlugin(h, "pg-plugin"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := validPluginBody("pg-plugin", "missing-sc")
	req := jsonRequest(http.MethodPut, "/v1/plugins/pg-plugin", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.serviceClassRefs" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestPluginHandler_Delete_Success(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	if rec := createPlugin(h, "pg-plugin"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	req := jsonRequest(http.MethodDelete, "/v1/plugins/pg-plugin", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

func TestPluginHandler_Delete_NotFound(t *testing.T) {
	h, _, _ := newTestPluginHandler()
	req := jsonRequest(http.MethodDelete, "/v1/plugins/missing", nil, "")
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

func TestPluginHandler_Delete_BlockedByCapability(t *testing.T) {
	capReg := registry.NewCapabilityRegistry()
	blocker := registry.NewCapabilityChildBlockerChecker(capReg)
	h, scReg, _ := newPluginHandlerWithBlocker(blocker)
	seedPluginServiceClass(t, scReg, "postgres")

	if rec := createPlugin(h, "pg-plugin"); rec.Code != http.StatusCreated {
		t.Fatalf("create Plugin status = %d", rec.Code)
	}
	cap := resources.Capability{
		APIVersion: resources.CapabilityAPIVersion,
		Kind:       resources.CapabilityKind,
		Metadata:   resources.Metadata{Name: "pg-provision"},
		Spec: resources.CapabilitySpec{
			PluginRef:       "pg-plugin",
			ServiceClassRef: "postgres",
			Operation:       resources.CapOpProvision,
			Supported:       true,
		},
		Status: resources.CapabilityStatus{Phase: resources.PhaseActive},
	}
	if _, err := capReg.CreateCapability(context.Background(), cap); err != nil {
		t.Fatalf("seed Capability: %v", err)
	}

	req := jsonRequest(http.MethodDelete, "/v1/plugins/pg-plugin", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeDeleteBlocked {
		t.Errorf("code = %q, want DELETE_BLOCKED", errBody.Code)
	}
	if !strings.Contains(errBody.Message, "Capability") {
		t.Errorf("message = %q, want it to mention Capability", errBody.Message)
	}
}

func TestPluginHandler_Delete_ZeroCapabilities(t *testing.T) {
	capReg := registry.NewCapabilityRegistry()
	blocker := registry.NewCapabilityChildBlockerChecker(capReg)
	h, scReg, _ := newPluginHandlerWithBlocker(blocker)
	seedPluginServiceClass(t, scReg, "postgres")

	if rec := createPlugin(h, "pg-plugin"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	req := jsonRequest(http.MethodDelete, "/v1/plugins/pg-plugin", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
}

func TestPluginHandler_NilBlockerAllows(t *testing.T) {
	h, scReg, _ := newPluginHandlerWithBlocker(nil)
	seedPluginServiceClass(t, scReg, "postgres")
	if rec := createPlugin(h, "pg-plugin"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	req := jsonRequest(http.MethodDelete, "/v1/plugins/pg-plugin", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
}

func TestPluginHandler_NilEmitterNoPanic(t *testing.T) {
	h, scReg, _ := newTestPluginHandler()
	seedPluginServiceClass(t, scReg, "postgres")
	if rec := createPlugin(h, "pg-plugin"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}
	body := validPluginBody("pg-plugin")
	req := jsonRequest(http.MethodPut, "/v1/plugins/pg-plugin", body, "application/json")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200", rec.Code)
	}
	del := jsonRequest(http.MethodDelete, "/v1/plugins/pg-plugin", nil, "")
	delRec := httptest.NewRecorder()
	h.HandleItem(delRec, del)
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204", delRec.Code)
	}
}

func TestPluginHandler_UnsupportedMethods(t *testing.T) {
	h, _, _ := newTestPluginHandler()
	req := jsonRequest(http.MethodPatch, "/v1/plugins", nil, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("collection status = %d, want 405", rec.Code)
	}

	itemReq := jsonRequest(http.MethodPatch, "/v1/plugins/pg-plugin", nil, "application/json")
	itemRec := httptest.NewRecorder()
	h.HandleItem(itemRec, itemReq)
	if itemRec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("item status = %d, want 405", itemRec.Code)
	}
}
