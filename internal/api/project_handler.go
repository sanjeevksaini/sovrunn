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

// ProjectHandler holds dependencies for Project CRUD endpoints. tenantLookup
// verifies that a referenced parent Tenant exists; the registry stores Project
// state. Both are injected by main.
type ProjectHandler struct {
	registry     registry.ProjectRegistryIface
	tenantLookup registry.TenantLookup
}

// NewProjectHandler constructs a ProjectHandler.
func NewProjectHandler(
	reg registry.ProjectRegistryIface,
	tenantLookup registry.TenantLookup,
) *ProjectHandler {
	return &ProjectHandler{registry: reg, tenantLookup: tenantLookup}
}

// HandleCollection dispatches POST → Create and GET → List.
func (h *ProjectHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
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
// Go 1.21-compatible path parsing. The composite key is carried in the path
// as "{organizationName}/{organizationUnitName}/{tenantName}/{name}".
func (h *ProjectHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	remainder := strings.TrimPrefix(r.URL.Path, "/v1/projects/")
	parts := strings.Split(remainder, "/")
	if remainder == "" || len(parts) != 4 ||
		parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" {
		writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "project not found", "", "")
		return
	}
	orgName := parts[0]
	ouName := parts[1]
	tenantName := parts[2]
	name := parts[3]

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r, orgName, ouName, tenantName, name)
	case http.MethodPut:
		h.Update(w, r, orgName, ouName, tenantName, name)
	case http.MethodDelete:
		h.Delete(w, r, orgName, ouName, tenantName, name)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// Create handles POST /v1/projects.
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	project, err := safeDecodeProject(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if errs := validation.ValidateProject(project); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if _, err := h.tenantLookup.GetTenant(ctx, project.Spec.OrganizationName, project.Spec.OrganizationUnitName, project.Spec.TenantName); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
				"parent Tenant not found: "+project.Spec.OrganizationName+"/"+project.Spec.OrganizationUnitName+"/"+project.Spec.TenantName,
				"spec.tenantName", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	project.APIVersion = resources.ProjectAPIVersion
	project.Kind = resources.ProjectKind
	project.Status.Phase = resources.PhaseActive
	project.Status.Message = ""

	created, err := h.registry.CreateProject(ctx, project)
	if err != nil {
		if errors.Is(err, registry.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, resources.ErrCodeResourceAlreadyExists, "project already exists", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	// TODO(FEATURE-0005): emit Operation record — type: CreateProject
	writeJSON(w, r, http.StatusCreated, created)
}

// Get handles GET /v1/projects/{organizationName}/{organizationUnitName}/{tenantName}/{name}.
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request, orgName, ouName, tenantName, name string) {
	if errs := validation.ValidateProjectPathSegments(orgName, ouName, tenantName, name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	project, err := h.registry.GetProject(r.Context(), orgName, ouName, tenantName, name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "project not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	writeJSON(w, r, http.StatusOK, project)
}

// projectListResponse is the list endpoint response shape.
type projectListResponse struct {
	Items []resources.Project `json:"items"`
}

// List handles GET /v1/projects.
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.registry.ListProjects(r.Context())
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if items == nil {
		items = []resources.Project{}
	}
	writeJSON(w, r, http.StatusOK, projectListResponse{Items: items})
}

// Update handles PUT /v1/projects/{organizationName}/{organizationUnitName}/{tenantName}/{name}.
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request, orgName, ouName, tenantName, name string) {
	if errs := validation.ValidateProjectPathSegments(orgName, ouName, tenantName, name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	project, err := safeDecodeProject(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if project.Metadata.Name == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name is required in request body", "metadata.name", "")
		return
	}
	if project.Metadata.Name != name {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name in body must match path", "metadata.name", "")
		return
	}
	if project.Spec.OrganizationName == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.organizationName is required in request body", "spec.organizationName", "")
		return
	}
	if project.Spec.OrganizationName != orgName {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.organizationName in body must match path", "spec.organizationName", "")
		return
	}
	if project.Spec.OrganizationUnitName == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.organizationUnitName is required in request body", "spec.organizationUnitName", "")
		return
	}
	if project.Spec.OrganizationUnitName != ouName {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.organizationUnitName in body must match path", "spec.organizationUnitName", "")
		return
	}
	if project.Spec.TenantName == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.tenantName is required in request body", "spec.tenantName", "")
		return
	}
	if project.Spec.TenantName != tenantName {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.tenantName in body must match path", "spec.tenantName", "")
		return
	}

	if errs := validation.ValidateProject(project); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	updated, err := h.registry.UpdateProject(r.Context(), project)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "project not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	// TODO(FEATURE-0005): emit Operation record — type: UpdateProject
	writeJSON(w, r, http.StatusOK, updated)
}

// Delete handles DELETE /v1/projects/{organizationName}/{organizationUnitName}/{tenantName}/{name}.
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request, orgName, ouName, tenantName, name string) {
	if errs := validation.ValidateProjectPathSegments(orgName, ouName, tenantName, name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if err := h.registry.DeleteProject(r.Context(), orgName, ouName, tenantName, name); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "project not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	// TODO(FEATURE-0005): emit Operation record — type: DeleteProject
	w.Header().Set("Content-Type", "application/json")
	if reqID := requestctx.RequestIDFromContext(r.Context()); reqID != "" {
		w.Header().Set("X-Sovrunn-Request-ID", reqID)
	}
	w.WriteHeader(http.StatusNoContent)
}
