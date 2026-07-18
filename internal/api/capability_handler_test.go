package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func newTestCapabilityHandler() (
	*CapabilityHandler,
	*registry.ServiceClassRegistry,
	*registry.PluginRegistry,
	*registry.CapabilityRegistry,
) {
	serviceClasses := registry.NewServiceClassRegistry()
	plugins := registry.NewPluginRegistry()
	capabilities := registry.NewCapabilityRegistry()
	handler := NewCapabilityHandler(capabilities, plugins, serviceClasses, nil)
	return handler, serviceClasses, plugins, capabilities
}

func seedCapabilityPlugin(
	t *testing.T,
	reg *registry.PluginRegistry,
	name string,
	serviceClassRefs ...string,
) {
	t.Helper()
	plugin := resources.Plugin{
		APIVersion: resources.PluginAPIVersion,
		Kind:       resources.PluginKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.PluginSpec{
			PluginType:       resources.PluginTypeDStoreOps,
			Version:          "1.0.0",
			ServiceClassRefs: serviceClassRefs,
			DeploymentMode:   resources.DeploymentModeCompiledIn,
		},
		Status: resources.PluginStatus{Phase: resources.PhaseActive},
	}
	if _, err := reg.CreatePlugin(context.Background(), plugin); err != nil {
		t.Fatalf("seedCapabilityPlugin(%s): %v", name, err)
	}
}

func validCapabilityBody(name, pluginRef, serviceClassRef string) map[string]any {
	return map[string]any{
		"metadata": map[string]any{"name": name},
		"spec": map[string]any{
			"pluginRef":       pluginRef,
			"serviceClassRef": serviceClassRef,
			"operation":       resources.CapOpProvision,
			"supported":       true,
		},
	}
}

func createCapability(
	h *CapabilityHandler,
	name, pluginRef, serviceClassRef string,
) *httptest.ResponseRecorder {
	req := jsonRequest(
		http.MethodPost,
		"/v1/capabilities",
		validCapabilityBody(name, pluginRef, serviceClassRef),
		"application/json",
	)
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	return rec
}

func decodeCapabilityList(t *testing.T, rec *httptest.ResponseRecorder) []resources.Capability {
	t.Helper()
	var response capabilityListResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode capability list: %v", err)
	}
	return response.Items
}

func capabilityNames(items []resources.Capability) []string {
	names := make([]string, len(items))
	for i := range items {
		names[i] = items[i].Metadata.Name
	}
	return names
}

func TestCapabilityHandler_Create_Valid(t *testing.T) {
	h, serviceClasses, plugins, _ := newTestCapabilityHandler()
	seedPluginServiceClass(t, serviceClasses, "postgres")
	seedCapabilityPlugin(t, plugins, "pg-plugin", "postgres")

	rec := createCapability(h, "pg-provision", "pg-plugin", "postgres")
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}

	var capability resources.Capability
	if err := json.NewDecoder(rec.Body).Decode(&capability); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if capability.APIVersion != resources.CapabilityAPIVersion ||
		capability.Kind != resources.CapabilityKind {
		t.Errorf(
			"apiVersion/kind = %q/%q, want server-owned values",
			capability.APIVersion,
			capability.Kind,
		)
	}
	if capability.Metadata.Name != "pg-provision" ||
		capability.Spec.PluginRef != "pg-plugin" ||
		capability.Spec.ServiceClassRef != "postgres" {
		t.Errorf("unexpected capability: %+v", capability)
	}
	if capability.Status.Phase != resources.PhaseActive {
		t.Errorf("phase = %q, want Active", capability.Status.Phase)
	}
}

func TestCapabilityHandler_Create_Duplicate(t *testing.T) {
	h, serviceClasses, plugins, _ := newTestCapabilityHandler()
	seedPluginServiceClass(t, serviceClasses, "postgres")
	seedCapabilityPlugin(t, plugins, "pg-plugin", "postgres")
	if rec := createCapability(h, "pg-provision", "pg-plugin", "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d", rec.Code)
	}

	rec := createCapability(h, "pg-provision", "pg-plugin", "postgres")
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceAlreadyExists {
		t.Errorf("code = %q, want RESOURCE_ALREADY_EXISTS", errBody.Code)
	}
}

