package api

import (
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/config"
	"github.com/sanjeevksaini/sovrunn/internal/health"
)

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
	writeJSON(w, r, http.StatusOK, map[string]string{"status": "ok"})
}

// Readyz returns readiness status.
func (h *BootstrapHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	if h.readiness.IsReady() {
		writeJSON(w, r, http.StatusOK, map[string]string{"status": "ready"})
		return
	}
	writeJSON(w, r, http.StatusServiceUnavailable, map[string]string{"status": "not-ready"})
}

// buildVersion is set via -ldflags at build time; defaults to "dev".
var buildVersion = "dev"

// Version returns build information.
func (h *BootstrapHandler) Version(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, http.StatusOK, map[string]string{
		"name":    "sovrunn-api",
		"version": buildVersion,
		"phase":   "1",
		"status":  "alpha",
	})
}
