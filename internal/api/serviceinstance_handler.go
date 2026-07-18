package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/requestctx"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
	"github.com/sanjeevksaini/sovrunn/internal/validation"
)

// ServiceInstanceHandler holds dependencies for ServiceInstance CRUD endpoints.
// All lookup and blocker fields are required (non-nil). emitter and logger may
// be nil.
type ServiceInstanceHandler struct {
	registry           registry.ServiceInstanceRegistryIface
	orgLookup          registry.OrganizationLookup
	ouLookup           registry.OrganizationUnitLookup
	tenantLookup       registry.TenantLookup
	projectLookup      registry.ProjectLookup
	serviceClassLookup registry.ServiceClassLookup
	servicePlanLookup  registry.ServicePlanLookup
	capabilityLookup   registry.CapabilityLookup
	bindingBlocker     registry.ServiceBindingInstanceBlocker
	emitter            OperationEmitter
	logger             *log.Logger
}

// NewServiceInstanceHandler constructs a ServiceInstanceHandler. emitter and
// logger may be nil; all other dependencies must be non-nil.
func NewServiceInstanceHandler(
	reg registry.ServiceInstanceRegistryIface,
	orgLookup registry.OrganizationLookup,
	ouLookup registry.OrganizationUnitLookup,
	tenantLookup registry.TenantLookup,
	projectLookup registry.ProjectLookup,
	serviceClassLookup registry.ServiceClassLookup,
	servicePlanLookup registry.ServicePlanLookup,
	capabilityLookup registry.CapabilityLookup,
	bindingBlocker registry.ServiceBindingInstanceBlocker,
	emitter OperationEmitter,
	logger *log.Logger,
) *ServiceInstanceHandler {
	return &ServiceInstanceHandler{
		registry:           reg,
		orgLookup:          orgLookup,
		ouLookup:           ouLookup,
		tenantLookup:       tenantLookup,
		projectLookup:      projectLookup,
		serviceClassLookup: serviceClassLookup,
		servicePlanLookup:  servicePlanLookup,
		capabilityLookup:   capabilityLookup,
		bindingBlocker:     bindingBlocker,
		emitter:            emitter,
		logger:             logger,
	}
}

