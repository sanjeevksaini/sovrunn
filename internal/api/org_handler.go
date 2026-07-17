package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/requestctx"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
	"github.com/sanjeevksaini/sovrunn/internal/validation"
)

// OrgHandler holds dependencies injected by main. emitter is optional
// (nil-safe) and records Operations after successful mutations in a later task.
type OrgHandler struct {
	registry registry.OrganizationRegistryIface
	blocker  registry.ChildBlockerChecker
	emitter  OperationEmitter
}

// NewOrgHandler constructs an OrgHandler. emitter may be nil.
func NewOrgHandler(
	reg registry.OrganizationRegistryIface,
	blocker registry.ChildBlockerChecker,
	emitter OperationEmitter,
) *OrgHandler {
	return &OrgHandler{registry: reg, blocker: blocker, emitter: emitter}
}

// HandleCollection dispatches POST → Create and GET → List.
func (h *OrgHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Create(w, r)
	case http.MethodGet:
		h.List(w, r)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// HandleItem dispatches GET → Get, PUT → Update, DELETE → Delete.
func (h *OrgHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/v1/organizations/")
	if name == "" {
		writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "organization not found", "", "")
		return
	}
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

// Create handles POST /v1/organizations.
func (h *OrgHandler) Create(w http.ResponseWriter, r *http.Request) {
	org, err := safeDecodeOrganization(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if errs := validation.ValidateOrganization(r.Context(), org); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	org.APIVersion = resources.OrgAPIVersion
	org.Kind = resources.OrgKind
	org.Status.Phase = resources.PhaseActive

	if err := h.registry.CreateOrganization(r.Context(), org); err != nil {
		if errors.Is(err, registry.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, resources.ErrCodeResourceAlreadyExists, "organization already exists", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(r.Context(), h.emitter, resources.OperationSpec{
		Type:             resources.OpCreateOrganization,
		ResourceKind:     resources.OrganizationKind,
		ResourceName:     org.Metadata.Name,
		OrganizationName: org.Metadata.Name,
		RequestID:        requestIDFromContext(r.Context()),
	})
	writeJSON(w, r, http.StatusCreated, org)
}

// Get handles GET /v1/organizations/{name}.
func (h *OrgHandler) Get(w http.ResponseWriter, r *http.Request, name string) {
	if errs := validation.ValidateNamePath(r.Context(), name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	org, err := h.registry.GetOrganization(r.Context(), name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "organization not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	writeJSON(w, r, http.StatusOK, org)
}

// organizationListResponse is the list endpoint response shape.
type organizationListResponse struct {
	Items []resources.Organization `json:"items"`
}

// List handles GET /v1/organizations.
func (h *OrgHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.registry.ListOrganizations(r.Context())
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if items == nil {
		items = []resources.Organization{}
	}
	writeJSON(w, r, http.StatusOK, organizationListResponse{Items: items})
}

// Update handles PUT /v1/organizations/{name}.
func (h *OrgHandler) Update(w http.ResponseWriter, r *http.Request, name string) {
	if errs := validation.ValidateNamePath(r.Context(), name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	org, err := safeDecodeOrganization(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if org.Metadata.Name == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name is required in request body", "metadata.name", "")
		return
	}
	if org.Metadata.Name != name {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name in body must match path", "metadata.name", "")
		return
	}

	if errs := validation.ValidateOrganization(r.Context(), org); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	updated, err := h.registry.UpdateOrganization(r.Context(), name, org)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "organization not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(r.Context(), h.emitter, resources.OperationSpec{
		Type:             resources.OpUpdateOrganization,
		ResourceKind:     resources.OrganizationKind,
		ResourceName:     updated.Metadata.Name,
		OrganizationName: updated.Metadata.Name,
		RequestID:        requestIDFromContext(r.Context()),
	})
	writeJSON(w, r, http.StatusOK, updated)
}

// Delete handles DELETE /v1/organizations/{name}.
func (h *OrgHandler) Delete(w http.ResponseWriter, r *http.Request, name string) {
	if errs := validation.ValidateNamePath(r.Context(), name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	_, err := h.registry.GetOrganization(r.Context(), name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "organization not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	blockers, err := h.blocker.BlockedByChildren(r.Context(), name)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if len(blockers) > 0 {
		msg := fmt.Sprintf("deletion blocked by %s resources", blockers[0].Kind)
		writeError(w, r, http.StatusConflict, resources.ErrCodeDeleteBlocked, msg, "", "")
		return
	}

	if err := h.registry.DeleteOrganization(r.Context(), name); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "organization not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(r.Context(), h.emitter, resources.OperationSpec{
		Type:             resources.OpDeleteOrganization,
		ResourceKind:     resources.OrganizationKind,
		ResourceName:     name,
		OrganizationName: name,
		RequestID:        requestIDFromContext(r.Context()),
	})
	w.Header().Set("Content-Type", "application/json")
	if reqID := requestctx.RequestIDFromContext(r.Context()); reqID != "" {
		w.Header().Set("X-Sovrunn-Request-ID", reqID)
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeValidationErrors(w http.ResponseWriter, r *http.Request, errs []resources.FieldError) {
	if len(errs) == 0 {
		return
	}
	field := errs[0].Field
	message := errs[0].Message
	details := ""
	if len(errs) > 1 {
		parts := make([]string, 0, len(errs)-1)
		for _, e := range errs[1:] {
			parts = append(parts, e.Field+": "+e.Message)
		}
		details = strings.Join(parts, "; ")
	}
	writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed, message, field, details)
}
