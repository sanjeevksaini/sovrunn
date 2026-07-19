package server

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/api"
	"github.com/sanjeevksaini/sovrunn/internal/config"
	"github.com/sanjeevksaini/sovrunn/internal/health"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Server owns HTTP server lifecycle.
type Server struct {
	cfg        config.Config
	httpServer *http.Server
	readiness  *health.ReadinessState
	logger     *log.Logger
}

// New constructs the Server with routes and middleware registered.
// ServiceClass, ServicePlan, Plugin, Capability, ServiceInstance, and
// ServiceBinding routes are registered only when the corresponding handlers
// are non-nil.
func New(
	cfg config.Config,
	org *api.OrgHandler,
	ou *api.OUHandler,
	tenant *api.TenantHandler,
	project *api.ProjectHandler,
	operation *api.OperationHandler,
	serviceClass *api.ServiceClassHandler,
	servicePlan *api.ServicePlanHandler,
	plugin *api.PluginHandler,
	capability *api.CapabilityHandler,
	serviceInstance *api.ServiceInstanceHandler,
	serviceBinding *api.ServiceBindingHandler,
	bootstrap *api.BootstrapHandler,
	readiness *health.ReadinessState,
) *Server {
	mux := http.NewServeMux()
	logger := log.New(os.Stdout, "", log.LstdFlags)

	chain := func(h http.Handler) http.Handler {
		return requestIDMiddleware(loggingMiddleware(logger)(contentTypeMiddleware(h)))
	}
	bootstrapChain := func(h http.Handler) http.Handler {
		return requestIDMiddleware(loggingMiddleware(logger)(methodGET(h)))
	}

	mux.Handle("/v1/organizations", chain(http.HandlerFunc(org.HandleCollection)))
	mux.Handle("/v1/organizations/", chain(http.HandlerFunc(org.HandleItem)))

	mux.Handle("/v1/organization-units", chain(http.HandlerFunc(ou.HandleCollection)))
	mux.Handle("/v1/organization-units/", chain(http.HandlerFunc(ou.HandleItem)))

	mux.Handle("/v1/tenants", chain(http.HandlerFunc(tenant.HandleCollection)))
	mux.Handle("/v1/tenants/", chain(http.HandlerFunc(tenant.HandleItem)))

	mux.Handle("/v1/projects", chain(http.HandlerFunc(project.HandleCollection)))
	mux.Handle("/v1/projects/", chain(http.HandlerFunc(project.HandleItem)))

	mux.Handle("/v1/operations", chain(http.HandlerFunc(operation.HandleCollection)))
	mux.Handle("/v1/operations/", chain(http.HandlerFunc(operation.HandleItem)))

	if serviceClass != nil {
		mux.Handle("/v1/service-classes", chain(http.HandlerFunc(serviceClass.HandleCollection)))
		mux.Handle("/v1/service-classes/", chain(http.HandlerFunc(serviceClass.HandleItem)))
	}
	if servicePlan != nil {
		mux.Handle("/v1/service-plans", chain(http.HandlerFunc(servicePlan.HandleCollection)))
		mux.Handle("/v1/service-plans/", chain(http.HandlerFunc(servicePlan.HandleItem)))
	}
	if plugin != nil {
		mux.Handle("/v1/plugins", chain(http.HandlerFunc(plugin.HandleCollection)))
		mux.Handle("/v1/plugins/", chain(http.HandlerFunc(plugin.HandleItem)))
	}
	if capability != nil {
		mux.Handle("/v1/capabilities", chain(http.HandlerFunc(capability.HandleCollection)))
		mux.Handle("/v1/capabilities/", chain(http.HandlerFunc(capability.HandleItem)))
	}
	if serviceInstance != nil {
		mux.Handle("/v1/service-instances", chain(http.HandlerFunc(serviceInstance.HandleCollection)))
		mux.Handle("/v1/service-instances/", chain(http.HandlerFunc(serviceInstance.HandleItem)))
	}
	if serviceBinding != nil {
		mux.Handle("/v1/service-bindings", chain(http.HandlerFunc(serviceBinding.HandleCollection)))
		mux.Handle("/v1/service-bindings/", chain(http.HandlerFunc(serviceBinding.HandleItem)))
	}

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		bootstrapChain(http.HandlerFunc(bootstrap.Healthz)).ServeHTTP(w, r)
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		bootstrapChain(http.HandlerFunc(bootstrap.Readyz)).ServeHTTP(w, r)
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		bootstrapChain(http.HandlerFunc(bootstrap.Version)).ServeHTTP(w, r)
	})

	return &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:              cfg.Addr(),
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		readiness: readiness,
		logger:    logger,
	}
}

// Start binds the listener, marks readiness true, and blocks until signal.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}

	s.readiness.SetReady(true)

	errCh := make(chan error, 1)
	go func() {
		s.logger.Printf("server listening on %s", s.cfg.Addr())
		if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-sigCtx.Done():
		timeout := time.Duration(s.cfg.Server.ShutdownTimeout) * time.Second
		if err := s.Shutdown(timeout); err != nil {
			return err
		}
		s.logger.Println("server shutdown complete")
		return nil
	case err := <-errCh:
		s.readiness.SetReady(false)
		return err
	}
}

// Shutdown stops accepting new connections and drains in-flight requests.
func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	s.readiness.SetShuttingDown()
	return s.httpServer.Shutdown(ctx)
}

func writeErrorBody(w http.ResponseWriter, code resources.ErrorCode, message, field, details string) error {
	envelope := resources.APIErrorEnvelope{
		Error: resources.APIError{
			Code:    code,
			Message: message,
			Field:   field,
			Details: details,
		},
	}
	return json.NewEncoder(w).Encode(envelope)
}