// HandleCollection dispatches POST → Create and GET → List on
// /v1/service-instances.
func (h *ServiceInstanceHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
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
func (h *ServiceInstanceHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	remainder := strings.TrimPrefix(r.URL.Path, "/v1/service-instances/")
	parts := strings.Split(remainder, "/")
	if len(parts) != 1 || parts[0] == "" {
		writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service instance not found", "", "")
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

// Create handles POST /v1/service-instances.
func (h *ServiceInstanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	si, err := safeDecodeServiceInstance(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if errs := validation.ValidateServiceInstance(si); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if !h.validateReferences(ctx, w, r, si) {
		return
	}

	h.warnIfNoActiveCapability(ctx, si)

	si.APIVersion = resources.ServiceInstanceAPIVersion
	si.Kind = resources.ServiceInstanceKind
	si.Status.Phase = "Ready"
	si.Status.Message = "Registered only; no real provisioning in Phase 1"

	created, err := h.registry.CreateServiceInstance(ctx, si)
	if err != nil {
		if errors.Is(err, registry.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, resources.ErrCodeResourceAlreadyExists, "service instance already exists", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:                 resources.OpCreateServiceInstance,
		ResourceKind:         resources.ServiceInstanceKind,
		ResourceName:         created.Metadata.Name,
		OrganizationName:     created.Spec.OrganizationRef,
		OrganizationUnitName: created.Spec.OrganizationUnitRef,
		TenantName:           created.Spec.TenantRef,
		ProjectName:          created.Spec.ProjectRef,
		ServiceInstanceName:  created.Metadata.Name,
		RequestID:            requestIDFromContext(ctx),
	})
	writeJSON(w, r, http.StatusCreated, created)
}

// Get handles GET /v1/service-instances/{name}.
func (h *ServiceInstanceHandler) Get(w http.ResponseWriter, r *http.Request, name string) {
	if errs := validation.ValidateServiceInstancePathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	si, err := h.registry.GetServiceInstance(r.Context(), name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service instance not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	writeJSON(w, r, http.StatusOK, si)
}

// serviceInstanceListResponse is the list endpoint response shape.
type serviceInstanceListResponse struct {
	Items []resources.ServiceInstance `json:"items"`
}

// List handles GET /v1/service-instances with optional tenantRef and projectRef
// query filters (AND logic when both are present).
func (h *ServiceInstanceHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantRef := r.URL.Query().Get("tenantRef")
	projectRef := r.URL.Query().Get("projectRef")

	items, err := h.registry.ListServiceInstances(r.Context(), tenantRef, projectRef)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if items == nil {
		items = []resources.ServiceInstance{}
	}
	writeJSON(w, r, http.StatusOK, serviceInstanceListResponse{Items: items})
}

// Update handles PUT /v1/service-instances/{name}.
func (h *ServiceInstanceHandler) Update(w http.ResponseWriter, r *http.Request, name string) {
	ctx := r.Context()

	if errs := validation.ValidateServiceInstancePathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	si, err := safeDecodeServiceInstance(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if si.Metadata.Name == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name is required in request body", "metadata.name", "")
		return
	}
	if si.Metadata.Name != name {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name in body must match path", "metadata.name", "")
		return
	}

	if errs := validation.ValidateServiceInstance(si); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	stored, err := h.registry.GetServiceInstance(ctx, name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service instance not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	if field, msg := immutableServiceInstanceFieldChanged(stored, si); field != "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	updated, err := h.registry.UpdateServiceInstance(ctx, name, si)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service instance not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:                 resources.OpUpdateServiceInstance,
		ResourceKind:         resources.ServiceInstanceKind,
		ResourceName:         updated.Metadata.Name,
		OrganizationName:     updated.Spec.OrganizationRef,
		OrganizationUnitName: updated.Spec.OrganizationUnitRef,
		TenantName:           updated.Spec.TenantRef,
		ProjectName:          updated.Spec.ProjectRef,
		ServiceInstanceName:  updated.Metadata.Name,
		RequestID:            requestIDFromContext(ctx),
	})
	writeJSON(w, r, http.StatusOK, updated)
}

// Delete handles DELETE /v1/service-instances/{name}.
func (h *ServiceInstanceHandler) Delete(w http.ResponseWriter, r *http.Request, name string) {
	ctx := r.Context()

	if errs := validation.ValidateServiceInstancePathSegment(name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	count, err := h.bindingBlocker.CountByServiceInstance(ctx, name)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if count > 0 {
		writeError(w, r, http.StatusConflict, resources.ErrCodeDeleteBlocked,
			"deletion blocked by ServiceBinding resources", "", "")
		return
	}

	stored, err := h.registry.GetServiceInstance(ctx, name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service instance not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	if err := h.registry.DeleteServiceInstance(ctx, name); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service instance not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:                 resources.OpDeleteServiceInstance,
		ResourceKind:         resources.ServiceInstanceKind,
		ResourceName:         name,
		OrganizationName:     stored.Spec.OrganizationRef,
		OrganizationUnitName: stored.Spec.OrganizationUnitRef,
		TenantName:           stored.Spec.TenantRef,
		ProjectName:          stored.Spec.ProjectRef,
		ServiceInstanceName:  name,
		RequestID:            requestIDFromContext(ctx),
	})
	w.Header().Set("Content-Type", "application/json")
	if reqID := requestctx.RequestIDFromContext(ctx); reqID != "" {
		w.Header().Set("X-Sovrunn-Request-ID", reqID)
	}
	w.WriteHeader(http.StatusNoContent)
}

// validateReferences checks governance and catalog references in create order.
// Returns false when an error response has already been written.
func (h *ServiceInstanceHandler) validateReferences(
	ctx context.Context, w http.ResponseWriter, r *http.Request, si resources.ServiceInstance,
) bool {
	if _, err := h.orgLookup.GetOrganization(ctx, si.Spec.OrganizationRef); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
				"referenced Organization not found: "+si.Spec.OrganizationRef,
				"spec.organizationRef", "")
			return false
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return false
	}

	if si.Spec.OrganizationUnitRef != "" {
		if _, err := h.ouLookup.GetOrganizationUnit(ctx, si.Spec.OrganizationRef, si.Spec.OrganizationUnitRef); err != nil {
			if errors.Is(err, registry.ErrNotFound) {
				writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
					"referenced OrganizationUnit not found: "+si.Spec.OrganizationRef+"/"+si.Spec.OrganizationUnitRef,
					"spec.organizationUnitRef", "")
				return false
			}
			writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
			return false
		}
	}

	if _, err := h.tenantLookup.GetTenant(ctx, si.Spec.OrganizationRef, si.Spec.OrganizationUnitRef, si.Spec.TenantRef); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
				"referenced Tenant not found: "+si.Spec.OrganizationRef+"/"+si.Spec.OrganizationUnitRef+"/"+si.Spec.TenantRef,
				"spec.tenantRef", "")
			return false
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return false
	}

	if _, err := h.projectLookup.GetProject(ctx, si.Spec.OrganizationRef, si.Spec.OrganizationUnitRef, si.Spec.TenantRef, si.Spec.ProjectRef); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
				"referenced Project not found: "+si.Spec.OrganizationRef+"/"+si.Spec.OrganizationUnitRef+"/"+si.Spec.TenantRef+"/"+si.Spec.ProjectRef,
				"spec.projectRef", "")
			return false
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return false
	}

	if _, err := h.serviceClassLookup.GetServiceClass(ctx, si.Spec.ServiceClassRef); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
				"referenced ServiceClass not found: "+si.Spec.ServiceClassRef,
				"spec.serviceClassRef", "")
			return false
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return false
	}

	if _, err := h.servicePlanLookup.GetServicePlan(ctx, si.Spec.ServiceClassRef, si.Spec.ServicePlanRef); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
				"referenced ServicePlan not found or does not belong to ServiceClass: "+si.Spec.ServiceClassRef+"/"+si.Spec.ServicePlanRef,
				"spec.servicePlanRef", "")
			return false
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return false
	}

	return true
}

