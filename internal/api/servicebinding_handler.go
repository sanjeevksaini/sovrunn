package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/requestctx"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
	"github.com/sanjeevksaini/sovrunn/internal/validation"
)

// Phase 1 stub secret reference stored on ServiceBinding status. No real
// credentials are generated.
const serviceBindingStubSecretRef = "stub-secret-ref"

// ServiceBindingHandler holds dependencies for ServiceBinding create, read,
// list, and delete endpoints. ServiceBinding does not support update (PUT→405).
// emitter is optional and nil-safe.
type ServiceBindingHandler struct {
	registry       registry.ServiceBindingRegistryIface
	instanceLookup registry.ServiceInstanceLookup
	emitter        OperationEmitter
}

// NewServiceBindingHandler constructs a ServiceBindingHandler. emitter may be nil.
func NewServiceBindingHandler(
	reg registry.ServiceBindingRegistryIface,
	instanceLookup registry.ServiceInstanceLookup,
	emitter OperationEmitter,
) *ServiceBindingHandler {
	return &ServiceBindingHandler{
		registry:       reg,
		instanceLookup: instanceLookup,
		emitter:        emitter,
	}
}

// HandleCollection dispatches POST → Create and GET → List on
// /v1/service-bindings.
func (h *ServiceBindingHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Create(w, r)
	case http.MethodGet:
		h.List(w, r)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// HandleItem dispatches GET → Get, PUT → 405, DELETE → Delete using
// Go 1.21-compatible path parsing. The identity is a single path segment:
// "{name}".
func (h *ServiceBindingHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	remainder := strings.TrimPrefix(r.URL.Path, "/v1/service-bindings/")
	parts := strings.Split(remainder, "/")
	if len(parts) != 1 || parts[0] == "" {
		writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service binding not found", "", "")
		return
	}
	name := parts[0]

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r, name)
	case http.MethodPut:
		writeError(
			w,
			r,
			http.StatusMethodNotAllowed,
			resources.ErrCodeMethodNotAllowed,
			"ServiceBinding does not support update; delete and recreate instead",
			"",
			"",
		)
	case http.MethodDelete:
		h.Delete(w, r, name)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// Create handles POST /v1/service-bindings.
func (h *ServiceBindingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sb, err := safeDecodeServiceBinding(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if errs := validation.ValidateServiceBinding(sb); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if _, err := h.instanceLookup.GetServiceInstance(ctx, sb.Spec.ServiceInstanceRef); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(
				w,
				r,
				http.StatusBadRequest,
				resources.ErrCodeValidationFailed,
				"referenced ServiceInstance not found: "+sb.Spec.ServiceInstanceRef,
				"spec.serviceInstanceRef",
				"",
			)
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	sb.APIVersion = resources.ServiceBindingAPIVersion
	sb.Kind = resources.ServiceBindingKind
	sb.Status.Phase = "Ready"
	sb.Status.SecretRef = serviceBindingStubSecretRef

	created, err := h.registry.CreateServiceBinding(ctx, sb)
	if err != nil {
		if errors.Is(err, registry.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, resources.ErrCodeResourceAlreadyExists, "service binding already exists", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:                resources.OpCreateServiceBinding,
		ResourceKind:        resources.ServiceBindingKind,
		ResourceName:        created.Metadata.Name,
		ServiceInstanceName: created.Spec.ServiceInstanceRef,
		ServiceBindingName:  created.Metadata.Name,
		RequestID:           requestIDFromContext(ctx),
	})
	writeJSON(w, r, http.StatusCreated, created)
}

// Get handles GET /v1/service-bindings/{name}.
func (h *ServiceBindingHandler) Get(w http.ResponseWriter, r *http.Request, name string) {
	if errs := validation.ValidateServiceBindingPathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	sb, err := h.registry.GetServiceBinding(r.Context(), name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service binding not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	writeJSON(w, r, http.StatusOK, sb)
}

// serviceBindingListResponse is the list endpoint response shape.
type serviceBindingListResponse struct {
	Items []resources.ServiceBinding `json:"items"`
}

// List handles GET /v1/service-bindings with an optional serviceInstanceRef
// query filter.
func (h *ServiceBindingHandler) List(w http.ResponseWriter, r *http.Request) {
	serviceInstanceRef := r.URL.Query().Get("serviceInstanceRef")

	items, err := h.registry.ListServiceBindings(r.Context(), serviceInstanceRef)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if items == nil {
		items = []resources.ServiceBinding{}
	}
	writeJSON(w, r, http.StatusOK, serviceBindingListResponse{Items: items})
}

// Delete handles DELETE /v1/service-bindings/{name}.
func (h *ServiceBindingHandler) Delete(w http.ResponseWriter, r *http.Request, name string) {
	ctx := r.Context()

	if errs := validation.ValidateServiceBindingPathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	stored, err := h.registry.GetServiceBinding(ctx, name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service binding not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	if err := h.registry.DeleteServiceBinding(ctx, name); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service binding not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:                resources.OpDeleteServiceBinding,
		ResourceKind:        resources.ServiceBindingKind,
		ResourceName:        name,
		ServiceInstanceName: stored.Spec.ServiceInstanceRef,
		ServiceBindingName:  name,
		RequestID:           requestIDFromContext(ctx),
	})
	w.Header().Set("Content-Type", "application/json")
	if reqID := requestctx.RequestIDFromContext(ctx); reqID != "" {
		w.Header().Set("X-Sovrunn-Request-ID", reqID)
	}
	w.WriteHeader(http.StatusNoContent)
}