func TestCapabilityHandler_Create_InvalidFields(t *testing.T) {
	h, serviceClasses, plugins, _ := newTestCapabilityHandler()
	seedPluginServiceClass(t, serviceClasses, "postgres")
	seedCapabilityPlugin(t, plugins, "pg-plugin", "postgres")

	tests := []struct {
		name  string
		body  map[string]any
		field string
	}{
		{
			name: "missing name",
			body: map[string]any{
				"metadata": map[string]any{},
				"spec": map[string]any{
					"pluginRef":       "pg-plugin",
					"serviceClassRef": "postgres",
					"operation":       resources.CapOpProvision,
				},
			},
			field: "metadata.name",
		},
		{
			name: "invalid pluginRef",
			body: map[string]any{
				"metadata": map[string]any{"name": "pg-provision"},
				"spec": map[string]any{
					"pluginRef":       "INVALID",
					"serviceClassRef": "postgres",
					"operation":       resources.CapOpProvision,
				},
			},
			field: "spec.pluginRef",
		},
		{
			name: "invalid serviceClassRef",
			body: map[string]any{
				"metadata": map[string]any{"name": "pg-provision"},
				"spec": map[string]any{
					"pluginRef":       "pg-plugin",
					"serviceClassRef": "INVALID",
					"operation":       resources.CapOpProvision,
				},
			},
			field: "spec.serviceClassRef",
		},
		{
			name: "invalid operation",
			body: map[string]any{
				"metadata": map[string]any{"name": "pg-provision"},
				"spec": map[string]any{
					"pluginRef":       "pg-plugin",
					"serviceClassRef": "postgres",
					"operation":       "Execute",
				},
			},
			field: "spec.operation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := jsonRequest(http.MethodPost, "/v1/capabilities", tt.body, "application/json")
			rec := httptest.NewRecorder()
			h.HandleCollection(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
			}
			errBody := decodeAPIError(t, rec)
			if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != tt.field {
				t.Errorf("error = %+v, want VALIDATION_FAILED field %q", errBody, tt.field)
			}
		})
	}
}