// warnIfNoActiveCapability logs a structured warning when no active Capability
// exists for the ServiceClass. Creation is never blocked by this check.
func (h *ServiceInstanceHandler) warnIfNoActiveCapability(ctx context.Context, si resources.ServiceInstance) {
	has, err := h.capabilityLookup.HasActiveCapabilityForServiceClass(ctx, si.Spec.ServiceClassRef)
	if err != nil {
		if h.logger != nil {
			h.logger.Printf("level=warn msg=%q service_instance=%q service_class=%q error=%q",
				"capability lookup failed; proceeding with ServiceInstance create",
				si.Metadata.Name, si.Spec.ServiceClassRef, err.Error())
		}
		return
	}
	if has {
		return
	}
	if h.logger != nil {
		h.logger.Printf("level=warn msg=%q service_instance=%q service_class=%q",
			"no active capability registered for ServiceClass",
			si.Metadata.Name, si.Spec.ServiceClassRef)
	}
}

// immutableServiceInstanceFieldChanged compares immutable governance and
// catalog fields. It returns the first changed field path and a message, or
// empty strings when all immutable fields match.
func immutableServiceInstanceFieldChanged(stored, incoming resources.ServiceInstance) (field, message string) {
	switch {
	case stored.Spec.OrganizationRef != incoming.Spec.OrganizationRef:
		return "spec.organizationRef", "spec.organizationRef is immutable"
	case stored.Spec.OrganizationUnitRef != incoming.Spec.OrganizationUnitRef:
		return "spec.organizationUnitRef", "spec.organizationUnitRef is immutable"
	case stored.Spec.TenantRef != incoming.Spec.TenantRef:
		return "spec.tenantRef", "spec.tenantRef is immutable"
	case stored.Spec.ProjectRef != incoming.Spec.ProjectRef:
		return "spec.projectRef", "spec.projectRef is immutable"
	case stored.Spec.ServiceClassRef != incoming.Spec.ServiceClassRef:
		return "spec.serviceClassRef", "spec.serviceClassRef is immutable"
	case stored.Spec.ServicePlanRef != incoming.Spec.ServicePlanRef:
		return "spec.servicePlanRef", "spec.servicePlanRef is immutable"
	default:
		return "", ""
	}
}
