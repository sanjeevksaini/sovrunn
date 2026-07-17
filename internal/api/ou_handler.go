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

// OUHandler holds dependencies for OrganizationUnit CRUD endpoints.
// orgLookup verifies that a referenced parent Organization exists; the
// registry stores OrganizationUnit state. blocker is optional and prevents
// deleting OrganizationUnits with child resources.
type OUHandler struct {
	registry  registry.OrganizationUnitRegistryIface
	orgLookup registry.OrganizationLookup
	blocker   registry.OUChildBlocker
}

// NewOUHandler constructs an OUHandler.
func NewOUHandler(
	reg registry.OrganizationUnitRegistryIface,
	orgLookup registry.OrganizationLookup,
	blocker registry.OUChildBlocker,
) *OUHandler {
	return &OUHandler{registry: reg, orgLookup: orgLookup, blocker: blocker}
}

// HandleCollection dispatches POST → Create and GET → List.
func (h *OUHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
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
// Go 1.21-compatible path parsing. The composite key is carried in the
// path as "{organizationName}/{name}".
func (h *OUHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	remainder := strings.TrimPrefix(r.URL.Path, "/v1/organization-units/")
	parts := strings.Split(remainder, "/")
	if remainder == "" || len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "organization unit not found", "", "")
		return
	}
	orgName := parts[0]
	name := parts[1]

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r, orgName, name)
	case http.MethodPut:
		h.Update(w, r, orgName, name)
	case http.MethodDelete:
		h.Delete(w, r, orgName, name)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// Create handles POST /v1/organization-units.
func (h *OUHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ou, err := safeDecodeOrganizationUnit(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if errs := validation.ValidateOrganizationUnit(ou); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if _, err := h.orgLookup.GetOrganization(ctx, ou.Spec.OrganizationName); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
				"referenced organization does not exist", "spec.organizationName", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	ou.APIVersion = resources.OUAPIVersion
	ou.Kind = resources.OUKind
	ou.Status.Phase = resources.PhaseActive

	created, err := h.registry.CreateOrganizationUnit(ctx, ou)
	if err != nil {
		if errors.Is(err, registry.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, resources.ErrCodeResourceAlreadyExists, "organization unit already exists", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	// TODO(FEATURE-0005): emit Operation record — type: CreateOrganizationUnit
	writeJSON(w, r, http.StatusCreated, created)
}

// Get handles GET /v1/organization-units/{organizationName}/{name}.
func (h *OUHandler) Get(w http.ResponseWriter, r *http.Request, orgName, name string) {
	if errs := validation.ValidateOUPathSegments(orgName, name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	ou, err := h.registry.GetOrganizationUnit(r.Context(), orgName, name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "organization unit not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	writeJSON(w, r, http.StatusOK, ou)
}

// organizationUnitListResponse is the list endpoint response shape.
type organizationUnitListResponse struct {
	Items []resources.OrganizationUnit `json:"items"`
}

// List handles GET /v1/organization-units.
func (h *OUHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.registry.ListOrganizationUnits(r.Context())
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if items == nil {
		items = []resources.OrganizationUnit{}
	}
	writeJSON(w, r, http.StatusOK, organizationUnitListResponse{Items: items})
}

// Update handles PUT /v1/organization-units/{organizationName}/{name}.
func (h *OUHandler) Update(w http.ResponseWriter, r *http.Request, orgName, name string) {
	if errs := validation.ValidateOUPathSegments(orgName, name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	ou, err := safeDecodeOrganizationUnit(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if ou.Metadata.Name == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name is required in request body", "metadata.name", "")
		return
	}
	if ou.Metadata.Name != name {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name in body must match path", "metadata.name", "")
		return
	}
	if ou.Spec.OrganizationName == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.organizationName is required in request body", "spec.organizationName", "")
		return
	}
	if ou.Spec.OrganizationName != orgName {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.organizationName in body must match path", "spec.organizationName", "")
		return
	}

	if errs := validation.ValidateOrganizationUnit(ou); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	updated, err := h.registry.UpdateOrganizationUnit(r.Context(), orgName, name, ou)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "organization unit not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	// TODO(FEATURE-0005): emit Operation record — type: UpdateOrganizationUnit
	writeJSON(w, r, http.StatusOK, updated)
}

// Delete handles DELETE /v1/organization-units/{organizationName}/{name}.
func (h *OUHandler) Delete(w http.ResponseWriter, r *http.Request, orgName, name string) {
	if errs := validation.ValidateOUPathSegments(orgName, name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if h.blocker != nil {
		blockers, err := h.blocker.BlockedByOUChildren(r.Context(), orgName, name)
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

	if err := h.registry.DeleteOrganizationUnit(r.Context(), orgName, name); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "organization unit not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	// TODO(FEATURE-0005): emit Operation record — type: DeleteOrganizationUnit
	w.Header().Set("Content-Type", "application/json")
	if reqID := requestctx.RequestIDFromContext(r.Context()); reqID != "" {
		w.Header().Set("X-Sovrunn-Request-ID", reqID)
	}
	w.WriteHeader(http.StatusNoContent)
}
