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

// ServicePlanHandler holds dependencies for ServicePlan CRUD endpoints.
// emitter and instanceBlocker are optional (nil-safe). serviceClassLookup
// verifies parent ServiceClass existence on create and update.
type ServicePlanHandler struct {
	registry           registry.ServicePlanRegistryIface
	serviceClassLookup registry.ServiceClassLookup
	instanceBlocker    registry.ServicePlanInstanceBlocker
	emitter            OperationEmitter
}

// NewServicePlanHandler constructs a ServicePlanHandler. emitter and the
// optional instanceBlocker may be nil.
func NewServicePlanHandler(
	reg registry.ServicePlanRegistryIface,
	serviceClassLookup registry.ServiceClassLookup,
	emitter OperationEmitter,
	instanceBlocker ...registry.ServicePlanInstanceBlocker,
) *ServicePlanHandler {
	h := &ServicePlanHandler{
		registry:           reg,
		serviceClassLookup: serviceClassLookup,
		emitter:            emitter,
	}
	if len(instanceBlocker) > 0 {
		h.instanceBlocker = instanceBlocker[0]
	}
	return h
}

// HandleCollection dispatches POST → Create and GET → List.
func (h *ServicePlanHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
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
// as "{serviceClassName}/{name}".
func (h *ServicePlanHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	remainder := strings.TrimPrefix(r.URL.Path, "/v1/service-plans/")
	parts := strings.Split(remainder, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service plan not found", "", "")
		return
	}
	serviceClassName := parts[0]
	name := parts[1]

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r, serviceClassName, name)
	case http.MethodPut:
		h.Update(w, r, serviceClassName, name)
	case http.MethodDelete:
		h.Delete(w, r, serviceClassName, name)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// Create handles POST /v1/service-plans.
func (h *ServicePlanHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sp, err := safeDecodeServicePlan(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if errs := validation.ValidateServicePlan(sp); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if h.serviceClassLookup == nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if _, err := h.serviceClassLookup.GetServiceClass(ctx, sp.Spec.ServiceClassName); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
				"parent ServiceClass not found: "+sp.Spec.ServiceClassName,
				"spec.serviceClassName", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	sp.APIVersion = resources.ServicePlanAPIVersion
	sp.Kind = resources.ServicePlanKind
	sp.Status.Phase = resources.PhaseActive
	sp.Status.Message = ""

	created, err := h.registry.CreateServicePlan(ctx, sp)
	if err != nil {
		if errors.Is(err, registry.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, resources.ErrCodeResourceAlreadyExists, "service plan already exists", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:             resources.OpCreateServicePlan,
		ResourceKind:     resources.ServicePlanKind,
		ResourceName:     created.Metadata.Name,
		ServiceClassName: created.Spec.ServiceClassName,
		ServicePlanName:  created.Metadata.Name,
		RequestID:        requestIDFromContext(ctx),
	})
	writeJSON(w, r, http.StatusCreated, created)
}

// Get handles GET /v1/service-plans/{serviceClassName}/{name}.
func (h *ServicePlanHandler) Get(w http.ResponseWriter, r *http.Request, serviceClassName, name string) {
	if errs := validation.ValidateServicePlanPathSegments(serviceClassName, name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	sp, err := h.registry.GetServicePlan(r.Context(), serviceClassName, name)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service plan not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	writeJSON(w, r, http.StatusOK, sp)
}

// servicePlanListResponse is the list endpoint response shape.
type servicePlanListResponse struct {
	Items []resources.ServicePlan `json:"items"`
}

// List handles GET /v1/service-plans.
func (h *ServicePlanHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.registry.ListServicePlans(r.Context())
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if items == nil {
		items = []resources.ServicePlan{}
	}
	writeJSON(w, r, http.StatusOK, servicePlanListResponse{Items: items})
}

// Update handles PUT /v1/service-plans/{serviceClassName}/{name}.
func (h *ServicePlanHandler) Update(w http.ResponseWriter, r *http.Request, serviceClassName, name string) {
	ctx := r.Context()

	if errs := validation.ValidateServicePlanPathSegments(serviceClassName, name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	sp, err := safeDecodeServicePlan(w, r)
	if err != nil {
		status, msg := mapDecodeError(err)
		field := ""
		if errors.Is(err, errStatusFieldPresent) {
			field = "status"
		}
		writeError(w, r, status, resources.ErrCodeValidationFailed, msg, field, "")
		return
	}

	if sp.Spec.ServiceClassName == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.serviceClassName is required in request body", "spec.serviceClassName", "")
		return
	}
	if sp.Spec.ServiceClassName != serviceClassName {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"spec.serviceClassName in body must match path", "spec.serviceClassName", "")
		return
	}
	if sp.Metadata.Name == "" {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name is required in request body", "metadata.name", "")
		return
	}
	if sp.Metadata.Name != name {
		writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
			"metadata.name in body must match path", "metadata.name", "")
		return
	}

	if errs := validation.ValidateServicePlan(sp); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if h.serviceClassLookup == nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if _, err := h.serviceClassLookup.GetServiceClass(ctx, serviceClassName); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusBadRequest, resources.ErrCodeValidationFailed,
				"parent ServiceClass not found: "+serviceClassName,
				"spec.serviceClassName", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	updated, err := h.registry.UpdateServicePlan(ctx, sp)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service plan not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:             resources.OpUpdateServicePlan,
		ResourceKind:     resources.ServicePlanKind,
		ResourceName:     updated.Metadata.Name,
		ServiceClassName: updated.Spec.ServiceClassName,
		ServicePlanName:  updated.Metadata.Name,
		RequestID:        requestIDFromContext(ctx),
	})
	writeJSON(w, r, http.StatusOK, updated)
}

// Delete handles DELETE /v1/service-plans/{serviceClassName}/{name}.
func (h *ServicePlanHandler) Delete(w http.ResponseWriter, r *http.Request, serviceClassName, name string) {
	ctx := r.Context()

	if errs := validation.ValidateServicePlanPathSegments(serviceClassName, name); len(errs) > 0 {
		writeValidationErrors(w, r, errs)
		return
	}

	if h.instanceBlocker != nil {
		blockers, err := h.instanceBlocker.BlockedByServicePlanInstances(ctx, serviceClassName, name)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
			return
		}
		if len(blockers) > 0 {
			writeError(w, r, http.StatusConflict, resources.ErrCodeDeleteBlocked,
				"deletion blocked by ServiceInstance resources", "", "")
			return
		}
	}

	if err := h.registry.DeleteServicePlan(ctx, serviceClassName, name); err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "service plan not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}

	emitOperation(ctx, h.emitter, resources.OperationSpec{
		Type:             resources.OpDeleteServicePlan,
		ResourceKind:     resources.ServicePlanKind,
		ResourceName:     name,
		ServiceClassName: serviceClassName,
		ServicePlanName:  name,
		RequestID:        requestIDFromContext(ctx),
	})
	w.Header().Set("Content-Type", "application/json")
	if reqID := requestctx.RequestIDFromContext(ctx); reqID != "" {
		w.Header().Set("X-Sovrunn-Request-ID", reqID)
	}
	w.WriteHeader(http.StatusNoContent)
}
