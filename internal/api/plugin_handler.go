package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/requestctx"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
	"github.com/sanjeevksaini/sovrunn/internal/validation"
)

// PluginHandler holds dependencies for Plugin CRUD endpoints.
// blocker and emitter are optional (nil-safe). serviceClassLookup verifies
// each entry in spec.serviceClassRefs on create and update.
type PluginHandler struct {
	registry           registry.PluginRegistryIface
	serviceClassLookup registry.ServiceClassLookup
	blocker            registry.PluginChildBlocker
	emitter            OperationEmitter
}

// NewPluginHandler constructs a PluginHandler. blocker and emitter may be nil.
func NewPluginHandler(
	reg registry.PluginRegistryIface,
	serviceClassLookup registry.ServiceClassLookup,
	blocker registry.PluginChildBlocker,
	emitter OperationEmitter,
) *PluginHandler {
	return &PluginHandler{
		registry:           reg,
		serviceClassLookup: serviceClassLookup,
		blocker:            blocker,
		emitter:            emitter,
	}
}

// HandleCollection dispatches POST → Create and GET → List.
func (h *PluginHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
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
func (h *PluginHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	remainder := strings.TrimPrefix(r.URL.Path, "/v1/plugins/")
	parts := strings.Split(remainder, "/")
	if len(parts) != 1 || parts[0] == "" {
		writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "plugin not found", "", "")
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

// Create handles POST /v1/plugins.
func (h *PluginHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	p, err := safeDecodePlugin(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if errs := validation.ValidatePlugin(p); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if err := h.verifyServiceClassRefs(ctx, w, r, p.Spec.ServiceClassRefs); err != nil {
		return
	}

	p.APIVersion = resources.PluginAPIVersion
	p.Kind = resources.PluginKind
	p.Status.Phase = resources.PhaseActive
	p.Status.Message = ""

	created, err := h.registry.CreatePlugin(ctx, p)
	if err != nil {
		if errors.Is(err, registry.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, resources.ErrCodeResourceAlreadyExists, "plugin already exists", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:         resources.OpCreatePlugin,
		ResourceKind: resources.PluginKind,
		ResourceName: created.Metadata.Name,
		PluginName:   created.Metadata.Name,
		RequestID:    requestIDFromContext(ctx),
	})
	writeJSON(w, r, http.StatusCreated, created)
}

// Get handles GET /v1/plugins/{name}.
func (h *PluginHandler) Get(w http.ResponseWriter, r *http.Request, name string) {
	if errs := validation.ValidatePluginPathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	p, err := h.registry.GetPlugin(r.Context(), name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "plugin not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	writeJSON(w, r, http.StatusOK, p)
}

// pluginListResponse is the list endpoint response shape.
type pluginListResponse struct {
	Items []resources.Plugin `json:"items"`
}

// List handles GET /v1/plugins.
func (h *PluginHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.registry.ListPlugins(r.Context())
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if items == nil {
		items = []resources.Plugin{}
	}
	writeJSON(w, r, http.StatusOK, pluginListResponse{Items: items})
}

// Update handles PUT /v1/plugins/{name}.
func (h *PluginHandler) Update(w http.ResponseWriter, r *http.Request, name string) {
	ctx := r.Context()

	if errs := validation.ValidatePluginPathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	p, err := safeDecodePlugin(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if p.Metadata.Name == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name is required in request body", "metadata.name", "")
		return
	}
	if p.Metadata.Name != name {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name in body must match path", "metadata.name", "")
		return
	}

	if errs := validation.ValidatePlugin(p); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if err := h.verifyServiceClassRefs(ctx, w, r, p.Spec.ServiceClassRefs); err != nil {
		return
	}

	updated, err := h.registry.UpdatePlugin(ctx, p)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "plugin not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:         resources.OpUpdatePlugin,
		ResourceKind: resources.PluginKind,
		ResourceName: updated.Metadata.Name,
		PluginName:   updated.Metadata.Name,
		RequestID:    requestIDFromContext(ctx),
	})
	writeJSON(w, r, http.StatusOK, updated)
}

// Delete handles DELETE /v1/plugins/{name}.
func (h *PluginHandler) Delete(w http.ResponseWriter, r *http.Request, name string) {
	ctx := r.Context()

	if errs := validation.ValidatePluginPathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if h.blocker != nil {
		blockers, err := h.blocker.BlockedByPluginChildren(ctx, name)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
			return
		}
		if len(blockers) > 0 {
			writeError(w, r, http.StatusConflict, resources.ErrCodeDeleteBlocked,
				"deletion blocked by Capability resources", "", "")
			return
		}
	}

	if err := h.registry.DeletePlugin(ctx, name); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "plugin not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:         resources.OpDeletePlugin,
		ResourceKind: resources.PluginKind,
		ResourceName: name,
		PluginName:   name,
		RequestID:    requestIDFromContext(ctx),
	})
	w.Header().Set("Content-Type", "application/json")
	if reqID := requestctx.RequestIDFromContext(ctx); reqID != "" {
		w.Header().Set("X-Sovrunn-Request-ID", reqID)
	}
	w.WriteHeader(http.StatusNoContent)
}

// verifyServiceClassRefs checks that each ServiceClass reference exists.
// On failure it writes the appropriate HTTP error and returns a non-nil error
// so callers can return immediately. A nil lookup is treated as a server fault.
func (h *PluginHandler) verifyServiceClassRefs(
	ctx context.Context, w http.ResponseWriter, r *http.Request, refs []string,
) error {
	if h.serviceClassLookup == nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return errors.New("service class lookup not configured")
	}
	for _, ref := range refs {
		if _, err := h.serviceClassLookup.GetServiceClass(ctx, ref); err != nil {
			if errors.Is(err, registry.ErrNotFound) {
				writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
					"referenced ServiceClass not found: "+ref,
					"spec.serviceClassRefs", "")
				return err
			}
			writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
			return err
		}
	}
	return nil
}
