package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// OperationHandler holds dependencies for read-only Operation endpoints.
type OperationHandler struct {
	registry registry.OperationRegistryIface
}

// NewOperationHandler constructs an OperationHandler.
func NewOperationHandler(reg registry.OperationRegistryIface) *OperationHandler {
	return &OperationHandler{registry: reg}
}

// HandleCollection dispatches GET -> List. POST and all other methods are
// rejected because Operations are system-generated only.
func (h *OperationHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.List(w, r)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// HandleItem dispatches GET -> Get using Go 1.21-compatible path parsing.
func (h *OperationHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	remainder := strings.TrimPrefix(r.URL.Path, "/v1/operations/")
	if remainder == "" || remainder == r.URL.Path || strings.Contains(remainder, "/") {
		writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "operation not found", "", "")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r, remainder)
	default:
		writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeValidationFailed, "method not allowed", "", "")
	}
}

// operationListResponse is the list endpoint response shape.
type operationListResponse struct {
	Items []resources.Operation `json:"items"`
}

// List handles GET /v1/operations.
func (h *OperationHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.registry.ListOperations(r.Context())
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	if items == nil {
		items = []resources.Operation{}
	}
	writeJSON(w, r, http.StatusOK, operationListResponse{Items: items})
}

// Get handles GET /v1/operations/{name}. Operation IDs are opaque and are not
// DNS-label validated.
func (h *OperationHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
	op, err := h.registry.GetOperation(r.Context(), id)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, resources.ErrCodeResourceNotFound, "operation not found", "", "")
			return
		}
		writeError(w, r, http.StatusInternalServerError, resources.ErrCodeInternalError, "internal error", "", "")
		return
	}
	writeJSON(w, r, http.StatusOK, op)
}
