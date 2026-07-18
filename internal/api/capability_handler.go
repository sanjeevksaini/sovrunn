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

// CapabilityHandler holds dependencies for Capability create, read, list, and
// delete endpoints. Capability is immutable, so update is not supported.
// emitter is optional and nil-safe.
type CapabilityHandler struct {
	registry           registry.CapabilityRegistryIface
	pluginLookup       registry.PluginLookup
	serviceClassLookup registry.ServiceClassLookup
	emitter            OperationEmitter
}

// NewCapabilityHandler constructs a CapabilityHandler. emitter may be nil.
func NewCapabilityHandler(
	reg registry.CapabilityRegistryIface,
	pluginLookup registry.PluginLookup,
	serviceClassLookup registry.ServiceClassLookup,
	emitter OperationEmitter,
) *CapabilityHandler {
	return &CapabilityHandler{
		registry:           reg,
		pluginLookup:       pluginLookup,
		serviceClassLookup: serviceClassLookup,
		emitter:            emitter,
	}
}

// HandleCollection dispatches POST to Create and GET to List.
func (h *CapabilityHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Create(w, r)
	case http.MethodGet:
		h.List(w, r)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// HandleItem dispatches GET to Get and DELETE to Delete using Go
// 1.21-compatible path parsing. Capability update is explicitly unsupported.
func (h *CapabilityHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	remainder := strings.TrimPrefix(r.URL.Path, "/v1/capabilities/")
	parts := strings.Split(remainder, "/")
	if len(parts) != 1 || parts[0] == "" {
		writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "capability not found", "", "")
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
			"Capability does not support update; delete and recreate instead",
			"",
			"",
		)
	case http.MethodDelete:
		h.Delete(w, r, name)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// Create handles POST /v1/capabilities.
func (h *CapabilityHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	capability, err := safeDecodeCapability(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if errs := validation.ValidateCapability(capability); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if h.pluginLookup == nil || h.serviceClassLookup == nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	plugin, err := h.pluginLookup.GetPlugin(ctx, capability.Spec.PluginRef)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(
				w,
				r,
				http.StatusBadRequest,
				resources.ErrCodeValidationFailed,
				"referenced Plugin not found: "+capability.Spec.PluginRef,
				"spec.pluginRef",
				"",
			)
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	if _, err := h.serviceClassLookup.GetServiceClass(ctx, capability.Spec.ServiceClassRef); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(
				w,
				r,
				http.StatusBadRequest,
				resources.ErrCodeValidationFailed,
				"referenced ServiceClass not found: "+capability.Spec.ServiceClassRef,
				"spec.serviceClassRef",
				"",
			)
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	declared := false
	for _, ref := range plugin.Spec.ServiceClassRefs {
		if ref == capability.Spec.ServiceClassRef {
			declared = true
			break
		}
	}
	if !declared {
		writeError(
			w,
			r,
			http.StatusBadRequest,
			resources.ErrCodeValidationFailed,
			"ServiceClass "+capability.Spec.ServiceClassRef+" is not declared by Plugin "+capability.Spec.PluginRef,
			"spec.serviceClassRef",
			"",
		)
		return
	}

	capability.APIVersion = resources.CapabilityAPIVersion
	capability.Kind = resources.CapabilityKind
	capability.Status.Phase = resources.PhaseActive
	capability.Status.Message = ""

	created, err := h.registry.CreateCapability(ctx, capability)
	if err != nil {
		if errors.Is(err, registry.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, resources.ErrCodeResourceAlreadyExists, "capability already exists", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:           resources.OpCreateCapability,
		ResourceKind:   resources.CapabilityKind,
		ResourceName:   created.Metadata.Name,
		PluginName:     created.Spec.PluginRef,
		CapabilityName: created.Metadata.Name,
		RequestID:      requestIDFromContext(ctx),
	})
	writeJSON(w, r, http.StatusCreated, created)
}

// Get handles GET /v1/capabilities/{name}.
func (h *CapabilityHandler) Get(w http.ResponseWriter, r *http.Request, name string) {
	if errs := validation.ValidateCapabilityPathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	capability, err := h.registry.GetCapability(r.Context(), name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "capability not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	writeJSON(w, r, http.StatusOK, capability)
}

type capabilityListResponse struct {
	Items []resources.Capability `json:"items"`
}

// List handles GET /v1/capabilities with optional pluginRef and
// serviceClassRef filters.
func (h *CapabilityHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	items, err := h.registry.ListCapabilities(
		r.Context(),
		query.Get("pluginRef"),
		query.Get("serviceClassRef"),
	)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if items == nil {
		items = []resources.Capability{}
	}
	writeJSON(w, r, http.StatusOK, capabilityListResponse{Items: items})
}

// Delete handles DELETE /v1/capabilities/{name}.
func (h *CapabilityHandler) Delete(w http.ResponseWriter, r *http.Request, name string) {
	ctx := r.Context()

	if errs := validation.ValidateCapabilityPathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	capability, err := h.registry.GetCapability(ctx, name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "capability not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	if err := h.registry.DeleteCapability(ctx, name); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "capability not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:           resources.OpDeleteCapability,
		ResourceKind:   resources.CapabilityKind,
		ResourceName:   name,
		PluginName:     capability.Spec.PluginRef,
		CapabilityName: name,
		RequestID:      requestIDFromContext(ctx),
	})
	w.Header().Set("Content-Type", "application/json")
	if reqID := requestctx.RequestIDFromContext(ctx); reqID != "" {
		w.Header().Set("X-Sovrunn-Request-ID", reqID)
	}
	w.WriteHeader(http.StatusNoContent)
}
