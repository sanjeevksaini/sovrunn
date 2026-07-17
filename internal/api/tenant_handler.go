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

// TenantHandler holds dependencies for Tenant CRUD endpoints. ouLookup
// verifies that a referenced parent OrganizationUnit exists; the registry
// stores Tenant state. Both are injected by main.
type TenantHandler struct {
	registry registry.TenantRegistryIface
	ouLookup registry.OrganizationUnitLookup
}

// NewTenantHandler constructs a TenantHandler.
func NewTenantHandler(
	reg registry.TenantRegistryIface,
	ouLookup registry.OrganizationUnitLookup,
) *TenantHandler {
	return &TenantHandler{registry: reg, ouLookup: ouLookup}
}

// HandleCollection dispatches POST → Create and GET → List.
func (h *TenantHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
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
// as "{organizationName}/{organizationUnitName}/{name}".
func (h *TenantHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	remainder := strings.TrimPrefix(r.URL.Path, "/v1/tenants/")
	parts := strings.Split(remainder, "/")
	if remainder == "" || len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "tenant not found", "", "")
		return
	}
	orgName := parts[0]
	ouName := parts[1]
	name := parts[2]

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r, orgName, ouName, name)
	case http.MethodPut:
		h.Update(w, r, orgName, ouName, name)
	case http.MethodDelete:
		h.Delete(w, r, orgName, ouName, name)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// Create handles POST /v1/tenants.
func (h *TenantHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenant, err := safeDecodeTenant(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if errs := validation.ValidateTenant(tenant); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if _, err := h.ouLookup.GetOrganizationUnit(ctx, tenant.Spec.OrganizationName, tenant.Spec.OrganizationUnitName); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
				"parent OrganizationUnit not found: "+tenant.Spec.OrganizationName+"/"+tenant.Spec.OrganizationUnitName,
				"spec.organizationUnitName", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	tenant.APIVersion = resources.TenantAPIVersion
	tenant.Kind = resources.TenantKind
	tenant.Status.Phase = resources.PhaseActive
	tenant.Status.Message = ""

	created, err := h.registry.CreateTenant(ctx, tenant)
	if err != nil {
		if errors.Is(err, registry.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, resources.ErrCodeResourceAlreadyExists, "tenant already exists", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	// TODO(FEATURE-0005): emit Operation record — type: CreateTenant
	writeJSON(w, r, http.StatusCreated, created)
}

// Get handles GET /v1/tenants/{organizationName}/{organizationUnitName}/{name}.
func (h *TenantHandler) Get(w http.ResponseWriter, r *http.Request, orgName, ouName, name string) {
	if errs := validation.ValidateTenantPathSegments(orgName, ouName, name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	tenant, err := h.registry.GetTenant(r.Context(), orgName, ouName, name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "tenant not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	writeJSON(w, r, http.StatusOK, tenant)
}

// tenantListResponse is the list endpoint response shape.
type tenantListResponse struct {
	Items []resources.Tenant `json:"items"`
}

// List handles GET /v1/tenants.
func (h *TenantHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.registry.ListTenants(r.Context())
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if items == nil {
		items = []resources.Tenant{}
	}
	writeJSON(w, r, http.StatusOK, tenantListResponse{Items: items})
}

// Update handles PUT /v1/tenants/{organizationName}/{organizationUnitName}/{name}.
func (h *TenantHandler) Update(w http.ResponseWriter, r *http.Request, orgName, ouName, name string) {
	if errs := validation.ValidateTenantPathSegments(orgName, ouName, name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	tenant, err := safeDecodeTenant(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if tenant.Metadata.Name == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name is required in request body", "metadata.name", "")
		return
	}
	if tenant.Metadata.Name != name {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name in body must match path", "metadata.name", "")
		return
	}
	if tenant.Spec.OrganizationName == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.organizationName is required in request body", "spec.organizationName", "")
		return
	}
	if tenant.Spec.OrganizationName != orgName {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.organizationName in body must match path", "spec.organizationName", "")
		return
	}
	if tenant.Spec.OrganizationUnitName == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.organizationUnitName is required in request body", "spec.organizationUnitName", "")
		return
	}
	if tenant.Spec.OrganizationUnitName != ouName {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.organizationUnitName in body must match path", "spec.organizationUnitName", "")
		return
	}

	if errs := validation.ValidateTenant(tenant); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	updated, err := h.registry.UpdateTenant(r.Context(), tenant)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "tenant not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	// TODO(FEATURE-0005): emit Operation record — type: UpdateTenant
	writeJSON(w, r, http.StatusOK, updated)
}

// Delete handles DELETE /v1/tenants/{organizationName}/{organizationUnitName}/{name}.
func (h *TenantHandler) Delete(w http.ResponseWriter, r *http.Request, orgName, ouName, name string) {
	if errs := validation.ValidateTenantPathSegments(orgName, ouName, name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if err := h.registry.DeleteTenant(r.Context(), orgName, ouName, name); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "tenant not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	// TODO(FEATURE-0005): emit Operation record — type: DeleteTenant
	w.Header().Set("Content-Type", "application/json")
	if reqID := requestctx.RequestIDFromContext(r.Context()); reqID != "" {
		w.Header().Set("X-Sovrunn-Request-ID", reqID)
	}
	w.WriteHeader(http.StatusNoContent)
}
