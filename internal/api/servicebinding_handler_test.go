package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

type sbTestEnv struct {
	handler    *ServiceBindingHandler
	bindingReg *registry.ServiceBindingRegistry
	siReg      *registry.ServiceInstanceRegistry
}

func newTestServiceBindingHandler() *sbTestEnv {
	bindingReg := registry.NewServiceBindingRegistry()
	siReg := registry.NewServiceInstanceRegistry()
	h := NewServiceBindingHandler(bindingReg, siReg, nil)
	return &sbTestEnv{
		handler:    h,
		bindingReg: bindingReg,
		siReg:      siReg,
	}
}

func (e *sbTestEnv) seedServiceInstance(t *testing.T, name string) {
	t.Helper()
	si := resources.ServiceInstance{
		APIVersion: resources.ServiceInstanceAPIVersion,
		Kind:       resources.ServiceInstanceKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.ServiceInstanceSpec{
			OrganizationRef: "nic",
			TenantRef:       "payments",
			ProjectRef:      "prod",
			ServiceClassRef: "postgres",
			ServicePlanRef:  "small",
		},
		Status: resources.ServiceInstanceStatus{Phase: "Ready"},
	}
	if _, err := e.siReg.CreateServiceInstance(context.Background(), si); err != nil {
		t.Fatalf("seedServiceInstance(%s): %v", name, err)
	}
}

func validServiceBindingBody(name, serviceInstanceRef string) map[string]any {
	return map[string]any{
		"metadata": map[string]any{"name": name},
		"spec": map[string]any{
			"serviceInstanceRef": serviceInstanceRef,
			"consumerRef": map[string]any{
				"kind": "Application",
				"name": "payments-app",
			},
			"bindingType": resources.BindingTypeCredentials,
		},
	}
}

func createServiceBinding(h *ServiceBindingHandler, name, serviceInstanceRef string) *httptest.ResponseRecorder {
	req := jsonRequest(
		http.MethodPost,
		"/v1/service-bindings",
		validServiceBindingBody(name, serviceInstanceRef),
		"application/json",
	)
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	return rec
}

func decodeServiceBindingList(t *testing.T, rec *httptest.ResponseRecorder) []resources.ServiceBinding {
	t.Helper()
	var response serviceBindingListResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode service binding list: %v", err)
	}
	return response.Items
}

func serviceBindingNames(items []resources.ServiceBinding) []string {
	names := make([]string, len(items))
	for i := range items {
		names[i] = items[i].Metadata.Name
	}
	return names
}

func TestServiceBindingHandler_Create_Valid(t *testing.T) {
	env := newTestServiceBindingHandler()
	env.seedServiceInstance(t, "pg-prod")

	rec := createServiceBinding(env.handler, "pg-binding", "pg-prod")
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}

	var sb resources.ServiceBinding
	if err := json.NewDecoder(rec.Body).Decode(&sb); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sb.APIVersion != resources.ServiceBindingAPIVersion || sb.Kind != resources.ServiceBindingKind {
		t.Errorf("apiVersion/kind = %q/%q, want server-owned values", sb.APIVersion, sb.Kind)
	}
	if sb.Metadata.Name != "pg-binding" || sb.Spec.ServiceInstanceRef != "pg-prod" {
		t.Errorf("unexpected binding: %+v", sb)
	}
	if sb.Status.Phase != "Ready" {
		t.Errorf("phase = %q, want Ready", sb.Status.Phase)
	}
	if sb.Status.SecretRef != "stub-secret-ref" {
		t.Errorf("secretRef = %q, want stub-secret-ref", sb.Status.SecretRef)
	}
}

func TestServiceBindingHandler_Create_Duplicate(t *testing.T) {
	env := newTestServiceBindingHandler()
	env.seedServiceInstance(t, "pg-prod")
	if rec := createServiceBinding(env.handler, "pg-binding", "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d", rec.Code)
	}

	rec := createServiceBinding(env.handler, "pg-binding", "pg-prod")
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceAlreadyExists {
		t.Errorf("code = %q, want RESOURCE_ALREADY_EXISTS", errBody.Code)
	}
}

func TestServiceBindingHandler_Create_DuplicateAcrossInstances(t *testing.T) {
	env := newTestServiceBindingHandler()
	env.seedServiceInstance(t, "pg-prod")
	env.seedServiceInstance(t, "pg-dev")
	if rec := createServiceBinding(env.handler, "pg-binding", "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d", rec.Code)
	}

	rec := createServiceBinding(env.handler, "pg-binding", "pg-dev")
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceAlreadyExists {
		t.Errorf("code = %q, want RESOURCE_ALREADY_EXISTS", errBody.Code)
	}
}