func TestCapabilityHandler_Create_StatusFieldRejected(t *testing.T) {
	h, _, _, _ := newTestCapabilityHandler()
	payload := `{"metadata":{"name":"pg-provision"},"spec":{"pluginRef":"pg-plugin","serviceClassRef":"postgres","operation":"Provision","supported":true},"status":{}}`
	req := withRequestID(
		httptest.NewRequest(http.MethodPost, "/v1/capabilities", strings.NewReader(payload)),
		"id",
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "status" {
		t.Errorf("error = %+v, want VALIDATION_FAILED field status", errBody)
	}
}

func TestCapabilityHandler_Create_BadJSON(t *testing.T) {
	h, _, _, _ := newTestCapabilityHandler()
	req := withRequestID(
		httptest.NewRequest(http.MethodPost, "/v1/capabilities", strings.NewReader("{")),
		"id",
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	if errBody := decodeAPIError(t, rec); errBody.Code != resources.ErrCodeValidationFailed {
		t.Errorf("code = %q, want VALIDATION_FAILED", errBody.Code)
	}
}

func TestCapabilityHandler_Create_UnknownField(t *testing.T) {
	h, _, _, _ := newTestCapabilityHandler()
	payload := `{"metadata":{"name":"pg-provision"},"spec":{"pluginRef":"pg-plugin","serviceClassRef":"postgres","operation":"Provision","supported":true},"bogus":true}`
	req := withRequestID(
		httptest.NewRequest(http.MethodPost, "/v1/capabilities", strings.NewReader(payload)),
		"id",
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	if errBody := decodeAPIError(t, rec); errBody.Code != resources.ErrCodeValidationFailed {
		t.Errorf("code = %q, want VALIDATION_FAILED", errBody.Code)
	}
}

func TestCapabilityHandler_Create_OversizedBody(t *testing.T) {
	h, _, _, _ := newTestCapabilityHandler()
	payload := fmt.Sprintf(
		`{"metadata":{"name":"pg-provision"},"spec":{"pluginRef":"pg-plugin","serviceClassRef":"postgres","operation":"Provision","supported":true,"description":"%s"}}`,
		strings.Repeat("a", 1<<20+1),
	)
	req := withRequestID(
		httptest.NewRequest(http.MethodPost, "/v1/capabilities", strings.NewReader(payload)),
		"id",
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

func TestCapabilityHandler_Create_MissingPluginRef(t *testing.T) {
	h, serviceClasses, _, _ := newTestCapabilityHandler()
	seedPluginServiceClass(t, serviceClasses, "postgres")

	rec := createCapability(h, "pg-provision", "missing-plugin", "postgres")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.pluginRef" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestCapabilityHandler_Create_MissingServiceClassRef(t *testing.T) {
	h, _, plugins, _ := newTestCapabilityHandler()
	seedCapabilityPlugin(t, plugins, "pg-plugin", "missing-sc")

	rec := createCapability(h, "pg-provision", "pg-plugin", "missing-sc")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.serviceClassRef" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestCapabilityHandler_Create_ServiceClassNotDeclaredByPlugin(t *testing.T) {
	h, serviceClasses, plugins, _ := newTestCapabilityHandler()
	seedPluginServiceClass(t, serviceClasses, "postgres")
	seedPluginServiceClass(t, serviceClasses, "redis")
	seedCapabilityPlugin(t, plugins, "pg-plugin", "postgres")

	rec := createCapability(h, "redis-provision", "pg-plugin", "redis")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed ||
		errBody.Field != "spec.serviceClassRef" ||
		!strings.Contains(errBody.Message, "is not declared by Plugin") {
		t.Errorf("error = %+v", errBody)
	}
}

func TestCapabilityHandler_Get_Exists(t *testing.T) {
	h, serviceClasses, plugins, _ := newTestCapabilityHandler()
	seedPluginServiceClass(t, serviceClasses, "postgres")
	seedCapabilityPlugin(t, plugins, "pg-plugin", "postgres")
	if rec := createCapability(h, "pg-provision", "pg-plugin", "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	req := jsonRequest(http.MethodGet, "/v1/capabilities/pg-provision", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var capability resources.Capability
	if err := json.NewDecoder(rec.Body).Decode(&capability); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if capability.Metadata.Name != "pg-provision" {
		t.Errorf("name = %q, want pg-provision", capability.Metadata.Name)
	}
}

func TestCapabilityHandler_Get_NotFound(t *testing.T) {
	h, _, _, _ := newTestCapabilityHandler()
	req := jsonRequest(http.MethodGet, "/v1/capabilities/missing", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if errBody := decodeAPIError(t, rec); errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}

func TestCapabilityHandler_Get_InvalidPathSegment(t *testing.T) {
	h, _, _, _ := newTestCapabilityHandler()
	req := jsonRequest(http.MethodGet, "/v1/capabilities/INVALID", nil, "")
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

func TestCapabilityHandler_Get_WrongPathShape(t *testing.T) {
	h, _, _, _ := newTestCapabilityHandler()
	for _, path := range []string{
		"/v1/capabilities/",
		"/v1/capabilities/pg-provision/extra",
	} {
		t.Run(path, func(t *testing.T) {
			req := jsonRequest(http.MethodGet, path, nil, "")
			rec := httptest.NewRecorder()
			h.HandleItem(rec, req)
			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want 404", rec.Code)
			}
		})
	}
}

func TestCapabilityHandler_List_Empty(t *testing.T) {
	h, _, _, _ := newTestCapabilityHandler()
	req := jsonRequest(http.MethodGet, "/v1/capabilities", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"items":[]}` {
		t.Errorf("body = %q, want {\"items\":[]}", got)
	}
}

func TestCapabilityHandler_List_SortedAndFiltered(t *testing.T) {
	h, serviceClasses, plugins, _ := newTestCapabilityHandler()
	seedPluginServiceClass(t, serviceClasses, "postgres")
	seedPluginServiceClass(t, serviceClasses, "redis")
	seedCapabilityPlugin(t, plugins, "plugin-a", "postgres", "redis")
	seedCapabilityPlugin(t, plugins, "plugin-b", "redis")

	for _, input := range []struct {
		name, pluginRef, serviceClassRef string
	}{
		{"zeta", "plugin-a", "redis"},
		{"alpha", "plugin-a", "postgres"},
		{"middle", "plugin-b", "redis"},
	} {
		rec := createCapability(h, input.name, input.pluginRef, input.serviceClassRef)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create %s status = %d; body=%s", input.name, rec.Code, rec.Body.String())
		}
	}

	tests := []struct {
		name string
		url  string
		want []string
	}{
		{
			name: "neither filter",
			url:  "/v1/capabilities",
			want: []string{"alpha", "middle", "zeta"},
		},
		{
			name: "pluginRef only",
			url:  "/v1/capabilities?pluginRef=plugin-a",
			want: []string{"alpha", "zeta"},
		},
		{
			name: "serviceClassRef only",
			url:  "/v1/capabilities?serviceClassRef=redis",
			want: []string{"middle", "zeta"},
		},
		{
			name: "both filters",
			url:  "/v1/capabilities?pluginRef=plugin-a&serviceClassRef=redis",
			want: []string{"zeta"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := jsonRequest(http.MethodGet, tt.url, nil, "")
			rec := httptest.NewRecorder()
			h.HandleCollection(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
			}
			got := capabilityNames(decodeCapabilityList(t, rec))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("names = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCapabilityHandler_Delete_Success(t *testing.T) {
	h, serviceClasses, plugins, _ := newTestCapabilityHandler()
	seedPluginServiceClass(t, serviceClasses, "postgres")
	seedCapabilityPlugin(t, plugins, "pg-plugin", "postgres")
	if rec := createCapability(h, "pg-provision", "pg-plugin", "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	req := jsonRequest(http.MethodDelete, "/v1/capabilities/pg-provision", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

func TestCapabilityHandler_Delete_NotFound(t *testing.T) {
	h, _, _, _ := newTestCapabilityHandler()
	req := jsonRequest(http.MethodDelete, "/v1/capabilities/missing", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if errBody := decodeAPIError(t, rec); errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}

func TestCapabilityHandler_Put_MethodNotAllowedRegardlessOfExistence(t *testing.T) {
	h, _, _, capabilities := newTestCapabilityHandler()
	stored := resources.Capability{
		APIVersion: resources.CapabilityAPIVersion,
		Kind:       resources.CapabilityKind,
		Metadata:   resources.Metadata{Name: "existing"},
		Spec: resources.CapabilitySpec{
			PluginRef:       "pg-plugin",
			ServiceClassRef: "postgres",
			Operation:       resources.CapOpProvision,
		},
		Status: resources.CapabilityStatus{Phase: resources.PhaseActive},
	}
	if _, err := capabilities.CreateCapability(context.Background(), stored); err != nil {
		t.Fatalf("seed capability: %v", err)
	}

	for _, name := range []string{"existing", "missing"} {
		t.Run(name, func(t *testing.T) {
			req := jsonRequest(
				http.MethodPut,
				"/v1/capabilities/"+name,
				validCapabilityBody(name, "pg-plugin", "postgres"),
				"application/json",
			)
			rec := httptest.NewRecorder()
			h.HandleItem(rec, req)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want 405; body=%s", rec.Code, rec.Body.String())
			}
			errBody := decodeAPIError(t, rec)
			if errBody.Code != resources.ErrCodeMethodNotAllowed {
				t.Errorf("code = %q, want METHOD_NOT_ALLOWED", errBody.Code)
			}
			if !strings.Contains(errBody.Message, "does not support update") {
				t.Errorf("message = %q, want update guidance", errBody.Message)
			}
		})
	}
}

func TestCapabilityHandler_NilEmitterNoPanic(t *testing.T) {
	h, serviceClasses, plugins, _ := newTestCapabilityHandler()
	seedPluginServiceClass(t, serviceClasses, "postgres")
	seedCapabilityPlugin(t, plugins, "pg-plugin", "postgres")

	if rec := createCapability(h, "pg-provision", "pg-plugin", "postgres"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	req := jsonRequest(http.MethodDelete, "/v1/capabilities/pg-provision", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
}
