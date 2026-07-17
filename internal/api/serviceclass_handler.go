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

// ServiceClassHandler holds dependencies for ServiceClass CRUD endpoints.
// blocker and emitter are optional (nil-safe).
type ServiceClassHandler struct {
	registry registry.ServiceClassRegistryIface
	blocker  registry.ServiceClassChildBlocker
	emitter  OperationEmitter
}

// NewServiceClassHandler constructs a ServiceClassHandler. blocker and emitter
// may be nil.
func NewServiceClassHandler(
	reg registry.ServiceClassRegistryIface,
	blocker registry.ServiceClassChildBlocker,
	emitter OperationEmitter,
) *ServiceClassHandler {
	return &ServiceClassHandler{registry: reg, blocker: blocker, emitter: emitter}
}

// HandleCollection dispatches POST → Create and GET → List.
func (h *ServiceClassHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Create(w, r)
	case http.MethodGet:
		h.List(w, r)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// HandleItem dispatches GET → Get, PUT → Update, DELETE → Delete using
// Go 1.21-compatible path parsing. The identity is a single path segment:
// "{name}".
func (h *ServiceClassHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	remainder := strings.TrimPrefix(r.URL.Path, "/v1/service-classes/")
	parts := strings.Split(remainder, "/")
	if len(parts) != 1 || parts[0] == "" {
		writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service class not found", "", "")
		return
	}
	name := parts[0]

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r, name)
	case http.MethodPut:
		h.Update(w, r, name)
	case http.MethodDelete:
		h.Delete(w, r, name)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// Create handles POST /v1/service-classes.
func (h *ServiceClassHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sc, err := safeDecodeServiceClass(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if errs := validation.ValidateServiceClass(sc); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	sc.APIVersion = resources.ServiceClassAPIVersion
	sc.Kind = resources.ServiceClassKind
	sc.Status.Phase = resources.PhaseActive
	sc.Status.Message = ""

	created, err := h.registry.CreateServiceClass(ctx, sc)
	if err != nil {
		if errors.Is(err, registry.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, resources.ErrCodeResourceAlreadyExists, "service class already exists", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:             resources.OpCreateServiceClass,
		ResourceKind:     resources.ServiceClassKind,
		ResourceName:     created.Metadata.Name,
		ServiceClassName: created.Metadata.Name,
		RequestID:        requestIDFromContext(ctx),
	})
	writeJSON(w, r, http.StatusCreated, created)
}

// Get handles GET /v1/service-classes/{name}.
func (h *ServiceClassHandler) Get(w http.ResponseWriter, r *http.Request, name string) {
	if errs := validation.ValidateServiceClassPathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	sc, err := h.registry.GetServiceClass(r.Context(), name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service class not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	writeJSON(w, r, http.StatusOK, sc)
}

// serviceClassListResponse is the list endpoint response shape.
type serviceClassListResponse struct {
	Items []resources.ServiceClass `json:"items"`
}

// List handles GET /v1/service-classes.
func (h *ServiceClassHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.registry.ListServiceClasses(r.Context())
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if items == nil {
		items = []resources.ServiceClass{}
	}
	writeJSON(w, r, http.StatusOK, serviceClassListResponse{Items: items})
}

// Update handles PUT /v1/service-classes/{name}.
func (h *ServiceClassHandler) Update(w http.ResponseWriter, r *http.Request, name string) {
	ctx := r.Context()

	if errs := validation.ValidateServiceClassPathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	sc, err := safeDecodeServiceClass(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if sc.Metadata.Name == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name is required in request body", "metadata.name", "")
		return
	}
	if sc.Metadata.Name != name {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name in body must match path", "metadata.name", "")
		return
	}

	if errs := validation.ValidateServiceClass(sc); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	updated, err := h.registry.UpdateServiceClass(ctx, sc)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service class not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:             resources.OpUpdateServiceClass,
		ResourceKind:     resources.ServiceClassKind,
		ResourceName:     updated.Metadata.Name,
		ServiceClassName: updated.Metadata.Name,
		RequestID:        requestIDFromContext(ctx),
	})
	writeJSON(w, r, http.StatusOK, updated)
}

// Delete handles DELETE /v1/service-classes/{name}.
func (h *ServiceClassHandler) Delete(w http.ResponseWriter, r *http.Request, name string) {
	ctx := r.Context()

	if errs := validation.ValidateServiceClassPathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if h.blocker != nil {
		blockers, err := h.blocker.BlockedByServiceClassChildren(ctx, name)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
			return
		}
		if len(blockers) > 0 {
			msg := "deletion blocked by " + blockers[0].Kind + " resources"
			writeError(w, r, http.StatusConflict, resources.ErrCodeDeleteBlocked, msg, "", "")
			return
		}
	}

	if err := h.registry.DeleteServiceClass(ctx, name); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service class not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:             resources.OpDeleteServiceClass,
		ResourceKind:     resources.ServiceClassKind,
		ResourceName:     name,
		ServiceClassName: name,
		RequestID:        requestIDFromContext(ctx),
	})
	w.Header().Set("Content-Type", "application/json")
	if reqID := requestctx.RequestIDFromContext(ctx); reqID != "" {
		w.Header().Set("X-Sovrunn-Request-ID", reqID)
	}
	w.WriteHeader(http.StatusNoContent)
}