func TestServiceBindingHandler_Create_InvalidFields(t *testing.T) {
	env := newTestServiceBindingHandler()
	env.seedServiceInstance(t, "pg-prod")

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
					"serviceInstanceRef": "pg-prod",
					"consumerRef":        map[string]any{"kind": "Application", "name": "payments-app"},
					"bindingType":        resources.BindingTypeCredentials,
				},
			},
			field: "metadata.name",
		},
		{
			name: "missing serviceInstanceRef",
			body: map[string]any{
				"metadata": map[string]any{"name": "pg-binding"},
				"spec": map[string]any{
					"consumerRef": map[string]any{"kind": "Application", "name": "payments-app"},
					"bindingType": resources.BindingTypeCredentials,
				},
			},
			field: "spec.serviceInstanceRef",
		},
		{
			name: "nil consumerRef",
			body: map[string]any{
				"metadata": map[string]any{"name": "pg-binding"},
				"spec": map[string]any{
					"serviceInstanceRef": "pg-prod",
					"bindingType":        resources.BindingTypeCredentials,
				},
			},
			field: "spec.consumerRef",
		},
		{
			name: "empty consumerRef.kind",
			body: map[string]any{
				"metadata": map[string]any{"name": "pg-binding"},
				"spec": map[string]any{
					"serviceInstanceRef": "pg-prod",
					"consumerRef":        map[string]any{"kind": "", "name": "payments-app"},
					"bindingType":        resources.BindingTypeCredentials,
				},
			},
			field: "spec.consumerRef.kind",
		},
		{
			name: "invalid consumerRef.name",
			body: map[string]any{
				"metadata": map[string]any{"name": "pg-binding"},
				"spec": map[string]any{
					"serviceInstanceRef": "pg-prod",
					"consumerRef":        map[string]any{"kind": "Application", "name": "INVALID"},
					"bindingType":        resources.BindingTypeCredentials,
				},
			},
			field: "spec.consumerRef.name",
		},
		{
			name: "invalid bindingType",
			body: map[string]any{
				"metadata": map[string]any{"name": "pg-binding"},
				"spec": map[string]any{
					"serviceInstanceRef": "pg-prod",
					"consumerRef":        map[string]any{"kind": "Application", "name": "payments-app"},
					"bindingType":        "endpoint",
				},
			},
			field: "spec.bindingType",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := jsonRequest(http.MethodPost, "/v1/service-bindings", tt.body, "application/json")
			rec := httptest.NewRecorder()
			env.handler.HandleCollection(rec, req)
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

func TestServiceBindingHandler_Create_StatusFieldRejected(t *testing.T) {
	env := newTestServiceBindingHandler()
	payload := `{"metadata":{"name":"pg-binding"},"spec":{"serviceInstanceRef":"pg-prod","consumerRef":{"kind":"Application","name":"payments-app"},"bindingType":"credentials"},"status":{}}`
	req := withRequestID(
		httptest.NewRequest(http.MethodPost, "/v1/service-bindings", strings.NewReader(payload)),
		"id",
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "status" {
		t.Errorf("error = %+v, want VALIDATION_FAILED field status", errBody)
	}
}

func TestServiceBindingHandler_Create_BadJSON(t *testing.T) {
	env := newTestServiceBindingHandler()
	req := withRequestID(
		httptest.NewRequest(http.MethodPost, "/v1/service-bindings", strings.NewReader("{")),
		"id",
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	if errBody := decodeAPIError(t, rec); errBody.Code != resources.ErrCodeValidationFailed {
		t.Errorf("code = %q, want VALIDATION_FAILED", errBody.Code)
	}
}

func TestServiceBindingHandler_Create_UnknownField(t *testing.T) {
	env := newTestServiceBindingHandler()
	payload := `{"metadata":{"name":"pg-binding"},"spec":{"serviceInstanceRef":"pg-prod","consumerRef":{"kind":"Application","name":"payments-app"},"bindingType":"credentials"},"bogus":true}`
	req := withRequestID(
		httptest.NewRequest(http.MethodPost, "/v1/service-bindings", strings.NewReader(payload)),
		"id",
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	if errBody := decodeAPIError(t, rec); errBody.Code != resources.ErrCodeValidationFailed {
		t.Errorf("code = %q, want VALIDATION_FAILED", errBody.Code)
	}
}

func TestServiceBindingHandler_Create_MissingServiceInstanceRef(t *testing.T) {
	env := newTestServiceBindingHandler()

	rec := createServiceBinding(env.handler, "pg-binding", "missing-instance")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.serviceInstanceRef" {
		t.Errorf("error = %+v, want VALIDATION_FAILED field spec.serviceInstanceRef", errBody)
	}
}

func TestServiceBindingHandler_Get_Exists(t *testing.T) {
	env := newTestServiceBindingHandler()
	env.seedServiceInstance(t, "pg-prod")
	if rec := createServiceBinding(env.handler, "pg-binding", "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	req := jsonRequest(http.MethodGet, "/v1/service-bindings/pg-binding", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var sb resources.ServiceBinding
	if err := json.NewDecoder(rec.Body).Decode(&sb); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sb.Metadata.Name != "pg-binding" {
		t.Errorf("name = %q, want pg-binding", sb.Metadata.Name)
	}
}

func TestServiceBindingHandler_Get_NotFound(t *testing.T) {
	env := newTestServiceBindingHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-bindings/missing", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if errBody := decodeAPIError(t, rec); errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}

func TestServiceBindingHandler_Get_InvalidPathSegment(t *testing.T) {
	env := newTestServiceBindingHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-bindings/INVALID", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "metadata.name" {
		t.Errorf("error = %+v", errBody)
	}
}

func TestServiceBindingHandler_Get_WrongPathShape(t *testing.T) {
	env := newTestServiceBindingHandler()
	for _, path := range []string{
		"/v1/service-bindings/",
		"/v1/service-bindings/pg-binding/extra",
	} {
		t.Run(path, func(t *testing.T) {
			req := jsonRequest(http.MethodGet, path, nil, "")
			rec := httptest.NewRecorder()
			env.handler.HandleItem(rec, req)
			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want 404", rec.Code)
			}
		})
	}
}

func TestServiceBindingHandler_List_Empty(t *testing.T) {
	env := newTestServiceBindingHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-bindings", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"items":[]}` {
		t.Errorf("body = %q, want {\"items\":[]}", got)
	}
}

func TestServiceBindingHandler_List_Filtered(t *testing.T) {
	env := newTestServiceBindingHandler()
	env.seedServiceInstance(t, "pg-prod")
	env.seedServiceInstance(t, "pg-dev")

	for _, input := range []struct {
		name, instanceRef string
	}{
		{"zeta", "pg-dev"},
		{"alpha", "pg-prod"},
		{"middle", "pg-prod"},
	} {
		rec := createServiceBinding(env.handler, input.name, input.instanceRef)
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
			name: "no filter",
			url:  "/v1/service-bindings",
			want: []string{"alpha", "middle", "zeta"},
		},
		{
			name: "serviceInstanceRef filter",
			url:  "/v1/service-bindings?serviceInstanceRef=pg-prod",
			want: []string{"alpha", "middle"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := jsonRequest(http.MethodGet, tt.url, nil, "")
			rec := httptest.NewRecorder()
			env.handler.HandleCollection(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
			}
			got := serviceBindingNames(decodeServiceBindingList(t, rec))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("names = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceBindingHandler_Put_MethodNotAllowedRegardlessOfExistence(t *testing.T) {
	env := newTestServiceBindingHandler()
	env.seedServiceInstance(t, "pg-prod")
	if rec := createServiceBinding(env.handler, "existing", "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("seed create status = %d", rec.Code)
	}

	for _, name := range []string{"existing", "missing"} {
		t.Run(name, func(t *testing.T) {
			req := jsonRequest(
				http.MethodPut,
				"/v1/service-bindings/"+name,
				validServiceBindingBody(name, "pg-prod"),
				"application/json",
			)
			rec := httptest.NewRecorder()
			env.handler.HandleItem(rec, req)
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

func TestServiceBindingHandler_Delete_Success(t *testing.T) {
	env := newTestServiceBindingHandler()
	env.seedServiceInstance(t, "pg-prod")
	if rec := createServiceBinding(env.handler, "pg-binding", "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	req := jsonRequest(http.MethodDelete, "/v1/service-bindings/pg-binding", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

func TestServiceBindingHandler_Delete_NotFound(t *testing.T) {
	env := newTestServiceBindingHandler()
	req := jsonRequest(http.MethodDelete, "/v1/service-bindings/missing", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if errBody := decodeAPIError(t, rec); errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}
