package api

import (
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/config"
	"github.com/sanjeevksaini/sovrunn/internal/health"
)

const (
	versionPhase  = "1"
	versionStatus = "alpha"
)

type healthResponse struct {
	Status string `json:"status"`
}

type readyResponse struct {
	Status string `json:"status"`
}

type notReadyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type versionResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Phase   string `json:"phase"`
	Status  string `json:"status"`
}

// BootstrapHandler serves health, readiness, and version endpoints.
type BootstrapHandler struct {
	cfg       config.Config
	readiness *health.ReadinessState
}

// NewBootstrapHandler constructs a BootstrapHandler.
func NewBootstrapHandler(cfg config.Config, r *health.ReadinessState) *BootstrapHandler {
	return &BootstrapHandler{cfg: cfg, readiness: r}
}

// Healthz returns liveness status.
func (h *BootstrapHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, http.StatusOK, healthResponse{Status: "ok"})
}

// Readyz returns readiness status.
func (h *BootstrapHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	if h.readiness.IsReady() {
		writeJSON(w, r, http.StatusOK, readyResponse{Status: "ready"})
		return
	}
	writeJSON(w, r, http.StatusServiceUnavailable, notReadyResponse{
		Status:  "not_ready",
		Message: h.readiness.Reason(),
	})
}

// buildVersion is set via -ldflags at build time; defaults to "dev".
var buildVersion = "dev"

// Version returns build information.
func (h *BootstrapHandler) Version(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, http.StatusOK, versionResponse{
		Name:    "sovrunn-api",
		Version: buildVersion,
		Phase:   versionPhase,
		Status:  versionStatus,
	})
}
